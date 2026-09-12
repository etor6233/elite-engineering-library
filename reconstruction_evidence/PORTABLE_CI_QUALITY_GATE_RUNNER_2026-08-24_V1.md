# Portable CI Quality Gate Runner — Reconstruction Evidence V1

## Scope

- Date: 2026-08-24
- Pack: `PORTABLE-CI-GATE-RUNNER` 0.1.0
- Runtime: Python stdlib
- Execution boundary: direct argv with `shell=False`; external environment inherited but never serialized by the plan

## Results

| Gate | Result | Detail |
|---|---:|---|
| clean materialization | PASS | 3 files; embedded SHA-256 verified |
| unit tests | PASS | 6/6 |
| nonzero failure | PASS | exit 7 recorded as required gate failure |
| timeout | PASS | process exceeded one second and was recorded `TIMEOUT` |
| coverage validation | PASS | missing required security category rejected |
| secret boundary | PASS | inline `env`/token-shaped plan rejected; output tokens redacted and bounded |
| unresolved example | PASS (negative) | invalid plan produced exit 2 and evidence status `INVALID` |
| instantiated runner | PASS | format, unit, static, build, security, supply-chain, integration and recovery all executed |

The instantiated plan SHA-256 was `d88c62f5d4274d841bd21667230aec5c179c2c9b738969b3939c841c20fc5758`; the resulting evidence contained timings, exit codes and bounded stdout/stderr for eight gates.

## Conditions

The test commands only prove orchestration. Each project must replace them with real compiler, test, scanner, SBOM, provenance, integration and restore commands; run on an isolated, pinned CI platform; inject secrets outside the plan; preserve immutable evidence; and exercise the actual deploy/rollback path.
