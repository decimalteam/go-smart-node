package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetHash(transactionNumber math.Int, tokenSymbol string, amount math.Int, recipient string,
	fromChain, destChain uint32) (Hash, error) {
	var hash [32]byte

	encoded := encodePacked(
		encodeUint256(transactionNumber.BigInt()),
		encodeUint256(amount.BigInt()),
		encodeString(tokenSymbol),
		encodeString(recipient),
		encodeUint256(sdk.NewInt(int64(fromChain)).BigInt()),
		encodeUint256(sdk.NewInt(int64(destChain)).BigInt()),
	)

	copy(hash[:], crypto.Keccak256(encoded))

	fmt.Println(hex.EncodeToString(hash[:]))

	copy(hash[:], crypto.Keccak256(encodePacked(
		encodeString(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(hash))),
		hash[:])))

	fmt.Println(hex.EncodeToString(hash[:]))

	return hash, nil
}

func encodePacked(input ...[]byte) []byte {
	return bytes.Join(input, nil)
}

func encodeString(v string) []byte {
	return []byte(v)
}

func encodeUint256(v *big.Int) []byte {
	return ethmath.U256Bytes(v)
}

func encodeUint8(v uint8) []byte {
	return new(big.Int).SetUint64(uint64(v)).Bytes()
}

func Ecrecover(sighash [32]byte, R, S, Vb *big.Int) (ethcmn.Address, error) {
	if Vb.BitLen() > 8 {
		return ethcmn.Address{}, errors.New("invalid sig")
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, true) {
		return ethcmn.Address{}, errors.New("invalid sig")
	}
	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return ethcmn.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return ethcmn.Address{}, errors.New("invalid public key")
	}
	var addr ethcmn.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
	//if vb.BitLen() > 8 {
	//	return ethcmn.Address{}, errors.New("invalid sig") // errors.InvalidCheckSig
	//}
	//v := byte(vb.Uint64() - 27)
	//if !crypto.ValidateSignatureValues(v, rb, sb, true) {
	//	return ethcmn.Address{}, errors.New("invalid sig") //errors.InvalidCheckSig
	//}
	//// encode the signature in uncompressed format
	//r, s := rb.Bytes(), sb.Bytes()
	//sig := make([]byte, 65)
	//copy(sig[32-len(r):32], r)
	//copy(sig[64-len(s):64], s)
	//sig[64] = v
	//// recover the public key from the signature
	//pub, err := crypto.Ecrecover(sighash[:], sig)
	//if err != nil {
	//	return ethcmn.Address{}, errors.New("invalid sig") //errors.FailedToRecoverPKFromSig
	//}
	//if len(pub) == 0 || pub[0] != 4 {
	//	return ethcmn.Address{}, errors.New("invalid sig") //errors.InvalidPubKey
	//}
	//// calculate address from the recovered public key
	//var addr ethcmn.Address
	//copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	//return addr, nil

	//if Vb.BitLen() > 8 {
	//	return ethcmn.Address{}, errors.New("invalid sig")
	//}
	//V := byte(Vb.Uint64() - 27)
	//if !crypto.ValidateSignatureValues(V, R, S, true) {
	//	return ethcmn.Address{}, errors.New("invalid sig")
	//}
	//// encode the signature in uncompressed format
	//r, s := R.Bytes(), S.Bytes()
	//sig := make([]byte, 65)
	//copy(sig[32-len(r):32], r)
	//copy(sig[64-len(s):64], s)
	//sig[64] = V
	//
	//// recover the public key from the signature
	//pub, err := crypto.SigToPub(sighash[:], sig)
	//if err != nil {
	//	return ethcmn.Address{}, err
	//}
	//
	//addr := crypto.PubkeyToAddress(*pub)
	//
	//return addr, nil
}
