package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	cosmosAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   store.StoreKey
	paramSpace paramtypes.Subspace

	accountKeeper auth.AccountKeeperI
	bankKeeper    bank.Keeper
}

// NewKeeper creates new Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey store.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper auth.AccountKeeperI,
	bankKeeper bank.Keeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	keeper := &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
	return keeper
}

// Returns true if account balance enought to send coins
func (k Keeper) CheckBalance(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) (bool, error) {
	balance := k.bankKeeper.GetAllBalances(ctx, address)

	return balance.IsAllGTE(coins), nil
}

func (k Keeper) UnlockFunds(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.SwapPool, address, coins)
}

func (k Keeper) LockFunds(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.SwapPool, coins)
}

func (k Keeper) CheckPoolFunds(ctx sdk.Context, coins sdk.Coins) bool {
	poolAddr := cosmosAuthTypes.NewModuleAddress(types.SwapPool)
	balance := k.bankKeeper.GetAllBalances(ctx, poolAddr)

	return balance.IsAllGTE(coins)
}

func (k Keeper) GetLockedFunds(ctx sdk.Context) sdk.Coins {
	poolAddr := cosmosAuthTypes.NewModuleAddress(types.SwapPool)
	return k.bankKeeper.GetAllBalances(ctx, poolAddr)
}

////////////////////////////////////////////////////////////////
// Params
////////////////////////////////////////////////////////////////

// GetParams returns the total set of the module parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the module parameters to the param space.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

////////////////////////////////////////////////////////////////
// Chain
////////////////////////////////////////////////////////////////

func (k Keeper) HasChain(ctx sdk.Context, chainNumber uint32) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetChainKey(chainNumber))
}

func (k *Keeper) SetChain(ctx sdk.Context, chain *types.Chain) {
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshalLengthPrefixed(chain)
	store.Set(types.GetChainKey(chain.Id), value)
}

func (k *Keeper) GetChain(ctx sdk.Context, chainNumber uint32) (types.Chain, bool) {
	var chain types.Chain

	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetChainKey(chainNumber))
	if len(value) == 0 {
		return types.Chain{}, false
	}
	k.cdc.MustUnmarshalLengthPrefixed(value, &chain)
	return chain, true
}

////////////////////////////////////////////////////////////////
// Swap
////////////////////////////////////////////////////////////////

func (k Keeper) SetSwap(ctx sdk.Context, hash types.Hash) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSwapKey(hash), []byte{})
}

func (k Keeper) HasSwap(ctx sdk.Context, hash types.Hash) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetSwapKey(hash))
}
