package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	mnemonicsFlag = "mnemonics_file"
	nodeFlag      = "node"
	turnOnDebug   = "debug"
	commitFlag    = "commit"
	customFee     = "customfee"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sendstorm",
		Short: "SendStorm is testing application for decimal blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	rootCmd.PersistentFlags().String(mnemonicsFlag, "mnemonics.cfg", "path to mnemonics file")
	rootCmd.PersistentFlags().String(nodeFlag, "localhost", "hostname or IP of decimal node without http:// or port")
	rootCmd.PersistentFlags().Bool(turnOnDebug, false, "write api requests/responses to sendstorm.log")
	rootCmd.PersistentFlags().Bool(commitFlag, false, "use broadcast_tx_commit (wait for block completion) for transaction sending (very slow)")
	rootCmd.PersistentFlags().Bool(customFee, false, "use custom coins for fee")

	rootCmd.AddCommand(
		cmdGenerate(),
		cmdFaucet(),
		cmdFaucetExt(),
		cmdRun(),
		cmdVerify(),
		cmdScenario(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
