package domainbind

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"elite.local/enterprise/internal/agenttools"
)

type tokenProvider struct{ token string }

func (p *tokenProvider) AccessToken(context.Context) (string, error) { return p.token, nil }

func testConfig(serverURL string) Config {
	return Config{
		BaseURL:         serverURL,
		TenantID:        "tenant-id",
		TenantCode:      "acme",
		OrganizationID:  "store",
		LeadID:          "lead-1",
		ServiceKinds:    map[string]string{"corte": "service", "prueba": "test-drive"},
		ProductVariants: map[string]string{"scooter": "variant-1"},
		PriceBookID:     "retail",
		ValidMinutes:    60,
		Timeout:         5 * time.Second,
		TokenProvider:   &tokenProvider{token: "token-1"},
	}
}

func TestConfigValidateAndResolve(t *testing.T) {
	cfg := testConfig("https://backend.example")
	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	if _, err := cfg.ResolveService("desconocido"); err == nil {
		t.Fatal("unknown service accepted")
	}
	if _, err := cfg.ResolveProduct("desconocido"); err == nil {
		t.Fatal("unknown product accepted")
	}
	noLead := testConfig("https://backend.example")
	noLead.LeadID = ""
	if err := noLead.Validate(); err == nil {
		t.Fatal("missing lead accepted")
	}
	dynamic := testConfig("https://backend.example")
	dynamic.OrganizationID, dynamic.LeadID = "", ""
	if err := dynamic.Validate(); err != nil {
		t.Fatalf("dynamic-scope config rejected: %v", err)
	}
}

func TestBookAppointmentRoundTrip(t *testing.T) {
	occurredAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/public/acme/store/appointments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["kind"] != "service" || body["lead_id"] != "lead-1" {
			t.Errorf("unexpected body: %+v", body)
		}
		if r.Header.Get("Idempotency-Key") != "meta-message-0001" {
			t.Errorf("missing idempotency key: %q", r.Header.Get("Idempotency-Key"))
		}
		if body["starts_at"] != occurredAt.Add(24*time.Hour).Format(time.RFC3339) {
			t.Errorf("relative time is not retry-stable: %+v", body)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"a1","state":"requested"}`))
	}))
	defer srv.Close()

	g, err := NewGateway(testConfig(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	out, err := g.BookAppointmentCommand(context.Background(), "tenant-id", CommandIdentity{IdempotencyKey: "meta-message-0001", OccurredAt: occurredAt}, agenttools.AppointmentInput{Service: "corte", When: "manana"})
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Fatal("empty confirmation")
	}
}

func TestBookAppointmentUnknownService(t *testing.T) {
	g, _ := NewGateway(testConfig("https://backend.example"))
	if _, err := g.BookAppointmentCommand(context.Background(), "tenant-id", CommandIdentity{IdempotencyKey: "meta-message-0002", OccurredAt: time.Now()}, agenttools.AppointmentInput{Service: "spa", When: "manana"}); err == nil {
		t.Fatal("unknown service accepted")
	}
	if _, err := g.BookAppointment(context.Background(), "tenant-id", agenttools.AppointmentInput{Service: "corte", When: "manana"}); !errors.Is(err, ErrIdempotencyRequired) {
		t.Fatal("legacy appointment path did not fail closed")
	}
}

func TestCreateQuoteRoundTrip(t *testing.T) {
	provider := &tokenProvider{token: "token-1"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/franchise/quotes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["variant_id"] != "variant-1" || body["price_book_id"] != "retail" {
			t.Errorf("unexpected body: %+v", body)
		}
		if r.Header.Get("Idempotency-Key") != "meta-message-quote-1" || r.Header.Get("Authorization") != "Bearer "+provider.token {
			t.Errorf("missing command identity or current bearer")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"q1","state":"draft"}`))
	}))
	defer srv.Close()

	cfg := testConfig(srv.URL)
	cfg.TokenProvider = provider
	g, _ := NewGateway(cfg)
	command := CommandIdentity{IdempotencyKey: "meta-message-quote-1", OccurredAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)}
	if _, err := g.CreateQuoteCommand(context.Background(), "tenant-id", command, agenttools.QuoteInput{Product: "scooter", Quantity: 1}); err != nil {
		t.Fatal(err)
	}
	provider.token = "token-2"
	if _, err := g.CreateQuoteCommand(context.Background(), "tenant-id", command, agenttools.QuoteInput{Product: "scooter", Quantity: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestOrderStatusRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/customer/orders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"items":[{"order_id":"ORD-1","state":"delivered"}]}`))
	}))
	defer srv.Close()

	g, _ := NewGateway(testConfig(srv.URL))
	out, err := g.OrderStatus(context.Background(), "tenant-id", "ORD-1")
	if err != nil {
		t.Fatal(err)
	}
	if out != "Pedido ORD-1: delivered." {
		t.Fatalf("unexpected: %q", out)
	}
}

func TestRequestReturnRejectsHandover(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/customer/journey":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"handovers":[{"id":"h1","order_id":"ORD-1","version":2}]}`))
		case "/v1/customer/handovers/h1/reject":
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"state":"rejected"}`))
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	g, _ := NewGateway(testConfig(srv.URL))
	if _, err := g.RequestReturn(context.Background(), "tenant-id", agenttools.ReturnInput{OrderID: "ORD-1", Reason: "no sirve"}); err != nil {
		t.Fatal(err)
	}
}

func TestDynamicContactScopeIsNotShared(t *testing.T) {
	seen := []string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		seen = append(seen, r.URL.Path+":"+body["lead_id"].(string))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()
	g, _ := NewGateway(testConfig(srv.URL))
	command := CommandIdentity{IdempotencyKey: "provider-message-0001", OccurredAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)}
	if _, err := g.BookAppointmentFor(context.Background(), "tenant-id", Scope{OrganizationID: "north", LeadID: "lead-a"}, command, agenttools.AppointmentInput{Service: "corte", When: "manana"}); err != nil {
		t.Fatal(err)
	}
	command.IdempotencyKey = "provider-message-0002"
	if _, err := g.BookAppointmentFor(context.Background(), "tenant-id", Scope{OrganizationID: "south", LeadID: "lead-b"}, command, agenttools.AppointmentInput{Service: "corte", When: "manana"}); err != nil {
		t.Fatal(err)
	}
	if len(seen) != 2 || seen[0] != "/v1/public/acme/north/appointments:lead-a" || seen[1] != "/v1/public/acme/south/appointments:lead-b" {
		t.Fatalf("contact scopes crossed: %v", seen)
	}
	if _, err := g.BookAppointmentFor(context.Background(), "tenant-id", Scope{OrganizationID: "north"}, command, agenttools.AppointmentInput{Service: "corte", When: "manana"}); err == nil {
		t.Fatal("incomplete contact scope accepted")
	}
}

func TestProtectedAndTenantScopeFailClosed(t *testing.T) {
	cfg := testConfig("https://backend.example")
	cfg.TokenProvider = nil
	g, _ := NewGateway(cfg)
	command := CommandIdentity{IdempotencyKey: "meta-message-quote-2", OccurredAt: time.Now()}
	if _, err := g.CreateQuoteCommand(context.Background(), "tenant-id", command, agenttools.QuoteInput{Product: "scooter", Quantity: 1}); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("missing token provider not rejected: %v", err)
	}
	if _, err := g.OrderStatus(context.Background(), "other-tenant", "ORD-1"); err == nil {
		t.Fatal("cross-tenant call accepted")
	}
}
