package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	stormScenario "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/scenario"
)

func cmdScenario() *cobra.Command {
	var subtokensCount int

	var cmd = &cobra.Command{
		Use:   "scenario",
		Short: "Run actions",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Flags().Parse(args)
			if err != nil {
				fmt.Println(err)
				return
			}
			reactor := stormReactor{}
			// init
			err = reactor.initApi(cmd.Flags())
			if err != nil {
				fmt.Println(err)
				return
			}
			err = reactor.initAccounts(cmd.Flags())
			if err != nil {
				fmt.Println(err)
				return
			}
			sc := stormScenario.NewNFTBlowScenario(reactor.api, reactor.accounts)
			sc.CreateNFTs(uint32(subtokensCount))
			time.Sleep(time.Second * 10)
			step := 0
			for {
				step++
				fmt.Printf("--------step: %d--------\n", step)
				sc.SendNFT()
				time.Sleep(time.Second * 5)
			}
		},
	}

	cmd.PersistentFlags().IntVar(&subtokensCount, "subtokens", 10, "count of subtokens")

	return cmd
}
