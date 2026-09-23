package postgres_test

// AUTHORED integration fixture: source-derived approved stored value changes
// only provider due, then the real SDK/callback/receipt journey completes.
import (
	"context"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
	"testing"
)

func fixtureStoredValueSetup(t *testing.T, pool *pgxpool.Pool, tenant, kind string, giftValue ...int64) (*db.StoredValue, identity.Principal, string) {
	t.Helper()
	ctx := context.Background()
	root, e := filepath.Abs("../../..")
	if e != nil {
		t.Fatal(e)
	}
	read := func(rel string) []byte {
		b, e := os.ReadFile(filepath.Join(root, rel))
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	var c map[string]any
	if json.Unmarshal(read("deploy/stored-value/profile.reference.json"), &c) != nil {
		t.Fatal("profile JSON")
	}
	c["tenant_id"] = tenant
	c["organization_id"] = "store"
	raw, e := json.Marshal(c)
	if e != nil {
		t.Fatal(e)
	}
	profile, e := sv.Load(raw, sv.Hash(raw))
	if e != nil {
		t.Fatal(e)
	}
	python := os.Getenv("HANDOVER_PROFILE_PYTHON")
	if python == "" {
		t.Fatal("explicit fixed Python required")
	}
	process := sv.Process{Python: python, Script: filepath.Join(root, "odoo_loyalty/run.py"), ScriptSHA256: sv.Hash(read("odoo_loyalty/run.py")), Manifest: filepath.Join(root, "odoo_loyalty/engine-lock.json"), ManifestSHA256: sv.Hash(read("odoo_loyalty/engine-lock.json"))}
	store, e := db.NewStoredValue(pool, profile, process)
	if e != nil {
		t.Fatal(e)
	}
	actor := identity.Principal{TenantID: tenant, Subject: "operator", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"stored_value:read": {}, "stored_value:request": {}, "stored_value:approve": {}, "stored_value:fund": {}}}
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	must(`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'gift-50','model','gift-50','Synthetic gift','{}','active')`, tenant)
	product, price := "gift-50", int64(5000)
	if len(giftValue) == 1 {
		price = giftValue[0]
	}
	if kind == "loyalty" {
		product, price = "variant", 20000
	}
	// The upstream issues/accrues at confirmed sale. This source-order fixture
	// is explicit; the target purchase uses actual quote/accept/stock writers.
	must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'source','store','customer','confirmed','ARS',$2,1)`, tenant, price)
	must(`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units)values($1,'source','source-line',$2,1,$3)`, tenant, product, price)
	return store, actor, "reference-" + kind
}

func fixtureStoredValue(t *testing.T, pool *pgxpool.Pool, tenant, kind string, giftValue ...int64) (*db.StoredValue, identity.Principal, string) {
	t.Helper()
	store, actor, program := fixtureStoredValueSetup(t, pool, tenant, kind, giftValue...)
	ctx := context.Background()
	review := actor
	review.Subject = "reviewer"
	op := "issue"
	if kind == "loyalty" {
		op = "accrue"
	}
	req := sv.Request{OperationID: "source-operation", Operation: op, OrganizationID: "store", OrderID: "source", ProgramID: program, ExpectedOrderVersion: 1, ProfileSHA256: store.ProfileSHA256()}
	v := storedValueHTTPPropose(t, store, actor, req)
	v = storedValueHTTPApprove(t, store, review, v)
	account := v.Payload.Entries[0].AccountID
	var version int64
	if e := pool.QueryRow(ctx, `select version from sales.customer_order where tenant_id=$1 and order_id='order'`, tenant).Scan(&version); e != nil {
		t.Fatal(e)
	}
	req = sv.Request{OperationID: "redemption", Operation: "redeem", OrganizationID: "store", OrderID: "order", ProgramID: program, ExpectedOrderVersion: version, AccountID: account, ProfileSHA256: store.ProfileSHA256()}
	v = storedValueHTTPPropose(t, store, actor, req)
	v = storedValueHTTPApprove(t, store, review, v)
	if v.State != "approved" {
		t.Fatal(v)
	}
	return store, actor, account
}

func TestStoredValueOfficialSDKTenderAndRelease(t *testing.T) {
	for _, kind := range []string{"gift_card", "loyalty"} {
		t.Run(kind, func(t *testing.T) {
			ctx := context.Background()
			var store *db.StoredValue
			var actor identity.Principal
			var account string
			var allocation sv.Allocation
			r := newConnectedRunWithAllocation(t, kind == "gift_card", func(t *testing.T, pool *pgxpool.Pool, tenant string) int64 {
				store, actor, account = fixtureStoredValue(t, pool, tenant, kind)
				var e error
				allocation, e = store.Allocation(ctx, actor, "order")
				if e != nil {
					t.Fatal(e)
				}
				expected := int64(118456)
				if kind == "loyalty" {
					expected = 117283
				}
				if allocation.GrossMinor != 123456 || allocation.ProviderDueMinor != expected {
					t.Fatal("gross/due", allocation)
				}
				// Legacy caller-supplied gross may not bypass the derived remaining due.
				_, e = db.NewCommerce(pool).RecordPaymentIntent(ctx, tenant, randomid.Generator{}.New(), "forged-gross-00001", commerce.PaymentAttempt{ID: "forged-gross", OrderID: "order", OrganizationID: "store", ProviderCode: "stripe", Currency: "ARS", AmountMinorUnits: 123456}, "operator")
				if e == nil {
					t.Fatal("gross charged after tender")
				}
				return expected
			})
			if code := r.callback(t, "evt_tender", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
				t.Fatal(code)
			}
			if n, e := r.processor.ProcessOnce(ctx); e != nil || n != 1 {
				t.Fatal("reconciliation", n, e)
			}
			var amount, gross int64
			var state, hash string
			e := r.pool.QueryRow(ctx, `select p.amount_minor_units,o.total_minor_units,p.state,ob.evidence_sha256_hex from payment.payment_attempt p join sales.customer_order o using(tenant_id,order_id) join payment.provider_observation ob using(tenant_id,payment_attempt_id) where p.tenant_id=$1`, r.tenant).Scan(&amount, &gross, &state, &hash)
			if e != nil || amount != allocation.ProviderDueMinor || gross != 123456 || state != "captured" {
				t.Fatal("captured tender", amount, gross, state, e)
			}
			// Once provider intent exists, another allocation cannot change its amount.
			_, e = store.Propose(ctx, actor, sv.Request{OperationID: "after-payment", Operation: "redeem", OrganizationID: "store", OrderID: "order", ProgramID: "reference-" + kind, ExpectedOrderVersion: allocation.OrderVersion, AccountID: account, ProfileSHA256: store.ProfileSHA256()})
			if e == nil {
				t.Fatal("allocation changed after dispatch")
			}
			assertConnectedCommercialRelease(t, r, hash, true)
			t.Logf("gross=%d provider=%d gift=%d discount=%d; SDK callback, recovery, handover/commercial receipt and partial-refund invalidation passed", gross, amount, allocation.GiftMinor, allocation.DiscountMinor)
		})
	}
}

func TestStoredValueHandoverRejectsProviderOnlyProfile(t *testing.T) {
	ctx := context.Background()
	r := newConnectedRunWithAllocation(t, false, func(t *testing.T, pool *pgxpool.Pool, tenant string) int64 {
		store, actor, _ := fixtureStoredValue(t, pool, tenant, "gift_card")
		a, e := store.Allocation(ctx, actor, "order")
		if e != nil {
			t.Fatal(e)
		}
		return a.ProviderDueMinor
	})
	if code := r.callback(t, "evt_tender_policy", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, e := r.processor.ProcessOnce(ctx); e != nil || n != 1 {
		t.Fatal(n, e)
	}
	var hash string
	if e := r.pool.QueryRow(ctx, `select evidence_sha256_hex from payment.provider_observation where tenant_id=$1`, r.tenant).Scan(&hash); e != nil {
		t.Fatal(e)
	}
	policy, _ := connectedCommercialProfile(t, r, false)
	svc, e := franchisejourney.NewHandoverPreparationService(db.NewFranchiseJourney(r.pool), randomid.Generator{}, policy)
	if e != nil {
		t.Fatal(e)
	}
	_, _, e = svc.Prepare(ctx, r.tenant, "operator", franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: r.payment.ID, ObservationSHA256: hash, IdempotencyKey: "provider-only-must-reject"})
	if e == nil {
		t.Fatal("provider-only profile accepted stored-value funding without explicit algorithm selection")
	}
	var n int
	if e = r.pool.QueryRow(ctx, `select count(*) from sales.delivery_handover_preparation where tenant_id=$1`, r.tenant).Scan(&n); e != nil || n != 0 {
		t.Fatal("unselected policy left an effect", n, e)
	}
}
