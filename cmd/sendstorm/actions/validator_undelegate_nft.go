package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type UndelegateNFTGenerator struct {
	knownDelegations []dscApi.Delegation
	rnd              *rand.Rand
}

type UndelegateNFTAction struct {
	delegatorAddress string
	validatorAddress string
	tokenID          string
	subTokenIDs      []uint32
}

func NewUndelegateNFTGenerator() *UndelegateNFTGenerator {
	return &UndelegateNFTGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *UndelegateNFTGenerator) Update(ui UpdateInfo) {
	gg.knownDelegations = ui.Delegations
}

func (gg *UndelegateNFTGenerator) Generate() Action {
	if len(gg.knownDelegations) == 0 {
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

	return &UndelegateNFTAction{
		delegatorAddress: stake.Delegator,
		validatorAddress: stake.Validator,
		tokenID:          stake.Stake.ID,
		subTokenIDs:      subs,
	}
}

func (ac *UndelegateNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (ac *UndelegateNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	val, err := sdk.ValAddressFromBech32(ac.validatorAddress)
	if err != nil {
		return nil, err
	}
	msg := dscTx.NewMsgUndelegateNFT(sa.Account().SdkAddress(), val, ac.tokenID, ac.subTokenIDs)

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *UndelegateNFTAction) String() string {
	return "UndelegateNFTAction"
}
