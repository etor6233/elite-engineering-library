# Go FinOps Core

## 1. Metadata

```yaml
pack_id: "GO-FINOPS-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el gobierno de costo de nube (COST-FINOPS): todo recurso etiquetado e inventariado (nada se escapa de la facturación), audit de drift/leak, budget por tenant fail-closed y teardown que contabiliza cada recurso."
stacks: ["Go 1.26.7", "PostgreSQL 18.6"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "PG-TX-FOUNDATION 0.1.x", "SECURE-OPS-DELIVERY-CORE 1.1.x"]
incompatible_with: ["recurso sin tags tenant/environment/owner", "gasto sin cap por tenant", "teardown sin inventario"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/"]
verified_at: "2026-09-02"
```

Gobernado por el contrato `COST-FINOPS` (48 superficies: tags/budgets/unit cost/alerts/caps/capacity plan/exit cost) y el corpus SRE (`SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`: IaC teardown, drift, error budget). El código es `AUTHORED`; no copia ningún producto cloud.

## 2. Applicability

Use este pack para que al crear/desplegar una franquicia **nada se escape**: cada recurso (Cloud Run, app, contenedor, DB, queue) queda etiquetado e inventariado, con budget por tenant y teardown verificable. Detecta untagged, drift y leak antes de teardown.

Rechace este pack para: recursos sin tags obligatorios; o gasto sin cap por tenant.

## 3. Architecture contract

- **Ownership**: `internal/finops` gobierna identidad/etiquetado, inventario, budget y audit. La migración `0043` gobierna `finops.*` durable.
- **Invariantes**: (1) tags `tenant`/`environment`/`owner` obligatorios (nada sin atribuir). (2) inventario único por (tenant, provider, id); duplicados rechazados. (3) budget por tenant fail-closed (`spent <= cap`). (4) audit detecta untagged, drift (desconocido) y leak (faltante); teardown exige audit limpio.
- **Data flow**: `Register` (validado) → inventario; `Audit(actual)` → untagged/unknown/missing; `Budget.Record` → fail-closed; `PlanTeardown` → contabiliza todo.
- **Failure modes**: recurso sin tag → `ErrMissingTag`; duplicado → `ErrDuplicate`; sobre-cap → `ErrOverBudget`; audit sucio → teardown bloqueado.
- **Seguridad/privacidad**: tenant-scoped; tags para atribución de billing.
- **Performance budget**: O(1) register/budget; O(N) audit por tenant.
- **Operación/migración/rollback**: migración `0043` up/down atómica; inventario reconstruible desde el provider.

## 4. Exact file manifest

```text
CREATE internal/finops/resource.go
CREATE internal/finops/inventory.go
CREATE internal/finops/budget.go
CREATE internal/finops/teardown.go
CREATE internal/finops/finops_test.go
CREATE db/migrations/0043_finops.up.sql
CREATE db/migrations/0043_finops.down.sql
CREATE db/tests/0043_finops.test.sql
```

## 5. Materialization blocks

### FILE: `internal/finops/resource.go`
```yaml
block_id: "GO-FINOPS-CORE:internal/finops/resource.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "dd138e2cadee55b2e5ff9571cd97e2f6046f0376788fc8d7c58c3480416ec181"
variables: []
secrets_allowed: false
```
````go
// Package finops provides cloud-cost governance: every resource is tagged and
// inventoried (nothing escapes billing attribution), spend is budgeted per
// tenant fail-closed, and teardown accounts for every resource. It is AUTHORED
// over the COST-FINOPS surface contract and the SRE corpus.
package finops

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	providerRe = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
	kindRe     = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)
	regionRe   = regexp.MustCompile(`^[a-z0-9-]{0,63}$`)

	ErrInvalidResource = errors.New("finops: invalid resource")
	ErrMissingTag      = errors.New("finops: missing required tag")
	ErrDuplicate       = errors.New("finops: duplicate resource")
)

// RequiredTags are the tags every resource must carry so cost is attributable
// (tenant, environment, owner). A resource without them cannot be billed
// correctly and is rejected.
var RequiredTags = []string{"tenant", "environment", "owner"}

// Resource is an inventoried cloud resource with a monthly cost estimate.
type Resource struct {
	TenantID               string
	Provider               string // aws | gcp | azure | cloudflare ...
	ID                     string // provider resource id
	Kind                   string // cloud_run | compute | database | storage ...
	Region                 string
	Tags                   map[string]string
	EstMonthlyMinorUnits   int64
}

// Validate enforces tagging and identity so nothing escapes attribution.
func (r Resource) Validate() error {
	if strings.TrimSpace(r.TenantID) == "" || len(r.TenantID) > 64 {
		return fmt.Errorf("%w: tenant", ErrInvalidResource)
	}
	if !providerRe.MatchString(r.Provider) {
		return fmt.Errorf("%w: provider", ErrInvalidResource)
	}
	if strings.TrimSpace(r.ID) == "" || len(r.ID) > 200 {
		return fmt.Errorf("%w: id", ErrInvalidResource)
	}
	if !kindRe.MatchString(r.Kind) {
		return fmt.Errorf("%w: kind", ErrInvalidResource)
	}
	if !regionRe.MatchString(r.Region) {
		return fmt.Errorf("%w: region", ErrInvalidResource)
	}
	for _, tag := range RequiredTags {
		if strings.TrimSpace(r.Tags[tag]) == "" {
			return fmt.Errorf("%w: %s", ErrMissingTag, tag)
		}
	}
	if r.EstMonthlyMinorUnits < 0 {
		return fmt.Errorf("%w: negative cost", ErrInvalidResource)
	}
	return nil
}

// Key returns the provider identity (dedup across tenants is enforced by the
// inventory, not by this key).
func (r Resource) Key() string {
	return r.Provider + "\x00" + r.ID
}
````

### FILE: `internal/finops/inventory.go`
```yaml
block_id: "GO-FINOPS-CORE:internal/finops/inventory.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "97b98838e50cc52f2a006e6c98cf35c13576ecb6181b4ce6477f638c2338ce82"
variables: []
secrets_allowed: false
```
````go
package finops

import (
	"sort"
	"sync"
)

// Inventory holds every deployed resource so nothing escapes billing. A
// resource is unique per (tenant, provider, id); duplicates are rejected.
type Inventory struct {
	mu        sync.RWMutex
	resources map[string]Resource // key: tenant + "\x00" + provider + "\x00" + id
}

// NewInventory returns an empty inventory.
func NewInventory() *Inventory {
	return &Inventory{resources: make(map[string]Resource)}
}

func invKey(tenant, provider, id string) string {
	return tenant + "\x00" + provider + "\x00" + id
}

// Register adds a validated, tagged resource; duplicates are rejected.
func (i *Inventory) Register(r Resource) error {
	if err := r.Validate(); err != nil {
		return err
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	k := invKey(r.TenantID, r.Provider, r.ID)
	if _, ok := i.resources[k]; ok {
		return ErrDuplicate
	}
	i.resources[k] = r
	return nil
}

// List returns the tenant's resources sorted by (provider, id).
func (i *Inventory) List(tenant string) []Resource {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var out []Resource
	for _, r := range i.resources {
		if r.TenantID == tenant {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a].Provider != out[b].Provider {
			return out[a].Provider < out[b].Provider
		}
		return out[a].ID < out[b].ID
	})
	return out
}

// AuditResult reports what escaped the inventory: untagged actual resources,
// unknown (drifted-in) resources, and missing (leaked) inventory resources.
type AuditResult struct {
	Untagged []Resource
	Unknown  []Resource
	Missing  []Resource
}

// Clean reports whether nothing escaped: no untagged, no drift, no leak.
func (a AuditResult) Clean() bool {
	return len(a.Untagged) == 0 && len(a.Unknown) == 0 && len(a.Missing) == 0
}

// Audit compares the tenant's actual cloud resources against the inventory.
func (i *Inventory) Audit(tenant string, actual []Resource) AuditResult {
	i.mu.RLock()
	defer i.mu.RUnlock()
	expected := make(map[string]Resource)
	for _, r := range i.resources {
		if r.TenantID == tenant {
			expected[r.Key()] = r
		}
	}
	seen := make(map[string]bool)
	var res AuditResult
	for _, r := range actual {
		if r.Validate() != nil {
			res.Untagged = append(res.Untagged, r)
			continue
		}
		if _, ok := expected[r.Key()]; !ok {
			res.Unknown = append(res.Unknown, r) // drifted-in, not inventoried
			continue
		}
		seen[r.Key()] = true
	}
	for k, r := range expected {
		if !seen[k] {
			res.Missing = append(res.Missing, r) // leaked out of inventory
		}
	}
	return res
}

// MonthlyCost returns the tenant's total estimated monthly cost.
func (i *Inventory) MonthlyCost(tenant string) int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var total int64
	for _, r := range i.resources {
		if r.TenantID == tenant {
			total += r.EstMonthlyMinorUnits
		}
	}
	return total
}
````

### FILE: `internal/finops/budget.go`
```yaml
block_id: "GO-FINOPS-CORE:internal/finops/budget.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "e56d8e03939c99c28483e77577261f18dc2be878567941238026ad1c864eef60"
variables: []
secrets_allowed: false
```
````go
package finops

import (
	"errors"
	"sync"
)

// ErrOverBudget reports a spend that would exceed the tenant's budget.
var ErrOverBudget = errors.New("finops: over budget")

// Budget tracks spend per tenant against a cap, fail-closed.
type Budget struct {
	mu    sync.Mutex
	caps  map[string]int64
	spent map[string]int64
}

// NewBudget returns an empty budget tracker.
func NewBudget() *Budget {
	return &Budget{caps: make(map[string]int64), spent: make(map[string]int64)}
}

// SetCap sets the tenant's spending cap (clamped to >= 0).
func (b *Budget) SetCap(tenant string, cap int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if cap < 0 {
		cap = 0
	}
	b.caps[tenant] = cap
}

// Record adds spend; it fails closed if the cap would be exceeded.
func (b *Budget) Record(tenant string, amount int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if amount < 0 {
		return errors.New("finops: negative spend")
	}
	if b.spent[tenant]+amount > b.caps[tenant] {
		return ErrOverBudget
	}
	b.spent[tenant] += amount
	return nil
}

// Remaining returns the tenant's remaining budget.
func (b *Budget) Remaining(tenant string) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.caps[tenant] - b.spent[tenant]
}

// OverBudget reports whether the tenant has no remaining budget.
func (b *Budget) OverBudget(tenant string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.spent[tenant] >= b.caps[tenant]
}
````

### FILE: `internal/finops/teardown.go`
```yaml
block_id: "GO-FINOPS-CORE:internal/finops/teardown.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "e80a4310601fd2090fba63cc9a4074ead5509caa437e36f3a5fb8d5eff10e116"
variables: []
secrets_allowed: false
```
````go
package finops

// TeardownPlan accounts for every tenant resource before teardown. Nothing may
// escape: teardown is safe only when the audit is clean (no untagged, no
// drift, no leak).
type TeardownPlan struct {
	Resources []Resource
	Audit     AuditResult
}

// PlanTeardown builds the teardown plan for a tenant from the actual cloud
// resource set. Complete means the audit is clean and every resource is
// accounted for.
func PlanTeardown(inv *Inventory, tenant string, actual []Resource) TeardownPlan {
	return TeardownPlan{
		Resources: inv.List(tenant),
		Audit:     inv.Audit(tenant, actual),
	}
}

// Complete reports whether teardown can proceed without escaping anything.
func (p TeardownPlan) Complete() bool {
	return p.Audit.Clean()
}
````

### FILE: `internal/finops/finops_test.go`
```yaml
block_id: "GO-FINOPS-CORE:internal/finops/finops_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "fcbd6c3a64d146ab7e7aa9c189f916a023f4d82d8e23399c2a88f100d2caa5d2"
variables: []
secrets_allowed: false
```
````go
package finops

import (
	"errors"
	"testing"
)

func res(tenant, provider, id string, tags map[string]string) Resource {
	if tags == nil {
		tags = map[string]string{"tenant": tenant, "environment": "prod", "owner": "platform"}
	}
	return Resource{
		TenantID: tenant, Provider: provider, ID: id, Kind: "cloud_run",
		Region: "us-central1", Tags: tags, EstMonthlyMinorUnits: 100,
	}
}

func TestResourceValidateTags(t *testing.T) {
	if err := res("t", "gcp", "r1", nil).Validate(); err != nil {
		t.Fatalf("valid resource rejected: %v", err)
	}
	bad := res("t", "gcp", "r1", map[string]string{"tenant": "t"}) // missing env/owner
	if err := bad.Validate(); !errors.Is(err, ErrMissingTag) {
		t.Fatalf("missing tag accepted: %v", err)
	}
	neg := res("t", "gcp", "r1", nil)
	neg.EstMonthlyMinorUnits = -1
	if err := neg.Validate(); !errors.Is(err, ErrInvalidResource) {
		t.Fatalf("negative cost accepted: %v", err)
	}
}

func TestInventoryDedupAndCost(t *testing.T) {
	inv := NewInventory()
	if err := inv.Register(res("t", "gcp", "r1", nil)); err != nil {
		t.Fatal(err)
	}
	if err := inv.Register(res("t", "gcp", "r1", nil)); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate accepted: %v", err)
	}
	_ = inv.Register(res("t", "aws", "r2", nil))
	if got := inv.MonthlyCost("t"); got != 200 {
		t.Fatalf("expected 200, got %d", got)
	}
	if len(inv.List("t")) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(inv.List("t")))
	}
}

func TestInventoryRejectsUntagged(t *testing.T) {
	inv := NewInventory()
	// Register rejects untagged fail-closed
	if err := inv.Register(res("t", "aws", "r2", map[string]string{"tenant": "t"})); !errors.Is(err, ErrMissingTag) {
		t.Fatalf("untagged resource accepted: %v", err)
	}
}

func TestAuditDetectsEscapes(t *testing.T) {
	inv := NewInventory()
	_ = inv.Register(res("t", "gcp", "r1", nil))

	// actual: r1 present (good), r2 untagged, r3 drifted-in; r1 expected but a
	// leak scenario is tested separately.
	actual := []Resource{
		res("t", "gcp", "r1", nil),
		Resource{TenantID: "t", Provider: "aws", ID: "r2", Kind: "cloud_run", Region: "us-east-1", Tags: map[string]string{"tenant": "t"}}, // untagged
		res("t", "gcp", "r3", nil), // drifted-in (not inventoried)
	}
	a := inv.Audit("t", actual)
	if len(a.Untagged) != 1 || len(a.Unknown) != 1 {
		t.Fatalf("expected 1 untagged + 1 unknown, got %+v", a)
	}
	if a.Clean() {
		t.Fatal("audit should not be clean")
	}
}

func TestAuditDetectsLeak(t *testing.T) {
	inv := NewInventory()
	_ = inv.Register(res("t", "gcp", "r1", nil))
	_ = inv.Register(res("t", "aws", "r2", nil))

	// actual only has r1 → r2 leaked out of inventory
	a := inv.Audit("t", []Resource{res("t", "gcp", "r1", nil)})
	if len(a.Missing) != 1 || a.Missing[0].ID != "r2" {
		t.Fatalf("expected r2 missing, got %+v", a.Missing)
	}
}

func TestBudgetFailClosed(t *testing.T) {
	b := NewBudget()
	b.SetCap("t", 100)
	if err := b.Record("t", 60); err != nil {
		t.Fatal(err)
	}
	if b.Remaining("t") != 40 {
		t.Fatalf("expected 40 remaining, got %d", b.Remaining("t"))
	}
	if err := b.Record("t", 41); !errors.Is(err, ErrOverBudget) {
		t.Fatalf("expected ErrOverBudget, got %v", err)
	}
	if b.OverBudget("t") {
		t.Fatal("should not be over budget at 60/100")
	}
	if err := b.Record("t", 40); err != nil {
		t.Fatal(err)
	}
	if !b.OverBudget("t") {
		t.Fatal("expected over-budget when spent == cap")
	}
}

func TestTeardownPlanAccountsAll(t *testing.T) {
	inv := NewInventory()
	_ = inv.Register(res("t", "gcp", "r1", nil))

	// actual matches inventory → clean teardown
	actual := []Resource{res("t", "gcp", "r1", nil)}
	if !PlanTeardown(inv, "t", actual).Complete() {
		t.Fatal("plan should be complete when actual == inventory")
	}
	// actual has an untagged drift → not complete
	drift := []Resource{
		res("t", "gcp", "r1", nil),
		Resource{TenantID: "t", Provider: "aws", ID: "x", Kind: "compute", Region: "us-east-1", Tags: map[string]string{"tenant": "t"}},
	}
	if PlanTeardown(inv, "t", drift).Complete() {
		t.Fatal("plan should be incomplete when something escaped")
	}
}
````

### FILE: `db/migrations/0043_finops.up.sql`
```yaml
block_id: "GO-FINOPS-CORE:db/migrations/0043_finops.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5736c7a328b34a450322fd7dfa9424779f97cdd8c0fd448cc7e347927b2c0f37"
variables: []
secrets_allowed: false
```
````sql
begin;

create schema if not exists finops;

create table finops.resource (
  tenant_id uuid not null,
  provider text not null,
  resource_id text not null,
  kind text not null,
  region text not null default '',
  tags jsonb not null,
  estimated_monthly_minor_units bigint not null check (estimated_monthly_minor_units >= 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, provider, resource_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (provider ~ '^[a-z][a-z0-9_-]{0,31}$'),
  check (kind ~ '^[a-z][a-z0-9_-]{0,63}$'),
  check (jsonb_typeof(tags) = 'object'),
  check (tags ? 'tenant' and tags ? 'environment' and tags ? 'owner')
);

create table finops.budget (
  tenant_id uuid not null,
  period text not null check (period ~ '^\d{4}-\d{2}$'),
  cap_minor_units bigint not null check (cap_minor_units >= 0),
  spent_minor_units bigint not null default 0 check (spent_minor_units >= 0),
  primary key (tenant_id, period),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (spent_minor_units <= cap_minor_units)
);

create table finops.cost_entry (
  tenant_id uuid not null,
  provider text not null,
  resource_id text not null,
  period text not null check (period ~ '^\d{4}-\d{2}$'),
  amount_minor_units bigint not null check (amount_minor_units >= 0),
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, provider, resource_id, period),
  foreign key (tenant_id) references platform.tenant (tenant_id)
);

commit;
````

### FILE: `db/migrations/0043_finops.down.sql`
```yaml
block_id: "GO-FINOPS-CORE:db/migrations/0043_finops.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b5c6c45a83180321dc40a23e9f131e1f5d43bb0fd6b2d4cd39e8b173fd1cd572"
variables: []
secrets_allowed: false
```
````sql
begin;

drop table if exists finops.cost_entry;
drop table if exists finops.budget;
drop table if exists finops.resource;
drop schema if exists finops;

commit;
````

### FILE: `db/tests/0043_finops.test.sql`
```yaml
block_id: "GO-FINOPS-CORE:db/tests/0043_finops.test.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a0530e2237db6c06171413df4d5ea885c5a28d064e1ab8b84a385303df19bd4f"
variables: []
secrets_allowed: false
```
````sql
-- 0043_finops.test.sql — verifica tags obligatorios (nada sin atribuir),
-- budget cap (spent <= cap) y cost entry por recurso.
-- Precondición: migración 0001 (platform.tenant) y 0043 aplicadas.

begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name) values
  ('11111111-1111-1111-1111-111111111111', 'tenant-a', 'A', 'Tenant A');

insert into finops.resource
  (tenant_id, provider, resource_id, kind, region, tags, estimated_monthly_minor_units)
values
  ('11111111-1111-1111-1111-111111111111', 'gcp', 'r1', 'cloud_run', 'us-central1',
   '{"tenant":"t","environment":"prod","owner":"platform"}', 100);

-- recurso sin environment/owner → check_violation
do $$
begin
  begin
    insert into finops.resource
      (tenant_id, provider, resource_id, kind, region, tags, estimated_monthly_minor_units)
    values
      ('11111111-1111-1111-1111-111111111111', 'aws', 'r2', 'compute', 'us-east-1',
       '{"tenant":"t"}', 100);
    raise exception 'untagged resource accepted';
  exception when check_violation then
    null; -- expected
  end;
end $$;

-- budget: spent <= cap
insert into finops.budget (tenant_id, period, cap_minor_units, spent_minor_units) values
  ('11111111-1111-1111-1111-111111111111', '2026-09', 100, 60);

do $$
begin
  begin
    update finops.budget set spent_minor_units = 101
     where tenant_id = '11111111-1111-1111-1111-111111111111' and period = '2026-09';
    raise exception 'over-cap accepted';
  exception when check_violation then
    null; -- expected
  end;
end $$;

-- cost entry por recurso
insert into finops.cost_entry (tenant_id, provider, resource_id, period, amount_minor_units) values
  ('11111111-1111-1111-1111-111111111111', 'gcp', 'r1', '2026-09', 40);

do $$
declare n int;
begin
  select count(*) into n from finops.resource
   where tenant_id = '11111111-1111-1111-1111-111111111111';
  if n <> 1 then raise exception 'expected 1 resource, got %', n; end if;
end $$;

rollback;
````


## 6. Configuration surface

Sin variables ni secretos. Los tags obligatorios (`tenant`/`environment`/`owner`) y el cap por tenant se configuran en código; el costo estimado es un dato por recurso.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | core Go | BSD-3-Clause | runtime | https://go.dev |
| PostgreSQL 18.6 | 18.6 | finops.* durable | PostgreSQL License | runtime | https://www.postgresql.org |

## 8. Apply order

1. Componer `PG-TX-FOUNDATION` (migración `0001`).
2. Aplicar migración `0043` sobre PostgreSQL 18.6.
3. Colocar los cinco archivos Go bajo `internal/finops/` y los tres SQL bajo `db/`.
4. Verificar: `psql ... -v ON_ERROR_STOP=1 -f db/tests/0043_finops.test.sql` y `go test ./... -count=1`.
5. Rollback: migración `0043` down y eliminar `internal/finops/`.

## 9. Verification

- `go test ./internal/finops/ -count=1`: 7/7 PASS (tags obligatorios, inventario dedup/costo, rechazo de untagged, audit detecta escapes y leak, budget fail-closed, teardown limpio vs sucio).
- `go test ./... -count=1` (13 paquetes): PASS.
- `go vet ./...`: exit 0.
- PostgreSQL 18.6 real (initdb → up → test → down): `ON_ERROR_STOP=1` exit 0; aserciones DO: untagged rechazado, over-cap rechazado; down deja `finops_ns=true`.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_FINOPS_CORE_2026-09-02_V190.md`.
