package postgres

import (
	"context"
	"elite.local/enterprise/internal/returnaccounting"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type businessCancelQuery struct {
	needle string
	fired  atomic.Bool
}

func (f *businessCancelQuery) TraceQueryStart(ctx context.Context, _ *pgx.Conn, d pgx.TraceQueryStartData) context.Context {
	if strings.Contains(d.SQL, f.needle) && f.fired.CompareAndSwap(false, true) {
		c, cancel := context.WithCancel(ctx)
		cancel()
		return c
	}
	return ctx
}
func (*businessCancelQuery) TraceQueryEnd(context.Context, *pgx.Conn, pgx.TraceQueryEndData) {}

type businessAccountingIDs struct{ value string }

func (i businessAccountingIDs) New() string { return i.value }

func TestBusinessAccountingScanErrorRecovery(t *testing.T) {
	for _, needle := range []string{"select r.status,r.total_debit_minor_units,count(*)", "select attempt_count,claimed_at from sales.return_effect_execution"} {
		t.Run(needle, func(t *testing.T) {
			ctx, pool, tenant := businessAccountingFixture(t)
			if _, err := pool.Exec(ctx, `update sales.return_effect_execution set available_at=clock_timestamp()+interval '1 hour' where tenant_id=$1 and request_id='accounting-2'`, tenant); err != nil {
				t.Fatal(err)
			}
			cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
			if err != nil {
				t.Fatal(err)
			}
			fault := &businessCancelQuery{needle: needle}
			cfg.ConnConfig.Tracer = fault
			traced, err := pgxpool.NewWithConfig(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}
			defer traced.Close()
			store := NewReturnAccounting(traced)
			processor, err := returnaccounting.NewProcessor(store, businessAccountingIDs{"business-claim"}, "business-worker", time.Minute, time.Second)
			if err != nil {
				t.Fatal(err)
			}
			_, err = processor.ProcessOne(ctx)
			if !fault.fired.Load() || err == nil || errors.Is(err, returnaccounting.ErrConflict) {
				t.Fatalf("transient swallowed: err=%v fired=%v", err, fault.fired.Load())
			}
			var status, code string
			if err = pool.QueryRow(ctx, `select status,last_error_code from sales.return_effect_execution where tenant_id=$1 and request_id='accounting-1'`, tenant).Scan(&status, &code); err != nil || status != "retry" || code != "ACCOUNTING_TRANSIENT_FAILURE" {
				t.Fatalf("persisted %s/%s err=%v", status, code, err)
			}
			if _, err = pool.Exec(ctx, `update sales.return_effect_execution set available_at=clock_timestamp() where tenant_id=$1 and request_id='accounting-1'`, tenant); err != nil {
				t.Fatal(err)
			}
			recovery, _ := returnaccounting.NewProcessor(store, businessAccountingIDs{"business-retry"}, "business-worker", time.Minute, time.Second)
			result, err := recovery.ProcessOne(ctx)
			if err != nil || result.RequestID != "accounting-1" {
				t.Fatalf("recovery=%+v err=%v", result, err)
			}
			var count int
			if err = pool.QueryRow(ctx, `select count(*) from accounting.journal where tenant_id=$1 and reversal_of='sale-1'`, tenant).Scan(&count); err != nil || count != 1 {
				t.Fatalf("reversals=%d err=%v", count, err)
			}
			if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='accounting-1'`, tenant).Scan(&status); err != nil || status != "succeeded" {
				t.Fatalf("recovered status=%s err=%v", status, err)
			}
		})
	}
}

func TestBusinessReturnClaimWaitsForPrerequisites(t *testing.T) {
	ctx, pool, tenant := businessAccountingFixture(t)
	// Fixtures establish only claim eligibility; they do not prove a complete return journey.
	if _, err := pool.Exec(ctx, `update sales.return_effect_execution set status='requested',result_sha256_hex=null where tenant_id=$1 and request_id like 'inventory-%'`, tenant); err != nil {
		t.Fatal(err)
	}
	store := NewReturnEffects(pool)
	work, err := store.Claim(ctx, "accounting", "eligibility", "claim-before", time.Minute)
	if err != nil || work != nil {
		t.Fatalf("claimed accounting before inventory: %+v %v", work, err)
	}
	var attempts int
	if err = pool.QueryRow(ctx, `select sum(attempt_count) from sales.return_effect_execution where tenant_id=$1 and request_id like 'accounting-%'`, tenant).Scan(&attempts); err != nil || attempts != 0 {
		t.Fatalf("premature attempts=%d %v", attempts, err)
	}
	if _, err = pool.Exec(ctx, `update sales.return_effect_execution set status='succeeded',result_sha256_hex=repeat('a',64) where tenant_id=$1 and request_id='inventory-1'`, tenant); err != nil {
		t.Fatal(err)
	}
	work, err = store.Claim(ctx, "accounting", "eligibility", "claim-after", time.Minute)
	if err != nil || work == nil || work.RequestID != "accounting-1" {
		t.Fatalf("eligible accounting=%+v %v", work, err)
	}
}
func businessAccountingFixture(t *testing.T) (context.Context, *pgxpool.Pool, string) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, url)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(pool.Close)
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2e144"
	cleanup := func() {
		tx, err := pool.Begin(ctx)
		if err != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"accounting.return_effect_posting", "sales.return_effect_resume", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "sales.return_authorization", "sales.delivery_exception", "sales.delivery_handover", "sales.customer_order_line", "sales.customer_order", "platform.outbox_event", "accounting.entry", "accounting.register", "accounting.journal_line", "accounting.journal", "accounting.period", "accounting.account", "inventory.stock_unit", "crm.customer_profile", "catalog.vehicle_variant", "catalog.vehicle_model", "org.organization", "platform.tenant"} {
			_, _ = tx.Exec(ctx, `delete from `+table+` where tenant_id=$1`, tenant)
		}
		_ = tx.Commit(ctx)
	}
	cleanup()
	t.Cleanup(cleanup)
	tx, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	defer tx.Rollback(ctx)
	_, e = tx.Exec(ctx, `set local session_replication_role=replica`)
	if e != nil {
		t.Fatal(e)
	}
	queries := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'return-accounting','Return Accounting','Return Accounting')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','return-accounting-store','Store','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','ra-model','Model','motorcycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','ra-variant','Variant','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name)values($1,'customer','Customer')`,
		`insert into accounting.account(tenant_id,account_code,display_name,account_type)values($1,'CASH','Cash','asset'),($1,'REVENUE','Revenue','revenue')`,
		`insert into accounting.period(tenant_id,period_id,starts_on,ends_on,status,version)values($1,'2026-08','2026-08-01','2026-09-01','open',1)`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at)values($1,'stock-1','store','variant','SERIAL-1','VIN-1','BAT-1','quarantine',2,clock_timestamp()),($1,'stock-2','store','variant','SERIAL-2','VIN-2','BAT-2','quarantine',2,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order-1','store','customer','delivered','ARS',10000,4),($1,'order-2','store','customer','delivered','ARS',20000,4)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id)values($1,'order-1','line-1','variant',1,10000,'stock-1'),($1,'order-2','line-2','variant',1,20000,'stock-2')`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,acceptance_evidence_sha256_hex,customer_accepted_at,version)values($1,'hand-1','store','order-1','customer','stock-1','accepted',repeat('a',64),clock_timestamp(),2),($1,'hand-2','store','order-2','customer','stock-2','accepted',repeat('b',64),clock_timestamp(),2)`,
		`insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version,resolved_at,resolved_by_subject)values($1,'ex-1','store','hand-1','customer','defect','Defect',repeat('c',64),'resolved',2,clock_timestamp(),'operator'),($1,'ex-2','store','hand-2','customer','defect','Defect',repeat('d',64),'resolved',2,clock_timestamp(),'operator')`,
		`insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject)values($1,'auth-1','store','ex-1','hand-1','order-1','stock-1','customer','return','order-1','authorized','operator'),($1,'auth-2','store','ex-2','hand-2','order-2','stock-2','customer','return','order-2','authorized','operator')`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject)values($1,'receipt-1','auth-1','store','order-1','stock-1','customer','SERIAL-1','damaged','Received',repeat('e',64),'operator'),($1,'receipt-2','auth-2','store','order-2','stock-2','customer','SERIAL-2','damaged','Received',repeat('f',64),'operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject)values($1,'disp-1','receipt-1','quarantine','refund','Return','operator'),($1,'disp-2','receipt-2','quarantine','refund','Return','operator')`,
	}
	for _, q := range queries {
		if _, e = tx.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	for _, v := range []struct {
		id     string
		amount int64
	}{{"1", 10000}, {"2", 20000}} {
		if _, e = tx.Exec(ctx, `insert into accounting.journal(tenant_id,journal_id,organization_id,period_id,source_type,source_id,currency,posting_date,status,total_debit_minor_units,total_credit_minor_units,version,posted_by,posted_at)values($1,$2,'store','2026-08','SALE',$3,'ARS','2026-08-15','posted',$4,$4,2,'controller',clock_timestamp())`, tenant, "sale-"+v.id, "order-"+v.id, v.amount); e != nil {
			t.Fatal(e)
		}
		if _, e = tx.Exec(ctx, `insert into accounting.journal_line(tenant_id,journal_id,line_no,account_code,description,debit_minor_units,credit_minor_units)values($1,$2,1,'CASH','Sale cash',$3,0),($1,$2,2,'REVENUE','Sale revenue',0,$3)`, tenant, "sale-"+v.id, v.amount); e != nil {
			t.Fatal(e)
		}
		if _, e = tx.Exec(ctx, `insert into accounting.register(tenant_id,register_id,journal_id,organization_id,period_id,posted_by,total_debit_minor_units,total_credit_minor_units)values($1,$2,$3,'store','2026-08','controller',$4,$4)`, tenant, "sale-register-"+v.id, "sale-"+v.id, v.amount); e != nil {
			t.Fatal(e)
		}
		if _, e = tx.Exec(ctx, `insert into accounting.entry(tenant_id,entry_id,register_id,journal_id,line_no,organization_id,period_id,account_code,currency,posting_date,description,debit_minor_units,credit_minor_units)values($1,$2,$3,$4,1,'store','2026-08','CASH','ARS','2026-08-15','Sale cash',$5,0),($1,$6,$3,$4,2,'store','2026-08','REVENUE','ARS','2026-08-15','Sale revenue',0,$5)`, tenant, "sale-entry-"+v.id+"-1", "sale-register-"+v.id, "sale-"+v.id, v.amount, "sale-entry-"+v.id+"-2"); e != nil {
			t.Fatal(e)
		}
	}
	if e = tx.Commit(ctx); e != nil {
		t.Fatal(e)
	}
	for _, id := range []string{"1", "2"} {
		if _, e = pool.Exec(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values($1,$2,$3,'inventory','inventory','requested',$4),($1,$5,$3,'refund','payment','requested',$6),($1,$7,$3,'accounting','accounting','requested',$8)`, tenant, "inventory-"+id, "disp-"+id, "inventory-key-000000"+id, "refund-"+id, "refund-key-000000000"+id, "accounting-"+id, "accounting-key-000000"+id); e != nil {
			t.Fatal(e)
		}
		if _, e = pool.Exec(ctx, `update sales.return_effect_execution set status='succeeded',result_sha256_hex=repeat('a',64) where tenant_id=$1 and request_id in($2,$3)`, tenant, "inventory-"+id, "refund-"+id); e != nil {
			t.Fatal(e)
		}
	}

	return ctx, pool, tenant
}

func TestBusinessReturnClaimStages(t *testing.T) {
	for _, tc := range []struct {
		owner, kind string
		required    []string
	}{
		{"fulfillment", "exchange", []string{"inventory"}},
		{"accounting", "accounting", []string{"inventory", "refund"}},
		{"fiscal", "fiscal", []string{"inventory", "refund", "accounting"}},
	} {
		t.Run(tc.owner, func(t *testing.T) {
			ctx, pool, tenant := businessAccountingFixture(t)
			if tc.kind == "exchange" {
				tx, err := pool.Begin(ctx)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback(ctx)
				if _, err = tx.Exec(ctx, `set local session_replication_role=replica`); err != nil {
					t.Fatal(err)
				}
				if _, err = tx.Exec(ctx, `update sales.return_disposition set customer_remedy='exchange' where tenant_id=$1 and disposition_id='disp-1'`, tenant); err != nil {
					t.Fatal(err)
				}
				if err = tx.Commit(ctx); err != nil {
					t.Fatal(err)
				}
			}
			request := tc.kind + "-1"
			if tc.kind != "accounting" {
				if _, err := pool.Exec(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values($1,$2,'disp-1',$3,$4,'requested',$5)`, tenant, request, tc.kind, tc.owner, "eligibility-key-"+tc.kind); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := pool.Exec(ctx, `update sales.return_effect_execution set status='requested',result_sha256_hex=null where tenant_id=$1`, tenant); err != nil {
				t.Fatal(err)
			}
			store := NewReturnEffects(pool)
			for _, prior := range tc.required {
				work, err := store.Claim(ctx, tc.owner, "stage-worker", "early-"+prior, time.Minute)
				if err != nil || work != nil {
					t.Fatalf("claimed %s before %s: %+v %v", tc.kind, prior, work, err)
				}
				if _, err = pool.Exec(ctx, `update sales.return_effect_execution set status='succeeded',result_sha256_hex=repeat('a',64) where tenant_id=$1 and request_id=$2`, tenant, prior+"-1"); err != nil {
					t.Fatal(err)
				}
			}
			work, err := store.Claim(ctx, tc.owner, "stage-worker", "ready-stage", time.Minute)
			if err != nil || work == nil || work.RequestID != request || work.Attempt != 1 {
				t.Fatalf("eligible=%+v %v", work, err)
			}
		})
	}
}
