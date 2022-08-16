package main

import (
	"encoding/base64"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisOld struct {
	AppState struct {
		Auth struct {
			Accounts []AccountOld `json:"accounts"`
		} `json:"auth"`
		Multisig struct {
			Transactions []TransactionOld `json:"txs"`
			Wallets      []WalletOld      `json:"wallets"`
		} `json:"multisig"`
		Coin struct {
			Coins []FullCoinOld `json:"coins"`
		} `json:"coin"`
		NFT struct {
			Collections map[string]CollectionOld `json:"collections"`
			SubTokens   []SubTokenOld            `json:"sub_tokens"`
		} `json:"nft"`
		Validator struct {
			Delegations    []DelegationOld      `json:"delegations"`
			DelegationsNFT []DelegationNFT      `json:"delegations_nft"`
			UndondingNFT   []UnbondingNFTRecord `json:"nft_unbonding_delegations"`
			Unbondings     []UnbondingRecord    `json:"unbonding_delegations"`
		} `json:"validator"`
	} `json:"app_state"`
}

///////////////////////////
// Account
///////////////////////////

const accountTypeRegular = "cosmos-sdk/Account"
const accountTypeModule = "cosmos-sdk/ModuleAccount"

type AccountOld struct {
	Typ   string `json:"type"` // cosmos-sdk/ModuleAccount |
	Value struct {
		Name          string    `json:"name"` //human-readable name for module accounts
		AccountNumber string    `json:"account_number"`
		Address       string    `json:"address"`
		Coins         sdk.Coins `json:"coins"`
		// public key value is base64 encoded bytes
		PublicKey interface{} `json:"public_key"` // key can be null, empty string for module account or map[string]interface{}
		Sequence  string      `json:"sequence"`
	} `json:"value"`
}

type CoinOld struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

func extractPubKey(v interface{}) ([]byte, error) {
	pkDesc, ok := v.(map[string]interface{})
	if !ok {
		return []byte{}, fmt.Errorf("no cast to map[string]string")
	}
	pkEnc, ok := pkDesc["value"].(string)
	if !ok {
		return []byte{}, fmt.Errorf("no value in public key")
	}
	return base64.RawStdEncoding.DecodeString(pkEnc)
}

///////////////////////////
// Coin
///////////////////////////

type FullCoinOld struct {
	CRR         string `json:"constant_reserve_ratio"`
	Creator     string `json:"creator"`
	Identity    string `json:"identity"`
	LimitVolume string `json:"limit_volume"`
	Reserve     string `json:"reserve"`
	Symbol      string `json:"symbol"`
	Title       string `json:"title"`
	Volume      string `json:"volume"`
}

///////////////////////////
// Multisig
///////////////////////////

type WalletOld struct {
	Address   string   `json:"address"`
	Owners    []string `json:"owners"`
	Threshold string   `json:"threshold"`
	Weights   []string `json:"weights"`
}

type TransactionOld struct {
	Coins     sdk.Coins `json:"coins"`
	ID        string    `json:"id"`
	Receiver  string    `json:"receiver"`
	CreatedAt string    `json:"string"`
	Wallet    string    `json:"wallet"`
	Signers   []string  `json:"signers"`
}

///////////////////////////
// NFT
///////////////////////////

type CollectionOld struct {
	Denom string            `json:"denom"`
	NFT   map[string]NFTOld `json:"nfts"`
}

type NFTOld struct {
	ID        string                `json:"id"`
	Creator   string                `json:"creator"`
	AllowMint bool                  `json:"allow_mint"`
	Reserve   string                `json:"reserve"`
	TokenURI  string                `json:"token_uri"`
	Owners    map[string][]OwnerOld `json:"owners"`
}

type OwnerOld struct {
	Address     string   `json:"address"`
	SubTokenIds []uint64 `json:"sub_token_ids"`
}

type SubTokenOld struct {
	Denom   string `json:"collection_denom"`
	NftID   string `json:"nft_id"`
	Reserve string `json:"reserve"`
	ID      string `json:"token_id"`
}

///////////////////////////
// Validator
///////////////////////////
type DelegationNFT struct {
	Coin             sdk.Coin `json:"coin"`
	DelegatorAddress string   `json:"delegator_address"`
	Denom            string   `json:"denom"`
	SubTokenIds      []string `json:"sub_token_ids"`
	TokenID          string   `json:"token_id"`
	Validator        string   `json:"validator_address"`
}

type UnbondingNFTRecord struct {
	DelegatorAddress string              `json:"delegator_address"`
	Validator        string              `json:"validator_address"`
	Entries          []UnbondingNFTEntry `json:"entries"`
}

type UnbondingNFTEntry struct {
	// TODO: complete
	Denom       string   `json:"denom"`
	TokenID     string   `json:"token_id"`
	SubTokenIds []string `json:"sub_token_ids"`
}

type DelegationOld struct {
	Coin             sdk.Coin `json:"coin"`
	DelegatorAddress string   `json:"delegator_address"`
	Validator        string   `json:"validator_address"`
	TokensBase       string   `json:"tokens_base"`
}

type UnbondingRecord struct {
	DelegatorAddress string           `json:"delegator_address"`
	Validator        string           `json:"validator_address"`
	Entries          []UnbondingEntry `json:"entries"`
}

type UnbondingEntry struct {
	// TODO: complete
	Value struct {
		Balance        sdk.Coin `json:"balance"`
		InitialBalance sdk.Coin `json:"initial_balance"`
	} `json:"value"`
}
