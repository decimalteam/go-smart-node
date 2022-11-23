package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CancelRedelegationNFTGenerator struct {
	knownRedelegations []dscApi.Redelegation
	rnd                *rand.Rand
}

type CancelRedelegationNFTAction struct {
	tokenID              string
	subTokenIDs          []uint32
	creationHeight       int64
	delegatorAddress     string
	fromValidatorAddress string
	toValidatorAddress   string
}

func NewCancelRedelegationNFTGenerator() *CancelRedelegationNFTGenerator {
	return &CancelRedelegationNFTGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *CancelRedelegationNFTGenerator) Update(ui UpdateInfo) {
	gg.knownRedelegations = ui.Redelegations
}

func (gg *CancelRedelegationNFTGenerator) Generate() Action {
	if len(gg.knownRedelegations) == 0 {
		return &EmptyAction{}
	}
	type redelegationRecord struct {
		tokenID              string
		subTokenIDs          []uint32
		creationHeight       int64
		delegatorAddress     string
		fromValidatorAddress string
		toValidatorAddress   string
	}
	var nftUndelegations = make([]redelegationRecord, 0)
	for _, redel := range gg.knownRedelegations {
		for _, entry := range redel.Entries {
			if entry.Stake.Type == dscApi.StakeType_NFT {
				nftUndelegations = append(nftUndelegations, redelegationRecord{
					tokenID:              entry.Stake.ID,
					subTokenIDs:          entry.Stake.SubTokenIDs,
					creationHeight:       entry.CreationHeight,
					delegatorAddress:     redel.Delegator,
					fromValidatorAddress: redel.ValidatorSrc,
					toValidatorAddress:   redel.ValidatorDst,
				})
			}
		}
	}
	if len(nftUndelegations) == 0 {
		return &EmptyAction{}
	}

	entry := RandomChoice(gg.rnd, nftUndelegations)
	subs := RandomSublist(gg.rnd, entry.subTokenIDs)
	if len(subs) == 0 {
		return &EmptyAction{}
	}

	return &CancelRedelegationNFTAction{
		tokenID:              entry.tokenID,
		subTokenIDs:          subs,
		creationHeight:       entry.creationHeight,
		delegatorAddress:     entry.delegatorAddress,
		fromValidatorAddress: entry.fromValidatorAddress,
		toValidatorAddress:   entry.toValidatorAddress,
	}
}

func (ac *CancelRedelegationNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (ac *CancelRedelegationNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
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
	msg := dscTx.NewMsgCancelRedelegationNFT(sa.Account().SdkAddress(), valSrc, valDst, ac.creationHeight, ac.tokenID, ac.subTokenIDs)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *CancelRedelegationNFTAction) String() string {
	return "CancelRedelegationNFTAction"
}
