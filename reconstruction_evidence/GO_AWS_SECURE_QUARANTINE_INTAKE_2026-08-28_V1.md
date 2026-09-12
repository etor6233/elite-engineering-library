# Go AWS Secure Quarantine Intake — Reconstruction Evidence V1

## Scope

Evidence for `GO-AWS-SECURE-QUARANTINE-INTAKE` 0.1.0 and `AWS_SECURE_DOCUMENT_INTAKE_PACK_PLAN.md`. This proves offline reconstruction and local contracts only. It does not claim an AWS account call, production IAM, live CORS, bucket policy, KMS, event delivery, malware scan, extraction accuracy or business-storage authority.

## Authorities

- AWS SDK for Go v2 S3 `v1.107.3`, official module source under commit `284a4846e7bb941926e2d15144b17c01ffc98c1d`;
- AWS official Go v2 presigning example and S3 upload documentation;
- AWS S3 SigV4 policy keys and Prescriptive Guidance for `s3:signatureAge`/temporary credentials;
- AWS S3 documentation states event delivery is at least once and may be duplicated/out of order, so event reconciliation remains a separate blocked capability;
- Go `go1.26.7.windows-amd64.zip`, 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`, identity re-read from `go.dev/dl/?mode=json&include=all` before extraction.

All materialized product glue, profile and tests are `AUTHORED` and licensed by the workspace owner. The AWS SDK remains an Apache-2.0 dependency; no local file is falsely represented as verbatim Amazon code.

## Reconstruction

Pack SHA-256: `f76a1aaf69afedc3ec1f9556436be530f5b4c3c56d29223bd5cd263706aa8832`.

| Path under `aws_secure_quarantine_intake/` | Bytes | SHA-256 |
|---|---:|---|
| `go.mod` | 828 | `abf577684fa2df625266739854c1bafdd6a04f7fbada2395c18594fcf4e6c37e` |
| `go.sum` | 2,369 | `25ae7474cd904d50162f0ba82565294d83020cbe3b24496ba9137b9fa9abeced` |
| `provider-profile.template.json` | 1,076 | `c169ed6d1f1f1101c1fcfba9bcf35d7a8dfc7dae2d1ef40cf8dd187a46e75497` |
| `awsintake/intake.go` | 10,931 | `96d40a72f85121626c7f2ada2e1db1e6ea973939b209452f6b35ed121ddc4039` |
| `awsintake/intake_test.go` | 10,184 | `a887df43d9c169184f82aabedd4e9b0b8c23713f2cdafec61e7058ea5132f7d8` |
| `README.md` | 1,211 | `da69d74b660f6bc914d20c8f525231f0d00596a6e3001f4b9c1e1a7995156085` |

A fresh materialization produced six files. Every SHA-256 matched the independent staging tree. No placeholder remained.

## Executed gates

With explicit portable Go 1.26.7 and an isolated module/build cache:

```text
go mod verify                                      PASS
GOPROXY=off go test -count=1 ./...                PASS
go vet ./...                                       PASS
go build ./...                                     PASS
```

The suite executed 6 tests and 8 named subtests. It covers exact content length/type/checksum, expected owner, SSE-KMS, create-only, HTTPS PUT, TTL, subject hash, document-class allowlist, persist-before-return, provider failure, exact version/ETag/head checksum, request-binding metadata, store failure and duplicate completion/CAS.

The complete library gates then returned:

```text
VERIFY_LIBRARY_PASS
packs=53 materialized_files=499 markdown_files=355
profile=AWS_SECURE_DOCUMENT_INTAKE_PACK_PLAN.md implementation_files=35

VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=53 upstream_sources=100 ...
```

The executable Audit materialized the new six-file pack and passed its fail-closed structural contract. It explicitly printed `AWS_SECURE_QUARANTINE_INTAKE_GO_GATES_SKIPPED reason=Go_1.26.7_not_available` because the host has no global Go; that is not substituted for the isolated portable Go suite above.

## Investigated alternative

Microsoft's official `Azure-Samples/functions-quickstart-dotnet-azd-eventgrid-blob` was fixed at signed commit `0a705781a7e66b996b81661b0c141b7b78cc066a`, MIT, archive 1,974,323 bytes/SHA-256 `e8d0007dd8ebbf69af880c46a0ef9e6330431ecdbe4bf8b450ffa29a79283157`. Microsoft .NET SDK 8.0.424 security release was downloaded from official metadata; its 285,090,820-byte Windows x64 ZIP matched published SHA-512 `1787ab90635c2950672ed7c6507b000e1b212ea7d9a22fcef37061344d37c64d4c4eda12b8742601eff5b45c8736485b31c55613892f240c300190e4e88a58b0`. The sample built in a short path with 0 warnings/0 errors; graph 6 direct/94 transitive and dated NuGet query reported zero known vulnerabilities.

It was not incorporated or promoted: it has no tests/lock, copies only PDFs between containers, and does not close authentication, secure content admission, session binding or multi-channel intake. Its code does not inherit this AWS pack's admission.

## Admission and remaining blockers

`REBUILD_VERIFIED / CONDITIONED`. The reusable code is immediately materializable, but live enablement remains blocked until the project proves:

- authenticated/authorized endpoint and real transactional `SessionStore` with concurrency tests;
- account, region, IAM, temporary credentials, costs, private bucket, Block Public Access, versioning, exact KMS and expected owner;
- bucket-policy `s3:signatureAge`, exact CORS origins/exposed headers and negative tests;
- real upload/finalize, abandoned-session expiry, duplicate/out-of-order event handling and reconciliation;
- isolated download of the exact version into Magika/ClamAV/YARA-X, load/quotas/telemetry/recovery/rollback;
- no persistence to business storage until routing, extraction evaluation and review gates independently authorize it.

Email receiving, SFTP, scanner/mobile, object-event workers and ERP/marketplace intake remain gaps. They are not inferred from this evidence.
