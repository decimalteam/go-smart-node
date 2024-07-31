package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Test to demonstrate small remainder after rewards pay
func TestRewards(t *testing.T) {
	t.SkipNow()
	rewards := GetRewardForBlock(40000)
	var powers = []int64{30198431, 30022398, 30000000}
	var totalPower = sdk.NewInt(powers[0] + powers[1] + powers[2])
	remainder := sdk.NewIntFromBigInt(rewards.BigInt())
	for _, power := range powers {
		r := sdk.ZeroInt()
		r = rewards.Mul(sdk.NewInt(power)).Quo(totalPower)
		remainder = remainder.Sub(r)
	}
	t.Fatalf("remainder = %s", remainder)
}
