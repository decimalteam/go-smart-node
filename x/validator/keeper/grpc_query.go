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
	valStore := prefix.NewStore(store, types.GetREDsFromValSrcIndexKey(valAddr))

	redelegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, valStore, req.Pagination,
		func(key []byte, red *types.Redelegation) (*types.Redelegation, error) {
			return red, nil
		}, func() *types.Redelegation {
			return &types.Redelegation{}
		})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	redels := []types.Redelegation{}
	for _, red := range redelegations {
		redels = append(redels, *red)
	}

	return &types.QueryValidatorRedelegationsResponse{
		Redelegations: redels,
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
	valStore := prefix.NewStore(store, types.GetUBDsByValIndexKey(valAddr))

	undelegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, valStore, req.Pagination,
		func(key []byte, ubd *types.Undelegation) (*types.Undelegation, error) {
			return ubd, nil
		}, func() *types.Undelegation {
			return &types.Undelegation{}
		})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	undels := []types.Undelegation{}
	for _, ubd := range undelegations {
		undels = append(undels, *ubd)
	}

	return &types.QueryValidatorUndelegationsResponse{
		Undelegations: undels,
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
		del, err := types.UnmarshalDelegation(k.cdc, iterator.Value())
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

// Undelegations queries undelegations info for given validator delegator pair.
func (k Querier) Undelegations(c context.Context, req *types.QueryUndelegationsRequest) (*types.QueryUndelegationsResponse, error) {
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

	var result []types.Undelegation
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetUBDByValIndexKey(delAddr, valAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		red, err := types.UnmarshalUBD(k.cdc, iterator.Value())
		if err != nil {
			return nil, err
		}
		result = append(result, red)
	}

	return &types.QueryUndelegationsResponse{Undelegations: result}, nil
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
