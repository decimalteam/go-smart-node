package legacy_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	multisig "bitbucket.org/decimalteam/go-smart-node/x/multisig"
	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nft "bitbucket.org/decimalteam/go-smart-node/x/nft"
	nftTypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/strings"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestInitGenesisForLegacy(t *testing.T) {
	app, ctx := getBaseApp(t)

	publicKey := []byte{0x3, 0x44, 0x8e, 0x6b, 0x3d, 0x50, 0xd6, 0xa3, 0x9c, 0xab, 0x3b, 0xab, 0xaa,
		0x4a, 0xa2, 0xb0, 0x88, 0x5f, 0x55, 0x6f, 0xe0, 0x5d, 0x71, 0x49, 0x88, 0x5a, 0x5, 0xa0, 0xe7, 0x94, 0xa, 0x7e, 0x4f}
	oldAddress, err := commonTypes.GetLegacyAddressFromPubKey(publicKey)
	require.NoError(t, err)
	newAddress := "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn"

	otherAddress := "dx1m3eg7v6pu0dga2knj9zm4683dk9c8800j9nfw0"

	wallet1, wallet2 := "dx108c4p0j7wqsawejfuuv43hj7nhyp36gt0296rs", "dx10fx59x9ytvf249axryvw0uh3eunwvgyfpm9jrp"

	legacyCoinPoolAddress, err := sdk.Bech32ifyAddressBytes(config.Bech32Prefix,
		cosmosAuthTypes.NewModuleAddress(types.LegacyCoinPool))
	require.NoError(t, err, "legacyCoinPoolAddress to bech32")

	bankGenesisState := &cosmosBankTypes.GenesisState{
		Params: cosmosBankTypes.DefaultParams(),
		Balances: []cosmosBankTypes.Balance{
			{
				Address: legacyCoinPoolAddress,
				Coins: sdk.Coins{
					{
						Denom:  "del",
						Amount: sdk.NewInt(100),
					},
					{
						Denom:  "foo",
						Amount: sdk.NewInt(100),
					},
				},
			},
		},
	}

	require.NoError(t, bankGenesisState.Validate(), "bankGenesisState")
	app.BankKeeper.InitGenesis(ctx, bankGenesisState)

	nftGenesisState := nftTypes.GenesisState{
		Collections: []nftTypes.Collection{
			{
				Denom: "a",
				NFTs:  nftTypes.SortedStringArray{"a1", "a2"},
			},
			{
				Denom: "b",
				NFTs:  nftTypes.SortedStringArray{"b1", "b2"},
			},
		},
		NFTs: []nftTypes.BaseNFT{
			{
				ID: "a1",
				Owners: nftTypes.TokenOwners{
					{Address: oldAddress, SubTokenIDs: []uint64{1}},
				},
				Creator:   newAddress,
				TokenURI:  "a1",
				Reserve:   sdk.NewCoin("del", sdk.NewInt(100)),
				AllowMint: false,
			},
			{
				ID: "a2",
				Owners: nftTypes.TokenOwners{
					{Address: oldAddress, SubTokenIDs: []uint64{1}},
				},
				Creator:   newAddress,
				TokenURI:  "a2",
				Reserve:   sdk.NewCoin("del", sdk.NewInt(100)),
				AllowMint: false,
			},
			{
				ID: "b1",
				Owners: nftTypes.TokenOwners{
					{Address: oldAddress, SubTokenIDs: []uint64{1, 2}},
				},
				Creator:   newAddress,
				TokenURI:  "b1",
				Reserve:   sdk.NewCoin("del", sdk.NewInt(100)),
				AllowMint: false,
			},
			{
				ID: "b2",
				Owners: nftTypes.TokenOwners{
					{Address: oldAddress, SubTokenIDs: []uint64{1, 2}},
				},
				Creator:   newAddress,
				TokenURI:  "b2",
				Reserve:   sdk.NewCoin("del", sdk.NewInt(100)),
				AllowMint: false,
			},
		},
		SubTokens: map[string]nftTypes.SubTokens{
			"a1": {
				SubTokens: []nftTypes.SubToken{
					{ID: 1, Reserve: sdk.NewCoin("del", sdk.NewInt(100))},
				},
			},
			"a2": {
				SubTokens: []nftTypes.SubToken{
					{ID: 1, Reserve: sdk.NewCoin("del", sdk.NewInt(100))},
				},
			},
			"b1": {
				SubTokens: []nftTypes.SubToken{
					{ID: 1, Reserve: sdk.NewCoin("del", sdk.NewInt(100))},
					{ID: 2, Reserve: sdk.NewCoin("del", sdk.NewInt(100))},
				},
			},
			"b2": {
				SubTokens: []nftTypes.SubToken{
					{ID: 1, Reserve: sdk.NewCoin("del", sdk.NewInt(100))},
					{ID: 2, Reserve: sdk.NewCoin("del", sdk.NewInt(100))},
				},
			},
		},
	}
	require.NoError(t, nftGenesisState.Validate(), "nftGenesisState")
	nft.InitGenesis(ctx, app.NFTKeeper, nftGenesisState)

	multisigGenesisState := multisigTypes.GenesisState{
		Wallets: []multisigTypes.Wallet{
			{
				Address:   wallet1,
				Owners:    []string{otherAddress, oldAddress},
				Weights:   []uint64{1, 2},
				Threshold: 3,
			},
			{
				Address:   wallet2,
				Owners:    []string{otherAddress, oldAddress},
				Weights:   []uint64{1, 2},
				Threshold: 3,
			},
		},
	}
	require.NoError(t, nftGenesisState.Validate(), "nftGenesisState")
	multisig.InitGenesis(ctx, app.MultisigKeeper, multisigGenesisState)

	legacyGenesisState := types.GenesisState{
		LegacyRecords: []types.LegacyRecord{
			{
				Address: oldAddress,
				Coins:   sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(50)), sdk.NewCoin("foo", sdk.NewInt(1))),
				Nfts:    []types.NFTRecord{{Denom: "a", Id: "a2"}, {Denom: "b", Id: "b1"}},
				Wallets: []string{wallet1},
			},
			{
				Address: otherAddress,
				Wallets: []string{wallet2},
			},
		},
	}

	require.NoError(t, legacyGenesisState.Validate(), "legacyGenesisState")
	legacy.InitGenesis(ctx, app.LegacyKeeper, legacyGenesisState)

	// init  genesis done
	// let's check
	app.LegacyKeeper.RestoreCache(ctx)
	require.True(t, app.LegacyKeeper.IsLegacyAddress(ctx, oldAddress), "check legacy address 1")
	require.False(t, app.LegacyKeeper.IsLegacyAddress(ctx, newAddress), "check legacy address 2")
	err = app.LegacyKeeper.ActualizeLegacy(ctx, publicKey)
	require.NoError(t, err, "ActualizeLegacy")

	// coins
	legacyRemain := app.BankKeeper.GetAllBalances(ctx, cosmosAuthTypes.NewModuleAddress(types.LegacyCoinPool))
	require.True(t, legacyRemain.IsEqual(sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(50)), sdk.NewCoin("foo", sdk.NewInt(99)))),
		"legacy coins remain")

	// nft no changes
	nft, err := app.NFTKeeper.GetNFT(ctx, "a", "a1")
	require.NoError(t, err, "nft-a-a1")
	require.NotNil(t, nft.GetOwners().GetOwner(oldAddress), "nft-a-a1 must not changed")
	nft, err = app.NFTKeeper.GetNFT(ctx, "b", "b2")
	require.NoError(t, err, "nft-b-b2")
	require.NotNil(t, nft.GetOwners().GetOwner(oldAddress), "nft-a-b2 must not changed")
	//  nft changes
	nft, err = app.NFTKeeper.GetNFT(ctx, "a", "a2")
	require.NoError(t, err, "nft-a-a2")
	require.Nil(t, nft.GetOwners().GetOwner(oldAddress), "nft-a-a2 must changed")
	require.NotNil(t, nft.GetOwners().GetOwner(newAddress), "nft-a-a2 must changed")
	nft, err = app.NFTKeeper.GetNFT(ctx, "b", "b1")
	require.NoError(t, err, "nft-b-b1")
	require.Nil(t, nft.GetOwners().GetOwner(oldAddress), "nft-b-b1 must changed")
	require.NotNil(t, nft.GetOwners().GetOwner(newAddress), "nft-b-b1 must changed")

	// wallet changes
	wallet, err := app.MultisigKeeper.GetWallet(ctx, wallet1)
	require.NoError(t, err, "multisig wallet 1")
	require.False(t, strings.StringInSlice(oldAddress, wallet.Owners), "wallet must have no old owner")
	require.True(t, strings.StringInSlice(otherAddress, wallet.Owners), "wallet must keep other owner")
	// wallet without changes
	wallet, err = app.MultisigKeeper.GetWallet(ctx, wallet2)
	require.NoError(t, err, "multisig wallet 2")
	require.True(t, strings.StringInSlice(oldAddress, wallet.Owners), "wallet must not changed")
	require.True(t, strings.StringInSlice(otherAddress, wallet.Owners), "wallet must not changed")

	// check kepeer end state
	_, err = app.LegacyKeeper.GetLegacyRecord(ctx, oldAddress)
	require.Error(t, err, "must no record")
	require.False(t, app.LegacyKeeper.IsLegacyAddress(ctx, oldAddress), "check legacy address at end")
}

func getBaseApp(t *testing.T) (*app.DSC, sdk.Context) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	return dsc, ctx
}
