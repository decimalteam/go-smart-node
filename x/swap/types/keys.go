package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "swap"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	// RouterKey defines module's messages routing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore key prefixes
var (
	KeyPrefixSwap  = []byte("swap")
	KeyPrefixChain = []byte("chain")
)

func GetSwapKey(hash Hash) []byte {
	return append(KeyPrefixSwap, hash[:]...)
}

func GetChainKey(chain uint32) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(chain))
	return append(KeyPrefixChain, buf...)
}
