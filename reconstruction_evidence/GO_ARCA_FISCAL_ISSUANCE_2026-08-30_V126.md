# GO ARCA Fiscal Issuance — V126 reconstruction evidence

Date: 2026-08-30  
Scope: durable local owner for ARCA WSFEv1 issuance; no production claim.

## Authority and identity

- Current ARCA WSFEv1 developer manual 4.6: `https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf`.
- Acquired PDF: 4,081,946 bytes; SHA-256 `809a68756c06a5cfa011a5e6db850e402d4fb33156ad0cb0029259a358d5305f`.
- The manual explicitly governs the three operations implemented as orchestration invariants: `FECompUltimoAutorizado`, `FECAESolicitar` and, after an uncertain communication outcome, `FECompConsultar` before any resend.
- WSAA certificate/relationship authority remains `https://www.arca.gob.ar/ws/documentacion/wsaa.asp`.
- Microsoft BCApps authority is pinned at commit `31a860b527f0dc72c7a44a255d7e7d403cfa4789`; it governs the separation between source document, VAT/accounting posting and reversal already recorded in V121. No Microsoft or ARCA source is embedded in this pack.

## Reconstructed artifact

- Pack: `GO-ARCA-FISCAL-ISSUANCE-API` 0.1.0.
- Pack SHA-256: `39326ca1d88505bf972e4515b5ac89c4a111c761823066c68c0312327cc75c0d`.
- Provenance: 11 `AUTHORED`, 0 `ADAPTED`, 0 `VERBATIM`.
- Application pack: `GO-ELECTROMOBILITY-APPLICATION` 1.3.0, SHA-256 `cf171340338a1abf406f00c1ad24a89f7db504c159493c0d1df92ff4ad624d64`.
- Backend plan SHA-256: `d9c34b7549acf1c7ca2859063e91c780d82d62e5d050bab6183ae40fbd5139f0`.
- Exact composition: 20 packs / 176 files; manifest and every recorded hash matched.

| File | SHA-256 |
|---|---|
| `internal/fiscal/service.go` | `aceed10fbf512aa4007588d0802147ba368d9c1ce68cc42809ffe4534e489d44` |
| `internal/fiscal/service_test.go` | `691368be42776d5dfe838902bea6e55c6e9fda79f7be5705a04a3c92ef76ed4c` |
| `internal/fiscal/processor.go` | `77fe23edacdc2235e8a521e8447feee16336f612cd0315224f58a21c20a4a313` |
| `internal/fiscal/processor_test.go` | `c09eeeaefa54177bdb6e6c4912ac6501210feb315bd95b237f1ada7c2843a7f5` |
| `internal/platform/postgres/fiscal.go` | `e945c1fc5ca95a69cb735ac28f7a79bb3d890458ebb9a2e3447044306b2a80f6` |
| `internal/platform/postgres/fiscal_integration_test.go` | `a9ee5533afeb7efae09191043040cf5f3fccf612703bbb489b664bd6e9dff526` |
| `internal/platform/httpapi/fiscal.go` | `76e2450f7a954ca34ee619a24b0ae9bd550b56d2e9ccbb72079dd2ea2b151e14` |
| `internal/platform/httpapi/fiscal_test.go` | `f1cb5f07c935d8dff7de0b8fe34cf04c4d1f3b5a85e761eec08f07ac1c67f39f` |
| `db/migrations/0012_arca_fiscal_issuance.up.sql` | `e8699a98923a88c8c221be247cc72f764dc8269996de3025b8239ae056f55415` |
| `db/migrations/0012_arca_fiscal_issuance.down.sql` | `5dc38808804f75fcca12f9a596abacd45280e1fa37c9ec446f9f9e493cdb6385` |
| `db/tests/0012_arca_fiscal_issuance.test.sql` | `192bd51fa6c180476d6e60edc3d10fe6a347af904ac7db12bd72060c1b6fed48` |

## Executed gates

- Toolchains: official Go 1.26.7 windows/amd64; PostgreSQL 18.6 x86-64 Windows with page checksums; PowerShell 7.6.4.
- Fresh database applied migrations 0001–0012 and every `db/tests/*.test.sql` with `ON_ERROR_STOP=1`: PASS.
- Migration 0012 down/up plus its SQL invariant test: PASS.
- Full composed backend with `TEST_DATABASE_URL`: every Go package PASS; fiscal repository integration executed against PostgreSQL, not skipped.
- The focal PostgreSQL test repeated twice after down/up: PASS.
- Full `go test ./...` twice, `go vet ./...`, `go build ./cmd/api` and `go build ./cmd/electromobility-api`: PASS.
- Fiscal tests prove CUIT checksum, amount/date boundary, exact-body idempotency, trailing-JSON rejection, organization authorization, payment/order/amount binding, replay, single active lane, exact assigned number, timeout to `reconcile_required`, consult-before-resend, immutable attempt history, CAE persistence and one authorized outbox event.
- PostgreSQL was stopped cleanly after the final run.

## Failures retained as lessons

- `LIB-FAIL-1304`: the first inspection guessed a compositor path; corrected by materializing and reading its exact manifest.
- `LIB-FAIL-1305`: the first test command selected a not-yet-created workdir; corrected by creating/copying first and entering it explicitly.
- `LIB-FAIL-1306`: the standalone materializer was assumed to emit `MATERIALIZATION_RECORD.md`; only the compositor does. Verification now uses direct hashes for a single pack and the compositor record for profiles.
- `LIB-FAIL-1307`: a long catalog-row patch used semantically equal but byte-different wording; the rejected atomic patch was followed by exact-line lookup and a minimal successful update.

## Exact condition

This cycle proves reusable local workflow code, not a completed Argentine fiscal deployment. Production remains fail-closed until a real project supplies an approved CUIT, X.509 certificate and ARCA relationship, homologation point of sale, current parameter tables, exact generated WSFE SOAP adapter, tax rules/rounding and voucher-class matrix approved by responsible professionals, CAE/QR/PDF delivery, notes and contingencies, load/recovery/security evidence and business/accounting acceptance. No token, sign, certificate, private key or production payload was used or stored.
