package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func commercialProfile(t *testing.T, tenant string, revision int) franchisejourney.HandoverReleaseContract {
	t.Helper()
	mode := false
	d := franchisejourney.HandoverProfileDocument{Schema: franchisejourney.HandoverProfileSchema, ProfileID: "commercial-fixture", Revision: revision, Algorithm: franchisejourney.HandoverSupportedAlgorithm, AlgorithmRevision: 2, Scope: "MATERIALIZED_PROFILE", TenantID: tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout", ExpectedLiveMode: &mode, MaximumObservationAgeSeconds: 300, Options: franchisejourney.SupportedCommercialReleaseOptions(), AuthorityReference: "docs/commercial-release.md", DecisionReference: "FIXTURE_ONLY"}
	raw, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	p, err := franchisejourney.LoadHandoverProfile(raw, franchisejourney.HandoverProfileActivation{Enabled: true, ProfileID: d.ProfileID, ProfileRevision: revision, DocumentSHA256: hex.EncodeToString(sum[:]), TenantID: tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout"})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func commercialFixture(t *testing.T, pool *pgxpool.Pool, accept bool) (string, *franchisejourney.HandoverPreparationService, franchisejourney.CommitCommercialReleaseCommand) {
	t.Helper()
	ctx := context.Background()
	tenant := initialHandoverFixture(t, pool)
	for _, sql := range []string{
		`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'checkout','stripe','FIXTURE_ONLY_NOT_SECRET')`,
		`insert into payment.provider_checkout(tenant_id,payment_attempt_id,provider_code,connection_id,account_ref,request_sha256_hex,live_mode,expires_at) values($1,'payment','stripe','checkout','acct_fixture',repeat('a',64),false,clock_timestamp()+interval '1 hour')`,
	} {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	svc, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, commercialProfile(t, tenant, 1))
	if err != nil {
		t.Fatal(err)
	}
	prepared, _, err := svc.Prepare(ctx, tenant, "operator", initialHandoverCommand())
	if err != nil {
		t.Fatal(err)
	}
	if accept {
		_, err = repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "release-checklist", OrganizationID: "store", Version: 1, Title: "Fixture", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, randomid.Generator{}.New())
		if err != nil {
			t.Fatal(err)
		}
		h, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", prepared.Handover.ID, 1, "release-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, randomid.Generator{}.New())
		if err != nil {
			t.Fatal(err)
		}
		_, err = repo.AcceptHandover(ctx, tenant, "store", "customer", h.ID, h.Version, "SERIAL-SYNTHETIC", "release-checklist", 1, strings.Repeat("e", 64), randomid.Generator{}.New())
		if err != nil {
			t.Fatal(err)
		}
	}
	return tenant, svc, franchisejourney.CommitCommercialReleaseCommand{OrganizationID: "store", HandoverID: prepared.Handover.ID, ObservationSHA256: strings.Repeat("b", 64), IdempotencyKey: "commercial-release-key"}
}

func assertNoCommercialRelease(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	var rows, events, keys int
	err := pool.QueryRow(context.Background(), `select (select count(*) from sales.commercial_release_receipt where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='commercial-release.recorded'),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='commercial-release')`, tenant).Scan(&rows, &events, &keys)
	if err != nil || rows != 0 || events != 0 || keys != 0 {
		t.Fatalf("partial commercial commit: rows%d events%d keys%d err%v", rows, events, keys, err)
	}
}

func TestCommercialReleaseConnectedPostgres(t *testing.T) {
	pool := initialHandoverPool(t)
	tenant, svc, c := commercialFixture(t, pool, true)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	type outcome struct {
		v      franchisejourney.CommercialReleaseReceipt
		replay bool
		err    error
	}
	out := make(chan outcome, 12)
	start := make(chan struct{})
	for range 12 {
		go func() {
			<-start
			v, r, e := svc.CommitCommercialRelease(ctx, tenant, "operator", c)
			out <- outcome{v, r, e}
		}()
	}
	close(start)
	var receipt franchisejourney.CommercialReleaseReceipt
	creates, replays := 0, 0
	for range 12 {
		o := <-out
		if o.err != nil {
			t.Fatal(o.err)
		}
		if o.replay {
			replays++
		} else {
			creates++
			receipt = o.v
		}
		if o.v.ObservationGeneration != 1 || !o.v.ValidUntil.After(o.v.RecordedAt) || o.v.AcceptanceSHA256 != strings.Repeat("e", 64) {
			t.Fatal("receipt binding")
		}
	}
	if creates != 1 || replays != 11 {
		t.Fatalf("creates%d replays%d", creates, replays)
	}
	current, err := svc.ValidateCommercialRelease(ctx, tenant, "store", c.HandoverID)
	if err != nil || !current.Current || current.Receipt.ID != receipt.ID {
		t.Fatal("current receipt", current, err)
	}
	other := c
	other.IdempotencyKey = "different-release-key"
	if _, _, err = svc.CommitCommercialRelease(ctx, tenant, "operator", other); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("second release accepted", err)
	}
	if _, _, err = svc.CommitCommercialRelease(ctx, tenant, "other-actor", c); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("divergent replay accepted", err)
	}
	if _, err = svc.CommercialReleaseResult(ctx, tenant, "other", c.HandoverID, c.IdempotencyKey); !errors.Is(err, franchisejourney.ErrReleaseConditioned) {
		t.Fatal("scope leak", err)
	}
	var rows, events, keys int
	var stock, payment string
	var amount int64
	err = pool.QueryRow(ctx, `select (select count(*) from sales.commercial_release_receipt where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='commercial-release.recorded'),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='commercial-release'),(select state from inventory.stock_unit where tenant_id=$1 and stock_unit_id='stock'),(select state from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment'),(select amount_minor_units from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment')`, tenant).Scan(&rows, &events, &keys, &stock, &payment, &amount)
	if err != nil || rows != 1 || events != 1 || keys != 1 || stock != "reserved" || payment != "captured" || amount != 123456 {
		t.Fatal("unexpected effects", rows, events, keys, stock, payment, amount, err)
	}
	for _, sql := range []string{`update sales.commercial_release_receipt set released_by_subject='forged' where tenant_id=$1`, `delete from sales.commercial_release_receipt where tenant_id=$1`} {
		if _, err = pool.Exec(ctx, sql, tenant); err == nil {
			t.Fatal("immutable history changed")
		}
	}
	if _, err = pool.Exec(ctx, `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING',generation=generation+1 where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	current, err = svc.ValidateCommercialRelease(ctx, tenant, "store", c.HandoverID)
	if err != nil || current.Current || current.Receipt.ID != receipt.ID {
		t.Fatal("post-callback receipt still current", err)
	}
	v, replay, err := svc.CommitCommercialRelease(ctx, tenant, "operator", c)
	if err != nil || !replay || v.ID != receipt.ID || v.ObservationGeneration != 1 {
		t.Fatal("historical recovery lost", err)
	}
	// Even an identical hash from a newer observation generation cannot reactivate
	// an old release. The current generation is part of its evidence boundary.
	if _, err = pool.Exec(ctx, `update payment.provider_observation set hold=false,hold_reason='',observed_at=clock_timestamp() where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	current, err = svc.ValidateCommercialRelease(ctx, tenant, "store", c.HandoverID)
	if err != nil || current.Current {
		t.Fatal("old generation reactivated", err)
	}
	t.Log("COMMERCIAL_RELEASE_PG_PASS creates=1 replays=11 atomic_receipt_outbox=true historical_recovery=true post_callback_current=false no_money_or_stock_effect=true")
}

func TestCommercialReleaseRejectsUnacceptedStaleHeldAndPolicyDrift(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	for _, name := range []string{"unaccepted", "stale", "held", "wrong-hash", "policy-drift", "account-drift"} {
		t.Run(name, func(t *testing.T) {
			tenant, svc, c := commercialFixture(t, pool, name != "unaccepted")
			sql := ""
			switch name {
			case "stale":
				sql = `update payment.provider_observation set observed_at=clock_timestamp()-interval '10 minutes' where tenant_id=$1`
			case "held":
				sql = `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING' where tenant_id=$1`
			case "account-drift":
				sql = `update payment.provider_observation set account_ref='acct_other' where tenant_id=$1`
			case "wrong-hash":
				c.ObservationSHA256 = strings.Repeat("a", 64)
			case "policy-drift":
				var err error
				svc, err = franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), randomid.Generator{}, commercialProfile(t, tenant, 2))
				if err != nil {
					t.Fatal(err)
				}
			}
			if sql != "" {
				if _, err := pool.Exec(ctx, sql, tenant); err != nil {
					t.Fatal(err)
				}
			}
			if _, _, err := svc.CommitCommercialRelease(ctx, tenant, "operator", c); !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatal("ineligible release accepted", err)
			}
			assertNoCommercialRelease(t, pool, tenant)
		})
	}
}

func TestCommercialReleaseWaitsForConcurrentObservation(t *testing.T) {
	pool := initialHandoverPool(t)
	for _, mode := range []string{"callback-hold", "became-stale"} {
		t.Run(mode, func(t *testing.T) {
			tenant, svc, c := commercialFixture(t, pool, true)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(ctx)
			sql := `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING',generation=generation+1 where tenant_id=$1`
			if mode == "became-stale" {
				sql = `update payment.provider_observation set observed_at=clock_timestamp()-interval '10 minutes' where tenant_id=$1`
			}
			if _, err = tx.Exec(ctx, sql, tenant); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { _, _, e := svc.CommitCommercialRelease(ctx, tenant, "operator", c); done <- e }()
			deadline := time.Now().Add(5 * time.Second)
			for {
				var waiting bool
				err = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like '%from payment.provider_observation where tenant_id=%')`).Scan(&waiting)
				if err != nil {
					t.Fatal(err)
				}
				if waiting {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("release never reached observation lock")
				}
				time.Sleep(10 * time.Millisecond)
			}
			if err = tx.Commit(ctx); err != nil {
				t.Fatal(err)
			}
			if err = <-done; !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatal("release ignored concurrent observation", err)
			}
			assertNoCommercialRelease(t, pool, tenant)
		})
	}
}

type commercialCollisionIDs struct{ calls int }

func (i *commercialCollisionIDs) New() string {
	i.calls++
	if i.calls%2 == 0 {
		return "018f4d4a-7b36-7a21-8d10-000000000001"
	}
	return randomid.Generator{}.New()
}
func TestCommercialReleaseRollbackAndLostResponseRecovery(t *testing.T) {
	pool := initialHandoverPool(t)
	tenant, svc, c := commercialFixture(t, pool, true)
	ctx := context.Background()
	collision, err := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), &commercialCollisionIDs{}, commercialProfile(t, tenant, 1))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = collision.CommitCommercialRelease(ctx, tenant, "operator", c); err == nil {
		t.Fatal("expected existing outbox UUID collision")
	}
	assertNoCommercialRelease(t, pool, tenant)
	v, _, err := svc.CommitCommercialRelease(ctx, tenant, "operator", c)
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := svc.CommercialReleaseResult(ctx, tenant, "store", c.HandoverID, c.IdempotencyKey)
	if err != nil || recovered.ID != v.ID {
		t.Fatal("lost response recovery", err)
	}
}

type commercialWaitingIDs struct {
	event string
	calls int
}

func (i *commercialWaitingIDs) New() string {
	i.calls++
	if i.calls%2 == 0 {
		return i.event
	}
	return randomid.Generator{}.New()
}

func TestCommercialReleaseRejectsExpiryDuringOutboxWait(t *testing.T) {
	pool := initialHandoverPool(t)
	tenant, _, c := commercialFixture(t, pool, true)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	event := randomid.Generator{}.New()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if err = writeJourneyOutbox(ctx, tx, tenant, event, "fixture", "outbox-lock", 1, "fixture.pending", map[string]any{"fixture": true}); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `update payment.provider_observation set observed_at=clock_timestamp()-interval '298 seconds' where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	svc, err := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), &commercialWaitingIDs{event: event}, commercialProfile(t, tenant, 1))
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { _, _, e := svc.CommitCommercialRelease(ctx, tenant, "operator", c); done <- e }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		var waiting bool
		err = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like '%insert into platform.outbox_event%')`).Scan(&waiting)
		if err != nil {
			t.Fatal(err)
		}
		if waiting {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("release never reached outbox lock")
		}
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(2200 * time.Millisecond)
	if err = tx.Rollback(ctx); err != nil {
		t.Fatal(err)
	}
	if err = <-done; !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("expired checkpoint committed after outbox wait", err)
	}
	assertNoCommercialRelease(t, pool, tenant)
}
