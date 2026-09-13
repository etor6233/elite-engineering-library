# Connected private locale proof

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-PRIVATE-LOCALE-PROOF"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["V402 full reference source owners"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-11"
```

## 2. Applicability

Local reference proof of existing UI/session/CMS/training owners; not production certification.

## 3. Architecture contract

Display language is never tenant/permission/currency/tax/timezone authority. Original wire, audit reasons and durable recovery references retained.

## 4. Exact file manifest

```text
CREATE config/private.locale.fixture.json
CREATE docs/PRIVATE_LOCALE_REFERENCE.md
CREATE internal/platform/postgres/private_locale_browser_integration_test.go
CREATE microsoft_playwright_browser_gate/tests/private-locale-connected.spec.mjs
CREATE src/platform/i18n/private-display.test.ts
```

## 5. Materialization blocks

### FILE: `config/private.locale.fixture.json`

```yaml
block_id: "GO-CONNECTED-PRIVATE-LOCALE-PROOF:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "706f5200fe8fef0f6dd8e9c6997f9f53fa4f00caf26f212ef1d66e836771b29a"
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
    "role_workspace": true
  }
}
````

### FILE: `docs/PRIVATE_LOCALE_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-PRIVATE-LOCALE-PROOF:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9db7ffb7d7f8f6e60ab6e074ceeb158d9b04b5af400e54fcfdc6a13419dadd2e"
variables: []
secrets_allowed: false
```

````markdown
# Private display locale

AUTHORED presentation glue, not code attributed to a company. Existing fixed Next,
React, Node Intl and original session/BFF owners supply runtime behavior.

Authenticated private displays use per-user/tenant cookie, then bounded browser
language, then configured es/en. The preference is not a permission or business
policy. Saving explicitly reloads the page with a warning about unsaved forms.
Anonymous public pages retain their configured locale.

Money uses the original integer/BigInt algorithm and currency digits. Display
dates use the configured timezone. UTC form command semantics remain unchanged.
Wire status codes, IDs, audit reasons, content and assessment answers are data.

Twenty-one same-release guides translate only when id/version/title/paragraphs
match original source. Five course translations additionally require the exact
profile/content hashes and course fields. Unknown or historical content remains
literal with lang=und; the original training profile/bundle and evidence stay intact.
Tenant CMS content retains its explicit source language. English help search
filters only guides returned by the authorized endpoint, after the original
strict query parser; no permission or API query-boundary bypass.

Reference proof: private-locale-connected.spec.mjs, TestPrivateLocaleBrowser,
private-display.test.ts, private-locale.test.ts, locale/route.test.ts.
Use the existing fixed local runtimes and an owned migrated PostgreSQL fixture.
The connected browser tests a lost create response, explicit language switch,
GET-only recovery and update/publication: exactly three durable writes. No live
account is needed. Operational readiness and provider credentials are separate.
````

### FILE: `internal/platform/postgres/private_locale_browser_integration_test.go`

```yaml
block_id: "GO-CONNECTED-PRIVATE-LOCALE-PROOF:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6896b31cebab908a0c1b84795c736cd2e858c8f5be4b954d658fb3e0cb290b35"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED CMS browser: only tenant/org seeded; all six article mutations from forms.
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

func TestPrivateLocaleBrowser(t *testing.T) {
	if os.Getenv("ELITE_PRIVATE_LOCALE_BROWSER") != "1" {
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

	if _, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'network-browser','Synthetic','Synthetic')`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'root','root','Synthetic','franchisor')`, tenant); e != nil {
		t.Fatal(e)
	}
	store, e := db.NewHelpCMS(pool)
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := catalogRoleIssuer(t)
	identities := map[string]any{}

	for _, name := range []string{"editor", "reader", "foreign", "unprivileged"} {
		permissions := []string{"help:read"}
		orgs := []string{"root"}
		switch name {
		case "editor":
			permissions = []string{"help:read", "help:write", "help:publish"}
		case "foreign":
			orgs = []string{"foreign"}
		case "unprivileged":
			permissions = []string{}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": orgs, "accessToken": token(name, tenant, permissions, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.HelpCMSModule{Service: store}.Register(mux, verifier)
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
	env = append(env, "ELITE_PRIVATE_LOCALE_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_LOCALE_IDENTITIES="+string(identitiesJSON), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL,
		"APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-private-locale-"+uuid.NewString(), "BUSINESS_CONFIG_FILE=private.locale.fixture.json",
		"CATALOG_RELEASE_ENABLED=false", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=0", "ENTERPRISE_TENANT_CODE=role-catalog", "ENTERPRISE_ORGANIZATION_CODE=role-org", "NEXT_TELEMETRY_DISABLED=1")
	artifacts, e := os.MkdirTemp(web, "private-locale-artifacts-")
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
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/private-locale-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
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
 (select count(*)from help.article where tenant_id=$1),
 (select count(*)from help.revision where tenant_id=$1),
 (select count(*)from platform.outbox_event where tenant_id=$1))::text`, tenant).Scan(&actual)
	if e != nil || actual != "[1, 3, 3]" {
		t.Fatal("durable effects", actual, e)
	}
	mu.Lock()
	countJSON, _ := json.Marshal(counts)
	total := 0
	for _, n := range counts {
		total += n
	}
	mu.Unlock()
	if total != 3 {
		t.Fatal("unexpected backend writes", total)
	}
	if e = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countJSON, 0600); e != nil {
		t.Fatal(e)
	}
	t.Logf("PRIVATE_LOCALE_BROWSER_PASS JWE_RS256_JWKS=true language_switch_after_lost_create_GET_recovery=true original_content_retained=true locale_cookie_subject_isolation=true backend_POSTs=3 artifacts=%s", artifacts)
}
````

### FILE: `microsoft_playwright_browser_gate/tests/private-locale-connected.spec.mjs`

```yaml
block_id: "GO-CONNECTED-PRIVATE-LOCALE-PROOF:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ed657d086a2caac4dc2218e2bac8ced1d8b16910dc0d43b0081aed9080cbd77a"
variables: []
secrets_allowed: false
```

````text
import{test,expect}from'@playwright/test';
import{createHash}from'node:crypto';import{createRequire}from'node:module';import{resolve}from'node:path';import{pathToFileURL}from'node:url';
test.use({locale:'es-AR'});
test('Private locale keeps durable commands and original content through English recovery',async({page:originalPage,context:originalContext,browser},info)=>{
 let page=originalPage,context=originalContext;
 if(process.env.ELITE_PRIVATE_LOCALE_BROWSER!=='1')throw new Error('explicit fixture required');const base=process.env.ELITE_BASE_URL;expect(new URL(base).hostname).toBe('127.0.0.1');expect(new URL(base).protocol).toBe('https:');page.setDefaultTimeout(12000);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href),identities=JSON.parse(process.env.ELITE_LOCALE_IDENTITIES);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[],posts=[];page.on('pageerror',e=>errors.push(e.message));page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/help/cms'))posts.push(r.postDataJSON())});
 const box=name=>page.getByRole('textbox',{name,exact:true}),button=name=>page.getByRole('button',{name,exact:true}),select=name=>page.getByRole('combobox',{name,exact:true});
 
 await identity('editor');await page.goto('/help/library');await expect(page.locator('html')).toHaveAttribute('lang','es-AR');
 await box('Título del artículo').fill('Inventario');await box('Texto del artículo').fill('Artículo original <script>window.LOCALE_UNSAFE=1</script>.');
 let dropped=false;await page.route('**/api/enterprise/help/cms**',async route=>{if(dropped||route.request().method()!=='POST'||route.request().postDataJSON().action!=='create'){await route.continue();return}dropped=true;const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')});
 await button('Guardar borrador').click();await expect(page.getByRole('status')).toContainText('Resultado sin confirmar');expect(posts).toHaveLength(1);const original=posts[0];expect(original.locale).toBe('es');expect(original.title).toBe('Inventario');
 await page.unrouteAll();await page.getByText('Idioma de la interfaz',{exact:true}).click();await select('Idioma').selectOption('en');await expect(page.getByText('Guardar recarga esta página.',{exact:false})).toBeVisible();await button('Guardar idioma y recargar').click();
 await expect(page.locator('html')).toHaveAttribute('lang','en-US');await expect(box('Article title')).toBeVisible();expect(posts).toHaveLength(1);
 await button('Check pending result').click();await expect(page.getByRole('status')).toContainText('Result retrieved without repeating the write.');expect(posts).toHaveLength(1);
 const id=await box('Article reference').inputValue();expect(id).toBe(original.article_id);
 await button('Check current article').click();await expect(page.getByRole('status')).toContainText('Current article retrieved.');await expect(page.getByRole('heading',{name:'Inventario',exact:true})).toHaveAttribute('lang','es');
 await expect(page.getByRole('article').getByText('Artículo original <script>window.LOCALE_UNSAFE=1</script>.',{exact:true})).toHaveAttribute('lang','es');expect(await page.evaluate(()=>window.LOCALE_UNSAFE)).toBeUndefined();
 await box('Article title').fill('Inventario revisado');await box('Article text').fill('Texto conservado en español.');await button('Save revision').click();await expect(page.getByRole('status')).toContainText('Operation recorded.');await button('Check current article').click();await expect(page.getByRole('status')).toContainText('Current article retrieved.');await button('Publish article').click();await expect(page.getByRole('status')).toContainText('Operation recorded.');
 expect(posts).toHaveLength(3);expect(posts.map(x=>x.action)).toEqual(['create','update','publish']);expect(posts[1]).not.toHaveProperty('locale');expect(posts[1].body).toBe('Texto conservado en español.');
 await page.goto('/help');await expect(page.getByRole('heading',{name:'Content guide',exact:true})).toBeVisible();await page.getByRole('searchbox',{name:'Search guides',exact:true}).fill('published versions');await button('Search').click();await expect(page.getByRole('heading',{name:'Content guide',exact:true})).toBeVisible();await expect(page.getByText('Published versions are not edited.',{exact:false})).toBeVisible();
 await page.goto('/guide/owner');await expect(page.getByRole('heading',{name:/Owner/})).toBeVisible();await page.screenshot({path:info.outputPath('private-locale-guide-en.png'),fullPage:true});
 await identity('reader');await page.goto('/help/library');await expect(page.locator('html')).toHaveAttribute('lang','es-AR');await expect(button('Guardar borrador')).toHaveCount(0);
 context=await browser.newContext({locale:'en-US',ignoreHTTPSErrors:true});page=await context.newPage();page.setDefaultTimeout(12000);page.on('pageerror',e=>errors.push(e.message));await identity('reader');await page.goto(base+'/help/library');await expect(page.locator('html')).toHaveAttribute('lang','en-US');await box('Article reference').fill(id);await button('Check current article').click();await expect(page.getByRole('heading',{name:'Inventario revisado',exact:true})).toHaveAttribute('lang','es');
 await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('private-locale-reader-mobile.png'),fullPage:true});
 await context.close();context=originalContext;page=originalPage;await identity('editor');await page.goto('/help/library');await expect(page.locator('html')).toHaveAttribute('lang','en-US'); // saved preference wins, belongs to editor only
 await identity('unprivileged');await page.goto('/help/library');await expect(page.locator('html')).toHaveAttribute('lang','es-AR');await expect(page.getByText('No tenés permiso para consultar este espacio.',{exact:true})).toBeVisible();
 expect(posts).toHaveLength(3);expect(errors).toEqual([]);
});
````

### FILE: `src/platform/i18n/private-display.test.ts`

```yaml
block_id: "GO-CONNECTED-PRIVATE-LOCALE-PROOF:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ee19288be054d7d3d83addfc14127fd7bfafd05ee7eff5084b22d0eed91c0466"
variables: []
secrets_allowed: false
```

````typescript
import{it,expect}from"vitest";
import{createElement as h}from"react";
import{renderToStaticMarkup}from"react-dom/server";
import{createHash}from"node:crypto";
import{readFileSync}from"node:fs";
import{privateMessage,controlledPrivateLabel}from"./private-catalog";
import messages from"./private-messages.json";
import{guideDisplay}from"./private-guide-display";
import{trainingCourseDisplay}from"./private-training-display";
import{ALL_GUIDES}from"@/platform/help/content";
import binding from"./private-training-binding.json";
import{courseView}from"@/platform/training/contract";
import{PrivateLocaleProvider}from"./private-provider";
import{RoleDashboard}from"@/components/role-dashboard";
import{TrainingWorkspace}from"@/components/training-workspace";
import{PrivateReadFailure}from"@/platform/backend/private-read-failure";
import{PortalPaging,portalPageHref}from"@/platform/backend/portal-paging";
import{minorAmountPresentation}from"./money";
const en={language:"en",locale:"en-US",timeZone:"America/Argentina/Buenos_Aires",source:"preference"}as const;
const render=(child:React.ReactNode)=>renderToStaticMarkup(h(PrivateLocaleProvider,{locale:en,children:child}));
function course(index:number){const row=binding.courses[index];if(!row)throw new Error("missing fixture course");return courseView.parse({course:row.source,articles:ALL_GUIDES.filter(g=>row.source.lessons.includes(g.id)).map(({id,version,title,paragraphs})=>({id,version,title,paragraphs})),profile_id:binding.profile_id,profile_revision:binding.profile_revision,profile_sha256:binding.profile_sha256,content_sha256:binding.content_sha256,method:"HUMAN_REVIEW_NO_GRANTS_V1"})}
it("retains all template parameters as literal text without HTML or recursive interpolation",()=>{
 const supplied="<img src=x onerror=alert(1)> {handover_state}";
 expect(privateMessage("en","p0629",{handover_id:supplied,handover_state:"accepted"})).toBe("Preparation retrieved: "+supplied+". Status accepted.");
 expect(privateMessage("es","p0629",{handover_id:"H1",handover_state:"accepted"})).toBe("Preparación recuperada: H1. Estado accepted.");
 expect(controlledPrivateLabel("en","Tenant-authored article")).toBe("Tenant-authored article");
 for(const value of Object.values(messages)){expect(value.es.trim()).not.toBe("");expect(value.en.trim()).not.toBe("");expect(value.en.match(/\{[a-z_]+\}/g)??[]).toEqual(value.es.match(/\{[a-z_]+\}/g)??[])}
});
it("translates every exact same-release guide while preserving the source and refusing modified or old content",()=>{
 const original=JSON.stringify(ALL_GUIDES);expect(ALL_GUIDES).toHaveLength(21);
 for(const guide of ALL_GUIDES){if(typeof guide.version!=="string")throw new Error("missing fixed guide version");const source={...guide,version:guide.version};const translated=guideDisplay(source,"en");expect(translated.displayLanguage).toBe("en");expect(translated.title).not.toBe(guide.title);expect(translated.paragraphs.every((p,i)=>p!==guide.paragraphs[i])).toBe(true);expect(guideDisplay(source,"es").paragraphs).toEqual(guide.paragraphs)}
 const known=ALL_GUIDES[0];expect(guideDisplay({...known,version:"0.0.1"},"en").title).toBe(known.title);
 const changed={...known,paragraphs:["Aceptar una cotización"]};expect(guideDisplay(changed,"en").paragraphs).toEqual(changed.paragraphs);expect(guideDisplay(changed,"en").displayLanguage).toBe("und");
 expect(guideDisplay({...known,summary:"Inventario"},"en").summary).toBe("Inventario");
 expect(JSON.stringify(ALL_GUIDES)).toBe(original);
});
it("binds five translated curricula to exact profile and course content without rewriting assessment evidence",()=>{
 const raw=readFileSync("deploy/training/reference.profile.json");expect(createHash("sha256").update(raw).digest("hex")).toBe(binding.profile_sha256);
 expect(createHash("sha256").update(readFileSync("training_content/help.bundle.json")).digest("hex")).toBe(binding.content_sha256);
 for(let i=0;i<5;i++){const value=course(i),before=JSON.stringify(value),view=trainingCourseDisplay(value,"en");expect(view.title).toBe(binding.courses[i]!.en.title);expect(view.prompts).toEqual(binding.courses[i]!.en.prompts);expect(JSON.stringify(value)).toBe(before);
  for(const altered of[{...value,profile_sha256:"0".repeat(64)},{...value,content_sha256:"0".repeat(64)},{...value,profile_revision:1},{...value,course:{...value.course,title:"Mi evaluación"}},{...value,course:{...value.course,prompts:[{id:"review",text:"Texto cambiado"}]}}])expect(trainingCourseDisplay(altered,"en").title).toBe(altered.course.title);
 }
});
it("renders private role links, training and error recovery in English without changing role grants or paging URLs",()=>{
 const sections=[{id:"guide",label:"Guía de uso",href:"/guide"},{id:"supply",label:"Suministro",href:"/supply"}],before=JSON.stringify(sections);
 const dashboard=render(h(RoleDashboard,{sections}));expect(dashboard).toContain("Available access");expect(dashboard).toContain('href="/supply"');expect(dashboard).not.toContain("Suministro");expect(JSON.stringify(sections)).toBe(before);
 const training=render(h(TrainingWorkspace,{courses:[course(0)],assessments:[],canLearn:true,canReview:false,scope:"tenant:user"}));expect(training).toContain("Register resources and recover a response");expect(training).toContain("Start practice");expect(training).not.toContain("Iniciar práctica");
 const failure=render(h(PrivateReadFailure,{path:"/customer"}));expect(failure).toContain("We could not load this information");expect(failure).toContain('href="/customer"');
 const href=portalPageHref("/customer",{orders_after:"opaque-reference"});const paging=render(h(PortalPaging,{label:"Orders",nextHref:href,firstHref:"/customer"}));expect(paging).toContain("View next page");expect(paging).toContain('href="/customer?orders_after=opaque-reference"');
});
it("formats safe minor units exactly in both languages without changing amounts or validity",()=>{
 const amount=9007199254740991;
 expect(minorAmountPresentation(amount,"ARS","en-US")).toEqual({amountLabel:"ARS 90,071,992,547,409.91",amountValid:true});
 expect(minorAmountPresentation(amount,"ARS","es-AR").amountLabel).toContain("90.071.992.547.409,91");
 expect(minorAmountPresentation(NaN,"ARS","en-US")).toEqual({amountLabel:"Amount cannot be verified; contact support.",amountValid:false});
 expect(minorAmountPresentation(100,"JPY","en-US").amountLabel).toBe("JPY 100");
});
````

## 6. Configuration surface

docs/PRIVATE_LOCALE_REFERENCE.md; per-user private preference, explicit save/reload, no provider key.

## 7. Dependency bill

AUTHORED integration and English translation of local source only. Existing fixed runtime/license owners unchanged.

## 8. Apply order

Core UI extensions belong to original packs. Proof requires full CMS/training/browser composition; smaller profiles retain their own compile scope.

## 9. Verification

6preference/auth tests,5display tests,90exact serialization/storage expressions, real browser es/en lost-response recovery and3durableCMSwrites, source-bound guides/courses, strict help query extraction.

## 10. Reconstruction evidence

PRIVATE_LOCALE_RELEASE_V402.md/json. T2804 local closure after exact composition; later controls remain pending.

