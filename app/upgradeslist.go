package app

// UpgradeList is application upgrade table. Different for different environments
var UpgradeList = []UpgradeCreator{
	{"https://repo.decimalchain.com/12830301", FixSendUpgradeHandlerCreator},
	{"https://repo.decimalchain.com/13798601", DummyUpgradeHandlerCreator},
}
