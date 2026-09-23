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
