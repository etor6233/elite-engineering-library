# Selected Provider Production Admission — V111

Date: 2026-08-29

## Scope and provenance boundary

This checkpoint adds a local fail-closed semantic adapter for the existing `PROVIDERS` production control. It adds no vendor source code and does not claim that Stripe, Mercado Pago, Amazon, Google, Meta, TikTok, Firebase or Mercado Libre authored the adapter. It consumes evidence from exact official SDKs/samples already fixed in the library's source lock and implementation packs.

Mercado Libre is deliberately different: the active adapter is Elite `AUTHORED` against current official HTTP contracts because the archived Mercado Libre SDK is not represented as maintained reusable official code. The validator records that source identity explicitly and never relabels it as vendor code.

No account, credential, sandbox, webhook, quota, cost, external write or provider reconciliation was exercised in this checkpoint. The 14 regressions use synthetic local evidence solely to prove that the contract fails closed.

## Official rules applied

- Stripe recommends official-library webhook signature verification over the raw payload, warns that events can duplicate and arrive out of order, and requires integration-side deduplication/reconciliation: <https://docs.stripe.com/webhooks>.
- Stripe documents idempotency keys for safe POST retries, indeterminate network/500 outcomes and exponential backoff; its rate-limit documentation distinguishes 429 rate/concurrency reasons: <https://docs.stripe.com/error-low-level> and <https://docs.stripe.com/rate-limits>.
- Mercado Pago requires `X-Idempotency-Key` for applicable payments/refunds and documents secret-signature Webhooks plus test URLs/simulation: <https://www.mercadopago.com.ar/developers/es/news/2023/01/04/Idempotency-key-usage-will-be-mandatory> and <https://www.mercadopago.com.ar/developers/es/docs/checkout-bricks/additional-content/your-integrations/notifications/webhooks>.
- AWS documents bounded SDK retry modes, throttling classification, retry quotas and exponential backoff with jitter: <https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html>.
- Google documents quotas and provider-directed throttling handling, including exponential backoff, while Merchant quotas remain operation/account specific: <https://developers.google.com/merchant/api/guides/quotas-limits/quotas> and <https://developers.google.com/workspace/drive/api/guides/limits>.

These authorities define required behavior; they do not supply project accounts or prove the project's integration.

## Exact selectable identities

The registry accepts only these declared identities and versions:

1. Stripe Go SDK 86.3.0, exact commit `a2df585a800a97fe8ec4ebf551b4449bdb3d90a1`;
2. Mercado Pago Go SDK 1.14.0, exact commit `f910ee53fbb6819e435eaf3d0f800cb1fe74ae09`;
3. Amazon Selling Partner API Python SDK 1.11.1, exact commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`;
4. Google Ads Python SDK 31.3.0, exact commit `f7bf312d26904ea50ad9ef02e395edccfaa1ebb7`;
5. Meta Business Python SDK 26.0.1, exact commit `788f363d15b1269ab5efb7cd00fb5e3b133cd99b`;
6. TikTok Business API Python SDK, exact signed commit `f809c396520df2d7b201a9ccc5378d822b728ed3`;
7. Meta WhatsApp Cloud API examples, exact signed commit `de70ee908a67026e642aaee3703d20464e2a9466`;
8. Firebase Admin Go 4.21.0, exact commit `eebb06f2a643fbb59b1cb262874a943584475128`;
9. Google Merchant Products Python 1.8.0, exact source commit `97d7b42cd74b41211f5ec8871cc0dd15debdb1a0`;
10. Elite Mercado Libre official-HTTP-contract adapter 0.1.0, explicitly `AUTHORED`, sourced from current official Mercado Libre documentation rather than an active vendor SDK.

The first nine identities are already governed by the acquisition lock and their focal packs. The tenth is governed by `GO_MERCADOLIBRE_MARKETPLACE_ADAPTER.md`, including its rejection of the archived SDK as reusable implementation.

## Executable semantic adapter

`SECURE-OPS-DELIVERY-CORE` 0.9.0 adds three `AUTHORED` files:

1. `provider-admission.template.json`;
2. `validate_provider_admission.py`;
3. `test_validate_provider_admission.py`.

The shipped template selects zero providers and cannot pass. A project must select only required operations and create an approved policy for every provider. That policy binds project/environment/release, adapter and source identity, exact operations/scopes, account and sandbox references, webhook/mutation mode, terms/version/review date, data regions, retention, quota/cost and three independent finance/security/reconciliation owners. Terms review older than 90 days at approval is rejected.

Every selected provider requires four ordered official-tool executions through the shared no-shell, hash-bound runner:

1. `SANDBOX_CONTRACT`: every approved operation needs a provider response, schema hash and proof of no unapproved production effect;
2. `WEBHOOK_AUTH`: positive signature plus altered/stale/replayed negatives, duplicate convergence, durable inbox and missed-event reconciliation; a provider without callbacks needs an exact reason and two reviewers;
3. `IDEMPOTENCY_RECONCILIATION`: key/dedup proof, one remote effect under replay, uncertain-outcome reconciliation, provider/local state agreement and recovery or compensation;
4. `RATE_LIMIT_RETRY`: observed throttling, provider-directed or jittered backoff, 2–10 bounded attempts, no retry of permanent errors, durable exhaustion and absence of retry storm.

All four executions for one provider must use the same exact tool/source hash. Arguments, target, observation, stdout and stderr are hash-bound. Secret-bearing environment-variable names are rejected unless they are reference variables ending in `_REF`.

The resulting receipt emits exactly the five assertions expected by the final production gate: `all_required_sandboxes_pass`, `webhook_auth_pass`, `idempotency_reconciliation_pass`, `rate_limit_retry_pass` and `terms_cost_scope_approved`.

## Reconstruction and gates

- secure operations 0.9.0 reconstructs 31/31 files;
- source-tree versus Markdown round trip: 31/31 paths and SHA-256 values match;
- provider adapter: 14/14 regressions PASS;
- complete production-admission group: 73 tests, 72 PASS and one symlink-only platform skip;
- `VERIFY_LIBRARY_PASS` is expected to govern 73 packs, 734 materializable files, 440 Markdown including this evidence, backend 16/120 and web 5/69;
- the source lock remains 120 exact sources and 15 source profiles; V111 introduces no new vendor archive;
- no final portable ZIP is created at this checkpoint.

Canonical checkpoint hashes before adding this evidence:

- `SECURE_OPERATIONS_DELIVERY_CORE.md`: `91b08f3c2152875a0836e7244e6cb81eb3c75c73389e05951cc546382ef8e025`;
- `OFFICIAL_UPSTREAM_ACQUISITION_CORE.md`: `0c6c3cfe2cc85a8f502820d77ddd9c3c2f5662307b84ab02fa43c873b2b77780`;
- `ENTERPRISE_BACKEND_PACK_PLAN.md`: `66d7f63a72e13b608fec9687386d7129c62f9dc645401a005cb441573936c54c`.

## Failure memory and remaining work

`LIB-FAIL-1220` preserves the rejected materialization caused by missing mandatory block metadata, and `LIB-FAIL-1221` preserves the atomically rejected correction that used stale hashes. `UP-FAIL-200` preserves the essential upstream limit: official provider code cannot prove a project's accounts, contracts, effects or reconciliation.

The library now contains 1,221 local failure identities plus 200 upstream conditions: 1,421 unique IDs. `PROVIDERS` is an executable semantic adapter, but no concrete project is `PRODUCTION_ADMITTED`. Two semantic adapters remain absent: `DEPLOY_ROLLBACK` and `BUSINESS_ACCEPTANCE`; all eight real receipts must still pass together for one immutable release digest.
