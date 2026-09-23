package postgres

import (
	"context"
	"errors"
 "fmt"
 "strings"
 "elite.local/enterprise/internal/fiscal"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/returneffects"
	"elite.local/enterprise/internal/returnfiscal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReturnFiscalRequestsOnceAndWaitsForAuthorization(t *testing.T) {
 runReturnFiscalLocalClosure(t, false)
}
func TestReturnFiscalProcessorReconcilesLostAuthorization(t *testing.T) {
 runReturnFiscalLocalClosure(t, true)
}
func runReturnFiscalLocalClosure(t *testing.T, connected bool) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2e146"
	cleanup := func() {
		tx, e := pool.Begin(ctx)
		if e != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"fiscal.return_credit_note_link", "fiscal.invoice_associated_voucher", "fiscal.issuance_attempt", "fiscal.invoice_vat", "fiscal.invoice_other_tax", "fiscal.invoice", "fiscal.point_of_sale", "accounting.return_effect_posting", "payment.return_refund_observation", "payment.return_refund", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "payment.payment_attempt", "sales.customer_order", "platform.outbox_event", "org.organization", "platform.tenant"} {
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
	queries := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'return-fiscal','Return Fiscal','Return Fiscal')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','return-fiscal-store','Store','store')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',12100,4)`,
		`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values($1,'payment','order','stripe','pay-ref','return-fiscal-payment','captured','ARS',12100,2)`,
		`insert into fiscal.point_of_sale(tenant_id,point_of_sale_id,organization_id,taxpayer_cuit,environment,point_of_sale_number,active,version)values($1,'pos','store','30715117564','homologation',1,true,1)`,
		`insert into fiscal.invoice(tenant_id,invoice_id,organization_id,order_id,payment_attempt_id,point_of_sale_id,voucher_type,concept,recipient_document_type,recipient_document,recipient_vat_condition_id,currency,provider_currency,total_minor_units,net_minor_units,vat_minor_units,exempt_minor_units,non_taxed_minor_units,other_tax_minor_units,issued_on,status,voucher_number,cae,cae_expires_on,authorized_at,request_hash,idempotency_key,version)values($1,'original','store','order','payment','pos',6,1,99,'0',5,'ARS','PES',12100,10000,2100,0,0,0,'2026-08-15','authorized',1,'12345678901234','2026-08-31',clock_timestamp(),repeat('a',64),'original-fiscal-0001',2)`,
		`insert into fiscal.invoice_vat(tenant_id,invoice_id,vat_id,base_minor_units,amount_minor_units)values($1,'original',5,10000,2100)`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject)values($1,'receipt','auth','store','order','stock','customer','SERIAL','damaged','Received',repeat('b',64),'operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject,decided_at)values($1,'disposition','receipt','quarantine','refund','Refund','operator','2026-08-31T12:00:00Z')`,
		`insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values($1,'inventory','disposition','inventory','inventory','requested','inventory-effect-0001'),($1,'refund','disposition','refund','payment','requested','refund-effect-000001'),($1,'accounting','disposition','accounting','accounting','requested','accounting-effect-01'),($1,'fiscal','disposition','fiscal','fiscal','requested','fiscal-effect-000001')`,
		`insert into sales.return_effect_execution(tenant_id,request_id,status,result_sha256_hex)values($1,'inventory','succeeded',repeat('c',64)),($1,'refund','succeeded',repeat('d',64)),($1,'accounting','succeeded',repeat('e',64)),($1,'fiscal','requested',null)`,
		`insert into payment.return_refund(tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,provider_refund_reference,provider_status,state,response_sha256_hex)values($1,'refund','payment','order','line','stock','stripe','pay-ref','refund-provider-0001','ARS',12100,'refund-ref','succeeded','succeeded',repeat('f',64))`,
		`insert into accounting.return_effect_posting(tenant_id,request_id,disposition_id,source_order_id,remedy,original_journal_id,reversal_journal_id,currency,reversed_minor_units,status,result_sha256_hex)values($1,'accounting','disposition','order','refund','sale-journal','reversal-journal','ARS',12100,'posted',repeat('1',64))`,
	}
	for _, query := range queries {
		if _, err = tx.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	store := NewReturnFiscal(pool)
	work, err := store.Claim(ctx, "fiscal", "worker", "claim-1", time.Minute)
	if err != nil || work == nil {
		t.Fatalf("work=%+v err=%v", work, err)
	}
	result, err := store.RequestOrObserve(ctx, *work, "worker")
	if !errors.Is(err, returnfiscal.ErrPending) || result.CreditInvoiceID == "" || result.Status != "queued" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if err = store.Finish(ctx, *work, "worker", returneffects.Completion{Outcome: "retry", ErrorCode: "FISCAL_AUTHORIZATION_PENDING", ProviderReference: result.CreditInvoiceID, RetryAfter: time.Second}); err != nil {
		t.Fatal(err)
	}
 if connected {
  fixture:= &closureFiscalProvider{t:t,invoiceID:result.CreditInvoiceID}
  processor,e:=fiscal.NewProcessor(NewFiscal(pool),fixture,&closureFiscalIDs{},"credit-worker",time.Minute)
  if e!=nil{t.Fatal(e)}
  if _,e=processor.ProcessOne(ctx);!errors.Is(e,closureFiscalLostReply){t.Fatalf("expected lost authorization: %v",e)}
  var status string
  if e=pool.QueryRow(ctx,`select status from fiscal.invoice where tenant_id=$1 and invoice_id=$2`,tenant,result.CreditInvoiceID).Scan(&status);e!=nil||status!="reconcile_required"{t.Fatalf("uncertain status %s %v",status,e)}
  authorized,e:=processor.ProcessOne(ctx)
  if e!=nil||authorized.Status!="authorized"||authorized.VoucherNumber!=2||fixture.authorizes!=1||fixture.consults!=1||fixture.sequences!=1{t.Fatalf("authorize/reconcile %+v counts=%+v err=%v",authorized,fixture,e)}
  if _,e=processor.ProcessOne(ctx);!errors.Is(e,fiscal.ErrNoWork){t.Fatalf("authorized credit was dispatched again: %v",e)}
  var events int
  if e=pool.QueryRow(ctx,`select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='fiscal-invoice.authorized'`,tenant,result.CreditInvoiceID).Scan(&events);e!=nil||events!=1{t.Fatalf("events=%d err=%v",events,e)}
 } else {
	if _, err = pool.Exec(ctx, `update fiscal.invoice set status='authorized',voucher_number=2,cae='22345678901234',cae_expires_on='2026-09-10',authorized_at=clock_timestamp(),version=version+1 where tenant_id=$1 and invoice_id=$2`, tenant, result.CreditInvoiceID); err != nil {
		t.Fatal(err)
	}
 }
	if _, err = pool.Exec(ctx, `update sales.return_effect_execution set available_at=clock_timestamp() where tenant_id=$1 and request_id='fiscal'`, tenant); err != nil {
		t.Fatal(err)
	}
	work, err = store.Claim(ctx, "fiscal", "worker", "claim-2", time.Minute)
	if err != nil || work == nil {
		t.Fatalf("work2=%+v err=%v", work, err)
	}
	result2, err := store.RequestOrObserve(ctx, *work, "worker")
	if err != nil || result2.Status != "authorized" || result2.CreditInvoiceID != result.CreditInvoiceID {
		t.Fatalf("result2=%+v err=%v", result2, err)
	}
	if err = store.Finish(ctx, *work, "worker", returneffects.Completion{Outcome: "succeeded", ProviderReference: result2.CreditInvoiceID, ResultSHA256: result2.ResultSHA256}); err != nil {
		t.Fatal(err)
	}
	var links, credits, attempts int
	var state string
	if err = pool.QueryRow(ctx, `select (select count(*) from fiscal.return_credit_note_link where tenant_id=$1),(select count(*) from fiscal.invoice where tenant_id=$1 and voucher_type=8),(select count(*) from sales.return_effect_attempt where tenant_id=$1 and request_id='fiscal'),(select status from sales.return_effect_execution where tenant_id=$1 and request_id='fiscal')`, tenant).Scan(&links, &credits, &attempts, &state); err != nil {
		t.Fatal(err)
	}
	if links != 1 || credits != 1 || attempts != 2 || state != "succeeded" {
		t.Fatalf("links=%d credits=%d attempts=%d state=%s", links, credits, attempts, state)
	}
	if _, err = pool.Exec(ctx, `update fiscal.return_credit_note_link set total_minor_units=1 where tenant_id=$1`, tenant); err == nil {
		t.Fatal("immutable return credit link update accepted")
	}
}

// Synthetic provider boundary; no credentials or homologation network call.
var closureFiscalLostReply=errors.New("synthetic authorization reply lost")
type closureFiscalIDs struct{n int}
func(i *closureFiscalIDs)New()string{i.n++;return fmt.Sprintf("closure-fiscal-%d",i.n)}
type closureFiscalProvider struct{t *testing.T;invoiceID string;authorizes,consults,sequences int}
func(p *closureFiscalProvider)validate(v fiscal.Invoice){
 p.t.Helper()
 if v.ID!=p.invoiceID||v.VoucherType!=8||v.TotalMinorUnits!=12100||v.Currency!="ARS"||len(v.AssociatedVouchers)!=1||v.AssociatedVouchers[0].InvoiceID!="original"{p.t.Fatalf("wrong credit source: %+v",v)}
}
func(p *closureFiscalProvider)LastAuthorized(_ context.Context,v fiscal.Invoice)(int64,string,error){p.validate(v);p.sequences++;return 1,strings.Repeat("a",64),nil}
func(p *closureFiscalProvider)Authorize(_ context.Context,v fiscal.Invoice)(fiscal.Authorization,error){p.validate(v);p.authorizes++;if v.VoucherNumber!=2{p.t.Fatalf("wrong voucher: %d",v.VoucherNumber)};return fiscal.Authorization{},closureFiscalLostReply}
func(p *closureFiscalProvider)Consult(_ context.Context,v fiscal.Invoice)(fiscal.Authorization,error){
 p.validate(v);p.consults++;if p.authorizes!=1||v.VoucherNumber!=2{p.t.Fatalf("consult was not same uncertain authorization: %+v",v)}
 expires:=time.Date(2026,10,10,0,0,0,0,time.UTC)
 return fiscal.Authorization{Found:true,Authorized:true,CAE:"22345678901234",CAEExpiresOn:&expires,ResponseHash:strings.Repeat("b",64)},nil
}
