package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestInboxAtomicDedupAndRollback(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28c01"
	event := "018f4d4a-7b36-7a21-8d10-2f4c54c28c02"
	_, _ = pool.Exec(ctx, `delete from platform.consumer_inbox where tenant_id=$1`, tenant)
	_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	defer func() {
		_, _ = pool.Exec(ctx, `delete from platform.consumer_inbox where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'inbox-test','Inbox Test','Inbox')`, tenant); err != nil {
		t.Fatal(err)
	}
	inbox := NewInbox(pool)
	calls := 0
	failing := func(context.Context, pgx.Tx) (string, error) { calls++; return "", errors.New("handler failed") }
	if _, err := inbox.Process(ctx, tenant, "inventory", event, failing); err == nil {
		t.Fatal("expected handler failure")
	}
	var count int
	if err := pool.QueryRow(ctx, `select count(*) from platform.consumer_inbox where tenant_id=$1`, tenant).Scan(&count); err != nil || count != 0 {
		t.Fatalf("failed handler left receipt: count=%d err=%v", count, err)
	}
	handler := func(ctx context.Context, tx pgx.Tx) (string, error) {
		calls++
		_, err := tx.Exec(ctx, `update platform.tenant set display_name='Inbox Applied' where tenant_id=$1`, tenant)
		return "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", err
	}
	processed, err := inbox.Process(ctx, tenant, "inventory", event, handler)
	if err != nil || !processed {
		t.Fatalf("first process: %v %v", processed, err)
	}
	processed, err = inbox.Process(ctx, tenant, "inventory", event, handler)
	if err != nil || processed {
		t.Fatalf("duplicate process: %v %v", processed, err)
	}
	if calls != 2 {
		t.Fatalf("handler calls=%d want 2 (one rollback, one commit)", calls)
	}
}
