package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
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
			addresses, err := reactor.api.AllAccounts()
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
			for _, adr := range addresses {
				coins, err := reactor.api.AddressBalance(adr)
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

func cmdValidators() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "validators",
		Short: "Show validators info",
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
			// validators info
			vals, err := reactor.api.Validators()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, val := range vals {
				dels, err := reactor.api.ValidatorDelegations(val.OperatorAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("moniker: %s, status: %d, online: %v, jailed: %v, stake: %d, rewards: %s, delegation: %d\n",
					val.Description.Moniker, val.Status, val.Online, val.Jailed, val.Stake, val.Rewards, len(dels))
			}
		},
	}

	return cmd
}

func cmdVerifyPools() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "verify-pools",
		Short: "Verify (un/re)delegations and validator pools",
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
			// validators info
			vals, err := reactor.api.Validators()
			if err != nil {
				fmt.Println(err)
				return
			}
			expectBondedPool := sdk.NewCoins()
			expectNotBondedPool := sdk.NewCoins()
			for _, val := range vals {
				// delegations
				dels, err := reactor.api.ValidatorDelegations(val.OperatorAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, del := range dels {
					if del.Stake.Type == dscApi.StakeType_Coin {
						switch val.Status {
						case dscApi.BondStatus_Bonded:
							expectBondedPool = expectBondedPool.Add(del.Stake.Stake)
						case dscApi.BondStatus_Unbonded, dscApi.BondStatus_Unbonding:
							expectNotBondedPool = expectNotBondedPool.Add(del.Stake.Stake)
						}
					}
				}
				// redelegations
				reds, err := reactor.api.ValidatorRedelegations(val.OperatorAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, red := range reds {
					for _, ent := range red.Entries {
						if ent.Stake.Type == dscApi.StakeType_Coin {
							expectNotBondedPool = expectNotBondedPool.Add(ent.Stake.Stake)
						}
					}
				}
				// undelegations
				ubds, err := reactor.api.ValidatorUndelegations(val.OperatorAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, ubd := range ubds {
					for _, ent := range ubd.Entries {
						if ent.Stake.Type == dscApi.StakeType_Coin {
							expectNotBondedPool = expectNotBondedPool.Add(ent.Stake.Stake)
						}
					}
				}
			}
			// TODO: check pool
		},
	}

	return cmd
}
