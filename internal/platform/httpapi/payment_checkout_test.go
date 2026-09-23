package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/platform/identity"
)

func TestPaymentRequestsCannotLeaveConfiguredHostScope(t *testing.T) {
	for _, tenant := range []string{"tenant", "other"} {
		organization := "other"
		if tenant == "other" {
			organization = "store"
		}
		p := identity.Principal{Subject: "operator", TenantID: tenant, Permissions: map[string]struct{}{"payment:create": {}, "payment:write": {}}, Organizations: map[string]struct{}{organization: {}}}
		mux := http.NewServeMux()
		CommerceModule{PaymentProvider: "stripe", PaymentTenantID: "tenant", PaymentOrganizationID: "store"}.Register(mux, paymentPrincipal{p})
		for _, path := range []string{"/v1/commerce/orders/order/payment-request", "/v1/commerce/orders/order/payments", "/v1/payments/payment/transitions"} {
			request := httptest.NewRequest("POST", path, strings.NewReader(`{"organization_id":"`+organization+`"}`))
			request.Header.Set("Authorization", "Bearer fixture")
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != 503 || !strings.Contains(response.Body.String(), "PAYMENT_SCOPE_DISABLED") {
				t.Fatal(path, response.Code, response.Body.String())
			}
		}
	}
}

type checkoutReadFixture struct {
	calls int
	value paymentbridge.CustomerCheckout
	err   error
	scope []string
}

func (f *checkoutReadFixture) CustomerCheckout(_ context.Context, tenant, org, subject, order string) (paymentbridge.CustomerCheckout, error) {
	f.calls++
	f.scope = []string{tenant, org, subject, order}
	if tenant != "tenant" || subject != "alice" || order != "order" {
		return paymentbridge.CustomerCheckout{}, paymentbridge.ErrCheckoutNotAvailable
	}
	return f.value, f.err
}
func TestCustomerCheckoutHTTPAuthorizationAndProviderURL(t *testing.T) {
	base := identity.Principal{Subject: "alice", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	for _, tc := range []struct {
		name, path, authorization, subject, tenant string
		permission                                 bool
		expected, calls                            int
	}{
		{"allowed", "/v1/customer/orders/order/checkout?organization_id=store", "Bearer fixture", "alice", "tenant", true, 200, 1},
		{"anonymous", "/v1/customer/orders/order/checkout?organization_id=store", "", "alice", "tenant", true, 401, 0},
		{"no-permission", "/v1/customer/orders/order/checkout?organization_id=store", "Bearer fixture", "alice", "tenant", false, 403, 0},
		{"foreign-org", "/v1/customer/orders/order/checkout?organization_id=other", "Bearer fixture", "alice", "tenant", true, 403, 0},
		{"foreign-customer", "/v1/customer/orders/order/checkout?organization_id=store", "Bearer fixture", "bob", "tenant", true, 404, 1},
		{"foreign-tenant", "/v1/customer/orders/order/checkout?organization_id=store", "Bearer fixture", "alice", "other", true, 404, 1},
		{"unknown-order", "/v1/customer/orders/other/checkout?organization_id=store", "Bearer fixture", "alice", "tenant", true, 404, 1},
		{"injected-subject", "/v1/customer/orders/order/checkout?organization_id=store&customer_id=alice", "Bearer fixture", "bob", "tenant", true, 400, 0},
		{"injected-url", "/v1/customer/orders/order/checkout?organization_id=store&url=https://evil.test", "Bearer fixture", "alice", "tenant", true, 400, 0},
		{"duplicate-org", "/v1/customer/orders/order/checkout?organization_id=store&organization_id=store", "Bearer fixture", "alice", "tenant", true, 400, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := base
			p.Subject = tc.subject
			p.TenantID = tc.tenant
			if !tc.permission {
				p.Permissions = nil
			}
			reader := &checkoutReadFixture{value: paymentbridge.CustomerCheckout{OrderID: "order", ProviderCode: "stripe", URL: "https://checkout.stripe.com/c/pay/cs_test_fixture", ExpiresAt: time.Now().Add(time.Hour)}}
			mux := http.NewServeMux()
			PaymentCheckoutModule{Reader: reader}.Register(mux, paymentPrincipal{p})
			request := httptest.NewRequest("GET", tc.path, nil)
			request.Header.Set("Authorization", tc.authorization)
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != tc.expected || reader.calls != tc.calls {
				t.Fatalf("status=%d calls=%d", response.Code, reader.calls)
			}
			if reader.calls > 0 && (reader.scope[0] != p.TenantID || reader.scope[2] != p.Subject) {
				t.Fatal("scope did not come from principal")
			}
			if response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("Referrer-Policy") != "no-referrer" {
				t.Fatal("missing private response policy")
			}
		})
	}
	for _, tc := range []struct {
		name, url string
		err       error
		nilReader bool
	}{
		{"invalid-provider-host", "https://checkout.stripe.com.evil.test/session", nil, false},
		{"unsafe-url", "javascript:alert(1)", nil, false},
		{"reader-error", "https://checkout.stripe.com/session", errors.New("PRIVATE_DATABASE_DETAIL"), false},
		{"disabled", "", nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := &checkoutReadFixture{value: paymentbridge.CustomerCheckout{OrderID: "order", ProviderCode: "stripe", URL: tc.url, ExpiresAt: time.Now().Add(time.Hour)}, err: tc.err}
			module := PaymentCheckoutModule{Reader: reader}
			if tc.nilReader {
				module.Reader = nil
			}
			mux := http.NewServeMux()
			module.Register(mux, paymentPrincipal{base})
			request := httptest.NewRequest("GET", "/v1/customer/orders/order/checkout?organization_id=store", nil)
			request.Header.Set("Authorization", "Bearer fixture")
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != 503 || strings.Contains(response.Body.String(), "PRIVATE_DATABASE_DETAIL") {
				t.Fatal(response.Code, response.Body.String())
			}
		})
	}
}
