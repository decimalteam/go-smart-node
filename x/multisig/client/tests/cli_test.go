package tests

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/x/multisig/client/cli"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCliCreateWallet(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	cmd := cli.NewCreateWalletCmd()
	ctx := setUpCmd(t, cmd, clientCtx, accs[0].GetAddress().String())

	// owners weights threshold
	cmd.SetArgs([]string{accs[0].GetAddress().String() + "," + accs[1].GetAddress().String(), "1,2", "3"})
	err := cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgCreateWallet)
	require.True(t, ok)
	require.Equal(t, []string{accs[0].GetAddress().String(), accs[1].GetAddress().String()}, msg.Owners)
	require.Equal(t, []uint64{1, 2}, msg.Weights)
	require.Equal(t, uint64(3), msg.Threshold)
}

func TestCliCreateTransaction(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	cmd := cli.NewCreateTransactionCmd()
	ctx := setUpCmd(t, cmd, clientCtx, accs[0].GetAddress().String())

	coins := sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(10)), sdk.NewCoin("btc", sdk.NewInt(100)))
	// wallet receiver coins
	cmd.SetArgs([]string{accs[0].GetAddress().String(), accs[1].GetAddress().String(), "10del,100btc"})
	err := cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgCreateTransaction)
	require.True(t, ok)
	require.Equal(t, accs[0].GetAddress().String(), msg.Wallet)
	require.Equal(t, accs[1].GetAddress().String(), msg.Receiver)
	require.True(t, msg.Coins.IsEqual(coins))
}

func TestCliSignTransaction(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	cmd := cli.NewSignTransactionCmd()
	ctx := setUpCmd(t, cmd, clientCtx, accs[0].GetAddress().String())

	// tx_id: dxmstx+1+....
	txID, err := sdk.Bech32ifyAddressBytes(types.MultisigTransactionIDPrefix, []byte{1, 2, 3})
	require.NoError(t, err)
	cmd.SetArgs([]string{txID})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgSignTransaction)
	require.True(t, ok)
	require.Equal(t, txID, msg.TxID)

}
