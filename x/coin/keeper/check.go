package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// IsCheckRedeemed returns true if provided check is marked as redeemed in KVStore.
func (k *Keeper) IsCheckRedeemed(ctx sdk.Context, check *types.Check) bool {
	checkHash := check.HashFull()
	store := ctx.KVStore(k.storeKey)
	key := types.GetCheckKey(checkHash[:])
	value := store.Get(key)
	if len(value) == 0 {
		return false
	}
	var c types.Check
	return k.cdc.UnmarshalLengthPrefixed(value, &c) == nil
}

// GetChecks returns all checks existing in KVStore.
func (k *Keeper) GetChecks(ctx sdk.Context) (checks []types.Check) {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, types.GetChecksKey())
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var check types.Check
		err := k.cdc.UnmarshalLengthPrefixed(it.Value(), &check)
		if err != nil {
			panic(err)
		}
		checks = append(checks, check)
	}

	return checks
}

// GetCheck returns the redeemed check if exists in KVStore.
func (k *Keeper) GetCheck(ctx sdk.Context, checkHash []byte) (check types.Check, err error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCheckKey(checkHash)
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.CheckDoesNotExist
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &check)
	return
}

// SetCheck writes check to KVStore.
func (k *Keeper) SetCheck(ctx sdk.Context, check *types.Check) {
	checkHash := check.HashFull()
	store := ctx.KVStore(k.storeKey)
	key := types.GetCheckKey(checkHash[:])
	value := k.cdc.MustMarshalLengthPrefixed(check)
	store.Set(key, value)
}
