package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMintNFT(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	expectedNFT := types.NewBaseNFT(
		firstID,
		sender.String(),
		firstTokenURI,
		firstReserve,
		true,
	)
	subTokenIDs := expectedNFT.GenSubTokenIDs(10)
	expectedNFT = expectedNFT.AddOwnerSubTokenIDs(expectedNFT.Creator, subTokenIDs)

	expectedCollection := types.NewCollection(firstDenom, []string{expectedNFT.ID})
	expectedOwnerCollections := []types.OwnerCollection{
		types.NewOwnerCollection(expectedCollection.Denom, []string{expectedNFT.ID}),
	}

	expectedSubTokens := make([]types.SubToken, subTokenIDs.Len())
	for i, subTokenID := range subTokenIDs {
		expectedSubTokens[i] = types.NewSubToken(subTokenID, expectedNFT.Reserve)
	}

	err := mintNFT(expectedCollection.Denom, int64(len(expectedSubTokens)), expectedNFT, dsc, ctx)
	require.NoError(t, err)

	storedCollection, found := dsc.NFTKeeper.GetCollection(ctx, expectedCollection.Denom)
	require.True(t, found)
	require.Equal(t, expectedCollection, storedCollection)

	storedNFT, err := dsc.NFTKeeper.GetNFT(ctx, expectedCollection.Denom, expectedNFT.ID)
	require.NoError(t, err)
	require.Equal(t, expectedNFT, storedNFT)

	storedOwnerCollections, err := dsc.NFTKeeper.GetOwnerCollections(ctx, sender)
	require.NoError(t, err)

	require.Len(t, storedOwnerCollections, len(expectedOwnerCollections))
	require.Equal(t, expectedOwnerCollections, storedOwnerCollections)

	storedSubTokens, err := dsc.NFTKeeper.GetSubTokens(ctx, expectedNFT.ID)
	require.NoError(t, err)

	require.Len(t, storedSubTokens, len(expectedSubTokens))
	require.Equal(t, expectedSubTokens, storedSubTokens)

	msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, msg)
}

func TestMintNFTValidation(t *testing.T) {
	t.Run("disabled mint", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		sender := app.GetAddrs(dsc, ctx, 1)[0]

		nft := types.NewBaseNFT(
			firstID,
			sender.String(),
			firstTokenURI,
			firstReserve,
			false,
		)

		// first nft mint with allowMint=false flag should be without error
		err := mintNFT(firstDenom, 10, nft, dsc, ctx)
		require.NoError(t, err)

		// next nft mint with allowMint=false flag should be with error
		err = mintNFT(firstDenom, 10, nft, dsc, ctx)
		require.Error(t, err)

		msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
		require.False(t, broken, msg)
	})

	t.Run("wrong sender account", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		addrs := app.GetAddrs(dsc, ctx, 2)
		creator := addrs[0]
		wrongCreator := addrs[1]

		firstNFT := types.NewBaseNFT(
			firstID,
			creator.String(),
			firstTokenURI,
			firstReserve,
			true,
		)

		err := mintNFT(firstDenom, 10, firstNFT, dsc, ctx)
		require.NoError(t, err)

		secondNFT := types.NewBaseNFT(
			firstNFT.ID,
			wrongCreator.String(),
			firstNFT.ID,
			firstNFT.Reserve,
			firstNFT.AllowMint,
		)

		// only nft creator can mint more nfts for nft tokenID
		err = mintNFT(firstDenom, 10, secondNFT, dsc, ctx)
		require.Error(t, err)

		msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
		require.False(t, broken, msg)
	})

	t.Run("invalid token id", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		sender := app.GetAddrs(dsc, ctx, 1)[0]

		firstNFT := types.NewBaseNFT(
			firstID,
			sender.String(),
			firstTokenURI,
			firstReserve,
			true,
		)

		err := mintNFT(firstDenom, 10, firstNFT, dsc, ctx)
		require.NoError(t, err)

		secondNFT := types.NewBaseNFT(
			firstID,
			sender.String(),
			secondTokenURI,
			secondReserve,
			true,
		)

		// nft token id must be unique
		err = mintNFT(secondDenom, 10, secondNFT, dsc, ctx)
		require.Error(t, err)

		msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
		require.False(t, broken, msg)
	})

	t.Run("invalid token uri", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		sender := app.GetAddrs(dsc, ctx, 1)[0]

		firstNFT := types.NewBaseNFT(
			firstID,
			sender.String(),
			firstTokenURI,
			firstReserve,
			true,
		)

		err := mintNFT(firstDenom, 10, firstNFT, dsc, ctx)
		require.NoError(t, err)

		secondNFT := types.NewBaseNFT(
			secondID,
			sender.String(),
			firstTokenURI,
			secondReserve,
			true,
		)

		// nft token uri must be unique
		err = mintNFT(secondDenom, 10, secondNFT, dsc, ctx)
		require.Error(t, err)

		msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
		require.False(t, broken, msg)
	})

	t.Run("invalid nft reserve", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		sender := app.GetAddrs(dsc, ctx, 1)[0]

		nft := types.NewBaseNFT(
			firstID,
			sender.String(),
			firstTokenURI,
			sdk.NewCoin(config.BaseDenom, sdkmath.ZeroInt()),
			true,
		)

		// nft reserve must be valid
		err := mintNFT(firstDenom, 10, nft, dsc, ctx)
		require.Error(t, err)

		msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
		require.False(t, broken, msg)
	})
}

func TestTransferNFT(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	addrs := app.GetAddrs(dsc, ctx, 2)
	fromOwner := addrs[0]
	toOwner := addrs[1]

	nft := types.NewBaseNFT(
		firstID,
		fromOwner.String(),
		firstTokenURI,
		firstReserve,
		true,
	)
	subTokenIDs := nft.GenSubTokenIDs(4)
	nft = nft.AddOwnerSubTokenIDs(nft.Creator, subTokenIDs)

	fromOwnerSubTokens := nft.GetOwners().GetOwner(fromOwner.String()).GetSubTokenIDs()
	require.Equal(t, types.SortedUintArray{1, 2, 3, 4}, fromOwnerSubTokens)

	const collectionDenom = firstDenom

	err := mintNFT(collectionDenom, int64(subTokenIDs.Len()), nft, dsc, ctx)
	require.NoError(t, err)

	msg := types.MsgTransferNFT{
		Sender:      nft.Creator,
		Recipient:   toOwner.String(),
		ID:          nft.ID,
		Denom:       collectionDenom,
		SubTokenIDs: []uint64{1, 2},
	}

	_, err = dsc.NFTKeeper.TransferNFT(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	storedNFT, err := dsc.NFTKeeper.GetNFT(ctx, collectionDenom, nft.ID)
	require.NoError(t, err)

	fromOwnerSubTokens = storedNFT.GetOwners().GetOwner(fromOwner.String()).GetSubTokenIDs()
	require.Equal(t, types.SortedUintArray{3, 4}, fromOwnerSubTokens)

	toOwnerSubTokens := storedNFT.GetOwners().GetOwner(toOwner.String()).GetSubTokenIDs()
	require.Equal(t, types.SortedUintArray{1, 2}, toOwnerSubTokens)

	invariantMsg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, invariantMsg)
}

func TestEditNFTMetadata(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	nft := types.NewBaseNFT(
		firstID,
		sender.String(),
		firstTokenURI,
		firstReserve,
		true,
	)
	subTokenIDs := nft.GenSubTokenIDs(4)
	nft = nft.AddOwnerSubTokenIDs(nft.Creator, subTokenIDs)

	require.Equal(t, nft.TokenURI, nft.GetTokenURI())

	const collectionDenom = firstDenom
	err := mintNFT(collectionDenom, int64(subTokenIDs.Len()), nft, dsc, ctx)
	require.NoError(t, err)

	msg := types.MsgEditNFTMetadata{
		Sender:   sender.String(),
		ID:       nft.ID,
		Denom:    collectionDenom,
		TokenURI: secondTokenURI,
	}

	_, err = dsc.NFTKeeper.EditNFTMetadata(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	storedNFT, err := dsc.NFTKeeper.GetNFT(ctx, collectionDenom, nft.ID)
	require.NoError(t, err)
	require.Equal(t, msg.TokenURI, storedNFT.GetTokenURI())

	invariantMsg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, invariantMsg)
}

func TestBurnNFT(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	nft := types.NewBaseNFT(
		firstID,
		sender.String(),
		firstTokenURI,
		firstReserve,
		true,
	)
	subTokenIDs := nft.GenSubTokenIDs(4)
	nft = nft.AddOwnerSubTokenIDs(nft.Creator, subTokenIDs)

	creatorSubTokens := nft.GetOwners().GetOwner(sender.String()).GetSubTokenIDs()
	require.Equal(t, types.SortedUintArray{1, 2, 3, 4}, creatorSubTokens)

	const collectionDenom = firstDenom
	err := mintNFT(collectionDenom, int64(subTokenIDs.Len()), nft, dsc, ctx)
	require.NoError(t, err)

	msg := types.MsgBurnNFT{
		Sender:      sender.String(),
		ID:          nft.ID,
		Denom:       collectionDenom,
		SubTokenIDs: subTokenIDs,
	}

	_, err = dsc.NFTKeeper.BurnNFT(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	// Check sub tokens after burn
	storedNFT, err := dsc.NFTKeeper.GetNFT(ctx, collectionDenom, nft.ID)
	require.NoError(t, err)

	creatorSubTokens = storedNFT.GetOwners().GetOwner(sender.String()).GetSubTokenIDs()
	require.Len(t, creatorSubTokens, 0)

	storedSubTokens, err := dsc.NFTKeeper.GetSubTokens(ctx, nft.ID)
	require.NoError(t, err)

	require.Len(t, storedSubTokens, 0)

	invariantMsg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, invariantMsg)
}

func TestUpdateNFTReserve(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	nft := types.NewBaseNFT(
		firstID,
		sender.String(),
		firstTokenURI,
		firstReserve,
		true,
	)
	subTokenIDs := nft.GenSubTokenIDs(4)
	nft = nft.AddOwnerSubTokenIDs(nft.Creator, subTokenIDs)

	const collectionDenom = firstDenom
	err := mintNFT(collectionDenom, int64(subTokenIDs.Len()), nft, dsc, ctx)
	require.NoError(t, err)

	// Check sub tokens reserve before update
	storedSubTokens, err := dsc.NFTKeeper.GetSubTokens(ctx, nft.ID)
	require.NoError(t, err)

	require.Len(t, storedSubTokens, subTokenIDs.Len())
	for i, storedSubToken := range storedSubTokens {
		require.Equal(t, types.NewSubToken(
			subTokenIDs[i],
			firstReserve,
		), storedSubToken)
	}

	msg := types.MsgUpdateReserveNFT{
		Sender:      sender.String(),
		ID:          nft.ID,
		Denom:       collectionDenom,
		SubTokenIDs: subTokenIDs,
		NewReserve:  secondReserve,
	}
	_, err = dsc.NFTKeeper.UpdateReserveNFT(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	// Check sub tokens reserve after update
	storedSubTokens, err = dsc.NFTKeeper.GetSubTokens(ctx, nft.ID)
	require.NoError(t, err)

	require.Len(t, storedSubTokens, subTokenIDs.Len())
	for i, storedSubToken := range storedSubTokens {
		require.Equal(t, types.NewSubToken(
			subTokenIDs[i],
			secondReserve,
		), storedSubToken)
	}

	invariantMsg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, invariantMsg)
}

func mintNFT(denom string, quantity int64, nft types.BaseNFT, dsc *app.DSC, ctx sdk.Context) error {
	msg := types.MsgMintNFT{
		Sender:    nft.Creator,
		Recipient: nft.Creator,
		ID:        nft.ID,
		Denom:     denom,
		Quantity:  sdk.NewInt(quantity),
		TokenURI:  nft.TokenURI,
		Reserve:   nft.Reserve,
		AllowMint: nft.AllowMint,
	}
	_, err := dsc.NFTKeeper.MintNFT(sdk.WrapSDKContext(ctx), &msg)
	return err
}
