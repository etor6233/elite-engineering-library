# Franchise Appointment Resources — Reconstruction Evidence V119

## Scope and provenance

V119 extends the single existing appointment/calendar owner. All changed Go/SQL is local `AUTHORED` workspace-owner code; no line is represented as Microsoft source. Exact governing authority remains `microsoft/BCApps@31a860b527f0dc72c7a44a255d7e7d403cfa4789` (archive SHA-256 `e3151b040df39cada83d41fcf5d9e6cdff1d8fddf226934c1ad21a2bdebc7210`, root MIT license SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`):

- `Resource.Table.al`, SHA-256 `4bed76e68bafe002d0f3c6715dd437833bb567660f672acb78414196b714466a`;
- `ResCapacityEntry.Table.al`, SHA-256 `13a442f84ffa90c38d2308271b9449358d4cf0909c5d0a17cc9cf69cb0954cbe`;
- `ResourceSkill.Table.al`, SHA-256 `92c66766acb27f99003f722dbf3e278bdf80704e4b36761fb87e4563f534dd82`;
- `ServiceResourceSkill.Codeunit.al`, SHA-256 `7e76a8aafc78520e9fd99953ad030fd712f26dafe878bcbe9b44552ea966ffd7`.

Google SRE governs proportional launch/release gates and exact-artifact testing; NASA governs separation of component verification, integrated validation and target acceptance. Those references govern method, never authorship.

## Materialized result

- `GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.4.0`, SHA-256 `c56936d669e72df02f5575b4a766f8f21a2d9d4f7666a95d3f4efa61a98a3cd7`, materializes 18 files.
- Migration 0009 creates organization-scoped resources, human/non-human identity invariants, explicit appointment skills and one resource assignment per appointment.
- Assignment serializes per resource, requires active same-organization skill and rejects overlapping requested/confirmed appointments.
- Confirmation requires assignment; completion/no-show cannot precede start; cancellation releases slot capacity through the existing state predicate. Version and outbox remain transactional.
- Protected HTTP operations use `resource:manage` and `appointment:manage` plus verified organization scope; server-owned ID/status/version are not accepted as request fields.
- `GO-ELECTROMOBILITY-APPLICATION 1.0.0`, SHA-256 `1c31e3a9654d6f56fee1d6127e6a861958505bf11cc87fb13ad7dd0704388b41`, advances the exact composition contract.
- Current inventory: 75 packs / 773 files; backend 17/147; provenance 636 `AUTHORED`, 32 `ADAPTED`, 105 `VERBATIM`.

## Exact executable proof

- Markdown-only backend composition produced 147 files from 17 packs.
- Official Go 1.26.7: full tests, vet and both command builds PASS.
- Official PostgreSQL 18.6: empty database, migrations 0001–0009 and every SQL invariant test PASS with `ON_ERROR_STOP=1`.
- Live repository integration proved skill mismatch rejection, cross-scope protection, stale version rejection, future completion rejection, required assignment, cancellation and overlapping-resource rejection.
- 0009 down/up/test plus PostgreSQL integration PASS.
- Global library gates are recorded after canonical documentation/count alignment.

## Failure memory

- `LIB-FAIL-1274`: the first compiler run exposed a missing `current` argument in the repository interface; the signature, fakes and call sites were aligned before clean rebuild.
- `LIB-FAIL-1275`: a gate tried to use a not-yet-created composition as its workdir; composition and execution were separated and the exact 17/147 artifact then passed.
- `LIB-FAIL-1276`: `NULLS NOT DISTINCT` would have limited an organization to one non-human resource; it was replaced by a partial unique index applying only to non-null human principals and guarded by a SQL regression.

## Production boundary

This proves reusable resource identity, skills, appointment assignment and lifecycle primitives. It does not prove working hours, leave, labor/privacy law, notification delivery, target IdP/edge, production recovery/load/security/deploy or business acceptance. Royalties/settlement, fiscality/accounting and financing remain separate unfinished capabilities.
