# Container Packaging Core — Reconstruction Evidence V1

## Scope

- Pack: `CONTAINER-PACKAGING-CORE` 0.1.0
- Files: Docker ignore, multi-stage API image, hardened local Compose, migration entrypoint and Python validator/tests
- Runtime available: Python 3.14.4
- Runtime unavailable: Docker, Podman, OpenTofu and kubectl

## Results

| Gate | Result |
|---|---:|
| six-file materialization and hashes | PASS |
| non-root distroless runtime and exec entrypoint policy | PASS |
| multi-stage/static/trimmed Go build policy | PASS |
| read-only API, no-new-privileges and all-capability drop | PASS |
| isolated data network and loopback-only local API publication | PASS |
| literal-secret negative test | PASS |
| root-runtime negative test | PASS |
| canonical validator and three tests | PASS |
| actual OCI build/run/smoke | NOT RUN — runtime unavailable |

The validator initially rejected safe variable interpolation because its secret regex allowed whitespace backtracking. The canonical regex was corrected and the pack reconstructed twice before all positive and negative tests passed.

## Conditions

On a container-enabled host: resolve all images to reviewed digests, build with SBOM/provenance, scan, execute migrations and API smoke with disposable OIDC, restart/persistence test and teardown. Production additionally requires secret manager/workload identity, TLS database transport, resource limits, orchestration/IaC and platform rollback evidence.
