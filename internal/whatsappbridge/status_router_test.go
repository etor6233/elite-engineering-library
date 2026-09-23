package whatsappbridge

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5/pgxpool"
)

func alterStatusBatch(t *testing.T, b SignedStatusWebhook, edit func(map[string]any)) SignedStatusWebhook {
	t.Helper()
	var body map[string]any
	if json.Unmarshal(b.Body, &body) != nil {
		t.Fatal("fixture JSON")
	}
	edit(body)
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	mac := hmac.New(sha256.New, []byte("test-app-secret"))
	mac.Write(raw)
	return SignedStatusWebhook{Body: raw, Signature: "sha256=" + hex.EncodeToString(mac.Sum(nil))}
}
func statusValue(body map[string]any) map[string]any {
	return body["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
}
func routerFixture(t *testing.T, batch SignedStatusWebhook) (*StatusRouter, *StatusObserver, *pgxpool.Pool, identity.Principal, AppointmentApprovalCommand, string, string) {
	t.Helper()
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	if _, err := pool.Exec(context.Background(), `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'wa-primary','meta-whatsapp','META_APP_SECRET')`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	config := receiverConfig(o)
	h, err := NewWebhookReceiver(config)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(h)
	defer server.Close()
	response, _ := webhookRequest(t, server, batch)
	if response.StatusCode != 200 {
		t.Fatalf("receiver %d", response.StatusCode)
	}
	directory := filepath.Join(o.Process.EvidenceDirectory, digest([]byte(p.TenantID+"\x00"+c.Message.DeliveryKey)))
	if err := os.MkdirAll(directory, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, "SEND_RECEIPT.json")
	if err := os.WriteFile(path, receipt, 0600); err != nil {
		t.Fatal(err)
	}
	r, err := NewStatusRouter(o, config.ConnectionID, config.RetentionApprovalSHA256, []byte("0123456789abcdef0123456789abcdef"), 4)
	if err != nil {
		t.Fatal(err)
	}
	return r, o, pool, p, c, "wa:" + digest(batch.Body), path
}

func TestStatusRouterRetainedToObservedAndReplay(t *testing.T) {
	batch := alterStatusBatch(t, statusBatch("delivered", "1603086314", nil), func(body map[string]any) {
		v := statusValue(body)
		v["statuses"] = append(v["statuses"].([]any), map[string]any{"id": "wamid.synthetic", "recipient_id": "5491112345678", "status": "sent", "timestamp": "1603086313"})
		// Case variants must not override the exact fields verified by Python.
		v["Statuses"] = []any{map[string]any{"id": "wamid.foreign", "recipient_id": "000", "status": "read", "timestamp": "1603086315"}}
	})
	r, o, pool, p, c, event, _ := routerFixture(t, batch)
	for _, expected := range []int{2, 0, 0} {
		n, err := r.ObserveRetained(context.Background(), p, event)
		if err != nil || n != expected {
			t.Fatalf("route inserted=%d expected=%d err=%v", n, expected, err)
		}
	}
	value, err := o.Approvals.ReadNotificationStatus(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID)
	if err != nil || value.DeliveryStatus != "observed_delivered" || value.ProviderEventCount != 2 {
		t.Fatalf("read: %+v %v", value, err)
	}
	var untouched bool
	if err := pool.QueryRow(context.Background(), `select
 (select state='received' and processed_at is null from integration.webhook_event where tenant_id=$1 and provider_event_id=$2)
 and (select count(*)=1 and bool_and(completed_at is null) from platform.job where tenant_id=$1 and queue='provider-events')
 and (select attempt_count=1 and state='accepted' from communication.outbound_delivery where tenant_id=$1 and delivery_key=$3)`, p.TenantID, event, c.Message.DeliveryKey).Scan(&untouched); err != nil || !untouched {
		t.Fatalf("router changed receipt/job/send: %v %v", untouched, err)
	}
}

func TestStatusRouterRejectsBeforeObservation(t *testing.T) {
	for _, mode := range []string{"unknown", "recipient", "mixed", "mixed-unknown", "permission", "organization", "tenant", "disabled", "connection-org", "retention", "profile", "signature", "payload", "count", "missing-receipt", "changed-receipt", "wrong-key", "busy", "scope-changed", "ambiguous"} {
		t.Run(mode, func(t *testing.T) {
			batch := statusBatch("delivered", "1603086314", nil)
			if mode == "unknown" {
				batch = statusBatch("delivered", "1603086314", map[string]any{"id": "wamid.unknown"})
			}
			if mode == "recipient" {
				batch = statusBatch("delivered", "1603086314", map[string]any{"recipient_id": "9999999"})
			}
			if mode == "mixed" || mode == "mixed-unknown" {
				batch = alterStatusBatch(t, batch, func(body map[string]any) {
					v := statusValue(body)
					if mode == "mixed" {
						v["messages"] = []any{map[string]any{"id": "wamid.inbound", "from": "5491112345678", "type": "text", "timestamp": "1603086315", "text": map[string]any{"body": "Synthetic"}}}
					} else {
						v["statuses"] = append(v["statuses"].([]any), map[string]any{"id": "wamid.unknown", "recipient_id": "5491112345678", "status": "sent", "timestamp": "1603086313"})
					}
				})
			}
			r, o, pool, p, c, event, path := routerFixture(t, batch)
			originalTenant := p.TenantID
			execSQL := func(sql string, args ...any) {
				t.Helper()
				if _, err := pool.Exec(context.Background(), sql, args...); err != nil {
					t.Fatal(err)
				}
			}
			switch mode {
			case "permission":
				p.Permissions = map[string]struct{}{}
			case "organization":
				p.Organizations = map[string]struct{}{"other": {}}
			case "tenant":
				p.TenantID = "00000000-0000-0000-0000-000000000001"
			case "disabled":
				execSQL(`update integration.provider_connection set state='disabled' where tenant_id=$1`, p.TenantID)
			case "connection-org", "scope-changed":
				execSQL(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-2','store-2','Synthetic other','store')`, p.TenantID)
				if mode == "connection-org" {
					execSQL(`update integration.provider_connection set organization_id='store-2' where tenant_id=$1`, p.TenantID)
				} else {
					execSQL(`update crm.lead set organization_id='store-2' where tenant_id=$1`, p.TenantID)
				}
			case "retention":
				r.retention = strings.Repeat("b", 64)
			case "profile":
				r.observer.Profile = append(r.observer.Profile, ' ')
			case "signature":
				r.observer.Secrets = statusSecretHook(func(context.Context) (string, error) { return "wrong-secret", nil })
			case "payload":
				execSQL(`update integration.webhook_event set payload=jsonb_set(payload,'{body}','"e30="') where tenant_id=$1`, p.TenantID)
			case "count":
				execSQL(`update integration.webhook_event set payload=jsonb_set(payload,'{verified_event_count}','2') where tenant_id=$1`, p.TenantID)
			case "missing-receipt":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "changed-receipt":
				if err := os.WriteFile(path, []byte("{}"), 0600); err != nil {
					t.Fatal(err)
				}
			case "wrong-key":
				r.key = []byte(strings.Repeat("x", 32))
			case "busy":
				r.busy <- struct{}{}
				defer func() { <-r.busy }()
			case "ambiguous":
				otherEvent := "00000000-0000-0000-0000-000000000002"
				otherKey := AppointmentConfirmationDeliveryKey(p.TenantID, otherEvent)
				execSQL(`insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject) values($1,$2,$3,'requested','confirmed','synthetic-duplicate')`, p.TenantID, otherEvent, c.AppointmentID)
				execSQL(`insert into communication.whatsapp_appointment_approval select tenant_id,$3,channel_code,appointment_id,appointment_version,organization_id,$2,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_id,consent_purpose,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,approved_at,expires_at from communication.whatsapp_appointment_approval where tenant_id=$1`, p.TenantID, otherEvent, otherKey)
				execSQL(`insert into communication.outbound_delivery select tenant_id,channel_code,$2,request_sha256_hex,recipient_hmac,state,attempt_count,locked_until,provider_message_hmac,evidence_sha256_hex,failure_code,accepted_at,created_at,updated_at from communication.outbound_delivery where tenant_id=$1`, p.TenantID, otherKey)
			}
			n, err := r.ObserveRetained(context.Background(), p, event)
			if err == nil || n != 0 {
				t.Fatalf("unsafe routing accepted: %d %v", n, err)
			}
			var count int
			if err := o.Approvals.pool.QueryRow(context.Background(), `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, originalTenant).Scan(&count); err != nil || count != 0 {
				t.Fatalf("partial write: %d %v", count, err)
			}
		})
	}
}

func TestStatusRouterProjectionBudgetAndExactKeys(t *testing.T) {
	batch := statusBatch("delivered", "1603086314", nil)
	ids, err := verifiedStatusIdentities(batch.Body, 1, 1)
	if err != nil || len(ids) != 1 || ids[0].Message != "wamid.synthetic" {
		t.Fatal(ids, err)
	}
	if _, err := verifiedStatusIdentities(batch.Body, 2, 1); err == nil {
		t.Fatal("count mismatch accepted")
	}
	if _, err := verifiedStatusIdentities(batch.Body, 1, 0); err == nil {
		t.Fatal("budget ignored")
	}
	upper := alterStatusBatch(t, batch, func(body map[string]any) {
		v := statusValue(body)
		v["Statuses"] = v["statuses"]
		delete(v, "statuses")
	})
	if _, err := verifiedStatusIdentities(upper.Body, 1, 1); err == nil {
		t.Fatal("case-only provider key accepted")
	}
	if _, err := NewStatusRouter(nil, "wa-primary", strings.Repeat("a", 64), []byte(strings.Repeat("x", 32)), 1); err == nil {
		t.Fatal("nil observer accepted")
	}
}

func TestStatusRouterMultipleMessagesRecoverPartialObservation(t *testing.T) {
	batch := alterStatusBatch(t, statusBatch("delivered", "1603086314", nil), func(body map[string]any) {
		v := statusValue(body)
		v["statuses"] = append(v["statuses"].([]any), map[string]any{"id": "wamid.zz-second", "recipient_id": "5491112345678", "status": "read", "timestamp": "1603086315"})
	})
	r, o, pool, p, c, event, path := routerFixture(t, batch)
	ctx := context.Background()
	otherEvent := "00000000-0000-0000-0000-000000000003"
	otherKey := AppointmentConfirmationDeliveryKey(p.TenantID, otherEvent)
	// A second explicitly synthetic, durably anchored send receipt. This is not
	// evidence of a second Meta API call or of live provider acceptance.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var receipt map[string]any
	if json.Unmarshal(raw, &receipt) != nil {
		t.Fatal("synthetic receipt JSON")
	}
	receipt["message_id_sha256"] = digest([]byte("wamid.zz-second"))
	secondRaw, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	messageHMAC, err := contactidentity.ExternalDigest(r.key, p.TenantID, "whatsapp", "wamid.zz-second")
	if err != nil {
		t.Fatal(err)
	}
	execSQL := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	execSQL(`insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject) values($1,$2,$3,'requested','confirmed','synthetic-second')`, p.TenantID, otherEvent, c.AppointmentID)
	execSQL(`insert into communication.whatsapp_appointment_approval select tenant_id,$3,channel_code,appointment_id,appointment_version,organization_id,$2,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_id,consent_purpose,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,approved_at,expires_at from communication.whatsapp_appointment_approval where tenant_id=$1`, p.TenantID, otherEvent, otherKey)
	execSQL(`insert into communication.outbound_delivery select tenant_id,channel_code,$2,request_sha256_hex,recipient_hmac,state,attempt_count,locked_until,$3,$4,failure_code,accepted_at,created_at,updated_at from communication.outbound_delivery where tenant_id=$1`, p.TenantID, otherKey, messageHMAC, digest(secondRaw))
	directory := filepath.Join(o.Process.EvidenceDirectory, digest([]byte(p.TenantID+"\x00"+otherKey)))
	if err := os.MkdirAll(directory, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "SEND_RECEIPT.json"), secondRaw, 0600); err != nil {
		t.Fatal(err)
	}
	// Signature validation, first observer, then fail the second observer.
	calls := 0
	r.observer.Secrets = statusSecretHook(func(context.Context) (string, error) {
		calls++
		if calls == 3 {
			return "", errors.New("synthetic secret source unavailable")
		}
		return "test-app-secret", nil
	})
	n, err := r.ObserveRetained(ctx, p, event)
	if err == nil || n != 1 || calls != 3 {
		t.Fatalf("partial result not preserved: %d %d %v", n, calls, err)
	}
	var pending bool
	if err := pool.QueryRow(ctx, `select state='received' and processed_at is null from integration.webhook_event where tenant_id=$1 and provider_event_id=$2`, p.TenantID, event).Scan(&pending); err != nil || !pending {
		t.Fatalf("premature inbox completion: %v %v", pending, err)
	}
	r.observer.Secrets = appSecretFixture{}
	for _, expected := range []int{1, 0} {
		n, err := r.ObserveRetained(ctx, p, event)
		if err != nil || n != expected {
			t.Fatalf("recovery: %d %v", n, err)
		}
	}
	var total int
	if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, p.TenantID).Scan(&total); err != nil || total != 2 {
		t.Fatalf("duplicated/lost observations: %d %v", total, err)
	}
}
