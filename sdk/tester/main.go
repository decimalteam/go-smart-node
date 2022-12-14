package main

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

// helper function
func formatAsJSON(obj interface{}) string {
	objStr, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s\n", objStr)
}

// TODO: split and document
func main() {
	universalMultiSig()
	return

	api, err := dscApi.NewAPI(dscApi.ConnectionOptions{EndpointHost: "127.0.0.1", Timeout: 40})
	if err != nil {
		fmt.Printf("connect error: %v\n", err)
		return
	}
	err = api.GetParameters()
	if err != nil {
		fmt.Printf("GetParameters error: %v\n", err)
		return
	}
	chainID := api.ChainID()
	fmt.Printf("ChainID=%s\n", chainID)
	fmt.Printf("BaseCoin=%s\n", api.BaseCoin())
	for _, addr := range []string{"dx1ag696xwwqlfec2p2v69498w034zw2udh9rr0nr"} {
		res, err := api.Address(addr)
		if err != nil {
			fmt.Printf("get Address(%s) error: %v\n", addr, err)
			continue
		}
		fmt.Printf("Address=%+v\n", res)
		blnc, err := api.AddressBalance(addr)
		if err != nil {
			fmt.Printf("get Balance(%s) error: %v\n", addr, err)
			continue
		}
		fmt.Printf("Balance=%+v\n", blnc)
	}
	fmt.Println("start coins")
	coins, err := api.Coins()
	if err != nil {
		fmt.Printf("api.Coins() error: %v\n", err)
	}
	fmt.Printf("api.Coins() result: %d == %v\n", len(coins), coins)
	/*
		///////////////
		faucetMnemonic := "domain kangaroo addict allow capital message young faculty diesel aware dry mirror tomato prepare census inflict diagram eye project modify question hip crater pelican"
		faucet, err := dscWallet.NewAccountFromMnemonicWords(faucetMnemonic, "")
		if err != nil {
			fmt.Printf("create wallet error: %v\n", err)
			return
		}
		an, as, err := api.AccountNumberAndSequence(faucet.Address())
		if err != nil {
			fmt.Printf("AccountNumberAndSequence error: %v\n", err)
			return
		}
		faucet = faucet.WithAccountNumber(an).WithSequence(as).WithChainID(chainID)
		mnemonics := []string{
			"possible hedgehog buddy desk smart camera frost vacant ridge robust seminar riot boost gauge jar aunt frozen morning system ordinary volcano rescue bind trust",
			"drip charge ridge between primary comic core fatigue evidence member fault tank tennis venue young lawsuit shock skull hybrid enlist shield opera please panther",
			"asthma science hawk hip piano enrich avoid myself divide seek number satoshi matter bunker question disease foster toward rare depth fame catch artefact woman",
		}
		for i, mn := range mnemonics {
			fmt.Printf("menmonic #%d\n", i)
			w, err := dscWallet.NewAccountFromMnemonicWords(mn, "")
			if err != nil {
				fmt.Printf("create wallet error: %v\n", err)
				return
			}
			tx, err := dscTx.BuildTransaction(faucet, []sdk.Msg{dscTx.NewMsgSendCoin(
				faucet.SdkAddress(),
				w.SdkAddress(),
				sdk.NewCoin(api.BaseCoin(), helpers.EtherToWei(sdk.NewInt(10))),
			)}, "some send", api.BaseCoin())
			if err != nil {
				fmt.Printf("BuildTransaction error: %v\n", err)
				return
			}
			err = tx.SignTransaction(faucet)
			if err != nil {
				fmt.Printf("SignTransaction error: %v\n", err)
				return
			}
			bytes, err := tx.BytesToSend()
			if err != nil {
				fmt.Printf("BytesToSend error: %v\n", err)
				return
			}
			r, err := api.BroadcastTxSync(bytes)
			if err != nil {
				fmt.Printf("err = %v\n", err)
				continue
			}
			fmt.Printf("result = %+v\n", r)
			faucet.IncrementSequence()
		}
	*/
	nftColls, err := api.NFTCollections()
	if err != nil {
		fmt.Printf("NFTCollections err = %v\n", err)
	} else {
		fmt.Printf("NFT collections:\n%s\n", formatAsJSON(nftColls))
	}

	// all nft
	nfts := make([]dscApi.NFTToken, 0)
	colls, err := api.NFTCollections()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, coll := range colls {
		for _, token := range coll.Tokens {
			nft, err := api.NFTToken(token.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
			nfts = append(nfts, nft)
		}
	}
	fmt.Printf("---\n%v\n---\n", nfts)

	// all multisig wallets by owner
	owners := []string{
		"dx1s4c5cak2u7l3ddu67jana9vtfwvj8ezdjdv7j4",
		"dx1fzulqsza5nqva7jesfjw7a3a2xwfq9kp24pm53",
		"dx1h5h43cmz892zaqfazxphnacgktzga9elmf0y27",
		"dx1yr7z7ts7v5gh688ay6pe0d384dm0y2hrh4madg",
		"dx10nd8yly2kzmezhlnka5s9et7hyvkvggcrl07rs",
	}
	for _, owner := range owners {
		wallets, err := api.MultisigWalletsByOwner(owner)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%v\n", wallets)
		}
	}
}

func universalMultiSig() {
	const faucetMnemonic = "load green feature bachelor interest glue amount typical define ankle blade enroll draft tiger lonely volcano slab brush arm lottery defy mountain exhaust always"
	const acc1mnemonic = "affair coral purse lounge fancy orbit region shine wagon fever frozen market equal coil mixed lottery will stand oil they pepper utility season fruit"
	const acc2mnemonic = "hard delay bag address subject dog flock cactus athlete legal arrange skull own elephant twelve switch sustain desert angle shop supply solid river aspect"
	const acc3mnemonic = "differ enter exhaust copy position gravity fun guide clump brisk confirm swarm salt stamp tape purpose country slam simple tourist fog load toddler warrior"
	const acc4mnemonic = "rural pause vacant couch dwarf soup isolate doll long market casino evolve employ reward barely laptop dilemma solar lesson pyramid oven trust organ mandate"

	api, err := dscApi.NewAPI(dscApi.ConnectionOptions{EndpointHost: "127.0.0.1", Timeout: 40})
	if err != nil {
		fmt.Printf("connect error: %v\n", err)
		return
	}
	err = api.GetParameters()
	if err != nil {
		fmt.Printf("GetParameters error: %v\n", err)
		return
	}

	faucet, _ := dscWallet.NewAccountFromMnemonicWords(faucetMnemonic, "")
	acc1, _ := dscWallet.NewAccountFromMnemonicWords(acc1mnemonic, "")
	acc2, _ := dscWallet.NewAccountFromMnemonicWords(acc2mnemonic, "")
	acc3, _ := dscWallet.NewAccountFromMnemonicWords(acc3mnemonic, "")
	acc4, _ := dscWallet.NewAccountFromMnemonicWords(acc4mnemonic, "")

	for _, acc := range []*dscWallet.Account{acc1, acc2, acc3} {
		bindAcc(api, faucet)
		msg := dscTx.NewMsgSendCoin(faucet.SdkAddress(), acc.SdkAddress(), sdk.NewCoin(api.BaseCoin(), helpers.EtherToWei(sdk.NewInt(100))))
		tx, _ := dscTx.BuildTransaction(faucet, []sdk.Msg{msg}, "", api.BaseCoin(), api.GetFeeCalculationOptions())
		tx.SignTransaction(faucet)
		bz, _ := tx.BytesToSend()
		res, _ := api.BroadcastTxCommit(bz)
		fmt.Printf("fill result: %#v\n\n", res)
	}

	{
		bindAcc(api, acc1)
		msg := dscTx.NewMsgCreateWallet(acc1.SdkAddress(), []string{acc1.Address(), acc2.Address(), acc3.Address()}, []uint32{1, 1, 1}, 3)
		tx, _ := dscTx.BuildTransaction(acc1, []sdk.Msg{msg}, "", api.BaseCoin(), api.GetFeeCalculationOptions())
		tx.SignTransaction(acc1)
		bz, _ := tx.BytesToSend()
		res, _ := api.BroadcastTxCommit(bz)
		fmt.Printf("create wallet result: %#v\n", res)
	}

	wallets, _ := api.MultisigWalletsByOwner(acc1.Address())
	if len(wallets) == 0 {
		fmt.Printf("no wallets\n\n")
		return
	}
	wal := wallets[0]
	fmt.Printf("wallet: %#v\n\n", wal)
	wAdr, _ := sdk.AccAddressFromBech32(wal.Address)
	{
		bindAcc(api, faucet)
		msg := dscTx.NewMsgSendCoin(faucet.SdkAddress(), wAdr, sdk.NewCoin(api.BaseCoin(), helpers.EtherToWei(sdk.NewInt(100))))
		tx, _ := dscTx.BuildTransaction(faucet, []sdk.Msg{msg}, "", api.BaseCoin(), api.GetFeeCalculationOptions())
		tx.SignTransaction(faucet)
		bz, _ := tx.BytesToSend()
		res, _ := api.BroadcastTxCommit(bz)
		fmt.Printf("fill wallet result: %#v\n\n", res)
	}
	{
		bindAcc(api, acc1)
		msg, _ := dscTx.NewMsgCreateTransaction(acc1.SdkAddress(), wal.Address,
			dscTx.NewMsgSendCoin(wAdr, acc4.SdkAddress(), sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(10)))),
		)
		tx, err := dscTx.BuildTransaction(acc1, []sdk.Msg{msg}, "", api.BaseCoin(), api.GetFeeCalculationOptions())
		if err != nil {
			fmt.Printf("NewMsgCreateUniversalTransaction: %v\n", err)
		}
		tx.SignTransaction(acc1)
		bz, _ := tx.BytesToSend()
		resp, err := api.SimulateTx(bz)
		fmt.Printf("simulate: %#v, error: %v", resp, err)
		return
		res, _ := api.BroadcastTxCommit(bz)
		fmt.Printf("create tx result: %#v\n\n", res)
	}

	mtxs, _ := api.MultisigTransactionsByWallet(wal.Address)
	for _, mtx := range mtxs {
		fmt.Printf("signing tx: %s\n\n", mtx.Id)
		fmt.Printf("signing tx: %#v\n\n", mtx)
		for _, acc := range []*dscWallet.Account{acc2, acc3} {
			bindAcc(api, acc)
			msg := dscTx.NewMsgSignTransaction(acc.SdkAddress(), mtx.Id)
			tx, _ := dscTx.BuildTransaction(acc, []sdk.Msg{msg}, "", api.BaseCoin(), api.GetFeeCalculationOptions())
			tx.SignTransaction(acc)
			bz, _ := tx.BytesToSend()
			res, _ := api.BroadcastTxCommit(bz)
			fmt.Printf("sign result: %#v\n\n", res)
		}
	}

	coins, _ := api.AddressBalance(acc4.Address())
	fmt.Printf("result balance: %s\n", coins.String())
}

func bindAcc(api *dscApi.API, acc *dscWallet.Account) error {
	an, as, err := api.AccountNumberAndSequence(acc.Address())
	if err != nil {
		return fmt.Errorf("%w: AccountNumberAndSequence", err)
	}
	acc.WithAccountNumber(an).WithSequence(as).WithChainID(api.ChainID())
	return nil
}
