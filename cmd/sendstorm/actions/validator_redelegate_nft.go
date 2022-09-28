package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RedelegateNFTGenerator struct {
	knownNFTStakes  []NFTStake
	knownValidators []string
	rnd             *rand.Rand
}

type RedelegateNFTAction struct {
	token                dscApi.NFTToken
	delegatorAddress     string
	fromValidatorAddress string
	toValidatorAddress   string
}

func NewRedelegateNFTGenerator() *RedelegateNFTGenerator {
	return &RedelegateNFTGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *RedelegateNFTGenerator) Update(ui UpdateInfo) {
	gg.knownNFTStakes = ui.NFTStakes
	gg.knownValidators = ui.Validators
}

func (gg *RedelegateNFTGenerator) Generate() Action {
	if len(gg.knownNFTStakes) == 0 {
		return &EmptyAction{}
	}
	if len(gg.knownValidators) < 2 {
		return &EmptyAction{}
	}
	stake := RandomChoice(gg.rnd, gg.knownNFTStakes)
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
	return &RedelegateNFTAction{
		delegatorAddress:     stake.Delegator,
		fromValidatorAddress: stake.Validator,
		toValidatorAddress:   toValidator,
	}
}

func (ac *RedelegateNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (ac *RedelegateNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	// TODO

	tx, err := dscTx.BuildTransaction(sa.Account(), []sdk.Msg{}, "", sa.FeeDenom(), feeConfig.DelPrice, feeConfig.Params)
	if err != nil {
		return nil, err
	}
	err = tx.SignTransaction(sa.Account())
	if err != nil {
		return nil, err
	}
	return tx.BytesToSend()
}

func (ac *RedelegateNFTAction) String() string {
	return "RedelegateNFTAction"
}
