# Architecture Pack — Public Web, Customer and Admin Portals

> **Estado:** `CANDIDATE_PACK`.  
> **Fuentes:** React/TypeScript/Chromium authorities, Microsoft Fluent UI/Playwright, Oracle JET candidate, Saleor Dashboard candidate, Shopify Hydrogen integration reference.  
> **Regla:** ninguna librería de UI decide dominio, autorización o accesibilidad por sí sola.

## 1. Superficies separadas

```text
public-web       anonymous-first, SEO, catalog, lead conversion
customer-portal authenticated self-service and delegated organization access
admin-portal    operational and privileged workflows
shared          contracts, tokens, primitive UI, telemetry helpers, test fixtures
```

Desplegar juntas o separadas es una decisión operacional. Mantener bundles/rutas/config/secrets diferenciados es una frontera de seguridad.

## 2. Baseline técnico

- React-family framework con SSR/static rendering para público;
- TypeScript strict, generated API clients y schema validation en boundaries;
- design tokens + accessible primitives; branded assets con licencia propia;
- server-side session/BFF cuando evita tokens sensibles en browser;
- Playwright para journeys y axe/accessibility tooling complementario;
- Content Security Policy, Trusted Types cuando aplique y dependency controls;
- image pipeline, font strategy y performance budgets.

Next.js/React Router/Vite u otro framework sólo se eligen después de deployment, rendering, team y lifecycle. No fijar framework por popularidad.

## 3. Public web contract

Rutas mínimas:

```text
/{locale}/
/{locale}/models
/{locale}/models/{slug}
/{locale}/compare
/{locale}/branches/{slug}
/{locale}/test-drive
/{locale}/quote
/{locale}/contact
/{locale}/service-recall
/{locale}/legal/{document}
```

Requisitos:

- URL/metadata/canonical/hreflang coherentes;
- structured data generado desde catalog version, no CMS libre;
- consentimiento separado por purpose;
- formularios con server validation, honeypot/rate limit/risk controls;
- idempotency token emitido por server;
- success no revela existencia de persona/cuenta;
- analytics se carga según consentimiento/jurisdicción;
- content/media fallback y cache invalidation por catalog version.

## 4. Customer portal contract

- login OIDC code+PKCE; session cookie `HttpOnly`, `Secure`, `SameSite` adecuada;
- account/customer/organization context resuelto en server;
- navigation no es autorización;
- downloads mediante short-lived authorized URL o streaming controlado;
- sensitive changes requieren recent auth/step-up;
- optimistic UI sólo para acciones reversibles y con reconciliation;
- estados inciertos de pago/fulfillment se muestran como tales.

## 5. Admin portal contract

- default deny; cada screen deriva de server capabilities;
- list queries reciben un authorization scope, no filtran resultados en cliente;
- high-risk mutation muestra resource, consequence y reason;
- maker-checker para publish, grants, refunds o recall cuando el risk model lo exige;
- impersonation/support access con banner, scope, expiry y audit;
- bulk actions tienen preview, bounded batch, progress y partial failure report;
- export es job auditado con expiry, no respuesta síncrona masiva.

## 6. Estado de cliente

Clasificación:

| Estado | Ubicación preferida |
|---|---|
| URL/navigation/filter | router/query string |
| server source of truth | query cache con invalidación explícita |
| form draft | component/form state; persistence sólo si requerida |
| session/capabilities | server/session boundary; mínimo en browser |
| feature flag result | server-evaluated bootstrap |
| secrets/tokens | nunca localStorage; server/secure storage |

No crear un global store para todos los datos. Cache cliente no confirma stock, autorización ni pago.

## 7. API client

Generated client desde OpenAPI/contract fijado:

```text
deadline/cancellation
request and trace IDs
idempotency key for commands
typed problem details
retry only safe/idempotent operations
auth/session refresh single-flight
schema mismatch telemetry
```

El cliente no reintenta automáticamente una mutation con efecto externo sin key y policy explícita.

## 8. Design system

Niveles:

```text
tokens → primitives → patterns → domain components → page compositions
```

- tokens versionados: color, type, spacing, motion, elevation, breakpoints;
- contrast, focus visible, reduced motion y forced colors;
- primitives no conocen entidades de dominio;
- domain components declaran states loading/empty/error/partial/forbidden;
- copy crítico/localización forma parte del contract;
- icon/font licensing incluida en provenance.

Fluent UI, Oracle JET o componentes Saleor pueden donar patrones/componentes sólo después de auditoría de assets, bundle, accessibility y fit visual.

## 9. Performance budgets

Por template y dispositivo/mercado:

- JS initial/route budget;
- image bytes/dimensions/formats;
- font count/weight/subsetting;
- LCP/INP/CLS field metrics;
- server response/cache hit;
- third-party scripts y CPU impact;
- admin table virtualization/pagination thresholds.

CI usa budgets sintéticos; producción usa RUM consentido y segmentado. Un promedio no oculta p75/p95 ni mercados lentos.

## 10. Seguridad frontend

- output encoding y no raw HTML salvo sanitizer allowlist;
- CSP nonce/hash y report-only rollout;
- CSRF defense según session model;
- clickjacking/frame ancestors;
- upload validation server-side aunque haya preview cliente;
- no PII/secrets en URL, analytics, error payload o source maps públicos;
- dependency pinning, lockfile, integrity/provenance y source-map access control;
- admin served from distinct origin si reduce riesgo y operación lo soporta.

## 11. Playwright journey suite

```text
public: catalog → compare → consented lead → confirmation
customer: login → order → document → service appointment
branch: login → lead assignment → quote → reserve serial
HQ: draft model → approve → publish → public visibility
factory: confirm PO → ASN → HQ receiving
security: cross-tenant URL/API, expired session, revoked grant, CSRF, replay
degraded: payment pending, provider down, stale search, slow API
```

Cada test usa browser context fresco, stable test IDs sólo donde semántica/role no alcanza, API seeding autorizado y artifacts redacted.

## 12. Gates

- [ ] keyboard-only y screen reader critical paths;
- [ ] responsive widths y zoom 200/400% aplicable;
- [ ] Chromium/Firefox/WebKit critical journeys;
- [ ] locale, long strings, RTL si aplica, timezone/currency;
- [ ] CSP/headers/cookie audit;
- [ ] no cross-surface privileged bundle leakage;
- [ ] performance budgets y RUM plan;
- [ ] offline/degraded states honestos;
- [ ] Playwright traces/screenshots/videos sin PII real.
