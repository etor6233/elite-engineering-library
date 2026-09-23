package accounting

import (
	"context"
	"elite.local/enterprise/internal/bcamounts"
	"errors"
	"regexp"
	"time"
)

var ErrInvalid = errors.New("invalid accounting command")
var ErrConflict = errors.New("accounting conflict")
var codePattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9._-]{0,31}$`)
var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

type Account struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`
}
type Period struct {
	ID       string    `json:"id"`
	StartsOn time.Time `json:"starts_on"`
	EndsOn   time.Time `json:"ends_on"`
	Status   string    `json:"status"`
	Version  int64     `json:"version"`
}
type Line struct {
	LineNo           int    `json:"line_no"`
	AccountCode      string `json:"account_code"`
	Description      string `json:"description"`
	DebitMinorUnits  int64  `json:"debit_minor_units"`
	CreditMinorUnits int64  `json:"credit_minor_units"`
}
type Journal struct {
	ID                    string    `json:"id"`
	OrganizationID        string    `json:"organization_id"`
	PeriodID              string    `json:"period_id"`
	SourceType            string    `json:"source_type"`
	SourceID              string    `json:"source_id"`
	Currency              string    `json:"currency"`
	PostingDate           time.Time `json:"posting_date"`
	Status                string    `json:"status"`
	TotalDebitMinorUnits  int64     `json:"total_debit_minor_units"`
	TotalCreditMinorUnits int64     `json:"total_credit_minor_units"`
	Version               int64     `json:"version"`
	ReversalOf            string    `json:"reversal_of,omitempty"`
	Lines                 []Line    `json:"lines,omitempty"`
}
type Balance struct {
	AccountCode      string `json:"account_code"`
	DebitMinorUnits  int64  `json:"debit_minor_units"`
	CreditMinorUnits int64  `json:"credit_minor_units"`
	NetMinorUnits    int64  `json:"net_minor_units"`
}
type Repository interface {
	CreateAccount(context.Context, string, string, Account) error
	OpenPeriod(context.Context, string, string, Period) error
	CreateJournal(context.Context, string, string, Journal) error
	PostJournal(context.Context, string, string, string, string, int64, string, string) (Journal, error)
	ReverseJournal(context.Context, string, string, string, string, string, string, int64, string, string) (Journal, error)
	ClosePeriod(context.Context, string, string, int64, string) (Period, error)
	TrialBalance(context.Context, string, string, string) ([]Balance, error)
}
type IDGenerator interface{ New() string }
type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

func (s *Service) CreateAccount(ctx context.Context, tenant string, value Account) (Account, error) {
	validType := value.Type == "asset" || value.Type == "liability" || value.Type == "equity" || value.Type == "revenue" || value.Type == "expense"
	if tenant == "" || !codePattern.MatchString(value.Code) || len(value.Name) < 2 || len(value.Name) > 120 || !validType {
		return value, ErrInvalid
	}
	if err := s.repository.CreateAccount(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}
func (s *Service) OpenPeriod(ctx context.Context, tenant string, value Period) (Period, error) {
	if tenant == "" || value.StartsOn.IsZero() || !value.EndsOn.After(value.StartsOn) {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	value.Status = "open"
	value.Version = 1
	if err := s.repository.OpenPeriod(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}
func (s *Service) CreateJournal(ctx context.Context, tenant string, value Journal) (Journal, error) {
	// Reserved for the typed conversion binding; generic input cannot forge it.
	if value.SourceType == "FX_CONVERSION" {
		return value, ErrInvalid
	}
	if tenant == "" || value.OrganizationID == "" || value.PeriodID == "" || !codePattern.MatchString(value.SourceType) || value.SourceID == "" || !currencyPattern.MatchString(value.Currency) || value.PostingDate.IsZero() || len(value.Lines) < 2 {
		return value, ErrInvalid
	}
	seen := map[int]struct{}{}
	entries := make([]bcamounts.Entry, 0, len(value.Lines))
	for _, line := range value.Lines {
		if line.LineNo < 1 || !codePattern.MatchString(line.AccountCode) || len(line.Description) > 250 || ((line.DebitMinorUnits > 0) == (line.CreditMinorUnits > 0)) {
			return value, ErrInvalid
		}
		if _, ok := seen[line.LineNo]; ok {
			return value, ErrInvalid
		}
		seen[line.LineNo] = struct{}{}
		entries = append(entries, bcamounts.Entry{Debit: line.DebitMinorUnits, Credit: line.CreditMinorUnits})
	}
	debit, credit, amountErr := bcamounts.JournalTotals(entries)
	if amountErr != nil {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	value.Status = "draft"
	value.Version = 1
	value.TotalDebitMinorUnits = debit
	value.TotalCreditMinorUnits = credit
	if err := s.repository.CreateJournal(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}
func (s *Service) PostJournal(ctx context.Context, tenant, organization, id, actor string, version int64) (Journal, error) {
	if tenant == "" || organization == "" || id == "" || actor == "" || version < 1 {
		return Journal{}, ErrInvalid
	}
	return s.repository.PostJournal(ctx, tenant, organization, id, actor, version, s.ids.New(), s.ids.New())
}
func (s *Service) ReverseJournal(ctx context.Context, tenant, organization, id, actor string, version int64, reason string) (Journal, error) {
	if tenant == "" || organization == "" || id == "" || actor == "" || version < 1 || len(reason) < 3 || len(reason) > 500 {
		return Journal{}, ErrInvalid
	}
	return s.repository.ReverseJournal(ctx, tenant, organization, id, s.ids.New(), s.ids.New(), s.ids.New(), version, actor, reason)
}
func (s *Service) ClosePeriod(ctx context.Context, tenant, id string, version int64) (Period, error) {
	if tenant == "" || id == "" || version < 1 {
		return Period{}, ErrInvalid
	}
	return s.repository.ClosePeriod(ctx, tenant, id, version, s.ids.New())
}
func (s *Service) TrialBalance(ctx context.Context, tenant, organization, period string) ([]Balance, error) {
	if tenant == "" || organization == "" || period == "" {
		return nil, ErrInvalid
	}
	return s.repository.TrialBalance(ctx, tenant, organization, period)
}
