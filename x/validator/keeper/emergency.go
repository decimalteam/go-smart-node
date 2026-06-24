package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// GetHaltInfo returns the current emergency hard-halt info and whether it is set.
func (k Keeper) GetHaltInfo(ctx sdk.Context) (types.HaltInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetHaltInfoKey())
	if bz == nil {
		return types.HaltInfo{}, false
	}
	var info types.HaltInfo
	k.cdc.MustUnmarshal(bz, &info)
	return info, true
}

// SetHaltInfo persists the emergency hard-halt info.
func (k Keeper) SetHaltInfo(ctx sdk.Context, info types.HaltInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetHaltInfoKey(), k.cdc.MustMarshal(&info))
}

// IsHaltScheduled reports whether a hard halt is currently scheduled (active).
func (k Keeper) IsHaltScheduled(ctx sdk.Context) bool {
	info, found := k.GetHaltInfo(ctx)
	return found && info.Active
}

// ShouldHardHalt reports whether the BeginBlocker must panic-halt the chain at
// the current height, returning the active HaltInfo when so. It always returns
// false when the node is booted with --unsafe-skip-halt (k.skipHalt), which is
// how operators get past a halt to perform a coordinated resume.
func (k Keeper) ShouldHardHalt(ctx sdk.Context) (types.HaltInfo, bool) {
	if k.skipHalt {
		return types.HaltInfo{}, false
	}
	info, found := k.GetHaltInfo(ctx)
	if found && info.Active && ctx.BlockHeight() >= info.Height {
		return info, true
	}
	return types.HaltInfo{}, false
}

// GetFreezeInfo returns the current emergency soft-freeze info and whether it is set.
func (k Keeper) GetFreezeInfo(ctx sdk.Context) (types.FreezeInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetFreezeInfoKey())
	if bz == nil {
		return types.FreezeInfo{}, false
	}
	var info types.FreezeInfo
	k.cdc.MustUnmarshal(bz, &info)
	return info, true
}

// SetFreezeInfo persists the emergency soft-freeze info.
func (k Keeper) SetFreezeInfo(ctx sdk.Context, info types.FreezeInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetFreezeInfoKey(), k.cdc.MustMarshal(&info))
}

// IsFrozen reports whether the chain is currently soft-frozen. It is consulted
// by the ante handler to reject non-emergency transactions while frozen.
func (k Keeper) IsFrozen(ctx sdk.Context) bool {
	info, found := k.GetFreezeInfo(ctx)
	return found && info.Frozen
}
