package main

import "fmt"

func copyParams(gs *GenesisNew, gsSource *GenesisNew) {
	gs.GenesisTime = gsSource.GenesisTime
	gs.AppHash = gsSource.AppHash
	gs.ChainID = gsSource.ChainID
	// initial height will be taken from last_height
	//gs.InitalHeight = gsSource.InitalHeight
	gs.ConsensusParam = gsSource.ConsensusParam
	// modules
	gs.AppState.Auth.Params = gsSource.AppState.Auth.Params
	gs.AppState.Coin.Params = gsSource.AppState.Coin.Params
	gs.AppState.Bank.Params = gsSource.AppState.Bank.Params
	gs.AppState.Validator.Params = gsSource.AppState.Validator.Params
	//
	gs.AppState.Genutil = gsSource.AppState.Genutil
	gs.AppState.Swap = gsSource.AppState.Swap
	gs.AppState.Authz = gsSource.AppState.Authz
	gs.AppState.Capability = gsSource.AppState.Capability
	gs.AppState.Crisis = gsSource.AppState.Crisis
	gs.AppState.Distribution = gsSource.AppState.Distribution
	gs.AppState.Evidence = gsSource.AppState.Evidence
	gs.AppState.Evm = gsSource.AppState.Evm
	gs.AppState.Feegrant = gsSource.AppState.Feegrant
	gs.AppState.Fee = gsSource.AppState.Fee
	gs.AppState.Gov = gsSource.AppState.Gov
	gs.AppState.Params = gsSource.AppState.Params
	gs.AppState.Slashing = gsSource.AppState.Slashing
	gs.AppState.Staking = gsSource.AppState.Staking
	gs.AppState.Upgrade = gsSource.AppState.Upgrade
	gs.AppState.Vesting = gsSource.AppState.Vesting
	gs.AppState.IBC = gsSource.AppState.IBC

	// Copy accounts and balances
	for _, acc := range gsSource.AppState.Auth.Accounts {
		var sourceAdr = extractAddress(acc)
		var accExists = false
		for _, accI := range gs.AppState.Auth.Accounts {
			if extractAddress(accI) == sourceAdr {
				accExists = true
				break
			}
		}
		if !accExists {
			fmt.Printf("copy account from source: %+v\n", acc)
			gs.AppState.Auth.Accounts = append(gs.AppState.Auth.Accounts, acc)
		} else {
			fmt.Printf("account '%s' exists. skip\n", sourceAdr)
		}
	}
	for _, bal := range gsSource.AppState.Bank.Balances {
		fmt.Printf("copy balance from source: %+v\n", bal)
		var balanceExists = false
		for i, b := range gs.AppState.Bank.Balances {
			if bal.Address == b.Address {
				gs.AppState.Bank.Balances[i].Coins = gs.AppState.Bank.Balances[i].Coins.Add(bal.Coins...)
				balanceExists = true
				break
			}
		}
		if !balanceExists {
			gs.AppState.Bank.Balances = append(gs.AppState.Bank.Balances, bal)
		}
	}
}

func extractAddress(acc interface{}) string {
	switch a := acc.(type) {
	case AccountNew:
		return a.BaseAccount.Address
	case ModuleAccountNew:
		return a.BaseAccount.Address
	case map[string]interface{}:
		v := a["base_account"].(map[string]interface{})
		return v["address"].(string)
	default:
		return ""
	}
}
