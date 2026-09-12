# Google Cloud Document AI Python 3.15.0 — source, runtime and admission audit V2

## Scope and claim boundary

This reconstruction audits the official Google Cloud Document AI Python client and the local `AUTHORED` orchestration that invokes its public API. Google owns `google-cloud-documentai`, generated clients, protobuf models and the cloud service. Elite owns the profile/security/evidence wrapper and its tests; none of that glue is labelled as Google code.

No Google Cloud credential, project, processor, billable request or business corpus was available or used. Therefore this evidence proves exact source/artifact identity, local installation, public SDK request construction and fail-closed orchestration. It does not prove extraction accuracy, a live processor, latency, quota, price, residency, review or production readiness.

## Official identity and artifacts

- repository: `googleapis/google-cloud-python`;
- tag: `google-cloud-documentai-v3.15.0`;
- commit: `5accbb42b3f8372c4d03ef59a6b0000b165f2200`;
- GitHub commit verification observed through the official API: `verified=true`, `reason=valid`;
- source archive: 202,963,530 bytes, SHA-256 `3d0af0f2e9c733a5bb5195e5f1e312dbed2163491140526efd192398ceebdfe2`;
- PyPI wheel: `google_cloud_documentai-3.15.0-py3-none-any.whl`, 310,694 bytes, SHA-256 `f040f4f9db43411184197a808b11fde52b580723a9fca336f43ea7c6a885bfd8`;
- PyPI sdist: `google_cloud_documentai-3.15.0.tar.gz`, 362,736 bytes, SHA-256 `d50b69a8a62aaf803b0d926f93b6579df440f5f779f45fec2c5972768a1cdbbd`;
- package metadata: Python `>=3.10`, Apache-2.0, Google LLC, PyPI Trusted Publishing;
- dated stable version check: PyPI still published 3.15.0 as the latest stable release on 2026-08-27.

The no-Git archive contained 52,474 files; `packages/google-cloud-documentai` contained 214. All 75 source/test files present in the sdist and all 67 Python files present in the wheel had counterparts with identical SHA-256 in the verified tag. There were zero missing or divergent compared files.

Exact package evidence fixed in `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.25`:

- root/package `LICENSE`: `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`;
- `setup.py`: `c2d59c375bdc5ef3085b0800a3852ce694beb82136dd874e2633fb2bc47c61eb`;
- `noxfile.py`: `21d8c0116d214fbddd5e5b4cf2735560382da6a9dec6c15c0101c65200ff3658`;
- CPython 3.12 test constraints: `6e042323ef99e9de87d626253dabb092956163c7070676fc67dcc95559a58783`.

Actual acquisition through the rebuilt no-Git runner returned `UPSTREAM_SOURCE_ACQUIRED id=google-cloud-documentai-python-3.15.0`, produced the exact receipt and reconstructed all 52,474 files. The lock validates 71/71 sources. The Google/document profiles now select 7/32 sources.

## Official test result retained without repair

The official sdist tests were copied outside the source root and executed against the exact wheel in a new CPython 3.12.13 Windows x86-64 environment. This avoided mixing the installed wheel and source namespace.

- protobuf `upb`: 1,879 PASS, 6 FAIL, 67 warnings;
- protobuf `python`: 1,879 PASS, 6 FAIL, 67 warnings;
- the same six tests fail in both runs: synchronous/async mTLS endpoint and certificate-source expectations across v1 and v1beta3 clients;
- resolved `google-auth`: 2.57.0, permitted by Google's open dependency range;
- no test was patched, skipped or relabelled as PASS.

This is retained as `UP-FAIL-075`. The package manifest provides ranges rather than an integral hash lock. Elite therefore admits the SDK only as `PINNED_CANDIDATE / CONDITIONED`; it does not claim the official suite is green.

## Productive target lock

For CPython 3.12/Windows x86-64, the runtime graph was resolved separately from audit tooling and frozen as 19 exact wheels in `requirements-windows-py312.lock`. Each line contains an exact version and wheel SHA-256. A clean global run executed:

```text
pip install --require-hashes -r requirements-windows-py312.lock
pip check
python -m unittest -v test_document_ai_runtime.py
```

The install, dependency check and six wrapper tests passed. A dated Google OSV batch query over all 19 exact PyPI package versions returned zero known matches at `2026-08-27T15:10:36Z`. This is a dated database observation, not a guarantee that no vulnerability exists or will be published later.

## Weakness demonstrated and corrected

`GOOGLE-DOCUMENT-AI-RUNTIME 0.1.0` passed its three declared tests but an adversarial probe showed that it accepted:

- an arbitrary processor version not selected by its profile;
- arbitrary bytes named with a `.pdf` suffix;
- no receipt from the secure local-file gate.

That weakness is `LIB-FAIL-231`. Version 0.1.2 closes it:

- the CLI accepts a document class, profile and security receipt, never a free processor resource;
- the class must occur exactly once with decision `REQUIRED`;
- profile schema/provider/SDK, processor type, exact processor version, schema version and `automatic_storage=false` are mandatory;
- input SHA-256, byte length and detected MIME must match an `ADMITTED` receipt containing policy, approval, ClamAV, Magika and YARA-X evidence;
- the exact Google `ProcessRequest` uses the profile-owned processor version;
- the complete `ProcessResponse` and an analysis receipt v2 are committed atomically;
- profile, security receipt, input and provider response hashes are linked;
- no field receives business-storage authority.

The distributed profile keeps supplier invoices at `CONFIGURATION_REQUIRED` and every other class blocked. Tests use a local explicit profile fixture; the template itself cannot call the provider.

## End-to-end composition and gates

`GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md` now composes five layers and 41 files:

1. official upstream acquisition;
2. secure local-file ingestion;
3. official document SDK artifact acquisition;
4. Google Document AI runtime;
5. strict document-field evaluation.

Observed final gates:

```text
VERIFY_LIBRARY_PASS
packs=41 materialized_files=366
profile=GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md implementation_files=41
UPSTREAM_LOCK_VALID sources=71 selected=71
SOURCE_PROFILE_VALID profile=google-secure-document-pipeline selected=7
SOURCE_PROFILE_VALID profile=document-intelligence-leaders selected=32
GOOGLE_DOCUMENT_AI wrapper tests: 6 PASS
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=41 upstream_sources=71
```

The executable runner always materializes and validates the Google lock. Offline mode emits an explicit `SKIPPED`; `-AllowNetwork` on the exact target creates a fresh venv, installs with hashes, runs `pip check` and executes all six wrapper tests.

## Retained lessons and remaining project authority

Local failures `LIB-FAIL-229..237` preserve wrong path assumptions, source/wheel namespace mixing, the original profile/security bypass, an OSV null-row parser error, patch/updater invocation mistakes and brittle negative-test output parsing. Each corrected condition was rerun; none was deleted from the ledger. `UP-FAIL-075` remains open because only Google can publish a corrected official suite/lock or compatible release.

Before the runtime can perform even a sandbox call, the project owner must supply and approve: GCP organization/project and billing, location/residency, exact processor and processor version, ADC identity and least-privilege IAM, schema version, representative authorized corpus and ground truth, review/correction route, quota/cost/latency limits, retention and redaction, load/failure tests, monitoring, rollback and field-level acceptance. These values must be requested; the agent is forbidden to invent them.

Official references: [Google Cloud Python tag](https://github.com/googleapis/google-cloud-python/tree/google-cloud-documentai-v3.15.0/packages/google-cloud-documentai), [PyPI 3.15.0](https://pypi.org/project/google-cloud-documentai/3.15.0/), [Document AI process method](https://docs.cloud.google.com/document-ai/docs/reference/rest/v1/projects.locations.processors/process), [processor versions](https://docs.cloud.google.com/document-ai/docs/manage-processor-versions).
