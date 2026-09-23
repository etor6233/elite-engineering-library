// AUTHORED bounded transport. Verified principal/scope and existing domain
// owners determine effects; request JSON cannot choose tenant or actor.
package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgconn"
)

type WarrantyService interface {
	Scope() (string, string)
	ProfileSHA256() string
	ProfileDocument(context.Context, identity.Principal) (json.RawMessage, error)
	Offer(context.Context, identity.Principal, string) (wc.Offer, error)
	BindOffer(context.Context, identity.Principal, wc.OfferRequest) (wc.Offer, error)
	Acknowledge(context.Context, identity.Principal, wc.AcknowledgeRequest) (wc.Offer, error)
	Activation(context.Context, identity.Principal, string) (wc.Activation, error)
	Activate(context.Context, identity.Principal, string) (wc.Activation, error)
	Claim(context.Context, identity.Principal, string) (wc.Claim, error)
	ClaimCommand(context.Context, identity.Principal, string, string) (wc.Step, error)
	OpenClaim(context.Context, identity.Principal, wc.OpenClaim) (wc.Step, error)
	Diagnose(context.Context, identity.Principal, wc.Diagnose) (wc.Step, error)
	PlanRepair(context.Context, identity.Principal, wc.Plan) (wc.Step, error)
	DecideRepair(context.Context, identity.Principal, wc.DecideRepair) (wc.Step, error)
	CompleteRepairWork(context.Context, identity.Principal, wc.CompleteWork) (wc.Step, error)
	RecordRepairQuality(context.Context, identity.Principal, wc.Quality) (wc.Step, error)
	AcceptRepair(context.Context, identity.Principal, wc.AcceptRepair) (wc.Step, error)
	ReconcileRepair(context.Context, identity.Principal, wc.ReconcileRepair) (wc.Step, error)
	CancelRepair(context.Context, identity.Principal, wc.CancelRepair) (wc.Step, error)
}
type WarrantyQuoteReader interface {
	QuoteForOffer(context.Context, identity.Principal, string) (wc.QuoteForOffer, error)
}
type WarrantyModule struct{ Service WarrantyService }
type warrantyProtection func(http.ResponseWriter, *http.Request, string, bool) (identity.Principal, bool)

func warrantyHTTPError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, wc.ErrNotFound), errors.Is(err, approval.ErrNotFound):
		writeProblem(w, 404, "WARRANTY_NOT_FOUND", "no warranty resource exists in this authorized scope")
	case errors.Is(err, wc.ErrInvalid):
		writeProblem(w, 400, "WARRANTY_INVALID", "the warranty command is invalid")
	case errors.Is(err, wc.ErrConflict), errors.Is(err, approval.ErrDuplicate), errors.Is(err, approval.ErrNotPending), errors.Is(err, approval.ErrSeparation), errors.Is(err, approval.ErrInvalidRequest):
		writeProblem(w, 409, "WARRANTY_CONFLICT", "consult the saved command and current case version before retrying")
	case errors.As(err, &pg) && (pg.Code == "23505" || pg.Code == "23514" || pg.Code == "P0001" || pg.Code == "40001" || pg.Code == "40P01"):
		writeProblem(w, 409, "WARRANTY_CONFLICT", "consult the saved command and current case version before retrying")
	default:
		writeProblem(w, 503, "WARRANTY_UNCONFIRMED", "the result is unconfirmed; consult its saved command reference")
	}
	return true
}
func warrantyDecode(w http.ResponseWriter, r *http.Request, out any) bool {
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if err != nil {
		writeProblem(w, 413, "WARRANTY_BODY_TOO_LARGE", "the command exceeds its bounded size")
		return false
	}
	// CanonicalPayload rejects duplicate keys, including case-insensitive ones,
	// before typed decoding could discard an ambiguous field.
	if _, _, err = approval.CanonicalPayload(raw); err != nil {
		writeProblem(w, 400, "WARRANTY_INVALID_JSON", "a bounded unambiguous JSON object is required")
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if decoder.Decode(out) != nil || decoder.Decode(new(any)) != io.EOF {
		writeProblem(w, 400, "WARRANTY_INVALID_JSON", "unknown fields or invalid typed values")
		return false
	}
	return true
}
func warrantyReply(w http.ResponseWriter, status int, value any) {
	replay := false
	switch v := value.(type) {
	case wc.Step:
		replay = v.Replay
	case wc.Offer:
		replay = v.Replay
	case wc.Activation:
		replay = v.Replay
	}
	if replay {
		status = 200
	}
	writeJSON(w, status, value)
}
func warrantyPost[T any](mux *http.ServeMux, route string, protect warrantyProtection, permission string, factory bool, id func(T) string, action func(context.Context, identity.Principal, T) (any, error)) {
	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, permission, factory)
		if !ok {
			return
		}
		var c T
		if !warrantyDecode(w, r, &c) {
			return
		}
		if path := r.PathValue("id"); path != "" && id(c) != path {
			warrantyHTTPError(w, wc.ErrInvalid)
			return
		}
		value, err := action(r.Context(), p, c)
		if warrantyHTTPError(w, err) {
			return
		}
		warrantyReply(w, 201, value)
	})
}
func (m WarrantyModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	api := franchiseJourneyAPI{verifier: verifier}
	protect := func(w http.ResponseWriter, r *http.Request, permission string, factory bool) (identity.Principal, bool) {
		query := r.URL.Query()
		if len(query["organization_id"]) != 1 || query.Get("organization_id") == "" {
			writeProblem(w, 400, "WARRANTY_SCOPE_REQUIRED", "one organization_id is required")
			return identity.Principal{}, false
		}
		org := query.Get("organization_id")
		p, ok := api.protected(w, r, permission, org)
		if !ok {
			return p, false
		}
		tenant, bound := m.Service.Scope()
		if p.TenantID != tenant || !factory && org != bound {
			writeProblem(w, 403, "FORBIDDEN", "the module is outside this request scope")
			return p, false
		}
		// Pass only the authority verified for this route/organization. A principal
		// with several organizations or '*' cannot silently switch request scope.
		p.Permissions = map[string]struct{}{permission: {}}
		p.Organizations = map[string]struct{}{org: {}}
		return p, true
	}
	mux.HandleFunc("GET /v1/franchise/warranty/profile", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "warranty:read", false)
		if !ok {
			return
		}
		raw, err := m.Service.ProfileDocument(r.Context(), p)
		if warrantyHTTPError(w, err) {
			return
		}
		writeJSON(w, 200, map[string]any{"profile_sha256": m.Service.ProfileSHA256(), "profile": raw})
	})
	if reader, ok := m.Service.(WarrantyQuoteReader); ok {
		mux.HandleFunc("GET /v1/franchise/warranty/quotes/{id}/source", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "warranty:offer", false)
			if !ok {
				return
			}
			v, e := reader.QuoteForOffer(r.Context(), p, r.PathValue("id"))
			if !warrantyHTTPError(w, e) {
				warrantyReply(w, 200, v)
			}
		})
	}
	for _, surface := range []struct{ path, permission string }{{"franchise", "warranty:read"}, {"customer", "warranty:self"}} {
		mux.HandleFunc("GET /v1/"+surface.path+"/warranty/quotes/{id}", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, surface.permission, false)
			if !ok {
				return
			}
			v, err := m.Service.Offer(r.Context(), p, r.PathValue("id"))
			if !warrantyHTTPError(w, err) {
				warrantyReply(w, 200, v)
			}
		})
		mux.HandleFunc("GET /v1/"+surface.path+"/warranty/handovers/{id}/activation", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, surface.permission, false)
			if !ok {
				return
			}
			v, err := m.Service.Activation(r.Context(), p, r.PathValue("id"))
			if !warrantyHTTPError(w, err) {
				warrantyReply(w, 200, v)
			}
		})
	}
	for _, surface := range []struct {
		path, permission string
		factory          bool
	}{{"franchise", "warranty:read", false}, {"customer", "warranty:self", false}, {"factory", "warranty:factory-read", true}} {
		mux.HandleFunc("GET /v1/"+surface.path+"/warranty/claims/{id}", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, surface.permission, surface.factory)
			if !ok {
				return
			}
			if command := r.URL.Query().Get("command_id"); command != "" {
				v, err := m.Service.ClaimCommand(r.Context(), p, r.PathValue("id"), command)
				if !warrantyHTTPError(w, err) {
					warrantyReply(w, 200, v)
				}
				return
			}
			v, err := m.Service.Claim(r.Context(), p, r.PathValue("id"))
			if !warrantyHTTPError(w, err) {
				warrantyReply(w, 200, v)
			}
		})
	}
	warrantyPost(mux, "POST /v1/franchise/warranty/quotes/{id}/terms", protect, "warranty:offer", false, func(c wc.OfferRequest) string { return c.QuoteID }, func(ctx context.Context, p identity.Principal, c wc.OfferRequest) (any, error) {
		return m.Service.BindOffer(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/customer/warranty/quotes/{id}/acknowledgements", protect, "warranty:self", false, func(c wc.AcknowledgeRequest) string { return c.QuoteID }, func(ctx context.Context, p identity.Principal, c wc.AcknowledgeRequest) (any, error) {
		return m.Service.Acknowledge(ctx, p, c)
	})
	type activate struct {
		HandoverID string `json:"handover_id"`
	}
	warrantyPost(mux, "POST /v1/franchise/warranty/handovers/{id}/activation", protect, "warranty:activate", false, func(c activate) string { return c.HandoverID }, func(ctx context.Context, p identity.Principal, c activate) (any, error) {
		return m.Service.Activate(ctx, p, c.HandoverID)
	})
	for _, surface := range []struct{ path, permission string }{{"franchise", "warranty:request"}, {"customer", "warranty:self"}} {
		warrantyPost(mux, "POST /v1/"+surface.path+"/warranty/claims", protect, surface.permission, false, func(c wc.OpenClaim) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.OpenClaim) (any, error) {
			return m.Service.OpenClaim(ctx, p, c)
		})
	}
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/diagnosis", protect, "warranty:diagnose", false, func(c wc.Diagnose) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.Diagnose) (any, error) {
		return m.Service.Diagnose(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/plan", protect, "warranty:plan", false, func(c wc.Plan) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.Plan) (any, error) {
		return m.Service.PlanRepair(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/decision", protect, "warranty:approve", false, func(c wc.DecideRepair) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.DecideRepair) (any, error) {
		return m.Service.DecideRepair(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/work", protect, "warranty:work", false, func(c wc.CompleteWork) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.CompleteWork) (any, error) {
		return m.Service.CompleteRepairWork(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/quality", protect, "warranty:quality", false, func(c wc.Quality) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.Quality) (any, error) {
		return m.Service.RecordRepairQuality(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/customer/warranty/claims/{id}/acceptance", protect, "warranty:self", false, func(c wc.AcceptRepair) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.AcceptRepair) (any, error) {
		return m.Service.AcceptRepair(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/factory/warranty/claims/{id}/reconciliation", protect, "warranty:reconcile", true, func(c wc.ReconcileRepair) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.ReconcileRepair) (any, error) {
		return m.Service.ReconcileRepair(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/cancellation", protect, "warranty:cancel", false, func(c wc.CancelRepair) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.CancelRepair) (any, error) {
		return m.Service.CancelRepair(ctx, p, c)
	})
}
