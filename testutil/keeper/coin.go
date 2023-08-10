package keeper

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	ethsrvflags "github.com/evmos/ethermint/server/flags"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	evmgeth "github.com/evmos/ethermint/x/evm/vm/geth"
	"github.com/spf13/cast"
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

func GetTestAppWithCoinKeeper(t *testing.T) (*codec.LegacyAmino, *app.DSC, sdk.Context) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	appOpts := simapp.EmptyAppOptions{}

	dsc.EvmKeeper = evmkeeper.NewKeeper(
		appCodec,
		dsc.GetKey(evmtypes.StoreKey),
		dsc.GetKey(evmtypes.TransientKey),
		dsc.GetSubspace(evmtypes.ModuleName),
		dsc.AccountKeeper,
		dsc.BankKeeper,
		&dsc.ValidatorKeeper,
		dsc.FeeKeeper,
		nil,
		evmgeth.NewEVM,
		cast.ToString(appOpts.Get(ethsrvflags.EVMTracer)),
	)

	dsc.CoinKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.AccountKeeper,
		&dsc.FeeKeeper,
		dsc.BankKeeper,
		dsc.EvmKeeper,
	)
	dsc.CoinKeeper.SetParams(ctx, types.DefaultParams())

	return codec.NewLegacyAmino(), dsc, ctx
}

// GenerateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func GenerateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}
