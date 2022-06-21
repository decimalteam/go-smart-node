package coin_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
)

// getBaseAppWithCustomKeeper Returns a simapp with custom CoinKeeper
// to avoid messing with the hooks.
func getBaseAppWithCustomKeeper() (*codec.LegacyAmino, *app.DSC, sdk.Context) {
	dsc := app.Setup(false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	dsc.CoinKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.AccountKeeper,
		dsc.BankKeeper,
	)
	dsc.CoinKeeper.SetParams(ctx, types.DefaultParams())

	return codec.NewLegacyAmino(), dsc, ctx
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}

// create invalid coin
func invalidCoin() sdk.Coin {
	return sdk.Coin{
		Denom:  "invalidDenom",
		Amount: sdk.NewInt(100000000),
	}
}

// create valid coin
func validCoin(denom string, amount int64) sdk.Coin {
	return sdk.Coin{
		Denom:  denom,
		Amount: helpers.EtherToWei(sdk.NewInt(amount)),
	}
}
