# Go Reliable Async Workers — Reconstruction Evidence V1

## Scope

- Date: 2026-08-24
- Pack: `GO-RELIABLE-ASYNC-WORKERS` 0.1.0
- Composition: Go enterprise backend 0.3.0 + PostgreSQL transactional foundation 0.1.0
- Runtime: Go 1.26.7, pgx 5.10.0, PostgreSQL 18.6

## Results

| Gate | Result | Detail |
|---|---:|---|
| clean worker materialization | PASS | 6 files; embedded hashes verified |
| clean composition with Go core | PASS | no path collisions; module `elite.local/enterprise` |
| formatting | PASS | `gofmt -d` produced no diff after canonical formatting was embedded |
| unit tests | PASS | outbox/job processors separate success from retry and forward event ID as provider idempotency key |
| inbox failure rollback | PASS | handler error left no receipt or local mutation |
| inbox deduplication | PASS | committed handler ran once; repeated event returned `processed=false` |
| job ownership | PASS | wrong worker could not complete a claimed job |
| retry and exhaustion | PASS | first failure requeued; second claim incremented attempt; final failure became terminal and was not reclaimed |
| existing repository/outbox tests | PASS | order repository and disjoint outbox claim/ownership/reclaim remained green |
| Go gates | PASS | `go test -count=1 ./...`, `go vet ./...`, `go build ./...` |

## Defects detected during reconstruction

The first composition rejected a placeholder module import and a missing worker argument. The first PostgreSQL run rejected two parameterized SQL commands combined into one prepared statement. Each defect was corrected in the canonical Markdown, hashes were regenerated and all gates were rerun from empty directories.

## Semantics and remaining conditions

Outbox is at-least-once, never exactly-once. Remote providers must honor the event ID as an idempotency key. Inbox atomically combines a local database handler and receipt; remote side effects are forbidden inside it. Production still needs concrete transports and handlers, provider contract tests, process shutdown/drain, metrics/SLO, retry classification with jitter/budgets, dead-letter operations, concurrency/load/crash tests and least-privilege database roles.
