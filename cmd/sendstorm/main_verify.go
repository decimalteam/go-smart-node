package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func cmdVerify() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify custom coins volume",
		Run: func(cmd *cobra.Command, args []string) {
			//
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
			// coins info from coin keeper
			coinsInfo, err := reactor.coins()
			if err != nil {
				fmt.Println(err)
				return
			}
			balances := sdk.NewCoins()
			for _, acc := range reactor.accounts {
				coins, err := reactor.api.AddressBalance(acc.Address())
				if err != nil {
					fmt.Println(err)
					return
				}
				balances = balances.Add(coins...)
			}
			for _, coinInfo := range coinsInfo {
				if coinInfo.Denom == reactor.api.BaseCoin() {
					continue
				}
				diff := "none"
				bal := balances.AmountOf(coinInfo.Denom)
				if !bal.Equal(coinInfo.Volume) {
					diff = fmt.Sprintf("volume=%s, balances=%s",
						coinInfo.Volume.String(), bal.String())
				}
				fmt.Printf("coin: (symbol: %s, volume: %s, reserve: %s), difference: %s\n",
					coinInfo.Denom, coinInfo.Volume, coinInfo.Reserve, diff)
			}
		},
	}

	return cmd
}
