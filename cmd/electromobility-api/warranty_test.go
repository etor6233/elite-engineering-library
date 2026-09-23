package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarrantyHostDisabledReadsNoProfile(t *testing.T) {
	for _, enabled := range []string{"", "false"} {
		m, err := selectedWarrantyModule(context.Background(), nil, func(key string) string {
			if key != "WARRANTY_ENABLED" {
				t.Fatal("disabled module read configuration", key)
			}
			return enabled
		})
		if err != nil || m != nil {
			t.Fatal(m, err)
		}
	}
	if _, err := selectedWarrantyModule(context.Background(), nil, func(string) string { return "TRUE" }); err == nil {
		t.Fatal("untyped activation admitted")
	}
}
func TestWarrantyHostProfileAndSchemaActivation(t *testing.T) {
	rawURL := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if rawURL == "" {
		t.Skip("owned reference PostgreSQL URL is required")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() != "127.0.0.1" || !strings.HasPrefix(parsed.Path, "/elite_payment_connected_") {
		t.Fatal("fixture database outside owned scope")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, rawURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	d := wc.ProfileDocument{Schema: "elite-warranty-profile/v1", Scope: "MATERIALIZED_PROFILE", Algorithm: "bc-inclusive-fixed-terms", AlgorithmRevision: 1, TenantID: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", OrganizationID: "store", FactoryOrganizationID: "factory", PolicyID: "synthetic", TermsVersion: "fixture-v1", TermsText: "Fixture profile only; no legal terms inferred.", BusinessTimeZone: "UTC", PartsDurationDays: 365, LaborDurationDays: 365, WorkReservationSeconds: 3600, FaultExclusions: []string{}, Settlement: "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT", AuthorityReference: "fixture-authority", DecisionReference: "fixture-decision"}
	raw, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "profile.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	values := map[string]string{"WARRANTY_ENABLED": "true", "WARRANTY_PROFILE_FILE": path, "WARRANTY_PROFILE_SHA256": wc.SHA(raw), "WARRANTY_TENANT_ID": d.TenantID, "WARRANTY_ORGANIZATION_ID": d.OrganizationID}
	lookup := func(key string) string { return values[key] }
	module, err := selectedWarrantyModule(ctx, pool, lookup)
	if err != nil || module == nil {
		t.Fatal("selected module", err)
	}
	mux := http.NewServeMux()
	module.Register(mux, nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest("GET", "/v1/franchise/warranty/profile?organization_id=store", nil))
	if response.Code != 401 {
		t.Fatal("route not protected/mounted", response.Code)
	}
	for _, v := range []struct{ key, value string }{{"WARRANTY_PROFILE_SHA256", strings.Repeat("f", 64)}, {"WARRANTY_PROFILE_FILE", "relative.json"}, {"WARRANTY_TENANT_ID", "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"}, {"WARRANTY_ORGANIZATION_ID", "other"}} {
		previous := values[v.key]
		values[v.key] = v.value
		_, err := selectedWarrantyModule(ctx, pool, lookup)
		values[v.key] = previous
		if err == nil {
			t.Fatal("mismatched configuration admitted", v.key)
		}
	}
	if err = os.WriteFile(path, append(raw, []byte(strings.Repeat(" ", 32769))...), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = selectedWarrantyModule(ctx, pool, lookup); err == nil {
		t.Fatal("oversized profile admitted")
	}
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	// This is the owned disposable fixture database; restore the guard even if
	// the negative activation assertion fails.
	if _, err = pool.Exec(ctx, `alter table service_ops.service_case disable trigger warranty_claim_case_guard`); err != nil {
		t.Fatal(err)
	}
	func() {
		defer func() {
			if _, err := pool.Exec(ctx, `alter table service_ops.service_case enable trigger warranty_claim_case_guard`); err != nil {
				t.Fatal("restore fixture guard", err)
			}
		}()
		if _, err := selectedWarrantyModule(ctx, pool, lookup); err == nil {
			t.Fatal("disabled command guard admitted")
		}
	}()
	t.Log("WARRANTY_HOST_PG_PASS exact_profile_scope_size_schema_guard=true routes_protected=true no_live_credentials=true")
}
