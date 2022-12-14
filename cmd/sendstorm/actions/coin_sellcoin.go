package actions

import (
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
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
	var amountToSell sdkmath.Int
	var coinName string
	if len(asg.knownCoins) == 1 {
		return &EmptyAction{}
	}
	for {
		coinName = RandomChoice(asg.rnd, asg.knownCoins)
		if coinName != asg.baseCoin {
			break
		}
	}

	for _, ci := range asg.knownFullCoins {
		if ci.Denom == coinName {
			coinInfo = ci
			break
		}
	}
	amountToSell = helpers.FinneyToWei(sdkmath.NewInt(RandomRange(asg.rnd, asg.bottomRange, asg.upperRange)))
	// respect limit volume to decrease amount of errors
	if coinInfo.Volume.Sub(amountToSell).LT(sdkmath.ZeroInt()) {
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

func (as *SellCoinAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].BalanceForCoin(as.coinToSell.Denom).LT(as.coinToSell.Amount) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (as *SellCoinAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgSellCoin(
		sender,
		as.coinToSell,
		as.minCoinToBuy,
	)

	return feeConfig.MakeTransaction(sa, msg)
}

func (as *SellCoinAction) String() string {
	return fmt.Sprintf("SellCoin{coinToSell:%s, minCoinToBuy:%s}", as.coinToSell.String(), as.minCoinToBuy.String())
}
