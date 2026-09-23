package postgres

// AUTHORED request-shape adaptation only; existing accounting owner executes all rules.
import (
	"context"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"encoding/json"
	"time"
)

func prepareFinanceAccounting(c FinanceCommand) ([]byte, func(context.Context, *accounting.Service, string, string) (any, error), error) {
	var payload any
	var call func(context.Context, *accounting.Service, string, string) (any, error)
	switch c.Action {
	case "accounting_account":
		var v accounting.Account
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil {
			return nil, nil, accounting.ErrInvalid
		}
		payload = v
		call = func(ctx context.Context, s *accounting.Service, tenant, actor string) (any, error) {
			return s.CreateAccount(ctx, tenant, v)
		}
	case "accounting_period":
		var v struct {
			StartsOn time.Time `json:"starts_on"`
			EndsOn   time.Time `json:"ends_on"`
		}
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil {
			return nil, nil, accounting.ErrInvalid
		}
		payload = v
		call = func(ctx context.Context, s *accounting.Service, tenant, actor string) (any, error) {
			return s.OpenPeriod(ctx, tenant, accounting.Period{StartsOn: v.StartsOn, EndsOn: v.EndsOn})
		}
	case "accounting_journal":
		type line struct {
			LineNo      int    `json:"line_no"`
			AccountCode string `json:"account_code"`
			Description string `json:"description"`
			Debit       *int64 `json:"debit_minor_units,string"`
			Credit      *int64 `json:"credit_minor_units,string"`
		}
		var v struct {
			PeriodID    string    `json:"period_id"`
			SourceType  string    `json:"source_type"`
			SourceID    string    `json:"source_id"`
			Currency    string    `json:"currency"`
			PostingDate time.Time `json:"posting_date"`
			Lines       []line    `json:"lines"`
		}
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil || len(v.Lines) < 2 || len(v.Lines) > 100 {
			return nil, nil, accounting.ErrInvalid
		}
		lines := make([]accounting.Line, 0, len(v.Lines))
		for _, l := range v.Lines {
			if l.Debit == nil || l.Credit == nil {
				return nil, nil, accounting.ErrInvalid
			}
			lines = append(lines, accounting.Line{LineNo: l.LineNo, AccountCode: l.AccountCode, Description: l.Description, DebitMinorUnits: *l.Debit, CreditMinorUnits: *l.Credit})
		}
		payload = v
		call = func(ctx context.Context, s *accounting.Service, tenant, actor string) (any, error) {
			return s.CreateJournal(ctx, tenant, accounting.Journal{OrganizationID: c.OrganizationID, PeriodID: v.PeriodID, SourceType: v.SourceType, SourceID: v.SourceID, Currency: v.Currency, PostingDate: v.PostingDate, Lines: lines})
		}
	case "accounting_post", "accounting_reverse":
		var v struct {
			JournalID       string `json:"journal_id"`
			ExpectedVersion int64  `json:"expected_version"`
			Reason          string `json:"reason,omitempty"`
		}
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil || c.Action == "accounting_post" && v.Reason != "" {
			return nil, nil, accounting.ErrInvalid
		}
		payload = v
		call = func(ctx context.Context, s *accounting.Service, tenant, actor string) (any, error) {
			if c.Action == "accounting_post" {
				return s.PostJournal(ctx, tenant, c.OrganizationID, v.JournalID, actor, v.ExpectedVersion)
			}
			return s.ReverseJournal(ctx, tenant, c.OrganizationID, v.JournalID, actor, v.ExpectedVersion, v.Reason)
		}
	case "accounting_close":
		var v struct {
			PeriodID        string `json:"period_id"`
			ExpectedVersion int64  `json:"expected_version"`
		}
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil {
			return nil, nil, accounting.ErrInvalid
		}
		payload = v
		call = func(ctx context.Context, s *accounting.Service, tenant, actor string) (any, error) {
			return s.ClosePeriod(ctx, tenant, v.PeriodID, v.ExpectedVersion)
		}
	default:
		return nil, nil, accounting.ErrInvalid
	}
	raw, e := json.Marshal(struct {
		Action         string `json:"action"`
		OrganizationID string `json:"organization_id"`
		Payload        any    `json:"payload"`
	}{c.Action, c.OrganizationID, payload})
	if e != nil {
		return nil, nil, e
	}
	raw, e = financeCanonical(raw)
	return raw, call, e
}
