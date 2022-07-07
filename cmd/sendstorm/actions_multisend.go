package main

import (
	"math/rand"
	"time"

	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgMultiSendCoin
type MultiSendCoinGenerator struct {
	// general values
	bottomRange, upperRange         int64 // bounds in 0.001 (10^15)
	sendCountBottom, sendCountUpper int64
	knownCoins                      []string
	knownAddresses                  []string
	rnd                             *rand.Rand
}

type MultiSendCoinAction struct {
	sends   []dscTx.OneSend
	summary sdk.Coins // for fast check
}

func NewMultiSendCoinGenerator(bottomRange, upperRange int64) *MultiSendCoinGenerator {
	return &MultiSendCoinGenerator{
		bottomRange: bottomRange,
		upperRange:  upperRange,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (asg *MultiSendCoinGenerator) Update(ui UpdateInfo) {
	asg.knownCoins = ui.Coins
	asg.knownAddresses = ui.Addresses
}

func (asg *MultiSendCoinGenerator) Generate() Action {
	n := randomRange(asg.rnd, asg.sendCountBottom, asg.sendCountUpper)
	sums := sdk.NewCoins()
	sends := make([]dscTx.OneSend, n)
	for i := int64(0); i < n; i++ {
		coin := sdk.NewCoin(
			randomChoice(asg.rnd, asg.knownCoins),
			helpers.FinneyToWei(sdk.NewInt(randomRange(asg.rnd, asg.bottomRange, asg.upperRange))),
		)
		sums = sums.Add(coin)
		sends[i] = dscTx.OneSend{
			Coin:     coin,
			Receiver: randomChoice(asg.rnd, asg.knownAddresses),
		}
	}
	return &MultiSendCoinAction{
		sends:   sends,
		summary: sums}
}

func (as *MultiSendCoinAction) CanPerform(sa *StormAccount) bool {
	if sa.IsDirty() {
		return false
	}
	for _, coin := range as.summary {
		if sa.BalanceForCoin(coin.Denom).LT(coin.Amount) {
			return false
		}
	}
	for _, send := range as.sends {
		if sa.Address() == send.Receiver {
			return false
		}
	}
	return true
}

func (as *MultiSendCoinAction) GenerateTx(sa *StormAccount) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgMultiSendCoin(sender, as.sends)
	tx, err := dscTx.BuildTransaction([]sdk.Msg{msg}, "", sa.FeeDenom(), sa.MaxGas())
	if err != nil {
		return nil, err
	}
	tx, err = tx.SignTransaction(sa.Account())
	if err != nil {
		return nil, err
	}
	return tx.BytesToSend()
}
