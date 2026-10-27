# Dependency modernization

The node is tied to an older generation of the Cosmos and EVM stacks. Its direct dependencies and build recipes need a coordinated modernization plan. Updating unrelated versions to their newest tags in one change would bypass the compatibility work required by custom staking, EVM hooks, migrations and asset accounting.

## Current baseline

Snapshot reviewed on 27 October 2026:

| Component | Repository baseline | Required decision |
| --- | --- | --- |
| Go | go.mod declares 1.21; pipeline downloads 1.21.4 | Select a supported toolchain and test it with the CGO cryptography and target platforms. Go supports a major release until two newer major releases exist; 1.21 is outside that window. |
| Cosmos SDK | 0.46.16 | Select a supported target release family and write the module/store/params migration sequence. |
| Consensus | Tendermint import path replaced with CometBFT 0.34.29 | Move with the SDK-supported CometBFT version; test ABCI, evidence, validator power and historical IBC headers. |
| IBC | ibc-go/v6 6.2.1 | Choose a compatible supported family before exposing transfer or other applications. Core dependency presence is not a complete IBC integration. |
| EVM integration | decimalteam/ethermint 0.22.13 | Audit fork changes and backports. Evaluate maintained Cosmos EVM separately; it is not a drop-in replacement for Decimal accounting. |
| Ethereum libraries | go-ethereum 1.13.12 | Upgrade together with the EVM integration and agreed chain rules. Test VM/gas behavior, signer rules and StateDB atomicity. |
| Numeric/protobuf stack | cosmossdk.io/math RC, errors beta, older protobuf/grpc families | Let the chosen SDK family define a tested compatible set instead of overriding transitive versions independently. |
| Build environments | Dockerfile uses golang:stretch; CI has older distribution recipes and Go download steps | Publish a supported platform matrix, pin builder/runtime digests, verify toolchain artifacts and retain reproducible build metadata. |

Version age is evidence of maintenance debt, not proof that a particular vulnerability is reachable in this binary. Associate advisories with exact module versions and reachable call paths before labelling a finding exploitable. Include the custom Ethermint and consensus fork changes in that assessment.

## Migration sequence

1. Preserve the current lock state, golden state snapshots, upgrade history and reproducible binaries. Replace CI dependency mutation (`go get ./...`) with download/verify and `-mod=readonly` builds.
2. Add source-aware vulnerability scanning, an SBOM, artifact checksums and a documented advisory/backport process. Track upstream patches against fork commits.
3. Establish Linux CGO CI, real-EVM rollback tests, multi-node upgrade/restart tests, and long historical replay checks before changing consensus dependencies.
4. Upgrade the toolchain and build infrastructure in a separately validated change. Do not assume successful unit compilation proves consensus equivalence.
5. Pick one mutually compatible SDK/CometBFT/IBC/EVM release family. Map removed APIs, custom ValidatorKeeper interfaces, state schemas, signing, fees and capability wiring.
6. Run state migration and contract compatibility rehearsals on isolated copies. Verify delegated amounts, active holds, pending redelegations, auto-unbond jobs, native/contract supply and validator power before and after migration.
7. Release through devnet and testnet, followed by a coordinated mainnet upgrade with explicit heights and checksummed binaries.

Keep old ABI support and migration handlers when required by historical replay. Removing generated bindings or legacy modules solely because they are large or old can break reconstruction of existing state.

## Sources for release decisions

- [Go release policy and history](https://go.dev/doc/devel/release)
- [Cosmos SDK releases](https://github.com/cosmos/cosmos-sdk/releases)
- [CometBFT releases](https://github.com/cometbft/cometbft/releases)
- [IBC-Go releases](https://github.com/cosmos/ibc-go/releases)
- [Cosmos EVM](https://github.com/cosmos/evm)
- [Ethereum Go releases](https://github.com/ethereum/go-ethereum/releases)
- [Go vulnerability management](https://go.dev/doc/security/vuln/)

Select and record actual release tags during the migration review. Development main branches and release-page headlines are not evidence of the versions deployed on another production chain.
