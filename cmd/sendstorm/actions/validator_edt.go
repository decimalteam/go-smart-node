package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EditValidatorGenerator struct {
	knownValidators []dscApi.Validator
	knownAddresses  []string
	rnd             *rand.Rand
}

type EditValidatorAction struct {
	creatorAddress   string
	validatorAddress string
	newRewardAddress string
	description      dscTx.Description
}

func NewEditValidatorGenerator() *EditValidatorGenerator {
	return &EditValidatorGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *EditValidatorGenerator) Update(ui UpdateInfo) {
	gg.knownValidators = ui.Validators
	gg.knownAddresses = ui.Addresses
}

func (gg *EditValidatorGenerator) Generate() Action {
	if len(gg.knownValidators) == 0 {
		return &EmptyAction{}
	}
	if len(gg.knownAddresses) == 0 {
		return &EmptyAction{}
	}
	val := RandomChoice(gg.rnd, gg.knownValidators)
	adr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	if err != nil {
		return &EmptyAction{}
	}
	return &EditValidatorAction{
		creatorAddress:   sdk.AccAddress(adr).String(),
		validatorAddress: val.OperatorAddress,
		newRewardAddress: RandomChoice(gg.rnd, gg.knownAddresses),
		description: dscTx.Description{
			Moniker:         RandomString(gg.rnd, 10, charsAll),
			Identity:        RandomString(gg.rnd, 10, charsAll),
			Website:         RandomString(gg.rnd, 10, charsAll),
			SecurityContact: RandomString(gg.rnd, 10, charsAll),
			Details:         RandomString(gg.rnd, 10, charsAll),
		},
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

	valAdr, err := sdk.ValAddressFromBech32(ac.validatorAddress)
	if err != nil {
		return nil, err
	}
	rewardAdr, err := sdk.AccAddressFromBech32(ac.newRewardAddress)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgEditValidator(valAdr, rewardAdr, ac.description)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *EditValidatorAction) String() string {
	return "EditValidatorAction"
}
