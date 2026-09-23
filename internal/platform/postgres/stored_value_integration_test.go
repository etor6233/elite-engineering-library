package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestStoredValueConnected(t *testing.T) {
	raw := os.Getenv("STORED_VALUE_DB_URL")
	if raw == "" {
		t.Skip("STORED_VALUE_DB_URL not set")
	}
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_stored_value_") {
		t.Fatal("dedicated loopback fixture required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Second)
	defer cancel()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := "11111111-1111-4111-8111-111111111111"
	org := "franchise-1"
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	root, e := filepath.Abs("../../..")
	if e != nil {
		t.Fatal(e)
	}
	read := func(p string) []byte {
		t.Helper()
		b, e := os.ReadFile(filepath.Join(root, p))
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	profile, e := sv.Load(read("deploy/stored-value/profile.reference.json"), sv.Hash(read("deploy/stored-value/profile.reference.json")))
	if e != nil {
		t.Fatal(e)
	}
	process := sv.Process{Python: `C:\Python314\python.exe`, Script: filepath.Join(root, "odoo_loyalty/run.py"), ScriptSHA256: sv.Hash(read("odoo_loyalty/run.py")), Manifest: filepath.Join(root, "odoo_loyalty/engine-lock.json"), ManifestSHA256: sv.Hash(read("odoo_loyalty/engine-lock.json"))}
	store, e := NewStoredValue(pool, profile, process)
	if e != nil {
		t.Fatal(e)
	}
	maker := identity.Principal{TenantID: tenant, Subject: "operator", Permissions: map[string]struct{}{"stored_value:request": {}, "stored_value:read": {}, "stored_value:approve": {}}, Organizations: map[string]struct{}{org: {}}}
	reviewer := maker
	reviewer.Subject = "reviewer"
	must(`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'stored-value','Fixture','Fixture')`, tenant)
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,$2,'Fixture','store')`, tenant, org)
	must(`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'m','m','Fixture','other','active')`, tenant)
	for _, v := range []string{"gift-50", "product-a"} {
		must(`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,$2,'m',$2,'Fixture','{}','active')`, tenant, v)
	}
	order := func(id, state, product string, qty, price int64) {
		t.Helper()
		must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,$2,$3,'customer',$4,'ARS',$5,1)`, tenant, id, org, state, qty*price)
		must(`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units)values($1,$2,'line',$3,$4,$5)`, tenant, id, product, qty, price)
	}
	request := func(id, op, order, program, account string, version int64) sv.Request {
		return sv.Request{OperationID: id, Operation: op, OrganizationID: org, OrderID: order, ProgramID: program, ExpectedOrderVersion: version, AccountID: account, ProfileSHA256: profile.SHA256()}
	}
	approve := func(v sv.ApprovalView) sv.ApprovalView {
		t.Helper()
		out, e := store.Decide(ctx, reviewer, sv.Decision{OrganizationID: org, ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: true, Reason: "Fixture reviewed exact snapshot"})
		if e != nil {
			t.Fatal(e)
		}
		return out
	}
	order("gift-source", "confirmed", "gift-50", 2, 5000)
	giftRequest := request("gift-issue", "issue", "gift-source", "reference-gift_card", "", 1)
	proposal, e := store.Propose(ctx, maker, giftRequest)
	if e != nil {

		tx, ee := pool.Begin(ctx)
		if ee == nil {
			oo, st, aa, oe := store.order(ctx, tx, giftRequest)
			t.Logf("fixture diagnostics: valid=%v allowed=%v order=%+v state=%s allocation=%+v ordererr=%v", giftRequest.Validate(profile), store.allowed(maker, "stored_value:request"), oo, st, aa, oe)
			bb, be := store.calculate(ctx, tx, giftRequest, maker.Subject, time.Now().UTC().Add(time.Hour).Format(time.RFC3339Nano), false)
			t.Logf("calculation result=%+v entries=%+v err=%v", bb.Result, bb.Entries, be)
			tx.Rollback(ctx)
		}
		t.Fatal("issue proposal", e)
	}
	if len(proposal.Payload.Entries) != 2 || proposal.Payload.Entries[0].PointsDelta != "50.000000" {
		t.Fatal(proposal)
	}
	wrong := reviewer
	wrong.TenantID = "22222222-2222-4222-8222-222222222222"
	if _, e = store.Decide(ctx, wrong, sv.Decision{OrganizationID: org, ApprovalID: proposal.ApprovalID, PayloadSHA256: proposal.PayloadSHA256, Approved: true, Reason: "wrong"}); e == nil {
		t.Fatal("cross tenant approved")
	}
	if _, e = store.Decide(ctx, maker, sv.Decision{OrganizationID: org, ApprovalID: proposal.ApprovalID, PayloadSHA256: proposal.PayloadSHA256, Approved: true, Reason: "self"}); !errors.Is(e, approval.ErrSeparation) {
		t.Fatal("separation", e)
	}
	approved := approve(proposal)
	if approved.State != "approved" || string(approved.Receipt) == "null" {
		t.Fatal(approved)
	}
	replay, e := store.Propose(ctx, maker, giftRequest)
	if e != nil || !replay.Replay || replay.PayloadSHA256 != proposal.PayloadSHA256 {
		t.Fatal("replay", replay, e)
	}
	different := giftRequest
	different.OperationID = "gift-again"
	if _, e = store.Propose(ctx, maker, different); e == nil {
		t.Fatal("duplicate issuance")
	}
	account := proposal.Payload.Entries[0].AccountID
	other := proposal.Payload.Entries[1].AccountID
	if _, e = store.Propose(ctx, maker, request("same-order", "redeem", "gift-source", "reference-gift_card", account, 1)); e == nil {
		t.Fatal("future card on issuing order")
	}
	order("purchase-a", "placed", "product-a", 1, 11500)
	order("purchase-b", "placed", "product-a", 1, 11500)
	var wg sync.WaitGroup
	result := make(chan sv.ApprovalView, 2)
	errs := make(chan error, 2)
	for _, id := range []string{"a", "b"} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			v, e := store.Propose(ctx, maker, request("spend-"+id, "redeem", "purchase-"+id, "reference-gift_card", account, 1))
			if e == nil {
				result <- v
			} else {
				errs <- e
			}
		}(id)
	}
	wg.Wait()
	close(result)
	close(errs)
	if len(result) != 1 || len(errs) != 1 {
		t.Fatal("double reservation", len(result), len(errs))
	}
	redeemed := approve(<-result)
	a, e := store.Allocation(ctx, maker, redeemed.Payload.Request.OrderID)
	if e != nil || a.GrossMinor != 11500 || a.GiftMinor != 5000 || a.ProviderDueMinor != 6500 || a.OrderVersion != 2 {
		t.Fatal(a, e)
	}
	// Rejecting a separate reservation returns availability without deleting history.
	order("small", "placed", "product-a", 1, 575)
	v, e := store.Propose(ctx, maker, request("spend-small-reject", "redeem", "small", "reference-gift_card", other, 1))
	if e != nil {
		t.Fatal(e)
	}
	_, e = store.Decide(ctx, reviewer, sv.Decision{OrganizationID: org, ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: false, Reason: "Fixture reject"})
	if e != nil {
		t.Fatal(e)
	}
	v, e = store.Propose(ctx, maker, request("spend-small", "redeem", "small", "reference-gift_card", other, 1))
	if e != nil {
		t.Fatal(e)
	}
	approve(v)
	a, e = store.Allocation(ctx, maker, "small")
	if e != nil || a.ProviderDueMinor != 0 || a.GiftMinor != 575 {
		t.Fatal(a, e)
	}
	// Source cancellation is an immutable inverse, including negative balance
	// when a previously issued card was already consumed elsewhere.
	must(`update sales.customer_order set state='cancelled',version=version+1 where tenant_id=$1 and order_id='gift-source'`, tenant)
	reverse := request("gift-reverse", "reverse", "gift-source", "reference-gift_card", "", 2)
	reverse.OriginalOperationID = "gift-issue"
	v, e = store.Propose(ctx, maker, reverse)
	if e != nil {
		t.Fatal("reverse proposal", e)
	}
	approve(v)
	var balance string
	if e = pool.QueryRow(ctx, `select sum(points_delta)::numeric(38,6)::text from stored_value.entry where tenant_id=$1 and account_id=$2`, tenant, account).Scan(&balance); e != nil || balance != "-50.000000" {
		t.Fatal(balance, e)
	}
	reverse.OperationID = "gift-reverse-again"
	if v, e = store.Propose(ctx, maker, reverse); e == nil {
		if _, e = store.Decide(ctx, reviewer, sv.Decision{OrganizationID: org, ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: true, Reason: "duplicate"}); e == nil {
			t.Fatal("double reversal")
		}
	}
	// Loyalty is bound to the customer and uses the admitted source default
	// money earning plus 5-percent/200-point reward, including minor rounding.
	order("earn", "confirmed", "product-a", 1, 20000)
	v, e = store.Propose(ctx, maker, request("loyalty-earn", "accrue", "earn", "reference-loyalty", "", 1))
	if e != nil {
		t.Fatal(e)
	}
	approve(v)
	loyalty := v.Payload.Entries[0].AccountID
	order("discount", "placed", "product-a", 1, 1001)
	v, e = store.Propose(ctx, maker, request("loyalty-spend", "redeem", "discount", "reference-loyalty", loyalty, 1))
	if e != nil {
		t.Fatal(e)
	}
	approve(v)
	a, e = store.Allocation(ctx, maker, "discount")
	if e != nil || a.DiscountMinor != 50 || a.ProviderDueMinor != 951 {
		t.Fatal(a, e)
	}

	// Seed a short-lived, correctly hash-bound proposal through the same typed
	// owner. Only the fixture lease is shorter than the deployment profile's
	// 60-second minimum; no production clock or validation is replaced.
	pending := func(id string, lifetime time.Duration) (sv.ApprovalView, time.Time) {
		t.Helper()
		order(id, "confirmed", "gift-50", 1, 5000)
		r := request(id, "issue", id, "reference-gift_card", "", 1)
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		defer tx.Rollback(ctx)
		var expiry time.Time
		if e = tx.QueryRow(ctx, `select clock_timestamp()+($1::bigint*interval '1 millisecond')`, lifetime.Milliseconds()).Scan(&expiry); e != nil {
			t.Fatal(e)
		}
		bound, e := store.calculate(ctx, tx, r, maker.Subject, expiry.UTC().Format(time.RFC3339Nano), false)
		if e != nil {
			t.Fatal(e)
		}
		raw, hash, e := sv.Canonical(bound)
		if e != nil {
			t.Fatal(e)
		}
		id = sv.StableID("svapproval", id)
		_, e = store.approvals.submitTx(ctx, tx, maker, HumanApprovalSpec{Request: approval.Request{TenantID: tenant, ID: id, Kind: approval.KindStoredValueOperation, SubjectID: r.OrderID, Requester: maker.Subject, EvidenceSHA: hash}, OrganizationID: org, Payload: raw}, "stored_value:request", nil)
		if e != nil {
			t.Fatal(e)
		}
		if e = tx.Commit(ctx); e != nil {
			t.Fatal(e)
		}
		return sv.ApprovalView{ApprovalID: id, PayloadSHA256: hash, Payload: bound}, expiry
	}
	t.Run("expiry_during_outbox_wait", func(t *testing.T) {
		v, expiry := pending("expiry-outbox", 2*time.Second)
		blocker, e := pool.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		defer blocker.Rollback(ctx)
		if _, e = blocker.Exec(ctx, `lock table platform.outbox_event in access exclusive mode`); e != nil {
			t.Fatal(e)
		}
		done := make(chan error, 1)
		go func() {
			_, e := store.Decide(ctx, reviewer, sv.Decision{OrganizationID: org, ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: true, Reason: "Exact short-lease fixture"})
			done <- e
		}()
		deadline := time.Now().Add(time.Second)
		for {
			var waiting bool
			if e = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like 'insert into platform.outbox_event%')`).Scan(&waiting); e != nil {
				t.Fatal(e)
			}
			if waiting {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("decision never reached blocked outbox")
			}
			time.Sleep(10 * time.Millisecond)
		}
		must(`select pg_sleep(greatest(0,extract(epoch from ($1::timestamptz-clock_timestamp())))+0.05)`, expiry)
		if e = blocker.Rollback(ctx); e != nil {
			t.Fatal(e)
		}
		if e = <-done; e == nil {
			t.Fatal("expired operation committed after outbox wait")
		}
		var operations, decisions, events int
		if e = pool.QueryRow(ctx, `select (select count(*) from stored_value.operation where tenant_id=$1 and approval_id=$2),(select count(*) from approval.decision where tenant_id=$1 and request_id=$2),(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='expiry-outbox')`, tenant, v.ApprovalID).Scan(&operations, &decisions, &events); e != nil || operations+decisions+events != 0 {
			t.Fatal("expired transaction left effects", operations, decisions, events, e)
		}
		got, e := store.Read(ctx, maker, v.ApprovalID)
		if e != nil || got.State != "pending" || string(got.Receipt) != "null" {
			t.Fatal(got, e)
		}
	})

	t.Run("proposal_expired_before_commit", func(t *testing.T) {
		order("expiry-proposal", "confirmed", "gift-50", 1, 5000)
		r := request("expiry-proposal", "issue", "expiry-proposal", "reference-gift_card", "", 1)
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		defer tx.Rollback(ctx)
		var expiry time.Time
		if e = tx.QueryRow(ctx, `select clock_timestamp()+interval '1 second'`).Scan(&expiry); e != nil {
			t.Fatal(e)
		}
		bound, e := store.calculate(ctx, tx, r, maker.Subject, expiry.UTC().Format(time.RFC3339Nano), false)
		if e != nil {
			t.Fatal(e)
		}
		raw, hash, e := sv.Canonical(bound)
		if e != nil {
			t.Fatal(e)
		}
		id := sv.StableID("svapproval", r.OperationID)
		_, e = store.approvals.submitTx(ctx, tx, maker, HumanApprovalSpec{Request: approval.Request{TenantID: tenant, ID: id, Kind: approval.KindStoredValueOperation, SubjectID: r.OrderID, Requester: maker.Subject, EvidenceSHA: hash}, OrganizationID: org, Payload: raw}, "stored_value:request", nil)
		if e != nil {
			t.Fatal(e)
		}
		if _, e = tx.Exec(ctx, `select pg_sleep(greatest(0,extract(epoch from ($1::timestamptz-clock_timestamp())))+0.05)`, expiry); e != nil {
			t.Fatal(e)
		}
		if e = tx.Commit(ctx); e == nil {
			t.Fatal("expired proposal committed")
		}
		if _, e = store.Read(ctx, maker, id); !errors.Is(e, approval.ErrNotFound) {
			t.Fatal("expired proposal retained", e)
		}
	})
	// Immutable register, one event per applied operation, no second ledger.
	if _, e = pool.Exec(ctx, `update stored_value.entry set points_delta=0 where tenant_id=$1`, tenant); e == nil {
		t.Fatal("mutable history")
	}
	var counts []byte
	e = pool.QueryRow(ctx, `select jsonb_build_object('operations',(select count(*) from stored_value.operation),'events',(select count(*) from platform.outbox_event where event_type='stored-value.committed'),'approvals',(select count(*) from approval.request where kind='stored_value_operation'),'reservations',(select count(*) from stored_value.reservation))`).Scan(&counts)
	if e != nil {
		t.Fatal(e)
	}
	var numbers map[string]int
	if json.Unmarshal(counts, &numbers) != nil || numbers["operations"] != numbers["events"] {
		t.Fatal(string(counts))
	}
	t.Log(string(counts))
}
