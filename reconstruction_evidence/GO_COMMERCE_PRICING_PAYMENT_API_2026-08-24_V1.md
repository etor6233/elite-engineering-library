# Go Commerce, Pricing and Payment API — Reconstruction Evidence V1

## Scope

- Date: 2026-08-24
- Pack: `GO-COMMERCE-PRICING-PAYMENT-API` 0.1.0
- Composition: Go core + catalog/CRM API + supply/factory/inventory API + migrations 0001/0002/0003
- Runtime: Go 1.26.7, pgx 5.10.0, PostgreSQL 18.6

## Results

| Gate | Result |
|---|---:|
| six-file materialization and embedded SHA-256 validation | PASS |
| canonical files are `gofmt`-clean after fresh reconstruction | PASS |
| price-book creation, activation and public current-price query | PASS |
| server-side order pricing and total invariant | PASS |
| optimistic order placement and stale-state rejection | PASS |
| same-organization/same-variant stock allocation | PASS |
| duplicate allocation and stale stock version rejection | PASS |
| durable payment-attempt creation and outbox event | PASS |
| allowed payment transition and stale transition rejection | PASS |
| HTTP public/protected boundary and permission separation | PASS |
| full composed `go test -count=1 ./...`, `go vet ./...` and both builds | PASS |

The first PostgreSQL run exposed ambiguous untyped parameters inside `jsonb_build_object`. The fix added explicit SQL text casts in the canonical Markdown, refreshed hashes, reconstructed into a new empty directory and repeated all dependent gates successfully. This is an exercised example of the agent error-recovery protocol.

## Conditions

This pack creates durable payment intent/state, not a real charge. Production still requires jurisdictional tax/accounting and price-approval rules, discounts/financing, immutable commercial audit, provider adapters and signed webhook inbox, idempotent response replay, reconciliation, partial capture/refund/dispute behavior, PCI-scope and fraud decisions, organization/resource authorization, concurrent multi-session and representative load tests.
