# Go Enterprise Query API

## 1. Metadata

```yaml
pack_id: "GO-ENTERPRISE-QUERY-API"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade lecturas paginadas y resource-scoped para portales administrativos, de fábrica y clientes sobre el backend empresarial Go/PostgreSQL."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.4", "GO-SUPPLY-FACTORY-INVENTORY-API 0.16.0", "GO-COMMERCE-PRICING-PAYMENT-API 0.6.1", "GO-FULFILLMENT-SERVICE-FRANCHISE-API 0.4.0", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.19 for the optional connected browser fixture"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/", "https://github.com/jackc/pgx"]
verified_at: "2026-09-11"
```

Este pack separa consultas de comandos y no expone un proxy SQL. Todas las rutas protegidas derivan tenant, sujeto y organizaciones del token verificado. La paginación es keyset y su límite máximo es 100.

## 2. Applicability

Use for bounded portal read models over the Go/PostgreSQL system when administrative, factory and customer views need distinct resource scopes. Reject it as an ad-hoc SQL proxy, analytics export or substitute for a separately owned search/warehouse workload.

## 3. Architecture contract

Queries are read-only adapters with keyset pagination and a hard limit of 100. Tenant, subject and organization scopes originate in verified claims and are repeated in PostgreSQL predicates. Customer service data is reachable only through owned/warrantied resources. Cursor/filter inputs are parsed strictly; invalid or cross-scope reads fail closed. Replica use requires an explicit consistency contract.

## 4. Exact file manifest

```text
CREATE internal/enterprisequery/service.go
CREATE internal/enterprisequery/service_test.go
CREATE internal/platform/postgres/enterprisequery.go
CREATE internal/platform/postgres/enterprisequery_integration_test.go
CREATE internal/platform/httpapi/enterprisequery.go
CREATE internal/platform/httpapi/enterprisequery_test.go
CREATE internal/platform/httpapi/admin_read_browser_test.go
CREATE internal/enterprisequery/factory_unit.go
CREATE internal/platform/httpapi/factory_unit_query.go
CREATE internal/platform/postgres/factory_browser_integration_test.go
CREATE internal/platform/postgres/factory_unit_query.go
```

## 5. Materialization blocks

### FILE: `internal/enterprisequery/service.go`

```yaml
block_id: "GO-ENTERPRISE-QUERY-API:internal-enterprisequery-service-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "f46e63feae5567f593bdaaf68472f54122beb1683f6119665810ec76a5ba6b7a"
variables: []
secrets_allowed: false
```

````go
package enterprisequery

import (
	"context"
	"errors"
	"fmt"
)

var ErrInvalid = errors.New("invalid query")

type Overview struct {
	OrganizationID  string `json:"organization_id"`
	Orders          int64  `json:"orders"`
	OpenLeads       int64  `json:"open_leads"`
	StockAvailable  int64  `json:"stock_available"`
	OpenCases       int64  `json:"open_cases"`
	ActiveShipments int64  `json:"active_shipments"`
}

type Order struct {
	ID              string `json:"id"`
	OrganizationID  string `json:"organization_id"`
	CustomerSubject string `json:"customer_subject"`
	State           string `json:"state"`
	Currency        string `json:"currency"`
	TotalMinorUnits int64  `json:"total_minor_units"`
	Version         int64  `json:"version"`
}

type FactoryUnit struct {
	ID              string `json:"id"`
	OrganizationID  string `json:"organization_id"`
	PurchaseOrderID string `json:"purchase_order_id"`
	VariantID       string `json:"variant_id"`
	SerialNumber    string `json:"serial_number"`
	State           string `json:"state"`
}

type ServiceCase struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	StockUnitID    string `json:"stock_unit_id"`
	State          string `json:"state"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	Version        int64  `json:"version"`
}

type Page[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type Repository interface {
	Overview(context.Context, string, string) (Overview, error)
	Orders(context.Context, string, string, string, int, string) (Page[Order], error)
	FactoryUnits(context.Context, string, string, int, string) (Page[FactoryUnit], error)
	ServiceCases(context.Context, string, string, string, int, string) (Page[ServiceCase], error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func validate(tenant, organization string, limit int) error {
	if tenant == "" || organization == "" || limit < 1 || limit > 100 {
		return fmt.Errorf("%w: scope or limit", ErrInvalid)
	}
	return nil
}

func (s *Service) Overview(ctx context.Context, tenant, organization string) (Overview, error) {
	if err := validate(tenant, organization, 1); err != nil {
		return Overview{}, err
	}
	return s.repository.Overview(ctx, tenant, organization)
}

func (s *Service) Orders(ctx context.Context, tenant, organization, customer string, limit int, after string) (Page[Order], error) {
	if err := validate(tenant, organization, limit); err != nil {
		return Page[Order]{}, err
	}
	return s.repository.Orders(ctx, tenant, organization, customer, limit, after)
}

func (s *Service) FactoryUnits(ctx context.Context, tenant, organization string, limit int, after string) (Page[FactoryUnit], error) {
	if err := validate(tenant, organization, limit); err != nil {
		return Page[FactoryUnit]{}, err
	}
	return s.repository.FactoryUnits(ctx, tenant, organization, limit, after)
}

func (s *Service) ServiceCases(ctx context.Context, tenant, organization, customer string, limit int, after string) (Page[ServiceCase], error) {
	if err := validate(tenant, organization, limit); err != nil {
		return Page[ServiceCase]{}, err
	}
	return s.repository.ServiceCases(ctx, tenant, organization, customer, limit, after)
}
````

### FILE: `internal/enterprisequery/service_test.go`

```yaml
block_id: "GO-ENTERPRISE-QUERY-API:internal-enterprisequery-service_test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "18271bb94394309a58f5273b8f24e0ee777e1e5173c356378c234bd954cde6d6"
variables: []
secrets_allowed: false
```

````go
package enterprisequery

import (
	"context"
	"testing"
)

type fakeRepository struct{ calls int }

func (f *fakeRepository) Overview(context.Context, string, string) (Overview, error) {
	f.calls++
	return Overview{}, nil
}
func (f *fakeRepository) Orders(context.Context, string, string, string, int, string) (Page[Order], error) {
	f.calls++
	return Page[Order]{}, nil
}
func (f *fakeRepository) FactoryUnits(context.Context, string, string, int, string) (Page[FactoryUnit], error) {
	f.calls++
	return Page[FactoryUnit]{}, nil
}
func (f *fakeRepository) ServiceCases(context.Context, string, string, string, int, string) (Page[ServiceCase], error) {
	f.calls++
	return Page[ServiceCase]{}, nil
}

func TestServiceRejectsUnboundedAndMissingScope(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)
	if _, err := service.Orders(context.Background(), "tenant", "org", "", 101, ""); err == nil {
		t.Fatal("unbounded page accepted")
	}
	if _, err := service.Overview(context.Background(), "tenant", ""); err == nil {
		t.Fatal("missing organization accepted")
	}
	if repository.calls != 0 {
		t.Fatal("invalid query reached repository")
	}
}
````

### FILE: `internal/platform/postgres/enterprisequery.go`

```yaml
block_id: "GO-ENTERPRISE-QUERY-API:internal-platform-postgres-enterprisequery-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "25e9b7988f09a485f11aeb613ad0367273c089c68ecd7c9818eead7fe12608d1"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"

	"elite.local/enterprise/internal/enterprisequery"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnterpriseQuery struct{ pool *pgxpool.Pool }

func NewEnterpriseQuery(pool *pgxpool.Pool) *EnterpriseQuery { return &EnterpriseQuery{pool: pool} }

func (r *EnterpriseQuery) Overview(ctx context.Context, tenant, organization string) (enterprisequery.Overview, error) {
	value := enterprisequery.Overview{OrganizationID: organization}
	err := r.pool.QueryRow(ctx, `select
  (select count(*) from sales.customer_order where tenant_id=$1 and organization_id=$2),
  (select count(*) from crm.lead where tenant_id=$1 and organization_id=$2 and lifecycle_state not in ('won','lost')),
  (select count(*) from inventory.stock_unit where tenant_id=$1 and organization_id=$2 and state='available'),
  (select count(*) from service_ops.service_case where tenant_id=$1 and organization_id=$2 and state not in ('closed','cancelled')),
  (select count(*) from logistics.shipment where tenant_id=$1 and (origin_organization_id=$2 or destination_organization_id=$2) and state not in ('delivered','cancelled'))`, tenant, organization).Scan(&value.Orders, &value.OpenLeads, &value.StockAvailable, &value.OpenCases, &value.ActiveShipments)
	return value, err
}

func (r *EnterpriseQuery) Orders(ctx context.Context, tenant, organization, customer string, limit int, after string) (enterprisequery.Page[enterprisequery.Order], error) {
	rows, err := r.pool.Query(ctx, `select order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version
from sales.customer_order where tenant_id=$1 and organization_id=$2 and ($3='' or customer_principal_id=$3) and ($4='' or order_id>$4)
order by order_id limit $5`, tenant, organization, customer, after, limit+1)
	if err != nil {
		return enterprisequery.Page[enterprisequery.Order]{}, err
	}
	defer rows.Close()
	items := make([]enterprisequery.Order, 0, limit)
	for rows.Next() {
		var value enterprisequery.Order
		if err = rows.Scan(&value.ID, &value.OrganizationID, &value.CustomerSubject, &value.State, &value.Currency, &value.TotalMinorUnits, &value.Version); err != nil {
			return enterprisequery.Page[enterprisequery.Order]{}, err
		}
		items = append(items, value)
	}
	if err = rows.Err(); err != nil {
		return enterprisequery.Page[enterprisequery.Order]{}, err
	}
	return orderPage(items, limit), nil
}

func orderPage(items []enterprisequery.Order, limit int) enterprisequery.Page[enterprisequery.Order] {
	page := enterprisequery.Page[enterprisequery.Order]{Items: items}
	if len(items) > limit {
		page.NextCursor = items[limit-1].ID
		page.Items = items[:limit]
	}
	return page
}

func (r *EnterpriseQuery) FactoryUnits(ctx context.Context, tenant, organization string, limit int, after string) (enterprisequery.Page[enterprisequery.FactoryUnit], error) {
	rows, err := r.pool.Query(ctx, `select u.production_unit_id,p.destination_organization_id,u.purchase_order_id,u.variant_id,u.serial_number,u.state
from factory.production_unit u join procurement.purchase_order p on p.tenant_id=u.tenant_id and p.purchase_order_id=u.purchase_order_id
where u.tenant_id=$1 and p.destination_organization_id=$2 and ($3='' or u.production_unit_id>$3)
order by u.production_unit_id limit $4`, tenant, organization, after, limit+1)
	if err != nil {
		return enterprisequery.Page[enterprisequery.FactoryUnit]{}, err
	}
	defer rows.Close()
	items := make([]enterprisequery.FactoryUnit, 0, limit)
	for rows.Next() {
		var value enterprisequery.FactoryUnit
		if err = rows.Scan(&value.ID, &value.OrganizationID, &value.PurchaseOrderID, &value.VariantID, &value.SerialNumber, &value.State); err != nil {
			return enterprisequery.Page[enterprisequery.FactoryUnit]{}, err
		}
		items = append(items, value)
	}
	if err = rows.Err(); err != nil {
		return enterprisequery.Page[enterprisequery.FactoryUnit]{}, err
	}
	page := enterprisequery.Page[enterprisequery.FactoryUnit]{Items: items}
	if len(items) > limit {
		page.NextCursor = items[limit-1].ID
		page.Items = items[:limit]
	}
	return page, nil
}

func (r *EnterpriseQuery) ServiceCases(ctx context.Context, tenant, organization, customer string, limit int, after string) (enterprisequery.Page[enterprisequery.ServiceCase], error) {
	rows, err := r.pool.Query(ctx, `select c.service_case_id,c.organization_id,c.stock_unit_id,c.state,c.severity,c.description,c.version
from service_ops.service_case c where c.tenant_id=$1 and c.organization_id=$2 and ($3='' or exists(select 1 from service_ops.warranty w where w.tenant_id=c.tenant_id and w.stock_unit_id=c.stock_unit_id and w.customer_principal_id=$3)) and ($4='' or c.service_case_id>$4)
order by c.service_case_id limit $5`, tenant, organization, customer, after, limit+1)
	if err != nil {
		return enterprisequery.Page[enterprisequery.ServiceCase]{}, err
	}
	defer rows.Close()
	items := make([]enterprisequery.ServiceCase, 0, limit)
	for rows.Next() {
		var value enterprisequery.ServiceCase
		if err = rows.Scan(&value.ID, &value.OrganizationID, &value.StockUnitID, &value.State, &value.Severity, &value.Description, &value.Version); err != nil {
			return enterprisequery.Page[enterprisequery.ServiceCase]{}, err
		}
		items = append(items, value)
	}
	if err = rows.Err(); err != nil {
		return enterprisequery.Page[enterprisequery.ServiceCase]{}, err
	}
	page := enterprisequery.Page[enterprisequery.ServiceCase]{Items: items}
	if len(items) > limit {
		page.NextCursor = items[limit-1].ID
		page.Items = items[:limit]
	}
	return page, nil
}
````

### FILE: `internal/platform/postgres/enterprisequery_integration_test.go`

```yaml
block_id: "GO-ENTERPRISE-QUERY-API:internal-platform-postgres-enterprisequery_integration_test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "c2712335a1a10e5206dfa0557e9d5385865499756ad7a1f625742cfc82330c8d"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestEnterpriseQueriesEnforceOrganizationCustomerAndCursor(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28901"
	cleanup := func() {
		for _, query := range []string{
			`delete from service_ops.service_case where tenant_id=$1`, `delete from service_ops.warranty where tenant_id=$1`,
			`delete from inventory.stock_unit where tenant_id=$1`, `delete from factory.production_unit where tenant_id=$1`,
			`delete from procurement.purchase_order where tenant_id=$1`, `delete from partner.supplier where tenant_id=$1`,
			`delete from sales.customer_order where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`,
			`delete from catalog.vehicle_model where tenant_id=$1`, `delete from crm.customer_profile where tenant_id=$1`,
			`delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`,
		} {
			_, _ = pool.Exec(ctx, query, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'query-api','Query','Query')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org-a','org-a','A','store'),($1,'org-b','org-b','B','store')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized,status) values($1,'customer-1','Customer 1','customer1@example.test','active'),($1,'customer-2','Customer 2','customer2@example.test','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status) values($1,'supplier','supplier','Supplier','active')`,
		`insert into procurement.purchase_order(tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version) values($1,'po','supplier','org-a','accepted','USD',100,1)`,
		`insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,state) values($1,'unit','po','variant','UNIT-QUERY','planned')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order-a','org-a','customer-1','draft','USD',100,1),($1,'order-b','org-a','customer-2','draft','USD',200,1),($1,'order-c','org-b','customer-1','draft','USD',300,1)`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'stock-a','org-a','variant','STOCK-QUERY-A','VIN-QUERY-A','BATTERY-QUERY-A','available',1,clock_timestamp()),($1,'stock-b','org-b','variant','STOCK-QUERY-B','VIN-QUERY-B','BATTERY-QUERY-B','available',1,clock_timestamp())`,
		`insert into service_ops.warranty(tenant_id,warranty_id,stock_unit_id,customer_principal_id,starts_at,ends_at,terms_version,status) values($1,'warranty','stock-a','customer-1',clock_timestamp(),clock_timestamp()+interval '1 year','v1','active')`,
		`insert into service_ops.service_case(tenant_id,service_case_id,stock_unit_id,organization_id,state,severity,description,version) values($1,'case-a','stock-a','org-a','opened','medium','inspection',1)`,
	}
	for _, query := range fixtures {
		if _, err = pool.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repository := NewEnterpriseQuery(pool)
	overview, err := repository.Overview(ctx, tenant, "org-a")
	if err != nil || overview.Orders != 2 || overview.StockAvailable != 1 || overview.OpenCases != 1 {
		t.Fatalf("overview=%+v err=%v", overview, err)
	}
	first, err := repository.Orders(ctx, tenant, "org-a", "", 1, "")
	if err != nil || len(first.Items) != 1 || first.NextCursor == "" {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	second, err := repository.Orders(ctx, tenant, "org-a", "", 1, first.NextCursor)
	if err != nil || len(second.Items) != 1 || second.Items[0].ID == first.Items[0].ID {
		t.Fatalf("second=%+v err=%v", second, err)
	}
	customerOrders, err := repository.Orders(ctx, tenant, "org-a", "customer-1", 10, "")
	if err != nil || len(customerOrders.Items) != 1 || customerOrders.Items[0].CustomerSubject != "customer-1" {
		t.Fatalf("customer orders=%+v err=%v", customerOrders, err)
	}
	units, err := repository.FactoryUnits(ctx, tenant, "org-a", 10, "")
	if err != nil || len(units.Items) != 1 || units.Items[0].OrganizationID != "org-a" {
		t.Fatalf("units=%+v err=%v", units, err)
	}
	cases, err := repository.ServiceCases(ctx, tenant, "org-a", "customer-1", 10, "")
	if err != nil || len(cases.Items) != 1 {
		t.Fatalf("cases=%+v err=%v", cases, err)
	}
	foreignCases, err := repository.ServiceCases(ctx, tenant, "org-a", "customer-2", 10, "")
	if err != nil || len(foreignCases.Items) != 0 {
		t.Fatalf("foreign cases=%+v err=%v", foreignCases, err)
	}
}
````

### FILE: `internal/platform/httpapi/enterprisequery.go`

```yaml
block_id: "GO-ENTERPRISE-QUERY-API:internal-platform-httpapi-enterprisequery-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "4a0f955754f95e2196cc931a780e45840eb65e999118fd26925ad243d6526dec"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"net/http"
	"strconv"

	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/platform/identity"
)

type EnterpriseQueryModule struct{ Service *enterprisequery.Service }

func (m EnterpriseQueryModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := enterpriseQueryAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("GET /v1/admin/overview", api.adminOverview)
	mux.HandleFunc("GET /v1/admin/orders", api.adminOrders)
	mux.HandleFunc("GET /v1/admin/service-cases", api.adminCases)
	mux.HandleFunc("GET /v1/factory/units", api.factoryUnits)
	mux.HandleFunc("GET /v1/factory/units/{id}", api.factoryUnit)
	mux.HandleFunc("GET /v1/customer/orders", api.customerOrders)
	mux.HandleFunc("GET /v1/customer/service-cases", api.customerCases)
}

type enterpriseQueryAPI struct {
	service  *enterprisequery.Service
	verifier identity.Verifier
}

func (a enterpriseQueryAPI) scope(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, string, bool) {
	principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return principal, "", false
	}
	if !principal.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return principal, "", false
	}
	organization := r.URL.Query().Get("organization_id")
	if !principal.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return principal, "", false
	}
	return principal, organization, true
}

func pageInput(w http.ResponseWriter, r *http.Request) (int, string, bool) {
	limit := 25
	if raw := r.URL.Query().Get("limit"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 1 || value > 100 {
			writeProblem(w, 400, "INVALID_PAGE", "limit must be between 1 and 100")
			return 0, "", false
		}
		limit = value
	}
	return limit, r.URL.Query().Get("after"), true
}

func (a enterpriseQueryAPI) adminOverview(w http.ResponseWriter, r *http.Request) {
	principal, organization, ok := a.scope(w, r, "admin:read")
	if !ok {
		return
	}
	value, err := a.service.Overview(r.Context(), principal.TenantID, organization)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "overview query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a enterpriseQueryAPI) adminOrders(w http.ResponseWriter, r *http.Request) {
	a.orders(w, r, "admin:read", "")
}

func (a enterpriseQueryAPI) customerOrders(w http.ResponseWriter, r *http.Request) {
	a.orders(w, r, "customer:self", "subject")
}

func (a enterpriseQueryAPI) orders(w http.ResponseWriter, r *http.Request, permission, customerMode string) {
	principal, organization, ok := a.scope(w, r, permission)
	if !ok {
		return
	}
	limit, after, ok := pageInput(w, r)
	if !ok {
		return
	}
	customer := ""
	if customerMode != "" {
		customer = principal.Subject
	}
	value, err := a.service.Orders(r.Context(), principal.TenantID, organization, customer, limit, after)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "orders query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a enterpriseQueryAPI) factoryUnits(w http.ResponseWriter, r *http.Request) {
	principal, organization, ok := a.scope(w, r, "factory:read")
	if !ok {
		return
	}
	limit, after, ok := pageInput(w, r)
	if !ok {
		return
	}
	value, err := a.service.FactoryUnits(r.Context(), principal.TenantID, organization, limit, after)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "factory query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a enterpriseQueryAPI) adminCases(w http.ResponseWriter, r *http.Request) {
	a.cases(w, r, "admin:read", "")
}

func (a enterpriseQueryAPI) customerCases(w http.ResponseWriter, r *http.Request) {
	a.cases(w, r, "customer:self", "subject")
}

func (a enterpriseQueryAPI) cases(w http.ResponseWriter, r *http.Request, permission, customerMode string) {
	principal, organization, ok := a.scope(w, r, permission)
	if !ok {
		return
	}
	limit, after, ok := pageInput(w, r)
	if !ok {
		return
	}
	customer := ""
	if customerMode != "" {
		customer = principal.Subject
	}
	value, err := a.service.ServiceCases(r.Context(), principal.TenantID, organization, customer, limit, after)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "service query failed")
		return
	}
	writeJSON(w, 200, value)
}
````

### FILE: `internal/platform/httpapi/enterprisequery_test.go`

```yaml
block_id: "GO-ENTERPRISE-QUERY-API:internal-platform-httpapi-enterprisequery_test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "0f6e18c24737636e62d3bb037e124e7bfea3cc71db661087caa975a938596a6c"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/platform/identity"
)

type queryRepository struct{ customer string }

func (*queryRepository) Overview(context.Context, string, string) (enterprisequery.Overview, error) {
	return enterprisequery.Overview{}, nil
}
func (r *queryRepository) Orders(_ context.Context, _, _, customer string, _ int, _ string) (enterprisequery.Page[enterprisequery.Order], error) {
	r.customer = customer
	return enterprisequery.Page[enterprisequery.Order]{}, nil
}
func (*queryRepository) FactoryUnits(context.Context, string, string, int, string) (enterprisequery.Page[enterprisequery.FactoryUnit], error) {
	return enterprisequery.Page[enterprisequery.FactoryUnit]{}, nil
}
func (*queryRepository) ServiceCases(context.Context, string, string, string, int, string) (enterprisequery.Page[enterprisequery.ServiceCase], error) {
	return enterprisequery.Page[enterprisequery.ServiceCase]{}, nil
}

type queryVerifier struct{}

func (queryVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "customer-1", TenantID: "tenant-1", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"org-a": {}}}, nil
}

func TestCustomerQueryBindsSubjectAndOrganization(t *testing.T) {
	repository := &queryRepository{}
	mux := http.NewServeMux()
	EnterpriseQueryModule{Service: enterprisequery.NewService(repository)}.Register(mux, queryVerifier{})
	request := httptest.NewRequest("GET", "/v1/customer/orders?organization_id=org-a", nil)
	request.Header.Set("Authorization", "Bearer valid")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || repository.customer != "customer-1" {
		t.Fatalf("status=%d customer=%q", response.Code, repository.customer)
	}
	request = httptest.NewRequest("GET", "/v1/customer/orders?organization_id=org-b", nil)
	request.Header.Set("Authorization", "Bearer valid")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("cross-organization query status=%d", response.Code)
	}
}
````

### FILE: `internal/platform/httpapi/admin_read_browser_test.go`
```yaml
block_id: "GO-ENTERPRISE-QUERY-API:PRIVATE-READ:internal-platform-httpapi-admin_read_browser_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read permission repair and exact existing quote-format extraction; see PRIVATE_PORTAL_READS_V382.md"
license: "LicenseRef-Workspace-Owner"
sha256: "e8b8fce98eb9f0fe818eee7ecd8613fcf69e4094894b960f59e945c15dbf9836"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"fmt"
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
	"sync/atomic"
	"testing"
	"time"
)

// AUTHORED read-only browser fixture. Requires the selected journey test issuer;
// it exercises actual RS256/JWKS, HTTP scopes, PostgreSQL queries and Next JWE.
func TestAdministrativeReadBrowserPostgres(t *testing.T) { runPrivateBrowserPostgres(t, false) }

func TestCustomerCancellationBrowserPostgres(t *testing.T) {
 if os.Getenv("ELITE_CUSTOMER_CANCEL_E2E") != "1" { t.Skip("explicit disposable customer cancellation gate not requested") }
 runPrivateBrowserPostgres(t, true)
}

func runPrivateBrowserPostgres(t *testing.T, cancellation bool) {
	if os.Getenv("ELITE_ADMIN_READ_E2E") != "1" {
		t.Skip("explicit disposable admin-read gate not requested")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback elite_confirmation_* database")
	}
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("absolute composed web root required")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'query-api','Query','Query')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org-a','org-a','A','store'),($1,'org-b','org-b','B','store')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized,status) values($1,'customer-1','Customer 1','customer1@example.test','active'),($1,'customer-2','Customer 2','customer2@example.test','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status) values($1,'supplier','supplier','Supplier','active')`,
		`insert into procurement.purchase_order(tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version) values($1,'po','supplier','org-a','accepted','USD',100,1)`,
		`insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,state) values($1,'unit','po','variant','UNIT-QUERY','planned')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order-a','org-a','customer-1','draft','USD',100,1),($1,'order-b','org-a','customer-2','draft','USD',200,1),($1,'order-c','org-b','customer-1','draft','USD',300,1)`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'stock-a','org-a','variant','STOCK-QUERY-A','VIN-QUERY-A','BATTERY-QUERY-A','available',1,clock_timestamp()),($1,'stock-b','org-b','variant','STOCK-QUERY-B','VIN-QUERY-B','BATTERY-QUERY-B','available',1,clock_timestamp())`,
		`insert into service_ops.warranty(tenant_id,warranty_id,stock_unit_id,customer_principal_id,starts_at,ends_at,terms_version,status) values($1,'warranty','stock-a','customer-1',clock_timestamp(),clock_timestamp()+interval '1 year','v1','active')`,
		`insert into service_ops.service_case(tenant_id,service_case_id,stock_unit_id,organization_id,state,severity,description,version) values($1,'case-a','stock-a','org-a','opened','medium','inspection',1)`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload) values($1,'lead-visible','org-a','customer-1','new','fixture','{}'),($1,'lead-foreign-org','org-b','customer-2','new','fixture','{}')`,
	}
	fixtures = append(fixtures,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status) values($1,'book','AR','USD',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode) values($1,'book','variant',100,'inclusive')`,
		`insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version) values($1,'quote-visible','org-a','lead-visible','customer-1','variant','book','USD',100,clock_timestamp()+interval '1 day','issued',1)`,
	)

	fixtures = append(fixtures,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) select $1,'zz-order-'||lpad(i::text,3,'0'),'org-a','customer-1','draft','USD',100,1 from generate_series(1,26) i`,
		`insert into service_ops.service_case(tenant_id,service_case_id,stock_unit_id,organization_id,state,severity,description,version) select $1,'zz-case-'||lpad(i::text,3,'0'),'stock-a','org-a','opened','medium','pagination fixture',1 from generate_series(1,26) i`,
		`insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,vin,battery_serial_number,state) select $1,'zz-unit-'||lpad(i::text,3,'0'),'po','variant','PAGE-UNIT-'||lpad(i::text,3,'0'),'PAGE-VIN-'||i::text,'PAGE-BATTERY-'||i::text,'planned' from generate_series(1,26) i`,
	)

	fixtures=append(fixtures,`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload) select $1,'zz-lead-'||lpad(i::text,3,'0'),'org-a','customer-1','new','fixture','{}' from generate_series(1,26) i`)

	fixtures = append(fixtures,
        `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,appointment_kind,starts_at,state,version) values($1,'time-winter','org-a','lead-visible','customer-1','consultation','2035-01-02T01:30:00Z','requested',1),($1,'time-summer','org-a','lead-visible','customer-1','service','2035-07-02T01:30:00Z','requested',1),($1,'time-foreign','org-b','lead-foreign-org','customer-2','delivery','2035-01-03T01:30:00Z','requested',1)`,
    )

	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}

 cancelFixtures := map[string]map[string]map[string]string{}
 if cancellation {
  for projectIndex, project := range []string{"chromium-desktop","chromium-mobile","firefox-desktop","webkit-desktop"} {
   cancelFixtures[project]=map[string]map[string]string{}
   for phaseIndex, phase := range []string{"normal","lost","malformed","stale"} {
    id:="cc-"+project+"-"+phase
    at:=time.Date(2035,1,10+projectIndex*4+phaseIndex,1,30,0,0,time.UTC).Format(time.RFC3339)
    cancelFixtures[project][phase]=map[string]string{"id":id,"instant":at}
    if _,err:=pool.Exec(ctx,`insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,appointment_kind,starts_at,state,version)values($1,$2,'org-a','lead-visible','customer-1','consultation',$3,'requested',1)`,tenant,id,at);err!=nil{t.Fatal(err)}
   }
  }
 }

	snapshot := func() string {
		var digest string
		q := `select md5(string_agg(payload,'|' order by payload)) from (
   select row_to_json(x)::text payload from sales.customer_order x where tenant_id=$1
   union all select row_to_json(x)::text from service_ops.service_case x where tenant_id=$1
   union all select row_to_json(x)::text from factory.production_unit x where tenant_id=$1
   union all select row_to_json(x)::text from crm.lead x where tenant_id=$1
   union all select row_to_json(x)::text from inventory.stock_unit x where tenant_id=$1
   union all select row_to_json(x)::text from sales.quotation x where tenant_id=$1
   union all select row_to_json(x)::text from crm.appointment x where tenant_id=$1
  )q`
        if cancellation {q=strings.Replace(q,"select row_to_json(x)::text from crm.appointment x","select (to_jsonb(x)-'state'-'version'-'updated_at')::text from crm.appointment x",1)}
		if err := pool.QueryRow(ctx, q, tenant).Scan(&digest); err != nil {
			t.Fatal(err)
		}
		return digest
	}
	before := snapshot()
	verifier, token := confirmationTestIssuer(t)
	identities := map[string]any{}
	tokens := map[string]string{}
	add := func(name string, permissions, organizations []string, tenantID string) {
		subject := name
		if name == "customer" {
			subject = "customer-1"
		}
		bearer := token(subject, tenantID, permissions, organizations)
		tokens[name] = bearer
		identities[name] = map[string]any{"subject": subject, "tenantId": tenantID, "permissions": permissions, "organizations": organizations, "accessToken": bearer}
	}
	for name, permissions := range map[string][]string{"admin": {"admin:read"}, "admin-lead": {"admin:read", "lead:read"}, "admin-case": {"admin:read", "Lead:Read"}, "wildcard": {"*"}, "lead": {"lead:read"}, "factory": {"factory:read"}, "forged-admin": {"admin:read"}} {
		add(name, permissions, []string{"org-a"}, tenant)
	}
	add("customer", []string{"customer:self"}, []string{"org-a"}, tenant)
	add("other-org", []string{"admin:read", "lead:read"}, []string{"org-b"}, tenant)
	add("other-tenant", []string{"admin:read", "lead:read"}, []string{"org-a"}, randomid.Generator{}.New())
	mux := http.NewServeMux()
	EnterpriseQueryModule{Service: enterprisequery.NewService(postgres.NewEnterpriseQuery(pool))}.Register(mux, verifier)
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	var writes, adminLeadReads, leadAdminReads atomic.Int64
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
            writes.Add(1)
            if !cancellation || r.Method!=http.MethodPost || !strings.HasPrefix(r.URL.Path,"/v1/customer/appointments/") || !strings.HasSuffix(r.URL.Path,"/cancel") {w.WriteHeader(405);return}
        }
		if r.Header.Get("Authorization") == "Bearer "+tokens["admin"] && r.URL.Path == "/v1/franchise/leads" {
			adminLeadReads.Add(1)
		}
		if r.Header.Get("Authorization") == "Bearer "+tokens["lead"] && strings.HasPrefix(r.URL.Path, "/v1/admin/") {
			leadAdminReads.Add(1)
		}
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()
	// Independent real HTTP negatives before browser counters are recorded.
	for _, probe := range []struct {
		path, bearer string
		status       int
	}{
		{"/v1/admin/overview?organization_id=org-a", "invalid", 401},
		{"/v1/admin/overview?organization_id=org-a", tokens["lead"], 403},
		{"/v1/franchise/leads?organization_id=org-a", tokens["forged-admin"], 403},
		{"/v1/admin/overview?organization_id=org-b", tokens["admin"], 403},
		{"/v1/factory/units?organization_id=org-a", tokens["admin"], 403},
	} {
		req, _ := http.NewRequestWithContext(ctx, "GET", api.URL+probe.path, nil)
		req.Header.Set("Authorization", "Bearer "+probe.bearer)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != probe.status {
			t.Fatalf("probe %s status=%d want=%d", probe.path, res.StatusCode, probe.status)
		}
	}
	baselineLeadAdmin := leadAdminReads.Load()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	env := []string{}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if name == "TEST_DATABASE_URL" || name == "DATABASE_URL" || strings.HasPrefix(name, "ELITE_") || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" || name == "APP_BASE_URL" {
			continue
		}
		env = append(env, item)
	}
	encoded, err := json.Marshal(identities)
	if err != nil {
		t.Fatal(err)
	}
	env = append(env, "ELITE_APPOINTMENT_TIME_EXPECTATIONS="+os.Getenv("ELITE_APPOINTMENT_TIME_EXPECTATIONS"), "ELITE_ADMIN_READ_E2E=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_ADMIN_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-"+randomid.Generator{}.New()+randomid.Generator{}.New(), "APP_BASE_URL="+edge.URL, "ELITE_BASE_URL="+edge.URL)
    if cancellation {fixtureJSON,err:=json.Marshal(cancelFixtures);if err!=nil{t.Fatal(err)};env=append(env,"ELITE_CUSTOMER_CANCEL_E2E=1","ELITE_CUSTOMER_CANCEL_FIXTURES="+string(fixtureJSON))}

	artifacts, err := os.MkdirTemp(web, "admin-read-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	nextlog, err := os.Create(filepath.Join(artifacts, "next.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer nextlog.Close()
	server := exec.CommandContext(ctx, "node", filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env = web, env
	server.Stdout, server.Stderr = nextlog, nextlog
	if err = server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = server.Process.Kill(); _ = server.Wait() }()
	ready := false
	client := &http.Client{Timeout: time.Second}
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		res, err := client.Get(target.String() + "/icon.svg")
		if err == nil {
			res.Body.Close()
			if res.StatusCode == 200 {
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
		t.Fatal("web did not start")
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
    selectedTest:="admin reads preserve independent grants";if cancellation {selectedTest="customer cancels with durable result recovery"}
	cmd := exec.CommandContext(ctx, "node", filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/role-workspace.spec.mjs", "-g", selectedTest, "--timeout=60000", "--output="+filepath.Join(artifacts, "browsers"))
	cmd.Dir, cmd.Env = gate, env
	output, err := cmd.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil {
		t.Fatalf("browser: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "4 passed") {
		t.Fatalf("missing four browser successes: %s", output)
	}

 if cancellation {
  var appointments,transitions,events,unchanged int
  checks:=[]struct{q string;dst *int;want int}{
   {`select count(*) from crm.appointment where tenant_id=$1 and appointment_id like 'cc-%' and state='cancelled' and version=2`,&appointments,16},
   {`select count(*) from crm.appointment_transition where tenant_id=$1 and appointment_id like 'cc-%' and from_state='requested' and to_state='cancelled' and actor_subject='customer-1' and reason_code='customer-request'`,&transitions,16},
   {`select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_type='appointment' and aggregate_id like 'cc-%' and aggregate_version=2 and event_type='appointment.cancelled' and payload->>'actor_subject'='customer-1' and payload->>'channel'='customer' and payload->>'reason_code'='customer-request'`,&events,16},
   {`select count(*) from crm.appointment where tenant_id=$1 and appointment_id not like 'cc-%' and state='requested' and version=1`,&unchanged,3},
  }
  for _,check:=range checks {if err:=pool.QueryRow(ctx,check.q,tenant).Scan(check.dst);err!=nil || *check.dst!=check.want {t.Fatalf("cancellation durable check got=%d want=%d err=%v",*check.dst,check.want,err)}}
  var allTransitions,allEvents int
  if err:=pool.QueryRow(ctx,`select (select count(*) from crm.appointment_transition where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1)`,tenant).Scan(&allTransitions,&allEvents);err!=nil || allTransitions!=16 || allEvents!=16 {t.Fatalf("unexpected durable effects transitions=%d events=%d err=%v",allTransitions,allEvents,err)}
  if writes.Load()!=28 || snapshot()!=before {t.Fatalf("cancellation writes=%d or unrelated/immutable data changed",writes.Load())}
  t.Logf("CUSTOMER_CANCELLATION_BROWSER_PASS browsers=4 cancelled=16 transitions=16 outbox=16 writes=28 immutable_snapshot=%s artifacts=%s",before,artifacts)
  return
 }

	if writes.Load() != 0 || adminLeadReads.Load() != 0 || leadAdminReads.Load() != baselineLeadAdmin {
		t.Fatalf("unexpected requests writes=%d admin_leads=%d lead_admin=%d", writes.Load(), adminLeadReads.Load(), leadAdminReads.Load()-baselineLeadAdmin)
	}
	if snapshot() != before {
		t.Fatal("read-only browser modified durable records")
	}
	t.Logf("ADMIN_READ_BROWSER_POSTGRES_PASS browsers=4 http_negatives=5 admin_only_lead_queries=0 non_admin_queries=0 writes=0 paging_records=81 admin_paging_records=82 appointment_times=2 snapshot_tables=7 durable_snapshot=%s artifacts=%s", before, artifacts)
	fmt.Fprintln(os.Stdout, "ADMIN_READ_OWNED_NEXT_STOP_ON_RETURN")
}
````

### FILE: `internal/enterprisequery/factory_unit.go`

```yaml
block_id: "FACTORY-DELTA-GO_ENTERPRISE_QUERY_API:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b7120e937c46fcd408d6dd1de0615764b731456c7e456dc0167a9bd43b4f42b6"
variables: []
secrets_allowed: false
```

````go
package enterprisequery

// AUTHORED read port extension over the existing factory projection.
import (
	"context"
	"errors"
)

var ErrFactoryUnitNotFound = errors.New("factory unit not found")

type FactoryUnitReader interface {
	FactoryUnit(context.Context, string, string, string) (FactoryUnit, error)
}

func (s *Service) FactoryUnit(ctx context.Context, tenant, organization, id string) (FactoryUnit, error) {
	if validate(tenant, organization, 1) != nil || id == "" || len(id) > 200 {
		return FactoryUnit{}, ErrInvalid
	}
	reader, ok := s.repository.(FactoryUnitReader)
	if !ok {
		return FactoryUnit{}, ErrInvalid
	}
	return reader.FactoryUnit(ctx, tenant, organization, id)
}
````

### FILE: `internal/platform/httpapi/factory_unit_query.go`

```yaml
block_id: "FACTORY-DELTA-GO_ENTERPRISE_QUERY_API:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "074278e5faef009a037c1f729698d1792674164c99b16161f1a4b24969353de7"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"elite.local/enterprise/internal/enterprisequery"
	"errors"
	"net/http"
)

// AUTHORED exact query; authorization is identical to the existing factory list.
func (a enterpriseQueryAPI) factoryUnit(w http.ResponseWriter, r *http.Request) {
	principal, organization, ok := a.scope(w, r, "factory:read")
	if !ok {
		return
	}
	value, err := a.service.FactoryUnit(r.Context(), principal.TenantID, organization, r.PathValue("id"))
	if errors.Is(err, enterprisequery.ErrFactoryUnitNotFound) {
		writeProblem(w, 404, "NOT_FOUND", "factory unit not found")
		return
	}
	if errors.Is(err, enterprisequery.ErrInvalid) {
		writeProblem(w, 400, "INVALID_QUERY", "factory unit query is invalid")
		return
	}
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "factory unit query failed")
		return
	}
	writeJSON(w, 200, value)
}
````

### FILE: `internal/platform/postgres/factory_browser_integration_test.go`

```yaml
block_id: "FACTORY-DELTA-GO_ENTERPRISE_QUERY_API:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "382770a2782237428c3447d057b97083426f1094ee3e5bc69e41864510d58914"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED local composition fixture over the existing operations owner.
import (
	"context"
	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
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

func TestFactoryBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_FACTORY_BROWSER") != "1" {
		t.Skip("explicit local factory browser fixture required")
	}
	rawDB := os.Getenv("FACTORY_BROWSER_DB_URL")
	parsed, err := url.Parse(rawDB)
	if err != nil || parsed.Hostname() != "127.0.0.1" || !strings.HasPrefix(parsed.Path, "/elite_payment_connected_") {
		t.Fatal("isolated loopback fixture database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool, err := pgxpool.New(ctx, rawDB)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	ids := randomid.Generator{}
	tenant := ids.New()
	otherTenant := ids.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'factory-browser','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Store','warehouse'),($1,'other','other','Other','warehouse')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status)values($1,'supplier','supplier','Supplier','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := db.NewOperations(pool)
	service := operations.NewService(repo, ids)
	if err = repo.CreatePurchaseOrder(ctx, tenant, ids.New(), operations.PurchaseOrder{ID: "po", SupplierID: "supplier", DestinationOrganizationID: "store", State: "draft", Currency: "USD", TotalMinorUnits: 100, Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err = repo.CreateProductionUnit(ctx, tenant, ids.New(), operations.ProductionUnit{ID: "unit", OrganizationID: "store", PurchaseOrderID: "po", VariantID: "variant", SerialNumber: "SERIAL-FACTORY-FIXTURE", State: "planned"}); err != nil {
		t.Fatal(err)
	}
	verifier, token := handoverBrowserIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"operator", "reader", "foreign-org", "foreign-tenant"} {
		permissions := []string{"factory:read", "factory:write"}
		orgs := []string{"store"}
		tenantID := tenant
		if name == "reader" {
			permissions = []string{"factory:read"}
		}
		if name == "foreign-org" {
			orgs = []string{"other"}
		}
		if name == "foreign-tenant" {
			tenantID = otherTenant
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenantID, "permissions": permissions, "organizations": orgs, "accessToken": token(name, tenantID, permissions, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.OperationsModule{Service: service}.Register(mux, verifier)
	httpapi.EnterpriseQueryModule{Service: enterprisequery.NewService(db.NewEnterpriseQuery(pool))}.Register(mux, verifier)
	controlToken := ids.New()
	mux.HandleFunc("POST /__fixture/factory-competing-advance", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Fixture-Token") != controlToken {
			w.WriteHeader(403)
			return
		}
		if err := service.TransitionProductionUnit(r.Context(), tenant, "store", "unit", "assembly", "quality"); err != nil {
			w.WriteHeader(409)
			return
		}
		w.WriteHeader(200)
	})
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && !strings.HasPrefix(r.URL.Path, "/__fixture/") {
			mu.Lock()
			counts[r.URL.Path]++
			mu.Unlock()
		}
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("absolute web root required")
	}
	env := []string{}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if strings.HasPrefix(name, "ELITE_") || name == "DATABASE_URL" || name == "TEST_DATABASE_URL" || name == "PAYMENT_CONNECTED_DB_URL" || name == "FACTORY_BROWSER_DB_URL" || name == "APP_BASE_URL" || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" {
			continue
		}
		env = append(env, item)
	}
	encoded, err := json.Marshal(identities)
	if err != nil {
		t.Fatal(err)
	}
	env = append(env, "ELITE_FACTORY_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_FACTORY_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-"+ids.New()+ids.New(), "ELITE_FACTORY_CONTROL="+api.URL+"/__fixture/factory-competing-advance", "ELITE_FACTORY_CONTROL_TOKEN="+controlToken)
	artifacts, err := os.MkdirTemp(web, "factory-browser-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	log, err := os.Create(filepath.Join(artifacts, "next.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	node := os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(node) {
		t.Fatal("exact Node executable required")
	}
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env = web, env
	server.Stdout, server.Stderr = log, log
	if err = server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { server.Process.Kill(); server.Wait() }()
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		resp, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
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
		t.Fatal("Next did not start")
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/factory-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("browser %v\n%s\nartifacts=%s", err, output, artifacts)
	}

	var state string
	var assembly, quality, stock int
	err = pool.QueryRow(ctx, `select (select state from factory.production_unit where tenant_id=$1 and production_unit_id='unit'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='production-unit.assembly'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='production-unit.quality'),(select count(*) from inventory.stock_unit where tenant_id=$1)`, tenant).Scan(&state, &assembly, &quality, &stock)
	if err != nil || state != "quality" || assembly != 1 || quality != 1 || stock != 0 {
		t.Fatal("durable outcomes", state, assembly, quality, stock, err)
	}
	mu.Lock()
	countJSON, _ := json.Marshal(counts)
	unitPosts := counts["/v1/factory/units/unit/transitions"]
	mu.Unlock()
	if unitPosts != 1 {
		t.Fatal("unexpected backend mutation attempts", unitPosts)
	}
	if err = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countJSON, 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("FACTORY_BROWSER_POSTGRES_PASS actual_BFF_SQL=true JWE_RS256_JWKS=true assembly_events=1 competing_quality_events=1 backend_transition_posts=1 response_loss_observation_only=true cross_org_tenant_role_stale_rejected=true stock_mutation=false artifacts=%s", artifacts)
}
````

### FILE: `internal/platform/postgres/factory_unit_query.go`

```yaml
block_id: "FACTORY-DELTA-GO_ENTERPRISE_QUERY_API:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d17e3cf3d679906a5725ad5aec5b3f2b7d92dc5bdc456bd61a34ccbb3032695b"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED exact read of the existing FactoryUnits projection; no state mutation.
import (
	"context"
	"elite.local/enterprise/internal/enterprisequery"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (r *EnterpriseQuery) FactoryUnit(ctx context.Context, tenant, organization, id string) (enterprisequery.FactoryUnit, error) {
	var v enterprisequery.FactoryUnit
	err := r.pool.QueryRow(ctx, `select u.production_unit_id,p.destination_organization_id,u.purchase_order_id,u.variant_id,u.serial_number,u.state
 from factory.production_unit u join procurement.purchase_order p on p.tenant_id=u.tenant_id and p.purchase_order_id=u.purchase_order_id
 where u.tenant_id=$1 and p.destination_organization_id=$2 and u.production_unit_id=$3`, tenant, organization, id).Scan(&v.ID, &v.OrganizationID, &v.PurchaseOrderID, &v.VariantID, &v.SerialNumber, &v.State)
	if errors.Is(err, pgx.ErrNoRows) {
		return enterprisequery.FactoryUnit{}, enterprisequery.ErrFactoryUnitNotFound
	}
	return v, err
}
````

## 6. Configuration surface

No new environment variables are introduced. The module receives the composed repository and verifier. Page size is client-selectable only within `1..100`; cursor, organization and filters are validated. Field visibility, export policy and any replica selection require project-owned configuration and tests.

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | `1.26.7` | query service/HTTP/tests | BSD-3-Clause | build/runtime | `go.dev` |
| PostgreSQL | `18.6` verified baseline | scoped reads | PostgreSQL | runtime/test | `postgresql.org` |
| `pgx` | `5.10.0` | query adapter | MIT | build/runtime | `github.com/jackc/pgx` |

## 8. Apply order

Compose after migrations and command-side modules so referenced tables and authorization types exist; then wire routes in the application. Run query isolation and pagination tests against PostgreSQL. Existing projects must align sort/cursor semantics and indexes. Rollback removes routes/binary wiring without deleting source data.

## 9. Verification

Componer con migrations 0001/0002/0003 y la aplicación Go. Ejecutar `gofmt -d`, unit/HTTP/PostgreSQL integration, `go vet` y ambos builds. Las pruebas deben demostrar: organización ajena rechazada, `customer:self` ligado al subject, lecturas de servicio ligadas a warranty, fábrica ligada a la organización destino de la compra y cursor sin duplicación.

Condiciones de producción: índices y planes `EXPLAIN` con cardinalidad representativa, filtros/ordenamientos definidos por producto, exportación/analytics separada, read replicas sólo con consistencia declarada, auditoría de lecturas sensibles, redacción, retención y autorización de campos por rol.

## 10. Reconstruction evidence

Clean materialization, PostgreSQL scope tests and build gates are recorded in `reconstruction_evidence/GO_ENTERPRISE_QUERY_API_2026-08-24_V1.md`; final integrated evidence rechecks version 0.1.1.

V382: exact selected-composition compatibility refreshed. Administrative lead reads are conditional on the existing lead grant; shared quote-format code corrects admin/customer minor-unit labels without repricing. TestAdministrativeReadBrowserPostgres uses real Go/JWKS/PostgreSQL and four browsers with synthetic identities; all production Go/SQL remains unchanged. The role implementation and quote acceptance/date/state logic retain exact source parity. See docs/private-portal-reads.md and reconstruction_evidence/PRIVATE_PORTAL_READS_V382.md. No full business, IdP, security, release or target admission.

V387: customer/factory cursor traversal connects three existing lists beyond25rows.81durable records per browser, four projects, reload/independent first-page recovery, zero writes and unchanged snapshots verified. Source is AUTHORED, business/API/SQL/dependencies unchanged. Supporting TEST02 progress, no whole-control or native security/release promotion. See reconstruction_evidence/PRIVATE_PORTAL_PAGINATION_V387.md.

V388: admin orders/cases/leads now traverse existing scoped keyset APIs using the shared navigator.276web tests and four browser projects exercise all six lists,163list records, no omissions/duplicates per list, permission removal, reload/first-page recovery and unchanged durable snapshots. All code remains AUTHORED/CONDITIONED; no full TEST02, native security or release promotion. See reconstruction_evidence/ADMIN_PORTAL_PAGINATION_V388.md.

V389: customer appointment times share trusted business locale/zone across SSR and interactive management, preserving original instants and cancellation behavior.282web tests plus eight connected browser/configuration runs, two appointment instants per view and seven-table unchanged snapshots. AUTHORED/CONDITIONED; no whole TEST02, security or release promotion. See reconstruction_evidence/CUSTOMER_APPOINTMENT_TIMEZONE_V389.md.

V390: customer cancellation receipt binding, synchronous single-submit and explicit GET recovery close the existing reference journey.294web tests;4actual cancellation browser projects,16unique durable cancellations/audits/outbox and28API writes;8read/time regression runs remain read-only. AUTHORED/CONDITIONED, no whole TEST02/security/release promotion. See reconstruction_evidence/CUSTOMER_CANCELLATION_RECOVERY_V390.md.

V402 composed delta: Connected factory tracking browser/BFF/Go/PostgreSQL gate; exact scoped unit read and planned-to-assembly call to existing owner; recovery observes current state without attributing an uncertain effect. AUTHORED glue; existing domain and fixed upstreams unchanged.

V402 factory browser evidence: reconstruction_evidence/FACTORY_BROWSER_V402.md.18BFF tests, Next production build with TypeScript, Go vet and real browser/API/PG PASS; current-state recovery after another operator advances. No idempotent receipt claim or domain/writer change.
