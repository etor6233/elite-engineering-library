package whatsappbridge

// AUTHORED explicit operator API for the immutable campaign and original effects.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
	"strings"
	"time"
)

func (c *Campaigns) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if c == nil {
		return
	}
	route := func(method, path, permission string, action func(http.ResponseWriter, *http.Request, identity.Principal)) {
		mux.HandleFunc(method+" /v1/franchise/marketing/campaigns"+path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			headers := r.Header.Values("Authorization")
			if verifier == nil || len(headers) != 1 || !strings.HasPrefix(headers[0], "Bearer ") || len(headers[0]) > 16391 || len(strings.Fields(headers[0])) != 2 {
				notificationProblem(w, 401, "UNAUTHENTICATED")
				return
			}
			p, e := verifier.Verify(r.Context(), strings.TrimPrefix(headers[0], "Bearer "))
			if e != nil {
				notificationProblem(w, 401, "UNAUTHENTICATED")
				return
			}
			if !c.authorized(p, permission) {
				notificationProblem(w, 403, "FORBIDDEN")
				return
			}
			if (r.URL.RawQuery != "" && path != "/workspace") || r.PathValue("id") != "" && (!cr.ValidID(r.PathValue("id")) || len(r.PathValue("id")) < 16) {
				notificationProblem(w, 400, "INVALID_TARGET")
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
			defer cancel()
			action(w, r.WithContext(ctx), p)
		})
	}
	fail := func(w http.ResponseWriter, e error) bool {
		if e == nil {
			return false
		}
		notificationProblem(w, 409, "CAMPAIGN_UNAVAILABLE_OR_DIVERGENT")
		return true
	}
	route("POST", "/prepare", "marketing:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in CampaignRequest
		if decodeSchedule(w, r, &in) != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := c.Prepare(r.Context(), p, in)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "", "marketing:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in CampaignSnapshot
		if decodeSchedule(w, r, &in) != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := c.Create(r.Context(), p, in)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("GET", "/{id}", "marketing:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		out, e := c.Status(r.Context(), p, r.PathValue("id"))
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/{id}/resume", "marketing:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash string `json:"campaign_sha256"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := c.Resume(r.Context(), p, r.PathValue("id"), in.Hash)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/{id}/decision", "marketing:approve", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash    string `json:"campaign_sha256"`
			Approve *bool  `json:"approve"`
			Reason  string `json:"reason"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) || in.Approve == nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := c.Review(r.Context(), p, r.PathValue("id"), in.Hash, *in.Approve, in.Reason)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/{id}/stop", "marketing:cancel", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash   string `json:"campaign_sha256"`
			Reason string `json:"reason"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		if !fail(w, c.Stop(r.Context(), p, r.PathValue("id"), in.Hash, in.Reason)) {
			replyJSON(w, map[string]bool{"stopped": true})
		}
	})
	c.registerWorkspace(route)
}
