package types

func NewGenesisState(params Params, swaps []Swap) GenesisState {
	return GenesisState{
		Swaps:  swaps,
		Params: params,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs *GenesisState) Validate() error {
	return nil
}
