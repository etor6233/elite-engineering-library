package httpapi

// AUTHORED HTTP composition; authenticated tenant/actor and organization
// permission bind every command and recovery. Money/stock are derived by SQL.
import (
	"context"
	"net/http"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
)

type CommercialReleaseService interface {
	CommitCommercialRelease(context.Context, string, string, franchisejourney.CommitCommercialReleaseCommand) (franchisejourney.CommercialReleaseReceipt, bool, error)
	CommercialReleaseResult(context.Context, string, string, string, string) (franchisejourney.CommercialReleaseReceipt, error)
	ValidateCommercialRelease(context.Context, string, string, string) (franchisejourney.CurrentCommercialRelease, error)
}
type CommercialReleaseModule struct{ Service CommercialReleaseService }

func (m CommercialReleaseModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	principal := franchiseJourneyAPI{verifier: verifier}
	mux.HandleFunc("POST /v1/franchise/handovers/{id}/commercial-release", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			OrganizationID    string `json:"organization_id"`
			ObservationSHA256 string `json:"observation_sha256"`
		}
		if !decodeStrict(w, r, &input) {
			return
		}
		actor, ok := principal.protected(w, r, "handover:manage", input.OrganizationID)
		if !ok {
			return
		}
		v, replay, err := m.Service.CommitCommercialRelease(r.Context(), actor.TenantID, actor.Subject, franchisejourney.CommitCommercialReleaseCommand{OrganizationID: input.OrganizationID, HandoverID: r.PathValue("id"), ObservationSHA256: input.ObservationSHA256, IdempotencyKey: r.Header.Get("Idempotency-Key")})
		if initialHandoverError(w, err) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		if replay {
			w.Header().Set("Idempotency-Replayed", "true")
		}
		writeJSON(w, http.StatusCreated, v)
	})
	mux.HandleFunc("GET /v1/franchise/handovers/{id}/commercial-release-result", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", org)
		if !ok {
			return
		}
		v, err := m.Service.CommercialReleaseResult(r.Context(), actor.TenantID, org, r.PathValue("id"), r.Header.Get("Idempotency-Key"))
		if initialHandoverError(w, err) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, http.StatusOK, v)
	})
	mux.HandleFunc("GET /v1/franchise/handovers/{id}/commercial-release-current", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", org)
		if !ok {
			return
		}
		v, err := m.Service.ValidateCommercialRelease(r.Context(), actor.TenantID, org, r.PathValue("id"))
		if initialHandoverError(w, err) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, http.StatusOK, v)
	})
}
