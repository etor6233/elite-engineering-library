# Failure Learning System — 2026-08-25 V1

## Scope

This evidence records the authored failure-memory extension requested by the library owner. It does not attribute the method, ledger entries or fixes to any external company.

Added authority:

- `markdown_system/FAILURE_LEARNING_CONTRACT.md`;
- `markdown_system/PROJECT_FAILURE_LESSONS_TEMPLATE.md`;
- `markdown_system/LIBRARY_FAILURE_LEARNING_LEDGER.md`.

The shared ledger initially indexed 23 own-library failure events already corrected with regression evidence and 17 retained upstream failures or external conditions. During the next official-source probe, `LIB-FAIL-024` captured a nested PowerShell interpolation error before retry; the direct-shell correction must be demonstrated before promotion to `REGRESSION_PROVEN`. Detailed commands and results remain in their linked evidence files. No retained upstream failure was relabeled as a PASS.

## Routing and readiness

The mandatory project startup now creates `PROJECT_FAILURE_LESSONS.md` before product implementation. `AGENTS.md`, `AGENT_SYSTEM_START.md`, `START_ANY_PROJECT.md`, the readiness gate and readiness record require agents to:

1. capture failures, material warnings, skips and invalid assumptions in the same cycle;
2. consult prior fingerprints before repeating a known approach;
3. correct canonical sources and add regressions;
4. rebuild from an empty destination and rerun dependent gates;
5. preserve recurrence and residual risk instead of deleting history;
6. block `READY_TO_BUILD` when a relevant HIGH/CRITICAL failure remains open.

## Structural enforcement

`VERIFY_LIBRARY.ps1` now fails closed if:

- any failure-learning authority file is absent;
- the contract or project template loses a required diagnostic/correction/evidence field;
- the shared ledger has no IDs or contains a duplicate ID;
- any mandatory startup/readiness router stops routing `PROJECT_FAILURE_LESSONS.md`.

## Verification

Clean command:

```powershell
pwsh -NoProfile -File .\VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit
```

Observed before this evidence file was added:

```text
VERIFY_LIBRARY_PASS
packs=23 materialized_files=224 markdown_files=150
profile=ENTERPRISE_BACKEND_PACK_PLAN.md implementation_files=95
profile=ENTERPRISE_WEB_PACK_PLAN.md implementation_files=44
UPSTREAM_ACQUISITION_TEST_PASS
UPSTREAM_LOCK_VALID sources=46 selected=46
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=23 upstream_sources=46
```

Adding this evidence changes only the Markdown count. The complete audit must be repeated and its final count retained in current console evidence.

## Claim boundary

This proves persistent routing, schema enforcement, unique library failure IDs and continued structural/materialization integrity. It does not prove that a future agent will never make a mistake. It makes omission detectable, preserves the lesson and requires a regression before closure.
