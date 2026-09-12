# Human Review Leader Code Reaudit — V3

Date: 2026-08-28  
Scope: exact Google Cloud Document AI human-review sample and current service status  
Governing result: `HITL_LANE_REJECTED_DEPRECATED`; repository remains `SAMPLE_ONLY` for its other samples

## Exact official source

- Owner/repository: Google Cloud, `GoogleCloudPlatform/document-ai-samples`.
- Exact commit: `001ba391ab4a2f40d001cc0387618cb3c3699523`.
- Commit verification: GitHub verified, reason `valid`, dated 2026-01-05.
- Repository: active, Apache-2.0, no release identity used.
- Existing exact archive: 155,461,201 bytes; SHA-256 `3131ff0604685967932cad1b689a64f7e50af1dc6b45524d99c5389342d34e8f`.
- Exact notebook: `hitl-custom-review/hitl-custom-review.ipynb`, 21,181 bytes; SHA-256 `82c730377cb483885ac7b56e4ffeca885ca5217029c435a6cd4ff5bbc4e40fa0`.
- Official repository: https://github.com/GoogleCloudPlatform/document-ai-samples
- Official deprecation register: https://docs.cloud.google.com/document-ai/docs/deprecation
- Official request-human-review sample: https://docs.cloud.google.com/document-ai/docs/samples/documentai-review-document

## What the notebook actually implements

| Capability | Exact behavior | Result |
|---|---|---|
| Load processed document | downloads one JSON object from Cloud Storage and parses `documentai.Document` | present |
| Low-confidence routing | iterates entities and compares global or per-entity thresholds | present, first match only |
| Human-review request | creates `ReviewDocumentRequest` and returns operation name | present |
| Schema validation | explicitly sets `enable_schema_validation=False` | disabled |
| Final decision retrieval | never polls or loads reviewed output | absent |
| Authenticated reviewer | no reviewer identity is read or retained | absent |
| Mandatory reason | no reason field or validation | absent |
| Corrections/evidence | no corrected fields, revision comparison or decision receipt | absent |
| Lease/CAS/request binding | no assignment expiry, expected version or conflict handling | absent |
| Dual control | no independent second reviewer/quorum | absent |
| Automated suite | notebook only; no focal tests or locked graph | absent |

## Executed verification

- Current official repository metadata resolved to the same exact verified commit already present in the source lock.
- Four Python cells compiled independently: `[6, 8, 10, 12]`; notebook magic cell `[4]` was correctly excluded from Python compilation.
- Embedded execution output records an unversioned `%pip install` resolving `google-cloud-documentai 2.10.0`; it is historical notebook output, not a reproducible lock.
- The embedded example output says no entities were under the threshold, so it does not prove a completed human-review round trip.
- Rebuilt acquisition suite passes after adding the notebook artifact and rejection assertion: `UPSTREAM_ACQUISITION_TEST_PASS negatives=9`.
- Source-profile suite remains `valid=14 negatives=5 positives=1`; Google selects 7 and document intelligence 41.
- No Google Cloud account, billing, processor call or deprecated HITL request was executed.

## Current-service decision

Google's official deprecation register lists Human in the Loop as deprecated on 2024-01-16. Therefore, a syntactically valid request sample is not a current production architecture. The library retains the broader repository for other exact samples but rejects this lane for any new implementation. The agent must ask the project to select a supported authenticated review workflow and prove identity, mandatory reason, corrections, expiry, CAS, durable evidence, dual control where required, recovery, security, load and cost before automatic storage can be enabled.

## Truth boundary

No Google code was rewritten or represented as Elite-authored. No replacement workflow was invented. The absence of a complete current official one-piece implementation remains an explicit blocker rather than an invitation to silently compose incomplete samples.

