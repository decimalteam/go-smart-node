package coin_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	appMain "bitbucket.org/decimalteam/go-smart-node/app"
	appAnte "bitbucket.org/decimalteam/go-smart-node/app/ante"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func TestConsistency(t *testing.T) {
	baseVolume := helpers.FinneyToWei(sdk.NewInt(100_000))   // = is sum 200_000 because 2 addresses
	baseReserve := helpers.FinneyToWei(sdk.NewInt(1000_000)) // = 1000 del
	limitVolume := helpers.FinneyToWei(sdk.NewInt(300_000))  // 3*baseVolume
	crr := uint64(99)

	app, ctx, adrs := initConsistencyApp(t, baseReserve, baseVolume, limitVolume, crr)
	runOpSequence(t, app, ctx, []coinOp{
		{opType: "sell", adr: adrs[0], amount: baseVolume},
		{opType: "sell", adr: adrs[1], amount: baseVolume},
	})

	app, ctx, adrs = initConsistencyApp(t, baseReserve, baseVolume, limitVolume, crr)
	runOpSequence(t, app, ctx, []coinOp{
		{opType: "buy", adr: adrs[0], amount: baseVolume},
		{opType: "buy", adr: adrs[1], amount: baseVolume},
	})

	app, ctx, adrs = initConsistencyApp(t, baseReserve, baseVolume, limitVolume, crr)
	runOpSequence(t, app, ctx, []coinOp{
		{opType: "sellAll", adr: adrs[0], amount: baseVolume},
		{opType: "sellAll", adr: adrs[1], amount: baseVolume},
	})

	app, ctx, adrs = initConsistencyApp(t, baseReserve, baseVolume, limitVolume, crr)
	runOpSequence(t, app, ctx, []coinOp{
		{opType: "fee", adr: adrs[0], amount: helpers.FinneyToWei(sdk.NewInt(100))},
		{opType: "fee", adr: adrs[1], amount: helpers.FinneyToWei(sdk.NewInt(100))},
		{opType: "sellAll", adr: adrs[0], amount: baseVolume},
		{opType: "sellAll", adr: adrs[1], amount: baseVolume},
	})

	app, ctx, adrs = initConsistencyApp(t, baseReserve, baseVolume, limitVolume, crr)
	runOpSequence(t, app, ctx, []coinOp{
		{opType: "fee", adr: adrs[0], amount: helpers.FinneyToWei(sdk.NewInt(100))},
		{opType: "fee", adr: adrs[1], amount: helpers.FinneyToWei(sdk.NewInt(100))},
		{opType: "sellAll", adr: adrs[0], amount: baseVolume},
		{opType: "sellAll", adr: adrs[1], amount: baseVolume},
		{opType: "validator", adr: nil, amount: sdk.ZeroInt()},
	})

	app, ctx, adrs = initConsistencyApp(t, baseReserve, baseVolume, limitVolume, crr)
	runOpSequence(t, app, ctx, []coinOp{
		{opType: "fee", adr: adrs[0], amount: helpers.FinneyToWei(sdk.NewInt(100))},
		{opType: "fee", adr: adrs[1], amount: helpers.FinneyToWei(sdk.NewInt(100))},
		{opType: "buy", adr: adrs[0], amount: baseVolume},
		{opType: "buy", adr: adrs[1], amount: baseVolume},
		{opType: "buy", adr: adrs[0], amount: baseVolume},
		{opType: "buy", adr: adrs[1], amount: baseVolume},
		{opType: "validator", adr: nil, amount: sdk.ZeroInt()},
		{opType: "sellAll", adr: adrs[0], amount: baseVolume},
		{opType: "sellAll", adr: adrs[1], amount: baseVolume},
	})

	app, ctx, adrs = initConsistencyApp(t, baseReserve, baseVolume, limitVolume, crr)
	runOpSequence(t, app, ctx, []coinOp{
		{opType: "fee", adr: adrs[0], amount: helpers.FinneyToWei(sdk.NewInt(100))},
		{opType: "fee", adr: adrs[1], amount: helpers.FinneyToWei(sdk.NewInt(100))},
		{opType: "buy", adr: adrs[0], amount: baseVolume},
		{opType: "buy", adr: adrs[1], amount: baseVolume},
		{opType: "buy", adr: adrs[0], amount: baseVolume},
		{opType: "buy", adr: adrs[1], amount: baseVolume},
		{opType: "validator", adr: nil, amount: sdk.ZeroInt()},
		{opType: "burn", adr: adrs[0], amount: baseVolume},
		{opType: "burn", adr: adrs[1], amount: baseVolume},
	})
}

type coinOp struct {
	opType string // fee, validator (=burn), buy, sell, sellAll
	adr    sdk.AccAddress
	amount sdk.Int
}

func runOpSequence(t *testing.T, app *appMain.DSC, ctx sdk.Context, seq []coinOp) {
	for i, op := range seq {
		fooCoin := sdk.NewCoin("foo", op.amount)
		switch op.opType {
		case "fee":
			appAnte.DeductFees(ctx, app.BankKeeper, &app.CoinKeeper, op.adr, fooCoin)
		case "validator":
			// fee burn (like validator module)
			coinInCollector := app.BankKeeper.GetBalance(ctx, sdkAuthTypes.NewModuleAddress(sdkAuthTypes.FeeCollectorName), "foo")
			coinInfo, err := app.CoinKeeper.GetCoin(ctx, "foo")
			require.NoError(t, err, "validator/GetCoin, step: %d", i)
			reserveDecrease := formulas.CalculateSaleReturn(coinInfo.Volume, coinInfo.Reserve, uint(coinInfo.CRR), coinInCollector.Amount)
			app.CoinKeeper.EditCoin(ctx, coinInfo, coinInfo.Reserve.Sub(reserveDecrease), coinInfo.Volume.Sub(coinInCollector.Amount))
			app.BankKeeper.BurnCoins(ctx, sdkAuthTypes.FeeCollectorName, sdk.NewCoins(coinInCollector))
		case "buy":
			app.CoinKeeper.BuyCoin(sdk.WrapSDKContext(ctx), types.NewMsgBuyCoin(
				op.adr,
				fooCoin,
				sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1000000))),
			))
		case "sell":
			app.CoinKeeper.SellCoin(sdk.WrapSDKContext(ctx), types.NewMsgSellCoin(
				op.adr,
				fooCoin,
				sdk.NewCoin("del", sdk.NewInt(0)),
			))
		case "sellAll":
			app.CoinKeeper.SellAllCoin(sdk.WrapSDKContext(ctx), types.NewMsgSellAllCoin(
				op.adr,
				fooCoin.Denom,
				sdk.NewCoin("del", sdk.NewInt(0)),
			))
		case "burn":
			app.CoinKeeper.BurnCoin(sdk.WrapSDKContext(ctx), types.NewMsgBurnCoin(
				op.adr,
				fooCoin,
			))
		}

		coinInfo, err := app.CoinKeeper.GetCoin(ctx, "foo")
		require.NoError(t, err, "GetCoin, step: %d", i)
		require.NoError(t, checkCoin(coinInfo), "step: %d", i)
	}

}

func initConsistencyApp(t *testing.T, reserve, volume, limitVolume sdk.Int, crr uint64) (*app.DSC, sdk.Context, []sdk.AccAddress) {
	app, ctx := bootstrapGenesisTest(t)

	// write genesis
	params := app.CoinKeeper.GetParams(ctx)
	adr1, err := sdk.Bech32ifyAddressBytes(config.Bech32Prefix, []byte("adr1"))
	require.NoError(t, err, "adr1 to bech32")

	adr2, err := sdk.Bech32ifyAddressBytes(config.Bech32Prefix, []byte("adr2"))
	require.NoError(t, err, "adr2 to bech32")

	coins := []types.Coin{
		{
			Title:  "del",
			Symbol: "del",
		},
		{
			Title:       "Foo coin",
			Symbol:      "foo",
			CRR:         crr,
			Reserve:     reserve,
			Volume:      volume.Mul(sdk.NewInt(2)),
			LimitVolume: limitVolume,
			Creator:     adr1,
			Identity:    "foo",
		},
	}

	bankGenesisState := &cosmosBankTypes.GenesisState{
		Params: cosmosBankTypes.DefaultParams(),
		Balances: []cosmosBankTypes.Balance{
			{
				Address: adr1,
				Coins: sdk.Coins{
					{
						Denom:  params.BaseSymbol,
						Amount: helpers.EtherToWei(sdk.NewInt(1000000)),
					},
					{
						Denom:  "foo",
						Amount: helpers.EtherToWei(volume),
					},
				},
			},
			{
				Address: adr2,
				Coins: sdk.Coins{
					{
						Denom:  params.BaseSymbol,
						Amount: helpers.EtherToWei(sdk.NewInt(1000000)),
					},
					{
						Denom:  "foo",
						Amount: helpers.EtherToWei(volume),
					},
				},
			},
		},
	}
	require.NoError(t, bankGenesisState.Validate(), "bankGenesisState")
	app.BankKeeper.InitGenesis(ctx, bankGenesisState)

	coinGenesisState := types.NewGenesisState(params, coins, []types.Check{})
	require.NoError(t, coinGenesisState.Validate(), "coinGenesisState")
	coin.InitGenesis(ctx, app.CoinKeeper, coinGenesisState)

	return app, ctx, []sdk.AccAddress{sdk.AccAddress("adr1"), sdk.AccAddress("adr2")}
}

func checkCoin(coinInfo types.Coin) error {
	if coinInfo.Volume.LT(types.MinCoinSupply) {
		return types.ErrTxBreaksMinVolumeLimit(coinInfo.Volume.String(), types.MinCoinSupply.String())
	}
	if coinInfo.Volume.GT(coinInfo.LimitVolume) {
		return types.ErrTxBreaksVolumeLimit(coinInfo.Volume.String(), coinInfo.LimitVolume.String())
	}
	if coinInfo.Reserve.LT(types.MinCoinReserve) {
		return types.ErrTxBreaksMinReserveRule(types.MinCoinReserve.String(), coinInfo.Reserve.String())
	}
	return nil
}
