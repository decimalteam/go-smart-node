package keeper_test

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/testvalidator"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// IsValSetSorted reports whether valset is sorted.
func IsValSetSorted(data []types.Validator) bool {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		if types.ValidatorsByVotingPower(data).Less(i, i-1) {
			return false
		}
	}
	return true
}

func TestHistoricalInfoSetGet(t *testing.T) {
	var err error
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 2, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	validators := make([]types.Validator, len(addrVals))

	for i := range addrVals {
		validators[i], err = types.NewValidator(
			addrVals[i],
			addrDels[i],
			PKs[i],
			types.Description{Moniker: "monik"},
			sdk.ZeroDec(),
		)
		validators[i].Stake = int64(100 - i)
		validators[i].Status = types.BondStatus_Bonded
		require.NoError(t, err)
	}

	hi := types.NewHistoricalInfo(ctx.BlockHeader(), validators)
	dsc.ValidatorKeeper.SetHistoricalInfo(ctx, 2, &hi)

	recv, found := dsc.ValidatorKeeper.GetHistoricalInfoDecimal(ctx, 2)
	require.True(t, found, "HistoricalInfo not found after set")
	require.Equal(t, hi, recv, "HistoricalInfo not equal")
	require.True(t, IsValSetSorted(recv.Valset), "HistoricalInfo validators is not sorted: %v \\ %v", hi, recv)

	dsc.ValidatorKeeper.DeleteHistoricalInfo(ctx, 2)

	recv, found = dsc.ValidatorKeeper.GetHistoricalInfoDecimal(ctx, 2)
	require.False(t, found, "HistoricalInfo found after delete")
	require.Equal(t, types.HistoricalInfo{}, recv, "HistoricalInfo is not empty")
}

func TestTrackHistoricalInfo(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 10, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	// set historical entries in params to 5
	params := types.DefaultParams()
	params.HistoricalEntries = 5
	dsc.ValidatorKeeper.SetParams(ctx, params)

	// set historical info at 5, 4 which should be pruned
	// and check that it has been stored
	h4 := tmproto.Header{
		ChainID: "HelloChain",
		Height:  4,
	}
	h5 := tmproto.Header{
		ChainID: "HelloChain",
		Height:  5,
	}
	valSet := []types.Validator{
		testvalidator.NewValidator(t, addrVals[0], PKs[0]),
		testvalidator.NewValidator(t, addrVals[1], PKs[1]),
	}
	hi4 := types.NewHistoricalInfo(h4, valSet)
	hi5 := types.NewHistoricalInfo(h5, valSet)
	dsc.ValidatorKeeper.SetHistoricalInfo(ctx, 4, &hi4)
	dsc.ValidatorKeeper.SetHistoricalInfo(ctx, 5, &hi5)
	recv, found := dsc.ValidatorKeeper.GetHistoricalInfoDecimal(ctx, 4)
	require.True(t, found)
	require.Len(t, recv.Valset, 2)
	require.Equal(t, hi4, recv)
	recv, found = dsc.ValidatorKeeper.GetHistoricalInfoDecimal(ctx, 5)
	require.True(t, found)
	require.Equal(t, hi5, recv)

	// genesis validator
	genesisVals := dsc.ValidatorKeeper.GetAllValidators(ctx)
	require.Len(t, genesisVals, 1)
	// fix
	dsc.ValidatorKeeper.SetLastValidatorPower(ctx, genesisVals[0].GetOperator(), genesisVals[0].Stake)

	// Set bonded validators in keeper
	val1 := testvalidator.NewValidator(t, addrVals[2], PKs[2])
	val1.Status = types.BondStatus_Bonded // when not bonded, consensus power is Zero
	val1.Stake = genesisVals[0].Stake - 1
	dsc.ValidatorKeeper.SetValidator(ctx, val1)
	dsc.ValidatorKeeper.SetValidatorRS(ctx, val1.GetOperator(), types.ValidatorRS{
		Stake: val1.Stake,
	})
	dsc.ValidatorKeeper.SetLastValidatorPower(ctx, val1.GetOperator(), val1.Stake)
	val2 := testvalidator.NewValidator(t, addrVals[3], PKs[3])
	val2.Status = types.BondStatus_Bonded
	val2.Stake = genesisVals[0].Stake + 1
	dsc.ValidatorKeeper.SetValidator(ctx, val2)
	dsc.ValidatorKeeper.SetValidatorRS(ctx, val2.GetOperator(), types.ValidatorRS{
		Stake: val2.Stake,
	})
	dsc.ValidatorKeeper.SetLastValidatorPower(ctx, val2.GetOperator(), val2.Stake)

	vals := []types.Validator{val2, genesisVals[0], val1}
	require.True(t, IsValSetSorted(vals), "powers: %d, %d, %d",
		vals[0].ConsensusPower(),
		vals[1].ConsensusPower(),
		vals[2].ConsensusPower(),
	)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// Set Header for BeginBlock context
	header := tmproto.Header{
		ChainID: "HelloChain",
		Height:  10,
	}
	ctx = ctx.WithBlockHeader(header)

	dsc.ValidatorKeeper.TrackHistoricalInfo(ctx)

	// Check HistoricalInfo at height 10 is persisted
	expected := types.HistoricalInfo{
		Header: header,
		Valset: vals,
	}

	recv, found = dsc.ValidatorKeeper.GetHistoricalInfoDecimal(ctx, 10)
	require.True(t, found, "GetHistoricalInfo failed after BeginBlock")
	require.Equal(t, expected.Header, recv.Header)
	require.Equal(t, len(expected.Valset), len(recv.Valset))
	for i := range expected.Valset {
		ve := expected.Valset[i]
		vr := recv.Valset[i]
		require.True(t, ve.OperatorAddress == vr.OperatorAddress &&
			reflect.DeepEqual(ve.ConsensusPubkey.Value, vr.ConsensusPubkey.Value) &&
			ve.Jailed == vr.Jailed && ve.Stake == vr.Stake,
			"diff at %d, actual powers: %d, %d, %d", i,
			recv.Valset[0].ConsensusPower(),
			recv.Valset[1].ConsensusPower(),
			recv.Valset[2].ConsensusPower(),
		)
	}
	//require.Equal(t, expected, recv, "GetHistoricalInfo returned unexpected result")

	// Check HistoricalInfo at height 5, 4 is pruned
	recv, found = dsc.ValidatorKeeper.GetHistoricalInfoDecimal(ctx, 4)
	require.False(t, found, "GetHistoricalInfo did not prune earlier height")
	require.Equal(t, types.HistoricalInfo{}, recv, "GetHistoricalInfo at height 4 is not empty after prune")
	recv, found = dsc.ValidatorKeeper.GetHistoricalInfoDecimal(ctx, 5)
	require.False(t, found, "GetHistoricalInfo did not prune first prune height")
	require.Equal(t, types.HistoricalInfo{}, recv, "GetHistoricalInfo at height 5 is not empty after prune")
}

func TestGetAllHistoricalInfo(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	// clear historical info
	infos := dsc.ValidatorKeeper.GetAllHistoricalInfo(ctx)
	require.Len(t, infos, 1)
	dsc.ValidatorKeeper.DeleteHistoricalInfo(ctx, infos[0].Header.Height)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 2, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	valSet := []types.Validator{
		testvalidator.NewValidator(t, addrVals[0], PKs[0]),
		testvalidator.NewValidator(t, addrVals[1], PKs[1]),
	}

	header1 := tmproto.Header{ChainID: "HelloChain", Height: 10}
	header2 := tmproto.Header{ChainID: "HelloChain", Height: 11}
	header3 := tmproto.Header{ChainID: "HelloChain", Height: 12}

	hist1 := types.HistoricalInfo{Header: header1, Valset: valSet}
	hist2 := types.HistoricalInfo{Header: header2, Valset: valSet}
	hist3 := types.HistoricalInfo{Header: header3, Valset: valSet}

	expHistInfos := []types.HistoricalInfo{hist1, hist2, hist3}

	for i, hi := range expHistInfos {
		dsc.ValidatorKeeper.SetHistoricalInfo(ctx, int64(10+i), &hi)
	}

	infos = dsc.ValidatorKeeper.GetAllHistoricalInfo(ctx)
	require.Equal(t, expHistInfos, infos)
}
