package httpapi

import (
	"net/http"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type warehouseAPI struct {
	service  *inventorycontrol.WarehouseService
	verifier identity.Verifier
}

func registerWarehouse(mux *http.ServeMux, service *inventorycontrol.WarehouseService, verifier identity.Verifier) {
	api := warehouseAPI{service: service, verifier: verifier}
	mux.HandleFunc("POST /v1/inventory/warehouse/receipts", api.postReceipt)
	mux.HandleFunc("POST /v1/inventory/warehouse/picks", api.createPick)
	mux.HandleFunc("POST /v1/inventory/warehouse/activities/{id}/register", api.register)
	mux.HandleFunc("POST /v1/inventory/warehouse/picks/{id}/cancel", api.cancelPick)
	mux.HandleFunc("POST /v1/inventory/warehouse/put-aways/{id}/cancel", api.cancelPutAway)
	mux.HandleFunc("POST /v1/inventory/warehouse/replenishments", api.createReplenishment)
	mux.HandleFunc("POST /v1/inventory/warehouse/replenishments/{id}/cancel", api.cancelReplenishment)
	mux.HandleFunc("POST /v1/inventory/warehouse/cross-dock-policies", api.configureCrossDock)
	mux.HandleFunc("POST /v1/inventory/warehouse/sales-bindings", api.configureSalesBinding)
	mux.HandleFunc("POST /v1/inventory/warehouse/customer-shipments", api.postCustomerShipment)
}

func (a warehouseAPI) postCustomerShipment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.CustomerShipmentCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.PostCustomerShipment(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_CUSTOMER_SHIPMENT") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) configureSalesBinding(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.SalesWarehouseBinding
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.ConfigureSalesBinding(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_SALES_WAREHOUSE_BINDING") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) configureCrossDock(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.WarehouseCrossDockPolicy
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.ConfigureCrossDock(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_CROSS_DOCK_POLICY") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) createReplenishment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.WarehouseReplenishmentCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CreateReplenishment(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_REPLENISHMENT") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) cancelReplenishment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
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
	value, err := a.service.CancelReplenishment(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_REPLENISHMENT_CANCELLATION") {
		return
	}
	writeJSON(w, 200, value)
}

func (a warehouseAPI) authorize(w http.ResponseWriter, r *http.Request) (identity.Principal, bool) {
	return (bulkInventoryAPI{verifier: a.verifier}).authorize(w, r, "inventory:write", true)
}

func (a warehouseAPI) postReceipt(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.WarehouseReceiptCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.PostReceipt(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_RECEIPT") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) createPick(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.WarehousePickCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CreatePick(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_PICK") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) register(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input struct {
		OrganizationID string    `json:"organization_id"`
		Version        int64     `json:"version"`
		PostingDate    time.Time `json:"posting_date"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.Register(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version, input.PostingDate)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_REGISTRATION") {
		return
	}
	writeJSON(w, 200, value)
}

func (a warehouseAPI) cancelPick(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
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
	value, err := a.service.CancelPick(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_PICK_CANCELLATION") {
		return
	}
	writeJSON(w, 200, value)
}

func (a warehouseAPI) cancelPutAway(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
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
	value, err := a.service.CancelPutAway(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_PUT_AWAY_CANCELLATION") {
		return
	}
	writeJSON(w, 200, value)
}
