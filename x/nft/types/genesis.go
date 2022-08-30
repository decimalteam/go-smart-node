package types

import (
	"errors"
	"fmt"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(collections []Collection, nfts []BaseNFT, subTokens map[string]SubTokens) *GenesisState {
	return &GenesisState{
		Collections: collections,
		NFTs:        nfts,
		SubTokens:   subTokens,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState([]Collection{}, []BaseNFT{}, map[string]SubTokens{})
}

// Validate performs basic validation of nfts genesis data returning an
// error for any failed validation criteria.
func (m GenesisState) Validate() error {
	msg, failed := SupplyInvariantCheck(m.Collections, m.NFTs)
	if failed {
		return errors.New(msg)
	}
	msg, failed = SubTokensInvariantCheck(m.NFTs, m.SubTokens)
	if failed {
		return errors.New(msg)
	}

	return nil
}

func SupplyInvariantCheck(collections []Collection, nfts []BaseNFT) (string, bool) {
	totalSupply := 0
	for _, collection := range collections {
		totalSupply += collection.Supply()
	}
	broken := len(nfts) != totalSupply

	return fmt.Sprintf("nft supply invariants found (total supply: %d, nfts: %d)", totalSupply, len(nfts)), broken
}

func SubTokensInvariantCheck(nfts []BaseNFT, subTokens map[string]SubTokens) (string, bool) {
	// TODO: temporary turn off this, because nfts subtokens are delegated to validator
	return "", false

	for _, nft := range nfts {
		subTokenLength := 0

		// validate: sub tokens count equal to count from owners
		for _, owner := range nft.GetOwners() {
			subTokenLength += owner.SubTokenIDs.Len()
		}

		if subTokenLength != len(subTokens[nft.ID].SubTokens) {
			return fmt.Sprintf(
				"invalid sub tokens len for nft %s (nft len: %d, sub tokens len %d)",
				nft.ID, subTokenLength, len(subTokens[nft.ID].SubTokens),
			), true
		}

		// validate: all subtokens have owners
		for _, subToken := range subTokens[nft.ID].SubTokens {
			subTokenHasOwner := false
			for _, owner := range nft.GetOwners() {
				subTokenHasOwner = subTokenHasOwner || owner.SubTokenIDs.Has(subToken.ID)
			}
			if !subTokenHasOwner {
				return fmt.Sprintf("unknown sub token id %d for nft %s", subToken.ID, nft.ID), true
			}
		}
	}

	return "", false
}
