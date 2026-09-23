package httpapi

import (
	"context"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type warehouseVerifier struct{ p identity.Principal }

func (v warehouseVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.p, nil
}

type warehouseReader struct {
	calls int
	scope inventorycontrol.WarehouseWorkspaceScope
	fail  bool
}

func (r *warehouseReader) WarehouseWorkspace(_ context.Context, s inventorycontrol.WarehouseWorkspaceScope, q inventorycontrol.WarehouseWorkspaceQuery) (inventorycontrol.WarehouseWorkspacePage, error) {
	r.calls++
	r.scope = s
	if r.fail {
		return inventorycontrol.WarehouseWorkspacePage{}, errors.New("private database failure")
	}
	return inventorycontrol.WarehouseWorkspacePage{Schema: "warehouse-workspace/v1", Dataset: q.Dataset, OrganizationID: s.Organization, Rows: []inventorycontrol.WarehouseWorkspaceRow{}}, nil
}
func TestWarehouseWorkspaceHTTP(t *testing.T) {
	cases := []struct {
		name, query, permission string
		want, calls             int
		fail                    bool
	}{{"read", "?organization_id=org-a", "inventory:read", 200, 1, false}, {"write_does_not_imply_read", "?organization_id=org-a", "inventory:write", 403, 0, false}, {"wrong_org", "?organization_id=other", "inventory:read", 403, 0, false}, {"duplicate_org", "?organization_id=org-a&organization_id=org-a", "inventory:read", 400, 0, false}, {"limit_bound", "?organization_id=org-a&limit=51", "inventory:read", 400, 0, false}, {"unknown_parameter", "?organization_id=org-a&tenant_id=other", "inventory:read", 400, 0, false}, {"cursor_invalid", "?organization_id=org-a&cursor=invalid", "inventory:read", 400, 0, false}, {"database_failure", "?organization_id=org-a", "inventory:read", 503, 1, true}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reader := &warehouseReader{fail: c.fail}
			mux := http.NewServeMux()
			p := identity.Principal{Subject: "actor-a", TenantID: "tenant-a", Permissions: map[string]struct{}{c.permission: {}}, Organizations: map[string]struct{}{"org-a": {}, "org-b": {}}}
			(WarehouseWorkspaceModule{Reader: reader}).Register(mux, warehouseVerifier{p})
			req := httptest.NewRequest("GET", "/v1/inventory/workspace/items"+c.query, nil)
			req.Header.Set("Authorization", "Bearer fixture")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != c.want || reader.calls != c.calls {
				t.Fatal(w.Code, w.Body.String(), reader.calls)
			}
			if c.calls > 0 && (reader.scope.Tenant != "tenant-a" || reader.scope.Organization != "org-a" || len(reader.scope.Organizations) != 2) {
				t.Fatal(reader.scope)
			}
			if strings.Contains(w.Body.String(), "private database") {
				t.Fatal("error leak")
			}
		})
	}
}

type warehouseRequestStore struct {
	calls                       int
	actor, tenant, org, request string
}

func (s *warehouseRequestStore) RequestWarehouseReservation(_ context.Context, tenant, actor string, c inventorycontrol.WarehouseReservationRequest) (inventorycontrol.WarehouseReservationReceipt, error) {
	s.calls++
	s.actor = actor
	s.tenant = tenant
	s.org = c.OrganizationID
	return inventorycontrol.WarehouseReservationReceipt{Schema: "warehouse-reservation-receipt/v1", RequestID: c.RequestID}, nil
}
func (s *warehouseRequestStore) ReadWarehouseReservationRequest(_ context.Context, tenant, org, actor, request string) (inventorycontrol.WarehouseReservationReceipt, error) {
	s.calls++
	s.actor = actor
	s.tenant = tenant
	s.org = org
	s.request = request
	return inventorycontrol.WarehouseReservationReceipt{}, inventorycontrol.ErrWarehouseRequestNotFound
}
func TestWarehouseReservationRequestHTTP(t *testing.T) {
	payload := `{"request_id":"request-a","organization_id":"org-a","bin_id":"bin-a","item_id":"item-a","demand_kind":"manual","demand_id":"demand-a","demand_line_id":"line-a","quantity":"1","cancellation_disallowed":false}`
	for _, tc := range []struct {
		name, method, path, permission, body string
		status, calls                        int
	}{{"create", "POST", "/v1/inventory/warehouse-workspace/reservations", "inventory:write", payload, 200, 1}, {"read_only_cannot_create", "POST", "/v1/inventory/warehouse-workspace/reservations", "inventory:read", payload, 403, 0}, {"malformed", "POST", "/v1/inventory/warehouse-workspace/reservations", "inventory:write", strings.Replace(payload, "\"1\"", "\"0\"", 1), 400, 0}, {"foreign_org", "POST", "/v1/inventory/warehouse-workspace/reservations", "inventory:write", strings.Replace(payload, "org-a", "org-b", 1), 403, 0}, {"unknown_fields", "POST", "/v1/inventory/warehouse-workspace/reservations", "inventory:write", strings.TrimSuffix(payload, "}") + `,"tenant":"other"}`, 400, 0}, {"get_scope", "GET", "/v1/inventory/warehouse-workspace/reservation-requests/request-a?organization_id=org-a", "inventory:read", "", 404, 1}} {
		t.Run(tc.name, func(t *testing.T) {
			store := &warehouseRequestStore{}
			mux := http.NewServeMux()
			(WarehouseWorkspaceModule{Reader: &warehouseReader{}, Reservations: store}).Register(mux, warehouseVerifier{identity.Principal{Subject: "actor-a", TenantID: "tenant-a", Permissions: map[string]struct{}{tc.permission: {}}, Organizations: map[string]struct{}{"org-a": {}}}})
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer fixture")
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tc.status || store.calls != tc.calls {
				t.Fatal(w.Code, w.Body.String(), store.calls)
			}
			if store.calls > 0 && (store.actor != "actor-a" || store.tenant != "tenant-a" || store.org != "org-a") {
				t.Fatal(store)
			}
		})
	}
}
