package wallet

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"

	config "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	ethermintHd "github.com/tharsis/ethermint/crypto/hd"
	ethermint "github.com/tharsis/ethermint/types"
)

const (
	derivationPath = "m/44'/60'/0'/0/0"
	addressPrefix  = "dx"
)

// Account contains private key of the account that allows to sign transactions to broadcast to the blockchain.
type Account struct {
	privateKeyTM *ethsecp256k1.PrivKey
	publicKeyTM  *ethsecp256k1.PubKey
	address      string

	// These fields are used only for signing transactions:
	chainID       string
	accountNumber int64
	sequence      int64
}

// NewAccount creates new account with random mnemonic.
func NewAccount(password string) (*Account, error) {
	mnemonic, err := NewMnemonic(password)
	if err != nil {
		return nil, err
	}
	return NewAccountFromMnemonicWords(mnemonic.Words(), password)
}

// NewAccountFromMnemonicWords creates account from mnemonic presented as set of words.
func NewAccountFromMnemonicWords(words string, password string) (*Account, error) {
	var result Account
	bz, err := ethermintHd.EthSecp256k1.Derive()(words, password, ethermint.BIP44HDPath)
	if err != nil {
		return nil, err
	}
	result.privateKeyTM = &ethsecp256k1.PrivKey{Key: bz}
	result.publicKeyTM = result.privateKeyTM.PubKey().(*ethsecp256k1.PubKey)
	result.address, err = bech32.ConvertAndEncode(config.Bech32Prefix, result.publicKeyTM.Address())
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// WithChainID sets chain ID of network.
func (acc *Account) WithChainID(chainID string) *Account {
	acc.chainID = chainID
	return acc
}

// WithAccountNumber sets accounts's number.
func (acc *Account) WithAccountNumber(accountNumber uint64) *Account {
	acc.accountNumber = int64(accountNumber)
	return acc
}

// WithSequence sets accounts's sequence (last used nonce).
func (acc *Account) WithSequence(sequence uint64) *Account {
	acc.sequence = int64(sequence)
	return acc
}

// Address returns accounts's address in bech32 format.
func (acc *Account) Address() string {
	return acc.address
}

// ChainID returns chain ID of network.
func (acc *Account) ChainID() string {
	return acc.chainID
}

// AccountNumber returns accounts's number.
func (acc *Account) AccountNumber() int64 {
	return acc.accountNumber
}

// Sequence returns accounts's sequence (last used nonce).
func (acc *Account) Sequence() int64 {
	return acc.sequence
}

/*
// CreateTransaction creates new transaction with specified messages and parameters.
func (acc *Account) CreateTransaction(msgs []sdk.Msg, fee auth.StdFee, memo string) sdk.Tx {
	return sdkTx.Tx{
		Body:
	}

}

// SignTransaction signs transaction and appends signature to transaction signatures.
func (acc *Account) SignTransaction(tx auth.StdTx) (auth.StdTx, error) {

	// Check chain ID, account number and sequence
	if len(acc.chainID) == 0 {
		return tx, errors.New("chain ID is not set up")
	}
	if acc.accountNumber < 0 || acc.sequence < 0 {
		return tx, errors.New("account number or sequence is not set up")
	}

	// Retrieve transaction bytes required to sign
	bytesToSign := auth.StdSignBytes(
		acc.chainID, uint64(acc.accountNumber), uint64(acc.sequence),
		tx.Fee, tx.Msgs, tx.Memo,
	)

	// Sign bytes prepared to sign
	signatureBytes, err := acc.privateKeyTM.Sign(bytesToSign)
	if err != nil {
		return tx, err
	}

	// Prepare auth.StdSignature object
	signature := auth.StdSignature{
		PubKey:    acc.publicKeyTM,
		Signature: signatureBytes,
	}

	// Copy input transaction and append signature to the list
	tx.Signatures = append(tx.Signatures, signature)

	return tx, err
}
*/
