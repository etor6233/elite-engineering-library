package postgres_test

// AUTHORED read-model proof. Synthetic source rows are explicit; stored value and
// survey deltas below use existing domain writers, never caller-supplied totals.
import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	rm "elite.local/enterprise/internal/rolemetrics"
	sc "elite.local/enterprise/internal/serialsupply"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func roleMetricPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("owned database required")
	}
	p, e := pgxpool.New(context.Background(), raw)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(p.Close)
	return p
}
func roleMetricSeed(t *testing.T, p *pgxpool.Pool) string {
	t.Helper()
	tenant := uuid.NewString()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'query-api','Query','Query')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store','store','A','store'),($1,'other','other','B','store')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized,status) values($1,'customer','Customer 1','customer1@example.test','active'),($1,'customer-2','Customer 2','customer2@example.test','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status) values($1,'supplier','supplier','Supplier','active')`,
		`insert into procurement.purchase_order(tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version) values($1,'po','supplier','store','accepted','USD',100,1)`,
		`insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,state) values($1,'unit','po','variant','UNIT-QUERY','planned')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order-a','store','customer','draft','USD',100,1),($1,'order-b','store','customer-2','draft','USD',200,1),($1,'order-c','other','customer','draft','USD',300,1)`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'stock-a','store','variant','STOCK-QUERY-A','VIN-QUERY-A','BATTERY-QUERY-A','available',1,clock_timestamp()),($1,'stock-b','other','variant','STOCK-QUERY-B','VIN-QUERY-B','BATTERY-QUERY-B','available',1,clock_timestamp())`,
		`insert into service_ops.warranty(tenant_id,warranty_id,stock_unit_id,customer_principal_id,starts_at,ends_at,terms_version,status) values($1,'warranty','stock-a','customer',clock_timestamp(),clock_timestamp()+interval '1 year','v1','active')`,
		`insert into service_ops.service_case(tenant_id,service_case_id,stock_unit_id,organization_id,state,severity,description,version) values($1,'case-a','stock-a','store','opened','medium','inspection',1)`,
	} {
		if _, e := p.Exec(context.Background(), strings.ReplaceAll(q, "'query-api'", "'metric-"+tenant+"'"), tenant); e != nil {
			t.Fatal(e)
		}
	}
	return tenant
}
func TestRoleMetricsConnectedScopeAndPrecision(t *testing.T) {
	ctx := context.Background()
	pool := roleMetricPool(t)
	tenant := roleMetricSeed(t, pool)
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	store, e := db.NewRoleMetrics(pool)
	if e != nil {
		t.Fatal(e)
	}
	p := identity.Principal{TenantID: tenant, Subject: "customer", Permissions: map[string]struct{}{"*": {}}, Organizations: map[string]struct{}{"store": {}}}
	kinds := []string{"orders", "own-orders", "leads", "stock", "cases", "own-cases", "shipments", "appointments", "own-appointments", "factory-destination", "factory-owned", "supply"}
	for _, kind := range kinds {
		v, e := store.Read(ctx, p, kind, "store")
		if e != nil || v.ObservedAt.IsZero() || v.Rows == nil {
			t.Fatal(kind, v, e)
		}
	}
	must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'large-a','store','customer','confirmed','ARS',9007199254740993,1),($1,'large-b','store','customer-2','confirmed','ARS',9007199254740993,1)`, tenant)
	v, e := store.Read(ctx, p, "orders", "store")
	if e != nil || len(v.Rows) != 2 || v.Rows[0].TotalMinor != "18014398509481986" || v.Rows[0].Currency != "ARS" || v.Rows[1].TotalMinor != "300" || v.Rows[1].Currency != "USD" {
		t.Fatal(v, e)
	}
	own := p
	own.Permissions = map[string]struct{}{"customer:self": {}}
	v, e = store.Read(ctx, own, "own-orders", "store")
	if e != nil || v.Rows[0].TotalMinor != "9007199254740993" || v.Rows[1].TotalMinor != "100" {
		t.Fatal("own isolation", v, e)
	}
	for _, kind := range []string{"orders", "stock", "cases", "leads"} {
		if _, e = store.Read(ctx, own, kind, "store"); e == nil {
			t.Fatal("permission widened", kind)
		}
	}
	own.Subject = "customer-2"
	v, e = store.Read(ctx, own, "own-cases", "store")
	if e != nil || len(v.Rows) != 0 {
		t.Fatal("foreign warranty", v, e)
	}
	foreign := p
	foreign.TenantID = uuid.NewString()
	v, e = store.Read(ctx, foreign, "orders", "store")
	if e != nil || len(v.Rows) != 0 {
		t.Fatal("foreign tenant", v, e)
	}
	must(`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','converted','fixture','{}')`, tenant)
	must(`insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,appointment_kind,starts_at,state,version)values($1,'appt','store','lead','customer','consultation',clock_timestamp()+interval '1 day','requested',1)`, tenant)
	v, e = store.Read(ctx, p, "appointments", "store")
	if e != nil || len(v.Rows) != 1 || v.Rows[0].Count != "1" {
		t.Fatal(v, e)
	}
	v, e = store.Read(ctx, own, "own-appointments", "store")
	if e != nil || len(v.Rows) != 0 {
		t.Fatal("customer appointment leak", v, e)
	}
	must(`insert into logistics.shipment(tenant_id,shipment_id,provider_code,origin_organization_id,destination_organization_id,state)values($1,'shipment','fixture','other','store','planned')`, tenant)
	v, e = store.Read(ctx, p, "shipments", "store")
	if e != nil || len(v.Rows) != 1 || v.Rows[0].Count != "1" {
		t.Fatal(v, e)
	}
	verifier, token := catalogRoleIssuer(t)
	mux := http.NewServeMux()
	httpapi.EnterpriseQueryModule{Service: enterprisequery.NewService(db.NewEnterpriseQuery(pool)), Metrics: store}.Register(mux, verifier)
	get := func(path, subject, tid string, permissions, orgs []string, want int) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest("GET", path, nil)
		r.Header.Set("Authorization", "Bearer "+token(subject, tid, permissions, orgs))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != want || !strings.Contains(w.Header().Get("Cache-Control"), "no-store") {
			t.Fatal(path, w.Code, w.Body.String())
		}
		return w
	}
	get("/v1/reporting/operations/orders?organization_id=store", "customer", tenant, []string{"customer:self"}, []string{"store"}, 403)
	get("/v1/reporting/operations/orders?organization_id=store", "admin", tenant, []string{"admin:read"}, []string{"other"}, 403)
	for _, q := range []string{"organization_id=store&organization_id=other", "organization_id=store&tenant_id=x", "organization_id=store;x=y"} {
		get("/v1/reporting/operations/orders?"+q, "admin", tenant, []string{"admin:read"}, []string{"store"}, 400)
	}
	get("/v1/reporting/operations/orders?organization_id=store", "admin", tenant, []string{"admin:read"}, []string{"store"}, 200)
	for n := 0; n < 101; n++ {
		must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,$2,'other','customer','draft',$3,1,1)`, tenant, fmt.Sprint("cardinality-", n), fmt.Sprintf("X%c%c", 'A'+n/26, 'A'+n%26))
	}
	if _, e = store.Read(ctx, p, "orders", "other"); !errors.Is(e, rm.ErrCardinality) {
		t.Fatal("silent truncation", e)
	}
	get("/v1/reporting/operations/orders?organization_id=other", "admin", tenant, []string{"admin:read"}, []string{"other"}, 409)
	pool.Close()
	get("/v1/reporting/operations/orders?organization_id=store", "admin", tenant, []string{"admin:read"}, []string{"store"}, 503)
	t.Log("ROLE_METRICS_SCOPE_PASS kinds=12 exact_sum_above_JS_safe_integer=true tenant_org_subject_permissions=true currency_state_separate=true cardinality_unavailable_not_zero=true")
}
func TestRoleMetricsStoredValueAndSurveyOwners(t *testing.T) {
	ctx := context.Background()
	pool := roleMetricPool(t)
	tenant := roleMetricSeed(t, pool)
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	store, p, program := fixtureStoredValueSetup(t, pool, tenant, "gift_card")
	read := func() db.StoredValueMetrics {
		t.Helper()
		v, e := store.Metrics(ctx, p, "store", program)
		if e != nil {
			t.Fatal(e)
		}
		return v
	}
	if len(read().Rows) != 0 {
		t.Fatal("empty ledger")
	}
	request := sv.Request{OperationID: "metric-issue", Operation: "issue", OrganizationID: "store", OrderID: "source", ProgramID: program, ExpectedOrderVersion: 1, ProfileSHA256: store.ProfileSHA256()}
	proposal := storedValueHTTPPropose(t, store, p, request)
	if len(read().Rows) != 0 {
		t.Fatal("pending counted as committed")
	}
	review := p
	review.Subject = "reviewer"
	storedValueHTTPApprove(t, store, review, proposal)
	v := read()
	if len(v.Rows) != 1 || v.Rows[0].Operation != "issue" || v.Rows[0].Operations != "1" || v.Rows[0].PointsDelta != "50.000000" {
		t.Fatal(v)
	}
	storedValueHTTPPropose(t, store, p, request)
	if read().Rows[0].Operations != "1" {
		t.Fatal("replay inflated metric")
	}
	// Explicit synthetic source lifecycle prerequisite; no production cancellation claim.
	must(`update sales.customer_order set state='cancelled',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id='source'`, tenant)
	var version int64
	if e := pool.QueryRow(ctx, `select version from sales.customer_order where tenant_id=$1 and order_id='source'`, tenant).Scan(&version); e != nil {
		t.Fatal(e)
	}
	reverse := sv.Request{OperationID: "metric-reverse", Operation: "reverse", OriginalOperationID: "metric-issue", OrganizationID: "store", OrderID: "source", ProgramID: program, ExpectedOrderVersion: version, ProfileSHA256: store.ProfileSHA256()}
	storedValueHTTPApprove(t, store, review, storedValueHTTPPropose(t, store, p, reverse))
	v = read()
	if len(v.Rows) != 2 || v.Rows[1].Operation != "reverse" || v.Rows[1].PointsDelta != "-50.000000" {
		t.Fatal("reversal hidden", v)
	}
	empty, e := store.Metrics(ctx, p, "store", "reference-loyalty")
	if e != nil || len(empty.Rows) != 0 {
		t.Fatal("cross program", empty, e)
	}
	foreign := p
	foreign.TenantID = uuid.NewString()
	if _, e = store.Metrics(ctx, foreign, "store", program); e == nil {
		t.Fatal("foreign profile tenant")
	}
	if _, e = store.Metrics(ctx, p, "other", program); e == nil {
		t.Fatal("foreign profile org")
	}
	service := customerfeedback.NewService(db.NewCustomerFeedback(pool))
	verifier, token := catalogRoleIssuer(t)
	mux := http.NewServeMux()
	httpapi.CustomerFeedbackModule{Service: service}.Register(mux, verifier)
	httpapi.StoredValueModule{Service: store}.Register(mux, verifier)
	get := func(path string, want int) map[string]any {
		t.Helper()
		r := httptest.NewRequest("GET", path, nil)
		r.Header.Set("Authorization", "Bearer "+token("operator", tenant, []string{"surveys:read", "stored_value:read"}, []string{"store"}))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != want || !strings.Contains(w.Header().Get("Cache-Control"), "no-store") {
			t.Fatal(path, w.Code, w.Body.String())
		}
		var v map[string]any
		if e := json.Unmarshal(w.Body.Bytes(), &v); e != nil {
			t.Fatal(e)
		}
		return v
	}
	get("/v1/franchise/stored-value/metrics?organization_id=store&program_id="+program, 200)
	get("/v1/franchise/stored-value/metrics?organization_id=other&program_id="+program, 403)
	must(`insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active)values($1,'store','survey','Fixture','v1','Fixture',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1 hour',clock_timestamp()+interval '1 day',2,true)`, tenant)
	path := "/v1/reporting/surveys/survey?organization_id=store"
	summary := get(path, 200)
	if summary["responses"] != "0" || summary["available"] != false || summary["nps"] != nil {
		t.Fatal(summary)
	}
	for n, name := range []string{"customer", "customer-2"} {
		score := 10
		if n == 1 {
			score = 2
		}
		if _, e := service.Submit(ctx, customerfeedback.Scope{Tenant: tenant, Organization: "store", Customer: name}, "survey", customerfeedback.Submission{Score: score, Consent: true, ConsentVersion: "v1"}); e != nil {
			t.Fatal(e)
		}
		summary = get(path, 200)
		if n == 0 && summary["available"] != false {
			t.Fatal("minimum bypass", summary)
		}
	}
	if summary["responses"] != "2" || summary["nps"] != float64(0) || summary["available"] != true {
		t.Fatal(summary)
	}
	get("/v1/reporting/surveys/survey?organization_id=other", 403)
	must(`insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active)values($1,'store','expired','Fixture','v1','Fixture',clock_timestamp()-interval '3 day',clock_timestamp()-interval '2 day',clock_timestamp()-interval '1 day',2,false)`, tenant)
	get("/v1/reporting/surveys/expired?organization_id=store", 404)
	t.Log("ROLE_METRICS_SOURCE_OWNERS_PASS source_approved_issue_reverse=true pending_replay_program_isolation=true survey_real_submit_threshold_retention=true")
}

func TestRoleMetricsFactorySourceScope(t *testing.T) {
	ctx := context.Background()
	pool := roleMetricPool(t)
	tenant := roleMetricSeed(t, pool)
	if _, e := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'factory','factory','Fixture','factory')`, tenant); e != nil {
		t.Fatal(e)
	}
	maker := identity.Principal{TenantID: tenant, Subject: "maker", Permissions: map[string]struct{}{"*": {}}, Organizations: map[string]struct{}{"store": {}, "factory": {}}}
	ops := operations.NewService(db.NewOperations(pool), randomid.Generator{})
	po, e := ops.CreatePurchaseOrder(ctx, tenant, operations.PurchaseOrder{SupplierID: "supplier", DestinationOrganizationID: "store", Currency: "ARS", TotalMinorUnits: 100})
	if e != nil {
		t.Fatal(e)
	}
	supply, e := db.NewSerialSupply(pool)
	if e != nil {
		t.Fatal(e)
	}
	receipt, e := supply.BindPlan(ctx, maker, sc.PlanRequest{PurchaseOrderID: po.ID, CommandID: "plan", FactoryOrganizationID: "factory", DemandReference: "metric source fixture", PolicyCode: sc.PolicyCode, EvidenceSHA256: strings.Repeat("a", 64), Lines: []sc.Line{{ID: "line", VariantID: "variant", Quantity: 1}}})
	if e != nil {
		t.Fatal(e)
	}
	for _, kind := range []string{"submit", "confirm", "start", "register"} {
		command := sc.Command{PurchaseOrderID: po.ID, CommandID: kind, ExpectedVersion: receipt.Version, Kind: kind, EvidenceSHA256: strings.Repeat("a", 64)}
		if kind == "register" {
			command.LineID = "line"
			command.SerialNumber = "KPI-SOURCE"
		}
		receipt, e = supply.Apply(ctx, maker, command)
		if e != nil {
			t.Fatal(kind, e)
		}
	}
	store, e := db.NewRoleMetrics(pool)
	if e != nil {
		t.Fatal(e)
	}
	p := identity.Principal{TenantID: tenant, Subject: "reader", Permissions: map[string]struct{}{"supply:factory-read": {}, "supply:read": {}, "factory:read": {}}, Organizations: map[string]struct{}{"store": {}, "other": {}, "factory": {}}}
	for _, test := range []struct {
		kind, org string
		count     int
	}{{"factory-owned", "factory", 1}, {"factory-owned", "store", 0}, {"factory-destination", "store", 1}, {"factory-destination", "other", 0}, {"supply", "store", 1}, {"supply", "other", 0}} {
		v, e := store.Read(ctx, p, test.kind, test.org)
		if e != nil || len(v.Rows) != test.count {
			t.Fatal(test, v, e)
		}
		if test.kind == "supply" && test.count == 1 && v.Rows[0].Count != "1" {
			t.Fatal("unconnected purchase counted", v)
		}
	}
	t.Log("ROLE_METRICS_FACTORY_SCOPE_PASS actual_plan_submit_confirm_start_register=true factory_destination_distinct=true connected_plans_only=true")
}
