package httpapi

import (
	"encoding/json"
	"io"
	"net/http"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/platform/identity"
)

type parameterAPI struct {
	registry *fiscal.ParameterRegistry
	verifier identity.Verifier
}

func (a parameterAPI) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/fiscal/parameter-schedules", a.configureSchedule)
	mux.HandleFunc("POST /v1/fiscal/parameters/refresh", a.refresh)
	mux.HandleFunc("GET /v1/fiscal/parameters/{id}", a.getSnapshot)
	mux.HandleFunc("POST /v1/fiscal/parameters/{id}/decisions", a.decide)
}

func (a parameterAPI) principal(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	return p, true
}

func parameterJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return false
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	if decoder.Decode(target) != nil || decoder.Decode(&struct{}{}) != io.EOF {
		writeProblem(w, 400, "INVALID_BODY", "body must contain exactly one value matching the contract")
		return false
	}
	return true
}

func (a parameterAPI) configureSchedule(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "fiscal:parameters:configure")
	if !ok {
		return
	}
	var input fiscal.ParameterSchedule
	if !parameterJSON(w, r, &input) || !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	value, err := a.registry.ConfigureSchedule(r.Context(), p.TenantID, p.Subject, input)
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}

func (a parameterAPI) refresh(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "fiscal:parameters:refresh")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string  `json:"organization_id"`
		TaxpayerCUIT   string  `json:"taxpayer_cuit"`
		Kind           string  `json:"kind"`
		VoucherClass   *string `json:"voucher_class"`
	}
	if !parameterJSON(w, r, &input) || !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	value, replayed, err := a.registry.Refresh(r.Context(), p.TenantID, input.OrganizationID, input.TaxpayerCUIT, input.Kind, input.VoucherClass)
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	if replayed {
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, 201, value)
}

func (a parameterAPI) getSnapshot(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "fiscal:parameters:read")
	if !ok {
		return
	}
	organization := r.URL.Query().Get("organization_id")
	if !fiscalScope(w, p, organization) {
		return
	}
	value, err := a.registry.Snapshot(r.Context(), p.TenantID, organization, r.PathValue("id"))
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}

func (a parameterAPI) decide(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "fiscal:parameters:approve")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Approved       bool   `json:"approved"`
		Reason         string `json:"reason"`
	}
	if !parameterJSON(w, r, &input) || !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	if err := a.registry.Decide(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), p.Subject, input.Approved, input.Reason); err != nil {
		writeFiscalResult(w, err)
		return
	}
	w.WriteHeader(204)
}
