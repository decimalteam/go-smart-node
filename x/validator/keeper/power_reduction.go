package keeper

import (
	"encoding/binary"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/evmos/ethermint/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// TokensToConsensusPower converts input tokens to potential consensus-engine power
func (k Keeper) TokensToConsensusPower(ctx sdk.Context, tokens sdkmath.Int) int64 {
	return sdk.TokensToConsensusPower(tokens, ethtypes.PowerReduction)
}

// TokensFromConsensusPower - convert input power to tokens
func (k Keeper) TokensFromConsensusPower(ctx sdk.Context, power int64) sdkmath.Int {
	return sdk.TokensFromConsensusPower(power, ethtypes.PowerReduction)
}

// GetValidatorsByPowerIndexKey creates the validator by power index.
// Power index is the key used in the power-store, and represents the relative power ranking of the validator.
func (k Keeper) GetValidatorsByPowerIndexKey(ctx sdk.Context, validator types.Validator) []byte {
	// NOTE the address doesn't need to be stored because counter bytes must always be different
	// NOTE the larger values are of higher value

	consensusPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(consensusPowerBytes, uint64(validator.GetConsensusPower()))

	powerBytes := consensusPowerBytes
	powerBytesLen := len(powerBytes) // 8

	addr := validator.GetOperator()
	operAddrInvr := sdk.CopyBytes(addr)
	addrLen := len(operAddrInvr)

	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}

	// key is of format prefix || powerbytes || addrLen (1byte) || addrBytes
	key := make([]byte, 1+powerBytesLen+1+addrLen)

	key[0] = types.GetValidatorsByPowerIndexKey()[0]
	copy(key[1:powerBytesLen+1], powerBytes)
	key[powerBytesLen+1] = byte(addrLen)
	copy(key[powerBytesLen+2:], operAddrInvr)

	return key
}
