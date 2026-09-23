package whatsappbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestNotificationStatusReadOnlyStates(t *testing.T) {
	pool := approvalPool(t)
	verifier, sign := notificationIssuer(t)
	config, err := pgxpool.ParseConfig(os.Getenv("ELITE_WHATSAPP_TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["default_transaction_read_only"] = "on"
	readPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer readPool.Close()
	var readOnly string
	if err := readPool.QueryRow(context.Background(), "show transaction_read_only").Scan(&readOnly); err != nil || readOnly != "on" {
		t.Fatal("read-only test session not demonstrated", readOnly, err)
	}
	for _, mode := range []string{"not_started", "sending", "stale_sending", "accepted", "unknown", "failed_terminal", "expired_accepted", "cancelled", "revoked", "request_mismatch", "recipient_mismatch"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			p, c, writer := approvalFixtureData(t, pool)
			if mode == "expired_accepted" {
				c.ExpiresAt = time.Now().Add(500 * time.Millisecond)
			}
			if _, err := writer.Approve(ctx, p, c); err != nil {
				t.Fatal(err)
			}
			store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
			if err != nil {
				t.Fatal(err)
			}
			hash, err := outbounddelivery.MessageSHA256(c.Message)
			if err != nil {
				t.Fatal(err)
			}
			want := mode
			if mode != "not_started" {
				if _, err := store.Claim(ctx, c.Message, hash); err != nil {
					t.Fatal(err)
				}
			}
			switch mode {
			case "unknown":
				if err := store.MarkUnknown(ctx, c.Message, hash, "PROVIDER_CALL_UNCERTAIN"); err != nil {
					t.Fatal(err)
				}
			case "failed_terminal":
				if err := store.MarkFailed(ctx, c.Message, hash, strings.Repeat("e", 64), "PROVEN_REJECTED"); err != nil {
					t.Fatal(err)
				}
			case "stale_sending":
				want = "sending"
				if _, err := pool.Exec(ctx, `update communication.outbound_delivery set locked_until=clock_timestamp()-interval '1 second' where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "accepted", "expired_accepted", "cancelled", "revoked", "request_mismatch", "recipient_mismatch":
				want = "accepted"
				if err := store.Complete(ctx, c.Message, hash, outbounddelivery.Receipt{ProviderMessageID: "wamid.status-fixture", EvidenceSHA256: strings.Repeat("a", 64), AcceptedAt: time.Now()}); err != nil {
					t.Fatal(err)
				}
				queries := map[string]string{
					"cancelled":          `update crm.appointment set state='cancelled',version=version+1 where tenant_id=$1`,
					"revoked":            `update communication.contact_channel_binding set state='revoked',pii_allowed=false,version=version+1 where tenant_id=$1`,
					"request_mismatch":   `update communication.outbound_delivery set request_sha256_hex=repeat('f',64) where tenant_id=$1`,
					"recipient_mismatch": `update communication.outbound_delivery set recipient_hmac=repeat('f',64) where tenant_id=$1`,
				}
				if q := queries[mode]; q != "" {
					if _, err := pool.Exec(ctx, q, p.TenantID); err != nil {
						t.Fatal(err)
					}
				}
				if mode == "expired_accepted" {
					time.Sleep(550 * time.Millisecond)
				}
			}
			reader, err := NewPostgresAppointmentApprovals(readPool, []byte("0123456789abcdef0123456789abcdef"), "fixture-appointment-notification", "policy-1")
			if err != nil {
				t.Fatal(err)
			}
			sender, _, marker := senderFixture(t, "accepted")
			sender.TenantID = p.TenantID
			module, err := NewAppointmentNotificationModule(reader, sender, store, receiverFixture{})
			if err != nil {
				t.Fatal(err)
			}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			var before int
			if err := pool.QueryRow(ctx, `select count(*) from communication.outbound_delivery_event where tenant_id=$1`, p.TenantID).Scan(&before); err != nil {
				t.Fatal(err)
			}
			for i := 0; i < 2; i++ {
				req := httptest.NewRequest(http.MethodGet, "/v1/franchise/appointments/appointment-1/whatsapp-confirmation?organization_id=store-1&confirmation_event_id="+c.ConfirmationEventID, nil)
				req.Header.Set("Authorization", "Bearer "+sign(p))
				response := httptest.NewRecorder()
				mux.ServeHTTP(response, req)
				if mode == "request_mismatch" || mode == "recipient_mismatch" {
					if response.Code != 409 || !strings.Contains(response.Body.String(), "NOTIFICATION_EVIDENCE_MISMATCH") {
						t.Fatal(response.Code, response.Body.String())
					}
					continue
				}
				if response.Code != 200 || response.Header().Get("Cache-Control") != "no-store" {
					t.Fatal(response.Code, response.Body.String())
				}
				var value NotificationStatus
				if err := json.Unmarshal(response.Body.Bytes(), &value); err != nil {
					t.Fatal(err)
				}
				if value.FenceState != want || value.DeliveryStatus != "not_observed_by_this_reader" || value.DeliveryKey != c.Message.DeliveryKey || value.ObservedAt.IsZero() {
					t.Fatal(value)
				}
				if value.ReconciliationRequired != (mode == "unknown" || mode == "stale_sending") {
					t.Fatal(value)
				}
				if value.ApprovalExpired != (mode == "expired_accepted") {
					t.Fatal(value)
				}
				if (value.AcceptedAt != nil) != (want == "accepted") {
					t.Fatal(value)
				}
				for _, secret := range []string{c.Message.ExternalID, p.Subject, c.ConsentID, c.EvidenceSHA256, "wamid.status-fixture", "appointment_confirmed", "synthetic-token"} {
					if strings.Contains(response.Body.String(), secret) {
						t.Fatal("sensitive field exposed")
					}
				}
			}
			var after int
			if err := pool.QueryRow(ctx, `select count(*) from communication.outbound_delivery_event where tenant_id=$1`, p.TenantID).Scan(&after); err != nil || after != before {
				t.Fatal("GET mutated events", before, after, err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("GET invoked provider")
			}
			if mode == "stale_sending" {
				var state string
				if err := pool.QueryRow(ctx, `select state from communication.outbound_delivery where tenant_id=$1`, p.TenantID).Scan(&state); err != nil || state != "sending" {
					t.Fatal("GET mutated stale fence", state, err)
				}
			}
		})
	}
}

func TestNotificationStatusScopeAndUnavailable(t *testing.T) {
	pool := approvalPool(t)
	verifier, sign := notificationIssuer(t)
	for _, mode := range []string{"no_token", "role", "organization", "tenant", "appointment", "event", "missing_grant", "reassigned", "duplicate_query", "unknown_query", "bad_uuid", "body", "unavailable"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			p, c, resolver := approvalFixtureData(t, pool)
			if mode != "missing_grant" {
				if _, err := resolver.Approve(ctx, p, c); err != nil {
					t.Fatal(err)
				}
			}
			sender, _, marker := senderFixture(t, "accepted")
			sender.TenantID = p.TenantID
			store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
			if err != nil {
				t.Fatal(err)
			}
			if mode == "unavailable" {
				closed := approvalPool(t)
				closed.Close()
				resolver, err = NewPostgresAppointmentApprovals(closed, []byte("0123456789abcdef0123456789abcdef"), "fixture-appointment-notification", "policy-1")
				if err != nil {
					t.Fatal(err)
				}
			}
			module, err := NewAppointmentNotificationModule(resolver, sender, store, receiverFixture{})
			if err != nil {
				t.Fatal(err)
			}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			query := url.Values{"organization_id": {"store-1"}, "confirmation_event_id": {c.ConfirmationEventID}}
			pathID := "appointment-1"
			want := 404
			var body *strings.Reader = strings.NewReader("")
			switch mode {
			case "no_token":
				want = 401
			case "role":
				p.Permissions = map[string]struct{}{}
				want = 403
			case "organization":
				p.Organizations = map[string]struct{}{}
				want = 403
			case "tenant":
				p.TenantID = "other-tenant"
				want = 403
			case "appointment":
				pathID = "other"
			case "event":
				query.Set("confirmation_event_id", "00000000-0000-0000-0000-000000000000")
			case "reassigned":
				if _, err := pool.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'other-lead','store-1','new','fixture','{}',true)`, p.TenantID); err != nil {
					t.Fatal(err)
				}
				if _, err := pool.Exec(ctx, `update crm.appointment set lead_id='other-lead' where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "duplicate_query":
				query.Add("organization_id", "store-1")
				want = 400
			case "unknown_query":
				query.Set("tenant_id", p.TenantID)
				want = 400
			case "bad_uuid":
				query.Set("confirmation_event_id", "not-an-event")
				want = 400
			case "body":
				body = strings.NewReader("{}")
				want = 400
			case "unavailable":
				want = 503
			}
			req := httptest.NewRequest(http.MethodGet, "/v1/franchise/appointments/"+pathID+"/whatsapp-confirmation?"+query.Encode(), body)
			if mode != "no_token" {
				req.Header.Set("Authorization", "Bearer "+sign(p))
			}
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, req)
			if response.Code != want || response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("Content-Type") != "application/problem+json" {
				t.Fatal(response.Code, response.Body.String())
			}
			if strings.Contains(response.Body.String(), "closed") || strings.Contains(response.Body.String(), c.Message.ExternalID) {
				t.Fatal("internal error leaked")
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("GET invoked child")
			}
		})
	}
}

func TestNotificationPostThenStatus(t *testing.T) {
	pool := approvalPool(t)
	p, c, resolver := approvalFixtureData(t, pool)
	verifier, sign := notificationIssuer(t)
	sender, _, marker := senderFixture(t, "accepted")
	sender.TenantID = p.TenantID
	store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	module, err := NewAppointmentNotificationModule(resolver, sender, store, receiverFixture{})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	module.Register(mux, verifier)
	server := httptest.NewServer(mux)
	defer server.Close()
	body, _ := json.Marshal(notificationRequest(c))
	token := sign(p)
	for _, method := range []string{http.MethodPost, http.MethodGet, http.MethodGet} {
		path := "/v1/franchise/appointments/appointment-1/whatsapp-confirmation"
		data := []byte(nil)
		if method == http.MethodPost {
			data = body
		} else {
			path += "?organization_id=store-1&confirmation_event_id=" + c.ConfirmationEventID
		}
		req, err := http.NewRequest(method, server.URL+path, bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		response, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != 200 {
			response.Body.Close()
			t.Fatal(response.StatusCode)
		}
		if method == http.MethodGet {
			var value NotificationStatus
			err = json.NewDecoder(response.Body).Decode(&value)
			if err != nil || value.FenceState != "accepted" || value.DeliveryStatus != "not_observed_by_this_reader" {
				response.Body.Close()
				t.Fatal(value, err)
			}
		}
		response.Body.Close()
	}
	calls, err := os.ReadFile(marker)
	if err != nil || string(calls) != "call\n" {
		t.Fatal("read caused duplicate", string(calls), err)
	}
}
