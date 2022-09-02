package main

import (
	"fmt"
	"sort"
	"strconv"

	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/strings"
)

func prepareAddressTable(gs *GenesisOld) (*AddressTable, error) {
	var table = NewAddressTable()
	var err error
	var pubkey []byte
	for _, acc := range gs.AppState.Auth.Accounts {
		if acc.Typ == accountTypeModule {
			continue
		}
		if acc.Value.Address == "" {
			continue
		}
		if acc.Value.PublicKey == nil {
			err = table.AddAddress(acc.Value.Address, []byte{})
		} else {
			pubkey, err = extractPubKey(acc.Value.PublicKey)
			if err != nil {
				return nil, err
			}
			err = table.AddAddress(acc.Value.Address, pubkey)
		}
		if err != nil {
			return nil, err
		}
	}
	for _, wallet := range gs.AppState.Multisig.Wallets {
		table.AddMultisig(wallet.Address)
	}
	return &table, nil
}

func convertAccounts(accsOld []AccountOld, addrTable *AddressTable) ([]interface{}, error) {
	var res []interface{}
	var err error
	for _, acc := range accsOld {
		var accNew interface{}
		if acc.Typ == accountTypeRegular {
			if acc.Value.Address == "" {
				continue
			}
			if acc.Value.PublicKey == nil {
				continue
			}
			accNew, err = AccountO2N(acc)
			if err != nil {
				return []interface{}{}, fmt.Errorf("account %s, error %s", acc.Value.Address, err.Error())
			}
		}
		if acc.Typ == accountTypeModule {
			accNew, err = ModuleAccountO2N(acc, addrTable)
			if err != nil {
				return []interface{}{}, fmt.Errorf("account %s, error %s", acc.Value.Address, err.Error())
			}
		}
		res = append(res, accNew)
	}
	return res, nil
}

// convert tDEL to DEL
func filterCoins(coins sdk.Coins, coinSymbols map[string]bool) sdk.Coins {
	var result = sdk.NewCoins()
	for _, coin := range coins {
		if !coinSymbols[coin.Denom] {
			continue
		}
		if coin.Denom == "tdel" {
			result = result.Add(sdk.NewCoin("del", coin.Amount))
		} else {
			result = result.Add(coin)
		}
	}
	return result
}

func convertBalances(accsOld []AccountOld, addrTable *AddressTable, legacyRecords *LegacyRecords, coins []FullCoinNew) ([]BalanceNew, error) {
	// coin symbol cache to skip unexisting coins
	var coinSymbols = make(map[string]bool)
	for _, c := range coins {
		coinSymbols[c.Symbol] = true
	}

	var res []BalanceNew
	var legacyBalance = sdk.NewCoins()
	for _, acc := range accsOld {
		if acc.Value.Address == "" {
			continue
		}
		if len(acc.Value.Coins) == 0 {
			continue
		}
		newAddress := addrTable.GetAddress(acc.Value.Address)
		if addrTable.IsMultisig(acc.Value.Address) {
			newAddress = acc.Value.Address
		}
		if acc.Typ == accountTypeModule {
			newAddress = addrTable.GetModule(acc.Value.Name).address
			if newAddress == "" {
				return []BalanceNew{}, fmt.Errorf("address %s: unknown module name '%s'", acc.Value.Address, acc.Value.Name)
			}
		}

		coins := filterCoins(acc.Value.Coins, coinSymbols)
		// TODO: return when correct staking starts work
		if acc.Value.Name == "not_bonded_tokens_pool" || acc.Value.Name == "bonded_tokens_pool" {
			fmt.Printf("set '%s' module account balance to zero\n", acc.Value.Name)
			coins = sdk.NewCoins()
		}

		if newAddress > "" {
			res = append(res, BalanceNew{Address: newAddress, Coins: coins})
		} else {
			// empty address: no multisig, no module
			legacyBalance = legacyBalance.Add(coins...)
			legacyRecords.AddCoins(acc.Value.Address, coins)
		}
	}
	// legacy_coin_pool
	res = append(res, BalanceNew{Address: addrTable.GetModule("legacy_coin_pool").address, Coins: legacyBalance})
	return res, nil
}

func validCoinLimits(coin FullCoinOld) bool {
	var result = true
	// volume
	v, _ := sdk.NewIntFromString(coin.Volume)
	if v.LT(coinTypes.MinCoinSupply) {
		fmt.Printf("coin %s: volume < MinCoinSupply\n", coin.Symbol)
		result = false
	}
	if v.GT(coinTypes.MaxCoinSupply) {
		fmt.Printf("coin %s: volume > MaxCoinSupply\n", coin.Symbol)
		result = false
	}
	// limit volume
	v, _ = sdk.NewIntFromString(coin.LimitVolume)
	if v.LT(coinTypes.MinCoinSupply) {
		fmt.Printf("coin %s: limit_volume < MinCoinSupply\n", coin.Symbol)
		result = false
	}
	if v.GT(coinTypes.MaxCoinSupply) {
		fmt.Printf("coin %s: limit_volume > MaxCoinSupply\n", coin.Symbol)
		result = false
	}
	// reserve
	v, _ = sdk.NewIntFromString(coin.Reserve)
	if v.LT(coinTypes.MinCoinReserve) {
		fmt.Printf("coin %s: reserve < MinCoinReserve\n", coin.Symbol)
		result = false
	}
	return result
}

func convertCoins(coinsOld []FullCoinOld, addrTable *AddressTable) ([]FullCoinNew, error) {
	var res []FullCoinNew
	for _, coin := range coinsOld {
		if coin.Symbol != "tdel" && coin.Symbol != "del" && !validCoinLimits(coin) {
			continue
		}
		res = append(res, FullCoinO2N(coin, addrTable))
	}
	return res, nil
}

func convertMultisigWallets(walletsOld []WalletOld, addrTable *AddressTable, legacyRecords *LegacyRecords) ([]WalletNew, error) {
	var res []WalletNew
	for _, wallet := range walletsOld {
		newWallet := WalletO2N(wallet, addrTable, legacyRecords)
		res = append(res, newWallet)
	}
	return res, nil
}

func isTxIncomplete(tx TransactionOld, wallet WalletOld) bool {
	threshold, _ := strconv.ParseUint(wallet.Threshold, 10, 64)
	wsum := uint64(0)
	for i := range tx.Signers {
		if tx.Signers[i] > "" {
			w, _ := strconv.ParseUint(wallet.Weights[i], 10, 64)
			wsum += w
		}
	}
	return wsum >= threshold
}

func convertMultisigTransactions(transactionsOld []TransactionOld, addrTable *AddressTable, walletsOld []WalletOld, coins []FullCoinNew) ([]TransactionNew, error) {
	// coin symbol cache to skip unexisting coins
	var coinSymbols = make(map[string]bool)
	for _, c := range coins {
		coinSymbols[c.Symbol] = true
	}
	var res []TransactionNew
	for _, txOld := range transactionsOld {
		if addrTable.GetAddress(txOld.Receiver) == "" {
			continue
		}
		wallet := WalletOld{}
		for _, w := range walletsOld {
			if txOld.Wallet == w.Address {
				wallet = w
				break
			}
		}
		if isTxIncomplete(txOld, wallet) {
			continue
		}
		newTx := TransactionO2N(txOld, addrTable, coinSymbols)
		if newTx.Coins.Empty() {
			fmt.Printf("skip multisig transaction %s - empty coins\n", txOld.ID)
			continue
		}
		res = append(res, newTx)
	}
	return res, nil
}

func convertNFT(collectionsOld map[string]CollectionOld, addrTable *AddressTable, legacyRecords *LegacyRecords) ([]CollectionNew, []NFTNew, error) {
	var collectionsNew []CollectionNew
	var nftsNew []NFTNew
	for _, colOld := range collectionsOld {
		colNew := CollectionNew{Denom: colOld.Denom}
		for _, nftOld := range colOld.NFT {
			creatorAddress := addrTable.GetAddress(nftOld.Creator)
			if creatorAddress == "" {
				return []CollectionNew{}, []NFTNew{}, fmt.Errorf("unknown creator %s for nft %s", nftOld.Creator, nftOld.ID)
			}
			reserve, _ := sdk.NewIntFromString(nftOld.Reserve)
			nftNew := NFTNew{
				ID:        nftOld.ID,
				AllowMint: nftOld.AllowMint,
				Creator:   creatorAddress,
				Reserve:   sdk.NewCoin("del", reserve),
				TokenURI:  nftOld.TokenURI,
			}
			owners := []OwnerNew{}
			for _, ownerOld := range nftOld.Owners["owners"] {
				if len(ownerOld.SubTokenIds) == 0 {
					continue
				}
				subs := make([]string, 0, len(ownerOld.SubTokenIds))
				for _, s := range ownerOld.SubTokenIds {
					// check subtoken id already owned
					owned := false
					subID := strconv.FormatUint(s, 10)
					for _, o := range owners {
						if strings.StringInSlice(subID, o.SubTokenIDs) {
							owned = true
							break
						}
					}
					if owned {
						fmt.Printf("ntf: %s, sub: %s already owned\n", nftOld.ID, subID)
					} else {
						subs = append(subs, subID)
					}
				}
				ownerAddress := addrTable.GetAddress(ownerOld.Address)
				if ownerAddress == "" {
					legacyRecords.AddNFT(ownerOld.Address, colOld.Denom, nftOld.ID)
					owners = append(owners, OwnerNew{Address: ownerOld.Address, SubTokenIDs: subs})
				} else {
					owners = append(owners, OwnerNew{Address: ownerAddress, SubTokenIDs: subs})
				}
			}
			if len(owners) == 0 {
				fmt.Printf("nft without owners: nft_id %s , collection %s\n", nftNew.ID, colNew.Denom)
				continue
			}
			nftNew.Owners = owners
			nftsNew = append(nftsNew, nftNew)
			colNew.NFTs = append(colNew.NFTs, nftNew.ID)
		}
		sort.Slice(colNew.NFTs, func(i, j int) bool {
			return colNew.NFTs[i] < colNew.NFTs[j]
		})
		collectionsNew = append(collectionsNew, colNew)
	}
	return collectionsNew, nftsNew, nil
}

func convertSubTokens(subsOld []SubTokenOld, nfts []NFTNew) (map[string]SubTokensNew, error) {
	// prepare existsing subtokens
	var existingSubs = make(map[string][]string)
	for _, nft := range nfts {
		for _, owner := range nft.Owners {
			existingSubs[nft.ID] = append(existingSubs[nft.ID], owner.SubTokenIDs...)
		}
	}
	var subsNew = make(map[string]SubTokensNew)
	for _, sub := range subsOld {
		newSub := subsNew[sub.NftID]
		reserve, _ := sdk.NewIntFromString(sub.Reserve)
		if !strings.StringInSlice(sub.ID, existingSubs[sub.NftID]) {
			fmt.Printf("skip (not owned) nft: %s, subtoken: %s\n", sub.NftID, sub.ID)
			continue
		}
		newSub.SubTokens = append(newSub.SubTokens, SubTokenNew{ID: sub.ID, Reserve: sdk.NewCoin("del", reserve)})
		subsNew[sub.NftID] = newSub
	}
	return subsNew, nil
}
