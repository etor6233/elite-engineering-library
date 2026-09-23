package httpapi

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type workspaceHTTPVerifier struct{}

func (workspaceHTTPVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	if token == "bad" {
		return identity.Principal{}, identity.ErrUnauthenticated
	}
	p := identity.Principal{Subject: "synthetic", TenantID: "tenant", Permissions: map[string]struct{}{}}
	if token == "verified" {
		p.Permissions = map[string]struct{}{"social:read": {}, "social:request": {}}
	}
	return p, nil
}

type workspaceHTTPService struct {
	socialTestService
	calls int
}

func (s *workspaceHTTPService) Workspace(context.Context, identity.Principal, string) (postgres.SocialWorkspacePage, error) {
	s.calls++
	return postgres.SocialWorkspacePage{OrganizationID: "store", Items: []postgres.SocialWorkspaceRow{}}, nil
}
func (s *workspaceHTTPService) WorkspaceStatus(context.Context, identity.Principal, string) (postgres.SocialWorkspaceDetail, error) {
	s.calls++
	return postgres.SocialWorkspaceDetail{OrganizationID: "store"}, nil
}
func (s *workspaceHTTPService) PrepareDraft(context.Context, identity.Principal, postgres.SocialDraft) (postgres.SocialPrepared, error) {
	s.calls++
	return postgres.SocialPrepared{}, nil
}
func TestSocialWorkspaceHTTPBoundaries(t *testing.T) {
	s := &workspaceHTTPService{}
	server := httptest.NewServer(NewSocialPublishing(s, workspaceHTTPVerifier{}))
	defer server.Close()
	call := func(method, path, token, body string) int {
		t.Helper()
		r, e := http.NewRequest(method, server.URL+path, strings.NewReader(body))
		if e != nil {
			t.Fatal(e)
		}
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("Content-Type", "application/json")
		res, e := server.Client().Do(r)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		io.Copy(io.Discard, res.Body)
		if res.Header.Get("Cache-Control") != "no-store" {
			t.Fatal("cache boundary")
		}
		return res.StatusCode
	}
	for _, tc := range []struct {
		method, path, token, body string
		want                      int
	}{{"GET", "/v1/social/workspace", "bad", "", 401}, {"GET", "/v1/social/workspace", "reader-without-grant", "", 403}, {"GET", "/v1/social/workspace?tenant_id=other", "verified", "", 400}, {"GET", "/v1/social/workspace?after=a&after=b", "verified", "", 400}, {"POST", "/v1/social/workspace/prepare", "verified", `{"approval_id":"a","Approval_ID":"b"}`, 400}, {"POST", "/v1/social/workspace/prepare", "verified", `{"page_id":"999"}`, 400}} {
		if got := call(tc.method, tc.path, tc.token, tc.body); got != tc.want {
			t.Fatalf("%s got %d want %d", tc.path, got, tc.want)
		}
	}
	if s.calls != 0 {
		t.Fatal("invalid request reached workspace")
	}
	if call("GET", "/v1/social/workspace", "verified", "") != 200 || call("GET", "/v1/social/workspace/publication-00001", "verified", "") != 200 {
		t.Fatal("mounted read routes")
	}
	if call("POST", "/v1/social/workspace/prepare", "verified", `{"approval_id":"publication-00001","operation":"publish","message":"Fixture","original_approval_id":"","scheduled_at":"2026-09-14T12:00:00Z","expires_at":"2026-09-15T12:00:00Z"}`) != 200 {
		t.Fatal("bounded prepare route")
	}
	if s.calls != 3 {
		t.Fatal("expected three exact requests")
	}
}
