package httpapi

import (
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
	"net/url"
	"strconv"
)

func (m DocumentModule) registerCatalogInbox(mux *http.ServeMux, v identity.Verifier) {
	mux.HandleFunc("GET /v1/documents/classes", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:read")
		if !ok {
			return
		}
		result, e := m.Store.Catalog(p)
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
	mux.HandleFunc("GET /v1/documents", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:read")
		if !ok {
			return
		}
		q, e := url.ParseQuery(r.URL.RawQuery)
		if e != nil {
			documentError(w, doc.ErrContract)
			return
		}
		for k, v := range q {
			if (k != "limit" && k != "cursor") || len(v) != 1 || v[0] == "" {
				documentError(w, doc.ErrContract)
				return
			}
		}
		limit := 25
		if q.Get("limit") != "" {
			limit, e = strconv.Atoi(q.Get("limit"))
			if e != nil {
				documentError(w, doc.ErrContract)
				return
			}
		}
		result, e := m.Store.List(r.Context(), p, q.Get("cursor"), limit)
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
}
