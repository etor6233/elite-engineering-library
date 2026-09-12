# Microsoft BC portable serial inventory — reconstruction V159

## Result

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.3.0 now materializes fifteen files. The enterprise backend composes 29 packs / 307 files and contains one inventory authority for serialized vehicles: reservations, sellable ATP and organization-to-organization transfers. Commerce allocation writes the same canonical reservation table; it does not retain a parallel allocation ledger.

This is an `ADAPTED` Go/PostgreSQL implementation governed by Microsoft source. It is not represented as verbatim Microsoft code or as a replacement for Business Central.

## Exact Microsoft authority

- repository: `microsoft/BCApps`;
- GitHub-verified commit: `2eae56d704a1fd035d104f333602aea7091b7749`;
- tree: `f9846fb1254c6311985c6131bb2af11c7f169e1b`;
- root license: MIT;
- exact files and SHA-256 values are recorded in materialized `docs/inventory/MICROSOFT_BC_DERIVATION.md`;
- source surfaces: `AvailabletoPromise`, `ReservationEngineMgt`, `CreateReservEntry`, `TransferLineReserve`, `WarehouseAvailabilityMgt` and `SCMTransferReservation`.

The transferred rules are narrow: ATP equals available serialized inventory plus due inbound transfers minus unallocated placed demand; a reservation uniquely binds supply to demand; shipment changes reserved units to in-transit; only the destination receives; receipt changes organization and returns the unit to availability.

## Materialized implementation

- migration `0025_serial_inventory_reservation_transfer` with reservation/transfer tables, partial uniqueness and `inventory.serial_atp`;
- domain port and validation under `internal/inventorycontrol`;
- PostgreSQL transactional adapter with outbox events;
- authenticated, organization-scoped HTTP routes for ATP, reservations and transfers;
- commerce allocation inserts the same cancellation-disallowed customer-order reservation;
- composition root wires the new module without adding a second inventory source of truth.

## Executed gates

- isolated official Go `1.26.7 windows/amd64`;
- `gofmt -d`: empty;
- full `go test -count=1 ./...`: PASS from the 29/307 Markdown composition;
- `go vet ./...`: PASS;
- builds `cmd/api` and `cmd/electromobility-api`: PASS;
- PostgreSQL `18.6` with page checksums, clean database and migrations 0001–0025: PASS;
- `TestCommercePriceOrderAllocationPaymentFlow`: PASS and canonical reservation row proven;
- `TestInventoryReservationATPTransferAndConcurrency`: PASS, including exactly one concurrent winner, release, ATP during transit, wrong receiver rejection, receipt, reverse transfer and cancellation;
- PostgreSQL stopped after each attempt; no production account or provider was touched.

## Failures retained

`LIB-FAIL-1492` through `LIB-FAIL-1503` record isolated-pack assumptions, Markdown JSON parsing, GitHub tree truncation, array parsing, toolchain discovery timeout, patch context, interface arity, pgx type identity, serial fixture constraints, Windows `pg_ctl` pipe inheritance, block newline preservation and a policy-rejected cleanup command. Every accepted fix was repeated from a clean reconstruction; the rejected cleanup mutated nothing and was replaced by an allowlisted, containment-checked procedure.

## Limits

The portable pack is deliberately serialized-stock scope. It does not claim full Microsoft planning/CTP, bulk quantities, lots, warehouse bins/picks, replenishment, calendars, substitution, landed cost or costing. Production remains blocked until the concrete project proves load, operational expiry/release policy, ERP/WMS ownership and reconciliation, identity, recovery, security, deploy/rollback and business acceptance.
