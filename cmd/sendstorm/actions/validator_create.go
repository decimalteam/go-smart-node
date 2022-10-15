package actions

import (
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CreateValidatorGenerator struct {
	initialStackBottom int64 // in 10^18
	initialStackUp     int64 // in 10^18
	rnd                *rand.Rand
}

type CreateValidatorAction struct {
	pubKey []byte
	rate   int64
	//
	moniker         string
	details         string
	identity        string
	securityContact string
	website         string
}

func NewCreateValidatorGenerator(
	initialStackBottom, initialStackUp int64) *CreateValidatorGenerator {
	return &CreateValidatorGenerator{
		initialStackBottom: initialStackBottom,
		initialStackUp:     initialStackUp,
		rnd:                rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *CreateValidatorGenerator) Update(ui UpdateInfo) {
}

func (gg *CreateValidatorGenerator) Generate() Action {
	return &CreateValidatorAction{
		// TODO: save or generate from mnemonic private key
		pubKey: []byte(RandomString(gg.rnd, 10, charsAll)),
		rate:   RandomRange(gg.rnd, 10, 100+1),
		//
		moniker:         RandomString(gg.rnd, 10, charsAll),
		details:         RandomString(gg.rnd, 10, charsAll),
		identity:        RandomString(gg.rnd, 10, charsAll),
		securityContact: RandomString(gg.rnd, 10, charsAll),
		website:         RandomString(gg.rnd, 10, charsAll),
	}
}

func (ac *CreateValidatorAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		// TODO: future checks
		res = append(res, saList[i])
	}
	return res
}

func (ac *CreateValidatorAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	_, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}
	// TODO
	/*
		msg := dscTx.NewMsgCreateCoin(
			sender,
			ac.symbol,
			ac.title,
			ac.crr,
			ac.initVolume,
			ac.initReserve,
			ac.limitVolume,
			ac.identity,
		)
		tx, err := dscTx.BuildTransaction(sa.Account(), []sdk.Msg{msg}, "", sa.FeeDenom(), feeConfig.DelPrice, feeConfig.Params)
	*/
	return feeConfig.MakeTransaction(sa, nil)
}

func (ac *CreateValidatorAction) String() string {
	return "CreateValidator"
}
