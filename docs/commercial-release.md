# Durable commercial release checkpoint

This owner persists a commercial release receipt after the existing handover has been accepted, its checklist completed, and its bound payment observation and reservation revalidated. The write has exactly three durable effects: one immutable receipt, one completed idempotency record, and one outbox event `commercial-release.recorded`. It does not modify money, stock, order delivery state, tax treatment, inventory costs, a carrier assignment, or a physical shipment journal.

The receipt is historical evidence. `recorded_at` records when its dependencies were checked inside the write transaction; `valid_until` is the earlier of observation expiry and reservation expiry. Neither the existence of a receipt nor an idempotent replay asserts present eligibility. A later callback may invalidate the receipt's current evidence while its history remains recoverable. `ValidateCommercialRelease` returns `current=false` after a hold, new observation generation, changed evidence/policy/acceptance, or expiry.

## Explicit profile activation

The existing `READ_ONLY_ELIGIBILITY` effect and algorithm revision 1 retain their behavior. They cannot call the durable release owner. To select the supported persistence effect, materialize the same bound handover profile with the additional option:

```text
--release-effect commercial-receipt
```

This emits algorithm revision 2 and `release_effect: COMMIT_COMMERCIAL_RELEASE_RECEIPT`. Tenant, organization, provider, account, connection, mode, profile revision and raw-document SHA-256 remain mandatory bindings. No new credentials or environment variables are introduced. The materializer remains disabled unless `--activate` is explicitly provided. A profile with revision 2 and the old read-only effect is rejected, as is the new effect under revision 1. Other business policy algorithms remain unsupported and fail closed.

Prepare the initial handover using that exact profile from the beginning: immutable preparation binds the profile hash, so changing only the release-stage profile cannot upgrade an existing preparation silently.

## HTTP composition

Mount `httpapi.CommercialReleaseModule{Service: service}` beside `InitialHandoverModule`, using the same `*franchisejourney.HandoverPreparationService` produced by the admitted loader. The optional module registers no routes when its service is absent. The existing host helper returns an `InitialHandoverModule` whose concrete service implements `httpapi.CommercialReleaseService`; the host may assert that interface and mount this module. No new database or SDK owner is constructed.

All routes require authenticated `handover:manage` permission and membership in the specified organization. Tenant and actor come from verified claims. The path identifies the handover; order, customer, payment reference, amount and stock are derived by the repository.

| Method/path | Inputs | Result |
|---|---|---|
| `POST /v1/franchise/handovers/{id}/commercial-release` | `Idempotency-Key`; JSON `organization_id`, `observation_sha256` | Immutable receipt; 201; replay header when recovering an identical command |
| `GET /v1/franchise/handovers/{id}/commercial-release-result` | `Idempotency-Key`; query `organization_id` | Historical receipt only |
| `GET /v1/franchise/handovers/{id}/commercial-release-current` | Query `organization_id` | Receipt plus current revalidation result and evaluation time |

Responses use `Cache-Control: no-store`. Unknown command fields, including forged customer, order, actor or amount, are rejected. A second idempotency key cannot create a second release for the same order; replaying a changed actor or command conflicts. A failed write leaves no receipt, idempotency record or release event.

## Transaction boundary and downstream effects

`lockCommercialRelease` reuses `lockInitialHandoverScope`: active bound provider connection, order, payment attempt, provider observation, customer, order line, serialized stock and reservation are held stable. Acceptance is checked under a handover share lock with NOWAIT to avoid inversion with the existing acceptance owner. Freshness is checked after these locks and again after writing the outbox, before committing. The receipt records the observation hash/generation, accepted handover version/evidence, checklist version and profile hash.

The `commercial-release-current` route is a read-only view. A subsequent effect cannot use that response as a durable authorization token. An effect implemented in the PostgreSQL owner must call the same locked checks and compare the receipt's generation, evidence, acceptance and deadline within its own transaction before writing its journal. No such physical consumer is claimed here. Existing Supply shipment posting requires a registered warehouse pick, UOM/FIFO cost and shipment facts; Fulfillment transport requires that actual shipment and a carrier. This receipt does not fabricate them.

## Provenance and evidence

All new code in this owner is `AUTHORED` composition, configuration, persistence, HTTP and test glue. It is not attributed to Microsoft, Stripe, Mercado Pago or any other company. No foreign source was copied and no dependency changed. Existing official SDKs remain pinned in their separate admitted owner; the integral test uses those SDKs through a local HTTP fixture, real callback authentication, durable inbox/queue and PostgreSQL reconciliation.

Existing reused owner paths are `internal/platform/postgres/handover_preparation.go`, `internal/platform/postgres/franchisejourney.go`, `internal/paymentbridge/checkout.go`, `internal/platform/postgres/payment_callback_processor.go` and the materialized profile loader. `sales_shipment.go` and `fulfillment.go` were inspected to establish their additional physical-effect prerequisites, and are not invoked by this receipt.

The candidate gates prove local infrastructure behavior with fixtures and hash-bound profiles. They do not prove production access, fiscal issuance, shipment completion or live account readiness.
