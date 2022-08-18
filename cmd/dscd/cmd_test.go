package main_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/app"
	dscd "bitbucket.org/decimalteam/go-smart-node/cmd/dscd"
)

func TestInitCmd(t *testing.T) {
	rootCmd, _ := dscd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",     // Test the init cmd
		"dsc-test", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
		fmt.Sprintf("--%s=%s", flags.FlagChainID, "decimal_202020-1"),
	})

	err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome)
	require.NoError(t, err)
}

func TestAddKeyLedgerCmd(t *testing.T) {
	rootCmd, _ := dscd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"keys",
		"add",
		"royalkey",
		fmt.Sprintf("--%s", flags.FlagUseLedger),
	})

	err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome)
	require.Error(t, err)
}
