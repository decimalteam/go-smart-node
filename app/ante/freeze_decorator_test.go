package ante

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

type mockFreezeKeeper struct{ frozen bool }

func (m mockFreezeKeeper) IsFrozen(_ sdk.Context) bool { return m.frozen }

type mockTx struct{ msgs []sdk.Msg }

func (m mockTx) GetMsgs() []sdk.Msg   { return m.msgs }
func (m mockTx) ValidateBasic() error { return nil }

func noopNext(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
	return ctx, nil
}

func TestFreezeDecorator(t *testing.T) {
	ctx := sdk.Context{}

	// Not frozen: a normal tx passes through.
	dec := NewFreezeDecorator(mockFreezeKeeper{frozen: false})
	_, err := dec.AnteHandle(ctx, mockTx{msgs: []sdk.Msg{&validatortypes.MsgDelegate{}}}, false, noopNext)
	require.NoError(t, err)

	// Frozen: a normal cosmos tx is rejected.
	dec = NewFreezeDecorator(mockFreezeKeeper{frozen: true})
	_, err = dec.AnteHandle(ctx, mockTx{msgs: []sdk.Msg{&validatortypes.MsgDelegate{}}}, false, noopNext)
	require.ErrorIs(t, err, ChainIsFrozen)

	// Frozen: an EVM transaction is also rejected (FreezeDecorator runs in the eth ante chain too).
	_, err = dec.AnteHandle(ctx, mockTx{msgs: []sdk.Msg{&evmtypes.MsgEthereumTx{}}}, false, noopNext)
	require.ErrorIs(t, err, ChainIsFrozen)

	// Frozen: emergency admin messages still pass so the admin can react.
	for _, msg := range []sdk.Msg{
		&validatortypes.MsgHaltChain{},
		&validatortypes.MsgResumeChain{},
		&validatortypes.MsgFreezeChain{},
		&validatortypes.MsgUnfreezeChain{},
	} {
		_, err = dec.AnteHandle(ctx, mockTx{msgs: []sdk.Msg{msg}}, false, noopNext)
		require.NoErrorf(t, err, "emergency msg %T must pass while frozen", msg)
	}
}

func TestIsEmergencyMsg(t *testing.T) {
	require.True(t, isEmergencyMsg(&validatortypes.MsgHaltChain{}))
	require.True(t, isEmergencyMsg(&validatortypes.MsgResumeChain{}))
	require.True(t, isEmergencyMsg(&validatortypes.MsgFreezeChain{}))
	require.True(t, isEmergencyMsg(&validatortypes.MsgUnfreezeChain{}))
	require.False(t, isEmergencyMsg(&validatortypes.MsgDelegate{}))
	require.False(t, isEmergencyMsg(&evmtypes.MsgEthereumTx{}))
}
