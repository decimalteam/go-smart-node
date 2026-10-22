package app

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
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
//     original owner back.
//
// It all happens in one block: users never see a paused or half-migrated contract.
// A contract that is already upgraded and migrated is not upgraded again, unless forced
// (testnet, whose contract got the holdStartTime implementation in July).

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

const delegationAdminABI = `[
 {"type":"function","name":"pause","stateMutability":"nonpayable","inputs":[],"outputs":[]},
 {"type":"function","name":"unpause","stateMutability":"nonpayable","inputs":[],"outputs":[]},
 {"type":"function","name":"paused","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"bool"}]},
 {"type":"function","name":"owner","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"address"}]},
 {"type":"function","name":"upgrade","stateMutability":"nonpayable","inputs":[{"name":"newImpl","type":"address"},{"name":"data","type":"bytes"}],"outputs":[]},
 {"type":"function","name":"migrateFrozenStakesLayout","stateMutability":"nonpayable","inputs":[{"name":"count","type":"uint256"}],"outputs":[]},
 {"type":"function","name":"frozenStakesMigrationProgress","stateMutability":"view","inputs":[],"outputs":[{"name":"migrated","type":"uint256"},{"name":"total","type":"uint256"}]}
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
// Without force, a contract that is already upgraded and migrated is left alone.
func UpgradeDelegationContract(ctx sdk.Context, app *DSC, force bool) (DelegationContractUpgradeReport, error) {
	var rep DelegationContractUpgradeReport
	if migrated, err := delegationContractMigrated(ctx, app); !force && err == nil && migrated {
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
