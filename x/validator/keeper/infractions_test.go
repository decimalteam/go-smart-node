package keeper_test

import (
	"testing"
	"time"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestStartHeightOperations(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	consAdr := sdk.GetConsAddress(PKs[0])

	require.Equal(t, int64(-1), dsc.ValidatorKeeper.GetStartHeight(ctx, consAdr))

	dsc.ValidatorKeeper.SetStartHeight(ctx, consAdr, 10)
	require.Equal(t, int64(10), dsc.ValidatorKeeper.GetStartHeight(ctx, consAdr))

	dsc.ValidatorKeeper.DeleteStartHeight(ctx, consAdr)
	require.Equal(t, int64(-1), dsc.ValidatorKeeper.GetStartHeight(ctx, consAdr))
}

func TestSigningInfo(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	consAdr := sdk.GetConsAddress(PKs[0])
	var window int64 = 24

	// before start
	dsc.ValidatorKeeper.AddMissedBlock(ctx, consAdr, 1)
	dsc.ValidatorKeeper.AddMissedBlock(ctx, consAdr, 2)
	// start
	dsc.ValidatorKeeper.SetStartHeight(ctx, consAdr, 5)
	// after start
	dsc.ValidatorKeeper.AddMissedBlock(ctx, consAdr, 5)
	dsc.ValidatorKeeper.AddMissedBlock(ctx, consAdr, 6)
	dsc.ValidatorKeeper.AddMissedBlock(ctx, consAdr, 8)
	// after 24 window
	dsc.ValidatorKeeper.AddMissedBlock(ctx, consAdr, 5+window)

	vsi := dsc.ValidatorKeeper.GetValidatorSigningInfo(ctx, consAdr, 2, 2+window)
	require.Equal(t, int64(5), vsi.StartHeight)
	require.Equal(t, int64(3), vsi.MissedBlocksCounter)

	vsi = dsc.ValidatorKeeper.GetValidatorSigningInfo(ctx, consAdr, 100, 100+window)
	require.Equal(t, int64(5), vsi.StartHeight)
	require.Equal(t, int64(0), vsi.MissedBlocksCounter)
}

func TestSigningInfoLimit(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	consAdr := sdk.GetConsAddress(PKs[0])
	var startH int64 = 5
	var window int64 = 24
	dsc.ValidatorKeeper.SetStartHeight(ctx, consAdr, startH)
	for i := int64(0); i < window+10; i++ {
		dsc.ValidatorKeeper.AddMissedBlock(ctx, consAdr, i)
	}
	height := window + 5
	vsi := dsc.ValidatorKeeper.GetValidatorSigningInfo(ctx, consAdr, height-window, height-1)
	require.Equal(t, window, vsi.MissedBlocksCounter)
}

func TestMissingSignature(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)
	//disable grace period
	cmdcfg.UpdatesInfo.LastBlock = 1000000

	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))
	consAdr := sdk.GetConsAddress(PKs[1])

	params := dsc.ValidatorKeeper.GetParams(ctx)

	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	msgCreate, err := types.NewMsgCreateValidator(vals[1], accs[1], PKs[1], types.Description{Moniker: "monik"},
		sdk.ZeroDec(), creatorStake)
	require.NoError(t, err)

	//
	goCtx := sdk.WrapSDKContext(ctx)
	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)

	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
	goCtx = sdk.WrapSDKContext(ctx)
	msgOnline := types.NewMsgSetOnline(vals[1])
	_, err = msgsrv.SetOnline(goCtx, msgOnline)
	require.NoError(t, err)

	for i := int64(0); i < params.SignedBlocksWindow; i++ {
		keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
		dsc.ValidatorKeeper.HandleValidatorSignature(ctx, consAdr.Bytes(), 0, false, params)
		val, found := dsc.ValidatorKeeper.GetValidatorByConsAddrDecimal(ctx, consAdr)
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, val.Status, "fail at step=%d", i)
	}
	// check min height not passed
	val, found := dsc.ValidatorKeeper.GetValidatorByConsAddrDecimal(ctx, consAdr)
	require.True(t, found)
	require.True(t, val.Online)
	require.False(t, val.Jailed)
	require.Equal(t, types.BondStatus_Bonded, val.Status)

	// min height passed, must be jailed and slashed
	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
	dsc.ValidatorKeeper.HandleValidatorSignature(ctx, consAdr.Bytes(), 0, false, params)

	val, found = dsc.ValidatorKeeper.GetValidatorByConsAddrDecimal(ctx, consAdr)
	require.True(t, found)
	require.False(t, val.Online)
	require.True(t, val.Jailed)
	// 1% slash
	del, found := dsc.ValidatorKeeper.GetDelegation(ctx, accs[1], vals[1], cmdcfg.BaseDenom)
	require.True(t, found)
	require.True(t, del.Stake.Stake.Equal(
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(99))),
	))
}

func TestDoubleSignature(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	//disable grace period
	cmdcfg.UpdatesInfo.LastBlock = 1000000

	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))
	consAdr := sdk.GetConsAddress(PKs[0])

	params := dsc.ValidatorKeeper.GetParams(ctx)

	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	msgCreate, err := types.NewMsgCreateValidator(vals[0], accs[0], PKs[0], types.Description{Moniker: "monik"},
		sdk.ZeroDec(), creatorStake)
	require.NoError(t, err)

	//
	goCtx := sdk.WrapSDKContext(ctx)
	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)

	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
	goCtx = sdk.WrapSDKContext(ctx)
	msgOnline := types.NewMsgSetOnline(vals[0])
	_, err = msgsrv.SetOnline(goCtx, msgOnline)
	require.NoError(t, err)

	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
	val, found := dsc.ValidatorKeeper.GetValidatorByConsAddrDecimal(ctx, consAdr)
	require.Equal(t, types.BondStatus_Bonded, val.Status)

	dsc.ValidatorKeeper.HandleDoubleSign(ctx, consAdr.Bytes(), ctx.BlockHeight(), time.Now(), 0, params)

	// min height passed, must be jailed and slashed
	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))

	val, found = dsc.ValidatorKeeper.GetValidatorByConsAddrDecimal(ctx, consAdr)
	require.True(t, found)
	require.False(t, val.Online)
	require.True(t, val.Jailed)
	// 5% slash
	del, found := dsc.ValidatorKeeper.GetDelegation(ctx, accs[0], vals[0], cmdcfg.BaseDenom)
	require.True(t, found)
	require.True(t, del.Stake.Stake.Equal(
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(95))),
	))
}
