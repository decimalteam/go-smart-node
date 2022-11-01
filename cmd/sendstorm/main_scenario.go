package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	stormScenario "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/scenario"
)

func cmdScenario() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "scenario",
		Short: "Run actions",
		Long:  "Blockchain test scenarios",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return nil
		},
		SuggestionsMinimumDistance: 1,
	}

	cmd.AddCommand(delegationsReadTime())
	cmd.AddCommand(redeemCheckScenario())

	return cmd
}

/*
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
*/
func redeemCheckScenario() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "redeem_check",
		Short: "Redeem check scenarion",
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
			sc := stormScenario.NewRedeemChecksScenario(reactor.api, reactor.accounts[0])
			time.Sleep(time.Second * 10)
			step := 0
			for {
				step++
				fmt.Printf("--------step: %d--------\n", step)
				sc.MakeCheck()
				time.Sleep(time.Second * 5)
			}
		},
	}

	return cmd
}

func delegationsReadTime() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "drs",
		Short: "Delegations Read time scenario",
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

			valCount, err := cmd.Flags().GetInt(stormScenario.ValCountFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
			coinCount, err := cmd.Flags().GetInt(stormScenario.CoinsCountFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
			nftsCount, err := cmd.Flags().GetInt(stormScenario.NftsCountFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
			delegationsCount, err := cmd.Flags().GetInt(stormScenario.DelegationsFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
			createValidators, err := cmd.Flags().GetBool(stormScenario.CreateValidators)
			if err != nil {
				fmt.Println(err)
				return
			}
			createCoins, err := cmd.Flags().GetBool(stormScenario.CreateCoins)
			if err != nil {
				fmt.Println(err)
				return
			}
			createNfts, err := cmd.Flags().GetBool(stormScenario.CreateNfts)
			if err != nil {
				fmt.Println(err)
				return
			}

			drs := stormScenario.NewDelegationsReadScenario(reactor.api, reactor.accounts)
			time.Sleep(time.Second * 10)
			step := 0
			for {
				step++
				fmt.Printf("--------step: %d--------\n", step)
				drs.Start(valCount, coinCount, nftsCount, delegationsCount, createValidators, createCoins, createNfts)
				time.Sleep(time.Second * 5)
			}
		},
	}

	cmd.Flags().Int(stormScenario.ValCountFlag, 50, "validators count in test")
	cmd.Flags().Int(stormScenario.CoinsCountFlag, 50, "custom coins in test")
	cmd.Flags().Int(stormScenario.NftsCountFlag, 50, "nfts in test")
	cmd.Flags().Int(stormScenario.DelegationsFlag, 150_000, "delegations(nft and coin) in test")
	cmd.Flags().Bool(stormScenario.CreateValidators, false, "create new validators?")
	cmd.Flags().Bool(stormScenario.CreateCoins, false, "create new coins?")
	cmd.Flags().Bool(stormScenario.CreateNfts, false, "create new nfts?")

	return cmd
}
