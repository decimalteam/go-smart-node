package fee_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/app/ante"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/fee"
	feeconfig "bitbucket.org/decimalteam/go-smart-node/x/fee/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

func TestSavePrice(t *testing.T) {
	dsc := app.Setup(t, false, nil)
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	dsc.FeeKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.BankKeeper,
		&dsc.CoinKeeper,
		dsc.AccountKeeper,
		config.BaseDenom,
		ante.CalculateFee,
	)

	gs := types.DefaultGenesisState()
	require.NoError(t, gs.Validate())
	fee.InitGenesis(ctx, dsc.FeeKeeper, gs)

	msgHandler := fee.NewHandler(dsc.FeeKeeper)

	prices := []types.CoinPrice{
		{
			Denom: "del",
			Quote: "usd",
			Price: sdk.NewDec(2),
		},
		{
			Denom: "del",
			Quote: "rub",
			Price: sdk.NewDec(2),
		},
	}
	// 1. invalid sender, must be error
	msg := types.NewMsgUpdateCoinPrices(gs.Params.Oracle+"0", prices)
	_, err := msgHandler(ctx, msg)
	require.Error(t, err)

	// 2. valid, must be no error
	msg = types.NewMsgUpdateCoinPrices(gs.Params.Oracle, prices)
	_, err = msgHandler(ctx, msg)
	require.NoError(t, err)
	// check saving
	storedPrices, err := dsc.FeeKeeper.GetPrices(ctx)
	require.NoError(t, err)
	require.Len(t, storedPrices, 2)
}

func TestCommissionCalculation(t *testing.T) {
	dsc := app.Setup(t, false, nil)
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	accs, _ := generateAddresses(dsc, ctx, 10, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))))

	appCodec := dsc.AppCodec()

	dsc.FeeKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.BankKeeper,
		&dsc.CoinKeeper,
		dsc.AccountKeeper,
		config.BaseDenom,
		ante.CalculateFee,
	)

	dsc.CoinKeeper.SetCoin(ctx,
		cointypes.Coin{
			Denom:       "testcoin",
			Title:       "testcoin",
			Creator:     accs[0].String(),
			CRR:         100,
			LimitVolume: helpers.EtherToWei(sdkmath.NewInt(10000)),
			Volume:      helpers.EtherToWei(sdkmath.NewInt(1000)),
			Reserve:     helpers.EtherToWei(sdkmath.NewInt(2000)),
		},
	)

	txHexBytes1 := "0a9c010a92010a1c2f646563696d616c2e636f696e2e76312e4d736753656e64436f696e12720a29647831746c796b79786e337a6464776d37773839726175727775767761356170763477333274683066122964783130647461766570683271303378333234346475766d643932676b7767796c6c35726c756c6d6e1a1a0a0364656c121331303030303030303030303030303030303030120568656c6c6f125b0a570a4f0a282f65746865726d696e742e63727970746f2e76312e657468736563703235366b312e5075624b657912230a2103915d3a632aaec661cc693adb5341a5f104661e6f7a85db9df1d8a7a332f781fe12040a02080112001a4159dc3cc63526e1a66e5ab6748ad5500f313e0553ecebe7e5fb8bfc34bccd63ed57735101a28d4fb1b8e1ebd83cb60c32c98100ee9d1858f1a4ccbb10475ee44400"

	params := dsc.FeeKeeper.GetModuleParams(ctx)
	delPrice, err := dsc.FeeKeeper.GetPrice(ctx, config.BaseDenom, feeconfig.DefaultQuote)
	require.NoError(t, err)
	comm, err := ante.CalculateFee(
		dsc.AppCodec(),
		[]sdk.Msg{cointypes.NewMsgSendCoin(
			sdk.MustAccAddressFromBech32("dx1tlykyxn3zddwm7w89raurwuvwa5apv4w32th0f"),
			sdk.MustAccAddressFromBech32("dx10dtaveph2q03x3244duvmd92gkwgyll5rlulmn"),
			sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1))),
		)},
		int64(len(txHexBytes1)/2),
		delPrice.Price,
		params,
	)

	goCtx := sdk.WrapSDKContext(ctx)

	// valid tx, base coin
	resp, err := dsc.FeeKeeper.CalculateCommission(goCtx, &types.QueryCalculateCommissionRequest{
		TxBytes: txHexBytes1,
		Denom:   config.BaseDenom,
	})
	require.NoError(t, err)
	require.False(t, resp.Commission.IsNil())
	require.False(t, resp.Commission.IsZero())
	require.True(t, resp.Commission.Equal(comm), "%s != %s", resp.Commission, comm)

	// valid tx, custom coin
	resp, err = dsc.FeeKeeper.CalculateCommission(goCtx, &types.QueryCalculateCommissionRequest{
		TxBytes: txHexBytes1,
		Denom:   "testcoin",
	})
	require.NoError(t, err)
	require.False(t, resp.Commission.IsNil())
	require.False(t, resp.Commission.IsZero())
	require.True(t, resp.Commission.Equal(comm.QuoRaw(2)))

	// valid tx, coin isn't exist
	resp, err = dsc.FeeKeeper.CalculateCommission(goCtx, &types.QueryCalculateCommissionRequest{
		TxBytes: txHexBytes1,
		Denom:   "notexists",
	})
	require.Error(t, err)

	// invalid transaction bytes
	txHexBytes2 := "0a9c010a92010a1c2f646563696d616c2e636f696e2e76312e4d736753656e64436f696e12720a29647831746c796b79786e337a6464776d37773839726175727775767761356170763477333274683066122964783130647461766570683271303378333234346475766d643932676b7767796c6c35726c756c6d6e1a1a0a0364656c121331303030303030303030303030303030303030120568656c6c6f125b0a570a4f0a282f65746865726d696e742e63727970746f2e76312e657468736563703235366b312e5075624b657912230a2103915d3a632aaec661cc693adb5341a5f104661e6f7a85db9df1d8a7a332f781fe12040a02080112001a4159dc3cc63526e1a66e5ab6748ad5500f313e0553ecebe7e5fb8bfc34bccd63ed57735101a28d4fb1b8e1ebd83cb60c32c98100"
	resp, err = dsc.FeeKeeper.CalculateCommission(goCtx, &types.QueryCalculateCommissionRequest{
		TxBytes: txHexBytes2,
		Denom:   config.BaseDenom,
	})
	require.Error(t, err)
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}
