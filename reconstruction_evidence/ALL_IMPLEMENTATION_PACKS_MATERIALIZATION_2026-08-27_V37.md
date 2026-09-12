# All Implementation Packs — Materialization Snapshot V37

Date: 2026-08-27. This snapshot supersedes V36 for current corpus counts and document orchestration. Earlier snapshots remain historical evidence of their exact state.

## Current verified corpus

- implementation packs: 45;
- materializable files across packs: 401;
- Markdown sources after this snapshot: 283;
- canonical PowerShell scripts: 8;
- composition profiles: 22;
- exact official source archives in the acquisition lock: 75;
- official source profiles: nine;
- provider adapters reported by the executable audit: nine;
- document orchestrators reported by the executable audit: one;
- provenance across exact materialization blocks: 394 `AUTHORED`, three `ADAPTED`, four `VERBATIM`.

`VERIFY_LIBRARY.ps1` passes manifest↔block↔SHA equality for every pack and composes all 22 governed profiles. Current key profiles are initialization 2/24, routing 2/22, durable document control plane 5/46, Azure document 7/63, Google document 7/63, AWS document 6/58, MarkItDown local 4/35, backend 16/95 and web/BFF 3/44.

## New durable document orchestration

`MICROSOFT-DURABLE-DOCUMENT-ORCHESTRATION 0.1.0` adds 13 files around Microsoft `durabletask 1.9.0` and the official Azure Durable Task Python agent skill. The exact signed commits, commit archives, PyPI wheel/sdist, licenses and skill hash are fixed. The official core suite produced 1,034 PASS, one skip and 14 deselected; the frozen six-wheel product graph returned zero OSV findings on the dated query.

The local glue is explicitly `AUTHORED`; the official skill is byte-`VERBATIM`; its license is `ADAPTED_FINAL_NEWLINE_ONLY`. A clean hash-required install passed `pip check`, compileall and 13 orchestration tests. The profile authorizes no effect by default. Ambiguous storage and business effects require reconciliation rather than blind retry; approval is result-hash bound; a final-evidence failure after commit preserves the transaction as `PARTIAL_EFFECT`.

The acquisition core is now V0.4.28 with 75 exact sources. Microsoft document selection is 16, document intelligence 35 and enterprise platform 12. The source/acquirer suite retains six negative gates; source profiles retain nine valid, five negative and one positive approval case.

## Failure learning and global audit

The ledger retains 321 local failures and 85 upstream conditions. This tranche preserves inferred artifact URLs, wrong license/path assumptions, indirect tag archive identity, missing optional-artifact hash, path/root and materializer/compositor mapping, three host/shell transport failures, one missing Linux fault-test prerequisite plus its initially incomplete runtime closure, incorrect materializer/verifier/workdir invocations, one pack reconstruction mismatch, two policy-rejected cleanup attempts, partial inventory paths, provenance parser overcount, corrected copied identities and one corrected Tessera example path. It also preserves Tessera's unsigned release identity, Windows/POSIX incompatibility and OSV-positive published graph. Corrective regressions bind exact URLs, path/hash pairs, commit archive identity, composition records, static contracts, fence-aware provenance counts and declared interfaces.

`VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit` passes without network with exact CPython 3.12.13 and Go 1.26.7 paths, reporting 45 packs, 75 sources, five document SDK artifacts, one Business Central artifact, nine provider adapters and one document orchestrator. In network-authorized Audit mode it creates clean environments, installs only hash-locked public wheels, and passes the 13 new orchestration tests plus the existing MarkItDown, Azure Blob, Azure Content Understanding, Google Document AI, routing, security and provider gates. No provider sandbox, credential, paid call, storage mutation or production authorization is inferred.

## Honest readiness

Status remains `NOT_READY_UNDER_EXPANDED_USER_STANDARD`. The reusable control plane and three cloud lanes now include cross-stage durable orchestration, but project-specific business adapters, real corpus/ground truth, human-review identity/UI, exact production DTS backend, cloud accounts/IAM/RBAC/KMS/retention, live duplicate/concurrency/load/soak, recovery and production rollout remain target evidence. The Microsoft event-loop warning and datetime deprecations remain open upstream. No portable ZIP was created.
