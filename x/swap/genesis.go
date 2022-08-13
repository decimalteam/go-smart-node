package swap

import (
	"bitbucket.org/decimalteam/go-smart-node/x/swap/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
