# Go Mercado Libre Marketplace Adapter

## 1. Metadata

```yaml
pack_id: "GO-MERCADOLIBRE-MARKETPLACE-ADAPTER"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa una frontera Go sin dependencias externas para lectura de ítems/órdenes/preguntas v4, reconciliación de preguntas sin responder, validación previa de publicaciones y normalización durable de notificaciones conforme a contratos HTTP oficiales vigentes de Mercado Libre, sin reutilizar sus SDKs archivados."
stacks: ["Go 1.26.7", "Mercado Libre REST API", "OAuth bearer token"]
compatible_with: ["GO-PROVIDER-INTEGRATION-CORE 0.1.x", "GO-ENTERPRISE-BACKEND 0.4.x"]
incompatible_with: ["SDK oficial Mercado Libre archivado", "access token en query/log/Markdown", "publicación o mutación automática", "notificaciones consultadas como URLs arbitrarias"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://developers.mercadolibre.com.ar/es_ar/autenticacion-y-autorizacion/autenticacion-y-autorizacion", "https://developers.mercadolibre.com.ar/validador-de-publicaciones", "https://developers.mercadolibre.com.ar/es_ar/productos-recibe-notificaciones", "https://developers.mercadolibre.com.ar/en_us/categories-and-attributes/manage-questions-and-answers", "https://github.com/mercadolibre/golang-sdk/tree/cd44c2846d1f7e0ff69315fbfaface6cb5a866b2"]
verified_at: "2026-09-05"
```

## 2. Applicability

Use cuando el proyecto seleccione Mercado Libre y necesite leer un ítem/orden/pregunta, reconciliar preguntas sin responder, validar un payload antes de publicar o convertir callbacks en trabajos idempotentes de consulta. Requiere aplicación y owner legal, seller grant, token seguro, scopes, IDs, políticas, límites y reconciliación probados. Rechácese como SDK oficial: Mercado Libre archivó sus SDKs y declara que no son funcionales/mantenidos. El código es `AUTHORED` contra documentación oficial. No cubre publicación/edición, stock/precio, OAuth refresh, respuestas a preguntas, claims, shipments write ni Mercado Pago.

## 3. Architecture contract

El token se resuelve únicamente desde la variable nombrada y viaja en `Authorization: Bearer`; el base URL está fijado en código a `https://api.mercadolibre.com`. Sólo se admiten `GET_ITEM`, `GET_ORDER`, `GET_QUESTION`, `SEARCH_UNANSWERED_QUESTIONS` y `VALIDATE_ITEM`; esta última exige scope/escritura, reconocimiento de que no hay preproducción y aprobación explícita. Las preguntas usan obligatoriamente `api_version=4`; la reconciliación fija seller y estado `UNANSWERED`. No existe operación de publicación ni de respuesta. Respuestas quedan atómicas y recibos sólo contienen hashes. Las notificaciones exigen control de origen en el edge, application/seller/topic/resource allowlisted, producen un `FETCH_JOB` durable y nunca consultan una URL absoluta provista por el callback. Duplicados convergen por hash del body. Rollback: detener wiring/revocar token y conservar inbox/receipts; cambios oficiales exigen reauditoría.

## 4. Exact file manifest

```text
CREATE mercadolibre_marketplace/go.mod
CREATE mercadolibre_marketplace/authority.lock.json
CREATE mercadolibre_marketplace/provider-profile.template.json
CREATE mercadolibre_marketplace/request.template.json
CREATE mercadolibre_marketplace/marketplace.go
CREATE mercadolibre_marketplace/notifications.go
CREATE mercadolibre_marketplace/marketplace_test.go
CREATE mercadolibre_marketplace/notifications_test.go
CREATE mercadolibre_marketplace/cmd/meli-marketplace/main.go
CREATE mercadolibre_marketplace/README.md
```

## 5. Materialization blocks

### FILE: `mercadolibre_marketplace/go.mod`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:gomod:v1"
operation: CREATE
provenance: AUTHORED
source: "dependency-free Go module"
license: "LicenseRef-Workspace-Owner"
sha256: "c942e5ce737c8c1e783b3f908bf61aaabc73470bbb86eeedaf5ceee7e38114fa"
variables: []
secrets_allowed: false
```
````go
module elite.local/mercadolibremarketplace

go 1.26.0
````

### FILE: `mercadolibre_marketplace/authority.lock.json`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:authority:v1"
operation: CREATE
provenance: AUTHORED
source: "official Mercado Libre documentation authority record and archived SDK rejection"
license: "LicenseRef-Workspace-Owner"
sha256: "a8a045121af6bd7cf53a719d661b9ba14a6e54b7f60be66682590577cc908bb8"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-http-authority-lock/v1",
  "provider": "Mercado Libre",
  "api_base_url": "https://api.mercadolibre.com",
  "official_sdk_status": "REJECTED_ARCHIVED_NOT_FUNCTIONAL",
  "archived_go_sdk": {
    "repository": "mercadolibre/golang-sdk",
    "commit": "cd44c2846d1f7e0ff69315fbfaface6cb5a866b2",
    "archived": true,
    "provider_notice": "not functional and no longer maintained"
  },
  "official_contracts": [
    {
      "claim": "OAuth bearer token in Authorization header",
      "url": "https://developers.mercadolibre.com.ar/es_ar/autenticacion-y-autorizacion/autenticacion-y-autorizacion",
      "last_updated": "2026-07-15"
    },
    {
      "claim": "POST /items/validate; 204 means valid; validation does not publish",
      "url": "https://developers.mercadolibre.com.ar/validador-de-publicaciones",
      "last_updated": "2025-12-30"
    },
    {
      "claim": "GET /orders/{id}",
      "url": "https://developers.mercadolibre.com.ar/es_ar/saldo-de-la-cuenta/gestion-packs",
      "last_updated": "2025-12-29"
    },
    {
      "claim": "notifications must be acknowledged quickly, queued, and fetched by resource; duplicates/retries exist",
      "url": "https://developers.mercadolibre.com.ar/es_ar/productos-recibe-notificaciones",
      "last_updated": "2026-07-30"
    },
    {
      "claim": "application scopes separate GET read from POST/PUT/DELETE write; Mercado Libre and Mercado Pago applications separate from 2026-08-30",
      "url": "https://developers.mercadolibre.com.ar/es_ar/metricas/crea-una-aplicacion-en-mercado-libre-es",
      "last_updated": "2026-08-06"
    },
    {
      "claim": "questions topic; GET /questions/{id} and /questions/search use api_version=4; BANNED text may be empty; answers are separate writes limited to 2000 characters",
      "url": "https://developers.mercadolibre.com.ar/en_us/categories-and-attributes/manage-questions-and-answers",
      "last_updated": "2026-01-15"
    }
  ],
  "verified_at": "2026-09-05"
}
````

### FILE: `mercadolibre_marketplace/provider-profile.template.json`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "fail-closed provider/access/policy profile"
license: "LicenseRef-Workspace-Owner"
sha256: "5d4ce698b2a3300d76a749f39af772b084d31de181f4bb95888bb065b89bd62e"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-mercadolibre-marketplace-profile/v1",
  "provider": "Mercado Libre",
  "decision": "BLOCKED",
  "official_documentation_reviewed": false,
  "archived_sdk_rejection_acknowledged": false,
  "application_and_legal_owner_proven": false,
  "seller_authorization_proven": false,
  "oauth_token_storage_proven": false,
  "read_scope_proven": false,
  "write_scope_proven": false,
  "no_preproduction_environment_acknowledged": false,
  "item_and_category_policy_approved": false,
  "notification_origin_control_proven": false,
  "notification_topics_approved": false,
  "quota_and_cost_approved": false,
  "reconciliation_approved": false,
  "site_id": "MLA",
  "seller_id": "",
  "application_id": "",
  "access_token_environment_variable": "MERCADOLIBRE_ACCESS_TOKEN",
  "approved_notification_topics": [],
  "validation_call_approved": false,
  "automatic_business_write": false
}
````

### FILE: `mercadolibre_marketplace/request.template.json`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:request:v1"
operation: CREATE
provenance: AUTHORED
source: "minimal allowlisted operation request"
license: "LicenseRef-Workspace-Owner"
sha256: "8b8b4003ff34b19fd9ec996bedfb07f55460a158feb5770d8edcc98434270ce3"
variables: []
secrets_allowed: false
```
````json
{
  "operation": "GET_ITEM",
  "resource_id": "MLA0000000000"
}
````

### FILE: `mercadolibre_marketplace/marketplace.go`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:adapter:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter governed by official Mercado Libre HTTP contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "c543958d1306c87c34fc502ec014574a975e55280e18aa6eba8ea029aa7d9f66"
variables: []
secrets_allowed: false
```
````go
package mercadolibremarketplace

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const APIBaseURL = "https://api.mercadolibre.com"

type Profile struct {
	Schema                                 string   `json:"schema"`
	Provider                               string   `json:"provider"`
	Decision                               string   `json:"decision"`
	OfficialDocumentationReviewed          bool     `json:"official_documentation_reviewed"`
	ArchivedSDKRejectionAcknowledged       bool     `json:"archived_sdk_rejection_acknowledged"`
	ApplicationAndLegalOwnerProven         bool     `json:"application_and_legal_owner_proven"`
	SellerAuthorizationProven              bool     `json:"seller_authorization_proven"`
	OAuthTokenStorageProven                bool     `json:"oauth_token_storage_proven"`
	ReadScopeProven                        bool     `json:"read_scope_proven"`
	WriteScopeProven                       bool     `json:"write_scope_proven"`
	NoPreproductionEnvironmentAcknowledged bool     `json:"no_preproduction_environment_acknowledged"`
	ItemAndCategoryPolicyApproved          bool     `json:"item_and_category_policy_approved"`
	NotificationOriginControlProven        bool     `json:"notification_origin_control_proven"`
	NotificationTopicsApproved             bool     `json:"notification_topics_approved"`
	QuotaAndCostApproved                   bool     `json:"quota_and_cost_approved"`
	ReconciliationApproved                 bool     `json:"reconciliation_approved"`
	SiteID                                 string   `json:"site_id"`
	SellerID                               string   `json:"seller_id"`
	ApplicationID                          string   `json:"application_id"`
	AccessTokenEnvironmentVariable         string   `json:"access_token_environment_variable"`
	ApprovedNotificationTopics             []string `json:"approved_notification_topics"`
	ValidationCallApproved                 bool     `json:"validation_call_approved"`
	AutomaticBusinessWrite                 bool     `json:"automatic_business_write"`
}

type Request struct {
	Operation  string          `json:"operation"`
	ResourceID string          `json:"resource_id"`
	Item       json.RawMessage `json:"item,omitempty"`
}

type Receipt struct {
	Schema                 string `json:"schema"`
	CreatedAt              string `json:"created_at"`
	Operation              string `json:"operation"`
	HTTPStatus             int    `json:"http_status"`
	Outcome                string `json:"outcome"`
	SellerIDSHA256         string `json:"seller_id_sha256"`
	ApplicationIDSHA256    string `json:"application_id_sha256"`
	ResourceIDSHA256       string `json:"resource_id_sha256"`
	RequestSHA256          string `json:"request_sha256"`
	ResponseSHA256         string `json:"response_sha256"`
	AutomaticBusinessWrite bool   `json:"automatic_business_write"`
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

var (
	sitePattern   = regexp.MustCompile(`^ML[A-Z]$`)
	digitsPattern = regexp.MustCompile(`^[0-9]{3,20}$`)
	itemPattern   = regexp.MustCompile(`^ML[A-Z][0-9]{6,20}$`)
	envPattern    = regexp.MustCompile(`^[A-Z][A-Z0-9_]{2,80}$`)
)

func LoadJSON(path string, target any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(io.LimitReader(file, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("JSON must contain exactly one value")
	}
	return nil
}

func Validate(profile Profile, request Request, token string) error {
	if profile.Schema != "elite-mercadolibre-marketplace-profile/v1" || profile.Provider != "Mercado Libre" || profile.Decision != "PROVEN" {
		return errors.New("provider profile is not PROVEN")
	}
	if !profile.OfficialDocumentationReviewed || !profile.ArchivedSDKRejectionAcknowledged || !profile.ApplicationAndLegalOwnerProven || !profile.SellerAuthorizationProven || !profile.OAuthTokenStorageProven || !profile.ReadScopeProven || !profile.QuotaAndCostApproved || !profile.ReconciliationApproved {
		return errors.New("required provider proof is missing")
	}
	if profile.AutomaticBusinessWrite {
		return errors.New("automatic business write is forbidden")
	}
	if !sitePattern.MatchString(profile.SiteID) || !digitsPattern.MatchString(profile.SellerID) || !digitsPattern.MatchString(profile.ApplicationID) || !envPattern.MatchString(profile.AccessTokenEnvironmentVariable) {
		return errors.New("provider identifiers are invalid")
	}
	if len(token) < 20 || len(token) > 4096 || strings.ContainsAny(token, "\r\n") {
		return errors.New("access token is invalid")
	}
	switch request.Operation {
	case "GET_ITEM":
		if !itemPattern.MatchString(request.ResourceID) || len(request.Item) != 0 {
			return errors.New("GET_ITEM request is invalid")
		}
	case "GET_ORDER":
		if !digitsPattern.MatchString(request.ResourceID) || len(request.Item) != 0 {
			return errors.New("GET_ORDER request is invalid")
		}
	case "GET_QUESTION":
		if !digitsPattern.MatchString(request.ResourceID) || len(request.Item) != 0 {
			return errors.New("GET_QUESTION request is invalid")
		}
	case "SEARCH_UNANSWERED_QUESTIONS":
		if request.ResourceID != "" || len(request.Item) != 0 {
			return errors.New("SEARCH_UNANSWERED_QUESTIONS request is invalid")
		}
	case "VALIDATE_ITEM":
		if request.ResourceID != "" || !profile.WriteScopeProven || !profile.NoPreproductionEnvironmentAcknowledged || !profile.ItemAndCategoryPolicyApproved || !profile.ValidationCallApproved {
			return errors.New("item validation is not approved")
		}
		if len(request.Item) < 2 || len(request.Item) > 1<<20 || request.Item[0] != '{' || !json.Valid(request.Item) {
			return errors.New("item JSON is invalid")
		}
	default:
		return errors.New("operation is not allowlisted")
	}
	return nil
}

func Execute(ctx context.Context, doer Doer, profile Profile, request Request, token, outputDirectory string) (Receipt, error) {
	if doer == nil {
		return Receipt{}, errors.New("HTTP client is required")
	}
	if err := Validate(profile, request, token); err != nil {
		return Receipt{}, err
	}
	method, path, body := http.MethodGet, "", []byte(nil)
	switch request.Operation {
	case "GET_ITEM":
		path = "/items/" + request.ResourceID
	case "GET_ORDER":
		path = "/orders/" + request.ResourceID
	case "GET_QUESTION":
		path = "/questions/" + request.ResourceID + "?api_version=4"
	case "SEARCH_UNANSWERED_QUESTIONS":
		path = "/questions/search?seller_id=" + profile.SellerID + "&status=UNANSWERED&api_version=4"
	case "VALIDATE_ITEM":
		method, path, body = http.MethodPost, "/items/validate", []byte(request.Item)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, method, APIBaseURL+path, bytes.NewReader(body))
	if err != nil {
		return Receipt{}, err
	}
	httpRequest.Header.Set("Authorization", "Bearer "+token)
	httpRequest.Header.Set("Accept", "application/json")
	if method == http.MethodPost {
		httpRequest.Header.Set("Content-Type", "application/json")
	}
	response, err := doer.Do(httpRequest)
	if err != nil {
		return Receipt{}, fmt.Errorf("Mercado Libre request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, (4<<20)+1))
	if err != nil {
		return Receipt{}, err
	}
	if len(responseBody) > 4<<20 {
		return Receipt{}, errors.New("provider response exceeds 4 MiB")
	}
	outcome := "ACCEPTED"
	if request.Operation == "VALIDATE_ITEM" && response.StatusCode == http.StatusBadRequest {
		outcome = "REJECTED_BY_VALIDATOR"
	} else if request.Operation == "VALIDATE_ITEM" && response.StatusCode != http.StatusNoContent {
		return Receipt{}, fmt.Errorf("unexpected validation status %d", response.StatusCode)
	} else if request.Operation != "VALIDATE_ITEM" && response.StatusCode != http.StatusOK {
		return Receipt{}, fmt.Errorf("unexpected read status %d", response.StatusCode)
	}
	abs, err := filepath.Abs(outputDirectory)
	if err != nil {
		return Receipt{}, err
	}
	if _, err := os.Stat(abs); !os.IsNotExist(err) {
		return Receipt{}, errors.New("output directory must not exist")
	}
	requestBytes, _ := json.Marshal(request)
	receipt := Receipt{Schema: "elite-mercadolibre-marketplace-receipt/v1", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano), Operation: request.Operation, HTTPStatus: response.StatusCode, Outcome: outcome, SellerIDSHA256: hash([]byte(profile.SellerID)), ApplicationIDSHA256: hash([]byte(profile.ApplicationID)), ResourceIDSHA256: hash([]byte(request.ResourceID)), RequestSHA256: hash(requestBytes), ResponseSHA256: hash(responseBody), AutomaticBusinessWrite: false}
	stage := abs + ".stage-" + fmt.Sprint(time.Now().UnixNano())
	if err := os.Mkdir(stage, 0o700); err != nil {
		return Receipt{}, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(stage)
		}
	}()
	if err := os.WriteFile(filepath.Join(stage, "PROVIDER_RESPONSE.json"), append(responseBody, '\n'), 0o600); err != nil {
		return Receipt{}, err
	}
	receiptBytes, _ := json.MarshalIndent(receipt, "", "  ")
	if err := os.WriteFile(filepath.Join(stage, "MARKETPLACE_RECEIPT.json"), append(receiptBytes, '\n'), 0o600); err != nil {
		return Receipt{}, err
	}
	if err := os.Rename(stage, abs); err != nil {
		return Receipt{}, err
	}
	committed = true
	return receipt, nil
}

func hash(value []byte) string { sum := sha256.Sum256(value); return hex.EncodeToString(sum[:]) }
````

### FILE: `mercadolibre_marketplace/notifications.go`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:notifications:v1"
operation: CREATE
provenance: AUTHORED
source: "notification queue/fetch contract based on official retry guidance"
license: "LicenseRef-Workspace-Owner"
sha256: "6fcaba883667b5338eca77d2d1a37c4c6d4f13b50a8798aa644774c8519c5121"
variables: []
secrets_allowed: false
```
````go
package mercadolibremarketplace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Notification struct {
	ID            string `json:"_id,omitempty"`
	Resource      string `json:"resource"`
	UserID        int64  `json:"user_id"`
	Topic         string `json:"topic"`
	ApplicationID int64  `json:"application_id"`
	Attempts      int    `json:"attempts"`
	Sent          string `json:"sent"`
	Received      string `json:"received"`
}

type FetchJob struct {
	Schema           string `json:"schema"`
	Topic            string `json:"topic"`
	Resource         string `json:"resource"`
	DeduplicationKey string `json:"deduplication_key"`
	ReceivedAt       string `json:"received_at"`
}

var resourcePattern = regexp.MustCompile(`^/(items/ML[A-Z][0-9]{6,20}|orders/[0-9]{3,20}|shipments/[0-9]{3,20}|questions/[0-9]{3,20})$`)

func NormalizeNotification(profile Profile, raw []byte) (FetchJob, error) {
	if profile.Decision != "PROVEN" || !profile.NotificationOriginControlProven || !profile.NotificationTopicsApproved || !profile.ReconciliationApproved {
		return FetchJob{}, errors.New("notification profile is not proven")
	}
	if len(raw) == 0 || len(raw) > 1<<20 {
		return FetchJob{}, errors.New("notification size is invalid")
	}
	var notification Notification
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&notification); err != nil {
		return FetchJob{}, err
	}
	if !resourcePattern.MatchString(notification.Resource) || notification.Attempts < 1 || notification.Attempts > 100 {
		return FetchJob{}, errors.New("notification envelope is invalid")
	}
	if notification.UserID <= 0 || notification.ApplicationID <= 0 || notification.UserID != parsePositive(profile.SellerID) || notification.ApplicationID != parsePositive(profile.ApplicationID) {
		return FetchJob{}, errors.New("notification identity mismatch")
	}
	allowed := false
	for _, topic := range profile.ApprovedNotificationTopics {
		if topic == notification.Topic {
			allowed = true
		}
	}
	if !allowed || !((notification.Topic == "items" && strings.HasPrefix(notification.Resource, "/items/")) || (notification.Topic == "orders_v2" && strings.HasPrefix(notification.Resource, "/orders/")) || (notification.Topic == "shipments" && strings.HasPrefix(notification.Resource, "/shipments/")) || (notification.Topic == "questions" && strings.HasPrefix(notification.Resource, "/questions/"))) {
		return FetchJob{}, errors.New("notification topic/resource is not approved")
	}
	if _, err := time.Parse(time.RFC3339Nano, notification.Sent); err != nil {
		return FetchJob{}, errors.New("notification sent timestamp is invalid")
	}
	sum := sha256.Sum256(raw)
	return FetchJob{Schema: "elite-mercadolibre-fetch-job/v1", Topic: notification.Topic, Resource: notification.Resource, DeduplicationKey: hex.EncodeToString(sum[:]), ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}, nil
}

func WriteFetchJob(job FetchJob, outputDirectory string) error {
	abs, err := filepath.Abs(outputDirectory)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); !os.IsNotExist(err) {
		return errors.New("output directory must not exist")
	}
	stage := abs + ".stage"
	if err := os.Mkdir(stage, 0o700); err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(stage)
		}
	}()
	data, _ := json.MarshalIndent(job, "", "  ")
	if err := os.WriteFile(filepath.Join(stage, "FETCH_JOB.json"), append(data, '\n'), 0o600); err != nil {
		return err
	}
	if err := os.Rename(stage, abs); err != nil {
		return err
	}
	committed = true
	return nil
}

func parsePositive(value string) int64 {
	var result int64
	for _, char := range value {
		if char < '0' || char > '9' {
			return -1
		}
		result = result*10 + int64(char-'0')
	}
	return result
}
````

### FILE: `mercadolibre_marketplace/marketplace_test.go`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local official-contract and negative tests"
license: "LicenseRef-Workspace-Owner"
sha256: "f6d8a1de31937e7846b5465731ea19bb9b89ff924d0deb7611f0a210d9e69c6e"
variables: []
secrets_allowed: false
```
````go
package mercadolibremarketplace

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeDoer struct {
	response *http.Response
	err      error
	request  *http.Request
	body     []byte
}

func (f *fakeDoer) Do(request *http.Request) (*http.Response, error) {
	f.request = request
	f.body, _ = io.ReadAll(request.Body)
	return f.response, f.err
}

func provenProfile() Profile {
	return Profile{Schema: "elite-mercadolibre-marketplace-profile/v1", Provider: "Mercado Libre", Decision: "PROVEN", OfficialDocumentationReviewed: true, ArchivedSDKRejectionAcknowledged: true, ApplicationAndLegalOwnerProven: true, SellerAuthorizationProven: true, OAuthTokenStorageProven: true, ReadScopeProven: true, WriteScopeProven: true, NoPreproductionEnvironmentAcknowledged: true, ItemAndCategoryPolicyApproved: true, NotificationOriginControlProven: true, NotificationTopicsApproved: true, QuotaAndCostApproved: true, ReconciliationApproved: true, SiteID: "MLA", SellerID: "123456789", ApplicationID: "987654321", AccessTokenEnvironmentVariable: "MERCADOLIBRE_ACCESS_TOKEN", ApprovedNotificationTopics: []string{"items", "orders_v2", "questions"}, ValidationCallApproved: true, AutomaticBusinessWrite: false}
}

const token = "APP_USR-test-token-value-1234567890"

func TestGetItemUsesExactOfficialContractAndRedactsReceipt(t *testing.T) {
	doer := &fakeDoer{response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"id":"MLA1234567890","title":"x"}`))}}
	output := filepath.Join(t.TempDir(), "evidence")
	receipt, err := Execute(context.Background(), doer, provenProfile(), Request{Operation: "GET_ITEM", ResourceID: "MLA1234567890"}, token, output)
	if err != nil {
		t.Fatal(err)
	}
	if doer.request.Method != "GET" || doer.request.URL.String() != APIBaseURL+"/items/MLA1234567890" || doer.request.Header.Get("Authorization") != "Bearer "+token {
		t.Fatal("official request contract mismatch")
	}
	serialized, _ := json.Marshal(receipt)
	if strings.Contains(string(serialized), "123456789") || strings.Contains(string(serialized), token) || receipt.AutomaticBusinessWrite {
		t.Fatal("receipt leaked identifiers or enabled write")
	}
}

func TestValidateItemUsesValidatorAndPreservesRejection(t *testing.T) {
	doer := &fakeDoer{response: &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(`{"message":"body.invalid_field_types"}`))}}
	receipt, err := Execute(context.Background(), doer, provenProfile(), Request{Operation: "VALIDATE_ITEM", Item: json.RawMessage(`{"title":"x"}`)}, token, filepath.Join(t.TempDir(), "evidence"))
	if err != nil {
		t.Fatal(err)
	}
	if doer.request.Method != "POST" || doer.request.URL.Path != "/items/validate" || receipt.Outcome != "REJECTED_BY_VALIDATOR" {
		t.Fatal("validator contract mismatch")
	}
}

func TestQuestionReadAndReconciliationUseOfficialV4Contract(t *testing.T) {
	tests := []struct {
		name     string
		request  Request
		path     string
		rawQuery string
	}{
		{name: "question by id", request: Request{Operation: "GET_QUESTION", ResourceID: "11751825075"}, path: "/questions/11751825075", rawQuery: "api_version=4"},
		{name: "unanswered seller reconciliation", request: Request{Operation: "SEARCH_UNANSWERED_QUESTIONS"}, path: "/questions/search", rawQuery: "seller_id=123456789&status=UNANSWERED&api_version=4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doer := &fakeDoer{response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"questions":[]}`))}}
			if _, err := Execute(context.Background(), doer, provenProfile(), tt.request, token, filepath.Join(t.TempDir(), "evidence")); err != nil {
				t.Fatal(err)
			}
			if doer.request.Method != http.MethodGet || doer.request.URL.Path != tt.path || doer.request.URL.RawQuery != tt.rawQuery {
				t.Fatalf("method=%s path=%s query=%s", doer.request.Method, doer.request.URL.Path, doer.request.URL.RawQuery)
			}
		})
	}
}

func TestProfileAndOperationsFailClosed(t *testing.T) {
	profile := provenProfile()
	profile.Decision = "BLOCKED"
	if err := Validate(profile, Request{Operation: "GET_ITEM", ResourceID: "MLA1234567890"}, token); err == nil {
		t.Fatal("blocked profile passed")
	}
	profile = provenProfile()
	profile.ValidationCallApproved = false
	if err := Validate(profile, Request{Operation: "VALIDATE_ITEM", Item: json.RawMessage(`{"title":"x"}`)}, token); err == nil {
		t.Fatal("unapproved validation passed")
	}
	if err := Validate(provenProfile(), Request{Operation: "PUBLISH_ITEM", Item: json.RawMessage(`{}`)}, token); err == nil {
		t.Fatal("unimplemented write passed")
	}
}

func TestProviderFailureAndUnexpectedStatusAreAtomic(t *testing.T) {
	parent := t.TempDir()
	output := filepath.Join(parent, "evidence")
	_, err := Execute(context.Background(), &fakeDoer{err: errors.New("provider unavailable")}, provenProfile(), Request{Operation: "GET_ORDER", ResourceID: "123456789"}, token, output)
	if err == nil || !strings.Contains(err.Error(), "provider unavailable") {
		t.Fatal("provider failure not retained")
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatal("partial output exists")
	}
	_, err = Execute(context.Background(), &fakeDoer{response: &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader(`{}`))}}, provenProfile(), Request{Operation: "GET_ORDER", ResourceID: "123456789"}, token, output)
	if err == nil {
		t.Fatal("unexpected status passed")
	}
}

func TestLoadJSONRejectsUnknownAndTrailing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "request.json")
	if err := os.WriteFile(path, []byte(`{"operation":"GET_ITEM","resource_id":"MLA1234567890","unknown":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var request Request
	if LoadJSON(path, &request) == nil {
		t.Fatal("unknown field passed")
	}
	if err := os.WriteFile(path, []byte(`{"operation":"GET_ITEM","resource_id":"MLA1234567890"} {}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if LoadJSON(path, &request) == nil {
		t.Fatal("trailing value passed")
	}
}
````

### FILE: `mercadolibre_marketplace/notifications_test.go`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:notificationtests:v1"
operation: CREATE
provenance: AUTHORED
source: "local notification identity/topic/resource/durability tests"
license: "LicenseRef-Workspace-Owner"
sha256: "39450c905b471857c240dceb0f381395039d71ec7797cb39352bb54fb8736204"
variables: []
secrets_allowed: false
```
````go
package mercadolibremarketplace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func notificationJSON(topic, resource string, user, application int64) []byte {
	value := Notification{ID: "provider-id", Resource: resource, UserID: user, Topic: topic, ApplicationID: application, Attempts: 1, Sent: "2026-08-26T10:00:00Z", Received: "2026-08-26T10:00:01Z"}
	result, _ := json.Marshal(value)
	return result
}

func TestNotificationProducesDurableFetchJob(t *testing.T) {
	job, err := NormalizeNotification(provenProfile(), notificationJSON("orders_v2", "/orders/123456789", 123456789, 987654321))
	if err != nil {
		t.Fatal(err)
	}
	if job.Resource != "/orders/123456789" || job.Topic != "orders_v2" || len(job.DeduplicationKey) != 64 {
		t.Fatal("invalid fetch job")
	}
	output := filepath.Join(t.TempDir(), "job")
	if err := WriteFetchJob(job, output); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(output, "FETCH_JOB.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "987654321") || strings.Contains(string(data), "provider-id") {
		t.Fatal("fetch job leaked envelope identifiers")
	}
}

func TestQuestionNotificationProducesAllowlistedFetchJob(t *testing.T) {
	job, err := NormalizeNotification(provenProfile(), notificationJSON("questions", "/questions/11751825075", 123456789, 987654321))
	if err != nil {
		t.Fatal(err)
	}
	if job.Resource != "/questions/11751825075" || job.Topic != "questions" || len(job.DeduplicationKey) != 64 {
		t.Fatalf("invalid question fetch job: %+v", job)
	}
}

func TestNotificationIdentityTopicAndResourceFailClosed(t *testing.T) {
	if _, err := NormalizeNotification(provenProfile(), notificationJSON("orders_v2", "/orders/123456789", 111111111, 987654321)); err == nil {
		t.Fatal("wrong seller passed")
	}
	if _, err := NormalizeNotification(provenProfile(), notificationJSON("unknown", "/orders/123456789", 123456789, 987654321)); err == nil {
		t.Fatal("unknown topic passed")
	}
	if _, err := NormalizeNotification(provenProfile(), notificationJSON("orders_v2", "https://evil.example/x", 123456789, 987654321)); err == nil {
		t.Fatal("untrusted resource passed")
	}
}

func TestNotificationRequiresOriginControl(t *testing.T) {
	profile := provenProfile()
	profile.NotificationOriginControlProven = false
	if _, err := NormalizeNotification(profile, notificationJSON("items", "/items/MLA1234567890", 123456789, 987654321)); err == nil {
		t.Fatal("unproven origin control passed")
	}
}
````

### FILE: `mercadolibre_marketplace/cmd/meli-marketplace/main.go`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:cli:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded CLI with environment-only token"
license: "LicenseRef-Workspace-Owner"
sha256: "1d090bcde11d76d5f96f9160e664fdf19f2f5575a3ee9a8e9c3e17caad757b12"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	marketplace "elite.local/mercadolibremarketplace"
)

func main() {
	profilePath := flag.String("profile", "", "approved provider profile JSON")
	requestPath := flag.String("request", "", "approved operation request JSON")
	output := flag.String("output", "", "new evidence directory")
	flag.Parse()
	if *profilePath == "" || *requestPath == "" || *output == "" {
		fatal("profile, request and output are required")
	}
	var profile marketplace.Profile
	var request marketplace.Request
	if err := marketplace.LoadJSON(*profilePath, &profile); err != nil {
		fatal(err.Error())
	}
	if err := marketplace.LoadJSON(*requestPath, &request); err != nil {
		fatal(err.Error())
	}
	token := os.Getenv(profile.AccessTokenEnvironmentVariable)
	client := &http.Client{Timeout: 30 * time.Second}
	if _, err := marketplace.Execute(context.Background(), client, profile, request, token, *output); err != nil {
		fatal(err.Error())
	}
	fmt.Println("MERCADOLIBRE_MARKETPLACE_OPERATION_PASS")
}

func fatal(message string) { fmt.Fprintln(os.Stderr, message); os.Exit(1) }
````

### FILE: `mercadolibre_marketplace/README.md`
```yaml
block_id: "GO-MERCADOLIBRE-MARKETPLACE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "operator guide and non-claims"
license: "LicenseRef-Workspace-Owner"
sha256: "959fe242a4daf7e57ec9e01f380ee5e07171341797cabc96025e955ca4ab7c73"
variables: []
secrets_allowed: false
```
````markdown
# Mercado Libre marketplace adapter

This authored Go adapter follows the current official Mercado Libre HTTP documentation. It does **not** use or revive Mercado Libre's archived SDK, which the provider marks not functional and unmaintained.

Implemented claims:

- authenticated `GET /items/{id}`;
- authenticated `GET /orders/{id}`;
- authenticated `GET /questions/{id}?api_version=4`;
- seller reconciliation with `GET /questions/search?seller_id=...&status=UNANSWERED&api_version=4`;
- approved `POST /items/validate`, preserving 204 acceptance or 400 validation details;
- notification envelope allowlisting for `items`, `orders_v2`, `shipments` and `questions`, durable fetch jobs and raw-body deduplication.

The template is blocked. Before use, prove the legal application owner, seller authorization, secure OAuth token storage, scopes, site/seller/application IDs, item/category rules, the lack of a preproduction environment, origin controls for notifications, topics, quotas/cost and reconciliation. The bearer token exists only in the named environment variable.

```powershell
go mod verify
go test ./...
go vet ./...
go run ./cmd/meli-marketplace -profile .\provider-profile.json -request .\request.json -output .\evidence\operation-001
```

Question responses can be passed to the `GO-OMNICHANNEL-LEAD-INGRESS` Mercado Libre v4 decoder. This pack deliberately does not publish, modify or delete listings; change price/stock; refresh OAuth tokens; answer questions; acknowledge claims; update shipments; or infer a provider response as a successful business write. Those capabilities require separate official-contract packs, an outbound fence and real-account gates.
````

## 6. Configuration surface

| Campo | Default seguro | Gate | Secreto | Efecto |
|---|---|---|---|---|
| `decision` | `BLOCKED` | sólo `PROVEN` | no | habilita frontera |
| proofs de docs/SDK rejection/app/seller/token/scopes/policies/origin/topic/quota/reconciliation | `false` | requeridos según operación | no | adopción/operación |
| `site_id`, `seller_id`, `application_id` | vacío/MLA template | patrones exactos | identificadores sensibles | routing/identity |
| `access_token_environment_variable` | nombre, sin valor | patrón de env | el valor sí | runtime/rotación |
| `validation_call_approved` | `false` | requerido para POST validate | no | llamada externa sin publicación |
| `approved_notification_topics` | vacío | items/orders_v2/shipments y resource matching | no | ingestión |
| `automatic_business_write` | `false` | debe permanecer false | no | inmutable |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go stdlib | Go 1.26.7 | HTTP/JSON/hash/atomic I/O/tests | BSD-3-Clause | build/runtime | https://go.dev/dl/ |
| Mercado Libre REST API | docs observadas 2025-12-30..2026-08-06 | endpoints/headers/notifications | términos de plataforma, no redistribuidos | externo | https://developers.mercadolibre.com.ar/devsite/api-docs |
| `mercadolibre/golang-sdk` | commit `cd44c284...` | sólo evidencia de rechazo | Apache-2.0 en archivos | no runtime | https://github.com/mercadolibre/golang-sdk |

## 8. Apply order

1. Completar gate de proyecto: proveedor, owner/app, seller grant, token store, scopes, site/IDs, términos, políticas, origen de callbacks, quota/costo y reconciliación.
2. Materializar diez archivos; resolver el token sólo en secret manager/environment del target.
3. Ejecutar `go mod verify`, `go test ./...`, `go vet ./...` y build con Go 1.26.7.
4. Probar primero GET con test user autorizado; luego `VALIDATE_ITEM` con aprobación explícita. No publicar.
5. Conectar callback sólo después del control de origen; persistir `FETCH_JOB`, responder rápido y consultar el recurso mediante worker/inbox.
6. Rollback: cortar routing, revocar token, preservar inbox/receipts y reauditar docs antes de reactivar.

## 9. Verification

- 10/10 archivos y hashes; reconstrucción byte-idéntica.
- `go mod verify`, ocho tests, `go vet` y build deben pasar.
- Requests prueban método, base/path y header Bearer exactos; validator conserva 204/400 sin publicar.
- Perfil bloqueado, operación desconocida/publicación, write no aprobado, token/IDs inválidos, status/fallo provider, JSON extra y output existente fallan cerrados.
- Notification seller/application/topic/resource/origin se valida; recursos absolutos/arbitrarios se rechazan; fetch job durable no expone envelope IDs.
- Cuenta/test user/token/endpoint real, quotas, retry-after, live callbacks, missed feeds y reconciliación siguen gates productivos.

## 10. Reconstruction evidence

`reconstruction_evidence/MERCADOLIBRE_MARKETPLACE_ADAPTER_2026-08-26_V1.md`: Go oficial 1.26.7, diez archivos, ocho tests, vet/build y reconstrucción limpia. No se contactó una cuenta ni se realizó publicación/mutación, por lo que permanece `CONDITIONED`.
