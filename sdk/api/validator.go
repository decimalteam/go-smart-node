package api

import (
	"context"

	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
)

type Validator = validatortypes.Validator
type Delegation = validatortypes.Delegation
type Stake = validatortypes.Stake

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
