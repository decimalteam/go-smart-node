package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetNFT sets nft to store
func (k Keeper) SetNFT(ctx sdk.Context, denom, id string, nft types.BaseNFT) error {
	_, found := k.GetCollection(ctx, denom)
	if !found {
		return errors.UnknownCollection
	}

	store := ctx.KVStore(k.storeKey)
	nftKey := types.GetNFTKey(id)

	bz := k.cdc.MustMarshalLengthPrefixed(&nft)
	store.Set(nftKey, bz)

	k.setTokenURI(ctx, nft.TokenURI)

	return nil
}

// GetNFT gets the entire NFT metadata struct for a uint64
func (k Keeper) GetNFT(ctx sdk.Context, denom, id string) (types.BaseNFT, error) {
	_, found := k.GetCollection(ctx, denom)
	if !found {
		return types.BaseNFT{}, errors.UnknownCollection
	}

	store := ctx.KVStore(k.storeKey)
	nftKey := types.GetNFTKey(id)
	bz := store.Get(nftKey)
	if bz == nil {
		return types.BaseNFT{}, errors.UnknownNFT
	}

	var nft types.BaseNFT

	k.cdc.MustUnmarshalLengthPrefixed(bz, &nft)

	return nft, nil
}

// GetNFTs returns all matched NFTs
func (k Keeper) GetNFTs(ctx sdk.Context) ([]types.BaseNFT, error) {
	var nfts []types.BaseNFT
	err := k.iterateNFTs(ctx, func(nft types.BaseNFT) (stop bool) {
		nfts = append(nfts, nft)
		return false
	})
	if err != nil {
		return nil, err
	}

	return nfts, nil
}

// HasTokenID check if nft exists
func (k Keeper) HasTokenID(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	nftKey := types.GetNFTKey(id)

	return store.Has(nftKey)
}

func (k Keeper) HasTokenURI(ctx sdk.Context, tokenURI string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenURIKey := types.GetTokenURIKey(tokenURI)

	return store.Has(tokenURIKey)
}

func (k Keeper) setTokenURI(ctx sdk.Context, tokenURI string) {
	store := ctx.KVStore(k.storeKey)
	tokenURIKey := types.GetTokenURIKey(tokenURI)

	store.Set(tokenURIKey, []byte{})
}

// iterateNFTs iterates over NFTs and performs a function
func (k Keeper) iterateNFTs(ctx sdk.Context, handler func(collection types.BaseNFT) (stop bool)) error {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.NFTKeyPrefix)
	for ; iterator.Valid(); iterator.Next() {
		var nft types.BaseNFT
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &nft)
		if handler(nft) {
			break
		}
	}

	err := iterator.Close()
	if err != nil {
		return err
	}

	return nil
}
