# Leader Privacy, Runtime and Cloud Portability Sources — 2026-08-25 V1

## Scope

This evidence adds three exact official public source archives to `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.3.7. The acquisition orchestration is authored locally; the acquired source remains under each upstream license.

## Presidio 2.2.364

- current repository: `data-privacy-stack/presidio` (Microsoft-origin project transferred to Data Privacy Stack);
- release: `2.2.364`;
- commit: `779dbd286d5ef4d1fbe2514275fb1bce358f2417`;
- archive: 108,784,615 bytes;
- archive SHA-256: `dbdd5a80d5f7e8b079f8b19eb5d75449589a80d3251beec1a7bbd19909feb3c9`;
- license: MIT, SHA-256 `f3e86ee59a49bcfb0d9a9547484d55224ea7b2d04f95b1947b4d18d17f6de535`;
- NOTICE SHA-256: `07997418ab87066ace8b0c9aceb76bea1291a77c62b16434b3d0b925af5b883f`;
- acquired tree: 878 files, 160,444,030 bytes, 268 files under test/e2e-test paths.

The source declares Python `>=3.10,<3.15` for the product packages. Python 3.14 compilation passed independently for:

```text
presidio
presidio-analyzer
presidio-anonymizer
presidio-cli
presidio-image-redactor
presidio-structured
```

The complete-tree compile did not pass and is retained honestly:

- two documentation sample files failed while Python attempted to create `__pycache__` under a path too deep for this Windows audit path;
- `presidio-anonymizer/tests/services/test_aes_cypher.py` emitted a Python 3.14 `SyntaxWarning` for a non-raw regex escape;
- the repository has no root dependency lock; `pyproject.toml` is configuration, not a frozen graph.

Therefore Presidio is `PINNED_CANDIDATE`, not a universal privacy guarantee. Project adoption still requires selected recognizers/languages, false-positive/false-negative evaluation, dependency lock, data-policy review and proof that anonymization preserves required document evidence.

## Dapr Runtime 1.18.3

- repository: `dapr/dapr` (Microsoft-origin, CNCF project);
- release: `v1.18.3`;
- commit: `2dcf2e3548faf3c740b3d29ee300e2132db65cff`;
- archive: 20,269,733 bytes;
- archive SHA-256: `dd64436b01c06bdf59db6c7c949dce81eff36c293581980af93e12a80bfd6681`;
- license: Apache-2.0, SHA-256 `9523eacf5421b65637420755c00b5175f7804b6a6eb736b86cd10ea6b3937dcf`;
- `go.sum` SHA-256: `e7c4c2dba67d067423e9cd85a8d13da8a86028cfd81d5b831a74b2ecb8c95e79`;
- acquired tree: 4,297 files, 34,315,342 bytes, 2,624 under test paths;
- declared Go version: 1.26.6.

Archive, license and lock acquisition passed. Build/tests were not run because this host has neither Go nor Docker. Dapr remains `CONDITIONAL_PLATFORM`; its presence does not justify microservices or sidecars without a blueprint need and target component/security/load/recovery evidence.

## Google Go Cloud 0.46.0

- repository: `google/go-cloud`;
- release: `v0.46.0`;
- commit: `60f0d749b071334179a08aa9b6c13d2fe239d120`;
- archive: 3,875,103 bytes;
- archive SHA-256: `4f9e26ff55800fb73dece7b6e4bf713f20ec98843cf88aa08d6c0fe06417e885`;
- license: Apache-2.0, SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`;
- `go.sum` SHA-256: `a5d9151661ea65b0e9711a7ce66ebd263944535d1abb65086a135234b5110896`;
- acquired tree: 1,079 files, 13,781,525 bytes;
- declared Go version: 1.25.0.

Archive, license and lock acquisition passed. Build/tests were not run because Go is absent. Adoption must select only required drivers and prove provider behavior, credentials, quotas, telemetry, recovery and exit semantics.

## Acquisition and library gates

The exact three-source acquisition returned:

```text
UPSTREAM_SOURCE_ACQUIRED id=presidio-2.2.364
UPSTREAM_SOURCE_ACQUIRED id=dapr-runtime-1.18.3
UPSTREAM_SOURCE_ACQUIRED id=google-go-cloud-0.46.0
```

After canonical pack synchronization, the clean library audit returned:

```text
VERIFY_LIBRARY_PASS
packs=23 materialized_files=224 markdown_files=151
profile=ENTERPRISE_BACKEND_PACK_PLAN.md implementation_files=95
profile=ENTERPRISE_WEB_PACK_PLAN.md implementation_files=44
UPSTREAM_ACQUISITION_TEST_PASS
UPSTREAM_LOCK_VALID sources=49 selected=49
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=23 upstream_sources=49
```

## Failure learning

The audit preserved failures rather than rewriting history:

- `LIB-FAIL-024`: nested PowerShell interpolation; corrected with direct-shell execution;
- `LIB-FAIL-025` and `LIB-FAIL-026`: Go and Docker absent, retained as external toolchain blockers;
- `LIB-FAIL-027`: omitted mandatory validate-only parameters; corrected explicitly;
- `LIB-FAIL-028`: array parameters crossing `pwsh -File`, recurrence count 2; corrected with the call operator and a real PowerShell array;
- `UP-FAIL-018` and `UP-FAIL-019`: Presidio full-tree path failure and Python warning retained.

Successful direct-shell probes, validate-only output, three-source acquisition and two-file pack synchronization provide the regression evidence for the corrected harness failures. The missing toolchains and upstream issues remain open/conditioned.

## Claim boundary

This evidence proves exact source acquisition and the stated local gates. It does not prove production privacy, distributed runtime suitability, cloud provider compatibility, document field accuracy or REVESTEX business behavior.
