package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRewardPerBlockOverride stores the contract-set per-block reward override
// (in PIP/wei). The override is consumed by GetBlockReward where it is clamped
// to the built-in schedule (reduce-only): it can only LOWER emission, never
// raise it above the scheduled reward. Storing zero is a valid, explicit
// "zero emission" override (distinct from clearing it).
func (k Keeper) SetRewardPerBlockOverride(ctx sdk.Context, amount sdkmath.Int) {
	store := ctx.KVStore(k.storeKey)

	bz, err := amount.Marshal()
	if err != nil {
		panic(err)
	}

	store.Set(types.GetRewardPerBlockOverrideKey(), bz)
}

// ClearRewardPerBlockOverride removes the override; the node returns to its
// built-in emission schedule.
func (k Keeper) ClearRewardPerBlockOverride(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetRewardPerBlockOverrideKey())
}

// GetRewardPerBlockOverride returns the stored per-block reward override and a
// flag indicating whether it is set (enabled). When not set, the returned
// amount is zero and the bool is false.
func (k Keeper) GetRewardPerBlockOverride(ctx sdk.Context) (sdkmath.Int, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetRewardPerBlockOverrideKey())
	if bz == nil {
		return sdk.ZeroInt(), false
	}

	amount := sdk.ZeroInt()
	if err := amount.Unmarshal(bz); err != nil {
		panic(err)
	}

	return amount, true
}

// GetBlockReward returns the effective per-block reward for the given height.
// The contract-set override is applied only when it is LOWER than the scheduled
// reward, so emission can never exceed the pre-update (scheduled) value. After
// the emission schedule ends (scheduled == 0), the result is zero regardless of
// the override.
func (k Keeper) GetBlockReward(ctx sdk.Context, height uint64) sdkmath.Int {
	scheduled := types.GetRewardForBlock(height)

	if amount, ok := k.GetRewardPerBlockOverride(ctx); ok && amount.LT(scheduled) {
		return amount
	}

	return scheduled
}
