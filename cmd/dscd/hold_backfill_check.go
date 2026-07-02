package main

import (
	"fmt"
	"math/big"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// HoldBackfillTestCmd is a DECODER-INDEPENDENT integration check for the
// hold-start-time backfill upgrade (app.BackfillHoldStartTimes). It runs the REAL
// routine against the real fork DB in a never-committed cache store, then re-derives
// each stakeId/slot with a second, independent implementation and reads the contract's
// Stake struct slots directly to prove:
//
//   - the derived stakeId hits a REAL contract stake: slot0 (validator) and slot1
//     (delegator) equal the expected EVM addresses, and slot6 (holdTimestamp) equals
//     the node-side HoldEndTime the id was keyed on (empirical proof the keying is right);
//   - after the routine runs, slot7 (holdStartTime) equals the expected earliest node
//     start (selectHoldStart), and slots 0/1/6 are unchanged (only offset 7 written);
//   - the routine's own `written` count equals the number of coin-hold buckets found.
//
// Nothing is committed. The node must be stopped.
func HoldBackfillTestCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hold-backfill-test",
		Short: "Run the real holdStartTime backfill on the fork DB and verify each write via raw contract slots (writes nothing)",
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
			// A real block time keeps the floor (upgrade block time) sane; real holds carry a
			// non-zero start so the floor only matters for zero-start holds.
			baseHeader := tmproto.Header{Height: height, ChainID: chainID, Time: time.Now()}
			// EVM contract-center calls resolve a coinbase = block proposer's validator, but a
			// hand-built context has an empty ProposerAddress. Seed it with any existing validator's
			// consensus address so GetAddressFromContractCenter can execute.
			ctx0 := sdk.NewContext(cms, baseHeader, false, serverCtx.Logger)
			var proposer []byte
			for _, v := range dscApp.ValidatorKeeper.GetAllValidators(ctx0) {
				ca, e := v.GetConsAddr()
				if e != nil {
					continue
				}
				if _, found := dscApp.ValidatorKeeper.GetValidatorByConsAddr(ctx0, ca); found {
					proposer = ca.Bytes()
					break
				}
			}
			if proposer == nil {
				return fmt.Errorf("no resolvable validator to serve as block proposer for EVM calls")
			}
			baseHeader.ProposerAddress = proposer
			ctx := sdk.NewContext(cms, baseHeader, false, serverCtx.Logger)
			floor := ctx.BlockTime().Unix()

			out := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(out, "holdStartTime backfill LIVE test\tchain-id %s\theight %d\tfloor %d\n\n", chainID, height, floor)

			// --- Independently resolve delegation contract + denom->token (same deterministic
			//     resolvers the routine uses; what we independently CHECK is the id/slot math
			//     and that a real stake lives at each derived id). ---
			delegationHex, err := contracts.GetAddressFromContractCenter(ctx, &dscApp.EvmKeeper, contracts.NameOfSlugForGetAddressDelegation)
			if err != nil {
				return fmt.Errorf("resolve delegation contract: %w", err)
			}
			delegationAddr := common.HexToAddress(delegationHex)
			if delegationAddr == (common.Address{}) {
				return fmt.Errorf("delegation contract address not resolved")
			}
			denomToToken := map[string]common.Address{}
			dscApp.CoinKeeper.IterateCoinDRC(ctx, func(denom, drc20 string) bool {
				if drc20 != "" {
					denomToToken[denom] = common.HexToAddress(drc20)
				}
				return false
			})
			baseDenom := dscApp.CoinKeeper.GetBaseDenom(ctx)
			wdelToken := common.Address{}
			if wdelHex, e := contracts.GetAddressFromContractCenter(ctx, &dscApp.EvmKeeper, contracts.NameOfSlugForGetAddressWDEL); e == nil {
				if a := common.HexToAddress(wdelHex); a != (common.Address{}) {
					denomToToken[baseDenom] = a
					wdelToken = a
				}
			}
			fmt.Fprintf(out, "delegation contract\t%s\nWDEL(base %q) token\t%s\n\n", delegationAddr.Hex(), baseDenom, wdelToken.Hex())

			// --- Independently rebuild the expected buckets (val,del,token,endTime) -> starts. ---
			type bkey struct {
				val, del string
				token    common.Address
				endTime  int64
			}
			startsByKey := map[bkey][]int64{}
			denomOf := map[bkey]string{}
			dscApp.ValidatorKeeper.IterateAllDelegations(ctx, func(d validatortypes.Delegation) bool {
				denom := d.Stake.Stake.Denom
				token, ok := denomToToken[denom]
				if !ok {
					return false
				}
				for _, h := range d.Stake.GetHolds() {
					if h.HoldEndTime == 0 {
						continue
					}
					k := bkey{d.Validator, d.Delegator, token, h.HoldEndTime}
					startsByKey[k] = append(startsByKey[k], h.HoldStartTime)
					denomOf[k] = denom
				}
				return false
			})
			expectedBuckets := len(startsByKey)

			// Snapshot each expected stake's slots BEFORE the routine runs.
			type row struct {
				k                       bkey
				stakeID                 common.Hash
				expStart                int64
				wantVal, wantDel        common.Address
				valB, delB              common.Address // read from slot0/slot1
				holdTsB                 int64          // slot6 before
				startB, startA          int64          // slot7 before/after
				valA, delA              common.Address // slot0/1 after
				holdTsA                 int64          // slot6 after
				hasStake                bool
			}
			rows := make([]row, 0, expectedBuckets)
			for k, starts := range startsByKey {
				id, err := deriveCoinHoldStakeID(k.val, k.del, k.token, k.endTime)
				if err != nil {
					return fmt.Errorf("derive stakeId for %s/%s: %w", k.val, k.del, err)
				}
				wv, wd, err := evmAddrsOf(k.val, k.del)
				if err != nil {
					return err
				}
				base := stakeStructBaseSlot(id)
				valB := slotAddr(ctx, dscApp, delegationAddr, base, 0)
				rows = append(rows, row{
					k: k, stakeID: id, expStart: selectHoldStartLocal(starts, floor),
					wantVal: wv, wantDel: wd,
					valB:    valB,
					delB:    slotAddr(ctx, dscApp, delegationAddr, base, 1),
					holdTsB: slotInt(ctx, dscApp, delegationAddr, base, 6),
					startB:  slotInt(ctx, dscApp, delegationAddr, base, 7),
					hasStake: valB != (common.Address{}),
				})
			}

			// --- Run the REAL routine. ---
			written, err := app.BackfillHoldStartTimes(ctx, dscApp)
			if err != nil {
				return fmt.Errorf("BackfillHoldStartTimes: %w", err)
			}

			// Snapshot AFTER.
			for i := range rows {
				base := stakeStructBaseSlot(rows[i].stakeID)
				rows[i].valA = slotAddr(ctx, dscApp, delegationAddr, base, 0)
				rows[i].delA = slotAddr(ctx, dscApp, delegationAddr, base, 1)
				rows[i].holdTsA = slotInt(ctx, dscApp, delegationAddr, base, 6)
				rows[i].startA = slotInt(ctx, dscApp, delegationAddr, base, 7)
			}

			// --- Verdicts. ---
			sort.Slice(rows, func(i, j int) bool {
				if rows[i].hasStake != rows[j].hasStake {
					return rows[i].hasStake // real stakes first
				}
				return rows[i].k.endTime < rows[j].k.endTime
			})

			realStakes, idMatch, holdTsMatch, startMatch, fieldsStable, phantom := 0, 0, 0, 0, 0, 0
			var firstFail string
			fail := func(msg string) {
				if firstFail == "" {
					firstFail = msg
				}
			}
			fmt.Fprintf(out, "written(routine)=%d\texpectedBuckets(independent)=%d\t%s\n\n",
				written, expectedBuckets, okBad(written == expectedBuckets))
			fmt.Fprintf(out, "-- per hold bucket (real stakes first) --\n")
			fmt.Fprintf(out, "denom\tendTime\tstake?\tid✔\thTs==end\tstartBefore\tstartAfter\texpStart\tstartOK\tfieldsStable\n")
			shown := 0
			for _, r := range rows {
				if !r.hasStake {
					phantom++
					continue
				}
				realStakes++
				idOK := r.valB == r.wantVal && r.delB == r.wantDel
				htsOK := r.holdTsB == r.k.endTime
				stOK := r.startA == r.expStart
				stable := r.valA == r.valB && r.delA == r.delB && r.holdTsA == r.holdTsB
				if idOK {
					idMatch++
				} else {
					fail(fmt.Sprintf("stake %s: slot0/1 %s/%s != expected %s/%s", r.stakeID.Hex(), r.valB.Hex(), r.delB.Hex(), r.wantVal.Hex(), r.wantDel.Hex()))
				}
				if htsOK {
					holdTsMatch++
				} else {
					fail(fmt.Sprintf("stake %s: slot6 holdTimestamp %d != node HoldEndTime %d", r.stakeID.Hex(), r.holdTsB, r.k.endTime))
				}
				if stOK {
					startMatch++
				} else {
					fail(fmt.Sprintf("stake %s: slot7 after=%d != expected %d", r.stakeID.Hex(), r.startA, r.expStart))
				}
				if stable {
					fieldsStable++
				} else {
					fail(fmt.Sprintf("stake %s: non-hold fields moved (val %s->%s del %s->%s hTs %d->%d)", r.stakeID.Hex(), r.valB.Hex(), r.valA.Hex(), r.delB.Hex(), r.delA.Hex(), r.holdTsB, r.holdTsA))
				}
				if shown < 25 {
					shown++
					fmt.Fprintf(out, "%s\t%d\t%s\t%s\t%s\t%d\t%d\t%d\t%s\t%s\n",
						denomOf[r.k], r.k.endTime, "yes", okBad(idOK), okBad(htsOK),
						r.startB, r.startA, r.expStart, okBad(stOK), okBad(stable))
				}
			}
			out.Flush()

			fmt.Fprintf(out, "\n-- summary --\n")
			fmt.Fprintf(out, "buckets total\t%d\n", expectedBuckets)
			fmt.Fprintf(out, "real contract stakes\t%d\n", realStakes)
			fmt.Fprintf(out, "node-hold w/o contract stake (phantom slot; harmless)\t%d\n", phantom)
			fmt.Fprintf(out, "stakeId matches (slot0/1 == expected val/del)\t%d/%d\n", idMatch, realStakes)
			fmt.Fprintf(out, "holdTimestamp(slot6) == node HoldEndTime\t%d/%d\n", holdTsMatch, realStakes)
			fmt.Fprintf(out, "holdStartTime(slot7) after == expected earliest start\t%d/%d\n", startMatch, realStakes)
			fmt.Fprintf(out, "non-hold fields unchanged by write\t%d/%d\n", fieldsStable, realStakes)

			pass := written == expectedBuckets &&
				idMatch == realStakes && holdTsMatch == realStakes &&
				startMatch == realStakes && fieldsStable == realStakes
			if realStakes == 0 {
				fmt.Fprintf(out, "\nRESULT: INCONCLUSIVE — no coin-hold contract stakes on this fork to check.\n")
				fmt.Fprintf(out, "  (Create a delegate-with-hold before running, or point at a fork that has holds.)\n")
			} else if pass {
				fmt.Fprintf(out, "\nRESULT: PASS — backfill wrote the correct earliest start into slot7 of every real coin-hold stake; keying (val/del/holdEnd) and non-hold fields verified.\n")
			} else {
				fmt.Fprintf(out, "\nRESULT: FAIL — %s\n", firstFail)
			}
			out.Flush()
			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory")
	return cmd
}

// --- second, independent implementation of the id/slot math (cross-checks the app pkg) ---

var hbtDelegationBase = func() *big.Int {
	v, _ := new(big.Int).SetString("c1dae510251b57b62087f142cc8746564a134308f8810580b66b08a29b6def00", 16)
	return v
}()
var hbtStakesMappingSlot = new(big.Int).Add(hbtDelegationBase, big.NewInt(1))

func hbtComputeStakeID(validator, delegator, token common.Address, tokenId, holdTimestamp *big.Int) common.Hash {
	packed := make([]byte, 0, 124)
	packed = append(packed, validator.Bytes()...)
	packed = append(packed, delegator.Bytes()...)
	packed = append(packed, token.Bytes()...)
	packed = append(packed, common.LeftPadBytes(tokenId.Bytes(), 32)...)
	packed = append(packed, common.LeftPadBytes(holdTimestamp.Bytes(), 32)...)
	return common.BytesToHash(crypto.Keccak256(packed))
}

func evmAddrsOf(valBech32, delBech32 string) (common.Address, common.Address, error) {
	valAddr, err := sdk.ValAddressFromBech32(valBech32)
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("validator %s: %w", valBech32, err)
	}
	delAddr, err := sdk.AccAddressFromBech32(delBech32)
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("delegator %s: %w", delBech32, err)
	}
	return common.BytesToAddress(valAddr.Bytes()), common.BytesToAddress(delAddr.Bytes()), nil
}

func deriveCoinHoldStakeID(valBech32, delBech32 string, token common.Address, endTime int64) (common.Hash, error) {
	v, d, err := evmAddrsOf(valBech32, delBech32)
	if err != nil {
		return common.Hash{}, err
	}
	return hbtComputeStakeID(v, d, token, big.NewInt(0), big.NewInt(endTime)), nil
}

// stakeStructBaseSlot = keccak256(stakeID ++ mappingSlot) — start of the Stake struct.
func stakeStructBaseSlot(stakeID common.Hash) *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256(
		append(stakeID.Bytes(), common.BigToHash(hbtStakesMappingSlot).Bytes()...),
	))
}

func slotAt(ctx sdk.Context, a *app.DSC, contract common.Address, base *big.Int, off int64) common.Hash {
	slot := common.BigToHash(new(big.Int).Add(base, big.NewInt(off)))
	return a.EvmKeeper.GetState(ctx, contract, slot)
}

func slotAddr(ctx sdk.Context, a *app.DSC, contract common.Address, base *big.Int, off int64) common.Address {
	return common.BytesToAddress(slotAt(ctx, a, contract, base, off).Bytes())
}

func slotInt(ctx sdk.Context, a *app.DSC, contract common.Address, base *big.Int, off int64) int64 {
	return slotAt(ctx, a, contract, base, off).Big().Int64()
}

func selectHoldStartLocal(starts []int64, floor int64) int64 {
	best := int64(0)
	for _, s := range starts {
		if s == 0 {
			s = floor
		}
		if best == 0 || s < best {
			best = s
		}
	}
	if best == 0 {
		return floor
	}
	return best
}

func okBad(b bool) string {
	if b {
		return "OK"
	}
	return "BAD"
}
