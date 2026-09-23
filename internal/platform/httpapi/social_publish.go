package httpapi

import (
	"context"
	"io"
	"net/http"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/socialbridge"
)

type SocialPublishingService interface {
	Submit(context.Context, identity.Principal, socialbridge.Request) (string, bool, error)
	Decide(context.Context, identity.Principal, string, string, bool, string) (approval.State, error)
	Status(context.Context, identity.Principal, string) (postgres.SocialStatus, error)
	Reconcile(context.Context, identity.Principal, string) error
}

// NewSocialPublishing exposes reviewed requests and read-only reconciliation.
// Principal is supplied only by the existing OIDC verifier, never request JSON.
func NewSocialPublishing(service SocialPublishingService, verifier identity.Verifier) http.Handler {
	mux := http.NewServeMux()
	principal := func(w http.ResponseWriter, r *http.Request) (identity.Principal, bool) {
		p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if e != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
			return p, false
		}
		return p, true
	}
	decode := func(w http.ResponseWriter, r *http.Request, v any) bool {
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "application/json required")
			return false
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
		if e != nil || socialbridge.Decode(raw, v) != nil {
			writeProblem(w, 400, "INVALID_BODY", "exact contract required")
			return false
		}
		return true
	}
	mux.HandleFunc("POST /v1/social/requests", func(w http.ResponseWriter, r *http.Request) {
		p, ok := principal(w, r)
		if !ok {
			return
		}
		var v socialbridge.Request
		if !decode(w, r, &v) {
			return
		}
		hash, replay, e := service.Submit(r.Context(), p, v)
		if e != nil {
			writeProblem(w, 409, "SOCIAL_BINDING_REJECTED", "request is unauthorized or conflicts with its exact binding")
			return
		}
		code := 201
		if replay {
			code = 200
		}
		writeJSON(w, code, map[string]any{"approval_id": v.Intent.ApprovalID, "request_sha256": hash, "replayed": replay})
	})
	mux.HandleFunc("POST /v1/social/requests/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		p, ok := principal(w, r)
		if !ok {
			return
		}
		var v struct {
			RequestSHA256 string `json:"request_sha256"`
			Approve       *bool  `json:"approve"`
			Reason        string `json:"reason"`
		}
		if !decode(w, r, &v) {
			return
		}
		if v.Approve == nil {
			writeProblem(w, 400, "INVALID_BODY", "explicit decision required")
			return
		}
		state, e := service.Decide(r.Context(), p, r.PathValue("id"), v.RequestSHA256, *v.Approve, v.Reason)
		if e != nil {
			writeProblem(w, 409, "SOCIAL_DECISION_REJECTED", "decision is unauthorized or no longer applicable")
			return
		}
		writeJSON(w, 200, map[string]any{"state": state})
	})
	mux.HandleFunc("GET /v1/social/requests/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := principal(w, r)
		if !ok {
			return
		}
		v, e := service.Status(r.Context(), p, r.PathValue("id"))
		if e != nil {
			writeProblem(w, 404, "SOCIAL_REQUEST_UNAVAILABLE", "request is unavailable in this scope")
			return
		}
		writeJSON(w, 200, v)
	})
	mux.HandleFunc("POST /v1/social/requests/{id}/reconcile", func(w http.ResponseWriter, r *http.Request) {
		p, ok := principal(w, r)
		if !ok {
			return
		}
		if e := service.Reconcile(r.Context(), p, r.PathValue("id")); e != nil {
			writeProblem(w, 409, "SOCIAL_RECONCILIATION_UNCONFIRMED", "no further write permitted; provider observation or an active lease remains unresolved")
			return
		}
		writeJSON(w, 200, map[string]string{"state": "reconciled"})
	})
	registerSocialWorkspace(mux, service, verifier)
	return mux
}
