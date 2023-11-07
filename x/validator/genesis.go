package validator

import (
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// WriteValidators returns a slice of bonded genesis validators.
func WriteValidators(ctx sdk.Context, keeper keeper.Keeper) (vals []tmtypes.GenesisValidator, err error) {
	keeper.IterateLastValidators(ctx, func(_ int64, validator types.ValidatorI) (stop bool) {
		pk, err := validator.ConsPubKey()
		if err != nil {
			return true
		}
		tmPk, err := cryptocodec.ToTmPubKeyInterface(pk)
		if err != nil {
			return true
		}

		vals = append(vals, tmtypes.GenesisValidator{
			Address: sdk.ConsAddress(tmPk.Address()).Bytes(),
			PubKey:  tmPk,
			Power:   validator.ConsensusPower(),
			Name:    validator.GetMoniker(),
		})

		return false
	})

	return
}

// ValidateGenesis validates the provided validator genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators, etc.).
func ValidateGenesis(data *types.GenesisState) error {
	if err := validateGenesisStateValidators(data.Validators); err != nil {
		return err
	}
	valMap := make(map[string]bool)
	for _, val := range data.Validators {
		valMap[val.OperatorAddress] = true
	}
	if err := validateGenesisStateDelegations(data.Delegations, valMap); err != nil {
		return err
	}

	return data.Params.Validate()
}

func validateGenesisStateValidators(validators []types.Validator) error {
	addrMap := make(map[string]bool, len(validators))

	for i := 0; i < len(validators); i++ {
		val := validators[i]
		consPk, err := val.ConsPubKey()
		if err != nil {
			return err
		}

		strKey := string(consPk.Bytes())

		if _, ok := addrMap[strKey]; ok {
			consAddr, err := val.GetConsAddr()
			if err != nil {
				return err
			}
			return fmt.Errorf("duplicate validator in genesis state: moniker %v, address %v", val.Description.Moniker, consAddr)
		}

		if val.Jailed && val.IsBonded() {
			consAddr, err := val.GetConsAddr()
			if err != nil {
				return err
			}
			return fmt.Errorf("validator is bonded and jailed in genesis state: moniker %v, address %v", val.Description.Moniker, consAddr)
		}

		addrMap[strKey] = true
	}

	return nil
}

func validateGenesisStateDelegations(delegations []types.Delegation, validators map[string]bool) error {
	type delegationKey struct {
		delegator string
		validator string
		id        string
	}
	delegMap := make(map[delegationKey]bool)
	for _, delegation := range delegations {
		if _, err := sdk.AccAddressFromBech32(delegation.Delegator); err != nil {
			return fmt.Errorf("delegator '%s' invalid bech32: %s", delegation.Delegator, err.Error())
		}
		key := delegationKey{delegation.Delegator, delegation.Validator, delegation.Stake.ID}
		if delegMap[key] {
			return fmt.Errorf("duplicate delegation record %#v", key)
		}
		delegMap[key] = true
		if !validators[delegation.Validator] {
			return fmt.Errorf("validator %s does not exists (record %#v)", delegation.Validator, key)
		}
	}
	return nil
}

func validateGenesisStateUndelegations(undelegations []types.Undelegation, validators map[string]bool) error {
	type undelegationKey struct {
		delegator string
		validator string
	}
	undelegMap := make(map[undelegationKey]bool)
	for _, ubd := range undelegations {
		if _, err := sdk.AccAddressFromBech32(ubd.Delegator); err != nil {
			return fmt.Errorf("(un)delegator '%s' invalid bech32: %s", ubd.Delegator, err.Error())
		}
		key := undelegationKey{delegator: ubd.Delegator, validator: ubd.Validator}
		if undelegMap[key] {
			return fmt.Errorf("duplicate undelegation record %#v", key)
		}
		undelegMap[key] = true
		if !validators[ubd.Validator] {
			return fmt.Errorf("validator %s does not exists (record %#v)", ubd.Validator, key)
		}
	}
	return nil
}
