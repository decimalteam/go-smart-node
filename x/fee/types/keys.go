package types

const (
	// ModuleName is the name of the module
	ModuleName = "basedenomfee"

	// StoreKey is the default store key for fee
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the fee store.
	QuerierRoute = ModuleName

	// RouterKey is the message route for the fee module
	RouterKey = ModuleName
)

const FeePrefix = 0x70

var (
	BaseDenomPriceKeyPrefix = []byte{FeePrefix, 0x00} // key for store
)

func GetBaseDenomPriceKey() []byte {
	return BaseDenomPriceKeyPrefix
}
