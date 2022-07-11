package main

import (
	"math/rand"
	"time"

	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

func (ac *CreateCoinAction) ChooseAccounts(saList []*StormAccount) []*StormAccount {
	var res []*StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].BalanceForCoin(saList[i].feeDenom).LT(ac.initReserve) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
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
	tx, err := dscTx.BuildTransaction(sa.Account(), []sdk.Msg{msg}, "", sa.FeeDenom())
	if err != nil {
		return nil, err
	}
	err = tx.SignTransaction(sa.Account())
	if err != nil {
		return nil, err
	}
	return tx.BytesToSend()
}

func (ac *CreateCoinAction) String() string {
	return ""
}
