# Microsoft BC packaging-aware pick and replenishment — reconstruction evidence V171

## Result and claim boundary

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.14.0 materializes 78 files and the enterprise backend composes 29 packs / 370 implementation files. V171 adds three `ADAPTED` SQL/Go files and updates twelve existing code, test and derivation owners. None is represented as Microsoft-authored Go or SQL. The portable implementation is governed by the fixed Microsoft authorities below and proves only the packaging-aware picking and bin-replenishment behavior in this record.

V171 does not claim the complete Microsoft Business Central WMS/ERP, target production readiness, operator/barcode UX, sales/service/production demand policy, planning/CTP, working-day calendars, accounting/landed cost, fiscal homologation or any live provider integration.

## Fixed official authority

- Repository: `microsoft/BCApps`
- Commit: `2eae56d704a1fd035d104f333602aea7091b7749`
- Tree: `f9846fb1254c6311985c6131bb2af11c7f169e1b`
- License: MIT
- Microsoft Learn:
  - <https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-enable-automatic-breaking-bulk-with-directed-put-away-and-pick>
  - <https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-how-to-calculate-bin-replenishment>
  - <https://learn.microsoft.com/en-us/dynamics365/business-central/warehouse-how-to-pick-items-for-warehouse-shipment>
  - <https://learn.microsoft.com/en-gb/dynamics365/business-central/inventory-how-setup-units-of-measure>

| Exact official owner/test | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Warehouse/Activity/CreatePick.Codeunit.al` | 285,538 | `6ca4acba3eebb46c719267ed788d8efc6ce097005db0d1a560489797e03b35e5` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/Replenishment.Codeunit.al` | 19,064 | `e6d5e3873b4828db8dca49755670bd9815faa7ad1b0396189014af57380b7c4e` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/WhseCreatePick.Codeunit.al` | 1,584 | `43ffaded138adb24d0cfc927c276d63ef3dfb2760f83f518bcbc8bbc217163bc` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMMovement.Codeunit.al` | 89,432 | `0e63996a0816d9379bdf76acc5712990f10624de81d875f2e499e6944b2955ed` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWMSItemUnitofMeasure.Codeunit.al` | 125,924 | `8faec658fad625a12c83b251bab88553319bc1cf564e7671104f6a19cadf4713` |

Microsoft's fixed replenishment owner selects the target UOM first and considers a larger source UOM only through its explicit `AllowBreakbulk` path. The fixed tests exercise unit-of-measure breakbulk, tracked inventory and activity registration. Microsoft remains authority for its AL implementation; Elite owns and must verify the portable adaptation.

## Executable behavior

- Pick and replenishment commands accept an explicit target/handling UOM and an independent `allow_breakbulk` decision.
- Planning prefers exact target-UOM composition. A larger source UOM is eligible only when breakbulk is explicitly allowed; otherwise the transaction fails without activity or reservation residue.
- Every line persists From/To UOM code, quantity and factor. An immutable stateful reservation binds the activity line to the physical balance and source-composition snapshot without preventing legitimate consumption of the mutable composition row.
- Open work reserves both base quantity and exact source packages. A concurrent manual conversion cannot consume the reserved package.
- Registration consumes both reservations, writes automatic breakbulk Take/Place plus source activity provenance, moves requested target composition and retains the package remainder exactly. Cancellation releases both reservation layers.
- Exact BOX→BOX work does not convert or silently flatten packaging. BOX→EA work requires authorization. Physical base quantity and packaging composition remain conserved.

## Executable proof

Pinned runtimes used on Windows:

- Go `1.26.7` from the admitted portable toolchain;
- PostgreSQL `18.6` archive, 343,808,005 bytes, SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`;
- `pgx` `5.10.0` from the locked backend module.

The following passed on `2026-09-01`:

1. isolated V169 supply reconstruction followed by V171 owner editing and Go formatting;
2. overlay compilation only after composing the complete enterprise backend module;
3. migrations `0001` through `0036` from a new PostgreSQL database;
4. all PostgreSQL integration tests on a separate clean gate database;
5. focal exact BOX→BOX pick, implicit BOX→EA rejection, authorized pick cancellation/registration, manual conversion exclusion and authorized replenishment;
6. migration `0036` down/up followed by the focal test;
7. updater regression test, then atomic update of fifteen materialization blocks;
8. Markdown-only supply reconstruction: 78 rebuilt files and exact SHA equality for all fifteen changed/new owners;
9. Markdown-only enterprise backend composition: 29 packs / 370 files;
10. canonical fresh-database migrations `0001`–`0036`, complete PostgreSQL suite, full uncached `go test -count=1 ./...`, `go vet ./...` and `go build ./...`;
11. canonical migration `0036` down/up and focal repetition.
12. root `VERIFY_LIBRARY_PASS`: 94 packs, 1,087 materialized files, 528 Markdown files and enterprise backend 370;
13. executable `Audit`: exit 0 with `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=94 upstream_sources=121`, including 62 readiness regressions and the governed browser/document/runtime component checks.

The final canonical corpus proves one exact package reaches SHIP as one BOX; an unauthorized six-EA pick or replenishment from one BOX fails atomically; an authorized pick reserves six base units and one BOX, blocks competing conversion, cancels cleanly, then registers six EA at source plus six EA at SHIP with immutable conversion provenance; replenishment repeats BOX→EA into its min/max target.

## Failure learning

V171 records local failures 1668 through 1689. Truncated reads, path/name assumptions, wildcard recurrence, partial-module compilation, an invalid FK to consumable state, cleanup order, fixture contracts, ambiguous SQL evidence, database reuse, PowerShell array binding and stale aggregated provenance counts were each invalidated before canonical promotion. No failed output is counted as PASS. Library memory at this closure is 1,689 local failures plus 202 upstream conditions: 1,891 unique IDs.

## Artifact identity

- Pack: `implementation_packs/GO_SUPPLY_FACTORY_INVENTORY_API.md`
- Version: `0.14.0`
- Bytes before root-gate evidence linkage: `584517`
- SHA-256 before root-gate evidence linkage: `1ae2d74c0fa07a16ddc472ef4d9d5e777961addbadb2e88e0edc6feda2a12df3`
- Materialized files: `78`
- Enterprise backend composition: `29 packs / 370 files`

Production remains conditioned on project-specific infrastructure, providers, policy, load/security/recovery evidence and acceptance. This proof closes the packaging-awareness gap only for the admitted pick and bin-replenishment journeys.
