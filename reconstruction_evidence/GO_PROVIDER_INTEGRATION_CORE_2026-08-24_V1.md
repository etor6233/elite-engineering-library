# Go Provider Integration Core — reconstruction evidence V1

Date: 2026-08-24  
Status: `REBUILD_VERIFIED / CONDITIONED`

## Result

The enterprise backend now includes a provider-neutral integration boundary with:

- connection-based tenant resolution; no tenant is trusted from webhook input;
- secret references separated from connection metadata and startup rejection of inline/missing/short secrets;
- an exact-body HMAC-SHA256 reference verifier with constant-time comparison and bounded replay window;
- strict and size-bounded JSON admission;
- transactional provider-event deduplication and durable job enqueue;
- stable internal/external resource mappings;
- discrepancy evidence and explicit reconciliation resolution;
- an application composition route ready for provider-specific adapters.

## Reconstruction

- `GO-PROVIDER-INTEGRATION-CORE 0.1.0`: 11 files independently materialized with SHA-256 verification.
- `GO-ELECTROMOBILITY-APPLICATION 0.6.0`: reference registry, verifier, PostgreSQL repository and HTTP module wired into the executable.
- compatible backend profile: 15 packs, 92 implementation files plus `MATERIALIZATION_RECORD.md`.
- clean target: `${TEMP}/elite-enterprise-profile-20260824-v10`.

## Gates

- migration `0004` down/up against PostgreSQL 18.6: PASS.
- all Go tests with PostgreSQL integration enabled: PASS.
- signature valid/tampered/expired, strict shape and secret-reference tests: PASS.
- first webhook receipt and atomic job creation: PASS.
- identical replay: PASS without a second event or job.
- repeated provider event ID with different body hash: conflict PASS.
- mapping upsert and reconciliation record/resolve: PASS.
- live HTTP delivery through the composed executable: `202`, replay `200`, invalid signature `401`; database counts `1 event / 1 job`.
- `gofmt -d`: PASS, no diff.
- `go vet ./...`: PASS.
- both Go executable builds: PASS.
- portable CI runner: 6 tests PASS.
- operational readiness validator: 6 tests PASS.
- packaging validator: 3 tests and `PACKAGING_VALID` PASS.

## Conditions

The HMAC envelope is a verified reference contract, not a claim about every vendor. Each selected payment, marketplace, advertising, carrier or communication provider still needs its official signature/authentication adapter, sandbox contract tests, outbound idempotency semantics, rate/quota handling, data-retention classification, reconciliation schedule and operational evidence. Credentials and production side effects are intentionally absent from this library.
