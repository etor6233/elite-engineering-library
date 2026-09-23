# Secure email MIME quarantine core

This component begins only after a provider has retained the exact raw MIME object and produced a version-bound admission receipt. It parses attachments into a new quarantine directory, names every byte payload from its part index and SHA-256 rather than the untrusted filename, and writes a deterministic manifest with `automatic_storage_authorized=false`.

The implementation is local `AUTHORED` integration code. It is not presented as AWS source. Its contract is derived from failures observed in official AWS samples and is intended to compose with the separately materialized AWS Powertools idempotent partial-batch component, Lambda Durable Execution component, retained S3 Object Lock lane, and Magika/ClamAV/YARA security gate.

Only base64 attachments and explicitly allowlisted claimed MIME types are accepted. This is deliberately fail-closed. A later security gate must identify the real content from bytes before any extraction or business persistence.

Run `pwsh ./verify_contract.ps1`. Live AWS remains blocked until the project proves Mail Manager/SES, DNS/MX, tenant routing, KMS, retained raw MIME, SQS/DLQ/redrive, DynamoDB idempotency, Lambda packaging, security engines, costs and outage/replay reconciliation.
