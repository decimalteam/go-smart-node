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
		return "0x62252958E017A95ddc8051c36a578F6CCE3B8b1A"
	}
}
