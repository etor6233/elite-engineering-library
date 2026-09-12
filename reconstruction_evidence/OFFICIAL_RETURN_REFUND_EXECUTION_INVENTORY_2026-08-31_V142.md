# Official Return Refund Execution — V142

Date: 2026-08-31  
Result: `REBUILD_VERIFIED / CONDITIONED`

## Scope and provenance

`GO-OFFICIAL-RETURN-REFUND-WORKER` 0.1.0 adds fourteen `AUTHORED` files. No local line is represented as copied Stripe, Mercado Pago or PostgreSQL source. The adapter invokes the exact official MIT SDK releases already pinned by Elite: `stripe-go/v86` 86.3.0 at commit `a2df585a800a97fe8ec4ebf551b4449bdb3d90a1` and `mercadopago/sdk-go` 1.14.0 at commit `f910ee53fbb6819e435eaf3d0f800cb1fe74ae09`.

The admitted claim is bounded and executable: a return with one exact allocated quantity-one order line and one unambiguous eligible payment can create and reconcile a partial Stripe or Mercado Pago refund. Amount and currency are server-derived. Ambiguous payments, missing mappings/configuration, divergent responses and unknown statuses block. Mercado Pago is limited to ARS/two decimals until another market has a separately admitted exponent contract.

## Reconstructed artifact

- pack round trip: 14/14 paths and SHA-256 byte-identical;
- library: 85 packs / 921 materialization blocks / 482 Markdown / 36 profiles;
- provenance: 779 `AUTHORED`, 37 `ADAPTED`, 105 `VERBATIM`;
- backend profile: 26 packs / 269 files composed from Markdown without collision;
- schema: PostgreSQL migrations 0001–0020 on PostgreSQL 18.6;
- module: frozen dependency graph, full `go test ./...`, `go vet ./...` and `cmd/return-refund-worker` build from Markdown;
- backend: full root Go tests/vet and all five actual main binaries build from Markdown;
- structural gate: `VERIFY_LIBRARY_PASS` reports 85 packs / 921 blocks / 482 Markdown and backend 269.

## Official provider contract evidence

Controlled local transports exercise the exact SDK clients and prove:

1. Stripe `V1Refunds.Create` sends the stable `Idempotency-Key`, exact PaymentIntent, amount, reason and return-request metadata; `Retrieve` reconciles the persisted refund ID;
2. Stripe response payment reference, refund ID, amount, currency and status must all match;
3. Mercado Pago retrieves the original payment through `payment.Client.Get`, verifies ID/currency/original amount and requires `approved/captured` before creation;
4. Mercado Pago `refund.Client.CreatePartialRefund` sends the exact ARS amount and caller-controlled `X-Idempotency-Key`; `Get` reconciles payment/refund IDs and amount;
5. retrieval accepts the official fully-refunded payment state while retaining all identity/amount/currency checks;
6. unsupported references, currencies, amounts, missing credentials and unknown provider states fail closed.

No real secret or provider account was used. These are deterministic SDK contract tests, not a claim of live account enablement.

## Real PostgreSQL evidence

The clean-database and down/up integrations prove:

1. preparation locks the return, exact allocated order line and candidate captured payments;
2. the first 1,000-minor-unit refund against a 2,000 payment remains partial and leaves payment `captured`;
3. the second confirmed 1,000 refund makes the exact successful sum 2,000 and atomically transitions payment to `refunded`;
4. two return-refund success events plus one payment-refunded event are written through the existing outbox;
5. provider observations and return-effect attempts remain immutable;
6. two eligible payments for the same returned line produce `REFUND_MAPPING_CONFLICT` and durable `blocked`, never an arbitrary choice;
7. completed and blocked terminal work is not reclaimed;
8. migration 0020 down removes only its objects, absence is verified, and up/SQL/integration pass again.

## Crash, retry and reconciliation contract

The provider call is outside the database transaction. The request's stable idempotency key survives every retry. A crash after remote acceptance but before local commit repeats the same create safely; after a provider refund ID is persisted, later claims retrieve instead of creating. Each provider response is hashed and appended before the mutable projection advances. Lease, attempt and claim token fence stale completions. Pending/action-required states remain retryable; exhaustion becomes explicit `blocked`.

## Governing official sources

- Stripe Go SDK exact source and refund service;
- Stripe Refund API and idempotent request contract;
- Mercado Pago Go SDK exact payment/refund clients and request-options implementation;
- Mercado Pago create/get refund APIs, refund conditions and transaction/refund status documentation;
- PostgreSQL 18 row locking, `SKIP LOCKED`, transactions, constraints and triggers.

## Failures converted into memory

`LIB-FAIL-1407` through `LIB-FAIL-1417` record path/glob discovery, Windows argv limits, PostgreSQL process-group shutdown, compositor signature drift, provenance-ledger enforcement and closing-check harness corrections. The PostgreSQL evidence accepts readiness only after status, port and client query agree in separate processes; the final server is stopped explicitly after closing all gates.

## Conditions not converted into PASS

No live Stripe or Mercado Pago credential/account, webhook delivery, provider sandbox refund, discount/tax/shipping allocation policy, accounting reversal, ARCA credit note, exchange shipment, CDN/WAF, IdP, load, offensive security, backup/PITR, deployment, rollback or business acceptance was executed against a selected target. This evidence admits an immediately materializable reusable refund owner, not a production franchise system.
