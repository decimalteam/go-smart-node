package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/testutil/network"
	gocontext "context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/stretchr/testify/suite"
)

type UpgradeTestSuite struct {
	suite.Suite

	cfg         network.Config
	app         *app.DSC
	ctx         sdk.Context
	queryClient types.QueryClient
	//network *network.Network
}

func bootstrapGRPCQueryTest() (*app.DSC, sdk.Context) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	return dsc, ctx
}

func (s *UpgradeTestSuite) SetupTest() {
	dsc, ctx := bootstrapGRPCQueryTest()
	s.app = dsc
	s.ctx = ctx //dsc.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, dsc.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, s.app.UpgradeKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	//s.T().Log("setting up integration test suite")
	//
	//cfg := network.DefaultConfig()
	//genesisState := cfg.GenesisState
	//cfg.NumValidators = 1
	//
	//plans := make([]types.Plan, 0)
	//plans = append(plans, types.Plan{})
	//
	////coinGenesisBz, err := cfg.Codec.MarshalJSON(&coinGenesis)
	////s.Require().NoError(err)
	////genesisState[types.ModuleName] = coinGenesisBz
	//cfg.GenesisState = genesisState
	//
	//s.cfg = cfg
	//
	//baseDir, err := ioutil.TempDir(s.T().TempDir(), cfg.ChainID)
	//require.NoError(s.T(), err)
	//s.T().Logf("created temporary directory: %s", baseDir)
	//
	//s.network, err = network.New(s.T(), baseDir, cfg)
	//s.Require().NoError(err)
	//
	//_, err = s.network.WaitForHeight(10)
	//s.Require().NoError(err)
}

func (suite *UpgradeTestSuite) TestQueryCurrentPlan() {
	var (
		req         *types.QueryCurrentPlanRequest
		expResponse types.QueryCurrentPlanResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"without current upgrade plan",
			func() {
				req = &types.QueryCurrentPlanRequest{}
				expResponse = types.QueryCurrentPlanResponse{}
			},
			true,
		},
		{
			"with current upgrade plan",
			func() {
				plan := types.Plan{Name: "test-plan", Height: 5}
				suite.app.UpgradeKeeper.ScheduleUpgrade(suite.ctx, plan)

				req = &types.QueryCurrentPlanRequest{}
				expResponse = types.QueryCurrentPlanResponse{Plan: &plan}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()

			res, err := suite.queryClient.CurrentPlan(gocontext.Background(), req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(&expResponse, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *UpgradeTestSuite) TestAppliedCurrentPlan() {
	var (
		req       *types.QueryAppliedPlanRequest
		expHeight int64
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"with non-existent upgrade plan",
			func() {
				req = &types.QueryAppliedPlanRequest{Name: "foo"}
			},
			true,
		},
		{
			"with applied upgrade plan",
			func() {
				expHeight = 5

				planName := "test-plan"
				plan := types.Plan{Name: planName, Height: expHeight}
				suite.app.UpgradeKeeper.ScheduleUpgrade(suite.ctx, plan)

				suite.ctx = suite.ctx.WithBlockHeight(expHeight)
				suite.app.UpgradeKeeper.SetUpgradeHandler(planName, func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
					return vm, nil
				})
				suite.app.UpgradeKeeper.ApplyUpgrade(suite.ctx, plan)

				req = &types.QueryAppliedPlanRequest{Name: planName}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()

			res, err := suite.queryClient.AppliedPlan(gocontext.Background(), req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expHeight, res.Height)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *UpgradeTestSuite) TestModuleVersions() {
	testCases := []struct {
		msg     string
		req     types.QueryModuleVersionsRequest
		single  bool
		expPass bool
	}{
		{
			msg:     "test full query",
			req:     types.QueryModuleVersionsRequest{},
			single:  false,
			expPass: true,
		},
		{
			msg:     "test single module",
			req:     types.QueryModuleVersionsRequest{ModuleName: "bank"},
			single:  true,
			expPass: true,
		},
		{
			msg:     "test non-existent module",
			req:     types.QueryModuleVersionsRequest{ModuleName: "abcdefg"},
			single:  true,
			expPass: false,
		},
	}

	vm := suite.app.UpgradeKeeper.GetModuleVersionMap(suite.ctx)
	mv := suite.app.UpgradeKeeper.GetModuleVersions(suite.ctx)

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			res, err := suite.queryClient.ModuleVersions(gocontext.Background(), &tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)

				if tc.single {
					// test that the single module response is valid
					suite.Require().Len(res.ModuleVersions, 1)
					// make sure we got the right values
					suite.Require().Equal(vm[tc.req.ModuleName], res.ModuleVersions[0].Version)
					suite.Require().Equal(tc.req.ModuleName, res.ModuleVersions[0].Name)
				} else {
					// check that the full response is valid
					suite.Require().NotEmpty(res.ModuleVersions)
					suite.Require().Equal(len(mv), len(res.ModuleVersions))
					for i, v := range res.ModuleVersions {
						suite.Require().Equal(mv[i].Version, v.Version)
						suite.Require().Equal(mv[i].Name, v.Name)
					}
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}
