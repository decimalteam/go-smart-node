package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"
	"testing"
)

func TestMintNFTAndCollection(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addr := app.GetAddrs(dsc, ctx, 1)[0]

	nft := types.NewBaseNFT(
		firstID,
		addr.String(),
		firstTokenURI,
		firstReserve,
		firstAllowMint,
	)
	subTokenIDs := []uint64{1, 2, 3}
	nft = nft.AddOwnerSubTokenIDs(nft.Creator, subTokenIDs)

	expectedCollection := types.NewCollection(firstDenom, []string{nft.ID})

	expectedOwnerCollection := types.NewOwnerCollection(expectedCollection.Denom, []string{nft.ID})

	err := dsc.NFTKeeper.MintNFTAndCollection(ctx, expectedCollection.Denom, nft.ID, nft.Reserve, nft.Creator, nft.Creator, nft.TokenURI, nft.AllowMint, subTokenIDs)
	require.NoError(t, err)

	storedCollection, found := dsc.NFTKeeper.GetCollection(ctx, expectedCollection.Denom)
	require.True(t, found)
	require.Equal(t, expectedCollection, storedCollection)

	storedNFT, err := dsc.NFTKeeper.GetNFT(ctx, expectedCollection.Denom, nft.ID)
	require.NoError(t, err)
	require.Equal(t, nft, storedNFT)

	storedOwnerCollection, found := dsc.NFTKeeper.GetOwnerCollectionByDenom(ctx, addr, expectedOwnerCollection.Denom)
	require.True(t, found)
	require.Equal(t, expectedOwnerCollection, storedOwnerCollection)
}

func TestMintExistingNFT(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	testCases := genNFTs(2, dsc, ctx)
	for _, tc := range testCases {
		fmt.Println("ITER", tc)

		nft := tc.nft

		err := dsc.NFTKeeper.MintNFTAndCollection(ctx, firstDenom, nft.ID, nft.Reserve, nft.Creator, nft.Creator, nft.TokenURI, nft.AllowMint, tc.newSubTokens)
		require.NoError(t, err)

		storedCollection, found := dsc.NFTKeeper.GetCollection(ctx, firstDenom)
		require.True(t, found)
		require.Equal(t, tc.expectedCollection, storedCollection)

		storedNFT, err := dsc.NFTKeeper.GetNFT(ctx, firstDenom, nft.ID)
		require.NoError(t, err)
		require.Equal(t, nft, storedNFT)

		storedOwnerCollection, found := dsc.NFTKeeper.GetOwnerCollectionByDenom(ctx, tc.newOwner, firstDenom)
		require.True(t, found)
		require.Equal(t, tc.expectedOwnerCollection, storedOwnerCollection)
	}
}

type testCase struct {
	nft                     types.BaseNFT
	newOwner                sdk.AccAddress
	newSubTokens            []uint64
	expectedCollection      types.Collection
	expectedOwnerCollection types.OwnerCollection
}

func genNFTs(amount int, dsc *app.DSC, ctx sdk.Context) []testCase {

	nftStates := make(map[string]types.BaseNFT)
	onwerCollectionStates := make(map[string][]string)
	var createdNFTIDs []string
	testCases := make([]testCase, amount)
	creator := app.GetAddrs(dsc, ctx, 1)[0]
	owners := append(app.GetAddrs(dsc, ctx, 1), creator)

	for i := 0; i < amount; i++ {
		nftID := fmt.Sprintf("%d", rand.Intn(3))
		ownerID := owners[rand.Intn(len(owners)-1)]

		nft, ok := nftStates[nftID]
		if !ok {
			nft = types.NewBaseNFT(
				nftID,
				creator.String(),
				firstTokenURI,
				firstReserve,
				firstAllowMint,
			)

			onwerCollectionStates[ownerID.String()] = append(onwerCollectionStates[ownerID.String()], nftID)
			createdNFTIDs = append(createdNFTIDs, nft.ID)
		}

		subTokenIDs := nft.GenSubTokenIDs(uint64(rand.Intn(3) + 1))

		fmt.Println(subTokenIDs)
		nft = nft.AddOwnerSubTokenIDs(ownerID.String(), subTokenIDs)

		nftStates[nftID] = nft

		testCases[i] = testCase{
			nft:                     nft,
			newOwner:                ownerID,
			newSubTokens:            subTokenIDs,
			expectedCollection:      types.NewCollection(firstDenom, types.SortedStringArray(createdNFTIDs).Sort()),
			expectedOwnerCollection: types.NewOwnerCollection(firstDenom, onwerCollectionStates[ownerID.String()]),
		}
	}

	fmt.Println(testCases[0].nft)
	fmt.Println(testCases[1].nft)

	return testCases
}
