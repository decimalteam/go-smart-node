package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CancelUndelegationGenerator struct {
	knownUndelegations []dscApi.Undelegation
	rnd                *rand.Rand
}

type CancelUndelegationAction struct {
	coin             sdk.Coin
	creationHeight   int64
	delegatorAddress string
	validatorAddress string
}

func NewCancelUndelegationGenerator() *CancelUndelegationGenerator {
	return &CancelUndelegationGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *CancelUndelegationGenerator) Update(ui UpdateInfo) {
	gg.knownUndelegations = ui.Undelegations
}

func (gg *CancelUndelegationGenerator) Generate() Action {
	if len(gg.knownUndelegations) == 0 {
		return &EmptyAction{}
	}
	type undelegationRecord struct {
		coin             sdk.Coin
		creationHeight   int64
		delegatorAddress string
		validatorAddress string
	}
	var coinUndelegations = make([]undelegationRecord, 0)
	for _, undel := range gg.knownUndelegations {
		for _, entry := range undel.Entries {
			if entry.Stake.Type == dscApi.StakeType_Coin {
				coinUndelegations = append(coinUndelegations, undelegationRecord{
					coin:             entry.Stake.Stake,
					creationHeight:   entry.CreationHeight,
					delegatorAddress: undel.Delegator,
					validatorAddress: undel.Validator,
				})
			}
		}
	}
	if len(coinUndelegations) == 0 {
		return &EmptyAction{}
	}

	entry := RandomChoice(gg.rnd, coinUndelegations)
	coin := ExtractPartCoin(gg.rnd, entry.coin)
	if coin.IsZero() {
		return &EmptyAction{}
	}

	return &CancelUndelegationAction{
		coin:             coin,
		creationHeight:   entry.creationHeight,
		delegatorAddress: entry.delegatorAddress,
		validatorAddress: entry.validatorAddress,
	}
}

func (ac *CancelUndelegationAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (ac *CancelUndelegationAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	val, err := sdk.ValAddressFromBech32(ac.validatorAddress)
	if err != nil {
		return nil, err
	}
	msg := dscTx.NewMsgCancelUndelegation(sa.Account().SdkAddress(), val, ac.creationHeight, ac.coin)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *CancelUndelegationAction) String() string {
	return "CancelUndelegationAction"
}
