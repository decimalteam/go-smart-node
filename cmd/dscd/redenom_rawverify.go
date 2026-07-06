package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"
	"text/tabwriter"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	legacytypes "bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	multisigtypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	swaptypes "bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/app/redenom"
)

// RedenomRawVerifyCmd is the DECODER-INDEPENDENT verification: instead of asking the
// redenomination's own logic where DEL lives, it walks every key/value in every module
// store and finds DEL amounts by their on-disk proto encoding. A base-coin amount is
// serialized as a Cosmos Coin {string denom = 1; string amount = 2;}, i.e. the byte
// sequence 0x0a 0x03 'd' 'e' 'l' 0x12 <varint len> <ascii-decimal-digits> — wherever it
// appears, in any record of any module. It records every such amount before the upgrade,
// runs the real Redenominate in a never-committed cache store, re-scans, and reports per
// store whether each DEL amount was scaled (÷1000), left unchanged, or changed otherwise.
//
// This catches DEL stored in places the upgrade's identification never looks (gov proposal
// deposits, multisig pending transactions, vesting accounts, fee/swap state, ...). Bank
// balances/supply (bare Int, denom in the key) and EVM contract slots (opaque binary) are
// not Coin-encoded and are counted only here — they are verified per-item by redenom-verify.
func RedenomRawVerifyCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redenom-rawverify",
		Short: "Decoder-independent raw-DB scan for DEL amounts; verify each is ÷1000 after the upgrade (writes nothing)",
		Long: `Walks every record in every module store, finds DEL amounts by their on-disk Coin proto
signature (independent of the redenomination's own identification logic), runs the full
redenomination in a cache-only context, and reports per store how many DEL amounts were
scaled, left unchanged, or changed otherwise — surfacing any DEL the upgrade does not touch.
The node must be stopped; nothing is committed.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			homePath := cast.ToString(appOpts.Get(flags.FlagHome))
			if homePath == "" {
				return fmt.Errorf("--home is required")
			}
			db, err := openRedenomAppDB(homePath, appOpts)
			if err != nil {
				return fmt.Errorf("open application db: %w", err)
			}
			defer db.Close()

			dscApp := app.NewDSC(serverCtx.Logger, db, nil, true, map[int64]bool{}, homePath, uint(1), encCfg, appOpts)
			height := dscApp.LastBlockHeight()
			chainID := readChainID(homePath)
			cms := dscApp.CommitMultiStore().CacheMultiStore()
			ctx := sdk.NewContext(cms, tmproto.Header{Height: height, ChainID: chainID}, false, serverCtx.Logger)

			base := dscApp.ValidatorKeeper.BaseDenom(ctx)
			div := redenom.DefaultDivisor
			out := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(out, "DEL redenomination RAW (decoder-independent) verify\tchain-id %s\theight %d\tbase %q\tdivisor %s\n\n",
				chainID, height, base, div)

			// Stores whose values are protobuf messages that may embed a DEL Coin. Each is
			// classified as scaled-by-redenom or untouched-by-redenom for the verdict.
			scanStores := []rawStore{
				{validatortypes.StoreKey, true, "delegations/undelegations/redelegations stakes"},
				{nfttypes.StoreKey, true, "token + sub-token reserves"},
				{legacytypes.StoreKey, true, "legacy record coins"},
				{govtypes.StoreKey, false, "proposal/deposit coins + MinDeposit param"},
				{authtypes.StoreKey, false, "accounts (incl. vesting schedules)"},
				{multisigtypes.StoreKey, false, "wallets + pending transactions"},
				{swaptypes.StoreKey, false, "swap state"},
				{feetypes.StoreKey, false, "fee module state"},
				{paramstypes.StoreKey, false, "module params subspace"},
				{upgradetypes.StoreKey, false, "upgrade plan"},
				{banktypes.StoreKey, false, "bank (Coin-encoded values, if any)"},
			}
			// Visited count-only (not Coin-encoded): bank balances/supply are bare Int with the
			// denom in the key; coin module uses separate Int reserve/volume; evm is opaque slots.
			countStores := []string{cointypes.StoreKey, evmtypes.StoreKey, capabilitytypes.StoreKey}

			// 1. BEFORE.
			before, statsBefore := scanAllStores(ctx, dscApp, scanStores)
			coverage := visitCount(ctx, dscApp, countStores)

			// 2. Run the real upgrade transform.
			keepers := redenom.Keepers{
				Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
				NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
				Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper, Fee: &dscApp.FeeKeeper,
			}
			storeKeys := redenom.StoreKeys{Bank: dscApp.GetKey(banktypes.StoreKey), EVM: dscApp.GetKey(evmtypes.StoreKey)}
			report, err := redenom.Redenominate(ctx, keepers, storeKeys, div)
			if err != nil {
				return fmt.Errorf("redenominate: %w", err)
			}
			_ = report

			// 3. AFTER (same cache ctx, post-mutation).
			after, _ := scanAllStores(ctx, dscApp, scanStores)

			// 4. Compare per store.
			reportRaw(out, scanStores, statsBefore, before, after, div, coverage)

			// 5. Multisig lifecycle breakdown: which unscaled DEL transactions are still
			// PENDING (signable to threshold -> would execute at the old x1000 amount).
			analyzeMultisig(out, dscApp, ctx)
			out.Flush()
			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory")
	return cmd
}

type rawStore struct {
	name   string
	scaled bool // does the redenomination scale DEL Coins in this store?
	desc   string
}

type storeStat struct {
	keys     int
	bytes    int64
	delCoins int
}

var delCoinSig = []byte{0x0a, 0x03, 'd', 'e', 'l', 0x12} // Coin{denom:"del", amount:...}

// scanDelCoinAmounts finds every DEL Coin amount embedded in a protobuf value by its wire
// signature, in field order. Returns the amounts (as big.Int) in the order found.
func scanDelCoinAmounts(val []byte) []*big.Int {
	var out []*big.Int
	for i := 0; i+len(delCoinSig) <= len(val); i++ {
		if !bytes.Equal(val[i:i+len(delCoinSig)], delCoinSig) {
			continue
		}
		j := i + len(delCoinSig)
		l, n := binary.Uvarint(val[j:])
		if n <= 0 || l == 0 || l > 40 { // amounts are short ascii-decimal strings
			continue
		}
		j += n
		if j+int(l) > len(val) {
			continue
		}
		digits := val[j : j+int(l)]
		if !allASCIIDigits(digits) {
			continue // e.g. a "del"-valued string field followed by a non-amount string
		}
		bi, ok := new(big.Int).SetString(string(digits), 10)
		if !ok {
			continue
		}
		out = append(out, bi)
	}
	return out
}

func allASCIIDigits(b []byte) bool {
	for _, c := range b {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(b) > 0
}

// scanAllStores iterates every key/value of each scan store, extracting DEL Coin amounts
// keyed by a stable locator (store|keyhex|ordinal). The ordinal is the position of the DEL
// Coin within the value in field order — stable across the upgrade (amounts change, count
// and order do not), so before/after locators line up.
func scanAllStores(ctx sdk.Context, a *app.DSC, stores []rawStore) (map[string]*big.Int, map[string]*storeStat) {
	m := map[string]*big.Int{}
	stats := map[string]*storeStat{}
	for _, s := range stores {
		st := &storeStat{}
		stats[s.name] = st
		kv := ctx.KVStore(a.GetKey(s.name))
		it := kv.Iterator(nil, nil)
		for ; it.Valid(); it.Next() {
			st.keys++
			v := it.Value()
			st.bytes += int64(len(v))
			amts := scanDelCoinAmounts(v)
			if len(amts) == 0 {
				continue
			}
			keyHex := fmt.Sprintf("%x", it.Key())
			for ord, amt := range amts {
				m[fmt.Sprintf("%s|%s|%d", s.name, keyHex, ord)] = amt
				st.delCoins++
			}
		}
		it.Close()
	}
	return m, stats
}

func visitCount(ctx sdk.Context, a *app.DSC, names []string) map[string]*storeStat {
	out := map[string]*storeStat{}
	for _, n := range names {
		st := &storeStat{}
		out[n] = st
		kv := ctx.KVStore(a.GetKey(n))
		it := kv.Iterator(nil, nil)
		for ; it.Valid(); it.Next() {
			st.keys++
			st.bytes += int64(len(it.Value()))
		}
		it.Close()
	}
	return out
}

func reportRaw(out *tabwriter.Writer, stores []rawStore, statsBefore map[string]*storeStat,
	before, after map[string]*big.Int, div sdkmath.Int, coverage map[string]*storeStat) {

	dv := div.BigInt()
	type agg struct {
		total, scaled, unchanged, other, vanished int
		unchangedDel                              *big.Int // sum of old amounts left unchanged
		firstOther                                string
	}
	byStore := map[string]*agg{}
	get := func(n string) *agg {
		if byStore[n] == nil {
			byStore[n] = &agg{unchangedDel: big.NewInt(0)}
		}
		return byStore[n]
	}

	for loc, old := range before {
		store := loc[:indexByte(loc, '|')]
		g := get(store)
		g.total++
		nv, ok := after[loc]
		if !ok {
			g.vanished++
			continue
		}
		want := new(big.Int).Quo(old, dv)
		switch {
		case nv.Cmp(old) == 0:
			g.unchanged++
			g.unchangedDel.Add(g.unchangedDel, old)
		case nv.Cmp(want) == 0:
			g.scaled++
		default:
			g.other++
			if g.firstOther == "" {
				g.firstOther = fmt.Sprintf("%s old=%s new=%s", loc, old, nv)
			}
		}
	}
	// Coins that APPEARED after (none expected for Coin amounts; refunds are bank-side Ints).
	appeared := map[string]int{}
	for loc := range after {
		if _, ok := before[loc]; !ok {
			appeared[loc[:indexByte(loc, '|')]]++
		}
	}

	fmt.Fprintf(out, "-- per-store DEL Coin scan (decoder-independent) --\n")
	fmt.Fprintf(out, "store\tscaled?\tkeys\tbytes\tdelCoins\t÷1000\tunchanged\tother\tverdict\n")
	var names []string
	for _, s := range stores {
		names = append(names, s.name)
	}
	scaledOf := map[string]bool{}
	descOf := map[string]string{}
	for _, s := range stores {
		scaledOf[s.name] = s.scaled
		descOf[s.name] = s.desc
	}
	totalFindings := 0
	for _, n := range names {
		g := get(n)
		sb := statsBefore[n]
		verdict := "ok"
		if scaledOf[n] {
			// Everything here should have scaled. Unchanged with a non-zero old amount that
			// does not floor to itself (i.e. old >= div) is a genuine MISS.
			if g.unchanged > 0 || g.other > 0 || g.vanished > 0 {
				verdict = fmt.Sprintf("CHECK: %d unchanged, %d other", g.unchanged, g.other)
			}
		} else {
			// redenom does not scale these. Any DEL Coin present is an unscaled DEL value.
			if g.total > 0 {
				verdict = fmt.Sprintf("GAP?: %d DEL amounts NOT scaled by redenom", g.total)
				totalFindings += g.total
			}
		}
		if g.other > 0 && g.firstOther != "" {
			verdict += " | e.g. " + g.firstOther
		}
		fmt.Fprintf(out, "%s\t%v\t%d\t%d\t%d\t%d\t%d\t%d\t%s\n",
			n, scaledOf[n], sb.keys, sb.bytes, sb.delCoins, g.scaled, g.unchanged, g.other, verdict)
	}

	// Coverage note for the non-Coin stores.
	fmt.Fprintf(out, "\n-- visited (not Coin-encoded; verified per-item by redenom-verify) --\n")
	var cn []string
	for n := range coverage {
		cn = append(cn, n)
	}
	sort.Strings(cn)
	for _, n := range cn {
		fmt.Fprintf(out, "%s\tkeys=%d\tbytes=%d\n", n, coverage[n].keys, coverage[n].bytes)
	}

	// Detail any store that redenom does NOT scale but where DEL amounts live.
	fmt.Fprintf(out, "\n-- findings: DEL amounts in stores the upgrade does not scale --\n")
	any := false
	for _, n := range names {
		if scaledOf[n] {
			continue
		}
		g := get(n)
		if g.total == 0 {
			continue
		}
		any = true
		fmt.Fprintf(out, "%s (%s): %d DEL Coin amount(s) present; %d unchanged (%s base-unit DEL total left unscaled)\n",
			n, descOf[n], g.total, g.unchanged, g.unchangedDel.String())
	}
	if !any {
		fmt.Fprintf(out, "none — every DEL Coin in the database lives in a store the upgrade scales.\n")
	}
	for n, c := range appeared {
		fmt.Fprintf(out, "NOTE: %d DEL Coin amount(s) APPEARED in %s after the upgrade (verify expected)\n", c, n)
	}
	fmt.Fprintf(out, "\nTotal DEL Coin amounts in unscaled stores: %d\n", totalFindings)
}

// analyzeMultisig classifies every multisig transaction carrying a DEL amount as PENDING
// (not yet completed) or completed, and sums the DEL at risk. A PENDING transaction created
// before the upgrade executes its STORED message (original x1000 amount) when it later
// reaches the signature threshold — the redenomination never rewrites it. (Run on the
// post-upgrade cache ctx; completion flags and stored amounts are unaffected by scaling.)
func analyzeMultisig(out *tabwriter.Writer, a *app.DSC, ctx sdk.Context) {
	txs, err := a.MultisigKeeper.GetAllTransactions(ctx)
	if err != nil {
		fmt.Fprintf(out, "\n-- multisig analysis: error: %v\n", err)
		return
	}
	pendingN, completedN := 0, 0
	pendingDel, completedDel := big.NewInt(0), big.NewInt(0)
	var examples []string
	for _, tx := range txs {
		amts := scanDelCoinAmounts(tx.Message.Value)
		if len(amts) == 0 {
			continue
		}
		sum := big.NewInt(0)
		for _, x := range amts {
			sum.Add(sum, x)
		}
		completed := a.MultisigKeeper.IsCompleted(ctx, tx.Id)
		if completed {
			completedN++
			completedDel.Add(completedDel, sum)
		} else {
			pendingN++
			pendingDel.Add(pendingDel, sum)
			if len(examples) < 8 {
				examples = append(examples, fmt.Sprintf("  tx %s wallet %s del=%s", tx.Id, tx.Wallet, sum))
			}
		}
	}
	fmt.Fprintf(out, "\n-- multisig DEL transaction lifecycle (the unscaled-DEL finding) --\n")
	fmt.Fprintf(out, "PENDING (signable -> would execute at OLD x1000 amount): %d tx, %s base-unit DEL at risk\n", pendingN, pendingDel.String())
	fmt.Fprintf(out, "completed (already executed; harmless):                   %d tx, %s base-unit DEL\n", completedN, completedDel.String())
	if len(examples) > 0 {
		fmt.Fprintf(out, "pending examples:\n")
		for _, e := range examples {
			fmt.Fprintf(out, "%s\n", e)
		}
	}
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return len(s)
}
