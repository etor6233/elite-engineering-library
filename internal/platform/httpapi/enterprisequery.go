package httpapi

import (
	"net/http"
	"strconv"

	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/platform/identity"
)

type EnterpriseQueryModule struct {
	Service *enterprisequery.Service
	Metrics RoleMetricsReader
}

func (m EnterpriseQueryModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Metrics != nil {
		RoleMetricsModule{Service: m.Metrics}.Register(mux, verifier)
	}
	api := enterpriseQueryAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("GET /v1/admin/overview", api.adminOverview)
	mux.HandleFunc("GET /v1/admin/orders", api.adminOrders)
	mux.HandleFunc("GET /v1/admin/service-cases", api.adminCases)
	mux.HandleFunc("GET /v1/factory/units", api.factoryUnits)
	mux.HandleFunc("GET /v1/factory/units/{id}", api.factoryUnit)
	mux.HandleFunc("GET /v1/customer/orders", api.customerOrders)
	mux.HandleFunc("GET /v1/customer/service-cases", api.customerCases)
}

type enterpriseQueryAPI struct {
	service  *enterprisequery.Service
	verifier identity.Verifier
}

func (a enterpriseQueryAPI) scope(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, string, bool) {
	principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return principal, "", false
	}
	if !principal.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return principal, "", false
	}
	organization := r.URL.Query().Get("organization_id")
	if !principal.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return principal, "", false
	}
	return principal, organization, true
}

func pageInput(w http.ResponseWriter, r *http.Request) (int, string, bool) {
	limit := 25
	if raw := r.URL.Query().Get("limit"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 1 || value > 100 {
			writeProblem(w, 400, "INVALID_PAGE", "limit must be between 1 and 100")
			return 0, "", false
		}
		limit = value
	}
	return limit, r.URL.Query().Get("after"), true
}

func (a enterpriseQueryAPI) adminOverview(w http.ResponseWriter, r *http.Request) {
	principal, organization, ok := a.scope(w, r, "admin:read")
	if !ok {
		return
	}
	value, err := a.service.Overview(r.Context(), principal.TenantID, organization)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "overview query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a enterpriseQueryAPI) adminOrders(w http.ResponseWriter, r *http.Request) {
	a.orders(w, r, "admin:read", "")
}

func (a enterpriseQueryAPI) customerOrders(w http.ResponseWriter, r *http.Request) {
	a.orders(w, r, "customer:self", "subject")
}

func (a enterpriseQueryAPI) orders(w http.ResponseWriter, r *http.Request, permission, customerMode string) {
	principal, organization, ok := a.scope(w, r, permission)
	if !ok {
		return
	}
	limit, after, ok := pageInput(w, r)
	if !ok {
		return
	}
	customer := ""
	if customerMode != "" {
		customer = principal.Subject
	}
	value, err := a.service.Orders(r.Context(), principal.TenantID, organization, customer, limit, after)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "orders query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a enterpriseQueryAPI) factoryUnits(w http.ResponseWriter, r *http.Request) {
	principal, organization, ok := a.scope(w, r, "factory:read")
	if !ok {
		return
	}
	limit, after, ok := pageInput(w, r)
	if !ok {
		return
	}
	value, err := a.service.FactoryUnits(r.Context(), principal.TenantID, organization, limit, after)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "factory query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a enterpriseQueryAPI) adminCases(w http.ResponseWriter, r *http.Request) {
	a.cases(w, r, "admin:read", "")
}

func (a enterpriseQueryAPI) customerCases(w http.ResponseWriter, r *http.Request) {
	a.cases(w, r, "customer:self", "subject")
}

func (a enterpriseQueryAPI) cases(w http.ResponseWriter, r *http.Request, permission, customerMode string) {
	principal, organization, ok := a.scope(w, r, permission)
	if !ok {
		return
	}
	limit, after, ok := pageInput(w, r)
	if !ok {
		return
	}
	customer := ""
	if customerMode != "" {
		customer = principal.Subject
	}
	value, err := a.service.ServiceCases(r.Context(), principal.TenantID, organization, customer, limit, after)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "service query failed")
		return
	}
	writeJSON(w, 200, value)
}
