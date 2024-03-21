package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strings"

	sdkmath "cosmossdk.io/math"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

var _ types.QueryServer = Querier{}

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

// Validators queries all validators that match the given status.
func (k Querier) Validators(c context.Context, req *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// validate the provided status, return all the validators if the status is empty
	if req.Status != "" && !(req.Status == types.BondStatus_Bonded.String() || req.Status == types.BondStatus_Unbonded.String() || req.Status == types.BondStatus_Unbonding.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator status %s", req.Status)
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	valStore := prefix.NewStore(store, types.GetValidatorsKey())

	validators, pageRes, err := query.GenericFilteredPaginate(k.cdc, valStore, req.Pagination, func(key []byte, val *types.Validator) (*types.Validator, error) {
		if req.Status != "" && !strings.EqualFold(val.GetStatus().String(), req.Status) {
			return nil, nil
		}
		rs, err := k.GetValidatorRS(ctx, val.GetOperator())
		if err != nil {
			val.Rewards = sdkmath.ZeroInt()
			val.TotalRewards = sdkmath.ZeroInt()
			val.Stake = 0
		} else {
			val.Rewards = rs.Rewards
			val.TotalRewards = rs.TotalRewards
			val.Stake = rs.Stake
		}
		accAddress, _ := sdk.AccAddressFromBech32(val.GetOperator().String())

		val.DRC20Contract = common.BytesToAddress(accAddress).String()
		return val, nil
	}, func() *types.Validator {
		return &types.Validator{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	vals := types.Validators{}
	for _, val := range validators {
		vals = append(vals, *val)
	}

	return &types.QueryValidatorsResponse{Validators: vals, Pagination: pageRes}, nil
}

// Validator queries validator info for given validator address.
func (k Querier) Validator(c context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "validator %s not found", req.Validator)
	}

	return &types.QueryValidatorResponse{Validator: validator}, nil
}

// ValidatorDelegations queries delegate info for given validator.
func (k Querier) ValidatorDelegations(c context.Context, req *types.QueryValidatorDelegationsRequest) (*types.QueryValidatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	valAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)
	byValPrefix := types.GetValidatorDelegationsKey(valAddr)
	valStore := prefix.NewStore(store, byValPrefix)
	if err != nil {
		return nil, err
	}

	var delegations []types.Delegation
	pageRes, err := query.Paginate(valStore, req.Pagination, func(key []byte, _ []byte) error {
		realKey := types.GetDelegationKeyFromValIndexKey(append(byValPrefix, key...))
		del, err := types.UnmarshalDelegation(k.cdc, store.Get(realKey))
		if err != nil {
			return err
		}
		delegations = append(delegations, del)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidatorDelegationsResponse{
		Delegations: delegations, Pagination: pageRes,
	}, nil
}

// ValidatorRedelegations queries redelegations of a validator.
func (k Querier) ValidatorRedelegations(c context.Context, req *types.QueryValidatorRedelegationsRequest) (*types.QueryValidatorRedelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	valAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)
	fromValSrcKey := types.GetREDsFromValSrcIndexKey(valAddr)
	valStore := prefix.NewStore(store, fromValSrcKey)

	var redelegations []types.Redelegation

	pageRes, err := query.Paginate(valStore, req.Pagination, func(key []byte, _ []byte) error {
		realKey := types.GetREDKeyFromValSrcIndexKey(append(fromValSrcKey, key...))
		value := store.Get(realKey)
		red, err := types.UnmarshalRED(k.cdc, value)
		if err != nil {
			return err
		}
		redelegations = append(redelegations, red)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidatorRedelegationsResponse{
		Redelegations: redelegations,
		Pagination:    pageRes,
	}, nil
}

// ValidatorUndelegations queries undelegations of a validator.
func (k Querier) ValidatorUndelegations(c context.Context, req *types.QueryValidatorUndelegationsRequest) (*types.QueryValidatorUndelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	valAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)
	byValKey := types.GetUBDsByValIndexKey(valAddr)
	valStore := prefix.NewStore(store, byValKey)

	ubds := []types.Undelegation{}

	pageRes, err := query.Paginate(valStore, req.Pagination,
		func(key []byte, _ []byte) error {
			realKey := types.GetUBDKeyFromValIndexKey(append(byValKey, key...))
			value := store.Get(realKey)
			ubd, err := types.UnmarshalUBD(k.cdc, value)
			if err != nil {
				return err
			}
			ubds = append(ubds, ubd)
			return nil
		})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidatorUndelegationsResponse{
		Undelegations: ubds,
		Pagination:    pageRes,
	}, nil
}

// Delegations queries delegations info for given validator delegator pair.
func (k Querier) Delegations(c context.Context, req *types.QueryDelegationsRequest) (*types.QueryDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}

	var result []types.Delegation
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetValidatorDelegatorDelegationsKey(valAddr, delAddr))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		realKey := types.GetDelegationKeyFromValIndexKey(iterator.Key())
		del, err := types.UnmarshalDelegation(k.cdc, store.Get(realKey))
		if err != nil {
			return nil, err
		}
		result = append(result, del)
	}

	return &types.QueryDelegationsResponse{Delegations: result}, nil
}

// Redelegations queries redelegations info for given validator delegator pair.
func (k Querier) Redelegations(c context.Context, req *types.QueryRedelegationsRequest) (*types.QueryRedelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, err
	}

	_, err = sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}

	var result []types.Redelegation
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetREDsKey(delAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		red, err := types.UnmarshalRED(k.cdc, iterator.Value())
		if err != nil {
			return nil, err
		}
		if red.ValidatorSrc != req.Validator {
			continue
		}
		result = append(result, red)
	}

	return &types.QueryRedelegationsResponse{Redelegations: result}, nil
}

// Undelegation queries undelegations info for given validator delegator pair.
func (k Querier) Undelegation(c context.Context, req *types.QueryUndelegationRequest) (*types.QueryUndelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}

	result, found := k.Keeper.GetUndelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, errors.UBDNotFound
	}
	return &types.QueryUndelegationResponse{Undelegation: result}, nil
}

// DelegatorDelegations queries all delegations of a give delegator address
func (k Querier) DelegatorDelegations(c context.Context, req *types.QueryDelegatorDelegationsRequest) (*types.QueryDelegatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var delegations types.Delegations
	ctx := sdk.UnwrapSDKContext(c)

	delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.GetDelegatorDelegationsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorDelegationsResponse{
		Delegations: delegations,
		Pagination:  pageRes,
	}, nil
}

// DelegatorRedelegations queries all redelegations of a given delegator address.
func (k Querier) DelegatorRedelegations(c context.Context, req *types.QueryDelegatorRedelegationsRequest) (*types.QueryDelegatorRedelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var redelegations types.Redelegations
	ctx := sdk.UnwrapSDKContext(c)

	delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.GetREDsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		red, err := types.UnmarshalRED(k.cdc, value)
		if err != nil {
			return err
		}
		redelegations = append(redelegations, red)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorRedelegationsResponse{
		Redelegations: redelegations,
		Pagination:    pageRes,
	}, nil
}

// DelegatorUndelegations queries all undelegations of a given delegator address.
func (k Querier) DelegatorUndelegations(c context.Context, req *types.QueryDelegatorUndelegationsRequest) (*types.QueryDelegatorUndelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var undelegations types.Undelegations
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, err
	}

	unbStore := prefix.NewStore(store, types.GetUBDsKey(delAddr))
	pageRes, err := query.Paginate(unbStore, req.Pagination, func(key []byte, value []byte) error {
		unbond, err := types.UnmarshalUBD(k.cdc, value)
		if err != nil {
			return err
		}
		undelegations = append(undelegations, unbond)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorUndelegationsResponse{
		Undelegations: undelegations, Pagination: pageRes,
	}, nil
}

// DelegatorValidators queries all validators info for given delegator address.
func (k Querier) DelegatorValidators(c context.Context, req *types.QueryDelegatorValidatorsRequest) (*types.QueryDelegatorValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var validators types.Validators
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, err
	}

	delStore := prefix.NewStore(store, types.GetDelegatorDelegationsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		validator, found := k.GetValidator(ctx, delegation.GetValidator())
		if !found {
			return types.ErrNoValidatorFound
		}

		validators = append(validators, validator)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorValidatorsResponse{
		Validators: validators,
		Pagination: pageRes,
	}, nil
}

// DelegatorValidator queries validator info for given delegator validator pair.
func (k Querier) DelegatorValidator(c context.Context, req *types.QueryDelegatorValidatorRequest) (*types.QueryDelegatorValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}

	validator, err := k.GetDelegatorValidator(ctx, delAddr, valAddr)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorValidatorResponse{
		Validator: validator,
	}, nil
}

// HistoricalInfo queries the historical info for given height.
func (k Querier) HistoricalInfo(c context.Context, req *types.QueryHistoricalInfoRequest) (*types.QueryHistoricalInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Height < 0 {
		return nil, status.Error(codes.InvalidArgument, "height cannot be negative")
	}
	ctx := sdk.UnwrapSDKContext(c)
	hi, found := k.GetHistoricalInfoDecimal(ctx, req.Height)
	if !found {
		return nil, status.Errorf(codes.NotFound, "historical info for height %d not found", req.Height)
	}

	return &types.QueryHistoricalInfoResponse{Hist: &hi}, nil
}

// Pool queries the pool info.
func (k Querier) Pool(c context.Context, _ *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bondedPool := k.GetBondedPool(ctx)
	notBondedPool := k.GetNotBondedPool(ctx)

	pool := types.NewPool(
		k.bankKeeper.GetAllBalances(ctx, bondedPool.GetAddress()),
		k.bankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress()),
	)

	return &types.QueryPoolResponse{Pool: pool}, nil
}

// CustomCoinPrice queries the total amount bonded custom coins.
func (k Querier) CustomCoinPrice(c context.Context, req *types.QueryCustomCoinPriceRequest) (*types.QueryCustomCoinPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	baseDenom := k.BaseDenom(ctx)
	if req.Denom == baseDenom {
		return nil, fmt.Errorf("is not custom coin")
	}

	customCoinStaked := k.GetCustomCoinStaked(ctx, req.Denom)
	if customCoinStaked.Equal(sdk.ZeroInt()) {
		customCoin, err := k.coinKeeper.GetCoin(ctx, req.Denom)
		if err != nil {
			panic(err)
		}
		amountInBaseCoin := formulas.CalculateSaleReturn(customCoin.Volume, customCoin.Reserve, uint(customCoin.CRR), customCoinStaked)

		return &types.QueryCustomCoinPriceResponse{
			Price: sdk.NewDecFromInt(amountInBaseCoin),
		}, nil
	}

	customCoinPrice := k.calculateCustomCoinPrice(ctx, req.Denom, customCoinStaked)

	baseAmount := sdk.NewDecFromInt(sdk.NewInt(1)).Mul(customCoinPrice)

	return &types.QueryCustomCoinPriceResponse{
		Price: baseAmount,
	}, nil
}

// TotalCustomCoin queries the total amount bonded custom coins.
func (k Querier) TotalCustomCoin(c context.Context, req *types.QueryTotalCustomCoinRequest) (*types.QueryTotalCustomCoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if req.Denom == k.BaseDenom(ctx) {
		return nil, fmt.Errorf("is not custom coin")
	}

	customCoinStaked := k.GetCustomCoinStaked(ctx, req.Denom)

	return &types.QueryTotalCustomCoinResponse{
		TotalAmount: customCoinStaked,
	}, nil
}

// Params queries the module params.
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}
