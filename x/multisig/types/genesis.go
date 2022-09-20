package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/errors"
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
			return errors.DuplicateWallet
		}
		seenWallets[wallet.Address] = true
		if _, err := sdk.AccAddressFromBech32(wallet.Address); err != nil {
			return errors.InvalidWallet
		}
		// Validation like MsgCreateWallet.ValidateBasic()
		// Validate owner count
		if len(wallet.Owners) < MinOwnerCount || len(wallet.Owners) > MaxOwnerCount {
			return errors.InvalidOwnerCount
		}
		// Validate weight count
		if len(wallet.Owners) != len(wallet.Weights) {
			return errors.InvalidWeightCount
		}
		// Validate owners (ensure there are no duplicates)
		seenOwners := make(map[string]bool, len(wallet.Owners))
		for i := 0; i < len(wallet.Owners); i++ {
			owner := wallet.Owners[i]
			if _, err := sdk.AccAddressFromBech32(owner); err != nil {
				return errors.InvalidOwner
			}
			if seenOwners[owner] {
				return errors.DuplicateOwner
			}
			seenOwners[owner] = true
		}
		walletOwners[wallet.Address] = wallet.Owners
		// Validate weights
		var sumOfWeights uint32
		for i := 0; i < len(wallet.Weights); i++ {
			if wallet.Weights[i] < MinWeight || wallet.Weights[i] > MaxWeight {
				return errors.InvalidWeight
			}
			sumOfWeights += wallet.Weights[i]
		}
		if sumOfWeights < wallet.Threshold {
			return errors.InvalidThreshold
		}
	}
	//
	seenTxs := make(map[string]bool)
	for _, tx := range gs.Transactions {
		if seenTxs[tx.Id] {
			return errors.DuplicateTxsID
		}
		seenTxs[tx.Id] = true
		// wallet must exist
		if !seenWallets[tx.Wallet] {
			return errors.UnknownWalletInTx
		}
		// validate signers
		if len(tx.Signers) != len(walletOwners[tx.Wallet]) {
			return errors.TxSignersNotEqualToWalletOwners
		}
		for _, signer := range tx.Signers {
			if signer == "" {
				continue
			}
			if !strings.StringInSlice(signer, walletOwners[tx.Wallet]) {
				return errors.UnknownSignerInTx
			}
		}
	}
	return nil
}
