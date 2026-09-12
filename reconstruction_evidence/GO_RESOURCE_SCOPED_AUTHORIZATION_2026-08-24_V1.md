# Go resource-scoped authorization — reconstruction evidence V1

Date: 2026-08-24  
Status: `REBUILD_VERIFIED / CONDITIONED`

## Scope

This reconstruction upgraded the compatible enterprise backend profile so a valid tenant token cannot mutate resources belonging to an organization outside its declared scope.

- OIDC principal accepts bounded `organization_ids` and exposes fail-closed `AllowedOrganization` checks.
- `*` remains the explicit tenant-wide administrative escape hatch.
- Order creation, procurement, factory, inventory, allocation, payment, shipment, service, recall-unit and franchise mutations enforce organization membership at HTTP boundaries.
- Existing-resource transitions also bind organization IDs in PostgreSQL predicates, preventing a caller from pairing an allowed organization with a resource owned by another one.
- Tenant-wide catalog, pricing publication, recall publication and provider callback permissions remain explicit separate capabilities.

## Clean reconstruction

Canonical versions:

- `GO-ENTERPRISE-BACKEND 0.4.0`
- `GO-SUPPLY-FACTORY-INVENTORY-API 0.2.0`
- `GO-COMMERCE-PRICING-PAYMENT-API 0.2.0`
- `GO-FULFILLMENT-SERVICE-FRANCHISE-API 0.2.0`
- `GO-ELECTROMOBILITY-APPLICATION 0.4.0`

The 13-pack enterprise profile materialized 75 implementation files plus `MATERIALIZATION_RECORD.md`. Individual affected packs also materialized independently with every declared SHA-256 verified.

## Executed gates

- Go 1.26.7 `gofmt -d`: PASS, no diff.
- `go test -count=1 ./...`: PASS.
- PostgreSQL 18.6 integration suite against migrations 0001/0002/0003: PASS.
- Negative SQL attempts using an unowned organization for purchase orders, factory units, stock, allocations, payments, shipments, service cases, recall units and franchise agreements: rejected.
- Negative HTTP order creation using a valid permission but an unowned organization: HTTP 403; repository not called.
- Principal assigned/unassigned/empty/wildcard organization tests: PASS.
- `go vet ./...`: PASS.
- builds for `cmd/api` and `cmd/electromobility-api`: PASS.
- CI runner tests: 6 PASS.
- operational readiness validator tests: 6 PASS.
- container packaging validator tests: 3 PASS; structural validation PASS.

## Recovery evidence

The first canonical block synchronization implementation used a regex replacement string. A Go `$` sequence was interpreted by PowerShell and corrupted a candidate block. SHA-256 preflight rejected the profile before materialization. The authoring tool was corrected to use a literal `MatchEvaluator`, affected sections were rebuilt from the formatted verified source, each pack materialized independently, and the complete profile was reconstructed and retested. No corrupted pack was admitted.

## Remaining conditions

Production admission still requires testing a selected real issuer/JWKS lifecycle, permission and membership administration/revocation, customer self-scope, complete query APIs, protected browser journeys, provider adapters and reconciliation, representative load/race/security testing, selected deployment/IaC, telemetry backend, restore/PITR/failover drills and project-specific legal/business rules.
