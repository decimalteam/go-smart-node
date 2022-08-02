package types

const (
	// ModuleName is the name of the module
	ModuleName = "basedenomfee"

	// StoreKey is the default store key for NFT
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the NFT store.
	QuerierRoute = ModuleName

	// RouterKey is the message route for the NFT module
	RouterKey = ModuleName
)

const FeePrefix = 0x70

var (
	BaseDenomPriceKeyPrefix = []byte{FeePrefix, 0x00} // key for NFT collections
)

// GetBaseDenomPriceKey gets the key of a collection
func GetBaseDenomPriceKey() []byte {
	return BaseDenomPriceKeyPrefix
}
