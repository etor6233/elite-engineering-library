# Initial handover reference

This is AUTHORED composition/persistence glue over FranchiseJourney, Commerce,
the serialized reservation owner, and the official-payment observation bridge.
It does not calculate price/tax/credit, post inventory, turn a payment request
into captured money, implement fiscal policy or declare goods delivered.

The initial handover is created by `HandoverPreparationService.Prepare`, using
`postgres.FranchiseJourney`. Callers provide an authorized tenant/operator,
organization, order/line, payment attempt, observed evidence hash and idempotency
key. Customer, stock, reservation, price and currency come from durable owners.
Neither a client-supplied customer/serial/amount nor `payment_attempt=captured`
without a matching provider observation is sufficient.

The retained fixture contract is explicit: reference-single-unit-observed-payment-v1,
scope LOCAL_FIXTURES, the SHA-256 of ReferenceHandoverContractDocument, and an
explicit finite maximum observation age. This fixture selects non-live mode.
Materialized profiles use the same owner through `LoadHandoverProfileFile`, an
exact-byte activation hash, supported algorithm/revision/options and explicit
tenant/organization/mode bindings. See `docs/handover-profile.md`. Both sandbox
and future live configuration require no Go edits; neither inherits production
acceptance from the fixture tests. Missing, altered or unsupported profiles fail.

Prepare locks order, payment attempt and observation in the reconciler's order,
then fixes the customer/line/stock/reservation relationship. It requires one
quantity-one line, reserved matching serialized stock, active customer and
reservation, full matching received amount/currency, no refund/hold, the exact
observation hash and current time after lock acquisition. Unknown, old, future,
partial, disconnected or mismatched facts fail closed. The reference deliberately
does not cover split payments, multiple lines or a new commercial policy.

One transaction creates delivery_handover, its immutable preparation projection,
the completed idempotency receipt and delivery-handover.prepared outbox event.
Sixteen competing identical commands produce one creation and fifteen replays.
A changed command or a second preparation key for that initial order line is
rejected. Existing exception/reschedule owners remain responsible for successors.

`Result` recovers a committed command without retrying its mutation. It preserves
historical preparation facts and hydrates the current handover acceptance and
checklist. Its receipt is not a new release permission: a later callback can
hold payment while the receipt remains readable.

`EvaluateRelease` rechecks the current payment fence and linked owners after
customer checklist/acceptance. It returns a momentary decision in the selected
LOCAL_FIXTURES or MATERIALIZED_PROFILE scope.
It never changes order to delivered or posts a shipment. Any future effect owner
must recheck the guard inside the effect transaction; it cannot consume this
read-only response as durable authorization. A callback hold after acceptance
rejects the next evaluation.

## HTTP composition

Mount `httpapi.InitialHandoverModule{Service: preparationService}` alongside the
existing FranchiseJourney module after constructing the explicit contract.
An unconfigured module mounts no routes. These commands require the existing
handover:manage permission and organization scope from the verified principal:

- POST /v1/franchise/orders/{id}/handover with organization_id, order_line_id,
  payment_attempt_id, observation_sha256 and Idempotency-Key header.
- GET /v1/franchise/orders/{id}/handover-result?organization_id=... with the
  original Idempotency-Key header. Recovery never resends POST automatically.
- GET /v1/franchise/handovers/{id}/release-check?organization_id=...&observation_sha256=...
  checks the current reference guard without moving stock or money.

The order is taken from the route, and tenant/operator from the verified token.
Unknown customer, order or amount body overrides are rejected. Readiness/conflict
errors do not expose database details. Existing checklist and customer acceptance
routes continue unchanged. No T2804 frontend is added by this increment.

## Scope of source reuse

The underlying serial reservation and customer-shipment owners preserve their
existing Microsoft BCApps-derived admission. Bulk CustomerShipment posts costs,
quantities and inventory; an initial serialized handover is a different projection.
This glue does not claim Microsoft wrote it or that shipment semantics were copied.
No new public upstream or Business Central source group is acquired here.
Owner-domain provenance remains governed by its own admission records; this
increment does not reclassify all Commerce or FranchiseJourney as ADAPTED.

## Focused evidence

The staging runner applies all 56 migrations in an isolated loopback PostgreSQL
instance, checks migration0056 down/up, then executes quote creation/acceptance,
stock allocation, a deliberately synthetic payment observation and handover
creation/checklist/acceptance. No fixture inserts a prepared handover. It exercises
17 invalid scope/observation cases, outbox failure rollback, observation expiry
while waiting for a lock, and recovery with exact metadata after server restart.
It runs bounded domain checks and Go vet. PostgreSQL is stopped afterwards.

The original synthetic-observation suite is not an SDK/webhook end-to-end claim.
The added connected suite separately proves Stripe and Mercado Pago SDK checkout,
signed callback, real queue processor and GET reconciliation through this owner,
without fixture insertion of captured payments or prepared handovers. The profile
suite materializes the policy with Python, loads its exact file/hash in Go and
proves the same PostgreSQL flow under the admitted materialized profile. These
are local infrastructure tests; host wiring and remaining gates have their own
evidence and no production acceptance is asserted here.
