package franchisejourney

// AUTHORED composition glue: persist the explicitly selected commercial
// checkpoint over existing acceptance/payment/stock owners. No shipment, money
// movement, tax treatment or inventory journal is created by this operation.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

const CommercialReleaseEffect = "COMMIT_COMMERCIAL_RELEASE_RECEIPT"
const CommercialReleaseAlgorithmRevision = 2

func SupportedCommercialReleaseOptions() HandoverProfileOptions {
	opts := SupportedHandoverProfileOptions()
	opts.ReleaseEffect = CommercialReleaseEffect
	return opts
}

// Explicit additive selection; revisions1/2 retain provider-only coverage.
func SupportedStoredValueReleaseOptions() HandoverProfileOptions {
	opts := SupportedCommercialReleaseOptions()
	opts.PaymentCoverage = "FULL_ORDER_WITH_STORED_VALUE"
	return opts
}
func (p HandoverReleaseContract) AllowsStoredValueFunding(tenant, org string) bool {
	return p.AllowsScope(tenant, org) && p.profile != nil && p.profile.storedValue
}
func supportedHandoverAlgorithm(revision int, options HandoverProfileOptions) bool {
	return revision == HandoverSupportedAlgorithmRevision && options == SupportedHandoverProfileOptions() ||
		revision == CommercialReleaseAlgorithmRevision && options == SupportedCommercialReleaseOptions() || revision == 3 && options == SupportedStoredValueReleaseOptions()
}

func (p HandoverReleaseContract) AllowsCommercialRelease(tenant, organization string) bool {
	return p.AllowsScope(tenant, organization) && p.profile != nil && p.profile.releaseEffect == CommercialReleaseEffect
}

type CommitCommercialReleaseCommand struct {
	OrganizationID    string `json:"organization_id"`
	HandoverID        string `json:"handover_id"`
	ObservationSHA256 string `json:"observation_sha256"`
	IdempotencyKey    string `json:"-"`
}

// CommercialReleaseReceipt is immutable historical evidence of one committed
// checkpoint. It is deliberately not a current shipping or payment authority.
type CommercialReleaseReceipt struct {
	ID                    string    `json:"id"`
	OrganizationID        string    `json:"organization_id"`
	HandoverID            string    `json:"handover_id"`
	OrderID               string    `json:"order_id"`
	PaymentAttemptID      string    `json:"payment_attempt_id"`
	FundingReceiptID      string    `json:"funding_receipt_id,omitempty"`
	ObservationSHA256     string    `json:"observation_sha256"`
	ObservationGeneration int64     `json:"observation_generation"`
	HandoverVersion       int       `json:"handover_version"`
	AcceptanceSHA256      string    `json:"acceptance_sha256"`
	ChecklistID           string    `json:"checklist_id"`
	ChecklistVersion      int       `json:"checklist_version"`
	ContractID            string    `json:"contract_id"`
	ContractSHA256        string    `json:"contract_sha256"`
	Effect                string    `json:"effect"`
	ReleasedBy            string    `json:"released_by"`
	RecordedAt            time.Time `json:"recorded_at"`
	ValidUntil            time.Time `json:"valid_until"`
}

type CurrentCommercialRelease struct {
	Receipt     CommercialReleaseReceipt `json:"receipt"`
	EvaluatedAt time.Time                `json:"evaluated_at"`
	Current     bool                     `json:"current"`
}

type CommercialReleaseRepository interface {
	CommitInitialCommercialRelease(context.Context, string, string, string, string, CommitCommercialReleaseCommand, HandoverReleaseContract, string) (CommercialReleaseReceipt, bool, error)
	InitialCommercialReleaseResult(context.Context, string, string, string, string) (CommercialReleaseReceipt, error)
	ValidateInitialCommercialRelease(context.Context, string, string, string, HandoverReleaseContract) (CurrentCommercialRelease, error)
}

func (s *HandoverPreparationService) commercialRepository(tenant, organization string) (CommercialReleaseRepository, error) {
	if s == nil || !s.contract.AllowsCommercialRelease(tenant, organization) {
		return nil, ErrReleaseConditioned
	}
	r, ok := s.repository.(CommercialReleaseRepository)
	if !ok {
		return nil, ErrReleaseConditioned
	}
	return r, nil
}

func (s *HandoverPreparationService) CommitCommercialRelease(ctx context.Context, tenant, actor string, c CommitCommercialReleaseCommand) (CommercialReleaseReceipt, bool, error) {
	var empty CommercialReleaseReceipt
	r, err := s.commercialRepository(tenant, c.OrganizationID)
	if err != nil {
		return empty, false, err
	}
	if tenant == "" || actor == "" || len(actor) > 256 || !validPreparationID(c.OrganizationID) || !validPreparationID(c.HandoverID) || !sha256Hex(c.ObservationSHA256) || len(c.IdempotencyKey) < 8 || len(c.IdempotencyKey) > 128 {
		return empty, false, ErrInvalid
	}
	b, err := json.Marshal(struct {
		Actor       string
		Command     CommitCommercialReleaseCommand
		ContractSHA string
	}{actor, c, s.contract.DocumentSHA256})
	if err != nil {
		return empty, false, err
	}
	h := sha256.Sum256(b)
	return r.CommitInitialCommercialRelease(ctx, tenant, actor, s.ids.New(), s.ids.New(), c, s.contract, hex.EncodeToString(h[:]))
}

func (s *HandoverPreparationService) CommercialReleaseResult(ctx context.Context, tenant, organization, handover, key string) (CommercialReleaseReceipt, error) {
	r, err := s.commercialRepository(tenant, organization)
	if err != nil {
		return CommercialReleaseReceipt{}, err
	}
	if !validPreparationID(handover) || len(key) < 8 || len(key) > 128 {
		return CommercialReleaseReceipt{}, ErrInvalid
	}
	return r.InitialCommercialReleaseResult(ctx, tenant, organization, handover, key)
}

// ValidateCommercialRelease rechecks current owners under their locks, then
// matches the receipt's generation and accepted handover version. Its response
// is read-only; a later physical effect must revalidate inside its own write.
func (s *HandoverPreparationService) ValidateCommercialRelease(ctx context.Context, tenant, organization, handover string) (CurrentCommercialRelease, error) {
	r, err := s.commercialRepository(tenant, organization)
	if err != nil {
		return CurrentCommercialRelease{}, err
	}
	if !validPreparationID(handover) {
		return CurrentCommercialRelease{}, ErrInvalid
	}
	return r.ValidateInitialCommercialRelease(ctx, tenant, organization, handover, s.contract)
}
