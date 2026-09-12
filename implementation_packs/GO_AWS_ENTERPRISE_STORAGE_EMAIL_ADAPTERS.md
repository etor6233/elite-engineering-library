# Go AWS Enterprise Storage and Email Adapters

## 1. Metadata

```yaml
pack_id: "GO-AWS-ENTERPRISE-STORAGE-EMAIL-ADAPTERS"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa adapters Go sobre SDKs oficiales Amazon S3 1.107.3 y SES v2 1.67.0 para upload con checksum/encryption, Object Lock verificado por versión y email transaccional con receipt sin PII."
stacks: ["Go 1.26.7", "AWS SDK for Go v2", "Amazon S3", "Amazon SES v2"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "GO-RELIABLE-ASYNC-WORKERS 0.2.x"]
incompatible_with: ["credencial en código/CLI", "overwrite silencioso", "bucket sin política", "email directo sin outbox", "auto-storage documental"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/aws/aws-sdk-go-v2/tree/284a4846e7bb941926e2d15144b17c01ffc98c1d/service/s3", "https://github.com/aws/aws-sdk-go-v2/tree/284a4846e7bb941926e2d15144b17c01ffc98c1d/service/sesv2", "https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/s3-checksums.html", "https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock-configure.html", "https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock-managing.html"]
verified_at: "2026-08-27"
```

## 2. Applicability

Use si el blueprint selecciona AWS para object storage y/o email, y existen cuenta, región, IAM y costos aprobados. S3 cubre un single `PutObject` de bytes controlados con checksum SHA-256, cifrado, owner esperado y retención Object Lock explícita/verificada sobre el VersionId creado. SES cubre email simple transaccional tras outbox/deduplicación propia. Rechazar como sustituto de multipart, malware scanning, habilitación/configuración de bucket, lifecycle, legal hold, delivery reconciliation, marketing consent o proveedor no seleccionado.

Los adapters son `AUTHORED` glue sobre SDKs oficiales AWS. No copian el SDK ni atribuyen a Amazon las reglas del negocio. El profile inicia bloqueado y exige las decisiones operativas reales.

## 3. Architecture contract

S3 recibe bucket/key relativo/content type/bytes, owner AWS esperado, modo `AES256` o `aws:kms`, `GOVERNANCE|COMPLIANCE` y fecha futura. Antes de escribir consulta `GetObjectLockConfiguration`; usa `IfNoneMatch: *`, `ChecksumSHA256`, metadata hash y exige ETag/VersionId. Después consulta `GetObjectRetention` sobre esa versión y sólo marca `retention_verified=true` si modo/fecha coinciden. Ante verificación post-write fallida devuelve error más receipt de la versión creada para impedir retry ciego. Nunca habilita Object Lock ni legal hold. SES recibe sender/recipients/contenido, requiere `MessageKey` de un outbox externo, usa `sesv2.SendEmail`, etiqueta el mensaje y retorna provider ID + hash de contenido sin guardar addresses/body en el receipt.

Input, provider o response inválidos fallan cerrados; un efecto remoto parcialmente verificado nunca se oculta. Las credenciales viven en la default credential chain. Retry/backoff del SDK no reemplaza idempotencia/reconciliation. Performance, cost, quota, privacy, malware, bucket controls/lifecycle, bounce/complaint, recovery y delivery se prueban en el target.

## 4. Exact file manifest

```text
CREATE aws_enterprise_adapters/go.mod
CREATE aws_enterprise_adapters/go.sum
CREATE aws_enterprise_adapters/provider-profile.template.json
CREATE aws_enterprise_adapters/awsenterprise/clients.go
CREATE aws_enterprise_adapters/awsenterprise/storage.go
CREATE aws_enterprise_adapters/awsenterprise/storage_test.go
CREATE aws_enterprise_adapters/awsenterprise/email.go
CREATE aws_enterprise_adapters/awsenterprise/email_test.go
CREATE aws_enterprise_adapters/README.md
```

## 5. Materialization blocks

### FILE: `aws_enterprise_adapters/go.mod`
```yaml
block_id: "GO-AWS-ENTERPRISE-ADAPTERS:gomod:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official AWS Go SDK module versions"
license: "LicenseRef-Workspace-Owner"
sha256: "e4f8102a121d560ed88b561a1b3cfe5ec7f5b147fdb98e0c90fe0987bb2fc567"
variables: []
secrets_allowed: false
```
````text
module example.com/elite/aws-enterprise-adapters

go 1.26.0

require (
	github.com/aws/aws-sdk-go-v2 v1.43.8
	github.com/aws/aws-sdk-go-v2/config v1.32.38
	github.com/aws/aws-sdk-go-v2/service/s3 v1.107.3
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.67.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.18 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.38 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.39 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.39 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.39 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.40 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.39 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.39 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.5.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.33.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.45.8 // indirect
	github.com/aws/smithy-go v1.27.10 // indirect
)
````

### FILE: `aws_enterprise_adapters/go.sum`
```yaml
block_id: "GO-AWS-ENTERPRISE-ADAPTERS:gosum:v1"
operation: CREATE
provenance: AUTHORED
source: "Go checksum database entries from exact direct pins"
license: "LicenseRef-Workspace-Owner"
sha256: "c040cd923e3c5b2c3503f7064c3bb9bd9a6a6c07448f13b60dd723268f123916"
variables: []
secrets_allowed: false
```
````text
github.com/aws/aws-sdk-go-v2 v1.43.8 h1:fpnrxwuwsoGIgjvgLeDU3y9w7YaHBxyF6AF3vQL8duw=
github.com/aws/aws-sdk-go-v2 v1.43.8/go.mod h1:j7gYSq8dL95QejkFXxvQNESH4I9WGHFI6iO+vhqEi5Q=
github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.18 h1:LAfOuhAH331fmOjTQpAaOlH+Ftn7RzSDJ2VFwjdMMy4=
github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.18/go.mod h1:4e5xhuXHx1e4U9EthvbPP1r/DIMp5c2823OL8karzcM=
github.com/aws/aws-sdk-go-v2/config v1.32.38 h1:n4yPHBjtQ3BrIIUyk0/LAqf/BL2iv0Tw6XZcMRzM0ps=
github.com/aws/aws-sdk-go-v2/config v1.32.38/go.mod h1:dencYsOS1R7rBy8zehCvwBYzdxxL4Q/nRK7In03wjN8=
github.com/aws/aws-sdk-go-v2/credentials v1.19.38 h1:Xf8j1+vzwPRCta9pFXjj0677BzXrRO2JbpAVNcdXnnI=
github.com/aws/aws-sdk-go-v2/credentials v1.19.38/go.mod h1:PGYzFTznwRAJ2q0m+oX+P8SlfZQKpBAKQCokNuMl3Sg=
github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.39 h1:9GLrXl8PKQ3+bMniXFg3vliMWJ+204bFcIvBCwJFglc=
github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.39/go.mod h1:MmlE5TLgq7+QbXKKUSzqUz4h0Uu5kz2SEe6iPX+ZFHI=
github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.39 h1:YrEI22hVQcqMpq934ZoPQyJjGNzX4CGdrSDCjBD59sI=
github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.39/go.mod h1:N8qOX83LkaCeizvrfiNjwkBOXkxHt6a74CiZn8qz9F8=
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.39 h1:Vo7UZzBjB6zS6feEOuBlpEgaj8iBTdiNlye+7w9ooGo=
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.39/go.mod h1:JgxtAO/77e95Rs9WMWUzz99hT182gqdAh7/DHuEMA/k=
github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.40 h1:oofDq8Y5M82fmDrxb8gsbP0LS73MqZ388qKVgs5ETYI=
github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.40/go.mod h1:LSfLmbvx50+T+/DoUZRqB1qS38v7lvNUebqIpidAWYM=
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.18 h1:+fiwOxNdE8bOK3SoVTln8hwP+OCyArbi2/InIr/A9AU=
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.18/go.mod h1:aua4m7EZSvQra/96b8zJxWHwtHxuXQ8bx4DiM92V044=
github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.31 h1:uZOinZb+h7lZw8IYzP1z1IuEnueB76/EFkcf/fEW4Ag=
github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.31/go.mod h1:NRtwAM/p5VRt03TlEUs0pH3TeWamWdf4YyJpSrzPYLc=
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.39 h1:inoUrqz4Lfpw1XwpUvQnBiAJ2tUzn3opZ0gduNLxo+8=
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.39/go.mod h1:Yx+RrmAF+XGZTccwhQ3o4K5V8qkZBsTAcq148Y8g57k=
github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.39 h1:HLPAVrlLDaN2boN0xJx7MgaQDNEO3Q+c9L6kl/8m47Q=
github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.39/go.mod h1:Pg/dVfsNkm1hsIDK/gMvCKtmyNfNTV12mrgHqVE/6Oo=
github.com/aws/aws-sdk-go-v2/service/s3 v1.107.3 h1:IKoCZqfWfZzSBi16QFQ+QcbQ3LRQ7QgB1S5tDAyPBQQ=
github.com/aws/aws-sdk-go-v2/service/s3 v1.107.3/go.mod h1:RBpRcXiM4s2pOInVs32GsBonnje+fiAj4mcrStRmlCA=
github.com/aws/aws-sdk-go-v2/service/sesv2 v1.67.0 h1:cvmzhKyIYHkR+ULgWBYK672NzybWJiANO31uOsv0Imo=
github.com/aws/aws-sdk-go-v2/service/sesv2 v1.67.0/go.mod h1:zHA87gWVfSnNYawE3e4ghWT3nxJOGSd1Ml0r+Epx/8o=
github.com/aws/aws-sdk-go-v2/service/signin v1.5.8 h1:bghrxelVQpGurGI1X94BT68h6p+hWQnlsu8nSmiSll4=
github.com/aws/aws-sdk-go-v2/service/signin v1.5.8/go.mod h1:gkwdIl9w+6LFKlGRLz3+Dw+cudc9dD1ViMDhHGmzOgk=
github.com/aws/aws-sdk-go-v2/service/sso v1.33.8 h1:/DbiPZ8maO03uFnXa6yEhFdWOTA5xObmGNfaEzt9Cac=
github.com/aws/aws-sdk-go-v2/service/sso v1.33.8/go.mod h1:mUywXl2WlN+gZD0vNeg1Hn0EMOifDQ79StJcdqXHkXo=
github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.8 h1:wv4pCyq/LkBYc5R4m/g5S+uGqF/DbL+bp9VXiQEnec4=
github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.8/go.mod h1:9AKVT0vADSCPXRuoZjziHwsbdLDFMGRExwWBQourCa8=
github.com/aws/aws-sdk-go-v2/service/sts v1.45.8 h1:oQrmuqpBAExYPEPJp8dkj9KLmc0y42iwvAV28OwlzF0=
github.com/aws/aws-sdk-go-v2/service/sts v1.45.8/go.mod h1:qNTXKrmzx2cC6VmM7PxHNasBMWKx3mfxgzcbVjcWVAU=
github.com/aws/smithy-go v1.27.10 h1:bw56MIx8bhTQZSdzucEJSKWLpwX0ju7hU8cVoa75dg8=
github.com/aws/smithy-go v1.27.10/go.mod h1:YE2RhdIuDbA5E5bTdciG9KrW3+TiEONeUWCqxX9i1Fc=
````

### FILE: `aws_enterprise_adapters/provider-profile.template.json`
```yaml
block_id: "GO-AWS-ENTERPRISE-ADAPTERS:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed AWS access/policy profile"
license: "LicenseRef-Workspace-Owner"
sha256: "fbcc5e7166e67c785b09592161bbdc12777050329924f0e23c2d4a78a375e79f"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-aws-enterprise-adapters-profile/v1",
  "provider": "Amazon Web Services",
  "region": "",
  "credentials": {"mechanism": "AWS SDK default credential chain", "secret_reference_required": true},
  "object_storage": {
    "decision": "BLOCKED_ACCESS_AND_POLICY_REQUIRED",
    "bucket": "",
    "key_prefix": "",
    "encryption": "aws:kms",
    "kms_key_id": "",
    "expected_bucket_owner": "",
    "versioning_required": true,
    "object_lock": {
      "decision": "BLOCKED_OWNER_MODE_RETENTION_REQUIRED",
      "mode": "",
      "retain_until_source": "",
      "legal_hold": "SEPARATE_EXPLICIT_APPROVAL_REQUIRED"
    },
    "lifecycle_decision_required": true,
    "cost_approved": false
  },
  "transactional_email": {
    "decision": "BLOCKED_ACCESS_AND_POLICY_REQUIRED",
    "verified_from_identity": "",
    "configuration_set": "",
    "sandbox_or_production_access": "",
    "suppression_and_bounce_route_required": true,
    "outbox_idempotency_required": true,
    "cost_approved": false
  }
}
````

### FILE: `aws_enterprise_adapters/awsenterprise/clients.go`
```yaml
block_id: "GO-AWS-ENTERPRISE-ADAPTERS:clients:v1"
operation: CREATE
provenance: AUTHORED
source: "local factory invoking AWS config and service client constructors"
license: "LicenseRef-Workspace-Owner"
sha256: "b93d2ad5a0fb038e925e8a26a67f0ab8a8588c0721fce434497c376a632c4e6c"
variables: []
secrets_allowed: false
```
````go
package awsenterprise

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

type Clients struct {
	S3    *s3.Client
	SESV2 *sesv2.Client
}

func NewClients(ctx context.Context, region string) (Clients, error) {
	if strings.TrimSpace(region) == "" || strings.ContainsAny(region, " /\\") {
		return Clients{}, errors.New("valid AWS region is required")
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return Clients{}, err
	}
	return Clients{S3: s3.NewFromConfig(cfg), SESV2: sesv2.NewFromConfig(cfg)}, nil
}
````

### FILE: `aws_enterprise_adapters/awsenterprise/storage.go`
```yaml
block_id: "GO-AWS-ENTERPRISE-ADAPTERS:storage:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter invoking AWS S3 SDK public API"
license: "LicenseRef-Workspace-Owner"
sha256: "757365fc7479475a103145deb89b019801c857401e3d050d49b8dd6b41a01277"
variables: []
secrets_allowed: false
```
````go
package awsenterprise

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const MaxSinglePutBytes = 64 * 1024 * 1024

type S3Client interface {
	GetObjectLockConfiguration(context.Context, *s3.GetObjectLockConfigurationInput, ...func(*s3.Options)) (*s3.GetObjectLockConfigurationOutput, error)
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObjectRetention(context.Context, *s3.GetObjectRetentionInput, ...func(*s3.Options)) (*s3.GetObjectRetentionOutput, error)
}

type PutObjectRequest struct {
	Bucket              string
	Key                 string
	ContentType         string
	Payload             []byte
	Encryption          string
	KMSKeyID            string
	ExpectedBucketOwner string
	RetentionMode       string
	RetainUntil         time.Time
}

type PutObjectReceipt struct {
	Provider            string `json:"provider"`
	SDKModule           string `json:"sdk_module"`
	SDKVersion          string `json:"sdk_version"`
	Bucket              string `json:"bucket"`
	Key                 string `json:"key"`
	Bytes               int    `json:"bytes"`
	SHA256              string `json:"sha256"`
	ChecksumSHA256      string `json:"checksum_sha256_base64"`
	ETag                string `json:"etag"`
	VersionID           string `json:"version_id,omitempty"`
	Encryption          string `json:"encryption"`
	ExpectedBucketOwner string `json:"expected_bucket_owner"`
	RetentionMode       string `json:"retention_mode"`
	RetainUntil         string `json:"retain_until"`
	RetentionVerified   bool   `json:"retention_verified"`
}

func validAWSAccountID(value string) bool {
	if len(value) != 12 {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

func PutImmutableObject(ctx context.Context, client S3Client, request PutObjectRequest) (PutObjectReceipt, error) {
	var receipt PutObjectReceipt
	if client == nil {
		return receipt, errors.New("S3 client is required")
	}
	if strings.TrimSpace(request.Bucket) == "" || strings.TrimSpace(request.Key) == "" || strings.HasPrefix(request.Key, "/") || strings.Contains(request.Key, "\\") {
		return receipt, errors.New("bucket and canonical relative key are required")
	}
	if len(request.Payload) == 0 || len(request.Payload) > MaxSinglePutBytes {
		return receipt, errors.New("payload is outside the single-put policy limit")
	}
	if strings.TrimSpace(request.ContentType) == "" {
		return receipt, errors.New("content type is required")
	}
	if !validAWSAccountID(request.ExpectedBucketOwner) {
		return receipt, errors.New("expected bucket owner must be a 12-digit AWS account ID")
	}
	retainUntil := request.RetainUntil.UTC().Truncate(time.Second)
	if retainUntil.Before(time.Now().UTC().Add(5 * time.Minute)) {
		return receipt, errors.New("retain-until must be at least five minutes in the future")
	}
	var lockMode types.ObjectLockMode
	var retentionMode types.ObjectLockRetentionMode
	switch request.RetentionMode {
	case "GOVERNANCE":
		lockMode = types.ObjectLockModeGovernance
		retentionMode = types.ObjectLockRetentionModeGovernance
	case "COMPLIANCE":
		lockMode = types.ObjectLockModeCompliance
		retentionMode = types.ObjectLockRetentionModeCompliance
	default:
		return receipt, errors.New("retention mode must be GOVERNANCE or COMPLIANCE")
	}
	lockConfiguration, err := client.GetObjectLockConfiguration(ctx, &s3.GetObjectLockConfigurationInput{Bucket: aws.String(request.Bucket), ExpectedBucketOwner: aws.String(request.ExpectedBucketOwner)})
	if err != nil {
		return receipt, fmt.Errorf("verify bucket object lock: %w", err)
	}
	if lockConfiguration == nil || lockConfiguration.ObjectLockConfiguration == nil || lockConfiguration.ObjectLockConfiguration.ObjectLockEnabled != types.ObjectLockEnabledEnabled {
		return receipt, errors.New("bucket Object Lock is not proven enabled")
	}
	sum := sha256.Sum256(request.Payload)
	checksum := base64.StdEncoding.EncodeToString(sum[:])
	input := &s3.PutObjectInput{
		Bucket: aws.String(request.Bucket), Key: aws.String(request.Key), Body: bytes.NewReader(request.Payload), ContentLength: aws.Int64(int64(len(request.Payload))), ContentType: aws.String(request.ContentType),
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256, ChecksumSHA256: aws.String(checksum), IfNoneMatch: aws.String("*"), Metadata: map[string]string{"content-sha256": hex.EncodeToString(sum[:])},
		ExpectedBucketOwner: aws.String(request.ExpectedBucketOwner), ObjectLockMode: lockMode, ObjectLockRetainUntilDate: aws.Time(retainUntil),
	}
	switch request.Encryption {
	case "AES256":
		if request.KMSKeyID != "" {
			return receipt, errors.New("KMS key is invalid with AES256")
		}
		input.ServerSideEncryption = types.ServerSideEncryptionAes256
	case "aws:kms":
		if strings.TrimSpace(request.KMSKeyID) == "" {
			return receipt, errors.New("KMS key is required for aws:kms")
		}
		input.ServerSideEncryption = types.ServerSideEncryptionAwsKms
		input.SSEKMSKeyId = aws.String(request.KMSKeyID)
	default:
		return receipt, errors.New("encryption must be AES256 or aws:kms")
	}
	output, err := client.PutObject(ctx, input)
	if err != nil {
		return receipt, err
	}
	if output == nil || aws.ToString(output.ETag) == "" || aws.ToString(output.VersionId) == "" {
		return receipt, errors.New("S3 Object Lock response requires ETag and VersionId")
	}
	if observed := aws.ToString(output.ChecksumSHA256); observed != "" && observed != checksum {
		return receipt, fmt.Errorf("S3 checksum mismatch: expected %s, got %s", checksum, observed)
	}
	receipt = PutObjectReceipt{Provider: "Amazon S3", SDKModule: "github.com/aws/aws-sdk-go-v2/service/s3", SDKVersion: "v1.107.3", Bucket: request.Bucket, Key: request.Key, Bytes: len(request.Payload), SHA256: hex.EncodeToString(sum[:]), ChecksumSHA256: checksum, ETag: aws.ToString(output.ETag), VersionID: aws.ToString(output.VersionId), Encryption: request.Encryption, ExpectedBucketOwner: request.ExpectedBucketOwner, RetentionMode: request.RetentionMode, RetainUntil: retainUntil.Format(time.RFC3339)}
	observedRetention, err := client.GetObjectRetention(ctx, &s3.GetObjectRetentionInput{Bucket: aws.String(request.Bucket), Key: aws.String(request.Key), VersionId: output.VersionId, ExpectedBucketOwner: aws.String(request.ExpectedBucketOwner)})
	if err != nil {
		return receipt, fmt.Errorf("verify object retention for created version %s: %w", receipt.VersionID, err)
	}
	if observedRetention == nil || observedRetention.Retention == nil || observedRetention.Retention.RetainUntilDate == nil || observedRetention.Retention.Mode != retentionMode || !observedRetention.Retention.RetainUntilDate.UTC().Equal(retainUntil) {
		return receipt, fmt.Errorf("object retention verification mismatch for created version %s", receipt.VersionID)
	}
	receipt.RetentionVerified = true
	return receipt, nil
}
````

### FILE: `aws_enterprise_adapters/awsenterprise/storage_test.go`
```yaml
block_id: "GO-AWS-ENTERPRISE-ADAPTERS:storage-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract/negative tests against official S3 types"
license: "LicenseRef-Workspace-Owner"
sha256: "ef872d300edd1389dfabfd9f718449e642c87b80a7650086d3eb86265844a7da"
variables: []
secrets_allowed: false
```
````go
package awsenterprise

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type fakeS3 struct {
	input           *s3.PutObjectInput
	output          *s3.PutObjectOutput
	err             error
	lockOutput      *s3.GetObjectLockConfigurationOutput
	lockErr         error
	retentionOutput *s3.GetObjectRetentionOutput
	retentionErr    error
	lockInput       *s3.GetObjectLockConfigurationInput
	retentionInput  *s3.GetObjectRetentionInput
}

func (f *fakeS3) GetObjectLockConfiguration(_ context.Context, input *s3.GetObjectLockConfigurationInput, _ ...func(*s3.Options)) (*s3.GetObjectLockConfigurationOutput, error) {
	f.lockInput = input
	return f.lockOutput, f.lockErr
}

func (f *fakeS3) PutObject(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	f.input = input
	return f.output, f.err
}

func (f *fakeS3) GetObjectRetention(_ context.Context, input *s3.GetObjectRetentionInput, _ ...func(*s3.Options)) (*s3.GetObjectRetentionOutput, error) {
	f.retentionInput = input
	return f.retentionOutput, f.retentionErr
}

func admittedS3(retainUntil time.Time) *fakeS3 {
	return &fakeS3{
		lockOutput:      &s3.GetObjectLockConfigurationOutput{ObjectLockConfiguration: &types.ObjectLockConfiguration{ObjectLockEnabled: types.ObjectLockEnabledEnabled}},
		output:          &s3.PutObjectOutput{ETag: aws.String("etag-1"), VersionId: aws.String("version-1")},
		retentionOutput: &s3.GetObjectRetentionOutput{Retention: &types.ObjectLockRetention{Mode: types.ObjectLockRetentionModeGovernance, RetainUntilDate: aws.Time(retainUntil)}},
	}
}

func TestPutImmutableObjectUsesChecksumEncryptionAndNoOverwrite(t *testing.T) {
	retainUntil := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	client := admittedS3(retainUntil)
	receipt, err := PutImmutableObject(context.Background(), client, PutObjectRequest{Bucket: "evidence", Key: "org/document.pdf", ContentType: "application/pdf", Payload: []byte("content"), Encryption: "aws:kms", KMSKeyID: "key-1", ExpectedBucketOwner: "123456789012", RetentionMode: "GOVERNANCE", RetainUntil: retainUntil})
	if err != nil {
		t.Fatal(err)
	}
	if client.input.ChecksumAlgorithm != types.ChecksumAlgorithmSha256 || aws.ToString(client.input.IfNoneMatch) != "*" || client.input.ServerSideEncryption != types.ServerSideEncryptionAwsKms || client.input.ObjectLockMode != types.ObjectLockModeGovernance || !aws.ToTime(client.input.ObjectLockRetainUntilDate).Equal(retainUntil) {
		t.Fatalf("unsafe input: %+v", client.input)
	}
	if _, err := base64.StdEncoding.DecodeString(receipt.ChecksumSHA256); err != nil {
		t.Fatal(err)
	}
	if receipt.VersionID != "version-1" || receipt.SHA256 == "" || !receipt.RetentionVerified || receipt.RetentionMode != "GOVERNANCE" || receipt.RetainUntil != retainUntil.Format(time.RFC3339) {
		t.Fatalf("unexpected receipt: %+v", receipt)
	}
	if aws.ToString(client.lockInput.ExpectedBucketOwner) != "123456789012" || aws.ToString(client.retentionInput.VersionId) != "version-1" {
		t.Fatalf("retention checks did not bind expected owner/version: lock=%+v retention=%+v", client.lockInput, client.retentionInput)
	}
}

func TestPutImmutableObjectFailsClosed(t *testing.T) {
	cases := []PutObjectRequest{
		{},
		{Bucket: "b", Key: "/absolute", ContentType: "x", Payload: []byte("x"), Encryption: "AES256"},
		{Bucket: "b", Key: "k", ContentType: "x", Payload: []byte("x"), Encryption: "aws:kms"},
		{Bucket: "b", Key: "k", ContentType: "x", Payload: []byte("x"), Encryption: "AES256", KMSKeyID: "wrong"},
	}
	for _, request := range cases {
		if _, err := PutImmutableObject(context.Background(), &fakeS3{}, request); err == nil {
			t.Fatalf("invalid request accepted: %+v", request)
		}
	}
	retainUntil := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	client := admittedS3(retainUntil)
	client.output.ChecksumSHA256 = aws.String("wrong")
	if _, err := PutImmutableObject(context.Background(), client, PutObjectRequest{Bucket: "b", Key: "k", ContentType: "x", Payload: []byte("x"), Encryption: "AES256", ExpectedBucketOwner: "123456789012", RetentionMode: "GOVERNANCE", RetainUntil: retainUntil}); err == nil {
		t.Fatal("checksum mismatch accepted")
	}
}

func TestPutImmutableObjectRejectsUnprovenObjectLock(t *testing.T) {
	retainUntil := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	request := PutObjectRequest{Bucket: "b", Key: "k", ContentType: "x", Payload: []byte("x"), Encryption: "AES256", ExpectedBucketOwner: "123456789012", RetentionMode: "GOVERNANCE", RetainUntil: retainUntil}
	client := admittedS3(retainUntil)
	client.lockOutput = &s3.GetObjectLockConfigurationOutput{}
	if _, err := PutImmutableObject(context.Background(), client, request); err == nil || client.input != nil {
		t.Fatal("upload proceeded without proven bucket Object Lock")
	}
	client = admittedS3(retainUntil)
	client.lockErr = errors.New("access denied")
	if _, err := PutImmutableObject(context.Background(), client, request); err == nil || client.input != nil {
		t.Fatal("upload proceeded after Object Lock probe error")
	}
}

func TestPutImmutableObjectReturnsCreatedVersionWhenRetentionVerificationFails(t *testing.T) {
	retainUntil := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	client := admittedS3(retainUntil)
	client.retentionOutput.Retention.Mode = types.ObjectLockRetentionModeCompliance
	receipt, err := PutImmutableObject(context.Background(), client, PutObjectRequest{Bucket: "b", Key: "k", ContentType: "x", Payload: []byte("x"), Encryption: "AES256", ExpectedBucketOwner: "123456789012", RetentionMode: "GOVERNANCE", RetainUntil: retainUntil})
	if err == nil || receipt.VersionID != "version-1" || receipt.RetentionVerified {
		t.Fatalf("partial provider effect was not exposed safely: receipt=%+v err=%v", receipt, err)
	}
}
````

### FILE: `aws_enterprise_adapters/awsenterprise/email.go`
```yaml
block_id: "GO-AWS-ENTERPRISE-ADAPTERS:email:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter invoking AWS SES v2 SDK public API"
license: "LicenseRef-Workspace-Owner"
sha256: "76768161e340e3cd2601d2d28ef6250b0517dedc2b7211a0bb2d300dade11054"
variables: []
secrets_allowed: false
```
````go
package awsenterprise

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/mail"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SESClient interface {
	SendEmail(context.Context, *sesv2.SendEmailInput, ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

type SendEmailRequest struct {
	From             string
	To               []string
	Subject          string
	TextBody         string
	HTMLBody         string
	ConfigurationSet string
	MessageKey       string
}

type SendEmailReceipt struct {
	Provider       string `json:"provider"`
	SDKModule      string `json:"sdk_module"`
	SDKVersion     string `json:"sdk_version"`
	ProviderID     string `json:"provider_id"`
	MessageKey     string `json:"message_key"`
	RecipientCount int    `json:"recipient_count"`
	ContentSHA256  string `json:"content_sha256"`
}

func SendTransactionalEmail(ctx context.Context, client SESClient, request SendEmailRequest) (SendEmailReceipt, error) {
	var receipt SendEmailReceipt
	if client == nil {
		return receipt, errors.New("SES client is required")
	}
	if strings.TrimSpace(request.MessageKey) == "" || len(request.MessageKey) > 256 {
		return receipt, errors.New("external outbox message key is required")
	}
	if !validAddress(request.From) || len(request.To) == 0 || len(request.To) > 50 {
		return receipt, errors.New("valid sender and 1..50 recipients are required")
	}
	seen := map[string]struct{}{}
	recipients := make([]string, 0, len(request.To))
	for _, raw := range request.To {
		address := strings.ToLower(strings.TrimSpace(raw))
		if !validAddress(address) {
			return receipt, errors.New("invalid recipient")
		}
		if _, duplicate := seen[address]; duplicate {
			return receipt, errors.New("duplicate recipient")
		}
		seen[address] = struct{}{}
		recipients = append(recipients, address)
	}
	if strings.TrimSpace(request.Subject) == "" || len(request.Subject) > 998 || (request.TextBody == "" && request.HTMLBody == "") {
		return receipt, errors.New("subject and at least one body are required")
	}
	if len(request.TextBody)+len(request.HTMLBody)+len(request.Subject) > 9*1024*1024 {
		return receipt, errors.New("message content exceeds conservative SES policy limit")
	}
	body := &types.Body{}
	if request.TextBody != "" {
		body.Text = &types.Content{Data: aws.String(request.TextBody), Charset: aws.String("UTF-8")}
	}
	if request.HTMLBody != "" {
		body.Html = &types.Content{Data: aws.String(request.HTMLBody), Charset: aws.String("UTF-8")}
	}
	input := &sesv2.SendEmailInput{FromEmailAddress: aws.String(request.From), Destination: &types.Destination{ToAddresses: recipients}, Content: &types.EmailContent{Simple: &types.Message{Subject: &types.Content{Data: aws.String(request.Subject), Charset: aws.String("UTF-8")}, Body: body}}, EmailTags: []types.MessageTag{{Name: aws.String("elite-message-key"), Value: aws.String(request.MessageKey)}}}
	if request.ConfigurationSet != "" {
		input.ConfigurationSetName = aws.String(request.ConfigurationSet)
	}
	output, err := client.SendEmail(ctx, input)
	if err != nil {
		return receipt, err
	}
	if output == nil || aws.ToString(output.MessageId) == "" {
		return receipt, errors.New("SES response contains no message ID")
	}
	hashInput := strings.Join([]string{request.Subject, request.TextBody, request.HTMLBody}, "\x00")
	sum := sha256.Sum256([]byte(hashInput))
	return SendEmailReceipt{Provider: "Amazon SES v2", SDKModule: "github.com/aws/aws-sdk-go-v2/service/sesv2", SDKVersion: "v1.67.0", ProviderID: aws.ToString(output.MessageId), MessageKey: request.MessageKey, RecipientCount: len(recipients), ContentSHA256: hex.EncodeToString(sum[:])}, nil
}

func validAddress(value string) bool {
	parsed, err := mail.ParseAddress(strings.TrimSpace(value))
	return err == nil && parsed.Address == strings.TrimSpace(value) && !slices.Contains([]string{"", "<>"}, parsed.Address)
}
````

### FILE: `aws_enterprise_adapters/awsenterprise/email_test.go`
```yaml
block_id: "GO-AWS-ENTERPRISE-ADAPTERS:email-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract/negative tests against official SES v2 types"
license: "LicenseRef-Workspace-Owner"
sha256: "7e4ba9584ac10d8ce13c1b0d2fac6a498bc23a446de40723a66e4faff04b6146"
variables: []
secrets_allowed: false
```
````go
package awsenterprise

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

type fakeSES struct {
	input  *sesv2.SendEmailInput
	output *sesv2.SendEmailOutput
	err    error
}

func (f *fakeSES) SendEmail(_ context.Context, input *sesv2.SendEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	f.input = input
	return f.output, f.err
}

func TestSendTransactionalEmailUsesOfficialSDKAndSafeReceipt(t *testing.T) {
	client := &fakeSES{output: &sesv2.SendEmailOutput{MessageId: aws.String("provider-1")}}
	receipt, err := SendTransactionalEmail(context.Background(), client, SendEmailRequest{From: "verified@example.com", To: []string{"customer@example.com"}, Subject: "Order update", TextBody: "Ready", HTMLBody: "<p>Ready</p>", ConfigurationSet: "events", MessageKey: "outbox-1"})
	if err != nil {
		t.Fatal(err)
	}
	if aws.ToString(client.input.FromEmailAddress) != "verified@example.com" || aws.ToString(client.input.Content.Simple.Subject.Data) != "Order update" || len(client.input.EmailTags) != 1 {
		t.Fatalf("unexpected official input: %+v", client.input)
	}
	if receipt.ProviderID != "provider-1" || receipt.ContentSHA256 == "" || receipt.RecipientCount != 1 {
		t.Fatalf("unexpected receipt: %+v", receipt)
	}
}

func TestSendTransactionalEmailFailsClosed(t *testing.T) {
	cases := []SendEmailRequest{
		{},
		{From: "bad", To: []string{"ok@example.com"}, Subject: "s", TextBody: "b", MessageKey: "k"},
		{From: "a@example.com", To: []string{"dup@example.com", "dup@example.com"}, Subject: "s", TextBody: "b", MessageKey: "k"},
		{From: "a@example.com", To: []string{"b@example.com"}, Subject: "", TextBody: "b", MessageKey: "k"},
	}
	for _, request := range cases {
		if _, err := SendTransactionalEmail(context.Background(), &fakeSES{}, request); err == nil {
			t.Fatalf("invalid email accepted: %+v", request)
		}
	}
	if _, err := SendTransactionalEmail(context.Background(), &fakeSES{output: &sesv2.SendEmailOutput{}}, SendEmailRequest{From: "a@example.com", To: []string{"b@example.com"}, Subject: "s", TextBody: "b", MessageKey: "k"}); err == nil {
		t.Fatal("missing provider ID accepted")
	}
}
````

### FILE: `aws_enterprise_adapters/README.md`
```yaml
block_id: "GO-AWS-ENTERPRISE-ADAPTERS:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "6d7f1785da8c804eff881fd4070777f545d8246200366514d62055cbb4d525f7"
variables: []
secrets_allowed: false
```
````markdown
# AWS enterprise storage and email adapters

This module uses exact official AWS SDK for Go v2 releases: S3 `v1.107.3`, SES v2 `v1.67.0`, and config `v1.32.38`.

`PutImmutableObject` first proves bucket Object Lock through `GetObjectLockConfiguration`, binding the expected 12-digit bucket-owner account. It sends a seekable byte payload through `s3.PutObject` with a precomputed SHA-256 checksum, `If-None-Match: *`, explicit server-side encryption and explicit GOVERNANCE or COMPLIANCE retention. It requires ETag and VersionId, then verifies the exact created version with `GetObjectRetention`. The receipt contains owner/version/checksum/retention and marks verification. If the post-write verification fails, it returns both an error and the created-version receipt so the caller can persist/reconcile the partial provider effect instead of retrying blindly. This is a single-put adapter with a conservative 64 MiB limit; larger objects need the official transfer manager and a separate multipart/checksum evidence lane. Legal hold is deliberately excluded from this path and requires separate explicit authority.

`SendTransactionalEmail` sends simple UTF-8 text/HTML through `sesv2.SendEmail`, requires a verified sender candidate and an external outbox message key, adds that key as an SES tag, and returns a receipt without recipient addresses or body content. SES has no claim here of provider-side idempotency: the enterprise outbox must prevent duplicate effects and process delivery/bounce/complaint events.

`NewClients` constructs both official service clients with `config.LoadDefaultConfig`, an explicit region, and the AWS default credential chain. Secrets are never request fields or receipts. Before a real call, complete region/IAM, expected bucket owner, bucket/versioning/Object Lock, retention owner/mode/date, lifecycle/encryption, SES verified identity/sandbox/configuration set/suppression/event route, privacy/residency, cost/quota/load/retry and recovery decisions. This code never enables Object Lock: AWS documents that enabling it is permanent for the bucket, so infrastructure mutation remains an explicit project operation.

```text
go mod verify
GOPROXY=off go test ./...
```
````

## 6. Configuration surface

| Variable | Tipo/default | Validación | Secreto | Efecto |
|---|---|---|---|---|
| AWS region/credentials | region + default chain | IAM mínimo y reference externa | sí | endpoint/identidad |
| S3 bucket/key | strings | bucket no vacío, key relativo | posible PII en key: prohibir | objeto destino |
| encryption/KMS key | enum/reference | AES256 sin key o aws:kms con key | key ID no secreto | cifrado |
| SES from/config set | email/string | identity verificada/config existente | no | sender/events |
| message key | string | 1..256, generado por outbox | no | deduplicación externa/tag |
| recipients/content | data | email válido, 1..50, tamaño conservador | PII/contenido | envío; no receipt |

El profile exige versioning/object-lock/lifecycle y suppression/bounce/outbox/costo antes de habilitar.

## 7. Dependency bill

| Package/tool | Pin | Uso | Licencia | Fase | Fuente |
|---|---|---|---|---|---|
| AWS config | v1.32.38 | default credential chain | Apache-2.0 | runtime | aws/aws-sdk-go-v2 |
| AWS S3 | v1.107.3 | `PutObject`/checksums | Apache-2.0 | runtime | aws/aws-sdk-go-v2 |
| AWS SES v2 | v1.67.0 | `SendEmail` | Apache-2.0 | runtime | aws/aws-sdk-go-v2 |
| Go | 1.26.7 | build/tests | BSD-3-Clause | build/runtime | go.dev |

`go.sum` fija transitivas. LICENSE/NOTICE y commit AWS exacto están en el upstream lock.

## 8. Apply order

Materializar, adquirir módulos exactos, verificar/tests offline, completar profile, crear IAM/bucket/identity/config set fuera del pack, probar sandbox y recién integrar mediante outbox/storage service. En existente, compositor rechaza colisiones y la migración añade adapters sin cambiar ownership de dominio.

Rollback: deshabilitar route/config, conservar receipts/object versions, pausar outbox, restaurar provider anterior y reconciliar efectos. No borrar objetos ni reenviar email automáticamente.

## 9. Verification

```text
go mod verify
GOPROXY=off go test -count=1 ./...
go list -m github.com/aws/aws-sdk-go-v2/config github.com/aws/aws-sdk-go-v2/service/s3 github.com/aws/aws-sdk-go-v2/service/sesv2
```

Éxito: 6 tests PASS; checksum/encryption/no-overwrite, preflight Object Lock, owner, retención exacta por VersionId, efecto parcial visible y SES input/tag/receipt/negativos; módulos exactos, vet y build PASS. Pendiente: sandbox, IAM, bucket controls/lifecycle, malware/restore, SES identity/events/delivery, privacy/costo/quota/load/retry/reconciliation.

## 10. Reconstruction evidence

- Go oficial 1.26.7 windows/amd64;
- 9/9 archivos reconstruibles/hash-verificados;
- módulos config v1.32.38, S3 v1.107.3, SES v2 v1.67.0;
- `go mod verify` + tests offline PASS;
- sin llamada AWS/credenciales/corpus;
- glue/tests `AUTHORED`, SDK Apache-2.0 adquirido como dependencia;
- 2026-08-26 / Codex; expediente `GO_AWS_ENTERPRISE_STORAGE_EMAIL_ADAPTERS_2026-08-26_V1.md`.
