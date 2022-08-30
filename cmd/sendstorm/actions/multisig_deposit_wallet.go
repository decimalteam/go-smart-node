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
)

// like SendCoin, but for multisig wallets as receivers
type DepositMultisigWalletGenerator struct {
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	knownAddresses          []string
	knownCoins              []string
	knownWallets            []dscApi.MultisigWallet
	rnd                     *rand.Rand
}

type DepositMultisigWalletAction struct {
	coin     sdk.Coin
	receiver string
}

func NewDepositMultisigWalletGenerator(bottomRange, upperRange int64) *DepositMultisigWalletGenerator {
	return &DepositMultisigWalletGenerator{
		bottomRange: bottomRange,
		upperRange:  upperRange,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *DepositMultisigWalletGenerator) Update(ui UpdateInfo) {
	gg.knownCoins = ui.Coins
	gg.knownAddresses = ui.Addresses
	gg.knownWallets = ui.MultisigWallets
}

func (gg *DepositMultisigWalletGenerator) Generate() Action {
	if len(gg.knownWallets) == 0 {
		return &EmptyAction{}
	}
	i := int(RandomRange(gg.rnd, 0, int64(len(gg.knownWallets))))
	receiver := gg.knownWallets[i].Address
	return &DepositMultisigWalletAction{
		coin: sdk.NewCoin(
			RandomChoice(gg.rnd, gg.knownCoins),
			helpers.FinneyToWei(sdk.NewInt(RandomRange(gg.rnd, gg.bottomRange, gg.upperRange))),
		),
		receiver: receiver,
	}
}

func (aa *DepositMultisigWalletAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].BalanceForCoin(aa.coin.Denom).LT(aa.coin.Amount) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *DepositMultisigWalletAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}
	receiver, err := sdk.AccAddressFromBech32(aa.receiver)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgSendCoin(sender, aa.coin, receiver)
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

func (aa *DepositMultisigWalletAction) String() string {
	return fmt.Sprintf("DepositMultisigWallet{receiver: %s, coin: %s}", aa.receiver, aa.coin.String())
}
