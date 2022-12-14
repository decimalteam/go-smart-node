package main

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func copyParams(gs *GenesisNew, gsSource *GenesisNew) {
	gs.GenesisTime = gsSource.GenesisTime
	gs.AppHash = gsSource.AppHash
	gs.ChainID = gsSource.ChainID
	// initial height will be taken from last_height
	// gs.InitalHeight = gsSource.InitalHeight
	gs.ConsensusParam = gsSource.ConsensusParam
	// modules
	gs.AppState.Auth.Params = gsSource.AppState.Auth.Params
	gs.AppState.Coin.Params = gsSource.AppState.Coin.Params
	gs.AppState.Bank.Params = gsSource.AppState.Bank.Params
	gs.AppState.Validator.Params = gsSource.AppState.Validator.Params
	gs.AppState.NFT.Params = gsSource.AppState.NFT.Params
	//
	// gs.AppState.Genutil = gsSource.AppState.Genutil
	gs.AppState.Swap = gsSource.AppState.Swap
	gs.AppState.Authz = gsSource.AppState.Authz
	gs.AppState.Capability = gsSource.AppState.Capability
	gs.AppState.Crisis = gsSource.AppState.Crisis
	gs.AppState.Distribution = gsSource.AppState.Distribution
	//gs.AppState.Evidence = gsSource.AppState.Evidence
	gs.AppState.Evm = gsSource.AppState.Evm
	gs.AppState.Feegrant = gsSource.AppState.Feegrant
	gs.AppState.Fee = gsSource.AppState.Fee
	gs.AppState.Gov = gsSource.AppState.Gov
	gs.AppState.Params = gsSource.AppState.Params
	//gs.AppState.Slashing = gsSource.AppState.Slashing
	//gs.AppState.Staking = gsSource.AppState.Staking
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

// fix bonded pool balance from staking
func fixBondedNotBondedPools(gs *GenesisNew) {
	// staking -> delegations[]: "delegator_address", "shares",
	bondings := sdk.NewCoins()
	notbondings := sdk.NewCoins()
	for _, del := range gs.AppState.Validator.Delegations {
		if del.Stake.Type == "STAKE_TYPE_COIN" {
			val := ValidatorNew{}
			for i, v := range gs.AppState.Validator.Validators {
				if v.OperatorAddress == del.Validator {
					val = gs.AppState.Validator.Validators[i]
					break
				}
			}
			switch val.Status {
			case "BOND_STATUS_BONDED":
				bondings = bondings.Add(del.Stake.Stake)
			case "BOND_STATUS_UNBONDED":
				notbondings = notbondings.Add(del.Stake.Stake)
			}
		}
	}
	for _, ubd := range gs.AppState.Validator.Undelegations {
		for _, entry := range ubd.Entries {
			if entry.Stake.Type == "STAKE_TYPE_COIN" {
				notbondings = notbondings.Add(entry.Stake.Stake)
			}
		}
	}

	for i := range gs.AppState.Bank.Balances {
		// "bonded_tokens_pool"
		if gs.AppState.Bank.Balances[i].Address == "d01fl48vsnmsdzcv85q5d2q4z5ajdha8yu3h9xgqa" {
			gs.AppState.Bank.Balances[i].Coins = bondings
		}
		// "not_bonded_tokens_pool"
		if gs.AppState.Bank.Balances[i].Address == "d01tygms3xhhs3yv487phx3dw4a95jn7t7lr96ekf" {
			gs.AppState.Bank.Balances[i].Coins = notbondings
		}
	}
}

func fixNFTPool(gs *GenesisNew) {
	expectedVolume := sdk.NewCoins()
	for _, coll := range gs.AppState.NFT.Collections {
		for _, token := range coll.Tokens {
			for _, sub := range token.SubTokens {
				expectedVolume = expectedVolume.Add(sub.Reserve)
			}
		}
	}
	for i := range gs.AppState.Bank.Balances {
		// "reserved_pool"
		if gs.AppState.Bank.Balances[i].Address == "d017epewz8nye288vvypc549pgr6mf52hlax5wry2" {
			gs.AppState.Bank.Balances[i].Coins = expectedVolume
		}
	}
}

func fixCoinVolumes(gs *GenesisNew) {
	summaryVolume := sdk.NewCoins()
	for i := range gs.AppState.Bank.Balances {
		summaryVolume = summaryVolume.Add(gs.AppState.Bank.Balances[i].Coins...)
	}
	for i, coinInfo := range gs.AppState.Coin.Coins {
		volume := summaryVolume.AmountOf(coinInfo.Symbol)
		if volume.IsZero() {
			fmt.Printf("ZERO amount for coin '%s'\n", coinInfo.Symbol)
			continue
		}
		gs.AppState.Coin.Coins[i].Volume = volume.String()
	}
}

func fixAccountNumbers(gs *GenesisNew) {
	for i, acc := range gs.AppState.Auth.Accounts {
		switch a := acc.(type) {
		case AccountNew:
			a.BaseAccount.AccountNumber = strconv.FormatInt(int64(i)+1, 10)
			gs.AppState.Auth.Accounts[i] = a
		case ModuleAccountNew:
			a.BaseAccount.AccountNumber = strconv.FormatInt(int64(i)+1, 10)
			gs.AppState.Auth.Accounts[i] = a
		case map[string]interface{}:
			a["base_account"].(map[string]interface{})["account_number"] = strconv.FormatInt(int64(i)+1, 10)
			gs.AppState.Auth.Accounts[i] = a
		}
	}

}

func fixDelegatedNFT(gs *GenesisNew, addrTable *AddressTable) {
	type nftKey struct {
		tokenID  string
		subtoken uint32
	}
	bpool := addrTable.GetModule("bonded_tokens_pool").address
	nbpool := addrTable.GetModule("not_bonded_tokens_pool").address
	delegationRecords := make(map[nftKey]string)
	validatorsStatus := make(map[string]bool) // true - online
	for _, val := range gs.AppState.Validator.Validators {
		validatorsStatus[val.OperatorAddress] = val.Online
	}

	for _, del := range gs.AppState.Validator.Delegations {
		if del.Stake.Type == "STAKE_TYPE_NFT" {
			owner := nbpool
			if validatorsStatus[del.Validator] {
				owner = bpool
			}
			for _, subID := range del.Stake.SubTokenIDs {
				delegationRecords[nftKey{del.Stake.ID, uint32(subID)}] = owner
			}
		}
	}
	for _, undel := range gs.AppState.Validator.Undelegations {
		for _, entry := range undel.Entries {
			if entry.Stake.Type == "STAKE_TYPE_NFT" {
				owner := nbpool
				for _, subID := range entry.Stake.SubTokenIDs {
					delegationRecords[nftKey{entry.Stake.ID, uint32(subID)}] = owner
				}
			}
		}
	}

	// fix
	for i, coll := range gs.AppState.NFT.Collections {
		for j, token := range coll.Tokens {
			for k, sub := range token.SubTokens {
				key := nftKey{token.ID, sub.ID}
				owner, ok := delegationRecords[key]
				if !ok {
					continue
				}
				gs.AppState.NFT.Collections[i].Tokens[j].SubTokens[k].Owner = owner
			}
		}
	}
}
