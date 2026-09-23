package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/returneffects"
	"elite.local/enterprise/internal/returnexchange"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReturnExchangeAtomicReservationAndNoStockRetry(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2e143"
	cleanup := func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"sales.return_exchange", "sales.return_effect_resume", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "sales.return_authorization", "sales.delivery_exception", "sales.delivery_handover", "sales.customer_order_line", "sales.customer_order", "platform.outbox_event", "inventory.stock_unit", "crm.customer_profile", "catalog.vehicle_variant", "catalog.vehicle_model", "org.organization", "platform.tenant"} {
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
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'return-exchange','Return Exchange','Return Exchange')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store','return-exchange-store','Store','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','return-exchange-model','Model','motorcycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','return-exchange-variant','Variant','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name) values($1,'customer','Customer')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'original-1','store','variant','ORIGINAL-1','VIN-O1','BAT-O1','quarantine',2,clock_timestamp()),($1,'original-2','store','variant','ORIGINAL-2','VIN-O2','BAT-O2','quarantine',2,clock_timestamp()),($1,'replacement-1','store','variant','REPLACEMENT-1','VIN-R1','BAT-R1','available',7,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'original-order-1','store','customer','delivered','ARS',500000,4),($1,'original-order-2','store','customer','delivered','ARS',500000,4)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id) values($1,'original-order-1','original-line-1','variant',1,500000,'original-1'),($1,'original-order-2','original-line-2','variant',1,500000,'original-2')`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,acceptance_evidence_sha256_hex,customer_accepted_at,version) values($1,'original-handover-1','store','original-order-1','customer','original-1','accepted','aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',clock_timestamp(),2),($1,'original-handover-2','store','original-order-2','customer','original-2','accepted','bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb',clock_timestamp(),2)`,
		`insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version,resolved_at,resolved_by_subject) values($1,'exception-1','store','original-handover-1','customer','defect','Defect','cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc','resolved',2,clock_timestamp(),'operator'),($1,'exception-2','store','original-handover-2','customer','defect','Defect','dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd','resolved',2,clock_timestamp(),'operator')`,
		`insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject) values($1,'auth-1','store','exception-1','original-handover-1','original-order-1','original-1','customer','exchange','original-order-1','authorized','operator'),($1,'auth-2','store','exception-2','original-handover-2','original-order-2','original-2','customer','exchange','original-order-2','authorized','operator')`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject) values($1,'receipt-1','auth-1','store','original-order-1','original-1','customer','ORIGINAL-1','damaged','Received','eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee','operator'),($1,'receipt-2','auth-2','store','original-order-2','original-2','customer','ORIGINAL-2','damaged','Received','ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff','operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject) values($1,'disposition-1','receipt-1','quarantine','exchange','Replace','operator'),($1,'disposition-2','receipt-2','quarantine','exchange','Replace','operator')`,
	}
	for _, query := range fixtures {
		if _, err = tx.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	for _, set := range []struct{ suffix string }{{"1"}, {"2"}} {
		if _, err = pool.Exec(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key) values($1,$2,$3,'inventory','inventory','requested',$4),($1,$5,$3,'exchange','fulfillment','requested',$6),($1,$7,$3,'accounting','accounting','requested',$8)`, tenant, "inventory-"+set.suffix, "disposition-"+set.suffix, "inventory-request-key-000"+set.suffix, "exchange-"+set.suffix, "exchange-request-key-000"+set.suffix, "accounting-"+set.suffix, "accounting-request-key-000"+set.suffix); err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, `update sales.return_effect_execution set status='succeeded',result_sha256_hex=repeat('a',64) where tenant_id=$1 and request_id=$2`, tenant, "inventory-"+set.suffix); err != nil {
			t.Fatal(err)
		}
	}

	store := NewReturnExchange(pool)
	work, err := store.Claim(ctx, "fulfillment", "worker-a", "claim-a", time.Minute)
	if err != nil || work == nil || work.RequestID != "exchange-1" {
		t.Fatalf("work=%+v err=%v", work, err)
	}
	result, err := store.PrepareExchange(ctx, *work, "worker-a")
	if err != nil || result.ReplacementOrderID == "" || result.ReplacementStockID != "replacement-1" || len(result.ResultSHA256) != 64 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	var orderState, stockState, executionState, exchangeStatus string
	var total, stockVersion int64
	var handovers, events, attempts int
	if err = pool.QueryRow(ctx, `select state,total_minor_units from sales.customer_order where tenant_id=$1 and order_id=$2`, tenant, result.ReplacementOrderID).Scan(&orderState, &total); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select state,version from inventory.stock_unit where tenant_id=$1 and stock_unit_id='replacement-1'`, tenant).Scan(&stockState, &stockVersion); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='exchange-1'`, tenant).Scan(&executionState); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_exchange where tenant_id=$1 and request_id='exchange-1'`, tenant).Scan(&exchangeStatus); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.delivery_handover where tenant_id=$1 and handover_id=$2 and state='prepared'`, tenant, result.ReplacementHandoverID).Scan(&handovers); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and payload->>'request_id'='exchange-1'`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.return_effect_attempt where tenant_id=$1 and request_id='exchange-1' and outcome='succeeded'`, tenant).Scan(&attempts); err != nil {
		t.Fatal(err)
	}
	if orderState != "allocated" || total != 0 || stockState != "reserved" || stockVersion != 8 || executionState != "succeeded" || exchangeStatus != "prepared" || handovers != 1 || events != 3 || attempts != 1 {
		t.Fatalf("order=%s/%d stock=%s/%d execution=%s exchange=%s handovers=%d events=%d attempts=%d", orderState, total, stockState, stockVersion, executionState, exchangeStatus, handovers, events, attempts)
	}

	work, err = store.Claim(ctx, "fulfillment", "worker-b", "claim-b", time.Minute)
	if err != nil || work == nil || work.RequestID != "exchange-2" {
		t.Fatalf("second work=%+v err=%v", work, err)
	}
	if _, err = store.PrepareExchange(ctx, *work, "worker-b"); !errors.Is(err, returnexchange.ErrStockUnavailable) {
		t.Fatalf("no-stock err=%v", err)
	}
	if err = store.Finish(ctx, *work, "worker-b", returneffects.Completion{Outcome: "retry", ErrorCode: "REPLACEMENT_STOCK_UNAVAILABLE", RetryAfter: time.Minute}); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='exchange-2'`, tenant).Scan(&executionState); err != nil {
		t.Fatal(err)
	}
	if executionState != "retry" {
		t.Fatalf("second execution=%s", executionState)
	}
	if _, err = pool.Exec(ctx, `update sales.return_exchange set status='prepared' where tenant_id=$1 and request_id='exchange-1'`, tenant); err == nil {
		t.Fatal("exchange mutation accepted")
	}
	if work, err = store.Claim(ctx, "fulfillment", "worker-c", "claim-c", time.Minute); err != nil || work != nil {
		t.Fatalf("terminal or delayed work reclaimed: %+v %v", work, err)
	}
}
