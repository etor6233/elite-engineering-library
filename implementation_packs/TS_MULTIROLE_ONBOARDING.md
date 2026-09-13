# TypeScript Multi-Role Onboarding

## 1. Metadata

```yaml
pack_id: "TS-MULTIROLE-ONBOARDING"
pack_version: "0.2.1"
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
CREATE src/app/api/enterprise/metrics/route.test.ts
CREATE src/app/api/enterprise/metrics/route.ts
CREATE src/components/role-metrics.tsx
CREATE src/platform/metrics/contract.test.ts
CREATE src/platform/metrics/contract.ts
```

## 5. Materialization blocks

### FILE: `src/platform/roles/role-visibility.ts`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:src/platform/roles/role-visibility.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5a3c7d4f6048df2b2409da6880c67825a5a4c31d802f63e0853ae94e71b6ce49"
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
  {id:"owner",label:"Dueño",blurb:"Accesos de supervisión disponibles en tu sesión.",sections:["overview","admin","sales","factory","catalog","help","training","catalog_editor","supply","warranty","network"]},
  {id:"admin",label:"Administrador",blurb:"Accesos disponibles para la operación diaria.",sections:["overview","admin","sales","catalog","help","training","catalog_editor","supply","warranty","network"]},
  {id:"employee",label:"Empleado",blurb:"Accesos disponibles para atención y operación.",sections:["overview","sales","catalog","help","training","catalog_editor","supply","warranty"]},
  {id:"customer",label:"Cliente",blurb:"Accesos disponibles para tus consultas y entregas.",sections:["overview","customer","catalog","locations","help","training","warranty"]}
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
sha256: "9a238f3f77e84bc2ee0f05e9c9c3c8563171e6aed132ba8ae41cb975a319dce7"
variables: []
secrets_allowed: false
```
````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";
import Link from "next/link";
import type { Route } from "next";
import type { Section } from "@/platform/roles/role-visibility";
export function RoleDashboard({sections}:{sections:readonly Section[]}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  return <section className="card" aria-label={t("p0752")}>
    <h2>{t("p0753")}</h2><p>{t("p0754")}</p>
    <nav aria-label={t("p0755")}><ul className="grid">{sections.map(section=><li key={section.id}><Link href={section.href as Route} className="button">{controlled(section.label)}</Link></li>)}</ul></nav>
    {sections.length===0?<p>{t("p0756")}</p>:null}
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
sha256: "efbaa20406fc8b1db8beecfe4b57233fc5a21152d1fbdc34e5db53f44da0cec6"
variables: []
secrets_allowed: false
```
````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { notFound, redirect } from "next/navigation";
import type { Route } from "next";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { visibleSections } from "@/platform/roles/role-visibility";
import { RoleDashboard } from "@/components/role-dashboard";
import { RoleMetrics } from "@/components/role-metrics";
export default async function DashboardPage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const config=await loadBusinessConfig();if(config.features.role_workspace!==true)notFound();
  const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/dashboard" as Route);
  return <><div className="eyebrow">{t("p0068")}</div><h1 className="pageTitle">{t("p0069")}</h1>
    <p className="lede">{t("p0070")}</p>
    <RoleDashboard sections={visibleSections(session.permissions,config)}/><RoleMetrics locale={privateLocale.language} organizations={session.organizations} permissions={session.permissions}/></>;
}
````

### FILE: `src/app/guide/page.tsx`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:src/app/guide/page.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "420e3b9cd1f736cf26e6502de4bc0b71135e4fa3f8b88100ee6025957a6d3c08"
variables: []
secrets_allowed: false
```
````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import Link from "next/link";
import type { Route } from "next";
import { notFound, redirect } from "next/navigation";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { ROLE_GUIDES } from "@/platform/roles/role-visibility";
export default async function GuidePage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const config=await loadBusinessConfig();if(config.features.role_workspace!==true)notFound();
  const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/guide" as Route);
  return <><div className="eyebrow">{t("p0100")}</div><h1 className="pageTitle">{t("p0101")}</h1>
    <p className="lede">{t("p0102")}</p>
    <div className="grid">{ROLE_GUIDES.map(guide=><article className="card" key={guide.id}><h2>{controlled(guide.label)}</h2><p>{controlled(guide.blurb)}</p><Link href={`/guide/${guide.id}` as Route} className="button">{t("p0103")} {controlled(guide.label)}</Link></article>)}</div>
    <p><Link href="/dashboard">{t("p0104")}</Link></p></>;
}
````

### FILE: `src/app/guide/[role]/page.tsx`
```yaml
block_id: "TS-MULTIROLE-ONBOARDING:src/app/guide/[role]/page.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "024aa5a6796d6063ad360a82d438843f2e24deb71509ace2b61fd7a8bca1896f"
variables: []
secrets_allowed: false
```
````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import Link from "next/link";
import type { Route } from "next";
import { notFound, redirect } from "next/navigation";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { roleGuide, sectionsForRole } from "@/platform/roles/role-visibility";
import { RoleDashboard } from "@/components/role-dashboard";
export default async function RoleGuidePage({params}:{params:Promise<{role:string}>}) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const config=await loadBusinessConfig();if(config.features.role_workspace!==true)notFound();
  const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/guide" as Route);
  const {role}=await params;const guide=roleGuide(role);if(!guide)notFound();
  return <><div className="eyebrow">{t("p0100")}</div><h1 className="pageTitle">{t("p0110")} {controlled(guide.label)}</h1>
    <p className="lede">{controlled(guide.blurb)} {t("p0111")}</p>
    <RoleDashboard sections={sectionsForRole(role,session.permissions,config)}/>
    <p>{t("p0112")}</p>
    <p><Link href="/help">{t("p0113")}</Link></p><p><Link href="/guide">{t("p0114")}</Link></p></>;
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

V402 composed delta: T2804 connected training: existing versioned help/audit/shared approvals/outbox/BFF; opt-in host and navigation; bounded body retains original default. No new dependency or automatic grant. TRAINING_CONNECTED_RELEASE_V402.md.

V402 composed delta: T2804 catalog role source/edit/review/publication transport reuses original model/price writers and catalog owner. Bounded PNG and text defaults retained. No new dependency. CATALOG_ROLE_AUTHORING_RELEASE_V402.md.

V402 composed delta: T2804 connected supply role: original Operations/BindPlan composed atomically; bounded forms, latest review/requester projection, nav and GET recovery. No migration/dependency/state-machine change. SUPPLY_ROLE_RELEASE_V402.md.

V402 composed delta: T2804 warranty roles: same authorized transaction projects existing immutable diagnosis/plan/work/quality/acceptance; quote version read, bounded forms and durable GET recovery. No migration/dependency/domain-rule change. WARRANTY_ROLE_RELEASE_V402.md.

V402 composed delta: T2804 network role: original four fulfillment SQL bodies extracted unchanged into one transaction with immutable result; optional host, forms and GET recovery. Migration0077, no dependency/domain-rule change. NETWORK_ROLE_RELEASE_V402.md.

V402 composed delta: T2804 role metrics use existing domain read models with exact strings, organization/customer/factory/program permissions, NPS minimum/retention and no zero on unavailable. FAIL868 converted lead and FAIL457 bounded generic body corrected. ROLE_METRICS_RELEASE_V402.md/json; no new dependency or corporate authorship.

### FILE: `src/app/api/enterprise/metrics/route.test.ts`

```yaml
block_id: "TS-MULTIROLE-ONBOARDING-METRIC-DELTA:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6dc73a40072a19312fa9b670eb158ee16008728e7da5bf60bb0b217ef2b1d247"
variables: []
secrets_allowed: false
```

````typescript
import{it,expect,vi,afterEach}from"vitest";import{NextRequest}from"next/server";
const m=vi.hoisted(()=>({session:vi.fn(),get:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes("*")||s.permissions.includes(p)}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:async()=>({features:{role_workspace:true}})}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get}));
import{GET}from"./route";afterEach(()=>vi.clearAllMocks());const r=(q:string)=>new NextRequest("https://portal.example.test/api/enterprise/metrics?"+q);
it("denies customer aggregate, extra and duplicate selectors before backend",async()=>{m.session.mockResolvedValue({permissions:["customer:self"],organizations:["org"]});expect((await GET(r("kind=orders&organization_id=org"))).status).toBe(403);expect((await GET(r("kind=own-orders&organization_id=other"))).status).toBe(403);expect((await GET(r("kind=own-orders&organization_id=org&organization_id=other"))).status).toBe(400);expect((await GET(r("kind=own-orders&organization_id=org&actor=other"))).status).toBe(400);expect(m.get).not.toHaveBeenCalled()});
it("failure and mismatched backend scope produce unavailable, never an empty metric",async()=>{m.session.mockResolvedValue({permissions:["admin:read"],organizations:["org"]});m.get.mockRejectedValueOnce(new Error("down"));expect((await GET(r("kind=orders&organization_id=org"))).status).toBe(503);m.get.mockResolvedValue({kind:"orders",organization_id:"other",scope:"organization",source:"sales.customer_order",basis:"current_registered_records_by_state",observed_at:"2026-09-12T01:00:00Z",rows:[]});expect((await GET(r("kind=orders&organization_id=org"))).status).toBe(503)});
````

### FILE: `src/app/api/enterprise/metrics/route.ts`

```yaml
block_id: "TS-MULTIROLE-ONBOARDING-METRIC-DELTA:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b22dbcaab1cc52e3cf4bb5b16a04deaee31d304136a18a9a607d6861ef709907"
variables: []
secrets_allowed: false
```

````typescript
import{NextResponse,type NextRequest}from"next/server";
import{readSession,allowed}from"@/platform/auth/session";
import{loadBusinessConfig}from"@/platform/config/load";
import{protectedGet}from"@/platform/backend/protected-client";
import{metricID,metricFamilies,parseMetric}from"@/platform/metrics/contract";
const response=(value:unknown,status=200)=>NextResponse.json(value,{status,headers:{"cache-control":"private, no-store",vary:"Cookie"}});
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.role_workspace!==true)return response({code:"NOT_FOUND"},404);
 const s=await readSession();if(!s)return response({code:"UNAUTHENTICATED"},401);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),org=q.get("organization_id"),id=q.get("id")??"",family=metricFamilies.find(f=>f.id===kind);
 if(r.nextUrl.search.length>600||q.size!==(kind==="stored-value"||kind==="survey"?3:2)||["kind","organization_id"].some(k=>q.getAll(k).length!==1)||q.getAll("id").length>1||!family||!metricID.safeParse(org).success||(kind==="stored-value"||kind==="survey")&&!metricID.safeParse(id).success||Array.from(q.keys()).some(k=>!["kind","organization_id","id"].includes(k)))return response({code:"INVALID_QUERY"},400);
 if(!allowed(s,family.permission)||!s.organizations.includes(org!))return response({code:"FORBIDDEN"},403);
 try{
  const path=kind==="stored-value"?"/v1/franchise/stored-value/metrics":kind==="survey"?"/v1/reporting/surveys/"+encodeURIComponent(id):"/v1/reporting/operations/"+family.id;
  const value=await protectedGet(s,path,{organization_id:org!,...(kind==="stored-value"?{program_id:id}:{})});
  return response(parseMetric(family.id,value,org!,id));
 }catch{return response({code:"METRIC_UNAVAILABLE"},503)}
}
````

### FILE: `src/components/role-metrics.tsx`

```yaml
block_id: "TS-MULTIROLE-ONBOARDING-METRIC-DELTA:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "617fb23a826961de282f7ba6df5798ee80d978efe0fdda790589626608cc8e80"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import{useState,useRef}from"react";
import{visibleMetricFamilies,parseMetric,type Metric}from"@/platform/metrics/contract";
export function RoleMetrics({organizations,permissions,locale="es"}:{organizations:string[];permissions:string[];locale?:"es"|"en"}){
 const families=visibleMetricFamilies(permissions),[org,setOrg]=useState(organizations[0]??""),[kind,setKind]=useState<string>(families[0]?.id??""),[id,setID]=useState(""),[value,setValue]=useState<Metric|null>(null),[status,setStatus]=useState<"idle"|"loading"|"unavailable">("idle"),generation=useRef(0),en=locale==="en";
 const reset=()=>{generation.current++;setValue(null);setStatus("idle")};
 async function read(){const current=++generation.current;setValue(null);setStatus("loading");try{const query=new URLSearchParams({kind,organization_id:org,...(kind==="survey"||kind==="stored-value"?{id}:{})}),r=await fetch("/api/enterprise/metrics?"+query,{cache:"no-store",signal:AbortSignal.timeout(6500)});if(!r.ok)throw new Error();const v=parseMetric(kind,await r.json(),org,id);if(current===generation.current){setValue(v);setStatus("idle")}}catch{if(current===generation.current)setStatus("unavailable")}}
 if(!families.length)return <p>{en?"No metric permission is assigned.":"No tenés permisos de indicadores asignados."}</p>;
 return <section aria-label={en?"Operational metrics":"Indicadores operativos"}>
 <h2>{en?"Operational metrics":"Indicadores operativos"}</h2>
 <p>{en?"Current registered records. Each consultation has its own observation time. Order amounts are not collected revenue; currencies and program points remain separate.":"Registros actuales. Cada consulta tiene su propia fecha de observación. Los importes de pedidos no son ingresos cobrados; monedas y puntos de programas permanecen separados."}</p>
 <form onSubmit={e=>{e.preventDefault();void read()}} style={{display:"grid",gap:"1rem",maxWidth:"50rem"}}>
 <label>{en?"Organization":"Organización"}<select value={org} onChange={e=>{reset();setOrg(e.target.value)}}>{organizations.map(o=><option key={o} value={o}>{o}</option>)}</select></label>
 <label>{en?"Indicator":"Indicador"}<select value={kind} onChange={e=>{reset();setKind(e.target.value);setID("")}}>{families.map(f=><option key={f.id} value={f.id}>{f[locale]}</option>)}</select></label>
 {(kind==="stored-value"||kind==="survey")&&<label>{kind==="survey"?(en?"Survey identifier":"Identificador de encuesta"):(en?"Program identifier":"Identificador de programa")}<input required pattern="[A-Za-z0-9][A-Za-z0-9._:-]{0,127}" maxLength={128} value={id} onChange={e=>{reset();setID(e.target.value)}}/></label>}
 <button type="submit" disabled={status==="loading"||!org}>{status==="loading"?(en?"Reading…":"Consultando…"):(en?"Consult indicator":"Consultar indicador")}</button></form>
 <div aria-live="polite">{status==="unavailable"&&<p role="alert">{en?"Consultation unavailable. No zero or previous result is shown. Check access and service availability, then retry.":"Consulta no disponible. No se muestra cero ni un resultado anterior. Revisá acceso y disponibilidad del servicio y volvé a consultar."}</p>}
 {value&&<article><h3>{families.find(f=>f.id===kind)?.[locale]}</h3><p>{en?"Organization":"Organización"}: {value.organization_id} · {en?"Observed":"Observado"}: <time dateTime={value.observed_at}>{value.observed_at}</time></p>
 {"scope"in value&&<p>{en?"Scope":"Alcance"}: {value.scope}</p>}
 {"responses"in value?<><p>{en?"Responses":"Respuestas"}: {value.responses}</p><p>{value.available?"NPS: "+value.nps:(en?"NPS unavailable: configured minimum not met.":"NPS no disponible: no se alcanzó el mínimo configurado.")}</p></>:
 <>{value.rows.length===0?<p>{en?"No registered records for this scope.":"Sin registros para este alcance."}</p>:<div style={{overflowX:"auto"}}><table><caption>{en?"Current values by status or operation":"Valores actuales por estado u operación"}</caption><thead><tr><th>{en?"State / operation":"Estado / operación"}</th><th>{en?"Count":"Cantidad"}</th><th>{en?"Currency":"Moneda"}</th><th>{en?"Amount in minor units":"Importe en unidades menores"}</th>{"program_id"in value&&<th>{en?"Point delta":"Variación de puntos"}</th>}</tr></thead><tbody>{"program_id"in value?value.rows.map(r=><tr key={r.operation}><td>{r.operation}</td><td>{r.operations}</td><td>{value.currency}</td><td>{r.applied_minor_units}</td><td>{r.points_delta}</td></tr>):value.rows.map(r=><tr key={r.state+":"+r.currency}><td>{r.state}</td><td>{r.count}</td><td>{r.currency??"—"}</td><td>{r.total_minor_units??"—"}</td></tr>)}</tbody></table></div>}</>}
 {"program_id"in value&&<p>{en?"Program":"Programa"}: {value.program_id} · {value.program_kind} · {value.currency}</p>}
 <details><summary>{en?"Definition and data source":"Definición y fuente de datos"}</summary><p>{value.basis}</p><p>{value.source}</p>{"profile_sha256"in value&&<p style={{overflowWrap:"anywhere"}}>SHA-256: {value.profile_sha256}</p>}</details></article>}</div></section>
}
````

### FILE: `src/platform/metrics/contract.test.ts`

```yaml
block_id: "TS-MULTIROLE-ONBOARDING-METRIC-DELTA:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9a321f4dac8dbb9a55f03ffc24b4d168f04c3f301c6afe64188f51567d13f56c"
variables: []
secrets_allowed: false
```

````typescript
import{it,expect}from"vitest";import{parseMetric,visibleMetricFamilies}from"./contract";
const base={kind:"orders",organization_id:"org",scope:"organization",source:"sales.customer_order",basis:"current_registered_records_by_state",observed_at:"2026-09-12T01:00:00Z",rows:[{state:"confirmed",currency:"ARS",count:"2",total_minor_units:"18014398509481986"}]};
it("keeps exact values and rejects unsafe number and wrong scope",()=>{expect(parseMetric("orders",base,"org","")).toEqual(base);for(const v of[{...base,organization_id:"other"},{...base,kind:"stock"},{...base,rows:[{...base.rows[0],total_minor_units:18014398509481986}]}])expect(()=>parseMetric("orders",v,"org","")).toThrow()});
it("customer choices exclude aggregate organization metrics",()=>{expect(visibleMetricFamilies(["customer:self"]).map(f=>f.id)).toEqual(["own-orders","own-cases","own-appointments"]);expect(visibleMetricFamilies([])).toEqual([])});
it("unavailable survey never becomes NPS zero",()=>{const v={organization_id:"org",survey_id:"s",source:"crm.survey_definition + crm.survey_response",basis:"retained_survey_population_with_configured_minimum",responses:"1",available:false,nps:null,observed_at:base.observed_at};expect(parseMetric("survey",v,"org","s")).toEqual(v);expect(()=>parseMetric("survey",{...v,nps:0},"org","s")).toThrow()});
````

### FILE: `src/platform/metrics/contract.ts`

```yaml
block_id: "TS-MULTIROLE-ONBOARDING-METRIC-DELTA:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2fa9cec09e2daa2a372ff002fcdefecb523f37408f58b8b7042460980b8a5feb"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED presentation/binding glue over existing domain read models.
import {z} from "zod";
export const metricID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/);
const count=z.string().regex(/^(0|[1-9][0-9]{0,37})$/),signed=z.string().regex(/^-?(0|[1-9][0-9]{0,37})$/);
const point=z.string().regex(/^-?(0|[1-9][0-9]{0,37})(\.[0-9]{1,6})?$/),stamp=z.iso.datetime({offset:true});
export const metricFamilies=[
 {id:"orders",permission:"admin:read",es:"Pedidos de la organización",en:"Organization orders"},
 {id:"own-orders",permission:"customer:self",es:"Mis pedidos",en:"My orders"},
 {id:"leads",permission:"lead:read",es:"Leads por estado",en:"Leads by status"},
 {id:"stock",permission:"admin:read",es:"Stock por estado",en:"Stock by status"},
 {id:"cases",permission:"admin:read",es:"Casos de servicio",en:"Service cases"},
 {id:"own-cases",permission:"customer:self",es:"Mis casos de servicio",en:"My service cases"},
 {id:"shipments",permission:"admin:read",es:"Envíos con origen o destino local",en:"Local origin or destination shipments"},
 {id:"appointments",permission:"appointment:manage",es:"Turnos de la organización",en:"Organization appointments"},
 {id:"own-appointments",permission:"customer:self",es:"Mis turnos",en:"My appointments"},
 {id:"factory-destination",permission:"factory:read",es:"Producción para este destino",en:"Production for this destination"},
 {id:"factory-owned",permission:"supply:factory-read",es:"Producción de esta fábrica",en:"This factory's production"},
 {id:"supply",permission:"supply:read",es:"Compras seriadas del destino",en:"Destination serial purchases"},
 {id:"stored-value",permission:"stored_value:read",es:"Puntos y valor por programa",en:"Points and value by program"},
 {id:"survey",permission:"surveys:read",es:"NPS de una encuesta",en:"Survey NPS"}
] as const;
export function visibleMetricFamilies(permissions:readonly string[]){return metricFamilies.filter(f=>permissions.includes("*")||permissions.includes(f.permission))}
export const operationsMetric=z.object({kind:metricID,organization_id:metricID,scope:z.enum(["organization","customer","origin_or_destination","destination","factory"]),source:z.string().min(1).max(200),basis:z.literal("current_registered_records_by_state"),observed_at:stamp,rows:z.array(z.object({state:metricID,currency:z.string().regex(/^[A-Z]{3}$/).optional(),count,total_minor_units:signed.optional()}).strict()).max(100)}).strict();
export const storedMetric=z.object({organization_id:metricID,program_id:metricID,program_kind:z.enum(["gift_card","loyalty"]),currency:z.string().regex(/^[A-Z]{3}$/),profile_sha256:z.string().regex(/^[a-f0-9]{64}$/),source:z.literal("stored_value.operation + entry + account"),basis:z.literal("committed_ledger_by_operation"),observed_at:stamp,rows:z.array(z.object({operation:z.enum(["issue","accrue","redeem","reverse"]),operations:count,points_delta:point,applied_minor_units:signed}).strict()).max(4)}).strict();
export const surveyMetric=z.object({organization_id:metricID,survey_id:metricID,source:z.literal("crm.survey_definition + crm.survey_response"),basis:z.literal("retained_survey_population_with_configured_minimum"),responses:count,available:z.boolean(),nps:z.number().finite().min(-100).max(100).nullable(),observed_at:stamp}).strict().refine(v=>v.available===(v.nps!==null));
export type Metric= z.infer<typeof operationsMetric>|z.infer<typeof storedMetric>|z.infer<typeof surveyMetric>;
export function parseMetric(kind:string,value:unknown,org:string,id:string){
 const v=kind==="stored-value"?storedMetric.parse(value):kind==="survey"?surveyMetric.parse(value):operationsMetric.parse(value);
 if(v.organization_id!==org||("kind"in v&&v.kind!==kind)||("program_id"in v&&v.program_id!==id)||("survey_id"in v&&v.survey_id!==id))throw new Error("metric scope mismatch");
 return v;
}
````


V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.
