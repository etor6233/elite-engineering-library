# Microsoft BC operational warehouse and FEFO evidence — V161

## Scope and provenance truth

V161 extends the existing inventory owner. It does not introduce a second WMS or inventory database. The portable code is explicitly `ADAPTED` from Microsoft BCApps invariants; it is not verbatim Microsoft Go/PostgreSQL code and is not presented as a Business Central replacement.

Authority: `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`, repository MIT license.

| Exact official path under `src/Layers/W1/BaseApp` | Git blob | Bytes | SHA-256 |
|---|---|---:|---|
| `Warehouse/Document/WhsePostReceipt.Codeunit.al` | `dbd6327a935a69c564d112203eaf42ee845ba393` | 105803 | `9aae2fa85c20e12e70b88b3214fd4b4defcb13b6e767c225251d413eda626989` |
| `Warehouse/Activity/CreatePutaway.Codeunit.al` | `df3011a0a74ace99e3d1788996ba750a7732a1e6` | 77370 | `b004c0663d237585d8b57a06356431acb2302cdfbbd30f7b82154d282f325b2b` |
| `Warehouse/Activity/CreatePick.Codeunit.al` | `07be15a7effeeda1ac749ef033d20d817bb5b218` | 285538 | `6ca4acba3eebb46c719267ed788d8efc6ce097005db0d1a560489797e03b35e5` |
| `Warehouse/Tracking/WhseItemTrackingFEFO.Codeunit.al` | `9568f0f687384a65ac6dc4135273cc44b7256c4e` | 17061 | `24f0b95de17d6d2a042b97fe2cefb0ded952df8db2f0ebfefdb75c47aca09965` |
| `Warehouse/Activity/WarehouseActivityHeader.Table.al` | `4e3e2cf15ed00cfdc627898a1fafb54a8c8eca29` | 47747 | `34f6f3b633aed6c77c1cf286e9528d449fb255dcc219531946b43b59ce862030` |
| `Warehouse/Activity/WarehouseActivityLine.Table.al` | `83d6dc89fa2b088e05d9ac224fe0ab9e3c1921b9` | 181260 | `f4ab727d452c158b1f397e83708150f04b9fac22262e378bbca0bd9a8ed6836b` |

Microsoft Learn authorities checked on 2026-08-31: warehouse management, inbound warehouse flow, advanced put-away and warehouse shipment picking. They state that receipt precedes put-away; Take/Place activities are separate; PUT AWAY, RECEIVE, SHIP and QC are not generic pick sources; bin ranking guides placement/pick; FEFO excludes unavailable stock and orders tracked inventory by expiration with deterministic tie-breaking.

## Materialized implementation

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.5.0 materializes 33 files. V161 adds nine files:

- migration 0027 up/down;
- exact Microsoft derivation record;
- warehouse domain and validation regression;
- PostgreSQL transaction boundary and integration regression;
- authenticated HTTP boundary and regression.

Six new blocks are `ADAPTED`; three are local `AUTHORED`. Three existing supply files add bin ranking/QC and HTTP composition. `GO-ELECTROMOBILITY-APPLICATION` 1.8.0 changes its authored composition root so the same PostgreSQL inventory repository serves serial, bulk and operational warehouse services.

The implementation provides posted immutable receipt evidence, inventory receipt into RECEIVE, open put-away planning across eligible bins using default/capacity/ranking and already-planned quantity, atomic Take/Place registration, PICK/PUTPICK-only selection, FEFO across lots, reservation creation per activity line, move-to-SHIP registration, open-pick cancellation, optimistic versions, tenant/organization authorization, unique request/source identities, append-only entries and transactional outbox.

## Executed gates

Authoring composition:

- Go official 1.26.7 windows/amd64;
- PostgreSQL 18.6 clean cluster;
- migrations 0001–0027 applied in order;
- receipt placed a six-unit lot in high-rank PUTPICK and a four-unit lot in lower-rank PUTPICK after capacity was exhausted;
- registration created paired movement evidence and made exactly ten units pickable;
- FEFO selected four early-expiry units before three later-expiry units despite lower bin ranking;
- two simultaneous three-unit picks against three remaining units produced exactly one winner;
- pick registration moved quantity and reservation to SHIP; cancellation released all outstanding reservation;
- receipt mutation returned SQLSTATE 55000;
- strict HTTP/organization tests, full `go test ./...`, `go vet ./...`, `cmd/api` and `cmd/electromobility-api` builds passed.

Canonical Markdown reconstruction:

- supply materialized 33/33 and application 5/5;
- 13 focal source/rebuilt hashes matched;
- backend composed 325 files from 29 packs;
- a new database applied migrations 0001–0027;
- full tests, vet and both builds passed again from the reconstructed tree.

## Failure learning

`LIB-FAIL-1513` through `LIB-FAIL-1521` retain the incorrect compositor parameter, a rejected compound cleanup, missing Go PATH, an ambiguous SQL join masked as conflict, non-UUID test events, a fixture that wrongly expected picking from PUT AWAY, an organization-scope mismatch, a PowerShell array binding error and a false multi-script success from null `LASTEXITCODE`. Each accepted correction was rerun; no failed result was promoted.

## Conditions retained

V161 is a narrow operational warehouse lane, not full Business Central WMS. It does not implement partial line handling, breakbulk/UOM conversion, weight/cubage, warehouse classes, cross-docking, replenishment, barcode/device workflows, serial FEFO, production/assembly/project picks, cross-organization bulk transfer, planning/CTP or the complete extension-event surface. Production still requires approved warehouse roles/processes, physical scan reconciliation, load/security/recovery testing, ERP/WMS cutover and business acceptance.
