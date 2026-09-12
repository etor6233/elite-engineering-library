# All Implementation Packs Materialization — 2026-08-27 V26

## Governing delta

This record succeeds V25. Historical evidence remains a snapshot of its own source state.

The current official AWS GenAI IDP release remains `v0.6.5` at commit `1b5fd74454e593de233342a02ee911af8ee38359`. Official open PR #669 states that committed `v0.6.5-advverify` results referenced a suite that existed only in the author's working tree; those published results were not reproducible. Because the proposed suite targets `develop` and is not merged or released, Elite does not download it or treat it as a correction.

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` advances from 0.4.21 to 0.4.22. The exact release/archive/license/lock remain unchanged; the lock and AWS/document-intelligence profiles add the benchmark blocker. All thirteen dependent plans now request 0.4.22.

## Evidence and ledger

`AWS_GENAI_IDP_BENCHMARK_REPRODUCIBILITY_2026-08-27_V1.md` fixes the release, main, PR/base/head identities and upstream statement. `UP-FAIL-069` prevents the missing benchmark suite from being silently treated as product evidence. Existing dependency findings and typecheck errors remain open; this delta does not erase any earlier condition.

## Executed gates

The rebuilt 0.4.22 acquisition pack materialized 18 files. Acquisition negatives 5/5, source profile suite 9 valid/5 negative/1 positive and 69/69 lock validation passed. AWS still selects nine sources and document intelligence thirty.

Before adding this record, the structural gate returned:

```text
VERIFY_LIBRARY_PASS
packs=41 materialized_files=364 markdown_files=254
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

After this governing record and the canonical count update, the final structural verifier is repeated against 255 Markdown files.

## Honest state

The library has not gained a newly promoted end-to-end document platform. It has gained a stronger fail-closed decision: brand, test count and published benchmark directories cannot substitute for committed benchmark inputs. AWS GenAI IDP v0.6.5 remains `PINNED_CANDIDATE / CONDITIONED` and is not copied or deployed.

Global state remains `NOT_READY_UNDER_EXPANDED_USER_STANDARD`.
