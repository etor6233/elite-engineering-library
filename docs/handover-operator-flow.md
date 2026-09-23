# Operator handover flow

The existing `/franchise` page now includes an order-bound handover panel for `handover:manage`. It reads an authoritative scoped context, prepares the initial handover, embeds the existing checklist completion owner with the server's handover ID/version, and records or recovers the commercial receipt after the customer's acceptance. The existing customer route `/customer/handovers` and `CustomerHandoverActions` continue to own acceptance, observed serial and discrepancy reporting. No operator can manufacture customer acceptance from this new panel.

## Host binding

Mount `httpapi.HandoverContextModule{Service: service}` beside `InitialHandoverModule` and `CommercialReleaseModule`, with the same admitted `*franchisejourney.HandoverPreparationService`. If using the existing helper, its `InitialHandoverModule.Service` concrete value implements `HandoverContextService`. No additional runtime credential, policy switch, SDK, database or dependency is introduced. The active profile remains hash-bound to tenant, organization, provider/account/connection and mode.

The single new Go route is `GET /v1/franchise/orders/{id}/handover-context?organization_id=...`. Verified claims and `handover:manage` membership authorize the read. SQL derives the allocated line, captured attempt and observed evidence, then invokes `lockInitialHandoverScope`. A hold, foreign scope, profile mismatch or incompatible order is not presented as ready. The result includes the current initial handover/version, if one exists. This read never replaces the write's own locked revalidation.

## BFF and browser behavior

`/api/enterprise/handovers` accepts only the supported actions `prepare` and `release`, an order/organization selection and a persisted request key. It obtains payment bindings server-to-server; browser-entered hashes, payment attempts, customer IDs, amounts and stock IDs are rejected. The context returned to the browser strips the internal observation digest, attempt ID and any unknown backend fields.

GET operations are `context`, `prepare-result`, `release-result` and `release-current`. Recovery uses the existing protected backend client with an optional validated `Idempotency-Key` header. The customer/operator access token stays in the server. Responses disable caching. Authentication, organization, duplicate/unknown query parameters, path identifiers and response bindings are checked before returning data.

The browser stores only an action, order ID, request UUID and (for release) handover ID in session storage scoped by the server-derived tenant/actor/organization digest. A storage failure blocks the mutation. An uncertain outcome retains its key and exposes GET recovery, without automatic POST replay or a new key. Keeping the handover ID allows historical recovery even when a later payment callback makes the current context unavailable. A receipt response contains no current authorization flag; current validity is an explicit separate query and is labeled with its evaluation time and expiry.

The existing checklist publisher and recovery controls remain available. The same checklist component supports optional read-only handover ID/version prefill for this per-order flow; its original independent workflow is preserved. No new checklist algorithm, pricing, shipment, tax or financial posting logic is implemented.

## Local evidence and provenance

All deltas are `AUTHORED` transport, persistence/read, presentation and test glue over the admitted owners. The exact existing Go, Next, React, Zod, TypeScript and Vitest dependency locks remain unchanged. No third-party code was copied, and no new dependency, license or notice was added.

Focal checks cover HTTP authorization, real PostgreSQL scope/profile/hold checks, BFF command derivation and recovery, cross-scope and forged fields, retained keys/storage failure, historical versus current receipts, existing checklist/customer acceptance BFF contracts, server rendering, protected-client recovery headers, type generation and TypeScript compilation. Provider reconciliation and the durable commercial write retain their separate already completed SDK/PostgreSQL evidence. No live account or production deployment is claimed.
