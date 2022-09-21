package types

import (
	"crypto/ecdsa"
	"fmt"
	"io"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
)

// HashLength represents constant length of a hash which equals to 32 bytes.
const HashLength = 32

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

func ParseCheck(buf []byte) (check *Check, err error) {
	err = rlp.DecodeBytes(buf, &check)
	if err != nil {
		return nil, errors.DecodeRLP
	}
	if check.S.BigInt() == nil || check.R.BigInt() == nil || check.V.BigInt() == nil {
		err = errors.InvalidCheckSig
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
		return nil, errors.UnableRecoverLockPkey
	}
	if len(pub) == 0 || pub[0] != 4 {
		return nil, errors.InvalidPubKey
	}
	return pub, nil
}

func (c *Check) Sign(prv *ecdsa.PrivateKey) error {
	h := c.Hash()
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return errors.UnableSignCheck
	}
	c.SetSignature(sig)
	return nil
}

func (c *Check) SetSignature(sig []byte) {
	c.R = sdk.NewIntFromBigInt(new(big.Int).SetBytes(sig[:32]))
	c.S = sdk.NewIntFromBigInt(new(big.Int).SetBytes(sig[32:64]))
	c.V = sdk.NewIntFromBigInt(new(big.Int).SetBytes([]byte{sig[64] + 27}))
}

func (c *Check) String() string {
	sender, _ := c.Sender()
	return fmt.Sprintf(
		"Check sender: %s nonce: %x, dueBlock: %d, value: %s",
		sender, c.Nonce, c.DueBlock, c.Coin,
	)
}

func (c *Check) HashWithoutLock() Hash {
	return rlpHash([]interface{}{c.ChainID, c.Coin.Denom, c.Coin.Amount.BigInt(), c.Nonce, c.DueBlock})
}

func (c *Check) Hash() Hash {
	return rlpHash([]interface{}{c.ChainID, c.Coin.Denom, c.Coin.Amount.BigInt(), c.Nonce, c.DueBlock, c.Lock})
}

func (c *Check) HashFull() Hash {
	return rlpHash([]interface{}{c.ChainID, c.Coin.Denom, c.Coin.Amount.BigInt(), c.Nonce, c.DueBlock, c.Lock, c.V, c.R, c.S})
}

func (c *Check) EncodeRLP(w io.Writer) error {
	if err := rlp.Encode(w, rlpCheck{
		ChainID:  c.ChainID,
		Coin:     c.Coin.Denom,
		Amount:   c.Coin.Amount.BigInt(),
		Nonce:    c.Nonce,
		DueBlock: c.DueBlock,
		Lock:     c.Lock,
		V:        c.V.BigInt(),
		R:        c.R.BigInt(),
		S:        c.S.BigInt(),
	}); err != nil {
		return errors.UnableRLPEncodeCheck
	}
	return nil
}

func (c *Check) DecodeRLP(st *rlp.Stream) error {
	var result rlpCheck
	if err := st.Decode(&result); err != nil {
		return err
	}

	c.ChainID = result.ChainID
	c.Coin = sdk.NewCoin(result.Coin, sdkmath.NewIntFromBigInt(result.Amount))
	c.Nonce = result.Nonce
	c.DueBlock = result.DueBlock
	c.Lock = result.Lock
	c.V = sdk.NewIntFromBigInt(result.V)
	c.R = sdk.NewIntFromBigInt(result.R)
	c.S = sdk.NewIntFromBigInt(result.S)

	return nil
}

type Checks []Check

func (c Checks) String() string {
	result := make([]string, len(c))
	for i, v := range c {
		result[i] = v.String()
	}

	return strings.Join(result, "\n")
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

func rlpHash(x interface{}) (h Hash) {
	hw := sha3.NewLegacyKeccak256()
	err := rlp.Encode(hw, x)
	if err != nil {
		panic(err)
	}
	hw.Sum(h[:0])
	return h
}

func recoverPlain(h Hash, rb, sb, vb *big.Int) (sdk.AccAddress, error) {
	if vb.BitLen() > 8 {
		return sdk.AccAddress{}, errors.InvalidCheckSig
	}
	v := byte(vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(v, rb, sb, true) {
		return sdk.AccAddress{}, errors.InvalidCheckSig
	}
	// encode the signature in uncompressed format
	r, s := rb.Bytes(), sb.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = v
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(h[:], sig)
	if err != nil {
		return sdk.AccAddress{}, errors.FailedToRecoverPKFromSig
	}
	if len(pub) == 0 || pub[0] != 4 {
		return sdk.AccAddress{}, errors.InvalidPubKey
	}
	// calculate address from the recovered public key
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return sdk.AccAddress(addr.Bytes()), nil
}
