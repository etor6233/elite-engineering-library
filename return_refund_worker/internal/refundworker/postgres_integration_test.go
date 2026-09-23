package refundworker

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type succeedingProvider struct{}

func (succeedingProvider) Create(_ context.Context, value Refund) (ProviderResult, error) {
	return ProviderResult{ProviderPaymentReference: value.ProviderPaymentReference, ProviderRefundReference: "provider-" + value.RequestID, ProviderStatus: "succeeded", Currency: value.Currency, AmountMinorUnits: value.AmountMinorUnits}, nil
}
func (succeedingProvider) Retrieve(_ context.Context, value Refund) (ProviderResult, error) {
	return ProviderResult{ProviderPaymentReference: value.ProviderPaymentReference, ProviderRefundReference: value.ProviderRefundReference, ProviderStatus: "succeeded", Currency: value.Currency, AmountMinorUnits: value.AmountMinorUnits}, nil
}

// This legacy compatibility fixture bypasses relational triggers when seeding
// and removing synthetic data. It does not validate the originating business
// workflow. Require a separate disposable database, never a project/audit DB.
func refundTestConfig(raw string) (*pgxpool.Config, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("explicit disposable database is required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_refund_test_") || len(strings.TrimPrefix(cfg.ConnConfig.Database, "elite_refund_test_")) < 16 {
		return nil, errors.New("requires a dedicated disposable loopback elite_refund_test_<unique> database")
	}
	for _, fallback := range cfg.ConnConfig.Fallbacks {
		if fallback.Host != "127.0.0.1" {
			return nil, errors.New("non-loopback fallback forbidden")
		}
	}
	return cfg, nil
}

func TestRefundDatabaseGuard(t *testing.T) {
	for _, value := range []string{"", "invalid", "postgres://u@remote/elite_refund_test_0123456789abcdef", "postgres://u@localhost/elite_refund_test_0123456789abcdef", "postgres://u@127.0.0.1/production", "postgres://u@127.0.0.1/elite_confirmation_audit", "postgres://u@127.0.0.1/elite_refund_test_short", "host=127.0.0.1,remote dbname=elite_refund_test_0123456789abcdef"} {
		if cfg, err := refundTestConfig(value); err == nil || cfg != nil {
			t.Fatal("unsafe target accepted")
		}
	}
	if cfg, err := refundTestConfig("postgres://u@127.0.0.1/elite_refund_test_0123456789abcdef"); err != nil || cfg == nil {
		t.Fatal("disposable target rejected", err)
	}
}

func TestPostgresRefundLineAllocationPartialThenFullAndAmbiguity(t *testing.T) {
	testRefundLineAllocationPartialThenFullAndAmbiguity(t, "stripe", succeedingProvider{})
}

func testRefundLineAllocationPartialThenFullAndAmbiguity(t *testing.T, sourceProvider string, provider Provider) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := refundTestConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2f201"
	cleanup := func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"payment.return_refund_observation", "payment.return_refund", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "sales.return_authorization", "payment.payment_attempt", "sales.customer_order_line", "sales.customer_order", "platform.outbox_event", "inventory.stock_unit", "org.organization", "platform.tenant"} {
			_, _ = tx.Exec(ctx, `delete from `+table+` where tenant_id=$1`, tenant)
		}
		_ = tx.Commit(ctx)
	}
	cleanup()
	defer cleanup()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `set local session_replication_role=replica`); err != nil {
		t.Fatal(err)
	}
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'return-refund','Return Refund','Return Refund')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store','return-refund-store','Store','store')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'stock-1','store','variant','SERIAL-1','VIN-1','BATTERY-1','sold',1,clock_timestamp()),($1,'stock-2','store','variant','SERIAL-2','VIN-2','BATTERY-2','sold',1,clock_timestamp()),($1,'stock-3','store','variant','SERIAL-3','VIN-3','BATTERY-3','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order-1','store','customer','delivered','ARS',2000,3),($1,'order-2','store','customer','delivered','ARS',1000,3)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id) values($1,'order-1','line-1','variant',1,1000,'stock-1'),($1,'order-1','line-2','variant',1,1000,'stock-2'),($1,'order-2','line-3','variant',1,1000,'stock-3')`,
		`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version) values($1,'payment-1','order-1','stripe','pi_payment_1','payment-key-000001','captured','ARS',2000,3),($1,'payment-2a','order-2','stripe','pi_payment_2a','payment-key-000002','captured','ARS',1000,3),($1,'payment-2b','order-2','stripe','pi_payment_2b','payment-key-000003','captured','ARS',1000,3)`,
		`insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject) values($1,'auth-1','store','exception-1','handover-1','order-1','stock-1','customer','return','order-1','authorized','operator'),($1,'auth-2','store','exception-2','handover-2','order-1','stock-2','customer','return','order-1','authorized','operator'),($1,'auth-3','store','exception-3','handover-3','order-2','stock-3','customer','return','order-2','authorized','operator')`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject) values($1,'receipt-1','auth-1','store','order-1','stock-1','customer','SERIAL-1','opened','received','aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa','operator'),($1,'receipt-2','auth-2','store','order-1','stock-2','customer','SERIAL-2','opened','received','bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb','operator'),($1,'receipt-3','auth-3','store','order-2','stock-3','customer','SERIAL-3','opened','received','cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc','operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject) values($1,'disposition-1','receipt-1','restock','refund','approved','operator'),($1,'disposition-2','receipt-2','restock','refund','approved','operator'),($1,'disposition-3','receipt-3','restock','refund','approved','operator')`,
	}
	for _, query := range fixtures {
		if _, err = tx.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if sourceProvider != "stripe" {
		if _, err = tx.Exec(ctx, `update payment.payment_attempt set provider_code=$2,provider_reference=case payment_attempt_id when 'payment-1' then '7186040733' when 'payment-2a' then '7186040734' else '7186040735' end where tenant_id=$1`, tenant, sourceProvider); err != nil {
			t.Fatal(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key) values($1,'refund-1','disposition-1','refund','payment','requested','return-refund-key-0001'),($1,'refund-2','disposition-2','refund','payment','requested','return-refund-key-0002'),($1,'refund-3','disposition-3','refund','payment','requested','return-refund-key-0003')`, tenant); err != nil {
		t.Fatal(err)
	}

	store := NewPostgresStore(pool)
	refundProvider := sourceProvider
	if sourceProvider == "mercadopago" {
		refundProvider = "mercado_pago"
	}
	processor, err := NewProcessor(store, map[string]Provider{refundProvider: provider}, "refund-worker", time.Minute, time.Millisecond, 5)
	if err != nil {
		t.Fatal(err)
	}
	if err = processor.Step(ctx, "claim-refund-1"); err != nil {
		t.Fatal(err)
	}
	var paymentState, executionState string
	var paymentVersion int64
	var observations, events int
	if err = pool.QueryRow(ctx, `select state,version from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment-1'`, tenant).Scan(&paymentState, &paymentVersion); err != nil {
		t.Fatal(err)
	}
	if paymentState != "captured" || paymentVersion != 3 {
		t.Fatalf("partial refund incorrectly closed payment: state=%s version=%d", paymentState, paymentVersion)
	}
	// The boundary alias must never rewrite historical payment identity or keys.
	var storedProvider, paymentKey string
	if err = pool.QueryRow(ctx, `select provider_code,idempotency_key from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment-1'`, tenant).Scan(&storedProvider, &paymentKey); err != nil || storedProvider != sourceProvider || paymentKey != "payment-key-000001" {
		t.Fatalf("source identity changed: provider=%s key=%s err=%v", storedProvider, paymentKey, err)
	}
	prepared, err := scanRefund(pool.QueryRow(ctx, `select tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,coalesce(provider_refund_reference,''),coalesce(provider_status,'') from payment.return_refund where tenant_id=$1 and request_id='refund-1'`, tenant))
	if err != nil || prepared.Provider != refundProvider || prepared.PaymentAttemptID != "payment-1" || prepared.RequestID != "refund-1" || prepared.IdempotencyKey != "return-refund-key-0001" || prepared.AmountMinorUnits != 1000 {
		t.Fatalf("refund identity or contract changed: %+v err=%v", prepared, err)
	}
	retrieved, err := provider.Retrieve(ctx, prepared)
	if err != nil || validateResult(prepared, retrieved) != nil || retrieved.ProviderRefundReference != prepared.ProviderRefundReference {
		t.Fatalf("persisted refund could not reconcile using its unchanged identity: %+v err=%v", retrieved, err)
	}
	if err = processor.Step(ctx, "claim-refund-2"); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select state,version from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment-1'`, tenant).Scan(&paymentState, &paymentVersion); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='refund-2'`, tenant).Scan(&executionState); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from payment.return_refund_observation where tenant_id=$1 and request_id in ('refund-1','refund-2')`, tenant).Scan(&observations); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type in ('return-refund.succeeded','payment.refunded')`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if paymentState != "refunded" || paymentVersion != 4 || executionState != "succeeded" || observations != 2 || events != 3 {
		t.Fatalf("state=%s version=%d execution=%s observations=%d events=%d", paymentState, paymentVersion, executionState, observations, events)
	}
	if err = processor.Step(ctx, "claim-refund-3"); !errors.Is(err, ErrMappingConflict) {
		t.Fatalf("ambiguous payment mapping accepted: %v", err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='refund-3'`, tenant).Scan(&executionState); err != nil {
		t.Fatal(err)
	}
	if executionState != "blocked" {
		t.Fatalf("ambiguous mapping status=%s", executionState)
	}
	if _, err = pool.Exec(ctx, `update payment.return_refund_observation set provider_status='invented' where tenant_id=$1 and request_id='refund-1'`, tenant); err == nil {
		t.Fatal("immutable refund observation was mutated")
	}
	if work, claimErr := store.Claim(ctx, "refund-worker", "claim-extra", time.Minute); claimErr != nil || work != nil {
		t.Fatalf("terminal refund work reclaimed: %+v %v", work, claimErr)
	}
}
