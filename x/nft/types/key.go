package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

const (
	// ModuleName defines the module name
	ModuleName = "nft"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	// RouterKey defines module's messages routing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// Lets assume parts of keys are:
//   <creator>         = <address>             [20 bytes]
//   <owner>           = <address>             [20 bytes]
//   <denom_hash>      = <hash(denom)>         [32 bytes]
//   <collection_id>   = <creator><denom_hash> [52 bytes]
//   <token_id>        = <hash(id)>            [32 bytes]
//   <token_uri>       = <hash(uri)>           [32 bytes]
//   <sub_token_id>    = <uint32>              [ 4 bytes]

// Main records containing all necessary information:
//   - Collections:          0x11<creator><denom_hash>           : <Collection>
//   - CollectionCounters:   0x12<creator><denom_hash>           : <CollectionCounter>
//   - Tokens:               0x21<token_id>                      : <Token>
//   - TokenCounters:        0x22<token_id>                      : <TokenCounter>
//   - SubTokens:            0x31<token_id><sub_token_id>        : <SubToken>

// There are also set of indexes for iterating over records:
//   - TokenURIs:            0xA1<token_uri>                     : []byte{}
//   - TokensByCollection:   0xA2<creator><denom_hash><token_id> : []byte{}
//   - SubTokensByOwner:     0xA3<owner><token_id><sub_token_id> : []byte{}

var (
	keyPrefixCollection         = []byte{0x11} // prefix for each key to NFT collections
	keyPrefixCollectionCounter  = []byte{0x12} // prefix for each key to NFT collection counter object
	keyPrefixToken              = []byte{0x21} // prefix for each key to NFT token
	keyPrefixTokenCounter       = []byte{0x22} // prefix for each key to NFT token counter object
	keyPrefixSubToken           = []byte{0x31} // prefix for each key to NFT sub-token
	keyPrefixTokenURI           = []byte{0xA1} // prefix for each key to index to NFT token URI
	keyPrefixTokensByCollection = []byte{0xA2} // prefix for each key to index for NFT tokens by collection creator and denom
	keyPrefixSubTokensByOwner   = []byte{0xA3} // prefix for each key to index for NFT sub-tokens by owner address
)

// GetCollectionsKey returns the key prefix of the NFT collections.
func GetCollectionsKey() []byte {
	return keyPrefixCollection
}

// GetCollectionsByCreatorKey returns the key prefix of the NFT collections created by specified creator.
func GetCollectionsByCreatorKey(creator sdk.AccAddress) []byte {
	return append(GetCollectionsKey(), creator.Bytes()...)
}

// GetCollectionKey returns the key of the NFT collection.
func GetCollectionKey(creator sdk.AccAddress, denom string) []byte {
	return append(GetCollectionsKey(), getCollectionID(creator, denom)...)
}

// GetCollectionCounterKey returns the key of the NFT collection counter.
func GetCollectionCounterKey(creator sdk.AccAddress, denom string) []byte {
	return append(keyPrefixCollectionCounter, getCollectionID(creator, denom)...)
}

// GetTokensKey returns the key prefix of the NFT tokens.
func GetTokensKey() []byte {
	return keyPrefixToken
}

// GetTokenKey returns the key of the NFT token.
func GetTokenKey(id string) []byte {
	return append(GetTokensKey(), helpers.CalcHashSHA256(id)...)
}

func GetTokenKeyByIDHash(id []byte) []byte {
	return append(GetTokensKey(), id...)
}

// GetTokenCounterKey returns the key of the NFT token counter.
func GetTokenCounterKey(id string) []byte {
	return append(keyPrefixTokenCounter, helpers.CalcHashSHA256(id)...)
}

// GetSubTokensKey returns the key prefix of the NFT sub-tokens.
func GetSubTokensKey(id string) []byte {
	return append(keyPrefixSubToken, helpers.CalcHashSHA256(id)...)
}

// GetSubTokenKey returns the key of the NFT sub-token.
func GetSubTokenKey(id string, index uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, index)
	return append(GetSubTokensKey(id), b...)
}

// GetTokenURIKey returns the key of the NFT token URI.
func GetTokenURIKey(tokenURI string) []byte {
	return append(keyPrefixTokenURI, helpers.CalcHashSHA256(tokenURI)...)
}

// GetTokensByCollectionKey returns the key prefix of the NFT tokens of specific collection.
func GetTokensByCollectionKey(creator sdk.AccAddress, denom string) []byte {
	return append(keyPrefixTokensByCollection, getCollectionID(creator, denom)...)
}

// GetTokenByCollectionKey returns the key prefix of the NFT token of specific collection.
func GetTokenByCollectionKey(creator sdk.AccAddress, denom string, id string) []byte {
	return append(GetTokensByCollectionKey(creator, denom), helpers.CalcHashSHA256(id)...)
}

// GetSubTokenByOwnerKey returns the key of the NFT sub-token of specific owner address NFT token and index.
func GetSubTokenByOwnerKey(owner sdk.AccAddress, id string, index uint32) []byte {
	return append(keyPrefixSubTokensByOwner, append(owner.Bytes(), getSubTokenID(id, index)...)...)
}

// getCollectionID returns the collection ID concatenated from creator address denom hash and c.
func getCollectionID(creator sdk.AccAddress, denom string) []byte {
	return append(creator.Bytes(), helpers.CalcHashSHA256(denom)...)
}

// getSubTokenID returns the sub-token ID concatenated from token ID and sub-token index.
func getSubTokenID(id string, index uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, index)
	return append(helpers.CalcHashSHA256(id), b...)
}
