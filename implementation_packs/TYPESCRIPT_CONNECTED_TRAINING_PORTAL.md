# Connected training portal

## 1. Metadata

```yaml
pack_id: "TS-CONNECTED-TRAINING-PORTAL"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Explicit same-release training content through durable participation and distinct human assessment; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804/J5 slice. Compose existing help, shared approval, audit/outbox and OIDC owners. Four role curricula guide public content; permissions remain independently administered. No automatic score, labor qualification, employee creation or grant.

## 3. Architecture contract

Immutable profile/content/actor/org identity; append-only participation facts in audit.event, assessment in existing approval.request/decision, atomic existing outbox. One assessment per attempt. Physical guards reject generic approval bypass and audit mutation. Same-release snapshot survives content revisions.

## 4. Exact file manifest

```text
CREATE config/training.browser.fixture.json
CREATE docs/TRAINING_WEB_REFERENCE.md
CREATE microsoft_playwright_browser_gate/tests/training-connected.spec.mjs
CREATE src/app/api/enterprise/training/route.test.ts
CREATE src/app/api/enterprise/training/route.ts
CREATE src/app/guide/training/page.tsx
CREATE src/components/training-workspace.tsx
CREATE src/platform/training/contract.ts
CREATE tools/export-training-content.mjs
CREATE src/platform/i18n/private-training-binding.json
CREATE src/platform/i18n/private-training-display.ts
```

## 5. Materialization blocks

### FILE: `config/training.browser.fixture.json`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3bdc6d006ed4c95813cc98a61fb4f4d60f28d7a6e6d34637ff704f046d77f119"
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
    "training_portal": true,
    "role_workspace": true
  }
}
````

### FILE: `docs/TRAINING_WEB_REFERENCE.md`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5be74c84fab09e80645980ea9d8830edc4a44dd1e1bf1eb4170b037ae5928eba"
variables: []
secrets_allowed: false
```

````markdown
# Portal de capacitación — misma release

Seleccionar TS-GO-API-WEB-BRIDGE, TS-OIDC-PORTAL-ADAPTER,
TS-FRANCHISE-JOURNEY-PORTALS y este pack. El perfil completo incorpora el módulo
Go de capacitación. Configurar features.training_portal=true explícitamente.
La navegación sólo muestra /guide/training a training:learn o training:review;
las comprobaciones del servidor siguen siendo obligatorias.

El BFF valida sesión/permisos antes de leer el cuerpo. Conserva el límite temporal
del boundedCommandBody existente y selecciona un tamaño acotado para respuestas.
Sólo permite acciones tipadas y rutas Go fijas. No toma URLs arbitrarias del cliente.
Los formularios de lectura, respuesta y evaluación escriben en los owners durables.
Ante pérdida de respuesta, consultar la referencia guardada; no repetir el POST.
sessionStorage guarda sólo identidad del intento/curso/hash, nunca respuestas.
Las guías son públicas. Una evaluación permanece sujeta a revisión humana.

Con las dependencias fijadas, reconstruir el bundle de ayuda a un destino ausente:
node tools/export-training-content.mjs src/platform/help/content.ts <ruta-absoluta-output.json>
Comparar content_sha256 y source_sha256 con el perfil Go antes de activarlo.
Un cambio de guías produce una nueva revisión y evidencia de contenido.
Para compilar con node_modules externo enlazado: next build --webpack.
La receta existente evita la restricción de rutas del builder Turbopack.

El fixture config/training.browser.fixture.json sólo activa la UI en pruebas;
no configura IdP ni concede permisos. La prueba real usa sesiones JWE y tokens
RS256/JWKS sintéticos, alumno/revisor/lector/otro alumno/otra organización.
Un resultado de capacitación no concede permisos ni acredita aptitud laboral.
````

### FILE: `microsoft_playwright_browser_gate/tests/training-connected.spec.mjs`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7d6e1fc4766026526f1bbed5a70f0d51444af253391399cda6f876b897eaf087"
variables: []
secrets_allowed: false
```

````text
import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
test('role training records exact answers and human assessment with read-only recovery after response loss',async({page,context},info)=>{
 if(process.env.ELITE_TRAINING_BROWSER!=='1')throw new Error('explicit local training fixture required');
 page.setDefaultTimeout(12000);
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const identities=JSON.parse(process.env.ELITE_TRAINING_IDENTITIES),require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[];page.on('pageerror',e=>errors.push(e.message));let startPosts=0,decisionPosts=0,reference;
 page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/training')){const body=r.postDataJSON();if(body.action==='start'){startPosts++;reference=body};if(body.action==='assess')decisionPosts++}});
 await identity('reader');await page.goto('/guide/employee');await expect(page.getByRole('link',{name:'Capacitación',exact:true})).toHaveCount(0);await page.goto('/guide/training');await expect(page.getByText('Tu sesión no permite acceder a capacitación.')).toBeVisible();
 await identity('foreign-org');await page.goto('/guide/training');await expect(page.getByRole('alert').filter({hasText:'No pudimos consultar el contenido'})).toBeVisible();
 await identity('learner');await page.goto('/guide/employee');await page.getByRole('link',{name:'Capacitación',exact:true}).first().click();await expect(page.getByRole('heading',{name:'Capacitación y evaluación'})).toBeVisible();
 await page.route('**/api/enterprise/training',async route=>{if(route.request().method()!=='POST'||route.request().postDataJSON().action!=='start'){await route.continue();return};const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')});
 await page.getByRole('button',{name:'Iniciar práctica',exact:true}).click();await expect(page.getByRole('status')).toContainText('La respuesta no quedó confirmada');expect(startPosts).toBe(1);await page.unrouteAll();
 await page.getByRole('button',{name:'Consultar intento guardado'}).click();await expect(page.getByRole('status')).toContainText('Estado consultado');expect(startPosts).toBe(1);
 const learning=page.getByRole('region',{name:'Mi capacitación'});await expect(learning).toContainText('Registrar un recurso · versión 1.0.0');await learning.getByRole('button',{name:'Registrar lectura',exact:true}).click();await expect(learning.getByText('Lectura registrada',{exact:true})).toBeVisible();
 await learning.getByLabel('Describí qué harías si se pierde la respuesta después de registrar un recurso.').fill('Consulto la referencia guardada y no creo otro recurso para forzar un reintento.');await learning.getByRole('button',{name:'Presentar respuestas'}).click();await expect(learning).toContainText('Pendiente de revisión humana');
 const recorded=await context.request.get(base+'/api/enterprise/training?attempt_id='+reference.attempt_id);expect(recorded.status()).toBe(200);const assessment=(await recorded.json()).assessment;
 const decision={action:'assess',request_id:assessment.request_id,payload_sha256:assessment.payload_sha256,approved:true,reason:'Intento de autorrevisión'};
 expect((await context.request.post(base+'/api/enterprise/training',{headers:{origin:base,'content-type':'application/json'},data:decision})).status()).toBe(403);
 expect((await context.request.post(base+'/api/enterprise/training',{headers:{origin:'https://example.invalid','content-type':'application/json'},data:decision})).status()).toBe(403);
 await identity('other-learner');expect((await context.request.get(base+'/api/enterprise/training?attempt_id='+reference.attempt_id)).status()).toBe(409);
 await identity('reviewer');await page.goto('/guide/training');const review=page.getByRole('region',{name:'Evaluaciones registradas'}).locator('article');await expect(review).toHaveCount(1);await expect(review).toContainText('Consulto la referencia guardada');await review.getByRole('button',{name:'Registrar evaluación humana'}).waitFor();
 expect((await context.request.post(base+'/api/enterprise/training',{headers:{origin:base,'content-type':'application/json'},data:{...decision,grant_permissions:['resource:manage']}})).status()).toBe(400);
 await review.getByLabel('Fundamento de la evaluación').fill('La respuesta describe la consulta de la referencia y evita duplicar el recurso.');await review.getByRole('combobox',{name:'Resultado',exact:true}).selectOption('approve');
 await page.route('**/api/enterprise/training',async route=>{if(route.request().method()!=='POST'||route.request().postDataJSON().action!=='assess'){await route.continue();return};const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')});
 await review.getByRole('button',{name:'Registrar evaluación humana'}).click();await expect(page.getByRole('status')).toContainText('La respuesta no quedó confirmada');expect(decisionPosts).toBe(1);await page.unrouteAll();await page.getByRole('link',{name:'Consultar evaluaciones',exact:true}).click();await expect(review).toContainText('Evaluación favorable');await expect(review.getByRole('button',{name:'Registrar evaluación humana'})).toHaveCount(0);expect(decisionPosts).toBe(1);
 await page.screenshot({path:info.outputPath('training-review-desktop.png'),fullPage:true});
 await identity('learner');await page.goto('/guide/training');await page.getByRole('button',{name:'Consultar intento guardado'}).click();await expect(learning).toContainText('Evaluación favorable');await expect(learning).toContainText('Revisó: reviewer');
 expect((await context.request.post(base+'/api/enterprise/training',{headers:{origin:base,'content-type':'application/json'},data:decision})).status()).toBe(403);
 await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('training-result-mobile.png'),fullPage:true});expect(errors).toEqual([]);
});
````

### FILE: `src/app/api/enterprise/training/route.test.ts`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "76476e260f389506c23686e2569d83a63dfbec8e601abfae5fb5ecaea4216755"
variables: []
secrets_allowed: false
```

````typescript
import{expect,it,vi,afterEach}from"vitest";
const m=vi.hoisted(()=>({get:vi.fn(),post:vi.fn(),session:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:async()=>({features:{training_portal:true}})}));
import{POST}from"./route";
afterEach(()=>vi.clearAllMocks());
it("bounds decoded chunked training input before downstream and cancels the stream",async()=>{
 m.session.mockResolvedValue({subject:"learner",permissions:["training:learn"]});
 let pulls=0,cancelled=false;
 const body=new ReadableStream<Uint8Array>({pull(c){pulls++;c.enqueue(new TextEncoder().encode("é".repeat(8193)))},cancel(){cancelled=true}});
 const r=new Request("https://portal.example.test/api/enterprise/training",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body,duplex:"half"}as RequestInit);
 const reply=await POST(r as never);expect(reply.status).toBe(413);expect(cancelled).toBe(true);expect(pulls).toBeLessThanOrEqual(5);
 expect(m.post).not.toHaveBeenCalled();expect(m.get).not.toHaveBeenCalled();
});
it("requires an authenticated role before consuming the body",async()=>{
 m.session.mockResolvedValue(null);
 const r=new Request("https://portal.example.test/api/enterprise/training",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body:"{}"});
 expect((await POST(r as never)).status).toBe(401);expect(r.bodyUsed).toBe(false);expect(m.post).not.toHaveBeenCalled();
});
````

### FILE: `src/app/api/enterprise/training/route.ts`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "89daf639e1b94204ad44a209830fea43a8fdbc2d2a136ed105ba27a73bce4a8e"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse,type NextRequest } from "next/server";
import { allowed,readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { protectedGet,protectedPost } from "@/platform/backend/protected-client";
import { loadBusinessConfig } from "@/platform/config/load";
import { attemptSchema,trainingCommand } from "@/platform/training/contract";
import { boundedCommandBody } from "@/platform/backend/bounded-command";
const response=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
async function enabled(){return (await loadBusinessConfig()).features.training_portal===true}
export async function GET(request:NextRequest){
 if(!await enabled())return response("NOT_FOUND",404);const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 if(!allowed(session,"training:learn"))return response("FORBIDDEN",403);
 const query=new URL(request.url).searchParams;if([...query.keys()].some(x=>x!=="attempt_id")||query.getAll("attempt_id").length!==1)return response("INVALID_QUERY",400);
 const id=query.get("attempt_id")!;if(!/^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/.test(id))return response("INVALID_QUERY",400);
 try{const value=attemptSchema.parse(await protectedGet(session,`/v1/training/attempts/${id}`,{}));if(!session.organizations.includes(value.organization_id)||value.learner_subject!==session.subject)return response("NOT_FOUND",404);return NextResponse.json(value,{headers:{"cache-control":"no-store"}})}catch{return response("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(request:NextRequest){
 if(!await enabled())return response("NOT_FOUND",404);let origin:string;try{origin=applicationBaseUrl().origin}catch{return response("UNAVAILABLE",503)};
 if(request.headers.get("origin")!==origin)return response("CROSS_ORIGIN_REJECTED",403);
 if(request.headers.get("content-type")!=="application/json")return response("UNSUPPORTED_MEDIA_TYPE",415);
 const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 if(!allowed(session,"training:learn")&&!allowed(session,"training:review"))return response("FORBIDDEN",403);
 let raw:string;try{raw=await boundedCommandBody(request,65536)}catch(error){return response(error instanceof Error&&error.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"INVALID_COMMAND",error instanceof Error&&error.message==="BODY_TOO_LARGE"?413:400)}
 let body:unknown;try{body=JSON.parse(raw)}catch{return response("INVALID_COMMAND",400)};
 const parsed=trainingCommand.safeParse(body);if(!parsed.success)return response("INVALID_COMMAND",400);
 const c=parsed.data;
 if(!allowed(session,c.action==="assess"?"training:review":"training:learn"))return response("FORBIDDEN",403);
 let path:string,payload:unknown;
 if(c.action==="start"){path="/v1/training/attempts";payload={attempt_id:c.attempt_id,course_id:c.course_id,profile_sha256:c.profile_sha256}}
 else if(c.action==="acknowledge"){path=`/v1/training/attempts/${c.attempt_id}/acknowledgements`;payload={lesson_id:c.lesson_id,profile_sha256:c.profile_sha256}}
 else if(c.action==="submit"){path=`/v1/training/attempts/${c.attempt_id}/submissions`;payload={answers:c.answers,profile_sha256:c.profile_sha256}}
 else{path=`/v1/training/assessments/${encodeURIComponent(c.request_id)}/decision`;payload={payload_sha256:c.payload_sha256,approved:c.approved,reason:c.reason}}
 try{return NextResponse.json(await protectedPost(session,path,payload),{headers:{"cache-control":"no-store"}})}catch{return response("CONSULT_RECORDED_STATE",409)}
}
````

### FILE: `src/app/guide/training/page.tsx`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9153e7eae2b2a7bad903df05d60db93360c87b3f0eb68de5600a026ffb3c8963"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { createHash } from "node:crypto";
import { redirect,notFound } from "next/navigation";
import type { Route } from "next";
import { allowed,readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { loadBusinessConfig } from "@/platform/config/load";
import { courseView,assessmentSchema } from "@/platform/training/contract";
import { TrainingWorkspace } from "@/components/training-workspace";
export default async function TrainingPage(){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.training_portal!==true)notFound();
 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/guide/training" as Route);
 const learn=allowed(session,"training:learn"),review=allowed(session,"training:review");
 if(!learn&&!review)return <><h1 className="pageTitle">{t("p0105")}</h1><p>{t("p0106")}</p></>;
 try{
  const [courses,assessments]=await Promise.all([protectedGet(session,"/v1/training/courses",{}),protectedGet(session,"/v1/training/assessments",{})]);
  const views=courseView.array().max(8).parse(courses),records=assessmentSchema.array().max(50).parse(assessments);
  if(records.some(x=>!session.organizations.includes(x.payload.organization_id)||(!review&&x.payload.learner_subject!==session.subject)))throw new Error("SCOPE");
  const scope=createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,session.organizations])).digest("hex");
  return <><h1 className="pageTitle">{t("p0107")}</h1><p>{t("p0108")}</p><TrainingWorkspace courses={views} assessments={records} canLearn={learn} canReview={review} scope={scope}/></>
 }catch{return <><h1 className="pageTitle">{t("p0105")}</h1><p role="alert">{t("p0109")}</p><a className="button" href="/guide/training">{t("p0033")}</a></>}
}
````

### FILE: `src/components/training-workspace.tsx`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "96daf8400eec68937a326719be2f0331a9e7fa983a28dd775d9a0550f9cfa449"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {guideDisplay} from "@/platform/i18n/private-guide-display";
import {trainingCourseDisplay} from "@/platform/i18n/private-training-display";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import { useEffect,useRef,useState,type FormEvent } from "react";
import { assessmentSchema,attemptSchema,trainingReference,type CourseView,type Assessment,type Attempt,type TrainingReference } from "@/platform/training/contract";


const endpoint="/api/enterprise/training";
export function TrainingWorkspace({courses,assessments,canLearn,canReview,scope}:{courses:CourseView[];assessments:Assessment[];canLearn:boolean;canReview:boolean;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const labels:Record<string,string>={owner:t("p0903"),admin:t("p0904"),employee:t("p0492"),customer:t("p0125")};

 const [selected,setSelected]=useState(courses[0]?.course.id??""),[reference,setReference]=useState<TrainingReference|null>(null),[attempt,setAttempt]=useState<Attempt|null>(null);
 const [ready,setReady]=useState(false),[busy,setBusy]=useState(false),[reading,setReading]=useState(false),[notice,setNotice]=useState("");
 const fence=useRef(false);const key=`elite-training:${scope}`;
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const saved=trainingReference.parse(JSON.parse(raw));setReference(saved);setSelected(saved.course_id);setNotice(t("p0908"))};setReady(true)}catch{fence.current=true;setNotice(t("p0909"))}},[key]);
 async function consult(ref=reference){if(!ref||reading)return;setReading(true);try{const r=await fetch(`${endpoint}?attempt_id=${encodeURIComponent(ref.attempt_id)}`,{cache:"no-store"});if(!r.ok)throw new Error();const value=attemptSchema.parse(await r.json());if(value.attempt_id!==ref.attempt_id||value.content.profile_sha256!==ref.profile_sha256)throw new Error();setAttempt(value);fence.current=false;setBusy(false);setNotice(t("p0910"))}catch{setNotice(t("p0911"))}finally{setReading(false)}}
 async function act(body:Record<string,unknown>){if(fence.current)return;fence.current=true;setBusy(true);setNotice("");try{const r=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(body)});if(!r.ok)throw new Error();const value:unknown=await r.json();if(body.action==="assess"){assessmentSchema.parse(value);window.location.reload();return};if(body.action==="submit"){assessmentSchema.parse(value);await consult();return};const current=attemptSchema.parse(value);setAttempt(current);fence.current=false;setBusy(false);setNotice(t("p0912"))}catch{setNotice(t("p0913"))}}
 function start(existing?:TrainingReference){const view=courses.find(c=>c.course.id===selected);if(!view||!ready)return;const ref=existing??{attempt_id:crypto.randomUUID(),course_id:view.course.id,profile_sha256:view.profile_sha256};try{sessionStorage.setItem(key,JSON.stringify(ref));setReference(ref)}catch{setNotice(t("p0914"));return};void act({action:"start",...ref})}
 function nextAttempt(){if(!attempt||(!attempt.assessment&&attempt.current_profile))return;try{sessionStorage.removeItem(key);setReference(null);setAttempt(null);fence.current=false;setBusy(false);setNotice(t("p0915"))}catch{setNotice(t("p0916"))}}
 function submit(event:FormEvent<HTMLFormElement>){event.preventDefault();if(!attempt)return;const form=new FormData(event.currentTarget),answers:Record<string,string>={};for(const prompt of attempt.content.course.prompts){const value=String(form.get(prompt.id)??"");if(!value.trim()||new TextEncoder().encode(value).length>2048){setNotice(t("p0917"));return};answers[prompt.id]=value};void act({action:"submit",attempt_id:attempt.attempt_id,profile_sha256:attempt.content.profile_sha256,answers})}
 const course=attempt?.content??courses.find(c=>c.course.id===selected);
 return <>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["training-role-view"]}/>
  <p role="status" aria-live="polite">{notice}</p>
  {canLearn&&<section className="card" aria-label={t("p0918")} style={{overflowWrap:"anywhere"}}><h2>{t("p0918")}</h2>
   {!reference&&<><label>{t("p0919")}<select value={selected} disabled={!ready||busy} onChange={event=>setSelected(event.target.value)}>{courses.map(c=><option key={c.course.id} value={c.course.id}>{labels[c.course.role]} {t("p0016")} {trainingCourseDisplay(c,privateLocale.language).title}</option>)}</select></label><p>{t("p0920")}</p><button className="button" disabled={!ready||busy} onClick={()=>start()}>{t("p0921")}</button></>}
   {reference&&<p><button className="button" disabled={reading} onClick={()=>void consult()}>{t("p0922")}</button>{!attempt&&!busy&&<button className="button" onClick={()=>start(reference)}>{t("p0923")}</button>}</p>}
   {course&&<><h3>{trainingCourseDisplay(course,privateLocale.language).title}</h3><p>{t("p0924")} {course.profile_id}{t("p0925")} {course.profile_revision} {t("p0926")} {labels[course.course.role]}</p></>}
   {attempt&&!attempt.current_profile&&<p role="alert">{t("p0927")}</p>}
   {attempt&&course?.articles.map(source=>{const article=guideDisplay(source,privateLocale.language);return <section lang={article.displayLanguage} key={article.id} aria-label={article.title}><h3>{article.title} {t("p0543")} {article.version}</h3>{article.paragraphs.map((text,index)=><p key={index}>{text}</p>)}{attempt.read_lessons.includes(article.id)?<p>{t("p0928")}</p>:<button className="button" disabled={busy||!attempt.current_profile} onClick={()=>void act({action:"acknowledge",attempt_id:attempt.attempt_id,lesson_id:article.id,profile_sha256:attempt.content.profile_sha256})}>{t("p0929")}</button>}</section>})}
   {attempt&&!attempt.assessment&&attempt.current_profile&&attempt.read_lessons.length===attempt.content.course.lessons.length&&<form onSubmit={submit}><h3>{t("p0930")}</h3>{trainingCourseDisplay(attempt.content,privateLocale.language).prompts.map(p=><label key={p.id} style={{display:"block",marginBottom:16}}>{p.text}<textarea name={p.id} required maxLength={2048} disabled={busy} style={{display:"block",width:"100%",minHeight:100}}/></label>)}<button className="button" disabled={busy}>{t("p0931")}</button></form>}
   {attempt?.assessment&&<AssessmentResult value={attempt.assessment}/>}
   {attempt&&(attempt.assessment?.state==="approved"||attempt.assessment?.state==="rejected"||!attempt.current_profile)&&<p><button className="button" disabled={busy} onClick={nextAttempt}>{t("p0932")}</button></p>}
  </section>}
  <section aria-label={t("p0933")}><h2>{canReview?t("p0934"):t("p0935")}</h2><p><a className="button" href="/guide/training">{t("p0936")}</a></p>
   {assessments.length===0?<p>{t("p0937")}</p>:assessments.map(value=><article className="card" key={value.request_id} style={{marginTop:16,overflowWrap:"anywhere"}}><h3>{trainingCourseDisplay(value.payload.content,privateLocale.language).title}</h3><p>{t("p0938")} {value.payload.learner_subject}</p><p>{t("p0939")} {value.payload.content.profile_revision}</p>
    {value.payload.content.articles.map(source=>{const article=guideDisplay(source,privateLocale.language);return <details lang={article.displayLanguage} key={article.id}><summary>{article.title} {t("p0543")} {article.version}</summary>{article.paragraphs.map((text,i)=><p key={i}>{text}</p>)}</details>})}
    <AssessmentResult value={value}/>
    {canReview&&value.state==="pending"&&<form onSubmit={event=>{event.preventDefault();const form=new FormData(event.currentTarget);void act({action:"assess",request_id:value.request_id,payload_sha256:value.payload_sha256,approved:form.get("decision")==="approve",reason:String(form.get("reason")??"")})}}><label>{t("p0940")}<textarea name="reason" required maxLength={2048} disabled={busy} style={{display:"block",width:"100%",minHeight:90}}/></label><label>{t("p0941")}<select name="decision" disabled={busy}><option value="reject">{t("p0906")}</option><option value="approve">{t("p0905")}</option></select></label><p><button className="button" disabled={busy}>{t("p0942")}</button></p></form>}
   </article>)}
  </section>
 </>;
}
function AssessmentResult({value}:{value:Assessment}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const result=(state:string)=>state==="approved"?t("p0905"):state==="rejected"?t("p0906"):t("p0907");
return <section aria-label={t("p0943")}><p><strong>{result(value.state)}</strong></p>{trainingCourseDisplay(value.payload.content,privateLocale.language).prompts.map(p=><div key={p.id}><p>{p.text}</p><p lang="und" style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{value.payload.answers[p.id]}</p></div>)}{value.reviewer&&<p>{t("p0944")} {value.reviewer}{t("p0008")} {value.reason}</p>}<p>{t("p0945")}</p></section>}
````

### FILE: `src/platform/training/contract.ts`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "16adfa52e4f77c3cfb5d57a247050b971ef725250a5557e1c18474e7df04b3ef"
variables: []
secrets_allowed: false
```

````typescript
import { z } from "zod";
const id=z.string().regex(/^[a-z][a-z0-9-]{0,63}$/),hash=z.string().regex(/^[0-9a-f]{64}$/),attemptId=z.uuid();
const article=z.object({id,version:z.string().regex(/^\d{1,6}\.\d{1,6}\.\d{1,6}$/),title:z.string().min(1).max(160),paragraphs:z.array(z.string().min(1).max(4000)).min(1).max(20)}).strict();
const course=z.object({id,title:z.string().min(1).max(160),role:z.enum(["owner","admin","employee","customer"]),lessons:z.array(id).min(1).max(4),prompts:z.array(z.object({id,text:z.string().min(1).max(1000)}).strict()).min(1).max(8)}).strict();
export const courseView=z.object({course,articles:z.array(article).min(1).max(4),profile_id:id,profile_revision:z.number().int().positive(),profile_sha256:hash,content_sha256:hash,method:z.literal("HUMAN_REVIEW_NO_GRANTS_V1")}).strict();
const answers=z.record(id,z.string().min(1).max(2048)).refine(value=>Object.keys(value).length>=1&&Object.keys(value).length<=8);
export const assessmentSchema=z.object({request_id:z.string().regex(/^training:[0-9a-f]{64}$/),payload_sha256:hash,state:z.enum(["pending","approved","rejected"]),payload:z.object({schema:z.literal("training-assessment/v1"),attempt_id:attemptId,learner_subject:z.string().min(1).max(128),organization_id:z.string().min(1).max(128),content:courseView,answers}).strict(),reviewer:z.string().max(128).optional(),reason:z.string().max(2048).optional(),approved:z.boolean().optional()}).strict();
export const attemptSchema=z.object({attempt_id:attemptId,learner_subject:z.string().min(1).max(128),organization_id:z.string().min(1).max(128),content:courseView,read_lessons:z.array(id).max(4),assessment:assessmentSchema.optional(),current_profile:z.boolean()}).strict();
export const trainingCommand=z.discriminatedUnion("action",[
 z.object({action:z.literal("start"),attempt_id:attemptId,course_id:id,profile_sha256:hash}).strict(),
 z.object({action:z.literal("acknowledge"),attempt_id:attemptId,lesson_id:id,profile_sha256:hash}).strict(),
 z.object({action:z.literal("submit"),attempt_id:attemptId,profile_sha256:hash,answers}).strict(),
 z.object({action:z.literal("assess"),request_id:z.string().regex(/^training:[0-9a-f]{64}$/),payload_sha256:hash,approved:z.boolean(),reason:z.string().trim().min(1).max(2048)}).strict(),
]);
export const trainingReference=z.object({attempt_id:attemptId,course_id:id,profile_sha256:hash}).strict();
export type CourseView=z.infer<typeof courseView>;
export type Assessment=z.infer<typeof assessmentSchema>;
export type Attempt=z.infer<typeof attemptSchema>;
export type TrainingReference=z.infer<typeof trainingReference>;
````

### FILE: `tools/export-training-content.mjs`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "87ceb5257403fcc312978577ce0301e783a2e00a364005f213880b25c6b06bba"
variables: []
secrets_allowed: false
```

````text
// AUTHORED build glue. Reads trusted, version-controlled help source only;
// TypeScript's pinned compiler evaluates no external provider or learner input.
import {readFile,writeFile,mkdtemp} from 'node:fs/promises';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {dirname,join,basename,resolve} from 'node:path';
import {spawnSync} from 'node:child_process';
import {pathToFileURL} from 'node:url';
const [source,out]=process.argv.slice(2);
if(!source||!out)throw new Error('usage: node tools/export-training-content.mjs trusted-help-source.ts absent-output.json');
const raw=await readFile(source);if(raw.length>65536)throw new Error('help source budget exceeded');
const require=createRequire(import.meta.url);const compiler=join(dirname(require.resolve('typescript/package.json')),'bin','tsc');
const compileDir=await mkdtemp(join(dirname(resolve(out)),'.training-content-compile-'));
const platformRoot=resolve(dirname(source),'..');
await writeFile(join(compileDir,'tsconfig.json'),JSON.stringify({compilerOptions:{target:'ES2022',module:'ESNext',moduleResolution:'Bundler',skipLibCheck:true,noEmitOnError:true,rootDir:platformRoot,outDir:compileDir,paths:{'@/*':[resolve(platformRoot,'..','*')]}},files:[resolve(source)]})+'\n',{encoding:'utf8',flag:'wx'});
const build=spawnSync(process.execPath,[compiler,'--project',join(compileDir,'tsconfig.json')],{encoding:'utf8',timeout:20000,windowsHide:true});
if(build.error||build.status!==0)throw new Error('help source compilation failed: '+build.stdout+build.stderr);
const compiledPath=join(compileDir,'help',basename(source).replace(/\.tsx?$/,'.js'));
const compiled=await readFile(compiledPath,'utf8');
const alias='"@/platform/notifications/status-contract"';
if(compiled.split(alias).length!==2)throw new Error('help import contract changed; exporter must be reviewed');
await writeFile(compiledPath,compiled.replace(alias,'"../notifications/status-contract.js"'),{encoding:'utf8'});
const mod=await import(pathToFileURL(compiledPath).href);
if(!Array.isArray(mod.ALL_GUIDES)||mod.ALL_GUIDES.length<1||mod.ALL_GUIDES.length>64)throw new Error('help catalog shape differs');
const articles=mod.ALL_GUIDES.map(({id,version,title,paragraphs})=>({id,version,title,paragraphs}));
const payload={schema:'elite-training-content/v1',source_sha256:createHash('sha256').update(raw).digest('hex'),articles};
const encoded=JSON.stringify(payload,null,2)+'\n';if(Buffer.byteLength(encoded)>65536)throw new Error('content bundle budget exceeded');
await writeFile(out,encoded,{encoding:'utf8',flag:'wx'});
process.stdout.write(JSON.stringify({source_sha256:payload.source_sha256,content_sha256:createHash('sha256').update(encoded).digest('hex'),article_count:articles.length})+'\n');
````

## 6. Configuration surface

docs/TRAINING_WEB_REFERENCE.md; features.training_portal opt-in; session and training permissions required before command body; fixed routes only. source exporter uses pinned compiler and trusted help source; no learner execution. Shared boundedCommandBody default4096 preserved.

## 7. Dependency bill

AUTHORED unavoidable typed configuration, persistence, authorization, UI and orchestration/test glue; no new upstream/dependency or corporate attribution. Existing pgx/UUID, Go/PG, Next/React/Zod/OIDC/Playwright pins and notices retained.

## 8. Apply order

Select TS-GO-API-WEB-BRIDGE, TS-OIDC-PORTAL-ADAPTER and TS-FRANCHISE-JOURNEY-PORTALS. Go counterpart required for enabled training. Narrow web has compile/import closure; only complete reference exercises the connected journey.

## 9. Verification

PG72migrations, concurrent start produces one fact; immutable answers/content and actor/org; self-review denied; generic approval/duplicate assessment and audit mutation blocked. Actual Next/BFF/Chromium learner+reviewer recover lost POST responses by GET:4facts/4outbox/1request/1decision/0resources/0grants. Host/path/hash/index/down proof, profile2s fuzz101053execs, static materialized curriculum matches exercised profile. Focused frontend3tests/full webpack build.

## 10. Reconstruction evidence

TRAINING_CONNECTED_RELEASE_V402.md/json binds source, before/after, RED locator and final PASS receipts, independent profile reconstruction and provenance. This closes training subclaim; other T2804 writes/CMS/private i18n remain open. No global readiness or production claim.


V402 composed delta: T2804 help CMS and21same-release guides;15oldguides unchanged,5bounded training curricula revision2, no automatic grants. New optional host, shared inline text, actual browser/PG proofs. HELP_CMS_RELEASE_V402.md.

V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.

### FILE: `src/platform/i18n/private-training-binding.json`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL-LOCALE-DELTA:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2ed6a184f9149cdf844117d4aa4ad67f493c6cc13156b65cac8a368239d98362"
variables: []
secrets_allowed: false
```

````json
{
  "profile_sha256": "2ce1e88eb44eef172bd4c1d8e09680e2c7c732aa479bee65d855e3251c53f937",
  "profile_id": "reference-onboarding",
  "profile_revision": 2,
  "content_sha256": "eb919b9731d56832c38233572d6f1fd71cc62e7da64ca3d53a97638ab6a3ea20",
  "courses": [
    {
      "source": {
        "id": "resource-onboarding",
        "title": "Registrar recursos y recuperar una respuesta",
        "role": "employee",
        "lessons": [
          "resource-create-view"
        ],
        "prompts": [
          {
            "id": "recovery",
            "text": "Describí qué harías si se pierde la respuesta después de registrar un recurso."
          }
        ]
      },
      "en": {
        "title": "Register resources and recover a response",
        "prompts": [
          {
            "id": "recovery",
            "text": "Describe what you would do if the response is lost after registering a resource."
          }
        ]
      }
    },
    {
      "source": {
        "id": "admin-onboarding",
        "title": "Revisión de recursos",
        "role": "admin",
        "lessons": [
          "resource-create-view",
          "supply-role-view",
          "warranty-role-view",
          "network-role-view"
        ],
        "prompts": [
          {
            "id": "review",
            "text": "Describí cómo recuperarías una operación incierta, revisarías contenido y separarías una evaluación de los permisos de acceso."
          }
        ]
      },
      "en": {
        "title": "Resource review",
        "prompts": [
          {
            "id": "review",
            "text": "Describe how you would recover an uncertain operation, review content and distinguish assessment from access permissions."
          }
        ]
      }
    },
    {
      "source": {
        "id": "owner-onboarding",
        "title": "Supervisión de la operación",
        "role": "owner",
        "lessons": [
          "operation-sections-view",
          "network-role-view",
          "training-role-view"
        ],
        "prompts": [
          {
            "id": "review",
            "text": "Describí cómo distinguís una sección no disponible de una lista vacía."
          }
        ]
      },
      "en": {
        "title": "Operations supervision",
        "prompts": [
          {
            "id": "review",
            "text": "Describe how you distinguish an unavailable section from an empty list."
          }
        ]
      }
    },
    {
      "source": {
        "id": "customer-onboarding",
        "title": "Revisar una cotización",
        "role": "customer",
        "lessons": [
          "quote-acceptance-view"
        ],
        "prompts": [
          {
            "id": "review",
            "text": "Describí qué confirma aceptar una cotización y qué se debe consultar ante una respuesta incierta."
          }
        ]
      },
      "en": {
        "title": "Review a quote",
        "prompts": [
          {
            "id": "review",
            "text": "Describe what accepting a quote confirms and what must be checked after an uncertain response."
          }
        ]
      }
    },
    {
      "source": {
        "id": "content-onboarding",
        "title": "Contenido, publicación y revisión humana",
        "role": "admin",
        "lessons": [
          "help-cms-view",
          "catalog-role-view",
          "training-role-view"
        ],
        "prompts": [
          {
            "id": "review",
            "text": "Explicá cómo publicarías contenido revisado, recuperarías una respuesta incierta y preservarías cursos anteriores sin conceder accesos."
          }
        ]
      },
      "en": {
        "title": "Content, publication and human review",
        "prompts": [
          {
            "id": "review",
            "text": "Explain how you would publish reviewed content, recover an uncertain response and preserve earlier courses without granting access."
          }
        ]
      }
    }
  ]
}
````

### FILE: `src/platform/i18n/private-training-display.ts`

```yaml
block_id: "TS-CONNECTED-TRAINING-PORTAL-LOCALE-DELTA:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7a6798824ca29ffa26accbc6815a913cd89bbf2377ae5ea81a6f6f35bf52d2a7"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED source-bound display only; assessment source and hashes stay unchanged.
import type{PrivateLanguage}from"./private-locale";
import binding from"./private-training-binding.json";
import type{CourseView}from"@/platform/training/contract";
export function trainingCourseDisplay(source:CourseView,language:PrivateLanguage){
 const record=binding.courses.find(x=>x.source.id===source.course.id);
 const same=record&&record.source.title===source.course.title&&record.source.role===source.course.role&&JSON.stringify(record.source.lessons)===JSON.stringify(source.course.lessons)&&record.source.prompts.length===source.course.prompts.length&&record.source.prompts.every((p,i)=>p.id===source.course.prompts[i]?.id&&p.text===source.course.prompts[i]?.text);
 const known=source.profile_sha256===binding.profile_sha256&&source.content_sha256===binding.content_sha256&&source.profile_id===binding.profile_id&&source.profile_revision===binding.profile_revision&&source.method==="HUMAN_REVIEW_NO_GRANTS_V1"&&same;
 return known&&language==="en"?{title:record.en.title,prompts:record.en.prompts,displayLanguage:"en"}:{title:source.course.title,prompts:source.course.prompts,displayLanguage:known?"es":"und"};
}
````

