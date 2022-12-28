package app

// UpgradeList is application upgrade table. Different for different environments
var UpgradeList = []UpgradeCreator{
	{"https://repo.decimalchain.com/12828401", FixSendUpgradeHandlerCreator},
}
