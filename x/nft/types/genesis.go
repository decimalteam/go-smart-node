package types

// NewGenesisState creates a new genesis state.
func NewGenesisState(collections []Collection, nfts []BaseNFT, subTokens map[string]SubTokens) *GenesisState {
	return &GenesisState{
		Collections: collections,
		Nfts:        nfts,
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
	// TODO validate
	return nil
}
