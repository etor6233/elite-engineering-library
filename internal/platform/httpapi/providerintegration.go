package httpapi

import (
	"errors"
	"io"
	"mime"
	"net/http"

	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/providerintegration"
)

type ProviderIntegrationModule struct{ Service *providerintegration.Service }

func (m ProviderIntegrationModule) Register(mux *http.ServeMux, _ identity.Verifier) {
	mux.HandleFunc("POST /v1/integrations/{provider}/connections/{connection}/webhooks", m.receive)
}

func (m ProviderIntegrationModule) receive(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 256<<10))
	if err != nil {
		writeProblem(w, 400, "INVALID_WEBHOOK", "webhook body is invalid or too large")
		return
	}
	replayed, err := m.Service.Receive(r.Context(), r.PathValue("provider"), r.PathValue("connection"), r.Header.Get("X-Elite-Webhook-Timestamp"), r.Header.Get("X-Elite-Webhook-Signature"), body)
	switch {
	case errors.Is(err, providerintegration.ErrUnauthenticated):
		writeProblem(w, 401, "WEBHOOK_UNAUTHENTICATED", "webhook signature is invalid")
	case errors.Is(err, providerintegration.ErrConnection):
		writeProblem(w, 404, "PROVIDER_CONNECTION_NOT_FOUND", "provider connection was not found")
	case errors.Is(err, providerintegration.ErrConflict):
		writeProblem(w, 409, "PROVIDER_EVENT_CONFLICT", "provider event identifier belongs to another body")
	case err != nil:
		writeProblem(w, 400, "INVALID_WEBHOOK", "webhook does not match the admitted contract")
	case replayed:
		w.Header().Set("Idempotency-Replayed", "true")
		writeJSON(w, 200, map[string]string{"status": "accepted"})
	default:
		writeJSON(w, 202, map[string]string{"status": "accepted"})
	}
}
