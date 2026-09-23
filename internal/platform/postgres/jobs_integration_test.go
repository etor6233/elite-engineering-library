package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"testing"
	"time"
)

// This fault-injection suite only uses a dedicated local test database.
func TestJobsExpiredAndReusedWorkerClaims(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_jobs_") {
		t.Fatal("dedicated local elite_jobs_ database required")
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	for _, takeover := range []bool{false, true} {
		for _, operation := range []string{"complete", "fail"} {
			name := operation
			if takeover {
				name += "-reused-worker"
			} else {
				name += "-expired"
			}
			t.Run(name, func(t *testing.T) {
				var tenant, job string
				if err := pool.QueryRow(ctx, `select gen_random_uuid()::text,gen_random_uuid()::text`).Scan(&tenant, &job); err != nil {
					t.Fatal(err)
				}
				if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Synthetic lease test','Lease')`, tenant, "lease-"+job); err != nil {
					t.Fatal(err)
				}
				if _, err := pool.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts) values($1,$2,$3,'lease.test',1,'{}',3)`, tenant, job, job); err != nil {
					t.Fatal(err)
				}
				store := NewJobs(pool)
				first, err := store.Claim(ctx, job, "same-worker", time.Minute, 1)
				if err != nil || len(first) != 1 {
					t.Fatalf("first claim: %v %v", first, err)
				}
				// Deterministic DB-clock expiry, not a timing-sensitive sleep.
				if _, err := pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, job); err != nil {
					t.Fatal(err)
				}
				if takeover {
					second, err := store.Claim(ctx, job, "same-worker", time.Minute, 1)
					if err != nil || len(second) != 1 || second[0].Attempts != 2 {
						t.Fatalf("takeover: %v %v", second, err)
					}
				}
				if _, err := store.RemainingLease(ctx, tenant, job, "same-worker", first[0].Attempts); err == nil {
					t.Fatal("stale preflight accepted")
				}
				if operation == "complete" {
					err = store.Complete(ctx, tenant, job, "same-worker", first[0].Attempts)
				} else {
					_, err = store.Fail(ctx, tenant, job, "same-worker", "STALE", first[0].Attempts, 0)
				}
				if err == nil {
					t.Fatal("stale execution changed job state")
				}
				var untouched bool
				if err := pool.QueryRow(ctx, `select completed_at is null and terminal_error_code is null and claimed_by='same-worker' from platform.job where tenant_id=$1 and job_id=$2`, tenant, job).Scan(&untouched); err != nil || !untouched {
					t.Fatalf("lost current claim: %v %v", untouched, err)
				}
				if takeover {
					if remaining, err := store.RemainingLease(ctx, tenant, job, "same-worker", 2); err != nil || remaining <= 0 || remaining > time.Minute {
						t.Fatalf("current lease %v %v", remaining, err)
					}
					if operation == "complete" {
						err = store.Complete(ctx, tenant, job, "same-worker", 2)
					} else {
						_, err = store.Fail(ctx, tenant, job, "same-worker", "CURRENT", 2, 0)
					}
					if err != nil {
						t.Fatalf("current generation rejected: %v", err)
					}
				}
			})
		}
	}
}

func TestJobsOwnershipRetryAndTerminal(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28d01"
	jobID := "018f4d4a-7b36-7a21-8d10-2f4c54c28d02"
	_, _ = pool.Exec(ctx, `delete from platform.job where tenant_id=$1`, tenant)
	_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	defer func() {
		_, _ = pool.Exec(ctx, `delete from platform.job where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	_, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'jobs-test','Jobs Test','Jobs')`, tenant)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts)values($1,$2,'email','welcome',1,'{}',2)`, tenant, jobID)
	if err != nil {
		t.Fatal(err)
	}
	store := NewJobs(pool)
	first, err := store.Claim(ctx, "email", "worker-a", time.Minute, 1)
	if err != nil || len(first) != 1 || first[0].Attempts != 1 {
		t.Fatalf("first claim %+v %v", first, err)
	}
	if err := store.Complete(ctx, tenant, jobID, "worker-b", first[0].Attempts); err == nil {
		t.Fatal("wrong worker completed job")
	}
	terminal, err := store.Fail(ctx, tenant, jobID, "worker-a", "TEMP", first[0].Attempts, 0)
	if err != nil || terminal {
		t.Fatalf("first failure terminal=%v err=%v", terminal, err)
	}
	second, err := store.Claim(ctx, "email", "worker-b", time.Minute, 1)
	if err != nil || len(second) != 1 || second[0].Attempts != 2 {
		t.Fatalf("second claim %+v %v", second, err)
	}
	terminal, err = store.Fail(ctx, tenant, jobID, "worker-b", "EXHAUSTED", second[0].Attempts, 0)
	if err != nil || !terminal {
		t.Fatalf("terminal failure=%v err=%v", terminal, err)
	}
	third, err := store.Claim(ctx, "email", "worker-c", time.Minute, 1)
	if err != nil || len(third) != 0 {
		t.Fatalf("terminal job reclaimed %+v %v", third, err)
	}
}

func TestJobsExpiryWhileWaitingForRowLock(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_jobs_") {
		t.Fatal("dedicated local elite_jobs_ database required")
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	for _, operation := range []string{"complete", "fail"} {
		t.Run(operation, func(t *testing.T) {
			var tenant, job string
			if err := pool.QueryRow(ctx, `select gen_random_uuid()::text,gen_random_uuid()::text`).Scan(&tenant, &job); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Synthetic lock wait','Lease')`, tenant, "wait-"+job); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts) values($1,$2,$3,'lease.test',1,'{}',3)`, tenant, job, job); err != nil {
				t.Fatal(err)
			}
			store := NewJobs(pool)
			claimed, err := store.Claim(ctx, job, "wait-worker", time.Second, 1)
			if err != nil || len(claimed) != 1 {
				t.Fatalf("claim %v %v", claimed, err)
			}
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(ctx)
			if _, err := tx.Exec(ctx, `select job_id from platform.job where tenant_id=$1 and job_id=$2 for update`, tenant, job); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() {
				if operation == "complete" {
					done <- store.Complete(ctx, tenant, job, "wait-worker", 1)
				} else {
					_, err := store.Fail(ctx, tenant, job, "wait-worker", "DELAYED", 1, 0)
					done <- err
				}
			}()
			// Confirm actual server lock contention, not merely scheduling a goroutine.
			for {
				var waiting bool
				if err := pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and pid<>pg_backend_pid() and state='active' and wait_event_type='Lock' and query like 'with owned as materialized (%')`).Scan(&waiting); err != nil {
					t.Fatal(err)
				}
				if waiting {
					break
				}
				select {
				case err := <-done:
					t.Fatalf("transition did not wait: %v", err)
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				case <-time.After(5 * time.Millisecond):
				}
			}
			// Hold the row lock through expiry without changing the row version.
			if _, err := tx.Exec(ctx, `select pg_sleep(greatest(0,extract(epoch from (claimed_until-clock_timestamp())))+0.02) from platform.job where tenant_id=$1 and job_id=$2`, tenant, job); err != nil {
				t.Fatal(err)
			}
			if err := tx.Commit(ctx); err != nil {
				t.Fatal(err)
			}
			if err := <-done; err == nil {
				t.Fatal("lease expired during lock wait but transition succeeded")
			}
			var untouched bool
			if err := pool.QueryRow(ctx, `select completed_at is null and terminal_error_code is null and claimed_by='wait-worker' from platform.job where tenant_id=$1 and job_id=$2`, tenant, job).Scan(&untouched); err != nil || !untouched {
				t.Fatalf("changed expired job: %v %v", untouched, err)
			}
		})
	}
}
