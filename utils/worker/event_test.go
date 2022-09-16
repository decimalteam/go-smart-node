package worker

import (
	"testing"

	utilsEvents "bitbucket.org/decimalteam/go-smart-node/utils/events"
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func TestCreateCoin(t *testing.T) {
	ea := NewEventAccumulator()
	tev := &coinTypes.EventCreateCoin{
		Sender:               "a",
		Symbol:               "test",
		Title:                "title",
		Crr:                  10,
		InitialVolume:        sdkmath.NewInt(100).String(),
		InitialReserve:       sdkmath.NewInt(200).String(),
		LimitVolume:          sdkmath.NewInt(1000).String(),
		Identity:             "ident",
		CommissionCreateCoin: sdk.NewCoin("del", sdkmath.NewInt(9000)).String(),
	}

	ev, err := utilsEvents.TypedEventToEvent(tev)
	require.NoError(t, err)
	err = processEventCreateCoin(ea, abci.Event(ev), "ABCD", 7)
	require.NoError(t, err)
	require.Equal(t, 1, len(ea.CoinsCreates))
	require.True(t, ea.BalancesChanges["a"]["del"].Equal(sdkmath.NewInt(9200).Neg())) //comission+reserve
	require.True(t, ea.BalancesChanges["a"]["test"].Equal(sdkmath.NewInt(100)))
	cc := ea.CoinsCreates[0]
	require.Equal(t, "a", cc.Creator)
	require.Equal(t, "test", cc.Symbol)
	require.Equal(t, "title", cc.Title)
	require.Equal(t, "ident", cc.Avatar)
	require.Equal(t, uint64(10), cc.CRR)
	require.Equal(t, "ABCD", cc.TxHash)
	require.Equal(t, int64(7), cc.BlockID)
	require.True(t, cc.Volume.Equal(sdkmath.NewInt(100)))
	require.True(t, cc.Reserve.Equal(sdkmath.NewInt(200)))
	require.True(t, cc.LimitVolume.Equal(sdkmath.NewInt(1000)))
}
