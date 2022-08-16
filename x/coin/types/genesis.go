package types

import (
	"encoding/hex"
	"fmt"
	"regexp"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, coins []Coin, checks []Check) GenesisState {
	return GenesisState{
		Params: params,
		Coins:  coins,
		Checks: checks,
	}
}

// DefaultGenesisState sets default evm genesis state with empty accounts and
// default params and chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Coins:  []Coin{},
		Checks: []Check{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs *GenesisState) Validate() error {
	// Check coin title maximum length
	if len(gs.Params.BaseTitle) > maxCoinNameBytes {
		return ErrInvalidCoinTitle(gs.Params.BaseTitle)
	}
	// Check coin symbol for correct regexp
	if match, _ := regexp.MatchString(allowedCoinSymbols, gs.Params.BaseSymbol); !match {
		return ErrInvalidCoinSymbol(gs.Params.BaseSymbol)
	}
	// Check coin initial volume to be correct
	if gs.Params.BaseInitialVolume.LT(MinCoinSupply) || gs.Params.BaseInitialVolume.GT(maxCoinSupply) {
		return ErrInvalidCoinInitialVolume(gs.Params.BaseInitialVolume.String())
	}
	// Check there are no coins with the same symbol
	seenSymbols := make(map[string]bool)
	for _, coin := range gs.Coins {
		if seenSymbols[coin.Symbol] {
			return fmt.Errorf("coin symbol duplicated on genesis: '%s'", coin.Symbol)
		}
		// Validate coin
		// if err := coin.Validate(); err != nil {
		// 	return err
		// }
		seenSymbols[coin.Symbol] = true
	}
	// Check there are no checks with the same hash
	seenChecks := make(map[string]bool)
	for _, check := range gs.Checks {
		checkHash := check.HashFull()
		checkHashStr := hex.EncodeToString(checkHash[:])
		if seenChecks[checkHashStr] {
			return fmt.Errorf("check hash duplicated on genesis: '%X'", checkHash[:])
		}
		// Validate check
		// if err := check.Validate(); err != nil {
		// 	return err
		// }
		seenChecks[checkHashStr] = true
	}

	// Validate params
	return gs.Params.Validate()
}
