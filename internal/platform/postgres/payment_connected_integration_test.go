package postgres_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/paymentbridge"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/platform/workers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

type connectedTransport struct {
	destination *url.URL
	host        string
}

func (s connectedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Scheme != "https" || r.URL.Host != s.host {
		return nil, fmt.Errorf("fixture denies remote destination")
	}
	copy := r.Clone(r.Context())
	u := *r.URL
	u.Scheme = s.destination.Scheme
	u.Host = s.destination.Host
	copy.URL = &u
	copy.Host = ""
	return http.DefaultTransport.RoundTrip(copy)
}

type connectedSecret struct{}

func (connectedSecret) PaymentWebhookSecret(context.Context, string, string) (string, error) {
	return "whsec_fixture", nil
}

type connectedProvider struct {
	mu                              sync.Mutex
	t                               *testing.T
	attempt, account, order         string
	expires                         int64
	posts, sessionGets, paymentGets int
	amount                          int64
	refund                          int64
	loseResponse                    bool
}

func (f *connectedProvider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	amount := f.amount
	if amount == 0 {
		amount = 123456
	}
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer sk_test_fixture" {
		f.t.Error("missing exact fixture credential")
	}
	switch r.URL.Path {
	case "/v1/account":
		if r.Method != "GET" {
			f.t.Error("account write")
		}
		fmt.Fprintf(w, `{"id":%q}`, f.account)
	case "/v1/checkout/sessions", "/v1/checkout/sessions/cs_test_fixture":
		if r.Method == "POST" {
			f.posts++
			if err := r.ParseForm(); err != nil {
				f.t.Error(err)
			}
			f.expires, _ = strconv.ParseInt(r.Form.Get("expires_at"), 10, 64)
			if r.Form.Get("metadata[payment_attempt_id]") != f.attempt || r.Form.Get("payment_intent_data[metadata][payment_attempt_id]") != f.attempt || r.Form.Get("metadata[order_id]") != "order" || r.Form.Get("line_items[0][price_data][unit_amount]") != strconv.FormatInt(amount, 10) || r.Header.Get("Idempotency-Key") != f.attempt {
				f.t.Error("checkout escaped order contract")
			}
			if f.loseResponse {
				_, _ = io.WriteString(w, `{"truncated":`)
				return
			}
		} else if r.Method == "GET" {
			f.sessionGets++
		} else {
			f.t.Error("invalid checkout method")
		}
		intent := `"pi_fixture"`
		checkoutURL := "null"
		if r.Method == "POST" {
			intent = "null"
			checkoutURL = `"https://checkout.stripe.com/c/pay/cs_test_fixture"`
		}
		fmt.Fprintf(w, `{"id":"cs_test_fixture","object":"checkout.session","mode":"payment","amount_total":%d,"currency":"ars","client_reference_id":%q,"metadata":{"order_id":%q,"payment_attempt_id":%q},"status":"complete","payment_status":"paid","livemode":false,"url":%s,"payment_intent":%s,"expires_at":%d}`, amount, f.attempt, f.order, f.attempt, checkoutURL, intent, f.expires)
	case "/v1/payment_intents/pi_fixture":
		f.paymentGets++
		if r.Method != "GET" || r.URL.Query().Get("expand[0]") != "latest_charge" {
			f.t.Error("payment must be expanded official GET")
		}
		fmt.Fprintf(w, `{"id":"pi_fixture","amount":%d,"amount_received":%d,"currency":"ars","status":"succeeded","livemode":false,"metadata":{"order_id":%q,"payment_attempt_id":%q},"latest_charge":{"id":"ch_fixture","amount":%d,"amount_captured":%d,"amount_refunded":%d,"currency":"ars","captured":true,"disputed":false,"livemode":false}}`, amount, amount, f.order, f.attempt, amount, amount, f.refund)
	default:
		f.t.Errorf("unrecognized provider HTTP %s", r.URL.Path)
		w.WriteHeader(500)
	}
}

func connectedPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("PAYMENT_CONNECTED_DB_URL not set")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_payment_connected_") {
		t.Fatal("requires disposable loopback connected database")
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}
func connectedInputs(t *testing.T, pool *pgxpool.Pool, provider string) (string, commerce.PaymentAttempt) {
	return connectedInputsWithAllocation(t, pool, provider, nil)
}
func connectedInputsWithAllocation(t *testing.T, pool *pgxpool.Pool, provider string, before func(*testing.T, *pgxpool.Pool, string) int64, quoteHook ...func(*testing.T, *pgxpool.Pool, string) int64) (string, commerce.PaymentAttempt) {
	t.Helper()
	ctx := context.Background()
	tenant := connectedSeedOrder(t, pool, provider, quoteHook...)
	sales := commerce.NewService(db.NewCommerce(pool), randomid.Generator{})
	expectedAmount := int64(123456)
	if before != nil {
		expectedAmount = before(t, pool, tenant)
	}
	payment, err := sales.RequestOrderPayment(ctx, tenant, "store", "order", provider, "connected-payment-key-0001", "operator")
	if err != nil || payment.State != "created" || payment.AmountMinorUnits != expectedAmount {
		t.Fatal(payment, err)
	}
	replay, err := sales.RequestOrderPayment(ctx, tenant, "store", "order", provider, "connected-payment-key-0001", "operator")
	if err != nil || replay.ID != payment.ID {
		t.Fatal("request replay", replay, err)
	}
	return tenant, payment
}

type connectedRun struct {
	pool      *pgxpool.Pool
	tenant    string
	payment   commerce.PaymentAttempt
	provider  *connectedProvider
	worker    paymentbridge.Worker
	handler   *db.PaymentCallbackProcessor
	processor workers.JobProcessor
	webhook   http.Handler
}

func newConnectedRun(t *testing.T, loss bool) *connectedRun {
	return newConnectedRunWithAllocation(t, loss, nil)
}
func newConnectedRunWithAllocation(t *testing.T, loss bool, before func(*testing.T, *pgxpool.Pool, string) int64, quoteHook ...func(*testing.T, *pgxpool.Pool, string) int64) *connectedRun {
	t.Helper()
	ctx := context.Background()
	pool := connectedPool(t)
	tenant, payment := connectedInputsWithAllocation(t, pool, "stripe", before, quoteHook...)
	f := &connectedProvider{t: t, attempt: payment.ID, account: "acct_fixture", order: "order", loseResponse: loss, amount: payment.AmountMinorUnits}
	server := httptest.NewServer(f)
	t.Cleanup(server.Close)
	destination, _ := url.Parse(server.URL)
	c := paymentbridge.Scope{TenantID: tenant, OrganizationID: "store", ConnectionID: "checkout", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}
	driver, err := paymentbridge.NewSDKDriver(paymentbridge.SDKDriverConfig{Scope: c, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", DisplayName: "Synthetic order"}, "sk_test_fixture", &http.Client{Transport: connectedTransport{destination, "api.stripe.com"}})
	if err != nil {
		t.Fatal(err)
	}
	store := db.NewPaymentCheckoutStore(pool)
	fence, err := db.NewOutboundDeliveryStore(pool, []byte(strings.Repeat("f", 32)), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	worker := paymentbridge.Worker{Scope: c, Store: store, Fence: &db.PaymentDispatchFence{OutboundDeliveryStore: fence, Scope: c}, Driver: driver}
	handler, err := db.NewPaymentCallbackProcessor(pool, worker, "connected-callback")
	if err != nil {
		t.Fatal(err)
	}
	jobs := &db.PaymentCallbackJobs{Jobs: db.NewJobs(pool), Scope: c}
	processor := workers.JobProcessor{Store: jobs, Handler: handler, Queue: "payment-provider-events", WorkerID: "connected-callback", Lease: 20 * time.Second, RetryDelay: time.Second, BatchSize: 1}
	webhook, err := paymentbridge.NewWebhook(paymentbridge.WebhookConfig{TenantID: tenant, ConnectionID: "checkout", ProviderCode: "stripe", Tolerance: time.Minute, MaxConcurrent: 4, Secrets: connectedSecret{}, Inbox: store})
	if err != nil {
		t.Fatal(err)
	}
	n, err := worker.ProcessOnce(ctx, 1)
	if err != nil || (!loss && n != 1) || (loss && n != 0) {
		t.Fatalf("dispatch %d %v", n, err)
	}
	if n, err = worker.ProcessOnce(ctx, 1); err != nil || n != 0 {
		t.Fatal("duplicate dispatch", n, err)
	}
	var state string
	if err = pool.QueryRow(ctx, `select state from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2`, tenant, payment.ID).Scan(&state); err != nil || state != "pending" {
		t.Fatal("unobserved money state", state, err)
	}
	return &connectedRun{pool: pool, tenant: tenant, payment: payment, provider: f, worker: worker, handler: handler, processor: processor, webhook: webhook}
}
func (r *connectedRun) callback(t *testing.T, event, kind, resource string, valid bool) int {
	t.Helper()
	object := "checkout.session"
	extra := ""
	if kind == "payment_intent.succeeded" {
		object = "payment_intent"
	}
	if kind == "charge.refunded" {
		object = "charge"
		extra = `,"payment_intent":"pi_fixture"`
	}
	body := []byte(fmt.Sprintf(`{"id":%q,"object":"event","api_version":%q,"type":%q,"livemode":false,"data":{"object":{"id":%q,"object":%q%s}}}`, event, stripe.APIVersion, kind, resource, object, extra))
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "whsec_fixture", Timestamp: time.Now()})
	signature := signed.Header
	if !valid {
		signature += "tampered"
	}
	request := httptest.NewRequest("POST", "https://api.example.test/payment", bytesReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Stripe-Signature", signature)
	response := httptest.NewRecorder()
	r.webhook.ServeHTTP(response, request)
	return response.Code
}
func bytesReader(raw []byte) *strings.Reader { return strings.NewReader(string(raw)) }

func TestPaymentConnectedReference(t *testing.T) {
	for _, loss := range []bool{false, true} {
		t.Run(fmt.Sprintf("lost-post-response=%t", loss), func(t *testing.T) {
			r := newConnectedRun(t, loss)
			ctx := context.Background()
			if code := r.callback(t, "evt_bad", "checkout.session.completed", "cs_test_fixture", false); code != 403 {
				t.Fatal("bad signature", code)
			}
			if code := r.callback(t, "evt_initial", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
				t.Fatal("callback", code)
			}
			if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
				t.Fatal("callback job", n, err)
			}
			var state, hash, fence string
			var hold bool
			if err := r.pool.QueryRow(ctx, `select p.state,o.evidence_sha256_hex,o.hold,f.state from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) join communication.outbound_delivery f on f.tenant_id=p.tenant_id and f.delivery_key=p.payment_attempt_id where p.tenant_id=$1`, r.tenant).Scan(&state, &hash, &hold, &fence); err != nil || state != "captured" || hold || fence != "accepted" || len(hash) != 64 {
				t.Fatal("observed projection", state, hold, fence, err)
			}
			if code := r.callback(t, "evt_initial", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
				t.Fatal("callback replay", code)
			}
			if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 0 {
				t.Fatal("replay job effect", n, err)
			}
			repo := db.NewFranchiseJourney(r.pool)
			sum := sha256.Sum256([]byte(franchisejourney.ReferenceHandoverContractDocument))
			contract := franchisejourney.HandoverReleaseContract{ID: "reference-single-unit-observed-payment-v1", DocumentSHA256: hex.EncodeToString(sum[:]), Scope: "LOCAL_FIXTURES", MaximumObservationAge: 5 * time.Minute}
			handover, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, contract)
			if err != nil {
				t.Fatal(err)
			}
			command := franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: r.payment.ID, ObservationSHA256: hash, IdempotencyKey: "connected-handover-key"}
			prepared, replayed, err := handover.Prepare(ctx, r.tenant, "operator", command)
			if err != nil || replayed {
				t.Fatal("prepare", err)
			}
			recovered, err := handover.Result(ctx, r.tenant, "store", "order", command.IdempotencyKey)
			if err != nil || recovered.Handover.ID != prepared.Handover.ID {
				t.Fatal("recover prepare", err)
			}
			ids := randomid.Generator{}
			if _, err = repo.PublishDeliveryChecklist(ctx, r.tenant, "operator", franchisejourney.DeliveryChecklist{ID: "initial-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
				t.Fatal(err)
			}
			presented, err := repo.CompleteDeliveryChecklist(ctx, r.tenant, "store", "operator", prepared.Handover.ID, 1, "initial-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
			if err != nil {
				t.Fatal(err)
			}
			accepted, err := repo.AcceptHandover(ctx, r.tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "initial-checklist", 1, strings.Repeat("e", 64), ids.New())
			if err != nil || accepted.State != "accepted" {
				t.Fatal(err)
			}
			eligible, err := handover.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash)
			if err != nil || !eligible.Eligible {
				t.Fatal("release reference", err)
			}
			r.provider.mu.Lock()
			posts, sessionGets, paymentGets := r.provider.posts, r.provider.sessionGets, r.provider.paymentGets
			r.provider.refund = 1
			r.provider.mu.Unlock()
			if posts != 1 || sessionGets != 1 || paymentGets != 1 {
				t.Fatalf("HTTP effects posts=%d checkoutgets=%d paymentgets=%d", posts, sessionGets, paymentGets)
			}
			if code := r.callback(t, "evt_refund", "charge.refunded", "ch_fixture", true); code != 200 {
				t.Fatal(code)
			}
			if _, err = handover.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash); !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatal("callback must hold before GET", err)
			}
			if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
				t.Fatal("refund reconciliation", n, err)
			}
			if _, err = handover.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash); !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatal("partial refund released", err)
			}
			var count int
			if err = r.pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.provider-observed'`, r.tenant).Scan(&count); err != nil || count != 1 {
				t.Fatal("duplicate money projection", count, err)
			}
			t.Logf("PAYMENT_CONNECTED_REFERENCE_PASS tenant=%s lost_post_response=%t quote_accept_allocate_request_sdk_signed_callback_job_get_handover_accept_release=true provider_posts=1 refund_hold=true production_claim=false", r.tenant, loss)
		})
	}
}

func TestPaymentCallbackRejectsForgedJobAndMetadata(t *testing.T) {
	r := newConnectedRun(t, false)
	ctx := context.Background()
	if err := r.handler.Handle(ctx, db.Job{TenantID: r.tenant, JobID: randomid.Generator{}.New(), Queue: "payment-provider-events", JobType: "provider.webhook.received", SchemaVersion: 1, Attempts: 1, Payload: json.RawMessage(`{"connection_id":"checkout","provider_code":"stripe","provider_event_id":"evt_fake","event_type":"payment.reconciliation-requested.v1"}`)}); err == nil {
		t.Fatal("unclaimed invented job accepted")
	}
	if code := r.callback(t, "evt_wrong", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	r.provider.mu.Lock()
	r.provider.order = "foreign-order"
	r.provider.mu.Unlock()
	if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 0 {
		t.Fatal("metadata mismatch accepted", n, err)
	}
	var captured int
	if err := r.pool.QueryRow(ctx, `select count(*) from payment.payment_attempt where tenant_id=$1 and state='captured'`, r.tenant).Scan(&captured); err != nil || captured != 0 {
		t.Fatal("foreign order paid", captured, err)
	}
	r.provider.mu.Lock()
	gets := r.provider.paymentGets
	posts := r.provider.posts
	r.provider.mu.Unlock()
	if gets != 0 || posts != 1 {
		t.Fatal("foreign metadata reached payment or POST", gets, posts)
	}
	t.Log("PAYMENT_CALLBACK_NEGATIVES_PASS forged_claim=true metadata_binding=true no_money_effect=true")
}

func TestPaymentCallbackFirstIntentAndClaimFencing(t *testing.T) {
	r := newConnectedRun(t, false)
	ctx := context.Background()
	if code := r.callback(t, "evt_first_intent", "payment_intent.succeeded", "pi_fixture", true); code != 200 {
		t.Fatal(code)
	}
	jobs := r.processor.Store
	claimed, err := jobs.Claim(ctx, "payment-provider-events", r.processor.WorkerID, 20*time.Second, 1)
	if err != nil || len(claimed) != 1 {
		t.Fatal("claim", err)
	}
	original := claimed[0]
	for name, edit := range map[string]func(*db.Job){
		"previous-generation": func(j *db.Job) { j.Attempts++ },
		"other-tenant":        func(j *db.Job) { j.TenantID = randomid.Generator{}.New() },
		"different-queue":     func(j *db.Job) { j.Queue = "provider-events" },
		"changed-payload": func(j *db.Job) {
			j.Payload = json.RawMessage(`{"connection_id":"checkout","provider_code":"stripe","provider_event_id":"evt_not_received","event_type":"payment.reconciliation-requested.v1"}`)
		},
	} {
		t.Run(name, func(t *testing.T) {
			modified := original
			edit(&modified)
			if err := r.handler.Handle(ctx, modified); err == nil {
				t.Fatal("invalid claim accepted")
			}
		})
	}
	r.provider.mu.Lock()
	gets := r.provider.paymentGets
	r.provider.mu.Unlock()
	if gets != 0 {
		t.Fatal("invalid claims issued GET")
	}
	if err = r.handler.Handle(ctx, original); err != nil {
		t.Fatal("first PI callback", err)
	}
	if err = jobs.Complete(ctx, original.TenantID, original.JobID, r.processor.WorkerID, original.Attempts); err != nil {
		t.Fatal(err)
	}
	if err = r.handler.Handle(ctx, original); err == nil {
		t.Fatal("completed claim accepted")
	}
	r.provider.mu.Lock()
	gets = r.provider.paymentGets
	posts := r.provider.posts
	r.provider.mu.Unlock()
	if gets != 2 || posts != 1 {
		t.Fatal("initial metadata GET and fresh reconciliation GET", gets, posts)
	}
	var state string
	var evidence int
	if err = r.pool.QueryRow(ctx, `select state,(select count(*) from payment.provider_observation_event where tenant_id=$1) from payment.payment_attempt where tenant_id=$1`, r.tenant).Scan(&state, &evidence); err != nil || state != "captured" || evidence != 1 {
		t.Fatal("first PI projection", state, evidence, err)
	}
	t.Log("PAYMENT_CALLBACK_FIRST_INTENT_AND_CLAIM_FENCING_PASS metadata_lookup_get=true fresh_reconciliation_get=true claim_fields=4 completed_claim_rejected=true")
}

func TestPaymentConnectedMercadoPagoWithoutInheritedMetadata(t *testing.T) {
	pool := connectedPool(t)
	ctx := context.Background()
	tenant, payment := connectedInputs(t, pool, "mercadopago")
	var mu sync.Mutex
	posts, gets := 0, 0
	externalReference := "order"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer TEST-fixture" {
			t.Error("MP credential mismatch")
		}
		switch r.URL.Path {
		case "/users/me":
			if r.Method != "GET" {
				t.Error("user write")
			}
			io.WriteString(w, `{"id":123,"country_id":"AR","site_id":"MLA"}`)
		case "/checkout/preferences":
			if r.Method != "POST" {
				t.Error("unexpected preference method")
			}
			posts++
			var body struct {
				ExternalReference string         `json:"external_reference"`
				Metadata          map[string]any `json:"metadata"`
				ExpirationDateTo  string         `json:"expiration_date_to"`
			}
			if json.NewDecoder(r.Body).Decode(&body) != nil || body.ExternalReference != "order" || body.Metadata["payment_attempt_id"] != payment.ID || body.ExpirationDateTo == "" {
				t.Error("MP request binding missing")
			}
			fmt.Fprintf(w, `{"id":"123-preference","external_reference":"order","metadata":{"order_id":"order","payment_attempt_id":%q},"collector_id":123,"items":[{"currency_id":"ARS","quantity":1,"unit_price":1234.56}],"sandbox_init_point":"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-preference","expiration_date_to":%q}`, payment.ID, body.ExpirationDateTo)
		case "/v1/payments/987":
			if r.Method != "GET" {
				t.Error("payment write")
			}
			gets++
			// Deliberately no metadata: a payment need not inherit preference metadata.
			fmt.Fprintf(w, `{"id":987,"transaction_amount":1234.56,"transaction_amount_refunded":0,"currency_id":"ARS","external_reference":%q,"status":"approved","captured":true,"collector_id":123,"live_mode":false}`, externalReference)
		default:
			t.Errorf("unexpected MP endpoint %s", r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	t.Cleanup(server.Close)
	destination, _ := url.Parse(server.URL)
	scope := paymentbridge.Scope{TenantID: tenant, OrganizationID: "store", ConnectionID: "checkout", ProviderCode: "mercadopago", AccountRef: "123", Currency: "ARS", MinorUnitExponent: 2}
	driver, err := paymentbridge.NewSDKDriver(paymentbridge.SDKDriverConfig{Scope: scope, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", NotificationURL: "https://api.example.test/payment", DisplayName: "Synthetic order"}, "TEST-fixture", &http.Client{Transport: connectedTransport{destination, "api.mercadopago.com"}})
	if err != nil {
		t.Fatal(err)
	}
	store := db.NewPaymentCheckoutStore(pool)
	base, err := db.NewOutboundDeliveryStore(pool, []byte(strings.Repeat("m", 32)), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	worker := paymentbridge.Worker{Scope: scope, Store: store, Fence: &db.PaymentDispatchFence{OutboundDeliveryStore: base, Scope: scope}, Driver: driver}
	if n, err := worker.ProcessOnce(ctx, 1); err != nil || n != 1 {
		t.Fatal("MP checkout", n, err)
	}
	handler, err := db.NewPaymentCallbackProcessor(pool, worker, "mp-callback")
	if err != nil {
		t.Fatal(err)
	}
	processor := workers.JobProcessor{Store: &db.PaymentCallbackJobs{Jobs: db.NewJobs(pool), Scope: scope}, Handler: handler, Queue: "payment-provider-events", WorkerID: "mp-callback", Lease: 20 * time.Second, RetryDelay: time.Second, BatchSize: 1}
	webhook, err := paymentbridge.NewWebhook(paymentbridge.WebhookConfig{TenantID: tenant, ConnectionID: "checkout", ProviderCode: "mercadopago", Tolerance: time.Minute, MaxConcurrent: 1, Secrets: connectedSecret{}, Inbox: store})
	if err != nil {
		t.Fatal(err)
	}
	callback := func(requestID string) {
		t.Helper()
		ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
		mac := hmac.New(sha256.New, []byte("whsec_fixture"))
		fmt.Fprintf(mac, "id:987;request-id:%s;ts:%s;", requestID, ts)
		// Body financial claims are unsigned and must not supply payment facts.
		request := httptest.NewRequest("POST", "https://api.example.test/payment?data.id=987&type=payment", strings.NewReader(`{"type":"payment","data":{"id":"987"},"status":"rejected","transaction_amount":999999}`))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Request-Id", requestID)
		request.Header.Set("X-Signature", "ts="+ts+",v1="+hex.EncodeToString(mac.Sum(nil)))
		response := httptest.NewRecorder()
		webhook.ServeHTTP(response, request)
		if response.Code != 200 {
			t.Fatal("MP signed callback", response.Code)
		}
	}
	// A genuine signed payment ID with a different GET order cannot bind.
	mu.Lock()
	externalReference = "foreign-order"
	mu.Unlock()
	callback("mp-wrong-order")
	if n, err := processor.ProcessOnce(ctx); err != nil || n != 0 {
		t.Fatal("foreign MP order accepted", n, err)
	}
	var captures int
	if err = pool.QueryRow(ctx, `select count(*) from payment.payment_attempt where tenant_id=$1 and state='captured'`, tenant).Scan(&captures); err != nil || captures != 0 {
		t.Fatal("foreign MP capture", captures, err)
	}
	mu.Lock()
	externalReference = "order"
	mu.Unlock()
	callback("mp-correct-order")
	if n, err := processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatal("MP reconciliation", n, err)
	}
	var hash, state string
	var hold bool
	if err = pool.QueryRow(ctx, `select p.state,o.evidence_sha256_hex,o.hold from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) where p.tenant_id=$1`, tenant).Scan(&state, &hash, &hold); err != nil || state != "captured" || hold {
		t.Fatal("MP observed", state, hold, err)
	}
	mu.Lock()
	postCount, getCount := posts, gets
	mu.Unlock()
	if postCount != 1 || getCount != 3 {
		t.Fatal("MP SDK effect count", postCount, getCount)
	}
	repo := db.NewFranchiseJourney(pool)
	sum := sha256.Sum256([]byte(franchisejourney.ReferenceHandoverContractDocument))
	contract := franchisejourney.HandoverReleaseContract{ID: "reference-single-unit-observed-payment-v1", DocumentSHA256: hex.EncodeToString(sum[:]), Scope: "LOCAL_FIXTURES", MaximumObservationAge: 5 * time.Minute}
	service, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, contract)
	if err != nil {
		t.Fatal(err)
	}
	prepared, _, err := service.Prepare(ctx, tenant, "operator", franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: payment.ID, ObservationSHA256: hash, IdempotencyKey: "mp-connected-handover-key"})
	if err != nil {
		t.Fatal("MP prepare", err)
	}
	ids := randomid.Generator{}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "mp-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	presented, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", prepared.Handover.ID, 1, "mp-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "mp-checklist", 1, strings.Repeat("e", 64), ids.New()); err != nil {
		t.Fatal(err)
	}
	release, err := service.EvaluateRelease(ctx, tenant, "store", prepared.Handover.ID, hash)
	if err != nil || !release.Eligible {
		t.Fatal("MP release check", err)
	}
	t.Logf("PAYMENT_CONNECTED_MP_PASS tenant=%s metadata_absent=true wrong_order_rejected=true official_get_money=true quote_to_release=true provider_posts=1", tenant)
}

func TestPaymentMaterializedHandoverProfile(t *testing.T) {
	r := newConnectedRun(t, false)
	ctx := context.Background()
	if code := r.callback(t, "evt_materialized", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	var hash string
	if err := r.pool.QueryRow(ctx, `select evidence_sha256_hex from payment.provider_observation where tenant_id=$1`, r.tenant).Scan(&hash); err != nil {
		t.Fatal(err)
	}
	python := os.Getenv("HANDOVER_PROFILE_PYTHON")
	if python == "" {
		t.Fatal("explicit admitted Python materializer runtime required")
	}
	load := func(mode string, refs ...string) franchisejourney.HandoverReleaseContract {
		t.Helper()
		account := "acct_fixture"
		if len(refs) > 0 {
			account = refs[0]
		}
		out := filepath.Join(t.TempDir(), "policy")
		cmd := exec.CommandContext(ctx, python, "-X", "utf8", "-B", filepath.Join("..", "..", "..", "tools", "materialize_handover_profile.py"), "--output", out, "--profile-id", "franchise-integral", "--revision", "1", "--tenant-id", r.tenant, "--organization-id", "store", "--expected-mode", mode, "--payment-provider", "stripe", "--payment-account-ref", account, "--payment-connection-id", "checkout", "--activate")
		if body, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("materializer %v %s", err, body)
		}
		raw, err := os.ReadFile(filepath.Join(out, "activation.json"))
		if err != nil {
			t.Fatal(err)
		}
		var activation franchisejourney.HandoverProfileActivation
		if err = json.Unmarshal(raw, &activation); err != nil {
			t.Fatal(err)
		}
		contract, err := franchisejourney.LoadHandoverProfileFile(filepath.Join(out, "profile.json"), activation)
		if err != nil {
			t.Fatal(err)
		}
		return contract
	}
	command := franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: r.payment.ID, ObservationSHA256: hash, IdempotencyKey: "materialized-profile-key"}
	repo := db.NewFranchiseJourney(r.pool)
	live := load("live")
	liveService, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, live)
	if err != nil {
		t.Fatal("supported future mode cannot load", err)
	}
	if _, _, err = liveService.Prepare(ctx, r.tenant, "operator", command); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("sandbox evidence admitted under live contract", err)
	}
	foreignAccount := load("sandbox", "acct_other")
	foreignService, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, foreignAccount)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = foreignService.Prepare(ctx, r.tenant, "operator", command); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("different provider account admitted", err)
	}
	policy := load("sandbox")
	service, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = r.pool.Exec(ctx, `update integration.provider_connection set state='disabled' where tenant_id=$1`, r.tenant); err != nil {
		t.Fatal(err)
	}
	if _, _, err = service.Prepare(ctx, r.tenant, "operator", command); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("disabled profile connection admitted", err)
	}
	if _, err = r.pool.Exec(ctx, `update integration.provider_connection set state='active' where tenant_id=$1`, r.tenant); err != nil {
		t.Fatal(err)
	}
	prepared, _, err := service.Prepare(ctx, r.tenant, "operator", command)
	if err != nil || prepared.ContractID != "franchise-integral@1" || prepared.ContractSHA256 != policy.DocumentSHA256 {
		t.Fatal("materialized profile preparation", prepared, err)
	}
	ids := randomid.Generator{}
	if _, err = repo.PublishDeliveryChecklist(ctx, r.tenant, "operator", franchisejourney.DeliveryChecklist{ID: "profile-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	presented, err := repo.CompleteDeliveryChecklist(ctx, r.tenant, "store", "operator", prepared.Handover.ID, 1, "profile-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AcceptHandover(ctx, r.tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "profile-checklist", 1, strings.Repeat("e", 64), ids.New()); err != nil {
		t.Fatal(err)
	}
	eligible, err := service.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash)
	if err != nil || !eligible.Eligible || eligible.Scope != "MATERIALIZED_PROFILE" {
		t.Fatal("materialized release", eligible, err)
	}
	if _, err = liveService.EvaluateRelease(ctx, r.tenant, "store", prepared.Handover.ID, hash); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("profile hash changed after preparation", err)
	}
	if _, err = service.Result(ctx, randomid.Generator{}.New(), "store", "order", command.IdempotencyKey); !errors.Is(err, franchisejourney.ErrReleaseConditioned) {
		t.Fatal("cross profile scope", err)
	}
	t.Logf("MATERIALIZED_HANDOVER_PROFILE_PG_PASS tenant=%s materializer_file_loader=true same_owners=true live_mode_mismatch_denied=true profile_hash_persisted=true live_proven=false", r.tenant)
}

func connectedSeedOrder(t *testing.T, pool *pgxpool.Pool, provider string, quoteHook ...func(*testing.T, *pgxpool.Pool, string) int64) string {
	t.Helper()
	ctx := context.Background()
	ids := randomid.Generator{}
	tenant := ids.New()
	statements := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'connected-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','fixture@example.test')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','model','new','fixture','{}')`,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'retail','AR','ARS',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','variant',123456,'inclusive')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SERIAL-SYNTHETIC','available',1,clock_timestamp())`,
	}
	for _, q := range statements {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, `insert into integration.provider_connection(tenant_id,connection_id,provider_code,organization_id,secret_ref)values($1,'checkout',$2,'store','FIXTURE_ONLY_NOT_A_SECRET')`, tenant, provider); err != nil {
		t.Fatal(err)
	}
	repo := db.NewFranchiseJourney(pool)
	if _, _, err := repo.CreateQuoteAs(ctx, tenant, "quote-key", franchisejourney.Quote{ID: "quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().UTC().Add(time.Hour)}, strings.Repeat("a", 64), ids.New(), "operator"); err != nil {
		t.Fatal(err)
	}
	version := int64(1)
	if len(quoteHook) > 1 {
		t.Fatal("one quote hook maximum")
	}
	if len(quoteHook) == 1 {
		version = quoteHook[0](t, pool, tenant)
	}
	if _, err := repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", version, strings.Repeat("a", 64), "order", "line", ids.New(), ids.New()); err != nil {
		t.Fatal(err)
	}
	sales := commerce.NewService(db.NewCommerce(pool), ids)
	if err := sales.AllocateStockAs(ctx, tenant, "store", "order", "line", "stock", 1, 1, "operator"); err != nil {
		t.Fatal(err)
	}
	return tenant
}
