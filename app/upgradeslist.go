package app

// UpgradeList is application upgrade table. Different for different environments
var UpgradeList = []UpgradeCreator{
	{"https://testnet-repo.decimalchain.com/6489301", FixSendUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/7377901", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/7490501", DummyUpgradeHandlerCreator},
	{"https://testnet-repo.decimalchain.com/9421801", DummyUpgradeHandlerCreator},
}
