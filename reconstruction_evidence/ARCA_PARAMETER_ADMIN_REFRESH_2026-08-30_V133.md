# ARCA parameter administration and durable refresh — V133

## Scope and provenance

This evidence governs `GO-ARCA-FISCAL-ISSUANCE-API` 0.5.0 and `GO-ELECTROMOBILITY-APPLICATION` 1.5.0. The added Go and PostgreSQL are local `AUTHORED` implementation, not code copied from ARCA, Microsoft, PostgreSQL or Google. ARCA WSFEv1 remains the protocol authority; the pinned Microsoft generated-client/credential/adapter chain remains the exact transport evidence; PostgreSQL 18 governs transaction, lease and `SKIP LOCKED` behavior. No fiscal table, credential, certificate, statutory cadence or provider response was invented.

Official governing sources:

- ARCA WSFE documentation and manual: <https://arca.gob.ar/ws/documentacion/ws-factura-electronica.asp> and <https://arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG-v4-0.pdf>.
- PostgreSQL 18 explicit locking: <https://www.postgresql.org/docs/18/explicit-locking.html> and `SELECT` locking/`SKIP LOCKED`: <https://www.postgresql.org/docs/18/sql-select.html>.
- Microsoft BCApps pinned evidence remains recorded in the admission ledger; it governs fiscal document/posting separation, not authorship of this Go/SQL.

## Materialized capability

1. An OIDC-protected administrative API exposes separately permissioned configuration, explicit refresh, scoped read and approve/reject-with-reason operations. Every protected operation binds the verified tenant and allowed organization; strict JSON rejects unknown fields and trailing values.
2. A durable schedule is configured explicitly per tenant, organization, CUIT, parameter kind and optional voucher class. Cadence is bounded to 900–2,592,000 seconds; there is no inferred legal cadence.
3. The worker claims one due schedule with a bounded lease and PostgreSQL `FOR UPDATE SKIP LOCKED`, calls the existing UDS-backed parameter provider, stores an immutable snapshot, advances or defers the schedule and emits an outbox event.
4. Acquisition and approval remain distinct. Neither direct refresh nor scheduled refresh can call the decision method. A new decision requires a separately authorized subject and a 3–500 character reason. Historical decisions remain upgrade-compatible; every new insert is protected by a database trigger.
5. Empty provider-code arrays are normalized to `[]`, never to fabricated codes and never to JSON `null`.

## Reconstruction and gates

- Fiscal pack materialized exactly 25 files; application pack exactly 5. Their manifests, blocks and SHA-256 values matched.
- The enterprise backend composed from Markdown into a new directory as 24 packs / 232 files.
- PostgreSQL 18.6 applied migrations 0001–0014 to a new database and all 13 SQL contract files passed. The focused real integration proved exclusive claim, no second worker claim, snapshot persistence with zero provider codes, successful completion, reasoned decision, scoped reload and isolated cleanup.
- The exact composed artifact passed empty `gofmt -d`, `go test ./... -count=1` including PostgreSQL integration, `go vet ./...` and `go build ./cmd/...`; four commands were built, including the new `arca-parameter-worker`.
- The governing library gate returned `VERIFY_LIBRARY_PASS` with 83 packs, 877 materialization blocks, 472 Markdown files, backend 24/232 and web 6/80. The independent executable audit returned exit 0 and `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` across all 83 packs and 121 locked upstream sources; target/live gates remain explicitly skipped rather than falsely passed where accounts, cost authorization or project evidence are absent.
- Failures 1340 through 1356 are retained in the failure ledger. They cover path assumptions, formatter invocation, Windows process lifecycle, fixture/schema assumptions, pgx prepared statements, immutable cleanup, the production `nil`-slice persistence defect, closure verification and blocked temporary cleanup. None was hidden or converted into a false pass.

## Admission boundary

This is `REBUILD_VERIFIED / CONDITIONED`. It is executable code available for immediate composition, but it is not evidence of live ARCA readiness. A concrete project still needs a real CUIT, certificate/private key in an approved secret store, WSAA association, point of sale, homologation account, approved fiscal/tax policy, official error handling, accountant acceptance, load/recovery/security evidence and then production authorization. Scheduled acquisition never makes a snapshot statutory policy automatically.
