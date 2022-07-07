package main

import (
	"math/rand"
	"time"

	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgBuyCoin
type BuyCoinGenerator struct {
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	// need full coin info to calculate price
	knownCoins     []string
	knownFullCoins []dscApi.Coin
	baseCoin       string
	rnd            *rand.Rand
}

type BuyCoinAction struct {
	coinToBuy     sdk.Coin
	maxCoinToSell sdk.Coin
}

func NewBuyCoinGenerator(
	bottomRange, upperRange int64,
	baseCoin string) *BuyCoinGenerator {
	return &BuyCoinGenerator{
		bottomRange: bottomRange,
		upperRange:  upperRange,
		baseCoin:    baseCoin,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (abg *BuyCoinGenerator) Update(ui UpdateInfo) {
	abg.knownFullCoins = ui.FullCoins
	abg.knownCoins = ui.Coins
}

func (abg *BuyCoinGenerator) Generate() Action {
	var coinInfo dscApi.Coin
	var amountToSell sdk.Int
	var coinName string
	if len(abg.knownCoins) == 1 {
		return &EmptyAction{}
	}
	for {
		coinName = randomChoice(abg.rnd, abg.knownCoins)
		if coinName != abg.baseCoin {
			break
		}
	}

	for _, ci := range abg.knownFullCoins {
		if ci.Symbol == coinName {
			coinInfo = ci
			break
		}
	}
	amountToBuy := helpers.FinneyToWei(sdk.NewInt(randomRange(abg.rnd, abg.bottomRange, abg.upperRange)))
	amountToSell = formulas.CalculatePurchaseAmount(coinInfo.Volume, coinInfo.Reserve, uint(coinInfo.CRR), amountToBuy)
	// respect limit volume to decrease amount of errors
	if coinInfo.Volume.Add(amountToBuy).GT(coinInfo.LimitVolume) {
		return &EmptyAction{}
	}

	return &BuyCoinAction{
		coinToBuy: sdk.NewCoin(
			coinName,
			amountToBuy,
		),
		maxCoinToSell: sdk.NewCoin(
			abg.baseCoin,
			amountToSell,
		),
	}
}

func (ab *BuyCoinAction) CanPerform(sa *StormAccount) bool {
	if sa.IsDirty() {
		return false
	}
	if sa.BalanceForCoin(sa.feeDenom).LT(ab.maxCoinToSell.Amount) {
		return false
	}
	return true
}

func (ab *BuyCoinAction) GenerateTx(sa *StormAccount) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgBuyCoin(
		sender,
		ab.coinToBuy,
		ab.maxCoinToSell,
	)
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
