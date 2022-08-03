//go:build norace
// +build norace

package rest_test

import (
	"bitbucket.org/decimalteam/go-smart-node/testutil/network"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/testcoin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"encoding/base64"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

var (
	coins = types.Coins{
		types.Coin{
			Title:       "Cosmos Hub Atom",
			Symbol:      "ATOM",
			CRR:         50,
			Reserve:     sdk.NewInt(1_000_000_000),
			Volume:      sdk.NewInt(1_000_000_000_0),
			LimitVolume: sdk.NewInt(1_000_000_000_000_000_000),
			Creator:     "uatom",
			Identity:    "dx1hs2wdrm87c92rzhq0vgmgrxr6u57xpr2lcygc2",
		},
		types.Coin{
			Title:       "Test Suite Token",
			Symbol:      "TST",
			CRR:         100,
			Reserve:     sdk.NewInt(1_000_000_000),
			Volume:      sdk.NewInt(1_000_000_000_0),
			LimitVolume: sdk.NewInt(1_000_000_000_000_000_000),
			Creator:     "uatom",
			Identity:    "dx1hs2wdrm87c92rzhq0vgmgrxr6u57xpr2lcygc2",
		},
	}

	initVolume, _ = sdk.NewIntFromString("10000000000000000000000000000000000000000")
	params        = types.Params{
		BaseTitle:         "DEL COIN TITLE",
		BaseInitialVolume: initVolume,
		BaseSymbol:        "DEL",
	}

	checks = types.Checks{}
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	genesisState := cfg.GenesisState
	cfg.NumValidators = 1

	checks = types.Checks{
		testcoin.CreateNewCheck(cfg.ChainID, "100000del", "9", "", 1),
		testcoin.CreateNewCheck(cfg.ChainID, "100000del", "10", "", 1),
	}

	var coinGenesis types.GenesisState
	s.Require().NoError(cfg.Codec.UnmarshalJSON(genesisState[types.ModuleName], &coinGenesis))

	coinGenesis.Coins = coins
	coinGenesis.Params = params
	coinGenesis.Checks = checks

	coins = append(coins, types.Coin{
		Title:       params.BaseTitle,
		Symbol:      params.BaseSymbol,
		Volume:      params.BaseInitialVolume,
		CRR:         0,
		Reserve:     sdk.NewInt(0),
		Creator:     "",
		LimitVolume: sdk.NewInt(0),
		Identity:    "",
	})

	coinGenesisBz, err := cfg.Codec.MarshalJSON(&coinGenesis)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = coinGenesisBz
	cfg.GenesisState = genesisState

	s.cfg = cfg

	baseDir, err := ioutil.TempDir(s.T().TempDir(), cfg.ChainID)
	require.NoError(s.T(), err)
	s.T().Logf("created temporary directory: %s", baseDir)

	s.network, err = network.New(s.T(), baseDir, cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(10)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestGetParamsRequestHandlerFn() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name      string
		url       string
		expHeight int64
		respType  fmt.Stringer
		expected  fmt.Stringer
	}{
		{
			name:      "Params",
			url:       fmt.Sprintf("%s/coin/parameters", baseURL),
			expHeight: -1,
			respType:  &types.Params{},
			expected:  &params,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			respJSON, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var resp = rest.ResponseWithHeight{}
			err = val.ClientCtx.LegacyAmino.UnmarshalJSON(respJSON, &resp)
			s.Require().NoError(err)

			// Check height.
			if tc.expHeight >= 0 {
				s.Require().Equal(resp.Height, tc.expHeight)
			} else {
				// To avoid flakiness, just test that height is positive.
				s.Require().Greater(resp.Height, int64(0))
			}

			// Check result.
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(resp.Result, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}

func (s *IntegrationTestSuite) TestQueryGetCoinRequestHandlerFn() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	invalidCoinSymbol := "SOPRANOT"

	testCases := []struct {
		name      string
		url       string
		expHeight int64
		respType  fmt.Stringer
		expected  fmt.Stringer
		error     bool
	}{
		{
			name:      "Get coin with valid denom",
			url:       fmt.Sprintf("%s/coin/coin/%s", baseURL, coins[0].Symbol),
			expHeight: -1,
			respType:  &types.Coin{},
			expected:  &coins[0],
			error:     false,
		},
		{
			name:      "Get coin with invalid denom",
			url:       fmt.Sprintf("%s/coin/coin/%s", baseURL, invalidCoinSymbol),
			expHeight: -1,
			error:     true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			respJSON, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var resp = rest.ResponseWithHeight{}
			err = val.ClientCtx.LegacyAmino.UnmarshalJSON(respJSON, &resp)
			s.Require().NoError(err)

			// Check result.
			err = val.ClientCtx.LegacyAmino.UnmarshalJSON(resp.Result, tc.respType)
			if !tc.error {
				// Check height.
				if tc.expHeight >= 0 {
					s.Require().Equal(resp.Height, tc.expHeight)
				} else {
					// To avoid flakiness, just test that height is positive.
					s.Require().Greater(resp.Height, int64(0))
				}

				s.Require().NoError(err)
				s.Require().Equal(tc.expected.String(), tc.respType.String())
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryGetCheckRequestHandlerFn() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	hash := checks[0].HashFull()
	checkHash1 := base64.URLEncoding.EncodeToString(hash[:])
	invalidHash := "sdkanfiojabd"

	testCases := []struct {
		name      string
		url       string
		expHeight int64
		respType  fmt.Stringer
		expected  fmt.Stringer
		error     bool
	}{
		{
			name:      "Get check by valid hash",
			url:       fmt.Sprintf("%s/coin/check/%s", baseURL, checkHash1),
			expHeight: -1,
			respType:  &types.Check{},
			expected:  &checks[0],
			error:     false,
		},
		{
			name:      "Get check by invalid hash",
			url:       fmt.Sprintf("%s/coin/check/%s", baseURL, invalidHash),
			expHeight: -1,
			error:     true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			respJSON, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var resp = rest.ResponseWithHeight{}
			err = val.ClientCtx.LegacyAmino.UnmarshalJSON(respJSON, &resp)
			s.Require().NoError(err)

			// Check result.
			err = val.ClientCtx.LegacyAmino.UnmarshalJSON(resp.Result, tc.respType)
			if !tc.error {
				// Check height.
				if tc.expHeight >= 0 {
					s.Require().Equal(resp.Height, tc.expHeight)
				} else {
					// To avoid flakiness, just test that height is positive.
					s.Require().Greater(resp.Height, int64(0))
				}

				s.Require().NoError(err)
				s.Require().Equal(tc.expected.String(), tc.respType.String())
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryAllCoinsRequestHandlerFn() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name      string
		url       string
		expHeight int64
		respType  *types.Coins
		expected  types.Coins
	}{
		{
			name:      "Get all coins",
			url:       fmt.Sprintf("%s/coin/all_coins", baseURL),
			expHeight: -1,
			respType:  &types.Coins{},
			expected:  coins,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			respJSON, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var resp = rest.ResponseWithHeight{}
			err = val.ClientCtx.LegacyAmino.UnmarshalJSON(respJSON, &resp)
			s.Require().NoError(err)

			// Check height.
			if tc.expHeight >= 0 {
				s.Require().Equal(resp.Height, tc.expHeight)
			} else {
				// To avoid flakiness, just test that height is positive.
				s.Require().Greater(resp.Height, int64(0))
			}

			// Check result.
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(resp.Result, tc.respType))

			s.Require().True(testcoin.CoinsEqual(coins, *tc.respType))
		})
	}

}

func (s *IntegrationTestSuite) TestQueryAllChecksRequestHandlerFn() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name      string
		url       string
		expHeight int64
		respType  *types.Checks
		expected  types.Checks
	}{
		{
			name:      "Get all checks",
			url:       fmt.Sprintf("%s/coin/all_checks", baseURL),
			expHeight: -1,
			respType:  &types.Checks{},
			expected:  checks,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			respJSON, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var resp = rest.ResponseWithHeight{}
			err = val.ClientCtx.LegacyAmino.UnmarshalJSON(respJSON, &resp)
			s.Require().NoError(err)

			// Check height.
			if tc.expHeight >= 0 {
				s.Require().Equal(resp.Height, tc.expHeight)
			} else {
				// To avoid flakiness, just test that height is positive.
				s.Require().Greater(resp.Height, int64(0))
			}

			// Check result.
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(resp.Result, tc.respType))

			s.Require().True(testcoin.ChecksEqual(checks, *tc.respType))
		})
	}
}
