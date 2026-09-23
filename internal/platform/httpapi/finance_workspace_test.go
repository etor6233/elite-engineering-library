package httpapi

import (
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFinanceHTTPBoundaryBeforeStore(t *testing.T) {
	p := identity.Principal{Subject: "author", TenantID: "tenant", Permissions: map[string]struct{}{"accounting:read": {}}, Organizations: map[string]struct{}{"store": {}}}
	m := http.NewServeMux()
	FinanceWorkspaceModule{}.Register(m, accountingVerifier{p: p})
	cases := []struct {
		method, path, body, token string
		status                    int
	}{
		{"GET", "/v1/finance/workspace?organization_id=store", "", "", 401},
		{"GET", "/v1/finance/workspace?organization_id=foreign", "", "valid", 403},
		{"GET", "/v1/finance/workspace?organization_id=store&limit=51&view=accounts", "", "valid", 400},
		{"GET", "/v1/finance/workspace?organization_id=store&organization_id=store", "", "valid", 400},
		{"GET", "/v1/finance/journals/j/lines?organization_id=foreign", "", "valid", 403},
		{"GET", "/v1/finance/journals/j/lines?organization_id=store&cursor=abc", "", "valid", 400},
		{"GET", "/v1/finance/trial-balance?organization_id=store&period_id=p&extra=1", "", "valid", 400},
		{"GET", "/v1/finance/trial-balance?organization_id=foreign&period_id=p", "", "valid", 403},
		{"GET", "/v1/finance/commands/result?organization_id=foreign", "", "valid", 403},
		{"POST", "/v1/finance/commands", `{"action":"accounting_account","organization_id":"store","payload":{"code":"CASH","name":"Caja","type":"asset"}}`, "valid", 403},
		{"POST", "/v1/finance/commands", `{"action":"accounting_account","organization_id":"foreign","payload":{}}`, "valid", 403},
		{"POST", "/v1/finance/commands", strings.Repeat("x", 16385), "valid", 413},
		{"POST", "/v1/finance/commands", `{"action":"accounting_account","organization_id":"store","tenant_id":"foreign","payload":{}}`, "valid", 400},
	}
	for _, c := range cases {
		t.Run(c.method+" "+c.path+" "+c.token, func(t *testing.T) {
			r := httptest.NewRequest(c.method, c.path, strings.NewReader(c.body))
			if c.token != "" {
				r.Header.Set("Authorization", "Bearer "+c.token)
			}
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Idempotency-Key", "finance-boundary-key")
			w := httptest.NewRecorder()
			m.ServeHTTP(w, r)
			if w.Code != c.status {
				t.Fatalf("status %d expected %d: %s", w.Code, c.status, w.Body.String())
			}
		})
	}
}
