package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
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
	if err := SupplyInvariantCheck(m.Collections, m.NFTs); err != nil {
		return err
	}

	if err := SubTokensInvariantCheck(m.NFTs, m.SubTokens); err != nil {
		return err
	}

	return nil
}

func SupplyInvariantCheck(collections []Collection, nfts []BaseNFT) error {
	totalSupply := 0
	for _, collection := range collections {
		totalSupply += collection.Supply()
	}
	if len(nfts) != totalSupply {
		return errors.NftSupply
	}

	return nil
}

func SubTokensInvariantCheck(nfts []BaseNFT, subTokens map[string]SubTokens) error {
	for _, nft := range nfts {
		subTokenLength := 0

		// validate: sub tokens count equal to count from owners
		for _, owner := range nft.GetOwners() {
			subTokenLength += owner.SubTokenIDs.Len()
		}

		if subTokenLength != len(subTokens[nft.ID].SubTokens) {
			return errors.InvalidSubTokensLen
		}

		// validate: all subtokens have owners
		for _, subToken := range subTokens[nft.ID].SubTokens {
			subTokenHasOwner := false
			for _, owner := range nft.GetOwners() {
				subTokenHasOwner = subTokenHasOwner || owner.SubTokenIDs.Has(subToken.ID)
			}
			if !subTokenHasOwner {
				return errors.UnknownSubTokenForNFT
			}
		}
	}

	return nil
}
