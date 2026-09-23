package httpapi

// AUTHORED exact namespace over admitted finance owners and scoped projections.
import (
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/royalty"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

type FinanceWorkspaceModule struct{ Store *postgres.FinanceStore }

func (m FinanceWorkspaceModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	auth := func(w http.ResponseWriter, r *http.Request) (identity.Principal, bool) {
		p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if e != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "valid bearer required")
			return p, false
		}
		if !p.Allowed("accounting:read") && !postgres.FinanceRoyaltyAllowed(p) {
			writeProblem(w, 403, "FORBIDDEN", "finance permission required")
			return p, false
		}
		return p, true
	}
	mux.HandleFunc("GET /v1/finance/workspace", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		q := r.URL.Query()
		for key, values := range q {
			if len(values) != 1 || (key != "organization_id" && key != "view" && key != "cursor" && key != "limit") {
				writeProblem(w, 400, "INVALID_QUERY", "bounded query required")
				return
			}
		}
		org := q.Get("organization_id")
		if !accountingOrg(w, p, org) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		if q.Get("view") == "" {
			if len(q) != 1 {
				writeProblem(w, 400, "INVALID_QUERY", "organization required")
				return
			}
			out, e := m.Store.WorkspaceAccess(r.Context(), p, org)
			if e != nil {
				financeProblem(w, e)
				return
			}
			writeJSON(w, 200, out)
			return
		}
		limit := 25
		if q.Get("limit") != "" {
			var e error
			limit, e = strconv.Atoi(q.Get("limit"))
			if e != nil || limit < 1 || limit > 50 {
				writeProblem(w, 400, "INVALID_QUERY", "bounded page required")
				return
			}
		}
		out, e := m.Store.WorkspacePage(r.Context(), p, org, q.Get("view"), q.Get("cursor"), limit)
		if e != nil {
			financeProblem(w, e)
			return
		}
		writeJSON(w, 200, out)
	})
	mux.HandleFunc("GET /v1/finance/journals/{id}/lines", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		q := r.URL.Query()
		for k, v := range q {
			if len(v) != 1 || (k != "organization_id" && k != "cursor" && k != "limit") {
				writeProblem(w, 400, "INVALID_QUERY", "exact page required")
				return
			}
		}
		org := q.Get("organization_id")
		if !accountingOrg(w, p, org) {
			return
		}
		cursor := 0
		limit := 25
		var e error
		if q.Get("cursor") != "" {
			cursor, e = strconv.Atoi(q.Get("cursor"))
			if e != nil {
				writeProblem(w, 400, "INVALID_QUERY", "numeric cursor required")
				return
			}
		}
		if q.Get("limit") != "" {
			limit, e = strconv.Atoi(q.Get("limit"))
			if e != nil {
				writeProblem(w, 400, "INVALID_QUERY", "limit required")
				return
			}
		}
		out, e := m.Store.JournalLines(r.Context(), p, org, r.PathValue("id"), cursor, limit)
		if e != nil {
			financeProblem(w, e)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, 200, out)
	})
	mux.HandleFunc("GET /v1/finance/trial-balance", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		q := r.URL.Query()
		if len(q) != 2 || len(q["organization_id"]) != 1 || len(q["period_id"]) != 1 {
			writeProblem(w, 400, "INVALID_QUERY", "organization and period required")
			return
		}
		if !accountingOrg(w, p, q.Get("organization_id")) {
			return
		}
		out, e := m.Store.TrialBalance(r.Context(), p, q.Get("organization_id"), q.Get("period_id"))
		if e != nil {
			financeProblem(w, e)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write(out)
	})
	mux.HandleFunc("GET /v1/finance/commands/result", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		q := r.URL.Query()
		if len(q) != 1 || len(q["organization_id"]) != 1 || len(r.Header.Values("Idempotency-Key")) != 1 {
			writeProblem(w, 400, "INVALID_QUERY", "exact organization and key required")
			return
		}
		if !accountingOrg(w, p, q.Get("organization_id")) {
			return
		}
		out, e := m.Store.CommandResult(r.Context(), p, q.Get("organization_id"), r.Header.Get("Idempotency-Key"))
		if e != nil {
			financeProblem(w, e)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, 200, out)
	})
	mux.HandleFunc("POST /v1/finance/commands", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		if r.URL.RawQuery != "" || len(r.Header.Values("Idempotency-Key")) != 1 || r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 400, "INVALID_CONTRACT", "exact command required")
			return
		}
		controller := http.NewResponseController(w)
		if e := controller.SetReadDeadline(time.Now().Add(5 * time.Second)); e != nil && !errors.Is(e, http.ErrNotSupported) {
			writeProblem(w, 503, "UNAVAILABLE", "request unavailable")
			return
		}
		body, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 16384))
		if e != nil {
			writeProblem(w, 413, "BODY_TOO_LARGE", "bounded command required")
			return
		}
		var c postgres.FinanceCommand
		if bcfx.DecodeExactJSON(body, &c) != nil {
			writeProblem(w, 400, "INVALID_CONTRACT", "exact command required")
			return
		}
		if !accountingOrg(w, p, c.OrganizationID) {
			return
		}
		permission := postgres.FinancePermission(c.Action)
		if permission == "" || !p.Allowed(permission) {
			writeProblem(w, 403, "FORBIDDEN", "action permission required")
			return
		}
		out, replay, e := m.Store.Command(r.Context(), p, r.Header.Get("Idempotency-Key"), c)
		if e != nil {
			financeProblem(w, e)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		if replay {
			w.Header().Set("Idempotent-Replay", "true")
		}
		writeJSON(w, 201, out)
	})
}
func financeProblem(w http.ResponseWriter, e error) {
	if errors.Is(e, postgres.ErrFinanceNotFound) {
		writeProblem(w, 404, "NOT_FOUND", "receipt not found")
		return
	}
	if errors.Is(e, royalty.ErrInvalid) || errors.Is(e, accounting.ErrInvalid) {
		writeProblem(w, 400, "INVALID_FINANCE_COMMAND", "command or scope rejected")
		return
	}
	if errors.Is(e, royalty.ErrConflict) || errors.Is(e, accounting.ErrConflict) {
		writeProblem(w, 409, "FINANCE_CONFLICT", "state or exact request conflict")
		return
	}
	writeProblem(w, 503, "FINANCE_UNAVAILABLE", "finance unavailable")
}
