# Franchise Serverless Pack Plan

Perfil de despliegue **serverless-first** (Cloud Run + servicios administrados) para
mínimo costo por uso: la capa de cómputo escala a cero cuando no hay tráfico y
sólo se factura lo consumido. Reutiliza los packs existentes; el adapter de
deploy cloud (`gcloud`) queda `CONDITIONED` (requiere cuenta/proyecto).

## Principio de costo

- **Cómputo**: Cloud Run (contenedor `CONTAINER-PACKAGING-CORE`, imagen mínima,
  non-root, read-only) con `min-instances=0` → **scale-to-zero**, factura por
  request/segundo de CPU/memoria, no por instancia ociosa.
- **Datos**: PostgreSQL administrado (Cloud SQL) con parada programada fuera de
  horario si el tráfico lo permite; backup/restore con `PG-BACKUP-RESTORE-CORE`.
- **Objetos/cola**: Cloud Storage + Pub/Sub administrados (por uso).
- **LLM**: siempre proveedor externo por uso (`GO-LLM-OPENAI-ADAPTER` +
  `GO-LLM-ECONOMY-CORE`), nunca GPU propia encendida (evita `GPU-ACCEL`).
- **Gobernanza**: `GO-FINOPS-CORE` — tags obligatorios, inventario, budget por
  tenant, audit de drift/leak y teardown que no deja huérfanos.

## Composición (reutiliza packs existentes)

Este perfil es materializable y construible como base: contiene los 39 packs únicos del perfil histórico más seis capacidades de costo/canales. No es autónomo ni omnicanal de extremo a extremo: omite tools, binding, aprobación, app wiring y el backbone durable de leads. `acknowledgeConditions` permite la composición deliberada, pero no demuestra esas conexiones, cuentas cloud, credenciales, billing, políticas de canal ni producción. Véase `POST_DEEPSEEK_FRANCHISE_REAUDIT_2026-09-04.md`.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md",
      "packId": "PG-TX-FOUNDATION",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md",
      "packId": "GO-ENTERPRISE-BACKEND",
      "version": "0.4.5",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/ELECTROMOBILITY_FRANCHISE_MODULES.md",
      "packId": "ELECTROMOBILITY-FRANCHISE-MODULES",
      "version": "0.1.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_RELIABLE_ASYNC_WORKERS.md",
      "packId": "GO-RELIABLE-ASYNC-WORKERS",
      "version": "0.3.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md",
      "packId": "GO-ELECTROMOBILITY-PUBLIC-CRM-API",
      "version": "0.2.3",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_SUPPLY_FACTORY_INVENTORY_API.md",
      "packId": "GO-SUPPLY-FACTORY-INVENTORY-API",
      "version": "0.16.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_BC_EXACT_AMOUNT_ADAPTER.md",
      "packId": "GO-BC-EXACT-AMOUNT-ADAPTER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_BC_SALES_CONTRACT_ADAPTER.md",
      "packId": "GO-BC-SALES-CONTRACT-ADAPTER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md",
      "packId": "GO-COMMERCE-PRICING-PAYMENT-API",
      "version": "0.6.4",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_FULFILLMENT_SERVICE_FRANCHISE_API.md",
      "packId": "GO-FULFILLMENT-SERVICE-FRANCHISE-API",
      "version": "0.4.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ENTERPRISE_QUERY_API.md",
      "packId": "GO-ENTERPRISE-QUERY-API",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
      "packId": "GO-FRANCHISE-CUSTOMER-JOURNEY-API",
      "version": "0.10.21",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_RETURN_EFFECT_EXECUTION_WORKER.md",
      "packId": "GO-RETURN-EFFECT-EXECUTION-WORKER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_OFFICIAL_RETURN_REFUND_WORKER.md",
      "packId": "GO-OFFICIAL-RETURN-REFUND-WORKER",
      "version": "0.1.6",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_RETURN_EXCHANGE_FULFILLMENT_WORKER.md",
      "packId": "GO-RETURN-EXCHANGE-FULFILLMENT-WORKER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_RETURN_ACCOUNTING_REVERSAL_WORKER.md",
      "packId": "GO-RETURN-ACCOUNTING-REVERSAL-WORKER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_RETURN_FISCAL_CREDIT_NOTE_WORKER.md",
      "packId": "GO-RETURN-FISCAL-CREDIT-NOTE-WORKER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_FRANCHISE_ROYALTY_SETTLEMENT_API.md",
      "packId": "GO-FRANCHISE-ROYALTY-SETTLEMENT-API",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ENTERPRISE_ACCOUNTING_LEDGER_API.md",
      "packId": "GO-ENTERPRISE-ACCOUNTING-LEDGER-API",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_ARCA_WSFE_GENERATED_CLIENT.md",
      "packId": "MICROSOFT-ARCA-WSFE-GENERATED-CLIENT",
      "version": "0.2.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_ARCA_WSAA_CREDENTIAL_CORE.md",
      "packId": "MICROSOFT-ARCA-WSAA-CREDENTIAL-CORE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "arca/credentials/README.md",
        "arca/credentials/src/Elite.Arca.Credentials/Elite.Arca.Credentials.csproj",
        "arca/credentials/src/Elite.Arca.Credentials/packages.lock.json",
        "arca/credentials/src/Elite.Arca.Credentials/WsaaContracts.cs",
        "arca/credentials/src/Elite.Arca.Credentials/LoginTicketRequestFactory.cs",
        "arca/credentials/src/Elite.Arca.Credentials/CmsRequestSigner.cs",
        "arca/credentials/src/Elite.Arca.Credentials/CertificateStoreLoader.cs",
        "arca/credentials/src/Elite.Arca.Credentials/WsaaCredentialProvider.cs"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_ARCA_WSFE_SOAP_ADAPTER.md",
      "packId": "MICROSOFT-ARCA-WSFE-SOAP-ADAPTER",
      "version": "0.3.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ARCA_FISCAL_ISSUANCE_API.md",
      "packId": "GO-ARCA-FISCAL-ISSUANCE-API",
      "version": "0.6.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/ARCA_WSFE_UDS_WORKER.md",
      "packId": "ARCA-WSFE-UDS-WORKER",
      "version": "0.3.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_PROVIDER_INTEGRATION_CORE.md",
      "packId": "GO-PROVIDER-INTEGRATION-CORE",
      "version": "0.1.3",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ELECTROMOBILITY_APPLICATION.md",
      "packId": "GO-ELECTROMOBILITY-APPLICATION",
      "version": "1.11.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md",
      "packId": "SECURE-OPS-DELIVERY-CORE",
      "version": "1.1.4",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/PORTABLE_CI_QUALITY_GATE_RUNNER.md",
      "packId": "PORTABLE-CI-GATE-RUNNER",
      "version": "0.1.1",
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
      "path": "implementation_packs/POSTGRES_BACKUP_RESTORE_CORE.md",
      "packId": "PG-BACKUP-RESTORE-CORE",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/CONTAINER_PACKAGING_CORE.md",
      "packId": "CONTAINER-PACKAGING-CORE",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ML_AI_FOUNDATION.md",
      "packId": "GO-ML-AI-FOUNDATION",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_SEARCH_CORE.md",
      "packId": "GO-SEARCH-CORE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_CACHE_CORE.md",
      "packId": "GO-CACHE-CORE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_CONVERSATIONAL_AGENT.md",
      "packId": "GO-CONVERSATIONAL-AGENT",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_DATA_ANALYTICS_CORE.md",
      "packId": "GO-DATA-ANALYTICS-CORE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
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
      "path": "implementation_packs/GO_LLM_OPENAI_ADAPTER.md",
      "packId": "GO-LLM-OPENAI-ADAPTER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_LLM_ECONOMY_CORE.md",
      "packId": "GO-LLM-ECONOMY-CORE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_CHANNELS_CORE.md",
      "packId": "GO-CHANNELS-CORE",
      "version": "0.4.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_MESSAGING_CHANNEL_ADAPTERS.md",
      "packId": "GO-MESSAGING-CHANNEL-ADAPTERS",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_SMS_CHANNEL_ADAPTER.md",
      "packId": "GO-SMS-CHANNEL-ADAPTER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_FINOPS_CORE.md",
      "packId": "GO-FINOPS-CORE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_CUSTOMER_SURVEY_API.md",
      "packId": "GO-CUSTOMER-SURVEY-API",
      "version": "0.1.1",
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
      "path": "implementation_packs/GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS.md",
      "packId": "GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS",
      "version": "0.3.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_PAYMENT_CHECKOUT_RUNTIME.md",
      "packId": "GO-PAYMENT-CHECKOUT-RUNTIME",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_INITIAL_HANDOVER_API.md",
      "packId": "GO-INITIAL-HANDOVER-API",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_BUSINESS_POLICY_PROFILE.md",
      "packId": "GO-BUSINESS-POLICY-PROFILE",
      "version": "0.1.0",
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
    },
    {
      "path": "implementation_packs/GO_EXACT_FX_SNAPSHOT_ACCOUNTING.md",
      "packId": "GO-EXACT-FX-SNAPSHOT-ACCOUNTING",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "cmd/electromobility-api/fx.go",
        "cmd/electromobility-api/fx_test.go",
        "config/fx/reference-profile.json",
        "config/fx/reference-rates.json",
        "db/migrations/0061_fx_conversion_receipt.down.sql",
        "db/migrations/0061_fx_conversion_receipt.up.sql",
        "docs/fx-conversion-runtime.md",
        "docs/provenance/BC_FX_DERIVATION.md",
        "docs/provenance/BC_FX_SOURCE_LOCK.json",
        "docs/provenance/BC_FX_THIRD_PARTY_NOTICES.md",
        "internal/accounting/fx.go",
        "internal/bcfx/exchange.go",
        "internal/bcfx/exchange_test.go",
        "internal/bcfx/rounding.go",
        "internal/bcfx/selection.go",
        "internal/bcfx/snapshot.go",
        "internal/bcfx/snapshot_test.go",
        "internal/bcfx/testdata/rounding-oracle.json",
        "internal/platform/httpapi/fx_conversion.go",
        "internal/platform/postgres/fx_connected_integration_test.go",
        "internal/platform/postgres/fx_conversion.go",
        "tools/generate_fx_rounding_oracle.py"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_PG_OUTBOUND_DELIVERY_FENCE.md",
      "packId": "GO-PG-OUTBOUND-DELIVERY-FENCE",
      "version": "0.2.1",
      "acknowledgeConditions": true,
      "files": [
        "db/migrations/0048_outbound_delivery_fence.down.sql",
        "db/migrations/0048_outbound_delivery_fence.up.sql",
        "db/tests/0048_outbound_delivery_fence.test.sql",
        "internal/channels/durable_registry.go",
        "internal/channels/durable_registry_test.go",
        "internal/outbounddelivery/delivery.go",
        "internal/outbounddelivery/delivery_test.go",
        "internal/platform/postgres/outbound_delivery.go"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_PG_CONTACT_CHANNEL_IDENTITY.md",
      "packId": "GO-PG-CONTACT-CHANNEL-IDENTITY",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "internal/contactidentity/identity.go",
        "internal/contactidentity/identity_test.go"
      ],
      "variables": {}
    }
  ]
}
```

| Pack | Rol en serverless |
|---|---|
| `GO-ENTERPRISE-BACKEND` (29 packs) | base backend empresarial; no sistema integral por sí solo |
| `GO-CONVERSATIONAL-AGENT` + `GO-LLM-OPENAI-ADAPTER` + `GO-LLM-ECONOMY-CORE` | chatbot por uso |
| `GO-CHANNELS-CORE` + `GO-MESSAGING-CHANNEL-ADAPTERS` + `GO-SMS-CHANNEL-ADAPTER` | canales |
| `GO-FINOPS-CORE` | tags/budget/teardown |
| `CONTAINER-PACKAGING-CORE` | imagen mínima para Cloud Run |
| `PORTABLE-CI-QUALITY-GATE-RUNNER` | CI |
| `GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE` + `MICROSOFT_PLAYWRIGHT_BROWSER_GATE` | web |

## Lo que queda `CONDITIONED`

- Adapter de deploy real (`gcloud run deploy`, Terraform/Config Connector) y su
  cuenta/proyecto/billing: runtime del proyecto.
- Escala mínima >0 si un path crítico no tolera cold start (trade-off explícito
  costo↔latencia, no asumido).
- Idle shutdown de Cloud SQL fuera de horario: política del proyecto.

## Gate de cierre (anti "que no se escape nada")

1. `GO-FINOPS-CORE`: todo recurso etiquetado + inventariado + budget por tenant.
2. `Audit` limpio (cero untagged/drift/leak) antes de teardown.
3. Teardown remueve cada recurso; nada queda facturando en segundo plano.

V400: explicitly selected customer survey reference; feature flags remain disabled by default. Project terms, identity grants, invitation distribution and runtime acceptance are conditions. See reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.

V402 composition update: exact FX snapshot receipt owner selected with its host; connected WhatsApp/OIDC/social owners selected only where their full application dependencies are present. The shared BC MIT license has one composed owner. Local fixture and source gates are scoped in COMMUNICATIONS_RUNTIME_V402.md; this is not a production claim.

V402 dependency closure: select8fence outputs and2contact-digest outputs for payment infrastructure. The separate omnichannel/leadstream integration test remains in the full franchise profile, because this independent profile does not select its leadstream owner. Payment integration tests remain selected. Compile test packages explicitly; exact reconstruction alone is not compilation.
