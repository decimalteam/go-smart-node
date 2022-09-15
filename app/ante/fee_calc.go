package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcTransfer "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibcCoreClientTypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	ibcCoreConnectionTypes "github.com/cosmos/ibc-go/v5/modules/core/03-connection/types"
	ibcCoreChannelTypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feeTypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nftTypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	swapTypes "bitbucket.org/decimalteam/go-smart-node/x/swap/types"
)

// Calculate fee in base coin
func CalculateFee(msgs []sdk.Msg, txBytesLen int64, delPrice sdk.Dec, params feeTypes.Params) (sdk.Int, error) {
	commissionInBaseCoin := sdk.ZeroInt()
	commissionInBaseCoin = commissionInBaseCoin.AddRaw(txBytesLen * int64(params.ByteFee))
	for _, msg := range msgs {
		switch m := msg.(type) {
		//coin
		case *coinTypes.MsgCreateCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinCreate))
		case *coinTypes.MsgSendCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinSend))
		case *coinTypes.MsgMultiSendCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinSend) + int64((len(m.Sends)-1)*int(params.CoinSendMultiAddition)))
		case *coinTypes.MsgBuyCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinBuy))
		case *coinTypes.MsgSellCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinSell))
		case *coinTypes.MsgSellAllCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinSell))
		case *coinTypes.MsgRedeemCheck:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *coinTypes.MsgUpdateCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *coinTypes.MsgBurnCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		// multisig
		case *multisigTypes.MsgCreateWallet:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.MultisigCreateWallet))
		case *multisigTypes.MsgCreateTransaction:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.MultisigCreateTransaction))
		case *multisigTypes.MsgSignTransaction:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.MultisigSignTransaction))
		case *swapTypes.MsgSwapInitialize:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *swapTypes.MsgSwapRedeem:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *swapTypes.MsgChainActivate:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *swapTypes.MsgChainDeactivate:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		// nft
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

		// fee
		case *feeTypes.MsgSaveBaseDenomPrice:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)

		// ibc client
		case *ibcCoreClientTypes.MsgCreateClient:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreClientTypes.MsgUpdateClient:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreClientTypes.MsgSubmitMisbehaviour:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreClientTypes.MsgUpgradeClient:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)

		// ibc connection
		case *ibcCoreConnectionTypes.MsgConnectionOpenInit:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreConnectionTypes.MsgConnectionOpenConfirm:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreConnectionTypes.MsgConnectionOpenAck:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreConnectionTypes.MsgConnectionOpenTry:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)

		//ibc channel
		case *ibcCoreChannelTypes.MsgChannelOpenInit:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreChannelTypes.MsgChannelOpenConfirm:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreChannelTypes.MsgChannelOpenAck:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreChannelTypes.MsgChannelOpenTry:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreChannelTypes.MsgChannelCloseInit:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *ibcCoreChannelTypes.MsgChannelCloseConfirm:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)

		// ibc transfer
		case *ibcTransfer.MsgTransfer:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		default:
			return sdk.NewInt(0), UnknownTransaction
		}
	}

	commissionInBaseCoin = helpers.FinneyToWei(commissionInBaseCoin)
	// change commission according to DEL price
	commissionInBaseCoin = sdk.OneDec().MulInt(commissionInBaseCoin).Quo(delPrice).RoundInt()
	// TODO: special gas value for special transactions
	return commissionInBaseCoin, nil
}
