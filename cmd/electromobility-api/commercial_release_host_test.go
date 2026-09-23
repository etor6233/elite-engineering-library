package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCommercialReleaseHostRequiresSelectedEffect(t *testing.T) {
	for _, commercial := range []bool{false, true} {
		runtime, env := handoverHostFixture(t)
		if commercial {
			raw, err := os.ReadFile(env["HANDOVER_PROFILE_FILE"])
			if err != nil {
				t.Fatal(err)
			}
			var doc franchisejourney.HandoverProfileDocument
			if err = json.Unmarshal(raw, &doc); err != nil {
				t.Fatal(err)
			}
			doc.AlgorithmRevision = franchisejourney.CommercialReleaseAlgorithmRevision
			doc.Options = franchisejourney.SupportedCommercialReleaseOptions()
			raw, err = json.Marshal(doc)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(env["HANDOVER_PROFILE_FILE"], raw, 0600); err != nil {
				t.Fatal(err)
			}
			sum := sha256.Sum256(raw)
			env["HANDOVER_PROFILE_SHA256"] = hex.EncodeToString(sum[:])
		}
		module, err := selectedInitialHandoverModule(&pgxpool.Pool{}, runtime, func(k string) string { return env[k] })
		if err != nil || module == nil {
			t.Fatal("host activation", err)
		}
		mux := http.NewServeMux()
		module.Register(mux, nil)
		// No credentials or database access: a mounted protected GET rejects
		// missing authorization, while a disabled route is absent altogether.
		for _, route := range []string{"commercial-release-result", "commercial-release-current"} {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest("GET", "/v1/franchise/handovers/fixture/"+route+"?organization_id=store", nil))
			want := http.StatusNotFound
			if commercial {
				want = http.StatusUnauthorized
			}
			if rec.Code != want {
				t.Fatalf("commercial=%v route=%s status=%d want=%d", commercial, route, rec.Code, want)
			}
		}
	}
}
