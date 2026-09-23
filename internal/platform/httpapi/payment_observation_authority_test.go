package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/platform/identity"
)

type observationAuthorityRepository struct {
	commerceRepo
	transitions int
}

func (r *observationAuthorityRepository) TransitionPayment(context.Context, string, string, string, string, string, int64, string, string) error {
	r.transitions++
	return nil
}

func TestProviderObservationOwnsFinancialTransitions(t *testing.T) {
	for _, c := range []struct {
		name, current, target, reference string
		status, calls                    int
	}{
		{"operator-cannot-authorize", "pending", "authorized", "pi_fixture", 403, 0},
		{"operator-cannot-capture", "authorized", "captured", "pi_fixture", 403, 0},
		{"operator-cannot-refund", "captured", "refunded", "pi_fixture", 403, 0},
		{"operator-cannot-dispute", "captured", "disputed", "pi_fixture", 403, 0},
		{"operator-cannot-fail-dispatched", "pending", "failed", "", 403, 0},
		{"cancel-does-not-invent-provider-reference", "created", "failed", "pi_forged", 403, 0},
		{"cancel-before-dispatch-remains-supported", "created", "failed", "", 200, 1},
	} {
		t.Run(c.name, func(t *testing.T) {
			repo := &observationAuthorityRepository{}
			p := identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{"payment:write": {}}, Organizations: map[string]struct{}{"org": {}}}
			mux := http.NewServeMux()
			CommerceModule{Service: commerce.NewService(repo, &commerceIDs{}), PaymentProvider: "stripe", ProviderObservedPayments: true}.Register(mux, paymentPrincipal{p})
			body, _ := json.Marshal(map[string]any{"organization_id": "org", "current": c.current, "target": c.target, "version": 1, "provider_reference": c.reference})
			request := httptest.NewRequest("POST", "/v1/payments/attempt/transitions", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer fixture")
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != c.status || repo.transitions != c.calls {
				t.Fatalf("status=%d repository_calls=%d body=%s", response.Code, repo.transitions, response.Body.String())
			}
		})
	}
}
