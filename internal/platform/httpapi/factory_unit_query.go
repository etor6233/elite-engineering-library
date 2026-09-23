package httpapi

import (
	"elite.local/enterprise/internal/enterprisequery"
	"errors"
	"net/http"
)

// AUTHORED exact query; authorization is identical to the existing factory list.
func (a enterpriseQueryAPI) factoryUnit(w http.ResponseWriter, r *http.Request) {
	principal, organization, ok := a.scope(w, r, "factory:read")
	if !ok {
		return
	}
	value, err := a.service.FactoryUnit(r.Context(), principal.TenantID, organization, r.PathValue("id"))
	if errors.Is(err, enterprisequery.ErrFactoryUnitNotFound) {
		writeProblem(w, 404, "NOT_FOUND", "factory unit not found")
		return
	}
	if errors.Is(err, enterprisequery.ErrInvalid) {
		writeProblem(w, 400, "INVALID_QUERY", "factory unit query is invalid")
		return
	}
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "factory unit query failed")
		return
	}
	writeJSON(w, 200, value)
}
