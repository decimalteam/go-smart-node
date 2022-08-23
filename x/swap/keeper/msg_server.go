package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/swap/errors"
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
		return nil, errors.ChainDoesNotExists
	}
	if !k.HasChain(ctx, msg.FromChain) {
		return nil, errors.ChainDoesNotExists
	}

	funds := sdk.NewCoins(sdk.NewCoin(strings.ToLower(msg.TokenSymbol), msg.Amount))

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errors.InvalidSenderAddress
	}

	ok, err := k.CheckBalance(ctx, sender, funds)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.InsufficientAccountFunds
	}

	err = k.LockFunds(ctx, sender, funds)
	if err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventSwapInitialize{
		Sender:            msg.Sender,
		From:              msg.Sender,
		DestChain:         msg.DestChain,
		Recipient:         msg.Recipient,
		Amount:            msg.Amount.String(),
		TransactionNumber: msg.TransactionNumber,
		TokenSymbol:       msg.TokenSymbol,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgSwapInitializeResponse{}, nil

}

func (k Keeper) SwapRedeem(goCtx context.Context, msg *types.MsgSwapRedeem) (*types.MsgSwapRedeemResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	transactionNumber, ok := sdk.NewIntFromString(msg.TransactionNumber)
	if !ok {
		return nil, errors.InvalidTransactionNumber
	}

	hash, err := types.GetHash(transactionNumber, msg.TokenSymbol, msg.Amount, msg.Recipient, msg.FromChain, msg.DestChain)
	if err != nil {
		return nil, err
	}

	if k.HasSwap(ctx, hash) {
		return nil, errors.AlreadyRedeemed
	}

	R := big.NewInt(0)
	R.SetBytes(msg.R[:])

	S := big.NewInt(0)
	S.SetBytes(msg.S[:])

	address, err := types.Ecrecover(hash, R, S, sdk.NewInt(int64(msg.V)).BigInt())
	if err != nil {
		return nil, err
	}

	params := k.GetParams(ctx)

	if hex.EncodeToString(address.Bytes()) != params.CheckingAddress {
		return nil, errors.InvalidServiceAddress
	}

	k.SetSwap(ctx, hash)

	funds := sdk.NewCoins(sdk.NewCoin(strings.ToLower(msg.TokenSymbol), msg.Amount))

	if !k.CheckPoolFunds(ctx, funds) {
		return nil, errors.InsufficientPoolFunds
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, errors.InvalidSenderAddress
	}

	err = k.UnlockFunds(ctx, recipient, funds)
	if err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventSwapRedeem{
		Sender:            msg.Sender,
		From:              msg.Sender,
		DestChain:         msg.DestChain,
		Recipient:         msg.Recipient,
		Amount:            msg.Amount.String(),
		TransactionNumber: msg.TransactionNumber,
		TokenSymbol:       msg.TokenSymbol,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgSwapRedeemResponse{}, nil

}

func (k Keeper) ChainActivate(goCtx context.Context, msg *types.MsgChainActivate) (*types.MsgChainActivateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if msg.Sender != params.ServiceAddress {
		return nil, errors.SenderIsNotSwapService
	}

	chain, found := k.GetChain(ctx, msg.ChainNumber)
	if found {
		chain.Active = true
	} else {
		chain = types.NewChain(msg.ChainNumber, msg.ChainName, true)
	}

	k.SetChain(ctx, &chain)

	err := ctx.EventManager().EmitTypedEvent(&types.EventChainActivate{
		ChainName:   msg.ChainName,
		ChainNumber: msg.ChainNumber,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgChainActivateResponse{}, nil
}

func (k Keeper) ChainDeactivate(goCtx context.Context, msg *types.MsgChainDeactivate) (*types.MsgChainDeactivateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if msg.Sender != params.ServiceAddress {
		return nil, errors.SenderIsNotSwapService
	}

	chain, found := k.GetChain(ctx, msg.ChainNumber)
	if !found {
		return nil, errors.ChainDoesNotExists
	}

	chain.Active = false
	k.SetChain(ctx, &chain)

	err := ctx.EventManager().EmitTypedEvent(&types.EventChainDeactivate{
		ChainNumber: msg.ChainNumber,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgChainDeactivateResponse{}, nil
}
