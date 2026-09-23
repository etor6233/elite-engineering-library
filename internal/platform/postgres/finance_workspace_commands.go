package postgres

// AUTHORED transactional receipt coordinator. Original royalty validation,
// calculations, state/version fences and outbox SQL execute unchanged.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/royalty"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"regexp"
	"strings"
	"time"
)

var ErrFinanceNotFound = errors.New("finance receipt not found")
var financeKey = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

type FinanceStore struct {
	pool     *pgxpool.Pool
	ids      royalty.IDGenerator
	Snapshot *bcfx.Snapshot
}

func NewFinanceStore(pool *pgxpool.Pool, ids royalty.IDGenerator, snapshot *bcfx.Snapshot) *FinanceStore {
	return &FinanceStore{pool: pool, ids: ids, Snapshot: snapshot}
}

type FinanceCommand struct {
	Action         string          `json:"action"`
	OrganizationID string          `json:"organization_id"`
	Payload        json.RawMessage `json:"payload"`
}
type FinanceReceipt struct {
	Schema         string          `json:"schema"`
	Outcome        string          `json:"outcome"`
	RequestKey     string          `json:"request_key"`
	Action         string          `json:"action"`
	OrganizationID string          `json:"organization_id"`
	RequestedBy    string          `json:"requested_by"`
	PayloadSHA256  string          `json:"payload_sha256"`
	Result         json.RawMessage `json:"result"`
}

func FinancePermission(action string) string {
	return map[string]string{"royalty_policy": "royalty:policy", "royalty_accrue": "royalty:post", "royalty_open": "royalty:settle", "royalty_close": "royalty:settle", "royalty_reverse": "royalty:reverse", "royalty_reconcile": "royalty:reconcile", "accounting_account": "accounting:manage", "accounting_period": "accounting:manage", "accounting_journal": "accounting:write", "accounting_post": "accounting:post", "accounting_reverse": "accounting:reverse", "accounting_close": "accounting:close"}[action]
}
func FinanceRoyaltyAllowed(p identity.Principal) bool {
	for _, action := range []string{"royalty_policy", "royalty_accrue", "royalty_open", "royalty_reverse", "royalty_reconcile"} {
		if p.Allowed(FinancePermission(action)) {
			return true
		}
	}
	return false
}
func financeScope(p identity.Principal, org string) bool {
	return p.TenantID != "" && p.Subject != "" && len(p.Subject) <= 256 && org != "" && p.AllowedOrganization(org)
}

type financeNested struct{ tx pgx.Tx }

func (n financeNested) Query(c context.Context, q string, args ...any) (pgx.Rows, error) {
	return n.tx.Query(c, q, args...)
}
func (n financeNested) QueryRow(c context.Context, q string, args ...any) pgx.Row {
	return n.tx.QueryRow(c, q, args...)
}
func (n financeNested) Begin(c context.Context) (pgx.Tx, error) { return n.tx.Begin(c) }
func (n financeNested) BeginTx(c context.Context, _ pgx.TxOptions) (pgx.Tx, error) {
	return n.tx.Begin(c)
}
func financeHash(raw []byte) string { v := sha256.Sum256(raw); return hex.EncodeToString(v[:]) }
func prepareFinance(c FinanceCommand) ([]byte, func(context.Context, *royalty.Service, string) (any, error), error) {
	var payload any
	var call func(context.Context, *royalty.Service, string) (any, error)
	switch c.Action {
	case "royalty_policy":
		var v struct {
			AgreementID     string     `json:"agreement_id"`
			Currency        string     `json:"currency"`
			RateBasisPoints int        `json:"rate_basis_points"`
			ValidFrom       time.Time  `json:"valid_from"`
			ValidUntil      *time.Time `json:"valid_until,omitempty"`
		}
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil {
			return nil, nil, royalty.ErrInvalid
		}
		payload = v
		call = func(ctx context.Context, s *royalty.Service, tenant string) (any, error) {
			return s.CreatePolicy(ctx, tenant, royalty.Policy{AgreementID: v.AgreementID, OrganizationID: c.OrganizationID, Currency: v.Currency, RateBasisPoints: v.RateBasisPoints, ValidFrom: v.ValidFrom, ValidUntil: v.ValidUntil})
		}
	case "royalty_open":
		var v struct {
			Currency    string    `json:"currency"`
			PeriodStart time.Time `json:"period_start"`
			PeriodEnd   time.Time `json:"period_end"`
		}
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil {
			return nil, nil, royalty.ErrInvalid
		}
		payload = v
		call = func(ctx context.Context, s *royalty.Service, tenant string) (any, error) {
			return s.OpenSettlement(ctx, tenant, royalty.Settlement{OrganizationID: c.OrganizationID, Currency: v.Currency, PeriodStart: v.PeriodStart, PeriodEnd: v.PeriodEnd})
		}
	case "royalty_accrue":
		var v struct {
			PaymentAttemptID       string    `json:"payment_attempt_id"`
			PaymentExpectedVersion int64     `json:"payment_expected_version"`
			State                  string    `json:"state"`
			OccurredAt             time.Time `json:"occurred_at"`
			SourceEventKey         string    `json:"source_event_key"`
		}
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil || !financeKey.MatchString(v.SourceEventKey) {
			return nil, nil, royalty.ErrInvalid
		}
		payload = v
		call = func(ctx context.Context, s *royalty.Service, tenant string) (any, error) {
			return s.AccruePayment(ctx, tenant, v.SourceEventKey, royalty.PaymentEvent{OrganizationID: c.OrganizationID, PaymentAttemptID: v.PaymentAttemptID, PaymentExpectedVersion: v.PaymentExpectedVersion, State: v.State, OccurredAt: v.OccurredAt})
		}
	case "royalty_close", "royalty_reverse":
		var v struct {
			SettlementID    string `json:"settlement_id"`
			ExpectedVersion int64  `json:"expected_version"`
			Reason          string `json:"reason,omitempty"`
		}
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil || c.Action == "royalty_close" && v.Reason != "" {
			return nil, nil, royalty.ErrInvalid
		}
		payload = v
		call = func(ctx context.Context, s *royalty.Service, tenant string) (any, error) {
			if c.Action == "royalty_close" {
				return s.CloseSettlement(ctx, tenant, c.OrganizationID, v.SettlementID, v.ExpectedVersion)
			}
			return s.ReverseSettlement(ctx, tenant, c.OrganizationID, v.SettlementID, v.ExpectedVersion, v.Reason)
		}
	case "royalty_reconcile":
		var v struct {
			SettlementID      string `json:"settlement_id"`
			ExternalReference string `json:"external_reference"`
			ActualMinorUnits  *int64 `json:"actual_minor_units,string"`
		}
		if bcfx.DecodeExactJSON(c.Payload, &v) != nil || v.ActualMinorUnits == nil {
			return nil, nil, royalty.ErrInvalid
		}
		payload = v
		call = func(ctx context.Context, s *royalty.Service, tenant string) (any, error) {
			return s.Reconcile(ctx, tenant, c.OrganizationID, "", royalty.Reconciliation{SettlementID: v.SettlementID, ExternalReference: v.ExternalReference, ActualMinorUnits: *v.ActualMinorUnits})
		}
	default:
		return nil, nil, royalty.ErrInvalid
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
func readFinanceReceipt(ctx context.Context, q fxReader, p identity.Principal, org, key string) (FinanceReceipt, error) {
	var out FinanceReceipt
	var raw []byte
	var hash string
	e := q.QueryRow(ctx, `select receipt_raw,receipt_sha256_hex from platform.finance_command_receipt where tenant_id=$1 and organization_id=$2 and requested_by_subject=$3 and request_key=$4`, p.TenantID, org, p.Subject, key).Scan(&raw, &hash)
	if errors.Is(e, pgx.ErrNoRows) {
		return out, ErrFinanceNotFound
	}
	if e != nil {
		return out, e
	}
	if financeHash(raw) != hash || json.Unmarshal(raw, &out) != nil || out.Schema != "finance-command-receipt/v1" || (out.Outcome != "CONFIRMED" && out.Outcome != "REJECTED") || out.OrganizationID != org || out.RequestedBy != p.Subject || out.RequestKey != key {
		return FinanceReceipt{}, royalty.ErrConflict
	}
	return out, nil
}
func (s *FinanceStore) CommandResult(ctx context.Context, p identity.Principal, org, key string) (FinanceReceipt, error) {
	if !financeScope(p, org) || !financeKey.MatchString(key) {
		return FinanceReceipt{}, royalty.ErrInvalid
	}
	out, e := readFinanceReceipt(ctx, s.pool, p, org, key)
	if e == nil && !p.Allowed(FinancePermission(out.Action)) {
		return FinanceReceipt{}, royalty.ErrInvalid
	}
	return out, e
}
func (s *FinanceStore) Command(ctx context.Context, p identity.Principal, key string, c FinanceCommand) (FinanceReceipt, bool, error) {
	var empty FinanceReceipt
	permission := FinancePermission(c.Action)
	if !financeScope(p, c.OrganizationID) || permission == "" || !p.Allowed(permission) || !financeKey.MatchString(key) {
		return empty, false, royalty.ErrInvalid
	}
	var canonical []byte
	var call func(context.Context, *royalty.Service, string) (any, error)
	var accountCall func(context.Context, *accounting.Service, string, string) (any, error)
	var e error
	if strings.HasPrefix(c.Action, "accounting_") {
		canonical, accountCall, e = prepareFinanceAccounting(c)
	} else {
		canonical, call, e = prepareFinance(c)
	}
	if e != nil {
		return empty, false, e
	}
	hash := financeHash(canonical)
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if e != nil {
		return empty, false, e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, p.TenantID+":"+c.OrganizationID+":"+p.Subject+":"+key); e != nil {
		return empty, false, e
	}
	if prior, e := readFinanceReceipt(ctx, tx, p, c.OrganizationID, key); e == nil {
		if prior.Action != c.Action || prior.PayloadSHA256 != hash {
			return empty, false, royalty.ErrConflict
		}
		return prior, true, tx.Commit(ctx)
	} else if !errors.Is(e, ErrFinanceNotFound) {
		return empty, false, e
	}
	// An outer savepoint owns every domain/outbox effect. Only acknowledged rollback
	// can produce an immutable negative receipt; uncertain database failures stay fenced.
	domainTx, e := tx.Begin(ctx)
	if e != nil {
		return empty, false, e
	}
	defer domainTx.Rollback(ctx)
	outcome := "CONFIRMED"
	service := royalty.NewService(&Royalty{pool: financeNested{domainTx}}, s.ids)
	var value any
	if accountCall != nil {
		value, e = accountCall(ctx, accounting.NewService(&Accounting{pool: financeNested{domainTx}}, s.ids), p.TenantID, p.Subject)
	} else if c.Action == "royalty_reconcile" {
		var v struct {
			SettlementID      string `json:"settlement_id"`
			ExternalReference string `json:"external_reference"`
			ActualMinorUnits  *int64 `json:"actual_minor_units,string"`
		}
		_ = bcfx.DecodeExactJSON(c.Payload, &v)
		value, e = service.Reconcile(ctx, p.TenantID, c.OrganizationID, p.Subject, royalty.Reconciliation{SettlementID: v.SettlementID, ExternalReference: v.ExternalReference, ActualMinorUnits: *v.ActualMinorUnits})
	} else {
		value, e = call(ctx, service, p.TenantID)
	}
	if e != nil {
		code := ""
		if errors.Is(e, royalty.ErrInvalid) || errors.Is(e, accounting.ErrInvalid) {
			code = "INVALID_COMMAND"
		} else if errors.Is(e, royalty.ErrConflict) || errors.Is(e, accounting.ErrConflict) {
			code = "STATE_CONFLICT"
		}
		if code == "" {
			return empty, false, e
		}
		if rollbackError := domainTx.Rollback(ctx); rollbackError != nil {
			return empty, false, rollbackError
		}
		outcome = "REJECTED"
		value = map[string]any{"code": code, "business_effect": false}
	} else if e = domainTx.Commit(ctx); e != nil {
		return empty, false, e
	}
	result, e := financeResultBytes(value)
	if e != nil {
		return empty, false, e
	}
	out := FinanceReceipt{Schema: "finance-command-receipt/v1", Outcome: outcome, RequestKey: key, Action: c.Action, OrganizationID: c.OrganizationID, RequestedBy: p.Subject, PayloadSHA256: hash, Result: result}
	raw, e := json.Marshal(out)
	if e != nil {
		return empty, false, e
	}
	_, e = tx.Exec(ctx, `insert into platform.finance_command_receipt(tenant_id,organization_id,requested_by_subject,request_key,action,payload_sha256_hex,receipt_raw,receipt_sha256_hex)values($1,$2,$3,$4,$5,$6,$7,$8)`, p.TenantID, c.OrganizationID, p.Subject, key, c.Action, hash, raw, financeHash(raw))
	if e != nil {
		return empty, false, royaltyConflict(e)
	}
	if e = tx.Commit(ctx); e != nil {
		return empty, false, royaltyConflict(e)
	}
	return out, false, nil
}

// Canonical maps sort keys without floating point; Go times have already been normalized.
func financeCanonical(raw []byte) ([]byte, error) {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var value any
	if e := d.Decode(&value); e != nil {
		return nil, e
	}
	return json.Marshal(value)
}

// Monetary integers cross JavaScript as exact decimal strings, without rounding.
func financeResultBytes(value any) ([]byte, error) {
	raw, e := json.Marshal(value)
	if e != nil {
		return nil, e
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var object any
	if e = d.Decode(&object); e != nil {
		return nil, e
	}
	var walk func(any)
	walk = func(node any) {
		switch v := node.(type) {
		case map[string]any:
			for key, item := range v {
				if strings.HasSuffix(key, "minor_units") {
					if n, ok := item.(json.Number); ok {
						v[key] = n.String()
					}
				} else {
					walk(item)
				}
			}
		case []any:
			for _, item := range v {
				walk(item)
			}
		}
	}
	walk(object)
	return json.Marshal(object)
}
