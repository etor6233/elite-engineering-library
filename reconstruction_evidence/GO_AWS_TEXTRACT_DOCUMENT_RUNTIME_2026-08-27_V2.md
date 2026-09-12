# Go AWS Textract Document Runtime — Reconstruction Evidence V2

Date: 2026-08-27  
Result: `REBUILD_VERIFIED / CONDITIONED`  
Production/admission result: **not promoted**; no AWS account, credential, service call, business document, corpus accuracy or storage authority was used or proven.

## Official AWS identity

- repository: `aws/aws-sdk-go-v2`;
- module/tag: `service/textract/v1.45.0`, published 2026-08-26;
- commit: `a30468cff35d6e385287a0cbff0ac11aa7202529`;
- Go proxy origin maps the module to that exact commit;
- annotated tag and commit verification: `verified=false`, reason `unsigned`; retained as `UP-FAIL-076`;
- full source archive: 132,399,807 bytes, SHA-256 `97d5ccc9cda85d8cd1566ec5a5a0812c3f033d419d0eb6561a83125e7592fecf`;
- module ZIP: 220,255 bytes, SHA-256 `c82c76b648012d198bc80c0aa303a1f6f4581af4bcb1234a3d5746f782df1cdc`;
- module `.mod`: 502 bytes, SHA-256 `68d45c9d48429d813dfb4812af636ae0a8a2ba8ab13ac23eb6e17f4dd179a90f`;
- Apache-2.0 `LICENSE.txt`: SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`;
- `NOTICE.txt`: SHA-256 `a6e830174d62dafad3a718384772ea1f4c0e27989f050932dc442a5e3ddad080`;
- package `go.sum`: SHA-256 `511be4a0ceeb8389d0748ba017d702618054887e226cdf830e926aea7f5964d4`.

Primary references: [AWS SDK source](https://github.com/aws/aws-sdk-go-v2/tree/a30468cff35d6e385287a0cbff0ac11aa7202529/service/textract), [AnalyzeExpense](https://docs.aws.amazon.com/textract/latest/APIReference/API_AnalyzeExpense.html), [AnalyzeDocument](https://docs.aws.amazon.com/textract/latest/APIReference/API_AnalyzeDocument.html), and [Textract fixed quotas](https://docs.aws.amazon.com/textract/latest/dg/limits-document.html).

## Official-source reconstruction

The 147 files published in the Go module were compared byte-for-byte against the `service/textract` subpath of the exact commit archive: 147 checked, zero missing, zero mismatches. The source package contains nine official `_test.go` files. With the official `go1.26.7 windows/amd64` archive (74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`):

- `go mod verify`: PASS;
- complete `go test -count=1 ./...`: PASS for `service/textract` and `internal/endpoints`;
- `go vet ./...`: PASS.

The rebuilt acquisition runner independently downloaded the same full archive and emitted `UPSTREAM_SOURCE_ACQUIRED id=aws-sdk-go-v2-textract-1.45.0` after checking archive, license, NOTICE, package manifest and package constraints. Lock tests: six negatives PASS. Source profiles: nine valid, five negatives and one positive PASS. The lock now contains 72 exact sources; AWS document profile selects 10 and document intelligence selects 33.

## Local wrapper gap and correction

The prior V0.1 wrapper accepted operation, region, features, Queries and adapter directly from CLI, treated PDF-like extensions as sufficient and did not require the secure local-file receipt. This allowed profile bypass and bytes not proven by Magika/ClamAV/YARA-X. It is retained as `LIB-FAIL-243` and was corrected rather than relabeled as AWS code.

`GO-AWS-TEXTRACT-DOCUMENT-RUNTIME 0.2.0` remains entirely explicit about provenance: AWS owns the SDK/client/types/service; the profile validator, security chaining, atomic writer and tests are local `AUTHORED` integration. The command now accepts only document class, approved profile, admitted security receipt, input/output and a non-expanding byte cap. It derives region, operation, features, Queries, adapter and response-schema version from a class that must appear exactly once as `REQUIRED`. The distributed V2 template enables zero classes.

Before any provider call, the runtime requires an `elite-secure-local-file-receipt/v1` with:

- `ADMITTED / ALL_SELECTED_GATES_PASSED`;
- approval ID and policy SHA-256;
- input SHA-256, byte count and MIME matching the exact file;
- complete ClamAV, Magika and YARA-X binary/version/results;
- zero YARA matches;
- explicit `business_storage_authorized=false`.

The V2 analysis receipt hash-links document profile, security receipt, input and complete modeled provider response. It requires a one-page synchronous result and still sets `automatic_storage_authorized=false`.

## Target graph and tests

Resolved graph: 15 versioned Go modules including `config v1.32.38`, AWS core `v1.44.0`, Textract `v1.45.0` and Smithy `v1.28.1`. `go.sum` is materialized. OSV query at `2026-08-27T15:54:05Z`: 15 modules, zero findings; this is dated evidence, not a future guarantee.

Clean runtime result:

- seven materialized files, all canonical SHA-256 verified;
- `go mod tidy`, `go mod verify`: PASS;
- seven tests plus three security subtests: PASS;
- `go vet ./...`: PASS;
- `go build ./cmd/analyze`: PASS;
- repeated `GOPROXY=off go test -count=1 ./...`: PASS;
- unsafe request fields `Mode`, `Features`, `Queries`, `AdapterID` and `AdapterVersion`: absent by reflection;
- blocked/duplicate class, region mismatch, rejected/wrong-MIME/incomplete security receipt, provider failure and unproved page: fail closed without committed output.

`VERIFY_LIBRARY.ps1` passed 41 packs/366 files and composes the AWS lane as 4 packs/36 files. `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit` passed with 72 sources and always checks AWS materialization, exact pin and zero-enabled template. Foundation mode now also executes module download only when network is authorized, then verify/tests/vet/build/version contract.

## Failures converted into regression

- `LIB-FAIL-239`: invalid PowerShell root selector; fixed with exact-one-root validation;
- `LIB-FAIL-240`: array lost through `pwsh -File`; fixed by in-process typed-array invocation and 7/7 reconstruction;
- `LIB-FAIL-241`: stale `$LASTEXITCODE` after a passing PowerShell test; fixed by exact PASS markers and exception propagation;
- `LIB-FAIL-242`: historical ID placed between backticks and read as duplicate; fixed by reserving that syntax for row declarations.
- `LIB-FAIL-244`: first temp cleanup encountered read-only source files; fixed by clearing only the `ReadOnly` bit beneath the revalidated exact audit root before deletion.
- `LIB-FAIL-245`: Go module cache also marked directories read-only; fixed by clearing that bit on descendant directories, deepest first, within the same exact root.

All six are `REGRESSION_PROVEN` in the governing ledger.

## Remaining project gates

This lane is executable immediately only after the owner supplies and approves AWS account, exact region, IAM role, quota/cost/residency, secure gate assets/policy, class-specific profile, representative corpus and ground truth. Queries/adapters require exact versions and their own train/test evidence. A live sandbox contract probe, concurrency/throttle/load, privacy, review/correction, idempotent persistence, recovery and rollback remain unproven. Therefore this V2 is not `REUSABLE_PACK` and cannot authorize automatic storage.
