package httpapi

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"net/http"
	"net/url"
	"time"
)

func registerStoredValueMetrics(m StoredValueModule, mux *http.ServeMux, protected func(http.ResponseWriter, *http.Request, string, string) (identity.Principal, bool)) {
	mux.HandleFunc("GET /v1/franchise/stored-value/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("Vary", "Authorization")
		q, e := url.ParseQuery(r.URL.RawQuery)
		if e != nil || len(r.URL.RawQuery) > 512 || len(q) != 2 || len(q["organization_id"]) != 1 || len(q["program_id"]) != 1 || !sv.ValidID(q.Get("organization_id")) || !sv.ValidID(q.Get("program_id")) {
			writeProblem(w, 400, "METRIC_INVALID", "one organization and selected program required")
			return
		}
		p, ok := protected(w, r, "stored_value:read", q.Get("organization_id"))
		if !ok {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()
		v, e := m.Service.Metrics(ctx, p, q.Get("organization_id"), q.Get("program_id"))
		if storedValueError(w, e) {
			return
		}
		writeJSON(w, 200, v)
	})
}
