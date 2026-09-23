package httpapi

// AUTHORED authenticated conversion and draft-binding routes.
// The existing accounting endpoints retain posting and reversal authority.
import (
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"io"
	"net/http"
)

type FXConversionModule struct{ Service *accounting.FXService }

func (m FXConversionModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	a := fxAPI{m.Service, verifier}
	mux.HandleFunc("POST /v1/accounting/fx/conversions", a.record)
	mux.HandleFunc("GET /v1/accounting/fx/conversions/result", a.result)
	mux.HandleFunc("POST /v1/accounting/fx/journals", a.prepareJournal)
	mux.HandleFunc("GET /v1/accounting/fx/journals/result", a.journalResult)
}

type fxAPI struct {
	service  *accounting.FXService
	verifier identity.Verifier
}

func (a fxAPI) principal(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, e := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if e != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "valid bearer required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", "permission required")
		return p, false
	}
	return p, true
}
func (a fxAPI) record(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "accounting:write")
	if !ok {
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "JSON required")
		return
	}
	if len(r.Header.Values("Idempotency-Key")) != 1 {
		writeProblem(w, 400, "INVALID_KEY", "one request key required")
		return
	}
	var c accounting.FXCommand
	raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	if e != nil {
		var limit *http.MaxBytesError
		if errors.As(e, &limit) {
			writeProblem(w, 413, "BODY_TOO_LARGE", "request is too large")
		} else {
			writeProblem(w, 400, "INVALID_REQUEST", "request is invalid")
		}
		return
	}
	if bcfx.DecodeExactJSON(raw, &c) != nil {
		writeProblem(w, 400, "INVALID_REQUEST", "request is invalid")
		return
	}
	c.IdempotencyKey = r.Header.Get("Idempotency-Key")
	if !accountingOrg(w, p, c.OrganizationID) {
		return
	}
	v, replay, e := a.service.Record(r.Context(), p.TenantID, p.Subject, c)
	if e != nil {
		fxProblem(w, e)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if replay {
		w.Header().Set("Idempotent-Replay", "true")
	}
	writeJSON(w, 201, v)
}
func (a fxAPI) result(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "accounting:read")
	if !ok {
		return
	}
	q := r.URL.Query()
	if len(q) != 1 || len(q["organization_id"]) != 1 || len(r.Header.Values("Idempotency-Key")) != 1 {
		writeProblem(w, 400, "INVALID_QUERY", "exact organization and request key required")
		return
	}
	organization := q.Get("organization_id")
	if !accountingOrg(w, p, organization) {
		return
	}
	v, e := a.service.Result(r.Context(), p.TenantID, organization, p.Subject, r.Header.Get("Idempotency-Key"))
	if e != nil {
		fxProblem(w, e)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 200, v)
}
func fxProblem(w http.ResponseWriter, e error) {
	if errors.Is(e, accounting.ErrFXNotFound) {
		writeProblem(w, 404, "NOT_FOUND", "conversion receipt not found")
		return
	}
	writeAccountingResult(w, e)
}
