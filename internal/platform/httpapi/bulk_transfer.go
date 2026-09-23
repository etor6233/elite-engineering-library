package httpapi

import (
	"net/http"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type bulkTransferAPI struct {
	service  *inventorycontrol.BulkTransferService
	verifier identity.Verifier
}

func registerBulkTransfer(mux *http.ServeMux, service *inventorycontrol.BulkTransferService, verifier identity.Verifier) {
	api := bulkTransferAPI{service: service, verifier: verifier}
	mux.HandleFunc("POST /v1/inventory/bulk-transfers", api.create)
	mux.HandleFunc("POST /v1/inventory/bulk-transfers/{id}/ship", api.ship)
	mux.HandleFunc("POST /v1/inventory/bulk-transfers/{id}/receive", api.receive)
	mux.HandleFunc("POST /v1/inventory/bulk-transfers/{id}/cancel", api.cancel)
}

func (a bulkTransferAPI) authorize(w http.ResponseWriter, r *http.Request) (identity.Principal, bool) {
	return (bulkInventoryAPI{verifier: a.verifier}).authorize(w, r, "inventory:write", true)
}
func allowedTransfer(p identity.Principal, from, to string) bool {
	return from != "" && to != "" && from != to && p.AllowedOrganization(from) && p.AllowedOrganization(to)
}

func (a bulkTransferAPI) create(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.BulkTransferCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !allowedTransfer(p, input.FromOrganizationID, input.ToOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token must be authorized for both transfer organizations")
		return
	}
	value, err := a.service.Create(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_TRANSFER") {
		return
	}
	writeJSON(w, 201, value)
}

type bulkTransferAction struct {
	FromOrganizationID  string    `json:"from_organization_id"`
	ToOrganizationID    string    `json:"to_organization_id"`
	Version             int64     `json:"version"`
	RequestID           string    `json:"request_id,omitempty"`
	Quantity            string    `json:"quantity,omitempty"`
	WarehouseActivityID string    `json:"warehouse_activity_id,omitempty"`
	PostingDate         time.Time `json:"posting_date"`
}

func (a bulkTransferAPI) action(w http.ResponseWriter, r *http.Request, kind string) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input bulkTransferAction
	if !decodeStrict(w, r, &input) {
		return
	}
	if !allowedTransfer(p, input.FromOrganizationID, input.ToOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token must be authorized for both transfer organizations")
		return
	}
	var value inventorycontrol.BulkTransfer
	var err error
	switch kind {
	case "ship":
		value, err = a.service.Ship(r.Context(), p.TenantID, r.PathValue("id"), input.FromOrganizationID, input.ToOrganizationID, input.Version, inventorycontrol.BulkTransferPostingCommand{RequestID: input.RequestID, Quantity: input.Quantity, WarehouseActivityID: input.WarehouseActivityID, PostingDate: input.PostingDate})
	case "receive":
		value, err = a.service.Receive(r.Context(), p.TenantID, r.PathValue("id"), input.FromOrganizationID, input.ToOrganizationID, input.Version, inventorycontrol.BulkTransferPostingCommand{RequestID: input.RequestID, Quantity: input.Quantity, WarehouseActivityID: input.WarehouseActivityID, PostingDate: input.PostingDate})
	default:
		value, err = a.service.Cancel(r.Context(), p.TenantID, r.PathValue("id"), input.FromOrganizationID, input.ToOrganizationID, input.Version)
	}
	if writeBulkError(w, err, "INVALID_BULK_TRANSFER_TRANSITION") {
		return
	}
	writeJSON(w, 200, value)
}
func (a bulkTransferAPI) ship(w http.ResponseWriter, r *http.Request)    { a.action(w, r, "ship") }
func (a bulkTransferAPI) receive(w http.ResponseWriter, r *http.Request) { a.action(w, r, "receive") }
func (a bulkTransferAPI) cancel(w http.ResponseWriter, r *http.Request)  { a.action(w, r, "cancel") }
