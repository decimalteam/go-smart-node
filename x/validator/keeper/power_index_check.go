package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// CheckPowerIndexConsistency verifies the ValidatorByPowerIndex is in a state the
// EndBlocker's ApplyAndReturnValidatorSetUpdates transition switch can process
// without hitting its default panic, and that every index key encodes the
// validator's authoritative consensus power (RS.Stake, overlaid into Validator.Stake
// by GetValidator; jailed validators are pinned to 0 by SetValidatorByPowerIndex).
//
// This is the exact C1 chain-halt invariant that the redenomination's
// repowerValidators must preserve: recomputing every validator's power from the
// scaled stakes and re-keying the power index so index-key-power == RS.Stake, never
// leaving an indexed validator in a state the switch cannot classify — notably
// {Unbonded, Online, Stake==0}, which the previous re-power path produced and which
// panics the next block's EndBlocker (see docs/redenom-full-audit-2026-07-02.md R7).
// The live EndBlocker recovers that panic into a truncated, inconsistent validator
// set, so it cannot be relied on to surface the fault; this check classifies each
// indexed validator explicitly instead.
//
// Returns nil when the index is consistent, otherwise an error naming every
// offending validator. It is read-only.
func (k Keeper) CheckPowerIndexConsistency(ctx sdk.Context) error {
	validators, powers, _ := k.GetAllValidatorsByPowerIndex(ctx)
	delCount := k.GetAllDelegationsCount(ctx)

	var panicStates, staleKeys []string
	for i, val := range validators {
		if !endBlockerCanClassify(val, delCount[val.OperatorAddress]) {
			panicStates = append(panicStates, fmt.Sprintf(
				"%s(online=%v,status=%d,stake=%d,dels=%d)",
				val.OperatorAddress, val.Online, val.Status, val.Stake, delCount[val.OperatorAddress]))
		}

		// The index key encodes PotentialConsensusPower()==Validator.Stake at the time
		// SetValidatorByPowerIndex ran; jailed validators are stored at power 0. A
		// mismatch against the current (scaled) RS.Stake means the key was not re-keyed.
		wantPower := val.Stake
		if val.Jailed {
			wantPower = 0
		}
		if powers[i] != wantPower {
			staleKeys = append(staleKeys, fmt.Sprintf(
				"%s(key=%d,stake=%d)", val.OperatorAddress, powers[i], wantPower))
		}
	}

	if len(panicStates) == 0 && len(staleKeys) == 0 {
		return nil
	}
	return fmt.Errorf(
		"power-index inconsistency: %d default-panic state(s) [%s]; %d stale key(s) [%s]",
		len(panicStates), strings.Join(panicStates, ", "),
		len(staleKeys), strings.Join(staleKeys, ", "))
}

// endBlockerCanClassify mirrors the state-transition switch in
// ApplyAndReturnValidatorSetUpdates (val_state_change.go). It returns false for
// exactly the (status, online, stake, delegationsCount) combinations that fall
// through to the switch's default case, which panics. Keep it in lockstep with that
// switch.
func endBlockerCanClassify(val types.Validator, delegationsCount uint32) bool {
	switch {
	case !val.IsUnbonding() && delegationsCount == 0: // -> begin unbonding
		return true
	case val.IsUnbonding() && delegationsCount > 0: // unbonding -> unbonded
		return true
	case val.IsUnbonded() && val.Online && val.Stake > 0: // -> bonded
		return true
	case val.IsBonded() && (!val.Online || val.Stake == 0): // -> unbonded
		return true
	case val.IsBonded() && val.Online && val.Stake > 0: // nothing to do
		return true
	case val.IsUnbonded() && !val.Online && val.Stake == 0: // nothing to do (idle)
		return true
	default:
		return false
	}
}
