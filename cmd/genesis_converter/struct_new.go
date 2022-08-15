package main

import (
	"encoding/base64"
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
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
		} `json:"auth"`
		Bank struct {
			Params   interface{}  `json:"params"`
			Balances []BalanceNew `json:"balances"`
		} `json:"bank"`
		Multisig struct {
			Transactions []TransactionNew `json:"transactions"`
			Wallets      []WalletNew      `json:"wallets"`
		} `json:"multisig"`
		Coin struct {
			Params interface{}   `json:"params"`
			Coins  []FullCoinNew `json:"coins"`
		} `json:"coin"`
		NFT struct {
			Collections []CollectionNew         `json:"collections"`
			NFTs        []NFTNew                `json:"nfts"`
			SubTokens   map[string]SubTokensNew `json:"subTokens"`
		} `json:"nft"`
		Legacy struct {
			LegacyRecords []LegacyRecordNew `json:"legacyRecords"`
		} `json:"legacy"`
		//
		Genutils struct {
			Gentxs []interface{} `json:"gen_txs"`
		} `json:"genutil"`
		// other modules
		Authz        interface{} `json:"authz"`
		Capability   interface{} `json:"capability"`
		Crisis       interface{} `json:"crisis"`
		Distribution interface{} `json:"distribution"`
		Epochs       interface{} `json:"epochs"`
		Erc20        interface{} `json:"erc20"`
		Evidence     interface{} `json:"evidence"`
		Evm          interface{} `json:"evm"`
		Feegrant     interface{} `json:"feegrant"`
		Feemarket    interface{} `json:"feemarket"`
		Gov          interface{} `json:"gov"`
		Incentives   interface{} `json:"incentives"`
		Inflation    interface{} `json:"inflation"`
		Params       interface{} `json:"params"`
		Slashing     interface{} `json:"slashing"`
		Staking      interface{} `json:"staking"`
		Upgrade      interface{} `json:"upgrade"`
		Vesting      interface{} `json:"vesting"`
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
		} `json:"pub_key"`
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
// Legacy
///////////////////////////
type LegacyRecordNew struct {
	Address string      `json:"address"`
	Coins   sdk.Coins   `json:"coins"`
	NFTs    []NFTRecord `json:"nfts"`
	Wallets []string    `json:"wallets"`
}

type NFTRecord struct {
	Denom string `json:"denom"`
	ID    string `json:"id"`
}

type LegacyRecords struct {
	data map[string]*LegacyRecordNew
}

func NewLegacyRecords() *LegacyRecords {
	return &LegacyRecords{make(map[string]*LegacyRecordNew)}
}

func (rs *LegacyRecords) AddCoins(address string, coins sdk.Coins) {
	rec, ok := rs.data[address]
	if !ok {
		rec = &LegacyRecordNew{Address: address}
	}
	rec.Coins = rec.Coins.Add(coins...)
	rs.data[address] = rec
}

func (rs *LegacyRecords) AddNFT(address string, denom, id string) {
	rec, ok := rs.data[address]
	if !ok {
		rec = &LegacyRecordNew{Address: address}
	}
	rec.NFTs = append(rec.NFTs, NFTRecord{Denom: denom, ID: id})
	rs.data[address] = rec
}

func (rs *LegacyRecords) AddWallet(address string, wallet string) {
	rec, ok := rs.data[address]
	if !ok {
		rec = &LegacyRecordNew{Address: address}
	}
	rec.Wallets = append(rec.Wallets, wallet)
	rs.data[address] = rec
}

///////////////////////////
// Multisig
///////////////////////////

type WalletNew struct {
	Address   string   `json:"address"`
	Owners    []string `json:"owners"`
	Threshold string   `json:"threshold"`
	Weights   []string `json:"weights"`
}

func WalletO2N(wallet WalletOld, addrTable *AddressTable, legacyRecords *LegacyRecords) WalletNew {
	var result = WalletNew{
		Address: wallet.Address,
	}
	result.Owners = make([]string, len(wallet.Owners))
	result.Weights = wallet.Weights
	result.Threshold = wallet.Threshold
	for i := range wallet.Owners {
		newAddress := addrTable.GetAddress(wallet.Owners[i])
		if newAddress == "" {
			result.Owners[i] = wallet.Owners[i]
			legacyRecords.AddWallet(wallet.Owners[i], wallet.Address)
		} else {
			result.Owners[i] = newAddress
		}
	}
	return result
}

type TransactionNew struct {
	Coins     sdk.Coins `json:"coins"`
	ID        string    `json:"id"`
	Receiver  string    `json:"receiver"`
	CreatedAt string    `json:"created_at"`
	Wallet    string    `json:"wallet"`
	Signers   []string  `json:"signers"`
}

func TransactionO2N(tx TransactionOld, addrTable *AddressTable) TransactionNew {
	var result = TransactionNew{
		Coins:     tx.Coins,
		ID:        tx.ID,
		Receiver:  addrTable.GetAddress(tx.Receiver),
		CreatedAt: "0", // field looking unused
		Wallet:    tx.Wallet,
	}
	result.Signers = make([]string, len(tx.Signers))
	for i, s := range tx.Signers {
		if s == "" {
			continue
		}
		result.Signers[i] = addrTable.GetAddress(s)
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

type SubTokensNew struct {
	SubTokens []SubTokenNew `json:"subTokens"`
}

type SubTokenNew struct {
	ID      string `json:"ID"`
	Reserve string `json:"reserve"`
}
