# Go Enterprise Query API — reconstruction evidence V1

Date: 2026-08-24  
Status: `REBUILD_VERIFIED / CONDITIONED`

## Result

The backend profile now provides bounded read paths for:

- administrative organization overview;
- administrative orders and service cases;
- factory production units;
- customer-owned orders and service cases.

All routes derive tenant and subject from verified identity, require an assigned `organization_id`, use a maximum page size of 100 and keyset cursors, and never expose a generic query/proxy surface. Customer service visibility is joined through the warranty owner rather than trusting a request field.

## Reconstruction

- `GO-ENTERPRISE-QUERY-API 0.1.0`: 6 files independently materialized with SHA-256 verification.
- `GO-ELECTROMOBILITY-APPLICATION 0.5.0`: query repository/service/module wired into the executable composition root.
- compatible backend profile: 14 packs, 81 implementation files plus `MATERIALIZATION_RECORD.md`.

## Gates

- Go 1.26.7 `gofmt -d`: PASS, no diff.
- `go test -count=1 ./...`: PASS.
- PostgreSQL 18.6 migrations 0001/0002/0003 and integration suite: PASS.
- organization isolation: PASS.
- customer subject isolation for orders and warranty-backed cases: PASS.
- factory organization linkage through destination purchase order: PASS.
- keyset page continuation without duplicate first item: PASS.
- HTTP unowned organization: 403 PASS.
- `go vet ./...`: PASS.
- `cmd/api` and `cmd/electromobility-api` builds: PASS.
- portable CI runner: 6 tests PASS.
- operational readiness validator: 6 tests PASS.
- container packaging validator: 3 tests and structural validation PASS.

## Conditions

Production still needs representative query plans/cardinalities, product-selected filters and sort orders, sensitive-read audit policy, field-level redaction, selected identity lifecycle, protected browser sessions, telemetry, provider effects and deployment evidence.
