# Return Effect Execution and Inventory Consumer — V141

Date: 2026-08-30  
Result: `REBUILD_VERIFIED / CONDITIONED`

## Scope and provenance

`GO-RETURN-EFFECT-EXECUTION-WORKER` 0.1.0 adds eight `AUTHORED` files. No line is represented as copied from AWS, Microsoft, Stripe, PostgreSQL or ARCA. Their official documentation governs idempotency, claim/concurrency, exact-source return linkage, provider reconciliation and fiscal homologation respectively.

The admitted claim is deliberately bounded: requests V140 now have durable execution state, attempt/resume evidence and a real inventory worker. Refund, exchange, accounting and fiscal effects remain requested until separate real owners prove target-specific execution and reconciliation.

## Reconstructed artifact

- pack round trip: 8/8 paths and SHA-256 byte-identical;
- backend profile: 25 packs / 255 files from Markdown, no collision;
- schema: PostgreSQL migrations 0001–0019 on PostgreSQL 18.6;
- binary: `cmd/return-effect-worker` plus the four pre-existing mains build successfully;
- Go: full `go test ./...`, `go vet ./...`, focused real-database integration and all actual main builds pass.
- global library: `VERIFY_LIBRARY_PASS` reports 84 packs / 907 blocks / 481 Markdown and `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` reports 84 packs / 121 upstream sources.

## Real PostgreSQL evidence

The focused integration proves:

1. two concurrent workers use `FOR UPDATE SKIP LOCKED` and receive different requests;
2. wrong worker/token completion is rejected;
3. inventory state, version, outbox event, immutable attempt and successful execution projection commit atomically;
4. a completed request is not claimed again;
5. blocked execution can resume only with a durable operator record;
6. a stale claim cannot complete after re-claim/fencing;
7. `quarantine` and `restock` map to real authoritative stock states;
8. attempts and resume evidence reject update/delete;
9. migration 0019 down removes only its projection/audit objects, up backfills existing requests, SQL contract and integration pass again.

Failures found before admission were recorded as `LIB-FAIL-1390` through `LIB-FAIL-1406`; the material ones became direct regression coverage: nullable claim-token uniqueness, transaction cleanup, outbox UUID decoupling, URL construction and command/path discovery.

## Governing official sources

- AWS Well-Architected Reliability Pillar: persisted idempotency, duplicate delivery and concurrency control.
- PostgreSQL 18: row locking, `SKIP LOCKED`, transactions and trigger semantics.
- Microsoft Business Central sales returns: original-sale/exact-cost linkage and separated return effects.
- Stripe refunds and idempotent requests: stable keys, asynchronous statuses and event/retrieval reconciliation.
- ARCA electronic invoicing help/manuals: credentials, point of sale, document mapping and homologation remain external target authority.

## Conditions not converted into PASS

No refund provider, replacement shipment, accounting reversal, fiscal credit note, CDN/WAF, IdP, load, offensive security, backup/PITR, deployment, rollback or business acceptance was executed against a production target. This evidence admits a reusable local component, not a production franchise system.
