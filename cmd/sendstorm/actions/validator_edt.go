package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EditValidatorGenerator struct {
	knownValidators []string
	rnd             *rand.Rand
}

type EditValidatorAction struct {
	creatorAddress   string
	validatorAddress string
}

func NewEditValidatorGenerator() *EditValidatorGenerator {
	return &EditValidatorGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *EditValidatorGenerator) Update(ui UpdateInfo) {
	gg.knownValidators = ui.Validators
}

func (gg *EditValidatorGenerator) Generate() Action {
	if len(gg.knownValidators) == 0 {
		return &EmptyAction{}
	}
	val := RandomChoice(gg.rnd, gg.knownValidators)
	return &EditValidatorAction{
		validatorAddress: val,
	}
}

func (ac *EditValidatorAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].Address() != ac.creatorAddress {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (ac *EditValidatorAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	// TODO

	return feeConfig.MakeTransaction(sa, nil)
}

func (ac *EditValidatorAction) String() string {
	return "RedelegateNFTAction"
}
