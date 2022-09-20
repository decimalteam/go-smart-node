package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

// GetTokens returns all NFT tokens from the NFT collection with specified creator and denom.
func (k *Keeper) GetTokens(ctx sdk.Context, creator sdk.AccAddress, denom string) (tokens []types.Token) {
	k.iterateTokens(ctx, creator, denom,
		func(token *types.Token) bool {
			tokens = append(tokens, *token)
			return false
		},
	)
	return
}

// GetToken returns the NFT token with specified ID.
func (k *Keeper) GetToken(ctx sdk.Context, id string) (token types.Token, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenKey(id)

	bz := store.Get(key)
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshalLengthPrefixed(bz, &token)

	// read token counter separately
	counter := k.getTokenCounter(ctx, id)
	token.Minted, token.Burnt = counter.Minted, counter.Burnt

	return token, true
}

// CreateToken writes the new NFT token to the KVStore.
func (k *Keeper) CreateToken(ctx sdk.Context, collection types.Collection, token types.Token) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenKey(token.ID)

	// write token but ignore minted, burnt and sub-tokens
	t := token
	t.Minted, t.Burnt, t.SubTokens = 0, 0, nil
	bz := k.cdc.MustMarshalLengthPrefixed(&t)
	store.Set(key, bz)

	// write token counter separately
	k.setTokenCounter(ctx, token.ID, types.TokenCounter{
		Minted: token.Minted,
		Burnt:  token.Burnt,
	})

	// write token URI index
	k.setTokenURI(ctx, token.URI)

	// write token by collection index
	creator := sdk.MustAccAddressFromBech32(token.Creator)
	k.setTokenByCollection(ctx, creator, collection.Denom, token.ID)
}

func (k *Keeper) setToken(ctx sdk.Context, t types.Token) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenKey(t.ID)

	bz := k.cdc.MustMarshalLengthPrefixed(&t)
	store.Set(key, bz)
}

// updateTokenURI removes previous NFT token URI and writes new one to the KVStore.
func (k *Keeper) updateTokenURI(ctx sdk.Context, oldTokenURI string, newTokenURI string) {
	k.removeTokenURI(ctx, oldTokenURI)
	k.setTokenURI(ctx, newTokenURI)
}

// iterateTokens iterates over all NFT tokens from the NFT collection with specified creator and denom.
func (k *Keeper) iterateTokens(ctx sdk.Context, creator sdk.AccAddress, denom string, handler func(token *types.Token) (stop bool)) error {
	store := ctx.KVStore(k.storeKey)

	rootKey := types.GetTokensByCollectionKey(creator, denom)
	it := sdk.KVStorePrefixIterator(store, rootKey)
	for ; it.Valid(); it.Next() {
		var token types.Token

		tokenKey := it.Key()[len(rootKey):]
		bz := store.Get(types.GetTokenKeyByIDHash(tokenKey))

		k.cdc.MustUnmarshalLengthPrefixed(bz, &token)

		// read token counter separately
		counter := k.getTokenCounter(ctx, token.ID)
		token.Minted, token.Burnt = counter.Minted, counter.Burnt

		if handler(&token) {
			break
		}
	}

	err := it.Close()
	if err != nil {
		return err
	}

	return nil
}

// getTokenCounter returns the NFT token counter.
func (k *Keeper) getTokenCounter(ctx sdk.Context, id string) (counter types.TokenCounter) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenCounterKey(id)

	bz := store.Get(key)
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshalLengthPrefixed(bz, &counter)
	return
}

// setTokenCounter writes the NFT token counter to the KVStore.
func (k *Keeper) setTokenCounter(ctx sdk.Context, id string, counter types.TokenCounter) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenCounterKey(id)

	bz := k.cdc.MustMarshalLengthPrefixed(&counter)
	store.Set(key, bz)
}

// hasTokenURI returns true if specified NFT token URI exists in the KVStore.
func (k *Keeper) hasTokenURI(ctx sdk.Context, tokenURI string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenURIKey(tokenURI)
	return store.Has(key)
}

// setTokenURI writes the NFT token URI hash index to the KVStore.
func (k *Keeper) setTokenURI(ctx sdk.Context, tokenURI string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenURIKey(tokenURI)
	store.Set(key, []byte{})
}

// removeTokenURI writes the NFT token URI hash index to the KVStore.
func (k *Keeper) removeTokenURI(ctx sdk.Context, tokenURI string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenURIKey(tokenURI)
	store.Delete(key)
}

// setTokenByCollection writes the NFT token index by NFT collection to the KVStore.
func (k *Keeper) setTokenByCollection(ctx sdk.Context, creator sdk.AccAddress, denom string, id string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenByCollectionKey(creator, denom, id)
	store.Set(key, []byte{})
}
