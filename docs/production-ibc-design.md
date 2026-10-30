# Production IBC integration

Decimal currently initializes IBC core, but app/app.go does not wire an ICS-20 TransferKeeper, a transfer application route or a transfer AppModule. Importing ibc-go and exposing clients/connections is therefore not sufficient to provide token transfers. This proposal defines the implementation and acceptance criteria for a production integration without silently changing token backing or validator security assumptions.

This is an architecture proposal. It does not enable transfers, create live clients/channels, change chain state or claim that relaying has already been tested.

## Baseline after the mainnet recovery

Reviewed against master `a46cb48d7965c1a64a495cebb22e735dfe4b9da3`, which schedules the delegation recovery at height 34342701. The recovery adds delegation-contract migration, node/contract reconciliation, legacy-event conversion, NFT exit pairing, hold-start backfill and reward synchronization. It does not add a TransferKeeper, an ICS-20 route or transfer module, and does not change the pinned SDK/IBC dependency family. The wiring and security work below is still required.

Use post-recovery state and the deployed delegation ABI as the integration baseline. Preserve recovery handlers and historical event behavior for replay. The recovery deliberately retains documented pre-incident delegation drift; an acceptance check must distinguish that baseline from new unexplained drift rather than require every historical node/contract difference to become zero. Reconcile balances, staking/hold liabilities and pending exits before introducing IBC escrow. Only spendable assets may be transferred; a held or staked claim must not become spendable through a voucher, wrapper or refund path.

ICS-20 receive, acknowledgement and timeout callbacks run through native application accounting. They must not rely on the EVM PostTxProcessing hook to mirror balances or enforce policy. Preserve the NFT multi-exit pairing and legacy-event fixes when adding a bank/EVM asset adapter, and exercise that adapter explicitly through both native IBC and EVM paths. NFT transfer remains outside this proposal's initial scope.

Treat reliability corrections in PR #1 as a separate change, not functionality already present in this master. Before production IBC activation, verify crash-before-commit upgrade replay, retry behavior for partial auto-unbond batches, and rollback on malformed financial events on the final integrated release. A successful recovery block alone does not prove those failure paths are safe.

## Initial delivery scope

The first production milestone should provide working IBC light clients, connection and unordered ICS-20 channel handshakes, native DEL transfers, receive/return paths, error acknowledgements and timeout refunds. Include a relayer runbook, monitoring and supply invariants.

Rate limiting should be included before opening significant transfer capacity. Packet forwarding, interchain accounts and incentivized relaying are subsequent capabilities with their own acceptance tests; a functioning ICS-20 application does not require every optional IBC application.

## Dependency and security prerequisites

The repository pins SDK 0.46.16, CometBFT 0.34.29 through a Tendermint replacement, and ibc-go/v6 6.2.1. That dependency includes transfer and ICA application packages, but their presence does not establish a supported production configuration.

Choose a supported SDK/CometBFT/IBC release family and verify compatibility with the Decimal Ethermint fork. An experimental v6 integration can establish local feasibility, but should not be presented as a recommendation to launch new production channels on an unreviewed legacy line. Record exact tags, fork patches, advisories and migration instructions.

Decimal replaces standard staking behavior with ValidatorKeeper. Verify every IBC-required staking interface and its semantics: historical consensus headers and validator sets, unbonding time, validator updates and evidence handling. Light-client trusting periods must be derived from actual unbonding/security settings, not copied from another chain. Test height revisions and the decimal_75-1 style chain ID explicitly.

## Application wiring

The implementation should include these changes using the APIs of the selected release family:

1. Add the transfer store key, params/genesis handling, module account, permissions, interface registration and application module.
2. Scope the transfer capability before CapabilityKeeper.Seal on release families that use capabilities. Do not reuse a general capability scope for unrelated modules.
3. Construct TransferKeeper with the correct channel/port keepers, bank/account keepers and selected ICS4 wrapper. Use the middleware stack as its sender when middleware must intercept sends.
4. Add the transfer application to the IBC router and install that router once construction is complete. Do not confuse the governance proposal router with the IBC application router.
5. Register module services, CLI, gRPC gateway, genesis import/export, simulation and required begin/end-block ordering.
6. Mount the new store in a coordinated upgrade and initialize its genesis/module version. Avoid reconstructing existing bank balances or IBC client state unnecessarily.
7. Add the transfer route to the transaction admission policy. Test actual messages against Decimal's custom ante handler; a registered Msg service alone does not make MsgTransfer admissible.

Conceptual stack:

```text
application send
  -> TransferKeeper
  -> selected outbound middleware wrappers
  -> IBC channel keeper

packet receive / acknowledgement / timeout
  -> IBC core verification
  -> selected inbound middleware
  -> ICS-20 application
  -> bank escrow, voucher mint/burn or refund
```

The middleware order must be documented for both send and receive paths. Do not copy constructor signatures from a newer IBC major version into the v6 application.

## Asset model and conservation

Define one canonical representation of each asset before accepting packets.

| Asset | Proposed initial policy | Required accounting checks |
| --- | --- | --- |
| Native DEL | Eligible for ICS-20 after explicit review of native bank/EVM balance equivalence | Outbound escrow, return unescrow, timeout refunds, native supply and spendable bank/EVM balances agree. |
| WDEL or another contract wrapper | Separate representation; no automatic equivalence with DEL | Wrapping/unwrapping and backing must remain conserved. A wrapper is not implicitly a second native asset. |
| Decimal custom coins | Excluded until each route and backing model is reviewed | Reserve-curve accounting, bank balances, DRC20 representation, symbol/address resolution and escrow agree. |
| Incoming IBC vouchers | Bank vouchers using the selected ICS-20 denomination-trace rules | Preserve full traces and channel provenance. Do not turn them automatically into reserve-backed coins or stakable assets. |
| NFTs | Outside initial ICS-20 scope | Separate protocol and custody design; ICS-20 is not NFT transfer. |

Document denomination units and metadata independently of display decimals. Test the effects of previous denomination changes: IBC wire amounts and escrow balances must use canonical base units. Avoid conversions through floating-point arithmetic.

Account for return paths: a voucher returning to its source is burned/unescrowed according to its trace, not minted as a fresh unrelated coin. The same ticker on different channels does not identify the same asset.

Add conservation checks for escrow, outstanding vouchers, mint/burn events and native supply. For representations that expose both bank and EVM balances, prove that one transfer cannot spend or create two independently backed balances.

## Transaction admission and authority

The current ante router admits selected native Cosmos messages and Ethereum extensions. Add explicit policies for transfer and IBC client/connection/channel messages. Decide whether ordinary users may create clients and channels or whether operations require a controlled workflow.

Validate every message in a mixed transaction. Do not admit an arbitrary transaction merely because one message is on the allowlist. Define how nested authz/group execution is inspected and authorized. Keep freeze/resume controls consistent across native, IBC and EVM paths.

Authz and feegrant are already registered in the application, but their presence does not prove they work through the custom ante policy. Test them before relying on fee sponsorship or delegated relayer operations.

Use separate operational accounts for relaying and emergency administration. Relayers must not have authority to replace light-client consensus roots, bypass proof verification or take custody of user escrow.

## Rate limits and optional applications

Rate limiting should use per-channel, per-denomination quotas with an explicit accounting window and governance policy. Test sends, receives, errors, acknowledgements and timeouts so failed transfers restore quotas correctly. Avoid policies that block the acknowledgement or refund required to settle an existing transfer.

Modern IBC rate limiting and packet-forward packages have moved across repositories/releases. Select the implementation matching the chosen IBC family rather than assuming the newest module is source compatible with v6.

Packet Forward Middleware requires multi-hop acknowledgement and timeout/refund tests. Establish what happens when the first hop succeeds and the second fails, and verify that source funds are not stranded.

ICA/ICS-27 should be added only after defining allowed host messages, controller policy and account authority. Enabling ICA does not automatically allow arbitrary EVM calls safely. Incentivized relaying/ICS-29 needs fee-asset and refund accounting compatible with the Decimal fee module.

## Required tests

These are proposed implementation acceptance tests, not tests already delivered by this documentation change.

| Scenario | Acceptance criterion |
| --- | --- |
| Two independent local chains, actual relayer | Client, connection and channel handshakes complete; packets are verified by the light client. |
| DEL outbound and return | Escrow and source balances match; the return path does not create a second asset or extra supply. |
| Foreign voucher outbound and return | Correct denomination trace, mint/burn and source provenance. |
| Invalid receiver / disabled asset | Error acknowledgement and exact refund; no stranded escrow. |
| Packet timeout | Exactly one refund; duplicate timeout/ack attempts cannot pay twice. |
| Process restart between send and acknowledgement | Commitments persist and the relayer resumes without duplicating transfer effects. |
| Upgrade with packets in flight | Client, channel, escrow and packet state survive migrations. |
| Recovery baseline and asset adapter | Documented historical drift remains classified; no new unexplained drift or double backing appears after sends, returns, refunds and stake/hold exits. Native callbacks enforce accounting without relying on an EVM transaction hook. |
| Recovery-era event compatibility | Legacy event conversion and multiple NFT exits retain their accounting behavior on the integrated release; IBC vouchers remain outside staking and NFT reserve accounting unless separately approved. |
| Duplicate packet or invalid proof | Rejected without side effects. |
| Expired or frozen client | Proof acceptance stops; recovery follows an explicit policy. |
| Validator-set change / evidence / unbonding | IBC verification remains consistent with the custom validator keeper. |
| Rate limit reached | New flow is restricted while permitted settlement/refunds remain possible. |
| Mixed and nested messages | Unapproved operations cannot enter through a permitted outer or sibling message. |
| Soft freeze and resume | IBC settlement policy is deliberate and cannot be bypassed through EVM/native routing. |
| Future middleware | Multi-hop and ICA tests run against actual supported counterparties. |

Use a version-compatible multi-chain framework, such as Interchaintest, rather than a mocked packet dispatcher as the sole acceptance test. Add state invariant assertions before and after each failure case.

## Rollout and operations

1. Agree on the dependency family, asset model, authority and exact activation upgrade.
2. Implement ICS-20 and conservation checks on isolated development chains.
3. Run testnet against at least one independent counterparty, including interrupted relaying and upgrades with packets in flight.
4. Publish supported client/channel configuration, denomination metadata, RPC/gRPC requirements, quotas and a relayer runbook.
5. Open mainnet channels with small limits after review; increase limits only after observing successful sends, refunds and returns.

Monitor client expiry/freeze, consensus height lag, relayer availability, pending commitments, acknowledgements, timeout eligibility, escrow totals, voucher supply and quota usage. Keep at least two independently operated relayers; this is an availability measure, not a substitute for proof verification.

## Completion criteria

The integration is complete when real clients and transfer channels are connected, native transfers and vouchers settle or refund correctly under the failure tests, supply invariants hold, node admission and custom staking interfaces are verified, and operators can upgrade and restart without losing packet state. A newly added dependency, a compiled TransferKeeper or a successful single happy-path transfer is not sufficient.

## References

- [IBC specifications](https://github.com/cosmos/ibc)
- [IBC-Go source and release notes](https://github.com/cosmos/ibc-go)
- [ICS-20 fungible token transfer specification](https://github.com/cosmos/ibc/blob/main/app/ics-020-fungible-token-transfer/README.md)
- [IBC application middleware](https://github.com/cosmos/ibc-apps)
- [Interchaintest](https://github.com/amygdala-labs/interchaintest)
