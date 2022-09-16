package worker

import (
	"testing"

	utilsEvents "bitbucket.org/decimalteam/go-smart-node/utils/events"
	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestMultisigCreateWallet(t *testing.T) {
	ea := NewEventAccumulator()
	tev := &multisigTypes.EventCreateWallet{
		Sender:    "aaa",
		Wallet:    "bbb",
		Owners:    []string{"a", "b", "c"},
		Weights:   []uint64{1, 2, 3},
		Threshold: 10,
	}
	ev, err := utilsEvents.TypedEventToEvent(tev)
	require.NoError(t, err)
	err = processEventCreateWallet(ea, abci.Event(ev), "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(ea.MultisigCreateWallets))
	require.Equal(t, 3, len(ea.MultisigCreateWallets[0].Owners))
}

func TestBalanceChanges(t *testing.T) {
	ea := NewEventAccumulator()
	ea.addBalanceChange("a", "coin1", sdkmath.NewInt(1))
	ea.addBalanceChange("b", "coin2", sdkmath.NewInt(2))
	ea.addBalanceChange("a", "coin1", sdkmath.NewInt(2))
	ea.addBalanceChange("a", "coin1", sdkmath.NewInt(1).Neg())
	ea.addBalanceChange("a", "coin2", sdkmath.NewInt(1).Neg())
	require.True(t, ea.BalancesChanges["a"]["coin1"].Equal(sdkmath.NewInt(2)))
	require.True(t, ea.BalancesChanges["a"]["coin2"].Equal(sdkmath.NewInt(-1)))
	require.True(t, ea.BalancesChanges["b"]["coin2"].Equal(sdkmath.NewInt(2)))
}
