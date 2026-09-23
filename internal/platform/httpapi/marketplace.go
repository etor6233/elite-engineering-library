package httpapi

// AUTHORED bounded role HTTP adapter. All authority comes from the verifier.
import (
	"context"
	"io"
	"net/http"
	"time"

	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/marketplacebridge"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
)

type MarketplacePublicationService interface {
	Prepare(context.Context, identity.Principal, mb.PrepareRequest) (mb.Intent, error)
	Submit(context.Context, identity.Principal, mb.Intent) (string, bool, error)
	Decide(context.Context, identity.Principal, string, string, bool, string) (approval.State, error)
	Status(context.Context, identity.Principal, string) (postgres.MarketplaceStatus, error)
	Send(context.Context, identity.Principal, string) error
	Reconcile(context.Context, identity.Principal, string) error
}
type MarketplaceModule struct {
	Service                  MarketplacePublicationService
	TenantID, OrganizationID string
}

func (m MarketplaceModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	route := func(method, path, permission string, action func(http.ResponseWriter, *http.Request, identity.Principal)) {
		mux.HandleFunc(method+" /v1/admin/marketplace"+path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
			if e != nil {
				writeProblem(w, 401, "UNAUTHENTICATED", "valid bearer required")
				return
			}
			if p.TenantID != m.TenantID || !p.Allowed(permission) || !p.AllowedOrganization(m.OrganizationID) {
				writeProblem(w, 403, "FORBIDDEN", "scope required")
				return
			}
			if r.URL.RawQuery != "" || r.PathValue("id") != "" && !cr.ValidID(r.PathValue("id")) {
				writeProblem(w, 400, "INVALID_TARGET", "exact target required")
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
			defer cancel()
			action(w, r.WithContext(ctx), p)
		})
	}
	decode := func(w http.ResponseWriter, r *http.Request, v any) bool {
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "application/json required")
			return false
		}
		if e := http.NewResponseController(w).SetReadDeadline(time.Now().Add(5 * time.Second)); e != nil {
			writeProblem(w, 503, "READ_BOUNDARY", "bounded read unavailable")
			return false
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
		if e != nil || mb.Decode(raw, v) != nil {
			writeProblem(w, 400, "INVALID_BODY", "exact bounded JSON required")
			return false
		}
		return true
	}
	failed := func(w http.ResponseWriter, e error) bool {
		if e == nil {
			return false
		}
		writeProblem(w, 409, "MARKETPLACE_UNCONFIRMED", "consult immutable request and provider reconciliation")
		return true
	}
	route("POST", "/prepare", "marketplace:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in mb.PrepareRequest
		if !decode(w, r, &in) {
			return
		}
		out, e := m.Service.Prepare(r.Context(), p, in)
		if !failed(w, e) {
			writeJSON(w, 200, out)
		}
	})
	route("POST", "/requests", "marketplace:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in mb.Intent
		if !decode(w, r, &in) {
			return
		}
		hash, replay, e := m.Service.Submit(r.Context(), p, in)
		if !failed(w, e) {
			writeJSON(w, 200, map[string]any{"approval_id": in.ApprovalID, "request_sha256": hash, "replay": replay})
		}
	})
	route("POST", "/requests/{id}/decision", "marketplace:approve", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash    string `json:"request_sha256"`
			Approve *bool  `json:"approve"`
			Reason  string `json:"reason"`
		}
		if !decode(w, r, &in) {
			return
		}
		if in.Approve == nil || !cr.ValidSHA(in.Hash) {
			writeProblem(w, 400, "INVALID_DECISION", "explicit bound decision required")
			return
		}
		out, e := m.Service.Decide(r.Context(), p, r.PathValue("id"), in.Hash, *in.Approve, in.Reason)
		if !failed(w, e) {
			writeJSON(w, 200, map[string]any{"state": out})
		}
	})
	route("GET", "/requests/{id}", "marketplace:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		out, e := m.Service.Status(r.Context(), p, r.PathValue("id"))
		if !failed(w, e) {
			writeJSON(w, 200, out)
		}
	})
	for _, v := range []struct {
		name, permission string
		call             func(context.Context, identity.Principal, string) error
	}{{"send", "marketplace:send", m.Service.Send}, {"reconcile", "marketplace:reconcile", m.Service.Reconcile}} {
		route("POST", "/requests/{id}/"+v.name, v.permission, func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
			var empty struct{}
			if !decode(w, r, &empty) {
				return
			}
			if !failed(w, v.call(r.Context(), p, r.PathValue("id"))) {
				writeJSON(w, 200, map[string]string{"outcome": "confirmed"})
			}
		})
	}
}
