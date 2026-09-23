package httpapi

// AUTHORED monetary-safe transport only. Identity and scope come from verified
// bearer claims; exact source calculations and durable effects stay with owners.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"strconv"
)

type StoredValueService interface {
	Metrics(context.Context, identity.Principal, string, string) (db.StoredValueMetrics, error)
	Scope() (string, string)
	ProfileDocument(context.Context, identity.Principal, string) (json.RawMessage, error)
	ProfileSHA256() string
	Approvals(context.Context, identity.Principal, string, string) (db.StoredValueApprovalPage, error)
	Propose(context.Context, identity.Principal, sv.Request) (sv.ApprovalView, error)
	Read(context.Context, identity.Principal, string) (sv.ApprovalView, error)
	Decide(context.Context, identity.Principal, sv.Decision) (sv.ApprovalView, error)
	Allocation(context.Context, identity.Principal, string) (sv.Allocation, error)
	FinalizeFunding(context.Context, identity.Principal, db.FinalizeStoredValueFunding) (db.LocalFundingResult, bool, error)
	FundingResult(context.Context, identity.Principal, string) (db.LocalFundingResult, error)
}
type StoredValueModule struct{ Service StoredValueService }

func storedValueVersion(raw string) (int64, error) {
	n, e := strconv.ParseInt(raw, 10, 64)
	if e != nil || n < 1 || strconv.FormatInt(n, 10) != raw {
		return 0, sv.ErrBinding
	}
	return n, nil
}
func storedValueError(w http.ResponseWriter, e error) bool {
	if e == nil {
		return false
	}
	var pg *pgconn.PgError
	switch {
	case errors.Is(e, approval.ErrNotFound) || errors.Is(e, pgx.ErrNoRows):
		writeProblem(w, 404, "STORED_VALUE_NOT_FOUND", "no result exists in this authorized scope")
	case errors.Is(e, sv.ErrBinding) || errors.Is(e, approval.ErrDuplicate) || errors.Is(e, approval.ErrNotPending) || errors.Is(e, approval.ErrSeparation) || errors.Is(e, approval.ErrInvalidRequest) || errors.As(e, &pg) && (pg.Code == "23514" || pg.Code == "23505"):
		writeProblem(w, 409, "STORED_VALUE_BINDING_REJECTED", "consult current order and approval state before retrying")
	default:
		writeProblem(w, 503, "STORED_VALUE_UNCONFIRMED", "consult the saved request reference; the result is not confirmed")
	}
	return true
}
func storedValueApprovalJSON(w http.ResponseWriter, status int, v sv.ApprovalView) {
	raw, hash, e := sv.Canonical(v.Payload)
	if e != nil || hash != v.PayloadSHA256 {
		storedValueError(w, sv.ErrBinding)
		return
	}
	entries := make([]map[string]string, 0, len(v.Payload.Entries))
	for _, e := range v.Payload.Entries {
		entries = append(entries, map[string]string{"account_id": e.AccountID, "points_delta": e.PointsDelta, "applied_minor_units": sv.Number(e.AppliedMinor)})
	}
	review := map[string]any{"operation_id": v.Payload.Request.OperationID, "program_id": v.Payload.Request.ProgramID, "original_operation_id": v.Payload.Request.OriginalOperationID, "currency": v.Payload.Order.Currency, "gross_minor_units": sv.Number(v.Payload.GrossMinor), "gift_minor_units": sv.Number(v.Payload.GiftMinor), "discount_minor_units": sv.Number(v.Payload.DiscountMinor), "entries": entries}
	// JSON document strings preserve signed int64 amounts across JavaScript.
	writeJSON(w, status, map[string]any{"review": review, "approval_id": v.ApprovalID, "payload_sha256": hash, "state": v.State, "replay": v.Replay, "payload_json": string(raw), "receipt_json": string(v.Receipt), "order_id": v.Payload.Request.OrderID, "organization_id": v.Payload.Request.OrganizationID, "operation": v.Payload.Request.Operation, "requester": v.Payload.Requester, "expires_at": v.Payload.ExpiresAt})
}
func storedValueFundingJSON(w http.ResponseWriter, status int, v db.LocalFundingResult) {
	raw, e := json.Marshal(v.Receipt)
	if e != nil {
		storedValueError(w, e)
		return
	}
	writeJSON(w, status, map[string]any{"funding_receipt_id": v.Receipt.ID, "receipt_sha256": v.SHA256, "receipt_json": string(raw), "organization_id": v.Receipt.Allocation.OrganizationID, "order_id": v.Receipt.Allocation.OrderID, "request_key": v.Receipt.RequestKey})
}
func (m StoredValueModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	api := franchiseJourneyAPI{verifier: verifier}
	protected := func(w http.ResponseWriter, r *http.Request, permission, org string) (identity.Principal, bool) {
		p, ok := api.protected(w, r, permission, org)
		if !ok {
			return p, false
		}
		tenant, bound := m.Service.Scope()
		if p.TenantID != tenant || org != bound {
			writeProblem(w, 403, "FORBIDDEN", "the selected program is outside this request scope")
			return p, false
		}
		return p, true
	}

	registerStoredValueMetrics(m, mux, protected)
	mux.HandleFunc("GET /v1/franchise/stored-value/profile", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		p, ok := protected(w, r, "stored_value:read", org)
		if !ok {
			return
		}
		raw, e := m.Service.ProfileDocument(r.Context(), p, org)
		if storedValueError(w, e) {
			return
		}
		writeJSON(w, 200, map[string]string{"profile_sha256": m.Service.ProfileSHA256(), "profile_json": string(raw), "organization_id": org})
	})
	mux.HandleFunc("GET /v1/franchise/stored-value/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		p, ok := protected(w, r, "stored_value:read", org)
		if !ok {
			return
		}
		a, e := m.Service.Allocation(r.Context(), p, r.PathValue("id"))
		if storedValueError(w, e) {
			return
		}
		page, e := m.Service.Approvals(r.Context(), p, r.PathValue("id"), r.URL.Query().Get("after"))
		if storedValueError(w, e) {
			return
		}
		writeJSON(w, 200, map[string]any{"organization_id": org, "order_id": a.OrderID, "order_version": sv.Number(a.OrderVersion), "currency": a.Currency, "gross_minor_units": sv.Number(a.GrossMinor), "gift_minor_units": sv.Number(a.GiftMinor), "discount_minor_units": sv.Number(a.DiscountMinor), "provider_due_minor_units": sv.Number(a.ProviderDueMinor), "approvals": page})
	})
	mux.HandleFunc("POST /v1/franchise/stored-value/proposals", func(w http.ResponseWriter, r *http.Request) {
		var c struct {
			OperationID          string `json:"operation_id"`
			Operation            string `json:"operation"`
			OrganizationID       string `json:"organization_id"`
			ProgramID            string `json:"program_id"`
			OrderID              string `json:"order_id"`
			ExpectedOrderVersion string `json:"expected_order_version"`
			AccountID            string `json:"account_id"`
			OriginalOperationID  string `json:"original_operation_id"`
			ProfileSHA256        string `json:"profile_sha256"`
		}
		if !decodeStrict(w, r, &c) {
			return
		}
		p, ok := protected(w, r, "stored_value:request", c.OrganizationID)
		if !ok {
			return
		}
		n, e := storedValueVersion(c.ExpectedOrderVersion)
		if storedValueError(w, e) {
			return
		}
		v, e := m.Service.Propose(r.Context(), p, sv.Request{OperationID: c.OperationID, Operation: c.Operation, OrganizationID: c.OrganizationID, ProgramID: c.ProgramID, OrderID: c.OrderID, ExpectedOrderVersion: n, AccountID: c.AccountID, OriginalOperationID: c.OriginalOperationID, ProfileSHA256: c.ProfileSHA256})
		if storedValueError(w, e) {
			return
		}
		storedValueApprovalJSON(w, 201, v)
	})
	mux.HandleFunc("GET /v1/franchise/stored-value/approvals/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protected(w, r, "stored_value:read", r.URL.Query().Get("organization_id"))
		if !ok {
			return
		}
		v, e := m.Service.Read(r.Context(), p, r.PathValue("id"))
		if storedValueError(w, e) {
			return
		}
		storedValueApprovalJSON(w, 200, v)
	})
	mux.HandleFunc("POST /v1/franchise/stored-value/approvals/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		var c struct {
			OrganizationID string `json:"organization_id"`
			PayloadSHA256  string `json:"payload_sha256"`
			Approved       *bool  `json:"approved"`
			Reason         string `json:"reason"`
		}
		if !decodeStrict(w, r, &c) {
			return
		}
		p, ok := protected(w, r, "stored_value:approve", c.OrganizationID)
		if !ok {
			return
		}
		if c.Approved == nil {
			storedValueError(w, sv.ErrBinding)
			return
		}
		v, e := m.Service.Decide(r.Context(), p, sv.Decision{OrganizationID: c.OrganizationID, ApprovalID: r.PathValue("id"), PayloadSHA256: c.PayloadSHA256, Approved: *c.Approved, Reason: c.Reason})
		if storedValueError(w, e) {
			return
		}
		storedValueApprovalJSON(w, 200, v)
	})
	mux.HandleFunc("POST /v1/franchise/stored-value/orders/{id}/funding", func(w http.ResponseWriter, r *http.Request) {
		var c struct {
			OrganizationID       string `json:"organization_id"`
			ExpectedOrderVersion string `json:"expected_order_version"`
		}
		if !decodeStrict(w, r, &c) {
			return
		}
		p, ok := protected(w, r, "stored_value:fund", c.OrganizationID)
		if !ok {
			return
		}
		n, e := storedValueVersion(c.ExpectedOrderVersion)
		if storedValueError(w, e) {
			return
		}
		v, replay, e := m.Service.FinalizeFunding(r.Context(), p, db.FinalizeStoredValueFunding{OrganizationID: c.OrganizationID, OrderID: r.PathValue("id"), ExpectedOrderVersion: n, RequestKey: r.Header.Get("Idempotency-Key")})
		if storedValueError(w, e) {
			return
		}
		if replay {
			w.Header().Set("Idempotency-Replayed", "true")
		}
		storedValueFundingJSON(w, 201, v)
	})
	mux.HandleFunc("GET /v1/franchise/stored-value/funding-result", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protected(w, r, "stored_value:read", r.URL.Query().Get("organization_id"))
		if !ok {
			return
		}
		v, e := m.Service.FundingResult(r.Context(), p, r.Header.Get("Idempotency-Key"))
		if storedValueError(w, e) {
			return
		}
		storedValueFundingJSON(w, 200, v)
	})
}
