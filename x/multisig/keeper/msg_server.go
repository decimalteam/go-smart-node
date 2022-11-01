package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

func (k Keeper) CreateWallet(goCtx context.Context, msg *types.MsgCreateWallet) (*types.MsgCreateWalletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Create new multisig wallet
	wallet, err := types.NewWallet(msg.Owners, msg.Weights, msg.Threshold, ctx.TxBytes())
	if err != nil {
		return nil, errors.UnableToCreateWallet
	}

	// Ensure multisig wallet with the address does not exist
	_, err = k.GetWallet(ctx, wallet.Address)
	if err == nil {
		return nil, errors.WalletAlreadyExists
	}

	adr, err := sdk.AccAddressFromBech32(wallet.Address)
	if err != nil {
		return nil, errors.InvalidWallet
	}
	// Ensure account with multisig address does not exist
	existingAccount := k.accountKeeper.GetAccount(ctx, adr)
	if existingAccount != nil && !existingAccount.GetAddress().Empty() {
		return nil, errors.AccountAlreadyExists
	}

	// Save created multisig wallet to the KVStore
	k.SetWallet(ctx, *wallet)

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventCreateWallet{
		Sender:    msg.Sender,
		Wallet:    wallet.Address,
		Owners:    msg.Owners,
		Weights:   msg.Weights,
		Threshold: msg.Threshold,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgCreateWalletResponse{
		Wallet: wallet.Address,
	}, nil
}

func (k Keeper) CreateTransaction(goCtx context.Context, msg *types.MsgCreateTransaction) (*types.MsgCreateTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve multisig wallet from the KVStore
	wallet, err := k.GetWallet(ctx, msg.Wallet)
	if err != nil {
		return nil, err
	}
	_, err = sdk.AccAddressFromBech32(wallet.Address)
	if err != nil {
		return nil, errors.InvalidWallet
	}

	// Create new multisig transaction
	transaction, err := types.NewTransaction(
		k.cdc,
		msg.Wallet,
		*msg.Content,
		len(wallet.Owners),
		ctx.BlockHeight(),
		ctx.TxBytes(),
	)
	if err != nil {
		return nil, err
	}

	// Check internal transaction
	var internal sdk.Msg
	err = k.cdc.UnpackAny(msg.Content, &internal)
	if err != nil {
		return nil, err
	}
	err = internal.ValidateBasic()
	if err != nil {
		return nil, err
	}
	walletInSigners := false
	for _, signer := range internal.GetSigners() {
		if signer.String() == msg.Wallet {
			walletInSigners = true
		}
	}
	if !walletInSigners {
		return nil, errors.WalletIsNotSignerInInternal
	}

	handler := k.router.Handler(internal)
	if handler == nil {
		return nil, errors.NoHandlerForInternal
	}

	// Save created multisig transaction to the KVStore
	k.SetTransaction(ctx, *transaction)

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventCreateTransaction{
		Sender:      msg.Sender,
		Wallet:      msg.Wallet,
		Transaction: transaction.Id,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	// Sign created multisig transaction by the creator
	_, err = k.SignTransaction(goCtx, &types.MsgSignTransaction{
		Sender: msg.Sender,
		ID:     transaction.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateTransactionResponse{
		ID: transaction.Id,
	}, nil

}

func (k Keeper) SignTransaction(goCtx context.Context, msg *types.MsgSignTransaction) (*types.MsgSignTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.IsCompleted(ctx, msg.ID) {
		return nil, errors.AlreadyEnoughSignatures
	}
	// Retrieve multisig transaction from the KVStore
	transaction, err := k.GetTransaction(ctx, msg.ID)
	if err != nil {
		return nil, err
	}

	// Retrieve multisig wallet from the KVStore
	wallet, err := k.GetWallet(ctx, transaction.Wallet)
	if err != nil {
		return nil, err
	}

	if k.IsSigned(ctx, msg.ID, msg.Sender) {
		return nil, errors.TransactionAlreadySigned
	}

	// Calculate current weight of signatures and check sender
	confirmations := uint32(0)
	senderIsOwner := false
	senderWeight := uint32(0)
	for i, owner := range wallet.Owners {
		if k.IsSigned(ctx, msg.ID, owner) {
			confirmations += wallet.Weights[i]
		}
		if owner == msg.Sender {
			senderIsOwner = true
			senderWeight = wallet.Weights[i]
		}
	}

	if !senderIsOwner {
		return nil, errors.SignerIsNotOwner
	}
	// Ensure current weight of signatures is not enough
	if confirmations >= wallet.Threshold {
		return nil, errors.AlreadyEnoughSignatures
	}

	// Append the signature to the multisig transaction
	k.SetSign(ctx, msg.ID, msg.Sender)

	confirmations += senderWeight

	// Check if new weight of signatures is enough to perform multisig transaction
	confirmed := confirmations >= wallet.Threshold

	// Emit transaction events
	err = events.EmitTypedEvent(ctx, &types.EventSignTransaction{
		Sender:        msg.Sender,
		Wallet:        wallet.Address,
		Transaction:   transaction.Id,
		SignerWeight:  senderWeight,
		Confirmations: confirmations,
		Confirmed:     confirmed,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	if confirmed {
		var msg sdk.Msg
		err = k.cdc.UnpackAny(&transaction.Message, &msg)
		if err != nil {
			return nil, err
		}
		handler := k.router.Handler(msg)
		if handler == nil {
			return nil, errors.NoHandlerForInternal
		}

		res, err := handler(ctx, msg)
		if err != nil {
			return nil, err
		}
		// pass events from handler
		for _, ev := range res.Events {
			ctx.EventManager().EmitEvent(sdk.Event(ev))
		}

		k.SetCompleted(ctx, transaction.Id)
	}

	return &types.MsgSignTransactionResponse{}, nil
}
