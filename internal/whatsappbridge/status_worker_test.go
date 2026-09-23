package whatsappbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/platform/postgres"
)

func TestStatusWorkerRejectsUninitializedRouter(t *testing.T) {
	for _, router := range []*StatusRouter{nil, {}} {
		if _, err := NewStatusWorker(router, "worker-main", time.Minute, 0); err == nil {
			t.Fatal("uninitialized router accepted")
		}
	}
}

func TestStatusWorkerScopeCompletionAndAudit(t *testing.T) {
	r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	var original postgres.Job
	if err := pool.QueryRow(ctx, `select tenant_id,job_id,queue,job_type,schema_version,payload,attempts,max_attempts from platform.job where tenant_id=$1`, p.TenantID).Scan(&original.TenantID, &original.JobID, &original.Queue, &original.JobType, &original.SchemaVersion, &original.Payload, &original.Attempts, &original.MaxAttempts); err != nil {
		t.Fatal(err)
	}
	var otherTenant string
	if err := pool.QueryRow(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values(gen_random_uuid(),'scope-'||gen_random_uuid()::text,'Synthetic scope','Scope') returning tenant_id`).Scan(&otherTenant); err != nil {
		t.Fatal(err)
	}
	var excluded []string
	for _, mode := range []string{"tenant", "queue", "job-type", "schema", "connection", "provider", "event-type"} {
		var payload map[string]string
		if json.Unmarshal(original.Payload, &payload) != nil {
			t.Fatal("fixture payload")
		}
		tenant, queue, kind, version := p.TenantID, original.Queue, original.JobType, 1
		switch mode {
		case "tenant":
			tenant = otherTenant
		case "queue":
			queue = "other"
		case "job-type":
			kind = "other"
		case "schema":
			version = 2
		case "connection":
			payload["connection_id"] = "other"
		case "provider":
			payload["provider_code"] = "other"
		case "event-type":
			payload["event_type"] = "other"
		}
		raw, _ := json.Marshal(payload)
		var id string
		if err := pool.QueryRow(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts) values($1,gen_random_uuid(),$2,$3,$4,$5,2) returning job_id`, tenant, queue, kind, version, raw).Scan(&id); err != nil {
			t.Fatal(err)
		}
		excluded = append(excluded, id)
	}
	w, err := NewStatusWorker(r, "worker-main", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	result, err := w.ProcessOnce(ctx, p)
	if err != nil || !result.Claimed || !result.Completed || result.InsertedObservations != 1 || result.FailureRecorded {
		t.Fatalf("completion: %+v %v", result, err)
	}
	// A new instance after the commit must not execute the job again.
	w, err = NewStatusWorker(r, "worker-restarted", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	result, err = w.ProcessOnce(ctx, p)
	if err != nil || result.Claimed || result.Completed {
		t.Fatalf("replay: %+v %v", result, err)
	}
	var completed, untouched bool
	if err := pool.QueryRow(ctx, `select j.completed_at is not null and j.attempts=1 and j.terminal_error_code is null and e.state='processed' and e.processed_at is not null from platform.job j join integration.webhook_event e on e.tenant_id=j.tenant_id and e.provider_event_id=$2 where j.tenant_id=$1 and j.job_id=$3`, p.TenantID, event, original.JobID).Scan(&completed); err != nil || !completed {
		t.Fatal("job/inbox completion", completed, err)
	}
	if err := pool.QueryRow(ctx, `select count(*)=7 and bool_and(attempts=0 and completed_at is null and terminal_error_code is null) from platform.job where job_id=any($1::uuid[])`, excluded).Scan(&untouched); err != nil || !untouched {
		t.Fatal("claimed another scope", untouched, err)
	}
	var audits int
	var safe bool
	if err := pool.QueryRow(ctx, `select count(*),bool_and(reason_code='STATUS_OBSERVED' and evidence ? 'attempt' and evidence::text not like '%wamid%' and evidence::text not like '%549111%') from audit.event where tenant_id=$1 and action='whatsapp.status_job.completed'`, p.TenantID).Scan(&audits, &safe); err != nil || audits != 1 || !safe {
		t.Fatal("audit", audits, safe, err)
	}
}

func TestStatusWorkerTransientAndTerminalFailures(t *testing.T) {
	for _, terminal := range []bool{false, true} {
		t.Run(fmt.Sprint(terminal), func(t *testing.T) {
			batch := statusBatch("delivered", "1603086314", nil)
			if terminal {
				batch = statusBatch("delivered", "1603086314", map[string]any{"id": "wamid.unresolved"})
			}
			r, _, pool, p, _, event, path := routerFixture(t, batch)
			ctx := context.Background()
			if _, err := pool.Exec(ctx, `update platform.job set max_attempts=2 where tenant_id=$1`, p.TenantID); err != nil {
				t.Fatal(err)
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !terminal {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			}
			w, err := NewStatusWorker(r, "worker-retry", time.Minute, 0)
			if err != nil {
				t.Fatal(err)
			}
			first, err := w.ProcessOnce(ctx, p)
			if err == nil || !first.Claimed || first.Completed || !first.FailureRecorded || first.Terminal {
				t.Fatalf("retry: %+v %v", first, err)
			}
			if !terminal {
				if err := os.WriteFile(path, raw, 0600); err != nil {
					t.Fatal(err)
				}
			}
			second, err := w.ProcessOnce(ctx, p)
			if terminal {
				if err == nil || !second.Terminal || !second.FailureRecorded || second.Completed {
					t.Fatalf("terminal: %+v %v", second, err)
				}
			} else if err != nil || !second.Completed || second.InsertedObservations != 1 {
				t.Fatalf("recovery: %+v %v", second, err)
			}
			third, err := w.ProcessOnce(ctx, p)
			if err != nil || third.Claimed {
				t.Fatalf("unbounded retry: %+v %v", third, err)
			}
			var state string
			var audits int
			if err := pool.QueryRow(ctx, `select state from integration.webhook_event where tenant_id=$1 and provider_event_id=$2`, p.TenantID, event).Scan(&state); err != nil {
				t.Fatal(err)
			}
			if (terminal && state != "failed") || (!terminal && state != "processed") {
				t.Fatal("inbox state", state)
			}
			if err := pool.QueryRow(ctx, `select count(*) from audit.event where tenant_id=$1 and action like 'whatsapp.status_job.%'`, p.TenantID).Scan(&audits); err != nil || audits != 2 {
				t.Fatal("missing attempt audit", audits, err)
			}
		})
	}
}

func TestStatusWorkerAtomicCompletionRollback(t *testing.T) {
	r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	// Fault injection is limited to this synthetic tenant in a disposable DB.
	if !notificationEventID.MatchString(p.TenantID) {
		t.Fatal("unsafe fixture tenant")
	}
	sql := fmt.Sprintf(`create function integration.reject_v277_completion() returns trigger language plpgsql as $body$ begin if new.tenant_id='%s'::uuid and new.state='processed' then raise exception 'synthetic ACK failure'; end if; return new; end $body$; create trigger reject_v277_completion before update on integration.webhook_event for each row execute function integration.reject_v277_completion();`, p.TenantID)
	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		if _, err := pool.Exec(ctx, `drop trigger if exists reject_v277_completion on integration.webhook_event; drop function if exists integration.reject_v277_completion();`); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(cleanup)
	w, err := NewStatusWorker(r, "worker-atomic", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	first, err := w.ProcessOnce(ctx, p)
	if err == nil || first.Completed || first.InsertedObservations != 1 || !first.FailureRecorded {
		t.Fatalf("failed ACK: %+v %v", first, err)
	}
	var pending bool
	if err := pool.QueryRow(ctx, `select j.completed_at is null and e.state='received' and e.processed_at is null from platform.job j join integration.webhook_event e on e.tenant_id=j.tenant_id and e.provider_event_id=$2 where j.tenant_id=$1`, p.TenantID, event).Scan(&pending); err != nil || !pending {
		t.Fatal("partial job/inbox commit", pending, err)
	}
	cleanup()
	second, err := w.ProcessOnce(ctx, p)
	if err != nil || !second.Completed || second.InsertedObservations != 0 {
		t.Fatalf("ACK recovery: %+v %v", second, err)
	}
	var count int
	if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, p.TenantID).Scan(&count); err != nil || count != 1 {
		t.Fatal("duplicate observation", count, err)
	}
}

func TestStatusWorkerReclaimsAndQuarantinesCrashedAttempt(t *testing.T) {
	for _, exhausted := range []bool{false, true} {
		t.Run(fmt.Sprint(exhausted), func(t *testing.T) {
			r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
			ctx := context.Background()
			attempt := 1
			if exhausted {
				attempt = 10
			}
			if _, err := pool.Exec(ctx, `update platform.job set attempts=$2,claimed_by='crashed',claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1`, p.TenantID, attempt); err != nil {
				t.Fatal(err)
			}
			w, err := NewStatusWorker(r, "worker-recovery", time.Minute, 0)
			if err != nil {
				t.Fatal(err)
			}
			result, err := w.ProcessOnce(ctx, p)
			if exhausted {
				if err == nil || result.Completed || !result.Terminal || !result.FailureRecorded {
					t.Fatalf("orphan: %+v %v", result, err)
				}
			} else if err != nil || !result.Completed {
				t.Fatalf("reclaimed: %+v %v", result, err)
			}
			var state string
			var count int
			if err := pool.QueryRow(ctx, `select state from integration.webhook_event where tenant_id=$1 and provider_event_id=$2`, p.TenantID, event).Scan(&state); err != nil {
				t.Fatal(err)
			}
			if exhausted && state != "failed" {
				t.Fatal(state)
			}
			if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, p.TenantID).Scan(&count); err != nil {
				t.Fatal(err)
			}
			if exhausted && count != 0 {
				t.Fatal("quarantine executed effect")
			}
			again, err := w.ProcessOnce(ctx, p)
			if err != nil || again.Claimed || again.Terminal {
				t.Fatalf("repeat: %+v %v", again, err)
			}
		})
	}
}

func TestStatusWorkerConcurrentClaims(t *testing.T) {
	r, _, pool, p, _, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	other, err := NewStatusRouter(&r.observer, r.connection, r.retention, r.key, r.maxRoutes)
	if err != nil {
		t.Fatal(err)
	}
	w1, err := NewStatusWorker(r, "worker-a", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	w2, err := NewStatusWorker(other, "worker-b", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	counts := make(chan bool, 2)
	errs := make(chan error, 2)
	for _, w := range []*StatusWorker{w1, w2} {
		wg.Add(1)
		go func(w *StatusWorker) {
			defer wg.Done()
			<-start
			res, err := w.ProcessOnce(context.Background(), p)
			counts <- res.Completed
			errs <- err
		}(w)
	}
	close(start)
	wg.Wait()
	close(counts)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	completed := 0
	for done := range counts {
		if done {
			completed++
		}
	}
	if completed != 1 {
		t.Fatal("completion count", completed)
	}
	var attempt int
	if err := pool.QueryRow(context.Background(), `select attempts from platform.job where tenant_id=$1`, p.TenantID).Scan(&attempt); err != nil || attempt != 1 {
		t.Fatal("claim race", attempt, err)
	}
}

func TestStatusWorkerClaimLostAfterObservation(t *testing.T) {
	r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	w, err := NewStatusWorker(r, "worker-old", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	r.observer.Secrets = statusSecretHook(func(context.Context) (string, error) {
		calls++
		if calls == 2 {
			if _, err := pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1`, p.TenantID); err != nil {
				t.Fatal(err)
			}
			match, _ := json.Marshal(map[string]string{"connection_id": r.connection, "provider_code": "meta-whatsapp", "event_type": "whatsapp.raw_webhook.received.v1"})
			jobs, err := w.jobs.ClaimScoped(ctx, "provider-events", "worker-new", time.Minute, 1, postgres.JobScope{TenantID: p.TenantID, JobType: "provider.webhook.received", SchemaVersion: 1, PayloadMatch: match})
			if err != nil || len(jobs) != 1 || jobs[0].Attempts != 2 {
				t.Fatal("takeover", jobs, err)
			}
		}
		return "test-app-secret", nil
	})
	result, err := w.ProcessOnce(ctx, p)
	if err == nil || result.Completed || result.FailureRecorded || result.InsertedObservations != 1 {
		t.Fatalf("stale result: %+v %v", result, err)
	}
	var retained bool
	if err := pool.QueryRow(ctx, `select j.completed_at is null and j.claimed_by='worker-new' and j.attempts=2 and e.state='received' from platform.job j join integration.webhook_event e on e.tenant_id=j.tenant_id and e.provider_event_id=$2 where j.tenant_id=$1`, p.TenantID, event).Scan(&retained); err != nil || !retained {
		t.Fatal("old worker damaged new claim", retained, err)
	}
	// After another simulated crash, a fresh claim completes without new effects.
	r.observer.Secrets = appSecretFixture{}
	if _, err := pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	result, err = w.ProcessOnce(ctx, p)
	if err != nil || !result.Completed || result.InsertedObservations != 0 {
		t.Fatalf("recovery after takeover: %+v %v", result, err)
	}
}

func TestStatusWorkerRejectsBeforeClaim(t *testing.T) {
	for _, mode := range []string{"permission", "tenant", "disabled", "cancelled", "busy"} {
		t.Run(mode, func(t *testing.T) {
			r, _, pool, p, _, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
			tenant := p.TenantID
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			w, err := NewStatusWorker(r, "worker-guard", time.Minute, 0)
			if err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "permission":
				p.Permissions = map[string]struct{}{}
			case "tenant":
				p.TenantID = "00000000-0000-0000-0000-000000000001"
			case "disabled":
				if _, err := pool.Exec(ctx, `update integration.provider_connection set state='disabled' where tenant_id=$1`, tenant); err != nil {
					t.Fatal(err)
				}
			case "cancelled":
				cancel()
			case "busy":
				w.busy <- struct{}{}
				defer func() { <-w.busy }()
			}
			result, err := w.ProcessOnce(ctx, p)
			if err == nil || result.Claimed || result.Completed {
				t.Fatalf("unauthorized claim: %+v %v", result, err)
			}
			var attempts int
			if err := pool.QueryRow(context.Background(), `select attempts from platform.job where tenant_id=$1`, tenant).Scan(&attempts); err != nil || attempts != 0 {
				t.Fatal("claim changed", attempts, err)
			}
		})
	}
}

func TestStatusWorkerClaimPayloadDriftCannotComplete(t *testing.T) {
	r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	w, err := NewStatusWorker(r, "worker-drift", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	r.observer.Secrets = statusSecretHook(func(context.Context) (string, error) {
		calls++
		if calls == 2 {
			if _, err := pool.Exec(ctx, `update platform.job set payload=jsonb_set(payload,'{unexpected}','true') where tenant_id=$1`, p.TenantID); err != nil {
				t.Fatal(err)
			}
		}
		return "test-app-secret", nil
	})
	result, err := w.ProcessOnce(ctx, p)
	if err == nil || result.Completed || result.FailureRecorded || result.InsertedObservations != 1 {
		t.Fatalf("drift accepted: %+v %v", result, err)
	}
	var pending bool
	if err := pool.QueryRow(ctx, `select j.completed_at is null and j.payload ? 'unexpected' and e.state='received' from platform.job j join integration.webhook_event e on e.tenant_id=j.tenant_id and e.provider_event_id=$2 where j.tenant_id=$1`, p.TenantID, event).Scan(&pending); err != nil || !pending {
		t.Fatal("changed claim overwritten", pending, err)
	}
}
