package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func TestNegativePower(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	params := dsc.ValidatorKeeper.GetParams(ctx)
	params.MaxEntries = 100
	dsc.ValidatorKeeper.SetParams(ctx, params)

	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))

	goCtx := sdk.WrapSDKContext(ctx)

	// 1. create coin
	msgCreateCoin := cointypes.NewMsgCreateCoin(accs[0], "negative", "negative", 10,
		helpers.EtherToWei(sdkmath.NewInt(1_000)),
		helpers.EtherToWei(sdkmath.NewInt(1_000)),
		helpers.EtherToWei(sdkmath.NewInt(10_000)),
		sdkmath.ZeroInt(),
		"negative")
	_, err := dsc.CoinKeeper.CreateCoin(goCtx, msgCreateCoin)
	require.NoError(t, err)

	// 2. create validator
	creatorStake := sdk.NewCoin("negative", helpers.EtherToWei(sdkmath.NewInt(1_000)))
	msgCreate, err := types.NewMsgCreateValidator(vals[0], accs[0], PKs[0], types.Description{Moniker: "monik"},
		sdk.ZeroDec(), creatorStake)
	require.NoError(t, err)
	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)
	_, err = msgsrv.SetOnline(goCtx, types.NewMsgSetOnline(vals[0]))
	require.NoError(t, err)

	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	// 3. undelegate 100 'negative' 9 times
	for i := 0; i < 9; i++ {
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		goCtx = sdk.WrapSDKContext(ctx)
		msg := types.NewMsgUndelegate(accs[0], vals[0], sdk.NewCoin("negative", helpers.EtherToWei(sdkmath.NewInt(100))))
		_, err = msgsrv.Undelegate(goCtx, msg)
		require.NoError(t, err)
		keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
		val, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.True(t, val.Stake >= 0, "step %d", i)
	}
}
