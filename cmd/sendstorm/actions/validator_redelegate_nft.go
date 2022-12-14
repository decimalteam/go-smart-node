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
	knownDelegations []dscApi.Delegation
	knownValidators  []dscApi.Validator
	rnd              *rand.Rand
}

type RedelegateNFTAction struct {
	delegatorAddress     string
	fromValidatorAddress string
	toValidatorAddress   string
	tokenID              string
	subTokenIDs          []uint32
}

func NewRedelegateNFTGenerator() *RedelegateNFTGenerator {
	return &RedelegateNFTGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *RedelegateNFTGenerator) Update(ui UpdateInfo) {
	gg.knownDelegations = ui.Delegations
	gg.knownValidators = ui.Validators
}

func (gg *RedelegateNFTGenerator) Generate() Action {
	if len(gg.knownDelegations) == 0 {
		return &EmptyAction{}
	}
	if len(gg.knownValidators) < 2 {
		return &EmptyAction{}
	}
	var nftDelegations = make([]dscApi.Delegation, 0)
	for _, del := range gg.knownDelegations {
		if del.Stake.Type == dscApi.StakeType_NFT {
			nftDelegations = append(nftDelegations, del)
		}
	}
	if len(nftDelegations) == 0 {
		return &EmptyAction{}
	}
	stake := RandomChoice(gg.rnd, nftDelegations)
	subs := RandomSublist(gg.rnd, stake.Stake.SubTokenIDs)
	if len(subs) == 0 {
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
	return &RedelegateNFTAction{
		delegatorAddress:     stake.Delegator,
		fromValidatorAddress: stake.Validator,
		toValidatorAddress:   toValidator,
		tokenID:              stake.Stake.ID,
		subTokenIDs:          subs,
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

	valSrc, err := sdk.ValAddressFromBech32(ac.fromValidatorAddress)
	if err != nil {
		return nil, err
	}
	valDst, err := sdk.ValAddressFromBech32(ac.toValidatorAddress)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgRedelegateNFT(sa.Account().SdkAddress(), valSrc, valDst, ac.tokenID, ac.subTokenIDs)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *RedelegateNFTAction) String() string {
	return "RedelegateNFTAction"
}
