package types

import (
	fmt "fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(wallets []Wallet) GenesisState {
	return GenesisState{
		Wallets: wallets,
	}
}

// DefaultGenesisState sets default evm genesis state with empty accounts and
// default params and chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Wallets: []Wallet{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs *GenesisState) Validate() error {
	seenWallets := make(map[string]bool)
	for _, wallet := range gs.Wallets {
		if seenWallets[wallet.Address] {
			return fmt.Errorf("multisig wallet address duplicated on genesis: '%s'", wallet.Address)
		}
		seenWallets[wallet.Address] = true
		if _, err := sdk.AccAddressFromBech32(wallet.Address); err != nil {
			return ErrInvalidOwner(wallet.Address)
		}
		// Validation like MsgCreateWallet.ValidateBasic()
		// Validate owner count
		if len(wallet.Owners) < MinOwnerCount || len(wallet.Owners) > MaxOwnerCount {
			return ErrInvalidOwnerCount(strconv.Itoa(len(wallet.Owners)), strconv.Itoa(MinOwnerCount), strconv.Itoa(MaxOwnerCount))
		}
		// Validate weight count
		if len(wallet.Owners) != len(wallet.Weights) {
			return ErrInvalidWeightCount(strconv.Itoa(len(wallet.Weights)), strconv.Itoa(len(wallet.Owners)))
		}
		// Validate owners (ensure there are no duplicates)
		seenOwners := make(map[string]bool, len(wallet.Owners))
		for i := 0; i < len(wallet.Owners); i++ {
			owner := wallet.Owners[i]
			if _, err := sdk.AccAddressFromBech32(owner); err != nil {
				return ErrInvalidOwner(owner)
			}
			if seenOwners[owner] {
				return ErrDuplicateOwner(owner)
			}
			seenOwners[owner] = true
		}
		// Validate weights
		var sumOfWeights uint64
		for i := 0; i < len(wallet.Weights); i++ {
			if wallet.Weights[i] < MinWeight {
				return ErrInvalidWeight(strconv.Itoa(MinWeight), "less")
			}
			if wallet.Weights[i] > MaxWeight {
				return ErrInvalidWeight(strconv.Itoa(MaxWeight), "greater")
			}
			sumOfWeights += wallet.Weights[i]
		}
		if sumOfWeights < wallet.Threshold {
			return ErrInvalidThreshold(strconv.FormatUint(sumOfWeights, 10), strconv.FormatUint(wallet.Threshold, 10))
		}
	}
	return nil
}
