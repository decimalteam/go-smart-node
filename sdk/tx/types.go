package tx

import (
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nftTypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

type (
	MsgCreateCoin    = coinTypes.MsgCreateCoin
	MsgUpdateCoin    = coinTypes.MsgUpdateCoin
	MsgMultiSendCoin = coinTypes.MsgMultiSendCoin
	MsgBuyCoin       = coinTypes.MsgBuyCoin
	MsgSellCoin      = coinTypes.MsgSellCoin
	MsgSellAllCoin   = coinTypes.MsgSellAllCoin
	MsgSendCoin      = coinTypes.MsgSendCoin
	MsgBurnCoin      = coinTypes.MsgBurnCoin
	MsgRedeemCheck   = coinTypes.MsgRedeemCheck

	OneSend = coinTypes.Send

	MsgMintNFT          = nftTypes.MsgMintNFT
	MsgBurnNFT          = nftTypes.MsgBurnNFT
	MsgUpdateReserveNFT = nftTypes.MsgUpdateReserveNFT
	MsgTransferNFT      = nftTypes.MsgTransferNFT
	MsgEditNFTMetadata  = nftTypes.MsgEditNFTMetadata

	MsgCreateWallet      = multisigTypes.MsgCreateWallet
	MsgCreateTransaction = multisigTypes.MsgCreateTransaction
	MsgSignTransaction   = multisigTypes.MsgSignTransaction
)

var (
	NewMsgCreateCoin    = coinTypes.NewMsgCreateCoin
	NewMsgUpdateCoin    = coinTypes.NewMsgUpdateCoin
	NewMsgMultiSendCoin = coinTypes.NewMsgMultiSendCoin
	NewMsgBuyCoin       = coinTypes.NewMsgBuyCoin
	NewMsgSellCoin      = coinTypes.NewMsgSellCoin
	NewMsgSellAllCoin   = coinTypes.NewMsgSellAllCoin
	NewMsgSendCoin      = coinTypes.NewMsgSendCoin
	NewMsgBurnCoin      = coinTypes.NewMsgBurnCoin
	NewMsgRedeemCheck   = coinTypes.NewMsgRedeemCheck

	NewMsgMintNFT          = nftTypes.NewMsgMintNFT
	NewMsgBurnNFT          = nftTypes.NewMsgBurnNFT
	NewMsgUpdateReserveNFT = nftTypes.NewMsgUpdateReserveNFT
	NewMsgTransferNFT      = nftTypes.NewMsgTransferNFT
	NewMsgEditNFTMetadata  = nftTypes.NewMsgEditNFTMetadata

	NewMsgCreateWallet      = multisigTypes.NewMsgCreateWallet
	NewMsgCreateTransaction = multisigTypes.NewMsgCreateTransaction
	NewMsgSignTransaction   = multisigTypes.NewMsgSignTransaction
)
