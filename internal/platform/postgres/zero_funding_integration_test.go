package postgres_test

import (
	"context"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestZeroProviderFundingReceiptThroughDelivery(t *testing.T) {
	ctx := context.Background()
	pool := connectedPool(t)
	tenant := connectedSeedOrder(t, pool, "stripe")
	store, actor, account := fixtureStoredValue(t, pool, tenant, "gift_card", 200000)
	allocation, e := store.Allocation(ctx, actor, "order")
	if e != nil || allocation.ProviderDueMinor != 0 || allocation.GiftMinor != 123456 {
		t.Fatal(allocation, e)
	}
	sales := commerce.NewService(db.NewCommerce(pool), randomid.Generator{})
	if _, e = sales.RequestOrderPayment(ctx, tenant, "store", "order", "stripe", "no-zero-provider-payment", "operator"); e == nil {
		t.Fatal("zero provider payment created")
	}
	// No provider setup can be used by this path, even though the profile records
	// the provider selection for other orders in the same local deployment.
	if _, e = pool.Exec(ctx, `update integration.provider_connection set state='disabled' where tenant_id=$1`, tenant); e != nil {
		t.Fatal(e)
	}
	command := db.FinalizeStoredValueFunding{OrganizationID: "store", OrderID: "order", ExpectedOrderVersion: allocation.OrderVersion, RequestKey: "full-funding-observation-0001"}
	denied := actor
	denied.Permissions = map[string]struct{}{"stored_value:read": {}}
	if _, _, e = store.FinalizeFunding(ctx, denied, command); e == nil {
		t.Fatal("missing funding permission accepted")
	}
	var wg sync.WaitGroup
	results := make(chan db.LocalFundingResult, 12)
	errors := make(chan error, 12)
	for range 12 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, _, e := store.FinalizeFunding(ctx, actor, command)
			if e != nil {
				errors <- e
			} else {
				results <- r
			}
		}()
	}
	wg.Wait()
	close(results)
	close(errors)
	if len(errors) != 0 || len(results) != 12 {
		t.Fatal("concurrent funding", len(results), len(errors), <-errors)
	}
	var funding db.LocalFundingResult
	for r := range results {
		if funding.SHA256 != "" && r.SHA256 != funding.SHA256 {
			t.Fatal("duplicate observations")
		}
		funding = r
	}
	if funding.Receipt.Allocation.ProviderMinor != 0 || funding.Receipt.Effect != db.LocalFundingEffect {
		t.Fatal(funding)
	}
	got, e := store.FundingResult(ctx, actor, command.RequestKey)
	if e != nil || got.SHA256 != funding.SHA256 {
		t.Fatal("lost response", got, e)
	}
	changed := command
	changed.ExpectedOrderVersion++
	if _, _, e = store.FinalizeFunding(ctx, actor, changed); e == nil {
		t.Fatal("changed request replay")
	}
	stranger := actor
	stranger.Subject = "other-actor"
	if _, _, e = store.FinalizeFunding(ctx, stranger, command); e == nil {
		t.Fatal("other actor adopted receipt")
	}
	// A new observation stalled on outbox must roll back its receipt atomically.
	blocker, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = blocker.Exec(ctx, `lock table platform.outbox_event in access exclusive mode`); e != nil {
		t.Fatal(e)
	}
	short, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	blocked := command
	blocked.RequestKey = "full-funding-cancelled-0001"
	_, _, e = store.FinalizeFunding(short, actor, blocked)
	cancel()
	if e == nil {
		t.Fatal("blocked funding committed")
	}
	if e = blocker.Rollback(ctx); e != nil {
		t.Fatal(e)
	}
	var count int
	if e = pool.QueryRow(ctx, `select count(*) from payment.local_funding_receipt where tenant_id=$1`, tenant).Scan(&count); e != nil || count != 1 {
		t.Fatal("orphan funding", count, e)
	}
	r := &connectedRun{pool: pool, tenant: tenant}
	old, _ := connectedCommercialProfile(t, r, false)
	policy, _ := connectedCommercialProfile(t, r, true)
	repo := db.NewFranchiseJourney(pool)
	ids := randomid.Generator{}
	svc, e := franchisejourney.NewHandoverPreparationService(repo, ids, policy)
	if e != nil {
		t.Fatal(e)
	}
	prepare := franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", FundingReceiptID: funding.Receipt.ID, ObservationSHA256: funding.SHA256, IdempotencyKey: "local-funding-prepare-0001"}
	oldSvc, e := franchisejourney.NewHandoverPreparationService(repo, ids, old)
	if e != nil {
		t.Fatal(e)
	}
	if _, _, e = oldSvc.Prepare(ctx, tenant, "operator", prepare); e == nil {
		t.Fatal("legacy profile accepted local funding")
	}
	wrong := prepare
	wrong.PaymentAttemptID = "invented-provider-payment"
	if _, _, e = svc.Prepare(ctx, tenant, "operator", wrong); e == nil {
		t.Fatal("mixed identity accepted")
	}
	wrong = prepare
	wrong.ObservationSHA256 = strings.Repeat("0", 64)
	if _, _, e = svc.Prepare(ctx, tenant, "operator", wrong); e == nil {
		t.Fatal("forged receipt hash")
	}
	prepared, replay, e := svc.Prepare(ctx, tenant, "operator", prepare)
	if e != nil || replay || prepared.PaymentAttemptID != "" || prepared.FundingReceiptID != funding.Receipt.ID {
		t.Fatal("full funding preparation", prepared, e)
	}
	recovered, e := svc.Result(ctx, tenant, "store", "order", prepare.IdempotencyKey)
	if e != nil || recovered.Handover.ID != prepared.Handover.ID {
		t.Fatal("prepare recovery", e)
	}
	view, e := svc.OperatorContext(ctx, tenant, "store", "order")
	if e != nil || view.PaymentAttemptID != "" || view.FundingReceiptID != funding.Receipt.ID {
		t.Fatal("operator funding context", view, e)
	}
	if _, e = repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "local-funded-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); e != nil {
		t.Fatal(e)
	}
	presented, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", prepared.Handover.ID, 1, "local-funded-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
	if e != nil {
		t.Fatal(e)
	}
	if _, e = repo.AcceptHandover(ctx, tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "local-funded-checklist", 1, strings.Repeat("e", 64), ids.New()); e != nil {
		t.Fatal(e)
	}
	commit := franchisejourney.CommitCommercialReleaseCommand{OrganizationID: "store", HandoverID: prepared.Handover.ID, ObservationSHA256: funding.SHA256, IdempotencyKey: "local-funding-release-0001"}
	receipt, replay, e := svc.CommitCommercialRelease(ctx, tenant, "operator", commit)
	if e != nil || replay || receipt.PaymentAttemptID != "" || receipt.FundingReceiptID != funding.Receipt.ID {
		t.Fatal("local funding release", receipt, e)
	}
	current, e := svc.ValidateCommercialRelease(ctx, tenant, "store", prepared.Handover.ID)
	if e != nil || !current.Current {
		t.Fatal("current receipt", current, e)
	}
	recoveredRelease, e := svc.CommercialReleaseResult(ctx, tenant, "store", prepared.Handover.ID, commit.IdempotencyKey)
	if e != nil || recoveredRelease.ID != receipt.ID {
		t.Fatal("release recovery", e)
	}
	// Negative source-state fixture, not a claim that a customer cancellation
	// endpoint is permitted after physical delivery. Existing truth invalidates.
	if _, e = pool.Exec(ctx, `update sales.customer_order set state='cancelled',version=version+1 where tenant_id=$1 and order_id='order'`, tenant); e != nil {
		t.Fatal(e)
	}
	current, e = svc.ValidateCommercialRelease(ctx, tenant, "store", prepared.Handover.ID)
	if e != nil || current.Current {
		t.Fatal("cancelled order remained releasable", e)
	}
	var version int64
	if e = pool.QueryRow(ctx, `select version from sales.customer_order where tenant_id=$1 and order_id='order'`, tenant).Scan(&version); e != nil {
		t.Fatal(e)
	}
	v, e := store.Propose(ctx, actor, sv.Request{OperationID: "reverse-funded-redemption", Operation: "reverse", OrganizationID: "store", OrderID: "order", ProgramID: "reference-gift_card", ExpectedOrderVersion: version, OriginalOperationID: "redemption", ProfileSHA256: store.ProfileSHA256()})
	if e != nil {
		t.Fatal("approved cancellation inverse", e)
	}
	reviewer := actor
	reviewer.Subject = "reviewer"
	if _, e = store.Decide(ctx, reviewer, sv.Decision{OrganizationID: "store", ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: true, Reason: "Exact source inverse after fixture cancellation"}); e != nil {
		t.Fatal(e)
	}
	var balance string
	if e = pool.QueryRow(ctx, `select sum(points_delta)::numeric(38,6)::text from stored_value.entry where tenant_id=$1 and account_id=$2`, tenant, account).Scan(&balance); e != nil || balance != "2000.000000" {
		t.Fatal("gift conservation", balance, e)
	}
	history, e := store.FundingResult(ctx, actor, command.RequestKey)
	if e != nil || history.SHA256 != funding.SHA256 {
		t.Fatal("historical receipt changed", e)
	}
	var snapshot []byte
	e = pool.QueryRow(ctx, `select jsonb_build_object('payments',(select count(*) from payment.payment_attempt where tenant_id=$1),'observations',(select count(*) from payment.provider_observation where tenant_id=$1),'funding',(select count(*) from payment.local_funding_receipt where tenant_id=$1),'funding_events',(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='order.local-funding-observed'),'releases',(select count(*) from sales.commercial_release_receipt where tenant_id=$1))`, tenant).Scan(&snapshot)
	var totals map[string]int
	if e != nil || json.Unmarshal(snapshot, &totals) != nil || totals["payments"] != 0 || totals["observations"] != 0 || totals["funding"] != 1 || totals["funding_events"] != 1 || totals["releases"] != 1 {
		t.Fatal("no fake payment", string(snapshot), e)
	}
	t.Log(string(snapshot))
}
