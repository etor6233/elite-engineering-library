# Microsoft BC bulk transfer / in-transit reconstruction evidence — V162

Date: 2026-08-31

## Outcome and provenance

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.6.0 adds nine materializable files: six `ADAPTED` Go/PostgreSQL migration/domain/adapter/test files and three local `AUTHORED` derivation/HTTP files. `GO-ELECTROMOBILITY-APPLICATION` 1.9.0 only wires the new service into the existing inventory repository and CSPRNG. No file is represented as Microsoft-authored Go, and this evidence does not claim the complete Business Central WMS, planning or costing product.

Official Microsoft authority is fixed to signed `microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`, MIT:

| Source | Blob | Bytes | SHA-256 |
|---|---|---:|---|
| `TransferHeader.Table.al` | `0f61480f5837ed9ac48d380c10c8434f7364a5d7` | 82,670 | `952ad36743a2ed14313eb3c99a52fd1efbdd8891344919a8cbbc0417c3b8f44c` |
| `TransferLine.Table.al` | `c78a9a3583498c2cb4ca0698ab20b473e93cfedc` | 108,502 | `49b906242b8ae707c94459301c560f489d8a885a229949f446cb40d30817fdd3` |
| `TransferOrderPostShipment.Codeunit.al` | `211bf8f4574b57a718524d659392ce29f874e888` | 53,015 | `1e9ad92a31e0e8751d1660e6eaadf48defbba903b85a60e9486ae5a06f7af708` |
| `TransferOrderPostReceipt.Codeunit.al` | `adede7876bfa59682c8693625676b3bd723fb4a7` | 49,811 | `5523af9ce81b5f972919681beaf48c75edb68fa7bfb92449ce557548c8d6e039` |

Microsoft Learn independently documents transfer order shipment, explicit in-transit quantity, later receipt, immutable tracking between shipment and receipt and the outbound-demand/inbound-supply duality:

- https://learn.microsoft.com/en-us/dynamics365/business-central/inventory-how-transfer-between-locations
- https://learn.microsoft.com/en-us/dynamics365/business-central/inventory-how-work-item-tracking
- https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-transfers-in-planning

## Executable scope

One existing PostgreSQL inventory owner now performs:

1. idempotent released transfer creation with source/destination organizations, receive bin and in-transit code;
2. existing `transfer-outbound` warehouse pick and SHIP registration;
3. versioned shipment with one concurrency winner, source quantity/reservation consumption and immutable cost applications;
4. durable `bulk_inventory_in_transit` evidence grouped by transfer/item/lot;
5. destination receipt that recreates every FIFO cost component without value invention;
6. atomic destination put-away generation and later registration through the existing warehouse workflow;
7. released-transfer cancellation that cancels open picks, releases reservations and prevents shipment;
8. strict JSON, tenant and dual-organization authorization at HTTP plus organization matching inside SQL.

Specific-cost transfer is deliberately rejected because V162 has no command field selecting the exact source receipt. Partial shipment/receipt, breakbulk, cross-docking, replenishment planning/CTP, landed cost, accounting periods, title/Incoterms and external WMS reconciliation remain conditioned rather than simulated.

## Reproducible gates executed

- Go official `go1.26.7 windows/amd64`;
- PostgreSQL 18.6 official EDB archive, 343,808,005 bytes, SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`, isolated cluster with page checksums;
- supply pack materialization: 42/42; application: 5/5; backend composition: 29 packs/334 files;
- clean migrations 0001–0028: PASS; migration 0028 down/up: PASS;
- exact author-to-Markdown hash comparison for ten supply paths plus application main: zero mismatches;
- full `go test ./... -count=1`, `go vet ./...`, builds `cmd/api` and `cmd/electromobility-api`: PASS from the final Markdown composition;
- PostgreSQL focal flow: idempotent replay; FEFO pick to SHIP; two concurrent shipments with exactly one winner; in-transit quantity `3.000000`, cost `30.0000`, original lot; one receipt; destination put-away; conservation `5 = 2 + 3`; cancelled pick/reservation and rejected post-cancel shipment: PASS;
- journey cleanup regression fixed in pack 0.8.1 and executed twice consecutively on one database: PASS.

Final root gates observed after this evidence entered the canonical tree:

- `VERIFY_LIBRARY_PASS packs=94 materialized_files=1051 markdown_files=518`;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=94 upstream_sources=121`;
- executable audit exit code: `0`.

Live provider, account, cost, cloud, recovery, load and production gates reported as `SKIPPED` remain explicit admission conditions; this evidence does not convert them into executed proof.

## Failure-derived hardening

V162 registered every observed failure through `LIB-FAIL-1540`. Material failures improved: parameter/path preflight, complete runtime extraction checks, request replay numeric equality, immutable teardown order, repeatable return-effect cleanup, and pgx cursor-before-mutation discipline. No failing run is counted as a PASS.

## Production boundary

This proves reconstructible code and local PostgreSQL behavior, not production readiness. A concrete project is not ready for production until its CDN/WAF, IdP, providers, PostgreSQL/recovery, representative load, offensive security, deploy/canary/rollback and business acceptance are demonstrated on the real target.
