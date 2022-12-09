package keeper

import (
	"testing"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestInGracePeriod(t *testing.T) {
	ctxWithHeight := func(height int64) sdk.Context {
		ctx := sdk.Context{}
		return ctx.WithBlockHeight(height)
	}

	//test overlapping grace periods
	{
		updatesInfo := cmdcfg.NewUpdatesInfo("")
		p0start := int64(10000)
		p0end := p0start + cmdcfg.GracePeriod
		p1start := p0end - cmdcfg.GracePeriod/2
		p1end := p1start + cmdcfg.GracePeriod
		updatesInfo.AddExecutedPlan("0", p0start)
		updatesInfo.PushNewPlanHeight(p1start)
		//
		require.False(t, inGracePeriod(ctxWithHeight(p0start-1), updatesInfo))
		require.False(t, inGracePeriod(ctxWithHeight(p1end+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p0start+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p0end+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p1start-1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p1end-1), updatesInfo))
	}

	//test non-overlapping grace periods
	{
		updatesInfo := cmdcfg.NewUpdatesInfo("")
		p0start := int64(10000)
		p0end := p0start + cmdcfg.GracePeriod
		p1start := p0end + cmdcfg.GracePeriod/2
		p1end := p1start + cmdcfg.GracePeriod
		updatesInfo.AddExecutedPlan("0", p0start)
		updatesInfo.AddExecutedPlan("1", p1start)
		updatesInfo.PushNewPlanHeight(p1start)
		//
		require.False(t, inGracePeriod(ctxWithHeight(p0start-1), updatesInfo))
		require.False(t, inGracePeriod(ctxWithHeight(p0end+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p0start+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p0end-1), updatesInfo))
		require.False(t, inGracePeriod(ctxWithHeight(p1start-1), updatesInfo))
		require.False(t, inGracePeriod(ctxWithHeight(p1end+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p1start+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p1end-1), updatesInfo))
	}
}
