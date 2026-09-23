package whatsappbridge

// AUTHORED reference proof: original CRM fixture, HTTP identity, durable jobs,
// unchanged Meta adapter on loopback, current approval/fence/status/recovery.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
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
	"sync/atomic"
	"testing"
	"time"
)

func TestCampaignPartialProgress(t *testing.T) {
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
	permissions := map[string]struct{}{"marketing:request": {}, "marketing:approve": {}, "marketing:read": {}, "marketing:cancel": {}, "notification:request": {}, "notification:approve": {}, "notification:read": {}, "notification:cancel": {}, "notification:dispatch": {}, "notification:reconcile": {}, "appointment:manage": {}, "whatsapp:process": {}}
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
		fmt.Fprintf(w, `{"messaging_product":"whatsapp","messages":[{"id":"wamid.campaign.%d"}]}`, n)
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

	campaigns, e := NewCampaigns(module, CampaignPolicy{Schema: "elite-whatsapp-campaign-policy/v1", ConsentPurpose: "fixture-marketing", PolicyVersion: "policy-1", MaxMembers: 20, MaxSteps: 3})
	if e != nil {
		t.Fatal(e)
	}
	if e = campaigns.VerifyInfrastructure(ctx); e != nil {
		t.Fatal(e)
	}
	seed := []string{
		`insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)values($1,'marketing-consent-1','lead-1','fixture-marketing','policy-1','granted',clock_timestamp()-interval '1 minute',repeat('a',64))`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Customer','customer@example.test')`,
		`update crm.lead set customer_principal_id='customer',model_id='model',updated_at=clock_timestamp() where tenant_id=$1 and lead_id='lead-1'`,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'retail','AR','ARS',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','variant',123456,'inclusive')`,
	}
	for _, q := range seed {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	verifier, token := notificationIssuer(t)
	mux := http.NewServeMux()
	module.Register(mux, verifier)
	campaigns.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	call := func(method, path string, p identity.Principal, in any) (int, []byte) {
		t.Helper()
		var body io.Reader
		if in != nil {
			raw, _ := json.Marshal(in)
			body = bytes.NewReader(raw)
		}
		req, _ := http.NewRequest(method, api.URL+"/v1/franchise/marketing/campaigns"+path, body)
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

	now := time.Now().UTC()
	request := CampaignRequest{ID: "campaign-partial-0001", Sources: []string{"fixture"}, States: []string{"new"}, Members: []CampaignRecipient{{LeadID: "lead-1", Recipient: old.Message.ExternalID}},
		Steps: []CampaignStep{
			{TemplateName: "order_update", LanguageCode: "es_AR", BodyParameters: []string{"partial", "first"}, NotBefore: now.Add(4 * time.Second), ExpiresAt: now.Add(5 * time.Second)},
			{TemplateName: "order_update", LanguageCode: "es_AR", BodyParameters: []string{"partial", "second"}, NotBefore: now.Add(7 * time.Second), ExpiresAt: now.Add(8 * time.Second)},
		}}
	if code, _ := call("POST", "/prepare", foreign, request); code != 403 {
		t.Fatal("foreign")
	}
	code, raw := call("POST", "/prepare", maker, request)
	var snapshot CampaignSnapshot
	if code != 200 || json.Unmarshal(raw, &snapshot) != nil {
		t.Fatalf("prepare %d %s", code, raw)
	}
	// A fixture trigger interrupts the second independent original approval.
	// The first request remains durable and the API must report partial state.
	fixture := `create function communication.campaign_fixture_fail_step() returns trigger language plpgsql as $f$
    begin
      if new.kind='whatsapp_schedule' and new.payload->'request'->>'campaign_id'='campaign-partial-0001' and new.payload->'request'->>'step'='2'
      then raise exception 'synthetic second-step persistence failure';end if;
      return new;
    end;$f$;
    create trigger campaign_fixture_fail_step before insert on approval.request for each row execute function communication.campaign_fixture_fail_step();`
	if _, e = pool.Exec(ctx, fixture); e != nil {
		t.Fatal(e)
	}
	restore := func() {
		_, _ = pool.Exec(context.Background(), "drop trigger if exists campaign_fixture_fail_step on approval.request;drop function if exists communication.campaign_fixture_fail_step()")
	}
	defer restore()
	code, raw = call("POST", "", maker, snapshot)
	var partial CampaignBatch
	if code != 200 || json.Unmarshal(raw, &partial) != nil || partial.Complete || len(partial.Items) != 2 || partial.Items[0].State != "PREPARED" || partial.Items[1].Code != "PREPARATION_UNCONFIRMED" {
		t.Fatalf("partial create %d %s", code, raw)
	}
	decision := map[string]any{"campaign_sha256": partial.CampaignSHA256, "approve": true, "reason": "explicit stored campaign batch"}
	code, raw = call("POST", "/"+request.ID+"/decision", reviewer, decision)
	var reviewed CampaignBatch
	if code != 200 || json.Unmarshal(raw, &reviewed) != nil || reviewed.Complete || reviewed.Items[0].State != "approved" || reviewed.Items[1].Code != "REQUEST_UNAVAILABLE" {
		t.Fatalf("partial decision %d %s", code, raw)
	}
	restore()
	code, raw = call("POST", "/"+request.ID+"/resume", maker, map[string]string{"campaign_sha256": partial.CampaignSHA256})
	var resumed CampaignBatch
	if code != 200 || json.Unmarshal(raw, &resumed) != nil || !resumed.Complete || resumed.Items[0].RequestSHA256 != partial.Items[0].RequestSHA256 {
		t.Fatalf("resume %d %s", code, raw)
	}
	code, raw = call("POST", "/"+request.ID+"/decision", reviewer, decision)
	if code != 200 || json.Unmarshal(raw, &reviewed) != nil || !reviewed.Complete {
		t.Fatalf("decision resume %d %s", code, raw)
	}
	code, raw = call("POST", "/"+request.ID+"/stop", maker, map[string]string{"campaign_sha256": partial.CampaignSHA256, "reason": "fixture finished before delivery"})
	if code != 200 {
		t.Fatalf("stop %d %s", code, raw)
	}
	for _, step := range request.Steps {
		if delay := time.Until(step.NotBefore) + 50*time.Millisecond; delay > 0 {
			time.Sleep(delay)
		}
		v, e := module.ProcessOnce(ctx, maker)
		if e != nil || v.Outcome != "SUPPRESSED" {
			t.Fatal(v, e)
		}
	}
	var requests, jobs, decisions int
	e = pool.QueryRow(ctx, `select(select count(*)from approval.request where tenant_id=$1),(select count(*)from platform.job where tenant_id=$1 and job_type='whatsapp.campaign.scheduled.v1'),(select count(*)from approval.decision where tenant_id=$1)`, tenant).Scan(&requests, &jobs, &decisions)
	if e != nil || requests != 2 || jobs != 2 || decisions != 2 || meta.Load() != 0 {
		t.Fatal("duplicate or unauthorized effect", requests, jobs, decisions, meta.Load(), e)
	}
	t.Logf("CAMPAIGN_PARTIAL_PROGRESS_PASS partial_create=visible partial_review=visible resumed_exactly=true requests=%d jobs=%d decisions=%d provider_posts=0 tenant=%s", requests, jobs, decisions, tenant)
}
