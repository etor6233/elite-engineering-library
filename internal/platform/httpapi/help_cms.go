// AUTHORED bounded transport; persisted article rules belong to helpcenter/helpcms.
package httpapi

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	cms "elite.local/enterprise/internal/helpcms"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type HelpCMSService interface {
	Execute(context.Context, identity.Principal, cms.Command) (cms.Receipt, error)
	Result(context.Context, identity.Principal, string, string) (cms.Receipt, error)
	Article(context.Context, identity.Principal, string, string, int64) (cms.Article, error)
	List(context.Context, identity.Principal, string, string, string, string) (cms.Page, error)
	History(context.Context, identity.Principal, string, string, int64) (cms.Page, error)
}
type HelpCMSModule struct{ Service HelpCMSService }

func helpCMSReply(w http.ResponseWriter, v any, e error) {
	if e != nil {
		switch {
		case errors.Is(e, cms.ErrInvalid):
			writeProblem(w, 400, "HELP_INVALID", "invalid article request")
		case errors.Is(e, cms.ErrNotFound):
			writeProblem(w, 404, "HELP_UNAVAILABLE", "article unavailable in this scope")
		case errors.Is(e, cms.ErrConflict):
			writeProblem(w, 409, "HELP_CONFLICT", "consult saved result and current version")
		default:
			writeProblem(w, 503, "HELP_UNCONFIRMED", "result unconfirmed; consult saved command")
		}
		return
	}
	writeJSON(w, 200, v)
}
func helpCMSQuery(r *http.Request, allowed ...string) (url.Values, bool) {
	if len(r.URL.RawQuery) > 2048 {
		return nil, false
	}
	q, e := url.ParseQuery(r.URL.RawQuery)
	if e != nil {
		return nil, false
	}
	if len(q["organization_id"]) != 1 || !cms.ID(q.Get("organization_id")) {
		return nil, false
	}
	for k, v := range q {
		found := k == "organization_id"
		for _, a := range allowed {
			found = found || k == a
		}
		if !found || len(v) != 1 {
			return nil, false
		}
	}
	return q, true
}
func helpCMSVersion(s string) (int64, bool) {
	if s == "" {
		return 0, true
	}
	v, e := strconv.ParseInt(s, 10, 64)
	return v, e == nil && v > 0 && strconv.FormatInt(v, 10) == s
}
func (m HelpCMSModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	auth := func(w http.ResponseWriter, r *http.Request, write bool) (identity.Principal, bool) {
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("Vary", "Authorization")
		p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if e != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "valid session required")
			return p, false
		}
		if !p.Allowed("help:write") && !p.Allowed("help:publish") && (write || !p.Allowed("help:read")) {
			writeProblem(w, 403, "FORBIDDEN", "help permission required")
			return p, false
		}
		return p, true
	}
	mux.HandleFunc("POST /v1/help/cms/commands", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, true)
		if !ok {
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "CONTENT_TYPE", "application/json required")
			return
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
		if e != nil {
			writeProblem(w, 413, "BODY_LIMIT", "article command exceeds32KiB")
			return
		}
		if _, _, e = approval.CanonicalPayload(raw); e != nil {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		var c cms.Command
		decoder := json.NewDecoder(bytes.NewReader(raw))
		decoder.DisallowUnknownFields()
		if decoder.Decode(&c) != nil || decoder.Decode(new(any)) != io.EOF || !c.Valid() {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		if !cms.Authorized(p, cms.Permission(c.Action), c.OrganizationID) {
			writeProblem(w, 403, "HELP_FORBIDDEN", "action or organization unavailable")
			return
		}
		v, e := m.Service.Execute(r.Context(), p, c)
		helpCMSReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/help/cms/commands/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, true)
		if !ok {
			return
		}
		q, ok := helpCMSQuery(r)
		if !ok || !cms.ID(r.PathValue("id")) {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		v, e := m.Service.Result(r.Context(), p, q.Get("organization_id"), r.PathValue("id"))
		helpCMSReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/help/cms/articles", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, false)
		if !ok {
			return
		}
		q, ok := helpCMSQuery(r, "locale", "q", "after")
		if !ok {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		v, e := m.Service.List(r.Context(), p, q.Get("organization_id"), q.Get("locale"), q.Get("q"), q.Get("after"))
		helpCMSReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/help/cms/articles/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, false)
		if !ok {
			return
		}
		q, ok := helpCMSQuery(r, "version")
		version, valid := helpCMSVersion(q.Get("version"))
		if !ok || !valid || !cms.ID(r.PathValue("id")) {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		v, e := m.Service.Article(r.Context(), p, q.Get("organization_id"), r.PathValue("id"), version)
		helpCMSReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/help/cms/articles/{id}/history", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, true)
		if !ok {
			return
		}
		q, ok := helpCMSQuery(r, "before")
		before, valid := helpCMSVersion(q.Get("before"))
		if !ok || !valid || !cms.ID(r.PathValue("id")) {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		v, e := m.Service.History(r.Context(), p, q.Get("organization_id"), r.PathValue("id"), before)
		helpCMSReply(w, v, e)
	})
}
