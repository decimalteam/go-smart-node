package dsc

import (
	"encoding/json"
	"fmt"
	"time"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

const basePower = 1000

// GenesisStateWithValSet returns a new genesis state with the validator set
func GenesisStateWithValSet(cdc codec.Codec, genesisState map[string]json.RawMessage,
	valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) (map[string]json.RawMessage, error) {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = cdc.MustMarshalJSON(authGenesis)

	validators := make([]validatortypes.Validator, 0, len(valSet.Validators))
	delegations := make([]validatortypes.Delegation, 0, len(valSet.Validators))

	bondSum := sdk.NewCoins()
	baseStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1000)))

	var lastPowers []validatortypes.LastValidatorPower
	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pubkey: %w", err)
		}

		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			return nil, fmt.Errorf("failed to create new any: %w", err)
		}

		validator := validatortypes.Validator{
			OperatorAddress: sdk.ValAddress(val.Address).String(),
			RewardAddress:   sdk.AccAddress(val.Address).String(),
			ConsensusPubkey: pkAny,
			Online:          true,
			Jailed:          false,
			Status:          validatortypes.BondStatus_Bonded,
			Description:     validatortypes.Description{},
			UnbondingHeight: int64(0),
			UnbondingTime:   time.Unix(0, 0).UTC(),
			Commission:      sdk.ZeroDec(),
		}
		validators = append(validators, validator)
		stake := validatortypes.NewStakeCoin(baseStake)
		bondSum = bondSum.Add(baseStake)
		delegations = append(delegations, validatortypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), stake))
		lastPowers = append(lastPowers, validatortypes.LastValidatorPower{
			Address: validator.OperatorAddress,
			Power:   basePower,
		})
	}
	// set validators and delegations
	validatorGenesis := validatortypes.NewGenesisState(validatortypes.DefaultParams(), validators, delegations, lastPowers)
	genesisState[validatortypes.ModuleName] = cdc.MustMarshalJSON(validatorGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	totalSupply = totalSupply.Add(bondSum...)

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(validatortypes.BondedPoolName).String(),
		Coins:   bondSum,
	})

	var sendActive []banktypes.SendEnabled
	sendActive = append(sendActive, banktypes.SendEnabled{
		Denom:   sdk.DefaultBondDenom,
		Enabled: true,
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, sendActive)
	genesisState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesis)

	return genesisState, nil
}
