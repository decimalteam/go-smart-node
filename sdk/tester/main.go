package main

import (
	"encoding/json"
	"fmt"
	"strings"

	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//helper function
func formatAsJSON(obj interface{}) string {
	objStr, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s\n", objStr)
}

var mult = sdk.NewInt(1000000000000000000)

// TODO: split and document
func main() {
	api, _ := dscApi.NewAPI(dscApi.ConnectionOptions{EndpointHost: "http://127.0.0.1", Timeout: 40})
	err := api.GetParameters()
	if err != nil {
		fmt.Printf("get ChainID error: %v\n", err)
		return
	}
	chainID := api.ChainID()
	fmt.Printf("ChainID=%s\n", chainID)
	for _, addr := range []string{"dx10n5429jflppqv9grprd3ywglzlheelarcvzsk6", "dx1r75hkp8q6ay535tskjdskynz3y3373pujqxyfr"} {
		res, err := api.Address(addr)
		if err != nil {
			fmt.Printf("get Address(%s) error: %v\n", addr, err)
			continue
		}
		fmt.Printf("Address=%s\n", formatAsJSON(res))
	}
	coins, err := api.Coins()
	if err != nil {
		fmt.Printf("api.Coins() error: %v\n", err)
	}
	fmt.Printf("api.Coins() result: %v\n", coins)
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
			sdk.NewCoin(api.BaseCoin(), helpers.EtherToWei(sdk.NewInt(10))),
			w.SdkAddress(),
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

	nftDenoms, err := api.NFTCollections()
	if err != nil {
		fmt.Printf("NFTCollections err = %v\n", err)
	} else {
		fmt.Printf("NFT collections:\n%s\n", strings.Join(nftDenoms, "\n"))
	}

	// all nft
	nfts := make([]dscApi.NFT, 0)
	denoms, err := api.NFTCollections()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, denom := range denoms {
		nftIds, err := api.NFTsByDenom(denom)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, id := range nftIds {
			nft, err := api.NFT(denom, id)
			if err != nil {
				fmt.Println(err)
				return
			}
			nfts = append(nfts, *nft)
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
