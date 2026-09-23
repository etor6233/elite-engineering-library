package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/platform/identity"
	wc "elite.local/enterprise/internal/warrantyclaim"
)

type warrantyTestVerifier struct{ p identity.Principal }

func (v warrantyTestVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	if token != "fixture" {
		return identity.Principal{}, identity.ErrUnauthenticated
	}
	return v.p, nil
}

type warrantyTestService struct {
	WarrantyService
	calls     int
	principal identity.Principal
	request   wc.OfferRequest
}

func (*warrantyTestService) Scope() (string, string) {
	return "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", "store"
}
func (s *warrantyTestService) BindOffer(_ context.Context, p identity.Principal, r wc.OfferRequest) (wc.Offer, error) {
	s.calls++
	s.principal = p
	s.request = r
	return wc.Offer{QuoteID: r.QuoteID, QuoteVersion: r.QuoteVersion, ProfileSHA256: r.ProfileSHA256}, nil
}
func TestWarrantyTransportPreservesExactVersionsAndRejectsAmbiguity(t *testing.T) {
	profile := strings.Repeat("a", 64)
	valid := `{"quote_id":"quote","quote_version":"9007199254740993","profile_sha256":"` + profile + `"}`
	for _, test := range []struct {
		name, body, query, token string
		code                     int
	}{
		{"exact-int64", valid, "organization_id=store", "fixture", 201},
		{"number-not-string", strings.Replace(valid, `"9007199254740993"`, `9007199254740993`, 1), "organization_id=store", "fixture", 400},
		{"overflow", strings.Replace(valid, "9007199254740993", "9223372036854775808", 1), "organization_id=store", "fixture", 400},
		{"duplicate-case", strings.Replace(valid, `"quote_id":"quote"`, `"quote_id":"quote","Quote_ID":"other"`, 1), "organization_id=store", "fixture", 400},
		{"unknown-actor", strings.Replace(valid, `{`, `{"actor":"forged",`, 1), "organization_id=store", "fixture", 400},
		{"path-body-mismatch", strings.Replace(valid, `"quote"`, `"other"`, 1), "organization_id=store", "fixture", 400},
		{"duplicate-scope", valid, "organization_id=store&organization_id=other", "fixture", 400},
		{"other-scope", valid, "organization_id=other", "fixture", 403},
		{"oversized", `{"padding":"` + strings.Repeat("x", 32768) + `"}`, "organization_id=store", "fixture", 413},
		{"unauthenticated", valid, "organization_id=store", "invalid", 401},
	} {
		t.Run(test.name, func(t *testing.T) {
			service := &warrantyTestService{}
			tenant, _ := service.Scope()
			principal := identity.Principal{TenantID: tenant, Subject: "operator", Organizations: map[string]struct{}{"store": {}, "other": {}}, Permissions: map[string]struct{}{"*": {}}}
			mux := http.NewServeMux()
			WarrantyModule{Service: service}.Register(mux, warrantyTestVerifier{p: principal})
			req := httptest.NewRequest("POST", "/v1/franchise/warranty/quotes/quote/terms?"+test.query, strings.NewReader(test.body))
			req.Header.Set("Authorization", "Bearer "+test.token)
			res := httptest.NewRecorder()
			mux.ServeHTTP(res, req)
			if res.Code != test.code {
				t.Fatal(res.Code, res.Body.String())
			}
			if test.code != 201 {
				if service.calls != 0 {
					t.Fatal("rejected body reached owner")
				}
				return
			}
			if service.calls != 1 || service.request.QuoteVersion != 9007199254740993 || service.principal.Subject != "operator" || service.principal.AllowedOrganization("other") || service.principal.Allowed("warranty:approve") {
				t.Fatal("precision or authority changed", service)
			}
			var body map[string]any
			if json.Unmarshal(res.Body.Bytes(), &body) != nil || body["quote_version"] != "9007199254740993" {
				t.Fatal("response lost int64", res.Body.String())
			}
		})
	}
}
