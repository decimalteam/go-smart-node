package keeper_test

import (
	"testing"

	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

type KeeperTestSuite struct {
	suite.Suite

	dsc         *app.DSC
	ctx         sdk.Context
	addrs       []sdk.AccAddress
	vals        []types.Validator
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	dsc := app.Setup(suite.T(), false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	querier := keeper.Querier{Keeper: dsc.ValidatorKeeper}

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, dsc.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	queryClient := types.NewQueryClient(queryHelper)

	suite.msgServer = keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	addrs, _, validators := createValidators(suite.T(), ctx, dsc, []int64{9, 8, 7})
	header := tmproto.Header{
		ChainID: "HelloChain",
		Height:  5,
	}

	// sort a copy of the validators, so that original validators does not
	// have its order changed
	sortedVals := make([]types.Validator, len(validators))
	copy(sortedVals, validators)
	hi := types.NewHistoricalInfo(header, sortedVals)
	dsc.ValidatorKeeper.SetHistoricalInfo(ctx, 5, &hi)

	suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals = dsc, ctx, queryClient, addrs, validators
}

func TestParams(t *testing.T) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	expParams := types.DefaultParams()

	// check that the empty keeper loads the default
	resParams := dsc.ValidatorKeeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))

	// modify a params, save, and retrieve
	expParams.MaxValidators = 777
	dsc.ValidatorKeeper.SetParams(ctx, expParams)
	resParams = dsc.ValidatorKeeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
