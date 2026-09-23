package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/platform/identity"
)

type parameterHTTPRepository struct {
	decisions int
	schedules int
	snapshots int
}

func (r *parameterHTTPRepository) StoreParameterSnapshot(_ context.Context, value fiscal.ParameterSnapshot, _ string) (fiscal.ParameterSnapshot, bool, error) {
	r.snapshots++
	return value, false, nil
}
func (r *parameterHTTPRepository) GetParameterSnapshot(context.Context, string, string, string) (fiscal.ParameterSnapshot, error) {
	return fiscal.ParameterSnapshot{ID: "snapshot"}, nil
}
func (r *parameterHTTPRepository) DecideParameterSnapshot(context.Context, string, string, string, string, bool, string, string) error {
	r.decisions++
	return nil
}
func (r *parameterHTTPRepository) ConfigureParameterSchedule(_ context.Context, value fiscal.ParameterSchedule, _, _ string) (fiscal.ParameterSchedule, error) {
	r.schedules++
	return value, nil
}

type parameterHTTPProvider struct{}

func (parameterHTTPProvider) FetchParameters(context.Context, string, string, *string) (fiscal.ParameterSnapshot, error) {
	return fiscal.ParameterSnapshot{TaxpayerCUIT: "33693450239", Kind: "vat_rate", Items: []fiscal.ParameterItem{{Code: "5"}}, ResponseHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}, nil
}

func TestFiscalParameterAdminSeparatesPermissionsAndScope(t *testing.T) {
	repository := &parameterHTTPRepository{}
	registry := fiscal.NewParameterRegistry(repository, parameterHTTPProvider{}, &fiscalIDs{}, time.Now)
	p := identity.Principal{Subject: "tax-owner", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29f12", Organizations: map[string]struct{}{"franchise": {}}, Permissions: map[string]struct{}{"fiscal:parameters:refresh": {}, "fiscal:parameters:approve": {}, "fiscal:parameters:configure": {}}}
	mux := http.NewServeMux()
	FiscalModule{Service: fiscal.NewService(&fiscalRepoFake{}, &fiscalIDs{}), Parameters: registry}.Register(mux, fiscalVerifier{p: p})
	cases := []struct {
		path, body string
		want       int
	}{
		{"/v1/fiscal/parameters/refresh", `{"organization_id":"other","taxpayer_cuit":"33693450239","kind":"vat_rate","voucher_class":null}`, 403},
		{"/v1/fiscal/parameters/refresh", `{"organization_id":"franchise","taxpayer_cuit":"33693450239","kind":"vat_rate","voucher_class":null}`, 201},
		{"/v1/fiscal/parameters/snapshot/decisions", `{"organization_id":"franchise","approved":true,"reason":"validated against ARCA configuration"}`, 204},
		{"/v1/fiscal/parameter-schedules", `{"organization_id":"franchise","taxpayer_cuit":"33693450239","kind":"vat_rate","voucher_class":null,"interval_seconds":900,"next_run_at":"2026-08-31T00:00:00Z"}`, 201},
	}
	for _, tc := range cases {
		request := httptest.NewRequest("POST", tc.path, strings.NewReader(tc.body))
		request.Header.Set("Authorization", "Bearer valid")
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		if response.Code != tc.want {
			t.Fatalf("%s status=%d body=%s", tc.path, response.Code, response.Body.String())
		}
	}
	if repository.snapshots != 1 || repository.decisions != 1 || repository.schedules != 1 {
		t.Fatalf("snapshots=%d decisions=%d schedules=%d", repository.snapshots, repository.decisions, repository.schedules)
	}
}
