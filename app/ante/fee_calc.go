package ante

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	coin "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	fee "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	multisig "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nft "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	swap "bitbucket.org/decimalteam/go-smart-node/x/swap/types"
)

var bigE15 = new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil)
var decE15 = sdk.NewDecFromBigInt(bigE15)

// Calculate fee in base coin
func CalculateFee(msgs []sdk.Msg, txBytesLen int64, delPrice sdk.Dec, params fee.Params) (sdkmath.Int, error) {
	commissionInBaseCoin := sdk.ZeroDec()
	commissionInBaseCoin = commissionInBaseCoin.Add(params.TxByteFee.Mul(sdk.NewDec(txBytesLen)))
	for _, msg := range msgs {
		switch m := msg.(type) {
		// coin
		case *coin.MsgCreateCoin:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.CoinCreate)
		case *coin.MsgSendCoin:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.CoinSend)
		case *coin.MsgMultiSendCoin:
			add := params.CoinSendAdd.Mul(sdk.NewDec(int64(len(m.Sends) - 1)))
			commissionInBaseCoin = commissionInBaseCoin.Add(params.CoinSend).Add(add)
		case *coin.MsgBuyCoin:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.CoinBuy)
		case *coin.MsgSellCoin:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.CoinSell)
		case *coin.MsgSellAllCoin:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.CoinSell)
		case *coin.MsgRedeemCheck:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.CoinRedeemCheck)
		case *coin.MsgUpdateCoin:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.CoinUpdate)
		case *coin.MsgBurnCoin:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.CoinBurn)
		// multisig
		case *multisig.MsgCreateWallet:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.MultisigCreateWallet)
		case *multisig.MsgCreateTransaction:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.MultisigCreateTransaction)
		case *multisig.MsgSignTransaction:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.MultisigSignTransaction)
		// swap
		case *swap.MsgInitializeSwap:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.SwapInitialize)
		case *swap.MsgRedeemSwap:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.SwapRedeem)
		case *swap.MsgActivateChain:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.SwapActivateChain)
		case *swap.MsgDeactivateChain:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.SwapDeactivateChain)
		// nft
		case *nft.MsgMintToken:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.NftMintToken)
		case *nft.MsgUpdateToken:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.NftUpdateToken)
		case *nft.MsgUpdateReserve:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.NftUpdateReserve)
		case *nft.MsgSendToken:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.NftSendToken)
		case *nft.MsgBurnToken:
			commissionInBaseCoin = commissionInBaseCoin.Add(params.NftBurnToken)
		// fee
		case *fee.MsgUpdateCoinPrices:
			commissionInBaseCoin = commissionInBaseCoin.Add(sdk.ZeroDec())
		default:
			return sdk.NewInt(0), UnknownTransaction
		}
	}

	// change commission according to DEL price
	commissionInBaseCoin = commissionInBaseCoin.Quo(delPrice).Mul(decE15)
	// TODO: special gas value for special transactions
	return commissionInBaseCoin.RoundInt(), nil
}
