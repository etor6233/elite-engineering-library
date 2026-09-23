package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/platform/identity"
)

type fiscalRepoFake struct{ requests int }

func (*fiscalRepoFake) ConfigurePointOfSale(context.Context, string, string, fiscal.PointOfSale) error {
	return nil
}
func (r *fiscalRepoFake) RequestInvoice(_ context.Context, _, _, _, _ string, value fiscal.Invoice) (fiscal.Invoice, bool, error) {
	r.requests++
	value.Status = "queued"
	value.TotalMinorUnits = value.NetMinorUnits + value.VATMinorUnits
	return value, false, nil
}
func (*fiscalRepoFake) GetInvoice(context.Context, string, string, string) (fiscal.Invoice, error) {
	return fiscal.Invoice{ID: "invoice"}, nil
}

type fiscalVerifier struct{ p identity.Principal }

func (v fiscalVerifier) Verify(context.Context, string) (identity.Principal, error) { return v.p, nil }

type fiscalIDs struct{ n int }

func (i *fiscalIDs) New() string { i.n++; return "id" + string(rune('0'+i.n)) }
func TestFiscalHTTPRejectsCrossScopeAndQueuesAuthorizedScope(t *testing.T) {
	repo := &fiscalRepoFake{}
	service := fiscal.NewService(repo, &fiscalIDs{})
	mux := http.NewServeMux()
	principal := identity.Principal{Subject: "controller", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29f12", Permissions: map[string]struct{}{"fiscal:issue": {}}, Organizations: map[string]struct{}{"franchise": {}}}
	FiscalModule{Service: service}.Register(mux, fiscalVerifier{p: principal})
	body := `{"organization_id":"other","order_id":"order","payment_attempt_id":"payment","point_of_sale_id":"pos","voucher_type":6,"concept":1,"recipient_document_type":99,"recipient_document":"0","recipient_vat_condition_id":5,"net_minor_units":10000,"vat_minor_units":2100,"vat_lines":[{"id":5,"base_minor_units":10000,"amount_minor_units":2100}],"issued_on":"2026-08-30T00:00:00Z"}`
	request := httptest.NewRequest("POST", "/v1/fiscal/invoices", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "1234567890abcdef")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.requests != 0 {
		t.Fatalf("status=%d requests=%d", response.Code, repo.requests)
	}
	request = httptest.NewRequest("POST", "/v1/fiscal/invoices", strings.NewReader(strings.Replace(body, `"other"`, `"franchise"`, 1)))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "1234567890abcdef")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 202 || repo.requests != 1 {
		t.Fatalf("status=%d requests=%d body=%s", response.Code, repo.requests, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/fiscal/invoices", strings.NewReader(strings.Replace(body, `"other"`, `"franchise"`, 1)+` {}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "1234567890abcdef")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 400 || repo.requests != 1 {
		t.Fatalf("trailing JSON status=%d requests=%d", response.Code, repo.requests)
	}
}
