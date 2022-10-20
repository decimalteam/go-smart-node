package legacy_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/libs/strings"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig"
	multisigtypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nftkeeper "bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
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
	reserve := sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(10)))
	nftGenesisState := nfttypes.GenesisState{
		Params: nfttypes.DefaultParams(),
		Collections: []nfttypes.Collection{
			{
				Creator: newAddress,
				Denom:   "aaa",
				Supply:  2,
				Tokens: []*nfttypes.Token{
					{
						Creator:   newAddress,
						Denom:     "aaa",
						ID:        "a1",
						URI:       "a1",
						Reserve:   reserve,
						AllowMint: true,
						Minted:    1,
						Burnt:     0,
						SubTokens: []*nfttypes.SubToken{
							{
								ID:      1,
								Owner:   oldAddress,
								Reserve: &reserve,
							},
						},
					},
					{
						Creator:   newAddress,
						Denom:     "aaa",
						ID:        "a2",
						URI:       "a2",
						Reserve:   reserve,
						AllowMint: true,
						Minted:    1,
						Burnt:     0,
						SubTokens: []*nfttypes.SubToken{
							{
								ID:      1,
								Owner:   oldAddress,
								Reserve: &reserve,
							},
						},
					},
				},
			},
		},
	}
	require.NoError(t, nftGenesisState.Validate(), "nftGenesisState")
	nftkeeper.InitGenesis(ctx, app.NFTKeeper, &nftGenesisState)

	multisigGenesisState := multisigtypes.GenesisState{
		Wallets: []multisigtypes.Wallet{
			{
				Address:   wallet1,
				Owners:    []string{otherAddress, oldAddress},
				Weights:   []uint32{1, 2},
				Threshold: 3,
			},
			{
				Address:   wallet2,
				Owners:    []string{otherAddress, oldAddress},
				Weights:   []uint32{1, 2},
				Threshold: 3,
			},
		},
	}
	require.NoError(t, nftGenesisState.Validate(), "nftGenesisState")
	multisig.InitGenesis(ctx, app.MultisigKeeper, &multisigGenesisState)

	legacyGenesisState := types.GenesisState{
		Records: []types.Record{
			{
				LegacyAddress: oldAddress,
				Coins:         sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(50)), sdk.NewCoin("foo", sdk.NewInt(1))),
				NFTs:          []string{"a2"},
				Wallets:       []string{wallet1},
			},
			{
				LegacyAddress: otherAddress,
				Wallets:       []string{wallet2},
			},
		},
	}

	require.NoError(t, legacyGenesisState.Validate(), "legacyGenesisState")
	legacy.InitGenesis(ctx, app.LegacyKeeper, &legacyGenesisState)

	// init  genesis done
	// let's check
	require.True(t, app.LegacyKeeper.IsLegacyAddress(ctx, oldAddress), "check legacy address 1")
	require.False(t, app.LegacyKeeper.IsLegacyAddress(ctx, newAddress), "check legacy address 2")
	err = app.LegacyKeeper.ActualizeLegacy(ctx, publicKey)
	require.NoError(t, err, "ActualizeLegacy")

	// coins
	legacyRemain := app.BankKeeper.GetAllBalances(ctx, cosmosAuthTypes.NewModuleAddress(types.LegacyCoinPool))
	require.True(t, legacyRemain.IsEqual(sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(50)), sdk.NewCoin("foo", sdk.NewInt(99)))),
		"legacy coins remain")

	// nft no changes
	subs := app.NFTKeeper.GetSubTokens(ctx, "a1")
	require.Len(t, subs, 1, "nft-a1")
	require.Equal(t, oldAddress, subs[0].Owner, "nft-a1 must not changed")
	//  nft changes
	subs = app.NFTKeeper.GetSubTokens(ctx, "a2")
	require.Len(t, subs, 1, "nft-a2")
	require.Equal(t, newAddress, subs[0].Owner, "nft-a2 must changed")

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
	//require.False(t, app.LegacyKeeper.IsLegacyAddress(ctx, oldAddress), "check legacy address at end")
}

func getBaseApp(t *testing.T) (*app.DSC, sdk.Context) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	return dsc, ctx
}
