package ante

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nftTypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Fee constants in units (10^15)
const (
	// completed
	sendCoinFee          = 10
	sendMultiCoinAddFee  = 5
	buyCoinFee           = 100
	sellCoinFee          = 100
	createCoinFee        = 100
	createWalletFee      = 100
	createTransactionFee = 100
	signTransactionFee   = 100
	// future
	declareCandidateFee = 10000
	editCandidateFee    = 10000
	delegateFee         = 200
	unbondFee           = 200
	setOnlineFee        = 100
	setOfflineFee       = 100

	burnFee = 10

	//redeemCheckFee = 30

	htltFee = 33000
)

// Calculate fee in base coin
func CalculateFee(msgs []sdk.Msg, txBytesLen int64, factor sdk.Dec) (sdk.Int, error) {
	commissionInBaseCoin := sdk.ZeroInt()
	commissionInBaseCoin = commissionInBaseCoin.AddRaw(txBytesLen * 2)
	for _, msg := range msgs {
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
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *coinTypes.MsgUpdateCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *coinTypes.MsgBurnCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *multisigTypes.MsgCreateWallet:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(createWalletFee)
		case *multisigTypes.MsgCreateTransaction:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(createTransactionFee)
		case *multisigTypes.MsgSignTransaction:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(signTransactionFee)
		case *nftTypes.MsgMintNFT:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *nftTypes.MsgBurnNFT:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *nftTypes.MsgTransferNFT:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *nftTypes.MsgUpdateReserveNFT:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *nftTypes.MsgEditNFTMetadata:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *govtypes.MsgSubmitProposal:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *govtypes.MsgDeposit:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *govtypes.MsgVote:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *govtypes.MsgVoteWeighted:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
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
