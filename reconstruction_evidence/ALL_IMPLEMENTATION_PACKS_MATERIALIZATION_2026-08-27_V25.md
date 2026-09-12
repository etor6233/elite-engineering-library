# All Implementation Packs Materialization — 2026-08-27 V25

## Governing delta

This record succeeds V24 for the current source tree. Historical records remain immutable evidence of their own snapshots.

The exact official `GoogleCloudPlatform/document-intake-accelerator` source at signed commit `956e3cc7900338d1146d82c5ee2bd4b2551f70bb` was rebuilt deeply without patching it, providing credentials or deploying billable GCP resources. Its official architecture is preserved, but its runtime, locks, tests, security, correctness, API boundary and IaC do not satisfy immediate enterprise adoption.

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` advances from 0.4.20 to 0.4.21. The immutable Google source remains acquirable only as evidence/sample; its lock and the Google/document-intelligence profiles now expose the rejection before network acquisition. All thirteen composition plans that depend on this core were aligned to 0.4.21.

## Exact Google evidence

`GOOGLE_DOCUMENT_INTAKE_ACCELERATOR_REBUILD_2026-08-26_V1.md` records:

- signed head commit, Apache-2.0, 43.108.717-byte archive and SHA-256 `38204b7365ed2df6b00633c63943d588fcb085d876e5decd7cd96d99ca2309bb`;
- exact historical CPython 3.7.17/Node 17.9.1 toolchains and independently verified support libraries;
- 23 Python requirements files, 264 entries, 51 unpinned and one file with unresolved merge-conflict markers;
- UI immutable install PASS, but Jest one suite failed before tests and CI build failed;
- npm audit: 115 vulnerable packages, including 9 critical and 68 high;
- OSV: 189 affected records, 111 unique package-version pairs, 337 primary IDs and severity up to 9.8;
- Python suite collection blocked by BigQuery ADC construction during import even with a local Firestore emulator;
- upload APIs without server token verification, content-type security, byte/rate limits or antimalware;
- aggregate/fuzzy autoapproval without a demonstrated closed-schema zero-FA/FD/FN/FP gate;
- broad IAM/org-policy weakening and no demonstrated idempotent cross-store transaction/outbox/DLQ contract.

The failure ledger retains those findings as `UP-FAIL-059` through `UP-FAIL-068`. Audit, harness and cleanup failures are retained through `LIB-FAIL-203`; no failed command is silently converted into evidence.

## Executed gates before this record

The updater regression passed, and `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.21` materialized 18 files. Its acquisition suite returned five expected negatives, its nine source profiles and five profile negatives passed, and the exact lock validated 69/69 sources. Google profile selection remains six sources; document intelligence remains thirty.

The structural verifier returned:

```text
VERIFY_LIBRARY_PASS
packs=41 materialized_files=364 markdown_files=252
20/20 composition profiles PASS
```

The executable audit returned:

```text
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit
packs=41
upstream_sources=69
document_sdk_artifacts=5
business_central_artifacts=1
provider_adapters=7
```

It additionally materialized and exercised acquisition, source profiles, five SDK artifact gates, Business Central config/runtime negatives, MarkItDown local 5 tests, secure-file ingress 11 tests and provider adapter contracts. Strict field evaluation materialized and syntax-compiled in Audit mode; its 12 runtime tests remain tied to the exact acquired Stickler 0.6.0 frozen environment recorded in the prior V1 evidence and were not rerun against an unfrozen global Python.

After adding the two evidence records and updating the canonical audit, the final structural verifier was repeated and returned `VERIFY_LIBRARY_PASS`, `packs=41`, `materialized_files=364`, `markdown_files=253` and all twenty profile counts. The final executable audit repeated the same structural count and returned `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` with 69 upstream sources, five document SDK artifacts, one Business Central artifact and seven provider adapters.

## Honest state

The library improved by preventing a leader-company sample from being mistaken for production code. It did not invent a Google fix or claim perfect extraction. The 41 packs/364 files and previously admitted lanes remain reconstructible according to their own evidence; the rejected Google accelerator supplies architecture evidence only.

The global state remains `NOT_READY_UNDER_EXPANDED_USER_STANDARD`. Document ingestion is not complete until at least one full official provider lane passes secure ingress, exact artifact/dependency reconstruction, field-level correctness on the user's approved multi-class corpus, authentication, idempotency, load, recovery and controlled DEV deployment.
