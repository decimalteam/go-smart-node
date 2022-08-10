package types

import (
	fmt "fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/strings"
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
	walletOwners := make(map[string][]string) // cache of wallet owners to validate transactions
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
		// Validate legacy addresses count
		// valid: len(wallet.LegacyOwners) = 0 or len(wallet.Owners) == len(wallet.LegacyOwners)
		if len(wallet.LegacyOwners) > 0 && len(wallet.Owners) != len(wallet.LegacyOwners) {
			return fmt.Errorf("legacy owners address count != owners count, %d != %d", len(wallet.LegacyOwners), len(wallet.Owners))
		}
		// Validate weight count
		if len(wallet.Owners) != len(wallet.Weights) {
			return ErrInvalidWeightCount(strconv.Itoa(len(wallet.Weights)), strconv.Itoa(len(wallet.Owners)))
		}
		// Validate owners (ensure there are no duplicates)
		seenOwners := make(map[string]bool, len(wallet.Owners))
		seenLegacyOwners := make(map[string]bool, len(wallet.LegacyOwners))
		for i := 0; i < len(wallet.Owners); i++ {
			owner := wallet.Owners[i]
			legacyOwner := ""
			// LegacyOwners may be empty in case of full address conversion
			if len(wallet.LegacyOwners) > 0 {
				legacyOwner = wallet.LegacyOwners[i]
			}
			// good cases:
			// owner != "" && legacyOwner == ""
			// owner == "" && legacyOwner != ""
			if owner != "" && legacyOwner != "" {
				return fmt.Errorf("wallet %s: both addresses are filled at pos %d", wallet.Address, i)
			}
			if owner == "" && legacyOwner == "" {
				return fmt.Errorf("wallet %s: both addresses are empty at pos %d", wallet.Address, i)
			}
			if owner > "" {
				if _, err := sdk.AccAddressFromBech32(owner); err != nil {
					return ErrInvalidOwner(owner)
				}
				if seenOwners[owner] {
					return ErrDuplicateOwner(owner)
				}
				seenOwners[owner] = true
			}
			if legacyOwner > "" {
				if _, err := sdk.AccAddressFromBech32(legacyOwner); err != nil {
					return ErrInvalidOwner(legacyOwner)
				}
				if seenLegacyOwners[legacyOwner] {
					return ErrDuplicateOwner(legacyOwner)
				}
				seenLegacyOwners[legacyOwner] = true
			}
		}
		walletOwners[wallet.Address] = wallet.Owners
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
	//
	seenTxs := make(map[string]bool)
	for _, tx := range gs.Transactions {
		if seenTxs[tx.Id] {
			return fmt.Errorf("multisig transaction id duplicated on genesis: '%s'", tx.Id)
		}
		seenTxs[tx.Id] = true
		// wallet must exist
		if !seenWallets[tx.Wallet] {
			return fmt.Errorf("multisig transaction '%s' unknown wallet: '%s'", tx.Id, tx.Wallet)
		}
		// validate signers
		if len(tx.Signers) != len(walletOwners[tx.Wallet]) {
			return fmt.Errorf("multisig transaction '%s' signers countr != wallet owners: %d != %d",
				tx.Id, len(tx.Signers), len(walletOwners[tx.Wallet]))
		}
		for _, signer := range tx.Signers {
			if signer == "" {
				continue
			}
			if !strings.StringInSlice(signer, walletOwners[tx.Wallet]) {
				return fmt.Errorf("multisig transaction '%s' unknown signer: '%s'", tx.Id, signer)
			}
		}
	}
	return nil
}
