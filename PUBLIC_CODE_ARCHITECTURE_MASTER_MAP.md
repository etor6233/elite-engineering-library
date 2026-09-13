# Public Code & Architecture — mapa maestro

> **Estado:** inventario y admisión transversal actualizados a V258 el 2026-09-05; 160 packs/1.395 archivos materializables, 51 perfiles y 121 fuentes oficiales gobernadas. La admisión continúa por claim, no por acumulación.
> **Propósito:** conectar los manuales de élite existentes con código y arquitecturas públicas reutilizables sin degradar sus estándares.
> **Gate:** `PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md`.
> **Regla:** una fila de este mapa no equivale a aprobación. Sólo `ELITE_REFERENCE` y `REUSABLE_PACK` pueden gobernar una implementación.
> **Compañías líderes:** `LEADING_COMPANY_PUBLIC_CODE_MATRIX.md`.
> **Composición completa:** `ENTERPRISE_FULL_STACK_BLUEPRINT.md`.

## 1. Resultado buscado

```text
necesidad de producto
→ manuales autoridad
→ arquitectura mínima suficiente
→ capacidades necesarias
→ candidatos públicos con licencia verificada
→ auditoría contra los manuales
→ reusable packs
→ project contract + tests + evidence
→ implementación y release
```

La biblioteca no será un generador de clones. Debe permitir ensamblar sistemas empresariales rápidamente conservando dominio, seguridad, rendimiento, operación y evolución.

Para contratos OpenAPI Go, la implementación vigente es `implementation_packs/MICROSOFT_KIOTA_OPENAPI_CLIENT_GATE.md`: Microsoft Kiota 1.35.0 exacto, binario autocontenido, generación no destructiva con receipt, determinismo de todo el código, compilación/vet y SCA. Su alcance es transporte tipado; nunca inventa auth, idempotencia, reconciliación ni reglas del proveedor.

Para SAST local de repositorios privados Go/TypeScript sin entitlement GitHub Code Security, el lane seleccionado es `implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md`: reconstruye el commit Microsoft firmado, conserva la adaptación SharpCompress declarada, ejecuta 300 tests y SCA, y sigue siendo sólo lint `ADAPTED / CONDITIONED`. `implementation_packs/GITLAB_OPENGREP_SIGNED_SAST_GATE.md` conserva OpenGrep 1.29.0 y reglas GitLab 2.9.3 como evidencia, pero está fuera del perfil ejecutable porque Cosign 3.1.3 reabrió su admisión por SCA. Ninguno sustituye análisis interprocedural, fuzz, DAST, threat model ni seguridad ofensiva.

## 2. Convención

- **Revisión:** snapshot observado para auditar; nunca significa que branch head sea el pin de producción.
- **Licencia:** expresión observada en esa revisión; cada proyecto vuelve a verificar paths, notices y dependencias.
- **Estado:** se interpreta exclusivamente según el Admission Standard.
- **Uso:** claim estrecho; no describe toda la plataforma propietaria del autor.

## 3. Línea base que ya existe

El corpus anterior ya fijó código industrial para Linux, PostgreSQL, RocksDB, FoundationDB, Kafka, Flink, Kubernetes, Chromium, React, TypeScript, Bazel, Buck2, AndroidX/AOSP, CUTLASS, NCCL, TensorRT, Faiss, DiskANN y publicaciones xAI. Sus claims permanecen en los manuales autoridad; este mapa no los duplica.

La nueva ola prioriza capacidades de ensamblaje de aplicaciones empresariales: identidad, autorización, workflows, commerce/ERP/CRM, observabilidad, edge, mensajería e infraestructura.

### 3.1 Cobertura de búsqueda: sistema completo

La ejecución conversacional usa además `implementation_packs/GO_PG_CONTACT_CHANNEL_IDENTITY.md`. Sus ocho archivos son `AUTHORED`, no código de NIST, PostgreSQL o un proveedor social: NIST SP 800-63 Rev. 4 gobierna únicamente el claim estrecho de sujeto, registros de proofing/consentimiento y bindings; PostgreSQL 18 gobierna concurrencia serializable. El pack exige policy/evidencia explícitas, HMAC tenant/canal, CAS, revocación y conserva cuentas/proveedores live como condiciones.

La salida utiliza `implementation_packs/GO_PG_OUTBOUND_DELIVERY_FENCE.md`, también `AUTHORED`: aplica el patrón de transactional outbox/idempotencia documentado por AWS y concurrencia PostgreSQL sin presentarse como código AWS. No promete exactly-once en un API externo; una llamada incierta se inmoviliza y requiere receipt/status de proveedor para reconciliar.

La búsqueda no se limita a backend. Una solución sólo puede declararse ensamblable cuando cubre o descarta explícitamente todo el flujo:

| Capa | Entregable buscado | Fuentes/patrones ya presentes o en admisión |
|---|---|---|
| producto y dominio | journeys, roles, invariantes, bounded contexts, failure cost | manuales de arquitectura, producto, backend y dominio vertical |
| web pública | catálogo, contenido, SEO, accesibilidad, leads, consentimiento y performance | perfil web/BFF 6/84 materializable, SafeValues/CSP, 38 tests, Playwright/Lighthouse; IdP/backend/edge/a11y AT y target reales condicionados |
| aplicaciones internas | portales admin, franquicia, fábrica, proveedor y operación | Saleor Dashboard como candidato; Oracle JET y Fluent UI como referencias condicionadas |
| clientes | portal de cuenta, pedidos, garantías, servicio, notificaciones y privacidad | frontend authority + identidad/autorización + contratos de API |
| mobile/desktop | lifecycle, offline, cámara/QR, firma, permisos y distribución | manual native; React Native/AOSP auditados; Camera Kit de Tesla en admisión |
| backend y workflows | API/BFF, módulos de dominio, transacciones, outbox, jobs y procesos largos | backend authority, PostgreSQL, Keycloak/OpenFGA, Temporal condicionado por necesidad |
| datos y analítica | source of truth, búsqueda, eventos, lineage, BI, retención y restore | PostgreSQL/streaming/data manuals; OpenTelemetry/Prometheus para telemetría, no ledger |
| integraciones | marketplaces, ads, pagos, CRM, ERP, carriers, fábricas y telemetría vehicular | adapters con ACL/idempotencia/reconciliation; Tesla Fleet en admisión |
| plataforma | IaC, secrets, CI/CD, provenance, rollout, rollback, SLO e incidentes | Kubernetes/OpenTofu/OTel/Prometheus/Envoy/NATS según decision gate |
| assurance | unit, integration, contract, E2E, load, security, resilience y license evidence | Engineering Kit + Admission Standard G0–G8 |

Un repositorio que sólo resuelve una fila no se presenta como “sistema completo”. Un pack full-stack deberá demostrar journeys verticales cruzando UI, API, dominio, datos, autorización, telemetría, deployment y recovery.

## 4. Primera ola — núcleo permisivo

### 4.1 Identidad y autorización

| Fuente fijada | Licencia observada | Capacidad | Estado | Decisión provisional |
|---|---|---|---|---|
| [`keycloak/keycloak@ccd6f59`](https://github.com/keycloak/keycloak/tree/ccd6f59fd8e3a27cebdf1b1f3be42a1d2d4888c4) | [Apache-2.0](https://github.com/keycloak/keycloak/blob/ccd6f59fd8e3a27cebdf1b1f3be42a1d2d4888c4/LICENSE.txt) | OIDC/SAML, SSO, federation, MFA, sessions, identity administration | `CANDIDATE` | referencia fuerte para IdP; no delegar autorización de dominio sólo a roles/tokens |
| [`openfga/openfga@60b3f45`](https://github.com/openfga/openfga/tree/60b3f451d496fb6c49ffbe233de160298cca7c16) | [Apache-2.0](https://github.com/openfga/openfga/blob/60b3f451d496fb6c49ffbe233de160298cca7c16/LICENSE) | ReBAC y autorización fina por organización/objeto/relación | `CANDIDATE` | complementa al IdP; pin de model ID, tuple lifecycle y consistency requieren diseño de aplicación |

**Evidencia oficial actual:** Keycloak 26.7.1 publica guías de producción, clustering, caches, OIDC/SAML, authorization services y separación de endpoints; OpenFGA publica modelado, migraciones y [operación en producción](https://openfga.dev/docs/best-practices/running-in-production). El control de acceso interno experimental de OpenFGA está documentado como no recomendado para producción: no se incorpora como default.

**Manuales autoridad:** `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`, `SOFTWARE_BACKEND_API_ENGINEERING.md`, `AI_SECURITY_GOVERNANCE_PRIVACY.md`, `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md`.

### 4.2 Ejecución durable

| Fuente fijada | Licencia observada | Capacidad | Estado | Decisión provisional |
|---|---|---|---|---|
| [`temporalio/temporal@baa0338`](https://github.com/temporalio/temporal/tree/baa033892586df799bb48ceb03d947ac87c35a13) | [MIT](https://github.com/temporalio/temporal/blob/baa033892586df799bb48ceb03d947ac87c35a13/LICENSE) | workflows durables, timers, retries y procesos largos | `CANDIDATE` | apto para fulfillment, onboarding, garantías y sincronización; no sustituye DB, transacción local ni hot path |

Temporal es un candidato de orquestación, no una autorización para introducir infraestructura distribuida en todo proyecto. Un modular monolith con outbox/jobs puede ser superior cuando el volumen, equipo y failure model no justifican un cluster adicional.

**Manuales autoridad:** `NETWORKING_DISTRIBUTED_STREAMING.md`, `SOFTWARE_BACKEND_API_ENGINEERING.md`, `DATABASE_STORAGE_INTERNALS.md`, `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md`.

### 4.3 Commerce headless

| Fuente fijada | Licencia observada | Capacidad | Estado | Decisión provisional |
|---|---|---|---|---|
| [`saleor/saleor@16e7ea9`](https://github.com/saleor/saleor/tree/16e7ea956a4cf181a5a2a5a8e73a5755c5c9a0b9) | [BSD-3-Clause](https://github.com/saleor/saleor/blob/16e7ea956a4cf181a5a2a5a8e73a5755c5c9a0b9/LICENSE) | catálogo, canales, stock, checkout, promociones, orders, payments y apps vía GraphQL/webhooks | `CANDIDATE` | referencia primaria permisiva para commerce; auditar release estable, dashboard/storefront y contratos de apps por separado |
| [`spree/spree@843bb21`](https://github.com/spree/spree/tree/843bb2153bfef6fb89582e2969428206723add15) | BSD-3-Clause según LICENSE de la revisión | commerce extensible Rails | `LICENSE_VERIFIED` | alternativa a evaluar; no se selecciona hasta comparar modelo, tests, upgrade y extensión contra Saleor |
| [`solidusio/solidus@1f5bf5c`](https://github.com/solidusio/solidus/tree/1f5bf5c638f5fe9fbda4d480bdb7ebcfa39a9a5e) | [BSD-3-Clause](https://github.com/solidusio/solidus/blob/1f5bf5c638f5fe9fbda4d480bdb7ebcfa39a9a5e/LICENSE.md) | commerce Rails derivado de Spree | `LICENSE_VERIFIED` | alternativa condicionada al stack Ruby y a evidencia comparativa |

El README de Saleor declara que `main` puede ser inestable y exige usar releases para producción. Esta advertencia forma parte del gate: el commit de esta tabla es evidencia de inspección, no dependency pin.

**Manuales autoridad:** `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md`, `SOFTWARE_BACKEND_API_ENGINEERING.md`, `DATABASE_STORAGE_INTERNALS.md`, `FRONTEND_PRODUCT_ENGINEERING_UX.md`, `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`.

### 4.4 Observabilidad

| Fuente fijada | Licencia observada | Capacidad | Estado | Decisión provisional |
|---|---|---|---|---|
| [`open-telemetry/opentelemetry-collector@3c16cbc`](https://github.com/open-telemetry/opentelemetry-collector/tree/3c16cbc614bc9d145f514d374b2bf7b77c10e08f) | [Apache-2.0](https://github.com/open-telemetry/opentelemetry-collector/blob/3c16cbc614bc9d145f514d374b2bf7b77c10e08f/LICENSE) | receivers → processors → exporters para traces/metrics/logs | `CANDIDATE` | pipeline portable; componentes alpha/unmaintained no entran por herencia |
| [`prometheus/prometheus@0ab2777`](https://github.com/prometheus/prometheus/tree/0ab2777312302453e037641dc3accdaba6687707) | [Apache-2.0](https://github.com/prometheus/prometheus/blob/0ab2777312302453e037641dc3accdaba6687707/LICENSE) | métricas, scraping, recording/alerting rules y PromQL | `CANDIDATE` | observabilidad numérica; no usar como ledger, billing o auditoría exacta |

La arquitectura oficial de OpenTelemetry advierte que un processor bloqueante compartido desde un receiver puede detener otras pipelines. El pack futuro debe incluir memory limiting, batching, queues, redaction y backpressure. Prometheus declara expresamente que no encaja cuando se necesita exactitud del 100 %, por ejemplo billing por request.

**Manuales autoridad:** `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`, `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md`, `DATA_ENGINEERING_ANALYTICS.md`.

### 4.5 Edge, mensajería y streaming

| Fuente fijada | Licencia observada | Capacidad | Estado | Decisión provisional |
|---|---|---|---|---|
| [`envoyproxy/envoy@d0ea4cd`](https://github.com/envoyproxy/envoy/tree/d0ea4cd7439bf4b6fb8638dee1c7280c9f3a0b70) | [Apache-2.0](https://github.com/envoyproxy/envoy/blob/d0ea4cd7439bf4b6fb8638dee1c7280c9f3a0b70/LICENSE) | proxy L4/L7, TLS, routing, retries, load balancing y observabilidad | `CANDIDATE` | útil cuando el data path y operación justifican proxy dedicado; no default para un proyecto pequeño |
| [`nats-io/nats-server@63b9cc0`](https://github.com/nats-io/nats-server/tree/63b9cc02cf88c7303575d64eb2e410e20499be4b) | [Apache-2.0](https://github.com/nats-io/nats-server/blob/63b9cc02cf88c7303575d64eb2e410e20499be4b/LICENSE) | mensajería request/reply, pub/sub y persistencia JetStream | `CANDIDATE` | candidato de menor peso operacional; garantías se fijan por modo y failure model |
| [`apache/kafka@61e0b5b`](https://github.com/apache/kafka/tree/61e0b5b04f5195c5721a1bc0fe116979a19c7586) | [Apache-2.0](https://github.com/apache/kafka/blob/61e0b5b04f5195c5721a1bc0fe116979a19c7586/LICENSE) | log distribuido, streaming y ecosistema | `ELITE_REFERENCE` para claims ya auditados | no introducir si outbox/queue simple satisface volumen, replay y operación |

Kafka conserva la auditoría profunda existente en `NETWORKING_DISTRIBUTED_STREAMING.md`. Envoy y NATS aún deben superar G3–G7 antes de promoción.

### 4.6 Plataforma e infraestructura

| Fuente fijada | Licencia observada | Capacidad | Estado | Decisión provisional |
|---|---|---|---|---|
| [`kubernetes/kubernetes@e81f39c`](https://github.com/kubernetes/kubernetes/tree/e81f39c0e03ce8ed8e2660c9147b391edd9e262b) | [Apache-2.0](https://github.com/kubernetes/kubernetes/blob/e81f39c0e03ce8ed8e2660c9147b391edd9e262b/LICENSE) | control loops, desired state y orquestación | `ELITE_REFERENCE` para claims ya auditados | no es requisito de una arquitectura de élite; adoptarlo sólo con necesidad operacional |
| [`opentofu/opentofu@5097d2d`](https://github.com/opentofu/opentofu/tree/5097d2de3a2c36c66f43b7f0e2f07cd82882a32e) | [MPL-2.0](https://github.com/opentofu/opentofu/blob/5097d2de3a2c36c66f43b7f0e2f07cd82882a32e/LICENSE) | infraestructura declarativa, plan/state/providers | `CANDIDATE` | MPL exige control por archivo; state contiene secretos y necesita backend/locking/recovery |
| [`backstage/backstage@1a705ca`](https://github.com/backstage/backstage/tree/1a705cad96f0fbd48d4f7fe7e92e8a45d6afbab7) | [Apache-2.0](https://github.com/backstage/backstage/blob/1a705cad96f0fbd48d4f7fe7e92e8a45d6afbab7/LICENSE) | catálogo de software, templates y portal de plataforma | `CANDIDATE` | útil para múltiples equipos/servicios; sobrecosto para una sola aplicación |

### 4.7 Empresas de Elon: xAI/X, Tesla y el límite de SpaceX

xAI/X ya tenía una auditoría material en `AI_ENGINEERING_MASTER_MAP.md`: `x-algorithm` aporta retrieval/ranking end-to-end, `grok-1` una referencia MoE, `grok-build` un runtime de agente con tools/sessions/permissions, y `xai-proto`/`xai-sdk-python` contratos y cliente. Sus licencias no son uniformes: los componentes Apache-2.0 pueden avanzar por el gate; `grok-prompts` es AGPL-3.0 y `xai-cookbook` tiene términos Beta restrictivos.

| Fuente fijada | Licencia observada | Capacidad | Estado | Decisión provisional |
|---|---|---|---|---|
| [`xai-org/x-algorithm@28e…`](https://github.com/xai-org/x-algorithm) | Apache-2.0 | retrieval, hydration, filters, multi-action ranking, selection e inference ejecutable | `ELITE_REFERENCE` para claims auditados en el mapa IA | reusable sólo después de separar datos/modelos/dependencias y reproducir el slice aplicable |
| [`teslamotors/fleet-telemetry@4767414`](https://github.com/teslamotors/fleet-telemetry/tree/4767414c1297c8db37a082561c2b761ce8613590) | [Apache-2.0](https://github.com/teslamotors/fleet-telemetry/blob/4767414c1297c8db37a082561c2b761ce8613590/LICENSE) | ingreso seguro de telemetría vehicular, WebSocket, mTLS, rate limiting y dispatch a Kafka/Kinesis/PubSub/MQTT/Redis/ZMQ | `CANDIDATE` | extraordinariamente relevante para flotas eléctricas; exige threat model de VIN, minimización de datos, HSM y tests de mensajes falsificados |
| [`teslamotors/vehicle-command@f97fa1e`](https://github.com/teslamotors/vehicle-command/tree/f97fa1e4bf617a364c72b85cb5d859528abeda67) | [Apache-2.0](https://github.com/teslamotors/vehicle-command/blob/f97fa1e4bf617a364c72b85cb5d859528abeda67/LICENSE) | SDK Go y proxy REST para comandos con OAuth, claves enroladas y autenticación extremo a extremo | `CANDIDATE` | adapter específico Tesla; API 0.x sin estabilidad garantizada, keys/roles/consentimiento y safety requieren gate propio |
| [`teslamotors/react-native-camera-kit@a2a81ce`](https://github.com/teslamotors/react-native-camera-kit/tree/a2a81cee8995c831aca472836e2be446994ba070) | [MIT](https://github.com/teslamotors/react-native-camera-kit/blob/a2a81cee8995c831aca472836e2be446994ba070/LICENSE) | cámara, QR/barcodes y detección en React Native | `LICENSE_VERIFIED` | componente útil para inventario/servicio, no arquitectura mobile completa; permisos, privacidad y lifecycle pendientes |
| [`teslamotors/fixed-containers@d0df5f6`](https://github.com/teslamotors/fixed-containers/tree/d0df5f67e102ba42025ef0ecf761d19a4510a543) | [MIT](https://github.com/teslamotors/fixed-containers/blob/d0df5f67e102ba42025ef0ecf761d19a4510a543/LICENSE) | contenedores C++20 de capacidad fija, `constexpr` y sin allocations dinámicas | `LICENSE_VERIFIED` | candidato sólo para paths C++ con necesidad demostrada de layout/latencia; irrelevante para la mayoría del SaaS empresarial |

Las publicaciones Tesla de Buildroot, Linux y coreboot son fuentes oficiales valiosas para firmware/kernel y cumplimiento, pero no revelan la arquitectura de producto completa del vehículo ni son un starter de aplicaciones empresariales. Se usan como `CONDITIONED` reference, respetando copyleft y separación por componente.

**SpaceX/Starlink:** no se verificó un repositorio oficial público que exponga la arquitectura de vuelo, ground systems o Starlink como código reutilizable. Starlink declara software propietario en object code y ofrece, por obligación de las licencias correspondientes, acceso a ciertos componentes open source mediante la cuenta del usuario. Eso es evidencia de compliance, no una licencia sobre la arquitectura propietaria. El usuario GitHub `SpaceX` pertenece a una persona y `r-spacex/SpaceX-API` es comunitario, no oficial y fue archivado; ambos quedan excluidos como autoridad SpaceX.

### 4.8 Oracle: frontend, backend, runtime, datos e infraestructura

| Fuente fijada | Licencia observada | Capacidad | Estado | Decisión provisional |
|---|---|---|---|---|
| [`oracle/oraclejet@44b10ad`](https://github.com/oracle/oraclejet/tree/44b10adf16d04049cbe4f5d800c16f56e95a6417) | [UPL-1.0](https://github.com/oracle/oraclejet/blob/44b10adf16d04049cbe4f5d800c16f56e95a6417/LICENSE.txt) | toolkit JavaScript modular para aplicaciones cliente empresariales | `CANDIDATE` | lane frontend/admin; comparar accesibilidad, design system, bundle, SSR y talento contra React stack antes de adoptar |
| [`helidon-io/helidon@9d258d2`](https://github.com/helidon-io/helidon/tree/9d258d29652aada6df0e003ff048f474c4220c05) | [Apache-2.0](https://github.com/helidon-io/helidon/blob/9d258d29652aada6df0e003ff048f474c4220c05/LICENSE.txt) | librerías Java para servicios, webserver, client, health, metrics, tracing, security y virtual threads | `CANDIDATE` | backend Java serio; no impone microservicios y debe compararse con el stack/equipo real |
| [`oracle/graal@3f1e5b5`](https://github.com/oracle/graal/tree/3f1e5b5f6ce52798791d2f8b494516c98fe83ffe) | componentes GPL-2.0, GPL-2.0 con Classpath Exception, UPL y BSD-3-Clause | Native Image, compiler, SubstrateVM, Truffle y runtimes | `CONDITIONED` | admisión path/component; medir build time, memoria, warmup, reflexión y observabilidad antes de elegir native image |
| [`oracle/coherence@bd13f35`](https://github.com/oracle/coherence/tree/bd13f3536e7f009fc6a487ca497d1c1ab90f8e0d) | [UPL-1.0](https://github.com/oracle/coherence/blob/bd13f3536e7f009fc6a487ca497d1c1ab90f8e0d/LICENSE.txt) | data grid distribuido, particionado, compute-on-data y clientes polyglot | `CANDIDATE` | sólo cuando escala/latencia/resiliencia justifican grid; PostgreSQL/cache simple sigue siendo el default empresarial |
| [`oracle/terraform-provider-oci@8de3c48`](https://github.com/oracle/terraform-provider-oci/tree/8de3c480b4bf08d8edd40ef0cc4ad0340bd01721) | [MPL-2.0](https://github.com/oracle/terraform-provider-oci/blob/8de3c480b4bf08d8edd40ef0cc4ad0340bd01721/LICENSE.txt) | recursos declarativos y discovery para Oracle Cloud Infrastructure | `LICENSE_VERIFIED` | adapter IaC si el proyecto selecciona OCI; controlar obligaciones por archivo, state, credentials y drift |

Oracle aporta código industrial real, pero no se convierte en stack obligatorio por marca. MySQL, OpenJDK, VirtualBox y demás familias se auditarán sólo si una decisión del proyecto los vuelve materiales; agregar productos sin una necesidad concreta aumenta superficie legal y operacional.

## 5. Licencias mixtas o copyleft — referencias condicionadas

| Fuente fijada | Términos observados | Valor técnico | Estado | Restricción de incorporación |
|---|---|---|---|---|
| [`medusajs/medusa@0451e06`](https://github.com/medusajs/medusa/tree/0451e065e7ae393e3b6a963894b04c43316eb366) | raíz: MIT excepto paths Enterprise regidos por `ENTERPRISE-LICENSE.md` | módulos aislados, workflows, commerce e integraciones | `CONDITIONED` | auditar cada path/package; nunca declarar todo el monorepo MIT |
| [`strapi/strapi@e5f3785`](https://github.com/strapi/strapi/tree/e5f3785bc92b28bb7819af9a1b2d6e678a1383a2) | Community MIT; paths `ee/` bajo licencia empresarial | CMS headless y extensibilidad | `CONDITIONED` | separar Community/Enterprise y revisar términos de cuenta/cloud |
| [`frappe/erpnext@5fa68dd`](https://github.com/frappe/erpnext/tree/5fa68dd06835f77a60ccefa033d7955d83294d3b) | GPL-3.0 | ERP: ventas, compras, stock, contabilidad, CRM y operación | `CONDITIONED` | gran valor de dominio; copyleft y adaptación requieren decisión explícita |
| [`frappe/frappe@4230fea`](https://github.com/frappe/frappe/tree/4230fea04f6b18db13dcd6ddcb74129b1b3ba785) | MIT | framework de metadata, permisos, workflows y apps | `LICENSE_VERIFIED` | licencia del framework no cambia GPL de ERPNext |
| [`odoo/odoo@b999d8d`](https://github.com/odoo/odoo/tree/b999d8d93ac4edb4710cf1e49c80e5ce8633df9f) | [LGPL-3.0](https://github.com/odoo/odoo/blob/b999d8d93ac4edb4710cf1e49c80e5ce8633df9f/LICENSE), con librerías/contribuciones compatibles adicionales | ERP/CRM/extensión modular | `CONDITIONED` | distinguir Community, Enterprise, addons y términos por módulo |
| [`vendure-ecommerce/vendure@cb52389`](https://github.com/vendure-ecommerce/vendure/tree/cb523893d534a19641b86f9302b130386f25d485) | GPL-3.0 o licencia comercial; excepción para plugins | commerce TypeScript/NestJS/GraphQL plugin-first | `CONDITIONED` | plugins pueden tener otra licencia bajo excepción; core/derivados no se tratan como permisivos |
| [`twentyhq/twenty@e6ecef8`](https://github.com/twentyhq/twenty/tree/e6ecef852decaf687eb84e6f1f17d6c25cacec9a) | mayoría AGPL-3.0, archivos Enterprise y paquetes MIT identificados | CRM y application platform | `CONDITIONED` | path-level license obligatorio; no copiar core a producto cerrado sin decisión |

Copyleft no significa baja calidad. Significa que la forma de composición, modificación, distribución o acceso por red puede imponer obligaciones. Estos proyectos pueden ser excelentes para adoptar completos, estudiar dominio o integrar por API; no entran automáticamente como bloques permisivos.

## 6. Rechazo legal preventivo

| Fuente fijada | Términos observados | Estado | Razón |
|---|---|---|---|
| [`directus/directus@2abb57f`](https://github.com/directus/directus/tree/2abb57f6dd79911ec7c331ad63267830f542fc46) | Monospace Sustainable Core License 1.0, restricción de `Competing Use`, conversión futura a GPL-3.0 | `REJECTED` como baseline reusable | código visible pero licencia actual no aprobada como open source permisivo; sólo evaluar para uso permitido/comercial específico |

La decisión no evalúa calidad técnica de Directus. Evita convertir una licencia source-available con restricción de campo de uso en “código público libre para incorporar”.

## 7. Reference applications — evidencia, no producción

| Fuente fijada | Licencia | Qué enseña | Estado | Límite explícito |
|---|---|---|---|---|
| [`dotnet/eShop@ae71a06`](https://github.com/dotnet/eShop/tree/ae71a0610162c6712c885f64e7706a9b7b54f2ee) | MIT | .NET 10, Aspire, servicios, tests y journeys de commerce | `CONDITIONED` | reference app; su despliegue documentado usa datos descartables y no es producción |
| [`GoogleCloudPlatform/microservices-demo@34ffea9`](https://github.com/GoogleCloudPlatform/microservices-demo/tree/34ffea9175946982c3088ed84994fe6019ad6e92) | Apache-2.0 | gRPC polyglot, Kubernetes, instrumentation y load generator | `CONDITIONED` | demo de 11 microservicios; no prueba que esa división sea adecuada |
| [`aws-containers/retail-store-sample-app@1a28474`](https://github.com/aws-containers/retail-store-sample-app/tree/1a28474f2461459f42e6b393db59e7d1434d4aec) | MIT-0 | containers, Helm/Terraform, telemetry y load | `REJECTED` como arquitectura de producto | el propio proyecto dice educativo, no productivo y deliberadamente sobre-ingenierizado |

Los samples pueden donar test harness, manifests o patrones acotados después de G0–G8. Nunca gobiernan bounded contexts, seguridad, consistencia o topología del proyecto.

## 8. Paquetes prioritarios

Orden por capacidad de desbloquear proyectos empresariales sin sobrearquitectura:

1. `FRANCHISE_FULL_STACK_FOUNDATION` — `CANDIDATE_PACK` creado; journeys verticales, web pública, portales, API, dominio, datos, IAM, tests y release.
2. `IDENTITY_AUTHORIZATION_KEYCLOAK_OPENFGA` — `CANDIDATE_PACK` creado; admisión profunda pendiente.
3. `POSTGRES_TRANSACTIONAL_MODULAR_MONOLITH` — `CANDIDATE_PACK` creado; source of truth, outbox, migrations y RLS condicionada.
4. `PUBLIC_WEB_ADMIN_PORTALS` — `CANDIDATE_PACK` creado; web pública, cliente, admin, accesibilidad, seguridad, performance y Playwright.
5. `MULTI_BRANCH_FRANCHISE_DOMAIN` — organizaciones, sucursales, territorios, empleados, inventario y permisos.
6. `COMMERCE_CATALOG_ORDERS_INVENTORY` — comparar Saleor/Medusa/implementación propia.
7. `DURABLE_SUPPLIER_FACTORY_FULFILLMENT` — outbox/jobs primero; Temporal sólo si supera el decision gate.
8. `ELECTRIC_FLEET_TELEMETRY_COMMANDS` — adapters vehiculares, consentimiento, keys, HSM, ingest y safety; Tesla sólo si aplica.
9. `INTEGRATION_PAYMENTS_MARKETPLACES_ADS` — Google Lead Form está materializado en `GO_OMNICHANNEL_LEAD_INGRESS`; Mercado Libre Questions agrega topic allowlisted, fetch/reconciliación v4, normalización al mismo owner y respuesta aprobada `POST /answers` con confirmación/read-after-write más reconciliación GET-only sobre el fence 0.2.x; Meta agrega señal firmada y durable antes del ACK más recuperación/import; TikTok Lead v1.3 agrega retrieval autenticado, evidencia hash-linked e import al mismo PostgreSQL/outbox. Promoción CRM/consentimiento sigue separada. Mercado Libre es `AUTHORED` contra su contrato HTTP porque su SDK oficial está archivado; el wrapper TikTok es `ADAPTED` sobre transporte genérico del wheel oficial. Cuentas/suscripciones live, autenticidad inbound, writes adicionales/postbacks y producción siguen condicionados.
10. `OBSERVABILITY_SECURE_DELIVERY` — `CANDIDATE_PACK` creado; OTel/Prometheus, redaction, SLO, provenance, canary, restore y rollback.

## 9. Arquitectura inicial para una plataforma de franquicias

Esto es un punto de partida a validar, no una plantilla universal:

```text
public web / admin / customer portal
              ↓
        BFF or versioned API
              ↓
modular application core + PostgreSQL source of truth
  ├─ identity adapter → Keycloak or managed OIDC
  ├─ authorization adapter → local policy/OpenFGA by demonstrated need
  ├─ catalog/products/vehicles/batteries
  ├─ organizations/franchises/branches/territories
  ├─ leads/customers/CRM
  ├─ suppliers/factories/purchases/logistics
  ├─ inventory/reservations/orders/payments
  ├─ warranties/service/maintenance
  ├─ documents/audit/notifications
  └─ integration outbox → marketplaces/ads/ERP/carriers

telemetry → OTel pipeline → metrics/traces/log backend
release → immutable artifact → staged rollout → rollback
```

Default: monolito modular desplegable como una unidad, fronteras internas estrictas y outbox. Separar un servicio sólo cuando ownership, SLO, escala, aislamiento o release independiente lo demuestre.

## 10. Próxima admisión

- [x] primer `REUSABLE_PACK`: readiness validator A–H, 48 superficies y journeys release-bound;
- [x] composición full-stack de referencia: franquicia 67 packs/743 archivos y backend/web PostgreSQL materializables; no implica cierre de todos los journeys;
- [x] modular monolith + PostgreSQL + outbox/inbox/jobs, observabilidad y secure delivery;
- [x] payments/marketplaces/ads con separación SDK/API/terms y dieciséis adapters auditados;
- [ ] recorrido navegador→backend→PostgreSQL→efecto→respuesta/recovery con ayuda/capacitación/soporte de la misma release;
- [ ] mutaciones Mercado Libre listing/stock/precio/postventa y receipts/reconciliación de notificaciones, sólo contra contratos oficiales actuales y sin reutilizar sus SDKs archivados;
- [ ] telemetría target-agnostic del journey sin PII y generación compacta de authority map/manifest;
- [ ] promoción por claim de foundations que cierren G0–G8 sin depender de evidencia imposible de portar;
- [ ] packs Tesla Fleet Telemetry/Vehicle Command sólo si un blueprint activa telemetría/comandos vehiculares;
- [ ] reauditoría continua de links, revisiones, SPDX/path overrides, releases y advisories.

## 11. Regla final

El objetivo no es acumular más software. Es reducir tiempo desde una necesidad hasta una implementación correcta sin importar una licencia incompatible, una topología de demo o un supuesto de escala ajeno.

### Node advisory engine y adapter acotado — V320

nodejs/is-my-node-vulnerable1.6.1/c37a56bad56e34fe5223ddd3cb223cc4158136ae MIT:
motor oficial conservado, getJsonADAPTED y7 archivosAUTHORED claramente separados.
NODE_OFFICIAL_RUNTIME_ADVISORY_GATE0.1.0 (1/10) reconstruido y probado; condiciones
de instancia demostradas con USE_REUSABLE_PACK, no aprobación del CLI raw que
falló corpus vacío/deadline enV319. Ver expediente V320 y perfil focal; no elevar
el claim a npm/pnpm, malware, otros runtimes o monitoring de producción.

### V334 — proyección pnpm canónica completada

PNPM-ARTIFACT-SELECTION-GATE0.1.0,4files/17tests dos veces y2proyecciones reales443files idénticas; G0–G8/USE_REUSABLE_PACK sólo tooling en este contexto. Perfil separado1/4 y2pasos de verificación integrados. No consumers/runtime/redistribución admitidos; condiciones explícitas y FAIL532 original conservado. Evidencia reconstruction_evidence/PNPM_SELECTION_GATE_V334.md.

### V335 — routing de consumidores pnpm fijados

PNPM-ARTIFACT-SELECTION-GATE0.2.0,6files/43tests fuente+rebuild. Prepare/verify produce recetas offline con configuración aislada, store/cache explícitos y perfiles exactos; nunca ejecuta ni admite runtime.3instalaciones desde rebuild PASS sin descargas/cambios de lock;115tests/1skip ybuild web PASS. FAIL555/557 corregidos sin omitir política; G0–G8 acotado. Perfil separado1/6; producto67/746 inalterado. Evidence PNPM_CONSUMER_ROUTING_V335.md.

### V337 — licencia exacta npm-lifecycle y transporte gobernado

OFFICIAL-UPSTREAM-ACQUISITION-CORE0.4.79:33files/125sources,78opaque checks/rebuild y18perfiles de adquisición. Nuevo perfil separado adquiere sólo npm-lifecycle1100.1.0; firma registry e integridad verificadas. LICENSE del tar exacto coincide con commit/sidecar,7de8archivos con Git blobs; único delta packageManager omitido.29planes de composición actualizados. Fuente licencia verificada, runtime/redistribución no admitidos. Ver reconstruction_evidence/PNPM_LIFECYCLE_LICENSE_V337.md.

V367: investigación oficial de entrenamiento prioriza TRL/SFT y PEFT opcional como DISCOVERED; torchtune/torchforge fuera de primera línea por mantenimiento, torchtitan diferido por scope/runtime. Identidades/artefactos sólo observados en metadata, no adquiridos/admitidos. RESEARCH_INCOMPLETE/FAIL663; TEST09 sigue bloqueado. Evidencia reconstruction_evidence/HISTORY_TRAINING_SOURCE_RESEARCH_V367.md y training_gap_v367/record.json. No pedir corpus privado para esta preparación de biblioteca.

V367 distribución181: los registros training_gap_v367 se conservan íntegros, con SHA-256, como secciones del reporte HISTORY_TRAINING_SOURCE_RESEARCH_V367.md; materializarlos sólo en raíz aislada de investigación. No son nuevos archivos sueltos del release ni un pack admitido.

V369: core de adquisición0.4.80,36files,127sources y19perfiles internos. Perfil aislado adquirió2sdists TRL1.12.0/PEFT0.20.0 y conserva7outputs inmutables.559archivos inspeccionados:544Git-equal,15metadata de packaging revisados,licencias raíz iguales.98checks transporte; no runtime/training admission. FAIL663 sigue investigación pendiente del grafo y pipeline. Evidencia reconstruction_evidence/TRAINING_SOURCE_ACQUISITION_V369.md.

V370: candidato59distribuciones/103relaciones,58baseTRL+PEFTopcional.59METADATA hash-verified/equivalentes; resolver pip confirma59y2negativos sin wheels/install/framework execution. OSV59versiones0hallazgos declarados;30provenance subjects coinciden,firmas no verificadas. Licencias de artefactos/nativos,adquisición gobernada y runtime/gates pendientes. DISCOVERED/FAIL663 RESEARCH_INCOMPLETE. Ver reconstruction_evidence/TRAINING_DEPENDENCY_GRAPH_V370.md.

V371: core0.4.81/39files;55quarantine checks+98opaque PASS.59wheels/214055762bytes adquiridos;23854files/23795RECORD hashes inspeccionados.226native files/3SBOMs;371identidades consultadas,6records=4avisos distintos (3security+1maintenance). FAIL675 mantiene cuarentena. Safetensors sourceb7c0f38 corrige pyo3/memmap2,45registry crates0OSV sólo metadata; build pendiente.43/48sin cambio. Ver reconstruction_evidence/TRAINING_WHEEL_QUARANTINE_V371.md.

V372 maintenance: reference telemetry implementation is returned to canonical worker0.1.5 and WINDOWS_REFERENCE_TELEMETRY_RUNTIME.md/0.1.0, with opt-in WINDOWS_REFERENCE_TELEMETRY_PACK_PLAN.md. Source-only acquirer0.4.82 remains separate from runtime admission. TEST05 closure requires the fresh integrated replay and gates in reconstruction_evidence/INTEGRATED_TELEMETRY_CONTROL_V372.md; no production acceptance or private training data inferred.

V373: opt-in HISTORY-MODEL-TRAINING-PIPELINE0.1.0/17files, HISTORY_MODEL_TRAINING_PACK_PLAN2packs/21files;26policy tests and actual earlier synthetic SFT/evaluation/rollback. Final exact reconstruction/bootstrap/E2E and TEST09 decision are recorded in reconstruction_evidence/HISTORY_MODEL_TRAINING_CONTROL_V373.md. NEW requires no history; EXISTING keeps data/model private with its own authority/quality/runtime gates. No automatic training, provider call or deployment. SDK component selection precedes AUTHORED glue; see HISTORY_TRAINING_SDK_QUALIFICATION_V373.md and HISTORY_TRAINING_SDK_GAP_V373.md. Integral franchise remains67/754.

V374: auditoría por claim21núcleos,228tests/306semillas/23campañas (39811771ejecuciones),6comparaciones nuevas de fuentes,42archivos reconstruidos,20packs/40archivos compuestos y candidato rechazado.20CONDITIONED/1CANDIDATE intactos; cero equivalencias empresariales/promociones inventadas.13negativos de enlace de evidencia integrados en VERIFY_LIBRARY. Preflight/cierre del control aún pendientes;45/48sin cambio. Ver reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md; la integración de negocio más amplia de FAIL385 mantiene sus owners.

V374 cierre203: TEST02 PASS como auditoría por claim21núcleos;228tests/306semillas,23targets finales (6389808ejecuciones),20packs/40files compuestos, candidato rechazado,13negativos de evidencia y Preflight202160pasos PASS.46/48controles;20CONDITIONED/1CANDIDATE y FAIL385/integración más amplia conservados.164/1506/819/55, integral67/754. No cambia48oráculos ni promueve producto. Ver CORE_CLAIM_ADMISSION_V374.md; continuar TEST03 y luego TEST07 según gates materiales.

V374 rectificación204 / FAIL710: el cierre203 de TEST02 se retracta. Fuente/tests/composición son válidos, pero FAIL385 aún exige equivalencia/integración material.45/48PASS; TEST02/03/07BLOCKED. Se preservan historia y mejoras; no se cambia el oracle ni se reduce el roadmap. Ver CORE_CLAIM_ADMISSION_V374.md, encabezado vigente.

V394: PNPM-ARTIFACT-SELECTION-GATE0.3.0/6files corrige FAIL782 (pins empresariales y Playwright obsoletos), entrega BLUEOAK-NOTICE.md ligado a5declaraciones exactas y verifica3recetas actuales;51tests y10mutaciones reales rechazadas. Payload442/22notices intactos; sin ejecución ni admisión runtime/redistribución. Evidencia: reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md. Los pendientes BlueOak se reducen a entrega downstream/revisión restante; ausencia de LICENSE original chownr no se falsea. Otros permisos/source/publicación de pnpm continúan abiertos.

V395: PNPM-ARTIFACT-SELECTION-GATE0.4.0/6files entrega QRCODE-NOTICE.md con encabezado de autor/modificación exacto y texto MIT completo, además del aviso BlueOak intacto.59tests PASS,3recetas reales y12negativos rechazados;10source blobs/10regiones de bundle fijados por separado, sin claim de equivalencia completa. FAIL783 entrega local resuelto, FAIL784 metadata ADAPTED reconstruida;payload442/22notices intactos. No ejecución ni admisión runtime/redistribución;45/48. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V396: cerrado descubrimiento y entrega local de licencia original semver-utils1.1.4: MIT OR Apache-2.0 explícito, opción MIT completa y1839bytes originales retenidos. Core0.4.88/51files/196sources/25profiles añade sólo cuarentena npm exact4193bytes/SHA512;140checks y gates anteriores PASS. Planner0.5.0/6files/67tests;3recetas reales/12negativos PASS, BlueOak/QRCode y442payloadfiles intactos. Clave registry vencida documentada; no firma vigente, equivalencia build, ejecución de recetas ni admisión global pnpm.45/48sin promoción; V386/Daybreak sigue diferido. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V397: PNPM-ARTIFACT-SELECTION-GATE0.6.0/7files entrega el conjunto completo de47textos retenidos/151253bytes en3recetas:141copias verificadas independientemente,12negativos reales y82tests PASS.22notices originales comparados contra payload y25evidencias conservan alcance; BlueOak/QRCode/semver y442payloadfiles intactos. Catálogo ADAPTED con términos por texto, sin algoritmos ni promoción de licencias/source/relinking/runtime.45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V398: PNPM-ARTIFACT-SELECTION-GATE0.7.0/8files entrega fuente original next-path1.0.0, manifiesto y MPL completa en3recetas:9copias exactas/12negativos reales/90tests PASS. Commit oficial y3Git blobs verificados;4sentencias comparadas bajo adaptadores explícitos,10mutaciones rechazadas. No equivalencia runtime ni reproducibilidad pnpm. Colección47, suplementos anteriores y442payloadfiles intactos.45/48; Daybreak diferido. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

V401: contención ZIP terminada en candidato pnpm aislado;16módulos adm-zip retirados,441archivos preservados/1bundle cambiado,4rechazos sinIO y3instalaciones offline reales PASS.475identidades conocidas/0hallazgos; SBOMrecursivo no cerrado;12571archivos node_modules idénticos. Pack0.8.0/9files,167/1588; pnpm general BLOCKED, TEST02/03/07 siguen abiertos,45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md.

V401 reconciliación final: comparación12571/12571idéntica completada antes de la solicitud de detenerla; no se detuvo proceso. FAIL807 separa las3instalaciones iniciales offline de un exec Next que descargó71paquetes por configuración omitida. Exec corregido con store/offline explícitos PASS sin descargas. Historial/log anterior retenido; no prueba global de ausencia de red. reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md

V402303 / GO_CONNECTED_MARKETPLACE_MUTATION0.1.0: AUTHORED glue on admitted owners and seven SHA-locked official HTTP contracts; G0–G8 narrow local evidence in reconstruction_evidence/MARKETPLACE_MUTATION_RELEASE_V402.md/json. No corporate authorship or whole-T2805 promotion.

V402304 / GO_CONNECTED_MARKETPLACE_MUTATION0.2.0: narrow initial source PNG→manual upload→manual CREATE→GET recovery PROVEN_LOCAL. AUTHORED glue,8fixed official contracts, G0–G8 in reconstruction_evidence/MARKETPLACE_INITIAL_RELEASE_V402.md/json. Existing content/feed/comms closure remains T2805.

V402305 / GO_CONNECTED_MARKETPLACE_MUTATION0.3.0: approved catalog/media→single-unsold-item CONTENT→GET recovery PROVEN_LOCAL. AUTHORED glue with original source owners, same eight official contracts. G0–G8 in MARKETPLACE_CONTENT_RELEASE_V402.md/json; T2805 feeds/comms remain.

V402306: GO_CONNECTED_GOOGLE_MERCHANT0.1.0: AUTHORED current-source/ATP/approval/fence/SDK-IPC/PG/queue binding. Unchanged official Google SDK dependency pin. G0–G8 local claim only; remaining T2805 explicit. MERCHANT_CONNECTED_RELEASE_V402.md/json.

V402307: GO_CONNECTED_SCHEDULED_WHATSAPP0.1.0 binds original CRM/consent/manual approval/jobs/fence/Meta adapter; G0–G8 narrow local claim. Campaign and later T280x remain open. SCHEDULED_COMMUNICATIONS_RELEASE_V402.md/json.

V402308: GO_CONNECTED_WHATSAPP_CAMPAIGNS0.1.0 uses selected source/approval/job/provider owners; G0-G8 local. Consolidates T2805; global source/security staysT2803. CAMPAIGN_CONNECTED_RELEASE_V402.md/json.

V402309: GO_OIDC_PORTAL_SESSION0.1.0 AUTHORED orchestration around admitted SDKs; G0-G8 narrow local session lifecycle. J5 administration and source/SCA remain T2803. IDENTITY_PORTAL_RELEASE_V402.md/json.

V402310: J5 explicit bootstrap and provider-side role-review contract proven locally; full source/SCA remains T2803. IDENTITY_J5_RELEASE_V402.md/json.

V402311: PNPM-ARTIFACT-SELECTION-GATE0.9.0 adds proven narrow offline local execution; reused official Node flags and existing selector/planner. AUTHORED glue only; no reputation-based admission. PNPM_LOCAL_RUNTIME_V402.md/json.

V402312: T2803 PROVEN_LOCAL identidad/source/SCA/lint sobre composición109/1481. Lint sólo evidencia estrecha; scope ARCA penúltimo y Daybreak diferido permanecen separados. COMPOSITION_SECURITY_RELEASE_V402.md/json.

V402313: T2806 referencia conectada PROVEN_LOCAL; AUTHORED sólo glue/fixtures, SDK DEPENDENCY_PIN y notices VERBATIM.20clases fuera de referencia explícitas. DOCUMENT_REFERENCE_RELEASE_V402.md/json.
