# TypeScript Enterprise Web Golden Path

V321 containment: this historical, unselected candidate still pins pnpm11.19.0,
whose published bundle includes brace-expansion5.0.8/GHSA-rgw5-rvv9-x895.
Do not execute or admit that tool for a new target. Preserve the old code/lock
for forensics; the tested pnpm11.25.0 correction belongs to the three active
BFF/Playwright/Lighthouse owners. This legacy pack requires its own dependency
and functional qualification before reuse; V321 does not transfer their PASS.

## 1. Metadata

```yaml
pack_id: "TS-ENTERPRISE-WEB"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CANDIDATE
claim: "Prototipo full-stack TypeScript histórico y reconstruible con catálogo/leads, route handlers, DB/outbox y shells de portales; no es el frontend/BFF recomendado ni un backend por defecto."
stacks: ["Node.js 24", "TypeScript 7", "Next.js 16", "React 19", "PostgreSQL/PGlite", "openid-client 6.8.5"]
compatible_with: ["PBC-CORE concepts", "PG-TX-FOUNDATION concepts"]
incompatible_with: ["adopción automática como frontend sobre el backend Go", "uso como backend por defecto", "edge runtimes sin Node.js APIs", "uso productivo de PGlite como reemplazo no evaluado de PostgreSQL gestionado"]
license_expression: "LicenseRef-Workspace-Owner AND third-party licenses in THIRD_PARTY_NOTICES.md"
upstream_sources: ["https://nextjs.org", "https://react.dev", "https://www.postgresql.org", "https://pglite.dev"]
verified_at: "2026-08-25"
```

Este pack preserva una implementación funcional creada en el workspace, pero combina UI, route handlers, DB y outbox TypeScript. Por eso queda como candidato histórico reconstruible y se excluye del perfil web empresarial vigente. “Golden Path” es sólo el nombre heredado: no es la foundation, el backend por defecto, el frontend/BFF recomendado ni el sistema empresarial completo.

El linaje de la primera extracción se conserva únicamente en la evidencia histórica V4. `golden_starter/` ya no existe ni es una dependencia: cada bloque declara este implementation pack Markdown como fuente canónica y autosuficiente.

## 2. Applicability

- Adoptar sólo si el blueprint elige deliberadamente un producto full-stack Node/React y admite su ownership de DB/outbox; no se selecciona automáticamente.
- Para frontend/BFF sobre el backend Go usar `TS-GO-API-WEB-BRIDGE 0.2.x` más `TS-OIDC-PORTAL-ADAPTER`.
- Rechazar cuando el blueprint exija runtime distinto, mobile nativo, hard real-time o límites que Next/Node no satisfagan con evidencia.
- PGlite es fallback local; producción exige PostgreSQL externo, secretos gestionados, identidad real y gates de producción.
- El claim cubre el código enumerado y sus tests; no cubre todos los módulos empresariales futuros.

## 3. Architecture contract

```text
Next App Router: public web + admin/customer/factory views + route handlers
  → typed platform services: config/auth/workflow/http/observability
  → modular domain services: catalog + CRM lead intake
  → PostgreSQL adapter or local PGlite fallback
  → transactional audit/outbox
```

Trust boundaries, invariants, degraded modes, production gaps and readiness are encoded in the included source, tests, `README.md`, `CODEX_START_HERE.md` and `SYSTEM_READINESS.md`. Authorization defaults deny. Unknown integration modes do not silently activate. Invalid business configuration stops bootstrap.

## 4. Exact file manifest

```text
CREATE .env.example
CREATE .gitignore
CREATE .github/workflows/verify.yml
CREATE CODEX_START_HERE.md
CREATE config/business.example.json
CREATE config/readiness.json
CREATE next-env.d.ts
CREATE next.config.ts
CREATE package.json
CREATE pnpm-lock.yaml
CREATE pnpm-workspace.yaml
CREATE README.md
CREATE scripts/check-readiness.ts
CREATE scripts/dispatch-outbox.ts
CREATE scripts/migrate.ts
CREATE scripts/seed.ts
CREATE scripts/verify-config.ts
CREATE src/app/admin/page.tsx
CREATE src/app/api/health/route.ts
CREATE src/app/api/leads/route.ts
CREATE src/app/api/platform/route.ts
CREATE src/app/customer/page.tsx
CREATE src/app/factory/page.tsx
CREATE src/app/globals.css
CREATE src/app/layout.tsx
CREATE src/app/models/[slug]/lead-form.tsx
CREATE src/app/models/[slug]/page.tsx
CREATE src/app/models/page.tsx
CREATE src/app/page.tsx
CREATE src/instrumentation.ts
CREATE src/modules/catalog/repository.ts
CREATE src/modules/crm/lead-service.test.ts
CREATE src/modules/crm/lead-service.ts
CREATE src/modules/enterprise/operations.test.ts
CREATE src/modules/enterprise/operations.ts
CREATE src/platform/auth/authorization.test.ts
CREATE src/platform/auth/authorization.ts
CREATE src/platform/auth/oidc.test.ts
CREATE src/platform/auth/oidc.ts
CREATE src/platform/bootstrap.ts
CREATE src/platform/config/load.ts
CREATE src/platform/config/registry.ts
CREATE src/platform/config/runtime-policy.ts
CREATE src/platform/config/schema.test.ts
CREATE src/platform/config/schema.ts
CREATE src/platform/db/client.ts
CREATE src/platform/db/migrations.ts
CREATE src/platform/db/types.ts
CREATE src/platform/http/problem.ts
CREATE src/platform/http/problem.test.ts
CREATE src/platform/integrations/contract.test.ts
CREATE src/platform/integrations/contract.ts
CREATE src/platform/integrations/http-adapter.test.ts
CREATE src/platform/integrations/http-adapter.ts
CREATE src/platform/integrations/registry.ts
CREATE src/platform/observability/events.ts
CREATE src/platform/outbox/dispatcher.ts
CREATE src/platform/seed.ts
CREATE src/platform/workflow/engine.test.ts
CREATE src/platform/workflow/engine.ts
CREATE SYSTEM_READINESS.md
CREATE THIRD_PARTY_NOTICES.md
CREATE tsconfig.json
CREATE vitest.config.ts
```

## 5. Materialization blocks

### FILE: `.env.example`

```yaml
block_id: "TS-ENTERPRISE-WEB:env-example:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "88a2012dcfda29a521cf5b9fd4402b52ac84753241cd9f639ecc77b5d96ecf78"
variables: []
secrets_allowed: false
```

````dotenv
# Local development uses embedded PostgreSQL-compatible PGlite.
DATABASE_URL=file://pglite
BUSINESS_CONFIG_FILE=business.example.json

# Demo identity is strictly local. Production refuses it.
ALLOW_DEMO_IDENTITY=true
APP_ENV=development

# Production examples:
# APP_ENV=production
# DATABASE_URL=postgresql://user:password@host:5432/database
# OIDC_ISSUER=https://identity.example.com/realms/business
# OIDC_CLIENT_ID=enterprise-platform
# OIDC_CLIENT_SECRET_REF=secret-manager://oidc/client-secret
# OIDC_AUDIENCE=enterprise-platform
# OIDC_JWKS_URI=https://identity.example.com/realms/business/protocol/openid-connect/certs
APP_BASE_URL=https://app.example.com
AUTH_SESSION_SECRET=replace-with-at-least-32-random-characters
OIDC_ISSUER=https://identity.example.com
OIDC_CLIENT_ID=enterprise-web
OIDC_CLIENT_SECRET=inject-from-secret-manager
OIDC_AUDIENCE=enterprise-api
OIDC_JWKS_URI=https://identity.example.com/.well-known/jwks.json
ENTERPRISE_API_BASE_URL=https://api.example.com
ENTERPRISE_TENANT_CODE=example-tenant
ENTERPRISE_ORGANIZATION_CODE=example-store
````

### FILE: `.gitignore`

```yaml
block_id: "TS-ENTERPRISE-WEB:gitignore:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "1cf6cf78f1a64047af43a0372ccae8ebb70abf5c99d80b37925c8f0aadad6286"
variables: []
secrets_allowed: false
```

````gitignore
node_modules/
.next/
.data/
coverage/
.env
.env.local
*.log
*.tsbuildinfo
````

### FILE: `.github/workflows/verify.yml`

```yaml
block_id: "TS-ENTERPRISE-WEB:github-workflows-verify-yml:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "5e75271088f74e7c14ddc9d9d8b3479bab8f85403e862320e9b6f337b3c08f3b"
variables: []
secrets_allowed: false
```

````yaml
name: verify

on:
  pull_request:
  push:
    branches: [main]

permissions:
  contents: read

concurrency:
  group: verify-${{ github.ref }}
  cancel-in-progress: true

jobs:
  verify:
    runs-on: ubuntu-24.04
    timeout-minutes: 20
    steps:
      - name: Checkout exact revision
        uses: actions/checkout@08eba0b27e820071cde6df949e0beb9ba4906955 # v4.3.0
        with:
          persist-credentials: false

      - name: Install pnpm
        uses: pnpm/action-setup@9fd676a19091d4595eefd76e4bd31c97133911f1 # v4.2.0
        with:
          version: 11.19.0
          run_install: false

      - name: Install Node.js
        uses: actions/setup-node@2028fbc5c25fe9cf00d9f06a71cc4710d4507903 # v6.0.0
        with:
          node-version: 24.14.1
          cache: pnpm

      - name: Install locked dependencies
        run: pnpm install --frozen-lockfile

      - name: Verify source and production build
        run: pnpm verify

      - name: Audit production dependency graph
        run: pnpm audit --prod
````

### FILE: `CODEX_START_HERE.md`

```yaml
block_id: "TS-ENTERPRISE-WEB:codex-start-here-md:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "b7e79105b1c4fd1ab7fd8b0f39db160d38e681b887504be49e42a84f213bef83"
variables: []
secrets_allowed: false
```

````markdown
# Codex — inicio obligatorio desde el sistema único

Usa los archivos materializados por este pack como adapter web ejecutable; no copies candidatos públicos directamente.

## Instrucción inicial reutilizable

```text
Trabaja sobre el workspace objetivo y trata este pack únicamente como adapter web opcional. El blueprint y los contratos stack-neutral gobiernan el sistema.

Primero lee el AGENTS.md raíz, AGENT_SYSTEM_START.md, CODEX_ELITE_PROJECT_BOOTSTRAP.md, README.md y esta instrucción. Ejecuta el intake y crea PROJECT_AUTHORITY_MAP.md más el manifest del Engineering Execution Kit. No amplíes código hasta identificar journeys, invariantes, jurisdicción, failure costs, SLO, RPO/RTO e integraciones reales.

Convierte las variantes normales del negocio en config/business.<proyecto>.json: mercados, organizaciones, roles, permisos, módulos, workflows, campos, features e integraciones. No bifurques el core por cliente. Si una regla no puede expresarse sin perder corrección, crea una extensión detrás de un contrato explícito y registra un ADR.

Antes de habilitar una integración exige licencia/terms compatibles, adapter, secret reference, idempotencia, límites, reconciliación, sandbox y contract tests. Nunca actives identidad demo en producción.

Conserva el vertical slice existente y expándelo journey por journey. Cada entrega debe pasar pnpm verify y los gates aplicables del Engineering Execution Kit. No declares production-ready, seguro, escalable o de baja latencia sin evidencia reproducible.
```

## Regla de evolución

```text
configuración cuando cambia política o composición
→ módulo cuando agrega capacidad de dominio
→ adapter cuando cambia un proveedor externo
→ nuevo servicio sólo ante evidencia de aislamiento, escala, SLO u ownership
```

Codex debe informar en cada hito qué quedó `READY`, `CONDITIONED` o `OPEN`. La compilación sola nunca convierte el sistema en producción.
````

### FILE: `config/business.example.json`

```yaml
block_id: "TS-ENTERPRISE-WEB:config-business-example-json:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "9e73ad1267c7c9e6207b05fc4f0f8d5bde0c2737684a51324cf2f1c351b0d43f"
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
      "locales": ["es-AR"],
      "timeZone": "America/Argentina/Buenos_Aires",
      "taxMode": "external"
    }
  ],
  "organizationTypes": [
    {"id": "hq", "label": "Casa central", "allowedParents": []},
    {"id": "franchise", "label": "Franquicia", "allowedParents": ["hq"]},
    {"id": "branch", "label": "Sucursal", "allowedParents": ["franchise", "hq"]},
    {"id": "factory", "label": "Fábrica", "allowedParents": ["hq"]},
    {"id": "supplier", "label": "Proveedor", "allowedParents": ["hq"]}
  ],
  "roles": [
    {"id": "hq_admin", "label": "Administrador central", "permissions": ["*"]},
    {"id": "branch_manager", "label": "Responsable de sucursal", "permissions": ["catalog:read", "lead:read", "lead:assign", "inventory:read", "inventory:reserve", "order:read", "order:transition"]},
    {"id": "sales", "label": "Ventas", "permissions": ["catalog:read", "lead:read", "lead:update", "inventory:read", "order:create", "order:read"]},
    {"id": "factory_operator", "label": "Operador de fábrica", "permissions": ["procurement:read", "procurement:transition", "shipment:create"]},
    {"id": "customer", "label": "Cliente", "permissions": ["catalog:read", "own:order:read", "own:service:read"]}
  ],
  "modules": {
    "catalog": {"enabled": true},
    "crm": {"enabled": true},
    "procurement": {"enabled": true},
    "inventory": {"enabled": true},
    "orders": {"enabled": true},
    "payments": {"enabled": true},
    "fulfillment": {"enabled": true},
    "service": {"enabled": true},
    "documents": {"enabled": true},
    "integrations": {"enabled": true}
  },
  "workflows": {
    "lead": {
      "initial": "new",
      "states": ["new", "assigned", "qualified", "converted", "lost"],
      "transitions": [
        {"from": "new", "to": "assigned", "permission": "lead:assign"},
        {"from": "assigned", "to": "qualified", "permission": "lead:update"},
        {"from": "qualified", "to": "converted", "permission": "lead:update"},
        {"from": "assigned", "to": "lost", "permission": "lead:update"},
        {"from": "qualified", "to": "lost", "permission": "lead:update"}
      ]
    },
    "order": {
      "initial": "draft",
      "states": ["draft", "placed", "confirmed", "paid", "allocated", "delivered", "cancelled"],
      "transitions": [
        {"from": "draft", "to": "placed", "permission": "order:create"},
        {"from": "placed", "to": "confirmed", "permission": "order:transition"},
        {"from": "confirmed", "to": "paid", "permission": "payment:reconcile"},
        {"from": "paid", "to": "allocated", "permission": "inventory:reserve"},
        {"from": "allocated", "to": "delivered", "permission": "order:transition"},
        {"from": "draft", "to": "cancelled", "permission": "order:transition"},
        {"from": "placed", "to": "cancelled", "permission": "order:transition"}
      ]
    }
  },
  "customFields": {
    "lead": [
      {"id": "preferred_vehicle_use", "label": "Uso principal", "type": "select", "required": false, "options": ["urban", "delivery", "recreation", "fleet"]}
    ],
    "catalog_model": [
      {"id": "estimated_range_km", "label": "Autonomía estimada (km)", "type": "number", "required": true}
    ]
  },
  "integrations": [
    {"id": "mercado_pago", "provider": "mercado_pago", "enabled": false, "mode": "sandbox", "capabilities": ["payments"], "credentialRefEnv": "MERCADO_PAGO_CREDENTIAL_REF"},
    {"id": "amazon_sp_api", "provider": "amazon_sp_api", "enabled": false, "mode": "sandbox", "capabilities": ["catalog", "orders", "fulfillment"], "credentialRefEnv": "AMAZON_SP_API_CREDENTIAL_REF"},
    {"id": "mercado_libre", "provider": "mercado_libre", "enabled": false, "mode": "sandbox", "capabilities": ["catalog", "orders", "fulfillment"], "credentialRefEnv": "MERCADO_LIBRE_CREDENTIAL_REF"},
    {"id": "google_ads", "provider": "google_ads", "enabled": false, "mode": "sandbox", "capabilities": ["ads", "conversions"], "credentialRefEnv": "GOOGLE_ADS_CREDENTIAL_REF"},
    {"id": "meta_ads", "provider": "meta_ads", "enabled": false, "mode": "sandbox", "capabilities": ["ads", "conversions"], "credentialRefEnv": "META_ADS_CREDENTIAL_REF"}
  ],
  "features": {
    "public_catalog": true,
    "lead_capture": true,
    "customer_portal": true,
    "factory_portal": true,
    "vehicle_telemetry": false
  }
}
````

### FILE: `config/readiness.json`

```yaml
block_id: "TS-ENTERPRISE-WEB:config-readiness-json:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "89d930915eb580114ebc2d216a15e426e87e4d92d165aaf2559fd4a5c721b890"
variables: []
secrets_allowed: false
```

````json
{
  "schemaVersion": "1.0.0",
  "claims": {
    "configuration_contract": "READY",
    "local_database": "READY",
    "public_catalog_lead_slice": "READY",
    "authorization_kernel": "READY",
    "workflow_kernel": "READY",
    "production_postgres": "CONDITIONED",
    "production_identity": "OPEN",
    "operational_domain_journeys": "OPEN",
    "external_integrations": "OPEN",
    "observability_slo_alerting": "OPEN",
    "performance_capacity": "OPEN",
    "backup_restore_recovery": "OPEN",
    "deployment_rollout_rollback": "OPEN",
    "source_code_license": "OPEN"
  },
  "targets": {
    "baseline": [
      "configuration_contract",
      "local_database",
      "public_catalog_lead_slice",
      "authorization_kernel",
      "workflow_kernel"
    ],
    "production": [
      "configuration_contract",
      "public_catalog_lead_slice",
      "authorization_kernel",
      "workflow_kernel",
      "production_postgres",
      "production_identity",
      "operational_domain_journeys",
      "external_integrations",
      "observability_slo_alerting",
      "performance_capacity",
      "backup_restore_recovery",
      "deployment_rollout_rollback",
      "source_code_license"
    ]
  }
}
````

### FILE: `next-env.d.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:next-env-d-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "1862ac4bbbc5192d4bf562161df66ea547ed3e67173100656ab606ae9797db2b"
variables: []
secrets_allowed: false
```

````typescript
/// <reference types="next" />
/// <reference types="next/image-types/global" />
import "./.next/types/routes.d.ts";
import "./.next/types/root-params.d.ts";

// NOTE: This file should not be edited
// see https://nextjs.org/docs/app/api-reference/config/typescript for more information.
````

### FILE: `next.config.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:next-config-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "d56644b2a9cebf8c8e58120abc1805596d681c08c714286971ae562f4d5e6c6a"
variables: []
secrets_allowed: false
```

````typescript
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  poweredByHeader: false,
  typedRoutes: true,
  // Native/WASM database drivers must retain Node's module realm on the server.
  serverExternalPackages: ["@electric-sql/pglite", "pg"],
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
          { key: "Permissions-Policy", value: "camera=(), microphone=(), geolocation=()" },
          { key: "X-Frame-Options", value: "DENY" },
          { key: "Cross-Origin-Opener-Policy", value: "same-origin" },
          { key: "Content-Security-Policy", value: "base-uri 'self'; object-src 'none'; frame-ancestors 'none'; form-action 'self'; upgrade-insecure-requests" }
        ]
      }
    ];
  }
};

export default nextConfig;
````

### FILE: `package.json`

```yaml
block_id: "TS-ENTERPRISE-WEB:package-json:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "79eb0dc6d31106aaa0d426390262d1240e41635548d2e098182966dfd188b68d"
variables: []
secrets_allowed: false
```

````json
{
  "name": "elite-configurable-enterprise-platform",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "packageManager": "pnpm@11.19.0",
  "engines": {
    "node": ">=24.0.0"
  },
  "scripts": {
    "dev": "next dev",
    "dev:ready": "pnpm bootstrap && next dev",
    "build": "next build",
    "start": "next start",
    "typecheck": "tsc --noEmit",
    "test": "vitest run",
    "test:watch": "vitest",
    "config:check": "tsx scripts/verify-config.ts",
    "readiness:baseline": "tsx scripts/check-readiness.ts baseline",
    "readiness:production": "tsx scripts/check-readiness.ts production",
    "db:migrate": "tsx scripts/migrate.ts",
    "db:seed": "tsx scripts/seed.ts",
    "outbox:once": "tsx scripts/dispatch-outbox.ts",
    "licenses:report": "pnpm licenses list --prod --json",
    "bootstrap": "pnpm config:check && pnpm db:migrate && pnpm db:seed",
    "verify": "pnpm config:check && pnpm readiness:baseline && pnpm typecheck && pnpm test && pnpm build"
  },
  "dependencies": {
    "@electric-sql/pglite": "0.5.7",
    "jose": "6.2.10",
    "next": "16.3.2",
    "openid-client": "6.8.5",
    "pg": "8.23.0",
    "react": "19.2.8",
    "react-dom": "19.2.8",
    "server-only": "0.0.1",
    "zod": "4.4.3"
  },
  "devDependencies": {
    "@types/node": "26.2.0",
    "@types/pg": "8.23.1",
    "@types/react": "19.2.18",
    "@types/react-dom": "19.2.5",
    "tsx": "4.23.12",
    "typescript": "7.0.2",
    "vitest": "4.1.11"
  }
}
````

### FILE: `pnpm-lock.yaml`

```yaml
block_id: "TS-ENTERPRISE-WEB:pnpm-lock-yaml:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "048413eae77d24fdeb89982894238203d538d60358d3b6abe6b1373e0afb2e0f"
variables: []
secrets_allowed: false
```

````yaml
lockfileVersion: '9.0'

settings:
  autoInstallPeers: true
  excludeLinksFromLockfile: false

importers:

  .:
    dependencies:
      '@electric-sql/pglite':
        specifier: 0.5.7
        version: 0.5.7
      jose:
        specifier: 6.2.10
        version: 6.2.10
      next:
        specifier: 16.3.2
        version: 16.3.2(@types/node@26.2.0)(react-dom@19.2.8(react@19.2.8))(react@19.2.8)
      openid-client:
        specifier: 6.8.5
        version: 6.8.5
      pg:
        specifier: 8.23.0
        version: 8.23.0
      react:
        specifier: 19.2.8
        version: 19.2.8
      react-dom:
        specifier: 19.2.8
        version: 19.2.8(react@19.2.8)
      server-only:
        specifier: 0.0.1
        version: 0.0.1
      zod:
        specifier: 4.4.3
        version: 4.4.3
    devDependencies:
      '@types/node':
        specifier: 26.2.0
        version: 26.2.0
      '@types/pg':
        specifier: 8.23.1
        version: 8.23.1
      '@types/react':
        specifier: 19.2.18
        version: 19.2.18
      '@types/react-dom':
        specifier: 19.2.5
        version: 19.2.5(@types/react@19.2.18)
      tsx:
        specifier: 4.23.12
        version: 4.23.12
      typescript:
        specifier: 7.0.2
        version: 7.0.2
      vitest:
        specifier: 4.1.11
        version: 4.1.11(@types/node@26.2.0)(vite@8.2.2(@types/node@26.2.0)(esbuild@0.28.2)(tsx@4.23.12))

packages:

  '@electric-sql/pglite@0.5.7':
    resolution: {integrity: sha512-6WBks0NFO6YUhGP03mCNkbJuVFWwrZPrzZ/vsnAPnsCD28pybdZ+yoTrlLlWFoZst3kdaju5pxIslV/jrr36nQ==}

  '@emnapi/runtime@1.11.3':
    resolution: {integrity: sha512-Xz4Tpyki7XyrpbUK1jR1AhdAdaXyhhY4lZ3neLodmhpuWfy2PAQN5B46sAiU4liOXGLkHypn/qU+jvfWSCYYLA==}

  '@esbuild/aix-ppc64@0.28.2':
    resolution: {integrity: sha512-XExcO+dvLKvVtNTibSTBej1NCAbaGhWn9Ww1ZPx80qsahhPFe/8jgWP0IchNe0F3HwkU7n8ejhH8bjonqht8mQ==}
    engines: {node: '>=18'}
    cpu: [ppc64]
    os: [aix]

  '@esbuild/android-arm64@0.28.2':
    resolution: {integrity: sha512-5YfKeeI8qWfBZIX+u2xZC3Zlb3Os/gLS2sbEKM+I4ZOcsWmHS2WLysCcQZDAFRslDUU5Oiq44gf6PYN1vGwG5A==}
    engines: {node: '>=18'}
    cpu: [arm64]
    os: [android]

  '@esbuild/android-arm@0.28.2':
    resolution: {integrity: sha512-kXXoiPVVGQcnIYGOeaovwOURpniDBpSq4A03qkQ+BMQqtGG6HYap3xne9C1O1yo4TR3qxlCX5IqqmX6fFo2Lqg==}
    engines: {node: '>=18'}
    cpu: [arm]
    os: [android]

  '@esbuild/android-x64@0.28.2':
    resolution: {integrity: sha512-O387ite7SzUyCcy3JQX4P4bLtEA7bLLkx+esve5JHnyYfNTxcVpXZo9jhdB0lTKN44gztELTdU7nS8Nr16Fs1Q==}
    engines: {node: '>=18'}
    cpu: [x64]
    os: [android]

  '@esbuild/darwin-arm64@0.28.2':
    resolution: {integrity: sha512-n4KqkOQrraxHJcgjM1RvwbigfQKIKJVpM7xp+KsxiyUSrRdIXnt73VhrPAx0fV44hgfmIVKjxMN9J1t5jySVkw==}
    engines: {node: '>=18'}
    cpu: [arm64]
    os: [darwin]

  '@esbuild/darwin-x64@0.28.2':
    resolution: {integrity: sha512-uq6suIWYP37qzGddBKPw5QEQPi6HiLGsO7UmkpfyaYNQ3D+rN6w6WfwH+nuqcGXWvawGwxOEroO4YGnFh95azw==}
    engines: {node: '>=18'}
    cpu: [x64]
    os: [darwin]

  '@esbuild/freebsd-arm64@0.28.2':
    resolution: {integrity: sha512-n+I0BTSRIoy+d6RPKnEVwql5UwBJolytvY4mAOIEJorKlqgPII8ix6slVVrfZ5Tnj7glIZvloylbB/EJPMWEXw==}
    engines: {node: '>=18'}
    cpu: [arm64]
    os: [freebsd]

  '@esbuild/freebsd-x64@0.28.2':
    resolution: {integrity: sha512-78XJTJkvPs0kz2w61301PJjXl4g7q3JqiYMZ/M/yVI73EHBrCRTgkhu9oqG7vPqq+a/yadEW8aD+agKlk5xrmg==}
    engines: {node: '>=18'}
    cpu: [x64]
    os: [freebsd]

  '@esbuild/linux-arm64@0.28.2':
    resolution: {integrity: sha512-pW4AC0P3it8c7do9MVM4p51FzHzdM/TZrerurgRcHJ2WTa1VQ1CIq18xncfpBJw4ojkiZZrKW2yIBWBP92j6Ug==}
    engines: {node: '>=18'}
    cpu: [arm64]
    os: [linux]

  '@esbuild/linux-arm@0.28.2':
    resolution: {integrity: sha512-XlDnu2q5yoqems+xay6wSAcg9DDD7K9RLKZEBOMZm3ckNpJBvOX20tSfby8KfrrhINDyv9V2YVZKY/SpoGJI8w==}
    engines: {node: '>=18'}
    cpu: [arm]
    os: [linux]

  '@esbuild/linux-ia32@0.28.2':
    resolution: {integrity: sha512-CYbnj78HsIeA+DhgUKgFCfvNsTHFhMMrinUrMZpDXJXKN8T3XViTZ/+wtHeVxEWY8ewSzTFN+nRmSwO2tZaLUQ==}
    engines: {node: '>=18'}
    cpu: [ia32]
    os: [linux]

  '@esbuild/linux-loong64@0.28.2':
    resolution: {integrity: sha512-buwkd8nsph4R+ajRvw0qM5Hja/TXQow3ptzWO2EbG/cqcIkHloRrdlBtQlshyYGTNFvfkfJ5tpPLVkY4DtsPfQ==}
    engines: {node: '>=18'}
    cpu: [loong64]
    os: [linux]

  '@esbuild/linux-mips64el@0.28.2':
    resolution: {integrity: sha512-ZVykbDyk7519VwiNb9Lcj9m8XM6v5V9uKPvrEMkkEedVewf+0itkhahp4HDpgERXhwLRpWFypsGbG/J8s0QjJA==}
    engines: {node: '>=18'}
    cpu: [mips64el]
    os: [linux]

  '@esbuild/linux-ppc64@0.28.2':
    resolution: {integrity: sha512-CAXl+Dtd9UUuJd8pKKdwh6MLm3MUMiqMPmhZ3tTSXPqfyQ3vDl6R5hZdZ/kYojK4ofXtdfSv1tFq8XzWx3heNQ==}
    engines: {node: '>=18'}
    cpu: [ppc64]
    os: [linux]

  '@esbuild/linux-riscv64@0.28.2':
    resolution: {integrity: sha512-GeXCej4IQtU1B+QlDV8W/RRvbzI3O/Stss+/bCXv4lZls5WGRtu2a+3JkA3i4qIUlMXpcHebWpF8AkJhATowuA==}
    engines: {node: '>=18'}
    cpu: [riscv64]
    os: [linux]

  '@esbuild/linux-s390x@0.28.2':
    resolution: {integrity: sha512-3H1weTYZPxt/WOhByszQZybS9w5lKzUn1FDMsgEChbHWQwHYQQRfBxgCcZvPhjHfKyJjIievvMmEUawJrdY9Dg==}
    engines: {node: '>=18'}
    cpu: [s390x]
    os: [linux]

  '@esbuild/linux-x64@0.28.2':
    resolution: {integrity: sha512-4xTZr1FUmSoQW4XIWmit3tzQrUTZM+N3P0XV8xROKYF50XfI7xeO90+1bZvNwxIufQ9hDQVRJH5YhgPVF8A/HQ==}
    engines: {node: '>=18'}
    cpu: [x64]
    os: [linux]

  '@esbuild/netbsd-arm64@0.28.2':
    resolution: {integrity: sha512-sSATRjPeDBg3pdgHoQfoYBob11Kk1FGa9lui5RIHZCoCkJa9QKlvl3/vKz2usCmYYjs7ymJR/2Nnsqe+Hjt5nw==}
    engines: {node: '>=18'}
    cpu: [arm64]
    os: [netbsd]

  '@esbuild/netbsd-x64@0.28.2':
    resolution: {integrity: sha512-lqnzCV+mM0gIADaKihiCg6ifgfU2L3h5E33rNQBN1Y4MaVGnzryzmvvf7UHxprpQdE8hpqLolJ9Rl+SkIRDpyw==}
    engines: {node: '>=18'}
    cpu: [x64]
    os: [netbsd]

  '@esbuild/openbsd-arm64@0.28.2':
    resolution: {integrity: sha512-AL2qJILH7lNjrDmCQDvdxMfAUIv8KMNZOvrwAQ8i8//ntL9FflhOyMJ8OZSMBb8/AWXe3/5v5S20y3zCoZWKoQ==}
    engines: {node: '>=18'}
    cpu: [arm64]
    os: [openbsd]

  '@esbuild/openbsd-x64@0.28.2':
    resolution: {integrity: sha512-QtiuPytchRyC4rwUKhexJdQKvDuZ6hWloi3igqPQNUJCS1/v9EiO3UTOXR6A3FoMo4fnAKbWJdqaIwhOzh8qEw==}
    engines: {node: '>=18'}
    cpu: [x64]
    os: [openbsd]

  '@esbuild/openharmony-arm64@0.28.2':
    resolution: {integrity: sha512-WkhYDmpTjLvGlScA1rwjRUmhl4k8oXR3cIbtqWmELgU/dFeHHlEllxDvdWcNJV9rbzCexB5vz8gtNewWLgCT7Q==}
    engines: {node: '>=18'}
    cpu: [arm64]
    os: [openharmony]

  '@esbuild/sunos-x64@0.28.2':
    resolution: {integrity: sha512-GPMSkTOtMnv2U2F8gxe4Io6qmVs+YKyp832Etqqxr0hFngmXQ3rzwytelm3GIn7T4VviRUlf3sOgBOiTdvaf7g==}
    engines: {node: '>=18'}
    cpu: [x64]
    os: [sunos]

  '@esbuild/win32-arm64@0.28.2':
    resolution: {integrity: sha512-PIhhEkE9uPBleRBrQEJpUn7MBnibZzbGzYWPmY3x+YoVg/95zbjB4CxPPOQ8l5tYYM4mMaCthF8/1DIfBQQyWQ==}
    engines: {node: '>=18'}
    cpu: [arm64]
    os: [win32]

  '@esbuild/win32-ia32@0.28.2':
    resolution: {integrity: sha512-YmJbfTlvU7Sdn9BB+4PRES4oB6pxgS37MAONj+hBr/cpXS1aBPKXxNnDbu+QCWPj0o9dgyxeq79g6c5P8KeuYA==}
    engines: {node: '>=18'}
    cpu: [ia32]
    os: [win32]

  '@esbuild/win32-x64@0.28.2':
    resolution: {integrity: sha512-5ebpxr3nWMzrL/rnUI755Jkuee0bHL/Gq0WTF9lvcpv73wAp5eu8MfBUgWK9bhWvZjj7yX8etf/8tI8Ney695g==}
    engines: {node: '>=18'}
    cpu: [x64]
    os: [win32]

  '@img/colour@1.1.0':
    resolution: {integrity: sha512-Td76q7j57o/tLVdgS746cYARfSyxk8iEfRxewL9h4OMzYhbW4TAcppl0mT4eyqXddh6L/jwoM75mo7ixa/pCeQ==}
    engines: {node: '>=18'}

  '@img/sharp-darwin-arm64@0.35.3':
    resolution: {integrity: sha512-RMnFX7YQsMoh7lWfcM4NEHHymBX/rLuKNPVM84XE9ONPcaSCDgE7CHIHpSgPcO2xcRthgBy1HfNO319mwhIAkg==}
    engines: {node: '>=20.9.0'}
    cpu: [arm64]
    os: [darwin]

  '@img/sharp-darwin-x64@0.35.3':
    resolution: {integrity: sha512-Xo+5uFBtLN0BKqieTxiFzFPQAUlBbbH5iBKyRX/z1JrbnYsHTfKJnUfL8+p2TPXr1pXqao4eeL4Rl144uDpK9w==}
    engines: {node: '>=20.9.0'}
    cpu: [x64]
    os: [darwin]

  '@img/sharp-freebsd-wasm32@0.35.3':
    resolution: {integrity: sha512-lUxcqWIj2wMQ9BrwNjngcr1gWUr5xgaGThBRqPPalIC2n67Cqj1uPh8NnA/ZhAg8hUbKl+kVHKwgUIwe6ZYPrg==}
    engines: {node: '>=20.9.0'}
    os: [freebsd]

  '@img/sharp-libvips-darwin-arm64@1.3.2':
    resolution: {integrity: sha512-9J6ypZFpQBj4YnePGoq/S38w6nz+vqg5WZLrLGY4YuSemdMq47GMLBPO42MzwdGwpg/agZ7xzZcFHa48xlywfg==}
    cpu: [arm64]
    os: [darwin]

  '@img/sharp-libvips-darwin-x64@1.3.2':
    resolution: {integrity: sha512-m2pW1n6cns9VaubNwsZ+c3CRYjxNQWgJ5gPlnL1nbBcpkBvFm6SCFN5o0psFHI8w9n11NKhFkeEDns98tiqbEw==}
    cpu: [x64]
    os: [darwin]

  '@img/sharp-libvips-linux-arm64@1.3.2':
    resolution: {integrity: sha512-dqVSFynCox4C/J8kT16V7SIFAns0IjgLwkvYT7p8LQVmJ5OS5b6tI9IGflxTeuBS//zXeFIUbwt5dwxyZ17cnA==}
    cpu: [arm64]
    os: [linux]
    libc: [glibc]

  '@img/sharp-libvips-linux-arm@1.3.2':
    resolution: {integrity: sha512-1eMLzy92I4J6rmi4mAT8yC3HxOtniyGELlzGbNMLLeqe052ahFQ0h6LFq+lh5DsDIdYViIDst08abvSbcEdLXQ==}
    cpu: [arm]
    os: [linux]
    libc: [glibc]

  '@img/sharp-libvips-linux-ppc64@1.3.2':
    resolution: {integrity: sha512-3z0NHDxD6n5I9gc05U1eW1AyRm+Gznzq3naMrthPNqE6oYykcogW0l/jfpJdjYnuNl8R7yI9pNbE1XiUeyq0Aw==}
    cpu: [ppc64]
    os: [linux]
    libc: [glibc]

  '@img/sharp-libvips-linux-riscv64@1.3.2':
    resolution: {integrity: sha512-bsb4rI+NldGOsXuej2r8OdSS8+zXDVaCWxyWrcv6kneTOlgAHtZABRzBBCwdsPiD90J4myNJuHpg6kA20ImW/w==}
    cpu: [riscv64]
    os: [linux]
    libc: [glibc]

  '@img/sharp-libvips-linux-s390x@1.3.2':
    resolution: {integrity: sha512-/ABshyj8gCpyIrNXnHn4LorDJ0HHm1VhXPBlxZ8zAtfVPAaSafXPGn+sUSIRiwaSBy0mmFjSjiXI5mkcwdChKQ==}
    cpu: [s390x]
    os: [linux]
    libc: [glibc]

  '@img/sharp-libvips-linux-x64@1.3.2':
    resolution: {integrity: sha512-ITPEtgffGJ0S6G9dRyw/366tJQqFRcHWPHhC+Stpg3Z8AEMrDrTr2lhdz4f/Y/HMbRh//7Z5mBzEpVdi62Oc3w==}
    cpu: [x64]
    os: [linux]
    libc: [glibc]

  '@img/sharp-libvips-linuxmusl-arm64@1.3.2':
    resolution: {integrity: sha512-zE9EdiUzUmg5mDT5a1rk5fYJ6GWPloTwWBYDS14naqHsL+EaMpDj1AWnpLgh3u0YCORv2Tt50wrcrpYqkP97Kw==}
    cpu: [arm64]
    os: [linux]
    libc: [musl]

  '@img/sharp-libvips-linuxmusl-x64@1.3.2':
    resolution: {integrity: sha512-m0lrLiUt+lBYnCFr8qV/65yMR4E/c7/wf78I5eKTdkEakFAlZ9QlzEM3QIhhAwVeUhLAHLcCq7a7Vszq/oFNZQ==}
    cpu: [x64]
    os: [linux]
    libc: [musl]

  '@img/sharp-linux-arm64@0.35.3':
    resolution: {integrity: sha512-QgKDspHPnrU+GQ55XPhGwyhC8acLVOOSyAvo1oVfFmrIXLkDNmGWzAfDZ4xK8oSA1qBQrALcHX0G5UZni/SuFQ==}
    engines: {node: '>=20.9.0'}
    cpu: [arm64]
    os: [linux]
    libc: [glibc]

  '@img/sharp-linux-arm@0.35.3':
    resolution: {integrity: sha512-affVWCTLooy8TSxbDx2qkzuDeaWLNVBA+P//FNBirHsXpP2fuBhk5AuboYUnrDnzoXes8GFjpTx0SBFOCRg+FA==}
    engines: {node: '>=20.9.0'}
    cpu: [arm]
    os: [linux]
    libc: [glibc]

  '@img/sharp-linux-ppc64@0.35.3':
    resolution: {integrity: sha512-sMd8rDxmpLOwv/7N44klFjOD5DUO7FLdjiXDI0hoxYaf7Ar262dQIEkosE98bps+5HPLtp/EvNqeqQtOycP/IA==}
    engines: {node: '>=20.9.0'}
    cpu: [ppc64]
    os: [linux]
    libc: [glibc]

  '@img/sharp-linux-riscv64@0.35.3':
    resolution: {integrity: sha512-0Eob78yjlYPfL5vMNWAW55l3R9Y6BQS/gOfe0ZcP9mEz9ohhKSt4im1hayiknXgf8AWrFqMvJcKIdmLmEe7yeQ==}
    engines: {node: '>=20.9.0'}
    cpu: [riscv64]
    os: [linux]
    libc: [glibc]

  '@img/sharp-linux-s390x@0.35.3':
    resolution: {integrity: sha512-KgAxQ0DxpNOq1rG2t5cgTgShJFGSuU7XO45cqC+1NVOuZnP6tlgZRuSYOfNupGkHID0o3cJOsw4DVeJpMovcGw==}
    engines: {node: '>=20.9.0'}
    cpu: [s390x]
    os: [linux]
    libc: [glibc]

  '@img/sharp-linux-x64@0.35.3':
    resolution: {integrity: sha512-8pqvxubL2PGdhlPy6GLqzDYMUjyRmKAwKHYKixpdJYBUK7PJ0C029XdsnpFIdgRZG68fZiGdHVWcKPvtiPB4cA==}
    engines: {node: '>=20.9.0'}
    cpu: [x64]
    os: [linux]
    libc: [glibc]

  '@img/sharp-linuxmusl-arm64@0.35.3':
    resolution: {integrity: sha512-Vz0iQjzzcSX3HCbfwFfCSG/9SCIqyO0mH2sXyiHaAYfBk0cRsCWXRyQYX0ovCK/PAQBbTzQ0dsPQHh5MAFL59w==}
    engines: {node: '>=20.9.0'}
    cpu: [arm64]
    os: [linux]
    libc: [musl]

  '@img/sharp-linuxmusl-x64@0.35.3':
    resolution: {integrity: sha512-6O1NPKcDVj9QEdg7Hx549EX8U0rp6yXQERqru6yRN7fGBn32UvIRJUlWnk+8xDCiG76hXVBbX82NZ/ZKr0euIg==}
    engines: {node: '>=20.9.0'}
    cpu: [x64]
    os: [linux]
    libc: [musl]

  '@img/sharp-wasm32@0.35.3':
    resolution: {integrity: sha512-cZ0XkcYGpHZkqW6iCkqTcmUC0CD9DhD5d/qeZlZkfRBn6GnHniZXLUo5+9xw8Iv76YE6LQFN9YNBlKREcCG76w==}
    engines: {node: '>=20.9.0'}

  '@img/sharp-webcontainers-wasm32@0.35.3':
    resolution: {integrity: sha512-2rnq7bX3NzeR2T4YWgz8qiG4h3TSdMe+vN1iQXpJleSJ3SM5zQ8Fy2SyyXAWlbxpEZ2Y+Z4u1BePgJEYbSy80Q==}
    engines: {node: '>=20.9.0'}
    cpu: [wasm32]

  '@img/sharp-win32-arm64@0.35.3':
    resolution: {integrity: sha512-4bPwFdMbeC4JQ8L8LOyWp6nsHcboP5fxkp6iPOXz2Vg49R42TuMs2whkJ5OAP4/Ul035qOzy0AecOF9VOscn4w==}
    engines: {node: '>=20.9.0'}
    cpu: [arm64]
    os: [win32]

  '@img/sharp-win32-ia32@0.35.3':
    resolution: {integrity: sha512-r53mXsBN6lFUDiST764SvgwUdHAqM4rPAiDzAmf4fLoB6X/rkfyTrLCg6+g17wJJiCmB3JYgHuUldCWUIRFSXw==}
    engines: {node: ^20.9.0}
    cpu: [ia32]
    os: [win32]

  '@img/sharp-win32-x64@0.35.3':
    resolution: {integrity: sha512-D4y1vNeZrIIJCN+uHaWVtH86B+aCrdMYYjicy9pXHvbGZeGYLLSd3wdVuC37FxVXlU1ARsk84eKWfWMXGYEqvA==}
    engines: {node: '>=20.9.0'}
    cpu: [x64]
    os: [win32]

  '@jridgewell/sourcemap-codec@1.5.5':
    resolution: {integrity: sha512-cYQ9310grqxueWbl+WuIUIaiUaDcj7WOq5fVhEljNVgRfOUhY9fy2zTvfoqWsnebh8Sl70VScFbICvJnLKB0Og==}

  '@next/env@16.3.2':
    resolution: {integrity: sha512-8k4YoG8cM7LWlkfzGNYCRBbFNlernLiMw4s0btVl+CmmWqn3VpYypA72/5Feb1UWdxe6tHqr5KHP4p4Y4m9luA==}

  '@next/swc-darwin-arm64@16.3.2':
    resolution: {integrity: sha512-ib5Llm93YCKoKWDh6ZaHq6QWTuOZ2bRkSnUwMmX8dsRIOkBNL1vVlSiUKSfixPL9SSh9pvukzqajk/klkn5vqg==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [darwin]

  '@next/swc-darwin-x64@16.3.2':
    resolution: {integrity: sha512-qd98fX2+I5nYJDioW2o7nSjoxM5KvWdeDefM80igia4+C/qSIEhH4MhTE+hO/7qKM7W37/Mq+dOWp8UePSyLHw==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [darwin]

  '@next/swc-linux-arm64-gnu@16.3.2':
    resolution: {integrity: sha512-vqsgb6FAOzcrCccsLXiKtAy5t8EzO+uOazuFaSkQxeY0tNONG3vpHYy8pyBafcI5SNFPTeyard6yTr6SzNGo2A==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [linux]
    libc: [glibc]

  '@next/swc-linux-arm64-musl@16.3.2':
    resolution: {integrity: sha512-xIe1eujfHUB2XcxHGddxJyu6TJRPjC5NpIkQYB/32ESkt5VkQyIAjmLRS38c+s6QY+qjtY/4KarVDzXRuD7lZQ==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [linux]
    libc: [musl]

  '@next/swc-linux-x64-gnu@16.3.2':
    resolution: {integrity: sha512-Fe0SA2j8X0kmc3aveuHD7UktO3AE2+mH3LguP60vGbz7u0z+MrDXbeb5iZFYAwR7EzzzXJ2Yk966w9mGTFMqfA==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [linux]
    libc: [glibc]

  '@next/swc-linux-x64-musl@16.3.2':
    resolution: {integrity: sha512-TFBipb+gyesI/2Ve4zVu7kGltBWN/R466G5/1gtt2lECfc22G1pjkTxu68Q9aFcOaXiRGTQfvDbQQFe7mYgxiQ==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [linux]
    libc: [musl]

  '@next/swc-win32-arm64-msvc@16.3.2':
    resolution: {integrity: sha512-rVtmnNpBYIosDnKD/96dKxFsJnwnn1WRGG/HioSe8XCm2ksSHNrd2R6+hSjvTBxeMNhJ9pYeu/90cWB1nQLuNA==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [win32]

  '@next/swc-win32-x64-msvc@16.3.2':
    resolution: {integrity: sha512-H4Y2o2/JcHu8LtwzD5CXfHhwxwz8gfsx2HXDEw46Mtev5xHnEmB7HNtZtmriw5ReUOjRtcDqo7XSbU01FT9NlA==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [win32]

  '@oxc-project/types@0.146.0':
    resolution: {integrity: sha512-XC0QsnnhVe7sLIWmYmdPw7x5P0h4W8vUU3Nv1ySgWXtvCz8NizoAEpGXA0sOYoJQV2Rl13LgURAHQ5cI5ILCSA==}

  '@rolldown/binding-android-arm-eabi@1.2.5':
    resolution: {integrity: sha512-DLe/i+l8ynIBY7XEQ191TeZvCoowIGa18R+dIV30GW7DiOtp74i/xX8hs8GUjW5ARV7VZuie3d6AumSmCwbeRA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm]
    os: [android]

  '@rolldown/binding-android-arm64@1.2.5':
    resolution: {integrity: sha512-zXcwKlQApYAOELHd8PwKDFkagYF9Wy4e0RJ+0qnzl9Pjnpj75TEG8ufv40p2J7kCEfwZAsNiuzRIyNNMWT38ig==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [android]

  '@rolldown/binding-darwin-arm64@1.2.5':
    resolution: {integrity: sha512-dK4QakI42nzWgJT5sm4y4y/O//D4OxM75/cH28RLV+nzIN9AY+YsbuUVrUTjlLjXR6vpyxFbSsbmNuJ6BP9sww==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [darwin]

  '@rolldown/binding-darwin-x64@1.2.5':
    resolution: {integrity: sha512-fqSALaUu1Wjd1nK2uW2kJDWdLCc8lx1IcY+MTY26Aurfdx19anlzhqXOgCFbBFQnlFDTn4TC1/7Nz4Bl2mLP3A==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [darwin]

  '@rolldown/binding-freebsd-x64@1.2.5':
    resolution: {integrity: sha512-/vCnNxlkxs9tKxNDcyWUePpJ/PgTzxIaVhoM5SmG8UV+GR/IcPam4VYxi7GIMo7PSDuNqlJqvprqii9NqqVCMw==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [freebsd]

  '@rolldown/binding-linux-arm-gnueabihf@1.2.5':
    resolution: {integrity: sha512-abk0NLA519LxRCszmbE0jYKuQ9YPocOXTiOXOo6Yr+YAT95VH+PtqYAjOJvGKt3viEd/x4qzabAlwd5bHOOARg==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm]
    os: [linux]

  '@rolldown/binding-linux-arm64-gnu@1.2.5':
    resolution: {integrity: sha512-Y7eALiJ8lr0M2HH103Js+g7V34wf6snlpZLAsHI90uLhr3PVlNsbFVAXJC9d/V6BnPyKtpSwI+NcB/RLxsQxuA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [linux]
    libc: [glibc]

  '@rolldown/binding-linux-arm64-musl@1.2.5':
    resolution: {integrity: sha512-xMvZgnbZg4YVnR/AX2b3oOPDTFYJvUVaJg5FedA/LuvexAtXibZQej4cnTkw3rjsJ/ggUROB64TdtETiim+FYA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [linux]
    libc: [musl]

  '@rolldown/binding-linux-ppc64-gnu@1.2.5':
    resolution: {integrity: sha512-GRjeqTUDHTo5GwntsLaAMcBahG3nlpjftXWZLN73HiYQlhwEowvarFgQnRnQZtIp4keXX7quXFbG38uPZBa2EA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [ppc64]
    os: [linux]
    libc: [glibc]

  '@rolldown/binding-linux-s390x-gnu@1.2.5':
    resolution: {integrity: sha512-vLNTR45F2Uwc8AufkNXPmB4VliaXs+FvcheEogIzOXzO4l+LzieXF5A/TWxLy5HtqpsRCHUfd0lPVrrdgXdLHQ==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [s390x]
    os: [linux]
    libc: [glibc]

  '@rolldown/binding-linux-x64-gnu@1.2.5':
    resolution: {integrity: sha512-Mgj59/HTuYeK9Gz2MA+mBWKnHsAgkBSec15ZMb1st3oIfFbX7gCjOae7GydHhzcyQi9Z/7M1QuN9bR3oFqF0jQ==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [linux]
    libc: [glibc]

  '@rolldown/binding-linux-x64-musl@1.2.5':
    resolution: {integrity: sha512-mY8AP0/ichsbhAxGnLa3d3+MwV0EfgrPND2bplI3Ym8T6R2pJ0N87bvrKVwNXmdy3jnr6eQBecdqx/HMknBmpA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [linux]
    libc: [musl]

  '@rolldown/binding-openharmony-arm64@1.2.5':
    resolution: {integrity: sha512-8SLssA2oweAxyRgDp789ACfRb/3P+zNRJpzZxSizxF9m8NUDQ4+3xjo8ttjhVGGw6Qxb70oZiEtIjaKikCO7Yw==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [openharmony]

  '@rolldown/binding-win32-arm64-msvc@1.2.5':
    resolution: {integrity: sha512-vGbruD5zquhoc8D9SViXgN2FBJtNdTyQ4DtG+SWiEGlJiAzoKcZ2xp+xuXCffhubVdt0NJlTZqkeRuERy7g8Cw==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [win32]

  '@rolldown/binding-win32-x64-msvc@1.2.5':
    resolution: {integrity: sha512-e/SXpgISz+IoqVcSSI0rx/d/he8zqLex+/rCWpnHpmVfmPIUjag9H6P7zotf0gJHwPUhQxZ/mF8tr6acebT9yw==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [win32]

  '@rolldown/pluginutils@1.0.1':
    resolution: {integrity: sha512-2j9bGt5Jh8hj+vPtgzPtl72j0yRxHAyumoo6TNfAjsLB04UtpSvPbPcDcBMxz7n+9CYB0c1GxQFxYRg2jimqGw==}

  '@standard-schema/spec@1.1.0':
    resolution: {integrity: sha512-l2aFy5jALhniG5HgqrD6jXLi/rUWrKvqN/qJx6yoJsgKhblVd+iqqU4RCXavm/jPityDo5TCvKMnpjKnOriy0w==}

  '@swc/helpers@0.5.23':
    resolution: {integrity: sha512-5lSsMOTXURePglDfvuAQUqkGek9Hg2kksOYay2m0+XR++b2NWYL/4sWyuvVBIs8oKnJaxkdi9whaL/sqN13afw==}

  '@types/chai@5.2.3':
    resolution: {integrity: sha512-Mw558oeA9fFbv65/y4mHtXDs9bPnFMZAL/jxdPFUpOHHIXX91mcgEHbS5Lahr+pwZFR8A7GQleRWeI6cGFC2UA==}

  '@types/deep-eql@4.0.2':
    resolution: {integrity: sha512-c9h9dVVMigMPc4bwTvC5dxqtqJZwQPePsWjPlpSOnojbor6pGqdk541lfA7AqFQr5pB1BRdq0juY9db81BwyFw==}

  '@types/estree@1.0.9':
    resolution: {integrity: sha512-GhdPgy1el4/ImP05X05Uw4cw2/M93BCUmnEvWZNStlCzEKME4Fkk+YpoA5OiHNQmoS7Cafb8Xa3Pya8m1Qrzeg==}

  '@types/node@26.2.0':
    resolution: {integrity: sha512-5IviulTZeRNp2vAJ514cc/HUlY5nZ9fCbq9DMyC52BrhFZACo3nI0R7qBxhQmo/d27NFe96ur/b7Wwxklda+kg==}

  '@types/pg@8.23.1':
    resolution: {integrity: sha512-fKVHpikPdg4GKks3JuLEhvwSyvwzF23hnabPy6DD8ljVbC7+6J5dQzdv4arV6jqq57djnMgs1HKBxX4P8aBI3A==}

  '@types/react-dom@19.2.5':
    resolution: {integrity: sha512-fMPwH9v7r/pp43yUd2/Mbiex5KouJwwR3dzHkhLREUC6764VyDsqxhAxv6OFEYR1RhjOyD1naqba8ECDBe7ZQg==}
    peerDependencies:
      '@types/react': ^19.2.0

  '@types/react@19.2.18':
    resolution: {integrity: sha512-AnzbBERsrLKtk2XSfTbYRLjQPdy116Sty4q+T+Bp3IC4l6jNBvreVPAHmpq9qhXQM7CXZPjLVmGMw9sy+hxQ3w==}

  '@typescript/typescript-aix-ppc64@7.0.2':
    resolution: {integrity: sha512-MTKKkWB7p/0E9xi1d1tHtZ5PiLkGEMIq88pK2CubZjOsLtYTLqhgIgi6zepFa+9GHZ6h05NMCkQxGKiPXMxXtQ==}
    engines: {node: '>=16.20.0'}
    cpu: [ppc64]
    os: [aix]

  '@typescript/typescript-darwin-arm64@7.0.2':
    resolution: {integrity: sha512-gowzar9MwS/aRWp6f3a4KUqzRjAZjOsmGNCM6LcTgXum+dBfgsBVMN+AgvOCCbguXyick6LJhpBszxMebJ8syA==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [darwin]

  '@typescript/typescript-darwin-x64@7.0.2':
    resolution: {integrity: sha512-SZ9xZInqApNlNGc9s0W1VSsktYSOe9cFqNOIqmN1Gs8SmkjKZYFt017G4VwPxASInODuAdbTW7sXiFUf893RgA==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [darwin]

  '@typescript/typescript-freebsd-arm64@7.0.2':
    resolution: {integrity: sha512-W5NH4y/J0plIIS5b2xvTEkU7JFxyqdMAOgf+Ilhl0vHQXKO5dZoxd+C/jEtq56c4F3wk71RB4BMRQ2XdI+bwYQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [freebsd]

  '@typescript/typescript-freebsd-x64@7.0.2':
    resolution: {integrity: sha512-UMGDx5sTpzNw3WiPebH7l90IWfJggEd+egHt/q6p7/Cm3zqoV7VxkGXt+3DxPIw8CcmvAB0j3sVVfbhX+M4Tpw==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [freebsd]

  '@typescript/typescript-linux-arm64@7.0.2':
    resolution: {integrity: sha512-Qh4eU4/y3yDjnfjjyPYihMj5/ODIlmt+Bzu17OI+fiSRDW57QmU5SiN63exPRNJPKUzcc1INa1NXdrJ+MqHjUQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [linux]

  '@typescript/typescript-linux-arm@7.0.2':
    resolution: {integrity: sha512-gffT3xPz9sR7j/YJExkyPntrI0P2EP9XbOyWzth2/Gs0RstK+90RBcO0ncXoXy/beYll1SXw846Nf2zdnEz0QQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm]
    os: [linux]

  '@typescript/typescript-linux-loong64@7.0.2':
    resolution: {integrity: sha512-uEHck9i8hoAzXPiYRib1O7miOnz23SxIeVl6F4LXox+qov1K35jHcEW6VHKvZI+pyvl7fZEP4MCU5LYvIq1GuQ==}
    engines: {node: '>=16.20.0'}
    cpu: [loong64]
    os: [linux]

  '@typescript/typescript-linux-mips64el@7.0.2':
    resolution: {integrity: sha512-R4KvAMnE43W5Qeqb0Ly56O3mWMWIAgsMyz36DCaycd5nbg/9kzm0liw3JocfRqyJY0KPmzFjbswozXyW0DnIYA==}
    engines: {node: '>=16.20.0'}
    cpu: [mips64el]
    os: [linux]

  '@typescript/typescript-linux-ppc64@7.0.2':
    resolution: {integrity: sha512-DORx5b3sd/4S7eayxm4FQv+A7CrkUIGRaHiwI8oiHTAI1fAPWhF4J0vAlkC8biAlHSVVwxMQ3tjZ2/DVbnQiiA==}
    engines: {node: '>=16.20.0'}
    cpu: [ppc64]
    os: [linux]

  '@typescript/typescript-linux-riscv64@7.0.2':
    resolution: {integrity: sha512-wf0jqEDOjrPRnKwYRyyJDRo11KMbvMFrU+q4zqKyChODBzvlkbhNQfKvLxQCcwTpdDaXSHZTVuh0JoCrKCUMHQ==}
    engines: {node: '>=16.20.0'}
    cpu: [riscv64]
    os: [linux]

  '@typescript/typescript-linux-s390x@7.0.2':
    resolution: {integrity: sha512-IkwJc3L7yhytWd/ewjyxNDfOmswCm9GWMJT/ue/dU4aZNbwZeYAetq42VyLmsmSjvoX7z74X6ZaYCtzAr0EuGw==}
    engines: {node: '>=16.20.0'}
    cpu: [s390x]
    os: [linux]

  '@typescript/typescript-linux-x64@7.0.2':
    resolution: {integrity: sha512-EYdf2cNg7rgCWJnxCdJ+F3V39O8ihb37eHAu1LK8oAFizgTQbPOK7zHHXbPt8rX24COqODXeI3sIf0fCXG7H/A==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [linux]

  '@typescript/typescript-netbsd-arm64@7.0.2':
    resolution: {integrity: sha512-+polYF4MF04aPpO5FTkHran9yUQDSXqy5GiSDKpsll5jy3l3+g9QLhpf39T+ePtefhXLOGrLl0QIjkQP6VnelA==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [netbsd]

  '@typescript/typescript-netbsd-x64@7.0.2':
    resolution: {integrity: sha512-8YIT0EHM/3dq10ZOVF/A7pc/YSMtbcecct4rWtexrnSCHOPcpC2KTLXfTCR6vDpnSiY12heNb1GiN/wu+T/FyA==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [netbsd]

  '@typescript/typescript-openbsd-arm64@7.0.2':
    resolution: {integrity: sha512-APT8+ClYnuYm1u9+kgGXoMj2VzWzcymwh2gNSQVySHfkRDGOTVkoWLjCmOQSaO+PoqQ57B0flRp9SA+7GnnkzQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [openbsd]

  '@typescript/typescript-openbsd-x64@7.0.2':
    resolution: {integrity: sha512-yX7s+Q0Dln0Dt9tEzZsAjXXR/+ytBM7AlglaqyeMPxQszJ1JhlJdZ6jLA+IzldHtflX81em7lDao1xXu+aRRkg==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [openbsd]

  '@typescript/typescript-sunos-x64@7.0.2':
    resolution: {integrity: sha512-dLJDGaLZ1D4HPQn62u1n8mBDkJREwMsAkCdkwd4Ieqw+x3TUyTsqY0YiBCtE6H6OzzgGk3iuZ3vFWRS+E8/d1g==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [sunos]

  '@typescript/typescript-win32-arm64@7.0.2':
    resolution: {integrity: sha512-Gyl1Vy6OsWesLzmq+EP0Fb7b4Nid5232AvcA2SFcdYreldpNtYFFofPjnt62y9hQy7VTaZp65ICJjuAQRaVcIQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [win32]

  '@typescript/typescript-win32-x64@7.0.2':
    resolution: {integrity: sha512-0BQ3HkAHHlKLSp1qRvf3SUhGpGsDuhB/jgFw75guyqbxJqEaS0Cw/VFO8i2nHglJUzQCRtMMR/IBAKE3ETMC4g==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [win32]

  '@vitest/expect@4.1.11':
    resolution: {integrity: sha512-VX2x5vNJXET47KAFzwERI+KRMtTTCSWTfSMKsW7JsUsXV4psq++e3DvZpuTDOpHcxytiDs6p2nhVb2tVDiiUYw==}

  '@vitest/mocker@4.1.11':
    resolution: {integrity: sha512-2XJVD55d1o5AZous5CCGKS74g/riOj9odEt2bQpCVZeblHyHdnMeFl4jl0XjU21stf4mbjUkew2eXQZt65g5CQ==}
    peerDependencies:
      msw: ^2.4.9
      vite: ^6.0.0 || ^7.0.0 || ^8.0.0
    peerDependenciesMeta:
      msw:
        optional: true
      vite:
        optional: true

  '@vitest/pretty-format@4.1.11':
    resolution: {integrity: sha512-yiZzPbGTS9Sr/JpFl8zHrcIkAofNbFV6k21vIgQN/cY/oxZeXhJv5sc/MBJ5jFKWmWs+oJHw0UXLZjmf931+Vw==}

  '@vitest/runner@4.1.11':
    resolution: {integrity: sha512-LztvUgdwMNJMIkj3hQnnxiC2Xy1zNxq928W/xhjCLaNCzqTZOudjwbQf6v9IntZGPw132i2Lq2rgTRZHD3JHNw==}

  '@vitest/snapshot@4.1.11':
    resolution: {integrity: sha512-pN7ikn1ON7h8ee4gIAp4AzyK+zBtJPzVbqOgu5LCEh4VaJVbPQcgYQYJIMGQPXVeJJq1fnfazis7a5pFNPahog==}

  '@vitest/spy@4.1.11':
    resolution: {integrity: sha512-apNa/prQy2qCeywhnixOHPRCgGNhvg7T4Dapfl1GahLp/R+uhBm5cPyFoNVyqsNd2h1nJxL6BqqdIjiABL60YA==}

  '@vitest/utils@4.1.11':
    resolution: {integrity: sha512-zTCVGpyFsGWBhllOyKlTw/vnr6D9qxsfSDyfbyZmTyjHw5N/VuvzHpHoQjm2ZJzn4RJgx5w4r7V0er69CmLgPQ==}

  assertion-error@2.0.1:
    resolution: {integrity: sha512-Izi8RQcffqCeNVgFigKli1ssklIbpHnCYc6AknXGYoB6grJqyeby7jv12JUQgmTAnIDnbck1uxksT4dzN3PWBA==}
    engines: {node: '>=12'}

  baseline-browser-mapping@2.11.18:
    resolution: {integrity: sha512-1iEmLEYSiE1SeBoAfPo/Mnx3PzfzHUkDK61ASkCpuk3YXugYLH5DYK1SzqV55F8FMI6s0F+/tCP7Polz1QRjxw==}
    engines: {node: '>=6.0.0'}
    hasBin: true

  caniuse-lite@1.0.30001809:
    resolution: {integrity: sha512-xxWVywk6a6Arlk+hymeycyn/VgqEfLDxupvhH/xiY5SJ/18kmi9o6MiO320DCUzypORHLtvh0I4i04tUhCNHNQ==}

  chai@6.2.2:
    resolution: {integrity: sha512-NUPRluOfOiTKBKvWPtSD4PhFvWCqOi0BGStNWs57X9js7XGTprSmFoz5F0tWhR4WPjNeR9jXqdC7/UpSJTnlRg==}
    engines: {node: '>=18'}

  client-only@0.0.1:
    resolution: {integrity: sha512-IV3Ou0jSMzZrd3pZ48nLkT9DA7Ag1pnPzaiQhpW7c3RbcqqzvzzVu+L8gfqMp/8IM2MQtSiqaCxrrcfu8I8rMA==}

  convert-source-map@2.0.0:
    resolution: {integrity: sha512-Kvp459HrV2FEJ1CAsi1Ku+MY3kasH19TFykTz2xWmMeq6bk2NU3XXvfJ+Q61m0xktWwt+1HSYf3JZsTms3aRJg==}

  csstype@3.2.3:
    resolution: {integrity: sha512-z1HGKcYy2xA8AGQfwrn0PAy+PB7X/GSj3UVJW9qKyn43xWa+gl5nXmU4qqLMRzWVLFC8KusUX8T/0kCiOYpAIQ==}

  detect-libc@2.1.2:
    resolution: {integrity: sha512-Btj2BOOO83o3WyH59e8MgXsxEQVcarkUOpEYrubB0urwnN10yQ364rsiByU11nZlqWYZm05i/of7io4mzihBtQ==}
    engines: {node: '>=8'}

  es-module-lexer@2.3.2:
    resolution: {integrity: sha512-poHGpORABojJJucnV9KbOavETW8lBVnphkW77ER5/BQ5Fz7oXSoCNek7IH3vR5nRjdsEz926ibFYX8KtLQmdyw==}

  esbuild@0.28.2:
    resolution: {integrity: sha512-HKVLS8dvII+xoKW9kmqxbRKrnWEXfJJr/FZhhJmiqIB0e053QNYFqOBouTMO/k5sID4MvCiUCvv8b9M4h32wIA==}
    engines: {node: '>=18'}
    hasBin: true

  estree-walker@3.0.3:
    resolution: {integrity: sha512-7RUKfXgSMMkzt6ZuXmqapOurLGPPfgj6l9uRZ7lRGolvk0y2yocc35LdcxKC5PQZdn2DMqioAQ2NoWcrTKmm6g==}

  expect-type@1.4.0:
    resolution: {integrity: sha512-KfYbmpRm0VbLjEvVa9yGwCi9GI34xvi7A/HXYWQO65CSD2u3MczUJSuwXKFIxlGsgBQizV9q5J9NHj4VG0n+pA==}
    engines: {node: '>=12.0.0'}

  fdir@6.5.0:
    resolution: {integrity: sha512-tIbYtZbucOs0BRGqPJkshJUYdL+SDH7dVM8gjy+ERp3WAUjLEFJE+02kanyHtwjWOnwrKYBiwAmM0p4kLJAnXg==}
    engines: {node: '>=12.0.0'}
    peerDependencies:
      picomatch: ^3 || ^4
    peerDependenciesMeta:
      picomatch:
        optional: true

  fsevents@2.3.3:
    resolution: {integrity: sha512-5xoDfX+fL7faATnagmWPpbFtwh/R77WmMMqqHGS65C3vvB0YHrgF+B1YmZ3441tMj5n63k0212XNoJwzlhffQw==}
    engines: {node: ^8.16.0 || ^10.6.0 || >=11.0.0}
    os: [darwin]

  jose@6.2.10:
    resolution: {integrity: sha512-iiW7J9qRFlGxvCOIBDBDxFePQSn7ZMAnrYGhrrOo6siO/MIqwfyilLR27pkfDgUk+raLuzADS8A3S/KLBisc0g==}

  lightningcss-android-arm64@1.33.0:
    resolution: {integrity: sha512-gEpRTalKdosp4Bb8qWtc2iOgE5SeIHlpS1up9bFq2wAyYhl1UdTObYiHe98zEM9SQvSoqQZ1IQD0JNpg3Ml5pg==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [android]

  lightningcss-darwin-arm64@1.33.0:
    resolution: {integrity: sha512-Sciaz8eenNTKn9b3t7+xr0ipTp9YxKQY4npwQ3mrRuL0BAVHBLyZxofhaKBAVtzmtRZ/zTyo0/to4B1uWG/Djg==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [darwin]

  lightningcss-darwin-x64@1.33.0:
    resolution: {integrity: sha512-Z5UPAxzrjlWNNyGy6i65cJzzvgJ5D3T6wMvs+gWpY9d7qRhANrxqAp6LhxIgZhWEw18RfJTGcRxjuLIBr+m8XQ==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [darwin]

  lightningcss-freebsd-x64@1.33.0:
    resolution: {integrity: sha512-QQM/Ti/hQajJwCY+RiWuCZ9sdtI/XQk7nDK5vC8kkdwixezOlDgvDx7+RT+QjK6FcFT4MpsuoBnHIo/O3StRRg==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [freebsd]

  lightningcss-linux-arm-gnueabihf@1.33.0:
    resolution: {integrity: sha512-N7FVBe6iS24MlM6R/4RBTxGhQheZGs7tiQ9U32UtF75NzP5Q7xWPRqLBCKxlRQRk3rY1jCIPLzx7WzOhuUIRLQ==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm]
    os: [linux]

  lightningcss-linux-arm64-gnu@1.33.0:
    resolution: {integrity: sha512-j2v/itmy4HlNxlc6voKXYgBqNi0Ng2LShg4z7GufpEgs05P+2suBVyi9I6YHq5uoVFx9ETin3eCEhLVyXGQnKg==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [linux]
    libc: [glibc]

  lightningcss-linux-arm64-musl@1.33.0:
    resolution: {integrity: sha512-yiO5ROMuYQgXbC60yjZU5CYSFZGKXL0HFATXt9mHJn1+zW55oCtMI9NfcVhYLMFDL7gV7oBPon/EmMMGg2OvtQ==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [linux]
    libc: [musl]

  lightningcss-linux-x64-gnu@1.33.0:
    resolution: {integrity: sha512-ar+Ju7LmcN0Jo4FpL4hpFybwNG9/3A/Br5KW2n2jyODg3MEZXaDYADdemoNS+BDNfMgKvylJLj4S5tyRActuAg==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [linux]
    libc: [glibc]

  lightningcss-linux-x64-musl@1.33.0:
    resolution: {integrity: sha512-RYiYbkokw0trfKqqzfF55lginwEPrD3OJDfTuJzFs1MK6iFnDenaz1fqLLtX4ITG3OktJQXOeTaw1awrBAlZPw==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [linux]
    libc: [musl]

  lightningcss-win32-arm64-msvc@1.33.0:
    resolution: {integrity: sha512-1K+MPfLSFVpphzpdbfkhlWk6wBrTObBzS2T6db10PNOZgR9GoVsAWzwNyuhUYYbTp23j+4RrncfujZ4uAzXvwA==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [win32]

  lightningcss-win32-x64-msvc@1.33.0:
    resolution: {integrity: sha512-OlEICDx/Xl0FqSp4bry8zFnCvGpig3Gl4gCquvYwHuqJKEC1+n9NgDniFvqHGmMv1ZkqDJrDqKKSykTDX+ehuA==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [win32]

  lightningcss@1.33.0:
    resolution: {integrity: sha512-WkUDrojuJs0xkgGf2udWxa3yGBRxPtxUkB79i6aCZLRgc7PM8fZe9TosfPDcvEpQZbuFASnHYmRLBLUbmLOIIA==}
    engines: {node: '>= 12.0.0'}

  magic-string@0.30.21:
    resolution: {integrity: sha512-vd2F4YUyEXKGcLHoq+TEyCjxueSeHnFxyyjNp80yg0XV4vUhnDer/lvvlqM/arB5bXQN5K2/3oinyCRyx8T2CQ==}

  nanoid@3.3.18:
    resolution: {integrity: sha512-DTg4MJbGMWkfi6VZFdNt2/caMbQy4Ou+Op/hJQvGEWcnVfoA1QA+xzRKAzw9jD6+GVOOeYr/mIcuDSdug6F6+w==}
    engines: {node: ^10 || ^12 || ^13.7 || ^14 || >=15.0.1}
    hasBin: true

  next@16.3.2:
    resolution: {integrity: sha512-/ZCaubUy17Lld1SiPWxuPbCk2ihqAxF2QNQaPZeEaEb7t1I58qhsJN187D7AfpapHAqUPXH0f/thtdW9dWgWFg==}
    engines: {node: '>=20.9.0'}
    hasBin: true
    peerDependencies:
      '@opentelemetry/api': ^1.1.0
      '@playwright/test': ^1.51.1
      babel-plugin-react-compiler: '*'
      react: ^18.2.0 || 19.0.0-rc-de68d2f4-20241204 || ^19.0.0
      react-dom: ^18.2.0 || 19.0.0-rc-de68d2f4-20241204 || ^19.0.0
      sass: ^1.3.0
    peerDependenciesMeta:
      '@opentelemetry/api':
        optional: true
      '@playwright/test':
        optional: true
      babel-plugin-react-compiler:
        optional: true
      sass:
        optional: true

  oauth4webapi@3.8.7:
    resolution: {integrity: sha512-4RxcKxXjuItDFZ20RRPf4YTw3kpeXJyCgJFxVzJ068A7PNJ18st2Dg90tlC1LkSDS0GecroagCLHYEIVUhCAkw==}

  obug@2.1.4:
    resolution: {integrity: sha512-4a+OsYv9UktOJKE+l1A4OufDgdRF9PifWj+tJnHURo/P+WOxpG4GzUFL9qCalmWauao6ogiG+QvnCovwPoyAWA==}
    engines: {node: '>=12.20.0'}

  openid-client@6.8.5:
    resolution: {integrity: sha512-jNGC/5wnTYwCcEUe2ss0IRUmVRQcgxM0A1nLb3eX/9llqNbMWOQd2xd+qDAgfVCpA5Qh96Y1cdnkfbva6+bSdA==}

  pathe@2.0.3:
    resolution: {integrity: sha512-WUjGcAqP1gQacoQe+OBJsFA7Ld4DyXuUIjZ5cc75cLHvJ7dtNsTugphxIADwspS+AraAUePCKrSVtPLFj/F88w==}

  pg-cloudflare@1.4.0:
    resolution: {integrity: sha512-Vo7z/6rrQYxpNRylp4Tlob2elzbh+N/MOQbxFVWCxS7oEx6jF53GTJFxK2WWpKuBRkmiin4Mt+xofFDjx09R0A==}

  pg-connection-string@2.14.0:
    resolution: {integrity: sha512-XwWDGcLRGCXAR8F/AM5bG7Q+A3Wm2s6QeEjlOKZLlH3UYcguiqCWKyWXVag5TLTIjR7oOJUY8kcADaZgWPyLeg==}

  pg-int8@1.0.1:
    resolution: {integrity: sha512-WCtabS6t3c8SkpDBUlb1kjOs7l66xsGdKpIPZsg4wR+B3+u9UAum2odSsF9tnvxg80h4ZxLWMy4pRjOsFIqQpw==}
    engines: {node: '>=4.0.0'}

  pg-pool@3.14.0:
    resolution: {integrity: sha512-gKtPkFdQPU3DksooVLi9LsjZxrsBUZIpa+7aVx+LV5pNh0KzP4Zleud2po+ConrxbuXGBJ6Hfer6hdgpIBpBaw==}
    peerDependencies:
      pg: '>=8.0'

  pg-protocol@1.16.0:
    resolution: {integrity: sha512-sILXutLVjCLjcDuOmvhX5e2Z4cS5qG/6Bu3VkpFwdf/633ElGLpEh9bgmuI5I4sqKqkifQiGyiCcx1HdtrK7tg==}

  pg-types@2.2.0:
    resolution: {integrity: sha512-qTAAlrEsl8s4OiEQY69wDvcMIdQN6wdz5ojQiOy6YRMuynxenON0O5oCpJI6lshc6scgAY8qvJ2On/p+CXY0GA==}
    engines: {node: '>=4'}

  pg@8.23.0:
    resolution: {integrity: sha512-Ip2EQCngowJLGOfCwkFhPXU7/ljlhn6Rxlmy4XYfL2Y+vyRM59+8uR2xqRWKdYmbXmxCFOAmKxBuSUCdF34qLg==}
    engines: {node: '>= 16.0.0'}
    peerDependencies:
      pg-native: '>=3.0.1'
    peerDependenciesMeta:
      pg-native:
        optional: true

  pgpass@1.0.5:
    resolution: {integrity: sha512-FdW9r/jQZhSeohs1Z3sI1yxFQNFvMcnmfuj4WBMUTxOrAyLMaTcE1aAMBiTlbMNaXvBCQuVi0R7hd8udDSP7ug==}

  picocolors@1.1.1:
    resolution: {integrity: sha512-xceH2snhtb5M9liqDsmEw56le376mTZkEX/jEb/RxNFyegNul7eNslCXP9FDj/Lcu0X8KEyMceP2ntpaHrDEVA==}

  picomatch@4.0.5:
    resolution: {integrity: sha512-RvwwcruNjI1ncT5xRakeyS9Lf8lcItv34KD+aif+VH9kduAyfYBipGh12274xtenIPZ119/R9BdTBa8gAwSh0A==}
    engines: {node: '>=12'}

  postcss@8.5.23:
    resolution: {integrity: sha512-g50586zr4bZmwFiTlflMu8E0bDTb5I5gertgwAKmsdUlTQIhZtunzUlD1WSzwcVWPoAVpsrA6vlfCD7oXvRwgg==}
    engines: {node: ^10 || ^12 || >=14}

  postcss@8.5.26:
    resolution: {integrity: sha512-u82N74LFzG8ca+dD8puPnplTXoGH4fTPpVGuIbt36G3qvNlkvfD0lEAZSxaly3KX8TS/L1A1gsCEmvKmBcVbkQ==}
    engines: {node: ^10 || ^12 || >=14}

  postgres-array@2.0.0:
    resolution: {integrity: sha512-VpZrUqU5A69eQyW2c5CA1jtLecCsN2U/bD6VilrFDWq5+5UIEVO7nazS3TEcHf1zuPYO/sqGvUvW62g86RXZuA==}
    engines: {node: '>=4'}

  postgres-bytea@1.0.1:
    resolution: {integrity: sha512-5+5HqXnsZPE65IJZSMkZtURARZelel2oXUEO8rH83VS/hxH5vv1uHquPg5wZs8yMAfdv971IU+kcPUczi7NVBQ==}
    engines: {node: '>=0.10.0'}

  postgres-date@1.0.7:
    resolution: {integrity: sha512-suDmjLVQg78nMK2UZ454hAG+OAW+HQPZ6n++TNDUX+L0+uUlLywnoxJKDou51Zm+zTCjrCl0Nq6J9C5hP9vK/Q==}
    engines: {node: '>=0.10.0'}

  postgres-interval@1.2.0:
    resolution: {integrity: sha512-9ZhXKM/rw350N1ovuWHbGxnGh/SNJ4cnxHiM0rxE4VN41wsg8P8zWn9hv/buK00RP4WvlOyr/RBDiptyxVbkZQ==}
    engines: {node: '>=0.10.0'}

  react-dom@19.2.8:
    resolution: {integrity: sha512-rVprimfGBG3DR+Tq0IQG2DT5PxKth1WIGDmj5yPmlzr4YBe7uyE+Du4oVqTDXZSHGGGXRtTJEGSSePyQCMBglQ==}
    peerDependencies:
      react: ^19.2.8

  react@19.2.8:
    resolution: {integrity: sha512-PWaYA1L/q9u2u7xYQi+Y3L3Yfnie7XyLeaJICV1MGD6LprsBxcAqGjYyr0eY3p+QdsA+x/Irkt4Qif8D63+Sbw==}
    engines: {node: '>=0.10.0'}

  rolldown@1.2.5:
    resolution: {integrity: sha512-VD2IE5PUG4Oj8zz2VGykiYd5wbnjdIiSsNQb8Qu5B+noEp+A78mu2iVvpp27g8es14Tk9rofNs5Tku9iQCS4fA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    hasBin: true

  scheduler@0.27.0:
    resolution: {integrity: sha512-eNv+WrVbKu1f3vbYJT/xtiF5syA5HPIMtf9IgY/nKg0sWqzAUEvqY/xm7OcZc/qafLx/iO9FgOmeSAp4v5ti/Q==}

  semver@7.8.5:
    resolution: {integrity: sha512-Y7/KDsb8LjooZpwaqGyulO6DQlksgCncchHGk+sZIY4SBvUocMBEFH5Ur1fI4dV+Jvl0w6cjvucaIi40puRioA==}
    engines: {node: '>=10'}
    hasBin: true

  server-only@0.0.1:
    resolution: {integrity: sha512-qepMx2JxAa5jjfzxG79yPPq+8BuFToHd1hm7kI+Z4zAq1ftQiP7HcxMhDDItrbtwVeLg/cY2JnKnrcFkmiswNA==}

  sharp@0.35.3:
    resolution: {integrity: sha512-ej0zVHuZGHCiABXcNxeYhpRnPNPAcvbG8RMdBAhDAxLKkCRVSpK3Iyu7qbqw3JMzoj0REeM6f3tJLtVwl0023Q==}
    engines: {node: '>=20.9.0'}
    peerDependencies:
      '@types/node': '*'
    peerDependenciesMeta:
      '@types/node':
        optional: true

  siginfo@2.0.0:
    resolution: {integrity: sha512-ybx0WO1/8bSBLEWXZvEd7gMW3Sn3JFlW3TvX1nREbDLRNQNaeNN8WK0meBwPdAaOI7TtRRRJn/Es1zhrrCHu7g==}

  source-map-js@1.2.1:
    resolution: {integrity: sha512-UXWMKhLOwVKb728IUtQPXxfYU+usdybtUrK/8uGE8CQMvrhOpwvzDBwj0QhSL7MQc7vIsISBG8VQ8+IDQxpfQA==}
    engines: {node: '>=0.10.0'}

  split2@4.2.0:
    resolution: {integrity: sha512-UcjcJOWknrNkF6PLX83qcHM6KHgVKNkV62Y8a5uYDVv9ydGQVwAHMKqHdJje1VTWpljG0WYpCDhrCdAOYH4TWg==}
    engines: {node: '>= 10.x'}

  stackback@0.0.2:
    resolution: {integrity: sha512-1XMJE5fQo1jGH6Y/7ebnwPOBEkIEnT4QF32d5R1+VXdXveM0IBMJt8zfaxX1P3QhVwrYe+576+jkANtSS2mBbw==}

  std-env@4.2.0:
    resolution: {integrity: sha512-oCUKSupKTHX53EyjDtuZQ64pjLJ6yYCtpmEw0goYxtjG9KpbRe8KAsl2tBUGU9DyMcJ0RwJ8GqJAFzMXcXW1Rw==}

  styled-jsx@5.1.6:
    resolution: {integrity: sha512-qSVyDTeMotdvQYoHWLNGwRFJHC+i+ZvdBRYosOFgC+Wg1vx4frN2/RG/NA7SYqqvKNLf39P2LSRA2pu6n0XYZA==}
    engines: {node: '>= 12.0.0'}
    peerDependencies:
      '@babel/core': '*'
      babel-plugin-macros: '*'
      react: '>= 16.8.0 || 17.x.x || ^18.0.0-0 || ^19.0.0-0'
    peerDependenciesMeta:
      '@babel/core':
        optional: true
      babel-plugin-macros:
        optional: true

  tinybench@2.9.0:
    resolution: {integrity: sha512-0+DUvqWMValLmha6lr4kD8iAMK1HzV0/aKnCtWb9v9641TnP/MFb7Pc2bxoxQjTXAErryXVgUOfv2YqNllqGeg==}

  tinyexec@1.3.0:
    resolution: {integrity: sha512-QKAl9m8gWWGHV8jZcPeym6j+XULi6tOf1mT83WYJ4Lk2ytW/uwAWkrP0uFsdoYMdueVJ0qs26wZ+23xeB4ibNQ==}
    engines: {node: '>=18'}

  tinyglobby@0.2.17:
    resolution: {integrity: sha512-wXR/dYpcqKmfWpEdZjiKJOwCNFndD0DMnrW/cYjVGttEkBfVgcLFHoNrlj47mjOVic9yyNu65alsgF4NQyTa2g==}
    engines: {node: '>=12.0.0'}

  tinyrainbow@3.1.1:
    resolution: {integrity: sha512-yau8yJdTt989Mm0Bd/236QnzEiPf2xLLTqUZRUJOo/3CB078LSwzei343DgtJVmfJKJE3TMINY1u42SQsP6mXw==}
    engines: {node: '>=14.0.0'}

  tslib@2.8.1:
    resolution: {integrity: sha512-oJFu94HQb+KVduSUQL7wnpmqnfmLsOA/nAh6b6EH0wCEoK0/mPeXU6c3wKDV83MkOuHPRHtSXKKU99IBazS/2w==}

  tsx@4.23.12:
    resolution: {integrity: sha512-FDf4L4sYzKtzWYhU/Xm0AQFdTjdIxNo9ElTf2mxXM6k8YMHXzYUe4yODVaXP4V9uMFbVg8c0qyBccK2OOxb45Q==}
    engines: {node: '>=18.0.0'}
    hasBin: true

  typescript@7.0.2:
    resolution: {integrity: sha512-8FYau96o3NKOhbjKi/qNvG/W5jhzxkbdm5sj9AbZ/5T5sWqn3hJgLfGx27sRKZWTvyzCP8dLRBTf5tBTSRVUNA==}
    engines: {node: '>=16.20.0'}
    hasBin: true

  undici-types@8.3.0:
    resolution: {integrity: sha512-j375ScV60dom+YkPFIfTLcOiPxkN/buHz5GobjLhixFuANaNs3C9l4GmrWqejgXWJ7BbJcFYpTEUkS1Ge8bpZQ==}

  vite@8.2.2:
    resolution: {integrity: sha512-cFKLV/PRgAUlIRm5WjMjJ86jrftzpqcgH+Us+DS8mI3CDNiH30Whrz8uHL3+MOLPAgqbMBAqWdAHAphOAM+z/Q==}
    engines: {node: ^20.19.0 || >=22.12.0}
    hasBin: true
    peerDependencies:
      '@types/node': ^20.19.0 || >=22.12.0
      '@vitejs/devtools': ^0.4.0 || ^0.5.0
      esbuild: ^0.27.0 || ^0.28.0
      jiti: '>=1.21.0'
      less: ^4.0.0
      sass: ^1.70.0
      sass-embedded: ^1.70.0
      stylus: '>=0.54.8'
      sugarss: ^5.0.0
      terser: ^5.16.0
      tsx: ^4.8.1
      yaml: ^2.4.2
    peerDependenciesMeta:
      '@types/node':
        optional: true
      '@vitejs/devtools':
        optional: true
      esbuild:
        optional: true
      jiti:
        optional: true
      less:
        optional: true
      sass:
        optional: true
      sass-embedded:
        optional: true
      stylus:
        optional: true
      sugarss:
        optional: true
      terser:
        optional: true
      tsx:
        optional: true
      yaml:
        optional: true

  vitest@4.1.11:
    resolution: {integrity: sha512-fhACrNXUidIbGSBr5FlbuBkO7VWC1ZyLl0DO4CU2DrQoAPxX84Ysxs+HeGQpii5lZWV1Q4gBZTTu49mF+A6Edw==}
    engines: {node: ^20.0.0 || ^22.0.0 || >=24.0.0}
    hasBin: true
    peerDependencies:
      '@edge-runtime/vm': '*'
      '@opentelemetry/api': ^1.9.0
      '@types/node': ^20.0.0 || ^22.0.0 || >=24.0.0
      '@vitest/browser-playwright': 4.1.11
      '@vitest/browser-preview': 4.1.11
      '@vitest/browser-webdriverio': 4.1.11
      '@vitest/coverage-istanbul': 4.1.11
      '@vitest/coverage-v8': 4.1.11
      '@vitest/ui': 4.1.11
      happy-dom: '*'
      jsdom: '*'
      vite: ^6.0.0 || ^7.0.0 || ^8.0.0
    peerDependenciesMeta:
      '@edge-runtime/vm':
        optional: true
      '@opentelemetry/api':
        optional: true
      '@types/node':
        optional: true
      '@vitest/browser-playwright':
        optional: true
      '@vitest/browser-preview':
        optional: true
      '@vitest/browser-webdriverio':
        optional: true
      '@vitest/coverage-istanbul':
        optional: true
      '@vitest/coverage-v8':
        optional: true
      '@vitest/ui':
        optional: true
      happy-dom:
        optional: true
      jsdom:
        optional: true

  why-is-node-running@2.3.0:
    resolution: {integrity: sha512-hUrmaWBdVDcxvYqnyh09zunKzROWjbZTiNy8dBEjkS7ehEDQibXJ7XvlmtbwuTclUiIyN+CyXQD4Vmko8fNm8w==}
    engines: {node: '>=8'}
    hasBin: true

  xtend@4.0.2:
    resolution: {integrity: sha512-LKYU1iAXJXUgAXn9URjiu+MWhyUXHsvfp7mcuYm9dSUKK0/CjtrUwFAxD82/mCWbtLsGjFIad0wIsod4zrTAEQ==}
    engines: {node: '>=0.4'}

  zod@4.4.3:
    resolution: {integrity: sha512-ytENFjIJFl2UwYglde2jchW2Hwm4GJFLDiSXWdTrJQBIN9Fcyp7n4DhxJEiWNAJMV1/BqWfW/kkg71UDcHJyTQ==}

snapshots:

  '@electric-sql/pglite@0.5.7': {}

  '@emnapi/runtime@1.11.3':
    dependencies:
      tslib: 2.8.1
    optional: true

  '@esbuild/aix-ppc64@0.28.2':
    optional: true

  '@esbuild/android-arm64@0.28.2':
    optional: true

  '@esbuild/android-arm@0.28.2':
    optional: true

  '@esbuild/android-x64@0.28.2':
    optional: true

  '@esbuild/darwin-arm64@0.28.2':
    optional: true

  '@esbuild/darwin-x64@0.28.2':
    optional: true

  '@esbuild/freebsd-arm64@0.28.2':
    optional: true

  '@esbuild/freebsd-x64@0.28.2':
    optional: true

  '@esbuild/linux-arm64@0.28.2':
    optional: true

  '@esbuild/linux-arm@0.28.2':
    optional: true

  '@esbuild/linux-ia32@0.28.2':
    optional: true

  '@esbuild/linux-loong64@0.28.2':
    optional: true

  '@esbuild/linux-mips64el@0.28.2':
    optional: true

  '@esbuild/linux-ppc64@0.28.2':
    optional: true

  '@esbuild/linux-riscv64@0.28.2':
    optional: true

  '@esbuild/linux-s390x@0.28.2':
    optional: true

  '@esbuild/linux-x64@0.28.2':
    optional: true

  '@esbuild/netbsd-arm64@0.28.2':
    optional: true

  '@esbuild/netbsd-x64@0.28.2':
    optional: true

  '@esbuild/openbsd-arm64@0.28.2':
    optional: true

  '@esbuild/openbsd-x64@0.28.2':
    optional: true

  '@esbuild/openharmony-arm64@0.28.2':
    optional: true

  '@esbuild/sunos-x64@0.28.2':
    optional: true

  '@esbuild/win32-arm64@0.28.2':
    optional: true

  '@esbuild/win32-ia32@0.28.2':
    optional: true

  '@esbuild/win32-x64@0.28.2':
    optional: true

  '@img/colour@1.1.0':
    optional: true

  '@img/sharp-darwin-arm64@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-darwin-arm64': 1.3.2
    optional: true

  '@img/sharp-darwin-x64@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-darwin-x64': 1.3.2
    optional: true

  '@img/sharp-freebsd-wasm32@0.35.3':
    dependencies:
      '@img/sharp-wasm32': 0.35.3
    optional: true

  '@img/sharp-libvips-darwin-arm64@1.3.2':
    optional: true

  '@img/sharp-libvips-darwin-x64@1.3.2':
    optional: true

  '@img/sharp-libvips-linux-arm64@1.3.2':
    optional: true

  '@img/sharp-libvips-linux-arm@1.3.2':
    optional: true

  '@img/sharp-libvips-linux-ppc64@1.3.2':
    optional: true

  '@img/sharp-libvips-linux-riscv64@1.3.2':
    optional: true

  '@img/sharp-libvips-linux-s390x@1.3.2':
    optional: true

  '@img/sharp-libvips-linux-x64@1.3.2':
    optional: true

  '@img/sharp-libvips-linuxmusl-arm64@1.3.2':
    optional: true

  '@img/sharp-libvips-linuxmusl-x64@1.3.2':
    optional: true

  '@img/sharp-linux-arm64@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-linux-arm64': 1.3.2
    optional: true

  '@img/sharp-linux-arm@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-linux-arm': 1.3.2
    optional: true

  '@img/sharp-linux-ppc64@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-linux-ppc64': 1.3.2
    optional: true

  '@img/sharp-linux-riscv64@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-linux-riscv64': 1.3.2
    optional: true

  '@img/sharp-linux-s390x@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-linux-s390x': 1.3.2
    optional: true

  '@img/sharp-linux-x64@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-linux-x64': 1.3.2
    optional: true

  '@img/sharp-linuxmusl-arm64@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-linuxmusl-arm64': 1.3.2
    optional: true

  '@img/sharp-linuxmusl-x64@0.35.3':
    optionalDependencies:
      '@img/sharp-libvips-linuxmusl-x64': 1.3.2
    optional: true

  '@img/sharp-wasm32@0.35.3':
    dependencies:
      '@emnapi/runtime': 1.11.3
    optional: true

  '@img/sharp-webcontainers-wasm32@0.35.3':
    dependencies:
      '@img/sharp-wasm32': 0.35.3
    optional: true

  '@img/sharp-win32-arm64@0.35.3':
    optional: true

  '@img/sharp-win32-ia32@0.35.3':
    optional: true

  '@img/sharp-win32-x64@0.35.3':
    optional: true

  '@jridgewell/sourcemap-codec@1.5.5': {}

  '@next/env@16.3.2': {}

  '@next/swc-darwin-arm64@16.3.2':
    optional: true

  '@next/swc-darwin-x64@16.3.2':
    optional: true

  '@next/swc-linux-arm64-gnu@16.3.2':
    optional: true

  '@next/swc-linux-arm64-musl@16.3.2':
    optional: true

  '@next/swc-linux-x64-gnu@16.3.2':
    optional: true

  '@next/swc-linux-x64-musl@16.3.2':
    optional: true

  '@next/swc-win32-arm64-msvc@16.3.2':
    optional: true

  '@next/swc-win32-x64-msvc@16.3.2':
    optional: true

  '@oxc-project/types@0.146.0': {}

  '@rolldown/binding-android-arm-eabi@1.2.5':
    optional: true

  '@rolldown/binding-android-arm64@1.2.5':
    optional: true

  '@rolldown/binding-darwin-arm64@1.2.5':
    optional: true

  '@rolldown/binding-darwin-x64@1.2.5':
    optional: true

  '@rolldown/binding-freebsd-x64@1.2.5':
    optional: true

  '@rolldown/binding-linux-arm-gnueabihf@1.2.5':
    optional: true

  '@rolldown/binding-linux-arm64-gnu@1.2.5':
    optional: true

  '@rolldown/binding-linux-arm64-musl@1.2.5':
    optional: true

  '@rolldown/binding-linux-ppc64-gnu@1.2.5':
    optional: true

  '@rolldown/binding-linux-s390x-gnu@1.2.5':
    optional: true

  '@rolldown/binding-linux-x64-gnu@1.2.5':
    optional: true

  '@rolldown/binding-linux-x64-musl@1.2.5':
    optional: true

  '@rolldown/binding-openharmony-arm64@1.2.5':
    optional: true

  '@rolldown/binding-win32-arm64-msvc@1.2.5':
    optional: true

  '@rolldown/binding-win32-x64-msvc@1.2.5':
    optional: true

  '@rolldown/pluginutils@1.0.1': {}

  '@standard-schema/spec@1.1.0': {}

  '@swc/helpers@0.5.23':
    dependencies:
      tslib: 2.8.1

  '@types/chai@5.2.3':
    dependencies:
      '@types/deep-eql': 4.0.2
      assertion-error: 2.0.1

  '@types/deep-eql@4.0.2': {}

  '@types/estree@1.0.9': {}

  '@types/node@26.2.0':
    dependencies:
      undici-types: 8.3.0

  '@types/pg@8.23.1':
    dependencies:
      '@types/node': 26.2.0
      pg-protocol: 1.16.0
      pg-types: 2.2.0

  '@types/react-dom@19.2.5(@types/react@19.2.18)':
    dependencies:
      '@types/react': 19.2.18

  '@types/react@19.2.18':
    dependencies:
      csstype: 3.2.3

  '@typescript/typescript-aix-ppc64@7.0.2':
    optional: true

  '@typescript/typescript-darwin-arm64@7.0.2':
    optional: true

  '@typescript/typescript-darwin-x64@7.0.2':
    optional: true

  '@typescript/typescript-freebsd-arm64@7.0.2':
    optional: true

  '@typescript/typescript-freebsd-x64@7.0.2':
    optional: true

  '@typescript/typescript-linux-arm64@7.0.2':
    optional: true

  '@typescript/typescript-linux-arm@7.0.2':
    optional: true

  '@typescript/typescript-linux-loong64@7.0.2':
    optional: true

  '@typescript/typescript-linux-mips64el@7.0.2':
    optional: true

  '@typescript/typescript-linux-ppc64@7.0.2':
    optional: true

  '@typescript/typescript-linux-riscv64@7.0.2':
    optional: true

  '@typescript/typescript-linux-s390x@7.0.2':
    optional: true

  '@typescript/typescript-linux-x64@7.0.2':
    optional: true

  '@typescript/typescript-netbsd-arm64@7.0.2':
    optional: true

  '@typescript/typescript-netbsd-x64@7.0.2':
    optional: true

  '@typescript/typescript-openbsd-arm64@7.0.2':
    optional: true

  '@typescript/typescript-openbsd-x64@7.0.2':
    optional: true

  '@typescript/typescript-sunos-x64@7.0.2':
    optional: true

  '@typescript/typescript-win32-arm64@7.0.2':
    optional: true

  '@typescript/typescript-win32-x64@7.0.2':
    optional: true

  '@vitest/expect@4.1.11':
    dependencies:
      '@standard-schema/spec': 1.1.0
      '@types/chai': 5.2.3
      '@vitest/spy': 4.1.11
      '@vitest/utils': 4.1.11
      chai: 6.2.2
      tinyrainbow: 3.1.1

  '@vitest/mocker@4.1.11(vite@8.2.2(@types/node@26.2.0)(esbuild@0.28.2)(tsx@4.23.12))':
    dependencies:
      '@vitest/spy': 4.1.11
      estree-walker: 3.0.3
      magic-string: 0.30.21
    optionalDependencies:
      vite: 8.2.2(@types/node@26.2.0)(esbuild@0.28.2)(tsx@4.23.12)

  '@vitest/pretty-format@4.1.11':
    dependencies:
      tinyrainbow: 3.1.1

  '@vitest/runner@4.1.11':
    dependencies:
      '@vitest/utils': 4.1.11
      pathe: 2.0.3

  '@vitest/snapshot@4.1.11':
    dependencies:
      '@vitest/pretty-format': 4.1.11
      '@vitest/utils': 4.1.11
      magic-string: 0.30.21
      pathe: 2.0.3

  '@vitest/spy@4.1.11': {}

  '@vitest/utils@4.1.11':
    dependencies:
      '@vitest/pretty-format': 4.1.11
      convert-source-map: 2.0.0
      tinyrainbow: 3.1.1

  assertion-error@2.0.1: {}

  baseline-browser-mapping@2.11.18: {}

  caniuse-lite@1.0.30001809: {}

  chai@6.2.2: {}

  client-only@0.0.1: {}

  convert-source-map@2.0.0: {}

  csstype@3.2.3: {}

  detect-libc@2.1.2: {}

  es-module-lexer@2.3.2: {}

  esbuild@0.28.2:
    optionalDependencies:
      '@esbuild/aix-ppc64': 0.28.2
      '@esbuild/android-arm': 0.28.2
      '@esbuild/android-arm64': 0.28.2
      '@esbuild/android-x64': 0.28.2
      '@esbuild/darwin-arm64': 0.28.2
      '@esbuild/darwin-x64': 0.28.2
      '@esbuild/freebsd-arm64': 0.28.2
      '@esbuild/freebsd-x64': 0.28.2
      '@esbuild/linux-arm': 0.28.2
      '@esbuild/linux-arm64': 0.28.2
      '@esbuild/linux-ia32': 0.28.2
      '@esbuild/linux-loong64': 0.28.2
      '@esbuild/linux-mips64el': 0.28.2
      '@esbuild/linux-ppc64': 0.28.2
      '@esbuild/linux-riscv64': 0.28.2
      '@esbuild/linux-s390x': 0.28.2
      '@esbuild/linux-x64': 0.28.2
      '@esbuild/netbsd-arm64': 0.28.2
      '@esbuild/netbsd-x64': 0.28.2
      '@esbuild/openbsd-arm64': 0.28.2
      '@esbuild/openbsd-x64': 0.28.2
      '@esbuild/openharmony-arm64': 0.28.2
      '@esbuild/sunos-x64': 0.28.2
      '@esbuild/win32-arm64': 0.28.2
      '@esbuild/win32-ia32': 0.28.2
      '@esbuild/win32-x64': 0.28.2

  estree-walker@3.0.3:
    dependencies:
      '@types/estree': 1.0.9

  expect-type@1.4.0: {}

  fdir@6.5.0(picomatch@4.0.5):
    optionalDependencies:
      picomatch: 4.0.5

  fsevents@2.3.3:
    optional: true

  jose@6.2.10: {}

  lightningcss-android-arm64@1.33.0:
    optional: true

  lightningcss-darwin-arm64@1.33.0:
    optional: true

  lightningcss-darwin-x64@1.33.0:
    optional: true

  lightningcss-freebsd-x64@1.33.0:
    optional: true

  lightningcss-linux-arm-gnueabihf@1.33.0:
    optional: true

  lightningcss-linux-arm64-gnu@1.33.0:
    optional: true

  lightningcss-linux-arm64-musl@1.33.0:
    optional: true

  lightningcss-linux-x64-gnu@1.33.0:
    optional: true

  lightningcss-linux-x64-musl@1.33.0:
    optional: true

  lightningcss-win32-arm64-msvc@1.33.0:
    optional: true

  lightningcss-win32-x64-msvc@1.33.0:
    optional: true

  lightningcss@1.33.0:
    dependencies:
      detect-libc: 2.1.2
    optionalDependencies:
      lightningcss-android-arm64: 1.33.0
      lightningcss-darwin-arm64: 1.33.0
      lightningcss-darwin-x64: 1.33.0
      lightningcss-freebsd-x64: 1.33.0
      lightningcss-linux-arm-gnueabihf: 1.33.0
      lightningcss-linux-arm64-gnu: 1.33.0
      lightningcss-linux-arm64-musl: 1.33.0
      lightningcss-linux-x64-gnu: 1.33.0
      lightningcss-linux-x64-musl: 1.33.0
      lightningcss-win32-arm64-msvc: 1.33.0
      lightningcss-win32-x64-msvc: 1.33.0

  magic-string@0.30.21:
    dependencies:
      '@jridgewell/sourcemap-codec': 1.5.5

  nanoid@3.3.18: {}

  next@16.3.2(@types/node@26.2.0)(react-dom@19.2.8(react@19.2.8))(react@19.2.8):
    dependencies:
      '@next/env': 16.3.2
      '@swc/helpers': 0.5.23
      baseline-browser-mapping: 2.11.18
      caniuse-lite: 1.0.30001809
      postcss: 8.5.23
      react: 19.2.8
      react-dom: 19.2.8(react@19.2.8)
      styled-jsx: 5.1.6(react@19.2.8)
    optionalDependencies:
      '@next/swc-darwin-arm64': 16.3.2
      '@next/swc-darwin-x64': 16.3.2
      '@next/swc-linux-arm64-gnu': 16.3.2
      '@next/swc-linux-arm64-musl': 16.3.2
      '@next/swc-linux-x64-gnu': 16.3.2
      '@next/swc-linux-x64-musl': 16.3.2
      '@next/swc-win32-arm64-msvc': 16.3.2
      '@next/swc-win32-x64-msvc': 16.3.2
      sharp: 0.35.3(@types/node@26.2.0)
    transitivePeerDependencies:
      - '@babel/core'
      - '@types/node'
      - babel-plugin-macros

  oauth4webapi@3.8.7: {}

  obug@2.1.4: {}

  openid-client@6.8.5:
    dependencies:
      jose: 6.2.10
      oauth4webapi: 3.8.7

  pathe@2.0.3: {}

  pg-cloudflare@1.4.0:
    optional: true

  pg-connection-string@2.14.0: {}

  pg-int8@1.0.1: {}

  pg-pool@3.14.0(pg@8.23.0):
    dependencies:
      pg: 8.23.0

  pg-protocol@1.16.0: {}

  pg-types@2.2.0:
    dependencies:
      pg-int8: 1.0.1
      postgres-array: 2.0.0
      postgres-bytea: 1.0.1
      postgres-date: 1.0.7
      postgres-interval: 1.2.0

  pg@8.23.0:
    dependencies:
      pg-connection-string: 2.14.0
      pg-pool: 3.14.0(pg@8.23.0)
      pg-protocol: 1.16.0
      pg-types: 2.2.0
      pgpass: 1.0.5
    optionalDependencies:
      pg-cloudflare: 1.4.0

  pgpass@1.0.5:
    dependencies:
      split2: 4.2.0

  picocolors@1.1.1: {}

  picomatch@4.0.5: {}

  postcss@8.5.23:
    dependencies:
      nanoid: 3.3.18
      picocolors: 1.1.1
      source-map-js: 1.2.1

  postcss@8.5.26:
    dependencies:
      nanoid: 3.3.18
      picocolors: 1.1.1
      source-map-js: 1.2.1

  postgres-array@2.0.0: {}

  postgres-bytea@1.0.1: {}

  postgres-date@1.0.7: {}

  postgres-interval@1.2.0:
    dependencies:
      xtend: 4.0.2

  react-dom@19.2.8(react@19.2.8):
    dependencies:
      react: 19.2.8
      scheduler: 0.27.0

  react@19.2.8: {}

  rolldown@1.2.5:
    dependencies:
      '@oxc-project/types': 0.146.0
      '@rolldown/pluginutils': 1.0.1
    optionalDependencies:
      '@rolldown/binding-android-arm-eabi': 1.2.5
      '@rolldown/binding-android-arm64': 1.2.5
      '@rolldown/binding-darwin-arm64': 1.2.5
      '@rolldown/binding-darwin-x64': 1.2.5
      '@rolldown/binding-freebsd-x64': 1.2.5
      '@rolldown/binding-linux-arm-gnueabihf': 1.2.5
      '@rolldown/binding-linux-arm64-gnu': 1.2.5
      '@rolldown/binding-linux-arm64-musl': 1.2.5
      '@rolldown/binding-linux-ppc64-gnu': 1.2.5
      '@rolldown/binding-linux-s390x-gnu': 1.2.5
      '@rolldown/binding-linux-x64-gnu': 1.2.5
      '@rolldown/binding-linux-x64-musl': 1.2.5
      '@rolldown/binding-openharmony-arm64': 1.2.5
      '@rolldown/binding-win32-arm64-msvc': 1.2.5
      '@rolldown/binding-win32-x64-msvc': 1.2.5

  scheduler@0.27.0: {}

  semver@7.8.5:
    optional: true

  server-only@0.0.1: {}

  sharp@0.35.3(@types/node@26.2.0):
    dependencies:
      '@img/colour': 1.1.0
      detect-libc: 2.1.2
      semver: 7.8.5
    optionalDependencies:
      '@img/sharp-darwin-arm64': 0.35.3
      '@img/sharp-darwin-x64': 0.35.3
      '@img/sharp-freebsd-wasm32': 0.35.3
      '@img/sharp-libvips-darwin-arm64': 1.3.2
      '@img/sharp-libvips-darwin-x64': 1.3.2
      '@img/sharp-libvips-linux-arm': 1.3.2
      '@img/sharp-libvips-linux-arm64': 1.3.2
      '@img/sharp-libvips-linux-ppc64': 1.3.2
      '@img/sharp-libvips-linux-riscv64': 1.3.2
      '@img/sharp-libvips-linux-s390x': 1.3.2
      '@img/sharp-libvips-linux-x64': 1.3.2
      '@img/sharp-libvips-linuxmusl-arm64': 1.3.2
      '@img/sharp-libvips-linuxmusl-x64': 1.3.2
      '@img/sharp-linux-arm': 0.35.3
      '@img/sharp-linux-arm64': 0.35.3
      '@img/sharp-linux-ppc64': 0.35.3
      '@img/sharp-linux-riscv64': 0.35.3
      '@img/sharp-linux-s390x': 0.35.3
      '@img/sharp-linux-x64': 0.35.3
      '@img/sharp-linuxmusl-arm64': 0.35.3
      '@img/sharp-linuxmusl-x64': 0.35.3
      '@img/sharp-webcontainers-wasm32': 0.35.3
      '@img/sharp-win32-arm64': 0.35.3
      '@img/sharp-win32-ia32': 0.35.3
      '@img/sharp-win32-x64': 0.35.3
      '@types/node': 26.2.0
    optional: true

  siginfo@2.0.0: {}

  source-map-js@1.2.1: {}

  split2@4.2.0: {}

  stackback@0.0.2: {}

  std-env@4.2.0: {}

  styled-jsx@5.1.6(react@19.2.8):
    dependencies:
      client-only: 0.0.1
      react: 19.2.8

  tinybench@2.9.0: {}

  tinyexec@1.3.0: {}

  tinyglobby@0.2.17:
    dependencies:
      fdir: 6.5.0(picomatch@4.0.5)
      picomatch: 4.0.5

  tinyrainbow@3.1.1: {}

  tslib@2.8.1: {}

  tsx@4.23.12:
    dependencies:
      esbuild: 0.28.2
    optionalDependencies:
      fsevents: 2.3.3

  typescript@7.0.2:
    optionalDependencies:
      '@typescript/typescript-aix-ppc64': 7.0.2
      '@typescript/typescript-darwin-arm64': 7.0.2
      '@typescript/typescript-darwin-x64': 7.0.2
      '@typescript/typescript-freebsd-arm64': 7.0.2
      '@typescript/typescript-freebsd-x64': 7.0.2
      '@typescript/typescript-linux-arm': 7.0.2
      '@typescript/typescript-linux-arm64': 7.0.2
      '@typescript/typescript-linux-loong64': 7.0.2
      '@typescript/typescript-linux-mips64el': 7.0.2
      '@typescript/typescript-linux-ppc64': 7.0.2
      '@typescript/typescript-linux-riscv64': 7.0.2
      '@typescript/typescript-linux-s390x': 7.0.2
      '@typescript/typescript-linux-x64': 7.0.2
      '@typescript/typescript-netbsd-arm64': 7.0.2
      '@typescript/typescript-netbsd-x64': 7.0.2
      '@typescript/typescript-openbsd-arm64': 7.0.2
      '@typescript/typescript-openbsd-x64': 7.0.2
      '@typescript/typescript-sunos-x64': 7.0.2
      '@typescript/typescript-win32-arm64': 7.0.2
      '@typescript/typescript-win32-x64': 7.0.2

  undici-types@8.3.0: {}

  vite@8.2.2(@types/node@26.2.0)(esbuild@0.28.2)(tsx@4.23.12):
    dependencies:
      lightningcss: 1.33.0
      picomatch: 4.0.5
      postcss: 8.5.26
      rolldown: 1.2.5
      tinyglobby: 0.2.17
    optionalDependencies:
      '@types/node': 26.2.0
      esbuild: 0.28.2
      fsevents: 2.3.3
      tsx: 4.23.12

  vitest@4.1.11(@types/node@26.2.0)(vite@8.2.2(@types/node@26.2.0)(esbuild@0.28.2)(tsx@4.23.12)):
    dependencies:
      '@vitest/expect': 4.1.11
      '@vitest/mocker': 4.1.11(vite@8.2.2(@types/node@26.2.0)(esbuild@0.28.2)(tsx@4.23.12))
      '@vitest/pretty-format': 4.1.11
      '@vitest/runner': 4.1.11
      '@vitest/snapshot': 4.1.11
      '@vitest/spy': 4.1.11
      '@vitest/utils': 4.1.11
      es-module-lexer: 2.3.2
      expect-type: 1.4.0
      magic-string: 0.30.21
      obug: 2.1.4
      pathe: 2.0.3
      picomatch: 4.0.5
      std-env: 4.2.0
      tinybench: 2.9.0
      tinyexec: 1.3.0
      tinyglobby: 0.2.17
      tinyrainbow: 3.1.1
      vite: 8.2.2(@types/node@26.2.0)(esbuild@0.28.2)(tsx@4.23.12)
      why-is-node-running: 2.3.0
    optionalDependencies:
      '@types/node': 26.2.0
    transitivePeerDependencies:
      - msw

  why-is-node-running@2.3.0:
    dependencies:
      siginfo: 2.0.0
      stackback: 0.0.2

  xtend@4.0.2: {}

  zod@4.4.3: {}
````

### FILE: `pnpm-workspace.yaml`

```yaml
block_id: "TS-ENTERPRISE-WEB:pnpm-workspace-yaml:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "8498920f6e5d52571f6528dc67ab7da339543d5e9de970d1ece50e3d1f90b14a"
variables: []
secrets_allowed: false
```

````yaml
allowBuilds:
  # Vitest/tsx depend on esbuild. This exact package is intentionally allowlisted.
  esbuild: true
minimumReleaseAgeExclude:
  - '@types/react-dom@19.2.5'
````

### FILE: `README.md`

```yaml
block_id: "TS-ENTERPRISE-WEB:readme-md:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "3cfd1cca40c7ab03c47a3e48f98c168d92d1fa047667bcd33c39fa0a86c48550"
variables: []
secrets_allowed: false
```

````markdown
# Elite Configurable Enterprise Platform

Estado: **adapter web ejecutable 0.1 y condicionado**, no foundation, backend por defecto, `REUSABLE_PACK` ni producción universal.

Este directorio es una variante web opcional que se materializa después de que el blueprint selecciona el stack. No es una colección de snippets: compila como una aplicación, valida una configuración empresarial y ofrece un vertical slice real de catálogo público y captación de interesados.

## Inicio local

Requisitos: Node.js 24+ y pnpm 11+.

```powershell
Copy-Item .env.example .env.local
pnpm install
pnpm dev:ready
```

Abrir `http://localhost:3000`. El entorno local usa PGlite, una distribución PostgreSQL embebida. Producción se niega a iniciar sin PostgreSQL externo, OIDC y referencias de secretos.

## Verificación reproducible

```powershell
pnpm verify
pnpm bootstrap
pnpm outbox:once
pnpm licenses:report
```

`verify` exige configuración válida, TypeScript estricto, tests y build de producción. El lockfile fija versiones; `pnpm-workspace.yaml` permite ejecutar únicamente el build script de `esbuild`.

`pnpm readiness:baseline` debe pasar para ampliar el adapter. `pnpm readiness:production` falla deliberadamente mientras exista cualquier claim `OPEN` o `CONDITIONED` en `config/readiness.json`; Codex no puede convertir esos estados en promesas.

## Qué se configura sin programar

- identidad visual básica del negocio y soporte;
- mercados, monedas, locales, zonas horarias y modo fiscal;
- tipos y jerarquías de organizaciones;
- roles y permisos;
- módulos y feature flags;
- estados y transiciones de workflows;
- campos personalizados;
- proveedores de integración, modo, capacidades y referencias de credenciales.

La fuente es `config/business.example.json`; el contrato ejecutable es `src/platform/config/schema.ts`. Para otra variante, crear otro archivo dentro de `config/` y seleccionar sólo su nombre con `BUSINESS_CONFIG_FILE`.

## Qué ya atraviesa el sistema

```text
web pública accesible
→ API con Problem Details e idempotencia
→ validación de mercado/campos
→ transacción PostgreSQL
→ lead + auditoría + outbox
→ evento despachable
→ logs operativos estructurados
```

Además existen schemas para organizaciones, membresías, catálogo, stock y pedidos; autorización por organización/objeto; workflows configurables; health/readiness; portales público, cliente, operación y fábrica; límites de abuso; seguridad de headers; y contrato fail-closed para adaptadores externos.

## Frontera honesta

Configurar no puede reemplazar toda lógica de dominio. Impuestos, facturación fiscal, pricing complejo, financiación, homologación, garantías, telemetría, contratos con fábricas y semántica concreta de proveedores requieren módulos o adaptadores con pruebas. OIDC, Mercado Pago, Mercado Libre, Amazon SP-API, Google Ads y Meta Ads están declarados pero deliberadamente no se activan sin implementación, credenciales y contract tests.

Los portales internos están cerrados y no muestran datos hasta integrar identidad real. PGlite es para desarrollo y tests; producción exige PostgreSQL administrado con backup/restore demostrado.

## Arquitectura

- modular monolith y un despliegue observable como baseline;
- PostgreSQL como fuente de verdad;
- contratos entre configuración, dominio, persistencia e integraciones;
- outbox transaccional para side effects;
- autorización separada de autenticación;
- fallos de configuración o seguridad bloquean el arranque productivo.

Las decisiones y restricciones se gobiernan desde los Markdown raíz de esta biblioteca. Leer `CODEX_START_HERE.md` antes de encargar un proyecto a Codex.
````

### FILE: `scripts/check-readiness.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:scripts-check-readiness-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "dd5b4f6cc7b6508e95c719f189abf1a6f7a26baa67a7407bcd99757095c16ad7"
variables: []
secrets_allowed: false
```

````typescript
import { readFile } from "node:fs/promises";
import { join } from "node:path";
import { z } from "zod";

const status = z.enum(["READY", "CONDITIONED", "OPEN"]);
const schema = z.object({
  schemaVersion: z.literal("1.0.0"),
  claims: z.record(z.string(), status),
  targets: z.record(z.string(), z.array(z.string()))
});

const target = process.argv[2] ?? "baseline";
const readiness = schema.parse(JSON.parse(await readFile(join(process.cwd(), "config", "readiness.json"), "utf8")));
const required = readiness.targets[target];
if (!required) throw new Error(`unknown readiness target: ${target}`);
const unresolved = required.filter((claim) => readiness.claims[claim] !== "READY").map((claim) => ({ claim, status: readiness.claims[claim] ?? "MISSING" }));
console.log(JSON.stringify({ target, status: unresolved.length === 0 ? "READY" : "NOT_READY", unresolved }, null, 2));
if (unresolved.length > 0) process.exitCode = 1;
````

### FILE: `scripts/dispatch-outbox.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:scripts-dispatch-outbox-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "f8c364061339926c4c87da032f53114a9db5bc6b787cbe5e7c0299a959f6f99b"
variables: []
secrets_allowed: false
```

````typescript
import { database, closeDatabase } from "../src/platform/db/client";
import { ensureBootstrapped } from "../src/platform/bootstrap";
import { dispatchOutboxBatch } from "../src/platform/outbox/dispatcher";

await ensureBootstrapped();
const count = await dispatchOutboxBatch(await database(), {
  "lead.created": async (event) => console.info(JSON.stringify({ event: event.event_type, eventId: event.event_id, payload: event.payload }))
});
console.info(`Dispatched ${count} outbox event(s).`);
await closeDatabase();
````

### FILE: `scripts/migrate.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:scripts-migrate-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "ed08abe274081b645d06d6e797afd8d201097734af8c8902bbd51aec4086731b"
variables: []
secrets_allowed: false
```

````typescript
import { createDatabase } from "../src/platform/db/client";
import { migrate } from "../src/platform/db/migrations";

const database = await createDatabase();
try {
  const applied = await migrate(database);
  console.log(JSON.stringify({ status: "migrated", applied }, null, 2));
} finally {
  await database.close();
}
````

### FILE: `scripts/seed.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:scripts-seed-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "1f58f394b5d56273509748aa045bf1eef7488903c837107d5f8a7a31aadcdce4"
variables: []
secrets_allowed: false
```

````typescript
import { loadBusinessConfig } from "../src/platform/config/load";
import { createDatabase } from "../src/platform/db/client";
import { migrate } from "../src/platform/db/migrations";
import { seed } from "../src/platform/seed";

const database = await createDatabase();
try {
  await migrate(database);
  await seed(database, await loadBusinessConfig());
  console.log(JSON.stringify({ status: "seeded" }, null, 2));
} finally {
  await database.close();
}
````

### FILE: `scripts/verify-config.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:scripts-verify-config-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "a485cdcc63d725a58ae374ae73670bec90578a1b3ea8bfa2462c1d122af584a7"
variables: []
secrets_allowed: false
```

````typescript
import { loadBusinessConfig } from "../src/platform/config/load";
import { enabledModules } from "../src/platform/config/registry";
import { validateRuntimePolicy } from "../src/platform/config/runtime-policy";
import { validateEnabledIntegrations } from "../src/platform/integrations/contract";
import { installedIntegrationAdapters } from "../src/platform/integrations/registry";

const config = await loadBusinessConfig({ bypassCache: true });
validateRuntimePolicy(config);
validateEnabledIntegrations(config, installedIntegrationAdapters);
console.log(JSON.stringify({
  status: "valid",
  schemaVersion: config.schemaVersion,
  business: config.business.id,
  markets: config.markets.map((market) => market.code),
  modules: enabledModules(config),
  enabledIntegrations: config.integrations.filter((item) => item.enabled).map((item) => item.id)
}, null, 2));
````

### FILE: `src/app/admin/page.tsx`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-admin-page-tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "561b66589954c71c88b4f663c5342a898bd6e0f8fdbf82e7bc0d4ba96cd26428"
variables: []
secrets_allowed: false
```

````tsx
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type OrderSummary, type Overview, type Page, type ServiceCaseSummary } from "@/platform/backend/protected-client";

export default async function AdminPage() {
  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/admin" as Route);
  if (!allowed(session, "admin:read")) return <><h1 className="pageTitle">Acceso denegado</h1><div className="notice">La sesión no posee <code>admin:read</code>.</div></>;
  const organization = session.organizations[0]!;
  const [overview, orders, cases] = await Promise.all([
    protectedGet<Overview>(session, "/v1/admin/overview", { organization_id: organization }),
    protectedGet<Page<OrderSummary>>(session, "/v1/admin/orders", { organization_id: organization, limit: "25" }),
    protectedGet<Page<ServiceCaseSummary>>(session, "/v1/admin/service-cases", { organization_id: organization, limit: "25" }),
  ]);
  return <>
    <div className="eyebrow">Portal operativo · {organization}</div><h1 className="pageTitle">Operación</h1>
    <div className="grid">
      <article className="card"><h2>Pedidos</h2><p>{overview.orders}</p></article>
      <article className="card"><h2>Stock disponible</h2><p>{overview.stock_available}</p></article>
      <article className="card"><h2>Casos abiertos</h2><p>{overview.open_cases}</p></article>
      <article className="card"><h2>Envíos activos</h2><p>{overview.active_shipments}</p></article>
    </div>
    <section><h2>Pedidos recientes</h2><div className="grid">{orders.items.map((order) => <article className="card" key={order.id}><h3>{order.id}</h3><p>{order.state} · {order.currency} {order.total_minor_units}</p></article>)}</div></section>
    <section><h2>Servicio</h2><div className="grid">{cases.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} · {item.severity}</p></article>)}</div></section>
  </>;
}
````

### FILE: `src/app/api/health/route.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-api-health-route-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "68ccba9169957a11d3a625c2ec9f572b728e2ef74e650f022e72d0152ae84cd8"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { loadBusinessConfig } from "@/platform/config/load";
import { database } from "@/platform/db/client";
import { migrate } from "@/platform/db/migrations";
import { errorResponse } from "@/platform/http/problem";
import { validateRuntimePolicy } from "@/platform/config/runtime-policy";

export const dynamic = "force-dynamic";

export async function GET() {
  const started = performance.now();
  try {
    const [config, db] = await Promise.all([loadBusinessConfig(), database()]);
    validateRuntimePolicy(config);
    await migrate(db);
    await db.query("select 1 as ready");
    return NextResponse.json({ status: "ready", configSchema: config.schemaVersion, businessId: config.business.id, checks: { configuration: "pass", database: "pass" }, durationMs: Math.round(performance.now() - started) }, { headers: { "cache-control": "no-store" } });
  } catch (error) {
    return errorResponse(error);
  }
}
````

### FILE: `src/app/api/leads/route.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-api-leads-route-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "5028b764244965e1b86d690898d7b85bc86c8bf272311837d23646efb83b176a"
variables: []
secrets_allowed: false
```

````typescript
import { loadBusinessConfig } from "@/platform/config/load";
import { database } from "@/platform/db/client";
import { ensureBootstrapped } from "@/platform/bootstrap";
import { errorResponse, problem } from "@/platform/http/problem";
import { emitOperationalEvent } from "@/platform/observability/events";
import { LeadService } from "@/modules/crm/lead-service";

export async function POST(request: Request) {
  const started = performance.now();
  const idempotencyKey = request.headers.get("idempotency-key");
  if (!idempotencyKey) return problem(400, "Falta idempotencia", "Envíe el encabezado Idempotency-Key.", "IDEMPOTENCY_KEY_REQUIRED");
  try {
    const config = await loadBusinessConfig();
    if (!config.modules.crm?.enabled || !config.features.lead_capture) return problem(404, "Función desactivada", "La captura de interesados no está habilitada.", "FEATURE_DISABLED");
    await ensureBootstrapped();
    const forwardedFor = request.headers.get("x-forwarded-for")?.split(",")[0]?.trim() ?? "unknown";
    const fingerprint = `${forwardedFor}:${request.headers.get("user-agent") ?? "unknown"}`;
    const result = await new LeadService(await database(), config).create(await request.json(), idempotencyKey, fingerprint);
    emitOperationalEvent({ event: "lead.create", outcome: "success", durationMs: Math.round(performance.now() - started), details: { replayed: result.replayed } });
    return Response.json(result, { status: result.replayed ? 200 : 201, headers: { "cache-control": "no-store" } });
  } catch (error) {
    emitOperationalEvent({ event: "lead.create", outcome: "failure", durationMs: Math.round(performance.now() - started) });
    return errorResponse(error);
  }
}
````

### FILE: `src/app/api/platform/route.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-api-platform-route-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "b8981a28235fe8e6b74b21dff6461ffa7619b355f21ccf3025c01f54d116fd7c"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { loadBusinessConfig } from "@/platform/config/load";
import { enabledModules } from "@/platform/config/registry";
import { errorResponse } from "@/platform/http/problem";
import { validateRuntimePolicy } from "@/platform/config/runtime-policy";

export async function GET() {
  try {
    const config = await loadBusinessConfig();
    validateRuntimePolicy(config);
    return NextResponse.json({
      schemaVersion: config.schemaVersion,
      business: { id: config.business.id, name: config.business.name, defaultLocale: config.business.defaultLocale, defaultMarket: config.business.defaultMarket },
      markets: config.markets.map(({ code, name, currency, locales }) => ({ code, name, currency, locales })),
      modules: enabledModules(config),
      features: config.features
    }, { headers: { "cache-control": "public, max-age=60, stale-while-revalidate=300" } });
  } catch (error) {
    return errorResponse(error);
  }
}
````

### FILE: `src/app/customer/page.tsx`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-customer-page-tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "a89873ded6b5522c5a0cb74016c552e8b36229318ff6b7820a17fb8939b3003d"
variables: []
secrets_allowed: false
```

````tsx
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type OrderSummary, type Page, type ServiceCaseSummary } from "@/platform/backend/protected-client";

export default async function CustomerPage() {
  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">Acceso denegado</h1><div className="notice">La sesión no posee <code>customer:self</code>.</div></>;
  const organization = session.organizations[0]!;
  const [orders, cases] = await Promise.all([
    protectedGet<Page<OrderSummary>>(session, "/v1/customer/orders", { organization_id: organization, limit: "25" }),
    protectedGet<Page<ServiceCaseSummary>>(session, "/v1/customer/service-cases", { organization_id: organization, limit: "25" }),
  ]);
  return <><div className="eyebrow">Portal de cliente</div><h1 className="pageTitle">Mi cuenta</h1>
    <section><h2>Mis pedidos</h2><div className="grid">{orders.items.map((order) => <article className="card" key={order.id}><h3>{order.id}</h3><p>{order.state} · {order.currency} {order.total_minor_units}</p></article>)}</div></section>
    <section><h2>Mis casos de servicio</h2><div className="grid">{cases.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} · {item.description}</p></article>)}</div></section>
  </>;
}
````

### FILE: `src/app/factory/page.tsx`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-factory-page-tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "94339c9261711ff7824f6c07631a3550f82c216d790060fc07c9db8e749cefd2"
variables: []
secrets_allowed: false
```

````tsx
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type FactoryUnitSummary, type Page } from "@/platform/backend/protected-client";

export default async function FactoryPage() {
  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/factory" as Route);
  if (!allowed(session, "factory:read")) return <><h1 className="pageTitle">Acceso denegado</h1><div className="notice">La sesión no posee <code>factory:read</code>.</div></>;
  const organization = session.organizations[0]!;
  const units = await protectedGet<Page<FactoryUnitSummary>>(session, "/v1/factory/units", { organization_id: organization, limit: "25" });
  return <><div className="eyebrow">Portal de fábrica · {organization}</div><h1 className="pageTitle">Abastecimiento y seguimiento</h1>
    <div className="grid">{units.items.map((unit) => <article className="card" key={unit.id}><h2>{unit.serial_number}</h2><p>{unit.state} · orden {unit.purchase_order_id}</p></article>)}</div>
  </>;
}
````

### FILE: `src/app/globals.css`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-globals-css:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "36ec4750b1912f516622166df632bc4a6109d2715b9b6822b2539c23eac1e516"
variables: []
secrets_allowed: false
```

````css
:root {
  color-scheme: light;
  --ink: #12221c;
  --muted: #5d6f67;
  --surface: #f4f7f5;
  --panel: #ffffff;
  --line: #d9e2dd;
  --accent: #006b4f;
  --accent-strong: #004c39;
  --warning: #7a4d00;
  font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

* { box-sizing: border-box; }
body { margin: 0; color: var(--ink); background: var(--surface); }
a { color: inherit; }
header { background: var(--panel); border-bottom: 1px solid var(--line); }
.shell { width: min(1120px, calc(100% - 2rem)); margin: 0 auto; }
.headerRow { min-height: 68px; display: flex; align-items: center; justify-content: space-between; gap: 1.5rem; }
.brand { font-weight: 800; text-decoration: none; letter-spacing: -0.03em; }
nav { display: flex; flex-wrap: wrap; gap: 1rem; }
nav a { color: var(--muted); text-decoration: none; font-weight: 650; }
nav a:hover, nav a:focus-visible { color: var(--accent); text-decoration: underline; }
main { padding: 3rem 0 5rem; }
.hero { display: grid; grid-template-columns: minmax(0, 1.5fr) minmax(280px, .8fr); gap: 2rem; align-items: center; }
.eyebrow { color: var(--accent); font-weight: 800; text-transform: uppercase; font-size: .78rem; letter-spacing: .12em; }
h1 { font-size: clamp(2.4rem, 7vw, 5.5rem); line-height: .95; letter-spacing: -.06em; margin: .6rem 0 1.2rem; }
.pageTitle { font-size: clamp(2.2rem, 5vw, 4rem); }
h2 { font-size: clamp(1.5rem, 3vw, 2.4rem); letter-spacing: -.035em; }
p { line-height: 1.65; }
.lede { font-size: 1.18rem; color: var(--muted); max-width: 68ch; }
.panel, .card { background: var(--panel); border: 1px solid var(--line); border-radius: 18px; padding: 1.4rem; box-shadow: 0 14px 35px rgba(18, 34, 28, .06); }
.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1rem; margin-top: 1.5rem; }
.metric { font-size: 2rem; font-weight: 850; color: var(--accent); }
.muted { color: var(--muted); }
.button { display: inline-flex; align-items: center; justify-content: center; padding: .8rem 1rem; border-radius: 10px; border: 0; background: var(--accent); color: white; font-weight: 750; text-decoration: none; cursor: pointer; }
.button:hover, .button:focus-visible { background: var(--accent-strong); }
label { display: grid; gap: .4rem; font-weight: 700; }
input, select { width: 100%; border: 1px solid #aebdb5; border-radius: 9px; padding: .75rem; font: inherit; background: white; }
input:focus, select:focus { outline: 3px solid rgba(0,107,79,.25); border-color: var(--accent); }
form { display: grid; gap: 1rem; }
.status { padding: .8rem; border-radius: 8px; background: #e7f6ef; color: var(--accent-strong); }
.notice { margin: 1rem 0 2rem; padding: 1rem; border: 1px solid #e3c98e; border-radius: 10px; background: #fff8e8; color: var(--warning); line-height: 1.6; }
.error { background: #fff0ed; color: #8a2717; }
footer { padding: 2rem 0; border-top: 1px solid var(--line); color: var(--muted); }
code { background: #e7ece9; border-radius: 5px; padding: .12rem .3rem; }
@media (max-width: 760px) { .hero { grid-template-columns: 1fr; } .headerRow { align-items: flex-start; flex-direction: column; padding: 1rem 0; } }
@media (prefers-reduced-motion: reduce) { *, *::before, *::after { scroll-behavior: auto !important; transition: none !important; } }
````

### FILE: `src/app/layout.tsx`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-layout-tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "d404abbac5a434c0ad185b14f965a5083e8a7cdd55e3a126df12d1a3f998f935"
variables: []
secrets_allowed: false
```

````tsx
import type { Metadata } from "next";
import Link from "next/link";
import "./globals.css";
import { loadBusinessConfig } from "@/platform/config/load";
import { enabledNavigation } from "@/platform/config/registry";

export async function generateMetadata(): Promise<Metadata> {
  const config = await loadBusinessConfig();
  return { title: { default: config.business.name, template: `%s · ${config.business.name}` }, description: "Plataforma empresarial configurable." };
}

export default async function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  const config = await loadBusinessConfig();
  const navigation = enabledNavigation(config);
  return (
    <html lang={config.business.defaultLocale}>
      <body>
        <header>
          <div className="shell headerRow">
            <Link className="brand" href="/">{config.business.name}</Link>
            <nav aria-label="Principal">
              {navigation.map((item) => <Link key={item.id} href={item.href as never}>{item.label}</Link>)}
            </nav>
          </div>
        </header>
        <main><div className="shell">{children}</div></main>
        <footer><div className="shell">Configuración <code>{config.schemaVersion}</code> · {config.business.supportEmail}</div></footer>
      </body>
    </html>
  );
}
````

### FILE: `src/app/models/[slug]/lead-form.tsx`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-models-slug-lead-form-tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "68614e6839e80087446483fb0c881a3e8fbae717d9007cae20f36e5d7b4172bc"
variables: []
secrets_allowed: false
```

````tsx
"use client";

import { useState, type FormEvent } from "react";

export function LeadForm({ modelId, marketCode }: { modelId: string; marketCode: string }) {
  const [status, setStatus] = useState<"idle" | "sending" | "success" | "error">("idle");
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setStatus("sending");
    const form = new FormData(event.currentTarget);
    const response = await fetch("/api/leads", {
      method: "POST",
      headers: { "content-type": "application/json", "idempotency-key": crypto.randomUUID() + crypto.randomUUID() },
      body: JSON.stringify({ name: form.get("name"), email: form.get("email"), modelId, marketCode, consentVersion: "lead-v1", source: "public-model", customFields: { preferred_vehicle_use: form.get("preferred_vehicle_use") } })
    });
    setStatus(response.ok ? "success" : "error");
    if (response.ok) event.currentTarget.reset();
  }
  return (
    <form onSubmit={submit}>
      <label>Nombre<input name="name" autoComplete="name" required minLength={2} /></label>
      <label>Email<input name="email" type="email" autoComplete="email" required /></label>
      <label>Uso principal<select name="preferred_vehicle_use" defaultValue="urban"><option value="urban">Urbano</option><option value="delivery">Reparto</option><option value="recreation">Recreación</option><option value="fleet">Flota</option></select></label>
      <label><span><input name="consent" type="checkbox" required style={{ width: "auto" }} /> Acepto ser contactado según el consentimiento lead-v1.</span></label>
      <button className="button" disabled={status === "sending"}>{status === "sending" ? "Enviando…" : "Solicitar contacto"}</button>
      {status === "success" && <div className="status" role="status">Solicitud registrada.</div>}
      {status === "error" && <div className="status error" role="alert">No pudimos registrarla. Intentá nuevamente.</div>}
    </form>
  );
}
````

### FILE: `src/app/models/[slug]/page.tsx`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-models-slug-page-tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "52fe7cfbaef31377ede8554c5558981a98b3ed17156f07f8d9e710e0d530d56b"
variables: []
secrets_allowed: false
```

````tsx
import { notFound } from "next/navigation";
import { database } from "@/platform/db/client";
import { migrate } from "@/platform/db/migrations";
import { seed } from "@/platform/seed";
import { loadBusinessConfig } from "@/platform/config/load";
import { findPublishedModel } from "@/modules/catalog/repository";
import { LeadForm } from "./lead-form";

export const dynamic = "force-dynamic";

export default async function ModelPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const config = await loadBusinessConfig();
  const db = await database();
  await migrate(db);
  await seed(db, config);
  const model = await findPublishedModel(db, slug, config.business.defaultMarket);
  if (!model) notFound();
  return (
    <div className="hero">
      <section>
        <div className="eyebrow">{model.marketCode} · catálogo publicado</div>
        <h1>{model.name}</h1>
        <p className="lede">{model.summary}</p>
        <div className="panel"><strong>Autonomía estimada</strong><div className="metric">{String(model.customFields.estimated_range_km ?? "—")} km</div></div>
      </section>
      <aside className="panel"><h2>Quiero conocerlo</h2><LeadForm modelId={model.modelId} marketCode={model.marketCode} /></aside>
    </div>
  );
}
````

### FILE: `src/app/models/page.tsx`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-models-page-tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "c74e7df21af3f3fd6bbf09a652062cc529473692bd21d3d75bce59dbed422588"
variables: []
secrets_allowed: false
```

````tsx
import type { Metadata } from "next";
import Link from "next/link";
import { database } from "@/platform/db/client";
import { migrate } from "@/platform/db/migrations";
import { seed } from "@/platform/seed";
import { loadBusinessConfig } from "@/platform/config/load";
import { listPublishedModels } from "@/modules/catalog/repository";

export const metadata: Metadata = { title: "Modelos" };
export const dynamic = "force-dynamic";

export default async function ModelsPage() {
  const config = await loadBusinessConfig();
  const db = await database();
  await migrate(db);
  await seed(db, config);
  const models = await listPublishedModels(db, config.business.defaultMarket);
  return (
    <>
      <div className="eyebrow">Catálogo público</div>
      <h1>Modelos</h1>
      <p className="lede">El contenido proviene de una versión publicada del catálogo y se filtra por mercado.</p>
      <div className="grid">
        {models.map((model) => (
          <article className="card" key={model.modelId}>
            <h2>{model.name}</h2>
            <p>{model.summary}</p>
            <p className="muted">Autonomía estimada: {String(model.customFields.estimated_range_km ?? "—")} km</p>
            <Link className="button" href={`/models/${model.slug}` as never}>Ver modelo</Link>
          </article>
        ))}
      </div>
    </>
  );
}
````

### FILE: `src/app/page.tsx`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-app-page-tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "9bfd77c3cc5e789987b24ff92ed7577055ef066a19c7cdcc5d7aff5b9fcd56c1"
variables: []
secrets_allowed: false
```

````tsx
import Link from "next/link";
import { loadBusinessConfig } from "@/platform/config/load";
import { enabledModules } from "@/platform/config/registry";

export default async function HomePage() {
  const config = await loadBusinessConfig();
  const modules = enabledModules(config);
  return (
    <>
      <section className="hero">
        <div>
          <div className="eyebrow">Adapter web configurable</div>
          <h1>Un núcleo. Muchos negocios.</h1>
          <p className="lede">Mercados, organizaciones, roles, permisos, módulos, workflows, campos e integraciones se validan desde configuración antes de iniciar el sistema.</p>
          <Link className="button" href="/models">Explorar modelos</Link>
        </div>
        <aside className="panel" aria-label="Estado de configuración">
          <div className="metric">{modules.length}</div>
          <p>Módulos activos</p>
          <p className="muted">{modules.join(" · ")}</p>
          <div className="metric">{config.markets.length}</div>
          <p>Mercados configurados</p>
        </aside>
      </section>
      <section>
        <h2>Capacidades configurables</h2>
        <div className="grid">
          <article className="card"><h3>Organización</h3><p>Casa central, franquicias, sucursales, fábricas y proveedores con jerarquías validadas.</p></article>
          <article className="card"><h3>Operación</h3><p>Catálogo, leads, compras, inventario, pedidos, pagos, fulfillment y service activables.</p></article>
          <article className="card"><h3>Control</h3><p>Permisos por organización, workflows explícitos, idempotencia, outbox y auditoría.</p></article>
        </div>
      </section>
    </>
  );
}
````

### FILE: `src/instrumentation.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-instrumentation-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "4e6c3a3f500d2fdc48596466b772cebd26c3ba9aa6c6b891b41522c15f676197"
variables: []
secrets_allowed: false
```

````typescript
export async function register(): Promise<void> {
  if (process.env.NEXT_RUNTIME !== "nodejs") return;
  const [{ loadBusinessConfig }, { validateRuntimePolicy }, { validateEnabledIntegrations }, { installedIntegrationAdapters }] = await Promise.all([
    import("@/platform/config/load"),
    import("@/platform/config/runtime-policy"),
    import("@/platform/integrations/contract"),
    import("@/platform/integrations/registry")
  ]);
  const config = await loadBusinessConfig();
  validateRuntimePolicy(config);
  validateEnabledIntegrations(config, installedIntegrationAdapters);
}
````

### FILE: `src/modules/catalog/repository.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-modules-catalog-repository-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "f12094b5c07df97410421ee448051c7463bcb23d0ea03f7fedcad5fbd2fa1f9c"
variables: []
secrets_allowed: false
```

````typescript
import type { SqlExecutor } from "@/platform/db/types";

export interface CatalogModel {
  modelId: string;
  slug: string;
  name: string;
  summary: string;
  marketCode: string;
  customFields: Record<string, unknown>;
}

interface CatalogModelRow extends Record<string, unknown> {
  model_id: string;
  slug: string;
  name: string;
  summary: string;
  market_code: string;
  custom_fields: Record<string, unknown>;
}

export async function listPublishedModels(database: SqlExecutor, marketCode: string): Promise<CatalogModel[]> {
  const result = await database.query<CatalogModelRow>(`
    select model_id, slug, name, summary, market_code, custom_fields
    from catalog.model
    where status = 'published' and market_code = $1
    order by name
  `, [marketCode]);
  return result.rows.map(mapModel);
}

export async function findPublishedModel(database: SqlExecutor, slug: string, marketCode: string): Promise<CatalogModel | undefined> {
  const result = await database.query<CatalogModelRow>(`
    select model_id, slug, name, summary, market_code, custom_fields
    from catalog.model
    where status = 'published' and slug = $1 and market_code = $2
  `, [slug, marketCode]);
  const row = result.rows[0];
  return row ? mapModel(row) : undefined;
}

function mapModel(row: CatalogModelRow): CatalogModel {
  return {
    modelId: row.model_id,
    slug: row.slug,
    name: row.name,
    summary: row.summary,
    marketCode: row.market_code,
    customFields: row.custom_fields
  };
}
````

### FILE: `src/modules/crm/lead-service.test.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-modules-crm-lead-service-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "2ac189797c6a5043b63fe3d41741da870a402d102a6bebd24e800746cf5c19c8"
variables: []
secrets_allowed: false
```

````typescript
import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { createDatabase } from "@/platform/db/client";
import type { Database } from "@/platform/db/types";
import { migrate } from "@/platform/db/migrations";
import { loadBusinessConfig } from "@/platform/config/load";
import { LeadService, IdempotencyConflictError, InvalidLeadInputError, RateLimitExceededError } from "./lead-service";

let db: Database;
let service: LeadService;

beforeEach(async () => {
  db = await createDatabase("memory://");
  await migrate(db);
  service = new LeadService(db, await loadBusinessConfig({ bypassCache: true }));
});

afterEach(async () => db.close());

describe("lead creation transaction", () => {
  const request = { name: "Ana Cliente", email: "ANA@example.com", marketCode: "AR", consentVersion: "2026-08", source: "public_web", customFields: { preferred_vehicle_use: "urban" } };

  it("writes the lead, audit and outbox atomically", async () => {
    const result = await service.create(request, "lead:test:00000001");
    expect(result.replayed).toBe(false);
    expect((await db.query("select * from crm.lead")).rows).toHaveLength(1);
    expect((await db.query("select * from platform.audit_event")).rows).toHaveLength(1);
    expect((await db.query("select * from platform.outbox_event")).rows).toHaveLength(1);
  });

  it("replays the original response for the same request", async () => {
    const first = await service.create(request, "lead:test:00000002");
    const replay = await service.create(request, "lead:test:00000002");
    expect(replay).toEqual({ leadId: first.leadId, state: first.state, replayed: true });
    expect((await db.query("select * from crm.lead")).rows).toHaveLength(1);
  });

  it("rejects reuse of a key with a different request", async () => {
    await service.create(request, "lead:test:00000003");
    await expect(service.create({ ...request, email: "otra@example.com" }, "lead:test:00000003")).rejects.toBeInstanceOf(IdempotencyConflictError);
  });

  it("limits abusive unique submissions by privacy-preserving fingerprint", async () => {
    for (let index = 0; index < 10; index += 1) {
      await service.create({ ...request, email: `limit-${index}@example.com` }, `lead:limit:${String(index).padStart(8, "0")}`, "same-client");
    }
    await expect(service.create({ ...request, email: "limit-final@example.com" }, "lead:limit:99999999", "same-client")).rejects.toBeInstanceOf(RateLimitExceededError);
  });

  it("classifies unknown custom fields as invalid client input", async () => {
    await expect(service.create({ ...request, customFields: { unknown_field: "value" } }, "lead:test:00000004")).rejects.toBeInstanceOf(InvalidLeadInputError);
  });
});
````

### FILE: `src/modules/crm/lead-service.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-modules-crm-lead-service-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "8de72348e485b8c01546c2b944593e217ca094f856241f8a24bc1863ac3c0be6"
variables: []
secrets_allowed: false
```

````typescript
import { createHash, randomUUID } from "node:crypto";
import { z } from "zod";
import type { BusinessConfig } from "@/platform/config/schema";
import type { Database } from "@/platform/db/types";

const leadInputSchema = z.object({
  name: z.string().trim().min(2).max(120),
  email: z.email().max(254).transform((value) => value.toLowerCase()),
  marketCode: z.string().length(2),
  modelId: z.uuid().optional(),
  consentVersion: z.string().min(1).max(50),
  source: z.string().min(1).max(80),
  customFields: z.record(z.string(), z.unknown()).default({})
});

export type LeadInput = z.input<typeof leadInputSchema>;

export interface LeadCreated {
  leadId: string;
  state: string;
  replayed: boolean;
}

export class IdempotencyConflictError extends Error {
  readonly code = "IDEMPOTENCY_CONFLICT";
}

export class RateLimitExceededError extends Error {
  readonly code = "RATE_LIMIT_EXCEEDED";
  readonly retryAfterSeconds = 3600;
}

export class InvalidLeadInputError extends Error {
  readonly code = "INVALID_LEAD_INPUT";
}

export class LeadService {
  constructor(private readonly database: Database, private readonly config: BusinessConfig) {}

  async create(input: LeadInput, idempotencyKey: string, clientFingerprint = "unspecified"): Promise<LeadCreated> {
    if (!/^[A-Za-z0-9_.:-]{16,128}$/.test(idempotencyKey)) throw new InvalidLeadInputError("invalid idempotency key");
    const parsed = leadInputSchema.parse(input);
    const market = this.config.markets.find((item) => item.code === parsed.marketCode);
    if (!market) throw new InvalidLeadInputError("unsupported market");
    const customFields = validateCustomFields(this.config, "lead", parsed.customFields);
    const requestHash = stableHash({ ...parsed, customFields });

    return this.database.transaction(async (transaction) => {
      const inserted = await transaction.query(`
        insert into platform.idempotency_record
          (scope, idempotency_key, request_hash, status, expires_at)
        values ('lead:create', $1, $2, 'processing', now() + interval '24 hours')
        on conflict do nothing
        returning idempotency_key
      `, [idempotencyKey, requestHash]);

      if (inserted.rows.length === 0) {
        const existing = await transaction.query<{
          request_hash: string;
          status: string;
          response_body: { leadId: string; state: string } | null;
        } & Record<string, unknown>>(`
          select request_hash, status, response_body
          from platform.idempotency_record
          where scope = 'lead:create' and idempotency_key = $1
        `, [idempotencyKey]);
        const record = existing.rows[0];
        if (!record || record.request_hash !== requestHash) throw new IdempotencyConflictError("idempotency key reused with a different request");
        if (record.status !== "completed" || !record.response_body) throw new Error("request is still processing");
        return { ...record.response_body, replayed: true };
      }

      const rate = await transaction.query<{ request_count: number } & Record<string, unknown>>(`
        insert into platform.rate_limit_window (scope, key_hash, window_started, request_count)
        values ('lead:create', $1, date_trunc('hour', now()), 1)
        on conflict (scope, key_hash) do update
        set request_count = case
              when platform.rate_limit_window.window_started = date_trunc('hour', now()) then platform.rate_limit_window.request_count + 1
              else 1
            end,
            window_started = date_trunc('hour', now())
        returning request_count
      `, [stableHash(clientFingerprint)]);
      if ((rate.rows[0]?.request_count ?? 0) > 10) throw new RateLimitExceededError("too many lead requests in the current window");

      const leadId = randomUUID();
      const state = this.config.workflows.lead?.initial;
      if (!state) throw new Error("lead workflow is not configured");
      await transaction.query(`
        insert into crm.lead
          (lead_id, market_code, email, name, model_id, consent_version, source, state, custom_fields)
        values ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb)
      `, [leadId, parsed.marketCode, parsed.email, parsed.name, parsed.modelId ?? null, parsed.consentVersion, parsed.source, state, JSON.stringify(customFields)]);

      await transaction.query(`
        insert into platform.outbox_event
          (event_id, aggregate_type, aggregate_id, aggregate_version, event_type, schema_version, payload)
        values ($1, 'lead', $2, 0, 'lead.created', 1, $3::jsonb)
      `, [randomUUID(), leadId, JSON.stringify({ leadId, marketCode: parsed.marketCode, source: parsed.source })]);

      await transaction.query(`
        insert into platform.audit_event
          (audit_event_id, actor_id, action, resource_type, resource_id, outcome, details)
        values ($1, 'anonymous', 'lead.create', 'lead', $2, 'success', $3::jsonb)
      `, [randomUUID(), leadId, JSON.stringify({ marketCode: parsed.marketCode, source: parsed.source })]);

      const response = { leadId, state };
      await transaction.query(`
        update platform.idempotency_record
        set status = 'completed', response_code = 201, response_body = $2::jsonb, resource_id = $3
        where scope = 'lead:create' and idempotency_key = $1
      `, [idempotencyKey, JSON.stringify(response), leadId]);
      return { ...response, replayed: false };
    });
  }
}

function validateCustomFields(config: BusinessConfig, entity: string, values: Record<string, unknown>): Record<string, unknown> {
  const definitions = config.customFields[entity] ?? [];
  const definitionMap = new Map(definitions.map((definition) => [definition.id, definition]));
  for (const key of Object.keys(values)) {
    if (!definitionMap.has(key)) throw new InvalidLeadInputError(`unknown custom field: ${key}`);
  }
  for (const definition of definitions) {
    const value = values[definition.id];
    if (definition.required && (value === undefined || value === null || value === "")) throw new InvalidLeadInputError(`required custom field missing: ${definition.id}`);
    if (value === undefined || value === null) continue;
    if (definition.type === "number" && typeof value !== "number") throw new InvalidLeadInputError(`custom field ${definition.id} must be number`);
    if (definition.type === "boolean" && typeof value !== "boolean") throw new InvalidLeadInputError(`custom field ${definition.id} must be boolean`);
    if ((definition.type === "string" || definition.type === "date" || definition.type === "select") && typeof value !== "string") throw new InvalidLeadInputError(`custom field ${definition.id} must be string`);
    if (definition.type === "select" && !definition.options.includes(value as string)) throw new InvalidLeadInputError(`custom field ${definition.id} has invalid option`);
  }
  return values;
}

function stableHash(value: unknown): string {
  return createHash("sha256").update(JSON.stringify(sortValue(value))).digest("hex");
}

function sortValue(value: unknown): unknown {
  if (Array.isArray(value)) return value.map(sortValue);
  if (value && typeof value === "object") {
    return Object.fromEntries(Object.entries(value).sort(([left], [right]) => left.localeCompare(right)).map(([key, item]) => [key, sortValue(item)]));
  }
  return value;
}
````

### FILE: `src/modules/enterprise/operations.test.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-modules-enterprise-operations-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "e6a8563c083bed77d6e8d8a20b071b498d952ecc4b90f5949bb3795c383605b3"
variables: []
secrets_allowed: false
```

````typescript
import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { createDatabase } from "@/platform/db/client";
import type { Database } from "@/platform/db/types";
import { migrate } from "@/platform/db/migrations";
import { seed } from "@/platform/seed";
import { loadBusinessConfig } from "@/platform/config/load";
import { EnterpriseConflictError, EnterpriseOperations } from "./operations";

let database: Database;
let operations: EnterpriseOperations;
let organizationId: string;
let modelId: string;

beforeEach(async () => {
  database = await createDatabase("memory://");
  await migrate(database);
  const config = await loadBusinessConfig({ bypassCache: true });
  await seed(database, config);
  operations = new EnterpriseOperations(database, config);
  organizationId = String((await database.query("select organization_id from org.organization limit 1")).rows[0]?.organization_id);
  modelId = String((await database.query("select model_id from catalog.model limit 1")).rows[0]?.model_id);
});

afterEach(async () => database.close());

describe("enterprise operational vertical", () => {
  it("tracks supplier, purchase order, shipment and serialized inventory atomically", async () => {
    const supplierId = await operations.createSupplier(organizationId, "factory-one", "Factory One Ltd.");
    const purchaseOrderId = await operations.createPurchaseOrder({
      organizationId,
      supplierId,
      currency: "USD",
      lines: [{ modelId, quantity: 2, unitMinorUnits: 100_000 }]
    });
    const receipt = await operations.receiveSerializedStock({
      organizationId,
      purchaseOrderId,
      externalReference: "shipment-001",
      modelId,
      serialNumbers: ["SERIAL-001", "SERIAL-002"]
    });

    expect(receipt.stockItemIds).toHaveLength(2);
    expect((await database.query("select * from inventory.stock_item")).rows).toHaveLength(2);
    expect((await database.query("select state from procurement.purchase_order where purchase_order_id = $1", [purchaseOrderId])).rows[0]?.state).toBe("received");
  });

  it("creates an order, enforces workflow permissions and allocates stock once", async () => {
    const supplierId = await operations.createSupplier(organizationId, "factory-two", "Factory Two Ltd.");
    const purchaseOrderId = await operations.createPurchaseOrder({ organizationId, supplierId, currency: "USD", lines: [{ modelId, quantity: 1, unitMinorUnits: 120_000 }] });
    const receipt = await operations.receiveSerializedStock({ organizationId, purchaseOrderId, externalReference: "shipment-002", modelId, serialNumbers: ["SERIAL-003"] });
    const customerId = await operations.createCustomer(organizationId, "customer:test", "customer@example.com", "Customer Test");
    const orderId = await operations.createOrder({ organizationId, customerId, customerPrincipalId: "customer:test", currency: "USD", modelId, unitMinorUnits: 120_000 });

    await operations.transitionOrder(orderId, "placed", new Set(["order:create"]));
    await operations.transitionOrder(orderId, "confirmed", new Set(["order:transition"]));
    await operations.transitionOrder(orderId, "paid", new Set(["payment:reconcile"]));
    await operations.allocateSerializedStock(orderId, receipt.stockItemIds[0]!, new Set(["inventory:reserve"]));

    expect((await database.query("select state from sales.customer_order where order_id = $1", [orderId])).rows[0]?.state).toBe("allocated");
    expect((await database.query("select state from inventory.stock_item where stock_item_id = $1", [receipt.stockItemIds[0]])).rows[0]?.state).toBe("reserved");
    await expect(operations.allocateSerializedStock(orderId, receipt.stockItemIds[0]!, new Set(["inventory:reserve"]))).rejects.toBeInstanceOf(EnterpriseConflictError);
  });
});
````

### FILE: `src/modules/enterprise/operations.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-modules-enterprise-operations-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "fec8b02c422ba9a83d19de057d8de0e8231b6a9f6571931a84a3a376acba7819"
variables: []
secrets_allowed: false
```

````typescript
import { randomUUID } from "node:crypto";
import type { BusinessConfig } from "@/platform/config/schema";
import type { Database, SqlExecutor } from "@/platform/db/types";
import { transitionWorkflow } from "@/platform/workflow/engine";

export interface PurchaseOrderLineInput {
  modelId: string;
  quantity: number;
  unitMinorUnits: number;
}

export class EnterpriseConflictError extends Error {
  readonly code = "ENTERPRISE_CONFLICT";
}

export class EnterpriseOperations {
  constructor(private readonly database: Database, private readonly config: BusinessConfig) {}

  async createSupplier(organizationId: string, supplierCode: string, legalName: string): Promise<string> {
    const supplierId = randomUUID();
    await this.database.query(`
      insert into procurement.supplier
        (supplier_id, organization_id, supplier_code, legal_name, status)
      values ($1, $2, $3, $4, 'approved')
    `, [supplierId, organizationId, supplierCode, legalName]);
    return supplierId;
  }

  async createPurchaseOrder(input: {
    organizationId: string;
    supplierId: string;
    currency: string;
    lines: PurchaseOrderLineInput[];
  }): Promise<string> {
    if (input.lines.length === 0) throw new EnterpriseConflictError("purchase order requires at least one line");
    const purchaseOrderId = randomUUID();
    const total = input.lines.reduce((sum, line) => sum + line.quantity * line.unitMinorUnits, 0);
    if (!Number.isSafeInteger(total) || total < 0) throw new EnterpriseConflictError("purchase order total is unsafe");

    await this.database.transaction(async (transaction) => {
      await transaction.query(`
        insert into procurement.purchase_order
          (purchase_order_id, organization_id, supplier_id, state, currency, total_minor_units)
        values ($1, $2, $3, 'draft', $4, $5)
      `, [purchaseOrderId, input.organizationId, input.supplierId, input.currency, total]);

      for (const [index, line] of input.lines.entries()) {
        if (!Number.isSafeInteger(line.quantity) || line.quantity <= 0) throw new EnterpriseConflictError("line quantity must be a positive safe integer");
        if (!Number.isSafeInteger(line.unitMinorUnits) || line.unitMinorUnits < 0) throw new EnterpriseConflictError("unit amount must be a non-negative safe integer");
        await transaction.query(`
          insert into procurement.purchase_order_line
            (purchase_order_id, line_number, model_id, quantity, unit_minor_units)
          values ($1, $2, $3, $4, $5)
        `, [purchaseOrderId, index + 1, line.modelId, line.quantity, line.unitMinorUnits]);
      }
      await this.record(transaction, input.organizationId, "purchase_order", purchaseOrderId, "purchase_order.created");
    });
    return purchaseOrderId;
  }

  async receiveSerializedStock(input: {
    organizationId: string;
    purchaseOrderId: string;
    externalReference: string;
    modelId: string;
    serialNumbers: string[];
  }): Promise<{ shipmentId: string; stockItemIds: string[] }> {
    if (input.serialNumbers.length === 0 || new Set(input.serialNumbers).size !== input.serialNumbers.length) {
      throw new EnterpriseConflictError("serial numbers must be non-empty and unique");
    }
    const shipmentId = randomUUID();
    const stockItemIds = input.serialNumbers.map(() => randomUUID());

    await this.database.transaction(async (transaction) => {
      await transaction.query(`
        insert into factory.shipment
          (shipment_id, organization_id, purchase_order_id, external_reference, state, departed_at, received_at)
        values ($1, $2, $3, $4, 'received', now(), now())
      `, [shipmentId, input.organizationId, input.purchaseOrderId, input.externalReference]);
      for (const [index, serialNumber] of input.serialNumbers.entries()) {
        await transaction.query(`
          insert into inventory.stock_item
            (stock_item_id, organization_id, model_id, serial_number, state)
          values ($1, $2, $3, $4, 'available')
        `, [stockItemIds[index], input.organizationId, input.modelId, serialNumber]);
      }
      await transaction.query(`
        update procurement.purchase_order
           set state = 'received', version = version + 1, updated_at = now()
         where purchase_order_id = $1 and organization_id = $2
      `, [input.purchaseOrderId, input.organizationId]);
      await this.record(transaction, input.organizationId, "shipment", shipmentId, "shipment.received");
    });
    return { shipmentId, stockItemIds };
  }

  async createCustomer(organizationId: string, principalId: string, email: string, name: string): Promise<string> {
    const customerId = randomUUID();
    await this.database.query(`
      insert into crm.customer
        (customer_id, organization_id, principal_id, email, name, status)
      values ($1, $2, $3, lower($4), $5, 'active')
    `, [customerId, organizationId, principalId, email, name]);
    return customerId;
  }

  async createOrder(input: {
    organizationId: string;
    customerId: string;
    customerPrincipalId: string;
    currency: string;
    modelId: string;
    unitMinorUnits: number;
  }): Promise<string> {
    if (!Number.isSafeInteger(input.unitMinorUnits) || input.unitMinorUnits < 0) throw new EnterpriseConflictError("order amount is unsafe");
    const orderId = randomUUID();
    const initial = this.config.workflows.order?.initial;
    if (!initial) throw new EnterpriseConflictError("order workflow is not configured");
    await this.database.transaction(async (transaction) => {
      await transaction.query(`
        insert into sales.customer_order
          (order_id, organization_id, customer_principal_id, customer_id, state, currency, total_minor_units)
        values ($1, $2, $3, $4, $5, $6, $7)
      `, [orderId, input.organizationId, input.customerPrincipalId, input.customerId, initial, input.currency, input.unitMinorUnits]);
      await transaction.query(`
        insert into sales.order_line
          (order_id, line_number, model_id, quantity, unit_minor_units)
        values ($1, 1, $2, 1, $3)
      `, [orderId, input.modelId, input.unitMinorUnits]);
      await this.record(transaction, input.organizationId, "order", orderId, "order.created");
    });
    return orderId;
  }

  async transitionOrder(orderId: string, target: string, permissions: Set<string>): Promise<void> {
    const workflow = this.config.workflows.order;
    if (!workflow) throw new EnterpriseConflictError("order workflow is not configured");
    await this.database.transaction(async (transaction) => {
      const result = await transaction.query<{ state: string; version: number; organization_id: string } & Record<string, unknown>>(`
        select state, version, organization_id from sales.customer_order where order_id = $1
      `, [orderId]);
      const order = result.rows[0];
      if (!order) throw new EnterpriseConflictError("order not found");
      const next = transitionWorkflow(workflow, { current: order.state, target, permissions });
      const updated = await transaction.query(`
        update sales.customer_order
           set state = $2, version = version + 1, updated_at = now()
         where order_id = $1 and version = $3 and state = $4
      `, [orderId, next, order.version, order.state]);
      if (updated.affectedRows !== 1) throw new EnterpriseConflictError("concurrent order update");
      await this.record(transaction, order.organization_id, "order", orderId, `order.${next}`);
    });
  }

  async allocateSerializedStock(orderId: string, stockItemId: string, permissions: Set<string>): Promise<void> {
    if (!permissions.has("*") && !permissions.has("inventory:reserve")) throw new EnterpriseConflictError("inventory:reserve permission is required");
    await this.database.transaction(async (transaction) => {
      const orderResult = await transaction.query<{ state: string; version: number; organization_id: string } & Record<string, unknown>>(`
        select state, version, organization_id from sales.customer_order where order_id = $1
      `, [orderId]);
      const order = orderResult.rows[0];
      if (!order || order.state !== "paid") throw new EnterpriseConflictError("order must be paid before allocation");
      const stock = await transaction.query(`
        update inventory.stock_item
           set state = 'reserved', reservation_id = $2, version = version + 1, updated_at = now()
         where stock_item_id = $1 and organization_id = $3 and state = 'available'
      `, [stockItemId, orderId, order.organization_id]);
      if (stock.affectedRows !== 1) throw new EnterpriseConflictError("stock is unavailable or belongs to another organization");
      const updated = await transaction.query(`
        update sales.customer_order
           set state = 'allocated', version = version + 1, updated_at = now()
         where order_id = $1 and version = $2 and state = 'paid'
      `, [orderId, order.version]);
      if (updated.affectedRows !== 1) throw new EnterpriseConflictError("concurrent order allocation");
      await this.record(transaction, order.organization_id, "order", orderId, "order.allocated");
    });
  }

  private async record(transaction: SqlExecutor, organizationId: string, resourceType: string, resourceId: string, eventType: string): Promise<void> {
    await transaction.query(`
      insert into platform.outbox_event
        (event_id, aggregate_type, aggregate_id, aggregate_version, event_type, schema_version, payload)
      values ($1, $2, $3, 1, $4, 1, $5::jsonb)
    `, [randomUUID(), resourceType, resourceId, eventType, JSON.stringify({ organizationId, resourceId })]);
    await transaction.query(`
      insert into platform.audit_event
        (audit_event_id, organization_id, actor_id, action, resource_type, resource_id, outcome)
      values ($1, $2, 'system:enterprise-operations', $3, $4, $5, 'success')
    `, [randomUUID(), organizationId, eventType, resourceType, resourceId]);
  }
}
````

### FILE: `src/platform/auth/authorization.test.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-auth-authorization-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "9b3c738aede400de2a9d1e452a4e22077e51c1af1523c773c9731c1f913271ad"
variables: []
secrets_allowed: false
```

````typescript
import { describe, expect, it } from "vitest";
import { isAuthorized } from "./authorization";

const roles = [
  { id: "manager", permissions: ["lead:read"] },
  { id: "customer", permissions: ["own:order:read"] }
];

describe("organization-scoped authorization", () => {
  it("authorizes a permission inside the membership organization", () => {
    expect(isAuthorized({ principal: { id: "p1", memberships: [{ organizationId: "org-a", roleIds: ["manager"] }] }, organizationId: "org-a", permission: "lead:read", roles })).toBe(true);
  });

  it("denies the same permission across organizations", () => {
    expect(isAuthorized({ principal: { id: "p1", memberships: [{ organizationId: "org-a", roleIds: ["manager"] }] }, organizationId: "org-b", permission: "lead:read", roles })).toBe(false);
  });

  it("requires ownership for own-resource grants", () => {
    const base = { principal: { id: "p2", memberships: [{ organizationId: "org-a", roleIds: ["customer"] }] }, organizationId: "org-a", permission: "order:read", roles };
    expect(isAuthorized({ ...base, ownsResource: false })).toBe(false);
    expect(isAuthorized({ ...base, ownsResource: true })).toBe(true);
  });
});
````

### FILE: `src/platform/auth/authorization.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-auth-authorization-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "4d45f092c948ec960499630b38f3f51b9890583759f29004562e0fcd0db741ec"
variables: []
secrets_allowed: false
```

````typescript
export interface Membership {
  organizationId: string;
  roleIds: string[];
}

export interface Principal {
  id: string;
  memberships: Membership[];
}

export interface RoleDefinition {
  id: string;
  permissions: string[];
}

export interface AuthorizationContext {
  principal: Principal;
  organizationId: string;
  permission: string;
  roles: RoleDefinition[];
  ownsResource?: boolean;
}

function permissionMatches(granted: string, required: string, ownsResource: boolean): boolean {
  if (granted === "*") return true;
  if (granted === required) return true;
  if (ownsResource && granted.startsWith("own:") && granted.slice(4) === required) return true;
  return false;
}

export function isAuthorized(context: AuthorizationContext): boolean {
  const membership = context.principal.memberships.find((item) => item.organizationId === context.organizationId);
  if (!membership) return false;

  const roleMap = new Map(context.roles.map((role) => [role.id, role]));
  return membership.roleIds.some((roleId) => {
    const role = roleMap.get(roleId);
    return role?.permissions.some((granted) => permissionMatches(granted, context.permission, context.ownsResource ?? false)) ?? false;
  });
}

export function requireAuthorized(context: AuthorizationContext): void {
  if (!isAuthorized(context)) throw new AuthorizationError(context.permission, context.organizationId);
}

export class AuthorizationError extends Error {
  readonly code = "AUTHORIZATION_DENIED";
  constructor(readonly permission: string, readonly organizationId: string) {
    super(`permission ${permission} denied for organization ${organizationId}`);
  }
}
````

### FILE: `src/platform/auth/oidc.test.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-auth-oidc-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "f28f435a754222dfb14f9468db0979ac022a935347efb18d28864e8933f6d600"
variables: []
secrets_allowed: false
```

````typescript
import { createLocalJWKSet, exportJWK, generateKeyPair, SignJWT } from "jose";
import { describe, expect, it } from "vitest";
import { AuthenticationError, principalFromVerifiedClaims, verifyOidcToken, type OidcSettings } from "./oidc";

const settings: OidcSettings = {
  issuer: "https://identity.example.com",
  audience: "enterprise-platform",
  jwksUri: "https://identity.example.com/jwks.json",
  algorithms: ["RS256"]
};

describe("OIDC boundary", () => {
  it("verifies signature, issuer, audience and maps scoped memberships", async () => {
    const { publicKey, privateKey } = await generateKeyPair("RS256");
    const publicJwk = await exportJWK(publicKey);
    publicJwk.kid = "test-key";
    const token = await new SignJWT({
      org_memberships: [{ organizationId: "018f4d4a-7b36-7a21-8d10-2f4c54c28a01", roleIds: ["tenant_admin"] }]
    })
      .setProtectedHeader({ alg: "RS256", kid: "test-key" })
      .setIssuer(settings.issuer)
      .setAudience(settings.audience)
      .setSubject("user-123")
      .setIssuedAt()
      .setExpirationTime("5m")
      .sign(privateKey);

    const claims = await verifyOidcToken(token, settings, createLocalJWKSet({ keys: [publicJwk] }));
    expect(principalFromVerifiedClaims(claims)).toEqual({
      id: "user-123",
      memberships: [{ organizationId: "018f4d4a-7b36-7a21-8d10-2f4c54c28a01", roleIds: ["tenant_admin"] }]
    });
  });

  it("rejects a token for a different audience without leaking verifier detail", async () => {
    const { publicKey, privateKey } = await generateKeyPair("RS256");
    const publicJwk = await exportJWK(publicKey);
    publicJwk.kid = "test-key";
    const token = await new SignJWT({ org_memberships: [] })
      .setProtectedHeader({ alg: "RS256", kid: "test-key" })
      .setIssuer(settings.issuer)
      .setAudience("different-service")
      .setSubject("user-123")
      .setIssuedAt()
      .setExpirationTime("5m")
      .sign(privateKey);

    await expect(verifyOidcToken(token, settings, createLocalJWKSet({ keys: [publicJwk] }))).rejects.toBeInstanceOf(AuthenticationError);
  });
});
````

### FILE: `src/platform/auth/oidc.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-auth-oidc-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "dcbdc3fbfdbb6bd7d3dddd027fa608eb51e320b97bb3eb407607e3e591c65ba2"
variables: []
secrets_allowed: false
```

````typescript
import { createRemoteJWKSet, jwtVerify, type JWTPayload, type JWTVerifyGetKey } from "jose";
import { z } from "zod";
import type { Principal } from "./authorization";

const membershipClaimSchema = z.array(z.object({
  organizationId: z.uuid(),
  roleIds: z.array(z.string().regex(/^[a-z][a-z0-9_]*$/)).min(1).max(20)
})).max(100);

export interface OidcSettings {
  issuer: string;
  audience: string;
  jwksUri: string;
  algorithms: readonly ("RS256" | "ES256")[];
}

export class AuthenticationError extends Error {
  readonly code = "AUTHENTICATION_FAILED";
}

export function loadOidcSettings(environment: NodeJS.ProcessEnv = process.env): OidcSettings {
  const issuer = requiredHttpsUrl(environment.OIDC_ISSUER, "OIDC_ISSUER");
  const jwksUri = requiredHttpsUrl(environment.OIDC_JWKS_URI, "OIDC_JWKS_URI");
  const audience = environment.OIDC_AUDIENCE?.trim();
  if (!audience || audience.length > 200) throw new AuthenticationError("OIDC_AUDIENCE is required");
  if (!jwksUri.startsWith(`${issuer}/`) && new URL(jwksUri).origin !== new URL(issuer).origin) {
    throw new AuthenticationError("OIDC_JWKS_URI must share the trusted issuer origin unless the adapter is explicitly extended");
  }
  return { issuer, audience, jwksUri, algorithms: ["RS256", "ES256"] };
}

export function bearerToken(request: Request): string {
  const header = request.headers.get("authorization");
  const match = header?.match(/^Bearer ([A-Za-z0-9._~-]+)$/);
  if (!match) throw new AuthenticationError("a single Bearer token is required");
  if (match[1]!.length > 16_384) throw new AuthenticationError("Bearer token exceeds the accepted limit");
  return match[1]!;
}

export async function authenticateRequest(request: Request, settings = loadOidcSettings()): Promise<Principal> {
  return principalFromVerifiedClaims(await verifyOidcToken(bearerToken(request), settings));
}

export async function verifyOidcToken(
  token: string,
  settings: OidcSettings,
  keySet: JWTVerifyGetKey = createRemoteJWKSet(new URL(settings.jwksUri), {
    timeoutDuration: 5_000,
    cooldownDuration: 30_000,
    cacheMaxAge: 600_000
  })
): Promise<JWTPayload> {
  try {
    const result = await jwtVerify(token, keySet, {
      issuer: settings.issuer,
      audience: settings.audience,
      algorithms: [...settings.algorithms],
      clockTolerance: 5,
      maxTokenAge: "15m",
      requiredClaims: ["sub", "iat", "exp"]
    });
    return result.payload;
  } catch {
    throw new AuthenticationError("token verification failed");
  }
}

export function principalFromVerifiedClaims(claims: JWTPayload): Principal {
  if (typeof claims.sub !== "string" || claims.sub.length < 1 || claims.sub.length > 255) {
    throw new AuthenticationError("verified token subject is invalid");
  }
  const memberships = membershipClaimSchema.safeParse(claims.org_memberships);
  if (!memberships.success) throw new AuthenticationError("verified token organization memberships are invalid");
  return {
    id: claims.sub,
    memberships: memberships.data.map((membership) => ({
      organizationId: membership.organizationId,
      roleIds: [...new Set(membership.roleIds)]
    }))
  };
}

function requiredHttpsUrl(raw: string | undefined, label: string): string {
  if (!raw) throw new AuthenticationError(`${label} is required`);
  let url: URL;
  try {
    url = new URL(raw);
  } catch {
    throw new AuthenticationError(`${label} must be an absolute URL`);
  }
  if (url.protocol !== "https:" || url.username || url.password || url.hash || url.search) {
    throw new AuthenticationError(`${label} must be a clean HTTPS URL`);
  }
  return url.toString().replace(/\/$/, "");
}
````

### FILE: `src/platform/bootstrap.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-bootstrap-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "be48957607bd8b432968834adef57e0ccec5e1c350aef5b98877c8b59c18d21b"
variables: []
secrets_allowed: false
```

````typescript
import { loadBusinessConfig } from "./config/load";
import { database } from "./db/client";
import { migrate } from "./db/migrations";
import { seed } from "./seed";
import { validateRuntimePolicy } from "./config/runtime-policy";
import { validateEnabledIntegrations } from "./integrations/contract";
import { installedIntegrationAdapters } from "./integrations/registry";

let bootstrapped: Promise<void> | undefined;

export function ensureBootstrapped(): Promise<void> {
  bootstrapped ??= (async () => {
    const [db, config] = await Promise.all([database(), loadBusinessConfig()]);
    validateRuntimePolicy(config);
    validateEnabledIntegrations(config, installedIntegrationAdapters);
    await migrate(db);
    await seed(db, config);
  })();
  return bootstrapped;
}

export function resetBootstrapForTests(): void {
  bootstrapped = undefined;
}
````

### FILE: `src/platform/config/load.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-config-load-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "b61d038dc12ccb77e9002a0fedc8c6310c3cdd1e7b6b5ba03c660b5258642920"
variables: []
secrets_allowed: false
```

````typescript
import { readFile } from "node:fs/promises";
import { basename, join } from "node:path";
import { businessConfigSchema, type BusinessConfig } from "./schema";

let cached: BusinessConfig | undefined;

export async function loadBusinessConfig(options?: { path?: string; bypassCache?: boolean }): Promise<BusinessConfig> {
  if (cached && !options?.bypassCache) return cached;

  const configFile = basename(options?.path ?? process.env.BUSINESS_CONFIG_FILE ?? "business.example.json");
  const configPath = join(process.cwd(), "config", configFile);
  const raw = JSON.parse(await readFile(configPath, "utf8")) as unknown;
  const parsed = businessConfigSchema.parse(raw);
  if (!options?.bypassCache) cached = parsed;
  return parsed;
}

export function clearBusinessConfigCache(): void {
  cached = undefined;
}
````

### FILE: `src/platform/config/registry.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-config-registry-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "cedca50639bb9423164e05c767144ed957dcf459d1309bf5c3255c7779bbf441"
variables: []
secrets_allowed: false
```

````typescript
import type { BusinessConfig } from "./schema";

export interface NavigationItem {
  id: string;
  label: string;
  href: string;
  module?: string;
  feature?: string;
  permission?: string;
}

const navigation: NavigationItem[] = [
  { id: "public_catalog", label: "Modelos", href: "/models", module: "catalog", feature: "public_catalog", permission: "catalog:read" },
  { id: "customer", label: "Mi cuenta", href: "/customer", feature: "customer_portal" },
  { id: "admin", label: "Operación", href: "/admin", permission: "lead:read" },
  { id: "factory", label: "Fábrica", href: "/factory", module: "procurement", feature: "factory_portal", permission: "procurement:read" }
];

export function enabledNavigation(config: BusinessConfig): NavigationItem[] {
  return navigation.filter((item) => {
    if (item.module && !config.modules[item.module]?.enabled) return false;
    if (item.feature && !config.features[item.feature]) return false;
    return true;
  });
}

export function enabledModules(config: BusinessConfig): string[] {
  return Object.entries(config.modules).filter(([, value]) => value.enabled).map(([name]) => name);
}
````

### FILE: `src/platform/config/runtime-policy.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-config-runtime-policy-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "682ff3bea7c4f754ac3bc01c7f4c1e7464d9299ab0582d46fa56156f30f8a65a"
variables: []
secrets_allowed: false
```

````typescript
import type { BusinessConfig } from "./schema";

export function validateRuntimePolicy(config: BusinessConfig, environment: Record<string, string | undefined> = process.env): void {
  const appEnvironment = environment.APP_ENV ?? "development";
  if (appEnvironment !== "production") return;
  const databaseUrl = environment.DATABASE_URL ?? "";
  if (!databaseUrl.startsWith("postgresql://") && !databaseUrl.startsWith("postgres://")) throw new Error("production requires an external PostgreSQL DATABASE_URL");
  if (environment.ALLOW_DEMO_IDENTITY === "true") throw new Error("demo identity is forbidden in production");
  if (!environment.OIDC_ISSUER || !environment.OIDC_CLIENT_ID || !environment.OIDC_CLIENT_SECRET_REF) throw new Error("production requires OIDC issuer, client id and secret reference");
  if (config.business.supportEmail.endsWith(".invalid")) throw new Error("production requires a deliverable support email");
}
````

### FILE: `src/platform/config/schema.test.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-config-schema-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "4afcdfc68821dc303ab051ca0303237fc6e1f0c61f928ccd2efc27639e0c0f9c"
variables: []
secrets_allowed: false
```

````typescript
import { describe, expect, it } from "vitest";
import { loadBusinessConfig } from "./load";
import { businessConfigSchema } from "./schema";
import { validateRuntimePolicy } from "./runtime-policy";

describe("business configuration", () => {
  it("accepts the reference electric mobility configuration", async () => {
    const config = await loadBusinessConfig({ bypassCache: true });
    expect(config.business.id).toBe("electric-mobility-network");
    expect(config.modules.orders?.enabled).toBe(true);
  });

  it("rejects a module whose dependency is disabled", async () => {
    const config = structuredClone(await loadBusinessConfig({ bypassCache: true }));
    config.modules.orders = { enabled: false };
    expect(() => businessConfigSchema.parse(config)).toThrow(/payments requires enabled module orders/);
  });

  it("fails closed for an unsafe production runtime", async () => {
    const config = await loadBusinessConfig({ bypassCache: true });
    expect(() => validateRuntimePolicy(config, { APP_ENV: "production", DATABASE_URL: "memory://", ALLOW_DEMO_IDENTITY: "true" })).toThrow(/external PostgreSQL/);
  });
});
````

### FILE: `src/platform/config/schema.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-config-schema-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "c8ad5a8d05122ca180fdb1be4098d04096f78b3039a672c241df703223584892"
variables: []
secrets_allowed: false
```

````typescript
import { z } from "zod";

const identifier = z.string().regex(/^[a-z][a-z0-9_]*$/, "use snake_case identifiers");
const permission = z.string().regex(/^\*$|^[a-z][a-z0-9_]*(?::[a-z][a-z0-9_]*){1,2}$/);

const fieldSchema = z.discriminatedUnion("type", [
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("string"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("number"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("boolean"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("date"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("select"), required: z.boolean(), options: z.array(identifier).min(1) })
]);

const workflowSchema = z.object({
  initial: identifier,
  states: z.array(identifier).min(1),
  transitions: z.array(z.object({ from: identifier, to: identifier, permission }))
});

export const businessConfigSchema = z.object({
  schemaVersion: z.literal("1.0.0"),
  business: z.object({
    id: z.string().regex(/^[a-z][a-z0-9-]*$/),
    name: z.string().min(1),
    defaultLocale: z.string().min(2),
    defaultMarket: z.string().length(2),
    supportEmail: z.email()
  }),
  markets: z.array(z.object({
    code: z.string().length(2),
    name: z.string().min(1),
    currency: z.string().length(3),
    locales: z.array(z.string().min(2)).min(1),
    timeZone: z.string().min(1),
    taxMode: z.enum(["internal", "external"])
  })).min(1),
  organizationTypes: z.array(z.object({
    id: identifier,
    label: z.string().min(1),
    allowedParents: z.array(identifier)
  })).min(1),
  roles: z.array(z.object({
    id: identifier,
    label: z.string().min(1),
    permissions: z.array(permission).min(1)
  })).min(1),
  modules: z.record(identifier, z.object({ enabled: z.boolean() })),
  workflows: z.record(identifier, workflowSchema),
  customFields: z.record(identifier, z.array(fieldSchema)),
  integrations: z.array(z.object({
    id: identifier,
    provider: identifier,
    enabled: z.boolean(),
    mode: z.enum(["sandbox", "production"]),
    capabilities: z.array(identifier).min(1),
    credentialRefEnv: z.string().regex(/^[A-Z][A-Z0-9_]+$/)
  })),
  features: z.record(identifier, z.boolean())
}).superRefine((config, context) => {
  const unique = (values: string[], path: (string | number)[], label: string) => {
    if (new Set(values).size !== values.length) {
      context.addIssue({ code: "custom", message: `${label} must be unique`, path });
    }
  };

  unique(config.markets.map((item) => item.code), ["markets"], "market codes");
  unique(config.organizationTypes.map((item) => item.id), ["organizationTypes"], "organization type ids");
  unique(config.roles.map((item) => item.id), ["roles"], "role ids");
  unique(config.integrations.map((item) => item.id), ["integrations"], "integration ids");

  const market = config.markets.find((item) => item.code === config.business.defaultMarket);
  if (!market) {
    context.addIssue({ code: "custom", message: "defaultMarket must reference an existing market", path: ["business", "defaultMarket"] });
  } else if (!market.locales.includes(config.business.defaultLocale)) {
    context.addIssue({ code: "custom", message: "defaultLocale must belong to defaultMarket", path: ["business", "defaultLocale"] });
  }

  const organizationTypeIds = new Set(config.organizationTypes.map((item) => item.id));
  config.organizationTypes.forEach((type, index) => {
    type.allowedParents.forEach((parent) => {
      if (!organizationTypeIds.has(parent)) {
        context.addIssue({ code: "custom", message: `unknown parent organization type: ${parent}`, path: ["organizationTypes", index, "allowedParents"] });
      }
      if (parent === type.id) {
        context.addIssue({ code: "custom", message: "organization type cannot parent itself", path: ["organizationTypes", index, "allowedParents"] });
      }
    });
  });

  for (const [name, workflow] of Object.entries(config.workflows)) {
    unique(workflow.states, ["workflows", name, "states"], `${name} states`);
    const states = new Set(workflow.states);
    if (!states.has(workflow.initial)) {
      context.addIssue({ code: "custom", message: "initial state must be declared", path: ["workflows", name, "initial"] });
    }
    workflow.transitions.forEach((transition, index) => {
      if (!states.has(transition.from) || !states.has(transition.to)) {
        context.addIssue({ code: "custom", message: "transition references an unknown state", path: ["workflows", name, "transitions", index] });
      }
      if (transition.from === transition.to) {
        context.addIssue({ code: "custom", message: "self transitions are not allowed", path: ["workflows", name, "transitions", index] });
      }
    });
  }

  for (const [entity, fields] of Object.entries(config.customFields)) {
    unique(fields.map((field) => field.id), ["customFields", entity], `${entity} custom field ids`);
  }

  const moduleDependencies: Record<string, string[]> = {
    payments: ["orders"],
    fulfillment: ["orders", "inventory"],
    service: ["catalog"],
    procurement: ["inventory"],
    integrations: ["catalog"]
  };
  for (const [moduleName, dependencies] of Object.entries(moduleDependencies)) {
    if (config.modules[moduleName]?.enabled) {
      dependencies.forEach((dependency) => {
        if (!config.modules[dependency]?.enabled) {
          context.addIssue({ code: "custom", message: `${moduleName} requires enabled module ${dependency}`, path: ["modules", moduleName] });
        }
      });
    }
  }
});

export type BusinessConfig = z.infer<typeof businessConfigSchema>;
export type WorkflowConfig = z.infer<typeof workflowSchema>;
export type CustomField = z.infer<typeof fieldSchema>;
````

### FILE: `src/platform/db/client.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-db-client-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "9dbd01ad36dcb2d4315f315de80e7572e0f39360e814e222d7f52cc1ccd55c11"
variables: []
secrets_allowed: false
```

````typescript
import { mkdir } from "node:fs/promises";
import { basename, join } from "node:path";
import { PGlite } from "@electric-sql/pglite";
import pg from "pg";
import type { Database, QueryResult, SqlExecutor } from "./types";

class PGliteDatabase implements Database {
  constructor(private readonly client: PGlite) {}

  async query<Row extends Record<string, unknown>>(sql: string, params: unknown[] = []): Promise<QueryResult<Row>> {
    const result = await this.client.query<Row>(sql, params);
    return { rows: result.rows, affectedRows: result.affectedRows };
  }

  async transaction<T>(callback: (transaction: SqlExecutor) => Promise<T>): Promise<T> {
    return this.client.transaction(async (transaction) => callback({
      query: async <Row extends Record<string, unknown>>(sql: string, params: unknown[] = []) => {
        const result = await transaction.query<Row>(sql, params);
        return { rows: result.rows, affectedRows: result.affectedRows };
      }
    }));
  }

  async close(): Promise<void> {
    await this.client.close();
  }
}

class PostgresDatabase implements Database {
  constructor(private readonly pool: pg.Pool) {}

  async query<Row extends Record<string, unknown>>(sql: string, params: unknown[] = []): Promise<QueryResult<Row>> {
    const result = await this.pool.query<Row>(sql, params);
    return { rows: result.rows, affectedRows: result.rowCount ?? undefined };
  }

  async transaction<T>(callback: (transaction: SqlExecutor) => Promise<T>): Promise<T> {
    const client = await this.pool.connect();
    try {
      await client.query("begin");
      const result = await callback({
        query: async <Row extends Record<string, unknown>>(sql: string, params: unknown[] = []) => {
          const queryResult = await client.query<Row>(sql, params);
          return { rows: queryResult.rows, affectedRows: queryResult.rowCount ?? undefined };
        }
      });
      await client.query("commit");
      return result;
    } catch (error) {
      await client.query("rollback");
      throw error;
    } finally {
      client.release();
    }
  }

  async close(): Promise<void> {
    await this.pool.end();
  }
}

let singleton: Promise<Database> | undefined;

export async function createDatabase(url = process.env.DATABASE_URL ?? "file://./.data/pglite"): Promise<Database> {
  if (url.startsWith("postgresql://") || url.startsWith("postgres://")) {
    return new PostgresDatabase(new pg.Pool({
      connectionString: url,
      max: Number(process.env.DATABASE_POOL_MAX ?? 10),
      connectionTimeoutMillis: 5_000,
      statement_timeout: 10_000,
      idle_in_transaction_session_timeout: 15_000
    }));
  }

  if (!url.startsWith("file://") && url !== "memory://") throw new Error("DATABASE_URL must use postgresql://, postgres://, file:// or memory://");
  if (url === "memory://") return new PGliteDatabase(new PGlite());

  const databaseName = basename(url.slice("file://".length));
  if (!/^[A-Za-z0-9_.-]+$/.test(databaseName)) throw new Error("embedded database name contains unsupported characters");
  const dataPath = join(process.cwd(), ".data", databaseName);
  await mkdir(join(process.cwd(), ".data"), { recursive: true });
  return new PGliteDatabase(new PGlite(`file://${dataPath.replaceAll("\\", "/")}`));
}

export function database(): Promise<Database> {
  singleton ??= createDatabase();
  return singleton;
}

export async function closeDatabase(): Promise<void> {
  if (singleton) await (await singleton).close();
  singleton = undefined;
}
````

### FILE: `src/platform/db/migrations.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-db-migrations-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "4aeb08519401ed47f9dcb3e1d28c15f59b7ce8c782245b58e23815a4a0073eed"
variables: []
secrets_allowed: false
```

````typescript
import type { Database } from "./types";

interface Migration {
  id: string;
  sql: string;
}

const migrations: Migration[] = [
  {
    id: "0001_platform_core",
    sql: `
      create schema if not exists platform;
      create schema if not exists org;
      create schema if not exists catalog;
      create schema if not exists crm;
      create schema if not exists inventory;
      create schema if not exists sales;

      create table if not exists platform.schema_migration (
        migration_id text primary key,
        applied_at timestamptz not null default now()
      );

      create table if not exists platform.config_snapshot (
        config_hash text primary key,
        schema_version text not null,
        business_id text not null,
        config jsonb not null,
        activated_at timestamptz not null default now()
      );

      create table if not exists platform.idempotency_record (
        scope text not null,
        idempotency_key text not null,
        request_hash text not null,
        status text not null check (status in ('processing','completed','failed_terminal')),
        response_code integer,
        response_body jsonb,
        resource_id uuid,
        created_at timestamptz not null default now(),
        expires_at timestamptz not null,
        primary key (scope, idempotency_key)
      );

      create table if not exists platform.outbox_event (
        event_id uuid primary key,
        aggregate_type text not null,
        aggregate_id uuid not null,
        aggregate_version bigint not null,
        event_type text not null,
        schema_version integer not null,
        tenant_id uuid,
        payload jsonb not null,
        occurred_at timestamptz not null default now(),
        available_at timestamptz not null default now(),
        published_at timestamptz,
        attempts integer not null default 0
      );
      create index if not exists outbox_ready_idx on platform.outbox_event (available_at, occurred_at) where published_at is null;

      create table if not exists platform.audit_event (
        audit_event_id uuid primary key,
        organization_id uuid,
        actor_id text,
        action text not null,
        resource_type text not null,
        resource_id text not null,
        outcome text not null,
        reason text,
        trace_id text,
        occurred_at timestamptz not null default now(),
        details jsonb not null default '{}'::jsonb
      );
    `
  },
  {
    id: "0002_business_modules",
    sql: `
      create table if not exists org.organization (
        organization_id uuid primary key,
        parent_organization_id uuid references org.organization(organization_id),
        organization_type text not null,
        code text not null unique,
        name text not null,
        status text not null check (status in ('draft','active','suspended','closed')),
        version bigint not null default 0,
        created_at timestamptz not null default now()
      );

      create table if not exists org.membership (
        organization_id uuid not null references org.organization(organization_id),
        principal_id text not null,
        role_id text not null,
        valid_from timestamptz not null default now(),
        valid_until timestamptz,
        primary key (organization_id, principal_id, role_id)
      );

      create table if not exists catalog.model (
        model_id uuid primary key,
        slug text not null unique,
        name text not null,
        summary text not null,
        status text not null check (status in ('draft','published','retired')),
        market_code char(2) not null,
        custom_fields jsonb not null default '{}'::jsonb,
        version bigint not null default 0,
        published_at timestamptz,
        created_at timestamptz not null default now()
      );

      create table if not exists crm.lead (
        lead_id uuid primary key,
        organization_id uuid references org.organization(organization_id),
        market_code char(2) not null,
        email text not null,
        name text not null,
        model_id uuid references catalog.model(model_id),
        consent_version text not null,
        source text not null,
        state text not null,
        custom_fields jsonb not null default '{}'::jsonb,
        version bigint not null default 0,
        created_at timestamptz not null default now(),
        updated_at timestamptz not null default now()
      );

      create table if not exists inventory.stock_item (
        stock_item_id uuid primary key,
        organization_id uuid not null references org.organization(organization_id),
        model_id uuid not null references catalog.model(model_id),
        serial_number text not null unique,
        state text not null check (state in ('in_transit','quality_hold','available','reserved','sold','service','retired')),
        reservation_id uuid,
        version bigint not null default 0,
        updated_at timestamptz not null default now(),
        check ((state = 'reserved') = (reservation_id is not null))
      );

      create table if not exists sales.customer_order (
        order_id uuid primary key,
        organization_id uuid not null references org.organization(organization_id),
        customer_principal_id text not null,
        state text not null,
        currency char(3) not null,
        total_minor_units bigint not null check (total_minor_units >= 0),
        version bigint not null default 0,
        created_at timestamptz not null default now(),
        updated_at timestamptz not null default now()
      );
    `
  },
  {
    id: "0003_public_abuse_controls",
    sql: `
      create table if not exists platform.rate_limit_window (
        scope text not null,
        key_hash text not null,
        window_started timestamptz not null,
        request_count integer not null check (request_count > 0),
        primary key (scope, key_hash)
      )
    `
  },
  {
    id: "0004_enterprise_operations",
    sql: `
      create schema if not exists procurement;
      create schema if not exists factory;
      create schema if not exists service;
      create schema if not exists integration;

      create table if not exists procurement.supplier (
        supplier_id uuid primary key,
        organization_id uuid not null references org.organization(organization_id),
        supplier_code text not null,
        legal_name text not null,
        status text not null check (status in ('candidate','approved','suspended','retired')),
        version bigint not null default 0,
        created_at timestamptz not null default now(),
        unique (organization_id, supplier_code)
      );

      create table if not exists procurement.purchase_order (
        purchase_order_id uuid primary key,
        organization_id uuid not null references org.organization(organization_id),
        supplier_id uuid not null references procurement.supplier(supplier_id),
        state text not null check (state in ('draft','submitted','accepted','in_production','shipped','received','cancelled')),
        currency char(3) not null,
        total_minor_units bigint not null check (total_minor_units >= 0),
        version bigint not null default 0,
        created_at timestamptz not null default now(),
        updated_at timestamptz not null default now()
      );

      create table if not exists procurement.purchase_order_line (
        purchase_order_id uuid not null references procurement.purchase_order(purchase_order_id),
        line_number integer not null check (line_number > 0),
        model_id uuid not null references catalog.model(model_id),
        quantity integer not null check (quantity > 0),
        unit_minor_units bigint not null check (unit_minor_units >= 0),
        primary key (purchase_order_id, line_number)
      );

      create table if not exists factory.shipment (
        shipment_id uuid primary key,
        organization_id uuid not null references org.organization(organization_id),
        purchase_order_id uuid references procurement.purchase_order(purchase_order_id),
        external_reference text not null,
        state text not null check (state in ('planned','in_transit','customs','received','exception')),
        departed_at timestamptz,
        received_at timestamptz,
        version bigint not null default 0,
        unique (organization_id, external_reference),
        check (received_at is null or departed_at is not null)
      );

      create table if not exists crm.customer (
        customer_id uuid primary key,
        organization_id uuid not null references org.organization(organization_id),
        principal_id text,
        email text not null,
        name text not null,
        status text not null check (status in ('active','blocked','erasure_pending','erased')),
        version bigint not null default 0,
        created_at timestamptz not null default now(),
        unique (organization_id, email)
      );

      alter table sales.customer_order
        add column if not exists customer_id uuid references crm.customer(customer_id);

      create table if not exists sales.order_line (
        order_id uuid not null references sales.customer_order(order_id),
        line_number integer not null check (line_number > 0),
        model_id uuid not null references catalog.model(model_id),
        quantity integer not null check (quantity > 0),
        unit_minor_units bigint not null check (unit_minor_units >= 0),
        primary key (order_id, line_number)
      );

      create table if not exists org.franchise_agreement (
        franchise_agreement_id uuid primary key,
        franchisor_organization_id uuid not null references org.organization(organization_id),
        franchisee_organization_id uuid not null references org.organization(organization_id),
        market_code char(2) not null,
        status text not null check (status in ('draft','active','suspended','terminated','expired')),
        valid_from date not null,
        valid_until date,
        version bigint not null default 0,
        check (franchisor_organization_id <> franchisee_organization_id),
        check (valid_until is null or valid_until >= valid_from)
      );

      create table if not exists service.warranty_claim (
        warranty_claim_id uuid primary key,
        organization_id uuid not null references org.organization(organization_id),
        stock_item_id uuid not null references inventory.stock_item(stock_item_id),
        customer_id uuid not null references crm.customer(customer_id),
        state text not null check (state in ('opened','triage','approved','rejected','repairing','resolved','closed')),
        issue_code text not null,
        description text not null,
        version bigint not null default 0,
        created_at timestamptz not null default now(),
        updated_at timestamptz not null default now()
      );

      create table if not exists integration.delivery (
        integration_id text not null,
        delivery_id uuid not null,
        operation text not null,
        idempotency_key text not null,
        state text not null check (state in ('pending','sent','acknowledged','retry','terminal')),
        request_payload jsonb not null,
        response_evidence jsonb,
        attempts integer not null default 0 check (attempts >= 0),
        next_attempt_at timestamptz not null default now(),
        created_at timestamptz not null default now(),
        primary key (integration_id, delivery_id),
        unique (integration_id, idempotency_key)
      )
    `
  }
];

export async function migrate(database: Database): Promise<string[]> {
  await executeStatements(database, `
    create schema if not exists platform;
    create table if not exists platform.schema_migration (
      migration_id text primary key,
      applied_at timestamptz not null default now()
    );
  `);

  const applied: string[] = [];
  for (const migration of migrations) {
    const existing = await database.query<{ migration_id: string }>(
      "select migration_id from platform.schema_migration where migration_id = $1",
      [migration.id]
    );
    if (existing.rows.length > 0) continue;

    await database.transaction(async (transaction) => {
      await executeStatements(transaction, migration.sql);
      await transaction.query("insert into platform.schema_migration (migration_id) values ($1)", [migration.id]);
    });
    applied.push(migration.id);
  }
  return applied;
}

async function executeStatements(executor: { query(sql: string): Promise<unknown> }, sql: string): Promise<void> {
  for (const statement of sql.split(";").map((item) => item.trim()).filter(Boolean)) {
    await executor.query(statement);
  }
}
````

### FILE: `src/platform/db/types.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-db-types-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "83eb5c05e3deecfad6e6300e27ecddcf1f691f43e4a571cbe8c79dc8d42c90c6"
variables: []
secrets_allowed: false
```

````typescript
export interface QueryResult<Row extends Record<string, unknown> = Record<string, unknown>> {
  rows: Row[];
  affectedRows: number | undefined;
}

export interface SqlExecutor {
  query<Row extends Record<string, unknown> = Record<string, unknown>>(sql: string, params?: unknown[]): Promise<QueryResult<Row>>;
}

export interface Database extends SqlExecutor {
  transaction<T>(callback: (transaction: SqlExecutor) => Promise<T>): Promise<T>;
  close(): Promise<void>;
}
````

### FILE: `src/platform/http/problem.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-http-problem-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "848e85d16e197b0a5473d2c9baf3bb6c28b256a17bfa8722b3dda80fa4c80f9d"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { ZodError } from "zod";
import { AuthorizationError } from "../auth/authorization";
import { IdempotencyConflictError, InvalidLeadInputError, RateLimitExceededError } from "@/modules/crm/lead-service";

export function problem(status: number, title: string, detail: string, code: string): NextResponse {
  return NextResponse.json({ type: `urn:elite-platform:problem:${code}`, title, status, detail, code }, {
    status,
    headers: { "content-type": "application/problem+json", "cache-control": "no-store" }
  });
}

export function errorResponse(error: unknown): NextResponse {
  if (error instanceof ZodError) return problem(400, "Solicitud inválida", error.issues.map((issue) => issue.message).join("; "), "VALIDATION_FAILED");
  if (error instanceof InvalidLeadInputError) return problem(400, "Solicitud inválida", error.message, error.code);
  if (error instanceof IdempotencyConflictError) return problem(409, "Conflicto de idempotencia", error.message, error.code);
  if (error instanceof RateLimitExceededError) {
    const response = problem(429, "Demasiadas solicitudes", "Intente nuevamente más tarde.", error.code);
    response.headers.set("retry-after", String(error.retryAfterSeconds));
    return response;
  }
  if (error instanceof AuthorizationError) return problem(403, "Acceso denegado", "La identidad no posee el permiso requerido.", error.code);
  const reference = crypto.randomUUID();
  console.error(JSON.stringify({ level: "error", event: "http.unhandled", reference, message: error instanceof Error ? error.message : "unknown" }));
  return problem(500, "Error interno", `Referencia: ${reference}`, "INTERNAL_ERROR");
}
````

### FILE: `src/platform/http/problem.test.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-http-problem-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "5164a9fc934527e0b749c274ba81c0836a8d17c250682daad4574ead4afbbafb"
variables: []
secrets_allowed: false
```

````typescript
import { describe, expect, it } from "vitest";
import { InvalidLeadInputError } from "@/modules/crm/lead-service";
import { errorResponse } from "./problem";

describe("HTTP problem classification", () => {
  it("returns a stable 400 problem for invalid domain input", async () => {
    const response = errorResponse(new InvalidLeadInputError("unknown custom field"));

    expect(response.status).toBe(400);
    await expect(response.json()).resolves.toMatchObject({
      status: 400,
      code: "INVALID_LEAD_INPUT",
      detail: "unknown custom field"
    });
  });
});
````

### FILE: `src/platform/integrations/contract.test.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-integrations-contract-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "613959ede429c0e8594e7d24415f75eafe2f45784c70ebb7baac9cbf00ef07ab"
variables: []
secrets_allowed: false
```

````typescript
import { describe, expect, it } from "vitest";
import { loadBusinessConfig } from "../config/load";
import { validateEnabledIntegrations } from "./contract";

describe("integration admission", () => {
  it("rejects an enabled provider without installed adapter", async () => {
    const config = structuredClone(await loadBusinessConfig({ bypassCache: true }));
    config.integrations[0]!.enabled = true;
    expect(() => validateEnabledIntegrations(config, [])).toThrow(/no installed adapter/);
  });
});
````

### FILE: `src/platform/integrations/contract.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-integrations-contract-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "a5e2de9da13a38e102b97f22c0ce428692e24ef9619b4e062cfbdecb928b8777"
variables: []
secrets_allowed: false
```

````typescript
import type { BusinessConfig } from "../config/schema";

export interface IntegrationCommand {
  commandId: string;
  capability: string;
  operation: string;
  payload: Record<string, unknown>;
}

export interface IntegrationResult {
  externalId: string;
  status: "accepted" | "completed";
}

export interface IntegrationAdapter {
  readonly provider: string;
  execute(command: IntegrationCommand): Promise<IntegrationResult>;
}

export function validateEnabledIntegrations(config: BusinessConfig, adapters: IntegrationAdapter[]): void {
  const providers = new Set(adapters.map((adapter) => adapter.provider));
  for (const integration of config.integrations.filter((item) => item.enabled)) {
    if (!providers.has(integration.provider)) throw new Error(`enabled integration has no installed adapter: ${integration.provider}`);
    if (!process.env[integration.credentialRefEnv]) throw new Error(`enabled integration is missing credential reference: ${integration.credentialRefEnv}`);
  }
}
````

### FILE: `src/platform/integrations/http-adapter.test.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-integrations-http-adapter-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "4f85335fa7cdca1d000558e353c013a649411d254bd38e52a1d56c7a5a8661fd"
variables: []
secrets_allowed: false
```

````typescript
import { describe, expect, it, vi } from "vitest";
import { SafeHttpAdapter } from "./http-adapter";

const credentials = { bearerToken: vi.fn(async () => "test-token") };
const settings = {
  provider: "marketplace",
  baseUrl: "https://provider.example.com/v1/",
  credentialReference: "secret-manager://marketplace/token",
  operations: { publish: { capability: "catalog", method: "POST" as const, path: "/v1/products" as const } },
  timeoutMs: 2_000,
  maxResponseBytes: 10_000
};

describe("safe integration HTTP adapter", () => {
  it("sends an admitted operation with idempotency and bounded authority", async () => {
    const fetcher = vi.fn(async () => Response.json({ id: "external-1", status: "accepted" }));
    const adapter = new SafeHttpAdapter(settings, credentials, fetcher as typeof fetch);
    await expect(adapter.execute({ commandId: "command:00000001", capability: "catalog", operation: "publish", payload: { sku: "E1" } })).resolves.toEqual({ externalId: "external-1", status: "accepted" });
    expect(fetcher).toHaveBeenCalledWith(new URL("https://provider.example.com/v1/products"), expect.objectContaining({ method: "POST", redirect: "error" }));
    const call = fetcher.mock.calls[0] as unknown as [URL, RequestInit];
    expect(call[1].headers).toMatchObject({ "idempotency-key": "command:00000001" });
  });

  it("rejects non-admitted operations before network access", async () => {
    const fetcher = vi.fn();
    const adapter = new SafeHttpAdapter(settings, credentials, fetcher as typeof fetch);
    await expect(adapter.execute({ commandId: "command:00000002", capability: "payments", operation: "publish", payload: {} })).rejects.toThrow(/not admitted/);
    expect(fetcher).not.toHaveBeenCalled();
  });

  it("classifies provider throttling as retryable without exposing its body", async () => {
    const fetcher = vi.fn(async () => new Response("sensitive upstream text", { status: 429 }));
    const adapter = new SafeHttpAdapter(settings, credentials, fetcher as typeof fetch);
    const failure = adapter.execute({ commandId: "command:00000003", capability: "catalog", operation: "publish", payload: {} });
    await expect(failure).rejects.toMatchObject({ retryable: true, status: 429 });
  });
});
````

### FILE: `src/platform/integrations/http-adapter.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-integrations-http-adapter-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "63ae6f7fc291c8a41f55728dfa929a9b9e212799c54b6447a7a6f453ef7d249c"
variables: []
secrets_allowed: false
```

````typescript
import type { IntegrationAdapter, IntegrationCommand, IntegrationResult } from "./contract";

export interface CredentialProvider {
  bearerToken(reference: string): Promise<string>;
}

export interface HttpOperation {
  capability: string;
  method: "POST" | "PUT" | "PATCH";
  path: `/${string}`;
}

export interface SafeHttpAdapterSettings {
  provider: string;
  baseUrl: string;
  credentialReference: string;
  operations: Record<string, HttpOperation>;
  timeoutMs: number;
  maxResponseBytes: number;
}

export class IntegrationTransportError extends Error {
  readonly code = "INTEGRATION_TRANSPORT_FAILED";
  constructor(message: string, readonly retryable: boolean, readonly status?: number) {
    super(message);
  }
}

export class SafeHttpAdapter implements IntegrationAdapter {
  readonly provider: string;
  private readonly baseUrl: URL;

  constructor(
    private readonly settings: SafeHttpAdapterSettings,
    private readonly credentials: CredentialProvider,
    private readonly fetcher: typeof fetch = fetch
  ) {
    this.provider = settings.provider;
    this.baseUrl = new URL(settings.baseUrl);
    if (this.baseUrl.protocol !== "https:" || this.baseUrl.username || this.baseUrl.password || this.baseUrl.search || this.baseUrl.hash) {
      throw new Error("integration baseUrl must be a clean HTTPS URL");
    }
    if (settings.timeoutMs < 100 || settings.timeoutMs > 30_000) throw new Error("integration timeout must be between 100 and 30000 ms");
    if (settings.maxResponseBytes < 1_024 || settings.maxResponseBytes > 5_000_000) throw new Error("integration response limit is invalid");
  }

  async execute(command: IntegrationCommand): Promise<IntegrationResult> {
    if (!/^[A-Za-z0-9_.:-]{16,128}$/.test(command.commandId)) throw new Error("integration commandId is invalid");
    const operation = this.settings.operations[command.operation];
    if (!operation || operation.capability !== command.capability) throw new Error("integration operation is not admitted");
    const body = JSON.stringify(command.payload);
    if (Buffer.byteLength(body, "utf8") > 1_000_000) throw new Error("integration payload exceeds one megabyte");
    const target = new URL(operation.path, this.baseUrl);
    if (target.origin !== this.baseUrl.origin || !target.pathname.startsWith(this.baseUrl.pathname)) throw new Error("integration path escaped the admitted base URL");
    const token = await this.credentials.bearerToken(this.settings.credentialReference);
    if (!token || /[\r\n]/.test(token)) throw new Error("credential provider returned an invalid token");

    let response: Response;
    try {
      response = await this.fetcher(target, {
        method: operation.method,
        redirect: "error",
        signal: AbortSignal.timeout(this.settings.timeoutMs),
        headers: {
          "accept": "application/json",
          "authorization": `Bearer ${token}`,
          "content-type": "application/json",
          "idempotency-key": command.commandId
        },
        body
      });
    } catch {
      throw new IntegrationTransportError("provider request failed before a trusted response", true);
    }

    if (!response.ok) {
      throw new IntegrationTransportError("provider returned a non-success status", response.status === 408 || response.status === 429 || response.status >= 500, response.status);
    }
    const declaredLength = Number(response.headers.get("content-length") ?? 0);
    if (declaredLength > this.settings.maxResponseBytes) throw new IntegrationTransportError("provider response exceeds admitted size", false, response.status);
    const text = await response.text();
    if (Buffer.byteLength(text, "utf8") > this.settings.maxResponseBytes) throw new IntegrationTransportError("provider response exceeds admitted size", false, response.status);
    let payload: unknown;
    try {
      payload = JSON.parse(text);
    } catch {
      throw new IntegrationTransportError("provider returned invalid JSON", false, response.status);
    }
    if (!payload || typeof payload !== "object") throw new IntegrationTransportError("provider returned an invalid result", false, response.status);
    const record = payload as Record<string, unknown>;
    const externalId = typeof record.externalId === "string" ? record.externalId : typeof record.id === "string" ? record.id : undefined;
    const status = record.status === "completed" ? "completed" : "accepted";
    if (!externalId || externalId.length > 300) throw new IntegrationTransportError("provider result lacks an external id", false, response.status);
    return { externalId, status };
  }
}
````

### FILE: `src/platform/integrations/registry.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-integrations-registry-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "0bf61c90c03db373713d35e28a6c799ec0bf3e6c6a53b32becefee24a2f591a8"
variables: []
secrets_allowed: false
```

````typescript
import type { IntegrationAdapter } from "./contract";

// Adapters enter this registry only after provider-specific contract and sandbox tests.
export const installedIntegrationAdapters: IntegrationAdapter[] = [];
````

### FILE: `src/platform/observability/events.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-observability-events-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "43736dcdf47517d0630f686b249a7bf55f4c7223c26f72a8090ddb072dc6f891"
variables: []
secrets_allowed: false
```

````typescript
export interface OperationalEvent {
  event: string;
  outcome: "success" | "failure";
  durationMs?: number;
  traceId?: string;
  details?: Record<string, unknown>;
}

export function emitOperationalEvent(value: OperationalEvent): void {
  console.info(JSON.stringify({ level: "info", timestamp: new Date().toISOString(), ...value }));
}
````

### FILE: `src/platform/outbox/dispatcher.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-outbox-dispatcher-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "506280a384ecfa18d57194f7903bb134148a04ad3d88d96136a7a208c33235db"
variables: []
secrets_allowed: false
```

````typescript
import type { Database } from "../db/types";

interface PendingEvent extends Record<string, unknown> {
  event_id: string;
  event_type: string;
  payload: Record<string, unknown>;
}

export type EventHandler = (event: PendingEvent) => Promise<void>;

export async function dispatchOutboxBatch(database: Database, handlers: Record<string, EventHandler>, limit = 50): Promise<number> {
  const pending = await database.query<PendingEvent>(`
    select event_id, event_type, payload
    from platform.outbox_event
    where published_at is null and available_at <= now()
    order by occurred_at
    limit $1
  `, [limit]);

  let dispatched = 0;
  for (const event of pending.rows) {
    const handler = handlers[event.event_type];
    if (!handler) continue;
    try {
      await handler(event);
      await database.query("update platform.outbox_event set published_at = now(), attempts = attempts + 1 where event_id = $1 and published_at is null", [event.event_id]);
      dispatched += 1;
    } catch {
      await database.query("update platform.outbox_event set attempts = attempts + 1, available_at = now() + interval '1 minute' where event_id = $1 and published_at is null", [event.event_id]);
    }
  }
  return dispatched;
}
````

### FILE: `src/platform/seed.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-seed-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "1183213195eb3c9cd6b197c9d5d01e78f5b4719faf76539517863e25d3d1c560"
variables: []
secrets_allowed: false
```

````typescript
import { createHash } from "node:crypto";
import type { BusinessConfig } from "./config/schema";
import type { Database } from "./db/types";

export const SEED_IDS = {
  hq: "00000000-0000-4000-8000-000000000001",
  franchise: "00000000-0000-4000-8000-000000000002",
  branch: "00000000-0000-4000-8000-000000000003",
  model: "00000000-0000-4000-8000-000000000101"
} as const;

export async function seed(database: Database, config: BusinessConfig): Promise<void> {
  const configJson = JSON.stringify(config);
  const hash = createHash("sha256").update(configJson).digest("hex");
  await database.query(`
    insert into platform.config_snapshot (config_hash, schema_version, business_id, config)
    values ($1, $2, $3, $4::jsonb)
    on conflict (config_hash) do nothing
  `, [hash, config.schemaVersion, config.business.id, configJson]);

  await database.query(`
    insert into org.organization (organization_id, parent_organization_id, organization_type, code, name, status)
    values
      ($1, null, 'hq', 'HQ', $4, 'active'),
      ($2, $1, 'franchise', 'FRANCHISE_NORTH', 'Franquicia Norte', 'active'),
      ($3, $2, 'branch', 'BRANCH_NORTH_01', 'Sucursal Norte 01', 'active')
    on conflict (organization_id) do nothing
  `, [SEED_IDS.hq, SEED_IDS.franchise, SEED_IDS.branch, config.business.name]);

  await database.query(`
    insert into catalog.model (model_id, slug, name, summary, status, market_code, custom_fields, published_at)
    values ($1, 'e1-city', 'E1 City', 'Movilidad eléctrica urbana, configurable por mercado.', 'published', $2, $3::jsonb, now())
    on conflict (model_id) do nothing
  `, [SEED_IDS.model, config.business.defaultMarket, JSON.stringify({ estimated_range_km: 120 })]);
}
````

### FILE: `src/platform/workflow/engine.test.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-workflow-engine-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "e5c6476d162d1c3522724e94342f68d41378a56e4f6f263bc170ac086908c201"
variables: []
secrets_allowed: false
```

````typescript
import { describe, expect, it } from "vitest";
import { loadBusinessConfig } from "../config/load";
import { transitionWorkflow } from "./engine";

describe("configuration-driven workflow", () => {
  it("allows a configured transition with its permission", async () => {
    const workflow = (await loadBusinessConfig({ bypassCache: true })).workflows.lead!;
    expect(transitionWorkflow(workflow, { current: "new", target: "assigned", permissions: new Set(["lead:assign"]) })).toBe("assigned");
  });

  it("rejects unconfigured transitions and missing permissions", async () => {
    const workflow = (await loadBusinessConfig({ bypassCache: true })).workflows.lead!;
    expect(() => transitionWorkflow(workflow, { current: "new", target: "converted", permissions: new Set(["*"]) })).toThrow(/not configured/);
    expect(() => transitionWorkflow(workflow, { current: "new", target: "assigned", permissions: new Set() })).toThrow(/permission lead:assign/);
  });
});
````

### FILE: `src/platform/workflow/engine.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:src-platform-workflow-engine-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "7d86b4ead011146bebb7c6069edb443be4c37f1db9976dbe88e777fdec06e6c7"
variables: []
secrets_allowed: false
```

````typescript
import type { WorkflowConfig } from "../config/schema";

export interface TransitionRequest {
  current: string;
  target: string;
  permissions: Set<string>;
}

export class WorkflowTransitionError extends Error {
  readonly code = "WORKFLOW_TRANSITION_REJECTED";
}

export function transitionWorkflow(workflow: WorkflowConfig, request: TransitionRequest): string {
  const transition = workflow.transitions.find((item) => item.from === request.current && item.to === request.target);
  if (!transition) throw new WorkflowTransitionError(`transition ${request.current} -> ${request.target} is not configured`);
  if (!request.permissions.has("*") && !request.permissions.has(transition.permission)) {
    throw new WorkflowTransitionError(`permission ${transition.permission} is required`);
  }
  return request.target;
}
````

### FILE: `SYSTEM_READINESS.md`

```yaml
block_id: "TS-ENTERPRISE-WEB:system-readiness-md:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "a5199725f832dd57ea82dcb3596f62f6e09807df99e34af53f178712bfef1514"
variables: []
secrets_allowed: false
```

````markdown
# Readiness del adapter web

Fecha de corte: 2026-08-24.

| Capacidad | Estado | Evidencia o bloqueo |
|---|---|---|
| contrato de configuración | READY para baseline | Zod + tests positivos/negativos |
| base local y PostgreSQL de producción | READY/CONDITIONED | PGlite probado; PostgreSQL externo requiere harness en entorno objetivo |
| catálogo público + lead | READY para desarrollo | build + transacción + idempotencia + audit + outbox tests |
| autorización organizacional | READY como kernel | unit tests; falta conectar identidad y probar E2E |
| workflows configurables | READY como kernel | unit tests; faltan servicios de dominio para todos los workflows |
| cliente/admin/fábrica | SHELL | identidad y journeys operativos pendientes |
| inventario/pedidos/pagos/service | SCHEMA/KERNEL | servicios, UI e invariantes verticales pendientes |
| marketplaces/ads/fiscal/logística | CONTRACT ONLY | adapters y sandboxes reales pendientes |
| observabilidad | BASELINE | eventos estructurados; collector, SLO y alertas pendientes |
| seguridad | BASELINE | fail-closed productivo, headers, rate limit; threat test/E2E/supply-chain pipeline pendientes |
| performance/low latency | OPEN | no hay workload, presupuesto ni benchmark del entorno objetivo |
| recovery | OPEN | backup/restore y RPO/RTO deben demostrarse con PostgreSQL objetivo |
| licencia del código propio | OPEN | el propietario debe elegir licencia; dependencias se reportan por SPDX |

Conclusión: existe un adapter web tangible sobre el cual configurar y ampliar una variante elegida por el blueprint. No es correcto usarlo como ERP/franquicia completo ni prometer producción universal. La siguiente promoción requiere corregir los invariantes auditados, completar los journeys seleccionados, identidad real, adaptadores necesarios, deployment, recovery y gates de evidencia.
````

### FILE: `THIRD_PARTY_NOTICES.md`

```yaml
block_id: "TS-ENTERPRISE-WEB:third-party-notices-md:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "a863b375d22d18c5f213bdf071a661893ae6ab3fcbbbd3775d8772e7c4cef236"
variables: []
secrets_allowed: false
```

````markdown
# Dependencias y licencias

Las dependencias directas están fijadas en `package.json` y el grafo completo en `pnpm-lock.yaml`.

| Dependencia directa | Versión | Licencia declarada por el paquete |
|---|---:|---|
| Next.js | 16.3.2 | MIT |
| React / React DOM | 19.2.8 | MIT |
| Zod | 4.4.3 | MIT |
| node-postgres | 8.23.0 | MIT |
| PGlite | 0.5.7 | Apache-2.0 |
| TypeScript | 7.0.2 | Apache-2.0 |
| Vitest | 4.1.11 | MIT |
| tsx | 4.23.12 | MIT |

Regenerar el inventario efectivo con `pnpm licenses:report`. El grafo transitivo observado también contiene Apache-2.0, MIT, ISC, BSD-3-Clause, 0BSD, CC-BY-4.0, MPL-2.0 y componentes binarios de imagen con términos adicionales. Este resumen no reemplaza conservar notices, revisar artefactos realmente distribuidos ni una revisión legal para el modo de entrega elegido.

No se asignó una licencia de distribución al código propio del adapter: esa decisión pertenece al propietario del repositorio.
````

### FILE: `tsconfig.json`

```yaml
block_id: "TS-ENTERPRISE-WEB:tsconfig-json:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "4c67dbbd61d014a203345b85086dc93ce3046260f229c733eafeeffe34808c0e"
variables: []
secrets_allowed: false
```

````json
{
  "compilerOptions": {
    "target": "ES2023",
    "lib": [
      "dom",
      "dom.iterable",
      "es2023"
    ],
    "allowJs": false,
    "skipLibCheck": true,
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "react-jsx",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": [
        "./src/*"
      ]
    }
  },
  "include": [
    "next-env.d.ts",
    "**/*.ts",
    "**/*.tsx",
    ".next/types/**/*.ts",
    ".next/dev/types/**/*.ts"
  ],
  "exclude": [
    "node_modules"
  ]
}
````

### FILE: `vitest.config.ts`

```yaml
block_id: "TS-ENTERPRISE-WEB:vitest-config-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local canonical Markdown pack"
license: "LicenseRef-Workspace-Owner"
sha256: "3a7bad10216bb4f319b026ca349a9f343bcc5e1f3b8a4b09892a4c9d91ef6a28"
variables: []
secrets_allowed: false
```

````typescript
import { defineConfig } from "vitest/config";
import { fileURLToPath } from "node:url";

export default defineConfig({
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
      "server-only": fileURLToPath(new URL("./src/test/server-only.ts", import.meta.url))
    }
  },
  test: {
    environment: "node",
    include: ["src/**/*.test.ts"],
    coverage: { reporter: ["text", "json", "html"] }
  }
});
````

## 6. Configuration surface

The exact non-secret surface is defined by `.env.example`, `config/business.example.json`, `config/readiness.json` and `src/platform/config/schema.ts`. Secrets must be injected by the target environment and must never be written into the business profile or Markdown.

## 7. Dependency bill

Exact direct and transitive pins are in `package.json` and `pnpm-lock.yaml`. Production license output on 2026-08-24 contained MIT, Apache-2.0, ISC, BSD-3-Clause, 0BSD, CC-BY-4.0 and an Apache-2.0 AND LGPL-3.0-or-later Sharp platform binary. Review redistribution obligations for the deployment artifact; do not infer approval from a successful build.

## 8. Apply order

1. Materialize every file exactly once into an empty directory.
2. Verify each SHA-256 before dependency installation.
3. Install Node.js 24 and pnpm compatible with the lockfile.
4. Run `pnpm install --frozen-lockfile`.
5. Run `pnpm verify`.
6. For production, satisfy `pnpm readiness:production`; never promote the PGlite fallback by accident.

## 9. Verification

V5 verification on 2026-08-24: canonical-source metadata, config and baseline readiness passed; TypeScript, 9 test files/22 tests, Next.js production build and productive dependency audit passed. `readiness:production` failed intentionally with nine unresolved production claims.

## 10. Reconstruction evidence

Estado vigente: `REBUILD_VERIFIED / CONDITIONED`, evidencia `TS-EW-20260824-V5`. V4 queda histórica.

- 64 archivos materializados desde Markdown con hashes verificados;
- frozen lockfile instaló 89 packages;
- config, readiness, TypeScript, 9 archivos/22 tests y build production: PASS;
- audit productivo: sin vulnerabilidades conocidas;
- production readiness: NOT_READY, rechazo esperado;
- incluye CI reproducible con GitHub Actions fijadas por SHA.

Las condiciones productivas por proyecto siguen siendo OIDC/JWKS real, PostgreSQL y roles, secretos, contratos sandbox, telemetry, deployment y recovery.
