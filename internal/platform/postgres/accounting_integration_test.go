package postgres

import (
	"context"
	"crypto/rand"
	"elite.local/enterprise/internal/accounting"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
	"testing"
	"time"
)

func TestAccountingPostingConcurrencyReversalAndClose(t *testing.T) {
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
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from accounting.entry where tenant_id=$1`, `delete from accounting.register where tenant_id=$1`, `delete from accounting.journal_line where tenant_id=$1`, `delete from accounting.journal where tenant_id=$1`, `delete from accounting.period where tenant_id=$1`, `delete from accounting.account where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			if _, e = tx.Exec(ctx, q, tenant); e != nil {
				t.Errorf("cleanup: %v", e)
				return
			}
		}
		if e = tx.Commit(ctx); e != nil {
			t.Errorf("cleanup commit: %v", e)
		}
	}()
	for _, q := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'acct-'||substring($1::text,1,8),'Accounting','Accounting')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise','Franchise','franchisee'),($1,'other','other','Other','franchisee')`} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewAccounting(pool)
	if err = repo.CreateAccount(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29e01", accounting.Account{Code: "CASH", Name: "Cash", Type: "asset"}); err != nil {
		t.Fatal(err)
	}
	if err = repo.CreateAccount(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29e02", accounting.Account{Code: "REVENUE", Name: "Revenue", Type: "revenue"}); err != nil {
		t.Fatal(err)
	}
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	until := from.AddDate(0, 1, 0)
	if err = repo.OpenPeriod(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29e03", accounting.Period{ID: "2026-08", StartsOn: from, EndsOn: until, Status: "open", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err = repo.OpenPeriod(ctx, tenant, "unused", accounting.Period{ID: "overlap", StartsOn: from.AddDate(0, 0, 15), EndsOn: until.AddDate(0, 0, 15), Status: "open", Version: 1}); !errors.Is(err, accounting.ErrConflict) {
		t.Fatalf("overlap accepted: %v", err)
	}
	journal := accounting.Journal{ID: "journal", OrganizationID: "franchise", PeriodID: "2026-08", SourceType: "SALE", SourceID: "order-1", Currency: "ARS", PostingDate: from.Add(12 * time.Hour), Status: "draft", TotalDebitMinorUnits: 10001, TotalCreditMinorUnits: 10001, Version: 1, Lines: []accounting.Line{{LineNo: 1, AccountCode: "CASH", Description: "cash", DebitMinorUnits: 10001}, {LineNo: 2, AccountCode: "REVENUE", Description: "sale", CreditMinorUnits: 10001}}}
	if err = repo.CreateJournal(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29e04", journal); err != nil {
		t.Fatal(err)
	}
	type result struct {
		journal accounting.Journal
		err     error
	}
	results := make(chan result, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			j, e := repo.PostJournal(ctx, tenant, "franchise", "journal", "controller", 1, []string{"register-a", "register-b"}[i], []string{"018f4d4a-7b36-7a21-8d10-2f4c54c29e05", "018f4d4a-7b36-7a21-8d10-2f4c54c29e06"}[i])
			results <- result{j, e}
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)
	success, conflict := 0, 0
	for got := range results {
		if got.err == nil {
			success++
		} else if errors.Is(got.err, accounting.ErrConflict) {
			conflict++
		} else {
			t.Fatal(got.err)
		}
	}
	if success != 1 || conflict != 1 {
		t.Fatalf("posting success=%d conflict=%d", success, conflict)
	}
	balances, err := repo.TrialBalance(ctx, tenant, "franchise", "2026-08")
	if err != nil || len(balances) != 2 {
		t.Fatalf("balances=%+v err=%v", balances, err)
	}
	reversal, err := repo.ReverseJournal(ctx, tenant, "franchise", "journal", "journal-reversal", "register-reversal", "018f4d4a-7b36-7a21-8d10-2f4c54c29e07", 2, "controller", "duplicate sale correction")
	if err != nil || reversal.ReversalOf != "journal" {
		t.Fatalf("reversal=%+v err=%v", reversal, err)
	}
	if _, err = repo.ReverseJournal(ctx, tenant, "franchise", "journal", "second", "register-second", "unused", 3, "controller", "again"); !errors.Is(err, accounting.ErrConflict) {
		t.Fatalf("second reversal accepted: %v", err)
	}
	balances, err = repo.TrialBalance(ctx, tenant, "franchise", "2026-08")
	if err != nil {
		t.Fatal(err)
	}
	for _, balance := range balances {
		if balance.NetMinorUnits != 0 {
			t.Fatalf("non-zero after reversal: %+v", balances)
		}
	}
	period, err := repo.ClosePeriod(ctx, tenant, "2026-08", 1, "018f4d4a-7b36-7a21-8d10-2f4c54c29e08")
	if err != nil || period.Status != "closed" || period.Version != 2 {
		t.Fatalf("period=%+v err=%v", period, err)
	}
	journal.ID = "after-close"
	journal.SourceID = "order-2"
	if err = repo.CreateJournal(ctx, tenant, "unused", journal); !errors.Is(err, accounting.ErrConflict) {
		t.Fatalf("journal in closed period accepted: %v", err)
	}
}
