package main

import (
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SubTokenRecord struct {
	Denom    string
	NFT      string
	SubToken string
}

type SubTokenCount struct {
	Subs      int
	Owned     int
	Delegated int
	Unbonding int
}

func verifySubtokens(subsOld []SubTokenOld, collectionsOld map[string]CollectionOld,
	delegatedNFT []DelegationNFTOld, unbondingNFT []UnbondingNFTRecordOld) map[SubTokenRecord]*SubTokenCount {
	var checkMap = make(map[SubTokenRecord]*SubTokenCount)
	for _, sub := range subsOld {
		checkMap[SubTokenRecord{Denom: sub.Denom, NFT: sub.NftID, SubToken: sub.ID}] = &SubTokenCount{Subs: 1}
	}
	for _, colOld := range collectionsOld {
		for _, nftOld := range colOld.NFT {
			for _, ownerOld := range nftOld.Owners["owners"] {
				if len(ownerOld.SubTokenIds) == 0 {
					continue
				}
				for _, s := range ownerOld.SubTokenIds {
					subID := strconv.FormatUint(s, 10)
					key := SubTokenRecord{Denom: colOld.Denom, NFT: nftOld.ID, SubToken: subID}
					cnt, ok := checkMap[key]
					if !ok {
						cnt = &SubTokenCount{}
					}
					cnt.Owned++
					checkMap[key] = cnt
				}
			}
		}
	}
	for _, del := range delegatedNFT {
		for _, subID := range del.SubTokenIds {
			key := SubTokenRecord{Denom: del.Denom, NFT: del.TokenID, SubToken: subID}
			cnt, ok := checkMap[key]
			if !ok {
				cnt = &SubTokenCount{}
			}
			cnt.Delegated++
			checkMap[key] = cnt
		}
	}
	for _, rec := range unbondingNFT {
		for _, ent := range rec.Entries {
			for _, subID := range ent.SubTokenIds {
				key := SubTokenRecord{Denom: ent.Denom, NFT: ent.TokenID, SubToken: subID}
				cnt, ok := checkMap[key]
				if !ok {
					cnt = &SubTokenCount{}
				}
				cnt.Unbonding++
				checkMap[key] = cnt
			}
		}
	}
	invalidSubtokens := make(map[SubTokenRecord]*SubTokenCount)
	for key, cnt := range checkMap {
		if cnt.Subs == 1 && (cnt.Delegated+cnt.Owned+cnt.Unbonding) == 1 {
			continue
		}
		invalidSubtokens[key] = cnt
	}
	return invalidSubtokens
}

type CoinDiff struct {
	Symbol string
	Volume math.Int
	BCSum  math.Int
}

func verifyCoinsVolume(coinsOld []FullCoinOld, accsOld []AccountOld,
	delegations []DelegationOld, unbondings []UnbondingRecordOld) []CoinDiff {
	fullSum := sdk.NewCoins()
	for _, acc := range accsOld {
		//if acc.Value.Name == "bonded_tokens_pool" || acc.Value.Name == "not_bonded_tokens_pool" {
		//	continue
		//}
		fullSum = fullSum.Add(acc.Value.Coins...)
	}
	//for _, del := range delegations {
	//	fullSum = fullSum.Add(del.Coin)
	//}
	//for _, rec := range unbondings {
	//	for _, ent := range rec.Entries {
	//		fullSum = fullSum.Add(ent.Value.Balance)
	//	}
	//}
	var result []CoinDiff
	for _, coin := range coinsOld {
		vol, _ := sdk.NewIntFromString(coin.Volume)
		diff := CoinDiff{Symbol: coin.Symbol, Volume: vol}
		for _, c := range fullSum {
			if c.Denom == diff.Symbol {
				diff.BCSum = c.Amount
				break
			}
		}
		result = append(result, diff)
	}
	return result
}

func verifyNFTSupply(collections []CollectionNew) {
	/*
		var nftDenoms = make(map[string]string)
		var collSupply int
		for _, coll := range collections {
			for _, nft := range coll.NFTs {
				_, ok := nftDenoms[nft]
				if ok {
					panic(fmt.Sprintf("duplicate nft %s", nft))
				}
				nftDenoms[nft] = coll.Denom
			}
			collSupply += len(coll.NFTs)
		}
		for _, nft := range nfts {
			_, ok := nftDenoms[nft.ID]
			if !ok {
				fmt.Printf("nft %s not found in collection\n", nft.ID)
			}
		}

		fmt.Printf("??? coll supply != len(ntfs): %d != %d\n", collSupply, len(nfts))
	*/
}

func verifyPools(balances []BalanceNew, validators []ValidatorNew, delegations []DelegationNew,
	undelegations []UndelegationNew, addrTable *AddressTable) {
	valStatuses := make(map[string]string)
	for _, val := range validators {
		valStatuses[val.OperatorAddress] = val.Status
	}
	bondedCoins := sdk.NewCoins()
	notBondedCoins := sdk.NewCoins()
	for _, del := range delegations {
		if del.Stake.Type == "STAKE_TYPE_COIN" {
			switch valStatuses[del.Validator] {
			case "BOND_STATUS_BONDED":
				bondedCoins = bondedCoins.Add(del.Stake.Stake)
			case "BOND_STATUS_UNBONDED", "BOND_STATUS_UNBONDING":
				notBondedCoins = notBondedCoins.Add(del.Stake.Stake)
			}
		}
	}
	for _, ubd := range undelegations {
		for _, entry := range ubd.Entries {
			if entry.Stake.Type == "STAKE_TYPE_COIN" {
				notBondedCoins = notBondedCoins.Add(entry.Stake.Stake)
			}
		}
	}
	bpool := addrTable.GetModule("bonded_tokens_pool").address
	nbpool := addrTable.GetModule("not_bonded_tokens_pool").address
	fmt.Printf("bonded_tokens_pool address = %s\n", bpool)
	fmt.Printf("not_bonded_tokens_pool address = %s\n", nbpool)
	for _, bal := range balances {
		if bal.Address == bpool {
			if !bal.Coins.IsEqual(bondedCoins) {
				denoms := make(map[string]bool)
				for _, c := range bal.Coins {
					denoms[c.Denom] = true
				}
				for _, c := range bondedCoins {
					denoms[c.Denom] = true
				}
				for denom := range denoms {
					b1 := bal.Coins.AmountOf(denom)
					b2 := bondedCoins.AmountOf(denom)
					if !b1.Equal(b2) {
						fmt.Printf("different bonded pool (module account <-> stakes): (%s) %s <-> %s\n", denom, b1, b2)
					}
				}
			}
		}
		if bal.Address == nbpool {
			if !bal.Coins.IsEqual(notBondedCoins) {
				denoms := make(map[string]bool)
				for _, c := range bal.Coins {
					denoms[c.Denom] = true
				}
				for _, c := range notBondedCoins {
					denoms[c.Denom] = true
				}
				for denom := range denoms {
					b1 := bal.Coins.AmountOf(denom)
					b2 := notBondedCoins.AmountOf(denom)
					if !b1.Equal(b2) {
						fmt.Printf("different not bonded pool (module account <-> stakes): (%s) %s <-> %s\n", denom, b1, b2)
					}
				}
			}
		}
	}
}

// all NFT delegations/undelegations/redelegations must have record in NFT genesis
// and must have proper owner
func verifyNFTDelegations(gs *GenesisNew, addrTable *AddressTable) {
	type nftKey struct {
		tokenID  string
		subtoken uint32
	}
	bpool := addrTable.GetModule("bonded_tokens_pool").address
	nbpool := addrTable.GetModule("not_bonded_tokens_pool").address

	nftRecords := make(map[nftKey]string)
	delegationRecords := make(map[nftKey]string)
	delegationCounts := make(map[nftKey]int)
	for _, coll := range gs.AppState.NFT.Collections {
		for _, token := range coll.Tokens {
			for _, sub := range token.SubTokens {
				nftRecords[nftKey{token.ID, sub.ID}] = sub.Owner
			}
		}
	}
	validatorsStatus := make(map[string]bool) // true - online
	for _, val := range gs.AppState.Validator.Validators {
		validatorsStatus[val.OperatorAddress] = val.Online
	}

	for _, del := range gs.AppState.Validator.Delegations {
		if del.Stake.Type == "STAKE_TYPE_NFT" {
			owner := nbpool
			if validatorsStatus[del.Validator] {
				owner = bpool
			}
			for _, subID := range del.Stake.SubTokenIDs {
				delegationRecords[nftKey{del.Stake.ID, uint32(subID)}] = owner
				delegationCounts[nftKey{del.Stake.ID, uint32(subID)}] = delegationCounts[nftKey{del.Stake.ID, uint32(subID)}] + 1
			}
		}
	}
	for _, undel := range gs.AppState.Validator.Undelegations {
		for _, entry := range undel.Entries {
			if entry.Stake.Type == "STAKE_TYPE_NFT" {
				owner := nbpool
				for _, subID := range entry.Stake.SubTokenIDs {
					delegationRecords[nftKey{entry.Stake.ID, uint32(subID)}] = owner
					delegationCounts[nftKey{entry.Stake.ID, uint32(subID)}] = delegationCounts[nftKey{entry.Stake.ID, uint32(subID)}] + 1
				}
			}
		}
	}
	for k, v := range delegationCounts {
		if v > 1 {
			fmt.Printf("NFT DELEGATIONS: too much delegations (%s, %d) = %d\n", k.tokenID, k.subtoken, v)
		}
	}
	for k, owner := range delegationRecords {
		nftowner, ok := nftRecords[k]
		if !ok {
			fmt.Printf("(un)delegation (%s, %03d) not found\n", k.tokenID, k.subtoken)
			continue
		}
		if owner != nftowner {
			fmt.Printf("(un)delegation (%s, %03d) owner (nft) %s != (del) %s\n", k.tokenID, k.subtoken, nftowner, owner)
		}
	}
}
