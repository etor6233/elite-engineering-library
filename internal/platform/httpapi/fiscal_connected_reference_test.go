package httpapi

// AUTHORED local composition fixture. Paid-source rows and provider responses are
// synthetic; actual HTTP handlers, PostgreSQL, Go worker, UDS and .NET owners run.
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/fiscal/wsfeipc"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type fiscalReferenceReceipt struct {
	Tenant  string         `json:"tenant"`
	Invoice string         `json:"invoice"`
	Credit  string         `json:"credit"`
	Stats   map[string]any `json:"stats"`
}

func TestFiscalConnectedReference(t *testing.T) {
	dsn, socket, path := os.Getenv("TEST_DATABASE_URL"), os.Getenv("ARCA_FIXTURE_SOCKET"), os.Getenv("ARCA_FIXTURE_RECEIPT")
	if dsn == "" || socket == "" || path == "" {
		t.Skip("explicit local fixture environment required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	ids := randomid.Generator{}
	tenant := ids.New()
	for _, sql := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'arca-'||$1::text,'Local ARCA Fixture','Local ARCA Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise','Fixture','franchisee')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'paid-order','franchise','fixture-customer','paid','ARS',12100,2)`,
		`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values($1,'captured-payment','paid-order','fixture','fixture-provider-reference','fixture-payment-key','captured','ARS',12100,3)`,
	} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := postgres.NewFiscal(pool)
	service := fiscal.NewService(repo, ids)
	principal := identity.Principal{Subject: "fixture-controller", TenantID: tenant, Permissions: map[string]struct{}{"fiscal:configure": {}, "fiscal:issue": {}, "fiscal:read": {}}, Organizations: map[string]struct{}{"franchise": {}}}
	mux := http.NewServeMux()
	FiscalModule{Service: service}.Register(mux, fiscalVerifier{p: principal})
	server := httptest.NewServer(mux)
	defer server.Close()
	call := func(method, url, key string, body []byte, want int) []byte {
		t.Helper()
		request, er := http.NewRequestWithContext(ctx, method, server.URL+url, bytes.NewReader(body))
		if er != nil {
			t.Fatal(er)
		}
		request.Header.Set("Authorization", "Bearer fixture-principal")
		request.Header.Set("Content-Type", "application/json")
		if key != "" {
			request.Header.Set("Idempotency-Key", key)
		}
		response, er := server.Client().Do(request)
		if er != nil {
			t.Fatal(er)
		}
		defer response.Body.Close()
		raw, er := io.ReadAll(io.LimitReader(response.Body, 65537))
		if er != nil || response.StatusCode != want {
			t.Fatalf("fiscal HTTP status=%d expected=%d error=%v body=%s", response.StatusCode, want, er, raw)
		}
		return raw
	}
	posRaw := call("POST", "/v1/fiscal/points-of-sale", "", []byte(`{"organization_id":"franchise","taxpayer_cuit":"30715117564","environment":"homologation","number":12}`), 201)
	var pos fiscal.PointOfSale
	if err = json.Unmarshal(posRaw, &pos); err != nil || pos.ID == "" {
		t.Fatalf("point-of-sale %v", err)
	}
	invoice := fiscal.Invoice{OrganizationID: "franchise", OrderID: "paid-order", PaymentAttemptID: "captured-payment", PointOfSaleID: pos.ID, VoucherType: 6, Concept: 1, RecipientDocumentType: 99, RecipientDocument: "0", RecipientVATConditionID: 5, NetMinorUnits: 10000, VATMinorUnits: 2100, VATLines: []fiscal.VATLine{{ID: 5, BaseMinorUnits: 10000, AmountMinorUnits: 2100}}, IssuedOn: time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)}
	body, _ := json.Marshal(invoice)
	raw := call("POST", "/v1/fiscal/invoices", "fixture-issue-once", body, 202)
	var queued fiscal.Invoice
	if err = json.Unmarshal(raw, &queued); err != nil {
		t.Fatal(err)
	}
	replay := call("POST", "/v1/fiscal/invoices", "fixture-issue-once", body, 200)
	var same fiscal.Invoice
	if err = json.Unmarshal(replay, &same); err != nil || same.ID != queued.ID {
		t.Fatal("invoice replay changed identity")
	}
	call("POST", "/v1/fiscal/invoices", "fixture-denied-org", bytes.Replace(body, []byte(`"franchise"`), []byte(`"other"`), 1), 403)
	provider, err := wsfeipc.NewUnixProvider(socket, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	processor, err := fiscal.NewProcessor(repo, provider, ids, "fixture-worker-a", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = processor.ProcessOne(ctx); err == nil {
		t.Fatal("lost provider response unexpectedly accepted")
	}
	pending, err := repo.GetInvoice(ctx, tenant, "franchise", queued.ID)
	if err != nil || pending.Status != "reconcile_required" || pending.VoucherNumber != 1 {
		t.Fatalf("ambiguous request was not durable: %+v %v", pending, err)
	}
	resumed, err := fiscal.NewProcessor(postgres.NewFiscal(pool), provider, ids, "fixture-worker-restarted", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	authorized, err := resumed.ProcessOne(ctx)
	if err != nil || authorized.ID != queued.ID || authorized.Status != "authorized" || authorized.CAE != "41124578989845" {
		t.Fatalf("consult/reconcile failed: %+v %v", authorized, err)
	}
	invoice.VoucherType = 8
	invoice.AssociatedVouchers = []fiscal.AssociatedVoucher{{InvoiceID: queued.ID}}
	invoice.IssuedOn = invoice.IssuedOn.AddDate(0, 0, 1)
	creditBody, _ := json.Marshal(invoice)
	raw = call("POST", "/v1/fiscal/invoices", "fixture-credit-once", creditBody, 202)
	var credit fiscal.Invoice
	if err = json.Unmarshal(raw, &credit); err != nil {
		t.Fatal(err)
	}
	credit, err = resumed.ProcessOne(ctx)
	if err != nil || credit.Status != "authorized" || len(credit.AssociatedVouchers) != 1 || credit.AssociatedVouchers[0].Number != 1 || credit.AssociatedVouchers[0].VoucherType != 6 {
		t.Fatalf("credit association did not cross actual bridge: %+v %v", credit, err)
	}
	if _, err = resumed.ProcessOne(ctx); !errors.Is(err, fiscal.ErrNoWork) {
		t.Fatalf("duplicate work left: %v", err)
	}
	for _, id := range []string{queued.ID, credit.ID} {
		raw = call("GET", "/v1/fiscal/invoices/"+id+"?organization_id=franchise", "", nil, 200)
		if !bytes.Contains(raw, []byte(`"status":"authorized"`)) {
			t.Fatal("authorized invoice not readable")
		}
	}
	for _, kind := range []string{"voucher_type", "concept", "document_type", "vat_rate", "other_tax", "point_of_sale", "recipient_vat_condition"} {
		var class *string
		if kind == "recipient_vat_condition" {
			v := "A"
			class = &v
		}
		snapshot, er := provider.FetchParameters(ctx, "30715117564", kind, class)
		if er != nil || len(snapshot.Items) != 1 || len(snapshot.ResponseHash) != 64 {
			t.Fatalf("parameter %s: %v", kind, er)
		}
	}
	dialer := net.Dialer{Timeout: time.Second}
	client := http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{DialContext: func(c context.Context, _, _ string) (net.Conn, error) { return dialer.DialContext(c, "unix", socket) }}}
	response, err := client.Get("http://arca-wsfe/fixture/stats")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var stats map[string]any
	if err = json.NewDecoder(response.Body).Decode(&stats); err != nil {
		t.Fatal(err)
	}
	if stats["authorize"] != float64(2) || stats["last"] != float64(2) || stats["consult"] != float64(1) || stats["cms"] != float64(1) {
		t.Fatalf("unexpected provider side effects: %+v", stats)
	}
	var rows, events int
	if err = pool.QueryRow(ctx, `select count(*) from fiscal.invoice where tenant_id=$1 and status='authorized'`, tenant).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='fiscal-invoice.authorized'`, tenant).Scan(&events); err != nil || rows != 2 || events != 2 {
		t.Fatalf("durable effects rows=%d events=%d %v", rows, events, err)
	}
	receipt, _ := json.Marshal(fiscalReferenceReceipt{Tenant: tenant, Invoice: queued.ID, Credit: credit.ID, Stats: stats})
	if err = os.WriteFile(path, receipt, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("ARCA_CONNECTED_HTTP_PG_GO_UDS_DOTNET_CMS_PASS invoices=2 authorized_events=2 authorize_calls=2 lost_response_reconciled=1 credential_refreshes=1 parameter_kinds=7")
	stopped, stopErr := client.Post("http://arca-wsfe/fixture/stop", "application/json", strings.NewReader("{}"))
	if stopErr != nil {
		t.Fatal(stopErr)
	}
	defer stopped.Body.Close()
	if stopped.StatusCode != 200 {
		t.Fatal("fixture graceful stop failed")
	}
}

func TestFiscalReferenceRecovery(t *testing.T) {
	path, dsn := os.Getenv("ARCA_FIXTURE_RECEIPT"), os.Getenv("TEST_DATABASE_URL")
	if path == "" || dsn == "" {
		t.Skip("explicit recovery fixture required")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var r fiscalReferenceReceipt
	if err = json.Unmarshal(raw, &r); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	repo := postgres.NewFiscal(pool)
	for _, id := range []string{r.Invoice, r.Credit} {
		v, er := repo.GetInvoice(ctx, r.Tenant, "franchise", id)
		if er != nil || v.Status != "authorized" || v.CAE != "41124578989845" || !strings.HasPrefix(v.Environment, "homo") {
			t.Fatalf("durable fiscal recovery %+v %v", v, er)
		}
	}
	t.Log("ARCA_POSTGRES_RESTART_PASS invoices=2 provider_calls=0")
}
