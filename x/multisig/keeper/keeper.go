package keeper

import (
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/x/multisig/errors"

	"github.com/cometbft/cometbft/libs/log"

	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey store.StoreKey
	ps       paramtypes.Subspace

	accountKeeper auth.AccountKeeperI
	bankKeeper    bank.Keeper
	router        *baseapp.MsgServiceRouter
}

// NewKeeper creates a multisig keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey store.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper auth.AccountKeeperI,
	bankKeeper bank.Keeper,
	router *baseapp.MsgServiceRouter,
) *Keeper {
	return &Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		ps:            ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		router:        router,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

////////////////////////////////////////////////////////////////
// Wallet
////////////////////////////////////////////////////////////////

// GetWallet returns multisig wallet metadata struct with specified address.
func (k *Keeper) GetWallet(ctx sdk.Context, address string) (wallet types.Wallet, err error) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixWallet, []byte(address)...)
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.WalletNotFound
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &wallet)
	return
}

// SetWallet sets the entire wallet metadata struct for a multisig wallet.
func (k *Keeper) SetWallet(ctx sdk.Context, wallet types.Wallet) {
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshalLengthPrefixed(&wallet)
	key := append(types.KeyPrefixWallet, []byte(wallet.Address)...)
	store.Set(key, value)
}

// GetWallets returns multisig wallets metadata struct for specified owner.
func (k *Keeper) GetWallets(ctx sdk.Context, owner string) (wallets []types.Wallet, err error) {
	store := ctx.KVStore(k.storeKey)

	for iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyPrefixWallet)); iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		value := iterator.Value()
		if len(value) == 0 {
			err = errors.EmptyValueInKVStore
			return
		}
		err = k.cdc.UnmarshalLengthPrefixed(value, &wallet)
		if err != nil {
			return
		}
		for _, o := range wallet.Owners {
			if o == owner {
				wallets = append(wallets, wallet)
				break
			}
		}
	}

	return
}

// GetAllWallets returns all multisig wallets metadata.
func (k *Keeper) GetAllWallets(ctx sdk.Context) (wallets []types.Wallet, err error) {
	store := ctx.KVStore(k.storeKey)

	for iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyPrefixWallet)); iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		value := iterator.Value()
		if len(value) == 0 {
			err = errors.EmptyValueInKVStore
			return
		}
		err = k.cdc.UnmarshalLengthPrefixed(value, &wallet)
		if err != nil {
			return
		}
		wallets = append(wallets, wallet)
	}

	return
}

// GetTransactions returns transactions for specified multisig wallet.
func (k *Keeper) GetTransactions(ctx sdk.Context, wallet string) (transactions []types.Transaction, err error) {
	store := ctx.KVStore(k.storeKey)

	for iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyPrefixTransaction)); iterator.Valid(); iterator.Next() {
		var tx types.Transaction
		value := iterator.Value()
		if len(value) == 0 {
			err = errors.EmptyValueInKVStore
			return
		}
		err = k.cdc.UnmarshalLengthPrefixed(value, &tx)
		if err != nil {
			return
		}
		if tx.Wallet == wallet {
			transactions = append(transactions, tx)
		}
	}

	return
}

func (k *Keeper) GetAllTransactions(ctx sdk.Context) (transactions []types.Transaction, err error) {
	store := ctx.KVStore(k.storeKey)

	for iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyPrefixTransaction)); iterator.Valid(); iterator.Next() {
		var tx types.Transaction
		value := iterator.Value()
		if len(value) == 0 {
			err = errors.EmptyValueInKVStore
			return
		}
		err = k.cdc.UnmarshalLengthPrefixed(value, &tx)
		if err != nil {
			return
		}
		transactions = append(transactions, tx)
	}

	return
}

////////////////////////////////////////////////////////////////
// Transaction
////////////////////////////////////////////////////////////////

// SetTransaction sets the entire multisig wallet universal transaction metadata struct for a multisig wallet.
func (k *Keeper) SetTransaction(ctx sdk.Context, transaction types.Transaction) error {
	store := ctx.KVStore(k.storeKey)
	value, err := k.cdc.MarshalLengthPrefixed(&transaction)
	if err != nil {
		return err
	}
	key := append(types.KeyPrefixTransaction, []byte(transaction.Id)...)
	store.Set(key, value)
	return nil
}

func (k *Keeper) GetTransaction(ctx sdk.Context, txID string) (transaction types.Transaction, err error) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixTransaction, []byte(txID)...)
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.TransactionNotFound
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &transaction)
	return
}

// SetSign mark signature for transaction and wallet owner.
func (k *Keeper) SetSign(ctx sdk.Context, txID, signer string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSignatureKey(txID, signer)
	store.Set(key, []byte{})
}

// IsSigned check signature for transaction.
func (k *Keeper) IsSigned(ctx sdk.Context, txID, signer string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetSignatureKey(txID, signer))
}

func (k *Keeper) SetCompleted(ctx sdk.Context, txID string) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCompletedTransaction, []byte(txID)...)
	store.Set(key, []byte{})
}

func (k *Keeper) IsCompleted(ctx sdk.Context, txID string) bool {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCompletedTransaction, []byte(txID)...)
	return store.Has(key)
}

// IsSigned check signature for transaction.
func (k *Keeper) GetSigners(ctx sdk.Context, txID string) []string {
	store := ctx.KVStore(k.storeKey)
	var result []string
	it := sdk.KVStorePrefixIterator(store, types.GetSignaturePrefixKey(txID))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		result = append(result, types.ExtractSignerFromKey(it.Key(), txID))
	}

	return result
}
