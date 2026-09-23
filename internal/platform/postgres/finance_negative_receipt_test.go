package postgres

import (
	"encoding/json"
	"testing"
)

func TestFinanceRejectedValidationHasDurableRecovery(t *testing.T) {
	ctx, p, tenant := financeFixture(t)
	store := NewFinanceStore(p, financeIDs{}, nil)
	actor := financePrincipal(tenant)
	actor.Permissions["accounting:manage"] = struct{}{}
	c := FinanceCommand{Action: "accounting_account", OrganizationID: "org", Payload: json.RawMessage(`{"code":"CASH","name":"Caja","type":"invalid"}`)}
	key := "finance-invalid-account"
	r, _, e := store.Command(ctx, actor, key, c)
	if e != nil {
		t.Fatalf("validation has no recoverable rejection: %v", e)
	}
	raw, _ := json.Marshal(r)
	var decoded map[string]any
	json.Unmarshal(raw, &decoded)
	if decoded["outcome"] != "REJECTED" {
		t.Fatalf("no explicit rejection %s", raw)
	}
	same, replay, e := store.Command(ctx, actor, key, c)
	if e != nil || !replay || same.PayloadSHA256 != r.PayloadSHA256 {
		t.Fatalf("negative replay %v", e)
	}
	recovered, e := store.CommandResult(ctx, actor, "org", key)
	if e != nil || string(recovered.Result) != string(r.Result) {
		t.Fatalf("negative recover %v", e)
	}
	var count int
	for _, q := range []string{`select count(*) from accounting.account where tenant_id=$1`, `select count(*) from platform.outbox_event where tenant_id=$1`} {
		if e = p.QueryRow(ctx, q, tenant).Scan(&count); e != nil || count != 0 {
			t.Fatalf("rejected effect %d %v", count, e)
		}
	}
	c.Payload = json.RawMessage(`{"code":"CASH","name":"Caja","type":"asset"}`)
	if _, _, e = store.Command(ctx, actor, key, c); e == nil {
		t.Fatal("changed negative key accepted")
	}
	fixed, _, e := store.Command(ctx, actor, "finance-corrected-account", c)
	if e != nil {
		t.Fatal(e)
	}
	raw, _ = json.Marshal(fixed)
	json.Unmarshal(raw, &decoded)
	if decoded["outcome"] != "CONFIRMED" {
		t.Fatalf("not corrected %s", raw)
	}
}
