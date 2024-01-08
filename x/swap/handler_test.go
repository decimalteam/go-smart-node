package swap_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/swap"
	"bitbucket.org/decimalteam/go-smart-node/x/swap/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestChainOperations(t *testing.T) {
	dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, 10,
		sdk.NewCoins(
			sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1000))),
		),
	)

	gs := types.DefaultGenesisState()

	swap.InitGenesis(ctx, dsc.SwapKeeper, gs)

	swapService, err := sdk.AccAddressFromBech32(gs.Params.ServiceAddress)
	require.NoError(t, err)

	// invalid sender to activate chain
	msg := types.NewMsgActivateChain(addrs[0], 1, "some chain")
	err = msg.ValidateBasic()
	require.NoError(t, err)
	goCtx := sdk.WrapSDKContext(ctx)
	_, err = dsc.SwapKeeper.ActivateChain(goCtx, msg)
	require.Error(t, err)

	// valid sender to activate chain
	msg = types.NewMsgActivateChain(swapService, 1, "some chain")
	err = msg.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = dsc.SwapKeeper.ActivateChain(goCtx, msg)
	require.NoError(t, err)

	chain, ok := dsc.SwapKeeper.GetChain(ctx, 1)
	require.True(t, ok)
	require.Equal(t, types.Chain{Id: 1, Name: "some chain", Active: true}, chain)

	// invalid sender to deactivate chain
	msgD := types.NewMsgDeactivateChain(addrs[0], 1)
	err = msgD.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = dsc.SwapKeeper.DeactivateChain(goCtx, msgD)
	require.Error(t, err)

	// valid sender to deactivate chain
	msgD = types.NewMsgDeactivateChain(swapService, 1)
	err = msg.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = dsc.SwapKeeper.DeactivateChain(goCtx, msgD)
	require.NoError(t, err)

	chain, ok = dsc.SwapKeeper.GetChain(ctx, 1)
	require.True(t, ok)
	require.Equal(t, types.Chain{Id: 1, Name: "some chain", Active: false}, chain)
}

// getBaseAppWithCustomKeeper Returns a simapp with custom keepers
// to avoid messing with the hooks.
func getBaseAppWithCustomKeeper(t *testing.T) (*app.DSC, sdk.Context) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	dsc.SwapKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.AccountKeeper,
		dsc.BankKeeper,
	)

	return dsc, ctx
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}
