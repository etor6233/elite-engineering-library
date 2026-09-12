# Return Fiscal Credit Note Execution Inventory — V146

Date: 2026-08-31  
Scope: durable full-refund-to-full-credit-note owner for admitted ARCA voucher classes A, B and C.

## Result

`GO-RETURN-FISCAL-CREDIT-NOTE-WORKER@0.1.0` adds eight `AUTHORED` files and composes with the existing fiscal, return-effect, refund and accounting owners. The backend profile now reconstructs 29 packs and 298 source files.

The worker claims only fiscal return effects. It requires successful inventory, refund and accounting executions; a succeeded refund whose amount and currency equal the one authorized original invoice; and a posted accounting reversal with the same total. It then maps only `1→3`, `6→8` or `11→13`, copies the exact authorized invoice amounts/taxes, delegates issuance to the existing fiscal owner and persists one immutable original-to-credit-note link. Pending authorization remains retryable, rejection blocks, and only an authorized credit note closes the fiscal effect. Deterministic IDs and the fiscal idempotency contract reconcile a crash after request creation without issuing a duplicate.

## Upstream authority and provenance

- ARCA voucher authority: https://www.arca.gob.ar/facturacion/comprobantes/
- ARCA WSFEv1 developer contract: https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf
- Microsoft Dynamics 365 return-order lifecycle: https://learn.microsoft.com/en-us/dynamics365/supply-chain/sales-marketing/sales-returns
- PostgreSQL 18 transaction/constraint authority: https://www.postgresql.org/docs/18/

All eight files are local `AUTHORED` orchestration governed by those locked contracts. None is represented as copied ARCA, Microsoft or PostgreSQL source code. The worker does not reimplement SOAP, credentials, fiscal sequencing, refund providers or ledger posting; it composes the already admitted owners.

## Executed gates

- pack materialization: 8/8 blocks;
- all eight reconstructed files matched the tested author tree byte-for-byte;
- clean profile composition: 29 packs / 298 source files;
- Go 1.26.7: `go test ./...` and `go vet ./...` PASS;
- all eight discovered Go main packages built independently;
- PostgreSQL 18.6 clean database applied migrations 0001–0024 with `ON_ERROR_STOP=1`;
- SQL contract `0024_return_fiscal_credit_note.test.sql` PASS;
- `TestReturnFiscalRequestsOnceAndWaitsForAuthorization` PASS against the Markdown-reconstructed database;
- integration observed one credit invoice, a pending result followed by authorization, two effect attempts, an immutable link and no duplicate request on replay;
- PostgreSQL was stopped explicitly after the gate.

## Failures converted to memory

`LIB-FAIL-1444` through `LIB-FAIL-1448` record the nonexistent assumed adapter path, ledger-anchor recovery, canonical plan discovery, empty PowerShell output and URL interpolation failure. Each failed closed before an unsupported PASS could be claimed.

## Boundary retained

This evidence does not claim live ARCA authorization, homologation acceptance, legal/tax correctness for a concrete taxpayer, partial credit notes, debit notes, invoice rendering, non-identical price/variant exchange settlement, production credentials or production readiness. A project must still prove its real CUIT, certificate, point of sale, approved fiscal mappings, refund/accounting policy, accountant/business acceptance and target controls. No project is production-ready until its real CDN/WAF, IdP, providers, PostgreSQL recovery, load, offensive security, deployment, rollback and business acceptance are demonstrated.
