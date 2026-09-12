# Google Document AI Terraform Evidence — 2026-08-25 V1

## Claim

This evidence proves exact acquisition and static inspection of Google's official Terraform module for creating Document AI processors. It supplies the infrastructure portion of the Google document lane. It does not prove a deployment because the audit host has no Terraform/Go/gcloud toolchain or billed GCP project.

## Exact source

```text
repository: GoogleCloudPlatform/terraform-google-document-ai
release: v0.0.1
commit: 4a07f60af61874c622d98c8d0f2193cd36c89a13
archive URL: https://github.com/GoogleCloudPlatform/terraform-google-document-ai/archive/4a07f60af61874c622d98c8d0f2193cd36c89a13.zip
archive bytes: 72254
archive SHA-256: a0b6c31df1805b8a87acabec0cfb10c1301865a9c43111818c267bc1eb8b0769
license: Apache-2.0
LICENSE SHA-256: cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30
test/integration/go.sum SHA-256: 831cc9b0d1bfa91c21db42b03d1a18288eb5bda775ed485c510bef07cdf5541e
canonical acquisition: UPSTREAM_SOURCE_ACQUIRED id=google-terraform-document-ai-0.0.1
files: 33
Terraform files: 13
test-tree files: 11
```

## Inspected contract

The module declares Terraform `>=1.3`, Google provider `>=3.53,<7`, Document AI processor name/location/type, optional KMS key, processor ID/name outputs and a simple example. Its integration test deploys the blueprint and verifies non-empty processor outputs and that `documentai.googleapis.com` is enabled.

The README requires an active billing account, permissions, `roles/documentai.admin`, enabled APIs and a GCP project. Those are project access conditions, not values that Elite can manufacture.

## Result

```text
archive/license identity: PASS
Terraform source/example inventory: PASS
Go dependency lock identity: PASS
terraform validate: NOT_RUN_MISSING_TOOLCHAIN
integration test: NOT_RUN_MISSING_GCP_ACCESS
admission: CONDITIONAL_PLATFORM
```

Before use, the project must select region/data residency, project/billing owner, service account and KMS policy; acquire this exact source; run `terraform init -lockfile=readonly`, formatting/validation/security policy checks, then the official integration test in a disposable project with destroy evidence. Processor creation is not document-field accuracy; corpus evaluation remains separately mandatory.
