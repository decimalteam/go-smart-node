package actions

import (
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/strings"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

type CreateMultisigTransactionGenerator struct {
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	knownAddresses          []string
	knownWallets            []dscApi.MultisigWallet
	knownMultisigBalances   map[string]sdk.Coins
	rnd                     *rand.Rand
}

type CreateMultisigTransactionAction struct {
	coin            sdk.Coin
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
	gg.knownAddresses = make([]string, 0, len(ui.Addresses)+len(ui.MultisigWallets))
	gg.knownAddresses = append(gg.knownAddresses, ui.Addresses...)
	for _, w := range ui.MultisigWallets {
		gg.knownAddresses = append(gg.knownAddresses, w.Address)
	}
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
	if upperLimit <= gg.bottomRange {
		return &EmptyAction{}
	}
	amount := helpers.FinneyToWei(sdkmath.NewInt(RandomRange(gg.rnd, gg.bottomRange, upperLimit)))
	coinToSend := sdk.NewCoin(balance[j].Denom, amount)
	return &CreateMultisigTransactionAction{
		coin:            coinToSend,
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

	wAdr, err := sdk.AccAddressFromBech32(aa.wallet)
	if err != nil {
		return nil, err
	}
	rAdr, err := sdk.AccAddressFromBech32(aa.receiver)
	if err != nil {
		return nil, err
	}

	msg, err := dscTx.NewMsgCreateTransaction(sender, aa.wallet, dscTx.NewMsgSendCoin(wAdr, rAdr, aa.coin))
	if err != nil {
		return nil, err
	}

	return feeConfig.MakeTransaction(sa, msg)
}

func (aa *CreateMultisigTransactionAction) String() string {
	return fmt.Sprintf("CreateMultisigTransaction{wallet: %s, receiver: %s, coin: %s}",
		aa.wallet, aa.receiver, aa.coin.String())
}
