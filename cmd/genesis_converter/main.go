package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	dscTypes "bitbucket.org/decimalteam/go-smart-node/types"
)

var globalBaseDenom string

func usage() {
	fmt.Println(`usage: ./genesis_converter -basedenom <del|tdel> -decimal <decimal_genesis_file> -source <dsc_params_source_genesis> -nftfix <nft_owners_fix_file>
		-nftdublicates <path_to_file> -nfturiprefix <https://wherebuynft.com/api/nfts/ | https://devnet-nft.decimalchain.com/api/nfts/ | https://testnet-nft.decimalchain.com/api/nfts/>
	 -result <dsc_result_genesis_file> -injectlegacy <true|false>`)
}

func main() {
	initConfig()

	var pathGenesisOld, pathGenesisSource, pathGenesisResult, pathNFTfix, pathExportNFTDublicates string
	var nftUriPrefix string
	var injectLegacy bool
	var setValidatorsOnline string
	var setOnline []string

	flag.StringVar(&globalBaseDenom, "basedenom", "del", "base denom for blockchain (del, tdel)")
	flag.StringVar(&pathGenesisOld, "decimal", "", "path to exported genesis from Decimal")
	flag.StringVar(&pathGenesisSource, "source", "", "path to source genesis from DSC to copy parameters and validators")
	flag.StringVar(&pathGenesisResult, "result", "", "path to result genesis for DSC")
	flag.StringVar(&pathNFTfix, "nftfix", "", "path to json with nftfix (may be empty)")
	flag.StringVar(&pathExportNFTDublicates, "nftdublicates", "", "path to json to export dublicates info (may be empty)")
	flag.StringVar(&nftUriPrefix, "nfturiprefix", "", "url prefix to fix URI dublicates (https://wherebuynft.com/api/nfts/ | https://devnet-nft.decimalchain.com/api/nfts/ | https://testnet-nft.decimalchain.com/api/nfts/)")
	flag.BoolVar(&injectLegacy, "injectlegacy", false, "generate legacy records")
	flag.StringVar(&setValidatorsOnline, "setonline", "", "comma separated list of validators to set online, all others will be offline")

	flag.Parse()

	fmt.Printf("baseDenom=%s\n", globalBaseDenom)
	fmt.Printf("pathGenesisOld=%s\n", pathGenesisOld)
	fmt.Printf("pathGenesisResult=%s\n", pathGenesisResult)
	fmt.Printf("pathGenesisSource=%s\n", pathGenesisSource)
	if globalBaseDenom == "" || pathGenesisOld == "" || pathGenesisResult == "" || pathGenesisSource == "" || nftUriPrefix == "" || pathExportNFTDublicates == "" {
		usage()
		os.Exit(1)
	}

	if len(setValidatorsOnline) > 0 {
		setOnline = strings.Split(setValidatorsOnline, ",")
	}

	gsOld := readGenesisOld(pathGenesisOld)
	gsSource := readGenesisNew(pathGenesisSource)
	fixNFTData := readNFTFix(pathNFTfix)
	gsNew, nftDublicatesRecords, err := convertGenesis(gsOld, fixNFTData, injectLegacy, nftUriPrefix, setOnline)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	copyParams(&gsNew, gsSource)
	fixBondedNotBondedPools(&gsNew)
	// fixNFTPool(&gsNew)
	fixCoinVolumes(&gsNew)
	fixAccountNumbers(&gsNew)
	writeGenesisNew(pathGenesisResult, &gsNew)
	exportNFTDublicates(pathExportNFTDublicates, nftDublicatesRecords)
}

// Init global cosmos sdk config
func initConfig() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(cmdcfg.Bech32PrefixAccAddr, cmdcfg.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(cmdcfg.Bech32PrefixValAddr, cmdcfg.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(cmdcfg.Bech32PrefixConsAddr, cmdcfg.Bech32PrefixConsPub)
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

func readNFTFix(fpath string) []NFTOwnerFixRecord {
	if fpath == "" {
		return []NFTOwnerFixRecord{}
	}
	f, err := os.Open(fpath)
	if err != nil {
		panic(fmt.Errorf("file open %s error: %s", fpath, err.Error()))
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Errorf("file read %s error: %s", fpath, err.Error()))
	}
	var res []NFTOwnerFixRecord
	err = json.Unmarshal(data, &res)
	if err != nil {
		panic(fmt.Errorf("unmarshal %s error: %s", fpath, err.Error()))
	}
	return res
}

func convertGenesis(gsOld *GenesisOld, fixNFTData []NFTOwnerFixRecord, injectLegacy bool, nftUriPrefix string, setOnline []string) (GenesisNew, []nftDublicatesRecord, error) {
	var gsNew GenesisNew
	var err error

	gsNew.InitalHeight = strconv.FormatInt(gsOld.AppState.LastHeight+1, 10)

	// old-new adresses table, multisig addresses table, module addresses
	addrTable, err := prepareAddressTable(gsOld)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	// modules
	addrTable.InitModules()
	// accounts
	accsNew, err := convertAccounts(gsOld.AppState.Auth.Accounts, addrTable)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	gsNew.AppState.Auth.Accounts = accsNew
	// coins
	gsNew.AppState.Coin.Coins, err = convertCoins(gsOld.AppState.Coin.Coins, addrTable)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	// balances
	legacyRecords := NewLegacyRecords()
	gsNew.AppState.Bank.Balances, err =
		convertBalances(gsOld.AppState.Auth.Accounts, addrTable, legacyRecords, gsNew.AppState.Coin.Coins)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	// multisig wallets
	gsNew.AppState.Multisig.Wallets, err = convertMultisigWallets(gsOld.AppState.Multisig.Wallets, addrTable, legacyRecords)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	// TODO: convert for new transaction type
	// transactions
	gsNew.AppState.Multisig.Transactions, err =
		convertMultisigTransactions(gsOld.AppState.Multisig.Transactions, addrTable, gsOld.AppState.Multisig.Wallets, gsNew.AppState.Coin.Coins)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}

	// nft
	delegationCache, err := createNFTDelegationCache(gsOld.AppState.Validator.DelegationsNFT, gsOld.AppState.Validator.UndondingsNFT,
		gsOld.AppState.Validator.Validators, addrTable)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	gsNew.AppState.NFT.Collections, err =
		convertNFT(gsOld.AppState.NFT.Collections, gsOld.AppState.NFT.SubTokens, addrTable, legacyRecords, fixNFTData, delegationCache)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	// fix URI dublicates
	nftDublicatesRecords := extractNFTDublicates(gsNew.AppState.NFT.Collections)
	generateReplacements(nftUriPrefix, &nftDublicatesRecords)
	fixNFTDublicates(&(gsNew.AppState.NFT.Collections), nftDublicatesRecords)

	// inject to legacy records
	// TODO: parametrize output?
	if injectLegacy {
		fmt.Printf("!!! LEGACY RECORDS INJECTING !!! REMOVE FROM PRODUCTION\n")
		_, err := os.Stat("legacy_test_mnemonics.txt")
		if errors.Is(err, os.ErrNotExist) {
			out, err := os.Create("legacy_test_mnemonics.txt")
			if err != nil {
				return GenesisNew{}, []nftDublicatesRecord{}, err
			}
			for i := 0; i < 40000; i++ {
				mn, _ := dscWallet.NewMnemonic("")
				acc, _ := dscWallet.NewAccountFromMnemonicWords(mn.Words(), "")
				legacy_address, _ := dscTypes.GetLegacyAddressFromPubKey(acc.PubKey().Bytes())
				address := acc.Address()
				out.WriteString(fmt.Sprintf("%s\t%s\t%s\n", legacy_address, address, mn.Words()))
			}
			out.Close()
		}
		inp, err := os.Open("legacy_test_mnemonics.txt")
		if err != nil {
			return GenesisNew{}, []nftDublicatesRecord{}, err
		}
		coins := sdk.NewCoins(sdk.NewCoin(globalBaseDenom, sdkmath.NewInt(3_141_592_653_589_793_238)))
		allCoins := sdk.NewCoins()
		fileScanner := bufio.NewScanner(inp)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {
			row := strings.Split(fileScanner.Text(), "\t")
			legacyRecords.AddCoins(row[0], coins)
			allCoins = allCoins.Add(coins...)
		}
		inp.Close()
		// add to pool
		legAddr := addrTable.GetModule("legacy_coin_pool").address
		for i := range gsNew.AppState.Bank.Balances {
			if gsNew.AppState.Bank.Balances[i].Address == legAddr {
				gsNew.AppState.Bank.Balances[i].Coins = gsNew.AppState.Bank.Balances[i].Coins.Add(allCoins...)
			}
		}
	}
	// validators
	gsNew.AppState.Validator.Validators, err =
		convertValidators(gsOld.AppState.Validator.Validators, addrTable, legacyRecords, setOnline)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	gsNew.AppState.Validator.Delegations, err =
		convertDelegations(gsOld.AppState.Validator.Delegations, gsOld.AppState.Validator.DelegationsNFT,
			gsNew.AppState.Coin.Coins, addrTable)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	gsNew.AppState.Validator.Undelegations, err =
		convertUnbondings(gsOld.AppState.Validator.Unbondings, gsOld.AppState.Validator.UndondingsNFT,
			gsNew.AppState.Coin.Coins, addrTable)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	gsNew.AppState.Validator.LastValidatorPowers, err =
		convertLastValidatorPowers(gsOld.AppState.Validator.LastValidatorPowers, gsNew.AppState.Validator.Validators, addrTable)
	if err != nil {
		return GenesisNew{}, []nftDublicatesRecord{}, err
	}
	for _, pwr := range gsNew.AppState.Validator.LastValidatorPowers {
		gsNew.AppState.Validator.LastTotalPower += pwr.Power
	}
	fixDelegatedNFT(&gsNew, addrTable)

	// legacy records
	var records []LegacyRecordNew
	for _, v := range legacyRecords.data {
		if v.Coins.Len() == 0 && len(v.NFTs) == 0 && len(v.Wallets) == 0 && len(v.Validators) == 0 {
			continue
		}
		records = append(records, *v)
	}
	gsNew.AppState.Legacy.LegacyRecords = records

	//////////////////////////////////////////
	// validate NFT subtokens
	/*
		invalidSubtokens := verifySubtokens(gsOld.AppState.NFT.SubTokens, gsOld.AppState.NFT.Collections,
			gsOld.AppState.Validator.DelegationsNFT, gsOld.AppState.Validator.UndondingsNFT)
		for key, cnt := range invalidSubtokens {
			fmt.Printf("invalid subtoken nft: %#v == %#v\n", key, *cnt)
		}
	*/
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
	verifyNFTSupply(gsNew.AppState.NFT.Collections)

	verifyPools(gsNew.AppState.Bank.Balances, gsNew.AppState.Validator.Validators, gsNew.AppState.Validator.Delegations,
		gsNew.AppState.Validator.Undelegations, addrTable)

	verifyNFTDelegations(&gsNew, addrTable)

	// DUMP OLD-NEW VALIDATORS
	fmt.Printf("DUMP OLD-NEW VALIDATORS\n")
	for _, oldVal := range gsOld.AppState.Validator.Validators {
		newAddress := addrTable.GetValidatorAddress(oldVal.ValAddress)
		fmt.Printf("%s\t%s\t%s\n", oldVal.ValAddress, newAddress, oldVal.Description.Moniker)
	}

	return gsNew, nftDublicatesRecords, nil
}
