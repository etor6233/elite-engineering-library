package postgres

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/paymentbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPaymentCustomerCheckoutAuthorizationPostgres(t *testing.T) {
	raw := os.Getenv("PAYMENT_BRIDGE_DB_URL")
	if raw == "" {
		t.Skip("PAYMENT_BRIDGE_DB_URL is not set")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_payment_") {
		t.Fatal("requires dedicated loopback elite_payment_ database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28771"
	must := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	must(`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'payment-customer-read','Fixture','Fixture')`, tenant)
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Fixture','store')`, tenant)
	must(`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name)values($1,'alice','Fixture')`, tenant)
	must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','alice','placed','ARS',1000,1)`, tenant)
	must(`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref)values($1,'checkout','stripe','FIXTURE_ONLY_NOT_A_SECRET')`, tenant)
	_, err = NewCommerce(pool).RecordOrderPayment(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28772", "payment-read-key-00001", commerce.PaymentAttempt{ID: "read-payment-attempt", OrderID: "order", OrganizationID: "store", ProviderCode: "stripe", Currency: "ARS", AmountMinorUnits: 1000}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	scope := paymentbridge.Scope{TenantID: tenant, OrganizationID: "store", ConnectionID: "checkout", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}
	fence, err := NewOutboundDeliveryStore(pool, []byte(strings.Repeat("f", 32)), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	worker := paymentbridge.Worker{Scope: scope, Store: NewPaymentCheckoutStore(pool), Fence: &PaymentDispatchFence{OutboundDeliveryStore: fence, Scope: scope}, Driver: &checkoutTransportFixture{}}
	if n, err := worker.ProcessOnce(ctx, 1); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	reader, err := NewCustomerCheckoutReader(pool, scope)
	if err != nil {
		t.Fatal(err)
	}
	if value, err := reader.CustomerCheckout(ctx, tenant, "store", "alice", "order"); err != nil || !value.Valid(time.Now()) {
		t.Fatal(value, err)
	}
	for _, args := range [][4]string{{"018f4d4a-7b36-7a21-8d10-2f4c54c28773", "store", "alice", "order"}, {tenant, "other", "alice", "order"}, {tenant, "store", "bob", "order"}, {tenant, "store", "alice", "other"}} {
		if _, err := reader.CustomerCheckout(ctx, args[0], args[1], args[2], args[3]); !errors.Is(err, paymentbridge.ErrCheckoutNotAvailable) {
			t.Fatal("foreign read allowed", err)
		}
	}
	must(`update crm.customer_profile set status='restricted' where tenant_id=$1`, tenant)
	if _, err := reader.CustomerCheckout(ctx, tenant, "store", "alice", "order"); !errors.Is(err, paymentbridge.ErrCheckoutNotAvailable) {
		t.Fatal("restricted customer allowed")
	}
	must(`update crm.customer_profile set status='active' where tenant_id=$1`, tenant)
	must(`update payment.provider_checkout set checkout_url='https://attacker.example/session' where tenant_id=$1`, tenant)
	if _, err := reader.CustomerCheckout(ctx, tenant, "store", "alice", "order"); !errors.Is(err, paymentbridge.ErrCheckoutNotAvailable) {
		t.Fatal("unsafe stored URL allowed")
	}
}
