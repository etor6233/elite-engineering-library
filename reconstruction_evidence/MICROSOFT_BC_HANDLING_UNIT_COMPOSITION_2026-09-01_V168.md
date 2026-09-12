# Microsoft BC handling-unit composition — reconstruction evidence V168

## Result and claim boundary

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.12.0 materializes 72 files and the enterprise backend composes 29 packs / 364 implementation files. V168 adds five `ADAPTED` SQL/Go/test files and one local `AUTHORED` derivation file; it also updates six existing owners. No Go or SQL file is represented as Microsoft-authored code. This evidence proves the explicit handling-unit composition primitive described below. It does not claim automatic Business Central breakbulk across warehouse activities, complete WMS equivalence or target production readiness.

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

The portable implementation transfers same-UOM/base-quantity authority, explicit larger-to-smaller breakbulk, smaller-to-larger gathering and paired Take/Place base-quantity evidence. The exact adaptation and non-claims are materialized in `docs/inventory/MICROSOFT_BC_BREAKBULK_DERIVATION.md`.

## Executable proof

Pinned runtimes:

- Go `1.26.7` from the already verified exact official archive;
- PostgreSQL `18.6` from the already verified exact official archive;
- `pgx` `5.10.0` from the locked backend module.

The following passed on Windows from a clean PostgreSQL 18.6 cluster on `2026-09-01`:

1. migrations `0001` through `0034` from an empty database;
2. focused PostgreSQL corpus `TestWarehouseHandlingUnitConservationIdempotencyAndConcurrency`;
3. full `go test -count=1 ./...` with `TEST_DATABASE_URL` set to that database;
4. `go vet ./...` and `go build ./...`;
5. migration `0034` down/up followed by the focal corpus;
6. separate legacy corpus: migrations through `0033`, a direct legacy `12.345 KG` physical balance with no UOM row, `0034` up, exact `12.345` composition, `0034` down preserving physical balance and removing only the recorded backfill, then `0034` up reproducing the exact result;
7. compositor self-tests and composition of the current backend profile from Markdown;
8. full Go/PostgreSQL test, vet and build from that Markdown-only 29-pack / 364-file composition;
9. supply pack reconstruction of 72/72 files with zero missing, extra or SHA differences.

The focal corpus proves `24 EA → 2 BOX`, exact idempotent replay, divergent replay conflict, `1 BOX → 12 EA`, and two concurrent commands over the final BOX with exactly one success and one conflict. It also proves physical quantity equals composed base quantity, exactly two immutable Take/Place lines per accepted conversion, atomic outbox evidence, SQLSTATE `55000` on evidence mutation and false composition, fail-closed implicit movement when explicit unpacking is required, and a valid fractional `0.025 KG` base movement.

## Failure learning

V168 recorded `LIB-FAIL-1612` through `LIB-FAIL-1624` before promotion. The material issues were: workdir/cached-test errors, a wrong authorization fixture, an integer-only precision assumption, missing legacy UOM backfill and several path/glob/patch mistakes. Each failed attempt was invalidated. The final proof uses explicit source/destination allowlists, `-count=1`, the composed module, an independent clean database, a separate legacy upgrade database, and Markdown reconstruction. The library memory now governs 1,624 local failures plus 202 upstream conditions with unique IDs.

## Artifact identity

- Pack: `implementation_packs/GO_SUPPLY_FACTORY_INVENTORY_API.md`
- Version: `0.12.0`
- Bytes at evidence creation: `495995`
- SHA-256 at evidence creation: `c04b148db9d3ecaa36d6eba98df0842394affb96ea02cf60b967a1f0772f9f4a`
- Materialized files: `72`
- Enterprise backend composition: `29 packs / 364 files`

Production remains conditioned on project-specific configuration and evidence: automatic breakbulk selection connected to receipt/put-away/replenishment/pick/registration/cancellation, packaging reservations and reversal, representative load, operator/barcode workflows, external WMS reconciliation, accounting/landed-cost policy, deployment, recovery, security and business acceptance.
