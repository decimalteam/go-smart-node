package worker

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdkmath "cosmossdk.io/math"
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
	// [address][coin_symbol]amount changes
	BalancesChanges map[string]map[string]sdkmath.Int `json:"balances_changes"`
	// [denom]vr struct
	CoinsVR map[string]UpdateCoinVR `json:"coins_vr"`
	// [coin_symbol]
	CoinsCreates []EventCreateCoin          `json:"-"`
	CoinUpdates  map[string]EventUpdateCoin `json:"-"`
	// replace legacy
	LegacyReown        map[string]string    `json:"-"`
	LegacyReturnNFT    []LegacyReturnNFT    `json:"-"`
	LegacyReturnWallet []LegacyReturnWallet `json:"-"`
	// multisig
	MultisigCreateWallets []MultisigCreateWallet `json:"-"`
	MultisigCreateTxs     []MultisigCreateTx     `json:"-"`
	MultisigSignTxs       []MultisigSignTx       `json:"-"`
	// nft
	//Collection    []EventUpdateCollection `json:"collection"`
	CreateToken   []EventCreateToken   `json:"-"`
	MintSubTokens []EventMintToken     `json:"-"`
	BurnSubTokens []EventBurnToken     `json:"-"`
	UpdateToken   []EventUpdateToken   `json:"-"`
	UpdateReserve []EventUpdateReserve `json:"-"`
	SendNFTs      []EventSendToken     `json:"-"`
	// swap
	ActivateChain   []EventActivateChain   `json:"-"`
	DeactivateChain []EventDeactivateChain `json:"-"`
	SwapInitialize  []EventSwapInitialize  `json:"-"`
	SwapRedeem      []EventSwapRedeem      `json:"-"`
}

func NewEventAccumulator() *EventAccumulator {
	return &EventAccumulator{
		BalancesChanges: make(map[string]map[string]sdkmath.Int),
		CoinUpdates:     make(map[string]EventUpdateCoin),
		CoinsVR:         make(map[string]UpdateCoinVR),
		LegacyReown:     make(map[string]string),
	}
}

func (ea *EventAccumulator) AddEvent(event abci.Event, txHash string) error {
	procFunc, ok := eventProcessors[event.Type]
	if !ok {
		return processStub(ea, event, txHash)
	}
	return procFunc(ea, event, txHash)
}

type processFunc func(ea *EventAccumulator, event abci.Event, txHash string) error

var eventProcessors = map[string]processFunc{
	// coins
	"decimal.coin.v1.EventCreateCoin":   processEventCreateCoin,
	"decimal.coin.v1.EventUpdateCoin":   processEventUpdateCoin,
	"decimal.coin.v1.EventUpdateCoinVR": processEventUpdateCoinVR,
	"decimal.coin.v1.EventSendCoin":     processEventSendCoin,
	"decimal.coin.v1.EventBuySellCoin":  processEventBuySellCoin,
	"decimal.coin.v1.EventBurnCoin":     processEventBurnCoin,
	"decimal.coin.v1.EventRedeemCheck":  processEventRedeemCheck,
	// fee
	"decimal.fee.v1.EventUpdateCoinPrices": processEventUpdatePrices,
	"decimal.fee.v1.EventPayCommission":    processEventPayCommission,
	// legacy
	"decimal.legacy.v1.EventReturnLegacyCoins":    processEventReturnLegacyCoins,
	"decimal.legacy.v1.EventReturnLegacySubToken": processEventReturnLegacySubToken,
	"decimal.legacy.v1.EventReturnMultisigWallet": processEventReturnMultisigWallet,
	// multisig
	"decimal.multisig.v1.EventCreateWallet":       processEventCreateWallet,
	"decimal.multisig.v1.EventCreateTransaction":  processEventCreateTransaction,
	"decimal.multisig.v1.EventSignTransaction":    processEventSignTransaction,
	"decimal.multisig.v1.EventConfirmTransaction": processEventConfirmTransaction,
	// nft
	"decimal.nft.v1.EventCreateCollection": processEventCreateCollection,
	"decimal.nft.v1.EventUpdateCollection": processEventCreateCollection,
	"decimal.nft.v1.EventCreateToken":      processEventCreateToken,
	"decimal.nft.v1.EventMintToken":        processEventMintNFT,
	"decimal.nft.v1.EventUpdateToken":      processEventUpdateToken,
	"decimal.nft.v1.EventUpdateReserve":    processEventUpdateReserve,
	"decimal.nft.v1.EventSendToken":        processEventSendNFT,
	"decimal.nft.v1.EventBurnToken":        processEventBurnNFT,
	// swap
	"decimal.swap.v1.EventActivateChain":   processEventActivateChain,
	"decimal.swap.v1.EventDeactivateChain": processEventDeactivateChain,
	"decimal.swap.v1.EventInitializeSwap":  processEventSwapInitialize,
	"decimal.swap.v1.EventRedeemSwap":      processEventSwapRedeem,
}

// stub to skip internal cosmos events
func processStub(ea *EventAccumulator, event abci.Event, txHash string) error {
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

// custom coin reserve or volume update
func (ea *EventAccumulator) addCoinVRChange(symbol string, vr UpdateCoinVR) {
	ea.CoinsVR[symbol] = vr
}

func (ea *EventAccumulator) addMintSubTokens(e EventMintToken) {
	ea.MintSubTokens = append(ea.MintSubTokens, e)
}

func (ea *EventAccumulator) addBurnSubTokens(e EventBurnToken) {
	ea.BurnSubTokens = append(ea.BurnSubTokens, e)
}

// set tx hash for some messages
func (ea *EventAccumulator) setTxHash(txs []Tx) error {
	// hash of CreateCoin
	// hash of MintNFT, TransferNFT

	return nil
}
