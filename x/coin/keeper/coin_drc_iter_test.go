package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdconfig "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdkmath "cosmossdk.io/math"
)

func TestIterateCoinDRC(t *testing.T) {
	_, dsc, ctx := testkeeper.GetTestAppWithCoinKeeper(t)

	// Seed some base balance so the keeper is set up (mirrors bootstrapKeeperTest usage)
	_ = sdk.Coins{{
		Denom:  cmdconfig.BaseDenom,
		Amount: helpers.EtherToWei(sdkmath.NewInt(1000000000000)),
	}}

	// Seed two CoinDRC records
	err := dsc.CoinKeeper.UpdateCoinDRC(ctx, "del", "0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	require.NoError(t, err)

	err = dsc.CoinKeeper.UpdateCoinDRC(ctx, "abc", "0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	require.NoError(t, err)

	// Collect all records via IterateCoinDRC
	got := map[string]string{}
	dsc.CoinKeeper.IterateCoinDRC(ctx, func(denom, drc20 string) bool {
		got[denom] = drc20
		return false
	})

	// Records are stored lowercased (setCoinDRC lowercases the address)
	require.Equal(t, "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", got["del"])
	require.Equal(t, "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", got["abc"])

	// Test early-stop: first call returns true, so only one record is visited
	visited := 0
	dsc.CoinKeeper.IterateCoinDRC(ctx, func(denom, drc20 string) bool {
		visited++
		return true // stop after first
	})
	require.Equal(t, 1, visited)
}
