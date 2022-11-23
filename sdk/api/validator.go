package api

import (
	"context"

	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
)

type (
	Validator    = validatortypes.Validator
	Delegation   = validatortypes.Delegation
	Redelegation = validatortypes.Redelegation
	Undelegation = validatortypes.Undelegation
	Stake        = validatortypes.Stake
	StakeType    = validatortypes.StakeType
	BondStatus   = validatortypes.BondStatus
)

const (
	// COIN defines the type for stakes in coin.
	StakeType_Coin StakeType = validatortypes.StakeType_Coin
	// NFT defines the type for stakes in NFT.
	StakeType_NFT StakeType = validatortypes.StakeType_NFT

	// UNSPECIFIED defines an invalid validator status.
	BondStatus_Unspecified BondStatus = validatortypes.BondStatus_Unspecified
	// UNBONDED defines a validator that is not bonded.
	BondStatus_Unbonded BondStatus = validatortypes.BondStatus_Unbonded
	// UNBONDING defines a validator that is unbonding.
	BondStatus_Unbonding BondStatus = validatortypes.BondStatus_Unbonding
	// BONDED defines a validator that is bonded.
	BondStatus_Bonded BondStatus = validatortypes.BondStatus_Bonded
)

func (api *API) Validators() ([]Validator, error) {
	client := validatortypes.NewQueryClient(api.grpcClient)
	vals := make([]Validator, 0)
	req := &validatortypes.QueryValidatorsRequest{
		Pagination: &query.PageRequest{Limit: queryLimit},
	}
	for {
		res, err := client.Validators(
			context.Background(),
			req,
		)
		if err != nil {
			return []Validator{}, err
		}
		if len(res.Validators) == 0 {
			break
		}
		vals = append(vals, res.Validators...)
		if len(res.Pagination.NextKey) == 0 {
			break
		}
		req.Pagination.Key = res.Pagination.NextKey
	}
	return vals, nil
}

func (api *API) ValidatorDelegations(validator string) ([]Delegation, error) {
	client := validatortypes.NewQueryClient(api.grpcClient)
	delegations := make([]Delegation, 0)
	req := &validatortypes.QueryValidatorDelegationsRequest{
		Validator:  validator,
		Pagination: &query.PageRequest{Limit: queryLimit},
	}
	for {
		res, err := client.ValidatorDelegations(
			context.Background(),
			req,
		)
		if err != nil {
			return []Delegation{}, err
		}
		delegations = append(delegations, res.Delegations...)
		if len(res.Pagination.NextKey) == 0 {
			break
		}
		req.Pagination.Key = res.Pagination.NextKey
	}
	return delegations, nil
}

func (api *API) ValidatorRedelegations(validator string) ([]Redelegation, error) {
	client := validatortypes.NewQueryClient(api.grpcClient)
	redelegations := make([]Redelegation, 0)
	req := &validatortypes.QueryValidatorRedelegationsRequest{
		Validator:  validator,
		Pagination: &query.PageRequest{Limit: queryLimit},
	}
	for {
		res, err := client.ValidatorRedelegations(
			context.Background(),
			req,
		)
		if err != nil {
			return []Redelegation{}, err
		}
		redelegations = append(redelegations, res.Redelegations...)
		if len(res.Pagination.NextKey) == 0 {
			break
		}
		req.Pagination.Key = res.Pagination.NextKey
	}
	return redelegations, nil
}

func (api *API) ValidatorUndelegations(validator string) ([]Undelegation, error) {
	client := validatortypes.NewQueryClient(api.grpcClient)
	undelegations := make([]Undelegation, 0)
	req := &validatortypes.QueryValidatorUndelegationsRequest{
		Validator:  validator,
		Pagination: &query.PageRequest{Limit: queryLimit},
	}
	for {
		res, err := client.ValidatorUndelegations(
			context.Background(),
			req,
		)
		if err != nil {
			return []Undelegation{}, err
		}
		undelegations = append(undelegations, res.Undelegations...)
		if len(res.Pagination.NextKey) == 0 {
			break
		}
		req.Pagination.Key = res.Pagination.NextKey
	}
	return undelegations, nil
}
