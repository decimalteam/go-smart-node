package types

const (
	// ModuleName is the name of the module
	ModuleName = "customfee"

	// StoreKey is the default store key for fee
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the fee store.
	QuerierRoute = ModuleName

	// RouterKey is the message route for the fee module
	RouterKey = ModuleName
)

const FeePrefix = 0x70
const separator = 0x01

var (
	PriceKeyPrefix = []byte{FeePrefix, 0x00} // key for store
)

// getCollectionID returns the collection ID concatenated from creator address denom hash and c.
func GetPriceKey(denom, quote string) []byte {
	key := []byte{}
	key = append(key, PriceKeyPrefix...)
	key = append(key, []byte(denom)...)
	key = append(key, separator)
	key = append(key, []byte(quote)...)
	return key
}
