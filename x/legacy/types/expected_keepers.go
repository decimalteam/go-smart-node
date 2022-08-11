package types

import (
	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nftTypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type NftKeeper interface {
	SetNFT(ctx sdk.Context, denom, id string, nft nftTypes.BaseNFT) error
	GetNFT(ctx sdk.Context, denom, id string) (nftTypes.BaseNFT, error)
}

type MultisigKeeper interface {
	GetWallet(ctx sdk.Context, address string) (wallet multisigTypes.Wallet, err error)
	SetWallet(ctx sdk.Context, wallet multisigTypes.Wallet)
}

type LegacyKeeper interface {
	IsLegacyAddress(ctx sdk.Context, address string) bool
	ActualizeLegacy(ctx sdk.Context, pubKeyBytes []byte) error
}
