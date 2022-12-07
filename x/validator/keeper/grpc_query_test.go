package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/testvalidator"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryValidators() {
	queryClient, vals := suite.queryClient, suite.vals
	var req *types.QueryValidatorsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		numVals  int
		hasNext  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryValidatorsRequest{}
			},
			true,

			len(vals) + 1, // +1 validator from genesis state
			false,
		},
		{
			"empty status returns all the validators",
			func() {
				req = &types.QueryValidatorsRequest{Status: ""}
			},
			true,
			len(vals) + 1, // +1 validator from genesis state
			false,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryValidatorsRequest{Status: "test"}
			},
			false,
			0,
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryValidatorsRequest{
					Status:     types.BondStatus_Bonded.String(),
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			1,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			valsResp, err := queryClient.Validators(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(valsResp)
				suite.Equal(tc.numVals, len(valsResp.Validators))
				suite.Equal(uint64(tc.numVals), valsResp.Pagination.Total)

				if tc.hasNext {
					suite.NotNil(valsResp.Pagination.NextKey)
				} else {
					suite.Nil(valsResp.Pagination.NextKey)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryValidator() {
	dsc, ctx, queryClient, vals := suite.dsc, suite.ctx, suite.queryClient, suite.vals
	validator, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0].GetOperator())
	suite.True(found)
	var req *types.QueryValidatorRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryValidatorRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryValidatorRequest{Validator: vals[0].OperatorAddress}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Validator(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.True(validator.Equal(&res.Validator))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorValidators() {
	dsc, ctx, queryClient, addrs := suite.dsc, suite.ctx, suite.queryClient, suite.addrs
	params := dsc.ValidatorKeeper.GetParams(ctx)
	delValidators := dsc.ValidatorKeeper.GetDelegatorValidators(ctx, addrs[0], params.MaxValidators)
	var req *types.QueryDelegatorValidatorsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorValidatorsRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorValidatorsRequest{
					Delegator:  addrs[0].String(),
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorValidators(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(1, len(res.Validators))
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(len(delValidators)), res.Pagination.Total)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorValidator() {
	queryClient, addrs, vals := suite.queryClient, suite.addrs, suite.vals
	addr := addrs[1]
	addrVal, addrVal1 := vals[0].OperatorAddress, vals[1].OperatorAddress
	var req *types.QueryDelegatorValidatorRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorValidatorRequest{}
			},
			false,
		},
		{
			"invalid delegator, validator pair",
			func() {
				req = &types.QueryDelegatorValidatorRequest{
					Delegator: addr.String(),
					Validator: addrVal,
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorValidatorRequest{
					Delegator: addr.String(),
					Validator: addrVal1,
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorValidator(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(addrVal1, res.Validator.OperatorAddress)
			} else {
				suite.Error(err, "resp=%#v", res)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegation() {
	dsc, ctx, queryClient, addrs, vals := suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc, addrAcc1 := addrs[0], addrs[1]
	addrVal := vals[0].OperatorAddress
	valAddr, err := sdk.ValAddressFromBech32(addrVal)
	suite.NoError(err)
	delegation, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc, valAddr, cmdcfg.BaseDenom)
	suite.True(found)
	var req *types.QueryDelegationsRequest

	testCases := []struct {
		msg         string
		malleate    func()
		emptyResult bool
		expPass     bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegationsRequest{}
			},
			true,
			false,
		},
		{
			"invalid validator, delegator pair",
			func() {
				req = &types.QueryDelegationsRequest{
					Delegator: addrAcc1.String(),
					Validator: addrVal,
				}
			},
			true,
			true,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegationsRequest{Delegator: addrAcc.String(), Validator: addrVal}
			},
			false,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			resp, err := queryClient.Delegations(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				if tc.emptyResult {
					suite.Len(resp.Delegations, 0)
				} else {
					suite.Len(resp.Delegations, 1)
					res := resp.Delegations[0]
					suite.Equal(delegation.Validator, res.Validator)
					suite.Equal(delegation.Delegator, res.Delegator)
					suite.Equal(delegation.Stake, res.Stake)
				}
			} else {
				suite.Error(err)
				suite.Nil(resp)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorDelegations() {
	dsc, ctx, queryClient, addrs, vals := suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc := addrs[0]
	addrVal1 := vals[0].OperatorAddress
	valAddr, err := sdk.ValAddressFromBech32(addrVal1)
	suite.NoError(err)
	delegation, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc, valAddr, cmdcfg.BaseDenom)
	suite.True(found)
	var req *types.QueryDelegatorDelegationsRequest

	testCases := []struct {
		msg       string
		malleate  func()
		onSuccess func(suite *KeeperTestSuite, response *types.QueryDelegatorDelegationsResponse)
		expErr    bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorDelegationsRequest{}
			},
			func(suite *KeeperTestSuite, response *types.QueryDelegatorDelegationsResponse) {},
			true,
		},
		{
			"valid request with no delegations",
			func() {
				req = &types.QueryDelegatorDelegationsRequest{Delegator: addrs[4].String()}
			},
			func(suite *KeeperTestSuite, response *types.QueryDelegatorDelegationsResponse) {
				suite.Equal(uint64(0), response.Pagination.Total)
				suite.Len(response.Delegations, 0)
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorDelegationsRequest{
					Delegator:  addrAcc.String(),
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			func(suite *KeeperTestSuite, response *types.QueryDelegatorDelegationsResponse) {
				suite.Equal(uint64(2), response.Pagination.Total)
				suite.Len(response.Delegations, 1)
				del := response.Delegations[0]
				suite.Equal(delegation.Delegator, del.Delegator)
				suite.Equal(delegation.Validator, del.Validator)
				suite.True(delegation.Stake.Equal(&del.Stake))
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorDelegations(gocontext.Background(), req)
			if tc.expErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
				tc.onSuccess(suite, res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryValidatorDelegations() {
	dsc, ctx, queryClient, addrs, vals := suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc := addrs[0]
	addrVal1 := vals[1].OperatorAddress
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	addrVal2 := valAddrs[4]
	valAddr, err := sdk.ValAddressFromBech32(addrVal1)
	suite.NoError(err)
	delegation, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc, valAddr, cmdcfg.BaseDenom)
	suite.True(found)

	var req *types.QueryValidatorDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryValidatorDelegationsRequest{}
			},
			false,
			true,
		},
		{
			"invalid validator delegator pair",
			func() {
				req = &types.QueryValidatorDelegationsRequest{Validator: addrVal2.String()}
			},
			false,
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryValidatorDelegationsRequest{
					Validator:  addrVal1,
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.ValidatorDelegations(gocontext.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.Delegations, 1)
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(2), res.Pagination.Total)
				suite.Equal(addrVal1, res.Delegations[0].Validator)
				suite.True(delegation.Stake.Equal(&res.Delegations[0].Stake))
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.Delegations)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryUnbondingDelegation() {
	dsc, ctx, queryClient, addrs, vals := suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc2 := addrs[1]
	addrVal2 := vals[1].OperatorAddress
	valAddr, err1 := sdk.ValAddressFromBech32(addrVal2)
	suite.NoError(err1)

	delegation, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc2, valAddr, cmdcfg.BaseDenom)
	suite.True(found)
	ubdStake := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10))))
	remainStake, err := dsc.ValidatorKeeper.CalculateRemainStake(ctx, delegation.Stake, ubdStake)
	suite.NoError(err)
	_, err = dsc.ValidatorKeeper.Undelegate(ctx, addrAcc2, valAddr, ubdStake, remainStake)
	suite.NoError(err)

	unbond, found := dsc.ValidatorKeeper.GetUndelegation(ctx, addrAcc2, valAddr)
	suite.True(found)
	var req *types.QueryUndelegationRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryUndelegationRequest{}
			},
			false,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryUndelegationRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryUndelegationRequest{
					Delegator: addrAcc2.String(), Validator: addrVal2,
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Undelegation(gocontext.Background(), req)
			if tc.expPass {
				suite.NotNil(res)
				suite.Equal(unbond, res.Undelegation)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorUndelegations() {
	dsc, ctx, queryClient, addrs, vals := suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc, addrAcc1 := addrs[0], addrs[1]
	addrVal, addrVal2 := vals[0].OperatorAddress, vals[1].OperatorAddress
	valAddr1, err1 := sdk.ValAddressFromBech32(addrVal)
	suite.NoError(err1)
	valAddr2, err1 := sdk.ValAddressFromBech32(addrVal2)
	suite.NoError(err1)

	// first undelegation
	delegation, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc, valAddr1, cmdcfg.BaseDenom)
	suite.True(found)
	ubdStake1 := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10))))
	remainStake1, err := dsc.ValidatorKeeper.CalculateRemainStake(ctx, delegation.Stake, ubdStake1)
	suite.NoError(err)
	_, err = dsc.ValidatorKeeper.Undelegate(ctx, addrAcc, valAddr1, ubdStake1, remainStake1)
	suite.NoError(err)

	// second undelegation
	delegation, found = dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc, valAddr2, cmdcfg.BaseDenom)
	suite.True(found)
	ubdStake2 := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10))))
	remainStake2, err := dsc.ValidatorKeeper.CalculateRemainStake(ctx, delegation.Stake, ubdStake2)
	suite.NoError(err)
	_, err = dsc.ValidatorKeeper.Undelegate(ctx, addrAcc, valAddr2, ubdStake2, remainStake2)
	suite.NoError(err)

	unbond, found := dsc.ValidatorKeeper.GetUndelegation(ctx, addrAcc, valAddr1)
	suite.True(found)
	var req *types.QueryDelegatorUndelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorUndelegationsRequest{}
			},
			false,
			true,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryDelegatorUndelegationsRequest{Delegator: addrAcc1.String()}
			},
			false,
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorUndelegationsRequest{
					Delegator:  addrAcc.String(),
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorUndelegations(gocontext.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(2), res.Pagination.Total)
				suite.Len(res.Undelegations, 1)
				suite.Equal(unbond, res.Undelegations[0])
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.Undelegations)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryPoolParameters() {
	dsc, ctx, queryClient := suite.dsc, suite.ctx, suite.queryClient

	// Query pool
	res, err := queryClient.Pool(gocontext.Background(), &types.QueryPoolRequest{})
	suite.NoError(err)
	bondedPool := dsc.ValidatorKeeper.GetBondedPool(ctx)
	notBondedPool := dsc.ValidatorKeeper.GetNotBondedPool(ctx)

	suite.True(dsc.BankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress()).IsEqual(res.Pool.NotBonded))
	suite.True(dsc.BankKeeper.GetAllBalances(ctx, bondedPool.GetAddress()).IsEqual(res.Pool.Bonded))

	// Query Params
	resp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.NoError(err)
	suite.Equal(dsc.ValidatorKeeper.GetParams(ctx), resp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryHistoricalInfo() {
	dsc, ctx, queryClient := suite.dsc, suite.ctx, suite.queryClient

	hi, found := dsc.ValidatorKeeper.GetHistoricalInfoDecimal(ctx, 5)
	suite.True(found)

	var req *types.QueryHistoricalInfoRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryHistoricalInfoRequest{}
			},
			false,
		},
		{
			"invalid request with negative height",
			func() {
				req = &types.QueryHistoricalInfoRequest{Height: -1}
			},
			false,
		},
		{
			"valid request with old height",
			func() {
				req = &types.QueryHistoricalInfoRequest{Height: 4}
			},
			false,
		},
		{
			"valid request with current height",
			func() {
				req = &types.QueryHistoricalInfoRequest{Height: 5}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.HistoricalInfo(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(res)
				suite.True(hi.Equal(res.Hist))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryRedelegations() {
	dsc, ctx, queryClient, addrs, vals := suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals

	addrAcc, addrAcc1 := addrs[0], addrs[1]
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	val1, val2, val3 := vals[0], vals[1], valAddrs[3]
	//applyValidatorSetUpdates(suite.T(), ctx, dsc.ValidatorKeeper, -1)
	val1.Status = types.BondStatus_Bonded
	dsc.ValidatorKeeper.SetValidator(ctx, val1)

	// create redelegation
	delegation, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc, val1.GetOperator(), cmdcfg.BaseDenom)
	suite.True(found)
	redStake := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10))))
	remainStake, err := dsc.ValidatorKeeper.CalculateRemainStake(ctx, delegation.Stake, redStake)
	suite.NoError(err)
	_, err = dsc.ValidatorKeeper.BeginRedelegation(ctx, addrAcc, val1.GetOperator(), val2.GetOperator(), redStake, remainStake)
	suite.NoError(err)

	redel, found := dsc.ValidatorKeeper.GetRedelegation(ctx, addrAcc, val1.GetOperator(), val2.GetOperator())
	suite.True(found)

	var req *types.QueryRedelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"request redelegations for non existent addr",
			func() {
				req = &types.QueryRedelegationsRequest{Delegator: addrAcc1.String()}
			},
			false,
			true,
		},
		{
			"request redelegations with non existent pairs",
			func() {
				req = &types.QueryRedelegationsRequest{
					Delegator: addrAcc.String(), Validator: val3.String(),
				}
			},
			false,
			false,
		},
		{
			"request redelegations with delegator, validator",
			func() {
				req = &types.QueryRedelegationsRequest{
					Delegator: addrAcc.String(), Validator: val1.OperatorAddress,
				}
			},
			true,
			false,
		},
		{
			"query redelegations with sourceValAddr only",
			func() {
				req = &types.QueryRedelegationsRequest{
					Validator: val1.GetOperator().String(),
				}
			},
			false,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Redelegations(gocontext.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.Redelegations, 1)
				suite.Len(res.Redelegations[0].Entries, 1)
				suite.Equal(redel.Delegator, res.Redelegations[0].Delegator)
				suite.Equal(redel.ValidatorSrc, res.Redelegations[0].ValidatorSrc)
				suite.Equal(redel.ValidatorDst, res.Redelegations[0].ValidatorDst)
				suite.Len(redel.Entries, len(res.Redelegations[0].Entries))
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.Redelegations)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorRedelegations() {
	dsc, ctx, queryClient, addrs, vals := suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals

	addrAcc, addrAcc1 := addrs[0], addrs[1]
	val1, val2 := vals[0], vals[1]
	//applyValidatorSetUpdates(suite.T(), ctx, dsc.ValidatorKeeper, -1)
	val1.Status = types.BondStatus_Bonded
	dsc.ValidatorKeeper.SetValidator(ctx, val1)

	// create redelegation
	delegation, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc, val1.GetOperator(), cmdcfg.BaseDenom)
	suite.True(found)
	redStake := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10))))
	remainStake, err := dsc.ValidatorKeeper.CalculateRemainStake(ctx, delegation.Stake, redStake)
	suite.NoError(err)
	_, err = dsc.ValidatorKeeper.BeginRedelegation(ctx, addrAcc, val1.GetOperator(), val2.GetOperator(), redStake, remainStake)
	suite.NoError(err)

	redel, found := dsc.ValidatorKeeper.GetRedelegation(ctx, addrAcc, val1.GetOperator(), val2.GetOperator())
	suite.True(found)

	var req *types.QueryDelegatorRedelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expCount int
		expErr   bool
	}{
		{
			"request redelegations for non existent redelegations for addr",
			func() {
				req = &types.QueryDelegatorRedelegationsRequest{Delegator: addrAcc1.String()}
			},
			0,
			false,
		},
		{
			"request redelegations with delegator",
			func() {
				req = &types.QueryDelegatorRedelegationsRequest{
					Delegator: addrAcc.String(),
				}
			},
			1,
			false,
		},
		{
			"request redelegations with invalid address delegator",
			func() {
				req = &types.QueryDelegatorRedelegationsRequest{
					Delegator: "d01asasadasd",
				}
			},
			0,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorRedelegations(gocontext.Background(), req)
			if tc.expCount > 0 && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.Redelegations, tc.expCount)
				suite.Len(res.Redelegations[0].Entries, 1)
				suite.Equal(redel.Delegator, res.Redelegations[0].Delegator)
				suite.Equal(redel.ValidatorSrc, res.Redelegations[0].ValidatorSrc)
				suite.Equal(redel.ValidatorDst, res.Redelegations[0].ValidatorDst)
				suite.Len(redel.Entries, len(res.Redelegations[0].Entries))
			} else if tc.expCount == 0 && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.Redelegations, 0)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryValidatorRedelegations() {
	dsc, ctx, queryClient, addrs, vals := suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals

	addrAcc := addrs[0]
	val1, val2 := vals[0], vals[1]
	//applyValidatorSetUpdates(suite.T(), ctx, dsc.ValidatorKeeper, -1)
	val1.Status = types.BondStatus_Bonded
	dsc.ValidatorKeeper.SetValidator(ctx, val1)

	// create redelegation
	delegation, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc, val1.GetOperator(), cmdcfg.BaseDenom)
	suite.True(found)
	redStake := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10))))
	remainStake, err := dsc.ValidatorKeeper.CalculateRemainStake(ctx, delegation.Stake, redStake)
	suite.NoError(err)
	_, err = dsc.ValidatorKeeper.BeginRedelegation(ctx, addrAcc, val1.GetOperator(), val2.GetOperator(), redStake, remainStake)
	suite.NoError(err)

	redel, found := dsc.ValidatorKeeper.GetRedelegation(ctx, addrAcc, val1.GetOperator(), val2.GetOperator())
	suite.True(found)

	var req *types.QueryValidatorRedelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expCount int
		expErr   bool
	}{
		{
			"request redelegations for non existent redelegations for addr",
			func() {
				req = &types.QueryValidatorRedelegationsRequest{Validator: val2.OperatorAddress}
			},
			0,
			false,
		},
		{
			"request redelegations with delegator",
			func() {
				req = &types.QueryValidatorRedelegationsRequest{
					Validator: val1.OperatorAddress,
				}
			},
			1,
			false,
		},
		{
			"request redelegations with invalid address delegator",
			func() {
				req = &types.QueryValidatorRedelegationsRequest{
					Validator: "d01asasadasd",
				}
			},
			0,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.ValidatorRedelegations(gocontext.Background(), req)
			if tc.expCount > 0 && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.Redelegations, tc.expCount)
				suite.Len(res.Redelegations[0].Entries, 1)
				suite.Equal(redel.Delegator, res.Redelegations[0].Delegator)
				suite.Equal(redel.ValidatorSrc, res.Redelegations[0].ValidatorSrc)
				suite.Equal(redel.ValidatorDst, res.Redelegations[0].ValidatorDst)
				suite.Len(redel.Entries, len(res.Redelegations[0].Entries))
			} else if tc.expCount == 0 && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.Redelegations, 0)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryValidatorUndelegations() {
	dsc, ctx, queryClient, addrs, vals := suite.dsc, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc1, _ := addrs[0], addrs[1]
	val1 := vals[0]
	valAddr1, err1 := sdk.ValAddressFromBech32(val1.OperatorAddress)
	suite.NoError(err1)

	// first undelegation
	delegation, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrAcc1, valAddr1, cmdcfg.BaseDenom)
	suite.True(found)
	ubdStake1 := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10))))
	remainStake1, err := dsc.ValidatorKeeper.CalculateRemainStake(ctx, delegation.Stake, ubdStake1)
	suite.NoError(err)
	_, err = dsc.ValidatorKeeper.Undelegate(ctx, addrAcc1, valAddr1, ubdStake1, remainStake1)
	suite.NoError(err)

	var req *types.QueryValidatorUndelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryValidatorUndelegationsRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryValidatorUndelegationsRequest{
					Validator:  val1.GetOperator().String(),
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.ValidatorUndelegations(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(uint64(1), res.Pagination.Total)
				suite.Equal(1, len(res.Undelegations))
				suite.Equal(res.Undelegations[0].Validator, val1.OperatorAddress)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func createValidators(t *testing.T, ctx sdk.Context, dsc *app.DSC, powers []int64) ([]sdk.AccAddress, []sdk.ValAddress, []types.Validator) {
	stake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	coins := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(300)))
	addrs := app.AddTestAddrsIncremental(dsc, ctx, 5, sdk.NewCoins(coins))
	valAddrs := app.ConvertAddrsToValAddrs(addrs)
	pks := simapp.CreateTestPubKeys(5)
	dsc.ValidatorKeeper = keeper.NewKeeper(
		dsc.AppCodec(),
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.AccountKeeper,
		dsc.BankKeeper,
		&dsc.NFTKeeper,
		&dsc.CoinKeeper,
		&dsc.MultisigKeeper,
	)

	val1 := testvalidator.NewValidator(t, valAddrs[0], pks[0])
	val2 := testvalidator.NewValidator(t, valAddrs[1], pks[1])
	vals := []types.Validator{val1, val2}

	dsc.ValidatorKeeper.SetValidator(ctx, val1)
	dsc.ValidatorKeeper.SetValidator(ctx, val2)
	dsc.ValidatorKeeper.SetValidatorByConsAddr(ctx, val1)
	dsc.ValidatorKeeper.SetValidatorByConsAddr(ctx, val2)
	dsc.ValidatorKeeper.SetNewValidatorByPowerIndex(ctx, val1)
	dsc.ValidatorKeeper.SetNewValidatorByPowerIndex(ctx, val2)

	err := dsc.ValidatorKeeper.Delegate(ctx, addrs[0], val1, types.NewStakeCoin(stake))
	require.NoError(t, err)
	err = dsc.ValidatorKeeper.Delegate(ctx, addrs[1], val2, types.NewStakeCoin(stake))
	require.NoError(t, err)
	err = dsc.ValidatorKeeper.Delegate(ctx, addrs[0], val2, types.NewStakeCoin(stake))
	require.NoError(t, err)
	//applyValidatorSetUpdates(t, ctx, app.StakingKeeper, -1)

	return addrs, valAddrs, vals
}
