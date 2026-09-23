package httpapi

import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/socialbridge"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type socialTestVerifier struct{}

func (socialTestVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	if token != "verified" {
		return identity.Principal{}, identity.ErrUnauthenticated
	}
	return identity.Principal{Subject: "verified-operator", TenantID: "tenant"}, nil
}

type socialTestService struct {
	decisions int
	principal string
	approved  bool
}

func (s *socialTestService) Submit(context.Context, identity.Principal, socialbridge.Request) (string, bool, error) {
	return "hash", false, nil
}
func (s *socialTestService) Decide(_ context.Context, p identity.Principal, _, _ string, a bool, _ string) (approval.State, error) {
	s.decisions++
	s.principal = p.Subject
	s.approved = a
	return approval.StateApproved, nil
}
func (s *socialTestService) Status(context.Context, identity.Principal, string) (postgres.SocialStatus, error) {
	return postgres.SocialStatus{}, nil
}
func (s *socialTestService) Reconcile(context.Context, identity.Principal, string) error { return nil }
func TestSocialHTTPExplicitDecisionAndVerifiedIdentity(t *testing.T) {
	s := &socialTestService{}
	h := NewSocialPublishing(s, socialTestVerifier{})
	call := func(token, body string) int {
		r := httptest.NewRequest(http.MethodPost, "/v1/social/requests/request/decision", strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code
	}
	for _, body := range []string{`{"request_sha256":"hash","approve":null,"reason":"x"}`, `{"request_sha256":"hash","reason":"x"}`, `{"request_sha256":"hash","approve":true,"Approve":false}`, `{"request_sha256":"hash","approve":true,"requester":"attacker"}`} {
		if code := call("verified", body); code != 400 {
			t.Fatalf("body %s code %d", body, code)
		}
	}
	if s.decisions != 0 {
		t.Fatal("invalid decision reached owner")
	}
	valid := `{"request_sha256":"hash","approve":true,"reason":"reviewed"}`
	if code := call("unverified", valid); code != 401 {
		t.Fatal(code)
	}
	if code := call("verified", valid); code != 200 || s.decisions != 1 || s.principal != "verified-operator" || !s.approved {
		t.Fatal(code, s)
	}
}
