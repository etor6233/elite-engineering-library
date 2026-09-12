# Franchise Handover Server Evidence — V136

Date: 2026-08-30  
Decision: `REBUILD_VERIFIED / CONDITIONED`

## Authority and provenance

Microsoft Business Central separates posted shipment and invoice records, updates ledgers during posting and requires reversal/undo for critical corrections: <https://learn.microsoft.com/en-us/dynamics365/business-central/ui-post-sales>. Microsoft Dynamics 365 Field Service inspections use published/versioned inspection definitions, required responses and read-only completion behavior: <https://learn.microsoft.com/en-us/dynamics365/field-service/inspections-overview>.

The V136 implementation is local `AUTHORED` Go/TypeScript governed by those boundaries; it is not presented as Microsoft source. It deliberately calls its result an authenticated acceptance record, not a qualified electronic signature or universal legal proof.

## Corrected trust boundary

The prior browser command accepted an arbitrary `evidence_sha256`. V136 removes that field. The browser can provide only:

- organization already admitted by the session;
- optimistic handover version;
- explicit receipt confirmation;
- the serial number physically observed by the customer.

The HTTP server calculates SHA-256 over the exact bounded JSON body. The authenticated subject comes from OIDC. PostgreSQL accepts the transition only when tenant, organization, customer, handover, version, presented state and durable `inventory.stock_unit.serial_number` all match. State, timestamp, evidence hash and outbox event commit in one transaction.

## Exact artifacts

- `GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.5.1`, SHA-256 `71468e8f4039b378cd910eadd0e83eb60260ba6124b4199783a3f85b6370c3af`, 24 files.
- `TS-FRANCHISE-JOURNEY-PORTALS 0.6.0`, SHA-256 `2b8269cc075e447fe69bd58a057454fe784e96550999d750355193728e6f0bcb`, 14 files.
- Backend profile: 24 packs / 238 implementation files.
- Web profile: 6 packs / 84 implementation files.

## Gates executed from Markdown

1. A fresh backend composition applied migrations 0001–0015 to PostgreSQL 18.6.
2. The real integration rejected `WRONG-SERIAL`, then accepted only `SERIAL-JOURNEY` for the exact customer and version.
3. Full `go test ./... -count=1`, `go vet ./...` and four command builds passed.
4. A fresh web composition installed from the frozen lock offline with zero downloads.
5. Nine Vitest files / 35 tests, `next typegen`, strict TypeScript, license tests and production build passed.
6. The production manifest contains `/customer/handovers`; strict Zod rejects the removed browser evidence field.
7. PostgreSQL was explicitly stopped after both real database runs.

## Remaining production conditions

A project must still select its legal acceptance statement, signature provider when legally required, checklist/version, attachment/photo retention, reversal authority, IdP and edge. V136 prevents fabricated evidence hashes and proves serial-bound authenticated acceptance; it does not invent those project/legal decisions.
