package app

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyContractPenalties(t *testing.T) {
	amt := big.NewInt(1000)
	out, applied := applyContractPenalties(amt, []uint16{10, 50}, big.NewInt(0))
	require.True(t, applied)
	require.Equal(t, int64(450), out.Int64()) // 1000 -100 = 900, 900 -450 = 450
	require.Equal(t, int64(1000), amt.Int64(), "input not mutated")

	out, applied = applyContractPenalties(amt, []uint16{10, 50}, big.NewInt(1))
	require.True(t, applied)
	require.Equal(t, int64(500), out.Int64(), "only penalties at or after the stake's index")

	out, applied = applyContractPenalties(amt, []uint16{10, 50}, big.NewInt(2))
	require.False(t, applied)
	require.Equal(t, int64(1000), out.Int64(), "all penalties already applied")

	out, applied = applyContractPenalties(big.NewInt(7), []uint16{33}, big.NewInt(0))
	require.True(t, applied)
	require.Equal(t, int64(5), out.Int64(), "integer division like solidity: 7 - 7*33/100 = 7 - 2")
}
