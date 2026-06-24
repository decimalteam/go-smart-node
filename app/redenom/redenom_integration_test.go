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
		Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper,
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

// TestRedenominate_RejectsBadDivisor guards the divisor precondition.
func TestRedenominate_RejectsBadDivisor(t *testing.T) {
	dscApp := app.Setup(t, false, nil)
	ctx := dscApp.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "redenom-test_1-1", Time: time.Now()})
	keepers := redenom.Keepers{
		Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
		NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
		Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper,
	}
	sk := redenom.StoreKeys{Bank: dscApp.GetKey(banktypes.StoreKey), EVM: dscApp.GetKey(evmtypes.StoreKey)}
	_, err := redenom.Redenominate(ctx, keepers, sk, sdkmath.OneInt())
	require.Error(t, err)
}
