package main

import (
	"fmt"
	"strconv"

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
	delegatedNFT []DelegationNFT, unbondingNFT []UnbondingNFTRecord) map[SubTokenRecord]*SubTokenCount {
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
	Volume sdk.Int
	BCSum  sdk.Int
}

func verifyCoinsVolume(coinsOld []FullCoinOld, accsOld []AccountOld,
	delegations []DelegationOld, unbondings []UnbondingRecord) []CoinDiff {
	fullSum := sdk.NewCoins()
	for _, acc := range accsOld {
		if acc.Value.Name == "bonded_tokens_pool" || acc.Value.Name == "not_bonded_tokens_pool" {
			continue
		}
		fullSum = fullSum.Add(acc.Value.Coins...)
	}
	for _, del := range delegations {
		fullSum = fullSum.Add(del.Coin)
	}
	for _, rec := range unbondings {
		for _, ent := range rec.Entries {
			fullSum = fullSum.Add(ent.Value.Balance)
		}
	}
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

func verifyNFTSupply(collections []CollectionNew, nfts []NFTNew) {
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

}
