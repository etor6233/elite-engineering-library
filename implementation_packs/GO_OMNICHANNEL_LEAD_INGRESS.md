# Go Omnichannel Lead Ingress

## 1. Metadata

```yaml
pack_id: "GO-OMNICHANNEL-LEAD-INGRESS"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa una frontera provider-neutral, el webhook oficial de Google Lead Form y el decoder de respuestas oficiales Mercado Libre Questions v4: autenticación/owner previos, hash del body fuente, payload durable gobernado, identidad idempotente, rechazo durable, candidato pending_policy y outbox en la misma transacción PostgreSQL."
stacks: ["Go 1.26.7", "PostgreSQL 18.6"]
compatible_with: ["GO-ENTERPRISE-BACKEND-CORE 0.4.x", "GO-RELIABLE-ASYNC-WORKERS 0.1.x", "GO-ELECTROMOBILITY-PUBLIC-CRM-API 0.2.x"]
incompatible_with: ["MemoryStore en producción", "contacto automático sin política/consentimiento", "payload provider inventado", "ack antes del commit", "exactly-once inferido"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources:
  - "https://developers.google.com/google-ads/webhook/docs/implementation"
  - "https://developers.google.com/google-ads/api/samples/add-lead-form-asset"
  - "https://developers.mercadolibre.com.ar/en_us/categories-and-attributes/manage-questions-and-answers"
  - "https://developers.mercadolibre.com.ar/es_ar/productos-recibe-notificaciones"
  - "https://github.com/cloudevents/spec/blob/ce@v1.0.2/cloudevents/spec.md"
  - "https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html"
  - "https://docs.cloud.google.com/pubsub/docs/subscription-retry-policy"
  - "https://www.postgresql.org/docs/18/transaction-iso.html"
  - "https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html"
verified_at: "2026-09-05"
```

## 2. Applicability

Use para recibir leads de Google Ads y normalizar la respuesta autoritativa obtenida tras una notificación `questions` de Mercado Libre. Google conserva la identidad SHA-256 del body fuente y una representación durable que elimina `google_key`; Mercado Libre exige seller esperado, API v4, estado conocido y conserva el body de respuesta con PII bajo el mismo owner durable. Ningún ingreso concede permiso de contacto por inferencia. La frontera `RawEvent` sólo se extiende cuando existe contrato oficial exacto y parser separado. Meta y TikTok mantienen sus adapters, firmas, permisos y postbacks propios.

## 3. Architecture contract

- **Ownership**: `internal/leadstream` gobierna autenticación/normalización del edge; `postgres.LeadIngress` gobierna commit atómico raw+candidato+outbox. `crm` conserva el ownership de lead/contacto posterior.
- **Flujo**: provider HTTP → verificación exacta → SHA-256 del body fuente → redacción específica si existe credencial → payload durable + SHA-256 + perfil de redacción → candidato o rechazo → transacción PostgreSQL → ACK. Google elimina `google_key`; Mercado Libre sólo entra desde `GET /questions/{id}?api_version=4` ejecutado con seller autorizado y rechaza seller divergente antes de persistir. Duplicado exacto es replay; misma identidad con distinto hash fuente falla cerrado.
- **Privacidad**: `google_key` no entra en PostgreSQL, logs ni outbox; el payload durable y el candidato todavía pueden contener PII y requieren cifrado, acceso y retención aprobados en el target. El outbox conserva sólo metadata operativa. `pending_policy` impide que la ingesta autorice contacto.
- **Entrega**: Google documenta reintentos sobre 5xx y duplicados posibles; por eso sólo se devuelve 200 después del commit durable. 4xx representa rechazo no reintentable.
- **Límites**: 1 MiB, POST JSON, tenant/organización requeridos, forward compatibility mediante unknown fields ignorados y raw preservado.
- **No claims**: no demuestra cuenta live, webhook público, consentimiento, respuesta automática, cita, venta, conversion postback, Meta/TikTok, observabilidad target ni exactly-once. `BANNED`, estados desconocidos, fecha/buyer ambiguos quedan como rechazo durable, nunca como lead inventado.

## 4. Exact file manifest

```text
CREATE internal/leadstream/event.go
CREATE internal/leadstream/google_ads.go
CREATE internal/leadstream/store.go
CREATE internal/leadstream/handler.go
CREATE internal/leadstream/leadstream_test.go
CREATE internal/leadstream/mercadolibre_questions.go
CREATE internal/leadstream/mercadolibre_questions_test.go
CREATE internal/platform/postgres/lead_ingress.go
CREATE internal/platform/postgres/lead_ingress_integration_test.go
CREATE db/migrations/0044_omnichannel_lead_ingress.up.sql
CREATE db/migrations/0044_omnichannel_lead_ingress.down.sql
CREATE db/tests/0044_omnichannel_lead_ingress.test.sql
```

## 5. Materialization blocks

### FILE: `internal/leadstream/event.go`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:internal/leadstream/event.go:v1"
operation: CREATE
provenance: ADAPTED
source: "CloudEvents 1.0.2 + official provider retry contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "60fee6649a1dc6f2cc65d5a11bd6e58bb6b4196d98c734ef8c248c189d547178"
variables: []
secrets_allowed: false
```
````go
// Package leadstream defines the provider-neutral, fail-closed boundary for
// lead events. Provider payloads remain byte-exact while normalized candidates
// are versioned separately; receiving a payload never grants permission to
// contact a person or to write a business lead automatically.
package leadstream

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const MaxPayloadBytes int64 = 1 << 20

var (
	ErrInvalidEvent       = errors.New("leadstream: invalid event")
	ErrInvalidCredential  = errors.New("leadstream: invalid provider credential")
	ErrInvalidPayload     = errors.New("leadstream: invalid provider payload")
	ErrDivergentDuplicate = errors.New("leadstream: divergent duplicate")
	providerRe            = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)
)

type RawEvent struct {
	TenantID         string
	OrganizationID   string
	Provider         string
	ProviderEventID  string
	EventType        string
	Source           string
	SchemaVersion    string
	OccurredAt       time.Time
	ReceivedAt       time.Time
	Payload          []byte // durable representation; provider credentials must be removed first
	SourceSHA256     string // hash of the exact provider body before redaction
	StoredSHA256     string // hash of Payload
	RedactionProfile string
}

func NewRawEvent(tenantID, organizationID, provider, eventID, eventType, source, schemaVersion string, occurredAt, receivedAt time.Time, payload []byte) (RawEvent, error) {
	return NewRedactedRawEvent(tenantID, organizationID, provider, eventID, eventType, source, schemaVersion, occurredAt, receivedAt, payload, payload, "none:v1")
}

// NewRedactedRawEvent binds a safe durable representation to the hash of the
// exact source body. It never retains sourcePayload and therefore prevents a
// provider credential removed by the caller from being persisted accidentally.
func NewRedactedRawEvent(tenantID, organizationID, provider, eventID, eventType, source, schemaVersion string, occurredAt, receivedAt time.Time, sourcePayload, storedPayload []byte, redactionProfile string) (RawEvent, error) {
	sourceHash := sha256.Sum256(sourcePayload)
	storedHash := sha256.Sum256(storedPayload)
	e := RawEvent{
		TenantID: tenantID, OrganizationID: organizationID, Provider: provider,
		ProviderEventID: eventID, EventType: eventType, Source: source,
		SchemaVersion: schemaVersion, OccurredAt: occurredAt.UTC(), ReceivedAt: receivedAt.UTC(),
		Payload: append([]byte(nil), storedPayload...), SourceSHA256: hex.EncodeToString(sourceHash[:]), StoredSHA256: hex.EncodeToString(storedHash[:]),
		RedactionProfile: redactionProfile,
	}
	return e, e.Validate()
}

func (e RawEvent) Validate() error {
	if strings.TrimSpace(e.TenantID) == "" || strings.TrimSpace(e.OrganizationID) == "" {
		return fmt.Errorf("%w: tenant and organization required", ErrInvalidEvent)
	}
	if !providerRe.MatchString(e.Provider) || strings.TrimSpace(e.ProviderEventID) == "" || len(e.ProviderEventID) > 256 {
		return fmt.Errorf("%w: provider identity", ErrInvalidEvent)
	}
	if strings.TrimSpace(e.EventType) == "" || strings.TrimSpace(e.Source) == "" || strings.TrimSpace(e.SchemaVersion) == "" {
		return fmt.Errorf("%w: event context", ErrInvalidEvent)
	}
	if e.ReceivedAt.IsZero() || e.OccurredAt.IsZero() || len(e.Payload) == 0 || int64(len(e.Payload)) > MaxPayloadBytes {
		return fmt.Errorf("%w: time or payload", ErrInvalidEvent)
	}
	if strings.TrimSpace(e.RedactionProfile) == "" || len(e.RedactionProfile) > 128 {
		return fmt.Errorf("%w: redaction profile", ErrInvalidEvent)
	}
	h := sha256.Sum256(e.Payload)
	if e.StoredSHA256 != hex.EncodeToString(h[:]) || !validSHA256(e.SourceSHA256) {
		return fmt.Errorf("%w: payload hash", ErrInvalidEvent)
	}
	return nil
}

func validSHA256(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

type Field struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type LeadCandidate struct {
	TenantID           string    `json:"tenant_id"`
	OrganizationID     string    `json:"organization_id"`
	Provider           string    `json:"provider"`
	ProviderLeadID     string    `json:"provider_lead_id"`
	FormID             string    `json:"form_id,omitempty"`
	CampaignID         string    `json:"campaign_id,omitempty"`
	AdGroupID          string    `json:"ad_group_id,omitempty"`
	CreativeID         string    `json:"creative_id,omitempty"`
	AssetGroupID       string    `json:"asset_group_id,omitempty"`
	ClickID            string    `json:"click_id,omitempty"`
	LeadStage          string    `json:"lead_stage,omitempty"`
	SourceKind         string    `json:"source_kind,omitempty"`
	SubmittedAt        time.Time `json:"submitted_at"`
	IsTest             bool      `json:"is_test"`
	Fields             []Field   `json:"fields"`
	ContactEligibility string    `json:"contact_eligibility"`
}

func (c LeadCandidate) Validate() error {
	if c.TenantID == "" || c.OrganizationID == "" || !providerRe.MatchString(c.Provider) || c.ProviderLeadID == "" || c.SubmittedAt.IsZero() {
		return fmt.Errorf("%w: candidate identity", ErrInvalidPayload)
	}
	if c.ContactEligibility != "pending_policy" {
		return fmt.Errorf("%w: ingress cannot grant contact eligibility", ErrInvalidPayload)
	}
	seen := map[string]struct{}{}
	for _, f := range c.Fields {
		if strings.TrimSpace(f.ID) == "" || strings.TrimSpace(f.Value) == "" {
			return fmt.Errorf("%w: empty field", ErrInvalidPayload)
		}
		if _, ok := seen[f.ID]; ok {
			return fmt.Errorf("%w: duplicate field %s", ErrInvalidPayload, f.ID)
		}
		seen[f.ID] = struct{}{}
	}
	return nil
}

func StableUUID(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80
	hexv := hex.EncodeToString(b)
	return hexv[0:8] + "-" + hexv[8:12] + "-" + hexv[12:16] + "-" + hexv[16:20] + "-" + hexv[20:32]
}
````

### FILE: `internal/leadstream/google_ads.go`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:internal/leadstream/google_ads.go:v1"
operation: CREATE
provenance: ADAPTED
source: "https://developers.google.com/google-ads/webhook/docs/implementation"
license: "LicenseRef-Workspace-Owner"
sha256: "42d83974aeeca63375d0c5a91506a7189e9040584f84832261c1262466280ad1"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type googleLeadColumn struct {
	ColumnName  string `json:"column_name"`
	StringValue string `json:"string_value"`
	ColumnID    string `json:"column_id"`
}

type googleWebhookLead struct {
	LeadID         string             `json:"lead_id"`
	UserData       []googleLeadColumn `json:"user_column_data"`
	APIVersion     string             `json:"api_version"`
	FormID         json.Number        `json:"form_id"`
	CampaignID     json.Number        `json:"campaign_id"`
	GoogleKey      string             `json:"google_key"`
	IsTest         bool               `json:"is_test"`
	GCLID          string             `json:"gcl_id"`
	AdGroupID      json.Number        `json:"adgroup_id"`
	CreativeID     json.Number        `json:"creative_id"`
	AssetGroupID   json.Number        `json:"asset_group_id"`
	LeadStage      string             `json:"lead_stage"`
	LeadSubmitTime string             `json:"lead_submit_time"`
	LeadSource     string             `json:"lead_source"`
}

type GoogleDecoded struct {
	Raw               RawEvent
	Candidate         *LeadCandidate
	NormalizationCode string
}

// DecodeGoogleWebhook implements the public Google Lead Form webhook contract.
// Unknown JSON fields are deliberately ignored for forward compatibility.
// Once the provider key is authenticated, a structurally invalid lead still
// returns a RawEvent so the durable store can retain rejected evidence.
func DecodeGoogleWebhook(payload []byte, expectedKey, tenantID, organizationID string, receivedAt time.Time) (GoogleDecoded, error) {
	if len(payload) == 0 || int64(len(payload)) > MaxPayloadBytes {
		return GoogleDecoded{}, fmt.Errorf("%w: payload size", ErrInvalidPayload)
	}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var in googleWebhookLead
	if err := dec.Decode(&in); err != nil {
		return GoogleDecoded{}, fmt.Errorf("%w: json", ErrInvalidPayload)
	}
	if err := ensureEOF(dec); err != nil {
		return GoogleDecoded{}, err
	}
	if !secretEqual(in.GoogleKey, expectedKey) {
		return GoogleDecoded{}, ErrInvalidCredential
	}
	storedPayload, err := redactGoogleCredential(payload)
	if err != nil {
		return GoogleDecoded{}, err
	}

	eventID := strings.TrimSpace(in.LeadID)
	if eventID == "" {
		eventID = "sha256:" + payloadHash(payload)
	}
	schema := strings.TrimSpace(in.APIVersion)
	if schema == "" {
		schema = "unspecified"
	}
	occurred := receivedAt.UTC()
	code := ""
	if in.LeadSubmitTime != "" {
		parsed, err := time.Parse(time.RFC3339, in.LeadSubmitTime)
		if err != nil {
			code = "INVALID_SUBMIT_TIME"
		} else {
			occurred = parsed.UTC()
		}
	}
	raw, err := NewRedactedRawEvent(tenantID, organizationID, "google_ads", eventID, "google.ads.lead.v"+schema, "https://googleads.googleapis.com/lead-form", schema, occurred, receivedAt, payload, storedPayload, "google_ads:remove-google_key:v1")
	if err != nil {
		return GoogleDecoded{}, err
	}
	if strings.TrimSpace(in.LeadID) == "" {
		return GoogleDecoded{Raw: raw, NormalizationCode: "MISSING_LEAD_ID"}, nil
	}
	fields := make([]Field, 0, len(in.UserData))
	for _, f := range in.UserData {
		fields = append(fields, Field{ID: strings.TrimSpace(f.ColumnID), Value: strings.TrimSpace(f.StringValue)})
	}
	candidate := LeadCandidate{
		TenantID: tenantID, OrganizationID: organizationID, Provider: "google_ads",
		ProviderLeadID: in.LeadID, FormID: in.FormID.String(), CampaignID: in.CampaignID.String(),
		AdGroupID: in.AdGroupID.String(), CreativeID: in.CreativeID.String(), AssetGroupID: in.AssetGroupID.String(),
		ClickID: in.GCLID, LeadStage: in.LeadStage, SourceKind: in.LeadSource,
		SubmittedAt: occurred, IsTest: in.IsTest, Fields: fields, ContactEligibility: "pending_policy",
	}
	if code == "" {
		if err := candidate.Validate(); err != nil {
			code = "INVALID_NORMALIZED_LEAD"
		}
	}
	if code != "" {
		return GoogleDecoded{Raw: raw, NormalizationCode: code}, nil
	}
	return GoogleDecoded{Raw: raw, Candidate: &candidate}, nil
}

func redactGoogleCredential(payload []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var object map[string]json.RawMessage
	if err := dec.Decode(&object); err != nil {
		return nil, fmt.Errorf("%w: redaction json", ErrInvalidPayload)
	}
	if err := ensureEOF(dec); err != nil {
		return nil, err
	}
	delete(object, "google_key")
	redacted, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("%w: redaction", ErrInvalidPayload)
	}
	return redacted, nil
}

func ensureEOF(dec *json.Decoder) error {
	var extra any
	if err := dec.Decode(&extra); err == io.EOF {
		return nil
	}
	return fmt.Errorf("%w: trailing json", ErrInvalidPayload)
}

func secretEqual(got, want string) bool {
	if got == "" || want == "" || len(got) != len(want) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}

func payloadHash(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}
````

### FILE: `internal/leadstream/store.go`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:internal/leadstream/store.go:v1"
operation: CREATE
provenance: ADAPTED
source: "AWS SQS at-least-once/idempotency + Google Cloud retry/dead-letter guidance"
license: "LicenseRef-Workspace-Owner"
sha256: "241b5a7a617910ac56326a8649c3cbf32fe6a5be0f763db72f4cc5f57070ae16"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"context"
	"errors"
	"sync"
)

type ReceiptState string

const (
	ReceiptNormalized ReceiptState = "normalized"
	ReceiptRejected   ReceiptState = "rejected"
	ReceiptDuplicate  ReceiptState = "duplicate"
)

type Receipt struct {
	State        ReceiptState
	SourceSHA256 string
}

type Store interface {
	Record(ctx context.Context, raw RawEvent, candidate *LeadCandidate, normalizationCode string) (Receipt, error)
}

type memoryRecord struct {
	raw       RawEvent
	candidate *LeadCandidate
	code      string
}

// MemoryStore is only a deterministic test/reference adapter. Production must
// inject a durable Store such as postgres.LeadIngress.
type MemoryStore struct {
	mu      sync.Mutex
	records map[string]memoryRecord
	Fail    error
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{records: map[string]memoryRecord{}} }

func eventKey(raw RawEvent) string {
	return raw.TenantID + "\x00" + raw.Provider + "\x00" + raw.ProviderEventID
}

func (s *MemoryStore) Record(_ context.Context, raw RawEvent, candidate *LeadCandidate, normalizationCode string) (Receipt, error) {
	if s == nil {
		return Receipt{}, errors.New("leadstream: nil store")
	}
	if err := raw.Validate(); err != nil {
		return Receipt{}, err
	}
	if candidate == nil && normalizationCode == "" {
		return Receipt{}, errors.New("leadstream: candidate or rejection required")
	}
	if candidate != nil {
		if normalizationCode != "" {
			return Receipt{}, errors.New("leadstream: candidate and rejection are mutually exclusive")
		}
		if err := candidate.Validate(); err != nil {
			return Receipt{}, err
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Fail != nil {
		return Receipt{}, s.Fail
	}
	key := eventKey(raw)
	if prior, ok := s.records[key]; ok {
		if prior.raw.SourceSHA256 != raw.SourceSHA256 {
			return Receipt{}, ErrDivergentDuplicate
		}
		return Receipt{State: ReceiptDuplicate, SourceSHA256: raw.SourceSHA256}, nil
	}
	var cloned *LeadCandidate
	if candidate != nil {
		copy := *candidate
		copy.Fields = append([]Field(nil), candidate.Fields...)
		cloned = &copy
	}
	s.records[key] = memoryRecord{raw: raw, candidate: cloned, code: normalizationCode}
	state := ReceiptNormalized
	if cloned == nil {
		state = ReceiptRejected
	}
	return Receipt{State: state, SourceSHA256: raw.SourceSHA256}, nil
}

func (s *MemoryStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.records)
}
````

### FILE: `internal/leadstream/handler.go`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:internal/leadstream/handler.go:v1"
operation: CREATE
provenance: ADAPTED
source: "https://developers.google.com/google-ads/webhook/docs/implementation"
license: "LicenseRef-Workspace-Owner"
sha256: "bb2eed2d0393407968bdb302f29aff868b28c4d1147ae0bc49c208230b640185"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"
)

type GoogleHandler struct {
	Store          Store
	TenantID       string
	OrganizationID string
	GoogleKey      string
	Clock          func() time.Time
	MaxBytes       int64
}

func (h GoogleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if r.Method != http.MethodPost || h.Store == nil || strings.TrimSpace(h.TenantID) == "" || strings.TrimSpace(h.OrganizationID) == "" || strings.TrimSpace(h.GoogleKey) == "" {
		writeGoogleError(w, http.StatusBadRequest, "invalid webhook configuration or request")
		return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("content-type"))
	if err != nil || mediaType != "application/json" {
		writeGoogleError(w, http.StatusUnsupportedMediaType, "content-type must be application/json")
		return
	}
	limit := h.MaxBytes
	if limit <= 0 || limit > MaxPayloadBytes {
		limit = MaxPayloadBytes
	}
	reader := http.MaxBytesReader(w, r.Body, limit)
	payload, err := io.ReadAll(reader)
	if err != nil {
		writeGoogleError(w, http.StatusRequestEntityTooLarge, "payload too large")
		return
	}
	now := time.Now().UTC()
	if h.Clock != nil {
		now = h.Clock().UTC()
	}
	decoded, err := DecodeGoogleWebhook(payload, h.GoogleKey, h.TenantID, h.OrganizationID, now)
	if errors.Is(err, ErrInvalidCredential) {
		writeGoogleError(w, http.StatusBadRequest, "provider verification failed")
		return
	}
	if err != nil {
		writeGoogleError(w, http.StatusBadRequest, "invalid webhook payload")
		return
	}
	_, err = h.Store.Record(r.Context(), decoded.Raw, decoded.Candidate, decoded.NormalizationCode)
	if errors.Is(err, ErrDivergentDuplicate) {
		writeGoogleError(w, http.StatusConflict, "provider event identity reused with different payload")
		return
	}
	if err != nil {
		writeGoogleError(w, http.StatusInternalServerError, "durable intake unavailable")
		return
	}
	if decoded.Candidate == nil {
		writeGoogleError(w, http.StatusBadRequest, "lead retained but normalization rejected")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}

func writeGoogleError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
````

### FILE: `internal/leadstream/leadstream_test.go`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:internal/leadstream/leadstream_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract tests against official Google webhook behavior"
license: "LicenseRef-Workspace-Owner"
sha256: "80c5bc3e9970536a1b453cfe62cb2abbd7151c7aaa7ba251e684d2f535202d94"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const validGooglePayload = `{"lead_id":"lead-123","user_column_data":[{"column_id":"FULL_NAME","string_value":"Ada Lovelace"},{"column_id":"EMAIL","string_value":"ada@example.com"}],"api_version":"3","form_id":9223372036854775806,"campaign_id":42,"google_key":"secret-key","is_test":false,"gcl_id":"click-1","adgroup_id":7,"creative_id":8,"asset_group_id":9,"lead_stage":"SUBMITTED","lead_submit_time":"2026-09-04T12:30:00Z","lead_source":"LEAD_FORM","future_field":{"kept":"in raw, ignored by parser"}}`

func fixedTime() time.Time { return time.Date(2026, 9, 4, 12, 31, 0, 0, time.UTC) }

func TestGoogleDecodePreservesSourceIdentityRedactsCredentialAndIgnoresUnknownFields(t *testing.T) {
	d, err := DecodeGoogleWebhook([]byte(validGooglePayload), "secret-key", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	if d.Candidate == nil || d.Candidate.FormID != "9223372036854775806" || d.Candidate.ContactEligibility != "pending_policy" {
		t.Fatalf("unexpected candidate: %+v", d.Candidate)
	}
	if d.Raw.ProviderEventID != "lead-123" || d.Raw.SourceSHA256 != payloadHash([]byte(validGooglePayload)) {
		t.Fatal("source hash or provider identity changed")
	}
	if d.Raw.StoredSHA256 != payloadHash(d.Raw.Payload) || d.Raw.RedactionProfile != "google_ads:remove-google_key:v1" {
		t.Fatal("stored payload evidence is not self-consistent")
	}
	if bytes.Contains(d.Raw.Payload, []byte("secret-key")) || bytes.Contains(d.Raw.Payload, []byte("google_key")) {
		t.Fatal("provider credential survived durable redaction")
	}
	if !bytes.Contains(d.Raw.Payload, []byte("future_field")) {
		t.Fatal("unknown provider field was not retained in redacted evidence")
	}
}

func TestGoogleDecodeRejectsWrongCredentialBeforeStorage(t *testing.T) {
	if _, err := DecodeGoogleWebhook([]byte(validGooglePayload), "wrong", "tenant-1", "store-1", fixedTime()); !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("expected credential error, got %v", err)
	}
}

func TestMemoryStoreExactReplayAndDivergence(t *testing.T) {
	d, _ := DecodeGoogleWebhook([]byte(validGooglePayload), "secret-key", "tenant-1", "store-1", fixedTime())
	s := NewMemoryStore()
	first, err := s.Record(context.Background(), d.Raw, d.Candidate, "")
	if err != nil || first.State != ReceiptNormalized {
		t.Fatalf("first record: %+v %v", first, err)
	}
	replay, err := s.Record(context.Background(), d.Raw, d.Candidate, "")
	if err != nil || replay.State != ReceiptDuplicate || s.Count() != 1 {
		t.Fatalf("replay: %+v %v count=%d", replay, err, s.Count())
	}
	changedPayload := strings.Replace(validGooglePayload, "Ada Lovelace", "Grace Hopper", 1)
	other, decodeErr := DecodeGoogleWebhook([]byte(changedPayload), "secret-key", "tenant-1", "store-1", fixedTime())
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	if _, err := s.Record(context.Background(), other.Raw, other.Candidate, ""); !errors.Is(err, ErrDivergentDuplicate) {
		t.Fatalf("expected divergent duplicate, got %v", err)
	}
}

func TestGoogleHandlerAcknowledgesOnlyAfterStore(t *testing.T) {
	s := NewMemoryStore()
	h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime}
	r := httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(validGooglePayload))
	r.Header.Set("content-type", "application/json; charset=utf-8")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK || strings.TrimSpace(w.Body.String()) != "{}" || s.Count() != 1 {
		t.Fatalf("status=%d body=%q count=%d", w.Code, w.Body.String(), s.Count())
	}
	s.Fail = errors.New("durability unavailable")
	r = httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(strings.Replace(validGooglePayload, "lead-123", "lead-124", 1)))
	r.Header.Set("content-type", "application/json")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("store failure must be retryable 5xx, got %d", w.Code)
	}
}

func TestGoogleHandlerExactReplayAndDivergentIdentity(t *testing.T) {
	s := NewMemoryStore()
	h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime}
	call := func(payload string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(payload))
		r.Header.Set("content-type", "application/json")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	if w := call(validGooglePayload); w.Code != http.StatusOK {
		t.Fatalf("first delivery status=%d body=%q", w.Code, w.Body.String())
	}
	if w := call(validGooglePayload); w.Code != http.StatusOK || s.Count() != 1 {
		t.Fatalf("exact replay status=%d count=%d", w.Code, s.Count())
	}
	changed := strings.Replace(validGooglePayload, "Ada Lovelace", "Grace Hopper", 1)
	if w := call(changed); w.Code != http.StatusConflict || s.Count() != 1 {
		t.Fatalf("divergent identity status=%d count=%d", w.Code, s.Count())
	}
}

func TestGoogleHandlerRequestBoundary(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		payload     string
		maxBytes    int64
		want        int
	}{
		{name: "method", method: http.MethodGet, contentType: "application/json", payload: validGooglePayload, want: http.StatusBadRequest},
		{name: "content type", method: http.MethodPost, contentType: "text/plain", payload: validGooglePayload, want: http.StatusUnsupportedMediaType},
		{name: "too large", method: http.MethodPost, contentType: "application/json", payload: validGooglePayload, maxBytes: 8, want: http.StatusRequestEntityTooLarge},
		{name: "malformed", method: http.MethodPost, contentType: "application/json", payload: `{`, want: http.StatusBadRequest},
		{name: "trailing json", method: http.MethodPost, contentType: "application/json", payload: validGooglePayload + `{}`, want: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryStore()
			h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime, MaxBytes: tt.maxBytes}
			r := httptest.NewRequest(tt.method, "/webhooks/google-ads", strings.NewReader(tt.payload))
			r.Header.Set("content-type", tt.contentType)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tt.want || s.Count() != 0 {
				t.Fatalf("status=%d want=%d stored=%d", w.Code, tt.want, s.Count())
			}
		})
	}
}

func TestGoogleHandlerRejectsCredentialAndRetainsAuthenticatedNormalizationFailure(t *testing.T) {
	s := NewMemoryStore()
	h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime}
	badKey := strings.Replace(validGooglePayload, "secret-key", "wrong-key", 1)
	r := httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(badKey))
	r.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest || s.Count() != 0 {
		t.Fatalf("untrusted payload stored: status=%d count=%d", w.Code, s.Count())
	}
	badTime := strings.Replace(validGooglePayload, "2026-09-04T12:30:00Z", "not-a-time", 1)
	r = httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(badTime))
	r.Header.Set("content-type", "application/json")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest || s.Count() != 1 {
		t.Fatalf("authenticated rejection not retained: status=%d count=%d", w.Code, s.Count())
	}
}

func TestCandidateRejectsAmbiguousDuplicateField(t *testing.T) {
	p := strings.Replace(validGooglePayload, `{"column_id":"EMAIL","string_value":"ada@example.com"}`, `{"column_id":"FULL_NAME","string_value":"Other"}`, 1)
	d, err := DecodeGoogleWebhook([]byte(p), "secret-key", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate != nil || d.NormalizationCode == "" {
		t.Fatalf("ambiguous lead must become retained rejection: %+v %v", d, err)
	}
}

func TestMissingProviderLeadIDIsRetainedAsRejectedEvidence(t *testing.T) {
	p := strings.Replace(validGooglePayload, `"lead_id":"lead-123"`, `"lead_id":""`, 1)
	d, err := DecodeGoogleWebhook([]byte(p), "secret-key", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate != nil || d.NormalizationCode != "MISSING_LEAD_ID" || !strings.HasPrefix(d.Raw.ProviderEventID, "sha256:") {
		t.Fatalf("unexpected decode: %+v %v", d, err)
	}
	s := NewMemoryStore()
	h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime}
	r := httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", bytes.NewReader([]byte(p)))
	r.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest || s.Count() != 1 {
		t.Fatalf("status=%d stored=%d", w.Code, s.Count())
	}
}

func TestMemoryStoreConcurrentExactReplayHasOneDurableIdentity(t *testing.T) {
	d, err := DecodeGoogleWebhook([]byte(validGooglePayload), "secret-key", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	s := NewMemoryStore()
	const workers = 32
	states := make(chan ReceiptState, workers)
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			receipt, recordErr := s.Record(context.Background(), d.Raw, d.Candidate, "")
			states <- receipt.State
			errs <- recordErr
		}()
	}
	wg.Wait()
	close(states)
	close(errs)
	for recordErr := range errs {
		if recordErr != nil {
			t.Fatal(recordErr)
		}
	}
	normalized, duplicates := 0, 0
	for state := range states {
		switch state {
		case ReceiptNormalized:
			normalized++
		case ReceiptDuplicate:
			duplicates++
		default:
			t.Fatalf("unexpected state %q", state)
		}
	}
	if normalized != 1 || duplicates != workers-1 || s.Count() != 1 {
		t.Fatalf("normalized=%d duplicates=%d count=%d", normalized, duplicates, s.Count())
	}
}
````

### FILE: `internal/platform/postgres/lead_ingress.go`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:internal/platform/postgres/lead_ingress.go:v1"
operation: CREATE
provenance: ADAPTED
source: "PostgreSQL 18 transaction/constraint semantics + existing Elite outbox owner"
license: "LicenseRef-Workspace-Owner"
sha256: "f84ea9d5940ab5bac04b14d71dfa2a8fa4e4d8041181162fc4acce278a04a842"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeadIngress struct{ pool *pgxpool.Pool }

func NewLeadIngress(pool *pgxpool.Pool) *LeadIngress { return &LeadIngress{pool: pool} }

func (r *LeadIngress) Record(ctx context.Context, raw leadstream.RawEvent, candidate *leadstream.LeadCandidate, normalizationCode string) (leadstream.Receipt, error) {
	if r == nil || r.pool == nil {
		return leadstream.Receipt{}, errors.New("postgres lead ingress: nil pool")
	}
	if err := raw.Validate(); err != nil {
		return leadstream.Receipt{}, err
	}
	if candidate == nil && normalizationCode == "" {
		return leadstream.Receipt{}, errors.New("postgres lead ingress: candidate or rejection required")
	}
	if candidate != nil {
		if normalizationCode != "" {
			return leadstream.Receipt{}, errors.New("postgres lead ingress: candidate and rejection conflict")
		}
		if err := candidate.Validate(); err != nil {
			return leadstream.Receipt{}, err
		}
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return leadstream.Receipt{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	state := leadstream.ReceiptNormalized
	if candidate == nil {
		state = leadstream.ReceiptRejected
	}
	result, err := tx.Exec(ctx, `insert into integration.lead_ingress_raw
(tenant_id,organization_id,provider,provider_event_id,event_type,source,schema_version,occurred_at,received_at,payload_redacted,source_payload_sha256,stored_payload_sha256,redaction_profile,state,normalization_error_code)
values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,nullif($15,'')) on conflict do nothing`,
		raw.TenantID, raw.OrganizationID, raw.Provider, raw.ProviderEventID, raw.EventType, raw.Source,
		raw.SchemaVersion, raw.OccurredAt, raw.ReceivedAt, raw.Payload, raw.SourceSHA256, raw.StoredSHA256, raw.RedactionProfile, string(state), normalizationCode)
	if err != nil {
		return leadstream.Receipt{}, err
	}
	if result.RowsAffected() == 0 {
		var hash string
		if err := tx.QueryRow(ctx, `select source_payload_sha256 from integration.lead_ingress_raw where tenant_id=$1 and provider=$2 and provider_event_id=$3`, raw.TenantID, raw.Provider, raw.ProviderEventID).Scan(&hash); err != nil {
			return leadstream.Receipt{}, err
		}
		if hash != raw.SourceSHA256 {
			return leadstream.Receipt{}, leadstream.ErrDivergentDuplicate
		}
		if err := tx.Commit(ctx); err != nil {
			return leadstream.Receipt{}, err
		}
		return leadstream.Receipt{State: leadstream.ReceiptDuplicate, SourceSHA256: hash}, nil
	}
	if candidate != nil {
		fields, err := json.Marshal(candidate.Fields)
		if err != nil {
			return leadstream.Receipt{}, err
		}
		_, err = tx.Exec(ctx, `insert into integration.lead_candidate
(tenant_id,provider,provider_event_id,organization_id,provider_lead_id,form_id,campaign_id,ad_group_id,creative_id,asset_group_id,click_id,lead_stage,source_kind,submitted_at,is_test,fields,contact_eligibility)
values($1,$2,$3,$4,$5,nullif($6,''),nullif($7,''),nullif($8,''),nullif($9,''),nullif($10,''),nullif($11,''),nullif($12,''),nullif($13,''),$14,$15,$16,'pending_policy')`,
			candidate.TenantID, candidate.Provider, raw.ProviderEventID, candidate.OrganizationID, candidate.ProviderLeadID,
			candidate.FormID, candidate.CampaignID, candidate.AdGroupID, candidate.CreativeID, candidate.AssetGroupID,
			candidate.ClickID, candidate.LeadStage, candidate.SourceKind, candidate.SubmittedAt, candidate.IsTest, fields)
		if err != nil {
			return leadstream.Receipt{}, err
		}
	}
	eventPayload, _ := json.Marshal(map[string]any{
		"provider": raw.Provider, "provider_event_id": raw.ProviderEventID,
		"organization_id": raw.OrganizationID, "state": state, "is_test": candidate != nil && candidate.IsTest,
	})
	outboxID := leadstream.StableUUID(raw.TenantID, raw.Provider, raw.ProviderEventID, "received")
	_, err = tx.Exec(ctx, `insert into platform.outbox_event
(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
values($1,$2,'provider-lead',$3,1,$4,1,clock_timestamp(),$5)`, raw.TenantID, outboxID, raw.Provider+":"+raw.ProviderEventID,
		"provider-lead."+string(state), eventPayload)
	if err != nil {
		return leadstream.Receipt{}, fmt.Errorf("write lead ingress outbox: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return leadstream.Receipt{}, err
	}
	return leadstream.Receipt{State: state, SourceSHA256: raw.SourceSHA256}, nil
}
````

### FILE: `internal/platform/postgres/lead_ingress_integration_test.go`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:internal/platform/postgres/lead_ingress_integration_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL integration evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "993a73808c130bf0da7a293f918181bd8e3778818a585ef3250c9d8adba1556d"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestLeadIngressDurabilityReplayDivergenceAndConcurrency(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	tenant := leadstream.StableUUID(t.Name(), time.Now().UTC().Format(time.RFC3339Nano))
	tenantCode := "lead-" + tenant[:8]
	const organization = "lead-store"
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
		values($1,$2,'Lead Runtime Test','Lead Runtime Test')`, tenant, tenantCode); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
		values($1,$2,'lead-store','Lead Store','store')`, tenant, organization); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event
			set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null
			where tenant_id=$1 and published_at is null`, tenant)
	}()

	repo := NewLeadIngress(pool)
	when := time.Date(2020, 9, 4, 18, 0, 0, 0, time.UTC)
	makeLead := func(eventID string, payload []byte) (leadstream.RawEvent, *leadstream.LeadCandidate) {
		raw, makeErr := leadstream.NewRawEvent(tenant, organization, "google_ads", eventID, "google.ads.lead.v3", "https://googleads.googleapis.com/lead-form", "3", when, when.Add(time.Second), payload)
		if makeErr != nil {
			t.Fatal(makeErr)
		}
		return raw, &leadstream.LeadCandidate{
			TenantID: tenant, OrganizationID: organization, Provider: "google_ads", ProviderLeadID: eventID,
			FormID: "9223372036854775806", CampaignID: "42", SubmittedAt: when, IsTest: true,
			Fields: []leadstream.Field{{ID: "EMAIL", Value: "fixture@example.invalid"}}, ContactEligibility: "pending_policy",
		}
	}

	googlePayload := []byte(`{"lead_id":"lead-db-1","user_column_data":[{"column_id":"EMAIL","string_value":"fixture@example.invalid"}],"api_version":"3","google_key":"integration-secret","is_test":true,"lead_submit_time":"2020-09-04T18:00:00Z"}`)
	decoded, err := leadstream.DecodeGoogleWebhook(googlePayload, "integration-secret", tenant, organization, when.Add(time.Second))
	if err != nil || decoded.Candidate == nil {
		t.Fatalf("decode candidate=%+v err=%v", decoded.Candidate, err)
	}
	raw, candidate := decoded.Raw, decoded.Candidate
	receipt, err := repo.Record(ctx, raw, candidate, "")
	if err != nil || receipt.State != leadstream.ReceiptNormalized {
		t.Fatalf("first record state=%q err=%v", receipt.State, err)
	}
	receipt, err = repo.Record(ctx, raw, candidate, "")
	if err != nil || receipt.State != leadstream.ReceiptDuplicate {
		t.Fatalf("exact replay state=%q err=%v", receipt.State, err)
	}
	changedPayload := []byte(`{"lead_id":"lead-db-1","user_column_data":[{"column_id":"EMAIL","string_value":"changed@example.invalid"}],"api_version":"3","google_key":"integration-secret","is_test":true,"lead_submit_time":"2020-09-04T18:00:00Z"}`)
	changed, decodeErr := leadstream.DecodeGoogleWebhook(changedPayload, "integration-secret", tenant, organization, when.Add(time.Second))
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	if _, err := repo.Record(ctx, changed.Raw, changed.Candidate, ""); !errors.Is(err, leadstream.ErrDivergentDuplicate) {
		t.Fatalf("expected divergent duplicate, got %v", err)
	}
	var rawCount, candidateCount, outboxCount int
	var sourceHash, storedHash, storedPayload, redactionProfile string
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider_event_id='lead-db-1'`, tenant).Scan(&rawCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_candidate where tenant_id=$1 and provider_event_id='lead-db-1' and contact_eligibility='pending_policy'`, tenant).Scan(&candidateCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='google_ads:lead-db-1'`, tenant).Scan(&outboxCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select source_payload_sha256,stored_payload_sha256,convert_from(payload_redacted,'UTF8'),redaction_profile
		from integration.lead_ingress_raw where tenant_id=$1 and provider='google_ads' and provider_event_id='lead-db-1'`, tenant).Scan(&sourceHash, &storedHash, &storedPayload, &redactionProfile); err != nil {
		t.Fatal(err)
	}
	if rawCount != 1 || candidateCount != 1 || outboxCount != 1 {
		t.Fatalf("raw=%d candidate=%d outbox=%d", rawCount, candidateCount, outboxCount)
	}
	if sourceHash != raw.SourceSHA256 || storedHash != raw.StoredSHA256 || redactionProfile != raw.RedactionProfile ||
		strings.Contains(storedPayload, "integration-secret") || strings.Contains(storedPayload, "google_key") {
		t.Fatalf("unsafe or inconsistent stored evidence source=%s stored=%s profile=%q payload=%q", sourceHash, storedHash, redactionProfile, storedPayload)
	}

	rejectedRaw, makeErr := leadstream.NewRawEvent(tenant, organization, "google_ads", "reject-db-1", "google.ads.lead.v3", "https://googleads.googleapis.com/lead-form", "3", when, when.Add(time.Second), []byte(`{"authenticated":"but-unmappable"}`))
	if makeErr != nil {
		t.Fatal(makeErr)
	}
	receipt, err = repo.Record(ctx, rejectedRaw, nil, "UNMAPPABLE_FIXTURE")
	if err != nil || receipt.State != leadstream.ReceiptRejected {
		t.Fatalf("rejection state=%q err=%v", receipt.State, err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_candidate where tenant_id=$1 and provider_event_id='reject-db-1'`, tenant).Scan(&candidateCount); err != nil || candidateCount != 0 {
		t.Fatalf("rejected candidate count=%d err=%v", candidateCount, err)
	}

	concurrentRaw, concurrentCandidate := makeLead("lead-db-concurrent", []byte(`{"lead_id":"lead-db-concurrent"}`))
	const workers = 16
	states := make(chan leadstream.ReceiptState, workers)
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, recordErr := repo.Record(ctx, concurrentRaw, concurrentCandidate, "")
			states <- got.State
			errs <- recordErr
		}()
	}
	wg.Wait()
	close(states)
	close(errs)
	for recordErr := range errs {
		if recordErr != nil {
			t.Fatal(recordErr)
		}
	}
	normalized, duplicates := 0, 0
	for state := range states {
		switch state {
		case leadstream.ReceiptNormalized:
			normalized++
		case leadstream.ReceiptDuplicate:
			duplicates++
		default:
			t.Fatalf("unexpected concurrent state %q", state)
		}
	}
	if normalized != 1 || duplicates != workers-1 {
		t.Fatalf("normalized=%d duplicates=%d", normalized, duplicates)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider_event_id='lead-db-concurrent'`, tenant).Scan(&rawCount); err != nil || rawCount != 1 {
		t.Fatalf("concurrent raw count=%d err=%v", rawCount, err)
	}
}
````

### FILE: `db/migrations/0044_omnichannel_lead_ingress.up.sql`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:db/migrations/0044_omnichannel_lead_ingress.up.sql:v1"
operation: CREATE
provenance: ADAPTED
source: "PostgreSQL 18 constraints, unique identity and trigger semantics"
license: "LicenseRef-Workspace-Owner"
sha256: "761153e744342c8bec1a77d5e93d81ee794189037276cbb051f216fdcd2a91cf"
variables: []
secrets_allowed: false
```
````sql
begin;

create table integration.lead_ingress_raw (
  tenant_id uuid not null,
  organization_id text not null,
  provider text not null check (provider ~ '^[a-z][a-z0-9_]{0,31}$'),
  provider_event_id text not null check (length(provider_event_id) between 1 and 256),
  event_type text not null,
  source text not null,
  schema_version text not null,
  occurred_at timestamptz not null,
  received_at timestamptz not null,
  payload_redacted bytea not null,
  source_payload_sha256 text not null check (source_payload_sha256 ~ '^[0-9a-f]{64}$'),
  stored_payload_sha256 text not null check (stored_payload_sha256 ~ '^[0-9a-f]{64}$'),
  redaction_profile text not null check (length(redaction_profile) between 1 and 128),
  state text not null check (state in ('normalized','rejected')),
  normalization_error_code text,
  primary key (tenant_id, provider, provider_event_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  check (octet_length(payload_redacted) between 1 and 1048576),
  check ((state='normalized' and normalization_error_code is null) or (state='rejected' and length(normalization_error_code) between 1 and 128))
);

create index lead_ingress_received_idx on integration.lead_ingress_raw
  (tenant_id, provider, received_at desc, provider_event_id);

create table integration.lead_candidate (
  tenant_id uuid not null,
  provider text not null,
  provider_event_id text not null,
  organization_id text not null,
  provider_lead_id text not null,
  form_id text,
  campaign_id text,
  ad_group_id text,
  creative_id text,
  asset_group_id text,
  click_id text,
  lead_stage text,
  source_kind text,
  submitted_at timestamptz not null,
  is_test boolean not null,
  fields jsonb not null,
  contact_eligibility text not null check (contact_eligibility='pending_policy'),
  primary key (tenant_id, provider, provider_event_id),
  foreign key (tenant_id, provider, provider_event_id)
    references integration.lead_ingress_raw (tenant_id, provider, provider_event_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  check (jsonb_typeof(fields)='array'),
  check (length(provider_lead_id) between 1 and 256)
);

create function integration.reject_lead_ingress_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='lead ingress evidence is append-only';
end;
$function$;

create trigger lead_ingress_raw_immutable before update or delete on integration.lead_ingress_raw
for each row execute function integration.reject_lead_ingress_mutation();

create trigger lead_candidate_immutable before update or delete on integration.lead_candidate
for each row execute function integration.reject_lead_ingress_mutation();

commit;
````

### FILE: `db/migrations/0044_omnichannel_lead_ingress.down.sql`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:db/migrations/0044_omnichannel_lead_ingress.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local reversible migration"
license: "LicenseRef-Workspace-Owner"
sha256: "1cd0fb4673bdcb0d21004fd16677b64f3713c47144bb49b67b0155e7813f5167"
variables: []
secrets_allowed: false
```
````sql
begin;

drop trigger if exists lead_candidate_immutable on integration.lead_candidate;
drop trigger if exists lead_ingress_raw_immutable on integration.lead_ingress_raw;
drop function if exists integration.reject_lead_ingress_mutation();
drop table if exists integration.lead_candidate;
drop table if exists integration.lead_ingress_raw;

commit;
````

### FILE: `db/tests/0044_omnichannel_lead_ingress.test.sql`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:db/tests/0044_omnichannel_lead_ingress.test.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL invariant tests"
license: "LicenseRef-Workspace-Owner"
sha256: "b4f1f2d8933424df91215abf69479c61e5a61bd6ffa0105996a26460085e7c10"
variables: []
secrets_allowed: false
```
````sql
begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('44444444-4444-4444-8444-444444444444','lead-ingress-test','Lead Ingress Test','Lead Ingress Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values('44444444-4444-4444-8444-444444444444','store-1','store-1','Store 1','store');

insert into integration.lead_ingress_raw
(tenant_id,organization_id,provider,provider_event_id,event_type,source,schema_version,occurred_at,received_at,payload_redacted,source_payload_sha256,stored_payload_sha256,redaction_profile,state)
values('44444444-4444-4444-8444-444444444444','store-1','google_ads','lead-1','google.ads.lead.v3','https://googleads.googleapis.com/lead-form','3',clock_timestamp(),clock_timestamp(),convert_to('{"lead_id":"lead-1"}','UTF8'),repeat('a',64),repeat('b',64),'google_ads:remove-google_key:v1','normalized');

insert into integration.lead_candidate
(tenant_id,provider,provider_event_id,organization_id,provider_lead_id,submitted_at,is_test,fields,contact_eligibility)
values('44444444-4444-4444-8444-444444444444','google_ads','lead-1','store-1','lead-1',clock_timestamp(),false,'[]','pending_policy');

do $test$
begin
  begin
    update integration.lead_ingress_raw set state='rejected'
    where tenant_id='44444444-4444-4444-8444-444444444444' and provider='google_ads' and provider_event_id='lead-1';
    raise exception 'raw update unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  begin
    update integration.lead_candidate set provider_lead_id='changed'
    where tenant_id='44444444-4444-4444-8444-444444444444' and provider='google_ads' and provider_event_id='lead-1';
    raise exception 'candidate update unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
end;
$test$;

rollback;
````

### FILE: `internal/leadstream/mercadolibre_questions.go`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:internal/leadstream/mercadolibre_questions.go:v1"
operation: CREATE
provenance: ADAPTED
source: "official Mercado Libre Questions API v4 contract + provider-neutral durable ingress boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "eaf8cc34902ca5a5045df56c7e7a7c7dc1b3f63dbe8b6fb03de1fa59cbec2639"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type mercadoLibreQuestionParty struct {
	ID        json.Number `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Phone     struct {
		AreaCode string `json:"area_code"`
		Number   string `json:"number"`
	} `json:"phone"`
}

type mercadoLibreQuestion struct {
	ID          json.Number                `json:"id"`
	SellerID    json.Number                `json:"seller_id"`
	BuyerID     json.Number                `json:"buyer_id"`
	ItemID      string                     `json:"item_id"`
	Status      string                     `json:"status"`
	Text        string                     `json:"text"`
	DateCreated string                     `json:"date_created"`
	From        *mercadoLibreQuestionParty `json:"from"`
}

type MercadoLibreQuestionDecoded struct {
	Raw               RawEvent
	Candidate         *LeadCandidate
	NormalizationCode string
}

// DecodeMercadoLibreQuestion normalizes the authoritative API v4 response
// fetched after a Mercado Libre `questions` notification. It does not verify a
// callback by itself and never grants permission to contact or answer a buyer.
func DecodeMercadoLibreQuestion(payload []byte, expectedSellerID, tenantID, organizationID string, receivedAt time.Time) (MercadoLibreQuestionDecoded, error) {
	if len(payload) == 0 || int64(len(payload)) > MaxPayloadBytes {
		return MercadoLibreQuestionDecoded{}, fmt.Errorf("%w: payload size", ErrInvalidPayload)
	}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var in mercadoLibreQuestion
	if err := dec.Decode(&in); err != nil {
		return MercadoLibreQuestionDecoded{}, fmt.Errorf("%w: json", ErrInvalidPayload)
	}
	if err := ensureEOF(dec); err != nil {
		return MercadoLibreQuestionDecoded{}, err
	}
	if strings.TrimSpace(expectedSellerID) == "" || in.SellerID.String() != expectedSellerID {
		return MercadoLibreQuestionDecoded{}, ErrInvalidCredential
	}
	questionID := in.ID.String()
	createdAt, parseErr := time.Parse(time.RFC3339Nano, in.DateCreated)
	if questionID == "" || questionID == "0" || strings.TrimSpace(in.ItemID) == "" || strings.TrimSpace(in.Status) == "" {
		return MercadoLibreQuestionDecoded{}, fmt.Errorf("%w: question identity", ErrInvalidPayload)
	}
	if parseErr != nil {
		createdAt = receivedAt.UTC()
	}
	raw, err := NewRedactedRawEvent(tenantID, organizationID, "mercadolibre", questionID, "mercadolibre.question.v4", "https://api.mercadolibre.com/questions/"+questionID+"?api_version=4", "4", createdAt, receivedAt, payload, payload, "mercadolibre:provider-response:v1")
	if err != nil {
		return MercadoLibreQuestionDecoded{}, err
	}
	if parseErr != nil {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "INVALID_DATE_CREATED"}, nil
	}
	status := strings.ToUpper(strings.TrimSpace(in.Status))
	if status == "BANNED" || status == "DELETED" || status == "DISABLED" || status == "UNDER_REVIEW" {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "QUESTION_" + status}, nil
	}
	if status != "UNANSWERED" && status != "ANSWERED" && status != "CLOSED_UNANSWERED" {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "UNKNOWN_QUESTION_STATUS"}, nil
	}
	text := strings.TrimSpace(in.Text)
	if text == "" {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "MISSING_QUESTION_TEXT"}, nil
	}
	buyerID := in.BuyerID.String()
	if in.From != nil && in.From.ID.String() != "" {
		if buyerID != "" && buyerID != "0" && buyerID != in.From.ID.String() {
			return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "AMBIGUOUS_BUYER_ID"}, nil
		}
		buyerID = in.From.ID.String()
	}
	if buyerID == "" || buyerID == "0" {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "MISSING_BUYER_ID"}, nil
	}
	fields := []Field{{ID: "ITEM_ID", Value: strings.TrimSpace(in.ItemID)}, {ID: "BUYER_ID", Value: buyerID}, {ID: "QUESTION_TEXT", Value: text}}
	if in.From != nil {
		appendField := func(id, value string) {
			if value = strings.TrimSpace(value); value != "" {
				fields = append(fields, Field{ID: id, Value: value})
			}
		}
		appendField("FIRST_NAME", in.From.FirstName)
		appendField("LAST_NAME", in.From.LastName)
		appendField("EMAIL", in.From.Email)
		appendField("PHONE_AREA_CODE", in.From.Phone.AreaCode)
		appendField("PHONE_NUMBER", in.From.Phone.Number)
	}
	candidate := LeadCandidate{
		TenantID: tenantID, OrganizationID: organizationID, Provider: "mercadolibre",
		ProviderLeadID: questionID, LeadStage: status, SourceKind: "MARKETPLACE_QUESTION",
		SubmittedAt: createdAt.UTC(), Fields: fields, ContactEligibility: "pending_policy",
	}
	if err := candidate.Validate(); err != nil {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "INVALID_NORMALIZED_LEAD"}, nil
	}
	return MercadoLibreQuestionDecoded{Raw: raw, Candidate: &candidate}, nil
}
````

### FILE: `internal/leadstream/mercadolibre_questions_test.go`
```yaml
block_id: "GO-OMNICHANNEL-LEAD-INGRESS:internal/leadstream/mercadolibre_questions_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact-contract, ownership, PII-policy and fail-closed regression tests"
license: "LicenseRef-Workspace-Owner"
sha256: "e2a2f240c96320b7bf38fcb047f2e618904711188b1bbcea9ffc2efd5e36be91"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

const validMercadoLibreQuestion = `{"id":11751825075,"seller_id":123456789,"buyer_id":56801932,"item_id":"MLA739200576","status":"UNANSWERED","text":"Necesito información","date_created":"2026-09-05T12:00:00Z","future_field":{"preserved":true}}`

func TestMercadoLibreQuestionNormalizesIntoCommonPendingPolicyCandidate(t *testing.T) {
	d, err := DecodeMercadoLibreQuestion([]byte(validMercadoLibreQuestion), "123456789", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	if d.Candidate == nil || d.Candidate.Provider != "mercadolibre" || d.Candidate.ProviderLeadID != "11751825075" || d.Candidate.ContactEligibility != "pending_policy" {
		t.Fatalf("unexpected candidate: %+v", d.Candidate)
	}
	if d.Raw.SourceSHA256 != payloadHash([]byte(validMercadoLibreQuestion)) || !bytes.Contains(d.Raw.Payload, []byte("future_field")) {
		t.Fatal("provider evidence identity or forward-compatible raw body was lost")
	}
	store := NewMemoryStore()
	receipt, err := store.Record(context.Background(), d.Raw, d.Candidate, "")
	if err != nil || receipt.State != ReceiptNormalized || store.Count() != 1 {
		t.Fatalf("receipt=%+v err=%v count=%d", receipt, err, store.Count())
	}
}

func TestMercadoLibreVehicleContactDataNeverGrantsContactEligibility(t *testing.T) {
	payload := `{"id":11949565740,"seller_id":123456789,"text":"Quiero coordinar","status":"UNANSWERED","item_id":"MLA595976788","date_created":"2026-09-05T12:07:18.109-04:00","from":{"id":21547449,"first_name":"Juan","last_name":"Lead","phone":{"number":"95712582","area_code":"9"},"email":"juan@example.com"}}`
	d, err := DecodeMercadoLibreQuestion([]byte(payload), "123456789", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate == nil {
		t.Fatalf("decode=%+v err=%v", d, err)
	}
	if d.Candidate.ContactEligibility != "pending_policy" || len(d.Candidate.Fields) != 8 {
		t.Fatalf("PII bypassed policy or was not normalized: %+v", d.Candidate)
	}
}

func TestMercadoLibreQuestionFailClosedBoundaries(t *testing.T) {
	if _, err := DecodeMercadoLibreQuestion([]byte(validMercadoLibreQuestion), "999999999", "tenant-1", "store-1", fixedTime()); !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("seller mismatch must fail before storage: %v", err)
	}
	banned := `{"id":11751825075,"seller_id":123456789,"buyer_id":56801932,"item_id":"MLA739200576","status":"BANNED","text":"","date_created":"2026-09-05T12:00:00Z"}`
	d, err := DecodeMercadoLibreQuestion([]byte(banned), "123456789", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate != nil || d.NormalizationCode != "QUESTION_BANNED" {
		t.Fatalf("banned question must be retained as rejection: %+v %v", d, err)
	}
	ambiguous := `{"id":11751825075,"seller_id":123456789,"buyer_id":56801932,"item_id":"MLA739200576","status":"UNANSWERED","text":"x","date_created":"2026-09-05T12:00:00Z","from":{"id":56801933}}`
	d, err = DecodeMercadoLibreQuestion([]byte(ambiguous), "123456789", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate != nil || d.NormalizationCode != "AMBIGUOUS_BUYER_ID" {
		t.Fatalf("ambiguous buyer must not normalize: %+v %v", d, err)
	}
}
````

## 6. Configuration surface

- `GoogleHandler.GoogleKey`: secreto de verificación recibido por canal de secretos; nunca se embebe ni registra.
- `TenantID` y `OrganizationID`: routing resuelto por configuración admitida, no por el payload del provider.
- `Store`: `postgres.NewLeadIngress(pool)` en runtime; `MemoryStore` sólo tests.
- `MaxBytes`: opcional, nunca mayor a 1 MiB.
- `TEST_DATABASE_URL`: sólo gates PostgreSQL; no es configuración productiva.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | Go 1.26.7 | HTTP, JSON, SHA-256, constant-time compare | BSD-3-Clause | runtime/build | https://go.dev |
| pgx/v5 | lock del backend compuesto | transacción PostgreSQL | MIT | runtime | https://github.com/jackc/pgx |
| PostgreSQL | 18.6 | evidencia raw/candidato/outbox | PostgreSQL | runtime/test | https://www.postgresql.org |

## 8. Apply order

1. Componer backend, PostgreSQL foundation, organización y outbox.
2. Materializar este pack sin colisiones y ejecutar migración 0044.
3. Inyectar `postgres.NewLeadIngress(pool)` en el composition root y montar `GoogleHandler` en un endpoint autenticado/observado.
4. Mantener el candidato en `pending_policy`; una capability separada y aprobada promueve a CRM/contacto.
5. Verificar focal, suite global, PostgreSQL limpio y round-trip down/up.
6. Rollback: retirar el route/worker y ejecutar 0044 down sólo bajo un plan de migración que preserve/exporte evidencia requerida.

## 9. Verification

- Go 1.26.7: tests focales de `leadstream` y `postgres` PASS.
- Frontera HTTP: método/tipo/tamaño/JSON/credencial/replay/divergencia/rechazo durable y ausencia de `google_key` en payload durable PASS.
- Concurrencia: 32 replays in-memory y 16 replays PostgreSQL producen una identidad durable.
- PostgreSQL 18.6 con checksums: 44 migraciones y 30 tests SQL PASS en base limpia.
- Suite compuesta con PostgreSQL real: 45 paquetes PASS; `go vet ./...` y `go build ./...` PASS.
- 0044 down/up/test PASS.
- `go test -race` sigue `BLOCKED_EXTERNAL` en este host por ausencia de CGO/GCC; es gate de promotion, no un PASS inferido.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_OMNICHANNEL_LEAD_INGRESS_2026-09-04_V224.md`.
