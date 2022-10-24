package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// MaxValidators returns maximum number of validators.
func (k Keeper) MaxValidators(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

// MaxDelegations returns maximum number of delegations.
func (k Keeper) MaxDelegations(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxDelegations, &res)
	return
}

// MaxEntries returns maximum number of simultaneous redelegations/undelegations (per trio/pair).
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxEntries, &res)
	return
}

// HistoricalEntries returns number of historical info entries to persist in store.
func (k Keeper) HistoricalEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyHistoricalEntries, &res)
	return
}

// RedelegationTime returns time duration needed to redelegate a stake.
func (k Keeper) RedelegationTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyRedelegationTime, &res)
	return
}

// UndelegationTime returns time duration needed to undelegate a stake.
func (k Keeper) UndelegationTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUndelegationTime, &res)
	return
}

// compatible method for IBCKeeper
func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUndelegationTime, &res)
	return
}

func (k Keeper) SignedBlocksWindow(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeySignedBlocksWindow, &res)
	return
}

func (k Keeper) SignedBlockWindow(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeySignedBlocksWindow, &res)
	return
}

func (k Keeper) MinSignedPerWindow(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMinSignedPerWindow, &res)
	return
}

func (k Keeper) SlashFractionDowntime(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeySlashFractionDowntime, &res)
	return
}

func (k Keeper) SlashFractionDoubleSign(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeySlashFractionDoubleSign, &res)
	return
}

func (k Keeper) BaseDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBaseDenom, &res)
	return
}

// GetParams returns all parameters as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.MaxValidators(ctx),
		k.MaxDelegations(ctx),
		k.MaxEntries(ctx),
		k.HistoricalEntries(ctx),
		k.RedelegationTime(ctx),
		k.UndelegationTime(ctx),
		k.BaseDenom(ctx),
		k.SignedBlocksWindow(ctx),
		k.MinSignedPerWindow(ctx),
		k.SlashFractionDowntime(ctx),
		k.SlashFractionDoubleSign(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
