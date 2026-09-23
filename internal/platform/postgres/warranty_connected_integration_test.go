package postgres_test

// AUTHORED fixture composes existing quote/order, official SDK, callback,
// durable reconciliation, checklist and handover with sold warranty terms.
import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgxpool"
)

func fixtureWarranty(t *testing.T, pool *pgxpool.Pool, tenant string, days int, leaseSeconds ...int) (*db.Warranty, identity.Principal, wc.Profile) {
	t.Helper()
	document := wc.ProfileDocument{Schema: "elite-warranty-profile/v1", Scope: "MATERIALIZED_PROFILE", Algorithm: "bc-inclusive-fixed-terms", AlgorithmRevision: 1, TenantID: tenant, OrganizationID: "store", FactoryOrganizationID: "factory", PolicyID: "synthetic", TermsVersion: "fixture-v1", TermsText: "Synthetic fixture terms; no jurisdiction or legal compliance claim.", BusinessTimeZone: "America/Argentina/Buenos_Aires", PartsDurationDays: days, LaborDurationDays: days, FaultExclusions: []string{"fixture-exclusion"}, WorkReservationSeconds: 3600, Settlement: "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT", AuthorityReference: "synthetic-authority", DecisionReference: "synthetic-acceptance"}
	if len(leaseSeconds) > 1 {
		t.Fatal("one fixture lease maximum")
	}
	if len(leaseSeconds) == 1 {
		document.WorkReservationSeconds = leaseSeconds[0]
	}
	raw, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	p, err := wc.LoadProfile(raw, wc.SHA(raw))
	if err != nil {
		t.Fatal(err)
	}
	store, err := db.NewWarranty(pool, p)
	if err != nil {
		t.Fatal(err)
	}
	actor := identity.Principal{TenantID: tenant, Subject: "operator", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"warranty:offer": {}, "warranty:read": {}, "warranty:activate": {}}}
	return store, actor, p
}
func TestWarrantySoldTermsThroughConnectedHandover(t *testing.T) {
	ctx := context.Background()
	ids := randomid.Generator{}
	var store *db.Warranty
	var operator, customer identity.Principal
	var sold wc.Profile
	var activated wc.Activation
	r := newConnectedRunWithAllocation(t, false, nil, func(t *testing.T, pool *pgxpool.Pool, tenant string) int64 {
		store, operator, sold = fixtureWarranty(t, pool, tenant, 365)
		customer = identity.Principal{TenantID: tenant, Subject: "customer", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"warranty:self": {}}}
		request := wc.OfferRequest{QuoteID: "quote", QuoteVersion: 1, ProfileSHA256: sold.Hash()}
		offer, err := store.BindOffer(ctx, operator, request)
		if err != nil || offer.QuoteVersion != 2 {
			t.Fatal("bind", offer, err)
		}
		replay, err := store.BindOffer(ctx, operator, request)
		if err != nil || !replay.Replay || replay.OfferedAt != offer.OfferedAt {
			t.Fatal("offer recovery", replay, err)
		}
		repo := db.NewFranchiseJourney(pool)
		if _, err = repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 2, strings.Repeat("a", 64), "order", "line", ids.New(), ids.New()); err == nil {
			t.Fatal("accepted without warranty consent")
		}
		var orders int
		if err = pool.QueryRow(ctx, `select count(*) from sales.customer_order where tenant_id=$1`, tenant).Scan(&orders); err != nil || orders != 0 {
			t.Fatal("failed acceptance leaked order", orders, err)
		}
		wrong := customer
		wrong.Subject = "different-customer"
		ack := wc.AcknowledgeRequest{QuoteID: "quote", QuoteVersion: 2, ProfileSHA256: sold.Hash(), EvidenceSHA256: strings.Repeat("d", 64)}
		if _, err = store.Acknowledge(ctx, wrong, ack); err == nil {
			t.Fatal("other customer acknowledged")
		}
		if _, err = store.Offer(ctx, wrong, "quote"); err == nil {
			t.Fatal("other customer read terms")
		}
		wrong = operator
		wrong.TenantID = ids.New()
		if _, err = store.BindOffer(ctx, wrong, request); err == nil {
			t.Fatal("other tenant bound terms")
		}
		wrong = operator
		wrong.Organizations = map[string]struct{}{"other": {}}
		if _, err = store.Offer(ctx, wrong, "quote"); err == nil {
			t.Fatal("other organization read terms")
		}
		bad := ack
		bad.ProfileSHA256 = strings.Repeat("f", 64)
		if _, err = store.Acknowledge(ctx, customer, bad); err == nil {
			t.Fatal("different terms acknowledged")
		}
		bad = ack
		bad.QuoteVersion = 3
		if _, err = store.Acknowledge(ctx, customer, bad); err == nil {
			t.Fatal("stale version acknowledged")
		}
		if _, err = pool.Exec(ctx, `update sales.quotation set total_minor_units=1 where tenant_id=$1 and quotation_id='quote'`, tenant); err == nil {
			t.Fatal("offered price mutated")
		}
		// Two uncertain duplicate submissions serialize to one immutable consent.
		var wg sync.WaitGroup
		results := make(chan error, 2)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); _, err := store.Acknowledge(ctx, customer, ack); results <- err }()
		}
		wg.Wait()
		close(results)
		for err := range results {
			if err != nil {
				t.Fatal("concurrent consent", err)
			}
		}
		offer, err = store.Offer(ctx, customer, "quote")
		if err != nil || !offer.Acknowledged || offer.EvidenceSHA256 != ack.EvidenceSHA256 {
			t.Fatal("consent recovery", offer, err)
		}
		bad = ack
		bad.EvidenceSHA256 = strings.Repeat("e", 64)
		if _, err = store.Acknowledge(ctx, customer, bad); err == nil {
			t.Fatal("changed replay accepted")
		}
		return offer.QuoteVersion
	})
	if code := r.callback(t, "evt_warranty", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	var hash string
	if err := r.pool.QueryRow(ctx, `select evidence_sha256_hex from payment.provider_observation where tenant_id=$1`, r.tenant).Scan(&hash); err != nil {
		t.Fatal(err)
	}
	assertConnectedCommercialReleaseWithHandoverHook(t, r, hash, func(t *testing.T, r *connectedRun, handover string) {
		// New configuration must not replace the profile the customer actually bought.
		changed, _, current := fixtureWarranty(t, r.pool, r.tenant, 7)
		if current.Hash() == sold.Hash() {
			t.Fatal("fixture requires distinct policy revisions")
		}
		// Fail after warranty/activation writes but before transaction commit.
		_, err := r.pool.Exec(ctx, `create function service_ops.fixture_warranty_outbox_failure() returns trigger language plpgsql as $$ begin if new.event_type='warranty.activated' then raise exception 'fixture commit failure'; end if;return new;end $$;
 create trigger fixture_warranty_outbox_failure before insert on platform.outbox_event for each row execute function service_ops.fixture_warranty_outbox_failure()`)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = changed.Activate(ctx, operator, handover); err == nil {
			t.Fatal("fault injection did not fail")
		}
		var count int
		if err = r.pool.QueryRow(ctx, `select count(*) from service_ops.warranty where tenant_id=$1`, r.tenant).Scan(&count); err != nil || count != 0 {
			t.Fatal("activation rollback leaked warranty", count, err)
		}
		if _, err = r.pool.Exec(ctx, `drop trigger fixture_warranty_outbox_failure on platform.outbox_event;drop function service_ops.fixture_warranty_outbox_failure()`); err != nil {
			t.Fatal(err)
		}
		activated, err = changed.Activate(ctx, operator, handover)
		if err != nil {
			t.Fatal("activate sold terms", err)
		}
		dates, err := sold.DatesAt(activated.AcceptedAt)
		if err != nil || activated.ProfileSHA256 != sold.Hash() || activated.Dates != dates {
			t.Fatal("current configuration replaced sold terms", activated, dates, err)
		}
		again, err := changed.Activate(ctx, operator, handover)
		if err != nil || !again.Replay || again.WarrantyID != activated.WarrantyID || !again.ActivatedAt.Equal(activated.ActivatedAt) {
			t.Fatal("activation recovery", again, err)
		}
		read, err := changed.Activation(ctx, customer, handover)
		if err != nil || read.WarrantyID != activated.WarrantyID {
			t.Fatal("customer receipt", read, err)
		}
		wrong := operator
		wrong.Subject = "another-operator"
		if _, err = changed.Activate(ctx, wrong, handover); err == nil {
			t.Fatal("different actor replay")
		}
		if _, err = r.pool.Exec(ctx, `update service_ops.warranty_activation set parts_end=parts_end+1 where tenant_id=$1`, r.tenant); err == nil {
			t.Fatal("activation rewritten")
		}
		wire := newWarrantyHTTPFixture(t, changed, true)
		assertWarrantyClaimLifecycle(t, r, wire, operator, customer, activated)
		wire.mu.Lock()
		drops, recovered := len(wire.drops), wire.recoveries
		workPosts := wire.posts["/v1/franchise/warranty/claims/claim-main/work"]
		reconciliationPosts := wire.posts["/v1/factory/warranty/claims/claim-main/reconciliation"]
		wire.mu.Unlock()
		if drops != 2 || recovered < 2 || workPosts != 3 || reconciliationPosts != 3 {
			t.Fatal("HTTP loss recovery accounting", drops, recovered, workPosts, reconciliationPosts)
		}
		t.Log("WARRANTY_HTTP_PG_PASS claims=4 dropped_commit_responses=2 recovery_gets_verified=true hidden_post_retry=false denied_factory_post=1 verifier_fixture_only=true")
	})
	var warranties, activations, events int
	var state string
	err := r.pool.QueryRow(ctx, `select (select count(*) from service_ops.warranty where tenant_id=$1),(select count(*) from service_ops.warranty_activation where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='warranty.activated'),(select status from service_ops.warranty where tenant_id=$1)`, r.tenant).Scan(&warranties, &activations, &events, &state)
	if err != nil || warranties != 1 || activations != 1 || events != 1 || state != "active" {
		t.Fatal("durable warranty or invented refund voiding", warranties, activations, events, state, err)
	}
	t.Log("WARRANTY_CONNECTED_PASS quote_terms_customer_consent_order_official_sdk_callback_reconcile_handover_commercial_receipt_activation=true sold_profile_preserved=true rollback=true scope=true exact_replay=true claim_workflow_connected=true live_proven=false")
}
