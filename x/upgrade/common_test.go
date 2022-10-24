package upgrade_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
)

// getBaseAppW  ithCustomKeeper Returns a simapp with custom CoinKeeper
// to avoid messing with the hooks.
//func getBaseAppWithCustomKeeper(skip map[int64]bool) (*codec.LegacyAmino, *app.DSC, sdk.Context) {
//	dsc := app.Setup(false, feemarkettypes.DefaultGenesisState())
//	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})
//
//	appCodec := dsc.AppCodec()
//
//	if skip == nil {
//		skip = make(map[int64]bool)
//	}
//	dsc.UpgradeKeeper = keeper.NewKeeper(
//		skip,
//		dsc.GetKey(types.StoreKey),
//		appCodec,
//		app.DefaultNodeHome,
//		dsc.BaseApp,
//		app.UpgraderAddress,
//	)
//
//	return codec.NewLegacyAmino(), dsc, ctx
//}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}
