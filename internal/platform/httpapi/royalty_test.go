package httpapi

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/royalty"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type royaltyRepo struct {
	policies int
	accruals int
}

func (r *royaltyRepo) CreatePolicy(context.Context, string, string, royalty.Policy) error {
	r.policies++
	return nil
}
func (r *royaltyRepo) AccruePayment(_ context.Context, _, _ string, value royalty.PaymentEvent, _ string) (royalty.Accrual, error) {
	r.accruals++
	return royalty.Accrual{ID: value.ID, OrganizationID: value.OrganizationID, SourceState: value.State}, nil
}
func (*royaltyRepo) OpenSettlement(context.Context, string, string, royalty.Settlement) error {
	return nil
}
func (*royaltyRepo) CloseSettlement(context.Context, string, string, string, string, int64) (royalty.Settlement, error) {
	return royalty.Settlement{Status: "closed"}, nil
}
func (*royaltyRepo) ReverseSettlement(context.Context, string, string, string, string, string, int64, string) (royalty.Settlement, error) {
	return royalty.Settlement{Status: "closed", ReversalOf: "s"}, nil
}
func (*royaltyRepo) Reconcile(_ context.Context, _, _ string, value royalty.Reconciliation, _, _ string) (royalty.Reconciliation, error) {
	value.Status = "matched"
	return value, nil
}

type royaltyIDs struct{ n int }

func (i *royaltyIDs) New() string { i.n++; return "id" + string(rune('0'+i.n)) }

type royaltyVerifier struct{ principal identity.Principal }

func (v royaltyVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, nil
}

func TestRoyaltyHTTPAuthorizationScopeAndStrictCommands(t *testing.T) {
	repo := &royaltyRepo{}
	service := royalty.NewService(repo, &royaltyIDs{})
	mux := http.NewServeMux()
	principal := identity.Principal{Subject: "controller", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29d00", Permissions: map[string]struct{}{"royalty:policy": {}, "royalty:post": {}}, Organizations: map[string]struct{}{"franchise": {}}}
	RoyaltyModule{Service: service}.Register(mux, royaltyVerifier{principal: principal})
	now := time.Now().UTC().Truncate(time.Second)
	request := httptest.NewRequest("POST", "/v1/royalties/policies", strings.NewReader(`{"agreement_id":"agreement","organization_id":"franchise","currency":"ARS","rate_basis_points":650,"valid_from":"`+now.Format(time.RFC3339)+`"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.policies != 1 {
		t.Fatalf("policy status=%d count=%d body=%s", response.Code, repo.policies, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/royalties/accruals/from-payment", strings.NewReader(`{"organization_id":"other","payment_attempt_id":"payment","payment_expected_version":3,"state":"captured","occurred_at":"`+now.Format(time.RFC3339)+`"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "provider:event")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.accruals != 0 {
		t.Fatalf("cross-scope status=%d accruals=%d", response.Code, repo.accruals)
	}
	request = httptest.NewRequest("POST", "/v1/royalties/settlements/s/close", strings.NewReader(`{"organization_id":"franchise","expected_version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("permission status=%d", response.Code)
	}
}
