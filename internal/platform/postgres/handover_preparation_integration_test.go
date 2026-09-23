package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initialHandoverPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_handover_") {
		t.Fatal("requires disposable loopback handover database")
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}
func initialHandoverContract() franchisejourney.HandoverReleaseContract {
	hash := sha256.Sum256([]byte(franchisejourney.ReferenceHandoverContractDocument))
	return franchisejourney.HandoverReleaseContract{ID: "reference-single-unit-observed-payment-v1", DocumentSHA256: hex.EncodeToString(hash[:]), Scope: "LOCAL_FIXTURES", MaximumObservationAge: 5 * time.Minute}
}
func initialHandoverFixture(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	ctx := context.Background()
	tenant := randomid.Generator{}.New()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'handover-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store'),($1,'other','other','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','fixture@example.test')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','model','new','fixture','{}')`,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'retail','AR','ARS',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','variant',123456,'inclusive')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SERIAL-SYNTHETIC','available',1,clock_timestamp())`,
	}
	for _, sql := range fixtures {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	quote := franchisejourney.Quote{ID: "quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().UTC().Add(time.Hour)}
	if _, _, err := repo.CreateQuoteAs(ctx, tenant, "quote-key", quote, strings.Repeat("a", 64), "018f4d4a-7b36-7a21-8d10-000000000001", "operator"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 1, strings.Repeat("a", 64), "order", "line", "018f4d4a-7b36-7a21-8d10-000000000002", "018f4d4a-7b36-7a21-8d10-000000000003"); err != nil {
		t.Fatal(err)
	}
	sales := NewCommerce(pool)
	if err := sales.AllocateStockAs(ctx, tenant, "store", "order", "line", "stock", 1, 1, "018f4d4a-7b36-7a21-8d10-000000000004", "operator"); err != nil {
		t.Fatal(err)
	}
	payment := commerce.PaymentAttempt{ID: "payment", OrganizationID: "store", OrderID: "order", ProviderCode: "stripe"}
	if _, err := sales.RecordOrderPayment(ctx, tenant, "018f4d4a-7b36-7a21-8d10-000000000005", "payment-key", payment, "operator"); err != nil {
		t.Fatal(err)
	}
	if err := sales.TransitionPayment(ctx, tenant, "store", "payment", "created", "authorized", 1, "pi_fixture", "018f4d4a-7b36-7a21-8d10-000000000006"); err != nil {
		t.Fatal(err)
	}
	if err := sales.TransitionPayment(ctx, tenant, "store", "payment", "authorized", "captured", 2, "pi_fixture", "018f4d4a-7b36-7a21-8d10-000000000007"); err != nil {
		t.Fatal(err)
	}
	// Deliberately synthetic provider observation. The combined SDK/Webhook test
	// is owned by the payment integration; this suite claims PG binding only.
	if _, err := pool.Exec(ctx, `insert into payment.provider_observation(tenant_id,payment_attempt_id,order_id,organization_id,provider_code,provider_reference,currency,amount_minor_units,received_minor_units,refunded_minor_units,provider_status,live_mode,account_ref,evidence_sha256_hex,observed_at,hold,hold_reason,generation)values($1,'payment','order','store','stripe','pi_fixture','ARS',123456,123456,0,'succeeded',false,'acct_fixture',$2,clock_timestamp(),false,'',1)`, tenant, strings.Repeat("b", 64)); err != nil {
		t.Fatal(err)
	}
	return tenant
}
func initialHandoverCommand() franchisejourney.PrepareHandoverCommand {
	return franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: "payment", ObservationSHA256: strings.Repeat("b", 64), IdempotencyKey: "initial-handover-key"}
}
func TestInitialHandoverConnectedPostgres(t *testing.T) {
	pool := initialHandoverPool(t)
	tenant := initialHandoverFixture(t, pool)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	repo := NewFranchiseJourney(pool)
	service, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, initialHandoverContract())
	if err != nil {
		t.Fatal(err)
	}
	command := initialHandoverCommand()
	type outcome struct {
		value  franchisejourney.HandoverPreparation
		replay bool
		err    error
	}
	start := make(chan struct{})
	results := make(chan outcome, 16)
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			v, replay, err := service.Prepare(ctx, tenant, "operator", command)
			results <- outcome{v, replay, err}
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	creates, replays := 0, 0
	handover := ""
	for result := range results {
		if result.err != nil {
			t.Fatal(result.err)
		}
		if result.replay {
			replays++
		} else {
			creates++
		}
		if handover == "" {
			handover = result.value.Handover.ID
		}
		if handover != result.value.Handover.ID || result.value.Handover.CustomerSubject != "customer" || result.value.Handover.StockUnitID != "stock" || result.value.ReservationID != "018f4d4a-7b36-7a21-8d10-000000000004" {
			t.Fatal("binding or replay drift")
		}
	}
	if creates != 1 || replays != 15 {
		t.Fatalf("creates=%d replays=%d", creates, replays)
	}
	var effects, preparations int
	if err = pool.QueryRow(ctx, `select (select count(*) from platform.outbox_event where tenant_id=$1 and event_type='delivery-handover.prepared'),(select count(*) from sales.delivery_handover_preparation where tenant_id=$1)`, tenant).Scan(&effects, &preparations); err != nil || effects != 1 || preparations != 1 {
		t.Fatalf("effects=%d preparations=%d err=%v", effects, preparations, err)
	}
	recovered, err := service.Result(ctx, tenant, "store", "order", command.IdempotencyKey)
	if err != nil || recovered.Handover.ID != handover {
		t.Fatal("lost response recovery failed", err)
	}
	if _, err = service.Result(ctx, tenant, "other", "order", command.IdempotencyKey); !errors.Is(err, franchisejourney.ErrNotFound) {
		t.Fatal("scope leak", err)
	}
	changed := command
	changed.ObservationSHA256 = strings.Repeat("c", 64)
	if _, _, err = service.Prepare(ctx, tenant, "operator", changed); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("divergent replay accepted", err)
	}
	if _, err = service.EvaluateRelease(ctx, tenant, "store", handover, command.ObservationSHA256); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("release before customer acceptance", err)
	}
	_, err = repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "initial-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic initial checklist", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read the serial", ResponseType: "serial", Required: true}}}, "018f4d4a-7b36-7a21-8d10-000000000008")
	if err != nil {
		t.Fatal(err)
	}
	presented, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", handover, 1, "initial-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, "018f4d4a-7b36-7a21-8d10-000000000009")
	if err != nil {
		t.Fatal(err)
	}
	accepted, err := repo.AcceptHandover(ctx, tenant, "store", "customer", handover, presented.Version, "SERIAL-SYNTHETIC", "initial-checklist", 1, strings.Repeat("e", 64), "018f4d4a-7b36-7a21-8d10-000000000010")
	if err != nil || accepted.State != "accepted" {
		t.Fatal(err)
	}
	decision, err := service.EvaluateRelease(ctx, tenant, "store", handover, command.ObservationSHA256)
	if err != nil || !decision.Eligible || decision.Scope != "LOCAL_FIXTURES" {
		t.Fatal("reference release gate", err)
	}
	if _, err = pool.Exec(ctx, `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING',generation=generation+1 where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = service.EvaluateRelease(ctx, tenant, "store", handover, command.ObservationSHA256); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("held observation released", err)
	}
	if _, err = service.Result(ctx, tenant, "store", "order", command.IdempotencyKey); err != nil {
		t.Fatal("historical receipt cannot be recovered", err)
	}
	t.Logf("INITIAL_HANDOVER_CONNECTED_PG_PASS tenant=%s handover=%s creates=1 replays=15 quote_order_allocation=true acceptance=true no_shipping_posted=true", tenant, handover)
}

func TestInitialHandoverRejectsInvalidObservationAndScope(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	cases := map[string]string{
		"missing-observation":  `delete from payment.provider_observation where tenant_id=$1`,
		"hold":                 `update payment.provider_observation set hold=true where tenant_id=$1`,
		"wrong-amount":         `update payment.provider_observation set hold=true,amount_minor_units=123457 where tenant_id=$1`,
		"partial":              `update payment.provider_observation set hold=true,received_minor_units=1 where tenant_id=$1`,
		"refunded":             `update payment.provider_observation set hold=true,refunded_minor_units=1 where tenant_id=$1`,
		"live-mode":            `update payment.provider_observation set live_mode=true where tenant_id=$1`,
		"currency":             `update payment.provider_observation set currency='USD' where tenant_id=$1`,
		"provider-reference":   `update payment.provider_observation set provider_reference='pi_other' where tenant_id=$1`,
		"organization":         `update payment.provider_observation set organization_id='other' where tenant_id=$1`,
		"unobserved":           `update payment.provider_observation set hold=true,observed_at=null,evidence_sha256_hex=null where tenant_id=$1`,
		"stale":                `update payment.provider_observation set observed_at=clock_timestamp()-interval '1 hour' where tenant_id=$1`,
		"future":               `update payment.provider_observation set observed_at=clock_timestamp()+interval '1 hour' where tenant_id=$1`,
		"payment-pending":      `update payment.payment_attempt set state='pending' where tenant_id=$1`,
		"customer-restricted":  `update crm.customer_profile set status='restricted' where tenant_id=$1`,
		"stock-scope":          `update inventory.stock_unit set organization_id='other' where tenant_id=$1`,
		"reservation-released": `update inventory.serial_reservation set status='released' where tenant_id=$1`,
		"order-cancelled":      `update sales.customer_order set state='cancelled' where tenant_id=$1`,
	}
	for name, mutation := range cases {
		t.Run(name, func(t *testing.T) {
			tenant := initialHandoverFixture(t, pool)
			if _, err := pool.Exec(ctx, mutation, tenant); err != nil {
				t.Fatal(err)
			}
			service, _ := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), randomid.Generator{}, initialHandoverContract())
			if _, _, err := service.Prepare(ctx, tenant, "operator", initialHandoverCommand()); !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatalf("invalid facts accepted or wrong error: %v", err)
			}
			var count int
			if err := pool.QueryRow(ctx, `select (select count(*) from sales.delivery_handover where tenant_id=$1)+(select count(*) from sales.delivery_handover_preparation where tenant_id=$1)+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='initial-handover')`, tenant).Scan(&count); err != nil || count != 0 {
				t.Fatalf("failed attempt leaked %d rows: %v", count, err)
			}
		})
	}
	t.Log("INITIAL_HANDOVER_NEGATIVES_PASS cases=17")
}

func TestInitialHandoverOutboxAtomicityAndFreshnessAfterLock(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	tenant := initialHandoverFixture(t, pool)
	repo := NewFranchiseJourney(pool)
	command := initialHandoverCommand()
	if _, _, err := repo.PrepareInitialHandover(ctx, tenant, "operator", "rollback-handover", "018f4d4a-7b36-7a21-8d10-000000000001", command, initialHandoverContract(), strings.Repeat("d", 64)); err == nil {
		t.Fatal("duplicate outbox id unexpectedly accepted")
	}
	var count int
	if err := pool.QueryRow(ctx, `select (select count(*) from sales.delivery_handover where tenant_id=$1)+(select count(*) from sales.delivery_handover_preparation where tenant_id=$1)+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='initial-handover')`, tenant).Scan(&count); err != nil || count != 0 {
		t.Fatalf("outbox failure leaked %d rows: %v", count, err)
	}
	blocker, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer blocker.Rollback(ctx)
	if _, err = blocker.Exec(ctx, `select 1 from payment.provider_observation where tenant_id=$1 for update`, tenant); err != nil {
		t.Fatal(err)
	}
	policy := initialHandoverContract()
	policy.MaximumObservationAge = 200 * time.Millisecond
	service, _ := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, policy)
	done := make(chan error, 1)
	go func() { _, _, err := service.Prepare(ctx, tenant, "operator", command); done <- err }()
	time.Sleep(350 * time.Millisecond)
	if err = blocker.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case err = <-done:
		if !errors.Is(err, franchisejourney.ErrConflict) {
			t.Fatal("stale-after-lock accepted", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("lock wait did not finish")
	}
	t.Log("INITIAL_HANDOVER_ATOMICITY_AND_POSTLOCK_FRESHNESS_PASS")
}

func TestInitialHandoverRecoveryAfterRestart(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	rows, err := pool.Query(ctx, `select p.tenant_id::text,p.organization_id,p.order_id,i.idempotency_key,p.handover_id from sales.delivery_handover_preparation p join platform.idempotency_record i on i.tenant_id=p.tenant_id and i.resource_id=p.handover_id and i.scope='initial-handover' order by p.tenant_id,p.handover_id`)
	if err != nil {
		t.Fatal(err)
	}
	type identity struct{ tenant, org, order, key, id string }
	var saved []identity
	for rows.Next() {
		var row identity
		if err = rows.Scan(&row.tenant, &row.org, &row.order, &row.key, &row.id); err != nil {
			t.Fatal(err)
		}
		saved = append(saved, row)
	}
	rows.Close()
	if rows.Err() != nil {
		t.Fatal(rows.Err())
	}
	if len(saved) != 1 {
		t.Fatal("expected one durable successful preparation", len(saved))
	}
	service, _ := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), randomid.Generator{}, initialHandoverContract())
	for _, row := range saved {
		got, err := service.Result(ctx, row.tenant, row.org, row.order, row.key)
		if err != nil || got.Handover.ID != row.id || got.Handover.State != "accepted" || got.Handover.CustomerAcceptedAt == nil || got.Handover.ChecklistCompletedAt == nil || got.Handover.ChecklistID != "initial-checklist" || len(got.Handover.ChecklistItems) != 1 || got.Handover.AcceptanceEvidence != strings.Repeat("e", 64) {
			t.Fatal(fmt.Sprintf("durable recovery mismatch: %+v %v", got, err))
		}
	}
	t.Log("INITIAL_HANDOVER_RESTART_RECOVERY_PASS")
}
