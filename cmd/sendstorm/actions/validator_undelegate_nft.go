package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type UndelegateNFTGenerator struct {
	knownNFTStakes []NFTStake
	rnd            *rand.Rand
}

type UndelegateNFTAction struct {
	token            dscApi.NFTToken
	delegatorAddress string
	validatorAddress string
}

func NewUndelegateNFTGenerator() *UndelegateNFTGenerator {
	return &UndelegateNFTGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *UndelegateNFTGenerator) Update(ui UpdateInfo) {
	gg.knownNFTStakes = ui.NFTStakes
}

func (gg *UndelegateNFTGenerator) Generate() Action {
	if len(gg.knownNFTStakes) == 0 {
		return &EmptyAction{}
	}
	stake := RandomChoice(gg.rnd, gg.knownNFTStakes)
	return &UndelegateNFTAction{
		delegatorAddress: stake.Delegator,
		validatorAddress: stake.Validator,
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

	// TODO

	return feeConfig.MakeTransaction(sa, nil)
}

func (ac *UndelegateNFTAction) String() string {
	return "UndelegateNFTAction"
}
