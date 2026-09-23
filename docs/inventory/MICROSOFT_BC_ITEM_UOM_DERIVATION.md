# Microsoft Business Central item UOM derivation

## Fixed authority

- Repository: `microsoft/BCApps`
- Commit: `2eae56d704a1fd035d104f333602aea7091b7749`
- Tree: `f9846fb1254c6311985c6131bb2af11c7f169e1b`
- License: MIT
- Product guidance:
  - <https://learn.microsoft.com/en-gb/dynamics365/business-central/inventory-how-setup-units-of-measure>
  - <https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-enable-automatic-breaking-bulk-with-directed-put-away-and-pick>
  - <https://learn.microsoft.com/en-ca/dynamics365/business-central/design-details-warehouse-management>

| Official file | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Inventory/Item/ItemUnitofMeasure.Table.al` | 20,689 | `ab785d08d4bbc9160c2e068b5660ada732f752d35ae7c94a6677ef0fb79b2009` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/CreatePutaway.Codeunit.al` | 77,370 | `b004c0663d237585d8b57a06356431acb2302cdfbbd30f7b82154d282f325b2b` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseChangeUnitofMeasure.Report.al` | 8,008 | `d1fca83ae830a8b615c558e15cb6ec3734560614ba66e86d801f3e9794b06186` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/Replenishment.Codeunit.al` | 19,064 | `e6d5e3873b4828db8dca49755670bd9815faa7ad1b0396189014af57380b7c4e` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/WhseCreatePick.Codeunit.al` | 1,584 | `43ffaded138adb24d0cfc927c276d63ef3dfb2760f83f518bcbc8bbc217163bc` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMMovement.Codeunit.al` | 89,432 | `0e63996a0816d9379bdf76acc5712990f10624de81d875f2e499e6944b2955ed` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWMSItemUnitofMeasure.Codeunit.al` | 125,924 | `8faec658fad625a12c83b251bab88553319bc1cf564e7671104f6a19cadf4713` |

## Narrow portable contract

This pack does not copy AL into Go and does not claim that Microsoft supports the adaptation. It transfers only the following demonstrated invariants:

1. Inventory authority is expressed in one base UOM. Its factor is exactly `1`.
2. Every alternate factor is positive and must align exactly with the base quantity rounding precision.
3. A conversion is `requested quantity × quantity per UOM`; it is admitted only when the result has no residual against base precision.
4. UOM identity, factor and precision are immutable in the portable owner. This is stricter than changing only when no warehouse or open-document evidence exists, and prevents historical reinterpretation.
5. Duplicate codes, unknown codes, misaligned factors and residual quantities fail closed before any inventory effect.

Microsoft BCApps further demonstrates same-UOM-first selection, explicitly enabled breakbulk, larger-package fallback, smaller-UOM gathering and paired warehouse Take/Place evidence. V168 adds an authoritative composition plus explicit breakbulk/gather with paired evidence; automatic selection inside receipt, movement, pick, registration and cancellation remains unclaimed until connected and tested end to end.

## Reconstruction and verification

Migration `0033` owns the immutable item-UOM table and one-base constraint. Item creation inserts its base record in the same transaction as the item and outbox event. The service exposes strict alternate configuration and exact conversion. PostgreSQL integration proves factor `12`, `2.5 BOX = 30 EA`, rejection of `0.1 BOX` for an integer base, rejection of factor `2.5`, duplicate rejection and SQLSTATE `55000` on mutation. Migration `0034` supplies the separately documented physical-composition owner and reversible legacy backfill. HTTP tests prove strict authenticated configuration and conversion surfaces.
