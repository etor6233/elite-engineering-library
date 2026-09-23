package httpapi

// AUTHORED bounded optional UI projection over the existing Page text owner.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/socialbridge"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SocialWorkspaceService interface {
	Workspace(context.Context, identity.Principal, string) (postgres.SocialWorkspacePage, error)
	WorkspaceStatus(context.Context, identity.Principal, string) (postgres.SocialWorkspaceDetail, error)
	PrepareDraft(context.Context, identity.Principal, postgres.SocialDraft) (postgres.SocialPrepared, error)
}

func registerSocialWorkspace(mux *http.ServeMux, service SocialPublishingService, verifier identity.Verifier) {
	ws, ok := service.(SocialWorkspaceService)
	if !ok {
		return
	}
	route := func(method, path, permission string, fn func(http.ResponseWriter, *http.Request, identity.Principal)) {
		mux.HandleFunc(method+" /v1/social/workspace"+path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			headers := r.Header.Values("Authorization")
			if len(headers) != 1 || len(headers[0]) > 16391 || !strings.HasPrefix(headers[0], "Bearer ") {
				writeProblem(w, 401, "UNAUTHENTICATED", "valid identity required")
				return
			}
			p, e := authenticate(r.Context(), headers[0], verifier)
			if e != nil {
				writeProblem(w, 401, "UNAUTHENTICATED", "valid identity required")
				return
			}
			if !p.Allowed(permission) {
				writeProblem(w, 403, "FORBIDDEN", "scope required")
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
			defer cancel()
			fn(w, r.WithContext(ctx), p)
		})
	}
	route("GET", "", "social:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		q, e := url.ParseQuery(r.URL.RawQuery)
		if e != nil {
			writeProblem(w, 400, "INVALID_QUERY", "cursor required")
			return
		}
		for k, v := range q {
			if k != "after" || len(v) != 1 {
				writeProblem(w, 400, "INVALID_QUERY", "cursor required")
				return
			}
		}
		after := q.Get("after")
		if after != "" && !cr.ValidID(after) {
			writeProblem(w, 400, "INVALID_QUERY", "cursor required")
			return
		}
		out, e := ws.Workspace(r.Context(), p, after)
		if e != nil {
			writeProblem(w, 409, "WORKSPACE_UNAVAILABLE", "scoped view unavailable")
			return
		}
		writeJSON(w, 200, out)
	})
	route("GET", "/{id}", "social:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		if r.URL.RawQuery != "" || !cr.ValidID(r.PathValue("id")) {
			writeProblem(w, 400, "INVALID_QUERY", "valid identifier required")
			return
		}
		out, e := ws.WorkspaceStatus(r.Context(), p, r.PathValue("id"))
		if e != nil {
			writeProblem(w, 404, "NOT_FOUND", "publication unavailable")
			return
		}
		writeJSON(w, 200, out)
	})
	route("POST", "/prepare", "social:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		if r.URL.RawQuery != "" {
			writeProblem(w, 400, "INVALID_QUERY", "query not supported")
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "JSON required")
			return
		}
		if e := http.NewResponseController(w).SetReadDeadline(time.Now().Add(5 * time.Second)); e != nil {
			writeProblem(w, 503, "READ_BOUNDARY", "bounded read required")
			return
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 24576))
		var in postgres.SocialDraft
		if e != nil || socialbridge.Decode(raw, &in) != nil {
			writeProblem(w, 400, "INVALID_BODY", "bounded exact draft required")
			return
		}
		out, e := ws.PrepareDraft(r.Context(), p, in)
		if e != nil {
			writeProblem(w, 409, "DRAFT_UNAVAILABLE", "review the draft and current permissions")
			return
		}
		writeJSON(w, 200, out)
	})
}
