package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k Keeper, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	coinInfo, err := k.GetCoin(ctx, k.GetBaseDenom(ctx))
	if err != nil {
		panic(fmt.Errorf("can't get base coin info: %s", err.Error()))
	}
	baseVolume := k.bankKeeper.GetSupply(ctx, coinInfo.Denom)
	if !coinInfo.Volume.Equal(baseVolume.Amount) {
		coinInfo.Volume = baseVolume.Amount
		err = k.UpdateCoinVR(ctx, coinInfo.Denom, coinInfo.Volume, sdkmath.ZeroInt())
		if err != nil {
			panic(fmt.Errorf("can't update coin volume and reserve: %s", err.Error()))
		}
	}
	return []abci.ValidatorUpdate{}
}
