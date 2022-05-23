package worker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Block struct {
	Header            interface{}        `json:"header"`
	Data              BlockData          `json:"data"`
	Evidence          interface{}        `json:"evidence"`
	LastCommit        interface{}        `json:"last_commit"`
	Emission          string             `json:"emission"`
	Rewards           []ProposerReward   `json:"rewards"`
	CommissionRewards []CommissionReward `json:"commission_rewards"`
	BeginBlockEvents  []Event            `json:"begin_block_events"`
	EndBlockEvents    []Event            `json:"end_block_events"`
	Size              int                `json:"size"`
}

type BlockData struct {
	Txs []Tx `json:"txs"`
}

type Tx struct {
	Hash      string        `json:"hash"`
	Log       []interface{} `json:"log"`
	Code      uint32        `json:"code"`
	Data      interface{}   `json:"data"`
	GasUsed   int64         `json:"gas_used"`
	GasWanted int64         `json:"gas_wanted"`
	Tx        interface{}   `json:"info"`
}

type FailedTxLog struct {
	Log string `json:"log"`
}

type TxMsg struct {
	Type   string      `json:"type"`
	Params interface{} `json:"params"`
	From   []string    `json:"from"`
}

type TxFee struct {
	Gas    uint64    `json:"gas"`
	Amount sdk.Coins `json:"amount"`
}

type ParsedTx struct {
	Msgs []TxMsg `json:"msgs"`
	Memo string  `json:"memo"`
	Fee  TxFee   `json:"fee"`
}

type ProposerReward struct {
	Reward    string `json:"reward"`
	Address   string `json:"address"`
	Delegator string `json:"delegator"`
}

type CommissionReward struct {
	Amount        string `json:"amount"`
	Validator     string `json:"validator"`
	RewardAddress string `json:"reward_address"`
}
