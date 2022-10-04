package worker

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	utilsEvents "bitbucket.org/decimalteam/go-smart-node/utils/events"
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feeTypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	legacyTypes "bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

func TestBalanceChanges(t *testing.T) {
	ea := NewEventAccumulator()
	ea.addBalanceChange("a", "coin1", sdkmath.NewInt(1))
	ea.addBalanceChange("b", "coin2", sdkmath.NewInt(2))
	ea.addBalanceChange("a", "coin1", sdkmath.NewInt(2))
	ea.addBalanceChange("a", "coin1", sdkmath.NewInt(1).Neg())
	ea.addBalanceChange("a", "coin2", sdkmath.NewInt(1).Neg())
	require.True(t, ea.BalancesChanges["a"]["coin1"].Equal(sdkmath.NewInt(2)))
	require.True(t, ea.BalancesChanges["a"]["coin2"].Equal(sdkmath.NewInt(-1)))
	require.True(t, ea.BalancesChanges["b"]["coin2"].Equal(sdkmath.NewInt(2)))
}

// coin
func TestCreateCoin(t *testing.T) {
	ea := NewEventAccumulator()
	tev := &coinTypes.EventCreateCoin{
		Sender:               "a",
		Denom:                "test",
		Title:                "title",
		CRR:                  10,
		InitialVolume:        sdkmath.NewInt(100).String(),
		InitialReserve:       sdkmath.NewInt(200).String(),
		LimitVolume:          sdkmath.NewInt(1000).String(),
		Identity:             "ident",
		CommissionCreateCoin: sdk.NewCoin("del", sdkmath.NewInt(9000)).String(),
	}

	ev, err := utilsEvents.TypedEventToEvent(tev)
	require.NoError(t, err)
	err = processEventCreateCoin(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)
	require.Equal(t, 1, len(ea.CoinsCreates))
	require.True(t, ea.BalancesChanges["a"]["del"].Equal(sdkmath.NewInt(9200).Neg())) // comission+reserve
	require.True(t, ea.BalancesChanges["a"]["test"].Equal(sdkmath.NewInt(100)))
	cc := ea.CoinsCreates[0]
	require.Equal(t, "a", cc.Creator)
	require.Equal(t, "test", cc.Denom)
	require.Equal(t, "title", cc.Title)
	require.Equal(t, "ident", cc.Avatar)
	require.Equal(t, uint64(10), cc.CRR)
	require.Equal(t, "ABCD", cc.TxHash)
	require.True(t, cc.Volume.Equal(sdkmath.NewInt(100)))
	require.True(t, cc.Reserve.Equal(sdkmath.NewInt(200)))
	require.True(t, cc.LimitVolume.Equal(sdkmath.NewInt(1000)))
}

func TestUpdateCoin(t *testing.T) {
	ea := NewEventAccumulator()
	ev, err := utilsEvents.TypedEventToEvent(&coinTypes.EventUpdateCoin{
		Sender:      "a",
		Denom:       "test",
		LimitVolume: "1000",
		Identity:    "aaa",
	})
	require.NoError(t, err)
	err = processEventUpdateCoin(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)
	ev, err = utilsEvents.TypedEventToEvent(&coinTypes.EventUpdateCoin{
		Sender:      "a",
		Denom:       "test",
		LimitVolume: "10000",
		Identity:    "bbb",
	})
	require.NoError(t, err)
	err = processEventUpdateCoin(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)

	require.Len(t, ea.CoinUpdates, 1)
	require.Equal(t, "10000", ea.CoinUpdates["test"].LimitVolume.String())
	require.Equal(t, "bbb", ea.CoinUpdates["test"].Avatar)
}

func TestEditCoin(t *testing.T) {
	ea := NewEventAccumulator()
	ev, err := utilsEvents.TypedEventToEvent(&coinTypes.EventUpdateCoinVR{
		Denom:   "test",
		Volume:  "1000",
		Reserve: "2000",
	})
	require.NoError(t, err)
	err = processEventUpdateCoinVR(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)
	ev, err = utilsEvents.TypedEventToEvent(&coinTypes.EventUpdateCoinVR{
		Denom:   "test",
		Volume:  "1001",
		Reserve: "2001",
	})
	require.NoError(t, err)
	err = processEventUpdateCoinVR(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)

	require.Len(t, ea.CoinsVR, 1)
	require.True(t, ea.CoinsVR["test"].Volume.Equal(sdk.NewInt(1001)))
	require.True(t, ea.CoinsVR["test"].Reserve.Equal(sdk.NewInt(2001)))
}

func TestSendCoin(t *testing.T) {
	ea := NewEventAccumulator()
	ev, err := utilsEvents.TypedEventToEvent(&coinTypes.EventSendCoin{
		Sender:    "a",
		Recipient: "b",
		Coin:      "1000abc",
	})
	require.NoError(t, err)
	err = processEventSendCoin(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)
	ev, err = utilsEvents.TypedEventToEvent(&coinTypes.EventSendCoin{
		Sender:    "b",
		Recipient: "a",
		Coin:      "1ggg",
	})
	require.NoError(t, err)
	err = processEventSendCoin(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)

	require.True(t, ea.BalancesChanges["a"]["abc"].Equal(sdk.NewInt(1000).Neg()))
	require.True(t, ea.BalancesChanges["b"]["abc"].Equal(sdk.NewInt(1000)))
	require.True(t, ea.BalancesChanges["a"]["ggg"].Equal(sdk.NewInt(1)))
	require.True(t, ea.BalancesChanges["b"]["ggg"].Equal(sdk.NewInt(1).Neg()))
}

func TestBuySellCoin(t *testing.T) {
	ea := NewEventAccumulator()
	ev, err := utilsEvents.TypedEventToEvent(&coinTypes.EventBuySellCoin{
		Sender:     "a",
		CoinToBuy:  "10aaa",
		CoinToSell: "20bbb",
	})
	require.NoError(t, err)
	err = processEventBuySellCoin(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)

	require.True(t, ea.BalancesChanges["a"]["aaa"].Equal(sdk.NewInt(10)))
	require.True(t, ea.BalancesChanges["a"]["bbb"].Equal(sdk.NewInt(20).Neg()))
}

func TestBurnCoin(t *testing.T) {
	ea := NewEventAccumulator()
	ev, err := utilsEvents.TypedEventToEvent(&coinTypes.EventBurnCoin{
		Sender: "a",
		Coin:   "10aaa",
	})
	require.NoError(t, err)
	err = processEventBurnCoin(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)

	require.True(t, ea.BalancesChanges["a"]["aaa"].Equal(sdk.NewInt(10).Neg()))
}

// legacy
func TestReturnLegacyCoins(t *testing.T) {
	ea := NewEventAccumulator()
	tev := &legacyTypes.EventReturnLegacyCoins{
		LegacyOwner: "a",
		Owner:       "b",
		Coins:       sdk.NewCoins(sdk.NewCoin("coin1", sdk.NewInt(1)), sdk.NewCoin("coin2", sdk.NewInt(2))),
	}
	ev, err := utilsEvents.TypedEventToEvent(tev)
	require.NoError(t, err)
	err = processEventReturnLegacyCoins(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)
	require.True(t, ea.BalancesChanges["b"]["coin1"].Equal(sdk.NewInt(1)))
	require.True(t, ea.BalancesChanges["b"]["coin2"].Equal(sdk.NewInt(2)))
}

func TestReturnMultisigWallet(t *testing.T) {
	ea := NewEventAccumulator()
	ev, err := utilsEvents.TypedEventToEvent(&legacyTypes.EventReturnMultisigWallet{
		LegacyOwner: "a",
		Owner:       "b",
		Wallet:      "c",
	})
	require.NoError(t, err)
	err = processEventReturnMultisigWallet(ea, abci.Event(ev), "ABCD")
	require.NoError(t, err)

	require.Len(t, ea.LegacyReturnWallet, 1)
	require.Equal(t, "a", ea.LegacyReturnWallet[0].LegacyOwner)
	require.Equal(t, "b", ea.LegacyReturnWallet[0].Owner)
	require.Equal(t, "c", ea.LegacyReturnWallet[0].Wallet)
}

// multisig
func TestMultisigCreateWallet(t *testing.T) {
	ea := NewEventAccumulator()
	tev := &multisigTypes.EventCreateWallet{
		Sender:    "aaa",
		Wallet:    "bbb",
		Owners:    []string{"a", "b", "c"},
		Weights:   []uint32{1, 2, 3},
		Threshold: 10,
	}
	ev, err := utilsEvents.TypedEventToEvent(tev)
	require.NoError(t, err)
	err = processEventCreateWallet(ea, abci.Event(ev), "")
	require.NoError(t, err)
	require.Equal(t, 1, len(ea.MultisigCreateWallets))
	require.Equal(t, 3, len(ea.MultisigCreateWallets[0].Owners))
}

func TestMultisigComplete(t *testing.T) {
	ea := NewEventAccumulator()
	tev := &multisigTypes.EventConfirmTransaction{
		Wallet:      "a",
		Transaction: "tx",
		Receiver:    "b",
		Coins:       sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(1)), sdk.NewCoin("btc", sdk.NewInt(2))),
	}
	ev, err := utilsEvents.TypedEventToEvent(tev)
	require.NoError(t, err)
	err = processEventConfirmTransaction(ea, abci.Event(ev), "")
	require.NoError(t, err)

	require.True(t, ea.BalancesChanges["a"]["del"].Equal(sdk.NewInt(1).Neg()))
	require.True(t, ea.BalancesChanges["a"]["btc"].Equal(sdk.NewInt(2).Neg()))
	require.True(t, ea.BalancesChanges["b"]["del"].Equal(sdk.NewInt(1)))
	require.True(t, ea.BalancesChanges["b"]["btc"].Equal(sdk.NewInt(2)))
}

func TestMultisigCreateTx(t *testing.T) {
	ea := NewEventAccumulator()
	coins := sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(1)), sdk.NewCoin("btc", sdk.NewInt(2)))
	tev := &multisigTypes.EventCreateTransaction{
		Sender:      "S",
		Wallet:      "a",
		Transaction: "tx",
		Receiver:    "b",
		Coins:       coins,
	}
	ev, err := utilsEvents.TypedEventToEvent(tev)
	require.NoError(t, err)
	err = processEventCreateTransaction(ea, abci.Event(ev), "")
	require.NoError(t, err)

	require.Len(t, ea.MultisigCreateTxs, 1)
	require.Equal(t, "S", ea.MultisigCreateTxs[0].Sender)
	require.Equal(t, "a", ea.MultisigCreateTxs[0].Wallet)
	require.Equal(t, "b", ea.MultisigCreateTxs[0].Receiver)
	require.Equal(t, "tx", ea.MultisigCreateTxs[0].Transaction)
	require.True(t, coins.IsEqual(ea.MultisigCreateTxs[0].Coins))
}

func TestMultisigSignTx(t *testing.T) {
	ea := NewEventAccumulator()
	tev := &multisigTypes.EventSignTransaction{
		Sender:        "S",
		Wallet:        "a",
		Transaction:   "tx",
		SignerWeight:  10,
		Confirmations: 20,
		Confirmed:     true,
	}
	ev, err := utilsEvents.TypedEventToEvent(tev)
	require.NoError(t, err)
	err = processEventSignTransaction(ea, abci.Event(ev), "")
	require.NoError(t, err)

	require.Len(t, ea.MultisigSignTxs, 1)
	require.Equal(t, "S", ea.MultisigSignTxs[0].Sender)
	require.Equal(t, "a", ea.MultisigSignTxs[0].Wallet)
	require.Equal(t, "tx", ea.MultisigSignTxs[0].Transaction)
	require.Equal(t, uint32(10), ea.MultisigSignTxs[0].SignerWeight)
	require.Equal(t, uint32(20), ea.MultisigSignTxs[0].Confirmations)
	require.Equal(t, true, ea.MultisigSignTxs[0].Confirmed)
}

func TestPayCommission(t *testing.T) {
	ea := NewEventAccumulator()
	ev, err := utilsEvents.TypedEventToEvent(&feeTypes.EventPayCommission{
		Payer: "pp",
		Coins: sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(1)), sdk.NewCoin("btc", sdk.NewInt(2))),
	})
	require.NoError(t, err)

	err = processEventPayCommission(ea, abci.Event(ev), "")
	require.NoError(t, err)

	require.True(t, ea.BalancesChanges["pp"]["del"].Equal(sdk.NewInt(1).Neg()))
	require.True(t, ea.BalancesChanges["pp"]["btc"].Equal(sdk.NewInt(2).Neg()))
}
