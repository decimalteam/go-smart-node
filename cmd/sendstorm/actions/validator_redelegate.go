package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RedelegateGenerator struct {
	knownStakes     []GenericStake
	knownValidators []string
	rnd             *rand.Rand
}

type RedelegateAction struct {
	delegatorAddress     string
	fromValidatorAddress string
	toValidatorAddress   string
	coin                 sdk.Coin
}

func NewRedelegateGenerator() *RedelegateGenerator {
	return &RedelegateGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *RedelegateGenerator) Update(ui UpdateInfo) {
	gg.knownStakes = ui.Stakes
	gg.knownValidators = ui.Validators
}

func (gg *RedelegateGenerator) Generate() Action {
	if len(gg.knownValidators) < 2 {
		return &EmptyAction{}
	}
	if len(gg.knownStakes) == 0 {
		return &EmptyAction{}
	}
	stake := RandomChoice(gg.rnd, gg.knownStakes)
	toValidator := ""
	for i := 0; i < 10; i++ {
		toValidator = RandomChoice(gg.rnd, gg.knownValidators)
		if toValidator != stake.Validator {
			break
		}
	}
	if toValidator == stake.Validator {
		return &EmptyAction{}
	}

	return &RedelegateAction{
		delegatorAddress:     stake.Delegator,
		fromValidatorAddress: stake.Validator,
		toValidatorAddress:   toValidator,
		coin:                 stake.Coin,
	}
}

func (ac *RedelegateAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].Address() != ac.delegatorAddress {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (ac *RedelegateAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	// TODO

	return feeConfig.MakeTransaction(sa, nil)
}

func (ac *RedelegateAction) String() string {
	return "RedelegateAction"
}
