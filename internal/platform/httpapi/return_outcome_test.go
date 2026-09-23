package httpapi

import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type outcomeHTTPRepo struct {
	journeyRepo
	calls             int
	tenant, org, auth string
}

func (r *outcomeHTTPRepo) ReturnOutcome(_ context.Context, tenant, org, auth string) (franchisejourney.ReturnOutcome, error) {
	r.calls++
	r.tenant, r.org, r.auth = tenant, org, auth
	return franchisejourney.ReturnOutcome{AuthorizationID: auth, OrganizationID: org, OrderID: "order", ObservedAt: time.Unix(1000, 0), Stages: []franchisejourney.ReturnOutcomeStage{}}, nil
}
func TestCommercialCareOutcomeHTTPAuthorization(t *testing.T) {
	verifier, token := confirmationTestIssuer(t)
	repo := &outcomeHTTPRepo{}
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(repo, nil, nil)}.Register(mux, verifier)
	for _, tc := range []struct {
		name, path  string
		perms, orgs []string
		want        int
	}{{"valid", "?organization_id=store", []string{"handover:manage"}, []string{"store"}, 200}, {"readOnly", "?organization_id=store", []string{"admin:read"}, []string{"store"}, 403}, {"wrongOrg", "?organization_id=other", []string{"handover:manage"}, []string{"store"}, 403}, {"browserTenant", "?organization_id=store&tenant_id=other", []string{"handover:manage"}, []string{"store"}, 400}, {"duplicateOrg", "?organization_id=store&organization_id=other", []string{"handover:manage"}, []string{"store"}, 400}} {
		t.Run(tc.name, func(t *testing.T) {
			before := repo.calls
			req := httptest.NewRequest("GET", "/v1/franchise/returns/auth/outcome"+tc.path, nil)
			req.Header.Set("Authorization", "Bearer "+token("operator", "018f4d4a-7b36-7a21-8d10-2f4c54c29b01", tc.perms, tc.orgs))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tc.want || w.Header().Get("Cache-Control") != "no-store" {
				t.Fatalf("status=%d body=%s", w.Code, w.Body)
			}
			if tc.want != 200 && repo.calls != before {
				t.Fatal("denied reached repository")
			}
			if tc.want == 200 && (repo.tenant != "018f4d4a-7b36-7a21-8d10-2f4c54c29b01" || repo.org != "store" || repo.auth != "auth") {
				t.Fatal("verified scope was not used")
			}
		})
	}
}
