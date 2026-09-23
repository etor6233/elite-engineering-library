package postgres

import (
	"encoding/json"
	"testing"
)

func TestFinanceMixedCurrencyLegacyRefusesMisleadingTotal(t *testing.T) {
	ctx, p, tenant := financeFixture(t)
	store := NewFinanceStore(p, financeIDs{}, nil)
	actor := financePrincipal(tenant)
	for _, v := range []string{"accounting:manage", "accounting:write", "accounting:post"} {
		actor.Permissions[v] = struct{}{}
	}
	run := func(action, key string, payload any) map[string]any {
		t.Helper()
		raw, _ := json.Marshal(payload)
		r, _, e := store.Command(ctx, actor, key, FinanceCommand{Action: action, OrganizationID: "org", Payload: raw})
		if e != nil {
			t.Fatalf("%s: %v", action, e)
		}
		var out map[string]any
		json.Unmarshal(r.Result, &out)
		return out
	}
	run("accounting_account", "balance-cash-account", map[string]any{"code": "CASH", "name": "Caja", "type": "asset"})
	run("accounting_account", "balance-sales-account", map[string]any{"code": "SALES", "name": "Ventas", "type": "revenue"})
	period := run("accounting_period", "balance-open-period", map[string]any{"starts_on": "2026-01-01T00:00:00Z", "ends_on": "2026-02-01T00:00:00Z"})
	for _, currency := range []string{"USD", "EUR"} {
		j := run("accounting_journal", "balance-journal-"+currency, map[string]any{"period_id": period["id"], "source_type": "MANUAL", "source_id": "sale-" + currency, "currency": currency, "posting_date": "2026-01-15T00:00:00Z", "lines": []any{map[string]any{"line_no": 1, "account_code": "CASH", "description": "Venta", "debit_minor_units": "10000", "credit_minor_units": "0"}, map[string]any{"line_no": 2, "account_code": "SALES", "description": "Venta", "debit_minor_units": "0", "credit_minor_units": "10000"}}})
		run("accounting_post", "balance-posted-"+currency, map[string]any{"journal_id": j["id"], "expected_version": 1})
		if currency == "USD" {
			one, e := NewAccounting(p).TrialBalance(ctx, tenant, "org", period["id"].(string))
			if e != nil || len(one) != 2 || one[0].DebitMinorUnits != 10000 {
				t.Fatalf("single currency changed: %v %v", one, e)
			}
		}
	}
	mixed, e := store.TrialBalance(ctx, actor, "org", period["id"].(string))
	if e != nil {
		t.Fatal(e)
	}
	var values []map[string]string
	if e = json.Unmarshal(mixed, &values); e != nil || len(values) != 4 {
		t.Fatalf("projection: %s %v", mixed, e)
	}
	seen := map[string]bool{}
	for _, v := range values {
		seen[v["account_code"]+":"+v["currency"]] = true
		if v["debit_minor_units"] != "10000" && v["credit_minor_units"] != "10000" {
			t.Fatalf("mixed money: %v", v)
		}
	}
	if !seen["CASH:USD"] || !seen["CASH:EUR"] {
		t.Fatal("currencies lost")
	}
	other := actor
	other.Subject = "other"
	other.Organizations = map[string]struct{}{"foreign": {}}
	if _, e = store.TrialBalance(ctx, other, "org", period["id"].(string)); e == nil {
		t.Fatal("foreign scope read")
	}
	invalid, e := NewAccounting(p).TrialBalance(ctx, tenant, "org", period["id"].(string))
	if e == nil {
		t.Fatalf("misleading mixed-currency total returned: %v", invalid)
	}
}
