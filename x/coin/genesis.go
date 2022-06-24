package coin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, gs types.GenesisState) {
	// Initialize params
	keeper.SetParams(ctx, gs.Params)

	// Initialize coins
	for _, coin := range gs.Coins {
		// TODO: Validate firstly
		keeper.SetCoin(ctx, coin)
		// TODO: Is that enough?
	}

	// Initialize checks
	for _, check := range gs.Checks {
		// TODO: Validate firstly
		keeper.SetCheck(ctx, &check)
		// TODO: Is that enough?
	}

	for _, balance := range gs.LegacyBalances {
		// TODO: Validate firstly
		keeper.SetLegacyBalance(ctx, balance)
		// TODO: Is that enough?
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:         k.GetParams(ctx),
		Coins:          k.GetCoins(ctx),
		Checks:         k.GetChecks(ctx),
		LegacyBalances: k.GetLegacyBalances(ctx),
	}
}
