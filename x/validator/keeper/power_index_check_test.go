package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	keeper "bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// setIndexedValidator persists a validator record + its rewards (RS.Stake becomes
// the validator's overlaid Stake) and inserts a power-index entry whose key encodes
// keyPower. Passing keyPower != rsStake simulates a stale (un-rekeyed) index key.
func setIndexedValidator(
	t *testing.T, k keeper.Keeper, ctx sdk.Context, op sdk.ValAddress, acc sdk.AccAddress,
	status validatortypes.BondStatus, online bool, rsStake, keyPower int64,
) {
	t.Helper()
	val, err := validatortypes.NewValidator(op, acc, PKs[0],
		validatortypes.NewDescription("m", "i", "w", "s", "d"), sdk.ZeroDec())
	require.NoError(t, err)
	val.Status = status
	val.Online = online
	val.DRC20Contract = common.BytesToAddress(val.GetOperator()).String()
	k.SetValidator(ctx, val)
	k.SetValidatorRS(ctx, op, validatortypes.ValidatorRS{
		Rewards: sdk.ZeroInt(), TotalRewards: sdk.ZeroInt(), Stake: rsStake,
	})
	// The index key encodes PotentialConsensusPower()==Stake at set time; force
	// keyPower to model a possibly-stale entry independent of RS.Stake.
	keyVal := val
	keyVal.Stake = keyPower
	k.SetValidatorByPowerIndex(ctx, keyVal)
}

// TestCheckPowerIndexConsistency is the R7/C1 integration regression
// (docs/redenom-full-audit-2026-07-02.md): the checker must pass on a healthy power
// index and fail on (a) an indexed validator in the EndBlocker default-panic state
// {Unbonded, Online, Stake==0, delegations>0} and (b) a stale index key whose
// encoded power != the validator's scaled RS.Stake.
func TestCheckPowerIndexConsistency(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	k := dsc.ValidatorKeeper
	accs, vals := generateAddresses(dsc, ctx, 10, defaultCoins)

	// Healthy: bonded, online, stake>0, index key matches RS.Stake.
	setIndexedValidator(t, k, ctx, vals[0], accs[0],
		validatortypes.BondStatus_Bonded, true, 5, 5)
	require.NoError(t, k.CheckPowerIndexConsistency(ctx),
		"a consistent power index must pass")

	// Stale key: index encodes 5000 (old ×1000 power) but RS.Stake is the scaled 5.
	setIndexedValidator(t, k, ctx, vals[1], accs[1],
		validatortypes.BondStatus_Bonded, true, 5, 5000)
	err := k.CheckPowerIndexConsistency(ctx)
	require.Error(t, err, "a stale index key must be detected")
	require.Contains(t, err.Error(), "stale key")

	// Panic state: {Unbonded, Online, Stake==0} with a delegation — the exact C1
	// EndBlocker default-panic condition. Key power 0 matches RS.Stake 0 (not stale),
	// so this isolates the classify-failure arm.
	setIndexedValidator(t, k, ctx, vals[2], accs[2],
		validatortypes.BondStatus_Unbonded, true, 0, 0)
	k.IncrementDelegationsCount(ctx, vals[2]) // delegationsCount > 0
	err = k.CheckPowerIndexConsistency(ctx)
	require.Error(t, err, "an EndBlocker default-panic state must be detected")
	require.Contains(t, err.Error(), "default-panic")
}
