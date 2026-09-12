# Go Payment Checkout Runtime

## 1. Metadata

```yaml
pack_id: "GO-PAYMENT-CHECKOUT-RUNTIME"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Durable hosted checkout, official SDK callback/GET reconciliation, scoped customer read and bounded host composition; no direct card handling or invented capture."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["GO-COMMERCE-PRICING-PAYMENT-API0.6.4", "GO-ELECTROMOBILITY-APPLICATION1.10.0", "GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS0.3.2", "TS-GO-API-WEB-BRIDGE0.5.17 where selected"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/stripe/stripe-go/tree/a2df585a800a97fe8ec4ebf551b4449bdb3d90a1", "https://github.com/mercadopago/sdk-go/tree/f910ee53fbb6819e435eaf3d0f800cb1fe74ae09", "https://docs.stripe.com/api/checkout/sessions/object"]
verified_at: "2026-09-11"
```

## 2. Applicability

Use in the current composed franchise reference with the exact declared owner revisions. Missing live credentials leave the provider lane disabled; no credentials are embedded or requested. The local fixture receipts prove the declared contract only. Select the supported materialized profile explicitly; incompatible business options are rejected, not silently invented.

## 3. Architecture contract

Existing Commerce, provider inbox, durable jobs, outbound delivery fence and PostgreSQL remain the state owners. Provider callbacks are signed resource hints; account/mode/order/attempt/currency/amount facts come from the fixed SDK GET. Send acceptance, request binding and pending state share one transaction. Unknown effects are never blindly retried. Handover uses a versioned exact-byte profile and current observations under locks; tenant/org/provider/account/connection/mode are bound. No new financial ledger, statutory decision or implicit physical shipment. Profile and customer routes use existing authorization and server-side identity. Network calls are bounded and outside database locks; queue batches are finite. Roll back the caller pack set before removing its new adapters; do not delete audit rows or provider effects.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/payment.go
CREATE cmd/electromobility-api/payment_test.go
CREATE db/migrations/0055_payment_provider_observation.down.sql
CREATE db/migrations/0055_payment_provider_observation.up.sql
CREATE docs/payment-callback-reference.md
CREATE docs/payment-hosted-checkout.md
CREATE internal/paymentbridge/checkout.go
CREATE internal/paymentbridge/customer_checkout.go
CREATE internal/paymentbridge/customer_checkout_fuzz_test.go
CREATE internal/paymentbridge/sdk_checkout_recovery.go
CREATE internal/paymentbridge/sdk_driver.go
CREATE internal/paymentbridge/sdk_driver_test.go
CREATE internal/paymentbridge/webhook.go
CREATE internal/paymentbridge/webhook_test.go
CREATE internal/platform/httpapi/payment_checkout.go
CREATE internal/platform/httpapi/payment_checkout_test.go
CREATE internal/platform/httpapi/payment_observation_authority_test.go
CREATE internal/platform/postgres/payment_callback_checkout_recovery.go
CREATE internal/platform/postgres/payment_callback_processor.go
CREATE internal/platform/postgres/payment_checkout.go
CREATE internal/platform/postgres/payment_checkout_integration_test.go
CREATE internal/platform/postgres/payment_checkout_read.go
CREATE internal/platform/postgres/payment_checkout_read_integration_test.go
CREATE internal/platform/postgres/payment_checkout_recovery_integration_test.go
CREATE internal/platform/postgres/payment_connected_integration_test.go
CREATE internal/platform/postgres/payment_dispatch_fence.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/payment.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5e58c03b1def8995e9b8ccbb61c8c1b60d5527099b47831a61deb05fd3a7b35c"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED host composition. Secrets are future environment inputs and never
// appear in configuration errors, URLs, logs, or customer responses.
import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/workers"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errPaymentConfiguration = errors.New("payment checkout configuration is incomplete or invalid")

type paymentConfiguration struct {
	Driver                         paymentbridge.SDKDriverConfig
	Token, WebhookSecret, WorkerID string
	HMACKey                        []byte
}

func loadPaymentConfiguration(lookup func(string) string) (*paymentConfiguration, error) {
	if lookup == nil {
		return nil, errPaymentConfiguration
	}
	switch lookup("PAYMENT_CHECKOUT_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errPaymentConfiguration
	}
	scope := paymentbridge.Scope{TenantID: lookup("PAYMENT_TENANT_ID"), OrganizationID: lookup("PAYMENT_ORGANIZATION_ID"), ConnectionID: lookup("PAYMENT_CONNECTION_ID"), ProviderCode: lookup("PAYMENT_REQUEST_PROVIDER"), AccountRef: lookup("PAYMENT_ACCOUNT_REF"), Currency: lookup("PAYMENT_CURRENCY")}
	exponent, err := strconv.Atoi(lookup("PAYMENT_MINOR_UNIT_EXPONENT"))
	if err != nil {
		return nil, errPaymentConfiguration
	}
	scope.MinorUnitExponent = exponent
	switch lookup("PAYMENT_LIVE_MODE") {
	case "true":
		scope.LiveMode = true
	case "false":
	default:
		return nil, errPaymentConfiguration
	}
	c := &paymentConfiguration{Driver: paymentbridge.SDKDriverConfig{Scope: scope, SuccessURL: lookup("PAYMENT_SUCCESS_URL"), CancelURL: lookup("PAYMENT_CANCEL_URL"), NotificationURL: lookup("PAYMENT_NOTIFICATION_URL"), DisplayName: lookup("PAYMENT_DISPLAY_NAME")}, Token: lookup("PAYMENT_PROVIDER_SECRET"), WebhookSecret: lookup("PAYMENT_WEBHOOK_SECRET"), WorkerID: lookup("PAYMENT_WORKER_ID")}
	c.HMACKey, err = base64.StdEncoding.Strict().DecodeString(lookup("PAYMENT_OUTBOUND_HMAC_KEY_BASE64"))
	if err != nil || len(c.HMACKey) < 32 || len(c.HMACKey) > 64 || len(c.WebhookSecret) < 16 || len(c.WebhookSecret) > 16384 || strings.TrimSpace(c.WebhookSecret) != c.WebhookSecret || strings.ContainsAny(c.WebhookSecret, "\r\n") || c.WorkerID == "" || len(c.WorkerID) > 128 || strings.TrimSpace(c.WorkerID) != c.WorkerID || strings.ContainsAny(c.WorkerID, "\r\n") {
		return nil, errPaymentConfiguration
	}
	// Construction checks scope, exact key mode and configured HTTPS redirects;
	// identity is subsequently probed through the official SDK before claims.
	if _, err = paymentbridge.NewSDKDriver(c.Driver, c.Token, nil); err != nil {
		return nil, errPaymentConfiguration
	}
	return c, nil
}

type scopedPaymentSecret struct{ provider, connection, secret string }

func (s scopedPaymentSecret) PaymentWebhookSecret(_ context.Context, provider, connection string) (string, error) {
	if provider != s.provider || connection != s.connection {
		return "", errPaymentConfiguration
	}
	return s.secret, nil
}

type paymentRuntime struct {
	module    httpapi.PaymentCheckoutModule
	worker    paymentbridge.Worker
	callbacks workers.JobProcessor
	driver    *paymentbridge.SDKDriver
}

func preparePaymentRuntime(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string, client *http.Client) (*paymentRuntime, error) {
	c, err := loadPaymentConfiguration(lookup)
	if err != nil || c == nil {
		return nil, err
	}
	if ctx == nil || pool == nil {
		return nil, errPaymentConfiguration
	}
	driver, err := paymentbridge.NewSDKDriver(c.Driver, c.Token, client)
	if err != nil {
		return nil, errPaymentConfiguration
	}
	probe, cancel := context.WithTimeout(ctx, 15*time.Second)
	err = driver.ValidateCredential(probe)
	cancel()
	if err != nil {
		return nil, errPaymentConfiguration
	}
	scope := c.Driver.Scope
	store := postgres.NewPaymentCheckoutStore(pool)
	baseFence, err := postgres.NewOutboundDeliveryStore(pool, c.HMACKey, 2*time.Minute)
	if err != nil {
		return nil, errPaymentConfiguration
	}
	worker := paymentbridge.Worker{Scope: scope, Store: store, Fence: &postgres.PaymentDispatchFence{OutboundDeliveryStore: baseFence, Scope: scope}, Driver: driver}
	callback, err := postgres.NewPaymentCallbackProcessor(pool, worker, c.WorkerID)
	if err != nil {
		return nil, errPaymentConfiguration
	}
	reader, err := postgres.NewCustomerCheckoutReader(pool, scope)
	if err != nil {
		return nil, errPaymentConfiguration
	}
	webhook, err := paymentbridge.NewWebhook(paymentbridge.WebhookConfig{TenantID: scope.TenantID, ConnectionID: scope.ConnectionID, ProviderCode: scope.ProviderCode, LiveMode: scope.LiveMode, Tolerance: 5 * time.Minute, MaxConcurrent: 8, Secrets: scopedPaymentSecret{scope.ProviderCode, scope.ConnectionID, c.WebhookSecret}, Inbox: store})
	if err != nil {
		return nil, errPaymentConfiguration
	}
	return &paymentRuntime{module: httpapi.PaymentCheckoutModule{Reader: reader, Webhook: webhook}, worker: worker, driver: driver, callbacks: workers.JobProcessor{Store: &postgres.PaymentCallbackJobs{Jobs: postgres.NewJobs(pool), Scope: scope}, Handler: postgres.PaymentCallbackRecoveryProcessor{Base: callback}, Queue: "payment-provider-events", WorkerID: c.WorkerID, Lease: 2 * time.Minute, RetryDelay: 30 * time.Second, BatchSize: 1}}, nil
}

// One bounded batch per iteration keeps callbacks responsive. Every provider
// mutation still passes the persisted request binding and outbound fence.
func (r *paymentRuntime) run(ctx context.Context) {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
		iteration, cancel := context.WithTimeout(ctx, 50*time.Second)
		_, sendErr := r.worker.ProcessOnce(iteration, 1)
		_, callbackErr := r.callbacks.ProcessOnce(iteration)
		cancel()
		if ctx.Err() != nil {
			return
		}
		delay := 500 * time.Millisecond
		if sendErr != nil || callbackErr != nil {
			slog.Warn("payment processing deferred; durable state retained")
			delay = 2 * time.Second
		}
		timer.Reset(delay)
	}
}
````

### FILE: `cmd/electromobility-api/payment_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a739eb14741686ebbec43e3b750fc8eb4a8229463ae321c7d30e3a4e32fda974"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"
)

func completePaymentEnvironment() map[string]string {
	return map[string]string{
		"PAYMENT_CHECKOUT_ENABLED": "true", "PAYMENT_TENANT_ID": "tenant", "PAYMENT_ORGANIZATION_ID": "store", "PAYMENT_CONNECTION_ID": "connection", "PAYMENT_REQUEST_PROVIDER": "stripe", "PAYMENT_ACCOUNT_REF": "acct_fixture", "PAYMENT_CURRENCY": "ARS", "PAYMENT_MINOR_UNIT_EXPONENT": "2", "PAYMENT_LIVE_MODE": "false", "PAYMENT_SUCCESS_URL": "https://portal.example.test/customer", "PAYMENT_CANCEL_URL": "https://portal.example.test/customer", "PAYMENT_DISPLAY_NAME": "Order", "PAYMENT_PROVIDER_SECRET": "sk_test_fixture", "PAYMENT_WEBHOOK_SECRET": "whsec_fixture_00000000000000", "PAYMENT_WORKER_ID": "fixture-worker", "PAYMENT_OUTBOUND_HMAC_KEY_BASE64": base64.StdEncoding.EncodeToString([]byte(strings.Repeat("f", 32))),
	}
}
func TestPaymentHostDisabledDoesNotRequireSecrets(t *testing.T) {
	for _, enabled := range []string{"", "false"} {
		lookup := func(key string) string {
			if key == "PAYMENT_CHECKOUT_ENABLED" {
				return enabled
			}
			return ""
		}
		value, err := preparePaymentRuntime(context.Background(), nil, lookup, nil)
		if err != nil || value != nil {
			t.Fatal(value, err)
		}
	}
}
func TestPaymentHostRejectsIncompleteConfiguration(t *testing.T) {
	base := completePaymentEnvironment()
	if _, err := loadPaymentConfiguration(func(k string) string { return base[k] }); err != nil {
		t.Fatal(err)
	}
	for key := range base {
		if key == "PAYMENT_CHECKOUT_ENABLED" {
			continue
		}
		t.Run(key, func(t *testing.T) {
			env := completePaymentEnvironment()
			delete(env, key)
			if _, err := loadPaymentConfiguration(func(k string) string { return env[k] }); err == nil {
				t.Fatal("missing input enabled checkout")
			}
		})
	}
	for _, tc := range []struct{ key, value string }{{"PAYMENT_CHECKOUT_ENABLED", "yes"}, {"PAYMENT_LIVE_MODE", ""}, {"PAYMENT_PROVIDER_SECRET", "sk_live_fixture"}, {"PAYMENT_MINOR_UNIT_EXPONENT", "-1"}, {"PAYMENT_SUCCESS_URL", "http://portal.example.test/"}, {"PAYMENT_ACCOUNT_REF", ""}, {"PAYMENT_OUTBOUND_HMAC_KEY_BASE64", "secret-not-base64"}} {
		env := completePaymentEnvironment()
		env[tc.key] = tc.value
		if _, err := loadPaymentConfiguration(func(k string) string { return env[k] }); err == nil {
			t.Fatal("invalid config accepted", tc.key)
		}
	}
}
````

### FILE: `db/migrations/0055_payment_provider_observation.down.sql`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8184300a02e5bed5d1c804bc2b93b834eb2c34d7ee0307519c36dc6ae30e2735"
variables: []
secrets_allowed: false
```

````sql
begin;
drop table payment.provider_observation_event;
drop table payment.provider_observation;
drop table payment.provider_checkout;
commit;
````

### FILE: `db/migrations/0055_payment_provider_observation.up.sql`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8cdf10bc792d38d8750f81344d2a42070d75a091744d65f55756018cd1377b8c"
variables: []
secrets_allowed: false
```

````sql
begin;

-- AUTHORED persistence glue for the existing payment owner and official SDKs.
-- It does not create a second balance, capture or money-movement ledger.
create table payment.provider_checkout (
  tenant_id uuid not null,
  payment_attempt_id text not null,
  provider_code text not null check (provider_code in ('stripe','mercadopago')),
  connection_id text not null,
  account_ref text not null check (length(account_ref) between 1 and 200),
  request_sha256_hex text not null check (request_sha256_hex ~ '^[0-9a-f]{64}$'),
  session_id text,
  payment_reference text,
  checkout_url text,
  checkout_state text,
  live_mode boolean not null,
  expires_at timestamptz not null,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,payment_attempt_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
  foreign key (tenant_id,connection_id) references integration.provider_connection(tenant_id,connection_id),
  unique (tenant_id,provider_code,session_id),
  unique (tenant_id,provider_code,payment_reference),
  check (session_id is null or length(session_id) between 1 and 200),
  check (payment_reference is null or length(payment_reference) between 1 and 200),
  check (checkout_url is null or length(checkout_url) between 1 and 4096),
  check ((session_id is null and checkout_url is null and checkout_state is null)
      or (session_id is not null and checkout_state is not null and
          (checkout_url is not null or (provider_code='stripe' and checkout_state in ('complete','expired')))))
);

create table payment.provider_observation (
  tenant_id uuid not null,
  payment_attempt_id text not null,
  order_id text not null,
  organization_id text not null,
  provider_code text not null check (provider_code in ('stripe','mercadopago')),
  provider_reference text not null check (length(provider_reference) between 1 and 200),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  amount_minor_units bigint not null check (amount_minor_units > 0),
  received_minor_units bigint not null default 0 check (received_minor_units >= 0),
  refunded_minor_units bigint not null default 0 check (refunded_minor_units >= 0),
  provider_status text not null default 'unobserved',
  live_mode boolean not null,
  account_ref text not null check (length(account_ref) between 1 and 200),
  evidence_sha256_hex text check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  observed_at timestamptz,
  hold boolean not null default true,
  hold_reason text not null default 'OBSERVATION_PENDING',
  generation bigint not null default 1 check (generation > 0),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,payment_attempt_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
  foreign key (tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
  unique (tenant_id,provider_code,provider_reference),
  check (hold or (observed_at is not null and evidence_sha256_hex is not null and received_minor_units=amount_minor_units and refunded_minor_units=0))
);

create table payment.provider_observation_event (
  tenant_id uuid not null,
  payment_attempt_id text not null,
  generation bigint not null,
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  result_code text not null,
  observed_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,payment_attempt_id,generation),
  foreign key (tenant_id,payment_attempt_id) references payment.provider_observation(tenant_id,payment_attempt_id)
);
create trigger payment_observation_event_immutable before update or delete on payment.provider_observation_event
for each row execute function communication.reject_outbound_delivery_event_mutation();
create index payment_observation_hold_idx on payment.provider_observation(tenant_id,provider_code,updated_at) where hold;

commit;
````

### FILE: `docs/payment-callback-reference.md`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5e18ffa19bc05e45ef427eedb86954a479784dba5689c9c8796287e89b10489f"
variables: []
secrets_allowed: false
```

````markdown
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
````

### FILE: `docs/payment-hosted-checkout.md`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "cc6df9d0f86d76f4d4ee3eb308a75a35d321c6966910fc40905b9e6039de9e11"
variables: []
secrets_allowed: false
```

````markdown
# Hosted checkout runtime

This infrastructure composes the existing Commerce request/outbox, pinned official payment SDKs, PostgreSQL outbound fence, provider inbox/jobs, and payment observation owner. New integration files are AUTHORED glue; no local checkout collects payer/card details or decides that a browser return means payment succeeded.

## Activation

The `cmd/electromobility-api` host reads `PAYMENT_CHECKOUT_ENABLED`; absent or `false` enables no payment worker, webhook route or new payment request. `true` requires every applicable input in `.env.example`. Secrets stay empty in this library. The host fails startup with a bounded configuration error if an enabled configuration is incomplete or the official account probe does not match. An inactive payment integration does not prevent other application capabilities from starting.

Each host selects one tenant/organization/provider/connection/account/currency/mode scope. Commerce request endpoints reject scopes outside that configured selection. The provider connection must already exist and be active in the materialized application's integration registry, with the selected tenant, organization and provider. Worker IDs must distinguish concurrent hosts. PostgreSQL claims and observation generations remain the authority for retries and completion.

| Input | Contract when enabled |
|---|---|
| `PAYMENT_TENANT_ID`, `PAYMENT_ORGANIZATION_ID`, `PAYMENT_CONNECTION_ID` | Exact server-side registry scope; not supplied by the browser |
| `PAYMENT_REQUEST_PROVIDER` | `stripe` or `mercadopago` |
| `PAYMENT_ACCOUNT_REF` | Stripe account ID or MercadoPago collector ID, checked using the credential's official account GET |
| `PAYMENT_CURRENCY`, `PAYMENT_MINOR_UNIT_EXPONENT` | Exact selected contract; MercadoPago reference lane requires ARS and exponent2 |
| `PAYMENT_LIVE_MODE` | Explicit `false` or `true`; Stripe key mode must match before any HTTP; payment GET mode must match before settlement projection |
| `PAYMENT_SUCCESS_URL`, `PAYMENT_CANCEL_URL` | Configured HTTPS portal return destinations; neither is a payment receipt |
| `PAYMENT_NOTIFICATION_URL` | MercadoPago HTTPS API callback destination `/v1/payment-provider/webhook`; unused for Stripe |
| `PAYMENT_DISPLAY_NAME`, `PAYMENT_WORKER_ID` | Bounded nonempty provider item label and unique worker identity |
| `PAYMENT_PROVIDER_SECRET`, `PAYMENT_WEBHOOK_SECRET` | Future secret-manager/environment inputs; never returned or logged |
| `PAYMENT_OUTBOUND_HMAC_KEY_BASE64` | Future base64-encoded32–64-byte key for the existing recipient digest fence; retain across restarts |

Checkout expiry is fixed once from the first durable `payment.requested` timestamp plus23hours. Retries reuse that exact expiry and the payment attempt idempotency key. The request owner rejects a dispatch with less than31minutes remaining. An ambiguous provider POST is held in the existing `unknown` fence; the worker does not blindly repeat the POST.

Stripe callbacks are registered with the provider for `/v1/payment-provider/webhook`; the selected reference uses the account's own endpoint, not Connect account routing. MercadoPago sets its configured notification URL on the preference. The endpoint passes raw bytes to the existing official signature normalizer. Acknowledgement occurs only after `PaymentCheckoutStore.AcceptWebhook` commits the receipt, scoped `payment-provider-events` job and applicable observation hold together. `PaymentCallbackProcessor` and the existing `workers.JobProcessor` consume that exact scoped queue and reconcile via official GETs. Shutdown cancels and waits for the bounded worker before closing PostgreSQL.

## Customer flow

An authorized operator requests payment through the existing Commerce command. The worker creates provider-hosted checkout. In the customer portal, “Continuar con el pago” calls the same-origin BFF, which forwards the session access token to `GET /v1/customer/orders/{id}/checkout?organization_id=...`. The API requires `customer:self`, an allowed organization and an active customer profile. Tenant and customer subject come exclusively from the verified token; PostgreSQL must find exactly one eligible checkout for that customer/order and configured account scope. Expired, foreign, ambiguous, completed and invalid provider URLs are not exposed.

The BFF and browser validate the provider redirect host and expiry. Responses use no-store and no-referrer policies. No token, provider payment reference, secret, price override or customer identity is accepted from a checkout URL request. Refund/dispute and commercial release remain governed by provider GET observations and the existing domain owners.

## Local evidence

Focused host/API configuration and authorization tests,13 BFF tests, Next route type generation and TypeScript checks pass. The dedicated PostgreSQL customer-read fixture rejects foreign tenant/organization/customer, restricted customer and tampered provider URL. Official SDK serialization, signed callbacks, refund/dispute handling and durable handover progression have separate connected integration receipts. The logger correction has an independent subprocess test: malformed Stripe responses fail without raw response samples in stdout/stderr. No live credential or provider endpoint was used for these checks.

Missing user credentials condition activation only. These local receipts are not production acceptance, country/fiscal approval or a live payment certification.
````

### FILE: `internal/paymentbridge/checkout.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4da77f0226c3c018d5553d98f35ee464bc0dc2c7ec64be49931d523816ae30ea"
variables: []
secrets_allowed: false
```

````go
package paymentbridge

// AUTHORED orchestration glue. Provider effects belong to the pinned SDK adapter;
// the existing outbound fence prevents an uncertain call from being sent twice.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
)

var ErrCheckoutConflict = errors.New("payment checkout binding conflict")
var ErrObservationStale = errors.New("payment observation superseded")
var ErrPaymentMismatch = errors.New("provider observation differs from payment contract")

type Scope struct {
	TenantID, OrganizationID, ConnectionID, ProviderCode, AccountRef string
	Currency                                                         string
	MinorUnitExponent                                                int
	LiveMode                                                         bool
}

func (s Scope) Validate() error {
	for _, v := range []string{s.TenantID, s.OrganizationID, s.ConnectionID, s.AccountRef} {
		if strings.TrimSpace(v) != v || v == "" || len(v) > 200 || strings.ContainsAny(v, "\r\n") {
			return ErrCheckoutConflict
		}
	}
	if (s.ProviderCode != "stripe" && s.ProviderCode != "mercadopago") || len(s.Currency) != 3 || s.Currency != strings.ToUpper(s.Currency) || s.MinorUnitExponent < 0 || s.MinorUnitExponent > 3 {
		return ErrCheckoutConflict
	}
	return nil
}

type Request struct {
	CheckoutExpiresAt                                                                            time.Time
	TenantID, OrganizationID, PaymentAttemptID, OrderID, CustomerSubject, ProviderCode, Currency string
	AmountMinor                                                                                  int64
}

type Checkout struct {
	ProviderCode, SessionID, URL, PaymentReference, Status string
	OrderID, PaymentAttemptID, Currency, AccountRef        string
	AmountMinor                                            int64
	LiveMode                                               bool
	ExpiresAt                                              time.Time
}
type Driver interface {
	CreateCheckout(context.Context, Request) (Checkout, error)
	RetrieveCheckout(context.Context, string) (Checkout, error)
	RetrievePayment(context.Context, string) (officialpayments.PaymentSnapshot, error)
}

type CheckoutStore interface {
	PendingRequests(context.Context, Scope, int) ([]Request, error)
	BindRequest(context.Context, Scope, Request, string) error
	SaveCheckout(context.Context, Scope, Request, string, Checkout) error
	Checkout(context.Context, Scope, string) (Request, Checkout, error)
	BeginObservation(context.Context, Scope, string, string) (Request, int64, error)
	RecordObservation(context.Context, Scope, Request, int64, officialpayments.PaymentSnapshot) error
}

type Worker struct {
	Scope  Scope
	Store  CheckoutStore
	Fence  outbounddelivery.Store
	Driver Driver
}

func (w Worker) valid() error {
	if w.Scope.Validate() != nil || w.Store == nil || w.Fence == nil || w.Driver == nil {
		return ErrCheckoutConflict
	}
	return nil
}

// message binds immutable database values, not caller-supplied prices or URLs.
func (r Request) message() (channels.Message, string, error) {
	if r.AmountMinor <= 0 || r.PaymentAttemptID == "" || r.OrderID == "" || r.CustomerSubject == "" {
		return channels.Message{}, "", ErrCheckoutConflict
	}
	data, err := json.Marshal(r)
	if err != nil {
		return channels.Message{}, "", err
	}
	m := channels.Message{ChannelCode: "payment_" + r.ProviderCode, TenantID: r.TenantID, ExternalID: r.CustomerSubject, ThreadID: r.OrderID, DeliveryKey: r.PaymentAttemptID, Direction: channels.DirectionOut, Text: string(data)}
	hash, err := outbounddelivery.MessageSHA256(m)
	return m, hash, err
}

func (w Worker) ProcessOnce(ctx context.Context, limit int) (int, error) {
	if w.valid() != nil || limit < 1 || limit > 100 {
		return 0, ErrCheckoutConflict
	}
	requests, err := w.Store.PendingRequests(ctx, w.Scope, limit)
	if err != nil {
		return 0, err
	}
	completed := 0
	for _, r := range requests {
		if ctx.Err() != nil {
			return completed, ctx.Err()
		}
		m, hash, err := r.message()
		if err != nil {
			return completed, err
		}
		if err = w.Store.BindRequest(ctx, w.Scope, r, hash); err != nil {
			return completed, err
		}
		claim, err := w.Fence.Claim(ctx, m, hash)
		if errors.Is(err, outbounddelivery.ErrUnknown) || errors.Is(err, outbounddelivery.ErrTerminal) || errors.Is(err, outbounddelivery.ErrInProgress) {
			continue
		}
		if err != nil {
			return completed, err
		}
		if claim.Replay {
			continue
		}
		checkout, err := w.Driver.CreateCheckout(ctx, r)
		if err != nil {
			mark := w.Fence.MarkUnknown(ctx, m, hash, "CHECKOUT_CALL_UNCERTAIN")
			if mark != nil {
				return completed, mark
			}
			continue
		}
		if checkout.ProviderCode != r.ProviderCode || checkout.SessionID == "" || checkout.URL == "" || checkout.Status == "" {
			mark := w.Fence.MarkUnknown(ctx, m, hash, "CHECKOUT_RESPONSE_INVALID")
			if mark != nil {
				return completed, mark
			}
			continue
		}
		if err = w.Store.SaveCheckout(ctx, w.Scope, r, hash, checkout); err != nil {
			_ = w.Fence.MarkUnknown(ctx, m, hash, "CHECKOUT_PERSISTENCE_UNCERTAIN")
			return completed, err
		}
		// URLs may contain opaque session tokens; only normalized identity is hashed.
		proof, _ := json.Marshal(struct{ Provider, Session, Payment, Status string }{checkout.ProviderCode, checkout.SessionID, checkout.PaymentReference, checkout.Status})
		sum := sha256.Sum256(proof)
		receipt := outbounddelivery.Receipt{ProviderMessageID: checkout.SessionID, EvidenceSHA256: hex.EncodeToString(sum[:]), AcceptedAt: time.Now().UTC()}
		if err = w.Fence.Complete(ctx, m, hash, receipt); err != nil {
			return completed, err
		}
		completed++
	}
	return completed, nil
}

// Reconcile starts a generation before GET. A later callback invalidates that
// generation, so a delayed older response cannot clear the handover hold.
func (w Worker) Reconcile(ctx context.Context, attempt, providerReference string) error {
	if w.valid() != nil || attempt == "" || providerReference == "" {
		return ErrCheckoutConflict
	}
	r, generation, err := w.Store.BeginObservation(ctx, w.Scope, attempt, providerReference)
	if err != nil {
		return err
	}
	observed, err := w.Driver.RetrievePayment(ctx, providerReference)
	if err != nil {
		return err
	}
	return w.Store.RecordObservation(ctx, w.Scope, r, generation, observed)
}
````

### FILE: `internal/paymentbridge/customer_checkout.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1238275631751549d8a44c6b7f390a856644c35c80ac76723cbc05227bdb60b7"
variables: []
secrets_allowed: false
```

````go
package paymentbridge

// AUTHORED read contract. This exposes a provider-hosted redirect, never payer
// data or a browser-authorized payment transition.
import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrCheckoutNotAvailable = errors.New("customer checkout is not available")

type CustomerCheckout struct {
	OrderID      string    `json:"order_id"`
	ProviderCode string    `json:"provider_code"`
	URL          string    `json:"url"`
	ExpiresAt    time.Time `json:"expires_at"`
}
type CustomerCheckoutReader interface {
	CustomerCheckout(context.Context, string, string, string, string) (CustomerCheckout, error)
}

func (c CustomerCheckout) Valid(now time.Time) bool {
	u, err := url.Parse(c.URL)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Fragment != "" || (u.Port() != "" && u.Port() != "443") || len(c.URL) > 8192 || !utf8.ValidString(c.URL) || strings.ContainsAny(c.URL, "\r\n") || !c.ExpiresAt.After(now) || c.OrderID == "" {
		return false
	}
	switch c.ProviderCode {
	case "stripe":
		return strings.EqualFold(u.Hostname(), "checkout.stripe.com")
	case "mercadopago":
		return strings.EqualFold(u.Hostname(), "www.mercadopago.com.ar") || strings.EqualFold(u.Hostname(), "sandbox.mercadopago.com.ar")
	}
	return false
}
````

### FILE: `internal/paymentbridge/customer_checkout_fuzz_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "183585b941bc6e2e59740320fe716bedf1b1359552f27e34f78393f2c17e968c"
variables: []
secrets_allowed: false
```

````go
package paymentbridge

// AUTHORED verification glue for the provider redirect boundary. These fixtures
// contain synthetic provider identifiers and do not contact either provider.
import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func FuzzCustomerCheckoutRedirect(f *testing.F) {
	for _, seed := range []struct{ url, provider string }{
		{"https://checkout.stripe.com/c/pay/cs_test_connected", "stripe"},
		{"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=pref_fixture", "mercadopago"},
		{"https://www.mercadopago.com.ar/checkout/v1/redirect?pref_id=pref_fixture", "mercadopago"},
		{"https://checkout.stripe.com.evil.invalid/c/pay/test", "stripe"},
		{"https://user@checkout.stripe.com/c/pay/test", "stripe"},
		{"http://checkout.stripe.com/c/pay/test", "stripe"},
		{"https://checkout.stripe.com/c/pay/test#other", "stripe"},
		{"https://checkout.stripe.com/c/pay/\xff", "stripe"},
		{"https://checkout.stripe.com/c/pay/test", "mercadopago"},
	} {
		f.Add(seed.url, seed.provider, int64(300))
	}
	f.Add("https://checkout.stripe.com/c/pay/expired", "stripe", int64(0))
	f.Fuzz(func(t *testing.T, rawURL, provider string, offset int64) {
		if len(rawURL) > 9000 || len(provider) > 128 {
			t.Skip()
		}
		now := time.Unix(1800000000, 0).UTC()
		c := CustomerCheckout{OrderID: "order-fixture", ProviderCode: provider, URL: rawURL, ExpiresAt: now.Add(time.Duration(offset%86400) * time.Second)}
		if !c.Valid(now) {
			return
		}
		// The HTTP consumer and JSON transport must agree with admission. A
		// parser discrepancy must not create a different redirect downstream.
		req, err := http.NewRequest(http.MethodGet, c.URL, nil)
		if err != nil {
			t.Fatalf("admitted redirect cannot be consumed: %v", err)
		}
		allowed := map[string]map[string]bool{
			"stripe":      {"checkout.stripe.com": true},
			"mercadopago": {"www.mercadopago.com.ar": true, "sandbox.mercadopago.com.ar": true},
		}
		if req.URL.Scheme != "https" || req.URL.Opaque != "" || req.URL.User != nil || req.URL.Fragment != "" || !allowed[provider][strings.ToLower(req.URL.Hostname())] {
			t.Fatal("admitted redirect escapes its provider origin")
		}
		wire, err := json.Marshal(c)
		if err != nil {
			t.Fatal(err)
		}
		var restored CustomerCheckout
		if json.Unmarshal(wire, &restored) != nil || restored.URL != c.URL || restored.ProviderCode != c.ProviderCode || restored.OrderID != c.OrderID || !restored.ExpiresAt.Equal(c.ExpiresAt) || !restored.Valid(now) {
			t.Fatal("JSON transport changes an admitted redirect")
		}
		expired := c
		expired.ExpiresAt = now
		if expired.Valid(now) {
			t.Fatal("expired redirect remains usable")
		}
		other := c
		if provider == "stripe" {
			other.ProviderCode = "mercadopago"
		} else {
			other.ProviderCode = "stripe"
		}
		if other.Valid(now) {
			t.Fatal("redirect can be transplanted to another provider")
		}
		for _, mutate := range []func(*url.URL){
			func(u *url.URL) { u.Host = "checkout.stripe.com.evil.invalid" },
			func(u *url.URL) { u.User = url.User("confused") },
			func(u *url.URL) { u.Scheme = "http" },
			func(u *url.URL) { u.Fragment = "changed" },
		} {
			u := *req.URL
			mutate(&u)
			changed := c
			changed.URL = u.String()
			if changed.Valid(now) {
				t.Fatal("redirect mutation crosses an admitted boundary")
			}
		}
	})
}
````

### FILE: `internal/paymentbridge/sdk_checkout_recovery.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f2f3566eedeef4481ef41744d64bd1d21d83652a420285d987f24a8078d3c89f"
variables: []
secrets_allowed: false
```

````go
package paymentbridge

// AUTHORED optional recovery wiring. The pinned provider SDK performs both
// preference Search and Get; this path never calls CreateCheckout.
import (
	"context"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
)

type CheckoutRecoveryDriver interface {
	RecoverCheckout(context.Context, Request) (Checkout, error)
}

func (d *SDKDriver) RecoverCheckout(ctx context.Context, r Request) (Checkout, error) {
	if d == nil || ctx == nil || d.mpCheckout == nil || d.config.Scope.ProviderCode != "mercadopago" || r.ProviderCode != d.config.Scope.ProviderCode || r.TenantID != d.config.Scope.TenantID || r.OrganizationID != d.config.Scope.OrganizationID || r.Currency != d.config.Scope.Currency {
		return Checkout{}, ErrCheckoutConflict
	}
	if err := d.ValidateCredential(ctx); err != nil {
		return Checkout{}, err
	}
	result, err := d.mpCheckout.RecoverCheckout(ctx, officialpayments.PreferenceRecoveryRequest{OrderID: r.OrderID, PaymentAttemptID: r.PaymentAttemptID, Currency: r.Currency, CollectorID: d.config.Scope.AccountRef, AmountMinor: r.AmountMinor, MinorUnitExponent: d.config.Scope.MinorUnitExponent, LiveMode: d.config.Scope.LiveMode, ExpiresAt: r.CheckoutExpiresAt})
	if err != nil {
		return Checkout{}, ErrSDKDriverUnavailable
	}
	return d.checkout(result)
}
````

### FILE: `internal/paymentbridge/sdk_driver.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fc26f4bdae3576b5cafd7e90f6d101a8ae4ae7582a696110a7bfe79993e02b17"
variables: []
secrets_allowed: false
```

````go
package paymentbridge

// AUTHORED composition only: provider HTTP serialization, authentication and
// resource GETs are performed by the exact official SDK dependency pins.
import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
	mpconfig "github.com/mercadopago/sdk-go/pkg/config"
	mpuser "github.com/mercadopago/sdk-go/pkg/user"
)

var ErrSDKDriverUnavailable = errors.New("payment SDK operation unavailable")

type SDKDriverConfig struct {
	Scope                                  Scope
	SuccessURL, CancelURL, NotificationURL string
	DisplayName                            string
}
type SDKDriver struct {
	config     SDKDriverConfig
	stripe     *officialpayments.StripeIntentClient
	mpPayment  *officialpayments.MercadoPagoPaymentClient
	mpCheckout *officialpayments.MercadoPagoCheckoutClient
	mpUser     mpuser.Client
}

func configuredHTTPS(value string) bool {
	u, err := url.Parse(value)
	return err == nil && u.Scheme == "https" && u.Hostname() != "" && u.User == nil && len(value) <= 2048 && !strings.ContainsAny(value, "\r\n") && (u.Port() == "" || u.Port() == "443")
}

// No network is used by construction. ValidateCredential must succeed at host
// startup before taking a durable send claim; CreateCheckout rechecks it.
func NewSDKDriver(c SDKDriverConfig, token string, client *http.Client) (*SDKDriver, error) {
	if c.Scope.Validate() != nil || strings.TrimSpace(token) == "" || len(token) > 16384 || !configuredHTTPS(c.SuccessURL) || !configuredHTTPS(c.CancelURL) || strings.TrimSpace(c.DisplayName) == "" || len(c.DisplayName) > 200 {
		return nil, ErrCheckoutConflict
	}
	if c.Scope.ProviderCode == "mercadopago" && (!configuredHTTPS(c.NotificationURL) || c.Scope.Currency != "ARS" || c.Scope.MinorUnitExponent != 2) {
		return nil, ErrCheckoutConflict
	}
	// Stripe's documented test/live secret and restricted key prefixes provide
	// an early configuration guard; the official account GET still proves identity.
	if c.Scope.ProviderCode == "stripe" {
		mode := "test_"
		if c.Scope.LiveMode {
			mode = "live_"
		}
		if (!strings.HasPrefix(token, "sk_"+mode) && !strings.HasPrefix(token, "rk_"+mode)) || len(token) <= len("sk_"+mode) || strings.TrimSpace(token) != token || strings.ContainsAny(token, "\r\n") {
			return nil, ErrPaymentMismatch
		}
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	bounded := *client
	bounded.Timeout = 10 * time.Second
	bounded.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	d := &SDKDriver{config: c}
	var err error
	if c.Scope.ProviderCode == "stripe" {
		d.stripe, err = officialpayments.NewStripeIntentClientWithHTTPClient(token, &bounded)
	} else {
		d.mpPayment, err = officialpayments.NewMercadoPagoPaymentClientWithHTTPClient(token, &bounded)
		if err != nil {
			return nil, ErrCheckoutConflict
		}
		d.mpCheckout, err = officialpayments.NewMercadoPagoCheckoutClientWithHTTPClient(token, !c.Scope.LiveMode, &bounded)
		if err != nil {
			return nil, ErrCheckoutConflict
		}
		cfg, configErr := mpconfig.New(token, mpconfig.WithHTTPClient(&bounded), mpconfig.WithMaxRetries(0))
		if configErr != nil {
			return nil, ErrCheckoutConflict
		}
		d.mpUser = mpuser.NewClient(cfg)
	}
	if err != nil {
		return nil, ErrCheckoutConflict
	}
	return d, nil
}
func (d *SDKDriver) ValidateCredential(ctx context.Context) error {
	if d == nil || ctx == nil {
		return ErrCheckoutConflict
	}
	if d.stripe != nil {
		id, err := d.stripe.ProbeAccount(ctx)
		if err != nil {
			return ErrSDKDriverUnavailable
		}
		if id != d.config.Scope.AccountRef {
			return ErrPaymentMismatch
		}
		return nil
	}
	if d.mpUser == nil {
		return ErrCheckoutConflict
	}
	account, err := d.mpUser.Get(ctx)
	if err != nil {
		return ErrSDKDriverUnavailable
	}
	if account == nil || strconv.Itoa(account.ID) != d.config.Scope.AccountRef || account.CountryID != "AR" || account.SiteID != "MLA" {
		return ErrPaymentMismatch
	}
	return nil
}
func (d *SDKDriver) CreateCheckout(ctx context.Context, r Request) (Checkout, error) {
	if d == nil || ctx == nil || r.TenantID != d.config.Scope.TenantID || r.OrganizationID != d.config.Scope.OrganizationID || r.ProviderCode != d.config.Scope.ProviderCode || r.Currency != d.config.Scope.Currency || r.AmountMinor <= 0 || r.PaymentAttemptID == "" || r.OrderID == "" || r.CheckoutExpiresAt.IsZero() {
		return Checkout{}, ErrPaymentMismatch
	}
	if err := d.ValidateCredential(ctx); err != nil {
		return Checkout{}, err
	}
	request := officialpayments.CheckoutRequest{OrderID: r.OrderID, PaymentAttemptID: r.PaymentAttemptID, AmountMinor: r.AmountMinor, Currency: r.Currency, MinorUnitExponent: d.config.Scope.MinorUnitExponent, DisplayName: d.config.DisplayName, SuccessURL: d.config.SuccessURL, CancelURL: d.config.CancelURL, NotificationURL: d.config.NotificationURL, IdempotencyKey: r.PaymentAttemptID, ExpiresAt: r.CheckoutExpiresAt}
	var result officialpayments.CheckoutResult
	var err error
	if d.stripe != nil {
		result, err = d.stripe.CreateCheckout(ctx, request)
	} else {
		result, err = d.mpCheckout.CreateCheckout(ctx, request)
	}
	if err != nil {
		return Checkout{}, ErrSDKDriverUnavailable
	}
	checkout, err := d.checkout(result)
	if err != nil || checkout.OrderID != r.OrderID || checkout.PaymentAttemptID != r.PaymentAttemptID || checkout.AmountMinor != r.AmountMinor || checkout.ExpiresAt.Unix() != r.CheckoutExpiresAt.Unix() {
		return Checkout{}, ErrPaymentMismatch
	}
	return checkout, nil
}
func (d *SDKDriver) RetrieveCheckout(ctx context.Context, id string) (Checkout, error) {
	if d == nil || ctx == nil || id == "" {
		return Checkout{}, ErrCheckoutConflict
	}
	if err := d.ValidateCredential(ctx); err != nil {
		return Checkout{}, err
	}
	var result officialpayments.CheckoutResult
	var err error
	if d.stripe != nil {
		result, err = d.stripe.RetrieveCheckout(ctx, id)
	} else {
		result, err = d.mpCheckout.RetrieveCheckout(ctx, id, d.config.Scope.MinorUnitExponent)
	}
	if err != nil {
		return Checkout{}, ErrSDKDriverUnavailable
	}
	return d.checkout(result)
}
func (d *SDKDriver) checkout(r officialpayments.CheckoutResult) (Checkout, error) {
	if r.ProviderCode != d.config.Scope.ProviderCode || r.Currency != d.config.Scope.Currency || r.LiveMode != d.config.Scope.LiveMode || r.OrderID == "" || r.PaymentAttemptID == "" || r.AmountMinor <= 0 || r.ExpiresAt.IsZero() {
		return Checkout{}, ErrPaymentMismatch
	}
	if r.ProviderCode == "mercadopago" && r.CollectorID != d.config.Scope.AccountRef {
		return Checkout{}, ErrPaymentMismatch
	}
	return Checkout{ProviderCode: r.ProviderCode, SessionID: r.SessionID, URL: r.URL, PaymentReference: r.PaymentReference, Status: r.Status, OrderID: r.OrderID, PaymentAttemptID: r.PaymentAttemptID, Currency: r.Currency, AccountRef: d.config.Scope.AccountRef, AmountMinor: r.AmountMinor, LiveMode: r.LiveMode, ExpiresAt: r.ExpiresAt}, nil
}
func (d *SDKDriver) RetrievePayment(ctx context.Context, id string) (officialpayments.PaymentSnapshot, error) {
	if d == nil || ctx == nil || id == "" {
		return officialpayments.PaymentSnapshot{}, ErrCheckoutConflict
	}
	if err := d.ValidateCredential(ctx); err != nil {
		return officialpayments.PaymentSnapshot{}, err
	}
	var result officialpayments.PaymentSnapshot
	var err error
	if d.stripe != nil {
		result, err = d.stripe.RetrieveIntent(ctx, id)
	} else {
		result, err = d.mpPayment.RetrievePayment(ctx, id, d.config.Scope.MinorUnitExponent)
	}
	if err != nil {
		return officialpayments.PaymentSnapshot{}, ErrSDKDriverUnavailable
	}
	if result.ProviderReference != id || result.ProviderCode != d.config.Scope.ProviderCode || result.Currency != d.config.Scope.Currency || result.LiveMode != d.config.Scope.LiveMode || (result.ProviderCode == "mercadopago" && result.CollectorID != d.config.Scope.AccountRef) {
		return officialpayments.PaymentSnapshot{}, ErrPaymentMismatch
	}
	return result, nil
}
````

### FILE: `internal/paymentbridge/sdk_driver_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "39aee9de29bcb33db266a00ab5502237867ff53ed4333e0c0cd7f553ab451144"
variables: []
secrets_allowed: false
```

````go
package paymentbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/providerintegration"
	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

type localProviderTransport struct {
	destination  *url.URL
	providerHost string
}

func (t localProviderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Scheme != "https" || r.URL.Host != t.providerHost {
		return nil, fmt.Errorf("fixture forbids network host")
	}
	cloned := r.Clone(r.Context())
	u := *r.URL
	u.Scheme = t.destination.Scheme
	u.Host = t.destination.Host
	cloned.URL = &u
	cloned.Host = ""
	return http.DefaultTransport.RoundTrip(cloned)
}
func fixtureHTTPClient(t *testing.T, host string, handler http.HandlerFunc) *http.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	parsed, _ := url.Parse(server.URL)
	return &http.Client{Transport: localProviderTransport{parsed, host}}
}
func TestSDKDriverHostedSessionCallbackAndRefundDisputeUseOfficialHTTP(t *testing.T) {
	var mu sync.Mutex
	account := "acct_fixture"
	refund := int64(0)
	disputed := false
	posts := 0
	expires := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	client := fixtureHTTPClient(t, "api.stripe.com", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer sk_test_fixture" {
			t.Error("missing fixture credential")
		}
		switch r.URL.Path {
		case "/v1/account":
			if r.Method != "GET" {
				t.Error("account write")
			}
			fmt.Fprintf(w, `{"id":%q}`, account)
		case "/v1/checkout/sessions", "/v1/checkout/sessions/cs_test_fixture":
			if r.Method == "POST" {
				posts++
				r.ParseForm()
				if r.Form.Get("metadata[payment_attempt_id]") != "attempt-1" || r.Form.Get("payment_intent_data[metadata][order_id]") != "order-1" || r.Form.Get("expires_at") != strconv.FormatInt(expires.Unix(), 10) || r.Header.Get("Idempotency-Key") != "attempt-1" {
					t.Errorf("wrong stable checkout request")
				}
			}
			fmt.Fprintf(w, `{"id":"cs_test_fixture","object":"checkout.session","mode":"payment","amount_total":125050,"currency":"ars","client_reference_id":"attempt-1","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"status":"open","payment_status":"unpaid","livemode":false,"url":"https://checkout.stripe.com/c/pay/cs_test_fixture","payment_intent":"pi_fixture","expires_at":%d}`, expires.Unix())
		case "/v1/payment_intents/pi_fixture":
			if r.Method != "GET" || r.URL.Query().Get("expand[0]") != "latest_charge" {
				t.Error("payment observation is not expanded GET")
			}
			fmt.Fprintf(w, `{"id":"pi_fixture","amount":125050,"amount_received":125050,"currency":"ars","status":"succeeded","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"latest_charge":{"id":"ch_fixture","amount":125050,"amount_captured":125050,"amount_refunded":%d,"currency":"ars","captured":true,"disputed":%t}}`, refund, disputed)
		default:
			t.Errorf("unexpected path%s", r.URL)
			w.WriteHeader(500)
		}
	})
	scope := Scope{TenantID: "tenant", OrganizationID: "store", ConnectionID: "connection", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}
	d, err := NewSDKDriver(SDKDriverConfig{Scope: scope, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", DisplayName: "Order"}, "sk_test_fixture", client)
	if err != nil {
		t.Fatal(err)
	}
	req := Request{TenantID: scope.TenantID, OrganizationID: scope.OrganizationID, ProviderCode: scope.ProviderCode, Currency: scope.Currency, PaymentAttemptID: "attempt-1", OrderID: "order-1", CustomerSubject: "customer", AmountMinor: 125050, CheckoutExpiresAt: expires}
	created, err := d.CreateCheckout(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	inbox := &fixtureInbox{records: map[string]providerintegration.Receipt{}}
	h, _ := NewWebhook(WebhookConfig{TenantID: scope.TenantID, ConnectionID: scope.ConnectionID, ProviderCode: "stripe", Tolerance: time.Minute, MaxConcurrent: 1, Secrets: fixtureSecret{}, Inbox: inbox})
	body := []byte(fmt.Sprintf(`{"id":"evt_session","object":"event","api_version":%q,"type":"checkout.session.completed","livemode":false,"data":{"object":{"id":"cs_test_fixture","object":"checkout.session"}}}`, stripe.APIVersion))
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "whsec_fixture", Timestamp: time.Now()})
	callback := httptest.NewRequest("POST", "https://api.example.test/payment-webhook", strings.NewReader(string(body)))
	callback.Header.Set("Content-Type", "application/json")
	callback.Header.Set("Stripe-Signature", signed.Header)
	response := httptest.NewRecorder()
	h.ServeHTTP(response, callback)
	if response.Code != 200 {
		t.Fatal(response.Code)
	}
	var payload InboxPayload
	if json.Unmarshal(inbox.records["evt_session"].Payload, &payload) != nil || payload.Notification.ResourceType != "checkout_session" {
		t.Fatal("wrong durable notice")
	}
	observedCheckout, err := d.RetrieveCheckout(context.Background(), payload.Notification.ProviderReference)
	if err != nil || observedCheckout != created {
		t.Fatalf("checkout %+v %v", observedCheckout, err)
	}
	paid, err := d.RetrievePayment(context.Background(), observedCheckout.PaymentReference)
	if err != nil || paid.ReceivedMinor != req.AmountMinor || paid.RefundedMinor != 0 || paid.Disputed {
		t.Fatalf("paid %+v %v", paid, err)
	}
	mu.Lock()
	refund = 1
	mu.Unlock()
	refunded, err := d.RetrievePayment(context.Background(), observedCheckout.PaymentReference)
	if err != nil || refunded.RefundedMinor != 1 {
		t.Fatalf("refund %+v %v", refunded, err)
	}
	mu.Lock()
	refund = 0
	disputed = true
	mu.Unlock()
	dispute, err := d.RetrievePayment(context.Background(), observedCheckout.PaymentReference)
	if err != nil || !dispute.Disputed {
		t.Fatalf("dispute %+v %v", dispute, err)
	}
	mu.Lock()
	account = "acct_other"
	mu.Unlock()
	if _, err = d.CreateCheckout(context.Background(), req); err == nil {
		t.Fatal("wrong credential account allowed")
	}
	mu.Lock()
	postCount := posts
	mu.Unlock()
	if postCount != 1 {
		t.Fatalf("unexpected provider writes%d", postCount)
	}
}
func TestSDKDriverMercadoPagoProbesAccountBeforeHostedCheckout(t *testing.T) {
	var mu sync.Mutex
	expires := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	posts := 0
	account := 123
	client := fixtureHTTPClient(t, "api.mercadopago.com", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/users/me" {
			fmt.Fprintf(w, `{"id":%d,"country_id":"AR","site_id":"MLA"}`, account)
			return
		}
		if r.URL.Path != "/checkout/preferences" || r.Method != "POST" {
			t.Error("unexpected MP effect")
			w.WriteHeader(500)
			return
		}
		posts++
		io.Copy(io.Discard, r.Body)
		fmt.Fprintf(w, `{"id":"123-preference","external_reference":"order-1","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"collector_id":123,"items":[{"currency_id":"ARS","quantity":1,"unit_price":1250.50}],"sandbox_init_point":"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-preference","expiration_date_to":%q}`, expires.Format(time.RFC3339))
	})
	scope := Scope{TenantID: "tenant", OrganizationID: "store", ConnectionID: "connection", ProviderCode: "mercadopago", AccountRef: "123", Currency: "ARS", MinorUnitExponent: 2}
	d, err := NewSDKDriver(SDKDriverConfig{Scope: scope, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", NotificationURL: "https://api.example.test/payment", DisplayName: "Order"}, "TEST-fixture", client)
	if err != nil {
		t.Fatal(err)
	}
	req := Request{TenantID: "tenant", OrganizationID: "store", ProviderCode: "mercadopago", Currency: "ARS", OrderID: "order-1", PaymentAttemptID: "attempt-1", CustomerSubject: "customer", AmountMinor: 125050, CheckoutExpiresAt: expires}
	result, err := d.CreateCheckout(context.Background(), req)
	if err != nil || result.AccountRef != "123" {
		t.Fatalf("%+v %v", result, err)
	}
	mu.Lock()
	account = 456
	mu.Unlock()
	if _, err = d.CreateCheckout(context.Background(), req); err == nil {
		t.Fatal("foreign account allowed")
	}
	mu.Lock()
	postCount := posts
	mu.Unlock()
	if postCount != 1 {
		t.Fatal("unexpected repeat effect")
	}
}

var _ Driver = (*SDKDriver)(nil)

func TestSDKDriverRejectsMismatchedStripeCredentialModeBeforeHTTP(t *testing.T) {
	config := SDKDriverConfig{Scope: Scope{TenantID: "tenant", OrganizationID: "store", ConnectionID: "connection", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", DisplayName: "Order"}
	client := fixtureHTTPClient(t, "api.stripe.com", func(w http.ResponseWriter, r *http.Request) {
		t.Error("constructor must not perform HTTP")
		w.WriteHeader(500)
	})
	for _, token := range []string{"sk_live_fixture", "rk_live_fixture", "pk_test_fixture", "sk_test_", "sk_test_fixture\n"} {
		if _, err := NewSDKDriver(config, token, client); err == nil {
			t.Error("invalid credential mode accepted")
		}
	}
	for _, token := range []string{"sk_test_fixture", "rk_test_fixture"} {
		if _, err := NewSDKDriver(config, token, client); err != nil {
			t.Fatal(err)
		}
	}
	config.Scope.LiveMode = true
	if _, err := NewSDKDriver(config, "sk_test_fixture", client); err == nil {
		t.Error("test credential accepted for live profile")
	}
	if _, err := NewSDKDriver(config, "rk_live_fixture", client); err != nil {
		t.Fatal(err)
	}
}
````

### FILE: `internal/paymentbridge/webhook.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2dc4f478994076dab3a47e35c44af7b05c4be2b013cf3081de9ac822e1c6d001"
variables: []
secrets_allowed: false
```

````go
package paymentbridge

// AUTHORED composition glue: signature implementations remain in the pinned
// SDK adapter; persistence reuses ProviderIntegration's durable inbox/job owner.
import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"elite.local/enterprise/internal/providerintegration"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
)

var ErrWebhookConfig = errors.New("payment webhook configuration invalid")

type SecretSource interface {
	PaymentWebhookSecret(context.Context, string, string) (string, error)
}
type WebhookConfig struct {
	TenantID, ConnectionID, ProviderCode string
	StripeAccountID                      string
	LiveMode                             bool
	Tolerance                            time.Duration
	MaxConcurrent                        int
	Secrets                              SecretSource
	Inbox                                providerintegration.Repository
}
type Webhook struct {
	config WebhookConfig
	slots  chan struct{}
}

func NewWebhook(c WebhookConfig) (*Webhook, error) {
	valid := func(s string) bool {
		return s != "" && len(s) <= 200 && strings.TrimSpace(s) == s && !strings.ContainsAny(s, "\r\n")
	}
	if !valid(c.TenantID) || !valid(c.ConnectionID) || (c.ProviderCode != "stripe" && c.ProviderCode != "mercadopago") || c.Tolerance <= 0 || c.Tolerance > 10*time.Minute || c.MaxConcurrent < 1 || c.MaxConcurrent > 32 || c.Secrets == nil || c.Inbox == nil {
		return nil, ErrWebhookConfig
	}
	return &Webhook{config: c, slots: make(chan struct{}, c.MaxConcurrent)}, nil
}

// InboxPayload contains no provider raw body, client secret, signature or payer
// data. BodySHA256 in the receipt binds original bytes. The normalized notice
// only queues GET reconciliation; it never authorizes a payment transition.
type InboxPayload struct {
	Schema       string                               `json:"schema"`
	Notification officialpayments.PaymentNotification `json:"notification"`
}

func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fail := func(code int) { w.WriteHeader(code); _, _ = io.WriteString(w, "PAYMENT_WEBHOOK_NOT_ACCEPTED\n") }
	if h == nil {
		fail(503)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		fail(405)
		return
	}
	select {
	case h.slots <- struct{}{}:
		defer func() { <-h.slots }()
	default:
		fail(503)
		return
	}
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" || len(r.Header.Values("Content-Type")) != 1 {
		fail(415)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	_ = http.NewResponseController(w).SetReadDeadline(time.Now().Add(10 * time.Second))
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, officialpayments.MaxPayloadBytes))
	if err != nil {
		fail(413)
		return
	}
	if len(body) == 0 {
		fail(400)
		return
	}
	one := func(key string) (string, bool) {
		values := r.Header.Values(key)
		return r.Header.Get(key), len(values) == 1 && len(values[0]) > 0 && len(values[0]) <= 16384
	}
	var sig, requestID, dataID string
	if h.config.ProviderCode == "stripe" {
		var ok bool
		sig, ok = one("Stripe-Signature")
		if !ok || r.URL.RawQuery != "" {
			fail(400)
			return
		}
	} else {
		var ok bool
		sig, ok = one("X-Signature")
		if !ok {
			fail(400)
			return
		}
		requestID, ok = one("X-Request-Id")
		if !ok || len(r.URL.RawQuery) > 2048 {
			fail(400)
			return
		}
		query, parseErr := url.ParseQuery(r.URL.RawQuery)
		if parseErr != nil || len(query["data.id"]) != 1 || len(query) > 2 {
			fail(400)
			return
		}
		for key, values := range query {
			if (key != "data.id" && key != "type") || len(values) != 1 {
				fail(400)
				return
			}
		}
		dataID = query.Get("data.id")
	}
	secret, err := h.config.Secrets.PaymentWebhookSecret(ctx, h.config.ProviderCode, h.config.ConnectionID)
	if err != nil || secret == "" || len(secret) > 16384 {
		fail(503)
		return
	}
	var n officialpayments.PaymentNotification
	if h.config.ProviderCode == "stripe" {
		n, err = officialpayments.VerifyStripePaymentNotification(body, sig, secret, h.config.Tolerance, h.config.LiveMode, h.config.StripeAccountID)
	} else {
		n, err = officialpayments.VerifyMercadoPagoPaymentNotification(body, sig, requestID, dataID, secret, h.config.Tolerance)
	}
	if err != nil {
		fail(403)
		return
	}
	normalized, err := json.Marshal(InboxPayload{Schema: "elite-payment-notification/v1", Notification: n})
	if err != nil {
		fail(503)
		return
	}
	_, err = h.config.Inbox.AcceptWebhook(ctx, providerintegration.Receipt{TenantID: h.config.TenantID, ConnectionID: h.config.ConnectionID, ProviderCode: n.ProviderCode, ProviderEventID: n.ProviderEventID, EventType: "payment.reconciliation-requested.v1", BodyHash: n.BodySHA256, Payload: normalized})
	if errors.Is(err, providerintegration.ErrConflict) {
		fail(409)
		return
	}
	if err != nil {
		fail(503)
		return
	}
	w.WriteHeader(200)
	_, _ = io.WriteString(w, "PAYMENT_WEBHOOK_ACCEPTED\n")
}
````

### FILE: `internal/paymentbridge/webhook_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5f57fd0684d1f2fe933eaeff636e2d41d9fb9578f8daebc9a6cd9a6878f9edb0"
variables: []
secrets_allowed: false
```

````go
package paymentbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/providerintegration"
	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

type fixtureSecret struct{}

func (fixtureSecret) PaymentWebhookSecret(context.Context, string, string) (string, error) {
	return "whsec_fixture", nil
}

type fixtureInbox struct {
	mu      sync.Mutex
	records map[string]providerintegration.Receipt
	err     error
}

func (f *fixtureInbox) AcceptWebhook(_ context.Context, r providerintegration.Receipt) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return false, f.err
	}
	old, ok := f.records[r.ProviderEventID]
	if ok {
		if old.BodyHash != r.BodyHash {
			return false, providerintegration.ErrConflict
		}
		return true, nil
	}
	f.records[r.ProviderEventID] = r
	return false, nil
}
func signedRequest(t *testing.T, body []byte) *http.Request {
	t.Helper()
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "whsec_fixture", Timestamp: time.Now()})
	r := httptest.NewRequest("POST", "https://fixture.test/webhooks/payment", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Stripe-Signature", signed.Header)
	return r
}
func paymentEvent() []byte {
	return []byte(fmt.Sprintf(`{"id":"evt_fixture","object":"event","api_version":%q,"type":"payment_intent.succeeded","livemode":false,"data":{"object":{"id":"pi_fixture","object":"payment_intent","client_secret":"must-never-persist","payer":{"email":"private@example.test"}}}}`, stripe.APIVersion))
}
func TestWebhookOfficialSignatureToScopedInboxAndReplay(t *testing.T) {
	inbox := &fixtureInbox{records: map[string]providerintegration.Receipt{}}
	h, err := NewWebhook(WebhookConfig{TenantID: "tenant-fixture", ConnectionID: "connection-fixture", ProviderCode: "stripe", Tolerance: 5 * time.Minute, MaxConcurrent: 2, Secrets: fixtureSecret{}, Inbox: inbox})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, signedRequest(t, paymentEvent()))
		if w.Code != 200 {
			t.Fatalf("code%d: %s", w.Code, w.Body.String())
		}
	}
	if len(inbox.records) != 1 {
		t.Fatal("duplicate inbox")
	}
	row := inbox.records["evt_fixture"]
	if row.TenantID != "tenant-fixture" || row.ConnectionID != "connection-fixture" || row.ProviderCode != "stripe" {
		t.Fatal("scope mismatch")
	}
	var normalized InboxPayload
	if json.Unmarshal(row.Payload, &normalized) != nil || normalized.Notification.ProviderReference != "pi_fixture" {
		t.Fatal("notification mismatch")
	}
	if strings.Contains(string(row.Payload), "secret") || strings.Contains(string(row.Payload), "private") || strings.Contains(string(row.Payload), "succeeded\",\"amount") {
		t.Fatal("raw data leaked")
	}
}
func TestWebhookRejectsBeforeInboxAndReportsCommitFailure(t *testing.T) {
	inbox := &fixtureInbox{records: map[string]providerintegration.Receipt{}}
	h, _ := NewWebhook(WebhookConfig{TenantID: "tenant", ConnectionID: "connection", ProviderCode: "stripe", Tolerance: time.Minute, MaxConcurrent: 1, Secrets: fixtureSecret{}, Inbox: inbox})
	cases := []struct {
		name string
		edit func(*http.Request)
		want int
	}{{"missing-signature", func(r *http.Request) { r.Header.Del("Stripe-Signature") }, 400}, {"duplicate-header", func(r *http.Request) { r.Header.Add("Stripe-Signature", r.Header.Get("Stripe-Signature")) }, 400}, {"wrong-signature", func(r *http.Request) { r.Header.Set("Stripe-Signature", "t=1,v1=bad") }, 403}, {"query-scope", func(r *http.Request) { r.URL.RawQuery = "tenant=other" }, 400}, {"method", func(r *http.Request) { r.Method = "GET" }, 405}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := signedRequest(t, paymentEvent())
			tc.edit(r)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("%d", w.Code)
			}
		})
	}
	if len(inbox.records) != 0 {
		t.Fatal("invalid receipt persisted")
	}
	inbox.err = errors.New("database unavailable with private info")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, signedRequest(t, paymentEvent()))
	if w.Code != 503 || strings.Contains(w.Body.String(), "private") {
		t.Fatal("failed commit accepted or leaked")
	}
}
````

### FILE: `internal/platform/httpapi/payment_checkout.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9f4f08f9f594f9cd9ca8b7206a9e29e8b94321aab7b387b7dd55f4e193e5d2bd"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED authorization/transport glue. Customer identity is always taken from
// the verified token; provider redirects come only from the scoped repository.
import (
	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const PaymentProviderWebhookPath = "/v1/payment-provider/webhook"

type PaymentCheckoutModule struct {
	Reader  paymentbridge.CustomerCheckoutReader
	Webhook http.Handler
}

func (m PaymentCheckoutModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	mux.HandleFunc("GET /v1/customer/orders/{id}/checkout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if verifier == nil {
			writeProblem(w, 503, "CHECKOUT_UNAVAILABLE", "checkout is unavailable")
			return
		}
		p, err := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if err != nil || p.Subject == "" || p.TenantID == "" {
			writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
			return
		}
		if !p.Allowed("customer:self") {
			writeProblem(w, 403, "FORBIDDEN", "customer:self permission is required")
			return
		}
		q, err := url.ParseQuery(r.URL.RawQuery)
		organization := q.Get("organization_id")
		order := r.PathValue("id")
		if err != nil || len(r.URL.RawQuery) > 1024 || len(q) != 1 || len(q["organization_id"]) != 1 || order == "" || len(order) > 200 || strings.ContainsAny(order, "/\\\r\n") || strings.TrimSpace(order) != order {
			writeProblem(w, 400, "INVALID_CHECKOUT_QUERY", "one organization scope and order are required")
			return
		}
		if !p.AllowedOrganization(organization) {
			writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "organization scope is required")
			return
		}
		if m.Reader == nil {
			writeProblem(w, 503, "CHECKOUT_UNAVAILABLE", "checkout is unavailable")
			return
		}
		value, err := m.Reader.CustomerCheckout(r.Context(), p.TenantID, organization, p.Subject, order)
		if errors.Is(err, paymentbridge.ErrCheckoutNotAvailable) {
			writeProblem(w, 404, "CHECKOUT_NOT_AVAILABLE", "checkout is not available")
			return
		}
		if err != nil || value.OrderID != order || !value.Valid(time.Now()) {
			writeProblem(w, 503, "CHECKOUT_UNAVAILABLE", "checkout is unavailable")
			return
		}
		writeJSON(w, 200, value)
	})
	if m.Webhook != nil {
		mux.Handle(PaymentProviderWebhookPath, m.Webhook)
	}
}
````

### FILE: `internal/platform/httpapi/payment_checkout_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "74bb67bcbb279fe1c6b43a3cc36f9f68f7f2073545e2c82c5f72e81b9a415655"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/platform/identity"
)

func TestPaymentRequestsCannotLeaveConfiguredHostScope(t *testing.T) {
	for _, tenant := range []string{"tenant", "other"} {
		organization := "other"
		if tenant == "other" {
			organization = "store"
		}
		p := identity.Principal{Subject: "operator", TenantID: tenant, Permissions: map[string]struct{}{"payment:create": {}, "payment:write": {}}, Organizations: map[string]struct{}{organization: {}}}
		mux := http.NewServeMux()
		CommerceModule{PaymentProvider: "stripe", PaymentTenantID: "tenant", PaymentOrganizationID: "store"}.Register(mux, paymentPrincipal{p})
		for _, path := range []string{"/v1/commerce/orders/order/payment-request", "/v1/commerce/orders/order/payments", "/v1/payments/payment/transitions"} {
			request := httptest.NewRequest("POST", path, strings.NewReader(`{"organization_id":"`+organization+`"}`))
			request.Header.Set("Authorization", "Bearer fixture")
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != 503 || !strings.Contains(response.Body.String(), "PAYMENT_SCOPE_DISABLED") {
				t.Fatal(path, response.Code, response.Body.String())
			}
		}
	}
}

type checkoutReadFixture struct {
	calls int
	value paymentbridge.CustomerCheckout
	err   error
	scope []string
}

func (f *checkoutReadFixture) CustomerCheckout(_ context.Context, tenant, org, subject, order string) (paymentbridge.CustomerCheckout, error) {
	f.calls++
	f.scope = []string{tenant, org, subject, order}
	if tenant != "tenant" || subject != "alice" || order != "order" {
		return paymentbridge.CustomerCheckout{}, paymentbridge.ErrCheckoutNotAvailable
	}
	return f.value, f.err
}
func TestCustomerCheckoutHTTPAuthorizationAndProviderURL(t *testing.T) {
	base := identity.Principal{Subject: "alice", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	for _, tc := range []struct {
		name, path, authorization, subject, tenant string
		permission                                 bool
		expected, calls                            int
	}{
		{"allowed", "/v1/customer/orders/order/checkout?organization_id=store", "Bearer fixture", "alice", "tenant", true, 200, 1},
		{"anonymous", "/v1/customer/orders/order/checkout?organization_id=store", "", "alice", "tenant", true, 401, 0},
		{"no-permission", "/v1/customer/orders/order/checkout?organization_id=store", "Bearer fixture", "alice", "tenant", false, 403, 0},
		{"foreign-org", "/v1/customer/orders/order/checkout?organization_id=other", "Bearer fixture", "alice", "tenant", true, 403, 0},
		{"foreign-customer", "/v1/customer/orders/order/checkout?organization_id=store", "Bearer fixture", "bob", "tenant", true, 404, 1},
		{"foreign-tenant", "/v1/customer/orders/order/checkout?organization_id=store", "Bearer fixture", "alice", "other", true, 404, 1},
		{"unknown-order", "/v1/customer/orders/other/checkout?organization_id=store", "Bearer fixture", "alice", "tenant", true, 404, 1},
		{"injected-subject", "/v1/customer/orders/order/checkout?organization_id=store&customer_id=alice", "Bearer fixture", "bob", "tenant", true, 400, 0},
		{"injected-url", "/v1/customer/orders/order/checkout?organization_id=store&url=https://evil.test", "Bearer fixture", "alice", "tenant", true, 400, 0},
		{"duplicate-org", "/v1/customer/orders/order/checkout?organization_id=store&organization_id=store", "Bearer fixture", "alice", "tenant", true, 400, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := base
			p.Subject = tc.subject
			p.TenantID = tc.tenant
			if !tc.permission {
				p.Permissions = nil
			}
			reader := &checkoutReadFixture{value: paymentbridge.CustomerCheckout{OrderID: "order", ProviderCode: "stripe", URL: "https://checkout.stripe.com/c/pay/cs_test_fixture", ExpiresAt: time.Now().Add(time.Hour)}}
			mux := http.NewServeMux()
			PaymentCheckoutModule{Reader: reader}.Register(mux, paymentPrincipal{p})
			request := httptest.NewRequest("GET", tc.path, nil)
			request.Header.Set("Authorization", tc.authorization)
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != tc.expected || reader.calls != tc.calls {
				t.Fatalf("status=%d calls=%d", response.Code, reader.calls)
			}
			if reader.calls > 0 && (reader.scope[0] != p.TenantID || reader.scope[2] != p.Subject) {
				t.Fatal("scope did not come from principal")
			}
			if response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("Referrer-Policy") != "no-referrer" {
				t.Fatal("missing private response policy")
			}
		})
	}
	for _, tc := range []struct {
		name, url string
		err       error
		nilReader bool
	}{
		{"invalid-provider-host", "https://checkout.stripe.com.evil.test/session", nil, false},
		{"unsafe-url", "javascript:alert(1)", nil, false},
		{"reader-error", "https://checkout.stripe.com/session", errors.New("PRIVATE_DATABASE_DETAIL"), false},
		{"disabled", "", nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := &checkoutReadFixture{value: paymentbridge.CustomerCheckout{OrderID: "order", ProviderCode: "stripe", URL: tc.url, ExpiresAt: time.Now().Add(time.Hour)}, err: tc.err}
			module := PaymentCheckoutModule{Reader: reader}
			if tc.nilReader {
				module.Reader = nil
			}
			mux := http.NewServeMux()
			module.Register(mux, paymentPrincipal{base})
			request := httptest.NewRequest("GET", "/v1/customer/orders/order/checkout?organization_id=store", nil)
			request.Header.Set("Authorization", "Bearer fixture")
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != 503 || strings.Contains(response.Body.String(), "PRIVATE_DATABASE_DETAIL") {
				t.Fatal(response.Code, response.Body.String())
			}
		})
	}
}
````

### FILE: `internal/platform/httpapi/payment_observation_authority_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file17:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "51ccf68a05177c18017bfff0ef37f2cf6d16d96cfddd1c83c32b9dff0b62d34a"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/platform/identity"
)

type observationAuthorityRepository struct {
	commerceRepo
	transitions int
}

func (r *observationAuthorityRepository) TransitionPayment(context.Context, string, string, string, string, string, int64, string, string) error {
	r.transitions++
	return nil
}

func TestProviderObservationOwnsFinancialTransitions(t *testing.T) {
	for _, c := range []struct {
		name, current, target, reference string
		status, calls                    int
	}{
		{"operator-cannot-authorize", "pending", "authorized", "pi_fixture", 403, 0},
		{"operator-cannot-capture", "authorized", "captured", "pi_fixture", 403, 0},
		{"operator-cannot-refund", "captured", "refunded", "pi_fixture", 403, 0},
		{"operator-cannot-dispute", "captured", "disputed", "pi_fixture", 403, 0},
		{"operator-cannot-fail-dispatched", "pending", "failed", "", 403, 0},
		{"cancel-does-not-invent-provider-reference", "created", "failed", "pi_forged", 403, 0},
		{"cancel-before-dispatch-remains-supported", "created", "failed", "", 200, 1},
	} {
		t.Run(c.name, func(t *testing.T) {
			repo := &observationAuthorityRepository{}
			p := identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{"payment:write": {}}, Organizations: map[string]struct{}{"org": {}}}
			mux := http.NewServeMux()
			CommerceModule{Service: commerce.NewService(repo, &commerceIDs{}), PaymentProvider: "stripe", ProviderObservedPayments: true}.Register(mux, paymentPrincipal{p})
			body, _ := json.Marshal(map[string]any{"organization_id": "org", "current": c.current, "target": c.target, "version": 1, "provider_reference": c.reference})
			request := httptest.NewRequest("POST", "/v1/payments/attempt/transitions", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer fixture")
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != c.status || repo.transitions != c.calls {
				t.Fatalf("status=%d repository_calls=%d body=%s", response.Code, repo.transitions, response.Body.String())
			}
		})
	}
}
````

### FILE: `internal/platform/postgres/payment_callback_checkout_recovery.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file18:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "365239d0421552a769c242bf6e4c73c4c12213020560b518fdc80653ed7797c3"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED composition around the existing callback/checkout/fence owners. The
// financial observation and preference identity remain separate receipts.
import (
	"context"
	"errors"

	"elite.local/enterprise/internal/paymentbridge"
)

var ErrCheckoutRecovery = errors.New("original checkout identity unresolved; no resend permitted")

type PaymentCallbackRecoveryProcessor struct{ Base *PaymentCallbackProcessor }

func (p PaymentCallbackRecoveryProcessor) Handle(ctx context.Context, j Job) error {
	if p.Base == nil {
		return ErrPaymentCallback
	}
	// The signed callback's payment resource is reconciled through the unchanged
	// processor first. Missing preference identity never fabricates payment facts.
	if err := p.Base.Handle(ctx, j); err != nil {
		return err
	}
	base := p.Base
	c := base.worker.Scope
	if c.ProviderCode != "mercadopago" {
		return nil
	}
	_, notice, err := base.notice(ctx, j)
	if err != nil {
		return err
	}
	var attempt, session, fence string
	err = base.pool.QueryRow(ctx, `select b.payment_attempt_id,coalesce(b.session_id,''),f.state from payment.provider_checkout b
 join payment.provider_observation o using(tenant_id,payment_attempt_id)
 join communication.outbound_delivery f on f.tenant_id=b.tenant_id and f.delivery_key=b.payment_attempt_id and f.channel_code='payment_'||b.provider_code and f.request_sha256_hex=b.request_sha256_hex
 where b.tenant_id=$1 and b.connection_id=$2 and b.provider_code=$3 and b.account_ref=$4 and b.live_mode=$5 and o.provider_reference=$6 and o.observed_at is not null and o.evidence_sha256_hex is not null`, c.TenantID, c.ConnectionID, c.ProviderCode, c.AccountRef, c.LiveMode, notice.ProviderReference).Scan(&attempt, &session, &fence)
	if err != nil {
		return errors.Join(ErrCheckoutRecovery, err)
	}
	if fence == "accepted" && session != "" {
		return nil
	}
	if fence != "unknown" && fence != "sending" {
		return ErrCheckoutRecovery
	}
	r, hash, known, err := base.binding(ctx, attempt)
	if err != nil {
		return err
	}
	var checkout paymentbridge.Checkout
	if known != "" {
		// Recover a crash between SaveCheckout and fence receipt without another
		// search or POST. Even a persisted session is re-observed through SDK GET.
		checkout, err = base.worker.Driver.RetrieveCheckout(ctx, known)
		if err == nil && checkout.SessionID != known {
			return ErrCheckoutRecovery
		}
	} else {
		driver, ok := base.worker.Driver.(paymentbridge.CheckoutRecoveryDriver)
		if !ok {
			return ErrCheckoutRecovery
		}
		checkout, err = driver.RecoverCheckout(ctx, r)
	}
	if err != nil {
		return errors.Join(ErrCheckoutRecovery, err)
	}
	if !checkoutMatches(c, r, checkout) {
		return ErrCheckoutRecovery
	}
	if _, _, err = base.notice(ctx, j); err != nil {
		return err
	}
	if err = base.worker.Store.SaveCheckout(ctx, c, r, hash, checkout); err != nil {
		return err
	}
	return base.recoverFence(ctx, r, hash, checkout)
}
````

### FILE: `internal/platform/postgres/payment_callback_processor.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file19:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e7fe432c2d5c492a817b39cf0828f699955e37cb3d15acf35e4f5987188b884b"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED composition glue. Signed callbacks are resource hints; the pinned
// SDK GET and existing checkout/observation owners establish payment facts.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/paymentbridge"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPaymentCallback = errors.New("payment callback scope or durable receipt invalid")

type PaymentCallbackProcessor struct {
	pool     *pgxpool.Pool
	worker   paymentbridge.Worker
	workerID string
}

func NewPaymentCallbackProcessor(pool *pgxpool.Pool, worker paymentbridge.Worker, workerID string) (*PaymentCallbackProcessor, error) {
	if pool == nil || worker.Scope.Validate() != nil || worker.Store == nil || worker.Driver == nil || worker.Fence == nil || strings.TrimSpace(workerID) == "" || len(workerID) > 128 {
		return nil, ErrPaymentCallback
	}
	return &PaymentCallbackProcessor{pool: pool, worker: worker, workerID: workerID}, nil
}

// Claim selects only this connection. A worker never exhausts another tenant's
// valid payment job merely because its configured credentials differ.
type PaymentCallbackJobs struct {
	*Jobs
	Scope paymentbridge.Scope
}

func (s *PaymentCallbackJobs) Claim(ctx context.Context, queue, worker string, lease time.Duration, limit int) ([]Job, error) {
	if s == nil || s.Jobs == nil || s.Scope.Validate() != nil || queue != "payment-provider-events" {
		return nil, ErrPaymentCallback
	}
	match, _ := json.Marshal(map[string]string{"connection_id": s.Scope.ConnectionID, "provider_code": s.Scope.ProviderCode, "event_type": "payment.reconciliation-requested.v1"})
	return s.Jobs.ClaimScoped(ctx, queue, worker, lease, limit, JobScope{TenantID: s.Scope.TenantID, JobType: "provider.webhook.received", SchemaVersion: 1, PayloadMatch: match})
}

type paymentCallbackJob struct {
	ConnectionID    string `json:"connection_id"`
	ProviderCode    string `json:"provider_code"`
	ProviderEventID string `json:"provider_event_id"`
	EventType       string `json:"event_type"`
}

func decodePaymentCallback(raw []byte, value any) error {
	if len(raw) == 0 || len(raw) > 16384 {
		return ErrPaymentCallback
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(value) != nil || d.Decode(new(any)) != io.EOF {
		return ErrPaymentCallback
	}
	return nil
}

func (p *PaymentCallbackProcessor) notice(ctx context.Context, j Job) (paymentCallbackJob, officialpayments.PaymentNotification, error) {
	var payload paymentCallbackJob
	var n officialpayments.PaymentNotification
	c := p.worker.Scope
	if j.TenantID != c.TenantID || j.Queue != "payment-provider-events" || j.JobType != "provider.webhook.received" || j.SchemaVersion != 1 || j.Attempts < 1 || decodePaymentCallback(j.Payload, &payload) != nil || payload.ConnectionID != c.ConnectionID || payload.ProviderCode != c.ProviderCode || payload.EventType != "payment.reconciliation-requested.v1" || payload.ProviderEventID == "" {
		return payload, n, ErrPaymentCallback
	}
	var raw []byte
	var digest string
	err := p.pool.QueryRow(ctx, `select e.payload,e.body_sha256_hex from platform.job j
 join integration.webhook_event e on e.tenant_id=j.tenant_id and e.connection_id=$9 and e.provider_event_id=$10 and e.provider_code=$11 and e.event_type=$12
 join integration.provider_connection c on c.tenant_id=e.tenant_id and c.connection_id=e.connection_id and c.provider_code=e.provider_code and c.state='active' and (c.organization_id is null or c.organization_id=$13)
 where j.tenant_id=$1 and j.job_id=$2 and j.claimed_by=$3 and j.attempts=$4 and j.queue=$5 and j.job_type=$6 and j.schema_version=$7 and j.payload=$8::jsonb
 and j.completed_at is null and j.terminal_error_code is null and j.claimed_until>clock_timestamp()`, j.TenantID, j.JobID, p.workerID, j.Attempts, j.Queue, j.JobType, j.SchemaVersion, j.Payload, c.ConnectionID, payload.ProviderEventID, c.ProviderCode, payload.EventType, c.OrganizationID).Scan(&raw, &digest)
	if err != nil {
		return payload, n, errors.Join(ErrPaymentCallback, err)
	}
	var inbox paymentbridge.InboxPayload
	if decodePaymentCallback(raw, &inbox) != nil || inbox.Schema != "elite-payment-notification/v1" {
		return payload, n, ErrPaymentCallback
	}
	n = inbox.Notification
	if n.ProviderCode != c.ProviderCode || n.ProviderEventID != payload.ProviderEventID || n.BodySHA256 != digest || len(digest) != 64 || n.ProviderReference == "" || n.Type == "" {
		return payload, n, ErrPaymentCallback
	}
	if _, err = hex.DecodeString(digest); err != nil {
		return payload, n, ErrPaymentCallback
	}
	if (c.ProviderCode == "stripe" && n.ResourceType != "checkout_session" && n.ResourceType != "payment_intent") || (c.ProviderCode == "mercadopago" && n.ResourceType != "payment") {
		return payload, n, ErrPaymentCallback
	}
	return payload, n, nil
}

// binding also accepts a durable pre-dispatch skeleton. This permits recovery
// after a lost POST response without issuing a second provider POST.
func (p *PaymentCallbackProcessor) binding(ctx context.Context, attempt string) (paymentbridge.Request, string, string, error) {
	var r paymentbridge.Request
	var hash, session string
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return r, hash, session, err
	}
	defer tx.Rollback(ctx)
	c := p.worker.Scope
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return r, hash, session, err
	}
	r, err = lockCheckoutRequest(ctx, tx, c, attempt)
	if err != nil {
		return r, hash, session, err
	}
	err = tx.QueryRow(ctx, `select b.request_sha256_hex,coalesce(b.session_id,'') from payment.provider_checkout b
 join communication.outbound_delivery f on f.tenant_id=b.tenant_id and f.delivery_key=b.payment_attempt_id and f.channel_code='payment_'||b.provider_code and f.request_sha256_hex=b.request_sha256_hex and f.state in ('sending','unknown','accepted')
 where b.tenant_id=$1 and b.payment_attempt_id=$2 and b.provider_code=$3 and b.connection_id=$4 and b.account_ref=$5 and b.live_mode=$6 and b.expires_at=$7`, c.TenantID, attempt, c.ProviderCode, c.ConnectionID, c.AccountRef, c.LiveMode, r.CheckoutExpiresAt).Scan(&hash, &session)
	if err != nil {
		return r, hash, session, errors.Join(ErrPaymentCallback, err)
	}
	_, actualHash, err := paymentCallbackMessage(r)
	if err != nil || actualHash != hash {
		return r, hash, session, ErrPaymentCallback
	}
	return r, hash, session, tx.Commit(ctx)
}

func paymentCallbackMessage(r paymentbridge.Request) (channels.Message, string, error) {
	raw, err := json.Marshal(r)
	if err != nil {
		return channels.Message{}, "", err
	}
	m := channels.Message{TenantID: r.TenantID, ChannelCode: "payment_" + r.ProviderCode, ExternalID: r.CustomerSubject, ThreadID: r.OrderID, DeliveryKey: r.PaymentAttemptID, Direction: channels.DirectionOut, Text: string(raw)}
	hash, err := outbounddelivery.MessageSHA256(m)
	return m, hash, err
}

func checkoutMatches(c paymentbridge.Scope, r paymentbridge.Request, v paymentbridge.Checkout) bool {
	return v.ProviderCode == c.ProviderCode && v.AccountRef == c.AccountRef && v.LiveMode == c.LiveMode && v.OrderID == r.OrderID && v.PaymentAttemptID == r.PaymentAttemptID && v.Currency == r.Currency && v.AmountMinor == r.AmountMinor && v.ExpiresAt.Equal(r.CheckoutExpiresAt) && v.SessionID != ""
}

func (p *PaymentCallbackProcessor) Handle(ctx context.Context, j Job) error {
	if p == nil || p.pool == nil {
		return ErrPaymentCallback
	}
	payload, n, err := p.notice(ctx, j)
	if err != nil {
		return err
	}
	c := p.worker.Scope
	attempt, reference := "", n.ProviderReference
	if n.ResourceType == "checkout_session" {
		observed, err := p.worker.Driver.RetrieveCheckout(ctx, reference)
		if err != nil {
			return err
		}
		if observed.SessionID != reference || observed.PaymentReference == "" {
			return ErrPaymentCallback
		}
		attempt = observed.PaymentAttemptID
		r, hash, session, err := p.binding(ctx, attempt)
		if err != nil {
			return err
		}
		if !checkoutMatches(c, r, observed) || (session != "" && session != observed.SessionID) {
			return ErrPaymentCallback
		}
		// Refresh hold before persisting newly learned session/payment linkage.
		// The payment observation itself still comes exclusively from Reconcile GET.
		if _, _, err = p.worker.Store.BeginObservation(ctx, c, attempt, observed.PaymentReference); err != nil {
			return err
		}
		if err = p.worker.Store.SaveCheckout(ctx, c, r, hash, observed); err != nil {
			return err
		}
		if err = p.recoverFence(ctx, r, hash, observed); err != nil {
			return err
		}
		reference = observed.PaymentReference
	} else {
		// The common refund/status path resolves the already bound provider ID.
		var matches int
		err = p.pool.QueryRow(ctx, `select count(*),coalesce(min(b.payment_attempt_id),'') from payment.provider_checkout b left join payment.provider_observation o using(tenant_id,payment_attempt_id)
 where b.tenant_id=$1 and b.connection_id=$2 and b.provider_code=$3 and b.account_ref=$4 and b.live_mode=$5 and (b.payment_reference=$6 or o.provider_reference=$6)`, c.TenantID, c.ConnectionID, c.ProviderCode, c.AccountRef, c.LiveMode, reference).Scan(&matches, &attempt)
		if err != nil {
			return err
		}
		if matches > 1 {
			return ErrPaymentCallback
		}
		if matches == 0 {
			// First payment callback may precede checkout completion. Metadata from an
			// authenticated provider GET locates the attempt; callback body cannot.
			snapshot, getErr := p.worker.Driver.RetrievePayment(ctx, reference)
			if getErr != nil {
				return getErr
			}
			attempt = snapshot.PaymentAttemptID
			// Mercado Pago preferences carry our order external_reference; payment
			// metadata is not guaranteed to inherit the preference metadata. Only
			// the authenticated GET may select a unique previously dispatched
			// whole-order request. Zero/ambiguous matches never pick an attempt.
			if attempt == "" && c.ProviderCode == "mercadopago" {
				var candidates int
				err = p.pool.QueryRow(ctx, `select count(*),coalesce(min(b.payment_attempt_id),'') from payment.provider_checkout b
 join payment.payment_attempt p using(tenant_id,payment_attempt_id)
 join sales.customer_order o using(tenant_id,order_id)
 join communication.outbound_delivery f on f.tenant_id=b.tenant_id and f.delivery_key=b.payment_attempt_id and f.channel_code='payment_'||b.provider_code and f.request_sha256_hex=b.request_sha256_hex and f.state in ('sending','unknown','accepted')
 where b.tenant_id=$1 and b.connection_id=$2 and b.provider_code=$3 and b.account_ref=$4 and b.live_mode=$5 and p.order_id=$6 and o.organization_id=$7 and p.currency=$8 and p.amount_minor_units=$9
 and (b.payment_reference is null or b.payment_reference=$10)`, c.TenantID, c.ConnectionID, c.ProviderCode, c.AccountRef, c.LiveMode, snapshot.OrderID, c.OrganizationID, snapshot.Currency, snapshot.AmountMinor, reference).Scan(&candidates, &attempt)
				if err != nil {
					return err
				}
				if candidates != 1 {
					return ErrPaymentCallback
				}
			}
			r, _, _, bindErr := p.binding(ctx, attempt)
			if bindErr != nil {
				return bindErr
			}
			if snapshot.ProviderCode != c.ProviderCode || snapshot.ProviderReference != reference || snapshot.OrderID != r.OrderID || snapshot.Currency != r.Currency || snapshot.AmountMinor != r.AmountMinor || snapshot.LiveMode != c.LiveMode || (c.ProviderCode == "mercadopago" && snapshot.CollectorID != c.AccountRef) {
				return ErrPaymentCallback
			}
		}
		if _, _, _, err = p.binding(ctx, attempt); err != nil {
			return err
		}
		var referenceMatches bool
		err = p.pool.QueryRow(ctx, `select (b.payment_reference is null or b.payment_reference=$3) and (o.provider_reference is null or o.provider_reference=$3) from payment.provider_checkout b left join payment.provider_observation o using(tenant_id,payment_attempt_id) where b.tenant_id=$1 and b.payment_attempt_id=$2`, c.TenantID, attempt, reference).Scan(&referenceMatches)
		if err != nil || !referenceMatches {
			return ErrPaymentCallback
		}
	}
	if _, _, err = p.notice(ctx, j); err != nil {
		return err
	}
	if err = p.worker.Reconcile(ctx, attempt, reference); err != nil {
		return err
	}
	// A crash before this receipt or the processor's subsequent ACK retries only
	// GETs. Existing observation/outbox owners suppress duplicate money effects.
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var owned bool
	err = tx.QueryRow(ctx, lockedJob+`select claimed_by=$3 and attempts=$4 and claimed_until>clock_timestamp() and payload=$5::jsonb from owned`, j.TenantID, j.JobID, p.workerID, j.Attempts, j.Payload).Scan(&owned)
	if err != nil || !owned {
		return ErrJobClaimLost
	}
	_, err = tx.Exec(ctx, `update integration.webhook_event set state='processed',processed_at=clock_timestamp(),last_error_code=null where tenant_id=$1 and connection_id=$2 and provider_event_id=$3 and provider_code=$4 and body_sha256_hex=$5`, c.TenantID, c.ConnectionID, payload.ProviderEventID, c.ProviderCode, n.BodySHA256)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (p *PaymentCallbackProcessor) recoverFence(ctx context.Context, r paymentbridge.Request, hash string, v paymentbridge.Checkout) error {
	var state string
	if err := p.pool.QueryRow(ctx, `select state from communication.outbound_delivery where tenant_id=$1 and channel_code=$2 and delivery_key=$3 and request_sha256_hex=$4`, r.TenantID, "payment_"+r.ProviderCode, r.PaymentAttemptID, hash).Scan(&state); err != nil {
		return err
	}
	if state == "accepted" {
		return nil
	}
	m, actual, err := paymentCallbackMessage(r)
	if err != nil || actual != hash {
		return ErrPaymentCallback
	}
	evidence, _ := json.Marshal(struct{ Provider, Session, Payment, Status string }{v.ProviderCode, v.SessionID, v.PaymentReference, v.Status})
	sum := sha256.Sum256(evidence)
	receipt := outbounddelivery.Receipt{ProviderMessageID: v.SessionID, EvidenceSHA256: hex.EncodeToString(sum[:]), AcceptedAt: time.Now().UTC()}
	if state == "sending" {
		return p.worker.Fence.Complete(ctx, m, hash, receipt)
	}
	reconciler, ok := p.worker.Fence.(interface {
		ReconcileAccepted(context.Context, channels.Message, string, outbounddelivery.Receipt) error
	})
	if state != "unknown" || !ok {
		return ErrPaymentCallback
	}
	return reconciler.ReconcileAccepted(ctx, m, hash, receipt)
}
````

### FILE: `internal/platform/postgres/payment_checkout.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file20:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d6aeaa87260bdedcda029c87118c0e8802f978e74a5e7134bcaac34e43dcc274"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED persistence/transaction glue around the existing Commerce payment,
// integration inbox and outbound fence; no new accounting or capture algorithm.
import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/providerintegration"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentCheckoutStore struct{ pool *pgxpool.Pool }

func NewPaymentCheckoutStore(pool *pgxpool.Pool) *PaymentCheckoutStore {
	return &PaymentCheckoutStore{pool: pool}
}

func (s *PaymentCheckoutStore) PendingRequests(ctx context.Context, c paymentbridge.Scope, limit int) ([]paymentbridge.Request, error) {
	if s == nil || s.pool == nil || c.Validate() != nil || limit < 1 || limit > 100 {
		return nil, paymentbridge.ErrCheckoutConflict
	}
	rows, err := s.pool.Query(ctx, `select p.tenant_id,o.organization_id,p.payment_attempt_id,p.order_id,o.customer_principal_id,p.provider_code,p.currency,p.amount_minor_units,
 (select date_trunc('second',min(e.occurred_at))+interval '23 hours' from platform.outbox_event e where e.tenant_id=p.tenant_id and e.aggregate_id=p.payment_attempt_id and e.event_type='payment.requested' and e.schema_version=1)
 from payment.payment_attempt p join sales.customer_order o using(tenant_id,order_id)
 join integration.provider_connection conn on conn.tenant_id=p.tenant_id and conn.connection_id=$4 and conn.provider_code=p.provider_code and conn.state='active' and (conn.organization_id is null or conn.organization_id=o.organization_id)
 left join payment.provider_checkout checkout on checkout.tenant_id=p.tenant_id and checkout.payment_attempt_id=p.payment_attempt_id
 where p.tenant_id=$1 and p.provider_code=$2 and o.organization_id=$3 and p.currency=$5 and p.state='created'
 and checkout.session_id is null and p.amount_minor_units=o.total_minor_units and p.currency=o.currency
 and o.state in ('placed','confirmed','allocated') and o.customer_principal_id is not null
 and exists(select 1 from platform.outbox_event e where e.tenant_id=p.tenant_id and e.aggregate_id=p.payment_attempt_id and e.event_type='payment.requested' and e.schema_version=1)
 and (select date_trunc('second',min(e.occurred_at))+interval '23 hours' from platform.outbox_event e where e.tenant_id=p.tenant_id and e.aggregate_id=p.payment_attempt_id and e.event_type='payment.requested' and e.schema_version=1)>clock_timestamp()+interval '31 minutes'
 and not exists(select 1 from communication.outbound_delivery f where f.tenant_id=p.tenant_id and f.channel_code='payment_'||p.provider_code and f.delivery_key=p.payment_attempt_id::text and f.state in ('accepted','unknown','failed_terminal'))
 order by p.updated_at,p.payment_attempt_id limit $6`, c.TenantID, c.ProviderCode, c.OrganizationID, c.ConnectionID, c.Currency, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []paymentbridge.Request{}
	for rows.Next() {
		var r paymentbridge.Request
		if err = rows.Scan(&r.TenantID, &r.OrganizationID, &r.PaymentAttemptID, &r.OrderID, &r.CustomerSubject, &r.ProviderCode, &r.Currency, &r.AmountMinor, &r.CheckoutExpiresAt); err != nil {
			return nil, err
		}
		r.CheckoutExpiresAt = r.CheckoutExpiresAt.UTC()
		result = append(result, r)
	}
	return result, rows.Err()
}

func lockCheckoutRequest(ctx context.Context, tx pgx.Tx, c paymentbridge.Scope, attempt string) (paymentbridge.Request, error) {
	var r paymentbridge.Request
	var order string
	err := tx.QueryRow(ctx, `select order_id from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2 and provider_code=$3`, c.TenantID, attempt, c.ProviderCode).Scan(&order)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, paymentbridge.ErrCheckoutConflict
	}
	if err != nil {
		return r, err
	}
	var organization, customer, currency string
	var total int64
	err = tx.QueryRow(ctx, `select organization_id,customer_principal_id,currency,total_minor_units from sales.customer_order where tenant_id=$1 and order_id=$2 and organization_id=$3 and state in ('placed','confirmed','allocated','delivered') for update`, c.TenantID, order, c.OrganizationID).Scan(&organization, &customer, &currency, &total)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, paymentbridge.ErrCheckoutConflict
	}
	if err != nil {
		return r, err
	}
	err = tx.QueryRow(ctx, `select tenant_id,payment_attempt_id,order_id,provider_code,currency,amount_minor_units from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2 and order_id=$3 for update`, c.TenantID, attempt, order).Scan(&r.TenantID, &r.PaymentAttemptID, &r.OrderID, &r.ProviderCode, &r.Currency, &r.AmountMinor)
	if err != nil {
		return r, err
	}
	if r.ProviderCode != c.ProviderCode || r.Currency != c.Currency || r.Currency != currency || r.AmountMinor != total || r.AmountMinor <= 0 {
		return r, paymentbridge.ErrCheckoutConflict
	}
	r.OrganizationID = organization
	r.CustomerSubject = customer
	err = tx.QueryRow(ctx, `select date_trunc('second',min(occurred_at))+interval '23 hours' from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='payment.requested' and schema_version=1`, c.TenantID, attempt).Scan(&r.CheckoutExpiresAt)
	if err != nil {
		return r, err
	}
	r.CheckoutExpiresAt = r.CheckoutExpiresAt.UTC()
	return r, nil
}

func checkoutConnection(ctx context.Context, tx pgx.Tx, c paymentbridge.Scope) error {
	var found string
	err := tx.QueryRow(ctx, `select connection_id from integration.provider_connection where tenant_id=$1 and connection_id=$2 and provider_code=$3 and state='active' and (organization_id is null or organization_id=$4) for share`, c.TenantID, c.ConnectionID, c.ProviderCode, c.OrganizationID).Scan(&found)
	if errors.Is(err, pgx.ErrNoRows) {
		return providerintegration.ErrConnection
	}
	return err
}

func (s *PaymentCheckoutStore) BindRequest(ctx context.Context, c paymentbridge.Scope, r paymentbridge.Request, hash string) error {
	if c.Validate() != nil || len(hash) != 64 {
		return paymentbridge.ErrCheckoutConflict
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return err
	}
	actual, err := lockCheckoutRequest(ctx, tx, c, r.PaymentAttemptID)
	if err != nil {
		return err
	}
	if actual != r {
		return paymentbridge.ErrCheckoutConflict
	}
	var dispatchable bool
	err = tx.QueryRow(ctx, `select p.state='created' and o.state in ('placed','confirmed','allocated') from payment.payment_attempt p join sales.customer_order o using(tenant_id,order_id) where p.tenant_id=$1 and p.payment_attempt_id=$2`, c.TenantID, r.PaymentAttemptID).Scan(&dispatchable)
	if err != nil {
		return err
	}
	if !dispatchable || !r.CheckoutExpiresAt.After(time.Now().Add(31*time.Minute)) {
		return paymentbridge.ErrCheckoutConflict
	}
	_, err = tx.Exec(ctx, `insert into payment.provider_checkout(tenant_id,payment_attempt_id,provider_code,connection_id,account_ref,request_sha256_hex,live_mode,expires_at) values($1,$2,$3,$4,$5,$6,$7,$8) on conflict do nothing`, c.TenantID, r.PaymentAttemptID, c.ProviderCode, c.ConnectionID, c.AccountRef, hash, c.LiveMode, r.CheckoutExpiresAt)
	if err != nil {
		return err
	}
	var exact bool
	err = tx.QueryRow(ctx, `select provider_code=$3 and connection_id=$4 and account_ref=$5 and request_sha256_hex=$6 and live_mode=$7 and expires_at=$8 from payment.provider_checkout where tenant_id=$1 and payment_attempt_id=$2 for update`, c.TenantID, r.PaymentAttemptID, c.ProviderCode, c.ConnectionID, c.AccountRef, hash, c.LiveMode, r.CheckoutExpiresAt).Scan(&exact)
	if err != nil {
		return err
	}
	if !exact {
		return paymentbridge.ErrCheckoutConflict
	}
	return tx.Commit(ctx)
}

func (s *PaymentCheckoutStore) SaveCheckout(ctx context.Context, c paymentbridge.Scope, r paymentbridge.Request, hash string, result paymentbridge.Checkout) error {
	if result.OrderID != r.OrderID || result.PaymentAttemptID != r.PaymentAttemptID || result.Currency != r.Currency || result.AmountMinor != r.AmountMinor || result.LiveMode != c.LiveMode || result.AccountRef != c.AccountRef || !result.ExpiresAt.Equal(r.CheckoutExpiresAt) {
		return paymentbridge.ErrCheckoutConflict
	}
	// Stripe only returns the hosted URL while a session is active. A completed
	// GET can recover the durable identity after a lost POST response without a URL.
	urlOptional := result.ProviderCode == "stripe" && (result.Status == "complete" || result.Status == "expired")
	if c.Validate() != nil || result.ProviderCode != c.ProviderCode || result.SessionID == "" || len(result.SessionID) > 200 || (result.URL == "" && !urlOptional) || len(result.URL) > 4096 || result.Status == "" {
		return paymentbridge.ErrCheckoutConflict
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return err
	}
	actual, err := lockCheckoutRequest(ctx, tx, c, r.PaymentAttemptID)
	if err != nil {
		return err
	}
	if actual != r {
		return paymentbridge.ErrCheckoutConflict
	}
	// A known session cannot be replaced by a later completion or a different key.
	updated, err := tx.Exec(ctx, `update payment.provider_checkout set session_id=$3,checkout_url=coalesce(nullif($4,''),checkout_url),checkout_state=$5,payment_reference=coalesce(nullif($11,''),payment_reference),updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2 and request_sha256_hex=$6 and provider_code=$7 and connection_id=$8 and account_ref=$9 and live_mode=$10 and (session_id is null or session_id=$3) and (payment_reference is null or $11='' or payment_reference=$11)`, c.TenantID, r.PaymentAttemptID, result.SessionID, result.URL, result.Status, hash, c.ProviderCode, c.ConnectionID, c.AccountRef, c.LiveMode, result.PaymentReference)
	if err != nil {
		return err
	}
	if updated.RowsAffected() != 1 {
		return paymentbridge.ErrCheckoutConflict
	}
	return tx.Commit(ctx)
}

func (s *PaymentCheckoutStore) Checkout(ctx context.Context, c paymentbridge.Scope, attempt string) (paymentbridge.Request, paymentbridge.Checkout, error) {
	var r paymentbridge.Request
	var result paymentbridge.Checkout
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return r, result, err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return r, result, err
	}
	r, err = lockCheckoutRequest(ctx, tx, c, attempt)
	if err != nil {
		return r, result, err
	}
	err = tx.QueryRow(ctx, `select provider_code,session_id,coalesce(checkout_url,''),checkout_state,coalesce(payment_reference,''),expires_at from payment.provider_checkout where tenant_id=$1 and payment_attempt_id=$2 and connection_id=$3 and account_ref=$4 and live_mode=$5 and session_id is not null`, c.TenantID, attempt, c.ConnectionID, c.AccountRef, c.LiveMode).Scan(&result.ProviderCode, &result.SessionID, &result.URL, &result.Status, &result.PaymentReference, &result.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, result, paymentbridge.ErrCheckoutConflict
	}
	if err != nil {
		return r, result, err
	}
	result.OrderID = r.OrderID
	result.PaymentAttemptID = r.PaymentAttemptID
	result.Currency = r.Currency
	result.AmountMinor = r.AmountMinor
	result.AccountRef = c.AccountRef
	result.LiveMode = c.LiveMode
	return r, result, tx.Commit(ctx)
}

func (s *PaymentCheckoutStore) BeginObservation(ctx context.Context, c paymentbridge.Scope, attempt, reference string) (paymentbridge.Request, int64, error) {
	var empty paymentbridge.Request
	if c.Validate() != nil || reference == "" || len(reference) > 200 {
		return empty, 0, paymentbridge.ErrCheckoutConflict
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return empty, 0, err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return empty, 0, err
	}
	r, err := lockCheckoutRequest(ctx, tx, c, attempt)
	if err != nil {
		return r, 0, err
	}
	var bound bool
	err = tx.QueryRow(ctx, `select exists(select 1 from payment.provider_checkout where tenant_id=$1 and payment_attempt_id=$2 and connection_id=$3 and provider_code=$4 and account_ref=$5 and live_mode=$6)`, c.TenantID, attempt, c.ConnectionID, c.ProviderCode, c.AccountRef, c.LiveMode).Scan(&bound)
	if err != nil {
		return r, 0, err
	}
	if !bound {
		return r, 0, paymentbridge.ErrCheckoutConflict
	}
	_, err = tx.Exec(ctx, `insert into payment.provider_observation(tenant_id,payment_attempt_id,order_id,organization_id,provider_code,provider_reference,currency,amount_minor_units,live_mode,account_ref) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) on conflict(tenant_id,payment_attempt_id) do nothing`, c.TenantID, attempt, r.OrderID, r.OrganizationID, c.ProviderCode, reference, r.Currency, r.AmountMinor, c.LiveMode, c.AccountRef)
	if err != nil {
		return r, 0, err
	}
	var generation int64
	err = tx.QueryRow(ctx, `update payment.provider_observation set provider_reference=$4,hold=true,hold_reason='GET_IN_PROGRESS',generation=generation+1,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2 and provider_code=$3 and (provider_reference=$4 or (observed_at is null and hold)) and account_ref=$5 and live_mode=$6 returning generation`, c.TenantID, attempt, c.ProviderCode, reference, c.AccountRef, c.LiveMode).Scan(&generation)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, 0, paymentbridge.ErrCheckoutConflict
	}
	if err != nil {
		return r, 0, err
	}
	return r, generation, tx.Commit(ctx)
}

func (s *PaymentCheckoutStore) RecordObservation(ctx context.Context, c paymentbridge.Scope, r paymentbridge.Request, generation int64, p officialpayments.PaymentSnapshot) error {
	if c.Validate() != nil || generation < 1 {
		return paymentbridge.ErrCheckoutConflict
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return err
	}
	actual, err := lockCheckoutRequest(ctx, tx, c, r.PaymentAttemptID)
	if err != nil {
		return err
	}
	if actual != r {
		return paymentbridge.ErrCheckoutConflict
	}
	var storedGeneration int64
	var previousRefunded int64
	var reference string
	err = tx.QueryRow(ctx, `select generation,provider_reference,refunded_minor_units from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2 for update`, c.TenantID, r.PaymentAttemptID).Scan(&storedGeneration, &reference, &previousRefunded)
	if err != nil {
		return err
	}
	if generation != storedGeneration {
		return paymentbridge.ErrObservationStale
	}
	mismatch := p.ProviderCode != c.ProviderCode || p.ProviderReference != reference || p.OrderID != r.OrderID || p.Currency != r.Currency || p.AmountMinor != r.AmountMinor || p.LiveMode != c.LiveMode || p.ReceivedMinor < 0 || p.ReceivedMinor > p.AmountMinor || p.RefundedMinor < 0 || p.RefundedMinor > p.ReceivedMinor
	if p.ProviderCode == "mercadopago" && p.CollectorID != c.AccountRef {
		mismatch = true
	}
	if p.PaymentAttemptID != "" && p.PaymentAttemptID != r.PaymentAttemptID {
		mismatch = true
	}
	if p.RefundedMinor < previousRefunded {
		mismatch = true
	}
	reason := "PROVIDER_NOT_CAPTURED"
	target := ""
	if mismatch {
		reason = "PROVIDER_BINDING_MISMATCH"
	} else if p.Disputed {
		reason = "PROVIDER_DISPUTED"
		target = "disputed"
	} else if p.RefundedMinor > 0 {
		reason = "PARTIAL_PROVIDER_REFUND"
		if p.RefundedMinor == r.AmountMinor && p.ReceivedMinor == r.AmountMinor {
			reason = "PROVIDER_REFUNDED"
			target = "refunded"
		}
	} else if (p.ProviderCode == "stripe" && p.Status == "succeeded" || p.ProviderCode == "mercadopago" && p.Status == "approved") && p.ReceivedMinor == r.AmountMinor {
		reason = "MATCHED_CAPTURE"
		target = "captured"
	}
	var financialState, boundReference string
	if err = tx.QueryRow(ctx, `select state,coalesce(provider_reference,'') from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2`, c.TenantID, r.PaymentAttemptID).Scan(&financialState, &boundReference); err != nil {
		return err
	}
	if boundReference != "" && boundReference != reference {
		mismatch = true
		target = ""
		reason = "BOUND_PROVIDER_REFERENCE_MISMATCH"
	}
	if (financialState == "refunded" || financialState == "disputed" || financialState == "failed") && target != "" && target != financialState {
		target = ""
		reason = "TERMINAL_STATE_REQUIRES_REVIEW"
	}
	hold := reason != "MATCHED_CAPTURE"
	evidence := p.SHA256()
	if mismatch {
		_, err = tx.Exec(ctx, `update payment.provider_observation set hold=true,hold_reason=$3,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2`, c.TenantID, r.PaymentAttemptID, reason)
	} else {
		_, err = tx.Exec(ctx, `update payment.provider_observation set received_minor_units=$3,refunded_minor_units=$4,provider_status=$5,evidence_sha256_hex=$6,observed_at=clock_timestamp(),hold=$7,hold_reason=$8,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2`, c.TenantID, r.PaymentAttemptID, p.ReceivedMinor, p.RefundedMinor, p.Status, evidence, hold, reason)
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into payment.provider_observation_event(tenant_id,payment_attempt_id,generation,evidence_sha256_hex,result_code) values($1,$2,$3,$4,$5)`, c.TenantID, r.PaymentAttemptID, generation, evidence, reason)
	if err != nil {
		return err
	}
	// These are provider vocabulary projections, not simulated authorization or
	// capture operations. One observed successful payment may skip earlier notices.
	if !mismatch && target != "" {
		_, err = tx.Exec(ctx, `with changed as(update payment.payment_attempt set state=$3,provider_reference=$4,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2 and (state<>$3 or provider_reference is distinct from $4) returning version)
   insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
   select $1,gen_random_uuid(),'payment',$2,version,'payment.provider-observed',1,clock_timestamp(),jsonb_build_object('state',$3::text,'evidence_sha256',$5::text,'provider_code',$6::text) from changed`, c.TenantID, r.PaymentAttemptID, target, reference, evidence, c.ProviderCode)
		if err != nil {
			return err
		}
	}
	reconciliationState := "matched"
	if mismatch {
		reconciliationState = "amount-mismatch"
	} else if hold {
		reconciliationState = "state-mismatch"
	}
	body, _ := json.Marshal(map[string]any{"reason": reason, "generation": generation, "evidence_sha256": evidence})
	_, err = tx.Exec(ctx, `insert into integration.reconciliation_item(tenant_id,reconciliation_id,provider_code,resource_type,internal_id,external_id,state,evidence,observed_at) values($1,gen_random_uuid(),$2,'payment',$3,$4,$5,$6,clock_timestamp())`, c.TenantID, c.ProviderCode, r.PaymentAttemptID, reference, reconciliationState, body)
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	if mismatch {
		return paymentbridge.ErrPaymentMismatch
	}
	return nil
}

// AcceptWebhook shares the established durable inbox transaction. Only a new
// authenticated notification advances a known observation generation/hold.
func (s *PaymentCheckoutStore) AcceptWebhook(ctx context.Context, value providerintegration.Receipt) (bool, error) {
	var body paymentbridge.InboxPayload
	if value.EventType != "payment.reconciliation-requested.v1" || json.Unmarshal(value.Payload, &body) != nil || body.Schema != "elite-payment-notification/v1" || body.Notification.ProviderCode != value.ProviderCode || body.Notification.ProviderEventID != value.ProviderEventID || body.Notification.BodySHA256 != value.BodyHash {
		return false, providerintegration.ErrInvalid
	}
	n := body.Notification
	if strings.TrimSpace(n.ProviderReference) == "" {
		return false, providerintegration.ErrInvalid
	}
	return NewProviderIntegration(s.pool).acceptWebhookWithEffect(ctx, value, "payment-provider-events", func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `update payment.provider_observation o set hold=true,hold_reason='PROVIDER_NOTIFICATION_PENDING',generation=o.generation+1,updated_at=clock_timestamp() from payment.provider_checkout c where o.tenant_id=$1 and o.provider_code=$2 and (o.provider_reference=$3 or c.session_id=$3) and c.tenant_id=o.tenant_id and c.payment_attempt_id=o.payment_attempt_id and c.connection_id=$4`, value.TenantID, value.ProviderCode, n.ProviderReference, value.ConnectionID)
		return err
	})
}
````

### FILE: `internal/platform/postgres/payment_checkout_integration_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file21:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7a104fbe2d840a8031bffd1ae724935fccae9aa5077ff685f62b9738fda71765"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/providerintegration"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type checkoutTransportFixture struct {
	calls        atomic.Int64
	observed     officialpayments.PaymentSnapshot
	loseResponse bool
}

func (f *checkoutTransportFixture) CreateCheckout(_ context.Context, r paymentbridge.Request) (paymentbridge.Checkout, error) {
	f.calls.Add(1)
	if f.loseResponse {
		return paymentbridge.Checkout{}, context.DeadlineExceeded
	}
	return paymentbridge.Checkout{ProviderCode: "stripe", SessionID: "cs_fixture", URL: "https://checkout.stripe.com/c/pay/cs_fixture", Status: "open", OrderID: r.OrderID, PaymentAttemptID: r.PaymentAttemptID, Currency: r.Currency, AmountMinor: r.AmountMinor, AccountRef: "acct_fixture", LiveMode: false, ExpiresAt: r.CheckoutExpiresAt}, nil
}
func (f *checkoutTransportFixture) RetrieveCheckout(context.Context, string) (paymentbridge.Checkout, error) {
	return paymentbridge.Checkout{ProviderCode: "stripe", SessionID: "cs_fixture", PaymentReference: "pi_fixture", Status: "complete"}, nil
}
func (f *checkoutTransportFixture) RetrievePayment(context.Context, string) (officialpayments.PaymentSnapshot, error) {
	return f.observed, nil
}

func TestPaymentCheckoutDurableProjection(t *testing.T) {
	raw := os.Getenv("PAYMENT_BRIDGE_DB_URL")
	if raw == "" {
		t.Skip("PAYMENT_BRIDGE_DB_URL is not set")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_payment_") {
		t.Fatal("requires a dedicated loopback elite_payment_ database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28471"
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	// Synthetic business inputs only. Payment is created through the actual owner.
	must(`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'payment-checkout','Fixture','Fixture')`, tenant)
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Fixture','store')`, tenant)
	must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','placed','ARS',1000,1)`, tenant)
	must(`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref)values($1,'checkout','stripe','FIXTURE_ONLY_NOT_A_SECRET')`, tenant)
	actual, err := NewCommerce(pool).RecordOrderPayment(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28472", "payment-key-00000001", commerce.PaymentAttempt{ID: "payment-attempt-0001", OrderID: "order", OrganizationID: "store", ProviderCode: "stripe", Currency: "ARS", AmountMinorUnits: 1000}, "operator-fixture")
	if err != nil || actual.State != "created" {
		t.Fatal(actual, err)
	}
	scope := paymentbridge.Scope{TenantID: tenant, OrganizationID: "store", ConnectionID: "checkout", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2, LiveMode: false}
	store := NewPaymentCheckoutStore(pool)
	fence, err := NewOutboundDeliveryStore(pool, []byte(strings.Repeat("f", 32)), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	transport := &checkoutTransportFixture{observed: officialpayments.PaymentSnapshot{ProviderCode: "stripe", ProviderReference: "pi_fixture", OrderID: "order", PaymentAttemptID: actual.ID, Currency: "ARS", AmountMinor: 1000, ReceivedMinor: 1000, Status: "succeeded", LiveMode: false}}
	worker := paymentbridge.Worker{Scope: scope, Store: store, Fence: &PaymentDispatchFence{OutboundDeliveryStore: fence, Scope: scope}, Driver: transport}
	t.Run("concurrent-request-one-call", func(t *testing.T) {
		var wg sync.WaitGroup
		out := make(chan error, 8)
		for range 8 {
			wg.Add(1)
			go func() { defer wg.Done(); _, e := worker.ProcessOnce(ctx, 8); out <- e }()
		}
		wg.Wait()
		close(out)
		for e := range out {
			var pg *pgconn.PgError
			if e != nil && !errors.Is(e, paymentbridge.ErrCheckoutConflict) && (!errors.As(e, &pg) || pg.Code != "40001") {
				t.Error(e)
			}
		}
		if transport.calls.Load() != 1 {
			t.Fatal("duplicate checkout effects", transport.calls.Load())
		}
		if _, e := worker.ProcessOnce(ctx, 8); e != nil {
			t.Fatal(e)
		}
		r, result, e := store.Checkout(ctx, scope, actual.ID)
		if e != nil || r.OrderID != "order" || result.SessionID != "cs_fixture" {
			t.Fatal(r, result, e)
		}
	})
	t.Run("cross-scope-request-denied", func(t *testing.T) {
		c := scope
		c.OrganizationID = "other"
		if _, _, e := store.Checkout(ctx, c, actual.ID); !errors.Is(e, paymentbridge.ErrCheckoutConflict) {
			t.Fatal(e)
		}
	})
	t.Run("capture-only-from-matched-observation", func(t *testing.T) {
		if e := worker.Reconcile(ctx, actual.ID, "pi_fixture"); e != nil {
			t.Fatal(e)
		}
		var state string
		var hold bool
		var received int64
		e := pool.QueryRow(ctx, `select p.state,o.hold,o.received_minor_units from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) where p.tenant_id=$1 and p.payment_attempt_id=$2`, tenant, actual.ID).Scan(&state, &hold, &received)
		if e != nil || state != "captured" || hold || received != 1000 {
			t.Fatal(state, hold, received, e)
		}
	})
	t.Run("partial-refund-preserves-payment-balance-owner", func(t *testing.T) {
		r, generation, e := store.BeginObservation(ctx, scope, actual.ID, "pi_fixture")
		if e != nil {
			t.Fatal(e)
		}
		snapshot := transport.observed
		snapshot.RefundedMinor = 100
		if e = store.RecordObservation(ctx, scope, r, generation, snapshot); e != nil {
			t.Fatal(e)
		}
		var state string
		var hold bool
		e = pool.QueryRow(ctx, `select p.state,o.hold from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) where p.tenant_id=$1 and p.payment_attempt_id=$2`, tenant, actual.ID).Scan(&state, &hold)
		if e != nil || state != "captured" || !hold {
			t.Fatal("partial refund became total", state, hold, e)
		}
	})
	t.Run("callback-hold-and-generation-fence", func(t *testing.T) {
		r, generation, e := store.BeginObservation(ctx, scope, actual.ID, "pi_fixture")
		if e != nil {
			t.Fatal(e)
		}
		notice := officialpayments.PaymentNotification{ProviderCode: "stripe", ProviderEventID: "evt_refund", ProviderReference: "pi_fixture", Type: "payment_intent.succeeded", BodySHA256: strings.Repeat("a", 64)}
		payload, _ := json.Marshal(paymentbridge.InboxPayload{Schema: "elite-payment-notification/v1", Notification: notice})
		receipt := providerintegration.Receipt{TenantID: tenant, ConnectionID: "checkout", ProviderCode: "stripe", ProviderEventID: notice.ProviderEventID, EventType: "payment.reconciliation-requested.v1", BodyHash: notice.BodySHA256, Payload: payload}
		replay, e := store.AcceptWebhook(ctx, receipt)
		if e != nil || replay {
			t.Fatal(replay, e)
		}
		if e = store.RecordObservation(ctx, scope, r, generation, transport.observed); !errors.Is(e, paymentbridge.ErrObservationStale) {
			t.Fatal("stale GET cleared hold", e)
		}
		var before, after int64
		var hold bool
		if e = pool.QueryRow(ctx, `select generation,hold from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2`, tenant, actual.ID).Scan(&before, &hold); e != nil || !hold {
			t.Fatal(before, hold, e)
		}
		replay, e = store.AcceptWebhook(ctx, receipt)
		if e != nil || !replay {
			t.Fatal(replay, e)
		}
		if e = pool.QueryRow(ctx, `select generation from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2`, tenant, actual.ID).Scan(&after); e != nil || before != after {
			t.Fatal(before, after, e)
		}
		var jobs int
		if e = pool.QueryRow(ctx, `select count(*) from platform.job where tenant_id=$1 and queue='payment-provider-events'`, tenant).Scan(&jobs); e != nil || jobs != 1 {
			t.Fatal(jobs, e)
		}
	})
	t.Run("mismatch-refund-dispute-never-clear-hold", func(t *testing.T) {
		for _, variant := range []string{"amount", "currency", "order", "attempt", "mode", "refund", "dispute"} {
			t.Run(variant, func(t *testing.T) {
				r, generation, e := store.BeginObservation(ctx, scope, actual.ID, "pi_fixture")
				if e != nil {
					t.Fatal(e)
				}
				snapshot := transport.observed
				switch variant {
				case "amount":
					snapshot.AmountMinor++
				case "currency":
					snapshot.Currency = "USD"
				case "order":
					snapshot.OrderID = "other"
				case "attempt":
					snapshot.PaymentAttemptID = "other"
				case "mode":
					snapshot.LiveMode = true
				case "refund":
					snapshot.RefundedMinor = 1000
				case "dispute":
					snapshot.Disputed = true
					snapshot.RefundedMinor = 1000
				}
				e = store.RecordObservation(ctx, scope, r, generation, snapshot)
				if variant != "refund" && variant != "dispute" && !errors.Is(e, paymentbridge.ErrPaymentMismatch) {
					t.Fatal(e)
				}
				if (variant == "refund" || variant == "dispute") && e != nil {
					t.Fatal(e)
				}
				var hold bool
				if e = pool.QueryRow(ctx, `select hold from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2`, tenant, actual.ID).Scan(&hold); e != nil || !hold {
					t.Fatal(hold, e)
				}
			})
		}
	})
	t.Run("terminal-state-never-regresses-to-captured", func(t *testing.T) {
		var before string
		if e := pool.QueryRow(ctx, `select state from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2`, tenant, actual.ID).Scan(&before); e != nil || (before != "refunded" && before != "disputed") {
			t.Fatal(before, e)
		}
		if e := worker.Reconcile(ctx, actual.ID, "pi_fixture"); e != nil && !errors.Is(e, paymentbridge.ErrPaymentMismatch) {
			t.Fatal(e)
		}
		var state string
		var hold bool
		e := pool.QueryRow(ctx, `select p.state,o.hold from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) where p.tenant_id=$1 and p.payment_attempt_id=$2`, tenant, actual.ID).Scan(&state, &hold)
		if e != nil || state != before || !hold {
			t.Fatal("terminal financial state regressed", state, hold, e)
		}
	})
	t.Run("connection-organization-is-authoritative", func(t *testing.T) {
		must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'other-store','other-store','Other fixture','store')`, tenant)
		must(`update integration.provider_connection set organization_id='other-store' where tenant_id=$1 and connection_id='checkout'`, tenant)
		defer must(`update integration.provider_connection set organization_id=null where tenant_id=$1 and connection_id='checkout'`, tenant)
		if _, _, e := store.Checkout(ctx, scope, actual.ID); !errors.Is(e, providerintegration.ErrConnection) {
			t.Fatal("connection for a different organization was accepted", e)
		}
		if e := worker.Reconcile(ctx, actual.ID, "pi_fixture"); !errors.Is(e, providerintegration.ErrConnection) {
			t.Fatal("cross-organization observation was accepted", e)
		}
	})
	t.Run("cancel-after-bind-rolls-back-dispatch", func(t *testing.T) {
		// Persist an old, never-dispatched request as fixture input. It must not
		// starve new work and must not be silently retried or marked financially paid.
		must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order-old-request','store','customer','placed','ARS',1000,1)`, tenant)
		must(`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,idempotency_key,state,currency,amount_minor_units,version)values($1,'payment-old-request','order-old-request','stripe','old-request-fixture','created','ARS',1000,1)`, tenant)
		must(`insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,gen_random_uuid(),'payment','payment-old-request',1,'payment.requested',1,clock_timestamp()-interval '24 hours','{}')`, tenant)
		must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order-cancelled-fixture','store','customer','placed','ARS',1000,1)`, tenant)
		attempt, e := NewCommerce(pool).RecordOrderPayment(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28473", "payment-key-cancelled", commerce.PaymentAttempt{ID: "payment-cancelled-fixture", OrderID: "order-cancelled-fixture", OrganizationID: "store", ProviderCode: "stripe", Currency: "ARS", AmountMinorUnits: 1000}, "operator-fixture")
		if e != nil {
			t.Fatal(e)
		}
		requests, e := store.PendingRequests(ctx, scope, 10)
		if e != nil || len(requests) != 1 || requests[0].PaymentAttemptID != attempt.ID {
			t.Fatal(requests, e)
		}
		r := requests[0]
		body, e := json.Marshal(r)
		if e != nil {
			t.Fatal(e)
		}
		message := channels.Message{ChannelCode: "payment_stripe", TenantID: tenant, ExternalID: r.CustomerSubject, ThreadID: r.OrderID, DeliveryKey: r.PaymentAttemptID, Direction: channels.DirectionOut, Text: string(body)}
		hash, e := outbounddelivery.MessageSHA256(message)
		if e != nil {
			t.Fatal(e)
		}
		if e = store.BindRequest(ctx, scope, r, hash); e != nil {
			t.Fatal(e)
		}
		// The operator's cancellation commits after preparation but before the send claim.
		must(`update payment.payment_attempt set state='failed',version=version+1 where tenant_id=$1 and payment_attempt_id=$2`, tenant, attempt.ID)
		if _, e = worker.Fence.Claim(ctx, message, hash); !errors.Is(e, paymentbridge.ErrCheckoutConflict) {
			t.Fatal("stale dispatch was accepted", e)
		}
		var state string
		var fences, events int
		e = pool.QueryRow(ctx, `select state,(select count(*) from communication.outbound_delivery f where f.tenant_id=p.tenant_id and f.delivery_key=p.payment_attempt_id),(select count(*) from platform.outbox_event e where e.tenant_id=p.tenant_id and e.aggregate_id=p.payment_attempt_id and e.event_type='payment.checkout-requested') from payment.payment_attempt p where tenant_id=$1 and payment_attempt_id=$2`, tenant, attempt.ID).Scan(&state, &fences, &events)
		if e != nil || state != "failed" || fences != 0 || events != 0 {
			t.Fatal("rejected dispatch left side effects", state, fences, events, e)
		}
	})
	t.Run("disabled-connection-no-observation", func(t *testing.T) {
		must(`update integration.provider_connection set state='disabled' where tenant_id=$1 and connection_id='checkout'`, tenant)
		if e := worker.Reconcile(ctx, actual.ID, "pi_fixture"); !errors.Is(e, providerintegration.ErrConnection) {
			t.Fatal(e)
		}
	})
	// This database is discarded by the runner. Immutable audit rows are never
	// disabled or deleted merely to make a test's cleanup convenient.
}
````

### FILE: `internal/platform/postgres/payment_checkout_read.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file22:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c05f4b63f75e16bbb9180e937e82b4ce88a4434c92c71bfa2fecace8b7264ec8"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED scoped read composition over existing payment/order owners.
import (
	"context"
	"net/url"
	"strings"
	"time"

	"elite.local/enterprise/internal/paymentbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerCheckoutReader struct {
	pool  *pgxpool.Pool
	scope paymentbridge.Scope
}

func NewCustomerCheckoutReader(pool *pgxpool.Pool, scope paymentbridge.Scope) (*CustomerCheckoutReader, error) {
	if pool == nil || scope.Validate() != nil {
		return nil, paymentbridge.ErrCheckoutConflict
	}
	return &CustomerCheckoutReader{pool, scope}, nil
}
func (s *CustomerCheckoutReader) CustomerCheckout(ctx context.Context, tenant, organization, subject, order string) (paymentbridge.CustomerCheckout, error) {
	var result paymentbridge.CustomerCheckout
	if s == nil || s.pool == nil || tenant != s.scope.TenantID || organization != s.scope.OrganizationID || strings.TrimSpace(subject) == "" || order == "" {
		return result, paymentbridge.ErrCheckoutNotAvailable
	}
	rows, err := s.pool.Query(ctx, `select o.order_id,b.provider_code,b.checkout_url,b.expires_at
 from sales.customer_order o
 join crm.customer_profile customer on customer.tenant_id=o.tenant_id and customer.customer_principal_id=o.customer_principal_id and customer.status='active'
 join payment.payment_attempt p on p.tenant_id=o.tenant_id and p.order_id=o.order_id and p.provider_code=$5 and p.state='pending' and p.amount_minor_units=o.total_minor_units and p.currency=o.currency
 join payment.provider_checkout b on b.tenant_id=p.tenant_id and b.payment_attempt_id=p.payment_attempt_id and b.provider_code=p.provider_code and b.connection_id=$6 and b.account_ref=$7 and b.live_mode=$8 and b.session_id is not null and b.checkout_url is not null and b.checkout_url<>'' and b.expires_at>clock_timestamp() and b.checkout_state in ('open','preference-created')
 join integration.provider_connection c on c.tenant_id=b.tenant_id and c.connection_id=b.connection_id and c.provider_code=b.provider_code and c.state='active' and (c.organization_id is null or c.organization_id=o.organization_id)
 where o.tenant_id=$1 and o.organization_id=$2 and o.customer_principal_id=$3 and o.order_id=$4 and o.state in ('placed','confirmed','allocated') and o.currency=$9
 limit 2`, tenant, organization, subject, order, s.scope.ProviderCode, s.scope.ConnectionID, s.scope.AccountRef, s.scope.LiveMode, s.scope.Currency)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
		if err = rows.Scan(&result.OrderID, &result.ProviderCode, &result.URL, &result.ExpiresAt); err != nil {
			return paymentbridge.CustomerCheckout{}, err
		}
	}
	if err = rows.Err(); err != nil {
		return paymentbridge.CustomerCheckout{}, err
	}
	if count != 1 || !result.Valid(time.Now()) {
		return paymentbridge.CustomerCheckout{}, paymentbridge.ErrCheckoutNotAvailable
	}
	if s.scope.ProviderCode == "mercadopago" {
		u, _ := url.Parse(result.URL)
		expected := "sandbox.mercadopago.com.ar"
		if s.scope.LiveMode {
			expected = "www.mercadopago.com.ar"
		}
		if !strings.EqualFold(u.Hostname(), expected) {
			return paymentbridge.CustomerCheckout{}, paymentbridge.ErrCheckoutNotAvailable
		}
	}
	return result, nil
}
````

### FILE: `internal/platform/postgres/payment_checkout_read_integration_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file23:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "44bf43ec7dd3d97d755994965d3c5ad43018f82abc021207d66501e9630ffe34"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/paymentbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPaymentCustomerCheckoutAuthorizationPostgres(t *testing.T) {
	raw := os.Getenv("PAYMENT_BRIDGE_DB_URL")
	if raw == "" {
		t.Skip("PAYMENT_BRIDGE_DB_URL is not set")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_payment_") {
		t.Fatal("requires dedicated loopback elite_payment_ database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28771"
	must := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	must(`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'payment-customer-read','Fixture','Fixture')`, tenant)
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Fixture','store')`, tenant)
	must(`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name)values($1,'alice','Fixture')`, tenant)
	must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','alice','placed','ARS',1000,1)`, tenant)
	must(`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref)values($1,'checkout','stripe','FIXTURE_ONLY_NOT_A_SECRET')`, tenant)
	_, err = NewCommerce(pool).RecordOrderPayment(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28772", "payment-read-key-00001", commerce.PaymentAttempt{ID: "read-payment-attempt", OrderID: "order", OrganizationID: "store", ProviderCode: "stripe", Currency: "ARS", AmountMinorUnits: 1000}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	scope := paymentbridge.Scope{TenantID: tenant, OrganizationID: "store", ConnectionID: "checkout", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}
	fence, err := NewOutboundDeliveryStore(pool, []byte(strings.Repeat("f", 32)), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	worker := paymentbridge.Worker{Scope: scope, Store: NewPaymentCheckoutStore(pool), Fence: &PaymentDispatchFence{OutboundDeliveryStore: fence, Scope: scope}, Driver: &checkoutTransportFixture{}}
	if n, err := worker.ProcessOnce(ctx, 1); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	reader, err := NewCustomerCheckoutReader(pool, scope)
	if err != nil {
		t.Fatal(err)
	}
	if value, err := reader.CustomerCheckout(ctx, tenant, "store", "alice", "order"); err != nil || !value.Valid(time.Now()) {
		t.Fatal(value, err)
	}
	for _, args := range [][4]string{{"018f4d4a-7b36-7a21-8d10-2f4c54c28773", "store", "alice", "order"}, {tenant, "other", "alice", "order"}, {tenant, "store", "bob", "order"}, {tenant, "store", "alice", "other"}} {
		if _, err := reader.CustomerCheckout(ctx, args[0], args[1], args[2], args[3]); !errors.Is(err, paymentbridge.ErrCheckoutNotAvailable) {
			t.Fatal("foreign read allowed", err)
		}
	}
	must(`update crm.customer_profile set status='restricted' where tenant_id=$1`, tenant)
	if _, err := reader.CustomerCheckout(ctx, tenant, "store", "alice", "order"); !errors.Is(err, paymentbridge.ErrCheckoutNotAvailable) {
		t.Fatal("restricted customer allowed")
	}
	must(`update crm.customer_profile set status='active' where tenant_id=$1`, tenant)
	must(`update payment.provider_checkout set checkout_url='https://attacker.example/session' where tenant_id=$1`, tenant)
	if _, err := reader.CustomerCheckout(ctx, tenant, "store", "alice", "order"); !errors.Is(err, paymentbridge.ErrCheckoutNotAvailable) {
		t.Fatal("unsafe stored URL allowed")
	}
}
````

### FILE: `internal/platform/postgres/payment_checkout_recovery_integration_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file24:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6874b569386e8a7192c8bb3e05f8e39dafdbb23250bd12cb78a628906c757329"
variables: []
secrets_allowed: false
```

````go
package postgres_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/paymentbridge"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/workers"
)

type recoveryFenceCrash struct {
	*db.PaymentDispatchFence
	fail bool
}

func (f *recoveryFenceCrash) ReconcileAccepted(ctx context.Context, m channels.Message, hash string, r outbounddelivery.Receipt) error {
	if f.fail {
		f.fail = false
		return errors.New("fixture crash before accepted receipt")
	}
	return f.PaymentDispatchFence.ReconcileAccepted(ctx, m, hash, r)
}

func runPreferenceRecovery(t *testing.T, scenario string) {
	t.Helper()
	ctx := context.Background()
	pool := connectedPool(t)
	tenant, payment := connectedInputs(t, pool, "mercadopago")
	var mu sync.Mutex
	posts, paymentGets, searchGets, preferenceGets := 0, 0, 0, 0
	expires := time.Time{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer TEST-fixture" {
			t.Error("wrong credential")
		}
		switch r.URL.Path {
		case "/users/me":
			if r.Method != "GET" {
				t.Error("unexpected write")
			}
			io.WriteString(w, `{"id":123,"country_id":"AR","site_id":"MLA"}`)
		case "/checkout/preferences":
			if r.Method != "POST" {
				t.Error("unexpected method")
			}
			posts++
			var body struct {
				ExternalReference string         `json:"external_reference"`
				Metadata          map[string]any `json:"metadata"`
				ExpirationDateTo  time.Time      `json:"expiration_date_to"`
			}
			if json.NewDecoder(r.Body).Decode(&body) != nil || body.ExternalReference != "order" || body.Metadata["payment_attempt_id"] != payment.ID || body.ExpirationDateTo.IsZero() {
				t.Error("request binding missing")
			}
			expires = body.ExpirationDateTo
			// Provider accepted the request but its response never reaches the worker
			// as a valid SDK result. Recovery must not issue another POST.
			io.WriteString(w, `{"incomplete":`)
		case "/v1/payments/987":
			paymentGets++
			if r.Method != "GET" {
				t.Error("unexpected payment effect")
			}
			io.WriteString(w, `{"id":987,"transaction_amount":1234.56,"transaction_amount_refunded":0,"currency_id":"ARS","external_reference":"order","status":"approved","captured":true,"collector_id":123,"live_mode":false}`)
		case "/checkout/preferences/search":
			searchGets++
			q := r.URL.Query()
			if r.Method != "GET" || len(q) != 4 || q.Get("external_reference") != "order" || q.Get("site_id") != "MLA" || q.Get("limit") != "2" || q.Get("offset") != "0" {
				t.Error("unbounded or unscoped search")
			}
			if scenario == "missing" {
				io.WriteString(w, `{"total":0,"next_offset":0,"elements":[]}`)
				return
			}
			if scenario == "ambiguous" {
				io.WriteString(w, `{"total":2,"next_offset":2,"elements":[{"id":"123-first"},{"id":"123-second"}]}`)
				return
			}
			collector, mode, expiration := 123, false, expires
			if scenario == "account" {
				collector = 456
			}
			if scenario == "mode" {
				mode = true
			}
			if scenario == "expiry" {
				expiration = expiration.Add(time.Second)
			}
			fmt.Fprintf(w, `{"total":1,"next_offset":1,"elements":[{"id":"123-preference","external_reference":"order","collector_id":%d,"site_id":"MLA","live_mode":%t,"expires":true,"expiration_date_to":%q}]}`, collector, mode, expiration.Format(time.RFC3339))
		case "/checkout/preferences/123-preference":
			preferenceGets++
			if r.Method != "GET" {
				t.Error("preference mutation")
			}
			attempt, amount := payment.ID, 1234.56
			if scenario == "metadata" {
				attempt = "other-attempt"
			}
			if scenario == "amount" {
				amount = 1234.57
			}
			fmt.Fprintf(w, `{"id":"123-preference","external_reference":"order","metadata":{"order_id":"order","payment_attempt_id":%q},"collector_id":123,"site_id":"MLA","expires":true,"items":[{"currency_id":"ARS","quantity":1,"unit_price":%.2f}],"sandbox_init_point":"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-preference","expiration_date_to":%q}`, attempt, amount, expires.Format(time.RFC3339))
		default:
			t.Errorf("unexpected HTTP path %s", r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	t.Cleanup(server.Close)
	destination, _ := url.Parse(server.URL)
	scope := paymentbridge.Scope{TenantID: tenant, OrganizationID: "store", ConnectionID: "checkout", ProviderCode: "mercadopago", AccountRef: "123", Currency: "ARS", MinorUnitExponent: 2}
	driver, err := paymentbridge.NewSDKDriver(paymentbridge.SDKDriverConfig{Scope: scope, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", NotificationURL: "https://api.example.test/payment", DisplayName: "Synthetic order"}, "TEST-fixture", &http.Client{Transport: connectedTransport{destination, "api.mercadopago.com"}})
	if err != nil {
		t.Fatal(err)
	}
	store := db.NewPaymentCheckoutStore(pool)
	baseFence, err := db.NewOutboundDeliveryStore(pool, []byte(strings.Repeat("r", 32)), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	fence := &recoveryFenceCrash{PaymentDispatchFence: &db.PaymentDispatchFence{OutboundDeliveryStore: baseFence, Scope: scope}, fail: scenario == "saved-crash"}
	worker := paymentbridge.Worker{Scope: scope, Store: store, Fence: fence, Driver: driver}
	if n, err := worker.ProcessOnce(ctx, 1); err != nil || n != 0 {
		t.Fatal("lost response", n, err)
	}
	callback, err := db.NewPaymentCallbackProcessor(pool, worker, "mp-recovery")
	if err != nil {
		t.Fatal(err)
	}
	processor := workers.JobProcessor{Store: &db.PaymentCallbackJobs{Jobs: db.NewJobs(pool), Scope: scope}, Handler: db.PaymentCallbackRecoveryProcessor{Base: callback}, Queue: "payment-provider-events", WorkerID: "mp-recovery", Lease: 20 * time.Second, RetryDelay: 0, BatchSize: 1}
	webhook, err := paymentbridge.NewWebhook(paymentbridge.WebhookConfig{TenantID: tenant, ConnectionID: "checkout", ProviderCode: "mercadopago", Tolerance: time.Minute, MaxConcurrent: 1, Secrets: connectedSecret{}, Inbox: store})
	if err != nil {
		t.Fatal(err)
	}
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	mac := hmac.New(sha256.New, []byte("whsec_fixture"))
	fmt.Fprintf(mac, "id:987;request-id:mp-recovery;ts:%s;", ts)
	request := httptest.NewRequest("POST", "https://api.example.test/payment?data.id=987&type=payment", strings.NewReader(`{"type":"payment","data":{"id":"987"}}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Request-Id", "mp-recovery")
	request.Header.Set("X-Signature", "ts="+ts+",v1="+hex.EncodeToString(mac.Sum(nil)))
	response := httptest.NewRecorder()
	webhook.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal(response.Code)
	}
	completed, err := processor.ProcessOnce(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var identity, state, financial string
	snapshot := func() {
		t.Helper()
		if err := pool.QueryRow(ctx, `select coalesce(b.session_id,''),f.state,p.state from payment.provider_checkout b join payment.payment_attempt p using(tenant_id,payment_attempt_id) join communication.outbound_delivery f on f.tenant_id=b.tenant_id and f.delivery_key=b.payment_attempt_id and f.channel_code='payment_mercadopago' where b.tenant_id=$1`, tenant).Scan(&identity, &state, &financial); err != nil {
			t.Fatal(err)
		}
	}
	snapshot()
	if financial != "captured" {
		t.Fatal("payment not observed through GET", financial)
	}
	switch scenario {
	case "unique":
		if completed != 1 || identity != "123-preference" || state != "accepted" {
			t.Fatal("identity recovery", completed, identity, state)
		}
	case "saved-crash":
		if completed != 0 || identity != "123-preference" || state != "unknown" {
			t.Fatal("crash point", completed, identity, state)
		}
		if completed, err = processor.ProcessOnce(ctx); err != nil || completed != 1 {
			t.Fatal("saved identity retry", completed, err)
		}
		snapshot()
		if state != "accepted" {
			t.Fatal("recovered receipt missing", state)
		}
	default:
		if completed != 0 || identity != "" || state != "unknown" {
			t.Fatal("unproven identity accepted", completed, identity, state)
		}
	}
	if n, err := worker.ProcessOnce(ctx, 1); err != nil || n != 0 {
		t.Fatal("recovery resent checkout", n, err)
	}
	var moneyEffects int
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.provider-observed'`, tenant).Scan(&moneyEffects); err != nil || moneyEffects != 1 {
		t.Fatal("duplicate money observation effect", moneyEffects, err)
	}
	mu.Lock()
	postCount, payCount, searchCount, getCount := posts, paymentGets, searchGets, preferenceGets
	mu.Unlock()
	if postCount != 1 || searchCount != 1 {
		t.Fatal("unexpected provider mutations/search", postCount, searchCount)
	}
	if scenario == "unique" && (payCount != 2 || getCount != 1) {
		t.Fatal("SDK query chain", payCount, getCount)
	}
	if scenario == "saved-crash" && (payCount != 3 || getCount != 2) {
		t.Fatal("retry should GET saved session without repeat search", payCount, getCount)
	}
	t.Logf("MP_PREFERENCE_RECOVERY_PASS scenario=%s posts=%d payment_gets=%d search_gets=%d preference_gets=%d state=%s", scenario, postCount, payCount, searchCount, getCount, state)
}

func TestMercadoPagoPreferenceRecoveryLostPost(t *testing.T) { runPreferenceRecovery(t, "unique") }
func TestMercadoPagoPreferenceRecoveryAfterSavedIdentityCrash(t *testing.T) {
	runPreferenceRecovery(t, "saved-crash")
}
func TestMercadoPagoPreferenceRecoveryRejectsUnprovenIdentity(t *testing.T) {
	for _, scenario := range []string{"missing", "ambiguous", "metadata", "amount", "account", "mode", "expiry"} {
		t.Run(scenario, func(t *testing.T) { runPreferenceRecovery(t, scenario) })
	}
}
````

### FILE: `internal/platform/postgres/payment_connected_integration_test.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file25:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3e91c0e323360965b6bb284a6df95c4c8d40d355a060a29f655570419a7aed04"
variables: []
secrets_allowed: false
```

````go
package postgres_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/paymentbridge"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/platform/workers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

type connectedTransport struct {
	destination *url.URL
	host        string
}

func (s connectedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Scheme != "https" || r.URL.Host != s.host {
		return nil, fmt.Errorf("fixture denies remote destination")
	}
	copy := r.Clone(r.Context())
	u := *r.URL
	u.Scheme = s.destination.Scheme
	u.Host = s.destination.Host
	copy.URL = &u
	copy.Host = ""
	return http.DefaultTransport.RoundTrip(copy)
}

type connectedSecret struct{}

func (connectedSecret) PaymentWebhookSecret(context.Context, string, string) (string, error) {
	return "whsec_fixture", nil
}

type connectedProvider struct {
	mu                              sync.Mutex
	t                               *testing.T
	attempt, account, order         string
	expires                         int64
	posts, sessionGets, paymentGets int
	refund                          int64
	loseResponse                    bool
}

func (f *connectedProvider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer sk_test_fixture" {
		f.t.Error("missing exact fixture credential")
	}
	switch r.URL.Path {
	case "/v1/account":
		if r.Method != "GET" {
			f.t.Error("account write")
		}
		fmt.Fprintf(w, `{"id":%q}`, f.account)
	case "/v1/checkout/sessions", "/v1/checkout/sessions/cs_test_fixture":
		if r.Method == "POST" {
			f.posts++
			if err := r.ParseForm(); err != nil {
				f.t.Error(err)
			}
			f.expires, _ = strconv.ParseInt(r.Form.Get("expires_at"), 10, 64)
			if r.Form.Get("metadata[payment_attempt_id]") != f.attempt || r.Form.Get("payment_intent_data[metadata][payment_attempt_id]") != f.attempt || r.Form.Get("metadata[order_id]") != "order" || r.Form.Get("line_items[0][price_data][unit_amount]") != "123456" || r.Header.Get("Idempotency-Key") != f.attempt {
				f.t.Error("checkout escaped order contract")
			}
			if f.loseResponse {
				_, _ = io.WriteString(w, `{"truncated":`)
				return
			}
		} else if r.Method == "GET" {
			f.sessionGets++
		} else {
			f.t.Error("invalid checkout method")
		}
		intent := `"pi_fixture"`
		checkoutURL := "null"
		if r.Method == "POST" {
			intent = "null"
			checkoutURL = `"https://checkout.stripe.com/c/pay/cs_test_fixture"`
		}
		fmt.Fprintf(w, `{"id":"cs_test_fixture","object":"checkout.session","mode":"payment","amount_total":123456,"currency":"ars","client_reference_id":%q,"metadata":{"order_id":%q,"payment_attempt_id":%q},"status":"complete","payment_status":"paid","livemode":false,"url":%s,"payment_intent":%s,"expires_at":%d}`, f.attempt, f.order, f.attempt, checkoutURL, intent, f.expires)
	case "/v1/payment_intents/pi_fixture":
		f.paymentGets++
		if r.Method != "GET" || r.URL.Query().Get("expand[0]") != "latest_charge" {
			f.t.Error("payment must be expanded official GET")
		}
		fmt.Fprintf(w, `{"id":"pi_fixture","amount":123456,"amount_received":123456,"currency":"ars","status":"succeeded","livemode":false,"metadata":{"order_id":%q,"payment_attempt_id":%q},"latest_charge":{"id":"ch_fixture","amount":123456,"amount_captured":123456,"amount_refunded":%d,"currency":"ars","captured":true,"disputed":false,"livemode":false}}`, f.order, f.attempt, f.refund)
	default:
		f.t.Errorf("unrecognized provider HTTP %s", r.URL.Path)
		w.WriteHeader(500)
	}
}

func connectedPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("PAYMENT_CONNECTED_DB_URL not set")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_payment_connected_") {
		t.Fatal("requires disposable loopback connected database")
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}
func connectedInputs(t *testing.T, pool *pgxpool.Pool, provider string) (string, commerce.PaymentAttempt) {
	t.Helper()
	ctx := context.Background()
	ids := randomid.Generator{}
	tenant := ids.New()
	statements := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'connected-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','fixture@example.test')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','model','new','fixture','{}')`,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'retail','AR','ARS',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','variant',123456,'inclusive')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SERIAL-SYNTHETIC','available',1,clock_timestamp())`,
	}
	for _, q := range statements {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, `insert into integration.provider_connection(tenant_id,connection_id,provider_code,organization_id,secret_ref)values($1,'checkout',$2,'store','FIXTURE_ONLY_NOT_A_SECRET')`, tenant, provider); err != nil {
		t.Fatal(err)
	}
	repo := db.NewFranchiseJourney(pool)
	if _, _, err := repo.CreateQuoteAs(ctx, tenant, "quote-key", franchisejourney.Quote{ID: "quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().UTC().Add(time.Hour)}, strings.Repeat("a", 64), ids.New(), "operator"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 1, strings.Repeat("a", 64), "order", "line", ids.New(), ids.New()); err != nil {
		t.Fatal(err)
	}
	sales := commerce.NewService(db.NewCommerce(pool), ids)
	if err := sales.AllocateStockAs(ctx, tenant, "store", "order", "line", "stock", 1, 1, "operator"); err != nil {
		t.Fatal(err)
	}
	payment, err := sales.RequestOrderPayment(ctx, tenant, "store", "order", provider, "connected-payment-key-0001", "operator")
	if err != nil || payment.State != "created" || payment.AmountMinorUnits != 123456 {
		t.Fatal(payment, err)
	}
	replay, err := sales.RequestOrderPayment(ctx, tenant, "store", "order", provider, "connected-payment-key-0001", "operator")
	if err != nil || replay.ID != payment.ID {
		t.Fatal("request replay", replay, err)
	}
	return tenant, payment
}

type connectedRun struct {
	pool      *pgxpool.Pool
	tenant    string
	payment   commerce.PaymentAttempt
	provider  *connectedProvider
	worker    paymentbridge.Worker
	handler   *db.PaymentCallbackProcessor
	processor workers.JobProcessor
	webhook   http.Handler
}

func newConnectedRun(t *testing.T, loss bool) *connectedRun {
	t.Helper()
	ctx := context.Background()
	pool := connectedPool(t)
	tenant, payment := connectedInputs(t, pool, "stripe")
	f := &connectedProvider{t: t, attempt: payment.ID, account: "acct_fixture", order: "order", loseResponse: loss}
	server := httptest.NewServer(f)
	t.Cleanup(server.Close)
	destination, _ := url.Parse(server.URL)
	c := paymentbridge.Scope{TenantID: tenant, OrganizationID: "store", ConnectionID: "checkout", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}
	driver, err := paymentbridge.NewSDKDriver(paymentbridge.SDKDriverConfig{Scope: c, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", DisplayName: "Synthetic order"}, "sk_test_fixture", &http.Client{Transport: connectedTransport{destination, "api.stripe.com"}})
	if err != nil {
		t.Fatal(err)
	}
	store := db.NewPaymentCheckoutStore(pool)
	fence, err := db.NewOutboundDeliveryStore(pool, []byte(strings.Repeat("f", 32)), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	worker := paymentbridge.Worker{Scope: c, Store: store, Fence: &db.PaymentDispatchFence{OutboundDeliveryStore: fence, Scope: c}, Driver: driver}
	handler, err := db.NewPaymentCallbackProcessor(pool, worker, "connected-callback")
	if err != nil {
		t.Fatal(err)
	}
	jobs := &db.PaymentCallbackJobs{Jobs: db.NewJobs(pool), Scope: c}
	processor := workers.JobProcessor{Store: jobs, Handler: handler, Queue: "payment-provider-events", WorkerID: "connected-callback", Lease: 20 * time.Second, RetryDelay: time.Second, BatchSize: 1}
	webhook, err := paymentbridge.NewWebhook(paymentbridge.WebhookConfig{TenantID: tenant, ConnectionID: "checkout", ProviderCode: "stripe", Tolerance: time.Minute, MaxConcurrent: 4, Secrets: connectedSecret{}, Inbox: store})
	if err != nil {
		t.Fatal(err)
	}
	n, err := worker.ProcessOnce(ctx, 1)
	if err != nil || (!loss && n != 1) || (loss && n != 0) {
		t.Fatalf("dispatch %d %v", n, err)
	}
	if n, err = worker.ProcessOnce(ctx, 1); err != nil || n != 0 {
		t.Fatal("duplicate dispatch", n, err)
	}
	var state string
	if err = pool.QueryRow(ctx, `select state from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2`, tenant, payment.ID).Scan(&state); err != nil || state != "pending" {
		t.Fatal("unobserved money state", state, err)
	}
	return &connectedRun{pool: pool, tenant: tenant, payment: payment, provider: f, worker: worker, handler: handler, processor: processor, webhook: webhook}
}
func (r *connectedRun) callback(t *testing.T, event, kind, resource string, valid bool) int {
	t.Helper()
	object := "checkout.session"
	extra := ""
	if kind == "payment_intent.succeeded" {
		object = "payment_intent"
	}
	if kind == "charge.refunded" {
		object = "charge"
		extra = `,"payment_intent":"pi_fixture"`
	}
	body := []byte(fmt.Sprintf(`{"id":%q,"object":"event","api_version":%q,"type":%q,"livemode":false,"data":{"object":{"id":%q,"object":%q%s}}}`, event, stripe.APIVersion, kind, resource, object, extra))
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "whsec_fixture", Timestamp: time.Now()})
	signature := signed.Header
	if !valid {
		signature += "tampered"
	}
	request := httptest.NewRequest("POST", "https://api.example.test/payment", bytesReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Stripe-Signature", signature)
	response := httptest.NewRecorder()
	r.webhook.ServeHTTP(response, request)
	return response.Code
}
func bytesReader(raw []byte) *strings.Reader { return strings.NewReader(string(raw)) }

func TestPaymentConnectedReference(t *testing.T) {
	for _, loss := range []bool{false, true} {
		t.Run(fmt.Sprintf("lost-post-response=%t", loss), func(t *testing.T) {
			r := newConnectedRun(t, loss)
			ctx := context.Background()
			if code := r.callback(t, "evt_bad", "checkout.session.completed", "cs_test_fixture", false); code != 403 {
				t.Fatal("bad signature", code)
			}
			if code := r.callback(t, "evt_initial", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
				t.Fatal("callback", code)
			}
			if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
				t.Fatal("callback job", n, err)
			}
			var state, hash, fence string
			var hold bool
			if err := r.pool.QueryRow(ctx, `select p.state,o.evidence_sha256_hex,o.hold,f.state from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) join communication.outbound_delivery f on f.tenant_id=p.tenant_id and f.delivery_key=p.payment_attempt_id where p.tenant_id=$1`, r.tenant).Scan(&state, &hash, &hold, &fence); err != nil || state != "captured" || hold || fence != "accepted" || len(hash) != 64 {
				t.Fatal("observed projection", state, hold, fence, err)
			}
			if code := r.callback(t, "evt_initial", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
				t.Fatal("callback replay", code)
			}
			if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 0 {
				t.Fatal("replay job effect", n, err)
			}
			repo := db.NewFranchiseJourney(r.pool)
			sum := sha256.Sum256([]byte(franchisejourney.ReferenceHandoverContractDocument))
			contract := franchisejourney.HandoverReleaseContract{ID: "reference-single-unit-observed-payment-v1", DocumentSHA256: hex.EncodeToString(sum[:]), Scope: "LOCAL_FIXTURES", MaximumObservationAge: 5 * time.Minute}
			handover, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, contract)
			if err != nil {
				t.Fatal(err)
			}
			command := franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: r.payment.ID, ObservationSHA256: hash, IdempotencyKey: "connected-handover-key"}
			prepared, replayed, err := handover.Prepare(ctx, r.tenant, "operator", command)
			if err != nil || replayed {
				t.Fatal("prepare", err)
			}
			recovered, err := handover.Result(ctx, r.tenant, "store", "order", command.IdempotencyKey)
			if err != nil || recovered.Handover.ID != prepared.Handover.ID {
				t.Fatal("recover prepare", err)
			}
			ids := randomid.Generator{}
			if _, err = repo.PublishDeliveryChecklist(ctx, r.tenant, "operator", franchisejourney.DeliveryChecklist{ID: "initial-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
				t.Fatal(err)
			}
			presented, err := repo.CompleteDeliveryChecklist(ctx, r.tenant, "store", "operator", prepared.Handover.ID, 1, "initial-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
			if err != nil {
				t.Fatal(err)
			}
			accepted, err := repo.AcceptHandover(ctx, r.tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "initial-checklist", 1, strings.Repeat("e", 64), ids.New())
			if err != nil || accepted.State != "accepted" {
				t.Fatal(err)
			}
			eligible, err := handover.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash)
			if err != nil || !eligible.Eligible {
				t.Fatal("release reference", err)
			}
			r.provider.mu.Lock()
			posts, sessionGets, paymentGets := r.provider.posts, r.provider.sessionGets, r.provider.paymentGets
			r.provider.refund = 1
			r.provider.mu.Unlock()
			if posts != 1 || sessionGets != 1 || paymentGets != 1 {
				t.Fatalf("HTTP effects posts=%d checkoutgets=%d paymentgets=%d", posts, sessionGets, paymentGets)
			}
			if code := r.callback(t, "evt_refund", "charge.refunded", "ch_fixture", true); code != 200 {
				t.Fatal(code)
			}
			if _, err = handover.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash); !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatal("callback must hold before GET", err)
			}
			if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
				t.Fatal("refund reconciliation", n, err)
			}
			if _, err = handover.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash); !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatal("partial refund released", err)
			}
			var count int
			if err = r.pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.provider-observed'`, r.tenant).Scan(&count); err != nil || count != 1 {
				t.Fatal("duplicate money projection", count, err)
			}
			t.Logf("PAYMENT_CONNECTED_REFERENCE_PASS tenant=%s lost_post_response=%t quote_accept_allocate_request_sdk_signed_callback_job_get_handover_accept_release=true provider_posts=1 refund_hold=true production_claim=false", r.tenant, loss)
		})
	}
}

func TestPaymentCallbackRejectsForgedJobAndMetadata(t *testing.T) {
	r := newConnectedRun(t, false)
	ctx := context.Background()
	if err := r.handler.Handle(ctx, db.Job{TenantID: r.tenant, JobID: randomid.Generator{}.New(), Queue: "payment-provider-events", JobType: "provider.webhook.received", SchemaVersion: 1, Attempts: 1, Payload: json.RawMessage(`{"connection_id":"checkout","provider_code":"stripe","provider_event_id":"evt_fake","event_type":"payment.reconciliation-requested.v1"}`)}); err == nil {
		t.Fatal("unclaimed invented job accepted")
	}
	if code := r.callback(t, "evt_wrong", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	r.provider.mu.Lock()
	r.provider.order = "foreign-order"
	r.provider.mu.Unlock()
	if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 0 {
		t.Fatal("metadata mismatch accepted", n, err)
	}
	var captured int
	if err := r.pool.QueryRow(ctx, `select count(*) from payment.payment_attempt where tenant_id=$1 and state='captured'`, r.tenant).Scan(&captured); err != nil || captured != 0 {
		t.Fatal("foreign order paid", captured, err)
	}
	r.provider.mu.Lock()
	gets := r.provider.paymentGets
	posts := r.provider.posts
	r.provider.mu.Unlock()
	if gets != 0 || posts != 1 {
		t.Fatal("foreign metadata reached payment or POST", gets, posts)
	}
	t.Log("PAYMENT_CALLBACK_NEGATIVES_PASS forged_claim=true metadata_binding=true no_money_effect=true")
}

func TestPaymentCallbackFirstIntentAndClaimFencing(t *testing.T) {
	r := newConnectedRun(t, false)
	ctx := context.Background()
	if code := r.callback(t, "evt_first_intent", "payment_intent.succeeded", "pi_fixture", true); code != 200 {
		t.Fatal(code)
	}
	jobs := r.processor.Store
	claimed, err := jobs.Claim(ctx, "payment-provider-events", r.processor.WorkerID, 20*time.Second, 1)
	if err != nil || len(claimed) != 1 {
		t.Fatal("claim", err)
	}
	original := claimed[0]
	for name, edit := range map[string]func(*db.Job){
		"previous-generation": func(j *db.Job) { j.Attempts++ },
		"other-tenant":        func(j *db.Job) { j.TenantID = randomid.Generator{}.New() },
		"different-queue":     func(j *db.Job) { j.Queue = "provider-events" },
		"changed-payload": func(j *db.Job) {
			j.Payload = json.RawMessage(`{"connection_id":"checkout","provider_code":"stripe","provider_event_id":"evt_not_received","event_type":"payment.reconciliation-requested.v1"}`)
		},
	} {
		t.Run(name, func(t *testing.T) {
			modified := original
			edit(&modified)
			if err := r.handler.Handle(ctx, modified); err == nil {
				t.Fatal("invalid claim accepted")
			}
		})
	}
	r.provider.mu.Lock()
	gets := r.provider.paymentGets
	r.provider.mu.Unlock()
	if gets != 0 {
		t.Fatal("invalid claims issued GET")
	}
	if err = r.handler.Handle(ctx, original); err != nil {
		t.Fatal("first PI callback", err)
	}
	if err = jobs.Complete(ctx, original.TenantID, original.JobID, r.processor.WorkerID, original.Attempts); err != nil {
		t.Fatal(err)
	}
	if err = r.handler.Handle(ctx, original); err == nil {
		t.Fatal("completed claim accepted")
	}
	r.provider.mu.Lock()
	gets = r.provider.paymentGets
	posts := r.provider.posts
	r.provider.mu.Unlock()
	if gets != 2 || posts != 1 {
		t.Fatal("initial metadata GET and fresh reconciliation GET", gets, posts)
	}
	var state string
	var evidence int
	if err = r.pool.QueryRow(ctx, `select state,(select count(*) from payment.provider_observation_event where tenant_id=$1) from payment.payment_attempt where tenant_id=$1`, r.tenant).Scan(&state, &evidence); err != nil || state != "captured" || evidence != 1 {
		t.Fatal("first PI projection", state, evidence, err)
	}
	t.Log("PAYMENT_CALLBACK_FIRST_INTENT_AND_CLAIM_FENCING_PASS metadata_lookup_get=true fresh_reconciliation_get=true claim_fields=4 completed_claim_rejected=true")
}

func TestPaymentConnectedMercadoPagoWithoutInheritedMetadata(t *testing.T) {
	pool := connectedPool(t)
	ctx := context.Background()
	tenant, payment := connectedInputs(t, pool, "mercadopago")
	var mu sync.Mutex
	posts, gets := 0, 0
	externalReference := "order"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer TEST-fixture" {
			t.Error("MP credential mismatch")
		}
		switch r.URL.Path {
		case "/users/me":
			if r.Method != "GET" {
				t.Error("user write")
			}
			io.WriteString(w, `{"id":123,"country_id":"AR","site_id":"MLA"}`)
		case "/checkout/preferences":
			if r.Method != "POST" {
				t.Error("unexpected preference method")
			}
			posts++
			var body struct {
				ExternalReference string         `json:"external_reference"`
				Metadata          map[string]any `json:"metadata"`
				ExpirationDateTo  string         `json:"expiration_date_to"`
			}
			if json.NewDecoder(r.Body).Decode(&body) != nil || body.ExternalReference != "order" || body.Metadata["payment_attempt_id"] != payment.ID || body.ExpirationDateTo == "" {
				t.Error("MP request binding missing")
			}
			fmt.Fprintf(w, `{"id":"123-preference","external_reference":"order","metadata":{"order_id":"order","payment_attempt_id":%q},"collector_id":123,"items":[{"currency_id":"ARS","quantity":1,"unit_price":1234.56}],"sandbox_init_point":"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-preference","expiration_date_to":%q}`, payment.ID, body.ExpirationDateTo)
		case "/v1/payments/987":
			if r.Method != "GET" {
				t.Error("payment write")
			}
			gets++
			// Deliberately no metadata: a payment need not inherit preference metadata.
			fmt.Fprintf(w, `{"id":987,"transaction_amount":1234.56,"transaction_amount_refunded":0,"currency_id":"ARS","external_reference":%q,"status":"approved","captured":true,"collector_id":123,"live_mode":false}`, externalReference)
		default:
			t.Errorf("unexpected MP endpoint %s", r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	t.Cleanup(server.Close)
	destination, _ := url.Parse(server.URL)
	scope := paymentbridge.Scope{TenantID: tenant, OrganizationID: "store", ConnectionID: "checkout", ProviderCode: "mercadopago", AccountRef: "123", Currency: "ARS", MinorUnitExponent: 2}
	driver, err := paymentbridge.NewSDKDriver(paymentbridge.SDKDriverConfig{Scope: scope, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", NotificationURL: "https://api.example.test/payment", DisplayName: "Synthetic order"}, "TEST-fixture", &http.Client{Transport: connectedTransport{destination, "api.mercadopago.com"}})
	if err != nil {
		t.Fatal(err)
	}
	store := db.NewPaymentCheckoutStore(pool)
	base, err := db.NewOutboundDeliveryStore(pool, []byte(strings.Repeat("m", 32)), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	worker := paymentbridge.Worker{Scope: scope, Store: store, Fence: &db.PaymentDispatchFence{OutboundDeliveryStore: base, Scope: scope}, Driver: driver}
	if n, err := worker.ProcessOnce(ctx, 1); err != nil || n != 1 {
		t.Fatal("MP checkout", n, err)
	}
	handler, err := db.NewPaymentCallbackProcessor(pool, worker, "mp-callback")
	if err != nil {
		t.Fatal(err)
	}
	processor := workers.JobProcessor{Store: &db.PaymentCallbackJobs{Jobs: db.NewJobs(pool), Scope: scope}, Handler: handler, Queue: "payment-provider-events", WorkerID: "mp-callback", Lease: 20 * time.Second, RetryDelay: time.Second, BatchSize: 1}
	webhook, err := paymentbridge.NewWebhook(paymentbridge.WebhookConfig{TenantID: tenant, ConnectionID: "checkout", ProviderCode: "mercadopago", Tolerance: time.Minute, MaxConcurrent: 1, Secrets: connectedSecret{}, Inbox: store})
	if err != nil {
		t.Fatal(err)
	}
	callback := func(requestID string) {
		t.Helper()
		ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
		mac := hmac.New(sha256.New, []byte("whsec_fixture"))
		fmt.Fprintf(mac, "id:987;request-id:%s;ts:%s;", requestID, ts)
		// Body financial claims are unsigned and must not supply payment facts.
		request := httptest.NewRequest("POST", "https://api.example.test/payment?data.id=987&type=payment", strings.NewReader(`{"type":"payment","data":{"id":"987"},"status":"rejected","transaction_amount":999999}`))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Request-Id", requestID)
		request.Header.Set("X-Signature", "ts="+ts+",v1="+hex.EncodeToString(mac.Sum(nil)))
		response := httptest.NewRecorder()
		webhook.ServeHTTP(response, request)
		if response.Code != 200 {
			t.Fatal("MP signed callback", response.Code)
		}
	}
	// A genuine signed payment ID with a different GET order cannot bind.
	mu.Lock()
	externalReference = "foreign-order"
	mu.Unlock()
	callback("mp-wrong-order")
	if n, err := processor.ProcessOnce(ctx); err != nil || n != 0 {
		t.Fatal("foreign MP order accepted", n, err)
	}
	var captures int
	if err = pool.QueryRow(ctx, `select count(*) from payment.payment_attempt where tenant_id=$1 and state='captured'`, tenant).Scan(&captures); err != nil || captures != 0 {
		t.Fatal("foreign MP capture", captures, err)
	}
	mu.Lock()
	externalReference = "order"
	mu.Unlock()
	callback("mp-correct-order")
	if n, err := processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatal("MP reconciliation", n, err)
	}
	var hash, state string
	var hold bool
	if err = pool.QueryRow(ctx, `select p.state,o.evidence_sha256_hex,o.hold from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) where p.tenant_id=$1`, tenant).Scan(&state, &hash, &hold); err != nil || state != "captured" || hold {
		t.Fatal("MP observed", state, hold, err)
	}
	mu.Lock()
	postCount, getCount := posts, gets
	mu.Unlock()
	if postCount != 1 || getCount != 3 {
		t.Fatal("MP SDK effect count", postCount, getCount)
	}
	repo := db.NewFranchiseJourney(pool)
	sum := sha256.Sum256([]byte(franchisejourney.ReferenceHandoverContractDocument))
	contract := franchisejourney.HandoverReleaseContract{ID: "reference-single-unit-observed-payment-v1", DocumentSHA256: hex.EncodeToString(sum[:]), Scope: "LOCAL_FIXTURES", MaximumObservationAge: 5 * time.Minute}
	service, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, contract)
	if err != nil {
		t.Fatal(err)
	}
	prepared, _, err := service.Prepare(ctx, tenant, "operator", franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: payment.ID, ObservationSHA256: hash, IdempotencyKey: "mp-connected-handover-key"})
	if err != nil {
		t.Fatal("MP prepare", err)
	}
	ids := randomid.Generator{}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "mp-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	presented, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", prepared.Handover.ID, 1, "mp-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "mp-checklist", 1, strings.Repeat("e", 64), ids.New()); err != nil {
		t.Fatal(err)
	}
	release, err := service.EvaluateRelease(ctx, tenant, "store", prepared.Handover.ID, hash)
	if err != nil || !release.Eligible {
		t.Fatal("MP release check", err)
	}
	t.Logf("PAYMENT_CONNECTED_MP_PASS tenant=%s metadata_absent=true wrong_order_rejected=true official_get_money=true quote_to_release=true provider_posts=1", tenant)
}

func TestPaymentMaterializedHandoverProfile(t *testing.T) {
	r := newConnectedRun(t, false)
	ctx := context.Background()
	if code := r.callback(t, "evt_materialized", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	var hash string
	if err := r.pool.QueryRow(ctx, `select evidence_sha256_hex from payment.provider_observation where tenant_id=$1`, r.tenant).Scan(&hash); err != nil {
		t.Fatal(err)
	}
	python := os.Getenv("HANDOVER_PROFILE_PYTHON")
	if python == "" {
		t.Fatal("explicit admitted Python materializer runtime required")
	}
	load := func(mode string, refs ...string) franchisejourney.HandoverReleaseContract {
		t.Helper()
		account := "acct_fixture"
		if len(refs) > 0 {
			account = refs[0]
		}
		out := filepath.Join(t.TempDir(), "policy")
		cmd := exec.CommandContext(ctx, python, "-X", "utf8", "-B", filepath.Join("..", "..", "..", "tools", "materialize_handover_profile.py"), "--output", out, "--profile-id", "franchise-integral", "--revision", "1", "--tenant-id", r.tenant, "--organization-id", "store", "--expected-mode", mode, "--payment-provider", "stripe", "--payment-account-ref", account, "--payment-connection-id", "checkout", "--activate")
		if body, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("materializer %v %s", err, body)
		}
		raw, err := os.ReadFile(filepath.Join(out, "activation.json"))
		if err != nil {
			t.Fatal(err)
		}
		var activation franchisejourney.HandoverProfileActivation
		if err = json.Unmarshal(raw, &activation); err != nil {
			t.Fatal(err)
		}
		contract, err := franchisejourney.LoadHandoverProfileFile(filepath.Join(out, "profile.json"), activation)
		if err != nil {
			t.Fatal(err)
		}
		return contract
	}
	command := franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: r.payment.ID, ObservationSHA256: hash, IdempotencyKey: "materialized-profile-key"}
	repo := db.NewFranchiseJourney(r.pool)
	live := load("live")
	liveService, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, live)
	if err != nil {
		t.Fatal("supported future mode cannot load", err)
	}
	if _, _, err = liveService.Prepare(ctx, r.tenant, "operator", command); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("sandbox evidence admitted under live contract", err)
	}
	foreignAccount := load("sandbox", "acct_other")
	foreignService, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, foreignAccount)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = foreignService.Prepare(ctx, r.tenant, "operator", command); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("different provider account admitted", err)
	}
	policy := load("sandbox")
	service, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = r.pool.Exec(ctx, `update integration.provider_connection set state='disabled' where tenant_id=$1`, r.tenant); err != nil {
		t.Fatal(err)
	}
	if _, _, err = service.Prepare(ctx, r.tenant, "operator", command); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("disabled profile connection admitted", err)
	}
	if _, err = r.pool.Exec(ctx, `update integration.provider_connection set state='active' where tenant_id=$1`, r.tenant); err != nil {
		t.Fatal(err)
	}
	prepared, _, err := service.Prepare(ctx, r.tenant, "operator", command)
	if err != nil || prepared.ContractID != "franchise-integral@1" || prepared.ContractSHA256 != policy.DocumentSHA256 {
		t.Fatal("materialized profile preparation", prepared, err)
	}
	ids := randomid.Generator{}
	if _, err = repo.PublishDeliveryChecklist(ctx, r.tenant, "operator", franchisejourney.DeliveryChecklist{ID: "profile-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	presented, err := repo.CompleteDeliveryChecklist(ctx, r.tenant, "store", "operator", prepared.Handover.ID, 1, "profile-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AcceptHandover(ctx, r.tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "profile-checklist", 1, strings.Repeat("e", 64), ids.New()); err != nil {
		t.Fatal(err)
	}
	eligible, err := service.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash)
	if err != nil || !eligible.Eligible || eligible.Scope != "MATERIALIZED_PROFILE" {
		t.Fatal("materialized release", eligible, err)
	}
	if _, err = liveService.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("profile hash changed after preparation", err)
	}
	if _, err = service.Result(ctx, randomid.Generator{}.New(), "store", "order", command.IdempotencyKey); !errors.Is(err, franchisejourney.ErrReleaseConditioned) {
		t.Fatal("cross profile scope", err)
	}
	t.Logf("MATERIALIZED_HANDOVER_PROFILE_PG_PASS tenant=%s materializer_file_loader=true same_owners=true live_mode_mismatch_denied=true profile_hash_persisted=true live_proven=false", r.tenant)
}
````

### FILE: `internal/platform/postgres/payment_dispatch_fence.go`

```yaml
block_id: "GO-PAYMENT-CHECKOUT-RUNTIME:file26:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "401c49f162c825e0b7f203a485eaed32388ab8b8f7d75c4214dab7530a4c565a"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED composition: use the existing outbound fence transaction as the
// command acceptance point. No network call occurs while database locks are held.
import (
	"context"
	"encoding/json"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/paymentbridge"
	"github.com/jackc/pgx/v5"
)

type PaymentDispatchFence struct {
	*OutboundDeliveryStore
	Scope paymentbridge.Scope
}

func (f *PaymentDispatchFence) Claim(ctx context.Context, m channels.Message, hash string) (outbounddelivery.Claim, error) {
	if f == nil || f.OutboundDeliveryStore == nil || f.Scope.Validate() != nil || m.TenantID != f.Scope.TenantID || m.ChannelCode != "payment_"+f.Scope.ProviderCode {
		return outbounddelivery.Claim{}, paymentbridge.ErrCheckoutConflict
	}
	computed, err := outbounddelivery.MessageSHA256(m)
	if err != nil || computed != hash {
		return outbounddelivery.Claim{}, paymentbridge.ErrCheckoutConflict
	}
	var requested paymentbridge.Request
	if json.Unmarshal([]byte(m.Text), &requested) != nil || requested.PaymentAttemptID != m.DeliveryKey || requested.OrderID != m.ThreadID || requested.CustomerSubject != m.ExternalID {
		return outbounddelivery.Claim{}, paymentbridge.ErrCheckoutConflict
	}
	return f.claimWithAdmission(ctx, m, hash, func(ctx context.Context, tx pgx.Tx) error {
		if err := checkoutConnection(ctx, tx, f.Scope); err != nil {
			return err
		}
		actual, err := lockCheckoutRequest(ctx, tx, f.Scope, requested.PaymentAttemptID)
		if err != nil {
			return err
		}
		if actual != requested {
			return paymentbridge.ErrCheckoutConflict
		}
		// Bind admission and the created→pending command acceptance in the same
		// transaction as the one allowed send. Later cancellation never creates a
		// second send; subsequent financial observations preserve terminal states.
		row, err := tx.Exec(ctx, `with accepted as (
   update payment.payment_attempt p set state='pending',version=p.version+1,updated_at=clock_timestamp()
   from sales.customer_order o,payment.provider_checkout c
   where p.tenant_id=$1 and p.payment_attempt_id=$2 and p.state='created'
     and o.tenant_id=p.tenant_id and o.order_id=p.order_id and o.organization_id=$3 and o.state in ('placed','confirmed','allocated')
     and c.tenant_id=p.tenant_id and c.payment_attempt_id=p.payment_attempt_id and c.connection_id=$4 and c.account_ref=$5 and c.request_sha256_hex=$6 and c.live_mode=$7 and c.session_id is null and c.expires_at>clock_timestamp()+interval '31 minutes'
   returning p.version)
   insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
   select $1,gen_random_uuid(),'payment',$2,version,'payment.checkout-requested',1,clock_timestamp(),jsonb_build_object('connection_id',$4::text,'request_sha256',$6::text) from accepted`, f.Scope.TenantID, requested.PaymentAttemptID, f.Scope.OrganizationID, f.Scope.ConnectionID, f.Scope.AccountRef, hash, f.Scope.LiveMode)
		if err != nil {
			return err
		}
		if row.RowsAffected() != 1 {
			return paymentbridge.ErrCheckoutConflict
		}
		return nil
	})
}
````

## 6. Configuration surface

Exact future environment and profile fields are listed in the materialized docs. Defaults disable PAYMENT_CHECKOUT_ENABLED and HANDOVER_ENABLED. Enabling requires complete scope and validation before traffic. Provider secrets, webhook secret and outbound HMAC key are external runtime inputs. Hosted return URLs must be explicit trusted HTTPS configuration. Profile ID/revision/SHA and supported algorithm/options bind the same deployment scope. No secret value is present in the pack.

## 7. Dependency bill

| Component | Fixed identity | Use / license |
|---|---|---|
| Go | 1.26.8 | standard library / BSD-3-Clause |
| PostgreSQL | 18.6 | existing transaction owner / PostgreSQL |
| pgx/v5 | exact existing root go.mod/go.sum | driver / MIT |
| official payment adapter | GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS0.3.2 | local SDK composition / workspace terms |
| Stripe-Go | 86.3.0, a2df585a800a97fe8ec4ebf551b4449bdb3d90a1 | external SDK dependency pin / MIT |
| MercadoPago Go SDK | 1.14.0, f910ee53fbb6819e435eaf3d0f800cb1fe74ae09 | external SDK dependency pin / MIT |
| Python | CPython3.14.4 actual reference | profile materializer only / PSF |

No vendor source is relabeled as authored product code. Exact source/license/SCA and SDK reconstruction receipts remain in PAYMENT_SDK_RECONCILIATION_V402.md.

## 8. Apply order

Compose the whole selected profile into an absent external destination. Apply migrations in numeric order;0055 precedes0056 and0057. Keep the SDK module replace path, its exact module files and the existing base owner modules. For an existing target, compare baseline and preserve durable records before coordinated migrations. The down scripts are for disposable reference/reversible pre-use schema checks, not authorization to erase production payment or handover evidence.

## 9. Verification

Run the affected Go unit/HTTP/host tests, bounded native fuzz where the selected input surface supports it, and the explicit loopback PostgreSQL connected-payment/profile/read tests with skips rejected. Run vet and build on the composed revision. Frontend selection additionally uses pinned Node/Next/Vitest/TypeScript from the admitted recipe, no unreviewed package-manager invocation. Verify wrong scope, provider mode, account, amount, URL, hash, generation, claim, duplicate callback, response loss and immutable terminal-state behavior. Compare a clean canonical reconstruction byte-for-byte before promotion. Tests of fixture contracts do not authorize live accounts or production.

## 10. Reconstruction evidence

V402: the complete77pack/915file reference was reconstructed into an absent external destination; every output matched the frozen source inventory. PostgreSQL57migrations and16 connected tests passed without skips, plus vet/build and the explicitly enumerated unit/host cases. Finite fuzz and frontend contract receipts are separately bound to identical files. See reconstruction_evidence/CONNECTED_PAYMENT_HANDOVER_V402.md and CONNECTED_DELTA_ASSURANCE_V402.md. Current SAST/SCA, whole-library operations/performance/portable release and remaining capabilities retain their own gates; no live production or complete TEST02 claim.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.
