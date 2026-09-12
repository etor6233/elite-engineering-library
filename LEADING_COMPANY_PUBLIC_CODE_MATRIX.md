# Leading Company Public Code Matrix

> **Corte:** 2026-08-26.  
> **Propósito:** seleccionar evidencia pública de compañías líderes por escala, valor, alcance o excelencia técnica para construir sistemas completos.  
> **Regla:** la compañía no es el argumento; el repositorio, la revisión, la licencia, la documentación, los tests y el fit sí lo son.

## 1. Cómo leer este mapa

Esta matriz no intenta copiar la plataforma interna completa de ninguna empresa. Las compañías publican piezas, no necesariamente el conjunto que opera su negocio. Cada fuente sólo respalda el claim estrecho indicado.

Estados:

- `ELITE_REFERENCE`: claim ya auditado en los manuales existentes;
- `CANDIDATE`: licencia y valor inicial verificados; falta G3–G7;
- `CONDITIONED`: licencia, términos, topología o coste exigen decisión explícita;
- `INTEGRATION_ONLY`: SDK/contrato válido sólo para consumir el servicio del proveedor;
- `EXCLUDED`: no oficial, archivado, educativo o no apto como baseline.

## 2. Matriz por compañía

| Compañía/ecosistema | Fuente pública fijada o autoridad | Licencia observada | Claim útil | Estado y uso permitido |
|---|---|---|---|---|
| Apple | [`apple/swift@7110e5a`](https://github.com/apple/swift/tree/7110e5aa544ab86e63837394eccf9d59a1c7f161) | Apache-2.0 | lenguaje, compiler/runtime y toolchain para cliente nativo | `ELITE_REFERENCE` sólo para claims ya auditados; seleccionar por producto/plataforma |
| Apple | [`apple/swift-nio@385c1a7`](https://github.com/apple/swift-nio/tree/385c1a7d48ebe018bcde205f8bd4506403dc889a) | Apache-2.0 | networking event-driven, backpressure, HTTP/TLS ecosystem | `CANDIDATE`; no usar NIO directamente si las APIs de plataforma resuelven el journey |
| Microsoft | TypeScript, VS Code y toolchain ya auditados | Apache-2.0/MIT por proyecto | contratos estáticos, tooling, language service y extensibilidad | `ELITE_REFERENCE` para claims existentes |
| Microsoft | [`microsoft/BCApps@2eae56d`](https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749) | MIT para source; runtime/servicio bajo términos propios | código real del producto Business Central: finanzas, compras, inventario/costing, fabricación, ventas, warehouse, service, first-party apps y tooling | `CONDITIONAL_PLATFORM`; commit firmado, 36.673 AL, Inventory 958, Warehouse 355 y 22 tests AL SCM-Reservation inventariados; snapshot de desarrollo que exige Business Central runtime/licencia/container y no prueba localización argentina |
| Microsoft | [`microsoft/fluentui@61f75c1`](https://github.com/microsoft/fluentui/tree/61f75c147e2d7fb42d5ea3395c813e06d877dd19) | MIT; fonts/icons con términos separados | componentes y accesibilidad para aplicaciones empresariales | `CONDITIONED`; código y assets se auditan por separado |
| Microsoft | [`microsoft/playwright@12b611d`](https://github.com/microsoft/playwright/tree/12b611da20d21db663ee0c6be399b1bb854dead0) | Apache-2.0 | E2E aislado sobre Chromium/Firefox/WebKit | `CANDIDATE`; base preferida para journeys web críticos |
| Microsoft/CNCF | [`dapr/dapr@151e08c`](https://github.com/dapr/dapr/tree/151e08cf63480b467dd55d8258ac7871ccf5d97e) | Apache-2.0 | service invocation, pub/sub, secrets, state, workflow y observabilidad mediante sidecar | `CONDITIONED`; no introducir hasta demostrar necesidad distribuida y coste operacional |
| Alphabet/Google | Chromium, Android/AOSP, Bazel, gRPC y Kubernetes ya auditados | BSD/Apache-2.0 por proyecto | browser, mobile, contracts, builds y control planes | `ELITE_REFERENCE` para claims existentes |
| Alphabet/Google | `GoogleCloudPlatform/microservices-demo` | Apache-2.0 | contratos gRPC, instrumentación, manifests y load generator | `CONDITIONED`; demo, no bounded contexts ni topología de producción |
| Alphabet/Google Firebase | [`firebase/firebase-admin-go@eebb06f`](https://github.com/firebase/firebase-admin-go/tree/eebb06f2a643fbb59b1cb262874a943584475128) | Apache-2.0 | SDK oficial server-side para FCM, incluidos `Send` y `SendDryRun` | `PINNED_CANDIDATE`; v4.21.0/source/licencia fijados, commit unsigned; pack/graph/cuenta/token policy/evidence pendientes |
| Meta | React, React Native, Buck2, PyTorch y RocksDB ya auditados | MIT/Apache/BSD por proyecto | UI, mobile, builds, AI y storage embebido | `ELITE_REFERENCE` para los claims documentados; nunca implica copiar la plataforma Meta |
| Amazon/AWS/CNCF | [`cedar-policy/cedar@b8fc36d`](https://github.com/cedar-policy/cedar/tree/b8fc36d497501d97b8be977d16ffed94b76a89b5) | Apache-2.0 | policy-as-code, RBAC/ABAC y validación formal/differential testing | `CANDIDATE`; alternativa local/embebida a OpenFGA, no mezclar dos engines sin necesidad |
| Amazon/AWS | [`smithy-lang/smithy@0f73231`](https://github.com/smithy-lang/smithy/tree/0f7323128b0606a1b94b1ac482c94d3800a22708) | Apache-2.0 | IDL, model validation y generación de clientes/servidores/docs | `CANDIDATE`; OpenAPI/Protobuf pueden ser suficientes |
| Amazon/AWS | [`firecracker-microvm/firecracker@40f2ed2`](https://github.com/firecracker-microvm/firecracker/tree/40f2ed2d3fd5896e378cf2fe07717efab66f1167) | Apache-2.0 con terceros/notices | aislamiento microVM usado para Lambda/Fargate | `CONDITIONED`; referencia de aislamiento, no dependencia de una app empresarial normal |
| Amazon/OpenSearch Foundation | [`opensearch-project/OpenSearch@82fbeca`](https://github.com/opensearch-project/OpenSearch/tree/82fbecab5483fe5f4a6982cbaf633e82ce521023) | Apache-2.0 | búsqueda distribuida, indexing y analytics | `CANDIDATE`; PostgreSQL search primero, cluster sólo por evidencia |
| Amazon Marketplace | [`amzn/selling-partner-api-models`](https://github.com/amzn/selling-partner-api-models) | Apache-2.0 | contratos OpenAPI de SP-API | `INTEGRATION_ONLY`; acuerdos, DPP, autorización y scopes siguen gobernando |
| NVIDIA | CUDA/CUTLASS/NCCL/TensorRT ya auditados | licencias por componente | GPU kernels, collectives e inference | `ELITE_REFERENCE` sólo para sistemas que necesiten GPU |
| Microsoft/Google/AWS/NVIDIA/Oracle | source lock documental exacto en `OFFICIAL_UPSTREAM_ACQUISITION_CORE.md` | MIT/Apache-2.0/UPL y términos de servicio por fuente | OCR, layout, tablas, invoice/receipt prebuilt, custom extraction/classification, revisión y SDKs oficiales | `CONDITIONED/INTEGRATION_ONLY/SAMPLE_ONLY` por fuente; usar `OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md`, nunca prometer exactitud sin corpus/ground truth |
| Oracle | Oracle JET, Helidon, GraalVM, Coherence y provider OCI | UPL/Apache/GPL+CE/BSD/MPL por componente | frontend, Java backend, native runtime, data grid e IaC | estados detallados en `PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md`; nunca tratar la familia como una licencia única |
| Tesla | Fleet Telemetry, Vehicle Command, Camera Kit y Fixed Containers | Apache-2.0/MIT | telemetría/commands vehiculares, mobile scanning y C++ fixed-capacity | `CANDIDATE` o `LICENSE_VERIFIED`; especialmente material para movilidad eléctrica |
| xAI/X | x-algorithm, grok-1, grok-build, xai-proto y SDK | Apache-2.0; AGPL/restrictiva en otros repos | ranking, inference, agent runtime y contratos | estados detallados en `AI_ENGINEERING_MASTER_MAP.md`; separar cada repositorio/licencia |
| SpaceX/Starlink | avisos/componentes OSS entregados por obligación | por componente | evidencia de compliance de dependencias | `EXCLUDED` como arquitectura SpaceX: no existe release oficial verificado del sistema propietario completo |
| Cloudflare | [`cloudflare/workerd@68459b7`](https://github.com/cloudflare/workerd/tree/68459b7e093084b14a35198f992f9196fc5a52ad) | Apache-2.0 | runtime JavaScript/Wasm server-first y programmable proxy | `CANDIDATE`; útil para edge/isolated functions, no default de backend |
| Cloudflare | [`cloudflare/pingora@0046038`](https://github.com/cloudflare/pingora/tree/0046038bd402bc82912da862dadf9a479f31e9f1) | Apache-2.0 | proxies HTTP/gRPC/WebSocket, TLS, LB, failover y reload | `CANDIDATE`; sólo para data paths que requieran proxy programable propio |
| Shopify | [`Shopify/hydrogen@2a2738b`](https://github.com/Shopify/hydrogen/tree/2a2738ba20487ccc07006815fe40e93b24cb5f08) | MIT | storefront React Router, codegen y contratos Shopify | `INTEGRATION_ONLY` para Shopify; recetas pueden informar UX pero no dominio agnóstico |
| Stripe | [`stripe/stripe-node@6639421`](https://github.com/stripe/stripe-node/tree/663942164d3876f04bb26617b88f498ba785e580) | MIT | cliente server-side, idempotency metadata, webhook signature y mocks | `INTEGRATION_ONLY`; PCI, API version, raw body y reconciliation son obligatorios |
| Mercado Pago | [`mercadopago/sdk-nodejs`](https://github.com/mercadopago/sdk-nodejs) | MIT | SDK oficial de pagos | `INTEGRATION_ONLY`; tokenización, webhooks, idempotencia y país/producto se fijan por proyecto |
| Mercado Libre | SDKs oficiales históricos | Apache-2.0 por SDK, pero archivados/deprecados | ejemplos de OAuth/client | `EXCLUDED` como dependencia; integrar la API vigente desde contratos/docs y un adapter propio |
| Google Ads | [`googleads/google-ads-java`](https://github.com/googleads/google-ads-java) | Apache-2.0 + terceros | cliente oficial, credentials y versiones generadas | `INTEGRATION_ONLY`; API deprecations, developer token, OAuth y account hierarchy requieren runbook |
| Meta Ads | [`facebook/facebook-python-business-sdk`](https://github.com/facebook/facebook-python-business-sdk) y [`facebook/facebook-nodejs-business-sdk`](https://github.com/facebook/facebook-nodejs-business-sdk) | Platform License limitada a APIs Meta | Marketing/Business/Pages/Instagram SDKs; Python Ads Insights read-only ya materializado | `CONDITIONED/INTEGRATION_ONLY`; no es licencia open source general, queda sujeta a Platform Policy y requiere cuenta/permisos/reconciliación |
| Meta WhatsApp | [`fbsamples/whatsapp-api-examples@de70ee90`](https://github.com/fbsamples/whatsapp-api-examples/tree/de70ee908a67026e642aaee3703d20464e2a9466) | Platform License limitada a APIs Meta | ejemplos oficiales template/webhook; tres códigos exactos y adaptación segura materializados | `CONDITIONED/INTEGRATION_ONLY`; SDK Node archivado excluido y bug sample retenido; exige terms/version/template/consent/inbox/reconciliación |
| TikTok/ByteDance | [`tiktok/tiktok-business-api-sdk@f809c396`](https://github.com/tiktok/tiktok-business-api-sdk/tree/f809c396520df2d7b201a9ccc5378d822b728ed3) | MIT | Business API generado; Python integrated reporting read-only ya materializado desde commit firmado | `CONDITIONED/INTEGRATION_ONLY`; sin release/tag/wheel, conflicto metadata/changelog y límite de paths Windows; exige source receipt, cuenta/scope/query y reconciliación |
| Linux Foundation/CNCF | Kubernetes, OpenTelemetry, Prometheus, Envoy, Valkey, Cedar | Apache/BSD por proyecto | primitives neutrales y gobernanza multi-vendor | preferibles cuando reducen lock-in, pero igualmente pasan el decision gate |
| Linux Foundation | [`valkey-io/valkey@94f61e4`](https://github.com/valkey-io/valkey/tree/94f61e4491e6c4422cb85d1a78d0437213b36fc0) | BSD-3-Clause | cache/ephemeral realtime data structures | `CANDIDATE`; nunca source of truth ni distributed lock sin proof específico |

## 3. Lo que se adopta como baseline

Para una plataforma empresarial nueva, la composición inicial preferida es deliberadamente menor que la matriz:

```text
web UI adapter selected per project (React-family is one candidate, not the foundation)
versioned HTTP API (OpenAPI) and events
compiled modular application core selected after blueprint (compare Go/JVM; Rust when justified)
PostgreSQL source of truth
OIDC provider + one authorization model
transactional outbox + bounded workers
object storage through an adapter
OpenTelemetry instrumentation
Playwright vertical journeys
immutable container artifact + IaC + staged release
```

La foundation y el dominio son stack-neutral. El adapter web puede usar TypeScript cuando sea la mejor elección para navegador/SSR, pero no impone el runtime del backend ni gobierna la arquitectura del sistema.

Todo lo demás es opt-in por evidencia: Valkey, OpenSearch, Temporal, Dapr, Kubernetes, Envoy, Pingora, workerd, Firecracker, Coherence y GraalVM no son requisitos para ser “élite”.

## 4. Reglas de composición

1. Un capability tiene un solo source of truth.
2. Un proyecto elige un engine de autorización primario, no Keycloak roles + OpenFGA + Cedar simultáneamente.
3. Los contratos externos viven detrás de anti-corruption layers versionadas.
4. SDK público no concede acceso a la API ni reemplaza sus términos.
5. Samples donan fixtures, tests o ideas; no gobiernan dominio ni seguridad.
6. Un servicio administrado puede ser mejor que autoalojar el repositorio que lo inspira.
7. Toda adopción fija release/commit, license expression, notices, SBOM, upgrade owner y exit plan.

## 5. Próximas profundizaciones

- comparar Cedar vs OpenFGA vs autorización local por dominio;
- contrato de storefront agnóstico frente a Saleor/Shopify/build propio;
- contratos de marketplace/ads por país y versión vigente;
- búsqueda PostgreSQL vs OpenSearch con dataset y SLO reales;
- cache Valkey con invalidación, stampede y degraded mode;
- edge workerd/Pingora sólo con benchmark y threat model.
