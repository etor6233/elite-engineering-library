package postgres

import (
	"context"
	"crypto/rand"
	"elite.local/enterprise/internal/royalty"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRoyaltyLedgerSettlementConcurrencyReversalAndReconciliation(t *testing.T) {
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
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		if _, cleanupErr = tx.Exec(ctx, `set local session_replication_role = replica`); cleanupErr != nil {
			t.Errorf("cleanup role: %v", cleanupErr)
			return
		}
		for _, query := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from royalty.reconciliation where tenant_id=$1`, `delete from royalty.settlement_line where tenant_id=$1`, `delete from royalty.settlement_run where tenant_id=$1`, `delete from royalty.accrual where tenant_id=$1`, `delete from royalty.policy where tenant_id=$1`, `delete from payment.payment_attempt where tenant_id=$1`, `delete from sales.customer_order where tenant_id=$1`, `delete from franchise.agreement where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			if _, cleanupErr = tx.Exec(ctx, query, tenant); cleanupErr != nil {
				t.Errorf("cleanup query: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("cleanup commit: %v", cleanupErr)
		}
	}()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'royalty-'||substring($1::text,1,8),'Royalty','Royalty')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise','Franchise','franchisee'),($1,'other','other','Other','franchisee')`,
		`insert into franchise.agreement(tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,status)values($1,'agreement','franchise','AR-BUE-ROYALTY','v1','2026-01-01','active')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','franchise','customer','confirmed','ARS',10001,1)`,
		`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values($1,'payment','order','provider','provider-payment','idem-payment','captured','ARS',10001,3)`,
	}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewRoyalty(pool)
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	until := from.AddDate(0, 1, 0)
	if err = repo.CreatePolicy(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29c01", royalty.Policy{ID: "policy", AgreementID: "agreement", OrganizationID: "franchise", Currency: "ARS", RateBasisPoints: 650, ValidFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	if err = repo.CreatePolicy(ctx, tenant, "unused", royalty.Policy{ID: "overlap", AgreementID: "agreement", OrganizationID: "franchise", Currency: "ARS", RateBasisPoints: 700, ValidFrom: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)}); !errors.Is(err, royalty.ErrConflict) {
		t.Fatalf("overlapping policy accepted: %v", err)
	}
	accrual, err := repo.AccruePayment(ctx, tenant, "provider:capture-1", royalty.PaymentEvent{ID: "accrual-capture", OrganizationID: "franchise", PaymentAttemptID: "payment", PaymentExpectedVersion: 3, State: "captured", OccurredAt: from.Add(time.Hour)}, "018f4d4a-7b36-7a21-8d10-2f4c54c29c02")
	if err != nil || accrual.RoyaltyMinorUnits != 650 || accrual.BasisMinorUnits != 10001 {
		t.Fatalf("accrual=%+v err=%v", accrual, err)
	}
	if _, err = repo.AccruePayment(ctx, tenant, "provider:capture-1", royalty.PaymentEvent{ID: "accrual-duplicate", OrganizationID: "franchise", PaymentAttemptID: "payment", PaymentExpectedVersion: 3, State: "captured", OccurredAt: from.Add(time.Hour)}, "unused"); !errors.Is(err, royalty.ErrConflict) {
		t.Fatalf("duplicate accrual accepted: %v", err)
	}
	if _, err = repo.AccruePayment(ctx, tenant, "provider:cross-scope", royalty.PaymentEvent{ID: "accrual-cross", OrganizationID: "other", PaymentAttemptID: "payment", PaymentExpectedVersion: 3, State: "captured", OccurredAt: from.Add(time.Hour)}, "unused"); !errors.Is(err, royalty.ErrConflict) {
		t.Fatalf("cross-scope accrual accepted: %v", err)
	}
	settlement := royalty.Settlement{ID: "settlement", OrganizationID: "franchise", Currency: "ARS", PeriodStart: from, PeriodEnd: until, Status: "draft", Version: 1}
	if err = repo.OpenSettlement(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29c03", settlement); err != nil {
		t.Fatal(err)
	}
	type closeResult struct {
		value royalty.Settlement
		err   error
	}
	results := make(chan closeResult, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			v, e := repo.CloseSettlement(ctx, tenant, "franchise", "settlement", []string{"018f4d4a-7b36-7a21-8d10-2f4c54c29c04", "018f4d4a-7b36-7a21-8d10-2f4c54c29c05"}[i], 1)
			results <- closeResult{v, e}
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)
	successes, conflicts := 0, 0
	var closed royalty.Settlement
	for result := range results {
		if result.err == nil {
			successes++
			closed = result.value
		} else if errors.Is(result.err, royalty.ErrConflict) {
			conflicts++
		} else {
			t.Fatal(result.err)
		}
	}
	if successes != 1 || conflicts != 1 || closed.ExpectedMinorUnits != 650 || closed.Version != 2 {
		t.Fatalf("success=%d conflicts=%d closed=%+v", successes, conflicts, closed)
	}
	reversal, err := repo.ReverseSettlement(ctx, tenant, "franchise", "settlement", "settlement-reversal", "018f4d4a-7b36-7a21-8d10-2f4c54c29c06", 2, "contract correction")
	if err != nil || reversal.ExpectedMinorUnits != -650 || reversal.ReversalOf != "settlement" {
		t.Fatalf("reversal=%+v err=%v", reversal, err)
	}
	if _, err = repo.ReverseSettlement(ctx, tenant, "franchise", "settlement", "second-reversal", "unused", 3, "duplicate reversal"); !errors.Is(err, royalty.ErrConflict) {
		t.Fatalf("second reversal accepted: %v", err)
	}
	reconciliation, err := repo.Reconcile(ctx, tenant, "franchise", royalty.Reconciliation{ID: "reconciliation", SettlementID: "settlement-reversal", ExternalReference: "bank-line-1", ActualMinorUnits: -649}, "controller", "018f4d4a-7b36-7a21-8d10-2f4c54c29c07")
	if err != nil || reconciliation.Status != "mismatch" || reconciliation.DifferenceMinorUnits != 1 {
		t.Fatalf("reconciliation=%+v err=%v", reconciliation, err)
	}
	var lineTotal int64
	var lines int
	if err = pool.QueryRow(ctx, `select coalesce(sum(amount_minor_units),0),count(*) from royalty.settlement_line where tenant_id=$1 and settlement_id in ('settlement','settlement-reversal')`, tenant).Scan(&lineTotal, &lines); err != nil {
		t.Fatal(err)
	}
	if lineTotal != 0 || lines != 2 {
		t.Fatalf("compensating history total=%d lines=%d", lineTotal, lines)
	}
}
