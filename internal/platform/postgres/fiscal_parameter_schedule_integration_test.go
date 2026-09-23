package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestFiscalParameterScheduleLeaseSnapshotAndExplicitDecision(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2f133"
	cleanup := []string{
		`delete from platform.outbox_event where tenant_id=$1`, `delete from fiscal.parameter_refresh_schedule where tenant_id=$1`, `delete from fiscal.parameter_decision where tenant_id=$1`, `delete from fiscal.parameter_item where tenant_id=$1`, `delete from fiscal.parameter_snapshot where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`,
	}
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
		for _, query := range cleanup {
			if _, e = tx.Exec(ctx, query, tenant); e != nil {
				t.Errorf("cleanup: %v", e)
				return
			}
		}
		if e = tx.Commit(ctx); e != nil {
			t.Errorf("cleanup commit: %v", e)
		}
	}()
	for _, query := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'fiscal-v133-integration','Fiscal V133','Fiscal V133')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise-v133','Franchise','franchisee')`} {
		if _, err = pool.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repository := NewFiscal(pool)
	now := time.Now().UTC()
	snapshot := fiscal.ParameterSnapshot{TenantID: tenant, ID: "snapshot-v133", OrganizationID: "franchise", TaxpayerCUIT: "33693450239", Kind: "vat_rate", Items: []fiscal.ParameterItem{{Code: "5"}}, ResponseHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", FetchedAt: now}
	if _, _, err = repository.StoreParameterSnapshot(ctx, snapshot, "018f4d4a-7b36-7a21-8d10-2f4c54c2f134"); err != nil {
		t.Fatal(err)
	}
	schedule := fiscal.ParameterSchedule{TenantID: tenant, ID: "schedule-v133", OrganizationID: "franchise", TaxpayerCUIT: "33693450239", Kind: "vat_rate", IntervalSeconds: 900, NextRunAt: now.Add(-time.Minute), Active: true, Version: 1}
	schedule, err = repository.ConfigureParameterSchedule(ctx, schedule, "tax-owner", "018f4d4a-7b36-7a21-8d10-2f4c54c2f135")
	if err != nil {
		t.Fatal(err)
	}
	claimed, err := repository.ClaimParameterSchedule(ctx, "worker-a", time.Minute)
	if err != nil || claimed.ID != schedule.ID {
		t.Fatalf("claim=%+v err=%v", claimed, err)
	}
	if _, err = repository.ClaimParameterSchedule(ctx, "worker-b", time.Minute); !errors.Is(err, fiscal.ErrNoParameterWork) {
		t.Fatalf("second claim=%v", err)
	}
	if err = repository.CompleteParameterSchedule(ctx, claimed, "worker-a", snapshot.ID, "018f4d4a-7b36-7a21-8d10-2f4c54c2f136"); err != nil {
		t.Fatal(err)
	}
	if err = repository.DecideParameterSnapshot(ctx, tenant, "franchise", snapshot.ID, "tax-owner", true, "validated against current ARCA configuration", "018f4d4a-7b36-7a21-8d10-2f4c54c2f137"); err != nil {
		t.Fatal(err)
	}
	loaded, err := repository.GetParameterSnapshot(ctx, tenant, "franchise", snapshot.ID)
	if err != nil || len(loaded.Items) != 1 || loaded.Items[0].Code != "5" {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
}
