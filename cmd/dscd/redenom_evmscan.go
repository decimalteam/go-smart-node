package main

import (
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
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/app/redenom"
	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
)

// RedenomEVMScanCmd is the decoder-independent EVM verification. EVM storage slots are
// opaque 32-byte words with no denom marker, so unlike Cosmos there is no "del" byte
// signature to scan for. Instead it treats EVERY storage slot of EVERY contract as an
// individual record: it snapshots the entire EVM storage namespace, runs the real
// redenomination in a never-committed cache, re-snapshots, and classifies every slot that
// changed as a clean ÷1000 floor, a floored-to-zero/void, or — the thing that must never
// happen — an "other" change (corruption). It does NOT use stakescan to decide which slots
// "are DEL", so it is independent of the upgrade's own identification logic.
//
// To catch DEL the upgrade MISSED (a slot that holds a DEL amount but was left unscaled),
// which a raw word can't reveal on its own, it anchors on the native "del" bank balance:
// it lists every contract that (a) holds native del, (b) has EVM storage, and (c) had ZERO
// slots changed by the upgrade — i.e. its del backing was divided but its internal ledger
// was never touched. System contracts (delegation/WDEL/collections/checks) change many
// slots and drop out; what remains are candidate missed DEL ledgers for manual review.
func RedenomEVMScanCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redenom-evmscan",
		Short: "Decoder-independent full EVM-storage diff across the upgrade (every slot; writes nothing)",
		Long: `Snapshots every slot of every EVM contract, runs the redenomination in a cache-only
context, re-snapshots, and verifies every changed slot is a clean ÷1000 floor or a void —
flagging any other change as corruption. Also lists DEL-holding contracts whose storage the
upgrade never modified (candidate missed DEL ledgers). The node must be stopped; nothing is committed.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			homePath := cast.ToString(appOpts.Get(flags.FlagHome))
			if homePath == "" {
				return fmt.Errorf("--home is required")
			}
			topN, _ := cmd.Flags().GetInt("top")
			dumpAddrs, _ := cmd.Flags().GetStringSlice("dump")
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
			evmKey := dscApp.GetKey(evmtypes.StoreKey)

			out := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(out, "DEL redenomination EVM full-storage diff\tchain-id %s\theight %d\tdivisor %s\n\n", chainID, height, div)

			// Optional: dump raw storage of specific contracts (to characterize "other" del holders).
			if len(dumpAddrs) > 0 {
				for _, a := range dumpAddrs {
					addr := common.HexToAddress(a)
					nat := dscApp.BankKeeper.GetBalance(ctx, sdk.AccAddress(addr.Bytes()), base).Amount.String()
					st := scanStorage(ctx.KVStore(evmKey), addr)
					fmt.Fprintf(out, "\n== %s  native del=%s  slots=%d ==\n", addr.Hex(), nat, len(st))
					var keys []string
					for s := range st {
						keys = append(keys, s.Hex())
					}
					sort.Strings(keys)
					for _, ks := range keys {
						s := common.HexToHash(ks)
						v := new(big.Int).SetBytes(st[s].Bytes())
						fmt.Fprintf(out, "  %s = %s\n", ks, v.String())
					}
				}
				out.Flush()
				return nil
			}

			// --- contracts that hold native del + which addresses have storage (pre-scale) ---
			delNative := map[common.Address]sdkmath.Int{}
			dscApp.BankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, c sdk.Coin) bool {
				if c.Denom == base && len(addr.Bytes()) == 20 {
					delNative[common.BytesToAddress(addr.Bytes())] = c.Amount
				}
				return false
			})

			// --- snapshot ALL evm storage slots BEFORE ---
			before := map[string][]byte{}
			storageAddrs := map[common.Address]struct{}{}
			scanEVMStorage(ctx, evmKey, func(key, val []byte) {
				k := string(key)
				v := make([]byte, len(val))
				copy(v, val)
				before[k] = v
				storageAddrs[common.BytesToAddress(key[:20])] = struct{}{}
			})
			fmt.Fprintf(out, "snapshot BEFORE: %d storage slots across %d contracts\n", len(before), len(storageAddrs))

			// --- run the real upgrade ---
			keepers := redenom.Keepers{
				Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
				NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
				Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper, Fee: &dscApp.FeeKeeper,
			}
			storeKeys := redenom.StoreKeys{Bank: dscApp.GetKey(banktypes.StoreKey), EVM: evmKey}
			if _, err := redenom.Redenominate(ctx, keepers, storeKeys, div); err != nil {
				return fmt.Errorf("redenominate: %w", err)
			}

			// --- diff: walk AFTER, compare to BEFORE (popping seen keys) ---
			dv := div.BigInt()
			type cAgg struct {
				floored, flooredZero, voided, other, changed int
				firstOther                                   string
			}
			per := map[common.Address]*cAgg{}
			ca := func(a common.Address) *cAgg {
				if per[a] == nil {
					per[a] = &cAgg{}
				}
				return per[a]
			}
			var unchanged, created int
			scanEVMStorage(ctx, evmKey, func(key, val []byte) {
				k := string(key)
				addr := common.BytesToAddress(key[:20])
				newBI := new(big.Int).SetBytes(val)
				ob, ok := before[k]
				if ok {
					delete(before, k) // mark seen
					oldBI := new(big.Int).SetBytes(ob)
					if oldBI.Cmp(newBI) == 0 {
						unchanged++
						return
					}
					g := ca(addr)
					g.changed++
					want := new(big.Int).Quo(oldBI, dv)
					if newBI.Cmp(want) == 0 && newBI.Sign() != 0 {
						g.floored++
					} else {
						g.other++
						if g.firstOther == "" {
							g.firstOther = fmt.Sprintf("slot %x old=%s new=%s want=%s", key[20:], oldBI, newBI, want)
						}
					}
					return
				}
				// present after, absent before => created from zero (no redenom path does this).
				if newBI.Sign() != 0 {
					created++
					g := ca(addr)
					g.changed++
					g.other++
					if g.firstOther == "" {
						g.firstOther = fmt.Sprintf("slot %x CREATED new=%s", key[20:], newBI)
					}
				}
			})
			// leftover BEFORE keys = slots deleted by the upgrade (new == 0).
			for k, ob := range before {
				addr := common.BytesToAddress([]byte(k)[:20])
				oldBI := new(big.Int).SetBytes(ob)
				if oldBI.Sign() == 0 {
					continue
				}
				g := ca(addr)
				g.changed++
				if new(big.Int).Quo(oldBI, dv).Sign() == 0 {
					g.flooredZero++ // small DEL amount floored to 0 -> slot deleted (legit)
				} else {
					g.voided++ // large value zeroed without flooring (expected only for checks void)
				}
			}

			// --- totals ---
			var tFloored, tFlooredZero, tVoided, tOther, tChanged int
			for _, g := range per {
				tFloored += g.floored
				tFlooredZero += g.flooredZero
				tVoided += g.voided
				tOther += g.other
				tChanged += g.changed
			}
			fmt.Fprintf(out, "\n-- global slot diff --\n")
			fmt.Fprintf(out, "unchanged\t%d\n", unchanged)
			fmt.Fprintf(out, "changed\t%d\n", tChanged)
			fmt.Fprintf(out, "  ÷1000 floored\t%d\n", tFloored)
			fmt.Fprintf(out, "  floored-to-zero (deleted)\t%d\n", tFlooredZero)
			fmt.Fprintf(out, "  voided (large->0; expect checks only)\t%d\n", tVoided)
			fmt.Fprintf(out, "  OTHER / corruption (must be 0)\t%d\n", tOther)
			fmt.Fprintf(out, "created-from-zero\t%d\n", created)

			// --- per-contract changed summary (top N) ---
			type row struct {
				addr common.Address
				g    *cAgg
			}
			var rows []row
			for a, g := range per {
				rows = append(rows, row{a, g})
			}
			sort.Slice(rows, func(i, j int) bool { return rows[i].g.changed > rows[j].g.changed })
			fmt.Fprintf(out, "\n-- contracts modified by the upgrade (top %d of %d) --\n", min(topN, len(rows)), len(rows))
			fmt.Fprintf(out, "contract\tnative del (post)\tchanged\tfloored\tflZero\tvoided\tother\n")
			for i, r := range rows {
				if i >= topN {
					break
				}
				nat := "-"
				if v, ok := delNative[r.addr]; ok {
					nat = floorI(v, div).String()
				}
				fmt.Fprintf(out, "%s\t%s\t%d\t%d\t%d\t%d\t%d\n",
					r.addr.Hex(), nat, r.g.changed, r.g.floored, r.g.flooredZero, r.g.voided, r.g.other)
			}
			if tOther > 0 {
				fmt.Fprintf(out, "\n!! OTHER changes detected (potential corruption):\n")
				n := 0
				for _, r := range rows {
					if r.g.other > 0 {
						fmt.Fprintf(out, "  %s: %d other; %s\n", r.addr.Hex(), r.g.other, r.g.firstOther)
						if n++; n >= 20 {
							break
						}
					}
				}
			}

			// --- independent missed-ledger detection: DEL-holding contracts not modified ---
			fmt.Fprintf(out, "\n-- DEL-holding contracts whose storage the upgrade did NOT modify --\n")
			fmt.Fprintf(out, "(EOAs have no storage and are excluded; system contracts change slots and drop out)\n")
			type miss struct {
				addr common.Address
				old  sdkmath.Int
			}
			var misses []miss
			for a, nat := range delNative {
				if _, hasStorage := storageAddrs[a]; !hasStorage {
					continue // EOA / no contract storage -> its del is just the (scaled) bank balance
				}
				if g, ok := per[a]; ok && g.changed > 0 {
					continue // upgrade touched this contract's storage
				}
				if !floorI(nat, div).IsPositive() {
					continue // native floored to 0 (already-known dropped-collection dust case)
				}
				misses = append(misses, miss{a, nat})
			}
			sort.Slice(misses, func(i, j int) bool { return misses[i].old.GT(misses[j].old) })
			if len(misses) == 0 {
				fmt.Fprintf(out, "none — every DEL-holding contract had its storage modified (or is an EOA / dust).\n")
				out.Flush()
				return nil
			}

			// Classify each missed contract independently. A DEL-custody contract that should
			// have been scaled is an NFTReserve-style contract holding a DEL `_reserve` entry
			// (reserveType slot == DEL(1), its token slot == address(0), amount slot != 0).
			// Custom-coin DRC20s are token-center beacon proxies whose del is the (cosmos-scaled)
			// coin reserve and whose ledger is the custom token — correctly left untouched.
			cc, _ := redenom.ContractCenterFor(chainID)
			ccStore := scanStorage(ctx.KVStore(evmKey), cc)
			tokenCenter := stakescan.ResolveAddressBySymbol(ccStore, "token-center")
			beaconSlot := common.HexToHash("0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50")
			delType := big.NewInt(1)

			type cls struct {
				kind        string // "drc20-coin" | "nft-collection" | "other"
				slots       int
				delReserves int // # of unscaled DEL _reserve entries (the concerning signal)
			}
			var coinProxies, nftWithReserve, otherC, withUnscaledReserve int
			classOf := map[common.Address]cls{}
			for _, m := range misses {
				st := scanStorage(ctx.KVStore(evmKey), m.addr)
				c := cls{kind: "other", slots: len(st)}
				beacon := st[beaconSlot]
				if beacon != (common.Hash{}) {
					if common.BytesToAddress(beacon.Bytes()) == tokenCenter {
						c.kind = "drc20-coin"
					} else {
						c.kind = "nft-collection"
					}
				}
				for slotHash, v := range st {
					if new(big.Int).SetBytes(v.Bytes()).Cmp(delType) != 0 {
						continue
					}
					sBig := new(big.Int).SetBytes(slotHash.Bytes())
					amountSlot := common.BigToHash(new(big.Int).Sub(sBig, big.NewInt(1)))
					tokenSlot := common.BigToHash(new(big.Int).Sub(sBig, big.NewInt(2)))
					if st[tokenSlot] != (common.Hash{}) {
						continue
					}
					if new(big.Int).SetBytes(st[amountSlot].Bytes()).Sign() != 0 {
						c.delReserves++
					}
				}
				switch c.kind {
				case "drc20-coin":
					coinProxies++
				case "nft-collection":
					nftWithReserve++
				default:
					otherC++
				}
				if c.delReserves > 0 {
					withUnscaledReserve++
				}
				classOf[m.addr] = c
			}

			fmt.Fprintf(out, "%d found. Classification:\n", len(misses))
			fmt.Fprintf(out, "  custom-coin DRC20 proxies (del = cosmos-scaled coin reserve; correctly untouched)\t%d\n", coinProxies)
			fmt.Fprintf(out, "  NFT-collection beacon proxies\t%d\n", nftWithReserve)
			fmt.Fprintf(out, "  other contracts\t%d\n", otherC)
			fmt.Fprintf(out, "  >> holding an UNSCALED DEL _reserve entry (potential MISS)\t%d\n", withUnscaledReserve)
			fmt.Fprintf(out, "\ntop %d by native del:\n", topN)
			fmt.Fprintf(out, "contract\tnative del (pre, base units)\tkind\tslots\tunscaledDelReserves\n")
			for i, m := range misses {
				if i >= topN {
					fmt.Fprintf(out, "... (%d more)\n", len(misses)-topN)
					break
				}
				c := classOf[m.addr]
				fmt.Fprintf(out, "%s\t%s\t%s\t%d\t%d\n", m.addr.Hex(), m.old.String(), c.kind, c.slots, c.delReserves)
			}
			if withUnscaledReserve > 0 {
				fmt.Fprintf(out, "\n!! contracts with an UNSCALED DEL _reserve (review — should the upgrade have scaled these?):\n")
				n := 0
				for _, m := range misses {
					c := classOf[m.addr]
					if c.delReserves == 0 {
						continue
					}
					fmt.Fprintf(out, "  %s  kind=%s  delReserves=%d  native=%s\n", m.addr.Hex(), c.kind, c.delReserves, m.old.String())
					if n++; n >= 40 {
						break
					}
				}
			}
			out.Flush()
			return nil
		},
	}
	cmd.Flags().Int("top", 40, "rows to print for the modified-contracts and missed-ledger tables")
	cmd.Flags().StringSlice("dump", nil, "comma-separated contract addresses: dump their raw storage and exit")
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory")
	return cmd
}

// scanEVMStorage calls fn(key, value) for every EVM storage slot, where key is addr(20)||slot(32).
func scanEVMStorage(ctx sdk.Context, evmKey storetypes.StoreKey, fn func(key, val []byte)) {
	ps := prefix.NewStore(ctx.KVStore(evmKey), evmtypes.KeyPrefixStorage)
	it := ps.Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		key := it.Key()
		if len(key) != 20+32 {
			continue
		}
		fn(key, it.Value())
	}
}

