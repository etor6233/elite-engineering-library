package httpapi

// AUTHORED HTTP glue. The module is mounted only with an explicitly configured
// service; claims, tenant, actor and organization are verified server-side.
import (
	"context"
	"errors"
	"net/http"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
)

type InitialHandoverService interface {
	Prepare(context.Context, string, string, franchisejourney.PrepareHandoverCommand) (franchisejourney.HandoverPreparation, bool, error)
	Result(context.Context, string, string, string, string) (franchisejourney.HandoverPreparation, error)
	EvaluateRelease(context.Context, string, string, string, string) (franchisejourney.HandoverReferenceRelease, error)
}
type InitialHandoverModule struct{ Service InitialHandoverService }

func (m InitialHandoverModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	principal := franchiseJourneyAPI{verifier: verifier}
	mux.HandleFunc("POST /v1/franchise/orders/{id}/handover", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			OrganizationID    string `json:"organization_id"`
			OrderLineID       string `json:"order_line_id"`
			PaymentAttemptID  string `json:"payment_attempt_id"`
			FundingReceiptID  string `json:"funding_receipt_id"`
			ObservationSHA256 string `json:"observation_sha256"`
		}
		if !decodeStrict(w, r, &input) {
			return
		}
		actor, ok := principal.protected(w, r, "handover:manage", input.OrganizationID)
		if !ok {
			return
		}
		value, replay, err := m.Service.Prepare(r.Context(), actor.TenantID, actor.Subject, franchisejourney.PrepareHandoverCommand{OrganizationID: input.OrganizationID, OrderID: r.PathValue("id"), OrderLineID: input.OrderLineID, PaymentAttemptID: input.PaymentAttemptID, FundingReceiptID: input.FundingReceiptID, ObservationSHA256: input.ObservationSHA256, IdempotencyKey: r.Header.Get("Idempotency-Key")})
		if initialHandoverError(w, err) {
			return
		}
		if replay {
			w.Header().Set("Idempotency-Replayed", "true")
		}
		writeJSON(w, http.StatusCreated, value)
	})
	mux.HandleFunc("GET /v1/franchise/orders/{id}/handover-result", func(w http.ResponseWriter, r *http.Request) {
		organization := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", organization)
		if !ok {
			return
		}
		value, err := m.Service.Result(r.Context(), actor.TenantID, organization, r.PathValue("id"), r.Header.Get("Idempotency-Key"))
		if initialHandoverError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, value)
	})
	mux.HandleFunc("GET /v1/franchise/handovers/{id}/release-check", func(w http.ResponseWriter, r *http.Request) {
		organization := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", organization)
		if !ok {
			return
		}
		value, err := m.Service.EvaluateRelease(r.Context(), actor.TenantID, organization, r.PathValue("id"), r.URL.Query().Get("observation_sha256"))
		if initialHandoverError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, value)
	})
}
func initialHandoverError(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, franchisejourney.ErrInvalid):
		writeProblem(w, 400, "INVALID_HANDOVER_PREPARATION", "handover request does not match its contract")
	case errors.Is(err, franchisejourney.ErrNotFound):
		writeProblem(w, 404, "HANDOVER_PREPARATION_NOT_FOUND", "no completed preparation exists for this scope and key")
	case errors.Is(err, franchisejourney.ErrConflict):
		writeProblem(w, 409, "HANDOVER_NOT_ELIGIBLE", "current order, stock, payment or acceptance evidence does not permit this operation")
	case errors.Is(err, franchisejourney.ErrReleaseConditioned):
		writeProblem(w, 503, "HANDOVER_CONTRACT_REQUIRED", "an admitted release contract must be configured")
	default:
		writeProblem(w, 500, "HANDOVER_PREPARATION_FAILED", "handover operation could not complete")
	}
	return true
}
