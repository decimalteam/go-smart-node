package keeper

//import (
//	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/stretchr/testify/require"
//	"testing"
//)
//
//func TestSetCollection(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	// create a new nft with id = "id" and owner = "address"
//	// MintNFT shouldn't fail when collection does not exist
//	nft := types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		[]int64{},
//		false,
//	)
//
//	//_, err := NFTKeeper.MintNFT(ctx, Denom1, nft)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	// collection should exist
//	collection, exists := NFTKeeper.GetCollection(ctx, Denom1)
//	require.True(t, exists)
//
//	// create a new NFT and add it to the collection created with the NFT mint
//	subTokenIDs := []int64{2, 5, 3, 6}
//	nft2 := types.NewBaseNFT(
//		ID2,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		subTokenIDs,
//		true,
//	)
//	collection2, err2 := collection.AddNFT(nft2)
//	require.NoError(t, err2)
//	NFTKeeper.SetCollection(ctx, Denom1, collection2)
//
//	collection2, exists = NFTKeeper.GetCollection(ctx, Denom1)
//	require.True(t, exists)
//	require.Len(t, collection2.NFTs, 2)
//
//	collectionNFT, _ := collection2.GetNFT(nft2.GetID())
//
//	require.Equal(t, len(subTokenIDs), len(collectionNFT.GetOwners().GetOwner(Addrs[0]).GetSubTokenIDs()))
//
//	// reset collection for invariant sanity
//	NFTKeeper.SetCollection(ctx, Denom1, collection)
//
//	msg, fail := SupplyInvariant(NFTKeeper)(ctx)
//	require.False(t, fail, msg)
//}
//func TestGetCollection(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	// collection shouldn't exist
//	collection, exists := NFTKeeper.GetCollection(ctx, Denom1)
//	require.Empty(t, collection)
//	require.False(t, exists)
//
//	// MintNFT shouldn't fail when collection does not exist
//	nft := types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	// collection should exist
//	collection, exists = NFTKeeper.GetCollection(ctx, Denom1)
//	require.True(t, exists)
//	require.NotEmpty(t, collection)
//
//	msg, fail := SupplyInvariant(NFTKeeper)(ctx)
//	require.False(t, fail, msg)
//}
//func TestGetCollections(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	// collections should be empty
//	collections := NFTKeeper.GetCollections(ctx)
//	require.Empty(t, collections)
//
//	// MintNFT shouldn't fail when collection does not exist
//	nft := types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	// collections should equal 1
//	collections = NFTKeeper.GetCollections(ctx)
//	require.NotEmpty(t, collections)
//	require.Equal(t, len(collections), 1)
//
//	msg, fail := SupplyInvariant(NFTKeeper)(ctx)
//	require.False(t, fail, msg)
//}
