# Franchise Serverless Pack Plan

Perfil vigente V402325: 66 packs / 997 archivos. ARCA conectado local y build portable; uso externo condicionado a credenciales.


Perfil vigente V402324: 66 packs / 965 archivos. Driver firmado local con dos builds independientes.


Perfil vigente V402322: 66 packs / 960 archivos. Avisos exactos de runtime; sin cambios funcionales ni nuevas dependencias.


Perfil vigente V402319: 66 packs / 936 archivos. El historial previo no gobierna este conteo.


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
      "version": "0.4.8",
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
      "version": "0.4.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_SUPPLY_FACTORY_INVENTORY_API.md",
      "packId": "GO-SUPPLY-FACTORY-INVENTORY-API",
      "version": "0.18.0",
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
      "version": "0.8.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_FULFILLMENT_SERVICE_FRANCHISE_API.md",
      "packId": "GO-FULFILLMENT-SERVICE-FRANCHISE-API",
      "version": "0.6.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ENTERPRISE_QUERY_API.md",
      "packId": "GO-ENTERPRISE-QUERY-API",
      "version": "0.3.0",
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
      "version": "0.1.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_ARCA_WSFE_GENERATED_CLIENT.md",
      "packId": "MICROSOFT-ARCA-WSFE-GENERATED-CLIENT",
      "version": "0.2.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_ARCA_WSAA_CREDENTIAL_CORE.md",
      "packId": "MICROSOFT-ARCA-WSAA-CREDENTIAL-CORE",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "arca/credentials/README.md",
        "arca/credentials/src/Elite.Arca.Credentials/Elite.Arca.Credentials.csproj",
        "arca/credentials/src/Elite.Arca.Credentials/packages.lock.json",
        "arca/credentials/src/Elite.Arca.Credentials/WsaaContracts.cs",
        "arca/credentials/src/Elite.Arca.Credentials/LoginTicketRequestFactory.cs",
        "arca/credentials/src/Elite.Arca.Credentials/CmsRequestSigner.cs",
        "arca/credentials/src/Elite.Arca.Credentials/CertificateStoreLoader.cs",
        "arca/credentials/src/Elite.Arca.Credentials/WsaaCredentialProvider.cs",
        "arca/credentials/integration-tests/Elite.Arca.Wsaa.Transport.Tests/Elite.Arca.Wsaa.Transport.Tests.csproj",
        "arca/credentials/integration-tests/Elite.Arca.Wsaa.Transport.Tests/Program.cs",
        "arca/credentials/integration-tests/Elite.Arca.Wsaa.Transport.Tests/packages.lock.json",
        "arca/credentials/integration/Elite.Arca.Wsaa.Transport/Elite.Arca.Wsaa.Transport.csproj",
        "arca/credentials/integration/Elite.Arca.Wsaa.Transport/GeneratedWsaaTransport.cs",
        "arca/credentials/integration/Elite.Arca.Wsaa.Transport/packages.lock.json",
        "arca/credentials/tests/Elite.Arca.Credentials.Tests/Elite.Arca.Credentials.Tests.csproj",
        "arca/credentials/tests/Elite.Arca.Credentials.Tests/Program.cs",
        "arca/credentials/tests/Elite.Arca.Credentials.Tests/packages.lock.json",
        "tools/test-arca-wsaa-credential-core.ps1"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_ARCA_WSFE_SOAP_ADAPTER.md",
      "packId": "MICROSOFT-ARCA-WSFE-SOAP-ADAPTER",
      "version": "0.3.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ARCA_FISCAL_ISSUANCE_API.md",
      "packId": "GO-ARCA-FISCAL-ISSUANCE-API",
      "version": "0.6.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/ARCA_WSFE_UDS_WORKER.md",
      "packId": "ARCA-WSFE-UDS-WORKER",
      "version": "0.3.1",
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
      "version": "1.22.3",
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
      "version": "0.1.9",
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
      "version": "0.1.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_ML_AI_FOUNDATION.md",
      "packId": "GO-ML-AI-FOUNDATION",
      "version": "0.2.1",
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
      "version": "0.15.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md",
      "packId": "TS-OIDC-PORTAL-ADAPTER",
      "version": "0.3.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md",
      "packId": "TS-FRANCHISE-JOURNEY-PORTALS",
      "version": "0.21.0",
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
      "version": "0.2.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TS_CUSTOMER_SURVEY_PORTAL.md",
      "packId": "TS-CUSTOMER-SURVEY-PORTAL",
      "version": "0.1.2",
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
      "version": "0.3.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_INITIAL_HANDOVER_API.md",
      "packId": "GO-INITIAL-HANDOVER-API",
      "version": "0.4.0",
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
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_EXACT_FX_SNAPSHOT_ACCOUNTING.md",
      "packId": "GO-EXACT-FX-SNAPSHOT-ACCOUNTING",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "cmd/electromobility-api/fx.go",
        "cmd/electromobility-api/fx_test.go",
        "config/fx/reference-profile.json",
        "config/fx/reference-rates.json",
        "db/migrations/0061_fx_conversion_receipt.down.sql",
        "db/migrations/0061_fx_conversion_receipt.up.sql",
        "db/migrations/0065_fx_journal_receipt.down.sql",
        "db/migrations/0065_fx_journal_receipt.up.sql",
        "docs/fx-conversion-runtime.md",
        "docs/fx-journal-runtime.md",
        "docs/provenance/BC_FX_DERIVATION.md",
        "docs/provenance/BC_FX_SOURCE_LOCK.json",
        "docs/provenance/BC_FX_THIRD_PARTY_NOTICES.md",
        "internal/accounting/fx.go",
        "internal/accounting/fx_journal.go",
        "internal/bcfx/exchange.go",
        "internal/bcfx/exchange_test.go",
        "internal/bcfx/rounding.go",
        "internal/bcfx/selection.go",
        "internal/bcfx/snapshot.go",
        "internal/bcfx/snapshot_test.go",
        "internal/bcfx/testdata/rounding-oracle.json",
        "internal/platform/httpapi/fx_conversion.go",
        "internal/platform/httpapi/fx_journal.go",
        "internal/platform/postgres/fx_connected_integration_test.go",
        "internal/platform/postgres/fx_conversion.go",
        "internal/platform/postgres/fx_journal.go",
        "internal/platform/postgres/fx_journal_integration_test.go",
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
    },
    {
      "path": "implementation_packs/TYPESCRIPT_PUBLISHED_CATALOG_STOREFRONT.md",
      "packId": "TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_CONNECTED_TRAINING_PORTAL.md",
      "packId": "TS-CONNECTED-TRAINING-PORTAL",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_CATALOG_AUTHORING_PORTAL.md",
      "packId": "TS-CATALOG-AUTHORING-PORTAL",
      "version": "0.1.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_SERIAL_SUPPLY_PORTAL.md",
      "packId": "TS-SERIAL-SUPPLY-PORTAL",
      "version": "0.1.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_WARRANTY_ROLE_PORTAL.md",
      "packId": "TS-WARRANTY-ROLE-PORTAL",
      "version": "0.1.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_NETWORK_ROLE_PORTAL.md",
      "packId": "TS-NETWORK-ROLE-PORTAL",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/TYPESCRIPT_CONNECTED_HELP_CMS_PORTAL.md",
      "packId": "TS-CONNECTED-HELP-CMS-PORTAL",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_AWS_TEXTRACT_DOCUMENT_RUNTIME.md",
      "packId": "GO-AWS-TEXTRACT-DOCUMENT-RUNTIME",
      "version": "0.2.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_HTTP_METRICS_REFERENCE.md",
      "packId": "GO-HTTP-METRICS-REFERENCE",
      "version": "0.1.16",
      "acknowledgeConditions": true,
      "files": [
        "reference_http_metrics/go.mod",
        "reference_http_metrics/go.sum",
        "reference_http_metrics/core/metrics.go"
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

V402 FX journal delta: FX0.2.0 and accounting0.1.2 connect conversion receipts to the existing draft/post/reversal lifecycle; added seven selected files. See FX_JOURNAL_CONNECTION_V402.md.

V402 stored-value connected delta: selected shared funding/handover owners and optional source-calculator/approved-tender pair; source/license bytes exact via compositor0.3.0. See STORED_VALUE_SOURCE_CLOSURE_V402.md; final reconstruction/admission pending.

V402 warranty connected delta: shared service/stock/approval transaction owners and sold terms through J4 claim. Full franchise/HTTP reference explicitly select date and connected packs with one MIT license owner. See WARRANTY_CONNECTED_RELEASE_V402.md for exact reconstruction and scoped admission. Earlier stored-value pending publication note is superseded by STORED_VALUE_CONNECTED_RELEASE_V402.md.

V402 optional serial identifier correction: inventory0.17.1/migration71 preserves unique known identifiers and allows multiple missing optional values. No previous migration/data changed; SERIAL_OPTIONAL_IDENTIFIERS_V402.md.

V402 connected serial J2 delta: GO-CONNECTED-SERIAL-SUPPLY0.1.0 and existing transaction/approval/optional host owners. See SERIAL_SUPPLY_CONNECTED_RELEASE_V402.md. Role UI and overall infrastructure closure remain separate.

V402 J3 connected catalog delta: exact source owners and approved storefront. See CATALOG_CONNECTED_RELEASE_V402.md. Role authoring T2804/provider SDK mappings T2805 remain separate.

V402 T2804 connected training: opt-in same-release public help, durable participation and human evaluation; existing identity/approval/outbox owners. TRAINING_CONNECTED_RELEASE_V402.md. No permission grants.

V402 T2804 catalog role source/PNG/draft/review/publication and recovery; original model/price services and catalog ownership preserved. CATALOG_ROLE_AUTHORING_RELEASE_V402.md.

V402 T2804 connected supply role order/quantities/factory/receiving/quality/recovery; original owners retained. SUPPLY_ROLE_RELEASE_V402.md.

V402 T2804 warranty role terms/consent/activation/repair/quality/customer acceptance/factory reconciliation with durable GET recovery; original owners retained. WARRANTY_ROLE_RELEASE_V402.md.

V402 T2804 network role organization/agreement/branch lifecycle and durable GET recovery; original owner rules retained. NETWORK_ROLE_RELEASE_V402.md.

V402 T2804 help CMS and21same-release guides, training revision2, original kernel retained. HELP_CMS_RELEASE_V402.md.

V402 T2804 source-connected role metrics and bounded generic BFF. ROLE_METRICS_RELEASE_V402.md/json.

V402 T2804 private locale/source-bound guide display. PRIVATE_LOCALE_RELEASE_V402.md/json.

V402 T2805 connected marketplace mutation slice. MARKETPLACE_MUTATION_RELEASE_V402.md/json; initial publication/media and other ordered work remain pending.

V402 Google Merchant connected local claim; MERCHANT_CONNECTED_RELEASE_V402.md/json. Other T2805 mappings/comms remain open.

V402 portal lifecycle connected; IDENTITY_PORTAL_RELEASE_V402.md/json. T2803 source/SCA and J5 administration remain open.

V402 J5 explicit IdP access review and narrow bootstrap; IDENTITY_J5_RELEASE_V402.md/json.

V402 composition security delta: fixed Go graph floor and bounded DevSkim failure diagnostics. COMPOSITION_SECURITY_RELEASE_V402.md/json.

V402 / T2806: documento referencia conectado; SDK original, gate de seguridad reutilizado, review humana y commit durable. DOCUMENT_REFERENCE_RELEASE_V402.md/json. FIXTURE explícito, no producción ni OCR probado.

V402315: T2807 connected domain/contact/runtime/evals and inherited historical execution proof. AI_CONNECTED_REFERENCE_RELEASE_V402.md/json; local/synthetic only.

V402316: local native build/deploy/migration/rollback entrypoints. LOCAL_REFERENCE_DELIVERY_V402.md/json; local/synthetic scope only.

V402317: actual API+Next local host, official closed metrics, alert/repair/concurrent replay/WAL recovery. LOCAL_REFERENCE_OPERATIONS_V402.md/json.

V402319: selected PostHog MIT NPS owner, explicit adaptation/typed caller and preserved source notices.
