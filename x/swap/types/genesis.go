package types

import (
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/x/swap/errors"
)

func NewGenesisState(params Params, swaps []Swap, chains []Chain) GenesisState {
	return GenesisState{
		Swaps:  swaps,
		Chains: chains,
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
	var seenIDs = make(map[uint32]bool)
	for _, chain := range gs.Chains {
		if seenIDs[chain.Id] {
			return fmt.Errorf("dublicate swap chain id %d", chain.Id)
		}
		if chain.Id == 0 {
			return errors.InvalidChainNumber
		}
		if chain.Name == "" {
			return errors.InvalidChainName
		}
	}
	return gs.Params.Validate()
}
