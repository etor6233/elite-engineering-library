# ARCA Associated Voucher Execution Inventory — V145

Date: 2026-08-31  
Scope: strict full-credit-note association lane for ARCA WSFEv1 voucher classes A, B and C.

## Result

The backend profile reconstructs 28 packs and 290 source files. Three coordinated packs now preserve and transport one exact associated voucher for admitted credit-note pairs `3→1`, `8→6` and `13→11`:

- `GO-ARCA-FISCAL-ISSUANCE-API@0.6.0` persists an immutable association to one local authorized original invoice;
- `ARCA-WSFE-UDS-WORKER@0.3.0` transports the exact identity through strict local JSON;
- `MICROSOFT-ARCA-WSFE-SOAP-ADAPTER@0.3.0` maps it to the Microsoft-generated `CbteAsoc` contract.

The profile also now selects all six files from `MICROSOFT-ARCA-WSFE-GENERATED-CLIENT@0.2.1`, so an empty reconstruction includes the generator and its regression instead of depending on an omitted local tool.

## Upstream authority and provenance

ARCA remains the protocol and legal authority:

- current electronic-invoice web-service index: https://www.arca.gob.ar/fe/ayuda/webservice.asp
- WSFEv1 developer manual: https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf
- ARCA voucher guidance: https://www.arca.gob.ar/facturacion/comprobantes/

The exact homologation WSDL was acquired under the existing byte/SHA lock. Microsoft `dotnet-svcutil` 8.0.0 generated `CbteAsoc` with `Tipo`, `PtoVta`, `Nro`, `Cuit` and `CbteFch`; the generated proxy hash matched the locked value. Microsoft documents the generator at https://learn.microsoft.com/dotnet/core/additional-tools/dotnet-svcutil-guide. PostgreSQL 18.6 governs the deferred constraints and transaction behavior: https://www.postgresql.org/docs/18/sql-createtrigger.html.

All 12 modified/new materialized files are `AUTHORED` local implementation. None is represented as copied ARCA, Microsoft or PostgreSQL code. The two generated proxies are derived build artifacts and remain outside the canonical Markdown/source count.

## Executed gates

- the three changed packs materialized 28/28, 14/14 and 13/13 blocks;
- 12 changed files matched the author tree byte-for-byte after Markdown reconstruction;
- full profile composition: 28 packs / 290 source files;
- official generator contract PASS, then offline homologation generation PASS for two clients with locked hashes;
- `go test ./...` and `go vet ./...` PASS on Go 1.26.7;
- seven discovered Go main packages built independently;
- SOAP adapter PASS: 15 invoice contracts and 15 parameter contracts, warning-free .NET 10.0.400 builds;
- UDS worker PASS: full Go suite, warning-free .NET build and 12 worker tests;
- PostgreSQL 18.6 clean database applied migrations 0001–0023 with `ON_ERROR_STOP=1`;
- SQL 0023 contract PASS, repository integration PASS, and author-tree 0023 down/up PASS;
- PostgreSQL was stopped explicitly after both author and Markdown gates.

The integration proves that the client supplies only the original local invoice ID; PostgreSQL resolves the authorized original identity, requires the same organization and taxpayer, stores exact type/POS/number/date, rejects unsupported voucher classes or missing association, makes the association immutable and rehydrates it for the durable fiscal worker.

## Boundary retained

This V145 evidence does not claim an ARCA live authorization, homologation acceptance, fiscal correctness of project amounts, partial credit notes, debit notes, legal invoice rendering, production credentials or production readiness. The next owner must create the credit-note request only after the return remedy and accounting/inventory evidence are complete, and a concrete project must still prove ARCA homologation with its real CUIT, certificate, point of sale, approved fiscal mappings and accountant/business acceptance.
