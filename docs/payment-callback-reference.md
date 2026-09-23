# Payment callback composition — local reference

Scope: library infrastructure using the pinned Stripe and Mercado Pago SDK
adapters. The callback processor and queue selector are AUTHORED composition
glue, with no claim of provider-authored business logic or live production.

The HTTP webhook verifies the provider signature, account/mode constraints and
body limits, then uses the existing integration inbox transaction to persist one
normalized notice and one `payment-provider-events` job. A new notification
invalidates a known observation before any subsequent provider GET. Raw payloads,
credentials and payer/card data are not persisted by this bridge.

`PaymentCallbackJobs` scopes claims by tenant, connection, provider, schema and
event type. `PaymentCallbackProcessor` requires the same worker identifier and
checks the durable job generation, lease, exact payload and matching inbox. A
checkout-session notice is only an authenticated resource hint. The processor
retrieves the session through the official SDK and matches its order, attempt,
amount, currency, account, mode, expiry and session identity to the persisted
checkout request and its outbound request hash. It then invokes the existing
payment worker's reconciliation, which starts a new observation generation and
retrieves the payment resource through the SDK. A notification arriving during
GET prevents an older response from clearing its hold.

An initial payment-intent notification may precede checkout completion. If its
resource is not already bound, a first SDK GET supplies lookup metadata, which
must match the durable request. A second fresh GET is performed after the
observation generation starts. Existing, conflicting provider bindings and
ambiguous matches fail closed.

A POST response may be lost after the provider creates the checkout. In that
case, the outbound fence remains `unknown`; processing never retries the POST.
A subsequently signed session notice permits SDK GET recovery of the original
session, subject to the same binding checks. The existing fence is reconciled
to accepted using only the observed session identity. The subsequent payment
GET is still mandatory; session completion alone never captures payment.
Stripe completed-session URLs may be null; recovery uses the immutable session
identity and preserves a previously known URL without inventing a new one.

For Mercado Pago payments, missing inherited preference metadata is supported:
the authenticated payment GET's `external_reference` selects exactly one durable
order request with the configured organization, account, currency, amount, mode,
connection and matching dispatch fence. Zero or ambiguous candidates fail. The
callback body's unsigned financial fields are ignored. The optional
`PaymentCallbackRecoveryProcessor` first completes that financial reconciliation,
then recovers a missing preference identity through the official SDK. It performs
one Search request with external_reference/site_id, limit=2 and offset=0. Exactly
one total result and one element are required, followed by a complete Get of that
ID. Order/attempt metadata, account, mode, currency, amount, site and expiry must
match the original persisted checkout request before saving identity and
reconciling the outbound fence. No recovery path calls CreateCheckout.

The official [preference Search contract](https://www.mercadopago.com.ar/developers/es/reference/online-payments/checkout-pro-preferences/search-preferences/get)
limits results to the last 90 days. Missing, ambiguous or mismatched results
retain unknown identity and use the existing bounded job retry/terminal handling;
they never prove absence of an earlier provider effect or authorize a new POST.
A crash between saving the recovered ID and its fence receipt retries Get of the
known ID without another Search or POST. The financial inbox receipt may already
be processed while the separate identity recovery job still requires a retry.
Payment observation and checkout identity remain distinct evidence.


The job inbox is marked processed only after reconciliation completes and the
lease/generation remain owned. A crash before the later queue acknowledgement
can replay GETs, but existing observation and outbox owners avoid another money
transition. Partial refunds, mismatches and stale responses retain the release
hold. The initial-handover reference contract still requires current payment
evidence and customer acceptance. Its release check is a local read-only
eligibility result; it does not post stock, transfer money, ship goods or claim
production acceptance.

Evidence is supplied by the dedicated PostgreSQL integration suite
`TestPaymentConnectedReference`, its first-intent/claim-fencing test and its
forged-job/metadata test. The connected suite seeds only synthetic business
inputs; quote acceptance, stock allocation, payment request, capture observation,
handover preparation and customer acceptance use their actual owners. Provider
HTTP is restricted to a loopback fixture through an allowlisted transport. The
full connected tests exercise Stripe (including null-URL lost-response recovery)
and Mercado Pago without inherited payment metadata, along with its wrong-order
negative. Three additional preference-recovery tests cover lost POST response,
a crash after saving identity, and seven absent/ambiguous/mismatched identity
cases. Each case issues exactly one provider POST and one money projection;
only SDK Get/Search responses can recover the original preference identity.

The exact fixture contract remains unchanged. Materialized profiles select the
same supported algorithm through `LoadHandoverProfileFile`, with a schema and
algorithm revision, exact document hash, tenant/organization/mode binding and
explicit activation. A policy with unsupported options or missing evidence is
not admitted. Configuration can select a future live provider mode without Go
edits; it does not assert live production readiness. See `docs/handover-profile.md`.
