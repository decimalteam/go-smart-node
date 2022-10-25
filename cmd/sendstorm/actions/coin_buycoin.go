package actions

import (
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"

	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
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
	var amountToSell sdkmath.Int
	var coinName string
	if len(abg.knownCoins) <= 1 {
		return &EmptyAction{}
	}
	for {
		coinName = RandomChoice(abg.rnd, abg.knownCoins)
		if coinName != abg.baseCoin {
			break
		}
	}

	for _, ci := range abg.knownFullCoins {
		if ci.Denom == coinName {
			coinInfo = ci
			break
		}
	}
	amountToBuy := helpers.FinneyToWei(sdkmath.NewInt(RandomRange(abg.rnd, abg.bottomRange, abg.upperRange)))
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

func (ab *BuyCoinAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].BalanceForCoin(saList[i].FeeDenom()).LT(ab.maxCoinToSell.Amount) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (ab *BuyCoinAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgBuyCoin(
		sender,
		ab.coinToBuy,
		ab.maxCoinToSell,
	)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ab *BuyCoinAction) String() string {
	return fmt.Sprintf("BuyCoin{ coinToBuy: %s,  maxCoinToSell: %s}", ab.coinToBuy.String(), ab.maxCoinToSell.String())
}
