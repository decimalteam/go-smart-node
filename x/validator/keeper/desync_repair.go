package keeper

import (
	"fmt"
	"sort"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Mainnet delegation-desync repair.
//
// Between the 32811496 upgrade and the repair upgrade the node could not decode the
// delegation contract's Stake-carrying events, so stakes, holds, withdrawals and
// redelegations made in that window exist only in the contract. The repair (planned
// in app/desync_repair.go from the contract's own storage) hands the keeper one
// DelegationTarget per (validator, delegator, denom) whose node delegation differs
// from what the contract implies. Amount changes go through Delegate/Unbond so
// validator power and the custom-coin staked totals move exactly as the EVM hooks
// would have moved them.

// EventTypeDelegationDesyncRepair marks a node delegation rewritten by the repair.
const EventTypeDelegationDesyncRepair = "delegation_desync_repair"

// DelegationTarget is the node delegation the contract implies for one key.
type DelegationTarget struct {
	Validator sdk.ValAddress
	Delegator sdk.AccAddress
	Denom     string
	Amount    sdkmath.Int
	// Holds maps a contract hold bucket (holdTimestamp) to its amount. Buckets
	// already matured are ignored when rebuilding node holds.
	Holds map[int64]sdkmath.Int
}

// DesyncRepairReport summarizes a repair run.
type DesyncRepairReport struct {
	Created      int
	Increased    int
	Decreased    int
	Removed      int
	HoldsChanged int
	Skipped      []string
}

func (r DesyncRepairReport) String() string {
	return fmt.Sprintf("created=%d increased=%d decreased=%d removed=%d holdsChanged=%d skipped=%d",
		r.Created, r.Increased, r.Decreased, r.Removed, r.HoldsChanged, len(r.Skipped))
}

// ApplyDelegationTarget rewrites one node delegation to match its target.
func (k Keeper) ApplyDelegationTarget(ctx sdk.Context, t DelegationTarget, rep *DesyncRepairReport) error {
	validator, found := k.GetValidator(ctx, t.Validator)
	if !found {
		rep.Skipped = append(rep.Skipped, fmt.Sprintf("%s/%s/%s: validator not on node", t.Validator, t.Delegator, t.Denom))
		return nil
	}
	current, hasDelegation := k.GetDelegation(ctx, t.Delegator, t.Validator, t.Denom)
	if hasDelegation && current.Stake.Type != types.StakeType_Coin {
		rep.Skipped = append(rep.Skipped, fmt.Sprintf("%s/%s/%s: node delegation is not a coin stake", t.Validator, t.Delegator, t.Denom))
		return nil
	}
	currentAmt := sdkmath.ZeroInt()
	if hasDelegation {
		currentAmt = current.Stake.Stake.Amount
	}
	// One event per repaired delegation, so explorers and indexers that rebuild
	// delegations from node events can follow the repair.
	ctx.EventManager().EmitEvent(sdk.NewEvent(EventTypeDelegationDesyncRepair,
		sdk.NewAttribute("validator", t.Validator.String()),
		sdk.NewAttribute("delegator", t.Delegator.String()),
		sdk.NewAttribute("denom", t.Denom),
		sdk.NewAttribute("old_amount", currentAmt.String()),
		sdk.NewAttribute("new_amount", t.Amount.String()),
	))

	switch {
	case t.Amount.GT(currentAmt):
		delta := types.NewStakeCoin(sdk.Coin{Denom: t.Denom, Amount: t.Amount.Sub(currentAmt)})
		if err := k.Delegate(ctx, t.Delegator, validator, delta); err != nil {
			return err
		}
		if hasDelegation {
			rep.Increased++
		} else {
			rep.Created++
		}
	case t.Amount.LT(currentAmt):
		delta := types.NewStakeCoin(sdk.Coin{Denom: t.Denom, Amount: currentAmt.Sub(t.Amount)})
		remain := types.NewStakeCoin(sdk.Coin{Denom: t.Denom, Amount: t.Amount})
		remain.Holds = current.Stake.Holds
		if err := k.Unbond(ctx, t.Delegator, t.Validator, delta, remain); err != nil {
			return err
		}
		// Unbond leaves the staked total for the unbonding queue to release; nothing
		// is queued here (the contract already holds or returned the tokens).
		k.SubCustomCoinStaked(ctx, delta.Stake)
		if t.Amount.IsZero() {
			rep.Removed++
			return nil
		}
		// Unbond drops a delegation whose remaining base value rounds to zero (a
		// near-worthless custom coin). The contract still holds those tokens, so keep
		// the delegation; it carries no power and the staked total is already right.
		if _, kept := k.GetDelegation(ctx, t.Delegator, t.Validator, t.Denom); !kept {
			k.SetDelegation(ctx, types.NewDelegation(t.Delegator, t.Validator, remain))
			k.IncrementDelegationsCount(ctx, t.Validator)
		}
		rep.Decreased++
	}

	delegation, found := k.GetDelegation(ctx, t.Delegator, t.Validator, t.Denom)
	if !found {
		return nil
	}
	holds, changed := RebuildHolds(delegation.Stake.Holds, t.Holds, ctx.BlockTime().Unix())
	if changed {
		delegation.Stake.Holds = holds
		k.SetDelegation(ctx, delegation)
		rep.HoldsChanged++
	}
	return nil
}

// HoldsMatch reports whether node holds already agree with the target buckets
// (ignoring matured buckets, which DeleteHoldMature prunes on its own).
func HoldsMatch(holds []*types.StakeHold, target map[int64]sdkmath.Int, now int64) bool {
	_, changed := RebuildHolds(cloneHolds(holds), target, now)
	return !changed
}

func cloneHolds(holds []*types.StakeHold) []*types.StakeHold {
	out := make([]*types.StakeHold, len(holds))
	for i, h := range holds {
		c := *h
		out[i] = &c
	}
	return out
}

// holdYear is the minimum hold length that earns the long-hold reward (reward.go).
const holdYear = 365 * 24 * 60 * 60

// restoredHoldStart is the start given to a hold segment the node never recorded. The
// deployed contract does not store when a hold began, so (owner decision D2) the hold
// is treated as committed for a full year: start = end - 365d, never later than now.
func restoredHoldStart(end, now int64) int64 {
	if start := end - holdYear; start < now {
		return start
	}
	return now
}

// RebuildHolds aligns node hold segments with the contract's hold buckets. Matured
// buckets are left alone. A bucket that grew gets a new segment starting at
// restoredHoldStart; a bucket that shrank loses amount the way a held withdrawal
// does. Node segments whose bucket is absent from the contract are removed.
func RebuildHolds(holds []*types.StakeHold, target map[int64]sdkmath.Int, now int64) ([]*types.StakeHold, bool) {
	ends := map[int64]bool{}
	for end := range target {
		ends[end] = true
	}
	for _, h := range holds {
		ends[h.HoldEndTime] = true
	}
	sorted := make([]int64, 0, len(ends))
	for end := range ends {
		sorted = append(sorted, end)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	changed := false
	for _, end := range sorted {
		if end <= now {
			continue
		}
		want, ok := target[end]
		if !ok {
			want = sdkmath.ZeroInt()
		}
		have := sdkmath.ZeroInt()
		for _, h := range holds {
			if h.HoldEndTime == end {
				have = have.Add(h.Amount)
			}
		}
		switch {
		case want.GT(have):
			holds = append(holds, &types.StakeHold{Amount: want.Sub(have), HoldStartTime: restoredHoldStart(end, now), HoldEndTime: end})
			changed = true
		case want.LT(have):
			holds, _, _ = applyTransferredHold(holds, &types.StakeHold{Amount: have.Sub(want), HoldEndTime: end})
			changed = true
		}
	}
	return holds, changed
}
