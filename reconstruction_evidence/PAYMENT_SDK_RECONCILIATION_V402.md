# V402 — official payment SDKs, hosted checkout and reconciliation boundary

2026-09-11. Local library maintenance. GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS0.3.0:12files,17tests PASS, exact Markdown rebuild12/12, go mod verify/vet/build PASS, OSV2.5.1 runtime graph3modules/0findings. These results cover this standalone adapter; franchise consumer/inbox, frontend, handover and release retain their own evidence.

## Concrete change

The existing pinned SDK module now supports Stripe hosted Checkout Session and MercadoPago Checkout Pro Preference Create/Get, canonical notification hints and official payment GET observations. Stripe GET expands latest_charge to preserve captured/refunded/disputed values. Account probe uses GET/v1/account with the exact credential. Checkout IDs, order/attempt metadata, amount/currency, HTTPS provider redirect hosts and bounded expiry are checked. No payer/card input is collected by hosted checkout. A browser return or created intent does not establish collection.

Notifications are authenticated by official SDK verification. Stripe session, payment-intent, refund and dispute events route to official GET; signed resource links never substitute the GET. MercadoPago signs resource/request/timestamp, so unsigned body status/amount/user_id is not used as financial authority. New normalization emits mercadopago to match Commerce while the legacy generic verifier API retains mercado_pago for compatibility. Snapshot hashing binds a local normalized observation, not a provider signature.

Production HTTP clients have a10second timeout, reject redirects and disable automatic network retries. WithHTTPClient is a trusted composition seam used for offline test transports. Durable claim/fence, ambiguous-effect handling, account/tenant/order/amount binding and local ledger transition belong to the existing application owners.

## G0–G8 and provenance

| Gate | Evidence / decision |
|---|---|
| G0 | Stripe86.3.0 commit a2df585a800a97fe8ec4ebf551b4449bdb3d90a1 and MercadoPago1.14.0 commit f910ee53fbb6819e435eaf3d0f800cb1fe74ae09 unchanged.13immutable official source/license files equal the module-cache bytes; Go module verification PASS. |
| G1 | Existing MIT SDK pins and license bytes preserved. No new dependency. New source is AUTHORED boundary glue under workspace owner license, not Stripe/MercadoPago-authored implementation. |
| G2 | Backend/API, distributed idempotency, explicit scope and source-admission authorities; no financial policy or pricing calculation added. Decimal SDK conversion rejects fractional minor units instead of rounding. |
| G3 | Create hosted checkout → durable owner receipt → authenticated notification hint → official GET snapshot → existing domain owner. Only SDK/client normalization belongs to this pack. |
| G4 |17tests including official serializers, hosted sessions, GET/charge/refund/dispute, account probe, callback signature/account/mode, unsigned MP claims and invalid URL/amount/identifier.12files rebuilt byte-identical. |
| G5 | No raw payment/card/payer data in normalized snapshot or notice; no client-secret in snapshot; bounded body/IDs; fixed SDK endpoints and trusted transport injection; production callback/storage gates remain in composition. |
| G6 |10s HTTP timeout, zero implicit retries, explicit context cancellation through SDKs; finite30min–24h hosted-session request expiry. Local transports only, no throughput/production claim. |
| G7 | Original0.2.0 pack and7materialized files preserved outside canonical tree. Rollback disconnects this adapter and preserves application receipts; it never reverses money or deletes inbox. |
| G8 | Standalone adapter delta admitted for local fixture composition under its documented conditions. Runtime graph3modules/0OSV findings. Main composition must demonstrate its own durable consumer/callback/reconciliation gates before claiming complete infrastructure. |

## Official authorities observed

The [Stripe SDK fixed service](https://github.com/stripe/stripe-go/blob/a2df585a800a97fe8ec4ebf551b4449bdb3d90a1/checkout_session_service.go) provides hosted session Create/Retrieve; the [current API contract](https://docs.stripe.com/api/checkout/sessions/create) defines session return URLs and expiry. The [MercadoPago pinned Preference client](https://github.com/mercadopago/sdk-go/blob/f910ee53fbb6819e435eaf3d0f800cb1fe74ae09/pkg/preference/client.go) supplies Create/Get. The selected MP redirect contract is the [Argentina Checkout Pro lane](https://www.mercadopago.com.ar/developers/en/docs/checkout-pro-preferences/create-payment-preference); no other country/custom-domain compatibility is inferred. Its CheckoutResult.LiveMode describes redirect selection; actual payment mode is observed only from PaymentSnapshot.LiveMode.

## Reproducible local evidence

Stage: C:/Users/NL/AppData/Local/Temp/elite-v402-library-infra.
Final source: payment-sdk-release-rebuilt/official_payment_webhooks.
Environment: Go1.26.8 exact previously admitted binary, GOPROXY=off/GOSUMDB=off/GOTOOLCHAIN=local. SDKs and runtime deps were unchanged. Tests invoke only fixture transports/loopback; no live accounts, user credentials or provider writes.

Earlier intermediate10file/11test results are superseded by the12file/17test final delta, not combined as additional progress. Companion webhook glue has two local tests plus five negative subcases and vet PASS in payment-bridge-test; root integrates its existing PostgreSQL inbox/hold hook separately.

Failures retained: one official raw-source fetch returned503 Backend.max_conn; a bounded retry remained503 for an additional file. Official GitHub Contents API at the same commit succeeded and its Git blob SHA/cache bytes matched. An inline PowerShell script repair had a parser error with no mutation; apply_patch replaced it. A standalone go mod tidy requested upstream test-only testify absent offline; fixed fixture go.mod lists the exact3runtime pins and uses -mod=readonly without fetching or upgrading dependencies. No product/SDK test failed and no check was disabled.

## Exact receipts

| Receipt | SHA-256 |
|---|---|
| `payment-source-receipts-complete.json` | `ea8a210a7d96ff61e9b36ade2764df791f9acceb6e2d04bcf2e348651f8505c7` |
| `payment-sdk-final-result.json` | `48fe4802c7eeb9522284ae46890d8ecd023429e7908b3cff2ba73453afc29773` |
| `payment-sdk-final-mod-verify.log` | `b4537ed75f533f993f371954de47e42a793b8e5b0587577de7e27fb3e50696bd` |
| `payment-sdk-final-tests.log` | `9b7aef52f9f5790850767a4ae548d2473d15254821ba2ed5c081921bc3bc476c` |
| `payment-sdk-final-vet.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `payment-sdk-final-build.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `payment-sdk-osv.json` | `01b8c26c57808538d47521d42f2ff3d93f78cae489d480f97c08a60c98b9d200` |
| `payment-bridge-test-result.json` | `a355bb0d045d9ee467280a0ab2ad892ec19211ebfe618a12022106bebe2c407d` |
| `payment-source-failures.log` | `1fefd67f38bcf575aa1b3030ddfdcfed2665aa582578f6f2806e7924f17ae30e` |
| `payment-bridge-tidy.log` | `5aba8790a1462b1a1bcb95f456cdb1fe374d60fff2261c27e8e5a45ea8894312` |

| Materialized source | SHA-256 |
|---|---|
| `go.mod` | `fc6106ab8e8e4ab2dde71a6210e5d4b31cc674e62088947e29eb8e025387b141` |
| `go.sum` | `6f564169198e386498d553a0c9a1e7bce16b619d4e88d640429ad1f3ff74dc06` |
| `README.md` | `63d63590cf39c65efa80488bf8aaf2ffb71ad2ba46fa47d1ebaf8e12a846bbaf` |
| `officialpayments/checkout.go` | `1533b728ca09a06016aa8e0c7b2828cb9ea630b9b3f1bab6834f2d70999b5338` |
| `officialpayments/checkout_test.go` | `754bef6d11f6d66574cf5555b95c111ed35cc70d69cbe44242797dfb7f21efa3` |
| `officialpayments/client.go` | `596c770ea1d6395f33cfaca3892c427963b46f6308f9cbbe1235907191f8bf6f` |
| `officialpayments/client_test.go` | `712ba65d51affd0dc5e661d2e31990ace3a744234ca0c257bcd4f4e9708989ef` |
| `officialpayments/notification.go` | `11ad6b9441b06c474ff2e5e7862c6fd118347b24d5993db8311b2ec631752da9` |
| `officialpayments/retrieve.go` | `5148b8b7d9195cb300bbdba3b4b6f16f7b161958770b2a395d214a12cae22b26` |
| `officialpayments/retrieve_test.go` | `0a6b926e852f8a5480f953dbf8c79df0622c70ab2b20f6405ed5f341c3da6913` |
| `officialpayments/verifier.go` | `64ba006387b43979648306f4b0acb0e4de3db7ea152f33f69ea0becd6ed74eac` |
| `officialpayments/verifier_test.go` | `e924989d2efbe1a3c2be3ecbbbc72a5bf4367e5e5b63e4a73902a6afea1b1939` |


## Connected SDK driver delta (staging, exact application interface)

The `payment-bridge-files/internal/paymentbridge/sdk_driver.go` composition implements the root application's existing `Driver` interface and its persisted `Request.CheckoutExpiresAt`. Request amounts, order, attempt, tenant, organization, currency and provider come from the bound owner request. Expiry and idempotency key remain stable across calls. Constructor configuration supplies return URLs and display name; no browser input supplies provider references or financial authority.

`ValidateCredential` observes Stripe GET `/v1/account`, or MercadoPago official SDK GET `/users/me` with the exact account plus Argentina `AR/MLA` country lane, before any checkout POST. Stripe's documented `sk_test_`/`rk_test_` and `sk_live_`/`rk_live_` prefixes reject an incompatible configured mode before HTTP; the prefix is not treated as account authentication. Hosted checkout metadata, amount, currency, mode and expiry are checked against the immutable request. MercadoPago preference redirect selection is not proof of the payment's live mode: actual payment GET must also satisfy `PaymentSnapshot.LiveMode`, account and amount binding in the storage owner.

Five top-level tests and five negative webhook subcases PASS through the pinned SDK HTTP serializers and a loopback HTTP server. They cover hosted Stripe checkout to signed session callback to session GET to expanded payment GET; refunded and disputed observations; MercadoPago account probe and hosted preference; wrong account rejected before another POST; Stripe key mode mismatch before HTTP; and callback signature/scope/commit errors. `go vet` PASS. The test module copies the actual root `checkout.go`; its tested hash is in the receipt. All SDK network requests are restricted by the fixture transport to the explicit provider hostname, then rewritten to a loopback server. No live endpoint or real credential was used. This evidence does not replace the separate PostgreSQL consumer, observation-generation fence or handover integration proof.

The two new driver files are AUTHORED composition/contract-test glue. No new dependency or copied vendor implementation was added; all provider serialization, user retrieval and payment operations use the same already pinned official SDK modules. Two additional MercadoPago user source files at commit `f910ee53fbb6819e435eaf3d0f800cb1fe74ae09` were obtained through the official GitHub Contents API, Git blob verified and found byte-identical to the module cache. Existing SDK licenses and redistribution obligations remain unchanged. Stripe key authority pages were captured with SHA-256; these documentation observations do not change an SDK pin.

| Staging receipt / source | SHA-256 |
|---|---|
| `payment-sdk-driver-result.json` | `4035ed65c1a7dc152a3db42ec52cf4a39214749cc3881f3d21a1dcb125c5c3d8` |
| `payment-sdk-driver-tests.log` | `ffe5f3e1715af7cb51416b581ed53ab8c15fb72baafe7fd21b57e681c3c0aaf5` |
| `payment-sdk-driver-vet.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `payment-driver-official-sources.json` | `23e01b8dd027275834873a86d6972091bda202376a9c431a8fcfed369413720a` |
| `payment-driver-key-authorities.json` | `fb85056e44707df861a2f9e334a31484ca94c7870f5adcb87b57e56cdd1c02af` |
| `payment-bridge-files/internal/paymentbridge/sdk_driver.go` | `fc26f4bdae3576b5cafd7e90f6d101a8ae4ae7582a696110a7bfe79993e02b17` |
| `payment-bridge-files/internal/paymentbridge/sdk_driver_test.go` | `39aee9de29bcb33db266a00ab5502237867ff53ed4333e0c0cd7f553ab451144` |

## Latest SDK patch — 0.3.1

The connected PostgreSQL/official SDK E2E observed the default Stripe logger writing a malformed response body sample. The per-backend official LevelNull option now suppresses raw SDK logging without changing the global logger; application errors remain bounded. A subprocess regression captures stdout/stderr and rejects raw-body disclosure. Twelve materialized files rebuilt exactly,18 SDK tests and vet/build PASS. Fixed Stripe `stripe.go` Git blob/cache bytes match the same admitted commit. No dependency changed; the earlier3-module/0-finding SCA result remains the applicable unchanged graph.

Receipt: `C:\Users\NL\AppData\Local\Temp\elite-v402-library-infra\payment-logger-result.json`, SHA-256 `a948997537d479a60368d270eaee7c95e91c8ccea9ecd019fca3600fb6715360`. Earlier0.3.0 results remain historical;0.3.1 is the current pack version.
