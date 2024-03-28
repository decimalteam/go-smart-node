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
		return "0xbF7BB2C9466D88D9224e9A4380287a22D2f2d3ac"
	}
}
