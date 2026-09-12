# Go Supply, Factory and Inventory API — Reconstruction Evidence V1

## Scope

- Date: 2026-08-24
- Pack: `GO-SUPPLY-FACTORY-INVENTORY-API` 0.1.0
- Composition: Go backend + public catalog/CRM API 0.1.1 + migrations 0001/0002/0003
- Runtime: Go 1.26.7, pgx 5.10.0, PostgreSQL 18.6

## Results

| Gate | Result |
|---|---:|
| six-file materialization and embedded hashes | PASS |
| modular HTTP registration without route collision | PASS |
| purchase-order creation and valid transition | PASS |
| stale purchase-order state/version rejection | PASS |
| production-unit registration and state transition | PASS |
| stock receipt and `received_at` activation | PASS |
| stale stock state/version rejection | PASS |
| skipped domain transitions rejected before repository | PASS |
| permission separation for procurement/factory/inventory | PASS |
| six transactional outbox events | PASS |
| `gofmt -d`, `go test -count=1 ./...`, `go vet ./...`, `go build ./...` | PASS |

Every transition event ID is supplied by the application ID generator; deterministic IDs derived from state/version were rejected during design because they could collide across aggregates.

## Conditions

The state machines are safe baselines, not jurisdictional or contractual truth. Production requires organization/resource-scoped authorization, supplier and PO line models, Incoterms and tax/legal decisions, quality/homologation rules, serial/VIN normalization, concurrent allocation tests, audit events, pagination/query APIs, operator tooling, metrics and representative crash/load tests.
