package actions

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgUpdateCoin
type UpdateCoinGenerator struct {
	// ranges for new limit delta
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	// need full coin info to calculate price
	knownCoins     []string
	knownFullCoins []dscApi.Coin
	baseCoin       string
	rnd            *rand.Rand
}

type UpdateCoinAction struct {
	creator        sdk.AccAddress
	symbol         string
	newLimitVolume sdkmath.Int
	newIdentity    string
}

func NewUpdateCoinGenerator(
	bottomRange, upperRange int64,
	baseCoin string) *UpdateCoinGenerator {
	return &UpdateCoinGenerator{
		bottomRange: bottomRange,
		upperRange:  upperRange,
		baseCoin:    baseCoin,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (aug *UpdateCoinGenerator) Update(ui UpdateInfo) {
	aug.knownCoins = ui.Coins
	aug.knownFullCoins = ui.FullCoins
}

func (aug *UpdateCoinGenerator) Generate() Action {
	var coinInfo dscApi.Coin
	var coinName string
	if len(aug.knownCoins) == 1 {
		return &EmptyAction{}
	}
	for {
		coinName = RandomChoice(aug.rnd, aug.knownCoins)
		if coinName != aug.baseCoin {
			break
		}
	}
	for _, ci := range aug.knownFullCoins {
		if ci.Symbol == coinName {
			coinInfo = ci
			break
		}
	}

	delta := helpers.FinneyToWei(sdk.NewInt(RandomRange(aug.rnd, aug.bottomRange, aug.upperRange)))

	creator, err := sdk.AccAddressFromBech32(coinInfo.Creator)
	if err != nil {
		return &EmptyAction{}
	}

	return &UpdateCoinAction{
		creator:        creator,
		symbol:         coinName,
		newLimitVolume: coinInfo.LimitVolume.Add(delta),
		newIdentity:    RandomString(aug.rnd, 30, charsAll),
	}
}

func (au *UpdateCoinAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if !saList[i].Account().SdkAddress().Equals(au.creator) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (au *UpdateCoinAction) GenerateTx(sa *stormTypes.StormAccount) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}
	msg := dscTx.NewMsgUpdateCoin(sender, au.symbol, au.newLimitVolume, au.newIdentity)
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

func (au *UpdateCoinAction) String() string {
	return fmt.Sprintf("UpdateCoin{Symbol: %s, new limit: %s, new idenity: %s}",
		au.symbol, au.newLimitVolume.String(), au.newIdentity)
}
