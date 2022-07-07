package main

import (
	"math/rand"
	"time"

	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgSellCoin
type SellCoinGenerator struct {
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	// need full coin info to calculate price
	knownCoins     []string
	knownFullCoins []dscApi.Coin
	baseCoin       string
	rnd            *rand.Rand
}

type SellCoinAction struct {
	coinToSell   sdk.Coin
	minCoinToBuy sdk.Coin
}

func NewSellCoinGenerator(
	bottomRange, upperRange int64,
	baseCoin string) *SellCoinGenerator {
	return &SellCoinGenerator{
		bottomRange: bottomRange,
		upperRange:  upperRange,
		baseCoin:    baseCoin,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (asg *SellCoinGenerator) Update(ui UpdateInfo) {
	asg.knownFullCoins = ui.FullCoins
	asg.knownCoins = ui.Coins
}

func (asg *SellCoinGenerator) Generate() Action {
	var coinInfo dscApi.Coin
	var amountToSell sdk.Int
	var coinName string
	if len(asg.knownCoins) == 1 {
		return &EmptyAction{}
	}
	for {
		coinName = randomChoice(asg.rnd, asg.knownCoins)
		if coinName != asg.baseCoin {
			break
		}
	}

	// TODO: random coin for fee

	for _, ci := range asg.knownFullCoins {
		if ci.Symbol == coinName {
			coinInfo = ci
			break
		}
	}
	amountToSell = helpers.FinneyToWei(sdk.NewInt(randomRange(asg.rnd, asg.bottomRange, asg.upperRange)))
	// respect limit volume to decrease amount of errors
	if coinInfo.Volume.Sub(amountToSell).LT(sdk.ZeroInt()) {
		return &EmptyAction{}
	}

	return &SellCoinAction{
		coinToSell: sdk.NewCoin(
			coinName,
			amountToSell,
		),
		minCoinToBuy: sdk.NewCoin(
			asg.baseCoin,
			sdk.ZeroInt(),
		),
	}
}

func (as *SellCoinAction) CanPerform(sa *StormAccount) bool {
	if sa.IsDirty() {
		return false
	}
	if sa.BalanceForCoin(as.coinToSell.Denom).LT(as.coinToSell.Amount) {
		return false
	}
	return true
}

func (as *SellCoinAction) GenerateTx(sa *StormAccount) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgSellCoin(
		sender,
		as.coinToSell,
		as.minCoinToBuy,
	)
	// TODO: fee in custom coin
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
