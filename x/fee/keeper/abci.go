package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

func EndBlocker(ctx sdk.Context, k Keeper, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	burningPool := k.authKeeper.GetModuleAccount(ctx, types.BurningPool)
	collectedCoins := k.bankKeeper.GetAllBalances(ctx, burningPool.GetAddress())
	if !collectedCoins.IsZero() {
		err := k.coinKeeper.BurnPoolCoins(ctx, types.BurningPool, collectedCoins)
		if err != nil {
			panic(err)
		}
	}
	return []abci.ValidatorUpdate{}
}
