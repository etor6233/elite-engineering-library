# TypeScript Multi-Role Onboarding

## 1. Metadata

```yaml
pack_id: "TS-MULTIROLE-ONBOARDING"
pack_version: "0.1.8"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Panel y cuatro vistas curadas de navegación, filtradas por sesión y configuración; no concede roles ni certifica capacitación, KPI financieros o activación empresarial."
stacks: ["Next.js 16.3.4", "React 19.2.8", "TypeScript"]
compatible_with: ["TS-GO-API-WEB-BRIDGE 0.5.16", "TS-OIDC-PORTAL-ADAPTER 0.2.4", "TS-FRANCHISE-JOURNEY-PORTALS 0.14.19"]
incompatible_with: ["rol como atajo de autorización (el acceso viene de la sesión, no de la etiqueta de rol)", "UI que muestre secciones sin permiso verificado"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://nextjs.org/docs"]
verified_at: "2026-09-11"
```

El rol es una vista curada de enlaces, no un conjunto de permisos concedidos. La sesión efectiva
y la configuración controlan navegación y páginas; cada destino conserva su autorización.
REBUILD_VERIFIED / CONDITIONED cubre sólo el boundary demostrado en V381.

## 2. Applicability

Usar para navegación simplificada de una composición compatible con permisos ya definidos.
No usar como alta de personas, evaluación de capacitación, KPI financiero ni política comercial.

## 3. Architecture contract

- Owner de navegación: registry del BFF; roles filtra esa misma fuente, sin duplicar permisos.
- Las cuatro vistas intersectan enlaces con permisos efectivos. La etiqueta nunca concede acceso.
- role_workspace debe estar explícitamente true tras seleccionar el pack; ausente/false devuelve404.
- Sin sesión válida:307 a login local. Rol desconocido o similar a prototipo:404.
- Sesión sólo server-side; no se serializan identidad, tenant, organización ni token.
- Las instrucciones operativas siguen en /help versionado. Ninguna nueva escritura/tabla/dependencia.
- Reservar /dashboard, /guide y /guide/[role]; una colisión impide componer.
- Rollback de acceso por flag false, sin migración de estado. Detalles en docs/role-workspace.md.

## 4. Exact file manifest

```text
CREATE src/platform/roles/role-visibility.ts
CREATE src/platform/roles/role-visibility.test.ts
CREATE src/components/role-dashboard.tsx
CREATE src/app/dashboard/page.tsx
CREATE src/app/guide/page.tsx
CREATE src/app/guide/[role]/page.tsx
CREATE docs/role-workspace.md
```

## 5. Materialization blocks

### FILE: `src/platform/roles/role-visibility.ts`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:src/platform/roles/role-visibility.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "bab58d5fd0678aabbd503fdd5297ffe73643dbc3947dd0a1a14e322c837ed36d"
variables: []
secrets_allowed: false
```
````ts
import type { BusinessConfig } from "@/platform/config/schema";
import { enabledNavigation } from "@/platform/config/registry";
export type Section = { id: string; label: string; href: string };

// Call only after readSession succeeds. Common links are explicit, never wildcard grants.
export function visibleSections(permissions: readonly string[], config: BusinessConfig): Section[] {
  const ids: Record<string,string> = { dashboard: "overview", franchise: "sales", public_catalog: "catalog" };
  const sections = enabledNavigation(config, permissions).map(({id,label,href})=>({id:ids[id]??id,label,href}));
  if (config.features.role_workspace === true) sections.push({id:"guide",label:"Guía de uso",href:"/guide"});
  return sections;
}
export const ROLE_GUIDES = [
  {id:"owner",label:"Dueño",blurb:"Accesos de supervisión disponibles en tu sesión.",sections:["overview","admin","sales","factory","catalog","help"]},
  {id:"admin",label:"Administrador",blurb:"Accesos disponibles para la operación diaria.",sections:["overview","admin","sales","catalog","help"]},
  {id:"employee",label:"Empleado",blurb:"Accesos disponibles para atención y operación.",sections:["overview","sales","catalog","help"]},
  {id:"customer",label:"Cliente",blurb:"Accesos disponibles para tus consultas y entregas.",sections:["overview","customer","catalog","locations","help"]}
] as const;
export function roleGuide(id: string) { return ROLE_GUIDES.find(guide=>guide.id===id); }
export function sectionsForRole(id: string, permissions: readonly string[], config: BusinessConfig): Section[] {
  const guide=roleGuide(id);if(!guide)return [];
  return visibleSections(permissions,config).filter(section=>(guide.sections as readonly string[]).includes(section.id));
}
````

### FILE: `src/platform/roles/role-visibility.test.ts`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:src/platform/roles/role-visibility.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c4e085922de8de753148bfd6a99127952032a967469b7a9170fad78f8d177f14"
variables: []
secrets_allowed: false
```
````ts
import { beforeEach, describe, expect, it, test } from "vitest";
import { readFileSync } from "node:fs";
import { businessConfigSchema, type BusinessConfig } from "@/platform/config/schema";
import { enabledNavigation } from "@/platform/config/registry";
import { ROLE_GUIDES, roleGuide, sectionsForRole, visibleSections } from "./role-visibility";
let config:BusinessConfig;
beforeEach(()=>{config=businessConfigSchema.parse(JSON.parse(readFileSync("config/business.example.json","utf8")));Object.assign(config.features,{role_workspace:true,customer_portal:true,factory_portal:true,public_catalog:true});for(const module of ["catalog","crm","procurement"])config.modules[module]!.enabled=true;});
const common=["dashboard","help","locations","public_catalog"];
const matrix=[
  {
    "name": "none",
    "permissions": [],
    "privateIds": []
  },
  {
    "name": "customer",
    "permissions": [
      "customer:self"
    ],
    "privateIds": [
      "customer"
    ]
  },
  {
    "name": "factory",
    "permissions": [
      "factory:read"
    ],
    "privateIds": [
      "factory"
    ]
  },
  {
    "name": "admin",
    "permissions": [
      "admin:read"
    ],
    "privateIds": [
      "admin",
      "franchise"
    ]
  },
  {
    "name": "lead",
    "permissions": [
      "lead:read"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "inventory",
    "permissions": [
      "inventory:allocate"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "payment",
    "permissions": [
      "payment:create"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "handover",
    "permissions": [
      "handover:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "resource",
    "permissions": [
      "resource:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "availability-read",
    "permissions": [
      "availability:read"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "availability-manage",
    "permissions": [
      "availability:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "appointment",
    "permissions": [
      "appointment:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "quote-only",
    "permissions": [
      "quote:write"
    ],
    "privateIds": []
  },
  {
    "name": "wrong-case",
    "permissions": [
      "Admin:Read",
      "owner",
      "lead:reader",
      "* "
    ],
    "privateIds": []
  },
  {
    "name": "wildcard",
    "permissions": [
      "*"
    ],
    "privateIds": [
      "admin",
      "customer",
      "factory",
      "franchise"
    ]
  }
] as const;
test.each(matrix)("header and workspace use actual permissions: $name",row=>{
 const nav=enabledNavigation(config,row.permissions);expect(nav.map(s=>s.id).sort()).toEqual([...common,...row.privateIds].sort());
 const map:Record<string,string>={dashboard:"overview",franchise:"sales",public_catalog:"catalog"};
 expect(visibleSections(row.permissions,config).map(s=>s.id).sort()).toEqual([...common,...row.privateIds].map(id=>map[id]??id).concat("guide").sort());
});
describe("common authenticated access and exact grants",()=>{
 it("retains common help and guide for non-wildcard sessions",()=>{for(const permissions of [[],["customer:self"],["lead:read"]])expect(visibleSections(permissions,config).map(s=>s.id)).toEqual(expect.arrayContaining(["overview","help","guide"]));});
 it("preserves the customer login entry only for guests",()=>{expect(enabledNavigation(config,null).map(s=>s.id).sort()).toEqual([...common,"customer"].sort());expect(enabledNavigation(config,[]).some(s=>s.id==="customer")).toBe(false);});
 it("feature flags win over wildcard and the workspace stays opt-in",()=>{Object.assign(config.features,{role_workspace:false,customer_portal:false,factory_portal:false,public_catalog:false});const ids=visibleSections(["*"],config).map(s=>s.id);for(const id of ["overview","guide","customer","factory","catalog"])expect(ids).not.toContain(id);expect(ids).toContain("help");});
 it("disabled modules remove their destinations even for wildcard",()=>{for(const module of ["catalog","crm","procurement"])config.modules[module]!.enabled=false;expect(enabledNavigation(config,["*"]).map(s=>s.id)).not.toEqual(expect.arrayContaining(["factory","franchise","locations","public_catalog"]));for(const id of ["factory","franchise","locations","public_catalog"])expect(enabledNavigation(config,["*"]).map(s=>s.id)).not.toContain(id);});
 it("keeps the four labels without treating any label as a permission",()=>{expect(ROLE_GUIDES.map(g=>g.id)).toEqual(["owner","admin","employee","customer"]);for(const role of ROLE_GUIDES){expect(sectionsForRole(role.id,[role.id],config).map(s=>s.id)).not.toEqual(expect.arrayContaining(["sales","factory","admin","customer"]));}});
 it("every role view is a subset of the same effective-permission workspace",()=>{for(const row of matrix)for(const role of ROLE_GUIDES){const visible=new Set(visibleSections(row.permissions,config).map(s=>s.id));for(const section of sectionsForRole(role.id,row.permissions,config))expect(visible.has(section.id)).toBe(true);}});
 it("unknown and prototype role names never select a guide",()=>{for(const role of ["unknown","__proto__","constructor","OWNER","../owner"]){expect(roleGuide(role)).toBeUndefined();expect(sectionsForRole(role,["*"],config)).toEqual([]);}});
});
````

### FILE: `src/components/role-dashboard.tsx`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:src/components/role-dashboard.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "ad78a3223dda60dc0d8d58de291324b6f6e7c698479bbac717be77222a8511b0"
variables: []
secrets_allowed: false
```
````tsx
import Link from "next/link";
import type { Route } from "next";
import type { Section } from "@/platform/roles/role-visibility";
export function RoleDashboard({sections}:{sections:readonly Section[]}) {
  return <section className="card" aria-label="Accesos de tu sesión">
    <h2>Accesos disponibles</h2><p>Las opciones dependen de tus accesos y de las funciones habilitadas.</p>
    <nav aria-label="Secciones disponibles"><ul className="grid">{sections.map(section=><li key={section.id}><Link href={section.href as Route} className="button">{section.label}</Link></li>)}</ul></nav>
    {sections.length===0?<p>No hay accesos disponibles para esta vista.</p>:null}
  </section>;
}
````

### FILE: `src/app/dashboard/page.tsx`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:src/app/dashboard/page.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "ea68ffbe9d01db38fa86ef29c951f25c418f0f7cfb0d3c0fa13332a0bc063278"
variables: []
secrets_allowed: false
```
````tsx
import { notFound, redirect } from "next/navigation";
import type { Route } from "next";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { visibleSections } from "@/platform/roles/role-visibility";
import { RoleDashboard } from "@/components/role-dashboard";
export default async function DashboardPage() {
  const config=await loadBusinessConfig();if(config.features.role_workspace!==true)notFound();
  const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/dashboard" as Route);
  return <><div className="eyebrow">Panel</div><h1 className="pageTitle">Tu espacio de trabajo</h1>
    <p className="lede">Elegí un acceso para continuar. Consultá la ayuda de cada recorrido antes de actuar.</p>
    <RoleDashboard sections={visibleSections(session.permissions,config)}/></>;
}
````

### FILE: `src/app/guide/page.tsx`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:src/app/guide/page.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "fa9abe04dfbfc1669e1c1f71d74a0a415fd4f87741d43c2198ca4e1e03d13e75"
variables: []
secrets_allowed: false
```
````tsx
import Link from "next/link";
import type { Route } from "next";
import { notFound, redirect } from "next/navigation";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { ROLE_GUIDES } from "@/platform/roles/role-visibility";
export default async function GuidePage() {
  const config=await loadBusinessConfig();if(config.features.role_workspace!==true)notFound();
  const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/guide" as Route);
  return <><div className="eyebrow">Guía de uso</div><h1 className="pageTitle">Elegí un enfoque</h1>
    <p className="lede">Podés consultar el recorrido de un rol. Tus accesos no cambian.</p>
    <div className="grid">{ROLE_GUIDES.map(guide=><article className="card" key={guide.id}><h2>{guide.label}</h2><p>{guide.blurb}</p><Link href={`/guide/${guide.id}` as Route} className="button">Ver recorrido de {guide.label}</Link></article>)}</div>
    <p><Link href="/dashboard">Volver a mi panel</Link></p></>;
}
````

### FILE: `src/app/guide/[role]/page.tsx`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:src/app/guide/[role]/page.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "449ded016bc423081ca197821a36715158036ba89118d93f53225ab93ec10d65"
variables: []
secrets_allowed: false
```
````tsx
import Link from "next/link";
import type { Route } from "next";
import { notFound, redirect } from "next/navigation";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { roleGuide, sectionsForRole } from "@/platform/roles/role-visibility";
import { RoleDashboard } from "@/components/role-dashboard";
export default async function RoleGuidePage({params}:{params:Promise<{role:string}>}) {
  const config=await loadBusinessConfig();if(config.features.role_workspace!==true)notFound();
  const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/guide" as Route);
  const {role}=await params;const guide=roleGuide(role);if(!guide)notFound();
  return <><div className="eyebrow">Guía de uso</div><h1 className="pageTitle">Recorrido: {guide.label}</h1>
    <p className="lede">{guide.blurb} Elegir este enfoque no cambia tus accesos.</p>
    <RoleDashboard sections={sectionsForRole(role,session.permissions,config)}/>
    <p>Consultá la ayuda de los recorridos disponibles. Leer esta página no acredita capacitación ni habilita operaciones.</p>
    <p><Link href="/help">Abrir ayuda versionada</Link></p><p><Link href="/guide">Volver a la guía</Link></p></>;
}
````


### FILE: `docs/role-workspace.md`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:WORKSPACE:docs-role-workspace.md:v1"
operation: CREATE
provenance: AUTHORED
source: "Local permission-filtered workspace integration; see ROLE_WORKSPACE_INTEGRATION_V381.md"
license: "LicenseRef-Workspace-Owner"
sha256: "71df15da7d76b30ef02d93ca4ac39aa428c6cc5d64e2efd8ba73fa47b71bceae"
variables: []
secrets_allowed: false
```
````markdown
# Role workspace — V381

This pack supplies session-filtered navigation and four curated entry views, not a role grant,
financial dashboard, employee activation system or training curriculum. Selecting owner, admin,
employee or customer changes only the displayed subset of links. The existing verified session
and destination page/backend remain authoritative. Public catalog/help/locations stay public.

Compose TS-MULTIROLE-ONBOARDING0.1.8 with TS-GO-API-WEB-BRIDGE0.5.16,
TS-OIDC-PORTAL-ADAPTER0.2.x and TS-FRANCHISE-JOURNEY-PORTALS0.14.19. Reserve
/dashboard, /guide and /guide/[role] before composition. Then explicitly set
features.role_workspace=true in the consumer business configuration. The default absent/false
flag hides the menu and returns404 for all three routes, including wildcard sessions. Do not
enable it in a profile that does not select this pack. Existing module/feature flags still apply.

An empty authenticated permission set retains only common public sections and the guide. A guest
gets the public navigation and the customer login entry but no workspace. The private customer,
factory and admin links require customer:self, factory:read and admin:read respectively. Franchise
uses the actual destination's any-permission predicate: inventory:allocate, payment:create,
handover:manage, admin:read, lead:read, resource:manage, availability:read, availability:manage,
appointment:manage. Exact wildcard is supported; wrong case, quote:write alone and role labels
grant no private link. Features/modules constrain wildcard too. Hiding links is not an API ACL.

Operational instructions remain under the existing versioned /help owner. No unsupported payment
approval, royalty or high-value policy remains in the curated role labels. Reading a view does
not train a model, enroll an employee, change permissions or record learning progress.

Qualification:22 policy tests;236 full web tests pass plus one pre-existing connected skip;
typecheck/build;15 explicit permission profiles across five actual pages per profile in each of
four browsers (300 page states), with exact header/panel/view link sets. Guest/expired/tampered
sessions require the actual307/Location using the same cookie context with maxRedirects:0; no
IdP follows. Real encrypted cookies are synthetic fixture identities, not proof of external IdP.
Four disabled-flag tests return404. Unknown/prototype-like guide names return404, customer direct
/admin is denied, owner view cannot widen customer access, and links reach the existing help.
No browser writes, subject/tenant/token disclosure or mobile horizontal overflow observed.

Run test:workspace twice in the admitted Playwright gate: ELITE_WORKSPACE_E2E=enabled and
disabled, with owned HTTPS loopback production Next, ELITE_WEB_ROOT, synthetic AUTH_SESSION_SECRET
and the corresponding runtime config. Use the existing help-loopback-proxy fixture; never deploy
it. Retries are zero. Retain traces and stop both owned processes. No dependency lock changed.

Rollback: set features.role_workspace=false; routes become404 without state migration. When
removing the selected pack, leave the flag off and remove only its seven materialized files.
No database, provider, confidential document, external identity or training dataset is required
for this bounded reference qualification. Consumer authorization/IdP/revocation, business
onboarding, assessed progress, localization and target acceptance remain separate requirements.

V382 refreshes these compatible versions after private-read qualification. The V381 workspace implementation and its fifteen-profile policy remain byte-identical; see private-portal-reads.md for the new administrative/customer read and amount evidence.

V383 refreshes compatible selected versions for private-read recovery; the role source and permission policy remain unchanged.

V387 preserves role navigation and updates compatible portal/query/browser revisions for cursor traversal. No role or domain-policy change.

V388 aligns selected versions after extending cursor traversal to admin lists; no new role or grant.

V389 aligns selected versions after correcting customer appointment presentation to the existing business time zone. No role or grant changes.

V390 aligns selected versions for the existing customer appointment management link and durable cancellation recovery; roles and grants unchanged.

V393 aligns the composed BFF version after omission of its unused optional image optimizer dependency. Role implementation and permission policy remain unchanged.
````

## 6. Configuration surface

features.role_workspace=true sólo después de seleccionar este pack; ausente/false mantiene404.
Consume sesión JWE existente y flags/módulos del BFF; sin secretos nuevos. Ver docs/role-workspace.md.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Next.js | 16.3.4 | App Router, server components | MIT | runtime | https://github.com/vercel/next.js |
| React | 19.2.8 | componentes | MIT | runtime | https://github.com/facebook/react |
| Vitest | (del perfil web) | tests | MIT | build | https://github.com/vitest-dev/vitest |

## 8. Apply order

1. Componer BFF0.5.9, OIDC0.2.x, portales0.14.13 y los siete archivos de este pack.
2. Activar role_workspace explícitamente en configuración del consumer y verificar colisiones.
3. Ejecutar instalación frozen/offline,22tests de rol, suite web, typecheck/build y test:workspace
   con ambos estados del flag en los cuatro navegadores admitidos. Detalles en el documento incluido.
4. Rollback: flag false, después retirar sólo los siete archivos propios si corresponde.

## 9. Verification

V381:22tests de rol y236tests web PASS/1skip conectado previo; typecheck/build PASS.
Quince perfiles por cinco páginas ×cuatro navegadores=300estados con enlaces esperados exactos;
cuatro tests adicionales del flag false, sesión invitada/vencida/adulterada,404 y negación admin.
No se usa proveedor/IdP externo ni datos privados. No prueba activación de personal ni aprendizaje.
Los fallos originales y la corrección del oracle307 se preservan en el reporte V381.

## 10. Reconstruction evidence

Historia: reconstruction_evidence/TS_MULTIROLE_ONBOARDING_2026-09-02_V185.md.
Vigente: reconstruction_evidence/ROLE_WORKSPACE_INTEGRATION_V381.md y docs/role-workspace.md.

V382: exact selected-composition compatibility refreshed. Administrative lead reads are conditional on the existing lead grant; shared quote-format code corrects admin/customer minor-unit labels without repricing. TestAdministrativeReadBrowserPostgres uses real Go/JWKS/PostgreSQL and four browsers with synthetic identities; all production Go/SQL remains unchanged. The role implementation and quote acceptance/date/state logic retain exact source parity. See docs/private-portal-reads.md and reconstruction_evidence/PRIVATE_PORTAL_READS_V382.md. No full business, IdP, security, release or target admission.

V383: private read failures expose a fixed alert and explicit GET reconsultation in admin/customer/factory, without error details or command replay. Production build plus Go/JWKS/PostgreSQL and four browser projects prove persistent denial with the insufficient bearer and recovery after session correction; zero writes and unchanged durable snapshots.263web tests and canonical parity. Existing role/portal/Go business sources unchanged. See reconstruction_evidence/PRIVATE_READ_RECOVERY_V383.md; no full-control or target promotion.

V387: customer/factory cursor traversal connects three existing lists beyond25rows.81durable records per browser, four projects, reload/independent first-page recovery, zero writes and unchanged snapshots verified. Source is AUTHORED, business/API/SQL/dependencies unchanged. Supporting TEST02 progress, no whole-control or native security/release promotion. See reconstruction_evidence/PRIVATE_PORTAL_PAGINATION_V387.md.

V388: admin orders/cases/leads now traverse existing scoped keyset APIs using the shared navigator.276web tests and four browser projects exercise all six lists,163list records, no omissions/duplicates per list, permission removal, reload/first-page recovery and unchanged durable snapshots. All code remains AUTHORED/CONDITIONED; no full TEST02, native security or release promotion. See reconstruction_evidence/ADMIN_PORTAL_PAGINATION_V388.md.

V389: customer appointment times share trusted business locale/zone across SSR and interactive management, preserving original instants and cancellation behavior.282web tests plus eight connected browser/configuration runs, two appointment instants per view and seven-table unchanged snapshots. AUTHORED/CONDITIONED; no whole TEST02, security or release promotion. See reconstruction_evidence/CUSTOMER_APPOINTMENT_TIMEZONE_V389.md.

V390: customer cancellation receipt binding, synchronous single-submit and explicit GET recovery close the existing reference journey.294web tests;4actual cancellation browser projects,16unique durable cancellations/audits/outbox and28API writes;8read/time regression runs remain read-only. AUTHORED/CONDITIONED, no whole TEST02/security/release promotion. See reconstruction_evidence/CUSTOMER_CANCELLATION_RECOVERY_V390.md.

V393: the reference does not use runtime image transformation. Next image optimization is explicitly disabled and pnpm11.25.0 ignores the optional sharp dependency; frozen lock removes only its exclusive closure. This is dependency omission, not a repaired/admitted Sharp binary. Before enabling runtime image optimization, re-admit a compatible image pipeline, its exact dependencies, licenses/native provenance/security and web performance. Existing V386 vulnerability reproduction remains deferred. A connected real HTTP regression requires /_next/image to return404. No new version or business policy; AUTHORED/CONDITIONED. See reconstruction_evidence/OPTIONAL_IMAGE_DEPENDENCY_CONTAINMENT_V393.md.
