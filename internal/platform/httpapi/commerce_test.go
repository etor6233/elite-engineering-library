package httpapi

import (
	"context"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type commerceRepo struct{ books int }

func TestPaymentRequestsDisabledWithoutSelection(t *testing.T) {
	for _, provider := range []string{"", "unsupported"} {
		mux := http.NewServeMux()
		p := identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{"payment:create": {}, "payment:write": {}}, Organizations: map[string]struct{}{"org": {}}}
		CommerceModule{Service: commerce.NewService(&commerceRepo{}, &commerceIDs{}), PaymentProvider: provider}.Register(mux, paymentPrincipal{p})
		for _, path := range []string{"/v1/commerce/orders/order/payment-request", "/v1/commerce/orders/order/payments", "/v1/payments/payment/transitions"} {
			r := httptest.NewRequest("POST", path, strings.NewReader(`{}`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", "Bearer synthetic")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != 503 || !strings.Contains(w.Body.String(), "PAYMENT_REQUEST_DISABLED") {
				t.Fatalf("disabled route %s status%d", path, w.Code)
			}
		}
	}
}

// The HTTP boundary uses a synthetic verified principal; real OIDC is covered
// by the connected portal gate, not claimed by this focused payment regression.
type paymentPrincipal struct{ identity.Principal }

func (v paymentPrincipal) Verify(context.Context, string) (identity.Principal, error) {
	return v.Principal, nil
}

type paymentFailureIDs struct {
	values []string
	n      int
}

func (g *paymentFailureIDs) New() string { v := g.values[g.n]; g.n++; return v }

func TestPaymentHTTPDurableReplayPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	u, err := url.Parse(dsn)
	if err != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_confirmation_") {
		t.Fatal("requires explicit loopback confirmation database")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { pool.Close() }()
	ids := randomid.Generator{}
	tenant := ids.New()
	for _, sql := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'test-'||($1::uuid)::text,'Synthetic payment test','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org','store','Store','store'),($1,'other','other','Other','store')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','org','synthetic','placed','ARS',1000,1),($1,'order-2','org','synthetic','placed','ARS',1000,1),($1,'order-3','org','synthetic','placed','ARS',1000,1),($1,'private','other','synthetic','placed','ARS',1000,1)`,
	} {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("synthetic tenant=%s", tenant)
	principal := identity.Principal{Subject: "payment-operator", TenantID: tenant, Permissions: map[string]struct{}{"payment:create": {}}, Organizations: map[string]struct{}{"org": {}}}
	var handler http.Handler
	rebind := func(p identity.Principal) {
		mux := http.NewServeMux()
		CommerceModule{Service: commerce.NewService(postgres.NewCommerce(pool), ids), PaymentProvider: "stripe"}.Register(mux, paymentPrincipal{p})
		handler = mux
	}
	rebind(principal)
	body := `{"organization_id":"org","provider_code":"stripe","currency":"ARS","amount_minor_units":1000}`
	request := func(target, key, payload string, auth bool) *httptest.ResponseRecorder {
		r := httptest.NewRequest("POST", "/v1/commerce/orders/"+target+"/payments", strings.NewReader(payload))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Idempotency-Key", key)
		if auth {
			r.Header.Set("Authorization", "Bearer synthetic")
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w
	}
	key := "payment-replay-00000001"
	first := request("order", key, body, true)
	if first.Code != 202 {
		t.Fatalf("first status=%d body=%s", first.Code, first.Body.String())
	}
	var original commerce.PaymentAttempt
	if err := json.Unmarshal(first.Body.Bytes(), &original); err != nil {
		t.Fatal(err)
	}
	// Discard the first response as a client that did not retain its receipt.
	// This is not claimed as an injected transport drop.
	second := request("order", key, body, true)
	if second.Code != 202 {
		t.Fatalf("durable replay status=%d want202 body=%s", second.Code, second.Body.String())
	}
	if second.Body.String() != first.Body.String() {
		t.Fatal("replay differs from durable original")
	}
	var wg sync.WaitGroup
	responses := make(chan *httptest.ResponseRecorder, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); responses <- request("order", key, body, true) }()
	}
	wg.Wait()
	close(responses)
	for w := range responses {
		if w.Code != 202 || w.Body.String() != first.Body.String() {
			t.Fatalf("concurrent replay status=%d", w.Code)
		}
	}
	for _, c := range []struct {
		name, target, key, body string
		want                    int
	}{
		{"amount", "order", key, strings.Replace(body, ":1000", ":999", 1), 409},
		{"currency", "order", key, strings.Replace(body, "ARS", "USD", 1), 409},
		{"other-order", "order-2", key, body, 409},
		{"private-order", "private", key, body, 409},
		{"organization", "order", key, strings.Replace(body, `"org"`, `"other"`, 1), 403},
		{"invalid-key", "order", "short", body, 400},
	} {
		t.Run(c.name, func(t *testing.T) {
			if w := request(c.target, c.key, c.body, true); w.Code != c.want {
				t.Fatalf("status%d want%d", w.Code, c.want)
			}
		})
	}
	if w := request("order", key, body, false); w.Code != 401 {
		t.Fatalf("anonymous%d", w.Code)
	}
	denied := principal
	denied.Permissions = map[string]struct{}{}
	rebind(denied)
	if w := request("order", key, body, true); w.Code != 403 {
		t.Fatalf("permission%d", w.Code)
	}
	other := principal
	other.TenantID = ids.New()
	rebind(other)
	if w := request("order", key, body, true); w.Code != 409 {
		t.Fatalf("tenant%d", w.Code)
	}
	rebind(principal)
	var rows, events int
	if err := pool.QueryRow(ctx, `select count(*) from payment.payment_attempt where tenant_id=$1`, tenant).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.requested' and payload->>'actor_subject'='payment-operator'`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if rows != 1 || events != 1 {
		t.Fatalf("rows%d audited events%d", rows, events)
	}
	// A different order must be able to await a provider reference concurrently.
	responses = make(chan *httptest.ResponseRecorder, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); responses <- request("order-2", "payment-replay-00000002", body, true) }()
	}
	wg.Wait()
	close(responses)
	receipt := ""
	for w := range responses {
		if w.Code != 202 {
			t.Fatalf("fresh concurrent status%d", w.Code)
		}
		if receipt == "" {
			receipt = w.Body.String()
		}
		if receipt != w.Body.String() {
			t.Fatal("fresh concurrent IDs differ")
		}
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.requested'`, tenant).Scan(&events); err != nil || events != 2 {
		t.Fatalf("fresh race events%d err%v", events, err)
	}
	if w := request("order-2", "another-distinct-key-0002", body, true); w.Code != 409 {
		t.Fatalf("legacy route bypassed initial payment exclusivity status%d", w.Code)
	}
	// Force an outbox unique failure after the payment INSERT, not before it.
	var duplicateEvent string
	if err := pool.QueryRow(ctx, `select event_id::text from platform.outbox_event where tenant_id=$1 and aggregate_id=$2`, tenant, original.ID).Scan(&duplicateEvent); err != nil {
		t.Fatal(err)
	}
	faultMux := http.NewServeMux()
	CommerceModule{Service: commerce.NewService(postgres.NewCommerce(pool), &paymentFailureIDs{values: []string{ids.New(), duplicateEvent}}), PaymentProvider: "stripe"}.Register(faultMux, paymentPrincipal{principal})
	handler = faultMux
	if w := request("order-3", "payment-replay-00000003", body, true); w.Code != 400 {
		t.Fatalf("outbox failure status%d", w.Code)
	}
	if err := pool.QueryRow(ctx, `select count(*) from payment.payment_attempt where tenant_id=$1 and order_id='order-3'`, tenant).Scan(&rows); err != nil || rows != 0 {
		t.Fatalf("partial intent after outbox failure rows%d err%v", rows, err)
	}
	rebind(principal)
	if w := request("order-3", "payment-replay-00000003", body, true); w.Code != 202 {
		t.Fatalf("retry after rollback status%d", w.Code)
	}
	if _, err := pool.Exec(ctx, `update payment.payment_attempt set state='pending',version=2 where tenant_id=$1 and payment_attempt_id=$2`, tenant, original.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `update sales.customer_order set state='cancelled',version=2 where tenant_id=$1 and order_id='order'`, tenant); err != nil {
		t.Fatal(err)
	}
	pool.Close()
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	rebind(principal)
	w := request("order", key, body, true)
	var recovered commerce.PaymentAttempt
	if w.Code != 202 || json.Unmarshal(w.Body.Bytes(), &recovered) != nil || recovered.ID != original.ID || recovered.State != "pending" || recovered.Version != 2 {
		t.Fatalf("reconnect replay status%d", w.Code)
	}
	if _, err := pool.Exec(ctx, `update payment.payment_attempt set provider_reference='same-external-id' where tenant_id=$1`, tenant); err == nil {
		t.Fatal("duplicate non-null provider reference accepted")
	}
	t.Log("PAYMENT_HTTP_REPLAY_PASS sequential+16 concurrent replay+16 fresh race; one audited event per intent; changed intent/scope denied; independent null references; atomic outbox failure/retry; reconnect/current-state recovery")
}

func (c *commerceRepo) CreatePriceBook(context.Context, string, string, commerce.PriceBook) error {
	c.books++
	return nil
}
func (c *commerceRepo) ActivatePriceBook(context.Context, string, string, string) error {
	return commerce.ErrConflict
}
func (c *commerceRepo) PublicPrice(context.Context, string, string, string) (commerce.PriceEntry, error) {
	return commerce.PriceEntry{VariantID: "v", AmountMinorUnits: 100, TaxMode: "inclusive"}, nil
}
func (c *commerceRepo) AddOrderLine(_ context.Context, _ string, _ string, v commerce.OrderLine, _ int64) (commerce.OrderLine, error) {
	return v, nil
}
func (c *commerceRepo) PlaceOrder(context.Context, string, string, string, int64, string) error {
	return nil
}
func (c *commerceRepo) AllocateStock(context.Context, string, string, string, string, string, int64, int64, string) error {
	return nil
}
func (c *commerceRepo) CreatePaymentAttempt(context.Context, string, string, string, commerce.PaymentAttempt) error {
	return nil
}
func (c *commerceRepo) TransitionPayment(context.Context, string, string, string, string, string, int64, string, string) error {
	return nil
}

type commerceIDs struct{ n int }

func (i *commerceIDs) New() string {
	i.n++
	return []string{"018f4d4a-7b36-7a21-8d10-2f4c54c28501", "018f4d4a-7b36-7a21-8d10-2f4c54c28502"}[i.n-1]
}

type commerceVerifier struct{}

type operationVerifier struct{ permission, organization string }

func (v operationVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{v.permission: {}}, Organizations: map[string]struct{}{v.organization: {}}}, nil
}

func TestOperationReadFailsClosedWithoutScopeOrReader(t *testing.T) {
	for _, c := range []struct {
		name, permission, organization, query string
		auth                                  bool
		status                                int
	}{
		{"anonymous", "inventory:allocate", "store", "store", false, 401},
		{"customer", "customer:self", "store", "store", true, 403},
		{"other-org", "inventory:allocate", "other", "store", true, 403},
		{"missing-org", "inventory:allocate", "store", "", true, 403},
		{"reader-unavailable", "inventory:allocate", "store", "store", true, 503},
	} {
		t.Run(c.name, func(t *testing.T) {
			mux := http.NewServeMux()
			CommerceModule{Service: commerce.NewService(&commerceRepo{}, &commerceIDs{})}.Register(mux, operationVerifier{c.permission, c.organization})
			r := httptest.NewRequest("GET", "/v1/commerce/orders?organization_id="+c.query, nil)
			if c.auth {
				r.Header.Set("Authorization", "Bearer fixture")
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != c.status {
				t.Fatalf("got%d expected%d", w.Code, c.status)
			}
		})
	}
}

func (commerceVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "admin", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28500", Permissions: map[string]struct{}{"pricing:write": {}}}, nil
}
func TestCommerceHTTPPublicAndProtected(t *testing.T) {
	repo := &commerceRepo{}
	service := commerce.NewService(repo, &commerceIDs{})
	mux := http.NewServeMux()
	CommerceModule{Service: service}.Register(mux, commerceVerifier{})
	request := httptest.NewRequest("GET", "/v1/public/acme/prices?market=AR&variant_id=v", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("price status=%d", response.Code)
	}
	request = httptest.NewRequest("POST", "/v1/pricing/books", strings.NewReader(`{"market":"AR","currency":"ARS","valid_from":"`+time.Now().UTC().Format(time.RFC3339)+`","entries":[{"variant_id":"v","amount_minor_units":100,"tax_mode":"inclusive"}]}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.books != 1 {
		t.Fatalf("book status=%d books=%d body=%s", response.Code, repo.books, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/commerce/orders/o/payments", strings.NewReader(`{}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("payment permission status=%d", response.Code)
	}
}
