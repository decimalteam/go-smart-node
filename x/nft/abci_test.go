package nft_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
)

func getBaseAppWithCustomKeeper() (*codec.LegacyAmino, *app.DSC, sdk.Context) {
	dsc := app.Setup(false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	dsc.NFTKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.BankKeeper,
		config.BaseDenom,
	)

	return codec.NewLegacyAmino(), dsc, ctx
}
func getAddrs(dsc *app.DSC, ctx sdk.Context, number int) []sdk.AccAddress {
	addrs := app.AddTestAddrsIncremental(dsc, ctx, number, sdk.Coins{
		{
			Denom:  "del",
			Amount: helpers.EtherToWei(sdk.NewInt(1000000000000)),
		},
	})

	return addrs
}
