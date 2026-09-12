# Franchise customer journey V114

## Scope

This evidence governs the 2026-08-29 V114 franchise/customer-journey increment. It supersedes V113 only for current library totals and the backend/web profiles; earlier evidence remains historical for unchanged packs.

The increment adds two packs and extends three existing composition owners without creating a second CRM, catalog, pricing, order, inventory, identity or HTTP-client implementation:

- `GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.1.1`: 9 `AUTHORED` files;
- `TS-FRANCHISE-JOURNEY-PORTALS 0.1.0`: 5 `AUTHORED` files;
- `GO-ELECTROMOBILITY-APPLICATION 0.7.0`: wires the new Go module;
- `TS-GO-API-WEB-BRIDGE 0.4.1`: aligns public/protected clients, configuration, deterministic metadata and local icon;
- `TS-OIDC-PORTAL-ADAPTER 0.2.0`: aligns journey types and allowed return paths.

Microsoft BCApps and Google Online Boutique are architecture/contract references fixed in pack metadata. NASA systems verification/integration and Google SRE launch/release guidance govern the evidence model. The fourteen journey files plus the icon are local `AUTHORED` implementation; no line is falsely attributed to Microsoft, Google, NASA or another vendor.

## Current inventory

- 75 implementation packs;
- 756 materialization blocks: 619 `AUTHORED`, 32 `ADAPTED`, 105 `VERBATIM`;
- 36 composition profiles;
- backend profile: 17 packs / 135 files;
- web profile: 6 packs / 75 files;
- 121 exact upstream sources / 16 source profiles remain unchanged;
- failure memory: 1,252 local IDs plus 202 upstream conditions, with no open local row.

`VERIFY_LIBRARY.ps1` passed before this evidence with 444 Markdown files and, after adding it, must report 445. Every selected pack version, exact manifest, materialization block, path and SHA-256 matched; all 36 profiles composed.

## Backend proof from Markdown

The compositor rebuilt 17 packs / 135 files into a new empty temporary directory. Official Go 1.26.7 for Windows/amd64 was invoked by absolute path from a verified ZIP whose SHA-256 is `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`.

The exact composed tree passed:

- `go test ./...`;
- `go vet ./...`;
- `go build ./cmd/...`;
- migrations 0001–0005 on a new PostgreSQL 18.6 database;
- every SQL invariant under `db/tests`;
- every Go test again with `TEST_DATABASE_URL` bound to that database;
- migration 0005 down, reapply and its SQL invariant again.

The PostgreSQL 18.6 Windows binary archive used only as a temporary test runtime was obtained from the official PostgreSQL Windows download route to EDB and measured as 343,808,005 bytes with SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`. It is not redistributed by the library.

The runtime proof exercised public location publication, request-hash appointment idempotency/replay, customer ownership, organization isolation, optimistic lead assignment/transitions, database transition enforcement, price-book-authoritative quotations, zero orphan outbox event when the scoped lead is absent, customer journey reads, subject-bound delivery acceptance, rollback and reapplication.

## Web proof from Markdown

The compositor rebuilt 6 packs / 75 files into a separate empty temporary directory. With Node 24.14.1 and pnpm 11.19.0, a frozen offline install downloaded zero packages and the exact tree passed:

- 8 Vitest files / 25 tests;
- Next type generation and TypeScript 7 strict typecheck;
- Next.js 16.3.2 production build;
- Microsoft Playwright 1.62.1 runtime smoke in Chromium desktop/mobile, Firefox desktop and WebKit desktop: 4/4;
- the exact production build in those four projects for semantic home, responsive overflow, security headers and fresh nonce/CSP: 8/8;
- Google Lighthouse 13.4.1 contract tests: 5/5;
- five target audits: performance median 0.99, accessibility 1.00, best practices 1.00 and SEO 1.00.

The first Lighthouse series was rejected because one run observed no meta description and a missing favicon returned 404, producing SEO minimum 0.90. The owner pack was corrected with deterministic root metadata and a local SVG icon, then the complete build/browser/Lighthouse sequence was repeated. This failure and every other V114 harness or implementation failure is retained in `markdown_system/LIBRARY_FAILURE_LEARNING_LEDGER.md` through `LIB-FAIL-1252`.

## What this proves

This proves a reconstructible, executable vertical slice from visitor contact into franchise operations and customer visibility. It proves the local contracts and negative cases listed above on the exact lab artefacts. It materially reduces implementation time because a new project can compose these files instead of regenerating the same foundation.

It does not prove a complete franchise business or production readiness. Current open product capabilities include authenticated mutation UX for franchise operators, organization/branch administration, scheduling capacity, quotation acceptance and order conversion, royalty/settlement, accounting/tax/fiscal adapters, financing, consent/attribution/analytics, full content/media/search, target IdP/roles, real provider accounts, CDN/WAF, workload and soak tests, offensive security, recovery/PITR, cloud deployment/canary/rollback and signed business acceptance. Those are work to implement or target evidence to obtain; none is converted into a positive claim by this document.

## Governing public references

- NASA Systems Engineering Handbook appendix: `https://www.nasa.gov/reference/system-engineering-handbook-appendix/`;
- Google SRE, Reliable Product Launches: `https://sre.google/sre-book/reliable-product-launches/`;
- Google SRE, Release Engineering: `https://sre.google/sre-book/release-engineering/`;
- Google Cloud Online Boutique exact reference: `https://github.com/GoogleCloudPlatform/microservices-demo/tree/5b3a712ab85ccb8f6f7cd5b720d36ba9a8d041eb`;
- Microsoft BCApps exact reference: `https://github.com/microsoft/BCApps/tree/31a860b527f0dc72c7a44a255d7e7d403cfa4789`;
- official Go downloads: `https://go.dev/dl/`;
- official PostgreSQL Windows downloads: `https://www.postgresql.org/download/windows/`.

## Result

`V114_REBUILD_VERIFIED / PROJECT_CONDITIONED`.

No project may claim production readiness until its exact CDN/WAF, IdP, providers, PostgreSQL recovery, load/resilience, offensive security, deployment/rollback and business acceptance evidence all pass together on the packaged release artefact.
