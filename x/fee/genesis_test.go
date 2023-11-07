package fee_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/app/ante"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

func TestDefaultGenesis(t *testing.T) {
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

	fee.InitGenesis(ctx, dsc.FeeKeeper, types.DefaultGenesisState())

	params := dsc.FeeKeeper.GetModuleParams(ctx)
	price, err := dsc.FeeKeeper.GetPrice(ctx, cmdcfg.BaseDenom, "usd")
	require.NoError(t, err)

	gs := types.DefaultGenesisState()
	// check proper genesis initialization
	require.Equal(t, types.DefaultParams(), params)
	require.True(t, price.Price.Equal(gs.Prices[0].Price))
}

func TestGenesisInit(t *testing.T) {
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
	gs.Params.Oracle = "d01qql8ag4cluz6r4dz28p3w00dnc9w8ueuak90zj"
	require.NoError(t, gs.Validate())
	fee.InitGenesis(ctx, dsc.FeeKeeper, gs)

	params := dsc.FeeKeeper.GetModuleParams(ctx)
	price, err := dsc.FeeKeeper.GetPrice(ctx, cmdcfg.BaseDenom, "usd")
	require.NoError(t, err)

	// check proper genesis initialization
	require.Equal(t, "d01qql8ag4cluz6r4dz28p3w00dnc9w8ueuak90zj", params.Oracle)
	require.True(t, price.Price.Equal(gs.Prices[0].Price))
}
