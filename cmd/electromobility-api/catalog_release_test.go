package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCatalogReleaseHost(t *testing.T) {
	ctx := context.Background()
	calls := 0
	off := func(k string) string {
		calls++
		if k != "CATALOG_RELEASE_ENABLED" {
			t.Fatal("disabled module read", k)
		}
		return "false"
	}
	if module, e := selectedCatalogReleaseModule(ctx, nil, nil, off); e != nil || module != nil || calls != 1 {
		t.Fatal(module, e, calls)
	}
	rawURL := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if rawURL == "" {
		t.Skip("owned fixture DB not configured")
	}
	pool, e := pgxpool.New(ctx, rawURL)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	profile := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893", OrganizationID: "j3-store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	raw, _ := json.Marshal(profile)
	sum := sha256.Sum256(raw)
	hash := hex.EncodeToString(sum[:])
	file := filepath.Join(t.TempDir(), "profile.json")
	if e = os.WriteFile(file, raw, 0600); e != nil {
		t.Fatal(e)
	}
	config := map[string]string{"CATALOG_RELEASE_ENABLED": "true", "CATALOG_RELEASE_PROFILE_FILE": file, "CATALOG_RELEASE_PROFILE_SHA256": hash}
	lookup := func(k string) string { return config[k] }
	price := postgres.NewCommerce(pool)
	if _, e = selectedCatalogReleaseModule(ctx, pool, nil, lookup); e == nil {
		t.Fatal("missing selected price owner admitted")
	}
	config["CATALOG_RELEASE_PROFILE_SHA256"] = ""
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("unbound profile admitted")
	}
	config["CATALOG_RELEASE_PROFILE_SHA256"] = hash
	if module, e := selectedCatalogReleaseModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("valid activation", e)
	}
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("activation made a provider request")
		w.WriteHeader(500)
	}))
	defer receiver.Close()
	feed := cr.FeedProfile{Schema: "elite-catalog-feed/v1", ReceiverID: "fixture-receiver", BaseURL: receiver.URL, Mode: "fixture"}
	feedRaw, _ := json.Marshal(feed)
	feedSum := sha256.Sum256(feedRaw)
	feedFile := filepath.Join(t.TempDir(), "feed.json")
	if e = os.WriteFile(feedFile, feedRaw, 0600); e != nil {
		t.Fatal(e)
	}
	config["CATALOG_FEED_ENABLED"] = "true"
	config["CATALOG_FEED_PROFILE_FILE"] = feedFile
	config["CATALOG_FEED_PROFILE_SHA256"] = hex.EncodeToString(feedSum[:])
	config["CATALOG_FEED_TOKEN"] = "synthetic-fixture-token"
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing feed fence key admitted")
	}
	config["CATALOG_FEED_HMAC_KEY"] = base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x42}, 32))
	if module, e := selectedCatalogReleaseModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("feed activation", e)
	} else if module.(httpapi.CatalogReleaseModule).Feed == nil {
		t.Fatal("feed absent")
	}
	delete(config, "CATALOG_FEED_TOKEN")
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing receiver token admitted")
	}
	config["CATALOG_FEED_TOKEN"] = "synthetic-fixture-token"
	config["CATALOG_FEED_PROFILE_SHA256"] = hash
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("feed hash mismatch admitted")
	}
	config["CATALOG_FEED_PROFILE_SHA256"] = hex.EncodeToString(feedSum[:])
	if _, e = pool.Exec(ctx, `alter table catalog.release_feed_intent disable trigger catalog_release_feed_immutable`); e != nil {
		t.Fatal(e)
	}
	func() {
		defer func() {
			if _, e := pool.Exec(context.Background(), `alter table catalog.release_feed_intent enable trigger catalog_release_feed_immutable`); e != nil {
				t.Error(e)
			}
		}()
		if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
			t.Fatal("disabled feed guard admitted")
		}
	}()
	if _, e = pool.Exec(ctx, `alter table approval.request disable trigger catalog_release_approval_guard`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table approval.request enable trigger catalog_release_approval_guard`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("disabled approval guard admitted")
	}
	t.Log("CATALOG_RELEASE_HOST_PASS disabled no-read; selected price owner, exact profile hashes,12base guards, feed immutable guard/key/token required; zero provider requests at startup")
}
