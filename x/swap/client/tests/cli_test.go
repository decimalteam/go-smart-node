package tests

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/x/swap/client/cli"
	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCliSwapInitialize(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)

	cmd := cli.NewSwapInitializeCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())

	// [recipient] [amount] [token_symbol] [tx_number] [from_chain] [dest_chain]
	cmd.SetArgs([]string{"0x12345", "1000", "btc", "1234567", "1", "3"})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgInitializeSwap)
	require.True(t, ok)
	require.Equal(t, "0x12345", msg.Recipient)
	require.True(t, msg.Amount.Equal(sdk.NewInt(1000)))
	require.Equal(t, "btc", msg.TokenSymbol)
	require.Equal(t, "1234567", msg.TransactionNumber)
	require.Equal(t, uint32(1), msg.FromChain)
	require.Equal(t, uint32(3), msg.DestChain)
}

func TestCliSwapRedeem(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)
	adr1, err := accs[1].GetAddress()
	require.NoError(t, err)

	cmd := cli.NewSwapRedeemCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())

	// [from] [recipient] [amount] [token_symbol] [tx_number] [from_chain] [dest_chain] [v] [r] [s]
	// r = 61 73 64 62 6e 31 32 33 38 37 67 61
	// s = 34 66 73 37
	cmd.SetArgs([]string{"0x12345", adr1.String(), "1000", "btc", "1234567", "3", "1", "9",
		"d8c0c8ff4a9b168be168f480bae61ead0a7f2b973f983a038f867621451fa553",
		"641ba9f5749afbb425e83b69ecacb3a0c6e32e2431609d474d4300a7cce5eb41"})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgRedeemSwap)
	require.True(t, ok)
	require.Equal(t, "0x12345", msg.From)
	require.Equal(t, adr1.String(), msg.Recipient)
	require.True(t, msg.Amount.Equal(sdk.NewInt(1000)))
	require.Equal(t, "btc", msg.TokenSymbol)
	require.Equal(t, "1234567", msg.TransactionNumber)
	require.Equal(t, uint32(3), msg.FromChain)
	require.Equal(t, uint32(1), msg.DestChain)
	require.Equal(t, uint32(9), msg.V)
	require.Equal(t, "d8c0c8ff4a9b168be168f480bae61ead0a7f2b973f983a038f867621451fa553", msg.R)
	require.Equal(t, "641ba9f5749afbb425e83b69ecacb3a0c6e32e2431609d474d4300a7cce5eb41", msg.S)

}

/*
func TestCliChainActivate(t *testing.T) {
	clientCtx, _, result := setUpCliTest(t, 2)

	cmd := cli.NewChainActivateCmd()
	// [number] [name]; sender must be swap service address
	ctx := setUpCmd(t, cmd, clientCtx, "dx1jqx7chw0faswfmw78cdejzzery5akzmk5zc5x5")

	// tx_id: dxmstx+1+....
	cmd.SetArgs([]string{"1", "Decimal AAA"})
	err := cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgChainActivate)
	require.True(t, ok)
	require.Equal(t, uint32(1), msg.ChainNumber)
	require.Equal(t, "Decimal AAA", msg.ChainName)
}


func TestCliChainDeactivate(t *testing.T) {
	clientCtx, _, result := setUpCliTest(t, 2)

	cmd := cli.NewChainDeactivateCmd()
	// [number] [name]; sender must be swap service address
	ctx := setUpCmd(t, cmd, clientCtx, "dx1jqx7chw0faswfmw78cdejzzery5akzmk5zc5x5")

	// tx_id: dxmstx+1+....
	cmd.SetArgs([]string{"1"})
	err := cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgChainDeactivate)
	require.True(t, ok)
	require.Equal(t, uint32(1), msg.ChainNumber)
}
*/
