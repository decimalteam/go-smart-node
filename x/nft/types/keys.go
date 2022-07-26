package types

import (
	"encoding/binary"
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
// - Colections: 0x00<denom_bytes_key> : <Collection>
//
// - NFTs: 0x01<token_id_bytes_key> : <NFT>
//
// - OwnerCollection: 0x02<address_bytes_key><denom_bytes_key> : <OwnerCollection>
//
// - SubTokens: 0x03<nft_id_bytes_key><sub_token_id_bytes_key> : <SubToken>
//
// - TokenURI: 0x04<token_uri_bytes_key> : <[]byte{}>

const NFTPrefix = 0x60

var (
	CollectionsKeyPrefix      = []byte{NFTPrefix, 0x00} // key for NFT collections
	NFTKeyPrefix              = []byte{NFTPrefix, 0x01} // key for NFTs
	OwnerCollectionsKeyPrefix = []byte{NFTPrefix, 0x02} // key for balance of NFTs held by an address
	SubTokenKeyPrefix         = []byte{NFTPrefix, 0x03}
	TokenURIKeyPrefix         = []byte{NFTPrefix, 0x04}
)

// GetCollectionKey gets the key of a collection
func GetCollectionKey(denom string) []byte {
	bs := getHash(denom)

	return append(CollectionsKeyPrefix, bs...)
}

// GetNFTKey gets the key of a nft
func GetNFTKey(id string) []byte {
	bs := getHash(id)

	return append(NFTKeyPrefix, bs...)
}

// GetOwnerCollectionsKey gets the key prefix for all the collections owned by an account address
func GetOwnerCollectionsKey(address sdk.AccAddress) []byte {
	return append(OwnerCollectionsKeyPrefix, address.Bytes()...)
}

// GetOwnerCollectionByDenomKey gets the key of a collection owned by an account address
func GetOwnerCollectionByDenomKey(address sdk.AccAddress, denom string) []byte {
	bs := getHash(denom)

	return append(GetOwnerCollectionsKey(address), bs...)
}

func GetSubTokensKey(nftID string) []byte {
	bsID := getHash(nftID)
	return append(SubTokenKeyPrefix, bsID...)
}

func GetSubTokenKey(nftID string, subTokenID uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, subTokenID)
	return append(GetSubTokensKey(nftID), b...)
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
