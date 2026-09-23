package main

import (
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
	"testing"
)

func TestCatalogSourceHostMigration(t *testing.T) {
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("owned fixture required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	profile := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893", OrganizationID: "j3-store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	body, _ := json.Marshal(profile)
	sum := sha256.Sum256(body)
	file := filepath.Join(t.TempDir(), "profile.json")
	if e = os.WriteFile(file, body, 0600); e != nil {
		t.Fatal(e)
	}
	config := map[string]string{"CATALOG_RELEASE_ENABLED": "true", "CATALOG_RELEASE_PROFILE_FILE": file, "CATALOG_RELEASE_PROFILE_SHA256": hex.EncodeToString(sum[:])}
	lookup := func(k string) string { return config[k] }
	price := postgres.NewCommerce(pool)
	if module, e := selectedCatalogReleaseModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("current source host rejected", e)
	}
	if _, e = pool.Exec(ctx, `alter table catalog.release_command rename constraint release_command_kind_check to source_kind_saved_fixture`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table catalog.release_command rename constraint source_kind_saved_fixture to release_command_kind_check`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing source migration admitted")
	}
	t.Log("CATALOG_SOURCE_HOST_PASS exact source migration required; earlier host guards retained; no provider call")
}
