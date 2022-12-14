package worker

import (
	web3common "github.com/ethereum/go-ethereum/common"
	web3hexutil "github.com/ethereum/go-ethereum/common/hexutil"
	web3types "github.com/ethereum/go-ethereum/core/types"

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
	ID                interface{}        `json:"id"`
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
	EVM               BlockEVM           `json:"evm"`
	StateChanges      EventAccumulator   `json:"state_changes"`
}

type TransactionEVM struct {
	Type             web3hexutil.Uint64  `json:"type"`
	Hash             web3common.Hash     `json:"hash"`
	Nonce            web3hexutil.Uint64  `json:"nonce"`
	BlockHash        web3common.Hash     `json:"blockHash"`
	BlockNumber      web3hexutil.Uint64  `json:"blockNumber"`
	TransactionIndex web3hexutil.Uint64  `json:"transactionIndex"`
	From             web3common.Address  `json:"from"`
	To               *web3common.Address `json:"to"`
	Value            *web3hexutil.Big    `json:"value"`
	Data             web3hexutil.Bytes   `json:"input"`
	Gas              web3hexutil.Uint64  `json:"gas"`
	GasPrice         *web3hexutil.Big    `json:"gasPrice"`

	// Optional
	ChainId    *web3hexutil.Big     `json:"chainId,omitempty"`              // EIP-155 replay protection
	AccessList web3types.AccessList `json:"accessList,omitempty"`           // EIP-2930 access list
	GasTipCap  *web3hexutil.Big     `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559 dynamic fee transactions
	GasFeeCap  *web3hexutil.Big     `json:"maxFeePerGas,omitempty"`         // EIP-1559 dynamic fee transactions
}

type BlockEVM struct {
	Header       *web3types.Header    `json:"header"`
	Transactions []*TransactionEVM    `json:"transactions"`
	Uncles       []*web3types.Header  `json:"uncles"`
	Receipts     []*web3types.Receipt `json:"receipts"`
}

type BlockData struct {
	Txs []Tx `json:"txs"`
}

type Tx struct {
	Hash      string        `json:"hash"`
	Log       []interface{} `json:"log"`
	Code      uint32        `json:"code"`
	Codespace string        `json:"codespace"`
	Data      interface{}   `json:"data"`
	GasUsed   int64         `json:"gas_used"`
	GasWanted int64         `json:"gas_wanted"`
	Info      TxInfo        `json:"info"`
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

type TxInfo struct {
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

////////////////////////////////////////////////////////////////

// type BlockEVM2 struct {
// 	Header       *BlockHeaderEVM   `json:"header"`
// 	Transactions []*TransactionEVM `json:"transactions"`
// 	Uncles       []*BlockHeaderEVM `json:"uncles"`
// 	Receipts     []*ReceiptEVM     `json:"receipts"`
// }

// type BlockHeaderEVM struct {
// 	ParentHash  web3common.Hash      `json:"parentHash"`
// 	UncleHash   web3common.Hash      `json:"sha3Uncles"`
// 	Coinbase    web3common.Address   `json:"miner"`
// 	Root        web3common.Hash      `json:"stateRoot"`
// 	TxHash      web3common.Hash      `json:"transactionsRoot"`
// 	ReceiptHash web3common.Hash      `json:"receiptsRoot"`
// 	Bloom       web3types.Bloom      `json:"logsBloom"`
// 	Difficulty  *big.Int             `json:"difficulty"`
// 	Number      *big.Int             `json:"number"`
// 	GasLimit    uint64               `json:"gasLimit"`
// 	GasUsed     uint64               `json:"gasUsed"`
// 	Time        uint64               `json:"timestamp"`
// 	Extra       []byte               `json:"extraData"`
// 	MixDigest   web3common.Hash      `json:"mixHash"`
// 	Nonce       web3types.BlockNonce `json:"nonce"`

// 	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
// 	BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`
// }

// type TransactionEVM struct {
// 	Type      uint8               `json:"type"`
// 	To        *web3common.Address `json:"to"`
// 	Value     *big.Int            `json:"value"`
// 	Input     []byte              `json:"input"`
// 	Nonce     uint64              `json:nonce`
// 	Gas       uint64              `json:gas`
// 	GasPrice  *big.Int            `json:"gasPrice"`
// 	GasTipCap *big.Int            `json:"maxPriorityFeePerGas"`
// 	GasFeeCap *big.Int            `json:"maxFeePerGas"`
// 	V         *big.Int            `json:"v"`
// 	R         web3common.Hash     `json:"r"`
// 	S         web3common.Hash     `json:"s"`
// 	Hash      web3common.Hash     `json:"hash"`
// }

// type ReceiptEVM struct {
// 	// Consensus fields: These fields are defined by the Yellow Paper
// 	Type              uint8            `json:"type,omitempty"`
// 	PostState         []byte           `json:"root"`
// 	Status            uint64           `json:"status"`
// 	CumulativeGasUsed uint64           `json:"cumulativeGasUsed" gencodec:"required"`
// 	Bloom             web3types.Bloom  `json:"logsBloom"         gencodec:"required"`
// 	Logs              []*web3types.Log `json:"logs"              gencodec:"required"`

// 	// Implementation fields: These fields are added by geth when processing a transaction.
// 	// They are stored in the chain database.
// 	TxHash          web3common.Hash    `json:"transactionHash" gencodec:"required"`
// 	ContractAddress web3common.Address `json:"contractAddress"`
// 	GasUsed         uint64             `json:"gasUsed" gencodec:"required"`

// 	// Inclusion information: These fields provide information about the inclusion of the
// 	// transaction corresponding to this receipt.
// 	BlockHash        web3common.Hash `json:"blockHash,omitempty"`
// 	BlockNumber      *big.Int        `json:"blockNumber,omitempty"`
// 	TransactionIndex uint            `json:"transactionIndex"`
// }

// func (w *Worker) convertBlockHeaderEVM(header *web3types.Header) *BlockHeaderEVM {
// 	return &BlockHeaderEVM{
// 		ParentHash:  header.ParentHash,
// 		UncleHash:   header.UncleHash,
// 		Coinbase:    header.Coinbase,
// 		Root:        header.Root,
// 		TxHash:      header.TxHash,
// 		ReceiptHash: header.ReceiptHash,
// 		Bloom:       header.Bloom,
// 		Difficulty:  header.Difficulty,
// 		Number:      header.Number,
// 		GasLimit:    header.GasLimit,
// 		GasUsed:     header.GasUsed,
// 		Time:        header.Time,
// 		Extra:       header.Extra,
// 		MixDigest:   header.MixDigest,
// 		Nonce:       header.Nonce,
// 		BaseFee:     header.BaseFee,
// 	}
// }

// func (w *Worker) convertBlockHeadersEVM(headers []*web3types.Header) (results []*BlockHeaderEVM) {
// 	for _, header := range headers {
// 		results = append(results, w.convertBlockHeaderEVM(header))
// 	}
// 	return
// }

// func (w *Worker) convertTransactionEVM(tx *web3types.Transaction) *TransactionEVM {
// 	v, r, s := tx.RawSignatureValues()
// 	return &TransactionEVM{
// 		Type:      tx.Type(),
// 		To:        tx.To(),
// 		Value:     tx.Value(),
// 		Input:     tx.Data(),
// 		Nonce:     tx.Nonce(),
// 		Gas:       tx.Gas(),
// 		GasPrice:  tx.GasPrice(),
// 		GasTipCap: tx.GasTipCap(),
// 		GasFeeCap: tx.GasFeeCap(),
// 		V:         v,
// 		R:         web3common.BigToHash(r),
// 		S:         web3common.BigToHash(s),
// 		Hash:      tx.Hash(),
// 	}
// }

// func (w *Worker) convertTransactionsEVM(txs []*web3types.Transaction) (results []*TransactionEVM) {
// 	for _, tx := range txs {
// 		results = append(results, w.convertTransactionEVM(tx))
// 	}
// 	return
// }

// func (w *Worker) convertReceiptEVM(receipt *web3types.Receipt) *ReceiptEVM {
// 	return &ReceiptEVM{
// 		Type:              receipt.Type,
// 		PostState:         receipt.PostState,
// 		Status:            receipt.Status,
// 		CumulativeGasUsed: receipt.CumulativeGasUsed,
// 		Bloom:             receipt.Bloom,
// 		Logs:              receipt.Logs,
// 		TxHash:            receipt.TxHash,
// 		ContractAddress:   receipt.ContractAddress,
// 		GasUsed:           receipt.GasUsed,
// 		BlockHash:         receipt.BlockHash,
// 		BlockNumber:       receipt.BlockNumber,
// 		TransactionIndex:  receipt.TransactionIndex,
// 	}
// }

// func (w *Worker) convertReceiptsEVM(receipts []*web3types.Receipt) (results []*ReceiptEVM) {
// 	for _, receipt := range receipts {
// 		results = append(results, w.convertReceiptEVM(receipt))
// 	}
// 	return
// }
