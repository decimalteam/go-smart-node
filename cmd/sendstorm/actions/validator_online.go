package actions

import (
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SetOnlineValidatorGenerator struct {
	knownValidators []dscApi.Validator
	rnd             *rand.Rand
}

type SetOnlineValidatorAction struct {
	creatorAddress   string
	validatorAddress string
}

func NewSetOnlineValidatorGenerator() *SetOnlineValidatorGenerator {
	return &SetOnlineValidatorGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *SetOnlineValidatorGenerator) Update(ui UpdateInfo) {
	gg.knownValidators = ui.Validators
}

func (gg *SetOnlineValidatorGenerator) Generate() Action {
	if len(gg.knownValidators) == 0 {
		return &EmptyAction{}
	}
	var val dscApi.Validator
	// 10 attempts to find offline validator
	for i := 0; i < 10; i++ {
		val = RandomChoice(gg.rnd, gg.knownValidators)
		if !val.Online {
			break
		}
	}
	if val.Online {
		return &EmptyAction{}
	}
	adr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	if err != nil {
		return &EmptyAction{}
	}
	return &SetOnlineValidatorAction{
		creatorAddress:   sdk.AccAddress(adr).String(),
		validatorAddress: val.OperatorAddress,
	}
}

func (ac *SetOnlineValidatorAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (ac *SetOnlineValidatorAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	valAdr, err := sdk.ValAddressFromBech32(ac.validatorAddress)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgSetOnline(valAdr)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *SetOnlineValidatorAction) String() string {
	return fmt.Sprintf("SetOnlineValidatorAction(%s)", ac.validatorAddress)
}
