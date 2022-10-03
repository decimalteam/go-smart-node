package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DelegateGenerator struct {
	stackBottom     int64 // in 10^15
	stackUp         int64 // in 10^15
	knownCoins      []string
	knownValidators []string
	rnd             *rand.Rand
}

type DelegateAction struct {
	coin             sdk.Coin
	validatorAddress string
}

func NewDelegateGenerator(
	stackBottom, stackUp int64) *DelegateGenerator {
	return &DelegateGenerator{
		stackBottom: stackBottom,
		stackUp:     stackUp,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *DelegateGenerator) Update(ui UpdateInfo) {
	gg.knownCoins = ui.Coins
	gg.knownValidators = ui.Validators
}

func (gg *DelegateGenerator) Generate() Action {
	if len(gg.knownCoins) == 0 {
		return &EmptyAction{}
	}
	if len(gg.knownValidators) == 0 {
		return &EmptyAction{}
	}
	denom := RandomChoice(gg.rnd, gg.knownCoins)
	validator := RandomChoice(gg.rnd, gg.knownValidators)
	amount := RandomRange(gg.rnd, gg.stackBottom, gg.stackUp)
	return &DelegateAction{
		coin:             sdk.NewCoin(denom, helpers.FinneyToWei(sdk.NewInt(amount))),
		validatorAddress: validator,
	}
}

func (ac *DelegateAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].BalanceForCoin(ac.coin.Denom).LT(ac.coin.Amount) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (ac *DelegateAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	// TODO

	tx, err := dscTx.BuildTransaction(sa.Account(), []sdk.Msg{}, "", sa.FeeDenom(), feeConfig)
	if err != nil {
		return nil, err
	}
	err = tx.SignTransaction(sa.Account())
	if err != nil {
		return nil, err
	}
	return tx.BytesToSend()
}

func (ac *DelegateAction) String() string {
	return "DelegateAction"
}
