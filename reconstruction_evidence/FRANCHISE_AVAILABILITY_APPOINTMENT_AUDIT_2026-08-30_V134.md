# Franchise Availability and Appointment Audit — V134

Date: 2026-08-30  
Decision: `REBUILD_VERIFIED / CONDITIONED`

## Claim boundary

`GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.5.0` adds explicit organization/resource working intervals, unavailable intervals with reason, appointment-slot and assignment enforcement, immutable appointment transition evidence and customer-owned future cancellation. It does not invent weekly recurrence, timezone, labor law, absence entitlement or cancellation policy. Those inputs remain governed project configuration.

All six new files and the changed local Go/SQL files are `AUTHORED`. They are governed by the official sources below; they are not represented as Microsoft or PostgreSQL source code.

## Official authority fixed

- Microsoft Business Central employee absence: <https://learn.microsoft.com/en-us/dynamics365/business-central/hr-how-manage-absence>
- Microsoft Business Central work center/shop calendars: <https://learn.microsoft.com/en-us/dynamics365/business-central/production-how-to-create-work-center-calendars>
- Microsoft Business Central base calendars and customized nonworking days: <https://learn.microsoft.com/en-gb/dynamics365/business-central/across-how-to-assign-base-calendars>
- Microsoft Employee Absence table contract: <https://learn.microsoft.com/en-us/dynamics365/business-central/application/base-application/table/microsoft.humanresources.absence.employee-absence>
- PostgreSQL 18 explicit/advisory locking: <https://www.postgresql.org/docs/18/explicit-locking.html>

The implementation uses explicit effective intervals because these authorities separate calendars, working/nonworking time and absence records. A concrete project must provide its approved recurrence/timezone and workforce rules before generating those intervals.

## Exact reconstructed artifact

- Journey pack SHA-256: `f56191e9cbce1a9a821b17031d56901d3828d1e9c775008cb5bb22cba71f24e7`.
- Application composition-root SHA-256 after compatibility bump: `9eadfd6348537a67fdf63d1a9a9a5e51abf7bbd8054e2f941db824311714bb9a`.
- Journey pack materialized 24 files into an empty directory; all 24 SHA-256 values matched the tested authoring tree.
- Enterprise backend plan composed 24 packs and 238 files into an empty directory, plus one `MATERIALIZATION_RECORD.md`.

## Gates executed

1. Go package tests for domain, HTTP and PostgreSQL passed after contract propagation.
2. A new PostgreSQL 18.6 cluster applied migrations 0001–0015 with `ON_ERROR_STOP=1`.
3. SQL test `0015_franchise_availability_and_appointment_audit.test.sql` passed.
4. Real PostgreSQL integration proved organization working time, slot enforcement, resource working time and skills, exactly-one concurrent capacity result, absence rejection over an active appointment, transition actor/reason persistence, customer identity isolation and customer cancellation.
5. Full `go test ./... -count=1`, `go vet ./...` and builds for `cmd/api`, `cmd/electromobility-api`, `cmd/arca-fiscal-worker` and `cmd/arca-parameter-worker` passed.
6. The same migrations, SQL test, full Go suite, vet and builds passed again from the 238-file Markdown composition, not the authoring tree.
7. Both PostgreSQL runs were stopped explicitly after verification.

## Negative evidence retained

- Slot outside organization working time: rejected by PostgreSQL.
- Resource without matching skill, working interval or free interval: rejected.
- Unavailability overlapping a requested/confirmed appointment: rejected.
- Working interval cancellation supporting an active appointment: rejected.
- Staff cancellation/no-show without a governed reason: rejected by the service/schema contract.
- Customer cancellation for another customer, stale version, past or terminal appointment: rejected.
- Transition rows cannot be updated or deleted in normal operation.

## Production condition

This proves the reusable code path and its reconstruction. It does not prove a concrete franchise is production-ready. The project still must connect its authoritative workforce/calendar source, approve timezone/recurrence and policy, and demonstrate real IdP, edge, load, security, recovery, deployment and business acceptance.
