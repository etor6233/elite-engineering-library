package httpapi

import (
	"errors"
	"net/http"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type InventoryControlModule struct {
	Service   *inventorycontrol.Service
	Bulk      *inventorycontrol.BulkService
	Warehouse *inventorycontrol.WarehouseService
	Transfer  *inventorycontrol.BulkTransferService
	RequireDurableReservations bool
}

func (m InventoryControlModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := inventoryControlAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("GET /v1/inventory/atp", api.atp)
	mux.HandleFunc("POST /v1/inventory/reservations", api.reserve)
	mux.HandleFunc("POST /v1/inventory/reservations/{id}/release", api.release)
	mux.HandleFunc("POST /v1/inventory/transfers", api.createTransfer)
	mux.HandleFunc("POST /v1/inventory/transfers/{id}/transitions", api.transitionTransfer)
	if m.Bulk != nil {
		registerBulkInventory(mux, m.Bulk, verifier, m.RequireDurableReservations)
	}
	if m.Warehouse != nil {
		registerWarehouse(mux, m.Warehouse, verifier)
	}
	if m.Transfer != nil {
		registerBulkTransfer(mux, m.Transfer, verifier)
	}
}

type inventoryControlAPI struct {
	service  *inventorycontrol.Service
	verifier identity.Verifier
}

func (a inventoryControlAPI) authorize(w http.ResponseWriter, r *http.Request, permission string, requireJSON bool) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if requireJSON && r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}

func (a inventoryControlAPI) atp(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:read", false)
	if !ok {
		return
	}
	organization, variant := r.URL.Query().Get("organization_id"), r.URL.Query().Get("variant_id")
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	horizon, err := time.Parse(time.RFC3339, r.URL.Query().Get("horizon"))
	if err != nil {
		writeProblem(w, 400, "INVALID_ATP_QUERY", "horizon must be RFC3339")
		return
	}
	value, err := a.service.AvailableToPromise(r.Context(), p.TenantID, organization, variant, horizon)
	if err != nil {
		writeProblem(w, 500, "ATP_FAILED", "availability could not be calculated")
		return
	}
	writeJSON(w, 200, value)
}

func (a inventoryControlAPI) reserve(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input struct {
		inventorycontrol.Reservation
		StockVersion int64 `json:"stock_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.Reserve(r.Context(), p.TenantID, input.Reservation, input.StockVersion)
	if errors.Is(err, inventorycontrol.ErrConflict) {
		writeProblem(w, 409, "RESERVATION_CONFLICT", "stock is no longer available")
		return
	}
	if err != nil {
		writeProblem(w, 400, "INVALID_RESERVATION", "reservation does not match contract")
		return
	}
	writeJSON(w, 201, value)
}

func (a inventoryControlAPI) release(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeInventoryTransition(w, a.service.Release(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version))
}

func (a inventoryControlAPI) createTransfer(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input struct {
		inventorycontrol.Transfer
		StockVersions map[string]int64 `json:"stock_versions"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.FromOrganizationID) || !p.AllowedOrganization(input.ToOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for both organizations")
		return
	}
	value, err := a.service.CreateTransfer(r.Context(), p.TenantID, input.Transfer, input.StockVersions)
	if errors.Is(err, inventorycontrol.ErrConflict) {
		writeProblem(w, 409, "TRANSFER_CONFLICT", "one or more stock units are unavailable")
		return
	}
	if err != nil {
		writeProblem(w, 400, "INVALID_TRANSFER", "transfer does not match contract")
		return
	}
	writeJSON(w, 201, value)
}

func (a inventoryControlAPI) transitionTransfer(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
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
	writeInventoryTransition(w, a.service.TransitionTransfer(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version))
}

func writeInventoryTransition(w http.ResponseWriter, err error) {
	if errors.Is(err, inventorycontrol.ErrConflict) {
		writeProblem(w, 409, "INVENTORY_CONFLICT", "state, scope or version conflict")
		return
	}
	if err != nil {
		writeProblem(w, 500, "INTERNAL_ERROR", "inventory transition failed")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "accepted"})
}
