package actions

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"math/rand"
	"time"

	//"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
)

type CreateValidatorGenerator struct {
	initialStackBottom int64 // in 10^18
	initialStackUp     int64 // in 10^18
	knownCoins         []string
	pks                map[string][]byte
	rnd                *rand.Rand
}

type CreateValidatorAction struct {
	pubKey cryptotypes.PubKey
	rate   int64
	stake  sdk.Coin
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
	gg.knownCoins = ui.Coins
}

func (gg *CreateValidatorGenerator) Generate() Action {
	if len(gg.knownCoins) == 0 {
		return &EmptyAction{}
	}

	denom := RandomChoice(gg.rnd, gg.knownCoins)
	amount := helpers.EtherToWei(sdk.NewInt(RandomRange(gg.rnd, gg.initialStackBottom, gg.initialStackUp)))

	return &CreateValidatorAction{
		rate:  RandomRange(gg.rnd, 10, 100+1),
		stake: sdk.NewCoin(denom, amount),
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
		if saList[i].BalanceForCoin(ac.stake.Denom).LT(ac.stake.Amount) {
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
	msg, err := dscTx.NewMsgCreateValidator(
		sdk.ValAddress(sa.Account().SdkAddress()),
		sa.Account().SdkAddress(),
		sa.Account().ValPubKey(),
		dscTx.Description{
			Moniker:         ac.moniker,
			Identity:        ac.identity,
			Website:         ac.website,
			SecurityContact: ac.securityContact,
			Details:         ac.details,
		},
		sdk.MustNewDecFromStr("0.1"),
		ac.stake,
	)
	if err != nil {
		return nil, err
	}

	return feeConfig.MakeTransaction(sa, msg)
}

func (ac *CreateValidatorAction) String() string {
	return "CreateValidator"
}
