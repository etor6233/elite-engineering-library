package httpapi

import (
	"context"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type accountingRepo struct{ journals int }

func (*accountingRepo) CreateAccount(context.Context, string, string, accounting.Account) error {
	return nil
}
func (*accountingRepo) OpenPeriod(context.Context, string, string, accounting.Period) error {
	return nil
}
func (r *accountingRepo) CreateJournal(_ context.Context, _, _ string, _ accounting.Journal) error {
	r.journals++
	return nil
}
func (*accountingRepo) PostJournal(context.Context, string, string, string, string, int64, string, string) (accounting.Journal, error) {
	return accounting.Journal{Status: "posted"}, nil
}
func (*accountingRepo) ReverseJournal(context.Context, string, string, string, string, string, string, int64, string, string) (accounting.Journal, error) {
	return accounting.Journal{Status: "posted"}, nil
}
func (*accountingRepo) ClosePeriod(context.Context, string, string, int64, string) (accounting.Period, error) {
	return accounting.Period{Status: "closed"}, nil
}
func (*accountingRepo) TrialBalance(context.Context, string, string, string) ([]accounting.Balance, error) {
	return []accounting.Balance{}, nil
}

type accountingIDs struct{ n int }

func (i *accountingIDs) New() string { i.n++; return "id" + string(rune('0'+i.n)) }

type accountingVerifier struct{ p identity.Principal }

func (v accountingVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.p, nil
}
func TestAccountingHTTPRejectsCrossScopeAndUnbalanced(t *testing.T) {
	repo := &accountingRepo{}
	service := accounting.NewService(repo, &accountingIDs{})
	mux := http.NewServeMux()
	p := identity.Principal{Subject: "controller", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29f00", Permissions: map[string]struct{}{"accounting:write": {}}, Organizations: map[string]struct{}{"franchise": {}}}
	AccountingModule{Service: service}.Register(mux, accountingVerifier{p: p})
	now := time.Now().UTC()
	body := `{"organization_id":"other","period_id":"p","source_type":"SALE","source_id":"o","currency":"ARS","posting_date":"` + now.Format(time.RFC3339) + `","lines":[{"line_no":1,"account_code":"CASH","debit_minor_units":100},{"line_no":2,"account_code":"REVENUE","credit_minor_units":100}]}`
	request := httptest.NewRequest("POST", "/v1/accounting/journals", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.journals != 0 {
		t.Fatalf("status=%d journals=%d", response.Code, repo.journals)
	}
	request = httptest.NewRequest("POST", "/v1/accounting/journals", strings.NewReader(strings.Replace(body, `"other"`, `"franchise"`, 1)))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.journals != 1 {
		t.Fatalf("status=%d journals=%d body=%s", response.Code, repo.journals, response.Body.String())
	}
}
