package upgrade

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

func GetAddressForUpdate(chainID string) string {
	if helpers.IsMainnet(chainID) {
		return "d01f4j9efhhy8ehjkfddqr80gv5kkl3e5yu7jxqvg"
	} else if helpers.IsTestnet(chainID) {
		return "d01y7sex8yvrazyd8pljjxvnvpndaavn99tk2j52y"
	} else {
		return "d01y7sex8yvrazyd8pljjxvnvpndaavn99tk2j52y"
	}
}
