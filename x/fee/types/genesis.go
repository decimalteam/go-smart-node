package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, initialPrice sdk.Dec) GenesisState {
	return GenesisState{
		Params:      params,
		IntialPrice: initialPrice,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		IntialPrice: sdk.OneDec(),
	}
}

func (gs *GenesisState) Validate() error {
	if gs.IntialPrice.LTE(sdk.ZeroDec()) {
		return ErrWrongPrice(gs.IntialPrice.String())
	}
	return gs.Params.Validate()
}
