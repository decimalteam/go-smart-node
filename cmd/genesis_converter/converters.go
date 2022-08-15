package main

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func convertBalances(accsOld []AccountOld, addrTable *AddressTable, legacyRecords *LegacyRecords) ([]BalanceNew, error) {
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
		coins := acc.Value.Coins
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

func convertCoins(coinsOld []FullCoinOld, addrTable *AddressTable) ([]FullCoinNew, error) {
	var res []FullCoinNew
	for _, coin := range coinsOld {
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

func convertMultisigTransactions(transactionsOld []TransactionOld, addrTable *AddressTable, walletsOld []WalletOld) ([]TransactionNew, error) {
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
		newTx := TransactionO2N(txOld, addrTable)
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
			nftNew := NFTNew{
				ID:        nftOld.ID,
				AllowMint: nftOld.AllowMint,
				Creator:   creatorAddress,
				Reserve:   nftOld.Reserve,
				TokenURI:  nftOld.TokenURI,
			}
			owners := []OwnerNew{}
			for _, ownerOld := range nftOld.Owners["owners"] {
				if len(ownerOld.SubTokenIds) == 0 {
					continue
				}
				subs := make([]string, len(ownerOld.SubTokenIds))
				for i, s := range ownerOld.SubTokenIds {
					subs[i] = strconv.FormatUint(s, 10)
				}
				ownerAddress := addrTable.GetAddress(ownerOld.Address)
				if ownerAddress == "" {
					legacyRecords.AddNFT(ownerOld.Address, colOld.Denom, nftOld.ID)
					owners = append(owners, OwnerNew{Address: ownerOld.Address, SubTokenIDs: subs})
				} else {
					owners = append(owners, OwnerNew{Address: ownerAddress, SubTokenIDs: subs})
				}
			}
			nftNew.Owners = owners
			nftsNew = append(nftsNew, nftNew)
			colNew.NFTs = append(colNew.NFTs, nftNew.ID)
		}
		collectionsNew = append(collectionsNew, colNew)
	}
	return collectionsNew, nftsNew, nil
}

func convertSubTokens(subsOld []SubTokenOld) (map[string]SubTokensNew, error) {
	var subsNew = make(map[string]SubTokensNew)
	for _, sub := range subsOld {
		newSub := subsNew[sub.NftID]
		newSub.SubTokens = append(newSub.SubTokens, SubTokenNew{ID: sub.ID, Reserve: sub.Reserve})
		subsNew[sub.NftID] = newSub
	}
	return subsNew, nil
}
