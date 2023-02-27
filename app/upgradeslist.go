package app

// UpgradeList is application upgrade table. Different for different environments
var UpgradeList = []UpgradeCreator{
	{"https://devnet-repo.decimalchain.com/968301", ExampleUpgradeHandlerCreator},
}
