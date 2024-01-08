package app

// UpgradeList is application upgrade table. Different for different environments
var UpgradeList = []UpgradeCreator{
	{"https://repo.decimalchain.com/522101", DummyUpgradeHandlerCreator},
}
