package actions

import (
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

// like SendCoin, but for multisig wallets as recipients
type DepositMultisigWalletGenerator struct {
	bottomRange, upperRange int64 // bounds in 0.001 (10^15)
	knownAddresses          []string
	knownCoins              []string
	knownWallets            []dscApi.MultisigWallet
	rnd                     *rand.Rand
}

type DepositMultisigWalletAction struct {
	recipient string
	coin      sdk.Coin
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
	recipient := gg.knownWallets[i].Address
	return &DepositMultisigWalletAction{
		recipient: recipient,
		coin: sdk.NewCoin(
			RandomChoice(gg.rnd, gg.knownCoins),
			helpers.FinneyToWei(sdkmath.NewInt(RandomRange(gg.rnd, gg.bottomRange, gg.upperRange))),
		),
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
	recipient, err := sdk.AccAddressFromBech32(aa.recipient)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgSendCoin(sender, recipient, aa.coin)

	return feeConfig.MakeTransaction(sa, msg)
}

func (aa *DepositMultisigWalletAction) String() string {
	return fmt.Sprintf("DepositMultisigWallet{recipient: %s, coin: %s}", aa.recipient, aa.coin.String())
}
