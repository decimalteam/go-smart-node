package wallet

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"

	config "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	decimalTypes "bitbucket.org/decimalteam/go-smart-node/types"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	ethermintHd "github.com/evmos/ethermint/crypto/hd"
	ethermint "github.com/evmos/ethermint/types"
)

// Account contains private key of the account that allows to sign transactions to broadcast to the blockchain.
type Account struct {
	privateKeyTM  *ethsecp256k1.PrivKey
	publicKeyTM   *ethsecp256k1.PubKey
	address       string
	legacyAddress string

	// These fields are used only for signing transactions:
	chainID       string
	accountNumber uint64
	sequence      uint64
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
	result.legacyAddress, err = decimalTypes.GetLegacyAddressFromPubKey(result.privateKeyTM.PubKey().Bytes())
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
	acc.accountNumber = accountNumber
	return acc
}

// WithSequence sets accounts's sequence (last used nonce).
func (acc *Account) WithSequence(sequence uint64) *Account {
	acc.sequence = sequence
	return acc
}

func (acc *Account) IncrementSequence() {
	acc.sequence++
}

// Address returns accounts's address in bech32 format.
func (acc *Account) Address() string {
	return acc.address
}

// Address returns accounts's address in bech32 format.
func (acc *Account) LegacyAddress() string {
	return acc.legacyAddress
}

// SdkAddress returns accounts's cosmos AccAddress ([]byte)
func (acc *Account) SdkAddress() sdk.AccAddress {
	return sdk.AccAddress(acc.publicKeyTM.Address())
}

// ChainID returns chain ID of network.
func (acc *Account) ChainID() string {
	return acc.chainID
}

// AccountNumber returns accounts's number.
func (acc *Account) AccountNumber() uint64 {
	return acc.accountNumber
}

// Sequence returns accounts's sequence (last used nonce).
func (acc *Account) Sequence() uint64 {
	return acc.sequence
}

// Sequence returns accounts's sequence (last used nonce).
func (acc *Account) PubKey() cryptoTypes.PubKey {
	return acc.publicKeyTM
}

// Sign data by private key and returns signature
func (acc *Account) Sign(bytesToSign []byte) ([]byte, error) {
	return acc.privateKeyTM.Sign(bytesToSign)
}
