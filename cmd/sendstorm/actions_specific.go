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

// MsgCreateCoin
type CreateCoinGenerator struct {
	// general values
	symbolLengthBottom, symbolLengthUp int64
	initVolumeBottom, initVolumeUp     int64 // in 10^18
	initReserveBottom, initReserveUp   int64 // in 10^18
	limitVolumeBottom, limitVolumeUp   int64 // in 10^18
	knownCoins                         []string
	rnd                                *rand.Rand
}

type CreateCoinAction struct {
	title       string
	symbol      string
	crr         uint64
	initVolume  sdk.Int
	initReserve sdk.Int
	limitVolume sdk.Int
	identity    string
}

func NewCreateCoinGenerator(
	symbolLengthBottom, symbolLengthUp,
	initVolumeBottom, initVolumeUp,
	initReserveBottom, initReserveUp,
	limitVolumeBottom, limitVolumeUp int64) *CreateCoinGenerator {
	return &CreateCoinGenerator{
		symbolLengthBottom: symbolLengthBottom,
		symbolLengthUp:     symbolLengthUp,
		initVolumeBottom:   initVolumeBottom,
		initVolumeUp:       initVolumeUp,
		initReserveBottom:  initReserveBottom,
		initReserveUp:      initReserveUp,
		limitVolumeBottom:  limitVolumeBottom,
		limitVolumeUp:      limitVolumeUp,
		rnd:                rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (acg *CreateCoinGenerator) Update(ui UpdateInfo) {
	acg.knownCoins = ui.Coins
}

func (acg *CreateCoinGenerator) Generate() Action {
	var symbol string
	doContinue := true
	for doContinue {
		symbol = randomString(acg.rnd, randomRange(acg.rnd, acg.symbolLengthBottom, acg.symbolLengthUp), charsAbc)
		doContinue = false
		for _, s := range acg.knownCoins {
			if symbol == s {
				doContinue = true
			}
		}
	}

	return &CreateCoinAction{
		title:       randomString(acg.rnd, 10, charsAll),
		symbol:      symbol,
		crr:         uint64(randomRange(acg.rnd, 10, 100+1)),
		initVolume:  helpers.EtherToWei(sdk.NewInt(randomRange(acg.rnd, acg.initVolumeBottom, acg.initVolumeUp))),
		initReserve: helpers.EtherToWei(sdk.NewInt(randomRange(acg.rnd, acg.initReserveBottom, acg.initReserveUp))),
		limitVolume: helpers.EtherToWei(sdk.NewInt(randomRange(acg.rnd, acg.limitVolumeBottom, acg.limitVolumeUp))),
		identity:    randomString(acg.rnd, 10, charsAll),
	}
}

func (ac *CreateCoinAction) CanPerform(sa *StormAccount) bool {
	if sa.IsDirty() {
		return false
	}
	if sa.BalanceForCoin(sa.feeDenom).LT(ac.initReserve) {
		return false
	}
	return true
}

func (ac *CreateCoinAction) GenerateTx(sa *StormAccount) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgCreateCoin(
		sender,
		ac.title,
		ac.symbol,
		ac.crr,
		ac.initVolume,
		ac.initReserve,
		ac.limitVolume,
		ac.identity,
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
	coinName := randomChoice(abg.rnd, abg.knownCoins)
	for _, ci := range abg.knownFullCoins {
		if ci.Symbol == coinName {
			coinInfo = ci
			break
		}
	}
	amountToBuy := sdk.NewInt(randomRange(abg.rnd, abg.bottomRange, abg.upperRange))
	amountToSell := formulas.CalculatePurchaseAmount(coinInfo.Volume, coinInfo.Reserve, uint(coinInfo.CRR), amountToBuy)

	return &BuyCoinAction{
		coinToBuy: sdk.NewCoin(
			coinName,
			helpers.FinneyToWei(amountToBuy),
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
