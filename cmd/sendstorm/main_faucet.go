package main

import (
	"fmt"
	"time"

	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	helpers "bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

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

// use external devnet/testnet faucet
func cmdFaucetExt() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "faucet-ext",
		Short: "Send some base coins to accounts",
		Run: func(cmd *cobra.Command, args []string) {
			//
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
			mnemonics, err := loadMnemonics(mnemonicsFile)
			if err != nil {
				fmt.Println(err)
				return
			}

			// do action
			client := resty.New().SetHostURL("https://devnet-dec2.console.decimalchain.com/api/faucet")
			for i, mnemonic := range mnemonics {
				acc, err := dscWallet.NewAccountFromMnemonicWords(mnemonic, "")
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("account (%d) %s query\n", i, acc.Address())
				resp, err := client.R().SetBody(map[string]string{"address": acc.Address()}).Post("https://devnet-dec2.console.decimalchain.com/api/faucet")
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("account (%d) %s result: %s\n", i, acc.Address(), resp.Status())
			}
		},
	}

	return cmd
}
