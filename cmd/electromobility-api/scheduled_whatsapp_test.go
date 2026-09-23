package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/whatsappbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AUTHORED host integration fixture derived from the existing local assembly fixture.
// Exercises exact scheduled migration, both owned loops, authenticated telemetry and shutdown.
func TestScheduledWhatsAppHost(t *testing.T) {
	python := os.Getenv("ELITE_WHATSAPP_HOST_PYTHON")
	if python == "" {
		t.Skip("explicit pinned Python runtime is required for this constructor gate")
	}
	dir := t.TempDir()
	hashFile := func(p string) string {
		t.Helper()
		b, e := os.ReadFile(p)
		if e != nil {
			t.Fatal(e)
		}
		h := sha256.Sum256(b)
		return hex.EncodeToString(h[:])
	}
	write := func(name string, b []byte) string {
		t.Helper()
		p := filepath.Join(dir, name)
		if e := os.WriteFile(p, b, 0600); e != nil {
			t.Fatal(e)
		}
		return p
	}
	adapter, e := filepath.Abs("../../whatsapp_cloud")
	if e != nil {
		t.Fatal(e)
	}
	provider := map[string]any{"provider": "Meta WhatsApp Cloud API", "source_id": "meta-whatsapp-api-examples", "source_commit": "de70ee908a67026e642aaee3703d20464e2a9466", "decision": "PROVEN", "business_account_id": "987654321", "graph_api_version": "v99.0", "phone_number_id": "123456789", "access_token_environment_variable": "WHATSAPP_ACCESS_TOKEN", "app_secret_environment_variable": "META_APP_SECRET", "verify_token_environment_variable": "WHATSAPP_VERIFY_TOKEN", "approved_templates": []any{map[string]any{"name": "order_update", "language_code": "es_AR", "body_parameter_count": 2}}, "data_retention": "synthetic local fixture only", "automatic_business_write": false}
	for _, k := range []string{"official_source_license_accepted", "platform_terms_accepted", "business_account_proven", "app_registration_proven", "phone_number_id_proven", "message_template_approval_proven", "webhook_subscription_proven", "test_recipient_consent_proven", "quota_and_cost_approved", "reconciliation_approved"} {
		provider[k] = true
	}
	providerRaw, _ := json.Marshal(provider)
	profilePath := write("synthetic-provider.json", providerRaw)
	c, values, reports := whatsappTLSFixture(t, 2)
	c.Schema = "elite-whatsapp-host/v1"
	c.TenantID = "11111111-1111-4111-8111-111111111111"
	c.TenantCode = "fixture"
	c.OrganizationID = "fixture-org"
	c.ConnectionID = "fixture-whatsapp"
	c.Purpose = "fixture-conversation"
	c.PolicyVersion = "fixture-policy-v1"
	c.RetentionApprovalSHA256 = strings.Repeat("1", 64)
	c.ProviderProfileFile = profilePath
	c.ProviderProfileSHA256 = hashFile(profilePath)
	c.Process = whatsappbridge.Process{PythonExecutable: python, PythonSHA256: hashFile(python), AdapterDirectory: adapter, AdapterSHA256: hashFile(filepath.Join(adapter, "whatsapp_cloud.py")), EvidenceDirectory: dir}
	c.ReconcilerSHA256 = hashFile(filepath.Join(adapter, "status_reconciliation.py"))
	c.WorkerID = "fixture-worker"
	c.PollSeconds = 1
	c.ConversationRetentionSeconds = 3600
	c.ConversationMaxAttempts = 2
	c.LLMModel = "fixture-model"
	c.LLMBaseURL = "https://llm.example.invalid/v1"
	c.LLMTokenBudget = 10000
	c.DomainBaseURL = "https://domain.example.invalid"
	c.ServiceKinds = map[string]string{"fixture": "service"}
	c.ProductVariants = map[string]string{"fixture": "variant"}
	c.PriceBookID = "fixture-price-book"
	c.Conversation = conversationruntime.Config{Instructions: "Synthetic fixture only", PromptCacheKey: "fixture", MaxOutputTokens: 32, MaxHistory: 1, MaxInputBytes: 4096, MaxToolBytes: 4096, StoreApproved: true, HandoffText: "Fixture human handoff"}
	c.ApprovalPolicy = approval.Policy{MaxOpenPerSubject: 5}
	for _, k := range []string{"WHATSAPP_ACCESS_TOKEN_FILE", "WHATSAPP_APP_SECRET_FILE", "WHATSAPP_VERIFY_TOKEN_FILE", "WHATSAPP_SERVICE_TOKEN_FILE", "LLM_API_KEY_FILE"} {
		values[k] = write(k, []byte("synthetic-fixture-secret"))
	}
	values["WHATSAPP_CONTACT_HMAC_KEY_HEX_FILE"] = write("hmac", []byte(strings.Repeat("a1", 32)))
	values["WHATSAPP_ENABLED"] = "true"
	values["WHATSAPP_SCHEDULE_ENABLED"] = "true"
	configRaw, _ := json.Marshal(c)
	values["WHATSAPP_HOST_PROFILE_FILE"] = write("host.json", configRaw)
	values["WHATSAPP_HOST_PROFILE_SHA256"] = hashFile(values["WHATSAPP_HOST_PROFILE_FILE"])
	p := identity.Principal{Subject: "fixture-service", TenantID: c.TenantID, Permissions: map[string]struct{}{"appointment:manage": {}, "whatsapp:process": {}, "notification:dispatch": {}}, Organizations: map[string]struct{}{c.OrganizationID: {}}}
	v := &hostWAIdentityVerifier{p: p, token: "synthetic-fixture-secret"}
	dsn := os.Getenv("ELITE_WHATSAPP_CONNECTED_DATABASE_URL")
	if dsn == "" {
		t.Skip("explicit isolated PostgreSQL fixture required")
	}
	u, e := url.Parse(dsn)
	if e != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_whatsapp_connected_") {
		t.Fatal("fixture database scope")
	}
	pool, e := pgxpool.New(context.Background(), dsn)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	if err := pool.QueryRow(context.Background(), "select tenant_id::text,organization_id,connection_id from integration.provider_connection where provider_code='meta-whatsapp' and state='active'").Scan(&c.TenantID, &c.OrganizationID, &c.ConnectionID); err != nil {
		t.Fatal(err)
	}
	v.p.TenantID = c.TenantID
	v.p.Organizations = map[string]struct{}{c.OrganizationID: {}}
	configRaw, _ = json.Marshal(c)
	values["WHATSAPP_HOST_PROFILE_FILE"] = write("connected-host.json", configRaw)
	values["WHATSAPP_HOST_PROFILE_SHA256"] = hashFile(values["WHATSAPP_HOST_PROFILE_FILE"])
	h, e := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] })
	if e != nil {
		t.Fatal(e)
	}
	defer h.close()
	if h.(*whatsappHost).application == nil || h.(*whatsappHost).application.Conversation == nil || h.(*whatsappHost).worker == nil || h.(*whatsappHost).ingress == nil || len(h.(*whatsappHost).modules) != 3 || h.(*whatsappHost).extraRun == nil {
		t.Fatal("missing existing owner")
	}
	mux := http.NewServeMux()
	h.Register(mux, v)
	for _, path := range []string{"/v1/franchise/whatsapp/replies", "/v1/providers/whatsapp/webhook", "/v1/franchise/notifications/scheduled/prepare"} {
		r := httptest.NewRequest("DELETE", path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code == 404 {
			t.Fatalf("unmounted %s", path)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- h.run(ctx) }()
	for i := 0; i < 2; i++ {
		select {
		case raw := <-reports:
			var report map[string]any
			if json.Unmarshal(raw, &report) != nil || report["outcome"] != "IDLE" {
				t.Fatalf("host loop report %s", raw)
			}
		case <-time.After(6 * time.Second):
			t.Fatal("both owned loops did not report")
		}
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal("shutdown", err)
		}
	case <-time.After(6 * time.Second):
		t.Fatal("host leaked worker")
	}
	saved := whatsappScheduleFactory
	whatsappScheduleFactory = nil
	if bad, err := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] }); err == nil {
		bad.close()
		t.Fatal("missing schedule pack activated")
	}
	whatsappScheduleFactory = saved
	values["WHATSAPP_SCHEDULE_ENABLED"] = "TRUE"
	if bad, err := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] }); err == nil {
		bad.close()
		t.Fatal("nonexact flag activated")
	}
	values["WHATSAPP_SCHEDULE_ENABLED"] = "true"
	delete(v.p.Permissions, "notification:dispatch")
	if bad, err := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] }); err == nil {
		bad.close()
		t.Fatal("unprivileged worker activated")
	}
	v.p.Permissions["notification:dispatch"] = struct{}{}
	if _, err := pool.Exec(context.Background(), "alter table communication.whatsapp_schedule disable trigger whatsapp_schedule_immutable"); err != nil {
		t.Fatal(err)
	}
	defer pool.Exec(context.Background(), "alter table communication.whatsapp_schedule enable trigger whatsapp_schedule_immutable")
	if bad, err := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] }); err == nil {
		bad.close()
		t.Fatal("disabled immutable guard activated")
	}
	t.Log("SCHEDULED_HOST_PASS mounted=3 loops=2 telemetry=mTLS shutdown=joined missing_pack_flag_permission_guard=rejected")
	// The normal profile remains blocked; the constructor cannot activate it.
	provider["decision"] = "BLOCKED_ACCESS_TERMS_AND_RECONCILIATION_REQUIRED"
	blocked, _ := json.Marshal(provider)
	c.ProviderProfileFile = write("blocked-provider.json", blocked)
	c.ProviderProfileSHA256 = hashFile(c.ProviderProfileFile)
	configRaw, _ = json.Marshal(c)
	values["WHATSAPP_HOST_PROFILE_FILE"] = write("blocked-host.json", configRaw)
	values["WHATSAPP_HOST_PROFILE_SHA256"] = hashFile(values["WHATSAPP_HOST_PROFILE_FILE"])
	if bad, e := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] }); e == nil {
		bad.close()
		t.Fatal("blocked provider activated")
	}
}
