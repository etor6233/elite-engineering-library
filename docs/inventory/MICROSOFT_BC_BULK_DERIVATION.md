# Microsoft BC-derived bulk, lot, bin and costing boundary

This portable implementation is `ADAPTED`; it is not verbatim Microsoft Go code and is not a Business Central runtime.

Authority is Microsoft `BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749`, root MIT license. The audited files are `TrackingSpecification.Table.al`, `LotNoInformation.Table.al`, `Bin.Table.al`, `BinContent.Table.al`, `WarehouseAvailabilityMgt.Codeunit.al`, `CostingMethod.Enum.al`, `AvgCostAdjmtEntryPoint.Table.al` and `ItemApplicationEntry.Table.al` under `src/Layers/W1/BaseApp`.

The portable contract preserves these narrow invariants:

1. lot identity is item-scoped and carries blocked/expiration state;
2. a bin is the smallest physical storage unit, while item/bin policy owns fixed, dedicated and default behavior;
3. sellable availability is physical quantity less reservations and excludes blocked movement, dedicated generic stock and blocked/expired lots;
4. quantity reservations are atomic and cannot exceed balance;
5. every issue is linked to one or more inbound cost entries through immutable application rows;
6. FIFO consumes the oldest remaining inbound layer; specific costing requires the exact inbound entry;
7. movements inside one organization preserve valuation and do not manufacture a cost;
8. serial-tracked vehicles remain in the existing serial owner and cannot enter the bulk API.

Business Central also defines LIFO, Average and Standard costing and a much larger adjustment engine. This portable boundary deliberately rejects those methods. It does not claim periodic average adjustment, manufacturing variance, expected cost, G/L posting, warehouse picks/put-aways, cross-organization bulk transfer or full CTP/planning.

The project must select accounting policy, currency/rounding, landed-cost treatment, closed periods, expiry policy and warehouse roles before production admission.
