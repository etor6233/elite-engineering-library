# Go AWS Enterprise Storage and Email Adapters — 2026-08-27 V2

Status: `REBUILD_VERIFIED / CONDITIONED`

## Official authority and admission decision

The executable dependency remains Amazon's official AWS SDK for Go v2: S3 `v1.107.3`, SES v2 `v1.67.0` and config `v1.32.38`, under Apache-2.0. The adapter is local `AUTHORED` glue and is not presented as Amazon source.

AWS's current documentation states that Object Lock requires versioning, implements WORM retention, and cannot be disabled after it is enabled on a bucket. AWS also documents that Batch Operations job identities use `s3:CreateJob`, `s3:DescribeJob`, `iam:PassRole`, manifest/report permissions and job resources. Governing sources:

- https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock-configure.html
- https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock-managing.html
- https://docs.aws.amazon.com/AmazonS3/latest/userguide/batch-ops-iam-role-policies.html
- https://docs.aws.amazon.com/AmazonS3/latest/userguide/security_iam_service-with-iam.html

The signed official AWS sample `aws-samples/sample-S3-Object-Lock-Automation` at commit `66e891d9f08af2937e6627c6825569195d1f07e4` was audited, not adopted. Exact archive: 28,287 bytes, SHA-256 `23bd8bcec943b4b7b79936aec576f0493bdc6a65c57e2e751e9d71f90ff09a9f`; MIT-0 LICENSE SHA-256 `5025c7bbedbc8da868b6dce2ab689225b3f33c43b3f03368b0a331809b6ada26`. GitHub reports the commit signature `verified=true`.

It is `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`: `upload_objects.sh` has a Python shebang and Python body and fails `bash -n` at line 25; the remaining three scripts pass syntax but use only `set -e`, a predictable `/tmp/iam-policy.json`, an unvalidated bucket substitution, fixed IAM role, permanent bucket mutations, interactive/no-confirmation batch jobs and no automated suite. The README's `s3control:*` IAM examples diverge from current official AWS `s3:*` policy actions. No sample file was copied or executed against AWS.

## V2 adapter behavior

`PutImmutableObject` now:

1. validates the exact expected 12-digit bucket owner, canonical key, bytes, MIME, encryption, retention mode and future timestamp;
2. calls official `GetObjectLockConfiguration` and refuses upload unless Object Lock is proven enabled;
3. calls official `PutObject` with SHA-256 checksum, `If-None-Match: *`, encryption, expected owner and explicit `GOVERNANCE|COMPLIANCE` retention;
4. requires ETag and VersionId;
5. calls official `GetObjectRetention` for that exact version and compares mode/date;
6. emits `retention_verified=true` only after equality;
7. if verification after the remote write fails, returns the created-version receipt together with the error so reconciliation can record the effect and must not retry blindly.

The adapter never enables Object Lock and never applies legal hold. Those remain separate, explicit, target-authorized infrastructure operations.

## Reproducible gates

The pack reconstructs nine files with exact manifest/SHA correspondence. The toolchain was the exact official `go1.26.7.windows-amd64.zip`, 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`.

```text
go mod verify
all modules verified
go test -count=1 -v ./...
6 PASS
go vet ./...
PASS
go build ./...
PASS
```

Tests cover checksum/encryption/no-overwrite, owner/lock preflight, exact retention verification, rejected disabled/unreachable Object Lock, visible partial provider effect and SES safe receipt/negative inputs.

## Remaining conditions

No AWS account, bucket, IAM policy, KMS key, invoice corpus or paid service was used. Production still requires project-approved owner/bucket/region/IAM/cost, existing Object Lock and versioning, retention/legal policy, malware/quarantine, lifecycle/restore, provider receipts, load, observability and reconciliation. Storage of an immutable original/evidence object never authorizes persistence of extracted business fields.
