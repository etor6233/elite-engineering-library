# Microsoft BC partial transfer posting — reconstruction V164

Date: `2026-08-31`

## Outcome

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.8.0 materializes 46 files and the enterprise backend composes 29 packs / 338 files. V164 adds migration 0030 up/down as two `ADAPTED` files and updates the existing transfer derivation, domain, PostgreSQL adapter, HTTP adapter and their tests. No Go or SQL file is represented as Microsoft-authored code and this evidence does not claim the complete Business Central WMS, planning or costing product.

The executable result supports cumulative partial shipment and receipt against one transfer line. A shipment must name one exact registered pick, owns a durable request/posting identity and can consume that pick once. A receipt owns an independent request/posting identity, creates its own put-away and slices only shipped-not-received cost components. Exact replay succeeds even with the pre-post version; changed payload, excessive quantity, reused pick and stale concurrent posting fail closed.

## Fixed official authority

- Microsoft `BCApps` commit: `2eae56d704a1fd035d104f333602aea7091b7749`.
- License: upstream Microsoft AL under `MIT`; portable Go/SQL remains `LicenseRef-Workspace-Owner AND MIT` and is declared `ADAPTED`.
- `TransferLine.Table.al`: explicit quantities to ship/receive and rejection of receipt beyond quantity in transit.
- `TransferOrderPostShipment.Codeunit.al`: cumulative shipped quantity and posted shipment creation.
- `TransferOrderPostReceipt.Codeunit.al`: cumulative received quantity and posted receipt creation.
- Microsoft Learn transfer-mode matrix: non-direct transfer through an in-transit location admits partial posting.

Pinned URLs:

- `https://github.com/microsoft/BCApps/blob/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Inventory/Transfer/TransferLine.Table.al`
- `https://github.com/microsoft/BCApps/blob/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Inventory/Transfer/TransferOrderPostShipment.Codeunit.al`
- `https://github.com/microsoft/BCApps/blob/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Inventory/Transfer/TransferOrderPostReceipt.Codeunit.al`
- `https://learn.microsoft.com/en-ca/dynamics365/business-central/inventory-how-transfer-between-locations`

## Files and reconstruction

Changed or added materialization paths:

- `db/migrations/0030_bulk_transfer_partial_posting.up.sql` — `ADAPTED`;
- `db/migrations/0030_bulk_transfer_partial_posting.down.sql` — `ADAPTED`;
- `docs/inventory/MICROSOFT_BC_TRANSFER_DERIVATION.md` — `AUTHORED`;
- `internal/inventorycontrol/bulk_transfer.go` and test — existing `ADAPTED` owner;
- `internal/platform/postgres/bulk_transfer.go` and integration test — existing `ADAPTED` owner;
- `internal/platform/httpapi/bulk_transfer.go` and test — existing `AUTHORED` boundary.

The canonical updater replaced nine blocks atomically. A fresh compositor self-test passed and a fresh backend reconstruction produced 338 files. SHA-256 comparison between all nine working-tree sources and the reconstructed outputs passed 9/9; `gofmt -d` returned empty.

## Executed gates

- Go archive: official `go1.26.7.windows-amd64.zip`, 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`.
- PostgreSQL archive: official 18.6 Windows x64 binaries, 343,808,005 bytes, SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`.
- PostgreSQL 18.6 clean cluster: migrations 0001–0030 applied with `ON_ERROR_STOP`.
- Focused partial-transfer integration: PASS twice consecutively on the same database.
- Full `go test ./... -count=1` with PostgreSQL: PASS, including the complete platform/postgres suite.
- `go vet ./...`: PASS.
- `go build ./cmd/...`: PASS.
- Migration 0030 clean down/up: PASS after the integration cleanup proved zero shipment/receipt rows.
- Final `VERIFY_LIBRARY_PASS`: 94 packs / 1,055 materialized files / 520 Markdown / backend 338 / web 84.
- Final `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 94 packs / 121 fixed upstream sources, including all declared source-profile and executable offline gates; live/provider/target gates remained explicitly skipped rather than inferred.

The focal proof posts one unit, receives one, concurrently attempts the remaining two through one exact pick and admits one winner, then receives one plus one. It proves two shipment rows, three receipt rows, in-transit quantity/value `2/20 → 1/10 → 0/0`, source/destination quantity `2/3`, destination transfer cost `30`, stable lot identity, exact request replay, divergent replay rejection, three registered put-away activities, cancellation before shipment and rejection after cancellation. The same test also preserves the V163 exact specific receipt cost `17` while leaving the unrelated `11` layer intact.

## Failure learning

The first V164 integration rejected the initial shipment. The exact cause was not hidden: warehouse pick reservations use `<transfer-line>:<sequence-6>`, while the new query compared that value to the unsuffixed line ID. The fix derives the exact reservation identity from the registered activity source line and immutable activity-line sequence. The failure is permanently recorded in `LIBRARY_FAILURE_LEARNING_LEDGER.md`; the two consecutive focal passes and the global suite are its non-recurrence evidence.

## Remaining conditions

This is reusable implementation evidence, not target production acceptance. Break-bulk, cross-docking, replenishment/CTP, freight and landed-cost allocation, approved accounting periods, ERP/WMS reconciliation, representative load, operator UX, cloud/network/IdP and live business acceptance remain project conditions. No project may claim production-ready until its exact packaged artifact proves those gates in its real target.
