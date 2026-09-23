package whatsappbridge

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNotificationHistoryReadOnlyAndScoped(t *testing.T) {
	r, o, pool, p, c, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	value, err := o.Approvals.ReadNotificationHistory(ctx, p, c.OrganizationID, c.AppointmentID)
	if err != nil || len(value.Items) != 1 || value.Items[0].ConfirmationEventID != c.ConfirmationEventID || value.Items[0].Status.DeliveryStatus != "not_observed_by_this_reader" {
		t.Fatal(err, value)
	}
	w, err := NewStatusWorker(r, "history-worker", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if result, err := w.ProcessOnce(ctx, p); err != nil || !result.Completed {
		t.Fatal(err, result)
	}
	value, err = o.Approvals.ReadNotificationHistory(ctx, p, c.OrganizationID, c.AppointmentID)
	if err != nil || len(value.Items) != 1 || value.Items[0].Status.DeliveryStatus != "observed_delivered" {
		t.Fatal(err, value)
	}
	for _, mode := range []string{"permission", "organization", "tenant", "appointment"} {
		t.Run(mode, func(t *testing.T) {
			bad := p
			appointment := c.AppointmentID
			switch mode {
			case "permission":
				bad.Permissions = map[string]struct{}{}
			case "organization":
				bad.Organizations = map[string]struct{}{}
			case "tenant":
				bad.TenantID = "00000000-0000-4000-8000-000000000000"
			case "appointment":
				appointment = "other"
			}
			if _, err := o.Approvals.ReadNotificationHistory(ctx, bad, c.OrganizationID, appointment); !errors.Is(err, ErrNotificationNotFound) {
				t.Fatal("foreign history", err)
			}
		})
	}
	other, command, empty := approvalFixtureData(t, pool)
	value, err = empty.ReadNotificationHistory(ctx, other, command.OrganizationID, command.AppointmentID)
	if err != nil || len(value.Items) != 0 {
		t.Fatal("empty history", err, value)
	}
	if _, err := pool.Exec(ctx, `update communication.outbound_delivery set recipient_hmac=repeat('f',64) where tenant_id=$1`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	if _, err := o.Approvals.ReadNotificationHistory(ctx, p, c.OrganizationID, c.AppointmentID); !errors.Is(err, ErrNotificationIntegrity) {
		t.Fatal("tampered history", err)
	}
}

func TestNotificationHistoryHTTPBoundary(t *testing.T) {
	_, o, _, p, c, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	v, sign := notificationIssuer(t)
	module := &AppointmentNotificationModule{approvals: o.Approvals, sender: &Sender{TenantID: p.TenantID}}
	mux := http.NewServeMux()
	module.Register(mux, v)
	api := httptest.NewServer(mux)
	defer api.Close()
	for _, test := range []struct {
		name, query, token string
		status             int
	}{{"valid", "organization_id=store-1", sign(p), 200}, {"duplicate", "organization_id=store-1&organization_id=store-1", sign(p), 400}, {"extra", "organization_id=store-1&tenant=other", sign(p), 400}, {"foreign", "organization_id=other", sign(p), 403}, {"no-token", "organization_id=store-1", "", 401}, {"no-permission", "organization_id=store-1", sign(identity.Principal{TenantID: p.TenantID, Subject: p.Subject}), 403}} {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", api.URL+"/v1/franchise/appointments/"+c.AppointmentID+"/whatsapp-confirmations?"+test.query, nil)
			if test.token != "" {
				req.Header.Set("Authorization", "Bearer "+test.token)
			}
			response, err := api.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != test.status || response.Header.Get("Cache-Control") != "no-store" {
				t.Fatal(response.StatusCode)
			}
			var body map[string]any
			if json.NewDecoder(response.Body).Decode(&body) != nil {
				t.Fatal("response JSON")
			}
			raw, _ := json.Marshal(body)
			for _, secret := range []string{"5491112345678", "wamid.synthetic", "test-app-secret"} {
				if strings.Contains(string(raw), secret) {
					t.Fatal("sensitive history")
				}
			}
		})
	}
}

func TestNotificationHistoryOverflowIsExplicit(t *testing.T) {
	_, o, pool, p, c, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	// Synthetic historical approvals test the read limit, not the send admission.
	for i := 0; i < 20; i++ {
		var id string
		if err := pool.QueryRow(ctx, `insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject) values($1,gen_random_uuid(),$2,'requested','confirmed','fixture-history') returning transition_id`, p.TenantID, c.AppointmentID).Scan(&id); err != nil {
			t.Fatal(err)
		}
		if _, err := pool.Exec(ctx, `insert into communication.whatsapp_appointment_approval(tenant_id,delivery_key,channel_code,appointment_id,appointment_version,organization_id,confirmation_event_id,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_id,consent_purpose,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,approved_at,expires_at) select tenant_id,'fixture-history-'||$2::text,channel_code,appointment_id,appointment_version,organization_id,$2::uuid,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_id,consent_purpose,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,approved_at,expires_at from communication.whatsapp_appointment_approval where tenant_id=$1 and delivery_key=$3`, p.TenantID, id, c.Message.DeliveryKey); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := o.Approvals.ReadNotificationHistory(ctx, p, c.OrganizationID, c.AppointmentID); !errors.Is(err, ErrNotificationHistoryLimit) {
		t.Fatal("silently truncated history", err)
	}
}
