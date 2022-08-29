package types

import (
	"crypto/sha256"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereumCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/assert"
)

func TestParseCheck(t *testing.T) {
	var (
		chain_id             = "soki_9000-2"
		coinAmount, _        = sdk.ParseCoinNormalized("10000del")
		nonce, _             = sdk.NewIntFromString("9")
		dueBlock      uint64 = 123
		pass                 = ""
	)

	priv, _ := ethsecp256k1.GenerateKey()

	passphraseHash := sha256.Sum256([]byte(pass))
	passphrasePrivKey, _ := ethereumCrypto.ToECDSA(passphraseHash[:])

	check := &Check{
		ChainID:  chain_id,
		Coin:     coinAmount.Denom,
		Amount:   coinAmount.Amount,
		Nonce:    nonce.BigInt().Bytes(),
		DueBlock: dueBlock,
	}

	checkHash := check.HashWithoutLock()
	lock, _ := ethereumCrypto.Sign(checkHash[:], passphrasePrivKey)
	check.Lock = lock

	// un armor key

	key, err := priv.ToECDSA()
	assert.NoError(t, err)

	checkHash = check.Hash()
	signature, err := ethereumCrypto.Sign(checkHash[:], key)
	assert.NoError(t, err)

	check.SetSignature(signature)

	checkBytes, err := rlp.EncodeToBytes(check)
	assert.NoError(t, err)

	//decode

	decodedCheck, err := ParseCheck(checkBytes)
	assert.NoError(t, err)

	assert.True(t, check.Equal(decodedCheck))

	sender, err := check.Sender()
	assert.NoError(t, err)

	senderFromPriv := sdk.AccAddress(priv.PubKey().Address().Bytes())

	assert.Equal(t, check.Hash(), decodedCheck.Hash(), "hash not equal")
	assert.Equal(t, check.HashFull(), decodedCheck.HashFull(), "fullHash not equal")
	assert.True(t, senderFromPriv.Equals(sender), "Addresses not equal")
}
