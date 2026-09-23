package httpapi

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type opsRepo struct{ created int }

func (o *opsRepo) CreatePurchaseOrder(context.Context, string, string, operations.PurchaseOrder) error {
	o.created++
	return nil
}
func (o *opsRepo) TransitionPurchaseOrder(context.Context, string, string, string, string, int64, string, string) error {
	return operations.ErrConflict
}
func (o *opsRepo) CreateProductionUnit(context.Context, string, string, operations.ProductionUnit) error {
	return nil
}
func (o *opsRepo) TransitionProductionUnit(context.Context, string, string, string, string, string, string) error {
	return nil
}
func (o *opsRepo) CreateStockUnit(context.Context, string, string, operations.StockUnit) error {
	return nil
}
func (o *opsRepo) TransitionStockUnit(context.Context, string, string, string, string, int64, string, string) error {
	return nil
}

type opsIDs struct{ n int }

func (i *opsIDs) New() string {
	i.n++
	return []string{"018f4d4a-7b36-7a21-8d10-2f4c54c28201", "018f4d4a-7b36-7a21-8d10-2f4c54c28202", "018f4d4a-7b36-7a21-8d10-2f4c54c28203"}[i.n-1]
}

type opsVerifier struct{}

func (opsVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "operator", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28200", Permissions: map[string]struct{}{"procurement:write": {}}, Organizations: map[string]struct{}{"o": {}}}, nil
}
func TestOperationsHTTPAuthorizationAndConflict(t *testing.T) {
	repo := &opsRepo{}
	service := operations.NewService(repo, &opsIDs{})
	module := OperationsModule{Service: service}
	mux := http.NewServeMux()
	module.Register(mux, opsVerifier{})
	request := httptest.NewRequest("POST", "/v1/procurement/purchase-orders", strings.NewReader(`{"supplier_id":"s","destination_organization_id":"o","currency":"USD","total_minor_units":100}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.created != 1 {
		t.Fatalf("create status=%d calls=%d body=%s", response.Code, repo.created, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/procurement/purchase-orders/po/transitions", strings.NewReader(`{"organization_id":"o","current":"draft","target":"submitted","version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 409 {
		t.Fatalf("conflict status=%d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/inventory/stock", strings.NewReader(`{}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("permission status=%d", response.Code)
	}
}
