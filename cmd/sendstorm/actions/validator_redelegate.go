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

type RedelegateGenerator struct {
	knownDelegations []dscApi.Delegation
	knownValidators  []dscApi.Validator
	rnd              *rand.Rand
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
	gg.knownDelegations = ui.Delegations
	gg.knownValidators = ui.Validators
}

func (gg *RedelegateGenerator) Generate() Action {
	if len(gg.knownValidators) < 2 {
		return &EmptyAction{}
	}
	if len(gg.knownDelegations) == 0 {
		return &EmptyAction{}
	}
	var coinDelegations = make([]dscApi.Delegation, 0)
	for _, del := range gg.knownDelegations {
		if del.Stake.Type == dscApi.StakeType_Coin {
			coinDelegations = append(coinDelegations, del)
		}
	}
	if len(coinDelegations) == 0 {
		return &EmptyAction{}
	}

	stake := RandomChoice(gg.rnd, coinDelegations)
	coin := ExtractPartCoin(gg.rnd, stake.Stake.Stake)
	if coin.IsZero() {
		return &EmptyAction{}
	}
	toValidator := ""
	for i := 0; i < 10; i++ {
		toValidator = RandomChoice(gg.rnd, gg.knownValidators).OperatorAddress
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
		coin:                 coin,
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

	valSrc, err := sdk.ValAddressFromBech32(ac.fromValidatorAddress)
	if err != nil {
		return nil, err
	}
	valDst, err := sdk.ValAddressFromBech32(ac.toValidatorAddress)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgRedelegate(sa.Account().SdkAddress(), valSrc, valDst, ac.coin)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *RedelegateAction) String() string {
	return fmt.Sprintf("RedelegateAction(%s->%s)", ac.fromValidatorAddress, ac.toValidatorAddress)
}
