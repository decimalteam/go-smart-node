package types

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, collections []Collection) *GenesisState {
	return &GenesisState{
		Params:      params,
		Collections: collections,
	}
}

// DefaultGenesisState returns a default genesis state.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), []Collection{})
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (m GenesisState) Validate() error {
	// TODO
	return nil
}
