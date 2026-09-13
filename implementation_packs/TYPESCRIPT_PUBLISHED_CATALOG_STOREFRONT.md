# Published catalog storefront

## 1. Metadata

```yaml
pack_id: "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Existing Next storefront renders the current approved tenant-bound catalog version, canonical detail pages and normalized media; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

Select with TS-GO-API-WEB-BRIDGE0.8.0. Disabled profile preserves existing storefront. Backend publication activation requires the separately selected Go owner.

## 3. Architecture contract

Typed Zod projection binds tenant/origin/model/variant identities; request-local React snapshot joins metadata and body. Existing LeadForm/Intl display reused. Current publication only: no draft fallback, fixed bounded PNG transport, source hash check, exact approved robots/sitemap paths.

## 4. Exact file manifest

```text
CREATE docs/CATALOG_STOREFRONT_REFERENCE.md
CREATE microsoft_playwright_browser_gate/tests/catalog-publication-connected.spec.mjs
CREATE src/app/api/public/catalog/media/[sha]/route.ts
CREATE src/app/models/[code]/page.tsx
CREATE src/platform/backend/catalog-client.ts
CREATE src/platform/catalog/load.ts
CREATE src/platform/catalog/schema.test.ts
CREATE src/platform/catalog/schema.ts
```

## 5. Materialization blocks

### FILE: `docs/CATALOG_STOREFRONT_REFERENCE.md`

```yaml
block_id: "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1205b9f5c1b2afb5eed0deba1db516df41e97458b8666716c9ba4cd278638ba8"
variables: []
secrets_allowed: false
```

````markdown
# Published catalog storefront

Select TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT together with the existing
TS-GO-API-WEB-BRIDGE. The complete franchise adds the connected Go owner.
No dependency/runtime/lock changes. All new source is AUTHORED binding; existing
Next/React/Zod/Intl/Playwright retain their admitted dependency provenance.

CATALOG_RELEASE_ENABLED=true reads the fixed backend /v1/public/catalog using
the existing validated ENTERPRISE_API_BASE_URL/TENANT_CODE/ORGANIZATION_CODE.
PUBLIC_SITE_ORIGIN must match the immutable publication HTTPS origin. Indexing
defaults disabled; PUBLIC_INDEXING_ENABLED=1 requires that explicit origin.
Runtime changes of a publication use no-store transport and request-local React
memoization, so metadata and body use one snapshot. There is no cross-request
catalog cache or mutable-draft fallback when the module is enabled.

The models index links to real /models/{code} pages only for current published
models. Details use captured name/variants/prices/media, the existing minor amount
presentation and LeadForm. Unsupported unsafe display amounts use the existing
explicit fallback without changing price math. Canonical, exact robots paths and
sitemap derive only from validated current codes; no private routes, fabricated
dates or translation claims. Unpublished models/media return404.

Media is fetched through the same-origin bounded PNG proxy. The proxy requires
current publication membership, checks decoded byte count and exact SHA256,
emits image/png/nosniff and no native Next image optimizer is enabled.

Reconstruction has no dependency installer. Reuse the admitted Node/modules
lock and build with the normal Next CLI. The existing local dependency junction
requires next build --webpack, with TypeScript validation active. A normal
in-root dependency install may use the default builder after its own gate.
No root widening, re-admission, network install or type-check bypass is needed.

The Go fixture TestCatalogReleaseStorefrontReference starts the actual Next host
and invokes the existing Playwright Chromium gate against the same PostgreSQL
publication journey. Explicit ELITE_CATALOG_STOREFRONT=1, absolute ELITE_WEB_ROOT
and ELITE_NODE_BIN, and an owned loopback PAYMENT_CONNECTED_DB_URL are required.
No live credentials. It verifies four generations, including rollback, against
their exact typed source, along with schema/SEO negative tests and Next build.
This scope supplies the public detail journey; role authoring remains T2804.
````

### FILE: `microsoft_playwright_browser_gate/tests/catalog-publication-connected.spec.mjs`

```yaml
block_id: "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3fabdb31afbc4050eb40cc094c68a00f953184df8a727b5efeeac717f4a3c56b"
variables: []
secrets_allowed: false
```

````text
// AUTHORED fixture assertions over the existing admitted Playwright gate.
import { test, expect } from '@playwright/test';
import { createHash } from 'node:crypto';
test('approved catalog publication is the actual public storefront', async ({ page, request }) => {
  if(process.env.ELITE_CATALOG_STOREFRONT!=='1')throw new Error('explicit catalog fixture required');
  const publication=JSON.parse(process.env.ELITE_CATALOG_EXPECTED);
  const model=publication.models[0], variant=publication.variants.find(v=>v.model_id===model.id);
  const browserErrors=[];page.on('pageerror',e=>browserErrors.push(e.message));
  await page.goto('/models');
  await expect(page.getByRole('link',{name:model.displayName,exact:true})).toHaveAttribute('href','/models/'+model.code);
  await page.getByRole('link',{name:model.displayName,exact:true}).click();
  await expect(page.getByRole('heading',{level:1,name:model.displayName,exact:true})).toBeVisible();
  await expect(page.locator('link[rel="canonical"]')).toHaveAttribute('href',model.canonical_url);
  await expect(page.locator('meta[name="robots"]')).toHaveAttribute('content','index, follow');
  const money=new Intl.NumberFormat('es-AR',{style:'currency',currency:publication.currency,currencyDisplay:'code'}).format(Number(variant.amount_minor_units)/100);
  await expect(page.getByText(money,{exact:true})).toBeVisible();
  const image=page.getByRole('img',{name:model.displayName,exact:true});
  await expect(image).toHaveAttribute('src','/api/public/catalog/media/'+model.media.sha256);
  await expect(image).toHaveJSProperty('naturalWidth',model.media.width);
  const png=await request.get('/api/public/catalog/media/'+model.media.sha256);
  expect(png.status()).toBe(200);expect(png.headers()['content-type']).toBe('image/png');
  expect(createHash('sha256').update(await png.body()).digest('hex')).toBe(model.media.sha256);
  const sitemap=await request.get('/sitemap.xml'),robots=await request.get('/robots.txt');
  expect(sitemap.status()).toBe(200);expect(await sitemap.text()).toContain('<loc>'+model.canonical_url+'</loc>');
  expect(await sitemap.text()).not.toContain('/admin');expect(await sitemap.text()).not.toContain('<lastmod>');
  expect(await robots.text()).toContain('Allow: /models/'+model.code+'$');
  expect(await robots.text()).toContain('Disallow: /');
  const absent=await request.get('/models/never-published');
  expect(absent.status()).toBe(404);
  expect((await request.get('/api/public/catalog/media/'+'0'.repeat(64))).status()).toBe(404);
  await expect(page.getByRole('button',{name:'Solicitar información',exact:true})).toBeVisible();
  expect(browserErrors).toEqual([]);
});
````

### FILE: `src/app/api/public/catalog/media/[sha]/route.ts`

```yaml
block_id: "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9c0973af671a73b77dee48c4844e4b649e971ca868f69f53359a1b4784c3d845"
variables: []
secrets_allowed: false
```

````typescript
import { createHash } from "node:crypto";
import { loadPublishedCatalog } from "@/platform/catalog/load";
import { readPublishedCatalogPNG } from "@/platform/backend/catalog-client";
export const dynamic = "force-dynamic";
export async function GET(_request: Request, { params }: { params: Promise<{ sha: string }> }) {
  const { sha } = await params;
  if (!/^[0-9a-f]{64}$/u.test(sha)) return new Response(null, { status: 404 });
  const catalog = await loadPublishedCatalog();
  if (!catalog?.models.some(model => model.media.sha256 === sha)) return new Response(null, { status: 404 });
  const response = await readPublishedCatalogPNG(sha);
  if (!response.ok || response.headers.get("content-type") !== "image/png" || !response.body) {
    if (response.body) void response.body.cancel().catch(() => {});
    return new Response(null, { status: response.status === 404 ? 404 : 502 });
  }
  const reader = response.body.getReader(), parts: Uint8Array[] = [];
  let size = 0, done = false;
  try {
    while (true) {
      const item = await reader.read(); if (item.done) { done = true; break; }
      size += item.value.byteLength; if (size > 4_194_304) return new Response(null, { status: 502 });
      parts.push(item.value);
    }
    const body = Buffer.concat(parts);
    if (createHash("sha256").update(body).digest("hex") !== sha) return new Response(null, { status: 502 });
    return new Response(new Uint8Array(body), { headers: { "content-type": "image/png", "x-content-type-options": "nosniff", "cache-control": "public, max-age=0, must-revalidate", etag: `"${sha}"` } });
  } finally {
    if (!done) void reader.cancel().catch(() => {});
    reader.releaseLock();
  }
}
````

### FILE: `src/app/models/[code]/page.tsx`

```yaml
block_id: "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ed092b2ae01d9a1cd193f1f6f4b7e527f9f68f4f69d6d58e2a9c90ae12edfe29"
variables: []
secrets_allowed: false
```

````tsx
import { notFound } from "next/navigation";
import { loadPublishedCatalog } from "@/platform/catalog/load";
import { publicModelMetadata } from "@/platform/seo/public-indexing";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { minorAmountPresentation } from "@/platform/i18n/money";
import { LeadForm } from "@/app/connected/lead-form";

export const dynamic = "force-dynamic";
async function selected(code: string) {
  const catalog = await loadPublishedCatalog();
  const model = catalog?.models.find(item => item.code === code);
  if (!catalog || !model) notFound();
  return { catalog, model };
}
export async function generateMetadata({ params }: { params: Promise<{ code: string }> }) {
  const { model } = await selected((await params).code);
  return publicModelMetadata(model.code, model.displayName, model.canonical_url);
}
export default async function ModelPage({ params }: { params: Promise<{ code: string }> }) {
  const [{ catalog, model }, locale] = await Promise.all([selected((await params).code), loadPublicLocale()]);
  return <section lang={locale.locale}>
    <h1 className="pageTitle">{model.displayName}</h1>
    {/* The same-origin bounded PNG route uses no native image optimizer. */}
    <img src={"/api/public/catalog/media/" + model.media.sha256} alt={model.displayName} width={model.media.width} height={model.media.height} style={{ maxWidth: "100%", height: "auto" }} />
    <div className="grid">{catalog.variants.filter(item => item.model_id === model.id).map(variant => {
      const amount = BigInt(variant.amount_minor_units);
      const display = minorAmountPresentation(amount <= BigInt(Number.MAX_SAFE_INTEGER) ? Number(amount) : NaN, catalog.currency);
      return <article className="card" key={variant.id}><h2>{variant.display_name}</h2><p>{display.amountLabel}</p></article>;
    })}</div>
    <LeadForm modelId={model.id} locale={locale.locale} />
  </section>;
}
````

### FILE: `src/platform/backend/catalog-client.ts`

```yaml
block_id: "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1136cd61981c41b8b7d7357f6fbb15db1a984047e38caa2adbb61177f67224cd"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import { publishedCatalogTransportSettings, readBackendResponse, BackendProblem } from "./public-client";
export async function readPublishedCatalogSource(): Promise<{ value: unknown; tenant: string } | null> {
  if (process.env.CATALOG_RELEASE_ENABLED !== "true") return null;
  const { base, tenant } = publishedCatalogTransportSettings();
  const response = await fetch(new URL("/v1/public/catalog", base), { cache: "no-store", redirect: "error", signal: AbortSignal.timeout(5000), headers: { accept: "application/json" } });
  if (response.status === 404) { if (response.body) void response.body.cancel().catch(() => {}); return null; }
  return { value: await readBackendResponse<unknown>(response), tenant };
}
export async function readPublishedCatalogPNG(sha: string): Promise<Response> {
  if (process.env.CATALOG_RELEASE_ENABLED !== "true" || !/^[0-9a-f]{64}$/u.test(sha)) throw new BackendProblem(404, "CATALOG_NOT_FOUND");
  const { base } = publishedCatalogTransportSettings();
  return fetch(new URL("/v1/public/catalog/media/" + sha, base), { cache: "no-store", redirect: "error", signal: AbortSignal.timeout(5000), headers: { accept: "image/png" } });
}
````

### FILE: `src/platform/catalog/load.ts`

```yaml
block_id: "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2e898b43e866df542905af3c5362571d72fcabf3750a4ad0ce09e8ce27a6d42e"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import { cache } from "react";
import { readPublishedCatalogSource } from "@/platform/backend/catalog-client";
import { readPublicIndexing } from "@/platform/seo/public-indexing";
import { parsePublishedCatalog } from "./schema";

// Request-local React memoization keeps metadata and page on the same snapshot.
export const loadPublishedCatalog = cache(async () => {
  const source = await readPublishedCatalogSource();
  return source === null ? null : parsePublishedCatalog(source.value, source.tenant, readPublicIndexing().origin);
});
````

### FILE: `src/platform/catalog/schema.test.ts`

```yaml
block_id: "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a6b421b1f5949d6b02f43b9f2c67a9cf15edf65e5749f6b93fc58749c1471bb8"
variables: []
secrets_allowed: false
```

````typescript
import { afterEach, describe, expect, it, vi } from "vitest";
import { parsePublishedCatalog } from "./schema";
import { publicModelMetadata, publicRobots, publicSitemap } from "../seo/public-indexing";
const origin = "https://catalog.example.invalid", hash = "a".repeat(64);
function publication() {
  return { tenant_code: "catalog-j3", generation: "2", source_sha256: hash,
    effective_price_book_id: "effective-2", market: "AR", currency: "ARS",
    models: [{ id: "one", code: "bicycle-one", displayName: "Approved", vehicleClass: "bicycle",
      specification: {}, media: { id: "png", original_sha256: hash, sha256: hash, width: 2, height: 2 },
      canonical_url: origin + "/models/bicycle-one" }],
    variants: [{ id: "v-one", model_id: "one", code: "v-one", display_name: "Approved variant",
      battery_specification: {}, homologation_state: "unknown", amount_minor_units: "120001", tax_mode: "inclusive" }] };
}
afterEach(() => vi.unstubAllEnvs());
describe("approved publication storefront boundary", () => {
  it("preserves quoted int64 money and binds the configured tenant and canonical origin", () => {
    const p = publication(); p.variants[0]!.amount_minor_units = "9223372036854775807";
    expect(parsePublishedCatalog(p, "catalog-j3", origin)).toEqual(p);
  });
  it("rejects another tenant, origin, draft-only field and a variant outside the snapshot", () => {
    expect(() => parsePublishedCatalog(publication(), "other-tenant", origin)).toThrow("tenant binding");
    expect(() => parsePublishedCatalog(publication(), "catalog-j3", "https://other.invalid")).toThrow("canonical binding");
    expect(() => parsePublishedCatalog({ ...publication(), profile: {} }, "catalog-j3", origin)).toThrow();
    const p = publication(); p.variants[0]!.model_id = "absent";
    expect(() => parsePublishedCatalog(p, "catalog-j3", origin)).toThrow("variant binding");
  });
  it.each(["/admin", "/models/other", "/models/bicycle-one?secret=1", "/models/bicycle-one#private"])("rejects mismatched canonical path %s", path => {
    const p = publication(); p.models[0]!.canonical_url = origin + path;
    expect(() => parsePublishedCatalog(p, "catalog-j3", origin)).toThrow();
  });
  it.each(["-1", "01", "1.1", "9223372036854775808"])("rejects money outside the source int64 contract %s", value => {
    const p = publication(); p.variants[0]!.amount_minor_units = value;
    expect(() => parsePublishedCatalog(p, "catalog-j3", origin)).toThrow();
  });
  it("exposes only approved real model paths and keeps the disabled profile out of crawling", () => {
    const settings = { enabled: true, origin };
    expect(publicSitemap(settings, ["bicycle-one"])).toEqual([{url: origin + "/"}, {url: origin + "/models"}, {url: origin + "/models/bicycle-one"}]);
    expect(publicRobots(settings, ["bicycle-one"]).rules).toEqual({userAgent:"*",disallow:"/",allow:["/$","/models$","/models/bicycle-one$","/_next/static/","/icon.svg$"]});
    expect(publicSitemap({ ...settings, enabled: false }, ["bicycle-one"])).toEqual([]);
    expect(() => publicSitemap(settings, ["../../admin"])).toThrow();
    vi.stubEnv("PUBLIC_SITE_ORIGIN", origin); vi.stubEnv("PUBLIC_INDEXING_ENABLED", "1");
    expect(publicModelMetadata("bicycle-one", "Approved", origin + "/models/bicycle-one")).toEqual({
      title: "Approved", robots: {index:true,follow:true}, alternates:{canonical:origin+"/models/bicycle-one"} });
    expect(() => publicModelMetadata("bicycle-one", "Approved", origin + "/admin")).toThrow();
  });
});
````

### FILE: `src/platform/catalog/schema.ts`

```yaml
block_id: "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f445cd5c15746ae9c5677b97ff9aee6a5a9211ee5ee4eed0a58f9f979d95a25e"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED typed projection validation; no local price or business rules.
import { z } from "zod";
const id = z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$/u);
const code = z.string().regex(/^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/u);
const digest = z.string().regex(/^[0-9a-f]{64}$/u);
const unsigned = z.string().regex(/^(0|[1-9][0-9]{0,18})$/u).refine(value => BigInt(value) <= 9223372036854775807n);
export const publishedCatalogSchema = z.object({
  tenant_code: code, generation: unsigned.refine(value => BigInt(value) > 0n),
  source_sha256: digest, effective_price_book_id: id,
  market: z.string().regex(/^[A-Z]{2}$/u), currency: z.string().regex(/^[A-Z]{3}$/u),
  models: z.array(z.object({
    id, code, displayName: z.string().min(2).max(160), vehicleClass: z.enum(["motorcycle", "bicycle", "scooter", "utility", "other"]),
    specification: z.record(z.string(), z.unknown()),
    media: z.object({ id, original_sha256: digest, sha256: digest, width: z.number().int().min(1).max(2048), height: z.number().int().min(1).max(2048) }).strict(),
    canonical_url: z.string().url()
  }).strict()).min(1).max(32),
  variants: z.array(z.object({
    id, model_id: id, code: z.string().min(1).max(128), display_name: z.string().min(1).max(160),
    battery_specification: z.record(z.string(), z.unknown()),
    homologation_state: z.enum(["unknown", "pending", "approved", "rejected", "expired"]),
    amount_minor_units: unsigned, tax_mode: z.enum(["inclusive", "exclusive", "not-applicable"])
  }).strict()).min(1).max(128)
}).strict();
export type PublishedCatalog = z.infer<typeof publishedCatalogSchema>;
export function parsePublishedCatalog(value: unknown, tenant: string, origin: string | null): PublishedCatalog {
  const result = publishedCatalogSchema.parse(value);
  if (result.tenant_code !== tenant) throw new Error("Catalog tenant binding mismatch");
  const ids = new Set<string>(), codes = new Set<string>(), variants = new Set<string>();
  for (const model of result.models) {
    if (ids.has(model.id) || codes.has(model.code) || model.media.width * model.media.height > 1_048_576) throw new Error("Invalid catalog model identity");
    ids.add(model.id); codes.add(model.code);
    const uri = new URL(model.canonical_url);
    if (uri.protocol !== "https:" || uri.username || uri.password || uri.search || uri.hash ||
        uri.pathname !== "/models/" + model.code || (origin !== null && uri.origin !== origin)) throw new Error("Catalog canonical binding mismatch");
  }
  for (const variant of result.variants) {
    if (!ids.has(variant.model_id) || variants.has(variant.id)) throw new Error("Invalid catalog variant binding");
    variants.add(variant.id);
  }
  return result;
}
````

## 6. Configuration surface

docs/CATALOG_STOREFRONT_REFERENCE.md covers runtime flags/origin, existing backend configuration, disabled indexing, exact CLI build and connected browser fixture.

## 7. Dependency bill

No new upstream or dependency. AUTHORED configuration, SQL/interface/orchestration/test glue around existing admitted owners. Go1.26.8 PNG/HTTP, PG18.6 and existing source adapters retain exact provenance.

## 8. Apply order

Materialize with MARKDOWN-COMPOSITOR0.3.0 and complete closure into an absent target. Apply migrations73/74 in order. Populated downgrade refuses history; empty74/73down and73/74up proven. Narrow original APIs preserved.

## 9. Verification

Actual PG71migrations,3snapshots/12reviews/5publications,1new7replay,13table rollback, source isolation, expiry/receipt recovery, current search/media, HTTP lost-response GET, feed rejection/unknown/GET reconciliation. Host12base+feed guard/hash/key/token gates; finite2s fuzz121059execs. Next/Chromium same four publication versions. Provider SDK mappings T2805 remain explicit.

## 10. Reconstruction evidence

CATALOG_CONNECTED_RELEASE_V402.md/json binds source, failed and successful receipts, exact affected profile reconstruction and notices. No global readiness or production claim from this one journey.


V402 composed delta: Use observed canonical pack_id TS-GO-API-WEB-BRIDGE; no change to tested product source.
