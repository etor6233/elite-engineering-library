# Go AWS Textract Document Runtime

## 1. Metadata

```yaml
pack_id: "GO-AWS-TEXTRACT-DOCUMENT-RUNTIME"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un runtime Go condicionado sobre el SDK oficial AWS v2 Textract 1.45.0; sólo ejecuta una clase REQUIRED del perfil aprobado y enlaza perfil, recibo de seguridad local, input y respuesta completa mediante SHA-256 y commit atómico."
stacks: ["Go 1.26.7", "AWS SDK for Go v2", "service/textract v1.45.0", "Amazon Textract"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION 0.1.x", "GO-ENTERPRISE-BACKEND-CORE 0.4.x"]
incompatible_with: ["operación/features/queries/adapter/región libres por CLI", "credenciales por CLI", "S3 obligatorio", "auto-storage sin corpus", "adapter sin versión", "documento síncrono mayor a 10 MiB", "input sin recibo de seguridad ADMITTED"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/aws/aws-sdk-go-v2/tree/a30468cff35d6e385287a0cbff0ac11aa7202529/service/textract", "https://proxy.golang.org/github.com/aws/aws-sdk-go-v2/service/textract/@v/v1.45.0.info", "https://docs.aws.amazon.com/textract/latest/APIReference/API_AnalyzeExpense.html", "https://docs.aws.amazon.com/textract/latest/APIReference/API_AnalyzeDocument.html", "https://docs.aws.amazon.com/textract/latest/dg/limits-document.html"]
verified_at: "2026-08-27"
```

## 2. Applicability

Use cuando el blueprint seleccione Amazon Textract y exista cuenta/región/IAM/costo aprobados. El perfil fija por clase `AnalyzeExpense` o `AnalyzeDocument`, feature types oficiales, Queries y opcionalmente AdapterId+Version exactos. El ejecutable no recibe esos valores libremente. Rechazar cuando se pretenda que un resultado no evaluado alimente la base empresarial, falte el recibo de seguridad, falte corpus/ground truth o una clase custom no tenga queries/adapter versionado.

El código es glue `AUTHORED` sobre tipos y clientes públicos del SDK AWS; servicio, modelos y respuesta son Amazon. No copia el SDK. Proformas, órdenes, packing lists, BOL, delivery notes, certificados y aduana permanecen bloqueados hasta completar la ruta oficial de Queries/Adapters con documentos representativos, entrenamiento/test y métricas.

## 3. Architecture contract

El runner exige una clase que aparezca exactamente una vez como `REQUIRED` en `elite-aws-textract-document-profile/v2`. El perfil fija región, SDK, operación, features, queries/alias, adapter exacto, schema de respuesta y `automatic_storage=false`; el template no habilita ninguna clase. Antes de crear el cliente, la CLI obtiene la región sólo de ese perfil. Antes de llamar AWS, el runtime verifica un receipt `elite-secure-local-file-receipt/v1` `ADMITTED`: approval/policy, hash/tamaño/MIME del input y evidencias completas ClamAV/Magika/YARA-X.

El input debe ser archivo regular no symlink, extensión coherente con el MIME detectado y tamaño dentro del límite síncrono AWS de 10 MiB. Envía `types.Document.Bytes`, conserva el output modelado completo como JSON, exige una página probada y emite receipt V2 con hashes de perfil, seguridad, input y respuesta. Autenticación usa la default credential chain del SDK; no recibe claves/tokens por CLI.

Falla antes de commit ante input/config/API/respuesta/output inválidos; staging se elimina y el destino final debe ser nuevo. No persiste campos en una base ni traduce confidence/blocks/expense fields a hechos. Rendimiento, precisión, cuota, costo, residencia y retry/throttle se prueban en el target. Rollback fija el adapter anterior o retira la ruta, conserva input/output/receipt y reevalúa mappings antes de reprocesar.

## 4. Exact file manifest

```text
CREATE aws_textract_runtime/go.mod
CREATE aws_textract_runtime/go.sum
CREATE aws_textract_runtime/document-profile.template.json
CREATE aws_textract_runtime/textractruntime/runtime.go
CREATE aws_textract_runtime/textractruntime/runtime_test.go
CREATE aws_textract_runtime/cmd/analyze/main.go
CREATE aws_textract_runtime/README.md
```

## 5. Materialization blocks

### FILE: `aws_textract_runtime/go.mod`
```yaml
block_id: "GO-AWS-TEXTRACT-RUNTIME:gomod:v2"
operation: CREATE
provenance: AUTHORED
source: "exact official AWS Go SDK module versions"
license: "LicenseRef-Workspace-Owner"
sha256: "591b58441d93dc9f65bbedd678437613d79cb6a93e4509522eb547bd4ba1f895"
variables: []
secrets_allowed: false
```
````text
module example.com/elite/aws-textract-document-runtime

go 1.26.0

require (
	github.com/aws/aws-sdk-go-v2 v1.44.0
	github.com/aws/aws-sdk-go-v2/config v1.32.38
	github.com/aws/aws-sdk-go-v2/service/textract v1.45.0
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.19.37 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.38 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.40 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.40 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.5.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.33.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.45.7 // indirect
	github.com/aws/smithy-go v1.28.1 // indirect
)
````

### FILE: `aws_textract_runtime/go.sum`
```yaml
block_id: "GO-AWS-TEXTRACT-RUNTIME:gosum:v2"
operation: CREATE
provenance: AUTHORED
source: "Go checksum database entries resolved from exact direct pins"
license: "LicenseRef-Workspace-Owner"
sha256: "903f85b219ba58e553bb7b6bafcc2b34b1593b7fba91699ede6b5998ad42d708"
variables: []
secrets_allowed: false
```
````text
github.com/aws/aws-sdk-go-v2 v1.44.0 h1:4IbaHhtzy+4h37z4JQyO9a2QsiCml3CNYHtq5hIHigo=
github.com/aws/aws-sdk-go-v2 v1.44.0/go.mod h1:bttEH6JqnUL8LepvDVfdrds/fZ5bCIxzpe3abyUrhDU=
github.com/aws/aws-sdk-go-v2/config v1.32.38 h1:n4yPHBjtQ3BrIIUyk0/LAqf/BL2iv0Tw6XZcMRzM0ps=
github.com/aws/aws-sdk-go-v2/config v1.32.38/go.mod h1:dencYsOS1R7rBy8zehCvwBYzdxxL4Q/nRK7In03wjN8=
github.com/aws/aws-sdk-go-v2/credentials v1.19.37 h1:FJ8Iz4/xISMB/rwLlgfWujfGDFWr0oneQgtA6KPcYLY=
github.com/aws/aws-sdk-go-v2/credentials v1.19.37/go.mod h1:Q6pWOgVUp49x4g5QVi29wHofUoICnZ+Zq4jHbRN/7ec=
github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.38 h1:Nqo2jU1wz5rnBM9XQyXfVD1RP8txkbP3EDx8hR/hbCE=
github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.38/go.mod h1:PzJFHhjR2vWFKHe8HmY5Lxhvwyxnr5MERtk0nDxWNbk=
github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.40 h1:UIXlbijuB2XK1Kr57fo8iIxCuaSHJzwZ1uo+2tbEYIk=
github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.40/go.mod h1:wcEsL6jscjZjVUinb0Q5qD/GXOG1yT3GNfmT9HuDwzU=
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.40 h1:xLQVRDs2NddDmK9BEyh5KSlJ1Gpy5/GIJXrV6WcVGAE=
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.40/go.mod h1:XRXnpFVFGLaEVK+olDdFIM1vNa04ETW452oFGEPUxAo=
github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39 h1:vo4xvMRs/F6h1E52qsgLqCQgWIQXgIJUauG6rlZEh4U=
github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39/go.mod h1:jB03R1ij/A+OE2e1dz6vgj076gd7vlYcfstAzj3HcnU=
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17 h1:OvYZOB3qA6zvfdRFiRFRzVSiElMYrz3GdntkXZxlp1o=
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17/go.mod h1:JgR/2Ew50ACfIWau1oeMRX59tMtC0kM+PYQGEaT04cY=
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38 h1:H/5TI1jqaHsNoDQ60UwvPvJBg4GURkinXI3Qga29t2w=
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38/go.mod h1:PTVFf+XH++7NJOky+RLBYQx0QA5NcaeEYFQ2fsi0nwo=
github.com/aws/aws-sdk-go-v2/service/signin v1.5.7 h1:YcczQ6zNH/ojIzD/ikDrO+RfW06wmdMp18d4NH5hXY4=
github.com/aws/aws-sdk-go-v2/service/signin v1.5.7/go.mod h1:nl9RVnb9ulgAYzOkjLq1NyFxmWcnH2maCUEuOdESy98=
github.com/aws/aws-sdk-go-v2/service/sso v1.33.7 h1:P+bMNiA93gyuYT3Oh+4dWtvrnGcu2bd9Uy5hRJM8BNo=
github.com/aws/aws-sdk-go-v2/service/sso v1.33.7/go.mod h1:zy+397isDFLvleg9H18Zq2MGzMso7uKyJyzR7DWSgFk=
github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.7 h1:WWkehGZ4nWtOKLMy0yi8+RqzzVqAGe60hGaxwF06JAw=
github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.7/go.mod h1:T8AI4SbQYm9ybcVmki2T3n7Qg1g3kfWoeQlNwNYOyO8=
github.com/aws/aws-sdk-go-v2/service/sts v1.45.7 h1:yU/9y2r7s9kSUPbHXbpQTa4LA8kt+CMgpu1OBrhx8p4=
github.com/aws/aws-sdk-go-v2/service/sts v1.45.7/go.mod h1:0lQTDEBArMevQXpxu443LVGjKxxEeSsSnrw9n8YiTMg=
github.com/aws/aws-sdk-go-v2/service/textract v1.45.0 h1:es1kFEIsARdt9US5dt5C6N4Nk8YXax5b75UTpDrrnIY=
github.com/aws/aws-sdk-go-v2/service/textract v1.45.0/go.mod h1:xFfDl2uBgGT5AoRgDjSzbZg1BZfqk8NlkIMOJaHulZU=
github.com/aws/smithy-go v1.28.1 h1:R/nXH00c8qcfCzQVELtRw+eLQWtzv+VAIEFJ1/xxXlQ=
github.com/aws/smithy-go v1.28.1/go.mod h1:YE2RhdIuDbA5E5bTdciG9KrW3+TiEONeUWCqxX9i1Fc=
````

### FILE: `aws_textract_runtime/document-profile.template.json`
```yaml
block_id: "GO-AWS-TEXTRACT-RUNTIME:profile:v2"
operation: CREATE
provenance: AUTHORED
source: "fail-closed mapping to AWS documented operations/features"
license: "LicenseRef-Workspace-Owner"
sha256: "6dd28e5cc1cd4e9bd20ffbc2c7d78b7ca423a09071ee4df99a194824e88c07e7"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-aws-textract-document-profile/v2",
  "provider": "Amazon Textract",
  "sdk": "github.com/aws/aws-sdk-go-v2/service/textract@v1.45.0",
  "region": "CONFIGURATION_REQUIRED",
  "classes": [
    {"id":"supplier-invoice","decision":"CONFIGURATION_REQUIRED","operation":"AnalyzeExpense","features":[],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"expense-receipt","decision":"BLOCKED_EVALUATION_REQUIRED","operation":"AnalyzeExpense","features":[],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"commercial-invoice","decision":"BLOCKED_EVALUATION_REQUIRED","operation":"AnalyzeExpense","features":[],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"proforma-invoice","decision":"BLOCKED_QUERIES_ADAPTER_REQUIRED","operation":"AnalyzeDocument","features":["QUERIES"],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"purchase-order","decision":"BLOCKED_QUERIES_ADAPTER_REQUIRED","operation":"AnalyzeDocument","features":["TABLES","FORMS","QUERIES"],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"packing-list","decision":"BLOCKED_QUERIES_ADAPTER_REQUIRED","operation":"AnalyzeDocument","features":["TABLES","FORMS","QUERIES"],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"bill-of-lading","decision":"BLOCKED_QUERIES_ADAPTER_REQUIRED","operation":"AnalyzeDocument","features":["TABLES","FORMS","QUERIES"],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"delivery-note","decision":"BLOCKED_QUERIES_ADAPTER_REQUIRED","operation":"AnalyzeDocument","features":["TABLES","FORMS","QUERIES"],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"certificate-of-origin","decision":"BLOCKED_QUERIES_ADAPTER_REQUIRED","operation":"AnalyzeDocument","features":["FORMS","QUERIES"],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"quality-inspection-certificate","decision":"BLOCKED_QUERIES_ADAPTER_REQUIRED","operation":"AnalyzeDocument","features":["TABLES","FORMS","QUERIES"],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false},
    {"id":"customs-declaration","decision":"BLOCKED_QUERIES_ADAPTER_REQUIRED","operation":"AnalyzeDocument","features":["TABLES","FORMS","QUERIES"],"queries":[],"adapter_id":"","adapter_version":"","response_schema_version":"CONFIGURATION_REQUIRED","automatic_storage":false}
  ]
}
````

### FILE: `aws_textract_runtime/textractruntime/runtime.go`
```yaml
block_id: "GO-AWS-TEXTRACT-RUNTIME:runtime:v2"
operation: CREATE
provenance: AUTHORED
source: "local wrapper invoking AWS SDK for Go v2 Textract public API"
license: "LicenseRef-Workspace-Owner"
sha256: "34ec475f9145b1e717f4fe1f929a6ec6f2e236b7074f00d306b64b4c8b5d7b88"
variables: []
secrets_allowed: false
```
````go
package textractruntime

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

const (
	ModeExpense           = "expense"
	ModeDocument          = "document"
	MaxSyncBytes          = 10 * 1024 * 1024
	profileSchema         = "elite-aws-textract-document-profile/v2"
	securityReceiptSchema = "elite-secure-local-file-receipt/v1"
	sdkModule             = "github.com/aws/aws-sdk-go-v2/service/textract"
	sdkVersion            = "v1.45.0"
)

var allowedMIMETypes = map[string]string{
	".pdf": "application/pdf", ".png": "image/png", ".jpg": "image/jpeg", ".jpeg": "image/jpeg",
	".tif": "image/tiff", ".tiff": "image/tiff",
}
var featureTypes = map[string]types.FeatureType{
	"TABLES": types.FeatureTypeTables, "FORMS": types.FeatureTypeForms, "QUERIES": types.FeatureTypeQueries,
	"SIGNATURES": types.FeatureTypeSignatures, "LAYOUT": types.FeatureTypeLayout,
}

type Client interface {
	AnalyzeExpense(context.Context, *textract.AnalyzeExpenseInput, ...func(*textract.Options)) (*textract.AnalyzeExpenseOutput, error)
	AnalyzeDocument(context.Context, *textract.AnalyzeDocumentInput, ...func(*textract.Options)) (*textract.AnalyzeDocumentOutput, error)
}

type Query struct {
	Text  string   `json:"text"`
	Alias string   `json:"alias"`
	Pages []string `json:"pages,omitempty"`
}

type Request struct {
	Region              string
	DocumentClass       string
	InputPath           string
	OutputDirectory     string
	ProfilePath         string
	SecurityReceiptPath string
	MaxBytes            int64
}

type Receipt struct {
	Schema                     string `json:"schema"`
	CreatedAt                  string `json:"created_at"`
	Provider                   string `json:"provider"`
	SDKModule                  string `json:"sdk_module"`
	SDKVersion                 string `json:"sdk_version"`
	Operation                  string `json:"operation"`
	Region                     string `json:"region"`
	DocumentClass              string `json:"document_class"`
	ResponseSchemaVersion      string `json:"response_schema_version"`
	DocumentProfileSHA256      string `json:"document_profile_sha256"`
	SecurityReceiptSHA256      string `json:"security_receipt_sha256"`
	ContentType                string `json:"content_type"`
	InputFilename              string `json:"input_filename"`
	InputBytes                 int    `json:"input_bytes"`
	InputSHA256                string `json:"input_sha256"`
	ProviderResponseSHA256     string `json:"provider_response_sha256"`
	PageCount                  int32  `json:"page_count"`
	ResultCount                int    `json:"result_count"`
	ModelVersion               string `json:"model_version,omitempty"`
	AdapterID                  string `json:"adapter_id,omitempty"`
	AdapterVersion             string `json:"adapter_version,omitempty"`
	AutomaticStorageAuthorized bool   `json:"automatic_storage_authorized"`
}

type documentProfile struct {
	Schema   string         `json:"schema"`
	Provider string         `json:"provider"`
	SDK      string         `json:"sdk"`
	Region   string         `json:"region"`
	Classes  []profileClass `json:"classes"`
}

type profileClass struct {
	ID                    string   `json:"id"`
	Decision              string   `json:"decision"`
	Operation             string   `json:"operation"`
	Features              []string `json:"features"`
	Queries               []Query  `json:"queries"`
	AdapterID             string   `json:"adapter_id"`
	AdapterVersion        string   `json:"adapter_version"`
	ResponseSchemaVersion string   `json:"response_schema_version"`
	AutomaticStorage      bool     `json:"automatic_storage"`
}

type securityReceipt struct {
	Schema       string `json:"schema"`
	Decision     string `json:"decision"`
	Reason       string `json:"reason"`
	StartedAt    string `json:"started_at"`
	CompletedAt  string `json:"completed_at"`
	ApprovalID   string `json:"approval_id"`
	PolicySHA256 string `json:"policy_sha256"`
	Input        struct {
		SHA256 string `json:"sha256"`
		Bytes  int    `json:"bytes"`
	} `json:"input"`
	Tools map[string]struct {
		SHA256  string `json:"sha256"`
		Version string `json:"version"`
	} `json:"tools"`
	BusinessStorageAuthorized *bool  `json:"business_storage_authorized"`
	SecurityClaim             string `json:"security_claim"`
	ClamAV                    struct {
		ExitCode     int    `json:"exit_code"`
		OutputSHA256 string `json:"output_sha256"`
	} `json:"clamav"`
	ContentType struct {
		Label    string  `json:"label"`
		MIMEType string  `json:"mime_type"`
		Score    float64 `json:"score"`
	} `json:"content_type"`
	YaraX struct {
		CompiledRulesSHA256 string   `json:"compiled_rules_sha256"`
		MatchingRules       []string `json:"matching_rules"`
	} `json:"yara_x"`
}

type approvedSelection struct {
	Region                string
	DocumentClass         string
	Mode                  string
	Features              []types.FeatureType
	Queries               []types.Query
	Adapter               *types.Adapter
	AdapterID             string
	AdapterVersion        string
	ResponseSchemaVersion string
	ProfileSHA256         string
}

func ApprovedRegion(profilePath, documentClass string) (string, error) {
	selection, err := selectProfile(profilePath, documentClass, "")
	if err != nil {
		return "", err
	}
	return selection.Region, nil
}

func AnalyzeToEvidence(ctx context.Context, client Client, request Request) (receipt Receipt, err error) {
	if client == nil {
		return receipt, errors.New("client is required")
	}
	if request.MaxBytes == 0 {
		request.MaxBytes = MaxSyncBytes
	}
	if request.MaxBytes < 1 || request.MaxBytes > MaxSyncBytes {
		return receipt, errors.New("max_bytes must be within the Textract synchronous 10 MiB limit")
	}
	selection, err := selectProfile(request.ProfilePath, request.DocumentClass, request.Region)
	if err != nil {
		return receipt, err
	}
	info, err := os.Lstat(request.InputPath)
	if err != nil {
		return receipt, err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return receipt, errors.New("input must be a regular non-symlink file")
	}
	absInput, err := filepath.Abs(request.InputPath)
	if err != nil {
		return receipt, err
	}
	mimeType, ok := allowedMIMETypes[strings.ToLower(filepath.Ext(absInput))]
	if !ok {
		return receipt, errors.New("unsupported document extension")
	}
	payload, err := os.ReadFile(absInput)
	if err != nil {
		return receipt, err
	}
	if len(payload) == 0 || int64(len(payload)) > request.MaxBytes {
		return receipt, errors.New("document byte size is outside the allowed range")
	}
	payloadSHA := hash(payload)
	securitySHA, err := verifySecurityReceipt(request.SecurityReceiptPath, payloadSHA, len(payload), mimeType)
	if err != nil {
		return receipt, err
	}
	absOutput, err := filepath.Abs(request.OutputDirectory)
	if err != nil {
		return receipt, err
	}
	if _, statErr := os.Stat(absOutput); statErr == nil {
		return receipt, errors.New("output_directory must not exist")
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return receipt, statErr
	}
	if err = os.MkdirAll(filepath.Dir(absOutput), 0o750); err != nil {
		return receipt, err
	}
	stage, err := os.MkdirTemp(filepath.Dir(absOutput), ".textract-analysis-stage-")
	if err != nil {
		return receipt, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(stage)
		}
	}()

	var output any
	var pages int32
	var resultCount int
	var modelVersion string
	switch selection.Mode {
	case ModeExpense:
		result, callErr := client.AnalyzeExpense(ctx, &textract.AnalyzeExpenseInput{Document: &types.Document{Bytes: payload}})
		if callErr != nil {
			return receipt, callErr
		}
		if result == nil || len(result.ExpenseDocuments) == 0 {
			return receipt, errors.New("provider response contains no expense documents")
		}
		output, resultCount = result, len(result.ExpenseDocuments)
		if result.DocumentMetadata != nil {
			pages = aws.ToInt32(result.DocumentMetadata.Pages)
		}
	case ModeDocument:
		input := &textract.AnalyzeDocumentInput{Document: &types.Document{Bytes: payload}, FeatureTypes: selection.Features}
		if len(selection.Queries) > 0 {
			input.QueriesConfig = &types.QueriesConfig{Queries: selection.Queries}
		}
		if selection.Adapter != nil {
			input.AdaptersConfig = &types.AdaptersConfig{Adapters: []types.Adapter{*selection.Adapter}}
		}
		result, callErr := client.AnalyzeDocument(ctx, input)
		if callErr != nil {
			return receipt, callErr
		}
		if result == nil || len(result.Blocks) == 0 {
			return receipt, errors.New("provider response contains no document blocks")
		}
		output, resultCount, modelVersion = result, len(result.Blocks), aws.ToString(result.AnalyzeDocumentModelVersion)
		if result.DocumentMetadata != nil {
			pages = aws.ToInt32(result.DocumentMetadata.Pages)
		}
	}
	if pages != 1 {
		return receipt, fmt.Errorf("synchronous provider response must prove exactly one page; got %d", pages)
	}
	responseBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return receipt, err
	}
	responseBytes = append(responseBytes, '\n')
	if err = os.WriteFile(filepath.Join(stage, "provider-response.json"), responseBytes, 0o600); err != nil {
		return receipt, err
	}
	receipt = Receipt{
		Schema: "elite-aws-textract-analysis-receipt/v2", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano), Provider: "Amazon Textract",
		SDKModule: sdkModule, SDKVersion: sdkVersion, Operation: map[string]string{ModeExpense: "AnalyzeExpense", ModeDocument: "AnalyzeDocument"}[selection.Mode], Region: selection.Region,
		DocumentClass: selection.DocumentClass, ResponseSchemaVersion: selection.ResponseSchemaVersion, DocumentProfileSHA256: selection.ProfileSHA256, SecurityReceiptSHA256: securitySHA, ContentType: mimeType,
		InputFilename: filepath.Base(absInput), InputBytes: len(payload), InputSHA256: payloadSHA, ProviderResponseSHA256: hash(responseBytes), PageCount: pages, ResultCount: resultCount, ModelVersion: modelVersion,
		AdapterID: selection.AdapterID, AdapterVersion: selection.AdapterVersion, AutomaticStorageAuthorized: false,
	}
	receiptBytes, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return receipt, err
	}
	receiptBytes = append(receiptBytes, '\n')
	if err = os.WriteFile(filepath.Join(stage, "ANALYSIS_RECEIPT.json"), receiptBytes, 0o600); err != nil {
		return receipt, err
	}
	if err = os.Rename(stage, absOutput); err != nil {
		return receipt, err
	}
	committed = true
	return receipt, nil
}

func selectProfile(path, documentClass, expectedRegion string) (approvedSelection, error) {
	var profile documentProfile
	body, err := readStrictJSON(path, "document profile", &profile)
	if err != nil {
		return approvedSelection{}, err
	}
	if profile.Schema != profileSchema || profile.Provider != "Amazon Textract" || profile.SDK != sdkModule+"@"+sdkVersion {
		return approvedSelection{}, errors.New("document profile identity is not admitted")
	}
	if profile.Region == "" || profile.Region == "CONFIGURATION_REQUIRED" || strings.ContainsAny(profile.Region, " /\\") {
		return approvedSelection{}, errors.New("document profile requires an exact AWS region")
	}
	if expectedRegion != "" && expectedRegion != profile.Region {
		return approvedSelection{}, errors.New("runtime region does not match document profile")
	}
	documentClass = strings.TrimSpace(documentClass)
	if documentClass == "" {
		return approvedSelection{}, errors.New("document class is required")
	}
	matches := make([]profileClass, 0, 1)
	for _, item := range profile.Classes {
		if item.ID == documentClass {
			matches = append(matches, item)
		}
	}
	if len(matches) != 1 {
		return approvedSelection{}, errors.New("document class must appear exactly once in profile")
	}
	item := matches[0]
	if item.Decision != "REQUIRED" {
		return approvedSelection{}, fmt.Errorf("document class %s is not REQUIRED", documentClass)
	}
	if item.AutomaticStorage {
		return approvedSelection{}, errors.New("document profile must deny automatic storage")
	}
	if item.ResponseSchemaVersion == "" || item.ResponseSchemaVersion == "CONFIGURATION_REQUIRED" {
		return approvedSelection{}, errors.New("exact response schema version is required")
	}
	mode, features, queries, adapter, err := validateAnalysis(item)
	if err != nil {
		return approvedSelection{}, err
	}
	return approvedSelection{Region: profile.Region, DocumentClass: documentClass, Mode: mode, Features: features, Queries: queries, Adapter: adapter, AdapterID: item.AdapterID, AdapterVersion: item.AdapterVersion, ResponseSchemaVersion: item.ResponseSchemaVersion, ProfileSHA256: hash(body)}, nil
}

func validateAnalysis(item profileClass) (string, []types.FeatureType, []types.Query, *types.Adapter, error) {
	operation := strings.TrimSpace(item.Operation)
	if operation == "AnalyzeExpense" {
		if len(item.Features) > 0 || len(item.Queries) > 0 || item.AdapterID != "" || item.AdapterVersion != "" {
			return "", nil, nil, nil, errors.New("AnalyzeExpense profile cannot contain document features, queries, or adapters")
		}
		return ModeExpense, nil, nil, nil, nil
	}
	if operation != "AnalyzeDocument" {
		return "", nil, nil, nil, errors.New("profile operation must be AnalyzeExpense or AnalyzeDocument")
	}
	if len(item.Features) == 0 {
		return "", nil, nil, nil, errors.New("AnalyzeDocument profile requires at least one feature")
	}
	features := make([]types.FeatureType, 0, len(item.Features))
	seen := map[string]struct{}{}
	for _, raw := range item.Features {
		name := strings.ToUpper(strings.TrimSpace(raw))
		value, ok := featureTypes[name]
		if !ok {
			return "", nil, nil, nil, fmt.Errorf("unsupported feature %q", raw)
		}
		if _, dup := seen[name]; dup {
			return "", nil, nil, nil, fmt.Errorf("duplicate feature %q", name)
		}
		seen[name] = struct{}{}
		features = append(features, value)
	}
	hasQueries := slices.Contains(features, types.FeatureTypeQueries)
	if hasQueries != (len(item.Queries) > 0) {
		return "", nil, nil, nil, errors.New("QUERIES feature and non-empty profile queries must be supplied together")
	}
	if len(item.Queries) > 15 {
		return "", nil, nil, nil, errors.New("synchronous Textract profile exceeds 15 queries per page")
	}
	queries := make([]types.Query, 0, len(item.Queries))
	aliases := map[string]struct{}{}
	for _, query := range item.Queries {
		text, alias := strings.TrimSpace(query.Text), strings.TrimSpace(query.Alias)
		if text == "" || alias == "" {
			return "", nil, nil, nil, errors.New("query text and alias are required")
		}
		if _, dup := aliases[alias]; dup {
			return "", nil, nil, nil, fmt.Errorf("duplicate query alias %q", alias)
		}
		for _, page := range query.Pages {
			if page != "1" {
				return "", nil, nil, nil, errors.New("synchronous profile query pages may only contain page 1")
			}
		}
		aliases[alias] = struct{}{}
		queries = append(queries, types.Query{Text: aws.String(text), Alias: aws.String(alias), Pages: query.Pages})
	}
	if (item.AdapterID == "") != (item.AdapterVersion == "") {
		return "", nil, nil, nil, errors.New("adapter ID and exact version must be supplied together")
	}
	var adapter *types.Adapter
	if item.AdapterID != "" {
		if !hasQueries {
			return "", nil, nil, nil, errors.New("an adapter requires the QUERIES feature")
		}
		if item.AdapterID == "CONFIGURATION_REQUIRED" || item.AdapterVersion == "CONFIGURATION_REQUIRED" {
			return "", nil, nil, nil, errors.New("adapter identity is not configured")
		}
		adapter = &types.Adapter{AdapterId: aws.String(item.AdapterID), Version: aws.String(item.AdapterVersion)}
	}
	return ModeDocument, features, queries, adapter, nil
}

func verifySecurityReceipt(path, payloadSHA string, payloadBytes int, mimeType string) (string, error) {
	var receipt securityReceipt
	body, err := readStrictJSON(path, "security receipt", &receipt)
	if err != nil {
		return "", err
	}
	if receipt.Schema != securityReceiptSchema {
		return "", errors.New("unsupported security receipt schema")
	}
	if receipt.Decision != "ADMITTED" || receipt.Reason != "ALL_SELECTED_GATES_PASSED" {
		return "", errors.New("security receipt did not admit the input")
	}
	if receipt.BusinessStorageAuthorized == nil || *receipt.BusinessStorageAuthorized {
		return "", errors.New("security receipt must explicitly deny business storage")
	}
	if receipt.ApprovalID == "" || !isSHA256(receipt.PolicySHA256) {
		return "", errors.New("security receipt approval or policy evidence is invalid")
	}
	if _, e := time.Parse(time.RFC3339Nano, receipt.StartedAt); e != nil {
		return "", errors.New("security receipt started_at is invalid")
	}
	if _, e := time.Parse(time.RFC3339Nano, receipt.CompletedAt); e != nil {
		return "", errors.New("security receipt completed_at is invalid")
	}
	if receipt.SecurityClaim == "" {
		return "", errors.New("security receipt claim is missing")
	}
	if receipt.Input.SHA256 != payloadSHA || receipt.Input.Bytes != payloadBytes {
		return "", errors.New("security receipt does not match input bytes")
	}
	if receipt.ContentType.MIMEType != mimeType || receipt.ContentType.Label == "" || receipt.ContentType.Score < 0 || receipt.ContentType.Score > 1 {
		return "", errors.New("security receipt content type does not match input")
	}
	for _, name := range []string{"clamscan", "magika", "yara_x"} {
		tool, ok := receipt.Tools[name]
		if !ok || !isSHA256(tool.SHA256) || tool.Version == "" {
			return "", fmt.Errorf("security receipt %s evidence is invalid", name)
		}
	}
	if receipt.ClamAV.ExitCode != 0 || !isSHA256(receipt.ClamAV.OutputSHA256) {
		return "", errors.New("security receipt ClamAV result is invalid")
	}
	if !isSHA256(receipt.YaraX.CompiledRulesSHA256) || receipt.YaraX.MatchingRules == nil || len(receipt.YaraX.MatchingRules) != 0 {
		return "", errors.New("security receipt YARA-X result is invalid")
	}
	return hash(body), nil
}

func readStrictJSON(path, label string, target any) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("%s must be a regular non-symlink file", label)
	}
	if info.Size() < 2 || info.Size() > 1024*1024 {
		return nil, fmt.Errorf("%s size is invalid", label)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(target); err != nil {
		return nil, fmt.Errorf("invalid %s: %w", label, err)
	}
	var trailing any
	if err = decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("invalid %s: trailing JSON value", label)
	}
	return body, nil
}

func isSHA256(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil && value == strings.ToLower(value)
}
func hash(value []byte) string { sum := sha256.Sum256(value); return hex.EncodeToString(sum[:]) }
````

### FILE: `aws_textract_runtime/textractruntime/runtime_test.go`
```yaml
block_id: "GO-AWS-TEXTRACT-RUNTIME:test:v2"
operation: CREATE
provenance: AUTHORED
source: "local contract/negative tests against official SDK types"
license: "LicenseRef-Workspace-Owner"
sha256: "a4f74d95884d2f50eb136afae6e1842fc83cc4dd445b16f0138a660db5afad54"
variables: []
secrets_allowed: false
```
````go
package textractruntime

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

const testHex = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

type fakeClient struct {
	expense       *textract.AnalyzeExpenseOutput
	document      *textract.AnalyzeDocumentOutput
	err           error
	expenseInput  *textract.AnalyzeExpenseInput
	documentInput *textract.AnalyzeDocumentInput
	calls         int
}

func (f *fakeClient) AnalyzeExpense(_ context.Context, input *textract.AnalyzeExpenseInput, _ ...func(*textract.Options)) (*textract.AnalyzeExpenseOutput, error) {
	f.calls++
	f.expenseInput = input
	return f.expense, f.err
}
func (f *fakeClient) AnalyzeDocument(_ context.Context, input *textract.AnalyzeDocumentInput, _ ...func(*textract.Options)) (*textract.AnalyzeDocumentOutput, error) {
	f.calls++
	f.documentInput = input
	return f.document, f.err
}

type fixtureState struct{ root, input, profile, security string }

func newFixture(t *testing.T, item profileClass) fixtureState {
	t.Helper()
	root := t.TempDir()
	input := filepath.Join(root, "invoice.pdf")
	if err := os.WriteFile(input, []byte("%PDF-1.7 fixture"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := fixtureState{root: root, input: input, profile: filepath.Join(root, "profile.json"), security: filepath.Join(root, "security.json")}
	writeProfile(t, state.profile, item)
	writeSecurity(t, state.security, state.input, "ADMITTED", "application/pdf", true)
	return state
}

func expenseClass() profileClass {
	return profileClass{ID: "supplier-invoice", Decision: "REQUIRED", Operation: "AnalyzeExpense", Features: []string{}, Queries: []Query{}, ResponseSchemaVersion: "supplier-invoice/v1", AutomaticStorage: false}
}

func writeProfile(t *testing.T, path string, item profileClass) {
	t.Helper()
	profile := documentProfile{Schema: profileSchema, Provider: "Amazon Textract", SDK: sdkModule + "@" + sdkVersion, Region: "us-east-1", Classes: []profileClass{item}}
	body, err := json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	body = append(body, '\n')
	if err = os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeSecurity(t *testing.T, path, input, decision, mime string, includeTools bool) {
	t.Helper()
	payload, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(payload)
	receipt := map[string]any{
		"schema": securityReceiptSchema, "decision": decision, "reason": "ALL_SELECTED_GATES_PASSED", "started_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano), "completed_at": time.Now().UTC().Format(time.RFC3339Nano),
		"approval_id": "security-approval-test", "policy_sha256": testHex, "input": map[string]any{"sha256": hex.EncodeToString(sum[:]), "bytes": len(payload)}, "business_storage_authorized": false,
		"security_claim": "test fixture matching the secure gate contract", "clamav": map[string]any{"exit_code": 0, "output_sha256": testHex},
		"content_type": map[string]any{"label": "pdf", "mime_type": mime, "score": 1.0}, "yara_x": map[string]any{"compiled_rules_sha256": testHex, "matching_rules": []string{}},
	}
	if includeTools {
		receipt["tools"] = map[string]any{"clamscan": map[string]any{"sha256": testHex, "version": "ClamAV 1.5.4/test"}, "magika": map[string]any{"sha256": testHex, "version": "magika 1.1.0"}, "yara_x": map[string]any{"sha256": testHex, "version": "1.20.0"}}
	}
	body, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	body = append(body, '\n')
	if err = os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func request(state fixtureState, output string) Request {
	return Request{Region: "us-east-1", DocumentClass: "supplier-invoice", InputPath: state.input, OutputDirectory: filepath.Join(state.root, output), ProfilePath: state.profile, SecurityReceiptPath: state.security}
}

func TestAnalyzeExpensePreservesOfficialOutputAndEvidenceChain(t *testing.T) {
	state := newFixture(t, expenseClass())
	client := &fakeClient{expense: &textract.AnalyzeExpenseOutput{DocumentMetadata: &types.DocumentMetadata{Pages: aws.Int32(1)}, ExpenseDocuments: []types.ExpenseDocument{{ExpenseIndex: aws.Int32(1), SummaryFields: []types.ExpenseField{{Type: &types.ExpenseType{Text: aws.String("INVOICE_RECEIPT_ID")}, ValueDetection: &types.ExpenseDetection{Text: aws.String("INV-1"), Confidence: aws.Float32(99)}}}}}}}
	receipt, err := AnalyzeToEvidence(context.Background(), client, request(state, "expense"))
	if err != nil {
		t.Fatal(err)
	}
	if receipt.SDKVersion != "v1.45.0" || receipt.DocumentClass != "supplier-invoice" || receipt.ResponseSchemaVersion != "supplier-invoice/v1" || receipt.AutomaticStorageAuthorized || receipt.PageCount != 1 || receipt.ResultCount != 1 {
		t.Fatalf("unexpected receipt: %+v", receipt)
	}
	if !isSHA256(receipt.DocumentProfileSHA256) || !isSHA256(receipt.SecurityReceiptSHA256) || string(client.expenseInput.Document.Bytes) != "%PDF-1.7 fixture" {
		t.Fatal("official request or evidence hashes not preserved")
	}
	var raw map[string]any
	body, err := os.ReadFile(filepath.Join(state.root, "expense", "provider-response.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(body, &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["ExpenseDocuments"]; !ok {
		t.Fatalf("complete modeled output missing: %s", body)
	}
}

func TestAnalyzeDocumentUsesOnlyProfileQueriesAndAdapter(t *testing.T) {
	item := profileClass{ID: "purchase-order", Decision: "REQUIRED", Operation: "AnalyzeDocument", Features: []string{"TABLES", "QUERIES"}, Queries: []Query{{Text: "What is the purchase order number?", Alias: "purchase_order_number", Pages: []string{"1"}}}, AdapterID: "adapter-1", AdapterVersion: "v7", ResponseSchemaVersion: "purchase-order/v3"}
	state := newFixture(t, item)
	client := &fakeClient{document: &textract.AnalyzeDocumentOutput{AnalyzeDocumentModelVersion: aws.String("1.0"), DocumentMetadata: &types.DocumentMetadata{Pages: aws.Int32(1)}, Blocks: []types.Block{{BlockType: types.BlockTypeQueryResult, Text: aws.String("PO-1"), Confidence: aws.Float32(98)}}}}
	req := request(state, "document")
	req.DocumentClass = "purchase-order"
	receipt, err := AnalyzeToEvidence(context.Background(), client, req)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.AdapterVersion != "v7" || receipt.ModelVersion != "1.0" || client.documentInput.QueriesConfig == nil || len(client.documentInput.QueriesConfig.Queries) != 1 || client.documentInput.AdaptersConfig == nil || len(client.documentInput.AdaptersConfig.Adapters) != 1 {
		t.Fatalf("profile selection not preserved: %+v", receipt)
	}
}

func TestSecurityReceiptMustCompletelyMatchAndAdmit(t *testing.T) {
	state := newFixture(t, expenseClass())
	client := &fakeClient{expense: &textract.AnalyzeExpenseOutput{DocumentMetadata: &types.DocumentMetadata{Pages: aws.Int32(1)}, ExpenseDocuments: []types.ExpenseDocument{{}}}}
	cases := []struct {
		name, decision, mime string
		tools                bool
		want                 string
	}{{"rejected", "REJECTED", "application/pdf", true, "did not admit"}, {"wrong-mime", "ADMITTED", "image/png", true, "content type does not match"}, {"missing-tools", "ADMITTED", "application/pdf", false, "clamscan evidence is invalid"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			writeSecurity(t, state.security, state.input, tc.decision, tc.mime, tc.tools)
			_, err := AnalyzeToEvidence(context.Background(), client, request(state, tc.name))
			if err == nil || !stringsContains(err.Error(), tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
	if client.calls != 0 {
		t.Fatalf("provider called for rejected security evidence: %d", client.calls)
	}
}

func TestProfileBlocksUnapprovedDuplicateAndRegionMismatch(t *testing.T) {
	state := newFixture(t, expenseClass())
	client := &fakeClient{}
	item := expenseClass()
	item.Decision = "CONFIGURATION_REQUIRED"
	writeProfile(t, state.profile, item)
	if _, err := AnalyzeToEvidence(context.Background(), client, request(state, "blocked")); err == nil || !stringsContains(err.Error(), "not REQUIRED") {
		t.Fatalf("blocked class accepted: %v", err)
	}
	profile := documentProfile{Schema: profileSchema, Provider: "Amazon Textract", SDK: sdkModule + "@" + sdkVersion, Region: "us-east-1", Classes: []profileClass{expenseClass(), expenseClass()}}
	body, _ := json.Marshal(profile)
	_ = os.WriteFile(state.profile, append(body, '\n'), 0o600)
	if _, err := AnalyzeToEvidence(context.Background(), client, request(state, "duplicate")); err == nil || !stringsContains(err.Error(), "exactly once") {
		t.Fatalf("duplicate class accepted: %v", err)
	}
	writeProfile(t, state.profile, expenseClass())
	req := request(state, "region")
	req.Region = "eu-west-1"
	if _, err := AnalyzeToEvidence(context.Background(), client, req); err == nil || !stringsContains(err.Error(), "does not match") {
		t.Fatalf("region mismatch accepted: %v", err)
	}
	if client.calls != 0 {
		t.Fatal("provider called before profile admission")
	}
}

func TestTemplateEnablesNoDocumentClass(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "document-profile.template.json"))
	if err != nil {
		t.Fatal(err)
	}
	var profile documentProfile
	if err = json.Unmarshal(body, &profile); err != nil {
		t.Fatal(err)
	}
	if profile.Region != "CONFIGURATION_REQUIRED" {
		t.Fatal("template region must require configuration")
	}
	for _, item := range profile.Classes {
		if item.Decision == "REQUIRED" || item.AutomaticStorage {
			t.Fatalf("template class enabled: %+v", item)
		}
	}
}

func TestRequestSurfaceCannotOverrideProfileOperation(t *testing.T) {
	typeOf := reflect.TypeOf(Request{})
	for _, name := range []string{"Mode", "Features", "Queries", "AdapterID", "AdapterVersion"} {
		if _, ok := typeOf.FieldByName(name); ok {
			t.Fatalf("unsafe request override field exists: %s", name)
		}
	}
}

func TestProviderFailureOrUnprovedPageLeavesNoOutput(t *testing.T) {
	state := newFixture(t, expenseClass())
	req := request(state, "failure")
	client := &fakeClient{err: errors.New("provider unavailable")}
	if _, err := AnalyzeToEvidence(context.Background(), client, req); err == nil {
		t.Fatal("provider error accepted")
	}
	if _, err := os.Stat(req.OutputDirectory); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("failed output committed")
	}
	client = &fakeClient{expense: &textract.AnalyzeExpenseOutput{ExpenseDocuments: []types.ExpenseDocument{{}}}}
	if _, err := AnalyzeToEvidence(context.Background(), client, req); err == nil || !stringsContains(err.Error(), "exactly one page") {
		t.Fatalf("unproved page accepted: %v", err)
	}
	if _, err := os.Stat(req.OutputDirectory); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("unproved output committed")
	}
}

func stringsContains(value, fragment string) bool {
	return len(fragment) == 0 || len(value) >= len(fragment) && func() bool {
		for i := 0; i+len(fragment) <= len(value); i++ {
			if value[i:i+len(fragment)] == fragment {
				return true
			}
		}
		return false
	}()
}
````

### FILE: `aws_textract_runtime/cmd/analyze/main.go`
```yaml
block_id: "GO-AWS-TEXTRACT-RUNTIME:cmd:v2"
operation: CREATE
provenance: AUTHORED
source: "local CLI over AWS default credential chain and official client"
license: "LicenseRef-Workspace-Owner"
sha256: "d4b0d15195726b22ec221a50a3f7ebe01263536ef2f6d3deb4b3538ce4d65108"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"example.com/elite/aws-textract-document-runtime/textractruntime"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

func main() {
	documentClass := flag.String("document-class", "", "class ID declared REQUIRED in the approved profile")
	profile := flag.String("profile", "", "approved document profile JSON")
	securityReceipt := flag.String("security-receipt", "", "ADMITTED secure local-file receipt JSON")
	input := flag.String("input", "", "local PDF/image already admitted by the security gate")
	output := flag.String("output", "", "new evidence directory")
	maxBytes := flag.Int64("max-bytes", textractruntime.MaxSyncBytes, "approved synchronous byte limit, at most 10 MiB")
	flag.Parse()
	region, err := textractruntime.ApprovedRegion(*profile, *documentClass)
	fatalIf(err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	fatalIf(err)
	client := textract.NewFromConfig(cfg)
	receipt, err := textractruntime.AnalyzeToEvidence(ctx, client, textractruntime.Request{Region: region, DocumentClass: *documentClass, InputPath: *input, OutputDirectory: *output, ProfilePath: *profile, SecurityReceiptPath: *securityReceipt, MaxBytes: *maxBytes})
	fatalIf(err)
	body, _ := json.Marshal(map[string]any{"status": "AWS_TEXTRACT_ANALYSIS_PASS", "receipt": receipt})
	fmt.Println(string(body))
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
````

### FILE: `aws_textract_runtime/README.md`
```yaml
block_id: "GO-AWS-TEXTRACT-RUNTIME:readme:v2"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0bdceffaaa6dce7f41db59771bd49f6b9741023a713a72b831cd53b56f9cf15b"
variables: []
secrets_allowed: false
```
````markdown
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
````

## 6. Configuration surface

| Variable/argument | Tipo | Default | Validación | Secreto | Mutabilidad/efecto |
|---|---|---|---|---|---|
| `document-class` | ID | none | existe exactamente una vez y decisión `REQUIRED` | no | sólo selecciona configuración aprobada |
| `profile` | JSON local aprobado | none | schema/provider/SDK/región/clase/operación/queries/adapter/schema de respuesta/storage deny | queries pueden ser sensibles | cambio exige evaluación y rollback |
| `security-receipt` | JSON local | none | `ADMITTED`, approval/policy, herramientas, resultados y SHA/bytes/MIME coincidentes | puede contener metadata del archivo | nuevo por input |
| default credential chain | reference | SDK | IAM del entorno | sí | nunca se serializa |
| `input/output` | paths | none | input regular/no symlink/admitido; output nuevo | posible PII | por ejecución |
| `max-bytes` | entero | 10 MiB | 1..10 MiB | no | nunca amplía cuota AWS |

No existe configuración para auto-storage. Cualquier combinación inválida falla cerrada.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| `github.com/aws/aws-sdk-go-v2/config` | v1.32.38 | credential/config chain | Apache-2.0 | runtime | AWS GitHub/Go module |
| `github.com/aws/aws-sdk-go-v2/service/textract` | v1.45.0 | client/types API | Apache-2.0 | runtime | AWS GitHub/Go module |
| `github.com/aws/aws-sdk-go-v2` | v1.44.0 | core SDK resuelto | Apache-2.0 | runtime | AWS GitHub/Go module |
| `github.com/aws/smithy-go` | v1.28.1 | transporte/modelado resuelto | Apache-2.0 | runtime | AWS GitHub/Go module |
| Go | 1.26.7 | build/tests | BSD-3-Clause | build/runtime | go.dev |

`go.sum` fija el grafo transitivo resuelto. El source archive AWS exacto `a30468cff35d6e385287a0cbff0ac11aa7202529` (132.399.807 bytes, SHA-256 `97d5ccc9cda85d8cd1566ec5a5a0812c3f033d419d0eb6561a83125e7592fecf`) y LICENSE/NOTICE quedan en el lock upstream; el tag/commit son unsigned y esa condición permanece abierta. Un proyecto mantiene notices/SBOM de módulos adquiridos.

## 8. Apply order

En workspace vacío: componer adquisición + seguridad local + runtime + evaluación estricta; adquirir/verificar source/dependencies oficiales; ejecutar verify/tests/vet/build offline; completar perfil; producir receipt de seguridad para cada archivo; aportar IAM/costo/corpus y llamar sandbox. En existente: compositor detecta colisiones; integrar como adapter aislado, no reemplazar persistencia/mapping sin migración y contract tests.

Rollback: retirar tráfico, restaurar adapter/version anterior o deshabilitar la ruta, conservar evidencia, revertir mapping y ejecutar corpus de regresión. Nunca borrar originales ni promover outputs parciales.

## 9. Verification

```text
go mod verify
GOPROXY=off go test ./...
go vet ./...
go build ./cmd/analyze
go list -m github.com/aws/aws-sdk-go-v2/config github.com/aws/aws-sdk-go-v2/service/textract
```

Éxito esperado: graph verificado, 7 tests + 3 subtests PASS online/offline, vet/build PASS, módulos exactos `v1.32.38`/`v1.45.0`, perfil sin overrides por caller, receipt de seguridad obligatorio, bytes oficiales, Queries/Adapters sólo desde perfil y respuesta/receipt V2 atómicos. OSV batch 2026-08-27: 15 módulos, cero findings. Antes de producción faltan AWS sandbox real, IAM/costos/cuotas, documentos representativos, métricas por campo/clase, adapter train/test, revisión, privacidad/residencia, concurrency/throttle/load y rollback.

## 10. Reconstruction evidence

- entorno limpio: fuente temporal nueva; Go oficial 1.26.7 windows/amd64;
- materialización: 7/7 archivos con SHA-256 verificado;
- source oficial: tag `service/textract/v1.45.0`, commit unsigned `a30468c...`; module ZIP 220.255 bytes/SHA-256 `c82c76b648012d198bc80c0aa303a1f6f4581af4bcb1234a3d5746f782df1cdc`; 147/147 archivos coinciden con archive; suite AWS, verify y vet PASS;
- módulos: AWS config v1.32.38 y Textract v1.45.0; `go mod verify` PASS; 15 módulos/0 OSV;
- tests locales: 7 tests + 3 subtests `GOPROXY=off go test ./...`, vet y build PASS;
- llamadas AWS: ninguna; no existían cuenta/IAM/corpus autorizados;
- divergencias: glue/tests `AUTHORED`; SDK no copiado, adquirido como módulos Apache-2.0 exactos;
- fecha/revisor: 2026-08-27 / Codex; expediente `GO_AWS_TEXTRACT_DOCUMENT_RUNTIME_2026-08-27_V2.md`.
