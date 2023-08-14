package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/precompile/drc20cosmos"
	"bytes"
	"context"
	"encoding/base64"
	"math/big"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	dscconfig "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/config"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feeconfig "bitbucket.org/decimalteam/go-smart-node/x/fee/config"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

var _ types.MsgServer = &Keeper{}

////////////////////////////////////////////////////////////////
// CreateCoin
////////////////////////////////////////////////////////////////

func (k Keeper) CreateCoin(goCtx context.Context, msg *types.MsgCreateCoin) (*types.MsgCreateCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	baseCoinDenom := k.GetBaseDenom(ctx)
	coinDenom := strings.ToLower(msg.Denom)

	// Create new coin instance
	var coin = types.Coin{
		Title:        msg.Title,
		Denom:        coinDenom,
		CRR:          msg.CRR,
		Reserve:      msg.InitialReserve,
		Volume:       msg.InitialVolume,
		LimitVolume:  msg.LimitVolume,
		MinVolume:    msg.MinVolume,
		Creator:      msg.Sender,
		Identity:     msg.Identity,
		Drc20Address: "",
	}

	// Ensure coin does not exist
	_, err := k.GetCoin(ctx, coinDenom)
	if err == nil {
		return nil, errors.CoinAlreadyExists
	}

	drc20, err := drc20cosmos.NewDrc20Cosmos(ctx, k.evm, k.bankKeeper, coin)
	if err != nil {
		ctx.Logger().Info(err.Error())
	}

	_, err = drc20.CreateContractIfNotSet()
	if err != nil {
		ctx.Logger().Info(err.Error())
	}

	coin = drc20.Coin

	// Calculate special fee for creating custom coin
	feeAmountBase, err := k.getCreateCoinCommission(ctx, coinDenom)
	if err != nil {
		return nil, err
	}

	feeAmount, feeDenom, err := k.GetCommission(ctx, feeAmountBase)
	if err != nil {
		return nil, err
	}
	feeCoin := sdk.NewCoin(feeDenom, feeAmount)

	// NOTE: It was already validated so no need to check error
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	// Retrieve sender's balances
	balanceBaseCoin := k.bankKeeper.GetBalance(ctx, sender, baseCoinDenom)
	balanceFeeCoin := balanceBaseCoin
	if feeDenom != baseCoinDenom {
		balanceFeeCoin = k.bankKeeper.GetBalance(ctx, sender, feeDenom)
	}

	// Ensure balances are enough
	if feeDenom == baseCoinDenom {
		feeAmountBaseTotal := feeAmount.Add(msg.InitialReserve)
		if balanceBaseCoin.Amount.LT(feeAmountBaseTotal) {
			return nil, errors.InsufficientFunds
		}
	} else {
		if balanceBaseCoin.Amount.LT(msg.InitialReserve) {
			return nil, errors.InsufficientFunds
		}
		if balanceFeeCoin.Amount.LT(feeAmount) {
			return nil, errors.InsufficientFunds
		}
	}

	// Send initial reserve to the module
	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, sender, types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(baseCoinDenom, msg.InitialReserve)),
	)
	if err != nil {
		return nil, err
	}

	// Send special fee to the module
	// TODO: Make sure it is correct way to get fees
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return nil, err
	}

	// Mint initial coins to the module and send to the coin creator
	initialCoins := sdk.NewCoins(sdk.NewCoin(coinDenom, msg.InitialVolume))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, initialCoins)
	if err != nil {
		return nil, err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, initialCoins)
	if err != nil {
		return nil, err
	}

	// Save coin to the storage
	k.SetCoin(ctx, coin)

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventCreateCoin{
		Sender:               sender.String(),
		Denom:                coinDenom,
		Title:                msg.Title,
		CRR:                  msg.CRR,
		InitialVolume:        msg.InitialVolume.String(),
		InitialReserve:       msg.InitialReserve.String(),
		LimitVolume:          msg.LimitVolume.String(),
		MinVolume:            msg.MinVolume.String(),
		Identity:             msg.Identity,
		CommissionCreateCoin: feeCoin.String(),
	})

	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgCreateCoinResponse{}, nil
}

////////////////////////////////////////////////////////////////
// UpdateCoin
////////////////////////////////////////////////////////////////

func (k Keeper) UpdateCoin(goCtx context.Context, msg *types.MsgUpdateCoin) (*types.MsgUpdateCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	coinDenom := strings.ToLower(msg.Denom)

	// Retrieve updating coin
	coin, err := k.GetCoin(ctx, coinDenom)
	if err != nil {
		return nil, err
	}

	// Ensure sender is the coin creator
	if strings.Compare(coin.Creator, msg.Sender) != 0 {
		return nil, errors.UpdateOnlyForCreator
	}

	// Ensure new limit volume is big enough
	if coin.Volume.GT(msg.LimitVolume) {
		return nil, errors.NewLimitVolumeLess
	}

	// Validate min emission if specified
	if coin.MinVolume.IsZero() != (msg.MinVolume.IsNil() || msg.MinVolume.IsZero()) {
		return nil, errors.UneditableCoinMinEmission
	}

	// Update coin metadata
	coin.LimitVolume = msg.LimitVolume
	coin.MinVolume = msg.MinVolume
	coin.Identity = msg.Identity

	// Save coin to the storage
	k.SetCoin(ctx, coin)

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventUpdateCoin{
		Sender:      msg.Sender,
		Denom:       coin.Denom,
		LimitVolume: msg.LimitVolume.String(),
		MinVolume:   msg.MinVolume.String(),
		Identity:    msg.Identity,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgUpdateCoinResponse{}, nil
}

////////////////////////////////////////////////////////////////
// SendCoin
////////////////////////////////////////////////////////////////

// SendCoin creates new custom coin.
func (k Keeper) SendCoin(goCtx context.Context, msg *types.MsgSendCoin) (*types.MsgSendCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	coinDenom := strings.ToLower(msg.Coin.Denom)

	// Retrieve sending coin
	_, err := k.GetCoin(ctx, coinDenom)
	if err != nil {
		return nil, err
	}

	// NOTE: It was already validated so no need to check error
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	recipient, _ := sdk.AccAddressFromBech32(msg.Recipient)

	// Send coins from the sender to the recipient
	err = k.bankKeeper.SendCoins(ctx, sender, recipient, sdk.NewCoins(msg.Coin))
	if err != nil {
		return nil, err
	}

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventSendCoin{
		Sender:    msg.Sender,
		Recipient: msg.Recipient,
		Coin:      msg.Coin.String(),
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgSendCoinResponse{}, nil
}

////////////////////////////////////////////////////////////////
// MultiSendCoin
////////////////////////////////////////////////////////////////

func (k Keeper) MultiSendCoin(goCtx context.Context, msg *types.MsgMultiSendCoin) (*types.MsgMultiSendCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// NOTE: It was already validated so no need to check error
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	for i := range msg.Sends {
		coinDenom := strings.ToLower(msg.Sends[i].Coin.Denom)

		// NOTE: It was already validated so no need to check error
		recipient, _ := sdk.AccAddressFromBech32(msg.Sends[i].Recipient)

		// Retrieve sending coin
		_, err := k.GetCoin(ctx, coinDenom)
		if err != nil {
			return nil, err
		}

		// Send coins from the sender to the recipient
		err = k.bankKeeper.SendCoins(ctx, sender, recipient, sdk.NewCoins(msg.Sends[i].Coin))
		if err != nil {
			return nil, err
		}

		// Emit transaction events
		err = events.EmitTypedEvent(ctx, &types.EventSendCoin{
			Sender:    msg.Sender,
			Recipient: msg.Sends[i].Recipient,
			Coin:      msg.Sends[i].Coin.String(),
		})
		if err != nil {
			return nil, errors.Internal.Wrapf("event err: %s", err.Error())
		}
	}

	return &types.MsgMultiSendCoinResponse{}, nil
}

////////////////////////////////////////////////////////////////
// BuyCoin
////////////////////////////////////////////////////////////////

func (k Keeper) BuyCoin(goCtx context.Context, msg *types.MsgBuyCoin) (*types.MsgBuyCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// NOTE: It was already validated so no need to check error
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	// Make buy
	err := k.buyCoin(ctx, sender, msg.CoinToBuy, msg.MaxCoinToSell)
	if err != nil {
		return nil, err
	}

	return &types.MsgBuyCoinResponse{}, nil
}

////////////////////////////////////////////////////////////////
// SellCoin
////////////////////////////////////////////////////////////////

func (k Keeper) SellCoin(goCtx context.Context, msg *types.MsgSellCoin) (*types.MsgSellCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// NOTE: It was already validated so no need to check error
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	// Make sale
	err := k.sellCoin(ctx, sender, msg.CoinToSell, msg.MinCoinToBuy, false)
	if err != nil {
		return nil, err
	}

	return &types.MsgSellCoinResponse{}, nil
}

////////////////////////////////////////////////////////////////
// SellAllCoin
////////////////////////////////////////////////////////////////

func (k Keeper) SellAllCoin(goCtx context.Context, msg *types.MsgSellAllCoin) (*types.MsgSellAllCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// NOTE: It was already validated so no need to check error
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	// Make sale
	err := k.sellCoin(ctx, sender, sdk.NewCoin(msg.CoinDenomToSell, sdkmath.ZeroInt()), msg.MinCoinToBuy, true)
	if err != nil {
		return nil, err
	}

	return &types.MsgSellAllCoinResponse{}, nil
}

////////////////////////////////////////////////////////////////
// BurnCoin
////////////////////////////////////////////////////////////////

func (k Keeper) BurnCoin(goCtx context.Context, msg *types.MsgBurnCoin) (*types.MsgBurnCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// NOTE: It was already validated so no need to check error
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	coin, err := k.GetCoin(ctx, msg.Coin.Denom)
	if err != nil {
		return nil, err
	}
	if !k.IsCoinBase(ctx, msg.Coin.Denom) {
		// check for limits
		err = k.CheckFutureVolumeChanges(ctx, coin, msg.Coin.Amount.Neg())
		if err != nil {
			return nil, err
		}
	}
	// send to coin module and burn
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(msg.Coin))
	if err != nil {
		return nil, err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(msg.Coin))
	if err != nil {
		return nil, err
	}
	if !k.IsCoinBase(ctx, msg.Coin.Denom) {
		// change coin volume
		err = k.UpdateCoinVR(ctx, coin.Denom, coin.Volume.Sub(msg.Coin.Amount), coin.Reserve)
		if err != nil {
			return nil, err
		}
	}

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventBurnCoin{
		Sender: msg.Sender,
		Coin:   msg.Coin.String(),
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgBurnCoinResponse{}, nil
}

////////////////////////////////////////////////////////////////
// RedeemCheck
////////////////////////////////////////////////////////////////

func (k Keeper) RedeemCheck(goCtx context.Context, msg *types.MsgRedeemCheck) (*types.MsgRedeemCheckResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	baseCoinDenom := k.GetBaseDenom(ctx)

	// NOTE: It was already validated so no need to check error
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	// Decode provided check from base58 format to raw bytes
	checkBytes := base58.Decode(msg.Check)
	if len(checkBytes) == 0 {
		return nil, errors.UnableDecodeCheckBase58
	}

	// Parse provided check from raw bytes to ensure it is valid
	check, err := types.ParseCheck(checkBytes)
	if err != nil {
		return nil, err
	}
	coinDenom := strings.ToLower(check.Coin.Denom)
	coinAmount := check.Coin.Amount

	// Decode provided proof from base64 format to raw bytes
	proof, err := base64.StdEncoding.DecodeString(msg.Proof)
	if err != nil {
		return nil, errors.UnableDecodeProofBase64
	}

	// Recover issuer address from check signature
	issuer, err := check.Sender()
	if err != nil {
		return nil, errors.UnableRecoverAddressFromCheck
	}

	// Retrieve the coin specified in the check
	coin, err := k.GetCoin(ctx, coinDenom)
	if err != nil {
		return nil, err
	}

	// Retrieve issuer's balance of issuing coins
	balance := k.bankKeeper.GetBalance(ctx, issuer, coinDenom)
	balanceFeeCoin := balance
	if coinDenom != baseCoinDenom {
		balanceFeeCoin = k.bankKeeper.GetBalance(ctx, issuer, baseCoinDenom)
	}

	// Calculate correct fee
	params := k.feeKeeper.GetModuleParams(ctx)
	delPrice, err := k.feeKeeper.GetPrice(ctx, k.GetBaseDenom(ctx), feeconfig.DefaultQuote)
	if err != nil {
		return nil, err
	}
	feeAmountBase := helpers.DecToDecWithE18(params.CoinRedeemCheck).Quo(delPrice.Price).RoundInt()
	feeAmount := feeAmountBase
	if coinDenom != baseCoinDenom {
		feeAmount = formulas.CalculateSaleAmount(coin.Volume, coin.Reserve, uint(coin.CRR), feeAmountBase)
	}
	feeCoin := sdk.NewCoin(coinDenom, feeAmount)
	// split to burning and collected part
	amountToBurn := sdk.NewDecFromInt(feeAmount).Mul(params.CommissionBurnFactor).RoundInt()
	amountToCollect := feeAmount.Sub(amountToBurn)

	// check case when fee pay will break reserve/volume limits
	err = k.CheckFutureChanges(ctx, coin, feeCoin.Amount.Neg())
	if err != nil {
		return nil, err
	}

	// Ensure that check issuer account holds enough coins
	if balance.Amount.LT(coinAmount) {
		return nil, errors.InsufficientFunds
	}
	// Ensure the check issuer account holds enough coins for fee pay
	if coinDenom != baseCoinDenom {
		if balanceFeeCoin.Amount.LT(feeAmount) {
			return nil, errors.InsufficientFunds
		}
	} else {
		if balance.Amount.LT(coinAmount.Add(feeAmount)) {
			return nil, errors.InsufficientFunds
		}
	}

	// Ensure the proper chain ID is specified in the check
	if check.ChainID != ctx.ChainID() {
		return nil, errors.InvalidChainID
	}

	// Ensure nonce length
	if len(check.Nonce) > 16 {
		return nil, errors.InvalidNonce
	}

	// Check block number
	if check.DueBlock < uint64(ctx.BlockHeight()) {
		return nil, errors.CheckExpired
	}

	// Ensure check is not redeemed yet
	if k.IsCheckRedeemed(ctx, check) {
		return nil, errors.CheckRedeemed
	}

	// Recover public key from check lock
	publicKeyA, err := check.LockPubKey()
	if err != nil {
		return nil, err
	}

	// Prepare bytes used to recover public key from provided proof
	senderAddressHash := make([]byte, 32)
	hw := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hw, []interface{}{sender})
	if err != nil {
		return nil, errors.UnableRLPEncodeAddress

	}
	hw.Sum(senderAddressHash[:0])

	// Recover public key from provided proof
	publicKeyB, err := crypto.Ecrecover(senderAddressHash[:], proof)
	if err != nil {
		return nil, errors.FailedToRecoverPKFromSig
	}

	// Compare both public keys to ensure provided proof is correct
	if !bytes.Equal(publicKeyA, publicKeyB) {
		return nil, errors.InvalidProof
	}

	// Write check to the storage
	k.SetCheck(ctx, check)

	// Send fee from issuer to the fee_collector
	// send to burn
	if amountToBurn.IsPositive() {
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, issuer, feetypes.BurningPool,
			sdk.NewCoins(sdk.NewCoin(coinDenom, amountToBurn)))
		if err != nil {
			return nil, err
		}
	}

	// send to collect
	if amountToCollect.IsPositive() {
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, issuer, sdkAuthTypes.FeeCollectorName,
			sdk.NewCoins(sdk.NewCoin(coinDenom, amountToCollect)))
		if err != nil {
			return nil, err
		}
	}

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &feetypes.EventPayCommission{
		Payer: issuer.String(),
		Coins: sdk.NewCoins(feeCoin),
		Burnt: sdk.NewCoins(sdk.NewCoin(coinDenom, amountToBurn)),
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	// Send check coins from issuer to the transaction sender
	err = k.bankKeeper.SendCoins(ctx, issuer, sender, sdk.NewCoins(sdk.NewCoin(coinDenom, coinAmount)))
	if err != nil {
		return nil, err
	}

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventRedeemCheck{
		Sender:                msg.Sender,
		Issuer:                issuer.String(),
		Coin:                  sdk.NewCoin(coinDenom, coinAmount).String(),
		Nonce:                 new(big.Int).SetBytes(check.Nonce).String(),
		DueBlock:              strconv.FormatUint(check.DueBlock, 10),
		CommissionRedeemCheck: feeCoin.String(),
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgRedeemCheckResponse{}, nil
}

////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////

func (k *Keeper) buyCoin(
	ctx sdk.Context,
	sender sdk.AccAddress,
	coin sdk.Coin,
	maxCoinToSell sdk.Coin,
) error {
	coinToBuyDenom := strings.ToLower(coin.Denom)
	coinToSellDenom := strings.ToLower(maxCoinToSell.Denom)

	// Retrieve buyer's balance of selling coins
	balance := k.bankKeeper.GetBalance(ctx, sender, coinToSellDenom)

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, coinToBuyDenom)
	if err != nil {
		return err
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, coinToSellDenom)
	if err != nil {
		return err
	}

	// Ensure supply limit of the coin to buy does not overflow
	err = k.CheckFutureChanges(ctx, coinToBuy, coin.Amount)
	if err != nil {
		return err
	}

	// Calculate amount of sell coins which buyer will receive
	amountToBuy := coin.Amount
	var amountToSell, amountInBaseCoin sdkmath.Int
	switch {
	case k.IsCoinBase(ctx, coinToSell.Denom):
		// Buyer buys custom coin for base coin
		amountToSell = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, uint(coinToBuy.CRR), amountToBuy)
		amountInBaseCoin = amountToSell
	case k.IsCoinBase(ctx, coinToBuy.Denom):
		// Buyer buys base coin for custom coin
		if coin.Amount.GT(coinToSell.Reserve) {
			return errors.InsufficientCoinReserve
		}
		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, uint(coinToSell.CRR), amountToBuy)
		amountInBaseCoin = amountToBuy
	default:
		// Buyer buys custom coin for custom coin
		amountInBaseCoin = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, uint(coinToBuy.CRR), amountToBuy)
		if amountInBaseCoin.GT(coinToSell.Reserve) {
			return errors.InsufficientCoinReserve
		}
		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, uint(coinToSell.CRR), amountInBaseCoin)
	}

	// Ensure maximum amount of coins to sell (price guard)
	if amountToSell.GT(maxCoinToSell.Amount) {
		return errors.MaximumValueToSellReached
	}

	// Ensure reserve of the coin to sell does not underflow
	err = k.CheckFutureChanges(ctx, coinToSell, amountToSell.Neg())
	if err != nil {
		return err
	}

	// Ensure that buyer account holds enough coins to sell
	if balance.Amount.LT(amountToSell) {
		return errors.InsufficientFunds
	}

	coinsToSell := sdk.NewCoins(sdk.NewCoin(coinToSellDenom, amountToSell))
	coinsToBuy := sdk.NewCoins(sdk.NewCoin(coinToBuyDenom, amountToBuy))

	// Send sold coins from the buyer to the module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coinsToSell)
	if err != nil {
		return err
	}

	// Burn sold coins from the module
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coinsToSell)
	if err != nil {
		return err
	}

	// Mint bought coins to the module
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coinsToBuy)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coinsToBuy)
	if err != nil {
		return err
	}

	// Update coins
	if !k.IsCoinBase(ctx, coinToSell.Denom) {
		err = k.UpdateCoinVR(ctx, coinToSell.Denom, coinToSell.Volume.Sub(amountToSell), coinToSell.Reserve.Sub(amountInBaseCoin))
		if err != nil {
			return err
		}
	}
	if !k.IsCoinBase(ctx, coinToBuy.Denom) {
		err = k.UpdateCoinVR(ctx, coinToBuy.Denom, coinToBuy.Volume.Add(amountToBuy), coinToBuy.Reserve.Add(amountInBaseCoin))
		if err != nil {
			return err
		}
	}

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventBuySellCoin{
		Sender:           sender.String(),
		CoinToBuy:        sdk.NewCoin(coinToBuyDenom, amountToBuy).String(),
		CoinToSell:       sdk.NewCoin(coinToSellDenom, amountToSell).String(),
		AmountInBaseCoin: amountInBaseCoin.String(),
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	return nil
}

func (k *Keeper) sellCoin(
	ctx sdk.Context,
	sender sdk.AccAddress,
	coin sdk.Coin,
	minCoinToBuy sdk.Coin,
	sellAll bool,
) error {
	coinToSellDenom := strings.ToLower(coin.Denom)
	coinToBuyDenom := strings.ToLower(minCoinToBuy.Denom)

	// Retrieve seller's balance of selling coins
	balance := k.bankKeeper.GetBalance(ctx, sender, coinToSellDenom)

	// Fill amount to sell in case of MsgSellAll
	if sellAll {
		coin.Amount = balance.Amount
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, coinToSellDenom)
	if err != nil {
		return err
	}

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, coinToBuyDenom)
	if err != nil {
		return err
	}

	// Ensure that seller account holds enough coins to sell
	if balance.Amount.LT(coin.Amount) {
		return errors.InsufficientFunds
	}

	err = k.CheckFutureChanges(ctx, coinToSell, coin.Amount.Neg())
	if err != nil {
		return err
	}

	// Calculate amount of buy coins which seller will receive
	amountToSell := coin.Amount
	var amountToBuy, amountInBaseCoin sdkmath.Int
	switch {
	case k.IsCoinBase(ctx, coinToBuy.Denom):
		// Seller sells custom coin for base coin
		amountToBuy = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, uint(coinToSell.CRR), amountToSell)
		amountInBaseCoin = amountToBuy
	case k.IsCoinBase(ctx, coinToSell.Denom):
		// Seller sells base coin for custom coin
		amountToBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, uint(coinToBuy.CRR), amountToSell)
		amountInBaseCoin = amountToSell
	default:
		// Seller sells custom coin for custom coin
		amountInBaseCoin = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, uint(coinToSell.CRR), amountToSell)
		amountToBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, uint(coinToBuy.CRR), amountInBaseCoin)
	}

	// Ensure minimum amount of coins to buy (price guard)
	if amountToBuy.LT(minCoinToBuy.Amount) {
		return errors.MinimumValueToBuyReached
	}

	// Ensure reserve of the coin to sell does not underflow
	if !k.IsCoinBase(ctx, coinToSell.Denom) {
		if coinToSell.Reserve.Sub(amountInBaseCoin).LT(config.MinCoinReserve) {
			return errors.TxBreaksMinReserveRule
		}
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !k.IsCoinBase(ctx, coinToBuy.Denom) {
		if coinToBuy.Volume.Add(amountToBuy).GT(coinToBuy.LimitVolume) {
			return errors.TxBreaksVolumeLimit
		}
	}

	coinsToSell := sdk.NewCoins(sdk.NewCoin(coinToSellDenom, amountToSell))
	coinsToBuy := sdk.NewCoins(sdk.NewCoin(coinToBuyDenom, amountToBuy))

	// Send sold coins from the seller to the module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coinsToSell)
	if err != nil {
		return err
	}

	// Burn sold coins from the module
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coinsToSell)
	if err != nil {
		return err
	}

	// Mint bought coins to the module
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coinsToBuy)
	if err != nil {
		return err
	}

	// Send bought coins from the module to the seller
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coinsToBuy)
	if err != nil {
		return err
	}

	// Update coins
	if !k.IsCoinBase(ctx, coinToSell.Denom) {
		err = k.UpdateCoinVR(ctx, coinToSell.Denom, coinToSell.Volume.Sub(amountToSell), coinToSell.Reserve.Sub(amountInBaseCoin))
		if err != nil {
			return err
		}
	}
	if !k.IsCoinBase(ctx, coinToBuy.Denom) {
		err = k.UpdateCoinVR(ctx, coinToBuy.Denom, coinToBuy.Volume.Add(amountToBuy), coinToBuy.Reserve.Add(amountInBaseCoin))
		if err != nil {
			return err
		}
	}

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventBuySellCoin{
		Sender:           sender.String(),
		CoinToBuy:        sdk.NewCoin(coinToBuyDenom, amountToBuy).String(),
		CoinToSell:       sdk.NewCoin(coinToSellDenom, amountToSell).String(),
		AmountInBaseCoin: amountInBaseCoin.String(),
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	return nil
}

func (k Keeper) getCreateCoinCommission(ctx sdk.Context, symbol string) (sdkmath.Int, error) {
	baseDenomPrice, err := k.feeKeeper.GetPrice(ctx, dscconfig.BaseDenom, feeconfig.DefaultQuote)
	if err != nil {
		return sdkmath.Int{}, err
	}

	params := k.feeKeeper.GetModuleParams(ctx)

	var createCoinFee sdk.Dec
	switch len(symbol) {
	case 3:
		createCoinFee = params.CoinCreateTicker3
	case 4:
		createCoinFee = params.CoinCreateTicker4
	case 5:
		createCoinFee = params.CoinCreateTicker5
	case 6:
		createCoinFee = params.CoinCreateTicker6
	default:
		createCoinFee = params.CoinCreateTicker7
	}

	return helpers.DecToIntWithE18(createCoinFee.Quo(baseDenomPrice.Price)), nil
}
