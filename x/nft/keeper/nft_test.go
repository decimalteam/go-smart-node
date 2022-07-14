package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
	"testing"
)

func getBaseAppWithCustomKeeper() (*codec.LegacyAmino, *app.DSC, sdk.Context) {
	dsc := app.Setup(false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	dsc.NftKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.BankKeeper,
		config.BaseDenom,
	)

	return codec.NewLegacyAmino(), dsc, ctx
}

func TestMintNFT(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := app.AddTestAddrsIncremental(dsc, ctx, 1, sdk.Coins{
		{
			Denom:  "del",
			Amount: helpers.EtherToWei(sdk.NewInt(1000000000000)),
		},
	})

	// MintNFT shouldn't fail when collection does not exist
	msg := types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom1,
		Quantity:  sdk.NewInt(1),
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err := dsc.NftKeeper.MintNFT(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	// MintNFT shouldn't fail when collection exists
	msg = types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID2,
		Denom:     Denom1,
		Quantity:  sdk.NewInt(1),
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err = dsc.NftKeeper.MintNFT(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)

	quantity := sdk.NewInt(50)
	msg = types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID3,
		Denom:     Denom1,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}
	_, err = dsc.NftKeeper.MintNFT(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)
}

//func TestGetNFT(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
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
//
//	require.NoError(t, err)
//
//	// GetNFT should get the NFT
//	receivedNFT, err := NFTKeeper.GetNFT(ctx, Denom1, ID1)
//	require.NoError(t, err)
//	require.Equal(t, receivedNFT.GetID(), ID1)
//	require.True(t, receivedNFT.GetCreator().Equals(Addrs[0]))
//	require.Equal(t, receivedNFT.GetTokenURI(), TokenURI1)
//
//	// MintNFT shouldn't fail when collection exists
//	nft2 := types.NewBaseNFT(
//		ID2,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		[]int64{},
//		true,
//	)
//	_, err = NFTKeeper.MintNFT(ctx, Denom1, nft2.GetID(), nft2.GetReserve(), sdk.NewInt(1), nft2.GetCreator(), Addrs[0], nft2.GetTokenURI(), nft2.GetAllowMint())
//	require.NoError(t, err)
//
//	// GetNFT should get the NFT when collection exists
//	receivedNFT2, err := NFTKeeper.GetNFT(ctx, Denom1, ID2)
//	require.NoError(t, err)
//	require.Equal(t, receivedNFT2.GetID(), ID2)
//	require.True(t, receivedNFT2.GetCreator().Equals(Addrs[0]))
//	require.Equal(t, receivedNFT2.GetTokenURI(), TokenURI1)
//
//	msg, fail := SupplyInvariant(NFTKeeper)(ctx)
//	require.False(t, fail, msg)
//}
//
//func TestUpdateNFT(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	subTokenIDs := []int64{}
//	reserve := sdk.NewInt(100)
//	quantity := sdk.NewInt(1)
//
//	nft := types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		reserve,
//		subTokenIDs,
//		true,
//	)
//
//	// UpdateNFT should fail when nft doesn't exist
//	//_, err := NFTKeeper.MintNFT(ctx, Denom1, ID2, reserve, quantity, Addrs[0], Addrs[0], TokenURI1, true)
//	//require.Error(t, err)
//
//	// MintNFT shouldn't fail when collection does not exist
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, ID1, nft.GetReserve(), quantity, nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	// UpdateNFT shouldn't fail when NFT exists
//	nft2 := types.NewBaseNFT(
//		ID1,
//		Addrs[1],
//		Addrs[1],
//		TokenURI2,
//		reserve,
//		subTokenIDs,
//		true,
//	)
//	err = NFTKeeper.UpdateNFT(ctx, Denom1, nft2)
//	require.NoError(t, err)
//
//	// GetNFT should get the NFT with new TokenURI1
//	receivedNFT, err := NFTKeeper.GetNFT(ctx, Denom1, ID1)
//	require.NoError(t, err)
//	require.Equal(t, receivedNFT.GetTokenURI(), TokenURI2)
//}
//
//func TestDeleteNFT(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	ctx = ctx.WithBlockHeight(updates.Update11Block)
//
//	// DeleteNFT should fail when NFT doesn't exist and collection doesn't exist
//	subTokenIDs := []int64{}
//	err := NFTKeeper.DeleteNFT(ctx, Denom1, ID1, subTokenIDs)
//	require.Error(t, err)
//
//	// MintNFT should not fail when collection does not exist
//	nft := types.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		subTokenIDs,
//		true,
//	)
//	_, err = NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//
//	require.NoError(t, err)
//
//	// DeleteNFT should fail when NFT doesn't exist but collection does exist
//	err = NFTKeeper.DeleteNFT(ctx, Denom1, ID2, subTokenIDs)
//	require.Error(t, err)
//
//	// DeleteNFT should fail when at least of nft's subtokenIds is not in the owner's subTokenIDs
//	err = NFTKeeper.DeleteNFT(ctx, Denom1, ID1, []int64{3})
//	require.Error(t, err)
//
//	// DeleteNFT should not fail when NFT and collection exist
//	err = NFTKeeper.DeleteNFT(ctx, Denom1, ID1, subTokenIDs)
//	require.NoError(t, err)
//
//	// NFT should no longer exist
//	isNFT := NFTKeeper.IsNFT(ctx, Denom1, ID1)
//	require.True(t, isNFT)
//
//	owner := NFTKeeper.GetOwner(ctx, Addrs[0])
//	require.Equal(t, 1, owner.Supply())
//}
//
//func TestIsNFT(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	// IsNFT should return false
//	isNFT := NFTKeeper.IsNFT(ctx, Denom1, ID1)
//	require.False(t, isNFT)
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
//
//	require.NoError(t, err)
//
//	// IsNFT should return true
//	isNFT = NFTKeeper.IsNFT(ctx, Denom1, ID1)
//	require.True(t, isNFT)
//}
