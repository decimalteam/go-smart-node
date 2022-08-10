package main

func copyParams(gs *GenesisNew, gsSource *GenesisNew) {
	gs.GenesisTime = gsSource.GenesisTime
	gs.AppHash = gsSource.AppHash
	gs.ChainID = gsSource.ChainID
	gs.InitalHeight = gsSource.InitalHeight
	gs.ConsensusParam = gsSource.ConsensusParam
	// modules
	gs.AppState.Auth.Params = gsSource.AppState.Auth.Params
}
