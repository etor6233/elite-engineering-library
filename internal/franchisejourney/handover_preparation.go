package franchisejourney

// AUTHORED composition glue. This projects existing order, allocation and
// payment observations into a handover; it neither prices an order nor posts a
// shipment. The explicitly selected reference contract is not a live policy.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrReleaseConditioned = errors.New("handover release contract not admitted")

// Exact reference contract, deliberately scoped to synthetic/offline execution.
// The binding conditions are selected by the reference, not a fiscal or credit rule.
const ReferenceHandoverContractDocument = "reference-single-unit-observed-payment-v1\nLOCAL_FIXTURES\nOne active customer order line, quantity one, allocated serialized stock and active matching reservation.\nThe configured provider observation must match the complete order amount and currency, with no refund, dispute, hold or unobserved state.\nPreparation does not post inventory, establish fiscal compliance, deliver goods or authorize production.\nThe observation maximum age is explicitly configured between one nanosecond and fifteen minutes and included in the command hash.\n"

type HandoverReleaseContract struct {
	ID                    string        `json:"id"`
	DocumentSHA256        string        `json:"document_sha256"`
	Scope                 string        `json:"scope"`
	ExpectedLiveMode      bool          `json:"expected_live_mode"`
	MaximumObservationAge time.Duration `json:"-"`
	profile               *handoverProfileBinding
}

func (p HandoverReleaseContract) Valid() bool {
	if p.profile != nil {
		b := p.profile
		return p.ID == b.id && p.DocumentSHA256 == b.sha && p.Scope == "MATERIALIZED_PROFILE" && p.ExpectedLiveMode == b.mode && p.MaximumObservationAge == b.age
	}
	decoded, err := hex.DecodeString(p.DocumentSHA256)
	expected := sha256.Sum256([]byte(ReferenceHandoverContractDocument))
	return p.ID == "reference-single-unit-observed-payment-v1" && p.Scope == "LOCAL_FIXTURES" && !p.ExpectedLiveMode &&
		err == nil && len(decoded) == 32 && strings.ToLower(p.DocumentSHA256) == p.DocumentSHA256 &&
		p.DocumentSHA256 == hex.EncodeToString(expected[:]) &&
		p.MaximumObservationAge > 0 && p.MaximumObservationAge <= 15*time.Minute
}

func validPreparationID(value string) bool {
	return strings.TrimSpace(value) == value && value != "" && len(value) <= 128
}

type PrepareHandoverCommand struct {
	OrganizationID    string `json:"organization_id"`
	OrderID           string `json:"order_id"`
	OrderLineID       string `json:"order_line_id"`
	PaymentAttemptID  string `json:"payment_attempt_id"`
	FundingReceiptID  string `json:"funding_receipt_id,omitempty"`
	ObservationSHA256 string `json:"observation_sha256"`
	IdempotencyKey    string `json:"-"`
}

type HandoverPreparation struct {
	Handover          Handover  `json:"handover"`
	OrderLineID       string    `json:"order_line_id"`
	ReservationID     string    `json:"reservation_id"`
	PaymentAttemptID  string    `json:"payment_attempt_id"`
	FundingReceiptID  string    `json:"funding_receipt_id,omitempty"`
	ObservationSHA256 string    `json:"observation_sha256"`
	ContractID        string    `json:"contract_id"`
	ContractSHA256    string    `json:"contract_sha256"`
	PreparedBy        string    `json:"prepared_by"`
	PreparedAt        time.Time `json:"prepared_at"`
}

type HandoverPreparationRepository interface {
	PrepareInitialHandover(context.Context, string, string, string, string, PrepareHandoverCommand, HandoverReleaseContract, string) (HandoverPreparation, bool, error)
	InitialHandoverResult(context.Context, string, string, string, string) (HandoverPreparation, error)
}

type HandoverReferenceRelease struct {
	HandoverID        string    `json:"handover_id"`
	OrganizationID    string    `json:"organization_id"`
	ObservationSHA256 string    `json:"observation_sha256"`
	ContractSHA256    string    `json:"contract_sha256"`
	EvaluatedAt       time.Time `json:"evaluated_at"`
	Scope             string    `json:"scope"`
	Eligible          bool      `json:"eligible"`
}

// EvaluateRelease reports a current reference decision. It is not a durable
// shipping authorization; the effect owner must revalidate inside its own write.
func (s *HandoverPreparationService) EvaluateRelease(ctx context.Context, tenant, organization, handover, observationSHA string) (HandoverReferenceRelease, error) {
	if s == nil || !s.contract.AllowsScope(tenant, organization) {
		return HandoverReferenceRelease{}, ErrReleaseConditioned
	}
	if tenant == "" || !validPreparationID(organization) || !validPreparationID(handover) || !sha256Hex(observationSHA) {
		return HandoverReferenceRelease{}, ErrInvalid
	}
	repo, ok := s.repository.(interface {
		EvaluateInitialHandoverRelease(context.Context, string, string, string, string, HandoverReleaseContract) (HandoverReferenceRelease, error)
	})
	if !ok {
		return HandoverReferenceRelease{}, ErrReleaseConditioned
	}
	return repo.EvaluateInitialHandoverRelease(ctx, tenant, organization, handover, observationSHA, s.contract)
}

type HandoverPreparationService struct {
	repository HandoverPreparationRepository
	ids        IDGenerator
	contract   HandoverReleaseContract
}

func NewHandoverPreparationService(repository HandoverPreparationRepository, ids IDGenerator, contract HandoverReleaseContract) (*HandoverPreparationService, error) {
	if repository == nil || ids == nil || !contract.Valid() {
		return nil, ErrReleaseConditioned
	}
	return &HandoverPreparationService{repository: repository, ids: ids, contract: contract}, nil
}

func (s *HandoverPreparationService) Prepare(ctx context.Context, tenant, actor string, command PrepareHandoverCommand) (HandoverPreparation, bool, error) {
	var empty HandoverPreparation
	if s == nil || !s.contract.AllowsScope(tenant, command.OrganizationID) {
		return empty, false, ErrReleaseConditioned
	}
	if tenant == "" || actor == "" || len(actor) > 256 || !validPreparationID(command.OrganizationID) || !validPreparationID(command.OrderID) ||
		!validPreparationID(command.OrderLineID) || !validHandoverFunding(command, s.contract, tenant) || len(command.IdempotencyKey) < 8 || len(command.IdempotencyKey) > 128 {
		return empty, false, ErrInvalid
	}
	hash, err := hex.DecodeString(command.ObservationSHA256)
	if err != nil || len(hash) != 32 || strings.ToLower(command.ObservationSHA256) != command.ObservationSHA256 {
		return empty, false, ErrInvalid
	}
	encoded, err := json.Marshal(struct {
		Actor                 string
		Command               PrepareHandoverCommand
		Contract              HandoverReleaseContract
		MaximumAgeNanoseconds int64
	}{actor, command, s.contract, int64(s.contract.MaximumObservationAge)})
	if err != nil {
		return empty, false, err
	}
	sum := sha256.Sum256(encoded)
	return s.repository.PrepareInitialHandover(ctx, tenant, actor, s.ids.New(), s.ids.New(), command, s.contract, hex.EncodeToString(sum[:]))
}

// Result recovers the committed receipt without resending the mutation. Scope
// identity is server-authenticated; no caller-supplied customer or stock is used.
func (s *HandoverPreparationService) Result(ctx context.Context, tenant, organization, order, key string) (HandoverPreparation, error) {
	if s == nil || !s.contract.AllowsScope(tenant, organization) {
		return HandoverPreparation{}, ErrReleaseConditioned
	}
	if tenant == "" || !validPreparationID(organization) || !validPreparationID(order) || len(key) < 8 || len(key) > 128 {
		return HandoverPreparation{}, ErrInvalid
	}
	return s.repository.InitialHandoverResult(ctx, tenant, organization, order, key)
}

func validHandoverFunding(c PrepareHandoverCommand, p HandoverReleaseContract, tenant string) bool {
	if c.FundingReceiptID != "" {
		return c.PaymentAttemptID == "" && validPreparationID(c.FundingReceiptID) && p.AllowsStoredValueFunding(tenant, c.OrganizationID)
	}
	return validPreparationID(c.PaymentAttemptID)
}
