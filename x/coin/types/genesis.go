package types

import (
	"encoding/hex"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/config"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, coins []Coin, checks []Check) *GenesisState {
	return &GenesisState{
		Params: params,
		Coins:  coins,
		Checks: checks,
	}
}

// DefaultGenesisState returns a default genesis state.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), []Coin{}, []Check{})
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs *GenesisState) Validate() error {
	// Check coin denom for correct regexp
	if !config.CoinDenomValidator.MatchString(gs.Params.BaseDenom) {
		return errors.InvalidCoinDenom
	}
	// Check coin title maximum length
	if len(gs.Params.BaseTitle) > config.MaxCoinTitleLength {
		return errors.InvalidCoinTitle
	}
	// Check coin initial volume to be correct
	if gs.Params.BaseVolume.LT(config.MinCoinSupply) || gs.Params.BaseVolume.GT(config.MaxCoinSupply) {
		return errors.InvalidCoinInitialVolume
	}
	// Check there are no coins with the same denom
	seenCoins := make(map[string]bool)
	for _, coin := range gs.Coins {
		if seenCoins[coin.Denom] {
			return errors.DuplicateCoinInGenesis
		}
		// Validate coin
		// if err := coin.Validate(); err != nil {
		// 	return err
		// }
		seenCoins[coin.Denom] = true
	}
	// Check there are no checks with the same hash
	seenChecks := make(map[string]bool)
	for _, check := range gs.Checks {
		checkHash := check.HashFull()
		checkHashStr := hex.EncodeToString(checkHash[:])
		if seenChecks[checkHashStr] {
			return errors.DuplicateCheckInGenesis
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
