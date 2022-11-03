package multisig

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs *types.GenesisState) {
	for _, wallet := range gs.Wallets {
		k.SetWallet(ctx, wallet)
	}
	// convert genesis transactions to coin.MsgMultiSendCoin
	for _, tx := range gs.Transactions {
		// 1. create trasnaction
		newTx := types.Transaction{
			Id:        tx.Id,
			Wallet:    tx.Wallet,
			CreatedAt: tx.CreatedAt,
		}
		msg := cointypes.MsgMultiSendCoin{
			Sender: tx.Wallet,
		}
		entries := []cointypes.MultiSendEntry{}
		for _, coin := range tx.Coins {
			entries = append(entries, cointypes.MultiSendEntry{
				Recipient: tx.Receiver,
				Coin:      coin,
			})
		}
		msg.Sends = entries
		anys, err := sdktx.SetMsgs([]sdk.Msg{&msg})
		if err != nil {
			panic(err)
		}
		newTx.Message = *anys[0]

		k.SetTransaction(ctx, newTx)

		// 2. add signatures
		for _, signer := range tx.Signers {
			if signer != "" {
				k.SetSign(ctx, tx.Id, signer)
			}
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	wallets, _ := k.GetAllWallets(ctx)
	// TODO: make good export for universal transactions
	//txs, _ := k.GetAllTransactions(ctx)
	return &types.GenesisState{
		Wallets: wallets,
		//Transactions: txs,
	}
}
