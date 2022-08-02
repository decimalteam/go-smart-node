package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"
)

// SetBaseDenomPrice sets the entire collection of a single denom
func (k Keeper) SetBaseDenomPrice(ctx sdk.Context, price float64) {
	store := ctx.KVStore(k.storeKey)
	collectionKey := types.GetBaseDenomPriceKey()

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz[:], math.Float64bits(price))

	store.Set(collectionKey, bz)
}

// GetBaseDenomPrice returns a collection of NFTs
func (k Keeper) GetBaseDenomPrice(ctx sdk.Context) (price float64, found bool) {
	store := ctx.KVStore(k.storeKey)
	baseDenomPriceKey := types.GetBaseDenomPriceKey()

	bz := store.Get(baseDenomPriceKey)
	if bz == nil {
		return
	}

	return math.Float64frombits(binary.BigEndian.Uint64(bz)), true
}
