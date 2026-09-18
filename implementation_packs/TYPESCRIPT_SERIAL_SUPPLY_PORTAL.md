# Connected serial supply role portal

## 1. Metadata

```yaml
pack_id: "TS-SERIAL-SUPPLY-PORTAL"
pack_version: "0.1.2"
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
CREATE config/supply.role.fixture.json
CREATE deploy/supply/role-goldens.json
CREATE docs/SUPPLY_ROLE_REFERENCE.md
CREATE microsoft_playwright_browser_gate/tests/supply-role-connected.spec.mjs
CREATE src/app/api/enterprise/supply/route.test.ts
CREATE src/app/api/enterprise/supply/route.ts
CREATE src/app/supply/page.tsx
CREATE src/components/supply-workspace.tsx
CREATE src/platform/supply/contract.test.ts
CREATE src/platform/supply/contract.ts
```

## 5. Materialization blocks

### FILE: `config/supply.role.fixture.json`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "45dc860d0ab7a16f1e55089ead136e63365dd36fd7a3624f9d875bd437b5bd85"
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
    "supply_portal": true
  }
}
````

### FILE: `deploy/supply/role-goldens.json`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6b48bcaa97ba8e6270defbbe413efc44e920d6030fc51c26e20c6d62a2ba6e0c"
variables: []
secrets_allowed: false
```

````json
[
  {
    "name": "create",
    "payload": {
      "purchase_order_id": "golden-po",
      "command_id": "golden-command",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "factory_organization_id": "factory",
      "demand_reference": "Demanda á <&> 
 interior",
      "policy_code": "strict-serial-reference/v1",
      "lines": [
        {
          "id": "line",
          "variant_id": "variant",
          "quantity": 2
        }
      ],
      "destination_organization_id": "store",
      "supplier_id": "supplier",
      "currency": "ARS",
      "total_minor_units": "9007199254740993"
    },
    "canonical": "{\"command_id\":\"golden-command\",\"currency\":\"ARS\",\"demand_reference\":\"Demanda á \\u003c\\u0026\\u003e \\u2028 interior\",\"destination_organization_id\":\"store\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"factory_organization_id\":\"factory\",\"lines\":[{\"id\":\"line\",\"quantity\":2,\"variant_id\":\"variant\"}],\"policy_code\":\"strict-serial-reference/v1\",\"purchase_order_id\":\"golden-po\",\"supplier_id\":\"supplier\",\"total_minor_units\":\"9007199254740993\"}",
    "sha256": "0a730937327a7b0cd4c205aa3245880dc1445f05d90865d049518c83b5ed412d"
  },
  {
    "name": "register",
    "payload": {
      "purchase_order_id": "golden-po",
      "command_id": "golden-command",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "kind": "register",
      "expected_version": "9007199254740993",
      "line_id": "line",
      "serial_number": "SERIE-Ñ",
      "vin": "VIN-1"
    },
    "canonical": "{\"command_id\":\"golden-command\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"expected_version\":\"9007199254740993\",\"kind\":\"register\",\"line_id\":\"line\",\"purchase_order_id\":\"golden-po\",\"serial_number\":\"SERIE-Ñ\",\"vin\":\"VIN-1\"}",
    "sha256": "e268f1fdc3b9aa1cd6aaf155a698d0284d8d78f80763dda2aaf908e37bb6c820"
  },
  {
    "name": "milestone",
    "payload": {
      "purchase_order_id": "golden-po",
      "command_id": "golden-command",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "kind": "milestone",
      "expected_version": "6",
      "unit_id": "unit",
      "target_state": "released",
      "reason": "Revisé <archivo> & contenido 
 interior"
    },
    "canonical": "{\"command_id\":\"golden-command\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"expected_version\":\"6\",\"kind\":\"milestone\",\"purchase_order_id\":\"golden-po\",\"reason\":\"Revisé \\u003carchivo\\u003e \\u0026 contenido \\u2029 interior\",\"target_state\":\"released\",\"unit_id\":\"unit\"}",
    "sha256": "822e177bcec096d0661b659a7decf3a39dcb50229bf06093b80126f170784195"
  },
  {
    "name": "ship",
    "payload": {
      "purchase_order_id": "golden-po",
      "command_id": "golden-command",
      "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "kind": "ship",
      "expected_version": "7",
      "shipment_id": "asn",
      "units": [
        "unit-2",
        "unit-1"
      ]
    },
    "canonical": "{\"command_id\":\"golden-command\",\"evidence_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"expected_version\":\"7\",\"kind\":\"ship\",\"purchase_order_id\":\"golden-po\",\"shipment_id\":\"asn\",\"units\":[\"unit-2\",\"unit-1\"]}",
    "sha256": "19ab1a7397f108f6b16b43bd43127cb8626fe13cb2107ca385bdeb0a2ff55108"
  }
]
````

### FILE: `docs/SUPPLY_ROLE_REFERENCE.md`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "07c77a0067497de2ac97dfc90a15ea5533d309e0b43fa6732e10550d75158c86"
variables: []
secrets_allowed: false
```

````markdown
# Suministro por rol — referencia de infraestructura

/supply conecta compras, fábrica, recepción y calidad con el owner J2 existente.
features.supply_portal es opt-in; requiere supply:read o supply:factory-read.
Las escrituras requieren el permiso explícito supply:plan, :factory, :receive,
:release o :inspect correspondiente, con tenant/organización y versión en Go.
El nombre de rol no concede permisos. La referencia del pedido se comparte
sólo con personas autorizadas. Maestros de proveedor/organizaciones/variantes
activas son precondiciones; el formulario crea la orden y plan de cantidades.

El alta usa Operations.CreatePurchaseOrder y BindPlan en una transacción,
SQL y reglas originales conservadas. ID de orden preasignado para recuperación,
eventos con UUID independiente. El lock transaccional por tenant/orden serializa
replay; una colisión del hash sólo serializa, no une IDs ni concede acceso.
El receipt planned existente liga actor/orden/comando/hash; no segundo ledger.
El importe y moneda son explícitos, sin política financiera nueva. Int64 se
transporta como texto. Admite32líneas y1000unidades; páginas de100series.

El formulario registra series y milestones, decisión humana distinta de la
solicitud, rechazo y reemplazo, ASN, recepción parcial en cuarentena, rechazo,
reinspección con evidencia nueva y liberación. El view proyecta el último
estado de revisión/requester para deshabilitar autoaprobación; Go y sus
constraints siguen siendo la autoridad. No cambia la máquina de estados.

El navegador guarda únicamente orden/comando/tipo/versión/hash en sessionStorage,
aislado por tenant/org/subject/vista. Un resultado incierto bloquea nuevas
escrituras hasta recuperar un receipt concordante mediante GET. Los archivos
de evidencia y motivos no se guardan en el marcador; se registra la huella,
no se afirma almacenar o clasificar documentos. BFF limita32KiB y tiempo,
autoriza antes de leer body y usa rutas backend fijas. Una respuesta409 incierta
se conserva para revisión autorizada; nunca se transforma en reintento automático.

La ayuda inline de esta versión explica el recorrido y su recuperación.
Nuevo nav.supply tiene texto español/inglés. La localización de todo el portal,
el CMS/general help y J5 siguen dentro de T2804; este delta no cierra todoT2804.

Pruebas nuevas: TestSerialSupplyCreateAtomic demuestra1nuevo+3replay concurrentes,
rollback de5tablas por variante inválida y fallo tardío de outbox, actor/hash/org,
lectura durable e int64 exacto. TestSupplyRoleCommandGoldens compara4comandos
Go/TS con Unicode/HTML/int64/orden de arrays. FuzzSupplyCreateBounds usa un
presupuesto finito de2s sobre source closure aislado. No nuevo runtime/migración.

TestSupplyRoleBrowser con ELITE_SUPPLY_ROLE_BROWSER=1, ELITE_WEB_ROOT y
ELITE_NODE_BIN absolutos, Next compilado--webpack y PAYMENT_CONNECTED_DB_URL
aislada usa Next/BFF/Go/PG reales y JWE/RS256/JWKS sintéticos. Sólo maestros
se siembran: sus23POST crean la orden/plan,3series(1rechazada),2stocks disponibles,
6solicitudes y6decisiones de calidad. Pérdidas de respuesta de alta y recepción
se recuperan porGET tras recargar, sinPOST adicional. Incluye permisos, foreign
org, autoaprobación deshabilitada, escritorio y390px sin desborde.

Seleccionar GO-CONNECTED-SERIAL-SUPPLY junto con su nuevo companion de alta
y TS-SERIAL-SUPPLY-PORTAL. El test Go de goldens usa el fixture del pack web;
el browser reutiliza el issuer sintético de GO-CONNECTED-CATALOG-AUTHORING,
seleccionado en la composición completa. El perfil web aislado prueba tipos,
no conectividad Go. Evidencias y reconstrucción en SUPPLY_ROLE_RELEASE_V402.

AUTHORED sólo composición/transporte/formularios/fixtures; sin nuevo upstream,
sin atribución corporativa. No aceptación productiva ni certificación de calidad.
ARCA penúltimo; Daybreak/libxml2 último diferido sin investigación.
````

### FILE: `microsoft_playwright_browser_gate/tests/supply-role-connected.spec.mjs`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f1c2de0c6ef70d81cb1ace5b911f8e9dda2f030ae3d8f19477a219b483855529"
variables: []
secrets_allowed: false
```

````text
import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
test('role supply order to factory replacement, partial receiving, reinspection and available stock',async({page,context},info)=>{
 if(process.env.ELITE_SUPPLY_ROLE_BROWSER!=='1')throw new Error('explicit fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');page.setDefaultTimeout(12000);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href),identities=JSON.parse(process.env.ELITE_SUPPLY_IDENTITIES);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[],posts=[];page.on('pageerror',e=>errors.push(e.message));page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/supply'))posts.push(r.url())});
 async function evidence(suffix='initial'){await page.getByLabel('Archivo de evidencia',{exact:true}).setInputFiles({name:'fixture.txt',mimeType:'text/plain',buffer:Buffer.from('Authorized synthetic evidence '+suffix)});await expect(page.getByText('Evidencia preparada.',{exact:true})).toBeVisible();await page.getByRole('textbox',{name:'Motivo de revisión',exact:true}).fill('Revisión explícita de la evidencia sintética '+suffix)}
 async function open(name,order=''){await identity(name);await page.goto('/supply'+(order?'?order='+encodeURIComponent(order):''));await expect(page.getByRole('heading',{name:'Suministro por serie',exact:true})).toBeVisible();if(order){await page.getByRole('button',{name:'Consultar orden',exact:true}).click();await expect(page.getByRole('status')).toContainText('Orden y versión consultadas')}await evidence(name)}
 async function click(name,within=page){await within.getByRole('button',{name,exact:true}).click();await expect(page.getByRole('status')).toContainText('Operación registrada');await expect(page.getByRole('button',{name:'Actualizar desde primera página',exact:true})).toBeEnabled()}
 function unit(serial){return page.getByRole('article',{name:'Serie '+serial,exact:true})}
 async function drop(kind){let dropped=false;await page.route('**/api/enterprise/supply**',async route=>{const r=route.request();if(dropped||r.method()!=='POST'||r.postDataJSON().kind!==kind){await route.continue();return}dropped=true;const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')})}
 async function recover(){await expect(page.getByRole('status')).toContainText('Resultado sin confirmar');const n=posts.length;await page.unrouteAll();await page.reload();await page.getByRole('button',{name:'Consultar resultado pendiente',exact:true}).click();await expect(page.getByRole('status')).toContainText('Resultado recuperado');expect(posts.length).toBe(n)}
 await identity('unprivileged');await page.goto('/supply');await expect(page.getByText('No tenés permiso para consultar este espacio.',{exact:true})).toBeVisible();
 await identity('reader');expect((await context.request.post(base+'/api/enterprise/supply',{headers:{origin:base,'content-type':'application/json'},data:{}})).status()).toBe(403);
 await open('buyer');
 for(const [name,value]of [['Proveedor','supplier'],['Organización de fábrica','factory'],['Moneda','ARS'],['Importe total en unidades menores','90000'],['Referencia de la demanda','Demanda de dos unidades'],['Variante','variant']])await page.getByRole('textbox',{name,exact:true}).fill(value);
 await page.getByRole('spinbutton',{name:'Cantidad',exact:true}).fill('2');await drop('create');await page.getByRole('button',{name:'Crear orden con cantidades',exact:true}).click();await recover();
 const order=await page.getByRole('textbox',{name:'Referencia de la orden',exact:true}).inputValue();expect(order).toMatch(/^[a-f0-9-]{36}$/);await evidence('buyer-confirm');await click('Enviar a fábrica');
 await open('maker',order);await click('Confirmar orden en fábrica');await click('Iniciar producción');
 async function make(serial){await page.getByRole('textbox',{name:'Número de serie',exact:true}).fill(serial);await click('Registrar serie');await click('Registrar montaje',unit(serial));await click('Solicitar revisión de fábrica',unit(serial));await expect(unit(serial).getByRole('button',{name:'Liberar serie en fábrica',exact:true})).toBeDisabled()}
 await make('ROLE-A');await make('ROLE-B');await open('factory-qa',order);await click('Liberar serie en fábrica',unit('ROLE-A'));await click('Rechazar serie en fábrica',unit('ROLE-B'));
 await open('maker',order);await make('ROLE-C');await open('factory-qa',order);await click('Liberar serie en fábrica',unit('ROLE-C'));
 await open('maker',order);for(const serial of ['ROLE-A','ROLE-C'])await unit(serial).getByRole('checkbox',{name:'Seleccionar '+serial,exact:true}).check();await page.getByRole('textbox',{name:'Referencia del envío',exact:true}).fill('ROLE-ASN');await click('Registrar despacho');
 await open('receiver',order);await unit('ROLE-A').getByRole('checkbox',{name:'Seleccionar ROLE-A',exact:true}).check();await page.getByRole('textbox',{name:'Referencia del envío',exact:true}).fill('ROLE-ASN');await drop('receive');await page.getByRole('button',{name:'Registrar recepción en cuarentena',exact:true}).click();await recover();await evidence('receiver');await expect(unit('ROLE-A')).toContainText('Stock: En cuarentena');await expect(unit('ROLE-A').getByRole('button',{name:'Liberar stock recibido',exact:true})).toBeDisabled();
 await open('receipt-qa',order);await click('Rechazar calidad recibida',unit('ROLE-A'));await expect(unit('ROLE-A')).toContainText('Revisión de recepción: Rechazada');await expect(unit('ROLE-A').getByRole('button',{name:'Liberar stock recibido',exact:true})).toBeDisabled();
 await open('receiver',order);await evidence('new-reinspection');await click('Solicitar reinspección',unit('ROLE-A'));await expect(unit('ROLE-A')).toContainText('Revisión de recepción: Pendiente');
 await open('receipt-qa',order);await click('Liberar stock recibido',unit('ROLE-A'));await expect(unit('ROLE-A')).toContainText('Stock: Disponible');
 await open('receiver',order);await unit('ROLE-C').getByRole('checkbox',{name:'Seleccionar ROLE-C',exact:true}).check();await page.getByRole('textbox',{name:'Referencia del envío',exact:true}).fill('ROLE-ASN');await click('Registrar recepción en cuarentena');
 await open('receipt-qa',order);await click('Liberar stock recibido',unit('ROLE-C'));await expect(unit('ROLE-A')).toContainText('Stock: Disponible');await expect(unit('ROLE-C')).toContainText('Stock: Disponible');await expect(page.getByText('Recibida',{exact:true})).toBeVisible();
 await page.screenshot({path:info.outputPath('supply-role-available-desktop.png'),fullPage:true});await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('supply-role-mobile.png'),fullPage:true});
 await identity('foreign');expect((await context.request.get(base+'/api/enterprise/supply?'+new URLSearchParams({organization_id:'foreign',purchase_order_id:order,surface:'franchise'}))).status()).toBe(409);
 expect(posts).toHaveLength(23);expect(errors).toEqual([]);
});
````

### FILE: `src/app/api/enterprise/supply/route.test.ts`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b00c6125c20b170a67870696bb70f47fca8d0364ed7f27bcdc00aa8157f82cda"
variables: []
secrets_allowed: false
```

````typescript
import {NextRequest} from "next/server";
import {it,expect,vi,afterEach} from "vitest";
const m=vi.hoisted(()=>({get:vi.fn(),post:vi.fn(),session:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:async()=>({features:{supply_portal:true}})}));
import {GET,POST} from "./route";
afterEach(()=>vi.clearAllMocks());
const req=(body:BodyInit)=>new NextRequest("https://portal.example.test/api/enterprise/supply",{method:"POST",body,headers:{origin:"https://portal.example.test","content-type":"application/json"},duplex:"half" as const});
it("authenticates before consuming an unbounded body",async()=>{m.session.mockResolvedValue(null);const r=req("{}");expect((await POST(r)).status).toBe(401);expect(r.bodyUsed).toBe(false);expect(m.post).not.toHaveBeenCalled()});
it("cancels oversized chunked bodies before downstream",async()=>{m.session.mockResolvedValue({subject:"buyer",permissions:["supply:plan","supply:read"],organizations:["store"]});let cancelled=false;const stream=new ReadableStream<Uint8Array>({pull(c){c.enqueue(new Uint8Array(32769))},cancel(){cancelled=true}});expect((await POST(req(stream))).status).toBe(413);expect(cancelled).toBe(true);expect(m.post).not.toHaveBeenCalled()});
it("does not let a factory writer use franchise operations or foreign organization",async()=>{m.session.mockResolvedValue({subject:"maker",permissions:["supply:factory","supply:factory-read"],organizations:["factory"]});const c={kind:"submit",organization_id:"factory",purchase_order_id:"po",command_id:"cmd",evidence_sha256:"a".repeat(64),expected_version:"1"};expect((await POST(req(JSON.stringify(c)))).status).toBe(403);expect(m.post).not.toHaveBeenCalled()});
it("rejects duplicate page/scope inputs and never recovers another actor's command",async()=>{m.session.mockResolvedValue({subject:"buyer",permissions:["supply:read"],organizations:["store"]});const base="https://portal.example.test/api/enterprise/supply?organization_id=store&purchase_order_id=po&surface=franchise";expect((await GET(new NextRequest(base+"&after_unit=a&after_unit=b"))).status).toBe(403);expect(m.get).not.toHaveBeenCalled();m.get.mockResolvedValue({purchase_order_id:"po",command_id:"cmd",kind:"planned",version:"1",actor:"other",request_sha256:"a".repeat(64),payload_sha256:"b".repeat(64),payload:{},recorded_at:"2026-09-12T00:00:00Z",replay:false});expect((await GET(new NextRequest(base+"&command_id=cmd"))).status).toBe(409);expect(m.post).not.toHaveBeenCalled()});
````

### FILE: `src/app/api/enterprise/supply/route.ts`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1edd12ad9c2df81046f17d648f0d709c787c8c70bb20ce4545d21a0e9d15e784"
variables: []
secrets_allowed: false
```

````typescript
import {NextResponse,type NextRequest} from "next/server";
import {allowed,readSession} from "@/platform/auth/session";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {protectedGet,protectedPost} from "@/platform/backend/protected-client";
import {boundedCommandBody} from "@/platform/backend/bounded-command";
import {loadBusinessConfig} from "@/platform/config/load";
import {supplyID,supplyCommand,supplyReceipt,supplyPlan,supplyPermission,supplySurface,supplyPayload,supplyReference,supplyMatches} from "@/platform/supply/contract";
const fail=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
const ok=(value:unknown)=>NextResponse.json(value,{headers:{"cache-control":"no-store"}});
const permissions=["supply:plan","supply:factory","supply:receive","supply:release","supply:inspect"];
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.supply_portal!==true)return fail("NOT_FOUND",404);
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);
 const q=r.nextUrl.searchParams,org=q.get("organization_id"),id=q.get("purchase_order_id"),surface=q.get("surface"),command=q.get("command_id"),after=q.get("after_unit");
 const keys=["organization_id","purchase_order_id","surface",...(command!==null?["command_id"]:[]),...(after!==null?["after_unit"]:[])];
 if(command!==null&&after!==null||q.size!==keys.length||keys.some(k=>q.getAll(k).length!==1)||!supplyID.safeParse(org).success||!s.organizations.includes(org!)||!supplyID.safeParse(id).success||!["factory","franchise"].includes(surface??"")||command!==null&&!supplyID.safeParse(command).success||after!==null&&!supplyID.safeParse(after).success)return fail("INVALID_SCOPE",403);
 if(!allowed(s,surface==="factory"?"supply:factory-read":"supply:read"))return fail("FORBIDDEN",403);
 try{
  const raw=await protectedGet(s,`/v1/${surface}/supply/orders/${encodeURIComponent(id!)}`,{organization_id:org!,command_id:command??undefined,after_unit:after??undefined});
  if(command!==null){const v=supplyReceipt.parse(raw);if(v.command_id!==command||v.purchase_order_id!==id||v.actor!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v)}
  const v=supplyPlan.parse(raw);if(v.purchase_order_id!==id||(surface==="factory"?v.factory_organization_id:v.destination_organization_id)!==org)return fail("SCOPE_MISMATCH",409);return ok(v);
 }catch{return fail("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(r:NextRequest){
 if((await loadBusinessConfig()).features.supply_portal!==true)return fail("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return fail("CROSS_ORIGIN_REJECTED",403)}catch{return fail("UNAVAILABLE",503)}
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);
 if(!permissions.some(p=>allowed(s,p)))return fail("FORBIDDEN",403);
 if(r.headers.get("content-type")!=="application/json")return fail("UNSUPPORTED_MEDIA_TYPE",415);
 if(r.nextUrl.searchParams.size)return fail("INVALID_COMMAND",400);
 try{
  let raw:unknown;try{raw=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return fail("INVALID_COMMAND",400)}
  const parsed=supplyCommand.safeParse(raw);if(!parsed.success)return fail("INVALID_COMMAND",400);const c=parsed.data,surface=supplySurface(c.kind);
  if(!s.organizations.includes(c.organization_id)||!allowed(s,supplyPermission(c.kind))||!allowed(s,surface==="factory"?"supply:factory-read":"supply:read"))return fail("FORBIDDEN",403);
  const expected=await supplyReference(c);
  const v=supplyReceipt.parse(await protectedPost(s,`/v1/${surface}/supply/orders/${encodeURIComponent(c.purchase_order_id)}/${c.kind}?${new URLSearchParams({organization_id:c.organization_id})}`,supplyPayload(c)));
  if(!supplyMatches(v,expected,s.subject))return fail("UNCONFIRMED",409);return ok(v);
 }catch(e){return fail(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
````

### FILE: `src/app/supply/page.tsx`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1666641d111c066c93cb5bdd2dc81a75e4995e2652d3363355423dd9ccc2a962"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import {createHash} from "node:crypto";
import {redirect,notFound} from "next/navigation";
import type {Route} from "next";
import {allowed,readSession} from "@/platform/auth/session";
import {loadBusinessConfig} from "@/platform/config/load";
import {supplyID} from "@/platform/supply/contract";
import {SupplyWorkspace} from "@/components/supply-workspace";
export default async function SupplyPage({searchParams}:{searchParams:Promise<Record<string,string|string[]|undefined>>}){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.supply_portal!==true)notFound();
 const s=await readSession();if(!s)redirect("/api/auth/login?return_to=/supply" as Route);
 if(!allowed(s,"supply:read")&&!allowed(s,"supply:factory-read"))return <><h1>{t("p0118")}</h1><p>{t("p0002")}</p></>;
 const q=await searchParams,organization=typeof q.organization==="string"&&s.organizations.includes(q.organization)?q.organization:s.organizations[0];
 if(!organization)return <p>{t("p0003")}</p>;
 const surface=q.surface==="factory"&&allowed(s,"supply:factory-read")?"factory":allowed(s,"supply:read")?"franchise":"factory";
 const initial=typeof q.order==="string"&&supplyID.safeParse(q.order).success?q.order:"";
 const scope=createHash("sha256").update(JSON.stringify([s.tenantId,organization,s.subject,surface])).digest("hex");
 return <><h1 className="pageTitle">{t("p0119")}</h1><p>{t("p0005")} {organization}{t("p0120")} {surface==="factory"?t("p0121"):t("p0122")}{t("p0008")}</p><SupplyWorkspace key={scope} scope={scope} organization={organization} subject={s.subject} permissions={s.permissions} surface={surface} initial={initial}/></>;
}
````

### FILE: `src/components/supply-workspace.tsx`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8e2e99a9be945faeee09e639344935ce7f2b3830c09ffc9c515ce94d414439d8"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import {useEffect,useRef,useState,type FormEvent} from "react";
import {supplyCommand,supplyReceipt,supplyPlan,pendingSupply,supplyDigest,supplyReference,supplyMatches,supplyPolicy,type SupplyCommand,type SupplyPlan,type PendingSupply} from "@/platform/supply/contract";
const endpoint="/api/enterprise/supply";

export function SupplyWorkspace({organization,subject,permissions,scope,surface,initial}:{organization:string;subject:string;permissions:string[];scope:string;surface:"factory"|"franchise";initial:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const labels:Record<string,string>={draft:t("p0648"),submitted:t("p0822"),accepted:t("p0823"),"in-production":t("p0824"),shipped:t("p0825"),received:t("p0826"),cancelled:t("p0056"),planned:t("p0827"),assembly:t("p0828"),quality:t("p0829"),released:t("p0830"),rejected:t("p0337"),"in-transit":t("p0831"),quarantine:t("p0832"),available:t("p0833"),retired:t("p0834")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p);
 const held=useRef(false),saved=useRef<PendingSupply|null>(null);
 const [ready,setReady]=useState(false),[busy,setBusy]=useState(false),[pending,setPending]=useState<PendingSupply|null>(null),[notice,setNotice]=useState("");
 const [id,setID]=useState(initial),[plan,setPlan]=useState<SupplyPlan|null>(null),[evidence,setEvidence]=useState(""),[reason,setReason]=useState("");
 const [supplier,setSupplier]=useState(""),[factory,setFactory]=useState(""),[currency,setCurrency]=useState(""),[total,setTotal]=useState(""),[demand,setDemand]=useState("");
 const [lines,setLines]=useState([{id:"line-1",variant_id:"",quantity:1}]),[line,setLine]=useState(""),[serial,setSerial]=useState(""),[vin,setVIN]=useState(""),[battery,setBattery]=useState("");
 const [selected,setSelected]=useState<string[]>([]),[shipment,setShipment]=useState("");
 const key=`elite-supply:${scope}`;
 function persist(value:PendingSupply|null){if(value){const raw=JSON.stringify(pendingSupply.parse(value));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage")}else sessionStorage.removeItem(key);saved.current=value;setPending(value)}
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const v=pendingSupply.parse(JSON.parse(raw));saved.current=v;setPending(v);setID(v.purchase_order_id)}setReady(true)}catch{setNotice(t("p0835"))}},[key]);
 async function get(order:string,extra:Record<string,string>={}){const r=await fetch(endpoint+"?"+new URLSearchParams({organization_id:organization,purchase_order_id:order,surface,...extra}),{cache:"no-store",signal:AbortSignal.timeout(10000)});if(!r.ok)throw new Error("unavailable");return r.json() as Promise<unknown>}
 async function readPlan(order:string,after=""){const v=supplyPlan.parse(await get(order,after?{after_unit:after}:{}));if(v.purchase_order_id!==order||(surface==="factory"?v.factory_organization_id:v.destination_organization_id)!==organization)throw new Error("scope");setPlan(v);setID(order);setSelected([]);setLine(v.lines[0]?.id??"")}
 async function consult(after=""){if(held.current)return;held.current=true;setBusy(true);setPlan(null);try{await readPlan(id,after);setNotice(t("p0836"))}catch{setNotice(t("p0837"))}finally{held.current=false;setBusy(false)}}
 async function send(input:unknown){
  if(held.current||!ready||saved.current)return;held.current=true;setBusy(true);let started=false;
  try{
   const c=supplyCommand.parse(input),ref=await supplyReference(c);persist(ref);started=true;
   const response=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c),signal:AbortSignal.timeout(10000)});
   if(!response.ok){if([400,401,403,413,415].includes(response.status)){persist(null);started=false}throw new Error("unconfirmed")}
   const value=supplyReceipt.parse(await response.json());if(!supplyMatches(value,ref,subject))throw new Error("identity");
   persist(null);started=false;setID(value.purchase_order_id);setPlan(null);setNotice(t("p0838"));
   try{await readPlan(value.purchase_order_id)}catch{setNotice(t("p0839"))}
  }catch{setNotice(started?t("p0660"):t("p0840"))}
  finally{held.current=false;setBusy(false)}
 }
 async function recover(){const p=saved.current;if(!p||held.current)return;held.current=true;setBusy(true);try{const value=supplyReceipt.parse(await get(p.purchase_order_id,{command_id:p.command_id}));if(!supplyMatches(value,p,subject))throw new Error("identity");persist(null);setPlan(null);setID(value.purchase_order_id);setNotice(t("p0841"));try{await readPlan(value.purchase_order_id)}catch{setNotice(t("p0842"))}}catch{setNotice(t("p0843"))}finally{held.current=false;setBusy(false)}}
 async function selectEvidence(file:File|null){setEvidence("");if(!file||file.size<1||file.size>1048576){setNotice(t("p0844"));return}setEvidence(await supplyDigest(new Uint8Array(await file.arrayBuffer())))}
 const blocked=busy||!ready||pending!==null;
 const command=(kind:SupplyCommand["kind"],extra:Record<string,unknown>={})=>{if(plan)void send({kind,organization_id:organization,purchase_order_id:plan.purchase_order_id,command_id:crypto.randomUUID(),evidence_sha256:evidence,expected_version:plan.version,...extra})};
 function create(e:FormEvent){e.preventDefault();void send({kind:"create",organization_id:organization,destination_organization_id:organization,purchase_order_id:crypto.randomUUID(),command_id:crypto.randomUUID(),supplier_id:supplier,factory_organization_id:factory,currency,total_minor_units:total,demand_reference:demand,policy_code:supplyPolicy,evidence_sha256:evidence,lines})}
 return <>
 <p role="status" aria-live="polite">{notice}</p>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["supply-role-view"]}/>
 {pending&&<section className="card" aria-label={t("p0166")}><h2>{t("p0166")}</h2><p>{t("p0845")} <code>{pending.purchase_order_id}</code></p><p>{t("p0846")} <code>{pending.command_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 <section className="card"><h2>{t("p0847")}</h2><label>{t("p0848")}<input type="file" disabled={blocked} onChange={e=>void selectEvidence(e.target.files?.[0]??null)}/></label><p>{t("p0849")}</p>{evidence&&<p>{t("p0850")}</p>}<label>{t("p0851")}<textarea value={reason} onChange={e=>setReason(e.target.value)} maxLength={2000} disabled={blocked}/></label></section>
 {surface==="franchise"&&can("supply:plan")&&<section className="card"><h2>{t("p0852")}</h2><form onSubmit={create}><fieldset disabled={blocked}>
 <label>{t("p0853")}<input value={supplier} onChange={e=>setSupplier(e.target.value)} maxLength={128} required/></label>
 <label>{t("p0854")}<input value={factory} onChange={e=>setFactory(e.target.value)} maxLength={128} required/></label>
 <label>{t("p0855")}<input value={currency} onChange={e=>setCurrency(e.target.value)} pattern="[A-Z]{3}" maxLength={3} required/></label>
 <label>{t("p0856")}<input value={total} onChange={e=>setTotal(e.target.value)} inputMode="numeric" pattern="[0-9]+" maxLength={19} required/></label>
 <label>{t("p0857")}<input value={demand} onChange={e=>setDemand(e.target.value)} maxLength={1024} required/></label>
 {lines.map((v,i)=><fieldset key={v.id}><legend>{t("p0858")} {i+1}</legend><label>{t("p0182")}<input value={v.variant_id} onChange={e=>setLines(old=>old.map((x,j)=>i===j?{...x,variant_id:e.target.value}:x))} required maxLength={128}/></label><label>{t("p0859")}<input type="number" min={1} max={1000} step={1} value={v.quantity} onChange={e=>setLines(old=>old.map((x,j)=>i===j?{...x,quantity:Number(e.target.value)}:x))} required/></label>{lines.length>1&&<button type="button" onClick={()=>setLines(old=>old.filter((_,j)=>j!==i))}>{t("p0860")} {i+1}</button>}</fieldset>)}
 <button type="button" className="button" disabled={lines.length>=32} onClick={()=>setLines(old=>[...old,{id:crypto.randomUUID(),variant_id:"",quantity:1}])}>{t("p0861")}</button>
 <button className="button" disabled={!evidence}>{t("p0862")}</button></fieldset></form></section>}
 <section className="card"><h2>{t("p0863")}</h2><form onSubmit={e=>{e.preventDefault();void consult()}}><label>{t("p0864")}<input value={id} onChange={e=>{setID(e.target.value);setPlan(null)}} required maxLength={128}/></label><button className="button" disabled={busy}>{t("p0863")}</button></form>
 {plan&&<><p>{t("p0865")} <strong>{labels[plan.state]}</strong>{t("p0674")} {plan.version}{t("p0008")}</p><p>{t("p0866")} {plan.supplier_id}{t("p0867")} {plan.factory_organization_id}{t("p0868")} {plan.destination_organization_id}{t("p0008")}</p><p>{t("p0275")} {plan.total_minor_units} {t("p0209")} {plan.currency}{t("p0869")} {plan.demand_reference}{t("p0008")}</p><p><a href={`/supply?order=${encodeURIComponent(plan.purchase_order_id)}&surface=factory`}>{t("p0870")}</a></p><ul>{plan.lines.map(l=><li key={l.id}>{l.variant_id}{t("p0213")} {l.quantity} {t("p0871")} {l.id}</li>)}</ul>
 {surface==="franchise"&&can("supply:plan")&&<><button className="button" disabled={blocked||!evidence||plan.state!=="draft"} onClick={()=>command("submit")}>{t("p0872")}</button><button className="button" disabled={blocked||!evidence||!["draft","submitted","accepted"].includes(plan.state)} onClick={()=>command("cancel")}>{t("p0873")}</button></>}
 {surface==="factory"&&can("supply:factory")&&<><button className="button" disabled={blocked||!evidence||plan.state!=="submitted"} onClick={()=>command("confirm")}>{t("p0874")}</button><button className="button" disabled={blocked||!evidence||plan.state!=="accepted"} onClick={()=>command("start")}>{t("p0875")}</button>
 <form onSubmit={e=>{e.preventDefault();command("register",{line_id:line,serial_number:serial,...(vin?{vin}:{}),...(battery?{battery_serial_number:battery}:{})});}}><fieldset disabled={blocked||plan.state!=="in-production"}><legend>{t("p0876")}</legend><label>{t("p0877")}<select value={line} onChange={e=>setLine(e.target.value)}>{plan.lines.map(l=><option key={l.id} value={l.id}>{l.variant_id} {t("p0016")} {l.id}</option>)}</select></label><label>{t("p0878")}<input value={serial} onChange={e=>setSerial(e.target.value)} maxLength={128} required/></label><label>{t("p0879")}<input value={vin} onChange={e=>setVIN(e.target.value)} maxLength={128}/></label><label>{t("p0880")}<input value={battery} onChange={e=>setBattery(e.target.value)} maxLength={128}/></label><button className="button" disabled={!evidence}>{t("p0881")}</button></fieldset></form></>}
 <div className="grid">{plan.units.map(u=><article className="card" key={u.id} aria-label={`Serie ${u.serial_number}`}><h3>{u.serial_number}</h3><p>{t("p0882")} {labels[u.state]}{t("p0008")} {u.stock_state?`Stock: ${labels[u.stock_state]??u.stock_state}.`:""}</p>{u.receipt_review_state&&<p>{t("p0883")} {u.receipt_review_state==="pending"?t("p0797"):u.receipt_review_state==="approved"?t("p0798"):t("p0337")}{t("p0008")}</p>}{u.shipment_id&&<p>{t("p0884")} {u.shipment_id}</p>}
 {((surface==="factory"&&can("supply:factory")&&u.state==="released")||(surface==="franchise"&&can("supply:receive")&&u.state==="shipped"))&&<label><input type="checkbox" disabled={blocked} checked={selected.includes(u.id)} onChange={e=>setSelected(old=>e.target.checked?[...old,u.id]:old.filter(x=>x!==u.id))}/>{t("p0885")} {u.serial_number}</label>}
 {surface==="factory"&&can("supply:factory")&&<>{u.state==="planned"&&<button disabled={blocked||!evidence} onClick={()=>command("milestone",{unit_id:u.id,target_state:"assembly"})}>{t("p0886")}</button>}{u.state==="assembly"&&<button disabled={blocked||!evidence} onClick={()=>command("milestone",{unit_id:u.id,target_state:"quality"})}>{t("p0887")}</button>}{u.state==="quality"&&<><p>{t("p0888")}</p><button disabled={blocked||!evidence||!reason.trim()||u.factory_review_state!=="pending"||u.factory_requester===subject} onClick={()=>command("milestone",{unit_id:u.id,target_state:"released",reason})}>{t("p0889")}</button><button disabled={blocked||!evidence||!reason.trim()||u.factory_review_state!=="pending"||u.factory_requester===subject} onClick={()=>command("milestone",{unit_id:u.id,target_state:"rejected",reason})}>{t("p0890")}</button></>}</>}
 {surface==="franchise"&&u.stock_state==="quarantine"&&<><p>{t("p0891")}</p>{can("supply:release")&&<><button disabled={blocked||!evidence||!reason.trim()||u.receipt_review_state!=="pending"||u.receipt_requester===subject} onClick={()=>command("quality",{unit_id:u.id,reason})}>{t("p0892")}</button><button disabled={blocked||!evidence||!reason.trim()||u.receipt_review_state!=="pending"||u.receipt_requester===subject} onClick={()=>command("quality-reject",{unit_id:u.id,reason})}>{t("p0893")}</button></>}{can("supply:inspect")&&<button disabled={blocked||!evidence||!reason.trim()||u.receipt_review_state!=="rejected"} onClick={()=>command("reinspect",{unit_id:u.id,reason})}>{t("p0894")}</button>}</>}
 </article>)}</div>
 {plan.next_unit_id&&<button className="button" disabled={busy} onClick={()=>void consult(plan.next_unit_id)}>{t("p0895")}</button>}<button className="button" disabled={busy} onClick={()=>void consult()}>{t("p0896")}</button>
 {((surface==="factory"&&can("supply:factory"))||(surface==="franchise"&&can("supply:receive")))&&<form onSubmit={e=>{e.preventDefault();command(surface==="factory"?"ship":"receive",{shipment_id:shipment,units:selected})}}><fieldset disabled={blocked}><legend>{surface==="factory"?t("p0897"):t("p0898")}</legend><label>{t("p0899")}<input value={shipment} onChange={e=>setShipment(e.target.value)} maxLength={128} required/></label><p>{selected.length} {t("p0900")}</p><button className="button" disabled={!evidence||selected.length===0}>{surface==="factory"?t("p0901"):t("p0902")}</button></fieldset></form>}
 </>}
 </section>
 </>;
}
````

### FILE: `src/platform/supply/contract.test.ts`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "258d8a0355dcc48d10359d34b39acdb9a538f464df0528c80e357a60423a3808"
variables: []
secrets_allowed: false
```

````typescript
import {readFileSync} from "node:fs";
import {it,expect} from "vitest";
import {supplyCommand,supplyCanonical,supplyPayload,supplyDigest,supplyReference,supplyMatches} from "./contract";
const cases=JSON.parse(readFileSync("deploy/supply/role-goldens.json","utf8")) as {name:string;payload:Record<string,unknown>;canonical:string;sha256:string}[];
it.each(cases)("matches Go typed $name including exact int64 and Unicode",async c=>{const parsed=supplyCommand.parse({organization_id:"store",kind:c.name,...c.payload});expect(supplyCanonical(supplyPayload(parsed))).toBe(c.canonical);expect(await supplyDigest(supplyCanonical(supplyPayload(parsed)))).toBe(c.sha256)});
it("binds recovery to actor/order/version/payload without retaining evidence contents",async()=>{const c=supplyCommand.parse({organization_id:"store",kind:"create",...cases[0]!.payload}),p=await supplyReference(c),r={...p,actor:"buyer",recorded_at:"2026-09-12T00:00:00Z",payload_sha256:"b".repeat(64),payload:{},replay:false};expect(supplyMatches(r,p,"buyer")).toBe(true);for(const change of [{actor:"other"},{purchase_order_id:"other"},{version:"2"},{request_sha256:"c".repeat(64)}])expect(supplyMatches({...r,...change},p,"buyer")).toBe(false);expect(Object.keys(p).sort()).toEqual(["command_id","kind","purchase_order_id","request_sha256","version"])});
it("rejects duplicate lines, fields that do not belong to an action, and unsafe integers",()=>{const c={organization_id:"store",kind:"create",...cases[0]!.payload};expect(supplyCommand.safeParse({...c,lines:[{id:"same",variant_id:"a",quantity:1},{id:"same",variant_id:"b",quantity:1}]}).success).toBe(false);expect(supplyCommand.safeParse({organization_id:"store",...cases[1]!.payload,units:["extra"]}).success).toBe(false);expect(()=>supplyCanonical({quantity:1.5})).toThrow();expect(()=>supplyCanonical({quantity:1001})).toThrow()});
````

### FILE: `src/platform/supply/contract.ts`

```yaml
block_id: "TS-SERIAL-SUPPLY-PORTAL:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7454b7aec65df1ca66034d90ff6106040836d6d3097bd46eef75e54b0dfe8df0"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED transport/form binding to the existing Go serial supply owner.
import {z} from "zod";
export const supplyID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/u);
export const supplySHA=z.string().regex(/^[0-9a-f]{64}$/u);
const unsigned=z.string().regex(/^(0|[1-9][0-9]{0,18})$/u).refine(v=>BigInt(v)<=9223372036854775807n);
const positive=unsigned.refine(v=>v!=="0");
const text=(max:number)=>z.string().min(1).refine(v=>new TextEncoder().encode(v).length<=max&&v.trim()===v&&!/[\p{Cc}]/u.test(v)&&new TextDecoder().decode(new TextEncoder().encode(v))===v);
export const supplyPolicy="strict-serial-reference/v1";
export const supplyKinds=["submit","confirm","start","cancel","register","milestone","ship","receive","quality","quality-reject","reinspect"] as const;
const line=z.object({id:supplyID,variant_id:supplyID,quantity:z.number().int().min(1).max(1000)}).strict();
const planFields={factory_organization_id:supplyID,demand_reference:text(1024),policy_code:z.literal(supplyPolicy),lines:z.array(line).min(1).max(32).refine(v=>new Set(v.map(l=>l.id)).size===v.length&&new Set(v.map(l=>l.variant_id)).size===v.length&&v.reduce((n,l)=>n+l.quantity,0)<=1000)};
const common={organization_id:supplyID,purchase_order_id:supplyID,command_id:supplyID,evidence_sha256:supplySHA};
const base={...common,expected_version:positive};
const simple=(kind:"submit"|"confirm"|"start"|"cancel")=>z.object({kind:z.literal(kind),...base}).strict();
export const supplyCommand=z.discriminatedUnion("kind",[
 z.object({kind:z.literal("create"),...common,...planFields,destination_organization_id:supplyID,supplier_id:supplyID,currency:z.string().regex(/^[A-Z]{3}$/u),total_minor_units:unsigned}).strict().refine(v=>v.organization_id===v.destination_organization_id),
 z.object({kind:z.literal("plan"),...common,...planFields}).strict(),
 simple("submit"),simple("confirm"),simple("start"),simple("cancel"),
 z.object({kind:z.literal("register"),...base,line_id:supplyID,serial_number:text(128),vin:text(128).optional(),battery_serial_number:text(128).optional()}).strict(),
 z.object({kind:z.literal("milestone"),...base,unit_id:supplyID,target_state:z.enum(["assembly","quality","released","rejected"]),reason:text(2000).optional()}).strict().refine(v=>["released","rejected"].includes(v.target_state)?v.reason!==undefined:v.reason===undefined),
 ...(["ship","receive"] as const).map(kind=>z.object({kind:z.literal(kind),...base,shipment_id:supplyID,units:z.array(supplyID).min(1).max(100).refine(v=>new Set(v).size===v.length)}).strict()),
 ...(["quality","quality-reject","reinspect"] as const).map(kind=>z.object({kind:z.literal(kind),...base,unit_id:supplyID,reason:text(2000)}).strict())
]);
export type SupplyCommand=z.infer<typeof supplyCommand>;
export const supplyReceipt=z.object({purchase_order_id:supplyID,command_id:supplyID,version:positive,kind:z.enum(["planned",...supplyKinds]),actor:z.string().min(1).max(256),request_sha256:supplySHA,payload_sha256:supplySHA,payload:z.unknown(),recorded_at:z.iso.datetime({offset:true}),replay:z.boolean()}).strict();
export type SupplyReceipt=z.infer<typeof supplyReceipt>;
const unit=z.object({factory_review_state:z.enum(["pending","approved","rejected"]).optional(),factory_requester:z.string().min(1).max(256).optional(),receipt_review_state:z.enum(["pending","approved","rejected"]).optional(),receipt_requester:z.string().min(1).max(256).optional(),id:supplyID,line_id:supplyID,state:z.enum(["planned","assembly","quality","released","rejected","shipped","received"]),serial_number:text(128),stock_unit_id:supplyID.optional(),stock_state:z.enum(["in-transit","quarantine","available","reserved","sold","service","retired"]).optional(),shipment_id:supplyID.optional()}).strict();
export const supplyPlan=z.object({purchase_order_id:supplyID,destination_organization_id:supplyID,factory_organization_id:supplyID,supplier_id:supplyID,state:z.enum(["draft","submitted","accepted","in-production","shipped","received","cancelled"]),purchase_version:positive,version:positive,currency:z.string().regex(/^[A-Z]{3}$/u),total_minor_units:unsigned,policy_code:z.literal(supplyPolicy),demand_reference:text(1024),lines:z.array(line).min(1).max(32),units:z.array(unit).max(100),latest:supplyReceipt.optional(),next_unit_id:supplyID.optional()}).strict();
export type SupplyPlan=z.infer<typeof supplyPlan>;
export const pendingSupply=z.object({purchase_order_id:supplyID,command_id:supplyID,kind:supplyReceipt.shape.kind,version:positive,request_sha256:supplySHA}).strict();
export type PendingSupply=z.infer<typeof pendingSupply>;
export function supplyPermission(kind:SupplyCommand["kind"]){return ["confirm","start","register","milestone","ship"].includes(kind)?"supply:factory":kind==="receive"?"supply:receive":["quality","quality-reject"].includes(kind)?"supply:release":kind==="reinspect"?"supply:inspect":"supply:plan"}
export function supplySurface(kind:SupplyCommand["kind"]){return supplyPermission(kind)==="supply:factory"?"factory":"franchise"}
export function supplyPayload(c:SupplyCommand){const{organization_id:_,...body}=c;if(body.kind==="create"||body.kind==="plan"){const{kind:__,...plan}=body;return plan}return body}
// Narrow Go encoding/json correspondence: ASCII keys, validated strings and
// whole quantities1..1000 only. Not a general-purpose JSON canonicalizer.
export function supplyCanonical(value:unknown):string{
 const literal=(v:string)=>JSON.stringify(v).replace(/[<>&\u2028\u2029]/gu,x=>({"<":"\\u003c",">":"\\u003e","&":"\\u0026","\u2028":"\\u2028","\u2029":"\\u2029"}[x]!));
 if(value===null)return"null";if(typeof value==="string")return literal(value);
 if(typeof value==="number"&&Number.isInteger(value)&&value>=1&&value<=1000)return String(value);
 if(Array.isArray(value))return"["+value.map(supplyCanonical).join(",")+"]";
 if(typeof value==="object"){const o=value as Record<string,unknown>,keys=Object.keys(o).sort();if(keys.some(k=>!/^[a-z_][a-z0-9_]*$/u.test(k)))throw new Error("unsupported key");return"{"+keys.map(k=>literal(k)+":"+supplyCanonical(o[k])).join(",")+"}"}
 throw new Error("unsupported supply input");
}
export async function supplyDigest(value:string|Uint8Array){const bytes=typeof value==="string"?new TextEncoder().encode(value):value;return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",new Uint8Array(bytes).buffer)),v=>v.toString(16).padStart(2,"0")).join("")}
export async function supplyReference(c:SupplyCommand):Promise<PendingSupply>{return pendingSupply.parse({purchase_order_id:c.purchase_order_id,command_id:c.command_id,kind:c.kind==="create"||c.kind==="plan"?"planned":c.kind,version:"expected_version"in c?(BigInt(c.expected_version)+1n).toString():"1",request_sha256:await supplyDigest(supplyCanonical(supplyPayload(c)))})}
export function supplyMatches(v:SupplyReceipt,p:PendingSupply,subject:string){return v.purchase_order_id===p.purchase_order_id&&v.command_id===p.command_id&&v.actor===subject&&v.kind===p.kind&&v.version===p.version&&v.request_sha256===p.request_sha256}
````

## 6. Configuration surface

docs/SUPPLY_ROLE_REFERENCE.md. features.supply_portal opt-in, supply:read/factory-read plus action permission. Source order IDs are preallocated recovery references, event IDs remain random. Explicit currency/total,32lines/1000units,100unit pages. No new migration.

## 7. Dependency bill

16new and7changed AUTHORED unavoidable composition/transport/UI/config/test glue. No new dependency/upstream/corporate attribution; existing pins and adapted-source licenses retained.

## 8. Apply order

Select TS-GO-API-WEB-BRIDGE, TS-OIDC-PORTAL-ADAPTER and TS-FRANCHISE-JOURNEY-PORTALS. Full Go serial supply/creation owners required when enabled; narrow web proves imports/types only. Shared command goldens owned here.

## 9. Verification

Actual Next/BFF/Go/PG/Chromium JWE/RS256/JWKS,23writes:1received order,3factory serials with1rejected/replaced,2available stock,6quality requests/decisions. Create/receive responses lost then recovered byGET after reload without extraPOST. Creation1new3replay,5table rollback including late outbox, actor/hash/org/int64;4Go/TS goldens;4BFF/6contract tests. Full webpack/types. Fuzz2s3seeds99596executions. Visual desktop/390px; nav key correction checked in exact web type closure.

## 10. Reconstruction evidence

SUPPLY_ROLE_RELEASE_V402.md/json: exact source/delta, RED/PASS receipts and profile rebuilds. Does not close allT2804; warranty/J5/CMS/KPIs/private locale remain. No production certification.


V402 composed delta: T2804 help CMS and21same-release guides;15oldguides unchanged,5bounded training curricula revision2, no automatic grants. New optional host, shared inline text, actual browser/PG proofs. HELP_CMS_RELEASE_V402.md.

V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.