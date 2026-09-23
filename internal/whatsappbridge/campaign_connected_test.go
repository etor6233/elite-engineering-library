package whatsappbridge

// AUTHORED reference proof: original CRM fixture, HTTP identity, durable jobs,
// unchanged Meta adapter on loopback, current approval/fence/status/recovery.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/leadstream"
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
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCampaignConnected(t *testing.T) {
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
	permissions := map[string]struct{}{"lead:read": {}, "marketing:request": {}, "marketing:approve": {}, "marketing:read": {}, "marketing:cancel": {}, "notification:request": {}, "notification:approve": {}, "notification:read": {}, "notification:cancel": {}, "notification:dispatch": {}, "notification:reconcile": {}, "appointment:manage": {}, "whatsapp:process": {}}
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

	if _, e = pool.Exec(ctx, `update crm.lead set contact_payload='{"name":"Contacto de prueba","phone":"+5491112345678"}' where tenant_id=$1 and lead_id='lead-1'`, tenant); e != nil {
		t.Fatal(e)
	}
	state := "new"
	prepare := func(id, last string, steps int) CampaignSnapshot {
		t.Helper()
		now := time.Now().UTC()
		request := CampaignRequest{ID: id, Sources: []string{"fixture"}, States: []string{state}, Members: []CampaignRecipient{{LeadID: "lead-1", Recipient: old.Message.ExternalID}}}
		for i := 0; i < steps; i++ {
			label := last
			if i > 0 {
				label = "followup"
			}
			request.Steps = append(request.Steps, CampaignStep{TemplateName: "order_update", LanguageCode: "es_AR", BodyParameters: []string{id, label}, NotBefore: now.Add(time.Duration(2+i*3) * time.Second), ExpiresAt: now.Add(time.Duration(4+i*3) * time.Second)})
		}
		if code, _ := call("POST", "/prepare", foreign, request); code != 403 {
			t.Fatal("foreign campaign", code)
		}
		code, raw := call("POST", "/prepare", maker, request)
		if code != 200 {
			t.Fatalf("prepare %d %s", code, raw)
		}
		var out CampaignSnapshot
		if json.Unmarshal(raw, &out) != nil || len(out.Members) != 1 || out.Members[0].ConsentID == old.ConsentID {
			t.Fatal("marketing-specific source", string(raw))
		}
		return out
	}
	create := func(plan CampaignSnapshot) CampaignBatch {
		t.Helper()
		_, h, e := campaignHash(plan)
		if e != nil {
			t.Fatal(e)
		}
		selection := CampaignSelection{ID: plan.Request.ID, Steps: plan.Request.Steps}
		for _, m := range plan.Request.Members {
			selection.LeadIDs = append(selection.LeadIDs, m.LeadID)
		}
		code, raw := call("POST", "/workspace/create", maker, map[string]any{"selection": selection, "campaign_sha256": h})
		if code != 200 {
			t.Fatalf("create %d %s", code, raw)
		}
		var out CampaignBatch
		if json.Unmarshal(raw, &out) != nil || !out.Complete || len(out.Items) != len(plan.Request.Steps) {
			t.Fatalf("partial preparation %s", raw)
		}
		return out
	}
	review := func(plan CampaignSnapshot, batch CampaignBatch) {
		t.Helper()
		decision := map[string]any{"campaign_sha256": batch.CampaignSHA256, "approve": true, "reason": "explicit review of all stored members and timed template steps"}
		if code, _ := call("POST", "/"+plan.Request.ID+"/decision", maker, decision); code != 409 {
			t.Fatal("self review", code)
		}
		code, raw := call("POST", "/"+plan.Request.ID+"/decision", reviewer, decision)
		var out CampaignBatch
		if code != 200 || json.Unmarshal(raw, &out) != nil || !out.Complete {
			t.Fatalf("review %d %s", code, raw)
		}
	}
	due := func(plan CampaignSnapshot, step int) {
		t.Helper()
		delay := time.Until(plan.Request.Steps[step-1].NotBefore) + 50*time.Millisecond
		if delay > 0 {
			time.Sleep(delay)
		}
	}
	dispatch := func(want string) ScheduledWorkResult {
		t.Helper()
		v, e := module.ProcessOnce(ctx, maker)
		if e != nil || v.Outcome != want {
			t.Fatalf("dispatch %+v %v want%s", v, e, want)
		}
		return v
	}

	a := prepare("campaign-accepted-001", "accepted", 2)
	ab := create(a)
	// NEXT-U7: the same persisted campaign is selectable without transcribing IDs.
	if code, raw := call("GET", "/workspace", maker, nil); code != 200 || !bytes.Contains(raw, []byte(a.Request.ID)) {
		t.Fatalf("workspace list %d %s", code, raw)
	}
	if code, _ := call("GET", "/workspace", foreign, nil); code != 403 {
		t.Fatal("workspace foreign org", code)
	}
	if code, _ := call("GET", "/workspace?after=bad%2Fcursor", maker, nil); code != 400 {
		t.Fatal("workspace cursor", code)
	}

	if code, raw := call("GET", "/workspace?kind=audience", maker, nil); code != 200 || !bytes.Contains(raw, []byte("Contacto de prueba")) || bytes.Contains(raw, []byte("5491112345678")) {
		t.Fatalf("workspace eligible masked audience %d %s", code, raw)
	}
	if code, raw := call("GET", "/workspace/templates", maker, nil); code != 200 || !bytes.Contains(raw, []byte("order_update")) || bytes.Contains(raw, []byte("secret")) {
		t.Fatalf("workspace templates %d %s", code, raw)
	}
	selection := map[string]any{"campaign_id": "campaign-selected-001", "lead_ids": []string{"lead-1"}, "steps": a.Request.Steps}
	if code, raw := call("POST", "/workspace/prepare", maker, selection); code != 200 || !bytes.Contains(raw, []byte("campaign-selected-001")) {
		t.Fatalf("workspace prepare %d %s", code, raw)
	}
	selection["lead_ids"] = []string{"foreign-lead"}
	if code, _ := call("POST", "/workspace/prepare", maker, selection); code != 409 {
		t.Fatal("foreign selection", code)
	}

	replay := create(a)
	if replay.CampaignSHA256 != ab.CampaignSHA256 {
		t.Fatal("lost create reply changed snapshot")
	}
	review(a, ab)
	if v, e := module.ProcessOnce(ctx, maker); e != nil || v.Claimed {
		t.Fatal("early campaign effect", v, e)
	}
	due(a, 1)
	inactive, e := NewScheduleApprovals(base, tenant, "store-1", "wa-primary", profile)
	if e != nil {
		t.Fatal(e)
	}
	disabled, e := NewScheduledNotifications(inactive, sender, store, VerifiedInboxChannel{}, "disabled-campaign-fixture")
	if e != nil {
		t.Fatal(e)
	}
	if v, e := disabled.ProcessOnce(ctx, maker); e != nil || v.Claimed {
		t.Fatal("disabled campaign owner consumed its job", v, e)
	}
	var wg sync.WaitGroup
	outcomes := make(chan ScheduledWorkResult, 12)
	errs := make(chan error, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); v, e := module.ProcessOnce(ctx, maker); outcomes <- v; errs <- e }()
	}
	wg.Wait()
	close(outcomes)
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	accepted := 0
	for v := range outcomes {
		if v.Outcome == "ACCEPTED" {
			accepted++
		}
	}
	if accepted != 1 || meta.Load() != 1 {
		t.Fatal("concurrent first step", accepted, meta.Load())
	}
	due(a, 2)
	dispatch("ACCEPTED")
	if meta.Load() != 2 {
		t.Fatal("drip second step")
	}
	if v, e := campaigns.Status(ctx, maker, a.Request.ID); e != nil || len(v.Conversions) != 0 {
		t.Fatal("delivery invented a conversion", v, e)
	}

	b := prepare("campaign-unknown-0001", "lost-stdout", 3)
	bb := create(b)
	review(b, bb)
	due(b, 1)
	dispatch("RECONCILIATION_REQUIRED")
	due(b, 2)
	dispatch("SUPPRESSED")
	if _, e := module.Recover(ctx, maker, bb.Items[0].DeliveryKey, bb.Items[0].RequestSHA256); e != nil {
		t.Fatal("original receipt recovery", e)
	}
	due(b, 3)
	dispatch("SUPPRESSED")
	if meta.Load() != 3 {
		t.Fatal("unknown predecessor/recovery resent", meta.Load())
	}

	stopped := prepare("campaign-stopped-0001", "stopped", 1)
	sb := create(stopped)
	review(stopped, sb)
	if code, raw := call("POST", "/"+stopped.Request.ID+"/stop", maker, map[string]string{"campaign_sha256": sb.CampaignSHA256, "reason": "operator stopped campaign"}); code != 200 {
		t.Fatalf("stop %d %s", code, raw)
	}
	due(stopped, 1)
	dispatch("SUPPRESSED")

	withdrawn := prepare("campaign-withdraw-001", "withdrawn", 1)
	wb := create(withdrawn)
	review(withdrawn, wb)
	if _, e = pool.Exec(ctx, `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)values($1,'marketing-withdrawn','lead-1','fixture-marketing','policy-1','withdrawn',clock_timestamp(),repeat('b',64))`, tenant); e != nil {
		t.Fatal(e)
	}
	due(withdrawn, 1)
	dispatch("SUPPRESSED")
	if _, e = pool.Exec(ctx, `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)values($1,'marketing-renewed','lead-1','fixture-marketing','policy-1','granted',clock_timestamp(),repeat('c',64))`, tenant); e != nil {
		t.Fatal(e)
	}

	changed := prepare("campaign-changed-0001", "changed", 1)
	cb := create(changed)
	review(changed, cb)
	if _, e = postgres.NewFranchiseJourney(pool).TransitionLeadAs(ctx, tenant, "store-1", "lead-1", "new", "contacted", 1, leadstream.StableUUID(tenant, "campaign-contacted"), maker.Subject); e != nil {
		t.Fatal(e)
	}
	due(changed, 1)
	dispatch("SUPPRESSED")
	state = "contacted"

	converted := prepare("campaign-convert-0001", "conversion", 2)
	vb := create(converted)
	review(converted, vb)
	due(converted, 1)
	dispatch("ACCEPTED")
	repo := postgres.NewFranchiseJourney(pool)
	quoteValue, _, e := repo.CreateQuoteAs(ctx, tenant, "campaign-quote-request-1", franchisejourney.Quote{ID: "campaign-quote", OrganizationID: "store-1", LeadID: "lead-1", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().Add(time.Hour), State: "issued", Version: 1}, strings.Repeat("d", 64), leadstream.StableUUID(tenant, "campaign-quote-issued"), maker.Subject)
	if e != nil {
		t.Fatal("original quote owner", e)
	}
	acceptedQuote, e := repo.AcceptQuote(ctx, tenant, "store-1", "customer", quoteValue.ID, quoteValue.Version, strings.Repeat("e", 64), "campaign-order", "campaign-order-line", leadstream.StableUUID(tenant, "campaign-quote-accepted"), leadstream.StableUUID(tenant, "campaign-order-placed"))
	if e != nil || acceptedQuote.OrderID != "campaign-order" {
		t.Fatal("original quote->order owner", acceptedQuote, e)
	}
	due(converted, 2)
	dispatch("SUPPRESSED")
	code, raw := call("GET", "/"+converted.Request.ID, maker, nil)
	var status CampaignStatus
	if code != 200 || json.Unmarshal(raw, &status) != nil || len(status.Conversions) != 1 || status.Conversions[0].OrderID != "campaign-order" || status.Conversions[0].QuotationID != quoteValue.ID || status.Conversions[0].EvidenceSHA256 != strings.Repeat("e", 64) {
		t.Fatalf("conversion observation %d %s", code, raw)
	}
	// Exact replay after source transition/expiry uses stored snapshot and creates no new sends.
	oldReplay := create(a)
	if oldReplay.CampaignSHA256 != ab.CampaignSHA256 {
		t.Fatal("historical replay diverged")
	}
	if code, raw := call("GET", "/workspace/result/"+a.Request.ID, maker, nil); code != 200 || !bytes.Contains(raw, []byte(a.Request.ID)) {
		t.Fatalf("result %d %s", code, raw)
	}
	if code, _ := call("GET", "/workspace/result/campaign-missing-001", maker, nil); code != 404 {
		t.Fatal("missing result", code)
	}
	altered := CampaignSelection{ID: a.Request.ID, Steps: a.Request.Steps, LeadIDs: []string{"wrong-lead"}}
	if code, _ := call("POST", "/workspace/create", maker, map[string]any{"selection": altered, "campaign_sha256": ab.CampaignSHA256}); code != 409 {
		t.Fatal("altered replay", code)
	}
	if meta.Load() != 4 {
		t.Fatal("extra provider effect", meta.Load())
	}
	for _, table := range []string{"communication.whatsapp_campaign", "communication.whatsapp_campaign_member", "communication.whatsapp_campaign_stop"} {
		if _, e = pool.Exec(ctx, "delete from "+table+" where tenant_id=$1", tenant); e == nil {
			t.Fatal("mutable campaign evidence", table)
		}
	}
	var campaignsN, approvalsN, jobsN, fencesN, resultsN, ordersN int
	e = pool.QueryRow(ctx, `select(select count(*)from communication.whatsapp_campaign where tenant_id=$1),(select count(*)from approval.request where tenant_id=$1 and kind='whatsapp_schedule'),(select count(*)from platform.job where tenant_id=$1 and job_type='whatsapp.campaign.scheduled.v1'),(select count(*)from communication.outbound_delivery where tenant_id=$1 and channel_code='whatsapp'),(select count(*)from communication.whatsapp_schedule_result where tenant_id=$1),(select count(*)from sales.customer_order where tenant_id=$1)`, tenant).Scan(&campaignsN, &approvalsN, &jobsN, &fencesN, &resultsN, &ordersN)
	if e != nil || campaignsN != 6 || approvalsN != 10 || jobsN != 10 || fencesN != 4 || resultsN != 10 || ordersN != 1 {
		t.Fatal("counts", campaignsN, approvalsN, jobsN, fencesN, resultsN, ordersN, e)
	}
	if os.Getenv("ELITE_FUNCTIONS_BROWSER_SCRIPT") != "" {
		qualifyFunctionsBrowser(t, "campaign", api.URL, token(maker), maker)
	}
	t.Logf("CAMPAIGN_CONNECTED_PASS campaigns=%d approvals=%d jobs=%d fences=%d results=%d provider_posts=%d original_quote_orders=%d concurrent=12 no_causal_attribution tenant=%s", campaignsN, approvalsN, jobsN, fencesN, resultsN, meta.Load(), ordersN, tenant)
}
