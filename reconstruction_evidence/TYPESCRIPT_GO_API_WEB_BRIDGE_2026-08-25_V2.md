# Reconstruction evidence — TS-GO-API-WEB-BRIDGE 0.2.0

## Scope

Reconstruction of the standalone TypeScript frontend/BFF composed from:

- `TS-GO-API-WEB-BRIDGE 0.2.0` — 31 files;
- `TS-OIDC-PORTAL-ADAPTER 0.1.2` — 10 files;
- `DEPENDENCY-LICENSE-EVIDENCE-CORE 0.1.1` — 3 files.

`TS-ENTERPRISE-WEB` was not selected and no file, runtime dependency, database client, migration, outbox or worker was inherited from it. The resulting web profile contains 44 non-colliding materialized files.

## Environment

- Date: 2026-08-25
- OS: Microsoft Windows 11 Pro
- PowerShell: 7.6.4
- Node.js: 24.14.1
- pnpm: 11.19.0
- Python: 3.14.4

## Reproducible procedure

1. Materialize `MARKDOWN-COMPOSITOR 0.2.1` with the root `materialize_markdown_pack.ps1`.
2. Run the compositor self-test suite.
3. Compose `markdown_system/ENTERPRISE_WEB_PACK_PLAN.md` into an empty destination.
4. From that destination run:

```powershell
pnpm install --frozen-lockfile --offline
pnpm typecheck
pnpm test
pnpm build
Push-Location supply_chain
python -m unittest test_license_gate.py
Pop-Location
```

## Observed results

- Root materializer: PASS, 31 bridge files reconstructed with exact manifest/block parity and SHA-256 validation.
- Compositor self-tests: PASS for success, CRLF, condition acknowledgement, integrity, manifest, CREATE-only, uniqueness, selection and reserved-path cases.
- Profile composition: PASS, 44 files from 3 packs with no collisions.
- Frozen offline install: PASS; lockfile accepted, 73 packages reused, 0 downloaded.
- TypeScript/Next type generation: PASS.
- Vitest: PASS, 5 files and 14 tests.
- Next production build: PASS; public home/models, public lead BFF, OIDC auth, customer, administration and factory surfaces emitted.
- License-gate unit tests: PASS, 2 tests.

## Boundary of the evidence

Offline installation proves that the pinned lockfile can be installed from the already populated local content-addressable store; it does not prove future registry availability. This reconstruction establishes deterministic materialization and local build/test integrity, not production readiness. A target project must still prove Go API contracts, selected IdP behavior, secrets and key rotation, browser role/tenant/object authorization, distributed antiabuse, final CSP/CSRF and TLS proxy policy, accessibility with browser/assistive technology, representative load and failure behavior, observability, legal review of the distributed artifact, rollout and rollback.
