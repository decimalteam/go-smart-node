package legacy

import (
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	for _, rec := range gs.LegacyRecords {
		k.SetLegacyRecord(ctx, rec)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		LegacyRecords: k.GetLegacyRecords(ctx),
	}
}
