package workers

import (
	"context"
	db "elite.local/enterprise/internal/platform/postgres"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type fakeOutbox struct {
	events           []db.OutboxEvent
	marked, released int
	remaining        time.Duration
	remainingByEvent map[string]time.Duration
}

func (f *fakeOutbox) Claim(context.Context, string, time.Duration, int) ([]db.OutboxEvent, error) {
	return f.events, nil
}
func (f *fakeOutbox) RemainingLease(_ context.Context, _, event, _ string, attempt int) (time.Duration, error) {
	if attempt != 1 {
		return 0, errors.New("wrong outbox preflight generation")
	}
	if f.remainingByEvent != nil {
		return f.remainingByEvent[event], nil
	}
	return f.remaining, nil
}
func (f *fakeOutbox) MarkPublished(_ context.Context, _, _, _ string, attempt int) error {
	if attempt != 1 {
		return errors.New("wrong outbox completion generation")
	}
	f.marked++
	return nil
}
func (f *fakeOutbox) Release(_ context.Context, _, _, _, _ string, attempt int, _ time.Duration) error {
	if attempt != 1 {
		return errors.New("wrong outbox release generation")
	}
	f.released++
	return nil
}

type fakePublisher struct{ fail string }

func (f fakePublisher) Publish(_ context.Context, e db.OutboxEvent, key string) error {
	if key != e.EventID {
		panic("missing idempotency key")
	}
	if e.EventID == f.fail {
		return errors.New("down")
	}
	return nil
}

type publisherFunc func(context.Context, db.OutboxEvent, string) error

func (f publisherFunc) Publish(ctx context.Context, e db.OutboxEvent, key string) error {
	return f(ctx, e, key)
}

func TestOutboxProcessorRejectsExpiredBudget(t *testing.T) {
	out := &fakeOutbox{events: []db.OutboxEvent{{TenantID: "t", EventID: "old", Attempts: 1}}, remaining: 0}
	called := false
	publisher := publisherFunc(func(context.Context, db.OutboxEvent, string) error { called = true; return nil })
	n, err := (OutboxProcessor{Store: out, Publisher: publisher, WorkerID: "w", Lease: time.Second, RetryDelay: time.Second, BatchSize: 1}).ProcessOnce(context.Background())
	if called || n != 0 || err == nil || out.marked != 0 || out.released != 0 {
		t.Fatalf("expired claim reached publisher: called=%v n=%d err=%v", called, n, err)
	}
}

func TestOutboxProcessorBoundsPublisherContext(t *testing.T) {
	out := &fakeOutbox{events: []db.OutboxEvent{{TenantID: "t", EventID: "fresh", Attempts: 1}}, remaining: 50 * time.Millisecond}
	publisher := publisherFunc(func(ctx context.Context, _ db.OutboxEvent, _ string) error {
		deadline, ok := ctx.Deadline()
		if !ok || time.Until(deadline) > 50*time.Millisecond {
			t.Error("publisher context exceeds remaining lease")
		}
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if n, err := (OutboxProcessor{Store: out, Publisher: publisher, WorkerID: "w", Lease: time.Second, RetryDelay: time.Second, BatchSize: 1}).ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatalf("fresh publish n=%d err=%v", n, err)
	}
}

func TestOutboxProcessorRechecksLaterBatchEntry(t *testing.T) {
	out := &fakeOutbox{events: []db.OutboxEvent{{TenantID: "t", EventID: "fresh", Attempts: 1}, {TenantID: "t", EventID: "late", Attempts: 1}}, remainingByEvent: map[string]time.Duration{"fresh": time.Second, "late": 0}}
	calls := 0
	publisher := publisherFunc(func(context.Context, db.OutboxEvent, string) error { calls++; return nil })
	n, err := (OutboxProcessor{Store: out, Publisher: publisher, WorkerID: "w", Lease: time.Second, RetryDelay: time.Second, BatchSize: 2}).ProcessOnce(context.Background())
	if n != 1 || err == nil || calls != 1 || out.marked != 1 {
		t.Fatalf("late batch entry published: n=%d err=%v calls=%d", n, err, calls)
	}
}

func TestOutboxProcessorPostgresRetryKeepsEventKey(t *testing.T) {
	// Separate database from the postgres package: Claim is intentionally global,
	// and Go can run package suites concurrently. Never use TEST_DATABASE_URL.
	raw := os.Getenv("OUTBOX_PROCESSOR_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("OUTBOX_PROCESSOR_TEST_DATABASE_URL is not set; separate isolated database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_outbox_test_") || len(strings.TrimPrefix(cfg.ConnConfig.Database, "elite_outbox_test_")) < 16 {
		t.Fatal("processor requires dedicated loopback elite_outbox_test_<unique> database")
	}
	for _, fallback := range cfg.ConnConfig.Fallbacks {
		if fallback.Host != "127.0.0.1" {
			t.Fatal("non-loopback fallback forbidden")
		}
	}
	if other := os.Getenv("OUTBOX_TEST_DATABASE_URL"); other != "" {
		otherCfg, e := pgxpool.ParseConfig(other)
		if e == nil && otherCfg.ConnConfig.Host == cfg.ConnConfig.Host && otherCfg.ConnConfig.Port == cfg.ConnConfig.Port && otherCfg.ConnConfig.Database == cfg.ConnConfig.Database {
			t.Fatal("processor and store tests require distinct databases")
		}
	}
	cfg.MaxConns = 4
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	const tenant = "018f4d4a-7b36-7a21-8d10-2f4c54c28506"
	const first = "018f4d4a-7b36-7a21-8d10-2f4c54c28507"
	const retry = "018f4d4a-7b36-7a21-8d10-2f4c54c28508"
	cleanup := func() {
		clean, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, e := pool.Exec(clean, `delete from platform.outbox_event where tenant_id=$1`, tenant); e != nil {
			t.Error(e)
		}
		if _, e := pool.Exec(clean, `delete from platform.tenant where tenant_id=$1`, tenant); e != nil {
			t.Error(e)
		}
	}
	cleanup()
	defer cleanup()
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'processor-v316','Synthetic','Synthetic')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'test','a',1,'test.created',1,clock_timestamp(),'{}'),($1,$3,'test','b',1,'test.created',1,clock_timestamp(),'{}')`, tenant, first, retry); err != nil {
		t.Fatal(err)
	}
	store := db.NewOutbox(pool)
	seen := map[string][]int{}
	publisher := publisherFunc(func(publishCtx context.Context, event db.OutboxEvent, key string) error {
		deadline, ok := publishCtx.Deadline()
		if !ok || time.Until(deadline) > time.Minute || publishCtx.Err() != nil {
			t.Error("publisher lacks current bounded context")
		}
		if key != event.EventID || event.TenantID != tenant {
			t.Error("publisher event binding changed")
		}
		seen[key] = append(seen[key], event.Attempts)
		if event.EventID == retry && event.Attempts == 1 {
			return errors.New("synthetic transport failure")
		}
		return nil
	})
	processor := OutboxProcessor{Store: store, Publisher: publisher, WorkerID: "processor-v316", Lease: time.Minute, RetryDelay: 0, BatchSize: 2}
	if n, err := processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatalf("first batch n=%d err=%v", n, err)
	}
	var published bool
	var code string
	var attempts int
	if err = pool.QueryRow(ctx, `select published_at is not null,coalesce(last_error_code,''),attempts from platform.outbox_event where tenant_id=$1 and event_id=$2`, tenant, retry).Scan(&published, &code, &attempts); err != nil || published || code != "PUBLISH_FAILED" || attempts != 1 {
		t.Fatal("retry state not durable", published, code, attempts, err)
	}
	if n, err := processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatalf("retry batch n=%d err=%v", n, err)
	}
	if len(seen[first]) != 1 || len(seen[retry]) != 2 || seen[retry][0] != 1 || seen[retry][1] != 2 {
		t.Fatal("retry key/generation changed", seen)
	}
	var count int
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and published_at is not null and claimed_by is null and claimed_until is null and last_error_code is null`, tenant).Scan(&count); err != nil || count != 2 {
		t.Fatal("ack state incomplete", count, err)
	}
	t.Log("OUTBOX_PROCESSOR_POSTGRES_PASS published=2 stable_retry_key=1 generations=1,2 provider=synthetic")
}

type fakeJobs struct {
	jobs              []db.Job
	completed, failed int
	remaining         time.Duration
	leaseErr          error
	generation        int
}

func (f *fakeJobs) Claim(context.Context, string, string, time.Duration, int) ([]db.Job, error) {
	return f.jobs, nil
}
func (f *fakeJobs) RemainingLease(_ context.Context, _, _, _ string, attempt int) (time.Duration, error) {
	if attempt != f.generation {
		return 0, errors.New("wrong preflight generation")
	}
	return f.remaining, f.leaseErr
}
func (f *fakeJobs) Complete(_ context.Context, _, _, _ string, attempt int) error {
	if attempt != f.generation {
		return errors.New("wrong completion generation")
	}
	f.completed++
	return nil
}
func (f *fakeJobs) Fail(_ context.Context, _, _, _, _ string, attempt int, _ time.Duration) (bool, error) {
	if attempt != f.generation {
		return false, errors.New("wrong failure generation")
	}
	f.failed++
	return false, nil
}

type fakeHandler struct{ fail string }

func (f fakeHandler) Handle(_ context.Context, j db.Job) error {
	if j.JobID == f.fail {
		return errors.New("bad")
	}
	return nil
}
func TestProcessorsSeparateSuccessAndRetry(t *testing.T) {
	ctx := context.Background()
	out := &fakeOutbox{events: []db.OutboxEvent{{TenantID: "t", EventID: "ok", Attempts: 1}, {TenantID: "t", EventID: "retry", Attempts: 1}}, remaining: time.Minute}
	n, err := (OutboxProcessor{Store: out, Publisher: fakePublisher{fail: "retry"}, WorkerID: "w", Lease: time.Minute, RetryDelay: time.Second, BatchSize: 2}).ProcessOnce(ctx)
	if err != nil || n != 1 || out.marked != 1 || out.released != 1 {
		t.Fatalf("outbox n=%d marked=%d released=%d err=%v", n, out.marked, out.released, err)
	}
	jobs := &fakeJobs{remaining: time.Minute, generation: 7, jobs: []db.Job{{TenantID: "t", JobID: "ok", Attempts: 7}, {TenantID: "t", JobID: "retry", Attempts: 7}}}
	n, err = (JobProcessor{Store: jobs, Handler: fakeHandler{fail: "retry"}, Queue: "q", WorkerID: "w", Lease: time.Minute, RetryDelay: time.Second, BatchSize: 2}).ProcessOnce(ctx)
	if err != nil || n != 1 || jobs.completed != 1 || jobs.failed != 1 {
		t.Fatalf("jobs n=%d completed=%d failed=%d err=%v", n, jobs.completed, jobs.failed, err)
	}
}

type handlerFunc func(context.Context, db.Job) error

func (f handlerFunc) Handle(ctx context.Context, job db.Job) error { return f(ctx, job) }

func TestJobProcessorLeaseBudgetAndLoss(t *testing.T) {
	for _, mode := range []string{"lost", "expired", "later-batch-expired", "bounded", "cancelled"} {
		t.Run(mode, func(t *testing.T) {
			store := &fakeJobs{generation: 3, remaining: time.Minute, jobs: []db.Job{{TenantID: "t", JobID: "one", Attempts: 3}}}
			if mode == "lost" {
				store.leaseErr = db.ErrJobClaimLost
			}
			if mode == "expired" {
				store.remaining = 0
			}
			if mode == "later-batch-expired" {
				store.jobs = append(store.jobs, db.Job{TenantID: "t", JobID: "two", Attempts: 3})
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if mode == "cancelled" {
				cancel()
			}
			calls := 0
			handler := handlerFunc(func(ctx context.Context, job db.Job) error {
				calls++
				deadline, ok := ctx.Deadline()
				if !ok || time.Until(deadline) <= 0 || time.Until(deadline) > time.Minute {
					t.Fatal("missing bounded lease deadline")
				}
				if mode == "later-batch-expired" {
					store.remaining = 0
				}
				return nil
			})
			n, err := (JobProcessor{Store: store, Handler: handler, Queue: "q", WorkerID: "w", Lease: time.Minute, BatchSize: 2}).ProcessOnce(ctx)
			if mode == "bounded" {
				if err != nil || n != 1 || calls != 1 || store.completed != 1 {
					t.Fatalf("bounded: %d %d %v", n, calls, err)
				}
			} else if mode == "later-batch-expired" {
				if err == nil || n != 1 || calls != 1 || store.completed != 1 {
					t.Fatalf("later batch: %d %d %v", n, calls, err)
				}
			} else if err == nil || n != 0 || calls != 0 || store.completed != 0 || store.failed != 0 {
				t.Fatalf("stale handler/effect: %d %d %v", n, calls, err)
			}
		})
	}
}
