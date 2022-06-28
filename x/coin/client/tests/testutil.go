package tests

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin/client/cli"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
)

func MsgCreateCoinExec(clientCtx client.Context, title, symbol, crr, initReserve, initVolume, limitVolume, identity, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{title, symbol, crr, initReserve, initVolume, limitVolume, identity, fmt.Sprintf("--%s=%s", flags.FlagFrom, from)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, cli.NewCreateCoinCmd(), args)
}

func MsgSendExec(clientCtx client.Context, from, to, amount fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{to.String(), amount.String(), fmt.Sprintf("--%s=%s", flags.FlagFrom, from)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, cli.NewSendCoinCmd(), args)
}

func MsgIssueCheckExec(clientCtx client.Context, from, amount, nonce, dueBlock, pass string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{amount, nonce, dueBlock, pass, fmt.Sprintf("--%s=%s", flags.FlagFrom, from)}
	args = append(args, extraArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, cli.NewIssueCheckCmd(), args)
}

func MsgCheckRedeemExec(clientCtx client.Context, check, pass, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{check, pass, fmt.Sprintf("--%s=%s", flags.FlagFrom, from)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, cli.NewRedeemCheckCmd(), args)
}
