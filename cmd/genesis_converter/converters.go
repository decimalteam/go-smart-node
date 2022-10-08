package main

import (
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	coinconfig "bitbucket.org/decimalteam/go-smart-node/x/coin/config"
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
		if coin.Denom == "tdel" {
			result = result.Add(sdk.NewCoin("del", coin.Amount))
		} else {
			if !coinSymbols[coin.Denom] {
				continue
			}
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

func validCoinParams(coin FullCoinOld) bool {
	var result = true
	// crr
	crr, err := strconv.ParseUint(coin.CRR, 10, 64)
	if err != nil {
		fmt.Printf("coin %s: crr '%s' error: %s\n", coin.Symbol, coin.CRR, err)
		result = false
	} else {
		if crr < 10 || crr > 100 {
			fmt.Printf("coin %s: invalid crr: %d\n", coin.Symbol, crr)
			result = false
		}
	}
	// volume
	v, _ := sdk.NewIntFromString(coin.Volume)
	if v.LT(coinconfig.MinCoinSupply) {
		fmt.Printf("coin %s: volume < MinCoinSupply\n", coin.Symbol)
		result = false
	}
	if v.GT(coinconfig.MaxCoinSupply) {
		fmt.Printf("coin %s: volume > MaxCoinSupply\n", coin.Symbol)
		result = false
	}
	// limit volume
	v, _ = sdk.NewIntFromString(coin.LimitVolume)
	if v.LT(coinconfig.MinCoinSupply) {
		fmt.Printf("coin %s: limit_volume < MinCoinSupply\n", coin.Symbol)
		result = false
	}
	if v.GT(coinconfig.MaxCoinSupply) {
		fmt.Printf("coin %s: limit_volume > MaxCoinSupply\n", coin.Symbol)
		result = false
	}
	// reserve
	v, _ = sdk.NewIntFromString(coin.Reserve)
	if v.LT(coinconfig.MinCoinReserve) {
		fmt.Printf("coin %s: reserve < MinCoinReserve\n", coin.Symbol)
		result = false
	}
	return result
}

func convertCoins(coinsOld []FullCoinOld, addrTable *AddressTable) ([]FullCoinNew, error) {
	var res []FullCoinNew
	for _, coin := range coinsOld {
		if coin.Symbol != "tdel" && coin.Symbol != "del" && !validCoinParams(coin) {
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

func convertNFT(collectionsOld map[string]CollectionOld, subsOld []SubTokenOld,
	addrTable *AddressTable, legacyRecords *LegacyRecords, fixNFTData []NFTOwnerFixRecord) ([]CollectionNew, error) {
	// prepare subtokens
	type subRecord struct {
		id      string
		reserve math.Int
	}
	preparedSubTokens := make(map[string][]subRecord)
	for _, sub := range subsOld {
		prep := preparedSubTokens[sub.NftID]
		reserve, ok := sdk.NewIntFromString(sub.Reserve)
		if !ok {
			return []CollectionNew{}, fmt.Errorf("cant parse reserve for subtoken %s nft %s", sub.ID, sub.NftID)
		}
		prep = append(prep, subRecord{id: sub.ID, reserve: reserve})
		preparedSubTokens[sub.NftID] = prep
	}
	// prepare collections
	type collectionKey struct {
		denom   string
		creator string
	}
	// URI uniqueness
	tokenURIs := make(map[string]bool)

	preparedColls := make(map[collectionKey]*CollectionNew)
	for _, colOld := range collectionsOld {
		for _, nftOld := range colOld.NFT {
			creatorAddress := addrTable.GetAddress(nftOld.Creator)
			if creatorAddress == "" {
				return []CollectionNew{}, fmt.Errorf("unknown creator (no pubkey) %s for nft %s", nftOld.Creator, nftOld.ID)
			}
			key := collectionKey{denom: colOld.Denom, creator: creatorAddress}
			if _, ok := preparedColls[key]; !ok {
				preparedColls[key] = &CollectionNew{Denom: colOld.Denom, Creator: creatorAddress}
			}
		}
	}
	for _, colOld := range collectionsOld {
		for _, nftOld := range colOld.NFT {
			// check URI uniq
			if tokenURIs[nftOld.TokenURI] {
				fmt.Printf("found yet another token URI: %s\n", nftOld.TokenURI)
				continue
			}
			tokenURIs[nftOld.TokenURI] = true

			creatorAddress := addrTable.GetAddress(nftOld.Creator)
			key := collectionKey{denom: colOld.Denom, creator: creatorAddress}
			collNew := preparedColls[key]
			// 2. subtokens
			subtokens := make([]SubTokenNew, 0)
			for _, sub := range preparedSubTokens[nftOld.ID] {
				id, err := strconv.ParseUint(sub.id, 10, 32)
				if err != nil {
					return []CollectionNew{}, fmt.Errorf("can't parse for nft '%s' subtoken id : %s", nftOld.ID, sub.id)
				}
				subtokens = append(subtokens, SubTokenNew{
					ID:      uint32(id),
					Owner:   "",
					Reserve: sdk.NewCoin("del", sub.reserve),
				})
			}
			// 3. owners for subtokens
			for _, ownerOld := range nftOld.Owners["owners"] {
				if len(ownerOld.SubTokenIds) == 0 {
					continue
				}
				ownerAddress := addrTable.GetAddress(ownerOld.Address)
				if ownerAddress == "" {
					legacyRecords.AddNFT(ownerOld.Address, colOld.Denom, nftOld.ID)
					ownerAddress = ownerOld.Address
				}
				for i := range subtokens {
					for _, id := range ownerOld.SubTokenIds {
						if id == uint64(subtokens[i].ID) {
							subtokens[i].Owner = ownerAddress
						}
					}
				}
			}
			// 3.5 fix owners
			for _, rec := range fixNFTData {
				if rec.TokenID != nftOld.ID {
					continue
				}
				ownerAddress := addrTable.GetAddress(rec.Owner)
				if ownerAddress == "" {
					return []CollectionNew{}, fmt.Errorf("impossible situation: lost nft for owner '%s'", rec.Owner)
				}
				for i := range subtokens {
					for _, id := range rec.SubTokens {
						if id == subtokens[i].ID {
							subtokens[i].Owner = ownerAddress
						}
					}
				}
			}
			// 3.9 TODO: empty owners for subtokens in testnet. Workaround with logging
			// NOTE: bech32 address for []byte{0} = "dx1qqjrdrw8",
			for i := range subtokens {
				if subtokens[i].Owner == "" {
					fmt.Printf("empty owner for collection '%s', creator '%s', nft '%s', sub token id '%d'\n",
						collNew.Denom, collNew.Creator, nftOld.ID, subtokens[i].ID)
					subtokens[i].Owner = "dx1qqjrdrw8"
				}
			}
			// 4. build nft and add to collection
			initialReserve, ok := sdk.NewIntFromString(nftOld.Reserve)
			if !ok {
				return []CollectionNew{}, fmt.Errorf("can't parse initial reserve for nft %s", nftOld.ID)
			}
			nftNew := TokenNew{
				Creator:   creatorAddress,
				Denom:     colOld.Denom,
				ID:        nftOld.ID,
				URI:       nftOld.TokenURI,
				Reserve:   sdk.NewCoin("del", initialReserve),
				AllowMint: nftOld.AllowMint,
				Minted:    uint32(len(subtokens)),
				Burnt:     0,
				SubTokens: subtokens,
			}
			// add to collection
			collNew.Supply++
			collNew.Tokens = append(collNew.Tokens, nftNew)
			preparedColls[key] = collNew
		}
	}
	var collectionsNew []CollectionNew
	for _, collNew := range preparedColls {
		collectionsNew = append(collectionsNew, *collNew)
	}
	return collectionsNew, nil
}

func convertValidators(valsOld []ValidatorOld, addrTable *AddressTable) ([]ValidatorNew, error) {
	var result []ValidatorNew
	for _, valOld := range valsOld {
		valNew, err := ValidatorO2N(valOld, addrTable)
		if err != nil {
			return []ValidatorNew{}, err
		}
		result = append(result, valNew)
	}
	return result, nil
}
