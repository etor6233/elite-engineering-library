package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/platform/identity"
)

type FiscalModule struct {
	Service    *fiscal.Service
	Parameters *fiscal.ParameterRegistry
}

func (m FiscalModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := fiscalAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("POST /v1/fiscal/points-of-sale", api.configurePointOfSale)
	mux.HandleFunc("POST /v1/fiscal/invoices", api.requestInvoice)
	mux.HandleFunc("GET /v1/fiscal/invoices/{id}", api.getInvoice)
	if m.Parameters != nil {
		parameterAPI{registry: m.Parameters, verifier: verifier}.Register(mux)
	}
}

type fiscalAPI struct {
	service  *fiscal.Service
	verifier identity.Verifier
}

func (a fiscalAPI) authorize(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
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
func fiscalScope(w http.ResponseWriter, p identity.Principal, organization string) bool {
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return false
	}
	return true
}
func (a fiscalAPI) configurePointOfSale(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "fiscal:configure")
	if !ok {
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	var input fiscal.PointOfSale
	if !decodeStrict(w, r, &input) || !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.ConfigurePointOfSale(r.Context(), p.TenantID, input)
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a fiscalAPI) requestInvoice(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "fiscal:issue")
	if !ok {
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
	if err != nil {
		writeProblem(w, 400, "INVALID_BODY", "body is invalid or too large")
		return
	}
	var input fiscal.Invoice
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil {
		writeProblem(w, 400, "INVALID_BODY", "body does not match the contract")
		return
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		writeProblem(w, 400, "INVALID_BODY", "body must contain exactly one JSON value")
		return
	}
	if !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	key := r.Header.Get("Idempotency-Key")
	sum := sha256.Sum256(body)
	value, replayed, err := a.service.RequestInvoice(r.Context(), p.TenantID, key, hex.EncodeToString(sum[:]), input)
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	if replayed {
		w.Header().Set("Idempotency-Replayed", "true")
		writeJSON(w, 200, value)
		return
	}
	writeJSON(w, 202, value)
}
func (a fiscalAPI) getInvoice(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "fiscal:read")
	if !ok {
		return
	}
	organization := r.URL.Query().Get("organization_id")
	if !fiscalScope(w, p, organization) {
		return
	}
	value, err := a.service.GetInvoice(r.Context(), p.TenantID, organization, r.PathValue("id"))
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func writeFiscalResult(w http.ResponseWriter, err error) {
	if errors.Is(err, fiscal.ErrInvalid) {
		writeProblem(w, 400, "INVALID_FISCAL_COMMAND", "fiscal command does not match contract")
		return
	}
	if errors.Is(err, fiscal.ErrNotFound) {
		writeProblem(w, 404, "FISCAL_NOT_FOUND", "fiscal record was not found in this scope")
		return
	}
	if errors.Is(err, fiscal.ErrConflict) {
		writeProblem(w, 409, "FISCAL_CONFLICT", "fiscal state, source, sequence, idempotency or version conflict")
		return
	}
	writeProblem(w, 500, "INTERNAL", "request failed")
}
