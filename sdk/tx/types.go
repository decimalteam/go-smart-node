package tx

import (
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

type (
	MsgCreateCoin          = coinTypes.MsgCreateCoin
	MsgUpdateCoin          = coinTypes.MsgUpdateCoin
	MsgMultiSendCoin       = coinTypes.MsgMultiSendCoin
	MsgBuyCoin             = coinTypes.MsgBuyCoin
	MsgSellCoin            = coinTypes.MsgSellCoin
	MsgSellAllCoin         = coinTypes.MsgSellAllCoin
	MsgSendCoin            = coinTypes.MsgSendCoin
	MsgRedeemCheck         = coinTypes.MsgRedeemCheck
	MsgReturnLegacyBalance = coinTypes.MsgReturnLegacyBalance

	OneSend = coinTypes.Send
)

var (
	NewMsgCreateCoin          = coinTypes.NewMsgCreateCoin
	NewMsgMultiSendCoin       = coinTypes.NewMsgMultiSendCoin
	NewMsgBuyCoin             = coinTypes.NewMsgBuyCoin
	NewMsgSellCoin            = coinTypes.NewMsgSellCoin
	NewMsgSellAllCoin         = coinTypes.NewMsgSellAllCoin
	NewMsgSendCoin            = coinTypes.NewMsgSendCoin
	NewMsgRedeemCheck         = coinTypes.NewMsgRedeemCheck
	NewMsgReturnLegacyBalance = coinTypes.NewMsgReturnLegacyBalance
)
