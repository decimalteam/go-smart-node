package app

import (
	"time"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

// is application upgrade table. Different for different environments

var UpgradeListDevnet = []UpgradeCreator{
	{"https://devnet-repo.decimalchain.com/523001", DummyUpgradeHandlerCreator},
}
// TODO(rewardPerBlock): at deploy, add the following entry once the contract upgrade height
// is known (must be at or after the master-validator contract upgrade that introduces the
// rewardPerBlock field/getter). The EVM hook keeps node state in sync afterwards; this
// one-time handler seeds any value already set on the contract:
//   UpgradeListTestnet: {"https://testnet-repo.decimalchain.com/<HEIGHT>", RewardPerBlockSyncHandlerCreator},
//   UpgradeListMainnet: {"https://repo.decimalchain.com/<HEIGHT>", RewardPerBlockSyncHandlerCreator},

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
	{"https://testnet-repo.decimalchain.com/17514701", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/17576701", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/17621701", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/17628701", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/17700701", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/19354701", CombinedTestnetUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/20546401", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/20562701", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/20611201", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/20800401", DummyUpgradeHandlerCreator},
	// DEL redenomination (÷1000). The second field is the coordinated UTC restart time:
	// the chain halts after this upgrade until that instant so off-chain providers can
	// migrate first. When uncommenting, set <HEIGHT> and the restart time, and add the
	// "time" import to this file.
	{"https://testnet-repo.decimalchain.com/20927201", RedenominationUpgradeHandlerCreator(time.Date(2026, time.June, 25, 15, 0, 0, 0, time.UTC))},
	{"https://testnet-repo.decimalchain.com/20971140", RewardPerBlockSyncHandlerCreator},
	// Oracle re-enablement: set the fee oracle to the production key and create
	// its account, paired with the ante change allowing MsgUpdateCoinPrices.
	// Set <HEIGHT> at deploy time:
	// {"https://testnet-repo.decimalchain.com/<HEIGHT>", SetOracleUpgradeHandlerCreator},
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
	{"https://repo.decimalchain.com/27259201", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/27916701", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/27994701", TransferDaoAndVals},
	{"https://repo.decimalchain.com/28728701", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/29512333", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/31049701", CombinedMainnetUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/31080201", DummyUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/31295301", DummyUpgradeHandlerCreator},
	// DEL redenomination (÷1000). The second field is the coordinated UTC restart time:
	// the chain halts after this upgrade until that instant so off-chain providers can
	// migrate first. When uncommenting, set <HEIGHT> and the restart time, and add the
	// "time" import to this file.
	// {"https://repo.decimalchain.com/<HEIGHT>", RedenominationUpgradeHandlerCreator(time.Date(2026, time.June, 25, 15, 0, 0, 0, time.UTC))},
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
