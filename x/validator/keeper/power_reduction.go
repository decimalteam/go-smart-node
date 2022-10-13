package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	sdkmath "cosmossdk.io/math"
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/evmos/ethermint/types"
)

// TokensToConsensusPower converts input tokens to potential consensus-engine power
func TokensToConsensusPower(tokens sdkmath.Int) int64 {
	return sdk.TokensToConsensusPower(tokens, ethtypes.PowerReduction)
}

// TokensFromConsensusPower - convert input power to tokens
func TokensFromConsensusPower(power int64) sdkmath.Int {
	return sdk.TokensFromConsensusPower(power, ethtypes.PowerReduction)
}

// GetValidatorByPowerIndexKey creates the validator by power index.
// Power index is the key used in the power-store, and represents the relative power ranking of the validator.
func (k Keeper) GetValidatorByPowerIndexKey(ctx sdk.Context, validator types.Validator, power int64) []byte {
	// NOTE the address doesn't need to be stored because counter bytes must always be different
	// NOTE the larger values are of higher value

	consensusPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(consensusPowerBytes, uint64(power))

	powerBytes := consensusPowerBytes
	powerBytesLen := len(powerBytes) // 8

	operAddrInvr := sdk.CopyBytes(validator.GetOperator())
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