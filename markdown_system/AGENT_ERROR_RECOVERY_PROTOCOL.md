# Agent Error Recovery Protocol

## Purpose

Errors are expected evidence during implementation. The agent must repair them autonomously when the correction is safe and inside scope; it must not stop at the first compiler, test, migration, security, performance or recovery failure.

## Mandatory loop

```text
capture exact failure and command
→ open or update its entry under FAILURE_LEARNING_CONTRACT.md
→ classify: contract | code | data | dependency | environment | operation | assumption
→ identify violated invariant and responsible layer
→ form the smallest differentiating hypothesis
→ reproduce with a focused test
→ correct the canonical Markdown/contract/source
→ add or strengthen regression coverage
→ regenerate manifests and hashes affected by the correction
→ reconstruct from a new empty destination
→ rerun focused gate
→ rerun every previously green dependent gate
→ record defect, correction and evidence
→ continue with the planned work
```

## Rules

- Preserve the original error and failed command in evidence; do not rewrite history as if the first attempt passed.
- Maintain `PROJECT_FAILURE_LESSONS.md` from `PROJECT_FAILURE_LESSONS_TEMPLATE.md`; a subsequent PASS closes the current incident only after regression evidence and never erases the failure.
- Consult `LIBRARY_FAILURE_LEARNING_LEDGER.md` before repeating a tool, source or lane with a known failure fingerprint.
- Correct the source of truth. A generated file may be used to diagnose, but a fix that is not returned to canonical Markdown will disappear on the next reconstruction.
- Change one causal layer at a time. Do not weaken an assertion, authorization check, invariant, type, migration or readiness gate merely to obtain green output.
- A retry without a hypothesis is not a correction. Repeated transient retries need a bounded budget and classification.
- After a local fix, rerun the focused regression first and then all gates that could have been invalidated.
- Never paste credentials, private customer payloads or unrestricted production telemetry into error evidence.
- If the error exposes corruption or ambiguous durable state, stop mutations, preserve evidence and follow the incident/recovery path.
- Ask the user only when correction requires a business/regulatory decision, credentials, external spend, destructive production action or another authority boundary listed in `AGENT_AUTONOMY_CONTRACT.md`.

## Upstream security exception

A vulnerability in admitted upstream code is not an ordinary trial-and-error defect. Reopen admission and apply `DEPENDENCY_UPDATE_CONTRACT.md`: contain exposure, identify the exact artifact and consumers, preserve the advisory/failure, consult official security/release/patch sources, and only then build a candidate. Prefer an official fixed release. If only an unreleased official patch exists, any pin or backport must preserve its exact upstream identity and declare the local result as `ADAPTED_PATCH` or `AUTHORED_PATCH`. If no official fix exists, keep production blocked or use an explicit containment/replacement/risk-decision path. Never silence a scanner or weaken a security assertion to call the incident fixed.

## Completion criteria

The agent may continue past the error only when one of these is true:

1. the canonical correction reconstructs cleanly and all affected gates pass;
2. a safe alternative with equal or stronger evidence replaces the failing component and the ADR/rollback path is recorded;
3. the capability is explicitly marked `BLOCKED` or `NONE_WITH_REASON` because an external decision is genuinely required.

An unexplained skip, disabled test, broadened permission, removed constraint or readiness downgrade is not recovery. A recurrence reopens the entry, increments its count and requires strengthening the preventive control.
