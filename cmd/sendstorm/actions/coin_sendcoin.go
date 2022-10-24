package actions

import (
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
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
			RandomChoice(asg.rnd, asg.knownCoins),
			helpers.FinneyToWei(sdkmath.NewInt(RandomRange(asg.rnd, asg.bottomRange, asg.upperRange))),
		),
		address: RandomChoice(asg.rnd, asg.knownAddresses)}
}

func (as *SendCoinAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].BalanceForCoin(as.coin.Denom).LT(as.coin.Amount) {
			continue
		}
		if saList[i].Address() == as.address {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (as *SendCoinAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}
	recipient, err := sdk.AccAddressFromBech32(as.address)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgSendCoin(sender, recipient, as.coin)

	return feeConfig.MakeTransaction(sa, msg)
}

func (as *SendCoinAction) String() string {
	return fmt.Sprintf("SendCoin{address: %s, coin: %s}", as.address, as.coin.String())
}
