package types

import (
	coin "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"encoding/binary"
	"fmt"
	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "nft"

	// StoreKey is the default store key for NFT
	StoreKey = coin.StoreKey

	// QuerierRoute is the querier route for the NFT store.
	QuerierRoute = ModuleName

	// RouterKey is the message route for the NFT module
	RouterKey = ModuleName

	// AddrLen defines a valid address length
	AddrLen = 20
)

// NFTs are stored as follow:
//
// - Colections: 0x00<denom_bytes_key> :<Collection>
//
// - Owners: 0x01<address_bytes_key><denom_bytes_key>: <Owner>

const NFTPrefix = 0x60

var (
	CollectionsKeyPrefix    = []byte{NFTPrefix, 0x00} // key for NFT collections
	OwnersKeyPrefix         = []byte{NFTPrefix, 0x01} // key for balance of NFTs held by an address
	SubTokenKeyPrefix       = []byte{NFTPrefix, 0x02}
	LastSubTokenIDKeyPrefix = []byte{NFTPrefix, 0x03}
	TokenURIKeyPrefix       = []byte{NFTPrefix, 0x04}
	TokenIDKeyPrefix        = []byte{NFTPrefix, 0x05}
)

const OwnerKeyHashLength = 54

// GetCollectionKey gets the key of a collection
func GetCollectionKey(denom string) []byte {
	bs := getHash(denom)

	return append(CollectionsKeyPrefix, bs...)
}

// SplitOwnerKey gets an address and denom from an owner key
func SplitOwnerKey(key []byte) (sdk.AccAddress, []byte) {
	if len(key) != OwnerKeyHashLength {
		panic(fmt.Sprintf("unexpected key length %d", len(key)))
	}
	address := key[2 : AddrLen+2]
	denomHashBz := key[AddrLen+1:]
	return sdk.AccAddress(address), denomHashBz
}

// GetOwnersKey gets the key prefix for all the collections owned by an account address
func GetOwnersKey(address sdk.AccAddress) []byte {
	return append(OwnersKeyPrefix, address.Bytes()...)
}

// GetOwnerKey gets the key of a collection owned by an account address
func GetOwnerKey(address sdk.AccAddress, denom string) []byte {
	bs := getHash(denom)

	return append(GetOwnersKey(address), bs...)
}

func GetSubTokenKey(denom, id string, subTokenID int64) []byte {
	bs := getHash(denom)
	bsID := getHash(id)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(subTokenID))
	return append(append(append(SubTokenKeyPrefix, bs...), bsID...), b...)
}

func GetLastSubTokenIDKey(denom, id string) []byte {
	bs := getHash(denom)
	bsID := getHash(id)
	return append(append(LastSubTokenIDKeyPrefix, bs...), bsID...)
}

func GetTokenURIKey(tokenURI string) []byte {
	return append(TokenURIKeyPrefix, []byte(tokenURI)...)
}

func GetTokenIDKey(id string) []byte {
	return append(TokenIDKeyPrefix, []byte(id)...)
}

func getHash(denom string) []byte {
	h := tmhash.New()
	_, err := h.Write([]byte(denom))
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
