package ante

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coin "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	fee "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	multisig "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nft "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	swap "bitbucket.org/decimalteam/go-smart-node/x/swap/types"
)

// Calculate fee in base coin
func CalculateFee(msgs []sdk.Msg, txBytesLen int64, delPrice sdk.Dec, params fee.Params) (sdkmath.Int, error) {
	commissionInBaseCoin := sdk.ZeroInt()
	commissionInBaseCoin = commissionInBaseCoin.AddRaw(txBytesLen * int64(params.ByteFee))
	for _, msg := range msgs {
		switch m := msg.(type) {
		// coin
		case *coin.MsgCreateCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinCreate))
		case *coin.MsgSendCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinSend))
		case *coin.MsgMultiSendCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinSend) + int64((len(m.Sends)-1)*int(params.CoinSendMultiAddition)))
		case *coin.MsgBuyCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinBuy))
		case *coin.MsgSellCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinSell))
		case *coin.MsgSellAllCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.CoinSell))
		case *coin.MsgRedeemCheck:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *coin.MsgUpdateCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *coin.MsgBurnCoin:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		// multisig
		case *multisig.MsgCreateWallet:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.MultisigCreateWallet))
		case *multisig.MsgCreateTransaction:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.MultisigCreateTransaction))
		case *multisig.MsgSignTransaction:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(params.MultisigSignTransaction))
		// swap
		case *swap.MsgSwapInitialize:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *swap.MsgSwapRedeem:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *swap.MsgChainActivate:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *swap.MsgChainDeactivate:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		// nft
		case *nft.MsgMintToken:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *nft.MsgUpdateToken:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *nft.MsgUpdateReserve:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *nft.MsgSendToken:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		case *nft.MsgBurnToken:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(0)
		// fee
		case *fee.MsgSaveBaseDenomPrice:
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
