# Go Fulfillment, Service and Franchise API — Reconstruction Evidence V1

## Scope

- Date: 2026-08-24
- Pack: `GO-FULFILLMENT-SERVICE-FRANCHISE-API` 0.1.0
- Composition: Go core + catalog/CRM + supply/factory/inventory + commerce + migrations 0001/0002/0003
- Runtime: Go 1.26.7, pgx 5.10.0, PostgreSQL 18.6

## Results

| Gate | Result |
|---|---:|
| six-file materialization and SHA-256 verification | PASS |
| shipment creation and controlled transition | PASS |
| service case organization/stock invariant and optimistic version | PASS |
| stale service transition rejection | PASS |
| recall activation and duplicate-unit rejection | PASS |
| franchise agreement transition and stale-state rejection | PASS |
| queued/sent/delivered communication flow | PASS |
| twelve transactional outbox events | PASS |
| domain and HTTP permission tests | PASS |
| fresh materialization and `gofmt -d` | PASS |

The first compile gate rejected malformed compact syntax; subsequent PostgreSQL runs rejected an invalid organization enum and an incomplete customer fixture. Each fault was corrected in the working source, all gates were rerun, and only the successful formatted files were embedded back into the canonical Markdown with new hashes.

## Conditions

Production still requires shipment/order-line linkage, addresses and custody proof, official carrier adapters and webhook reconciliation, service parts/labor/appointments, warranty claims, regulatory recall evidence, territory conflict/approval rules, consent and suppression policy, provider templates, resource authorization, pagination, audit and concurrency/load testing.
