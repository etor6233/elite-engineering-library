package main

// AUTHORED activation contract. Pool is lazy and no database is contacted here;
// the connected PG/browser suite separately verifies every mounted operation.
import (
	"context"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestStoredValueOptionalActivation(t *testing.T) {
	ctx := context.Background()
	root, e := filepath.Abs("../..")
	if e != nil {
		t.Fatal(e)
	}
	read := func(rel string) []byte {
		t.Helper()
		b, e := os.ReadFile(filepath.Join(root, rel))
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	config := map[string]string{"STORED_VALUE_ENABLED": "true", "STORED_VALUE_PROFILE_FILE": filepath.Join(root, "deploy/stored-value/profile.reference.json"), "STORED_VALUE_PROFILE_SHA256": sv.Hash(read("deploy/stored-value/profile.reference.json")), "STORED_VALUE_TENANT_ID": "11111111-1111-4111-8111-111111111111", "STORED_VALUE_ORGANIZATION_ID": "franchise-1", "STORED_VALUE_PYTHON": os.Getenv("HANDOVER_PROFILE_PYTHON"), "STORED_VALUE_SCRIPT": filepath.Join(root, "odoo_loyalty/run.py"), "STORED_VALUE_SCRIPT_SHA256": sv.Hash(read("odoo_loyalty/run.py")), "STORED_VALUE_MANIFEST": filepath.Join(root, "odoo_loyalty/engine-lock.json"), "STORED_VALUE_MANIFEST_SHA256": sv.Hash(read("odoo_loyalty/engine-lock.json"))}
	if config["STORED_VALUE_PYTHON"] == "" {
		t.Fatal("explicit fixed Python required for source preflight")
	}
	pool, e := pgxpool.New(ctx, "postgres://fixture@127.0.0.1:1/fixture?sslmode=disable&connect_timeout=1")
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	for _, which := range []string{"valid", "disabled", "invalid-flag", "profile-sha", "organization", "script-sha", "manifest-sha", "relative-profile"} {
		t.Run(which, func(t *testing.T) {
			v := map[string]string{}
			for k, val := range config {
				v[k] = val
			}
			switch which {
			case "disabled":
				v["STORED_VALUE_ENABLED"] = "false"
			case "invalid-flag":
				v["STORED_VALUE_ENABLED"] = "1"
			case "profile-sha":
				v["STORED_VALUE_PROFILE_SHA256"] = "bad"
			case "organization":
				v["STORED_VALUE_ORGANIZATION_ID"] = "other"
			case "script-sha":
				v["STORED_VALUE_SCRIPT_SHA256"] = "bad"
			case "manifest-sha":
				v["STORED_VALUE_MANIFEST_SHA256"] = "bad"
			case "relative-profile":
				v["STORED_VALUE_PROFILE_FILE"] = "profile.json"
			}
			module, e := storedValueModuleFactory(ctx, pool, func(k string) string { return v[k] })
			if which == "disabled" {
				if module != nil || e != nil {
					t.Fatal("disabled module mounted", e)
				}
				return
			}
			if which != "valid" {
				if e == nil || module != nil {
					t.Fatal("invalid activation admitted")
				}
				return
			}
			if e != nil || module == nil {
				t.Fatal("valid activation", e)
			}
			mux := http.NewServeMux()
			module.Register(mux, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, httptest.NewRequest("GET", "/v1/franchise/stored-value/profile?organization_id=franchise-1", nil))
			if w.Code != 401 {
				t.Fatal("route was not mounted behind authentication", w.Code)
			}
		})
	}
}
