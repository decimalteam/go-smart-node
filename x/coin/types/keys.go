package types

import "strings"

const (
	// ModuleName defines the module name
	ModuleName = "coin"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	// RouterKey defines module's messages routing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore prefix bytes
const (
	prefixCoin  byte = iota + 1 // 1
	prefixCheck                 // 2
	prefixLegacy
)

// KVStore key prefixes
var (
	KeyPrefixCoin   = []byte{prefixCoin}   // 0x01
	KeyPrefixCheck  = []byte{prefixCheck}  // 0x02
	KeyPrefixLegacy = []byte{prefixLegacy} // 0x03
)

func CoinStoreKey(symbol string) []byte {
	return append(KeyPrefixCoin, []byte(strings.ToLower(symbol))...)
}

func CheckStoreKey(checkHash []byte) []byte {
	return append(KeyPrefixCheck, checkHash...)
}

func LegacyBalanceStoreKey(legacyAddress string) []byte {
	return append(KeyPrefixLegacy, []byte(legacyAddress)...)
}
