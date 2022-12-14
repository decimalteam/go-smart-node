package legacy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/legacy/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs *types.GenesisState) {
	for _, rec := range gs.Records {
		k.SetLegacyRecord(ctx, rec)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Records: k.GetLegacyRecords(ctx),
	}
}
