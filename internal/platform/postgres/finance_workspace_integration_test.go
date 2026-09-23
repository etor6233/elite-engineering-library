package postgres

import (
	"bytes"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/royalty"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
)

func financePrincipal(tenant string) identity.Principal {
	return identity.Principal{TenantID: tenant, Subject: "accountant", Organizations: map[string]struct{}{"org": {}}, Permissions: map[string]struct{}{"royalty:policy": {}, "royalty:settle": {}, "royalty:reverse": {}, "royalty:reconcile": {}, "royalty:post": {}, "accounting:read": {}}}
}
func TestFinanceRoyaltyDurableReplayPostgres(t *testing.T) {
	ctx, p, tenant := financeFixture(t)
	store := NewFinanceStore(p, financeIDs{}, nil)
	actor := financePrincipal(tenant)
	policy := FinanceCommand{Action: "royalty_policy", OrganizationID: "org", Payload: json.RawMessage(`{"agreement_id":"agreement","currency":"USD","rate_basis_points":500,"valid_from":"2026-01-01T00:00:00Z"}`)}
	a, replay, e := store.Command(ctx, actor, "finance-policy-key-0001", policy)
	if e != nil || replay {
		t.Fatalf("create %v replay=%v", e, replay)
	}
	reordered := policy
	reordered.Payload = json.RawMessage(`{"valid_from":"2026-01-01T00:00:00Z","rate_basis_points":500,"currency":"USD","agreement_id":"agreement"}`)
	b, replay, e := store.Command(ctx, actor, a.RequestKey, reordered)
	if e != nil || !replay || !bytes.Equal(a.Result, b.Result) {
		t.Fatalf("exact canonical replay %v %v", e, replay)
	}
	got, e := store.CommandResult(ctx, actor, "org", a.RequestKey)
	if e != nil || !bytes.Equal(got.Result, a.Result) {
		t.Fatalf("GET recover %v", e)
	}
	changed := policy
	changed.Payload = json.RawMessage(`{"agreement_id":"agreement","currency":"USD","rate_basis_points":600,"valid_from":"2026-01-01T00:00:00Z"}`)
	if _, _, e = store.Command(ctx, actor, a.RequestKey, changed); !errors.Is(e, royalty.ErrConflict) {
		t.Fatalf("changed body accepted %v", e)
	}
	open := FinanceCommand{Action: "royalty_open", OrganizationID: "org", Payload: json.RawMessage(`{"currency":"USD","period_start":"2026-01-01T00:00:00Z","period_end":"2026-02-01T00:00:00Z"}`)}
	x, replay, e := store.Command(ctx, actor, "finance-settlement-key-01", open)
	if e != nil || replay {
		t.Fatal(e)
	}
	y, replay, e := store.Command(ctx, actor, x.RequestKey, open)
	if e != nil || !replay || !bytes.Equal(x.Result, y.Result) {
		t.Fatalf("open replay %v", e)
	}
	other := actor
	other.Subject = "other-accountant"
	if _, e = store.CommandResult(ctx, other, "org", a.RequestKey); !errors.Is(e, ErrFinanceNotFound) {
		t.Fatalf("actor leak %v", e)
	}
	other = actor
	other.TenantID = "018f4d4a-7b36-7a21-8d10-2f4c54c28903"
	if _, e = store.CommandResult(ctx, other, "org", a.RequestKey); !errors.Is(e, ErrFinanceNotFound) {
		t.Fatalf("tenant leak %v", e)
	}
	if _, e = store.CommandResult(ctx, actor, "other", a.RequestKey); !errors.Is(e, royalty.ErrInvalid) {
		t.Fatalf("organization leak %v", e)
	}
	var policies, settlements, events, receipts int
	for i, q := range []string{`select count(*) from royalty.policy where tenant_id=$1`, `select count(*) from royalty.settlement_run where tenant_id=$1`, `select count(*) from platform.outbox_event where tenant_id=$1`, `select count(*) from platform.finance_command_receipt where tenant_id=$1`} {
		dest := []*int{&policies, &settlements, &events, &receipts}
		if e = p.QueryRow(ctx, q, tenant).Scan(dest[i]); e != nil {
			t.Fatal(e)
		}
	}
	if policies != 1 || settlements != 1 || events != 2 || receipts != 2 {
		t.Fatalf("duplicate effect %d/%d/%d/%d", policies, settlements, events, receipts)
	}
}
func TestFinanceReceiptFailureRollsBackDomainAndOutbox(t *testing.T) {
	ctx, p, tenant := financeFixture(t)
	_, e := p.Exec(ctx, `create function royalty.test_fail_finance_receipt()returns trigger language plpgsql as $$begin if NEW.request_key='finance-injected-failure' then raise exception 'receipt injected failure';end if;return NEW;end;$$;create trigger test_finance_receipt_fail before insert on platform.finance_command_receipt for each row execute function royalty.test_fail_finance_receipt()`)
	if e != nil {
		t.Fatal(e)
	}
	defer p.Exec(ctx, `drop trigger test_finance_receipt_fail on platform.finance_command_receipt;drop function royalty.test_fail_finance_receipt()`)
	store := NewFinanceStore(p, financeIDs{}, nil)
	c := FinanceCommand{Action: "royalty_policy", OrganizationID: "org", Payload: json.RawMessage(`{"agreement_id":"agreement","currency":"EUR","rate_basis_points":700,"valid_from":"2026-01-01T00:00:00Z"}`)}
	if _, _, e = store.Command(ctx, financePrincipal(tenant), "finance-injected-failure", c); e == nil {
		t.Fatal("injected failure not triggered")
	}
	for _, q := range []string{`select count(*) from royalty.policy where tenant_id=$1`, `select count(*) from platform.outbox_event where tenant_id=$1`, `select count(*) from platform.finance_command_receipt where tenant_id=$1`} {
		var count int
		if e = p.QueryRow(ctx, q, tenant).Scan(&count); e != nil || count != 0 {
			t.Fatalf("partial transaction %d %v", count, e)
		}
	}
}
func TestFinanceConcurrentRequestRecovery(t *testing.T) {
	ctx, p, tenant := financeFixture(t)
	store := NewFinanceStore(p, financeIDs{}, nil)
	actor := financePrincipal(tenant)
	c := FinanceCommand{Action: "royalty_open", OrganizationID: "org", Payload: json.RawMessage(`{"currency":"USD","period_start":"2026-03-01T00:00:00Z","period_end":"2026-04-01T00:00:00Z"}`)}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, _, _ = store.Command(ctx, actor, "finance-concurrency-key-1", c) }()
	}
	wg.Wait()
	first, e := store.CommandResult(ctx, actor, "org", "finance-concurrency-key-1")
	if e != nil {
		t.Fatal(e)
	}
	retry, replayed, e := store.Command(ctx, actor, first.RequestKey, c)
	if e != nil || !replayed || !bytes.Equal(first.Result, retry.Result) {
		t.Fatalf("eventual exact replay %v", e)
	}
	var n int
	if e = p.QueryRow(ctx, `select count(*) from royalty.settlement_run where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 1 {
		t.Fatalf("concurrent duplicate %d %v", n, e)
	}
}
func TestFinanceAuthorizedPagesAndScope(t *testing.T) {
	ctx, p, tenant := financeFixture(t)
	store := NewFinanceStore(p, financeIDs{}, nil)
	actor := financePrincipal(tenant)
	for i := 0; i < 3; i++ {
		_, e := p.Exec(ctx, `insert into accounting.account(tenant_id,account_code,display_name,account_type)values($1,$2,$3,'asset')`, tenant, fmt.Sprintf("ASSET%d", i), fmt.Sprintf("Cuenta %d", i))
		if e != nil {
			t.Fatal(e)
		}
	}
	a, e := store.WorkspaceAccess(ctx, actor, "org")
	if e != nil || a.OrganizationName != "Local centro" || a.FX.Enabled {
		t.Fatalf("access %v %+v", e, a)
	}
	first, e := store.WorkspacePage(ctx, actor, "org", "accounts", "", 2)
	if e != nil || len(first.Items) != 2 || first.NextCursor == nil {
		t.Fatalf("page1 %v %+v", e, first)
	}
	second, e := store.WorkspacePage(ctx, actor, "org", "accounts", *first.NextCursor, 2)
	if e != nil || len(second.Items) != 1 || second.NextCursor != nil {
		t.Fatalf("page2 %v %+v", e, second)
	}
	for _, view := range []string{"periods", "agreements", "policies", "settlements", "payments", "conversions", "journals"} {
		if _, e = store.WorkspacePage(ctx, actor, "org", view, "", 25); e != nil {
			t.Fatalf("%s %v", view, e)
		}
	}
	actor.Permissions = map[string]struct{}{"royalty:settle": {}}
	if _, e = store.WorkspacePage(ctx, actor, "org", "accounts", "", 25); !errors.Is(e, royalty.ErrInvalid) {
		t.Fatal("accounting permission leaked")
	}
	if _, e = store.WorkspacePage(ctx, actor, "other", "settlements", "", 25); !errors.Is(e, royalty.ErrInvalid) {
		t.Fatal("organization leaked")
	}
	if _, e = store.WorkspacePage(ctx, actor, "org", "payments", "", 25); !errors.Is(e, royalty.ErrInvalid) {
		t.Fatal("payment grant inferred")
	}
}
