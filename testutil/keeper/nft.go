package keeper

//func GetBaseAppWithCustomKeeper(t *testing.T) (*app.DSC, sdk.Context) {
//	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
//	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})
//
//	appCodec := dsc.AppCodec()
//
//	dsc.NFTKeeper = *keeper.NewKeeper(
//		appCodec,
//		dsc.GetKey(types.StoreKey),
//		dsc.BankKeeper,
//		cmdcfg.BaseDenom,
//	)
//
//	return dsc, ctx
//}
