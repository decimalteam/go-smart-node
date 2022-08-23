package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetNFT(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	addrs := app.GetAddrs(dsc, ctx, 2)

	nft1 := types.NewBaseNFT(
		firstID,
		addrs[0].String(),
		firstTokenURI,
		firstReserve,
		true,
	)

	nft1 = nft1.AddOwnerSubTokenIDs(addrs[0].String(), []uint64{1, 2, 3})

	nft2 := types.NewBaseNFT(
		secondID,
		addrs[0].String(),
		firstTokenURI,
		firstReserve,
		true,
	)

	nft2 = nft2.AddOwnerSubTokenIDs(addrs[1].String(), []uint64{4, 5, 6})

	// SetNFT method must return an error if there is no collection for this denom
	err := dsc.NFTKeeper.SetNFT(ctx, firstDenom, firstID, nft1)
	require.Error(t, err, firstID)

	nftsToStore := []types.BaseNFT{nft1, nft2}
	nftIDsToStore := make([]string, len(nftsToStore))
	for i, nft := range nftsToStore {
		nftIDsToStore[i] = nft.GetID()
	}

	collection := types.NewCollection(firstDenom, nftIDsToStore)
	dsc.NFTKeeper.SetCollection(ctx, collection.Denom, collection)

	for _, nft := range nftsToStore {
		err := dsc.NFTKeeper.SetNFT(ctx, collection.GetDenom(), nft.GetID(), nft)
		require.NoError(t, err, nft.GetID())
	}

	// Check throw GetNFTs method
	storedNFTs, err := dsc.NFTKeeper.GetNFTs(ctx)
	require.NoError(t, err)

	require.Len(t, nftsToStore, len(storedNFTs))

	for _, nftToStore := range nftsToStore {
		var stored bool
		for _, storedNFT := range storedNFTs {
			if nftToStore.GetID() == storedNFT.GetID() {
				stored = true
				require.Equal(t, nftToStore, storedNFT)
				break
			}
		}

		require.True(t, stored, nftToStore.GetID())
	}

	// Check throw GetNFT method
	for _, nftToStore := range nftsToStore {
		storedNFT, err := dsc.NFTKeeper.GetNFT(ctx, collection.GetDenom(), nftToStore.GetID())
		require.NoError(t, err, nftToStore.GetID())
		require.Equal(t, nftToStore, storedNFT)
	}
}

func TestHasTokenID(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	addrs := app.GetAddrs(dsc, ctx, 1)

	nft1 := types.NewBaseNFT(
		firstID,
		addrs[0].String(),
		firstTokenURI,
		firstReserve,
		true,
	)

	exists := dsc.NFTKeeper.HasTokenID(ctx, nft1.GetID())
	require.False(t, exists, nft1.GetID())

	// Сan not save nft without collection
	collection := types.NewCollection(firstDenom, []string{nft1.GetID()})
	dsc.NFTKeeper.SetCollection(ctx, collection.Denom, collection)

	err := dsc.NFTKeeper.SetNFT(ctx, collection.Denom, nft1.GetID(), nft1)
	require.NoError(t, err, nft1.GetID())

	exists = dsc.NFTKeeper.HasTokenID(ctx, nft1.GetID())
	require.True(t, exists, nft1.GetID())
}

func TestHasTokenURI(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	addrs := app.GetAddrs(dsc, ctx, 1)

	nft1 := types.NewBaseNFT(
		firstID,
		addrs[0].String(),
		firstTokenURI,
		firstReserve,
		true,
	)

	exists := dsc.NFTKeeper.HasTokenURI(ctx, nft1.GetTokenURI())
	require.False(t, exists, nft1.GetTokenURI())

	// Сan not save nft without collection
	collection := types.NewCollection(firstDenom, []string{nft1.GetID()})
	dsc.NFTKeeper.SetCollection(ctx, collection.Denom, collection)

	err := dsc.NFTKeeper.SetNFT(ctx, collection.Denom, nft1.GetID(), nft1)
	require.NoError(t, err, nft1.GetID())

	exists = dsc.NFTKeeper.HasTokenURI(ctx, nft1.GetTokenURI())
	require.True(t, exists, nft1.GetTokenURI())
}
