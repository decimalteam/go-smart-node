package scenario

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/actions"
	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
)

const (
	ValCountFlag     = "val_count"
	CoinsCountFlag   = "coin_count"
	NftsCountFlag    = "nfts_count"
	DelegationsFlag  = "delegations_count"
	CreateValidators = "create_validators"
	CreateCoins      = "create_coins"
	CreateNfts       = "create_nfts"
)

//const charsAll = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
//const charsAbc = "abcdefghijklmnopqrstuvwxyz"
//const charsUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type DelegationsReadScenario struct {
	accs       []*stormTypes.StormAccount
	nfts       []dscApi.NFTToken
	coins      []dscApi.Coin
	validators []dscApi.Validator

	// generators
	coinGenerator            *actions.CreateCoinGenerator
	sendCoinGenerator        *actions.SendCoinGenerator
	nftGenerator             *actions.MintNFTGenerator
	createValidatorGenerator *actions.CreateValidatorGenerator
	delegateCoinGenerator    *actions.DelegateGenerator
	delegateNftGenerator     *actions.DelegateNFTGenerator

	faucet    *stormTypes.StormAccount
	feeConfig *stormTypes.FeeConfiguration
	api       *dscApi.API
	rnd       *rand.Rand
}

func NewDelegationsReadScenario(api *dscApi.API, accs []*stormTypes.StormAccount) DelegationsReadScenario {
	coinGenerator := actions.NewCreateCoinGenerator(
		6,
		8,
		10000,
		100000,
		10000,
		12000,
		100000,
		100000000,
	)

	sendCoinGenerator := actions.NewSendCoinGenerator(500, 20000)

	nftGenerator := actions.NewMintNFTGenerator(
		1,
		100,
		100,
		1000,
		1,
		10,
	)
	createValidatorGenerator := actions.NewCreateValidatorGenerator(1, 10)
	delegateCoin := actions.NewDelegateGenerator(1, 3)
	delegateNft := actions.NewDelegateNFTGenerator()

	return DelegationsReadScenario{
		accs:       accs,
		nfts:       make([]dscApi.NFTToken, 0),
		coins:      make([]dscApi.Coin, 0),
		validators: make([]dscApi.Validator, 0),

		coinGenerator:            coinGenerator,
		sendCoinGenerator:        sendCoinGenerator,
		nftGenerator:             nftGenerator,
		createValidatorGenerator: createValidatorGenerator,
		delegateCoinGenerator:    delegateCoin,
		delegateNftGenerator:     delegateNft,

		feeConfig: stormTypes.NewFeeConfiguration(false),
		api:       api,
		rnd:       rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (d *DelegationsReadScenario) CreateCoins(coinCount int) error {
	ui, err := d.GenerateUpdatesInfo()
	if err != nil {
		panic(err)
	}

	for i := range d.accs {
		d.accs[i].UpdateBalance()
		d.accs[i].UpdateNumberSequence()
	}

	d.coinGenerator.Update(ui)
	for i := 0; i < coinCount; i++ {
		accIndex, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(d.accs))))
		acc := d.accs[accIndex.Int64()]

		action := d.coinGenerator.Generate()
		bz, err := action.GenerateTx(acc, d.feeConfig)
		if err != nil {
			return err
		}

		resp, err := d.api.BroadcastTxSync(bz)
		if err != nil {
			fmt.Printf("CreateNFTs-BytesToSend err: %s\n", err.Error())
			continue
		}
		if resp.Code != 0 {
			fmt.Printf("TxResult code=%d, codespace=%s, msg=%s\n", resp.Code, resp.Codespace, resp.Log)
		}
		acc.UpdateBalance()
		acc.IncrementSequence()
		if i%10 == 0 {
			time.Sleep(5 * time.Second)
		}
	}

	coins, err := d.api.Coins()
	if err != nil {
		return err
	}

	//for _,coin := range coins {
	//	accIndex := rand.Int63n(int64(len(d.accs)))
	//	acc := d.accs[accIndex]
	//	acc.UpdateBalance()
	//	acc.UpdateNumberSequence()
	//
	//}
	d.coins = coins

	fmt.Println("coins created!")
	return nil
}

func (d *DelegationsReadScenario) CreateNFTs(nftCount int) error {
	ui, err := d.GenerateUpdatesInfo()
	if err != nil {
		panic(err)
	}

	for i := range d.accs {
		d.accs[i].UpdateBalance()
		d.accs[i].UpdateNumberSequence()
	}

	d.nftGenerator.Update(ui)
	for i := 0; i < nftCount; i++ {
		accIndex, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(d.accs))))
		acc := d.accs[accIndex.Int64()]

		action := d.nftGenerator.Generate()
		bz, err := action.GenerateTx(acc, d.feeConfig)
		if err != nil {
			return err
		}

		resp, err := d.api.BroadcastTxSync(bz)
		if err != nil {
			fmt.Printf("CreateNFTs-BytesToSend err: %s\n", err.Error())
			continue
		}
		if resp.Code != 0 {
			fmt.Printf("TxResult code=%d, codespace=%s, msg=%s\n", resp.Code, resp.Codespace, resp.Log)
		}

		d.accs[i].UpdateBalance()
		acc.IncrementSequence()

		if i%10 == 0 {
			time.Sleep(5 * time.Second)
		}
	}

	//collections, err := d.api.NFTCollections()
	//if err != nil {
	//	return err
	//}
	//fmt.Println(collections)
	//tokens := make([]dscApi.NFTToken, len(collections[0].Tokens))
	//for i, v := range collections[0].Tokens {
	//	token, err := d.api.NFTToken(v.ID)
	//	if err != nil {
	//		return err
	//	}
	//	tokens[i] = token
	//}
	//d.nfts = tokens

	fmt.Println("nfts created!")
	return nil
}

func (d *DelegationsReadScenario) CreateValidators(validatorCount int) error {
	ui, err := d.GenerateUpdatesInfo()
	if err != nil {
		panic(err)
	}
	d.createValidatorGenerator.Update(ui)

	for i := range d.accs {
		d.accs[i].UpdateBalance()
		d.accs[i].UpdateNumberSequence()
	}

	valOwners := make(map[string]*stormTypes.StormAccount)
	for i := 0; i < validatorCount; i++ {
		accIndex, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(d.accs))))
		acc := d.accs[accIndex.Int64()]
		acc.UpdateBalance()
		acc.UpdateNumberSequence()

		action := d.createValidatorGenerator.Generate()
		bz, err := action.GenerateTx(acc, d.feeConfig)
		if err != nil {
			return err
		}
		resp, err := d.api.BroadcastTxSync(bz)
		if err != nil {
			fmt.Printf("CreateValidators-BytesToSend err: %s\n", err.Error())
			continue
		}
		if resp.Code != 0 {
			fmt.Printf("CreateValidators TxResult code=%d, codespace=%s, msg=%s\n", resp.Code, resp.Codespace, resp.Log)
		}
		acc.IncrementSequence()
		valOwners[sdk.ValAddress(acc.Account().SdkAddress()).String()] = acc
	}

	fmt.Println("validators created!")
	validators, err := d.api.Validators()
	if err != nil {
		return err
	}
	d.validators = validators
	time.Sleep(5 * time.Second)
	fmt.Println("change validators status to online")
	err = d.ValidatorToOnline(valOwners)
	if err != nil {
		return err
	}

	validators, err = d.api.Validators()
	if err != nil {
		return err
	}
	d.validators = validators
	fmt.Println("validators statuses updated!")
	return nil
}

func (d *DelegationsReadScenario) CreateAndSendDelegations(delegationsCount int) {
	for i := range d.accs {
		d.accs[i].UpdateBalance()
		d.accs[i].UpdateNumberSequence()
	}

	// update balances
	for i := 0; i < delegationsCount; i++ {
		// update infos
		ui, err := d.GenerateUpdatesInfo()
		if err != nil {
			fmt.Println(err)
		}
		d.delegateCoinGenerator.Update(ui)
		d.delegateNftGenerator.Update(ui)

		// generate action suite
		var action actions.Action
		if d.rnd.Intn(100)%2 == 0 {
			action = d.delegateCoinGenerator.Generate()
		} else {
			action = d.delegateNftGenerator.Generate()
		}
		// choose accounts
		accs := action.ChooseAccounts(d.accs)
		if len(accs) == 0 {
			fmt.Sprintf("no one account with this stake %s\n", action.String())
			continue
		}
		accIndex, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(accs)))) //rand.Int63n(int64(len(accs)/2)) + rand.Int63n(int64(len(accs)))/rand.Int63n(int64(len(accs)))
		acc := accs[int(accIndex.Int64())]
		acc.UpdateBalance()
		//acc.UpdateNumberSequence()
		// generate txs
		bz, err := action.GenerateTx(acc, d.feeConfig)
		if err != nil {
			fmt.Println(err)
		}

		// send tx
		_, err = d.api.BroadcastTxAsync(bz)
		if err != nil {
			fmt.Printf("Send Delegate-BytesToSend err: %s\n", err.Error())
			continue
		}
		//if resp.Code != 0 {
		//	fmt.Printf("Send Delegate TxResult code=%d, codespace=%s, msg=%s\n", resp.Code, resp.Codespace, resp.Log)
		//}
		acc.IncrementSequence()
		//if i%100 == 0 {
		//	time.Sleep(time.Second * 7)
		//}
	}
}

func (d *DelegationsReadScenario) GenerateUpdatesInfo() (actions.UpdateInfo, error) {

	ui := actions.UpdateInfo{}
	ui.MultisigBalances = make(map[string]sdk.Coins)

	coins, err := d.api.Coins()
	if err != nil {
		return ui, err
	}
	for _, c := range coins {
		ui.Coins = append(ui.Coins, c.Denom)
		ui.FullCoins = append(ui.FullCoins, c)
	}
	for _, acc := range d.accs {
		ui.Addresses = append(ui.Addresses, acc.Address())
	}
	// nft
	nfts := make([]*dscApi.NFTToken, 0)
	colls, err := d.api.NFTCollections()
	if err != nil {
		return ui, err
	}
	for _, coll := range colls {
		collWithTokens, err := d.api.NFTCollection(coll.Creator, coll.Denom)
		if err != nil {
			continue
		}
		nfts = append(nfts, collWithTokens.Tokens...)
	}
	ui.NFTs = nfts

	// nft subtokens
	ui.NFTSubTokenReserves = make(map[actions.NFTSubTokenKey]sdk.Coin)
	for j := range ui.NFTs {
		nft := ui.NFTs[j]
		tok, err := d.api.NFTToken(nft.ID)
		if err != nil {
			continue
		}
		for i := range tok.SubTokens {
			ui.NFTSubTokenReserves[actions.NFTSubTokenKey{TokenID: nft.ID, ID: tok.SubTokens[i].ID}] = *tok.SubTokens[i].Reserve
		}
		ui.NFTs[j] = &tok
	}

	// multisig wallets
	for _, owner := range ui.Addresses {
		wallets, err := d.api.MultisigWalletsByOwner(owner)
		if err != nil {
			return ui, err
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
		txs, err := d.api.MultisigTransactionsByWallet(wallet.Address)
		if err != nil {
			return ui, err
		}
		for _, tx := range txs {
			txInfo, err := d.api.MultisigTransactionsByID(tx.Id)
			if err != nil {
				return ui, err
			}
			ui.MultisigTransactions = append(ui.MultisigTransactions, txInfo)
		}
	}
	// multisig balances
	for _, wallet := range ui.MultisigWallets {
		balance, err := d.api.AddressBalance(wallet.Address)
		if err != nil {
			return ui, err
		}
		ui.MultisigBalances[wallet.Address] = balance
	}

	vals, err := d.api.Validators()
	for _, v := range vals {
		ui.Validators = append(ui.Validators, v)
	}

	return ui, err
}

func (d *DelegationsReadScenario) ValidatorToOnline(valOwners map[string]*stormTypes.StormAccount) error {
	for _, val := range d.validators {
		acc, ok := valOwners[val.OperatorAddress]
		if !ok {
			continue
		}
		acc.UpdateBalance()
		//acc.UpdateNumberSequence()

		msg := dscTx.NewMsgSetOnline(
			val.GetOperator(),
		)
		tx, err := dscTx.BuildTransaction(acc.Account(), []sdk.Msg{msg}, "", "del", d.api.GetFeeCalculationOptions())
		if err != nil {
			fmt.Printf("SetOnline-BuildTransaction err: %s\n", err.Error())
			continue
		}
		err = tx.SignTransaction(acc.Account())
		if err != nil {
			fmt.Printf("SetOnline-SignTransaction err: %s\n", err.Error())
			continue
		}
		bz, err := tx.BytesToSend()
		if err != nil {
			fmt.Printf("SetOnline-BytesToSend err: %s\n", err.Error())
			continue
		}
		resp, err := d.api.BroadcastTxSync(bz)
		if err != nil {
			fmt.Printf("SetOnline-BytesToSend err: %s\n", err.Error())
			continue
		}
		if resp.Code != 0 {
			fmt.Printf("TxResult code=%d, codespace=%s, msg=%s\n", resp.Code, resp.Codespace, resp.Log)
		}
	}
	return nil

}

func (d *DelegationsReadScenario) SendCoinsForAllUsers(count int) {
	ui, err := d.GenerateUpdatesInfo()
	if err != nil {
		panic(err)
	}
	d.sendCoinGenerator.Update(ui)
	for i := 0; i < count; i++ {
		accIndex := rand.Int63n(int64(len(d.accs)))
		acc := d.accs[accIndex]

		action := d.sendCoinGenerator.Generate()
		// generate txs
		bz, err := action.GenerateTx(acc, d.feeConfig)
		if err != nil {
			fmt.Println(err)
		}

		// send tx
		resp, err := d.api.BroadcastTxSync(bz)
		if err != nil {
			fmt.Printf("CreateNFTs-BytesToSend err: %s\n", err.Error())
			continue
		}
		if resp.Code != 0 {
			fmt.Printf("TxResult code=%d, codespace=%s, msg=%s\n", resp.Code, resp.Codespace, resp.Log)
		}
	}
}

func (d *DelegationsReadScenario) Start(validatorsCount, coinsCount, nftsCount, delegationsCount int, createValidators, createCoins, createNfts bool) {
	fmt.Println("Start scenario reconstruct")
	err := d.feeConfig.Update(d.api)
	if err != nil {
		panic(err)
	}

	if !createValidators {
		fmt.Println("start create validators")
		err = d.CreateValidators(validatorsCount)
		if err != nil {
			panic(err)
		}
	}

	if !createCoins {
		time.Sleep(6)
		fmt.Println("start create coins")
		err = d.CreateCoins(coinsCount)
		if err != nil {
			panic(err)
		}
	}

	if !createNfts {
		time.Sleep(6)
		fmt.Println("start create nfts")
		err = d.CreateNFTs(nftsCount)
		if err != nil {
			panic(err)
		}
		time.Sleep(7)
	}

	fmt.Println("start distribute delegations")
	d.CreateAndSendDelegations(delegationsCount)
}
