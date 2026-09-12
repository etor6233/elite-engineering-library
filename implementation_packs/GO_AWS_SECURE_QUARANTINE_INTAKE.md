# Go AWS Secure Quarantine Intake

## 1. Metadata

```yaml
pack_id: "GO-AWS-SECURE-QUARANTINE-INTAKE"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa una sesión Go de carga portal/API hacia cuarentena S3 mediante presigned PUT del SDK oficial AWS, con checksum/tamaño/MIME/KMS/owner/create-only firmados, estado server-side y verificación de la versión exacta antes del gate de seguridad."
stacks: ["Go 1.26.7", "AWS SDK for Go v2", "Amazon S3"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "GO-AWS-ENTERPRISE-STORAGE-EMAIL-ADAPTERS 0.2.x"]
incompatible_with: ["anonymous API", "public bucket", "permissive CORS", "permanent credentials", "overwrite", "business storage before security", "caller-owned session state"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND BSD-3-Clause"
upstream_sources:
  - "https://github.com/aws/aws-sdk-go-v2/tree/284a4846e7bb941926e2d15144b17c01ffc98c1d/service/s3"
  - "https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/go_s3_code_examples.html"
  - "https://docs.aws.amazon.com/AmazonS3/latest/userguide/PresignedUrlUploadObject.html"
  - "https://docs.aws.amazon.com/AmazonS3/latest/developerguide/bucket-policy-s3-sigv4-conditions.html"
  - "https://docs.aws.amazon.com/prescriptive-guidance/latest/presigned-url-best-practices/additional-guardrails.html"
  - "https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html"
verified_at: "2026-08-28"
```

El SDK y las guías son oficiales de AWS. Los seis archivos materializados son glue/configuración/tests `AUTHORED`; no se atribuyen a Amazon ni se presentan como código verbatim suyo.

## 2. Applicability

Use cuando una aplicación autenticada necesite entregar un archivo desde portal/API directamente a un bucket S3 privado de cuarentena sin pasar los bytes por el proceso web. El caller ya debe haber resuelto identidad/autorización y entrega sólo una referencia opaca SHA-256 del subject. La sesión emitida se persiste server-side antes de devolver la URL.

Rechazar si faltan cuenta/IAM/costo, credenciales temporales, Block Public Access, KMS, versioning, policy `s3:signatureAge`, orígenes CORS exactos, bucket/prefix de cuarentena aislado, allowlists de clase/MIME/tamaño o un store transaccional con compare-and-set. No cubre multipart, email/SFTP/scanner/mobile, creación de infraestructura, evento S3, antimalware, extracción, verdad de campos ni almacenamiento de negocio.

## 3. Architecture contract

`Issue` valida configuración y request, genera 128 bits aleatorios, fija un key bajo `quarantine/<class>/<upload-id>`, llama `PresignPutObject` con `ContentLength`, `ContentType`, `ChecksumSHA256`, owner esperado, SSE-KMS, `If-None-Match: *` y metadata de binding. Sólo acepta HTTPS+PUT, persiste `ISSUED` y devuelve URL más todos los headers firmados; URL/headers no se registran ni almacenan.

`Finalize` recibe sólo upload ID, version ID y ETag del response S3. Recupera la sesión desde store confiable, ejecuta `HeadObject` con `ChecksumMode=ENABLED` sobre esa versión y compara identidad, bytes, checksum, MIME, KMS y metadata. El receipt se publica sólo si `Complete` hace CAS de `ISSUED` a `QUARANTINED_VERIFIED`; duplicados, timeout, mismatch o store failure quedan bloqueados y se reconcilian, nunca se reintenta una carga ciegamente.

Este receipt sólo prueba recepción esperada en cuarentena. La siguiente transición obligatoria es `SECURE-LOCAL-FILE-INGESTION-GATE`; no autoriza clasificación, extracción ni persistencia empresarial. El target debe añadir expiración/abandono de sesiones, event delivery at-least-once, reconciliación, métricas sin URL/PII, límites, carga, recovery y rollback.

## 4. Exact file manifest

```text
CREATE aws_secure_quarantine_intake/go.mod
CREATE aws_secure_quarantine_intake/go.sum
CREATE aws_secure_quarantine_intake/provider-profile.template.json
CREATE aws_secure_quarantine_intake/awsintake/intake.go
CREATE aws_secure_quarantine_intake/awsintake/intake_test.go
CREATE aws_secure_quarantine_intake/README.md
```

## 5. Materialization blocks

### FILE: `aws_secure_quarantine_intake/go.mod`

```yaml
block_id: "GO-AWS-SECURE-QUARANTINE-INTAKE:gomod:v1"
operation: CREATE
provenance: AUTHORED
source: "local module pinned to official AWS SDK for Go v2"
license: "LicenseRef-Workspace-Owner"
sha256: "abf577684fa2df625266739854c1bafdd6a04f7fbada2395c18594fcf4e6c37e"
variables: []
secrets_allowed: false
```

````text
module example.com/elite/aws-secure-quarantine-intake

go 1.26.0

require (
	github.com/aws/aws-sdk-go-v2 v1.43.8
	github.com/aws/aws-sdk-go-v2/service/s3 v1.107.3
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.38 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.38 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.39 // indirect
	github.com/aws/smithy-go v1.27.10 // indirect
)
````

### FILE: `aws_secure_quarantine_intake/go.sum`

```yaml
block_id: "GO-AWS-SECURE-QUARANTINE-INTAKE:gosum:v1"
operation: CREATE
provenance: AUTHORED
source: "Go checksum lock resolved from exact module pins"
license: "LicenseRef-Workspace-Owner"
sha256: "25ae7474cd904d50162f0ba82565294d83020cbe3b24496ba9137b9fa9abeced"
variables: []
secrets_allowed: false
```

````text
github.com/aws/aws-sdk-go-v2 v1.43.8 h1:fpnrxwuwsoGIgjvgLeDU3y9w7YaHBxyF6AF3vQL8duw=
github.com/aws/aws-sdk-go-v2 v1.43.8/go.mod h1:j7gYSq8dL95QejkFXxvQNESH4I9WGHFI6iO+vhqEi5Q=
github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.18 h1:LAfOuhAH331fmOjTQpAaOlH+Ftn7RzSDJ2VFwjdMMy4=
github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.18/go.mod h1:4e5xhuXHx1e4U9EthvbPP1r/DIMp5c2823OL8karzcM=
github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.38 h1:MBMg0zJ6i4TkAJ0dVFLKKn2cOkY6FkicmUDM67BRr6g=
github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.38/go.mod h1:9MWuJbyiUyj6eA7W1/zm1zuePDPSB3g+xcgRQeMWsXc=
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.38 h1:lHm4jPf3k1Lz5ZWc+Vcn3MKVwym+26kWCba9FkJ4f0Y=
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.38/go.mod h1:Rn+P2XR+FbyZzjmWKjg/KUZNxmGfr5oZwh5jQiE+CzI=
github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39 h1:vo4xvMRs/F6h1E52qsgLqCQgWIQXgIJUauG6rlZEh4U=
github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39/go.mod h1:jB03R1ij/A+OE2e1dz6vgj076gd7vlYcfstAzj3HcnU=
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17 h1:OvYZOB3qA6zvfdRFiRFRzVSiElMYrz3GdntkXZxlp1o=
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17/go.mod h1:JgR/2Ew50ACfIWau1oeMRX59tMtC0kM+PYQGEaT04cY=
github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.31 h1:uZOinZb+h7lZw8IYzP1z1IuEnueB76/EFkcf/fEW4Ag=
github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.31/go.mod h1:NRtwAM/p5VRt03TlEUs0pH3TeWamWdf4YyJpSrzPYLc=
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38 h1:H/5TI1jqaHsNoDQ60UwvPvJBg4GURkinXI3Qga29t2w=
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38/go.mod h1:PTVFf+XH++7NJOky+RLBYQx0QA5NcaeEYFQ2fsi0nwo=
github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.39 h1:HLPAVrlLDaN2boN0xJx7MgaQDNEO3Q+c9L6kl/8m47Q=
github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.39/go.mod h1:Pg/dVfsNkm1hsIDK/gMvCKtmyNfNTV12mrgHqVE/6Oo=
github.com/aws/aws-sdk-go-v2/service/s3 v1.107.3 h1:IKoCZqfWfZzSBi16QFQ+QcbQ3LRQ7QgB1S5tDAyPBQQ=
github.com/aws/aws-sdk-go-v2/service/s3 v1.107.3/go.mod h1:RBpRcXiM4s2pOInVs32GsBonnje+fiAj4mcrStRmlCA=
github.com/aws/smithy-go v1.27.10 h1:bw56MIx8bhTQZSdzucEJSKWLpwX0ju7hU8cVoa75dg8=
github.com/aws/smithy-go v1.27.10/go.mod h1:YE2RhdIuDbA5E5bTdciG9KrW3+TiEONeUWCqxX9i1Fc=
````

### FILE: `aws_secure_quarantine_intake/provider-profile.template.json`

```yaml
block_id: "GO-AWS-SECURE-QUARANTINE-INTAKE:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed configuration constrained by AWS authorities"
license: "LicenseRef-Workspace-Owner"
sha256: "c169ed6d1f1f1101c1fcfba9bcf35d7a8dfc7dae2d1ef40cf8dd187a46e75497"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-aws-quarantine-intake-profile/v1",
  "enabled": false,
  "decision": "BLOCKED_ACCESS_COST_AND_POLICY_REQUIRED",
  "aws": {
    "account_id": "",
    "region": "",
    "temporary_credentials_verified": false,
    "quarantine_bucket": "",
    "quarantine_prefix": "quarantine/",
    "expected_bucket_owner": "",
    "kms_key_arn": "",
    "bucket_versioning_verified": false,
    "block_public_access_verified": false,
    "signature_age_policy_verified": false,
    "signature_age_max_seconds": 600,
    "quarantine_boundary_verified": false,
    "cors_exact_origins": [],
    "cors_exposed_headers": ["ETag", "x-amz-version-id", "x-amz-checksum-sha256"]
  },
  "upload": {
    "max_bytes": 0,
    "max_ttl_seconds": 600,
    "allowed_content_types": [],
    "allowed_document_classes": [],
    "malware_gate_route": "",
    "automatic_business_storage_authorized": false
  },
  "approvals": {
    "account_iam_cost": false,
    "bucket_kms_network": false,
    "content_policy": false,
    "privacy_retention": false,
    "sandbox_execution": false
  }
}
````

### FILE: `aws_secure_quarantine_intake/awsintake/intake.go`

```yaml
block_id: "GO-AWS-SECURE-QUARANTINE-INTAKE:intake:v1"
operation: CREATE
provenance: AUTHORED
source: "local glue invoking official AWS S3 SDK public APIs and official guardrails"
license: "LicenseRef-Workspace-Owner"
sha256: "96d40a72f85121626c7f2ada2e1db1e6ea973939b209452f6b35ed121ddc4039"
variables: []
secrets_allowed: false
```

````go
package awsintake

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	safeClass = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)
	ownerID   = regexp.MustCompile(`^[0-9]{12}$`)
	hexHash   = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

type PresignAPI interface {
	PresignPutObject(context.Context, *s3.PutObjectInput, ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type HeadAPI interface {
	HeadObject(context.Context, *s3.HeadObjectInput, ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

type SessionStore interface {
	Create(context.Context, Session) error
	Get(context.Context, string) (Session, error)
	Complete(context.Context, string, Receipt) error
}

type Config struct {
	Bucket                       string
	Prefix                       string
	ExpectedBucketOwner          string
	KMSKeyARN                    string
	MaxBytes                     int64
	MaxTTL                       time.Duration
	AllowedContentTypes          map[string]struct{}
	AllowedDocumentClasses       map[string]struct{}
	SignatureAgePolicyVerified   bool
	TemporaryCredentialsVerified bool
	QuarantineBoundaryVerified   bool
}

type IssueRequest struct {
	SubjectRefSHA256 string
	DocumentClass    string
	ContentType      string
	ContentLength    int64
	ChecksumSHA256   string
	TTL              time.Duration
}

type Session struct {
	UploadID         string    `json:"upload_id"`
	ObjectKey        string    `json:"object_key"`
	SubjectRefSHA256 string    `json:"subject_ref_sha256"`
	DocumentClass    string    `json:"document_class"`
	ContentType      string    `json:"content_type"`
	ContentLength    int64     `json:"content_length"`
	ChecksumSHA256   string    `json:"checksum_sha256"`
	IssuedAt         time.Time `json:"issued_at"`
	ExpiresAt        time.Time `json:"expires_at"`
	State            string    `json:"state"`
}

type UploadInstructions struct {
	Session Session     `json:"session"`
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
}

type FinalizeRequest struct {
	UploadID  string
	VersionID string
	ETag      string
}

type Receipt struct {
	UploadID       string    `json:"upload_id"`
	ObjectKey      string    `json:"object_key"`
	VersionID      string    `json:"version_id"`
	ETag           string    `json:"etag"`
	ChecksumSHA256 string    `json:"checksum_sha256"`
	ContentLength  int64     `json:"content_length"`
	VerifiedAt     time.Time `json:"verified_at"`
	State          string    `json:"state"`
}

type Service struct {
	Config    Config
	Presigner PresignAPI
	Head      HeadAPI
	Store     SessionStore
	Random    io.Reader
	Now       func() time.Time
}

func (s Service) Issue(ctx context.Context, request IssueRequest) (UploadInstructions, error) {
	if err := s.validateConfig(); err != nil {
		return UploadInstructions{}, err
	}
	if s.Presigner == nil {
		return UploadInstructions{}, errors.New("presigner is required")
	}
	if s.Store == nil {
		return UploadInstructions{}, errors.New("session store is required")
	}
	if !hexHash.MatchString(request.SubjectRefSHA256) {
		return UploadInstructions{}, errors.New("subject reference must be a lowercase SHA-256")
	}
	if !safeClass.MatchString(request.DocumentClass) {
		return UploadInstructions{}, errors.New("document class is invalid")
	}
	if _, ok := s.Config.AllowedDocumentClasses[request.DocumentClass]; !ok {
		return UploadInstructions{}, errors.New("document class is not approved")
	}
	if _, ok := s.Config.AllowedContentTypes[request.ContentType]; !ok {
		return UploadInstructions{}, errors.New("content type is not approved")
	}
	if request.ContentLength < 1 || request.ContentLength > s.Config.MaxBytes {
		return UploadInstructions{}, errors.New("content length is outside the approved limit")
	}
	checksumBytes, err := base64.StdEncoding.Strict().DecodeString(request.ChecksumSHA256)
	if err != nil || len(checksumBytes) != 32 {
		return UploadInstructions{}, errors.New("checksum must be a canonical base64 SHA-256")
	}
	if request.TTL < time.Minute || request.TTL > s.Config.MaxTTL {
		return UploadInstructions{}, errors.New("TTL is outside the approved range")
	}
	random := s.Random
	if random == nil {
		random = rand.Reader
	}
	idBytes := make([]byte, 16)
	if _, err := io.ReadFull(random, idBytes); err != nil {
		return UploadInstructions{}, fmt.Errorf("generate upload id: %w", err)
	}
	uploadID := hex.EncodeToString(idBytes)
	objectKey := s.Config.Prefix + request.DocumentClass + "/" + uploadID
	now := s.now()
	metadata := map[string]string{
		"upload-id":          uploadID,
		"subject-ref-sha256": request.SubjectRefSHA256,
		"document-class":     request.DocumentClass,
	}
	presigned, err := s.Presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(s.Config.Bucket),
		Key:                  aws.String(objectKey),
		ContentLength:        aws.Int64(request.ContentLength),
		ContentType:          aws.String(request.ContentType),
		ChecksumSHA256:       aws.String(request.ChecksumSHA256),
		ExpectedBucketOwner:  aws.String(s.Config.ExpectedBucketOwner),
		ServerSideEncryption: types.ServerSideEncryptionAwsKms,
		SSEKMSKeyId:          aws.String(s.Config.KMSKeyARN),
		IfNoneMatch:          aws.String("*"),
		Metadata:             metadata,
	}, func(options *s3.PresignOptions) {
		options.Expires = request.TTL
	})
	if err != nil {
		return UploadInstructions{}, fmt.Errorf("presign quarantine upload: %w", err)
	}
	parsed, err := url.Parse(presigned.URL)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || presigned.Method != "PUT" {
		return UploadInstructions{}, errors.New("presigner returned a non-HTTPS PUT request")
	}
	session := Session{
		UploadID: uploadID, ObjectKey: objectKey, SubjectRefSHA256: request.SubjectRefSHA256,
		DocumentClass: request.DocumentClass, ContentType: request.ContentType,
		ContentLength: request.ContentLength, ChecksumSHA256: request.ChecksumSHA256,
		IssuedAt: now, ExpiresAt: now.Add(request.TTL), State: "ISSUED",
	}
	if err := s.Store.Create(ctx, session); err != nil {
		return UploadInstructions{}, fmt.Errorf("persist upload session: %w", err)
	}
	return UploadInstructions{Session: session, Method: presigned.Method, URL: presigned.URL, Headers: presigned.SignedHeader.Clone()}, nil
}

func (s Service) Finalize(ctx context.Context, request FinalizeRequest) (Receipt, error) {
	if err := s.validateConfig(); err != nil {
		return Receipt{}, err
	}
	if s.Head == nil {
		return Receipt{}, errors.New("head client is required")
	}
	if s.Store == nil {
		return Receipt{}, errors.New("session store is required")
	}
	if request.UploadID == "" {
		return Receipt{}, errors.New("upload ID is required")
	}
	session, err := s.Store.Get(ctx, request.UploadID)
	if err != nil {
		return Receipt{}, fmt.Errorf("load upload session: %w", err)
	}
	if session.State != "ISSUED" || session.UploadID != request.UploadID {
		return Receipt{}, errors.New("session is not issuable")
	}
	if request.VersionID == "" || request.ETag == "" {
		return Receipt{}, errors.New("version ID and ETag returned by S3 are required")
	}
	if !strings.HasPrefix(session.ObjectKey, s.Config.Prefix) {
		return Receipt{}, errors.New("session object escaped quarantine prefix")
	}
	if s.now().After(session.ExpiresAt.Add(5 * time.Minute)) {
		return Receipt{}, errors.New("session finalization window expired")
	}
	output, err := s.Head.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket:              aws.String(s.Config.Bucket),
		Key:                 aws.String(session.ObjectKey),
		VersionId:           aws.String(request.VersionID),
		ExpectedBucketOwner: aws.String(s.Config.ExpectedBucketOwner),
		ChecksumMode:        types.ChecksumModeEnabled,
	})
	if err != nil {
		return Receipt{}, fmt.Errorf("verify uploaded version: %w", err)
	}
	if output.VersionId == nil || *output.VersionId != request.VersionID || output.ETag == nil || *output.ETag != request.ETag {
		return Receipt{}, errors.New("S3 version identity mismatch")
	}
	if output.ContentLength == nil || *output.ContentLength != session.ContentLength || output.ContentType == nil || *output.ContentType != session.ContentType {
		return Receipt{}, errors.New("S3 content metadata mismatch")
	}
	if output.ChecksumSHA256 == nil || *output.ChecksumSHA256 != session.ChecksumSHA256 {
		return Receipt{}, errors.New("S3 SHA-256 checksum mismatch")
	}
	if output.ServerSideEncryption != types.ServerSideEncryptionAwsKms || output.SSEKMSKeyId == nil || *output.SSEKMSKeyId != s.Config.KMSKeyARN {
		return Receipt{}, errors.New("S3 encryption identity mismatch")
	}
	if output.Metadata["upload-id"] != session.UploadID || output.Metadata["subject-ref-sha256"] != session.SubjectRefSHA256 || output.Metadata["document-class"] != session.DocumentClass {
		return Receipt{}, errors.New("S3 request-binding metadata mismatch")
	}
	receipt := Receipt{
		UploadID: session.UploadID, ObjectKey: session.ObjectKey,
		VersionID: request.VersionID, ETag: request.ETag,
		ChecksumSHA256: session.ChecksumSHA256, ContentLength: session.ContentLength,
		VerifiedAt: s.now(), State: "QUARANTINED_VERIFIED",
	}
	if err := s.Store.Complete(ctx, session.UploadID, receipt); err != nil {
		return Receipt{}, fmt.Errorf("complete upload session: %w", err)
	}
	return receipt, nil
}

func (s Service) validateConfig() error {
	c := s.Config
	if !c.SignatureAgePolicyVerified || !c.TemporaryCredentialsVerified || !c.QuarantineBoundaryVerified {
		return errors.New("AWS upload guardrails are not verified")
	}
	if c.Bucket == "" || !ownerID.MatchString(c.ExpectedBucketOwner) || !strings.HasPrefix(c.KMSKeyARN, "arn:aws:kms:") {
		return errors.New("AWS storage identity is invalid")
	}
	if c.Prefix == "" || !strings.HasSuffix(c.Prefix, "/") || strings.HasPrefix(c.Prefix, "/") || strings.Contains(c.Prefix, "..") || strings.Contains(c.Prefix, "\\") {
		return errors.New("quarantine prefix is invalid")
	}
	if c.MaxBytes < 1 || c.MaxTTL < time.Minute || c.MaxTTL > 15*time.Minute || len(c.AllowedContentTypes) == 0 || len(c.AllowedDocumentClasses) == 0 {
		return errors.New("upload policy is invalid")
	}
	for contentType := range c.AllowedContentTypes {
		if contentType == "" || contentType != strings.ToLower(contentType) || strings.Contains(contentType, "*") {
			return errors.New("content type allowlist must be exact and lowercase")
		}
	}
	for documentClass := range c.AllowedDocumentClasses {
		if !safeClass.MatchString(documentClass) {
			return errors.New("document class allowlist contains an invalid value")
		}
	}
	return nil
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}
````

### FILE: `aws_secure_quarantine_intake/awsintake/intake_test.go`

```yaml
block_id: "GO-AWS-SECURE-QUARANTINE-INTAKE:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline contract and negative tests over official AWS SDK types"
license: "LicenseRef-Workspace-Owner"
sha256: "a887df43d9c169184f82aabedd4e9b0b8c23713f2cdafec61e7058ea5132f7d8"
variables: []
secrets_allowed: false
```

````go
package awsintake

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type fakePresigner struct {
	input   *s3.PutObjectInput
	expires time.Duration
	err     error
}

func (f *fakePresigner) PresignPutObject(_ context.Context, input *s3.PutObjectInput, options ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	f.input = input
	configured := s3.PresignOptions{}
	for _, option := range options {
		option(&configured)
	}
	f.expires = configured.Expires
	if f.err != nil {
		return nil, f.err
	}
	return &v4.PresignedHTTPRequest{
		Method: "PUT", URL: "https://quarantine.example.test/object?signature=redacted",
		SignedHeader: http.Header{"Content-Type": []string{aws.ToString(input.ContentType)}, "X-Amz-Checksum-Sha256": []string{aws.ToString(input.ChecksumSHA256)}},
	}, nil
}

type fakeHead struct {
	input  *s3.HeadObjectInput
	output *s3.HeadObjectOutput
	err    error
}

func (f *fakeHead) HeadObject(_ context.Context, input *s3.HeadObjectInput, _ ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	f.input = input
	return f.output, f.err
}

type fakeStore struct {
	sessions    map[string]Session
	createErr   error
	getErr      error
	completeErr error
	completed   *Receipt
}

func (f *fakeStore) Create(_ context.Context, session Session) error {
	if f.createErr != nil {
		return f.createErr
	}
	if _, exists := f.sessions[session.UploadID]; exists {
		return errors.New("duplicate")
	}
	f.sessions[session.UploadID] = session
	return nil
}

func (f *fakeStore) Get(_ context.Context, uploadID string) (Session, error) {
	if f.getErr != nil {
		return Session{}, f.getErr
	}
	session, ok := f.sessions[uploadID]
	if !ok {
		return Session{}, errors.New("not found")
	}
	return session, nil
}

func (f *fakeStore) Complete(_ context.Context, uploadID string, receipt Receipt) error {
	if f.completeErr != nil {
		return f.completeErr
	}
	session, ok := f.sessions[uploadID]
	if !ok || session.State != "ISSUED" {
		return errors.New("state conflict")
	}
	session.State = "QUARANTINED_VERIFIED"
	f.sessions[uploadID] = session
	f.completed = &receipt
	return nil
}

func validService() (Service, *fakePresigner, *fakeHead, *fakeStore) {
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	presigner := &fakePresigner{}
	head := &fakeHead{}
	store := &fakeStore{sessions: map[string]Session{}}
	return Service{
		Config: Config{
			Bucket: "elite-quarantine", Prefix: "quarantine/", ExpectedBucketOwner: "123456789012",
			KMSKeyARN: "arn:aws:kms:us-east-1:123456789012:key/00000000-0000-0000-0000-000000000001",
			MaxBytes:  10 << 20, MaxTTL: 10 * time.Minute,
			AllowedContentTypes:        map[string]struct{}{"application/pdf": {}},
			AllowedDocumentClasses:     map[string]struct{}{"invoice": {}},
			SignatureAgePolicyVerified: true, TemporaryCredentialsVerified: true, QuarantineBoundaryVerified: true,
		},
		Presigner: presigner, Head: head, Store: store, Random: bytes.NewReader(make([]byte, 16)), Now: func() time.Time { return now },
	}, presigner, head, store
}

func validIssueRequest() IssueRequest {
	return IssueRequest{
		SubjectRefSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		DocumentClass:    "invoice", ContentType: "application/pdf", ContentLength: 4096,
		ChecksumSHA256: base64.StdEncoding.EncodeToString(make([]byte, 32)), TTL: 5 * time.Minute,
	}
}

func TestIssueBindsQuarantineUpload(t *testing.T) {
	service, presigner, _, store := validService()
	issued, err := service.Issue(context.Background(), validIssueRequest())
	if err != nil {
		t.Fatal(err)
	}
	if issued.Method != "PUT" || issued.Session.State != "ISSUED" || presigner.expires != 5*time.Minute {
		t.Fatalf("unexpected issue result: %#v", issued)
	}
	if aws.ToString(presigner.input.Bucket) != service.Config.Bucket || aws.ToString(presigner.input.ExpectedBucketOwner) != service.Config.ExpectedBucketOwner {
		t.Fatal("bucket owner was not bound")
	}
	if aws.ToString(presigner.input.IfNoneMatch) != "*" || presigner.input.ServerSideEncryption != types.ServerSideEncryptionAwsKms || aws.ToString(presigner.input.SSEKMSKeyId) != service.Config.KMSKeyARN {
		t.Fatal("create-only KMS constraints were not bound")
	}
	if aws.ToString(presigner.input.ChecksumSHA256) != validIssueRequest().ChecksumSHA256 || aws.ToInt64(presigner.input.ContentLength) != 4096 {
		t.Fatal("content checksum or length was not bound")
	}
	if presigner.input.Metadata["upload-id"] != issued.Session.UploadID || issued.Headers.Get("Content-Type") != "application/pdf" {
		t.Fatal("request binding metadata or signed headers missing")
	}
	if store.sessions[issued.Session.UploadID].ObjectKey != issued.Session.ObjectKey {
		t.Fatal("server-side session was not persisted")
	}
}

func TestIssueFailsClosed(t *testing.T) {
	tests := map[string]func(*Service, *IssueRequest){
		"guardrail":    func(s *Service, _ *IssueRequest) { s.Config.SignatureAgePolicyVerified = false },
		"content-type": func(_ *Service, r *IssueRequest) { r.ContentType = "text/html" },
		"size":         func(_ *Service, r *IssueRequest) { r.ContentLength = 0 },
		"checksum":     func(_ *Service, r *IssueRequest) { r.ChecksumSHA256 = "not-a-sha" },
		"subject":      func(_ *Service, r *IssueRequest) { r.SubjectRefSHA256 = "person@example.com" },
		"class":        func(_ *Service, r *IssueRequest) { r.DocumentClass = "../invoice" },
		"unapproved":   func(_ *Service, r *IssueRequest) { r.DocumentClass = "packing-list" },
		"ttl":          func(_ *Service, r *IssueRequest) { r.TTL = 30 * time.Second },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			service, _, _, _ := validService()
			request := validIssueRequest()
			mutate(&service, &request)
			if _, err := service.Issue(context.Background(), request); err == nil {
				t.Fatal("expected fail-closed error")
			}
		})
	}
}

func TestIssuePropagatesPresignerFailure(t *testing.T) {
	service, presigner, _, _ := validService()
	presigner.err = errors.New("provider down")
	if _, err := service.Issue(context.Background(), validIssueRequest()); err == nil {
		t.Fatal("expected provider error")
	}
}

func TestFinalizeVerifiesExactVersion(t *testing.T) {
	service, _, head, store := validService()
	issued, err := service.Issue(context.Background(), validIssueRequest())
	if err != nil {
		t.Fatal(err)
	}
	head.output = &s3.HeadObjectOutput{
		VersionId: aws.String("version-1"), ETag: aws.String("\"etag-1\""),
		ContentLength: aws.Int64(issued.Session.ContentLength), ContentType: aws.String(issued.Session.ContentType),
		ChecksumSHA256: aws.String(issued.Session.ChecksumSHA256), ServerSideEncryption: types.ServerSideEncryptionAwsKms,
		SSEKMSKeyId: aws.String(service.Config.KMSKeyARN),
		Metadata:    map[string]string{"upload-id": issued.Session.UploadID, "subject-ref-sha256": issued.Session.SubjectRefSHA256, "document-class": issued.Session.DocumentClass},
	}
	receipt, err := service.Finalize(context.Background(), FinalizeRequest{UploadID: issued.Session.UploadID, VersionID: "version-1", ETag: "\"etag-1\""})
	if err != nil {
		t.Fatal(err)
	}
	if receipt.State != "QUARANTINED_VERIFIED" || aws.ToString(head.input.VersionId) != "version-1" || head.input.ChecksumMode != types.ChecksumModeEnabled {
		t.Fatalf("unexpected receipt or verification input: %#v", receipt)
	}
	if store.completed == nil || store.sessions[issued.Session.UploadID].State != "QUARANTINED_VERIFIED" {
		t.Fatal("server-side completion was not atomic")
	}
}

func TestFinalizeRejectsMismatchAndBlindRetry(t *testing.T) {
	service, _, head, _ := validService()
	issued, err := service.Issue(context.Background(), validIssueRequest())
	if err != nil {
		t.Fatal(err)
	}
	head.output = &s3.HeadObjectOutput{
		VersionId: aws.String("version-1"), ETag: aws.String("\"etag-1\""),
		ContentLength: aws.Int64(issued.Session.ContentLength), ContentType: aws.String(issued.Session.ContentType),
		ChecksumSHA256: aws.String("different"), ServerSideEncryption: types.ServerSideEncryptionAwsKms,
		SSEKMSKeyId: aws.String(service.Config.KMSKeyARN), Metadata: map[string]string{"upload-id": issued.Session.UploadID},
	}
	if _, err := service.Finalize(context.Background(), FinalizeRequest{UploadID: issued.Session.UploadID, VersionID: "version-1", ETag: "\"etag-1\""}); err == nil {
		t.Fatal("expected checksum mismatch rejection")
	}
	head.err = errors.New("timeout")
	if _, err := service.Finalize(context.Background(), FinalizeRequest{UploadID: issued.Session.UploadID, VersionID: "version-1", ETag: "\"etag-1\""}); err == nil {
		t.Fatal("expected unknown provider effect to remain unverified")
	}
}

func TestStoreFailuresAndDuplicateCompletionFailClosed(t *testing.T) {
	service, _, head, store := validService()
	store.createErr = errors.New("database down")
	if _, err := service.Issue(context.Background(), validIssueRequest()); err == nil {
		t.Fatal("expected session persistence failure")
	}
	store.createErr = nil
	service.Random = bytes.NewReader(bytes.Repeat([]byte{1}, 16))
	issued, err := service.Issue(context.Background(), validIssueRequest())
	if err != nil {
		t.Fatal(err)
	}
	head.output = &s3.HeadObjectOutput{
		VersionId: aws.String("version-1"), ETag: aws.String("\"etag-1\""),
		ContentLength: aws.Int64(issued.Session.ContentLength), ContentType: aws.String(issued.Session.ContentType),
		ChecksumSHA256: aws.String(issued.Session.ChecksumSHA256), ServerSideEncryption: types.ServerSideEncryptionAwsKms,
		SSEKMSKeyId: aws.String(service.Config.KMSKeyARN),
		Metadata:    map[string]string{"upload-id": issued.Session.UploadID, "subject-ref-sha256": issued.Session.SubjectRefSHA256, "document-class": issued.Session.DocumentClass},
	}
	request := FinalizeRequest{UploadID: issued.Session.UploadID, VersionID: "version-1", ETag: "\"etag-1\""}
	if _, err := service.Finalize(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Finalize(context.Background(), request); err == nil {
		t.Fatal("duplicate completion must fail")
	}
}
````

### FILE: `aws_secure_quarantine_intake/README.md`

```yaml
block_id: "GO-AWS-SECURE-QUARANTINE-INTAKE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration guide constrained by AWS authorities"
license: "LicenseRef-Workspace-Owner"
sha256: "da69d74b660f6bc914d20c8f525231f0d00596a6e3001f4b9c1e1a7995156085"
variables: []
secrets_allowed: false
```

````markdown
# AWS secure quarantine intake

This module adds one bounded intake lane: an authenticated application may issue a short-lived S3 presigned `PUT` for a private quarantine bucket and later verify the exact created version before the file enters the local security gate.

It uses the official AWS SDK for Go v2. The service glue, profile, and tests are locally authored and are not represented as Amazon source code.

The profile is disabled. Before enabling it, prove account/IAM/cost authority, temporary credentials, Block Public Access, versioning, KMS, exact CORS origins, a bucket-policy `s3:signatureAge` guardrail, content policy, quarantine isolation, event delivery, malware routing, reconciliation, retention, monitoring, load, and rollback.

The returned URL and signed headers are ephemeral and must not be persisted or logged. Persist the `Session` in the application database, return the instructions only to the authenticated subject, and require the S3 response `ETag` and `x-amz-version-id` for `Finalize`. A successful receipt means only that S3 stored the expected bytes in quarantine. It does not mean the document is safe, classified, extracted, correct, or authorized for business storage.
````

## 6. Configuration surface

| Variable | Tipo/default | Validación | Secreto | Efecto |
|---|---|---|---|---|
| enabled/decision | false/BLOCKED | sólo habilitar tras approvals y probes | no | fail-closed global |
| account/region/credentials | references | cuenta esperada + credenciales temporales | referencia sí | autoridad AWS |
| bucket/prefix/owner | strings | privados, prefix relativo, owner 12 dígitos | no; nombres pueden ser sensibles | cuarentena exacta |
| KMS key ARN | ARN | key aprobada y accesible | no | cifrado firmado/verificado |
| signature age | 600 s | policy real probada; TTL 60..900 s | no | reduce replay |
| CORS origins/headers | listas | origins exactos; expone ETag/version/checksum | no | browser finalize |
| max bytes/MIME/classes | allowlists | valores exactos, sin wildcard | no | admission previa |
| SessionStore | interfaz | create/get/complete CAS y durable | datos operativos | request binding/idempotencia |

No hay credenciales, endpoints firmados ni datos personales en el profile. `automatic_business_storage_authorized=false` es invariante.

## 7. Dependency bill

| Package/tool | Pin | Uso | Licencia | Fase | Fuente oficial |
|---|---|---|---|---|---|
| AWS SDK for Go v2 core | v1.43.8 | tipos/signing | Apache-2.0 | runtime | aws/aws-sdk-go-v2 |
| AWS S3 | v1.107.3 | presign/head/checksum/KMS | Apache-2.0 | runtime | aws/aws-sdk-go-v2/service/s3 |
| smithy-go y módulos internos | `go.sum` exacto | transporte/protocolo | Apache-2.0 | runtime | AWS |
| Go | 1.26.7 | build/tests | BSD-3-Clause | build/runtime | go.dev/dl |

## 8. Apply order

Materializar en workspace vacío; ejecutar verify/tests offline después de adquirir el graph; completar profile; implementar `SessionStore` sobre el backend transaccional; probar bucket/IAM/KMS/policy/CORS en sandbox; conectar endpoint autenticado; conectar versión verificada al gate local de seguridad; añadir eventos/reconciliación y recién habilitar por configuración.

En workspace existente, el compositor rechaza colisiones. Rollback: deshabilitar emisión, conservar sesiones/receipts/versiones, cancelar sesiones no entregadas, reconciliar uploads con efecto desconocido y revertir route; nunca borrar objetos ni marcar seguridad aprobada.

## 9. Verification

```text
go mod verify
GOPROXY=off go test -count=1 ./...
go vet ./...
go build ./...
```

Éxito local: 6 tests con 8 subtests PASS; checksum/size/MIME/KMS/owner/create-only, HTTPS PUT, allowlists, persist-before-return, exact version HEAD, metadata binding, store failures y duplicate CAS están cubiertos. Go vet y build PASS. Pendiente y bloqueante: sandbox AWS, IAM negative, bucket/KMS/policy/CORS probes, event duplicates/order, malware route, DB store real, concurrency/load, observability, abandoned sessions, recovery/rollback y costo.

## 10. Reconstruction evidence

- staging limpio con Go oficial 1.26.7 Windows amd64, archive 74.955.002 bytes/SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`;
- AWS SDK core v1.43.8, S3 v1.107.3 y graph fijado por `go.sum`;
- seis archivos creados desde staging; materialización y comparación SHA se ejecutan después de insertar los bloques;
- 6 tests/8 subtests, vet y build PASS; cero llamadas AWS/credenciales;
- código/config/tests `AUTHORED`, SDK oficial sólo dependencia Apache-2.0;
- Azure Event Grid commit oficial `0a705781…` se investigó y compiló aparte, pero no se incorpora a este pack AWS ni hereda su admisión;
- 2026-08-28 / Codex.
