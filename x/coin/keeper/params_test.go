package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.CoinKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
