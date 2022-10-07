package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	dscTypes "bitbucket.org/decimalteam/go-smart-node/types"
)

func main() {
	initConfig()
	if len(os.Args) < 4 {
		fmt.Println("usage: ./genesis_converter <decimal_genesis_file> <dsc_params_source_genesis> <nft_owners_fix_file> <dsc_genesis_file>")
		os.Exit(1)
	}
	gsOld := readGenesisOld(os.Args[1])
	gsSource := readGenesisNew(os.Args[2])
	fixNFTData := readNFTFix(os.Args[3])
	gsNew, _, err := convertGenesis(gsOld, fixNFTData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	copyParams(&gsNew, gsSource)
	writeGenesisNew(os.Args[4], &gsNew)
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

type Statistic struct {
	countRegularAccounts            uint64
	countRegularAccountsNoPublicKey uint64
}

func convertGenesis(gsOld *GenesisOld, fixNFTData []NFTOwnerFixRecord) (GenesisNew, Statistic, error) {
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
	// coins
	gsNew.AppState.Coin.Coins, err = convertCoins(gsOld.AppState.Coin.Coins, addrTable)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	// balances
	legacyRecords := NewLegacyRecords()
	gsNew.AppState.Bank.Balances, err =
		convertBalances(gsOld.AppState.Auth.Accounts, addrTable, legacyRecords, gsNew.AppState.Coin.Coins)
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
		convertMultisigTransactions(gsOld.AppState.Multisig.Transactions, addrTable, gsOld.AppState.Multisig.Wallets, gsNew.AppState.Coin.Coins)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	// nft
	gsNew.AppState.NFT.Collections, err =
		convertNFT(gsOld.AppState.NFT.Collections, gsOld.AppState.NFT.SubTokens, addrTable, legacyRecords, fixNFTData)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
	// inject to legacy records
	// TODO: parametrize output?
	if true {
		fmt.Printf("!!! LEGACY RECORDS INJECTING !!! REMOVE FROM PRODUCTION\n")
		out, err := os.Create("legacy_test_mnemonics.txt")
		if err != nil {
			return GenesisNew{}, Statistic{}, err
		}
		coins := sdk.NewCoins(sdk.NewCoin("del", sdkmath.NewInt(3_141_592_653_589_793_238)))
		for i := 0; i < 40000; i++ {
			mn, _ := dscWallet.NewMnemonic("")
			acc, _ := dscWallet.NewAccountFromMnemonicWords(mn.Words(), "")
			legacy_address, _ := dscTypes.GetLegacyAddressFromPubKey(acc.PubKey().Bytes())
			address := acc.Address()
			out.WriteString(fmt.Sprintf("%s\t%s\t%s\n", legacy_address, address, mn.Words()))
			legacyRecords.AddCoins(legacy_address, coins)
		}
		out.Close()
	}
	// legacy records
	var records []LegacyRecordNew
	for _, v := range legacyRecords.data {
		if v.Coins.Len() == 0 && len(v.NFTs) == 0 && len(v.Wallets) == 0 {
			continue
		}
		records = append(records, *v)
	}
	gsNew.AppState.Legacy.LegacyRecords = records
	// validators
	gsNew.AppState.Validator.Validators, err =
		convertValidators(gsOld.AppState.Validator.Validators, addrTable)
	if err != nil {
		return GenesisNew{}, Statistic{}, err
	}
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
	verifyNFTSupply(gsNew.AppState.NFT.Collections)

	return gsNew, Statistic{}, nil
}
