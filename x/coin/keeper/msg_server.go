package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bytes"
	"context"
	sdkmath "cosmossdk.io/math"
	"encoding/base64"
	"math/big"
	"strconv"
	"strings"

	"github.com/cosmos/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"

	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

var _ types.MsgServer = &Keeper{}

////////////////////////////////////////////////////////////////
// CreateCoin
////////////////////////////////////////////////////////////////

func (k Keeper) CreateCoin(goCtx context.Context, msg *types.MsgCreateCoin) (*types.MsgCreateCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	baseCoinDenom := k.GetBaseDenom(ctx)
	coinDenom := strings.ToLower(msg.Symbol)

	// Create new coin instance
	var coin = types.Coin{
		Title:       msg.Title,
		Symbol:      coinDenom,
		CRR:         msg.CRR,
		Reserve:     msg.InitialReserve,
		Volume:      msg.InitialVolume,
		LimitVolume: msg.LimitVolume,
		Creator:     msg.Sender,
		Identity:    msg.Identity,
	}

	// Ensure coin does not exist
	_, err := k.GetCoin(ctx, coinDenom)
	if err == nil {
		return nil, errors.CoinAlreadyExists
	}

	// Calculate special fee for creating custom coin
	feeAmountBase := helpers.EtherToWei(getCreateCoinCommission(coinDenom))
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
	err = ctx.EventManager().EmitTypedEvent(&types.EventCreateCoin{
		Sender:               sender.String(),
		Symbol:               coinDenom,
		Title:                msg.Title,
		Crr:                  msg.CRR,
		InitialVolume:        msg.InitialVolume.String(),
		InitialReserve:       msg.InitialReserve.String(),
		LimitVolume:          msg.LimitVolume.String(),
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
	coinDenom := strings.ToLower(msg.Symbol)

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
	if coin.LimitVolume.GT(msg.LimitVolume) {
		return nil, errors.NewLimitVolumeLess
	}

	// Update coin metadata
	coin.LimitVolume = msg.LimitVolume
	coin.Identity = msg.Identity

	// Save coin to the storage
	k.SetCoin(ctx, coin)

	// Emit transaction events
	err = ctx.EventManager().EmitTypedEvent(&types.EventUpdateCoin{
		Sender:      msg.Sender,
		Symbol:      coin.Symbol,
		LimitVolume: msg.LimitVolume.String(),
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
	receiver, _ := sdk.AccAddressFromBech32(msg.Receiver)

	// Send coins from the sender to the recipient
	err = k.bankKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(msg.Coin))
	if err != nil {
		return nil, err
	}

	// Emit transaction events
	err = ctx.EventManager().EmitTypedEvent(&types.EventSendCoin{
		Sender:   msg.Sender,
		Receiver: msg.Receiver,
		Coin:     msg.Coin.String(),
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
		receiver, _ := sdk.AccAddressFromBech32(msg.Sends[i].Receiver)

		// Retrieve sending coin
		_, err := k.GetCoin(ctx, coinDenom)
		if err != nil {
			return nil, err
		}

		// Send coins from the sender to the recipient
		err = k.bankKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(msg.Sends[i].Coin))
		if err != nil {
			return nil, err
		}

		// Emit transaction events
		err = ctx.EventManager().EmitTypedEvent(&types.EventSendCoin{
			Sender:   msg.Sender,
			Receiver: msg.Sends[i].Receiver,
			Coin:     msg.Sends[i].Coin.String(),
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
	err := k.sellCoin(ctx, sender, sdk.NewCoin(msg.CoinSymbolToSell, sdk.ZeroInt()), msg.MinCoinToBuy, true)
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
		k.EditCoin(ctx, coin, coin.Reserve, coin.Volume.Sub(msg.Coin.Amount))
	}

	// Emit transaction events
	err = ctx.EventManager().EmitTypedEvent(&types.EventBurnCoin{
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
	coinDenom := strings.ToLower(check.Coin)
	coinAmount := check.Amount

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
	feeAmountBase := helpers.FinneyToWei(sdk.NewIntFromUint64(30))
	feeAmount := feeAmountBase
	if coinDenom != baseCoinDenom {
		feeAmount = formulas.CalculateSaleAmount(coin.Volume, coin.Reserve, uint(coin.CRR), feeAmountBase)
	}
	feeCoin := sdk.NewCoin(coinDenom, feeAmount)

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
		return nil, errors.UnableRPLEncodeAddress

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

	// Send fee from issuer to the module
	// TODO: Make sure it is correct way to get fees
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, issuer, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return nil, err

	}

	// Send check coins from issuer to the transaction sender
	err = k.bankKeeper.SendCoins(ctx, issuer, sender, sdk.NewCoins(sdk.NewCoin(coinDenom, coinAmount)))
	if err != nil {
		return nil, err
	}

	// Emit transaction events
	err = ctx.EventManager().EmitTypedEvent(&types.EventRedeemCheck{
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

func (k Keeper) buyCoin(
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
	amountToSell := sdk.ZeroInt()
	amountInBaseCoin := sdk.ZeroInt()
	switch {
	case k.IsCoinBase(ctx, coinToSell.Symbol):
		// Buyer buys custom coin for base coin
		amountToSell = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, uint(coinToBuy.CRR), amountToBuy)
		amountInBaseCoin = amountToSell
	case k.IsCoinBase(ctx, coinToBuy.Symbol):
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
	if !k.IsCoinBase(ctx, coinToSell.Symbol) {
		if coinToSell.Reserve.Sub(amountInBaseCoin).LT(types.MinCoinReserve) {
			return errors.TxBreaksMinReserveRule
		}
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
	if !k.IsCoinBase(ctx, coinToSell.Symbol) {
		err = k.EditCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountInBaseCoin), coinToSell.Volume.Sub(amountToSell))
		if err != nil {
			return err
		}
	}
	if !k.IsCoinBase(ctx, coinToBuy.Symbol) {
		err = k.EditCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountInBaseCoin), coinToBuy.Volume.Add(amountToBuy))
		if err != nil {
			return err
		}
	}

	// Emit transaction events
	err = ctx.EventManager().EmitTypedEvent(&types.EventBuySellCoin{
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

func (k Keeper) sellCoin(
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
	amountToSell, amountToBuy, amountInBaseCoin := coin.Amount, sdk.ZeroInt(), sdk.ZeroInt()
	switch {
	case k.IsCoinBase(ctx, coinToBuy.Symbol):
		// Seller sells custom coin for base coin
		amountToBuy = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, uint(coinToSell.CRR), amountToSell)
		amountInBaseCoin = amountToBuy
	case k.IsCoinBase(ctx, coinToSell.Symbol):
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
	if !k.IsCoinBase(ctx, coinToSell.Symbol) {
		if coinToSell.Reserve.Sub(amountInBaseCoin).LT(types.MinCoinReserve) {
			return errors.TxBreaksMinReserveRule
		}
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !k.IsCoinBase(ctx, coinToBuy.Symbol) {
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
	if !k.IsCoinBase(ctx, coinToSell.Symbol) {
		err = k.EditCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountInBaseCoin), coinToSell.Volume.Sub(amountToSell))
		if err != nil {
			return err
		}
	}
	if !k.IsCoinBase(ctx, coinToBuy.Symbol) {
		err = k.EditCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountInBaseCoin), coinToBuy.Volume.Add(amountToBuy))
		if err != nil {
			return err
		}
	}

	// Emit transaction events
	err = ctx.EventManager().EmitTypedEvent(&types.EventBuySellCoin{
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

func getCreateCoinCommission(symbol string) sdkmath.Int {
	switch len(symbol) {
	case 3:
		return sdk.NewInt(1_000_000)
	case 4:
		return sdk.NewInt(100_000)
	case 5:
		return sdk.NewInt(10_000)
	case 6:
		return sdk.NewInt(1000)
	}
	return sdk.NewInt(100)
}
