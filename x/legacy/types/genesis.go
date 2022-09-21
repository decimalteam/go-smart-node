package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/errors"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

func (gs *GenesisState) Validate() error {
	// Check there are repeated addresses in legacy balances
	// and validate balances
	seenLegacy := make(map[string]bool)
	for _, rec := range gs.Records {
		if seenLegacy[rec.LegacyAddress] {
			return errors.LegacyAddressesDuplicatedOnGenesis
		}
		seenLegacy[rec.LegacyAddress] = true
		err := rec.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
