package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// TestEndBlockerCanClassify is the R7/C1 regression
// (docs/redenom-full-audit-2026-07-02.md): endBlockerCanClassify must return false
// for exactly the (status, online, stake, delegations) combinations that fall through
// the ApplyAndReturnValidatorSetUpdates transition switch to its default panic — most
// importantly {Unbonded, Online, Stake==0} with delegations, the stale-power-index
// state the redenomination's repowerValidators must never leave behind.
func TestEndBlockerCanClassify(t *testing.T) {
	val := func(status types.BondStatus, online bool, stake int64) types.Validator {
		return types.Validator{Status: status, Online: online, Stake: stake}
	}

	cases := []struct {
		name  string
		v     types.Validator
		dels  uint32
		valid bool
	}{
		// The default-panic states the C1 fix exists to prevent.
		{"unbonded+online+zero+dels PANIC", val(types.BondStatus_Unbonded, true, 0), 1, false},
		{"unbonding+online+zero+no-dels PANIC", val(types.BondStatus_Unbonding, true, 0), 0, false},

		// The six classifiable (safe) states from the switch.
		{"any+offline+no-dels -> begin unbonding", val(types.BondStatus_Bonded, false, 5), 0, true},
		{"unbonding+dels -> unbonded", val(types.BondStatus_Unbonding, false, 5), 1, true},
		{"unbonded+online+power -> bonded", val(types.BondStatus_Unbonded, true, 5), 1, true},
		{"bonded+offline -> unbonded", val(types.BondStatus_Bonded, false, 5), 1, true},
		{"bonded+zero -> unbonded", val(types.BondStatus_Bonded, true, 0), 1, true},
		{"bonded+online+power -> nothing", val(types.BondStatus_Bonded, true, 5), 1, true},
		{"unbonded+offline+zero -> idle nothing", val(types.BondStatus_Unbonded, false, 0), 1, true},
	}
	for _, c := range cases {
		require.Equalf(t, c.valid, endBlockerCanClassify(c.v, c.dels), "case %q", c.name)
	}
}
