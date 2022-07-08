package ante

import (
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Fee constants in units (10^15)
const (
	// completed
	sendCoinFee         = 10
	sendMultiCoinAddFee = 5
	buyCoinFee          = 100
	sellCoinFee         = 100
	createCoinFee       = 100
	// future
	declareCandidateFee = 10000
	editCandidateFee    = 10000
	delegateFee         = 200
	unbondFee           = 200
	setOnlineFee        = 100
	setOfflineFee       = 100

	burnFee = 10

	//redeemCheckFee = 30

	createWalletFee      = 100
	createTransactionFee = 100
	signTransactionFee   = 100

	htltFee = 33000
)

// Calculate fee in base coin
func CalculateFee(tx sdk.Tx, txBytesLen int64, factor sdk.Dec) (sdk.Int, error) {
	commissionInBaseCoin := sdk.ZeroInt()
	commissionInBaseCoin = commissionInBaseCoin.AddRaw(txBytesLen * 2)
	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		case *coinTypes.MsgCreateCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(createCoinFee)
		case *coinTypes.MsgSendCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sendCoinFee)
		case *coinTypes.MsgMultiSendCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sendCoinFee + int64((len(m.Sends)-1)*sendMultiCoinAddFee))
		case *coinTypes.MsgBuyCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(buyCoinFee)
		case *coinTypes.MsgSellCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sellCoinFee)
		case *coinTypes.MsgSellAllCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sellCoinFee)
		case *coinTypes.MsgRedeemCheck:
			commissionInBaseCoin = sdk.ZeroInt()
		case *coinTypes.MsgUpdateCoin:
			commissionInBaseCoin = sdk.ZeroInt()
		default:
			return sdk.NewInt(0), ErrUnknownTransaction(fmt.Sprintf("%T", msg))
		}
	}

	commissionInBaseCoin = helpers.FinneyToWei(commissionInBaseCoin)
	// change commission according to factor
	commissionInBaseCoin = factor.MulInt(commissionInBaseCoin).RoundInt()
	// TODO: special gas value for special transactions
	return commissionInBaseCoin, nil
}
