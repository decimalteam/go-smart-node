package main

func copyParams(gs *GenesisNew, gsSource *GenesisNew) {
	gs.GenesisTime = gsSource.GenesisTime
	gs.AppHash = gsSource.AppHash
	gs.ChainID = gsSource.ChainID
	gs.InitalHeight = gsSource.InitalHeight
	gs.ConsensusParam = gsSource.ConsensusParam
	// modules
	gs.AppState.Auth.Params = gsSource.AppState.Auth.Params
	gs.AppState.Coin.Params = gsSource.AppState.Coin.Params
	//
	gs.AppState.Swap = gsSource.AppState.Swap
	gs.AppState.Authz = gsSource.AppState.Authz
	gs.AppState.Capability = gsSource.AppState.Capability
	gs.AppState.Crisis = gsSource.AppState.Crisis
	gs.AppState.Distribution = gsSource.AppState.Distribution
	gs.AppState.Epochs = gsSource.AppState.Epochs
	gs.AppState.Erc20 = gsSource.AppState.Erc20
	gs.AppState.Evidence = gsSource.AppState.Evidence
	gs.AppState.Evm = gsSource.AppState.Evm
	gs.AppState.Feegrant = gsSource.AppState.Feegrant
	gs.AppState.Feemarket = gsSource.AppState.Feemarket
	gs.AppState.Gov = gsSource.AppState.Gov
	gs.AppState.Incentives = gsSource.AppState.Incentives
	gs.AppState.Inflation = gsSource.AppState.Inflation
	gs.AppState.Params = gsSource.AppState.Params
	gs.AppState.Slashing = gsSource.AppState.Slashing
	gs.AppState.Staking = gsSource.AppState.Staking
	gs.AppState.Upgrade = gsSource.AppState.Upgrade
	gs.AppState.Vesting = gsSource.AppState.Vesting
}
