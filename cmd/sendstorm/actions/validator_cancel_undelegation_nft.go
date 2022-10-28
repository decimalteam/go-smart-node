package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CancelUndelegationNFTGenerator struct {
	knownUndelegations []dscApi.Undelegation
	rnd                *rand.Rand
}

type CancelUndelegationNFTAction struct {
	tokenID          string
	subTokensIDs     []uint32
	creationHeight   int64
	delegatorAddress string
	validatorAddress string
}

func NewCancelUndelegationNFTGenerator() *CancelUndelegationNFTGenerator {
	return &CancelUndelegationNFTGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *CancelUndelegationNFTGenerator) Update(ui UpdateInfo) {
	gg.knownUndelegations = ui.Undelegations
}

func (gg *CancelUndelegationNFTGenerator) Generate() Action {
	if len(gg.knownUndelegations) == 0 {
		return &EmptyAction{}
	}
	type undelegationRecord struct {
		tokenID          string
		subTokensIDs     []uint32
		creationHeight   int64
		delegatorAddress string
		validatorAddress string
	}
	var nftUndelegations = make([]undelegationRecord, 0)
	for _, undel := range gg.knownUndelegations {
		for _, entry := range undel.Entries {
			if entry.Stake.Type == dscApi.StakeType_NFT {
				nftUndelegations = append(nftUndelegations, undelegationRecord{
					tokenID:          entry.Stake.ID,
					subTokensIDs:     entry.Stake.SubTokenIDs,
					creationHeight:   entry.CreationHeight,
					delegatorAddress: undel.Delegator,
					validatorAddress: undel.Validator,
				})
			}
		}
	}
	if len(nftUndelegations) == 0 {
		return &EmptyAction{}
	}

	entry := RandomChoice(gg.rnd, nftUndelegations)
	subs := RandomSublist(gg.rnd, entry.subTokensIDs)
	if len(subs) == 0 {
		return &EmptyAction{}
	}

	return &CancelUndelegationNFTAction{
		tokenID:          entry.tokenID,
		subTokensIDs:     subs,
		creationHeight:   entry.creationHeight,
		delegatorAddress: entry.delegatorAddress,
		validatorAddress: entry.validatorAddress,
	}
}

func (ac *CancelUndelegationNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (ac *CancelUndelegationNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	val, err := sdk.ValAddressFromBech32(ac.validatorAddress)
	if err != nil {
		return nil, err
	}
	msg := dscTx.NewMsgCancelUndelegationNFT(sa.Account().SdkAddress(), val, ac.creationHeight, ac.tokenID, ac.subTokensIDs)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *CancelUndelegationNFTAction) String() string {
	return "CancelUndelegationNFTAction"
}
