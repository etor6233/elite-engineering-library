# PostgreSQL Backup/Restore Core — Reconstruction Evidence V1

## Scope

- Date: 2026-08-24
- Pack: `PG-BACKUP-RESTORE-CORE` 0.1.0
- Source database: `elite_composed_v2`, a previous clean composition containing `PG-TX-FOUNDATION`, portable configuration, Go backend migrations and the electromobility vertical
- Runtime: PostgreSQL 18.6 official Windows binaries; PowerShell 7
- Secret handling: local trust configuration; no password or connection URI written to the pack, command output or evidence

This evidence authorizes clean reconstruction and the demonstrated logical restore workflow. It does not authorize production, PITR, failover, representative-volume RPO/RTO or a particular cloud backup service.

## Results

| Gate | Result | Evidence |
|---|---:|---|
| materialize five files and validate embedded hashes | PASS | `Materialized 5 files` |
| `pg_dump` custom format, zstd level 6 | PASS | PostgreSQL 18.6; 28,867-byte artifact |
| machine-readable manifest | PASS | schema `elite.postgres.backup-manifest.v1`; 56 TOC entries |
| backup integrity | PASS | SHA-256 `4f1fbbe433ec7e7ee0de91769a18022efaee31e0c5eeb62d5ab493ba8209cd4b` |
| isolated restore | PASS | generated `restore_drill_YYYYMMDDHHMMSS_<nonce>` database; `--single-transaction --exit-on-error` |
| restored invariants | PASS | required relations, tenant/profile integrity, outbox and jobs |
| cleanup | PASS | zero `restore_drill_%` databases after the run |
| corrupt-artifact negative test | PASS | altered copy rejected by SHA-256 before database creation; database count remained zero |

Measured restore plus invariant and cleanup duration was 0.548 seconds on a small local test database. This is a smoke measurement, not an operational RTO.

## Defect found during reconstruction

The first invariant draft referenced hypothetical job columns (`attempt_count`, `lease_owner`, `lease_expires_at`). The actual admitted foundation exposes `attempts`, `claimed_by` and `claimed_until`. The first restore correctly failed, cleanup removed the drill database, the invariant was corrected to the materialized contract, hashes were regenerated and the complete run passed. This failure is retained because it demonstrates that the drill detects schema/contract drift.

## Remaining production conditions

- platform-specific secret injection and least-privilege backup/restore roles;
- encrypted, immutable and access-audited backup storage;
- retention/deletion plus cross-account or cross-region copies;
- WAL archiving and point-in-time recovery;
- scheduled restores using backups produced by the actual service;
- representative data volume and dependency-aware recovery;
- monitored backup freshness/failure and measured RPO/RTO;
- failover and application reconnection exercises.
