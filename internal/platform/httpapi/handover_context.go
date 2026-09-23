package httpapi

import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
)

type HandoverContextService interface {
	OperatorContext(context.Context, string, string, string) (franchisejourney.HandoverOperatorContext, error)
}
type HandoverContextModule struct{ Service HandoverContextService }

func (m HandoverContextModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	principal := franchiseJourneyAPI{verifier: verifier}
	mux.HandleFunc("GET /v1/franchise/orders/{id}/handover-context", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", org)
		if !ok {
			return
		}
		v, err := m.Service.OperatorContext(r.Context(), actor.TenantID, org, r.PathValue("id"))
		if initialHandoverError(w, err) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, 200, v)
	})
}
