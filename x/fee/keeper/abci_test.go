package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/app/ante"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

//func TestBurningPool(t *testing.T) {
//	_, dsc, ctx := createTestInput(t)
//	accs, _ := generateAddresses(dsc, ctx, 5, sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100_000)))))
//
//	//
//	goCtx := sdk.WrapSDKContext(ctx)
//	_, err := dsc.CoinKeeper.CreateCoin(goCtx, cointypes.NewMsgCreateCoin(
//		accs[0],
//		"testdenom",
//		"title",
//		10,
//		helpers.EtherToWei(sdkmath.NewInt(2000)),
//		helpers.EtherToWei(sdkmath.NewInt(1000)),
//		helpers.EtherToWei(sdkmath.NewInt(2_000_000)),
//		sdkmath.ZeroInt(),
//		"",
//	))
//	require.NoError(t, err)
//
//	coinsToBurn := sdk.NewCoins(
//		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
//		sdk.NewCoin("testdenom", helpers.EtherToWei(sdkmath.NewInt(100))),
//	)
//
//	coinInfoBaseBefore, err := dsc.CoinKeeper.GetCoin(ctx, cmdcfg.BaseDenom)
//	require.NoError(t, err)
//	coinInfoTestBefore, err := dsc.CoinKeeper.GetCoin(ctx, "testdenom")
//	require.NoError(t, err)
//
//	err = dsc.BankKeeper.SendCoinsFromAccountToModule(ctx, accs[0], types.BurningPool, coinsToBurn)
//	require.NoError(t, err)
//	// burn coins in EndBlocker
//	keeper.EndBlocker(ctx, dsc.FeeKeeper, abci.RequestEndBlock{})
//
//	// check coins
//	coinInfoBaseAfter, err := dsc.CoinKeeper.GetCoin(ctx, cmdcfg.BaseDenom)
//	require.NoError(t, err)
//	coinInfoTestAfter, err := dsc.CoinKeeper.GetCoin(ctx, "testdenom")
//	require.NoError(t, err)
//	require.True(t, coinInfoBaseBefore.Volume.Sub(coinInfoBaseAfter.Volume).Equal(coinsToBurn.AmountOf(cmdcfg.BaseDenom)))
//	require.True(t, coinInfoTestBefore.Volume.Sub(coinInfoTestAfter.Volume).Equal(coinsToBurn.AmountOf("testdenom")))
//
//	// check pool
//	bpAddr := dsc.AccountKeeper.GetModuleAddress(types.BurningPool)
//	burnPoolBalance := dsc.BankKeeper.GetAllBalances(ctx, bpAddr)
//	require.True(t, burnPoolBalance.Empty())
//}

func createTestInput(t *testing.T) (*codec.LegacyAmino, *app.DSC, sdk.Context) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	dsc.FeeKeeper = *keeper.NewKeeper(
		dsc.AppCodec(),
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.BankKeeper,
		&dsc.CoinKeeper,
		dsc.AccountKeeper,
		ante.CalculateFee,
	)
	return dsc.LegacyAmino(), dsc, ctx
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}
