# Delivery Exception and Return Authorization — V139

Date: 2026-08-30  
Admission: `REBUILD_VERIFIED / CONDITIONED`

## Scope and provenance

V139 closes the technical delivery-discrepancy boundary without creating a second order, inventory, payment, accounting, fiscal, identity or customer owner. The implementation is `AUTHORED`; it is not presented as Microsoft code. Microsoft Business Central governs return/cancellation traceability and exact-cost source linkage, Microsoft Field Service governs inspection discrepancy reporting and follow-up, and the public Microsoft BCApps repository is an inspectable code reference:

- https://learn.microsoft.com/en-us/dynamics365/business-central/sales-how-process-sales-returns-cancellations
- https://learn.microsoft.com/en-us/dynamics365/field-service/inspections-reporting
- https://github.com/microsoft/BCApps/blob/main/src/Layers/BE/BaseApp/Sales/History/SalesReturnReceipt.Report.al
- https://github.com/microsoft/BCApps/blob/main/src/Layers/IT/BaseApp/Sales/History/ReturnReceiptHeader.Table.al
- https://github.com/microsoft/BCApps/blob/main/src/Layers/CH/BaseApp/Sales/Document/SalesReturnOrder.Page.al

The project must still approve reason codes, return/exchange policy, physical receipt/inspection, refund, inventory, exact-cost accounting and fiscal execution. An authorization is deliberately not represented as completion of those downstream effects.

## Materialized contract

- `GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.7.0`: 30 files; pack SHA-256 `b267d10655bf3a37a765a9d3f5e3b6aa0a27ae635186fff7aa3ff6f0316bd8ca`.
- `TS-FRANCHISE-JOURNEY-PORTALS 0.8.0`: 14 files; pack SHA-256 `c10f687c9038f22b374b8757d04c4b10f3bca2d393e47ff638bca23f9584f6a1`.
- backend composition: 24 packs / 244 files; record SHA-256 `120c6533aa259dd1a2b5c8f22601f651f88389a41de802e34f7542e147dbc20d`.
- web composition: 6 packs / 84 files; record SHA-256 `0aa6333c32fa0608e1d4457e83832bedd6f09fe805af6a196a7edf5ce6ededb2`.

New PostgreSQL files:

| File | Bytes | SHA-256 |
|---|---:|---|
| `0017_delivery_exception_and_return_authorization.up.sql` | 6,060 | `b7c0c066915a649eed82ea811ac87fe90b54ca1ad4f7166f4235979ac84a82fc` |
| `0017_delivery_exception_and_return_authorization.down.sql` | 1,044 | `bb486c2eb7f8a889fdaf036f162547594bfce9341263723f270ecf0dc28111b6` |
| `0017_delivery_exception_and_return_authorization.test.sql` | 903 | `a7643f33eb07d288067467a2c61b1788fe28ed067c03135ae7ecfff22e0b6dc8` |

## Demonstrated behavior

1. Only the verified customer can reject a presented handover in its tenant/organization/customer scope; the server owns command digest, exception identity and event identity.
2. The rejection creates a durable versioned exception and outbox event; stale, cross-customer and repeated commands fail closed.
3. `correct-and-represent` creates a new prepared handover linked through `supersedes_handover_id`; the rejected handover remains immutable history.
4. `return` and `exchange` use one closed resolution contract and create an immutable authorization linked to the original exception, handover, order, stock unit, customer and exact-cost source order.
5. PostgreSQL rejected direct mutation of the return authorization. The authorization does not silently change stock, refund money or post accounting/fiscal effects.
6. The web BFF admits only the closed customer rejection/operator resolution fields, enforces server-side session/permissions and displays authorization without claiming downstream completion.

## Reconstruction and gates

- Individual packs materialized 30/30 and 14/14; all 44 paths were byte-identical to the tested author trees.
- Clean backend profile composed 24/244. `go test ./...`, `go vet ./...` and four command builds passed.
- Fresh PostgreSQL 18.6 databases applied migrations 0001–0017 with `ON_ERROR_STOP=1`; the SQL structure test and full repository integration passed from both author and Markdown-reconstructed trees.
- The real repository integration proved customer-scope rejection, successor creation, return authorization, exact original order/stock/cost-source linkage, one-time resolution and authorization immutability.
- Clean web profile composed 6/84; `pnpm install --frozen-lockfile --offline` downloaded zero packages, strict TypeScript passed, nine Vitest files/37 tests passed, Next.js 16.3.2 production build and license report passed.
- Both temporary PostgreSQL runs were explicitly stopped and `pg_ctl status` confirmed no server running.

## Learned failures

Library failure rows 1380 through 1388 preserve two atomic patch rejections, historical path/name assumptions, two failed PostgreSQL launch/diagnostic stages, one PowerShell URL interpolation defect, one catalog-path recurrence and one undersized verifier window. None was hidden or counted as PASS; each corrected repetition is demonstrated above, including the persistent-session structural PASS and the full executable Audit PASS.

## Remaining condition

This block is reusable technical infrastructure. A project is not production-ready until its owner configures and proves the complete physical return receipt/inspection, inventory disposition, exchange fulfillment, refund/payment reconciliation, exact-cost accounting and jurisdiction-approved fiscal reversal across their real providers and policies.
