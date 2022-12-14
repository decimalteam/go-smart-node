package main

import (
	"bufio"
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
	feeConfig     *stormTypes.FeeConfiguration
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

	useCustomFee, err := flags.GetBool(customFee)
	if err != nil {
		return err
	}
	reactor.feeConfig = stormTypes.NewFeeConfiguration(useCustomFee)
	err = reactor.feeConfig.Update(reactor.api)
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
		err = reactor.actionReactor.Add(generatorName, weight, reactor.api.BaseCoin())
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
	err := reactor.feeConfig.Update(reactor.api)
	if err != nil {
		fmt.Println(err)
		return
	}
	// update info
	ui := stormActions.UpdateInfo{}
	ui.MultisigBalances = make(map[string]sdk.Coins)

	coins, err := reactor.api.Coins()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range coins {
		ui.Coins = append(ui.Coins, c.Denom)
		ui.FullCoins = append(ui.FullCoins, c)
	}
	for _, acc := range reactor.accounts {
		ui.Addresses = append(ui.Addresses, acc.Address())
	}
	// nft
	nfts := make([]*dscApi.NFTToken, 0)
	colls, err := reactor.api.NFTCollections()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, coll := range colls {
		collWithTokens, err := reactor.api.NFTCollection(coll.Creator, coll.Denom)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nfts = append(nfts, collWithTokens.Tokens...)
	}
	ui.NFTs = nfts
	// nft subtokens
	ui.NFTSubTokenReserves = make(map[stormActions.NFTSubTokenKey]sdk.Coin)
	for j := range ui.NFTs {
		nft := ui.NFTs[j]
		tok, err := reactor.api.NFTToken(nft.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for i := range tok.SubTokens {
			ui.NFTSubTokenReserves[stormActions.NFTSubTokenKey{TokenID: nft.ID, ID: tok.SubTokens[i].ID}] = *tok.SubTokens[i].Reserve
		}
		ui.NFTs[j] = &tok
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
	// TODO: rework
	for _, wallet := range ui.MultisigWallets {
		txs, err := reactor.api.MultisigTransactionsByWallet(wallet.Address)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, tx := range txs {
			txInfo, err := reactor.api.MultisigTransactionsByID(tx.Id)
			if err != nil {
				fmt.Println(err)
				return
			}
			ui.MultisigTransactions = append(ui.MultisigTransactions, txInfo)
		}
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

	// validators
	ui.Validators, err = reactor.api.Validators()
	if err != nil {
		fmt.Println(err)
		return
	}

	// delegations
	for _, val := range ui.Validators {
		dels, err := reactor.api.ValidatorDelegations(val.OperatorAddress)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ui.Delegations = append(ui.Delegations, dels...)
	}
	// undelegations
	for _, val := range ui.Validators {
		undels, err := reactor.api.ValidatorUndelegations(val.OperatorAddress)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ui.Undelegations = append(ui.Undelegations, undels...)
	}
	// redelegations
	for _, val := range ui.Validators {
		redels, err := reactor.api.ValidatorRedelegations(val.OperatorAddress)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ui.Redelegations = append(ui.Redelegations, redels...)
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

type stepsCounter struct {
	limit   int
	counter int
}

func NewStepsCounter(flags *pflag.FlagSet) *stepsCounter {
	limit, err := flags.GetInt(stepsCount)
	if err != nil || limit < 0 {
		limit = 0
	}
	return &stepsCounter{
		limit:   limit,
		counter: 0,
	}
}

func (sc *stepsCounter) increment() bool {
	if sc.limit == 0 {
		return true
	}
	sc.counter++
	if sc.counter > sc.limit {
		return false
	}
	return true
}
