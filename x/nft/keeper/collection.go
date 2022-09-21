package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

// GetCollections returns all the NFTs collections.
func (k *Keeper) GetCollections(ctx sdk.Context) (collections []types.Collection) {
	k.iterateCollections(ctx,
		func(collection *types.Collection) bool {
			collections = append(collections, *collection)
			return false
		},
	)
	return
}

// GetCollection returns the NFT collection.
func (k *Keeper) GetCollection(ctx sdk.Context, creator sdk.AccAddress, denom string) (collection types.Collection, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCollectionKey(creator, denom)

	bz := store.Get(key)
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshalLengthPrefixed(bz, &collection)

	// read collection counter separately
	counter := k.getCollectionCounter(ctx, creator, denom)
	collection.Supply = counter.Supply

	return collection, true
}

// SetCollection writes the NFT collection to the KVStore.
func (k *Keeper) SetCollection(ctx sdk.Context, collection types.Collection) {
	creator := sdk.MustAccAddressFromBech32(collection.Creator)

	store := ctx.KVStore(k.storeKey)
	key := types.GetCollectionKey(creator, collection.Denom)

	// write only creator address and denom to the main record
	bz := k.cdc.MustMarshalLengthPrefixed(&types.Collection{
		Creator: collection.Creator,
		Denom:   collection.Denom,
	})
	store.Set(key, bz)

	// write collection counter separately
	k.setCollectionCounter(ctx, creator, collection.Denom, types.CollectionCounter{
		Supply: collection.Supply,
	})
}

// iterateCollections iterates over all NFT collections created.
func (k *Keeper) iterateCollections(ctx sdk.Context, handler func(collection *types.Collection) (stop bool)) error {
	store := ctx.KVStore(k.storeKey)

	it := sdk.KVStorePrefixIterator(store, types.GetCollectionsKey())
	for ; it.Valid(); it.Next() {
		var collection types.Collection
		k.cdc.MustUnmarshalLengthPrefixed(it.Value(), &collection)

		// read collection counter separately
		counter := k.getCollectionCounter(ctx, sdk.MustAccAddressFromBech32(collection.Creator), collection.Denom)
		collection.Supply = counter.Supply

		if handler(&collection) {
			break
		}
	}

	err := it.Close()
	if err != nil {
		return err
	}

	return nil
}

// getCollectionCounter returns the NFT collection counter.
func (k *Keeper) getCollectionCounter(ctx sdk.Context, creator sdk.AccAddress, denom string) (counter types.CollectionCounter) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCollectionCounterKey(creator, denom)

	bz := store.Get(key)
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshalLengthPrefixed(bz, &counter)
	return
}

// setCollectionCounter writes the NFT collection counter to the KVStore.
func (k *Keeper) setCollectionCounter(ctx sdk.Context, creator sdk.AccAddress, denom string, counter types.CollectionCounter) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCollectionCounterKey(creator, denom)

	bz := k.cdc.MustMarshalLengthPrefixed(&counter)
	store.Set(key, bz)
}
