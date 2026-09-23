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
