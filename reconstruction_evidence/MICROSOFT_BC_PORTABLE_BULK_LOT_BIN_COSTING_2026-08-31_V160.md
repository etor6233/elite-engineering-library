# Microsoft BC portable bulk, lot, bin and costing evidence — V160

## Scope and truth of provenance

V160 extends the existing inventory owner; it does not create a competing inventory service. The new portable code is an explicit `ADAPTED` implementation governed by Microsoft BCApps invariants. It is not verbatim Microsoft Go/PostgreSQL code and is not a Business Central replacement.

Authority: `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`, root MIT license. Inventory subtree `7b2ce465bcc461bd7f0c5c05f998327cfe0133c3`; Warehouse subtree `5a61730670ca182300c1030a4b8643d64f4a2468`.

| Exact official file under `src/Layers/W1/BaseApp` | Git blob | Bytes | SHA-256 |
|---|---|---:|---|
| `Inventory/Tracking/LotNoInformation.Table.al` | `9d22fd3615ebe6cd5851a4919f3b993b87358134` | 8,728 | `fbad735cc2e5fdd4c51d2665fb61391f70e06f8be39be18499860384664f51dc` |
| `Inventory/Tracking/TrackingSpecification.Table.al` | `c423a0a08e2de9a5f6d705aeeca25f48a1859875` | 62,963 | `266092d86843f957ff447ce38e06845581a9a33ab2193d0f7899fe43da48e85c` |
| `Warehouse/Structure/Bin.Table.al` | `f05964937f99dafb8e26f27433989b4c7a75fae1` | 29,194 | `d44007a5a306974b6fcc489428bb731cbf856c12b459e47edab810dd1eda043e` |
| `Warehouse/Structure/BinContent.Table.al` | `dbc4bafd207a306b7e271f767ad53fe31d07c080` | 80,930 | `841f172d7e92a0c581624f7e3b466e6e2ca59807d8d32201d82e1d38e3baa904` |
| `Warehouse/Availability/WarehouseAvailabilityMgt.Codeunit.al` | `946f3604a69e9bbaeab514a2d1fab42f2dc01e6f` | 62,348 | `1744161a3d50d013c0ca07f1f4e2e0dd657d81b4fab1f733cf9867325a04c096` |
| `Inventory/Item/CostingMethod.Enum.al` | `a51e18654ef3348189e59d8757ebdc946f958af4` | 720 | `8f2bac8308ad53a0d3e1b43515a336df3b862e99408fe6fc61921ec23c57f780` |
| `Inventory/Costing/AvgCostAdjmtEntryPoint.Table.al` | `ca8b892ebf0ef0b53d638a82f16df527d217bcc8` | 12,333 | `ab1988b35fde043f575dae032e5beff5169cc486fb0b25f1d9c31325ebb99450` |
| `Inventory/Ledger/ItemApplicationEntry.Table.al` | `17dd82953a3615e541a617d0db4ca6da2f685430` | 49,965 | `3df93f7a3ade1c6c73e715a991514f1d8a56fa5565552d1d3921c2bedf84ee16` |

Microsoft Learn authorities were checked on 2026-08-31: item tracking, item tracking in warehouse, warehouse management, inventory costing and average cost. They confirm lot/serial tracking, bin content, availability net of warehouse activity and explicit costing/application behavior.

## Materialized implementation

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.4.0 materializes 24 files. V160 adds nine:

- migration 0026 up/down;
- Microsoft derivation record;
- bulk domain and regression;
- PostgreSQL transactions and integration regression;
- authenticated HTTP boundary and regression.

Six new blocks are `ADAPTED`; three are local `AUTHORED`. The application pack 1.7.0 changes its existing authored composition root to inject one PostgreSQL repository into both serial and bulk services.

The implementation provides tenant/organization-scoped items, lot identity and expiry, physical bins, per-item bin policy, exact `numeric` balances, atomic quantity reservations, availability net of reservations, internal bin movements, immutable inventory entries, inbound cost layers and outbound cost applications. FIFO consumes oldest remaining layers. Specific costing requires the exact receipt entry. Serial vehicles remain in the prior serial owner.

LIFO, Average and Standard are recognized official methods but deliberately rejected by this portable lane; no synthetic cost adjustment was invented.

## Executed gates

From the authoring tree:

- supply pack materialized 24/24; application pack 5/5;
- 26 migrations applied to a clean PostgreSQL 18.6 instance;
- bulk, commerce and serial focal integration tests passed;
- full `go test ./...`, `go vet ./...`, `cmd/api` build and `cmd/electromobility-api` build passed;
- concurrency admitted exactly one of two reservations exceeding one lot's capacity;
- FIFO consumed 5 units at 100 plus 2 at 120 and posted exact cost 740 through two immutable applications;
- specific cost rejected a missing inbound identity and then posted exact 999.99 when supplied;
- entry mutation returned SQLSTATE 55000;
- movement preserved total quantity and created paired zero-cost internal entries;
- test teardown removed its tenant without contaminating the global outbox suite.

From canonical Markdown:

- backend composed 316 files from 29 packs;
- 12 focal source/rebuilt hashes matched, with 316 materialization-record rows;
- clean migrations 0001–0026, full Go tests, vet and both builds passed again from the reconstructed tree.

PostgreSQL was stopped after verification.

## Failures retained

`LIB-FAIL-1504` through `LIB-FAIL-1512` preserve wrong GitHub path extrapolation, obsolete materializer parameter/script names, a guessed HTTP path, a non-UUID fixture, outbox test isolation, pgx multi-command teardown and two Windows/search-expression failures. No failed result was promoted. Accepted corrections were rerun from clean composition and database state.

## Conditions that remain open

V160 does not implement the entire Business Central costing/planning or warehouse engine. Periodic average adjustment, standard variances, LIFO policy, expected versus actual cost, item charges/landed cost, G/L posting, closed inventory periods, picks/put-aways, cross-organization bulk transfer, FEFO orchestration, replenishment and CTP remain conditioned. Production also requires approved accounting/currency/rounding rules, real warehouse roles, load/security/recovery tests, ERP/WMS reconciliation and business acceptance.
