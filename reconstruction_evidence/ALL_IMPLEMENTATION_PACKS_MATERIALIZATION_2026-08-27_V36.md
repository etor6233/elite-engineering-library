# All Implementation Packs — Materialization Snapshot V36

Date: 2026-08-27. This snapshot supersedes V35 for current counts and the Azure retained-evidence lane; earlier evidence remains historical.

## Current verified corpus

- implementation packs: 44;
- materializable files across packs: 388;
- Markdown sources after this snapshot: 279;
- canonical PowerShell scripts: 8;
- composition profiles: 21;
- exact official source archives in the acquisition lock: 73;
- provider adapters reported by the executable audit: 9;
- provenance across materialization blocks: 383 `AUTHORED`, 2 `ADAPTED`, 3 `VERBATIM`.

`VERIFY_LIBRARY.ps1` passes manifest↔block↔SHA equality for all packs and composes every profile. The Azure document profile now composes acquisition, secure local ingestion, official SDK artifacts, Azure Content Understanding runtime, strict field evaluation and Azure Blob immutable evidence as 6 packs / 50 files. Google remains 6/50, AWS document 5/45, backend 16/95 and web/BFF 3/44.

## New executable Azure retained-evidence primitive

`PYTHON-AZURE-BLOB-IMMUTABLE-EVIDENCE-ADAPTER 0.1.0` adds nine files around Microsoft `azure-storage-blob 12.30.0`. Exact wheel/sdist/archive hashes, MIT license/notice, unsigned tag condition, 79/79 wheel↔sdist and 193/193 sdist↔commit file equality are retained. The 12-wheel runtime lock installed with required hashes; OSV returned zero findings; 15 tests, compileall and pip-check passed from a clean Markdown materialization.

The adapter requires a private versioned immutable container and exact non-overridable encryption scope before it can write. It binds fresh authority to the configured target, requires access/cost/retention/irreversibility approval, creates only with locked WORM and optional separately approved legal hold, validates MD5 transport plus SHA-256 evidence, and verifies the exact returned version. Unknown or partially verified effects are retained for reconciliation and never retried blindly. It does not configure Azure infrastructure.

## Failure learning and global audit

The ledger retains 285 local failures and 80 upstream conditions. The new Azure work preserves wrong-reference, enum normalization, OSV parsing, PowerShell array invocation, stale process status, hardcoded count, wrong-object assertion, duplicate failure-ID narrative, provider materialization order, evidence-name, Windows-glob, provenance-syntax and embedded-fixture counting failures. Each corrective rule was exercised; the unsigned Microsoft tag remains explicit.

`VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit` passes with exact CPython 3.12.13 and Go 1.26.7 paths. In network-authorized audit mode it creates a clean Python environment, installs the exact Azure lock, runs pip-check, compileall and 15 tests, while the global runner continues to execute acquisition, composition, document, security and provider gates. No provider sandbox, credential, paid call, storage mutation or production authorization is inferred.

## Honest readiness

Status remains `NOT_READY_UNDER_EXPANDED_USER_STANDARD`. AWS, Google and Azure now have executable retained-evidence storage primitives, but truly local immutable storage, cross-stage document orchestration, human-review workflow, business persistence handoff, end-to-end duplicate/partial-failure tests, target load, recovery and real-provider evidence remain unfinished. No portable ZIP was created.
