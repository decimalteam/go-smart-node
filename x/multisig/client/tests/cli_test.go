package tests

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/x/multisig/client/cli"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

func TestCliCreateWallet(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)
	adr1, err := accs[1].GetAddress()
	require.NoError(t, err)

	cmd := cli.NewCreateWalletCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())

	// owners weights threshold
	cmd.SetArgs([]string{adr0.String() + "," + adr1.String(), "1,2", "3"})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgCreateWallet)
	require.True(t, ok)
	require.Equal(t, []string{adr0.String(), adr1.String()}, msg.Owners)
	require.Equal(t, []uint32{1, 2}, msg.Weights)
	require.Equal(t, uint32(3), msg.Threshold)
}

/*
func TestCliCreateTransaction(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)
	adr1, err := accs[1].GetAddress()
	require.NoError(t, err)

	cmd := cli.NewCreateTransactionCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())

	coins := sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, sdk.NewInt(10)), sdk.NewCoin("btc", sdk.NewInt(100)))
	// wallet receiver coins
	cmd.SetArgs([]string{adr0.String(), adr1.String(), "10" + cmdcfg.BaseDenom + ",100btc"})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgCreateTransaction)
	require.True(t, ok)
	require.Equal(t, adr0.String(), msg.Wallet)
	require.Equal(t, adr1.String(), msg.Receiver)
	require.True(t, msg.Coins.IsEqual(coins))
}
*/

func TestCliSignTransaction(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)

	cmd := cli.NewSignTransactionCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())

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
	require.Equal(t, txID, msg.ID)

}
