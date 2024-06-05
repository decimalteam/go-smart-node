package keeper

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func GetBaseAppWithCustomKeeper(t *testing.T) (*app.DSC, sdk.Context) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	dsc.NFTKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(nfttypes.ModuleName),
		dsc.BankKeeper,
		&dsc.EvmKeeper,
	)

	return dsc, ctx
}
