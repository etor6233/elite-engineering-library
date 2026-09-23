// AUTHORED protected transport for the connected serial supply owner.
package httpapi

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"net/http"
)

type SerialSupplyService interface {
	Plan(context.Context, identity.Principal, string, string) (sc.Plan, error)
	CommandReceipt(context.Context, identity.Principal, string, string) (sc.Receipt, error)
	BindPlan(context.Context, identity.Principal, sc.PlanRequest) (sc.Receipt, error)
	Apply(context.Context, identity.Principal, sc.Command) (sc.Receipt, error)
}
type SerialSupplyCreator interface {
	Create(context.Context, identity.Principal, sc.CreateRequest) (sc.Receipt, error)
}
type SerialSupplyModule struct{ Service SerialSupplyService }

func supplyHTTPError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, sc.ErrInvalid):
		writeProblem(w, 400, "SUPPLY_INVALID", "invalid serial supply command")
	case errors.Is(err, sc.ErrNotFound), errors.Is(err, approval.ErrNotFound):
		writeProblem(w, 404, "SUPPLY_NOT_FOUND", "resource is outside the authorized scope")
	case errors.Is(err, sc.ErrConflict), errors.Is(err, operations.ErrConflict), errors.Is(err, approval.ErrSeparation), errors.Is(err, approval.ErrNotPending), errors.Is(err, approval.ErrInvalidRequest), errors.Is(err, approval.ErrDuplicate):
		writeProblem(w, 409, "SUPPLY_CONFLICT", "consult the saved command and current order version")
	case errors.As(err, &pg) && (pg.Code == "23505" || pg.Code == "23514" || pg.Code == "P0001" || pg.Code == "40001" || pg.Code == "40P01"):
		writeProblem(w, 409, "SUPPLY_CONFLICT", "consult the saved command and current order version")
	default:
		writeProblem(w, 503, "SUPPLY_UNCONFIRMED", "result unconfirmed; consult the saved command reference")
	}
	return true
}
func supplyDecode(w http.ResponseWriter, r *http.Request, out any) bool {
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if err != nil {
		writeProblem(w, 413, "SUPPLY_BODY_TOO_LARGE", "command exceeds its bounded size")
		return false
	}
	if _, _, err = approval.CanonicalPayload(raw); err != nil {
		writeProblem(w, 400, "SUPPLY_INVALID_JSON", "unambiguous JSON object required")
		return false
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(out) != nil || d.Decode(new(any)) != io.EOF {
		writeProblem(w, 400, "SUPPLY_INVALID_JSON", "unknown fields or invalid typed values")
		return false
	}
	return true
}
func supplyReply(w http.ResponseWriter, v sc.Receipt) {
	status := http.StatusCreated
	if v.Replay {
		status = http.StatusOK
	}
	writeJSON(w, status, v)
}
func (m SerialSupplyModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	api := franchiseJourneyAPI{verifier: verifier}
	protect := func(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
		w.Header().Set("Cache-Control", "no-store")
		query := r.URL.Query()
		for key, values := range query {
			if len(values) != 1 || values[0] == "" || (key != "organization_id" && (r.Method != "GET" || (key != "command_id" && key != "after_unit"))) {
				writeProblem(w, 400, "SUPPLY_INVALID_QUERY", "one explicit organization and bounded query parameters required")
				return identity.Principal{}, false
			}
		}
		if query.Get("organization_id") == "" || (query.Get("command_id") != "" && query.Get("after_unit") != "") {
			writeProblem(w, 400, "SUPPLY_SCOPE_REQUIRED", "one organization_id and one recovery or page cursor required")
			return identity.Principal{}, false
		}
		org := query.Get("organization_id")
		p, ok := api.protected(w, r, permission, org)
		if !ok {
			return p, false
		}
		p.Permissions = map[string]struct{}{permission: {}}
		p.Organizations = map[string]struct{}{org: {}}
		return p, true
	}
	for _, surface := range []struct{ path, permission string }{{"franchise", "supply:read"}, {"factory", "supply:factory-read"}} {
		mux.HandleFunc("GET /v1/"+surface.path+"/supply/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, surface.permission)
			if !ok {
				return
			}
			if key := r.URL.Query().Get("command_id"); key != "" {
				v, err := m.Service.CommandReceipt(r.Context(), p, r.PathValue("id"), key)
				if !supplyHTTPError(w, err) {
					writeJSON(w, 200, v)
				}
				return
			}
			v, err := m.Service.Plan(r.Context(), p, r.PathValue("id"), r.URL.Query().Get("after_unit"))
			if !supplyHTTPError(w, err) {
				writeJSON(w, 200, v)
			}
		})
	}
	mux.HandleFunc("POST /v1/franchise/supply/orders/{id}/plan", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "supply:plan")
		if !ok {
			return
		}
		var c sc.PlanRequest
		if !supplyDecode(w, r, &c) {
			return
		}
		if c.PurchaseOrderID != r.PathValue("id") {
			supplyHTTPError(w, sc.ErrInvalid)
			return
		}
		v, err := m.Service.BindPlan(r.Context(), p, c)
		if !supplyHTTPError(w, err) {
			supplyReply(w, v)
		}
	})

	if creator, ok := m.Service.(SerialSupplyCreator); ok {
		mux.HandleFunc("POST /v1/franchise/supply/orders/{id}/create", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "supply:plan")
			if !ok {
				return
			}
			var c sc.CreateRequest
			if !supplyDecode(w, r, &c) {
				return
			}
			if c.PurchaseOrderID != r.PathValue("id") || c.DestinationOrganizationID != r.URL.Query().Get("organization_id") {
				supplyHTTPError(w, sc.ErrInvalid)
				return
			}
			v, err := creator.Create(r.Context(), p, c)
			if !supplyHTTPError(w, err) {
				supplyReply(w, v)
			}
		})
	}
	routes := []struct{ surface, kind, permission string }{
		{"franchise", "submit", "supply:plan"}, {"franchise", "cancel", "supply:plan"},
		{"factory", "confirm", "supply:factory"}, {"factory", "start", "supply:factory"},
		{"factory", "register", "supply:factory"}, {"factory", "milestone", "supply:factory"}, {"factory", "ship", "supply:factory"},
		{"franchise", "receive", "supply:receive"}, {"franchise", "quality", "supply:release"}, {"franchise", "quality-reject", "supply:release"}, {"franchise", "reinspect", "supply:inspect"},
	}
	for _, route := range routes {
		mux.HandleFunc("POST /v1/"+route.surface+"/supply/orders/{id}/"+route.kind, func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, route.permission)
			if !ok {
				return
			}
			var c sc.Command
			if !supplyDecode(w, r, &c) {
				return
			}
			if c.PurchaseOrderID != r.PathValue("id") || c.Kind != route.kind {
				supplyHTTPError(w, sc.ErrInvalid)
				return
			}
			v, err := m.Service.Apply(r.Context(), p, c)
			if !supplyHTTPError(w, err) {
				supplyReply(w, v)
			}
		})
	}
}
