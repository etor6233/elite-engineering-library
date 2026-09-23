package postgres

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestFiscalPostgresSourceBindingLeaseReconciliationAndImmutability(t *testing.T) {
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
	var raw [16]byte
	if _, err = rand.Read(raw[:]); err != nil {
		t.Fatal(err)
	}
	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80
	tenant := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", raw[0:4], raw[4:6], raw[6:8], raw[8:10], raw[10:16])
	defer func() {
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Errorf("cleanup begin: %v", e)
			return
		}
		defer tx.Rollback(ctx)
		if _, e = tx.Exec(ctx, `set local session_replication_role=replica`); e != nil {
			t.Errorf("cleanup role: %v", e)
			return
		}
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from fiscal.issuance_attempt where tenant_id=$1`, `delete from fiscal.invoice_associated_voucher where tenant_id=$1`, `delete from fiscal.invoice_vat where tenant_id=$1`, `delete from fiscal.invoice_other_tax where tenant_id=$1`, `delete from fiscal.invoice where tenant_id=$1`, `delete from fiscal.point_of_sale where tenant_id=$1`, `delete from payment.payment_attempt where tenant_id=$1`, `delete from sales.customer_order where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			if _, e = tx.Exec(ctx, q, tenant); e != nil {
				t.Errorf("cleanup: %v", e)
				return
			}
		}
		if e = tx.Commit(ctx); e != nil {
			t.Errorf("cleanup commit: %v", e)
		}
	}()
	fixtures := []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'fiscal-'||substring($1::text,1,8),'Fiscal','Fiscal')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise','Franchise','franchisee'),($1,'other','other','Other','franchisee')`, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','franchise','customer','paid','ARS',12100,2),($1,'wrong-order','franchise','customer','paid','ARS',12000,1)`, `insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values($1,'payment','order','provider','reference','payment-fiscal','captured','ARS',12100,3),($1,'wrong-payment','wrong-order','provider','wrong-reference','wrong-payment-fiscal','captured','ARS',12000,1)`}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFiscal(pool)
	pos := fiscal.PointOfSale{ID: "pos", OrganizationID: "franchise", TaxpayerCUIT: "30715117564", Environment: "homologation", Number: 1, Active: true, Version: 1}
	if err = repo.ConfigurePointOfSale(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c2a001", pos); err != nil {
		t.Fatal(err)
	}
	request := fiscal.Invoice{ID: "invoice", OrganizationID: "franchise", OrderID: "order", PaymentAttemptID: "payment", PointOfSaleID: "pos", VoucherType: 6, Concept: 1, RecipientDocumentType: 99, RecipientDocument: "0", RecipientVATConditionID: 5, NetMinorUnits: 10000, VATMinorUnits: 2100, VATLines: []fiscal.VATLine{{ID: 5, BaseMinorUnits: 10000, AmountMinorUnits: 2100}}, IssuedOn: time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC), Status: "queued", Version: 1}
	requestHash := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	created, replayed, err := repo.RequestInvoice(ctx, tenant, "1234567890abcdef", requestHash, "018f4d4a-7b36-7a21-8d10-2f4c54c2a002", request)
	if err != nil || replayed || created.TotalMinorUnits != 12100 {
		t.Fatalf("created=%+v replay=%v err=%v", created, replayed, err)
	}
	if len(created.VATLines) != 1 || created.VATLines[0].ID != 5 {
		t.Fatalf("VAT detail not retained: %+v", created.VATLines)
	}
	if created.RecipientVATConditionID != 5 {
		t.Fatalf("recipient VAT condition not retained: %+v", created)
	}
	if _, replayed, err = repo.RequestInvoice(ctx, tenant, "1234567890abcdef", requestHash, "unused", request); err != nil || !replayed {
		t.Fatalf("replay=%v err=%v", replayed, err)
	}
	request.ID = "wrong"
	request.PaymentAttemptID = "wrong-payment"
	if _, _, err = repo.RequestInvoice(ctx, tenant, "1234567890abcdeg", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "unused", request); !errors.Is(err, fiscal.ErrConflict) {
		t.Fatalf("source mismatch accepted: %v", err)
	}
	claimed, err := repo.ClaimInvoice(ctx, "worker", time.Minute, "attempt-claim")
	if err != nil || claimed.TenantID != tenant || claimed.ClaimedFrom != "queued" {
		t.Fatalf("claimed=%+v err=%v", claimed, err)
	}
	if _, err = repo.ClaimInvoice(ctx, "other-worker", time.Minute, "attempt-other"); !errors.Is(err, fiscal.ErrNoWork) {
		t.Fatalf("second lane claim accepted: %v", err)
	}
	claimed, err = repo.AssignVoucherNumber(ctx, claimed, "worker", 1, "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc", "attempt-sequence")
	if err != nil || claimed.VoucherNumber != 1 || claimed.Status != "authorizing" {
		t.Fatalf("assigned=%+v err=%v", claimed, err)
	}
	if err = repo.DeferInvoice(ctx, claimed, "worker", "authorize-timeout", "", true); err != nil {
		t.Fatal(err)
	}
	claimed, err = repo.ClaimInvoice(ctx, "worker", time.Minute, "attempt-reconcile")
	if err != nil || claimed.ClaimedFrom != "reconcile_required" || claimed.VoucherNumber != 1 {
		t.Fatalf("reclaimed=%+v err=%v", claimed, err)
	}
	expires := time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)
	finished, err := repo.FinishInvoice(ctx, claimed, "worker", fiscal.Authorization{Found: true, Authorized: true, CAE: "12345678901234", CAEExpiresOn: &expires, ResponseHash: "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", ProviderCodes: nil}, "attempt-finish")
	if err != nil || finished.Status != "authorized" || finished.CAE != "12345678901234" {
		t.Fatalf("finished=%+v err=%v", finished, err)
	}
	creditRequest := fiscal.Invoice{ID: "credit", OrganizationID: "franchise", OrderID: "order", PaymentAttemptID: "payment", PointOfSaleID: "pos", VoucherType: 8, Concept: 1, RecipientDocumentType: 99, RecipientDocument: "0", RecipientVATConditionID: 5, NetMinorUnits: 10000, VATMinorUnits: 2100, VATLines: []fiscal.VATLine{{ID: 5, BaseMinorUnits: 10000, AmountMinorUnits: 2100}}, AssociatedVouchers: []fiscal.AssociatedVoucher{{InvoiceID: "invoice"}}, IssuedOn: time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC), Status: "queued", Version: 1}
	credit, replayed, err := repo.RequestInvoice(ctx, tenant, "1234567890credit", "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "018f4d4a-7b36-7a21-8d10-2f4c54c2a003", creditRequest)
	if err != nil || replayed || len(credit.AssociatedVouchers) != 1 || credit.AssociatedVouchers[0].VoucherType != 6 || credit.AssociatedVouchers[0].PointOfSale != 1 || credit.AssociatedVouchers[0].Number != 1 || credit.AssociatedVouchers[0].TaxpayerCUIT != pos.TaxpayerCUIT {
		t.Fatalf("credit=%+v replay=%v err=%v", credit, replayed, err)
	}
	credit, err = repo.ClaimInvoice(ctx, "credit-worker", time.Minute, "attempt-credit-claim")
	if err != nil || credit.ID != "credit" || len(credit.AssociatedVouchers) != 1 || credit.AssociatedVouchers[0].InvoiceID != "invoice" {
		t.Fatalf("claimed credit=%+v err=%v", credit, err)
	}
	if _, err = pool.Exec(ctx, `update fiscal.issuance_attempt set outcome='deferred' where tenant_id=$1 and attempt_id='attempt-claim'`, tenant); err == nil {
		t.Fatal("immutable attempt update accepted")
	}
	var attempts, authorizedEvents int
	if err = pool.QueryRow(ctx, `select count(*) from fiscal.issuance_attempt where tenant_id=$1`, tenant).Scan(&attempts); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='fiscal-invoice.authorized'`, tenant).Scan(&authorizedEvents); err != nil {
		t.Fatal(err)
	}
	if attempts != 6 || authorizedEvents != 1 {
		t.Fatalf("attempts=%d authorized_events=%d", attempts, authorizedEvents)
	}
}
