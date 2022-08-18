package types

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	// Return a new account with the next account number and the specified address. Does not save the new account to the store.
	NewAccountWithAddress(sdk.Context, sdk.AccAddress) types.AccountI

	// Return a new account with the next account number. Does not save the new account to the store.
	NewAccount(sdk.Context, types.AccountI) types.AccountI

	// Check if an account exists in the store.
	HasAccount(sdk.Context, sdk.AccAddress) bool

	// Retrieve an account from the store.
	GetAccount(sdk.Context, sdk.AccAddress) types.AccountI

	// Set an account in the store.
	SetAccount(sdk.Context, types.AccountI)

	// Remove an account from the store.
	RemoveAccount(sdk.Context, types.AccountI)

	// Iterate over all accounts, calling the provided function. Stop iteration when it returns true.
	IterateAccounts(sdk.Context, func(types.AccountI) bool)

	// Fetch the public key of an account at a specified address
	GetPubKey(sdk.Context, sdk.AccAddress) (cryptotypes.PubKey, error)

	// Fetch the sequence of an account at a specified address.
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)

	// Fetch the next account number, and increment the internal counter.
	GetNextAccountNumber(sdk.Context) uint64
}

// BankKeeper defines the expected bank keeper.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

// CoinKeeper defines the exported coin keeper.
type CoinKeeper interface {
	GetCoin(ctx sdk.Context, symbol string) (coin Coin, err error)
	GetCoins(ctx sdk.Context) (coins []Coin)
	SetCoin(ctx sdk.Context, coin Coin)
	EditCoin(ctx sdk.Context, coin Coin, reserve sdk.Int, volume sdk.Int) error

	IsCheckRedeemed(ctx sdk.Context, check *Check) bool
	GetCheck(ctx sdk.Context, checkHash []byte) (check Check, err error)
	GetChecks(ctx sdk.Context) (checks []Check)
	SetCheck(ctx sdk.Context, check *Check)

	GetParams(ctx sdk.Context) (params Params)
	SetParams(ctx sdk.Context, params Params)

	// need for fee deduction
	IsCoinBase(ctx sdk.Context, symbol string) bool
	GetBaseDenom(ctx sdk.Context) string
	CheckFutureChanges(ctx sdk.Context, coinInfo Coin, amount sdk.Int) error
}
