package coin

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(
	ctx sdk.Context,
	keeper keeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	genesisState types.GenesisState,
) {
	keeper.SetParams(ctx, genesisState.Params)

	// ensure the module account is set on genesis
	if acc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		panic(fmt.Sprintf("the %s module account has not been set", types.ModuleName))
	}

	for _, coin := range genesisState.Coins {
		// TODO: Validate firstly
		keeper.SetCoin(ctx, coin)
		// TODO: Is that enough?
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
		Coins:  k.GetAllCoins(ctx),
	}
}
