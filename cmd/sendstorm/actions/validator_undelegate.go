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

type UndelegateGenerator struct {
	knownDelegations []dscApi.Delegation
	rnd              *rand.Rand
}

type UndelegateAction struct {
	coin             sdk.Coin
	delegatorAddress string
	validatorAddress string
}

func NewUndelegateGenerator() *UndelegateGenerator {
	return &UndelegateGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *UndelegateGenerator) Update(ui UpdateInfo) {
	gg.knownDelegations = ui.Delegations
}

func (gg *UndelegateGenerator) Generate() Action {
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

	return &UndelegateAction{
		delegatorAddress: stake.Delegator,
		validatorAddress: stake.Validator,
		coin:             coin,
	}
}

func (ac *UndelegateAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (ac *UndelegateAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	val, err := sdk.ValAddressFromBech32(ac.validatorAddress)
	if err != nil {
		return nil, err
	}
	msg := dscTx.NewMsgUndelegate(sa.Account().SdkAddress(), val, ac.coin)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *UndelegateAction) String() string {
	return fmt.Sprintf("UndelegateAction(val:%s,del:%s)", ac.validatorAddress, ac.delegatorAddress)
}
