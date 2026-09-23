// AUTHORED bounded transport for original organization/agreement owners.
package httpapi

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	nr "elite.local/enterprise/internal/networkrole"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type NetworkRoleService interface {
	Execute(context.Context, identity.Principal, nr.Command) (nr.Receipt, error)
	Result(context.Context, identity.Principal, string, string) (nr.Receipt, error)
	Entity(context.Context, identity.Principal, string, string, string) (nr.Entity, error)
}
type NetworkRoleModule struct{ Service NetworkRoleService }

type networkRoleOptions interface {
 Options(context.Context, identity.Principal, string) ([]nr.Option,error)
}

func networkReply(w http.ResponseWriter, v any, e error) {
	if e != nil {
		switch {
		case errors.Is(e, nr.ErrInvalid):
			writeProblem(w, 400, "NETWORK_INVALID", "invalid network command")
		case errors.Is(e, nr.ErrNotFound):
			writeProblem(w, 404, "NETWORK_UNAVAILABLE", "result unavailable in this scope")
		case errors.Is(e, nr.ErrConflict):
			writeProblem(w, 409, "NETWORK_CONFLICT", "consult saved result and current version")
		default:
			writeProblem(w, 503, "NETWORK_UNCONFIRMED", "result not confirmed; consult saved command")
		}
		return
	}
	writeJSON(w, 200, v)
}
func (m NetworkRoleModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	auth := func(w http.ResponseWriter, r *http.Request) (identity.Principal, bool) {
		p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if e != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "valid session required")
			return p, false
		}
		if !p.Allowed("network:admin") && !p.Allowed("franchise:write") {
			writeProblem(w, 403, "FORBIDDEN", "network or franchise permission required")
			return p, false
		}
		return p, true
	}
	mux.HandleFunc("GET /v1/franchise/network/options/{kind}",func(w http.ResponseWriter,r *http.Request){
	 p,ok:=auth(w,r);if !ok{return}
	 kind:=r.PathValue("kind")
	 if r.URL.RawQuery!=""||(kind!="organizations"&&kind!="entities"){networkReply(w,nil,nr.ErrInvalid);return}
	 source,ok:=m.Service.(networkRoleOptions);if !ok{networkReply(w,nil,nr.ErrNotFound);return}
	 value,e:=source.Options(r.Context(),p,kind);w.Header().Set("Cache-Control","no-store");networkReply(w,value,e)
	})
	mux.HandleFunc("POST /v1/franchise/network/commands", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "CONTENT_TYPE", "application/json required")
			return
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
		if e != nil {
			writeProblem(w, 413, "BODY_LIMIT", "network command exceeds32KiB")
			return
		}
		if _, _, e = approval.CanonicalPayload(raw); e != nil {
			networkReply(w, nil, nr.ErrInvalid)
			return
		}
		var c nr.Command
		decoder := json.NewDecoder(bytes.NewReader(raw))
		decoder.DisallowUnknownFields()
		if decoder.Decode(&c) != nil || decoder.Decode(new(any)) != io.EOF || !c.Valid() {
			networkReply(w, nil, nr.ErrInvalid)
			return
		}
		if !nr.Authorized(p, c.Action, c.ScopeOrganizationID) {
			writeProblem(w, 403, "NETWORK_FORBIDDEN", "action or organization unavailable")
			return
		}
		v, e := m.Service.Execute(r.Context(), p, c)
		networkReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/franchise/network/commands/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		q := r.URL.Query()
		// protectedGet omits an empty root scope. Result still requires original root authority.
		if (len(q) != 0 && (len(q) != 1 || len(q["scope_organization_id"]) != 1)) || !nr.ID(r.PathValue("id")) {
			networkReply(w, nil, nr.ErrInvalid)
			return
		}
		v, e := m.Service.Result(r.Context(), p, q.Get("scope_organization_id"), r.PathValue("id"))
		networkReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/franchise/network/entities/{kind}/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		q := r.URL.Query()
		if len(q) != 1 || len(q["scope_organization_id"]) != 1 || !nr.ID(r.PathValue("id")) {
			networkReply(w, nil, nr.ErrInvalid)
			return
		}
		v, e := m.Service.Entity(r.Context(), p, r.PathValue("kind"), q.Get("scope_organization_id"), r.PathValue("id"))
		networkReply(w, v, e)
	})
}
