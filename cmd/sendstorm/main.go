package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	helpers "bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"
)

const (
	mnemonicsFlag  = "mnemonics_file"
	nodeFlag       = "node"
	tendermintPort = "tport"
	restPort       = "rport"
	turnOnDebug    = "debug"
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
	rootCmd.PersistentFlags().String(nodeFlag, "http://localhost", "hostname of decimal node as http://... without port")
	rootCmd.PersistentFlags().Int(tendermintPort, 26657, "tendermint RPC port of decimal node")
	rootCmd.PersistentFlags().Int(restPort, 1317, "REST port of decimal node")
	rootCmd.PersistentFlags().Bool(turnOnDebug, false, "write api requests/responses to sendstorm.log")

	rootCmd.AddCommand(
		cmdGenerate(),
		cmdFaucet(),
		cmdRun(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cmdGenerate() *cobra.Command {
	var mnemonicsCount int

	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate mnemonics and save them to file",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Flags().Parse(args)
			if err != nil {
				fmt.Println(err)
				return
			}
			mnemonicsFile, err := cmd.Flags().GetString(mnemonicsFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
			f, err := os.Create(mnemonicsFile)
			if err != nil {
				fmt.Printf("create mnemonic file error: %s\n", err.Error())
				return
			}
			defer f.Close()
			fmt.Printf("create %d mnemonics\n", mnemonicsCount)
			for i := 0; i < mnemonicsCount; i++ {
				mn, err := wallet.NewMnemonic("")
				if err != nil {
					fmt.Printf("create mnemonic error: %s\n", err.Error())
					continue
				}
				_, err = fmt.Fprintf(f, "%s\n", mn.Words())
				if err != nil {
					fmt.Printf("write to mnemonic file error: %s\n", err.Error())
				}
			}
		},
	}

	cmd.PersistentFlags().IntVar(&mnemonicsCount, "count", 10, "count of mnemonics")

	return cmd
}

func cmdFaucet() *cobra.Command {
	var onlyEmpty bool
	var amountToSend int64

	var cmd = &cobra.Command{
		Use:   "faucet",
		Short: "Send some base coins to accounts from mnemonics",
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
			err = reactor.initFaucet(cmd.Flags())
			if err != nil {
				fmt.Println(err)
				return
			}
			// do action
			for i, acc := range reactor.accounts {
				err := acc.UpdateBalance()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("account: (%d) %s, balance: %s\n", i, acc.Address(), acc.BalanceForCoin(reactor.api.BaseCoin()))
				if onlyEmpty && acc.BalanceForCoin(reactor.api.BaseCoin()).GT(sdk.ZeroInt()) {
					continue
				}
				msg := dscTx.NewMsgSendCoin(
					reactor.faucetAccount.SdkAddress(),
					sdk.NewCoin(reactor.api.BaseCoin(), helpers.EtherToWei(sdk.NewInt(amountToSend))),
					acc.Account().SdkAddress(),
				)
				tx, err := dscTx.BuildTransaction(reactor.faucetAccount, []sdk.Msg{msg}, "", reactor.api.BaseCoin())
				if err != nil {
					fmt.Println(err)
					continue
				}
				err = tx.SignTransaction(reactor.faucetAccount)
				if err != nil {
					fmt.Println(err)
					continue
				}
				data, err := tx.BytesToSend()
				if err != nil {
					fmt.Println(err)
					continue
				}
				res, err := reactor.api.BroadcastTxSync(data)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if res.Code != 0 {
					fmt.Println(res)
					time.Sleep(time.Second * 6) // wait ~block
					// reset faucet
					err = reactor.initFaucet(cmd.Flags())
					if err != nil {
						fmt.Println(err)
						return
					}
					continue
				}
				reactor.faucetAccount.IncrementSequence()
			}
		},
	}

	cmd.PersistentFlags().String("faucet_mnemonic", "", "faucet mnemonic")
	cmd.PersistentFlags().BoolVar(&onlyEmpty, "only_empty", true, "send coins to account with zero balance")
	cmd.PersistentFlags().Int64Var(&amountToSend, "amount", 10000, "amount of base coins to send")

	return cmd
}

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
			for _, acc := range reactor.accounts {
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
			n := 0
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
				bytesToSend, err := action.GenerateTx(acc)
				if err != nil {
					fmt.Println(err)
					return
				}
				res, err := reactor.api.BroadcastTxSync(bytesToSend)
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
				//fmt.Printf("%T:: tx hash: %s\n", action, res.Hash)
				acc.IncrementSequence()
				go acc.UpdateBalance()
				n++
				if n >= 100 {
					go reactor.updateGeneratorsInfo()
					n = 0
				}
			}
		},
	}

	cmd.PersistentFlags().StringSlice("actions", []string{},
		"actions list in format: action1=weight1,action2=weight2,... weight must be integer")
	cmd.PersistentFlags().Int64("tps", 1, "transactions per second")

	return cmd
}

// helpers
func loadMnemonics(fname string) ([]string, error) {
	var result []string
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		mn := fileScanner.Text()
		result = append(result, mn)
	}
	return result, nil
}
