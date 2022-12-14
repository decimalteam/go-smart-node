package tx

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"

	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// Issue check for redeem
func IssueCheck(acc *dscWallet.Account, denom string, amount math.Int, nonce math.Int, dueBlock uint64, passphrase string) (string, error) {
	coin := sdk.NewCoin(denom, amount)

	// Prepare private key from passphrase
	passphraseHash := sha256.Sum256([]byte(passphrase))
	passphrasePrivKey, err := crypto.ToECDSA(passphraseHash[:])
	if err != nil {
		return "", fmt.Errorf("unable to create private key from passphrase")
	}

	// Prepare check without lock
	check := &types.Check{
		ChainID:  acc.ChainID(),
		Coin:     coin,
		Nonce:    nonce.BigInt().Bytes(),
		DueBlock: dueBlock,
	}

	// Prepare check lock
	checkHash := check.HashWithoutLock()
	lock, err := crypto.Sign(checkHash[:], passphrasePrivKey)
	if err != nil {
		return "", fmt.Errorf("unable to sign check hash")
	}

	// Fill check with prepared lock
	check.Lock = lock

	// Sign check by check issuer
	checkHash = check.Hash()
	signature, err := acc.Sign(checkHash[:])
	if err != nil {
		return "", fmt.Errorf("unable to sign check")
	}

	check.SetSignature(signature)

	checkBytes, err := rlp.EncodeToBytes(check)
	if err != nil {
		return "", fmt.Errorf("unable to encode check")
	}

	return base58.Encode(checkBytes), nil
}

func CreateRedeemCheck(acc *dscWallet.Account, checkBase58 string, passphrase string) (*MsgRedeemCheck, error) {
	// Decode provided check from base58 format to raw bytes
	checkBytes := base58.Decode(checkBase58)
	if len(checkBytes) == 0 {
		return nil, fmt.Errorf("unable to decode check from base58")
	}

	// Parse provided check from raw bytes to ensure it is valid
	_, err := types.ParseCheck(checkBytes)
	if err != nil {
		return nil, err
	}

	// Prepare private key from passphrase
	passphraseHash := sha256.Sum256([]byte(passphrase))
	passphrasePrivKey, err := crypto.ToECDSA(passphraseHash[:])
	if err != nil {
		return nil, fmt.Errorf("unable to create private key from passphrase")
	}

	// Prepare bytes to sign by private key generated from passphrase
	receiverAddressHash := make([]byte, 32)
	hw := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hw, []interface{}{
		acc.SdkAddress(),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to encode address in rpl")
	}
	hw.Sum(receiverAddressHash[:0])

	// Sign receiver address by private key generated from passphrase
	signature, err := crypto.Sign(receiverAddressHash, passphrasePrivKey)
	if err != nil {
		return nil, fmt.Errorf("unable to sign check")
	}
	proofBase64 := base64.StdEncoding.EncodeToString(signature)

	// Prepare redeem check message
	return types.NewMsgRedeemCheck(acc.SdkAddress(), checkBase58, proofBase64), nil
}
