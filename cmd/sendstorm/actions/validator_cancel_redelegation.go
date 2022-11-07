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

type CancelRedelegationGenerator struct {
	knownRedelegations []dscApi.Redelegation
	rnd                *rand.Rand
}

type CancelRedelegationAction struct {
	coin                 sdk.Coin
	creationHeight       int64
	delegatorAddress     string
	fromValidatorAddress string
	toValidatorAddress   string
}

func NewCancelRedelegationGenerator() *CancelRedelegationGenerator {
	return &CancelRedelegationGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *CancelRedelegationGenerator) Update(ui UpdateInfo) {
	gg.knownRedelegations = ui.Redelegations
}

func (gg *CancelRedelegationGenerator) Generate() Action {
	if len(gg.knownRedelegations) == 0 {
		return &EmptyAction{}
	}
	type redelegationRecord struct {
		coin                 sdk.Coin
		creationHeight       int64
		delegatorAddress     string
		fromValidatorAddress string
		toValidatorAddress   string
	}
	var coinUndelegations = make([]redelegationRecord, 0)
	for _, redel := range gg.knownRedelegations {
		for _, entry := range redel.Entries {
			if entry.Stake.Type == dscApi.StakeType_Coin {
				coinUndelegations = append(coinUndelegations, redelegationRecord{
					coin:                 entry.Stake.Stake,
					creationHeight:       entry.CreationHeight,
					delegatorAddress:     redel.Delegator,
					fromValidatorAddress: redel.ValidatorSrc,
					toValidatorAddress:   redel.ValidatorDst,
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

	return &CancelRedelegationAction{
		coin:                 coin,
		creationHeight:       entry.creationHeight,
		delegatorAddress:     entry.delegatorAddress,
		fromValidatorAddress: entry.fromValidatorAddress,
		toValidatorAddress:   entry.toValidatorAddress,
	}
}

func (ac *CancelRedelegationAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (ac *CancelRedelegationAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
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
	msg := dscTx.NewMsgCancelRedelegation(sa.Account().SdkAddress(), valSrc, valDst, ac.creationHeight, ac.coin)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *CancelRedelegationAction) String() string {
	return fmt.Sprintf("CancelRedelegationAction(%s->%s, h=%d, id=%s)", ac.fromValidatorAddress, ac.toValidatorAddress, ac.creationHeight, ac.coin.Denom)
}
