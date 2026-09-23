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
