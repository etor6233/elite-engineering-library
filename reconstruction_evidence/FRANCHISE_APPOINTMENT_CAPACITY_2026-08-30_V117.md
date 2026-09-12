# Franchise Appointment Capacity — Reconstruction Evidence V117

## Scope and provenance

This milestone extends the single existing franchise journey, appointment owner and web BFF. It creates no second CRM, calendar, organization or HTTP client. Every changed application line is `AUTHORED` and remains licensed as workspace-owner code; no local Go, SQL or TypeScript line is represented as Microsoft source.

The architecture/test authority is the exact public Microsoft `BCApps` commit `31a860b527f0dc72c7a44a255d7e7d403cfa4789`. The acquired archive receipt fixes that commit and archive URL. The two directly reviewed files were:

- `src/Layers/W1/DemoTool/CreateResourceCapacityEntry.Codeunit.al`, SHA-256 `9a5bb85eb93bbd82d9ef7b043f9aa060653adee5148ed4fd4555a0e8b4048e42`: dated resource-capacity entries;
- `src/Layers/W1/Tests/Resource/ResourceMatrixManagement.Codeunit.al`, SHA-256 `fa9e96b6a884c6fd4396f93cbf642ded6cb18272c0300dce22c2b6e48dcd754a`: capacity, zero-capacity, positive/zero/negative availability and period tests.

Google SRE launch/release discipline governs proportional gates and exact packaged-artifact testing; NASA verification/validation discipline governs separation of local component proof from project production acceptance. Those sources govern method and claims, not authorship.

## Materialized change

- `GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.3.0`: 15 files. Migration 0007 adds organization/kind slots, bounded capacity, server-derived end time, indexes, binding constraints and a serialized PostgreSQL non-overlap guard. Public reads expose only future, open and non-full published slots. Booking locks the exact slot and performs capacity check, appointment insert, idempotency and outbox in one transaction.
- `GO-ELECTROMOBILITY-APPLICATION 0.8.0`: unchanged composition code with compatibility advanced to journey 0.3.x.
- `TS-GO-API-WEB-BRIDGE 0.5.0`: its existing public client now reads bounded server-owned slots.
- `TS-FRANCHISE-JOURNEY-PORTALS 0.4.0`: still 10 files. Visitors choose an exact server slot instead of inventing a datetime; authorized franchise operators publish bounded capacity through the existing same-origin command BFF.
- Exact profiles: backend 17 packs / 141 files; web 6 packs / 80 files.

## Exact executable evidence

From new Markdown-only compositions:

- official Go 1.26.7 archive SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`: `gofmt -d`, full `go test -count=1 ./...`, `go vet ./...` and builds for the two materialized commands PASS;
- official PostgreSQL 18.6 EDB archive SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`: empty database, migrations 0001–0007 and every SQL invariant test with `ON_ERROR_STOP=1` PASS;
- live integration: overlapping open interval rejected, full slot removed from public results, over-capacity request rolled back without idempotency residue, and two simultaneous requests at capacity one produced exactly one persisted booking plus one conflict;
- 0007 down/up/test and 0007→0006 down followed by 0006→0007 up/tests PASS;
- exact web frozen offline install downloaded zero packages; nine Vitest files / 32 tests, strict TypeScript and Next.js 16.3.2 production build with 17 routes/pages PASS;
- Microsoft Playwright 1.62.1 runtime 4 and target 8 PASS across Chromium desktop/mobile, Firefox and WebKit;
- Google Lighthouse 13.4.1 five governed runs PASS: performance median `0.99`, accessibility/best-practices/SEO medians `1.00`.

The two temporary databases were dropped and the isolated PostgreSQL server was stopped after the gates.

## Failures converted to memory

`LIB-FAIL-1258` through `LIB-FAIL-1264` preserve the missing global Go path, PostgreSQL child lifecycle, PowerShell array boundary recurrence, Next build lock after wrapper timeout, atomic patch anchor mismatch, assumed Go command target and rejected compound background launch. Each has an explicit correction and repeated proof; no failed attempt is counted as a PASS.

## Honest production boundary

This proves the reusable local capacity/calendar primitive and its exact web/backend integration. It does not prove employee/resource rosters, holidays, cancellation/no-show policy, notification delivery, target IdP/CDN/WAF, production recovery, target load, offensive security, deployment/canary/rollback or business acceptance. A concrete project remains blocked from production until those target receipts exist together.
