# Connected scheduled WhatsApp reminders

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-SCHEDULED-WHATSAPP"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "PROVEN_LOCAL appointment reminders on original CRM/approval/jobs/WhatsApp/fence/host, with cancel/reschedule/consent suppression, exact receipt recovery and signed status. AUTHORED glue; no whole T2805 or live production claim."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["GO-CONNECTED-WHATSAPP-HOST", "GO-HUMAN-APPROVAL-CORE", "GO-RELIABLE-ASYNC-WORKERS", "GO-PG-OUTBOUND-DELIVERY-FENCE", "PYTHON-META-WHATSAPP-CLOUD-ADAPTER", "GO-ELECTROMOBILITY-APPLICATION"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["whatsapp_cloud/official-source.lock.json"]
verified_at: "2026-09-13"
```

## 2. Applicability

Selected appointment reminder source and original WhatsApp host; local fixtures and future user credentials.

## 3. Architecture contract

Original CRM appointment/version/consent -> prepared exact template -> distinct manual approval -> durable due job -> original outbound fence/actual Python adapter -> immutable receipt/status. No new provider retry or source algorithm.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/scheduled_whatsapp.go
CREATE cmd/electromobility-api/scheduled_whatsapp_test.go
CREATE db/migrations/0083_scheduled_whatsapp.down.sql
CREATE db/migrations/0083_scheduled_whatsapp.up.sql
CREATE docs/SCHEDULED_COMMUNICATIONS_CONTRACT_PLAN.md
CREATE docs/SCHEDULED_WHATSAPP_REFERENCE.md
CREATE internal/platform/postgres/scheduled_outbound_binding.go
CREATE internal/whatsappbridge/scheduled_boundary_test.go
CREATE internal/whatsappbridge/scheduled_connected_test.go
CREATE internal/whatsappbridge/scheduled_contract.go
CREATE internal/whatsappbridge/scheduled_exhaustion.go
CREATE internal/whatsappbridge/scheduled_module.go
CREATE internal/whatsappbridge/scheduled_recovery.go
CREATE internal/whatsappbridge/scheduled_store.go
CREATE internal/whatsappbridge/scheduled_worker.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/scheduled_whatsapp.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bde351e87c6956b7a5bd3c73ba175ed032d0aae540e1f32229b72697b6871775"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED optional same-host composition; all credential and identity paths
// remain the original WhatsApp host's verified runtime sources.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/whatsappbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() { whatsappScheduleFactory = selectedWhatsAppSchedule }
func selectedWhatsAppSchedule(ctx context.Context, h *whatsappHost, c whatsappHostConfig, pool *pgxpool.Pool, base *whatsappbridge.PostgresAppointmentApprovals, sender *whatsappbridge.Sender, fence *postgres.OutboundDeliveryStore) error {
	if h == nil || pool == nil {
		return errWhatsAppHost
	}
	p, e := h.identity.Resolve(ctx)
	if e != nil || !p.Allowed("notification:dispatch") {
		return errWhatsAppHost
	}
	approvals, e := whatsappbridge.NewScheduleApprovals(base, c.TenantID, c.OrganizationID, c.ConnectionID, sender.Profile)
	if e != nil {
		return e
	}
	if e = approvals.VerifyInfrastructure(ctx); e != nil {
		return e
	}
	module, e := whatsappbridge.NewScheduledNotifications(approvals, sender, fence, whatsappbridge.VerifiedInboxChannel{}, c.WorkerID+"-scheduled")
	if e != nil {
		return e
	}
	if c.CampaignPolicyFile != "" || c.CampaignPolicySHA256 != "" {
		if whatsappCampaignFactory == nil {
			return errWhatsAppHost
		}
		campaign, e := whatsappCampaignFactory(ctx, module, c)
		if e != nil {
			return e
		}
		h.modules = append(h.modules, campaign)
	}
	h.modules = append(h.modules, module)
	h.extraRun = func(ctx context.Context) error { return module.Run(ctx, h.identity, h.reporter, h.interval) }
	return nil
}

var whatsappCampaignFactory func(context.Context, *whatsappbridge.ScheduledNotifications, whatsappHostConfig) (httpapi.EnterpriseModule, error)
````

### FILE: `cmd/electromobility-api/scheduled_whatsapp_test.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6addc9f7ec6c97c229c19b62390eab65c89af06a29a466a17753299df496e2bd"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `db/migrations/0083_scheduled_whatsapp.down.sql`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6da1ab5c7f1b10d209878e3b31b0f8cc5b570572bc38d35731ab76e6b38ba42f"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
 if exists(select 1 from approval.request where kind='whatsapp_schedule') or exists(select 1 from communication.whatsapp_schedule)
 then raise exception 'cannot remove scheduled WhatsApp evidence';end if;
end $$;
create or replace view communication.whatsapp_delivery_approval as
 select g.tenant_id,g.channel_code,g.delivery_key,g.organization_id,g.appointment_id,g.confirmation_event_id::text as approval_event_id,g.lead_id,g.external_id_hmac,g.profile_sha256,g.message_sha256,g.expires_at
 from communication.whatsapp_appointment_approval g
 union all
 select r.tenant_id,'whatsapp',r.request_id,r.organization_id,''::text,r.request_id,
 r.payload->>'lead_id',r.payload->>'external_id_hmac',r.payload->>'profile_sha256',r.payload->>'message_sha256',(r.payload->>'expires_at')::timestamptz
 from approval.request r
 where r.kind='whatsapp_reply' and r.state='approved' and r.payload->>'schema'='elite-whatsapp-reply-approval/v1'
 and r.payload->>'organization_id'=r.organization_id;
drop table communication.whatsapp_schedule_result;
drop table communication.whatsapp_schedule_cancellation;
drop table communication.whatsapp_schedule;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

commit;
````

### FILE: `db/migrations/0083_scheduled_whatsapp.up.sql`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1d23e7405dac18866e8091d500530c5f2da7d5864798afb0062200b59f1fb915"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;


create table communication.whatsapp_schedule(
 tenant_id uuid not null,delivery_key text not null,job_id uuid not null,
 organization_id text not null,appointment_id text not null,appointment_version bigint not null,
 not_before timestamptz not null,expires_at timestamptz not null,
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 primary key(tenant_id,delivery_key),unique(tenant_id,job_id),
 unique(tenant_id,appointment_id,appointment_version,not_before),
 foreign key(tenant_id,delivery_key) references approval.request(tenant_id,request_id),
 foreign key(tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id),
 check(expires_at>not_before)
);
create trigger whatsapp_schedule_immutable before update or delete on communication.whatsapp_schedule
 for each row execute function catalog.release_immutable();
create table communication.whatsapp_schedule_cancellation(
 tenant_id uuid not null,delivery_key text not null,actor_subject text not null,reason text not null,
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 cancelled_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,delivery_key),
 foreign key(tenant_id,delivery_key) references communication.whatsapp_schedule(tenant_id,delivery_key),
 check(length(actor_subject)between 1 and 128),check(length(reason)between 1 and 2048)
);
create trigger whatsapp_schedule_cancellation_immutable before update or delete on communication.whatsapp_schedule_cancellation
 for each row execute function catalog.release_immutable();
create table communication.whatsapp_schedule_result(
 tenant_id uuid not null,delivery_key text not null,job_id uuid not null,
 outcome text not null check(outcome in('ACCEPTED','SUPPRESSED','RECONCILIATION_REQUIRED','FAILED_TERMINAL')),
 completed_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,delivery_key),
 foreign key(tenant_id,delivery_key) references communication.whatsapp_schedule(tenant_id,delivery_key),
 foreign key(tenant_id,job_id) references platform.job(tenant_id,job_id)
);
create trigger whatsapp_schedule_result_immutable before update or delete on communication.whatsapp_schedule_result
 for each row execute function catalog.release_immutable();
create or replace view communication.whatsapp_delivery_approval as
 select g.tenant_id,g.channel_code,g.delivery_key,g.organization_id,g.appointment_id,g.confirmation_event_id::text as approval_event_id,g.lead_id,g.external_id_hmac,g.profile_sha256,g.message_sha256,g.expires_at
 from communication.whatsapp_appointment_approval g
 union all
 select r.tenant_id,'whatsapp',r.request_id,r.organization_id,''::text,r.request_id,
 r.payload->>'lead_id',r.payload->>'external_id_hmac',r.payload->>'profile_sha256',r.payload->>'message_sha256',(r.payload->>'expires_at')::timestamptz
 from approval.request r
 where r.kind='whatsapp_reply' and r.state='approved' and r.payload->>'schema'='elite-whatsapp-reply-approval/v1'
 and r.payload->>'organization_id'=r.organization_id
 union all
 select r.tenant_id,'whatsapp',r.request_id,r.organization_id,''::text,r.request_id,
 r.payload->>'lead_id',r.payload->>'external_id_hmac',r.payload->>'profile_sha256',
 r.payload->>'message_sha256',(r.payload->>'expires_at')::timestamptz
 from approval.request r join communication.whatsapp_schedule s on s.tenant_id=r.tenant_id and s.delivery_key=r.request_id
 where r.kind='whatsapp_schedule' and r.state='approved' and r.payload->>'schema'='elite-whatsapp-schedule/v1'
 and r.payload->>'organization_id'=r.organization_id;
commit;
````

### FILE: `docs/SCHEDULED_COMMUNICATIONS_CONTRACT_PLAN.md`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "15169d28429af5c03c3a29caba4c516b7d31adae39916a6e7edbf69fe036dc5f"
variables: []
secrets_allowed: false
```

````markdown
# Scheduled communications composition plan

Historical plan before revision307. Appointment scheduling is now PROVEN_LOCAL; see SCHEDULED_WHATSAPP_REFERENCE.md. The campaign extension and its separate proof are described in CAMPAIGN_CONNECTED_REFERENCE.md. Original source owners remain fixed by their admission receipts.
First connected claim: appointment reminder prepared from current original
CRM appointment/version, lead/org, contact HMAC/binding and latest explicit
consent/policy. Distinct manual approval binds exact template request and time.
The original platform.job.available_at and finite claim generation schedule it;
no in-memory reminder core or new clock/scheduler algorithm is adopted.

Immutable schedule references the original approval and one original job.
Cancellation is append-only. Rescheduling creates a newly reviewed source/time
binding; it never mutates an existing approval or resets attempts. Before a
provider attempt, recheck appointment state/version/start, tenant/org/contact/
consent and cancellation. Suppress ineligible work without claiming delivery.
One durable outbound fence owns the attempt. UNKNOWN is never retried by a job.
Provider/API lost replies reuse original acceptance-receipt recovery and
signed status/inbox projection. A cancellation cannot recall an in-flight send.

Keep the existing WhatsApp runtime/secrets/profile, authenticated service
identity, sender and verified status worker. Add only optional scheduling
configuration/worker/module to that host. Disabled composition stays compatible.
Source and job parameters remain bounded; worker failures are reported and
readable. Failed or suppressed jobs never make a false delivery/conversion claim.

Then extend this same source-bound owner for a bounded explicit campaign audience
snapshot, reviewed drip times/templates, current opt-out/consent suppression,
and conversion observations from original quotation_acceptance/order receipts.
No invented attribution model, finance or autonomous consent. Existing selected
lead adapters and Page publishing remain intact. Provider credentials are future
user input, never requested or needed for fixture proof.
````

### FILE: `docs/SCHEDULED_WHATSAPP_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "12216f030e6ac9671a76dc6df3dd5db1c7977d7fd71981940df1c18e088be50a"
variables: []
secrets_allowed: false
```

````markdown
# Scheduled WhatsApp reference infrastructure

Scope: local/fixture proof of appointment reminders on the same CRM, identity,
approval, jobs, provider adapter, delivery fence and status host. No live account
is configured by this composition. Campaign audience/drip attribution is a
separate connected T2805 claim.

Materialize the selected profile and apply migrations through
0083_scheduled_whatsapp.up.sql. Configure the existing WhatsApp host from its
hash-locked profile, then set WHATSAPP_ENABLED=true and
WHATSAPP_SCHEDULE_ENABLED=true. The schedule flag is effective only while the
WhatsApp host is enabled. The service identity additionally needs
notification:dispatch in its configured tenant and organization. No new secret
file or provider credential mechanism is introduced.

The host refuses a missing selected schedule implementation, a nonexact flag,
an unprivileged worker or missing/disabled immutable schema guards. It starts
both the original status worker and the scheduled worker, serializes reports
over the original authenticated mTLS connection and joins both on shutdown.

The operator API under /v1/franchise/notifications/scheduled provides:

- POST /prepare with request_id, appointment_id, recipient, template_name,
  language_code, body_parameters, not_before and expires_at.
- POST /requests with the exact prepared context, followed by
  POST /requests/{delivery_key}/decision by a distinct authorized reviewer.
- GET /requests/{delivery_key} for approval, durable job and provider status.
- POST /requests/{delivery_key}/cancel for an unsent scheduled job.
- POST /requests/{delivery_key}/reconcile for a preserved acceptance receipt.

Use notification:request, notification:approve, notification:read,
notification:cancel and notification:reconcile for their respective operations;
reconciliation also requires read permission. IDs, hashes, tenant, organization,
appointment version/time, contact binding, latest consent, approved template and
provider profile remain bound to the reviewed context. JSON is limited to 32 KiB
with duplicate/case-collision and unknown-field rejection. Body reads have a
5-second deadline and requests a 25-second deadline.

Approval atomically records one durable job, due in at most 30 days. Expiry must
follow the due time by at most 24 hours and cannot exceed the appointment start.
Dispatch uses the original lease-generation worker with three attempts and one
job per poll. Source drift, cancellation or withdrawn consent suppress an
unclaimed provider effect. A rescheduled appointment requires a new prepared
context and distinct review; the old approval is never edited.

The original fence serializes the first effect. A source change after its claim
is treated conservatively as uncertain. No cancellation or consent update can
recall a provider effect already accepted. An uncertain effect requires review;
it is never retried as a fresh message. Exact text and template receipts use the
same original payload builders and recover acceptance without network traffic.
The immutable historical worker result and the current delivery-fence state
are shown separately. Provider acceptance does not mean delivery; signed
observations remain observed_delivered or observed_read.

Expired final leases are quarantined as LEASE_EXHAUSTED with terminal job
evidence, not a false completion. Approval, schedule, cancellation and result
records are immutable. Retention policy and archival for this data must be
included in the deployment's operational policy; these tests do not authorize
deleting evidence. A populated downgrade is refused. Empty down/up is verified.

Local verification: TestScheduledWhatsAppConnected uses PostgreSQL 18.6, the
actual pinned Python adapter and a loopback provider fixture; seven approvals
and jobs, three provider POSTs, six immutable results, twelve concurrent
dispatchers, cancellation, CRM reschedule, consent withdrawal, signed status
and lost-stdout recovery. TestScheduledWhatsAppHost proves actual activation,
both loops, mTLS and shutdown. Boundary regressions and a finite fuzz profile
exercise the new input parser. Original official Meta bytes and source lock are
unchanged; the new scheduling and recovery composition is AUTHORED glue.

Later user setup uses the existing WhatsApp account/phone, approved templates,
provider credentials and identity configuration. The materialized code contains
no account values, access tokens or private keys.
````

### FILE: `internal/platform/postgres/scheduled_outbound_binding.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7fec26c3a9562f2d04c03afeb927f9d5472ef94dd868be27f105161de1288c46"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED exported composition point for the existing transactional fence.
// A missing domain admission callback is never silently accepted.
import (
	"context"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"github.com/jackc/pgx/v5"
)

func (s *OutboundDeliveryStore) ClaimWithDomainAdmission(ctx context.Context, m channels.Message, h string, admit func(context.Context, pgx.Tx) error) (outbounddelivery.Claim, error) {
	if admit == nil {
		return outbounddelivery.Claim{}, outbounddelivery.ErrInvalid
	}
	return s.claimWithAdmission(ctx, m, h, admit)
}
````

### FILE: `internal/whatsappbridge/scheduled_boundary_test.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0a162c0c2e502a358d6e9d62da1b88bfeb1b37e4e6aaba38dc902cb9d9ed3b7e"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED regressions of the new HTTP boundary; no database or provider calls.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type scheduledBoundaryVerifier struct {
	principal identity.Principal
	err       error
}

func (v scheduledBoundaryVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, v.err
}

type scheduledDecodeWriter struct{ *httptest.ResponseRecorder }

func (scheduledDecodeWriter) SetReadDeadline(time.Time) error { return nil }

func TestScheduledRequestBoundary(t *testing.T) {
	for _, raw := range []string{
		`{"request_id":"one","REQUEST_ID":"two"}`,
		`{"extra":true}`, `{}{}`, `{"not_before":"invalid"}`,
		strings.Repeat(" ", 32769), `{"request":{"recipient":"one","Recipient":"two"}}`,
	} {
		r := httptest.NewRequest("POST", "/", strings.NewReader(raw))
		r.Header.Set("Content-Type", "application/json")
		if decodeSchedule(scheduledDecodeWriter{httptest.NewRecorder()}, r, new(ScheduleRequest)) == nil {
			t.Fatalf("accepted ambiguous or oversized request %q", raw[:min(len(raw), 80)])
		}
	}
	for _, headers := range []struct{ content, url string }{{"text/plain", "/"}, {"application/json", "/?override=1"}} {
		r := httptest.NewRequest("POST", headers.url, strings.NewReader("{}"))
		r.Header.Set("Content-Type", headers.content)
		if decodeSchedule(scheduledDecodeWriter{httptest.NewRecorder()}, r, new(ScheduleRequest)) == nil {
			t.Fatal("unsupported request metadata")
		}
	}
}

func TestScheduledAuthorizationBoundary(t *testing.T) {
	approved := identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{"notification:request": {}}, Organizations: map[string]struct{}{"store": {}}}
	for _, tc := range []struct {
		name string
		p    identity.Principal
		auth []string
		err  error
		want int
	}{
		{"missing bearer", approved, nil, nil, 401},
		{"duplicate bearer", approved, []string{"Bearer one", "Bearer two"}, nil, 401},
		{"malformed bearer", approved, []string{"Bearer one extra"}, nil, 401},
		{"rejected token", approved, []string{"Bearer one"}, errors.New("fixture"), 401},
		{"other tenant", identity.Principal{Subject: "operator", TenantID: "other", Permissions: approved.Permissions, Organizations: approved.Organizations}, []string{"Bearer one"}, nil, 403},
		{"other org", identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: approved.Permissions, Organizations: map[string]struct{}{"other": {}}}, []string{"Bearer one"}, nil, 403},
		{"no permission", identity.Principal{Subject: "operator", TenantID: "tenant", Organizations: approved.Organizations}, []string{"Bearer one"}, nil, 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := &ScheduledNotifications{approvals: &ScheduleApprovals{tenant: "tenant", org: "store"}}
			mux := http.NewServeMux()
			m.Register(mux, scheduledBoundaryVerifier{tc.p, tc.err})
			r := httptest.NewRequest("POST", "/v1/franchise/notifications/scheduled/prepare", strings.NewReader("{}"))
			r.Header["Authorization"] = tc.auth
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("status %d want %d", w.Code, tc.want)
			}
		})
	}
}

func FuzzScheduledRequestBoundary(f *testing.F) {
	for _, seed := range []string{`{}`, `{"request_id":"schedule-0000001","appointment_id":"appointment-1","recipient":"5491112345678","template_name":"order_update","language_code":"es_AR","body_parameters":["A-1","recordatorio"],"not_before":"2026-09-15T12:00:00Z","expires_at":"2026-09-15T13:00:00Z"}`, `{"request_id":"one","REQUEST_ID":"two"}`, `{"body_parameters":["\ud800"]}`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32769 {
			t.Skip()
		}
		decode := func(body []byte) (ScheduleRequest, error) {
			r := httptest.NewRequest("POST", "/", strings.NewReader(string(body)))
			r.Header.Set("Content-Type", "application/json")
			var request ScheduleRequest
			err := decodeSchedule(scheduledDecodeWriter{httptest.NewRecorder()}, r, &request)
			return request, err
		}
		first, e := decode(raw)
		if e != nil {
			return
		}
		canonical, e := json.Marshal(first)
		if e != nil {
			t.Fatal(e)
		}
		second, e := decode(canonical)
		if e != nil || !reflect.DeepEqual(first, second) {
			t.Fatal("accepted request cannot round-trip through its exact contract")
		}
	})
}
````

### FILE: `internal/whatsappbridge/scheduled_connected_test.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0fb9bd3e2f8ddf6e1fecadbbb5f595e54a5cd4ecd64ce053e7bf299bdb858181"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `internal/whatsappbridge/scheduled_contract.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dbd51921dbf76e3634490fe7d85d8eb05fa30dd1e431739da6faead479448e4f"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED source/time/request binding over the existing CRM and template owner.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

// Optional source resolver shares the original approval/job/fence machinery.
type scheduledSourceExtension interface {
	Source(context.Context, pgx.Tx, ScheduleRequest, bool) (ScheduleContext, error)
	Enqueue(context.Context, pgx.Tx, ScheduleContext, string) error
	Admit(context.Context, pgx.Tx, ScheduleContext) error
}
type ScheduleRequest struct {
	CampaignID     string    `json:"campaign_id,omitempty"`
	LeadID         string    `json:"lead_id,omitempty"`
	Step           int       `json:"step,omitempty"`
	RequestID      string    `json:"request_id"`
	AppointmentID  string    `json:"appointment_id"`
	Recipient      string    `json:"recipient"`
	TemplateName   string    `json:"template_name"`
	LanguageCode   string    `json:"language_code"`
	BodyParameters []string  `json:"body_parameters"`
	NotBefore      time.Time `json:"not_before"`
	ExpiresAt      time.Time `json:"expires_at"`
}
type ScheduleContext struct {
	SourceSHA256          string           `json:"source_sha256,omitempty"`
	MemberSHA256          string           `json:"member_sha256,omitempty"`
	Schema                string           `json:"schema"`
	Request               ScheduleRequest  `json:"request"`
	DeliveryKey           string           `json:"delivery_key"`
	OrganizationID        string           `json:"organization_id"`
	ConnectionID          string           `json:"connection_id"`
	AppointmentID         string           `json:"appointment_id"`
	AppointmentVersion    int64            `json:"appointment_version,string"`
	AppointmentStartsAt   time.Time        `json:"appointment_starts_at"`
	LeadID                string           `json:"lead_id"`
	SubjectID             string           `json:"subject_id"`
	ExternalIDHMAC        string           `json:"external_id_hmac"`
	BindingVersion        int64            `json:"binding_version,string"`
	PolicyVersion         string           `json:"policy_version"`
	ConsentID             string           `json:"consent_id"`
	ConsentPurpose        string           `json:"consent_purpose"`
	ConsentEvidenceSHA256 string           `json:"consent_evidence_sha256"`
	ProfileSHA256         string           `json:"profile_sha256"`
	MessageSHA256         string           `json:"message_sha256"`
	NotBefore             time.Time        `json:"not_before"`
	ExpiresAt             time.Time        `json:"expires_at"`
	Message               channels.Message `json:"message"`
}
type ScheduleApprovals struct {
	extension               scheduledSourceExtension
	base                    *PostgresAppointmentApprovals
	tenant, org, connection string
	profile                 json.RawMessage
	human                   *postgres.HumanApprovals
}

func NewScheduleApprovals(base *PostgresAppointmentApprovals, tenant, org, connection string, profile json.RawMessage) (*ScheduleApprovals, error) {
	if base == nil || base.pool == nil || !notificationEventID.MatchString(tenant) || !cr.ValidID(org) || !webhookConnection.MatchString(connection) || len(profile) > 32768 || !json.Valid(profile) {
		return nil, ErrApproval
	}
	return &ScheduleApprovals{base: base, tenant: tenant, org: org, connection: connection, profile: append(json.RawMessage(nil), profile...), human: postgres.NewHumanApprovals(base.pool)}, nil
}
func (s *ScheduleApprovals) authorized(p identity.Principal, permission string) bool {
	return s != nil && p.TenantID == s.tenant && p.Subject != "" && len(p.Subject) <= 128 && p.Allowed(permission) && p.AllowedOrganization(s.org)
}
func scheduleKey(tenant, id string) string { return "wa-schedule:" + digest([]byte(tenant+"\x00"+id)) }
func validScheduleKey(key string) bool {
	return strings.HasPrefix(key, "wa-schedule:") && validDigest(strings.TrimPrefix(key, "wa-schedule:"))
}
func (s *ScheduleApprovals) source(ctx context.Context, tx pgx.Tx, r ScheduleRequest, preparing bool) (ScheduleContext, error) {
	if r.CampaignID != "" {
		if s.extension == nil || r.AppointmentID != "" {
			return ScheduleContext{}, ErrApproval
		}
		return s.extension.Source(ctx, tx, r, preparing)
	}
	if r.LeadID != "" || r.Step != 0 {
		return ScheduleContext{}, ErrApproval
	}
	var c ScheduleContext
	var now time.Time
	if !cr.ValidID(r.RequestID) || len(r.RequestID) < 16 || !cr.ValidID(r.AppointmentID) || r.BodyParameters == nil || len(r.BodyParameters) > 20 || r.NotBefore.IsZero() || !r.ExpiresAt.After(r.NotBefore) || r.ExpiresAt.Sub(r.NotBefore) > 24*time.Hour {
		return c, ErrApproval
	}
	for _, v := range r.BodyParameters {
		if !cr.ValidText(v, 1024) {
			return c, ErrApproval
		}
	}
	var profile struct {
		Templates []struct {
			Name     string `json:"name"`
			Language string `json:"language_code"`
			Count    int    `json:"body_parameter_count"`
		} `json:"approved_templates"`
	}
	if json.Unmarshal(s.profile, &profile) != nil {
		return c, ErrApproval
	}
	matches := 0
	for _, t := range profile.Templates {
		if t.Name == r.TemplateName && t.Language == r.LanguageCode && t.Count == len(r.BodyParameters) {
			matches++
		}
	}
	if matches != 1 {
		return c, ErrApproval
	}
	h, e := contactidentity.ExternalDigest(s.base.hmacKey, s.tenant, "whatsapp", r.Recipient)
	if e != nil {
		return c, ErrApproval
	}
	e = tx.QueryRow(ctx, `select clock_timestamp(),a.version,a.starts_at,l.lead_id,b.subject_id,b.version,b.policy_version,consent.consent_id,consent.evidence_sha256_hex
 from crm.appointment a join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id and l.organization_id=a.organization_id
 join org.organization o on o.tenant_id=a.tenant_id and o.organization_id=a.organization_id and o.status='active'
 join platform.tenant tenant on tenant.tenant_id=a.tenant_id and tenant.status='active'
 join integration.provider_connection pc on pc.tenant_id=a.tenant_id and pc.connection_id=$4 and pc.organization_id=a.organization_id and pc.provider_code='meta-whatsapp' and pc.state='active'
 join communication.contact_channel_binding b on b.tenant_id=a.tenant_id and b.lead_id=a.lead_id and b.channel_code='whatsapp' and b.external_id_hmac=$5 and b.state='active' and b.pii_allowed and b.effective_at<=statement_timestamp()
 join crm.consent_evidence consent on consent.tenant_id=l.tenant_id and consent.lead_id=l.lead_id and consent.purpose_code=$6 and consent.policy_version=$7 and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 where a.tenant_id=$1 and a.organization_id=$2 and a.appointment_id=$3 and a.state='confirmed' and a.starts_at>statement_timestamp() and b.policy_version=$7
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())
 for share of a,l,b,consent,o,tenant,pc`, s.tenant, s.org, r.AppointmentID, s.connection, h, s.base.purpose, s.base.policy).Scan(&now, &c.AppointmentVersion, &c.AppointmentStartsAt, &c.LeadID, &c.SubjectID, &c.BindingVersion, &c.PolicyVersion, &c.ConsentID, &c.ConsentEvidenceSHA256)
	if errors.Is(e, pgx.ErrNoRows) {
		return c, ErrApproval
	}
	if e != nil {
		return c, e
	}
	if e = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return c, e
	}
	if !r.ExpiresAt.After(now) || r.ExpiresAt.After(c.AppointmentStartsAt) || preparing && !r.NotBefore.After(now) || r.NotBefore.After(now.Add(30*24*time.Hour)) {
		return c, ErrApproval
	}
	body, _ := json.Marshal(struct {
		Recipient      string   `json:"recipient"`
		TemplateName   string   `json:"template_name"`
		LanguageCode   string   `json:"language_code"`
		BodyParameters []string `json:"body_parameters"`
	}{r.Recipient, r.TemplateName, r.LanguageCode, r.BodyParameters})
	c.Schema = "elite-whatsapp-schedule/v1"
	c.Request = r
	c.DeliveryKey = scheduleKey(s.tenant, r.RequestID)
	c.OrganizationID = s.org
	c.ConnectionID = s.connection
	c.AppointmentID = r.AppointmentID
	c.ExternalIDHMAC = h
	c.ConsentPurpose = s.base.purpose
	c.ProfileSHA256 = digest(s.profile)
	c.NotBefore = r.NotBefore
	c.ExpiresAt = r.ExpiresAt
	c.Message = channels.Message{TenantID: s.tenant, ChannelCode: "whatsapp", Direction: channels.DirectionOut, DeliveryKey: c.DeliveryKey, ExternalID: r.Recipient, Text: string(body)}
	c.MessageSHA256, e = outbounddelivery.MessageSHA256(c.Message)
	if e != nil {
		return ScheduleContext{}, ErrApproval
	}
	return c, nil
}
func (s *ScheduleApprovals) Prepare(ctx context.Context, p identity.Principal, r ScheduleRequest) (ScheduleContext, error) {
	if !s.authorized(p, "notification:request") {
		return ScheduleContext{}, ErrApproval
	}
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return ScheduleContext{}, e
	}
	defer tx.Rollback(ctx)
	out, e := s.source(ctx, tx, r, true)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
func (s *ScheduleApprovals) guard(ctx context.Context, tx pgx.Tx, c ScheduleContext) error {
	rebuilt, e := s.source(ctx, tx, c.Request, false)
	if e != nil {
		return e
	}
	raw, _ := json.Marshal(c)
	again, _ := json.Marshal(rebuilt)
	_, h, e := approval.CanonicalPayload(raw)
	_, other, e2 := approval.CanonicalPayload(again)
	if e != nil || e2 != nil || h != other {
		return ErrApproval
	}
	var cancelled bool
	if e = tx.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_schedule_cancellation where tenant_id=$1 and delivery_key=$2)`, s.tenant, c.DeliveryKey).Scan(&cancelled); e != nil {
		return e
	}
	if cancelled {
		return ErrApproval
	}
	return nil
}
func (s *ScheduleApprovals) Submit(ctx context.Context, p identity.Principal, c ScheduleContext) (string, bool, error) {
	if !s.authorized(p, "notification:request") {
		return "", false, ErrApproval
	}
	raw, e := json.Marshal(c)
	if e != nil {
		return "", false, e
	}
	_, h, e := approval.CanonicalPayload(raw)
	if e != nil {
		return "", false, e
	}
	subject := c.AppointmentID
	if c.Request.CampaignID != "" {
		subject = c.LeadID
	}
	replay, e := s.human.Submit(ctx, p, postgres.HumanApprovalSpec{Request: approval.Request{TenantID: s.tenant, ID: c.DeliveryKey, Kind: approval.KindWhatsAppSchedule, SubjectID: subject, Requester: p.Subject, EvidenceSHA: h}, OrganizationID: s.org, Payload: raw}, "notification:request", func(ctx context.Context, tx pgx.Tx) error { return s.guard(ctx, tx, c) })
	return h, replay, e
}
func (s *ScheduleApprovals) request(ctx context.Context, key string) (ScheduleContext, string, error) {
	var c ScheduleContext
	var raw []byte
	var h string
	if !validScheduleKey(key) {
		return c, "", ErrApproval
	}
	e := s.base.pool.QueryRow(ctx, `select payload,evidence_sha from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind='whatsapp_schedule'`, s.tenant, key, s.org).Scan(&raw, &h)
	if e != nil {
		return c, "", ErrApproval
	}
	_, actual, e := approval.CanonicalPayload(raw)
	if e != nil || h != actual || json.Unmarshal(raw, &c) != nil || c.OrganizationID != s.org || c.ConnectionID != s.connection || c.DeliveryKey != key || c.ProfileSHA256 != digest(s.profile) {
		return c, "", ErrApproval
	}
	return c, h, nil
}
````

### FILE: `internal/whatsappbridge/scheduled_exhaustion.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "01be1e469220d46c27f48c94f0e48c8773bbacd3544545613efdf91cccd3d441"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED selection of the existing owner's crashed-final-attempt transition.
// This never resets a lease/attempt or creates a replacement provider effect.
import (
	"context"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (m *ScheduledNotifications) quarantineExpired(ctx context.Context, match json.RawMessage, jobType string) (ScheduledWorkResult, error) {
	var out ScheduledWorkResult
	tx, e := m.approvals.base.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	var job postgres.Job
	e = tx.QueryRow(ctx, `select tenant_id,job_id,queue,job_type,schema_version,payload,attempts,max_attempts from platform.job
 where tenant_id=$1 and queue='whatsapp-scheduled' and job_type=$3 and schema_version=1 and payload @> $2::jsonb
 and completed_at is null and terminal_error_code is null and attempts>=max_attempts and claimed_until<clock_timestamp()
 order by available_at,job_id for update skip locked limit 1`, m.approvals.tenant, match, jobType).Scan(&job.TenantID, &job.JobID, &job.Queue, &job.JobType, &job.SchemaVersion, &job.Payload, &job.Attempts, &job.MaxAttempts)
	if errors.Is(e, pgx.ErrNoRows) {
		return out, nil
	}
	if e != nil {
		return out, e
	}
	if e = postgres.ExhaustJobInTx(ctx, tx, job); e != nil {
		return out, e
	}
	var payload struct {
		DeliveryKey string `json:"delivery_key"`
	}
	if json.Unmarshal(job.Payload, &payload) != nil || !validScheduleKey(payload.DeliveryKey) {
		return out, ErrApproval
	}
	out = ScheduledWorkResult{Claimed: true, DeliveryKey: payload.DeliveryKey, Outcome: "LEASE_EXHAUSTED"}
	return out, tx.Commit(ctx)
}
````

### FILE: `internal/whatsappbridge/scheduled_module.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "61b413388058596ac5b4a3ab6d80f210eefc88abd17dabc40147441d1e865da4"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED bounded operator interface, mounted on the existing role API.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

func decodeSchedule(w http.ResponseWriter, r *http.Request, v any) error {
	if r.Header.Get("Content-Type") != "application/json" || r.URL.RawQuery != "" {
		return ErrApproval
	}
	if http.NewResponseController(w).SetReadDeadline(time.Now().Add(5*time.Second)) != nil {
		return ErrApproval
	}
	raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if e != nil {
		return ErrApproval
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return ErrApproval
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(v) != nil || d.Decode(new(any)) != io.EOF {
		return ErrApproval
	}
	return nil
}
func (m *ScheduledNotifications) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m == nil {
		return
	}
	route := func(method, path, permission string, action func(http.ResponseWriter, *http.Request, identity.Principal)) {
		mux.HandleFunc(method+" /v1/franchise/notifications/scheduled"+path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			header := r.Header.Values("Authorization")
			if verifier == nil || len(header) != 1 || !strings.HasPrefix(header[0], "Bearer ") || len(header[0]) > 16391 || len(strings.Fields(header[0])) != 2 {
				notificationProblem(w, 401, "UNAUTHENTICATED")
				return
			}
			p, e := verifier.Verify(r.Context(), strings.TrimPrefix(header[0], "Bearer "))
			if e != nil {
				notificationProblem(w, 401, "UNAUTHENTICATED")
				return
			}
			if !m.approvals.authorized(p, permission) {
				notificationProblem(w, 403, "FORBIDDEN")
				return
			}
			if r.URL.RawQuery != "" || r.PathValue("key") != "" && !validScheduleKey(r.PathValue("key")) {
				notificationProblem(w, 400, "INVALID_TARGET")
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
			defer cancel()
			action(w, r.WithContext(ctx), p)
		})
	}
	fail := func(w http.ResponseWriter, e error) bool {
		if e == nil {
			return false
		}
		notificationProblem(w, 409, "SCHEDULE_NOT_CURRENT_OR_DIVERGENT")
		return true
	}
	route("POST", "/prepare", "notification:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in ScheduleRequest
		if e := decodeSchedule(w, r, &in); e != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := m.approvals.Prepare(r.Context(), p, in)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/requests", "notification:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in ScheduleContext
		if e := decodeSchedule(w, r, &in); e != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		h, replay, e := m.approvals.Submit(r.Context(), p, in)
		if !fail(w, e) {
			replyJSON(w, map[string]any{"delivery_key": in.DeliveryKey, "request_sha256": h, "replay": replay})
		}
	})
	route("GET", "/requests/{key}", "notification:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		out, e := m.approvals.Status(r.Context(), p, r.PathValue("key"))
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/requests/{key}/decision", "notification:approve", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash    string `json:"request_sha256"`
			Approve *bool  `json:"approve"`
			Reason  string `json:"reason"`
		}
		if decodeSchedule(w, r, &in) != nil || in.Approve == nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := m.approvals.Decide(r.Context(), p, r.PathValue("key"), in.Hash, *in.Approve, in.Reason)
		if !fail(w, e) {
			replyJSON(w, map[string]any{"state": out})
		}
	})
	route("POST", "/requests/{key}/cancel", "notification:cancel", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash   string `json:"request_sha256"`
			Reason string `json:"reason"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		if !fail(w, m.approvals.Cancel(r.Context(), p, r.PathValue("key"), in.Hash, in.Reason)) {
			replyJSON(w, map[string]bool{"cancelled": true})
		}
	})
	route("POST", "/requests/{key}/reconcile", "notification:reconcile", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash string `json:"request_sha256"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := m.Recover(r.Context(), p, r.PathValue("key"), in.Hash)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
}

func (m *ScheduledNotifications) Run(ctx context.Context, source StatusPrincipalSource, reporter StatusReporter, interval time.Duration) error {
	if m == nil || source == nil || reporter == nil || interval < time.Second || interval > time.Minute {
		return ErrStatusHost
	}
	process := func(ctx context.Context, p identity.Principal) (StatusWorkResult, error) {
		v, e := m.ProcessOnce(ctx, p)
		return StatusWorkResult{Claimed: v.Claimed, Completed: v.Outcome != "" && v.Outcome != "LEASE_EXHAUSTED", FailureRecorded: v.Outcome == "RECONCILIATION_REQUIRED" || v.Outcome == "FAILED_TERMINAL" || v.Outcome == "LEASE_EXHAUSTED", Terminal: v.Outcome == "FAILED_TERMINAL" || v.Outcome == "LEASE_EXHAUSTED"}, e
	}
	return runStatusLoop(ctx, process, source, reporter, interval, waitStatusPoll)
}
func (s *ScheduleApprovals) VerifyInfrastructure(ctx context.Context) error {
	var valid bool
	e := s.base.pool.QueryRow(ctx, `select
 (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('communication.whatsapp_schedule') and tgname='whatsapp_schedule_immutable'
 or tgrelid=to_regclass('communication.whatsapp_schedule_cancellation') and tgname='whatsapp_schedule_cancellation_immutable'
 or tgrelid=to_regclass('communication.whatsapp_schedule_result') and tgname='whatsapp_schedule_result_immutable'))=3
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_request()')) like '%whatsapp_schedule%'
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_decision()')) like '%whatsapp_schedule%'
 and pg_get_viewdef(to_regclass('communication.whatsapp_delivery_approval')) like '%whatsapp_schedule%'`).Scan(&valid)
	if e != nil || !valid {
		return ErrApproval
	}
	return nil
}
````

### FILE: `internal/whatsappbridge/scheduled_recovery.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3bd02c1c320ddac7517b69048f5c2f83bc20706dd8e1cdcae6dede2a7c52509a"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

import (
	"context"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"time"
)

func (m *ScheduledNotifications) Recover(ctx context.Context, p identity.Principal, key, hash string) (outbounddelivery.Receipt, error) {
	if m == nil || !m.approvals.authorized(p, "notification:reconcile") {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	value, e := m.approvals.Status(ctx, p, key)
	if e != nil || value.RequestSHA256 != hash || value.ApprovalState != "approved" || value.Notification == nil {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	state := value.Notification.FenceState
	if state != "accepted" && state != "sending" && state != "unknown" {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	var approvedAt time.Time
	var reviewer string
	s := m.approvals
	e = s.base.pool.QueryRow(ctx, `select r.decided_at,d.reviewer from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.approved and d.reviewer<>r.requester where r.tenant_id=$1 and r.request_id=$2 and r.organization_id=$3 and r.kind='whatsapp_schedule' and r.state='approved' and (select count(*) from approval.decision d2 where d2.tenant_id=r.tenant_id and d2.request_id=r.request_id)=1`, s.tenant, key, s.org).Scan(&approvedAt, &reviewer)
	if e != nil || reviewer == "" {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	c := value.Context
	if c.NotBefore.After(approvedAt) {
		approvedAt = c.NotBefore
	}
	return recoverBoundProviderReceipt(ctx, m.sender, m.store, s.base.pool, s.tenant, key, boundReceiptContext{c.Message, c.MessageSHA256, c.ExpiresAt}, approvedAt, state)
}
````

### FILE: `internal/whatsappbridge/scheduled_store.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8796cc331af44d0200ba173f786691ad8970f49487bcefec7ca9f6a8b37a271e"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED transaction wiring. Approval decisions and queue ownership remain
// in their original owners; this module binds their exact source and timing.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"time"
)

func (s *ScheduleApprovals) Decide(ctx context.Context, p identity.Principal, key, hash string, yes bool, reason string) (approval.State, error) {
	if !s.authorized(p, "notification:approve") {
		return "", ErrApproval
	}
	c, h, e := s.request(ctx, key)
	if e != nil || h != hash {
		return "", ErrApproval
	}
	guard := func(ctx context.Context, tx pgx.Tx) error {
		if !yes {
			return nil
		}
		if e := s.guard(ctx, tx, c); e != nil {
			return e
		}
		var now time.Time
		if e := tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil || !c.NotBefore.After(now) {
			return ErrApproval
		}
		var e error
		if c.Request.CampaignID != "" {
			if s.extension == nil {
				return ErrApproval
			}
			e = s.extension.Enqueue(ctx, tx, c, h)
			if e != nil {
				return e
			}
		} else {
			_, e = tx.Exec(ctx, `insert into communication.whatsapp_schedule(tenant_id,delivery_key,job_id,organization_id,appointment_id,appointment_version,not_before,expires_at,request_sha256)
 values($1,$2,gen_random_uuid(),$3,$4,$5,$6,$7,$8) on conflict(tenant_id,delivery_key) do nothing`, s.tenant, key, s.org, c.AppointmentID, c.AppointmentVersion, c.NotBefore, c.ExpiresAt, h)
			if e != nil {
				return e
			}
		}
		jobType := "whatsapp.scheduled.v1"
		if c.Request.CampaignID != "" {
			jobType = "whatsapp.campaign.scheduled.v1"
		}
		_, e = tx.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,available_at,max_attempts)
 select tenant_id,job_id,'whatsapp-scheduled',$6::text,1,
 jsonb_build_object('delivery_key',delivery_key,'connection_id',$3::text,'profile_sha256',$4::text),
 not_before,3 from communication.whatsapp_schedule
 where tenant_id=$1 and delivery_key=$2 and request_sha256=$5 on conflict(tenant_id,job_id) do nothing`, s.tenant, key, s.connection, digest(s.profile), h, jobType)
		return e
	}
	return s.human.Decide(ctx, p, s.tenant, key, s.org, hash, yes, reason, "notification:approve", guard)
}
func (s *ScheduleApprovals) Cancel(ctx context.Context, p identity.Principal, key, hash, reason string) error {
	if !s.authorized(p, "notification:cancel") || reason == "" || len(reason) > 2048 {
		return ErrApproval
	}
	_, h, e := s.request(ctx, key)
	if e != nil || h != hash {
		return ErrApproval
	}
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	var found string
	e = tx.QueryRow(ctx, `select request_sha256 from communication.whatsapp_schedule where tenant_id=$1 and delivery_key=$2 for update`, s.tenant, key).Scan(&found)
	if e != nil || found != h {
		return ErrApproval
	}
	var effect bool
	e = tx.QueryRow(ctx, `select exists(select 1 from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2 and channel_code='whatsapp' and state in('sending','unknown','accepted'))`, s.tenant, key).Scan(&effect)
	if e != nil || effect {
		return ErrApproval
	}
	_, e = tx.Exec(ctx, `insert into communication.whatsapp_schedule_cancellation(tenant_id,delivery_key,actor_subject,reason,request_sha256)values($1,$2,$3,$4,$5)on conflict do nothing`, s.tenant, key, p.Subject, reason, h)
	if e != nil {
		return e
	}
	var same bool
	e = tx.QueryRow(ctx, `select actor_subject=$3 and reason=$4 and request_sha256=$5 from communication.whatsapp_schedule_cancellation where tenant_id=$1 and delivery_key=$2`, s.tenant, key, p.Subject, reason, h).Scan(&same)
	if e != nil || !same {
		return ErrApproval
	}
	return tx.Commit(ctx)
}

type ScheduleStatus struct {
	Notification  *NotificationStatus `json:"notification,omitempty"`
	Context       ScheduleContext     `json:"context"`
	RequestSHA256 string              `json:"request_sha256"`
	ApprovalState string              `json:"approval_state"`
	DeliveryState string              `json:"delivery_state"`
	Cancelled     bool                `json:"cancelled"`
	Outcome       string              `json:"outcome"`
	Attempts      int                 `json:"attempts"`
	JobCompleted  bool                `json:"job_completed"`
	TerminalError string              `json:"terminal_error"`
}

func (s *ScheduleApprovals) Status(ctx context.Context, p identity.Principal, key string) (ScheduleStatus, error) {
	var out ScheduleStatus
	if !s.authorized(p, "notification:read") {
		return out, ErrApproval
	}
	c, h, e := s.request(ctx, key)
	if e != nil {
		return out, e
	}
	out.Context = c
	out.RequestSHA256 = h
	e = s.base.pool.QueryRow(ctx, `select r.state,coalesce(d.state,''),cancel.delivery_key is not null,coalesce(result.outcome,''),
 coalesce(j.attempts,0),j.completed_at is not null,coalesce(j.terminal_error_code,'')
 from approval.request r left join communication.whatsapp_schedule s on s.tenant_id=r.tenant_id and s.delivery_key=r.request_id
 left join communication.whatsapp_schedule_cancellation cancel on cancel.tenant_id=s.tenant_id and cancel.delivery_key=s.delivery_key
 left join communication.whatsapp_schedule_result result on result.tenant_id=s.tenant_id and result.delivery_key=s.delivery_key
 left join platform.job j on j.tenant_id=s.tenant_id and j.job_id=s.job_id
 left join communication.outbound_delivery d on d.tenant_id=r.tenant_id and d.delivery_key=r.request_id and d.channel_code='whatsapp'
 where r.tenant_id=$1 and r.request_id=$2 and r.organization_id=$3 and r.kind='whatsapp_schedule'`, s.tenant, key, s.org).Scan(&out.ApprovalState, &out.DeliveryState, &out.Cancelled, &out.Outcome, &out.Attempts, &out.JobCompleted, &out.TerminalError)
	if e != nil {
		return out, e
	}
	if out.ApprovalState == "approved" {
		status, err := readNotificationStatus(ctx, s.base.pool, p, s.org, "", key)
		if err != nil {
			return out, err
		}
		out.Notification = &status
	}
	return out, nil
}
func (s *ScheduleApprovals) ResolveWhatsAppApproval(ctx context.Context, tenant, key string) (Approval, error) {
	var out Approval
	if tenant != s.tenant {
		return out, ErrApproval
	}
	c, h, e := s.request(ctx, key)
	if e != nil {
		return out, e
	}
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	if e = s.admit(ctx, tx, c, h); e != nil {
		return out, e
	}
	e = tx.QueryRow(ctx, `select d.reviewer from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.approved and d.reviewer<>r.requester where r.tenant_id=$1 and r.request_id=$2 and r.state='approved'`, tenant, key).Scan(&out.ApprovedBy)
	if e != nil {
		return out, e
	}
	out.MessageSHA256 = c.MessageSHA256
	out.ProfileSHA256 = c.ProfileSHA256
	out.EvidenceSHA256 = h
	out.NotBefore = c.NotBefore
	out.ExpiresAt = c.ExpiresAt
	return out, tx.Commit(ctx)
}

// Keep the canonical approval payload check available to the bounded HTTP
// wrapper; unknown/duplicate fields are rejected by its strict typed decoder.
func scheduleCanonical(v any) ([]byte, string, error) {
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, "", e
	}
	return approval.CanonicalPayload(raw)
}
````

### FILE: `internal/whatsappbridge/scheduled_worker.go`

```yaml
block_id: "GO-CONNECTED-SCHEDULED-WHATSAPP:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2d499a8df1d170d88ebe961ccfcfd472fabfc39aac45e521f4041dee981998ac"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED composition of the existing job claim generation and send fence.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"sync/atomic"
	"time"
)

func (s *ScheduleApprovals) admit(ctx context.Context, tx pgx.Tx, c ScheduleContext, h string) error {
	var stored string
	e := tx.QueryRow(ctx, `select request_sha256 from communication.whatsapp_schedule where tenant_id=$1 and delivery_key=$2 and organization_id=$3 for update`, s.tenant, c.DeliveryKey, s.org).Scan(&stored)
	if errors.Is(e, pgx.ErrNoRows) {
		return ErrApproval
	}
	if e != nil {
		return e
	}
	if stored != h {
		return ErrApproval
	}
	if e = s.guard(ctx, tx, c); e != nil {
		return e
	}
	if _, e = postgres.LookupHumanApproval(ctx, tx, s.tenant, c.DeliveryKey, s.org, approval.KindWhatsAppSchedule, h); e != nil {
		return e
	}
	var now time.Time
	if e = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return e
	}
	if c.NotBefore.After(now) || !c.ExpiresAt.After(now) {
		return ErrApproval
	}
	if c.Request.CampaignID != "" {
		if s.extension == nil {
			return ErrApproval
		}
		return s.extension.Admit(ctx, tx, c)
	}
	var ambiguous bool
	e = tx.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_schedule s
 join communication.outbound_delivery d on d.tenant_id=s.tenant_id and d.delivery_key=s.delivery_key and d.channel_code='whatsapp'
 where s.tenant_id=$1 and s.appointment_id=$2 and s.delivery_key<>$3 and d.state in('unknown','sending'))`, s.tenant, c.AppointmentID, c.DeliveryKey).Scan(&ambiguous)
	if e != nil {
		return e
	}
	if ambiguous {
		return outbounddelivery.ErrUnknown
	}
	return nil
}

type scheduledBound struct {
	*postgres.OutboundDeliveryStore
	approvals *ScheduleApprovals
	c         ScheduleContext
	hash      string
}

func (b *scheduledBound) Claim(ctx context.Context, m channels.Message, h string) (outbounddelivery.Claim, error) {
	if m.DeliveryKey != b.c.DeliveryKey || h != b.c.MessageSHA256 {
		return outbounddelivery.Claim{}, ErrApproval
	}
	return b.OutboundDeliveryStore.ClaimWithDomainAdmission(ctx, m, h, func(ctx context.Context, tx pgx.Tx) error { return b.approvals.admit(ctx, tx, b.c, b.hash) })
}

type ScheduledNotifications struct {
	turn      atomic.Uint32
	approvals *ScheduleApprovals
	sender    *Sender
	store     *postgres.OutboundDeliveryStore
	receiver  channels.Channel
	jobs      *postgres.Jobs
	worker    string
}

func NewScheduledNotifications(s *ScheduleApprovals, sender *Sender, store *postgres.OutboundDeliveryStore, receiver channels.Channel, worker string) (*ScheduledNotifications, error) {
	if s == nil || sender == nil || sender.TenantID != s.tenant || digest(sender.Profile) != digest(s.profile) || sender.Tokens == nil || store == nil || receiver == nil || receiver.Code() != "whatsapp" || !webhookConnection.MatchString(worker) {
		return nil, ErrApproval
	}
	frozen := *sender
	frozen.Profile = append(json.RawMessage(nil), sender.Profile...)
	frozen.Approvals = s
	return &ScheduledNotifications{approvals: s, sender: &frozen, store: store, receiver: receiver, jobs: postgres.NewJobs(s.base.pool), worker: worker}, nil
}

type ScheduledWorkResult struct {
	Claimed     bool   `json:"claimed"`
	DeliveryKey string `json:"delivery_key"`
	Outcome     string `json:"outcome"`
}

func (m *ScheduledNotifications) ProcessOnce(ctx context.Context, p identity.Principal) (ScheduledWorkResult, error) {
	var out ScheduledWorkResult
	if m == nil || !m.approvals.authorized(p, "notification:dispatch") {
		return out, ErrApproval
	}
	s := m.approvals
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	match, _ := json.Marshal(map[string]string{"connection_id": s.connection, "profile_sha256": digest(s.profile)})
	kinds := []string{"whatsapp.scheduled.v1"}
	if s.extension != nil {
		kinds = append(kinds, "whatsapp.campaign.scheduled.v1")
		if m.turn.Add(1)%2 == 0 {
			kinds[0], kinds[1] = kinds[1], kinds[0]
		}
	}
	var jobs []postgres.Job
	var e error
	for _, kind := range kinds {
		if expired, err := m.quarantineExpired(ctx, match, kind); err != nil || expired.Claimed {
			return expired, err
		}
		jobs, e = m.jobs.ClaimScoped(ctx, "whatsapp-scheduled", m.worker, 2*time.Minute, 1, postgres.JobScope{TenantID: s.tenant, JobType: kind, SchemaVersion: 1, PayloadMatch: match})
		if e != nil {
			return out, e
		}
		if len(jobs) > 0 {
			break
		}
	}
	if len(jobs) == 0 {
		return out, nil
	}

	job := jobs[0]
	out.Claimed = true
	var payload struct {
		DeliveryKey   string `json:"delivery_key"`
		ConnectionID  string `json:"connection_id"`
		ProfileSHA256 string `json:"profile_sha256"`
	}
	if json.Unmarshal(job.Payload, &payload) != nil || payload.ConnectionID != s.connection || payload.ProfileSHA256 != digest(s.profile) {
		_, e = m.jobs.Fail(ctx, job.TenantID, job.JobID, m.worker, "SCHEDULE_JOB_BINDING", job.Attempts, time.Second)
		return out, e
	}
	out.DeliveryKey = payload.DeliveryKey
	c, h, e := s.request(ctx, payload.DeliveryKey)
	if e != nil {
		_, failed := m.jobs.Fail(ctx, job.TenantID, job.JobID, m.worker, "SCHEDULE_REQUEST_UNAVAILABLE", job.Attempts, time.Second)
		if failed != nil {
			return out, failed
		}
		return out, e
	}
	var same bool
	e = s.base.pool.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_schedule where tenant_id=$1 and delivery_key=$2 and job_id=$3 and request_sha256=$4)`, s.tenant, c.DeliveryKey, job.JobID, h).Scan(&same)
	if e != nil || !same {
		_, failed := m.jobs.Fail(ctx, job.TenantID, job.JobID, m.worker, "SCHEDULE_JOB_DIVERGENT", job.Attempts, time.Second)
		if failed != nil {
			return out, failed
		}
		return out, ErrApproval
	}
	bound := &scheduledBound{OutboundDeliveryStore: m.store, approvals: s, c: c, hash: h}
	channel := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: m.receiver, Sender: m.sender, Store: bound}
	sendErr := channel.Send(ctx, c.Message)
	var state string
	e = s.base.pool.QueryRow(ctx, `select coalesce((select state from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2 and channel_code='whatsapp'),'')`, s.tenant, c.DeliveryKey).Scan(&state)
	if e != nil {
		return out, e
	}
	switch state {
	case "accepted":
		out.Outcome = "ACCEPTED"
	case "failed_terminal":
		out.Outcome = "FAILED_TERMINAL"
	case "unknown", "sending":
		out.Outcome = "RECONCILIATION_REQUIRED"
	default:
		if errors.Is(sendErr, ErrApproval) || errors.Is(sendErr, approval.ErrInvalidRequest) || errors.Is(sendErr, approval.ErrSeparation) {
			out.Outcome = "SUPPRESSED"
		} else {
			_, e = m.jobs.Fail(ctx, job.TenantID, job.JobID, m.worker, "SCHEDULE_EFFECT_UNCONFIRMED", job.Attempts, time.Second)
			if e != nil {
				return out, e
			}
			return out, sendErr
		}
	}
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	_, e = tx.Exec(ctx, `insert into communication.whatsapp_schedule_result(tenant_id,delivery_key,job_id,outcome)values($1,$2,$3,$4)on conflict do nothing`, s.tenant, c.DeliveryKey, job.JobID, out.Outcome)
	if e != nil {
		return out, e
	}
	if e = postgres.CompleteJobInTx(ctx, tx, job, m.worker); e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
````

## 6. Configuration surface

docs/SCHEDULED_WHATSAPP_REFERENCE.md; WHATSAPP_ENABLED=true and WHATSAPP_SCHEDULE_ENABLED=true, original host profile and credentials, explicit notification:dispatch identity.

## 7. Dependency bill

Original Meta source de70ee908a67026e642aaee3703d20464e2a9466 and license retained. No dependency or upstream update. New files AUTHORED unavoidable composition glue.

## 8. Apply order

Existing human approval, CRM, identity/contact, reliable jobs, PG fence, WhatsApp adapter and connected host. Apply0083; populated downgrade refused, empty down/up verified.

## 9. Verification

Seven approval/jobs, three actual adapter POSTs, six results,12concurrent dispatchers ->1first effect; no early delivery; cancel/source drift/consent suppression; receipt recovery without resend; exhausted lease terminal. Host two loops/mTLS/shutdown; scoped boundaries; unit/vet/build and2s fuzz.

## 10. Reconstruction evidence

reconstruction_evidence/SCHEDULED_COMMUNICATIONS_RELEASE_V402.md/json. Narrow local claim; campaign segmentation/drip/conversion and later ordered controls remain open.


V402 composed delta: Optional campaign source/typed job scope and fixed policy host hook; appointment compatibility proven; CAMPAIGN_CONNECTED_RELEASE_V402.md/json.
