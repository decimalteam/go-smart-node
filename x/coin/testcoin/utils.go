package testcoin

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"crypto/sha256"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereumCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

func CreateNewCheck(chainID, coinAmountStr, nonceStr, password string, dueBlock uint64) types.Check {
	var (
		coinAmount, _ = sdk.ParseCoinNormalized(coinAmountStr)
		nonce, _      = sdk.NewIntFromString(nonceStr)
	)

	priv, _ := ethsecp256k1.GenerateKey()

	passphraseHash := sha256.Sum256([]byte(password))
	passphrasePrivKey, _ := ethereumCrypto.ToECDSA(passphraseHash[:])

	check := &types.Check{
		ChainID:  chainID,
		Coin:     coinAmount.Denom,
		Amount:   coinAmount.Amount,
		Nonce:    nonce.BigInt().Bytes(),
		DueBlock: dueBlock,
	}

	checkHash := check.HashWithoutLock()
	lock, _ := ethereumCrypto.Sign(checkHash[:], passphrasePrivKey)
	check.Lock = lock

	// un armor key
	key, _ := priv.ToECDSA()

	checkHash = check.Hash()
	signature, _ := ethereumCrypto.Sign(checkHash[:], key)

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

func CoinsEqual(whichCoins, withCoins types.Coins) bool {
	if len(whichCoins) != len(withCoins) {
		return false
	}

	withCoinsMap := make(map[string]types.Coin)
	for _, v := range withCoins {
		withCoinsMap[v.Symbol] = v
	}

	for _, v := range whichCoins {
		with, ok := withCoinsMap[v.Symbol]
		if !ok {
			return false
		}

		if !v.Equal(with) {
			return false
		}
	}

	return true
}
