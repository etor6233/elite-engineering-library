package httpapi

import (
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
)

type AccountingModule struct{ Service *accounting.Service }

func (m AccountingModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := accountingAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("POST /v1/accounting/accounts", api.createAccount)
	mux.HandleFunc("POST /v1/accounting/periods", api.openPeriod)
	mux.HandleFunc("POST /v1/accounting/journals", api.createJournal)
	mux.HandleFunc("POST /v1/accounting/journals/{id}/post", api.postJournal)
	mux.HandleFunc("POST /v1/accounting/journals/{id}/reverse", api.reverseJournal)
	mux.HandleFunc("POST /v1/accounting/periods/{id}/close", api.closePeriod)
	mux.HandleFunc("GET /v1/accounting/trial-balance", api.trialBalance)
}

type accountingAPI struct {
	service  *accounting.Service
	verifier identity.Verifier
}

func (a accountingAPI) auth(w http.ResponseWriter, r *http.Request, permission string, jsonBody bool) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if jsonBody && r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}
func accountingOrg(w http.ResponseWriter, p identity.Principal, id string) bool {
	if !p.AllowedOrganization(id) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return false
	}
	return true
}
func (a accountingAPI) createAccount(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:manage", true)
	if !ok {
		return
	}
	var input accounting.Account
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.CreateAccount(r.Context(), p.TenantID, input)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a accountingAPI) openPeriod(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:manage", true)
	if !ok {
		return
	}
	var input accounting.Period
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.OpenPeriod(r.Context(), p.TenantID, input)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a accountingAPI) createJournal(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:write", true)
	if !ok {
		return
	}
	var input accounting.Journal
	if !decodeStrict(w, r, &input) {
		return
	}
	if !accountingOrg(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.CreateJournal(r.Context(), p.TenantID, input)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a accountingAPI) postJournal(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:post", true)
	if !ok {
		return
	}
	var input struct {
		OrganizationID  string `json:"organization_id"`
		ExpectedVersion int64  `json:"expected_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !accountingOrg(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.PostJournal(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), p.Subject, input.ExpectedVersion)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func (a accountingAPI) reverseJournal(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:reverse", true)
	if !ok {
		return
	}
	var input struct {
		OrganizationID  string `json:"organization_id"`
		ExpectedVersion int64  `json:"expected_version"`
		Reason          string `json:"reason"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !accountingOrg(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.ReverseJournal(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), p.Subject, input.ExpectedVersion, input.Reason)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a accountingAPI) closePeriod(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:close", true)
	if !ok {
		return
	}
	var input struct {
		ExpectedVersion int64 `json:"expected_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.ClosePeriod(r.Context(), p.TenantID, r.PathValue("id"), input.ExpectedVersion)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func (a accountingAPI) trialBalance(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:read", false)
	if !ok {
		return
	}
	organization := r.URL.Query().Get("organization_id")
	if !accountingOrg(w, p, organization) {
		return
	}
	value, err := a.service.TrialBalance(r.Context(), p.TenantID, organization, r.URL.Query().Get("period_id"))
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func writeAccountingResult(w http.ResponseWriter, err error) {
	if errors.Is(err, accounting.ErrInvalid) {
		writeProblem(w, 400, "INVALID_ACCOUNTING_COMMAND", "accounting command does not match contract")
		return
	}
	if errors.Is(err, accounting.ErrConflict) {
		writeProblem(w, 409, "ACCOUNTING_CONFLICT", "accounting state, scope, balance or version conflict")
		return
	}
	writeProblem(w, 500, "INTERNAL", "request failed")
}
