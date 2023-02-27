package ante

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bank "github.com/cosmos/cosmos-sdk/x/bank/types" // add this for keplr
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coin "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	fee "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	legacy "bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	multisig "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nft "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	swap "bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	validator "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// CalculateFee calculates fee in base coin

func CalculateFee(cdc codec.BinaryCodec, msgs []sdk.Msg, txBytesLen int64, delPrice sdk.Dec, params fee.Params) (sdkmath.Int, error) {

	internalFee := sdkmath.ZeroInt()

	// Do not place commission for tx bytes to end because of RedeemCheck case
	commission := helpers.DecToDecWithE18(params.TxByteFee.MulInt64(txBytesLen))

	msgsFee := sdk.ZeroDec()
	for _, msg := range msgs {
		switch m := msg.(type) {
		// cosmos
		case *bank.MsgSend:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinSend))
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
			// NOTE: for redeem check commission will be payed by check issuer in keeper
			// Here commission will be set to zero to enable redeem for new accounts
			// without coins
			msgsFee = sdk.ZeroDec()
			commission = sdk.ZeroDec()
		case *coin.MsgUpdateCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinUpdate))
		case *coin.MsgBurnCoin:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.CoinBurn))
		// multisig
		case *multisig.MsgCreateWallet:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigCreateWallet))
		case *multisig.MsgCreateTransaction:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.MultisigCreateTransaction))
			var internal sdk.Msg
			err := cdc.UnpackAny(m.Content, &internal)
			if err != nil {
				return sdkmath.ZeroInt(), err
			}
			// calculate fee of internal transaction excluding fee for bytes
			internalFee, err = CalculateFee(cdc, []sdk.Msg{internal}, 0, delPrice, params)
			if err != nil {
				return sdkmath.ZeroInt(), err
			}
		case *multisig.MsgSignTransaction:
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
		// validator
		case *validator.MsgCreateValidator:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorCreateValidator))
		case *validator.MsgEditValidator:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorEditValidator))
		case *validator.MsgSetOnline:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorSetOnline))
		case *validator.MsgSetOffline:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorSetOffline))
		case *validator.MsgDelegate:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorDelegate))
		case *validator.MsgDelegateNFT:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorDelegateNFT))
		case *validator.MsgUndelegate:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorUndelegate))
		case *validator.MsgUndelegateNFT:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorUndelegateNFT))
		case *validator.MsgRedelegate:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorRedelegate))
		case *validator.MsgRedelegateNFT:
			msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorRedelegateNFT))
		/*
			case *validator.MsgCancelUndelegation:
				msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorUndelegate))
			case *validator.MsgCancelUndelegationNFT:
				msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorUndelegateNFT))
			case *validator.MsgCancelRedelegation:
				msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorRedelegate))
			case *validator.MsgCancelRedelegationNFT:
				msgsFee = msgsFee.Add(helpers.DecToDecWithE18(params.ValidatorRedelegateNFT))
		*/
		// fee
		case *fee.MsgUpdateCoinPrices:
		case *upgrade.MsgSoftwareUpgrade:
		case *upgrade.MsgCancelUpgrade:
		// legacy
		case *legacy.MsgReturnLegacy:
		default:
			return sdkmath.ZeroInt(), UnknownTransaction
		}
	}

	commission = commission.Add(msgsFee)

	// change commission according to DEL price
	commissionInBaseCoin := commission.Quo(delPrice).RoundInt()
	commissionInBaseCoin = commissionInBaseCoin.Add(internalFee)

	return commissionInBaseCoin, nil
}
