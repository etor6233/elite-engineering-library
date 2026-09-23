package whatsappbridge

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func approvalPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	raw := os.Getenv("ELITE_WHATSAPP_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("requires dedicated local elite_whatsapp_ database and 50 migrations")
	}
	// Same local/dedicated guard as the existing process integration test.
	if !strings.HasPrefix(raw, "postgres://postgres@127.0.0.1:55959/elite_whatsapp_") {
		t.Fatal("unsafe test database target")
	}
	pool, err := pgxpool.New(context.Background(), raw)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func approvalFixtureData(t *testing.T, pool *pgxpool.Pool) (identity.Principal, AppointmentApprovalCommand, *PostgresAppointmentApprovals) {
	t.Helper()
	ctx := context.Background()
	tenant := leadstream.StableUUID(t.Name(), time.Now().Format(time.RFC3339Nano))
	event := leadstream.StableUUID(tenant, "confirmation")
	statements := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1::uuid,'wa-'||$1::text,'Synthetic approval','Synthetic approval')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-1','store-1','Store 1','store')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'lead-1','store-1','new','fixture','{}',true)`,
		`insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,appointment_kind,starts_at,state,version) values($1,'appointment-1','store-1','lead-1','consultation',clock_timestamp()+interval '1 day','confirmed',3)`,
		`insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex) values($1,'consent-1','lead-1','fixture-appointment-notification','policy-1','granted',clock_timestamp()-interval '1 minute',repeat('c',64))`,
	}
	for _, q := range statements {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, `insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject) values($1,$2,'appointment-1','requested','confirmed','fixture-operator')`, tenant, event); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'appointment','appointment-1',3,'appointment.confirmed',1,clock_timestamp(),'{}')`, tenant, event); err != nil {
		t.Fatal(err)
	}
	key := []byte("0123456789abcdef0123456789abcdef")
	binding, err := postgres.NewContactIdentityStore(pool, key)
	if err != nil {
		t.Fatal(err)
	}
	m := messageFixture()
	m.TenantID = tenant
	m.DeliveryKey = AppointmentConfirmationDeliveryKey(tenant, event)
	_, err = binding.Apply(ctx, contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: m.ExternalID, LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(-time.Minute), RequestID: "binding-1"})
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := NewPostgresAppointmentApprovals(pool, key, "fixture-appointment-notification", "policy-1")
	if err != nil {
		t.Fatal(err)
	}
	p := identity.Principal{Subject: "fixture-operator", TenantID: tenant, Permissions: map[string]struct{}{"appointment:manage": {}}, Organizations: map[string]struct{}{"store-1": {}}}
	c := AppointmentApprovalCommand{Message: m, OrganizationID: "store-1", AppointmentID: "appointment-1", ConfirmationEventID: event, AppointmentVersion: 3, BindingVersion: 1, PolicyVersion: "policy-1", ConsentID: "consent-1", ConsentPurpose: "fixture-appointment-notification", ConsentEvidenceSHA256: strings.Repeat("c", 64), ProfileSHA256: digest([]byte(`{"synthetic_fixture":true}`)), EvidenceSHA256: strings.Repeat("d", 64), ExpiresAt: time.Now().Add(time.Hour)}
	return p, c, resolver
}

func TestAppointmentApprovalDurableConcurrentAndFenced(t *testing.T) {
	pool := approvalPool(t)
	p, c, resolver := approvalFixtureData(t, pool)
	ctx := context.Background()
	var wg sync.WaitGroup
	errs := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := resolver.Approve(ctx, p, c); errs <- err }()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	var count int
	if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_appointment_approval where tenant_id=$1`, p.TenantID).Scan(&count); err != nil || count != 1 {
		t.Fatal(count, err)
	}
	if _, err := resolver.ResolveWhatsAppApproval(ctx, leadstream.StableUUID("other-tenant"), c.Message.DeliveryKey); !errors.Is(err, ErrApproval) {
		t.Fatal("foreign tenant approval", err)
	}
	changed := c
	changed.Message.Text += " "
	if _, err := resolver.Approve(ctx, p, changed); !errors.Is(err, ErrApproval) {
		t.Fatal("divergent replay accepted", err)
	}
	for _, mutation := range []string{`update communication.whatsapp_appointment_approval set approved_by='changed' where tenant_id=$1`, `delete from communication.whatsapp_appointment_approval where tenant_id=$1`} {
		if _, err := pool.Exec(ctx, mutation, p.TenantID); err == nil {
			t.Fatal("immutable evidence mutated")
		}
	}
	sender, _, marker := senderFixture(t, "accepted")
	sender.TenantID = p.TenantID
	sender.Approvals = resolver
	for i := 0; i < 2; i++ {
		store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Second)
		if err != nil {
			t.Fatal(err)
		}
		channel := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiverFixture{}, Sender: sender, Store: store}
		if err = channel.Send(ctx, c.Message); err != nil {
			t.Fatal(err)
		}
	}
	calls, err := os.ReadFile(marker)
	if err != nil || string(calls) != "call\n" {
		t.Fatal("duplicate child send", err)
	}
	var state string
	if err = pool.QueryRow(ctx, `select state from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2`, p.TenantID, c.Message.DeliveryKey).Scan(&state); err != nil || state != "accepted" {
		t.Fatal(state, err)
	}
}

func TestAppointmentApprovalRequiresExactProof(t *testing.T) {
	pool := approvalPool(t)
	for _, mode := range []string{"role", "organization", "tenant", "recipient", "event", "appointment_version", "binding_version", "policy", "consent_id", "consent_hash", "consent_purpose", "profile_hash", "expired", "key", "no_outbox", "requested", "withdrawn", "tied_consent", "future_binding"} {
		t.Run(mode, func(t *testing.T) {
			p, c, resolver := approvalFixtureData(t, pool)
			ctx := context.Background()
			switch mode {
			case "role":
				p.Permissions = map[string]struct{}{}
			case "organization":
				p.Organizations = map[string]struct{}{}
			case "tenant":
				p.TenantID = leadstream.StableUUID("different")
			case "recipient":
				c.Message.ExternalID = "5491199999999"
			case "event":
				c.ConfirmationEventID = leadstream.StableUUID("other-event")
				c.Message.DeliveryKey = AppointmentConfirmationDeliveryKey(p.TenantID, c.ConfirmationEventID)
			case "appointment_version":
				c.AppointmentVersion = 2
			case "binding_version":
				c.BindingVersion = 2
			case "policy":
				c.PolicyVersion = "other-policy"
			case "consent_id":
				c.ConsentID = "other-consent"
			case "consent_hash":
				c.ConsentEvidenceSHA256 = strings.Repeat("a", 64)
			case "consent_purpose":
				c.ConsentPurpose = "other-purpose"
			case "profile_hash":
				c.ProfileSHA256 = ""
			case "expired":
				c.ExpiresAt = time.Now().Add(-time.Second)
			case "key":
				c.Message.DeliveryKey = strings.Repeat("k", 64)
			case "no_outbox":
				if _, err := pool.Exec(ctx, `delete from platform.outbox_event where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "requested":
				if _, err := pool.Exec(ctx, `update crm.appointment set state='requested' where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "withdrawn", "tied_consent":
				timing := "clock_timestamp()"
				if mode == "tied_consent" {
					timing = "occurred_at"
				}
				q := `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex) select tenant_id,'newer',lead_id,purpose_code,policy_version,'withdrawn',` + timing + `,evidence_sha256_hex from crm.consent_evidence where tenant_id=$1`
				if _, err := pool.Exec(ctx, q, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "future_binding":
				if _, err := pool.Exec(ctx, `update communication.contact_channel_binding set effective_at=clock_timestamp()+interval '1 day' where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := resolver.Approve(ctx, p, c); !errors.Is(err, ErrApproval) {
				t.Fatalf("missing proof accepted: %v", err)
			}
			var count int
			if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_appointment_approval where tenant_id=$1`, c.Message.TenantID).Scan(&count); err != nil || count != 0 {
				t.Fatal("unexpected durable approval", count, err)
			}
		})
	}
}

func TestAppointmentApprovalRechecksRevocationAndChanges(t *testing.T) {
	pool := approvalPool(t)
	for _, mode := range []string{"revoked", "pii", "version", "cancelled", "consent_withdrawn", "expiry"} {
		t.Run(mode, func(t *testing.T) {
			p, c, resolver := approvalFixtureData(t, pool)
			ctx := context.Background()
			if mode == "expiry" {
				c.ExpiresAt = time.Now().Add(300 * time.Millisecond)
			}
			if _, err := resolver.Approve(ctx, p, c); err != nil {
				t.Fatal(err)
			}
			queries := map[string]string{
				"revoked":           `update communication.contact_channel_binding set state='revoked',pii_allowed=false,version=version+1 where tenant_id=$1`,
				"pii":               `update communication.contact_channel_binding set pii_allowed=false,version=version+1 where tenant_id=$1`,
				"version":           `update crm.appointment set version=version+1 where tenant_id=$1`,
				"cancelled":         `update crm.appointment set state='cancelled',version=version+1 where tenant_id=$1`,
				"consent_withdrawn": `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex) select tenant_id,'withdrawal',lead_id,purpose_code,policy_version,'withdrawn',clock_timestamp(),evidence_sha256_hex from crm.consent_evidence where tenant_id=$1`,
			}
			if mode == "expiry" {
				time.Sleep(350 * time.Millisecond)
			} else {
				if _, err := pool.Exec(ctx, queries[mode], p.TenantID); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := resolver.ResolveWhatsAppApproval(ctx, p.TenantID, c.Message.DeliveryKey); !errors.Is(err, ErrApproval) {
				t.Fatal("stale approval returned", err)
			}
			sender, _, marker := senderFixture(t, "accepted")
			sender.TenantID = p.TenantID
			sender.Approvals = resolver
			if _, err := sender.SendWithReceipt(ctx, c.Message); !errors.Is(err, ErrBridge) {
				t.Fatal("stale approval sent", err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("child invoked after revocation")
			}
		})
	}
}
