package httpapi

// AUTHORED typed HTTP binding. Existing posting/reversal endpoints keep authority.
import (
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"io"
	"net/http"
)

func (a fxAPI) prepareJournal(w http.ResponseWriter, r *http.Request) {
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
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	if err != nil {
		writeProblem(w, 413, "BODY_TOO_LARGE", "bounded JSON required")
		return
	}
	var c accounting.FXJournalCommand
	if bcfx.DecodeExactJSON(raw, &c) != nil {
		writeProblem(w, 400, "INVALID_REQUEST", "exact FX journal request required")
		return
	}
	c.IdempotencyKey = r.Header.Get("Idempotency-Key")
	if !accountingOrg(w, p, c.OrganizationID) {
		return
	}
	result, replay, err := a.service.PrepareJournal(r.Context(), p.TenantID, p.Subject, c)
	if err != nil {
		fxProblem(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if replay {
		w.Header().Set("Idempotent-Replay", "true")
	}
	writeJSON(w, 201, result)
}
func (a fxAPI) journalResult(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "accounting:read")
	if !ok {
		return
	}
	q := r.URL.Query()
	if len(q) != 1 || len(q["organization_id"]) != 1 || len(r.Header.Values("Idempotency-Key")) != 1 {
		writeProblem(w, 400, "INVALID_QUERY", "exact organization and request key required")
		return
	}
	org := q.Get("organization_id")
	if !accountingOrg(w, p, org) {
		return
	}
	result, err := a.service.JournalResult(r.Context(), p.TenantID, org, p.Subject, r.Header.Get("Idempotency-Key"))
	if err != nil {
		fxProblem(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 200, result)
}
