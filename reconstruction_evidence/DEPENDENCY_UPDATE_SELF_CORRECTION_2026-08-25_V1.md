# Dependency Update and Agent Self-Correction — 2026-08-25 V1

## Scope

This evidence records the authored dependency-lifecycle and authority-freshness system plus the exact admission and execution of Google OSV-Scanner. The contracts are Elite Engineering Library policy; OSV-Scanner source/binary remains Google upstream code under Apache-2.0.

## New mandatory project artifacts

- `PROJECT_DEPENDENCY_UPDATE_RECORD.md` from `PROJECT_DEPENDENCY_UPDATE_RECORD_TEMPLATE.md`;
- `PROJECT_AUTHORITY_FRESHNESS_RECORD.md` from `PROJECT_AUTHORITY_FRESHNESS_TEMPLATE.md`.

The startup routers and readiness gate now require both artifacts together with readiness and failure learning. The dependency contract covers discovery, immutable candidate identity, license/notices/advisories, isolated lock resolution, baseline↔candidate gates, SBOM/provenance, canary, rollback and retained rejection evidence. The freshness contract covers temporal claims, primary-source hierarchy, contradictions, superseding records, explicit correction of prior agent statements and rerunning dependent gates.

`VERIFY_LIBRARY.ps1` fails closed if either contract/template disappears, mandatory terms are removed or any startup/readiness router stops routing the two project records.

## Google OSV-Scanner 2.5.1 source

- repository: `google/osv-scanner`;
- release: `v2.5.1`, published 2026-08-17;
- commit: `c84fa4568f2526d0333e9a914ea8a0a5f74ad68b`;
- source archive: 13,461,636 bytes;
- source SHA-256: `688ceb4ab62f38b5fd8fb4c49a9220e7589d49b88eee422aa3061ab6e339995a`;
- license: Apache-2.0, SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`;
- `go.sum` SHA-256: `660d54260464f85ae1800b9ee17ca8dbd440270e7cea3a98d2fd417e4a50b331`.

The canonical acquisition runner returned:

```text
UPSTREAM_ACQUISITION_TEST_PASS
UPSTREAM_SOURCE_ACQUIRED id=google-osv-scanner-2.5.1
```

The global lock is now 52 exact source archives.

## Official Windows binary

- `osv-scanner_windows_amd64.exe`: 58,970,112 bytes;
- SHA-256: `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`;
- official `osv-scanner_SHA256SUMS`: 554 bytes, SHA-256 `40120ab7d670199f6663d76f9806b6d11ab3c402df808e5cc4e13c5cb0643356`;
- published checksum entry equals the observed binary digest;
- `multiple.intoto.jsonl`: 24,443 bytes, SHA-256 `664273601759eb6540f35bfdbfa6ff083026c09e2e9bb7ad3cded795967dd462`.

Local execution returned:

```text
osv-scanner version: 2.5.1
osv-scalibr version: 0.5.2
commit: c84fa4568f2526d0333e9a914ea8a0a5f74ad68b
built at: 2026-08-17T03:44:26Z
```

The artifact lock preserves exact bytes and requires cryptographic provenance verification against an approved trust policy before product promotion; hashing the provenance file alone is not an authenticity proof.

## Actual web dependency scan

The canonical 3-pack/44-file web profile was reconstructed into an empty destination. OSV-Scanner inspected its `pnpm-lock.yaml` and reported:

```text
found 150 packages
OSV_EXIT=0
results: []
report SHA-256: 54345d5194d238b6017d0f3181f392109ce35023cc6b4bdfa8c134943dcb578e
```

This means no known vulnerability match was returned by the queried OSV data for detected packages at that time. It does not prove absence of unknown vulnerabilities, malware, unreachable/undetected components, compromised maintainer behavior or future advisories.

## Safety boundary

The `fix` command is prohibited by default. Official documentation warns it may execute package managers/scripts or follow project registries. Findings enter `DEPENDENCY_UPDATE_CONTRACT.md` as candidates; updates occur only in a disposable credential-free workspace and must pass all baseline/candidate gates before promotion.

## Clean library audit before this evidence file

```text
VERIFY_LIBRARY_PASS
packs=23 materialized_files=224 markdown_files=158
profile=ENTERPRISE_BACKEND_PACK_PLAN.md implementation_files=95
profile=ENTERPRISE_WEB_PACK_PLAN.md implementation_files=44
UPSTREAM_ACQUISITION_TEST_PASS
UPSTREAM_LOCK_VALID sources=52 selected=52
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=23 upstream_sources=52
```

Adding this evidence changes only the Markdown count; the final audit is repeated after canonical synchronization.
