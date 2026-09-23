# Public indexing

This BFF uses the exact installed Next metadata/robots/sitemap APIs. It publishes
only the existing home and model-list pages; model detail URLs do not exist in
this composition. Private portals inherit noindex,nofollow. Robots directives
are crawler advice and do not replace authentication, remove indexed pages or
guarantee indexing, ranking or search-engine support.

The approved pages and their public Next static assets/icon remain crawlable;
otherwise a crawler could see the HTML but be unable to render CSS/JavaScript.

Indexing is disabled by default: no sitemap URLs, no WebSite JSON-LD, noindex on
pages and Disallow: / in robots.txt. The operator must explicitly set both
PUBLIC_INDEXING_ENABLED=1 and PUBLIC_SITE_ORIGIN to the public HTTPS origin for
the approved deployment, for example https://www.example.com. A malformed
origin or flag is a configuration error; errors never echo its value.

PUBLIC_SITE_ORIGIN is independent of the identity provider APP_BASE_URL.
Only an origin is accepted: no credentials, path prefix, non-default port,
query, fragment, whitespace or percent encoding. Request Host and forwarded
headers never define canonical URLs. Runtime routes read the current setting;
do not bake a staging host into a production sitemap. Deployment must confirm
domain ownership, TLS, redirects, canonical origin, cache behavior and target
crawler acceptance before enabling indexing.

Sitemap entries do not invent publication/modification dates, translated routes,
offers, prices or review ratings. Homepage JSON-LD contains configured business
name and locale only, is escaped for HTML and uses the existing CSP nonce.
This is a single configured site, not a domain resolver for multiple tenants.
Full localization, CMS publishing and per-model SEO require their own routes,
sources and integration evidence before extending the exact public allowlist.
