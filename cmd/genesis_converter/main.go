package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Init global cosmos sdk config
func initConfig() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(config.Bech32PrefixAccAddr, config.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(config.Bech32PrefixValAddr, config.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(config.Bech32PrefixConsAddr, config.Bech32PrefixConsPub)
}

func readGenesisNew(fpath string) *GenesisNew {
	f, err := os.Open(fpath)
	if err != nil {
		panic(fmt.Errorf("file open %s error: %s", fpath, err.Error()))
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Errorf("file read %s error: %s", fpath, err.Error()))
	}
	var gs GenesisNew
	err = json.Unmarshal(data, &gs)
	if err != nil {
		panic(fmt.Errorf("unmarshal %s error: %s", fpath, err.Error()))
	}
	return &gs
}

func writeGenesisNew(fpath string, gs *GenesisNew) {
	f, err := os.Create(fpath)
	if err != nil {
		panic(fmt.Errorf("file %s error: %s", fpath, err.Error()))
	}
	defer f.Close()
	data, err := json.MarshalIndent(gs, "", "  ")
	if err != nil {
		panic(fmt.Errorf("marshal new genesis error: %s", err.Error()))
	}
	_, err = f.Write(data)
	if err != nil {
		panic(fmt.Errorf("file %s error: %s", fpath, err.Error()))
	}
}

func readGenesisOld(fpath string) *GenesisOld {
	f, err := os.Open(fpath)
	if err != nil {
		panic(fmt.Errorf("file open %s error: %s", fpath, err.Error()))
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Errorf("file read %s error: %s", fpath, err.Error()))
	}
	var gs GenesisOld
	err = json.Unmarshal(data, &gs)
	if err != nil {
		panic(fmt.Errorf("unmarshal %s error: %s", fpath, err.Error()))
	}
	return &gs
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: ./genesis_converter <decimal_genesis_file> <dsc_genesis_file> <dsc_params_source_genesis>")
		os.Exit(1)
	}
	gsOld := readGenesisOld(os.Args[1])
	gsSource := readGenesisNew(os.Args[3])
	gsNew, _, err := convertGenesis(gsOld)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	copyParams(&gsNew, gsSource)
	writeGenesisNew(os.Args[2], &gsNew)
}

type Statistic struct {
	countRegularAccounts            uint64
	countRegularAccountsNoPublicKey uint64
}

func convertGenesis(gsOld *GenesisOld) (GenesisNew, Statistic, error) {
	var gsNew GenesisNew
	var err error
	// old-new adresses table, multisig addresses table, module addresses
	addrTable, err := prepareAddressTable(gsOld)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	// modules
	addrTable.InitModules()
	// accounts
	accsNew, err := convertAccounts(gsOld.AppState.Auth.Accounts, addrTable)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	gsNew.AppState.Auth.Accounts = accsNew
	// balances
	legacyRecords := NewLegacyRecords()
	gsNew.AppState.Bank.Balances, err =
		convertBalances(gsOld.AppState.Auth.Accounts, addrTable, legacyRecords)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	// coins
	gsNew.AppState.Coin.Coins, err = convertCoins(gsOld.AppState.Coin.Coins, addrTable)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	// multisig wallets
	gsNew.AppState.Multisig.Wallets, err = convertMultisigWallets(gsOld.AppState.Multisig.Wallets, addrTable, legacyRecords)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	// transactions
	gsNew.AppState.Multisig.Transactions, err =
		convertMultisigTransactions(gsOld.AppState.Multisig.Transactions, addrTable, gsOld.AppState.Multisig.Wallets)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	// nft
	gsNew.AppState.NFT.Collections, gsNew.AppState.NFT.NFTs, err =
		convertNFT(gsOld.AppState.NFT.Collections, addrTable, legacyRecords)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	gsNew.AppState.NFT.SubTokens, err =
		convertSubTokens(gsOld.AppState.NFT.SubTokens, gsNew.AppState.NFT.NFTs)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	// legacy records
	var records []LegacyRecordNew
	for _, v := range legacyRecords.data {
		records = append(records, *v)
	}
	gsNew.AppState.Legacy.LegacyRecords = records
	// validate NFT subtokens
	invalidSubtokens := verifySubtokens(gsOld.AppState.NFT.SubTokens, gsOld.AppState.NFT.Collections,
		gsOld.AppState.Validator.DelegationsNFT, gsOld.AppState.Validator.UndondingNFT)
	for key, cnt := range invalidSubtokens {
		fmt.Printf("invalid subtoken nft: %#v == %#v\n", key, *cnt)
	}
	// validate coins
	coinDiffs := verifyCoinsVolume(gsOld.AppState.Coin.Coins, gsOld.AppState.Auth.Accounts,
		gsOld.AppState.Validator.Delegations, gsOld.AppState.Validator.Unbondings)
	for _, diff := range coinDiffs {
		if !diff.BCSum.Equal(diff.Volume) {
			fmt.Printf("%s invalid coin volume (sum in blockchain, volume in storage):  %s != %s (dif:%v)\n",
				diff.Symbol, diff.BCSum.String(), diff.Volume.String(), diff.BCSum.Sub(diff.Volume))
		}
	}
	// validate NFT colections
	verifyNFTSupply(gsNew.AppState.NFT.Collections, gsNew.AppState.NFT.NFTs)

	return gsNew, Statistic{}, nil
}
