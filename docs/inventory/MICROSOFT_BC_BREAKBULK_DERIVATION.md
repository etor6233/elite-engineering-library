# Microsoft Business Central handling-unit composition derivation

This portable implementation is `ADAPTED`. Its Go, SQL and tests are not Microsoft-authored code and Microsoft does not support this adaptation. The governing source is the MIT-licensed `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749` (tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`) plus the linked Microsoft Learn product guidance.

## Fixed official inputs

- [Set up units of measure](https://learn.microsoft.com/en-gb/dynamics365/business-central/inventory-how-setup-units-of-measure)
- [Enable automatic breaking bulk](https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-enable-automatic-breaking-bulk-with-directed-put-away-and-pick)
- [Warehouse management design details](https://learn.microsoft.com/en-ca/dynamics365/business-central/design-details-warehouse-management)

| BCApps path | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Inventory/Item/ItemUnitofMeasure.Table.al` | 20,689 | `ab785d08d4bbc9160c2e068b5660ada732f752d35ae7c94a6677ef0fb79b2009` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/CreatePutaway.Codeunit.al` | 77,370 | `b004c0663d237585d8b57a06356431acb2302cdfbbd30f7b82154d282f325b2b` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseChangeUnitofMeasure.Report.al` | 8,008 | `d1fca83ae830a8b615c558e15cb6ec3734560614ba66e86d801f3e9794b06186` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/Replenishment.Codeunit.al` | 19,064 | `e6d5e3873b4828db8dca49755670bd9815faa7ad1b0396189014af57380b7c4e` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/WhseCreatePick.Codeunit.al` | 1,584 | `43ffaded138adb24d0cfc927c276d63ef3dfb2760f83f518bcbc8bbc217163bc` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMMovement.Codeunit.al` | 89,432 | `0e63996a0816d9379bdf76acc5712990f10624de81d875f2e499e6944b2955ed` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWMSItemUnitofMeasure.Codeunit.al` | 125,924 | `8faec658fad625a12c83b251bab88553319bc1cf564e7671104f6a19cadf4713` |

## Official invariants transferred

The fixed replenishment owner searches the requested UOM first. Its breakbulk fallback is explicit, considers a larger package factor, retains base quantity as the authority and writes source and destination UOM/factor evidence. The fixed Microsoft tests aggregate Take and Place activity lines in base quantity, verify warehouse entries by UOM in base quantity, cover lot tracking, expose breakbulk lines during pick registration and also exercise gathering from a smaller UOM into a larger UOM.

The portable contract therefore admits only these narrower behaviors:

1. `inventory.bulk_balance.quantity` is the sole physical quantity in base UOM. Packaging rows are a composition, never a second inventory authority.
2. Every composition row carries an immutable UOM factor snapshot and exact base quantity. A deferred database invariant requires the sum of composition base quantities to equal the physical balance at commit.
3. `breakbulk` converts a larger factor to a smaller factor; `gather` converts a smaller factor to a larger factor. Both require exact source and target precision and preserve base quantity without rounding.
4. Each accepted command is serializable and idempotent by tenant/request. Divergent replay conflicts. Source composition is conditionally decremented so concurrent commands cannot consume the same package twice.
5. The immutable conversion header and its exactly two lines record Take from the source UOM and Place into the target UOM with equal base quantity. The outbox event commits in the same transaction.
6. Ordinary balance increases enter the base composition. Ordinary decreases consume only base composition. If stock is packaged only in a larger UOM, implicit unpacking is rejected and the caller must execute an explicit conversion first.
7. The upgrade migrates a legacy article lacking UOM evidence to factor `1` at the existing six-decimal physical storage precision, records the exact backfill set and reverses only that set on migration rollback.
8. A warehouse receipt accepts either an unambiguous base quantity or an unambiguous handling UOM/quantity pair. PostgreSQL resolves the immutable factor, verifies handling and base precision exactly, and keeps unit cost denominated in base UOM.
9. Put-away planning retains the received handling UOM only when every activity line is an exact multiple. Otherwise it fails closed unless the caller explicitly authorizes automatic breakbulk; the authorized transaction records equal-base Take/Place evidence linked to the posted receipt.
10. Every open put-away reserves both the singular physical balance and its matching composition. Registration releases/consumes both reservations, preserves the handling UOM across bins when possible, and commits movement entries in the same transaction. Cancellation releases both reservations without deleting receipt evidence or stock.
11. A manual conversion may consume only physical and composition quantities not reserved by an open activity. Transfer receipts use the same base-UOM reservation contract, so they do not bypass warehouse registration.

## Executable proof

PostgreSQL 18.6 from migrations `0001` through `0035` retains the V168 proof and adds an end-to-end V169 corpus. Two received BOX are preserved as one BOX on each of two 12-EA put-away lines; physical `24` and composition `2 BOX` are both fully reserved, manual conversion conflicts, registration leaves one unreserved BOX in each target and conservation remains exact. A one-BOX receipt split into two 6-EA targets rolls back without residue when authorization is absent; with authorization it commits two base-UOM lines, immutable Take/Place, receipt source binding and one outbox event. Cancelling an open one-BOX put-away leaves `12 EA = 1 BOX`, releases both reservations and permits the subsequent explicit conversion.

The full PostgreSQL suite also proves transfer-receipt put-away under the same reservation contract. Go tests, vet and build pass. Migration proof creates legacy receipt/activity rows under `0034`, applies `0035`, verifies exact `EA / 12 / factor 1` snapshots and zero packaging reservation, rolls back while preserving physical/receipt/activity quantities, reapplies and reproduces the same snapshots. A post-cycle V169 corpus passes again.

## Explicit non-claims

V169 proves automatic handling-UOM selection only for warehouse receipt → put-away → registration/cancellation and base-UOM transfer-receipt put-away. It does not claim packaging-aware pick or replenishment selection, sales/service/production demand parity, handheld/barcode execution, cubage/weight constraints, external WMS reconciliation or Business Central feature parity. Those lanes remain fail-closed until independently connected and tested. No target production readiness is inferred.
