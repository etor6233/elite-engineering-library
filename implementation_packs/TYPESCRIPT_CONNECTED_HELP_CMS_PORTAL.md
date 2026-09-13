# Connected help CMS and release guides

## 1. Metadata

```yaml
pack_id: "TS-CONNECTED-HELP-CMS-PORTAL"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Organization-scoped durable help authoring/publication/history and recovery, with six operation guides shared by UI/help and explicitly versioned human training."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804: existing GO-HELP-CENTER-CORE pure article lifecycle, PG/identity/outbox/approval owners. Public operation guides remain release bound; private CMS content requires explicit reviewed curriculum incorporation.

## 3. Architecture contract

Organization/locale immutable; draft edits only, publish/archive preserve immutable version and actor/command/hash receipt plus outbox in one transaction. Archive withdraws prior published versions from ordinary readers. Explicit current read before another UI write; lost response GET recovery. No public CMS or automatic grant.

## 4. Exact file manifest

```text
CREATE config/help.cms.fixture.json
CREATE deploy/help/cms-goldens.json
CREATE docs/HELP_CMS_REFERENCE.md
CREATE microsoft_playwright_browser_gate/tests/help-cms-connected.spec.mjs
CREATE src/app/api/enterprise/help/cms/route.test.ts
CREATE src/app/api/enterprise/help/cms/route.ts
CREATE src/app/help/library/page.tsx
CREATE src/components/help-cms-workspace.tsx
CREATE src/platform/help/cms-contract.test.ts
CREATE src/platform/help/cms-contract.ts
CREATE src/platform/help/release-guides.test.ts
```

## 5. Materialization blocks

### FILE: `config/help.cms.fixture.json`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f308f40b4a68d25f800eca7f50719c53b6a540c2161f8f0e41361a9d055d5163"
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
    "help_cms": true
  }
}
````

### FILE: `deploy/help/cms-goldens.json`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a8bf04fd5d1f41b1539bb78c2da7057d5ea512c015d5476d72ffb27e8cece4cd"
variables: []
secrets_allowed: false
```

````json
[
  {
    "command": {
      "command_id": "create",
      "action": "create",
      "article_id": "guide",
      "organization_id": "org",
      "locale": "es",
      "category": "operations",
      "title": "Guía <&> segura",
      "body": "Texto   entre   líneas: ñ 😀."
    },
    "canonical": "{\"action\":\"create\",\"article_id\":\"guide\",\"body\":\"Texto \\u2028 entre \\u2029 líneas: ñ 😀.\",\"category\":\"operations\",\"command_id\":\"create\",\"locale\":\"es\",\"organization_id\":\"org\",\"title\":\"Guía \\u003c\\u0026\\u003e segura\"}",
    "sha256": "bc69ad670e36838d0dafc18aeaba1242c30f7ee5100aa12cb8a040e96f9631f4"
  },
  {
    "command": {
      "command_id": "update",
      "action": "update",
      "article_id": "guide",
      "organization_id": "org",
      "version": "9007199254740993",
      "title": "Revisión <&> segura",
      "body": "Texto\ncon\tcampos."
    },
    "canonical": "{\"action\":\"update\",\"article_id\":\"guide\",\"body\":\"Texto\\ncon\\tcampos.\",\"command_id\":\"update\",\"organization_id\":\"org\",\"title\":\"Revisión \\u003c\\u0026\\u003e segura\",\"version\":\"9007199254740993\"}",
    "sha256": "30e06351c5e8062ebf895d98abd679f3427f4b77749d4c187d2db512c173a2d8"
  },
  {
    "command": {
      "command_id": "publish",
      "action": "publish",
      "article_id": "guide",
      "organization_id": "org",
      "version": "9007199254740993"
    },
    "canonical": "{\"action\":\"publish\",\"article_id\":\"guide\",\"command_id\":\"publish\",\"organization_id\":\"org\",\"version\":\"9007199254740993\"}",
    "sha256": "d084b2753ec7d52263db9904895a62fdaec0891f7588bdb081d378bb799a5526"
  },
  {
    "command": {
      "command_id": "archive",
      "action": "archive",
      "article_id": "guide",
      "organization_id": "org",
      "version": "9007199254740993"
    },
    "canonical": "{\"action\":\"archive\",\"article_id\":\"guide\",\"command_id\":\"archive\",\"organization_id\":\"org\",\"version\":\"9007199254740993\"}",
    "sha256": "ead85c81c69706774bb1356ef9b7060932dcc9b0d7a8d5129fd1a2764a217882"
  }
]
````

### FILE: `docs/HELP_CMS_REFERENCE.md`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bf2401a7d0ede163536a5aace78576e5b02af9f8ee383a6e5cfab5198721656c"
variables: []
secrets_allowed: false
```

````markdown
# CMS y guías de la misma revisión

Infraestructura local de biblioteca. HELP_CMS_ENABLED=true activa el host sólo
con migración0078,5guards y PostgreSQL18/UTF8/pg_unicode_fast. Deshabilitado lee
sólo la bandera. features.help_cms habilita /help/library y su BFF. Sin nuevas
cuentas/dependencias. help:read, help:write y help:publish exigen organización
explícita; wildcard conserva la autoridad ya existente, sin permisos heredados.

GO-HELP-CENTER-CORE conserva Store y12pruebas/16semillas. Se extraen funciones
puras de validación/version/draft→published→archived para el adapter durable.
El perfil CMS restringe edición a borradores. Cada transición guarda snapshot
inmutable, actor/comando/hash/version e incorpora outbox en una sola transacción.
Una respuesta perdida se recupera por GET del mismo actor/scope/hash; la UI
conserva sólo referencia/hash/versión bajo tenant+subject, nunca el texto.
El resultado histórico no habilita la siguiente transición sin lectura actual.

Crear/editar/publicar/archivar requiere acción explícita. Publicar no vuelve
público el artículo: sólo lectores de esa organización acceden al publicado.
Archivar retira también consultas a revisiones publicadas anteriores para
lectores; editores conservan historial. No edición silenciosa de publicados:
se prepara otro artículo y se revisa el reemplazo. Texto literal, sin HTML
ejecutable ni ingestión/clasificación de archivos empresariales.
Idioma es/en inmutable por artículo. Búsqueda substring con lower Unicode
pg_unicode_fast de PostgreSQL18.6 ya fijado; no afirma equivalencia exacta con
Go Search ni eliminación de acentos/normalización Unicode. No ICU/dependencia nueva.
Referencia oficial: https://www.postgresql.org/docs/18/collation.html
Contenido≤16KiB, comando codificado≤24KiB; bodyHTTP≤32KiB. Listas/historial50,
cursor por ID/versión, scope y búsqueda acotados; int64 como texto exacto.

Las6guías nuevas de catálogo/supply/garantía/red/CMS/capacitación comparten source
entre UI, /help y training_content/help.bundle.json. Las15anteriores son intactas.
Perfil reference-onboarding revisión2 fija SHA de21guías y5cursos, cada uno≤4
lecciones, mismos límites. Nuevas guías son AUTHORED instrucciones de operación
y extracción exacta de párrafos previamente publicados, no fuente corporativa.
No incorporación automática de artículos privados al curso ni a instrucciones
de IA; nueva versión de curso requiere revisión explícita del contenido.
Evaluación humana por otra persona y cero grants. Una activación futura conserva
el intento/assessment con su revisión original, sin mutación.

Pruebas: CMS1nuevo+3replay, rollback3tablas por outbox tardío, actor/org/tenant,
stale, permisos de publicación, inmutabilidad y retiro, paginación52artículos.
Browser realNext/BFF/Go/PG/Chromium conJWE/RS256/JWKS:6POST para2artículos/6
versiones/6eventos; crear y archivar pierden respuesta y recuperanGET trasreload,
0POST adicional. Reader no ve borradores/historial/retirados, textoHTML literal,
inglés/español por contenido, capturasdesktop390px.
5guards faltantes rechazan host, downgradevacío/reapply y rechazo poblado.
4goldens Go/TS,11webtests y7source/bundle/rolechecks; HTTP niega antesdelbody,
duplicados/oversize/versions. Fuzz2s3seeds328160ejecuciones. Next/tipos PASS.
Capacitación nueva:3guías→3lecturas→respuestas→revisor distinto;6facts6events,
1request1decision0resources/grants; futura revisión3 preserva historial2.

Owned PG con75migraciones; PAYMENT_CONNECTED_DB_URL prefixelite_payment_connected_
paraCMS; ELITE_TRAINING_DATABASE_URL prefixelite_training_connected_ para
TestTrainingReleaseGuidesConnected. Browser requiereELITE_HELP_CMS_BROWSER=1,
ELITE_WEB_ROOT/ELITE_NODE_BIN absolutos yNextbuild--webpack. Usa issuer sintético
de catálogo existente. Perfil web estrecho prueba imports/tipos, no Go remoto.
Usar node tools/export-training-content.mjs con outputausente bajo consumidor
para resolver la dependencia existente; conservar output/SHA después.

FAIL864 búsqueda C-locale corrigida con collationUnicode explícita;865output
del exporter fuera de árbol de módulos;866acceso índice testJSON;867currículo
superó4lecciones y se dividió en5cursos sin ampliar límites. REDs conservados.
HELP_CMS_RELEASE_V402.md/json contiene evidencia y reconstrucciones exactas.
T2804 aún requiereKPIs e i18n privado; resto del orden continúa. No READY global.
````

### FILE: `microsoft_playwright_browser_gate/tests/help-cms-connected.spec.mjs`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fb655ed8c75c2ff4a8b9a1f8d6526ee8cbaa81fc18c8c206767da8db2f25c522"
variables: []
secrets_allowed: false
```

````text
import{test,expect}from'@playwright/test';
import{createHash}from'node:crypto';import{createRequire}from'node:module';import{resolve}from'node:path';import{pathToFileURL}from'node:url';
test('CMS publish read history withdraw with actual role forms',async({page,context},info)=>{
 if(process.env.ELITE_HELP_CMS_BROWSER!=='1')throw new Error('explicit fixture required');const base=process.env.ELITE_BASE_URL;expect(new URL(base).hostname).toBe('127.0.0.1');expect(new URL(base).protocol).toBe('https:');page.setDefaultTimeout(12000);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href),identities=JSON.parse(process.env.ELITE_HELP_IDENTITIES);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[],posts=[];page.on('pageerror',e=>errors.push(e.message));page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/help/cms'))posts.push(r.url())});
 const box=name=>page.getByRole('textbox',{name,exact:true}),button=name=>page.getByRole('button',{name,exact:true}),select=name=>page.getByRole('combobox',{name,exact:true});
 async function click(name){await button(name).click();await expect(page.getByRole('status')).toContainText('Operación registrada.');await expect(button('Consultar artículo actual')).toBeEnabled()}
 async function consult(id){if(id)await box('Referencia del artículo').fill(id);await button('Consultar artículo actual').click();await expect(page.getByRole('status')).toContainText('Artículo actual consultado.');await expect(button('Consultar artículo actual')).toBeEnabled()}
 async function drop(action){let done=false;await page.route('**/api/enterprise/help/cms**',async route=>{if(done||route.request().method()!=='POST'||route.request().postDataJSON().action!==action){await route.continue();return}done=true;const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')})}
 async function recover(){await expect(page.getByRole('status')).toContainText('Resultado sin confirmar');const n=posts.length;await page.unrouteAll();await page.reload();await button('Consultar resultado pendiente').click();await expect(page.getByRole('status')).toContainText('Resultado recuperado');expect(posts.length).toBe(n)}
 await identity('unprivileged');await page.goto('/help/library');await expect(page.getByText('No tenés permiso para consultar este espacio.',{exact:true})).toBeVisible();
 await identity('editor');await page.goto('/help/library');await box('Título del artículo').fill('Recepción segura');await box('Texto del artículo').fill('Contenido <script>window.CMS_UNSAFE=1</script> literal.');await drop('create');await button('Guardar borrador').click();await recover();const es=await box('Referencia del artículo').inputValue();expect(es).toMatch(/^[a-f0-9-]{36}$/);await expect(button('Publicar artículo')).toBeDisabled();
 await identity('reader');await page.goto('/help/library');await box('Referencia del artículo').fill(es);await button('Consultar artículo actual').click();await expect(page.getByRole('status')).toContainText('Artículo no disponible');await expect(button('Guardar borrador')).toHaveCount(0);
 await identity('editor');await page.goto('/help/library');await consult(es);await box('Título del artículo').fill('Recepción revisada');await box('Texto del artículo').fill('Respetar la aprobación humana. <script>window.CMS_UNSAFE=1</script>');await click('Guardar revisión');await consult();await click('Publicar artículo');await consult();
 await button('Consultar historial').click();await expect(page.getByRole('status')).toContainText('Historial consultado');await button('Ver versión 1 · Borrador').click();await expect(page.getByRole('status')).toContainText('Versión histórica consultada');await expect(button('Guardar revisión')).toBeDisabled();await consult();await expect(button('Guardar revisión')).toBeDisabled();
 await button('Preparar artículo nuevo').click();await select('Idioma del contenido').selectOption('en');await box('Título del artículo').fill('Safe receiving');await box('Texto del artículo').fill('Review the serial and retain evidence.');await click('Guardar borrador');const en=await box('Referencia del artículo').inputValue();await consult();await click('Publicar artículo');
 await identity('reader');await page.goto('/help/library');await box('Buscar en título y texto').fill('APROBACIÓN');await button('Buscar artículos').click();await expect(button('Recepción revisada · Publicado')).toBeVisible();await button('Recepción revisada · Publicado').click();await expect(page.getByRole('status')).toContainText('Artículo actual consultado');await expect(page.getByText('Respetar la aprobación humana. <script>window.CMS_UNSAFE=1</script>',{exact:true})).toBeVisible();expect(await page.evaluate(()=>window.CMS_UNSAFE)).toBeUndefined();await expect(button('Consultar historial')).toHaveCount(0);
 await select('Idioma del contenido').selectOption('en');await box('Buscar en título y texto').fill('');await button('Buscar artículos').click();await expect(button('Safe receiving · Publicado')).toBeVisible();await expect(button('Recepción revisada · Publicado')).toHaveCount(0);
 await page.screenshot({path:info.outputPath('help-reader-desktop.png'),fullPage:true});await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('help-reader-mobile.png'),fullPage:true});
 await identity('editor');await page.goto('/help/library');await consult(es);await drop('archive');await button('Archivar artículo').click();await recover();await consult();await expect(page.getByText('Estado: Archivado. Versión 4. Idioma es.',{exact:true})).toBeVisible();
 await identity('reader');await page.goto('/help/library');await box('Referencia del artículo').fill(es);await button('Consultar artículo actual').click();await expect(page.getByRole('status')).toContainText('Artículo no disponible');
 expect((await context.request.get(base+'/api/enterprise/help/cms?'+new URLSearchParams({kind:'article',organization_id:'root',id:es,version:'3'}))).status()).toBe(409);
 await identity('foreign');await page.goto('/help/library');await box('Organización del contenido').fill('root');await expect(button('Buscar artículos')).toBeDisabled();expect((await context.request.get(base+'/api/enterprise/help/cms?'+new URLSearchParams({kind:'article',organization_id:'root',id:en}))).status()).toBe(403);
 expect(posts).toHaveLength(6);expect(errors).toEqual([]);
});
````

### FILE: `src/app/api/enterprise/help/cms/route.test.ts`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a1d0d6cda6be529e5e3d6ed88cb78bdc19a58998a4db87ed25edaf54ce830443"
variables: []
secrets_allowed: false
```

````typescript
import{NextRequest}from"next/server";import{it,expect,vi,afterEach}from"vitest";
const m=vi.hoisted(()=>({session:vi.fn(),get:vi.fn(),post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes("*")||s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:async()=>({features:{help_cms:true}})}));
import{GET,POST}from"./route";
afterEach(()=>vi.clearAllMocks());
const req=(body:BodyInit)=>new NextRequest("https://portal.example.test/api/enterprise/help/cms",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body,duplex:"half"as const});
it("authenticates and authorizes before consuming body",async()=>{for(const session of[null,{permissions:["help:read"],organizations:["org"]}]){m.session.mockResolvedValue(session);const r=req("{}");expect([401,403]).toContain((await POST(r)).status);expect(r.bodyUsed).toBe(false)}expect(m.post).not.toHaveBeenCalled()});
it("bounds and cancels streamed input",async()=>{m.session.mockResolvedValue({subject:"editor",permissions:["help:write"],organizations:["org"]});let cancelled=false;const body=new ReadableStream<Uint8Array>({pull(c){c.enqueue(new Uint8Array(32769))},cancel(){cancelled=true}});expect((await POST(req(body))).status).toBe(413);expect(cancelled).toBe(true);expect(m.post).not.toHaveBeenCalled()});
it("writer cannot publish; foreign scope cannot write",async()=>{m.session.mockResolvedValue({subject:"editor",permissions:["help:write"],organizations:["org"]});expect((await POST(req(JSON.stringify({command_id:"pub",action:"publish",article_id:"guide",organization_id:"org",version:"1"})))).status).toBe(403);expect((await POST(req(JSON.stringify({command_id:"create",action:"create",article_id:"guide",organization_id:"foreign",locale:"es",category:"operations",title:"Title",body:"Body"})))).status).toBe(403);expect(m.post).not.toHaveBeenCalled()});
it("reader cannot request editor history; duplicate scope is rejected",async()=>{m.session.mockResolvedValue({subject:"reader",permissions:["help:read"],organizations:["org"]});expect((await GET(new NextRequest("https://portal.example.test/api/enterprise/help/cms?kind=history&id=guide&organization_id=org"))).status).toBe(403);expect((await GET(new NextRequest("https://portal.example.test/api/enterprise/help/cms?kind=article&id=guide&organization_id=org&organization_id=other"))).status).toBe(403);expect(m.get).not.toHaveBeenCalled()});
it("validates published-only server projection before exposing it",async()=>{m.session.mockResolvedValue({subject:"reader",permissions:["help:read"],organizations:["org"]});m.get.mockResolvedValue({items:[{id:"guide",locale:"es",category:"operations",title:"Hidden",state:"draft",version:"1"}]});expect((await GET(new NextRequest("https://portal.example.test/api/enterprise/help/cms?kind=list&organization_id=org&locale=es"))).status).toBe(409);expect(m.post).not.toHaveBeenCalled()});
````

### FILE: `src/app/api/enterprise/help/cms/route.ts`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0a2a2245f998348ec3a40554fd11a88b4548907577d8d37b30556d2391b06eb2"
variables: []
secrets_allowed: false
```

````typescript
import{NextResponse,type NextRequest}from"next/server";
import{allowed,readSession,type PortalSession}from"@/platform/auth/session";
import{applicationBaseUrl}from"@/platform/auth/oidc-client";
import{protectedGet,protectedPost}from"@/platform/backend/protected-client";
import{boundedCommandBody}from"@/platform/backend/bounded-command";
import{loadBusinessConfig}from"@/platform/config/load";
import{helpID,helpCommand,helpReceipt,helpArticle,helpPage,helpPermission,helpReference,helpMatches}from"@/platform/help/cms-contract";
const fail=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"private, no-store"}});
const ok=(value:unknown)=>NextResponse.json(value,{headers:{"cache-control":"private, no-store"}});
const scoped=(s:PortalSession,org:string)=>allowed(s,"*")||s.organizations.includes(org);
const editor=(s:PortalSession)=>allowed(s,"help:write")||allowed(s,"help:publish");
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.help_cms!==true)return fail("NOT_FOUND",404);
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!editor(s)&&!allowed(s,"help:read"))return fail("FORBIDDEN",403);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),org=q.get("organization_id"),id=q.get("id");
 const keys=kind==="list"?["kind","organization_id","locale","q","after"]:kind==="article"?["kind","organization_id","id","version"]:kind==="history"?["kind","organization_id","id","before"]:["kind","organization_id","id"];
 if(r.url.length>2048||!["list","article","history","result"].includes(kind??"")||[...q.keys()].some(k=>!keys.includes(k)||q.getAll(k).length!==1)||!helpID.safeParse(org).success||!scoped(s,org!))return fail("INVALID_SCOPE",403);
 if(kind!=="list"&&!helpID.safeParse(id).success)return fail("INVALID_QUERY",400);
 if(["history","result"].includes(kind!)&&!editor(s))return fail("FORBIDDEN",403);
 const query:Record<string,string>={organization_id:org!};for(const k of keys.filter(k=>!["kind","organization_id","id"].includes(k))){if(q.has(k))query[k]=q.get(k)!}
 try{
  if(kind==="result"){
   const v=helpReceipt.parse(await protectedGet(s,`/v1/help/cms/commands/${encodeURIComponent(id!)}`,query));
   if(v.command_id!==id||v.article.organization_id!==org||v.actor!==s.subject||!allowed(s,helpPermission(v.action)))return fail("SCOPE_MISMATCH",409);return ok(v)
  }
  if(kind==="list"||kind==="history"){
   const path=kind==="list"?"/v1/help/cms/articles":`/v1/help/cms/articles/${encodeURIComponent(id!)}/history`;
   const v=helpPage.parse(await protectedGet(s,path,query));if(kind==="history"&&v.items.some(a=>a.id!==id)||kind==="list"&&v.items.some(a=>a.locale!==query.locale)||!editor(s)&&v.items.some(a=>a.state!=="published"))return fail("SCOPE_MISMATCH",409);return ok(v)
  }
  const v=helpArticle.parse(await protectedGet(s,`/v1/help/cms/articles/${encodeURIComponent(id!)}`,query));if(v.id!==id||v.organization_id!==org||query.version&&v.version!==query.version||!editor(s)&&v.state!=="published")return fail("SCOPE_MISMATCH",409);return ok(v)
 }catch{return fail("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(r:NextRequest){
 if((await loadBusinessConfig()).features.help_cms!==true)return fail("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return fail("CROSS_ORIGIN_REJECTED",403)}catch{return fail("UNAVAILABLE",503)}
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!editor(s))return fail("FORBIDDEN",403);
 if(r.headers.get("content-type")!=="application/json")return fail("UNSUPPORTED_MEDIA_TYPE",415);if(r.nextUrl.searchParams.size)return fail("INVALID_COMMAND",400);
 try{
  let raw:unknown;try{raw=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return fail("INVALID_COMMAND",400)}
  const parsed=helpCommand.safeParse(raw);if(!parsed.success)return fail("INVALID_COMMAND",400);const c=parsed.data;if(!allowed(s,helpPermission(c.action))||!scoped(s,c.organization_id))return fail("FORBIDDEN",403);
  const p=await helpReference(c),v=helpReceipt.parse(await protectedPost(s,"/v1/help/cms/commands",c));if(!await helpMatches(v,p,s.subject))return fail("UNCONFIRMED",409);return ok(v)
 }catch(e){return fail(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
````

### FILE: `src/app/help/library/page.tsx`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d03183928add89d4a7281df6cd22c717dad1a6f90a6da63a38d6e8b7b868d155"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import{createHash}from"node:crypto";
import{notFound,redirect}from"next/navigation";
import type{Route}from"next";
import{readSession,allowed}from"@/platform/auth/session";
import{loadBusinessConfig}from"@/platform/config/load";
import{HelpCMSWorkspace}from"@/components/help-cms-workspace";
export default async function HelpLibraryPage(){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.help_cms!==true)notFound();const s=await readSession();if(!s)redirect("/api/auth/login?return_to=/help/library"as Route);
 if(!["help:read","help:write","help:publish"].some(p=>allowed(s,p)))return <><h1>{t("p0115")}</h1><p>{t("p0002")}</p></>;
 const scope=createHash("sha256").update(JSON.stringify([s.tenantId,s.subject])).digest("hex");
 return <><h1 className="pageTitle">{t("p0115")}</h1><HelpCMSWorkspace scope={scope} subject={s.subject} permissions={s.permissions} organizations={s.organizations}/></>
}
````

### FILE: `src/components/help-cms-workspace.tsx`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d4167d06b7e305adf8481e570de8cecebcef8d7bd994bc8d99a81ef022643a63"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import Link from"next/link";
import{useEffect,useRef,useState,type FormEvent}from"react";
import{helpCommand,helpArticle,helpPage,helpReceipt,pendingHelp,helpReference,helpMatches,type HelpArticle,type PendingHelp}from"@/platform/help/cms-contract";
const endpoint="/api/enterprise/help/cms";
export function HelpCMSWorkspace({scope,subject,permissions,organizations}:{scope:string;subject:string;permissions:string[];organizations:string[]}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const labels:Record<string,string>={draft:t("p0648"),published:t("p0649"),archived:t("p0650")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p),editor=can("help:write")||can("help:publish"),canOrg=(v:string)=>can("*")||organizations.includes(v);
 const held=useRef(false),saved=useRef<PendingHelp|null>(null);const[ready,setReady]=useState(false),[busy,setBusy]=useState(false),[notice,setNotice]=useState(""),[pending,setPending]=useState<PendingHelp|null>(null);
 const[org,setOrg]=useState(organizations[0]??""),[locale,setLocale]=useState("es"),[query,setQuery]=useState(""),[items,setItems]=useState<ReturnType<typeof helpPage.parse>|null>(null);
 const[id,setID]=useState(""),[article,setArticle]=useState<HelpArticle|null>(null),[current,setCurrent]=useState(false),[history,setHistory]=useState<ReturnType<typeof helpPage.parse>|null>(null);
 const[category,setCategory]=useState("operations"),[title,setTitle]=useState(""),[body,setBody]=useState(""),key=`elite-help-cms:${scope}`;
 function persist(v:PendingHelp|null){if(v){const raw=JSON.stringify(pendingHelp.parse(v));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage")}else sessionStorage.removeItem(key);saved.current=v;setPending(v)}
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const v=pendingHelp.parse(JSON.parse(raw));saved.current=v;setPending(v)}setReady(true)}catch{setNotice(t("p0651"))}},[key]);
 const blocked=!ready||busy||pending!==null;
 async function get(params:Record<string,string>){const r=await fetch(endpoint+"?"+new URLSearchParams(params),{cache:"no-store",signal:AbortSignal.timeout(10000)});if(!r.ok)throw new Error("unavailable");return r.json()as Promise<unknown>}
 function show(v:HelpArticle,isCurrent:boolean){setArticle(v);setOrg(v.organization_id);setID(v.id);setCurrent(isCurrent);setTitle(v.title);setBody(v.body);setCategory(v.category);setLocale(v.locale);setHistory(null)}
 async function read(selected=id,version=""){if(held.current)return;held.current=true;setBusy(true);setArticle(null);setCurrent(false);try{const v=helpArticle.parse(await get({kind:"article",organization_id:org,id:selected,...(version?{version}:{})}));if(v.id!==selected||v.organization_id!==org||version&&v.version!==version)throw new Error("scope");show(v,!version);setNotice(version?t("p0652"):t("p0653"))}catch{setNotice(t("p0654"))}finally{held.current=false;setBusy(false)}}
 async function search(after=""){if(held.current)return;held.current=true;setBusy(true);try{const v=helpPage.parse(await get({kind:"list",organization_id:org,locale,q:query,...(after?{after}:{})}));setItems(v);setNotice(t("p0655"))}catch{setItems(null);setNotice(t("p0656"))}finally{held.current=false;setBusy(false)}}
 async function versions(before=""){if(!article||held.current)return;held.current=true;setBusy(true);try{setHistory(helpPage.parse(await get({kind:"history",organization_id:article.organization_id,id:article.id,...(before?{before}:{})})));setNotice(t("p0657"))}catch{setHistory(null);setNotice(t("p0658"))}finally{held.current=false;setBusy(false)}}
 async function accept(raw:unknown,p:PendingHelp){const v=helpReceipt.parse(raw);if(!await helpMatches(v,p,subject))throw new Error("receipt");persist(null);show(v.article,false)}
 async function send(input:unknown){if(held.current||!ready||saved.current)return;held.current=true;setBusy(true);let sent=false;try{const c=helpCommand.parse(input),p=await helpReference(c);persist(p);sent=true;const r=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c),signal:AbortSignal.timeout(10000)});if(!r.ok){if([400,401,403,413,415].includes(r.status)){persist(null);sent=false}throw new Error("unconfirmed")};await accept(await r.json(),p);sent=false;setNotice(t("p0659"))}catch{setNotice(sent?t("p0660"):t("p0661"))}finally{held.current=false;setBusy(false)}}
 async function recover(){const p=saved.current;if(!p||held.current)return;held.current=true;setBusy(true);try{await accept(await get({kind:"result",organization_id:p.organization_id,id:p.command_id}),p);setNotice(t("p0662"))}catch{setNotice(t("p0663"))}finally{held.current=false;setBusy(false)}}
 function save(e:FormEvent){e.preventDefault();if(article){void send({command_id:crypto.randomUUID(),action:"update",article_id:article.id,organization_id:article.organization_id,version:article.version,title,body})}else void send({command_id:crypto.randomUUID(),action:"create",article_id:crypto.randomUUID(),organization_id:org,locale,category,title,body})}
 return <>
 <p role="status" aria-live="polite">{notice}</p>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["help-cms-view"]}/>
 {pending&&<section className="card"><h2>{t("p0166")}</h2><p>{t("p0167")} <code>{pending.command_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 <section className="card"><h2>{t("p0664")}</h2><fieldset disabled={busy}>
 <label>{t("p0665")}<input maxLength={128} value={org} onChange={e=>{setOrg(e.target.value);setArticle(null);setCurrent(false);setItems(null);setHistory(null);setTitle("");setBody("")}}/></label>
 <label>{t("p0666")}<select value={locale} onChange={e=>{setLocale(e.target.value);setItems(null);setArticle(null);setID("");setCurrent(false);setHistory(null);setTitle("");setBody("")}}><option value="es">{t("p0667")}</option><option value="en">{t("p0668")}</option></select></label>
 <form onSubmit={e=>{e.preventDefault();void search()}}><label>{t("p0669")}<input value={query} onChange={e=>{setQuery(e.target.value);setItems(null)}} maxLength={160}/></label><button className="button" disabled={!canOrg(org)}>{t("p0670")}</button></form>
 {items&&<><ul>{items.items.map(a=><li key={a.id}><button className="button" onClick={()=>void read(a.id)}>{a.title} {t("p0016")} {labels[a.state]}</button></li>)}</ul>{items.next&&<button className="button" onClick={()=>void search(items.next)}>{t("p0671")}</button>}</>}
 <form onSubmit={e=>{e.preventDefault();void read()}}><label>{t("p0672")}<input value={id} onChange={e=>{setID(e.target.value);setArticle(null);setCurrent(false);setHistory(null)}} maxLength={128} required/></label><button className="button" disabled={!canOrg(org)}>{t("p0673")}</button></form>
 </fieldset></section>
 {article&&<article className="card"><h2 lang={article.locale}>{article.title}</h2><p>{t("p0247")} {labels[article.state]}{t("p0674")} {article.version}{t("p0675")} {article.locale}{t("p0008")}</p><p>{t("p0167")} <code>{article.id}</code>{t("p0008")}</p><p lang={article.locale} style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{article.body}</p>{!current&&<p>{t("p0676")}</p>}
 {can("help:publish")&&["draft","published"].includes(article.state)&&<button className="button" disabled={blocked||!current} onClick={()=>void send({command_id:crypto.randomUUID(),action:article.state==="draft"?"publish":"archive",article_id:article.id,organization_id:article.organization_id,version:article.version})}>{article.state==="draft"?t("p0677"):t("p0678")}</button>}
 {editor&&<button className="button" disabled={busy} onClick={()=>void versions()}>{t("p0679")}</button>}
 {history&&<><ul>{history.items.map(a=><li key={a.version}><button className="button" disabled={busy} onClick={()=>void read(a.id,a.version)}>{t("p0680")} {a.version} {t("p0016")} {labels[a.state]}</button></li>)}</ul>{history.next&&<button className="button" disabled={busy} onClick={()=>void versions(history.next)}>{t("p0681")}</button>}</>}
 </article>}
 {can("help:write")&&<section className="card"><h2>{article?t("p0682"):t("p0683")}</h2><button className="button" disabled={blocked} onClick={()=>{setArticle(null);setID("");setTitle("");setBody("");setCurrent(false);setHistory(null)}}>{t("p0684")}</button>
 <form onSubmit={save}><fieldset disabled={blocked||article!==null&&(!current||article.state!=="draft")}>
 {!article&&<label>{t("p0685")}<input value={category} onChange={e=>setCategory(e.target.value)} pattern="[a-z][a-z0-9_-]*" maxLength={64} required/></label>}
 <label>{t("p0686")}<input value={title} onChange={e=>setTitle(e.target.value)} maxLength={200} required/></label><label>{t("p0687")}<textarea value={body} onChange={e=>setBody(e.target.value)} maxLength={16384} rows={8} required/></label><button className="button" disabled={!canOrg(org)}>{article?t("p0688"):t("p0689")}</button></fieldset></form></section>}
 </>
}
````

### FILE: `src/platform/help/cms-contract.test.ts`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8a274884ec0b8a840410d69a4e6d6f556191f45459c5c9000b191023adf81510"
variables: []
secrets_allowed: false
```

````typescript
import{it,expect}from"vitest";
import rows from"../../../deploy/help/cms-goldens.json";
import{helpCommand,helpReference,helpMatches}from"./cms-contract";
import{catalogCanonical,catalogSHA}from"@/platform/catalog/authoring";
const row=(i:number)=>{const v=rows[i];if(!v)throw new Error("required golden missing");return v};
it.each(rows)("typed Go/TS golden $command.action",async r=>{const c=helpCommand.parse(r.command);expect(catalogCanonical(c)).toBe(r.canonical);expect(await catalogSHA(r.canonical)).toBe(r.sha256);expect((await helpReference(c)).request_sha256).toBe(r.sha256)});
it("rejects invalid Unicode, encoded expansion and unsafe versions",()=>{for(const body of["a".repeat(16385),"<".repeat(5000),"bad\u0000body","\ud800"]){expect(helpCommand.safeParse({...row(0).command,body}).success).toBe(false)}for(const version of["0","01","-1","9223372036854775807","9223372036854775808"]){expect(helpCommand.safeParse({...row(1).command,version}).success).toBe(false)}});
it("binds recovered state, exact version, actor and full article digest",async()=>{const c=helpCommand.parse(row(0).command),p=await helpReference(c),article={id:"guide",organization_id:"org",locale:"es"as const,category:"operations",title:"Guía <&> segura",body:"Texto \u2028 entre \u2029 líneas: ñ 😀.",state:"draft"as const,version:"1",updated_at:"2026-09-12T01:02:03.123456Z"},v={command_id:"create",action:"create"as const,actor:"editor",request_sha256:p.request_sha256,article_sha256:await catalogSHA(catalogCanonical(article)),article,replay:true};expect(await helpMatches(v,p,"editor")).toBe(true);expect(await helpMatches({...v,article:{...article,body:"Changed"}},p,"editor")).toBe(false);expect(await helpMatches(v,p,"other")).toBe(false);expect(await helpMatches(v,{...p,version:"2"},"editor")).toBe(false)});
````

### FILE: `src/platform/help/cms-contract.ts`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e5c7baf40523c4f628d6df5948119633503ad43488728d6b0ff6cad037a43f31"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED bounded transport binding; existing Go kernel owns article lifecycle.
import{z}from"zod";
import{catalogCanonical,catalogSHA}from"@/platform/catalog/authoring";
export const helpID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/u);
const version=z.string().regex(/^[1-9][0-9]{0,18}$/u).refine(v=>BigInt(v)<=9223372036854775807n);
const expected=version.refine(v=>BigInt(v)<9223372036854775807n);
const text=(max:number,multiline=false)=>z.string().refine(v=>v.trim()!==""&&new TextEncoder().encode(v).length<=max&&new TextDecoder().decode(new TextEncoder().encode(v))===v&&!(multiline?/[\u0000-\u0008\u000b\u000c\u000e-\u001f\u007f-\u009f]/u:/[\u0000-\u001f\u007f-\u009f]/u).test(v));
const base={command_id:helpID,article_id:helpID,organization_id:helpID};
const locale=z.enum(["es","en"]),category=z.string().regex(/^[a-z][a-z0-9_-]{0,63}$/u),sha=z.string().regex(/^[a-f0-9]{64}$/u);
export const helpCommand=z.discriminatedUnion("action",[
 z.object({...base,action:z.literal("create"),locale,category,title:text(200),body:text(16384,true)}).strict(),
 z.object({...base,action:z.literal("update"),version:expected,title:text(200),body:text(16384,true)}).strict(),
 z.object({...base,action:z.literal("publish"),version:expected}).strict(),
 z.object({...base,action:z.literal("archive"),version:expected}).strict()
]).refine(c=>new TextEncoder().encode(catalogCanonical(c)).length<=24576);
export type HelpCommand=z.infer<typeof helpCommand>;
const state=z.enum(["draft","published","archived"]);
export const helpArticle=z.object({id:helpID,organization_id:helpID,locale,category,title:text(200),body:text(16384,true),state,version,updated_at:z.iso.datetime({offset:true})}).strict();
export type HelpArticle=z.infer<typeof helpArticle>;
export const helpReceipt=z.object({command_id:helpID,action:z.enum(["create","update","publish","archive"]),actor:z.string().min(1).max(256),request_sha256:sha,article_sha256:sha,article:helpArticle,replay:z.boolean()}).strict();
export const helpPage=z.object({items:z.array(helpArticle.omit({organization_id:true,body:true,updated_at:true})).max(50),next:helpID.optional()}).strict();
export const pendingHelp=z.object({command_id:helpID,action:helpReceipt.shape.action,organization_id:helpID,article_id:helpID,version,state,request_sha256:sha}).strict();
export type PendingHelp=z.infer<typeof pendingHelp>;
export const helpPermission=(action:string)=>["publish","archive"].includes(action)?"help:publish":"help:write";
export async function helpReference(c:HelpCommand):Promise<PendingHelp>{return pendingHelp.parse({command_id:c.command_id,action:c.action,organization_id:c.organization_id,article_id:c.article_id,version:c.action==="create"?"1":(BigInt(c.version)+1n).toString(),state:c.action==="publish"?"published":c.action==="archive"?"archived":"draft",request_sha256:await catalogSHA(catalogCanonical(c))})}
export async function helpMatches(v:z.infer<typeof helpReceipt>,p:PendingHelp,subject:string){return v.command_id===p.command_id&&v.action===p.action&&v.actor===subject&&v.request_sha256===p.request_sha256&&v.article.id===p.article_id&&v.article.organization_id===p.organization_id&&v.article.version===p.version&&v.article.state===p.state&&v.article_sha256===await catalogSHA(catalogCanonical(v.article))}
````

### FILE: `src/platform/help/release-guides.test.ts`

```yaml
block_id: "TS-CONNECTED-HELP-CMS-PORTAL:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d8b7e455e2974a75d532a561bbe33de3ac50f24b440f63301475fe61d62a8c5b"
variables: []
secrets_allowed: false
```

````typescript
import{it,expect,vi}from"vitest";
vi.mock("server-only",()=>({}));
import{readFileSync}from"node:fs";import{createHash}from"node:crypto";
import{availableGuides}from"./catalog";import{ALL_GUIDES,RELEASE_GUIDES}from"./content";
it.each([["catalog:read","catalog-role-view"],["supply:read","supply-role-view"],["warranty:read","warranty-role-view"],["network:admin","network-role-view"],["help:read","help-cms-view"],["training:learn","training-role-view"]])("new guide limited to %s",(permission,id)=>{const s={subject:"reader",tenantId:"tenant",permissions:[permission!],organizations:["org"],accessToken:"synthetic"};expect(availableGuides(s,false).map(v=>v.id)).toEqual([id]);expect(availableGuides({...s,permissions:[]},false)).toEqual([])});
it("all21 guides and six inline sources share exact release-bound training content",()=>{const source=readFileSync("src/platform/help/content.ts"),raw=readFileSync("training_content/help.bundle.json"),bundle=JSON.parse(raw.toString()),profile=JSON.parse(readFileSync("deploy/training/reference.profile.json","utf8"));expect(ALL_GUIDES).toHaveLength(21);expect(bundle.source_sha256).toBe(createHash("sha256").update(source).digest("hex"));expect(bundle.articles).toEqual(ALL_GUIDES.map(({id,version,title,paragraphs})=>({id,version,title,paragraphs})));expect(profile.content_sha256).toBe(createHash("sha256").update(raw).digest("hex"));expect(profile.revision).toBe(2);expect(profile.courses).toHaveLength(5);for(const [name,id]of[["supply-workspace","supply-role-view"],["warranty-workspace","warranty-role-view"],["network-workspace","network-role-view"],["help-cms-workspace","help-cms-view"],["catalog-authoring","catalog-role-view"],["training-workspace","training-role-view"]]){expect(readFileSync("src/components/"+name+".tsx","utf8")).toContain('guide={RELEASE_GUIDES["'+id+'"]}')}expect(Object.keys(RELEASE_GUIDES)).toHaveLength(6)});
````

## 6. Configuration surface

docs/HELP_CMS_REFERENCE.md. HELP_CMS_ENABLED host and features.help_cms UI opt-in, explicit help:read/write/publish and organization. PostgreSQL18 UTF8 pg_unicode_fast,5guard checks, bounded24KiB encoded command/16KiB body/32KiBHTTP,50row keysets.

## 7. Dependency bill

23new AUTHORED glue blocks;15selected owner revisions plus pure kernel extraction. No new dependency/runtime/upstream. Existing licenses/notices preserved. Immutable history and populated downgrade guard; original isolated kernel12tests16seeds retained.

## 8. Apply order

TS-GO-API-WEB-BRIDGE0.14.0, TS-OIDC-PORTAL-ADAPTER, TS-FRANCHISE-JOURNEY-PORTALS0.20.0 and catalog canonical helper. Full Go CMS/PG needed when enabled; narrow web proves imports/types only. Shared goldens and browser fixture selected here.

## 9. Verification

CMS PG1new3replay,3table rollback, role/org/actor/tenant/stale/immutability/history/withdrawal/pagination. Browser actual Next/BFF/Go/PG/Chromium6POST2articles6revisions6events, create/archive lost response GET afterreload0extraPOST; desktop390px inspected.5host guards/down/reapply.4goldens/11webtests/7guidealignment checks, HTTP boundaries and Next/types. Fuzz2s3seeds328160executions. Training3newlessons→humanreview6facts6events1request1decision0grants, futureprofile retains old attempt.

## 10. Reconstruction evidence

HELP_CMS_RELEASE_V402.md/json and exact profile rebuilds. FAIL864Unicode search,865export module root,866test indices,867bounded curriculum corrected with RED history and no relaxed constraints. T2804 remains active forKPI/privatei18n; later controls ordered.


V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.
