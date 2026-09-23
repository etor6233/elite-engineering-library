package accounting

// AUTHORED binding of an existing immutable conversion to existing accounting.
// The accountant selects accounts; this adapter selects no business/fiscal rule.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
	"unicode/utf8"
)

type FXJournalCommand struct {
	OrganizationID       string `json:"organization_id"`
	ConversionID         string `json:"conversion_id"`
	ConversionRequestKey string `json:"conversion_request_key"`
	PeriodID             string `json:"period_id"`
	DebitAccount         string `json:"debit_account"`
	CreditAccount        string `json:"credit_account"`
	Description          string `json:"description"`
	IdempotencyKey       string `json:"-"`
}
type FXJournalReceipt struct {
	JournalID      string    `json:"journal_id"`
	ConversionID   string    `json:"conversion_id"`
	OrganizationID string    `json:"organization_id"`
	RequestedBy    string    `json:"requested_by"`
	RequestKey     string    `json:"request_key"`
	PeriodID       string    `json:"period_id"`
	PostingDate    string    `json:"posting_date"`
	Currency       string    `json:"currency"`
	AmountMinor    int64     `json:"amount_minor,string"`
	DebitAccount   string    `json:"debit_account"`
	CreditAccount  string    `json:"credit_account"`
	Description    string    `json:"description"`
	RecordedAt     time.Time `json:"recorded_at"`
	Effect         string    `json:"effect"`
}
type FXJournalResult struct {
	Receipt        FXJournalReceipt `json:"receipt"`
	CurrentStatus  string           `json:"current_status"`
	CurrentVersion int64            `json:"current_version"`
}

func (s *FXService) PrepareJournal(ctx context.Context, tenant, actor string, c FXJournalCommand) (FXJournalResult, bool, error) {
	if s == nil || !s.snapshot.Allows(tenant, c.OrganizationID) || len(actor) < 1 || len(actor) > 256 || !fxKey.MatchString(c.IdempotencyKey) || !fxKey.MatchString(c.ConversionRequestKey) || len(c.ConversionID) < 1 || len(c.ConversionID) > 128 || len(c.PeriodID) < 1 || len(c.PeriodID) > 128 || !codePattern.MatchString(c.DebitAccount) || !codePattern.MatchString(c.CreditAccount) || !utf8.ValidString(c.Description) || len(c.Description) > 250 {
		return FXJournalResult{}, false, ErrInvalid
	}
	raw, err := json.Marshal(struct {
		Tenant, Actor string
		Command       FXJournalCommand
	}{tenant, actor, c})
	if err != nil {
		return FXJournalResult{}, false, ErrInvalid
	}
	hash := sha256.Sum256(raw)
	return s.repository.PrepareFXJournal(ctx, tenant, actor, s.ids.New(), s.ids.New(), s.ids.New(), c, hex.EncodeToString(hash[:]))
}
func (s *FXService) JournalResult(ctx context.Context, tenant, organization, actor, key string) (FXJournalResult, error) {
	if s == nil || !s.snapshot.Allows(tenant, organization) || len(actor) < 1 || len(actor) > 256 || !fxKey.MatchString(key) {
		return FXJournalResult{}, ErrInvalid
	}
	return s.repository.FXJournalResult(ctx, tenant, organization, actor, key)
}
