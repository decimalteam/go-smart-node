package keeper

//import (
//	"testing"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//
//	"github.com/stretchr/testify/require"
//
//	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
//)
//
//func TestGetOwners(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	quantity := sdk.NewInt(1)
//	reserve := sdk.NewInt(101)
//
//	nft := types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		reserve,
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), quantity, nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//
//	require.NoError(t, err)
//
//	require.NoError(t, err)
//	nft2 := types.NewBaseNFT(
//		ID2,
//		Addrs[1],
//		Addrs[1],
//		TokenURI1,
//		reserve,
//		[]int64{},
//		true,
//	)
//	_, err = NFTKeeper.MintNFT(ctx, Denom1, nft2.GetID(), nft2.GetReserve(), quantity, nft2.GetCreator(), Addrs[1], nft2.GetTokenURI(), nft2.GetAllowMint())
//
//	nft3 := types.NewBaseNFT(
//		ID3,
//		Addrs[2],
//		Addrs[2],
//		TokenURI1,
//		reserve,
//		[]int64{},
//		true,
//	)
//	_, err = NFTKeeper.MintNFT(ctx, Denom1, nft3.GetID(), nft3.GetReserve(), quantity, nft3.GetCreator(), Addrs[2], nft3.GetTokenURI(), nft3.GetAllowMint())
//
//	require.NoError(t, err)
//
//	owners := NFTKeeper.GetOwners(ctx)
//	require.Equal(t, 3, len(owners))
//
//	nft = types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		reserve,
//		[]int64{},
//		true,
//	)
//	_, err = NFTKeeper.MintNFT(ctx, Denom2, nft.GetID(), nft.GetReserve(), quantity, nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//
//	require.NoError(t, err)
//
//	nft2 = types.NewBaseNFT(
//		ID2,
//		Addrs[1],
//		Addrs[1],
//		TokenURI1,
//		reserve,
//		[]int64{},
//		true,
//	)
//	_, err = NFTKeeper.MintNFT(ctx, Denom2, nft2.GetID(), nft2.GetReserve(), sdk.NewInt(1), nft2.GetCreator(), Addrs[1], nft2.GetTokenURI(), nft2.GetAllowMint())
//	require.NoError(t, err)
//
//	nft3 = types.NewBaseNFT(
//		ID3,
//		Addrs[2],
//		Addrs[2],
//		TokenURI1,
//		reserve,
//		[]int64{},
//		true,
//	)
//	_, err = NFTKeeper.MintNFT(ctx, Denom2, nft3.GetID(), nft3.GetReserve(), quantity, nft3.GetCreator(), Addrs[3], nft3.GetTokenURI(), nft3.GetAllowMint())
//
//	require.NoError(t, err)
//
//	owners = NFTKeeper.GetOwners(ctx)
//	require.Equal(t, 3, len(owners))
//
//	msg, fail := SupplyInvariant(NFTKeeper)(ctx)
//	require.False(t, fail, msg)
//}
//
//func TestSetOwner(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	nft := types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(1),
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//
//	require.NoError(t, err)
//
//	idCollection := types.NewIDCollection(Denom1, []string{ID1, ID2, ID3})
//	owner := types.NewOwner(Addrs[0], idCollection)
//
//	oldOwner := NFTKeeper.GetOwner(ctx, Addrs[0])
//
//	NFTKeeper.SetOwner(ctx, owner)
//
//	newOwner := NFTKeeper.GetOwner(ctx, Addrs[0])
//	require.NotEqual(t, oldOwner.String(), newOwner.String())
//	require.Equal(t, owner.String(), newOwner.String())
//
//	NFTKeeper.SetOwner(ctx, oldOwner)
//
//	msg, fail := SupplyInvariant(NFTKeeper)(ctx)
//	require.False(t, fail, msg)
//}
//
//func TestSetOwners(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	// create NFT where ID1 = "ID1" with owner = "Addrs[0]"
//	nft := types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(1),
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//
//	require.NoError(t, err)
//
//	// create NFT where ID1 = "ID2" with owner = "Addrs[1]"
//	nft = types.NewBaseNFT(
//		ID2,
//		Addrs[1],
//		Addrs[1],
//		TokenURI1,
//		sdk.NewInt(1),
//		[]int64{},
//		true,
//	)
//	_, err = NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[1], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	// create two owners (Addrs[0] and Addrs[1]) with the same ID1 collections of "ID1", "ID2"  "ID3"
//	idCollection := types.NewIDCollection(Denom1, []string{ID1, ID2, ID3})
//	owner := types.NewOwner(Addrs[0], idCollection)
//	owner2 := types.NewOwner(Addrs[1], idCollection)
//
//	// get both owners that were created during the NFT mint process
//	oldOwner := NFTKeeper.GetOwner(ctx, Addrs[0])
//	oldOwner2 := NFTKeeper.GetOwner(ctx, Addrs[1])
//
//	// replace previous old owners with updated versions (that have multiple ids)
//	NFTKeeper.SetOwners(ctx, []types.Owner{owner, owner2})
//
//	newOwner := NFTKeeper.GetOwner(ctx, Addrs[0])
//	require.NotEqual(t, oldOwner.String(), newOwner.String())
//	require.Equal(t, owner.String(), newOwner.String())
//
//	newOwner2 := NFTKeeper.GetOwner(ctx, Addrs[1])
//	require.NotEqual(t, oldOwner2.String(), newOwner2.String())
//	require.Equal(t, owner2.String(), newOwner2.String())
//
//	// replace old owners for invariance sanity
//	NFTKeeper.SetOwners(ctx, []types.Owner{oldOwner, oldOwner2})
//
//	msg, fail := SupplyInvariant(NFTKeeper)(ctx)
//	require.False(t, fail, msg)
//}
//
//func TestSwapOwners(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	nft := types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(1),
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//
//	require.NoError(t, err)
//
//	err = NFTKeeper.SwapOwners(ctx, Denom1, ID1, Addrs[0], Addrs[1])
//	require.NoError(t, err)
//
//	err = NFTKeeper.SwapOwners(ctx, Denom1, ID1, Addrs[0], Addrs[1])
//	require.Error(t, err)
//
//	err = NFTKeeper.SwapOwners(ctx, Denom2, ID1, Addrs[0], Addrs[1])
//	require.Error(t, err)
//
//	msg, fail := SupplyInvariant(NFTKeeper)(ctx)
//	require.False(t, fail, msg)
//}
