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
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/paymentbridge"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/workers"
)

type recoveryFenceCrash struct {
	*db.PaymentDispatchFence
	fail bool
}

func (f *recoveryFenceCrash) ReconcileAccepted(ctx context.Context, m channels.Message, hash string, r outbounddelivery.Receipt) error {
	if f.fail {
		f.fail = false
		return errors.New("fixture crash before accepted receipt")
	}
	return f.PaymentDispatchFence.ReconcileAccepted(ctx, m, hash, r)
}

func runPreferenceRecovery(t *testing.T, scenario string) {
	t.Helper()
	ctx := context.Background()
	pool := connectedPool(t)
	tenant, payment := connectedInputs(t, pool, "mercadopago")
	var mu sync.Mutex
	posts, paymentGets, searchGets, preferenceGets := 0, 0, 0, 0
	expires := time.Time{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer TEST-fixture" {
			t.Error("wrong credential")
		}
		switch r.URL.Path {
		case "/users/me":
			if r.Method != "GET" {
				t.Error("unexpected write")
			}
			io.WriteString(w, `{"id":123,"country_id":"AR","site_id":"MLA"}`)
		case "/checkout/preferences":
			if r.Method != "POST" {
				t.Error("unexpected method")
			}
			posts++
			var body struct {
				ExternalReference string         `json:"external_reference"`
				Metadata          map[string]any `json:"metadata"`
				ExpirationDateTo  time.Time      `json:"expiration_date_to"`
			}
			if json.NewDecoder(r.Body).Decode(&body) != nil || body.ExternalReference != "order" || body.Metadata["payment_attempt_id"] != payment.ID || body.ExpirationDateTo.IsZero() {
				t.Error("request binding missing")
			}
			expires = body.ExpirationDateTo
			// Provider accepted the request but its response never reaches the worker
			// as a valid SDK result. Recovery must not issue another POST.
			io.WriteString(w, `{"incomplete":`)
		case "/v1/payments/987":
			paymentGets++
			if r.Method != "GET" {
				t.Error("unexpected payment effect")
			}
			io.WriteString(w, `{"id":987,"transaction_amount":1234.56,"transaction_amount_refunded":0,"currency_id":"ARS","external_reference":"order","status":"approved","captured":true,"collector_id":123,"live_mode":false}`)
		case "/checkout/preferences/search":
			searchGets++
			q := r.URL.Query()
			if r.Method != "GET" || len(q) != 4 || q.Get("external_reference") != "order" || q.Get("site_id") != "MLA" || q.Get("limit") != "2" || q.Get("offset") != "0" {
				t.Error("unbounded or unscoped search")
			}
			if scenario == "missing" {
				io.WriteString(w, `{"total":0,"next_offset":0,"elements":[]}`)
				return
			}
			if scenario == "ambiguous" {
				io.WriteString(w, `{"total":2,"next_offset":2,"elements":[{"id":"123-first"},{"id":"123-second"}]}`)
				return
			}
			collector, mode, expiration := 123, false, expires
			if scenario == "account" {
				collector = 456
			}
			if scenario == "mode" {
				mode = true
			}
			if scenario == "expiry" {
				expiration = expiration.Add(time.Second)
			}
			fmt.Fprintf(w, `{"total":1,"next_offset":1,"elements":[{"id":"123-preference","external_reference":"order","collector_id":%d,"site_id":"MLA","live_mode":%t,"expires":true,"expiration_date_to":%q}]}`, collector, mode, expiration.Format(time.RFC3339))
		case "/checkout/preferences/123-preference":
			preferenceGets++
			if r.Method != "GET" {
				t.Error("preference mutation")
			}
			attempt, amount := payment.ID, 1234.56
			if scenario == "metadata" {
				attempt = "other-attempt"
			}
			if scenario == "amount" {
				amount = 1234.57
			}
			fmt.Fprintf(w, `{"id":"123-preference","external_reference":"order","metadata":{"order_id":"order","payment_attempt_id":%q},"collector_id":123,"site_id":"MLA","expires":true,"items":[{"currency_id":"ARS","quantity":1,"unit_price":%.2f}],"sandbox_init_point":"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-preference","expiration_date_to":%q}`, attempt, amount, expires.Format(time.RFC3339))
		default:
			t.Errorf("unexpected HTTP path %s", r.URL.Path)
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
	baseFence, err := db.NewOutboundDeliveryStore(pool, []byte(strings.Repeat("r", 32)), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	fence := &recoveryFenceCrash{PaymentDispatchFence: &db.PaymentDispatchFence{OutboundDeliveryStore: baseFence, Scope: scope}, fail: scenario == "saved-crash"}
	worker := paymentbridge.Worker{Scope: scope, Store: store, Fence: fence, Driver: driver}
	if n, err := worker.ProcessOnce(ctx, 1); err != nil || n != 0 {
		t.Fatal("lost response", n, err)
	}
	callback, err := db.NewPaymentCallbackProcessor(pool, worker, "mp-recovery")
	if err != nil {
		t.Fatal(err)
	}
	processor := workers.JobProcessor{Store: &db.PaymentCallbackJobs{Jobs: db.NewJobs(pool), Scope: scope}, Handler: db.PaymentCallbackRecoveryProcessor{Base: callback}, Queue: "payment-provider-events", WorkerID: "mp-recovery", Lease: 20 * time.Second, RetryDelay: 0, BatchSize: 1}
	webhook, err := paymentbridge.NewWebhook(paymentbridge.WebhookConfig{TenantID: tenant, ConnectionID: "checkout", ProviderCode: "mercadopago", Tolerance: time.Minute, MaxConcurrent: 1, Secrets: connectedSecret{}, Inbox: store})
	if err != nil {
		t.Fatal(err)
	}
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	mac := hmac.New(sha256.New, []byte("whsec_fixture"))
	fmt.Fprintf(mac, "id:987;request-id:mp-recovery;ts:%s;", ts)
	request := httptest.NewRequest("POST", "https://api.example.test/payment?data.id=987&type=payment", strings.NewReader(`{"type":"payment","data":{"id":"987"}}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Request-Id", "mp-recovery")
	request.Header.Set("X-Signature", "ts="+ts+",v1="+hex.EncodeToString(mac.Sum(nil)))
	response := httptest.NewRecorder()
	webhook.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal(response.Code)
	}
	completed, err := processor.ProcessOnce(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var identity, state, financial string
	snapshot := func() {
		t.Helper()
		if err := pool.QueryRow(ctx, `select coalesce(b.session_id,''),f.state,p.state from payment.provider_checkout b join payment.payment_attempt p using(tenant_id,payment_attempt_id) join communication.outbound_delivery f on f.tenant_id=b.tenant_id and f.delivery_key=b.payment_attempt_id and f.channel_code='payment_mercadopago' where b.tenant_id=$1`, tenant).Scan(&identity, &state, &financial); err != nil {
			t.Fatal(err)
		}
	}
	snapshot()
	if financial != "captured" {
		t.Fatal("payment not observed through GET", financial)
	}
	switch scenario {
	case "unique":
		if completed != 1 || identity != "123-preference" || state != "accepted" {
			t.Fatal("identity recovery", completed, identity, state)
		}
	case "saved-crash":
		if completed != 0 || identity != "123-preference" || state != "unknown" {
			t.Fatal("crash point", completed, identity, state)
		}
		if completed, err = processor.ProcessOnce(ctx); err != nil || completed != 1 {
			t.Fatal("saved identity retry", completed, err)
		}
		snapshot()
		if state != "accepted" {
			t.Fatal("recovered receipt missing", state)
		}
	default:
		if completed != 0 || identity != "" || state != "unknown" {
			t.Fatal("unproven identity accepted", completed, identity, state)
		}
	}
	if n, err := worker.ProcessOnce(ctx, 1); err != nil || n != 0 {
		t.Fatal("recovery resent checkout", n, err)
	}
	var moneyEffects int
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.provider-observed'`, tenant).Scan(&moneyEffects); err != nil || moneyEffects != 1 {
		t.Fatal("duplicate money observation effect", moneyEffects, err)
	}
	mu.Lock()
	postCount, payCount, searchCount, getCount := posts, paymentGets, searchGets, preferenceGets
	mu.Unlock()
	if postCount != 1 || searchCount != 1 {
		t.Fatal("unexpected provider mutations/search", postCount, searchCount)
	}
	if scenario == "unique" && (payCount != 2 || getCount != 1) {
		t.Fatal("SDK query chain", payCount, getCount)
	}
	if scenario == "saved-crash" && (payCount != 3 || getCount != 2) {
		t.Fatal("retry should GET saved session without repeat search", payCount, getCount)
	}
	t.Logf("MP_PREFERENCE_RECOVERY_PASS scenario=%s posts=%d payment_gets=%d search_gets=%d preference_gets=%d state=%s", scenario, postCount, payCount, searchCount, getCount, state)
}

func TestMercadoPagoPreferenceRecoveryLostPost(t *testing.T) { runPreferenceRecovery(t, "unique") }
func TestMercadoPagoPreferenceRecoveryAfterSavedIdentityCrash(t *testing.T) {
	runPreferenceRecovery(t, "saved-crash")
}
func TestMercadoPagoPreferenceRecoveryRejectsUnprovenIdentity(t *testing.T) {
	for _, scenario := range []string{"missing", "ambiguous", "metadata", "amount", "account", "mode", "expiry"} {
		t.Run(scenario, func(t *testing.T) { runPreferenceRecovery(t, scenario) })
	}
}
