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
	SetNFT(ctx sdk.Context, denom string, id string, token nfttypes.Token) error
	GetNFT(ctx sdk.Context, denom string, id string) (nfttypes.Token, error)
}

type MultisigKeeper interface {
	GetWallet(ctx sdk.Context, address string) (wallet multisigtypes.Wallet, err error)
	SetWallet(ctx sdk.Context, wallet multisigtypes.Wallet)
}
