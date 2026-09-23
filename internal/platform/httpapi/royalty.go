package httpapi

import (
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/royalty"
	"errors"
	"net/http"
)

type RoyaltyModule struct {
	Service *royalty.Service
	// RequireDurableCommands contains legacy create endpoints without replay receipts.
	RequireDurableCommands bool
}

func (m RoyaltyModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := royaltyAPI{service: m.Service, verifier: verifier, requireDurableCommands: m.RequireDurableCommands}
	mux.HandleFunc("POST /v1/royalties/policies", api.createPolicy)
	mux.HandleFunc("POST /v1/royalties/accruals/from-payment", api.accruePayment)
	mux.HandleFunc("POST /v1/royalties/settlements", api.openSettlement)
	mux.HandleFunc("POST /v1/royalties/settlements/{id}/close", api.closeSettlement)
	mux.HandleFunc("POST /v1/royalties/settlements/{id}/reverse", api.reverseSettlement)
	mux.HandleFunc("POST /v1/royalties/settlements/{id}/reconciliations", api.reconcile)
}

type royaltyAPI struct {
	service  *royalty.Service
	verifier identity.Verifier
	requireDurableCommands bool
}

func (a royaltyAPI) auth(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}
func royaltyOrganization(w http.ResponseWriter, p identity.Principal, organization string) bool {
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return false
	}
	return true
}

func (a royaltyAPI) createPolicy(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:policy")
	if !ok {
		return
	}
	var input royalty.Policy
	if !decodeStrict(w, r, &input) {
		return
	}
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	if a.requireDurableCommands {
		writeProblem(w, 409, "DURABLE_COMMAND_REQUIRED", "use POST /v1/finance/commands with action royalty_policy and its explicit command contract; legacy payload is not forwarded")
		return
	}
	value, err := a.service.CreatePolicy(r.Context(), p.TenantID, input)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a royaltyAPI) accruePayment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:post")
	if !ok {
		return
	}
	var input royalty.PaymentEvent
	if !decodeStrict(w, r, &input) {
		return
	}
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.AccruePayment(r.Context(), p.TenantID, r.Header.Get("Idempotency-Key"), input)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a royaltyAPI) openSettlement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:settle")
	if !ok {
		return
	}
	var input royalty.Settlement
	if !decodeStrict(w, r, &input) {
		return
	}
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	if a.requireDurableCommands {
		writeProblem(w, 409, "DURABLE_COMMAND_REQUIRED", "use POST /v1/finance/commands with action royalty_open and its explicit command contract; legacy payload is not forwarded")
		return
	}
	value, err := a.service.OpenSettlement(r.Context(), p.TenantID, input)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a royaltyAPI) closeSettlement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:settle")
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
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.CloseSettlement(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.ExpectedVersion)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func (a royaltyAPI) reverseSettlement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:reverse")
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
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.ReverseSettlement(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.ExpectedVersion, input.Reason)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a royaltyAPI) reconcile(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:reconcile")
	if !ok {
		return
	}
	var input struct {
		OrganizationID    string `json:"organization_id"`
		ExternalReference string `json:"external_reference"`
		ActualMinorUnits  int64  `json:"actual_minor_units"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.Reconcile(r.Context(), p.TenantID, input.OrganizationID, p.Subject, royalty.Reconciliation{SettlementID: r.PathValue("id"), ExternalReference: input.ExternalReference, ActualMinorUnits: input.ActualMinorUnits})
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func writeRoyaltyResult(w http.ResponseWriter, err error) {
	if errors.Is(err, royalty.ErrInvalid) {
		writeProblem(w, 400, "INVALID_ROYALTY_COMMAND", "royalty command does not match contract")
		return
	}
	if errors.Is(err, royalty.ErrConflict) {
		writeProblem(w, 409, "ROYALTY_CONFLICT", "royalty state, scope, idempotency or version conflict")
		return
	}
	writeProblem(w, 500, "INTERNAL", "request failed")
}
