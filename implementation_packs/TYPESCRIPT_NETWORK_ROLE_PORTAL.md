# Connected network role portal

## 1. Metadata

```yaml
pack_id: "TS-NETWORK-ROLE-PORTAL"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Recoverable organization/agreement/branch lifecycle by authorized role, using original network owners; local infrastructure reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804/J5: existing fulfillment organization/agreement owners, tenant/hierarchy/territory guards, approval CanonicalPayload and identity/BFF. No in-memory onboarding-core promotion.

## 3. Architecture contract

Four original PG writer SQL bodies extracted exactly into transaction adapter. Same service states and outbox. Immutable command receipt binds actor/scope/hash/exact result; command lock and preallocated aggregate IDs, independent random event IDs. No parent permission inheritance or automatic grants.

## 4. Exact file manifest

```text
CREATE config/network.role.fixture.json
CREATE deploy/network/role-goldens.json
CREATE docs/NETWORK_ROLE_REFERENCE.md
CREATE microsoft_playwright_browser_gate/tests/network-role-connected.spec.mjs
CREATE src/app/api/enterprise/network/route.test.ts
CREATE src/app/api/enterprise/network/route.ts
CREATE src/app/network/page.tsx
CREATE src/components/network-workspace.tsx
CREATE src/platform/network/contract.test.ts
CREATE src/platform/network/contract.ts
```

## 5. Materialization blocks

### FILE: `config/network.role.fixture.json`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "52f2f3390b0d63f6cf2cecee449f53916232d972aa8d4aba5067a6cea3a547a5"
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
    "network_portal": true
  }
}
````

### FILE: `deploy/network/role-goldens.json`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1d0edc10239ab49c148da6274cc7fdc5d8170b5030b8971a36e89a962b83bdc4"
variables: []
secrets_allowed: false
```

````json
[
  {
    "command": {
      "command_id": "root",
      "action": "create-organization",
      "scope_organization_id": "",
      "entity_id": "root-id",
      "code": "fixture-root",
      "display_name": "Raíz <&> 
 interior 
 final",
      "type": "franchisor"
    },
    "canonical": "{\"action\":\"create-organization\",\"code\":\"fixture-root\",\"command_id\":\"root\",\"display_name\":\"Raíz \\u003c\\u0026\\u003e \\u2028 interior \\u2029 final\",\"entity_id\":\"root-id\",\"scope_organization_id\":\"\",\"type\":\"franchisor\"}",
    "sha256": "25b991dcaa776a384186dabdf6e27f96f5e4d56ab4670835c6cece97965ec726"
  },
  {
    "command": {
      "command_id": "agreement",
      "action": "create-agreement",
      "scope_organization_id": "franchise",
      "entity_id": "agreement-id",
      "territory_code": "fixture",
      "terms_version": "fixture-v1",
      "starts_on": "2026-01-01",
      "ends_on": "2027-01-01"
    },
    "canonical": "{\"action\":\"create-agreement\",\"command_id\":\"agreement\",\"ends_on\":\"2027-01-01\",\"entity_id\":\"agreement-id\",\"scope_organization_id\":\"franchise\",\"starts_on\":\"2026-01-01\",\"terms_version\":\"fixture-v1\",\"territory_code\":\"fixture\"}",
    "sha256": "d26ec61385c14a20d32d938108bbe4bcc154bd0ada813d0d29562952cff0d97c"
  },
  {
    "command": {
      "command_id": "transition-org",
      "action": "transition-organization",
      "scope_organization_id": "org",
      "entity_id": "org",
      "current": "active",
      "target": "suspended",
      "version": "9007199254740993"
    },
    "canonical": "{\"action\":\"transition-organization\",\"command_id\":\"transition-org\",\"current\":\"active\",\"entity_id\":\"org\",\"scope_organization_id\":\"org\",\"target\":\"suspended\",\"version\":\"9007199254740993\"}",
    "sha256": "7ae8a3c26ffa7398636e1f41ae5e28d2038b3f195d754b1f9c5854c68c9b5f9c"
  },
  {
    "command": {
      "command_id": "transition-agreement",
      "action": "transition-agreement",
      "scope_organization_id": "org",
      "entity_id": "agreement-id",
      "current": "draft",
      "target": "active",
      "version": "1"
    },
    "canonical": "{\"action\":\"transition-agreement\",\"command_id\":\"transition-agreement\",\"current\":\"draft\",\"entity_id\":\"agreement-id\",\"scope_organization_id\":\"org\",\"target\":\"active\",\"version\":\"1\"}",
    "sha256": "2ef6724a3d34161098b5e6d4d2c3badc279137161e05a35116dd219f654bda3a"
  }
]
````

### FILE: `docs/NETWORK_ROLE_REFERENCE.md`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "851b3086ef55940c1f52572ddd11a295511121133fec6993cf66653e7bcb2053"
variables: []
secrets_allowed: false
```

````markdown
# Red de franquicia por rol — infraestructura de referencia

/network es opt-in features.network_portal; NETWORK_ROLE_ENABLED=true activa
el host después de verificar migración0077,3triggers,PK e índice del receipt.
No lee cuentas/secretos al estar deshabilitado. UI requiere network:admin o
franchise:write; cada comando conserva el permiso original y el scope exacto.
Sólo el permiso global existente permite crear organizaciones raíz. No se
hereda autoridad del padre para operar una organización hija. El alta no
otorga permisos; el proveedor de identidad gobierna admins/sesiones/accesos.

Fulfillment.Service conserva validaciones/estados y cuatro writers PostgreSQL
conservan exactamente SQL/guards. Se extrae sólo su frontera transaccional.
El adapter permite esas cuatro operaciones y rechaza otras. Organización e ID
de acuerdo se preasignan para recuperar; IDs de eventos independientes aleatorios.
El receipt inmutable liga comando/actor/scope/hash y entidad/version/estado en
la misma transacción del writer y outbox. No segundo ledger de negocio.
GET histórico no sustituye lectura actual antes de una nueva transición.

La estructura activa permite acuerdos y sucursales, no producción/despliegue.
Términos, territorio y fechas son inputs explícitos; no se inventa contrato,
jurisdicción, regalía o criterio de habilitación comercial. Se conservan
no-superposición de territorio, jerarquía y bloqueos para cerrar padres con
hijos/acuerdos activos. Estado cerrado no afirma revocar sesiones del IdP.

BFF autoriza antes del body y limita32KiB/tiempo. El navegador persiste sólo
identidad del comando/hash/resultado esperado bajo tenant+subject, nunca nombres,
términos ni secretos. Una respuesta incierta se recupera GET, sin POST automático.
Si aún no puede confirmarse, conserva referencia para revisión autorizada.
La consulta de un resultado root admite scope vacío/omitido porque protectedGet
omite valores vacíos; el receipt sólo vuelve al actor con autoridad root original.
Importes no se usan; versiones int64 viajan como texto exacto.

TestNetworkRoleAtomic:1nuevo+3replay concurrentes, cuatro tablas revierten en
fallo tardío del receipt. Actor/hash/tenant/org/stale y cierre prematuro rechazados;
territorio superpuesto impedido;12receipts/12outbox finales, histórico/current
separados. TestNetworkRoleBrowser usa Next/BFF/Go/PG/Chromium con JWE/RS256/JWKS;
siembra sólo tenant.13POST crean/activan raíz, franquicia, acuerdo y sucursal,
suspenden/reactivan sucursal y cierran acuerdo/sucursal/franquicia en orden.
Respuestas de alta root y cierre sucursal perdidas y recuperadas tras reload
con0POST adicional. Scope limitado no hereda permisos. Capturas desktop/390px.
TestNetworkRoleHostGuards rechaza ausencia de cada uno de5guards; downgrade
vacío/reapply PASS y downgrade poblado rechaza perder evidencia. Fuzz2s.
10pruebas web y4goldens Go/TS; build/types exactos, sin nueva dependencia.

Pruebas Go requieren PAYMENT_CONNECTED_DB_URL owned/loopback con74migraciones.
Browser añade ELITE_NETWORK_ROLE_BROWSER=1, ELITE_WEB_ROOT/ELITE_NODE_BIN absolutos,
Next compilado--webpack y issuer sintético existente del catálogo. No IdP live.
El companion portal posee los goldens Go/TS; composición completa requerida
para pruebas conectadas, perfil web aislado sólo imports/tipos.

J5: frontend de estructura/acuerdo por rol y ResourceCreate/progreso/evaluación
humana existentes son piezas conectables. Identidad/admin/revocación/rotación,
composición final y readiness de despliegue siguen sus controles T2803/8/1;
este delta no los certifica. Nuevas guías/CMS/KPIs/i18n privados siguen T2804.
AUTHORED sólo glue declarado; ninguna atribución a empresas famosas.
````

### FILE: `microsoft_playwright_browser_gate/tests/network-role-connected.spec.mjs`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f235a1a78ba0577f226164ed930ffe884a4ccf3ce2be79bb73017699bb88950b"
variables: []
secrets_allowed: false
```

````text
import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
test('network organization agreement and branch lifecycle through actual forms',async({page,context},info)=>{
 if(process.env.ELITE_NETWORK_ROLE_BROWSER!=='1')throw new Error('explicit fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');page.setDefaultTimeout(12000);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href),identities=JSON.parse(process.env.ELITE_NETWORK_IDENTITIES);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[],posts=[];page.on('pageerror',e=>errors.push(e.message));page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/network'))posts.push(r.url())});
 const box=name=>page.getByRole('textbox',{name,exact:true}),button=name=>page.getByRole('button',{name,exact:true}),select=name=>page.getByRole('combobox',{name,exact:true});
 async function click(name){await button(name).click();await expect(page.getByRole('status')).toContainText('Operación registrada.');await expect(button('Consultar estado actual')).toBeEnabled()}
 async function consult(){await button('Consultar estado actual').click();await expect(page.getByRole('status')).toContainText('Estado actual consultado.');await expect(button('Consultar estado actual')).toBeEnabled()}
 async function drop(action){let dropped=false;await page.route('**/api/enterprise/network**',async route=>{const r=route.request();if(dropped||r.method()!=='POST'||r.postDataJSON().action!==action){await route.continue();return}dropped=true;const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')})}
 async function recover(){await expect(page.getByRole('status')).toContainText('Resultado sin confirmar');const n=posts.length;await page.unrouteAll();await page.reload();await button('Consultar resultado pendiente').click();await expect(page.getByRole('status')).toContainText('Resultado recuperado');expect(posts.length).toBe(n)}
 async function organization(type,code,parent,loss=false){await select('Tipo de organización').selectOption(type);if(parent)await box('Organización padre').fill(parent);await box('Código de organización').fill(code);await box('Nombre visible').fill('Organización sintética '+code);if(loss){await drop('create-organization');await button('Registrar organización').click();await recover()}else await click('Registrar organización');const id=await box('Referencia a consultar').inputValue();expect(id).toMatch(/^[a-f0-9-]{36}$/);await expect(button('Activar organización')).toBeDisabled();await consult();await click('Activar organización');await consult();return id}
 await identity('unprivileged');await page.goto('/network');await expect(page.getByText('No tenés permiso para consultar este espacio.',{exact:true})).toBeVisible();expect((await context.request.post(base+'/api/enterprise/network',{headers:{origin:base,'content-type':'application/json'},data:{}})).status()).toBe(403);
 await identity('bootstrap');await page.goto('/network');await expect(page.getByRole('heading',{name:'Red de franquicia',exact:true})).toBeVisible();
 const root=await organization('franchisor','root',null,true);
 const franchise=await organization('franchisee','franchise',root);
 await box('Organización del franquiciado').fill(franchise);await box('Territorio').fill('SYNTHETIC');await box('Versión de los términos').fill('fixture-v1');
 await page.getByLabel('Fecha de inicio',{exact:true}).fill('2026-01-01');await page.getByLabel('Fecha de fin opcional',{exact:true}).fill('2027-01-01');await click('Registrar acuerdo en borrador');
 const agreement=await box('Referencia a consultar').inputValue();await consult();await click('Activar acuerdo');await consult();
 const branch=await organization('store','branch',franchise);await click('Suspender organización');await consult();await click('Activar organización');await consult();
 await select('Tipo de consulta').selectOption('agreement');await box('Referencia a consultar').fill(agreement);await box('Organización de la consulta').fill(franchise);await consult();await click('Terminar acuerdo');await consult();
 await select('Tipo de consulta').selectOption('organization');await box('Referencia a consultar').fill(branch);await consult();await drop('transition-organization');await button('Cerrar organización').click();await recover();await consult();await expect(page.getByText('Estado: Cerrada. Versión 5.',{exact:true})).toBeVisible();
 await box('Referencia a consultar').fill(franchise);await consult();await click('Cerrar organización');await consult();await expect(page.getByText('Estado: Cerrada. Versión 3.',{exact:true})).toBeVisible();
 await page.screenshot({path:info.outputPath('network-closed-desktop.png'),fullPage:true});await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('network-closed-mobile.png'),fullPage:true});
 await identity('limited');await page.goto('/network');await expect(page.getByRole('heading',{name:'Registrar acuerdo',exact:true})).toHaveCount(0);await box('Organización padre').fill(root);await box('Código de organización').fill('forbidden');await box('Nombre visible').fill('No permitido');await expect(button('Registrar organización')).toBeDisabled();
 expect((await context.request.get(base+'/api/enterprise/network?'+new URLSearchParams({kind:'organization',id:root,scope_organization_id:root}))).status()).toBe(403);
 expect(posts).toHaveLength(13);expect(errors).toEqual([]);
});
````

### FILE: `src/app/api/enterprise/network/route.test.ts`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "705b9d778c2e80f54a4ccd05810b0b3df0289a895acc180f9133ce4e11b73cd7"
variables: []
secrets_allowed: false
```

````typescript
import{NextRequest}from"next/server";
import{it,expect,vi,afterEach}from"vitest";
const m=vi.hoisted(()=>({session:vi.fn(),get:vi.fn(),post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes("*")||s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:async()=>({features:{network_portal:true}})}));
import{GET,POST}from"./route";
afterEach(()=>vi.clearAllMocks());
const req=(body:BodyInit)=>new NextRequest("https://portal.example.test/api/enterprise/network",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body,duplex:"half"as const});
it("authenticates and authorizes before body",async()=>{for(const session of[null,{permissions:[],organizations:["org"]}]){m.session.mockResolvedValue(session);const r=req("{}");expect([401,403]).toContain((await POST(r)).status);expect(r.bodyUsed).toBe(false)}expect(m.post).not.toHaveBeenCalled()});
it("bounds streaming commands and cancels excess input",async()=>{m.session.mockResolvedValue({subject:"admin",permissions:["network:admin"],organizations:["org"]});let cancelled=false;const body=new ReadableStream<Uint8Array>({pull(c){c.enqueue(new Uint8Array(32769))},cancel(){cancelled=true}});expect((await POST(req(body))).status).toBe(413);expect(cancelled).toBe(true);expect(m.post).not.toHaveBeenCalled()});
it("rejects inherited parent authority and ambiguous queries",async()=>{m.session.mockResolvedValue({subject:"admin",permissions:["network:admin"],organizations:["parent"]});const c={command_id:"cmd",action:"transition-organization",scope_organization_id:"child",entity_id:"child",current:"active",target:"suspended",version:"2"};expect((await POST(req(JSON.stringify(c)))).status).toBe(403);expect((await GET(new NextRequest("https://portal.example.test/api/enterprise/network?kind=organization&id=parent&scope_organization_id=parent&kind=agreement"))).status).toBe(403);expect(m.post).not.toHaveBeenCalled();expect(m.get).not.toHaveBeenCalled()});
it("rejects another actor's root recovery",async()=>{m.session.mockResolvedValue({subject:"admin",permissions:["*"],organizations:[]});m.get.mockResolvedValue({command_id:"cmd",action:"create-organization",scope_organization_id:"",actor:"other",request_sha256:"a".repeat(64),entity:{kind:"organization",id:"root",organization_id:"root",version:"1",state:"provisioning"},recorded_at:"2026-09-12T00:00:00Z",replay:true});expect((await GET(new NextRequest("https://portal.example.test/api/enterprise/network?kind=result&id=cmd&scope_organization_id="))).status).toBe(409);expect(m.get).toHaveBeenCalledWith(expect.anything(),"/v1/franchise/network/commands/cmd",{scope_organization_id:""});expect(m.post).not.toHaveBeenCalled()});
````

### FILE: `src/app/api/enterprise/network/route.ts`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "379977fa88e1dafdaad80b5751d399812bb15d293a5120ae89837cdc8806b1e3"
variables: []
secrets_allowed: false
```

````typescript
import {NextResponse,type NextRequest} from "next/server";
import {allowed,readSession,type PortalSession} from "@/platform/auth/session";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {protectedGet,protectedPost} from "@/platform/backend/protected-client";
import {boundedCommandBody} from "@/platform/backend/bounded-command";
import {loadBusinessConfig} from "@/platform/config/load";
import {networkID,networkCommand,networkReceipt,networkEntity,networkPermission,networkPayload,networkReference,networkMatches} from "@/platform/network/contract";
const fail=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
const ok=(value:unknown)=>NextResponse.json(value,{headers:{"cache-control":"no-store"}});
const scoped=(s:PortalSession,org:string)=>allowed(s,"*")||org!==""&&s.organizations.includes(org);
const anyPermission=(s:PortalSession)=>allowed(s,"network:admin")||allowed(s,"franchise:write");
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.network_portal!==true)return fail("NOT_FOUND",404);
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!anyPermission(s))return fail("FORBIDDEN",403);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),id=q.get("id"),org=q.get("scope_organization_id");
 if(q.size!==3||["kind","id","scope_organization_id"].some(k=>q.getAll(k).length!==1)||!["result","organization","agreement"].includes(kind??"")||!networkID.safeParse(id).success||org===null||org!==""&&!networkID.safeParse(org).success||!scoped(s,org))return fail("INVALID_SCOPE",403);
 try{
  if(kind==="result"){const v=networkReceipt.parse(await protectedGet(s,`/v1/franchise/network/commands/${encodeURIComponent(id!)}`,{scope_organization_id:org}));if(v.command_id!==id||v.scope_organization_id!==org||v.actor!==s.subject||!allowed(s,networkPermission(v.action)))return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(!allowed(s,networkPermission("transition-"+kind)))return fail("FORBIDDEN",403);
  const v=networkEntity.parse(await protectedGet(s,`/v1/franchise/network/entities/${kind}/${encodeURIComponent(id!)}`,{scope_organization_id:org}));if(v.id!==id||v.kind!==kind||v.organization_id!==org)return fail("SCOPE_MISMATCH",409);return ok(v);
 }catch{return fail("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(r:NextRequest){
 if((await loadBusinessConfig()).features.network_portal!==true)return fail("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return fail("CROSS_ORIGIN_REJECTED",403)}catch{return fail("UNAVAILABLE",503)}
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!anyPermission(s))return fail("FORBIDDEN",403);
 if(r.headers.get("content-type")!=="application/json")return fail("UNSUPPORTED_MEDIA_TYPE",415);if(r.nextUrl.searchParams.size)return fail("INVALID_COMMAND",400);
 try{
  let raw:unknown;try{raw=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return fail("INVALID_COMMAND",400)}
  const parsed=networkCommand.safeParse(raw);if(!parsed.success)return fail("INVALID_COMMAND",400);const c=parsed.data;
  if(!allowed(s,networkPermission(c.action))||!scoped(s,c.scope_organization_id))return fail("FORBIDDEN",403);
  const expected=await networkReference(c),v=networkReceipt.parse(await protectedPost(s,"/v1/franchise/network/commands",networkPayload(c)));
  if(!networkMatches(v,expected,s.subject))return fail("UNCONFIRMED",409);return ok(v);
 }catch(e){return fail(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
````

### FILE: `src/app/network/page.tsx`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5c0ba948ca0abad605c4a94ef01bf92a24bae72432b0646414a3caab87a70055"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import {createHash} from "node:crypto";
import {notFound,redirect} from "next/navigation";
import type {Route} from "next";
import {readSession,allowed} from "@/platform/auth/session";
import {loadBusinessConfig} from "@/platform/config/load";
import {NetworkWorkspace} from "@/components/network-workspace";
export default async function NetworkPage(){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.network_portal!==true)notFound();
 const s=await readSession();if(!s)redirect("/api/auth/login?return_to=/network"as Route);
 if(!allowed(s,"network:admin")&&!allowed(s,"franchise:write"))return <><h1>{t("p0116")}</h1><p>{t("p0002")}</p></>;
 const scope=createHash("sha256").update(JSON.stringify([s.tenantId,s.subject])).digest("hex");
 return <><h1 className="pageTitle">{t("p0116")}</h1><p>{t("p0117")}</p><NetworkWorkspace scope={scope} subject={s.subject} permissions={s.permissions} organizations={s.organizations}/></>;
}
````

### FILE: `src/components/network-workspace.tsx`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7a2da943d49fe4df44ea931022f71a0fe74c806af0128e790fef48a71d9a1d98"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import Link from "next/link";
import {useEffect,useRef,useState,type FormEvent} from "react";
import {networkCommand,networkEntity,networkReceipt,networkReference,networkMatches,pendingNetwork,networkTargets,type NetworkEntity,type PendingNetwork} from "@/platform/network/contract";
const endpoint="/api/enterprise/network";


export function NetworkWorkspace({scope,subject,permissions,organizations}:{scope:string;subject:string;permissions:string[];organizations:string[]}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const states:Record<string,string>={provisioning:t("p0057"),active:t("p0705"),suspended:t("p0706"),closed:t("p0707"),draft:t("p0648"),terminated:t("p0708"),expired:t("p0709")};
const actions:Record<string,string>={active:t("p0710"),suspended:t("p0711"),closed:t("p0712"),terminated:t("p0713"),expired:t("p0714")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p),canOrg=(id:string)=>permissions.includes("*")||organizations.includes(id);
 const held=useRef(false),saved=useRef<PendingNetwork|null>(null);const[ready,setReady]=useState(false),[busy,setBusy]=useState(false),[notice,setNotice]=useState(""),[pending,setPending]=useState<PendingNetwork|null>(null);
 const[kind,setKind]=useState("organization"),[id,setID]=useState(""),[org,setOrg]=useState(organizations[0]??""),[entity,setEntity]=useState<NetworkEntity|null>(null),[current,setCurrent]=useState(false);
 const[parent,setParent]=useState(organizations[0]??""),[type,setType]=useState("franchisee"),[code,setCode]=useState(""),[name,setName]=useState("");
 const[agreementOrg,setAgreementOrg]=useState(organizations[0]??""),[territory,setTerritory]=useState(""),[terms,setTerms]=useState(""),[start,setStart]=useState(""),[end,setEnd]=useState("");
 const key=`elite-network:${scope}`;
 function persist(v:PendingNetwork|null){if(v){const raw=JSON.stringify(pendingNetwork.parse(v));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage")}else sessionStorage.removeItem(key);saved.current=v;setPending(v)}
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const v=pendingNetwork.parse(JSON.parse(raw));saved.current=v;setPending(v)}setReady(true)}catch{setNotice(t("p0651"))}},[key]);
 const blocked=!ready||busy||pending!==null;
 async function get(kind:string,id:string,org:string){const r=await fetch(endpoint+"?"+new URLSearchParams({kind,id,scope_organization_id:org}),{cache:"no-store",signal:AbortSignal.timeout(10000)});if(!r.ok)throw new Error("unavailable");return r.json()as Promise<unknown>}
 function show(v:NetworkEntity,isCurrent:boolean){setEntity(v);setKind(v.kind);setID(v.id);setOrg(v.organization_id);setCurrent(isCurrent)}
 async function consult(){if(held.current)return;held.current=true;setBusy(true);setCurrent(false);setEntity(null);try{const v=networkEntity.parse(await get(kind,id,kind==="organization"?id:org));if(v.id!==id||v.kind!==kind||v.organization_id!==(kind==="organization"?id:org))throw new Error("scope");show(v,true);setNotice(t("p0715"))}catch{setNotice(t("p0716"))}finally{held.current=false;setBusy(false)}}
 function accept(raw:unknown,p:PendingNetwork){const v=networkReceipt.parse(raw);if(!networkMatches(v,p,subject))throw new Error("receipt");persist(null);show(v.entity,false);return v.entity}
 async function send(input:unknown){
  if(held.current||!ready||saved.current)return;held.current=true;setBusy(true);let sent=false;
  try{const c=networkCommand.parse(input),p=await networkReference(c);persist(p);sent=true;
   const r=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c),signal:AbortSignal.timeout(10000)});
   if(!r.ok){if([400,401,403,413,415].includes(r.status)){persist(null);sent=false}throw new Error("unconfirmed")}
   accept(await r.json(),p);sent=false;setNotice(t("p0717"));
  }catch{setNotice(sent?t("p0660"):t("p0718"))}finally{held.current=false;setBusy(false)}
 }
 async function recover(){const p=saved.current;if(!p||held.current)return;held.current=true;setBusy(true);try{accept(await get("result",p.command_id,p.scope_organization_id),p);setNotice(t("p0719"))}catch{setNotice(t("p0663"))}finally{held.current=false;setBusy(false)}}
 const createOrganization=(e:FormEvent)=>{e.preventDefault();void send({command_id:crypto.randomUUID(),action:"create-organization",entity_id:crypto.randomUUID(),scope_organization_id:["enterprise","franchisor"].includes(type)?"":parent,code,display_name:name,type})};
 const createAgreement=(e:FormEvent)=>{e.preventDefault();void send({command_id:crypto.randomUUID(),action:"create-agreement",entity_id:crypto.randomUUID(),scope_organization_id:agreementOrg,territory_code:territory,terms_version:terms,starts_on:start,...(end?{ends_on:end}:{})})};
 return <>
 <p role="status" aria-live="polite">{notice}</p>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["network-role-view"]}/>
 {pending&&<section className="card"><h2>{t("p0166")}</h2><p>{t("p0167")} <code>{pending.command_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 {can("network:admin")&&<section className="card"><h2>{t("p0720")}</h2><form onSubmit={createOrganization}><fieldset disabled={blocked}><label>{t("p0721")}<select value={type} onChange={e=>setType(e.target.value)}>{(can("network:bootstrap")?["enterprise","franchisor","franchisee","factory","warehouse","store","service_center"]:["franchisee","factory","warehouse","store","service_center"]).map(v=><option key={v} value={v}>{({enterprise:t("p0722"),franchisor:t("p0723"),franchisee:t("p0724"),factory:t("p0121"),warehouse:t("p0725"),store:t("p0726"),service_center:t("p0727")}as Record<string,string>)[v]}</option>)}</select></label>
 {!["enterprise","franchisor"].includes(type)&&<label>{t("p0728")}<input value={parent} onChange={e=>setParent(e.target.value)} maxLength={128} required/></label>}
 <label>{t("p0729")}<input value={code} onChange={e=>setCode(e.target.value)} maxLength={128} pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*" required/></label><p>{t("p0730")}</p><label>{t("p0731")}<input value={name} onChange={e=>setName(e.target.value)} maxLength={100} required/></label>
 <button className="button" disabled={!["enterprise","franchisor"].includes(type)&&!canOrg(parent)}>{t("p0720")}</button></fieldset></form></section>}
 {can("franchise:write")&&<section className="card"><h2>{t("p0732")}</h2><form onSubmit={createAgreement}><fieldset disabled={blocked}>
 <label>{t("p0733")}<input value={agreementOrg} onChange={e=>setAgreementOrg(e.target.value)} maxLength={128} required/></label><label>{t("p0734")}<input value={territory} onChange={e=>setTerritory(e.target.value)} maxLength={128} required/></label><label>{t("p0735")}<input value={terms} onChange={e=>setTerms(e.target.value)} maxLength={128} required/></label><label>{t("p0736")}<input type="date" value={start} onChange={e=>setStart(e.target.value)} required/></label><label>{t("p0737")}<input type="date" value={end} onChange={e=>setEnd(e.target.value)}/></label><button className="button" disabled={!canOrg(agreementOrg)}>{t("p0738")}</button>
 </fieldset></form></section>}
 <section className="card"><h2>{t("p0739")}</h2><form onSubmit={e=>{e.preventDefault();void consult()}}><fieldset disabled={busy}><label>{t("p0740")}<select value={kind} onChange={e=>{setKind(e.target.value);setEntity(null);setCurrent(false)}}><option value="organization">{t("p0091")}</option><option value="agreement">{t("p0741")}</option></select></label><label>{t("p0742")}<input value={id} onChange={e=>{setID(e.target.value);setEntity(null);setCurrent(false)}} maxLength={128} required/></label>{kind==="agreement"&&<label>{t("p0743")}<input value={org} onChange={e=>{setOrg(e.target.value);setEntity(null);setCurrent(false)}} maxLength={128} required/></label>}<button className="button">{t("p0307")}</button></fieldset></form>
 {entity&&<article><h3>{entity.kind==="organization"?t("p0744"):t("p0745")}</h3><p>{t("p0167")} <code>{entity.id}</code>{t("p0008")}</p><p>{t("p0247")} {states[entity.state]??entity.state}{t("p0674")} {entity.version}{t("p0008")}</p>{entity.display_name&&<p>{entity.display_name} {t("p0016")} {entity.code}</p>}{entity.territory_code&&<p>{t("p0734")} {entity.territory_code}{t("p0746")} {entity.terms_version}{t("p0747")} {entity.starts_on} {t("p0748")} {entity.ends_on??t("p0749")}{t("p0008")}</p>}
 {!current&&<p>{t("p0750")}</p>}
 {networkTargets(entity).map(target=><button className="button" key={target} disabled={blocked||!current||!canOrg(entity.organization_id)||!can(entity.kind==="organization"?"network:admin":"franchise:write")} onClick={()=>void send({command_id:crypto.randomUUID(),action:"transition-"+entity.kind,entity_id:entity.id,scope_organization_id:entity.organization_id,current:entity.state,target,version:entity.version})}>{actions[target]} {entity.kind==="organization"?t("p0751"):"acuerdo"}</button>)}
 </article>}
 </section>
 </>;
}
````

### FILE: `src/platform/network/contract.test.ts`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2194f857443aca17aeb5c715ad27effd022d0d79f7224008575deb5329b28506"
variables: []
secrets_allowed: false
```

````typescript
import {it,expect} from "vitest";
import {readFileSync} from "node:fs";
import {networkCommand,networkReference,networkMatches,networkPayload} from "./contract";
import {catalogCanonical} from "@/platform/catalog/authoring";
const rows=JSON.parse(readFileSync("deploy/network/role-goldens.json","utf8"))as{command:unknown;canonical:string;sha256:string}[];
for(const row of rows)it("matches Go typed command "+String((row.command as {action:string}).action),async()=>{const c=networkCommand.parse(row.command);expect(catalogCanonical(networkPayload(c))).toBe(row.canonical);expect((await networkReference(c)).request_sha256).toBe(row.sha256)});
it("rejects wrong scope, schema codes, overflow and invalid date order",()=>{for(const change of [{scope_organization_id:"foreign"},{version:"9223372036854775807"},{version:"01"}])expect(networkCommand.safeParse({...rows[2]!.command as object,...change}).success).toBe(false);expect(networkCommand.safeParse({...rows[0]!.command as object,code:"UPPER"}).success).toBe(false);expect(networkCommand.safeParse({...rows[1]!.command as object,ends_on:"2025-01-01"}).success).toBe(false)});
it("binds recovered result to actor scope entity and exact effect",async()=>{const c=networkCommand.parse(rows[2]!.command),p=await networkReference(c),v={command_id:p.command_id,action:p.action,scope_organization_id:p.scope_organization_id,actor:"actor",request_sha256:p.request_sha256,entity:{kind:"organization"as const,id:p.entity_id,organization_id:p.entity_id,version:p.version,state:p.state},recorded_at:"2026-09-12T00:00:00Z",replay:true};expect(networkMatches(v,p,"actor")).toBe(true);expect(networkMatches(v,p,"other")).toBe(false);expect(networkMatches({...v,entity:{...v.entity,version:"2"}},p,"actor")).toBe(false)});
````

### FILE: `src/platform/network/contract.ts`

```yaml
block_id: "TS-NETWORK-ROLE-PORTAL:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "194a795b9a526cdee72020c7a51abc20ce24292c7481628a10510ac727706585"
variables: []
secrets_allowed: false
```

````typescript
import {z} from "zod";
import {catalogCanonical,catalogSHA} from "@/platform/catalog/authoring";
export const networkID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/u);
const scope=networkID.or(z.literal(""));
const version=z.string().regex(/^[1-9][0-9]{0,18}$/u).refine(v=>BigInt(v)<9223372036854775807n);
const text=(max:number)=>z.string().min(2).refine(v=>v.trim()===v&&!/[\u0000-\u001f\u007f-\u009f]/u.test(v)&&new TextEncoder().encode(v).length<=max);
const date=z.string().regex(/^[0-9]{4}-[0-9]{2}-[0-9]{2}$/u).refine(v=>v.slice(0,4)!=="0000"&&Number.isFinite(Date.parse(v+"T00:00:00Z"))&&new Date(v+"T00:00:00Z").toISOString().slice(0,10)===v);
const base={command_id:networkID,scope_organization_id:scope,entity_id:networkID};
export const networkCommand=z.discriminatedUnion("action",[
 z.object({...base,action:z.literal("create-organization"),code:z.string().max(128).regex(/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/u),display_name:text(100),type:z.enum(["enterprise","franchisor","franchisee","factory","warehouse","store","service_center"])}).strict(),
 z.object({...base,action:z.literal("create-agreement"),scope_organization_id:networkID,territory_code:networkID,terms_version:networkID,starts_on:date,ends_on:date.optional()}).strict(),
 z.object({...base,action:z.literal("transition-organization"),scope_organization_id:networkID,current:networkID,target:networkID,version}).strict(),
 z.object({...base,action:z.literal("transition-agreement"),scope_organization_id:networkID,current:networkID,target:networkID,version}).strict()
]).refine(c=>c.action!=="transition-organization"||c.scope_organization_id===c.entity_id)
 .refine(c=>c.action!=="create-organization"||(["enterprise","franchisor"].includes(c.type)?c.scope_organization_id==="":c.scope_organization_id!==""))
 .refine(c=>c.action!=="create-agreement"||!c.ends_on||c.ends_on>c.starts_on);
export type NetworkCommand=z.infer<typeof networkCommand>;
export const networkEntity=z.object({
 kind:z.enum(["organization","agreement"]),id:networkID,organization_id:networkID,version:z.string().regex(/^[1-9][0-9]{0,18}$/u).refine(v=>BigInt(v)<=9223372036854775807n),
 state:networkID,code:networkID.optional(),display_name:text(100).optional(),type:networkID.optional(),parent_organization_id:networkID.optional(),
 territory_code:networkID.optional(),terms_version:networkID.optional(),starts_on:date.optional(),ends_on:date.optional()
}).strict();
export type NetworkEntity=z.infer<typeof networkEntity>;
export const networkReceipt=z.object({command_id:networkID,action:z.enum(["create-organization","create-agreement","transition-organization","transition-agreement"]),scope_organization_id:scope,actor:z.string().min(1).max(256),request_sha256:z.string().regex(/^[a-f0-9]{64}$/u),entity:networkEntity,recorded_at:z.string().datetime({offset:true}),replay:z.boolean()}).strict();
export const pendingNetwork=z.object({command_id:networkID,action:networkReceipt.shape.action,scope_organization_id:scope,entity_id:networkID,version:networkEntity.shape.version,state:networkID,request_sha256:networkReceipt.shape.request_sha256}).strict();
export type PendingNetwork=z.infer<typeof pendingNetwork>;
export function networkPermission(action:string){return action.endsWith("-organization")?"network:admin":"franchise:write"}
export function networkPayload(c:NetworkCommand){return Object.fromEntries(Object.entries(c).filter(([k,v])=>v!==undefined&&(v!==""||k==="scope_organization_id")))}
export async function networkReference(c:NetworkCommand):Promise<PendingNetwork>{return pendingNetwork.parse({command_id:c.command_id,action:c.action,scope_organization_id:c.scope_organization_id,entity_id:c.entity_id,version:"version"in c?(BigInt(c.version)+1n).toString():"1",state:"target"in c?c.target:c.action==="create-organization"?"provisioning":"draft",request_sha256:await catalogSHA(catalogCanonical(networkPayload(c)))})}
export function networkMatches(v:z.infer<typeof networkReceipt>,p:PendingNetwork,subject:string){
 return v.command_id===p.command_id&&v.action===p.action&&v.scope_organization_id===p.scope_organization_id&&v.actor===subject&&v.request_sha256===p.request_sha256&&v.entity.id===p.entity_id&&v.entity.version===p.version&&v.entity.state===p.state&&v.entity.kind===(p.action.endsWith("-organization")?"organization":"agreement")&&v.entity.organization_id===(v.entity.kind==="organization"?p.entity_id:p.scope_organization_id);
}
// Existing fulfillment.Service transition tables, surfaced for forms; Go remains authority.
export function networkTargets(v:NetworkEntity){const table:Record<string,string[]>=v.kind==="organization"?{provisioning:["active","closed"],active:["suspended","closed"],suspended:["active","closed"]}:{draft:["active","terminated"],active:["suspended","terminated","expired"],suspended:["active","terminated","expired"]};return table[v.state]??[]}
````

## 6. Configuration surface

docs/NETWORK_ROLE_REFERENCE.md. NETWORK_ROLE_ENABLED optional guard and features.network_portal opt-in. Existing network:admin/franchise:write and organization scope. Root authority remains existing wildcard permission. Explicit terms/territory/dates; org.active is not production readiness.

## 7. Dependency bill

21new and6changed AUTHORED glue/fixtures/config. No new source/dependency/runtime. Migration0077 preserves evidence with immutable receipt and guarded downgrade. Source license/notices unchanged.

## 8. Apply order

Select TS-GO-API-WEB-BRIDGE, TS-OIDC-PORTAL-ADAPTER and TS-FRANCHISE-JOURNEY-PORTALS. Full Go network role and fulfillment owners required when enabled; narrow web proves imports/types only. Shared command goldens owned here.

## 9. Verification

Atomic proof1new3replay,4table late-failure rollback, actor/hash/tenant/org/stale/territory/parent-close guards and historical/current results. Actual Next/BFF/Go/PG/Chromium13POST create/activate root/franchise/agreement/branch, suspend/resume and close in order. Two lost responses GET after reload,0extraPOST. Five host guards, empty downgrade/reapply and populated downgrade refusal.10web tests,4Go/TSgoldens, Next/types, fuzz2s3seeds286419executions. Desktop390px screenshots inspected.

## 10. Reconstruction evidence

NETWORK_ROLE_RELEASE_V402.md/json includes source SQL hash correspondence, RED/PASS and exact profile rebuilds. FAIL861 template extraction,862 qualified types,863 retained SQL lowercase-code fixture corrected without weakening owner guards. WholeT2804 and J5 identity/readiness remain open.


V402 composed delta: T2804 help CMS and21same-release guides;15oldguides unchanged,5bounded training curricula revision2, no automatic grants. New optional host, shared inline text, actual browser/PG proofs. HELP_CMS_RELEASE_V402.md.

V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.

V402 composed delta: Explicit network:admin + network:bootstrap, original scoped ongoing access; IDENTITY_J5_RELEASE_V402.md/json. IdP owns role grants and deprovisioning.