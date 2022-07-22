package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetOwnerCollectionByDenom sets a collection of NFT IDs owned by an address
func (k Keeper) SetOwnerCollectionByDenom(ctx sdk.Context, owner sdk.AccAddress, denom string, ownerCollection types.OwnerCollection) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOwnerCollectionByDenomKey(owner, denom)

	store.Set(key, k.cdc.MustMarshalLengthPrefixed(&ownerCollection))
}

// GetOwnerCollections gets all the ID Collections owned by an address
func (k Keeper) GetOwnerCollections(ctx sdk.Context, address sdk.AccAddress) []types.OwnerCollection {
	var collections []types.OwnerCollection
	k.iterateOwnerCollections(ctx, address, func(collection types.OwnerCollection) (stop bool) {
		collections = append(collections, collection)
		return false
	},
	)
	return collections
}

// GetOwnerCollectionByDenom gets the ID Collection owned by an address of a specific denom
func (k Keeper) GetOwnerCollectionByDenom(ctx sdk.Context, address sdk.AccAddress, denom string) (oc types.OwnerCollection, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOwnerCollectionByDenomKey(address, denom))
	if bz == nil {
		return
	}

	var collection types.OwnerCollection

	k.cdc.MustUnmarshalLengthPrefixed(bz, &collection)

	return collection, true
}

// iterateOwners iterates over the OwnerCollection by Owner and performs a function
func (k Keeper) iterateOwnerCollections(
	ctx sdk.Context,
	address sdk.AccAddress,
	handler func(ownerCollection types.OwnerCollection) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetOwnerCollectionsKey(address))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var idCollection types.OwnerCollection
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &idCollection)
		if handler(idCollection) {
			break
		}
	}
}
