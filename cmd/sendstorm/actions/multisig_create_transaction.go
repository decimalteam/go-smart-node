package actions

import (
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/strings"
)

type CreateMultisigTransactionGenerator struct {
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	knownAddresses          []string
	knownWallets            []dscApi.MultisigWallet
	knownMultisigBalances   map[string]sdk.Coins
	rnd                     *rand.Rand
}

type CreateMultisigTransactionAction struct {
	coins           sdk.Coins
	receiver        string
	wallet          string
	possibleSenders []string // wallet owners
}

func NewCreateMultisigTransactionGenerator(bottomRange, upperRange int64) *CreateMultisigTransactionGenerator {
	return &CreateMultisigTransactionGenerator{
		bottomRange: bottomRange,
		upperRange:  upperRange,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *CreateMultisigTransactionGenerator) Update(ui UpdateInfo) {
	gg.knownAddresses = ui.Addresses
	gg.knownWallets = ui.MultisigWallets
	gg.knownMultisigBalances = ui.MultisigBalances
}

func (gg *CreateMultisigTransactionGenerator) Generate() Action {
	if len(gg.knownWallets) == 0 {
		return &EmptyAction{}
	}
	i := int(RandomRange(gg.rnd, 0, int64(len(gg.knownWallets))))
	wallet := gg.knownWallets[i]
	balance, ok := gg.knownMultisigBalances[wallet.Address]
	if !ok || balance.IsZero() {
		return &EmptyAction{}
	}
	j := int(RandomRange(gg.rnd, 0, int64(balance.Len())))
	upperLimit := helpers.WeiToFinney(balance[j].Amount).Int64()
	if upperLimit > gg.upperRange {
		upperLimit = gg.upperRange
	}
	if upperLimit < gg.bottomRange {
		return &EmptyAction{}
	}
	amount := helpers.FinneyToWei(sdk.NewInt(RandomRange(gg.rnd, gg.bottomRange, upperLimit)))
	coinToSend := sdk.NewCoin(balance[j].Denom, amount)
	return &CreateMultisigTransactionAction{
		coins:           sdk.NewCoins(coinToSend),
		receiver:        RandomChoice(gg.rnd, gg.knownAddresses),
		wallet:          wallet.Address,
		possibleSenders: wallet.Owners,
	}
}

func (aa *CreateMultisigTransactionAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if !strings.StringInSlice(saList[i].Address(), aa.possibleSenders) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *CreateMultisigTransactionAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgCreateTransaction(sender, aa.wallet, aa.receiver, aa.coins)
	tx, err := dscTx.BuildTransaction(sa.Account(), []sdk.Msg{msg}, "", sa.FeeDenom(), feeConfig.DelPrice, feeConfig.Params)
	if err != nil {
		return nil, err
	}
	err = tx.SignTransaction(sa.Account())
	if err != nil {
		return nil, err
	}
	return tx.BytesToSend()
}

func (aa *CreateMultisigTransactionAction) String() string {
	return fmt.Sprintf("CreateMultisigTransaction{wallet: %s, receiver: %s, coin: %s}",
		aa.wallet, aa.receiver, aa.coins.String())
}
