# Catalog authoring role portal

## 1. Metadata

```yaml
pack_id: "TS-CATALOG-AUTHORING-PORTAL"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Role catalog authoring from initial model/variants/prices and PNG through existing approval/publication/storefront; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804/J3. Existing catalog release, model/price services, shared approvals, OIDC/BFF and Next are required. Source form creates one model with1–16variants; frozen drafts are immutable, publication supports rollback to a reviewed version.

## 3. Architecture contract

Original model SQL/event extracted into its owner transaction method; original model and Commerce services are bound to one transaction through narrow repositories. Existing variant table, catalog command ledger and outbox; source remains draft/private/unknown homologation. Existing publication and review rules unchanged.

## 4. Exact file manifest

```text
CREATE config/catalog.authoring.fixture.json
CREATE deploy/catalog/authoring-goldens.json
CREATE docs/CATALOG_ROLE_AUTHORING_REFERENCE.md
CREATE microsoft_playwright_browser_gate/tests/catalog-authoring-connected.spec.mjs
CREATE src/app/admin/catalog/page.tsx
CREATE src/app/api/enterprise/catalog/route.test.ts
CREATE src/app/api/enterprise/catalog/route.ts
CREATE src/components/catalog-authoring.tsx
CREATE src/platform/catalog/authoring.test.ts
CREATE src/platform/catalog/authoring.ts
```

## 5. Materialization blocks

### FILE: `config/catalog.authoring.fixture.json`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "920fd6d05b650b293a9c2275dd3cdb9110ce32b9d390961ad229207c4c42411a"
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
    "catalog_editor": true,
    "role_workspace": true
  }
}
````

### FILE: `deploy/catalog/authoring-goldens.json`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "17f42b25225652def7e0b2fa728e71130ce369ecf6661383f2e810d491c4c9b4"
variables: []
secrets_allowed: false
```

````json
[
  {
    "name": "source",
    "payload": {
      "command_id": "golden-source",
      "model": {
        "id": "",
        "code": "bicycle-one",
        "displayName": "Modelo á <&>",
        "vehicleClass": "bicycle",
        "specification": {
          "description": "Guía   fija"
        }
      },
      "variants": [
        {
          "code": "variant-one",
          "display_name": "Variante",
          "battery_specification": {
            "description": "Referencia   fija"
          },
          "amount_minor_units": "9007199254740993",
          "tax_mode": "not-applicable"
        }
      ],
      "valid_from": "2026-09-12T00:00:00Z",
      "valid_until": "2026-10-12T00:00:00.12Z"
    },
    "canonical": "{\"command_id\":\"golden-source\",\"model\":{\"code\":\"bicycle-one\",\"displayName\":\"Modelo á \\u003c\\u0026\\u003e\",\"id\":\"\",\"specification\":{\"description\":\"Guía \\u2028 fija\"},\"vehicleClass\":\"bicycle\"},\"valid_from\":\"2026-09-12T00:00:00Z\",\"valid_until\":\"2026-10-12T00:00:00.12Z\",\"variants\":[{\"amount_minor_units\":\"9007199254740993\",\"battery_specification\":{\"description\":\"Referencia \\u2029 fija\"},\"code\":\"variant-one\",\"display_name\":\"Variante\",\"tax_mode\":\"not-applicable\"}]}",
    "sha256": "7072b34182c9fca16217271eb597ec0df98cad5730ac5e3bccc05c0ba2d5cb88"
  },
  {
    "name": "review",
    "payload": {
      "command_id": "golden-review",
      "draft_id": "draft-one",
      "stage": "legal",
      "snapshot_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "approved": true,
      "reason": "Revisión <documento> & alcance",
      "evidence_sha256": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
    },
    "canonical": "{\"approved\":true,\"command_id\":\"golden-review\",\"draft_id\":\"draft-one\",\"evidence_sha256\":\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\",\"reason\":\"Revisión \\u003cdocumento\\u003e \\u0026 alcance\",\"snapshot_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"stage\":\"legal\"}",
    "sha256": "79a960fc55b25b6d2a01f68f7521736f18ff50d0f2a4615fbc722a8304211326"
  },
  {
    "name": "publish",
    "payload": {
      "command_id": "golden-publish",
      "draft_id": "draft-one",
      "snapshot_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "expected_generation": "9007199254740993",
      "reason": "Publicación aprobada"
    },
    "canonical": "{\"command_id\":\"golden-publish\",\"draft_id\":\"draft-one\",\"expected_generation\":\"9007199254740993\",\"reason\":\"Publicación aprobada\",\"snapshot_sha256\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\"}",
    "sha256": "c3aaafed180ffd0c5d032542c403873f22ab6b5b9cffa36306df15766348e67b"
  },
  {
    "name": "media",
    "payload": {
      "command_id": "golden-media",
      "original_sha256": "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
    },
    "canonical": "{\"command_id\":\"golden-media\",\"original_sha256\":\"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\"}",
    "sha256": "cf668e67145a569401a2246a52d01d1d42a28381b63a53a376b675265d31fe56"
  }
]
````

### FILE: `docs/CATALOG_ROLE_AUTHORING_REFERENCE.md`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "755d68908cd6ba335608fe6277bfbdd7c77d2695d195997ca63137fe1238fb6a"
variables: []
secrets_allowed: false
```

````markdown
# Catálogo por rol — referencia de infraestructura

La página /admin/catalog permite crear un modelo con1–16variantes, precios
explícitos y vigencia UTC, cargar PNG, congelar una versión, revisar sus cuatro
etapas y publicarla. El formulario de versión compone un modelo con su libro;
el contrato backend admite hasta32modelos y128variantes en un snapshot. No se
afirma un editor masivo multi-modelo ni edición de borradores ya congelados.
Para cambiar una versión publicada se crea y revisa otro snapshot; restaurar
una versión aprobada crea una publicación nueva con un libro efectivo nuevo.

Seleccionar los packs de autoría y catálogo conectado junto con la composición
Go/PG/Next/BFF/identidad existente. La referencia completa incorpora todo.
features.catalog_editor=true muestra el portal; la sesión debe tener catalog:read.
catalog:draft crea fuentes/PNG/snapshot. catalog:review:legal, :technical, :media y
:publication permiten las revisiones correspondientes. catalog:publish publica;
el creador del snapshot no puede revisarlo ni publicarlo. Cada permiso sigue
limitado por tenant/organización real en Go; el nombre de rol no autoriza.

Aplicar migration0076 después de0075. El host del catálogo exige ese contrato
además de sus guards anteriores y perfil SHA fijado. La migración down rechaza
historial source; en base vacía el down/up funciona. Los modelos/variantes y
libros iniciales permanecen draft, modelos no públicos y homologación unknown.
La revisión humana no inventa acreditación fiscal ni homologación oficial.

El escritor de modelo conserva SQL/evento originales de Electromobility;
sus services y los de Commerce se componen en la misma transacción mediante
repositories acotados. La variante usa la tabla existente; release_command
conserva actor/identidad/hash/referencias y outbox. No segundo ledger.
La transacción completa se revierte si falla el outbox; cuatro solicitudes
simultáneas de la misma identidad producen un alta y tres replays.

El navegador conserva sólo referencias/hashes de recuperación y comprobantes
de fuente/imagen, ligados a tenant/org/subject. No guarda razones de revisión,
archivos de evidencia ni PNG en sessionStorage. Se verifica actor/command/hash/
snapshot antes de resolver una respuesta perdida; la recuperación sólo consulta.
Las huellas de source/review/publish/media se comparan entre Go y TypeScript,
incluidos caracteres no ASCII, escapes HTML y valores int64 como texto.
Este contrato estrecho no es un serializador JSON/JCS genérico.

Los JSON están limitados a32KiB y el PNG a1MiB, con plazo de lectura y cancelación
del reader. El helper textual conserva sus presupuestos y comportamiento UTF-8.
El BFF verifica sesión/permiso antes de consumir el cuerpo; una publicación
vigente se consulta con autorización privada, no desde otro tenant público.

Ejecutar TestCatalogSourceAuthoringAtomic con PAYMENT_CONNECTED_DB_URL aislada,
TestCatalogRoleCommandGolden y la prueba de navegador
TestCatalogAuthoringRoleBrowser con ELITE_CATALOG_ROLE_BROWSER=1,
ELITE_WEB_ROOT/ELITE_NODE_BIN absolutos y Next compilado --webpack.
La fixture sólo crea tenant/org; las16escrituras salen del navegador:
2modelos/2variantes/2snapshots/8reviews/3publicaciones, con rollback a la primera.
Los dos replays por pérdida de respuesta son GET y no agregan un POST.
Las5price_books finales son2de fuente y3efectivas, sin reactivar una antigua.
config/catalog.authoring.fixture.json y el goldens JSON son públicos sintéticos.

La revisión visual verificó escritorio y390px sin desborde. Los códigos de
estado/tributación/homologación reflejan sus owners; la localización privada
completa continúa bajo T2804. Tampoco se cierra supply/warranty/J5/CMS/KPIs ni
el mapping de marketplace T2805 por este recorrido.
Procedencia: AUTHORED glue; ninguna dependencia o atribución corporativa nueva.
````

### FILE: `microsoft_playwright_browser_gate/tests/catalog-authoring-connected.spec.mjs`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "02fa76757f8ec1b1577729e882d8c9e2cfc566d6368a7608b7794873551c29d1"
variables: []
secrets_allowed: false
```

````text
import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
test('catalog role source to publication and rollback preserves command identities after response loss',async({page,context},info)=>{
 if(process.env.ELITE_CATALOG_ROLE_BROWSER!=='1')throw new Error('explicit local fixture required');
 page.setDefaultTimeout(12000);
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 const identities=JSON.parse(process.env.ELITE_CATALOG_IDENTITIES);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[];page.on('pageerror',e=>errors.push(e.message));const posts=[];page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/catalog'))posts.push(r.url())});
 async function editor(draft=''){await page.goto('/admin/catalog'+(draft?'?draft='+encodeURIComponent(draft):''));await expect(page.getByRole('heading',{name:'Edición y publicación del catálogo'})).toBeVisible()}
 async function dropOnce(action){
  let dropped=false;
  await page.route('**/api/enterprise/catalog**',async route=>{
   const r=route.request();if(dropped||r.method()!=='POST'||r.headers()['content-type']!=='application/json'||r.postDataJSON().action!==action){await route.continue();return}
   dropped=true;const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed');
  });
 }
 async function recover(){await expect(page.getByRole('status')).toContainText('Resultado sin confirmar');const before=posts.length;await page.unrouteAll();await page.reload();await page.getByRole('button',{name:'Consultar resultado pendiente',exact:true}).click();await expect(page.getByRole('status')).toContainText('Resultado recuperado');expect(posts.length).toBe(before)}
 await identity('unprivileged');await page.goto('/admin/catalog');await expect(page.getByText('No tenés permiso para consultar este espacio.')).toBeVisible();
 await identity('foreign-org');await editor();await page.getByRole('button',{name:'Consultar publicación vigente',exact:true}).click();await expect(page.getByRole('status')).toContainText('No pudimos confirmar');await expect(page.getByText('Todavía no hay una publicación.',{exact:true})).toHaveCount(0);
 await identity('reader');expect((await context.request.post(base+'/api/enterprise/catalog',{headers:{origin:base,'content-type':'application/json'},data:{}})).status()).toBe(403);
 async function createVersion(suffix,name,amount,loseSource,upload){
  await identity('maker');await editor();
  await page.getByLabel('Código del modelo',{exact:true}).fill('role-'+suffix);
  await page.getByLabel('Nombre del modelo',{exact:true}).fill(name);
  await page.getByRole('combobox',{name:'Clase de vehículo',exact:true}).selectOption('bicycle');
  await page.getByLabel('Descripción técnica',{exact:true}).fill('Ficha técnica de referencia <verificada> & sintética.');
  await page.getByLabel('Vigencia desde (UTC)',{exact:true}).fill(process.env.ELITE_CATALOG_VALID_FROM);
  await page.getByLabel('Código de variante',{exact:true}).fill('variant-'+suffix);
  await page.getByLabel('Nombre de variante',{exact:true}).fill('Variante '+suffix);
  await page.getByLabel('Especificación de batería',{exact:true}).fill('Batería de referencia, sin datos privados.');
  await page.getByLabel('Importe en unidades menores',{exact:true}).fill(amount);
  await page.getByRole('combobox',{name:'Tratamiento del precio',exact:true}).selectOption('not-applicable');
  if(loseSource)await dropOnce('source');
  await page.getByRole('button',{name:'Guardar modelo y precios',exact:true}).click();
  if(loseSource)await recover();else await expect(page.getByRole('status')).toContainText('Modelo, variantes y precios guardados');
  expect(await page.getByLabel('Referencia del modelo',{exact:true}).inputValue()).not.toBe('');
  if(upload){
   await page.getByLabel('Imagen PNG',{exact:true}).setInputFiles({name:'reference.png',mimeType:'image/png',buffer:Buffer.from(process.env.ELITE_CATALOG_PNG,'base64')});
   await page.getByRole('button',{name:'Guardar imagen',exact:true}).click();await expect(page.getByRole('status')).toContainText('Imagen recibida');
  }
  expect(await page.getByLabel('Referencia de la imagen',{exact:true}).inputValue()).not.toBe('');
  await page.getByRole('button',{name:'Guardar versión para revisión',exact:true}).click();await expect(page.getByRole('status')).toContainText('Versión guardada');
  const id=await page.getByLabel('Referencia del borrador',{exact:true}).inputValue();expect(id).not.toBe('');
  await page.getByRole('button',{name:'Consultar borrador',exact:true}).click();await expect(page.getByRole('heading',{name,exact:true})).toBeVisible();
  await expect(page.getByText('Otra persona autorizada debe revisar y publicar esta versión.',{exact:true})).toBeVisible();await expect(page.getByRole('button',{name:'Aprobar etapa revisada',exact:true})).toHaveCount(0);
  return id;
 }
 async function review(id,expectedName){
  await identity('reviewer');await editor(id);
  await page.getByRole('button',{name:'Consultar borrador',exact:true}).click();await expect(page.getByRole('heading',{name:expectedName,exact:true})).toBeVisible();
  await page.getByLabel('Evidencia de la revisión',{exact:true}).setInputFiles({name:'review.txt',mimeType:'text/plain',buffer:Buffer.from('Synthetic evidence of explicit reference review.')});
  for(const stage of ['legal','technical','media','publication']){
   await page.getByRole('combobox',{name:'Etapa de revisión',exact:true}).selectOption(stage);
   await page.getByRole('textbox',{name:'Motivo de la decisión',exact:true}).fill('Revisé el contenido sintético para la etapa '+stage+'.');
   await page.getByRole('button',{name:'Aprobar etapa revisada',exact:true}).click();
   const label={legal:'Legal',technical:'Técnica',media:'Imagen',publication:'Publicación'}[stage];
   await expect(page.getByLabel('Estado de revisiones').getByText(label+': approved',{exact:true})).toBeVisible();
  }
 }
 async function publish(name,generation,lose=false){
  await page.getByRole('button',{name:'Consultar publicación vigente',exact:true}).click();await expect(page.getByRole('status')).toContainText('Publicación vigente consultada');
  await page.getByRole('textbox',{name:'Motivo de publicación',exact:true}).fill('Publicar la versión revisada de referencia.');
  if(lose)await dropOnce('publish');
  await page.getByRole('button',{name:'Publicar versión revisada',exact:true}).click();
  if(lose)await recover();else await expect(page.getByRole('status')).toContainText('Publicación registrada');
  await page.getByRole('button',{name:'Consultar publicación vigente',exact:true}).click();await expect(page.getByText('Versión publicada '+generation+'.',{exact:true})).toBeVisible();
  await page.goto('/models');await expect(page.getByRole('heading',{name,exact:true})).toBeVisible();
 }
 const first=await createVersion('one','Modelo Rol Uno','123456',true,true);
 await review(first,'Modelo Rol Uno');await publish('Modelo Rol Uno','1',true);
 const second=await createVersion('two','Modelo Rol Dos','234567',false,false);
 await review(second,'Modelo Rol Dos');await publish('Modelo Rol Dos','2');
 await editor(first);await page.getByRole('button',{name:'Consultar borrador',exact:true}).click();await expect(page.getByRole('heading',{name:'Modelo Rol Uno',exact:true})).toBeVisible();
 await publish('Modelo Rol Uno','3');await expect(page.getByRole('heading',{name:'Modelo Rol Dos',exact:true})).toHaveCount(0);
 await page.screenshot({path:info.outputPath('catalog-role-published-desktop.png'),fullPage:true});
 await editor(first);await page.getByRole('button',{name:'Consultar borrador',exact:true}).click();await page.setViewportSize({width:390,height:844});
 expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
 await page.screenshot({path:info.outputPath('catalog-role-review-mobile.png'),fullPage:true});expect(errors).toEqual([]);
 // 16 successful writes; the deliberate reader POST was rejected by the BFF.
 expect(posts).toHaveLength(16);
});
````

### FILE: `src/app/admin/catalog/page.tsx`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "071a161f121f68617606d0070ada376cdcb43409e478f78b923c636d49af27c2"
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
import {catalogID} from "@/platform/catalog/authoring";
import {CatalogAuthoring} from "@/components/catalog-authoring";
export default async function CatalogAuthoringPage({searchParams}:{searchParams:Promise<Record<string,string|string[]|undefined>>}){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.catalog_editor!==true)notFound();
 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/admin/catalog" as Route);
 if(!allowed(session,"catalog:read"))return <><h1>{t("p0001")}</h1><p>{t("p0002")}</p></>;
 const q=await searchParams,requested=q.organization;
 const organization=typeof requested==="string"&&session.organizations.includes(requested)?requested:session.organizations[0];
 if(!organization)return <><h1>{t("p0001")}</h1><p>{t("p0003")}</p></>;
 const initialDraft=typeof q.draft==="string"&&catalogID.safeParse(q.draft).success?q.draft:"";
 const scope=createHash("sha256").update(JSON.stringify([session.tenantId,organization,session.subject])).digest("hex");
 return <><h1 className="pageTitle">{t("p0004")}</h1><p>{t("p0005")} {organization}</p><CatalogAuthoring key={scope} scope={scope} organization={organization} tenant={session.tenantId} subject={session.subject} permissions={session.permissions} initialDraft={initialDraft}/></>;
}
````

### FILE: `src/app/api/enterprise/catalog/route.test.ts`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "20e2b50709b8653057c0170963bd59d15acf2919daa3292126c2498176de41fa"
variables: []
secrets_allowed: false
```

````typescript
import {NextRequest} from "next/server";
import {expect,it,vi,afterEach} from "vitest";
const m=vi.hoisted(()=>({get:vi.fn(),post:vi.fn(),png:vi.fn(),session:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post,protectedPostPNG:m.png}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:async()=>({features:{catalog_editor:true}})}));
import {POST,GET} from "./route";
afterEach(()=>vi.clearAllMocks());
const req=(body:BodyInit,tail="",type="application/json")=>{
 const options={method:"POST",headers:{origin:"https://portal.example.test","content-type":type},body,duplex:"half" as const};
 return new NextRequest("https://portal.example.test/api/enterprise/catalog"+tail,options);
};
it("rejects unauthenticated media before body consumption",async()=>{
 m.session.mockResolvedValue(null);const r=req(new Uint8Array([1,2,3]),"?organization_id=org&command_id=media","image/png");
 expect((await POST(r)).status).toBe(401);expect(r.bodyUsed).toBe(false);expect(m.png).not.toHaveBeenCalled();
});
it("bounds chunked binary media and cancels before downstream",async()=>{
 m.session.mockResolvedValue({subject:"maker",organizations:["org"],permissions:["catalog:draft","catalog:read"]});
 let cancelled=false,pulls=0;const stream=new ReadableStream<Uint8Array>({pull(c){pulls++;c.enqueue(new Uint8Array(262145))},cancel(){cancelled=true}});
 const r=req(stream,"?organization_id=org&command_id=media","image/png");
 expect((await POST(r)).status).toBe(413);expect(cancelled).toBe(true);expect(pulls).toBeLessThanOrEqual(5);expect(m.png).not.toHaveBeenCalled();
});
it("rejects duplicate and foreign scope before upload",async()=>{
 m.session.mockResolvedValue({subject:"maker",organizations:["org"],permissions:["catalog:draft","catalog:read"]});
 for(const tail of ["?organization_id=other&command_id=media","?organization_id=org&organization_id=org&command_id=media"]){const r=req(new Uint8Array([1]),tail,"image/png");expect((await POST(r)).status).toBe(403);expect(r.bodyUsed).toBe(false)}
 expect(m.png).not.toHaveBeenCalled();
});
it("uses private authorized current publication and rejects foreign draft projection",async()=>{
 m.session.mockResolvedValue({subject:"reader",tenantId:"50f38793-8a22-4f6b-983f-81dd0fca8208",organizations:["org"],permissions:["catalog:read"]});
 m.get.mockRejectedValue(new Error("unavailable"));
 expect((await GET(new NextRequest("https://portal.example.test/api/enterprise/catalog?kind=current&organization_id=org"))).status).toBe(409);
 expect(m.get).toHaveBeenCalledWith(expect.anything(),"/v1/admin/catalog/current",{organization_id:"org"});
 expect(m.post).not.toHaveBeenCalled();
});
````

### FILE: `src/app/api/enterprise/catalog/route.ts`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "56221e799afcd402170692b14772c2cb427daa94c343a70623b67f5b09db5d2a"
variables: []
secrets_allowed: false
```

````typescript
import {NextResponse,type NextRequest} from "next/server";
import {readSession,allowed} from "@/platform/auth/session";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {protectedGet,protectedPost,protectedPostPNG} from "@/platform/backend/protected-client";
import {BackendProblem} from "@/platform/backend/public-client";
import {boundedCommandBody,boundedRequestBytes} from "@/platform/backend/bounded-command";
import {loadBusinessConfig} from "@/platform/config/load";
import {catalogCommand,catalogDraft,catalogReceipt,catalogID,commandPayload} from "@/platform/catalog/authoring";
import {publishedCatalogSchema} from "@/platform/catalog/schema";
const response=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
const result=(v:unknown)=>NextResponse.json(v,{headers:{"cache-control":"no-store"}});
const permits=["catalog:draft","catalog:publish","catalog:review:legal","catalog:review:technical","catalog:review:media","catalog:review:publication"];
async function enabled(){return(await loadBusinessConfig()).features.catalog_editor===true}
export async function GET(r:NextRequest){
 if(!await enabled())return response("NOT_FOUND",404);
 const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 if(!allowed(session,"catalog:read"))return response("FORBIDDEN",403);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),org=q.get("organization_id"),id=q.get("id");
 const keys=kind==="current"?["kind","organization_id"]:["kind","organization_id","id"];
 if([...q.keys()].length!==keys.length||keys.some(k=>q.getAll(k).length!==1)||![...q.keys()].every(k=>keys.includes(k))||!catalogID.safeParse(org).success||!session.organizations.includes(org!))return response("INVALID_SCOPE",403);
 if(!["current","draft","command"].includes(kind??"")||kind!=="current"&&!catalogID.safeParse(id).success)return response("INVALID_REFERENCE",400);
 try{
  if(kind==="current"){
   try{return result({publication:publishedCatalogSchema.parse(await protectedGet(session,"/v1/admin/catalog/current",{organization_id:org!}))})}
   catch(e){if(e instanceof BackendProblem&&e.status===404)return result({publication:null});throw e}
  }
  const raw=await protectedGet(session,`/v1/admin/catalog/${kind==="draft"?"drafts":"commands"}/${encodeURIComponent(id!)}`,{organization_id:org!});
  if(kind==="draft"){const value=catalogDraft.parse(raw);if(value.id!==id||value.snapshot.profile.tenant_id!==session.tenantId||value.snapshot.profile.organization_id!==org)return response("SCOPE_MISMATCH",409);return result(value)}
  const value=catalogReceipt.parse(raw);if(value.command_id!==id||value.actor!==session.subject)return response("SCOPE_MISMATCH",409);return result(value);
 }catch{return response("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(r:NextRequest){
 if(!await enabled())return response("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return response("CROSS_ORIGIN_REJECTED",403)}catch{return response("UNAVAILABLE",503)}
 const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 if(!allowed(session,"catalog:read")||!permits.some(p=>allowed(session,p)))return response("FORBIDDEN",403);
 try{
  const type=r.headers.get("content-type");let value:unknown,commandID:string;
  if(type==="image/png"){
   if(!allowed(session,"catalog:draft"))return response("FORBIDDEN",403);
   const q=r.nextUrl.searchParams,org=q.get("organization_id"),id=q.get("command_id");
   if([...q.keys()].length!==2||q.getAll("organization_id").length!==1||q.getAll("command_id").length!==1||!catalogID.safeParse(org).success||!catalogID.safeParse(id).success||!session.organizations.includes(org!))return response("INVALID_SCOPE",403);
   commandID=id!;const bytes=await boundedRequestBytes(r,1048576);
   value=await protectedPostPNG(session,`/v1/admin/catalog/media/${encodeURIComponent(id!)}?${new URLSearchParams({organization_id:org!})}`,bytes);
  }else{
   if(type!=="application/json")return response("UNSUPPORTED_MEDIA_TYPE",415);
   if(r.nextUrl.searchParams.size!==0)return response("INVALID_COMMAND",400);
   let body:unknown;try{body=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return response("INVALID_COMMAND",400)}
   const parsed=catalogCommand.safeParse(body);if(!parsed.success)return response("INVALID_COMMAND",400);
   const c=parsed.data;commandID=c.command_id;
   if(!session.organizations.includes(c.organization_id))return response("INVALID_SCOPE",403);
   const permission=c.action==="review"?`catalog:review:${c.stage}`:c.action==="publish"?"catalog:publish":"catalog:draft";
   if(!allowed(session,permission))return response("FORBIDDEN",403);
   const tail=c.action==="source"?"sources":c.action==="draft"?"drafts":`drafts/${encodeURIComponent(c.draft_id)}/${c.action==="review"?"review/"+c.stage:"publish"}`;
   value=await protectedPost(session,`/v1/admin/catalog/${tail}?${new URLSearchParams({organization_id:c.organization_id})}`,commandPayload(c));
  }
  const v=catalogReceipt.parse(value);
  if(v.command_id!==commandID||v.actor!==session.subject)return response("UNCONFIRMED",409);
  return result(v);
 }catch(e){return response(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
````

### FILE: `src/components/catalog-authoring.tsx`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bef83d3f597c19d9a17ac6e8ca2ac980d1c5ff3ac5d095cf1dd318f12cfbd022"
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
import {catalogCommand,catalogDraft,catalogReceipt,pendingCatalog,commandReference,catalogCanonical,catalogSHA,receiptMatches,type CatalogCommand,type CatalogDraft,type CatalogReceipt,type PendingCatalog} from "@/platform/catalog/authoring";
import {publishedCatalogSchema,type PublishedCatalog} from "@/platform/catalog/schema";
const endpoint="/api/enterprise/catalog";
const capsule=z.object({pending:pendingCatalog.nullable(),source:catalogReceipt.nullable(),media:catalogReceipt.nullable()}).strict();
type Capsule=z.infer<typeof capsule>;
type VariantInput={code:string;name:string;battery:string;amount:string;tax:"inclusive"|"exclusive"|"not-applicable"};
const emptyVariant=():VariantInput=>({code:"",name:"",battery:"",amount:"",tax:"not-applicable"});
const stages=["legal","technical","media","publication"] as const;

export function CatalogAuthoring({organization,tenant,subject,permissions,scope,initialDraft}:{organization:string;tenant:string;subject:string;permissions:string[];scope:string;initialDraft:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const stageLabel={legal:t("p0144"),technical:t("p0145"),media:t("p0146"),publication:t("p0147")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p);
 const held=useRef(false),saved=useRef<Capsule>({pending:null,source:null,media:null});
 const [storageReady,setStorageReady]=useState(false),[busy,setBusy]=useState(false),[pending,setPending]=useState<PendingCatalog|null>(null),[notice,setNotice]=useState("");
 const [modelCode,setModelCode]=useState(""),[modelName,setModelName]=useState(""),[vehicleClass,setVehicleClass]=useState("bicycle"),[description,setDescription]=useState("");
 const [variants,setVariants]=useState<VariantInput[]>([emptyVariant()]),[from,setFrom]=useState(""),[until,setUntil]=useState("");
 const [modelID,setModelID]=useState(""),[mediaID,setMediaID]=useState(""),[bookID,setBookID]=useState(""),[file,setFile]=useState<File|null>(null);
 const [draftID,setDraftID]=useState(initialDraft),[draft,setDraft]=useState<CatalogDraft|null>(null),[current,setCurrent]=useState<PublishedCatalog|null>(null),[currentKnown,setCurrentKnown]=useState(false);
 const [stage,setStage]=useState<typeof stages[number]>("legal"),[reason,setReason]=useState(""),[evidence,setEvidence]=useState(""),[publishReason,setPublishReason]=useState("");
 const key=`elite-catalog:${scope}`;
 function persist(value:Capsule){const raw=JSON.stringify(capsule.parse(value));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage");saved.current=value;setPending(value.pending)}
 function accepted(v:CatalogReceipt){
  const next={...saved.current,pending:null,...(v.kind==="source"?{source:v}:{}),...(v.kind==="media"?{media:v}:{})};persist(next);
  if(v.kind==="source"){setModelID(v.resource_id);setBookID(v.source_price_book_id!)}
  if(v.kind==="media")setMediaID(v.resource_id);
  if(v.kind==="draft"){setDraftID(v.resource_id);setDraft(null)}
 }
 useEffect(()=>{try{
  const raw=sessionStorage.getItem(key);if(raw){const v=capsule.parse(JSON.parse(raw));if(v.source&&v.source.actor!==subject||v.media&&v.media.actor!==subject)throw new Error("scope");saved.current=v;setPending(v.pending);if(v.source){setModelID(v.source.resource_id);setBookID(v.source.source_price_book_id!)}if(v.media)setMediaID(v.media.resource_id)}
  setStorageReady(true);
 }catch{setNotice(t("p0148"))}},[key,subject]);
 async function get(kind:string,id?:string){const r=await fetch(endpoint+"?"+new URLSearchParams({kind,organization_id:organization,...(id?{id}:{})}),{cache:"no-store"});if(!r.ok)throw new Error("unavailable");return r.json() as Promise<unknown>}
 async function readDraftValue(id:string){const value=catalogDraft.parse(await get("draft",id));if(value.id!==id||value.snapshot.profile.tenant_id!==tenant||value.snapshot.profile.organization_id!==organization)throw new Error("scope");return value}
 async function loadDraft(){
  if(held.current)return;held.current=true;setBusy(true);setDraft(null);
  try{setDraft(await readDraftValue(draftID));setNotice(t("p0149"))}catch{setNotice(t("p0150"))}finally{held.current=false;setBusy(false)}
 }
 async function loadCurrent(){
  if(held.current)return;held.current=true;setBusy(true);setCurrentKnown(false);
  try{const v=z.object({publication:publishedCatalogSchema.nullable()}).strict().parse(await get("current"));setCurrent(v.publication);setCurrentKnown(true);setNotice(t("p0151"))}catch{setNotice(t("p0152"))}finally{held.current=false;setBusy(false)}
 }
 async function send(c:CatalogCommand|{action:"media";command_id:string;bytes:Uint8Array}){
  if(held.current||!storageReady||saved.current.pending)return;held.current=true;setBusy(true);setNotice("");
  let started=false;
  try{
   let ref:PendingCatalog,url=endpoint,body:BodyInit,type:string;
   if(c.action==="media"){
    const original=await catalogSHA(c.bytes);ref={command_id:c.command_id,kind:"media",request_sha256:await catalogSHA(catalogCanonical({command_id:c.command_id,original_sha256:original}))};
    url+="?"+new URLSearchParams({organization_id:organization,command_id:c.command_id});body=new Uint8Array(c.bytes).buffer;type="image/png";
   }else{const parsed=catalogCommand.parse(c);ref=await commandReference(parsed);body=JSON.stringify(parsed);type="application/json"}
   persist({...saved.current,pending:ref});started=true;
   const r=await fetch(url,{method:"POST",headers:{"content-type":type},body});
   if(!r.ok){
    if([400,401,403,413,415].includes(r.status)){persist({...saved.current,pending:null});started=false;throw new Error("rejected-before-write")}
    throw new Error("unconfirmed");
   }
   const v=catalogReceipt.parse(await r.json());if(!receiptMatches(v,ref,subject))throw new Error("receipt identity");
   accepted(v);started=false;setNotice(v.kind==="source"?t("p0153"):v.kind==="media"?t("p0154"):v.kind==="draft"?t("p0155"):v.kind==="review"?t("p0156"):t("p0157"));
   if(v.kind==="review"){try{setDraft(await readDraftValue(v.resource_id))}catch{setDraft(null);setNotice(t("p0158"))}}
   if(v.kind==="publish"){setCurrentKnown(false);setCurrent(null)}
  }catch{
   setNotice(started?t("p0159"):t("p0160"));
  }finally{held.current=false;setBusy(false)}
 }
 async function recover(){
  const p=saved.current.pending;if(!p||held.current)return;held.current=true;setBusy(true);
  try{const value=catalogReceipt.parse(await get("command",p.command_id));if(!receiptMatches(value,p,subject))throw new Error("identity");accepted(value);setNotice(t("p0161"));if(value.kind==="review"){setDraft(null);setDraftID(value.resource_id)}if(value.kind==="publish"){setCurrentKnown(false);setCurrent(null)}}
  catch{setNotice(t("p0162"))}finally{held.current=false;setBusy(false)}
 }
 const blocked=busy||!storageReady||pending!==null;
 const prevent=(e:FormEvent)=>e.preventDefault();
 const source=async(e:FormEvent)=>{prevent(e);try{
  const value=catalogCommand.parse({action:"source",organization_id:organization,command_id:crypto.randomUUID(),model:{id:"",code:modelCode,displayName:modelName,vehicleClass,specification:{description}},variants:variants.map(v=>({code:v.code,display_name:v.name,battery_specification:{description:v.battery},amount_minor_units:v.amount,tax_mode:v.tax})),valid_from:new Date(from+"Z").toISOString(),valid_until:until?new Date(until+"Z").toISOString():null});
  await send(value);
 }catch{setNotice(t("p0163"))}};
 function changeVariant(i:number,key:keyof VariantInput,value:string){setVariants(old=>old.map((v,j)=>j===i?{...v,[key]:value}:v))}
 const draftAction=(e:FormEvent)=>{prevent(e);void send({action:"draft",organization_id:organization,command_id:crypto.randomUUID(),price_book_id:bookID,models:[{model_id:modelID,media_id:mediaID}]})};
 async function upload(e:FormEvent){prevent(e);if(!file||file.type!=="image/png"||file.size<1||file.size>1048576){setNotice(t("p0164"));return}await send({action:"media",command_id:crypto.randomUUID(),bytes:new Uint8Array(await file.arrayBuffer())})}
 async function evidenceFile(v:File|null){setEvidence("");if(!v||v.size<1||v.size>1048576){setNotice(t("p0165"));return}setEvidence(await catalogSHA(new Uint8Array(await v.arrayBuffer())))}
 return <>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["catalog-role-view"]}/>
 <p role="status" aria-live="polite">{notice}</p>
 {pending&&<section className="card" aria-label={t("p0166")}><h2>{t("p0166")}</h2><p>{t("p0167")} <code>{pending.command_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 {can("catalog:draft")&&<section className="card"><h2>{t("p0169")}</h2><p>{t("p0170")}</p>
 <form onSubmit={source}><fieldset disabled={blocked}>
 <label>{t("p0171")}<input value={modelCode} onChange={e=>setModelCode(e.target.value)} required maxLength={64}/></label>
 <label>{t("p0172")}<input value={modelName} onChange={e=>setModelName(e.target.value)} required maxLength={160}/></label>
 <label>{t("p0173")}<select value={vehicleClass} onChange={e=>setVehicleClass(e.target.value)}><option value="bicycle">{t("p0174")}</option><option value="motorcycle">{t("p0175")}</option><option value="scooter">{t("p0176")}</option><option value="utility">{t("p0177")}</option><option value="other">{t("p0178")}</option></select></label>
 <label>{t("p0179")}<textarea value={description} onChange={e=>setDescription(e.target.value)} required maxLength={4000}/></label>
 <label>{t("p0180")}<input type="datetime-local" value={from} onChange={e=>setFrom(e.target.value)} required/></label>
 <label>{t("p0181")}<input type="datetime-local" value={until} onChange={e=>setUntil(e.target.value)}/></label>
 {variants.map((v,i)=><fieldset key={i}><legend>{t("p0182")} {i+1}</legend>
 <label>{t("p0183")}<input value={v.code} onChange={e=>changeVariant(i,"code",e.target.value)} required maxLength={64}/></label>
 <label>{t("p0184")}<input value={v.name} onChange={e=>changeVariant(i,"name",e.target.value)} required maxLength={160}/></label>
 <label>{t("p0185")}<textarea value={v.battery} onChange={e=>changeVariant(i,"battery",e.target.value)} required maxLength={4000}/></label>
 <label>{t("p0186")}<input inputMode="numeric" pattern="[0-9]+" value={v.amount} onChange={e=>changeVariant(i,"amount",e.target.value)} required maxLength={19}/></label>
 <label>{t("p0187")}<select value={v.tax} onChange={e=>changeVariant(i,"tax",e.target.value)}><option value="not-applicable">{t("p0188")}</option><option value="inclusive">{t("p0189")}</option><option value="exclusive">{t("p0190")}</option></select></label>
 {variants.length>1&&<button type="button" onClick={()=>setVariants(old=>old.filter((_,j)=>i!==j))}>{t("p0191")} {i+1}</button>}
 </fieldset>)}
 <button type="button" className="button" disabled={variants.length>=16} onClick={()=>setVariants(old=>[...old,emptyVariant()])}>{t("p0192")}</button>
 <button className="button">{t("p0193")}</button>
 </fieldset></form></section>}
 {can("catalog:draft")&&<section className="card"><h2>{t("p0194")}</h2><form onSubmit={upload}><label>{t("p0195")}<input type="file" accept="image/png" disabled={blocked} onChange={e=>setFile(e.target.files?.[0]??null)} required/></label><button className="button" disabled={blocked||!file}>{t("p0196")}</button></form>
 <form onSubmit={draftAction}><fieldset disabled={blocked}>
 <label>{t("p0197")}<input value={modelID} onChange={e=>setModelID(e.target.value)} required maxLength={64}/></label>
 <label>{t("p0198")}<input value={mediaID} onChange={e=>setMediaID(e.target.value)} required maxLength={64}/></label>
 <label>{t("p0199")}<input value={bookID} onChange={e=>setBookID(e.target.value)} required maxLength={64}/></label>
 <button className="button">{t("p0200")}</button></fieldset></form></section>}
 <section className="card"><h2>{t("p0201")}</h2>
 <form onSubmit={e=>{prevent(e);void loadDraft()}}><label>{t("p0202")}<input value={draftID} onChange={e=>{setDraftID(e.target.value);setDraft(null)}} required maxLength={64}/></label><button className="button" disabled={busy}>{t("p0203")}</button></form>
 {draft&&<><p>{t("p0204")} {draft.maker}{t("p0008")}</p><a href={`/admin/catalog?draft=${encodeURIComponent(draft.id)}`}>{t("p0205")}</a>
 {draft.snapshot.models.map(m=><article key={m.id}><h3>{m.displayName}</h3><p>{t("p0206")} {m.code}{t("p0207")} {m.vehicleClass}{t("p0008")}</p><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{JSON.stringify(m.specification,null,2)}</pre></article>)}
 {draft.snapshot.variants.map(v=><article key={v.id}><h3>{v.display_name}</h3><p>{t("p0208")} {v.amount_minor_units} {t("p0209")} {draft.snapshot.profile.currency}{t("p0210")} {v.tax_mode}{t("p0211")} {v.homologation_state}{t("p0008")}</p></article>)}
 <ul aria-label={t("p0212")}>{stages.map(s=><li key={s}>{stageLabel[s]}{t("p0213")} {draft.reviews[s]}</li>)}</ul>
 <details><summary>{t("p0214")}</summary><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{JSON.stringify(draft.snapshot,null,2)}</pre><p>{t("p0215")} <code>{draft.sha256}</code></p></details>
 {draft.maker===subject?<p>{t("p0216")}</p>:<><label>{t("p0217")}<select value={stage} onChange={e=>setStage(e.target.value as typeof stage)}>{stages.filter(s=>can("catalog:review:"+s)).map(s=><option key={s} value={s}>{stageLabel[s]}</option>)}</select></label>
 <label>{t("p0218")}<input type="file" disabled={blocked} onChange={e=>void evidenceFile(e.target.files?.[0]??null)}/></label><p>{t("p0219")}</p>
 <label>{t("p0220")}<textarea value={reason} onChange={e=>setReason(e.target.value)} maxLength={2000}/></label>
 {[true,false].map(approved=><button key={String(approved)} className="button" disabled={blocked||!can("catalog:review:"+stage)||!evidence||!reason.trim()||draft.reviews[stage]!=="pending"||approved&&stage==="publication"&&stages.some(s=>s!=="publication"&&draft.reviews[s]!=="approved")} onClick={()=>void send({action:"review",organization_id:organization,command_id:crypto.randomUUID(),draft_id:draft.id,stage,snapshot_sha256:draft.sha256,approved,reason,evidence_sha256:evidence})}>{approved?t("p0221"):t("p0222")}</button>)}
 </>}
 </>}
 </section>
 <section className="card"><h2>{t("p0223")}</h2><button className="button" disabled={busy} onClick={()=>void loadCurrent()}>{t("p0224")}</button>
 {currentKnown&&(current?<><p>{t("p0225")} {current.generation}{t("p0008")}</p><ul>{current.models.map(m=><li key={m.id}>{m.displayName}</li>)}</ul><a href="/models">{t("p0226")}</a></>:<p>{t("p0227")}</p>)}
 {draft&&can("catalog:publish")&&draft.maker!==subject&&<><p>{t("p0228")}</p>
 <label>{t("p0229")}<textarea value={publishReason} onChange={e=>setPublishReason(e.target.value)} maxLength={2000}/></label>
 <button className="button" disabled={blocked||!currentKnown||!publishReason.trim()||stages.some(s=>draft.reviews[s]!=="approved")} onClick={()=>void send({action:"publish",organization_id:organization,command_id:crypto.randomUUID(),draft_id:draft.id,snapshot_sha256:draft.sha256,expected_generation:current?.generation??"0",reason:publishReason})}>{t("p0230")}</button>
 </>}
 </section>
 </>;
}
````

### FILE: `src/platform/catalog/authoring.test.ts`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "88fccc10c25271a1d886390490c35e2fb69daeee9513f8b0636e124a73db9d77"
variables: []
secrets_allowed: false
```

````typescript
import {readFileSync} from "node:fs";
import {expect,it} from "vitest";
import {catalogCanonical,catalogSHA,catalogCommand,commandPayload,receiptMatches} from "./authoring";
const cases=JSON.parse(readFileSync("deploy/catalog/authoring-goldens.json","utf8")) as {name:string;payload:Record<string,unknown>;canonical:string;sha256:string}[];
it.each(cases)("matches typed Go $name hashes without narrowing int64 values",async c=>{
 const body=c.name==="media"?c.payload:commandPayload(catalogCommand.parse({action:c.name,organization_id:"org",...c.payload}));
 expect(catalogCanonical(body)).toBe(c.canonical);expect(await catalogSHA(catalogCanonical(body))).toBe(c.sha256);
});
it("normalizes source instants to the same Go wire representation",()=>{
 const source=cases.find(c=>c.name==="source")!;
 const c=catalogCommand.parse({action:"source",organization_id:"org",...source.payload,valid_from:"2026-09-12T00:00:00.000Z",valid_until:"2026-10-12T00:00:00.120Z"});
 expect(catalogCanonical(commandPayload(c))).toBe(source.canonical);
});
it("never accepts a receipt for another actor, payload or draft",()=>{
 const p={command_id:"command",kind:"publish" as const,request_sha256:"a".repeat(64),resource_id:"draft",snapshot_sha256:"b".repeat(64)};
 const receipt={...p,actor:"reviewer",generation:"1",replay:false,effective_price_book_id:"book"};
 expect(receiptMatches(receipt,p,"reviewer")).toBe(true);
 for(const change of [{actor:"other"},{request_sha256:"c".repeat(64)},{resource_id:"other"},{snapshot_sha256:"c".repeat(64)}])expect(receiptMatches({...receipt,...change},p,"reviewer")).toBe(false);
});
it("rejects unsupported transport shape and invalid source before a write",()=>{
 const source=cases.find(c=>c.name==="source")!;
 const input={action:"source",organization_id:"org",...source.payload};
 expect(catalogCommand.safeParse({...input,extra:true}).success).toBe(false);
 expect(()=>catalogCanonical({amount:42})).toThrow();
 expect(()=>catalogCanonical({"non-ascii-é":"value"})).toThrow();
});
````

### FILE: `src/platform/catalog/authoring.ts`

```yaml
block_id: "TS-CATALOG-AUTHORING-PORTAL:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e08488800f393611e8ff220c498866f85d3335cf97b79d70b5882c76c380a67e"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED typed form/receipt binding. Existing Go owners govern all effects.
import {z} from "zod";
import {publishedCatalogSchema} from "./schema";
export const catalogID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$/u);
const sha=z.string().regex(/^[0-9a-f]{64}$/u);
const code=z.string().max(64).regex(/^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/u);
const text=(max:number)=>z.string().min(1).max(max).refine(v=>v.trim()===v&&!/[\x00-\x1f\x7f]/u.test(v)&&new TextDecoder().decode(new TextEncoder().encode(v))===v);
const unsigned=z.string().regex(/^(0|[1-9][0-9]{0,18})$/u).refine(v=>BigInt(v)<=9223372036854775807n);
const instant=z.iso.datetime({offset:true}).transform(v=>new Date(v).toISOString().replace(/\.000Z$/u,"Z").replace(/(\.\d*?[1-9])0+Z$/u,"$1Z"));
export const reviewStage=z.enum(["legal","technical","media","publication"]);
const common={organization_id:catalogID,command_id:catalogID};
export const sourceForm=z.object({
 action:z.literal("source"),...common,
 model:z.object({id:z.literal(""),code,displayName:text(160).refine(v=>v.length>=2),vehicleClass:z.enum(["motorcycle","bicycle","scooter","utility","other"]),specification:z.object({description:text(4000)}).strict()}).strict(),
 variants:z.array(z.object({code,display_name:text(160),battery_specification:z.object({description:text(4000)}).strict(),amount_minor_units:unsigned,tax_mode:z.enum(["inclusive","exclusive","not-applicable"])}).strict()).min(1).max(16),
 valid_from:instant,valid_until:instant.nullable()
}).strict().refine(v=>v.valid_until===null||Date.parse(v.valid_until)>Date.parse(v.valid_from)).refine(v=>new Set(v.variants.map(x=>x.code)).size===v.variants.length);
export const catalogCommand=z.discriminatedUnion("action",[
 sourceForm,
 z.object({action:z.literal("draft"),...common,price_book_id:catalogID,models:z.array(z.object({model_id:catalogID,media_id:catalogID}).strict()).min(1).max(32)}).strict(),
 z.object({action:z.literal("review"),...common,draft_id:catalogID,stage:reviewStage,snapshot_sha256:sha,approved:z.boolean(),reason:text(2000),evidence_sha256:sha}).strict(),
 z.object({action:z.literal("publish"),...common,draft_id:catalogID,snapshot_sha256:sha,expected_generation:unsigned,reason:text(2000)}).strict()
]);
export type CatalogCommand=z.infer<typeof catalogCommand>;
export const catalogReceipt=z.object({command_id:catalogID,actor:z.string().min(1).max(256),kind:z.enum(["media","source","draft","review","publish"]),resource_id:catalogID,request_sha256:sha,snapshot_sha256:sha.optional(),generation:unsigned,replay:z.boolean(),effective_price_book_id:catalogID.optional(),source_price_book_id:catalogID.optional(),source_variant_ids:z.array(catalogID).min(1).max(16).optional()}).strict().superRefine((v,c)=>{
 if(v.kind==="source"&&(!v.source_price_book_id||!v.source_variant_ids))c.addIssue({code:"custom",message:"missing source references"});
 if(["draft","review","publish"].includes(v.kind)&&!v.snapshot_sha256)c.addIssue({code:"custom",message:"missing snapshot"});
 if(v.kind==="publish"&&(!v.effective_price_book_id||v.generation==="0"))c.addIssue({code:"custom",message:"missing publication"});
});
export type CatalogReceipt=z.infer<typeof catalogReceipt>;
const snapshot=z.object({tenant_code:code,schema:z.literal("elite-catalog-snapshot/v1"),profile:z.object({schema:z.literal("elite-catalog-publication/v1"),tenant_id:z.uuid(),organization_id:catalogID,market:z.string().regex(/^[A-Z]{2}$/u),currency:z.string().regex(/^[A-Z]{3}$/u),origin:z.url(),policy_code:z.literal("catalog-reference-v1")}).strict(),price_book_id:catalogID,valid_from:z.iso.datetime({offset:true}),valid_until:z.iso.datetime({offset:true}).nullable(),models:publishedCatalogSchema.shape.models,variants:publishedCatalogSchema.shape.variants}).strict();
export const catalogDraft=z.object({id:catalogID,maker:z.string().min(1).max(256),sha256:sha,snapshot,reviews:z.object({legal:z.enum(["pending","approved","rejected"]),technical:z.enum(["pending","approved","rejected"]),media:z.enum(["pending","approved","rejected"]),publication:z.enum(["pending","approved","rejected"])}).strict()}).strict();
export type CatalogDraft=z.infer<typeof catalogDraft>;
export const pendingCatalog=z.object({command_id:catalogID,kind:catalogReceipt.shape.kind,request_sha256:sha,resource_id:catalogID.optional(),snapshot_sha256:sha.optional()}).strict();
export type PendingCatalog=z.infer<typeof pendingCatalog>;
export function commandPayload(c:CatalogCommand){const{action:_,organization_id:__,...body}=c;return body}
// Transport correspondence with Go encoding/json for this explicitly bounded
// string/boolean/null schema and ASCII object keys; not a generic JSON/JCS engine.
export function catalogCanonical(value:unknown):string{
 const literal=(v:string)=>JSON.stringify(v).replace(/[<>&\u2028\u2029]/gu,x=>({"<":"\\u003c",">":"\\u003e","&":"\\u0026","\u2028":"\\u2028","\u2029":"\\u2029"}[x]!));
 if(value===null)return"null";
 if(typeof value==="string")return literal(value);
 if(typeof value==="boolean")return String(value);
 if(Array.isArray(value))return"["+value.map(catalogCanonical).join(",")+"]";
 if(typeof value==="object"){
  const o=value as Record<string,unknown>;const keys=Object.keys(o).sort();
  if(keys.some(k=>!/^[A-Za-z_][A-Za-z0-9_]*$/u.test(k)))throw new Error("unsupported catalog key");
  return"{"+keys.map(k=>literal(k)+":"+catalogCanonical(o[k])).join(",")+"}";
 }
 throw new Error("unsupported catalog value");
}
export async function catalogSHA(value:Uint8Array|string){const bytes=typeof value==="string"?new TextEncoder().encode(value):value;return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",new Uint8Array(bytes).buffer)),x=>x.toString(16).padStart(2,"0")).join("")}
export async function commandReference(c:CatalogCommand):Promise<PendingCatalog>{const body=commandPayload(c);return pendingCatalog.parse({command_id:c.command_id,kind:c.action,request_sha256:await catalogSHA(catalogCanonical(body)),...("draft_id"in c?{resource_id:c.draft_id,snapshot_sha256:c.snapshot_sha256}:{})})}
export function receiptMatches(v:CatalogReceipt,p:PendingCatalog,subject:string){return v.command_id===p.command_id&&v.actor===subject&&v.kind===p.kind&&v.request_sha256===p.request_sha256&&(!p.resource_id||p.resource_id===v.resource_id)&&(!p.snapshot_sha256||p.snapshot_sha256===v.snapshot_sha256)}
````

## 6. Configuration surface

docs/CATALOG_ROLE_AUTHORING_REFERENCE.md; opt-in features.catalog_editor with catalog:read and independent draft/review/publish permissions. Go host also requires migration0076. Exact existing catalog profile and provider conditions retained.

## 7. Dependency bill

20new/10changed AUTHORED unavoidable transport/form/configuration/transaction/test glue; no new upstream or dependency. Existing ADAPTED business sources and pinned runtimes/licenses retained; no corporate attribution.

## 8. Apply order

Select TS-GO-API-WEB-BRIDGE, TS-OIDC-PORTAL-ADAPTER, TS-FRANCHISE-JOURNEY-PORTALS and TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT. Go source/publication owners required for enabled catalog; narrow web only proves import/type closure. Shared goldens live in this pack.

## 9. Verification

Actual empty-catalog role browser:2models/2variants/5pricebooks/1media/2snapshots/8reviews/3publications/16commands. Source and publication response losses recover via GET with zero extra POST. Four concurrent source requests create one aggregate;6table rollback and durable receipt recovery.14focused web tests, strict types/build,4Go/TS command goldens, source-only finite fuzz. Host requires migration; populated down refuses16commands/3publications/2models, empty73migrations/down/up PASS.

## 10. Reconstruction evidence

CATALOG_ROLE_AUTHORING_RELEASE_V402.md/json binds exact source/delta, all RED/PASS receipts, reconstruction and notices. Supply/warranty/J5/CMS/KPIs/private locale remain explicit T2804 owners. No whole T2804, TEST02 or overall readiness promotion.


V402 composed delta: T2804 help CMS and21same-release guides;15oldguides unchanged,5bounded training curricula revision2, no automatic grants. New optional host, shared inline text, actual browser/PG proofs. HELP_CMS_RELEASE_V402.md.

V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.
