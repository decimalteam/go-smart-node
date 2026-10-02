package main

import (
	"path/filepath"

	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"fmt"
	cometbftdb "github.com/cometbft/cometbft-db"
	cmtstore "github.com/tendermint/tendermint/store"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
)

// DesyncPlanCmd runs the mainnet delegation-desync repair against an offline node
// home in a cache-only context: it prints the plan, applies it, re-plans to prove
// the node then matches the contract, and optionally writes every diff as JSON.
// Nothing is committed; the node must be stopped.
func DesyncPlanCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "desync-plan",
		Short: "Dry-run the delegation-desync repair on an offline node home (writes nothing)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			homePath := cast.ToString(appOpts.Get(flags.FlagHome))
			if homePath == "" {
				return fmt.Errorf("--home is required")
			}
			unix, _ := cmd.Flags().GetInt64("time")
			outPath, _ := cmd.Flags().GetString("out")
			apply, _ := cmd.Flags().GetBool("apply")

			db, err := openRedenomAppDB(homePath, appOpts)
			if err != nil {
				return fmt.Errorf("open application db: %w", err)
			}
			defer db.Close()

			dscApp := app.NewDSC(serverCtx.Logger, db, nil, true, map[int64]bool{}, homePath, uint(1), encCfg, appOpts)
			height := dscApp.LastBlockHeight()
			chainID := readChainID(homePath)
			if unix == 0 {
				// Hold maturity follows block time, which on a fast fork can run ahead
				// of the wall clock; use the last committed block's time.
				unix = time.Now().Unix()
				if blockDB, err := cometbftdb.NewGoLevelDB("blockstore", filepath.Join(homePath, "data")); err == nil {
					bs := cmtstore.NewBlockStore(blockDB)
					if meta := bs.LoadBlockMeta(bs.Height()); meta != nil {
						unix = meta.Header.Time.Unix()
					}
					blockDB.Close()
				}
			}
			cms := dscApp.CommitMultiStore().CacheMultiStore()
			ctx := sdk.NewContext(cms, tmproto.Header{Height: height + 1, ChainID: chainID, Time: time.Unix(unix, 0).UTC()}, false, serverCtx.Logger)
			// EVM calls resolve the coinbase from the block proposer.
			for _, v := range dscApp.ValidatorKeeper.GetAllValidators(ctx) {
				if v.IsBonded() {
					cons, err := v.GetConsAddr()
					if err == nil {
						ctx = ctx.WithBlockHeader(tmproto.Header{Height: height + 1, ChainID: chainID, Time: time.Unix(unix, 0).UTC(), ProposerAddress: cons})
						break
					}
				}
			}

			out := cmd.OutOrStdout()
			plan, err := app.PlanDelegationRepair(ctx, dscApp)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "height %d time %d\nplan: %s\n", height, unix, plan.Stats.String())
			migrated, merr := app.DelegationContractMigrated(ctx, dscApp)
			fmt.Fprintf(out, "contract upgraded+migrated (hold-start backfill will run): %t %v\n", migrated, merr)
			{
				v := app.ValuePlan(ctx, dscApp, plan)
				del := func(x sdkmath.Int) string { return sdkmath.LegacyNewDecFromIntWithPrec(x, 18).String() }
				fmt.Fprintf(out, "value in base coin: created +%s increased +%s decreased -%s removed -%s\n", del(v.Created), del(v.Increased), del(v.Decreased), del(v.Removed))
			}
			if qf, _ := cmd.Flags().GetInt64("queue-from"); qf >= 0 {
				n, total, err := app.AutoUnbondQueueValue(ctx, dscApp, uint64(qf))
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "auto-unbond queue from #%d: %d entries, %s base coin\n", qf, n, sdkmath.LegacyNewDecFromIntWithPrec(total, 18).String())
			}
			for tok, n := range plan.Stats.UnknownTokenCounts() {
				fmt.Fprintf(out, "  unknown token %s: %d stakes\n", tok.Hex(), n)
			}
			{
				touched := 0
				for _, d := range plan.Diffs {
					for k := range plan.ZeroTokenQueued {
						if k.Validator() == d.Target.Validator.String() && k.Delegator() == d.Target.Delegator.String() {
							touched++
							fmt.Fprintf(out, "  zero-token queue overlaps diff: %s %s %q node=%s found=%t target=%s\n", d.Target.Validator, d.Target.Delegator, d.Target.Denom, d.NodeAmount, d.NodeFound, d.Target.Amount)
						}
					}
				}
				fmt.Fprintf(out, "zero-token queue pairs=%d overlapping diffs=%d\n", len(plan.ZeroTokenQueued), touched)
				if os.Getenv("DESYNC_ZERO_TOKEN_DETAIL") != "" {
					// For every 0x0 entry: which coin stake of the same (validator, delegator)
					// has exactly that amount, and was that coin missing its DRC20 address?
					noDRC := map[string]bool{}
					for _, d := range plan.ResolvedBySymbol {
						noDRC[d] = true
					}
					matched, matchedNoDRC, total := 0, 0, 0
					denoms := map[string]int{}
					for k, entries := range plan.ZeroTokenQueued {
						var dels []string
						addr, _ := sdk.AccAddressFromBech32(k.Delegator())
						valAddr, _ := sdk.ValAddressFromBech32(k.Validator())
						for _, dl := range dscApp.ValidatorKeeper.GetDelegatorDelegations(ctx, addr, 1000) {
							if dl.Validator == valAddr.String() {
								dels = append(dels, dl.Stake.Stake.Denom)
							}
						}
						for _, e := range entries {
							total++
							amt := sdkmath.NewIntFromBigInt(e.Amount)
							found := ""
							for _, d := range plan.ResolvedBySymbol {
								c, holds := plan.ContractAmount(k.Validator(), k.Delegator(), d)
								h := e.HoldTimestamp.Int64()
								if (h == 0 && c.GTE(amt)) || (h != 0 && holds[h].GTE(amt)) {
									found = d
									break
								}
							}
							if found != "" {
								matched++
								if noDRC[found] {
									matchedNoDRC++
								}
								denoms[found]++
							}
						}
						_ = dels
					}
					fmt.Fprintf(out, "zero-token entries=%d matched to a stake of a coin without DRC20 on the node=%d denoms=%v\n", total, matchedNoDRC, denoms)
					_ = matched
				}
			}
			if outPath != "" {
				type row struct {
					Validator  string            `json:"validator"`
					Delegator  string            `json:"delegator"`
					Denom      string            `json:"denom"`
					NodeFound  bool              `json:"node_found"`
					NodeAmount string            `json:"node_amount"`
					Target     string            `json:"target"`
					Holds      map[int64]string  `json:"holds,omitempty"`
					NodeHolds  map[string]string `json:"node_holds,omitempty"`
					Queued     bool              `json:"queued,omitempty"`
				}
				rows := make([]row, 0, len(plan.Diffs))
				for _, d := range plan.Diffs {
					r := row{Validator: d.Target.Validator.String(), Delegator: d.Target.Delegator.String(), Denom: d.Target.Denom,
						NodeFound: d.NodeFound, NodeAmount: d.NodeAmount.String(), Target: d.Target.Amount.String(), Holds: map[int64]string{}, Queued: d.Queued}
					for end, amt := range d.Target.Holds {
						r.Holds[end] = amt.String()
					}
					if del, ok := dscApp.ValidatorKeeper.GetDelegation(ctx, d.Target.Delegator, d.Target.Validator, d.Target.Denom); ok {
						r.NodeHolds = map[string]string{}
						for i, h := range del.Stake.Holds {
							r.NodeHolds[fmt.Sprintf("%d#%d@%d", h.HoldEndTime, i, h.HoldStartTime)] = h.Amount.String()
						}
					}
					rows = append(rows, r)
				}
				bz, _ := json.MarshalIndent(rows, "", " ")
				if err := os.WriteFile(outPath, bz, 0o644); err != nil {
					return err
				}
			}
			if p, _ := cmd.Flags().GetString("expired-holds-out"); p != "" {
				bz, _ := json.MarshalIndent(plan.ExpiredHolds, "", " ")
				if err := os.WriteFile(p, bz, 0o644); err != nil {
					return err
				}
				fmt.Fprintf(out, "wrote %d expired hold buckets to %s\n", len(plan.ExpiredHolds), p)
			}
			if !apply {
				return nil
			}
			if dbg := os.Getenv("DESYNC_DEBUG_DELEGATOR"); dbg != "" {
				for _, d := range plan.Diffs {
					if d.Target.Delegator.String() == dbg {
						fmt.Fprintf(out, "  debug plan: %s %q node=%s found=%t target=%s holds=%v\n", d.Target.Validator, d.Target.Denom, d.NodeAmount, d.NodeFound, d.Target.Amount, d.Target.Holds)
					}
				}
			}
			rep, err := app.ApplyDelegationRepair(ctx, dscApp, plan)
			if dbg := os.Getenv("DESYNC_DEBUG_DELEGATOR"); dbg != "" {
				addr, _ := sdk.AccAddressFromBech32(dbg)
				for _, d := range dscApp.ValidatorKeeper.GetDelegatorDelegations(ctx, addr, 1000) {
					fmt.Fprintf(out, "  debug after: %s %q %s\n", d.Validator, d.Stake.Stake.Denom, d.Stake.Stake.Amount)
				}
			}
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "applied: %s\n", rep.String())
			for _, s := range rep.Skipped {
				fmt.Fprintf(out, "  skipped: %s\n", s)
			}
			again, err := app.PlanDelegationRepair(ctx, dscApp)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "re-plan: %s\n", again.Stats.String())
			for i, d := range again.Diffs {
				if i == 20 {
					break
				}
				fmt.Fprintf(out, "  residual: %s %s %q node=%s(found=%t) target=%s holds=%v\n", d.Target.Validator, d.Target.Delegator, d.Target.Denom, d.NodeAmount, d.NodeFound, d.Target.Amount, d.Target.Holds)
			}
			if again.Stats.Differing != 0 {
				return fmt.Errorf("re-plan still differs on %d keys", again.Stats.Differing)
			}
			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory")
	cmd.Flags().Int64("time", 0, "block time (unix) for hold maturity; default now")
	cmd.Flags().String("out", "", "write every differing key as JSON")
	cmd.Flags().Int64("queue-from", -1, "also price auto-unbond queue entries from this index")
	cmd.Flags().String("expired-holds-out", "", "write expired, never-reset contract hold buckets (input for reset-expired-holds.ts)")
	cmd.Flags().Bool("apply", true, "apply the plan in the cache context and re-plan")
	return cmd
}
