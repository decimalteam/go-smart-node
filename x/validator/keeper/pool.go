package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// GetBondedPool returns the bonded coins pool's module account.
func (k Keeper) GetBondedPool(ctx sdk.Context) (bondedPool authtypes.ModuleAccountI) {
	return k.authKeeper.GetModuleAccount(ctx, types.BondedPoolName)
}

// GetNotBondedPool returns the not bonded coins pool's module account.
func (k Keeper) GetNotBondedPool(ctx sdk.Context) (notBondedPool authtypes.ModuleAccountI) {
	return k.authKeeper.GetModuleAccount(ctx, types.NotBondedPoolName)
}

// BondedTotal returns the total staking coin supply which is bonded.
func (k Keeper) BondedTotal(ctx sdk.Context, denom string) sdkmath.Int {
	bondedPool := k.GetBondedPool(ctx)
	return k.bankKeeper.GetBalance(ctx, bondedPool.GetAddress(), denom).Amount
}

// BondedRatio returns the fraction of the staking coin which are currently bonded.
func (k Keeper) BondedRatio(ctx sdk.Context, denom string) sdk.Dec {
	stakeSupply := k.bankKeeper.GetSupply(ctx, denom).Amount
	if stakeSupply.IsPositive() {
		return sdk.NewDecFromInt(k.BondedTotal(ctx, denom)).QuoInt(stakeSupply)
	}
	return sdk.ZeroDec()
}

// sendCoinsToBonded transfers coins from the not bonded to the bonded pool within staking.
func (k Keeper) sendCoinsToBonded(ctx sdk.Context, coins sdk.Coins) {
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.NotBondedPoolName, types.BondedPoolName, coins); err != nil {
		panic(err)
	}
}

// sendCoinsToNotBonded transfers coins from the bonded to the not bonded pool within staking.
func (k Keeper) sendCoinsToNotBonded(ctx sdk.Context, coins sdk.Coins) {
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.BondedPoolName, types.NotBondedPoolName, coins); err != nil {
		panic(err)
	}
}

// burnCoinsFromBonded removes coins from the bonded pool module account.
func (k Keeper) burnCoinsFromBonded(ctx sdk.Context, coins sdk.Coins) error {
	if coins.IsZero() {
		return nil
	}
	return k.bankKeeper.BurnCoins(ctx, types.BondedPoolName, coins)
}

// burnCoinsFromNotBonded removes coins from the not bonded pool module account.
func (k Keeper) burnCoinsFromNotBonded(ctx sdk.Context, coins sdk.Coins) error {
	if coins.IsZero() {
		return nil
	}
	return k.bankKeeper.BurnCoins(ctx, types.NotBondedPoolName, coins)
}
