package wallet

import "github.com/cosmos/go-bip39"

const entropyBitsCount = 256

// Mnemonic contains entropy, seed and mnemonic words array
// which can be used for hierarchical deterministic extended keys.
type Mnemonic struct {
	entropy []byte
	words   string
	seed    []byte
}

// NewMnemonic creates a new random (crypto safe) Mnemonic. Use 128 bits for a 12 words code or 256 bits for a 24 words.
func NewMnemonic(password string) (*Mnemonic, error) {
	entropy, err := bip39.NewEntropy(entropyBitsCount)
	if err != nil {
		return nil, err
	}
	return NewMnemonicFromEntropy(entropy, password)
}

// NewMnemonicFromWords creates a Mnemonic based on a known list of words.
func NewMnemonicFromWords(words string, password string) (*Mnemonic, error) {
	entropy, err := bip39.MnemonicToByteArray(words)
	if err != nil {
		return nil, err
	}
	// (for bip39 from cosmos sdk) last byte contains checksum, so we need to cut last byte
	return NewMnemonicFromEntropy(entropy[:len(entropy)-1], password)
}

// NewMnemonicFromEntropy creates a Mnemonic based on a known entropy.
func NewMnemonicFromEntropy(entropy []byte, password string) (*Mnemonic, error) {
	words, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return &Mnemonic{
		entropy: entropy,
		words:   words,
		seed:    bip39.NewSeed(words, password),
	}, nil
}

// Entropy returns the entropy of the Mnemonic.
func (m *Mnemonic) Entropy() []byte {
	return m.entropy
}

// Words returns the words from the Mnemonic.
func (m *Mnemonic) Words() string {
	return m.words
}

// Seed returns the seed of the Mnemonic.
func (m *Mnemonic) Seed() []byte {
	return m.seed
}
