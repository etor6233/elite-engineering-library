package httpapi

// AUTHORED role/scoping and bounded HTTP glue. No user-supplied evidence receipt.
import (
	"bytes"
	"elite.local/enterprise/internal/approval"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type DocumentModule struct{ Store *postgres.Documents }

func documentJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func documentError(w http.ResponseWriter, e error) {
	status, code := 503, "UNAVAILABLE"
	switch {
	case errors.Is(e, postgres.ErrDocumentScope):
		status, code = 404, "NOT_FOUND"
	case errors.Is(e, doc.ErrContract):
		status, code = 400, "INVALID_CONTRACT"
	case errors.Is(e, doc.ErrSecurity):
		status, code = 422, "QUARANTINED"
	case errors.Is(e, postgres.ErrDocumentConflict), errors.Is(e, approval.ErrDuplicate), errors.Is(e, approval.ErrNotPending), errors.Is(e, approval.ErrSeparation):
		status, code = 409, "CONSULT_RECORDED_STATE"
	}
	documentJSON(w, status, map[string]string{"code": code})
}
func (m DocumentModule) principal(w http.ResponseWriter, r *http.Request, v identity.Verifier, permission string) (identity.Principal, bool) {
	values := r.Header.Values("Authorization")
	if v == nil || len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") || len(values[0]) > 16391 || len(strings.Fields(values[0])) != 2 {
		documentJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return identity.Principal{}, false
	}
	p, e := v.Verify(r.Context(), strings.TrimPrefix(values[0], "Bearer "))
	if e != nil {
		documentJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return p, false
	}
	allowed := m.Store != nil && (m.Store.Allowed(p, permission) || permission == "documents:read" && (m.Store.Allowed(p, "documents:write") || m.Store.Allowed(p, "documents:review") || m.Store.Allowed(p, "documents:process")))
	if !allowed {
		documentJSON(w, 403, map[string]string{"code": "FORBIDDEN"})
		return p, false
	}
	if r.URL.RawQuery != "" && !(r.Method == "GET" && r.URL.Path == "/v1/documents") {
		documentJSON(w, 400, map[string]string{"code": "INVALID_QUERY"})
		return p, false
	}
	return p, true
}
func documentBody(w http.ResponseWriter, r *http.Request, out any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return doc.ErrContract
	}
	b, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if e != nil {
		return doc.ErrContract
	}
	canonical, _, e := approval.CanonicalPayload(b)
	if e != nil {
		return doc.ErrContract
	}
	d := json.NewDecoder(bytes.NewReader(canonical))
	d.DisallowUnknownFields()
	if d.Decode(out) != nil {
		return doc.ErrContract
	}
	return nil
}
func (m DocumentModule) Register(mux *http.ServeMux, v identity.Verifier) {
	m.registerCatalogInbox(mux, v)
	mux.HandleFunc("GET /v1/documents/{id}/evidence/{part}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:read")
		if !ok {
			return
		}
		b, e := m.Store.EvidencePart(r.Context(), p, r.PathValue("id"), r.PathValue("part"))
		if e != nil {
			documentError(w, e)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=\"evidence.json\"")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Document-Evidence-SHA256", doc.Hash(b))
		_, _ = w.Write(b)
	})

	mux.HandleFunc("PUT /v1/documents/{id}/original", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:write")
		if !ok {
			return
		}
		if r.Header.Get("Content-Type") != "application/octet-stream" {
			documentError(w, doc.ErrContract)
			return
		}
		for _, h := range []string{"X-Document-Name", "X-Document-SHA256", "X-Document-Profile-SHA256"} {
			if len(r.Header.Values(h)) != 1 {
				documentError(w, doc.ErrContract)
				return
			}
		}
		for _, h := range []string{"X-Document-Class", "X-Document-Schema-Version"} {
			if len(r.Header.Values(h)) > 1 {
				documentError(w, doc.ErrContract)
				return
			}
		}
		b, e := io.ReadAll(http.MaxBytesReader(w, r.Body, doc.MaxBytes))
		if e != nil {
			documentError(w, doc.ErrContract)
			return
		}
		result, e := m.Store.Receive(r.Context(), p, r.PathValue("id"), r.Header.Get("X-Document-Profile-SHA256"), doc.Original{Name: r.Header.Get("X-Document-Name"), SHA256: r.Header.Get("X-Document-SHA256"), Bytes: b, ClassID: r.Header.Get("X-Document-Class"), SchemaVersion: r.Header.Get("X-Document-Schema-Version")})
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 201, result)
	})
	mux.HandleFunc("GET /v1/documents/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:read")
		if !ok {
			return
		}
		result, e := m.Store.Read(r.Context(), p, r.PathValue("id"))
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
	mux.HandleFunc("GET /v1/documents/{id}/original", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:read")
		if !ok {
			return
		}
		result, e := m.Store.Original(r.Context(), p, r.PathValue("id"))
		if e != nil {
			documentError(w, e)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\"original.bin\"")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Length", strconv.Itoa(len(result.Bytes)))
		w.Header().Set("X-Document-SHA256", result.SHA256)
		_, _ = w.Write(result.Bytes)
	})
	mux.HandleFunc("POST /v1/documents/{id}/process", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:process")
		if !ok {
			return
		}
		var b struct{}
		if documentBody(w, r, &b) != nil {
			documentError(w, doc.ErrContract)
			return
		}
		result, e := m.Store.Process(r.Context(), p, r.PathValue("id"))
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
	mux.HandleFunc("POST /v1/documents/{id}/review", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:write")
		if !ok {
			return
		}
		var b struct {
			EvidenceSHA string          `json:"evidence_sha256"`
			Fields      json.RawMessage `json:"fields"`
		}
		if documentBody(w, r, &b) != nil {
			documentError(w, doc.ErrContract)
			return
		}
		result, e := m.Store.SubmitFields(r.Context(), p, r.PathValue("id"), b.EvidenceSHA, b.Fields)
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
	mux.HandleFunc("POST /v1/documents/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:review")
		if !ok {
			return
		}
		var b struct {
			SHA      string `json:"payload_sha256"`
			Approved *bool  `json:"approved"`
			Reason   string `json:"reason"`
		}
		if documentBody(w, r, &b) != nil || b.Approved == nil {
			documentError(w, doc.ErrContract)
			return
		}
		result, e := m.Store.Decide(r.Context(), p, r.PathValue("id"), b.SHA, *b.Approved, b.Reason)
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
}
