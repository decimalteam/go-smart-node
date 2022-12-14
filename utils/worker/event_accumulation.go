package worker

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	evmtypes "github.com/evmos/ethermint/x/evm/types"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	legacytypes "bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	swaptypes "bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
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

var pool = map[string]bool{
	mustConvertAndEncode(authtypes.NewModuleAddress(authtypes.FeeCollectorName)):     false,
	mustConvertAndEncode(authtypes.NewModuleAddress(distrtypes.ModuleName)):          false,
	mustConvertAndEncode(authtypes.NewModuleAddress(stakingtypes.BondedPoolName)):    false,
	mustConvertAndEncode(authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName)): false,
	mustConvertAndEncode(authtypes.NewModuleAddress(govtypes.ModuleName)):            false,
	mustConvertAndEncode(authtypes.NewModuleAddress(evmtypes.ModuleName)):            false,
	mustConvertAndEncode(authtypes.NewModuleAddress(cointypes.ModuleName)):           false,
	mustConvertAndEncode(authtypes.NewModuleAddress(nfttypes.ReservedPool)):          false,
	mustConvertAndEncode(authtypes.NewModuleAddress(legacytypes.LegacyCoinPool)):     false,
	mustConvertAndEncode(authtypes.NewModuleAddress(swaptypes.SwapPool)):             false,
	mustConvertAndEncode(authtypes.NewModuleAddress(validatortypes.ModuleName)):      false,
	mustConvertAndEncode(authtypes.NewModuleAddress(feetypes.BurningPool)):           false,
}

type EventAccumulator struct {
	// [address][coin_symbol]amount changes
	BalancesChanges map[string]map[string]sdkmath.Int `json:"balances_changes"`
	// [denom]vr struct
	CoinsVR       map[string]UpdateCoinVR `json:"coins_vr"`
	PayCommission []EventPayCommission    `json:"pay_commission"`
	CoinsStaked   map[string]sdkmath.Int  `json:"coin_staked"`
	// [coin_symbol]
	CoinsCreates []EventCreateCoin          `json:"-"`
	CoinUpdates  map[string]EventUpdateCoin `json:"-"`
	// replace legacy
	LegacyReown        map[string]string    `json:"legacy_reown"`
	LegacyReturnNFT    []LegacyReturnNFT    `json:"legacy_return_nft"`
	LegacyReturnWallet []LegacyReturnWallet `json:"legacy_return_wallet"`
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
		CoinsStaked:     make(map[string]sdkmath.Int),
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
	// validator
	"decimal.validator.v1.EventDelegate":           processEventDelegate,
	"decimal.validator.v1.EventUndelegateComplete": processEventUndelegateComplete,
	"decimal.validator.v1.EventForceUndelegate":    processEventUndelegateComplete,
	"decimal.validator.v1.EventRedelegateComplete": processEventRedelegateComplete,
	"decimal.validator.v1.EventUpdateCoinsStaked":  processEventUpdateCoinsStaked,

	banktypes.EventTypeTransfer: processEventTransfer,
}

// stub to skip internal cosmos events
func processStub(ea *EventAccumulator, event abci.Event, txHash string) error {
	return nil
}

func processEventTransfer(ea *EventAccumulator, event abci.Event, txHash string) error {
	var (
		err        error
		coins      sdk.Coins
		sender     string
		receipient string
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case banktypes.AttributeKeySender:
			sender = string(attr.Value)
		case banktypes.AttributeKeyRecipient:
			receipient = string(attr.Value)
		case sdk.AttributeKeyAmount:
			coins, err = sdk.ParseCoinsNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse coins: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		}
	}

	for _, coin := range coins {
		if _, ok := pool[sender]; !ok {
			ea.addBalanceChange(sender, coin.Denom, coin.Amount.Neg())
		}
		if _, ok := pool[receipient]; !ok {
			ea.addBalanceChange(receipient, coin.Denom, coin.Amount)
		}
	}

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

func (ea *EventAccumulator) addCoinsStaked(e EventUpdateCoinsStaked) {
	ea.CoinsStaked[e.denom] = e.amount
}

func mustConvertAndEncode(address sdk.AccAddress) string {
	res, err := bech32.ConvertAndEncode(cmdcfg.Bech32PrefixAccAddr, address)
	if err != nil {
		panic(err)
	}

	return res
}
