package types

import "fmt"

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, coins []Coin) GenesisState {
	return GenesisState{
		Params: params,
		Coins:  coins,
	}
}

// DefaultGenesisState sets default evm genesis state with empty accounts and
// default params and chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Coins:  []Coin{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs *GenesisState) Validate() error {
	seenSymbols := make(map[string]bool)
	for _, c := range gs.Coins {
		if seenSymbols[c.Symbol] {
			return fmt.Errorf("coin symbol duplicated on genesis: '%s'", c.Symbol)
		}
		// if err := c.Validate(); err != nil {
		// 	return err
		// }
		seenSymbols[c.Symbol] = true
	}
	return gs.Params.Validate()
}
