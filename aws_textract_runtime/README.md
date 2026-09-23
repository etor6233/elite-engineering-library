# AWS Textract document runtime

This component calls Amazon Textract with the exact official AWS SDK for Go v2 module `service/textract v1.45.0`. AWS owns the SDK, generated clients, API models and service. The profile enforcement, security-receipt chaining, atomic evidence writer and tests in this directory are local `AUTHORED` orchestration and are not presented as AWS code.

The command does not accept an operation, feature list, query list, adapter, region or credentials chosen freely by the caller. It requires a document class declared exactly once as `REQUIRED` in an approved profile. That profile fixes the AWS region, official SDK version, operation, features, queries, optional exact adapter version, response-schema version and denies automatic storage. The distributed template enables no class.

Before any AWS request, the exact input SHA-256, byte count and MIME type must match an `ADMITTED` `elite-secure-local-file-receipt/v1` containing approval/policy evidence and complete ClamAV, Magika and YARA-X results. The wrapper preserves the complete modeled provider response, hashes profile/security/input/response, requires a proved single-page synchronous result and commits a new evidence directory atomically. It never writes predictions to a business database.

```text
go mod verify
GOPROXY=off go test ./...
go vet ./...
go build ./cmd/analyze
```

Real execution still requires an explicitly provided AWS account, IAM role, quota/cost/residency decision, exact region, approved class profile, security gate output, representative corpus and field-level evaluation. A provider response is evidence for evaluation; it is not authorization to store a business fact.
