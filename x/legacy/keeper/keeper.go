package keeper

import (
	"context"
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.MsgServer = &Keeper{}

// Keeper implements the module data storaging.
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
		k.addressCache[rec.Address] = true
	}
}

func (k *Keeper) IsLegacyAddress(ctx sdk.Context, address string) bool {
	return k.addressCache[address]
}

func (k *Keeper) GetLegacyRecords(ctx sdk.Context) []types.LegacyRecord {
	var result []types.LegacyRecord
	store := ctx.KVStore(k.storeKey)
	it := store.Iterator(nil, nil)

	for ; it.Valid(); it.Next() {
		var rec types.LegacyRecord
		err := k.cdc.UnmarshalLengthPrefixed(it.Value(), &rec)
		if err != nil {
			panic(err)
		}
		result = append(result, rec)
	}
	it.Close()

	return result
}

func (k *Keeper) GetLegacyRecord(ctx sdk.Context, legacyAddress string) (types.LegacyRecord, error) {
	var result types.LegacyRecord
	store := ctx.KVStore(k.storeKey)
	key := []byte(legacyAddress)
	value := store.Get(key)
	if len(value) == 0 {
		return result, fmt.Errorf("address '%s' is not found in the key-value store", legacyAddress)
	}
	err := k.cdc.UnmarshalLengthPrefixed(value, &result)
	return result, err
}

func (k *Keeper) SetLegacyRecord(ctx sdk.Context, record types.LegacyRecord) {
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshalLengthPrefixed(&record)
	key := []byte(record.Address)
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
		return types.ErrCannnotGetAddressFromPublicKey(err.Error())
	}
	if !k.addressCache[legacyAddress] {
		return nil
	}
	actualSdkAddress := sdk.AccAddress(ethsecp256k1.PubKey{Key: pubKeyBytes}.Address())
	actualAddress, err := bech32.ConvertAndEncode(config.Bech32Prefix, actualSdkAddress)
	if err != nil {
		return types.ErrCannnotGetAddressFromPublicKey(err.Error())
	}

	record, err := k.GetLegacyRecord(ctx, legacyAddress)
	// only error - there is no record
	// so just stop here and return
	if err != nil {
		return nil
	}

	// 1. send coins
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.LegacyCoinPool, actualSdkAddress, record.Coins)
	if err != nil {
		return types.ErrInternal(err.Error())
	}

	// Emit send event
	err = ctx.EventManager().EmitTypedEvent(&types.EventLegacyReturnCoin{
		OldAddress: legacyAddress,
		NewAddress: actualAddress,
		Coins:      record.Coins.String(),
	})
	if err != nil {
		return types.ErrInternal(err.Error())
	}

	// 2. update nft owners
	for _, nftRecord := range record.Nfts {
		nft, err := k.nftKeeper.GetNFT(ctx, nftRecord.Denom, nftRecord.Id)
		// may be nft already burned
		if err != nil {
			continue
		}
		if nft.ID == "" {
			continue
		}
		for i := range nft.Owners {
			if nft.Owners[i].Address == legacyAddress {
				nft.Owners[i].Address = actualAddress
			}
		}
		k.nftKeeper.SetNFT(ctx, nftRecord.Denom, nftRecord.Id, nft)
		// Emit nft event
		err = ctx.EventManager().EmitTypedEvent(&types.EventLegacyReturnNFT{
			OldAddress: legacyAddress,
			NewAddress: actualAddress,
			Denom:      nftRecord.Denom,
			TokenId:    nftRecord.Id,
		})
		if err != nil {
			return types.ErrInternal(err.Error())
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
		// Emit nft event
		err = ctx.EventManager().EmitTypedEvent(&types.EventLegacyReturnWallet{
			OldAddress: legacyAddress,
			NewAddress: actualAddress,
			Wallet:     walletAddress,
		})
		if err != nil {
			return types.ErrInternal(err.Error())
		}
	}

	// all complete, delete
	k.DeleteLegacyRecord(ctx, legacyAddress)

	// NOTE: BE CAREFUL WITH CACHES, update only during delivery step
	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		k.RestoreCache(ctx)
	}
	return nil
}

// Stub
func (k Keeper) Stub(c context.Context, req *types.QueryStubRequest) (*types.QueryStubResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryStubResponse{}, nil
}
