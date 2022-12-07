package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

// GetSubTokens returns existing NFT sub-tokens from the NFT token with specified ID.
func (k *Keeper) GetSubTokens(ctx sdk.Context, id string) (subTokens []types.SubToken) {
	k.iterateSubTokens(ctx, id,
		func(subToken *types.SubToken) bool {
			subTokens = append(subTokens, *subToken)
			return false
		},
	)
	return
}

// GetSubToken returns the NFT sub-token with specified ID.
func (k *Keeper) GetSubToken(ctx sdk.Context, id string, index uint32) (subToken types.SubToken, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSubTokenKey(id, index)

	bz := store.Get(key)
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshalLengthPrefixed(bz, &subToken)
	return subToken, true
}

// setSubToken writes the NFT sub-token to the KVStore.
func (k *Keeper) SetSubToken(ctx sdk.Context, id string, subToken types.SubToken) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSubTokenKey(id, subToken.ID)

	bz := k.cdc.MustMarshalLengthPrefixed(&subToken)
	store.Set(key, bz)
}

// ReplaceSubTokenOwner replaces sub token owner and make changes in indexes.
// need for legacy module, for common reason there is TransferSubTokens transaction
// newOwner must be valid bech32 address, subtoken.Owner MUST BE legacy address with prefix 'dx'
func (k *Keeper) ReplaceSubTokenOwner(ctx sdk.Context, id string, index uint32, newOwner string) error {
	subtoken, found := k.GetSubToken(ctx, id, index)
	if !found {
		return errors.SubTokenDoesNotExists
	}
	newOwnerSdk, err := sdk.AccAddressFromBech32(newOwner)
	if err != nil {
		return err
	}
	oldOwnerSdk, err := sdk.GetFromBech32(subtoken.Owner, "dx")
	if err != nil {
		return err
	}
	subtoken.Owner = newOwner
	k.SetSubToken(ctx, id, subtoken)
	k.transferSubToken(ctx, oldOwnerSdk, newOwnerSdk, id, index)
	return nil
}

// removeSubToken deletes the NFT sub-token from the KVStore.
func (k *Keeper) removeSubToken(ctx sdk.Context, id string, index uint32) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSubTokenKey(id, index)
	store.Delete(key)
}

// transferSubToken transfers NFT sub-token, removes previous NFT sub-token by owner index and writes new one to the KVStore.
func (k *Keeper) transferSubToken(ctx sdk.Context, oldOwner sdk.AccAddress, newOwner sdk.AccAddress, id string, index uint32) {
	k.removeSubTokenByOwner(ctx, oldOwner, id, index)
	k.setSubTokenByOwner(ctx, newOwner, id, index)
}

// iterateSubTokens iterates over existing NFT sub-tokens from the NFT token with specified ID.
func (k *Keeper) iterateSubTokens(ctx sdk.Context, id string, handler func(subToken *types.SubToken) (stop bool)) error {
	store := ctx.KVStore(k.storeKey)

	it := sdk.KVStorePrefixIterator(store, types.GetSubTokensKey(id))
	for ; it.Valid(); it.Next() {
		var subToken types.SubToken
		k.cdc.MustUnmarshalLengthPrefixed(it.Value(), &subToken)

		if handler(&subToken) {
			break
		}
	}

	err := it.Close()
	if err != nil {
		return err
	}

	return nil
}

// setSubTokenByOwner writes the NFT sub-token index by owner address to the KVStore.
func (k *Keeper) setSubTokenByOwner(ctx sdk.Context, owner sdk.AccAddress, id string, index uint32) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSubTokenByOwnerKey(owner, id, index)
	store.Set(key, []byte{})
}

// removeSubTokenByOwner removes the NFT sub-token index by owner address from the KVStore.
func (k *Keeper) removeSubTokenByOwner(ctx sdk.Context, owner sdk.AccAddress, id string, index uint32) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSubTokenByOwnerKey(owner, id, index)
	store.Delete(key)
}
