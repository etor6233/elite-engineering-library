package httpapi

import (
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
)

type OperationsModule struct{ Service *operations.Service }

func (m OperationsModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := operationsAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("POST /v1/procurement/purchase-orders", api.createPO)
	mux.HandleFunc("POST /v1/procurement/purchase-orders/{id}/transitions", api.transitionPO)
	mux.HandleFunc("POST /v1/factory/units", api.createUnit)
	mux.HandleFunc("POST /v1/factory/units/{id}/transitions", api.transitionUnit)
	mux.HandleFunc("POST /v1/inventory/stock", api.createStock)
	mux.HandleFunc("POST /v1/inventory/stock/{id}/transitions", api.transitionStock)
}

type operationsAPI struct {
	service  *operations.Service
	verifier identity.Verifier
}

func (a operationsAPI) principal(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
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
func (a operationsAPI) createPO(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "procurement:write")
	if !ok {
		return
	}
	var input operations.PurchaseOrder
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.DestinationOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CreatePurchaseOrder(r.Context(), p.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_PURCHASE_ORDER", "purchase order does not match contract")
		return
	}
	writeJSON(w, 201, value)
}
func (a operationsAPI) transitionPO(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "procurement:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeTransition(w, a.service.TransitionPurchaseOrder(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version))
}
func (a operationsAPI) createUnit(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "factory:write")
	if !ok {
		return
	}
	var input operations.ProductionUnit
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.RegisterProductionUnit(r.Context(), p.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_PRODUCTION_UNIT", "unit does not match contract")
		return
	}
	writeJSON(w, 201, value)
}
func (a operationsAPI) transitionUnit(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "factory:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeTransition(w, a.service.TransitionProductionUnit(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target))
}
func (a operationsAPI) createStock(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "inventory:write")
	if !ok {
		return
	}
	var input operations.StockUnit
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.ReceiveStockUnit(r.Context(), p.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_STOCK_UNIT", "stock unit does not match contract")
		return
	}
	writeJSON(w, 201, value)
}
func (a operationsAPI) transitionStock(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "inventory:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeTransition(w, a.service.TransitionStockUnit(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version))
}
func writeTransition(w http.ResponseWriter, err error) {
	if errors.Is(err, operations.ErrConflict) {
		writeProblem(w, 409, "TRANSITION_CONFLICT", "state or version conflict")
		return
	}
	if err != nil {
		writeProblem(w, 500, "INTERNAL_ERROR", "transition failed")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "accepted"})
}
