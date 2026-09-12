# All Implementation Packs Materialization — 2026-08-26 V24

## Governing delta

This record succeeds V23 for the current source tree. Historical records remain immutable evidence of their own snapshots.

The exact AWS sample `aws-samples/aws-textract-document-data-extraction-platform` at commit `c7fed5353abcae184fbd705d5cf6a8512af0e9c4` was rebuilt more deeply with supported, isolated toolchains. Its installation rewrites three Poetry locks, its official Windows and Linux build/test paths do not complete, unversioned Projen/codegen resolution escapes the root lock, and the original dependency locks contain critical/high known vulnerabilities. It is now explicitly `SAMPLE_ONLY / REJECTED_FOR_IMMEDIATE_ADOPTION`; no local product patch was added.

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` advances from 0.4.19 to 0.4.20. Its AWS and document-intelligence source profiles now expose those blockers before approval or acquisition, so the agent can skip the rejected sample instead of spending project time discovering the same failures. The immutable source identity, license and archive hash remain available for evidence.

## Exact evidence

`AWS_TEXTRACT_EXTRACTION_PLATFORM_REBUILD_2026-08-26_V1.md` records:

- archive 1.124.136 bytes/SHA-256 `52755b28d8afd2c8cee23ecba76b11c364edd3ecad066b2e4a51eb2f77219bb1` and pnpm lock SHA-256 `240755c9c46f7d8e41985f50d141c7f2f93d15bbae461c595f21cd655b147e22`;
- Node 18.20.8, pnpm 8.15.9, Python 3.11.9, Poetry 1.8.3, Corretto 17.0.20.10.1 and OSV-Scanner 2.5.1;
- install exit 0 but three mutated Python locks with different Windows/Linux resolutions;
- build exit 1 on Windows and Linux; test exit 1 on Linux;
- `pnpm audit --prod`: 1.485 production dependencies and 153 vulnerability occurrences, including 5 critical and 65 high;
- OSV over the four original locks: 2.434 detected packages, 91 affected records and 218 unique vulnerability IDs.

The failure ledger now retains 182 local failures and 58 upstream conditions. The new local failures are closed with demonstrated regressions; the five new upstream conditions remain open or externally blocked and prevent promotion.

## Executed gates

Before this record was added, the structural verifier returned:

```text
VERIFY_LIBRARY_PASS
packs=41 materialized_files=364 markdown_files=250
20/20 composition profiles PASS
```

The rebuilt acquisition pack materialized 18 files. Its exact 69-source lock, all nine source profiles, five profile negatives, one positive offline case and five acquisition negatives passed. The AWS profile selected nine IDs and document intelligence selected thirty while carrying the new fail-closed blockers.

The full executable audit returned:

```text
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit
packs=41
upstream_sources=69
document_sdk_artifacts=5
business_central_artifacts=1
provider_adapters=7
```

The canonical source tree contains 251 Markdown files after adding this record. A final structural and executable audit is repeated against that exact count before handoff.

## Honest state

The library improved by rejecting a leader-company sample that does not satisfy the user's current standard, not by disguising it as ready code. Other admitted sources and materialized packs keep their prior classifications. The global state remains `NOT_READY_UNDER_EXPANDED_USER_STANDARD`: official code is available immediately only where its exact lane is admitted; a missing or rejected lane remains blocked until an official corrected source passes the same evidence chain.
