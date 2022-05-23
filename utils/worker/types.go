package worker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

////////////////////////////////////////////////////////////////
// FetchTxs
////////////////////////////////////////////////////////////////

type FetchTxsTxResult struct {
	Code      uint64        `json:"code"`
	Data      interface{}   `json:"data"`
	Log       string        `json:"log"`
	Info      string        `json:"info"`
	GasWanted string        `json:"gasWanted"`
	GasUsed   string        `json:"gasUsed"`
	Events    []interface{} `json:"events"`
	Codespace string        `json:"codespace"`
}

type FetchTxsTx struct {
	Hash     string           `json:"hash"`
	Height   string           `json:"height"`
	Index    uint64           `json:"index"`
	Tx       string           `json:"tx"`
	TxResult FetchTxsTxResult `json:"tx_result"`
	Proof    interface{}      `json:"proof"`
}

type FetchTxsResult struct {
	Count string       `json:"total_count"`
	Txs   []FetchTxsTx `json:"txs"`
}

type FetchTxsResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      int64          `json:"id"`
	Result  FetchTxsResult `json:"result"`
}

////////////////////////////////////////////////////////////////
// BlockInfo
////////////////////////////////////////////////////////////////

type BlockInfoResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int64           `json:"id"`
	Result  BlockInfoResult `json:"result"`
	Error   interface{}     `json:"error"`
}

type BlockInfoResult struct {
	BlockId interface{}    `json:"block_id"`
	Block   BlockInfoBlock `json:"block"`
}

type BlockInfoBlock struct {
	Header     interface{}        `json:"header"`
	Data       BlockInfoBlockData `json:"data"`
	Evidence   interface{}        `json:"evidence"`
	LastCommit interface{}        `json:"last_commit"`
}

type BlockInfoBlockData struct {
	Txs []string `json:"txs"`
}

////////////////////////////////////////////////////////////////
// BlockResults
////////////////////////////////////////////////////////////////

type BlockResultsResponse struct {
	Jsonrpc string             `json:"jsonrpc"`
	Id      int64              `json:"id"`
	Result  BlockResultsResult `json:"result"`
}

type BlockResultsResult struct {
	Height           string      `json:"height"`
	TxResults        interface{} `json:"tx_results"`
	BeginBlockEvents []Event     `json:"begin_block_events"`
	EndBlockEvents   []Event     `json:"end_block_events"`
}

type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

////////////////////////////////////////////////////////////////
// BlockSize
////////////////////////////////////////////////////////////////

type BlockSizeResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int64           `json:"id"`
	Result  BlockSizeResult `json:"result"`
}

type BlockSizeResult struct {
	LastHeight string      `json:"last_height"`
	BlockMetas []BlockMeta `json:"block_metas"`
}

type BlockMeta struct {
	BlockId   interface{} `json:"block_id"`
	BlockSize string      `json:"block_size"`
	NumTxs    string      `json:"num_txs"`
	Header    interface{} `json:"header"`
}

////////////////////////////////////////////////////////////////
// GetStatus
////////////////////////////////////////////////////////////////

type GetStatusResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int64           `json:"id"`
	Result  GetStatusResult `json:"result"`
}

type GetStatusResult struct {
	NodeInfo interface{}       `json:"node_info"`
	SyncInfo GetStatusSyncInfo `json:"sync_info"`
}

type GetStatusSyncInfo struct {
	LatestBlockHash   string `json:"latest_block_hash"`
	LatestAppHash     string `json:"latest_app_hash"`
	LatestBlockHeight string `json:"latest_block_height"`
}

////////////////////////////////////////////////////////////////
// Black and transactions structures
////////////////////////////////////////////////////////////////

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
