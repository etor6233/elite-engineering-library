package httpapi

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	rm "elite.local/enterprise/internal/rolemetrics"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type RoleMetricsReader interface {
	Read(context.Context, identity.Principal, string, string) (rm.Snapshot, error)
}
type RoleMetricsModule struct{ Service RoleMetricsReader }

func (m RoleMetricsModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	mux.HandleFunc("GET /v1/reporting/operations/{kind}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("Vary", "Authorization")
		p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if e != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "valid session required")
			return
		}
		q, e := url.ParseQuery(r.URL.RawQuery)
		kind := r.PathValue("kind")
		if e != nil || len(r.URL.RawQuery) > 256 || len(q) != 1 || len(q["organization_id"]) != 1 || !rm.ID(q.Get("organization_id")) || rm.Permission(kind) == "" {
			writeProblem(w, 400, "METRIC_INVALID", "one organization and a supported metric kind required")
			return
		}
		org := q.Get("organization_id")
		if !rm.Authorized(p, kind, org) {
			writeProblem(w, 403, "FORBIDDEN", "metric outside session authority")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()
		v, e := m.Service.Read(ctx, p, kind, org)
		if e != nil {
			if errors.Is(e, rm.ErrCardinality) {
				writeProblem(w, 409, "METRIC_CARDINALITY", "reference report exceeds bounded result; refine reporting profile")
				return
			}
			writeProblem(w, 503, "METRIC_UNAVAILABLE", "measurement unavailable; no zero substituted")
			return
		}
		writeJSON(w, 200, v)
	})
}
