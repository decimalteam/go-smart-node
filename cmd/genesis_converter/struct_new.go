package main

import (
	"encoding/base64"
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

type GenesisNew struct {
	GenesisTime    string      `json:"genesis_time"`
	ChainID        string      `json:"chain_id"`
	InitalHeight   string      `json:"initial_height"`
	AppHash        string      `json:"app_hash"`
	ConsensusParam interface{} `json:"consensus_params"`

	AppState struct {
		Auth struct {
			// can be EthAccount or ModuleAccount
			Accounts []interface{} `json:"accounts"`
			Params   interface{}   `json:"params"`
			/*
				Params   struct {
					MaxMemoCharacters      string `json:"max_memo_characters"`       // 256
					SigVerifyCostEd25519   string `json:"sig_verify_cost_ed25519"`   // 590
					SigVerifyCostSecp256k1 string `json:"sig_verify_cost_secp256k1"` // 1000
					TxSigLimit             string `json:"tx_sig_limit"`              // 7
					TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`     // 10
				} `json:"params"`
			*/
		} `json:"auth"`
		Bank struct {
			Params   interface{}  `json:"params"`
			Balances []BalanceNew `json:"balances"`
		} `json:"bank"`
		Multisig struct {
			Wallets []WalletNew `json:"wallets"`
		} `json:"multisig"`
		Coin struct {
			Params         interface{}        `json:"params"`
			Coins          []FullCoinNew      `json:"coins"`
			LegacyBalances []LegacyBalanceNew `json:"legacyBalances"`
		} `json:"coin"`
		NFT struct {
			Collections []CollectionNew `json:"collections"`
			NFTs        []NFTNew        `json:"nfts"`
		} `json:"nft"`
	} `json:"app_state"`
}

///////////////////////////
// Account
///////////////////////////

type AccountNew struct {
	Typ         string `json:"@type"`
	BaseAccount struct {
		AccountNumber string `json:"account_number"`
		Address       string `json:"address"`
		PublicKey     struct {
			Typ string `json:"@type"`
			Key string `json:"key"`
		}
		Sequence string `json:"sequence"`
	} `json:"base_account"`
	CodeHash string `json:"code_hash"`
}

// convert only regular accounts
func AccountO2N(acc AccountOld) (AccountNew, error) {
	var res = AccountNew{
		Typ:      "/ethermint.types.v1.EthAccount",
		CodeHash: "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470", // keccak256 of empty string
	}
	res.BaseAccount.AccountNumber = acc.Value.AccountNumber
	res.BaseAccount.Sequence = acc.Value.Sequence
	// public key
	res.BaseAccount.PublicKey.Typ = "/ethermint.crypto.v1.ethsecp256k1.PubKey"
	pubkeyBytes, err := extractPubKey(acc.Value.PublicKey)
	if err != nil {
		return AccountNew{}, err
	}
	res.BaseAccount.PublicKey.Key = base64.RawStdEncoding.EncodeToString(pubkeyBytes)
	// address
	pubkey := ethsecp256k1.PubKey{Key: pubkeyBytes}
	res.BaseAccount.Address, err = bech32.ConvertAndEncode(config.Bech32PrefixAccAddr, pubkey.Address())
	if err != nil {
		return AccountNew{}, err
	}
	return res, nil
}

type ModuleAccountNew struct {
	Typ         string `json:"@type"` // "/cosmos.auth.v1beta1.ModuleAccount"
	BaseAccount struct {
		AccountNumber string      `json:"account_number"`
		Address       string      `json:"address"`
		PublicKey     interface{} `json:"pub_key"` // null,
		Sequence      string      `json:"sequence"`
	} `json:"base_account"`
	Name       string   `json:"name"`
	Permission []string `json:"permissions"`
}

func ModuleAccountO2N(acc AccountOld, addrTable *AddressTable) (ModuleAccountNew, error) {
	modInfo := addrTable.GetModule(acc.Value.Name)
	if modInfo.address == "" {
		return ModuleAccountNew{}, fmt.Errorf("unknown module name '%s'", acc.Value.Name)
	}
	var res = ModuleAccountNew{
		Typ:        "/cosmos.auth.v1beta1.ModuleAccount",
		Name:       acc.Value.Name,
		Permission: modInfo.permissions,
	}
	res.BaseAccount.AccountNumber = acc.Value.AccountNumber
	res.BaseAccount.Address = modInfo.address
	res.BaseAccount.PublicKey = nil
	res.BaseAccount.Sequence = "0"
	return res, nil
}

///////////////////////////
// Bank
///////////////////////////

type BalanceNew struct {
	Address string    `json:"address"`
	Coins   sdk.Coins `json:"coins"`
}

///////////////////////////
// Coins
///////////////////////////

type FullCoinNew struct {
	CRR         string `json:"crr"`
	Creator     string `json:"creator"`
	Identity    string `json:"identity"`
	LimitVolume string `json:"limit_volume"`
	Reserve     string `json:"reserve"`
	Symbol      string `json:"symbol"`
	Title       string `json:"title"`
	Volume      string `json:"volume"`
}

type LegacyBalanceNew struct {
	Address string    `json:"address"`
	Coins   sdk.Coins `json:"coins"`
}

func FullCoinO2N(coin FullCoinOld, addrTable *AddressTable) FullCoinNew {
	return FullCoinNew{
		CRR:         coin.CRR,
		Creator:     addrTable.GetAddress(coin.Creator),
		Identity:    coin.Identity,
		LimitVolume: coin.LimitVolume,
		Reserve:     coin.Reserve,
		Symbol:      coin.Symbol,
		Title:       coin.Title,
		Volume:      coin.Volume,
	}
}

///////////////////////////
// Multisig
///////////////////////////

type WalletNew struct {
	Address      string   `json:"address"`
	Owners       []string `json:"owners"`
	Threshold    string   `json:"threshold"`
	Weights      []string `json:"weights"`
	LegacyOwners []string `json:"legacyOwners,omitempty"`
}

func WalletO2N(wallet WalletOld, addrTable *AddressTable) WalletNew {
	var hasLegacy = false
	var result = WalletNew{
		Address: wallet.Address,
	}
	result.Owners = make([]string, len(wallet.Owners))
	result.LegacyOwners = make([]string, len(wallet.Owners))
	result.Weights = wallet.Weights
	result.Threshold = wallet.Threshold
	for i := range wallet.Owners {
		newAddress := addrTable.GetAddress(wallet.Owners[i])
		if newAddress == "" {
			result.LegacyOwners[i] = wallet.Owners[i]
			hasLegacy = true
		} else {
			result.Owners[i] = newAddress
		}
	}
	if !hasLegacy {
		result.LegacyOwners = nil
	}
	return result
}

///////////////////////////
// NFT
///////////////////////////

type CollectionNew struct {
	Denom string   `json:"denom"`
	NFTs  []string `json:"NFTs"`
}

type NFTNew struct {
	ID        string     `json:"ID"`
	AllowMint bool       `json:"allowMint"`
	Creator   string     `json:"creator"`
	Reserve   string     `json:"reserve"`
	TokenURI  string     `json:"tokenURI"`
	Owners    []OwnerNew `json:"owners"`
}

type OwnerNew struct {
	SubTokenIDs []string `json:"SubTokenIDs"`
	Address     string   `json:"address"`
}
