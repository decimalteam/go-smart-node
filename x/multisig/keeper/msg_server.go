package keeper

import (
	"context"
	"strconv"
	"strings"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateWallet(goCtx context.Context, msg *types.MsgCreateWallet) (*types.MsgCreateWalletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Create new multisig wallet
	wallet, err := types.NewWallet(msg.Owners, msg.Weights, msg.Threshold, ctx.TxBytes())
	if err != nil {
		return nil, types.ErrUnableToCreateWallet(err.Error())
	}

	// Ensure multisig wallet with the address does not exist
	_, err = k.GetWallet(ctx, wallet.Address)
	if err == nil {
		return nil, types.ErrWalletAlreadyExists(wallet.Address)
	}

	adr, err := sdk.AccAddressFromBech32(wallet.Address)
	if err != nil {
		return nil, types.ErrInvalidWallet(wallet.Address)
	}
	// Ensure account with multisig address does not exist
	existingAccount := k.accountKeeper.GetAccount(ctx, adr)
	if existingAccount != nil && !existingAccount.GetAddress().Empty() {
		return nil, types.ErrAccountAlreadyExists(wallet.Address)
	}

	// Save created multisig wallet to the KVStore
	k.SetWallet(ctx, *wallet)

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyOwners, strings.Join(msg.Owners, ",")),
		sdk.NewAttribute(types.AttributeKeyWeights, helpers.JoinUints64(msg.Weights)),
		sdk.NewAttribute(types.AttributeKeyThreshold, strconv.FormatUint(msg.Threshold, 10)),
		sdk.NewAttribute(types.AttributeKeyWallet, wallet.Address),
	))

	return &types.MsgCreateWalletResponse{
		Wallet: wallet.Address,
	}, nil
}

func (k Keeper) CreateTransaction(goCtx context.Context, msg *types.MsgCreateTransaction) (*types.MsgCreateTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve multisig wallet from the KVStore
	wallet, err := k.GetWallet(ctx, msg.Wallet)
	if err != nil {
		return nil, types.ErrWalletAccountNotFound(msg.Wallet)
	}
	adr, err := sdk.AccAddressFromBech32(wallet.Address)
	if err != nil {
		return nil, types.ErrInvalidWallet(wallet.Address)
	}
	// Retrieve coins hold on the multisig wallet
	walletCoins := k.bankKeeper.GetAllBalances(ctx, adr)

	// Ensure there are enough coins on the multisig wallet
	for _, coin := range msg.Coins {
		coinName := strings.ToLower(coin.Denom)
		if walletCoins.AmountOf(coinName).LT(coin.Amount) {
			return nil, types.ErrInsufficientFunds(coin.String(), sdk.NewCoin(coinName, walletCoins.AmountOf(coinName)).String())
		}
	}

	// Create new multisig transaction
	transaction, err := types.NewTransaction(
		msg.Wallet,
		msg.Receiver,
		msg.Coins,
		len(wallet.Owners),
		ctx.BlockHeight(),
		ctx.TxBytes(),
	)
	if err != nil {
		return nil, types.ErrUnableToCreateTransaction(err.Error())
	}

	// Save created multisig transaction to the KVStore
	k.SetTransaction(ctx, *transaction)

	// Sign created multisig transaction by the creator
	_, err = k.SignTransaction(goCtx, &types.MsgSignTransaction{
		Sender: msg.Sender,
		TxID:   transaction.Id,
	})
	if err != nil {
		return nil, err
	}

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyWallet, msg.Wallet),
		sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		sdk.NewAttribute(types.AttributeKeyCoins, msg.Coins.String()),
		sdk.NewAttribute(types.AttributeKeyTransaction, transaction.Id),
	))

	return &types.MsgCreateTransactionResponse{
		TxID: transaction.Id,
	}, nil
}

func (k Keeper) SignTransaction(goCtx context.Context, msg *types.MsgSignTransaction) (*types.MsgSignTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve multisig transaction from the KVStore
	transaction, err := k.GetTransaction(ctx, msg.TxID)
	if err != nil {
		return nil, types.ErrTransactionNotFound(msg.TxID)
	}

	// Retrieve multisig wallet from the KVStore
	wallet, err := k.GetWallet(ctx, transaction.Wallet)
	if err != nil {
		return nil, types.ErrWalletAccountNotFound(transaction.Wallet)
	}

	// Calculate current weight of signatures
	confirmations := uint64(0)
	for i := 0; i < len(wallet.Owners); i++ {
		if transaction.Signers[i] != "" {
			confirmations += wallet.Weights[i]
		}
	}

	// Ensure current weight of signatures is not enough
	if confirmations >= wallet.Threshold {
		return nil, types.ErrAlreadyEnoughSignatures(strconv.FormatUint(confirmations, 10), strconv.FormatUint(wallet.Threshold, 10))
	}

	// Append the signature to the multisig transaction
	weight := uint64(0)
	signed := false
	for i := 0; i < len(wallet.Owners); i++ {
		if wallet.Owners[i] != msg.Sender {
			continue
		}
		if transaction.Signers[i] != "" {
			return nil, types.ErrTransactionAlreadySigned(msg.Sender)
		}
		signed = true
		weight = wallet.Weights[i]
		confirmations += weight
		transaction.Signers[i] = msg.Sender
		break
	}
	if !signed {
		return nil, types.ErrSignerIsNotOwner(msg.Sender, transaction.Wallet)
	}

	// Save updated multisig transaction to the KVStore
	k.SetTransaction(ctx, transaction)

	// Check if new weight of signatures is enough to perform multisig transaction
	confirmed := confirmations >= wallet.Threshold
	if confirmed {
		wAdr, err := sdk.AccAddressFromBech32(wallet.Address)
		if err != nil {
			return nil, types.ErrInvalidWallet(wallet.Address)
		}
		rAdr, err := sdk.AccAddressFromBech32(transaction.Receiver)
		if err != nil {
			return nil, types.ErrInvalidReceiver(transaction.Receiver)
		}
		// Perform transaction
		err = k.bankKeeper.SendCoins(ctx, wAdr, rAdr, transaction.Coins)
		if err != nil {
			return nil, types.ErrUnablePreformTransaction(transaction.Id, err.Error())
		}
	}

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyWallet, wallet.Address),
		sdk.NewAttribute(types.AttributeKeyTransaction, msg.TxID),
		sdk.NewAttribute(types.AttributeKeySignerWeight, strconv.FormatUint(uint64(weight), 10)),
		sdk.NewAttribute(types.AttributeKeyConfirmations, strconv.FormatUint(uint64(confirmations), 10)),
		sdk.NewAttribute(types.AttributeKeyConfirmed, strconv.FormatBool(confirmed)),
	))

	return &types.MsgSignTransactionResponse{}, nil
}