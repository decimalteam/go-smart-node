package fee

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs *types.GenesisState) {
	k.SetModuleParams(ctx, gs.Params)
	err := k.SavePrice(ctx, gs.InitialPrice)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	price, err := k.GetPrice(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params:       k.GetModuleParams(ctx),
		InitialPrice: price,
	}
}
