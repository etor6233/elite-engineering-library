# Connected warranty role portal

## 1. Metadata

```yaml
pack_id: "TS-WARRANTY-ROLE-PORTAL"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Sold terms, explicit customer acknowledgement, activation and complete warranty repair workflow through authorized role forms and durable recovery; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804/J1/J4: existing sold warranty/coverage/repair, approval/inventory, Next/OIDC/BFF/PG owners. Optional feature enabled only with configured warranty profile.

## 3. Architecture contract

Original writers, domain guards and hash verification retained. RoleContext projected only after existing Claim authorization in its repeatable-read transaction. Quote source read requires warranty:offer; monetary/version int64 as strings. No new state machine.

## 4. Exact file manifest

```text
CREATE config/warranty.role.fixture.json
CREATE deploy/warranty/role-goldens.json
CREATE docs/WARRANTY_ROLE_REFERENCE.md
CREATE microsoft_playwright_browser_gate/tests/warranty-role-connected.spec.mjs
CREATE src/app/api/enterprise/warranty/route.test.ts
CREATE src/app/api/enterprise/warranty/route.ts
CREATE src/app/warranty/page.tsx
CREATE src/components/warranty-workspace.tsx
CREATE src/platform/warranty/contract.test.ts
CREATE src/platform/warranty/contract.ts
```

## 5. Materialization blocks

### FILE: `config/warranty.role.fixture.json`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f79ed8e6bed3d96933021eebb2cf198b11d83aa5773fc0b0a452da540f9ec787"
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
    "warranty_portal": true
  }
}
````

### FILE: `deploy/warranty/role-goldens.json`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a9b29e716842d90e67e0345047ba73e372242c5b8c92fa0d432252b8d8cfead1"
variables: []
secrets_allowed: false
```

````json
[
  {
    "name": "offer",
    "payload": {
      "quote_id": "quote",
      "quote_version": "1",
      "profile_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    },
    "canonical": "{\"profile_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"quote_id\":\"quote\",\"quote_version\":\"1\"}",
    "sha256": "b8b0eab55f3643cdeb85b74ee3bc7f2881723a60266955bc9054dcbca3f4b504"
  },
  {
    "name": "acknowledge",
    "payload": {
      "quote_id": "quote",
      "quote_version": "2",
      "profile_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    },
    "canonical": "{\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"profile_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"quote_id\":\"quote\",\"quote_version\":\"2\"}",
    "sha256": "589926c9ba55da6d4b68f8869f445af1bec991010d921a65dbaeb2b021e83d05"
  },
  {
    "name": "activate",
    "payload": {
      "handover_id": "delivery"
    },
    "canonical": "{\"handover_id\":\"delivery\"}",
    "sha256": "fedffebef6066e493ed048207a0992800221a70beb825c37ca4511228fc14dac"
  },
  {
    "name": "open",
    "payload": {
      "case_id": "case",
      "handover_id": "delivery",
      "appointment_id": "appointment",
      "severity": "medium",
      "description": "Falla á <&>   interior",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    },
    "canonical": "{\"appointment_id\":\"appointment\",\"case_id\":\"case\",\"description\":\"Falla á \\u003c\\u0026\\u003e \\u2028 interior\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"handover_id\":\"delivery\",\"severity\":\"medium\"}",
    "sha256": "7860f9b2bd1c9fe0a883e4a136d9e2af41ca46c2deca1accd5c3297ec4ff9fbe"
  },
  {
    "name": "diagnose",
    "payload": {
      "case_id": "case",
      "command_id": "command",
      "expected_version": "9007199254740993",
      "fault_code": "brake",
      "description": "Diagnóstico <&>   interior",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    },
    "canonical": "{\"case_id\":\"case\",\"command_id\":\"command\",\"description\":\"Diagnóstico \\u003c\\u0026\\u003e \\u2029 interior\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"expected_version\":\"9007199254740993\",\"fault_code\":\"brake\"}",
    "sha256": "3f19f73f1ba0c2cbd6ce16974992676baa986dd2a58266e1348c265216e013a8"
  },
  {
    "name": "plan",
    "payload": {
      "case_id": "case",
      "command_id": "command",
      "expected_version": "9007199254740993",
      "labor_work": "Trabajo sintético",
      "parts": [
        {
          "line_id": "part-line",
          "item_id": "part",
          "bin_id": "bin",
          "lot_id": "lot",
          "quantity": "7"
        }
      ]
    },
    "canonical": "{\"case_id\":\"case\",\"command_id\":\"command\",\"expected_version\":\"9007199254740993\",\"labor_work\":\"Trabajo sintético\",\"parts\":[{\"bin_id\":\"bin\",\"item_id\":\"part\",\"line_id\":\"part-line\",\"lot_id\":\"lot\",\"quantity\":\"7\"}]}",
    "sha256": "cf4117e410e5288f7694a577d3bceacbc1238f740ef006518d2dc7ad20326a83"
  },
  {
    "name": "decide",
    "payload": {
      "case_id": "case",
      "command_id": "command",
      "expected_version": "9007199254740993",
      "payload_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "approved": true,
      "reason": "Revisión explícita"
    },
    "canonical": "{\"approved\":true,\"case_id\":\"case\",\"command_id\":\"command\",\"expected_version\":\"9007199254740993\",\"payload_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"reason\":\"Revisión explícita\"}",
    "sha256": "6752f1a357050813bdd22406966e1926f34161892c99cf0d43099951c77428f7"
  },
  {
    "name": "work",
    "payload": {
      "case_id": "case",
      "command_id": "command",
      "expected_version": "9007199254740993",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    },
    "canonical": "{\"case_id\":\"case\",\"command_id\":\"command\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"expected_version\":\"9007199254740993\"}",
    "sha256": "a09f07c9bb818e7d40720f09d51b934d7af648ab983d54c8c42c163f40ccea76"
  },
  {
    "name": "quality",
    "payload": {
      "case_id": "case",
      "command_id": "command",
      "expected_version": "9007199254740993",
      "passed": false,
      "work_evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "correction_evidence_sha256": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
    },
    "canonical": "{\"case_id\":\"case\",\"command_id\":\"command\",\"correction_evidence_sha256\":\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"expected_version\":\"9007199254740993\",\"passed\":false,\"work_evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\"}",
    "sha256": "c99ebe42a71c8d36eb69565dd4ed1933aa4b7dc02a2f9c9ce4a6416933355714"
  },
  {
    "name": "accept",
    "payload": {
      "case_id": "case",
      "command_id": "command",
      "expected_version": "9007199254740993",
      "quality_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    },
    "canonical": "{\"case_id\":\"case\",\"command_id\":\"command\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"expected_version\":\"9007199254740993\",\"quality_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\"}",
    "sha256": "e85554ddeaf259c9d32a519e16fe2faa1e05f4696e85b3fa3a4af6f8f5c64034"
  },
  {
    "name": "reconcile",
    "payload": {
      "case_id": "case",
      "command_id": "command",
      "expected_version": "9007199254740993",
      "acceptance_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    },
    "canonical": "{\"acceptance_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"case_id\":\"case\",\"command_id\":\"command\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"expected_version\":\"9007199254740993\"}",
    "sha256": "ff26411fbca58b6ac6408f962963cf588bf308a2e084ab1aaa4406f5fc83e4ea"
  },
  {
    "name": "cancel",
    "payload": {
      "case_id": "case",
      "command_id": "command",
      "expected_version": "9007199254740993",
      "reason": "Cancelación explícita",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    },
    "canonical": "{\"case_id\":\"case\",\"command_id\":\"command\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"expected_version\":\"9007199254740993\",\"reason\":\"Cancelación explícita\"}",
    "sha256": "210987fa3f278f3bfcc206443cc987c75327db0fad975a80f0f31fca324f3777"
  }
]
````

### FILE: `docs/WARRANTY_ROLE_REFERENCE.md`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9567164fa67eec96074b8852e858d06f04ff8d4ce526973f6429ac1f6de998e6"
variables: []
secrets_allowed: false
```

````markdown
# Garantía por rol — infraestructura de referencia

/warranty es opt-in features.warranty_portal. Conserva los owners de términos
vendidos, cobertura, reclamos, revisión humana, FIFO, calidad y conciliación.
El nombre de rol no concede permisos. Se exige warranty:read, warranty:self o
warranty:factory-read según vista, más el permiso concreto para cada escritura.
La consulta de cotización exige warranty:offer y entrega la versión actual,
importe exacto como texto y vigencia; BindOffer sigue bloqueando y validando
en PostgreSQL. RoleContext lee los pasos inmutables autorizados dentro de la
misma transacción repeatable-read de Claim, verificando hash y case_id.

El operador revisa términos configurados y los vincula a la cotización.
El cliente lee esos términos y registra recepción explícita. La entrega
aceptada activa la garantía vendida. Los reclamos requieren entrega y atención.
Diagnóstico→plan con repuestos/labor→revisión independiente→trabajo y emisión
FIFO→calidad distinta del técnico→aceptación del cliente→conciliación de fábrica.
El UI presenta las referencias verificadas del paso anterior; no exige copiar
payloads/hashes internos. Exclusión bloquea aprobación, revisión fallida exige
corrección, autoaprobación/autocalidad permanecen rechazadas por el owner.
Cancelación previa al trabajo y rechazo por exclusión usan las rutas originales.
La conciliación es un acuse interno, no un pago ni documento fiscal.

El BFF limita32KiB/tiempo, autentica y autoriza antes del body, fija rutas y
valida recibos contra tenant/org/actor/referencia/versión/hash. Command recovery
guarda sólo referencia y hash, aislados por tenant/org/subject/vista. Guarda
antes de POST y recupera por GET tras una respuesta incierta incluso con reload;
no reenvía la escritura. Ni archivos ni motivos se conservan en ese marcador.
La evidencia se representa por su huella; el portal no afirma almacenar o
clasificar documentos. Una referencia aún incierta permanece para revisión.

TestWarrantyRoleBrowser usa Next/BFF/Go/PG/Chromium y sesiones JWE/RS256/JWKS
sintéticas. Habilitar ELITE_WARRANTY_ROLE_BROWSER=1, ELITE_WEB_ROOT,
ELITE_NODE_BIN y PAYMENT_CONNECTED_DB_URL aislada elite_payment_connected_*.
HANDOVER_PROFILE_PYTHON fija el builder admitido. La primera fase realiza
términos/consentimiento en UI antes de que el fixture continúe los owners
quote/order/SDK/callback/reconciliación/checklist/handover/release ya admitidos.
Los owners de atención y stock preparan tres turnos y dos capas FIFO. La
segunda fase ejecuta activación y todos los pasos de tres casos desde forms.
18POST,15pasos,1caso cerrado/2cancelados,3unidades restantes, salida de inventario
-740.0000 y costo positivo740.0000 presentado en el servicio. Dos respuestas
perdidas recuperadas con GET y ceroPOST adicional. Calidad fallida/corrección,
exclusión, cancelación, foreign y otro cliente incluidos. Desktop/390px.

12goldens Go/TS verifican los comandos tipados, Unicode/interiores/HTML,
int64 fuera del rango seguro de JavaScript y arrays;18tests web, build/types.
No nuevo algoritmo, migración, runtime o dependencia. El fuzz de dominio y
guards de host/migración de WARRANTY_CONNECTED_RELEASE_V402 permanecen.
El nuevo fixture de browser comparte el issuer sintético del pack de catálogo
y los goldens son compartidos con el portal. Composición completa obligatoria
para la prueba conectada; perfil web aislado sólo prueba imports/tipos.

FAIL859 fue una expectativa positiva en el ledger de salida, después de que
ambas fases browser pasaron. Se conserva el RED y se corrige la expectativa;
la base ya guardada se verificó read-only sin repetir escrituras.
Ayuda inline1.0.0 en esta release; integración CMS/capacitación e i18n privados
continúan en T2804. AUTHORED sólo proyección/transporte/UI/fixtures; fuentes
BC/Go/PG y sus atribuciones exactas conservadas, sin atribución empresarial nueva.
No cierre de todo T2804, infraestructura global ni producción.
````

### FILE: `microsoft_playwright_browser_gate/tests/warranty-role-connected.spec.mjs`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8ef9aa926d1ea87494ad1071a09ec5b8807682a2525ace97df2465701aea0541"
variables: []
secrets_allowed: false
```

````text
import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
test('sold terms and connected warranty through role forms',async({page,context},info)=>{
 if(process.env.ELITE_WARRANTY_ROLE_BROWSER!=='1')throw new Error('explicit fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');page.setDefaultTimeout(12000);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href),identities=JSON.parse(process.env.ELITE_WARRANTY_IDENTITIES),inputs=JSON.parse(process.env.ELITE_WARRANTY_INPUTS);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[],posts=[];page.on('pageerror',e=>errors.push(e.message));page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/warranty'))posts.push(r.url())});
 const textbox=name=>page.getByRole('textbox',{name,exact:true});
 const button=name=>page.getByRole('button',{name,exact:true});
 async function evidence(suffix){await page.getByLabel('Archivo de evidencia',{exact:true}).setInputFiles({name:'fixture.txt',mimeType:'text/plain',buffer:Buffer.from('Authorized synthetic evidence '+suffix)});await expect(page.getByText('Evidencia preparada.',{exact:true})).toBeVisible();await textbox('Motivo de la decisión').fill('Revisión de la evidencia sintética '+suffix)}
 async function open(name,params={}){await identity(name);await page.goto('/warranty?'+new URLSearchParams(params));await expect(page.getByRole('heading',{name:'Garantía y reparación',exact:true})).toBeVisible();if(params.case){await button('Consultar caso').click();await expect(page.getByRole('status')).toContainText('Consulta actualizada.')}await evidence(name)}
 async function click(name){await button(name).click();await expect(page.getByRole('status')).toContainText('Operación registrada.');await expect(button('Consultar caso')).toBeEnabled()}
 async function consult(name){await button(name).click();await expect(page.getByRole('status')).toContainText('Consulta actualizada.');await expect(button(name)).toBeEnabled()}
 async function drop(action){let dropped=false;await page.route('**/api/enterprise/warranty**',async route=>{const r=route.request();if(dropped||r.method()!=='POST'||r.postDataJSON().action!==action){await route.continue();return}dropped=true;const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')})}
 async function recover(){await expect(page.getByRole('status')).toContainText('Resultado sin confirmar');const n=posts.length;await page.unrouteAll();await page.reload();await button('Consultar resultado pendiente').click();await expect(page.getByRole('status')).toContainText('Resultado recuperado');await expect(button('Consultar caso')).toBeEnabled();expect(posts.length).toBe(n)}
 if(process.env.ELITE_WARRANTY_PHASE==='terms'){
  await identity('unprivileged');await page.goto('/warranty');await expect(page.getByText('No tenés permiso para consultar este espacio.',{exact:true})).toBeVisible();
  await identity('reader');expect((await context.request.post(base+'/api/enterprise/warranty',{headers:{origin:base,'content-type':'application/json'},data:{}})).status()).toBe(403);
  await open('operator',{quote:inputs.quote});await consult('Consultar cotización para términos');await consult('Consultar política configurada');await drop('offer');await button('Vincular términos revisados').click();await recover();await expect(page.getByText('Recepción de términos pendiente.',{exact:true})).toBeVisible();
  await open('customer',{quote:inputs.quote});await consult('Consultar términos ofrecidos');await expect(button('Registrar recepción de términos')).toBeDisabled();await page.getByRole('checkbox',{name:'Leí los términos mostrados de esta cotización',exact:true}).check();await click('Registrar recepción de términos');await expect(page.getByText('Recepción de términos registrada.',{exact:true})).toBeVisible();
  expect(posts).toHaveLength(2);
 }else{
  expect(process.env.ELITE_WARRANTY_PHASE).toBe('claims');
  await open('operator',{handover:inputs.handover});await click('Activar garantía de entrega aceptada');await expect(page.getByRole('heading',{name:'Garantía activa',exact:true})).toBeVisible();
  async function claim(appointment){await open('customer',{handover:inputs.handover});await consult('Consultar garantía activa');await textbox('Referencia de la atención completada').fill(appointment);await textbox('Descripción del problema').fill('Problema sintético de freno');await click('Abrir reclamo');const id=await textbox('Referencia del caso').inputValue();expect(id).toMatch(/^[a-f0-9-]{36}$/);return id}
  async function diagnosis(id,fault){await open('operator',{case:id});await textbox('Código de falla').fill(fault);await textbox('Descripción del diagnóstico').fill('Inspección sintética de freno');await click('Registrar diagnóstico')}
  async function plan(withParts){await textbox('Trabajo propuesto').fill('Ajuste sintético de freno');if(withParts){await button('Agregar repuesto').click();for(const [name,value]of [['Artículo','warranty-part'],['Ubicación de stock','warranty-bin'],['Lote opcional','warranty-lot'],['Cantidad de repuesto','7']])await textbox(name).fill(value)}await click('Proponer reparación');await expect(button('Aprobar reparación')).toBeDisabled()}
  const main=await claim(inputs.main);await diagnosis(main,'fixture-wear');await plan(true);
  await open('reviewer',{case:main});await expect(page.getByRole('heading',{name:'Plan de reparación',exact:true})).toBeVisible();await expect(page.getByText('warranty-part: 7 · warranty-bin · lote warranty-lot',{exact:true})).toBeVisible();await click('Aprobar reparación');
  await open('operator',{case:main});await drop('work');await button('Registrar trabajo realizado').click();await recover();await evidence('after-work');await expect(button('Aprobar calidad del trabajo')).toBeDisabled();
  await open('quality',{case:main});await click('Registrar calidad fallida');await expect(button('Aprobar calidad del trabajo')).toBeDisabled();await page.getByLabel('Evidencia de corrección',{exact:true}).setInputFiles({name:'correction.txt',mimeType:'text/plain',buffer:Buffer.from('Synthetic corrected brake')});await expect(button('Aprobar calidad del trabajo')).toBeEnabled();await click('Aprobar calidad del trabajo');
  await open('customer',{case:main});await click('Aceptar reparación recibida');await expect(page.getByText('Reparación aceptada por el cliente.',{exact:true})).toBeVisible();
  await open('factory',{case:main});await click('Conciliar servicio de garantía');await expect(page.getByText('Cerrado',{exact:true})).toBeVisible();await expect(page.getByText('Costo de inventario registrado: 740.0000.',{exact:true})).toBeVisible();
  await page.screenshot({path:info.outputPath('warranty-closed-desktop.png'),fullPage:true});await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('warranty-closed-mobile.png'),fullPage:true});
  const cancelled=await claim(inputs.cancel);await open('operator',{case:cancelled});await click('Cancelar reparación sin trabajo emitido');await expect(page.getByText('Cancelado',{exact:true})).toBeVisible();
  const excluded=await claim(inputs.excluded);await diagnosis(excluded,'fixture-exclusion');await plan(false);await open('reviewer',{case:excluded});await expect(button('Aprobar reparación')).toBeDisabled();await click('Rechazar reparación');await expect(page.getByText('Cancelado',{exact:true})).toBeVisible();
  await identity('other-customer');expect((await context.request.get(base+'/api/enterprise/warranty?'+new URLSearchParams({kind:'claim',organization_id:'store',id:main,surface:'customer'}))).status()).toBe(409);
  await identity('foreign');expect((await context.request.get(base+'/api/enterprise/warranty?'+new URLSearchParams({kind:'claim',organization_id:'foreign',id:main,surface:'franchise'}))).status()).toBe(409);
  expect(posts).toHaveLength(16);
 }
 expect(errors).toEqual([]);
});
````

### FILE: `src/app/api/enterprise/warranty/route.test.ts`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fb2752c90468a542d570cf4b37a9e4b5bfdaef67a0f11f0aad13b52414ac00c6"
variables: []
secrets_allowed: false
```

````typescript
import {NextRequest} from "next/server";
import {it,expect,vi,afterEach} from "vitest";
const m=vi.hoisted(()=>({session:vi.fn(),get:vi.fn(),post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:async()=>({features:{warranty_portal:true}})}));
import{GET,POST}from"./route";
afterEach(()=>vi.clearAllMocks());
const req=(body:BodyInit)=>new NextRequest("https://portal.example.test/api/enterprise/warranty",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body,duplex:"half" as const});
it("authenticates before consuming bodies",async()=>{m.session.mockResolvedValue(null);const r=req("{}");expect((await POST(r)).status).toBe(401);expect(r.bodyUsed).toBe(false);expect(m.post).not.toHaveBeenCalled()});
it("bounds and cancels chunked commands",async()=>{m.session.mockResolvedValue({subject:"customer",permissions:["warranty:self"],organizations:["store"]});let cancelled=false;const body=new ReadableStream<Uint8Array>({pull(c){c.enqueue(new Uint8Array(32769))},cancel(){cancelled=true}});expect((await POST(req(body))).status).toBe(413);expect(cancelled).toBe(true);expect(m.post).not.toHaveBeenCalled()});
it("rejects customer work, ambiguous scope and foreign organization before backend",async()=>{m.session.mockResolvedValue({subject:"customer",permissions:["warranty:self"],organizations:["store"]});const c={surface:"customer",organization_id:"store",action:"work",case_id:"case",command_id:"cmd",expected_version:"4",evidence_sha256:"a".repeat(64)};expect((await POST(req(JSON.stringify(c)))).status).toBe(403);for(const tail of["&organization_id=store","&surprise=x"]){expect((await GET(new NextRequest("https://portal.example.test/api/enterprise/warranty?kind=claim&id=case&surface=customer&organization_id=store"+tail))).status).toBe(403)};expect(m.post).not.toHaveBeenCalled();expect(m.get).not.toHaveBeenCalled()});
it("does not return another actor's recovery receipt",async()=>{m.session.mockResolvedValue({subject:"customer",permissions:["warranty:self"],organizations:["store"]});m.get.mockResolvedValue({case_id:"case",command_id:"open",version:"1",kind:"opened",state:"opened",actor:"other",request_sha256:"a".repeat(64),payload_sha256:"b".repeat(64),payload:{},recorded_at:"2026-09-12T00:00:00Z",replay:false});expect((await GET(new NextRequest("https://portal.example.test/api/enterprise/warranty?kind=command&id=case&command_id=open&surface=customer&organization_id=store"))).status).toBe(409);expect(m.post).not.toHaveBeenCalled()});
````

### FILE: `src/app/api/enterprise/warranty/route.ts`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3bb4a9beef56534cc450c52f6063b05bbad784107716dcbbb50f527683d6ff02"
variables: []
secrets_allowed: false
```

````typescript
import {NextResponse,type NextRequest} from "next/server";
import {readSession,allowed} from "@/platform/auth/session";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {protectedGet,protectedPost} from "@/platform/backend/protected-client";
import {boundedCommandBody} from "@/platform/backend/bounded-command";
import {loadBusinessConfig} from "@/platform/config/load";
import {warrantyID,warrantySurface,warrantyCommand,warrantyPermission,warrantyReadPermission,warrantyPath,warrantyPayload,warrantyOffer,warrantyActivation,warrantyClaim,warrantyStep,warrantyStepReference,warrantyStepMatches,warrantyProfileView,warrantyQuote} from "@/platform/warranty/contract";
const fail=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
const ok=(value:unknown)=>NextResponse.json(value,{headers:{"cache-control":"no-store"}});
const writePermissions=["warranty:self","warranty:offer","warranty:activate","warranty:request","warranty:diagnose","warranty:plan","warranty:approve","warranty:work","warranty:quality","warranty:reconcile","warranty:cancel"];
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.warranty_portal!==true)return fail("NOT_FOUND",404);
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),org=q.get("organization_id"),id=q.get("id"),command=q.get("command_id"),surface=warrantySurface.safeParse(q.get("surface"));
 const keys=["kind","organization_id","surface",...(kind!=="profile"?["id"]:[]),...(kind==="command"?["command_id"]:[])];
 if(!surface.success||!["profile","quote","offer","activation","claim","command"].includes(kind??"")||q.size!==keys.length||keys.some(k=>q.getAll(k).length!==1)||!warrantyKeys(q,keys)||!warrantyID.safeParse(org).success||!s.organizations.includes(org!)||kind!=="profile"&&!warrantyID.safeParse(id).success||kind==="command"&&!warrantyID.safeParse(command).success)return fail("INVALID_SCOPE",403);
 const role=surface.data;if(kind==="quote"&&(role!=="franchise"||!allowed(s,"warranty:offer")))return fail("FORBIDDEN",403);if(!allowed(s,warrantyReadPermission(role))||kind==="profile"&&role!=="franchise"||["offer","activation"].includes(kind!)&&role==="factory")return fail("FORBIDDEN",403);
 try{
  const tail=kind==="quote"?`quotes/${encodeURIComponent(id!)}/source`:kind==="profile"?"profile":kind==="offer"?`quotes/${encodeURIComponent(id!)}`:kind==="activation"?`handovers/${encodeURIComponent(id!)}/activation`:`claims/${encodeURIComponent(id!)}`;
  const raw=await protectedGet(s,`/v1/${role}/warranty/${tail}`,{organization_id:org!,command_id:command??undefined});
  if(kind==="quote"){const v=warrantyQuote.parse(raw);if(v.quote_id!==id||v.organization_id!==org)return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(kind==="profile"){const v=warrantyProfileView.parse(raw);if(v.profile.tenant_id!==s.tenantId||v.profile.organization_id!==org)return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(kind==="offer"){const v=warrantyOffer.parse(raw);if(v.quote_id!==id||v.profile.tenant_id!==s.tenantId||v.profile.organization_id!==org||role==="customer"&&v.customer_subject!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(kind==="activation"){const v=warrantyActivation.parse(raw);if(v.handover_id!==id||v.organization_id!==org||role==="customer"&&v.customer_subject!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(kind==="command"){const v=warrantyStep.parse(raw);if(v.case_id!==id||v.command_id!==command||v.actor!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v)}
  const v=warrantyClaim.parse(raw);if(v.case_id!==id||(role==="factory"?v.factory_organization_id:v.organization_id)!==org||role==="customer"&&v.customer_subject!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v);
 }catch{return fail("CONSULTATION_UNAVAILABLE",409)}
}
function warrantyKeys(q:URLSearchParams,keys:string[]){return [...q.keys()].every(k=>keys.includes(k))}
export async function POST(r:NextRequest){
 if((await loadBusinessConfig()).features.warranty_portal!==true)return fail("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return fail("CROSS_ORIGIN_REJECTED",403)}catch{return fail("UNAVAILABLE",503)}
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!writePermissions.some(p=>allowed(s,p)))return fail("FORBIDDEN",403);
 if(r.headers.get("content-type")!=="application/json")return fail("UNSUPPORTED_MEDIA_TYPE",415);if(r.nextUrl.searchParams.size)return fail("INVALID_COMMAND",400);
 try{
  let raw:unknown;try{raw=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return fail("INVALID_COMMAND",400)}
  const parsed=warrantyCommand.safeParse(raw);if(!parsed.success)return fail("INVALID_COMMAND",400);const c=parsed.data,permission=warrantyPermission(c);
  if(!permission||!s.organizations.includes(c.organization_id)||!allowed(s,permission)||!allowed(s,warrantyReadPermission(c.surface)))return fail("FORBIDDEN",403);
  const result=await protectedPost(s,warrantyPath(c)+"?"+new URLSearchParams({organization_id:c.organization_id}),warrantyPayload(c));
  if(c.action==="offer"||c.action==="acknowledge"){
   const v=warrantyOffer.parse(result);if(v.quote_id!==c.quote_id||v.profile_sha256!==c.profile_sha256||v.profile.tenant_id!==s.tenantId||v.profile.organization_id!==c.organization_id||v.quote_version!==(c.action==="offer"?(BigInt(c.quote_version)+1n).toString():c.quote_version))return fail("UNCONFIRMED",409);
   if(c.action==="offer"?v.offered_by!==s.subject:v.customer_subject!==s.subject||!v.acknowledged||v.evidence_sha256!==c.evidence_sha256)return fail("UNCONFIRMED",409);return ok(v);
  }
  if(c.action==="activate"){const v=warrantyActivation.parse(result);if(v.handover_id!==c.handover_id||v.organization_id!==c.organization_id||v.activated_by!==s.subject)return fail("UNCONFIRMED",409);return ok(v)}
  const v=warrantyStep.parse(result);if(!warrantyStepMatches(v,await warrantyStepReference(c),s.subject))return fail("UNCONFIRMED",409);return ok(v);
 }catch(e){return fail(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
````

### FILE: `src/app/warranty/page.tsx`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "58f9fe967decfc656e1dcaf358357168d23d9d4f8065915789a28e1ad39f78de"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import {createHash} from "node:crypto";
import {redirect,notFound} from "next/navigation";
import type {Route} from "next";
import {readSession,allowed} from "@/platform/auth/session";
import {loadBusinessConfig} from "@/platform/config/load";
import {warrantyID,warrantySurface,warrantyReadPermission} from "@/platform/warranty/contract";
import {WarrantyWorkspace} from "@/components/warranty-workspace";
export default async function WarrantyPage({searchParams}:{searchParams:Promise<Record<string,string|string[]|undefined>>}){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.warranty_portal!==true)notFound();const s=await readSession();if(!s)redirect("/api/auth/login?return_to=/warranty" as Route);
 const available=(["franchise","customer","factory"] as const).filter(role=>allowed(s,warrantyReadPermission(role)));if(!available.length)return <><h1>{t("p0123")}</h1><p>{t("p0002")}</p></>;
 const q=await searchParams,requested=warrantySurface.safeParse(q.surface),surface=requested.success&&available.includes(requested.data)?requested.data:available[0]!;
 const organization=typeof q.organization==="string"&&s.organizations.includes(q.organization)?q.organization:s.organizations[0];if(!organization)return <p>{t("p0003")}</p>;
 const initial=(key:string)=>typeof q[key]==="string"&&warrantyID.safeParse(q[key]).success?q[key] as string:"";
 const scope=createHash("sha256").update(JSON.stringify([s.tenantId,organization,s.subject,surface])).digest("hex");
 return <><h1 className="pageTitle">{t("p0124")}</h1><p>{t("p0005")} {organization}{t("p0120")} {surface==="franchise"?t("p0010"):surface==="customer"?t("p0125"):t("p0121")}{t("p0008")}</p><WarrantyWorkspace key={scope} scope={scope} organization={organization} tenant={s.tenantId} subject={s.subject} permissions={s.permissions} surface={surface} initialClaim={initial("case")} initialQuote={initial("quote")} initialHandover={initial("handover")}/></>;
}
````

### FILE: `src/components/warranty-workspace.tsx`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bfbda5f77dd9c515ecd0e691c338ac1901976037d06c13ca9f6b22baacb0da15"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import {useEffect,useRef,useState,type FormEvent} from "react";
import {z} from "zod";
import {catalogSHA} from "@/platform/catalog/authoring";
import {warrantyID,warrantySHA,warrantyCommand,warrantyOffer,warrantyActivation,warrantyClaim,warrantyStep,warrantyStepReference,warrantyStepMatches,pendingWarrantyStep,warrantyProfileView,warrantyQuote,type WarrantyCommand,type WarrantyClaim,type WarrantyOffer} from "@/platform/warranty/contract";
const endpoint="/api/enterprise/warranty";
const offerPending=z.object({family:z.literal("offer"),action:z.enum(["offer","acknowledge"]),quote_id:warrantyID,quote_version:z.string().regex(/^[1-9][0-9]{0,18}$/u),profile_sha256:warrantySHA,evidence_sha256:warrantySHA.optional()}).strict();
const activationPending=z.object({family:z.literal("activation"),handover_id:warrantyID}).strict();
const pendingSchema=z.discriminatedUnion("family",[pendingWarrantyStep,offerPending,activationPending]);
type Pending=z.infer<typeof pendingSchema>;
type PartInput={line_id:string;item_id:string;bin_id:string;lot_id:string;quantity:string;specific_receipt_entry:string};
const newPart=():PartInput=>({line_id:crypto.randomUUID(),item_id:"",bin_id:"",lot_id:"",quantity:"1",specific_receipt_entry:""});

export function WarrantyWorkspace({scope,organization,subject,tenant,permissions,surface,initialClaim,initialQuote,initialHandover}:{scope:string;organization:string;subject:string;tenant:string;permissions:string[];surface:"franchise"|"customer"|"factory";initialClaim:string;initialQuote:string;initialHandover:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const statusLabel={opened:t("p0603"),diagnosis:t("p0946"),repair:t("p0611"),quality:t("p0947"),closed:t("p0948"),cancelled:t("p0326")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p);
 const held=useRef(false),saved=useRef<Pending|null>(null);const[ready,setReady]=useState(false),[busy,setBusy]=useState(false),[pending,setPending]=useState<Pending|null>(null),[notice,setNotice]=useState("");
 const[id,setID]=useState(initialClaim),[claim,setClaim]=useState<WarrantyClaim|null>(null),[quoteID,setQuoteID]=useState(initialQuote),[quote,setQuote]=useState<z.infer<typeof warrantyQuote>|null>(null),[offer,setOffer]=useState<WarrantyOffer|null>(null),[profile,setProfile]=useState<z.infer<typeof warrantyProfileView>|null>(null),[ack,setAck]=useState(false);
 const[handoverID,setHandoverID]=useState(initialHandover),[activation,setActivation]=useState<z.infer<typeof warrantyActivation>|null>(null),[appointment,setAppointment]=useState(""),[severity,setSeverity]=useState("medium"),[description,setDescription]=useState(""),[fault,setFault]=useState(""),[diagnosis,setDiagnosis]=useState(""),[labor,setLabor]=useState(""),[parts,setParts]=useState<PartInput[]>([]),[reason,setReason]=useState("");
 const[evidence,setEvidence]=useState(""),[correction,setCorrection]=useState("");const key=`elite-warranty:${scope}`;
 function persist(value:Pending|null){if(value){const raw=JSON.stringify(pendingSchema.parse(value));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage")}else sessionStorage.removeItem(key);saved.current=value;setPending(value)}
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const p=pendingSchema.parse(JSON.parse(raw));saved.current=p;setPending(p);if(p.family==="claim")setID(p.case_id);if(p.family==="offer")setQuoteID(p.quote_id);if(p.family==="activation")setHandoverID(p.handover_id)}setReady(true)}catch{setNotice(t("p0835"))}},[key]);
 async function get(kind:string,id?:string,command?:string){const r=await fetch(endpoint+"?"+new URLSearchParams({kind,organization_id:organization,surface,...(id?{id}:{}),...(command?{command_id:command}:{})}),{cache:"no-store",signal:AbortSignal.timeout(10000)});if(!r.ok)throw new Error("unavailable");return r.json() as Promise<unknown>}
 async function readClaim(id:string){const v=warrantyClaim.parse(await get("claim",id));if(v.case_id!==id||(surface==="factory"?v.factory_organization_id:v.organization_id)!==organization||surface==="customer"&&v.customer_subject!==subject)throw new Error("scope");setClaim(v);setID(id)}
 async function consult(kind:"claim"|"quote"|"offer"|"profile"|"activation"){
  if(held.current)return;held.current=true;setBusy(true);if(kind==="claim")setClaim(null);if(kind==="quote")setQuote(null);if(kind==="offer"){setOffer(null);setAck(false)}
  try{
   if(kind==="claim")await readClaim(id);
   if(kind==="quote")setQuote(warrantyQuote.parse(await get(kind,quoteID)));
   if(kind==="offer")setOffer(warrantyOffer.parse(await get(kind,quoteID)));
   if(kind==="profile")setProfile(warrantyProfileView.parse(await get(kind)));
   if(kind==="activation")setActivation(warrantyActivation.parse(await get(kind,handoverID)));
   setNotice(t("p0949"));
  }catch{setNotice(t("p0950"))}finally{held.current=false;setBusy(false)}
 }
 async function reference(c:WarrantyCommand):Promise<Pending>{
  if(c.action==="offer"||c.action==="acknowledge")return offerPending.parse({family:"offer",action:c.action,quote_id:c.quote_id,quote_version:c.action==="offer"?(BigInt(c.quote_version)+1n).toString():c.quote_version,profile_sha256:c.profile_sha256,...(c.action==="acknowledge"?{evidence_sha256:c.evidence_sha256}:{})});
  if(c.action==="activate")return{family:"activation",handover_id:c.handover_id};return warrantyStepReference(c);
 }
 function acceptReceipt(raw:unknown,p:Pending){
  if(p.family==="claim"){const v=warrantyStep.parse(raw);if(!warrantyStepMatches(v,p,subject))throw new Error("identity");persist(null);setClaim(null);setID(v.case_id);return v.case_id}
  if(p.family==="offer"){const v=warrantyOffer.parse(raw);if(v.quote_id!==p.quote_id||v.quote_version!==p.quote_version||v.profile_sha256!==p.profile_sha256||v.profile.tenant_id!==tenant||v.profile.organization_id!==organization||(p.action==="offer"?v.offered_by!==subject:!v.acknowledged||v.customer_subject!==subject||v.evidence_sha256!==p.evidence_sha256))throw new Error("identity");persist(null);setOffer(v);setQuote(null);setAck(false);return null}
  const v=warrantyActivation.parse(raw);if(v.handover_id!==p.handover_id||v.organization_id!==organization||v.activated_by!==subject)throw new Error("identity");persist(null);setActivation(v);return null;
 }
 async function send(input:unknown){
  if(held.current||!ready||saved.current)return;held.current=true;setBusy(true);let started=false;
  try{const c=warrantyCommand.parse(input),p=await reference(c);persist(p);started=true;
   const r=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c),signal:AbortSignal.timeout(10000)});
   if(!r.ok){if([400,401,403,413,415].includes(r.status)){persist(null);started=false}throw new Error("unconfirmed")}
   const caseID=acceptReceipt(await r.json(),p);started=false;setNotice(t("p0838"));if(caseID){try{await readClaim(caseID)}catch{setNotice(t("p0951"))}}
  }catch{setNotice(started?t("p0660"):t("p0718"))}finally{held.current=false;setBusy(false)}
 }
 async function recover(){const p=saved.current;if(!p||held.current)return;held.current=true;setBusy(true);try{const raw=p.family==="claim"?await get("command",p.case_id,p.command_id):p.family==="offer"?await get("offer",p.quote_id):await get("activation",p.handover_id);const caseID=acceptReceipt(raw,p);setNotice(t("p0841"));if(caseID){try{await readClaim(caseID)}catch{setNotice(t("p0952"))}}}catch{setNotice(t("p0843"))}finally{held.current=false;setBusy(false)}}
 async function choose(file:File|null,correct=false){const set=correct?setCorrection:setEvidence;set("");if(!file||file.size<1||file.size>1048576){setNotice(t("p0953"));return}set(await catalogSHA(new Uint8Array(await file.arrayBuffer())))}
 const blocked=busy||!ready||pending!==null;
 const base={organization_id:organization,surface};
 const command=(action:WarrantyCommand["action"],extra:Record<string,unknown>={})=>{if(claim)void send({...base,action,case_id:claim.case_id,command_id:crypto.randomUUID(),expected_version:claim.version,...extra})};
 function open(e:FormEvent){e.preventDefault();void send({...base,action:"open",case_id:crypto.randomUUID(),handover_id:handoverID,appointment_id:appointment,severity,description,evidence_sha256:evidence})}
 function planRepair(e:FormEvent){e.preventDefault();command("plan",{labor_work:labor,parts:parts.map(({lot_id,specific_receipt_entry,...p})=>({...p,...(lot_id?{lot_id}:{}),...(specific_receipt_entry?{specific_receipt_entry}:{})}))})}
 function partChange(i:number,k:keyof PartInput,v:string){setParts(old=>old.map((p,j)=>j===i?{...p,[k]:v}:p))}
 const ctx=claim?.context;
 return <>
 <p role="status" aria-live="polite">{notice}</p>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["warranty-role-view"]}/>
 {pending&&<section className="card" aria-label={t("p0166")}><h2>{t("p0166")}</h2><p>{t("p0167")} <code>{pending.family==="claim"?pending.case_id:pending.family==="offer"?pending.quote_id:pending.handover_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 <section className="card"><h2>{t("p0954")}</h2><label>{t("p0848")}<input type="file" disabled={blocked} onChange={e=>void choose(e.target.files?.[0]??null)}/></label><p>{t("p0955")}</p>{evidence&&<p>{t("p0850")}</p>}<label>{t("p0220")}<textarea value={reason} onChange={e=>setReason(e.target.value)} maxLength={2048} disabled={blocked}/></label></section>
 {surface!=="factory"&&<section className="card"><h2>{t("p0956")}</h2><label>{t("p0957")}<input value={quoteID} onChange={e=>{setQuoteID(e.target.value);setQuote(null);setOffer(null);setAck(false)}} maxLength={120}/></label>
 {surface==="franchise"&&can("warranty:offer")&&<><button className="button" disabled={busy||!quoteID} onClick={()=>void consult("quote")}>{t("p0958")}</button><button className="button" disabled={busy} onClick={()=>void consult("profile")}>{t("p0959")}</button>
 {quote&&<p>{t("p0960")} {quote.total_minor_units} {t("p0209")} {quote.currency}{t("p0008")} {quote.current?t("p0961"):t("p0962")}{t("p0008")}</p>}{profile&&<article><h3>{t("p0963")} {profile.profile.terms_version}</h3><p style={{whiteSpace:"pre-wrap"}}>{profile.profile.terms_text}</p><p>{t("p0964")} {profile.profile.parts_duration_days} {t("p0965")} {profile.profile.labor_duration_days} {t("p0966")}</p></article>}
 <button className="button" disabled={blocked||!quote||!profile||quote.state!=="issued"||!quote.current||!quote.customer_subject} onClick={()=>{if(quote&&profile)void send({...base,action:"offer",quote_id:quote.quote_id,quote_version:quote.quote_version,profile_sha256:profile.profile_sha256})}}>{t("p0967")}</button></>}
 <button className="button" disabled={busy||!quoteID} onClick={()=>void consult("offer")}>{t("p0968")}</button>
 {offer&&<article><h3>{t("p0969")} {offer.profile.terms_version}</h3><p style={{whiteSpace:"pre-wrap"}}>{offer.profile.terms_text}</p><p>{offer.acknowledged?t("p0970"):t("p0971")}</p>{surface==="customer"&&!offer.acknowledged&&<><label><input type="checkbox" checked={ack} onChange={e=>setAck(e.target.checked)} disabled={blocked}/>{t("p0972")}</label><button className="button" disabled={blocked||!ack||!evidence} onClick={()=>void send({...base,action:"acknowledge",quote_id:offer.quote_id,quote_version:offer.quote_version,profile_sha256:offer.profile_sha256,evidence_sha256:evidence})}>{t("p0973")}</button></>}</article>}
 </section>}
 {surface!=="factory"&&<section className="card"><h2>{t("p0974")}</h2><label>{t("p0975")}<input value={handoverID} onChange={e=>{setHandoverID(e.target.value);setActivation(null)}} maxLength={120}/></label><button className="button" disabled={busy||!handoverID} onClick={()=>void consult("activation")}>{t("p0976")}</button>{surface==="franchise"&&can("warranty:activate")&&<button className="button" disabled={blocked||!handoverID} onClick={()=>void send({...base,action:"activate",handover_id:handoverID})}>{t("p0977")}</button>}
 {activation&&<article><h3>{t("p0978")}</h3><p>{t("p0167")} {activation.warranty_id}{t("p0979")} {activation.terms_version}{t("p0008")}</p><p>{t("p0980")} {activation.dates.parts_start} {t("p0748")} {activation.dates.parts_end}{t("p0981")} {activation.dates.labor_start} {t("p0748")} {activation.dates.labor_end}{t("p0008")}</p></article>}
 {(surface==="customer"||can("warranty:request"))&&<form onSubmit={open}><fieldset disabled={blocked}><legend>{t("p0982")}</legend><label>{t("p0983")}<input value={appointment} onChange={e=>setAppointment(e.target.value)} maxLength={120} required/></label><label>{t("p0984")}<select value={severity} onChange={e=>setSeverity(e.target.value)}><option value="low">{t("p0612")}</option><option value="medium">{t("p0985")}</option><option value="high">{t("p0986")}</option><option value="safety">{t("p0987")}</option></select></label><label>{t("p0988")}<textarea value={description} onChange={e=>setDescription(e.target.value)} maxLength={4000} required/></label><button className="button" disabled={!evidence||!handoverID}>{t("p0989")}</button></fieldset></form>}
 </section>}
 <section className="card"><h2>{t("p0990")}</h2><form onSubmit={e=>{e.preventDefault();void consult("claim")}}><label>{t("p0991")}<input value={id} onChange={e=>{setID(e.target.value);setClaim(null)}} maxLength={120} required/></label><button className="button" disabled={busy}>{t("p0990")}</button></form>
 {claim&&ctx&&<><p>{t("p0992")} <strong>{statusLabel[claim.state]}</strong>{t("p0674")} {claim.version}{t("p0008")}</p><p>{ctx.description}</p><p>{t("p0993")} {claim.parts_covered?t("p0994"):t("p0995")}{t("p0996")} {claim.labor_covered?t("p0994"):t("p0995")}{t("p0008")}</p>{ctx.diagnosis&&<article><h3>{t("p0997")}</h3><p>{ctx.diagnosis.fault_code}{t("p0213")} {ctx.diagnosis.description}</p><p>{ctx.diagnosis.excluded?t("p0998"):t("p0999")}</p></article>}
 {surface==="franchise"&&can("warranty:diagnose")&&claim.state==="opened"&&<form onSubmit={e=>{e.preventDefault();command("diagnose",{fault_code:fault,description:diagnosis,evidence_sha256:evidence})}}><fieldset disabled={blocked}><legend>{t("p1000")}</legend><label>{t("p1001")}<input value={fault} onChange={e=>setFault(e.target.value)} maxLength={120} required/></label><label>{t("p1002")}<textarea value={diagnosis} onChange={e=>setDiagnosis(e.target.value)} maxLength={4000} required/></label><button className="button" disabled={!evidence}>{t("p1000")}</button></fieldset></form>}
 {surface==="franchise"&&can("warranty:plan")&&claim.state==="diagnosis"&&!ctx.plan&&<form onSubmit={planRepair}><fieldset disabled={blocked}><legend>{t("p1003")}</legend><label>{t("p1004")}<textarea value={labor} onChange={e=>setLabor(e.target.value)} maxLength={4000}/></label>{parts.map((p,i)=><fieldset key={p.line_id}><legend>{t("p1005")} {i+1}</legend>{([['item_id',t("p1006")],['bin_id',t("p1007")],['lot_id',t("p1008")],['quantity',t("p1009")],['specific_receipt_entry',t("p1010")]] as const).map(([k,label])=><label key={k}>{label}<input value={p[k]} onChange={e=>partChange(i,k,e.target.value)} maxLength={k==="quantity"?6:120} required={k==="item_id"||k==="bin_id"||k==="quantity"}/></label>)}<button type="button" onClick={()=>setParts(old=>old.filter((_,j)=>i!==j))}>{t("p1011")} {i+1}</button></fieldset>)}<button className="button" type="button" disabled={parts.length>=16} onClick={()=>setParts(old=>[...old,newPart()])}>{t("p1012")}</button><button className="button" disabled={parts.length===0&&!labor}>{t("p1003")}</button></fieldset></form>}
 {ctx.plan&&<article><h3>{t("p1013")}</h3><p>{t("p1014")} {ctx.plan.requester}{t("p1015")} {ctx.plan.approval_state==="pending"?t("p0797"):ctx.plan.approval_state==="approved"?t("p0798"):t("p0337")}{t("p0008")}</p><p>{ctx.plan.request.labor_work}</p><ul>{ctx.plan.request.parts.map(p=><li key={p.line_id}>{p.item_id}{t("p0213")} {p.quantity} {t("p0016")} {p.bin_id}{p.lot_id?` · lote ${p.lot_id}`:""}</li>)}</ul>
 {surface==="franchise"&&can("warranty:approve")&&claim.state==="diagnosis"&&ctx.plan.approval_state==="pending"&&<><p>{t("p1016")} {ctx.plan.expires_at}{t("p0008")}</p>{[true,false].map(approved=><button key={String(approved)} className="button" disabled={blocked||reason.trim().length<3||ctx.plan!.requester===subject||approved&&(ctx.plan!.excluded||!claim.parts_covered&&!claim.labor_covered)} onClick={()=>command("decide",{payload_sha256:ctx.plan!.payload_sha256,approved,reason})}>{approved?t("p1017"):t("p1018")}</button>)}</>}
 </article>}
 {surface==="franchise"&&can("warranty:work")&&claim.state==="repair"&&<button className="button" disabled={blocked||!evidence} onClick={()=>command("work",{evidence_sha256:evidence})}>{t("p1019")}</button>}
 {ctx.work&&<p>{t("p1020")} {ctx.work.actor}{t("p0008")}</p>}
 {surface==="franchise"&&can("warranty:quality")&&claim.state==="quality"&&ctx.work&&!ctx.acceptance&&!ctx.quality?.passed&&<article><h3>{t("p1021")}</h3>{ctx.quality&&!ctx.quality.passed&&<><p>{t("p1022")}</p><label>{t("p1023")}<input type="file" disabled={blocked} onChange={e=>void choose(e.target.files?.[0]??null,true)}/></label></>}{[true,false].map(passed=><button className="button" key={String(passed)} disabled={blocked||!evidence||ctx.work!.actor===subject||!!ctx.quality&&!correction} onClick={()=>command("quality",{passed,work_evidence_sha256:ctx.work!.payload_sha256,evidence_sha256:evidence,...(correction?{correction_evidence_sha256:correction}:{})})}>{passed?t("p1024"):t("p1025")}</button>)}</article>}
 {ctx.quality&&<p>{t("p1026")} {ctx.quality.passed?t("p0798"):t("p1027")}{t("p0008")}</p>}
 {surface==="customer"&&claim.state==="quality"&&ctx.quality?.passed&&!ctx.acceptance&&<button className="button" disabled={blocked||!evidence} onClick={()=>command("accept",{quality_sha256:ctx.quality!.payload_sha256,evidence_sha256:evidence})}>{t("p1028")}</button>}
 {ctx.acceptance&&<p>{t("p1029")}</p>}
 {surface==="factory"&&can("warranty:reconcile")&&claim.state==="quality"&&ctx.acceptance&&<button className="button" disabled={blocked||!evidence} onClick={()=>command("reconcile",{acceptance_sha256:ctx.acceptance!.payload_sha256,evidence_sha256:evidence})}>{t("p1030")}</button>}
 {ctx.reconciliation&&<article><h3>{t("p1031")}</h3><p>{t("p1032")} {ctx.reconciliation.recorded_inventory_cost}{t("p0008")}</p><p>{t("p1033")}</p></article>}
 {surface==="franchise"&&can("warranty:cancel")&&["opened","diagnosis","repair"].includes(claim.state)&&<button className="button" disabled={blocked||!evidence||reason.trim().length<3||ctx.plan?.approval_state==="pending"&&ctx.plan.requester===subject} onClick={()=>command("cancel",{reason,evidence_sha256:evidence})}>{t("p1034")}</button>}
 </>}
 </section>
 </>;
}
````

### FILE: `src/platform/warranty/contract.test.ts`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "58773406f1b68241bc11bee70d697a2a7a956c0cf1d3b77ed8bab3d227cfba55"
variables: []
secrets_allowed: false
```

````typescript
import {readFileSync} from "node:fs";
import {it,expect} from "vitest";
import {catalogCanonical,catalogSHA} from "@/platform/catalog/authoring";
import {warrantyCommand,warrantyPayload,warrantyPermission,warrantyPath,warrantyStepReference,warrantyStepMatches} from "./contract";
const cases=JSON.parse(readFileSync("deploy/warranty/role-goldens.json","utf8"))as{name:string;payload:Record<string,unknown>;canonical:string;sha256:string}[];
const role=(a:string)=>["acknowledge","accept","open"].includes(a)?"customer":a==="reconcile"?"factory":"franchise";
it.each(cases)("matches exact Go $name wire without rounding int64",async c=>{const v=warrantyCommand.parse({action:c.name,organization_id:"store",surface:role(c.name),...c.payload});expect(catalogCanonical(warrantyPayload(v))).toBe(c.canonical);expect(await catalogSHA(catalogCanonical(warrantyPayload(v)))).toBe(c.sha256);expect(warrantyPermission(v)).toBeTruthy();expect(warrantyPath(v)).toMatch(/^\/v1\/(franchise|customer|factory)\/warranty\//)});
it("binds open recovery to fixed open command and rejects other identities and versions",async()=>{const c=cases.find(c=>c.name==="open")!,v=warrantyCommand.parse({action:"open",surface:"customer",organization_id:"store",...c.payload}),p=await warrantyStepReference(v);expect(p.command_id).toBe("open");expect(p.version).toBe("1");const r={...p,actor:"customer",state:"opened" as const,payload_sha256:"a".repeat(64),payload:{},recorded_at:"2026-09-12T00:00:00Z",replay:false};expect(warrantyStepMatches(r,p,"customer")).toBe(true);for(const patch of [{actor:"other"},{version:"2"},{request_sha256:"b".repeat(64)},{case_id:"other"}])expect(warrantyStepMatches({...r,...patch},p,"customer")).toBe(false)});
it("does not map staff acceptance or customer work to an authorized write",()=>{for(const[name,surface]of[["accept","franchise"],["work","customer"],["reconcile","customer"]]){const c=cases.find(c=>c.name===name)!,v=warrantyCommand.parse({action:name,surface,organization_id:"store",...c.payload});expect(warrantyPermission(v)).toBeNull()}});
````

### FILE: `src/platform/warranty/contract.ts`

```yaml
block_id: "TS-WARRANTY-ROLE-PORTAL:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8ff31770fe19afd2868bd56c6d027ece48e31351fac7758fa5f2d1570dfc56dc"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED bounded form/receipt correspondence; existing Go owners decide effects.
import {z} from "zod";
import {catalogCanonical,catalogSHA} from "@/platform/catalog/authoring";
export const warrantyID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,119}$/u);
export const warrantySHA=z.string().regex(/^[0-9a-f]{64}$/u);
const positive=z.string().regex(/^[1-9][0-9]{0,18}$/u).refine(v=>BigInt(v)<=9223372036854775807n);
const text=(max:number,min=1)=>z.string().refine(v=>v.trim().length>=min&&new TextEncoder().encode(v).length<=max&&new TextDecoder().decode(new TextEncoder().encode(v))===v);
export const warrantySurface=z.enum(["franchise","customer","factory"]);
const base={organization_id:warrantyID,surface:warrantySurface};
const command={case_id:warrantyID,command_id:warrantyID,expected_version:positive};
const part=z.object({line_id:warrantyID,item_id:warrantyID,bin_id:warrantyID,lot_id:warrantyID.optional(),quantity:z.string().regex(/^[1-9][0-9]{0,5}$/u),specific_receipt_entry:warrantyID.optional()}).strict();
const plan=z.object({...command,labor_work:z.string().refine(v=>new TextEncoder().encode(v).length<=4000),parts:z.array(part).max(16).refine(v=>new Set(v.map(p=>p.line_id)).size===v.length)}).strict().refine(v=>v.parts.length>0||v.labor_work.length>0);
export const warrantyCommand=z.discriminatedUnion("action",[
 z.object({action:z.literal("offer"),...base,quote_id:warrantyID,quote_version:positive,profile_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("acknowledge"),...base,quote_id:warrantyID,quote_version:positive,profile_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("activate"),...base,handover_id:warrantyID}).strict(),
 z.object({action:z.literal("open"),...base,case_id:warrantyID,handover_id:warrantyID,appointment_id:warrantyID,severity:z.enum(["low","medium","high","safety"]),description:text(4000,3),evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("diagnose"),...base,...command,fault_code:warrantyID,description:text(4000,3),evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("plan"),...base,...plan.shape}).strict().refine(v=>v.parts.length>0||v.labor_work.length>0),
 z.object({action:z.literal("decide"),...base,...command,payload_sha256:warrantySHA,approved:z.boolean(),reason:text(2048,3)}).strict(),
 z.object({action:z.literal("work"),...base,...command,evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("quality"),...base,...command,passed:z.boolean(),work_evidence_sha256:warrantySHA,evidence_sha256:warrantySHA,correction_evidence_sha256:warrantySHA.optional()}).strict(),
 z.object({action:z.literal("accept"),...base,...command,quality_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("reconcile"),...base,...command,acceptance_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("cancel"),...base,...command,reason:text(2048,3),evidence_sha256:warrantySHA}).strict()
]);
export type WarrantyCommand=z.infer<typeof warrantyCommand>;
export const warrantyProfile=z.object({schema:z.literal("elite-warranty-profile/v1"),scope:z.literal("MATERIALIZED_PROFILE"),algorithm:z.literal("bc-inclusive-fixed-terms"),algorithm_revision:z.literal(1),tenant_id:z.uuid(),organization_id:warrantyID,factory_organization_id:warrantyID,policy_id:warrantyID,terms_version:warrantyID,terms_text:text(16000),business_time_zone:z.string().min(1).max(128),parts_duration_days:z.number().int().min(1).max(3652500),labor_duration_days:z.number().int().min(1).max(3652500),work_reservation_seconds:z.number().int().min(60).max(86400),fault_exclusions:z.array(warrantyID).max(64),settlement:z.literal("INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT"),authority_reference:warrantyID,decision_reference:warrantyID}).strict();
export const warrantyQuote=z.object({quote_id:warrantyID,organization_id:warrantyID,quote_version:positive,customer_subject:z.string().max(256),state:z.string().min(1).max(64),current:z.boolean(),currency:z.string().regex(/^[A-Z]{3}$/u),total_minor_units:z.string().regex(/^[0-9]{1,19}$/u)}).strict();
export const warrantyProfileView=z.object({profile_sha256:warrantySHA,profile:warrantyProfile}).strict();
export const warrantyOffer=z.object({quote_id:warrantyID,quote_version:positive,customer_subject:z.string().min(1).max(256),profile_sha256:warrantySHA,profile:warrantyProfile,offered_by:z.string().min(1).max(256),offered_at:z.iso.datetime({offset:true}),acknowledged:z.boolean(),evidence_sha256:warrantySHA.optional(),replay:z.boolean()}).strict();
export type WarrantyOffer=z.infer<typeof warrantyOffer>;
export const warrantyActivation=z.object({warranty_id:warrantyID,handover_id:warrantyID,quote_id:warrantyID,order_id:warrantyID,stock_unit_id:warrantyID,organization_id:warrantyID,customer_subject:z.string().min(1).max(256),profile_sha256:warrantySHA,terms_version:warrantyID,dates:z.object({parts_start:z.iso.date(),parts_end:z.iso.date(),labor_start:z.iso.date(),labor_end:z.iso.date()}).strict(),accepted_at:z.iso.datetime({offset:true}),activated_by:z.string().min(1).max(256),activated_at:z.iso.datetime({offset:true}),replay:z.boolean()}).strict();
const kinds=z.enum(["opened","diagnosed","planned","decision","work","quality","accepted","reconciled","cancelled"]);
const states=z.enum(["opened","diagnosis","repair","quality","closed","cancelled"]);
export const warrantyStep=z.object({case_id:warrantyID,command_id:warrantyID,version:positive,kind:kinds,state:states,actor:z.string().min(1).max(256),request_sha256:warrantySHA,payload_sha256:warrantySHA,payload:z.unknown(),recorded_at:z.iso.datetime({offset:true}),replay:z.boolean()}).strict();
export type WarrantyStep=z.infer<typeof warrantyStep>;
const roleContext=z.object({description:text(4000,3),severity:z.enum(["low","medium","high","safety"]),diagnosis:z.object({fault_code:warrantyID,description:text(4000,3),excluded:z.boolean(),evidence_sha256:warrantySHA}).strict().optional(),plan:z.object({request:plan,requester:z.string().min(1).max(256),payload_sha256:warrantySHA,excluded:z.boolean(),expires_at:z.iso.datetime({offset:true}),approval_state:z.enum(["pending","approved","rejected"])}).strict().optional(),work:z.object({actor:z.string().min(1).max(256),payload_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict().optional(),quality:z.object({passed:z.boolean(),payload_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict().optional(),acceptance:z.object({actor:z.string().min(1).max(256),payload_sha256:warrantySHA}).strict().optional(),reconciliation:z.object({settlement:z.literal("INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT"),recorded_inventory_cost:z.string().regex(/^[0-9]+\.[0-9]{4}$/u)}).strict().optional()}).strict();
export const warrantyClaim=z.object({case_id:warrantyID,warranty_id:warrantyID,handover_id:warrantyID,appointment_id:warrantyID,organization_id:warrantyID,factory_organization_id:warrantyID,customer_subject:z.string().min(1).max(256),stock_unit_id:warrantyID,service_date:z.iso.date(),state:states,version:positive,profile_sha256:warrantySHA,parts_covered:z.boolean(),labor_covered:z.boolean(),latest:warrantyStep,context:roleContext}).strict();
export type WarrantyClaim=z.infer<typeof warrantyClaim>;
export const pendingWarrantyStep=z.object({family:z.literal("claim"),case_id:warrantyID,command_id:warrantyID,kind:kinds,version:positive,request_sha256:warrantySHA}).strict();
export type PendingWarrantyStep=z.infer<typeof pendingWarrantyStep>;
export function warrantyPayload(c:WarrantyCommand){const{action:_,organization_id:__,surface:___,...payload}=c;return payload}
export function warrantyPermission(c:WarrantyCommand){
 const action=c.action;
 if(action==="open")return c.surface==="customer"?"warranty:self":c.surface==="franchise"?"warranty:request":null;
 if(action==="acknowledge"||action==="accept")return c.surface==="customer"?"warranty:self":null;
 if(action==="reconcile")return c.surface==="factory"?"warranty:reconcile":null;
 if(c.surface!=="franchise")return null;
 return {offer:"warranty:offer",activate:"warranty:activate",diagnose:"warranty:diagnose",plan:"warranty:plan",decide:"warranty:approve",work:"warranty:work",quality:"warranty:quality",cancel:"warranty:cancel"}[action];
}
export const warrantyReadPermission=(surface:z.infer<typeof warrantySurface>)=>surface==="customer"?"warranty:self":surface==="factory"?"warranty:factory-read":"warranty:read";
export function warrantyPath(c:WarrantyCommand){const base=`/v1/${c.surface}/warranty/`;
 if(c.action==="offer"||c.action==="acknowledge")return base+`quotes/${encodeURIComponent(c.quote_id)}/${c.action==="offer"?"terms":"acknowledgements"}`;
 if(c.action==="activate")return base+`handovers/${encodeURIComponent(c.handover_id)}/activation`;
 if(c.action==="open")return base+"claims";
 return base+`claims/${encodeURIComponent(c.case_id)}/`+{diagnose:"diagnosis",plan:"plan",decide:"decision",work:"work",quality:"quality",accept:"acceptance",reconcile:"reconciliation",cancel:"cancellation"}[c.action];
}
export async function warrantyStepReference(c:WarrantyCommand):Promise<PendingWarrantyStep>{
 if(!("case_id"in c))throw new Error("claim command required");
 const kind={open:"opened",diagnose:"diagnosed",plan:"planned",decide:"decision",work:"work",quality:"quality",accept:"accepted",reconcile:"reconciled",cancel:"cancelled"}[c.action];
 return pendingWarrantyStep.parse({family:"claim",case_id:c.case_id,command_id:c.action==="open"?"open":c.command_id,kind,version:c.action==="open"?"1":(BigInt(c.expected_version)+1n).toString(),request_sha256:await catalogSHA(catalogCanonical(warrantyPayload(c)))});
}
export function warrantyStepMatches(v:WarrantyStep,p:PendingWarrantyStep,subject:string){return v.case_id===p.case_id&&v.command_id===p.command_id&&v.kind===p.kind&&v.version===p.version&&v.actor===subject&&v.request_sha256===p.request_sha256}
````

## 6. Configuration surface

docs/WARRANTY_ROLE_REFERENCE.md. features.warranty_portal opt-in. Role read and explicit action permissions. Evidence hash and read-only predecessor context; up to16parts and32KiB BFF. Pending result reference scoped to tenant/org/actor/view, GET only recovery.

## 7. Dependency bill

14new and7changed AUTHORED unavoidable projection/transport/form/config/fixture glue. No new dependency/runtime; retained BC date adaptation and other official source licenses, no corporate attribution.

## 8. Apply order

Select TS-GO-API-WEB-BRIDGE, TS-OIDC-PORTAL-ADAPTER and TS-FRANCHISE-JOURNEY-PORTALS. Full Go warranty claim and role-view owners required when enabled; narrow web proves imports/types only. Shared command goldens owned here.

## 9. Verification

18web tests,12typed Go/TS goldens, full Next webpack/types. Both actual browser phases pass:18POST,15steps,3claims,1closed/2cancelled,3remaining parts, signed ledger -740.0000 and positive service cost740.0000. Two response losses recovered GET after reload with0extraPOST; review/quality separation, exclusion/rejection/cancellation and foreign/other-customer checks. Screenshot desktop/390px inspected.

## 10. Reconstruction evidence

WARRANTY_ROLE_RELEASE_V402.md/json records exact source/rebuilds and original RED wrapper FAIL859. Only final ledger-sign fixture assertion corrected; existing owned DB verified read-only without replay. Original domain fuzz/host/transaction guards retained, no new algorithm/migration. Does not close wholeT2804 or production.


V402 composed delta: T2804 help CMS and21same-release guides;15oldguides unchanged,5bounded training curricula revision2, no automatic grants. New optional host, shared inline text, actual browser/PG proofs. HELP_CMS_RELEASE_V402.md.

V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.
