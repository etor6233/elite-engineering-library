package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/returneffects"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReturnEffectsInventoryClaimFencingAndResume(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2e101"
	cleanup := func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"sales.return_effect_resume", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "sales.return_authorization", "platform.outbox_event", "inventory.stock_unit", "org.organization", "platform.tenant"} {
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
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'return-effects','Return Effects','Return Effects')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store','return-effects-store','Store','store')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'stock-1','store','variant','SERIAL-1','VIN-1','BATTERY-1','sold',1,clock_timestamp()),($1,'stock-2','store','variant','SERIAL-2','VIN-2','BATTERY-2','sold',1,clock_timestamp())`,
		`insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject) values($1,'auth-1','store','exception-1','handover-1','order-1','stock-1','customer','return','order-1','authorized','operator'),($1,'auth-2','store','exception-2','handover-2','order-2','stock-2','customer','return','order-2','authorized','operator')`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject) values($1,'receipt-1','auth-1','store','order-1','stock-1','customer','SERIAL-1','damaged','received','aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa','operator'),($1,'receipt-2','auth-2','store','order-2','stock-2','customer','SERIAL-2','opened','received','bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb','operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject) values($1,'disposition-1','receipt-1','quarantine','refund','inspect','operator'),($1,'disposition-2','receipt-2','restock','refund','inspect','operator')`,
	}
	for _, query := range fixtures {
		if _, err = tx.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	if _, err = pool.Exec(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key) values($1,'effect-1','disposition-1','inventory','inventory','requested','return-effect-key-0001'),($1,'effect-2','disposition-2','inventory','inventory','requested','return-effect-key-0002')`, tenant); err != nil {
		t.Fatal(err)
	}
	store := NewReturnEffects(pool)

	type claimResult struct {
		work *returneffects.Work
		err  error
	}
	start := make(chan struct{})
	claims := make(chan claimResult, 2)
	for _, input := range []struct{ worker, token string }{{"worker-a", "claim-a"}, {"worker-b", "claim-b"}} {
		go func(worker, token string) {
			<-start
			work, claimErr := store.Claim(ctx, "inventory", worker, token, time.Minute)
			claims <- claimResult{work, claimErr}
		}(input.worker, input.token)
	}
	close(start)
	claimed := map[string]*returneffects.Work{}
	for range 2 {
		result := <-claims
		if result.err != nil || result.work == nil {
			t.Fatalf("claim=%+v err=%v", result.work, result.err)
		}
		if _, exists := claimed[result.work.RequestID]; exists {
			t.Fatal("same request claimed concurrently")
		}
		claimed[result.work.RequestID] = result.work
	}
	if len(claimed) != 2 {
		t.Fatalf("claimed=%d", len(claimed))
	}

	first := claimed["effect-1"]
	if first == nil {
		t.Fatal("effect-1 not claimed")
	}
	firstWorker := "worker-a"
	if first.ClaimToken == "claim-b" {
		firstWorker = "worker-b"
	}
	if _, err = store.ApplyInventory(ctx, *first, "wrong-worker"); !errors.Is(err, returneffects.ErrInventoryState) {
		t.Fatalf("wrong worker applied effect: %v", err)
	}
	result, err := store.ApplyInventory(ctx, *first, firstWorker)
	if err != nil || result.InventoryState != "quarantine" || result.StockVersion != 2 || len(result.ResultSHA256) != 64 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	var state, executionStatus string
	var version int64
	var attempts, events int
	if err = pool.QueryRow(ctx, `select state,version from inventory.stock_unit where tenant_id=$1 and stock_unit_id='stock-1'`, tenant).Scan(&state, &version); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='effect-1'`, tenant).Scan(&executionStatus); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.return_effect_attempt where tenant_id=$1 and request_id='effect-1'`, tenant).Scan(&attempts); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_type='stock-unit' and aggregate_id='stock-1' and event_type='stock-unit.return-quarantine' and payload->>'request_id'='effect-1'`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if state != "quarantine" || version != 2 || executionStatus != "succeeded" || attempts != 1 || events != 1 {
		t.Fatalf("state=%s version=%d status=%s attempts=%d events=%d", state, version, executionStatus, attempts, events)
	}
	if work, claimErr := store.Claim(ctx, "inventory", "worker-c", "claim-c", time.Minute); claimErr != nil || work != nil {
		t.Fatalf("completed work reclaimed: %+v %v", work, claimErr)
	}

	second := claimed["effect-2"]
	if second == nil {
		t.Fatal("effect-2 not claimed")
	}
	secondWorker := "worker-a"
	if second.ClaimToken == "claim-b" {
		secondWorker = "worker-b"
	}
	if err = store.Finish(ctx, *second, secondWorker, returneffects.Completion{Outcome: "blocked", ErrorCode: "INVENTORY_STATE_CONFLICT"}); err != nil {
		t.Fatal(err)
	}
	if err = store.ResumeBlocked(ctx, tenant, "effect-2", "resume-1", "operator", "INVENTORY_RECONCILED"); err != nil {
		t.Fatal(err)
	}
	reclaimed, err := store.Claim(ctx, "inventory", "worker-c", "claim-c", time.Minute)
	if err != nil || reclaimed == nil || reclaimed.RequestID != "effect-2" || reclaimed.Attempt != 2 {
		t.Fatalf("reclaimed=%+v err=%v", reclaimed, err)
	}
	if err = store.Finish(ctx, *second, secondWorker, returneffects.Completion{Outcome: "failed", ErrorCode: "STALE_CLAIM"}); !errors.Is(err, returneffects.ErrInventoryState) {
		t.Fatalf("stale claim completed: %v", err)
	}
	result, err = store.ApplyInventory(ctx, *reclaimed, "worker-c")
	if err != nil || result.InventoryState != "available" || result.StockVersion != 2 {
		t.Fatalf("restock result=%+v err=%v", result, err)
	}
	if _, err = pool.Exec(ctx, `update sales.return_effect_attempt set error_code='MUTATED' where tenant_id=$1 and request_id='effect-2'`, tenant); err == nil {
		t.Fatal("attempt audit mutation accepted")
	}
	if _, err = pool.Exec(ctx, `delete from sales.return_effect_resume where tenant_id=$1 and request_id='effect-2'`, tenant); err == nil {
		t.Fatal("resume audit deletion accepted")
	}
}
