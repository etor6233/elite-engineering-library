# Franchise Network Administration — Reconstruction Evidence V118

## Scope, method and honest provenance

This milestone extends the existing `org.organization` hierarchy and `franchise.agreement` owner. It creates no duplicate organization, franchise, authorization or outbox module. All nine materialized files in `GO-FULFILLMENT-SERVICE-FRANCHISE-API 0.3.0` are local `AUTHORED` code licensed as workspace-owner code. No Go or SQL line is represented as Microsoft source.

The exact public design/test authority is `microsoft/BCApps` commit `31a860b527f0dc72c7a44a255d7e7d403cfa4789`; its acquired archive is 225,939,109 bytes with SHA-256 `e3151b040df39cada83d41fcf5d9e6cdff1d8fddf226934c1ad21a2bdebc7210`, root MIT license SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`. The reviewed authorities are:

| Exact path | SHA-256 | Governing concern |
|---|---|---|
| `src/Layers/W1/BaseApp/Inventory/Location/ResponsibilityCenter.Table.al` | `864e237ccb860bd3b919a1cda6f983855ac2383be3357b6cb30dccab105bcdbc` | bounded responsibility-center identity and relationships |
| `src/Layers/W1/BaseApp/CRM/Team/Team.Table.al` | `1583f382721c9e8a58d0f3f18d633b67eacf52881ddc56c5781cda6872acad91` | explicit team owner and lifecycle surface |
| `src/Layers/W1/BaseApp/CRM/Team/TeamSalesperson.Table.al` | `854e3f5648ad959026347c4e1dbc2d347984d353227755132750d5d20e361a9f` | explicit membership boundary rather than duplicated person state |
| `src/Layers/W1/BaseApp/Finance/Intercompany/Partner/ICPartner.Table.al` | `8b0e3ac7e99b217138fd39a57d8e9d9de3dbe48cd0bbb5b30763756947292da7` | partner identity/status and intercompany boundary |
| `src/Layers/W1/BaseApp/Finance/Intercompany/Outbox/ICOutboxTransaction.Table.al` | `b8ebd60f3e3ca5aa1b9eecbced944a33c94d31d198d6686fc5e28a3e6c5da011` | durable outgoing transaction state |
| `src/Layers/W1/BaseApp/Finance/Intercompany/Inbox/ICInboxTransaction.Table.al` | `23fbdfb3e8a783ed86496ec9381df0184242f98cc1b0e4db1851c711c05a82d2` | durable incoming transaction state |
| `src/Layers/W1/Tests/ERM/ERMIntercompany.Codeunit.al` | `4c2c125a2de2f2fce65b5ada37688bd862ab2eb104ba4f13204e6a6da026c181` | positive/negative intercompany workflow testing |

Google SRE reliable-launch and release-engineering guidance governs proportional gates and testing the exact packaged artifact. NASA systems engineering guidance governs the separate statements of component verification, integrated validation and target operational acceptance. These sources govern the method and claim boundaries; they do not make the local implementation Google-, NASA- or Microsoft-authored.

## Materialized capability

- `GO-FULFILLMENT-SERVICE-FRANCHISE-API 0.3.0`, pack SHA-256 `870b93149cae2402f7d9aca8ea05bf5eb0cce23f9d1d163f97bfcdd351360643`, materializes nine files.
- Migration 0008 adds optimistic versions to organizations/agreements, rejects pre-existing hierarchy cycles before admission, serializes hierarchy mutations per tenant, rejects unavailable parents/new cycles and serializes active-territory activation per tenant/territory.
- The existing service/repository/HTTP owner creates allowed network node types, applies explicit lifecycle transitions, prevents closing a parent with live children or agreements, scopes actions by verified tenant/organization/permission and emits the state change through the existing transactional outbox.
- Agreement activation requires an active franchisee, exact current state/version and a non-overlapping active date range. PostgreSQL constraint errors become the domain conflict contract instead of an unclassified client error.
- `GO-ELECTROMOBILITY-APPLICATION 0.9.0`, pack SHA-256 `fb16935e6db5d245d2577c21b7373aa0158c6686f6fb6302c3ed6a8d35e9f214`, advances only its verified compatibility/composition contract.
- Exact backend profile: 17 packs / 144 files. Library inventory: 75 packs / 770 files, of which 633 are `AUTHORED`, 32 `ADAPTED` and 105 `VERBATIM`.

## Exact executable evidence

From a new Markdown-only backend composition:

- official Go 1.26.7 archive SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`: every source was `gofmt`-clean; full `go test -count=1 ./...`, `go vet ./...` and builds of `./cmd/api` plus `./cmd/electromobility-api` PASS;
- official PostgreSQL 18.6 EDB archive SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`: an empty database accepted migrations 0001–0008 and every SQL invariant test with `ON_ERROR_STOP=1`;
- direct SQL tests prove a proposed hierarchy cycle and an overlapping active territory are rejected;
- repository integration proves organization creation, activation, stale-version rejection, protected parent closure, agreement scope and stale-version rejection, and transactional outbox cardinality;
- two simultaneous activations for overlapping agreements in the same territory produced exactly one success and one normalized conflict;
- 0008 down removed the two version columns/functions/triggers/indexes; 0008 up plus its SQL test and PostgreSQL integration passed again;
- the materializer reconstructed all nine pack files byte-identically to the tested staging tree;
- `VERIFY_LIBRARY_PASS` reported 75 packs / 770 materialized files / backend 17/144.

## Failures converted to reusable memory

`LIB-FAIL-1265` through `LIB-FAIL-1273` preserve the interface/test-double drift, bounded ID fixture, multi-directory formatter misuse, missing test-path suffix, unavailable race instrumentation, invalid concurrent fixture, unsupported PowerShell parameter, blocked compound cleanup and repeated `$LASTEXITCODE` misuse after a PowerShell script. None of those attempts is counted as proof; each correction was followed by a clean Markdown reconstruction and dependent gates.

## Honest production boundary

This proves the reusable local network-administration primitive. It does not prove employee/resource rosters, jurisdiction-specific franchise contracts, royalty basis/settlement, accounting/tax, financing, target IdP/CDN/WAF, production recovery, target load, offensive security, deployment/canary/rollback or business acceptance. Those capabilities and receipts remain explicit conditions; no local test is renamed production evidence.
