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
