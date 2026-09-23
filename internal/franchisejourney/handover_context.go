package franchisejourney

// AUTHORED read composition for the operator UI. Every write still revalidates
// the actual owners; this view grants no authority to deliver or move money.
import (
	"context"
	"time"
)

type HandoverOperatorContext struct {
	OrganizationID    string    `json:"organization_id"`
	OrderID           string    `json:"order_id"`
	OrderLineID       string    `json:"order_line_id"`
	PaymentAttemptID  string    `json:"payment_attempt_id"`
	FundingReceiptID  string    `json:"funding_receipt_id,omitempty"`
	ObservationSHA256 string    `json:"observation_sha256"`
	Handover          *Handover `json:"handover,omitempty"`
	CanPrepare        bool      `json:"can_prepare"`
	ReleaseEffect     string    `json:"release_effect"`
	EvaluatedAt       time.Time `json:"evaluated_at"`
}

func (s *HandoverPreparationService) OperatorContext(ctx context.Context, tenant, org, order string) (HandoverOperatorContext, error) {
	if s == nil || !s.contract.AllowsScope(tenant, org) {
		return HandoverOperatorContext{}, ErrReleaseConditioned
	}
	if !validPreparationID(order) {
		return HandoverOperatorContext{}, ErrInvalid
	}
	r, ok := s.repository.(interface {
		InitialHandoverOperatorContext(context.Context, string, string, string, HandoverReleaseContract) (HandoverOperatorContext, error)
	})
	if !ok {
		return HandoverOperatorContext{}, ErrReleaseConditioned
	}
	return r.InitialHandoverOperatorContext(ctx, tenant, org, order, s.contract)
}
