# All implementation packs — materialization evidence V28

Date: 2026-08-27  
Governing status: `NOT_READY_UNDER_EXPANDED_USER_STANDARD`  
Supersedes: V27 for the current library source tree. Historical records remain evidence of their dated snapshots.

## Delta admitted in this snapshot

The official Python implementation of Microsoft Azure Content Understanding is no longer represented only by a PyPI artifact and tutorial/sample. `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.24` fixes the signed GA tag/commit `129a5cbb…`, a 191429136-byte source archive, root/package license and notice, package pyproject/assets and service CI hashes. Actual acquisition through the rebuilt no-Git runner produced 52748 files and an exact receipt. The lock now contains 70 sources; Microsoft/document profiles select 13/31.

The exact wheel/sdist remain Microsoft artifacts. All 114 sdist files corresponding to signed repository package files were byte-identical. The exact wheel passed 183 isolated official unit tests on CPython 3.12.13. The full 268-test source collection retains `UP-FAIL-074`: Microsoft test tooling requires a `.git` checkout plus test-proxy assets before any test executes. No live Azure service, model or corpus was claimed.

`AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME 0.1.2` remains visibly authored glue. It now derives analyzers only from `REQUIRED` profile classes, requires an `ADMITTED` local security receipt matching SHA/bytes/MIME, rejects provider analyzer/API drift, hash-links the complete evidence chain and keeps storage unauthorized. It adds an exact nine-wheel CPython 3.12/Windows x86-64 lock. The global executable runner materializes the pack on every audit, reports an explicit offline skip, and with network plus the exact target creates an empty venv, installs with `--require-hashes`, passes `pip check` and runs five tests.

`AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md` now composes five layers/41 files: official source acquisition, secure local ingestion, SDK artifact acquisition, Azure runtime and strict field evaluation. This makes the code path tangible without pretending that account, corpus, custom analyzers or business storage are complete.

Detailed evidence: `MICROSOFT_AZURE_CONTENT_UNDERSTANDING_PYTHON_2026-08-27_V1.md`.

## Reproducible gates

```text
VERIFY_LIBRARY_PASS
packs=41 materialized_files=365
profile=AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md implementation_files=41

UPSTREAM_LOCK_VALID sources=70 selected=70
UPSTREAM_ACQUISITION_TEST_PASS negatives=6
SOURCE_PROFILE_VALID profile=microsoft-secure-document-pipeline selected=13
SOURCE_PROFILE_VALID profile=document-intelligence-leaders selected=31
SOURCE_PROFILE_TEST_PASS valid=9 negatives=5 positives=1

official Microsoft unit tests: 183 PASS
Azure authored wrapper tests: 5 PASS
product graph: 9 packages, pip check PASS, OSV 0 known records
actual source acquisition: 52748 files and all recorded package/CI hashes PASS

VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit
```

The executable audit was run twice after the runner change: offline, where the Azure install/test stage was explicitly marked skipped, and with `-AllowNetwork -PythonExecutable <CPython 3.12.13>`, where the new frozen environment and five tests passed. No provider sandbox was contacted.

## Failure memory

The ledger retains 228 local failures and 74 upstream conditions. This cycle added `LIB-FAIL-218..228` and `UP-FAIL-074`; it did not erase corrected mistakes. The final upstream admission is a candidate, not a promotion to `REUSABLE_PACK` or automatic storage.

## Remaining blockers

- no complete target blueprint/authority/access record for a real company system;
- no Azure tenant/resource/identity/region/quota/model/cost approval;
- no authorized target document corpus, independent ground truth or custom analyzer for proforma, packing list, BOL, customs and the remaining enterprise document classes;
- no live field-level proof, review operation, idempotent storage, reconciliation, load, privacy, offensive security, backup/recovery, canary or rollback evidence for that target;
- other implementation packs keep their existing project and upstream conditions.

Therefore the library is structurally coherent and materially stronger, but it is not declared complete under the user's expanded whole-enterprise standard.
