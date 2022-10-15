package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"golang.org/x/crypto/sha3"
)

// NewWallet returns a new Wallet.
func NewWallet(owners []string, weights []uint32, threshold uint32, salt []byte) (*Wallet, error) {
	// temporary wallet structure to create address
	w := Wallet{
		Address:   sdk.AccAddress(salt).String(), // use field Address as salt
		Owners:    owners,
		Weights:   weights,
		Threshold: threshold,
	}
	bz := sha3.Sum256(sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&w)))
	address := sdk.AccAddress(bz[12:])

	return &Wallet{
		Address:   address.String(),
		Owners:    owners,
		Weights:   weights,
		Threshold: threshold,
	}, nil
}

// NewTransaction returns a new Transaction.
func NewTransaction(wallet, receiver string, coins sdk.Coins, signersCount int, height int64, salt []byte) (*Transaction, error) {

	// temporary transaction struct to create TxID
	t := Transaction{
		Id:        sdk.AccAddress(salt).String(),
		Wallet:    wallet,
		Receiver:  receiver,
		Coins:     coins,
		Signers:   make([]string, signersCount),
		CreatedAt: height,
	}

	bz := sha3.Sum256(sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&t)))
	id, err := bech32.ConvertAndEncode(MultisigTransactionIDPrefix, bz[12:])
	if err != nil {
		return nil, errors.UnableToCreateTransaction
	}

	return &Transaction{
		Id:       id,
		Wallet:   wallet,
		Receiver: receiver,
		Coins:    coins,
		// create transaction withlist of empty strings; filled string mean 'signed by owner'
		Signers:   make([]string, signersCount),
		CreatedAt: height,
	}, nil
}
