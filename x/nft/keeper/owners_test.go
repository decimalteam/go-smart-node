package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestGetOwners(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 3)
	quantity := sdk.NewInt(1)

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

	msg = types.MsgMintNFT{
		Sender:    addrs[1].String(),
		Recipient: addrs[1].String(),
		ID:        ID2,
		Denom:     Denom1,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	msg = types.MsgMintNFT{
		Sender:    addrs[2].String(),
		Recipient: addrs[2].String(),
		ID:        ID3,
		Denom:     Denom1,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	owners := dsc.NftKeeper.GetOwners(ctx)
	require.Equal(t, 3, len(owners))

	msg = types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom2,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	msg = types.MsgMintNFT{
		Sender:    addrs[1].String(),
		Recipient: addrs[1].String(),
		ID:        ID2,
		Denom:     Denom2,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	msg = types.MsgMintNFT{
		Sender:    addrs[2].String(),
		Recipient: addrs[2].String(),
		ID:        ID3,
		Denom:     Denom2,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	owners = dsc.NftKeeper.GetOwners(ctx)
	require.Equal(t, 3, len(owners))

	invariantMsg, fail := keeper.SupplyInvariant(dsc.NftKeeper)(ctx)
	require.False(t, fail, invariantMsg)
}

func TestSetOwner(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)
	quantity := sdk.NewInt(1)

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

	idCollection := types.NewIDCollection(Denom1, []string{ID1, ID2, ID3})
	owner := types.NewOwner(addrs[0].String(), idCollection)

	oldOwner := dsc.NftKeeper.GetOwner(ctx, addrs[0])

	dsc.NftKeeper.SetOwner(ctx, owner)

	newOwner := dsc.NftKeeper.GetOwner(ctx, addrs[0])
	require.NotEqual(t, oldOwner.String(), newOwner.String())
	require.Equal(t, owner.String(), newOwner.String())

	dsc.NftKeeper.SetOwner(ctx, oldOwner)

	invariantMsg, fail := keeper.SupplyInvariant(dsc.NftKeeper)(ctx)
	require.False(t, fail, invariantMsg)
}

func TestSetOwners(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 2)
	quantity := sdk.NewInt(1)

	// create NFT where ID1 = "ID1" with owner = "Addrs[0]"
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

	// create NFT where ID1 = "ID2" with owner = "Addrs[1]"
	msg = types.MsgMintNFT{
		Sender:    addrs[1].String(),
		Recipient: addrs[1].String(),
		ID:        ID2,
		Denom:     Denom1,
		Quantity:  quantity,
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}
	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	// create two owners (Addrs[0] and Addrs[1]) with the same ID1 collections of "ID1", "ID2"  "ID3"
	idCollection := types.NewIDCollection(Denom1, []string{ID1, ID2, ID3})
	owner := types.NewOwner(addrs[0].String(), idCollection)
	owner2 := types.NewOwner(addrs[1].String(), idCollection)

	// get both owners that were created during the NFT mint process
	oldOwner := dsc.NftKeeper.GetOwner(ctx, addrs[0])
	oldOwner2 := dsc.NftKeeper.GetOwner(ctx, addrs[1])

	// replace previous old owners with updated versions (that have multiple ids)
	dsc.NftKeeper.SetOwners(ctx, []types.Owner{owner, owner2})

	newOwner := dsc.NftKeeper.GetOwner(ctx, addrs[0])
	require.NotEqual(t, oldOwner.String(), newOwner.String())
	require.Equal(t, owner.String(), newOwner.String())

	newOwner2 := dsc.NftKeeper.GetOwner(ctx, addrs[1])
	require.NotEqual(t, oldOwner2.String(), newOwner2.String())
	require.Equal(t, owner2.String(), newOwner2.String())

	// replace old owners for invariance sanity
	dsc.NftKeeper.SetOwners(ctx, []types.Owner{oldOwner, oldOwner2})

	invariantMsg, fail := keeper.SupplyInvariant(dsc.NftKeeper)(ctx)
	require.False(t, fail, invariantMsg)
}
