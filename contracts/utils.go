package contracts

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts/tokenCenter"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmkeeper "github.com/decimalteam/ethermint/x/evm/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// evm coin center events
const (
	NameOfSlugForGetAddressTokenCenter     = "token-center"
	NameOfSlugForGetAddressNftCenter       = "nft-center"
	NameOfSlugForGetAddressDelegation      = "delegation"
	NameOfSlugForGetAddressWDEL            = "wdel"
	NameOfSlugForGetAddressMasterValidator = "master-validator"
	EventChangeTokenCenter                 = "ContractAdded"

	// DRC20MethodCreateToken defines the create method for DRC20 token
	DRC20MethodCreateToken = "createToken"
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

func GetAddressFromContractCenter(
	ctx sdk.Context,
	evmKeeper *evmkeeper.Keeper,
	nameOfAddress string,
) (string, error) {
	contractCenter, _ := tokenCenter.TokenMetaData.GetAbi()
	contract := common.HexToAddress(GetContractCenter(ctx.ChainID()))
	methodCall := "getAddress"
	// Address token center
	res, err := evmKeeper.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, methodCall, nameOfAddress)
	if err != nil {
		return new(common.Address).Hex(), err
	}
	data, err := contractCenter.Unpack(methodCall, res.Ret)
	fmt.Println(data)
	if len(data) == 0 {
		return new(common.Address).Hex(), err
	}
	return data[0].(common.Address).String(), err
}

func GetContractCenter(chainID string) string {
	if helpers.IsMainnet(chainID) {
		return "0x2r32432"
	} else if helpers.IsTestnet(chainID) {
		return "0x464eB51b5965f4520B7180E2cC7805c55f9cefDA"
	} else {
		return "0x0e877d02eac44a0d4c2a33c70d44591e2081397e"
	}
}
