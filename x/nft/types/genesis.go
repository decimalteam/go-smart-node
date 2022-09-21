package types

import fmt "fmt"

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
	for _, coll := range m.Collections {
		if err := coll.Validate(); err != nil {
			fmt.Printf("error in check collection '%s'\n", coll.Denom)
			return err
		}
		for _, token := range coll.Tokens {
			if err := token.Validate(); err != nil {
				fmt.Printf("error in check nft: collection '%s', token '%s'\n", coll.Denom, token.ID)
				return err
			}
			// iterate subtokens
			for _, subToken := range token.SubTokens {
				// validate subtoken
				if err := subToken.Validate(); err != nil {
					fmt.Printf("error in check subtoken: collection '%s', token '%s', subtoken: %d\n", coll.Denom, token.ID, subToken.ID)
					return err
				}
			}
		}
	}
	return nil
}
