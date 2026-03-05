package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// SetValidatorOfflineSince stores the time when a validator went offline.
func (k Keeper) SetValidatorOfflineSince(ctx sdk.Context, valAddr sdk.ValAddress, t time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorOfflineSinceKey(valAddr), sdk.FormatTimeBytes(t))
}

// GetValidatorOfflineSince returns the time when a validator went offline and whether the entry exists.
func (k Keeper) GetValidatorOfflineSince(ctx sdk.Context, valAddr sdk.ValAddress) (time.Time, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorOfflineSinceKey(valAddr))
	if len(bz) == 0 {
		return time.Time{}, false
	}
	t, err := sdk.ParseTimeBytes(bz)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// DeleteValidatorOfflineSince removes the offline-since entry for a validator.
func (k Keeper) DeleteValidatorOfflineSince(ctx sdk.Context, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorOfflineSinceKey(valAddr))
}

// IterateValidatorOfflineSince iterates over all offline-since entries.
// The callback should return true to stop iteration.
func (k Keeper) IterateValidatorOfflineSince(ctx sdk.Context, fn func(valAddr sdk.ValAddress, t time.Time) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetAllValidatorOfflineSinceKey()
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	prefixLen := len(prefix)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		// key format: [prefix][len_byte][20-byte addr]
		// KVStorePrefixIterator does NOT strip the prefix from keys.
		if len(key) <= prefixLen {
			continue
		}
		remainder := key[prefixLen:]
		addrLen := int(remainder[0])
		if len(remainder) < 1+addrLen {
			continue
		}
		valAddr := sdk.ValAddress(remainder[1 : 1+addrLen])
		t, err := sdk.ParseTimeBytes(iter.Value())
		if err != nil {
			continue
		}
		if fn(valAddr, t) {
			break
		}
	}
}
