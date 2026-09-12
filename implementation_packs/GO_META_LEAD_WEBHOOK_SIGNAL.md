# Go Meta Lead Webhook Signal

## 1. Metadata

```yaml
pack_id: "GO-META-LEAD-WEBHOOK-SIGNAL"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el endpoint inmediato Meta Lead Ads con verificación SHA-256 sobre bytes crudos, challenge, batches completos, persistencia PostgreSQL previa al ACK, identidad idempotente, divergencia fail-closed y outbox hacia recuperación del lead."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0", "Meta Graph Webhooks"]
compatible_with: ["GO-OMNICHANNEL-LEAD-INGRESS 0.2.x", "PYTHON-META-LEAD-RECONCILIATION-ADAPTER 0.1.x", "GO-META-LEAD-EVIDENCE-IMPORT 0.1.x"]
incompatible_with: ["webhook sin firma", "ACK antes de commit", "primer elemento solamente", "log de payload o secreto", "webhook tratado como lead completo", "permisos o ownership inferidos"]
license_expression: "LicenseRef-Workspace-Owner AND LicenseRef-Meta-Platform-API-Only AND MIT"
upstream_sources: ["https://github.com/fbsamples/messenger-platform-samples/tree/354ee221ac1d081cc6105a1515a8468cc44f6710", "https://github.com/fbsamples/lead-ads-webhook-sample/tree/c0843165ad0f6acd82731e0c9da96a22b4162b35", "https://www.postgresql.org/docs/18/mvcc.html"]
verified_at: "2026-09-04"
```

Los dos archivos `internal/leadstream/meta_webhook*` son `ADAPTED`: combinan la implementación y regresiones oficiales actuales de firma `X-Hub-Signature-256` de Meta con el envelope `page/entry/changes/leadgen` del sample oficial Lead Ads. Los siete archivos restantes son `AUTHORED`. Ninguno se presenta como producto interno de Meta.

## 2. Applicability

Use para recibir la señal inmediata de formularios Meta. El webhook no contiene el lead completo: guarda el lote exacto, crea una señal `pending_retrieval` por `leadgen_id` y emite outbox. La recuperación de campos continúa en `PYTHON-META-LEAD-RECONCILIATION-ADAPTER`; la importación durable continúa en `GO-META-LEAD-EVIDENCE-IMPORT`.

El ejecutable exige `DATABASE_URL`, `META_APP_SECRET`, `META_VERIFY_TOKEN`, `TENANT_ID` y `ORGANIZATION_ID`. Nunca registre sus valores. Cuenta, permisos, Page/form ownership, suscripción, test lead y delivery live permanecen gates del proyecto.

## 3. Architecture contract

- GET acepta únicamente `hub.mode=subscribe`, challenge no vacío y verify token exacto en tiempo constante.
- POST acepta sólo JSON acotado, verifica `X-Hub-Signature-256` contra los bytes crudos antes de parsear y recorre todas las entradas/cambios `leadgen`.
- El lote exacto y cada señal se insertan en una transacción serializable; el `200 EVENT_RECEIVED` sale sólo tras commit.
- Replay semántico no duplica signal/outbox; el mismo `leadgen_id` con valor distinto se rechaza.
- Las tablas son append-only y tenant-scoped. El outbox no concede consentimiento ni crea un lead de negocio.
- El código archivado Lead Ads se usa sólo como contrato/fixture. Su controller inseguro no se incorpora.

## 4. Exact file manifest

```text
CREATE internal/leadstream/meta_webhook.go
CREATE internal/leadstream/meta_webhook_test.go
CREATE internal/platform/postgres/meta_lead_webhook.go
CREATE internal/platform/postgres/meta_lead_webhook_integration_test.go
CREATE cmd/meta-lead-webhook/main.go
CREATE db/migrations/0049_meta_lead_webhook_signal.up.sql
CREATE db/migrations/0049_meta_lead_webhook_signal.down.sql
CREATE db/tests/0049_meta_lead_webhook_signal.test.sql
CREATE third_party/meta-lead-webhook/official-source-lock.json
```

## 5. Materialization blocks

### FILE: `internal/leadstream/meta_webhook.go`
```yaml
block_id: "GO-META-LEAD-WEBHOOK-SIGNAL:internal/leadstream/meta_webhook.go:v1"
operation: CREATE
provenance: ADAPTED
source: "fbsamples/messenger-platform-samples@354ee221 and fbsamples/lead-ads-webhook-sample@c0843165; hardened local composition"
license: "LicenseRef-Workspace-Owner AND LicenseRef-Meta-Platform-API-Only AND MIT"
sha256: "48c65cd1b9a170de738b51b998240f942beb02dcf22b52c592510850b006cdc3"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"
)

var ErrInvalidMetaSignature = errors.New("leadstream: invalid Meta signature")

type MetaLeadSignal struct {
	LeadgenID   string
	FormID      string
	PageID      string
	AdID        string
	AdGroupID   string
	OccurredAt  time.Time
	ValueSHA256 string
}

type MetaWebhookBatch struct {
	TenantID       string
	OrganizationID string
	ReceivedAt     time.Time
	Payload        []byte
	PayloadSHA256  string
	Signals        []MetaLeadSignal
}

func (b MetaWebhookBatch) Validate() error {
	if strings.TrimSpace(b.TenantID) == "" || strings.TrimSpace(b.OrganizationID) == "" || b.ReceivedAt.IsZero() || len(b.Payload) == 0 || int64(len(b.Payload)) > MaxPayloadBytes || len(b.Signals) == 0 {
		return fmt.Errorf("%w: Meta batch identity", ErrInvalidPayload)
	}
	if payloadHash(b.Payload) != b.PayloadSHA256 {
		return fmt.Errorf("%w: Meta payload hash", ErrInvalidPayload)
	}
	seen := map[string]string{}
	for _, signal := range b.Signals {
		if strings.TrimSpace(signal.LeadgenID) == "" || len(signal.LeadgenID) > 256 || strings.TrimSpace(signal.FormID) == "" || len(signal.FormID) > 256 || strings.TrimSpace(signal.PageID) == "" || len(signal.PageID) > 256 || len(signal.AdID) > 256 || len(signal.AdGroupID) > 256 || signal.OccurredAt.IsZero() || !validSHA256(signal.ValueSHA256) {
			return fmt.Errorf("%w: Meta lead signal", ErrInvalidPayload)
		}
		if prior, ok := seen[signal.LeadgenID]; ok && prior != signal.ValueSHA256 {
			return ErrDivergentDuplicate
		}
		seen[signal.LeadgenID] = signal.ValueSHA256
	}
	return nil
}

type MetaWebhookReceipt struct {
	NewSignals       int
	DuplicateSignals int
}

type MetaWebhookStore interface {
	RecordMetaWebhook(context.Context, MetaWebhookBatch) (MetaWebhookReceipt, error)
}

type metaWebhookEnvelope struct {
	Object  string             `json:"object"`
	Entries []metaWebhookEntry `json:"entry"`
}

type metaWebhookEntry struct {
	PageID  string              `json:"id"`
	Time    json.Number         `json:"time"`
	Changes []metaWebhookChange `json:"changes"`
}

type metaWebhookChange struct {
	Field string          `json:"field"`
	Value json.RawMessage `json:"value"`
}

type metaLeadgenValue struct {
	FormID      string      `json:"form_id"`
	LeadgenID   string      `json:"leadgen_id"`
	CreatedTime json.Number `json:"created_time"`
	PageID      string      `json:"page_id"`
	AdID        string      `json:"ad_id"`
	AdGroupID   string      `json:"adgroup_id"`
}

// DecodeMetaLeadWebhook combines the current Meta X-Hub-Signature-256
// verification boundary with the official Lead Ads page/leadgen envelope.
// It records only the notification signal; full lead fields are retrieved by
// the separately pinned Meta Business SDK reconciliation adapter.
func DecodeMetaLeadWebhook(payload []byte, signatureHeader, appSecret, tenantID, organizationID string, receivedAt time.Time) (MetaWebhookBatch, error) {
	if len(payload) == 0 || int64(len(payload)) > MaxPayloadBytes || strings.TrimSpace(appSecret) == "" {
		return MetaWebhookBatch{}, fmt.Errorf("%w: payload or secret", ErrInvalidPayload)
	}
	if !verifyMetaSHA256(payload, signatureHeader, appSecret) {
		return MetaWebhookBatch{}, ErrInvalidMetaSignature
	}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var envelope metaWebhookEnvelope
	if err := dec.Decode(&envelope); err != nil {
		return MetaWebhookBatch{}, fmt.Errorf("%w: Meta JSON", ErrInvalidPayload)
	}
	if err := ensureEOF(dec); err != nil {
		return MetaWebhookBatch{}, err
	}
	if envelope.Object != "page" || len(envelope.Entries) == 0 {
		return MetaWebhookBatch{}, fmt.Errorf("%w: Meta object", ErrInvalidPayload)
	}
	signals := make([]MetaLeadSignal, 0)
	seen := map[string]string{}
	for _, entry := range envelope.Entries {
		for _, change := range entry.Changes {
			if change.Field != "leadgen" {
				continue
			}
			var value metaLeadgenValue
			valueDecoder := json.NewDecoder(bytes.NewReader(change.Value))
			valueDecoder.UseNumber()
			if err := valueDecoder.Decode(&value); err != nil {
				return MetaWebhookBatch{}, fmt.Errorf("%w: Meta leadgen value", ErrInvalidPayload)
			}
			if err := ensureEOF(valueDecoder); err != nil {
				return MetaWebhookBatch{}, err
			}
			pageID := strings.TrimSpace(value.PageID)
			if pageID == "" {
				pageID = strings.TrimSpace(entry.PageID)
			}
			created, err := value.CreatedTime.Int64()
			if err != nil || created <= 0 {
				return MetaWebhookBatch{}, fmt.Errorf("%w: Meta created_time", ErrInvalidPayload)
			}
			canonical, err := canonicalJSON(change.Value)
			if err != nil {
				return MetaWebhookBatch{}, err
			}
			signal := MetaLeadSignal{
				LeadgenID: strings.TrimSpace(value.LeadgenID), FormID: strings.TrimSpace(value.FormID), PageID: pageID,
				AdID: strings.TrimSpace(value.AdID), AdGroupID: strings.TrimSpace(value.AdGroupID),
				OccurredAt: time.Unix(created, 0).UTC(), ValueSHA256: payloadHash(canonical),
			}
			if prior, exists := seen[signal.LeadgenID]; exists {
				if prior != signal.ValueSHA256 {
					return MetaWebhookBatch{}, ErrDivergentDuplicate
				}
				continue
			}
			seen[signal.LeadgenID] = signal.ValueSHA256
			signals = append(signals, signal)
		}
	}
	batch := MetaWebhookBatch{TenantID: tenantID, OrganizationID: organizationID, ReceivedAt: receivedAt.UTC(), Payload: append([]byte(nil), payload...), PayloadSHA256: payloadHash(payload), Signals: signals}
	return batch, batch.Validate()
}

func verifyMetaSHA256(payload []byte, header, secret string) bool {
	const prefix = "sha256="
	if !strings.HasPrefix(header, prefix) || len(header) != len(prefix)+sha256.Size*2 || secret == "" {
		return false
	}
	provided, err := hex.DecodeString(header[len(prefix):])
	if err != nil || len(provided) != sha256.Size {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	return hmac.Equal(provided, mac.Sum(nil))
}

func canonicalJSON(raw []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var value any
	if err := dec.Decode(&value); err != nil {
		return nil, fmt.Errorf("%w: Meta canonical JSON", ErrInvalidPayload)
	}
	if err := ensureEOF(dec); err != nil {
		return nil, err
	}
	canonical, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%w: Meta canonical JSON", ErrInvalidPayload)
	}
	return canonical, nil
}

type MetaLeadWebhookHandler struct {
	Store          MetaWebhookStore
	TenantID       string
	OrganizationID string
	AppSecret      string
	VerifyToken    string
	Clock          func() time.Time
	MaxBytes       int64
}

func (h MetaLeadWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.verifySubscription(w, r)
	case http.MethodPost:
		h.receive(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h MetaLeadWebhookHandler) verifySubscription(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	got := r.URL.Query().Get("hub.verify_token")
	if r.URL.Query().Get("hub.mode") != "subscribe" || challenge == "" || !secretEqual(got, h.VerifyToken) {
		http.Error(w, "verification failed", http.StatusForbidden)
		return
	}
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, challenge)
}

func (h MetaLeadWebhookHandler) receive(w http.ResponseWriter, r *http.Request) {
	if h.Store == nil || strings.TrimSpace(h.TenantID) == "" || strings.TrimSpace(h.OrganizationID) == "" || h.AppSecret == "" {
		http.Error(w, "webhook unavailable", http.StatusServiceUnavailable)
		return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("content-type"))
	if err != nil || mediaType != "application/json" {
		http.Error(w, "content-type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	limit := h.MaxBytes
	if limit <= 0 || limit > MaxPayloadBytes {
		limit = MaxPayloadBytes
	}
	payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, limit))
	if err != nil {
		http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
		return
	}
	now := time.Now().UTC()
	if h.Clock != nil {
		now = h.Clock().UTC()
	}
	batch, err := DecodeMetaLeadWebhook(payload, r.Header.Get("X-Hub-Signature-256"), h.AppSecret, h.TenantID, h.OrganizationID, now)
	if errors.Is(err, ErrInvalidMetaSignature) {
		http.Error(w, "signature verification failed", http.StatusForbidden)
		return
	}
	if errors.Is(err, ErrDivergentDuplicate) {
		http.Error(w, "divergent provider identity", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "invalid webhook payload", http.StatusBadRequest)
		return
	}
	if _, err = h.Store.RecordMetaWebhook(r.Context(), batch); errors.Is(err, ErrDivergentDuplicate) {
		http.Error(w, "divergent provider identity", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, "durable intake unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "EVENT_RECEIVED")
}

type MemoryMetaWebhookStore struct {
	mu      sync.Mutex
	values  map[string]string
	batches map[string]struct{}
	Fail    error
}

func NewMemoryMetaWebhookStore() *MemoryMetaWebhookStore {
	return &MemoryMetaWebhookStore{values: map[string]string{}, batches: map[string]struct{}{}}
}

func (s *MemoryMetaWebhookStore) RecordMetaWebhook(_ context.Context, batch MetaWebhookBatch) (MetaWebhookReceipt, error) {
	if s == nil {
		return MetaWebhookReceipt{}, errors.New("leadstream: nil Meta webhook store")
	}
	if err := batch.Validate(); err != nil {
		return MetaWebhookReceipt{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Fail != nil {
		return MetaWebhookReceipt{}, s.Fail
	}
	for _, signal := range batch.Signals {
		key := batch.TenantID + "\x00" + signal.LeadgenID
		if prior, ok := s.values[key]; ok && prior != signal.ValueSHA256 {
			return MetaWebhookReceipt{}, ErrDivergentDuplicate
		}
	}
	receipt := MetaWebhookReceipt{}
	for _, signal := range batch.Signals {
		key := batch.TenantID + "\x00" + signal.LeadgenID
		if _, ok := s.values[key]; ok {
			receipt.DuplicateSignals++
			continue
		}
		s.values[key] = signal.ValueSHA256
		receipt.NewSignals++
	}
	s.batches[batch.TenantID+"\x00"+batch.PayloadSHA256] = struct{}{}
	return receipt, nil
}

func (s *MemoryMetaWebhookStore) SignalCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.values)
}
````

### FILE: `internal/leadstream/meta_webhook_test.go`
```yaml
block_id: "GO-META-LEAD-WEBHOOK-SIGNAL:internal/leadstream/meta_webhook_test.go:v1"
operation: CREATE
provenance: ADAPTED
source: "Meta official webhook security regressions and official Lead Ads Postman payload; local Go tests"
license: "LicenseRef-Workspace-Owner AND LicenseRef-Meta-Platform-API-Only AND MIT"
sha256: "6ea0810589c09f882e29c8688637da51be52433c605350bc6af87a2b558c0412"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const validMetaLeadWebhook = `{"object":"page","entry":[{"id":"138124712925158","time":1630087927,"changes":[{"value":{"form_id":"172036835026334","leadgen_id":"233370042045345","created_time":1630087926,"page_id":"138124712925158"},"field":"leadgen"}]}]}`

func metaSignature(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func metaHandler(store MetaWebhookStore) MetaLeadWebhookHandler {
	return MetaLeadWebhookHandler{Store: store, TenantID: "tenant-1", OrganizationID: "store-1", AppSecret: "app-secret", VerifyToken: "verify-token", Clock: fixedTime}
}

func TestDecodeMetaLeadWebhookMatchesOfficialLeadgenEnvelope(t *testing.T) {
	batch, err := DecodeMetaLeadWebhook([]byte(validMetaLeadWebhook), metaSignature(validMetaLeadWebhook, "app-secret"), "app-secret", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	if batch.PayloadSHA256 != payloadHash([]byte(validMetaLeadWebhook)) || len(batch.Signals) != 1 {
		t.Fatalf("batch=%+v", batch)
	}
	signal := batch.Signals[0]
	if signal.LeadgenID != "233370042045345" || signal.FormID != "172036835026334" || signal.PageID != "138124712925158" || !signal.OccurredAt.Equal(time.Unix(1630087926, 0).UTC()) {
		t.Fatalf("signal=%+v", signal)
	}
}

func TestMetaSignatureUsesExactRawBodyAndSHA256(t *testing.T) {
	valid := metaSignature(validMetaLeadWebhook, "app-secret")
	if _, err := DecodeMetaLeadWebhook([]byte(validMetaLeadWebhook), valid, "app-secret", "tenant-1", "store-1", fixedTime()); err != nil {
		t.Fatal(err)
	}
	for _, signature := range []string{"", "sha1=bad", "sha256=bad", metaSignature(validMetaLeadWebhook+" ", "app-secret")} {
		if _, err := DecodeMetaLeadWebhook([]byte(validMetaLeadWebhook), signature, "app-secret", "tenant-1", "store-1", fixedTime()); !errors.Is(err, ErrInvalidMetaSignature) {
			t.Fatalf("signature %q: %v", signature, err)
		}
	}
}

func TestMetaWebhookSubscriptionVerification(t *testing.T) {
	handler := metaHandler(NewMemoryMetaWebhookStore())
	request := httptest.NewRequest(http.MethodGet, "/webhooks/meta-leads?hub.mode=subscribe&hub.challenge=challenge-123&hub.verify_token=verify-token", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Body.String() != "challenge-123" {
		t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
	}
	request = httptest.NewRequest(http.MethodGet, "/webhooks/meta-leads?hub.mode=subscribe&hub.challenge=challenge-123&hub.verify_token=wrong", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden || strings.Contains(response.Body.String(), "wrong") {
		t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
	}
}

func TestMetaWebhookAcknowledgesOnlyAfterDurableRecord(t *testing.T) {
	store := NewMemoryMetaWebhookStore()
	handler := metaHandler(store)
	call := func(payload, signature string) *httptest.ResponseRecorder {
		request := httptest.NewRequest(http.MethodPost, "/webhooks/meta-leads", strings.NewReader(payload))
		request.Header.Set("content-type", "application/json")
		request.Header.Set("X-Hub-Signature-256", signature)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		return response
	}
	if response := call(validMetaLeadWebhook, metaSignature(validMetaLeadWebhook, "app-secret")); response.Code != http.StatusOK || response.Body.String() != "EVENT_RECEIVED" || store.SignalCount() != 1 {
		t.Fatalf("status=%d body=%q count=%d", response.Code, response.Body.String(), store.SignalCount())
	}
	if response := call(validMetaLeadWebhook, metaSignature(validMetaLeadWebhook, "app-secret")); response.Code != http.StatusOK || store.SignalCount() != 1 {
		t.Fatalf("replay status=%d count=%d", response.Code, store.SignalCount())
	}
	store.Fail = errors.New("database unavailable")
	second := strings.Replace(validMetaLeadWebhook, "233370042045345", "233370042045346", 1)
	if response := call(second, metaSignature(second, "app-secret")); response.Code != http.StatusInternalServerError || store.SignalCount() != 1 {
		t.Fatalf("failure status=%d count=%d", response.Code, store.SignalCount())
	}
}

func TestMetaWebhookRejectsInvalidBoundaryWithoutStorage(t *testing.T) {
	for _, test := range []struct {
		name, payload, signature, contentType string
		maxBytes                              int64
		want                                  int
	}{
		{name: "signature", payload: validMetaLeadWebhook, signature: "sha256=bad", contentType: "application/json", want: http.StatusForbidden},
		{name: "content type", payload: validMetaLeadWebhook, signature: metaSignature(validMetaLeadWebhook, "app-secret"), contentType: "text/plain", want: http.StatusUnsupportedMediaType},
		{name: "malformed", payload: `{`, signature: metaSignature(`{`, "app-secret"), contentType: "application/json", want: http.StatusBadRequest},
		{name: "oversize", payload: validMetaLeadWebhook, signature: metaSignature(validMetaLeadWebhook, "app-secret"), contentType: "application/json", maxBytes: 8, want: http.StatusRequestEntityTooLarge},
		{name: "wrong object", payload: `{"object":"user","entry":[]}`, signature: metaSignature(`{"object":"user","entry":[]}`, "app-secret"), contentType: "application/json", want: http.StatusBadRequest},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := NewMemoryMetaWebhookStore()
			handler := metaHandler(store)
			handler.MaxBytes = test.maxBytes
			request := httptest.NewRequest(http.MethodPost, "/webhooks/meta-leads", strings.NewReader(test.payload))
			request.Header.Set("content-type", test.contentType)
			request.Header.Set("X-Hub-Signature-256", test.signature)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.want || store.SignalCount() != 0 {
				t.Fatalf("status=%d want=%d stored=%d", response.Code, test.want, store.SignalCount())
			}
		})
	}
}

func TestMetaWebhookBatchedEntriesAndDivergence(t *testing.T) {
	second := `{"object":"page","entry":[{"id":"138124712925158","time":1630087927,"changes":[{"value":{"form_id":"172036835026334","leadgen_id":"233370042045345","created_time":1630087926,"page_id":"138124712925158"},"field":"leadgen"}]},{"id":"138124712925158","time":1630087928,"changes":[{"value":{"form_id":"172036835026334","leadgen_id":"233370042045346","created_time":1630087927,"page_id":"138124712925158"},"field":"leadgen"}]}]}`
	batch, err := DecodeMetaLeadWebhook([]byte(second), metaSignature(second, "app-secret"), "app-secret", "tenant-1", "store-1", fixedTime())
	if err != nil || len(batch.Signals) != 2 {
		t.Fatalf("signals=%d err=%v payload=%s", len(batch.Signals), err, second)
	}
	store := NewMemoryMetaWebhookStore()
	if receipt, err := store.RecordMetaWebhook(t.Context(), batch); err != nil || receipt.NewSignals != 2 {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
	changed := strings.Replace(validMetaLeadWebhook, "172036835026334", "different-form", 1)
	changedBatch, err := DecodeMetaLeadWebhook([]byte(changed), metaSignature(changed, "app-secret"), "app-secret", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordMetaWebhook(t.Context(), changedBatch); !errors.Is(err, ErrDivergentDuplicate) {
		t.Fatalf("expected divergence, got %v", err)
	}
}

func TestMemoryMetaWebhookConcurrentReplay(t *testing.T) {
	batch, err := DecodeMetaLeadWebhook([]byte(validMetaLeadWebhook), metaSignature(validMetaLeadWebhook, "app-secret"), "app-secret", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	store := NewMemoryMetaWebhookStore()
	const workers = 24
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, recordErr := store.RecordMetaWebhook(t.Context(), batch)
			errs <- recordErr
		}()
	}
	wg.Wait()
	close(errs)
	for recordErr := range errs {
		if recordErr != nil {
			t.Fatal(recordErr)
		}
	}
	if store.SignalCount() != 1 {
		t.Fatalf("signals=%d", store.SignalCount())
	}
}
````

### FILE: `internal/platform/postgres/meta_lead_webhook.go`
```yaml
block_id: "GO-META-LEAD-WEBHOOK-SIGNAL:internal/platform/postgres/meta_lead_webhook.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL durable boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "e4291888a02c3ce3ba96cf59f3d29ea2f23854b92e27a33a930be379f4a9c19b"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MetaLeadWebhookStore struct{ pool *pgxpool.Pool }

func NewMetaLeadWebhookStore(pool *pgxpool.Pool) *MetaLeadWebhookStore {
	return &MetaLeadWebhookStore{pool: pool}
}

func (s *MetaLeadWebhookStore) RecordMetaWebhook(ctx context.Context, batch leadstream.MetaWebhookBatch) (leadstream.MetaWebhookReceipt, error) {
	if s == nil || s.pool == nil {
		return leadstream.MetaWebhookReceipt{}, errors.New("postgres Meta webhook: nil pool")
	}
	if err := batch.Validate(); err != nil {
		return leadstream.MetaWebhookReceipt{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return leadstream.MetaWebhookReceipt{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	batchResult, err := tx.Exec(ctx, `insert into integration.meta_lead_webhook_batch
(tenant_id,organization_id,payload_sha256,payload,received_at,signal_count)
values($1,$2,$3,$4,$5,$6) on conflict do nothing`, batch.TenantID, batch.OrganizationID, batch.PayloadSHA256, batch.Payload, batch.ReceivedAt, len(batch.Signals))
	if err != nil {
		return leadstream.MetaWebhookReceipt{}, err
	}
	if batchResult.RowsAffected() == 0 {
		var organizationID string
		if err = tx.QueryRow(ctx, `select organization_id from integration.meta_lead_webhook_batch where tenant_id=$1 and payload_sha256=$2`, batch.TenantID, batch.PayloadSHA256).Scan(&organizationID); err != nil {
			return leadstream.MetaWebhookReceipt{}, err
		}
		if organizationID != batch.OrganizationID {
			return leadstream.MetaWebhookReceipt{}, leadstream.ErrDivergentDuplicate
		}
	}
	receipt := leadstream.MetaWebhookReceipt{}
	for _, signal := range batch.Signals {
		result, execErr := tx.Exec(ctx, `insert into integration.meta_lead_webhook_signal
(tenant_id,leadgen_id,organization_id,batch_payload_sha256,value_sha256,form_id,page_id,ad_id,ad_group_id,occurred_at,state)
values($1,$2,$3,$4,$5,$6,$7,nullif($8,''),nullif($9,''),$10,'pending_retrieval') on conflict do nothing`,
			batch.TenantID, signal.LeadgenID, batch.OrganizationID, batch.PayloadSHA256, signal.ValueSHA256,
			signal.FormID, signal.PageID, signal.AdID, signal.AdGroupID, signal.OccurredAt)
		if execErr != nil {
			return leadstream.MetaWebhookReceipt{}, execErr
		}
		if result.RowsAffected() == 0 {
			var hash string
			if err = tx.QueryRow(ctx, `select value_sha256 from integration.meta_lead_webhook_signal where tenant_id=$1 and leadgen_id=$2`, batch.TenantID, signal.LeadgenID).Scan(&hash); err != nil {
				return leadstream.MetaWebhookReceipt{}, err
			}
			if hash != signal.ValueSHA256 {
				return leadstream.MetaWebhookReceipt{}, leadstream.ErrDivergentDuplicate
			}
			receipt.DuplicateSignals++
			continue
		}
		eventPayload, _ := json.Marshal(map[string]any{
			"provider": "meta_lead_ads", "leadgen_id": signal.LeadgenID, "form_id": signal.FormID,
			"page_id": signal.PageID, "organization_id": batch.OrganizationID, "state": "pending_retrieval",
			"source_payload_sha256": batch.PayloadSHA256,
		})
		eventID := leadstream.StableUUID(batch.TenantID, "meta_lead_ads", signal.LeadgenID, "signal")
		if _, err = tx.Exec(ctx, `insert into platform.outbox_event
(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
values($1,$2,'provider-lead-signal',$3,1,'meta-lead.signal-received',1,$4,$5)`, batch.TenantID, eventID, "meta_lead_ads:"+signal.LeadgenID, signal.OccurredAt, eventPayload); err != nil {
			return leadstream.MetaWebhookReceipt{}, err
		}
		receipt.NewSignals++
	}
	if err = tx.Commit(ctx); err != nil {
		return leadstream.MetaWebhookReceipt{}, err
	}
	return receipt, nil
}
````

### FILE: `internal/platform/postgres/meta_lead_webhook_integration_test.go`
```yaml
block_id: "GO-META-LEAD-WEBHOOK-SIGNAL:internal/platform/postgres/meta_lead_webhook_integration_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL integration regression"
license: "LicenseRef-Workspace-Owner"
sha256: "2dc87bd2640b06bbcdc7e8db617ee29ad9802fc78c4840763652543a252fc57f"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMetaLeadWebhookDurabilityReplayDivergenceAndTenantIsolation(t *testing.T) {
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

	tenantA := leadstream.StableUUID(t.Name(), "a", time.Now().UTC().Format(time.RFC3339Nano))
	tenantB := leadstream.StableUUID(t.Name(), "b", time.Now().UTC().Format(time.RFC3339Nano))
	for index, tenant := range []string{tenantA, tenantB} {
		code := "meta-" + tenant[:8]
		if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,$3,$3)`, tenant, code, "Meta Test "+string(rune('A'+index))); err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-1','store-1','Store 1','store')`, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-2','store-2','Store 2','store')`, tenantA); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null where tenant_id=any($1::uuid[]) and published_at is null`, []string{tenantA, tenantB})
	}()

	payload := `{"object":"page","entry":[{"id":"page-1","time":1630087927,"changes":[{"value":{"form_id":"form-1","leadgen_id":"lead-1","created_time":1630087926,"page_id":"page-1"},"field":"leadgen"}]}]}`
	sign := func(value string) string {
		mac := hmac.New(sha256.New, []byte("secret"))
		_, _ = mac.Write([]byte(value))
		return "sha256=" + hex.EncodeToString(mac.Sum(nil))
	}
	decode := func(tenant, value string) leadstream.MetaWebhookBatch {
		batch, decodeErr := leadstream.DecodeMetaLeadWebhook([]byte(value), sign(value), "secret", tenant, "store-1", time.Now().UTC())
		if decodeErr != nil {
			t.Fatal(decodeErr)
		}
		return batch
	}
	store := NewMetaLeadWebhookStore(pool)
	first, err := store.RecordMetaWebhook(ctx, decode(tenantA, payload))
	if err != nil || first.NewSignals != 1 {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	replay, err := store.RecordMetaWebhook(ctx, decode(tenantA, payload))
	if err != nil || replay.DuplicateSignals != 1 {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	other, err := store.RecordMetaWebhook(ctx, decode(tenantB, payload))
	if err != nil || other.NewSignals != 1 {
		t.Fatalf("tenant B=%+v err=%v", other, err)
	}
	changed := strings.Replace(payload, `"form_id":"form-1"`, `"form_id":"form-2"`, 1)
	if _, err = store.RecordMetaWebhook(ctx, decode(tenantA, changed)); !errors.Is(err, leadstream.ErrDivergentDuplicate) {
		t.Fatalf("expected divergence, got %v", err)
	}
	crossOrganization := decode(tenantA, payload)
	crossOrganization.OrganizationID = "store-2"
	if _, err = store.RecordMetaWebhook(ctx, crossOrganization); !errors.Is(err, leadstream.ErrDivergentDuplicate) {
		t.Fatalf("expected organization divergence, got %v", err)
	}
	var signals, batches, outbox int
	if err = pool.QueryRow(ctx, `select count(*) from integration.meta_lead_webhook_signal where tenant_id=$1 and leadgen_id='lead-1' and state='pending_retrieval'`, tenantA).Scan(&signals); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from integration.meta_lead_webhook_batch where tenant_id=$1`, tenantA).Scan(&batches); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='meta-lead.signal-received'`, tenantA).Scan(&outbox); err != nil {
		t.Fatal(err)
	}
	if signals != 1 || batches != 1 || outbox != 1 {
		t.Fatalf("signals=%d batches=%d outbox=%d", signals, batches, outbox)
	}
}
````

### FILE: `cmd/meta-lead-webhook/main.go`
```yaml
block_id: "GO-META-LEAD-WEBHOOK-SIGNAL:cmd/meta-lead-webhook/main.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local executable wiring"
license: "LicenseRef-Workspace-Owner"
sha256: "b16502cb272fa65088f7ae8b6114980348ea688d03db22c6efb9be53f1bc8e24"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	required := map[string]string{
		"DATABASE_URL": os.Getenv("DATABASE_URL"), "META_APP_SECRET": os.Getenv("META_APP_SECRET"),
		"META_VERIFY_TOKEN": os.Getenv("META_VERIFY_TOKEN"), "TENANT_ID": os.Getenv("TENANT_ID"),
		"ORGANIZATION_ID": os.Getenv("ORGANIZATION_ID"),
	}
	for name, value := range required {
		if value == "" {
			slog.Error("required configuration is missing", "name", name)
			os.Exit(2)
		}
	}
	pool, err := pgxpool.New(ctx, required["DATABASE_URL"])
	if err != nil {
		slog.Error("database configuration failed")
		os.Exit(2)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}
	handler := leadstream.MetaLeadWebhookHandler{
		Store: postgres.NewMetaLeadWebhookStore(pool), TenantID: required["TENANT_ID"], OrganizationID: required["ORGANIZATION_ID"],
		AppSecret: required["META_APP_SECRET"], VerifyToken: required["META_VERIFY_TOKEN"],
	}
	mux := http.NewServeMux()
	mux.Handle("/webhooks/meta-leads", handler)
	address := os.Getenv("LISTEN_ADDR")
	if address == "" {
		address = ":8082"
	}
	server := &http.Server{Addr: address, Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	slog.Info("Meta lead webhook listening", "address", address)
	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Meta lead webhook failed")
		os.Exit(1)
	}
}
````

### FILE: `db/migrations/0049_meta_lead_webhook_signal.up.sql`
```yaml
block_id: "GO-META-LEAD-WEBHOOK-SIGNAL:db/migrations/0049_meta_lead_webhook_signal.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL schema"
license: "LicenseRef-Workspace-Owner"
sha256: "611f0941f49796ae85f902c021092dffb5618256b136569d369c054fb9219f0e"
variables: []
secrets_allowed: false
```
````sql
begin;

create table integration.meta_lead_webhook_batch (
  tenant_id uuid not null references platform.tenant(tenant_id),
  organization_id text not null,
  payload_sha256 text not null check (payload_sha256 ~ '^[0-9a-f]{64}$'),
  payload bytea not null check (octet_length(payload) between 1 and 1048576),
  received_at timestamptz not null,
  signal_count integer not null check (signal_count > 0),
  primary key (tenant_id,payload_sha256),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id)
);

create table integration.meta_lead_webhook_signal (
  tenant_id uuid not null references platform.tenant(tenant_id),
  leadgen_id text not null check (length(leadgen_id) between 1 and 256),
  organization_id text not null,
  batch_payload_sha256 text not null,
  value_sha256 text not null check (value_sha256 ~ '^[0-9a-f]{64}$'),
  form_id text not null check (length(form_id) between 1 and 256),
  page_id text not null check (length(page_id) between 1 and 256),
  ad_id text,
  ad_group_id text,
  occurred_at timestamptz not null,
  state text not null check (state='pending_retrieval'),
  primary key (tenant_id,leadgen_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,batch_payload_sha256) references integration.meta_lead_webhook_batch(tenant_id,payload_sha256)
);

create index meta_lead_webhook_pending_idx on integration.meta_lead_webhook_signal
  (tenant_id,state,occurred_at,leadgen_id);

create function integration.reject_meta_lead_webhook_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='Meta lead webhook evidence is append-only';
end;
$function$;

create trigger meta_lead_webhook_batch_immutable before update or delete on integration.meta_lead_webhook_batch
for each row execute function integration.reject_meta_lead_webhook_mutation();

create trigger meta_lead_webhook_signal_immutable before update or delete on integration.meta_lead_webhook_signal
for each row execute function integration.reject_meta_lead_webhook_mutation();

commit;
````

### FILE: `db/migrations/0049_meta_lead_webhook_signal.down.sql`
```yaml
block_id: "GO-META-LEAD-WEBHOOK-SIGNAL:db/migrations/0049_meta_lead_webhook_signal.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL rollback"
license: "LicenseRef-Workspace-Owner"
sha256: "d5847b6e6bcbf4a92958037dd1c6a0e610c0285f13e6deb1162ecf1270aae3bd"
variables: []
secrets_allowed: false
```
````sql
begin;
drop trigger if exists meta_lead_webhook_signal_immutable on integration.meta_lead_webhook_signal;
drop trigger if exists meta_lead_webhook_batch_immutable on integration.meta_lead_webhook_batch;
drop table if exists integration.meta_lead_webhook_signal;
drop table if exists integration.meta_lead_webhook_batch;
drop function if exists integration.reject_meta_lead_webhook_mutation();
commit;
````

### FILE: `db/tests/0049_meta_lead_webhook_signal.test.sql`
```yaml
block_id: "GO-META-LEAD-WEBHOOK-SIGNAL:db/tests/0049_meta_lead_webhook_signal.test.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local SQL regression"
license: "LicenseRef-Workspace-Owner"
sha256: "e78e910c73ad48668dc4fdd09d792ff1021063b46a7073b72f8c423d12002b06"
variables: []
secrets_allowed: false
```
````sql
begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('49494949-4949-4949-8949-494949494949','meta-webhook-test','Meta Webhook Test','Meta Webhook Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values('49494949-4949-4949-8949-494949494949','store-1','store-1','Store 1','store');

insert into integration.meta_lead_webhook_batch
(tenant_id,organization_id,payload_sha256,payload,received_at,signal_count)
values('49494949-4949-4949-8949-494949494949','store-1',repeat('a',64),convert_to('{"object":"page"}','UTF8'),clock_timestamp(),1);

insert into integration.meta_lead_webhook_signal
(tenant_id,leadgen_id,organization_id,batch_payload_sha256,value_sha256,form_id,page_id,occurred_at,state)
values('49494949-4949-4949-8949-494949494949','lead-1','store-1',repeat('a',64),repeat('b',64),'form-1','page-1',clock_timestamp(),'pending_retrieval');

do $test$
begin
  begin
    update integration.meta_lead_webhook_signal set form_id='changed'
    where tenant_id='49494949-4949-4949-8949-494949494949' and leadgen_id='lead-1';
    raise exception 'signal update unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  begin
    delete from integration.meta_lead_webhook_batch
    where tenant_id='49494949-4949-4949-8949-494949494949' and payload_sha256=repeat('a',64);
    raise exception 'batch delete unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
end;
$test$;

rollback;
````

### FILE: `third_party/meta-lead-webhook/official-source-lock.json`
```yaml
block_id: "GO-META-LEAD-WEBHOOK-SIGNAL:third_party/meta-lead-webhook/official-source-lock.json:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official source identities and limitations"
license: "LicenseRef-Workspace-Owner"
sha256: "cb9233db3c54d703aa2482047d22e61592229731540437d88518fe13c73b375d"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-official-meta-lead-webhook-source-lock/v1",
  "observed_at": "2026-09-04",
  "sources": [
    {
      "repository": "fbsamples/messenger-platform-samples",
      "commit": "354ee221ac1d081cc6105a1515a8468cc44f6710",
      "tree": "a711d923d3797d202cdeaf415d7de088da0ee424",
      "archive_bytes": 7674719,
      "archive_sha256": "65fc927d030d97ba7b8cc8f891fc5369a89b820e73326ee12703b9052af5202c",
      "license": "LicenseRef-Meta-Platform-API-Only",
      "license_path": "LICENSE",
      "license_bytes": 1045,
      "license_sha256": "b274a3f84cf720f8250784b4f8b8758d019626d4496f67a8870424a5924a9894",
      "focal_artifacts": [
        {"path": "messenger-api/messenger-api-and-webhooks/app.py", "bytes": 2765, "sha256": "8aa10d7b370d8e9eed9f05c4259b06b4ee48cf115cc98e2d8063964924608d42"},
        {"path": "messenger-api/messenger-api-and-webhooks/tests/test_webhook_security.py", "bytes": 8919, "sha256": "120e8a742b2eb1da66dd8b3029e045b464de3e04a3a4bd42aa449280fb95d4e2"}
      ],
      "claim": "raw-body X-Hub-Signature-256 HMAC verification, constant-time comparison, GET subscription verification and security regressions"
    },
    {
      "repository": "fbsamples/lead-ads-webhook-sample",
      "commit": "c0843165ad0f6acd82731e0c9da96a22b4162b35",
      "tree": "93c3556abd9bee22bb92db66673e75d67e0f2449",
      "archive_bytes": 7604183,
      "archive_sha256": "8da8485a8df58aa7ebf9e004e9ab3ebb44e019bc426be6b535ab0646785a67c8",
      "license": "MIT",
      "license_path": "LICENSE",
      "license_bytes": 1085,
      "license_sha256": "a80c79b8f80e824c8ea0612b7abadb9c92c48673cf0adc20b2a7f43e2b19faa0",
      "focal_artifacts": [
        {"path": "prjFBLeadAds/Controllers/WebhooksController.cs", "bytes": 5516, "sha256": "12dcdf4561318eb39d2e2ff1ffc0cf7620627d35a2ff7f5815e54b06c911d26b"},
        {"path": "postman/FB Lead Ads (Part 1 - The Webhook).postman_collection.json", "bytes": 10686, "sha256": "3d2a5516949fcab8203af1d9cb65b7596fd3034504635b139892ae304dde10a4"}
      ],
      "claim": "official page/entry/changes/leadgen notification envelope and leadgen_id retrieval handoff"
    }
  ],
  "limitations": [
    "Both repositories are official samples, not production certification.",
    "lead-ads-webhook-sample is archived and its controller is not adopted as a security boundary.",
    "messenger-platform-samples code may be used only with Meta/Facebook platform APIs.",
    "Live permissions, Page/form ownership, subscription, delivery and Graph API retrieval remain project gates."
  ]
}
````

## 6. Configuration surface

- `DATABASE_URL`, `META_APP_SECRET`, `META_VERIFY_TOKEN`, `TENANT_ID` y `ORGANIZATION_ID` son obligatorios; secretos sólo desde el mecanismo seguro del target.
- `LISTEN_ADDR` es opcional y usa `:8082` únicamente como default local.
- El límite de body nunca supera 1 MiB; CDN/WAF puede imponer uno menor.
- Suscripción, Page/form y permisos se fijan en el registro de proyecto, no en el pack.

## 7. Dependency bill

- Go stdlib y `github.com/jackc/pgx/v5 v5.10.0`, ya fijado por el backend.
- PostgreSQL 18.6 y esquemas `platform`, `org`, `integration`/outbox de los packs previos.
- Código fuente Meta oficial fijado por commit/árbol/archive SHA; la licencia del sample Messenger limita el uso a APIs Meta/Facebook.
- `PYTHON-META-LEAD-RECONCILIATION-ADAPTER` para recuperar campos completos; este pack no duplica el SDK.

## 8. Apply order

Aplicar `0001`–`0048`, luego `0049`. Montar `/webhooks/meta-leads` detrás del ingress público protegido y conectar `meta-lead.signal-received` al lane de recuperación/reconciliación. No promover el signal como lead completo.

## 9. Verification

1. Materializar junto a `GO-OMNICHANNEL-LEAD-INGRESS`, `PYTHON-META-LEAD-RECONCILIATION-ADAPTER` y `GO-META-LEAD-EVIDENCE-IMPORT`.
2. Ejecutar las 49 migraciones en PostgreSQL limpio y `db/tests/0049_meta_lead_webhook_signal.test.sql`.
3. Ejecutar dos veces tests focales con `TEST_DATABASE_URL`, después `go test -count=1 ./...`, `go vet ./...` y `go build ./...`.
4. En el proyecto, probar challenge, firma válida/alterada, lote múltiple, replay/divergencia, error de commit y test lead oficial.

Los gates live enumerados en la sección siguiente no se sustituyen por estas pruebas locales.

### Production gates

- app/business/cuenta y permisos Lead Retrieval demostrados;
- Page/form ownership, suscripción `leadgen`, verify token y app secret por gestor de secretos;
- callback HTTPS detrás de CDN/WAF con límites y observabilidad sin PII;
- test lead oficial que produce signal, fetch SDK, import, promotion y journey real;
- política de consentimiento, contacto, retención, exportación/borrado y región aprobada;
- cuotas, backoff, gaps, reconciliación, alertas y rollback probados.

## 10. Reconstruction evidence

`reconstruction_evidence/GO_META_LEAD_WEBHOOK_SIGNAL_2026-09-04_V239.md`
