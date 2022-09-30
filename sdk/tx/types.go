package tx

import (
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	multisigtypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

type (
	FeeParams = feetypes.Params
)

type (
	MsgCreateCoin    = cointypes.MsgCreateCoin
	MsgUpdateCoin    = cointypes.MsgUpdateCoin
	MsgSendCoin      = cointypes.MsgSendCoin
	MultiSendEntry   = cointypes.MultiSendEntry
	MsgMultiSendCoin = cointypes.MsgMultiSendCoin
	MsgBuyCoin       = cointypes.MsgBuyCoin
	MsgSellCoin      = cointypes.MsgSellCoin
	MsgSellAllCoin   = cointypes.MsgSellAllCoin
	MsgBurnCoin      = cointypes.MsgBurnCoin
	MsgRedeemCheck   = cointypes.MsgRedeemCheck

	MsgMintToken     = nfttypes.MsgMintToken
	MsgUpdateToken   = nfttypes.MsgUpdateToken
	MsgUpdateReserve = nfttypes.MsgUpdateReserve
	MsgSendToken     = nfttypes.MsgSendToken
	MsgBurnToken     = nfttypes.MsgBurnToken

	MsgCreateWallet      = multisigtypes.MsgCreateWallet
	MsgCreateTransaction = multisigtypes.MsgCreateTransaction
	MsgSignTransaction   = multisigtypes.MsgSignTransaction

	MsgCreateUniversalTransaction = multisigtypes.MsgCreateUniversalTransaction
	MsgSignUniversalTransaction   = multisigtypes.MsgSignUniversalTransaction
)

var (
	NewMsgCreateCoin    = cointypes.NewMsgCreateCoin
	NewMsgUpdateCoin    = cointypes.NewMsgUpdateCoin
	NewMsgMultiSendCoin = cointypes.NewMsgMultiSendCoin
	NewMsgBuyCoin       = cointypes.NewMsgBuyCoin
	NewMsgSellCoin      = cointypes.NewMsgSellCoin
	NewMsgSellAllCoin   = cointypes.NewMsgSellAllCoin
	NewMsgSendCoin      = cointypes.NewMsgSendCoin
	NewMsgBurnCoin      = cointypes.NewMsgBurnCoin
	NewMsgRedeemCheck   = cointypes.NewMsgRedeemCheck

	NewMsgMintToken     = nfttypes.NewMsgMintToken
	NewMsgUpdateToken   = nfttypes.NewMsgUpdateToken
	NewMsgUpdateReserve = nfttypes.NewMsgUpdateReserve
	NewMsgSendToken     = nfttypes.NewMsgSendToken
	NewMsgBurnToken     = nfttypes.NewMsgBurnToken

	NewMsgCreateWallet      = multisigtypes.NewMsgCreateWallet
	NewMsgCreateTransaction = multisigtypes.NewMsgCreateTransaction
	NewMsgSignTransaction   = multisigtypes.NewMsgSignTransaction

	NewMsgCreateUniversalTransaction = multisigtypes.NewMsgCreateUniversalTransaction
	NewMsgSignUniversalTransaction   = multisigtypes.NewMsgSignUniversalTransaction
)
