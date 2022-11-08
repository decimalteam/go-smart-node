package fee_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"github.com/stretchr/testify/require"
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
