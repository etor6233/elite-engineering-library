package accounting

// AUTHORED conversion-receipt glue. RecordConversion creates only a receipt.
// PrepareJournal binds that receipt to the existing separate posting lifecycle.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/bcfx"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

var ErrFXNotFound = errors.New("FX conversion receipt not found")

type FXCommand struct {
	OrganizationID string `json:"organization_id"`
	FromCurrency   string `json:"from_currency"`
	ToCurrency     string `json:"to_currency"`
	AmountMinor    *int64 `json:"amount_minor,string"`
	ConversionDate string `json:"conversion_date"`
	IdempotencyKey string `json:"-"`
}
type FXReceipt struct {
	ID          string                `json:"conversion_id"`
	RequestedBy string                `json:"requested_by"`
	RequestKey  string                `json:"request_key"`
	RecordedAt  time.Time             `json:"recorded_at"`
	Snapshot    bcfx.SnapshotIdentity `json:"snapshot"`
	Conversion  bcfx.Conversion       `json:"conversion"`
	Effect      string                `json:"effect"`
}
type FXRepository interface {
	RecordFXConversion(context.Context, string, string, string, string, FXCommand, *bcfx.Snapshot, string) (FXReceipt, bool, error)
	FXConversionResult(context.Context, string, string, string, string) (FXReceipt, error)
	PrepareFXJournal(context.Context, string, string, string, string, string, FXJournalCommand, string) (FXJournalResult, bool, error)
	FXJournalResult(context.Context, string, string, string, string) (FXJournalResult, error)
}
type FXService struct {
	repository FXRepository
	ids        IDGenerator
	snapshot   *bcfx.Snapshot
}

var fxKey = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

func NewFXService(repository FXRepository, ids IDGenerator, snapshot *bcfx.Snapshot) (*FXService, error) {
	if repository == nil || ids == nil || snapshot == nil || snapshot.Identity().ProfileSHA256 == "" {
		return nil, ErrInvalid
	}
	return &FXService{repository, ids, snapshot}, nil
}
func (s *FXService) Record(ctx context.Context, tenant, actor string, c FXCommand) (FXReceipt, bool, error) {
	if s == nil || c.AmountMinor == nil || len(actor) < 1 || len(actor) > 256 || !fxKey.MatchString(c.IdempotencyKey) || !s.snapshot.Allows(tenant, c.OrganizationID) || s.snapshot.ValidateRequest(c.FromCurrency, c.ToCurrency, c.ConversionDate) != nil {
		return FXReceipt{}, false, ErrInvalid
	}
	amount := *c.AmountMinor
	c.AmountMinor = &amount
	identity := s.snapshot.Identity()
	payload, _ := json.Marshal(struct {
		Tenant, Actor, Profile, Source string
		Command                        FXCommand
	}{tenant, actor, identity.ProfileSHA256, identity.SourceSHA256, c})
	sum := sha256.Sum256(payload)
	return s.repository.RecordFXConversion(ctx, tenant, actor, s.ids.New(), s.ids.New(), c, s.snapshot, hex.EncodeToString(sum[:]))
}
func (s *FXService) Result(ctx context.Context, tenant, organization, actor, key string) (FXReceipt, error) {
	if s == nil || !s.snapshot.Allows(tenant, organization) || len(actor) < 1 || len(actor) > 256 || !fxKey.MatchString(key) {
		return FXReceipt{}, ErrInvalid
	}
	return s.repository.FXConversionResult(ctx, tenant, organization, actor, key)
}
