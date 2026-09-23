package whatsappbridge

// AUTHORED reference proof: original CRM fixture, HTTP identity, durable jobs,
// unchanged Meta adapter on loopback, current approval/fence/status/recovery.
import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestScheduledWhatsAppConnected(t *testing.T) {
	ctx := context.Background()
	database := os.Getenv("ELITE_WHATSAPP_CONNECTED_DATABASE_URL")
	u, e := url.Parse(database)
	if e != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_whatsapp_connected_") {
		t.Fatal("owned fixture required")
	}
	pool, e := pgxpool.New(ctx, database)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	maker, old, base := approvalFixtureData(t, pool)
	tenant := maker.TenantID
	if _, e = pool.Exec(ctx, `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref,organization_id)values($1,'wa-primary','meta-whatsapp','META_APP_SECRET','store-1')`, tenant); e != nil {
		t.Fatal(e)
	}
	permissions := map[string]struct{}{"notification:request": {}, "notification:approve": {}, "notification:read": {}, "notification:cancel": {}, "notification:dispatch": {}, "notification:reconcile": {}, "appointment:manage": {}, "whatsapp:process": {}}
	maker.Permissions = permissions
	reviewer := maker
	reviewer.Subject = "distinct-reviewer"
	foreign := maker
	foreign.Organizations = map[string]struct{}{"other": {}}
	key := []byte("0123456789abcdef0123456789abcdef")
	python := os.Getenv("ELITE_WHATSAPP_PYTHON")
	if python == "" {
		t.Fatal("fixed Python required")
	}
	actual, err := filepath.Abs("../../whatsapp_cloud")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := exec.Command(python, "-I", "-B", "-c", `import sys,json;sys.path.insert(0,sys.argv[1]);from test_whatsapp_cloud import profile;print(json.dumps(profile()))`, actual).Output()
	if err != nil {
		t.Fatal("profile fixture", err)
	}
	var profile json.RawMessage = raw
	profile = append(json.RawMessage(nil), []byte(strings.TrimSpace(string(profile)))...)
	var meta atomic.Int32
	metaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/messages" || r.Header.Get("Authorization") != "Bearer synthetic-token" {
			t.Error("unexpected fixture Meta transport")
		}
		body, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(body, &m)
		if m["type"] != "template" || m["to"] != "5491112345678" || m["messaging_product"] != "whatsapp" {
			t.Errorf("provider payload: %s", body)
		}
		n := meta.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"messaging_product":"whatsapp","messages":[{"id":"wamid.scheduled.%d"}]}`, n)
	}))
	t.Cleanup(metaServer.Close)
	// The production owner runs unchanged inside a hash-locked fixture wrapper.
	// Its only transport maps the exact official URL to this loopback HTTP server.
	// No live network/credential is used; lost stdout happens AFTER real receipt write.
	wrap := t.TempDir()
	quote := func(s string) string { b, _ := json.Marshal(s); return string(b) }
	actualCode, _ := os.ReadFile(filepath.Join(actual, "whatsapp_cloud.py"))
	statusCode, _ := os.ReadFile(filepath.Join(actual, "status_reconciliation.py"))
	script := `import sys,json,hashlib,pathlib,importlib.util,urllib.request
root=pathlib.Path(` + quote(actual) + `)
assert hashlib.sha256((root/'whatsapp_cloud.py').read_bytes()).hexdigest()==` + quote(digest(actualCode)) + `
assert hashlib.sha256((root/'status_reconciliation.py').read_bytes()).hexdigest()==` + quote(digest(statusCode)) + `
spec=importlib.util.spec_from_file_location('whatsapp_cloud',root/'whatsapp_cloud.py');m=importlib.util.module_from_spec(spec);sys.modules['whatsapp_cloud']=m;spec.loader.exec_module(m)
raw=sys.stdin.buffer.read(13*1024*1024)
def transport(method,url,headers,body,timeout):
 assert method=='POST' and url=='https://graph.facebook.com/v99.0/123456789/messages'
 req=urllib.request.Request(` + quote(metaServer.URL+"/messages") + `,data=body,headers=headers,method=method)
 with urllib.request.urlopen(req,timeout=timeout) as r:return r.status,dict(r.headers),r.read(65537)
mode=sys.argv[1]
if mode=='--send-bridge':
 reply=m.send_bridge(raw,transport)
 if json.loads(raw)['request'].get('body_parameters',[])[-1:] == ['lost-stdout']:sys.exit(2)
elif mode=='--ingress-bridge':reply=m.ingress_bridge(raw)
elif mode=='--recover-send-bridge':reply=m.recover_send_bridge(raw)
else:
 spec=importlib.util.spec_from_file_location('status_reconciliation',root/'status_reconciliation.py');q=importlib.util.module_from_spec(spec);spec.loader.exec_module(q);reply=q.status_bridge(raw)
print(json.dumps(reply))
`
	if err = os.WriteFile(filepath.Join(wrap, "whatsapp_cloud.py"), []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(wrap, "status_reconciliation.py"), statusCode, 0600); err != nil {
		t.Fatal(err)
	}
	pythonCode, _ := os.ReadFile(python)
	process := Process{PythonExecutable: python, PythonSHA256: digest(pythonCode), AdapterDirectory: wrap, AdapterSHA256: digest([]byte(script)), EvidenceDirectory: t.TempDir()}

	approvals, err := NewScheduleApprovals(base, tenant, "store-1", "wa-primary", profile)
	if err != nil {
		t.Fatal(err)
	}
	if e = approvals.VerifyInfrastructure(ctx); e != nil {
		t.Fatal(e)
	}
	store, err := postgres.NewOutboundDeliveryStore(pool, key, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	sender := &Sender{TenantID: tenant, Profile: profile, Process: process, Tokens: tokenFixture{}}
	module, err := NewScheduledNotifications(approvals, sender, store, VerifiedInboxChannel{}, "scheduled-fixture")
	if err != nil {
		t.Fatal(err)
	}
	verifier, token := notificationIssuer(t)
	mux := http.NewServeMux()
	module.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	call := func(method, path string, p identity.Principal, in any) (int, []byte) {
		t.Helper()
		var body io.Reader
		if in != nil {
			raw, _ := json.Marshal(in)
			body = bytes.NewReader(raw)
		}
		req, _ := http.NewRequest(method, api.URL+"/v1/franchise/notifications/scheduled"+path, body)
		req.Header.Set("Authorization", "Bearer "+token(p))
		if in != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		res, e := http.DefaultClient.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		raw, _ := io.ReadAll(res.Body)
		return res.StatusCode, raw
	}
	makeRequest := func(id, last string) (ScheduleContext, string) {
		t.Helper()
		r := ScheduleRequest{RequestID: id, AppointmentID: "appointment-1", Recipient: old.Message.ExternalID, TemplateName: "order_update", LanguageCode: "es_AR", BodyParameters: []string{"appointment-1", last}, NotBefore: time.Now().UTC().Add(2 * time.Second), ExpiresAt: time.Now().UTC().Add(time.Hour)}
		status, raw := call("POST", "/prepare", maker, r)
		if status != 200 {
			t.Fatalf("prepare %d %s", status, raw)
		}
		var c ScheduleContext
		if json.Unmarshal(raw, &c) != nil {
			t.Fatal("context")
		}
		if c.AppointmentVersion < 3 || c.ConsentID != "consent-1" || c.ProfileSHA256 != digest(profile) {
			t.Fatal("source not bound", c)
		}
		status, raw = call("POST", "/requests", maker, c)
		if status != 200 {
			t.Fatalf("submit %d %s", status, raw)
		}
		var reply struct {
			Hash string `json:"request_sha256"`
		}
		if json.Unmarshal(raw, &reply) != nil || !validDigest(reply.Hash) {
			t.Fatal("hash")
		}
		if code, _ := call("POST", "/requests/"+c.DeliveryKey+"/decision", maker, map[string]any{"request_sha256": reply.Hash, "approve": true, "reason": "self"}); code != 409 {
			t.Fatal("self approval", code)
		}
		status, raw = call("POST", "/requests/"+c.DeliveryKey+"/decision", reviewer, map[string]any{"request_sha256": reply.Hash, "approve": true, "reason": "synthetic scheduled reminder"})
		if status != 200 {
			t.Fatalf("decision %d %s", status, raw)
		}
		return c, reply.Hash
	}
	dispatch := func(want string) ScheduledWorkResult {
		t.Helper()
		out, e := module.ProcessOnce(ctx, maker)
		if e != nil || out.Outcome != want {
			t.Fatalf("dispatch %+v %v want%s", out, e, want)
		}
		return out
	}
	due := func(c ScheduleContext) {
		t.Helper()
		delay := time.Until(c.NotBefore) + 50*time.Millisecond
		if delay > 0 {
			time.Sleep(delay)
		}
	}
	if code, _ := call("POST", "/prepare", foreign, ScheduleRequest{}); code != 403 {
		t.Fatal("foreign scope")
	}
	first, h := makeRequest("scheduled-normal-0001", "normal")
	if out, e := module.ProcessOnce(ctx, maker); e != nil || out.Claimed || meta.Load() != 0 {
		t.Fatal("early effect", out, e)
	}
	due(first)
	var wg sync.WaitGroup
	results := make(chan ScheduledWorkResult, 12)
	errs := make(chan error, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); o, e := module.ProcessOnce(ctx, maker); results <- o; errs <- e }()
	}
	wg.Wait()
	close(results)
	close(errs)
	accepted := 0
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	for out := range results {
		if out.Outcome == "ACCEPTED" {
			accepted++
		}
	}
	if accepted != 1 || meta.Load() != 1 {
		t.Fatal("duplicate/absent provider effect", accepted, meta.Load())
	}
	if code, _ := call("POST", "/requests/"+first.DeliveryKey+"/cancel", maker, map[string]string{"request_sha256": h, "reason": "too late"}); code != 409 {
		t.Fatal("accepted effect recalled")
	}
	// Signed provider status uses the original observer and original status owner.
	receiptPath := filepath.Join(process.EvidenceDirectory, digest([]byte(tenant+"\x00"+first.DeliveryKey)), "SEND_RECEIPT.json")
	receipt, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatal(err)
	}
	statusBody, _ := json.Marshal(map[string]any{"object": "whatsapp_business_account", "entry": []any{map[string]any{"id": "987654321", "changes": []any{map[string]any{"field": "messages", "value": map[string]any{"messaging_product": "whatsapp", "metadata": map[string]string{"phone_number_id": "123456789"}, "statuses": []any{map[string]any{"id": "wamid.scheduled.1", "status": "delivered", "timestamp": fmt.Sprint(time.Now().Unix()), "recipient_id": old.Message.ExternalID}}}}}}}})
	sig := hmac.New(sha256.New, []byte("test-app-secret"))
	sig.Write(statusBody)
	observer := &StatusObserver{TenantID: tenant, Profile: profile, Process: process, ReconcilerSHA256: digest(statusCode), Secrets: appSecretFixture{}, Approvals: base}
	count, e := observer.Observe(ctx, maker, "store-1", "", first.DeliveryKey, receipt, []SignedStatusWebhook{{Body: statusBody, Signature: "sha256=" + hex.EncodeToString(sig.Sum(nil))}})
	if e != nil || count != 1 {
		t.Fatal("signed status", count, e)
	}
	visible, e := approvals.Status(ctx, maker, first.DeliveryKey)
	if e != nil || visible.Notification == nil || visible.Notification.DeliveryStatus != "observed_delivered" || !visible.JobCompleted {
		t.Fatal("durable status", visible, e)
	}
	cancelled, ch := makeRequest("scheduled-cancel-0001", "cancel")
	if code, _ := call("POST", "/requests/"+cancelled.DeliveryKey+"/cancel", maker, map[string]string{"request_sha256": ch, "reason": "operator cancellation"}); code != 200 {
		t.Fatal("cancel", code)
	}
	due(cancelled)
	dispatch("SUPPRESSED")
	if meta.Load() != 1 {
		t.Fatal("cancelled send")
	}
	changed, _ := makeRequest("scheduled-source-0001", "source-change")
	// Simulate a separately committed original CRM reschedule; preserve the
	// historical reviewed reminder and require a new source/version approval.
	if _, e = pool.Exec(ctx, `update crm.appointment set version=version+1,starts_at=starts_at+interval '1 hour',updated_at=clock_timestamp() where tenant_id=$1 and appointment_id='appointment-1'`, tenant); e != nil {
		t.Fatal(e)
	}
	due(changed)
	dispatch("SUPPRESSED")
	replacement, _ := makeRequest("scheduled-newsource-01", "rescheduled")
	due(replacement)
	dispatch("ACCEPTED")
	if meta.Load() != 2 {
		t.Fatal("replacement")
	}
	lost, lh := makeRequest("scheduled-recovery-01", "lost-stdout")
	due(lost)
	dispatch("RECONCILIATION_REQUIRED")
	if meta.Load() != 3 {
		t.Fatal("lost result effect")
	}
	time.Sleep(1100 * time.Millisecond)
	if code, raw := call("POST", "/requests/"+lost.DeliveryKey+"/reconcile", maker, map[string]string{"request_sha256": lh}); code != 200 {
		t.Fatalf("recover %d %s", code, raw)
	}
	if meta.Load() != 3 {
		t.Fatal("recovery resent")
	}
	exhausted, _ := makeRequest("scheduled-exhausted-1", "exhausted")
	due(exhausted)
	if _, e = pool.Exec(ctx, `update platform.job j set attempts=max_attempts,claimed_by='crashed-fixture',claimed_until=clock_timestamp()-interval '1 second' from communication.whatsapp_schedule s where s.tenant_id=j.tenant_id and s.job_id=j.job_id and s.tenant_id=$1 and s.delivery_key=$2`, tenant, exhausted.DeliveryKey); e != nil {
		t.Fatal(e)
	}
	dispatch("LEASE_EXHAUSTED")
	visible, e = approvals.Status(ctx, maker, exhausted.DeliveryKey)
	if e != nil || visible.JobCompleted || visible.TerminalError != "LEASE_EXHAUSTED" || meta.Load() != 3 {
		t.Fatal("exhaustion hidden", visible, e)
	}
	withdrawn, _ := makeRequest("scheduled-withdrawn1", "withdrawn")
	if _, e = pool.Exec(ctx, `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)values($1,'consent-withdrawn','lead-1','fixture-appointment-notification','policy-1','withdrawn',clock_timestamp(),repeat('f',64))`, tenant); e != nil {
		t.Fatal(e)
	}
	due(withdrawn)
	dispatch("SUPPRESSED")
	if meta.Load() != 3 {
		t.Fatal("opt-out ignored")
	}
	for _, table := range []string{"communication.whatsapp_schedule", "communication.whatsapp_schedule_cancellation", "communication.whatsapp_schedule_result"} {
		if _, e = pool.Exec(ctx, "delete from "+table+" where tenant_id=$1", tenant); e == nil {
			t.Fatal("mutable evidence", table)
		}
	}
	var approvalsN, jobsN, fencesN, resultsN int
	if e = pool.QueryRow(ctx, `select (select count(*) from approval.request where tenant_id=$1 and kind='whatsapp_schedule'),(select count(*) from platform.job where tenant_id=$1 and queue='whatsapp-scheduled'),(select count(*) from communication.outbound_delivery where tenant_id=$1 and channel_code='whatsapp'),(select count(*) from communication.whatsapp_schedule_result where tenant_id=$1)`, tenant).Scan(&approvalsN, &jobsN, &fencesN, &resultsN); e != nil {
		t.Fatal(e)
	}
	if approvalsN != 7 || jobsN != 7 || fencesN != 3 || resultsN != 6 {
		t.Fatal("counts", approvalsN, jobsN, fencesN, resultsN)
	}
	t.Logf("SCHEDULED_COMMUNICATIONS_PASS approvals=%d jobs=%d fences=%d results=%d provider_posts=%d concurrent=12 accepted=1 signed_delivered=1 no_cancel_reschedule_optout_replay tenant=%s", approvalsN, jobsN, fencesN, resultsN, meta.Load(), tenant)
}
