package httpapi

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type cmsVerifier struct{ p identity.Principal }

func (v cmsVerifier) Verify(context.Context, string) (identity.Principal, error) { return v.p, nil }

type cmsUnread struct{ reads int }

func (v *cmsUnread) Read([]byte) (int, error) { v.reads++; return 0, io.EOF }
func (v *cmsUnread) Close() error             { return nil }
func TestHelpCMSHTTPBoundaries(t *testing.T) {
	for _, permissions := range []map[string]struct{}{{"help:read": {}}, {}} {
		mux := http.NewServeMux()
		HelpCMSModule{}.Register(mux, cmsVerifier{identity.Principal{TenantID: "tenant", Subject: "reader", Permissions: permissions}})
		body := &cmsUnread{}
		r := httptest.NewRequest("POST", "/v1/help/cms/commands", body)
		r.Header.Set("Authorization", "Bearer fixture")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != 403 || body.reads != 0 {
			t.Fatal("unauthorized body consumed", w.Code, body.reads)
		}
	}
	p := identity.Principal{TenantID: "tenant", Subject: "writer", Permissions: map[string]struct{}{"help:write": {}}, Organizations: map[string]struct{}{"org": {}}}
	cases := []struct {
		method, path, body string
		status             int
	}{
		{"POST", "/v1/help/cms/commands", strings.Repeat("x", 32769), 413},
		{"POST", "/v1/help/cms/commands", `{"action":"create","action":"archive"}`, 400},
		{"POST", "/v1/help/cms/commands", `{"command_id":"x","action":"publish","article_id":"guide","organization_id":"org","version":"1"}`, 403},
		{"POST", "/v1/help/cms/commands", `{"command_id":"x","action":"update","article_id":"guide","organization_id":"org","version":1,"title":"Title","body":"Body"}`, 400},
		{"GET", "/v1/help/cms/articles/guide?organization_id=org&organization_id=other", "", 400},
		{"GET", "/v1/help/cms/articles/guide?organization_id=org&version=01", "", 400},
		{"GET", "/v1/help/cms/articles/guide?organization_id=org&version=9223372036854775808", "", 400},
	}
	for _, c := range cases {
		mux := http.NewServeMux()
		HelpCMSModule{}.Register(mux, cmsVerifier{p})
		r := httptest.NewRequest(c.method, c.path, strings.NewReader(c.body))
		r.Header.Set("Authorization", "Bearer fixture")
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != c.status {
			t.Fatal(c.path, w.Code, c.status, w.Body.String())
		}
		if !strings.Contains(w.Header().Get("Cache-Control"), "no-store") {
			t.Fatal("private content cached")
		}
	}
}
