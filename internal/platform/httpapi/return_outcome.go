package httpapi

import (
	"elite.local/enterprise/internal/franchisejourney"
	"errors"
	"net/http"
)

// Read-only; verified principal owns tenant, existing handover permission owns org scope.
func (a franchiseJourneyAPI) returnOutcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	q := r.URL.Query()
	if len(q) != 1 || len(q["organization_id"]) != 1 || q.Get("organization_id") == "" {
		writeProblem(w, 400, "INVALID_RETURN_QUERY", "one organization_id is required")
		return
	}
	organization := q.Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	v, e := a.service.ReturnOutcome(r.Context(), p.TenantID, organization, r.PathValue("id"))
	switch {
	case errors.Is(e, franchisejourney.ErrNotFound):
		writeProblem(w, 404, "RETURN_NOT_FOUND", "no return in this authorized scope")
	case errors.Is(e, franchisejourney.ErrInvalid):
		writeProblem(w, 400, "INVALID_RETURN_QUERY", "invalid return reference")
	case e != nil:
		writeProblem(w, 503, "RETURN_OUTCOME_UNAVAILABLE", "saved outcome cannot currently be confirmed")
	default:
		writeJSON(w, 200, v)
	}
}
