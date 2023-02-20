package coin_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/testcoin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

func bootstrapGenesisTest(t *testing.T) (*app.DSC, sdk.Context) {
	_, dsc, ctx := testkeeper.GetTestAppWithCoinKeeper(t)

	return dsc, ctx
}

var (
	atomCoin = types.Coin{
		Denom:       "atom",
		Title:       "Cosmos Hub Atom",
		Creator:     "uatom",
		CRR:         50,
		LimitVolume: sdkmath.NewInt(1_000_000_000_000_000),
		MinVolume:   sdkmath.ZeroInt(),
		Identity:    "d01hs2wdrm87c92rzhq0vgmgrxr6u57xpr2ml8an4",
		Reserve:     sdkmath.NewInt(1_000_000_000),
		Volume:      helpers.EtherToWei(sdkmath.NewInt(1000000000000)),
	}
	tstCoin = types.Coin{
		Denom:       "tst",
		Title:       "Test Suite Token",
		Creator:     "uatom",
		CRR:         100,
		LimitVolume: sdkmath.NewInt(1_000_000_000_000_000_000),
		MinVolume:   sdkmath.ZeroInt(),
		Identity:    "d01hs2wdrm87c92rzhq0vgmgrxr6u57xpr2ml8an4",
		Reserve:     sdkmath.NewInt(1_000_000_000),
		Volume:      sdkmath.NewInt(1_000_000_000_0),
	}

	check1 types.Check
	check2 types.Check
)

func TestAppModuleBasic_InitGenesis(t *testing.T) {
	app, ctx := bootstrapGenesisTest(t)

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

	genesisState := types.NewGenesisState(params, coins, checks)
	coin.InitGenesis(ctx, app.CoinKeeper, genesisState)

	// export genesis

	coins = append(coins, types.Coin{
		Denom:       params.BaseDenom,
		Title:       params.BaseTitle,
		Creator:     "",
		CRR:         0,
		LimitVolume: sdkmath.NewInt(0),
		MinVolume:   sdkmath.ZeroInt(),
		Identity:    "",
		Volume:      params.BaseVolume,
		Reserve:     sdkmath.NewInt(0),
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
		{"params coin is regexp", func(data *types.GenesisState) {
			data.Params.BaseDenom = "laSK;DM"
		}, true},
		{"params coin title > 64 characters", func(data *types.GenesisState) {
			data.Params.BaseTitle = "vsafa;jkdfndsj;anf;asdnf;dsjfkldasfkmsdkalmf;alkdsmflmasl;dkmf;lds"
		}, true},
		{"params volume is < min", func(data *types.GenesisState) {
			data.Params.BaseVolume = sdkmath.NewInt(0)
		}, true},
		{"params volume is > max", func(data *types.GenesisState) {
			data.Params.BaseVolume = helpers.EtherToWei(sdkmath.NewInt(1000000000000002))
		}, true},
		{"valid coins is repeated", func(data *types.GenesisState) {
			data.Coins = append(data.Coins, data.Coins[0])
		}, true},
		{"valid checks is repeated", func(data *types.GenesisState) {
			data.Checks = append(data.Checks, data.Checks[0])
		}, true},
		//{"invalid coin", func(data *types.GenesisState) {
		//	data.Coins = append(data.Coins, types.Coin{
		//		Denom:       "sdjfn",
		//		Title:       "vsafa;jkdfndsj;anf;asdnf;dsjfkldasfkmsdkalmf;alkdsmflmasl;dkmf;lds",
		//		Creator:     "",
		//		CRR:         0,
		//		LimitVolume: sdkmath.NewInt(0),
		//		Identity:    "",
		//		Volume:      helpers.EtherToWei(sdkmath.NewInt(1000000000000002)),
		//		Reserve:     sdkmath.NewInt(0),
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
