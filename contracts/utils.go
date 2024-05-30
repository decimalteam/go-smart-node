package contracts

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

func GetContractCenter(chainID string) string {
	if helpers.IsMainnet(chainID) {
		return "0x2r32432"
	} else if helpers.IsTestnet(chainID) {
		return "0x464eB51b5965f4520B7180E2cC7805c55f9cefDA"
	} else {
		return "0xa052da26a526e251db6390834009464ab0398ddc"
	}
}
