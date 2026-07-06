package redenom_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/app/redenom"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// TestRedenominate_BankAndSupply validates the riskiest piece end-to-end through the
// real bank read path: the direct-store balance rewrite + supply reconciliation. It
// funds accounts with assorted "del" amounts (including sub-divisor dust), runs the
// full Redenominate, and asserts every balance is floor(old/1000), supply equals the
// sum of balances, and supply never exceeds floor(old supply).
func TestRedenominate_BankAndSupply(t *testing.T) {
	dscApp := app.Setup(t, false, nil)
	ctx := dscApp.BaseApp.NewContext(false, tmproto.Header{
		Height: 1,
		// A mapped chain-id (devnet) so scaleEVM runs its real path; with no EVM contracts
		// deployed in the test app it resolves nothing and is a graceful no-op. (scaleEVM
		// now HARD-FAILS on an unknown chain-id, so a fake id would make Redenominate error.)
		ChainID: "decimal_20202020-1",
		Time:    time.Now(),
	})
	base := cmdcfg.BaseDenom

	// Fund several accounts with assorted del amounts (mint via the coin module, which
	// holds Minter perm, then move to fresh accounts).
	amounts := []sdkmath.Int{
		sdkmath.NewInt(999),                 // dust -> 0
		sdkmath.NewInt(1000),                // -> 1
		sdkmath.NewInt(1500),                // -> 1
		sdkmath.NewInt(1_000_000),           // -> 1000
		sdkmath.NewIntWithDecimal(53_2, 18), // a "real" balance
	}
	var addrs []sdk.AccAddress
	for i, amt := range amounts {
		coins := sdk.NewCoins(sdk.NewCoin(base, amt))
		require.NoError(t, dscApp.BankKeeper.MintCoins(ctx, cointypes.ModuleName, coins))
		raw := make([]byte, 20) // 20-byte eth-style account address
		copy(raw, "redenomtestacct")
		raw[19] = byte(i)
		acc := sdk.AccAddress(raw)
		require.NoError(t, dscApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, cointypes.ModuleName, acc, coins))
		addrs = append(addrs, acc)
	}

	oldSupply := dscApp.BankKeeper.GetSupply(ctx, base).Amount

	keepers := redenom.Keepers{
		Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
		NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
		Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper, Fee: &dscApp.FeeKeeper,
	}
	storeKeys := redenom.StoreKeys{
		Bank: dscApp.GetKey(banktypes.StoreKey),
		EVM:  dscApp.GetKey(evmtypes.StoreKey),
	}

	div := sdkmath.NewInt(1000)
	report, err := redenom.Redenominate(ctx, keepers, storeKeys, div)
	require.NoError(t, err)

	// Every funded account is floor(old/1000).
	for i, acc := range addrs {
		got := dscApp.BankKeeper.GetBalance(ctx, acc, base).Amount
		want := amounts[i].Quo(div)
		require.Truef(t, got.Equal(want), "acct %d: got %s want %s", i, got, want)
	}

	// supply == report.NewSupplyDel == sum of all del balances, and <= floor(old).
	newSupply := dscApp.BankKeeper.GetSupply(ctx, base).Amount
	require.True(t, newSupply.Equal(report.NewSupplyDel), "supply %s != report %s", newSupply, report.NewSupplyDel)

	sum := sdkmath.ZeroInt()
	dscApp.BankKeeper.IterateAllBalances(ctx, func(_ sdk.AccAddress, c sdk.Coin) bool {
		if c.Denom == base {
			sum = sum.Add(c.Amount)
		}
		return false
	})
	require.True(t, newSupply.Equal(sum), "supply %s != sum(balances) %s", newSupply, sum)
	require.True(t, newSupply.LTE(oldSupply.Quo(div)), "new supply %s > floor(old) %s", newSupply, oldSupply.Quo(div))
	// RoundingDust = floor(old/div) - new >= 0 and is tiny relative to the reduction.
	require.False(t, report.RoundingDust.IsNegative(), "rounding dust negative: %s", report.RoundingDust)
	require.True(t, report.SupplyRemoved.Equal(oldSupply.Sub(newSupply)))
}

// TestRedenominate_BaseCoinRecord validates that the base coin's whole record is
// scaled: Volume -> new bank supply, Reserve -> floored, and (the E2 fix) LimitVolume
// -> floored. LimitVolume lives in the main coin record (not the CoinVR sub-record that
// UpdateCoinVR writes), and it is the denominator of the x/validator reward
// "percentForHold" split; leaving it ×1000 while stakes are scaled ÷1000 pins
// percentForHold to its 90% cap and diverts ~90% of every block's reward.
func TestRedenominate_BaseCoinRecord(t *testing.T) {
	dscApp := app.Setup(t, false, nil)
	ctx := dscApp.BaseApp.NewContext(false, tmproto.Header{
		Height:  1,
		ChainID: "decimal_20202020-1",
		Time:    time.Now(),
	})
	base := cmdcfg.BaseDenom

	// Fund one account so the base coin has a non-zero bank supply to reconcile against.
	bal := sdkmath.NewIntWithDecimal(1_000_000, 18)
	coins := sdk.NewCoins(sdk.NewCoin(base, bal))
	require.NoError(t, dscApp.BankKeeper.MintCoins(ctx, cointypes.ModuleName, coins))
	acc := sdk.AccAddress(append([]byte("redenombasecoin"), make([]byte, 5)...)[:20])
	require.NoError(t, dscApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, cointypes.ModuleName, acc, coins))

	// Establish a base coin record with a non-trivial Reserve and LimitVolume.
	reserve := sdkmath.NewIntWithDecimal(700_000, 18)
	limit := sdkmath.NewIntWithDecimal(108_000_000, 18) // emission cap, far above supply
	dscApp.CoinKeeper.SetCoin(ctx, cointypes.Coin{
		Denom:       base,
		Title:       "Decimal",
		CRR:         100,
		Reserve:     reserve,
		Volume:      bal,
		LimitVolume: limit,
		MinVolume:   sdkmath.ZeroInt(),
	})

	keepers := redenom.Keepers{
		Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
		NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
		Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper, Fee: &dscApp.FeeKeeper,
	}
	storeKeys := redenom.StoreKeys{
		Bank: dscApp.GetKey(banktypes.StoreKey),
		EVM:  dscApp.GetKey(evmtypes.StoreKey),
	}

	div := sdkmath.NewInt(1000)
	report, err := redenom.Redenominate(ctx, keepers, storeKeys, div)
	require.NoError(t, err)

	got, err := dscApp.CoinKeeper.GetCoin(ctx, base)
	require.NoError(t, err)

	// Volume tracks the new bank supply; Reserve and LimitVolume are floored by div.
	require.True(t, got.Volume.Equal(report.NewSupplyDel), "volume %s != new supply %s", got.Volume, report.NewSupplyDel)
	require.True(t, got.Reserve.Equal(reserve.Quo(div)), "reserve %s != floor %s", got.Reserve, reserve.Quo(div))
	require.True(t, got.LimitVolume.Equal(limit.Quo(div)), "LimitVolume %s != floor %s", got.LimitVolume, limit.Quo(div))
}

// TestRedenominate_StakeHoldAmounts is the R3 regression
// (redenom-full-audit-2026-07-02.md): every StakeHold.Amount inside base-denom
// delegation/undelegation/redelegation stakes must be floored together with the
// stake amount it is a part of — otherwise ≥1yr holds keep ×1000 reward weight
// (PayRewards weighs holds by hold.Amount) and auto-unbond computes
// scaledStake − Σ(×1000 holds) < 0, stranding the non-held portion. Custom-coin
// stakes and their holds stay untouched.
func TestRedenominate_StakeHoldAmounts(t *testing.T) {
	dscApp := app.Setup(t, false, nil)
	ctx := dscApp.BaseApp.NewContext(false, tmproto.Header{
		Height:  1,
		ChainID: "decimal_20202020-1",
		Time:    time.Now(),
	})
	base := cmdcfg.BaseDenom

	delegator := sdk.AccAddress(append([]byte("redenomholddeleg"), make([]byte, 4)...)[:20]).String()
	validator := sdk.ValAddress(append([]byte("redenomholdval"), make([]byte, 6)...)[:20]).String()
	validator2 := sdk.ValAddress(append([]byte("redenomholdval2"), make([]byte, 5)...)[:20]).String()

	holds := func(amts ...int64) []*validatortypes.StakeHold {
		var hs []*validatortypes.StakeHold
		for i, a := range amts {
			hs = append(hs, &validatortypes.StakeHold{
				Amount:        sdkmath.NewInt(a),
				HoldStartTime: 1000 + int64(i),
				HoldEndTime:   2000 + int64(i),
			})
		}
		return hs
	}
	coinStake := func(denom string, amt int64, hs []*validatortypes.StakeHold) validatortypes.Stake {
		return validatortypes.Stake{
			Type:  validatortypes.StakeType_Coin,
			ID:    denom,
			Stake: sdk.NewCoin(denom, sdkmath.NewInt(amt)),
			Holds: hs,
		}
	}

	// DEL delegation with two holds; custom-coin delegation with one hold.
	dscApp.ValidatorKeeper.SetDelegation(ctx, validatortypes.Delegation{
		Delegator: delegator, Validator: validator,
		Stake: coinStake(base, 5000, holds(3000, 1000)),
	})
	dscApp.ValidatorKeeper.SetDelegation(ctx, validatortypes.Delegation{
		Delegator: delegator, Validator: validator,
		Stake: coinStake("custcoin", 7000, holds(7000)),
	})
	// DEL undelegation and redelegation entries with holds.
	dscApp.ValidatorKeeper.SetUndelegation(ctx, validatortypes.Undelegation{
		Delegator: delegator, Validator: validator,
		Entries: []validatortypes.UndelegationEntry{{
			CreationHeight: 1, CompletionTime: time.Now().Add(time.Hour),
			Stake: coinStake(base, 4000, holds(2000)),
		}},
	})
	dscApp.ValidatorKeeper.SetRedelegation(ctx, validatortypes.Redelegation{
		Delegator: delegator, ValidatorSrc: validator, ValidatorDst: validator2,
		Entries: []validatortypes.RedelegationEntry{{
			CreationHeight: 1, CompletionTime: time.Now().Add(time.Hour),
			Stake: coinStake(base, 6000, holds(5000)),
		}},
	})

	keepers := redenom.Keepers{
		Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
		NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
		Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper, Fee: &dscApp.FeeKeeper,
	}
	storeKeys := redenom.StoreKeys{
		Bank: dscApp.GetKey(banktypes.StoreKey),
		EVM:  dscApp.GetKey(evmtypes.StoreKey),
	}

	div := sdkmath.NewInt(1000)
	_, err := redenom.Redenominate(ctx, keepers, storeKeys, div)
	require.NoError(t, err)

	requireHolds := func(st validatortypes.Stake, want ...int64) {
		require.Len(t, st.Holds, len(want))
		sum := sdkmath.ZeroInt()
		for i, w := range want {
			require.Truef(t, st.Holds[i].Amount.Equal(sdkmath.NewInt(w)),
				"hold %d of %s stake: got %s want %d", i, st.Stake.Denom, st.Holds[i].Amount, w)
			sum = sum.Add(st.Holds[i].Amount)
		}
		require.Truef(t, sum.LTE(st.Stake.Amount),
			"Σholds %s > stake %s for %s", sum, st.Stake.Amount, st.Stake.Denom)
	}

	seen := 0
	for _, d := range dscApp.ValidatorKeeper.GetAllDelegations(ctx) {
		if d.Delegator != delegator {
			continue // ignore genesis self-delegations from app.Setup
		}
		seen++
		switch d.Stake.Stake.Denom {
		case base:
			require.True(t, d.Stake.Stake.Amount.Equal(sdkmath.NewInt(5)))
			requireHolds(d.Stake, 3, 1)
		case "custcoin":
			require.True(t, d.Stake.Stake.Amount.Equal(sdkmath.NewInt(7000)), "custom stake must be untouched")
			requireHolds(d.Stake, 7000)
		}
	}
	require.Equal(t, 2, seen, "both test delegations must be found")

	dscApp.ValidatorKeeper.IterateUndelegations(ctx, func(_ int64, ubd validatortypes.Undelegation) bool {
		if ubd.Delegator != delegator {
			return false
		}
		require.True(t, ubd.Entries[0].Stake.Stake.Amount.Equal(sdkmath.NewInt(4)))
		requireHolds(ubd.Entries[0].Stake, 2)
		return false
	})
	dscApp.ValidatorKeeper.IterateRedelegations(ctx, func(_ int64, red validatortypes.Redelegation) bool {
		if red.Delegator != delegator {
			return false
		}
		require.True(t, red.Entries[0].Stake.Stake.Amount.Equal(sdkmath.NewInt(6)))
		requireHolds(red.Entries[0].Stake, 5)
		return false
	})
}

// TestRedenominate_RejectsBadDivisor guards the divisor precondition.
func TestRedenominate_RejectsBadDivisor(t *testing.T) {
	dscApp := app.Setup(t, false, nil)
	ctx := dscApp.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "redenom-test_1-1", Time: time.Now()})
	keepers := redenom.Keepers{
		Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
		NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
		Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper, Fee: &dscApp.FeeKeeper,
	}
	sk := redenom.StoreKeys{Bank: dscApp.GetKey(banktypes.StoreKey), EVM: dscApp.GetKey(evmtypes.StoreKey)}
	_, err := redenom.Redenominate(ctx, keepers, sk, sdkmath.OneInt())
	require.Error(t, err)
}

// TestRedenominate_FeeOraclePrice validates the fee-layer fix (fork rehearsal
// 2026-07-06 F-9/F-10): the base coin's fiat oracle price must be MULTIPLIED by the
// divisor (one new DEL is worth div× more), which in turn divides the derived EVM
// base fee / min gas price (x/fee GetMinGasPrice = EvmGasPrice / price) by div.
// Leaving the record unscaled keeps every EVM and Cosmos fee at div× its intended
// real value (measured live on the mainnet fork: 21k-gas transfer = 1.125 new DEL
// instead of 0.001125).
func TestRedenominate_FeeOraclePrice(t *testing.T) {
	dscApp := app.Setup(t, false, nil)
	ctx := dscApp.BaseApp.NewContext(false, tmproto.Header{
		Height: 1, ChainID: "decimal_20202020-1", Time: time.Now(),
	})
	base := cmdcfg.BaseDenom

	// The live mainnet value at the rehearsal snapshot: del/usd = 0.04.
	require.NoError(t, dscApp.FeeKeeper.SavePrice(ctx, feetypes.CoinPrice{
		Denom: base, Quote: "usd", Price: sdk.MustNewDecFromStr("0.04"),
	}))
	minGasBefore := dscApp.FeeKeeper.GetMinGasPrice(ctx)

	keepers := redenom.Keepers{
		Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
		NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
		Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper, Fee: &dscApp.FeeKeeper,
	}
	sk := redenom.StoreKeys{Bank: dscApp.GetKey(banktypes.StoreKey), EVM: dscApp.GetKey(evmtypes.StoreKey)}
	div := sdkmath.NewInt(1000)
	report, err := redenom.Redenominate(ctx, keepers, sk, div)
	require.NoError(t, err)
	require.GreaterOrEqual(t, report.FeePricesScaled, 1)

	price, err := dscApp.FeeKeeper.GetPrice(ctx, base, "usd")
	require.NoError(t, err)
	require.True(t, price.Price.Equal(sdk.MustNewDecFromStr("40")), "price %s != 40", price.Price)

	// The derived EVM min gas price (== base fee) drops by exactly div.
	minGasAfter := dscApp.FeeKeeper.GetMinGasPrice(ctx)
	require.True(t, minGasAfter.MulInt(div).Equal(minGasBefore),
		"minGasPrice %s != before %s ÷ %s", minGasAfter, minGasBefore, div)

	// A missing fee keeper must hard-fail, never silently skip the fee layer.
	badKeepers := keepers
	badKeepers.Fee = nil
	_, err = redenom.Redenominate(ctx, badKeepers, sk, div)
	require.ErrorContains(t, err, "fee keeper not wired")
}
