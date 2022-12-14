package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DelegateNFTGenerator struct {
	knownNFT        []*dscApi.NFTToken
	knownValidators []dscApi.Validator
	rnd             *rand.Rand
}

type DelegateNFTAction struct {
	token            *dscApi.NFTToken
	subToken         *dscApi.SubToken
	validatorAddress string
}

func NewDelegateNFTGenerator() *DelegateNFTGenerator {
	return &DelegateNFTGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *DelegateNFTGenerator) Update(ui UpdateInfo) {
	gg.knownNFT = ui.NFTs
	gg.knownValidators = ui.Validators
}

func (gg *DelegateNFTGenerator) Generate() Action {
	if len(gg.knownNFT) == 0 {
		return &EmptyAction{}
	}
	if len(gg.knownValidators) == 0 {
		return &EmptyAction{}
	}
	token := RandomChoice(gg.rnd, gg.knownNFT)
	if len(token.SubTokens) == 0 {
		return &EmptyAction{}
	}
	sub := RandomChoice(gg.rnd, token.SubTokens)
	return &DelegateNFTAction{
		token:            token,
		subToken:         sub,
		validatorAddress: RandomChoice(gg.rnd, gg.knownValidators).OperatorAddress,
	}
}

func (ac *DelegateNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if ac.subToken.Owner != saList[i].Address() {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (ac *DelegateNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	valAdr, err := sdk.ValAddressFromBech32(ac.validatorAddress)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgDelegateNFT(sa.Account().SdkAddress(), valAdr, ac.token.ID, []uint32{ac.subToken.ID})

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *DelegateNFTAction) String() string {
	return "DelegateNFTAction"
}
