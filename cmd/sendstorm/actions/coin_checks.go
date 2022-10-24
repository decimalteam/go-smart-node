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

type RedeemCheckGenerator struct {
	// general values
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	knownCoins              []string
	knownAddresses          []string
	rnd                     *rand.Rand
}

type RedeemCheckAction struct {
	coin     sdk.Coin
	issuer   string
	receiver string
	// need cache sender account because check need signature of issuer
	issuerAcc *stormTypes.StormAccount
}

func NewRedeemCheckGenerator(bottomRange, upperRange int64) *RedeemCheckGenerator {
	return &RedeemCheckGenerator{
		bottomRange: bottomRange,
		upperRange:  upperRange,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *RedeemCheckGenerator) Update(ui UpdateInfo) {
	gg.knownCoins = ui.Coins
	gg.knownAddresses = ui.Addresses
}

func (gg *RedeemCheckGenerator) Generate() Action {
	return &RedeemCheckAction{
		coin: sdk.NewCoin(
			RandomChoice(gg.rnd, gg.knownCoins),
			helpers.FinneyToWei(sdkmath.NewInt(RandomRange(gg.rnd, gg.bottomRange, gg.upperRange))),
		),
		issuer:   RandomChoice(gg.rnd, gg.knownAddresses),
		receiver: RandomChoice(gg.rnd, gg.knownAddresses),
	}
}

func (ac *RedeemCheckAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].Address() == ac.issuer {
			ac.issuerAcc = saList[i]
		}
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].Address() != ac.receiver {
			continue
		}
		if saList[i].BalanceForCoin(ac.coin.Denom).LT(ac.coin.Amount) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (ac *RedeemCheckAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	const password = "der_password"

	if ac.issuerAcc == nil {
		return nil, fmt.Errorf("empty issuer account")
	}
	nonce := sdkmath.NewInt(rand.Int63())
	dueBlock := sa.LastHeight()
	if dueBlock == 0 {
		return nil, fmt.Errorf("block height is 0")
	}

	dueBlock += rand.Int63n(100) + 100

	checkBase58, err := dscTx.IssueCheck(ac.issuerAcc.Account(), ac.coin.Denom, ac.coin.Amount, nonce, uint64(dueBlock), password)
	if err != nil {
		return nil, err
	}

	msg, err := dscTx.CreateRedeemCheck(sa.Account(), checkBase58, password)
	if err != nil {
		return nil, err
	}

	// Redeem check has fixed fee, zero at fee decorator
	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *RedeemCheckAction) String() string {
	return fmt.Sprintf("RedeemCheck{issuer: %s, receiver: %s, coin: %s}", ac.issuer, ac.receiver, ac.coin.String())
}
