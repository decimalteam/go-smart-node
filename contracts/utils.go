package contracts

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts/center"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmkeeper "github.com/decimalteam/ethermint/x/evm/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

// evm coin center events
const (
	NameOfSlugForGetAddressTokenCenter     = "token-center"
	NameOfSlugForGetAddressNftCenter       = "nft-center"
	NameOfSlugForGetAddressDelegation      = "delegation"
	NameOfSlugForGetAddressDelegationNft   = "delegation-nft"
	NameOfSlugForGetAddressWDEL            = "wdel"
	NameOfSlugForGetAddressMasterValidator = "master-validator"
	EventChangeTokenCenter                 = "ContractAdded"
)

// MasterValidatorValidatorAddedMeta represents a ReserveUpdated event raised by the Contracts contract.
type MasterValidatorValidatorAddedMeta struct {
	OperatorAddress string `json:"operator_address"`
	RewardAddress   string `json:"reward_address"`
	ConsensusPubkey string `json:"consensus_pubkey"`
	Description     struct {
		Moniker         string `json:"moniker"`
		Identity        string `json:"identity"`
		Website         string `json:"website"`
		SecurityContact string `json:"security_contact"`
		Details         string `json:"details"`
	} `json:"description"`
	Commission string `json:"commission"`
}

func UnpackInputsData(v interface{}, inputs abi.Arguments, data []byte) error {
	unpacked, err := inputs.Unpack(data)
	if err != nil {
		return err
	}
	return inputs.Copy(v, unpacked)
}

// UnpackLog unpacks a retrieved log into the provided output structure.
func UnpackLog(abiUsed *abi.ABI, out interface{}, event string, log *ethTypes.Log) error {
	if len(log.Data) > 0 {
		if err := abiUsed.UnpackIntoInterface(out, event, log.Data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range abiUsed.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return abi.ParseTopics(out, indexed, log.Topics[1:])
}

func GetAddressFromContractCenter(
	ctx sdk.Context,
	evmKeeper *evmkeeper.Keeper,
	nameOfAddress string,
) (string, error) {
	contractCenter, _ := center.CenterMetaData.GetAbi()
	contract := common.HexToAddress(GetContractCenter(ctx.ChainID()))
	methodCall := "getContractAddress"
	// Address token center
	res, err := evmKeeper.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, methodCall, nameOfAddress)
	if err != nil {
		return new(common.Address).Hex(), err
	}
	data, err := contractCenter.Unpack(methodCall, res.Ret)
	if len(data) == 0 {
		return new(common.Address).Hex(), err
	}
	return data[0].(common.Address).String(), err
}

func GetIsMigration(
	ctx sdk.Context,
	evmKeeper *evmkeeper.Keeper,
) (bool, error) {
	contractCenter, _ := center.CenterMetaData.GetAbi()
	contract := common.HexToAddress(GetContractCenter(ctx.ChainID()))
	methodCall := "isMigrating"
	// Address token center
	res, err := evmKeeper.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, methodCall)
	if err != nil {
		return false, err
	}
	data, err := contractCenter.Unpack(methodCall, res.Ret)
	fmt.Println(data)
	if len(data) == 0 {
		return false, err
	}
	return data[0].(bool), err
}

func GetContractCenter(chainID string) string {
	if helpers.IsMainnet(chainID) {
		return "0xc108715a06f76caa96fa2c943ebf05159c29a87d"
	} else if helpers.IsTestnet(chainID) {
		return "0xc6e67e6b7fa068dca595d2b6d29378d1d2158b6f"
	} else {
		return "0xe5268fd6a4d041f20cbb92c662ceff1efe4c861e"
	}
}
