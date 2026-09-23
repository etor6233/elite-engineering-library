package postgres

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/socialbridge"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"testing"
	"time"
)

type publishingFixture struct {
	service             *SocialPublishing
	pool                *pgxpool.Pool
	requester, reviewer identity.Principal
	profile             socialbridge.Config
	calls               func() []map[string]any
}

func publishingFixtureFor(t *testing.T, steps []sdkStep) *publishingFixture {
	t.Helper()
	raw := os.Getenv("TEST_DATABASE_URL")
	cfg, e := pgxpool.ParseConfig(raw)
	if e != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_social_") {
		t.Fatal("owned social PG required")
	}
	pool, e := pgxpool.NewWithConfig(context.Background(), cfg)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(pool.Close)
	tenant, org := uuid.NewString(), "store"
	_, e = pool.Exec(context.Background(), `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Synthetic','Synthetic')`, tenant, "social-"+uuid.NewString())
	if e != nil {
		t.Fatal(e)
	}
	_, e = pool.Exec(context.Background(), `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,$2,'store','Synthetic','store')`, tenant, org)
	if e != nil {
		t.Fatal(e)
	}
	c := socialbridge.Config{Schema: "elite.meta-page-publishing.v1", TenantID: tenant, OrganizationID: org, PageID: "123", Queue: "social_fixture", Publish: true, Revoke: true, Review: "one_distinct_human", LeaseSeconds: 60, RetrySeconds: 1, PollSeconds: 1, MaxAttempts: 2}
	b, _ := json.Marshal(c)
	p, e := socialbridge.Load(b, socialbridge.Hash(b))
	if e != nil {
		t.Fatal(e)
	}
	process, calls := sdkProcess(t, steps)
	s, e := NewSocialPublishing(pool, p, []byte(strings.Repeat("k", 32)), process)
	if e != nil {
		t.Fatal(e)
	}
	actor := func(subject string) identity.Principal {
		return identity.Principal{TenantID: tenant, Subject: subject, Organizations: map[string]struct{}{org: {}}, Permissions: map[string]struct{}{"social:read": {}, "social:request": {}, "social:approve": {}, "social:reconcile": {}}}
	}
	return &publishingFixture{s, pool, actor("requester"), actor("reviewer"), c, calls}
}
func (f *publishingFixture) request(id string) socialbridge.Request {
	i := socialbridge.Intent{TenantID: f.profile.TenantID, PageID: f.profile.PageID, ProfileSHA256: f.service.profile.SHA256(), ApprovalID: id, Operation: "publish", Message: "Exact reviewed fixture bytes"}
	i.ContentSHA256 = socialbridge.Hash([]byte(i.Message))
	i.DeliveryKey = socialbridge.DeliveryKey(i)
	return socialbridge.Request{Intent: i, ScheduledAt: time.Now().UTC().Add(-time.Second), ExpiresAt: time.Now().UTC().Add(time.Hour)}
}
func TestPublishingWorkspaceStatusRejectsOtherProfile(t *testing.T) {
	f := publishingFixtureFor(t, nil)
	ctx := context.Background()
	r := f.request(uuid.NewString())
	if _, _, e := f.service.Submit(ctx, f.requester, r); e != nil {
		t.Fatal(e)
	}
	c := f.profile
	c.PageID = "999"
	b, _ := json.Marshal(c)
	p, e := socialbridge.Load(b, socialbridge.Hash(b))
	if e != nil {
		t.Fatal(e)
	}
	other, e := NewSocialPublishing(f.pool, p, []byte(strings.Repeat("k", 32)), f.service.adapter)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = other.Status(ctx, f.requester, r.Intent.ApprovalID); e == nil {
		t.Fatal("CONFIRMED: old Status exposed a request from a different configured Page/profile")
	}
}
func TestPublishingWorkspacePrepareReviewWorkerReceiptRevoke(t *testing.T) {
	observed := map[string]any{"id": "123_200", "from": map[string]string{"id": "123"}, "message": "Exact reviewed fixture bytes", "is_published": true}
	f := publishingFixtureFor(t, []sdkStep{{Method: "POST", Path: "/v26.0/123/feed", Value: map[string]string{"id": "123_200"}}, {Method: "GET", Path: "/v26.0/123_200", Value: observed}, {Method: "GET", Path: "/v26.0/123_200", Value: observed}, {Method: "DELETE", Path: "/v26.0/123_200", Value: map[string]bool{"success": true}}})
	ctx := context.Background()
	now := time.Now().UTC()
	draft := SocialDraft{ID: uuid.NewString(), Operation: "publish", Message: "Exact reviewed fixture bytes", ScheduledAt: now.Add(-time.Second), ExpiresAt: now.Add(time.Hour)}
	bad := draft
	bad.ID = "short"
	if _, e := f.service.PrepareDraft(ctx, f.requester, bad); e == nil {
		t.Fatal("new workspace ID boundary")
	}
	empty, e := f.service.Workspace(ctx, f.requester, "")
	if e != nil || len(empty.Items) != 0 || empty.OrganizationID != f.profile.OrganizationID {
		t.Fatalf("empty %+v %v", empty, e)
	}
	prepared, e := f.service.PrepareDraft(ctx, f.requester, draft)
	if e != nil {
		t.Fatal(e)
	}
	if len(f.calls()) != 0 {
		t.Fatal("prepare contacted provider")
	}
	empty, e = f.service.Workspace(ctx, f.requester, "")
	if e != nil || len(empty.Items) != 0 {
		t.Fatal("prepare persisted approval")
	}
	h, replayed, e := f.service.Submit(ctx, f.requester, prepared.Request)
	if e != nil || replayed || h != prepared.Hash {
		t.Fatalf("exact prepared submit %v", e)
	}
	if _, replayed, e = f.service.Submit(ctx, f.requester, prepared.Request); e != nil || !replayed {
		t.Fatal("exact replay")
	}
	list, e := f.service.Workspace(ctx, f.requester, "")
	if e != nil || len(list.Items) != 1 || !list.Items[0].RequestedBySelf || list.Items[0].Message != draft.Message {
		t.Fatalf("list %+v %v", list, e)
	}
	detail, e := f.service.WorkspaceStatus(ctx, f.reviewer, draft.ID)
	if e != nil || detail.RequestedBySelf || len(detail.Status.Receipt) != 0 || detail.Status.DeliveryState != "" {
		t.Fatalf("pending is not published %+v %v", detail, e)
	}
	foreign := f.reviewer
	foreign.Organizations = map[string]struct{}{"other": {}}
	if _, e = f.service.Workspace(ctx, foreign, ""); e == nil {
		t.Fatal("foreign org list")
	}
	if _, e = f.service.WorkspaceStatus(ctx, foreign, draft.ID); e == nil {
		t.Fatal("foreign org detail")
	}
	if _, e = f.service.Decide(ctx, f.requester, draft.ID, h, true, ""); e == nil {
		t.Fatal("self review")
	}
	if _, e = f.service.Decide(ctx, f.reviewer, draft.ID, h, true, "Exact content reviewed"); e != nil {
		t.Fatal(e)
	}
	run := func() {
		t.Helper()
		jobs, e := f.service.Claim(ctx, "workspace-worker")
		if e != nil || len(jobs) != 1 {
			t.Fatalf("claim %d %v", len(jobs), e)
		}
		if e = f.service.Process(ctx, jobs[0], "workspace-worker"); e != nil {
			t.Fatal(e)
		}
	}
	run()
	detail, e = f.service.WorkspaceStatus(ctx, f.reviewer, draft.ID)
	if e != nil || detail.Status.DeliveryState != "accepted" {
		t.Fatalf("published %+v %v", detail, e)
	}
	var receipt socialbridge.Receipt
	if socialbridge.Decode(detail.Status.Receipt, &receipt) != nil || receipt.Validate(prepared.Request.Intent) != nil || receipt.State != "published" {
		t.Fatal("verified published receipt")
	}
	revoke, e := f.service.PrepareDraft(ctx, f.requester, SocialDraft{ID: uuid.NewString(), Operation: "revoke", OriginalID: draft.ID, ScheduledAt: now.Add(-time.Second), ExpiresAt: now.Add(time.Hour)})
	if e != nil {
		t.Fatal(e)
	}
	if revoke.Request.Intent.Message != draft.Message || revoke.Request.ProviderReference != "123_200" {
		t.Fatal("revoke not bound to receipt")
	}
	rh, _, e := f.service.Submit(ctx, f.requester, revoke.Request)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.service.Decide(ctx, f.reviewer, revoke.Request.Intent.ApprovalID, rh, true, "Remove exact reviewed publication"); e != nil {
		t.Fatal(e)
	}
	run()
	detail, e = f.service.WorkspaceStatus(ctx, f.reviewer, revoke.Request.Intent.ApprovalID)
	if e != nil || socialbridge.Decode(detail.Status.Receipt, &receipt) != nil || receipt.State != "revoked" || receipt.Validate(revoke.Request.Intent) != nil {
		t.Fatalf("revoke receipt %v", e)
	}
	if len(f.calls()) != 4 {
		t.Fatalf("SDK wire calls %d", len(f.calls()))
	}
	if e = f.service.Reconcile(ctx, f.reviewer, revoke.Request.Intent.ApprovalID); e != nil {
		t.Fatal(e)
	}
	if len(f.calls()) != 4 {
		t.Fatal("accepted reconcile repeated provider effect")
	}
	t.Log("PASS prepare (no persistence/send), exact submit/replay, distinct review, worker SDK POST+GET, verified receipt, receipt-bound revoke GET+DELETE, accepted reconcile no duplicate effect")
}
