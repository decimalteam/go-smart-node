package types

import fmt "fmt"

func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

func (gs *GenesisState) Validate() error {
	// Check there are repeated addresses in legacy balances
	// and validate balances
	seenLegacy := make(map[string]bool)
	for _, rec := range gs.LegacyRecords {
		if seenLegacy[rec.Address] {
			return fmt.Errorf("legacy address duplicated on genesis: '%s'", rec.Address)
		}
		seenLegacy[rec.Address] = true
		err := rec.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
