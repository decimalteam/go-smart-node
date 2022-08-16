package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, initialPrice sdk.Dec) GenesisState {
	return GenesisState{
		Params:       params,
		InitialPrice: initialPrice,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		InitialPrice: sdk.OneDec(),
	}
}

func (gs *GenesisState) Validate() error {
	if gs.InitialPrice.LTE(sdk.ZeroDec()) {
		return ErrWrongPrice(gs.InitialPrice.String())
	}
	return gs.Params.Validate()
}
