package postgres

// AUTHORED read projections over existing finance owners; no financial calculation.
import (
	"context"
	"elite.local/enterprise/internal/bcfx"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/royalty"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type FinancePage struct {
	Schema         string            `json:"schema"`
	View           string            `json:"view"`
	OrganizationID string            `json:"organization_id"`
	Items          []json.RawMessage `json:"items"`
	NextCursor     *string           `json:"next_cursor"`
}
type FinanceFXOptions struct {
	Enabled       bool                     `json:"enabled"`
	ProfileID     string                   `json:"profile_id,omitempty"`
	SourceID      string                   `json:"source_id,omitempty"`
	LocalCurrency string                   `json:"local_currency,omitempty"`
	DateFrom      string                   `json:"date_from,omitempty"`
	DateThrough   string                   `json:"date_through,omitempty"`
	Currencies    []bcfx.CurrencyPrecision `json:"currencies"`
	Rates         []bcfx.RateRow           `json:"rates"`
}
type FinanceAccess struct {
	Schema           string           `json:"schema"`
	OrganizationID   string           `json:"organization_id"`
	OrganizationName string           `json:"organization_name"`
	Subject          string           `json:"subject"`
	Actions          []string         `json:"actions"`
	FX               FinanceFXOptions `json:"fx"`
}

func (s *FinanceStore) WorkspaceAccess(ctx context.Context, p identity.Principal, org string) (FinanceAccess, error) {
	var out FinanceAccess
	if !financeScope(p, org) || !p.Allowed("accounting:read") && !FinanceRoyaltyAllowed(p) {
		return out, royalty.ErrInvalid
	}
	out = FinanceAccess{Schema: "finance-workspace-access/v1", OrganizationID: org, Subject: p.Subject, Actions: []string{}, FX: FinanceFXOptions{Currencies: []bcfx.CurrencyPrecision{}, Rates: []bcfx.RateRow{}}}
	if e := s.pool.QueryRow(ctx, `select display_name from org.organization where tenant_id=$1 and organization_id=$2`, p.TenantID, org).Scan(&out.OrganizationName); e != nil {
		return out, e
	}
	for _, action := range []string{"royalty_policy", "royalty_accrue", "royalty_open", "royalty_close", "royalty_reverse", "royalty_reconcile"} {
		if p.Allowed(FinancePermission(action)) {
			out.Actions = append(out.Actions, action)
		}
	}
	for _, action := range []string{"accounting_account", "accounting_period", "accounting_journal", "accounting_post", "accounting_reverse", "accounting_close"} {
		if p.Allowed("accounting:read") && p.Allowed(FinancePermission(action)) {
			out.Actions = append(out.Actions, action)
		}
	}
	if p.Allowed("accounting:read") {
		out.Actions = append(out.Actions, "accounting_read")
	}
	if p.Allowed("accounting:read") && p.Allowed("accounting:write") {
		out.Actions = append(out.Actions, "fx_conversion", "fx_journal")
	}
	if s.Snapshot != nil && p.Allowed("accounting:read") && s.Snapshot.Allows(p.TenantID, org) && s.Snapshot.Current(time.Now()) {
		profile, source := s.Snapshot.Bytes()
		var v bcfx.ProfileDocument
		var r bcfx.SourceDocument
		if json.Unmarshal(profile, &v) != nil || json.Unmarshal(source, &r) != nil {
			return out, errors.New("invalid snapshot")
		}
		out.FX = FinanceFXOptions{Enabled: true, ProfileID: v.ID, SourceID: v.SourceID, LocalCurrency: v.LocalCurrency, DateFrom: v.DateFrom, DateThrough: v.DateThrough, Currencies: v.Currencies, Rates: r.Rates}
	}
	return out, nil
}
func (s *FinanceStore) WorkspacePage(ctx context.Context, p identity.Principal, org, view, cursor string, limit int) (FinancePage, error) {
	out := FinancePage{Schema: "finance-workspace-list/v1", View: view, OrganizationID: org, Items: []json.RawMessage{}}
	if !financeScope(p, org) || limit < 1 || limit > 50 || len(cursor) > 256 || strings.ContainsAny(cursor, "\r\n\x00") {
		return out, royalty.ErrInvalid
	}
	accounting := p.Allowed("accounting:read")
	royaltyRead := FinanceRoyaltyAllowed(p)
	var query string
	args := []any{p.TenantID, org, cursor, limit + 1}
	switch view {
	case "accounts":
		if !accounting {
			return out, royalty.ErrInvalid
		}
		query = `select account_code,jsonb_build_object('id',account_code,'label',display_name,'account_type',account_type,'active',active) from accounting.account where tenant_id=$1 and $2::text<>'' and account_code>$3 order by account_code limit $4`
	case "periods":
		if !accounting {
			return out, royalty.ErrInvalid
		}
		query = `select period_id,jsonb_build_object('id',period_id,'starts_on',starts_on,'ends_on',ends_on,'status',status,'version',version::text) from accounting.period where tenant_id=$1 and $2::text<>'' and period_id>$3 order by period_id limit $4`
	case "agreements":
		if !royaltyRead {
			return out, royalty.ErrInvalid
		}
		query = `select agreement_id,jsonb_build_object('id',agreement_id,'label',territory_code,'terms_version',terms_version,'status',status,'starts_on',starts_on,'ends_on',ends_on) from franchise.agreement where tenant_id=$1 and franchise_organization_id=$2 and agreement_id>$3 order by agreement_id limit $4`
	case "policies":
		if !royaltyRead {
			return out, royalty.ErrInvalid
		}
		query = `select policy_id,jsonb_build_object('id',policy_id,'agreement_id',agreement_id,'currency',currency,'rate_basis_points',rate_basis_points,'valid_from',valid_from,'valid_until',valid_until) from royalty.policy where tenant_id=$1 and franchise_organization_id=$2 and policy_id>$3 order by policy_id limit $4`
	case "settlements":
		if !royaltyRead {
			return out, royalty.ErrInvalid
		}
		query = `select settlement_id,jsonb_build_object('id',settlement_id,'currency',currency,'period_start',period_start,'period_end',period_end,'status',status,'version',version::text,'expected_minor_units',expected_minor_units::text,'reversal_of',reversal_of) from royalty.settlement_run where tenant_id=$1 and franchise_organization_id=$2 and settlement_id>$3 order by settlement_id limit $4`
	case "payments":
		if !p.Allowed("royalty:post") {
			return out, royalty.ErrInvalid
		}
		query = `select p.payment_attempt_id,jsonb_build_object('id',p.payment_attempt_id,'order_id',p.order_id,'state',p.state,'version',p.version::text,'currency',p.currency,'amount_minor_units',p.amount_minor_units::text,'updated_at',p.updated_at) from payment.payment_attempt p join sales.customer_order o on o.tenant_id=p.tenant_id and o.order_id=p.order_id where p.tenant_id=$1 and o.organization_id=$2 and p.payment_attempt_id>$3 and p.state in('captured','refunded') order by p.payment_attempt_id limit $4`
	case "conversions":
		if !accounting {
			return out, royalty.ErrInvalid
		}
		query = `select conversion_id,receipt_raw from accounting.fx_conversion_receipt where tenant_id=$1 and organization_id=$2 and conversion_id>$3 and requested_by_subject=$5 order by conversion_id limit $4`
		args = append(args, p.Subject)
	case "journals":
		if !accounting {
			return out, royalty.ErrInvalid
		}
		query = `select j.journal_id,jsonb_build_object('id',j.journal_id,'period_id',j.period_id,'currency',j.currency,'posting_date',j.posting_date,'status',j.status,'version',j.version::text,'description',coalesce((select min(l.description) from accounting.journal_line l where l.tenant_id=j.tenant_id and l.journal_id=j.journal_id),''),'total_debit_minor_units',j.total_debit_minor_units::text,'total_credit_minor_units',j.total_credit_minor_units::text,'source_type',j.source_type,'source_id',j.source_id,'requested_by',coalesce(r.requested_by_subject,''),'request_key',coalesce(r.request_key,'')) from accounting.journal j left join accounting.fx_journal_receipt r using(tenant_id,journal_id) where j.tenant_id=$1 and j.organization_id=$2 and j.journal_id>$3 order by j.journal_id limit $4`
	default:
		return out, royalty.ErrInvalid
	}
	rows, e := s.pool.Query(ctx, query, args...)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	var last string
	for rows.Next() {
		var id string
		var raw []byte
		if e = rows.Scan(&id, &raw); e != nil {
			return out, e
		}
		if len(out.Items) == limit {
			out.NextCursor = &last
			break
		}
		if len(raw) > 32768 || !json.Valid(raw) {
			return out, errors.New("invalid finance row")
		}
		if view == "conversions" {
			var key struct {
				RequestKey string `json:"request_key"`
			}
			if json.Unmarshal(raw, &key) != nil {
				return out, errors.New("invalid conversion receipt")
			}
			value, _, err := readFXReceipt(ctx, s.pool, p.TenantID, org, p.Subject, key.RequestKey)
			if err != nil {
				return out, err
			}
			raw, e = json.Marshal(value)
			if e != nil {
				return out, e
			}
		}
		out.Items = append(out.Items, json.RawMessage(raw))
		last = id
	}
	return out, rows.Err()
}

func (s *FinanceStore) JournalLines(ctx context.Context, p identity.Principal, org, id string, after, limit int) (FinancePage, error) {
	out := FinancePage{Schema: "finance-workspace-list/v1", View: "journal_lines", OrganizationID: org, Items: []json.RawMessage{}}
	if !financeScope(p, org) || !p.Allowed("accounting:read") || id == "" || len(id) > 128 || after < 0 || limit < 1 || limit > 50 {
		return out, royalty.ErrInvalid
	}
	rows, e := s.pool.Query(ctx, `select l.line_no,jsonb_build_object('line_no',l.line_no,'account_code',l.account_code,'label',a.display_name,'description',l.description,'debit_minor_units',l.debit_minor_units::text,'credit_minor_units',l.credit_minor_units::text)from accounting.journal_line l join accounting.journal j using(tenant_id,journal_id)join accounting.account a using(tenant_id,account_code)where l.tenant_id=$1 and j.organization_id=$2 and l.journal_id=$3 and l.line_no>$4 order by l.line_no limit $5`, p.TenantID, org, id, after, limit+1)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	last := 0
	for rows.Next() {
		var no int
		var raw []byte
		if e = rows.Scan(&no, &raw); e != nil {
			return out, e
		}
		if len(out.Items) == limit {
			cursor := strconv.Itoa(last)
			out.NextCursor = &cursor
			break
		}
		out.Items = append(out.Items, json.RawMessage(raw))
		last = no
	}
	return out, rows.Err()
}

// Read-only currency-preserving projection over posted ledger entries; no FX rule or conversion.
func (s *FinanceStore) TrialBalance(ctx context.Context, p identity.Principal, org, period string) ([]byte, error) {
	if !financeScope(p, org) || !p.Allowed("accounting:read") || period == "" || len(period) > 128 {
		return nil, royalty.ErrInvalid
	}
	rows, e := s.pool.Query(ctx, `select account_code,currency,sum(debit_minor_units)::text,sum(credit_minor_units)::text,(sum(debit_minor_units)-sum(credit_minor_units))::text from accounting.entry where tenant_id=$1 and organization_id=$2 and period_id=$3 group by account_code,currency order by account_code,currency limit 10001`, p.TenantID, org, period)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	values := []map[string]string{}
	for rows.Next() {
		var account, currency, debit, credit, net string
		if e = rows.Scan(&account, &currency, &debit, &credit, &net); e != nil {
			return nil, e
		}
		if len(values) >= 10000 {
			return nil, errors.New("balance result too large")
		}
		values = append(values, map[string]string{"account_code": account, "currency": currency, "debit_minor_units": debit, "credit_minor_units": credit, "net_minor_units": net})
	}
	if e = rows.Err(); e != nil {
		return nil, e
	}
	return json.Marshal(values)
}
