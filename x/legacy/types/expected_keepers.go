package types

//go:generate mockgen -destination=../testutil/expected_keepers_mock.go -package=testutil . BankKeeper,NftKeeper,MultisigKeeper,ValidatorKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	multisigtypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type NftKeeper interface {
	GetSubTokens(ctx sdk.Context, id string) (subTokens []nfttypes.SubToken)
	ReplaceSubTokenOwner(ctx sdk.Context, id string, index uint32, newOwner string) error
	GetToken(ctx sdk.Context, id string) (token nfttypes.Token, found bool)
}

type MultisigKeeper interface {
	GetWallet(ctx sdk.Context, address string) (wallet multisigtypes.Wallet, err error)
	SetWallet(ctx sdk.Context, wallet multisigtypes.Wallet)
}

type ValidatorKeeper interface {
	GetValidator(ctx sdk.Context, validator sdk.ValAddress) (validatortypes.Validator, bool)
	SetValidator(ctx sdk.Context, val validatortypes.Validator)
}
