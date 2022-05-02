package types

const (
	// ModuleName defines the module name
	ModuleName = "coin"

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_coin"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines module's messages routing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore prefix bytes
const (
	prefixCoin byte = iota + 1
	prefixCheck
)

// KVStore key prefixes
var (
	KeyPrefixCoin  = []byte{prefixCoin}
	KeyPrefixCheck = []byte{prefixCheck}
)
