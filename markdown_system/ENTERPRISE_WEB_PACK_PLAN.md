# Enterprise Web BFF Pack Plan

Este perfil materializa una unidad frontend/BFF autónoma del golden path TypeScript histórico. Incluye web pública configurable, catálogo y captación, portales protegidos cliente/operaciones/fábrica, OIDC server-side, gate de licencias, navegador Microsoft Playwright y calidad web Google Lighthouse en módulos separados. Requiere las APIs Go públicas y de consulta del perfil backend; no contiene dominio, SQL, migrations, persistencia ni workers.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md",
      "packId": "TS-GO-API-WEB-BRIDGE",
      "version": "0.7.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md",
      "packId": "TS-OIDC-PORTAL-ADAPTER",
      "version": "0.2.5",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md",
      "packId": "TS-FRANCHISE-JOURNEY-PORTALS",
      "version": "0.16.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/DEPENDENCY_LICENSE_EVIDENCE_CORE.md",
      "packId": "DEPENDENCY-LICENSE-EVIDENCE-CORE",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md",
      "packId": "MICROSOFT-PLAYWRIGHT-BROWSER-GATE",
      "version": "0.1.37",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE.md",
      "packId": "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE",
      "version": "0.1.5",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TS_CUSTOMER_SURVEY_PORTAL.md",
      "packId": "TS-CUSTOMER-SURVEY-PORTAL",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_PAYMENT_CHECKOUT_PORTAL.md",
      "packId": "TS-PAYMENT-CHECKOUT-PORTAL",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

La composición exacta produce 84 archivos sin colisiones: 37 del BFF, 10 del overlay OIDC, 14 del journey de franquicia/cliente, 3 del gate de licencias, 10 de Playwright y 10 de Lighthouse. Incluye administración de jornadas/ausencias/recursos, asignación/transición auditable, cancelación propia y recepción de entrega con serie/hash server-side; ninguna UI inventa recurrencia, política, customer subject ni digest de evidencia. Incluye Google SafeValues 1.2.0, CSP estricta con nonce fresco y cuatro regresiones con Google CSP Evaluator 1.1.8. `TS-ENTERPRISE-WEB` queda fuera: es un pack full-stack legado con admisión `CANDIDATE`, no foundation. Los secretos OIDC y de sesión se inyectan en el target. `acknowledgeConditions` no elimina IdP/edge/roles/backend reales, CSRF/origin y antiabuse, accesibilidad con tecnología asistiva, RUM/load, observabilidad, seguridad, rollout ni rollback; SafeValues y CSP Evaluator tienen alcance limitado y Lighthouse aporta únicamente mediciones automatizadas de laboratorio. La política nonce obliga render dinámico y debe probarse nuevamente frente al CDN/WAF/hosting reales.

V400: explicitly selected customer survey reference; feature flags remain disabled by default. Project terms, identity grants, invitation distribution and runtime acceptance are conditions. See reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.

V402 composition update: exact FX snapshot receipt owner selected with its host; connected WhatsApp/OIDC/social owners selected only where their full application dependencies are present. The shared BC MIT license has one composed owner. Local fixture and source gates are scoped in COMMUNICATIONS_RUNTIME_V402.md; this is not a production claim.
