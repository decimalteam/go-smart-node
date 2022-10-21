package types

const (
	// ModuleName defines the module name
	ModuleName = "multisig"

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
	prefixWallet               byte = iota + 1 // 1
	prefixTransaction                          // 2
	prefixUniversalTransaction                 // 3
	prefixUniversalSign                        // 4
	prefixCompleted                            // 5
)

// KVStore key prefixes
var (
	KeyPrefixWallet      = []byte{prefixWallet}      // 0x01
	KeyPrefixTransaction = []byte{prefixTransaction} // 0x02

	KeyPrefixUniversalTransaction = []byte{prefixUniversalTransaction} // 0x02
	KeyPrefixCompletedTransaction = []byte{prefixCompleted}
)

func GetSignatureKey(txID string, signer string) []byte {
	key := []byte{prefixUniversalSign}
	key = append(key, []byte(txID)...)
	key = append(key, prefixUniversalSign)
	key = append(key, []byte(signer)...)
	return key
}

func GetSignaturePrefixKey(txID string) []byte {
	key := []byte{prefixUniversalSign}
	key = append(key, []byte(txID)...)
	key = append(key, prefixUniversalSign)
	return key
}

func ExtractSignerFromKey(key []byte, txID string) string {
	skip := 1 + len(txID) + 1
	return string(key[skip:])
}
