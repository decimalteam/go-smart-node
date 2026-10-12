package app

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/token"
	"bitbucket.org/decimalteam/go-smart-node/contracts/validator"
	validatorkeeper "bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// DesyncRepairPlanName is the mainnet upgrade plan that runs the delegation-desync
// repair. Set the real height at release; a registered but unscheduled name never runs.
var DesyncRepairPlanName = "https://repo.decimalchain.com/33990001"

// Contract token types that the node mirrors as coin delegations.
const (
	contractTokenTypeDRC20 = 1
	contractTokenTypeDEL   = 4
)

// DelegationDesyncRepairHandlerCreator is the single mainnet recovery upgrade. In one
// block, with no owner transactions before or after it, it:
//  1. upgrades the delegation contract (deploy, upgrade, frozen-queue migration; see
//     UpgradeDelegationContract), unless the owner already did;
//  2. resets every expired, never-reset contract hold (permissionless resetHolds);
//  3. reconciles node delegations with the contract (the custodian of the staked
//     tokens) using only on-chain state (see PlanDelegationRepair);
//  4. switches on decoding of both event layouts;
//  5. copies node hold start times into the contract;
//  6. re-syncs the per-block reward.
var DelegationDesyncRepairHandlerCreator = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", plan.Name)

		up, err := UpgradeDelegationContract(ctx, app)
		if err != nil {
			return nil, fmt.Errorf("delegation contract upgrade: %w", err)
		}
		logger.Info("delegation contract upgrade: " + up.String())

		// The repair plan lists the expired holds; resetting them changes the contract
		// side, so plan again afterwards.
		pre, err := PlanDelegationRepair(ctx, app)
		if err != nil {
			return nil, err
		}
		reset, skipped, err := ResetExpiredHolds(ctx, app, pre.ExpiredHolds)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("expired holds reset: %d of %d (skipped %d)", reset, len(pre.ExpiredHolds), skipped))

		p, err := PlanDelegationRepair(ctx, app)
		if err != nil {
			return nil, err
		}
		logger.Info("delegation desync repair plan: " + p.Stats.String())
		pst, err := ApplyPreIncidentDrift(&p)
		if err != nil {
			return nil, err
		}
		logger.Info("pre-incident drift kept: " + pst.String())
		rep, err := ApplyDelegationRepair(ctx, app, p)
		if err != nil {
			return nil, err
		}
		logger.Info("delegation desync repair applied: " + rep.String())
		for _, s := range rep.Skipped {
			logger.Info("delegation desync repair skipped: " + s)
		}
		app.ValidatorKeeper.EnableLegacyDelegationEvents(ctx)

		// After the repair, so restored holds get their start time too.
		written, err := BackfillHoldStartTimes(ctx, app)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("hold-start-time backfill wrote %d contract slots", written))

		amount, enabled, err := SyncRewardPerBlockOverride(ctx, app)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("reward-per-block override synced from contract: enabled=%t amount=%s", enabled, amount.String()))
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// delegationContractMigrated reports whether the delegation contract runs the
// holdStartTime implementation and has migrated its whole frozen queue.
func delegationContractMigrated(ctx sdk.Context, app *DSC) (bool, error) {
	hex, err := contracts.GetAddressFromContractCenter(ctx, &app.EvmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return false, err
	}
	// The node's delegation binding predates this view, so use a one-method ABI.
	a, err := abi.JSON(strings.NewReader(`[{"type":"function","name":"frozenStakesMigrationProgress","stateMutability":"view","inputs":[],"outputs":[{"name":"migrated","type":"uint256"},{"name":"total","type":"uint256"}]}]`))
	if err != nil {
		return false, err
	}
	const method = "frozenStakesMigrationProgress"
	res, err := app.EvmKeeper.CallEVM(ctx, a, common.Address(validatortypes.ModuleAddress), common.HexToAddress(hex), false, method)
	if err != nil {
		return false, err
	}
	if res.Failed() {
		return false, fmt.Errorf("%s reverted (old implementation)", method)
	}
	out, err := a.Unpack(method, res.Ret)
	if err != nil || len(out) != 2 {
		return false, fmt.Errorf("unpack %s: %v", method, err)
	}
	migrated, _ := out[0].(*big.Int)
	total, _ := out[1].(*big.Int)
	if migrated == nil || total == nil {
		return false, fmt.Errorf("unexpected %s result", method)
	}
	// >=: after unpause new frozen stakes are appended to the migrated queue.
	return migrated.Cmp(total) >= 0, nil
}

// RepairStats describes what the contract and node state look like and what differs.
type RepairStats struct {
	ContractStakes           int // active coin stakes read from contract storage
	ContractStakesOther      int // non-coin (NFT) stakes, not mirrored by the node
	ContractZeroAmount       int
	PenaltiesApplied         int // stakes whose amount includes not-yet-applied penalties
	FrozenTransfers          int // still-frozen redelegations credited to their destination
	AutoUnbondEntries        int // unprocessed queue entries (node already removed them)
	NodePendingRedelegs      int // node redelegation entries still to be credited by the node
	UnknownTokens            int // contract tokens with no node coin
	ExpiredHoldStakes        int // contract hold stakes past holdTimestamp, never reset (resetHolds sweep needed)
	QueuedKeysRepairedToFull int // keys with queued auto-unbond entries compared against the full contract stake
	unknown                  map[common.Address]int
	Keys                     int // distinct (validator, delegator, denom)
	Matching                 int
	Differing                int
	DiffNodeOnly             int // node delegation, nothing on the contract
	DiffContractOnly         int // contract stake, no node delegation
	DiffAmount               int
	DiffHoldsOnly            int
}

func (s RepairStats) String() string {
	return fmt.Sprintf("expiredHolds=%d contractStakes=%d other=%d zero=%d penalized=%d frozenTransfers=%d autoUnbond=%d nodePendingRedelegations=%d unknownTokens=%d keys=%d matching=%d differing=%d [nodeOnly=%d contractOnly=%d amount=%d holdsOnly=%d]",
		s.ExpiredHoldStakes, s.ContractStakes, s.ContractStakesOther, s.ContractZeroAmount, s.PenaltiesApplied, s.FrozenTransfers, s.AutoUnbondEntries,
		s.NodePendingRedelegs, s.UnknownTokens, s.Keys, s.Matching, s.Differing, s.DiffNodeOnly, s.DiffContractOnly, s.DiffAmount, s.DiffHoldsOnly)
}

// RepairDiff is one key whose node delegation differs from the contract.
type RepairDiff struct {
	Target     validatorkeeper.DelegationTarget
	NodeAmount sdkmath.Int
	NodeFound  bool
	Queued     bool // the key has unprocessed auto-unbond entries
}

// RepairPlan is the full set of differing keys, in deterministic order.
type RepairPlan struct {
	Stats RepairStats
	Diffs []RepairDiff
	// ZeroTokenQueued holds auto-unbond entries enqueued with token 0x0 (coins whose
	// DRC20 address was unknown to the node at enqueue time), keyed by (validator,
	// delegator) with an empty denom.
	ZeroTokenQueued map[repairKey][]stakescan.AutoUnbondEntry
	// ExpiredHolds are contract hold stakes past their holdTimestamp that were never
	// reset. The node already treats them as unheld, so an auto-unbond entry for them
	// is enqueued with holdTimestamp=0 and cannot draw on them until resetHolds moves
	// them into the base stake (permissionless; scripts/reset-expired-holds.ts).
	ExpiredHolds []ExpiredHold
	// ResolvedBySymbol lists coins whose DRC20 address the node had never recorded.
	ResolvedBySymbol []string
	// ContractAmount returns the contract-side amount and hold buckets for a key
	// (diagnostics only).
	ContractAmount func(val, del, denom string) (sdkmath.Int, map[int64]sdkmath.Int)
}

// ExpiredHold names one expired, unreset contract hold bucket.
type ExpiredHold struct {
	Validator     common.Address `json:"validator"`
	Delegator     common.Address `json:"delegator"`
	Token         common.Address `json:"token"`
	HoldTimestamp int64          `json:"holdTimestamp"`
}

type repairKey struct {
	val, del, denom string
}

type repairPosition struct {
	amount sdkmath.Int
	holds  map[int64]sdkmath.Int
}

func (p *repairPosition) add(holdEnd int64, amt sdkmath.Int) {
	p.amount = p.amount.Add(amt)
	if holdEnd != 0 {
		cur, ok := p.holds[holdEnd]
		if !ok {
			cur = sdkmath.ZeroInt()
		}
		p.holds[holdEnd] = cur.Add(amt)
	}
}

// minus returns p - q, clamped at zero per bucket.
func (p *repairPosition) minus(q *repairPosition) *repairPosition {
	out := &repairPosition{amount: p.amount.Sub(sdkmath.MinInt(p.amount, q.amount)), holds: map[int64]sdkmath.Int{}}
	for end, amt := range p.holds {
		if qa, ok := q.holds[end]; ok {
			amt = amt.Sub(sdkmath.MinInt(amt, qa))
		}
		out.holds[end] = amt
	}
	return out
}

// PlanDelegationRepair derives, from on-chain state only, the node delegation every
// (validator, delegator, denom) should have and returns the keys that differ.
//
// Contract side (raw storage of the delegation contract, the token custodian):
//   - every active coin stake, reduced by validator penalties it has not applied yet
//     (the contract applies penalties lazily; the node slashes eagerly);
//   - plus every still-frozen redelegation, credited to its destination (the node
//     completes redelegations from its own queue and ignores the contract's
//     completion event);
//   - minus every unprocessed auto-unbond queue entry (the node removed those
//     delegations when it enqueued them; the contract keeps the stake until processed).
//
// Node side: coin delegations plus node redelegation entries still in flight (the
// node will credit those to the destination by itself when they mature).
func PlanDelegationRepair(ctx sdk.Context, app *DSC) (RepairPlan, error) {
	var plan RepairPlan
	st := &plan.Stats
	vk := app.ValidatorKeeper

	delegationHex, err := contracts.GetAddressFromContractCenter(ctx, &app.EvmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return plan, fmt.Errorf("resolve delegation contract: %w", err)
	}
	delegationAddr := common.HexToAddress(delegationHex)
	if delegationAddr == (common.Address{}) {
		return plan, fmt.Errorf("delegation contract address not resolved")
	}
	storage := readEVMContractStorage(ctx, app, delegationAddr)
	res := stakescan.Reconstruct(storage, stakescan.DelegationBase)
	// stakescan walks a Go map; fix the order so every node resolves (and repairs)
	// coin addresses in the same sequence.
	sort.Slice(res.CoinStakes, func(i, j int) bool {
		return bytes.Compare(res.CoinStakes[i].BaseSlot.Bytes(), res.CoinStakes[j].BaseSlot.Bytes()) < 0
	})

	denoms := map[common.Address]string{}
	denomOf := func(token common.Address) (string, bool) {
		if d, ok := denoms[token]; ok {
			return d, d != ""
		}
		coin, err := app.CoinKeeper.GetCoinByDRC(ctx, token.String())
		d := ""
		if err == nil {
			d = coin.Denom
		} else if symbol, serr := tokenSymbol(ctx, app, token); serr == nil {
			// Same lazy repair the EVM hooks apply: a coin whose DRC20 address was never
			// recorded is resolved by the token's symbol and the address is stored.
			if c, cerr := app.CoinKeeper.GetCoin(ctx, symbol); cerr == nil && c.DRC20Contract == "" {
				plan.ResolvedBySymbol = append(plan.ResolvedBySymbol, c.Denom)
				_ = app.CoinKeeper.UpdateCoinDRC(ctx, symbol, token.String())
				c.DRC20Contract = token.String()
				app.CoinKeeper.SetCoin(ctx, c)
				d = c.Denom
			}
		}
		denoms[token] = d
		return d, d != ""
	}
	keyOf := func(val, del common.Address, denom string) repairKey {
		return repairKey{sdk.ValAddress(val.Bytes()).String(), sdk.AccAddress(del.Bytes()).String(), denom}
	}

	penalties := newPenaltyCache(ctx, app)
	contract := map[repairKey]*repairPosition{}
	pos := func(k repairKey) *repairPosition {
		p, ok := contract[k]
		if !ok {
			p = &repairPosition{amount: sdkmath.ZeroInt(), holds: map[int64]sdkmath.Int{}}
			contract[k] = p
		}
		return p
	}

	for _, cs := range res.CoinStakes {
		if cs.TokenType != contractTokenTypeDRC20 && cs.TokenType != contractTokenTypeDEL {
			st.ContractStakesOther++
			continue
		}
		if cs.Amount.Sign() == 0 {
			st.ContractZeroAmount++
			continue
		}
		denom, ok := denomOf(cs.Token)
		if !ok {
			st.noteUnknown(cs.Token)
			continue
		}
		amount, penalized, err := penalties.apply(cs.Validator, cs.Amount, cs.PenaltyIndex)
		if err != nil {
			return plan, err
		}
		if penalized {
			st.PenaltiesApplied++
		}
		st.ContractStakes++
		if h := cs.HoldTimestamp.Int64(); h != 0 && h <= ctx.BlockTime().Unix() {
			st.ExpiredHoldStakes++
			plan.ExpiredHolds = append(plan.ExpiredHolds, ExpiredHold{Validator: cs.Validator, Delegator: cs.Delegator, Token: cs.Token, HoldTimestamp: h})
		}
		pos(keyOf(cs.Validator, cs.Delegator, denom)).add(cs.HoldTimestamp.Int64(), sdkmath.NewIntFromBigInt(amount))
	}

	// The deployed contract keeps its frozen queue in the pre-hold-start layout
	// (FrozenDeprecated). migrateFrozenStakesLayout copies it into FrozenLive and leaves
	// the old array in place, so once FrozenLive exists it is the only live queue.
	// Indices are preserved, so a partially migrated queue is the migrated prefix plus
	// the not-yet-copied tail of the old array.
	frozen := res.FrozenDeprecated
	if n := len(res.FrozenLive); n > 0 {
		frozen = append([]stakescan.FrozenStake{}, res.FrozenLive...)
		if n < len(res.FrozenDeprecated) {
			frozen = append(frozen, res.FrozenDeprecated[n:]...)
		}
	}
	for _, fz := range frozen {
		if fz.FreezeStatus != 1 || fz.FreezeType != 2 { // Frozen, Transfer
			continue
		}
		if fz.Stake.TokenType != contractTokenTypeDRC20 && fz.Stake.TokenType != contractTokenTypeDEL {
			continue
		}
		denom, ok := denomOf(fz.Stake.Token)
		if !ok {
			st.noteUnknown(fz.Stake.Token)
			continue
		}
		st.FrozenTransfers++
		pos(keyOf(fz.Stake.Validator, fz.Stake.Delegator, denom)).add(fz.Stake.HoldTimestamp.Int64(), sdkmath.NewIntFromBigInt(fz.Stake.Amount))
	}

	// Unprocessed auto-unbond entries. The contract keeps the stake until the bot
	// processes the entry; the node either still holds the delegation (offline
	// validator path, removed when the processed WithdrawRequest arrives) or already
	// removed it (max-delegations trim). Both are consistent.
	queued := map[repairKey]*repairPosition{}
	zeroTokenQueued := map[repairKey][]stakescan.AutoUnbondEntry{}
	for _, e := range stakescan.ReadAutoUnbondQueueAt(storage, stakescan.AutoUnbondQueueSlotMain) {
		denom, ok := denomOf(e.Token)
		if !ok {
			st.noteUnknown(e.Token)
			if e.Token == (common.Address{}) {
				zeroTokenQueued[repairKey{sdk.ValAddress(e.Validator.Bytes()).String(), sdk.AccAddress(e.Delegator.Bytes()).String(), ""}] = append(
					zeroTokenQueued[repairKey{sdk.ValAddress(e.Validator.Bytes()).String(), sdk.AccAddress(e.Delegator.Bytes()).String(), ""}], e)
			}
			continue
		}
		st.AutoUnbondEntries++
		k := keyOf(e.Validator, e.Delegator, denom)
		q, ok := queued[k]
		if !ok {
			q = &repairPosition{amount: sdkmath.ZeroInt(), holds: map[int64]sdkmath.Int{}}
			queued[k] = q
		}
		q.add(e.HoldTimestamp.Int64(), sdkmath.NewIntFromBigInt(e.Amount))
	}

	plan.ZeroTokenQueued = zeroTokenQueued
	plan.ContractAmount = func(val, del, denom string) (sdkmath.Int, map[int64]sdkmath.Int) {
		if p, ok := contract[repairKey{val, del, denom}]; ok {
			return p.amount, p.holds
		}
		return sdkmath.ZeroInt(), nil
	}

	// Node side.
	node := map[repairKey]validatortypes.Delegation{}
	vk.IterateAllDelegations(ctx, func(d validatortypes.Delegation) bool {
		if d.Stake.Type == validatortypes.StakeType_Coin {
			node[repairKey{d.Validator, d.Delegator, d.Stake.Stake.Denom}] = d
		}
		return false
	})
	pending := map[repairKey]*repairPosition{}
	vk.IterateRedelegations(ctx, func(_ int64, red validatortypes.Redelegation) bool {
		for _, e := range red.Entries {
			if e.Stake.Type != validatortypes.StakeType_Coin {
				continue
			}
			st.NodePendingRedelegs++
			k := repairKey{red.ValidatorDst, red.Delegator, e.Stake.Stake.Denom}
			p, ok := pending[k]
			if !ok {
				p = &repairPosition{amount: sdkmath.ZeroInt(), holds: map[int64]sdkmath.Int{}}
				pending[k] = p
			}
			p.add(0, e.Stake.Stake.Amount)
			for _, h := range e.Stake.Holds {
				cur, ok := p.holds[h.HoldEndTime]
				if !ok {
					cur = sdkmath.ZeroInt()
				}
				p.holds[h.HoldEndTime] = cur.Add(h.Amount)
			}
		}
		return false
	})
	// A redelegation the node is still going to credit is already counted in the
	// contract's frozen transfers; take it out of the target so it is not credited twice.
	for k, p := range pending {
		c := pos(k)
		c.amount = c.amount.Sub(sdkmath.MinInt(c.amount, p.amount))
		for end, amt := range p.holds {
			if cur, ok := c.holds[end]; ok {
				c.holds[end] = cur.Sub(sdkmath.MinInt(cur, amt))
			}
		}
	}

	keys := map[repairKey]bool{}
	for k := range contract {
		keys[k] = true
	}
	for k := range node {
		keys[k] = true
	}
	sorted := make([]repairKey, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Slice(sorted, func(i, j int) bool {
		a, b := sorted[i], sorted[j]
		if a.val != b.val {
			return a.val < b.val
		}
		if a.del != b.del {
			return a.del < b.del
		}
		return a.denom < b.denom
	})

	now := ctx.BlockTime().Unix()
	for _, k := range sorted {
		c, hasContract := contract[k]
		if !hasContract {
			c = &repairPosition{amount: sdkmath.ZeroInt(), holds: map[int64]sdkmath.Int{}}
		}
		d, hasNode := node[k]
		nodeAmt := sdkmath.ZeroInt()
		if hasNode {
			nodeAmt = d.Stake.Stake.Amount
		}
		if !hasNode && c.amount.IsZero() {
			continue
		}
		st.Keys++
		holdsMatch := func(p *repairPosition) bool {
			return !hasNode || validatorkeeper.HoldsMatch(d.Stake.Holds, p.holds, now)
		}
		target := c
		queuedRepaired := false
		if q, ok := queued[k]; ok {
			// A processed entry arrives as a WithdrawRequest that the node subtracts from
			// whatever delegation it holds for the key, or skips when there is none. So
			// the only states that stay consistent through processing are: the node holds
			// the full contract stake C (offline path), or it holds nothing and the queue
			// covers all of C (trim path). Anything else is repaired to C.
			if !hasNode && c.minus(q).amount.IsZero() {
				st.Matching++
				continue
			}
			st.QueuedKeysRepairedToFull++
			queuedRepaired = true
		}
		amountEqual := nodeAmt.Equal(target.amount)
		if amountEqual && holdsMatch(target) {
			st.Matching++
			continue
		}
		st.Differing++
		switch {
		case hasNode && target.amount.IsZero():
			st.DiffNodeOnly++
		case !hasNode:
			st.DiffContractOnly++
		case !amountEqual:
			st.DiffAmount++
		default:
			st.DiffHoldsOnly++
		}
		valAddr, err := sdk.ValAddressFromBech32(k.val)
		if err != nil {
			return plan, err
		}
		delAddr, err := sdk.AccAddressFromBech32(k.del)
		if err != nil {
			return plan, err
		}
		plan.Diffs = append(plan.Diffs, RepairDiff{
			Target: validatorkeeper.DelegationTarget{
				Validator: valAddr, Delegator: delAddr, Denom: k.denom,
				Amount: target.amount, Holds: target.holds,
			},
			NodeAmount: nodeAmt,
			NodeFound:  hasNode,
			Queued:     queuedRepaired,
		})
	}
	return plan, nil
}

// ApplyDelegationRepair applies every planned target.
func ApplyDelegationRepair(ctx sdk.Context, app *DSC, plan RepairPlan) (validatorkeeper.DesyncRepairReport, error) {
	var rep validatorkeeper.DesyncRepairReport
	for _, d := range plan.Diffs {
		if err := app.ValidatorKeeper.ApplyDelegationTarget(ctx, d.Target, &rep); err != nil {
			return rep, fmt.Errorf("repair %s/%s/%s: %w", d.Target.Validator, d.Target.Delegator, d.Target.Denom, err)
		}
	}
	return rep, nil
}

func readEVMContractStorage(ctx sdk.Context, app *DSC, addr common.Address) stakescan.Storage {
	ps := prefix.NewStore(ctx.KVStore(app.GetKey(evmtypes.StoreKey)), evmtypes.AddressStoragePrefix(addr))
	it := ps.Iterator(nil, nil)
	defer it.Close()
	s := stakescan.Storage{}
	for ; it.Valid(); it.Next() {
		s[common.BytesToHash(it.Key())] = common.BytesToHash(it.Value())
	}
	return s
}

// penaltyCache applies a validator's penalty list the way the contract's
// _applyPenaltiesToStake does: amount -= amount*pct/100 for every penalty index the
// stake has not applied yet.
type penaltyCache struct {
	ctx    sdk.Context
	app    *DSC
	master common.Address
	byVal  map[common.Address][]uint16
}

func newPenaltyCache(ctx sdk.Context, app *DSC) *penaltyCache {
	return &penaltyCache{ctx: ctx, app: app, byVal: map[common.Address][]uint16{}}
}

func (c *penaltyCache) list(val common.Address) ([]uint16, error) {
	if l, ok := c.byVal[val]; ok {
		return l, nil
	}
	if c.master == (common.Address{}) {
		hex, err := contracts.GetAddressFromContractCenter(c.ctx, &c.app.EvmKeeper, contracts.NameOfSlugForGetAddressMasterValidator)
		if err != nil {
			return nil, fmt.Errorf("resolve master-validator contract: %w", err)
		}
		c.master = common.HexToAddress(hex)
	}
	a, err := validator.ValidatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	const method = "getValidatorPenalties"
	res, err := c.app.EvmKeeper.CallEVM(c.ctx, *a, common.Address(validatortypes.ModuleAddress), c.master, false, method, val)
	if err != nil {
		return nil, fmt.Errorf("call %s(%s): %w", method, val.Hex(), err)
	}
	out, err := a.Unpack(method, res.Ret)
	if err != nil {
		return nil, fmt.Errorf("unpack %s: %w", method, err)
	}
	l, ok := out[0].([]uint16)
	if !ok {
		return nil, fmt.Errorf("unexpected %s type %T", method, out[0])
	}
	c.byVal[val] = l
	return l, nil
}

func (c *penaltyCache) apply(val common.Address, amount, index *big.Int) (*big.Int, bool, error) {
	l, err := c.list(val)
	if err != nil {
		return nil, false, err
	}
	out, applied := applyContractPenalties(amount, l, index)
	return out, applied, nil
}

// applyContractPenalties mirrors DecimalDelegation._applyPenaltiesToStake: every penalty
// at or after the stake's penalty index cuts amount*pct/100 (integer division) in turn.
func applyContractPenalties(amount *big.Int, penalties []uint16, index *big.Int) (*big.Int, bool) {
	if !index.IsUint64() || index.Uint64() >= uint64(len(penalties)) || amount.Sign() == 0 {
		return amount, false
	}
	cur := new(big.Int).Set(amount)
	for i := index.Uint64(); i < uint64(len(penalties)); i++ {
		cut := new(big.Int).Mul(cur, big.NewInt(int64(penalties[i])))
		cut.Quo(cut, big.NewInt(100))
		cur.Sub(cur, cut)
	}
	return cur, true
}

// UnknownTokenCounts lists contract tokens with no node coin and how many stakes use them.
func (s RepairStats) UnknownTokenCounts() map[common.Address]int { return s.unknown }

// tokenSymbol reads a DRC20's symbol(), returning an error instead of panicking when
// the address does not answer like a token.
func tokenSymbol(ctx sdk.Context, app *DSC, tokenAddr common.Address) (string, error) {
	a, err := token.TokenMetaData.GetAbi()
	if err != nil {
		return "", err
	}
	res, err := app.EvmKeeper.CallEVM(ctx, *a, common.Address(validatortypes.ModuleAddress), tokenAddr, false, "symbol")
	if err != nil {
		return "", err
	}
	if res.Failed() {
		return "", fmt.Errorf("symbol() reverted: %s", res.VmError)
	}
	out, err := a.Unpack("symbol", res.Ret)
	if err != nil || len(out) != 1 {
		return "", fmt.Errorf("symbol() undecodable")
	}
	sym, ok := out[0].(string)
	if !ok {
		return "", fmt.Errorf("symbol() type %T", out[0])
	}
	return sym, nil
}

func (s *RepairStats) noteUnknown(tok common.Address) {
	s.UnknownTokens++
	if s.unknown == nil {
		s.unknown = map[common.Address]int{}
	}
	s.unknown[tok]++
}

// Validator returns the key's validator (bech32).
func (k repairKey) Validator() string { return k.val }

// Delegator returns the key's delegator (bech32).
func (k repairKey) Delegator() string { return k.del }

// RepairValue sums, in base coin at current prices (the node's ToBaseCoin, as used
// for validator power), how much delegated value each class of the plan adds or removes.
type RepairValue struct {
	Created, Increased, Decreased, Removed sdkmath.Int
}

// ValuePlan prices every planned change.
func ValuePlan(ctx sdk.Context, app *DSC, plan RepairPlan) RepairValue {
	v := RepairValue{sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt()}
	for _, d := range plan.Diffs {
		delta := d.Target.Amount.Sub(d.NodeAmount)
		base := app.ValidatorKeeper.ToBaseCoin(ctx, sdk.Coin{Denom: d.Target.Denom, Amount: delta.Abs()}).Amount
		switch {
		case !d.NodeFound:
			v.Created = v.Created.Add(base)
		case d.Target.Amount.IsZero():
			v.Removed = v.Removed.Add(base)
		case delta.IsPositive():
			v.Increased = v.Increased.Add(base)
		case delta.IsNegative():
			v.Decreased = v.Decreased.Add(base)
		}
	}
	return v
}

// AutoUnbondQueueValue prices the auto-unbond queue entries with index >= from.
func AutoUnbondQueueValue(ctx sdk.Context, app *DSC, from uint64) (int, sdkmath.Int, error) {
	hex, err := contracts.GetAddressFromContractCenter(ctx, &app.EvmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return 0, sdkmath.ZeroInt(), err
	}
	storage := readEVMContractStorage(ctx, app, common.HexToAddress(hex))
	n, total := 0, sdkmath.ZeroInt()
	for _, e := range stakescan.ReadAutoUnbondQueueAt(storage, stakescan.AutoUnbondQueueSlotMain) {
		if e.Index < from {
			continue
		}
		coin, err := app.CoinKeeper.GetCoinByDRC(ctx, e.Token.String())
		if err != nil {
			continue
		}
		n++
		total = total.Add(app.ValidatorKeeper.ToBaseCoin(ctx, sdk.Coin{Denom: coin.Denom, Amount: sdkmath.NewIntFromBigInt(e.Amount)}).Amount)
	}
	return n, total, nil
}

// DelegationContractMigrated reports whether the recovery upgrade will run the
// hold-start backfill (contract upgraded and frozen queue fully migrated).
func DelegationContractMigrated(ctx sdk.Context, app *DSC) (bool, error) {
	return delegationContractMigrated(ctx, app)
}
