package actions

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgMultiSendCoin
type MultiSendCoinGenerator struct {
	// general values
	bottomRange, upperRange         int64 // bounds in 0.001 (10^15)
	sendCountBottom, sendCountUpper int64
	knownCoins                      []string
	knownAddresses                  []string
	rnd                             *rand.Rand
}

type MultiSendCoinAction struct {
	sends   []dscTx.MultiSendEntry
	summary sdk.Coins // for fast check
}

func NewMultiSendCoinGenerator(bottomRange, upperRange, sendCountBottom, sendCountUpper int64) *MultiSendCoinGenerator {
	return &MultiSendCoinGenerator{
		bottomRange:     bottomRange,
		upperRange:      upperRange,
		sendCountBottom: sendCountBottom,
		sendCountUpper:  sendCountUpper,
		rnd:             rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (asg *MultiSendCoinGenerator) Update(ui UpdateInfo) {
	asg.knownCoins = ui.Coins
	asg.knownAddresses = ui.Addresses
}

func (asg *MultiSendCoinGenerator) Generate() Action {
	n := RandomRange(asg.rnd, asg.sendCountBottom, asg.sendCountUpper)
	sums := sdk.NewCoins()
	sends := make([]dscTx.MultiSendEntry, n)
	for i := int64(0); i < n; i++ {
		coin := sdk.NewCoin(
			RandomChoice(asg.rnd, asg.knownCoins),
			helpers.FinneyToWei(sdkmath.NewInt(RandomRange(asg.rnd, asg.bottomRange, asg.upperRange))),
		)
		sums = sums.Add(coin)
		sends[i] = dscTx.MultiSendEntry{
			Recipient: RandomChoice(asg.rnd, asg.knownAddresses),
			Coin:      coin,
		}
	}
	return &MultiSendCoinAction{
		sends:   sends,
		summary: sums}
}

func (as *MultiSendCoinAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		doAdd := true
		for _, coin := range as.summary {
			if saList[i].BalanceForCoin(coin.Denom).LT(coin.Amount) {
				doAdd = false
			}
		}
		for _, send := range as.sends {
			if saList[i].Address() == send.Recipient {
				doAdd = false
			}
		}
		if !doAdd {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (as *MultiSendCoinAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgMultiSendCoin(sender, as.sends)

	return feeConfig.MakeTransaction(sa, msg)
}

func (as *MultiSendCoinAction) String() string {
	var sb strings.Builder
	sb.WriteString("MultiSendCoin {")
	for i, s := range as.sends {
		sb.WriteString(fmt.Sprintf("Send{Recipient: %s, Coin: %s}", s.Recipient, s.Coin.String()))
		if i < len(as.sends)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}
