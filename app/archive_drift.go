package app

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"

	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// DriftRow compares one key the repair would change today with the same key at a
// historical height (raw amounts: node stake vs sum of the key's contract stakes).
type DriftRow struct {
	Validator string `json:"validator"`
	Delegator string `json:"delegator"`
	Denom     string `json:"denom"`
	Class     string `json:"class"`
	Queued    bool   `json:"queued,omitempty"`

	NodeNow      string `json:"node_now"`
	ContractNow  string `json:"contract_now"`
	TargetNow    string `json:"target_now"` // what the repair sets the node to
	NodeThen     string `json:"node_then"`
	ContractThen string `json:"contract_then"`
	DriftNow     string `json:"drift_now"`  // node - contract
	DriftThen    string `json:"drift_then"` // node - contract
	DriftNowBase string `json:"drift_now_base_coin"`
}

// Drift classes.
const (
	DriftIncident        = "incident"          // node matched the contract at the historical height
	DriftOldUnchanged    = "old_unchanged"     // already differed then, nothing moved since on either side
	DriftOldNodeMoved    = "old_node_moved"    // already differed then; node side changed since
	DriftOldContractMove = "old_contract_moved" // already differed then; contract side changed since
)

// ArchiveDrift classifies every key of today's repair plan by whether it already
// differed at a historical height. get reads one store key at that height (it must
// be a versioned read). Tooling only: not used by the upgrade handler.
func ArchiveDrift(ctx sdk.Context, app *DSC, getBatch func(stores []string, keys [][]byte) ([][]byte, error), workers int) ([]DriftRow, error) {
	cctx, _ := ctx.CacheContext() // the planner records coin DRC addresses it resolves
	plan, err := PlanDelegationRepair(cctx, app)
	if err != nil {
		return nil, err
	}
	proxyHex, err := contracts.GetAddressFromContractCenter(cctx, &app.EvmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return nil, err
	}
	proxy := common.HexToAddress(proxyHex)
	res := stakescan.Reconstruct(readEVMContractStorage(cctx, app, proxy), stakescan.DelegationBase)

	type key struct{ val, del, denom string }
	diffs := map[key]RepairDiff{}
	for _, d := range plan.Diffs {
		diffs[key{d.Target.Validator.String(), d.Target.Delegator.String(), d.Target.Denom}] = d
	}
	denoms := map[common.Address]string{}
	denomOf := func(tok common.Address) string {
		if d, ok := denoms[tok]; ok {
			return d
		}
		d := ""
		if c, err := app.CoinKeeper.GetCoinByDRC(cctx, tok.String()); err == nil {
			d = c.Denom
		}
		denoms[tok] = d
		return d
	}
	stakesOf := map[key][]stakescan.CoinStake{}
	for _, cs := range res.CoinStakes {
		if cs.TokenType != contractTokenTypeDRC20 && cs.TokenType != contractTokenTypeDEL {
			continue
		}
		k := key{sdk.ValAddress(cs.Validator.Bytes()).String(), sdk.AccAddress(cs.Delegator.Bytes()).String(), denomOf(cs.Token)}
		if _, ok := diffs[k]; ok {
			stakesOf[k] = append(stakesOf[k], cs)
		}
	}

	// Historical reads, in parallel.
	type job struct {
		store string
		key   []byte
	}
	var jobs []job
	keys := make([]key, 0, len(diffs))
	for k := range diffs {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].val != keys[j].val {
			return keys[i].val < keys[j].val
		}
		if keys[i].del != keys[j].del {
			return keys[i].del < keys[j].del
		}
		return keys[i].denom < keys[j].denom
	})
	for _, k := range keys {
		d := diffs[k]
		jobs = append(jobs, job{validatortypes.StoreKey, validatortypes.GetDelegationKey(d.Target.Delegator, d.Target.Validator, k.denom)})
		for _, cs := range stakesOf[k] {
			jobs = append(jobs, job{evmtypes.StoreKey, append(evmtypes.AddressStoragePrefix(proxy), stakescan.StakeAmountSlot(cs.BaseSlot).Bytes()...)})
		}
	}
	const batch = 100
	values := make([][]byte, len(jobs))
	var errs []error
	var emu sync.Mutex
	ch := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for start := range ch {
				end := start + batch
				if end > len(jobs) {
					end = len(jobs)
				}
				stores := make([]string, 0, end-start)
				ks := make([][]byte, 0, end-start)
				for _, j := range jobs[start:end] {
					stores = append(stores, j.store)
					ks = append(ks, j.key)
				}
				vals, err := getBatch(stores, ks)
				if err != nil {
					emu.Lock()
					errs = append(errs, fmt.Errorf("historical reads %d..%d: %w", start, end, err))
					emu.Unlock()
					continue
				}
				copy(values[start:end], vals)
			}
		}()
	}
	for start := 0; start < len(jobs); start += batch {
		ch <- start
	}
	close(ch)
	wg.Wait()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	rows := make([]DriftRow, 0, len(keys))
	i := 0
	for _, k := range keys {
		d := diffs[k]
		nodeThen := sdkmath.ZeroInt()
		if bz := values[i]; len(bz) > 0 {
			del := validatortypes.MustUnmarshalDelegation(app.appCodec, bz)
			nodeThen = del.Stake.Stake.Amount
		}
		i++
		contractNow, contractThen := new(big.Int), new(big.Int)
		for _, cs := range stakesOf[k] {
			contractNow.Add(contractNow, cs.Amount)
			contractThen.Add(contractThen, new(big.Int).SetBytes(values[i]))
			i++
		}
		nodeNow := sdkmath.ZeroInt()
		if del, ok := app.ValidatorKeeper.GetDelegation(cctx, d.Target.Delegator, d.Target.Validator, k.denom); ok {
			nodeNow = del.Stake.Stake.Amount
		}
		cNow, cThen := sdkmath.NewIntFromBigInt(contractNow), sdkmath.NewIntFromBigInt(contractThen)
		driftNow, driftThen := nodeNow.Sub(cNow), nodeThen.Sub(cThen)
		class := DriftIncident
		if !driftThen.IsZero() {
			switch {
			case nodeThen.Equal(nodeNow) && cThen.Equal(cNow):
				class = DriftOldUnchanged
			case !nodeThen.Equal(nodeNow):
				class = DriftOldNodeMoved
			default:
				class = DriftOldContractMove
			}
		}
		base := sdkmath.ZeroInt()
		if !driftNow.IsZero() && k.denom != "" {
			abs := driftNow.Abs()
			base = app.ValidatorKeeper.ToBaseCoin(cctx, sdk.Coin{Denom: k.denom, Amount: abs}).Amount
		}
		rows = append(rows, DriftRow{
			Validator: k.val, Delegator: k.del, Denom: k.denom, Class: class, Queued: d.Queued,
			NodeNow: nodeNow.String(), ContractNow: cNow.String(), TargetNow: d.Target.Amount.String(),
			NodeThen: nodeThen.String(), ContractThen: cThen.String(),
			DriftNow: driftNow.String(), DriftThen: driftThen.String(), DriftNowBase: base.String(),
		})
	}
	return rows, nil
}
