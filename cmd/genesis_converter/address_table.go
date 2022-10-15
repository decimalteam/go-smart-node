package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	cosmosAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

// helper structure for old to new address conversion
type AddressTable struct {
	data      map[string]string
	multisigs map[string]bool
	// name_of_module -> new address
	modules map[string]moduleInfo
}

type moduleInfo struct {
	address     string
	permissions []string
}

func NewAddressTable() AddressTable {
	return AddressTable{make(map[string]string), make(map[string]bool), make(map[string]moduleInfo)}
}

func (at *AddressTable) AddAddress(oldAddress string, pubKey []byte) error {
	var newAddress string
	var err error
	if len(pubKey) > 0 {
		newPubKey := ethsecp256k1.PubKey{Key: pubKey}
		newAddress, err = bech32.ConvertAndEncode("dx", newPubKey.Address())
		if err != nil {
			return err
		}
	}
	at.data[oldAddress] = newAddress
	return nil
}

func (at *AddressTable) AddMultisig(oldAddress string) {
	at.multisigs[oldAddress] = true
}

func (at *AddressTable) GetAddress(oldAddress string) string {
	return at.data[oldAddress]
}

func (at *AddressTable) IsMultisig(oldAddress string) bool {
	return at.multisigs[oldAddress]
}

func (at *AddressTable) InitModules() {
	// known module accounts
	at.modules = map[string]moduleInfo{
		"erc20": {
			address:     moduleNameToAddress("erc20"),
			permissions: []string{"minter", "burner"},
		},
		"bonded_tokens_pool": {
			address:     moduleNameToAddress("bonded_tokens_pool"),
			permissions: []string{"burner", "staking"},
		},
		"not_bonded_tokens_pool": {
			address:     moduleNameToAddress("not_bonded_tokens_pool"),
			permissions: []string{"burner", "staking"},
		},
		"inflation": {
			address:     moduleNameToAddress("inflation"),
			permissions: []string{"minter"},
		},
		"gov": {
			address:     moduleNameToAddress("gov"),
			permissions: []string{"burner"},
		},
		"distribution": {
			address:     moduleNameToAddress("distribution"),
			permissions: []string{},
		},
		"incentives": {
			address:     moduleNameToAddress("incentives"),
			permissions: []string{"minter", "burner"},
		},
		"coin": {
			address:     moduleNameToAddress("coin"),
			permissions: []string{"minter", "burner"},
		},
		"fee_collector": {
			address:     moduleNameToAddress("fee_collector"),
			permissions: []string{"burner", "minter"},
		},
		"validator": {
			address:     moduleNameToAddress("validator"),
			permissions: []string{"burner", "minter"},
		},
		"reserved_pool": {
			address:     moduleNameToAddress("reserved_pool"),
			permissions: []string{"minter", "burner"},
		},
		"legacy_coin_pool": {
			address:     moduleNameToAddress("legacy_coin_pool"),
			permissions: []string{},
		},
		"atomic_swap_pool": {
			address:     moduleNameToAddress("atomic_swap_pool"),
			permissions: []string{"minter", "burner"},
		},
	}
}

func (at *AddressTable) GetModule(name string) moduleInfo {
	return at.modules[name]
}

func moduleNameToAddress(name string) string {
	address, err := bech32.ConvertAndEncode("dx", cosmosAuthTypes.NewModuleAddress(name))
	if err != nil {
		panic(fmt.Sprintf("moduleNameToAddress(%s) = %s", name, err.Error()))
	}
	return address
}
