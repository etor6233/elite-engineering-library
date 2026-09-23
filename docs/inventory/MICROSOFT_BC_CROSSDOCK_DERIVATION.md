# Microsoft Business Central warehouse cross-dock derivation

This portable implementation is `ADAPTED`; its Go, SQL and tests are not Microsoft-authored code. It is governed by the MIT-licensed `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749` (tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`) and by Microsoft Learn, [Cross-dock items](https://learn.microsoft.com/en-gb/dynamics365/business-central/warehouse-how-to-cross-dock-items). The upstream code and documentation remain the authority; this file states the smaller contract that the portable implementation proves.

## Fixed official inputs

| BCApps path | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseCrossDockManagement.Codeunit.al` | 47,115 | `85a1eede4b2132f8ef22c853be1bd66a23a29b5176875a0432f500be150609fd` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseCrossDockOpportunity.Table.al` | 18,740 | `778cc8833ea7c28262a9ae695136281c66002d190845785e6f41e100db274ecb` |
| `src/Layers/W1/BaseApp/Warehouse/Document/WarehouseReceiptHeader.Table.al` | 24,887 | `d02b7f1fa0be0847e10e688d98125fdf34e1f157089a99d4632024d16fd2c0fe` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWarehouseV.Codeunit.al` | 306,725 | `faee2bf1ef5bae41211f1374b8e1b7d49e9e7e75f2672e70ed0c5b181136f839` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWarehouseOrders.Codeunit.al` | 238,014 | `79f2e97c278763aa7c140dcc6dffaa53f7e6e6fe2d55100b7f42bad2ccd0a1da` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWarehouseUnitTests.Codeunit.al` | 170,444 | `9af357c3eddaa6f6acaa50623aec43753d212efe7d412e6a948f35c87521ac3f` |

## Admitted executable contract

- an organization/item has at most one enabled policy, an explicit due-date horizon and an explicit fixed, unblocked PICK/PUTPICK bin marked `cross_dock`;
- only released outbound bulk-transfer demand from the same organization and item is eligible, ordered by due date and stable identities;
- demand beyond the horizon, specific-cost transfer lines and demand with an existing open or registered pick are excluded;
- required quantity is outstanding transfer quantity minus every existing cross-dock allocation, whether or not its receipt has already been put away;
- receipt allocation is the minimum of received quantity, remaining demand and available cross-dock-bin capacity; the unallocated receipt quantity follows normal put-away rules;
- an allocation is not pickable until its exact put-away line is registered; generic picks exclude cross-dock bins;
- a matching transfer pick consumes its own available allocation before normal storage, atomically reserves balance and allocation quantity, and records the exact pick-line link;
- pick cancellation releases both reservations; pick registration moves stock to the ship bin and converts reserved allocation quantity to picked quantity;
- organization/item advisory locking and serializable transactions prevent concurrent double allocation and double reservation.

The PostgreSQL corpus proves partial cross-dock plus ordinary storage, existing-opportunity subtraction, due-date exclusion, unavailable-before-put-away rejection, exact transfer prioritization, cancellation/retry, final pick registration, ledger balances and zero leaked tenant data.

## Explicit non-claims

This version does not claim complete Business Central WMS equivalence. Sales, service and production demand; base-calendar working-day calculation; UOM conversion; breakbulk; cubage/weight; warehouse class; handheld/barcode execution; cross-dock time-zone policy; carrier routing; and external-WMS reconciliation remain conditioned until their authoritative mappings and executable tests exist. No production readiness is inferred from pack admission.
