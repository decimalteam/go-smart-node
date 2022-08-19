package actions

import (
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgBurnCoin
type BurnCoinGenerator struct {
	// general values
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	knownCoins              []string
	rnd                     *rand.Rand
}

type BurnCoinAction struct {
	coin sdk.Coin
}

func NewBurnCoinGenerator(bottomRange, upperRange int64) *BurnCoinGenerator {
	return &BurnCoinGenerator{
		bottomRange: bottomRange,
		upperRange:  upperRange,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *BurnCoinGenerator) Update(ui UpdateInfo) {
	gg.knownCoins = ui.Coins
}

func (gg *BurnCoinGenerator) Generate() Action {
	return &BurnCoinAction{
		coin: sdk.NewCoin(
			RandomChoice(gg.rnd, gg.knownCoins),
			helpers.FinneyToWei(sdk.NewInt(RandomRange(gg.rnd, gg.bottomRange, gg.upperRange))),
		),
	}
}

func (aa *BurnCoinAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].BalanceForCoin(aa.coin.Denom).LT(aa.coin.Amount) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *BurnCoinAction) GenerateTx(sa *stormTypes.StormAccount) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgBurnCoin(sender, aa.coin)
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

func (aa *BurnCoinAction) String() string {
	return fmt.Sprintf("BurnCoin{coin: %s}", aa.coin.String())
}