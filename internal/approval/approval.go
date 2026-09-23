// Package approval provides a governed human-approval queue for business
// effects (reservations, sales, refunds, payments): separation of duties,
// dual control for high-value requests, velocity-based fraud flags, immutable
// decisions and fail-closed money handling.
package approval

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Kind is the bounded set of approvable business effects.
type Kind string

const (
	KindReservation          Kind = "reservation"
	KindSale                 Kind = "sale"
	KindRefund               Kind = "refund"
	KindPayment              Kind = "payment"
	KindWhatsAppReply        Kind = "whatsapp_reply"
	KindSocialPublish        Kind = "social_publish"
	KindSocialRevoke         Kind = "social_revoke"
	KindStoredValueOperation Kind = "stored_value_operation"
	KindWarrantyRepair       Kind = "warranty_repair"
	KindSerialQuality        Kind = "serial_quality"
	KindCatalogReview        Kind = "catalog_review"
	KindTrainingAssessment   Kind = "training_assessment"
	KindMarketplaceMutation  Kind = "marketplace_mutation"
	KindWhatsAppSchedule     Kind = "whatsapp_schedule"
	KindDocumentReview       Kind = "document_review"
)

// State is the lifecycle of an approval request.
type State string

const (
	StatePending  State = "pending"
	StateApproved State = "approved"
	StateRejected State = "rejected"
)

var hex64Re = regexp.MustCompile(`^[0-9a-f]{64}$`)

var (
	ErrInvalidRequest    = errors.New("approval: invalid request")
	ErrDuplicate         = errors.New("approval: duplicate request")
	ErrSubjectFlagged    = errors.New("approval: subject flagged (velocity)")
	ErrNotFound          = errors.New("approval: not found")
	ErrNotPending        = errors.New("approval: not pending")
	ErrSeparation        = errors.New("approval: reviewer must differ from requester")
	ErrDuplicateApprover = errors.New("approval: reviewer already approved")
)

// Policy configures auto-approval and dual-control thresholds.
type Policy struct {
	AutoApproveMinorUnits int64 // <= threshold auto-approves; 0 disables auto-approval
	DualControlMinorUnits int64 // >= threshold requires two distinct reviewers
	MaxOpenPerSubject     int   // velocity bound per subject; 0 disables
}

// Request is an immutable, tenant-scoped approval request.
type Request struct {
	TenantID         string
	ID               string
	Kind             Kind
	SubjectID        string
	AmountMinorUnits int64
	Requester        string
	EvidenceSHA      string
}

// Validate enforces the request contract.
func (r Request) Validate() error {
	if strings.TrimSpace(r.TenantID) == "" {
		return fmt.Errorf("%w: tenant", ErrInvalidRequest)
	}
	if strings.TrimSpace(r.ID) == "" || len(r.ID) > 128 {
		return fmt.Errorf("%w: id", ErrInvalidRequest)
	}
	switch r.Kind {
	case KindReservation, KindSale, KindRefund, KindPayment, KindWhatsAppReply, KindSocialPublish, KindSocialRevoke, KindStoredValueOperation, KindWarrantyRepair, KindSerialQuality, KindCatalogReview, KindTrainingAssessment, KindMarketplaceMutation, KindWhatsAppSchedule, KindDocumentReview:
	default:
		return fmt.Errorf("%w: kind", ErrInvalidRequest)
	}
	if strings.TrimSpace(r.SubjectID) == "" || len(r.SubjectID) > 128 {
		return fmt.Errorf("%w: subject", ErrInvalidRequest)
	}
	if r.AmountMinorUnits < 0 {
		return fmt.Errorf("%w: amount", ErrInvalidRequest)
	}
	if strings.TrimSpace(r.Requester) == "" || len(r.Requester) > 128 {
		return fmt.Errorf("%w: requester", ErrInvalidRequest)
	}
	if !hex64Re.MatchString(strings.ToLower(r.EvidenceSHA)) {
		return fmt.Errorf("%w: evidence sha", ErrInvalidRequest)
	}
	return nil
}

// Decision is an immutable approval/rejection record.
type Decision struct {
	RequestID string
	Reviewer  string
	Approved  bool
	Reason    string
	At        time.Time
}
