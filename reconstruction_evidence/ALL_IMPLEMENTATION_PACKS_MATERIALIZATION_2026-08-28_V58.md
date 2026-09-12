# All Implementation Packs Materialization — V58

Date: 2026-08-28  
Scope: governing V64 source-lock and Microsoft Durable Task JS human-review reaudit

## Canonical state

- 50 implementation packs.
- 482 materializable files.
- `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.49.
- 97 exact official source archives.
- 14 fail-closed source profiles.
- Microsoft source profile: 17 sources.
- Document-intelligence source profile: 42 sources.
- Microsoft Durable Task JS v0.4.0 is exact/hash-locked as `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION` for its human-review examples.

## Microsoft source result

- Archive: 1,372,198 bytes; SHA-256 `9b0b01aa6aab53f5865897ccf6dec9ab98dd847e358a147466e2d237208a9ba5`.
- License: MIT; SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`.
- Stable v0.4.0 is a lightweight tag over unsigned commit `051bd5f4bb8a6a4c3408865625847abfc4659115`; signed current `main` is newer and was not substituted.
- npm 10.9.2 exact-lock install PASS; complete build PASS; lint PASS; 98/98 suites and 1,503/1,503 unit tests PASS.
- npm 11.19.0 rejects the same lock; exact npm audit retains 8 findings — 4 high, 3 moderate, 1 low.
- Neither human-interaction example is covered by the official unit tests.
- The managed example supplies only `approved:boolean`; the classic example trusts a client-supplied name and omits `yield` before `whenAny`.
- Neither example supplies authenticated human identity, mandatory reason, corrections, lease, CAS, request-bound evidence or dual control.
- No Docker emulator, Azure account, billed service or cloud deployment was used.

## Canonical reconstruction

The source pack was materialized into a new empty temporary directory. All four changed artifacts matched the audited source tree byte for byte.

```text
V64_REBUILD_PASS files=4 materialized=23
UPSTREAM_ACQUISITION_TEST_PASS negatives=9
SOURCE_PROFILE_TEST_PASS valid=14 negatives=5 positives=1
UPSTREAM_LOCK_VALID sources=97 selected=97
SOURCE_PROFILE_VALID profile=microsoft-secure-document-pipeline selected=17
SOURCE_PROFILE_VALID profile=document-intelligence-leaders selected=42
```

## Global gates

```text
VERIFY_LIBRARY_PASS
packs=50 materialized_files=482 markdown_files=337

VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=50 upstream_sources=97 document_sdk_artifacts=10 official_invoice_samples=1 official_google_process_samples=1 official_dapr_outboxes=1 official_pg_durable_handoffs=1 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1
EXEC_AUDIT_EXIT=0
```

The `markdown_files=337` value was measured before adding this evidence file. A final structural verification after insertion governs the distribution count.

## Admission boundary

The Microsoft runtime is real and materially tested, but the user-requested enterprise human review remains unresolved. Elite did not patch either demo or represent authored composition as Microsoft code. No official one-piece current source proves authenticated identity, mandatory reason, corrections, lease/expiry, optimistic conflict, request binding, durable evidence and dual control. The project must still provide authorized corpus, accounts, schemas, policies and live operational gates.
