package httpapi

import (
	"errors"
	"net/http"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type bulkInventoryAPI struct {
	service  *inventorycontrol.BulkService
	verifier identity.Verifier
	requireDurableReservations bool
}

func registerBulkInventory(mux *http.ServeMux, service *inventorycontrol.BulkService, verifier identity.Verifier, requireDurable ...bool) {
	api := bulkInventoryAPI{service: service, verifier: verifier}
	api.requireDurableReservations = len(requireDurable) > 0 && requireDurable[0]
	mux.HandleFunc("POST /v1/inventory/bulk/items", api.createItem)
	mux.HandleFunc("POST /v1/inventory/bulk/item-units", api.configureUnitOfMeasure)
	mux.HandleFunc("POST /v1/inventory/bulk/uom-conversions", api.convertUnitOfMeasure)
	mux.HandleFunc("POST /v1/inventory/bulk/handling-unit-conversions", api.convertHandlingUnits)
	mux.HandleFunc("POST /v1/inventory/bulk/bins", api.createBin)
	mux.HandleFunc("POST /v1/inventory/bulk/bin-policies", api.configureBin)
	mux.HandleFunc("POST /v1/inventory/bulk/receipts", api.receive)
	mux.HandleFunc("GET /v1/inventory/bulk/availability", api.availability)
	mux.HandleFunc("POST /v1/inventory/bulk/reservations", api.reserve)
	mux.HandleFunc("POST /v1/inventory/bulk/reservations/{id}/release", api.release)
	mux.HandleFunc("POST /v1/inventory/bulk/movements", api.move)
	mux.HandleFunc("POST /v1/inventory/bulk/issues", api.issue)
}

func (a bulkInventoryAPI) convertHandlingUnits(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.PackagingConversionCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.ConvertHandlingUnits(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_HANDLING_UNIT_CONVERSION") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) authorize(w http.ResponseWriter, r *http.Request, permission string, requireJSON bool) (identity.Principal, bool) {
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

func writeBulkError(w http.ResponseWriter, err error, invalidCode string) bool {
	if errors.Is(err, inventorycontrol.ErrConflict) {
		writeProblem(w, 409, "BULK_INVENTORY_CONFLICT", "inventory state, identity, capacity or costing conflict")
		return true
	}
	if err != nil {
		writeProblem(w, 400, invalidCode, "request does not match the bulk inventory contract")
		return true
	}
	return false
}

func (a bulkInventoryAPI) createItem(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkItem
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.CreateItem(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_ITEM") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) createBin(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.WarehouseBin
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CreateBin(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_BIN") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) configureBin(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.ItemBinPolicy
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.ConfigureBin(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_ITEM_BIN_POLICY") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) receive(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkReceipt
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.Receive(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_RECEIPT") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) availability(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:read", false)
	if !ok {
		return
	}
	organization, item := r.URL.Query().Get("organization_id"), r.URL.Query().Get("item_id")
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	values, err := a.service.Availability(r.Context(), p.TenantID, organization, item)
	if err != nil {
		writeProblem(w, 400, "INVALID_BULK_AVAILABILITY_QUERY", "query does not match the bulk inventory contract")
		return
	}
	writeJSON(w, 200, map[string]any{"availability": values})
}

func (a bulkInventoryAPI) reserve(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkReservation
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	if a.requireDurableReservations {
		writeProblem(w, 409, "DURABLE_RESERVATION_REQUIRED", "use POST /v1/inventory/warehouse-workspace/reservations with request_id and its explicit contract; legacy payload is not forwarded")
		return
	}
	value, err := a.service.Reserve(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_RESERVATION") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) release(w http.ResponseWriter, r *http.Request) {
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
	err := a.service.Release(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version)
	if writeBulkError(w, err, "INVALID_BULK_RELEASE") {
		return
	}
	writeJSON(w, 200, map[string]string{"status": "released"})
}

func (a bulkInventoryAPI) move(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkMovement
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	err := a.service.Move(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_MOVEMENT") {
		return
	}
	writeJSON(w, 200, map[string]string{"status": "moved"})
}

func (a bulkInventoryAPI) issue(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkIssue
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.Issue(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_ISSUE") {
		return
	}
	writeJSON(w, 201, value)
}
