package types

import (
	"strings"
)

// NewCollection creates a new NFT Collection
func NewCollection(denom string, nfts []string) Collection {
	return Collection{
		Denom: strings.TrimSpace(denom),
		NFTs:  nfts,
	}
}

// AddNFT adds an NFT to the collection
func (collection Collection) AddNFT(id string) Collection {
	if collection.NFTs.Has(id) {
		return collection
	}

	collection.NFTs = append(collection.NFTs, id).Sort()

	return collection
}

// Supply gets the total supply of NFTs of a collection
func (collection Collection) Supply() int {
	return len(collection.NFTs)
}
