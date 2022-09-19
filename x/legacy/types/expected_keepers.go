package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	multisigtypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type NftKeeper interface {
	GetSubTokens(ctx sdk.Context, id string) (subTokens []nfttypes.SubToken)
	SetSubToken(ctx sdk.Context, id string, subToken nfttypes.SubToken)
	GetToken(ctx sdk.Context, id string) (token nfttypes.Token, found bool)
}

type MultisigKeeper interface {
	GetWallet(ctx sdk.Context, address string) (wallet multisigtypes.Wallet, err error)
	SetWallet(ctx sdk.Context, wallet multisigtypes.Wallet)
}
