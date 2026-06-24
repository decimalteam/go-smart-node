package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func TestShouldHardHalt(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	vk := &dsc.ValidatorKeeper

	// No halt info recorded => never halts.
	_, halt := vk.ShouldHardHalt(ctx.WithBlockHeight(10))
	require.False(t, halt)

	// Active halt: triggers at and after the target height, not before.
	vk.SetHaltInfo(ctx, types.HaltInfo{Active: true, Height: 5, Reason: "incident"})
	_, halt = vk.ShouldHardHalt(ctx.WithBlockHeight(4))
	require.False(t, halt, "must not halt before the target height")
	info, halt := vk.ShouldHardHalt(ctx.WithBlockHeight(5))
	require.True(t, halt, "must halt at the target height")
	require.Equal(t, "incident", info.Reason)
	_, halt = vk.ShouldHardHalt(ctx.WithBlockHeight(10))
	require.True(t, halt, "must stay halted after the target height")

	// Inactive halt info => never halts.
	vk.SetHaltInfo(ctx, types.HaltInfo{Active: false, Height: 5})
	_, halt = vk.ShouldHardHalt(ctx.WithBlockHeight(10))
	require.False(t, halt)

	// --unsafe-skip-halt overrides an active halt.
	vk.SetHaltInfo(ctx, types.HaltInfo{Active: true, Height: 5})
	vk.SetSkipHalt(true)
	_, halt = vk.ShouldHardHalt(ctx.WithBlockHeight(10))
	require.False(t, halt, "skipHalt must suppress the halt")
	vk.SetSkipHalt(false)
	_, halt = vk.ShouldHardHalt(ctx.WithBlockHeight(10))
	require.True(t, halt)
}

func TestBeginBlockerHaltPanics(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	vk := &dsc.ValidatorKeeper

	vk.SetHaltInfo(ctx, types.HaltInfo{Active: true, Height: 5, Reason: "boom"})

	// At/after the halt height the BeginBlocker must panic to stop consensus.
	require.Panics(t, func() {
		keeper.BeginBlocker(ctx.WithBlockHeight(5), dsc.ValidatorKeeper, abci.RequestBeginBlock{})
	})

	// With --unsafe-skip-halt the halt check must not panic.
	vk.SetSkipHalt(true)
	require.NotPanics(t, func() {
		_, halt := dsc.ValidatorKeeper.ShouldHardHalt(ctx.WithBlockHeight(5))
		require.False(t, halt)
	})
}

func TestHaltChainAuthorization(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgServer := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	admin := types.GetHaltAdmin(ctx.ChainID())
	require.NotEmpty(t, admin, "test chain must have a configured halt-admin")

	addrs, _ := generateAddresses(dsc, ctx, 1, sdk.NewCoins())
	notAdmin := addrs[0].String()
	require.NotEqual(t, admin, notAdmin)

	// A non-admin sender is rejected and no halt is recorded.
	_, err := msgServer.HaltChain(sdk.WrapSDKContext(ctx), &types.MsgHaltChain{
		Sender: notAdmin, Height: 0, Reason: "nope",
	})
	require.Error(t, err)
	require.False(t, dsc.ValidatorKeeper.IsHaltScheduled(ctx))

	// The admin can schedule a halt; height 0 resolves to the next block.
	ctxH := ctx.WithBlockHeight(100)
	_, err = msgServer.HaltChain(sdk.WrapSDKContext(ctxH), &types.MsgHaltChain{
		Sender: admin, Height: 0, Reason: "incident",
	})
	require.NoError(t, err)
	info, found := dsc.ValidatorKeeper.GetHaltInfo(ctxH)
	require.True(t, found)
	require.True(t, info.Active)
	require.Equal(t, int64(101), info.Height)
	require.Equal(t, "incident", info.Reason)

	// A past/equal explicit height is rejected.
	_, err = msgServer.HaltChain(sdk.WrapSDKContext(ctxH), &types.MsgHaltChain{
		Sender: admin, Height: 100, Reason: "too late",
	})
	require.Error(t, err)

	// The admin can clear the halt.
	_, err = msgServer.ResumeChain(sdk.WrapSDKContext(ctxH), &types.MsgResumeChain{Sender: admin})
	require.NoError(t, err)
	require.False(t, dsc.ValidatorKeeper.IsHaltScheduled(ctxH))
}

func TestFreezeUnfreeze(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgServer := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	admin := types.GetHaltAdmin(ctx.ChainID())
	require.NotEmpty(t, admin)

	require.False(t, dsc.ValidatorKeeper.IsFrozen(ctx))

	// Non-admin cannot freeze.
	addrs, _ := generateAddresses(dsc, ctx, 1, sdk.NewCoins())
	_, err := msgServer.FreezeChain(sdk.WrapSDKContext(ctx), &types.MsgFreezeChain{
		Sender: addrs[0].String(), Reason: "nope",
	})
	require.Error(t, err)
	require.False(t, dsc.ValidatorKeeper.IsFrozen(ctx))

	// Admin can freeze and unfreeze.
	_, err = msgServer.FreezeChain(sdk.WrapSDKContext(ctx), &types.MsgFreezeChain{
		Sender: admin, Reason: "anomaly",
	})
	require.NoError(t, err)
	require.True(t, dsc.ValidatorKeeper.IsFrozen(ctx))

	_, err = msgServer.UnfreezeChain(sdk.WrapSDKContext(ctx), &types.MsgUnfreezeChain{Sender: admin})
	require.NoError(t, err)
	require.False(t, dsc.ValidatorKeeper.IsFrozen(ctx))
}
