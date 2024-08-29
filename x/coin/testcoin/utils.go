package testcoin

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/decimalteam/ethermint/crypto/ethsecp256k1"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

func CreateNewCheck(chainID, coinStr, nonceStr, password string, dueBlock uint64) types.Check {
	var (
		coin, _  = sdk.ParseCoinNormalized(coinStr)
		nonce, _ = sdk.NewIntFromString(nonceStr)
	)

	priv, _ := ethsecp256k1.GenerateKey()

	passphraseHash := sha256.Sum256([]byte(password))
	passphrasePrivKey, _ := crypto.ToECDSA(passphraseHash[:])

	check := &types.Check{
		ChainID:  chainID,
		Coin:     coin,
		Nonce:    nonce.BigInt().Bytes(),
		DueBlock: dueBlock,
	}

	checkHash := check.HashWithoutLock()
	lock, _ := crypto.Sign(checkHash[:], passphrasePrivKey)
	check.Lock = lock

	// un armor key
	key, _ := priv.ToECDSA()

	checkHash = check.Hash()
	signature, _ := crypto.Sign(checkHash[:], key)

	check.SetSignature(signature)

	return *check
}

func ChecksEqual(whichChecks, withChecks types.Checks) bool {
	if len(whichChecks) != len(withChecks) {
		return false
	}

	withChecksMap := make(map[[32]byte]types.Check)
	for _, v := range withChecks {
		withChecksMap[v.HashFull()] = v
	}

	for _, v := range whichChecks {
		with, ok := withChecksMap[v.HashFull()]
		if !ok {
			return false
		}

		if !v.Equal(with) {
			return false
		}
	}

	return true
}

func CoinsEqual(whichCoins, withCoins []types.Coin) bool {
	if len(whichCoins) != len(withCoins) {
		return false
	}

	withCoinsMap := make(map[string]types.Coin)
	for _, v := range withCoins {
		withCoinsMap[v.Denom] = v
	}

	for _, v := range whichCoins {
		with, ok := withCoinsMap[v.Denom]
		if !ok {
			return false
		}

		if !v.Equal(with) {
			return false
		}
	}

	return true
}
