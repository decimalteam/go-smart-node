package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
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

	addrs := getAddrs(dsc, ctx)

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

	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
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

	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
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
	lastSubTokenID, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)
	require.Equal(t, quantity.AddRaw(1).Int64(), lastSubTokenID)
}

func TestGetNFT(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx)

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

	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	// GetNFT should get the NFT
	receivedNFT, err := dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, receivedNFT.GetID(), ID1)
	require.True(t, receivedNFT.GetCreator() == addrs[0].String())
	require.Equal(t, receivedNFT.GetTokenURI(), TokenURI1)

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

	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	// GetNFT should get the NFT when collection exists
	receivedNFT2, err := dsc.NftKeeper.GetNFT(ctx, Denom1, ID2)
	require.NoError(t, err)
	require.Equal(t, receivedNFT2.GetID(), ID2)
	require.True(t, receivedNFT2.GetCreator() == addrs[0].String())
	require.Equal(t, receivedNFT2.GetTokenURI(), TokenURI1)

	invariantMsg, fail := keeper.SupplyInvariant(dsc.NftKeeper)(ctx)
	require.False(t, fail, invariantMsg)
}

func TestUpdateNFT(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx)

	var subTokenIDs []int64
	reserve := sdk.NewInt(100)

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

	// MintNFT shouldn't fail when collection does not exist
	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	// UpdateNFT shouldn't fail when NFT exists
	nft := types.NewBaseNFT(
		ID1,
		addrs[1].String(),
		addrs[1].String(),
		TokenURI2,
		reserve,
		subTokenIDs,
		true,
	)
	err = dsc.NftKeeper.UpdateNFT(ctx, Denom1, nft)
	require.NoError(t, err)

	// GetNFT should get the NFT with new TokenURI1
	receivedNFT, err := dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, receivedNFT.GetTokenURI(), TokenURI2)
}

func TestDeleteNFT(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx)

	// DeleteNFT should fail when NFT doesn't exist and collection doesn't exist
	var subTokenIDs []int64
	err := dsc.NftKeeper.DeleteNFT(ctx, Denom1, ID1, subTokenIDs)
	require.Error(t, err)

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

	// MintNFT should not fail when collection does not exist
	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	// DeleteNFT should fail when NFT doesn't exist but collection does exist
	err = dsc.NftKeeper.DeleteNFT(ctx, Denom1, ID2, subTokenIDs)
	require.Error(t, err)

	// DeleteNFT should fail when at least of nft's subtokenIds is not in the owner's subTokenIDs
	err = dsc.NftKeeper.DeleteNFT(ctx, Denom1, ID1, []int64{3})
	require.Error(t, err)

	// DeleteNFT should not fail when NFT and collection exist
	err = dsc.NftKeeper.DeleteNFT(ctx, Denom1, ID1, subTokenIDs)
	require.NoError(t, err)

	// NFT should no longer exist
	isNFT := dsc.NftKeeper.IsNFT(ctx, Denom1, ID1)
	require.True(t, isNFT)

	owner := dsc.NftKeeper.GetOwner(ctx, addrs[0])
	require.Equal(t, 1, owner.Supply())
}

func TestIsNFT(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx)

	// IsNFT should return false
	isNFT := dsc.NftKeeper.IsNFT(ctx, Denom1, ID1)
	require.False(t, isNFT)

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
	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)

	require.NoError(t, err)

	// IsNFT should return true
	isNFT = dsc.NftKeeper.IsNFT(ctx, Denom1, ID1)
	require.True(t, isNFT)
}
