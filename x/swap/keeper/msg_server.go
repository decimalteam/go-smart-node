package keeper

import (
	"context"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"

	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) SwapInitialize(goCtx context.Context, msg *types.MsgSwapInitialize) (*types.MsgSwapInitializeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasChain(ctx, msg.DestChain) {
		return nil, types.ErrChainDoesNotExists(strconv.FormatUint(uint64(msg.DestChain), 10))
	}
	if !k.HasChain(ctx, msg.FromChain) {
		return nil, types.ErrChainDoesNotExists(strconv.FormatUint(uint64(msg.FromChain), 10))
	}

	funds := sdk.NewCoins(sdk.NewCoin(strings.ToLower(msg.TokenSymbol), msg.Amount))

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, types.ErrInvalidSenderAddress(msg.Sender)
	}

	ok, err := k.CheckBalance(ctx, sender, funds)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, types.ErrInsufficientAccountFunds(msg.Sender, funds.String())
	}

	err = k.LockFunds(ctx, sender, funds)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventSwapInitialize{
		Sender:            msg.Sender,
		From:              msg.Sender,
		DestChain:         msg.DestChain,
		Recipient:         msg.Recipient,
		Amount:            msg.Amount.String(),
		TransactionNumber: msg.TransactionNumber,
		TokenSymbol:       msg.TokenSymbol,
	})

	return &types.MsgSwapInitializeResponse{}, nil

}

func (k Keeper) SwapRedeem(goCtx context.Context, msg *types.MsgSwapRedeem) (*types.MsgSwapRedeemResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	transactionNumber, ok := sdk.NewIntFromString(msg.TransactionNumber)
	if !ok {
		return nil, types.ErrInvalidTransactionNumber(msg.TransactionNumber)
	}

	hash, err := types.GetHash(transactionNumber, msg.TokenSymbol, msg.Amount, msg.Recipient, msg.FromChain, msg.DestChain)
	if err != nil {
		return nil, err
	}

	if k.HasSwap(ctx, hash) {
		return nil, types.ErrAlreadyRedeemed(hash.String())
	}

	R := big.NewInt(0)
	R.SetBytes(msg.R[:])

	S := big.NewInt(0)
	S.SetBytes(msg.S[:])

	address, err := types.Ecrecover(hash, R, S, sdk.NewInt(int64(msg.V)).BigInt())
	if err != nil {
		return nil, err
	}

	if hex.EncodeToString(address.Bytes()) != types.CheckingAddress {
		return nil, types.ErrInvalidServiceAddress(types.CheckingAddress, hex.EncodeToString(address.Bytes()))
	}

	k.SetSwap(ctx, hash)

	funds := sdk.NewCoins(sdk.NewCoin(strings.ToLower(msg.TokenSymbol), msg.Amount))

	if !k.CheckPoolFunds(ctx, funds) {
		return nil, types.ErrInsufficientPoolFunds(funds.String(), k.GetLockedFunds(ctx).String())
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, types.ErrInvalidSenderAddress(msg.Recipient)
	}

	err = k.UnlockFunds(ctx, recipient, funds)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventSwapRedeem{
		Sender:            msg.Sender,
		From:              msg.Sender,
		DestChain:         msg.DestChain,
		Recipient:         msg.Recipient,
		Amount:            msg.Amount.String(),
		TransactionNumber: msg.TransactionNumber,
		TokenSymbol:       msg.TokenSymbol,
	})

	return &types.MsgSwapRedeemResponse{}, nil

}

func (k Keeper) ChainActivate(goCtx context.Context, msg *types.MsgChainActivate) (*types.MsgChainActivateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain, found := k.GetChain(ctx, msg.ChainNumber)
	if found {
		chain.Active = true
	} else {
		chain = types.NewChain(msg.ChainNumber, msg.ChainName, true)
	}

	k.SetChain(ctx, &chain)

	ctx.EventManager().EmitTypedEvent(&types.EventChainActivate{
		ChainName:   msg.ChainName,
		ChainNumber: msg.ChainNumber,
	})

	return &types.MsgChainActivateResponse{}, nil
}

func (k Keeper) ChainDeactivate(goCtx context.Context, msg *types.MsgChainDeactivate) (*types.MsgChainDeactivateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	chain, found := k.GetChain(ctx, msg.ChainNumber)
	if !found {
		return nil, types.ErrChainDoesNotExists(strconv.FormatUint(uint64(msg.ChainNumber), 10))
	}

	chain.Active = false
	k.SetChain(ctx, &chain)

	ctx.EventManager().EmitTypedEvent(&types.EventChainDeactivate{
		ChainNumber: msg.ChainNumber,
	})

	return &types.MsgChainDeactivateResponse{}, nil
}
