# Connected supply order creation

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-SUPPLY-CREATION"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Connected role supply from initial order and quantities through factory QA/replacement, receipt/reinspection/availability and durable recovery; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804/J2. Existing serial supply, Operations, approvals/inventory/outbox and OIDC/BFF/Next required. Only active tenant/org/supplier/catalog masters precede the actual order/quantity form.

## 3. Architecture contract

Original Operations service/repository and extracted BindPlan shared transaction; same planned receipt, tenant/order advisory lock, exact actor/hash/version. Quality view projects existing approval states and requester; state machine and human separation remain in existing owners.

## 4. Exact file manifest

```text
CREATE internal/platform/postgres/serial_supply_create.go
CREATE internal/platform/postgres/serial_supply_create_integration_test.go
CREATE internal/platform/postgres/serial_supply_role_browser_integration_test.go
CREATE internal/serialsupply/create.go
CREATE internal/serialsupply/create_fuzz_test.go
CREATE internal/serialsupply/role_golden_test.go
```

## 5. Materialization blocks

### FILE: `internal/platform/postgres/serial_supply_create.go`

```yaml
block_id: "GO-CONNECTED-SUPPLY-CREATION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "86f2a13701ae83f484d23dabfaf99f542f68da5af1c05372a788c17bc2dc6653"
variables: []
secrets_allowed: false
```

````go
// AUTHORED single-transaction composition; no new pricing or stock rules.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"errors"
	"github.com/jackc/pgx/v5"
)

// Operations requests an aggregate ID first and then an event ID. Only the
// aggregate is a preallocated recovery reference; events retain random IDs.
type supplyCreationIDs struct{ first string }

func (g *supplyCreationIDs) New() string {
	if g.first != "" {
		v := g.first
		g.first = ""
		return v
	}
	return randomid.Generator{}.New()
}
func (s *SerialSupply) Create(ctx context.Context, p identity.Principal, r sc.CreateRequest) (sc.Receipt, error) {
	if s == nil || !r.Valid() || p.TenantID == "" {
		return sc.Receipt{}, sc.ErrInvalid
	}
	if !approvalPrincipal(p, p.TenantID, r.DestinationOrganizationID, "supply:plan") {
		return sc.Receipt{}, sc.ErrNotFound
	}
	_, hash, err := sc.Canonical(r)
	if err != nil {
		return sc.Receipt{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return sc.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	// Separate SQL parameters and a namespace; hash collisions only serialize
	// unrelated orders, never grant access or merge their durable identities.
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended('serial-supply-create:'||$1::text||':'||$2::text,0))`, p.TenantID, r.PurchaseOrderID); err != nil {
		return sc.Receipt{}, err
	}
	old, err := readSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID)
	if err == nil {
		plan, e := readSupplyPlan(ctx, tx, p.TenantID, r.PurchaseOrderID, false)
		if e != nil {
			return sc.Receipt{}, e
		}
		if plan.DestinationOrganizationID != r.DestinationOrganizationID || old.Actor != p.Subject || old.Kind != "planned" || old.RequestSHA256 != hash {
			return sc.Receipt{}, sc.ErrConflict
		}
		old.Replay = true
		return old, tx.Commit(ctx)
	}
	if !errors.Is(err, sc.ErrNotFound) {
		return sc.Receipt{}, err
	}
	svc := operations.NewService(operationsTxRepository{tx: tx}, &supplyCreationIDs{first: r.PurchaseOrderID})
	_, err = svc.CreatePurchaseOrder(ctx, p.TenantID, operations.PurchaseOrder{SupplierID: r.SupplierID, DestinationOrganizationID: r.DestinationOrganizationID, Currency: r.Currency, TotalMinorUnits: r.TotalMinorUnits})
	if err != nil {
		return sc.Receipt{}, supplyError(err)
	}
	receipt, err := bindSupplyPlanInTx(ctx, tx, p, r.PlanRequest, hash)
	if err != nil {
		return sc.Receipt{}, err
	}
	return receipt, supplyError(tx.Commit(ctx))
}
````

### FILE: `internal/platform/postgres/serial_supply_create_integration_test.go`

```yaml
block_id: "GO-CONNECTED-SUPPLY-CREATION:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8ac9c625e0afb5f0ee64d95e348d89e6effaa7880765cecb79d5356e3406957a"
variables: []
secrets_allowed: false
```

````go
// AUTHORED focused proof for the new order/plan transaction composition.
package postgres_test

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestSerialSupplyCreateAtomic(t *testing.T) {
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("owned fixture required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'supply-create','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Fixture','store'),($1,'factory','factory','Fixture','factory')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status)values($1,'supplier','supplier','Fixture','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Fixture','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Fixture','{}','active')`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	store, e := db.NewSerialSupply(pool)
	if e != nil {
		t.Fatal(e)
	}
	p := identity.Principal{TenantID: tenant, Subject: "buyer", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"supply:plan": {}, "supply:read": {}}}
	r := sc.CreateRequest{PlanRequest: sc.PlanRequest{PurchaseOrderID: uuid.NewString(), CommandID: uuid.NewString(), FactoryOrganizationID: "factory", DemandReference: "Demand <á> & fixture", PolicyCode: sc.PolicyCode, EvidenceSHA256: strings.Repeat("a", 64), Lines: []sc.Line{{ID: "line-1", VariantID: "variant", Quantity: 2}}}, DestinationOrganizationID: "store", SupplierID: "supplier", Currency: "ARS", TotalMinorUnits: 9007199254740993}
	snapshot := func() string {
		t.Helper()
		out := ""
		for _, table := range []string{"procurement.purchase_order", "procurement.serial_supply_plan", "procurement.serial_supply_line", "procurement.serial_supply_step", "platform.outbox_event"} {
			var text string
			if e = pool.QueryRow(ctx, "select coalesce(jsonb_agg(to_jsonb(x) order by to_jsonb(x)::text),'[]'::jsonb)::text from "+table+" x where tenant_id=$1", tenant).Scan(&text); e != nil {
				t.Fatal(e)
			}
			out += text
		}
		return out
	}
	before := snapshot()
	bad := r
	bad.Lines = []sc.Line{{ID: "line-1", VariantID: "missing", Quantity: 2}}
	if _, e = store.Create(ctx, p, bad); e == nil {
		t.Fatal("invalid variant accepted")
	}
	if snapshot() != before {
		t.Fatal("order/outbox escaped failed plan")
	}
	// Failure after the plan receipt also rolls back its purchase owner/outbox.
	if _, e = pool.Exec(ctx, `create function public.supply_role_fail()returns trigger language plpgsql as $$begin if new.aggregate_type='serial-supply' and new.payload->>'command_id'='forced-rollback' then raise exception 'fixture';end if;return new;end$$;create trigger supply_role_fail before insert on platform.outbox_event for each row execute function public.supply_role_fail()`); e != nil {
		t.Fatal(e)
	}
	bad = r
	bad.CommandID = "forced-rollback"
	if _, e = store.Create(ctx, p, bad); e == nil {
		t.Fatal("forced outbox accepted")
	}
	if snapshot() != before {
		t.Fatal("receipt rollback leaked")
	}
	if _, e = pool.Exec(ctx, `drop trigger supply_role_fail on platform.outbox_event;drop function public.supply_role_fail()`); e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	out := make(chan sc.Receipt, 4)
	errs := make(chan error, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); v, e := store.Create(ctx, p, r); out <- v; errs <- e }()
	}
	wg.Wait()
	close(out)
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	fresh, replay := 0, 0
	for v := range out {
		if v.Replay {
			replay++
		} else {
			fresh++
		}
		if v.PurchaseOrderID != r.PurchaseOrderID || v.Version != 1 || v.Actor != p.Subject || v.Kind != "planned" {
			t.Fatal(v)
		}
	}
	if fresh != 1 || replay != 3 {
		t.Fatal(fresh, replay)
	}
	durable := snapshot()
	for _, variation := range []string{"actor", "hash", "org"} {
		other := p
		bad = r
		switch variation {
		case "actor":
			other.Subject = "other"
		case "hash":
			bad.TotalMinorUnits++
		case "org":
			bad.DestinationOrganizationID = "other"
		}
		if _, e = store.Create(ctx, other, bad); e == nil {
			t.Fatal("replay mismatch", variation)
		}
	}
	if snapshot() != durable {
		t.Fatal("negative replay changed history")
	}
	plan, e := store.Plan(ctx, p, r.PurchaseOrderID, "")
	if e != nil || plan.TotalMinorUnits != r.TotalMinorUnits || plan.Currency != "ARS" || len(plan.Lines) != 1 || plan.Version != 1 {
		t.Fatal(plan, e)
	}
	rebuilt, e := db.NewSerialSupply(pool)
	if e != nil {
		t.Fatal(e)
	}
	receipt, e := rebuilt.CommandReceipt(ctx, p, r.PurchaseOrderID, r.CommandID)
	_, hash, _ := sc.Canonical(r)
	if e != nil || receipt.RequestSHA256 != hash {
		t.Fatal("durable recovery", receipt, e)
	}
	var counts string
	if e = pool.QueryRow(ctx, `select jsonb_build_array((select count(*)from procurement.purchase_order where tenant_id=$1),(select count(*)from procurement.serial_supply_plan where tenant_id=$1),(select count(*)from procurement.serial_supply_step where tenant_id=$1),(select count(*)from platform.outbox_event where tenant_id=$1))::text`, tenant).Scan(&counts); e != nil || counts != "[1, 1, 1, 2]" {
		t.Fatal(counts, e)
	}
	rawPlan, _ := json.Marshal(plan)
	if !strings.Contains(string(rawPlan), `"total_minor_units":"9007199254740993"`) {
		t.Fatal("unsafe numeric projection", string(rawPlan))
	}
	t.Log("SUPPLY_CREATE_ATOMIC_PASS concurrency=1new3replay tables5_rollback=true actor_hash_scope_bound=true durable_recovery=true exact_int64=true")
}
````

### FILE: `internal/platform/postgres/serial_supply_role_browser_integration_test.go`

```yaml
block_id: "GO-CONNECTED-SUPPLY-CREATION:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c074229f95474604eeb44338a39d6885cd7e179642e05b2241ba378f080279c4"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED role browser: only tenant/org/supplier/catalog masters seeded; all J2 writes from forms.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSupplyRoleBrowser(t *testing.T) {
	if os.Getenv("ELITE_SUPPLY_ROLE_BROWSER") != "1" {
		t.Skip("explicit local browser fixture required")
	}
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Fatal("owned database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := uuid.NewString()

	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'supply-role','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Fixture','store'),($1,'factory','factory','Fixture','factory')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status)values($1,'supplier','supplier','Fixture','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Fixture','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Fixture','{}','active')`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	store, e := db.NewSerialSupply(pool)
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := catalogRoleIssuer(t)
	identities := map[string]any{}

	for _, name := range []string{"buyer", "maker", "factory-qa", "receiver", "receipt-qa", "reader", "foreign", "unprivileged"} {
		permissions := []string{"supply:read"}
		orgs := []string{"store"}
		switch name {
		case "buyer":
			permissions = append(permissions, "supply:plan")
		case "maker", "factory-qa":
			permissions = []string{"supply:factory", "supply:factory-read"}
			orgs = []string{"factory"}
		case "receiver":
			permissions = append(permissions, "supply:receive", "supply:inspect", "supply:release")
		case "receipt-qa":
			permissions = append(permissions, "supply:release")
		case "foreign":
			orgs = []string{"foreign"}
		case "unprivileged":
			permissions = []string{}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": orgs, "accessToken": token(name, tenant, permissions, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.SerialSupplyModule{Service: store}.Register(mux, verifier)
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			mu.Lock()
			counts[r.URL.Path]++
			mu.Unlock()
		}
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()
	web, node := os.Getenv("ELITE_WEB_ROOT"), os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("fixed runtime/root required")
	}
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	env := []string{}
	for _, v := range os.Environ() {
		k := strings.ToUpper(strings.SplitN(v, "=", 2)[0])
		if strings.HasPrefix(k, "ELITE_") || strings.HasPrefix(k, "CATALOG_") || strings.HasPrefix(k, "PUBLIC_") || strings.HasPrefix(k, "ENTERPRISE_") || strings.Contains(k, "DATABASE_URL") || k == "APP_BASE_URL" || k == "AUTH_SESSION_SECRET" || k == "BUSINESS_CONFIG_FILE" {
			continue
		}
		env = append(env, v)
	}
	identitiesJSON, _ := json.Marshal(identities)
	env = append(env, "ELITE_SUPPLY_ROLE_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_SUPPLY_IDENTITIES="+string(identitiesJSON), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL,
		"APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-supply-role-"+uuid.NewString(), "BUSINESS_CONFIG_FILE=supply.role.fixture.json",
		"CATALOG_RELEASE_ENABLED=false", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=0", "ENTERPRISE_TENANT_CODE=role-catalog", "ENTERPRISE_ORGANIZATION_CODE=role-org", "NEXT_TELEMETRY_DISABLED=1")
	artifacts, e := os.MkdirTemp(web, "supply-role-artifacts-")
	if e != nil {
		t.Fatal(e)
	}
	log, e := os.Create(filepath.Join(artifacts, "next.log"))
	if e != nil {
		t.Fatal(e)
	}
	defer log.Close()
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env, server.Stdout, server.Stderr = web, env, log, log
	if e = server.Start(); e != nil {
		t.Fatal(e)
	}
	defer func() { server.Process.Kill(); server.Wait() }()
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		r, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			r.Body.Close()
			if r.StatusCode == 200 {
				ready = true
				break
			}
		}
		select {
		case <-ctx.Done():
			t.Fatal("startup cancelled")
		case <-time.After(100 * time.Millisecond):
		}
	}
	if !ready {
		t.Fatal("Next startup", artifacts)
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/supply-role-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, e := command.CombinedOutput()
	if x := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); x != nil {
		t.Fatal(x)
	}
	if e != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("catalog role browser %v\n%s\nartifacts=%s", e, output, artifacts)
	}
	var actual string

	e = pool.QueryRow(ctx, `select jsonb_build_array(
 (select count(*)from procurement.purchase_order where tenant_id=$1 and state='received'),
 (select count(*)from factory.production_unit where tenant_id=$1),
 (select count(*)from inventory.stock_unit where tenant_id=$1 and state='available'),
 (select count(*)from procurement.serial_supply_step where tenant_id=$1),
 (select count(*)from approval.request where tenant_id=$1),
 (select count(*)from approval.decision where tenant_id=$1))::text`, tenant).Scan(&actual)
	if e != nil || actual != "[1, 3, 2, 23, 6, 6]" {
		t.Fatal("durable effects", actual, e)
	}
	mu.Lock()
	countJSON, _ := json.Marshal(counts)
	total := 0
	for _, n := range counts {
		total += n
	}
	mu.Unlock()
	if total != 23 {
		t.Fatal("unexpected backend writes", total)
	}
	if e = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countJSON, 0600); e != nil {
		t.Fatal(e)
	}
	t.Logf("SUPPLY_ROLE_BROWSER_PASS actual_order_creation=true JWE_RS256_JWKS=true create_and_receive_response_loss_GET_only=true three_factory_serials_one_rejected_two_available=true quality_reinspection=true backend_POSTs=23 artifacts=%s", artifacts)
}
````

### FILE: `internal/serialsupply/create.go`

```yaml
block_id: "GO-CONNECTED-SUPPLY-CREATION:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "437c494f2b98c6172fda0ef36dec4ca8b7362ac53f08dc724165ad57daba8c6d"
variables: []
secrets_allowed: false
```

````go
// AUTHORED typed composition input. Operations remains the purchase owner.
package serialsupply

import "regexp"

type CreateRequest struct {
	PlanRequest
	DestinationOrganizationID string `json:"destination_organization_id"`
	SupplierID                string `json:"supplier_id"`
	Currency                  string `json:"currency"`
	TotalMinorUnits           int64  `json:"total_minor_units,string"`
}

func (r CreateRequest) Valid() bool {
	return r.PlanRequest.Valid() && ValidID(r.DestinationOrganizationID) && ValidID(r.SupplierID) && regexp.MustCompile(`^[A-Z]{3}$`).MatchString(r.Currency) && r.TotalMinorUnits >= 0
}
````

### FILE: `internal/serialsupply/create_fuzz_test.go`

```yaml
block_id: "GO-CONNECTED-SUPPLY-CREATION:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "661e109bdf7387359af8d23b207476f44af2d207141e441e2763963f8c493918"
variables: []
secrets_allowed: false
```

````go
package serialsupply

import (
	"elite.local/enterprise/internal/approval"
	"encoding/json"
	"testing"
)

func FuzzSupplyCreateBounds(f *testing.F) {
	f.Add([]byte(`{"command_id":"golden-command","currency":"ARS","demand_reference":"Demanda á \u003c\u0026\u003e \u2028 interior","destination_organization_id":"store","evidence_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","factory_organization_id":"factory","lines":[{"id":"line","quantity":2,"variant_id":"variant"}],"policy_code":"strict-serial-reference/v1","purchase_order_id":"golden-po","supplier_id":"supplier","total_minor_units":"9007199254740993"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"command_id":"a","command_id":"b"}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32768 {
			return
		}
		canon, _, e := approval.CanonicalPayload(raw)
		if e != nil {
			return
		}
		var value CreateRequest
		if json.Unmarshal(canon, &value) != nil {
			return
		}
		before, _ := json.Marshal(value)
		valid := value.Valid()
		after, _ := json.Marshal(value)
		if string(before) != string(after) {
			t.Fatal("validation mutated input")
		}
		if valid {
			if !value.PlanRequest.Valid() || value.TotalMinorUnits < 0 || !ValidID(value.SupplierID) || !ValidID(value.DestinationOrganizationID) {
				t.Fatal("unbounded source")
			}
			total := 0
			for _, line := range value.Lines {
				total += line.Quantity
			}
			if total < 1 || total > 1000 {
				t.Fatal("quantity scope")
			}
		}
	})
}
````

### FILE: `internal/serialsupply/role_golden_test.go`

```yaml
block_id: "GO-CONNECTED-SUPPLY-CREATION:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9a3c477f84d07941cf3715b4deb13a7c24a4c3b4b534ea20730b69b41424dd41"
variables: []
secrets_allowed: false
```

````go
package serialsupply

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSupplyRoleCommandGoldens(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/supply/role-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var cases []struct {
		Name      string
		Payload   json.RawMessage
		Canonical string
		SHA256    string
	}
	if e = json.Unmarshal(raw, &cases); e != nil {
		t.Fatal(e)
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var v any
			if c.Name == "create" {
				v = new(CreateRequest)
			} else {
				v = new(Command)
			}
			if e = json.Unmarshal(c.Payload, v); e != nil {
				t.Fatal(e)
			}
			switch r := v.(type) {
			case *CreateRequest:
				if !r.Valid() {
					t.Fatal("invalid create")
				}
			case *Command:
				if !r.Valid() {
					t.Fatal("invalid command")
				}
			}
			raw, hash, e := Canonical(v)
			if e != nil || string(raw) != c.Canonical || hash != c.SHA256 {
				t.Fatal(string(raw), hash, e)
			}
		})
	}
}
````

## 6. Configuration surface

docs/SUPPLY_ROLE_REFERENCE.md. features.supply_portal opt-in, supply:read/factory-read plus action permission. Source order IDs are preallocated recovery references, event IDs remain random. Explicit currency/total,32lines/1000units,100unit pages. No new migration.

## 7. Dependency bill

16new and7changed AUTHORED unavoidable composition/transport/UI/config/test glue. No new dependency/upstream/corporate attribution; existing pins and adapted-source licenses retained.

## 8. Apply order

Select creation companion with GO-CONNECTED-SERIAL-SUPPLY0.2.0; full reference also selects TS-SERIAL-SUPPLY-PORTAL holding shared Go/TS goldens and GO-CONNECTED-CATALOG-AUTHORING holding the reused synthetic browser issuer.

## 9. Verification

Actual Next/BFF/Go/PG/Chromium JWE/RS256/JWKS,23writes:1received order,3factory serials with1rejected/replaced,2available stock,6quality requests/decisions. Create/receive responses lost then recovered byGET after reload without extraPOST. Creation1new3replay,5table rollback including late outbox, actor/hash/org/int64;4Go/TS goldens;4BFF/6contract tests. Full webpack/types. Fuzz2s3seeds99596executions. Visual desktop/390px; nav key correction checked in exact web type closure.

## 10. Reconstruction evidence

SUPPLY_ROLE_RELEASE_V402.md/json: exact source/delta, RED/PASS receipts and profile rebuilds. Does not close allT2804; warranty/J5/CMS/KPIs/private locale remain. No production certification.

