package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetCollection(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)
	quantity := sdk.NewInt(50)

	// create a new nft with id = "id" and owner = "address"
	// MintNFT shouldn't fail when collection does not exist
	msg := types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom1,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	// collection should exist
	collection, exists := dsc.NftKeeper.GetCollection(ctx, Denom1)
	require.True(t, exists)
	require.NotEmpty(t, collection)

	// create a new NFT and add it to the collection created with the NFT mint
	subTokenIDs := []int64{2, 5, 3, 6}
	nft2 := types.NewBaseNFT(
		ID2,
		addrs[0].String(),
		addrs[0].String(),
		TokenURI1,
		sdk.NewInt(100),
		subTokenIDs,
		true,
	)
	collection2, err2 := collection.AddNFT(nft2)
	require.NoError(t, err2)
	dsc.NftKeeper.SetCollection(ctx, Denom1, collection2)

	collection2, exists = dsc.NftKeeper.GetCollection(ctx, Denom1)
	require.True(t, exists)
	require.NotEmpty(t, collection)
	require.Len(t, collection2.NFTs, 2)

	collectionNFT, err := collection2.GetNFT(nft2.GetID())
	require.NoError(t, err)

	orderedSubTokens := []int64{2, 3, 5, 6}
	require.Equal(t, orderedSubTokens, collectionNFT.GetOwners().GetOwner(addrs[0].String()).GetSubTokenIDs())

	// reset collection for invariant sanity
	dsc.NftKeeper.SetCollection(ctx, Denom1, collection)

	invariantMsg, fail := keeper.SupplyInvariant(dsc.NftKeeper)(ctx)
	require.False(t, fail, invariantMsg)
}

func TestGetCollection(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)
	quantity := sdk.NewInt(50)

	// collection shouldn't exist
	collection, exists := dsc.NftKeeper.GetCollection(ctx, Denom1)
	require.Empty(t, collection)
	require.False(t, exists)

	// MintNFT shouldn't fail when collection does not exist
	msg := types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom1,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	// collection should exist
	collection, exists = dsc.NftKeeper.GetCollection(ctx, Denom1)
	require.True(t, exists)
	require.NotEmpty(t, collection)

	invariantMsg, fail := keeper.SupplyInvariant(dsc.NftKeeper)(ctx)
	require.False(t, fail, invariantMsg)
}

func TestGetCollections(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)
	quantity := sdk.NewInt(50)

	// collections should be empty
	collections := dsc.NftKeeper.GetCollections(ctx)
	require.Empty(t, collections)

	// MintNFT shouldn't fail when collection does not exist
	msg := types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom1,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	// collections should equal 1
	collections = dsc.NftKeeper.GetCollections(ctx)
	require.NotEmpty(t, collections)
	require.Equal(t, len(collections), 1)

	invariantMsg, fail := keeper.SupplyInvariant(dsc.NftKeeper)(ctx)
	require.False(t, fail, invariantMsg)
}
