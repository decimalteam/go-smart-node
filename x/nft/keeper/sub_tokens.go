package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetSubTokens(ctx sdk.Context, id string) (subTokens []types.SubToken) {
	k.iterateSubTokens(ctx, id, func(subToken types.SubToken) (stop bool) {
		subTokens = append(subTokens, subToken)
		return false
	})

	return
}

func (k Keeper) GetSubToken(ctx sdk.Context, id string, subTokenID uint64) (types.SubToken, bool) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(id, subTokenID)
	bz := store.Get(subTokenKey)
	if bz == nil {
		return types.SubToken{}, false
	}

	var subToken types.SubToken
	k.cdc.MustUnmarshalLengthPrefixed(bz, &subToken)

	return subToken, true
}

func (k Keeper) SetSubToken(ctx sdk.Context, id string, subToken types.SubToken) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(id, subToken.ID)

	bz := k.cdc.MustMarshalLengthPrefixed(&subToken)

	store.Set(subTokenKey, bz)
}

func (k Keeper) RemoveSubToken(ctx sdk.Context, id string, subTokenID uint64) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(id, subTokenID)
	store.Delete(subTokenKey)
}

func (k Keeper) iterateSubTokens(ctx sdk.Context, id string, handler func(subToken types.SubToken) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetSubTokensKey(id))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var subToken types.SubToken
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &subToken)
		if handler(subToken) {
			break
		}
	}
}
