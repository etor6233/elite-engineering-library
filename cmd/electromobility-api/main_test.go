package main

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/providerintegration"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

type rejectIdentity struct{}

func TestPaymentRequestSelection(t *testing.T) {
	for _, v := range []string{"", "stripe", "mercadopago", "unknown", "true", "STRIPE", " stripe"} {
		got, err := paymentRequestProvider(func(string) string { return v })
		valid := v == "" || v == "stripe" || v == "mercadopago"
		if (err == nil) != valid || (valid && got != v) {
			t.Fatalf("selection %q result %q err%v", v, got, err)
		}
	}
}

func (rejectIdentity) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{}, identity.ErrUnauthenticated
}

func TestOptionalProviderConfigurationFailsClosedWithoutLeaking(t *testing.T) {
	entry := `{"tenant_id":"tenant","connection_id":"primary","provider_code":"sandbox","secret_env":"TEST_SECRET"}`
	for _, tc := range []struct {
		name, config, secret string
		valid, present       bool
	}{
		{"absent", "", "", true, false},
		{"empty-registry", "[]", "", true, false},
		{"malformed", "{", "", false, false},
		{"inline-secret", `[{"secret":"DO_NOT_LEAK_VALUE"}]`, "", false, false},
		{"missing-reference", `[{}]`, "", false, false},
		{"missing-secret", "[" + entry + "]", "", false, false},
		{"short-secret", "[" + entry + "]", "DO_NOT_LEAK_VALUE", false, false},
		{"valid", "[" + entry + "]", strings.Repeat("s", 32), true, true},
		{"duplicate", "[" + entry + "," + entry + "]", strings.Repeat("s", 32), false, false},
		{"trailing", "[" + entry + "]{}", strings.Repeat("s", 32), false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			registry, err := providerintegration.LoadHMACConnections([]byte(tc.config), func(string) (string, bool) { calls++; return tc.secret, tc.secret != "" })
			if (err == nil) != tc.valid {
				t.Fatal("unexpected configuration outcome")
			}
			if err != nil && (strings.Contains(err.Error(), "DO_NOT_LEAK_VALUE") || (tc.secret != "" && strings.Contains(err.Error(), tc.secret))) {
				t.Fatal("secret exposed in error")
			}
			if err == nil {
				_, found := registry.Lookup("sandbox", "primary")
				if found != tc.present {
					t.Fatal("unexpected provider activated")
				}
			}
			if (tc.config == "" || tc.config == "[]") && calls != 0 {
				t.Fatal("inactive integration requested secrets")
			}
		})
	}
}

func TestFiscalActivationRequiresExplicitValidConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name, flag, socket    string
		wantModule, wantError bool
	}{
		{"absent", "", "", false, false},
		{"deferred", "false", "", false, false},
		{"deferred-ignores-unused-socket", "false", "must-not-connect", false, false},
		{"invalid", "yes", "", false, true},
		{"whitespace", " true ", "", false, true},
		{"enabled-missing-socket", "true", "", false, true},
		{"enabled-relative-socket", "true", "relative.sock", false, true},
		{"enabled-valid-socket", "true", filepath.Join(t.TempDir(), "wsfe.sock"), true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			socketLookups := 0
			m, err := selectedFiscalModule(nil, randomid.Generator{}, func(k string) string {
				if k == "ARCA_ENABLED" {
					return tc.flag
				}
				socketLookups++
				return tc.socket
			})
			if (err != nil) != tc.wantError || (m != nil) != tc.wantModule {
				t.Fatalf("unexpected module/error presence")
			}
			if !tc.wantModule && !tc.wantError && socketLookups != 0 {
				t.Fatal("deferred path read fiscal configuration")
			}
			mux := http.NewServeMux()
			mux.HandleFunc("GET /independent", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
			if m != nil {
				m.Register(mux, rejectIdentity{})
			}
			independent := httptest.NewRecorder()
			mux.ServeHTTP(independent, httptest.NewRequest("GET", "/independent", nil))
			if independent.Code != 204 {
				t.Fatal("independent route affected")
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, httptest.NewRequest("POST", "/v1/fiscal/invoices", nil))
			wantStatus := http.StatusNotFound
			if tc.wantModule {
				wantStatus = http.StatusUnauthorized
			}
			if w.Code != wantStatus {
				t.Fatalf("fiscal route status = %d, want %d", w.Code, wantStatus)
			}
		})
	}
}
