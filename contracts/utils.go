package contracts

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

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
		return "0xbE8f533E4A894a8ef0E467296B3B76d309aeF3c8"
	}
}
