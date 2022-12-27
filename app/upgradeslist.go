package app

// UpgradeList is application upgrade table. Different for different environments
var UpgradeList = []UpgradeCreator{
	{"https://testnet-repo.decimalchain.com/6489301", FixSendUpgradeHandlerCreator},
}
