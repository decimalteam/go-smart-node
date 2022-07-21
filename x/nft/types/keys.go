package types

import (
	"encoding/binary"
	"fmt"
	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "nft"

	// StoreKey is the default store key for NFT
	StoreKey = ModuleName

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
	CollectionsKeyPrefix = []byte{NFTPrefix, 0x00} // key for NFT collections
	NFTKeyPrefix         = []byte{NFTPrefix, 0x01} // key for NFTs
	OwnersKeyPrefix      = []byte{NFTPrefix, 0x02} // key for balance of NFTs held by an address
	SubTokenKeyPrefix    = []byte{NFTPrefix, 0x03}
	TokenURIKeyPrefix    = []byte{NFTPrefix, 0x05}
)

const OwnerKeyHashLength = 54

// GetCollectionKey gets the key of a collection
func GetCollectionKey(denom string) []byte {
	bs := getHash(denom)

	return append(CollectionsKeyPrefix, bs...)
}

// GetNFTKey gets the key of a nft
func GetNFTKey(id string) []byte {
	bs := getHash(id)

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

// GetOwnerCollectionsKey gets the key prefix for all the collections owned by an account address
func GetOwnerCollectionsKey(address sdk.AccAddress) []byte {
	return append(OwnersKeyPrefix, address.Bytes()...)
}

// GetOwnerCollectionByDenomKey gets the key of a collection owned by an account address
func GetOwnerCollectionByDenomKey(address sdk.AccAddress, denom string) []byte {
	bs := getHash(denom)

	return append(GetOwnerCollectionsKey(address), bs...)
}

func GetSubTokensKey(id string) []byte {
	bsID := getHash(id)
	return append(SubTokenKeyPrefix, bsID...)
}

func GetSubTokenKey(id string, subTokenID uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, subTokenID)
	return append(GetSubTokensKey(id), b...)
}

func GetTokenURIKey(tokenURI string) []byte {
	return append(TokenURIKeyPrefix, []byte(tokenURI)...)
}

func getHash(denom string) []byte {
	h := tmhash.New()
	_, err := h.Write([]byte(denom))
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
