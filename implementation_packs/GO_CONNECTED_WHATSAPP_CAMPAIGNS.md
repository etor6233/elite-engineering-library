# Connected permitted WhatsApp campaigns

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "PROVEN_LOCAL permitted CRM lead audiences -> exact reviewed timed templates -> original approval/jobs/adapter/fence -> stop/consent/source/prior-step/quote-conversion suppression and observed original quote/order. AUTHORED glue; no causal attribution or live certification."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["GO-CONNECTED-SCHEDULED-WHATSAPP", "GO-CONNECTED-WHATSAPP-HOST", "GO-HUMAN-APPROVAL-CORE", "GO-RELIABLE-ASYNC-WORKERS", "GO-PG-OUTBOUND-DELIVERY-FENCE", "PYTHON-META-WHATSAPP-CLOUD-ADAPTER", "GO-ELECTROMOBILITY-APPLICATION"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["whatsapp_cloud/official-source.lock.json"]
verified_at: "2026-09-13"
```

## 2. Applicability

Explicit CRM lead audience/source/lifecycle filters, marketing consent,20members/3nonoverlapping exact template steps; same scheduled WhatsApp runtime.

## 3. Architecture contract

Immutable campaign/member snapshot; original distinct per-message approvals and finite jobs/fence; visible/resumable partial create/review; current-source/consent/stop/prior-effect/conversion gates; quote/order observation without attribution formula.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/campaign_whatsapp.go
CREATE cmd/electromobility-api/campaign_whatsapp_test.go
CREATE config/whatsapp.campaign.reference.json
CREATE db/migrations/0084_whatsapp_campaign.down.sql
CREATE db/migrations/0084_whatsapp_campaign.up.sql
CREATE docs/CAMPAIGN_CONNECTED_CONTRACT_PLAN.md
CREATE docs/CAMPAIGN_CONNECTED_REFERENCE.md
CREATE internal/whatsappbridge/campaign_boundary_test.go
CREATE internal/whatsappbridge/campaign_connected_test.go
CREATE internal/whatsappbridge/campaign_contract.go
CREATE internal/whatsappbridge/campaign_decision_replay_test.go
CREATE internal/whatsappbridge/campaign_module.go
CREATE internal/whatsappbridge/campaign_operations.go
CREATE internal/whatsappbridge/campaign_partial_test.go
CREATE internal/whatsappbridge/campaign_schedule_source.go
CREATE internal/whatsappbridge/campaign_status.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/campaign_whatsapp.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "81eb93f2143376b53a92de95f7121b574a64da122c12c0f4394f54669fe0da18"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED optional fixed-policy configuration; no new token or provider path.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/whatsappbridge"
	"encoding/json"
	"io"
)

func init() { whatsappCampaignFactory = selectedWhatsAppCampaign }
func selectedWhatsAppCampaign(ctx context.Context, schedule *whatsappbridge.ScheduledNotifications, c whatsappHostConfig) (httpapi.EnterpriseModule, error) {
	raw, e := readWhatsAppExact(c.CampaignPolicyFile, c.CampaignPolicySHA256, 4096)
	if e != nil {
		return nil, e
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return nil, errWhatsAppHost
	}
	var policy whatsappbridge.CampaignPolicy
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&policy) != nil || d.Decode(new(any)) != io.EOF {
		return nil, errWhatsAppHost
	}
	module, e := whatsappbridge.NewCampaigns(schedule, policy)
	if e != nil {
		return nil, e
	}
	if e = module.VerifyInfrastructure(ctx); e != nil {
		return nil, e
	}
	return module, nil
}
````

### FILE: `cmd/electromobility-api/campaign_whatsapp_test.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0959ee1ee8fd8fbc0f6adf6ed482b17afed0a973e0534d2c9a66d0e2f2c5448e"
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
func TestCampaignHost(t *testing.T) {
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
	campaignPolicy, _ := json.Marshal(whatsappbridge.CampaignPolicy{Schema: "elite-whatsapp-campaign-policy/v1", ConsentPurpose: "fixture-marketing", PolicyVersion: "policy-1", MaxMembers: 20, MaxSteps: 3})
	c.CampaignPolicyFile = write("campaign-policy.json", campaignPolicy)
	c.CampaignPolicySHA256 = hashFile(c.CampaignPolicyFile)

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
	if err := pool.QueryRow(context.Background(), "select tenant_id::text,organization_id,connection_id from integration.provider_connection where provider_code='meta-whatsapp' and state='active' and tenant_id='a582fdb1-25db-59e7-ace2-6cfbc6e155d8'").Scan(&c.TenantID, &c.OrganizationID, &c.ConnectionID); err != nil {
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
	if h.(*whatsappHost).application == nil || h.(*whatsappHost).application.Conversation == nil || h.(*whatsappHost).worker == nil || h.(*whatsappHost).ingress == nil || len(h.(*whatsappHost).modules) != 4 || h.(*whatsappHost).extraRun == nil {
		t.Fatal("missing existing owner")
	}
	mux := http.NewServeMux()
	h.Register(mux, v)
	for _, path := range []string{"/v1/franchise/whatsapp/replies", "/v1/providers/whatsapp/webhook", "/v1/franchise/notifications/scheduled/prepare", "/v1/franchise/marketing/campaigns/prepare"} {
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

	campaignSaved := whatsappCampaignFactory
	whatsappCampaignFactory = nil
	if bad, err := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] }); err == nil {
		bad.close()
		t.Fatal("missing campaign pack activated")
	}
	whatsappCampaignFactory = campaignSaved
	if _, err := pool.Exec(context.Background(), "alter table communication.whatsapp_campaign disable trigger whatsapp_campaign_immutable"); err != nil {
		t.Fatal(err)
	}
	defer pool.Exec(context.Background(), "alter table communication.whatsapp_campaign enable trigger whatsapp_campaign_immutable")
	if bad, err := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] }); err == nil {
		bad.close()
		t.Fatal("campaign guard disabled")
	}
	if _, err := pool.Exec(context.Background(), "alter table communication.whatsapp_campaign enable trigger whatsapp_campaign_immutable"); err != nil {
		t.Fatal(err)
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
	t.Log("CAMPAIGN_HOST_PASS mounted=4 loops=2 telemetry=mTLS shutdown=joined missing_pack_flag_permission_guard=rejected")
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

### FILE: `config/whatsapp.campaign.reference.json`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b04c4bec7bbb45fb6982932267cfc2c594a8eab0d04e7d5f33d601e0e1f49429"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-whatsapp-campaign-policy/v1",
  "consent_purpose": "reference-marketing-consent",
  "policy_version": "reference-policy-1",
  "max_members": 20,
  "max_steps": 3
}
````

### FILE: `db/migrations/0084_whatsapp_campaign.down.sql`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2f2a1ec1e548a0fec0d0c6f620b36a83b50f9a6d87e44f017fe5576f69ee57c3"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
 if exists(select 1 from communication.whatsapp_campaign) or exists(select 1 from communication.whatsapp_schedule where campaign_id is not null)
 then raise exception 'cannot remove campaign evidence';end if;
end $$;
alter table communication.whatsapp_schedule drop constraint whatsapp_schedule_campaign_member_fk,
 drop constraint whatsapp_schedule_source_check,drop constraint whatsapp_schedule_campaign_step_unique,
 drop column campaign_id,drop column lead_id,drop column campaign_step;
alter table communication.whatsapp_schedule alter column appointment_id set not null;
drop table communication.whatsapp_campaign_stop;
drop table communication.whatsapp_campaign_member;
drop table communication.whatsapp_campaign;
commit;
````

### FILE: `db/migrations/0084_whatsapp_campaign.up.sql`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ed47a46d9a7325bf6859ce725785d70371a356f39e622d3de3dc99cd23c28d98"
variables: []
secrets_allowed: false
```

````sql
begin;
create table communication.whatsapp_campaign(
 tenant_id uuid not null, campaign_id text not null, organization_id text not null,
 creator_subject text not null, campaign_sha256 text not null,
 payload jsonb not null, created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,campaign_id),
 unique(tenant_id,campaign_id,campaign_sha256),
 foreign key(tenant_id,organization_id)references org.organization(tenant_id,organization_id),
 check(length(campaign_id)between 16 and 128),
 check(length(creator_subject)between 1 and 128),
 check(campaign_sha256~'^[0-9a-f]{64}$'),
 check(jsonb_typeof(payload)='object' and octet_length(payload::text)<=131072)
);
create trigger whatsapp_campaign_immutable before update or delete on communication.whatsapp_campaign
 for each row execute function catalog.release_immutable();
create table communication.whatsapp_campaign_member(
 tenant_id uuid not null,campaign_id text not null,lead_id text not null,
 member_sha256 text not null check(member_sha256~'^[0-9a-f]{64}$'),
 primary key(tenant_id,campaign_id,lead_id),
 foreign key(tenant_id,campaign_id)references communication.whatsapp_campaign(tenant_id,campaign_id),
 foreign key(tenant_id,lead_id)references crm.lead(tenant_id,lead_id)
);
create trigger whatsapp_campaign_member_immutable before update or delete on communication.whatsapp_campaign_member
 for each row execute function catalog.release_immutable();
create table communication.whatsapp_campaign_stop(
 tenant_id uuid not null,campaign_id text not null,campaign_sha256 text not null,
 actor_subject text not null,reason text not null,
 stopped_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,campaign_id),
 foreign key(tenant_id,campaign_id,campaign_sha256)references communication.whatsapp_campaign(tenant_id,campaign_id,campaign_sha256),
 check(length(actor_subject)between 1 and 128),check(length(reason)between 1 and 2048)
);
create trigger whatsapp_campaign_stop_immutable before update or delete on communication.whatsapp_campaign_stop
 for each row execute function catalog.release_immutable();
alter table communication.whatsapp_schedule alter column appointment_id drop not null;
alter table communication.whatsapp_schedule add column campaign_id text,add column lead_id text,add column campaign_step integer,
 add constraint whatsapp_schedule_campaign_member_fk foreign key(tenant_id,campaign_id,lead_id)
 references communication.whatsapp_campaign_member(tenant_id,campaign_id,lead_id),
 add constraint whatsapp_schedule_source_check check(
 (appointment_id is not null and appointment_version>0 and campaign_id is null and lead_id is null and campaign_step is null)
 or(appointment_id is null and appointment_version=0 and campaign_id is not null and lead_id is not null and campaign_step between 1 and 3)),
 add constraint whatsapp_schedule_campaign_step_unique unique(tenant_id,campaign_id,lead_id,campaign_step);
commit;
````

### FILE: `docs/CAMPAIGN_CONNECTED_CONTRACT_PLAN.md`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ea11c911521b8289f4b9cbbed2d17d524f2dad0055d4e7b2940b9e45e4ab0cc2"
variables: []
secrets_allowed: false
```

````markdown
# Connected campaign composition plan

Historical implementation plan; the resulting local proof and runtime contract are in CAMPAIGN_CONNECTED_REFERENCE.md and CAMPAIGN_CONNECTED_RELEASE_V402.md/json. Reuse the existing
scheduled WhatsApp worker, original approval/job/fence/provider source, CRM
lead/contact/consent and quotation_acceptance/customer_order. No new marketing
engine, provider scheduler, consent authority or attribution formula.

One immutable campaign captures a bounded explicit audience: source/lifecycle
filters, selected lead IDs and recipient bindings, exact approved template
steps/times and source snapshots. At most20members and3steps. No arbitrary
predicate SQL or dynamic unreviewed recipient expansion. Marketing consent uses
an explicit configured purpose/policy distinct from appointment notification
consent. Future policy/account values are configuration, never inferred from a
synthetic fixture or acquisition of an upstream.

Each member/step gets its own original pending WhatsApp schedule approval.
Create/resume is idempotent; partial preparation is visible. The reviewer
explicitly approves the exact stored campaign hash and its pending member/step
hashes, using the original distinct-actor decision owner. Partial review is
visible/resumable; it never manufactures an approval principal. Original
service workers perform only the effects actually approved.

Dispatch rechecks the immutable campaign snapshot, current source/lifecycle,
contact binding, marketing consent, append-only campaign stop, original quote
acceptance and any preceding step. Follow-ups require the preceding step's
accepted fence; an uncertain or rejected predecessor suppresses that follow-up.
Step windows cannot overlap. No source or consent update recalls an accepted
provider effect. All unknown effects use the existing receipt recovery; no resend.

Conversion is a read-only observation from this campaign member's original
quotation_acceptance/order after an accepted campaign delivery. Report quote,
order, acceptance time and original evidence hash. Never claim marketing caused
the sale, charge twice, alter the order or upload a provider conversion.

Mount the campaign API on the existing authenticated role API and host through
an optional exact policy profile. Reuse the existing worker, mTLS reporter,
shutdown and status endpoint. Apply one schema delta; preserve appointment
schedule compatibility and populated rollback refusal.
````

### FILE: `docs/CAMPAIGN_CONNECTED_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6c60d9f13bff6149dd71b1e72dc530bd107621355bb592bb38c3850c9b2a841d"
variables: []
secrets_allowed: false
```

````markdown
# Connected campaign reference infrastructure

Selected scope: explicit permitted CRM lead audiences, exact reviewed WhatsApp
template steps, current consent/source/stop/conversion suppression, durable
provider evidence and observed quotation/order conversion. The immutable
campaign binds the existing CRM, contact identity, approval, jobs, adapter,
outbound fence and host. All new composition code is AUTHORED.

Apply migrations through0084_whatsapp_campaign.up.sql. Keep the existing
WhatsApp host and scheduled worker enabled. In its exact host profile, set
campaign_policy_file and campaign_policy_sha256 to an absolute policy file and
its SHA-256. A complete reference shape is in
config/whatsapp.campaign.reference.json; its reference purpose/policy values are
configuration examples, never evidence of real consent. The selected campaign
pack must be present. Missing/disabled schema guards prevent activation.

The policy must declare a marketing consent purpose distinct from the host's
appointment-notification purpose, and the actual contact/consent policy version.
The profile caps an audience at20members and a campaign at3steps. No new provider
credential path is introduced. The same existing account, exact approved
templates, identity and credential providers are used when the user configures
them. Local fixtures use only synthetic values and loopback provider traffic.

The authenticated operator API is /v1/franchise/marketing/campaigns:

- POST /prepare: campaign_id, explicit sources/states, members containing
  lead_id/recipient, and exact template steps with nonoverlapping due/expiry
  windows. Returns the current source/consent snapshot.
- POST the exact snapshot to the collection to store it and create the original
  pending schedule approvals. This operation can return complete:false with
  per-item codes; no failure is concealed as a fully prepared campaign.
- POST /{id}/resume with campaign_sha256 to resume preparation as its creator.
- POST /{id}/decision with campaign_sha256, approve and reason as a distinct
  reviewer. It explicitly reviews the stored member/step batch. Partial review
  remains visible and can be resumed.
- GET /{id}: immutable audience, stop state, each original approval/job/fence/
  signed notification status and first observed downstream quote/order.
- POST /{id}/stop with campaign_sha256 and reason stops future eligible work.

Use marketing:request plus notification:request for create/resume;
marketing:approve plus notification:approve for review; marketing:read plus
notification:read for state; marketing:cancel for stop. The service retains
notification:dispatch. Receipt reconciliation uses the original scheduled
notification endpoint and its existing notification:reconcile/read permissions.
Nothing manufactures the creator's or reviewer's identity.

Each member must match the explicitly selected source and one of new/contacted/
qualified, an active contact binding, current marketing consent and configured
organization. Source snapshots are immutable. A lifecycle/source/binding/consent
change suppresses the old reviewed effect and requires a new campaign snapshot.
No recipients are dynamically added after review. Requests/snapshots are bounded
to32KiB; template arity/body limits come from the existing approved profile.
Each due time is within30days, each expiry window within24hours, and steps do not
overlap. A follow-up requires an accepted predecessor; uncertain/rejected/
suppressed predecessors never authorize another provider attempt.

The same original worker claims typed jobs. Campaign jobs use
whatsapp.campaign.scheduled.v1 so a host without the campaign extension leaves
them untouched. When enabled, the existing worker alternates the two admitted
types. It retains the original finite lease generations, one job per poll,
durable fence and exact receipt recovery. A stop or opt-out cannot recall an
accepted/in-flight effect. Unknown effects are not retried as new messages.

An interrupted batch retry reads an already-stored decision only when request
hash, reviewer, reason, decision and its unique immutable record match. It
creates neither another approval nor another job. A changed decision/reviewer/
reason remains unresolved rather than rewriting historical approval.

Conversion is the first original quotation_acceptance/order for this member
after an accepted campaign message. The API reports its original IDs, acceptance
time, order state and evidence hash. This is observed sequence, not causal
attribution, payment proof, attributed revenue or a provider conversion upload.
Existing acceptance suppresses subsequent campaign steps.

Reference proofs: six campaigns,10jobs,4actual Python adapter POSTs,10durable
results and one original quote-to-order conversion;12concurrent workers give
one first effect. Further proof covers partial create/review/resume with2unique
requests/jobs/decisions and zero sends. Four decision-replay negatives are
read-only against those exact receipts. The host mounts both capabilities with
two existing worker loops, mTLS reports and joined shutdown. Populated rollback
is refused; empty down/up is tested. Original appointment behavior has its
specific compatibility regression after this source extension.

Retention/archival for campaign and recipient evidence belongs in the deployment
operational policy and T2809. These fixtures do not authorize deleting immutable
history or certify live delivery, production security or jurisdictional consent.
````

### FILE: `internal/whatsappbridge/campaign_boundary_test.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "719a334a7f25faf937749d1da387271de548cf50e50712acadb586385280ddc9"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED finite command-envelope and source-selection invariants.
import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func campaignBoundaryFixture() (*Campaigns, CampaignRequest) {
	c := &Campaigns{policy: CampaignPolicy{MaxMembers: 20, MaxSteps: 3}, schedule: &ScheduledNotifications{approvals: &ScheduleApprovals{profile: json.RawMessage(`{"approved_templates":[{"name":"order_update","language_code":"es_AR","body_parameter_count":2}]}`)}}}
	at := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	r := CampaignRequest{ID: "campaign-boundary-001", Sources: []string{"fixture"}, States: []string{"new"}, Members: []CampaignRecipient{{LeadID: "lead-1", Recipient: "5491112345678"}}, Steps: []CampaignStep{{TemplateName: "order_update", LanguageCode: "es_AR", BodyParameters: []string{"a", "b"}, NotBefore: at, ExpiresAt: at.Add(time.Hour)}}}
	return c, r
}

type campaignDecodeWriter struct{ *httptest.ResponseRecorder }

func (campaignDecodeWriter) SetReadDeadline(time.Time) error { return nil }
func TestCampaignSelectionBoundary(t *testing.T) {
	c, r := campaignBoundaryFixture()
	if c.validate(r) != nil {
		t.Fatal("valid bounded selection")
	}
	for name, change := range map[string]func(*CampaignRequest){
		"implicit source":     func(r *CampaignRequest) { r.Sources = nil },
		"converted audience":  func(r *CampaignRequest) { r.States = []string{"converted"} },
		"duplicate member":    func(r *CampaignRequest) { r.Members = append(r.Members, r.Members[0]) },
		"unreviewed template": func(r *CampaignRequest) { r.Steps[0].TemplateName = "invented" },
		"wrong arity":         func(r *CampaignRequest) { r.Steps[0].BodyParameters = []string{"only-one"} },
		"overlap":             func(r *CampaignRequest) { r.Steps = append(r.Steps, r.Steps[0]) },
		"unbounded steps":     func(r *CampaignRequest) { r.Steps = make([]CampaignStep, 4) },
		"unbounded audience":  func(r *CampaignRequest) { r.Members = make([]CampaignRecipient, 21) },
	} {
		t.Run(name, func(t *testing.T) {
			_, r := campaignBoundaryFixture()
			change(&r)
			if c.validate(r) == nil {
				t.Fatal("invalid selection admitted")
			}
		})
	}
	for _, raw := range []string{`{"campaign_id":"one","CAMPAIGN_ID":"two"}`, `{"states":["new"],"predicate_sql":"select anything"}`, `{}{}`, strings.Repeat(" ", 32769)} {
		r := httptest.NewRequest("POST", "/", strings.NewReader(raw))
		r.Header.Set("Content-Type", "application/json")
		if decodeSchedule(campaignDecodeWriter{httptest.NewRecorder()}, r, new(CampaignRequest)) == nil {
			t.Fatal("ambiguous command")
		}
	}
	base := &ScheduledNotifications{approvals: &ScheduleApprovals{base: &PostgresAppointmentApprovals{purpose: "appointment"}}}
	policy := CampaignPolicy{Schema: "elite-whatsapp-campaign-policy/v1", ConsentPurpose: "appointment", PolicyVersion: "policy-1", MaxMembers: 20, MaxSteps: 3}
	if _, e := NewCampaigns(base, policy); e == nil {
		t.Fatal("appointment consent reused for marketing")
	}
}
func FuzzCampaignSelection(f *testing.F) {
	_, r := campaignBoundaryFixture()
	raw, _ := json.Marshal(r)
	f.Add(raw)
	f.Add([]byte(`{"campaign_id":"one","CAMPAIGN_ID":"two"}`))
	f.Add([]byte(`{"steps":[{"not_before":"invalid"}]}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32769 {
			t.Skip()
		}
		req := httptest.NewRequest("POST", "/", strings.NewReader(string(raw)))
		req.Header.Set("Content-Type", "application/json")
		var value CampaignRequest
		if decodeSchedule(campaignDecodeWriter{httptest.NewRecorder()}, req, &value) != nil {
			return
		}
		canonical, e := json.Marshal(value)
		if e != nil {
			t.Fatal(e)
		}
		var again CampaignRequest
		if json.Unmarshal(canonical, &again) != nil || !reflect.DeepEqual(value, again) {
			t.Fatal("round-trip lost reviewed command")
		}
		c, _ := campaignBoundaryFixture()
		if c.validate(value) == nil {
			if len(value.Members) < 1 || len(value.Members) > 20 || len(value.Steps) < 1 || len(value.Steps) > 3 {
				t.Fatal("unbounded admitted campaign")
			}
			for i, s := range value.Steps {
				if !s.ExpiresAt.After(s.NotBefore) || i > 0 && !s.NotBefore.After(value.Steps[i-1].ExpiresAt) {
					t.Fatal("overlapping admitted steps")
				}
			}
		}
	})
}
````

### FILE: `internal/whatsappbridge/campaign_connected_test.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "cd3fde8df059637a54c25bc078d7198e88dcb3ed12bd134c748a5fcd2d043257"
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
		code, raw := call("POST", "", maker, plan)
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
	t.Logf("CAMPAIGN_CONNECTED_PASS campaigns=%d approvals=%d jobs=%d fences=%d results=%d provider_posts=%d original_quote_orders=%d concurrent=12 no_causal_attribution tenant=%s", campaignsN, approvalsN, jobsN, fencesN, resultsN, meta.Load(), ordersN, tenant)
}
````

### FILE: `internal/whatsappbridge/campaign_contract.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "66d516419feb852522dfd2b1e5ac8a6d7164acc30787840ac57313d349013d63"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED bounded audience/source binding around existing CRM and consent rows.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"sort"
	"time"
)

type CampaignPolicy struct {
	Schema         string `json:"schema"`
	ConsentPurpose string `json:"consent_purpose"`
	PolicyVersion  string `json:"policy_version"`
	MaxMembers     int    `json:"max_members"`
	MaxSteps       int    `json:"max_steps"`
}
type CampaignRecipient struct {
	LeadID    string `json:"lead_id"`
	Recipient string `json:"recipient"`
}
type CampaignStep struct {
	TemplateName   string    `json:"template_name"`
	LanguageCode   string    `json:"language_code"`
	BodyParameters []string  `json:"body_parameters"`
	NotBefore      time.Time `json:"not_before"`
	ExpiresAt      time.Time `json:"expires_at"`
}
type CampaignRequest struct {
	ID      string              `json:"campaign_id"`
	Sources []string            `json:"sources"`
	States  []string            `json:"states"`
	Members []CampaignRecipient `json:"members"`
	Steps   []CampaignStep      `json:"steps"`
}
type CampaignMember struct {
	LeadID         string    `json:"lead_id"`
	Recipient      string    `json:"recipient"`
	Source         string    `json:"source"`
	Lifecycle      string    `json:"lifecycle"`
	LeadUpdatedAt  time.Time `json:"lead_updated_at"`
	SubjectID      string    `json:"subject_id"`
	ExternalIDHMAC string    `json:"external_id_hmac"`
	BindingVersion int64     `json:"binding_version,string"`
	BindingPolicy  string    `json:"binding_policy"`
	ConsentID      string    `json:"consent_id"`
	ConsentSHA256  string    `json:"consent_sha256"`
}
type CampaignSnapshot struct {
	Schema         string           `json:"schema"`
	Request        CampaignRequest  `json:"request"`
	Requester      string           `json:"requester"`
	TenantID       string           `json:"tenant_id"`
	OrganizationID string           `json:"organization_id"`
	ConnectionID   string           `json:"connection_id"`
	ProfileSHA256  string           `json:"profile_sha256"`
	Policy         CampaignPolicy   `json:"policy"`
	Members        []CampaignMember `json:"members"`
}
type Campaigns struct {
	schedule *ScheduledNotifications
	policy   CampaignPolicy
}

func NewCampaigns(schedule *ScheduledNotifications, policy CampaignPolicy) (*Campaigns, error) {
	if schedule == nil || schedule.approvals == nil || schedule.approvals.extension != nil ||
		policy.Schema != "elite-whatsapp-campaign-policy/v1" || !cr.ValidID(policy.ConsentPurpose) || !cr.ValidID(policy.PolicyVersion) ||
		policy.ConsentPurpose == schedule.approvals.base.purpose || policy.MaxMembers < 1 || policy.MaxMembers > 20 || policy.MaxSteps < 1 || policy.MaxSteps > 3 {
		return nil, ErrApproval
	}
	c := &Campaigns{schedule: schedule, policy: policy}
	schedule.approvals.extension = c
	return c, nil
}
func (c *Campaigns) authorized(p identity.Principal, permission string) bool {
	return c != nil && c.schedule.approvals.authorized(p, permission)
}
func campaignHash(v any) ([]byte, string, error) {
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, "", e
	}
	return approval.CanonicalPayload(raw)
}
func (c *Campaigns) validate(r CampaignRequest) error {
	if !cr.ValidID(r.ID) || len(r.ID) < 16 || len(r.Members) < 1 || len(r.Members) > c.policy.MaxMembers || len(r.Steps) < 1 || len(r.Steps) > c.policy.MaxSteps || len(r.Sources) < 1 || len(r.Sources) > 8 || len(r.States) < 1 || len(r.States) > 3 {
		return ErrApproval
	}
	unique := func(v []string) bool {
		seen := map[string]bool{}
		for _, x := range v {
			if !cr.ValidID(x) || seen[x] {
				return false
			}
			seen[x] = true
		}
		return true
	}
	if !unique(r.Sources) || !unique(r.States) {
		return ErrApproval
	}
	for _, state := range r.States {
		if state != "new" && state != "qualified" && state != "contacted" {
			return ErrApproval
		}
	}
	seen := map[string]bool{}
	for _, m := range r.Members {
		if !cr.ValidID(m.LeadID) || seen[m.LeadID] {
			return ErrApproval
		}
		seen[m.LeadID] = true
	}
	var previous time.Time
	var profile struct {
		Templates []struct {
			Name     string `json:"name"`
			Language string `json:"language_code"`
			Count    int    `json:"body_parameter_count"`
		} `json:"approved_templates"`
	}
	if json.Unmarshal(c.schedule.approvals.profile, &profile) != nil {
		return ErrApproval
	}
	for _, s := range r.Steps {
		if s.NotBefore.IsZero() || !s.ExpiresAt.After(s.NotBefore) || s.ExpiresAt.Sub(s.NotBefore) > 24*time.Hour || !previous.IsZero() && !s.NotBefore.After(previous) || s.BodyParameters == nil || len(s.BodyParameters) > 20 {
			return ErrApproval
		}
		matches := 0
		for _, t := range profile.Templates {
			if t.Name == s.TemplateName && t.Language == s.LanguageCode && t.Count == len(s.BodyParameters) {
				matches++
			}
		}
		if matches != 1 {
			return ErrApproval
		}
		for _, v := range s.BodyParameters {
			if !cr.ValidText(v, 1024) {
				return ErrApproval
			}
		}
		previous = s.ExpiresAt
	}
	return nil
}
func (c *Campaigns) member(ctx context.Context, tx pgx.Tx, recipient CampaignRecipient) (CampaignMember, error) {
	s := c.schedule.approvals
	m := CampaignMember{LeadID: recipient.LeadID, Recipient: recipient.Recipient}
	h, e := contactidentity.ExternalDigest(s.base.hmacKey, s.tenant, "whatsapp", recipient.Recipient)
	if e != nil {
		return m, ErrApproval
	}
	e = tx.QueryRow(ctx, `select l.source_code,l.lifecycle_state,l.updated_at,b.subject_id,b.version,b.policy_version,consent.consent_id,consent.evidence_sha256_hex
 from crm.lead l join org.organization o on o.tenant_id=l.tenant_id and o.organization_id=l.organization_id
 join platform.tenant tenant on tenant.tenant_id=l.tenant_id
 join communication.contact_channel_binding b on b.tenant_id=l.tenant_id and b.lead_id=l.lead_id
 join crm.consent_evidence consent on consent.tenant_id=l.tenant_id and consent.lead_id=l.lead_id
 join integration.provider_connection pc on pc.tenant_id=l.tenant_id and pc.organization_id=l.organization_id
 where l.tenant_id=$1 and l.organization_id=$2 and l.lead_id=$3 and o.status='active' and tenant.status='active'
 and b.channel_code='whatsapp' and b.external_id_hmac=$4 and b.state='active' and b.pii_allowed
 and b.effective_at<=statement_timestamp()
 and b.policy_version=$5 and consent.purpose_code=$6 and consent.policy_version=$5 and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())
 and pc.connection_id=$7 and pc.provider_code='meta-whatsapp' and pc.state='active'
 for share of l,o,tenant,b,consent,pc`, s.tenant, s.org, m.LeadID, h, c.policy.PolicyVersion, c.policy.ConsentPurpose, s.connection).Scan(&m.Source, &m.Lifecycle, &m.LeadUpdatedAt, &m.SubjectID, &m.BindingVersion, &m.BindingPolicy, &m.ConsentID, &m.ConsentSHA256)
	if errors.Is(e, pgx.ErrNoRows) {
		return m, ErrApproval
	}
	m.ExternalIDHMAC = h
	return m, e
}
func contains(v []string, s string) bool {
	for _, x := range v {
		if x == s {
			return true
		}
	}
	return false
}
func (c *Campaigns) prepare(ctx context.Context, tx pgx.Tx, p identity.Principal, r CampaignRequest) (CampaignSnapshot, error) {
	var out CampaignSnapshot
	if c.validate(r) != nil {
		return out, ErrApproval
	}
	// Copy before sorting: a caller's request is never mutated.
	r.Members = append([]CampaignRecipient(nil), r.Members...)
	sort.Slice(r.Members, func(i, j int) bool { return r.Members[i].LeadID < r.Members[j].LeadID })
	r.Sources = append([]string(nil), r.Sources...)
	sort.Strings(r.Sources)
	r.States = append([]string(nil), r.States...)
	sort.Strings(r.States)
	s := c.schedule.approvals
	out = CampaignSnapshot{Schema: "elite-whatsapp-campaign/v1", Request: r, Requester: p.Subject, TenantID: s.tenant, OrganizationID: s.org, ConnectionID: s.connection, ProfileSHA256: digest(s.profile), Policy: c.policy}
	for _, selected := range r.Members {
		m, e := c.member(ctx, tx, selected)
		if e != nil {
			return out, e
		}
		if !contains(r.Sources, m.Source) || !contains(r.States, m.Lifecycle) {
			return out, ErrApproval
		}
		out.Members = append(out.Members, m)
	}
	var now time.Time
	if e := tx.QueryRow(ctx, "select clock_timestamp()").Scan(&now); e != nil {
		return out, e
	}
	for _, step := range r.Steps {
		if !step.NotBefore.After(now) || step.NotBefore.After(now.Add(30*24*time.Hour)) {
			return out, ErrApproval
		}
	}
	if raw, _, e := campaignHash(out); e != nil || len(raw) > 32768 {
		return out, ErrApproval
	}
	return out, nil
}
func (c *Campaigns) Prepare(ctx context.Context, p identity.Principal, r CampaignRequest) (CampaignSnapshot, error) {
	if !c.authorized(p, "marketing:request") {
		return CampaignSnapshot{}, ErrApproval
	}
	tx, e := c.schedule.approvals.base.pool.Begin(ctx)
	if e != nil {
		return CampaignSnapshot{}, e
	}
	defer tx.Rollback(ctx)
	out, e := c.prepare(ctx, tx, p, r)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
````

### FILE: `internal/whatsappbridge/campaign_decision_replay_test.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d7b0b5f76f1aa34ccbdc40b7a66856b57d8dac740d7d017a6471e64747b94696"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED negative probe of the already-created immutable partial-review receipt.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestCampaignDecisionReplayBinding(t *testing.T) {
	dsn, tenant := os.Getenv("ELITE_WHATSAPP_CONNECTED_DATABASE_URL"), os.Getenv("ELITE_CAMPAIGN_PARTIAL_TENANT")
	if dsn == "" || tenant == "" {
		t.Skip("exact owned partial-progress fixture required")
	}
	u, e := url.Parse(dsn)
	if e != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_whatsapp_connected_") {
		t.Fatal("fixture scope")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, dsn)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	var key, hash, reviewer, reason string
	e = pool.QueryRow(ctx, `select r.request_id,r.evidence_sha,d.reviewer,d.reason from approval.request r join approval.decision d using(tenant_id,request_id)
 where r.tenant_id=$1 and r.kind='whatsapp_schedule' and r.payload->'request'->>'campaign_id'='campaign-partial-0001' and r.payload->'request'->>'step'='1' and r.state='approved'`, tenant).Scan(&key, &hash, &reviewer, &reason)
	if e != nil {
		t.Fatal(e)
	}
	c := &Campaigns{schedule: &ScheduledNotifications{approvals: &ScheduleApprovals{tenant: tenant, org: "store-1", base: &PostgresAppointmentApprovals{pool: pool}}}}
	p := identity.Principal{Subject: reviewer, TenantID: tenant}
	if state, e := c.decisionReplay(ctx, p, key, hash, true, reason); e != nil || state != "approved" {
		t.Fatal("exact replay", state, e)
	}
	for _, v := range []struct {
		actor, hash, reason string
		yes                 bool
	}{
		{"another-reviewer", hash, reason, true}, {reviewer, strings.Repeat("f", 64), reason, true}, {reviewer, hash, reason + " changed", true}, {reviewer, hash, reason, false},
	} {
		q := p
		q.Subject = v.actor
		if _, e := c.decisionReplay(ctx, q, key, v.hash, v.yes, v.reason); e == nil {
			t.Fatal("unbound decision replay accepted")
		}
	}
	t.Log("CAMPAIGN_DECISION_REPLAY_BOUNDARY_PASS exact=1 actor_hash_reason_decision_negatives=4 no_writes")
}
````

### FILE: `internal/whatsappbridge/campaign_module.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "578516e7fc7180b6ec17856f8936f1a648f8021e049776fb7fd7fe3455f9fb72"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED explicit operator API for the immutable campaign and original effects.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
	"strings"
	"time"
)

func (c *Campaigns) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if c == nil {
		return
	}
	route := func(method, path, permission string, action func(http.ResponseWriter, *http.Request, identity.Principal)) {
		mux.HandleFunc(method+" /v1/franchise/marketing/campaigns"+path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			headers := r.Header.Values("Authorization")
			if verifier == nil || len(headers) != 1 || !strings.HasPrefix(headers[0], "Bearer ") || len(headers[0]) > 16391 || len(strings.Fields(headers[0])) != 2 {
				notificationProblem(w, 401, "UNAUTHENTICATED")
				return
			}
			p, e := verifier.Verify(r.Context(), strings.TrimPrefix(headers[0], "Bearer "))
			if e != nil {
				notificationProblem(w, 401, "UNAUTHENTICATED")
				return
			}
			if !c.authorized(p, permission) {
				notificationProblem(w, 403, "FORBIDDEN")
				return
			}
			if r.URL.RawQuery != "" || r.PathValue("id") != "" && (!cr.ValidID(r.PathValue("id")) || len(r.PathValue("id")) < 16) {
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
		notificationProblem(w, 409, "CAMPAIGN_UNAVAILABLE_OR_DIVERGENT")
		return true
	}
	route("POST", "/prepare", "marketing:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in CampaignRequest
		if decodeSchedule(w, r, &in) != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := c.Prepare(r.Context(), p, in)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "", "marketing:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in CampaignSnapshot
		if decodeSchedule(w, r, &in) != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := c.Create(r.Context(), p, in)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("GET", "/{id}", "marketing:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		out, e := c.Status(r.Context(), p, r.PathValue("id"))
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/{id}/resume", "marketing:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash string `json:"campaign_sha256"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := c.Resume(r.Context(), p, r.PathValue("id"), in.Hash)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/{id}/decision", "marketing:approve", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash    string `json:"campaign_sha256"`
			Approve *bool  `json:"approve"`
			Reason  string `json:"reason"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) || in.Approve == nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := c.Review(r.Context(), p, r.PathValue("id"), in.Hash, *in.Approve, in.Reason)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/{id}/stop", "marketing:cancel", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash   string `json:"campaign_sha256"`
			Reason string `json:"reason"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		if !fail(w, c.Stop(r.Context(), p, r.PathValue("id"), in.Hash, in.Reason)) {
			replyJSON(w, map[string]bool{"stopped": true})
		}
	})
}
````

### FILE: `internal/whatsappbridge/campaign_operations.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9c541b333bf2df88385c1ce12b574f0d9a22767dc10fc2895c783d329630e1dd"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED explicit, resumable composition of original per-message approvals.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
)

type CampaignItem struct {
	LeadID        string `json:"lead_id"`
	Step          int    `json:"step"`
	DeliveryKey   string `json:"delivery_key"`
	RequestSHA256 string `json:"request_sha256,omitempty"`
	State         string `json:"state"`
	Code          string `json:"code,omitempty"`
}
type CampaignBatch struct {
	CampaignID     string         `json:"campaign_id"`
	CampaignSHA256 string         `json:"campaign_sha256"`
	Complete       bool           `json:"complete"`
	Items          []CampaignItem `json:"items"`
}

func (c *Campaigns) Create(ctx context.Context, p identity.Principal, in CampaignSnapshot) (CampaignBatch, error) {
	var out CampaignBatch
	if !c.authorized(p, "marketing:request") || !c.authorized(p, "notification:request") || in.Requester != p.Subject {
		return out, ErrApproval
	}
	raw, h, e := campaignHash(in)
	if e != nil || len(raw) > 32768 {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	var exists bool
	e = tx.QueryRow(ctx, "select exists(select 1 from communication.whatsapp_campaign where tenant_id=$1 and campaign_id=$2)", s.tenant, in.Request.ID).Scan(&exists)
	if e != nil {
		return out, e
	}
	if !exists {
		fresh, e := c.prepare(ctx, tx, p, in.Request)
		if e != nil {
			return out, e
		}
		_, actual, e := campaignHash(fresh)
		if e != nil || actual != h {
			return out, ErrApproval
		}
		_, e = tx.Exec(ctx, `insert into communication.whatsapp_campaign(tenant_id,campaign_id,organization_id,creator_subject,campaign_sha256,payload)values($1,$2,$3,$4,$5,$6)on conflict(tenant_id,campaign_id)do nothing`, s.tenant, in.Request.ID, s.org, p.Subject, h, json.RawMessage(raw))
		if e != nil {
			return out, e
		}
		stored, actual, e := c.load(ctx, tx, in.Request.ID)
		if e != nil {
			return out, e
		}
		if actual != h || stored.Requester != p.Subject {
			return out, ErrApproval
		}
		for _, member := range fresh.Members {
			_, mh, e := campaignHash(member)
			if e != nil {
				return out, e
			}
			_, e = tx.Exec(ctx, `insert into communication.whatsapp_campaign_member(tenant_id,campaign_id,lead_id,member_sha256)values($1,$2,$3,$4)on conflict do nothing`, s.tenant, in.Request.ID, member.LeadID, mh)
			if e != nil {
				return out, e
			}
		}
	} else {
		stored, actual, e := c.load(ctx, tx, in.Request.ID)
		if e != nil {
			return out, e
		}
		if actual != h || stored.Requester != p.Subject {
			return out, ErrApproval
		}
	}
	if e = tx.Commit(ctx); e != nil {
		return out, e
	}
	return c.Resume(ctx, p, in.Request.ID, h)
}
func (c *Campaigns) Resume(ctx context.Context, p identity.Principal, id, hash string) (CampaignBatch, error) {
	var out CampaignBatch
	if !c.authorized(p, "marketing:request") || !c.authorized(p, "notification:request") || !validDigest(hash) {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	plan, h, e := c.load(ctx, s.base.pool, id)
	if e != nil {
		return out, e
	}
	if hash != h || plan.Requester != p.Subject {
		return out, ErrApproval
	}
	out = CampaignBatch{CampaignID: id, CampaignSHA256: h, Complete: true}
	for _, m := range plan.Members {
		for step := 1; step <= len(plan.Request.Steps); step++ {
			request := campaignScheduleRequest(plan, m, step)
			item := CampaignItem{LeadID: m.LeadID, Step: step, DeliveryKey: scheduleKey(s.tenant, request.RequestID)}
			var exists bool
			if e = s.base.pool.QueryRow(ctx, "select exists(select 1 from approval.request where tenant_id=$1 and request_id=$2)", s.tenant, item.DeliveryKey).Scan(&exists); e != nil {
				return out, e
			}
			if exists {
				bound, bh, err := s.request(ctx, item.DeliveryKey)
				if err != nil || bound.SourceSHA256 != h || bound.Request.CampaignID != id || bound.LeadID != m.LeadID || bound.Request.Step != step {
					item.Code = "EXISTING_REQUEST_UNVERIFIED"
					out.Complete = false
				} else {
					item.RequestSHA256 = bh
					item.State = "PREPARED"
				}
			} else {
				bound, err := s.Prepare(ctx, p, request)
				if err != nil {
					item.Code = "SOURCE_NOT_CURRENT"
					out.Complete = false
				} else {
					bh, _, err := s.Submit(ctx, p, bound)
					if err != nil {
						item.Code = "PREPARATION_UNCONFIRMED"
						out.Complete = false
					} else {
						item.RequestSHA256 = bh
						item.State = "PREPARED"
					}
				}
			}
			out.Items = append(out.Items, item)
		}
	}
	return out, nil
}
func (c *Campaigns) Review(ctx context.Context, p identity.Principal, id, hash string, yes bool, reason string) (CampaignBatch, error) {
	var out CampaignBatch
	if !c.authorized(p, "marketing:approve") || !c.authorized(p, "notification:approve") || !validDigest(hash) || !cr.ValidText(reason, 2048) {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	plan, h, e := c.load(ctx, s.base.pool, id)
	if e != nil {
		return out, e
	}
	if h != hash || plan.Requester == p.Subject {
		return out, ErrApproval
	}
	out = CampaignBatch{CampaignID: id, CampaignSHA256: h, Complete: true}
	for _, m := range plan.Members {
		for step := 1; step <= len(plan.Request.Steps); step++ {
			request := campaignScheduleRequest(plan, m, step)
			item := CampaignItem{LeadID: m.LeadID, Step: step, DeliveryKey: scheduleKey(s.tenant, request.RequestID)}
			bound, bh, err := s.request(ctx, item.DeliveryKey)
			if err != nil || bound.SourceSHA256 != h || bound.Request.CampaignID != id || bound.LeadID != m.LeadID || bound.Request.Step != step {
				item.Code = "REQUEST_UNAVAILABLE"
				out.Complete = false
			} else {
				state, err := s.Decide(ctx, p, item.DeliveryKey, bh, yes, reason)
				if errors.Is(err, approval.ErrNotPending) {
					state, err = c.decisionReplay(ctx, p, item.DeliveryKey, bh, yes, reason)
				}
				item.RequestSHA256 = bh
				item.State = string(state)
				if err != nil {
					item.Code = "DECISION_UNCONFIRMED"
					out.Complete = false
				}
			}
			out.Items = append(out.Items, item)
		}
	}
	return out, nil
}
func (c *Campaigns) Stop(ctx context.Context, p identity.Principal, id, hash, reason string) error {
	if !c.authorized(p, "marketing:cancel") || !validDigest(hash) || !cr.ValidText(reason, 2048) {
		return ErrApproval
	}
	s := c.schedule.approvals
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	var actual string
	if e = tx.QueryRow(ctx, `select campaign_sha256 from communication.whatsapp_campaign where tenant_id=$1 and campaign_id=$2 and organization_id=$3 for update`, s.tenant, id, s.org).Scan(&actual); e != nil {
		return e
	}
	if actual != hash {
		return ErrApproval
	}
	_, e = tx.Exec(ctx, `insert into communication.whatsapp_campaign_stop(tenant_id,campaign_id,campaign_sha256,actor_subject,reason)values($1,$2,$3,$4,$5)on conflict do nothing`, s.tenant, id, hash, p.Subject, reason)
	if e != nil {
		return e
	}
	var actor, storedReason string
	e = tx.QueryRow(ctx, "select actor_subject,reason from communication.whatsapp_campaign_stop where tenant_id=$1 and campaign_id=$2", s.tenant, id).Scan(&actor, &storedReason)
	if e != nil {
		return e
	}
	if actor != p.Subject || storedReason != reason {
		return ErrApproval
	}
	return tx.Commit(ctx)
}

// A replay reads the original immutable decision; it creates no approval/job.
func (c *Campaigns) decisionReplay(ctx context.Context, p identity.Principal, key, hash string, yes bool, reason string) (approval.State, error) {
	s := c.schedule.approvals
	want := approval.StateRejected
	if yes {
		want = approval.StateApproved
	}
	var exact bool
	e := s.base.pool.QueryRow(ctx, `select exists(
      select 1 from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id
      where r.tenant_id=$1 and r.organization_id=$2 and r.request_id=$3 and r.evidence_sha=$4 and r.kind='whatsapp_schedule' and r.state=$5
      and d.reviewer=$6 and d.reviewer<>r.requester and d.approved=$7 and d.reason=$8
      and(select count(*)from approval.decision d2 where d2.tenant_id=r.tenant_id and d2.request_id=r.request_id)=1)`, s.tenant, s.org, key, hash, want, p.Subject, yes, reason).Scan(&exact)
	if e != nil {
		return "", e
	}
	if !exact {
		return "", ErrApproval
	}
	return want, nil
}
````

### FILE: `internal/whatsappbridge/campaign_partial_test.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ba3c5e501e5edb8e0ff4380d0c5b19520279455b964517caf5da265e93cebd0a"
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
````

### FILE: `internal/whatsappbridge/campaign_schedule_source.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "95d5f414020ec31220b8b430950eaed4341574f6fda023d1fe14fc69a1c249c9"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED campaign source implements the existing schedule extension.
import (
	"context"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type campaignQuery interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func (c *Campaigns) load(ctx context.Context, db campaignQuery, id string) (CampaignSnapshot, string, error) {
	var out CampaignSnapshot
	var raw []byte
	var h string
	s := c.schedule.approvals
	e := db.QueryRow(ctx, `select payload,campaign_sha256 from communication.whatsapp_campaign where tenant_id=$1 and organization_id=$2 and campaign_id=$3 for share`, s.tenant, s.org, id).Scan(&raw, &h)
	if errors.Is(e, pgx.ErrNoRows) {
		return out, "", ErrApproval
	}
	if e != nil {
		return out, "", e
	}
	if json.Unmarshal(raw, &out) != nil {
		return out, "", ErrApproval
	}
	_, actual, e := campaignHash(out)
	if e != nil || actual != h || out.Schema != "elite-whatsapp-campaign/v1" || out.TenantID != s.tenant || out.OrganizationID != s.org || out.ConnectionID != s.connection || out.ProfileSHA256 != digest(s.profile) || out.Policy != c.policy || out.Request.ID != id || len(out.Members) != len(out.Request.Members) || c.validate(out.Request) != nil {
		return out, "", ErrApproval
	}
	return out, h, nil
}
func campaignScheduleRequest(plan CampaignSnapshot, member CampaignMember, step int) ScheduleRequest {
	r := plan.Request.Steps[step-1]
	return ScheduleRequest{RequestID: "campaign-" + digest([]byte(plan.Request.ID+"\x00"+member.LeadID+"\x00"+fmt.Sprint(step))), CampaignID: plan.Request.ID, LeadID: member.LeadID, Step: step, Recipient: member.Recipient, TemplateName: r.TemplateName, LanguageCode: r.LanguageCode, BodyParameters: append([]string{}, r.BodyParameters...), NotBefore: r.NotBefore, ExpiresAt: r.ExpiresAt}
}
func (c *Campaigns) stoppedOrConverted(ctx context.Context, db campaignQuery, id, lead string) (bool, error) {
	var stopped bool
	s := c.schedule.approvals
	e := db.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_campaign_stop where tenant_id=$1 and campaign_id=$2)
 or exists(select 1 from sales.quotation q join sales.quotation_acceptance a on a.tenant_id=q.tenant_id and a.quotation_id=q.quotation_id
 join sales.customer_order o on o.tenant_id=a.tenant_id and o.order_id=a.order_id
 where q.tenant_id=$1 and q.organization_id=$3 and q.lead_id=$4 and q.order_id=a.order_id and a.accepted_at<=clock_timestamp())`, s.tenant, id, s.org, lead).Scan(&stopped)
	return stopped, e
}
func (c *Campaigns) Source(ctx context.Context, tx pgx.Tx, r ScheduleRequest, preparing bool) (ScheduleContext, error) {
	var out ScheduleContext
	plan, h, e := c.load(ctx, tx, r.CampaignID)
	if e != nil {
		return out, e
	}
	if r.Step < 1 || r.Step > len(plan.Request.Steps) || r.AppointmentID != "" {
		return out, ErrApproval
	}
	var selected *CampaignMember
	for i := range plan.Members {
		if plan.Members[i].LeadID == r.LeadID {
			selected = &plan.Members[i]
			break
		}
	}
	if selected == nil {
		return out, ErrApproval
	}
	_, requestHash, e := campaignHash(r)
	_, expected, e2 := campaignHash(campaignScheduleRequest(plan, *selected, r.Step))
	if e != nil || e2 != nil || requestHash != expected {
		return out, ErrApproval
	}
	fresh, e := c.member(ctx, tx, CampaignRecipient{selected.LeadID, selected.Recipient})
	if e != nil {
		return out, e
	}
	_, memberHash, e := campaignHash(fresh)
	_, before, e2 := campaignHash(*selected)
	if e != nil || e2 != nil || memberHash != before || !contains(plan.Request.Sources, fresh.Source) || !contains(plan.Request.States, fresh.Lifecycle) {
		return out, ErrApproval
	}
	if stopped, e := c.stoppedOrConverted(ctx, tx, r.CampaignID, r.LeadID); e != nil {
		return out, e
	} else if stopped {
		return out, ErrApproval
	}
	var now time.Time
	if e = tx.QueryRow(ctx, "select clock_timestamp()").Scan(&now); e != nil {
		return out, e
	}
	if !r.ExpiresAt.After(now) || preparing && !r.NotBefore.After(now) {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	body, _ := json.Marshal(struct {
		Recipient      string   `json:"recipient"`
		TemplateName   string   `json:"template_name"`
		LanguageCode   string   `json:"language_code"`
		BodyParameters []string `json:"body_parameters"`
	}{r.Recipient, r.TemplateName, r.LanguageCode, r.BodyParameters})
	out = ScheduleContext{Schema: "elite-whatsapp-schedule/v1", Request: r, DeliveryKey: scheduleKey(s.tenant, r.RequestID), OrganizationID: s.org, ConnectionID: s.connection, LeadID: r.LeadID, SubjectID: fresh.SubjectID, ExternalIDHMAC: fresh.ExternalIDHMAC, BindingVersion: fresh.BindingVersion, PolicyVersion: fresh.BindingPolicy, ConsentID: fresh.ConsentID, ConsentPurpose: c.policy.ConsentPurpose, ConsentEvidenceSHA256: fresh.ConsentSHA256, ProfileSHA256: digest(s.profile), NotBefore: r.NotBefore, ExpiresAt: r.ExpiresAt, SourceSHA256: h, MemberSHA256: memberHash}
	out.Message = channels.Message{TenantID: s.tenant, ChannelCode: "whatsapp", Direction: channels.DirectionOut, DeliveryKey: out.DeliveryKey, ExternalID: r.Recipient, Text: string(body)}
	out.MessageSHA256, e = outbounddelivery.MessageSHA256(out.Message)
	return out, e
}
func (c *Campaigns) Enqueue(ctx context.Context, tx pgx.Tx, b ScheduleContext, h string) error {
	s := c.schedule.approvals
	_, e := tx.Exec(ctx, `insert into communication.whatsapp_schedule(tenant_id,delivery_key,job_id,organization_id,appointment_id,appointment_version,not_before,expires_at,request_sha256,campaign_id,lead_id,campaign_step)
 select m.tenant_id,$2,gen_random_uuid(),p.organization_id,null,0,$3,$4,$5,m.campaign_id,m.lead_id,$6
 from communication.whatsapp_campaign_member m join communication.whatsapp_campaign p using(tenant_id,campaign_id)
 where m.tenant_id=$1 and m.campaign_id=$7 and m.lead_id=$8 and m.member_sha256=$9 and p.campaign_sha256=$10 and p.organization_id=$11
 on conflict(tenant_id,delivery_key)do nothing`, s.tenant, b.DeliveryKey, b.NotBefore, b.ExpiresAt, h, b.Request.Step, b.Request.CampaignID, b.LeadID, b.MemberSHA256, b.SourceSHA256, s.org)
	if e != nil {
		return e
	}
	var same bool
	e = tx.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_schedule where tenant_id=$1 and delivery_key=$2 and request_sha256=$3 and campaign_id=$4 and lead_id=$5 and campaign_step=$6)`, s.tenant, b.DeliveryKey, h, b.Request.CampaignID, b.LeadID, b.Request.Step).Scan(&same)
	if e != nil {
		return e
	}
	if !same {
		return ErrApproval
	}
	return nil
}
func (c *Campaigns) Admit(ctx context.Context, tx pgx.Tx, b ScheduleContext) error {
	if b.Request.Step <= 1 {
		return nil
	}
	plan, _, e := c.load(ctx, tx, b.Request.CampaignID)
	if e != nil {
		return e
	}
	for _, m := range plan.Members {
		if m.LeadID == b.LeadID {
			previous := campaignScheduleRequest(plan, m, b.Request.Step-1)
			var accepted bool
			e = tx.QueryRow(ctx, `select exists(select 1 from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2 and channel_code='whatsapp' and state='accepted')`, b.Message.TenantID, scheduleKey(b.Message.TenantID, previous.RequestID)).Scan(&accepted)
			if e != nil {
				return e
			}
			if !accepted {
				return ErrApproval
			}
			return nil
		}
	}
	return ErrApproval
}
````

### FILE: `internal/whatsappbridge/campaign_status.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-CAMPAIGNS:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "40f065fcfa00fb2f213cc91f94b2665d8d5c156371b723814feb71699960216c"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED read-only observation of original message and quotation/order evidence.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type CampaignConversion struct {
	LeadID         string    `json:"lead_id"`
	QuotationID    string    `json:"quotation_id"`
	OrderID        string    `json:"order_id"`
	OrderState     string    `json:"order_state"`
	AcceptedAt     time.Time `json:"accepted_at"`
	EvidenceSHA256 string    `json:"evidence_sha256"`
}
type CampaignItemStatus struct {
	LeadID      string          `json:"lead_id"`
	Step        int             `json:"step"`
	DeliveryKey string          `json:"delivery_key"`
	Preparation string          `json:"preparation"`
	Schedule    *ScheduleStatus `json:"schedule,omitempty"`
}
type CampaignStatus struct {
	CampaignSHA256    string               `json:"campaign_sha256"`
	Snapshot          CampaignSnapshot     `json:"snapshot"`
	Stopped           bool                 `json:"stopped"`
	Items             []CampaignItemStatus `json:"items"`
	Conversions       []CampaignConversion `json:"conversions"`
	ConversionMeaning string               `json:"conversion_meaning"`
}

func (c *Campaigns) Status(ctx context.Context, p identity.Principal, id string) (CampaignStatus, error) {
	var out CampaignStatus
	if !c.authorized(p, "marketing:read") || !c.authorized(p, "notification:read") {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	plan, h, e := c.load(ctx, s.base.pool, id)
	if e != nil {
		return out, e
	}
	out = CampaignStatus{CampaignSHA256: h, Snapshot: plan, ConversionMeaning: "First original quote/order acceptance after an accepted campaign message; observed sequence, not causal attribution or payment proof", Conversions: []CampaignConversion{}}
	if e = s.base.pool.QueryRow(ctx, "select exists(select 1 from communication.whatsapp_campaign_stop where tenant_id=$1 and campaign_id=$2)", s.tenant, id).Scan(&out.Stopped); e != nil {
		return out, e
	}
	for _, m := range plan.Members {
		for step := 1; step <= len(plan.Request.Steps); step++ {
			request := campaignScheduleRequest(plan, m, step)
			key := scheduleKey(s.tenant, request.RequestID)
			item := CampaignItemStatus{LeadID: m.LeadID, Step: step, DeliveryKey: key, Preparation: "NOT_PREPARED"}
			var exists bool
			e = s.base.pool.QueryRow(ctx, "select exists(select 1 from approval.request where tenant_id=$1 and request_id=$2)", s.tenant, key).Scan(&exists)
			if e != nil {
				return out, e
			}
			if exists {
				status, e := s.Status(ctx, p, key)
				if e != nil {
					return out, e
				}
				if status.Context.SourceSHA256 != h {
					return out, ErrApproval
				}
				item.Preparation = "PREPARED"
				item.Schedule = &status
			}
			out.Items = append(out.Items, item)
		}
		conversion := CampaignConversion{LeadID: m.LeadID}
		e = s.base.pool.QueryRow(ctx, `select a.quotation_id,a.order_id,o.state,a.accepted_at,a.evidence_sha256_hex
 from sales.quotation q join sales.quotation_acceptance a on a.tenant_id=q.tenant_id and a.quotation_id=q.quotation_id and a.order_id=q.order_id
 join sales.customer_order o on o.tenant_id=a.tenant_id and o.order_id=a.order_id and o.organization_id=q.organization_id and o.customer_principal_id=a.customer_principal_id
 where q.tenant_id=$1 and q.organization_id=$2 and q.lead_id=$3 and a.accepted_at<=clock_timestamp()
 and exists(select 1 from communication.whatsapp_schedule s join communication.outbound_delivery d on d.tenant_id=s.tenant_id and d.delivery_key=s.delivery_key and d.channel_code='whatsapp'
 where s.tenant_id=q.tenant_id and s.campaign_id=$4 and s.lead_id=q.lead_id and d.state='accepted' and d.accepted_at<=a.accepted_at)
 order by a.accepted_at,a.quotation_id limit 1`, s.tenant, s.org, m.LeadID, id).Scan(&conversion.QuotationID, &conversion.OrderID, &conversion.OrderState, &conversion.AcceptedAt, &conversion.EvidenceSHA256)
		if e != nil && !errors.Is(e, pgx.ErrNoRows) {
			return out, e
		}
		if e == nil {
			out.Conversions = append(out.Conversions, conversion)
		}
	}
	return out, nil
}
func (c *Campaigns) VerifyInfrastructure(ctx context.Context) error {
	var valid bool
	e := c.schedule.approvals.base.pool.QueryRow(ctx, `select
 (select count(*)from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('communication.whatsapp_campaign')and tgname='whatsapp_campaign_immutable'
 or tgrelid=to_regclass('communication.whatsapp_campaign_member')and tgname='whatsapp_campaign_member_immutable'
 or tgrelid=to_regclass('communication.whatsapp_campaign_stop')and tgname='whatsapp_campaign_stop_immutable'))=3
 and(select count(*)from pg_constraint where conrelid=to_regclass('communication.whatsapp_schedule') and convalidated
 and conname in('whatsapp_schedule_source_check','whatsapp_schedule_campaign_member_fk','whatsapp_schedule_campaign_step_unique'))=3`).Scan(&valid)
	if e != nil || !valid {
		return ErrApproval
	}
	return nil
}
````

## 6. Configuration surface

docs/CAMPAIGN_CONNECTED_REFERENCE.md; exact optional campaign_policy_file/campaign_policy_sha256 in the existing WhatsApp host. Explicit policy purpose separate from appointment consent. Original account/identity/credentials reused.

## 7. Dependency bill

No dependency or upstream upgrade. Original Meta source/payload builder, human approval, jobs/fence, CRM/contact/consent and quote/order sources preserved. All new source/policy/SQL/API/host/test binding AUTHORED.

## 8. Apply order

Apply0084 after scheduled0083. Same host/worker with optional campaign source and typed job scope. Missing pack/policy/immutable guard refused; disabled extension leaves campaign jobs untouched.

## 9. Verification

6campaigns/10jobs/4actual adapterPOST/10results/1original quote-order;12concurrentfirst sends;2steps,unknown recovery/prior suppression,stop,optout/source/converted suppression. Partialcreate/review/resume2unique requests/jobs/decisions and0POST;4read-only replay negatives. Host4modules/2loops/mTLS/shutdown,rollback/appointment compatibility,boundaries/vet/build and2sfuzz.

## 10. Reconstruction evidence

reconstruction_evidence/CAMPAIGN_CONNECTED_RELEASE_V402.md/json. Closes selected T2805 alongside unchanged linked ML/Merchant/WA/Page/scheduled proofs. T2803 and later ordered controls remain separate.

