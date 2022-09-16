package worker

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	abci "github.com/tendermint/tendermint/abci/types"
)

/*
Event accumulation:
Prepare data for backend.
Goal: good update set for decimal-models.

1) Balance changes
records like (address, address_type, coin_symbol, amount, operation_type)
address_type - 'single' - common from mnemonic, 'multisig'(?) - multisignature wallet...
operation_type: "+"/"-"/"fee"
!!! +fee deduction
from coin sends, multisig sends, nft reserving, swaps, etc...

2) Coin changes:
last event EditCoin group by coin

3) NFT owners changes ()
(denom,id,subtokens,old owner,new owner)

4) Replace owner (legacy events)

5) multisig events as is
*/

type EventAccumulator struct {
	MultisigCreateWallets []MultisigCreateWallet `json:"multisig_create_wallets"`
	CoinsCreates          []EventCreateCoin      `json:"coin_creates"`
	// [address][coin_symbol]amount changes
	BalancesChanges map[string]map[string]sdkmath.Int `json:"balances_changes"`
	// [coin_symbol]
	CoinUpdates map[string]EventUpdateCoin `json:"coin_updates"`
	CoinEdits   map[string]EventEditCoin   `json:"coin_edits"`
	// replace legacy
	LegacyReown        map[string]string    `json:"legacy_reown"`
	LegacyReturnNFT    []LegacyReturnNFT    `json:"legacy_return_nft"`
	LegacyReturnWallet []LegacyReturnWallet `json:"legacy_return_multisig"`
	// nft
	NFTMints []EventMintNFT `json:"nft_mints"`
}

func NewEventAccumulator() *EventAccumulator {
	return &EventAccumulator{
		BalancesChanges: make(map[string]map[string]sdkmath.Int),
		CoinUpdates:     make(map[string]EventUpdateCoin),
		CoinEdits:       make(map[string]EventEditCoin),
		LegacyReown:     make(map[string]string),
	}
}

func (ea *EventAccumulator) AddEvent(event abci.Event, txHash string, blockId int64) error {
	procFunc, ok := eventProcessors[event.Type]
	if !ok {
		return fmt.Errorf("processor for event '%s' not found", event.Type)
	}
	return procFunc(ea, event, txHash, blockId)
}

type processFunc func(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error

var eventProcessors = map[string]processFunc{
	// coins
	"decimal.coin.v1.EventCreateCoin":  processEventCreateCoin,
	"decimal.coin.v1.EventUpdateCoin":  processEventUpdateCoin,
	"decimal.coin.v1.EventEditCoin":    processEventEditCoin,
	"decimal.coin.v1.EventSendCoin":    processEventSendCoin,
	"decimal.coin.v1.EventBuySellCoin": processEventBuySellCoin,
	"decimal.coin.v1.EventBurnCoin":    processEventBurnCoin,
	"decimal.coin.v1.EventRedeemCheck": processEventRedeemCheck,
	// fee
	"decimal.fee.v1.EventPayCommission": processEventPayCommission,
	// legacy
	"decimal.legacy.v1.EventLegacyReturnCoin":   processEventLegacyReturnCoin,
	"decimal.legacy.v1.EventLegacyReturnNFT":    processEventLegacyReturnNFT,
	"decimal.legacy.v1.EventLegacyReturnWallet": processEventLegacyReturnWallet,
	// multisig
	"decimal.multisig.v1.EventCreateWallet": processEventCreateWallet,
	// nft
	"decimal.nft.v1.EventMintNFT":     processEventMintNFT,
	"decimal.nft.v1.EventTransferNFT": processEventTransferNFT,
	// swap
	// stub for cosmos events
	"coin_spent":      processStub,
	"coin_received":   processStub,
	"transfer":        processStub,
	"message":         processStub,
	"proposer_reward": processStub,
	"commission":      processStub,
	"rewards":         processStub,
	"tx":              processStub,
}

// stub to skip internal cosmos events
func processStub(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	return nil
}

// account balances
// gathered from all transactions, amount can be negative
func (ea *EventAccumulator) addBalanceChange(address string, symbol string, amount sdkmath.Int) {
	balance, ok := ea.BalancesChanges[address]
	if !ok {
		ea.BalancesChanges[address] = map[string]sdkmath.Int{symbol: amount}
		return
	}
	knownChange, ok := balance[symbol]
	if !ok {
		knownChange = sdkmath.ZeroInt()
	}
	balance[symbol] = knownChange.Add(amount)
	ea.BalancesChanges[address] = balance
}

// set tx hash for some messages
func (ea *EventAccumulator) setTxHash(txs []Tx) error {
	// hash of CreateCoin
	// hash of MintNFT, TransferNFT

	return nil
}
