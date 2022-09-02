package fee_test

import (
	"math/big"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
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
		dsc.BaseApp,
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

func TestFeeLimitByConsensus(t *testing.T) {
	dsc := app.Setup(t, false, nil)
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	dsc.FeeKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.BankKeeper,
		dsc.BaseApp,
		config.BaseDenom,
	)

	gs := *types.DefaultGenesisState()
	require.NoError(t, gs.Validate())
	fee.InitGenesis(ctx, dsc.FeeKeeper, gs)

	msgHandler := fee.NewHandler(dsc.FeeKeeper)

	dsc.BaseApp.StoreConsensusParams(ctx, &abci.ConsensusParams{
		Block: &abci.BlockParams{
			MaxGas: 20202020,
		},
	})

	// 1. price 1.0
	{
		price, err := sdk.NewDecFromStr("1.0")
		require.NoError(t, err)
		msg := types.NewMsgSaveBaseDenomPrice(gs.Params.OracleAddress, "del", price)
		_, err = msgHandler(ctx, msg)
		require.NoError(t, err)
		// Check fee. Must be default 1_000_000_000
		baseFee := dsc.FeeKeeper.GetBaseFee(ctx)
		require.Equal(t, 0, baseFee.Cmp(big.NewInt(1_000_000_000)))
	}
	// 2. price 0.04
	{
		price, err := sdk.NewDecFromStr("0.04")
		require.NoError(t, err)
		msg := types.NewMsgSaveBaseDenomPrice(gs.Params.OracleAddress, "del", price)
		_, err = msgHandler(ctx, msg)
		require.NoError(t, err)
		// Check fee. Must be limited by consensus max gas
		baseFee := dsc.FeeKeeper.GetBaseFee(ctx)
		require.Equal(t, 0, baseFee.Cmp(big.NewInt(20_202_020_000)))
	}

	// 3. null block in consensus param
	{
		dsc.BaseApp.StoreConsensusParams(ctx, &abci.ConsensusParams{
			Block: nil,
		})
		price, err := sdk.NewDecFromStr("0.04")
		require.NoError(t, err)
		msg := types.NewMsgSaveBaseDenomPrice(gs.Params.OracleAddress, "del", price)
		_, err = msgHandler(ctx, msg)
		require.NoError(t, err)
		dsc.BaseApp.StoreConsensusParams(ctx, nil)
		// Check fee. Must be according to price 0.04
		baseFee := dsc.FeeKeeper.GetBaseFee(ctx)
		require.Equal(t, 0, baseFee.Cmp(big.NewInt(25_000_000_000)), baseFee.String())
	}
}
