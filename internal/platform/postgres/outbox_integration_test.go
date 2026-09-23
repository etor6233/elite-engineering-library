package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOutboxClaimsAreDisjointAndOwned(t *testing.T) {
	ctx, pool := outboxTestPool(t)
	const tenantID = "018f4d4a-7b36-7a21-8d10-2f4c54c28b02"
	cleanup := func() {
		_, _ = pool.Exec(ctx, `delete from platform.outbox_event where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenantID)
	}
	cleanup()
	defer cleanup()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'outbox-test','Outbox Test S.A.','Outbox Test')`, tenantID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values ($1,'018f4d4a-7b36-7a21-8d10-2f4c54c28b10','test','a',1,'test.created',1,clock_timestamp(),'{}'),($1,'018f4d4a-7b36-7a21-8d10-2f4c54c28b11','test','b',1,'test.created',1,clock_timestamp(),'{}')`, tenantID); err != nil {
		t.Fatal(err)
	}
	store := NewOutbox(pool)
	first, err := store.Claim(ctx, "worker-a", time.Minute, 1)
	if err != nil || len(first) != 1 {
		t.Fatalf("first claim failed: events=%+v err=%v", first, err)
	}
	second, err := store.Claim(ctx, "worker-b", time.Minute, 1)
	if err != nil || len(second) != 1 || second[0].EventID == first[0].EventID {
		t.Fatalf("second claim was not disjoint: first=%+v second=%+v err=%v", first, second, err)
	}
	if err := store.MarkPublished(ctx, first[0].TenantID, first[0].EventID, "worker-b", first[0].Attempts); err == nil {
		t.Fatal("wrong worker marked an event")
	}
	if err := store.MarkPublished(ctx, first[0].TenantID, first[0].EventID, "worker-a", first[0].Attempts); err != nil {
		t.Fatal(err)
	}
	if err := store.Release(ctx, second[0].TenantID, second[0].EventID, "worker-b", "PROVIDER_TIMEOUT", second[0].Attempts, 0); err != nil {
		t.Fatal(err)
	}
	reclaimed, err := store.Claim(ctx, "worker-c", time.Minute, 1)
	if err != nil || len(reclaimed) != 1 || reclaimed[0].EventID != second[0].EventID {
		t.Fatalf("released event was not reclaimed: events=%+v err=%v", reclaimed, err)
	}
}

func outboxTestConfig(raw string) (*pgxpool.Config, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("explicit disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_outbox_test_") || len(strings.TrimPrefix(cfg.ConnConfig.Database, "elite_outbox_test_")) < 16 {
		return nil, errors.New("requires dedicated loopback elite_outbox_test_<unique> database")
	}
	for _, fallback := range cfg.ConnConfig.Fallbacks {
		if fallback.Host != "127.0.0.1" {
			return nil, errors.New("non-loopback fallback forbidden")
		}
	}
	cfg.MaxConns = 6
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second
	return cfg, nil
}

func outboxTestPool(t *testing.T) (context.Context, *pgxpool.Pool) {
	t.Helper()
	raw := os.Getenv("OUTBOX_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("OUTBOX_TEST_DATABASE_URL is not set; isolated database required")
	}
	cfg, err := outboxTestConfig(raw)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err = pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	return ctx, pool
}

func outboxLeaseFixture(t *testing.T) (context.Context, *pgxpool.Pool, *Outbox) {
	t.Helper()
	ctx, pool := outboxTestPool(t)
	const tenant = "018f4d4a-7b36-7a21-8d10-2f4c54c28504"
	cleanup := func() {
		clean, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := pool.Exec(clean, `delete from platform.outbox_event where tenant_id=$1`, tenant); err != nil {
			t.Error(err)
		}
		if _, err := pool.Exec(clean, `delete from platform.tenant where tenant_id=$1`, tenant); err != nil {
			t.Error(err)
		}
	}
	cleanup()
	t.Cleanup(cleanup)
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'outbox-v316','Synthetic','Synthetic')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,'018f4d4a-7b36-7a21-8d10-2f4c54c28505','test','v316',1,'test.created',1,clock_timestamp(),'{}')`, tenant); err != nil {
		t.Fatal(err)
	}
	return ctx, pool, NewOutbox(pool)
}

func finishOutboxClaim(ctx context.Context, store *Outbox, event OutboxEvent, release bool) error {
	if release {
		return store.Release(ctx, event.TenantID, event.EventID, "same-worker", "PUBLISH_FAILED", event.Attempts, 0)
	}
	return store.MarkPublished(ctx, event.TenantID, event.EventID, "same-worker", event.Attempts)
}

func TestOutboxExpiredClaimCannotFinalize(t *testing.T) {
	for _, release := range []bool{false, true} {
		t.Run(map[bool]string{false: "publish", true: "release"}[release], func(t *testing.T) {
			ctx, pool, store := outboxLeaseFixture(t)
			events, err := store.Claim(ctx, "same-worker", time.Minute, 1)
			if err != nil || len(events) != 1 {
				t.Fatalf("claim: %v %d", err, len(events))
			}
			e := events[0]
			if _, err = pool.Exec(ctx, `update platform.outbox_event set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and event_id=$2`, e.TenantID, e.EventID); err != nil {
				t.Fatal(err)
			}
			if err = finishOutboxClaim(ctx, store, e, release); err == nil {
				t.Fatal("expired claim finalized")
			}
		})
	}
}

func TestOutboxReclaimedGenerationCannotFinalize(t *testing.T) {
	for _, release := range []bool{false, true} {
		t.Run(map[bool]string{false: "publish", true: "release"}[release], func(t *testing.T) {
			ctx, pool, store := outboxLeaseFixture(t)
			first, err := store.Claim(ctx, "same-worker", time.Minute, 1)
			if err != nil || len(first) != 1 {
				t.Fatal("first claim", err)
			}
			e := first[0]
			if _, err = pool.Exec(ctx, `update platform.outbox_event set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and event_id=$2`, e.TenantID, e.EventID); err != nil {
				t.Fatal(err)
			}
			second, err := store.Claim(ctx, "same-worker", time.Minute, 1)
			if err != nil || len(second) != 1 || second[0].EventID != e.EventID || second[0].Attempts != e.Attempts+1 {
				t.Fatal("reclaim", err)
			}
			if err = finishOutboxClaim(ctx, store, e, release); err == nil {
				t.Fatal("stale generation finalized reassigned work")
			}
			if err = finishOutboxClaim(ctx, store, second[0], release); err != nil {
				t.Fatal("current generation failed", err)
			}
		})
	}
}

func TestOutboxFinalizationChecksClockAfterLock(t *testing.T) {
	for _, release := range []bool{false, true} {
		t.Run(map[bool]string{false: "publish", true: "release"}[release], func(t *testing.T) {
			ctx, pool, store := outboxLeaseFixture(t)
			events, err := store.Claim(ctx, "same-worker", 500*time.Millisecond, 1)
			if err != nil || len(events) != 1 {
				t.Fatal("claim", err)
			}
			e := events[0]
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(context.Background())
			if _, err = tx.Exec(ctx, `select event_id from platform.outbox_event where tenant_id=$1 and event_id=$2 for update`, e.TenantID, e.EventID); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { done <- finishOutboxClaim(ctx, store, e, release) }()
			deadline := time.Now().Add(2 * time.Second)
			blocked := false
			for time.Now().Before(deadline) {
				// Observe from fresh transactions: pg_stat_activity can retain a
				// statistics snapshot for the lifetime of the lock-holding tx.
				if err = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and pid<>pg_backend_pid())`).Scan(&blocked); err != nil {
					t.Fatal(err)
				}
				if blocked {
					break
				}
				time.Sleep(5 * time.Millisecond)
			}
			if !blocked {
				t.Fatal("finalization did not wait on row lock")
			}
			expired := false
			for time.Now().Before(deadline) {
				if err = tx.QueryRow(ctx, `select claimed_until<=clock_timestamp() from platform.outbox_event where tenant_id=$1 and event_id=$2`, e.TenantID, e.EventID).Scan(&expired); err != nil {
					t.Fatal(err)
				}
				if expired {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			if !expired {
				t.Fatal("lease did not expire while blocked")
			}
			if err = tx.Commit(ctx); err != nil {
				t.Fatal(err)
			}
			select {
			case err = <-done:
				if err == nil {
					t.Fatal("claim finalized after lock wait exceeded lease")
				}
			case <-time.After(2 * time.Second):
				t.Fatal("finalization did not return")
			}
		})
	}
}

func TestOutboxDatabaseGuard(t *testing.T) {
	for _, raw := range []string{"", "postgres://postgres@127.0.0.1/main?sslmode=disable", "postgres://postgres@localhost/elite_outbox_test_0123456789abcdef?sslmode=disable", "postgres://postgres@192.0.2.1/elite_outbox_test_0123456789abcdef?sslmode=disable", "postgres://postgres@127.0.0.1/elite_outbox_test_short?sslmode=disable", "host=127.0.0.1,192.0.2.1 dbname=elite_outbox_test_0123456789abcdef sslmode=disable"} {
		if _, err := outboxTestConfig(raw); err == nil {
			t.Fatal("unsafe fixture target accepted")
		}
	}
	if _, err := outboxTestConfig("postgres://postgres@127.0.0.1/elite_outbox_test_0123456789abcdef?sslmode=disable"); err != nil {
		t.Fatal(err)
	}
}

func TestOutboxRemainingLeaseRequiresCurrentOwner(t *testing.T) {
	ctx, pool, store := outboxLeaseFixture(t)
	events, err := store.Claim(ctx, "same-worker", time.Minute, 1)
	if err != nil || len(events) != 1 {
		t.Fatal("claim", err)
	}
	e := events[0]
	if remaining, err := store.RemainingLease(ctx, e.TenantID, e.EventID, "same-worker", e.Attempts); err != nil || remaining <= 0 || remaining > time.Minute {
		t.Fatal("invalid remaining lease", remaining, err)
	}
	for _, tc := range []struct {
		worker  string
		attempt int
	}{{"other-worker", e.Attempts}, {"same-worker", 0}, {"same-worker", e.Attempts + 1}} {
		if _, err = store.RemainingLease(ctx, e.TenantID, e.EventID, tc.worker, tc.attempt); !errors.Is(err, ErrOutboxClaimLost) {
			t.Fatal("invalid owner/generation retained lease", err)
		}
	}
	if _, err = pool.Exec(ctx, `update platform.outbox_event set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and event_id=$2`, e.TenantID, e.EventID); err != nil {
		t.Fatal(err)
	}
	if _, err = store.RemainingLease(ctx, e.TenantID, e.EventID, "same-worker", e.Attempts); !errors.Is(err, ErrOutboxClaimLost) {
		t.Fatal("expired lease retained budget", err)
	}
	second, err := store.Claim(ctx, "same-worker", time.Minute, 1)
	if err != nil || len(second) != 1 {
		t.Fatal("reclaim", err)
	}
	if _, err = store.RemainingLease(ctx, e.TenantID, e.EventID, "same-worker", e.Attempts); !errors.Is(err, ErrOutboxClaimLost) {
		t.Fatal("stale generation retained budget", err)
	}
	if err = finishOutboxClaim(ctx, store, second[0], false); err != nil {
		t.Fatal(err)
	}
	if _, err = store.RemainingLease(ctx, e.TenantID, e.EventID, "same-worker", second[0].Attempts); !errors.Is(err, ErrOutboxClaimLost) {
		t.Fatal("published event retained budget", err)
	}
}
