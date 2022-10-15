package coin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs *types.GenesisState) {
	// Initialize params
	k.SetParams(ctx, gs.Params)

	// Initialize base coin
	coin := types.Coin{
		Denom:  gs.Params.BaseDenom,
		Title:  gs.Params.BaseTitle,
		Volume: gs.Params.BaseVolume,
	}
	k.SetCoin(ctx, coin)

	// Initialize coins
	for _, coin := range gs.Coins {
		// TODO: Validate firstly
		k.SetCoin(ctx, coin)
		// TODO: Is that enough?
	}

	// Initialize checks
	for _, check := range gs.Checks {
		// TODO: Validate firstly
		k.SetCheck(ctx, &check)
		// TODO: Is that enough?
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
		Coins:  k.GetCoins(ctx),
		Checks: k.GetChecks(ctx),
	}
}
