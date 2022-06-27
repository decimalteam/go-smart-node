package tests

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin/testcoin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/testutil"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gogo/protobuf/proto"
)

func (s *IntegrationTestSuite) TestGetParamsGRPCHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			name: "Get params on grpc",
			url:  fmt.Sprintf("%s/coin/v1/params", baseURL),
			headers: map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			respType: &types.QueryParamsResponse{},
			expected: &types.QueryParamsResponse{
				Params: params,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}

func (s *IntegrationTestSuite) TestGetCheckByHashGrpcHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	hash := checks[0].HashFull()
	checkHash1 := base64.URLEncoding.EncodeToString(hash[:])

	testCases := []struct {
		name     string
		url      string
		respType proto.Message
		expected proto.Message
		error    bool
	}{
		{
			name:     "Get check by hash on grpc with valid hash param",
			url:      baseURL + "/coin/v1/check/" + checkHash1,
			respType: &types.QueryCheckResponse{},
			expected: &types.QueryCheckResponse{
				Check: checks[0],
			},
			error: false,
		},
		{
			name:  "Get check by hash on grpc with invalid param",
			url:   fmt.Sprintf("%s/coin/v1/check/%s", baseURL, "skjadfn"),
			error: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			err = val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType)
			if !tc.error {
				s.Require().NoError(err)
				s.Require().Equal(tc.expected.String(), tc.respType.String())
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCoinByDenomGrpcHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
		error    bool
	}{
		{
			name: "test query get coin by denom with valid param",
			url:  fmt.Sprintf("%s/coin/v1/coin/%s", baseURL, coins[0].Symbol),
			headers: map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			respType: &types.QueryCoinResponse{},
			expected: &types.QueryCoinResponse{
				Coin: coins[0],
			},
			error: false,
		},

		{
			name: "test query get coin by denom with invalid param",
			url:  fmt.Sprintf("%s/coin/v1/coin/%s", baseURL, "sdfj"),
			headers: map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			error: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			err = val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType)

			if !tc.error {
				s.Require().NoError(err)
				s.Require().Equal(tc.expected.String(), tc.respType.String())
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestAllCoinsGrpcHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			name: "Get coins all on grpc",
			url:  fmt.Sprintf("%s/coin/v1/coins", baseURL),
			headers: map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			respType: &types.QueryCoinsResponse{},
			expected: &types.QueryCoinsResponse{
				Coins: coins,
				Pagination: &query.PageResponse{
					Total: 3,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType))
			respCoins, ok := tc.respType.(*types.QueryCoinsResponse)
			s.Require().True(ok)
			s.Require().True(testcoin.CoinsEqual(coins, respCoins.Coins))
		})
	}
}

func (s *IntegrationTestSuite) TestAllChecksGrpcHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			name: "Get all checks on grpc",
			url:  fmt.Sprintf("%s/coin/v1/checks", baseURL),
			headers: map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			respType: &types.QueryChecksResponse{},
			expected: &types.QueryChecksResponse{
				Checks: checks,
				Pagination: &query.PageResponse{
					Total: 2,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType))
			respChecks, ok := tc.respType.(*types.QueryChecksResponse)
			s.Require().True(ok)

			s.Require().True(testcoin.ChecksEqual(checks, respChecks.Checks))
		})
	}
}
