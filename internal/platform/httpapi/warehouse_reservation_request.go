package httpapi

import (
	"context"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
	"time"
)

func registerWarehouseReservationRequests(mux *http.ServeMux, store inventorycontrol.WarehouseReservationStore, verifier identity.Verifier) {
	handle := func(w http.ResponseWriter, r *http.Request) {
		permission := "inventory:read"
		if r.Method == "POST" {
			permission = "inventory:write"
		}
		p, ok := (inventoryControlAPI{verifier: verifier}).authorize(w, r, permission, r.Method == "POST")
		if !ok {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 7*time.Second)
		defer cancel()
		var v inventorycontrol.WarehouseReservationReceipt
		var e error
		if r.Method == "POST" {
			if r.URL.RawQuery != "" {
				writeProblem(w, 400, "INVALID_WAREHOUSE_REQUEST", "query is not accepted")
				return
			}
			var c inventorycontrol.WarehouseReservationRequest
			if !decodeStrict(w, r, &c) {
				return
			}
			if !p.AllowedOrganization(c.OrganizationID) {
				writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "organization permission required")
				return
			}
			if c.Validate() != nil {
				writeProblem(w, 400, "INVALID_WAREHOUSE_REQUEST", "invalid reservation request")
				return
			}
			v, e = store.RequestWarehouseReservation(ctx, p.TenantID, p.Subject, c)
		} else {
			q := r.URL.Query()
			if len(q) != 1 || len(q["organization_id"]) != 1 || !p.AllowedOrganization(q.Get("organization_id")) {
				writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "organization permission required")
				return
			}
			id := r.PathValue("request")
			if len(id) < 1 || len(id) > 128 {
				writeProblem(w, 400, "INVALID_WAREHOUSE_REQUEST", "invalid request")
				return
			}
			v, e = store.ReadWarehouseReservationRequest(ctx, p.TenantID, q.Get("organization_id"), p.Subject, id)
		}
		if errors.Is(e, inventorycontrol.ErrConflict) {
			writeProblem(w, 409, "WAREHOUSE_RESERVATION_CONFLICT", "reservation state or request differs")
			return
		}
		if errors.Is(e, inventorycontrol.ErrWarehouseRequestNotFound) {
			writeProblem(w, 404, "WAREHOUSE_REQUEST_NOT_FOUND", "request was not found")
			return
		}
		if errors.Is(e, inventorycontrol.ErrWarehouseWorkspaceQuery) {
			writeProblem(w, 400, "INVALID_WAREHOUSE_REQUEST", "invalid request")
			return
		}
		if e != nil {
			writeProblem(w, 503, "WAREHOUSE_RESERVATION_UNCERTAIN", "query this request before retrying")
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, 200, v)
	}
	mux.HandleFunc("POST /v1/inventory/warehouse-workspace/reservations", handle)
	mux.HandleFunc("GET /v1/inventory/warehouse-workspace/reservation-requests/{request}", handle)
}
