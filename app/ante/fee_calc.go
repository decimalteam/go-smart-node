package ante

import (
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feeTypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	nftTypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Calculate fee in base coin
func CalculateFee(msgs []sdk.Msg, txBytesLen int64, factor sdk.Dec, params feeTypes.Params) (sdk.Int, error) {
	commissionInBaseCoin := sdk.ZeroInt()
	commissionInBaseCoin = commissionInBaseCoin.AddRaw(txBytesLen * int64(params.ByteFee))
	for _, msg := range msgs {
		switch m := msg.(type) {
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
