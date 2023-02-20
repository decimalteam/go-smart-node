package actions

import (
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
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
	initVolume  sdkmath.Int
	initReserve sdkmath.Int
	limitVolume sdkmath.Int
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
		symbol = RandomString(acg.rnd, RandomRange(acg.rnd, acg.symbolLengthBottom, acg.symbolLengthUp), charsAbc)
		doContinue = false
		for _, s := range acg.knownCoins {
			if symbol == s {
				doContinue = true
			}
		}
	}

	return &CreateCoinAction{
		title:       RandomString(acg.rnd, 10, charsAll),
		symbol:      symbol,
		crr:         uint64(RandomRange(acg.rnd, 10, 100+1)),
		initVolume:  helpers.EtherToWei(sdkmath.NewInt(RandomRange(acg.rnd, acg.initVolumeBottom, acg.initVolumeUp))),
		initReserve: helpers.EtherToWei(sdkmath.NewInt(RandomRange(acg.rnd, acg.initReserveBottom, acg.initReserveUp))),
		limitVolume: helpers.EtherToWei(sdkmath.NewInt(RandomRange(acg.rnd, acg.limitVolumeBottom, acg.limitVolumeUp))),
		identity:    RandomString(acg.rnd, 10, charsAll),
	}
}

func (ac *CreateCoinAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].BalanceForCoin(saList[i].FeeDenom()).LT(ac.initReserve) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (ac *CreateCoinAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgCreateCoin(
		sender,
		ac.symbol,
		ac.title,
		ac.crr,
		ac.initVolume,
		ac.initReserve,
		ac.limitVolume,
		sdkmath.ZeroInt(),
		ac.identity,
	)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *CreateCoinAction) String() string {
	return fmt.Sprintf("CreateCoin{title: %s, symbol: %s, identity: %s, crr: %d, init volume: %s, init reserve: %s, limit volume: %s}",
		ac.title, ac.symbol, ac.identity, ac.crr,
		ac.initVolume.String(), ac.initReserve.String(), ac.limitVolume.String())
}
