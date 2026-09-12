# Dapr Official Transactional Outbox + Go Revalidation — 2026-08-28 V1

## Outcome

`Dapr Runtime v1.18.3` is admitted only as `PINNED_CANDIDATE_CONDITIONED` for the transactional-outbox boundary. The exact signed release, full Go graph, offline verification, official fake/unit tests and a focused reachability scan passed. The new pack materializes 21 files: 13 exact Apache-2.0 upstream files and 8 local admission/acquisition files.

This is not production authorization. Dapr documents at-least-once delivery. Every project must select and prove a transactional/outbox-capable state store, pub/sub component/account, consumer idempotency, duplicate delivery, recovery and rollback.

`Dapr Go SDK v1.14.2` is rejected for immediate adoption. Its official graph contains two vulnerabilities that govulncheck traced as reachable. The signed `main` commit observed on 2026-08-28 still pins an affected gRPC version. No dependency was locally bumped and no local result is called official Dapr code.

## Official identities

| Source | Official identity | Signature | Exact artifact |
|---|---|---|---|
| Dapr Runtime | release `v1.18.3`, commit `2dcf2e3548faf3c740b3d29ee300e2132db65cff` | GitHub verification `valid` | ZIP 20,269,733 bytes; SHA-256 `dd64436b01c06bdf59db6c7c949dce81eff36c293581980af93e12a80bfd6681` |
| Dapr Docs | branch `v1.18`, commit `5958a7e19a04e199326a6c5321bdd7714ee83b4c` | GitHub verification `valid` | outbox Markdown 26,022 bytes; SHA-256 `60487ccdd57bd7c308c0fc87f5e48e7843fb09dbb723e1532a0d8d412672e385` |
| Dapr Go SDK rejected release | `v1.14.2`, commit `2523198e6f8a3ab51fd5677dd525710060d4ce2d` | GitHub verification `valid` | ZIP 8,100,135 bytes; SHA-256 `1e3d46cbe820f3046cee4a06740824cca2a1b6ef7f7987ddba3decbe84e73bd8` |
| Dapr Go SDK rejected main | commit `9d745f76a46e23ae0d022f1abb098132b67603e5` | GitHub verification `valid` | `go.mod` observed with `github.com/dapr/dapr v1.18.0` and `google.golang.org/grpc v1.80.0` |

Official sources: [Dapr Runtime](https://github.com/dapr/dapr), [Dapr Go SDK](https://github.com/dapr/go-sdk), [Dapr outbox documentation](https://docs.dapr.io/developing-applications/building-blocks/state-management/howto-outbox/), [Dapr security advisory GHSA-85gx-3qv6-4463](https://github.com/dapr/dapr/security/advisories/GHSA-85gx-3qv6-4463).

## Runtime bytes retained

The pack preserves these exact runtime files from commit `2dcf2e3...`:

| Path | Bytes | SHA-256 |
|---|---:|---|
| `LICENSE` | 11,412 | `9523eacf5421b65637420755c00b5175f7804b6a6eb736b86cd10ea6b3937dcf` |
| `pkg/outbox/outbox.go` | 1,151 | `a0d48867049b7404e51852c8d688d7d9ca918cf2eecd9fc8a686636c783a985c` |
| `pkg/outbox/fake/fake.go` | 2,838 | `bf653d08463f8ca3f8a10658cf4bb1d661e7f244ef3cd8852574932e9e0545bd` |
| `pkg/outbox/fake/fake_test.go` | 696 | `cac9dfe14be633bdf1034c7b24ca133415f57a8c17338023708a8ab6ba1a4f61` |
| `pkg/runtime/pubsub/outbox.go` | 10,717 | `8a0a060481e107f3f0617fc1d6baf86c602105d3ef29c429a83abbc0c04923fc` |
| `pkg/runtime/pubsub/outbox_test.go` | 29,271 | `fac77524bbbbe8e2a3449ba7ec062f8979a1760085cf56e598acd8670cdd3cc4` |
| HTTP integration `basic.go` | 4,701 | `ce9ab635518713343b8e3db3fe98b12357a55086102ae2343fe80d1c34078bfe` |
| HTTP integration `delete.go` | 4,787 | `222dd239f765a04accac212563e56b200884cb77246816388ec80521efd2dbfb` |
| HTTP integration `projection.go` | 5,017 | `4e49a51289c0a2f809f496ef7a78e72e963a3d2b5211cbcdcfd246101392c216` |
| gRPC integration `basic.go` | 3,398 | `42e5524ebe243896d893a00a86afe2c9846a8642f6bf6b9c514467613d05eab3` |
| gRPC integration `projection.go` | 3,843 | `45849b9bc17678967886a3249b861097d702a14d5e99aa6bbb468d055b2d31bf` |

Dapr Docs license SHA-256 is `0b9cab20a5e2ae7e44f40a5ee6b8416f12d2135a547f9fef00e5b61f8d5be99a`. Runtime `go.mod` SHA-256 is `4f1941ca734128e36ae3c059d862b9d92759dab63de5f860b4e3a525c0998d9d`; `go.sum` SHA-256 is `e7c4c2dba67d067423e9cd85a8d13da8a86028cfd81d5b831a74b2ecb8c95e79`.

## Complete runtime execution

Toolchain: official Go `1.26.7` Windows amd64, ZIP 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`.

- Full module graph: 1,338 modules.
- `go mod download all`: PASS once against official proxy/sumdb.
- `GOPROXY=off go mod verify`: `all modules verified`.
- `go test ./pkg/outbox/... -count=1 -v`: interface package has no tests; official fake `Test_Fake` PASS.
- Focused `pkg/runtime/pubsub` regex selected exactly six main tests: `TestNewOutbox`, `TestEnabled`, `TestAddOrUpdateOutbox`, `TestPublishInternal`, `TestSubscribeToInternalTopics`, `TestOutboxTopic`.
- Those six main tests and twenty subtests PASS, including projection, overridden/default metadata, JSON content type, missing store, no transactions, pub/sub failure, discard and namespaced topic behavior.
- `govulncheck v1.7.0 ./pkg/runtime/pubsub`: `No vulnerabilities found.` DB update timestamp: `2026-08-27T19:52:28Z`.

The empty focused result is dated evidence, not a future-security guarantee. It does not replace whole-graph monitoring, project component tests or a later reachability scan.

## Rejected predecessor and SDK

### Dapr Runtime 1.18.2

Signed release `v1.18.2`, commit `d2dd3747145e38e4d85db4714e26d7569c9c37db`, passed the same focused tests but govulncheck traced reachable `GO-2026-6061` from `pkg/runtime/pubsub/default_bulkpub.go` through `google.golang.org/grpc@v1.80.0`. It is retained as `UP-FAIL-125` and rejected. Runtime 1.18.3 uses gRPC `v1.82.1` and was revalidated instead of locally modifying 1.18.2.

### Dapr Go SDK 1.14.2

The signed SDK release fixed a 504-module graph and `go mod verify` passed offline. Its official `TestStateTransactions` exercised six branches successfully, but the package suite closed non-zero on Windows because upstream `TestMain` unconditionally binds Unix socket `/tmp/dapr.socket`; this partial result is not promoted.

More importantly, govulncheck v1.7.0 found two reachable vulnerabilities:

- `GO-2026-6061` in `google.golang.org/grpc@v1.78.0`, fixed in `v1.82.1`;
- `GO-2026-5242` in `github.com/dapr/dapr@v1.17.0`, fixed in `v1.17.5`.

The signed main commit inspected later still pins gRPC `v1.80.0`. Therefore both release and observed main are rejected; no WSL rerun or local dependency change can make that upstream version admissible. A later official Dapr SDK must repeat commit/signature, archive/license, full graph, offline suite on a supported OS and reachability scan.

## Pack and acquisition gates

`implementation_packs/DAPR_OFFICIAL_TRANSACTIONAL_OUTBOX.md` materialized 21/21 files from Markdown.

- 13 embedded upstream files matched exact SHA-256.
- Runtime markers `outbox.projection`, `PublishInternal`, `TestNewOutbox` and `TestPublishInternal` were present.
- Incomplete project selection was rejected.
- Acquisition negative suite PASS: approval false, wrong lock binding, network not authorized and destination already present.
- Positive cache acquisition PASS from the exact 20,269,733-byte archive; root commit, license, `go.mod` and `go.sum` were verified before move.
- Source-lock regression PASS: 90 sources and nine negatives; Dapr classification/evidence are asserted.
- Composition profile materializes 21 implementation files plus record.

## Residual conditions

This evidence does not prove a project store, broker, account, component configuration, cross-process transaction, duplicate consumer behavior, load, restart, network partition, backup/restore, canary or rollback. Those are deliberately required by `project-selection.template.json` and remain false in the template. No credentials or infrastructure defaults are embedded.

`UP-FAIL-125` and `UP-FAIL-126` remain open upstream conditions. Local failures `LIB-FAIL-517` through `LIB-FAIL-530` retain the exact investigation mistakes and their regressions; their statuses are closed only after the final global rebuild.
