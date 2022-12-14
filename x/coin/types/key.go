package types

const (
	// ModuleName defines the module name
	ModuleName = "coin"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key_
	MemStoreKey = "mem_" + ModuleName

	// RouterKey defines module's messages routing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// Coins and checks are stored as follow:
// - Coins:        0x11<denom_bytes>  : <Coin>
// - CoinVRs:      0x12<denom_bytes>  : <CoinVR>
// - Checks:       0x21<check_hash>   : <Check>

// KVStore key prefixes
var (
	keyPrefixCoin   = []byte{0x11} // prefix for each key to a coin
	keyPrefixCoinVR = []byte{0x12} // prefix for each key to a record containing coin volume and reserve
	keyPrefixCheck  = []byte{0x21} // prefix for each key to a redeemed check
)

// GetCoinsKey returns the key prefix of the coins.
func GetCoinsKey() []byte {
	return keyPrefixCoin
}

// GetCoinKey returns the key of the coin.
func GetCoinKey(denom string) []byte {
	return append(GetCoinsKey(), []byte(denom)...)
}

// GetCoinVRKey returns the key of the record containing coin volume and reserve.
func GetCoinVRKey(denom string) []byte {
	return append(keyPrefixCoinVR, []byte(denom)...)
}

// GetChecksKey returns the key prefix of the redeemed check.
func GetChecksKey() []byte {
	return keyPrefixCheck
}

// GetCheckKey returns the key of the redeemed check.
func GetCheckKey(hash []byte) []byte {
	return append(keyPrefixCheck, hash...)
}
