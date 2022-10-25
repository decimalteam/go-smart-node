package types

import (
	"encoding/binary"

	"bitbucket.org/decimalteam/go-smart-node/x/multisig/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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
func NewTransaction(unpacker codectypes.AnyUnpacker, wallet string, txContent codectypes.Any, signersCount int, height int64, salt []byte) (*Transaction, error) {

	h := make([]byte, 8)
	binary.BigEndian.PutUint64(h, uint64(height))
	idSource := []byte{}
	idSource = append(idSource, salt...)
	idSource = append(idSource, []byte(wallet)...)
	idSource = append(idSource, txContent.GetValue()...)
	idSource = append(idSource, h...)

	bz := sha3.Sum256(idSource)
	id, err := bech32.ConvertAndEncode(MultisigTransactionIDPrefix, bz[12:])
	if err != nil {
		return nil, errors.UnableToCreateTransaction
	}

	var msg sdk.Msg
	err = unpacker.UnpackAny(&txContent, &msg)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Id:        id,
		Wallet:    wallet,
		Message:   txContent,
		CreatedAt: height,
	}, nil
}
