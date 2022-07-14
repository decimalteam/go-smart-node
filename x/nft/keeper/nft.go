package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/exported"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IsNFT returns whether an NFT exists
func (k Keeper) IsNFT(ctx sdk.Context, denom, id string) (exists bool) {
	_, err := k.GetNFT(ctx, denom, id)
	return err == nil
}

// GetNFT gets the entire NFT metadata struct for a uint64
func (k Keeper) GetNFT(ctx sdk.Context, denom, id string) (exported.NFT, error) {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return nil, types.ErrUnknownCollection(denom)
	}

	nft, err := collection.GetNFT(id)
	if err != nil {
		return nil, err
	}

	return nft, nil
}

func (k Keeper) GetSubToken(ctx sdk.Context, denom, id string, subTokenID int64) (sdk.Int, bool) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(denom, id, subTokenID)
	bz := store.Get(subTokenKey)
	if bz == nil {
		return sdk.Int{}, false
	}

	reserve := sdk.ZeroInt()

	err := reserve.Unmarshal(bz)
	if err != nil {
		panic(err)
	}

	return reserve, true
}

func (k Keeper) SetSubToken(ctx sdk.Context, denom, id string, subTokenID int64, reserve sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(denom, id, subTokenID)

	bz, err := reserve.Marshal()
	if err != nil {
		panic(err)
	}

	store.Set(subTokenKey, bz)
}

func (k Keeper) RemoveSubToken(ctx sdk.Context, denom, id string, subTokenID int64) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(denom, id, subTokenID)
	store.Delete(subTokenKey)
}

func (k Keeper) GetLastSubTokenID(ctx sdk.Context, denom, id string) int64 {
	store := ctx.KVStore(k.storeKey)
	lastSubTokenIDKey := types.GetLastSubTokenIDKey(denom, id)

	bz := store.Get(lastSubTokenIDKey)
	if bz == nil {
		return 0
	}

	return int64(binary.LittleEndian.Uint64(bz))
}

func (k Keeper) SetLastSubTokenID(ctx sdk.Context, denom, id string, lastSubTokenID int64) {
	store := ctx.KVStore(k.storeKey)
	lastSubTokenIDKey := types.GetLastSubTokenIDKey(denom, id)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, uint64(lastSubTokenID))
	store.Set(lastSubTokenIDKey, bz)
}

func (k Keeper) SetTokenURI(ctx sdk.Context, tokenURI string) {
	store := ctx.KVStore(k.storeKey)
	tokenURIKey := types.GetTokenURIKey(tokenURI)

	store.Set(tokenURIKey, []byte{})
}

func (k Keeper) ExistTokenURI(ctx sdk.Context, tokenURI string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenURIKey := types.GetTokenURIKey(tokenURI)

	return store.Has(tokenURIKey)
}

func (k Keeper) SetTokenIDIndex(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	tokenIDKey := types.GetTokenIDKey(id)

	store.Set(tokenIDKey, []byte{})
}

func (k Keeper) ExistTokenID(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenIDKey := types.GetTokenIDKey(id)

	return store.Has(tokenIDKey)
}

// mintNFT mints an NFT and manages that NFTs existence within Collections and Owners
func (k Keeper) mintNFT(
	ctx sdk.Context,
	denom, id string,
	reserve, quantity sdk.Int,
	creator, owner string,
	tokenURI string,
	allowMint bool,
) (int64, error) {
	// TODO abcd move from here
	nft, err := k.GetNFT(ctx, denom, id)
	if err == nil {
		reserve = nft.GetReserve()
	}

	// TODO last sub token хранит не last, а следующий достпный для сохранения саб токен id
	lastSubTokenID := k.GetLastSubTokenID(ctx, denom, id)

	if lastSubTokenID == 0 {
		lastSubTokenID = 1
	}

	tempSubTokenID := lastSubTokenID
	subTokenIDs := make([]int64, quantity.Int64())
	for i := int64(0); i < quantity.Int64(); i++ {
		subTokenIDs[i] = tempSubTokenID
		tempSubTokenID++
	}

	nft = types.NewBaseNFT(id, creator, owner, tokenURI, reserve, subTokenIDs, allowMint)

	collection, found := k.GetCollection(ctx, denom)
	if found {
		collection, err = collection.AddNFT(nft)
		if err != nil {
			return 0, err
		}
	} else {
		collection = types.NewCollection(denom, types.NewNFTs(nft.(types.BaseNFT)))
	}
	k.SetCollection(ctx, denom, collection)

	k.SetTokenIDIndex(ctx, id)

	newLastSubTokenID := lastSubTokenID + quantity.Int64()

	for i := lastSubTokenID; i < newLastSubTokenID; i++ {
		k.SetSubToken(ctx, denom, nft.GetID(), i, nft.GetReserve())
	}

	k.SetLastSubTokenID(ctx, denom, nft.GetID(), newLastSubTokenID)

	creatorAddress, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return 0, err
	}

	reservedCoin := sdk.NewCoin(
		*k.BaseDenom,
		reserve.Mul(quantity), // reserve * quantity
	)

	err = k.ReserveTokens(ctx, sdk.NewCoins(reservedCoin), creatorAddress)
	if err != nil {
		return 0, err
	}

	ownerIDCollection, _ := k.GetIDCollectionByDenom(ctx, creatorAddress, denom)
	// TODO abcd Зачем сохранять одинаковые айди списком? IDs: [1, 1, 1, 1]. Мб лучше [1] ?
	ownerIDCollection = ownerIDCollection.AddID(nft.GetID())
	k.SetIDCollectionByDenom(ctx, creatorAddress, denom, ownerIDCollection.IDs)

	return newLastSubTokenID, err
}

// UpdateNFT updates an already existing NFTs
func (k Keeper) UpdateNFT(ctx sdk.Context, denom string, nft exported.NFT) (err error) {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return types.ErrUnknownCollection(denom)
	}

	oldNFT, err := collection.GetNFT(nft.GetID())
	if err != nil {
		return err
	}

	collection.NFTs, _ = collection.NFTs.Update(oldNFT.GetID(), nft)

	k.SetCollection(ctx, denom, collection)
	return nil
}

// EditNFT edits an existing NFT meta info
func (k Keeper) EditNFT(ctx sdk.Context, denom, id string, tokenURI string) error {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return types.ErrUnknownCollection(denom)
	}

	nft, err := collection.GetNFT(id)
	if err != nil {
		return err
	}

	nft = nft.EditMetadata(tokenURI)

	collection, err = collection.UpdateNFT(nft)
	if err != nil {
		return err
	}

	k.SetCollection(ctx, denom, collection)

	return nil
}

// DeleteNFT deletes an existing NFT from store
func (k Keeper) DeleteNFT(ctx sdk.Context, denom, id string, subTokenIDsToDelete []int64) error {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return types.ErrUnknownCollection(denom)
	}

	nft, err := collection.GetNFT(id)
	if err != nil {
		return err
	}

	reserveForReturn := sdk.ZeroInt()

	owner := nft.GetOwners().GetOwner(nft.GetCreator())
	if owner == nil {
		return types.ErrNotAllowedBurn()
	}

	ownerSubTokenIDs := types.SortedIntArray(owner.GetSubTokenIDs())
	for _, subTokenIDToDelete := range subTokenIDsToDelete {
		if ownerSubTokenIDs.Find(subTokenIDToDelete) == -1 {

			return sdkerrors.Wrap(
				types.ErrNotAllowedBurn(),
				fmt.Sprintf(
					"owner %s has only %s tokens", nft.GetCreator(),
					types.SortedIntArray(nft.GetOwners().GetOwner(nft.GetCreator()).GetSubTokenIDs()).String(),
				),
			)
		}
		owner = owner.RemoveSubTokenID(subTokenIDToDelete)
		reserve, ok := k.GetSubToken(ctx, denom, id, subTokenIDToDelete)
		if !ok {
			return fmt.Errorf("subToken with ID = %d not found", subTokenIDToDelete)
		}
		reserveForReturn = reserveForReturn.Add(reserve)
		k.RemoveSubToken(ctx, denom, id, subTokenIDToDelete)
	}

	tokenOwners, err := nft.GetOwners().SetOwner(owner)
	if err != nil {
		return err
	}

	nft = nft.SetOwners(tokenOwners)

	collection, err = collection.UpdateNFT(nft)
	if err != nil {
		return err
	}

	k.SetCollection(ctx, denom, collection)

	ownerAddress, err := sdk.AccAddressFromBech32(owner.GetAddress())
	if err != nil {
		return err
	}

	coinsToReturn := sdk.NewCoins(sdk.NewCoin(*k.BaseDenom, reserveForReturn))

	err = k.ReturnTokensTo(ctx, coinsToReturn, ownerAddress)
	if err != nil {
		return err
	}

	return nil
}

//UpdateNFTReserve function to increase the minimum reserve of the NFT token
func (k Keeper) UpdateNFTReserve(ctx sdk.Context, denom, id string, subTokenIDs []int64, newReserve sdk.Int) error {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return types.ErrUnknownCollection(denom)
	}

	nft, err := collection.GetNFT(id)
	if err != nil {
		return err
	}

	owner := nft.GetOwners().GetOwner(nft.GetCreator())
	ownerSubTokenIDs := types.SortedIntArray(owner.GetSubTokenIDs())

	reserveForRefill := sdk.NewInt(0)

	for _, subTokenID := range subTokenIDs {
		if ownerSubTokenIDs.Find(subTokenID) == -1 {
			return sdkerrors.Wrap(types.ErrNotAllowedUpdateReserve(),
				fmt.Sprintf("owner %s has only %s tokens", nft.GetCreator(),
					types.SortedIntArray(nft.GetOwners().GetOwner(nft.GetCreator()).GetSubTokenIDs()).String()))
		}
		reserve, _ := k.GetSubToken(ctx, denom, id, subTokenID)
		if reserve.Equal(newReserve) {
			return types.ErrNotSetValueLowerNow()
		}

		if reserve.GT(newReserve) {
			return types.ErrNotSetValueLowerNow()
		}

		reserveForRefill = reserveForRefill.Add(newReserve.Sub(reserve))

		k.SetSubToken(ctx, denom, id, subTokenID, newReserve)
	}

	ownerAddress, err := sdk.AccAddressFromBech32(owner.GetAddress())
	if err != nil {
		return err
	}

	coinsToReserve := sdk.NewCoins(sdk.NewCoin(*k.BaseDenom, reserveForRefill))
	err = k.ReserveTokens(ctx, coinsToReserve, ownerAddress)
	if err != nil {
		return types.ErrNotEnoughFunds(reserveForRefill.String())
	}

	return err
}
