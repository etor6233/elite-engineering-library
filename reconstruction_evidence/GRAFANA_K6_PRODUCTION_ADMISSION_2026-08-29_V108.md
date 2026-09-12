# Grafana k6 and production admission — V108

Date: 2026-08-29  
Authority: current V108 reconstruction evidence  
Admission: `REBUILD_VERIFIED / CONDITIONED`

## Outcome

Elite now separates three different facts that must never be conflated:

1. an exact official tool was selected and executed;
2. its project-specific semantic result passed;
3. all eight production controls jointly authorize the same immutable release.

`SECURE-OPS-DELIVERY-CORE 0.4.0` materializes 16 files. The generic runner verifies the selected binary SHA-256, requires explicit execution authorization, uses an argv array with `shell=False`, restricts environment variables, applies a timeout, captures bounded stdout/stderr and records canonical argv, working directory and output hashes. Its receipt is execution evidence only.

The Grafana-specific adapter consumes three distinct exact executions—`BASELINE`, `SOAK`, and `SATURATION_RECOVERY`—plus an independently approved error budget. It verifies project/environment/release identity, k6 identity, argv, workload-script bytes, real target, minimum declared duration, official exit code, k6 v2 summary version, iteration count, HTTP failure rate and HTTP duration p95. It then emits the exact five semantic assertions required for `LOAD_RESILIENCE`. Missing soak, nonzero exit, script or argv drift, legacy summary, budget breach or fewer than two distinct budget approvers fails closed.

## Exact official Grafana identity

| Item | Exact value |
|---|---|
| repository | `grafana/k6` |
| release | `v2.2.0`, release ID `367956260`, published `2026-08-10T14:01:35Z` |
| tag object | `9b7cf93462a32ef2dc74c388f58435cba327c931` (`unsigned`) |
| commit | `00a9a1b7f552d6bb4337278b10ae25aac0f4e666` (`verified=true`) |
| tree | `995b8621afd0f8457c306823fe75a255e76c128e` |
| source ZIP | 13,962,538 bytes; SHA-256 `0ed19683455d33569ee7b0d420b6625f7ce2897108d6db06e77910b9c14d3706` |
| license | `AGPL-3.0-only`; `LICENSE.md` SHA-256 `45bd5efa5e66840e253253f3cc172d07213cb812206cda7bb71cd5b9dbd5a985` |
| `go.mod` | SHA-256 `25d254eefc7d3c0c6b09f214b0b17f3fc34229d3fde4dc1275710a5a5928c5d6` |
| Windows release ZIP | asset `508824505`; 31,362,450 bytes; SHA-256 `ceb2b1e1cf9dbe1303c6c33ec83ffda86dda5c610b4def92064d3c7ebae8d9f4` |
| checksums | asset `508824465`; 746 bytes; SHA-256 `cecbd810434043492a60fc6f2d32f60036118927e0164757271823f6c96c64f2` |
| SPDX | asset `508824497`; 483,669 bytes; SHA-256 `34085f8f058d799416e26d4cb42051896a4489661d812521073649418add3d2d` |
| `k6.exe` | 66,343,632 bytes; SHA-256 `87dfa91bc3e47bc4bd77911d59d7ff79f25cd76fb8322c97072f08f08a0da5ed` |
| version output | `k6.exe v2.2.0 (commit/00a9a1b7f5, go1.26.5, windows/amd64)` |

The exact executable ran through the generic boundary with exit `0`. Its stdout was 59 bytes/SHA-256 `87044a01f6e000c8660bb84bffb603b782726ef750d5733e951820744beebad7`; stderr was empty/SHA-256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`; canonical argv `['version']` was SHA-256 `1c7814b63a0c9eef6a2e127918c2c30610720d11fe6111c1f497b64167a2ae34`.

PostgreSQL 18.6 is additionally fixed at mirror commit `724edf9bde9d356724ad384a2e196edc3c9f80f7`, tree `69f81c582d01133710c1eeb9c12fadf7f47633e4`, archive SHA-256 `6b49182dba7aaece2f9c5fc7bd14d9f4e1ea3f767396ae13cfab7e1b09059063` and official tar.bz2 SHA-256 `555610c24d53e4316da5b7d3fc25c279d96856d5e0e23ee308c328c5fa881d9f`. Its local adapter requires five bound runs and measured RPO/RTO; eight regressions pass. The lightweight tag/unsigned commit and absence of a real project restore remain explicit.

OWASP ZAP 2.17.0 is fixed at verified commit `8a1bff313f4d183dba5aa154ecbe89ad751c9153`, tree `5e1b5e4842803871286950fa4082d8dd504a8c34`, archive SHA-256 `44d3e5d6b9a777652ff82ede1a7de64e4ef0aded453c1f161cee03c6b5f3faec`, Core ZIP SHA-256 `0cb73b7f72d12c263fb61de304edb82a455d7aa4e1813c216c061765c306f5b7` and release SBOM SHA-256 `c91137d66c34a7e2803cbd4fbaf41f2e16201f754b9fc7c954293939a32371bf`. Its tag is unsigned. The local adapter requires authorized authenticated web/API runs and manual review; eight regressions pass. No target was scanned.

## Executed reconstruction

- upstream acquisition: `UPSTREAM_ACQUISITION_TEST_PASS negatives=9`;
- source profiles: 14 valid, 5 negative, 1 positive; enterprise selection `17`;
- secure-operations round-trip: 16/16 files materialized;
- production-admission/runner/k6/PostgreSQL/ZAP suite: 38 cases, 37 PASS and one symlink-only platform skip;
- global structure: `VERIFY_LIBRARY_PASS`, 73 packs, 724 files, 36 profiles, backend 16/111 and web 5/69.
- executable audit: `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=73 upstream_sources=118`;
- Foundation with official Go 1.26.7: backend Go tests, web frozen install/typecheck/21 tests/build, Playwright runtime 4 + target 8 and Lighthouse five-run gate all PASS; final exit `0` and `PROJECT_CONDITIONED` preserved.

The skipped symlink case is not presented as a PASS. Symlink rejection remains implemented and is also covered by library-level path gates.

## Failure learning

`LIB-FAIL-1194` through `LIB-FAIL-1200` record the exact local failures and corrections from source/test discovery, array invocation, cleanup boundary, parameter recurrence, k6 summary interpretation and missing argv binding. `UP-FAIL-195` through `UP-FAIL-197` record Grafana, PostgreSQL and ZAP conditions. Current memory is 1,200 local failures plus 197 upstream conditions, 1,397 unique IDs.

## Non-claims and remaining production work

No REVESTEX or other real target was loaded in this cycle. Therefore this evidence does not assert throughput, latency, capacity or `PRODUCTION_ADMITTED` for any project. A project must provide an authorized target, immutable release digest, approved workload scripts and error budget, representative data, observability, baseline/soak/saturation-recovery executions and fresh evidence.

The final production authority still requires five additional project-specific semantic controls alongside the implemented load, PostgreSQL-recovery and offensive-security adapters: edge/CDN/WAF, identity/authorization, providers, deploy/canary/rollback and business acceptance. The library must not claim a concrete project production-ready until all eight receipts bind to the same target release and pass the final validator.
