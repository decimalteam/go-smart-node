package keeper_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMintNFT(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	expectedCollection := types.Collection{
		Creator: sender.String(),
		Denom:   firstDenom,
		Supply:  1,
		Tokens: []*types.Token{
			{
				Creator:   sender.String(),
				Denom:     firstDenom,
				ID:        firstID,
				URI:       firstTokenURI,
				Reserve:   firstReserve,
				AllowMint: true,
				Minted:    2,
				Burnt:     0,
				SubTokens: []*types.SubToken{
					{ID: 1, Owner: sender.String(), Reserve: &firstReserve},
					{ID: 2, Owner: sender.String(), Reserve: &firstReserve},
				},
			},
		},
	}

	err := mintNFT(firstDenom, 2, *expectedCollection.Tokens[0], dsc, ctx)
	require.NoError(t, err)

	storedCollection, found := dsc.NFTKeeper.GetCollection(ctx, sender, firstDenom)
	require.True(t, found)
	require.Equal(t, expectedCollection, storedCollection)

	storedNFT, found := dsc.NFTKeeper.GetToken(ctx, firstID)
	require.True(t, found)
	require.Equal(t, *expectedCollection.Tokens[0], storedNFT)

	storedSubTokens := dsc.NFTKeeper.GetSubTokens(ctx, firstID)

	require.Len(t, storedSubTokens, len(expectedCollection.Tokens[0].SubTokens))
	require.Equal(t, expectedCollection.Tokens[0].SubTokens, storedSubTokens)

	msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, msg)
}

func TestMintNFTValidation(t *testing.T) {
	t.Run("disabled mint", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		sender := app.GetAddrs(dsc, ctx, 1)[0]

		nft := types.Token{
			Creator:   sender.String(),
			Denom:     firstDenom,
			ID:        firstID,
			URI:       firstTokenURI,
			Reserve:   firstReserve,
			AllowMint: false,
		}

		// first nft mint with allowMint=false flag should be without error
		err := mintNFT(firstDenom, 2, nft, dsc, ctx)
		require.NoError(t, err)

		// next nft mint with allowMint=false flag should be with error
		err = mintNFT(firstDenom, 1, nft, dsc, ctx)
		require.Error(t, err)

		msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
		require.False(t, broken, msg)
	})

	t.Run("wrong sender account", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		addrs := app.GetAddrs(dsc, ctx, 2)
		creator := addrs[0]
		wrongCreator := addrs[1]

		firstNFT := types.Token{
			Creator:   creator.String(),
			Denom:     firstDenom,
			ID:        firstID,
			URI:       firstTokenURI,
			Reserve:   firstReserve,
			AllowMint: false,
		}

		err := mintNFT(firstDenom, 1, firstNFT, dsc, ctx)
		require.NoError(t, err)

		secondNFT := types.Token{
			Creator:   wrongCreator.String(),
			Denom:     firstDenom,
			ID:        firstID,
			URI:       firstTokenURI,
			Reserve:   firstReserve,
			AllowMint: false,
		}

		// only nft creator can mint more nfts for nft tokenID
		err = mintNFT(firstDenom, 1, secondNFT, dsc, ctx)
		require.Error(t, err)

		msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
		require.False(t, broken, msg)
	})

	t.Run("invalid token id", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		sender := app.GetAddrs(dsc, ctx, 1)[0]

		firstNFT := types.Token{
			Creator:   sender.String(),
			Denom:     firstDenom,
			ID:        firstID,
			URI:       firstTokenURI,
			Reserve:   firstReserve,
			AllowMint: true,
		}

		err := mintNFT(firstDenom, 2, firstNFT, dsc, ctx)
		require.NoError(t, err)

		secondNFT := types.Token{
			Creator:   sender.String(),
			Denom:     firstDenom,
			ID:        firstID,
			URI:       secondTokenURI,
			Reserve:   secondReserve,
			AllowMint: true,
		}

		// nft token id must be unique
		err = mintNFT(secondDenom, 10, secondNFT, dsc, ctx)
		require.Error(t, err)

		msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
		require.False(t, broken, msg)
	})

	t.Run("invalid token uri", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		sender := app.GetAddrs(dsc, ctx, 1)[0]

		firstNFT := types.Token{
			Creator:   sender.String(),
			Denom:     firstDenom,
			ID:        firstID,
			URI:       firstTokenURI,
			Reserve:   firstReserve,
			AllowMint: true,
		}

		err := mintNFT(firstDenom, 10, firstNFT, dsc, ctx)
		require.NoError(t, err)

		secondNFT := types.Token{
			Creator:   sender.String(),
			Denom:     secondDenom,
			ID:        secondID,
			URI:       firstTokenURI,
			Reserve:   secondReserve,
			AllowMint: true,
		}

		// nft token uri must be unique
		err = mintNFT(secondDenom, 10, secondNFT, dsc, ctx)
		require.Error(t, err)

		msg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
		require.False(t, broken, msg)
	})

	t.Run("invalid nft reserve", func(t *testing.T) {
		dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

		sender := app.GetAddrs(dsc, ctx, 1)[0]

		nft := types.Token{
			Creator:   sender.String(),
			Denom:     firstDenom,
			ID:        firstID,
			URI:       firstTokenURI,
			Reserve:   sdk.NewCoin(config.BaseDenom, sdkmath.ZeroInt()),
			AllowMint: true,
		}

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

	nft := types.Token{
		Creator:   fromOwner.String(),
		Denom:     firstDenom,
		ID:        firstID,
		URI:       firstTokenURI,
		Reserve:   firstReserve,
		AllowMint: true,
	}

	err := mintNFT(firstDenom, 4, nft, dsc, ctx)
	require.NoError(t, err)

	msg := types.MsgSendToken{
		Sender:      nft.Creator,
		Recipient:   toOwner.String(),
		TokenID:     nft.ID,
		SubTokenIDs: []uint32{1, 2},
	}

	_, err = dsc.NFTKeeper.SendToken(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	_, found := dsc.NFTKeeper.GetToken(ctx, nft.ID)
	require.True(t, found)

	subTokens := dsc.NFTKeeper.GetSubTokens(ctx, nft.ID)
	require.Len(t, subTokens, 4)

	for _, sub := range subTokens {
		switch sub.ID {
		case 1:
			require.Equal(t, toOwner.String(), sub.Owner)
		case 2:
			require.Equal(t, toOwner.String(), sub.Owner)
		case 3:
			require.Equal(t, fromOwner.String(), sub.Owner)
		case 4:
			require.Equal(t, fromOwner.String(), sub.Owner)
		}
	}

	invariantMsg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, invariantMsg)
}

func TestEditNFTMetadata(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	nft := types.Token{
		Creator:   sender.String(),
		Denom:     firstDenom,
		ID:        firstID,
		URI:       firstTokenURI,
		Reserve:   firstReserve,
		AllowMint: true,
	}

	err := mintNFT(firstDenom, 2, nft, dsc, ctx)
	require.NoError(t, err)

	msg := types.MsgUpdateToken{
		Sender:   sender.String(),
		TokenID:  nft.ID,
		TokenURI: secondTokenURI,
	}

	_, err = dsc.NFTKeeper.UpdateToken(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	storedNFT, found := dsc.NFTKeeper.GetToken(ctx, nft.ID)
	require.True(t, found)
	require.Equal(t, msg.TokenURI, storedNFT.URI)

	invariantMsg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, invariantMsg)
}

func TestBurnNFT(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	nft := types.Token{
		Creator:   sender.String(),
		Denom:     firstDenom,
		ID:        firstID,
		URI:       firstTokenURI,
		Reserve:   firstReserve,
		AllowMint: true,
	}

	const collectionDenom = firstDenom
	err := mintNFT(firstDenom, 4, nft, dsc, ctx)
	require.NoError(t, err)

	msg := types.MsgBurnToken{
		Sender:      sender.String(),
		TokenID:     nft.ID,
		SubTokenIDs: []uint32{2, 3},
	}

	_, err = dsc.NFTKeeper.BurnToken(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	// Check sub tokens after burn
	subTokens := dsc.NFTKeeper.GetSubTokens(ctx, nft.ID)
	require.Len(t, subTokens, 2)

	for _, sub := range subTokens {
		switch sub.ID {
		case 1:
			require.Equal(t, sender.String(), sub.Owner)
		case 4:
			require.Equal(t, sender.String(), sub.Owner)
		default:
			t.Fail()
		}
	}

	invariantMsg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, invariantMsg)
}

func TestUpdateNFTReserve(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	sender := app.GetAddrs(dsc, ctx, 1)[0]

	nft := types.Token{
		Creator:   sender.String(),
		Denom:     firstDenom,
		ID:        firstID,
		URI:       firstTokenURI,
		Reserve:   firstReserve,
		AllowMint: true,
	}

	err := mintNFT(firstDenom, 4, nft, dsc, ctx)
	require.NoError(t, err)

	// Check sub tokens reserve before update
	storedSubTokens := dsc.NFTKeeper.GetSubTokens(ctx, nft.ID)

	require.Len(t, storedSubTokens, 4)
	for _, storedSubToken := range storedSubTokens {
		require.True(t, storedSubToken.Reserve.Equal(firstReserve))
	}

	msg := types.MsgUpdateReserve{
		Sender:      sender.String(),
		TokenID:     nft.ID,
		SubTokenIDs: []uint32{1, 3},
		Reserve:     secondReserve,
	}
	_, err = dsc.NFTKeeper.UpdateReserve(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	// Check sub tokens reserve after update
	storedSubTokens = dsc.NFTKeeper.GetSubTokens(ctx, nft.ID)
	require.Len(t, storedSubTokens, 4)

	for _, storedSubToken := range storedSubTokens {
		switch storedSubToken.ID {
		case 1:
			require.True(t, storedSubToken.Reserve.Equal(secondReserve))
		case 2:
			require.True(t, storedSubToken.Reserve.Equal(firstReserve))
		case 3:
			require.True(t, storedSubToken.Reserve.Equal(secondReserve))
		case 4:
			require.True(t, storedSubToken.Reserve.Equal(firstReserve))
		}
	}

	invariantMsg, broken := keeper.AllInvariants(dsc.NFTKeeper)(ctx)
	require.False(t, broken, invariantMsg)
}

func mintNFT(denom string, quantity uint32, nft types.Token, dsc *app.DSC, ctx sdk.Context) error {
	msg := types.MsgMintToken{
		Sender:    nft.Creator,
		Denom:     denom,
		TokenID:   nft.ID,
		TokenURI:  nft.URI,
		AllowMint: nft.AllowMint,
		Recipient: nft.Creator,
		Quantity:  quantity,
		Reserve:   nft.Reserve,
	}
	_, err := dsc.NFTKeeper.MintToken(sdk.WrapSDKContext(ctx), &msg)
	return err
}
