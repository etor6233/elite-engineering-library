package postgres

import (
	"encoding/json"
	"testing"
)

func TestFinanceAccountingOwnerLifecycleRecovery(t *testing.T) {
	ctx, p, tenant := financeFixture(t)
	store := NewFinanceStore(p, financeIDs{}, nil)
	actor := financePrincipal(tenant)
	for _, v := range []string{"accounting:manage", "accounting:write", "accounting:post", "accounting:reverse", "accounting:close"} {
		actor.Permissions[v] = struct{}{}
	}
	n := 0
	run := func(action, payload string) map[string]any {
		t.Helper()
		n++
		key := fmtFinanceKey(n)
		c := FinanceCommand{Action: action, OrganizationID: "org", Payload: json.RawMessage(payload)}
		a, replay, e := store.Command(ctx, actor, key, c)
		if e != nil || replay {
			t.Fatalf("%s: %v replay %v", action, e, replay)
		}
		b, replay, e := store.Command(ctx, actor, key, c)
		if e != nil || !replay || string(a.Result) != string(b.Result) {
			t.Fatalf("%s replay: %v", action, e)
		}
		var out map[string]any
		if e = json.Unmarshal(a.Result, &out); e != nil {
			t.Fatal(e)
		}
		return out
	}
	run("accounting_account", `{"code":"CASH","name":"Caja","type":"asset"}`)
	run("accounting_account", `{"code":"SALES","name":"Ventas","type":"revenue"}`)
	period := run("accounting_period", `{"starts_on":"2026-01-01T00:00:00Z","ends_on":"2026-02-01T00:00:00Z"}`)
	body := map[string]any{"period_id": period["id"], "source_type": "MANUAL", "source_id": "fixture-transaction-1", "currency": "USD", "posting_date": "2026-01-15T00:00:00Z", "lines": []any{map[string]any{"line_no": 1, "account_code": "CASH", "description": "Fixture", "debit_minor_units": "10000", "credit_minor_units": "0"}, map[string]any{"line_no": 2, "account_code": "SALES", "description": "Fixture", "debit_minor_units": "0", "credit_minor_units": "10000"}}}
	raw, _ := json.Marshal(body)
	draft := run("accounting_journal", string(raw))
	raw, _ = json.Marshal(map[string]any{"journal_id": draft["id"], "expected_version": 1})
	posted := run("accounting_post", string(raw))
	if posted["status"] != "posted" || posted["total_debit_minor_units"] != "10000" {
		t.Fatalf("not exact posted %v", posted)
	}
	raw, _ = json.Marshal(map[string]any{"journal_id": draft["id"], "expected_version": 2, "reason": "Corrección fixture"})
	reversed := run("accounting_reverse", string(raw))
	if reversed["reversal_of"] != draft["id"] {
		t.Fatal("no reversal binding")
	}
	raw, _ = json.Marshal(map[string]any{"period_id": period["id"], "expected_version": 1})
	closed := run("accounting_close", string(raw))
	if closed["status"] != "closed" {
		t.Fatal("period not closed")
	}
	lines, e := store.JournalLines(ctx, actor, "org", draft["id"].(string), 0, 1)
	if e != nil || len(lines.Items) != 1 || lines.NextCursor == nil {
		t.Fatalf("lines page %v %v", lines, e)
	}
	next, e := store.JournalLines(ctx, actor, "org", draft["id"].(string), 1, 1)
	if e != nil || len(next.Items) != 1 || next.NextCursor != nil {
		t.Fatalf("next line %v %v", next, e)
	}
	bal, e := store.TrialBalance(ctx, actor, "org", period["id"].(string))
	if e != nil || !json.Valid(bal) {
		t.Fatalf("balance %s %v", bal, e)
	}
	var nEntries, nReceipts int
	var balance int64
	if e := p.QueryRow(ctx, `select count(*),sum(debit_minor_units-credit_minor_units)::bigint from accounting.entry where tenant_id=$1`, tenant).Scan(&nEntries, &balance); e != nil || nEntries != 4 || balance != 0 {
		t.Fatalf("ledger effects %d %d %v", nEntries, balance, e)
	}
	if e := p.QueryRow(ctx, `select count(*) from platform.finance_command_receipt where tenant_id=$1`, tenant).Scan(&nReceipts); e != nil || nReceipts != 7 {
		t.Fatalf("receipts %d %v", nReceipts, e)
	}
}
func fmtFinanceKey(n int) string {
	return []string{"", "finance-account-cash", "finance-account-sales", "finance-open-period", "finance-create-journal", "finance-post-journal", "finance-reverse-journal", "finance-close-period"}[n]
}
func TestFinanceRoyaltyAllCommandsRecoverable(t *testing.T) {
	ctx, p, tenant := financeFixture(t)
	store := NewFinanceStore(p, financeIDs{}, nil)
	actor := financePrincipal(tenant)
	for _, q := range []string{`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','org','customer','confirmed','USD',10000,1)`, `insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values($1,'payment','order','fixture','fixture-payment','payment-key','captured','USD',10000,1)`} {
		if _, e := p.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	run := func(action, key, payload string) map[string]any {
		t.Helper()
		c := FinanceCommand{Action: action, OrganizationID: "org", Payload: json.RawMessage(payload)}
		a, _, e := store.Command(ctx, actor, key, c)
		if e != nil {
			t.Fatalf("%s %v", action, e)
		}
		b, replay, e := store.Command(ctx, actor, key, c)
		if e != nil || !replay || string(a.Result) != string(b.Result) {
			t.Fatalf("replay %s %v", action, e)
		}
		var value map[string]any
		json.Unmarshal(a.Result, &value)
		return value
	}
	run("royalty_policy", "royalty-all-policy", `{"agreement_id":"agreement","currency":"USD","rate_basis_points":500,"valid_from":"2026-01-01T00:00:00Z"}`)
	accrual := run("royalty_accrue", "royalty-all-accrue", `{"payment_attempt_id":"payment","payment_expected_version":1,"state":"captured","occurred_at":"2026-01-15T00:00:00Z","source_event_key":"payment-event-key-001"}`)
	if accrual["royalty_minor_units"] != "500" {
		t.Fatal("royalty amount changed")
	}
	settlement := run("royalty_open", "royalty-all-open-key", `{"currency":"USD","period_start":"2026-01-01T00:00:00Z","period_end":"2026-02-01T00:00:00Z"}`)
	raw, _ := json.Marshal(map[string]any{"settlement_id": settlement["id"], "expected_version": 1})
	closed := run("royalty_close", "royalty-all-close-key", string(raw))
	if closed["expected_minor_units"] != "500" {
		t.Fatal("settlement amount changed")
	}
	raw, _ = json.Marshal(map[string]any{"settlement_id": settlement["id"], "expected_version": 2, "reason": "Corrección de prueba"})
	reverse := run("royalty_reverse", "royalty-all-reverse-key", string(raw))
	raw, _ = json.Marshal(map[string]any{"settlement_id": reverse["id"], "external_reference": "bank-fixture-001", "actual_minor_units": "-500"})
	recon := run("royalty_reconcile", "royalty-all-reconcile-key", string(raw))
	if recon["status"] != "matched" || recon["difference_minor_units"] != "0" {
		t.Fatalf("reconciliation changed %v", recon)
	}
	var count int
	if e := p.QueryRow(ctx, `select count(*) from platform.finance_command_receipt where tenant_id=$1`, tenant).Scan(&count); e != nil || count != 6 {
		t.Fatalf("receipt count %d %v", count, e)
	}
}
