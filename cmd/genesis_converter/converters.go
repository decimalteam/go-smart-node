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

func convertBalances(accsOld []AccountOld, addrTable *AddressTable) ([]BalanceNew, []LegacyBalanceNew, error) {
	var res []BalanceNew
	var legacyBalance = sdk.NewCoins()
	var legacyRecords []LegacyBalanceNew
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
				return []BalanceNew{}, []LegacyBalanceNew{}, fmt.Errorf("address %s: unknown module name '%s'", acc.Value.Address, acc.Value.Name)
			}
		}
		coins := sdk.NewCoins()
		for _, c := range acc.Value.Coins {
			amount, ok := sdk.NewIntFromString(c.Amount)
			if !ok {
				return []BalanceNew{}, []LegacyBalanceNew{}, fmt.Errorf("address %s: cannot convert '%s' to sdk.Int", acc.Value.Address, c.Amount)
			}
			coins = coins.Add(sdk.NewCoin(c.Denom, amount))
		}
		if newAddress > "" {
			res = append(res, BalanceNew{Address: newAddress, Coins: coins})
		} else {
			// empty address: no multisig, no module
			legacyBalance = legacyBalance.Add(coins...)
			legacyRecords = append(legacyRecords, LegacyBalanceNew{Address: acc.Value.Address, Coins: coins})
		}
	}
	// legacy_coin_pool
	res = append(res, BalanceNew{Address: addrTable.GetModule("legacy_coin_pool").address, Coins: legacyBalance})
	return res, legacyRecords, nil
}

func convertCoins(coinsOld []FullCoinOld, addrTable *AddressTable) ([]FullCoinNew, error) {
	var res []FullCoinNew
	for _, coin := range coinsOld {
		res = append(res, FullCoinO2N(coin, addrTable))
	}
	return res, nil
}

func convertMultisig(walletsOld []WalletOld, addrTable *AddressTable) ([]WalletNew, error) {
	var res []WalletNew
	for _, wallet := range walletsOld {
		newWallet := WalletO2N(wallet, addrTable)
		res = append(res, newWallet)
	}
	return res, nil
}

func convertNFT(collectionsOld map[string]CollectionOld, addrTable *AddressTable) ([]CollectionNew, []NFTNew, error) {
	var collectionsNew []CollectionNew
	var nftsNew []NFTNew
	unknownOwners := 0
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
					unknownOwners++
					//	return []CollectionNew{}, []NFTNew{}, fmt.Errorf("unknown owner %s for nft %s", ownerOld.Address, nftOld.ID)
				}
				owners = append(owners, OwnerNew{Address: ownerAddress, SubTokenIDs: subs})
			}
			nftNew.Owners = owners
			nftsNew = append(nftsNew, nftNew)
			colNew.NFTs = append(colNew.NFTs, nftNew.ID)
		}
		collectionsNew = append(collectionsNew, colNew)
	}
	fmt.Printf("unknownOwners=%d\n", unknownOwners)
	return collectionsNew, nftsNew, nil
}
