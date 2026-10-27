# Upgrade and EVM accounting reliability

Three failure paths can leave the node's accounting inconsistent with its committed state or with EVM contracts. This change removes local upgrade journals from consensus decisions and adds activation-gated handling for auto-unbond failures and financial event decoding.

## Reproduced failures and corrections

| Failure | Original result | Corrected result |
| --- | --- | --- |
| Process exits after an upgrade handler returns but before block commit | The JSON journal survives while EVM writes do not. Replay skips the pending migration. | UpgradeKeeper's committed state remains authoritative. Replay executes the lost EVM migration, then does not repeat it after commit. |
| Second auto-unbond enqueue reverts | The first EVM enqueue persists and the retry marker is removed. | The validator batch uses a CacheContext. All writes roll back on error and the marker remains. A successful retry commits queue writes and marker removal together. |
| Trusted contract emits truncated or empty RewardPerBlockUpdated data | The decoder error is ignored, or empty data bypasses decoding. Enabled defaults to false and the existing override is cleared. | Validate the event envelope and propagate decoding errors before node accounting changes. Ethermint marks the transaction failed and discards EVM storage writes, value transfer, logs and hook changes. |
| A diagnostic upgrade journal is overwritten with shorter JSON | Trailing bytes make the next Load fail. | Write a complete temporary file, sync it, close it and replace the old file. Sync the containing directory on supported Unix filesystems. |

The journal replacement is independent of the main replay fix. Adding O_TRUNC alone would not prevent journal-induced migration skips. Application startup no longer reads or writes the diagnostic journal.

## Real EVM regression tests

The tests instantiate the full DSC application, its actual Ethermint keeper, the registered Coin/NFT/Validator hooks, and on-disk LevelDB. They do not replace the EVM executor or financial hook with mocks.

`TestIntegrationUpgradeCrashRealEVM` runs a child test process. The child changes contract storage through CallEVM, completes BeginBlock, then calls os.Exit(83) without EndBlock, Commit or database cleanup. The parent reopens LevelDB and verifies that the contract migration is replayed once. A second reopen verifies that a committed migration is not repeated.

`TestIntegrationAutoUnbondRealEVM` uses a Solidity fixture whose second hold enqueue reverts. The first real EVM write survives in the legacy control. With the correction activated, the failed batch leaves zero queue writes, retains OfflineSince, and writes the two holds and base amount exactly once on retry.

`TestIntegrationEventRollbackRealEVM` signs Ethereum transactions calling a Solidity fixture that writes storage, receives a native value transfer, and emits malformed RewardPerBlockUpdated data. The activated path must fail specifically with ErrPostTxProcessing, not an unrelated VM error. Storage, transferred value and the existing reward override are checked after failure. A valid event is then required to commit successfully. Legacy controls confirm the original malformed-event behavior.

The Solidity fixtures are deliberately small ABI-compatible fault injectors, not the production delegation implementation. The tests verify execution, cache, hook and recovery behavior; they do not establish that every deployed Solidity version has identical queue semantics. Gas fees and the outer transaction ante path are not exercised: signed transactions enter ApplyTransaction directly with a zero gas price.

Validated against master 9e6c6d718d662083c4a524376a66d2c50bd4bc77 and stage 88dd556783e8dbb56c95ac12324122ae1f1d1aa8. The uncorrected algorithms fail the active regression expectations. After the changes, both the regression scenarios and package tests pass:

```sh
CGO_ENABLED=1 go test -mod=readonly -count=1 -v -run '^TestIntegration.*RealEVM$' ./app
CGO_ENABLED=1 go test -mod=readonly -count=1 -timeout 120s ./app ./x/validator/keeper ./x/upgrade ./cmd/config
```

Use a supported CGO compiler. Validation used Go 1.24.4 with a Windows GNU-targeting C compiler. Linux release CI must also run these tests; Windows success is not a Linux deployment certificate.

Fixture bytecode is checked in so Go tests do not require Node.js. To reproduce the fixtures:

```sh
cd app/testdata/reliability
npm ci --ignore-scripts
node compile.cjs
```

The compiler is pinned to Solidity 0.8.20, optimizer runs 200, EVM target Paris. Regeneration must produce identical JSON artifacts.

## Activation and deployment

`ReliabilityFixActivationHeights` is deliberately empty. Reviewers must select a coordinated height for each exact chain ID before releasing the consensus-changing auto-unbond and event rules. Do not populate this table from local flags, environment variables or wall-clock time. The tests verify both disabled and activated behavior.

Before activation:

1. Reconcile existing queue entries against failed batches from the old implementation. A new retry rule does not undo already committed partial enqueues, and replaying those blindly may enqueue them again.
2. Rehearse on a fork with the actual deployed delegation and validator contracts, including owner changes and ABI versions.
3. Run Linux CI and historical replay checks, and publish the selected height and checksummed binaries to operators.
4. Keep the upgrade handler history required for replay.

The corrected batch is atomic per validator but still unbounded. A large validator or a permanently invalid delegation can repeatedly roll back its whole batch. A follow-up should introduce bounded persistent jobs, explicit progress and diagnosis of poison entries. This patch is not an assertion that the queue is now scalable without limits.

## Remaining findings

These findings are separate from the three regression-tested fixes; their impact should not be presented as a demonstrated live exploit.

| Finding | Evidence | Suggested correction |
| --- | --- | --- |
| NFT supply invariant is a placeholder | x/nft/types/types.go returns broken=false unconditionally | Implement conservation checks and fault-injection tests; avoid presenting the existing route as effective coverage. |
| Upgrade downloader blocks BeginBlock on an external HEAD request without an explicit timeout | x/upgrade/abci.go and x/upgrade/upgrade.go | Move downloads outside consensus execution; use explicit deadlines, limits and retry state. |
| Failed download is not retried in-process after downloadStat is set | downloadStat is set before the asynchronous download completes | Track success/failure separately; use bounded retries outside BeginBlock. |
| Binary replacement unlinks the current executable before renaming its successor | x/upgrade/upgrade.go changeBinary | Prepare versioned binaries and switch through an external supervisor, retaining the previous executable. |
| Empty per-platform checksum list can reach hashes[0] | x/upgrade/abci.go and changeBinary index a JSON-derived slice | Validate schema, nonempty list, hash length and algorithm before indexing. |
| Upgrade test suite is largely commented out | x/upgrade/abci_test.go | Restore applicable tests; retain the new subprocess/EVM regressions in CI. |
| Emergency authority documentation contradicts implementation | x/validator/types/emergency.go says mainnet is disabled, but returns an address; the same address is used across environments | Confirm actual ownership/policy, separate environments, document authority and rehearse recovery. An address alone does not establish whether the signer is a single key or threshold account. |
| TPS view counts observations rather than their transaction values | app/tps_counter.go uses view.Count with M(n) | Use an appropriate sum/counter and test a known transaction batch. |
| Debug output in financial hooks | fmt.Println in x/validator/keeper/evm_hooks.go and contracts/utils.go | Use structured, level-controlled logs. |
| Migration changes contract ownership slots directly | app/delegation_contract_upgrade.go on stage | Treat as a privileged state migration; verify storage layout, owner restoration and failure atomicity. Do not remove historical handlers casually. |
| Dependency fork overwrites transaction post-processing diagnostics | Ethermint v0.22.13 x/evm/keeper/state_transition.go assigns res.PostTxError="test" | Correct in the dependency fork and add an upstream regression. The EVM tests above establish rollback using VmError and actual state, not this overwritten field. |

See [Dependency modernization](dependency-modernization.md) for the compatibility and release plan. New IBC functionality should be reviewed separately from these corrections.
