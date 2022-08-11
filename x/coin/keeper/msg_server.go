package keeper

import (
	"bytes"
	"context"
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
		return nil, types.ErrCoinAlreadyExists(coinDenom)
	}

	// Calculate special fee for creating custom coin
	feeAmountBase := helpers.EtherToWei(getCreateCoinCommission(coinDenom))
	feeAmount, feeDenom, err := k.GetCommission(ctx, feeAmountBase)
	if err != nil {
		return nil, types.ErrCalculateCommission(err.Error())
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
			return nil, types.ErrInsufficientFunds(
				sdk.NewCoin(baseCoinDenom, feeAmountBaseTotal).String(),
				balanceBaseCoin.String(),
			)
		}
	} else {
		if balanceBaseCoin.Amount.LT(msg.InitialReserve) {
			return nil, types.ErrInsufficientCoinReserve()
		}
		if balanceFeeCoin.Amount.LT(feeAmount) {
			return nil, types.ErrInsufficientFundsToPayCommission(feeAmount.String())
		}
	}

	// Send initial reserve to the module
	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, sender, types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(baseCoinDenom, msg.InitialReserve)),
	)
	if err != nil {
		// TODO: Change error
		return nil, types.ErrUpdateBalance(sender.String(), err.Error())
	}

	// Send special fee to the module
	// TODO: Make sure it is correct way to get fees
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}

	// Mint initial coins to the module and send to the coin creator
	initialCoins := sdk.NewCoins(sdk.NewCoin(coinDenom, msg.InitialVolume))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, initialCoins)
	if err != nil {
		// TODO: Change error
		return nil, types.ErrUpdateBalance(sender.String(), err.Error())
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, initialCoins)
	if err != nil {
		// TODO: Change error
		return nil, types.ErrUpdateBalance(sender.String(), err.Error())
	}

	// Save coin to the storage
	k.SetCoin(ctx, coin)

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeSymbol, coinDenom),
		sdk.NewAttribute(types.AttributeTitle, msg.Title),
		sdk.NewAttribute(types.AttributeCRR, strconv.FormatUint(msg.CRR, 10)),
		sdk.NewAttribute(types.AttributeInitVolume, msg.InitialVolume.String()),
		sdk.NewAttribute(types.AttributeInitReserve, msg.InitialReserve.String()),
		sdk.NewAttribute(types.AttributeLimitVolume, msg.LimitVolume.String()),
		sdk.NewAttribute(types.AttributeIdentity, msg.Identity),
		sdk.NewAttribute(types.AttributeCommissionCreateCoin, feeCoin.String()),
	))

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
		return nil, types.ErrCoinDoesNotExist(coinDenom)
	}

	// Ensure sender is the coin creator
	if strings.Compare(coin.Creator, msg.Sender) != 0 {
		return nil, types.ErrUpdateOnlyForCreator()
	}

	// Ensure new limit volume is big enough
	if coin.LimitVolume.GT(msg.LimitVolume) {
		return nil, types.ErrLimitVolumeBroken(coin.LimitVolume.String(), msg.LimitVolume.String())
	}

	// Update coin metadata
	coin.LimitVolume = msg.LimitVolume
	coin.Identity = msg.Identity

	// Save coin to the storage
	k.SetCoin(ctx, coin)

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeSymbol, coin.Symbol),
		sdk.NewAttribute(types.AttributeLimitVolume, coin.LimitVolume.String()),
		sdk.NewAttribute(types.AttributeIdentity, coin.Identity),
	))

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
		return nil, types.ErrCoinDoesNotExist(coinDenom)
	}

	// NOTE: It was already validated so no need to check error
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	receiver, _ := sdk.AccAddressFromBech32(msg.Receiver)

	// Send coins from the sender to the recipient
	err = k.bankKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(msg.Coin))
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeReceiver, msg.Receiver),
		sdk.NewAttribute(types.AttributeCoin, msg.Coin.String()),
	))

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
			return nil, types.ErrCoinDoesNotExist(coinDenom)
		}

		// Send coins from the sender to the recipient
		err = k.bankKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(msg.Sends[i].Coin))
		if err != nil {
			return nil, types.ErrInternal(err.Error())
		}

		// Emit transaction events
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeReceiver, msg.Sends[i].Receiver),
			sdk.NewAttribute(types.AttributeCoin, msg.Sends[i].Coin.String()),
		))
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
	err := k.sellCoin(ctx, sender, msg.CoinToSell, msg.MinCoinToBuy, true)
	if err != nil {
		return nil, err
	}

	return &types.MsgSellAllCoinResponse{}, nil
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
		return nil, types.ErrUnableDecodeCheck(msg.Check)
	}

	// Parse provided check from raw bytes to ensure it is valid
	check, err := types.ParseCheck(checkBytes)
	if err != nil {
		return nil, types.ErrInvalidCheck(err.Error())
	}
	coinDenom := strings.ToLower(check.Coin)
	coinAmount := check.Amount

	// Decode provided proof from base64 format to raw bytes
	proof, err := base64.StdEncoding.DecodeString(msg.Proof)
	if err != nil {
		return nil, types.ErrUnableDecodeProof()
	}

	// Recover issuer address from check signature
	issuer, err := check.Sender()
	if err != nil {
		return nil, types.ErrUnableRecoverAddress(err.Error())
	}

	// Retrieve the coin specified in the check
	coin, err := k.GetCoin(ctx, coinDenom)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(coinDenom)
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
		return nil, types.ErrInsufficientFunds(
			sdk.NewCoin(coinDenom, coinAmount).String(),
			balance.String(),
		)
	}
	if coinDenom != baseCoinDenom {
		if balanceFeeCoin.Amount.LT(feeAmount) {
			return nil, types.ErrInsufficientFunds(
				sdk.NewCoin(coinDenom, feeAmount).String(),
				balanceFeeCoin.String())
		}
	} else {
		if balance.Amount.LT(coinAmount.Add(feeAmount)) {
			return nil, types.ErrInsufficientFunds(
				sdk.NewCoin(coinDenom, coinAmount).String(),
				balance.String(),
			)
		}
	}

	// Ensure the proper chain ID is specified in the check
	if check.ChainID != ctx.ChainID() {
		return nil, types.ErrInvalidChainID(ctx.ChainID(), check.ChainID)
	}

	// Ensure nonce length
	if len(check.Nonce) > 16 {
		return nil, types.ErrInvalidNonce()
	}

	// Check block number
	if check.DueBlock < uint64(ctx.BlockHeight()) {
		return nil, types.ErrCheckExpired(
			strconv.FormatInt(int64(check.DueBlock), 10))
	}

	// Ensure check is not redeemed yet
	if k.IsCheckRedeemed(ctx, check) {
		return nil, types.ErrCheckRedeemed()
	}

	// Recover public key from check lock
	publicKeyA, err := check.LockPubKey()
	if err != nil {
		return nil, types.ErrUnableRecoverLockPkey(err.Error())
	}

	// Prepare bytes used to recover public key from provided proof
	senderAddressHash := make([]byte, 32)
	hw := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hw, []interface{}{sender})
	if err != nil {
		return nil, types.ErrUnableRPLEncodeCheck(err.Error())
	}
	hw.Sum(senderAddressHash[:0])

	// Recover public key from provided proof
	publicKeyB, err := crypto.Ecrecover(senderAddressHash[:], proof)

	// Compare both public keys to ensure provided proof is correct
	if !bytes.Equal(publicKeyA, publicKeyB) {
		return nil, types.ErrInvalidProof("public keys to ensure provided proof is not correct")
	}

	// Write check to the storage
	k.SetCheck(ctx, check)

	// Send fee from issuer to the module
	// TODO: Make sure it is correct way to get fees
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, issuer, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}

	// Send check coins from issuer to the transaction sender
	err = k.bankKeeper.SendCoins(ctx, issuer, sender, sdk.Coins{})
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeIssuer, issuer.String()),
		sdk.NewAttribute(types.AttributeCoin, sdk.NewCoin(coinDenom, coinAmount).String()),
		sdk.NewAttribute(types.AttributeNonce, new(big.Int).SetBytes(check.Nonce).String()),
		sdk.NewAttribute(types.AttributeDueBlock, strconv.FormatUint(check.DueBlock, 10)),
		sdk.NewAttribute(types.AttributeCommissionRedeemCheck, feeCoin.String()),
	))

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
		return types.ErrCoinDoesNotExist(coinToBuyDenom)
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, coinToSellDenom)
	if err != nil {
		return types.ErrCoinDoesNotExist(coinToSellDenom)
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
			return types.ErrInsufficientCoinReserve()
		}
		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, uint(coinToSell.CRR), amountToBuy)
		amountInBaseCoin = amountToBuy
	default:
		// Buyer buys custom coin for custom coin
		amountInBaseCoin = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, uint(coinToBuy.CRR), amountToBuy)
		if amountInBaseCoin.GT(coinToSell.Reserve) {
			return types.ErrInsufficientCoinReserve()
		}
		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, uint(coinToSell.CRR), amountInBaseCoin)
	}

	// Ensure maximum amount of coins to sell (price guard)
	if amountToSell.GT(maxCoinToSell.Amount) {
		return types.ErrMaximumValueToSellReached(maxCoinToSell.Amount.String(), amountToSell.String())
	}

	// Ensure reserve of the coin to sell does not underflow
	if !k.IsCoinBase(ctx, coinToSell.Symbol) {
		if coinToSell.Reserve.Sub(amountInBaseCoin).LT(types.MinCoinReserve) {
			return types.ErrTxBreaksMinReserveRule(types.MinCoinReserve.String(), amountInBaseCoin.String())
		}
	}

	// Ensure that buyer account holds enough coins to sell
	if balance.Amount.LT(amountToSell) {
		return types.ErrInsufficientFunds(
			sdk.NewCoin(coinToSellDenom, amountToSell).String(),
			balance.String(),
		)
	}

	coinsToSell := sdk.NewCoins(sdk.NewCoin(coinToSellDenom, amountToSell))
	coinsToBuy := sdk.NewCoins(sdk.NewCoin(coinToBuyDenom, amountToBuy))

	// Send sold coins from the buyer to the module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coinsToSell)
	if err != nil {
		// TODO: Change error
		return types.ErrUpdateBalance(sender.String(), err.Error())
	}

	// Burn sold coins from the module
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coinsToSell)
	if err != nil {
		// TODO: Change error
		return types.ErrUpdateBalance(sender.String(), err.Error())
	}

	// Mint bought coins to the module
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coinsToBuy)
	if err != nil {
		// TODO: Change error
		return types.ErrUpdateBalance(sender.String(), err.Error())
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coinsToBuy)
	if err != nil {
		// TODO: Change error
		return types.ErrUpdateBalance(sender.String(), err.Error())
	}

	// Update coins
	if !k.IsCoinBase(ctx, coinToSell.Symbol) {
		k.EditCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountInBaseCoin), coinToSell.Volume.Sub(amountToSell))
	}
	if !k.IsCoinBase(ctx, coinToBuy.Symbol) {
		k.EditCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountInBaseCoin), coinToBuy.Volume.Add(amountToBuy))
	}

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeCoinToBuy, sdk.NewCoin(coinToBuyDenom, amountToBuy).String()),
		sdk.NewAttribute(types.AttributeCoinToSell, sdk.NewCoin(coinToSellDenom, amountToSell).String()),
		sdk.NewAttribute(types.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
	))

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
		return types.ErrCoinDoesNotExist(coinToSellDenom)
	}

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, coinToBuyDenom)
	if err != nil {
		return types.ErrCoinDoesNotExist(coinToBuyDenom)
	}

	// Ensure that seller account holds enough coins to sell
	if balance.Amount.LT(coin.Amount) {
		return types.ErrInsufficientFunds(coin.String(), balance.String())
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
		return types.ErrMinimumValueToBuyReached(amountToBuy.String(), minCoinToBuy.Amount.String())
	}

	// Ensure reserve of the coin to sell does not underflow
	if !k.IsCoinBase(ctx, coinToSell.Symbol) {
		if coinToSell.Reserve.Sub(amountInBaseCoin).LT(types.MinCoinReserve) {
			return types.ErrTxBreaksMinReserveRule(types.MinCoinReserve.String(), amountInBaseCoin.String())
		}
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !k.IsCoinBase(ctx, coinToBuy.Symbol) {
		if coinToBuy.Volume.Add(amountToBuy).GT(coinToBuy.LimitVolume) {
			return types.ErrTxBreaksVolumeLimit(coinToBuy.Volume.Add(amountToBuy).String(), coinToBuy.LimitVolume.String())
		}
	}

	coinsToSell := sdk.NewCoins(sdk.NewCoin(coinToSellDenom, amountToSell))
	coinsToBuy := sdk.NewCoins(sdk.NewCoin(coinToBuyDenom, amountToBuy))

	// Send sold coins from the seller to the module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coinsToSell)
	if err != nil {
		// TODO: Change error
		return types.ErrUpdateBalance(sender.String(), err.Error())
	}

	// Burn sold coins from the module
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coinsToSell)
	if err != nil {
		// TODO: Change error
		return types.ErrUpdateBalance(sender.String(), err.Error())
	}

	// Mint bought coins to the module
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coinsToBuy)
	if err != nil {
		// TODO: Change error
		return types.ErrUpdateBalance(sender.String(), err.Error())
	}

	// Send bought coins from the module to the seller
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coinsToBuy)
	if err != nil {
		// TODO: Change error
		return types.ErrUpdateBalance(sender.String(), err.Error())
	}

	// Update coins
	if !k.IsCoinBase(ctx, coinToSell.Symbol) {
		k.EditCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountInBaseCoin), coinToSell.Volume.Sub(amountToSell))
	}
	if !k.IsCoinBase(ctx, coinToBuy.Symbol) {
		k.EditCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountInBaseCoin), coinToBuy.Volume.Add(amountToBuy))
	}

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeCoinToSell, sdk.NewCoin(coinToSellDenom, amountToSell).String()),
		sdk.NewAttribute(types.AttributeCoinToBuy, sdk.NewCoin(coinToBuyDenom, amountToBuy).String()),
		sdk.NewAttribute(types.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
	))

	return nil
}

func getCreateCoinCommission(symbol string) sdk.Int {
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
