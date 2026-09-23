package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
)

type initialHandoverSpy struct {
	calls                     int
	tenant, actor, order, key string
	replay                    bool
	err                       error
}

func (s *initialHandoverSpy) Prepare(_ context.Context, tenant, actor string, c franchisejourney.PrepareHandoverCommand) (franchisejourney.HandoverPreparation, bool, error) {
	s.calls++
	s.tenant = tenant
	s.actor = actor
	s.order = c.OrderID
	s.key = c.IdempotencyKey
	return franchisejourney.HandoverPreparation{}, s.replay, s.err
}
func (s *initialHandoverSpy) Result(_ context.Context, tenant, organization, order, key string) (franchisejourney.HandoverPreparation, error) {
	s.calls++
	s.tenant = tenant
	s.order = order
	s.key = key
	return franchisejourney.HandoverPreparation{}, s.err
}
func (s *initialHandoverSpy) EvaluateRelease(_ context.Context, tenant, organization, handover, sha string) (franchisejourney.HandoverReferenceRelease, error) {
	s.calls++
	s.tenant = tenant
	return franchisejourney.HandoverReferenceRelease{Scope: "LOCAL_FIXTURES", Eligible: true}, s.err
}
func initialHandoverPrincipal() identity.Principal {
	return identity.Principal{TenantID: "verified-tenant", Subject: "verified-actor", Permissions: map[string]struct{}{"handover:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
}
func TestInitialHandoverHTTPBindsPrincipalAndRecovery(t *testing.T) {
	spy := &initialHandoverSpy{replay: true}
	mux := http.NewServeMux()
	InitialHandoverModule{Service: spy}.Register(mux, journeyVerifier{principal: initialHandoverPrincipal()})
	request := httptest.NewRequest("POST", "/v1/franchise/orders/server-order/handover", strings.NewReader(`{"organization_id":"store","order_line_id":"line","payment_attempt_id":"payment","observation_sha256":"`+strings.Repeat("a", 64)+`"}`))
	request.Header.Set("Authorization", "Bearer fixture")
	request.Header.Set("Idempotency-Key", "recovery-key")
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	if recorder.Code != 201 || recorder.Header().Get("Idempotency-Replayed") != "true" || spy.tenant != "verified-tenant" || spy.actor != "verified-actor" || spy.order != "server-order" || spy.key != "recovery-key" {
		t.Fatalf("binding failed: %d %+v", recorder.Code, spy)
	}
	request = httptest.NewRequest("GET", "/v1/franchise/orders/server-order/handover-result?organization_id=store", nil)
	request.Header.Set("Authorization", "Bearer fixture")
	request.Header.Set("Idempotency-Key", "recovery-key")
	recorder = httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	if recorder.Code != 200 || spy.calls != 2 || spy.key != "recovery-key" {
		t.Fatalf("GET recovery failed: %d %+v", recorder.Code, spy)
	}
}
func TestInitialHandoverHTTPRejectsForgedInputsAndUnauthorizedScope(t *testing.T) {
	valid := `{"organization_id":"store","order_line_id":"line","payment_attempt_id":"payment","observation_sha256":"` + strings.Repeat("a", 64) + `"}`
	for _, name := range []string{"unauthenticated", "wrong-permission", "wrong-organization", "forged-customer", "forged-order", "forged-amount"} {
		t.Run(name, func(t *testing.T) {
			principal := initialHandoverPrincipal()
			body := valid
			authorization := "Bearer fixture"
			want := 400
			switch name {
			case "unauthenticated":
				authorization = ""
				want = 401
			case "wrong-permission":
				principal.Permissions = map[string]struct{}{}
				want = 403
			case "wrong-organization":
				principal.Organizations = map[string]struct{}{"other": {}}
				want = 403
			default:
				field := map[string]string{"forged-customer": "customer_subject", "forged-order": "order_id", "forged-amount": "amount_minor_units"}[name]
				body = strings.TrimSuffix(valid, "}") + `,"` + field + `":"forged"}`
			}
			spy := &initialHandoverSpy{}
			mux := http.NewServeMux()
			InitialHandoverModule{Service: spy}.Register(mux, journeyVerifier{principal: principal})
			request := httptest.NewRequest("POST", "/v1/franchise/orders/order/handover", strings.NewReader(body))
			request.Header.Set("Authorization", authorization)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, request)
			if recorder.Code != want || spy.calls != 0 {
				t.Fatalf("rejected request reached owner: status%d calls%d", recorder.Code, spy.calls)
			}
		})
	}
}
func TestInitialHandoverHTTPUnconfiguredModuleAndErrors(t *testing.T) {
	mux := http.NewServeMux()
	InitialHandoverModule{}.Register(mux, journeyVerifier{})
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest("GET", "/v1/franchise/orders/order/handover-result", nil))
	if recorder.Code != 404 {
		t.Fatal("unconfigured route mounted")
	}
	for _, item := range []struct {
		err    error
		status int
	}{{franchisejourney.ErrConflict, 409}, {franchisejourney.ErrNotFound, 404}, {franchisejourney.ErrReleaseConditioned, 503}, {errors.New("private database information"), 500}} {
		spy := &initialHandoverSpy{err: item.err}
		mux = http.NewServeMux()
		InitialHandoverModule{Service: spy}.Register(mux, journeyVerifier{principal: initialHandoverPrincipal()})
		request := httptest.NewRequest("GET", "/v1/franchise/handovers/handover/release-check?organization_id=store&observation_sha256="+strings.Repeat("a", 64), nil)
		request.Header.Set("Authorization", "Bearer fixture")
		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, request)
		if recorder.Code != item.status || strings.Contains(recorder.Body.String(), "private database information") {
			t.Fatalf("unsafe mapping %d %s", recorder.Code, recorder.Body.String())
		}
	}
}
