package keeper_test

import (
	gocontext "context"
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (s *KeeperTestSuite) TestGRPCQueryCoinRecords() {
	ctx, k, queryClient := s.ctx, s.legacyKeeper, s.queryClient
	require := s.Require()

	hits := make(map[string]types.Record)
	records := []types.Record{
		defaultRecord,
		{
			LegacyAddress: "dx1m3eg7v6pu0dga2knj9zm4683dk9c8800j9nfw0",
			Coins:         sdk.Coins{},
			Wallets:       []string{},
			NFTs:          []string{},
		},
	}

	for _, v := range records {
		k.SetLegacyRecord(ctx, v)
	}

	var req *types.QueryRecordsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryRecordsRequest{
					Pagination: &query.PageRequest{
						Offset: 0,
						Limit:  100,
					},
				}
				hits = map[string]types.Record{
					records[0].LegacyAddress: records[0],
					records[1].LegacyAddress: records[1],
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Records(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				for _, resRecord := range res.GetRecords() {
					if record, ok := hits[resRecord.LegacyAddress]; ok {
						require.Equal(record.String(), resRecord.String()) // TODO replace to equal after regenerate proto
						delete(hits, resRecord.LegacyAddress)
					} else {
						s.T().Fatal("record does not set, but it was included in the resp")
					}
				}
				require.Equal(0, len(hits), "not all records were returned")
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryCoinRecord() {
	ctx, k, queryClient := s.ctx, s.legacyKeeper, s.queryClient
	require := s.Require()

	k.SetLegacyRecord(ctx, defaultRecord)

	var req *types.QueryRecordRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryRecordRequest{LegacyAddress: defaultRecord.LegacyAddress}
			},
			true,
		},
		{
			"not exists record",
			func() {
				req = &types.QueryRecordRequest{LegacyAddress: invalidOldAddress}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Record(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				require.Equal(defaultRecord.String(), res.Record.String()) // TODO replace to equal after regenerate proto
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryCheck() {
	ctx, k, queryClient := s.ctx, s.legacyKeeper, s.queryClient
	require := s.Require()

	k.SetLegacyRecord(ctx, defaultRecord)

	var req *types.QueryCheckRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryCheckRequest{Pubkey: publicKey}
			},
			true,
		},
		{
			"not exist record",
			func() {
				req = &types.QueryCheckRequest{Pubkey: invalidAccPk.Bytes()}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Check(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				require.Equal(defaultRecord.String(), res.Record.String()) // TODO replace to equal after regenerate proto
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}
