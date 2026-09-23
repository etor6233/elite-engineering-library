package postgres

import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
	"time"
)

// Synthetic owner records seed a disposable database. Replica mode is fixture-only.
// Reads and the one inventory worker transition execute with normal triggers.
func TestCommercialCareOutcomePostgres(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	ctx := context.Background()
	p, e := pgxpool.New(ctx, dsn)
	if e != nil {
		t.Fatal(e)
	}
	defer p.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2c900"
	fixture := func(queries ...string) {
		t.Helper()
		tx, e := p.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		defer tx.Rollback(ctx)
		if _, e = tx.Exec(ctx, `set local session_replication_role=replica`); e != nil {
			t.Fatal(e)
		}
		for _, q := range queries {
			if _, e = tx.Exec(ctx, q, tenant); e != nil {
				t.Fatalf("fixture failed: %v query=%s", e, q)
			}
		}
		if e = tx.Commit(ctx); e != nil {
			t.Fatal(e)
		}
	}
	fixture(
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'commercial-care','Synthetic Care','Synthetic Care')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','care-store','Store','store')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at)values($1,'stock','store','variant','SERIAL-CARE','VIN-CARE','BAT-CARE','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',10000,4)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,acceptance_evidence_sha256_hex,customer_accepted_at,version)values($1,'hand','store','order','customer','stock','accepted',repeat('a',64),clock_timestamp(),2)`,
		`insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version,resolved_at,resolved_by_subject)values($1,'exception','store','hand','customer','defect','Synthetic defect',repeat('b',64),'resolved',2,clock_timestamp(),'operator')`,
		`insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject)values($1,'auth','store','exception','hand','order','stock','customer','return','order','authorized','operator')`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject)values($1,'receipt','auth','store','order','stock','customer','SERIAL-CARE','damaged','Synthetic receipt',repeat('c',64),'operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject)values($1,'disp','receipt','quarantine','refund','Synthetic disposition','operator')`,
		`insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values($1,'inventory','disp','inventory','inventory','requested','care-inventory-key'),($1,'refund','disp','refund','payment','requested','care-refund-key000'),($1,'accounting','disp','accounting','accounting','requested','care-account-key0'),($1,'fiscal','disp','fiscal','fiscal','requested','care-fiscal-key00')`,
		`insert into sales.return_effect_execution(tenant_id,request_id,status)select tenant_id,request_id,'requested' from sales.return_effect_request where tenant_id=$1`)
	repo := NewFranchiseJourney(p)
	read := func() franchisejourney.ReturnOutcome {
		t.Helper()
		v, e := repo.ReturnOutcome(ctx, tenant, "store", "auth")
		if e != nil {
			t.Fatal(e)
		}
		return v
	}
	if v := read(); len(v.Stages) != 4 || v.Stages[0].Status != "requested" || v.Stages[1].Owner != nil {
		t.Fatalf("initial projection %v", v)
	}
	t.Run("scope", func(t *testing.T) {
		for _, scope := range [][2]string{{"018f4d4a-7b36-7a21-8d10-2f4c54c2c901", "store"}, {tenant, "other"}} {
			if _, e := repo.ReturnOutcome(ctx, scope[0], scope[1], "auth"); !errors.Is(e, franchisejourney.ErrNotFound) {
				t.Fatalf("wrong scope e=%v", e)
			}
		}
	})
	t.Run("actualInventoryWorker", func(t *testing.T) {
		effects := NewReturnEffects(p)
		work, e := effects.Claim(ctx, "inventory", "care-worker", "care-claim", time.Minute)
		if e != nil || work == nil {
			t.Fatalf("claim %v %v", work, e)
		}
		if _, e = effects.ApplyInventory(ctx, *work, "care-worker"); e != nil {
			t.Fatal(e)
		}
		v := read()
		if v.Stages[0].Status != "succeeded" || v.Stages[0].ResultSHA256 == "" || v.Stages[1].Status != "requested" {
			t.Fatalf("actual transition missing %v", v)
		}
	})
	fixture(`insert into payment.return_refund(tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,provider_refund_reference,provider_status,state,response_sha256_hex)values($1,'refund','attempt','order','line','stock','stripe','pi_fixture','care-provider-key0','ARS',10000,'re_fixture','pending','pending',repeat('d',64))`,
		`update sales.return_effect_execution set status='retry',last_error_code='PROVIDER_PENDING' where tenant_id=$1 and request_id='refund'`)
	t.Run("pendingProviderIsNotSuccess", func(t *testing.T) {
		v := read()
		if v.Stages[1].Status != "retry" || v.Stages[1].Owner.State != "pending" || v.Stages[1].Owner.AmountMinorUnits != "10000" {
			t.Fatalf("pending erased %v", v)
		}
	})
	fixture(`update payment.return_refund set state='succeeded',provider_status='succeeded',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and request_id='refund'`,
		`update sales.return_effect_execution set status='succeeded',last_error_code=null,result_sha256_hex=repeat('d',64) where tenant_id=$1 and request_id='refund'`,
		`insert into accounting.return_effect_posting(tenant_id,request_id,disposition_id,source_order_id,remedy,original_journal_id,reversal_journal_id,currency,reversed_minor_units,status,result_sha256_hex)values($1,'accounting','disp','order','refund','original','reversal','ARS',10000,'posted',repeat('e',64))`,
		`update sales.return_effect_execution set status='succeeded',result_sha256_hex=repeat('e',64) where tenant_id=$1 and request_id='accounting'`,
		`insert into fiscal.invoice(tenant_id,invoice_id,organization_id,order_id,payment_attempt_id,point_of_sale_id,voucher_type,concept,recipient_document_type,recipient_document,recipient_vat_condition_id,currency,provider_currency,total_minor_units,net_minor_units,vat_minor_units,exempt_minor_units,non_taxed_minor_units,other_tax_minor_units,issued_on,status,request_hash,idempotency_key,version)values($1,'credit','store','order','attempt','pos',3,1,99,'',5,'ARS','PES',10000,10000,0,0,0,0,'2026-09-14','queued',repeat('f',64),'care-credit-key00',1)`,
		`insert into fiscal.return_credit_note_link(tenant_id,request_id,disposition_id,original_invoice_id,credit_invoice_id,currency,total_minor_units,request_sha256_hex)values($1,'fiscal','disp','original-invoice','credit','ARS',10000,repeat('f',64))`,
		`update sales.return_effect_execution set status='retry',last_error_code='FISCAL_PENDING' where tenant_id=$1 and request_id='fiscal'`)
	t.Run("queuedCreditIsNotAuthorization", func(t *testing.T) {
		v := read()
		if v.Stages[3].Status != "retry" || v.Stages[3].Owner.State != "queued" || v.Stages[2].Owner.Reference != "reversal" {
			t.Fatalf("owner records lost %v", v)
		}
	})
	fixture(`update fiscal.invoice set status='authorized',voucher_number=1,cae='00000000',cae_expires_on='2026-09-30',authorized_at=clock_timestamp() where tenant_id=$1 and invoice_id='credit'`, `update sales.return_effect_execution set status='succeeded',last_error_code=null,result_sha256_hex=repeat('f',64) where tenant_id=$1 and request_id='fiscal'`)
	t.Run("allStoredOwnersAndReadOnly", func(t *testing.T) {
		counts := func() string {
			var s string
			if e := p.QueryRow(ctx, `select jsonb_build_array((select count(*)from platform.outbox_event where tenant_id=$1),(select sum(attempt_count)from sales.return_effect_execution where tenant_id=$1),(select count(*)from sales.return_effect_attempt where tenant_id=$1),(select sum(version)from payment.return_refund where tenant_id=$1),(select sum(version)from fiscal.invoice where tenant_id=$1))::text`, tenant).Scan(&s); e != nil {
				t.Fatal(e)
			}
			return s
		}
		before := counts()
		for i := 0; i < 3; i++ {
			for _, x := range read().Stages {
				if x.Status != "succeeded" {
					t.Fatalf("not confirmed %v", x)
				}
			}
		}
		if after := counts(); after != before {
			t.Fatalf("GET mutated %s => %s", before, after)
		}
	})
	fixture(`delete from payment.return_refund where tenant_id=$1`, `delete from sales.return_effect_execution where tenant_id=$1 and request_id in('refund','fiscal')`, `delete from sales.return_effect_request where tenant_id=$1 and request_id in('refund','fiscal')`, `update sales.return_authorization set disposition='exchange' where tenant_id=$1 and authorization_id='auth'`, `update sales.return_disposition set customer_remedy='exchange' where tenant_id=$1 and disposition_id='disp'`,
		`insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values($1,'exchange','disp','exchange','fulfillment','requested','care-exchange-key0')`,
		`insert into sales.return_effect_execution(tenant_id,request_id,status,result_sha256_hex)values($1,'exchange','succeeded',repeat('a',64))`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'replacement-hand','store','replacement-order','customer','replacement-stock','prepared',1)`,
		`insert into sales.return_exchange(tenant_id,request_id,disposition_id,original_order_id,original_line_id,original_stock_unit_id,replacement_order_id,replacement_line_id,replacement_stock_unit_id,replacement_handover_id,accounting_request_id,currency,original_unit_price_minor_units,settlement_mode,status,result_sha256_hex)values($1,'exchange','disp','order','line','stock','replacement-order','replacement-line','replacement-stock','replacement-hand','accounting','ARS',10000,'even-exchange-zero-balance','prepared',repeat('a',64))`)
	t.Run("preparedExchangeIsNotDeliveryAcceptance", func(t *testing.T) {
		v := read()
		if len(v.Stages) != 3 || v.Stages[1].Kind != "exchange" || v.Stages[1].Owner.State != "prepared" || v.Stages[1].Owner.HandoverState != "prepared" || v.Stages[1].Owner.Reference != "replacement-hand" {
			t.Fatalf("exchange meaning lost %v", v)
		}
	})
	t.Run("missingWorkerEvidenceIsUnavailable", func(t *testing.T) {
		fixture(`update sales.return_effect_execution set status='requested',result_sha256_hex=null where tenant_id=$1 and request_id='exchange'`, `delete from sales.return_effect_execution where tenant_id=$1 and request_id='exchange'`)
		if _, e := repo.ReturnOutcome(ctx, tenant, "store", "auth"); !errors.Is(e, franchisejourney.ErrConflict) {
			t.Fatalf("missing worker shown as empty/pending: %v", e)
		}
	})
}
