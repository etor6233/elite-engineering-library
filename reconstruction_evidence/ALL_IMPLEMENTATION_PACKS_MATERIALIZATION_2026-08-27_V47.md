# All implementation packs — materialization evidence V47

## Governing snapshot

Date: `2026-08-27`.

This V47 supersedes V46 for the current library state. Historical evidence remains immutable.

| Metric | V47 |
|---|---:|
| implementation packs | 46 |
| materializable files | 414 |
| composition plans | 23 |
| official source profiles | 11 |
| upstream sources in exact lock | 83 |
| Markdown files after this evidence | 305 |
| PowerShell scripts | 8 |
| retained local failures | 375 |
| retained upstream conditions | 113 |
| block provenance | 402 `AUTHORED`, 3 `ADAPTED`, 9 `VERBATIM` |

No source code from the new AWS audit was embedded as an authored pack block. The 414-file materialization/provenance total is unchanged.

## Delta from V46

- `OFFICIAL_UPSTREAM_ACQUISITION_CORE` advanced from `0.4.37` to `0.4.38`.
- Exact source lock advanced from 82 to 83 entries.
- Added `aws-transactional-outbox-pattern-23e0519` at verified signed commit `23e0519c7a8048c6db4acd8f3d7de534a9d5deac`, MIT-0, with archive/license/README/lock/manifest and 11 additional artifact hashes.
- The new source is `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`, absent from all automatic source profiles.
- All 18 composition plans that select the acquisition pack now require `0.4.38`.
- Added four upstream conditions `UP-FAIL-110` through `UP-FAIL-113` for coverage, correctness/failure handling, dependency security and IaC security/recovery.
- Added local failure lessons `LIB-FAIL-367` through `LIB-FAIL-375`; browser/API routing, non-hermetic tests, HTTP-client ambiguity, three vulnerable graphs, JSON empty-array counting, exception-class precision and path discovery are retained rather than erased.
- Added exact technical evidence in `AWS_TRANSACTIONAL_OUTBOX_REAUDIT_2026-08-27_V1.md` and updated readiness, capability, admission, notices and implementation indexes.

## Exact AWS verdict

AWS Prescriptive Guidance remains an accepted official architecture reference for transactional outbox. The corresponding sample repository is not admitted as ready-to-copy code:

- both Java modules compile/package only with tests skipped;
- relational `contextLoads` requires an external PostgreSQL instance and has no behavior assertions;
- CDC `contextLoads` fails on two SDK HTTP implementations;
- the CDK Jest test reports one pass while all assertions/imports are commented;
- relational partial SQS batch failures are not inspected before deleting all selected rows;
- CDC conversion omits the ID required by the sender, absorbs serialization `IOException`, and imports an SDK-v1 service exception beside an SDK-v2 client;
- infrastructure uses broad IAM and destructive/deletion-unprotected defaults;
- OSV reports 28 unique npm advisories, 117 relational Maven advisories and 139 CDC Maven advisories.

This rejects weak code without inventing a replacement or attributing local glue to AWS. `UP-FAIL-109` remains open: a maintained, proven durable outbox/reconciler is still required before the AWS approval source can provide a complete decision audit.

## Canonical reconstruction gates

Fresh materialization of `OFFICIAL_UPSTREAM_ACQUISITION_CORE.md` produced 20 files.

```text
UPSTREAM_ACQUISITION_TEST_PASS negatives=7
UPSTREAM_LOCK_VALID sources=83 selected=83
SOURCE_PROFILE_TEST_PASS valid=11 negatives=5 positives=1
V47_ACQUISITION_CANONICAL_PASS sources=83 artifacts=11
```

Every existing profile retained its exact selection count: secure file 6, Linux archive 1, AWS documents 10, Google documents 7, Microsoft documents 16, document intelligence 35, strict field 1, commerce 25, enterprise 12, Tessera 1 and AWS authenticated approval 1.

Global structure/composition gate:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=414 markdown_files=305
```

The final rerun after adding this V47 evidence measured the 305-file Markdown corpus.

Global executable audit:

```text
VERIFY_EXECUTABLE_LIBRARY_PASS
mode=Audit
packs=46
upstream_sources=83
document_sdk_artifacts=5
business_central_artifacts=1
provider_adapters=9
document_orchestrators=1
evidence_logs=1
```

The audit preserves explicit skips for network/account/Linux/toolchain-dependent runtime lanes; it does not reinterpret them as production PASS.

The exact temporary audit root was resolved as one direct child of the Windows temp directory, inventoried at 11,417 files / 1,030,161,009 bytes, deleted with the .NET directory API and confirmed absent: `V47_TEMP_CLEANUP_PASS targets=1`.

## Current boundary

The library remains consistent and executable at its documented scope, but is not declared universally complete. Specifically:

- AWS authenticated approval remains conditioned because its audit event is written after the decision transaction.
- The audited AWS transactional-outbox sample cannot safely close that gap.
- No AWS account, PostgreSQL service, Kinesis, DynamoDB, SQS, ECS, WAF, PITR/restore, replay or live failure injection was authorized in this audit.
- A future source/candidate must prove atomic domain+outbox write, per-entry acknowledgment, concurrency control, idempotency, redelivery/DLQ/replay, reconciliation, immutable receipt, least privilege, restore and clean dependencies.

V47 therefore improves the agent's speed and safety by preventing automatic adoption of an official but insufficient sample. It does not add unproven production code.
