package main

import (
	"bufio"
	sdkmath "cosmossdk.io/math"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	stormActions "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/actions"
	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/pflag"
)

type stormReactor struct {
	accounts      []*stormTypes.StormAccount
	faucetAccount *dscWallet.Account
	api           *dscApi.API
	actionReactor *stormActions.ActionReactor
	limiter       *stormActions.TPSLimiter
}

func (reactor *stormReactor) initApi(flags *pflag.FlagSet) error {
	nodeHost, err := flags.GetString(nodeFlag)
	if err != nil {
		return err
	}
	debug, err := flags.GetBool(turnOnDebug)
	if err != nil {
		return err
	}
	// TODO: make this dirty debug more nice
	if debug {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		logfile, err := os.Create("sendstorm.log")
		if err != nil {
			log.Fatalln("Cannot create log file")
		}
		//defer logfile.Close()
		log.SetOutput(logfile)
	}

	reactor.api, err = dscApi.NewAPI(dscApi.ConnectionOptions{
		EndpointHost: nodeHost,
		Timeout:      10,
	})
	if err != nil {
		return err
	}
	err = reactor.api.GetParameters()
	if err != nil {
		return err
	}
	return nil
}

func (reactor *stormReactor) initAccounts(flags *pflag.FlagSet) error {
	mnemonicsFile, err := flags.GetString(mnemonicsFlag)
	if err != nil {
		return err
	}
	mnemonics, err := loadMnemonics(mnemonicsFile)
	if err != nil {
		return err
	}
	reactor.accounts = make([]*stormTypes.StormAccount, 0)
	for _, mn := range mnemonics {
		acc, err := stormTypes.NewStormAccount(mn, reactor.api)
		if err != nil {
			return err
		}
		reactor.accounts = append(reactor.accounts, acc)
	}
	return nil
}

func (reactor *stormReactor) initFaucet(flags *pflag.FlagSet) error {
	faucetMnemonic, err := flags.GetString("faucet_mnemonic")
	if err != nil {
		return err
	}
	reactor.faucetAccount, err = dscWallet.NewAccountFromMnemonicWords(faucetMnemonic, "")
	if err != nil {
		return err
	}
	fmt.Printf("faucet address: %s\n", reactor.faucetAccount.Address())
	an, as, err := reactor.api.AccountNumberAndSequence(reactor.faucetAccount.Address())
	if err != nil {
		return err
	}
	reactor.faucetAccount = reactor.faucetAccount.WithAccountNumber(an).WithSequence(as).WithChainID(reactor.api.ChainID())
	return nil
}

func (reactor *stormReactor) initActionReactor(flags *pflag.FlagSet) error {
	// simple action parser
	actions, err := flags.GetStringSlice("actions")
	if err != nil {
		return err
	}
	reactor.actionReactor = &stormActions.ActionReactor{}
	for _, act := range actions {
		ss := strings.Split(act, "=")
		if len(ss) != 2 {
			return fmt.Errorf("'%s' must be actionName=weight", act)
		}
		generatorName := ss[0]
		weight, err := strconv.ParseInt(ss[1], 10, 64)
		if err != nil {
			return fmt.Errorf("'%s' weight must be integer, go error '%s'", act, err.Error())
		}
		err = reactor.actionReactor.Add(generatorName, weight)
		if err != nil {
			return fmt.Errorf("'%s': %s", act, err.Error())
		}
	}
	return nil
}

func (reactor *stormReactor) initLimiter(flags *pflag.FlagSet) error {
	limit, err := flags.GetInt64("tps")
	if err != nil {
		return err
	}
	reactor.limiter = stormActions.NewTPSLimiter(limit)
	return nil
}

func (reactor *stormReactor) updateGeneratorsInfo() {
	// update info
	ui := stormActions.UpdateInfo{}
	ui.MultisigBalances = make(map[string]sdk.Coins)

	coins, err := reactor.api.Coins()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range coins {
		ui.Coins = append(ui.Coins, c.Symbol)
		ui.FullCoins = append(ui.FullCoins, c)
	}
	for _, acc := range reactor.accounts {
		ui.Addresses = append(ui.Addresses, acc.Address())
	}
	// nft
	nfts := make([]dscApi.NFT, 0)
	denoms, err := reactor.api.NFTCollections()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, denom := range denoms {
		coll, err := reactor.api.NFTCollection(denom)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, id := range coll.NFTs {
			nft, err := reactor.api.NFT(denom, id)
			if err != nil {
				fmt.Println(err)
				return
			}
			nfts = append(nfts, nft)
		}
	}
	ui.NFTs = nfts
	// nft subtokens
	ui.NFTSubTokenReserves = make(map[stormActions.NFTSubTokenKey]sdkmath.Int)
	for _, nft := range ui.NFTs {
		subIds := []uint64{}
		for _, o := range nft.Owners {
			subIds = append(subIds, o.SubTokenIDs...)
		}
		subtokens, err := reactor.api.NFTSubTokens(nft.Denom, nft.ID, subIds)
		if err != nil {
			fmt.Println(err)
			return
		}
		for i := range subtokens {
			ui.NFTSubTokenReserves[stormActions.NFTSubTokenKey{Denom: nft.Denom, TokenID: nft.ID, ID: subtokens[i].ID}] = subtokens[i].Reserve.Amount
		}
	}
	// multisig wallets
	for _, owner := range ui.Addresses {
		wallets, err := reactor.api.MultisigWalletsByOwner(owner)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, wallet := range wallets {
			doAdd := true
			for _, w := range ui.MultisigWallets {
				if wallet.Address == w.Address {
					doAdd = false
					break
				}
			}
			if doAdd {
				ui.MultisigWallets = append(ui.MultisigWallets, wallet)
			}
		}
	}
	// multisig transactions
	for _, wallet := range ui.MultisigWallets {
		txs, err := reactor.api.MultisigTransactionsByWallet(wallet.Address)
		if err != nil {
			fmt.Println(err)
			return
		}
		ui.MultisigTransactions = append(ui.MultisigTransactions, txs...)
	}
	// multisig balances
	for _, wallet := range ui.MultisigWallets {
		balance, err := reactor.api.AddressBalance(wallet.Address)
		if err != nil {
			fmt.Println(err)
			return
		}
		ui.MultisigBalances[wallet.Address] = balance
	}
	reactor.actionReactor.Update(ui)
}

func (reactor *stormReactor) coins() ([]dscApi.Coin, error) {
	coins, err := reactor.api.Coins()
	if err != nil {
		return []dscApi.Coin{}, err
	}
	return coins, nil
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
