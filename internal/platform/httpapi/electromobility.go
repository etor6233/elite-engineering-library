package httpapi

import (
	"crypto/sha256"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/order"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type ElectromobilityAPI struct {
	service  *electromobility.Service
	verifier identity.Verifier
}

type EnterpriseModule interface {
	Register(*http.ServeMux, identity.Verifier)
}

func NewEnterprise(orders *order.Service, verifier identity.Verifier, service *electromobility.Service, modules ...EnterpriseModule) http.Handler {
	api := &ElectromobilityAPI{service: service, verifier: verifier}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/public/{tenantCode}/models", api.listModels)
	mux.HandleFunc("POST /v1/public/{tenantCode}/{organizationCode}/leads", api.captureLead)
	mux.HandleFunc("POST /v1/catalog/models", api.createModel)
	var media http.Handler
	for _, module := range modules {
		if module == nil {
			continue
		}
		if owner, ok := module.(interface{ DeferredPublicMedia() http.Handler }); ok {
			if handler := owner.DeferredPublicMedia(); handler != nil {
				media = handler
			}
		}
		module.Register(mux, verifier)
	}
	mux.Handle("/", New(orders, verifier))
	handler := recoverMiddleware(mux)
	if media != nil {
		handler = catalogMediaGate{next: handler, media: media}
	}
	return handler
}

type catalogMediaGate struct {
	next  http.Handler
	media http.Handler
}

func (g catalogMediaGate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const prefix = "/v1/public/catalog/media/"
	rest := strings.TrimPrefix(r.URL.Path, prefix)
	if r.Method == http.MethodGet && rest != r.URL.Path && rest != "" && !strings.Contains(rest, "/") {
		g.media.ServeHTTP(w, r)
		return
	}
	g.next.ServeHTTP(w, r)
}
func (a *ElectromobilityAPI) listModels(w http.ResponseWriter, r *http.Request) {
	models, err := a.service.ListPublicModels(r.Context(), r.PathValue("tenantCode"))
	if err != nil {
		writeProblem(w, 400, "INVALID_TENANT", "tenant code is invalid")
		return
	}
	writeJSON(w, 200, map[string]any{"models": models})
}
func (a *ElectromobilityAPI) createModel(w http.ResponseWriter, r *http.Request) {
	principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return
	}
	if !principal.Allowed("catalog:write") {
		writeProblem(w, 403, "FORBIDDEN", "catalog:write permission is required")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	var input electromobility.Model
	if !decodeStrict(w, r, &input) {
		return
	}
	model, err := a.service.CreateModel(r.Context(), principal.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_MODEL", "model does not match the contract")
		return
	}
	writeJSON(w, 201, model)
}
func (a *ElectromobilityAPI) captureLead(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
	if err != nil {
		writeProblem(w, 400, "INVALID_BODY", "body is invalid or too large")
		return
	}
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if len(idempotencyKey) < 16 || len(idempotencyKey) > 128 {
		writeProblem(w, 400, "IDEMPOTENCY_KEY_REQUIRED", "Idempotency-Key must contain 16 to 128 characters")
		return
	}
	var input struct {
		ModelID        string          `json:"model_id"`
		SourceCode     string          `json:"source_code"`
		Contact        json.RawMessage `json:"contact"`
		ConsentGranted bool            `json:"consent_granted"`
	}
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil || !input.ConsentGranted {
		writeProblem(w, 400, "CONSENT_REQUIRED", "valid explicit consent is required")
		return
	}
	hash := sha256.Sum256(body)
	lead, replayed, err := a.service.CaptureLead(r.Context(), r.PathValue("tenantCode"), r.PathValue("organizationCode"), input.ModelID, input.SourceCode, input.Contact, hex.EncodeToString(hash[:]), idempotencyKey)
	if err != nil {
		if errors.Is(err, electromobility.ErrConflict) {
			writeProblem(w, 409, "IDEMPOTENCY_CONFLICT", "idempotency key is processing or belongs to another request")
			return
		}
		writeProblem(w, 400, "INVALID_LEAD", "lead does not match the contract")
		return
	}
	status := 202
	if replayed {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, map[string]string{"lead_id": lead.ID, "status": "accepted"})
}
func decodeStrict(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	if decoder.Decode(target) != nil {
		writeProblem(w, 400, "INVALID_BODY", "body does not match the contract")
		return false
	}
	return true
}
