# Portable serial inventory — Microsoft Business Central derivation record

## Authority and license

The portable module is an **ADAPTED** implementation, not verbatim Microsoft Go code and not a Business Central runtime. Its authority is Microsoft `BCApps` at GitHub-verified commit `2eae56d704a1fd035d104f333602aea7091b7749` (tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`), licensed under MIT by the root `License.txt`.

| Official Microsoft source | Git blob | Bytes | SHA-256 |
|---|---:|---:|---|
| `src/Layers/W1/BaseApp/Inventory/Availability/AvailabletoPromise.Codeunit.al` | `440b6ed4dfc6c9297dd3115d9fee682aab6aaf17` | 36,958 | `76d582db57a228559e2458920e11f70776597af25fd976863dbbccac1079d29f` |
| `src/Layers/W1/BaseApp/Inventory/Tracking/ReservationEngineMgt.Codeunit.al` | `b4943be474b5c2023ecc1c5590ae0ce4e21e6158` | 57,538 | `7b0b4c4d1f58b74587e7582a763208342db6376684371f0e6b6448a5f78358db` |
| `src/Layers/W1/BaseApp/Inventory/Tracking/CreateReservEntry.Codeunit.al` | `9fe395c7ac00b36ea628eaac98b4af2b839a9529` | 64,760 | `363c0aefc1350867cf3891ae34ff46967905f2493774492af7a481da35836998` |
| `src/Layers/W1/BaseApp/Inventory/Transfer/TransferLineReserve.Codeunit.al` | `0fd198a2c51219ee76eacee26c9c76bde5b58483` | 65,044 | `81f33a5d8cfe079420bbfde13f4fb8b589422fc4e9ea010d4ab409b301974a9e` |
| `src/Layers/W1/BaseApp/Warehouse/Availability/WarehouseAvailabilityMgt.Codeunit.al` | `946f3604a69e9bbaeab514a2d1fab42f2dc01e6f` | 62,348 | `1744161a3d50d013c0ca07f1f4e2e0dd657d81b4fab1f733cf9867325a04c096` |
| `src/Layers/W1/Tests/SCM-Reservation/SCMTransferReservation.Codeunit.al` | `ec88292b6cfe26db06e19f93598fc3cf48db121d` | 161,402 | `b41c605cd18bca08bdafae7680f04434c14853667be2fb44f05f7e390a846607` |

## Narrow rules transferred

1. ATP follows the Microsoft shape: available inventory plus scheduled receipts minus gross requirement, while avoiding subtraction of demand already bound by reservation. The portable function applies this only to serialized sellable units, due inbound transfers and unallocated placed/confirmed order lines.
2. A reservation binds one supply unit to one demand and cannot be silently duplicated. Cancellation is versioned; customer-order allocations are cancellation-disallowed through the generic endpoint.
3. Transfer reservation is direction-aware. Draft creation reserves outbound units, shipment moves them to `in-transit`, only the destination organization can receive them, and receipt changes ownership before returning the unit to `available`.
4. Every mutation is tenant- and organization-scoped, optimistic-versioned, transactional with its outbox event and tested for conflict/replay boundaries.

## Explicit non-claims

This serial-inventory boundary does not itself implement full Business Central planning, CTP, warehouse picks/put-aways, replenishment, calendars, substitutions, costing adjustment or AL extension events. The composed pack's separate `MICROSOFT_BC_BULK_DERIVATION.md` governs the admitted lot/bin/FIFO-specific boundary without changing these serial rules. Neither boundary claims Microsoft authored or supports the Go/PostgreSQL translation.

## Executed regression mapping

- concurrent attempts to reserve one stock unit produce exactly one winner;
- customer allocation creates the same canonical reservation row;
- ATP reports inbound stock only after shipment and within the requested horizon;
- source organization cannot receive its own outbound transfer;
- receipt changes organization and permits a later reverse transfer;
- all migrations 0001–0025, full Go tests, vet and both API builds execute with the pinned toolchains.
