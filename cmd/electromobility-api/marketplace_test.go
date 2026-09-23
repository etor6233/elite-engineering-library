package main

import (
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/marketplacebridge"
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

func TestMarketplaceHostGuards(t *testing.T) {
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
	p := mb.Profile{Schema: "elite-marketplace-publication/v1", TenantID: cp.TenantID, OrganizationID: cp.OrganizationID, SellerID: "123456789", SiteID: "MLA", Currency: "ARS", CurrencyDigits: 2, MediaOrigin: cp.Origin, Bindings: []mb.Binding{{VariantID: "fixture-variant", SKU: "fixture-bicycle", CategoryID: "MLA3530", ListingTypeID: "gold_special", StockMode: "item", Attributes: []mb.Attribute{{ID: "ITEM_CONDITION", ValueID: "2230284"}}}}}
	config := map[string]string{"MARKETPLACE_ENABLED": "true", "CATALOG_RELEASE_ENABLED": "true", "MARKETPLACE_HMAC_KEY": base64.StdEncoding.EncodeToString([]byte(strings.Repeat("f", 32))), "MERCADOLIBRE_ACCESS_TOKEN": "synthetic-fixture-token-no-live"}
	for _, v := range []struct {
		prefix string
		value  any
	}{{"MARKETPLACE", p}, {"CATALOG_RELEASE", cp}} {
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
	if module, e := selectedMarketplaceModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("current host", e)
	}
	for _, key := range []string{"MARKETPLACE_PROFILE_SHA256", "CATALOG_RELEASE_ENABLED", "MERCADOLIBRE_ACCESS_TOKEN", "MARKETPLACE_HMAC_KEY"} {
		old := config[key]
		config[key] = ""
		if _, e := selectedMarketplaceModule(ctx, pool, price, lookup); e == nil {
			t.Fatal("missing binding admitted", key)
		}
		config[key] = old
	}
	if _, e = pool.Exec(ctx, `alter table catalog.marketplace_observation disable trigger marketplace_observation_immutable`); e != nil {
		t.Fatal(e)
	}
	_, rejected := selectedMarketplaceModule(ctx, pool, price, lookup)
	if _, e = pool.Exec(ctx, `alter table catalog.marketplace_observation enable trigger marketplace_observation_immutable`); e != nil {
		t.Fatal(e)
	}
	if rejected == nil {
		t.Fatal("missing observation guard admitted")
	}
	if _, e = pool.Exec(ctx, `alter table catalog.marketplace_effect disable trigger marketplace_effect_immutable`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table catalog.marketplace_effect enable trigger marketplace_effect_immutable`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedMarketplaceModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing immutable effect guard admitted")
	}
	t.Log("MARKETPLACE_HOST_PASS exact profile/hash/catalog activation/token/HMAC/migration required; no provider network call")
}
