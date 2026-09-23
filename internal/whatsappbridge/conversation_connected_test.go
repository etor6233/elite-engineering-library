package whatsappbridge

import (
	"context"
	"encoding/json"
	"fmt"
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

	"elite.local/enterprise/internal/app"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type replyDomainToken struct{}

func (replyDomainToken) AccessToken(context.Context) (string, error) {
	return "fixture-domain-token", nil
}

type replyVerifier map[string]identity.Principal

func (v replyVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	p, ok := v[token]
	if !ok {
		return p, ErrApproval
	}
	return p, nil
}

type connectedReplyFixture struct {
	pool                             *pgxpool.Pool
	observer                         *StatusObserver
	worker                           *StatusWorker
	router                           *StatusRouter
	module                           *ReplyModule
	receiver                         *httptest.Server
	api                              *httptest.Server
	service, human                   identity.Principal
	metaCalls, llmCalls, domainCalls *atomic.Int32
}

func newConnectedReplyFixture(t *testing.T) *connectedReplyFixture {
	t.Helper()
	ctx := context.Background()
	db := os.Getenv("ELITE_WHATSAPP_CONNECTED_DATABASE_URL")
	u, err := url.Parse(db)
	if err != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_whatsapp_connected_") {
		t.Fatal("dedicated loopback WhatsApp database required")
	}
	pool, err := pgxpool.New(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	tenant := leadstream.StableUUID(t.Name(), time.Now().Format(time.RFC3339Nano))
	key := []byte("0123456789abcdef0123456789abcdef")
	for _, sql := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'wa-'||$1::text,'Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store-1','store-1','Fixture','store')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required)values($1,'lead-1','store-1','new','fixture','{}',true)`,
		`insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)values($1,'consent-1','lead-1','fixture-conversation','policy-1','granted',clock_timestamp()-interval '1 minute',repeat('c',64))`,
		`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref,organization_id)values($1,'wa-primary','meta-whatsapp','META_APP_SECRET','store-1')`,
	} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	contacts, err := postgres.NewContactIdentityStore(pool, key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = contacts.Apply(ctx, contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(-time.Minute), RequestID: "binding-1"}); err != nil {
		t.Fatal(err)
	}
	base, err := NewPostgresAppointmentApprovals(pool, key, "fixture-conversation", "policy-1")
	if err != nil {
		t.Fatal(err)
	}
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
	var meta, llmCount, domainCount atomic.Int32
	metaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/messages" || r.Header.Get("Authorization") != "Bearer synthetic-token" {
			t.Error("unexpected fixture Meta transport")
		}
		body, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(body, &m)
		if m["type"] != "text" || m["to"] != "5491112345678" || m["messaging_product"] != "whatsapp" {
			t.Errorf("provider payload: %s", body)
		}
		n := meta.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"messaging_product":"whatsapp","messages":[{"id":"wamid.reply.%d"}]}`, n)
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
 if json.loads(raw)['request'].get('source_message_id')=='wamid.lost':sys.exit(2)
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
	observer := &StatusObserver{TenantID: tenant, Profile: profile, Process: process, ReconcilerSHA256: digest(statusCode), Secrets: appSecretFixture{}, Approvals: base}
	replies, err := NewPostgresReplyApprovals(base, tenant, "store-1", "wa-primary", profile)
	if err != nil {
		t.Fatal(err)
	}
	store, err := postgres.NewOutboundDeliveryStore(pool, key, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	sender := &Sender{TenantID: tenant, Profile: profile, Process: process, Tokens: tokenFixture{}}
	module, err := NewReplyModule(replies, sender, store)
	if err != nil {
		t.Fatal(err)
	}
	service := identity.Principal{Subject: "fixture-service", TenantID: tenant, Organizations: map[string]struct{}{"store-1": {}}, Permissions: map[string]struct{}{"appointment:manage": {}, "whatsapp:process": {}}}
	human := identity.Principal{Subject: "fixture-human", TenantID: tenant, Organizations: map[string]struct{}{"store-1": {}}, Permissions: map[string]struct{}{"whatsapp:approve": {}, "whatsapp:send": {}}}
	llm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := llmCount.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if n%2 == 1 {
			fmt.Fprintf(w, `{"id":"resp%d","status":"completed","output":[{"type":"function_call","call_id":"call%d","name":"create_quote","arguments":"{\"product\":\"scooter\",\"quantity\":1}"}],"usage":{"input_tokens":12,"output_tokens":4,"total_tokens":16}}`, n, n)
			return
		}
		fmt.Fprintf(w, `{"id":"resp%d","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Cotización preparada para revisar."}]}],"usage":{"input_tokens":8,"output_tokens":5,"total_tokens":13}}`, n)
	}))
	t.Cleanup(llm.Close)
	domain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainCount.Add(1)
		if r.URL.Path != "/v1/franchise/quotes" || r.Header.Get("Authorization") != "Bearer fixture-domain-token" || len(r.Header.Get("Idempotency-Key")) != 64 {
			t.Error("unbound domain tool")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"quotation_id":"fixture-quote"}`))
	}))
	t.Cleanup(domain.Close)
	turns, err := postgres.NewConversationStore(pool, time.Minute, 48*time.Hour, 3)
	if err != nil {
		t.Fatal(err)
	}
	registry := channels.NewRegistry()
	_ = registry.Register(VerifiedInboxChannel{})
	runtime, err := app.New(app.Config{LLMAPIKey: "synthetic-llm-key", LLMModel: "fixture-model", LLMBaseURL: llm.URL, DomainBaseURL: domain.URL, DomainTokenProvider: replyDomainToken{}, TenantID: tenant, TenantCode: "fixture", ServiceKinds: map[string]string{"consulta": "consultation"}, ProductVariants: map[string]string{"scooter": "scooter-1"}, PriceBookID: "fixture-book", ConversationStore: turns, ContactResolver: ScopedContactResolver{TenantID: tenant, OrganizationID: "store-1", Resolver: contacts}, ChannelRegistry: registry, LLMTokenBudget: 100000, Conversation: conversationruntime.Config{Instructions: "Synthetic fixture; propose and use only admitted tools.", PromptCacheKey: "fixture-wa", MaxOutputTokens: 128, MaxHistory: 4, MaxInputBytes: 16384, MaxToolBytes: 16384, StoreApproved: true, HandoffText: "Requiere revisión humana."}})
	if err != nil {
		t.Fatal(err)
	}
	router, err := NewStatusRouter(observer, "wa-primary", strings.Repeat("a", 64), key, 8)
	if err != nil {
		t.Fatal(err)
	}
	if err = router.EnableConversation(ConversationRoute{OrganizationID: "store-1", Runtime: runtime.Conversation, Proposals: replies}); err != nil {
		t.Fatal(err)
	}
	worker, err := NewStatusWorker(router, "wa-connected", 2*time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	receiver, err := NewWebhookReceiver(receiverConfig(observer))
	if err != nil {
		t.Fatal(err)
	}
	webhook := httptest.NewServer(receiver)
	t.Cleanup(webhook.Close)
	mux := http.NewServeMux()
	module.Register(mux, replyVerifier{"human": human, "service": service})
	api := httptest.NewServer(mux)
	t.Cleanup(api.Close)
	return &connectedReplyFixture{pool: pool, observer: observer, worker: worker, router: router, module: module, receiver: webhook, api: api, service: service, human: human, metaCalls: &meta, llmCalls: &llmCount, domainCalls: &domainCount}
}
func (f *connectedReplyFixture) inbound(t *testing.T, id string, at time.Time, recipient string) SignedStatusWebhook {
	t.Helper()
	return alterStatusBatch(t, statusBatch("sent", fmt.Sprint(at.Unix()), nil), func(body map[string]any) {
		v := statusValue(body)
		delete(v, "statuses")
		v["messages"] = []any{map[string]any{"id": id, "from": recipient, "type": "text", "timestamp": fmt.Sprint(at.Unix()), "text": map[string]any{"body": "Cotizame un scooter"}}}
	})
}
func (f *connectedReplyFixture) ingest(t *testing.T, b SignedStatusWebhook) {
	t.Helper()
	res, _ := webhookRequest(t, f.receiver, b)
	if res.StatusCode != 200 {
		t.Fatalf("ingress: %d", res.StatusCode)
	}
}
func (f *connectedReplyFixture) process(t *testing.T) {
	t.Helper()
	r, err := f.worker.ProcessOnce(context.Background(), f.service)
	if err != nil || !r.Completed {
		t.Fatalf("worker %+v %v", r, err)
	}
}
func (f *connectedReplyFixture) command(t *testing.T, key, op, token string, body any) int {
	t.Helper()
	raw, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", f.api.URL+"/v1/franchise/whatsapp/replies/"+key+"/"+op, strings.NewReader(string(raw)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := f.api.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	t.Logf("%s %d %s", op, res.StatusCode, string(b))
	return res.StatusCode
}
func (f *connectedReplyFixture) proposal(t *testing.T, id string) ReplyProposal {
	t.Helper()
	v, err := f.module.approvals.Read(context.Background(), f.human, ReplyDeliveryKey(f.service.TenantID, "wa-primary", id))
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func TestWhatsAppConnectedHumanReplyAndStatus(t *testing.T) {
	f := newConnectedReplyFixture(t)
	b := f.inbound(t, "wamid.inbound", time.Now().Add(-time.Second), "5491112345678")
	f.ingest(t, b)
	f.ingest(t, b)
	f.process(t)
	proposal := f.proposal(t, "wamid.inbound")
	if proposal.State != "pending" || f.metaCalls.Load() != 0 || f.llmCalls.Load() != 2 || f.domainCalls.Load() != 1 {
		t.Fatal("runtime did not stop at proposal", proposal.State, f.metaCalls.Load(), f.llmCalls.Load(), f.domainCalls.Load())
	}
	if f.command(t, proposal.RequestID, "send", "human", map[string]any{"payload_sha256": proposal.PayloadSHA256}) != 409 {
		t.Fatal("pending sent")
	}
	if f.command(t, proposal.RequestID, "decision", "human", map[string]any{"payload_sha256": strings.Repeat("f", 64), "approved": true, "reason": "review"}) != 409 {
		t.Fatal("changed hash approved")
	}
	if f.command(t, proposal.RequestID, "decision", "human", map[string]any{"payload_sha256": proposal.PayloadSHA256, "approved": true, "reason": "Reviewed exact response"}) != 200 {
		t.Fatal("approval failed")
	}
	for i := 0; i < 2; i++ {
		if f.command(t, proposal.RequestID, "send", "human", map[string]any{"payload_sha256": proposal.PayloadSHA256}) != 200 {
			t.Fatal("send/replay failed")
		}
	}
	if f.metaCalls.Load() != 1 {
		t.Fatal("provider send duplicated")
	}
	status := statusBatch("delivered", fmt.Sprint(time.Now().Unix()), map[string]any{"id": "wamid.reply.1"})
	f.ingest(t, status)
	f.process(t)
	f.ingest(t, status)
	result, err := f.worker.ProcessOnce(context.Background(), f.service)
	if err != nil || result.Claimed {
		t.Fatal("status replay reprocessed", result, err)
	}
	var observations, attempts int
	var state string
	if err = f.pool.QueryRow(context.Background(), `select state,attempt_count from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2`, f.service.TenantID, proposal.RequestID).Scan(&state, &attempts); err != nil {
		t.Fatal(err)
	}
	if err = f.pool.QueryRow(context.Background(), `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, f.service.TenantID).Scan(&observations); err != nil {
		t.Fatal(err)
	}
	read := f.proposal(t, "wamid.inbound")
	if read.Status == nil || read.Status.FenceState != "accepted" || read.Status.DeliveryStatus != "observed_delivered" {
		t.Fatal("human status projection missing", read.Status)
	}
	if state != "accepted" || attempts != 1 || observations != 1 || f.domainCalls.Load() != 1 {
		t.Fatal("unexpected connected effects", state, attempts, observations)
	}
}
func TestWhatsAppConnectedLostResponseRecoveryNoSecondPost(t *testing.T) {
	f := newConnectedReplyFixture(t)
	f.ingest(t, f.inbound(t, "wamid.lost", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	v := f.proposal(t, "wamid.lost")
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Reviewed"}) != 200 {
		t.Fatal("approval")
	}
	for i := 0; i < 2; i++ {
		if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
			t.Fatal("uncertain accepted")
		}
	}
	if f.metaCalls.Load() != 1 {
		t.Fatal("expected exactly one fixture POST", f.metaCalls.Load())
	}
	if f.command(t, v.RequestID, "recover", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 200 {
		t.Fatal("evidence recovery")
	}
	f.ingest(t, statusBatch("read", fmt.Sprint(time.Now().Unix()), map[string]any{"id": "wamid.reply.1"}))
	f.process(t)
	if f.metaCalls.Load() != 1 {
		t.Fatal("reconciliation sent again")
	}
}
func TestWhatsAppConnectedRejectWindowAndUnboundContact(t *testing.T) {
	f := newConnectedReplyFixture(t)
	f.ingest(t, f.inbound(t, "wamid.reject", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	v := f.proposal(t, "wamid.reject")
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": false, "reason": "Rejected"}) != 200 {
		t.Fatal("reject")
	}
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Retry approval"}) != 409 {
		t.Fatal("reject mutated")
	}
	if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("reject sent")
	}
	f.ingest(t, f.inbound(t, "wamid.unknown", time.Now().Add(-time.Second), "5491199999999"))
	f.process(t)
	if f.llmCalls.Load() != 2 || f.domainCalls.Load() != 1 || f.metaCalls.Load() != 0 {
		t.Fatal("unbound contact reached model/tool/provider")
	}
	var handed bool
	if err := f.pool.QueryRow(context.Background(), `select state='handed_off' and failure_code='CONTACT_SCOPE_UNAVAILABLE' from communication.conversation_turn where tenant_id=$1 and provider_message_id='wamid.unknown'`, f.service.TenantID).Scan(&handed); err != nil || !handed {
		t.Fatal("missing durable handoff", err)
	}
	// Local clock boundary for expired text is also enforced independently before
	// the provider process, even with an otherwise exact typed message.
	expired := v.Context.Message
	var request ReplyRequest
	_ = json.Unmarshal([]byte(expired.Text), &request)
	request.LastInboundAt = time.Now().Add(-25 * time.Hour).Unix()
	request.WindowExpiresAt = request.LastInboundAt + 86400
	raw, _ := json.Marshal(request)
	expired.Text = string(raw)
	if _, err := replyRequestFor(expired); err == nil {
		t.Fatal("closed window admitted")
	}
}

// Failure injection at the existing durable Store boundary: provider returns a
// receipt, then process loss prevents Complete. No SQL fabricates a send state.
type replyCrashStore struct {
	*postgres.OutboundDeliveryStore
}

func (replyCrashStore) Complete(context.Context, channels.Message, string, outbounddelivery.Receipt) error {
	return outbounddelivery.ErrUnknown
}
func TestWhatsAppConnectedLeaseCrashAndMissingEvidence(t *testing.T) {
	f := newConnectedReplyFixture(t)
	f.ingest(t, f.inbound(t, "wamid.crash", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	v := f.proposal(t, "wamid.crash")
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Reviewed"}) != 200 {
		t.Fatal("approval")
	}
	f.module.channel.Store = replyCrashStore{f.module.store}
	if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("crash not surfaced")
	}
	dir := filepath.Join(f.module.sender.Process.EvidenceDirectory, digest([]byte(f.service.TenantID+"\x00"+v.RequestID)))
	receipt := filepath.Join(dir, "SEND_RECEIPT.json")
	saved, err := os.ReadFile(receipt)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(receipt, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	if f.command(t, v.RequestID, "recover", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("unproved receipt accepted")
	}
	if err = os.WriteFile(receipt, saved, 0600); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1100 * time.Millisecond)
	for i := 0; i < 2; i++ {
		if f.command(t, v.RequestID, "recover", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 200 {
			t.Fatal("expired lease recovery/replay failed")
		}
	}
	if f.metaCalls.Load() != 1 {
		t.Fatal("crash recovery duplicated POST")
	}
}

func TestWhatsAppFutureContactCannotReachRuntime(t *testing.T) {
	f := newConnectedReplyFixture(t)
	binding, err := postgres.NewContactIdentityStore(f.pool, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = binding.Apply(context.Background(), contactidentity.Command{TenantID: f.service.TenantID, ChannelCode: "whatsapp", ExternalID: "5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(time.Hour), ExpectedVersion: 1, RequestID: "binding-future"})
	if err != nil {
		t.Fatal(err)
	}
	f.ingest(t, f.inbound(t, "wamid.future", time.Now().Add(-time.Second), "5491112345678"))
	result, err := f.worker.ProcessOnce(context.Background(), f.service)
	if f.llmCalls.Load() != 0 || f.domainCalls.Load() != 0 || f.metaCalls.Load() != 0 {
		t.Fatalf("future contact escaped before effective_at: model=%d domain=%d provider=%d", f.llmCalls.Load(), f.domainCalls.Load(), f.metaCalls.Load())
	}
	if err != nil || !result.Completed {
		t.Fatal("future contact should produce durable handoff", result, err)
	}
}

func TestWhatsAppReplyApprovalCannotSurviveScopeOrBindingDrift(t *testing.T) {
	f := newConnectedReplyFixture(t)
	ctx := context.Background()
	f.ingest(t, f.inbound(t, "wamid.binding", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	v := f.proposal(t, "wamid.binding")
	self := f.service
	self.Permissions = map[string]struct{}{"whatsapp:approve": {}}
	if _, err := f.module.approvals.Decide(ctx, self, v.RequestID, v.PayloadSHA256, true, "self review"); err == nil {
		t.Fatal("service requester approved own reply")
	}
	foreign := f.human
	foreign.Organizations = map[string]struct{}{"other": {}}
	if _, err := f.module.approvals.Read(ctx, foreign, v.RequestID); err == nil {
		t.Fatal("foreign organization read exact reply")
	}
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Reviewed"}) != 200 {
		t.Fatal("approval")
	}
	store, err := postgres.NewContactIdentityStore(f.pool, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.Apply(ctx, contactidentity.Command{TenantID: f.service.TenantID, ChannelCode: "whatsapp", ExternalID: "5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: false, State: contactidentity.StateRevoked, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(-time.Second), ExpectedVersion: 1, RequestID: "revoke-contact"})
	if err != nil {
		t.Fatal(err)
	}
	if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("revoked binding sent")
	}
	_, err = store.Apply(ctx, contactidentity.Command{TenantID: f.service.TenantID, ChannelCode: "whatsapp", ExternalID: "5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(-time.Second), ExpectedVersion: 2, RequestID: "rebind-contact"})
	if err != nil {
		t.Fatal(err)
	}
	if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("new binding silently renewed old approval")
	}
	if f.metaCalls.Load() != 0 {
		t.Fatal("drift reached provider")
	}
}
