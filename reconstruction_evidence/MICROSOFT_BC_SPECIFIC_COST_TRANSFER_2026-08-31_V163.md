# Microsoft BC exact specific-cost transfer reconstruction evidence — V163

Date: 2026-08-31

## Outcome and provenance

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.7.0 materializes 44 files and the enterprise backend composes 29 packs / 336 files. V163 adds migration 0029 up/down as two `ADAPTED` files and updates the existing transfer domain, PostgreSQL adapter, integration test and derivation document. No Go or SQL file is represented as Microsoft-authored code.

The governing implementation is an explicit translation of Microsoft Business Central contracts fixed at signed `microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`, MIT:

| Source | Blob | Bytes | SHA-256 |
|---|---|---:|---|
| `TransferLine.Table.al` | `c78a9a3583498c2cb4ca0698ab20b473e93cfedc` | 108,502 | `49b906242b8ae707c94459301c560f489d8a885a229949f446cb40d30817fdd3` |
| `TransferOrderPostShipment.Codeunit.al` | `211bf8f4574b57a718524d659392ce29f874e888` | 53,015 | `1e9ad92a31e0e8751d1660e6eaadf48defbba903b85a60e9486ae5a06f7af708` |
| `TransferOrderPostReceipt.Codeunit.al` | `adede7876bfa59682c8693625676b3bd723fb4a7` | 49,811 | `5523af9ce81b5f972919681beaf48c75edb68fa7bfb92449ce557548c8d6e039` |
| `ItemApplicationEntry.Table.al` | `17dd82953a3615e541a617d0db4ca6da2f685430` | 49,965 | `3df93f7a3ade1c6c73e715a991514f1d8a56fa5565552d1d3921c2bedf84ee16` |

Microsoft Learn independently states that transfer shipment precedes receipt, partial posting depends on transfer mode, and serial/lot identities shipped from one location must be received unchanged:

- https://learn.microsoft.com/en-ca/dynamics365/business-central/inventory-how-transfer-between-locations
- https://learn.microsoft.com/en-us/dynamics365/business-central/inventory-how-work-item-tracking
- https://learn.microsoft.com/en-us/dynamics365/business-central/reservation-entry-table-features-that-update-the-table

## Executable contract

For `specific` costing, the command must name `specific_receipt_entry`. The selector is stored on the transfer line, included in idempotent replay equality and constrained to the tenant's immutable inventory entry. Shipment locks and consumes only the selected open cost layer for the picked item and unchanged lot. Missing selector, FIFO-with-selector, changed replay, wrong organization/item/lot, insufficient selected quantity, stale version and duplicate posting fail before commit. Destination receipt recreates the exact selected cost component and then enters the existing put-away workflow.

V163 also closes the uncovered nullable-lot defect: transfer receipt reads `lot_id` through `coalesce` and persists it through `nullif`, so an admitted non-lot specific-cost item no longer fails during receipt.

## Reproducible gates executed

- Go official `go1.26.7 windows/amd64`; archive 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`;
- PostgreSQL 18.6 official EDB archive; 343,808,005 bytes, SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`;
- supply reconstruction: 44/44; backend canonical composition: 29 packs / 336 files;
- migrations 0001–0029 on a clean page-checksummed cluster: PASS;
- migration 0029 down/up: PASS;
- complete Go suite without PostgreSQL: PASS;
- focal PostgreSQL transfer test twice consecutively: PASS;
- complete Go suite with PostgreSQL: PASS;
- `go vet ./...`: PASS;
- builds `cmd/api` and `cmd/electromobility-api`: PASS.

The specific-cost proof created two source receipts at costs `11.0000` and `17.0000`, selected only the second, preserved first-layer remaining quantity `1.000000`, reduced the selected layer to `0.000000` and received exactly `17.0000` at the destination. FIFO/lotted transfer and cancellation evidence from V162 continued passing in the same test.

## Failure-derived hardening

`LIB-FAIL-1541` through `LIB-FAIL-1549` record the exact maintenance, packaging and executable failures observed in V163. The material defect was `LIB-FAIL-1549`; the new non-lot corpus prevents recurrence. All entries are closed by demonstrated regressions rather than prose-only disposition.

## Deliberate boundary

This is not proof of partial shipment/receipt, breakbulk/UOM conversion, cross-docking, replenishment planning, average/standard/LIFO, landed cost, accounting periods, external WMS reconciliation, load, offensive security or production operation. Those capabilities remain fail-closed until their own official derivation and executable evidence exist.

Final root gates observed after this evidence entered the canonical tree:

- `VERIFY_LIBRARY_PASS packs=94 materialized_files=1053 markdown_files=519`;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=94 upstream_sources=121`;
- both processes exited `0`.

Provider/account/cost/cloud/target gates reported as `SKIPPED` remain explicit conditions and are not converted into executed production evidence.
