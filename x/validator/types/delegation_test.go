package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func TestDelegationEqual(t *testing.T) {
	d1 := types.NewDelegation(sdk.AccAddress(valAddr1), valAddr2, types.NewStakeCoin(sdk.NewInt64Coin("aaa", 10)))
	d2 := d1

	ok := d1.String() == d2.String()
	require.True(t, ok)

	d2.Validator = valAddr3.String()
	d2.Stake.Stake.Amount = sdk.NewInt(200)

	ok = d1.String() == d2.String()
	require.False(t, ok)
}

func TestDelegationString(t *testing.T) {
	d := types.NewDelegation(sdk.AccAddress(valAddr1), valAddr2, types.NewStakeCoin(sdk.NewInt64Coin("aaa", 10)))
	require.NotEmpty(t, d.String())
}

func TestUnbondingDelegationEqual(t *testing.T) {
	ubd1 := types.NewUndelegation(sdk.AccAddress(valAddr1), valAddr2, 0,
		time.Unix(0, 0), types.NewStakeCoin(sdk.NewInt64Coin("aaa", 10)))
	ubd2 := ubd1

	ok := ubd1.String() == ubd2.String()
	require.True(t, ok)

	ubd2.Validator = valAddr3.String()

	ubd2.Entries[0].CompletionTime = time.Unix(20*20*2, 0)
	ok = (ubd1.String() == ubd2.String())
	require.False(t, ok)
}

func TestUnbondingDelegationString(t *testing.T) {
	ubd := types.NewUndelegation(sdk.AccAddress(valAddr1), valAddr2, 0,
		time.Unix(0, 0), types.NewStakeCoin(sdk.NewInt64Coin("aaa", 10)))

	require.NotEmpty(t, ubd.String())
}

func TestRedelegationEqual(t *testing.T) {
	r1 := types.NewRedelegation(sdk.AccAddress(valAddr1), valAddr2, valAddr3, 0,
		time.Unix(0, 0), types.NewStakeCoin(sdk.NewInt64Coin("aaa", 10)))
	r2 := types.NewRedelegation(sdk.AccAddress(valAddr1), valAddr2, valAddr3, 0,
		time.Unix(0, 0), types.NewStakeCoin(sdk.NewInt64Coin("aaa", 10)))

	ok := r1.String() == r2.String()
	require.True(t, ok)

	r2.Entries[0].CompletionTime = time.Unix(20*20*2, 0)

	ok = r1.String() == r2.String()
	require.False(t, ok)
}

func TestRedelegationString(t *testing.T) {
	r := types.NewRedelegation(sdk.AccAddress(valAddr1), valAddr2, valAddr3, 0,
		time.Unix(0, 0), types.NewStakeCoin(sdk.NewInt64Coin("aaa", 10)))

	require.NotEmpty(t, r.String())
}
