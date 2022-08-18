package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetCollection sets the entire collection of a single denom
func (k Keeper) SetCollection(ctx sdk.Context, denom string, collection types.Collection) {
	store := ctx.KVStore(k.storeKey)
	collectionKey := types.GetCollectionKey(denom)

	bz := k.cdc.MustMarshalLengthPrefixed(&collection)
	store.Set(collectionKey, bz)
}

// GetCollection returns a collection of NFTs
func (k Keeper) GetCollection(ctx sdk.Context, denom string) (collection types.Collection, found bool) {
	store := ctx.KVStore(k.storeKey)
	collectionKey := types.GetCollectionKey(denom)

	bz := store.Get(collectionKey)
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshalLengthPrefixed(bz, &collection)
	return collection, true
}

// GetCollections returns all the NFTs collections
func (k Keeper) GetCollections(ctx sdk.Context) (collections []types.Collection) {
	k.iterateCollections(ctx,
		func(collection types.Collection) (stop bool) {
			collections = append(collections, collection)
			return false
		},
	)
	return
}

// GetDenoms returns all the NFT denoms
func (k Keeper) GetDenoms(ctx sdk.Context) ([]string, error) {
	var denoms []string
	err := k.iterateCollections(ctx,
		func(collection types.Collection) (stop bool) {
			denoms = append(denoms, collection.Denom)
			return false
		},
	)
	if err != nil {
		return nil, err
	}

	return denoms, nil
}

func (k Keeper) iterateCollections(ctx sdk.Context, handler func(collection types.Collection) (stop bool)) error {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CollectionsKeyPrefix)
	for ; iterator.Valid(); iterator.Next() {
		var collection types.Collection
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &collection)
		if handler(collection) {
			break
		}
	}

	err := iterator.Close()
	if err != nil {
		return err
	}

	return nil
}
