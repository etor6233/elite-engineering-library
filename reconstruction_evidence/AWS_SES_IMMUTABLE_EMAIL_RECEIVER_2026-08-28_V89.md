# AWS SES Immutable Email Receiver — Reconstruction Evidence V89

## Scope and provenance

- Pack: `AWS-SES-IMMUTABLE-EMAIL-RECEIVER` 0.1.0.
- Pack SHA-256: `239a8194b8ed8b04ce0c87d76a8dea938b9df57e7373691cb930105e900dfe81`.
- Nine `AUTHORED` integration files. They are not attributed to AWS.
- Runtime dependency: official Apache-2.0 Boto3/Botocore 1.43.83 graph from PyPI, seven exact wheels with SHA-256 in `requirements.txt`.
- Contract authorities: official AWS SES receipt/S3 action, S3 event ordering/duplicates/structure, SQS partial batch and visibility timeout, S3 Object Lock, Lambda Python dependency packaging and CloudFormation references, verified 2026-08-28.

## Materialized identity

| File | Bytes | SHA-256 |
|---|---:|---|
| `build.ps1` | 3218 | `339692a57a15f8ebe864e1d4f5a6f17abbca9c0aaf017e79414ade4328650a97` |
| `deploy.ps1` | 5043 | `2e1893e326c719518007e2ad619fbdb4074794b735206b21c5e8490c0816ad95` |
| `handler.py` | 24351 | `e87642aacfb16d23bcb3b760f44df5b9395cb30cb46358fa3618d4c0f9e60859` |
| `provider-profile.template.json` | 523 | `9e1bd96e9a29e760826e5c806a92596184a72ede2c54d070ddf4fac2f9d40091` |
| `README.md` | 3217 | `58941a53bb7809b8787ca2b09e41d52e693dea5a4607df99a65b948f5248b349` |
| `requirements.txt` | 719 | `95c37a7ffe3d5e3e3dc9c188d4370a934e7d853c2c3494c99cc7884e54df24cc` |
| `template.yaml` | 14465 | `9a74bf3abefa56e58ad67adc752181e15f95e7e0ce82dedaf8502d2ca6593381` |
| `test_handler.py` | 14383 | `1c42d9766421e0ad543442afe660812dfa58b67ffd3df1dd534b9407b805113f` |
| `verify_contract.ps1` | 5255 | `d802546d46b0b2ed4121c6b507e8fca61b380e1e9ad9a34e73f2dd3e09351a96` |

## Executed evidence

- Canonical materializer: 9/9 files in a new directory; declared and reconstructed hashes matched.
- Python AST: PASS.
- Receiver unit/negative suite: 11/11 PASS.
- Existing MIME core suite: 18/18 PASS.
- Failure paths exercised: wrong bucket/prefix/version/size; virus verdict; missing Object Lock; wrong manifest authority; busy and expired leases; stale-owner completion; duplicate event; post-retention attachment outage; per-record SQS retry.
- Build with `--require-hashes --only-binary=:all:`: seven packages installed, Boto3/Botocore imports both reported 1.43.83, artifact receipt covered 2,191 files.
- Dated OSV batch query: seven empty result objects, zero records. This is not a future guarantee.
- cfn-lint 1.55.1: exit 0 after correcting the `ScalingConfig.MaximumConcurrency` minimum to 2.
- `VERIFY_LIBRARY_PASS`: 63 packs, 594 files, 406 Markdown before this evidence; email profile 43 files.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: `aws_ses_immutable_email_receivers=1`, 63 packs and 111 upstream sources.
- Failure memory: 998 local + 175 upstream = 1,173 unique rows before this evidence; zero open local failures.

## Proven behavior and non-claims

The reconstructed path binds one tenant to a receipt rule/prefix, fetches the exact S3 VersionId, verifies owner/etag/size, uses owner-fenced DynamoDB leases, retains the raw MIME content-addressed with checksum/KMS/Object Lock and verifies the created version retention, extracts through the bounded MIME core, uploads hash-named quarantine objects and commits the manifest last. Durable manifest detection prevents re-execution after DynamoDB TTL; retries reconcile already-created retained/quarantine objects instead of overwriting them. Logs use event/message hashes and error codes, not addresses or content.

The template/build/tests do not prove an AWS account, SES identity, MX, quota, cost, target KMS/IAM, live delivery, DLQ replay, canary, security engines or document-field accuracy. The fail-closed profile enables none of those. `deploy.ps1` requires completed ownership fields, live-effects approval in the profile and on the command line, STS account match, SES identity and exact regional MX before deployment; then it activates and verifies the receipt rule set. Attachments remain `PENDING` with `automatic_storage_authorized=false` until the separately composed security, extraction, evaluation and review gates pass.
