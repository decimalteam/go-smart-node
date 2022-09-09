package keeper_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestQueryCollectionSupply(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	collectionDenom := firstDenom

	firstNFT := types.NewBaseNFT(
		firstID,
		sender.String(),
		firstTokenURI,
		firstReserve,
		true,
	)

	err := mintNFT(collectionDenom, 10, firstNFT, dsc, ctx)
	require.NoError(t, err)

	err = mintNFT(collectionDenom, 10, firstNFT, dsc, ctx)
	require.NoError(t, err)

	secondNFT := types.NewBaseNFT(
		secondID,
		sender.String(),
		secondTokenURI,
		firstReserve,
		true,
	)

	err = mintNFT(collectionDenom, 10, secondNFT, dsc, ctx)
	require.NoError(t, err)

	expectedResponse := types.QueryCollectionSupplyResponse{
		Supply: 2,
	}

	req := types.QueryCollectionSupplyRequest{
		Denom: collectionDenom,
	}

	res, err := dsc.NFTKeeper.QueryCollectionSupply(sdk.WrapSDKContext(ctx), &req)
	require.NoError(t, err)
	require.Equal(t, expectedResponse, *res)

	msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, msg)
}

func TestQueryOwnerCollections(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	firstCollectionDenom := firstDenom
	firstNFT := types.NewBaseNFT(
		firstID,
		sender.String(),
		firstTokenURI,
		firstReserve,
		true,
	)
	err := mintNFT(firstCollectionDenom, 10, firstNFT, dsc, ctx)
	require.NoError(t, err)

	secondNFT := types.NewBaseNFT(
		secondID,
		sender.String(),
		secondTokenURI,
		firstReserve,
		true,
	)
	err = mintNFT(firstCollectionDenom, 10, secondNFT, dsc, ctx)
	require.NoError(t, err)

	thirdNFT := types.NewBaseNFT(
		thirdID,
		sender.String(),
		thirdTokenURI,
		firstReserve,
		true,
	)
	secondCollectionDenom := secondDenom
	err = mintNFT(secondCollectionDenom, 10, thirdNFT, dsc, ctx)
	require.NoError(t, err)

	msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, msg)

	t.Run("without denom", func(t *testing.T) {
		expectedResponse := types.QueryOwnerCollectionsResponse{
			Owner: types.Owner{
				Address: sender.String(),
				Collections: []types.OwnerCollection{
					{
						Denom: firstCollectionDenom,
						NFTs:  types.SortedStringArray{firstNFT.ID, secondNFT.ID},
					},
					{
						Denom: secondCollectionDenom,
						NFTs:  types.SortedStringArray{thirdNFT.ID},
					},
				},
			},
		}

		req := types.QueryOwnerCollectionsRequest{
			Owner: sender.String(),
		}

		res, err := dsc.NFTKeeper.QueryOwnerCollections(sdk.WrapSDKContext(ctx), &req)
		require.NoError(t, err)
		require.Equal(t, expectedResponse, *res)
	})

}

func TestQueryCollection(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	firstCollectionDenom := firstDenom
	firstNFT := types.NewBaseNFT(
		firstID,
		sender.String(),
		firstTokenURI,
		firstReserve,
		true,
	)
	err := mintNFT(firstCollectionDenom, 10, firstNFT, dsc, ctx)
	require.NoError(t, err)

	secondNFT := types.NewBaseNFT(
		secondID,
		sender.String(),
		secondTokenURI,
		firstReserve,
		true,
	)
	err = mintNFT(firstCollectionDenom, 10, secondNFT, dsc, ctx)
	require.NoError(t, err)

	expectedResponse := types.QueryCollectionResponse{
		Collection: types.Collection{
			Denom: firstCollectionDenom,
			NFTs:  types.SortedStringArray{firstNFT.ID, secondNFT.ID},
		},
	}

	req := types.QueryCollectionRequest{
		Denom: firstCollectionDenom,
	}

	res, err := dsc.NFTKeeper.QueryCollection(sdk.WrapSDKContext(ctx), &req)
	require.NoError(t, err)
	require.Equal(t, expectedResponse, *res)

	msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, msg)
}

func TestQueryNFT(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	collectionDenom := firstDenom
	nft := types.NewBaseNFT(
		firstID,
		sender.String(),
		firstTokenURI,
		firstReserve,
		true,
	)

	subTokenIDs := nft.GenSubTokenIDs(5)
	nft = nft.AddOwnerSubTokenIDs(sender.String(), subTokenIDs)

	err := mintNFT(collectionDenom, int64(subTokenIDs.Len()), nft, dsc, ctx)
	require.NoError(t, err)

	expectedResponse := types.QueryNFTResponse{
		NFT: nft,
	}

	req := types.QueryNFTRequest{
		Denom:   collectionDenom,
		TokenId: nft.ID,
	}

	res, err := dsc.NFTKeeper.QueryNFT(sdk.WrapSDKContext(ctx), &req)
	require.NoError(t, err)
	require.Equal(t, expectedResponse, *res)

	msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, msg)
}

func TestQuerySubTokens(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	collectionDenom := firstDenom
	nft := types.NewBaseNFT(
		firstID,
		sender.String(),
		firstTokenURI,
		firstReserve,
		true,
	)

	subTokenIDs := nft.GenSubTokenIDs(5)

	err := mintNFT(collectionDenom, int64(subTokenIDs.Len()), nft, dsc, ctx)
	require.NoError(t, err)

	subTokens := make([]types.SubToken, subTokenIDs.Len())
	for i, subTokenID := range subTokenIDs {
		subTokens[i] = types.NewSubToken(subTokenID, nft.Reserve)
	}

	expectedResponse := types.QuerySubTokensResponse{
		SubTokens: subTokens,
	}

	req := types.QuerySubTokensRequest{
		Denom:   collectionDenom,
		TokenID: nft.ID,
	}

	res, err := dsc.NFTKeeper.QuerySubTokens(sdk.WrapSDKContext(ctx), &req)
	require.NoError(t, err)
	require.Equal(t, expectedResponse, *res)

	msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, msg)
}
