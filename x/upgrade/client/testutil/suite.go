package testutil

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"fmt"
	"github.com/stretchr/testify/require"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
	"io/ioutil"

	"bitbucket.org/decimalteam/go-smart-node/testutil/network"
	"bitbucket.org/decimalteam/go-smart-node/x/upgrade/client/cli"
	"bitbucket.org/decimalteam/go-smart-node/x/upgrade/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

type IntegrationTestSuite struct {
	suite.Suite

	app     *app.DSC
	cfg     network.Config
	network *network.Network
	ctx     sdk.Context
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	dsc := app.Setup(false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})
	s.app = dsc
	s.ctx = ctx

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	baseDir, err := ioutil.TempDir(s.T().TempDir(), cfg.ChainID)
	require.NoError(s.T(), err)
	s.T().Logf("created temporary directory: %s", baseDir)

	s.network, err = network.New(s.T(), baseDir, cfg)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestModuleVersionsCLI() {
	testCases := []struct {
		msg     string
		req     types.QueryModuleVersionsRequest
		single  bool
		expPass bool
	}{
		{
			msg:     "test full query",
			req:     types.QueryModuleVersionsRequest{ModuleName: ""},
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

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// avoid printing as yaml from CLI command
	clientCtx.OutputFormat = "JSON"

	vm := s.app.UpgradeKeeper.GetModuleVersionMap(s.ctx)
	mv := s.app.UpgradeKeeper.GetModuleVersions(s.ctx)
	s.Require().NotEmpty(vm)

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {

			expect := mv
			if tc.expPass {
				if tc.single {
					expect = []*types.ModuleVersion{{Name: tc.req.ModuleName, Version: vm[tc.req.ModuleName]}}
				}
				// setup expected response
				pm := types.QueryModuleVersionsResponse{
					ModuleVersions: expect,
				}
				jsonVM, _ := clientCtx.Codec.MarshalJSON(&pm)
				expectedRes := string(jsonVM)
				// append new line to match behaviour of PrintProto
				expectedRes += "\n"

				// get actual module versions list response from cli
				cmd := cli.GetModuleVersionsCmd()
				outVM, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{tc.req.ModuleName})
				s.Require().NoError(err)

				s.Require().Equal(expectedRes, outVM.String())
			} else {
				cmd := cli.GetModuleVersionsCmd()
				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{tc.req.ModuleName})
				s.Require().Error(err)
			}
		})
	}
}
