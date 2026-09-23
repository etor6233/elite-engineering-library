package whatsappbridge

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"elite.local/enterprise/internal/platform/postgres"
)

type verifyTokenFixture struct{}

func (verifyTokenFixture) WhatsAppVerifyToken(context.Context) (string, error) {
	return "test-verify-token", nil
}

func receiverConfig(o *StatusObserver) WebhookReceiverConfig {
	return WebhookReceiverConfig{TenantID: o.TenantID, ConnectionID: "wa-primary", RetentionApprovalSHA256: strings.Repeat("a", 64), Profile: o.Profile, Process: o.Process, Secrets: appSecretFixture{}, Verification: verifyTokenFixture{}, Store: postgres.NewProviderIntegration(o.Approvals.pool), MaxConcurrent: 4}
}

func webhookRequest(t *testing.T, server *httptest.Server, batch SignedStatusWebhook) (*http.Response, []byte) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(batch.Body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature-256", batch.Signature)
	res, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return res, body
}

func TestWebhookReceiverPostgresReceiptAndObservation(t *testing.T) {
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
	res, err := server.Client().Get(server.URL + "?hub.mode=subscribe&hub.verify_token=test-verify-token&hub.challenge=12345")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil || res.StatusCode != 200 || string(raw) != "12345" {
		t.Fatal("challenge", res.StatusCode, string(raw), err)
	}
	batch := statusBatch("delivered", "1603086314", nil)
	// Drop the transport only after the receiver has produced a successful ACK.
	// Its 200 was written to a recorder, never sent to the provider-side client.
	lost := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured := httptest.NewRecorder()
		h.ServeHTTP(captured, r)
		if captured.Code != 200 {
			http.Error(w, "fixture did not commit", 500)
			return
		}
		connection, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error("fixture hijack", err)
			return
		}
		_ = connection.Close()
	}))
	defer lost.Close()
	lostRequest, _ := http.NewRequest(http.MethodPost, lost.URL, bytes.NewReader(batch.Body))
	lostRequest.Header.Set("Content-Type", "application/json")
	lostRequest.Header.Set("X-Hub-Signature-256", batch.Signature)
	if lostResponse, err := lost.Client().Do(lostRequest); err == nil {
		lostResponse.Body.Close()
		t.Fatal("fixture unexpectedly delivered acknowledgement")
	}
	var committed int
	if err := pool.QueryRow(context.Background(), `select count(*) from integration.webhook_event where tenant_id=$1`, p.TenantID).Scan(&committed); err != nil || committed != 1 {
		t.Fatal("lost ACK without commit", committed, err)
	}
	var wg sync.WaitGroup
	results := make(chan int, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(batch.Body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Hub-Signature-256", batch.Signature)
			response, err := server.Client().Do(req)
			if err != nil {
				results <- 0
				return
			}
			defer response.Body.Close()
			_, _ = io.Copy(io.Discard, response.Body)
			results <- response.StatusCode
		}()
	}
	wg.Wait()
	close(results)
	for status := range results {
		if status != 200 {
			t.Fatalf("concurrent receipt status %d", status)
		}
	}
	var events, jobs int
	if err := pool.QueryRow(context.Background(), `select (select count(*) from integration.webhook_event where tenant_id=$1 and state='received' and processed_at is null),(select count(*) from platform.job where tenant_id=$1 and job_type='provider.webhook.received' and completed_at is null)`, p.TenantID).Scan(&events, &jobs); err != nil || events != 1 || jobs != 1 {
		t.Fatal("not durable or duplicated", events, jobs, err)
	}
	var stored []byte
	if err := pool.QueryRow(context.Background(), `select payload from integration.webhook_event where tenant_id=$1 and connection_id='wa-primary'`, p.TenantID).Scan(&stored); err != nil {
		t.Fatal(err)
	}
	var retained struct {
		Schema        string `json:"schema"`
		Body          []byte `json:"body"`
		Signature     string `json:"signature"`
		ProfileHash   string `json:"profile_sha256"`
		RetentionHash string `json:"retention_approval_sha256"`
		Count         int    `json:"verified_event_count"`
	}
	if json.Unmarshal(stored, &retained) != nil || retained.Schema != "elite-whatsapp-retained-webhook/v1" || !bytes.Equal(retained.Body, batch.Body) || retained.Signature != batch.Signature || retained.ProfileHash != digest(o.Profile) || retained.RetentionHash != config.RetentionApprovalSHA256 || retained.Count != 1 {
		t.Fatal("retained original differs")
	}
	if strings.Contains(string(stored), "test-app-secret") || strings.Contains(string(stored), "test-verify-token") {
		t.Fatal("secret persisted")
	}
	// Exercise the existing observer using the actual reloaded signed original.
	// Target routing is supplied by this synthetic fixture, not an automatic worker.
	if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{{Body: retained.Body, Signature: retained.Signature}}); n != 1 || err != nil {
		t.Fatal("observation from retained original", n, err)
	}
	value, err := o.Approvals.ReadNotificationStatus(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID)
	if err != nil || value.DeliveryStatus != "observed_delivered" {
		t.Fatal(value, err)
	}
	// A lost HTTP acknowledgement is recoverable by replaying the exact original.
	res, raw = webhookRequest(t, server, batch)
	if res.StatusCode != 200 || string(raw) != "EVENT_RECEIVED\n" || res.Header.Get("Cache-Control") != "no-store" {
		t.Fatal(res.StatusCode, string(raw))
	}
	if _, err := pool.Exec(context.Background(), `update integration.provider_connection set state='disabled' where tenant_id=$1`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	res, _ = webhookRequest(t, server, batch)
	if res.StatusCode != 503 {
		t.Fatal("disabled acknowledged", res.StatusCode)
	}
	if err := pool.QueryRow(context.Background(), `select (select count(*) from integration.webhook_event where tenant_id=$1),(select count(*) from platform.job where tenant_id=$1 and job_type='provider.webhook.received')`, p.TenantID).Scan(&events, &jobs); err != nil || events != 1 || jobs != 1 {
		t.Fatal(events, jobs, err)
	}
}

func TestWebhookReceiverPostgresRejectsBeforeReceipt(t *testing.T) {
	pool := approvalPool(t)
	o, p, _, _ := statusObserverFixture(t, pool)
	if _, err := pool.Exec(context.Background(), `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'wa-primary','meta-whatsapp','META_APP_SECRET')`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	config := receiverConfig(o)
	bad := config
	bad.RetentionApprovalSHA256 = ""
	if _, err := NewWebhookReceiver(bad); err == nil {
		t.Fatal("retention approval missing")
	}
	bad = config
	bad.MaxConcurrent = 0
	if _, err := NewWebhookReceiver(bad); err == nil {
		t.Fatal("missing concurrency budget")
	}
	h, err := NewWebhookReceiver(config)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(h)
	defer server.Close()
	for _, mode := range []string{"tampered", "foreign_waba", "duplicate_signature", "media", "encoding", "oversize", "query", "invalid_json", "runtime_hash", "busy", "get_token", "get_duplicate", "get_malformed", "method"} {
		t.Run(mode, func(t *testing.T) {
			batch := statusBatch("read", "1603086315", nil)
			req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(batch.Body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Hub-Signature-256", batch.Signature)
			want := 403
			sign := func(raw []byte) {
				mac := hmac.New(sha256.New, []byte("test-app-secret"))
				mac.Write(raw)
				req.Body = io.NopCloser(bytes.NewReader(raw))
				req.ContentLength = int64(len(raw))
				req.Header.Set("X-Hub-Signature-256", "sha256="+hex.EncodeToString(mac.Sum(nil)))
			}
			switch mode {
			case "tampered":
				req.Header.Set("X-Hub-Signature-256", "sha256="+strings.Repeat("0", 64))
			case "foreign_waba":
				sign(bytes.ReplaceAll(batch.Body, []byte("987654321"), []byte("987654322")))
			case "duplicate_signature":
				req.Header.Add("X-Hub-Signature-256", batch.Signature)
				want = 400
			case "media":
				req.Header.Set("Content-Type", "text/plain")
				want = 415
			case "encoding":
				req.Header.Set("Content-Encoding", "gzip")
				want = 415
			case "oversize":
				req.Body = io.NopCloser(strings.NewReader(strings.Repeat("x", (1<<20)+1)))
				req.ContentLength = (1 << 20) + 1
				want = 413
			case "query":
				req.URL.RawQuery = "tenant_id=other"
				want = 400
			case "invalid_json":
				sign([]byte(`{"object":`))
			case "runtime_hash":
				previous := h.config.Process.AdapterSHA256
				h.config.Process.AdapterSHA256 = strings.Repeat("0", 64)
				defer func() { h.config.Process.AdapterSHA256 = previous }()
			case "busy":
				for i := 0; i < cap(h.slots); i++ {
					h.slots <- struct{}{}
				}
				defer func() {
					for i := 0; i < cap(h.slots); i++ {
						<-h.slots
					}
				}()
				want = 503
			case "get_token":
				req.Method = http.MethodGet
				req.URL.RawQuery = "hub.mode=subscribe&hub.verify_token=wrong&hub.challenge=12345"
			case "get_duplicate":
				req.Method = http.MethodGet
				req.URL.RawQuery = "hub.mode=subscribe&hub.verify_token=test-verify-token&hub.challenge=1&hub.challenge=2"
				want = 400
			case "get_malformed":
				req.Method = http.MethodGet
				req.URL.RawQuery = "hub.mode=subscribe&hub.verify_token=%XX&hub.challenge=1"
				want = 400
			case "method":
				req.Method = http.MethodPut
				want = 405
			}
			res, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			raw, err := io.ReadAll(res.Body)
			if err != nil || res.StatusCode != want || string(raw) != "WEBHOOK_NOT_ACCEPTED\n" {
				t.Fatal(mode, res.StatusCode, string(raw), err)
			}
		})
	}
	var count int
	if err := pool.QueryRow(context.Background(), `select count(*) from integration.webhook_event where tenant_id=$1`, p.TenantID).Scan(&count); err != nil || count != 0 {
		t.Fatal("rejected callback persisted", count, err)
	}
}
