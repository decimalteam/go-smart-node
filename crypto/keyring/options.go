package keyring

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	ethhd "github.com/evmos/ethermint/crypto/hd"
)

var (
	// SupportedAlgorithms defines the list of signing algorithms used on Decimal:
	//  - secp256k1     (Cosmos SDK)
	//  - eth_secp256k1 (Ethereum)
	SupportedAlgorithms = keyring.SigningAlgoList{hd.Secp256k1, ethhd.EthSecp256k1}

	// SupportedAlgorithmsLedger defines the list of signing algorithms used on Decimal for the Ledger device:
	//  - secp256k1     (Cosmos SDK)
	//  - eth_secp256k1 (Ethereum)
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{hd.Secp256k1, ethhd.EthSecp256k1}
)

// Option defines a function keys options for the ethereum Secp256k1 curve.
// It supports eth_secp256k1 keys for accounts.
func Option() keyring.Option {
	return func(options *keyring.Options) {
		options.SupportedAlgos = SupportedAlgorithms
		options.SupportedAlgosLedger = SupportedAlgorithmsLedger
	}
}
