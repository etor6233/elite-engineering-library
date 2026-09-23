package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type inventoryRepo struct{ atp, reserve int }

func (f *inventoryRepo) AvailableToPromise(context.Context, string, string, string, time.Time) (inventorycontrol.ATP, error) {
	f.atp++
	return inventorycontrol.ATP{OrganizationID: "a", VariantID: "v", AvailableToPromise: 3}, nil
}
func (f *inventoryRepo) Reserve(_ context.Context, _, _, _ string, value inventorycontrol.Reservation, _ int64) (inventorycontrol.Reservation, error) {
	f.reserve++
	return value, nil
}
func (f *inventoryRepo) Release(context.Context, string, string, string, int64, string) error {
	return inventorycontrol.ErrConflict
}
func (f *inventoryRepo) CreateTransfer(_ context.Context, _, _ string, value inventorycontrol.Transfer, _ map[string]int64, _ string) (inventorycontrol.Transfer, error) {
	return value, nil
}
func (f *inventoryRepo) TransitionTransfer(context.Context, string, string, string, string, int64, string, string) error {
	return nil
}

type inventoryIDs struct{}

func (inventoryIDs) New() string { return "018f4d4a-7b36-7a21-8d10-2f4c54c29999" }

type inventoryVerifier struct{}

func (inventoryVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "operator", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29990", Permissions: map[string]struct{}{"inventory:read": {}, "inventory:write": {}}, Organizations: map[string]struct{}{"a": {}, "b": {}}}, nil
}

func TestInventoryControlHTTPATPAndScope(t *testing.T) {
	repo := &inventoryRepo{}
	mux := http.NewServeMux()
	InventoryControlModule{Service: inventorycontrol.NewService(repo, inventoryIDs{})}.Register(mux, inventoryVerifier{})
	request := httptest.NewRequest("GET", "/v1/inventory/atp?organization_id=a&variant_id=v&horizon=2030-01-01T00:00:00Z", nil)
	request.Header.Set("Authorization", "Bearer valid")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || repo.atp != 1 || !strings.Contains(response.Body.String(), `"available_to_promise":3`) {
		t.Fatalf("ATP status=%d calls=%d body=%s", response.Code, repo.atp, response.Body.String())
	}
	request = httptest.NewRequest("GET", "/v1/inventory/atp?organization_id=forbidden&variant_id=v&horizon=2030-01-01T00:00:00Z", nil)
	request.Header.Set("Authorization", "Bearer valid")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.atp != 1 {
		t.Fatalf("scope status=%d calls=%d", response.Code, repo.atp)
	}
}

func TestInventoryControlHTTPRejectsCustomerOrderAndMapsConflict(t *testing.T) {
	repo := &inventoryRepo{}
	mux := http.NewServeMux()
	InventoryControlModule{Service: inventorycontrol.NewService(repo, inventoryIDs{})}.Register(mux, inventoryVerifier{})
	request := httptest.NewRequest("POST", "/v1/inventory/reservations", strings.NewReader(`{"organization_id":"a","stock_unit_id":"s","variant_id":"v","demand_kind":"customer-order","demand_id":"o","stock_version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 400 || repo.reserve != 0 {
		t.Fatalf("customer-order status=%d calls=%d", response.Code, repo.reserve)
	}
	request = httptest.NewRequest("POST", "/v1/inventory/reservations/r/release", strings.NewReader(`{"organization_id":"a","version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 409 {
		t.Fatalf("release conflict status=%d body=%s", response.Code, response.Body.String())
	}
}
