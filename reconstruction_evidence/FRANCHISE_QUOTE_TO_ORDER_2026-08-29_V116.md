# Franchise Quote-to-Order — Reconstruction Evidence V116

## Scope

This milestone extends the existing franchise journey and existing canonical order owner; it does not create a second commerce model. The implementation is `AUTHORED`, governed by the pinned Microsoft BCApps and Google Online Boutique references already declared by the packs, plus NASA systems verification and Google SRE launch/release discipline. No local line is attributed to those companies.

## Materialized change

- `GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.2.0`: 12 files, including migration 0006, SQL invariant test, authenticated quote acceptance and atomic PostgreSQL conversion.
- `GO-ELECTROMOBILITY-APPLICATION 0.7.1`: unchanged composition code, compatibility contract advanced to journey 0.2.x.
- `TS-FRANCHISE-JOURNEY-PORTALS 0.3.0`: 10 files, including same-origin BFF mapping and `/customer/quotes` action UI.
- Exact backend profile: 17 packs / 138 files.
- Exact web profile: 6 packs / 80 files.

The browser cannot submit price, currency or an evidence digest. The Go HTTP boundary hashes the exact strict JSON command after authentication; PostgreSQL locks the customer-owned, organization-scoped, issued and unexpired quotation and derives currency and total from its server record.

## Exact executable evidence

Official Go 1.26.7 archive SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11` and PostgreSQL 18.6 EDB archive SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c` were reused from the verified local toolchain.

From new Markdown-only compositions:

- backend materialization: 17 packs / 138 files;
- full `go test ./...`, `go vet ./...` and both command builds: PASS;
- empty PostgreSQL database, migrations 0001–0006 and every SQL test with `ON_ERROR_STOP=1`: PASS;
- `TestFranchiseJourneyPersistenceIsolationAndReplay` against PostgreSQL 18.6: PASS;
- successful acceptance created one `placed` order, one line at the server price, one immutable acceptance row and two outbox events;
- stale version, cross-customer and expired quotation: rejected;
- forced duplicate event on the second outbox write: entire transaction rolled back, leaving quote issued and zero order/acceptance/accepted event;
- 0006 down/up/test and full 0006→0005 down then 0005→0006 up/tests: PASS;
- web materialization: 6 packs / 80 files; frozen offline install downloaded zero packages;
- nine Vitest files / 30 tests, strict TypeScript and Next.js 16.3.2 production build: PASS;
- build route `/customer/quotes`: emitted;
- Microsoft Playwright 1.62.1: runtime 4 and target 8 tests PASS across Chromium desktop/mobile, Firefox and WebKit;
- Google Lighthouse 13.4.1: five governed runs PASS; performance median 0.99, accessibility/best-practices/SEO medians 1.00.

## Failures converted to memory

`LIB-FAIL-1253` through `LIB-FAIL-1257` record, respectively: isolated multi-package Go execution, wrong active module, stale journey cardinality, incorrect `$LASTEXITCODE` handling for a PowerShell script and an empty-output compound final command. Each was reproduced, corrected and re-proven by isolated exact gates.

## Honest boundary

This proves local reconstruction and the quote-to-order transaction. It does not prove legal acceptance wording, tax/fiscal rules, payment capture, financing, target IdP, target CDN/WAF, production recovery, target load, offensive security or business acceptance. Those remain project-specific production gates and cannot inherit this PASS.
