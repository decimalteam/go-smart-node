package app

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	validatorkeeper "bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
)

// Owner decision (2026-10-12): node/contract drift that already existed before the
// 32811496 upgrade is left as it was; only the incident is repaired. The list was
// built by `dscd archive-drift` (versioned per-key archive reads at 32811495, the
// last block before the upgrade) against mainnet 34070580:
//
//   - mode "keep": the key differed at 32811495 exactly as it differs now and neither
//     side moved since. Not touched, as long as the planner still sees the recorded
//     node amount and target; if either moved (new activity before H) the key falls
//     back to "offset".
//   - mode "offset": the key already differed at 32811495 and moved since. The
//     incident part is repaired, the old drift kept: target = contract target + drift,
//     where drift = node - contract at 32811495 (raw amounts).

//go:embed delegationupgrade/pre_incident_drift.json
var preIncidentDriftJSON []byte

type preIncidentDrift struct {
	Validator string `json:"v"`
	Delegator string `json:"d"`
	Denom     string `json:"denom"`
	Mode      string `json:"mode"`
	Drift     string `json:"drift"`
	Node      string `json:"node"`
	Target    string `json:"target"`
}

// PreIncidentStats reports how the pre-incident drift list changed the plan.
type PreIncidentStats struct {
	Listed       int
	Kept         int // "keep" keys left untouched
	Offset       int // repaired with the old drift preserved
	FellBack     int // "keep" keys that moved before H, repaired with offset
	NoLongerDiff int // listed keys the planner no longer reports as differing
	HoldsClamped int // offset targets whose hold buckets exceeded the amount
}

func (s PreIncidentStats) String() string {
	return fmt.Sprintf("listed=%d kept=%d offset=%d keepFellBackToOffset=%d notDiffering=%d holdsClamped=%d",
		s.Listed, s.Kept, s.Offset, s.FellBack, s.NoLongerDiff, s.HoldsClamped)
}

func loadPreIncidentDrift() (map[repairKey]preIncidentDrift, error) {
	var list []preIncidentDrift
	if err := json.Unmarshal(preIncidentDriftJSON, &list); err != nil {
		return nil, fmt.Errorf("decode pre-incident drift list: %w", err)
	}
	m := make(map[repairKey]preIncidentDrift, len(list))
	for _, e := range list {
		if e.Mode != "keep" && e.Mode != "offset" {
			return nil, fmt.Errorf("pre-incident drift %s/%s/%s: bad mode %q", e.Validator, e.Delegator, e.Denom, e.Mode)
		}
		m[repairKey{e.Validator, e.Delegator, e.Denom}] = e
	}
	return m, nil
}

func mustInt(s string) (sdkmath.Int, error) {
	v, ok := sdkmath.NewIntFromString(s)
	if !ok {
		return sdkmath.Int{}, fmt.Errorf("bad integer %q", s)
	}
	return v, nil
}

// ApplyPreIncidentDrift rewrites plan.Diffs in place following the list above.
// Deterministic: it only depends on the plan (already ordered) and the embedded list.
func ApplyPreIncidentDrift(plan *RepairPlan) (PreIncidentStats, error) {
	list, err := loadPreIncidentDrift()
	if err != nil {
		return PreIncidentStats{}, err
	}
	st := PreIncidentStats{Listed: len(list)}
	seen := 0
	out := plan.Diffs[:0]
	for _, d := range plan.Diffs {
		k := repairKey{d.Target.Validator.String(), d.Target.Delegator.String(), d.Target.Denom}
		e, ok := list[k]
		if !ok {
			out = append(out, d)
			continue
		}
		seen++
		drift, err := mustInt(e.Drift)
		if err != nil {
			return st, err
		}
		if e.Mode == "keep" {
			node, err := mustInt(e.Node)
			if err != nil {
				return st, err
			}
			target, err := mustInt(e.Target)
			if err != nil {
				return st, err
			}
			if d.NodeAmount.Equal(node) && d.Target.Amount.Equal(target) {
				st.Kept++
				continue
			}
			st.FellBack++
		} else {
			st.Offset++
		}
		amount := d.Target.Amount.Add(drift)
		if amount.IsNegative() {
			amount = sdkmath.ZeroInt()
		}
		d.Target.Amount = amount
		if clampHolds(d.Target.Holds, amount) {
			st.HoldsClamped++
		}
		// The offset can land exactly on the node amount; ApplyDelegationTarget then
		// only reconciles hold buckets.
		out = append(out, d)
	}
	plan.Diffs = out
	st.NoLongerDiff = st.Listed - seen
	return st, nil
}

// ResidualAfterPreIncident applies the pre-incident list to plan and returns the
// diffs the repair would still change (amount or hold buckets). Zero after the
// upgrade means the node matches the contract except for the deliberately kept drift.
func ResidualAfterPreIncident(ctx sdk.Context, app *DSC, plan RepairPlan) ([]RepairDiff, PreIncidentStats, error) {
	p := plan
	p.Diffs = append([]RepairDiff{}, plan.Diffs...)
	st, err := ApplyPreIncidentDrift(&p)
	if err != nil {
		return nil, st, err
	}
	var res []RepairDiff
	for _, d := range p.Diffs {
		if !d.Target.Amount.Equal(d.NodeAmount) {
			res = append(res, d)
			continue
		}
		del, ok := app.ValidatorKeeper.GetDelegation(ctx, d.Target.Delegator, d.Target.Validator, d.Target.Denom)
		if ok && !validatorkeeper.HoldsMatch(del.Stake.Holds, d.Target.Holds, ctx.BlockTime().Unix()) {
			res = append(res, d)
		}
	}
	return res, st, nil
}

// clampHolds trims hold buckets (latest end first) so they do not exceed amount.
func clampHolds(holds map[int64]sdkmath.Int, amount sdkmath.Int) bool {
	total := sdkmath.ZeroInt()
	ends := make([]int64, 0, len(holds))
	for end, a := range holds {
		total = total.Add(a)
		ends = append(ends, end)
	}
	if total.LTE(amount) {
		return false
	}
	sort.Slice(ends, func(i, j int) bool { return ends[i] > ends[j] })
	excess := total.Sub(amount)
	for _, end := range ends {
		if excess.IsZero() {
			break
		}
		a := holds[end]
		cut := sdkmath.MinInt(a, excess)
		if a.Equal(cut) {
			delete(holds, end)
		} else {
			holds[end] = a.Sub(cut)
		}
		excess = excess.Sub(cut)
	}
	return true
}
