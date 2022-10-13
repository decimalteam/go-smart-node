package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

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
	valStore := prefix.NewStore(store, types.GetValidatorDelegationsKey(valAddr))
	if err != nil {
		return nil, err
	}
	k.GetValidatorDelegations(ctx, valAddr)
	delegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, valStore, req.Pagination, func(key []byte, delegation *types.Delegation) (*types.Delegation, error) {
		return delegation, nil
	}, func() *types.Delegation {
		return &types.Delegation{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dels := types.Delegations{}
	for _, d := range delegations {
		dels = append(dels, *d)
	}

	return &types.QueryValidatorDelegationsResponse{
		Delegations: dels, Pagination: pageRes,
	}, nil
}

// ValidatorRedelegations queries redelegations of a validator.
func (k Querier) ValidatorRedelegations(c context.Context, req *types.QueryValidatorRedelegationsRequest) (*types.QueryValidatorRedelegationsResponse, error) {
	//if req == nil {
	//	return nil, status.Error(codes.InvalidArgument, "empty request")
	//}
	//if req.Validator == "" {
	//	return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	//}
	//ctx := sdk.UnwrapSDKContext(c)
	//
	//store := ctx.KVStore(k.storeKey)
	//redels, pagination, err := queryRedelegationsFromSrcValidator(store, k, req.Validator, req.Pagination)
	//if err != nil {
	//	return nil, err
	//}

	return &types.QueryValidatorRedelegationsResponse{
		//Redelegations: redels,
		//Pagination:    pagination,
	}, nil
}

// ValidatorUndelegations queries undelegations of a validator.
func (k Querier) ValidatorUndelegations(c context.Context, req *types.QueryValidatorUndelegationsRequest) (*types.QueryValidatorUndelegationsResponse, error) {
	//if req == nil {
	//	return nil, status.Error(codes.InvalidArgument, "empty request")
	//}
	//
	//if req.ValidatorAddr == "" {
	//	return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	//}
	//var ubds types.Undelegations
	//ctx := sdk.UnwrapSDKContext(c)
	//
	//store := ctx.KVStore(k.storeKey)
	//
	//valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//srcValPrefix := types.GetUBDsByValIndexKey(valAddr)
	//ubdStore := prefix.NewStore(store, srcValPrefix)
	//pageRes, err := query.Paginate(ubdStore, req.Pagination, func(key []byte, value []byte) error {
	//	storeKey := types.GetUBDKeyFromValIndexKey(append(srcValPrefix, key...))
	//	storeValue := store.Get(storeKey)
	//
	//	ubd, err := types.UnmarshalUBD(k.cdc, storeValue)
	//	if err != nil {
	//		return err
	//	}
	//	ubds = append(ubds, ubd)
	//	return nil
	//})
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//
	//return &types.QueryValidatorUndelegationsResponse{
	//	UnbondingResponses: ubds,
	//	Pagination:         pageRes,
	//}, nil
	return &types.QueryValidatorUndelegationsResponse{}, nil
}

// Delegations queries delegations info for given validator delegator pair.
func (k Querier) Delegations(c context.Context, req *types.QueryDelegationsRequest) (*types.QueryDelegationsResponse, error) {
	//if req == nil {
	//	return nil, status.Error(codes.InvalidArgument, "empty request")
	//}
	//
	//if req.DelegatorAddr == "" {
	//	return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	//}
	//if req.ValidatorAddr == "" {
	//	return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	//}
	//
	//ctx := sdk.UnwrapSDKContext(c)
	//delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	//if !found {
	//	return nil, status.Errorf(
	//		codes.NotFound,
	//		"delegation with delegator %s not found for validator %s",
	//		req.DelegatorAddr, req.ValidatorAddr)
	//}
	//
	//delResponse, err := DelegationToDelegationResponse(ctx, k.Keeper, delegation)
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//
	//return &types.QueryDelegationsResponse{DelegationResponse: &delResponse}, nil
	return &types.QueryDelegationsResponse{}, nil
}

// Redelegations queries redelegations info for given validator delegator pair.
func (k Querier) Redelegations(c context.Context, req *types.QueryRedelegationsRequest) (*types.QueryRedelegationsResponse, error) {
	//if req == nil {
	//	return nil, status.Error(codes.InvalidArgument, "empty request")
	//}
	//
	//var redels types.Redelegations
	//var pageRes *query.PageResponse
	//var err error
	//
	//ctx := sdk.UnwrapSDKContext(c)
	//store := ctx.KVStore(k.storeKey)
	//switch {
	//case req.DelegatorAddr != "" && req.SrcValidatorAddr != "" && req.DstValidatorAddr != "":
	//	redels, err = queryRedelegation(ctx, k, req)
	//case req.DelegatorAddr == "" && req.SrcValidatorAddr != "" && req.DstValidatorAddr == "":
	//	redels, pageRes, err = queryRedelegationsFromSrcValidator(store, k, req)
	//default:
	//	redels, pageRes, err = queryAllRedelegations(store, k, req)
	//}
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//redelResponses, err := RedelegationsToRedelegationResponses(ctx, k.Keeper, redels)
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	return &types.QueryRedelegationsResponse{}, nil
}

// Undelegations queries undelegations info for given validator delegator pair.
func (k Querier) Undelegations(c context.Context, req *types.QueryUndelegationsRequest) (*types.QueryUndelegationsResponse, error) {
	//if req == nil {
	//	return nil, status.Errorf(codes.InvalidArgument, "empty request")
	//}
	//
	//if req.DelegatorAddr == "" {
	//	return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	//}
	//if req.ValidatorAddr == "" {
	//	return nil, status.Errorf(codes.InvalidArgument, "validator address cannot be empty")
	//}
	//
	//ctx := sdk.UnwrapSDKContext(c)
	//
	//delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//unbond, found := k.GetUndelegation(ctx, delAddr, valAddr)
	//if !found {
	//	return nil, status.Errorf(
	//		codes.NotFound,
	//		"unbonding delegation with delegator %s not found for validator %s",
	//		req.DelegatorAddr, req.ValidatorAddr)
	//}

	return &types.QueryUndelegationsResponse{}, nil
}

// DelegatorDelegations queries all delegations of a give delegator address
func (k Querier) DelegatorDelegations(c context.Context, req *types.QueryDelegatorDelegationsRequest) (*types.QueryDelegatorDelegationsResponse, error) {
	//if req == nil {
	//	return nil, status.Errorf(codes.InvalidArgument, "empty request")
	//}
	//
	//if req.DelegatorAddr == "" {
	//	return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	//}
	//var delegations types.Delegations
	//ctx := sdk.UnwrapSDKContext(c)
	//
	//delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//store := ctx.KVStore(k.storeKey)
	//delStore := prefix.NewStore(store, types.GetDelegatorDelegationsKey(delAddr))
	//pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
	//	delegation, err := types.UnmarshalDelegation(k.cdc, value)
	//	if err != nil {
	//		return err
	//	}
	//	delegations = append(delegations, delegation)
	//	return nil
	//})
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//
	//delegationResps, err := DelegationsToDelegationResponses(ctx, k.Keeper, delegations)
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	return &types.QueryDelegatorDelegationsResponse{}, nil
}

// DelegatorRedelegations queries all redelegations of a given delegator address.
func (k Querier) DelegatorRedelegations(c context.Context, req *types.QueryDelegatorRedelegationsRequest) (*types.QueryDelegatorRedelegationsResponse, error) {
	return &types.QueryDelegatorRedelegationsResponse{}, nil
}

// DelegatorUndelegations queries all undelegations of a given delegator address.
func (k Querier) DelegatorUndelegations(c context.Context, req *types.QueryDelegatorUndelegationsRequest) (*types.QueryDelegatorUndelegationsResponse, error) {
	//if req == nil {
	//	return nil, status.Error(codes.InvalidArgument, "empty request")
	//}
	//
	//if req.DelegatorAddr == "" {
	//	return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	//}
	//var undelegations types.Undelegations
	//ctx := sdk.UnwrapSDKContext(c)
	//
	//store := ctx.KVStore(k.storeKey)
	//delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//unbStore := prefix.NewStore(store, types.GetUBDsKey(delAddr))
	//pageRes, err := query.Paginate(unbStore, req.Pagination, func(key []byte, value []byte) error {
	//	unbond, err := types.UnmarshalUBD(k.cdc, value)
	//	if err != nil {
	//		return err
	//	}
	//	undelegations = append(undelegations, unbond)
	//	return nil
	//})
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	return &types.QueryDelegatorUndelegationsResponse{
		//UnbondingResponses: undelegations, Pagination: pageRes,
	}, nil
}

// DelegatorValidators queries all validators info for given delegator address.
func (k Querier) DelegatorValidators(c context.Context, req *types.QueryDelegatorValidatorsRequest) (*types.QueryDelegatorValidatorsResponse, error) {
	//if req == nil {
	//	return nil, status.Error(codes.InvalidArgument, "empty request")
	//}
	//
	//if req.DelegatorAddr == "" {
	//	return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	//}
	//var validators types.Validators
	//ctx := sdk.UnwrapSDKContext(c)
	//
	//store := ctx.KVStore(k.storeKey)
	//delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//delStore := prefix.NewStore(store, types.GetDelegatorDelegationsKey(delAddr))
	//pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
	//	delegation, err := types.UnmarshalDelegation(k.cdc, value)
	//	if err != nil {
	//		return err
	//	}
	//
	//	validator, found := k.GetValidator(ctx, delegation.GetValidatorAddr())
	//	if !found {
	//		return types.ErrNoValidatorFound
	//	}
	//
	//	validators = append(validators, validator)
	//	return nil
	//})
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	return &types.QueryDelegatorValidatorsResponse{}, nil
}

// DelegatorValidator queries validator info for given delegator validator pair.
func (k Querier) DelegatorValidator(c context.Context, req *types.QueryDelegatorValidatorRequest) (*types.QueryDelegatorValidatorResponse, error) {
	//if req == nil {
	//	return nil, status.Error(codes.InvalidArgument, "empty request")
	//}
	//
	//if req.DelegatorAddr == "" {
	//	return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	//}
	//if req.ValidatorAddr == "" {
	//	return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	//}
	//
	//ctx := sdk.UnwrapSDKContext(c)
	//delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//validator, err := k.GetDelegatorValidator(ctx, delAddr, valAddr)
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//
	return &types.QueryDelegatorValidatorResponse{}, nil
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
		k.bankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress()),
		k.bankKeeper.GetAllBalances(ctx, bondedPool.GetAddress()),
	)

	return &types.QueryPoolResponse{Pool: pool}, nil
}

// Params queries the module params.
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

//func queryRedelegation(ctx sdk.Context, k Querier, req *types.QueryRedelegationsRequest) (redels types.Redelegations, err error) {
//	//delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//srcValAddr, err := sdk.ValAddressFromBech32(req.SrcValidator)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//dstValAddr, err := sdk.ValAddressFromBech32(req.DstValidator)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//redel, found := k.GetRedelegation(ctx, delAddr, srcValAddr, dstValAddr)
//	//if !found {
//	//	return nil, status.Errorf(
//	//		codes.NotFound,
//	//		"redelegation not found for delegator address %s from validator address %s",
//	//		req.DelegatorAddr, req.SrcValidatorAddr)
//	//}
//	//redels = []types.Redelegation{redel}
//	//
//	//return redels, err
//}
//
//func queryRedelegationsFromSrcValidator(store sdk.KVStore, k Querier, srcValidator string, pagination *query.PageRequest) (redels types.Redelegations, res *query.PageResponse, err error) {
//	valAddr, err := sdk.ValAddressFromBech32(srcValidator)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	srcValPrefix := types.GetREDsFromValSrcIndexKey(valAddr)
//	redStore := prefix.NewStore(store, srcValPrefix)
//	res, err = query.Paginate(redStore, pagination, func(key []byte, value []byte) error {
//		storeKey := types.GetREDKeyFromValSrcIndexKey(append(srcValPrefix, key...))
//		storeValue := store.Get(storeKey)
//		red, err := types.UnmarshalRED(k.cdc, storeValue)
//		if err != nil {
//			return err
//		}
//		redels = append(redels, red)
//		return nil
//	})
//
//	return redels, res, err
//}
//
//func queryUndelegationsFromSrcValidator(store sdk.KVStore, k Querier, srcValidator string, pagination *query.PageRequest) (redels types.Undelegations, res *query.PageResponse, err error) {
//	valAddr, err := sdk.ValAddressFromBech32(srcValidator)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	srcValPrefix := types.GetUBDsByValIndexKey(valAddr)
//	redStore := prefix.NewStore(store, srcValPrefix)
//	res, err = query.Paginate(redStore, pagination, func(key []byte, value []byte) error {
//		storeKey := types.GetUBDKeyFromValIndexKey(append(srcValPrefix, key...))
//		storeValue := store.Get(storeKey)
//		red, err := types.UnmarshalRED(k.cdc, storeValue)
//		if err != nil {
//			return err
//		}
//		redels = append(redels, red)
//		return nil
//	})
//
//	return redels, res, err
//}
//
//func queryAllRedelegations(store sdk.KVStore, k Querier, req *types.QueryRedelegationsRequest) (redels types.Redelegations, res *query.PageResponse, err error) {
//	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	redStore := prefix.NewStore(store, types.GetREDsKey(delAddr))
//	res, err = query.Paginate(redStore, req.Pagination, func(key []byte, value []byte) error {
//		redelegation, err := types.UnmarshalRED(k.cdc, value)
//		if err != nil {
//			return err
//		}
//		redels = append(redels, redelegation)
//		return nil
//	})
//
//	return redels, res, err
//}
