package multisig

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/multisig/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs *types.GenesisState) {
	for _, wallet := range gs.Wallets {
		k.SetWallet(ctx, wallet)
	}
	for _, tx := range gs.Transactions {
		k.SetTransaction(ctx, tx)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	wallets, _ := k.GetAllWallets(ctx)
	txs, _ := k.GetAllTransactions(ctx)
	return &types.GenesisState{
		Wallets:      wallets,
		Transactions: txs,
	}
}
