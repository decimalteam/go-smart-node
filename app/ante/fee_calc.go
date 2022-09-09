package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feeTypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nftTypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	swapTypes "bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	sdkmath "cosmossdk.io/math"
)

// CalculateFee calculates fee in base coin
func CalculateFee(msgs []sdk.Msg, txBytesLen int64, delPrice sdk.Dec, params feeTypes.Params) (sdkmath.Int, error) {
	params = feeTypes.DefaultParams()

	var msgsFee sdk.Dec
	for _, msg := range msgs {
		switch m := msg.(type) {
		// coin
		case *coinTypes.MsgCreateCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinCreate))
		case *coinTypes.MsgSendCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinSend))
		case *coinTypes.MsgMultiSendCoin:
			multiAdditionFee := params.CoinSendMultiAddition.MulInt64(int64(len(m.Sends) - 1))
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinSend.Add(multiAdditionFee)))
		case *coinTypes.MsgBuyCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinBuy))
		case *coinTypes.MsgSellCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinSell))
		case *coinTypes.MsgSellAllCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinSell))
		case *coinTypes.MsgRedeemCheck:
		case *coinTypes.MsgUpdateCoin:
		case *coinTypes.MsgBurnCoin:
		// multisig
		case *multisigTypes.MsgCreateWallet:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigCreateWallet))
		case *multisigTypes.MsgCreateTransaction:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigCreateTransaction))
		case *multisigTypes.MsgSignTransaction:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigSignTransaction))
		case *swapTypes.MsgSwapInitialize:
		case *swapTypes.MsgSwapRedeem:
		case *swapTypes.MsgChainActivate:
		case *swapTypes.MsgChainDeactivate:
		// nft
		case *nftTypes.MsgMintNFT:
		case *nftTypes.MsgBurnNFT:
		case *nftTypes.MsgTransferNFT:
		case *nftTypes.MsgUpdateReserveNFT:
		case *nftTypes.MsgEditNFTMetadata:

		// fee
		case *feeTypes.MsgSaveBaseDenomPrice:
		default:
			return sdkmath.ZeroInt(), UnknownTransaction
		}
	}

	bytesFee := helpers.DecToDecWithE18(params.ByteFee.MulInt64(txBytesLen))

	commission := bytesFee.Add(msgsFee)

	// change commission according to DEL price
	commissionInBaseCoin := commission.Quo(delPrice).RoundInt()

	return commissionInBaseCoin, nil
}
