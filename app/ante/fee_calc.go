package ante

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coin "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	fee "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	multisig "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nft "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	swap "bitbucket.org/decimalteam/go-smart-node/x/swap/types"
)

// CalculateFee calculates fee in base coin
func CalculateFee(cdc codec.BinaryCodec, msgs []sdk.Msg, txBytesLen int64, delPrice sdk.Dec, params fee.Params) (sdkmath.Int, error) {
	params = fee.DefaultParams()

	msgsFee := sdk.ZeroDec()
	for _, msg := range msgs {
		switch m := msg.(type) {
		// coin
		case *coin.MsgCreateCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinCreate))
		case *coin.MsgSendCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinSend))
		case *coin.MsgMultiSendCoin:
			multiAdditionFee := params.CoinSendAdd.MulInt64(int64(len(m.Sends) - 1))
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinSend.Add(multiAdditionFee)))
		case *coin.MsgBuyCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinBuy))
		case *coin.MsgSellCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinSell))
		case *coin.MsgSellAllCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinSell))
		case *coin.MsgRedeemCheck:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinRedeemCheck))
		case *coin.MsgUpdateCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinUpdate))
		case *coin.MsgBurnCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinBurn))
		// multisig
		case *multisig.MsgCreateWallet:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigCreateWallet))
		case *multisig.MsgCreateTransaction:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigCreateTransaction))
		case *multisig.MsgSignTransaction:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigSignTransaction))
		case *multisig.MsgCreateUniversalTransaction:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigCreateTransaction))
			var internal sdk.Msg
			err := cdc.UnpackAny(m.Content, &internal)
			if err != nil {
				return sdkmath.ZeroInt(), err
			}
			// calculate fee of internal transaction exxcluding fee for bytes
			internalFee, err := CalculateFee(cdc, []sdk.Msg{internal}, 0, delPrice, params)
			if err != nil {
				return sdkmath.ZeroInt(), err
			}
			msgsFee = msgsFee.Add(sdk.NewDecFromInt(internalFee))
		case *multisig.MsgSignUniversalTransaction:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigSignTransaction))
		// swap
		case *swap.MsgInitializeSwap:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.SwapInitialize))
		case *swap.MsgRedeemSwap:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.SwapRedeem))
		case *swap.MsgActivateChain:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.SwapActivateChain))
		case *swap.MsgDeactivateChain:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.SwapDeactivateChain))
		// nft
		case *nft.MsgMintToken:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.NftMintToken))
		case *nft.MsgUpdateToken:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.NftUpdateToken))
		case *nft.MsgUpdateReserve:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.NftUpdateReserve))
		case *nft.MsgSendToken:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.NftSendToken))
		case *nft.MsgBurnToken:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.NftBurnToken))
		// fee
		case *fee.MsgUpdateCoinPrices:
		default:
			return sdkmath.ZeroInt(), UnknownTransaction
		}
	}

	bytesFee := helpers.DecToDecWithE18(params.TxByteFee.MulInt64(txBytesLen))

	commission := bytesFee.Add(msgsFee)

	// change commission according to DEL price
	commissionInBaseCoin := commission.Quo(delPrice).RoundInt()

	return commissionInBaseCoin, nil
}
