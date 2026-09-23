package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/providerintegration"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestProviderWebhookConnectionAndReplayScope(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	if !strings.HasPrefix(url, "postgres://postgres@127.0.0.1:55959/elite_provider_") {
		t.Fatal("requires dedicated loopback elite_provider_ test database")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28d41"
	defer func() {
		for _, q := range []string{
			`delete from platform.job where tenant_id=$1`,
			`delete from integration.webhook_event where tenant_id=$1`,
			`delete from integration.provider_connection where tenant_id=$1`,
			`delete from platform.tenant where tenant_id=$1`,
		} {
			if _, err := pool.Exec(ctx, q, tenant); err != nil {
				t.Error("fixture cleanup", err)
			}
		}
	}()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'provider-scope','Synthetic Scope','Synthetic Scope')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'primary','sandbox','SANDBOX_SECRET')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewProviderIntegration(pool)
	base := providerintegration.Receipt{TenantID: tenant, ConnectionID: "primary", ProviderCode: "sandbox", ProviderEventID: "evt-scope", EventType: "payment.updated", BodyHash: strings.Repeat("a", 64), Payload: json.RawMessage(`{"state":"paid"}`)}
	if replay, err := repo.AcceptWebhook(ctx, base); replay || err != nil {
		t.Fatal(replay, err)
	}
	for _, mode := range []string{"disabled_new", "disabled_replay", "event_type", "payload", "provider", "connection"} {
		t.Run(mode, func(t *testing.T) {
			if _, err := pool.Exec(ctx, `update integration.provider_connection set state='active' where tenant_id=$1`, tenant); err != nil {
				t.Fatal(err)
			}
			value := base
			want := providerintegration.ErrConflict
			switch mode {
			case "disabled_new", "disabled_replay":
				if _, err := pool.Exec(ctx, `update integration.provider_connection set state='disabled' where tenant_id=$1`, tenant); err != nil {
					t.Fatal(err)
				}
				if mode == "disabled_new" {
					value.ProviderEventID = "evt-new"
				}
				want = providerintegration.ErrConnection
			case "event_type":
				value.EventType = "refund.created"
			case "payload":
				value.Payload = json.RawMessage(`{"state":"refunded"}`)
			case "provider":
				value.ProviderCode = "other"
				want = providerintegration.ErrConnection
			case "connection":
				value.ConnectionID = "missing"
				want = providerintegration.ErrConnection
			}
			if replay, err := repo.AcceptWebhook(ctx, value); replay || !errors.Is(err, want) {
				t.Errorf("scope accepted or wrong failure: replay=%v error=%v want=%v", replay, err, want)
			}
		})
	}
	if _, err := pool.Exec(ctx, `update integration.provider_connection set state='active' where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	base.Payload = json.RawMessage(`{ "state" : "paid" }`)
	if replay, err := repo.AcceptWebhook(ctx, base); !replay || err != nil {
		t.Fatal("equivalent JSON replay rejected", replay, err)
	}
	t.Run("concurrent_disable", func(t *testing.T) {
		bounded, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		change, err := pool.Begin(bounded)
		if err != nil {
			t.Fatal(err)
		}
		defer change.Rollback(context.Background())
		if _, err := change.Exec(bounded, `update integration.provider_connection set state='disabled' where tenant_id=$1`, tenant); err != nil {
			t.Fatal(err)
		}
		value := base
		value.ProviderEventID = "evt-concurrent-disable"
		type outcome struct {
			replay bool
			err    error
		}
		done := make(chan outcome, 1)
		go func() { replay, err := repo.AcceptWebhook(bounded, value); done <- outcome{replay, err} }()
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		probe, stop := context.WithTimeout(ctx, 5*time.Second)
		defer stop()
		for {
			var blocked bool
			if err := pool.QueryRow(probe, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like 'select connection_id from integration.provider_connection%')`).Scan(&blocked); err != nil {
				t.Fatal("lock probe", err)
			}
			if blocked {
				break
			}
			select {
			case result := <-done:
				t.Fatal("receipt did not wait for current connection", result)
			case <-probe.Done():
				t.Fatal("connection lock not observed")
			case <-ticker.C:
			}
		}
		if err := change.Commit(bounded); err != nil {
			t.Fatal(err)
		}
		select {
		case result := <-done:
			if result.replay || !errors.Is(result.err, providerintegration.ErrConnection) {
				t.Fatal("concurrent disable bypassed", result)
			}
		case <-bounded.Done():
			t.Fatal("receipt did not resume after disable")
		}
	})
	var events, jobs int
	if err := pool.QueryRow(ctx, `select (select count(*) from integration.webhook_event where tenant_id=$1),(select count(*) from platform.job where tenant_id=$1)`, tenant).Scan(&events, &jobs); err != nil || events != 1 || jobs != 1 {
		t.Fatal("scope writes escaped", events, jobs, err)
	}
}

func TestProviderWebhookDeduplicationAndJobAreAtomic(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28d40"
	defer func() {
		_, _ = pool.Exec(ctx, `delete from platform.job where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from integration.webhook_event where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from integration.reconciliation_item where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from integration.external_mapping where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from integration.provider_connection where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'provider-test','Provider Test','Provider')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'primary','sandbox','SANDBOX_SECRET')`, tenant); err != nil {
		t.Fatal(err)
	}
	repository := NewProviderIntegration(pool)
	receipt := providerintegration.Receipt{TenantID: tenant, ConnectionID: "primary", ProviderCode: "sandbox", ProviderEventID: "evt-1", EventType: "payment.updated", BodyHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Payload: json.RawMessage(`{"state":"paid"}`)}
	if replayed, err := repository.AcceptWebhook(ctx, receipt); err != nil || replayed {
		t.Fatalf("first replay=%v err=%v", replayed, err)
	}
	if replayed, err := repository.AcceptWebhook(ctx, receipt); err != nil || !replayed {
		t.Fatalf("second replay=%v err=%v", replayed, err)
	}
	receipt.BodyHash = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	if _, err := repository.AcceptWebhook(ctx, receipt); !errors.Is(err, providerintegration.ErrConflict) {
		t.Fatalf("conflict err=%v", err)
	}
	var events, jobs int
	if err = pool.QueryRow(ctx, `select count(*) from integration.webhook_event where tenant_id=$1`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.job where tenant_id=$1 and job_type='provider.webhook.received'`, tenant).Scan(&jobs); err != nil {
		t.Fatal(err)
	}
	if events != 1 || jobs != 1 {
		t.Fatalf("events=%d jobs=%d", events, jobs)
	}
	reconciliation := providerintegration.NewReconciliationService(repository)
	if err = reconciliation.UpsertMapping(ctx, providerintegration.Mapping{TenantID: tenant, ProviderCode: "sandbox", ResourceType: "offer", InternalID: "variant-1", ExternalID: "external-1", VersionToken: "v1"}); err != nil {
		t.Fatal(err)
	}
	if err = reconciliation.Record(ctx, providerintegration.Reconciliation{TenantID: tenant, ID: "recon-1", ProviderCode: "sandbox", ResourceType: "offer", InternalID: "variant-1", ExternalID: "external-1", State: "state-mismatch", Evidence: json.RawMessage(`{"internal":"active","external":"paused"}`), ObservedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if err = reconciliation.Resolve(ctx, tenant, "sandbox", "recon-1"); err != nil {
		t.Fatal(err)
	}
	var mapping, resolved int
	if err = pool.QueryRow(ctx, `select count(*) from integration.external_mapping where tenant_id=$1 and internal_id='variant-1' and external_id='external-1'`, tenant).Scan(&mapping); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from integration.reconciliation_item where tenant_id=$1 and reconciliation_id='recon-1' and state='resolved' and resolved_at is not null`, tenant).Scan(&resolved); err != nil {
		t.Fatal(err)
	}
	if mapping != 1 || resolved != 1 {
		t.Fatalf("mapping=%d resolved=%d", mapping, resolved)
	}
}
