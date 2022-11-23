package keeper

import (
	"fmt"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// HandleValidatorSignature handles a validator signature, must be called once per validator per block.
func (k Keeper) HandleValidatorSignature(ctx sdk.Context, addr cryptotypes.Address, power int64, signed bool, params types.Params) {
	logger := k.Logger(ctx)
	height := ctx.BlockHeight()

	// fetch the validator public key
	consAddr := sdk.ConsAddress(addr)
	validator, found := k.GetValidatorByConsAddrDecimal(ctx, consAddr)
	if !found {
		panic(fmt.Sprintf("Validator by consensus address %s not found", sdk.ConsAddress(addr)))
	}

	if validator.IsJailed() {
		return
	}

	// TODO: impossible situation?
	if !validator.Online {
		k.DeleteStartHeight(ctx, consAddr)
		return
	}

	// fetch signing info
	signInfo := k.GetValidatorSigningInfo(ctx, consAddr, height-params.SignedBlocksWindow, height-1)
	if signInfo.StartHeight == -1 {
		logger.Debug("Expected signing info for validator but not found",
			"validator", consAddr.String(),
		)
		return
	}

	// TODO: grace period
	if !signed {
		k.AddMissedBlock(ctx, consAddr, height)
		signInfo.MissedBlocksCounter++

		events.EmitTypedEvent(ctx, &types.EventLiveness{
			Validator:       validator.OperatorAddress,
			ConsensusPubkey: consAddr.String(),
			MissedBlocks:    uint32(signInfo.MissedBlocksCounter),
		})

		logger.Debug(
			"absent validator",
			"height", height,
			"validator", consAddr.String(),
			"missed", signInfo.MissedBlocksCounter,
			"threshold", params.MinSignedPerWindow,
		)
	}

	minHeight := signInfo.StartHeight + params.SignedBlocksWindow
	// max missed = SignedBlocksWindow * (1 - MinSignedPerWindow)
	maxMissed := sdk.NewDec(params.SignedBlocksWindow).Mul(sdk.OneDec().Sub(params.MinSignedPerWindow)).RoundInt64()

	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
		validator := k.ValidatorByConsAddr(ctx, consAddr)
		if validator != nil && !validator.IsJailed() {
			// Downtime confirmed: slash and jail the validator
			// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height,
			// and subtract an additional 1 since this is the LastCommit.
			// Note that this *can* result in a negative "distributionHeight" up to -ValidatorUpdateDelay-1,
			// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
			// That's fine since this is just used to filter unbonding delegations & redelegations.
			distributionHeight := height - sdk.ValidatorUpdateDelay - 1

			k.Slash(ctx, consAddr, distributionHeight, power, params.SlashFractionDowntime)
			k.Jail(ctx, consAddr)

			logger.Info(
				"slashing and jailing validator due to liveness fault",
				"height", height,
				"validator", consAddr.String(),
				"min_height", minHeight,
				"threshold", params.MinSignedPerWindow,
				"slashed", params.SlashFractionDowntime.String(),
				"jailed_until", signInfo.JailedUntil,
			)
		} else {
			// validator was (a) not found or (b) already jailed so we do not slash
			logger.Info(
				"validator would have been slashed for downtime, but was either not found in store or already jailed",
				"validator", consAddr.String(),
			)
		}
	}
}

// handle a validator signing two blocks at the same height
// power: power of the double-signing validator at the height of infraction
func (k Keeper) HandleDoubleSign(ctx sdk.Context, addr crypto.Address, infractionHeight int64, timestamp time.Time, power int64, params types.Params) {
	logger := k.Logger(ctx)

	// calculate the age of the evidence
	t := ctx.BlockHeader().Time
	age := t.Sub(timestamp)

	// fetch the validator public key
	consAddr := sdk.ConsAddress(addr)
	validator, found := k.GetValidatorByConsAddrDecimal(ctx, consAddr)
	if !found {
		panic(fmt.Sprintf("Validator %s not found", consAddr))
	}

	if !validator.Online {
		// Defensive.
		// Simulation doesn't take unbonding periods into account, and
		// Tendermint might break this assumption at some point.
		return
	}

	// fetch the validator signing info
	signInfo := k.GetValidatorSigningInfo(ctx, consAddr, ctx.BlockHeight()-params.SignedBlocksWindow-1, ctx.BlockHeight()-1)
	if signInfo.StartHeight == -1 {
		logger.Debug("Expected signing info for validator but not found",
			"validator", consAddr.String(),
		)
		return
	}

	// validator is already tombstoned
	if signInfo.Tombstoned {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, validator already tombstoned", consAddr, infractionHeight))
		return
	}

	// double sign confirmed
	logger.Info(fmt.Sprintf("Confirmed double sign from %s at height %d, age of %d", consAddr, infractionHeight, age))

	// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height.
	// Note that this *can* result in a negative "distributionHeight", up to -ValidatorUpdateDelay,
	// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
	// That's fine since this is just used to filter unbonding delegations & redelegations.
	distributionHeight := infractionHeight - sdk.ValidatorUpdateDelay

	// Slash validator
	// `power` is the int64 power of the validator as provided to/by
	// Tendermint. This value is validator.Tokens as sent to Tendermint via
	// ABCI, and now received as evidence.
	// The fraction is passed in to separately to slash unbonding and rebonding delegations.
	k.Slash(ctx, consAddr, distributionHeight, power, params.SlashFractionDoubleSign)

	// Jail validator if not already jailed
	// begin unbonding validator if not already unbonding (tombstone)
	if !validator.IsJailed() {
		k.Jail(ctx, consAddr)
	}
}

// fromHeight must bee less toHeight [fromHeight, toHeight]
func (k Keeper) GetValidatorSigningInfo(ctx sdk.Context, addr sdk.ConsAddress, fromHeight, toHeight int64) types.ValidatorSigningInfo {
	store := ctx.KVStore(k.storeKey)
	var result types.ValidatorSigningInfo
	// get ValidatorSigningInfo
	result.StartHeight = k.GetStartHeight(ctx, addr)

	// iterate missing blocks
	var missedBlocksCounter int64

	h := fromHeight
	if h < result.StartHeight {
		h = result.StartHeight
	}
	if h < 0 {
		h = 0
	}

	startKey := types.GetMissedBlockKey(addr, h)
	endKey := types.GetMissedBlockKey(addr, toHeight+1)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		missedBlocksCounter++
	}

	result.MissedBlocksCounter = missedBlocksCounter
	return result
}

func (k Keeper) AddMissedBlock(ctx sdk.Context, addr sdk.ConsAddress, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetMissedBlockKey(addr, height), []byte{})
}

func (k Keeper) SetStartHeight(ctx sdk.Context, addr sdk.ConsAddress, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetStartHeightKey(addr), sdk.Uint64ToBigEndian(uint64(height)))
}

func (k Keeper) GetStartHeight(ctx sdk.Context, addr sdk.ConsAddress) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStartHeightKey(addr))
	if len(bz) == 0 {
		return -1
	}
	return int64(sdk.BigEndianToUint64(bz))
}

func (k Keeper) DeleteStartHeight(ctx sdk.Context, addr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetStartHeightKey(addr))
}
