package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/platform/identity"
)

type queryRepository struct{ customer string }

func (*queryRepository) Overview(context.Context, string, string) (enterprisequery.Overview, error) {
	return enterprisequery.Overview{}, nil
}
func (r *queryRepository) Orders(_ context.Context, _, _, customer string, _ int, _ string) (enterprisequery.Page[enterprisequery.Order], error) {
	r.customer = customer
	return enterprisequery.Page[enterprisequery.Order]{}, nil
}
func (*queryRepository) FactoryUnits(context.Context, string, string, int, string) (enterprisequery.Page[enterprisequery.FactoryUnit], error) {
	return enterprisequery.Page[enterprisequery.FactoryUnit]{}, nil
}
func (*queryRepository) ServiceCases(context.Context, string, string, string, int, string) (enterprisequery.Page[enterprisequery.ServiceCase], error) {
	return enterprisequery.Page[enterprisequery.ServiceCase]{}, nil
}

type queryVerifier struct{}

func (queryVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "customer-1", TenantID: "tenant-1", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"org-a": {}}}, nil
}

func TestCustomerQueryBindsSubjectAndOrganization(t *testing.T) {
	repository := &queryRepository{}
	mux := http.NewServeMux()
	EnterpriseQueryModule{Service: enterprisequery.NewService(repository)}.Register(mux, queryVerifier{})
	request := httptest.NewRequest("GET", "/v1/customer/orders?organization_id=org-a", nil)
	request.Header.Set("Authorization", "Bearer valid")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || repository.customer != "customer-1" {
		t.Fatalf("status=%d customer=%q", response.Code, repository.customer)
	}
	request = httptest.NewRequest("GET", "/v1/customer/orders?organization_id=org-b", nil)
	request.Header.Set("Authorization", "Bearer valid")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("cross-organization query status=%d", response.Code)
	}
}
