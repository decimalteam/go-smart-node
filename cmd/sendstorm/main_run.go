package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"bitbucket.org/decimalteam/go-smart-node/sdk/api"
	"github.com/spf13/cobra"
)

func cmdRun() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Run actions",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Flags().Parse(args)
			if err != nil {
				fmt.Println(err)
				return
			}
			reactor := stormReactor{}
			// init
			doCommitTx, err := cmd.Flags().GetBool(commitFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
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
			err = reactor.initActionReactor(cmd.Flags())
			if err != nil {
				fmt.Println(err)
				return
			}
			err = reactor.initLimiter(cmd.Flags())
			if err != nil {
				fmt.Println(err)
				return
			}
			for i, acc := range reactor.accounts {
				fmt.Printf("load balances(%d): %s\n", i, acc.Address())
				err = acc.UpdateNumberSequence()
				if err != nil {
					fmt.Println(err)
					return
				}
				err = acc.UpdateBalance()
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			reactor.updateGeneratorsInfo()
			// infinite loop
			n := int64(0)
			for {
				action := reactor.actionReactor.Generate()
				accs := action.ChooseAccounts(reactor.accounts)
				if len(accs) == 0 {
					continue
				}
				acc := accs[rand.Intn(len(accs))]
				if !reactor.limiter.CanMake() {
					continue
				}
				n++
				if n >= reactor.limiter.Limit()*5 {
					go reactor.updateGeneratorsInfo()
					runtime.GC()
					n = 0
				}
				bytesToSend, err := action.GenerateTx(acc, reactor.feeConfig)
				if err != nil {
					fmt.Println(err)
					return
				}
				var res *api.TxSyncResponse
				if doCommitTx {
					res, err = reactor.api.BroadcastTxCommit(bytesToSend)
				} else {
					res, err = reactor.api.BroadcastTxSync(bytesToSend)
				}
				if err != nil {
					fmt.Println(err)
					acc.MarkDirty()
					// TODO: Update() returns error
					go acc.UpdateNumberSequence()
					continue
				}
				if res.Code != 0 {
					fmt.Printf("account: %s, action: %s, result: %#v\n", acc.Address(), action, res)
					if res.Code == 111222 {
						os.Exit(1)
					}
					acc.MarkDirty()
					// TODO: Update() returns error
					go acc.UpdateNumberSequence()
					continue
				}
				if doCommitTx {
					time.Sleep(time.Second)
					txRes, err := reactor.api.Transaction(res.Hash)
					if err != nil {
						fmt.Printf("api.Transaction = %v\n", err)
					}
					fmt.Printf("TxHash= %s\n", res.Hash)
					fmt.Printf("TxResult = %v\n\n\n", txRes)
				}
				acc.IncrementSequence()
				go acc.UpdateBalance()
			}
		},
	}

	cmd.PersistentFlags().StringSlice("actions", []string{},
		"actions list in format: action1=weight1,action2=weight2,... weight must be integer")
	cmd.PersistentFlags().Int64("tps", 1, "transactions per second")

	return cmd
}
