package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/socialbridge"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sdkStep struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Value   any    `json:"value,omitempty"`
	Timeout bool   `json:"timeout,omitempty"`
}

func sdkProcess(t *testing.T, steps []sdkStep) (socialbridge.Process, func() []map[string]any) {
	t.Helper()
	python := os.Getenv("SOCIAL_TEST_PYTHON")
	if python == "" {
		t.Fatal("SOCIAL_TEST_PYTHON exact admitted runtime required")
	}
	root, e := filepath.Abs("../../..")
	if e != nil {
		t.Fatal(e)
	}
	d := t.TempDir()
	state := filepath.Join(d, "steps.json")
	log := filepath.Join(d, "calls.jsonl")
	script := filepath.Join(d, "fixture.py")
	raw, _ := json.Marshal(steps)
	if os.WriteFile(state, raw, 0600) != nil {
		t.Fatal("fixture state")
	}
	quote := func(v string) string { b, _ := json.Marshal(v); return string(b) }
	code := `# Test-only real SDK transport; no production flag enables this transport.
import sys,json
from pathlib import Path
sys.path.insert(0, ` + quote(root) + `)
from meta_page_write.bridge import main
from meta_page_write.adapter import create_api
from requests.adapters import BaseAdapter
from requests import Response, Timeout
from urllib.parse import urlparse,parse_qs
state=Path(` + quote(state) + `)
log=Path(` + quote(log) + `)
class Fixture(BaseAdapter):
 def send(self,request,**kwargs):
  url=urlparse(request.url)
  assert url.scheme=='https' and url.netloc=='graph.facebook.com'
  assert kwargs['timeout']==10
  steps=json.loads(state.read_text()); assert steps
  step=steps.pop(0)
  assert request.method==step['method'] and url.path.rstrip('/')==step['path'],(request.method,url.path,step)
  state.write_text(json.dumps(steps))
  body=request.body.decode() if isinstance(request.body,bytes) else request.body
  with log.open('a') as f:f.write(json.dumps({'method':request.method,'path':url.path,'body':parse_qs(body or '')})+'\n')
  if step.get('timeout'):raise Timeout('fixture')
  response=Response();response.status_code=200;response.url=request.url;response.request=request;response._content=json.dumps(step.get('value',{})).encode();response.headers['Content-Type']='application/json';return response
 def close(self):pass
def api():
 value=create_api('321','fixture-secret','fixture-token');value._session.requests.mount('https://',Fixture());value._session.requests.mount('http://',Fixture());return value
main(api)
`
	if e = os.WriteFile(script, []byte(code), 0600); e != nil {
		t.Fatal(e)
	}
	read := func() []map[string]any {
		raw, e := os.ReadFile(log)
		if os.IsNotExist(e) {
			return nil
		}
		if e != nil {
			t.Fatal(e)
		}
		var values []map[string]any
		for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
			var v map[string]any
			if json.Unmarshal([]byte(line), &v) != nil {
				t.Fatal("fixture log invalid")
			}
			values = append(values, v)
		}
		return values
	}
	return socialbridge.Process{Python: python, Script: script, ScriptSHA256: socialbridge.Hash([]byte(code))}, read
}
func TestSocialPublishingConnectedSDKApprovalAndRecovery(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("dedicated PostgreSQL required")
	}
	cfg, e := pgxpool.ParseConfig(url)
	if e != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_social_") {
		t.Fatal("isolated local database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	pool, e := pgxpool.NewWithConfig(ctx, cfg)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant, org := uuid.NewString(), uuid.NewString()
	_, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Synthetic Social','Fixture')`, tenant, "social-"+uuid.NewString())
	if e != nil {
		t.Fatal(e)
	}
	config := socialbridge.Config{Schema: "elite.meta-page-publishing.v1", TenantID: tenant, OrganizationID: org, PageID: "123", Queue: "social_fixture", Publish: true, Revoke: true, Review: "one_distinct_human", LeaseSeconds: 60, RetrySeconds: 1, PollSeconds: 1, MaxAttempts: 2}
	raw, _ := json.Marshal(config)
	profile, e := socialbridge.Load(raw, socialbridge.Hash(raw))
	if e != nil {
		t.Fatal(e)
	}
	actor := func(subject string) identity.Principal {
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"social:request": {}, "social:approve": {}, "social:read": {}, "social:reconcile": {}}, Organizations: map[string]struct{}{org: {}}}
	}
	requester, reviewer := actor("requester"), actor("reviewer")
	makeRequest := func(id string) socialbridge.Request {
		i := socialbridge.Intent{TenantID: tenant, PageID: "123", ProfileSHA256: profile.SHA256(), ApprovalID: id, Operation: "publish", Message: "Exact reviewed fixture bytes"}
		i.ContentSHA256 = socialbridge.Hash([]byte(i.Message))
		i.DeliveryKey = socialbridge.DeliveryKey(i)
		return socialbridge.Request{Intent: i, ScheduledAt: time.Now().UTC().Add(-time.Second), ExpiresAt: time.Now().UTC().Add(time.Hour)}
	}
	observed := func(ref string) any {
		return map[string]any{"id": ref, "from": map[string]string{"id": "123"}, "message": "Exact reviewed fixture bytes", "is_published": true}
	}
	process, calls := sdkProcess(t, []sdkStep{
		{Method: "POST", Path: "/v26.0/123/feed", Value: map[string]string{"id": "123_100"}}, {Method: "GET", Path: "/v26.0/123_100", Value: observed("123_100")},
		{Method: "POST", Path: "/v26.0/123/feed", Value: map[string]string{"id": "123_101"}}, {Method: "GET", Path: "/v26.0/123_101", Timeout: true}, {Method: "GET", Path: "/v26.0/123_101", Value: observed("123_101")},
		{Method: "GET", Path: "/v26.0/123_100", Value: observed("123_100")}, {Method: "DELETE", Path: "/v26.0/123_100", Value: map[string]bool{"success": true}},
		{Method: "POST", Path: "/v26.0/123/feed", Timeout: true},
		{Method: "POST", Path: "/v26.0/123/feed", Value: map[string]string{"id": "123_102"}}, {Method: "GET", Path: "/v26.0/123_102", Timeout: true}, {Method: "GET", Path: "/v26.0/123_102", Timeout: true}, {Method: "GET", Path: "/v26.0/123_102", Value: observed("123_102")},
	})
	service, e := NewSocialPublishing(pool, profile, []byte(strings.Repeat("k", 32)), process)
	if e != nil {
		t.Fatal(e)
	}
	claim := func() Job {
		t.Helper()
		v, e := service.Claim(ctx, "worker")
		if e != nil || len(v) != 1 {
			t.Fatalf("claim %+v %v", v, e)
		}
		return v[0]
	}
	approve := func(r socialbridge.Request) string {
		t.Helper()
		hash, replay, e := service.Submit(ctx, requester, r)
		if e != nil || replay {
			t.Fatalf("submit %v %v", replay, e)
		}
		if _, e = service.Decide(ctx, reviewer, r.Intent.ApprovalID, hash, true, "Reviewed exact bytes"); e != nil {
			t.Fatal(e)
		}
		return hash
	}
	r := makeRequest("published")
	hash, replay, e := service.Submit(ctx, requester, r)
	if e != nil || replay {
		t.Fatal(e)
	}
	if jobs, e := service.Claim(ctx, "worker"); e != nil || len(jobs) != 0 {
		t.Fatal("unapproved request enqueued")
	}
	if _, e = service.Decide(ctx, requester, r.Intent.ApprovalID, hash, true, ""); !errors.Is(e, approval.ErrSeparation) {
		t.Fatalf("self review %v", e)
	}
	if _, e = service.Decide(ctx, reviewer, r.Intent.ApprovalID, strings.Repeat("0", 64), true, ""); e == nil {
		t.Fatal("wrong hash accepted")
	}
	_, replay, e = service.Submit(ctx, requester, r)
	if e != nil || !replay {
		t.Fatal("exact submit replay rejected")
	}
	changed := r
	changed.Intent.Message = "other"
	changed.Intent.ContentSHA256 = socialbridge.Hash([]byte("other"))
	if _, _, e = service.Submit(ctx, requester, changed); e == nil {
		t.Fatal("changed content under same id")
	}
	foreign := requester
	foreign.TenantID = uuid.NewString()
	if _, _, e = service.Submit(ctx, foreign, r); e == nil {
		t.Fatal("foreign tenant")
	}
	badPage := r
	badPage.Intent.PageID = "999"
	badPage.Intent.DeliveryKey = socialbridge.DeliveryKey(badPage.Intent)
	if _, _, e = service.Submit(ctx, requester, badPage); e == nil {
		t.Fatal("foreign page")
	}
	if _, e = service.Decide(ctx, reviewer, r.Intent.ApprovalID, hash, true, "Exact reviewed bytes"); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `update approval.request set payload='{}' where tenant_id=$1 and request_id=$2`, tenant, r.Intent.ApprovalID); e == nil {
		t.Fatal("approval payload mutable")
	}
	if _, e = pool.Exec(ctx, `update approval.decision set reviewer='intruder' where tenant_id=$1 and request_id=$2`, tenant, r.Intent.ApprovalID); e == nil {
		t.Fatal("decision mutable")
	}
	job := claim()
	if jobs, e := service.Claim(ctx, "other-worker"); e != nil || len(jobs) != 0 {
		t.Fatal("concurrent claim won")
	}
	stale := job
	stale.Attempts++
	if e = service.Process(ctx, stale, "worker"); !errors.Is(e, ErrJobClaimLost) {
		t.Fatalf("wrong generation %v", e)
	}
	if len(calls()) != 0 {
		t.Fatal("stale generation sent")
	}
	if e = service.Process(ctx, job, "worker"); e != nil {
		t.Fatal(e)
	}
	if len(calls()) != 2 {
		t.Fatal("SDK POST/GET not observed")
	}
	_ = service.Process(ctx, job, "worker")
	if len(calls()) != 2 {
		t.Fatal("replay repeated POST")
	}
	status, e := service.Status(ctx, reviewer, r.Intent.ApprovalID)
	if e != nil || status.DeliveryState != "accepted" || len(status.Receipt) == 0 {
		t.Fatalf("status %+v %v", status, e)
	}
	if _, e = service.Status(ctx, foreign, r.Intent.ApprovalID); e == nil {
		t.Fatal("foreign read")
	}
	uncertain := makeRequest("unknown-known-id")
	approve(uncertain)
	first := claim()
	if e = service.Process(ctx, first, "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatalf("expected unknown %v", e)
	}
	if len(calls()) != 4 {
		t.Fatal("unknown POST/GET count")
	}
	// A new instance proves receipt identity and retry behavior are in PostgreSQL.
	restarted, e := NewSocialPublishing(pool, profile, []byte(strings.Repeat("k", 32)), process)
	if e != nil {
		t.Fatal(e)
	}
	service = restarted
	if _, e = pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, first.JobID); e != nil {
		t.Fatal(e)
	}
	if e = service.Process(ctx, claim(), "worker"); e != nil {
		t.Fatal(e)
	}
	if len(calls()) != 5 || calls()[4]["method"] != "GET" {
		t.Fatal("recovery repeated POST")
	}
	revoke := makeRequest("revoke-reviewed")
	revoke.Intent.Operation = "revoke"
	revoke.Intent.DeliveryKey = socialbridge.DeliveryKey(revoke.Intent)
	revoke.OriginalApprovalID = r.Intent.ApprovalID
	revoke.ProviderReference = "123_999"
	if _, _, e = service.Submit(ctx, requester, revoke); e == nil {
		t.Fatal("arbitrary delete reference")
	}
	revoke.ProviderReference = "123_100"
	approve(revoke)
	if e = service.Process(ctx, claim(), "worker"); e != nil {
		t.Fatal(e)
	}
	if len(calls()) != 7 || calls()[6]["method"] != "DELETE" {
		t.Fatal("governed revoke absent")
	}
	unknown := makeRequest("unknown-no-id")
	approve(unknown)
	lost := claim()
	if e = service.Process(ctx, lost, "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatalf("POST uncertainty %v", e)
	}
	if _, e = pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, lost.JobID); e != nil {
		t.Fatal(e)
	}
	if e = service.Process(ctx, claim(), "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatal("missing reference guessed")
	}
	if e = service.Reconcile(ctx, reviewer, unknown.Intent.ApprovalID); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatal("missing ref reconciled")
	}
	if len(calls()) != 8 {
		t.Fatal("unknown write duplicated")
	}
	terminal := makeRequest("terminal-known-id")
	approve(terminal)
	attempt := claim()
	if e = service.Process(ctx, attempt, "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, attempt.JobID); e != nil {
		t.Fatal(e)
	}
	if e = service.Process(ctx, claim(), "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatal(e)
	}
	var attempts int
	var terminalCode string
	if e = pool.QueryRow(ctx, `select attempts,terminal_error_code from platform.job where tenant_id=$1 and job_id=$2`, tenant, attempt.JobID).Scan(&attempts, &terminalCode); e != nil || attempts != 2 || terminalCode == "" {
		t.Fatalf("terminal budget %d %s %v", attempts, terminalCode, e)
	}
	if e = service.Reconcile(ctx, reviewer, terminal.Intent.ApprovalID); e != nil {
		t.Fatal("terminal GET recovery", e)
	}
	var complete bool
	if e = pool.QueryRow(ctx, `select attempts=2 and completed_at is not null and terminal_error_code is null from platform.job where tenant_id=$1 and job_id=$2`, tenant, attempt.JobID).Scan(&complete); e != nil || !complete {
		t.Fatal("reconciliation reset generation or failed ack", e)
	}
	if len(calls()) != 12 || calls()[11]["method"] != "GET" {
		t.Fatal("terminal reconciliation wrote provider")
	}
	crash := makeRequest("crash-after-durable-receipt")
	approve(crash)
	crashedJob := claim()
	m := crash.Message()
	mh, _ := outbounddelivery.MessageSHA256(m)
	if _, e = service.fence.Claim(ctx, m, mh); e != nil {
		t.Fatal(e)
	}
	crashReceipt := socialbridge.Receipt{TenantID: tenant, PageID: "123", ProfileSHA256: profile.SHA256(), ApprovalID: crash.Intent.ApprovalID, DeliveryKey: crash.Intent.DeliveryKey, ProviderReference: "123_103", State: "published", ContentSHA256: crash.Intent.ContentSHA256, EvidenceSHA256: strings.Repeat("a", 64)}
	if e = service.event(ctx, crash, "social.receipt", crashReceipt); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `update communication.outbound_delivery set locked_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and channel_code='facebook_pages' and delivery_key=$2`, tenant, crash.Intent.DeliveryKey); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, crashedJob.JobID); e != nil {
		t.Fatal(e)
	}
	if e = service.Reconcile(ctx, reviewer, crash.Intent.ApprovalID); e != nil {
		t.Fatal("durable receipt crash recovery", e)
	}
	if len(calls()) != 12 {
		t.Fatal("durable receipt caused provider call")
	}
	expiring := makeRequest("expiry-while-lock-wait")
	expiring.ExpiresAt = time.Now().UTC().Add(500 * time.Millisecond)
	approve(expiring)
	expiryJob := claim()
	block, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = block.Exec(ctx, `select tenant_id from platform.tenant where tenant_id=$1 for update`, tenant); e != nil {
		t.Fatal(e)
	}
	done := make(chan error, 1)
	go func() { done <- service.Process(ctx, expiryJob, "worker") }()
	time.Sleep(time.Until(expiring.ExpiresAt) + 100*time.Millisecond)
	if e = block.Commit(ctx); e != nil {
		t.Fatal(e)
	}
	if e = <-done; e == nil {
		t.Fatal("expiry during authority lock wait admitted")
	}
	if len(calls()) != 12 {
		t.Fatal("expired approval sent")
	}
	// This rejected claim is intentionally released without changing its request.
	_, _ = service.jobs.Fail(ctx, tenant, expiryJob.JobID, "worker", "APPROVAL_EXPIRED", expiryJob.Attempts, time.Hour)
	scheduled := makeRequest("scheduled-future")
	scheduled.ScheduledAt = time.Now().UTC().Add(30 * time.Minute)
	approve(scheduled)
	if jobs, e := service.Claim(ctx, "worker"); e != nil || len(jobs) != 0 {
		t.Fatal("scheduled early")
	}
	// Tampering the queue timestamp is still constrained by the approved schedule.
	if _, e = pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, socialJobID(tenant, scheduled.Intent.ApprovalID)); e != nil {
		t.Fatal(e)
	}
	if e = service.Process(ctx, claim(), "worker"); e == nil {
		t.Fatal("approved not-before bypassed")
	}
	if len(calls()) != 12 {
		t.Fatal("early schedule effect")
	}
	t.Logf("real official SDK: POST=4 GET=7 DELETE=1; known-reference/terminal recovery GET-only; receipt-before-fence crash and expiry-under-lock covered; no network; tenant=%s", tenant)
	_ = fmt.Sprintf("%s", hash)
}
