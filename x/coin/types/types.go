package types

// import (
// 	"bytes"
// 	"crypto/ecdsa"
// 	"crypto/sha256"
// 	"errors"
// 	"fmt"
// 	"math/big"
// 	"strings"

// 	"golang.org/x/crypto/ripemd160"
// 	"golang.org/x/crypto/sha3"

// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/rlp"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// ////////////////////////////////////////////////////////////////
// // Coin
// ////////////////////////////////////////////////////////////////

// type Coin struct {
// 	Title       string         `json:"title" yaml:"title"`                                   // Full coin title (Bitcoin)
// 	CRR         uint           `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
// 	Symbol      string         `json:"symbol" yaml:"symbol"`                                 // Short coin title (BTC)
// 	Reserve     sdk.Int        `json:"reserve" yaml:"reserve"`
// 	LimitVolume sdk.Int        `json:"limit_volume" yaml:"limit_volume"` // How many coins can be issued
// 	Volume      sdk.Int        `json:"volume" yaml:"volume"`
// 	Creator     sdk.AccAddress `json:"creator" yaml:"creator"`
// 	Identity    string         `json:"identity" yaml:"identity"`
// }

// func (c Coin) String() string {
// 	return strings.TrimSpace(fmt.Sprintf(`Title: %s
// 		CRR: %d
// 		Symbol: %s
// 		Reserve: %s
// 		LimitVolume: %s
// 		Volume: %s
// 		Creator: %s
// 	`, c.Title, c.CRR, c.Symbol, c.Reserve.String(), c.LimitVolume.String(), c.Volume.String(), c.Creator.String()))
// }

// func (c Coin) IsBase() bool {
// 	// if strings.HasPrefix(config.ChainID, "decimal-testnet") {
// 	// 	return c.Symbol == config.SymbolTestBaseCoin
// 	// } else {
// 	// 	return c.Symbol == config.SymbolBaseCoin
// 	// }
// 	return c.Symbol == "del"
// }

// ////////////////////////////////////////////////////////////////
// // Check
// ////////////////////////////////////////////////////////////////

// const HashLength = 32

// // Hash represents the 32 byte Keccak256 hash of arbitrary data.
// type Hash [HashLength]byte

// var (
// 	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
// )

// type Check struct {
// 	ChainID  string
// 	Coin     string
// 	Amount   *big.Int
// 	Nonce    []byte
// 	DueBlock uint64
// 	Lock     *big.Int
// 	V        *big.Int
// 	R        *big.Int
// 	S        *big.Int
// }

// func (check *Check) Sender() (sdk.AccAddress, error) {
// 	return recoverPlain(check.Hash(), check.R, check.S, check.V)
// }

// func (check *Check) LockPubKey() ([]byte, error) {
// 	sig := check.Lock.Bytes()

// 	if len(sig) < 65 {
// 		sig = append(make([]byte, 65-len(sig)), sig...)
// 	}

// 	hash := check.HashWithoutLock()

// 	pub, err := crypto.Ecrecover(hash[:], sig)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(pub) == 0 || pub[0] != 4 {
// 		return nil, errors.New("invalid public key")
// 	}

// 	return pub, nil
// }

// func (check *Check) HashWithoutLock() Hash {
// 	return rlpHash([]interface{}{
// 		check.ChainID,
// 		check.Coin,
// 		check.Amount,
// 		check.Nonce,
// 		check.DueBlock,
// 	})
// }

// func (check *Check) Hash() Hash {
// 	return rlpHash([]interface{}{
// 		check.ChainID,
// 		check.Coin,
// 		check.Amount,
// 		check.Nonce,
// 		check.DueBlock,
// 		check.Lock,
// 	})
// }

// func (check *Check) HashFull() Hash {
// 	return rlpHash([]interface{}{
// 		check.ChainID,
// 		check.Coin,
// 		check.Amount,
// 		check.Nonce,
// 		check.DueBlock,
// 		check.Lock,
// 		check.V,
// 		check.R,
// 		check.S,
// 	})
// }

// func (check *Check) Sign(prv *ecdsa.PrivateKey) error {
// 	h := check.Hash()
// 	sig, err := crypto.Sign(h[:], prv)
// 	if err != nil {
// 		return err
// 	}

// 	check.SetSignature(sig)

// 	return nil
// }

// func (check *Check) SetSignature(sig []byte) {
// 	check.R = new(big.Int).SetBytes(sig[:32])
// 	check.S = new(big.Int).SetBytes(sig[32:64])
// 	check.V = new(big.Int).SetBytes([]byte{sig[64] + 27})
// }

// func (check *Check) String() string {
// 	sender, _ := check.Sender()

// 	return fmt.Sprintf("Check sender: %s nonce: %x, dueBlock: %d, value: %s %s", sender.String(), check.Nonce,
// 		check.DueBlock, check.Amount.String(), check.Coin)
// 	// return fmt.Sprintf("Check nonce: %x, dueBlock: %d, value: %s %s", check.Nonce, check.DueBlock, check.Amount.String(), check.Coin)
// }

// func ParseCheck(buf []byte) (*Check, error) {
// 	var check Check
// 	err := rlp.Decode(bytes.NewReader(buf), &check)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if check.S == nil || check.R == nil || check.V == nil {
// 		return nil, errors.New("incorrect tx signature")
// 	}

// 	return &check, nil
// }

// func rlpHash(x interface{}) (h Hash) {
// 	hw := sha3.NewLegacyKeccak256()
// 	err := rlp.Encode(hw, x)
// 	if err != nil {
// 		panic(err)
// 	}
// 	hw.Sum(h[:0])
// 	return h
// }

// func recoverPlain(sighash Hash, R, S, Vb *big.Int) (sdk.AccAddress, error) {
// 	if Vb.BitLen() > 8 {
// 		return sdk.AccAddress{}, ErrInvalidSig
// 	}
// 	V := byte(Vb.Uint64() - 27)
// 	if !crypto.ValidateSignatureValues(V, R, S, true) {
// 		return sdk.AccAddress{}, ErrInvalidSig
// 	}
// 	// encode the snature in uncompressed format
// 	r, s := R.Bytes(), S.Bytes()
// 	sig := make([]byte, 65)
// 	copy(sig[32-len(r):32], r)
// 	copy(sig[64-len(s):64], s)
// 	sig[64] = V
// 	// recover the public key from the snature
// 	pub, err := crypto.Ecrecover(sighash[:], sig)
// 	if err != nil {
// 		return sdk.AccAddress{}, err
// 	}
// 	if len(pub) == 0 || pub[0] != 4 {
// 		return sdk.AccAddress{}, errors.New("invalid public key")
// 	}
// 	pub2, err := crypto.UnmarshalPubkey(pub)
// 	if err != nil {
// 		return sdk.AccAddress{}, err
// 	}
// 	pub3 := crypto.CompressPubkey(pub2)
// 	hasherSHA256 := sha256.New()
// 	hasherSHA256.Write(pub3)
// 	sha := hasherSHA256.Sum(nil)
// 	hasherRIPEMD160 := ripemd160.New()
// 	hasherRIPEMD160.Write(sha)
// 	return sdk.AccAddress(hasherRIPEMD160.Sum(nil)), nil
// }
