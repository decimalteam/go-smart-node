package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetSubToken(ctx sdk.Context, nftID string, subToken types.SubToken) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(nftID, subToken.ID)

	bz := k.cdc.MustMarshalLengthPrefixed(&subToken)

	store.Set(subTokenKey, bz)
}

func (k Keeper) GetSubTokens(ctx sdk.Context, nftID string) ([]types.SubToken, error) {
	var subTokens []types.SubToken

	err := k.iterateSubTokens(ctx, nftID, func(subToken types.SubToken) (stop bool) {
		subTokens = append(subTokens, subToken)
		return false
	})
	if err != nil {
		return nil, err
	}

	return subTokens, nil
}

func (k Keeper) GetSubToken(ctx sdk.Context, nftID string, subTokenID uint64) (types.SubToken, bool) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(nftID, subTokenID)
	bz := store.Get(subTokenKey)
	if bz == nil {
		return types.SubToken{}, false
	}

	var subToken types.SubToken
	k.cdc.MustUnmarshalLengthPrefixed(bz, &subToken)

	return subToken, true
}

func (k Keeper) RemoveSubToken(ctx sdk.Context, nftID string, subTokenID uint64) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(nftID, subTokenID)
	store.Delete(subTokenKey)
}

func (k Keeper) iterateSubTokens(ctx sdk.Context, nftID string, handler func(subToken types.SubToken) (stop bool)) error {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetSubTokensKey(nftID))
	for ; iterator.Valid(); iterator.Next() {
		var subToken types.SubToken
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &subToken)
		if handler(subToken) {
			break
		}
	}
	err := iterator.Close()
	if err != nil {
		return err
	}

	return nil
}
