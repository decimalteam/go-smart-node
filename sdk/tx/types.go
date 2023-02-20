package tx

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	multisigtypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type (
	FeeParams = feetypes.Params
)

type FeeCalculationOptions struct {
	DelPrice  sdk.Dec
	FeeParams feetypes.Params
	AppCodec  codec.BinaryCodec
}

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

	MsgCreateValidator       = validatortypes.MsgCreateValidator
	MsgEditValidator         = validatortypes.MsgEditValidator
	MsgSetOnline             = validatortypes.MsgSetOnline
	MsgSetOffline            = validatortypes.MsgSetOffline
	MsgDelegate              = validatortypes.MsgDelegate
	MsgDelegateNFT           = validatortypes.MsgDelegateNFT
	MsgUndelegate            = validatortypes.MsgUndelegate
	MsgUndelegateNFT         = validatortypes.MsgUndelegateNFT
	MsgRedelegate            = validatortypes.MsgRedelegate
	MsgRedelegateNFT         = validatortypes.MsgRedelegateNFT
	MsgCancelUndelegation    = validatortypes.MsgCancelUndelegation
	MsgCancelUndelegationNFT = validatortypes.MsgCancelUndelegationNFT
	MsgCancelRedelegation    = validatortypes.MsgCancelRedelegation
	MsgCancelRedelegationNFT = validatortypes.MsgCancelRedelegationNFT

	Description = validatortypes.Description

	MsgSoftwareUpgrade = upgradetypes.MsgSoftwareUpgrade
	MsgCancelUpgrade   = upgradetypes.MsgCancelUpgrade
	Plan               = upgradetypes.Plan
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

	NewMsgCreateValidator       = validatortypes.NewMsgCreateValidator
	NewMsgEditValidator         = validatortypes.NewMsgEditValidator
	NewMsgSetOnline             = validatortypes.NewMsgSetOnline
	NewMsgSetOffline            = validatortypes.NewMsgSetOffline
	NewMsgDelegate              = validatortypes.NewMsgDelegate
	NewMsgDelegateNFT           = validatortypes.NewMsgDelegateNFT
	NewMsgUndelegate            = validatortypes.NewMsgUndelegate
	NewMsgUndelegateNFT         = validatortypes.NewMsgUndelegateNFT
	NewMsgRedelegate            = validatortypes.NewMsgRedelegate
	NewMsgRedelegateNFT         = validatortypes.NewMsgRedelegateNFT
	NewMsgCancelUndelegation    = validatortypes.NewMsgCancelUndelegation
	NewMsgCancelUndelegationNFT = validatortypes.NewMsgCancelUndelegationNFT
	NewMsgCancelRedelegation    = validatortypes.NewMsgCancelRedelegation
	NewMsgCancelRedelegationNFT = validatortypes.NewMsgCancelRedelegationNFT
)
