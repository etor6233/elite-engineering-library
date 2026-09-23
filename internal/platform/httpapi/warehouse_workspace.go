// AUTHORED workspace reads; POSTs stay in InventoryControlModule and existing owners.
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type WarehouseWorkspaceModule struct {
	Reader inventorycontrol.WarehouseWorkspaceReader
	Reservations inventorycontrol.WarehouseReservationStore
}

func (m WarehouseWorkspaceModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Reservations != nil { registerWarehouseReservationRequests(mux, m.Reservations, verifier) }
	mux.HandleFunc("GET /v1/inventory/workspace/{dataset}", func(w http.ResponseWriter, r *http.Request) {
		api := inventoryControlAPI{verifier: verifier}
		p, ok := api.authorize(w, r, "inventory:read", false)
		if !ok {
			return
		}
		params := r.URL.Query()
		allowed := map[string]bool{"organization_id": true, "search": true, "item_id": true, "parent_id": true, "id": true, "request_id": true, "cursor": true, "limit": true}
		for k, v := range params {
			if !allowed[k] || len(v) != 1 {
				writeProblem(w, 400, "INVALID_WAREHOUSE_QUERY", "invalid query")
				return
			}
		}
		org := params.Get("organization_id")
		if !p.AllowedOrganization(org) {
			writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "organization permission required")
			return
		}
		limit := 25
		if v := params.Get("limit"); v != "" {
			n, e := strconv.Atoi(v)
			if e != nil {
				writeProblem(w, 400, "INVALID_WAREHOUSE_QUERY", "invalid limit")
				return
			}
			limit = n
		}
		orgs := []string{}
		for id := range p.Organizations {
			orgs = append(orgs, id)
		}
		sort.Strings(orgs)
		_, all := p.Permissions["*"]
		scope := inventorycontrol.WarehouseWorkspaceScope{Tenant: p.TenantID, Organization: org, Organizations: orgs, AllOrganizations: all}
		q := inventorycontrol.WarehouseWorkspaceQuery{Dataset: r.PathValue("dataset"), Organization: org, Search: params.Get("search"), ItemID: params.Get("item_id"), ParentID: params.Get("parent_id"), ID: params.Get("id"), RequestID: params.Get("request_id"), Cursor: params.Get("cursor"), Limit: limit}
		if q.Validate(scope) != nil {
			writeProblem(w, 400, "INVALID_WAREHOUSE_QUERY", "invalid query or cursor scope")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		result, err := m.Reader.WarehouseWorkspace(ctx, scope, q)
		if errors.Is(err, inventorycontrol.ErrWarehouseWorkspaceQuery) {
			writeProblem(w, 400, "INVALID_WAREHOUSE_QUERY", "unsupported query")
			return
		}
		if err != nil {
			writeProblem(w, 503, "WAREHOUSE_UNAVAILABLE", "warehouse is temporarily unavailable")
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		writeJSON(w, 200, result)
	})
}
