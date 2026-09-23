package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/providerintegration"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type checkoutTransportFixture struct {
	calls        atomic.Int64
	observed     officialpayments.PaymentSnapshot
	loseResponse bool
}

func (f *checkoutTransportFixture) CreateCheckout(_ context.Context, r paymentbridge.Request) (paymentbridge.Checkout, error) {
	f.calls.Add(1)
	if f.loseResponse {
		return paymentbridge.Checkout{}, context.DeadlineExceeded
	}
	return paymentbridge.Checkout{ProviderCode: "stripe", SessionID: "cs_fixture", URL: "https://checkout.stripe.com/c/pay/cs_fixture", Status: "open", OrderID: r.OrderID, PaymentAttemptID: r.PaymentAttemptID, Currency: r.Currency, AmountMinor: r.AmountMinor, AccountRef: "acct_fixture", LiveMode: false, ExpiresAt: r.CheckoutExpiresAt}, nil
}
func (f *checkoutTransportFixture) RetrieveCheckout(context.Context, string) (paymentbridge.Checkout, error) {
	return paymentbridge.Checkout{ProviderCode: "stripe", SessionID: "cs_fixture", PaymentReference: "pi_fixture", Status: "complete"}, nil
}
func (f *checkoutTransportFixture) RetrievePayment(context.Context, string) (officialpayments.PaymentSnapshot, error) {
	return f.observed, nil
}

func TestPaymentCheckoutDurableProjection(t *testing.T) {
	raw := os.Getenv("PAYMENT_BRIDGE_DB_URL")
	if raw == "" {
		t.Skip("PAYMENT_BRIDGE_DB_URL is not set")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_payment_") {
		t.Fatal("requires a dedicated loopback elite_payment_ database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28471"
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	// Synthetic business inputs only. Payment is created through the actual owner.
	must(`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'payment-checkout','Fixture','Fixture')`, tenant)
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Fixture','store')`, tenant)
	must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','placed','ARS',1000,1)`, tenant)
	must(`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref)values($1,'checkout','stripe','FIXTURE_ONLY_NOT_A_SECRET')`, tenant)
	actual, err := NewCommerce(pool).RecordOrderPayment(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28472", "payment-key-00000001", commerce.PaymentAttempt{ID: "payment-attempt-0001", OrderID: "order", OrganizationID: "store", ProviderCode: "stripe", Currency: "ARS", AmountMinorUnits: 1000}, "operator-fixture")
	if err != nil || actual.State != "created" {
		t.Fatal(actual, err)
	}
	scope := paymentbridge.Scope{TenantID: tenant, OrganizationID: "store", ConnectionID: "checkout", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2, LiveMode: false}
	store := NewPaymentCheckoutStore(pool)
	fence, err := NewOutboundDeliveryStore(pool, []byte(strings.Repeat("f", 32)), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	transport := &checkoutTransportFixture{observed: officialpayments.PaymentSnapshot{ProviderCode: "stripe", ProviderReference: "pi_fixture", OrderID: "order", PaymentAttemptID: actual.ID, Currency: "ARS", AmountMinor: 1000, ReceivedMinor: 1000, Status: "succeeded", LiveMode: false}}
	worker := paymentbridge.Worker{Scope: scope, Store: store, Fence: &PaymentDispatchFence{OutboundDeliveryStore: fence, Scope: scope}, Driver: transport}
	t.Run("concurrent-request-one-call", func(t *testing.T) {
		var wg sync.WaitGroup
		out := make(chan error, 8)
		for range 8 {
			wg.Add(1)
			go func() { defer wg.Done(); _, e := worker.ProcessOnce(ctx, 8); out <- e }()
		}
		wg.Wait()
		close(out)
		for e := range out {
			var pg *pgconn.PgError
			if e != nil && !errors.Is(e, paymentbridge.ErrCheckoutConflict) && (!errors.As(e, &pg) || pg.Code != "40001") {
				t.Error(e)
			}
		}
		if transport.calls.Load() != 1 {
			t.Fatal("duplicate checkout effects", transport.calls.Load())
		}
		if _, e := worker.ProcessOnce(ctx, 8); e != nil {
			t.Fatal(e)
		}
		r, result, e := store.Checkout(ctx, scope, actual.ID)
		if e != nil || r.OrderID != "order" || result.SessionID != "cs_fixture" {
			t.Fatal(r, result, e)
		}
	})
	t.Run("cross-scope-request-denied", func(t *testing.T) {
		c := scope
		c.OrganizationID = "other"
		if _, _, e := store.Checkout(ctx, c, actual.ID); !errors.Is(e, paymentbridge.ErrCheckoutConflict) {
			t.Fatal(e)
		}
	})
	t.Run("capture-only-from-matched-observation", func(t *testing.T) {
		if e := worker.Reconcile(ctx, actual.ID, "pi_fixture"); e != nil {
			t.Fatal(e)
		}
		var state string
		var hold bool
		var received int64
		e := pool.QueryRow(ctx, `select p.state,o.hold,o.received_minor_units from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) where p.tenant_id=$1 and p.payment_attempt_id=$2`, tenant, actual.ID).Scan(&state, &hold, &received)
		if e != nil || state != "captured" || hold || received != 1000 {
			t.Fatal(state, hold, received, e)
		}
	})
	t.Run("partial-refund-preserves-payment-balance-owner", func(t *testing.T) {
		r, generation, e := store.BeginObservation(ctx, scope, actual.ID, "pi_fixture")
		if e != nil {
			t.Fatal(e)
		}
		snapshot := transport.observed
		snapshot.RefundedMinor = 100
		if e = store.RecordObservation(ctx, scope, r, generation, snapshot); e != nil {
			t.Fatal(e)
		}
		var state string
		var hold bool
		e = pool.QueryRow(ctx, `select p.state,o.hold from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) where p.tenant_id=$1 and p.payment_attempt_id=$2`, tenant, actual.ID).Scan(&state, &hold)
		if e != nil || state != "captured" || !hold {
			t.Fatal("partial refund became total", state, hold, e)
		}
	})
	t.Run("callback-hold-and-generation-fence", func(t *testing.T) {
		r, generation, e := store.BeginObservation(ctx, scope, actual.ID, "pi_fixture")
		if e != nil {
			t.Fatal(e)
		}
		notice := officialpayments.PaymentNotification{ProviderCode: "stripe", ProviderEventID: "evt_refund", ProviderReference: "pi_fixture", Type: "payment_intent.succeeded", BodySHA256: strings.Repeat("a", 64)}
		payload, _ := json.Marshal(paymentbridge.InboxPayload{Schema: "elite-payment-notification/v1", Notification: notice})
		receipt := providerintegration.Receipt{TenantID: tenant, ConnectionID: "checkout", ProviderCode: "stripe", ProviderEventID: notice.ProviderEventID, EventType: "payment.reconciliation-requested.v1", BodyHash: notice.BodySHA256, Payload: payload}
		replay, e := store.AcceptWebhook(ctx, receipt)
		if e != nil || replay {
			t.Fatal(replay, e)
		}
		if e = store.RecordObservation(ctx, scope, r, generation, transport.observed); !errors.Is(e, paymentbridge.ErrObservationStale) {
			t.Fatal("stale GET cleared hold", e)
		}
		var before, after int64
		var hold bool
		if e = pool.QueryRow(ctx, `select generation,hold from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2`, tenant, actual.ID).Scan(&before, &hold); e != nil || !hold {
			t.Fatal(before, hold, e)
		}
		replay, e = store.AcceptWebhook(ctx, receipt)
		if e != nil || !replay {
			t.Fatal(replay, e)
		}
		if e = pool.QueryRow(ctx, `select generation from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2`, tenant, actual.ID).Scan(&after); e != nil || before != after {
			t.Fatal(before, after, e)
		}
		var jobs int
		if e = pool.QueryRow(ctx, `select count(*) from platform.job where tenant_id=$1 and queue='payment-provider-events'`, tenant).Scan(&jobs); e != nil || jobs != 1 {
			t.Fatal(jobs, e)
		}
	})
	t.Run("mismatch-refund-dispute-never-clear-hold", func(t *testing.T) {
		for _, variant := range []string{"amount", "currency", "order", "attempt", "mode", "refund", "dispute"} {
			t.Run(variant, func(t *testing.T) {
				r, generation, e := store.BeginObservation(ctx, scope, actual.ID, "pi_fixture")
				if e != nil {
					t.Fatal(e)
				}
				snapshot := transport.observed
				switch variant {
				case "amount":
					snapshot.AmountMinor++
				case "currency":
					snapshot.Currency = "USD"
				case "order":
					snapshot.OrderID = "other"
				case "attempt":
					snapshot.PaymentAttemptID = "other"
				case "mode":
					snapshot.LiveMode = true
				case "refund":
					snapshot.RefundedMinor = 1000
				case "dispute":
					snapshot.Disputed = true
					snapshot.RefundedMinor = 1000
				}
				e = store.RecordObservation(ctx, scope, r, generation, snapshot)
				if variant != "refund" && variant != "dispute" && !errors.Is(e, paymentbridge.ErrPaymentMismatch) {
					t.Fatal(e)
				}
				if (variant == "refund" || variant == "dispute") && e != nil {
					t.Fatal(e)
				}
				var hold bool
				if e = pool.QueryRow(ctx, `select hold from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2`, tenant, actual.ID).Scan(&hold); e != nil || !hold {
					t.Fatal(hold, e)
				}
			})
		}
	})
	t.Run("terminal-state-never-regresses-to-captured", func(t *testing.T) {
		var before string
		if e := pool.QueryRow(ctx, `select state from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2`, tenant, actual.ID).Scan(&before); e != nil || (before != "refunded" && before != "disputed") {
			t.Fatal(before, e)
		}
		if e := worker.Reconcile(ctx, actual.ID, "pi_fixture"); e != nil && !errors.Is(e, paymentbridge.ErrPaymentMismatch) {
			t.Fatal(e)
		}
		var state string
		var hold bool
		e := pool.QueryRow(ctx, `select p.state,o.hold from payment.payment_attempt p join payment.provider_observation o using(tenant_id,payment_attempt_id) where p.tenant_id=$1 and p.payment_attempt_id=$2`, tenant, actual.ID).Scan(&state, &hold)
		if e != nil || state != before || !hold {
			t.Fatal("terminal financial state regressed", state, hold, e)
		}
	})
	t.Run("connection-organization-is-authoritative", func(t *testing.T) {
		must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'other-store','other-store','Other fixture','store')`, tenant)
		must(`update integration.provider_connection set organization_id='other-store' where tenant_id=$1 and connection_id='checkout'`, tenant)
		defer must(`update integration.provider_connection set organization_id=null where tenant_id=$1 and connection_id='checkout'`, tenant)
		if _, _, e := store.Checkout(ctx, scope, actual.ID); !errors.Is(e, providerintegration.ErrConnection) {
			t.Fatal("connection for a different organization was accepted", e)
		}
		if e := worker.Reconcile(ctx, actual.ID, "pi_fixture"); !errors.Is(e, providerintegration.ErrConnection) {
			t.Fatal("cross-organization observation was accepted", e)
		}
	})
	t.Run("cancel-after-bind-rolls-back-dispatch", func(t *testing.T) {
		// Persist an old, never-dispatched request as fixture input. It must not
		// starve new work and must not be silently retried or marked financially paid.
		must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order-old-request','store','customer','placed','ARS',1000,1)`, tenant)
		must(`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,idempotency_key,state,currency,amount_minor_units,version)values($1,'payment-old-request','order-old-request','stripe','old-request-fixture','created','ARS',1000,1)`, tenant)
		must(`insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,gen_random_uuid(),'payment','payment-old-request',1,'payment.requested',1,clock_timestamp()-interval '24 hours','{}')`, tenant)
		must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order-cancelled-fixture','store','customer','placed','ARS',1000,1)`, tenant)
		attempt, e := NewCommerce(pool).RecordOrderPayment(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28473", "payment-key-cancelled", commerce.PaymentAttempt{ID: "payment-cancelled-fixture", OrderID: "order-cancelled-fixture", OrganizationID: "store", ProviderCode: "stripe", Currency: "ARS", AmountMinorUnits: 1000}, "operator-fixture")
		if e != nil {
			t.Fatal(e)
		}
		requests, e := store.PendingRequests(ctx, scope, 10)
		if e != nil || len(requests) != 1 || requests[0].PaymentAttemptID != attempt.ID {
			t.Fatal(requests, e)
		}
		r := requests[0]
		body, e := json.Marshal(r)
		if e != nil {
			t.Fatal(e)
		}
		message := channels.Message{ChannelCode: "payment_stripe", TenantID: tenant, ExternalID: r.CustomerSubject, ThreadID: r.OrderID, DeliveryKey: r.PaymentAttemptID, Direction: channels.DirectionOut, Text: string(body)}
		hash, e := outbounddelivery.MessageSHA256(message)
		if e != nil {
			t.Fatal(e)
		}
		if e = store.BindRequest(ctx, scope, r, hash); e != nil {
			t.Fatal(e)
		}
		// The operator's cancellation commits after preparation but before the send claim.
		must(`update payment.payment_attempt set state='failed',version=version+1 where tenant_id=$1 and payment_attempt_id=$2`, tenant, attempt.ID)
		if _, e = worker.Fence.Claim(ctx, message, hash); !errors.Is(e, paymentbridge.ErrCheckoutConflict) {
			t.Fatal("stale dispatch was accepted", e)
		}
		var state string
		var fences, events int
		e = pool.QueryRow(ctx, `select state,(select count(*) from communication.outbound_delivery f where f.tenant_id=p.tenant_id and f.delivery_key=p.payment_attempt_id),(select count(*) from platform.outbox_event e where e.tenant_id=p.tenant_id and e.aggregate_id=p.payment_attempt_id and e.event_type='payment.checkout-requested') from payment.payment_attempt p where tenant_id=$1 and payment_attempt_id=$2`, tenant, attempt.ID).Scan(&state, &fences, &events)
		if e != nil || state != "failed" || fences != 0 || events != 0 {
			t.Fatal("rejected dispatch left side effects", state, fences, events, e)
		}
	})
	t.Run("disabled-connection-no-observation", func(t *testing.T) {
		must(`update integration.provider_connection set state='disabled' where tenant_id=$1 and connection_id='checkout'`, tenant)
		if e := worker.Reconcile(ctx, actual.ID, "pi_fixture"); !errors.Is(e, providerintegration.ErrConnection) {
			t.Fatal(e)
		}
	})
	// This database is discarded by the runner. Immutable audit rows are never
	// disabled or deleted merely to make a test's cleanup convenient.
}
