# Microsoft BC receipt-to-put-away packaging — reconstruction evidence V169

## Result and claim boundary

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.13.0 materializes 75 files and the enterprise backend composes 29 packs / 367 implementation files. V169 adds three `ADAPTED` SQL/Go/test files and updates nine existing owners. None of those files is represented as Microsoft-authored code. The portable implementation is governed by the exact Microsoft authorities below and proves only the receipt-to-put-away packaging behavior in this record.

V169 does not claim the complete Microsoft Business Central WMS/ERP, packaging-aware pick or replenishment, target production readiness, operator/barcode UX, planning/CTP, accounting/landed cost, fiscal homologation or any external provider integration.

## Fixed official authority

- Repository: `microsoft/BCApps`
- Commit: `2eae56d704a1fd035d104f333602aea7091b7749`
- Tree: `f9846fb1254c6311985c6131bb2af11c7f169e1b`
- License: MIT
- Microsoft Learn:
  - <https://learn.microsoft.com/en-gb/dynamics365/business-central/inventory-how-setup-units-of-measure>
  - <https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-enable-automatic-breaking-bulk-with-directed-put-away-and-pick>
  - <https://learn.microsoft.com/en-ca/dynamics365/business-central/design-details-warehouse-management>

| Exact official owner/test | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Inventory/Item/ItemUnitofMeasure.Table.al` | 20,689 | `ab785d08d4bbc9160c2e068b5660ada732f752d35ae7c94a6677ef0fb79b2009` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/CreatePutaway.Codeunit.al` | 77,370 | `b004c0663d237585d8b57a06356431acb2302cdfbbd30f7b82154d282f325b2b` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseChangeUnitofMeasure.Report.al` | 8,008 | `d1fca83ae830a8b615c558e15cb6ec3734560614ba66e86d801f3e9794b06186` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/Replenishment.Codeunit.al` | 19,064 | `e6d5e3873b4828db8dca49755670bd9815faa7ad1b0396189014af57380b7c4e` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/WhseCreatePick.Codeunit.al` | 1,584 | `43ffaded138adb24d0cfc927c276d63ef3dfb2760f83f518bcbc8bbc217163bc` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMMovement.Codeunit.al` | 89,432 | `0e63996a0816d9379bdf76acc5712990f10624de81d875f2e499e6944b2955ed` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWMSItemUnitofMeasure.Codeunit.al` | 125,924 | `8faec658fad625a12c83b251bab88553319bc1cf564e7671104f6a19cadf4713` |

The exact transfer and the non-claims are materialized in `docs/inventory/MICROSOFT_BC_BREAKBULK_DERIVATION.md`. Microsoft remains the authority for its AL implementation; Elite owns and must verify the portable Go/PostgreSQL adaptation.

## Executable behavior

The command accepts exactly one unambiguous quantity representation: legacy base quantity, or handling UOM plus handling quantity. Unit cost remains per base unit. A put-away preserves its handling UOM only when every placement is an exact package quantity. A non-exact split rolls back unless `allow_breakbulk` is explicit; authorized breakbulk records immutable Take/Place evidence and an outbox event bound to the exact warehouse receipt.

An open put-away reserves both physical base quantity and packaging composition. Registration moves preserved packages to destination bins or consumes the explicitly authorized base composition. Cancellation releases both reservations. A manual conversion cannot consume either kind of reserved stock. Transfer-receipt put-away follows the same base-UOM reservation contract.

## Executable proof

Pinned runtimes used on Windows:

- Go `1.26.7` from the already admitted portable toolchain;
- PostgreSQL `18.6` archive, exactly 343,808,005 bytes, SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`;
- `pgx` `5.10.0` from the locked backend module.

The following passed on `2026-09-01`:

1. migrations `0001` through `0035` from an empty checksummed PostgreSQL 18.6 database;
2. focal `TestWarehouseReceiptPackagingPreservationBreakbulkAndCancellation`;
3. combined warehouse, UOM and transfer suites after correcting the transfer-receipt reservation regression;
4. full `go test -count=1 ./...`, `go vet ./...` and `go build ./...` against PostgreSQL;
5. migration `0035` down/up on an empty schema followed by the focal corpus;
6. a separate legacy corpus under `0034`, then `0035` up/down/up, preserving quantity and reproducing exact base-UOM composition;
7. Markdown-only supply reconstruction: 75 source files, 75 rebuilt files, zero path or hash differences;
8. Markdown-only backend composition: V168 29/364 plus the three V169 owners equals 29/367; full Go/PostgreSQL tests, vet and build passed;
9. root library verification before adding this evidence: 94 packs, 1,084 materialized files and 525 Markdown files;
10. executable Audit: exit 0 with `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=94 upstream_sources=121` and all governed component markers present.
11. root library verification after adding this evidence: `VERIFY_LIBRARY_PASS`, 94 packs, 1,084 materialized files and 526 Markdown files.

The focal corpus proves two `BOX` receipts become two placements of `12 EA = 1 BOX`, both reservations block manual conversion, registration preserves one package per destination, an unauthorized `6+6 EA` split fails atomically, an authorized split emits Take/Place/outbox evidence bound to the receipt, and cancellation releases physical and composition reservations so a later conversion can succeed.

## Failure learning

V169 records `LIB-FAIL-1625` through `LIB-FAIL-1656`. Every failed build, migration, fixture, reconstruction, gate, path assumption and truncated output was invalidated before promotion. The final proof uses uncached Go tests, an isolated checksummed PostgreSQL cluster, a separate legacy corpus, source/rebuilt path and SHA equality, Markdown-only composition, the root verifier and the executable auditor. The library memory at evidence creation governs 1,656 local failures plus 202 upstream conditions: 1,858 unique IDs.

## Artifact identity

- Pack: `implementation_packs/GO_SUPPLY_FACTORY_INVENTORY_API.md`
- Version: `0.13.0`
- Bytes at evidence creation: `542922`
- SHA-256 at evidence creation: `f99fd4111245b88b257d2737a10e822723b2d617a2d8d7cf937a733219e42c6c`
- Materialized files: `75`
- Enterprise backend composition: `29 packs / 367 files`

Production remains conditioned on project-specific targets and evidence. The immediate next warehouse gap is packaging-aware pick and replenishment; it must not be inferred from receipt-to-put-away proof.
