package main

import (
	"encoding/base64"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
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
			Collections []CollectionNew `json:"collections"`
			Params      interface{}     `json:"params"`
		} `json:"nft"`
		Legacy struct {
			LegacyRecords []LegacyRecordNew `json:"records"`
		} `json:"legacy"`
		//
		Genutil interface{} `json:"genutil"`
		// other modules
		Swap         interface{} `json:"swap"`
		Authz        interface{} `json:"authz"`
		Capability   interface{} `json:"capability"`
		Crisis       interface{} `json:"crisis"`
		Distribution interface{} `json:"distribution"`
		//Evidence     interface{} `json:"evidence"`
		Evm      interface{} `json:"evm"`
		Feegrant interface{} `json:"feegrant"`
		Fee      interface{} `json:"customfee"`
		Gov      interface{} `json:"gov"`
		IBC      interface{} `json:"ibc"`
		Params   interface{} `json:"params"`
		//Slashing     interface{} `json:"slashing"`
		//Staking      interface{} `json:"staking"`
		Upgrade   interface{} `json:"upgrade"`
		Vesting   interface{} `json:"vesting"`
		Validator struct {
			Validators          []ValidatorNew          `json:"validators"`
			Delegations         []DelegationNew         `json:"delegations"`
			Undelegations       []UndelegationNew       `json:"undelegations"`
			LastValidatorPowers []LastValidatorPowerNew `json:"last_validator_powers"`
			LastTotalPower      int64                   `json:"last_total_power"`
			Params              interface{}             `json:"params"`
		} `json:"validator"`
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
	res.BaseAccount.Address, err = bech32.ConvertAndEncode(cmdcfg.Bech32PrefixAccAddr, pubkey.Address())
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
	CRR         int64  `json:"crr"`
	Creator     string `json:"creator"`
	Identity    string `json:"identity"`
	LimitVolume string `json:"limit_volume"`
	Reserve     string `json:"reserve"`
	Symbol      string `json:"denom"`
	Title       string `json:"title"`
	Volume      string `json:"volume"`
}

func FullCoinO2N(coin FullCoinOld, addrTable *AddressTable) FullCoinNew {
	symbol := coin.Symbol
	if symbol == "tdel" || symbol == "del" {
		symbol = globalBaseDenom
	}
	crr, _ := strconv.ParseInt(coin.CRR, 10, 64)
	return FullCoinNew{
		CRR:         crr,
		Creator:     addrTable.GetAddress(coin.Creator),
		Identity:    coin.Identity,
		LimitVolume: coin.LimitVolume,
		Reserve:     coin.Reserve,
		Symbol:      symbol,
		Title:       coin.Title,
		Volume:      coin.Volume,
	}
}

// /////////////////////////
// Legacy
// /////////////////////////
type LegacyRecordNew struct {
	Address string    `json:"legacy_address"`
	Coins   sdk.Coins `json:"coins"`
	NFTs    []string  `json:"nfts"`
	Wallets []string  `json:"wallets"`
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
	rec.NFTs = append(rec.NFTs, id)
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
		if addrTable.IsMultisig(wallet.Owners[i]) {
			newAddress = wallet.Owners[i]
		}
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

func TransactionO2N(tx TransactionOld, addrTable *AddressTable, coinSymbols map[string]bool) TransactionNew {
	var result = TransactionNew{
		Coins:     filterCoins(tx.Coins, coinSymbols),
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
	Creator string     `json:"creator"`
	Denom   string     `json:"denom"`
	Supply  uint32     `json:"supply"`
	Tokens  []TokenNew `json:"tokens"`
}

type TokenNew struct {
	Creator   string        `json:"creator"`
	Denom     string        `json:"denom"`
	ID        string        `json:"id"`
	URI       string        `json:"uri"`
	Reserve   sdk.Coin      `json:"reserve"`
	AllowMint bool          `json:"allow_mint"`
	Minted    uint32        `json:"minted"`
	Burnt     uint32        `json:"burnt"`
	SubTokens []SubTokenNew `json:"sub_tokens"`
}

type SubTokenNew struct {
	ID      uint32   `json:"id"`
	Owner   string   `json:"owner"`
	Reserve sdk.Coin `json:"reserve"`
}

type NFTOwnerFixRecord struct {
	TokenID   string   `json:"token_id"`
	Owner     string   `json:"owner"`
	SubTokens []uint32 `json:"sub_tokens"`
}

// /////////////////////////
// Validator
// /////////////////////////
type ValidatorNew struct {
	OperatorAddress string `json:"operator_address"` // dxvaloper1
	RewardAddress   string `json:"reward_address"`   // dx1
	ConsensusPubKey struct {
		Type string `json:"@type"`
		Key  []byte `json:"key"`
	} `json:"consensus_pubkey"`
	Description struct {
		Details         string `json:"details"`
		Identity        string `json:"identity"`
		Moniker         string `json:"moniker"`
		SecurityContact string `json:"security_contact"`
		Website         string `json:"website"`
	} `json:"description"`
	Commission      string `json:"commission"`
	Status          string `json:"status"` // BOND_STATUS
	Online          bool   `json:"online"`
	Jailed          bool   `json:"jailed"`
	UnbondingHeight string `json:"unbonding_height"`
	UnbondingTime   string `json:"unbonding_time"`
}

func ValidatorO2N(valOld ValidatorOld, addrTable *AddressTable) (ValidatorNew, error) {
	var result ValidatorNew
	result.OperatorAddress = valOld.ValAddress
	newRewardAdr := addrTable.GetAddress(valOld.RewardAddress)
	result.RewardAddress = newRewardAdr
	// pubkey
	result.ConsensusPubKey.Type = "/cosmos.crypto.ed25519.PubKey"
	bz, err := sdk.GetFromBech32(valOld.PubKey, "dxvalconspub")
	if err != nil {
		return ValidatorNew{}, err
	}
	pk, err := legacy.PubKeyFromBytes(bz)
	if err != nil {
		return ValidatorNew{}, err
	}
	result.ConsensusPubKey.Key = pk.Bytes()
	// description
	result.Description.Details = valOld.Description.Details
	result.Description.Identity = valOld.Description.Identity
	result.Description.Moniker = valOld.Description.Moniker
	result.Description.SecurityContact = valOld.Description.SecurityContact
	result.Description.Website = valOld.Description.Website
	//
	result.Commission = valOld.Commission
	/*
		Unbonded  BondStatus = 0x00 -- BOND_STATUS_UNBONDED
		Unbonding BondStatus = 0x01 -- BOND_STATUS_UNBONDING
		Bonded    BondStatus = 0x02 -- BOND_STATUS_BONDED
	*/
	switch valOld.Status {
	case 0:
		result.Status = "BOND_STATUS_UNBONDED"
	case 1:
		result.Status = "BOND_STATUS_UNBONDING"
	case 2:
		result.Status = "BOND_STATUS_BONDED"
	default:
		return ValidatorNew{}, fmt.Errorf("unknown status code: %d", valOld.Status)
	}

	result.Online = valOld.Online
	result.Jailed = valOld.Jailed
	result.UnbondingHeight = valOld.UnbondingHeight
	result.UnbondingTime = valOld.UnbondingCompletionTime

	return result, nil
}

type DelegationNew struct {
	Delegator string   `json:"delegator"`
	Validator string   `json:"validator"`
	Stake     StakeNew `json:"stake"`
}

type StakeNew struct {
	Type        string   `json:"type"`
	ID          string   `json:"id"`
	Stake       sdk.Coin `json:"stake"`
	SubTokenIDs []int64  `json:"sub_token_ids"`
}

type UndelegationNew struct {
	Delegator string                 `json:"delegator"`
	Validator string                 `json:"validator"`
	Entries   []UndelegationEntryNew `json:"entries"`
}
type UndelegationEntryNew struct {
	CreationHeight int64    `json:"creation_height"`
	CompletionTime string   `json:"completion_time"`
	Stake          StakeNew `json:"stake"`
}

type RedelegationNew struct{}
type LastValidatorPowerNew struct {
	Address string `json:"address"`
	Power   int64  `json:"power"`
}

/*
  // COIN defines the type for stakes in coin.
  STAKE_TYPE_COIN = 1 [ (gogoproto.enumvalue_customname) = "Coin" ];
  // NFT defines the type for stakes in NFT.
  STAKE_TYPE_NFT = 2 [ (gogoproto.enumvalue_customname) = "NFT" ];
*/

func DelegationO2NCoin(delOld DelegationOld, coinSymbols map[string]bool, addrTable *AddressTable) (DelegationNew, error) {
	var delNew DelegationNew
	coins := filterCoins(sdk.NewCoins(delOld.Coin), coinSymbols)
	if coins.Len() == 0 {
		fmt.Printf("delegation: (%s, %s) unexisting coin '%s'\n", delOld.DelegatorAddress, delOld.Validator, delOld.Coin.Denom)
		return DelegationNew{}, nil
	}
	coin := coins[0]
	delNew.Validator = delOld.Validator
	newAdr := addrTable.GetAddress(delOld.DelegatorAddress)
	if newAdr == "" {
		return DelegationNew{}, fmt.Errorf("delegator '%s' has no new address", delOld.DelegatorAddress)
	}
	delNew.Delegator = newAdr
	delNew.Stake.ID = coin.Denom
	delNew.Stake.Stake = coin
	delNew.Stake.Type = "STAKE_TYPE_COIN"

	return delNew, nil
}

func DelegationO2NNFT(delOld DelegationNFTOld, addrTable *AddressTable) (DelegationNew, error) {
	var delNew DelegationNew

	delNew.Validator = delOld.Validator
	newAdr := addrTable.GetAddress(delOld.DelegatorAddress)
	if newAdr == "" {
		return DelegationNew{}, fmt.Errorf("delegator '%s' has no new address", delOld.DelegatorAddress)
	}
	delNew.Delegator = newAdr
	delNew.Stake.ID = delOld.TokenID
	delNew.Stake.Stake = delOld.Coin
	// fix for testnet
	if delNew.Stake.Stake.Denom == "tdel" || delNew.Stake.Stake.Denom == "del" {
		delNew.Stake.Stake.Denom = globalBaseDenom
	}
	for _, s := range delOld.SubTokenIds {
		id, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return DelegationNew{}, fmt.Errorf("subtoken of token '%s' error: %s", delOld.TokenID, err.Error())
		}
		delNew.Stake.SubTokenIDs = append(delNew.Stake.SubTokenIDs, id)
	}
	delNew.Stake.Type = "STAKE_TYPE_NFT"

	return delNew, nil
}

func UnbondingO2NCoin(ubdOld UnbondingRecordOld, coinSymbols map[string]bool, addrTable *AddressTable) (UndelegationNew, error) {
	var err error
	var ubdNew UndelegationNew

	newAdr := addrTable.GetAddress(ubdOld.DelegatorAddress)
	if newAdr == "" {
		return UndelegationNew{}, fmt.Errorf("delegator '%s' has no new address", ubdOld.DelegatorAddress)
	}
	ubdNew.Delegator = newAdr
	ubdNew.Validator = ubdOld.Validator

	for _, entryOld := range ubdOld.Entries {
		coins := filterCoins(sdk.NewCoins(entryOld.Value.Balance), coinSymbols)
		if coins.Len() == 0 {
			fmt.Printf("undonding: (%s, %s) unexisting coin '%s'\n", ubdOld.DelegatorAddress, ubdOld.Validator, entryOld.Value.Balance.Denom)
			continue
		}
		coin := coins[0]

		var entryNew UndelegationEntryNew
		entryNew.CompletionTime = entryOld.Value.CompletionTime
		entryNew.CreationHeight, err = strconv.ParseInt(entryOld.Value.CreationHeight, 10, 64)
		if err != nil {
			return UndelegationNew{}, fmt.Errorf("delegator '%s' creation_height error: %s", ubdOld.DelegatorAddress, err.Error())
		}
		entryNew.Stake.ID = coin.Denom
		entryNew.Stake.Stake = coin
		entryNew.Stake.Type = "STAKE_TYPE_COIN"

		ubdNew.Entries = append(ubdNew.Entries, entryNew)
	}
	return ubdNew, nil
}

func UnbondingO2NNFT(ubdOld UnbondingNFTRecordOld, addrTable *AddressTable) (UndelegationNew, error) {
	var err error
	var ubdNew UndelegationNew

	newAdr := addrTable.GetAddress(ubdOld.DelegatorAddress)
	if newAdr == "" {
		return UndelegationNew{}, fmt.Errorf("delegator '%s' has no new address", ubdOld.DelegatorAddress)
	}
	ubdNew.Delegator = newAdr
	ubdNew.Validator = ubdOld.Validator

	for _, entryOld := range ubdOld.Entries {
		var entryNew UndelegationEntryNew
		entryNew.CompletionTime = entryOld.CompletionTime
		entryNew.CreationHeight, err = strconv.ParseInt(entryOld.CreationHeight, 10, 64)
		if err != nil {
			return UndelegationNew{}, fmt.Errorf("delegator '%s' creation_height error: %s", ubdOld.DelegatorAddress, err.Error())
		}
		entryNew.Stake.ID = entryOld.TokenID
		entryNew.Stake.Stake = sdk.NewCoin(globalBaseDenom, entryOld.Balance.Amount)
		entryNew.Stake.Type = "STAKE_TYPE_NFT"
		for _, s := range entryOld.SubTokenIds {
			id, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return UndelegationNew{}, fmt.Errorf("subtoken of token '%s' error: %s", entryOld.TokenID, err.Error())
			}
			entryNew.Stake.SubTokenIDs = append(entryNew.Stake.SubTokenIDs, id)
		}
		ubdNew.Entries = append(ubdNew.Entries, entryNew)
	}
	return ubdNew, nil
}

func LastValidatorPowerO2N(pwrOld LastValidatorPowerOld) (LastValidatorPowerNew, error) {
	var err error
	var pwrNew LastValidatorPowerNew
	pwrNew.Address = pwrOld.Address
	pwrNew.Power, err = strconv.ParseInt(pwrOld.Power, 10, 64)
	if err != nil {
		return LastValidatorPowerNew{}, fmt.Errorf("validator power '%s' error: %s", pwrOld.Address, err.Error())
	}
	return pwrNew, nil
}
