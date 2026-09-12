# Official Document Upstream Freshness — 2026-08-26 V1

## Scope

Read-only revalidation of the exact official document-processing sources already present in `OFFICIAL-UPSTREAM-ACQUISITION-CORE`. No Git repository was created, no branch source was copied, no dependency was installed and no lock entry was promoted.

## Primary-source result

GitHub's official repository and commit APIs returned:

| Official repository | Default branch head at audit | Locked identity | Result |
|---|---|---|---|
| `Azure-Samples/azure-ai-content-understanding-python` | `fbc18880179397729abfa9506bda865d009dc711` | same commit; MIT | `EXACT_HEAD` |
| `Azure-Samples/data-extraction-using-azure-content-understanding` | `0461a5549104ca769b8ec082c05894468997f300` | same commit; MIT | `EXACT_HEAD` |
| `GoogleCloudPlatform/document-ai-samples` | `001ba391ab4a2f40d001cc0387618cb3c3699523` | same commit; Apache-2.0 | `EXACT_HEAD` |
| `aws-samples/aws-textract-document-data-extraction-platform` | `c7fed5353abcae184fbd705d5cf6a8512af0e9c4` | same commit; Apache-2.0 | `EXACT_HEAD` |
| `aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws` | `720f3052a57c22ae0d1e50886965c065ee86d60b` | tag `v0.6.5` → `1b5fd74454e593de233342a02ee911af8ee38359` | `LATEST_TAG_STILL_LOCKED`; branch not promoted |
| `aws-samples/amazon-textract-textractor` | `c0fe379340108c1992944303d102c49aae560152` | tag `v1.10.0` → `8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92` | `LATEST_TAG_STILL_LOCKED`; branch not promoted |

All six repositories reported `archived=false`. The expected repository licenses were returned by the official metadata API. The Azure Content Understanding samples also state that production applications should use the GA SDK; Elite already pins and verifies that official `1.1.0` SDK artifact separately instead of misclassifying the REST sample wrapper as production runtime.

## Release rule retained

A newer default-branch commit is discovery evidence, not admission. Textractor remains fixed to its latest published tag. AWS GenAI IDP remains fixed to its latest tag and retains its recorded OSV/npm/typecheck blockers. Neither branch snapshot inherits a release's license graph, tests, artifact identity or project approval.

## Failure learning

The first PowerShell query failed before network access because a `foreach` statement was piped directly into `ConvertTo-Json`, producing an empty pipeline element parser error. `LIB-FAIL-169` records it. The corrected command accumulated typed objects first and serialized the completed list; it returned all six repositories and was rerun without modifying source state.

## Outcome

The official document sources selected for immediate acquisition are not stale as of this audit. This does not authorize deployment or automatic storage: real provider access, representative project documents, closed schemas, independent ground truth, field-level evaluation, privacy, human review, load, recovery and cost controls remain project evidence.
