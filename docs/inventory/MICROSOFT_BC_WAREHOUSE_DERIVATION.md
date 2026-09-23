# Microsoft BC operational warehouse derivation

This portable Go/PostgreSQL boundary is `ADAPTED`; it is not Microsoft-authored Go and is not a replacement for Business Central. Its behavioral authority is Microsoft BCApps at commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`, under the repository MIT license.

## Exact official source identities

| BCApps path under `src/Layers/W1/BaseApp` | Git blob | Bytes | SHA-256 |
|---|---|---:|---|
| `Warehouse/Document/WhsePostReceipt.Codeunit.al` | `dbd6327a935a69c564d112203eaf42ee845ba393` | 105803 | `9aae2fa85c20e12e70b88b3214fd4b4defcb13b6e767c225251d413eda626989` |
| `Warehouse/Activity/CreatePutaway.Codeunit.al` | `df3011a0a74ace99e3d1788996ba750a7732a1e6` | 77370 | `b004c0663d237585d8b57a06356431acb2302cdfbbd30f7b82154d282f325b2b` |
| `Warehouse/Activity/CreatePick.Codeunit.al` | `07be15a7effeeda1ac749ef033d20d817bb5b218` | 285538 | `6ca4acba3eebb46c719267ed788d8efc6ce097005db0d1a560489797e03b35e5` |
| `Warehouse/Tracking/WhseItemTrackingFEFO.Codeunit.al` | `9568f0f687384a65ac6dc4135273cc44b7256c4e` | 17061 | `24f0b95de17d6d2a042b97fe2cefb0ded952df8db2f0ebfefdb75c47aca09965` |
| `Warehouse/Activity/WarehouseActivityHeader.Table.al` | `4e3e2cf15ed00cfdc627898a1fafb54a8c8eca29` | 47747 | `34f6f3b633aed6c77c1cf286e9528d449fb255dcc219531946b43b59ce862030` |
| `Warehouse/Activity/WarehouseActivityLine.Table.al` | `83d6dc89fa2b088e05d9ac224fe0ab9e3c1921b9` | 181260 | `f4ab727d452c158b1f397e83708150f04b9fac22262e378bbca0bd9a8ed6836b` |

Microsoft Learn authorities checked on 2026-08-31:

- `https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-warehouse-management`
- `https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-inbound-warehouse-flow`
- `https://learn.microsoft.com/en-us/dynamics365/business-central/warehouse-how-to-put-items-away-with-warehouse-put-aways`
- `https://learn.microsoft.com/en-us/dynamics365/business-central/warehouse-how-to-pick-items-for-warehouse-shipment`

## Narrow invariants transferred

- Posting a warehouse receipt creates inventory in a RECEIVE bin before a separate put-away is registered.
- Put-away planning selects only eligible PUT AWAY or PUTPICK bins, accounts for maximum quantity and already-open placements, and prefers default/high-ranking bins.
- Registration records paired Take/Place movement evidence; an open activity can transition once to registered or cancelled under optimistic version control.
- Picks use only PICK or PUTPICK bins. RECEIVE, SHIP, PUT AWAY, QC, blocked and unapproved dedicated bins are excluded.
- FEFO excludes expired/blocked tracking identities and existing reservations. It orders by expiration, lot identity, registration age and then bin ranking; selection and reservation commit atomically.
- A registered pick moves both quantity and its reservation to the SHIP bin, keeping issuance separate. Cancellation is allowed only while the pick is open and releases every reservation atomically.
- A pick resolves its requested target UOM first. It consumes exact target-UOM composition before considering a larger source UOM, and that fallback is admitted only when `allow_breakbulk=true`, matching the explicit `AllowBreakbulk` branch in the fixed BCApps replenishment/pick behavior.
- Every packaging-aware line persists From/To UOM, quantity and factor. Open work reserves both physical base quantity and the exact source composition; registration consumes the source package, records immutable breakbulk Take/Place evidence when conversion is required, moves the requested target composition and retains any package remainder as exact base composition. Cancellation releases both reservations.
- Receipt evidence, activity lines and inventory entries are append-only. Request IDs and source identities prevent silent duplicate effects.

## Deliberate limits

The portable lane admits only the packaging-aware customer/transfer picking and bin-replenishment behavior demonstrated by migrations 0033–0036 and its PostgreSQL tests. Cubage/weight, warehouse classes, barcode/device flows, partial line handling, serial tracking, picks for production/assembly/projects, planning/CTP and the full Business Central extension-event surface remain rejected until their own official derivation and executable evidence exist. Cross-docking and cross-organization transfer are governed separately by their fixed derivations and tests.

## Regression mapping

Migration 0027 enforces bin types, ranking, immutable receipt/activity evidence and one-way state transitions. Migration 0036 adds immutable From/To UOM evidence and stateful physical/composition reservations without binding history to consumable composition rows. Domain tests reject malformed quantities and unsupported demands before persistence. PostgreSQL integration proves clean receipt-to-put-away, capacity/ranking, FEFO across two lots and bins, one winner under concurrent pick pressure, exact BOX→BOX picking, explicit rejection of implicit BOX→EA, authorized automatic breakbulk, registration conservation, cancellation release and SQLSTATE 55000 immutability. HTTP tests prove strict JSON, organization scope and executable registration routing.
