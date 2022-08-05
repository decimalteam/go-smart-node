package fee_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestSavePrice(t *testing.T) {
	dsc := app.Setup(false, nil)
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	dsc.FeeKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.BankKeeper,
		config.BaseDenom,
	)

	gs := *types.DefaultGenesisState()
	require.NoError(t, gs.Validate())
	fee.InitGenesis(ctx, dsc.FeeKeeper, gs)

	msgHandler := fee.NewHandler(dsc.FeeKeeper)

	// 1. invalid sender, must be error
	msg := types.NewMsgSaveBaseDenomPrice(gs.Params.OracleAddress+"0", "del", sdk.NewDec(2))
	_, err := msgHandler(ctx, msg)
	require.Error(t, err)

	// 2. valid, must be no error
	msg = types.NewMsgSaveBaseDenomPrice(gs.Params.OracleAddress, "del", sdk.NewDec(2))
	_, err = msgHandler(ctx, msg)
	require.NoError(t, err)
	// check saving
	price, err := dsc.FeeKeeper.GetPrice(ctx)
	require.NoError(t, err)
	require.True(t, price.Equal(sdk.NewDec(2)))
}
