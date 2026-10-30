package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	validatorkeeper "bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
)

func TestPreIncidentDriftListDecodes(t *testing.T) {
	m, err := loadPreIncidentDrift()
	if err != nil {
		t.Fatal(err)
	}
	keep, offset := 0, 0
	for _, e := range m {
		if e.Mode == "keep" {
			keep++
		} else {
			offset++
		}
		for _, s := range []string{e.Drift, e.Node, e.Target} {
			if _, err := mustInt(s); err != nil {
				t.Fatal(err)
			}
		}
	}
	if keep != 3718 || offset != 203 {
		t.Fatalf("keep=%d offset=%d", keep, offset)
	}
}

func TestApplyPreIncidentDriftRules(t *testing.T) {
	list, err := loadPreIncidentDrift()
	if err != nil {
		t.Fatal(err)
	}
	var keep, offset preIncidentDrift
	for _, e := range list {
		if e.Mode == "keep" && keep.Mode == "" {
			keep = e
		}
		if e.Mode == "offset" && offset.Mode == "" {
			offset = e
		}
	}
	diffFor := func(e preIncidentDrift, node, target sdkmath.Int) RepairDiff {
		val, _ := sdk.ValAddressFromBech32(e.Validator)
		del, _ := sdk.AccAddressFromBech32(e.Delegator)
		return RepairDiff{
			Target:     validatorkeeper.DelegationTarget{Validator: val, Delegator: del, Denom: e.Denom, Amount: target, Holds: map[int64]sdkmath.Int{}},
			NodeAmount: node,
		}
	}
	n, _ := mustInt(keep.Node)
	tg, _ := mustInt(keep.Target)
	od, _ := mustInt(offset.Drift)
	ot, _ := mustInt(offset.Target)
	unlisted := diffFor(preIncidentDrift{Validator: keep.Validator, Delegator: keep.Delegator, Denom: "not-listed"}, sdkmath.ZeroInt(), sdkmath.NewInt(5))

	plan := RepairPlan{Diffs: []RepairDiff{
		diffFor(keep, n, tg),                   // unchanged since the snapshot: kept
		diffFor(offset, sdkmath.ZeroInt(), ot), // offset
		unlisted,                               // untouched by the list
	}}
	st, err := ApplyPreIncidentDrift(&plan)
	if err != nil {
		t.Fatal(err)
	}
	if st.Kept != 1 || st.Offset != 1 || len(plan.Diffs) != 2 {
		t.Fatalf("stats %s, diffs %d", st, len(plan.Diffs))
	}
	want := ot.Add(od)
	if want.IsNegative() {
		want = sdkmath.ZeroInt()
	}
	if !plan.Diffs[0].Target.Amount.Equal(want) {
		t.Fatalf("offset target %s, want %s", plan.Diffs[0].Target.Amount, want)
	}
	if plan.Diffs[1].Target.Denom != "not-listed" || !plan.Diffs[1].Target.Amount.Equal(sdkmath.NewInt(5)) {
		t.Fatal("unlisted diff changed")
	}

	// A "keep" key that moved before H falls back to offset.
	moved := RepairPlan{Diffs: []RepairDiff{diffFor(keep, n, tg.AddRaw(1))}}
	st, err = ApplyPreIncidentDrift(&moved)
	if err != nil {
		t.Fatal(err)
	}
	kd, _ := mustInt(keep.Drift)
	want = tg.AddRaw(1).Add(kd)
	if want.IsNegative() {
		want = sdkmath.ZeroInt()
	}
	if st.FellBack != 1 || len(moved.Diffs) != 1 || !moved.Diffs[0].Target.Amount.Equal(want) {
		t.Fatalf("fallback: %s target %s want %s", st, moved.Diffs[0].Target.Amount, want)
	}
}

func TestClampHolds(t *testing.T) {
	h := map[int64]sdkmath.Int{100: sdkmath.NewInt(30), 200: sdkmath.NewInt(50)}
	if !clampHolds(h, sdkmath.NewInt(60)) {
		t.Fatal("expected clamp")
	}
	if !h[100].Equal(sdkmath.NewInt(30)) || !h[200].Equal(sdkmath.NewInt(30)) {
		t.Fatalf("got %v", h)
	}
	if clampHolds(h, sdkmath.NewInt(60)) {
		t.Fatal("no clamp expected")
	}
	if !clampHolds(h, sdkmath.NewInt(10)) || len(h) != 1 || !h[100].Equal(sdkmath.NewInt(10)) {
		t.Fatalf("got %v", h)
	}
}
