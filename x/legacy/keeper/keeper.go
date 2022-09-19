package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

var _ types.MsgServer = &Keeper{}

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey store.StoreKey

	bankKeeper     types.BankKeeper
	nftKeeper      types.NftKeeper
	multisigKeeper types.MultisigKeeper

	addressCache map[string]bool
}

// NewKeeper creates new Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey store.StoreKey,
	bankKeeper types.BankKeeper,
	nftKeeper types.NftKeeper,
	multisigKeeper types.MultisigKeeper,
) *Keeper {
	keeper := &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		bankKeeper:     bankKeeper,
		nftKeeper:      nftKeeper,
		multisigKeeper: multisigKeeper,
		addressCache:   make(map[string]bool),
	}
	return keeper
}

func (k *Keeper) RestoreCache(ctx sdk.Context) {
	k.addressCache = make(map[string]bool)
	for _, rec := range k.GetLegacyRecords(ctx) {
		k.addressCache[rec.LegacyAddress] = true
	}
}

func (k *Keeper) IsLegacyAddress(ctx sdk.Context, address string) bool {
	return k.addressCache[address]
}

func (k *Keeper) GetLegacyRecords(ctx sdk.Context) []types.Record {
	var result []types.Record
	store := ctx.KVStore(k.storeKey)
	it := store.Iterator(nil, nil)

	for ; it.Valid(); it.Next() {
		var rec types.Record
		err := k.cdc.UnmarshalLengthPrefixed(it.Value(), &rec)
		if err != nil {
			panic(err)
		}
		result = append(result, rec)
	}
	it.Close()

	return result
}

func (k *Keeper) GetLegacyRecord(ctx sdk.Context, legacyAddress string) (types.Record, error) {
	var result types.Record
	store := ctx.KVStore(k.storeKey)
	key := []byte(legacyAddress)
	value := store.Get(key)
	if len(value) == 0 {
		return result, errors.NotFoundLegacyAddress
	}
	err := k.cdc.UnmarshalLengthPrefixed(value, &result)
	return result, err
}

func (k *Keeper) SetLegacyRecord(ctx sdk.Context, record types.Record) {
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshalLengthPrefixed(&record)
	key := []byte(record.LegacyAddress)
	store.Set(key, value)
}

func (k *Keeper) DeleteLegacyRecord(ctx sdk.Context, legacyAddress string) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(legacyAddress)
	store.Delete(key)
}

func (k *Keeper) ActualizeLegacy(ctx sdk.Context, pubKeyBytes []byte) error {
	legacyAddress, err := commonTypes.GetLegacyAddressFromPubKey(pubKeyBytes)
	if err != nil {
		return errors.CannotGetLegacyAddressFromPublicKey
	}
	if !k.addressCache[legacyAddress] {
		return nil
	}
	actualSdkAddress := sdk.AccAddress(ethsecp256k1.PubKey{Key: pubKeyBytes}.Address())
	actualAddress, err := bech32.ConvertAndEncode(config.Bech32Prefix, actualSdkAddress)
	if err != nil {
		return errors.CannotGetActualAddressFromPublicKey
	}

	record, err := k.GetLegacyRecord(ctx, legacyAddress)
	// error - only if there is no record
	// so just stop here and return
	if err != nil {
		return nil
	}

	// 1. send coins
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.LegacyCoinPool, actualSdkAddress, record.Coins)
	if err != nil {
		return err
	}

	// Emit send event
	err = events.EmitTypedEvent(ctx, &types.EventReturnLegacyCoins{
		LegacyOwner: legacyAddress,
		Owner:       actualAddress,
		Coins:       record.Coins,
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	// 2. update nft owners
	for _, tokenId := range record.NFTs {
		token, found := k.nftKeeper.GetToken(ctx, tokenId)
		if !found {
			continue
		}
		subTokens := k.nftKeeper.GetSubTokens(ctx, tokenId)
		// may be nft already burned
		if len(subTokens) == 0 {
			continue
		}
		returnedSubTokens := make([]uint32, 0)
		for i := range subTokens {
			if subTokens[i].Owner == legacyAddress {
				subTokens[i].Owner = actualAddress
				k.nftKeeper.SetSubToken(ctx, tokenId, subTokens[i])
				returnedSubTokens = append(returnedSubTokens, subTokens[i].ID)
			}
		}
		if len(returnedSubTokens) > 0 {
			// Emit nft event
			err = events.EmitTypedEvent(ctx, &types.EventReturnLegacySubToken{
				LegacyOwner: legacyAddress,
				Owner:       actualAddress,
				Denom:       token.Denom,
				ID:          tokenId,
				SubTokenIDs: returnedSubTokens,
			})
			if err != nil {
				return errors.Internal.Wrapf("err: %s", err.Error())
			}
		}
	}

	// 3. update mutisig wallet owners
	for _, walletAddress := range record.Wallets {
		wallet, err := k.multisigKeeper.GetWallet(ctx, walletAddress)
		// if only wallet not found
		if err != nil {
			continue
		}
		for i, owner := range wallet.Owners {
			if owner == legacyAddress {
				wallet.Owners[i] = actualAddress
			}
		}
		k.multisigKeeper.SetWallet(ctx, wallet)
		// Emit multisig event
		err = events.EmitTypedEvent(ctx, &types.EventReturnMultisigWallet{
			LegacyOwner: legacyAddress,
			Owner:       actualAddress,
			Wallet:      walletAddress,
		})
		if err != nil {
			return errors.Internal.Wrapf("err: %s", err.Error())
		}
	}

	// all complete, delete
	k.DeleteLegacyRecord(ctx, legacyAddress)

	// NOTE: BE CAREFUL WITH CACHES, update only during delivery step
	// NOTE: cache restores every block, so no need to restore after every ActualizeLegacy
	// if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
	//	k.RestoreCache(ctx)
	// }
	return nil
}
