package nft_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// nolint: deadcode unused
const (
	firstDenom  = "first_denom"
	secondDenom = "second_denom"
	thirdDenom  = "third_denom"
)

func TestInitGenesis(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	genesisState := types.DefaultGenesisState()
	require.Equal(t, 0, len(genesisState.Collections))
	require.Equal(t, 0, len(genesisState.NFTs))
	require.Equal(t, 0, len(genesisState.SubTokens))

	denoms := []string{firstDenom, secondDenom, thirdDenom}
	const nftsAmount = 2

	// prepare collections
	collections := make([]types.Collection, len(denoms))
	for i, denom := range denoms {
		nftIDs := make([]string, nftsAmount)
		for j := 0; j < nftsAmount; j++ {
			nftIDs[j] = time.Now().String()
		}

		collections[i] = types.NewCollection(denom, nftIDs)
	}

	const subTokenAmount = 3
	// prepare nfts
	nfts := make([]types.BaseNFT, len(denoms)*nftsAmount)
	subTokens := make(map[string]types.SubTokens, len(nfts))
	for i, collection := range collections {
		for j, nftID := range collection.NFTs {
			address := app.GetAddrs(dsc, ctx, 1)[0]
			tokenURI := time.Now().String()
			baseNFT := types.NewBaseNFT(
				nftID,
				address.String(),
				tokenURI,
				sdk.NewCoin(config.BaseDenom, types.MinReserve),
				true,
			)

			subTokenIDs := baseNFT.GenSubTokenIDs(subTokenAmount)
			firstTokenOwner := types.NewTokenOwner(address.String(), subTokenIDs)
			baseNFT = baseNFT.SetOwners(baseNFT.GetOwners().SetOwner(firstTokenOwner))

			nfts[i*nftsAmount+j] = baseNFT

			// prepare sub tokens map
			subTokens[baseNFT.ID] = types.SubTokens{
				SubTokens: make([]types.SubToken, subTokenIDs.Len()),
			}
			for q, subTokenID := range subTokenIDs {
				subTokens[baseNFT.ID].SubTokens[q] = types.NewSubToken(subTokenID, baseNFT.Reserve)
			}
		}
	}

	genesisState = types.NewGenesisState(collections, nfts, subTokens)
	err := genesisState.Validate()
	require.NoError(t, err)

	nft.InitGenesis(ctx, dsc.NFTKeeper, *genesisState)

	storedCollections := dsc.NFTKeeper.GetCollections(ctx)
	compareCollections(t, storedCollections, collections)

	storedNFTs, err := dsc.NFTKeeper.GetNFTs(ctx)
	require.NoError(t, err)

	compareNFTs(t, storedNFTs, nfts)

	exportedGenesisState := nft.ExportGenesis(ctx, dsc.NFTKeeper)
	err = exportedGenesisState.Validate()
	require.NoError(t, err)

	require.Len(t, exportedGenesisState.Collections, len(collections))
	compareCollections(t, exportedGenesisState.Collections, collections)

	require.Len(t, exportedGenesisState.NFTs, len(nfts))
	compareNFTs(t, exportedGenesisState.NFTs, nfts)

	require.Equal(t, exportedGenesisState.SubTokens, subTokens)
}

func compareCollections(t *testing.T, storedCollections, collectionsToStore []types.Collection) {
	require.Len(t, storedCollections, len(collectionsToStore))
	for _, storedCollection := range storedCollections {
		var matched bool
		for _, collectionToStore := range collectionsToStore {
			if collectionToStore.Denom == storedCollection.Denom {
				matched = true
				require.Equal(t, collectionToStore, storedCollection)
				break
			}
		}

		require.True(t, matched, storedCollection.Denom)
	}
}

func compareNFTs(t *testing.T, storedNFTs, nftsToStore []types.BaseNFT) {
	require.Len(t, storedNFTs, len(nftsToStore))
	for _, storedNFT := range storedNFTs {
		var matched bool
		for _, nftToStore := range nftsToStore {
			if nftToStore.ID == storedNFT.ID {
				matched = true
				require.Equal(t, nftToStore, storedNFT)
				break
			}
		}

		require.True(t, matched, storedNFT.ID)
	}
}
