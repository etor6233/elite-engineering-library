# Return Receipt, Disposition and Effect Requests — V140

Date: 2026-08-30  
Admission: `REBUILD_VERIFIED / CONDITIONED`

## Scope and provenance

V140 turns an immutable return/exchange authorization into physical receipt evidence, one disposition and explicit requests to the existing downstream owners. The implementation is `AUTHORED`; it is not presented as Microsoft code. It is governed by current Microsoft Business Central sales-return/exact-cost guidance, Microsoft Field Service inspection-response/follow-up guidance and public MIT-licensed Microsoft BCApps return-order/receipt code:

- https://learn.microsoft.com/en-us/dynamics365/business-central/sales-how-process-sales-returns-cancellations
- https://learn.microsoft.com/en-us/dynamics365/field-service/inspections-reporting
- https://github.com/microsoft/BCApps/blob/main/src/Layers/CH/BaseApp/Sales/Document/SalesReturnOrder.Page.al
- https://github.com/microsoft/BCApps/blob/main/src/Layers/IT/BaseApp/Sales/History/ReturnReceiptHeader.Table.al
- https://github.com/microsoft/BCApps/blob/main/src/Layers/BE/BaseApp/Sales/History/SalesReturnReceipt.Report.al

The local contract does not copy those files and does not claim that a `requested` downstream effect has already moved inventory, refunded money, fulfilled an exchange or posted accounting/fiscal documents.

## Materialized contract

- `GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.8.0`: 33 files; pack SHA-256 `160ad590953b07a87c3d2d3b195034dd4a747221268be03797fe3b75e933c2e2`.
- `TS-FRANCHISE-JOURNEY-PORTALS 0.9.0`: 14 files; pack SHA-256 `26f2c6b48ca1525f6c50a2fb1784ee92b5987784012ede045b3e4d598410dfdc`.
- backend composition: 24 packs / 247 files; record SHA-256 `3ba164ad69d3d4f2f2a76a2d8bf1d059639525b8458ff11c19f759a05349f8d3`.
- web composition: 6 packs / 84 files; record SHA-256 `6e93001d35705f63e07ec9e5418af6865d44f39105fed1827c92b790a07b8444`.

New PostgreSQL files:

| File | Bytes | SHA-256 |
|---|---:|---|
| `0018_return_receipt_disposition_effects.up.sql` | 6,030 | `db8b62aeab5d3bc23881ce78f1edb6df35dff6c617a6d4223feb8faef6d8ff4a` |
| `0018_return_receipt_disposition_effects.down.sql` | 866 | `a62f62aafb65fa0c5f6ed2047b972ee260d1c7b5e12c4348eeaf61a479673b19` |
| `0018_return_receipt_disposition_effects.test.sql` | 773 | `74ac0b6ded07730909db45d8547f372445c46a82f12704bb2c1f878f4e2ac47b` |

## Demonstrated behavior

1. Receipt derives organization, order, stock unit and customer from the authorization and rejects a serial different from the durable inventory record.
2. Condition is closed to `sealed/opened/damaged/incomplete`; notes, evidence digest, receiver and timestamp are immutable.
3. Exactly one receipt is allowed per authorization and exactly one disposition per receipt.
4. Disposition accepts only `quarantine/restock/repair/scrap`; refund versus exchange is derived from the authorization, not the browser.
5. The same transaction creates immutable idempotent requests for inventory, refund or exchange, accounting and—when refund applies—fiscal owners, plus an outbox event.
6. The franchise BFF hashes the admitted receipt command server-side and rejects browser-owned evidence, request IDs, refund/payment, stock, accounting, fiscal or completion-state fields.
7. The UI reports every downstream effect as `requested` and explicitly does not represent it as completed.

## Reconstruction and gates

- Individual packs materialized 33/33 and 14/14; all 47 paths were byte-identical to tested author trees.
- Clean backend profile composed 24/247. `go test ./...`, `go vet ./...` and four command builds passed.
- PostgreSQL 18.6 applied migrations 0001–0018 with `ON_ERROR_STOP=1` in author and clean Markdown databases. SQL structure tests and the complete repository integration passed.
- Negative paths proved wrong serial, duplicate receipt, duplicate disposition and in-place receipt mutation are rejected. Positive assertions proved four refund-path owner requests and complete return-case projection.
- Clean web profile composed 6/84; frozen offline installation downloaded zero packages, strict TypeScript passed, nine Vitest files/38 tests passed, production build and license report passed.
- Both temporary PostgreSQL runs were explicitly stopped and status confirmed no server running.

## Learned failure

Library failure row 1389 records the rejected multi-directory `go fmt` harness command. The corrected `gofmt` invocation and all subsequent tests passed; the failed command was not counted as evidence.

## Remaining condition

V140 completes the durable technical handoff, not the external effects. A project remains conditioned until its selected inventory, payment or exchange fulfillment, accounting and fiscal adapters consume each request idempotently, reconcile provider truth, expose failures/retries, and pass real credentials, sandbox/live canaries, load, security, rollback and business acceptance.
