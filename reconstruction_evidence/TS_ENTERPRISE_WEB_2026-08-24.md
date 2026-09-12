# Reconstruction Evidence — TS Enterprise Web

```yaml
evidence_id: "TS-EW-20260824-V2"
pack_id: "TS-ENTERPRISE-WEB"
pack_version: "0.1.0"
pack_sha256: "89876dfe826f63ab4e6eebd941aae7f687007e6e90bcd0e274db7b756312c3f9"
materialized_file_count: 57
ordered_file_hash_aggregate_sha256: "3785ad78fe75092b31189ac6577663e9a15c9a4c52b72794c0a97a9683094143"
environment: "Windows / clean target directory"
node: "v24.14.1"
pnpm: "11.19.0"
verified_at: "2026-08-24T17:20:39-03:00"
result: "PASS"
```

## Independence

The target `rebuild_verification/ts-enterprise-web-v2` began empty. `materialize_markdown_pack.ps1` read `implementation_packs/TYPESCRIPT_ENTERPRISE_WEB_GOLDEN_PATH.md`, verified each declared SHA-256 and wrote 57 files. Build and tests ran from the reconstructed target, not from `golden_starter/`.

## Commands and results

| Gate | Result |
|---|---|
| materialize 57 files + verify hashes | PASS |
| `pnpm install --frozen-lockfile` | PASS; 88 packages from locked graph |
| `pnpm config:check` | PASS |
| `pnpm readiness:baseline` | READY |
| `pnpm typecheck` | PASS |
| `pnpm test` | PASS; 6 files / 15 tests |
| `pnpm build` | PASS |
| `pnpm audit --prod` | no known vulnerabilities |

## Smoke journeys

| Journey | Expected | Result |
|---|---|---|
| `/`, `/models`, `/admin`, `/customer`, `/factory` | HTTP 200 | PASS |
| `/api/health`, `/api/platform` | HTTP 200 | PASS |
| unknown lead custom field | HTTP 400 | PASS |
| valid lead creation | HTTP 201 | PASS |
| same idempotency key + same body | HTTP 200 + same resource | PASS |

## Defect discovered and corrected

The first smoke pass found that an unknown custom field was rejected internally but classified as HTTP 500. The domain now raises `InvalidLeadInputError`, the HTTP boundary maps it to stable problem code `INVALID_LEAD_INPUT`/400, and unit/HTTP tests cover the behavior. The pack was regenerated and reconstructed again after the fix.

## Remaining conditions

This evidence supports `REBUILD_VERIFIED`, not unrestricted production readiness. External OIDC, managed PostgreSQL, secret manager, least-privilege DB roles, provider sandboxes, exported telemetry, deployment, backup/restore and business-module coverage remain separate gates.

