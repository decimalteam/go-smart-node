package app

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"bitbucket.org/decimalteam/go-smart-node/contracts"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// Atomic delegation-contract upgrade, run inside the recovery upgrade handler so the
// release needs no owner transactions, pause window or migration scripts:
//
//  1. deploy DelegationNFTLib and the new DecimalDelegation implementation (the
//     flattened build that delegation/scripts/upgrade.ts deploys, embedded below)
//     from the validator module account;
//  2. make that account the proxy owner for the duration of the handler, pause,
//     upgrade the proxy, migrate the frozen-stake queue to the holdStartTime layout
//     with the contract's own migrateFrozenStakesLayout, unpause, and write the
//     original owner back;
//  3. reset the expired, never-reset hold buckets with the permissionless resetHolds,
//     so the auto-unbond queue (holdTimestamp=0 entries) can withdraw them.
//
// It all happens in one block: users never see a paused or half-migrated contract.
// A contract that is already upgraded and migrated is not upgraded again.

//go:embed delegationupgrade/delegation_flattened.json
var delegationFlattenedJSON []byte

type linkRef struct {
	Start  int `json:"start"`
	Length int `json:"length"`
}

type delegationFlattened struct {
	SourceSha256 string `json:"sourceSha256"`
	Lib          struct {
		Bytecode         string `json:"bytecode"`
		DeployedBytecode string `json:"deployedBytecode"`
	} `json:"lib"`
	Impl struct {
		Bytecode               string                          `json:"bytecode"`
		DeployedBytecode       string                          `json:"deployedBytecode"`
		LinkReferences         map[string]map[string][]linkRef `json:"linkReferences"`
		DeployedLinkReferences map[string]map[string][]linkRef `json:"deployedLinkReferences"`
	} `json:"impl"`
}

// OpenZeppelin v5 OwnableUpgradeable ERC-7201 storage slot.
var ozOwnableSlot = common.HexToHash("0x9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300")

// EIP-1967 implementation slot.
var eip1967ImplSlot = common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")

// migrateFrozenChunk keeps one migrateFrozenStakesLayout call (~200k gas per entry,
// cold storage; 30 entries ≈ 6M) under the per-call EVM gas cap, which is
// min(block max_gas, 25M) = 10M on mainnet. Running out of gas
// inside the library delegatecall surfaces as a bare revert, so a reverted chunk is
// retried at half size before giving up.
const migrateFrozenChunk = 30

// resetHoldsChunk bounds one resetHolds call.
const resetHoldsChunk = 20

const delegationAdminABI = `[
 {"type":"function","name":"pause","stateMutability":"nonpayable","inputs":[],"outputs":[]},
 {"type":"function","name":"unpause","stateMutability":"nonpayable","inputs":[],"outputs":[]},
 {"type":"function","name":"paused","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"bool"}]},
 {"type":"function","name":"owner","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"address"}]},
 {"type":"function","name":"upgrade","stateMutability":"nonpayable","inputs":[{"name":"newImpl","type":"address"},{"name":"data","type":"bytes"}],"outputs":[]},
 {"type":"function","name":"migrateFrozenStakesLayout","stateMutability":"nonpayable","inputs":[{"name":"count","type":"uint256"}],"outputs":[]},
 {"type":"function","name":"frozenStakesMigrationProgress","stateMutability":"view","inputs":[],"outputs":[{"name":"migrated","type":"uint256"},{"name":"total","type":"uint256"}]},
 {"type":"function","name":"resetHolds","stateMutability":"nonpayable","inputs":[{"name":"validator","type":"address"},{"name":"delegator","type":"address"},{"name":"token","type":"address"},{"name":"holdTimestamps","type":"uint256[]"}],"outputs":[]}
]`

// DelegationContractUpgradeReport summarizes the atomic contract upgrade.
type DelegationContractUpgradeReport struct {
	AlreadyUpgraded  bool
	OldImpl          common.Address
	Lib, Impl        common.Address
	SourceSha256     string
	FrozenMigrated   uint64
	MigrationCalls   int
	HoldsReset       int
	HoldResetSkipped int
}

func (r DelegationContractUpgradeReport) String() string {
	if r.AlreadyUpgraded {
		return "contract already upgraded and migrated, upgrade skipped"
	}
	return fmt.Sprintf("impl %s -> %s (lib %s, source sha256 %s), frozen migrated %d in %d calls",
		r.OldImpl.Hex(), r.Impl.Hex(), r.Lib.Hex(), r.SourceSha256, r.FrozenMigrated, r.MigrationCalls)
}

type delegationAdmin struct {
	app   *DSC
	abi   abi.ABI
	from  common.Address
	proxy common.Address
}

func newDelegationAdmin(ctx sdk.Context, app *DSC) (*delegationAdmin, error) {
	proxyHex, err := contracts.GetAddressFromContractCenter(ctx, &app.EvmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return nil, fmt.Errorf("resolve delegation contract: %w", err)
	}
	a, err := abi.JSON(strings.NewReader(delegationAdminABI))
	if err != nil {
		return nil, err
	}
	return &delegationAdmin{app: app, abi: a, from: common.Address(validatortypes.ModuleAddress), proxy: common.HexToAddress(proxyHex)}, nil
}

func (d *delegationAdmin) call(ctx sdk.Context, method string, args ...interface{}) ([]interface{}, error) {
	commit := d.abi.Methods[method].StateMutability != "view"
	res, err := d.app.EvmKeeper.CallEVM(ctx, d.abi, d.from, d.proxy, commit, method, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}
	if res.Failed() {
		return nil, fmt.Errorf("%s reverted: %s", method, res.VmError)
	}
	return d.abi.Unpack(method, res.Ret)
}

// UpgradeDelegationContract runs steps 1 and 2. Any error aborts the upgrade handler
// (chain halt at H), so the chain never continues with a half-upgraded contract.
func UpgradeDelegationContract(ctx sdk.Context, app *DSC) (DelegationContractUpgradeReport, error) {
	var rep DelegationContractUpgradeReport
	if migrated, err := delegationContractMigrated(ctx, app); err == nil && migrated {
		rep.AlreadyUpgraded = true
		return rep, nil
	}
	d, err := newDelegationAdmin(ctx, app)
	if err != nil {
		return rep, err
	}
	var art delegationFlattened
	if err := json.Unmarshal(delegationFlattenedJSON, &art); err != nil {
		return rep, fmt.Errorf("decode embedded delegation build: %w", err)
	}
	rep.SourceSha256 = art.SourceSha256
	ek := &app.EvmKeeper

	// The deployer's nonce determines the deployed addresses, so it needs an account;
	// it must stay a module account (the per-block reward mints through it).
	app.AccountKeeper.GetModuleAccount(ctx, validatortypes.ModuleName)
	rep.OldImpl = common.BytesToAddress(ek.GetState(ctx, d.proxy, eip1967ImplSlot).Bytes())

	// 1. deploy library and implementation
	lib, err := d.deploy(ctx, common.FromHex(art.Lib.Bytecode))
	if err != nil {
		return rep, fmt.Errorf("deploy DelegationNFTLib: %w", err)
	}
	implCode, err := linkLibrary(art.Impl.Bytecode, art.Impl.LinkReferences, lib)
	if err != nil {
		return rep, err
	}
	impl, err := d.deploy(ctx, implCode)
	if err != nil {
		return rep, fmt.Errorf("deploy DecimalDelegation: %w", err)
	}
	wantRuntime, err := linkLibrary(art.Impl.DeployedBytecode, art.Impl.DeployedLinkReferences, lib)
	if err != nil {
		return rep, err
	}
	if got := ek.GetCode(ctx, common.BytesToHash(ek.GetAccount(ctx, impl).CodeHash)); !bytes.Equal(got, wantRuntime) {
		return rep, fmt.Errorf("deployed implementation runtime (%d bytes) differs from the embedded build (%d bytes)", len(got), len(wantRuntime))
	}
	rep.Lib, rep.Impl = lib, impl

	// 2. take ownership, pause, upgrade, migrate, unpause, write the owner back
	owner := ek.GetState(ctx, d.proxy, ozOwnableSlot)
	ek.SetState(ctx, d.proxy, ozOwnableSlot, common.BytesToHash(d.from.Bytes()).Bytes())

	out, err := d.call(ctx, "paused")
	if err != nil {
		return rep, err
	}
	wasPaused := out[0].(bool)
	if !wasPaused {
		if _, err := d.call(ctx, "pause"); err != nil {
			return rep, err
		}
	}
	if _, err := d.call(ctx, "upgrade", impl, []byte{}); err != nil {
		return rep, err
	}
	if got := common.BytesToAddress(ek.GetState(ctx, d.proxy, eip1967ImplSlot).Bytes()); got != impl {
		return rep, fmt.Errorf("proxy implementation is %s after upgrade, want %s", got.Hex(), impl.Hex())
	}
	chunk := int64(migrateFrozenChunk)
	for {
		out, err := d.call(ctx, "frozenStakesMigrationProgress")
		if err != nil {
			return rep, err
		}
		migrated, total := out[0].(*big.Int), out[1].(*big.Int)
		if migrated.Cmp(total) >= 0 {
			rep.FrozenMigrated = migrated.Uint64()
			break
		}
		if rep.MigrationCalls > int(total.Uint64())+2 {
			return rep, fmt.Errorf("frozen-stake migration not converging at %s/%s", migrated, total)
		}
		if _, err := d.call(ctx, "migrateFrozenStakesLayout", big.NewInt(chunk)); err != nil {
			if chunk == 1 {
				return rep, err
			}
			chunk /= 2
			continue
		}
		rep.MigrationCalls++
	}
	if !wasPaused {
		if _, err := d.call(ctx, "unpause"); err != nil {
			return rep, err
		}
	}
	ek.SetState(ctx, d.proxy, ozOwnableSlot, owner.Bytes())
	out, err = d.call(ctx, "owner")
	if err != nil {
		return rep, err
	}
	if got := out[0].(common.Address); common.BytesToHash(got.Bytes()) != owner {
		return rep, fmt.Errorf("owner not restored: %s", got.Hex())
	}
	if migrated, err := delegationContractMigrated(ctx, app); err != nil || !migrated {
		return rep, fmt.Errorf("contract not migrated after upgrade (err=%v)", err)
	}
	return rep, nil
}

func (d *delegationAdmin) deploy(ctx sdk.Context, code []byte) (common.Address, error) {
	nonce, err := d.app.AccountKeeper.GetSequence(ctx, d.from.Bytes())
	if err != nil {
		return common.Address{}, err
	}
	addr := crypto.CreateAddress(d.from, nonce)
	if acc := d.app.EvmKeeper.GetAccount(ctx, addr); acc != nil && acc.IsContract() {
		return common.Address{}, fmt.Errorf("deploy address %s already has code", addr.Hex())
	}
	res, err := d.app.EvmKeeper.CallEVMWithData(ctx, d.from, nil, code, true)
	if err != nil {
		return common.Address{}, err
	}
	if res.Failed() {
		return common.Address{}, fmt.Errorf("create reverted: %s", res.VmError)
	}
	acc := d.app.EvmKeeper.GetAccount(ctx, addr)
	if acc == nil || !acc.IsContract() {
		return common.Address{}, fmt.Errorf("no contract code at expected address %s", addr.Hex())
	}
	return addr, nil
}

// linkLibrary writes the library address into every link placeholder, like
// utils/flattenedBuild.ts linkBytecode.
func linkLibrary(hexCode string, refs map[string]map[string][]linkRef, lib common.Address) ([]byte, error) {
	if strings.Count(hexCode, "__$") == 0 {
		return nil, fmt.Errorf("embedded bytecode has no link placeholders")
	}
	// Placeholders are not valid hex; blank them before decoding.
	code := []byte(strings.TrimPrefix(hexCode, "0x"))
	n := 0
	for _, libs := range refs {
		for name, rr := range libs {
			if name != "DelegationNFTLib" {
				return nil, fmt.Errorf("unexpected linked library %s", name)
			}
			for _, r := range rr {
				if r.Length != 20 || 2*(r.Start+20) > len(code) {
					return nil, fmt.Errorf("bad link reference %+v", r)
				}
				copy(code[2*r.Start:2*(r.Start+20)], []byte(strings.ToLower(lib.Hex()[2:])))
				n++
			}
		}
	}
	if n == 0 || strings.Contains(string(code), "__$") {
		return nil, fmt.Errorf("unresolved link placeholder")
	}
	return common.FromHex("0x" + string(code)), nil
}

// ResetExpiredHolds moves the given expired hold buckets into their base stakes with
// the permissionless resetHolds, in deterministic order. A batch that reverts is split
// down to single buckets; a single bucket that still reverts is counted as skipped
// and never aborts the upgrade.
func ResetExpiredHolds(ctx sdk.Context, app *DSC, holds []ExpiredHold) (reset, skipped int, err error) {
	d, err := newDelegationAdmin(ctx, app)
	if err != nil {
		return 0, 0, err
	}
	type group struct {
		val, del, tok common.Address
		ends          []int64
	}
	groups := map[string]*group{}
	for _, h := range holds {
		k := h.Validator.Hex() + h.Delegator.Hex() + h.Token.Hex()
		g, ok := groups[k]
		if !ok {
			g = &group{val: h.Validator, del: h.Delegator, tok: h.Token}
			groups[k] = g
		}
		g.ends = append(g.ends, h.HoldTimestamp)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sweep func(g *group, ends []int64)
	sweep = func(g *group, ends []int64) {
		ts := make([]*big.Int, len(ends))
		for i, e := range ends {
			ts[i] = big.NewInt(e)
		}
		// A reverted call leaves no state behind (CallEVM commits only on success).
		if _, err := d.call(ctx, "resetHolds", g.val, g.del, g.tok, ts); err == nil {
			reset += len(ends)
			return
		}
		if len(ends) == 1 {
			skipped++
			return
		}
		mid := len(ends) / 2
		sweep(g, ends[:mid])
		sweep(g, ends[mid:])
	}
	for _, k := range keys {
		g := groups[k]
		sort.Slice(g.ends, func(i, j int) bool { return g.ends[i] < g.ends[j] })
		g.ends = dedupInt64(g.ends)
		for i := 0; i < len(g.ends); i += resetHoldsChunk {
			j := i + resetHoldsChunk
			if j > len(g.ends) {
				j = len(g.ends)
			}
			sweep(g, g.ends[i:j])
		}
	}
	return reset, skipped, nil
}

func dedupInt64(s []int64) []int64 {
	out := s[:0]
	for i, v := range s {
		if i == 0 || v != s[i-1] {
			out = append(out, v)
		}
	}
	return out
}
