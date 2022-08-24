package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"fmt"
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey store.StoreKey // Unexposed key to access store from sdk.Context

	cdc codec.BinaryCodec // The amino codec for binary encoding/decoding.

	bankKeeper keeper.Keeper

	BaseDenom *string
}

// NewKeeper creates new instances of the nft Keeper
func NewKeeper(cdc codec.BinaryCodec, storeKey store.StoreKey, bankKeeper keeper.Keeper, baseDenom string) *Keeper {
	return &Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		bankKeeper: bankKeeper,
		BaseDenom:  &baseDenom,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Mint(
	ctx sdk.Context,
	denom, id string,
	reserve sdk.Coin,
	quantity sdkmath.Int,
	creator, owner string,
	tokenURI string,
	allowMint bool,
) ([]uint64, error) {
	subTokenIDs, err := k.GenAndMintSubTokens(ctx, denom, id, reserve, quantity, creator)
	if err != nil {
		return []uint64{}, err
	}

	err = k.MintNFTAndCollection(ctx, denom, id, reserve, creator, owner, tokenURI, allowMint, subTokenIDs)
	if err != nil {
		return []uint64{}, err
	}

	return subTokenIDs, nil
}

func (k Keeper) GenAndMintSubTokens(
	ctx sdk.Context,
	denom, id string,
	reserve sdk.Coin,
	quantity sdkmath.Int,
	creator string,
) (types.SortedUintArray, error) {
	nft, err := k.GetNFT(ctx, denom, id)
	if err == nil {
		reserve = nft.GetReserve()
	}

	subTokenIDs := nft.GenSubTokenIDs(quantity.Uint64())
	for _, subTokenID := range subTokenIDs {
		subToken := types.SubToken{
			ID:      subTokenID,
			Reserve: reserve,
		}

		k.SetSubToken(ctx, id, subToken)
	}

	creatorAddress, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return nil, err
	}

	coinToReserve := sdk.NewCoin(
		reserve.Denom,
		reserve.Amount.Mul(quantity), // reserve * quantity
	)

	err = k.ReserveTokens(ctx, sdk.NewCoins(coinToReserve), creatorAddress)
	if err != nil {
		return nil, err
	}

	return subTokenIDs, nil
}

func (k Keeper) MintSubTokens(
	ctx sdk.Context,
	id string,
	subTokens []types.SubToken,
) {
	for _, subToken := range subTokens {
		k.SetSubToken(ctx, id, subToken)
	}
}

// MintNFTAndCollection mints an NFT and manages that NFTs existence within Collections and Owners
func (k Keeper) MintNFTAndCollection(
	ctx sdk.Context,
	denom, id string,
	reserve sdk.Coin,
	creator, owner string,
	tokenURI string,
	allowMint bool,
	subTokenIDs []uint64,
) error {
	nft, err := k.GetNFT(ctx, denom, id)
	if err == nil {
		// add sub tokens
		reserve = nft.GetReserve()
	} else {
		nft = types.NewBaseNFT(id, creator, tokenURI, reserve, allowMint)
	}

	nft = nft.AddOwnerSubTokenIDs(owner, subTokenIDs)

	collection, found := k.GetCollection(ctx, denom)
	if found {
		collection = collection.AddNFT(nft.ID)
	} else {
		collection = types.NewCollection(denom, []string{nft.ID})
	}

	k.SetCollection(ctx, denom, collection)

	err = k.SetNFT(ctx, denom, id, nft)
	if err != nil {
		return err
	}

	creatorAddress, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return err
	}

	ownerCollection, found := k.GetOwnerCollectionByDenom(ctx, creatorAddress, denom)
	if !found {
		ownerCollection = types.NewOwnerCollection(denom, []string{nft.GetID()})
	} else {
		ownerCollection = ownerCollection.AddID(nft.GetID())
	}

	k.SetOwnerCollectionByDenom(ctx, creatorAddress, denom, ownerCollection)

	return nil
}

func (k Keeper) Transfer(ctx sdk.Context, denom, id string, sender, recipient string, subTokenIDsToTransfer []uint64) (types.BaseNFT, error) {
	nft, err := k.GetNFT(ctx, denom, id)
	if err != nil {
		return types.BaseNFT{}, err
	}

	senderOwner := nft.GetOwners().GetOwner(sender)

	for _, idToTransfer := range subTokenIDsToTransfer {
		if senderOwner.GetSubTokenIDs().Find(idToTransfer) == -1 {
			return types.BaseNFT{}, errors.OwnerDoesNotOwnSubTokenID
		}
		senderOwner = senderOwner.RemoveSubTokenID(idToTransfer)
	}
	nft = nft.SetOwners(nft.GetOwners().SetOwner(senderOwner))

	recipientOwner := nft.GetOwners().GetOwner(recipient)
	if recipientOwner == nil {
		recipientOwner = types.NewTokenOwner(recipient, subTokenIDsToTransfer)
	} else {
		for _, subTokenID := range subTokenIDsToTransfer {
			recipientOwner = recipientOwner.SetSubTokenID(subTokenID)
		}
	}
	nft = nft.SetOwners(nft.GetOwners().SetOwner(recipientOwner))

	err = k.SetNFT(ctx, denom, id, nft)
	if err != nil {
		return types.BaseNFT{}, err
	}

	return nft, nil
}

// EditNFT edits an existing NFT meta info
func (k Keeper) EditNFT(ctx sdk.Context, denom, id string, tokenURI string) error {
	nft, err := k.GetNFT(ctx, denom, id)
	if err != nil {
		return err
	}

	nft = nft.EditMetadata(tokenURI)

	err = k.SetNFT(ctx, denom, id, nft)
	if err != nil {
		return err
	}

	return nil
}

// DeleteNFTSubTokens deletes an NFT sub tokens from store
func (k Keeper) DeleteNFTSubTokens(ctx sdk.Context, denom, id string, subTokenIDsToDelete []uint64) error {
	nft, err := k.GetNFT(ctx, denom, id)
	if err != nil {
		return err
	}

	reserveForReturn := sdk.ZeroInt()

	owner := nft.GetOwners().GetOwner(nft.GetCreator())
	if owner == nil {
		return errors.NotAllowedBurn
	}

	ownerSubTokenIDs := owner.GetSubTokenIDs()
	for _, subTokenIDToDelete := range subTokenIDsToDelete {
		if ownerSubTokenIDs.Find(subTokenIDToDelete) == -1 {

			return errors.NotAllowedBurn
		}
		owner = owner.RemoveSubTokenID(subTokenIDToDelete)

		subToken, ok := k.GetSubToken(ctx, id, subTokenIDToDelete)
		if !ok {
			return errors.SubTokenDoesNotExists
		}
		reserveForReturn = reserveForReturn.Add(subToken.Reserve.Amount)
		k.RemoveSubToken(ctx, id, subTokenIDToDelete)
	}

	nft = nft.SetOwners(nft.GetOwners().SetOwner(owner))

	err = k.SetNFT(ctx, denom, id, nft)
	if err != nil {
		return err
	}

	ownerAddress, err := sdk.AccAddressFromBech32(owner.GetAddress())
	if err != nil {
		return err
	}

	coinsToReturn := sdk.NewCoins(
		sdk.NewCoin(nft.Reserve.Denom, reserveForReturn),
	)

	err = k.ReturnTokensTo(ctx, coinsToReturn, ownerAddress)
	if err != nil {
		return err
	}

	return nil
}

// UpdateNFTReserve increases the minimum reserve of the NFT token
func (k Keeper) UpdateNFTReserve(ctx sdk.Context, denom, id string, subTokenIDs []uint64, newReserve sdk.Coin) error {
	nft, err := k.GetNFT(ctx, denom, id)
	if err != nil {
		return err
	}

	if nft.Reserve.Denom != newReserve.Denom {
		return errors.WrongReserveCoinDenom
	}

	owner := nft.GetOwners().GetOwner(nft.GetCreator())
	ownerSubTokenIDs := owner.GetSubTokenIDs()

	var reserveForRefill sdk.Coin

	for _, subTokenID := range subTokenIDs {
		if ownerSubTokenIDs.Find(subTokenID) == -1 {
			return errors.NotAllowedUpdateReserve
		}
		subToken, _ := k.GetSubToken(ctx, id, subTokenID)
		if subToken.Reserve.Equal(newReserve) {
			return errors.NotSetValueLowerNow
		}

		if newReserve.IsLT(subToken.Reserve) {
			return errors.NotSetValueLowerNow
		}

		reserveForRefill = newReserve.Sub(subToken.Reserve)

		k.SetSubToken(ctx, id, types.SubToken{
			ID:      subTokenID,
			Reserve: newReserve,
		})
	}

	ownerAddress, err := sdk.AccAddressFromBech32(owner.GetAddress())
	if err != nil {
		return err
	}

	coinsToReserve := sdk.NewCoins(reserveForRefill)
	err = k.ReserveTokens(ctx, coinsToReserve, ownerAddress)
	if err != nil {
		return errors.InsufficientFunds
	}

	return err
}
