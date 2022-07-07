package main

import (
	"math/rand"
	"time"

	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgSendCoin
type SendCoinGenerator struct {
	// general values
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	knownCoins              []string
	knownAddresses          []string
	rnd                     *rand.Rand
}

type SendCoinAction struct {
	coin    sdk.Coin
	address string
}

func NewSendCoinGenerator(bottomRange, upperRange int64) *SendCoinGenerator {
	return &SendCoinGenerator{
		bottomRange: bottomRange,
		upperRange:  upperRange,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (asg *SendCoinGenerator) Update(ui UpdateInfo) {
	asg.knownCoins = ui.Coins
	asg.knownAddresses = ui.Addresses
}

func (asg *SendCoinGenerator) Generate() Action {
	return &SendCoinAction{
		coin: sdk.NewCoin(
			randomChoice(asg.rnd, asg.knownCoins),
			helpers.FinneyToWei(sdk.NewInt(randomRange(asg.rnd, asg.bottomRange, asg.upperRange))),
		),
		address: randomChoice(asg.rnd, asg.knownAddresses)}
}

func (as *SendCoinAction) CanPerform(sa *StormAccount) bool {
	if sa.IsDirty() {
		return false
	}
	if sa.BalanceForCoin(as.coin.Denom).LT(as.coin.Amount) {
		return false
	}
	if sa.Address() == as.address {
		return false
	}
	return true
}

func (as *SendCoinAction) GenerateTx(sa *StormAccount) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}
	receiver, err := sdk.AccAddressFromBech32(as.address)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgSendCoin(sender, as.coin, receiver)
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
