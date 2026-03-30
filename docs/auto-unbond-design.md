# Auto-Unbond System Design

## Problem

When a validator goes offline and never returns, delegators' stakes are locked indefinitely. They must manually discover the situation and undelegate. This is bad UX and a capital efficiency problem.

## Goal

If a validator remains offline for a configurable duration (default: 30 days), automatically unbond **all** delegations from that validator, returning stakes to delegators after the normal unbonding period.

---

## Architecture Overview

```
                    SetOffline / Jail
Validator Online ──────────────────────► Validator Offline
     ▲                                       │
     │                                       │ Store OfflineSince = now
     │ SetOnline                             │
     │ (clears OfflineSince)                 ▼
     │                              ┌─────────────────────┐
     └──────────────────────────────│  EndBlocker          │
                                    │  (every 120 blocks)  │
                                    │                      │
                                    │  Check:              │
                                    │  now - OfflineSince  │
                                    │  > AutoUnbondTimeout?│
                                    └──────────┬───────────┘
                                               │ YES
                                               ▼
                                    ┌──────────────────────┐
                                    │  PHASE 1 (EndBlocker)│
                                    │                      │
                                    │  For each delegation:│
                                    │  CallEVM →           │
                                    │  autoUnbondEnqueue() │
                                    │                      │
                                    │  Writes to EVM       │
                                    │  contract storage    │
                                    │  (no logs, no Cosmos │
                                    │   state changes)     │
                                    └──────────┬───────────┘
                                               │
                                               │ Queue populated
                                               ▼
                                    ┌──────────────────────┐
                                    │  PHASE 2 (Real EVM   │
                                    │  tx, permissionless) │
                                    │                      │
                                    │  Anyone calls        │
                                    │  processAutoUnbond   │
                                    │  (index)             │
                                    │  (1 entry per call,  │
                                    │   skip held stakes)  │
                                    │                      │
                                    │  → _withdrawFor()    │
                                    │  → _freezeStake()    │
                                    │  → WithdrawRequest   │
                                    │    (The Graph sees)  │
                                    │                      │
                                    │  PostTxProcessing    │
                                    │  hook fires:         │
                                    │  → RequestWithdraw() │
                                    │  → Undelegate()      │
                                    │  → UBD queue         │
                                    │    (normal unbonding │
                                    │     period)          │
                                    └──────────────────────┘
```

---

## 1. New State: `OfflineSince` Timestamp

**Storage key:** `ValidatorOfflineSinceKey(valAddr) → timestamp`

**Location:** `x/validator/types/keys.go` (new key prefix), `x/validator/keeper/` (new getter/setter)

```go
// keys.go - new prefix
var ValidatorOfflineSinceKey = []byte{0x45} // new unique prefix

func GetValidatorOfflineSinceKey(valAddr sdk.ValAddress) []byte {
    return append(ValidatorOfflineSinceKey, address.MustLengthPrefix(valAddr)...)
}
```

```go
// keeper/offline_tracker.go - new file
func (k Keeper) SetValidatorOfflineSince(ctx sdk.Context, valAddr sdk.ValAddress, t time.Time)
func (k Keeper) GetValidatorOfflineSince(ctx sdk.Context, valAddr sdk.ValAddress) (time.Time, bool)
func (k Keeper) DeleteValidatorOfflineSince(ctx sdk.Context, valAddr sdk.ValAddress)
func (k Keeper) IterateValidatorOfflineSince(ctx sdk.Context, fn func(valAddr sdk.ValAddress, t time.Time) bool)
```

**Lifecycle:**

| Event | Action |
|---|---|
| `SetOffline` (msg_server.go:261) | `SetValidatorOfflineSince(ctx, valAddr, ctx.BlockTime())` |
| `Jail` (slash.go:196) | `SetValidatorOfflineSince(ctx, valAddr, ctx.BlockTime())` — only if not already set |
| `SetOnline` (msg_server.go:196) | `DeleteValidatorOfflineSince(ctx, valAddr)` |
| Auto-unbond triggered | `DeleteValidatorOfflineSince(ctx, valAddr)` after processing |

**Important:** When setting `OfflineSince` on `Jail`, check if already set (validator could have been offline before being jailed — don't reset the timer).

---

## 2. New Parameter: `AutoUnbondTimeout`

### Proto change (`proto/custom/decimal/validator/v1/params.proto`)

```protobuf
message Params {
  // ... existing fields 1-14 ...

  // auto_unbond_timeout defines the duration after which all stakes are
  // force-unbonded from a continuously offline validator.
  // Set to 0 to disable auto-unbond.
  google.protobuf.Duration auto_unbond_timeout = 15
      [ (gogoproto.nullable) = false, (gogoproto.stdduration) = true ];
}
```

### Go defaults (`x/validator/types/params.go`)

```go
// DefaultAutoUnbondTimeout = 30 days
DefaultAutoUnbondTimeout time.Duration = time.Hour * 24 * 30
```

Add to `DefaultParams()`, `NewParams()`, `Validate()`, `ParamSetPairs()`.

### EVM integration (`x/validator/keeper/abci.go`)

In the `height%120 == 0` block, add a call to read auto-unbond timeout from the EVM contract (similar to how `UndelegationTime` and `RedelegationTime` are already fetched from EVM contracts). This requires a new contract method like `GetTimeAutoUnbond`.

---

## 3. Dual-Layer Event Architecture (CRITICAL)

The system has a **dual-layer architecture** that must be respected:

```
┌─────────────────────────────────┐     ┌──────────────────────────────┐
│         EVM Layer               │     │       Cosmos Layer           │
│                                 │     │                              │
│  DecimalDelegation.sol          │     │  x/validator/keeper          │
│  - StakeUpdated (log)           │     │  - EventDelegate             │
│  - WithdrawRequest (log)        │     │  - EventUndelegate           │
│  - WithdrawCompleted (log)      │     │  - EventUndelegateComplete   │
│  - TransferRequest (log)        │     │  - EventForceUndelegate      │
│                                 │     │                              │
│  Indexed by: The Graph          │     │  Indexed by: Cosmos indexers │
└─────────────────────────────────┘     └──────────────────────────────┘
```

### The `CallEVM` log visibility problem

**`CallEVM` (used by Cosmos modules) does NOT produce The Graph-visible logs.**

Proof — ethermint `CallEVMWithData` flow:
1. Calls `ApplyMessage` (NOT `ApplyTransaction`)
2. `ApplyMessage` → `ApplyMessageWithConfig` creates `statedb.NewEmptyTxConfig(...)` — **no real tx hash**
3. With `commit=true`, `stateDB.Commit()` persists **EVM state changes** (storage, balances)
4. But **logs are only returned in `MsgEthereumTxResponse.Logs`** — they are **never stored as a block receipt**
5. The Graph indexes **block receipts** → it will **never see** logs from `CallEVM`

This means:
- `CallEVM` with `commit=true` → **EVM contract storage IS updated** (stakes, balances)
- `CallEVM` logs → **NOT in any block receipt** → **The Graph CANNOT see them**
- The existing `forceWithdraw()` in `CheckDelegations` already has this same problem

### Normal undelegation flow (EVM→Cosmos) — works correctly:

1. User sends real EVM tx → calls delegation contract `withdraw()`
2. `ApplyTransaction` processes it → creates **receipt with logs** → stored in block
3. Contract emits `WithdrawRequest` + `StakeUpdated` → The Graph sees it (from receipt)
4. `PostTxProcessing` EVM hook catches logs → calls `keeper.RequestWithdraw()` → calls `Undelegate()`
5. Cosmos: creates UBD entry, queues for completion, emits `EventUndelegate`
6. After unbonding period: `CompleteUnbonding()` emits `EventUndelegateComplete`

**Key insight:** The Graph only sees EVM logs from **real user-initiated EVM transactions** that go through `ApplyTransaction` and produce block receipts. Internal `CallEVM` from BeginBlocker/EndBlocker uses `ApplyMessage` which does NOT produce receipts.

### Correct approach: two-phase auto-unbond

Since we cannot produce The Graph-visible EVM logs from EndBlocker, we need a **two-phase approach** where the actual unbonding happens entirely through the standard EVM withdrawal flow in Phase 2:

**Phase 1 — Cosmos EndBlocker (enqueue only):**
- Detect expired validators (offline > timeout)
- Write pending entries to EVM contract queue via `CallEVM` with `commit=true`
- **No Cosmos state changes** — no `Unbond()`, no Cosmos events
- Phase 1 is just the decision + queue population

**Phase 2 — Trustless flush (real EVM tx, permissionless):**
- **Anyone** calls `processAutoUnbond()` on the delegation contract (1 entry per call)
- Contract pops the next queued entry and runs the **standard `_withdraw()` flow**:
  - `_freezeStake()` → `_removeStake()` → push to `_frozenStakes` → emit `WithdrawRequest`
  - Internally emits `StakeAmountUpdated` + `StakeUpdated` + `WithdrawRequest`
- This is a **real EVM transaction** → produces a **real receipt** with logs
- **The Graph** indexes the logs from the receipt
- **`PostTxProcessing` hook** catches `WithdrawRequest` → calls `RequestWithdraw()` → `Undelegate()` on Cosmos side
- Cosmos creates UBD entries in the **normal unbonding queue** with standard unbonding period
- Standard `EventUndelegate` flow via Cosmos indexers

```
Phase 1 (EndBlocker via CallEVM):        Phase 2 (Real EVM tx, anyone can call):
┌──────────────────────┐                 ┌──────────────────────────────────────┐
│ Detect expired       │                 │ Anyone calls                         │
│ validators           │                 │ processAutoUnbond(index)             │
│ Write queue entries  │    ──then──►    │ (skips held stakes by index choice)  │
│ to EVM contract      │                 │                                      │
│ storage              │                 │ Contract: swap-and-pop entry         │
│                      │                 │ → _withdrawFor() per entry           │
│ No Cosmos changes    │                 │ → Emits WithdrawRequest (standard)   │
│ No EVM logs          │                 │ → PostTxProcessing hook fires        │
│                      │                 │ → RequestWithdraw → Undelegate()     │
│                      │                 │ → UBD queue (normal unbonding)       │
│                      │                 │ The Graph + Cosmos indexers see it   │
└──────────────────────┘                 └──────────────────────────────────────┘
```

**Why this is the right approach:**
- **Single source of truth**: unbonding goes through the exact same path as user-initiated withdrawals
- **Both layers updated atomically** in Phase 2 (EVM tx → hook → Cosmos state)
- **The Graph sees standard `WithdrawRequest`** — no special indexer changes needed
- **Cosmos indexers see standard `Undelegate`** — `RequestWithdraw()` → `Undelegate()` flow
- **Normal unbonding period** applies — delegations enter the UBD queue, complete via `CompleteUnbonding()`
- **Trustless**: caller has zero influence on what gets processed — data comes from contract storage

**Why unbonding queue (not immediate return):**
- Goes through the exact same code path as normal user withdrawals — zero special cases
- `PostTxProcessing` hook already handles `WithdrawRequest` → `Undelegate()` → UBD queue
- No need for separate Cosmos event emission — the existing flow handles everything
- Slashing protection during unbonding still works correctly

---

## 4. Auto-Unbond Processing Logic

### Phase 1: Cosmos EndBlocker (`x/validator/keeper/auto_unbond.go`)

```go
// ProcessAutoUnbond checks all offline validators and enqueues their
// delegations for auto-unbond in the EVM delegation contract.
// The actual unbonding happens in Phase 2 when anyone calls
// processAutoUnbond() on the EVM contract.
func (k Keeper) ProcessAutoUnbond(ctx sdk.Context) {
    params := k.GetParams(ctx)

    // Feature disabled
    if params.AutoUnbondTimeout == 0 {
        return
    }

    cutoff := ctx.BlockTime().Add(-params.AutoUnbondTimeout)

    k.IterateValidatorOfflineSince(ctx, func(valAddr sdk.ValAddress, offlineSince time.Time) bool {
        if offlineSince.After(cutoff) {
            return false // not yet expired
        }

        k.enqueueAutoUnbond(ctx, valAddr)
        k.DeleteValidatorOfflineSince(ctx, valAddr)
        return false
    })
}

// enqueueAutoUnbond writes all delegations for a validator into the
// EVM contract's auto-unbond queue. No Cosmos state changes here —
// the actual unbonding happens in Phase 2 via the standard withdrawal flow.
func (k Keeper) enqueueAutoUnbond(ctx sdk.Context, valAddr sdk.ValAddress) {
    delegations := k.GetValidatorDelegations(ctx, valAddr)

    for _, delegation := range delegations {
        err := k.ExecuteAutoUnbondEnqueue(ctx, delegation)
        if err != nil {
            ctx.Logger().Error("auto-unbond: enqueue failed",
                "validator", valAddr,
                "delegator", delegation.Delegator,
                "err", err)
            continue
        }
    }

    ctx.Logger().Info("auto-unbond: enqueued delegations for withdrawal",
        "validator", valAddr,
        "delegations_count", len(delegations),
    )
}
```

### Cosmos-side EVM caller (`x/validator/keeper/evm.go`)

```go
func (k *Keeper) ExecuteAutoUnbondEnqueue(ctx sdk.Context, del types.Delegation) error {
    delegationAddress, err := contracts.GetAddressFromContractCenter(
        ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressDelegation)
    if err != nil {
        return err
    }

    valAddr, _ := sdk.ValAddressFromBech32(del.Validator)
    delAddr, _ := sdk.AccAddressFromBech32(del.Delegator)
    amount := del.GetStake().GetStake().Amount.BigInt()
    coin, _ := k.coinKeeper.GetCoin(ctx, del.GetStake().GetStake().Denom)

    contractDelegation, _ := delegation.DelegationMetaData.GetAbi()
    _, err = k.evmKeeper.CallEVM(ctx, *contractDelegation,
        common.Address{}, // msg.sender == address(0)
        common.HexToAddress(delegationAddress),
        true, // commit — writes queue to EVM contract storage
        "autoUnbondEnqueue",
        common.BytesToAddress(valAddr.Bytes()),
        common.BytesToAddress(delAddr.Bytes()),
        amount,
        common.HexToAddress(coin.DRC20Contract),
    )
    return err
}
```

### Phase 2: Trustless Flush — EVM Contract (`DecimalDelegation.sol`)

```solidity
// ============================================================
// Auto-Unbond Queue (populated by Phase 1 via CallEVM)
// ============================================================

struct AutoUnbondEntry {
    address validator;
    address delegator;
    uint256 amount;
    address token;
    uint256 holdTimestamp;
}

AutoUnbondEntry[] private _autoUnbondQueue;

/// @notice Called by Cosmos EndBlocker (Phase 1) via CallEVM.
/// Only enqueues — does NOT modify stakes or emit logs.
/// msg.sender == address(0) ensures only callable from Cosmos module.
function autoUnbondEnqueue(
    address validator,
    address delegator,
    uint256 amount,
    address token
) external {
    require(msg.sender == address(0), "Not authorized");

    _autoUnbondQueue.push(AutoUnbondEntry({
        validator: validator,
        delegator: delegator,
        amount: amount,
        token: token,
        holdTimestamp: 0
    }));
}

/// @notice Trustless Phase 2: anyone can call to process a pending auto-unbond entry.
/// Processes exactly 1 entry per call to keep gas predictable and avoid
/// PostTxProcessing hook complexity with multiple WithdrawRequest logs.
/// Caller specifies the index to process — this allows skipping entries
/// whose holdTimestamp has not yet expired (they stay in the queue).
/// Uses swap-and-pop for O(1) removal.
/// Calls the standard _withdraw() flow, which:
///   1. Freezes stake (_freezeStake → _removeStake)
///   2. Emits StakeAmountUpdated + StakeUpdated + WithdrawRequest
///   3. PostTxProcessing hook catches WithdrawRequest → Cosmos Undelegate()
function processAutoUnbond(uint256 index) external {
    uint256 len = _autoUnbondQueue.length;
    require(index < len, "Invalid index");

    AutoUnbondEntry memory entry = _autoUnbondQueue[index];

    // Check hold — revert if not expired, caller should skip this entry
    require(
        entry.holdTimestamp == 0 || entry.holdTimestamp < block.timestamp,
        "Hold not expired"
    );

    // Swap with last element and pop (O(1) removal, order doesn't matter)
    _autoUnbondQueue[index] = _autoUnbondQueue[len - 1];
    _autoUnbondQueue.pop();

    // Use the standard withdrawal flow — same as user-initiated withdraw()
    // This emits WithdrawRequest which PostTxProcessing hook will catch
    _withdrawFor(
        entry.validator,
        entry.delegator,
        entry.token,
        0,              // tokenId (0 for DRC20)
        entry.amount,
        entry.holdTimestamp
    );
}

/// @notice View: how many pending entries in the queue
function autoUnbondQueueLength() external view returns (uint256) {
    return _autoUnbondQueue.length;
}

/// @notice View: read a specific queue entry (for bots to check hold status)
function getAutoUnbondEntry(uint256 index) external view returns (AutoUnbondEntry memory) {
    require(index < _autoUnbondQueue.length, "Invalid index");
    return _autoUnbondQueue[index];
}
```

### New internal method: `_withdrawFor` (modification of `_withdraw`)

The existing `_withdraw()` uses `msg.sender` as the delegator. We need a variant
that accepts delegator as a parameter (callable only from `processAutoUnbond`):

```solidity
/// @dev Internal: same as _withdraw but with explicit delegator
/// (msg.sender is the flush caller, not the delegator)
function _withdrawFor(
    address validator,
    address delegator,
    address token,
    uint256 tokenId,
    uint256 amount,
    uint256 holdTimestamp
) internal {
    bool isHoldExpired = holdTimestamp == 0 || holdTimestamp < block.timestamp;
    if (!isHoldExpired) {
        revert HoldNotExpired(holdTimestamp, block.timestamp);
    }

    bytes32 stakeId = _getStakeId(
        validator,
        delegator,    // <-- explicit delegator instead of msg.sender
        token,
        tokenId,
        holdTimestamp
    );
    FrozenStake memory frozenStake = _freezeStake(
        stakeId,
        amount,
        FreezeType.Withdraw
    );
    FrozenStake[] storage frozenStakes = _getDelegationStorage()._frozenStakes;
    frozenStakes.push(frozenStake);

    uint256 stakeIndex = frozenStakes.length - 1;
    emit WithdrawRequest(stakeId, stakeIndex, frozenStake);
}
```

### Hook into EndBlocker (`x/validator/keeper/abci.go`)

```go
func EndBlocker(ctx sdk.Context, k Keeper, req abci.RequestEndBlock) []abci.ValidatorUpdate {
    // ... existing code ...

    if height%120 == 0 {
        // ... existing 120-block logic ...

        // Phase 1: Enqueue auto-unbonds for expired validators
        k.ProcessAutoUnbond(ctx)
    }

    return updates
}
```

### What happens when Phase 2 runs

1. Anyone reads `autoUnbondQueueLength()` and `getAutoUnbondEntry(i)` to find processable entries
2. Caller sends `processAutoUnbond(index)` as a real EVM tx
3. Contract checks hold — if active, reverts (caller picks a different index)
4. Entry is swap-and-popped from the queue, then `_withdrawFor()` runs the standard withdrawal:
   - `_freezeStake()` → `_removeStake()` (updates EVM stake storage)
   - Emits `StakeAmountUpdated`, `StakeUpdated`, `WithdrawRequest`
5. EVM tx produces a receipt with these logs → **The Graph indexes them**
6. `PostTxProcessing` hook fires for the `WithdrawRequest` log:
   - Calls `k.RequestWithdraw()` → `k.Undelegate()` on Cosmos side
   - Creates UBD entry with standard unbonding period
   - Entry enters the normal unbonding queue
7. After unbonding period, `CompleteUnbonding()` runs in EndBlocker:
   - Emits `EventUndelegateComplete` → Cosmos indexers see it

### Frozen (held) stakes

Stakes with an active `holdTimestamp` (hold not yet expired) **stay in the queue** until the hold expires. They are not lost — the bot or any caller simply skips them and processes other entries first. Once `block.timestamp > holdTimestamp`, the entry becomes processable.

This avoids the queue getting stuck on a held entry, while still respecting hold semantics.

### Phase 2 trigger

`processAutoUnbond(index)` is **permissionless** — anyone can call it:
- Processes exactly **1 entry per call** (predictable gas, clean PostTxProcessing hook handling)
- Caller specifies **which entry** to process (skip held stakes, no stuck queue)
- Still **trustless** — index only selects *which* entry, caller cannot influence withdrawal parameters
- **Delegators themselves** — inspect queue and process their own entries
- **Protocol bot** — scan queue, skip held entries, process all eligible
- **Any user** — costs only gas, no trust implications

---

## 5. Event Summary

### Phase 1 (EndBlocker — enqueue only):

| Layer | What happens | Visible to |
|---|---|---|
| EVM storage | `CallEVM` → `autoUnbondEnqueue()` writes queue entries | Contract storage only |
| Cosmos state | **No changes** | — |
| EVM logs | **None** (CallEVM → ApplyMessage → no receipt) | — |
| Cosmos events | **None** | — |

### Phase 2 (Real EVM tx — full standard withdrawal flow):

| Layer | What happens | Visible to |
|---|---|---|
| EVM storage | `_withdrawFor()` → `_removeStake()` updates stakes, `_frozenStakes` populated, queue cleared | On-chain state |
| EVM logs | `StakeAmountUpdated` + `StakeUpdated` + `WithdrawRequest` (standard events) | **The Graph** (real receipt) |
| Cosmos state | `PostTxProcessing` → `RequestWithdraw()` → `Undelegate()` → UBD queue entry | On-chain state |
| Cosmos events | Standard flow from `Undelegate()` | Cosmos indexers |
| Caller | Anyone — no trust required | — |

### After unbonding period (standard EndBlocker):

| Layer | What happens | Visible to |
|---|---|---|
| Cosmos state | `CompleteUnbonding()` removes UBD entry | On-chain state |
| Cosmos events | `EventUndelegateComplete` | Cosmos indexers |

### End result: identical to user-initiated withdrawal
- The Graph sees standard `WithdrawRequest` — no special handling needed
- Cosmos indexers see standard `Undelegate` + `UndelegateComplete` — no special handling needed
- No new event types on either layer
- No trusted party needed

---

## 6. Genesis State

### Add to genesis export/import

```protobuf
// genesis.proto - add to GenesisState
repeated ValidatorOfflineSince validator_offline_since = X;

message ValidatorOfflineSince {
    string validator_address = 1;
    google.protobuf.Timestamp offline_since = 2
        [(gogoproto.nullable) = false, (gogoproto.stdtime) = true];
}
```

Export all `OfflineSince` entries on genesis export, restore them on init.

---

## 7. Migration Considerations

### Existing offline validators

On upgrade, validators that are currently offline have no `OfflineSince` timestamp.

In the migration handler, set `OfflineSince = upgradeTime` for **all currently offline validators**. This gives them the full timeout period from the upgrade moment. After `AutoUnbondTimeout` passes, any that are still offline will be auto-unbonded.

```go
func MigrateAutoUnbond(ctx sdk.Context, k keeper.Keeper) {
    validators := k.GetAllValidators(ctx)
    for _, val := range validators {
        if !val.Online {
            k.SetValidatorOfflineSince(ctx, val.GetOperator(), ctx.BlockTime())
        }
    }
}
```

---

## 8. Key Design Decisions & Rationale

| Decision | Rationale |
|---|---|
| **Two-phase approach** | `CallEVM` from EndBlocker uses `ApplyMessage` (not `ApplyTransaction`) — logs never stored in block receipts. Phase 1 enqueues, Phase 2 executes via real EVM tx. |
| **Unbond via standard EVM withdrawal** | Phase 2 `processAutoUnbond()` calls `_withdrawFor()` — the exact same `_freezeStake` → `_removeStake` → `WithdrawRequest` flow as user-initiated withdrawals. Zero special cases. |
| **PostTxProcessing does the Cosmos sync** | `WithdrawRequest` log → hook → `RequestWithdraw()` → `Undelegate()` → UBD queue. Same path as every other EVM withdrawal. No manual Cosmos event emission needed. |
| **Trustless Phase 2** | Queue lives in EVM contract storage. `processAutoUnbond()` is permissionless, takes no trust-sensitive args. Data comes from contract storage written by Phase 1. |
| **Normal unbonding period** | Goes through the standard UBD queue. Same code path as user withdrawals — slashing protection, indexer compatibility, no special cases. |
| **Check every 120 blocks** | Reuses existing periodic processing cycle. Auto-unbond timeout is 30+ days so 120-block (~10 min) granularity is fine. |
| **Track OfflineSince as KV store** | Simple, deterministic. No need for a time-based queue since we check periodically. |
| **Don't slash on auto-unbond** | This is not misbehavior — it's abandonment. Slashing already happens via the existing missed-blocks mechanism. Auto-unbond is a protective measure for delegators. |
| **`timeout = 0` disables feature** | Simple kill switch. |

---

## 9. Files to Modify/Create

### Cosmos side (go-smart-node)

| File | Action |
|---|---|
| `proto/custom/decimal/validator/v1/params.proto` | Add `auto_unbond_timeout` field (15) |
| `proto/custom/decimal/validator/v1/genesis.proto` | Add `ValidatorOfflineSince` |
| `x/validator/types/params.go` | Add default, validation, param key |
| `x/validator/types/keys.go` | Add `ValidatorOfflineSinceKey` prefix |
| **`x/validator/keeper/offline_tracker.go`** | **New** — CRUD for OfflineSince |
| **`x/validator/keeper/auto_unbond.go`** | **New** — ProcessAutoUnbond + enqueueAutoUnbond (Phase 1 only) |
| `x/validator/keeper/abci.go` | Hook ProcessAutoUnbond into EndBlocker (120-block cycle) |
| `x/validator/keeper/msg_server.go` | Set/clear OfflineSince on SetOnline/SetOffline |
| `x/validator/keeper/evm_hooks.go` | Set/clear OfflineSince on SetOnlineFromEvm/SetOfflineFromEvm |
| `x/validator/keeper/slash.go` | Set OfflineSince on Jail (if not already set) |
| `x/validator/keeper/genesis.go` | Export/import OfflineSince state |
| `contracts/` | Add `GetTimeAutoUnbond` EVM contract call |

### EVM side (decimal-smart-contracts)

| File | Action |
|---|---|
| `delegation/contracts/DecimalDelegation.sol` | Add `autoUnbondEnqueue()` (Phase 1, `msg.sender==0`), `processAutoUnbond()` (Phase 2, permissionless), `_withdrawFor()` (internal, like `_withdraw` but with explicit delegator), `AutoUnbondEntry[]` storage |
| `delegation/contracts/interfaces/IDecimalDelegationCommon.sol` | No changes needed (reuses existing events) |

### Off-chain (optional, trustless)

| Component | Action |
|---|---|
| Any EOA / bot | Calls `processAutoUnbond(0)` when `autoUnbondQueueLength() > 0`. Permissionless, no trust. |

---

## 10. Testing Strategy

1. **Unit tests** (`x/validator/keeper/auto_unbond_test.go`):
   - Validator offline for < timeout → no auto-unbond
   - Validator offline for > timeout → all delegations force-unbonded
   - Validator comes back online before timeout → timer reset, no auto-unbond
   - Validator jailed then unjailed+online before timeout → no auto-unbond
   - Validator jailed (already offline) → timer doesn't reset
   - Mixed coin + NFT delegations → both handled correctly
   - `timeout = 0` → feature disabled
   - Genesis export/import preserves OfflineSince state

2. **Integration tests**:
   - Full lifecycle: create validator → delegate → go offline → wait → verify auto-unbond
   - Verify EVM force withdrawal is called
   - Verify correct events emitted

---

## 11. Phase 2 Bot Script: `processAutoUnbond()` Transaction Sender

### Overview

A standalone Hardhat script that monitors the auto-unbond queue on the EVM delegation contract and sends `processAutoUnbond()` transactions to drain it. The script is **trustless** — anyone can run it, and it only processes entries already enqueued by Phase 1.

### Location

`decimal-smart-contracts/delegation/scripts/process-auto-unbond.ts`

### Implementation

```typescript
import { ethers, network } from 'hardhat'
import { DecimalContractCenter, DecimalDelegation } from '../typechain-types'
import { contractCenters, Networks, getGasPrice } from '../../config'

const CONFIG = {
  // How long to wait between queue checks when queue is empty
  POLL_INTERVAL_MS: 60_000, // 1 minute
  // Delay between consecutive processAutoUnbond() txs (avoid nonce issues)
  TX_DELAY_MS: 5_000, // 5 seconds
  // Max retries for a single tx
  MAX_TX_RETRIES: 3,
  // Whether to run continuously or exit after draining
  DAEMON_MODE: true,
}

async function getDelegationContract(): Promise<DecimalDelegation> {
  const DecimalContractCenter = await ethers.getContractFactory('DecimalContractCenter')
  const contractCenter = DecimalContractCenter.attach(
    contractCenters[network.config.chainId!]
  ) as DecimalContractCenter

  const delegationAddress = await contractCenter.getContractAddress('delegation')
  const DecimalDelegation = await ethers.getContractFactory('DecimalDelegation')
  return DecimalDelegation.attach(delegationAddress) as DecimalDelegation
}

/// Scan queue and return indices of entries whose hold has expired (processable now).
/// Entries with active holds are skipped — they stay in the queue for later.
async function findProcessableIndices(delegation: DecimalDelegation): Promise<number[]> {
  const queueLen = Number(await delegation.autoUnbondQueueLength())
  const now = Math.floor(Date.now() / 1000)
  const processable: number[] = []

  for (let i = 0; i < queueLen; i++) {
    const entry = await delegation.getAutoUnbondEntry(i)
    const holdTimestamp = Number(entry.holdTimestamp)
    if (holdTimestamp === 0 || holdTimestamp < now) {
      processable.push(i)
    }
  }

  return processable
}

async function processOne(
  delegation: DecimalDelegation,
  index: number,
  networkName: Networks
): Promise<boolean> {
  for (let attempt = 1; attempt <= CONFIG.MAX_TX_RETRIES; attempt++) {
    try {
      const tx = await delegation.processAutoUnbond(index, {
        gasPrice: await getGasPrice(networkName),
      })
      const receipt = await tx.wait()
      console.log(
        `[OK] processAutoUnbond(${index}) tx=${receipt!.hash} block=${receipt!.blockNumber} gas=${receipt!.gasUsed}`
      )
      return true
    } catch (err: any) {
      if (err.message?.includes('Invalid index') || err.message?.includes('Hold not expired')) {
        console.log(`[SKIP] index=${index}: ${err.message}`)
        return false
      }
      console.error(`[ERR] index=${index} attempt ${attempt}/${CONFIG.MAX_TX_RETRIES}: ${err.message}`)
      if (attempt === CONFIG.MAX_TX_RETRIES) {
        throw err
      }
      await sleep(CONFIG.TX_DELAY_MS)
    }
  }
  return false
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

async function run() {
  const networkName = network.name as Networks
  console.log(`Starting auto-unbond processor on ${networkName}`)

  const delegation = await getDelegationContract()
  console.log(`Delegation contract: ${delegation.target}`)

  const [signer] = await ethers.getSigners()
  console.log(`Signer: ${signer.address}`)

  do {
    const queueLen = Number(await delegation.autoUnbondQueueLength())
    console.log(`Queue length: ${queueLen}`)

    if (queueLen === 0) {
      if (!CONFIG.DAEMON_MODE) {
        console.log('Queue empty, exiting.')
        return
      }
      console.log(`Queue empty, sleeping ${CONFIG.POLL_INTERVAL_MS / 1000}s...`)
      await sleep(CONFIG.POLL_INTERVAL_MS)
      continue
    }

    // Find entries that are processable (hold expired or no hold)
    const indices = await findProcessableIndices(delegation)
    console.log(`Processable: ${indices.length}/${queueLen} (${queueLen - indices.length} held)`)

    if (indices.length === 0) {
      console.log(`All entries held, sleeping ${CONFIG.POLL_INTERVAL_MS / 1000}s...`)
      await sleep(CONFIG.POLL_INTERVAL_MS)
      continue
    }

    // Process from highest index to lowest — swap-and-pop won't shift
    // lower indices when we remove a higher one first
    for (const index of indices.reverse()) {
      await processOne(delegation, index, networkName)
      await sleep(CONFIG.TX_DELAY_MS)
    }

    console.log(`Processed ${indices.length} entries`)
  } while (CONFIG.DAEMON_MODE)
}

run()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error('Fatal error:', error)
    process.exit(1)
  })
```

### How to run

```bash
# One-shot: drain queue and exit
cd delegation
npx hardhat run scripts/process-auto-unbond.ts --network dsc_main

# Daemon mode (default): continuously poll and drain
# Set DAEMON_MODE = true in CONFIG (default)
npx hardhat run scripts/process-auto-unbond.ts --network dsc_main
```

### Alternative: Hardhat task (for ad-hoc use)

Also register as a task in `delegation/tasks/processAutoUnbond.ts` for one-off invocation:

```typescript
import { task } from 'hardhat/config'
import { Networks, contractCenters, getGasPrice } from '../../config'
import { DecimalContractCenter, DecimalDelegation } from '../typechain-types'

task('process-auto-unbond')
  .addOptionalParam('max', 'Max entries to process (0 = all)', '0')
  .setAction(async ({ max }, { ethers, network }) => {
    const maxEntries = parseInt(max)
    const networkName = network.name as Networks

    const DecimalContractCenter = await ethers.getContractFactory('DecimalContractCenter')
    const contractCenter = DecimalContractCenter.attach(
      contractCenters[network.config.chainId!]
    ) as DecimalContractCenter

    const delegationAddress = await contractCenter.getContractAddress('delegation')
    const DecimalDelegation = await ethers.getContractFactory('DecimalDelegation')
    const delegation = DecimalDelegation.attach(delegationAddress) as DecimalDelegation

    const queueLen = Number(await delegation.autoUnbondQueueLength())
    console.log(`Queue length: ${queueLen}`)

    // Find processable entries (skip held stakes)
    const now = Math.floor(Date.now() / 1000)
    const processable: number[] = []
    for (let i = 0; i < queueLen; i++) {
      const entry = await delegation.getAutoUnbondEntry(i)
      const hold = Number(entry.holdTimestamp)
      if (hold === 0 || hold < now) processable.push(i)
    }
    console.log(`Processable: ${processable.length}/${queueLen}`)

    const toProcess = maxEntries > 0
      ? processable.slice(0, maxEntries)
      : processable

    // Process highest index first (swap-and-pop safety)
    for (const [n, index] of toProcess.reverse().entries()) {
      const tx = await delegation.processAutoUnbond(index, {
        gasPrice: await getGasPrice(networkName),
      })
      const receipt = await tx.wait()
      console.log(`[${n + 1}/${toProcess.length}] index=${index} tx=${receipt!.hash} gas=${receipt!.gasUsed}`)
    }

    console.log(`Done. Processed ${toProcess.length} entries.`)
  })
```

Register in `hardhat.config.ts`:
```typescript
import "./tasks/processAutoUnbond";
```

Usage:
```bash
# Process all pending
npx hardhat process-auto-unbond --network dsc_main

# Process max 10
npx hardhat process-auto-unbond --max 10 --network dsc_main
```

### Operational notes

- **Gas cost**: Each `processAutoUnbond()` call processes 1 entry. Gas cost is roughly equivalent to a single `withdraw()` call (~100-200k gas). The caller pays gas.
- **Who runs it**: Anyone. In practice, the Decimal team runs the daemon script. Delegators can also self-serve by calling `processAutoUnbond()` directly.
- **Failure handling**: If a `processAutoUnbond()` tx reverts (e.g., stake already withdrawn by another path), the entry is still popped from the queue. The next call processes the next entry.
- **No MEV risk**: The function takes no parameters. The caller cannot influence which entry gets processed or extract any value — data comes entirely from contract storage.
- **Concurrency**: Multiple callers can run simultaneously. Each call atomically pops one entry. No double-processing possible due to EVM transaction ordering.
- **Monitoring**: Check `autoUnbondQueueLength()` via any RPC/explorer. If it grows and stays non-zero for extended periods, the bot may be down.
