package types

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io"
	"math/big"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const HashLength = 32

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
)

func ParseCheck(buf []byte) (check *Check, err error) {
	err = rlp.DecodeBytes(buf, &check)
	if check.S.BigInt() == nil || check.R.BigInt() == nil || check.V.BigInt() == nil {
		err = errors.New("incorrect tx signature")
		return
	}
	return
}

func (c *Check) Sender() (sdk.AccAddress, error) {
	return recoverPlain(c.Hash(), c.R.BigInt(), c.S.BigInt(), c.V.BigInt())
}

func (c *Check) LockPubKey() ([]byte, error) {
	sig := c.Lock
	if len(sig) < 65 {
		sig = append(make([]byte, 65-len(sig)), sig...)
	}
	hash := c.HashWithoutLock()
	pub, err := crypto.Ecrecover(hash[:], sig)
	if err != nil {
		return nil, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return nil, errors.New("invalid public key")
	}
	return pub, nil
}

func (c *Check) Sign(prv *ecdsa.PrivateKey) error {
	h := c.Hash()
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return err
	}
	c.SetSignature(sig)
	return nil
}

func (c *Check) SetSignature(sig []byte) {
	r := sdk.NewIntFromBigInt(new(big.Int).SetBytes(sig[:32]))
	s := sdk.NewIntFromBigInt(new(big.Int).SetBytes(sig[32:64]))
	v := sdk.NewIntFromBigInt(new(big.Int).SetBytes([]byte{sig[64] + 27}))
	c.R, c.S, c.V = &r, &s, &v
}

func (c *Check) String() string {
	sender, _ := c.Sender()
	return fmt.Sprintf(
		"Check sender: %s nonce: %x, dueBlock: %d, value: %s%s",
		sender, c.Nonce, c.DueBlock, c.Amount, c.Coin,
	)
}

func (c *Check) HashWithoutLock() Hash {
	return rlpHash([]interface{}{c.ChainID, c.Coin, c.Amount.BigInt(), c.Nonce, c.DueBlock})
}

func (c *Check) Hash() Hash {
	return rlpHash([]interface{}{c.ChainID, c.Coin, c.Amount.BigInt(), c.Nonce, c.DueBlock, c.Lock})
}

func (c *Check) HashFull() Hash {
	return rlpHash([]interface{}{c.ChainID, c.Coin, c.Amount.BigInt(), c.Nonce, c.DueBlock, c.Lock, c.V, c.R, c.S})
}

func rlpHash(x interface{}) (h Hash) {
	hw := sha3.NewLegacyKeccak256()
	err := rlp.Encode(hw, x)
	if err != nil {
		panic(err)
	}
	hw.Sum(h[:0])
	return h
}

func recoverPlain(sighash Hash, rb, sb, vb *big.Int) (sdk.AccAddress, error) {
	if vb.BitLen() > 8 {
		return sdk.AccAddress{}, ErrInvalidSig
	}
	v := byte(vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(v, rb, sb, true) {
		return sdk.AccAddress{}, ErrInvalidSig
	}
	// encode the signature in uncompressed format
	r, s := rb.Bytes(), sb.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = v
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return sdk.AccAddress{}, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return sdk.AccAddress(addr.Bytes()), nil
}

type rlpCheck struct {
	ChainID  string
	Coin     string
	Amount   *big.Int
	Nonce    []byte
	DueBlock uint64
	Lock     []byte
	V        *big.Int
	R        *big.Int
	S        *big.Int
}

func (c *Check) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpCheck{
		ChainID:  c.ChainID,
		Coin:     c.Coin,
		Amount:   c.Amount.BigInt(),
		Nonce:    c.Nonce,
		DueBlock: c.DueBlock,
		Lock:     c.Lock,
		V:        c.V.BigInt(),
		R:        c.R.BigInt(),
		S:        c.S.BigInt(),
	})
}

func (c *Check) DecodeRLP(st *rlp.Stream) error {
	var result rlpCheck
	if err := st.Decode(&result); err != nil {
		return err
	}

	v, r, s := sdk.NewIntFromBigInt(result.V), sdk.NewIntFromBigInt(result.R), sdk.NewIntFromBigInt(result.S)

	c.ChainID = result.ChainID
	c.Coin = result.Coin
	c.Amount = sdk.NewIntFromBigInt(result.Amount)
	c.Nonce = result.Nonce
	c.DueBlock = result.DueBlock
	c.Lock = result.Lock
	c.V = &v
	c.R = &r
	c.S = &s

	return nil
}
