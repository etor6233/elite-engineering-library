# Connected role metrics proof

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-ROLE-METRICS-PROOF"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "T2804 role metrics use existing domain read models with exact strings, organization/customer/factory/program permissions, NPS minimum/retention and no zero on unavailable. FAIL868 converted lead and FAIL457 bounded generic body corrected. ROLE_METRICS_RELEASE_V402.md/json; no new dependency or corporate authorship."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["V402 full reference source owners"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-11"
```

## 2. Applicability

Local reference proof of existing query, stored value, feedback, supply and role UI owners; no production certification.

## 3. Architecture contract

Read-only report semantics; currencies/status/programs remain separate, exact strings, profile-bound sources, private no-store, unavailable distinct from empty. Existing guards retained.

## 4. Exact file manifest

```text
CREATE config/role.metrics.fixture.json
CREATE docs/ROLE_METRICS_REFERENCE.md
CREATE internal/platform/postgres/role_metrics_browser_integration_test.go
CREATE internal/platform/postgres/role_metrics_connected_integration_test.go
CREATE microsoft_playwright_browser_gate/tests/role-metrics-connected.spec.mjs
```

## 5. Materialization blocks

### FILE: `config/role.metrics.fixture.json`

```yaml
block_id: "GO-CONNECTED-ROLE-METRICS-PROOF:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6882d401e11b953af100fd1d3aa9eb4f783fee15c5cb2257ee57fc938bd87716"
variables: []
secrets_allowed: false
```

````json
{
  "schemaVersion": "1.0.0",
  "business": {
    "id": "electric-mobility-network",
    "name": "Electric Mobility Network",
    "defaultLocale": "es-AR",
    "defaultMarket": "AR",
    "supportEmail": "soporte@example.invalid"
  },
  "markets": [
    {
      "code": "AR",
      "name": "Argentina",
      "currency": "ARS",
      "locales": [
        "es-AR"
      ],
      "timeZone": "America/Argentina/Buenos_Aires",
      "taxMode": "external"
    }
  ],
  "organizationTypes": [
    {
      "id": "hq",
      "label": "Casa central",
      "allowedParents": []
    },
    {
      "id": "franchise",
      "label": "Franquicia",
      "allowedParents": [
        "hq"
      ]
    },
    {
      "id": "branch",
      "label": "Sucursal",
      "allowedParents": [
        "franchise",
        "hq"
      ]
    },
    {
      "id": "factory",
      "label": "Fábrica",
      "allowedParents": [
        "hq"
      ]
    },
    {
      "id": "supplier",
      "label": "Proveedor",
      "allowedParents": [
        "hq"
      ]
    }
  ],
  "roles": [
    {
      "id": "hq_admin",
      "label": "Administrador central",
      "permissions": [
        "*"
      ]
    },
    {
      "id": "branch_manager",
      "label": "Responsable de sucursal",
      "permissions": [
        "admin:read",
        "lead:read",
        "lead:assign",
        "lead:update",
        "quote:write",
        "catalog:write",
        "pricing:write",
        "order:create",
        "order:write",
        "inventory:allocate",
        "payment:create",
        "payment:write",
        "service:write",
        "communication:write"
      ]
    },
    {
      "id": "sales",
      "label": "Ventas",
      "permissions": [
        "admin:read",
        "lead:read",
        "lead:assign",
        "lead:update",
        "quote:write",
        "catalog:write",
        "order:create",
        "order:write",
        "payment:create"
      ]
    },
    {
      "id": "factory_operator",
      "label": "Operador de fábrica",
      "permissions": [
        "factory:read",
        "procurement:write",
        "factory:write",
        "inventory:write",
        "logistics:write"
      ]
    },
    {
      "id": "customer",
      "label": "Cliente",
      "permissions": [
        "customer:self"
      ]
    }
  ],
  "modules": {
    "catalog": {
      "enabled": true
    },
    "crm": {
      "enabled": true
    },
    "procurement": {
      "enabled": true
    },
    "inventory": {
      "enabled": true
    },
    "orders": {
      "enabled": true
    },
    "payments": {
      "enabled": true
    },
    "fulfillment": {
      "enabled": true
    },
    "service": {
      "enabled": true
    },
    "documents": {
      "enabled": true
    },
    "integrations": {
      "enabled": true
    }
  },
  "workflows": {
    "lead": {
      "initial": "new",
      "states": [
        "new",
        "contacted",
        "qualified",
        "converted",
        "lost"
      ],
      "transitions": [
        {
          "from": "new",
          "to": "contacted",
          "permission": "lead:update"
        },
        {
          "from": "contacted",
          "to": "qualified",
          "permission": "lead:update"
        },
        {
          "from": "qualified",
          "to": "converted",
          "permission": "lead:update"
        },
        {
          "from": "contacted",
          "to": "lost",
          "permission": "lead:update"
        },
        {
          "from": "qualified",
          "to": "lost",
          "permission": "lead:update"
        }
      ]
    },
    "order": {
      "initial": "draft",
      "states": [
        "draft",
        "placed",
        "confirmed",
        "paid",
        "allocated",
        "delivered",
        "cancelled"
      ],
      "transitions": [
        {
          "from": "draft",
          "to": "placed",
          "permission": "order:create"
        },
        {
          "from": "placed",
          "to": "confirmed",
          "permission": "order:transition"
        },
        {
          "from": "confirmed",
          "to": "paid",
          "permission": "payment:reconcile"
        },
        {
          "from": "paid",
          "to": "allocated",
          "permission": "inventory:reserve"
        },
        {
          "from": "allocated",
          "to": "delivered",
          "permission": "order:transition"
        },
        {
          "from": "draft",
          "to": "cancelled",
          "permission": "order:transition"
        },
        {
          "from": "placed",
          "to": "cancelled",
          "permission": "order:transition"
        }
      ]
    }
  },
  "customFields": {
    "lead": [
      {
        "id": "preferred_vehicle_use",
        "label": "Uso principal",
        "type": "select",
        "required": false,
        "options": [
          "urban",
          "delivery",
          "recreation",
          "fleet"
        ]
      }
    ],
    "catalog_model": [
      {
        "id": "estimated_range_km",
        "label": "Autonomía estimada (km)",
        "type": "number",
        "required": true
      }
    ]
  },
  "integrations": [
    {
      "id": "mercado_pago",
      "provider": "mercado_pago",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "payments"
      ],
      "credentialRefEnv": "MERCADO_PAGO_CREDENTIAL_REF"
    },
    {
      "id": "amazon_sp_api",
      "provider": "amazon_sp_api",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "catalog",
        "orders",
        "fulfillment"
      ],
      "credentialRefEnv": "AMAZON_SP_API_CREDENTIAL_REF"
    },
    {
      "id": "mercado_libre",
      "provider": "mercado_libre",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "catalog",
        "orders",
        "fulfillment"
      ],
      "credentialRefEnv": "MERCADO_LIBRE_CREDENTIAL_REF"
    },
    {
      "id": "google_ads",
      "provider": "google_ads",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "ads",
        "conversions"
      ],
      "credentialRefEnv": "GOOGLE_ADS_CREDENTIAL_REF"
    },
    {
      "id": "meta_ads",
      "provider": "meta_ads",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "ads",
        "conversions"
      ],
      "credentialRefEnv": "META_ADS_CREDENTIAL_REF"
    }
  ],
  "features": {
    "public_catalog": true,
    "lead_capture": true,
    "customer_portal": true,
    "factory_portal": true,
    "vehicle_telemetry": false,
    "training_portal": false,
    "catalog_editor": false,
    "supply_portal": false,
    "warranty_portal": false,
    "network_portal": true,
    "help_cms": true,
    "role_workspace": true,
    "customer_surveys": true
  }
}
````

### FILE: `docs/ROLE_METRICS_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-ROLE-METRICS-PROOF:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "401db5b17c969763bc4d4e4c55242127325baaf2454699416bde6c7fd8f82c84"
variables: []
secrets_allowed: false
```

````markdown
# Indicadores conectados por rol

Infraestructura de biblioteca, local/fixtures; AUTHORED read-model/transport/UI
glue sobre los owners admitidos, sin nueva dependencia ni autoría corporativa.
GO-DASHBOARDS-CORE permanece alternativa aislada, no se promueve por estos tests.

El host electromobility integra RoleMetrics en EnterpriseQuery. El panel
/dashboard requiere features.role_workspace=true y sesión. Sus accesos a las
operaciones siguen visibleSections del owner por rol. El BFF sólo recibe kind,
organization_id y, para puntos/NPS, identificador de programa/encuesta.
Nunca recibe tenant, actor, SQL ni valores calculados por el cliente.

| Familia | Fuente | Permiso / alcance |
|---|---|---|
| orders, stock, cases, shipments | sales.customer_order, inventory.stock_unit, service_ops.service_case, logistics.shipment | admin:read + organización; envío origen o destino |
| own-orders, own-cases, own-appointments | pedido, caso vinculado a warranty propia, crm.appointment | customer:self + organización + subject |
| leads / appointments | crm.lead / crm.appointment | lead:read / appointment:manage + organización |
| factory-destination | factory.production_unit + procurement.purchase_order | factory:read + destino |
| factory-owned | factory.production_unit + serial_supply_plan | supply:factory-read + fábrica |
| supply | purchase_order + serial_supply_plan | supply:read + destino; excluye compras sin plan |
| stored-value | immutable operation + entry + account | stored_value:read + tenant/organización/programa del perfil activo |
| survey | Summary original de crm.survey_definition/response | surveys:read + organización; mínimo y retención originales |

Las métricas operativas son registros actuales por estado, no ventas mensuales
ni ingresos cobrados, utilidad o contabilidad. Importes se suman como numeric y
se transmiten como strings exactos por moneda y estado, nunca como JS number.
La suma supera int64 sin redondeo (wire hasta38dígitos); un resultado fuera de
la capacidad del contrato se declara no disponible. Sin conversión de monedas.
Los puntos se separan por programa/moneda y operación: issue/accrue/redeem/reverse.
Sólo se cuentan operaciones con entradas de ledger confirmadas, nunca propuestas
pendientes. No se convierten puntos a dinero ni se suman programas distintos.
NPS preserva el cálculo y mínimo existentes; cero respuestas no es NPS0.

Lecturas de operación/puntos usan transacción read-only repeatable-read,
statement_timeout3s; endpoint4s; NPS usa el SELECT original con deadline3s.
Cada consulta tiene su propio observed_at, no hay snapshot global entre familias.
Máximo100grupos operativos; el grupo101 produce409 y no truncamiento silencioso.
Wire puntos≤4operaciones. Autorización/servicio/perfil no disponible jamás se
presentan como cero. Cambiar selector borra resultados previos; generaciones
evitan que una respuesta tardía restaure datos de otra selección.
No-store en API/BFF. El BFF retenido usa transporte con deadline y bytes acotados.

FAIL868 Overview excluía won, estado ajeno al dominio. El cambio mínimo excluye
converted y lost; regresión RED→PASS preservada. FAIL457 del BFF genérico:
sesión/permisos antes del body, límite65536bytes, UTF8 fatal, deadline2.5s,
cancelación y permisos originales por acción. No nuevos grants ni cambios
de transición de negocio. No prueba de seguridad de producción.

Fixture config/role.metrics.fixture.json habilita el panel/encuestas para el
ensayo local; los flags de usuario siguen explícitos. Browser Chromium real
usa JWE+JWT/JWKS sintéticos, Go/PG y BFF; dos respuestas guardadas actualizan
NPS, distintas monedas/estado/importes exactos, roles, vacío y unavailable.
TestRoleMetricsConnectedScopeAndPrecision, StoredValueAndSurveyOwners,
FactorySourceScope y Browser prueban el delta; no se repiten recorridos previos.
FuzzMetricBoundary tiene seis semillas y presupuesto finito2s, no reemplaza SAST.
ROLE_METRICS_RELEASE_V402.md/json registra hashes/evidencia/reconstrucciones.
T2804 continúa con i18n privado. No READY global ni producción live.
````

### FILE: `internal/platform/postgres/role_metrics_browser_integration_test.go`

```yaml
block_id: "GO-CONNECTED-ROLE-METRICS-PROOF:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c79ebf032e1318890150d29cc78b476875d56555ca80002fc2d4b747deec5206"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED browser composition of scoped read models and existing survey writer.
import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/enterprisequery"
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

func TestRoleMetricsBrowser(t *testing.T) {
	if os.Getenv("ELITE_ROLE_METRICS_BROWSER") != "1" {
		t.Skip("explicit local browser fixture required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool, e := pgxpool.New(ctx, os.Getenv("PAYMENT_CONNECTED_DB_URL"))
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := roleMetricSeed(t, pool)
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'large','store','customer','confirmed','ARS',9007199254740993,1)`, tenant)
	stored, p, program := fixtureStoredValueSetup(t, pool, tenant, "gift_card")
	_ = p
	_ = program
	must(`insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active)values($1,'store','survey','Fixture','v1','Fixture',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1 hour',clock_timestamp()+interval '1 day',2,true)`, tenant)
	metrics, e := db.NewRoleMetrics(pool)
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := catalogRoleIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"admin", "customer", "customer-2", "none", "foreign"} {
		perms := []string{"customer:self"}
		orgs := []string{"store"}
		if name == "admin" {
			perms = []string{"admin:read", "surveys:read", "stored_value:read"}
			orgs = []string{"store", "other"}
		}
		if name == "none" {
			perms = []string{}
		}
		if name == "foreign" {
			perms = []string{"admin:read"}
			orgs = []string{"other"}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": perms, "organizations": orgs, "accessToken": token(name, tenant, perms, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.EnterpriseQueryModule{Service: enterprisequery.NewService(db.NewEnterpriseQuery(pool)), Metrics: metrics}.Register(mux, verifier)
	httpapi.StoredValueModule{Service: stored}.Register(mux, verifier)
	httpapi.CustomerFeedbackModule{Service: customerfeedback.NewService(db.NewCustomerFeedback(pool))}.Register(mux, verifier)
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		counts[r.Method+" "+r.URL.Path]++
		mu.Unlock()
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
	env = append(env, "ELITE_ROLE_METRICS_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_METRIC_IDENTITIES="+string(identitiesJSON), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL,
		"APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-role-metrics-"+uuid.NewString(), "BUSINESS_CONFIG_FILE=role.metrics.fixture.json",
		"CATALOG_RELEASE_ENABLED=false", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=0", "ENTERPRISE_TENANT_CODE=role-catalog", "ENTERPRISE_ORGANIZATION_CODE=role-org", "NEXT_TELEMETRY_DISABLED=1")
	artifacts, e := os.MkdirTemp(web, "role-metrics-artifacts-")
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
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/role-metrics-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, e := command.CombinedOutput()
	if x := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); x != nil {
		t.Fatal(x)
	}
	if e != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("role metrics browser %v\n%s\nartifacts=%s", e, output, artifacts)
	}

	var n int
	if e = pool.QueryRow(ctx, `select count(*)from crm.survey_response where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 2 {
		t.Fatal("durable survey writes", n, e)
	}
	mu.Lock()
	rawCounts, _ := json.Marshal(counts)
	posts := 0
	for k, n := range counts {
		if strings.HasPrefix(k, "POST ") {
			posts += n
		}
	}
	mu.Unlock()
	if posts != 2 {
		t.Fatal("unexpected write", posts)
	}
	if e = os.WriteFile(filepath.Join(artifacts, "api-counts.json"), rawCounts, 0600); e != nil {
		t.Fatal(e)
	}
	t.Logf("ROLE_METRICS_BROWSER_PASS JWE_RS256_JWKS=true PG_BFF_exact_amounts=true scope_and_unavailable=true real_survey_deltas=2 metric_POSTs=0 artifacts=%s", artifacts)
}
````

### FILE: `internal/platform/postgres/role_metrics_connected_integration_test.go`

```yaml
block_id: "GO-CONNECTED-ROLE-METRICS-PROOF:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ad03b1415465a1ff8bf4c1042b4b3adcc47bccf49fe6a0c6b779604ba4d0e5b8"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `microsoft_playwright_browser_gate/tests/role-metrics-connected.spec.mjs`

```yaml
block_id: "GO-CONNECTED-ROLE-METRICS-PROOF:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9d6c8ed0b558a2e74d6116a230f7436ef4aaa7038bb843d59f2d8512a44d2466"
variables: []
secrets_allowed: false
```

````text
import{test,expect}from'@playwright/test';
import{createHash}from'node:crypto';import{createRequire}from'node:module';import{resolve}from'node:path';import{pathToFileURL}from'node:url';
test('Role metrics from PG preserve exactness scope failures and survey updates',async({page,context},info)=>{
 if(process.env.ELITE_ROLE_METRICS_BROWSER!=='1')throw new Error('explicit fixture required');const base=process.env.ELITE_BASE_URL;expect(new URL(base).hostname).toBe('127.0.0.1');expect(new URL(base).protocol).toBe('https:');page.setDefaultTimeout(12000);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href),identities=JSON.parse(process.env.ELITE_METRIC_IDENTITIES);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}]);await page.goto('/dashboard')}
 const errors=[];page.on('pageerror',e=>errors.push(e.message));
 const area=()=>page.getByRole('region',{name:'Indicadores operativos'}),select=name=>area().getByRole('combobox',{name,exact:true});
 async function read(kind,id){if(kind)await select('Indicador').selectOption(kind);if(id)await area().getByRole('textbox').fill(id);const response=page.waitForResponse(r=>r.url().includes('/api/enterprise/metrics?'));await area().getByRole('button',{name:'Consultar indicador',exact:true}).click();return response}
 await identity('none');await expect(page.getByText('No tenés permisos de indicadores asignados.',{exact:true})).toBeVisible();
 await identity('customer');expect(await select('Indicador').locator('option').evaluateAll(xs=>xs.map(x=>x.value))).toEqual(['own-orders','own-cases','own-appointments']);
 expect((await read()).status()).toBe(200);await expect(area().getByRole('cell',{name:'9007199254745993',exact:true})).toBeVisible();await expect(area().getByRole('cell',{name:'100',exact:true})).toBeVisible();
 const forbidden=await page.evaluate(async()=>{const r=await fetch('/api/enterprise/metrics?kind=orders&organization_id=store');return r.status});expect(forbidden).toBe(403);
 await identity('admin');expect((await read()).status()).toBe(200);await expect(area().getByRole('cell',{name:'300',exact:true})).toBeVisible();await expect(area().getByRole('cell',{name:'9007199254745993',exact:true})).toBeVisible();
 await select('Organización').selectOption('other');await expect(area().getByRole('table')).toHaveCount(0);expect((await read()).status()).toBe(200);await expect(area().getByRole('cell',{name:'300',exact:true})).toBeVisible();
 await select('Organización').selectOption('store');expect((await read('stored-value','reference-gift_card')).status()).toBe(200);await expect(area().getByText('Sin registros para este alcance.',{exact:true})).toBeVisible();
 expect((await read('survey','survey')).status()).toBe(200);await expect(area().getByText('Respuestas: 0',{exact:true})).toBeVisible();await expect(area().getByText('NPS no disponible: no se alcanzó el mínimo configurado.',{exact:true})).toBeVisible();
 for(const [name,score]of[['customer',10],['customer-2',2]]){await identity(name);expect(await page.evaluate(async score=>{const r=await fetch('/api/enterprise/surveys?organizationId=store&surveyId=survey',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({score,consent:true,consent_version:'v1'})});return r.status},score)).toBe(201);await identity('admin');expect((await read('survey','survey')).status()).toBe(200);if(score===10)await expect(area().getByText('NPS no disponible: no se alcanzó el mínimo configurado.',{exact:true})).toBeVisible()}
 await expect(area().getByText('Respuestas: 2',{exact:true})).toBeVisible();await expect(area().getByText('NPS: 0',{exact:true})).toBeVisible();
 await page.route('**/api/enterprise/metrics?**',route=>route.fulfill({status:503,contentType:'application/json',body:'{"code":"METRIC_UNAVAILABLE"}'}));expect((await read()).status()).toBe(503);await expect(area().getByRole('alert')).toBeVisible();await expect(area().getByText('NPS: 0',{exact:true})).toHaveCount(0);await page.unrouteAll();
 expect((await read('orders')).status()).toBe(200);await page.screenshot({path:info.outputPath('metrics-desktop.png'),fullPage:true});await page.setViewportSize({width:390,height:844});await page.screenshot({path:info.outputPath('metrics-mobile.png'),fullPage:true});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);expect(errors).toEqual([]);
});
````

## 6. Configuration surface

docs/ROLE_METRICS_REFERENCE.md; features.role_workspace and provider-independent fixture. No credential required.

## 7. Dependency bill

AUTHORED glue and proof only. Existing Go/PG/Next dependencies and licenses are unchanged.

## 8. Apply order

Selected full reference owners must be present; optional browser test requires built Next and contained Playwright runtime. Exact profile rebuilds distinguish compile-only smaller profiles.

## 9. Verification

New metric scope/precision, source-approved issue/reverse, original survey Submit/minimum/retention, minimal original supply writers and actual browser reads/updates. Bounded-body regression and finite fuzz.

## 10. Reconstruction evidence

ROLE_METRICS_RELEASE_V402.md/json; receipts retain RED fixture setup failures and scoped corrected tests. Private locale and later ordered controls remain pending.

