package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
)

type commercialReleaseSpy struct {
	calls         int
	tenant, actor string
	command       franchisejourney.CommitCommercialReleaseCommand
}

func (s *commercialReleaseSpy) CommitCommercialRelease(_ context.Context, tenant, actor string, c franchisejourney.CommitCommercialReleaseCommand) (franchisejourney.CommercialReleaseReceipt, bool, error) {
	s.calls++
	s.tenant = tenant
	s.actor = actor
	s.command = c
	return franchisejourney.CommercialReleaseReceipt{Effect: franchisejourney.CommercialReleaseEffect}, true, nil
}
func (s *commercialReleaseSpy) CommercialReleaseResult(_ context.Context, tenant, org, handover, key string) (franchisejourney.CommercialReleaseReceipt, error) {
	s.calls++
	s.tenant = tenant
	return franchisejourney.CommercialReleaseReceipt{}, nil
}
func (s *commercialReleaseSpy) ValidateCommercialRelease(_ context.Context, tenant, org, handover string) (franchisejourney.CurrentCommercialRelease, error) {
	s.calls++
	s.tenant = tenant
	return franchisejourney.CurrentCommercialRelease{Current: false}, nil
}
func TestCommercialReleaseHTTPAuthorizationAndHistoricalReplay(t *testing.T) {
	valid := `{"organization_id":"store","observation_sha256":"` + strings.Repeat("a", 64) + `"}`
	for _, name := range []string{"authorized", "unauthenticated", "permission", "organization", "forged-actor", "forged-money", "forged-order"} {
		t.Run(name, func(t *testing.T) {
			p := initialHandoverPrincipal()
			auth := "Bearer fixture"
			body := valid
			want := 201
			switch name {
			case "unauthenticated":
				auth = ""
				want = 401
			case "permission":
				p.Permissions = map[string]struct{}{}
				want = 403
			case "organization":
				p.Organizations = map[string]struct{}{}
				want = 403
			case "forged-actor", "forged-money", "forged-order":
				field := map[string]string{"forged-actor": "actor_subject", "forged-money": "amount_minor_units", "forged-order": "order_id"}[name]
				body = strings.TrimSuffix(valid, "}") + `,"` + field + `":"forged"}`
				want = 400
			}
			s := &commercialReleaseSpy{}
			mux := http.NewServeMux()
			CommercialReleaseModule{Service: s}.Register(mux, journeyVerifier{principal: p})
			r := httptest.NewRequest("POST", "/v1/franchise/handovers/server-handover/commercial-release", strings.NewReader(body))
			r.Header.Set("Authorization", auth)
			r.Header.Set("Idempotency-Key", "release-key")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != want {
				t.Fatalf("got %d want %d", w.Code, want)
			}
			if want != 201 {
				if s.calls != 0 {
					t.Fatal("unauthorized owner call")
				}
				return
			}
			if s.tenant != "verified-tenant" || s.actor != "verified-actor" || s.command.HandoverID != "server-handover" || s.command.IdempotencyKey != "release-key" || w.Header().Get("Idempotency-Replayed") != "true" || strings.Contains(w.Body.String(), `"current":true`) {
				t.Fatal("scope or replay semantics")
			}
			for _, route := range []string{"commercial-release-result", "commercial-release-current"} {
				r = httptest.NewRequest("GET", "/v1/franchise/handovers/server-handover/"+route+"?organization_id=store", nil)
				r.Header.Set("Authorization", auth)
				r.Header.Set("Idempotency-Key", "release-key")
				w = httptest.NewRecorder()
				mux.ServeHTTP(w, r)
				if w.Code != 200 || w.Header().Get("Cache-Control") != "no-store" || strings.Contains(w.Body.String(), `"current":true`) {
					t.Fatal("recovery grants current authority", w.Code)
				}
			}
		})
	}
}
