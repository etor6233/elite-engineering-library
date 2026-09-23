package whatsappbridge

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type appSecretFixture struct{}

type statusSecretHook func(context.Context) (string, error)

func (f statusSecretHook) WhatsAppAppSecret(ctx context.Context) (string, error) { return f(ctx) }

func (appSecretFixture) WhatsAppAppSecret(context.Context) (string, error) {
	return "test-app-secret", nil
}

func statusObserverFixture(t *testing.T, pool *pgxpool.Pool) (*StatusObserver, identity.Principal, AppointmentApprovalCommand, []byte) {
	t.Helper()
	python := os.Getenv("ELITE_WHATSAPP_PYTHON")
	if python == "" {
		t.Fatal("ELITE_WHATSAPP_PYTHON required for actual isolated status verifier")
	}
	directory, err := filepath.Abs("../../whatsapp_cloud")
	if err != nil {
		t.Fatal(err)
	}
	// Reuse the Python send-bridge fixture; its only provider transport is local.
	code := `import sys,json,base64
sys.path.insert(0,sys.argv[1])
from test_status_reconciliation import StatusReconciliationTests,profile
c=StatusReconciliationTests();c.setUp()
try: print(json.dumps({"profile":profile(),"receipt":base64.b64encode(c.receipt).decode(),"message":__import__('test_whatsapp_cloud').message()}))
finally: c.doCleanups()
`
	raw, err := exec.Command(python, "-I", "-B", "-c", code, directory).Output()
	if err != nil {
		t.Fatal("synthetic send fixture", err)
	}
	var data struct {
		Profile json.RawMessage `json:"profile"`
		Receipt []byte          `json:"receipt"`
		Message json.RawMessage `json:"message"`
	}
	if json.Unmarshal(raw, &data) != nil {
		t.Fatal("fixture JSON")
	}
	p, c, approvals := approvalFixtureData(t, pool)
	c.Message.Text = string(data.Message)
	c.ProfileSHA256 = digest(data.Profile)
	if _, err = approvals.Approve(context.Background(), p, c); err != nil {
		t.Fatal(err)
	}
	store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	hash, _ := outbounddelivery.MessageSHA256(c.Message)
	if _, err = store.Claim(context.Background(), c.Message, hash); err != nil {
		t.Fatal(err)
	}
	if err = store.Complete(context.Background(), c.Message, hash, outbounddelivery.Receipt{ProviderMessageID: "wamid.synthetic", EvidenceSHA256: digest(data.Receipt), AcceptedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	fileHash := func(path string) string {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return digest(b)
	}
	o := &StatusObserver{TenantID: p.TenantID, Profile: data.Profile, Approvals: approvals, Secrets: appSecretFixture{}, ReconcilerSHA256: fileHash(filepath.Join(directory, "status_reconciliation.py")), Process: Process{PythonExecutable: python, PythonSHA256: fileHash(python), AdapterDirectory: directory, AdapterSHA256: fileHash(filepath.Join(directory, "whatsapp_cloud.py")), EvidenceDirectory: t.TempDir()}}
	return o, p, c, data.Receipt
}

func statusBatch(status, timestamp string, extra map[string]any) SignedStatusWebhook {
	event := map[string]any{"id": "wamid.synthetic", "recipient_id": "5491112345678", "status": status, "timestamp": timestamp}
	for k, v := range extra {
		event[k] = v
	}
	raw, _ := json.Marshal(map[string]any{"object": "whatsapp_business_account", "entry": []any{map[string]any{"id": "987654321", "changes": []any{map[string]any{"field": "messages", "value": map[string]any{"messaging_product": "whatsapp", "metadata": map[string]any{"phone_number_id": "123456789"}, "statuses": []any{event}}}}}}})
	mac := hmac.New(sha256.New, []byte("test-app-secret"))
	mac.Write(raw)
	return SignedStatusWebhook{Body: raw, Signature: "sha256=" + hex.EncodeToString(mac.Sum(nil))}
}

func TestStatusObserverPostgresConcurrentAndRead(t *testing.T) {
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	batches := []SignedStatusWebhook{statusBatch("read", "1603086315", nil), statusBatch("sent", "1603086313", nil), statusBatch("delivered", "1603086314", nil)}
	var wg sync.WaitGroup
	errs := make(chan error, 4)
	counts := make(chan int, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, e := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches)
			counts <- n
			errs <- e
		}()
	}
	wg.Wait()
	close(errs)
	close(counts)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	total := 0
	for n := range counts {
		total += n
	}
	if total != 3 {
		t.Fatalf("inserted=%d", total)
	}
	value, err := o.Approvals.ReadNotificationStatus(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID)
	if err != nil || value.FenceState != "accepted" || value.DeliveryStatus != "observed_read" || value.ProviderEventCount != 3 || value.ProviderTimestamp == nil || *value.ProviderTimestamp != 1603086315 {
		t.Fatal(value, err)
	}
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
	readApprovals := *o.Approvals
	readApprovals.pool = readPool
	module := &AppointmentNotificationModule{approvals: &readApprovals, sender: &Sender{TenantID: p.TenantID}}
	mux := http.NewServeMux()
	module.registerNotificationStatus(mux, verifier)
	server := httptest.NewServer(mux)
	defer server.Close()
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/v1/franchise/appointments/"+c.AppointmentID+"/whatsapp-confirmation?organization_id="+c.OrganizationID+"&confirmation_event_id="+c.ConfirmationEventID, nil)
	req.Header.Set("Authorization", "Bearer "+sign(p))
	response, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var httpValue NotificationStatus
	if response.StatusCode != 200 || response.Header.Get("Cache-Control") != "no-store" || json.NewDecoder(response.Body).Decode(&httpValue) != nil || httpValue.DeliveryStatus != "observed_read" || httpValue.ProviderEventCount != 3 {
		t.Fatal("HTTP read-only observation", httpValue, response.StatusCode)
	}
	if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{statusBatch("failed", "1603086315", nil)}); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	value, err = o.Approvals.ReadNotificationStatus(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID)
	if err != nil || value.DeliveryStatus != "ambiguous_latest_timestamp" || value.ProviderEventCount != 4 {
		t.Fatal(value, err)
	}
	var events, attempts int
	if err = pool.QueryRow(context.Background(), `select attempt_count,(select count(*) from communication.outbound_delivery_event e where e.tenant_id=d.tenant_id) from communication.outbound_delivery d where tenant_id=$1`, p.TenantID).Scan(&attempts, &events); err != nil || attempts != 1 || events != 2 {
		t.Fatal("send fence changed", attempts, events, err)
	}
	// Evidence is append-only and does not leak into the public response.
	if _, err = pool.Exec(context.Background(), `update communication.whatsapp_status_observation set provider_status='sent' where tenant_id=$1`, p.TenantID); err == nil {
		t.Fatal("immutable evidence updated")
	}
	encoded, _ := json.Marshal(value)
	for _, private := range []string{"5491112345678", "wamid.synthetic", "test-app-secret", digest(receipt)} {
		if strings.Contains(string(encoded), private) {
			t.Fatal("private data in response")
		}
	}
}

func TestStatusObserverPostgresRejections(t *testing.T) {
	pool := approvalPool(t)
	for _, mode := range []string{"signature", "receipt", "tenant", "permission", "organization", "recipient", "adapter", "reconciler", "unknown", "reassigned", "profile"} {
		t.Run(mode, func(t *testing.T) {
			o, p, c, receipt := statusObserverFixture(t, pool)
			batches := []SignedStatusWebhook{statusBatch("delivered", "1603086314", nil)}
			switch mode {
			case "signature":
				batches[0].Signature = "sha256=" + strings.Repeat("0", 64)
			case "receipt":
				receipt = append(receipt, ' ')
			case "tenant":
				p.TenantID = "00000000-0000-0000-0000-000000000000"
			case "permission":
				p.Permissions = map[string]struct{}{}
			case "organization":
				c.OrganizationID = "other"
			case "recipient":
				batches[0] = statusBatch("delivered", "1603086314", map[string]any{"recipient_id": "5491112345679"})
			case "adapter":
				o.Process.AdapterSHA256 = strings.Repeat("0", 64)
			case "reconciler":
				o.ReconcilerSHA256 = strings.Repeat("0", 64)
			case "profile":
				o.Profile = json.RawMessage(`{"wrong":true}`)
			case "unknown":
				if _, err := pool.Exec(context.Background(), `update communication.outbound_delivery set state='unknown' where tenant_id=$1`, o.TenantID); err != nil {
					t.Fatal(err)
				}
			case "reassigned":
				if _, err := pool.Exec(context.Background(), `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'other','other','Other','store')`, o.TenantID); err != nil {
					t.Fatal(err)
				}
				if _, err := pool.Exec(context.Background(), `update crm.lead set organization_id='other' where tenant_id=$1`, o.TenantID); err != nil {
					t.Fatal(err)
				}
			}
			if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches); !errors.Is(err, ErrStatusObservation) || n != 0 {
				t.Fatal(n, err)
			}
			var count int
			if err := pool.QueryRow(context.Background(), `select count(*) from communication.whatsapp_status_batch where tenant_id=$1`, o.TenantID).Scan(&count); err != nil || count != 0 {
				t.Fatal("partial evidence", count, err)
			}
		})
	}
}

func TestStatusObserverPostgresConflictRollsBackBatch(t *testing.T) {
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	original := statusBatch("failed", "1603086314", map[string]any{"errors": []any{map[string]any{"code": 1}}})
	if _, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{original}); err != nil {
		t.Fatal(err)
	}
	divergent := statusBatch("failed", "1603086314", map[string]any{"errors": []any{map[string]any{"code": 2}}})
	if _, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{statusBatch("sent", "1603086313", nil), divergent}); err == nil {
		t.Fatal("divergence accepted")
	}
	var events, batches int
	if err := pool.QueryRow(context.Background(), `select (select count(*) from communication.whatsapp_status_observation where tenant_id=$1),(select count(*) from communication.whatsapp_status_batch where tenant_id=$1)`, p.TenantID).Scan(&events, &batches); err != nil || events != 1 || batches != 1 {
		t.Fatal("not atomic", events, batches, err)
	}
}

// The secret is requested after the first database anchor lookup and before the
// isolated verifier runs. This hook changes durable state in that exact window;
// no sleeps, replacement verifier or weakened production guard are involved.
func TestStatusObserverPostgresRechecksAfterVerification(t *testing.T) {
	pool := approvalPool(t)
	for _, mode := range []string{"lead_scope", "fence_state", "receipt_anchor"} {
		t.Run(mode, func(t *testing.T) {
			o, p, c, receipt := statusObserverFixture(t, pool)
			called := false
			var mutationErr error
			o.Secrets = statusSecretHook(func(ctx context.Context) (secret string, err error) {
				defer func() { mutationErr = err }()
				called = true
				switch mode {
				case "lead_scope":
					if _, err := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'other','other','Other','store')`, p.TenantID); err != nil {
						return "", err
					}
					if _, err := pool.Exec(ctx, `update crm.lead set organization_id='other' where tenant_id=$1`, p.TenantID); err != nil {
						return "", err
					}
				case "fence_state":
					if _, err := pool.Exec(ctx, `update communication.outbound_delivery set state='unknown' where tenant_id=$1`, p.TenantID); err != nil {
						return "", err
					}
				case "receipt_anchor":
					if _, err := pool.Exec(ctx, `update communication.outbound_delivery set evidence_sha256_hex=repeat('f',64) where tenant_id=$1`, p.TenantID); err != nil {
						return "", err
					}
				}
				return "test-app-secret", nil
			})
			n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{statusBatch("read", "1603086315", nil)})
			if !called || mutationErr != nil || n != 0 || !errors.Is(err, ErrStatusObservation) {
				t.Fatal("stale authority accepted or mutation not exercised", called, mutationErr, n, err)
			}
			assertStatusRows(t, pool, p.TenantID, 0, 0)
		})
	}
}

func assertStatusRows(t *testing.T, pool *pgxpool.Pool, tenant string, wantEvents, wantBatches int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var events, batches int
	err := pool.QueryRow(ctx, `select (select count(*) from communication.whatsapp_status_observation where tenant_id=$1),(select count(*) from communication.whatsapp_status_batch where tenant_id=$1)`, tenant).Scan(&events, &batches)
	if err != nil || events != wantEvents || batches != wantBatches {
		t.Fatal("unexpected durable observations", events, batches, err)
	}
}

func TestStatusObserverPostgresCanceledWriteAndRecovery(t *testing.T) {
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	batches := []SignedStatusWebhook{statusBatch("read", "1603086315", nil)}
	lock, err := pool.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Rollback(context.Background())
	// Block the event INSERT after its parent batch has been inserted. The lock
	// is confined to the dedicated loopback test database guarded by approvalPool.
	if _, err := lock.Exec(context.Background(), `lock table communication.whatsapp_status_observation in exclusive mode`); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	type outcome struct {
		n   int
		err error
	}
	done := make(chan outcome, 1)
	go func() {
		n, err := o.Observe(ctx, p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches)
		done <- outcome{n, err}
	}()
	probe, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		var blocked bool
		err := pool.QueryRow(probe, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like 'insert into communication.whatsapp_status_observation%')`).Scan(&blocked)
		if err != nil {
			t.Fatal("could not prove blocked INSERT", err)
		}
		if blocked {
			break
		}
		select {
		case early := <-done:
			t.Fatal("observer ended before blocked INSERT", early)
		case <-probe.Done():
			t.Fatal("blocked INSERT not observed")
		case <-ticker.C:
		}
	}
	cancel()
	select {
	case result := <-done:
		if result.n != 0 || !errors.Is(result.err, ErrStatusObservation) {
			t.Fatal("canceled write reported success", result)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancellation did not release observer")
	}
	if err := lock.Rollback(context.Background()); err != nil {
		t.Fatal(err)
	}
	assertStatusRows(t, pool, p.TenantID, 0, 0)
	if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches); n != 1 || err != nil {
		t.Fatal("recovery failed", n, err)
	}
	assertStatusRows(t, pool, p.TenantID, 1, 1)
}

func TestStatusObserverPostgresReplayAcrossConnectionRestart(t *testing.T) {
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	batches := []SignedStatusWebhook{statusBatch("delivered", "1603086314", nil)}
	if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches); n != 1 || err != nil {
		t.Fatal(n, err)
	}
	pool.Close()
	// Fresh database connections and service values; no in-memory dedup cache.
	reopened := approvalPool(t)
	approvals := *o.Approvals
	approvals.pool = reopened
	observer := *o
	observer.Approvals = &approvals
	if n, err := observer.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches); n != 0 || err != nil {
		t.Fatal("replay after connection restart", n, err)
	}
	assertStatusRows(t, reopened, p.TenantID, 1, 1)
	var attempts int
	if err := reopened.QueryRow(context.Background(), `select attempt_count from communication.outbound_delivery where tenant_id=$1`, p.TenantID).Scan(&attempts); err != nil || attempts != 1 {
		t.Fatal("replay changed send attempts", attempts, err)
	}
}
