package main

import (
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/merchantbridge"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMerchantHostGuards(t *testing.T) {
	database := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if database == "" {
		t.Skip("owned fixture required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, database)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	cp := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: "d84c64f9-1252-4db5-bf5a-3af8026190f1", OrganizationID: "marketplace-store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}

	var variant string
	if e = pool.QueryRow(ctx, `select payload->>'variant_id' from approval.request where tenant_id=$1 and payload->>'account_id'='123456' order by created_at desc limit 1`, cp.TenantID).Scan(&variant); e != nil {
		t.Fatal(e)
	}
	p := mb.Profile{Schema: "elite-connected-merchant/v1", TenantID: cp.TenantID, OrganizationID: cp.OrganizationID, AccountID: "123456", DataSourceID: "789012", Language: "es", FeedLabel: "AR", Currency: "ARS", CurrencyDigits: 2, Origin: cp.Origin, RefreshCadenceDays: 30, Bindings: []mb.Binding{{VariantID: variant, OfferID: "fixture-bicycle-standard", Condition: "NEW", GTINs: []string{}}}}
	fileHash := func(path string) string {
		t.Helper()
		raw, e := os.ReadFile(path)
		if e != nil {
			t.Fatal(e)
		}
		return mb.Hash(raw)
	}
	script := filepath.Join(os.Getenv("MERCHANT_TEST_ROOT"), "google_merchant_product_sync", "connected_worker.py")
	python := os.Getenv("MERCHANT_TEST_PYTHON")
	runtime := merchantRuntime{Schema: "elite-merchant-runtime/v1", Python: python, PythonSHA256: fileHash(python), Script: script, ScriptSHA256: fileHash(script), OwnerSHA256: fileHash(filepath.Join(filepath.Dir(script), "sync_product.py")), RuntimeLockSHA256: fileHash(filepath.Join(filepath.Dir(script), "connected-runtime.lock.json")), Mode: "LOCAL_FIXTURES", FixtureOrigin: "http://127.0.0.1:1"}
	config := map[string]string{"MERCHANT_ENABLED": "true", "CATALOG_RELEASE_ENABLED": "true", "MERCHANT_HMAC_KEY": base64.StdEncoding.EncodeToString([]byte(strings.Repeat("f", 32)))}
	raw, _ := json.Marshal(runtime)
	f := filepath.Join(t.TempDir(), "runtime.json")
	if e = os.WriteFile(f, raw, 0600); e != nil {
		t.Fatal(e)
	}
	config["MERCHANT_RUNTIME_FILE"] = f
	config["MERCHANT_RUNTIME_SHA256"] = mb.Hash(raw)
	for _, v := range []struct {
		prefix string
		value  any
	}{{"MERCHANT", p}, {"CATALOG_RELEASE", cp}} {
		raw, _ := json.Marshal(v.value)
		sum := sha256.Sum256(raw)
		file := filepath.Join(t.TempDir(), v.prefix+".json")
		if e = os.WriteFile(file, raw, 0600); e != nil {
			t.Fatal(e)
		}
		config[v.prefix+"_PROFILE_FILE"] = file
		config[v.prefix+"_PROFILE_SHA256"] = hex.EncodeToString(sum[:])
	}
	lookup := func(k string) string { return config[k] }
	price := postgres.NewCommerce(pool)
	if module, e := selectedMerchantModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("current merchant host", e)
	}

	catalog, e := postgres.NewCatalogRelease(pool, cp, price)
	if e != nil {
		t.Fatal(e)
	}
	process := mb.Process{Python: runtime.Python, PythonSHA256: runtime.PythonSHA256, Script: runtime.Script, ScriptSHA256: runtime.ScriptSHA256, OwnerSHA256: runtime.OwnerSHA256, RuntimeLockSHA256: runtime.RuntimeLockSHA256, Mode: runtime.Mode, FixtureOrigin: runtime.FixtureOrigin}
	service, e := postgres.NewMerchantPublication(catalog, p, []byte(strings.Repeat("f", 32)), process)
	if e != nil {
		t.Fatal(e)
	}
	principal := identity.Principal{TenantID: p.TenantID, Subject: "queue-reader", Permissions: map[string]struct{}{"merchant:read": {}}, Organizations: map[string]struct{}{p.OrganizationID: {}}}
	queue, e := service.Queue(ctx, principal)
	if e != nil || len(queue.Items) != 1 || queue.Items[0].Action != "PROVIDER_REJECTED" || queue.Items[0].RefreshDueAt == nil || queue.Items[0].Generation != 1 || queue.Items[0].DesiredProductSHA256 == "" {
		t.Fatalf("queue: %+v %v", queue, e)
	}
	principal.Organizations = map[string]struct{}{"foreign": {}}
	if _, e = service.Queue(ctx, principal); e == nil {
		t.Fatal("foreign queue admitted")
	}
	// The queue and disabled/unconfigured host must never contact the fixture
	// origin (port1 deliberately has no provider) or mutate the proved journey.
	for _, key := range []string{"MERCHANT_PROFILE_SHA256", "CATALOG_RELEASE_ENABLED", "MERCHANT_RUNTIME_SHA256", "MERCHANT_HMAC_KEY"} {
		old := config[key]
		config[key] = ""
		if _, e := selectedMerchantModule(ctx, pool, price, lookup); e == nil {
			t.Fatal("missing binding admitted", key)
		}
		config[key] = old
	}
	if _, e = pool.Exec(ctx, `alter table catalog.merchant_observation disable trigger merchant_observation_immutable`); e != nil {
		t.Fatal(e)
	}
	_, rejected := selectedMerchantModule(ctx, pool, price, lookup)
	if _, e = pool.Exec(ctx, `alter table catalog.merchant_observation enable trigger merchant_observation_immutable`); e != nil {
		t.Fatal(e)
	}
	if rejected == nil {
		t.Fatal("missing observation guard admitted")
	}
	if _, e = pool.Exec(ctx, `alter table catalog.merchant_effect disable trigger merchant_effect_immutable`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table catalog.merchant_effect enable trigger merchant_effect_immutable`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedMerchantModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing immutable effect guard admitted")
	}
	t.Log("MERCHANT_HOST_PASS exact profile/hash/catalog activation/runtime/HMAC/migration required; queue exact and scoped; no provider network call")
}
