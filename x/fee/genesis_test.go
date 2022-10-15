package fee_test

//
//func TestDefaultGenesis(t *testing.T) {
//	dsc := app.Setup(t, false, nil)
//	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})
//
//	appCodec := dsc.AppCodec()
//
//	dsc.FeeKeeper = *keeper.NewKeeper(
//		appCodec,
//		dsc.GetKey(types.StoreKey),
//		dsc.GetSubspace(types.ModuleName),
//		dsc.BankKeeper,
//		config.BaseDenom,
//	)
//
//	fee.InitGenesis(ctx, dsc.FeeKeeper, types.DefaultGenesisState())
//
//	params := dsc.FeeKeeper.GetModuleParams(ctx)
//	price, err := dsc.FeeKeeper.GetPrice(ctx, "del", "usd")
//	require.NoError(t, err)
//
//	gs := types.DefaultGenesisState()
//	// check proper genesis initialization
//	require.Equal(t, types.DefaultParams(), params)
//	require.True(t, price.Price.Equal(gs.Prices[0].Price))
//}

//func TestGenesisInit(t *testing.T) {
//	dsc := app.Setup(t, false, nil)
//	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})
//
//	appCodec := dsc.AppCodec()
//
//	dsc.FeeKeeper = *keeper.NewKeeper(
//		appCodec,
//		dsc.GetKey(types.StoreKey),
//		dsc.GetSubspace(types.ModuleName),
//		dsc.BankKeeper,
//		config.BaseDenom,
//	)
//
//	gs := types.DefaultGenesisState()
//	gs.Params.Oracle = "dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd"
//	require.NoError(t, gs.Validate())
//	fee.InitGenesis(ctx, dsc.FeeKeeper, gs)
//
//	params := dsc.FeeKeeper.GetModuleParams(ctx)
//	price, err := dsc.FeeKeeper.GetPrice(ctx, "del", "usd")
//	require.NoError(t, err)
//
//	// check proper genesis initialization
//	require.Equal(t, "dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd", params.Oracle)
//	require.True(t, price.Price.Equal(gs.Prices[0].Price))
//}
