package coin_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/testcoin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	//tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func bootstrapGenesisTest() (*app.DSC, sdk.Context) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	return dsc, ctx
}

var (
	atomCoin = types.Coin{
		Title:       "Cosmos Hub Atom",
		Symbol:      "ATOM",
		CRR:         50,
		Reserve:     sdk.NewInt(1_000_000_000),
		Volume:      helpers.EtherToWei(sdk.NewInt(1000000000000)),
		LimitVolume: sdk.NewInt(1_000_000_000_000_000),
		Creator:     "uatom",
		Identity:    "dx1hs2wdrm87c92rzhq0vgmgrxr6u57xpr2lcygc2",
	}
	tstCoin = types.Coin{
		Title:       "Test Suite Token",
		Symbol:      "TST",
		CRR:         100,
		Reserve:     sdk.NewInt(1_000_000_000),
		Volume:      sdk.NewInt(1_000_000_000_0),
		LimitVolume: sdk.NewInt(1_000_000_000_000_000_000),
		Creator:     "uatom",
		Identity:    "dx1hs2wdrm87c92rzhq0vgmgrxr6u57xpr2lcygc2",
	}

	check1 types.Check
	check2 types.Check
)

func TestAppModuleBasic_InitGenesis(t *testing.T) {
	app, ctx := bootstrapGenesisTest()

	// write genesis
	params := app.CoinKeeper.GetParams(ctx)

	coins := []types.Coin{
		atomCoin,
		tstCoin,
	}
	check1 = testcoin.CreateNewCheck(ctx.ChainID(), "100000del", "9", "", 1)
	check2 = testcoin.CreateNewCheck(ctx.ChainID(), "100000del", "10", "", 1)
	checks := []types.Check{
		check1,
		check2,
	}

	genesisState := types.NewGenesisState(params, coins, checks, []types.LegacyBalance{})
	coin.InitGenesis(ctx, app.CoinKeeper, genesisState)

	// export genesis

	coins = append(coins, types.Coin{
		Title:       params.BaseTitle,
		Symbol:      params.BaseSymbol,
		Volume:      params.BaseInitialVolume,
		CRR:         0,
		Reserve:     sdk.NewInt(0),
		Creator:     "",
		LimitVolume: sdk.NewInt(0),
		Identity:    "",
	})

	exportedGenesis := coin.ExportGenesis(ctx, app.CoinKeeper)
	require.Equal(t, params, exportedGenesis.Params)
	require.True(t, testcoin.CoinsEqual(coins, exportedGenesis.Coins))
	require.True(t, testcoin.ChecksEqual(checks, exportedGenesis.Checks))
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*types.GenesisState)
		wantErr bool
	}{
		{"default", func(*types.GenesisState) {}, false},
		// validate genesis validators
		{"params coin title > 64 symbols", func(data *types.GenesisState) {
			data.Params.BaseTitle = "vsafa;jkdfndsj;anf;asdnf;dsjfkldasfkmsdkalmf;alkdsmflmasl;dkmf;lds"
		}, true},
		{"params symbol is regexp", func(data *types.GenesisState) {
			data.Params.BaseSymbol = "laSK;DM"
		}, true},
		{"params init volume is < min", func(data *types.GenesisState) {
			data.Params.BaseInitialVolume = sdk.NewInt(0)
		}, true},
		{"params init volume is > max", func(data *types.GenesisState) {
			data.Params.BaseInitialVolume = helpers.EtherToWei(sdk.NewInt(1000000000000002))
		}, true},
		{"valid coins is repeated", func(data *types.GenesisState) {
			data.Coins = append(data.Coins, data.Coins[0])
		}, true},
		{"valid checks is repeated", func(data *types.GenesisState) {
			data.Checks = append(data.Checks, data.Checks[0])
		}, true},
		//{"invalid coin", func(data *types.GenesisState) {
		//	data.Coins = append(data.Coins, types.Coin{
		//		Title:       "vsafa;jkdfndsj;anf;asdnf;dsjfkldasfkmsdkalmf;alkdsmflmasl;dkmf;lds",
		//		Symbol:      "sdjfn",
		//		Volume:      helpers.EtherToWei(sdk.NewInt(1000000000000002)),
		//		CRR:         0,
		//		Reserve:     sdk.NewInt(0),
		//		Creator:     "",
		//		LimitVolume: sdk.NewInt(0),
		//		Identity:    "",
		//	})
		//}, true},
		//{"invalid check", func(data *types.GenesisState) {
		//	data.Checks = append(data.Checks, types.Check{})
		//}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			genesisState := types.DefaultGenesisState()
			genesisState.Coins = append(genesisState.Coins, atomCoin)
			genesisState.Checks = append(genesisState.Checks, check1)

			tt.mutate(genesisState)
			if tt.wantErr {
				require.Error(t, genesisState.Validate())
			} else {
				require.NoError(t, genesisState.Validate())
			}
		})
	}
}

func TestInitGenesisForLegacy(t *testing.T) {
	app, ctx := bootstrapGenesisTest()

	// write genesis
	params := app.CoinKeeper.GetParams(ctx)

	coins := []types.Coin{
		{
			Title:  "del",
			Symbol: "del",
		},
		{
			Title:       "Test coin",
			Symbol:      "foo",
			CRR:         50,
			Reserve:     sdk.NewInt(1_000_000_000),
			Volume:      sdk.NewInt(1_000_000_000_0),
			LimitVolume: sdk.NewInt(1_000_000_000_000_000_000),
			Creator:     "uatom",
			Identity:    "dx1hs2wdrm87c92rzhq0vgmgrxr6u57xpr2lcygc2",
		},
	}

	legacyCoinPoolAddress, err := sdk.Bech32ifyAddressBytes(config.Bech32Prefix, cosmosAuthTypes.NewModuleAddress(types.LegacyCoinPool))
	require.NoError(t, err, "legacyCoinPoolAddress to bech32")

	otherAddress, err := sdk.Bech32ifyAddressBytes(config.Bech32Prefix, sdk.AccAddress("someotheraddressfortest"))
	require.NoError(t, err, "other address to bech32")

	bankGenesisState := &cosmosBankTypes.GenesisState{
		Params: cosmosBankTypes.DefaultParams(),
		Balances: []cosmosBankTypes.Balance{
			{
				Address: legacyCoinPoolAddress,
				Coins: sdk.Coins{
					{
						Denom:  params.BaseSymbol,
						Amount: sdk.NewInt(100),
					},
					{
						Denom:  "foo",
						Amount: sdk.NewInt(100),
					},
				},
			},
			{
				Address: otherAddress,
				Coins: sdk.Coins{
					{
						Denom:  params.BaseSymbol,
						Amount: sdk.NewInt(200),
					},
					{
						Denom:  "foo",
						Amount: sdk.NewInt(200),
					},
				},
			},
		},
	}
	require.NoError(t, bankGenesisState.Validate(), "bankGenesisState")
	app.BankKeeper.InitGenesis(ctx, bankGenesisState)

	const oldAddress = "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy"
	coinGenesisState := types.NewGenesisState(params, coins, []types.Check{}, []types.LegacyBalance{
		{
			LegacyAddress: oldAddress,
			Coins: sdk.Coins{
				{
					Denom:  "foo",
					Amount: sdk.NewInt(100),
				},
				{
					Denom:  "del",
					Amount: sdk.NewInt(150),
				},
			},
		},
		{
			// second address to check iterator
			LegacyAddress: "dx1lw2q66zph22x3hzmc527em25kd4zfydnx7arw7",
			Coins: sdk.Coins{
				{
					Denom:  "foo",
					Amount: sdk.NewInt(1),
				},
				{
					Denom:  "del",
					Amount: sdk.NewInt(1),
				},
			},
		},
	})
	require.NoError(t, coinGenesisState.Validate(), "coinGenesisState")
	coin.InitGenesis(ctx, app.CoinKeeper, coinGenesisState)

	// here bank and coin must be proper initialized
	// so we check state
	coinKeeper := app.CoinKeeper
	balance, err := coinKeeper.GetLegacyBalance(ctx, oldAddress)
	require.NoError(t, err, "GetLegacyBalance")
	require.Equal(t, oldAddress, balance.LegacyAddress)
	require.Equal(t, 2, len(balance.Coins))
	for _, coin := range balance.Coins {
		if coin.Denom == "del" {
			require.True(t, coin.Amount.Equal(sdk.NewInt(150)), "balance for del")
		}
		if coin.Denom == "foo" {
			require.True(t, coin.Amount.Equal(sdk.NewInt(100)), "balance for foo")
		}
	}

	// otherAddress must be full
	otherCoins := app.BankKeeper.GetAllBalances(ctx, sdk.AccAddress("someotheraddressfortest"))
	require.Equal(t, 2, len(otherCoins))

	// chekc itertor by GetLegacyBalances
	balances := coinKeeper.GetLegacyBalances(ctx)
	require.Equal(t, 2, len(balances), "something wrong in iterator")
}
