package app

import "bitbucket.org/decimalteam/go-smart-node/utils/helpers"

// is application upgrade table. Different for different environments

var UpgradeListDevnet = []UpgradeCreator{
	{"https://devnet-repo.decimalchain.com/523001", DummyUpgradeHandlerCreator},
}
var UpgradeListTestnet = []UpgradeCreator{
	{"https://testnet-repo.decimalchain.com/6489301", FixSendUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/7377901", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/7490501", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/9421801", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/9434301", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/10229801", MigrationUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/10328801", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/10337801", MigrationUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/15069701", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/15586201", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/15698601", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/16379201", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/16406401", TransferDaoAndVals},
	{"https://testnet-repo.decimalchain.com/17515701", DummyUpgradeHandlerCreator},

}
var UpgradeListMainnet = []UpgradeCreator{
	{"https://repo.decimalchain.com/12830301", FixSendUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/13798601", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/15656601", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/22280701", MigrationUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/22372801", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/22466601", ValidatorDuplicatesHandlerCreator},
	{"https://repo.decimalchain.com/23003701", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/23116601", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/24537501", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/25812301", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/27239601", DummyUpgradeHandlerCreator},
}

func GetUpgradeList(chainID string) []UpgradeCreator {
	if helpers.IsMainnet(chainID) {
		return UpgradeListMainnet
	} else if helpers.IsTestnet(chainID) {
		return UpgradeListTestnet
	} else {
		return UpgradeListDevnet
	}
}
