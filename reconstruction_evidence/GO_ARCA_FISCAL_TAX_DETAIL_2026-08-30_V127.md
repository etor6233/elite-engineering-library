# GO ARCA Fiscal Tax Detail — V127 reconstruction evidence

Date: 2026-08-30  
Scope: exact line-level VAT and other-tax persistence for the durable WSFEv1 owner; no production claim.

## Authority and correction

- Current ARCA WSFEv1 developer manual 4.6: `https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf`.
- Acquired official PDF: 4,081,946 bytes; SHA-256 `809a68756c06a5cfa011a5e6db850e402d4fb33156ad0cb0029259a358d5305f`.
- Exact generated interface inspection proved that `FECAESolicitar` requires `AlicIva[]` and `Tributo[]`, not only aggregate amounts. V126 therefore remained insufficient for an exact SOAP request even though its aggregate accounting invariant was valid.
- V127 stores every VAT identifier/base/amount and every other-tax identifier/description/base/rate/amount. Service validation and deferred PostgreSQL constraint triggers independently require their sums and bases to equal the invoice aggregates before commit. Detail rows are immutable after insertion.
- The eleven files remain local `AUTHORED` code governed by ARCA's public contract. No ARCA, Microsoft or PostgreSQL source is embedded or falsely attributed.

## Reconstructed artifact

- Pack: `GO-ARCA-FISCAL-ISSUANCE-API` 0.2.0.
- Pack: 74,507 bytes; SHA-256 `24464ea105d75999f034720b86a37fc53dcc215bfdd5e5d90715d489ee836746`.
- Provenance: 11 `AUTHORED`, 0 `ADAPTED`, 0 `VERBATIM`; library totals remain 691/37/105 over 833 blocks.
- Application: `GO-ELECTROMOBILITY-APPLICATION` 1.3.0, SHA-256 `f235f61056953b06aa51d645edd35c2976d1461d653d04fdeed5dad58df31d30`; its compatibility range now requires fiscal 0.2.x.
- Backend plan SHA-256: `0cdecb69d28d4395e4b44a9b84cd5486fb428d870d3070836f463ea69749ef71`.
- Exact composition: 20 packs / 176 implementation files plus `MATERIALIZATION_RECORD.md`; every expected/actual hash matched.

| File | SHA-256 |
|---|---|
| `db/migrations/0012_arca_fiscal_issuance.down.sql` | `5dc38808804f75fcca12f9a596abacd45280e1fa37c9ec446f9f9e493cdb6385` |
| `db/migrations/0012_arca_fiscal_issuance.up.sql` | `05d4c2fa4600c5ceeecdc255ae58c861f3e6f669ae32e87a0b23db7d9bcf33d1` |
| `db/tests/0012_arca_fiscal_issuance.test.sql` | `9ab6e906f314fcf17a1726d8a7449414fb4110fed48cf2037d403a1c5b90eed8` |
| `internal/fiscal/processor_test.go` | `c09eeeaefa54177bdb6e6c4912ac6501210feb315bd95b237f1ada7c2843a7f5` |
| `internal/fiscal/processor.go` | `77fe23edacdc2235e8a521e8447feee16336f612cd0315224f58a21c20a4a313` |
| `internal/fiscal/service_test.go` | `c0c3469cb4e56416bff7e9928bf924863852f884bb09e05846452998ab48bf69` |
| `internal/fiscal/service.go` | `542b70561bcf9f22dcd89a41413cbb58b28a9fb77234e3a56e2222a63a89b311` |
| `internal/platform/httpapi/fiscal_test.go` | `773065be65f37e34d71d749a939ecec520b2c92977501391f7783c92da771a71` |
| `internal/platform/httpapi/fiscal.go` | `76e2450f7a954ca34ee619a24b0ae9bd550b56d2e9ccbb72079dd2ea2b151e14` |
| `internal/platform/postgres/fiscal_integration_test.go` | `8366b7fcb0f9198d8b5bc207720cd5088ebd860cde9552c93cf8b6402841b511` |
| `internal/platform/postgres/fiscal.go` | `d10b4eb74b09c6dbc024964b3c573d723a50a82960dce80b8c41c6785675e722` |

## Executed gates

- Standalone pack materialized 11/11 from Markdown; composed backend materialized 20/176 with exact record.
- Official Go 1.26.7: formatting clean, full `go test -count=1 ./...` twice with PostgreSQL enabled, `go vet ./...` and `go build ./...`: PASS.
- PostgreSQL 18.6: migration 0012 down/up from the composed artifact, SQL invariant test, and exact fiscal integration test twice: PASS. The output contained two real `--- PASS` events; a prior no-test filter is not counted.
- Tests reject missing/mismatched VAT details, duplicate identifiers, inconsistent other-tax totals, detail mutation, aggregate mutation, wrong payment/order amounts and ambiguous resend. Replay reloads the exact persisted lines.
- `VERIFY_LIBRARY_PASS`: 81 packs / 833 materialization blocks / 36 profiles; backend 20/176 and web 6/80.

## Failures retained as lessons

- `LIB-FAIL-1308`: a prior .NET temp root was only a locked partial residue and could not be reused; a fresh official SDK was hash-verified in a new isolated root.
- `LIB-FAIL-1309`: aggregate tax totals were initially mistaken for an exact WSFE request; generated types exposed the missing arrays and forced this correction before building the SOAP bridge.
- `LIB-FAIL-1310`: V127 initially guessed materializer/compositor parameter names; their exact `param` blocks were read and the empty destinations were reused safely.
- `LIB-FAIL-1311`: a test preflight invoked global `go`, which is intentionally absent; all accepted gates use the exact isolated Go 1.26.7 executable.
- `LIB-FAIL-1312`: a focal `-run` pattern matched no tests and was rejected as evidence; the exact test name was then executed twice and the two PASS events asserted.

## Exact condition

This closes the representational gap required to build an exact `FECAESolicitar` request; it does not approve the values. Production remains fail-closed until a real project supplies an approved CUIT/certificate/relationship/point of sale, current ARCA parameter tables, tax/rounding/voucher rules approved by responsible professionals, exact SOAP adapter, homologation receipts, QR/PDF/notes/contingencies, load/recovery/security evidence and business/accounting acceptance.
