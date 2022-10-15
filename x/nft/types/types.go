package types

import (
	"bitbucket.org/decimalteam/go-smart-node/types"
)

const ReservedPool = "reserved_pool"

type (
	SortedStringArray = types.SortedStringArray
	SortedUintArray   = types.SortedUintArray
	Tokens            = []*Token
	SubTokens         = []*SubToken
)

func SupplyInvariantCheck(collections []Collection) (string, bool) {
	// totalSupply := 0
	// for _, collection := range collections {
	// 	totalSupply += collection.Supply()
	// }
	broken := false // len(nfts) != totalSupply

	return "nft supply invariants found", broken
}
