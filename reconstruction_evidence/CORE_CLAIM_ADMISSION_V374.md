# V374 — auditoría por claim de los 21 núcleos

## Current checkpoint207 — operations correction verified

**45/48 controls passed; TEST02/03/07 remain blocked.** FullPreflight206 completed162executed steps PASS and composed all55plans. This includes the exact source-engine rule regression. Inventory164packs/1507files/819Markdown/55plans; integral67/755. Only Docker availability remains missing in that receipt. The source audit below and real rule correction are supporting progress, not closure of the unresolved functional integration control. Evidence for207 is appended after the preserved205candidate history.

## Successor205 — actual low-traffic alert defect corrected

Current acceptance remains45/48; TEST02/03/07 blocked. SECURE-OPS-DELIVERY-CORE1.1.4 fixes a real generic defect under the existing operations owner: dividing by a1request/second denominator floor hid100% failures at low traffic. The rule now divides by actual traffic while retaining the1% threshold,10minute hold, labels and annotations. It remains an error-ratio alert, not multi-window SRE burn-rate or a connected HTTP exporter.

Canonical source-engine verification:9scenarios/13assertions PASS with exact admitted V372 Prometheus3.14.0 promtool. Sparse failure and counter reset were red against1.1.3. Healthy/no traffic, exact/above threshold, pending delay, a real20% short spike and recovery after firing are covered. The first recovery fixture had a mutable-oracle error; its rejected bytes/log remain in the evidence. Canonical LF reconstruction was executed independently after candidate CRLF transport was diagnosed.

The reusable fixture is included in all three existing consumers. VERIFY_EXECUTABLE_LIBRARY accepts explicit PromtoolExecutable+PromtoolSha256, verifies identity and runs check/test rules through the existing60second bounded runner. Without both inputs it reports NOT_EXECUTED;5invalid-input cases reject before execution/materialization. Tool hash selection is not new SCA/provenance or production admission. The V372 exact engine retains its prior narrow security binding; no new dependencies were acquired.

Inventory becomes164packs/1507files/819Markdown/55plans, provenance1260AUTHORED/140ADAPTED/107VERBATIM, integral67/755. The earlier21core comparison snapshot and its SECUREOPS1.1.3 hash are preserved as historical inputs, not overwritten as if they were the current1.1.4 owner. The21cores themselves are unchanged by this fix. Full successor Preflight is pending below; no control is promoted.

2026-09-10. Mantenimiento de biblioteca. Checkpoint201 validado antes de cambios.



## Estado vigente corregido — checkpoint204, TEST-02 BLOCKED

45/48controles PASS. Continúan TEST-02, TEST-03 y TEST-07.
La auditoría21,228tests,306semillas,23campañas finales, composición20 y13negativos
son resultados reales preservados. No satisfacen por sí solos la equivalencia/
integración pendiente de FAIL385. La promoción203 fue prematura y se revierte por
un sucesor, sin borrar eventos ni modificar el oracle; ver FAIL710.
No se declaran completos los controles mediante un dictamen parcial.

## Cierre203 retractado — auditoría insuficiente para cerrar TEST-02

46/48controles PASS; quedan TEST-03 y TEST-07. El oracle de TEST-02 se conserva
y se cumple como auditoría de los claims actuales, no como21admisiones empresariales.
20CONDITIONED y1CANDIDATE conservados; cero promociones por reputación.
FAIL385 y los macrofrentes de integración mantienen sus propios pendientes;
no se declaran completos ni se ocultan como configuración ya aprobada.

Preflight202:160pasos ejecutados PASS,55planes compuestos;164packs/1506archivos/
819Markdown. Disponibilidad BLOCKED sólo Docker, aparte de esta auditoría.
228tests ordinarios,306semillas,20packs/40archivos compuestos y13negativos
permanentes PASS. La campaña inicial23x10s registra39811771ejecuciones.
La comprobación adicional sobre los42archivos finales exactos,23x3s, registra
6389808ejecuciones PASS. Ambas conservan su identidad y presupuesto;
no se suman como porcentaje de calidad ni sustituyen seguridad integral.

## Integración previa preservada

Las21decisiones por claim están contrastadas con fuente, revisión, licencia, límites y pruebas.
TEST-02 continúa BLOCKED hasta completar el Preflight y el cierre del checkpoint.
Los48oráculos/acciones se conservan exactamente. No se reclasifica una capability como OPTIONAL/NONE.

## Qué demuestra esta auditoría

El control audita los claims actuales y decide su admisión. Una comparación puede
demostrar que dos cosas no son equivalentes: eso impide una sustitución incorrecta,
no implementa la capability que falta. Los20núcleos CONDITIONED y el candidato
conservan sus estados. No hay21promociones ni21módulos empresariales completos;
la implementación sigue AUTHORED.

V280/V293/V365 identificaron correctamente integraciones pendientes. Este sucesor
completa el contraste por claim con ocho snapshots primarios, fuentes exactas
Go1.26.8, pruebas comparativas nuevas, suites actuales y composición real.
No convierte la ausencia de un módulo del perfil en una función ausente ni elimina
los pendientes funcionales de sus owners. FAIL385 conserva su alcance de integración
más amplio; este expediente resuelve trazabilidad/admisión del claim estrecho de
TEST02, no suple integración empresarial.

El oracle se conserva: cada claim tiene fuente estrecha, licencia, revisión,
limitaciones y prueba; ningún CANDIDATE se promueve por reputación.
No se usa un G0–G8 parcialmente condicionado como promoción ELITE_REFERENCE/REUSABLE_PACK.

## Resultado ejecutado

228tests ordinarios y306semillas explícitas PASS, sin skips/fallos;
23campañas nativas de10s/target, 39811771ejecuciones observadas, PASS.
Vet/build PASS.42archivos reconstruidos coinciden con el árbol probado.
20packs/40archivos se componen y pasan sus suites; GO-OBSERVABILITY-CORE sigue
rechazado por el compositor incluso con acknowledgeConditions=true y no deja
archivos de producto.

El grafo ejecutado contiene stdlib y21paquetes locales, sin dependencias externas.
Fuzz prueba presupuestos finitos, no exhaustividad, races ni seguridad integral.
Los tests nuevos contrastan NPS con población tabulada, fórmula SRE y atributos
aceptados con JSONHandler oficial; mantienen explícita la no equivalencia CLDR
y la política SEO propia. No pretenden probar un sistema de negocio completo.

Cambian seis archivos de tests y sólo comentarios en dos implementaciones;
no cambia comportamiento/API de producto. Los21packs incrementan patch de metadata
y registran Go1.26.8 como runtime efectivamente probado. Sus claims, licencias,
admisión y selección integral permanecen intactos. FAIL707 corrige comentarios
que implicaban ventanas SRE/lecturas de stores que no existen.

## Autoridades y licencias

Go1.26.8: los bytes del toolchain, LICENSE, especificación y archivos de stdlib
están fijados en el registro. V281 conserva adquisición/distribución oficial y
suites de origen. La página web actual de especificación declara1.27; se utiliza
el snapshot local1.26, no se hereda su versión. BSD-3-Clause sólo licencia Go;
LICENSE.md gobierna los42bloques AUTHORED para uso del dueño del workspace.
No se concede redistribución pública de la biblioteca ni se copian módulos
empresariales por citar go.dev.

Comparaciones primarias consultadas2026-09-10, con SHA del snapshot HTML conservado:

- [Bain: cálculo NPS](https://www.netpromotersystem.com/about/measuring-your-net-promoter-score/):
  clases y resta de porcentajes, no validez de muestra o proceso de encuesta.
- [Google SRE: alertas sobre SLO](https://sre.google/workbook/alerting-on-slos/):
  fórmula de burn rate; presupuestos acumulados no implementan ventanas ni operación.
- [Unicode CLDR: plurales](https://cldr.unicode.org/index/cldr-spec/plural-rules):
  francés0 pertenece a one; la regla binaria local retorna other y no admite equivalencia general.
- [Google: sitemap](https://developers.google.com/search/docs/crawling-indexing/sitemaps/build-sitemap),
  [títulos](https://developers.google.com/search/docs/appearance/title-link) y
  [descripciones](https://developers.google.com/search/docs/appearance/snippet):
  formatos/publicación son otra frontera y60/160es política local, no límite oficial.
- [OpenTelemetry: modelo](https://opentelemetry.io/docs/specs/otel/metrics/data-model/) y
  [datos sensibles](https://opentelemetry.io/docs/security/handling-sensitive-data/):
  contraste de semántica/fronteras; no admiten el logger candidato ni garantizan reconocer toda PII.

No se adquieren nuevos SDKs, se copian datos CLDR ni se redistribuyen textos de esas
páginas. El registro de fuentes aporta comparación, nunca licencia por asociación.
Cada regla de negocio sigue siendo AUTHORED y corresponde al owner del proyecto.

## Dictamen por núcleo

| Núcleo | Comportamiento verificado | Comparación / límite | Condición obligatoria y owner | Tests |
|---|---|---|---|---|
| GO-DASHBOARDS-CORE | Agregados aportados por caller; NPS tabulado, media de ratings y validación de mes/conteos. | No lee PostgreSQL ni demuestra pertenencia de Input. Queries de pedidos/leads/stock no equivalen a estos KPIs. | Origen y autorización de cada valor, moneda/período coherentes, tamaño máximo de arrays. T2802/T2804 | 8 |
| GO-FX-CORE | Producto binary64 positivo/finito por tasa por tenant/par; no representación decimal. | Unidades menores o monedas del owner no realizan esta conversión. El helper tampoco tiene rates de proveedor, vigencia, asientos o redondeo monetario. | Fuente/timestamp de tasa, política monetaria, autenticación y consistencia con ledger antes de uso financiero. T2802 | 9 |
| GO-GIFT-CARDS-CORE | Saldo int64 local; canje único por tenant y sin sobregiro; duplicado rechaza, no replay de éxito. | Intento de pago no es saldo de tarjeta. Este Store no ofrece transacción durable, reversión ni checkout. | Autorización, moneda, almacenamiento/transacción, reversión y conciliación en owner financiero. T2802 | 9 |
| GO-HELP-CENTER-CORE | Estado de artículos local, versión optimista, búsquedas substring sobre publicados. | Ayuda versionada del journey no es un CMS general; este Store no autentica actores ni publica una web. | Autorización editorial, persistencia, render seguro, búsqueda/orden y versión de ayuda del consumidor. T2804 | 12 |
| GO-I18N-CORE | Fallback exacto/base/default/key y plural binario uno/otro. | Atributo HTML lang no implementa catálogo. Caso francés cero contradice equivalencia del helper con CLDR, demostrado por prueba. | Sólo idiomas/números compatibles con regla binaria; motor CLDR admitido y journeys localizados si se requiere más. T2804 | 9 |
| GO-LOYALTY-CORE | Ledger local int64 por tupla, movimientos únicos y crédito con overflow rechazado. | Pedido/pago durable no acredita puntos. El core no enlaza compra, reversión, política ni persistencia. | Reglas aprobadas, actor/tenant autorizados, enlace durable compra/movimiento, reversión e idempotencia de negocio. T2802 | 10 |
| GO-MARKETING-CORE | Campaña y acuse local scoped, deduplicación de destinatarios y completado local. | Leads/reporting/adapter no constituyen drip o campaña. Send no llama a red ni reclama lease ni impone ScheduledAt. | Consentimiento/supresión, tiempo de envío, reserva durable y provider receipt/fence; mantener datos privados. T2805/T2807 | 12 |
| GO-OBSERVABILITY-CORE | Política de eventos cerrada, JSONHandler oficial, contador e histograma local con límites. | Igualdad de atributos aceptados con slog no es SDK OTel. V372 usa su runtime admitido y no este core candidato; no hereda aquel PASS. | CANDIDATE: prohibida composición productiva. Readmitir sólo con política/integración/carga/retención y gates sobre este artefacto exacto. T2809 | 29 |
| GO-ONBOARDING-CORE | Checklist local con deduplicación, orden de alta y completado repetible. | Organización/acuerdo/login no equivale a onboarding de persona. El core no exige prerequisites ni evalúa o activa franquiciados. | Autorización, pasos/política aprobados, persistencia y activación/competencia evaluadas en el owner. T2803/T2804 | 9 |
| GO-PAYROLL-CORE | Aritmética local int64 de deducciones fijas/bps, prioridad de errores y truncado individual. | Ledger contable no calcula nómina. Execute no valida legalidad, liquida ni paga; duplicado rechaza, no replay de éxito. | Jurisdicción y reglas laborales/fiscales aprobadas, moneda, recibos, revisión, ledger durable y autorización. T2802 | 9 |
| GO-POS-CORE | Total/vuelto int64 con overflow y tender insuficiente rechazados, venta local por tupla. | Owners ya gobiernan pedidos/pagos/fiscalidad; helper no los llama ni implementa caja/dispositivo/offline. No sustituirlos. | Precio de fuente autorizada, moneda/cobro/efectos durables y UX de caja requeridos por journey. T2802/T2804 | 11 |
| GO-PROMOTIONS-CORE | Cálculo fijo/bps local y consumo de contador con máximo y vencimiento. | Price book no es cupón; core Redeem no conoce pedido o clave de replay ni revierte un descuento. | Reglas de acumulación/redondeo, autoridad comercial, vínculo durable pedido/cupón y reserva/reversión. T2802 | 9 |
| GO-REFERRALS-CORE | Tuplas tenant/identidad, no auto-referido y estados locales en secuencia. | Origen del lead no es referido. Convert/Reward locales no prueban conversión ni emiten recompensa. | Prueba de conversión/identidad, antifraude, regla aprobada, persistencia y recompensa durable por owner. T2802/T2805 | 10 |
| GO-REMINDERS-CORE | Schedule/Due/MarkSent con tenant explícito; Due snapshot y acuse local idempotente. | Jobs/fence son owners durables pero no llaman a este core. MarkSent no es receipt ni Due es lease; no sustituir worker. | Scheduling durable/cancelación, autorización/consentimiento, fence/receipt/reconciliación y presupuesto del provider. T2805/T2809 | 11 |
| GO-REVIEWS-CORE | Ratings válidos, unicidad por autor/sujeto/tenant, estados y promedio de aprobadas. | Identidad CRM no es reseña moderada; el core no valida compra, elegibilidad, antiabuso o publicación. | Actor/moderador autorizados, compra si se requiere, protección antiabuso, persistencia y render/publicación seguros. T2802/T2804 | 9 |
| GO-SEO-CORE | URLs por tenant, robots sin CR/LF/NUL, JSON-LD y política local 60/160. | Metadata web no prueba sitemap/robots/JSON-LD. Este core retorna URLs, no sitemap XML; longitudes son locales, no límites Google. | Dominio autorizado, serialización/formato/publicación y tests de URLs/canonical/robots del frontend; no garantía de indexación. T2804 | 14 |
| GO-SLO-CORE | Fórmulas de disponibilidad/presupuesto/burn de acumuladores y comparación booleana de dos budgets. | Reglas/alertas V372 del runtime no importan este Budget. Fórmula coincide; no hay ventanas temporales, ingestion o paginación en el helper. | SLI/ventanas alineadas y acotadas, thresholds validados, capacidad int64, señales y respuesta probadas sobre el runtime real. T2809 | 10 |
| GO-SOCIAL-POSTING-CORE | Agenda scoped con estados locales scheduled/posted/failed. | Reporting read-only no publica. MarkPosted/Failed sólo cambia memoria; no prueba proveedor, lease, cuotas o receipt. | Adapter oficial write, aprobación, cuota, fencing, receipt/reconciliación y revocación antes de cualquier efecto. T2805 | 11 |
| GO-SURVEYS-CORE | Respuestas únicas por tupla, clases NPS 0..10 y fórmula contrastada con Bain. | Reporting no captura respuestas. Fórmula no demuestra cuestionario, consentimiento, representatividad o invitación. | Autorización/consentimiento, encuesta y población pertinentes, privacidad/retención, captura durable y UX. T2802/T2804 | 9 |
| GO-WAITLIST-CORE | Posición FIFO entre pending; estados locales, Seat permite pending o notified. | Agenda/capacidad no equivale a espera. Notify no envía; no reserva cupo, oferta/expiración ni recupera después de reinicio. | Política de turnos aprobada, reserva/concurrencia durable, oferta/expiración y notificación con receipt. T2802 | 10 |
| GO-WARRANTY-CLAIMS-CORE | Estados locales de reclamo y notas requeridas antes de rechazar/resolver. | Casos/devoluciones durables no prueban términos de garantía. Core no decide elegibilidad ni ejecuta reembolso/reemplazo. | Términos/regulación y evidencia de compra, autorización, transición durable, efecto único y recovery en owner de servicio. T2802/T2804 | 8 |

## Registro exacto de evidencia

Los tests están ligados a cada paquete y a los hashes de sus dos fuentes. La matriz G0–G8 conserva CONDITIONED y NOT_PROMOTED donde corresponde. Las referencias de owner identifican fronteras para comparar, no evidencia de equivalencia completa.

```json
{
  "schema_version": 1,
  "audit": "TEST-02",
  "scope": "Current narrow claims, source attribution, reference tests and admission decisions; not full business integration or production assurance.",
  "summary": {
    "cores": 21,
    "source_files": 42,
    "ordinary_tests": 228,
    "fuzz_targets": 23,
    "explicit_seed_cases": 306,
    "fuzz_executions": 39811771,
    "composed_reference_packs": 20,
    "composed_reference_files": 40,
    "candidate_composition_rejected": true,
    "external_packages_in_isolated_graph": 0,
    "enterprise_equivalence_claims_approved": 0,
    "global_admission_promotions": 0,
    "all48_action_oracles_unchanged": true
  },
  "sources": [
    {
      "id": "go-1.26.8",
      "revision": "go1.26.8",
      "url": "https://go.dev/dl/go1.26.8.windows-amd64.zip",
      "provenance_evidence": "reconstruction_evidence/OFFICIAL_STRUCTURED_LOGGING_AUDIT_V281.md",
      "license": "BSD-3-Clause",
      "scope": "Exact runtime/language/stdlib only; none of the 21 business implementations is Go/Google upstream source.",
      "files": [
        {
          "path": "VERSION",
          "sha256": "25b231b380ff42bb887b29e7c2eb176167820f3416665cd4fbc2d3b116bc417d"
        },
        {
          "path": "LICENSE",
          "sha256": "911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad"
        },
        {
          "path": "doc/go_spec.html",
          "sha256": "e5806989624c1e9dd32bddaedd2f70a1dd208295cb6d0ff72b02625b2cd6f592"
        },
        {
          "path": "src/sync/mutex.go",
          "sha256": "3dec4264678744c4ce15a05ba0cc88e5e8b155f071b71d74eb7ed229ca72c114"
        },
        {
          "path": "src/sync/rwmutex.go",
          "sha256": "4ecc9b2fa7771cbe37e1670e4caefd7a28939683bfe1614fcf5e9d545261d17f"
        },
        {
          "path": "src/math/big/int.go",
          "sha256": "7689a5447c12d8c4cb350697cd7cdb4c827b3bdd3f0fb875441588cde3b18676"
        },
        {
          "path": "src/log/slog/json_handler.go",
          "sha256": "a10503560615180e018a44a0ebd5ab0c4d0a79f74ab7ee0248ca24277da259df"
        },
        {
          "path": "src/encoding/json/encode.go",
          "sha256": "4012d78aa8236f3067a6330ddf2f719c7e59b262fcf4fff5e51b4abd1ccba235"
        },
        {
          "path": "src/time/time.go",
          "sha256": "e8b6ea9e6e7cb890aeed1d24ee7163a2237b6fa78ffc564efb8688dba717f02a"
        },
        {
          "path": "bin/go.exe",
          "sha256": "21761eceb9302062c9623fb699f332c8c7fe000f15f70efe8da01a2cfbbc16b9"
        }
      ]
    },
    {
      "id": "bain-nps",
      "url": "https://www.netpromotersystem.com/about/measuring-your-net-promoter-score/",
      "observed_url": "https://www.netpromotersystem.com/about/measuring-your-net-promoter-score/",
      "observed_at": "2026-09-10T15:24:12.287591+00:00",
      "sha256": "997beb3ce69d8a2423581aadef7360c48b3b5dd8f83f110d7aebb6a86f0d7d28",
      "bytes": 67726,
      "license_scope": "Documentation comparison only; no text/code/data redistributed or relicensed."
    },
    {
      "id": "google-sre",
      "url": "https://sre.google/workbook/alerting-on-slos/",
      "observed_url": "https://sre.google/workbook/alerting-on-slos/",
      "observed_at": "2026-09-10T15:24:12.990770+00:00",
      "sha256": "5a2195ae9280d92489ec6c5be9781ef191d73137514da7c716ebabf5479cfcf5",
      "bytes": 67534,
      "license_scope": "Documentation comparison only; no text/code/data redistributed or relicensed."
    },
    {
      "id": "cldr",
      "url": "https://cldr.unicode.org/index/cldr-spec/plural-rules",
      "observed_url": "https://cldr.unicode.org/index/cldr-spec/plural-rules",
      "observed_at": "2026-09-10T15:24:13.114090+00:00",
      "sha256": "e750efdab5ffaca7567269ffd7278cc474d7a1c4045bad8db9ce58d8e8697427",
      "bytes": 27010,
      "license_scope": "Documentation comparison only; no text/code/data redistributed or relicensed."
    },
    {
      "id": "google-sitemap",
      "url": "https://developers.google.com/search/docs/crawling-indexing/sitemaps/build-sitemap",
      "observed_url": "https://developers.google.com/search/docs/crawling-indexing/sitemaps/build-sitemap",
      "observed_at": "2026-09-10T15:24:13.961483+00:00",
      "sha256": "f324ff601bb86aa464f93166c06e284dfb14b8d47ddf5ad28b622bf68513ec9e",
      "bytes": 190914,
      "license_scope": "Documentation comparison only; no text/code/data redistributed or relicensed."
    },
    {
      "id": "google-title",
      "url": "https://developers.google.com/search/docs/appearance/title-link",
      "observed_url": "https://developers.google.com/search/docs/appearance/title-link",
      "observed_at": "2026-09-10T15:24:15.387934+00:00",
      "sha256": "81ba2e4bc9c4654bb547cad3767b7bbb533fbb8e60db07d4ee0fe8beb30dafe4",
      "bytes": 188604,
      "license_scope": "Documentation comparison only; no text/code/data redistributed or relicensed."
    },
    {
      "id": "google-snippet",
      "url": "https://developers.google.com/search/docs/appearance/snippet",
      "observed_url": "https://developers.google.com/search/docs/appearance/snippet",
      "observed_at": "2026-09-10T15:24:16.073305+00:00",
      "sha256": "ad90ab93f237fbbc4ded0078e24ec81cd9b943cbf96868a84d4f45d3afeaed34",
      "bytes": 181837,
      "license_scope": "Documentation comparison only; no text/code/data redistributed or relicensed."
    },
    {
      "id": "otel-data-model",
      "url": "https://opentelemetry.io/docs/specs/otel/metrics/data-model/",
      "observed_url": "https://opentelemetry.io/docs/specs/otel/metrics/data-model/",
      "observed_at": "2026-09-10T15:24:17.342188+00:00",
      "sha256": "1ea2b3d384a61a040fe2fc9c953b818c244c81e5aff1cac40ede13385064fdea",
      "bytes": 166644,
      "license_scope": "Documentation comparison only; no text/code/data redistributed or relicensed."
    },
    {
      "id": "otel-sensitive",
      "url": "https://opentelemetry.io/docs/security/handling-sensitive-data/",
      "observed_url": "https://opentelemetry.io/docs/security/handling-sensitive-data/",
      "observed_at": "2026-09-10T15:24:17.997797+00:00",
      "sha256": "2a8b1b8ca989423ec8fb841bfffba97b6cd76c77f3da5ae46aaf2e9698b8f57b",
      "bytes": 413272,
      "license_scope": "Documentation comparison only; no text/code/data redistributed or relicensed."
    }
  ],
  "cores": [
    {
      "path": "implementation_packs/GO_DASHBOARDS_CORE.md",
      "id": "GO-DASHBOARDS-CORE",
      "version": "0.1.2",
      "sha256": "cb17f0d2093ce54d3fc554de5119b4ed51dda5076355189a79f5d5379cd5c00e",
      "files": [
        {
          "path": "internal/dashboards/dashboards.go",
          "sha256": "c4ec3fb2c2e63340d1d6b0963a3d8ed30f1a57bf9423c0ab54f146db1c02773a",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/dashboards/dashboards_test.go",
          "sha256": "235dafdda91235e9eac8c99a0334d3c0078dd9d8437b75a1f8249e84960efe34",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8",
        "bain-nps"
      ],
      "reviewed_behavior": "Agregados aportados por caller; NPS tabulado, media de ratings y validación de mes/conteos.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_ENTERPRISE_QUERY_API.md",
          "sha256": "84601b5f6b3917f4f540ffb90d9d569b6cf3e9a8e0fa4955a736c2c4e9a38f70"
        },
        {
          "path": "implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md",
          "sha256": "9f9aebf972ed734459ded6de84dbf569f1e36d05957c166be3d46ace2c554646"
        }
      ],
      "comparison": "No lee PostgreSQL ni demuestra pertenencia de Input. Queries de pedidos/leads/stock no equivalen a estos KPIs.",
      "conditions": "Origen y autorización de cada valor, moneda/período coherentes, tamaño máximo de arrays.",
      "roadmap_owner": "T2802/T2804",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Referencia AUTHORED de agregación pura de KPIs suministrados por el caller: valida forma YYYY-MM y mes 01-12, conteos y scores; etiqueta tenant sin demostrar procedencia o autorización de los datos.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestBuildSnapshot",
        "TestEmptyReviewsAndNPS",
        "TestInvalidCalendarMonth",
        "TestInvalidInput",
        "TestNegativeAggregatesAndInvalidScores",
        "TestSnapshotOwnsAggregatedValues",
        "TestSourceNPSGoldenClassesV374",
        "TestValidCalendarMonthsPreserveShape"
      ],
      "reference_test_count": 8,
      "seed_case_count": 18,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_FX_CORE.md",
      "id": "GO-FX-CORE",
      "version": "0.1.3",
      "sha256": "329c2609f493666c204e2a0c36803f9b583e60eccd4daec781d805660ff50ab5",
      "files": [
        {
          "path": "internal/fx/fx.go",
          "sha256": "d7e7aa5764c35a72dd99856f4cb93e3eeef67ac3b15740d7e63e7e43b0318bf9",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/fx/fx_test.go",
          "sha256": "8c9a9bb0b45c284c6c171d9f1306d971a23ac0f120fbc8fd5206306fc1d8b557",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Producto binary64 positivo/finito por tasa por tenant/par; no representación decimal.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md",
          "sha256": "0e7b07bba6c4046ee369f52f2a2964990b54b4510645e29a8e40259edc4c7ed3"
        },
        {
          "path": "implementation_packs/GO_ENTERPRISE_ACCOUNTING_LEDGER_API.md",
          "sha256": "d6a788b518904c56b6c2c0a35c745827964f9d3caafa9567719784dec36b9daa"
        }
      ],
      "comparison": "Unidades menores o monedas del owner no realizan esta conversión. El helper tampoco tiene rates de proveedor, vigencia, asientos o redondeo monetario.",
      "conditions": "Fuente/timestamp de tasa, política monetaria, autenticación y consistencia con ledger antes de uso financiero.",
      "roadmap_owner": "T2802",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Materializa un helper AUTHORED en memoria para multiplicar montos binary64 por tasas positivas finitas, aislado por tenant/par, con rechazo sin tasa y error ante resultado no finito o cero. No implementa política monetaria, proveedor de tasas, redondeo decimal ni ledger durable.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestConvert",
        "TestDenyByDefault",
        "TestInvalidInput",
        "TestInvalidRate",
        "TestNonFiniteAmountRejected",
        "TestNonFiniteRateDoesNotReplacePrior",
        "TestNonRepresentableProductRejected",
        "TestRepresentablePositiveResultsPreserved",
        "TestTenantIsolation"
      ],
      "reference_test_count": 9,
      "seed_case_count": 14,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_GIFT_CARDS_CORE.md",
      "id": "GO-GIFT-CARDS-CORE",
      "version": "0.1.3",
      "sha256": "f885c813feda3fb78bf75291cef76b23ab85e23396b8fff5c94e3f488e66d453",
      "files": [
        {
          "path": "internal/giftcards/giftcards.go",
          "sha256": "a088ad1335038bc2f5e6658197b1aa8d3e8a9456cba5d9527523f6874e51a878",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/giftcards/giftcards_test.go",
          "sha256": "3d3d247d2f3531707ccc4088e97aa2c89d11c995c6a6b9344ceece566670fe7a",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Saldo int64 local; canje único por tenant y sin sobregiro; duplicado rechaza, no replay de éxito.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md",
          "sha256": "0e7b07bba6c4046ee369f52f2a2964990b54b4510645e29a8e40259edc4c7ed3"
        }
      ],
      "comparison": "Intento de pago no es saldo de tarjeta. Este Store no ofrece transacción durable, reversión ni checkout.",
      "conditions": "Autorización, moneda, almacenamiento/transacción, reversión y conciliación en owner financiero.",
      "roadmap_owner": "T2802",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Referencia AUTHORED en memoria para emitir saldo positivo y canjear sin sobreconsumo, con identificador no vacío de canje único por tenant. Rechaza duplicados; no hace top-up, replay exitoso, persistencia, autorización ni checkout productivo.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestBlankRedeemIdentityRejected",
        "TestDuplicateAndInvalid",
        "TestFailedRedemptionCanRetryAndConcurrentDebit",
        "TestIdempotencyByRedeemID",
        "TestIssueRedeemBalance",
        "TestRedeemIdentityScopedByTenant",
        "TestRedeemInsufficientFailsClosed",
        "TestRedemptionTupleDoesNotAlias",
        "TestTenantIsolation"
      ],
      "reference_test_count": 9,
      "seed_case_count": 6,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_HELP_CENTER_CORE.md",
      "id": "GO-HELP-CENTER-CORE",
      "version": "0.1.2",
      "sha256": "7669198e3df9d82e4fa76ddb53ec730ac642a6133ffae1ed93547e8e1d60d47e",
      "files": [
        {
          "path": "internal/helpcenter/helpcenter.go",
          "sha256": "5c8e958ab4ef1262fbfed2583b99f23d50cce20f27715031cb347e2ac47cd844",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/helpcenter/helpcenter_test.go",
          "sha256": "86b8d27917677867cfdc5f9c7536e0bfb4ddb3196d3756825c046861246ecbc6",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Estado de artículos local, versión optimista, búsquedas substring sobre publicados.",
      "comparison_owners": [
        {
          "path": "implementation_packs/TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md",
          "sha256": "7f2892e0e436bfd4fffe665dc81bdfbd408c9507ef9acd11f2cc6bdf2f358be8"
        },
        {
          "path": "implementation_packs/PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md",
          "sha256": "5fc99e1797556b86dc14a10efc482b3582958ddce13dda99f5bfa580847095af"
        }
      ],
      "comparison": "Ayuda versionada del journey no es un CMS general; este Store no autentica actores ni publica una web.",
      "conditions": "Autorización editorial, persistencia, render seguro, búsqueda/orden y versión de ayuda del consumidor.",
      "roadmap_owner": "T2804",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Referencia AUTHORED de artículos en memoria por tupla tenant/id, concurrencia optimista acotada a int64 y estados draft→published→archived; búsqueda local substring, sin adapter externo ni autorización del actor.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestBadTransition",
        "TestConcurrentVersionHasSingleWinner",
        "TestCreatePublishArchive",
        "TestDistinctArticleTuples",
        "TestForeignArticleTransitionsRejected",
        "TestForeignArticleUpdateRejected",
        "TestInvalidAndTenantIsolation",
        "TestOptimisticConcurrency",
        "TestPublishedSearchCopyAndArchivedUpdate",
        "TestSearchPublishedOnly",
        "TestVersionBoundaryAndErrorPrecedence",
        "TestVersionExhaustionRejected"
      ],
      "reference_test_count": 12,
      "seed_case_count": 16,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_I18N_CORE.md",
      "id": "GO-I18N-CORE",
      "version": "0.1.2",
      "sha256": "6eb6c2b27f0f0d7506787eefc4f37192356a4d054934430c0ae6785b1211a8e7",
      "files": [
        {
          "path": "internal/i18n/i18n.go",
          "sha256": "818df4f6310c82a45520c2fb858745702c590a9c4d37db14decdd3ecc44b1bac",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/i18n/i18n_test.go",
          "sha256": "53fe04950e4f95d65315983964a90375c83160005101d57947db1924a06e1092",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8",
        "cldr"
      ],
      "reviewed_behavior": "Fallback exacto/base/default/key y plural binario uno/otro.",
      "comparison_owners": [
        {
          "path": "implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md",
          "sha256": "9f9aebf972ed734459ded6de84dbf569f1e36d05957c166be3d46ace2c554646"
        }
      ],
      "comparison": "Atributo HTML lang no implementa catálogo. Caso francés cero contradice equivalencia del helper con CLDR, demostrado por prueba.",
      "conditions": "Sólo idiomas/números compatibles con regla binaria; motor CLDR admitido y journeys localizados si se requiere más.",
      "roadmap_owner": "T2804",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Catálogo AUTHORED con fallback exacto→base→default→key y regla binaria n==1 one/else other; rechaza formas vacías. No implementa CLDR, validación completa de tags ni interpolación.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestBaseLocale",
        "TestConcurrentCatalogReadWrite",
        "TestInvalidKeyAndLocale",
        "TestPlural",
        "TestPluralErrorPrecedence",
        "TestPluralRejectsBlankWithoutReplacing",
        "TestRejectedPluralDoesNotCreateKey",
        "TestSourceCLDRNonEquivalenceV374",
        "TestTWithFallback"
      ],
      "reference_test_count": 9,
      "seed_case_count": 30,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_LOYALTY_CORE.md",
      "id": "GO-LOYALTY-CORE",
      "version": "0.1.2",
      "sha256": "001a1eb4dbc47b2aa1f564ba648c15b586ebd731f6071f0676c38f74618f5b24",
      "files": [
        {
          "path": "internal/loyalty/loyalty.go",
          "sha256": "14a339db71422e0c247a74dcd73998a62c1c7e1da4549acb0e3f0ada317cabf1",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/loyalty/loyalty_test.go",
          "sha256": "105a99ffc6c9a12fb8ec0ae7c1627aec807ea68ad1eb2d125686a8081274d3a1",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Ledger local int64 por tupla, movimientos únicos y crédito con overflow rechazado.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md",
          "sha256": "0e7b07bba6c4046ee369f52f2a2964990b54b4510645e29a8e40259edc4c7ed3"
        }
      ],
      "comparison": "Pedido/pago durable no acredita puntos. El core no enlaza compra, reversión, política ni persistencia.",
      "conditions": "Reglas aprobadas, actor/tenant autorizados, enlace durable compra/movimiento, reversión e idempotencia de negocio.",
      "roadmap_owner": "T2802",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Referencia AUTHORED en memoria de puntos int64: earn/burn con saldo no negativo ni overflow, cuentas por tupla tenant/cliente y movimientos únicos por tenant. Historial devuelto como copia; no ledger durable, autorización ni política de recompensas.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestAccountTupleIsolation",
        "TestBurnInsufficientFailsClosed",
        "TestConcurrentDuplicateCreditAndHistoryCopy",
        "TestCreditOverflowAtomicAndRetry",
        "TestDuplicateEntryRejected",
        "TestEarnBurnBalance",
        "TestEntryIdentityScopedByTenant",
        "TestHistoryImmutableOrder",
        "TestInvalidPointsRejected",
        "TestTenantIsolation"
      ],
      "reference_test_count": 10,
      "seed_case_count": 6,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_MARKETING_CORE.md",
      "id": "GO-MARKETING-CORE",
      "version": "0.2.1",
      "sha256": "d8cc6c006c3bb739e1350eb9619b39fc6d3cd6154f74c8b5dbb39d047b126d80",
      "files": [
        {
          "path": "internal/marketing/marketing.go",
          "sha256": "1a1af2821d5b9bce83b4454842d8af0ad837b5dbb3f23b1705fe802c8a1e811f",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/marketing/marketing_test.go",
          "sha256": "2e90bd08651a7eed724bb824b744bb0c05e5139db294fb397a99d1e7fab803af",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Campaña y acuse local scoped, deduplicación de destinatarios y completado local.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md",
          "sha256": "5ff3325a3d1e2e33f79f80034b92fd771e97c104416e96d67e902b64c8ebc242"
        },
        {
          "path": "implementation_packs/PYTHON_META_ADS_REPORTING_ADAPTER.md",
          "sha256": "2ed2981302f3c7c6caf69eaead0205cd9cf2e54b6799b719664241cb6339e53f"
        },
        {
          "path": "implementation_packs/PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md",
          "sha256": "5fc99e1797556b86dc14a10efc482b3582958ddce13dda99f5bfa580847095af"
        }
      ],
      "comparison": "Leads/reporting/adapter no constituyen drip o campaña. Send no llama a red ni reclama lease ni impone ScheduledAt.",
      "conditions": "Consentimiento/supresión, tiempo de envío, reserva durable y provider receipt/fence; mantener datos privados.",
      "roadmap_owner": "T2805/T2807",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Registro AUTHORED de campañas en memoria, consultas scoped por tenant y acuses locales por destinatario; no implementa envío al proveedor.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestCampaignIDsIndependentPerTenant",
        "TestConcurrentCampaignAndAcknowledgementScope",
        "TestCreateDeduplicatesRecipients",
        "TestDeliveryIdentityBoundary",
        "TestDueCanRestrictCallerTenant",
        "TestDueRecipientsDoNotAliasStore",
        "TestDuplicateCampaignRejected",
        "TestInvalidCampaign",
        "TestScopedCampaignStateAndScheduleSemantics",
        "TestSendAndCompletion",
        "TestSendDedupPerRecipient",
        "TestSendUnknownRecipient"
      ],
      "reference_test_count": 12,
      "seed_case_count": 12,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_OBSERVABILITY_CORE.md",
      "id": "GO-OBSERVABILITY-CORE",
      "version": "0.3.3",
      "sha256": "fab3317ce4e85fe49b69217d7642b2630d8f1f042d5a25ccef35b58f097ae681",
      "files": [
        {
          "path": "internal/observability/observability.go",
          "sha256": "0425ec30fa22dbec4f579274281e2668e12638eb7252d74eaf6772a0c801b2c5",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/observability/observability_test.go",
          "sha256": "0e817dbcd2c97d6622416bc83038464b7708b5512955bc4f7605637e745a9ae5",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CANDIDATE",
      "source_refs": [
        "go-1.26.8",
        "otel-data-model",
        "otel-sensitive"
      ],
      "reviewed_behavior": "Política de eventos cerrada, JSONHandler oficial, contador e histograma local con límites.",
      "comparison_owners": [
        {
          "path": "implementation_packs/WINDOWS_REFERENCE_TELEMETRY_RUNTIME.md",
          "sha256": "82aa19096bb60a60ac834a33185f2dea6d79637bfdf9d009b76427131ca22ecd"
        },
        {
          "path": "implementation_packs/GO_OFFICIAL_RETURN_REFUND_WORKER.md",
          "sha256": "8d81308a434615f9d8e6de9225e0b2b0aaaae8e3186dff4c00f46b9294221ae0"
        }
      ],
      "comparison": "Igualdad de atributos aceptados con slog no es SDK OTel. V372 usa su runtime admitido y no este core candidato; no hereda aquel PASS.",
      "conditions": "CANDIDATE: prohibida composición productiva. Readmitir sólo con política/integración/carga/retención y gates sobre este artefacto exacto.",
      "roadmap_owner": "T2809",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Candidato local: logger de política cerrada, contador monotónico e histograma validado con overflow separado y snapshot coherente. No detector universal de PII, SDK OTel ni observabilidad integral.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestConcurrentJSONRecords",
        "TestCounter",
        "TestCounterConcurrentAcceptedAndRejected",
        "TestCounterConcurrentOverflow",
        "TestCounterRejectNegative",
        "TestCounterRejectOverflow",
        "TestCounterTryIncStatusAndRecovery",
        "TestDeliveryFailureContainedWithoutRetry",
        "TestHistogramConcurrentSnapshots",
        "TestHistogramConfigurationAndSnapshot",
        "TestHistogramEmpty",
        "TestHistogramPercentile",
        "TestHistogramPercentileLargeCount",
        "TestHistogramRegressionCounterOverflow",
        "TestHistogramRegressionNonFinite",
        "TestHistogramRegressionOverflow",
        "TestHistogramRegressionZeroPercentile",
        "TestHistogramRejectAndRecovery",
        "TestJSONApprovedEventEndToEnd",
        "TestJSONDeliveryBoundaryFullWrite",
        "TestJSONDeliveryBoundaryRejectsIncompleteWrite",
        "TestLegacyUnconfiguredDoesNotDisclose",
        "TestLevelFiltering",
        "TestPolicySnapshotAndSinkSnapshot",
        "TestPolicyValidation",
        "TestReferenceClosedFileAndFreshDestination",
        "TestReferenceHTTPFileTelemetry",
        "TestRejectWholeEventWithoutFormatting",
        "TestSourceSlogAcceptedSerializationV374"
      ],
      "reference_test_count": 29,
      "seed_case_count": 17,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CANDIDATE_NO_COMPOSITION",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_ONBOARDING_CORE.md",
      "id": "GO-ONBOARDING-CORE",
      "version": "0.1.2",
      "sha256": "cd249cafc4f753ac328de11f5692ab7ecbad4679ad6304e5bd9ab5945407d6ac",
      "files": [
        {
          "path": "internal/onboarding/onboarding.go",
          "sha256": "57bc85e39c18cbea30ed7f3477160fdc0faf2f317ce807223f08bd9be2215e53",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/onboarding/onboarding_test.go",
          "sha256": "94506a1638046a14f30fce623a017ef986e18962815e6c42f40fd3f5a007ee72",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Checklist local con deduplicación, orden de alta y completado repetible.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_FULFILLMENT_SERVICE_FRANCHISE_API.md",
          "sha256": "b647b19aa840b1f6723fa9e349425f8d042084579da6cda4cace91c28bb021a6"
        },
        {
          "path": "implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md",
          "sha256": "353f05c1f2c942417d6406662e6742f5e930183f610093316b73dabffb05c183"
        }
      ],
      "comparison": "Organización/acuerdo/login no equivale a onboarding de persona. El core no exige prerequisites ni evalúa o activa franquiciados.",
      "conditions": "Autorización, pasos/política aprobados, persistencia y activación/competencia evaluadas en el owner.",
      "roadmap_owner": "T2803/T2804",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Materializa el onboarding del franquiciado: pasos ordenados, checklist con progreso y completado idempotente, tenant-scoped.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestAliasedOnboardingCannotReadOrComplete",
        "TestChecklistCopiesInputAndAllowsAnyKnownStep",
        "TestCompleteIdempotent",
        "TestConcurrentChecklistStartAndIdempotentCompletion",
        "TestDistinctOnboardingTuplesDoNotCollide",
        "TestProgressAndDone",
        "TestStartDedupAndInvalid",
        "TestTenantIsolation",
        "TestUnknownStepRejected"
      ],
      "reference_test_count": 9,
      "seed_case_count": 12,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_PAYROLL_CORE.md",
      "id": "GO-PAYROLL-CORE",
      "version": "0.1.3",
      "sha256": "4510fb7d4f1a672b64fadc62c77c470f775b1c1be6fd158a08bbe0655ef9d059",
      "files": [
        {
          "path": "internal/payroll/payroll.go",
          "sha256": "5a9b0f24dd5f18eb87f3285a715e0595a57bf8e42e866a644d6635f2469b6b9d",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/payroll/payroll_test.go",
          "sha256": "0d44af51784430e279fb9ab796079387b9f728f5b2c555fed829baa7049f8457",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Aritmética local int64 de deducciones fijas/bps, prioridad de errores y truncado individual.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_ENTERPRISE_ACCOUNTING_LEDGER_API.md",
          "sha256": "d6a788b518904c56b6c2c0a35c745827964f9d3caafa9567719784dec36b9daa"
        }
      ],
      "comparison": "Ledger contable no calcula nómina. Execute no valida legalidad, liquida ni paga; duplicado rechaza, no replay de éxito.",
      "conditions": "Jurisdicción y reglas laborales/fiscales aprobadas, moneda, recibos, revisión, ledger durable y autorización.",
      "roadmap_owner": "T2802",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Materializa el cálculo de nómina tenant-scoped config-driven: deducciones fijas o porcentuales inyectadas en runtime, net fail-closed (nunca negativo) e idempotente. Sin tablas impositivas hardcodeadas.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestConcurrentMaximumPayrollRunUniquePerTenant",
        "TestDuplicateAndTenantIsolation",
        "TestFixedDeductions",
        "TestFullRangePercentageNet",
        "TestInvalidDeduction",
        "TestInvalidDeductionPrecedenceAfterExcess",
        "TestNegativeNetCannotWrapToStoredSuccess",
        "TestNegativeNetRejected",
        "TestPercentDeduction"
      ],
      "reference_test_count": 9,
      "seed_case_count": 24,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_POS_CORE.md",
      "id": "GO-POS-CORE",
      "version": "0.2.1",
      "sha256": "340e5f461d975258827961da1330b16f0620259b0e7ff9546799462d0700155b",
      "files": [
        {
          "path": "internal/pos/pos.go",
          "sha256": "4d3be0f26614884e59ba2c052c5d2a29fcb5ce0200719c7da611e8750951ee41",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/pos/pos_test.go",
          "sha256": "ee7fb029c87abb828b911149f722c9f364558c62d4bf908f23be3cd32e30555c",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Total/vuelto int64 con overflow y tender insuficiente rechazados, venta local por tupla.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md",
          "sha256": "0e7b07bba6c4046ee369f52f2a2964990b54b4510645e29a8e40259edc4c7ed3"
        },
        {
          "path": "implementation_packs/GO_ARCA_FISCAL_ISSUANCE_API.md",
          "sha256": "2d8dd3e13e5ae38986cb7232cf345adb70778f94e968d83b3155bd8f21d2a83d"
        }
      ],
      "comparison": "Owners ya gobiernan pedidos/pagos/fiscalidad; helper no los llama ni implementa caja/dispositivo/offline. No sustituirlos.",
      "conditions": "Precio de fuente autorizada, moneda/cobro/efectos durables y UX de caja requeridos por journey.",
      "roadmap_owner": "T2802/T2804",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Referencia AUTHORED en memoria: total/vuelto con errores explícitos, registro de venta por tenant/id y copia de líneas; no implementa cobro ni checkout.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestCheckedArithmeticAndErrorPrecedence",
        "TestCompleteAndDuplicate",
        "TestCompletedSaleOwnsItsLines",
        "TestConcurrentSaleAndDebugRead",
        "TestDistinctSaleIdentityTuples",
        "TestInvalidSale",
        "TestNegativeTenderCannotWrapToPositiveChange",
        "TestOverflowingSaleRejectedBeforeRecord",
        "TestShortTenderFailsClosed",
        "TestTenantIsolation",
        "TestTotalAndChange"
      ],
      "reference_test_count": 11,
      "seed_case_count": 14,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_PROMOTIONS_CORE.md",
      "id": "GO-PROMOTIONS-CORE",
      "version": "0.1.3",
      "sha256": "79f76aec4fc7dd5776ad9b8e85a8cfeb05179e251d83fa7bbcb5075eb1a72412",
      "files": [
        {
          "path": "internal/promotions/promotions.go",
          "sha256": "81da2ae2b9743532f1bcbd2a691980ad71b797ec6dc8ed0a9d24cb1bdcda8467",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/promotions/promotions_test.go",
          "sha256": "837205419f3c2e4b11ceaf1810e1fac600d1600239305760de21394bc8c5e795",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Cálculo fijo/bps local y consumo de contador con máximo y vencimiento.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md",
          "sha256": "0e7b07bba6c4046ee369f52f2a2964990b54b4510645e29a8e40259edc4c7ed3"
        }
      ],
      "comparison": "Price book no es cupón; core Redeem no conoce pedido o clave de replay ni revierte un descuento.",
      "conditions": "Reglas de acumulación/redondeo, autoridad comercial, vínculo durable pedido/cupón y reserva/reversión.",
      "roadmap_owner": "T2802",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Materializa los cupones/descuentos tenant-scoped: percent o fixed, con expiración, monto mínimo, límite de uso y redención fail-closed.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestConcurrentCouponLimitAtMaximumAmount",
        "TestCouponBoundaryAndRejectedUsePreserved",
        "TestExpiredAndInactiveAndMinimum",
        "TestFullRangePercentageDiscount",
        "TestRedeemFixedCappedAtOrder",
        "TestRedeemPercent",
        "TestTenantIsolation",
        "TestUsageLimitExhausted",
        "TestValidateRejectsBadValues"
      ],
      "reference_test_count": 9,
      "seed_case_count": 13,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_REFERRALS_CORE.md",
      "id": "GO-REFERRALS-CORE",
      "version": "0.1.2",
      "sha256": "98c52e351453046d3c684f791630634d4b9b5db23a7ee710c3c794e3b5f2688a",
      "files": [
        {
          "path": "internal/referrals/referrals.go",
          "sha256": "23f3fb2be5f341caf748b8bcf364ffcff4e571e234549c84e8b312a0e9d83231",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/referrals/referrals_test.go",
          "sha256": "9314db0d1bf1730673bf51105d23959367309c671aa29d3abf9f7a493e399219",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Tuplas tenant/identidad, no auto-referido y estados locales en secuencia.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md",
          "sha256": "5ff3325a3d1e2e33f79f80034b92fd771e97c104416e96d67e902b64c8ebc242"
        },
        {
          "path": "implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md",
          "sha256": "69b0f715eab64ef831f510505ceb4b59d2a36444c7e2c557524eb44d424e6557"
        }
      ],
      "comparison": "Origen del lead no es referido. Convert/Reward locales no prueban conversión ni emiten recompensa.",
      "conditions": "Prueba de conversión/identidad, antifraude, regla aprobada, persistencia y recompensa durable por owner.",
      "roadmap_owner": "T2802/T2805",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Materializa el programa de referidos tenant-scoped: códigos únicos, sin auto-referido, un referido por persona y máquina de estados pending→converted→rewarded fail-closed.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestAbsentAliasedRefereeIsNotFound",
        "TestAliasedLookupCannotTransitionAnotherReferee",
        "TestConcurrentReferralLifecycleAtMostOneTransition",
        "TestCreateAndStateMachine",
        "TestDistinctReferralTuplesDoNotCollide",
        "TestDuplicateCodeAndReferee",
        "TestRejectedCreatesDoNotConsumeIdentity",
        "TestRewardRequiresConverted",
        "TestSelfReferralRejected",
        "TestTenantIsolation"
      ],
      "reference_test_count": 10,
      "seed_case_count": 12,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_REMINDERS_CORE.md",
      "id": "GO-REMINDERS-CORE",
      "version": "0.2.1",
      "sha256": "267837d492c81538dea6f0bb0cac6cf5a5be18a4396298db8acd05708ad47a91",
      "files": [
        {
          "path": "internal/reminders/reminders.go",
          "sha256": "0771a0a86db811051231a9094726803c4783210dae3e1964e45e9071e44bbd2e",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/reminders/reminders_test.go",
          "sha256": "0054fd9b5af5282a36a0a80dfa5d836c69d1b392bfb4858cc291fdb7a1def8ba",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Schedule/Due/MarkSent con tenant explícito; Due snapshot y acuse local idempotente.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_RELIABLE_ASYNC_WORKERS.md",
          "sha256": "8ab8136709da437dba73e7335805a349fd1dcb9cd61622165362b899ee087aad"
        },
        {
          "path": "implementation_packs/GO_PG_OUTBOUND_DELIVERY_FENCE.md",
          "sha256": "a3084eec92ce1a918c9240f815999724daf44d010fa53f45173e85173367b061"
        },
        {
          "path": "implementation_packs/PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md",
          "sha256": "5fc99e1797556b86dc14a10efc482b3582958ddce13dda99f5bfa580847095af"
        }
      ],
      "comparison": "Jobs/fence son owners durables pero no llaman a este core. MarkSent no es receipt ni Due es lease; no sustituir worker.",
      "conditions": "Scheduling durable/cancelación, autorización/consentimiento, fence/receipt/reconciliación y presupuesto del provider.",
      "roadmap_owner": "T2805/T2809",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Registro AUTHORED en memoria de recordatorios, con identidad (tenant,ID), Due(tenant) y MarkSent(tenant,ID) explícitos. Devuelve copias y conserva flag local; no reserva despacho ni demuestra envío, entrega exactamente una vez, persistencia o autorización.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestConcurrentScopedScheduleAndAcknowledgement",
        "TestDueAndMarkSent",
        "TestDueNeverReturnsOtherTenant",
        "TestDuplicateRejected",
        "TestMarkSentNotFound",
        "TestOtherTenantCannotMarkReminder",
        "TestRejectedSchedulePreservesIdentityAndDueBoundary",
        "TestSameReminderIDIndependentTenants",
        "TestScheduleRejectsInvalid",
        "TestScheduleRejectsPast",
        "TestScopedAcknowledgementAndDetachedDueReads"
      ],
      "reference_test_count": 11,
      "seed_case_count": 6,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_REVIEWS_CORE.md",
      "id": "GO-REVIEWS-CORE",
      "version": "0.1.2",
      "sha256": "d5cb5b438e1cccc6dbb6c6d913b84da687fe86591044cc25f0a9b869da325ee9",
      "files": [
        {
          "path": "internal/reviews/reviews.go",
          "sha256": "6dc3aa976fb612ee7157a009c7de4638ece3558b9e83f0eeccb4bf4df2bd908c",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/reviews/reviews_test.go",
          "sha256": "71b05f0a282a482d103219e3bbbf101cf78b033eef4cbb5453bed9c38d73e395",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Ratings válidos, unicidad por autor/sujeto/tenant, estados y promedio de aprobadas.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md",
          "sha256": "69b0f715eab64ef831f510505ceb4b59d2a36444c7e2c557524eb44d424e6557"
        }
      ],
      "comparison": "Identidad CRM no es reseña moderada; el core no valida compra, elegibilidad, antiabuso o publicación.",
      "conditions": "Actor/moderador autorizados, compra si se requiere, protección antiabuso, persistencia y render/publicación seguros.",
      "roadmap_owner": "T2802/T2804",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Materializa reseñas/ratings tenant-scoped con moderación: rating 1-5, una por autor por sujeto, estados pending→approved/rejected y promedio sólo sobre aprobadas.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestAliasedModerationCannotPublishOrReject",
        "TestApprovedSnapshotCannotMutateStore",
        "TestBadTransition",
        "TestConcurrentModerationAndDebugRead",
        "TestDistinctModerationTuplesDoNotCollide",
        "TestDuplicateAuthorRejected",
        "TestInvalidRating",
        "TestSubmitModerationAndAverage",
        "TestTenantIsolation"
      ],
      "reference_test_count": 9,
      "seed_case_count": 12,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_SEO_CORE.md",
      "id": "GO-SEO-CORE",
      "version": "0.1.2",
      "sha256": "8fb03f30829698ba859f96340837092f1cf8790de51e8605e6a77a499416109e",
      "files": [
        {
          "path": "internal/seo/seo.go",
          "sha256": "d5aecce729bce9f8ab3717d81795cdebdc2aeb65f5a8b74169f1408c07a5ebff",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/seo/seo_test.go",
          "sha256": "b934a766ae605d2d7c7ea725b79e88fadc646f3e4406f19751beab9619ee3e44",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8",
        "google-sitemap",
        "google-title",
        "google-snippet"
      ],
      "reviewed_behavior": "URLs por tenant, robots sin CR/LF/NUL, JSON-LD y política local 60/160.",
      "comparison_owners": [
        {
          "path": "implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md",
          "sha256": "9f9aebf972ed734459ded6de84dbf569f1e36d05957c166be3d46ace2c554646"
        },
        {
          "path": "implementation_packs/GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE.md",
          "sha256": "2d9f4742402dcbdcf6aecd49cb55eb0697998417bd3797cbf5a9c2d6f00221f5"
        }
      ],
      "comparison": "Metadata web no prueba sitemap/robots/JSON-LD. Este core retorna URLs, no sitemap XML; longitudes son locales, no límites Google.",
      "conditions": "Dominio autorizado, serialización/formato/publicación y tests de URLs/canonical/robots del frontend; no garantía de indexación.",
      "roadmap_owner": "T2804",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Referencia AUTHORED: lista URL por tenant exacto, robots sin CR/LF/NUL y JSON-LD con URL HTTP(S)/hostname; política local de meta60/160 code points. No sitemap XML, dominio autorizado ni certificación Google.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestConcurrentSitemapDedup",
        "TestJSONLD",
        "TestJSONLDEscapesScriptText",
        "TestLocalBusinessAbsoluteURL",
        "TestRobotsRejectLineInjection",
        "TestRobotsRender",
        "TestSitemapDedupAndOrder",
        "TestSitemapExactTenant",
        "TestSitemapInvalidURL",
        "TestSitemapRejectsHostlessURL",
        "TestSitemapReturnedSliceCannotMutateStore",
        "TestSourceMetaPolicyIsLocalV374",
        "TestTenantIsolation",
        "TestValidateMeta"
      ],
      "reference_test_count": 14,
      "seed_case_count": 16,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_SLO_CORE.md",
      "id": "GO-SLO-CORE",
      "version": "0.1.3",
      "sha256": "ca27fa75c2e20c2de18552ad5b5cac03f695efbd35c7fc2a77f9a4f9a21efe81",
      "files": [
        {
          "path": "internal/slo/slo.go",
          "sha256": "900b34d5ae9ef847176fda3b6f33b131df02cfe88002222e8cb54492ec520962",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/slo/slo_test.go",
          "sha256": "8a704820fa654b0844ef4746838e5e769a4c95c314e6d03961254e9602aac08a",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8",
        "google-sre"
      ],
      "reviewed_behavior": "Fórmulas de disponibilidad/presupuesto/burn de acumuladores y comparación booleana de dos budgets.",
      "comparison_owners": [
        {
          "path": "implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md",
          "sha256": "8b2b457d8c5232462b1c4cf3d5e4c7bdb021f752e107f0fc665229b254bd79ef"
        },
        {
          "path": "implementation_packs/WINDOWS_REFERENCE_TELEMETRY_RUNTIME.md",
          "sha256": "82aa19096bb60a60ac834a33185f2dea6d79637bfdf9d009b76427131ca22ecd"
        }
      ],
      "comparison": "Reglas/alertas V372 del runtime no importan este Budget. Fórmula coincide; no hay ventanas temporales, ingestion o paginación en el helper.",
      "conditions": "SLI/ventanas alineadas y acotadas, thresholds validados, capacidad int64, señales y respuesta probadas sobre el runtime real.",
      "roadmap_owner": "T2809",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Calculadora AUTHORED en memoria de disponibilidad/presupuesto/burn rate para target finito en (0,1), más comparación de dos budgets suministrados por el caller; sin ventanas temporales ni alertas desplegadas.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestAvailability",
        "TestBurnRate",
        "TestConcurrentBudgetRecords",
        "TestErrorBudgetRemaining",
        "TestExhausted",
        "TestFiniteTargetBoundariesAndEmptyBudget",
        "TestInvalidTarget",
        "TestMultiWindowAlert",
        "TestNaNTargetRejected",
        "TestSourceBurnRateGoldenV374"
      ],
      "reference_test_count": 10,
      "seed_case_count": 20,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_SOCIAL_POSTING_CORE.md",
      "id": "GO-SOCIAL-POSTING-CORE",
      "version": "0.2.1",
      "sha256": "2bb5438f96951fee8951882b7a7398ee8736fe912c0751ed30d73e440b7d2831",
      "files": [
        {
          "path": "internal/social/social.go",
          "sha256": "56bb88d4f04ff1e8e30bf6f9376603a02449f2cef7548c6e33721e43e5e56102",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/social/social_test.go",
          "sha256": "49a04c495857257f9485979e1e28de065d0229a81edebdef007dace20eabfafd",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Agenda scoped con estados locales scheduled/posted/failed.",
      "comparison_owners": [
        {
          "path": "implementation_packs/PYTHON_META_ADS_REPORTING_ADAPTER.md",
          "sha256": "2ed2981302f3c7c6caf69eaead0205cd9cf2e54b6799b719664241cb6339e53f"
        },
        {
          "path": "implementation_packs/PYTHON_TIKTOK_ADS_REPORTING_ADAPTER.md",
          "sha256": "8804c81cfe58e25a062830f5f75f8f664b6bb1b097565f5d3201cbf01056b787"
        }
      ],
      "comparison": "Reporting read-only no publica. MarkPosted/Failed sólo cambia memoria; no prueba proveedor, lease, cuotas o receipt.",
      "conditions": "Adapter oficial write, aprobación, cuota, fencing, receipt/reconciliación y revocación antes de cualquier efecto.",
      "roadmap_owner": "T2805",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Agenda AUTHORED en memoria por tenant/id, consulta Due y registro local scheduled→posted/failed; no conexión, entrega, retry ni publicación real a redes.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestConcurrentTenantPostIdentity",
        "TestDueHasNoTenantBoundary",
        "TestDuplicateRejected",
        "TestForeignPostTransitionsRejected",
        "TestMarkFailed",
        "TestMutationAndReadRequireTenantArgument",
        "TestPostDueCopyAndEqualityBoundary",
        "TestRejectsPastAndInvalid",
        "TestSamePostIDAcrossTenants",
        "TestScheduleDueAndMarkPosted",
        "TestTenantIsolation"
      ],
      "reference_test_count": 11,
      "seed_case_count": 16,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_SURVEYS_CORE.md",
      "id": "GO-SURVEYS-CORE",
      "version": "0.1.2",
      "sha256": "8a03eb00d9bdabae1772830a855ba1a330773c091ae24be5dab8843ea4ac926c",
      "files": [
        {
          "path": "internal/surveys/surveys.go",
          "sha256": "cc7cc6be6f8ee6ff330911311060500228e090def644ead9979a433fa380665d",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/surveys/surveys_test.go",
          "sha256": "c2b2efd795a156f22874e6d3dc37054a300517227ab76ceb2c9cb9a155d165e2",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8",
        "bain-nps"
      ],
      "reviewed_behavior": "Respuestas únicas por tupla, clases NPS 0..10 y fórmula contrastada con Bain.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_ENTERPRISE_QUERY_API.md",
          "sha256": "84601b5f6b3917f4f540ffb90d9d569b6cf3e9a8e0fa4955a736c2c4e9a38f70"
        }
      ],
      "comparison": "Reporting no captura respuestas. Fórmula no demuestra cuestionario, consentimiento, representatividad o invitación.",
      "conditions": "Autorización/consentimiento, encuesta y población pertinentes, privacidad/retención, captura durable y UX.",
      "roadmap_owner": "T2802/T2804",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Materializa encuestas/NPS: respuestas con score 0-10, una por cliente por encuesta, tenant-scoped y NPS estándar (promotores 9-10, detractores 0-6, pasivos 7-8).",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestConcurrentSurveyUniquenessAndDebugRead",
        "TestDistinctSurveyTuplesDoNotCollide",
        "TestDuplicateCustomerRejected",
        "TestInvalidScore",
        "TestNPSCalculation",
        "TestNPSEmptyIsZero",
        "TestSourceNPSGoldenPopulationV374",
        "TestSurveySnapshotAndScoreBoundaries",
        "TestTenantIsolation"
      ],
      "reference_test_count": 9,
      "seed_case_count": 12,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_WAITLIST_CORE.md",
      "id": "GO-WAITLIST-CORE",
      "version": "0.1.3",
      "sha256": "d80a1f623c6933fde10d4c3c1d0933072d3c767441db042f5df27b5f06ae6845",
      "files": [
        {
          "path": "internal/waitlist/waitlist.go",
          "sha256": "a83273d7fa047571ebc5953e1a97a97c30170973bafddda6f97024c0fb907794",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/waitlist/waitlist_test.go",
          "sha256": "89a167dd428fcec6e91beb640e1e5b96a8e8dd929333deee6be452ed61b54963",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Posición FIFO entre pending; estados locales, Seat permite pending o notified.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
          "sha256": "9f64509e4ca6d658c16800b6923cc55af947f4eb2bfe40c75a234c508b458fea"
        }
      ],
      "comparison": "Agenda/capacidad no equivale a espera. Notify no envía; no reserva cupo, oferta/expiración ni recupera después de reinicio.",
      "conditions": "Política de turnos aprobada, reserva/concurrencia durable, oferta/expiración y notificación con receipt.",
      "roadmap_owner": "T2802",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Materializa la lista de espera tenant-scoped: join con posición por orden de llegada, y máquina de estados pending→notified→seated/cancelled fail-closed.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestAliasedWaitlistCannotNotifySeatOrCancel",
        "TestAliasedWaitlistLookupDoesNotReadForeignState",
        "TestBadTransition",
        "TestConcurrentWaitlistJoinAndTerminalTransition",
        "TestDuplicateSubject",
        "TestInvalidEntry",
        "TestJoinAndPosition",
        "TestNotifySeatExcludesFromPosition",
        "TestPendingFIFOAndTerminalIdentityPreserved",
        "TestTenantIsolation"
      ],
      "reference_test_count": 10,
      "seed_case_count": 12,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    },
    {
      "path": "implementation_packs/GO_WARRANTY_CLAIMS_CORE.md",
      "id": "GO-WARRANTY-CLAIMS-CORE",
      "version": "0.1.2",
      "sha256": "acd508e61b9325564af1a1981f2a29a51bcec7f8400936a122f5b65a055d9696",
      "files": [
        {
          "path": "internal/warranty/warranty.go",
          "sha256": "d3f7799a357dc018014ea03bf48b8e9f2cc314f754fa325b9f74f23be7c34bbf",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        },
        {
          "path": "internal/warranty/warranty_test.go",
          "sha256": "a77d369a631557626ecfabf22f83bfa3ded9ebbce094faec5ce90dcdc2fabe20",
          "provenance": "AUTHORED",
          "license": "LicenseRef-Workspace-Owner"
        }
      ],
      "admission": "CONDITIONED",
      "source_refs": [
        "go-1.26.8"
      ],
      "reviewed_behavior": "Estados locales de reclamo y notas requeridas antes de rechazar/resolver.",
      "comparison_owners": [
        {
          "path": "implementation_packs/GO_FULFILLMENT_SERVICE_FRANCHISE_API.md",
          "sha256": "b647b19aa840b1f6723fa9e349425f8d042084579da6cda4cace91c28bb021a6"
        },
        {
          "path": "implementation_packs/GO_ENTERPRISE_QUERY_API.md",
          "sha256": "84601b5f6b3917f4f540ffb90d9d569b6cf3e9a8e0fa4955a736c2c4e9a38f70"
        },
        {
          "path": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
          "sha256": "9f64509e4ca6d658c16800b6923cc55af947f4eb2bfe40c75a234c508b458fea"
        }
      ],
      "comparison": "Casos/devoluciones durables no prueban términos de garantía. Core no decide elegibilidad ni ejecuta reembolso/reemplazo.",
      "conditions": "Términos/regulación y evidencia de compra, autorización, transición durable, efecto único y recovery en owner de servicio.",
      "roadmap_owner": "T2802/T2804",
      "business_equivalence": "NOT_DEMONSTRATED_DO_NOT_SUBSTITUTE",
      "candidate_promotion": false,
      "claim": "Materializa el ciclo de garantía/reclamos tenant-scoped: open→in_review→approved/rejected→resolved/closed, fail-closed y con nota obligatoria al rechazar/resolver.",
      "license_ref": "LICENSE.md",
      "license_sha256": "8c8228f1430a714aedd95ad5eaf5232017128238d7e560cfe42789fd0c4fa6fe",
      "authority": "SUPPORTED_REFERENCE",
      "reference_test_names": [
        "TestAliasedClaimCannotReadForeignState",
        "TestAliasedClaimCannotTransition",
        "TestConcurrentClaimLifecycle",
        "TestDuplicateAndTenantIsolation",
        "TestHappyPath",
        "TestRejectRequiresNote",
        "TestRequiredNotePrecedesLookupAndDoesNotTransition",
        "TestResolveRequiresNoteAndBadTransition"
      ],
      "reference_test_count": 8,
      "seed_case_count": 18,
      "source_is_enterprise_upstream": false,
      "admission_decision": "REMAIN_CONDITIONED_NARROW_REFERENCE",
      "gates": {
        "G0": "VERIFIED_AUTHORED_EXACT_HASH",
        "G1": "VERIFIED_WORKSPACE_OWNER_USE_ONLY",
        "G2": "REVIEWED_NARROW_BEHAVIOR_AND_NON_EQUIVALENCE",
        "G3": "RECONSTRUCTED_MODULE_BOUNDARIES",
        "G4": "EXECUTED_REFERENCE_SUITE_AND_FINITE_FUZZ",
        "G5": "CONDITIONED_CALLER_AUTHORIZATION_AND_PRIVATE_BOUNDED_INPUT",
        "G6": "CONDITIONED_NO_PRODUCTION_SLO_OR_DURABILITY_CLAIM",
        "G7": "CONDITIONED_TARGET_OPERATION_AND_RECOVERY",
        "G8": "NOT_PROMOTED_REFERENCE_ONLY"
      }
    }
  ]
}
```

## Reproducción y continuidad

Stage: `<LOCALAPPDATA>/Temp/elite-v374-8966287c156543bdb575918394dfd339`. Materializar cada pack a destinos ausentes. Reunir sus42rutas en un módulo aislado Go1.26.8 sólo para investigación del candidato. Para uso/composición, respetar los estados canónicos:20condicionados y1rechazado. El perfil sintético20no cambia los55planes admitidos ni el integral67/754.

| Recibo | SHA-256 |
|---|---|
| rebuild374.py | `f64812f21969fef745070e448dd9caad9345cb030d177c53c16ace6bff564340` |
| suites374.py | `aa73c2600452e5ff7b47ccd7ff6c285dacf2169b4b3551020afe8081c1b96fd2` |
| comparisons374.py | `a9ad59a1948c4ab57c6ecba9130b9499726c96d11c917e8862a6a0ed68ebf69d` |
| run-comparisons374.py | `6a54e3f1d6cb0c5a47552528a871ee953542c22c4e561b910c388bf3d05eb70a` |
| fix-comments374.py | `15bc9aae029f2f8d14fbe4ac38d6052814dca37be0242d1d5cd6570cca6c07c3` |
| integrate-cores374.py | `498c3b3b7bbc5cf6eb1d4264b85779628ac7791abf2fbb362079c83ce415990f` |
| qualify-final374.py | `ff2c9ad882e3e55b79c86545977aba605ee310a0a6c89579d272f5f8f50669aa` |
| decisions374.py | `0b12f5d64943116f45dbf09b1e1ed30bd673255394a9323588db7c7269413d70` |
| assemble-audit374.py | `a881f4d1796659d8bf6f24bc0efe1571a4f904bfc56000fd210031dc549abef0` |
| fuzz-profile.json | `33ccb65dfb8b00d3b061b2257f58870b910b504409096c2cdb33532b4be33de8` |
| fuzz23-receipt.json | `421acf69c6c98d6a6194585ad5e29f7318fcdef59a85b1db8f566ca42f030745` |
| fuzz23.log | `7b3dd5b7b8d332da615765e49baa0e0a0bef93d96c4654a890c9397bada87a8b` |
| final-suite21.log | `dbb078326356cc0f7e911d0305f91e5d0553c9ae132a0af5ded33a95b2d3d966` |
| final-vet21.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| final-build21.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| final-deps21.log | `5fd6ed946c7f6b54a3fbe55eed6ed0ea46e9fe64608c7afef38040913b5a055c` |
| final-materialized21.json | `820c4b7823466f8de466473e656c3bbee137998f82b9e711885430060f5f456e` |
| source-comparisons6.log | `8cf75dce211037698423f8b0d55cd604c8e35502a1f1fa193f4c46aac580aa1d` |
| reference20.md | `8f59e419d553ecbefa9db1dd8921e836884f1b5a5dd3d6c3180f0da7dfb63cad` |
| composed20-suite.log | `c6d26a53af4abcbceb891d60bc241a9e346f45be61e9e54f04515c049355afb2` |
| candidate-denial.md | `9b076831b9140d3aace7aefafa852708a95275fe074e6f04ac30d801acadc402` |
| candidate-denial.log | `5140ca816bfe82c48fc260e9a70cd5c380d36f812b8295b1e41a514623c7ef47` |
| canonical-changes.json | `9cc04751299981d087b5a9866708c18d41062ab2c81727a748080b3b22f7722d` |
| authority-receipt.json | `38fec5f9e3b50275b42f656be84e32102e9e910a768f8f6977b2aff9e3491320` |
| before-hash-review.json | `a5338d955b09046ec0b16f3a9625b7955c763aae07dc722e474e6078745f932f` |

## Recibos portables

Compresión reversible de scripts, salidas y planes sintéticos; no contiene datos de consumidores, claves ni snapshots HTML externos.

```json
{
  "rebuild374.py": {
    "sha256": "f64812f21969fef745070e448dd9caad9345cb030d177c53c16ace6bff564340",
    "gzip_base64": "H4sIAAAAAAAC/1VTTW/bOBC9C9j/QAQLkGxkqs6h6MbQYeF4T21qOA56iFOBFsc2NzJJkFTt9Nd3RnLR5iTOaD7evHmzi/7Igs6Hzm6ZPQYfM1uiWVze/yfvyghl6rch+hZSKn0qDzpRQvFQU6xomp3toGmkipB89x2EVEFHcHm2GiJUezJCFsuvNZ/fVo8JYqruP1Wq1e0BqtYbOE9i77I9QrqYIdqjjq+/3JWBAM6Aay2GOJ3td6iCP2GlA3RdFU7poOAMvIj+lGqCrTqvTRLioeKtjzChKZOiP5yAatNkOGch5cz3uaaoPhLmmyknjzq+GBsR9c5H5vQRmHWMit8WLNSrilwzqlCHP6vNgjV1BJVAx/YgIv8WdPvSWHPLrsTTt6vna3nFS4pFWtVn+TR9LtjJ5gNDoAKTr/kkwra3ncEB9lxK5XFywc9bLplObIf9Wax/L0QhReJp+bXkk3u/jJ52wdH4b/imHMWq4kedIVrd2R/QIK8vxp9cQ8hUSFMuMXyJxu+UQK47SNkS10jZ4L1QmXiFQOUz+gxxt6MHxDegHtZ3Xx7XJe2OQv55LwtEj5vPLCJfuY+OFl3X70usVTBiOUI3kAxqZ53RXUf0zVeLf9cLJtS1/PtP4oiGFNv6DaYKS8xMyjX2HN6Xls5nhm5UiE054ZbIGCV62fNopHodeyiHsMa/DNYYfIo2Q7N9zZAEth1XPppSFoL68b1XR29QXWPwIAiOnr4DBh26cKOt7qpBjtocbUpI7UT3xuaN27i9Z1N180F93KBEC7wAh/nzx9Vqcb++mTarxfzL/cN69ThfL+74KMxAhCW8VDBiBIFJEJ3uEMa+81vB31Xv1N6jkJAxu2NBkXIVHlMi2QmOOHE+irhtPR6b6wEVPjbfuLquGb8miXfD0TXZUx+pdGqCT/ZM4w/b66wbbuTNOagUcHD6hTzRyhABWSplHfOIQPBd71qGms2vAei7yYsYCfEIgxLkX8VPNrBFm6wEAAA="
  },
  "suites374.py": {
    "sha256": "aa73c2600452e5ff7b47ccd7ff6c285dacf2169b4b3551020afe8081c1b96fd2",
    "gzip_base64": "H4sIAAAAAAAC/41U33OiSBB+p+r+B94GNjgIukk0xYNrjJs6Ey3Fy95m96gBBmWDM+zMoLlc7f9+PciquVzV3QMw3fSPb7r760zwjVkStS7y2Mw3JRfKnIFoNGcuHVnFpeAJldL5JjlzBHVUvqHOmkjtZSwC7WBFUZYXNIpsLKjkxZZaNi6JoExdzWsLnOxSy74Kg4WLkkroP76HjNlDgIZ9dympkO79xMUJSdbUTXhKn1uiYjqXbMRS5Bsi/vypdlNaUpZSluRgwojKt9Qt+Q4irWlRuOVOrjF9psgYT18nGZTlNVHEnfCEFG5IN6VLi1zR1ta/9FqddpK1e+8v2pe9tNvrpfFFnPViP/XOz3vxRRq7DfyW4rxI1iRn7oq7cf3Z56NsG6R5oiwuMZxzAXUbT8PpdDL8OLi9D1ChEyPQPUznvwaIZ5kWZvPpp98P0mJ5d/2hkYbjaTS6H3yYjK4D1NZ/7wafwHy4CFBXizeTwRjOLSg5KQpaBLV2OBh+HP3fq3e6nZZHuzSmpO13O91u3PUvaZpc+F6SXnj++3OPwA1bdYNagiZ8SwVNkW2kNDOhJxZziFhJBxodzOsh4ZUKOu223TdMqYhQgVbiDWdccZYnlm2Yu1ytTWvhWuwM4YKvkG1jDm210HOMbJNIM+t/D45DiHWiQxp4HF1seBypUp0v0wcqxKnPIryeLsMDpOYL2WGimALc32FqVSWYnjNH8Iql1j+htuobOL7tZEUl10EoKmpfEQmFVeapfxC0HWZomCirXl5aG6KoyEmRv1DkPM4eHNS65zPBNV8QCDf1VyphzV10YhzBrD+lfMeikiRPuJQessF8BsIrF6BqQTcwj0AAvjeWLszL/SC8/W0U3Sw/f47Gg3CEN2kd4JpKlbPaeB8DCFkDfcoVsr/ae+iy0kNxJKrzOJ46SIEvYMYuxlhjT6BUKvD0US+HvQrGotxfL0BnTYJaS1bU9zBUHtI4YZNoS9WbNKA7ZDlaxlVepG9sa+2/WJ8uh4Nxkdf4W/BTnoA+8YUmr6iSweNXI+PCLM2cmRJWIU0tK4RiM+gPA+7aeFXw2ELv3HeRrgpecRhdGHTtxciGakdBcZazFBhpCfRHVrHEtG6g1F92Z/YXCzkljA1Jwf9ZWTZsVnxn9xsAmJQav/UX0g2FyqE+wERnZbNVwbOoN16kuBXamMio5DJ/hjhIp0d9/f5hG82EFkCpJrQdBH7HOLS96RWua2HjnYDG7yFpDU6rTSkBhgTWb0gEbZR6cvqegygjcQEboK+54NTRIk0bgOq1dX2bhOjnpX44uW6JCnzbPjLE7/wHL04GFBZQVOdZAVFcCBGdKo4kEfwbTdScc7UPETba7E3Q19cHqzlNaF6qV1Z+Ry88rT7ajfnomSaV0jXQy/brYb/0YOUZ+9WCBpOJ70XD5Xw+ug+jxfI2HC32nJwNFgt0ukx+Mf4GvOhkvY8HAAA="
  },
  "comparisons374.py": {
    "sha256": "a9ad59a1948c4ab57c6ecba9130b9499726c96d11c917e8862a6a0ed68ebf69d",
    "gzip_base64": "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"
  },
  "run-comparisons374.py": {
    "sha256": "6a54e3f1d6cb0c5a47552528a871ee953542c22c4e561b910c388bf3d05eb70a",
    "gzip_base64": "H4sIAAAAAAAC/41RXW+bQBB8t9T/wNtBBBxg4o9YPDh25ER1TBQTNVXlWnAs9lX4ltwddn5+jzRNIlWt+nDand09zcxuJfFgNbne17yw+KFBqa07A3uvuWqLRiIDpVxUvXXS9ezttuI1bLeOL0FhfQTb8ZtcgtCTLFlTwvBgIFcoopBMFr8+kdkFfVAgFV0t6bRp5rnO6RJZXtMMDg2FmmvwjtEo9PoBq4Lx+TAYjct4PC6LYVGNi6gMB4NxMSwLylrZsXkasWb7nAu6Q1pwQZweiGNScqZtVL7JuUThLtIsTZez6+nNKiF1R0lM7Ut6/zkhWFUduLtPH7++ofXD7fzyFc0W6fZqNb1cXs0TEnTd2+mjGZ+tExJ3cDadXV8l/2mvH/e9EGIoIA+iuB/HRRyNoGTDKGTlMIzOB2FuzHgsZ3vwJDA8goTSGHtK3m/hy1bY35SW9oKSHVYH7cMzEMcl3om4Zy+dxrEqlFZjcWHZGSVcaJDCWHf8XY2FTc7o2VaD0v4OibPZuN3mzHM1PwC2OukHziRXxpK2BGrryRxbt1IwLKF34npv2ebUClvJwHu/uBr4Ne4MCTYgbPJcEMfKlVVd/EP/b/GdGuISn/q+b6Jnpkz4npny+oWnKzJshU7CLv1h+MjGZacyyd7kK1126qsuASk/sq6zefqQvRkMR4HTa6RZjP2nD+J+NPzXTXzq/QQ1qRYqQgMAAA=="
  },
  "fix-comments374.py": {
    "sha256": "15bc9aae029f2f8d14fbe4ac38d6052814dca37be0242d1d5cd6570cca6c07c3",
    "gzip_base64": "H4sIAAAAAAAC/51UTY/bNhC9L9D/MDfZgCSjORVZ+LDZGBsju97FWkVyCCDQ4kgiQnFYkrLr/voOJTlR3PSSgyFTGL2Z9zGsHXVgRWi1OoDqLLkAL3y82a/jY1GWtdJYlsvcoSd9xMUyt8KhCbfFer9KKur4qDyZN78nN3ZdrBJlAjoj9Mprir+8oeTWry0jCFkG/Dsslnz2fLZaVLhIVivuWX0VDQLXg3V0VBI9CAP7x+cVOkcuO/SywQAVOXwL4iiUFgelVTin4LATyijTfDEMNZTDVC6M5L/OgBMBU7C69/BA1GiE/esGul4HlZ2UkXRKp1Msz2L5gCY0usDQOeyDZJEyMvqcw92fxYfn1837PEn/d/xLDSiTddiRO/8w9zDbnNvYz6nQdhhUlcO7yyRQk+PhBEviEEbNUcKJS2dsbiG0yrObwywDWsufGIKgOpxYglY1VudKsxr7xy3P1qAPigywaJIdoTMjf2OdLG+urXqKKn0awO5iWYyNxo4j8YvSvoXgeoQoLJxaNHAg5hValrONeYw6aTINjGB+NLQWnnPGZcIMaFyvHEfBW6yCOiK/4ENLWnpYiCMpGaVpuB93GpEzzXUSvFVf0S8nK/9Db5LbQzgRVELzy8z31mrF31Z99GXoN5rIFjUcRh9mQ32fJIeCaY0gQCfjJ05gyQ44ZFLWRTUm6jkyv/h1OyGyxe+INPJ2tKgtA0nCaPPMCRBzKy5CT/7GgmiszU8cNhxX0qdoKpJctU76UGd/JKnBk1YG18kXE6t/2G0pfHsg4aSf/f3Zpv90y79/Mlt2CGiECZmvyLKyH1+24I2wLFyI+YfaCVMx+zE8DEtTejzz8ynnV6LjJzMUmtfLEk/r4+1wVMipEUd03DyF3cs+Kjtm0A5lQ3pZ+aZx2HA+JXfjizFmUFK8XMAHihlY2N4xRM1rnFWaPMrl1RUw43boVUyfgPjRFaEIfx2mo9A9+nwY7DonnnpXIfzVo1OR7SgWiD5wktU/Q3SGvHAgFPEGh+DUoY+vZ8GxlyEviYlOXVFkqVkpVZ+ZShCxO8vaKvsLkXEs7SK5f3562uyKfXn3uH3Ybd6Xn7bFh3LzebsvtruH8v55V7ze3RfJ8rebfwF5bddmjwYAAA=="
  },
  "integrate-cores374.py": {
    "sha256": "498c3b3b7bbc5cf6eb1d4264b85779628ac7791abf2fbb362079c83ce415990f",
    "gzip_base64": "H4sIAAAAAAAC/3VWzXLjNgy+e6bvwGY7Q2ksS3Ha7o89Ongc7zYz+Rvb2T3EWZWW4JiJRCokFdvt9GH6AHvqI+TFCkhKNpttTxZBAAQ+fAC8MrpgpXDrXC6ZLEptHDvHY6f9vrFaBQYCWy1Lo1OwNlgLS9qBXVdO5p1pTPphusk8fzirD16SrGQOSeKHBqzO78Hzw1IYUK5z/inm40F0YcHY6PQ4ClORriFKdQbbnqmUkwXY9lgaWQizexRHGZSgMlCpRBUlnLyHqNQb9LSGPI/KjV2HsAXeMXpjYwo9zLXIrOfNIp5qAz3K1IZ0wyk2kSUOts7z/WG6FuoabHx51Vlpw5QogEnFyNOgw8p4GpFo6OLyud1wCagNseuwewxDahUbCC0Ik649wz+XIr1N2psB2/MuP+9ddf09HjgENTzxL/tXQwTG2bhVCm2ZS+fxkPvNxWWvfxVbZzypnPck8bt9f6hgE6NieKOlaq78DqNUm6B69LjFPIvbTBoPttK6RN/Gc1OBP2yqF6a63FGtvDJ4aRiVIaWMPknBUmIrqTKR55TZeDoZzSfMC7v+T0/ptChmhCIjGA3khGLtAGFk1qRxXYsC45VYhoM+j1BpqPOsvqgMsaSVooFckU0D+XLnwHr+jzEqfyMhz2wj3Zqy93hVZsJBj3drLqIfv06ky5EN19z3Q4008vh2yX0mLFsN7uKv9A6RbN7l+aeA9071udEUOsfD+/qXCjGNOLLyNtMbldiddVBEzYtJXe0VdlTiDEBY2j730fQcxV/NSxLNdGVSmGrtGuFLTPz2RcsR2PwKdTJduXhFH2DM83hn88Ozi3lADUIq7/Z9QkNYbDDHlHbsDrFylVHUU+TtsUihKKmdaoA67CWxXezwXOYiBY+/oDHvtt9djlT+/haJWd/0a7/P/FhHzBqwy70PmvXDg9fhm70r9PC9/C3JG3tqqGpJnMM35EpClgiH79S9RO9/Kz7YP3jd23/X6+/XbZZqHB5xP1jl4rrm8Ak57cac84V69YqNqkw6bR6+CIYDj6W5kAX7+PObXxZqoY4FMvheWvnwj2Jo0WTGObOAF6lWWLtqB2zHSmyqpWCg2AfdJBB9wmbB6cFQja0qZDXYhUI+pVJgYbEnbkSG9NuxpvIiRfhQXkehlTPC6ZB9BLNQvz+9lTpUSjAkmoIQjc+mk2R8PDo6SUaHJ0ez2dHZaULBh0X2O8vAZsByzEE8fPljwMjIyUxkwULlMqU5KliuU5EHj/HT41LdIw8Fxhu0cWO6xcPfaJsKDDeTqcN2UuifwV0l7zEdcrVQiEqBdhIU3lqgpDJtQzbJGapcV+IGIiQwrRqlGWHQuEdHlWIKrjViwKAocWtQCHm4UKeapaJYYqSArdbUJ2Cj8yOMGUyB/ULmBKG2CC0VSklMusYdQ21AzWiYO7gmlFEl/I/SPlEay7tjtrKsni7sroJMKEYxgbmnhKjKONdgW0Im6wQILCO1Qb9ojusi3BiJA6HuJRcA9Z5U1zGv3Kr3lgdIo1wqiJGCHOnY7p7HhvyT05riAxpaQTuTH/uLD9qPgIsVPvpVjj6flO1aHPz6mg/aVR02Z6+5DetwcCX74Rq2mcSXseEf/f2PafnNwH1h2Q6UeulbPmiPf/mdeq4JpRUyJ+89ptmu32cQ1as6q4rSeq1SIGnRu/jA7zYg4X8BXIB8PDo9Oz0aj44P+sno4vBonpxM5qPD0XyEg7QqPOSZt718EdGVXy+jLa2i1j9GPTu7mI4nydk0mU9m8+T90fFklox/G51+mBxy/4fOv0OHsQUbCQAA"
  },
  "qualify-final374.py": {
    "sha256": "ff2c9ad882e3e55b79c86545977aba605ee310a0a6c89579d272f5f8f50669aa",
    "gzip_base64": "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"
  },
  "decisions374.py": {
    "sha256": "0b12f5d64943116f45dbf09b1e1ed30bd673255394a9323588db7c7269413d70",
    "gzip_base64": "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"
  },
  "assemble-audit374.py": {
    "sha256": "a881f4d1796659d8bf6f24bc0efe1571a4f904bfc56000fd210031dc549abef0",
    "gzip_base64": "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"
  },
  "fuzz-profile.json": {
    "sha256": "33ccb65dfb8b00d3b061b2257f58870b910b504409096c2cdb33532b4be33de8",
    "gzip_base64": "H4sIAAAAAAAC/7VVTW/bMAy9F+h/CHIu1uU27JZmC1YswYJkH4ehKBiJsYnIokfJSd2i/320nRXYQcHsoSdDpPye+PX4dHkxGo2DybGA+wNKIPbj96PJVWtHD1uHVg1RKuxsu+rx8T5SgWodT96GcWeOIBnGoMafzXk0euo+6irB7CFr77+5Jh9RPLhrCyHfMog9IbRXPXS4cyWZgUNvQabefl7dLtmiG3cXn6/+jWL3kICek6eIM/Z/Iu6Hm9EumjMvX6PFoowKvIkQcQkmJ489SXJ0pcHmlGCZSiTj8BOFyFIPyQ9N3vkE+A15kHrlKgE3B+e2CtAT3XENLtYJggXaDOU/8lOA7DGSz5LtU5RAmZ+aveeja+gKTeeQPPE2oBxgS46S8azYkakXnGlUM8cB7XdwFYZXoJpx1VxfsufInswrULRNlQkU0ywTzLRIfUl8O9tnypOj2TulWYkSYQhDClNCLexcqiSd9+MDmPgBbWWaiexbkJJTU75RfWqxp0IxL7QX+1aiFC64e1QiApRGAvRPrXnZW6cEdyg6wGmd6vwb/FWhNzikBIIFeasymuD4qjvEx/Xp1oJ2aGpVrd4sB8JjOo7G27xeoEnnkDgC8tkIvq0XQTfRmrccF6pWffsoOD67jG4qFahB6hTYEKRmYGO4RLvR/W6rl13Rl6CSA9bJMWi9awzan2FQCx2BYqMECYIfJ/f8dv5lGLz2uE+rqQMqXvryL4Lmc3d58fwbnPt5+yUJAAA="
  },
  "fuzz23-receipt.json": {
    "sha256": "421acf69c6c98d6a6194585ad5e29f7318fcdef59a85b1db8f566ca42f030745",
    "gzip_base64": "H4sIAAAAAAAC/71W204bMRB9R+Ifojxz2WsufQspaVFDiRKgElWFxvbsxorXTm1vQkD8e72bgGiljXhweVp5PDvHZy7Hfjo8aLXahs6xgPsVasOVbH9qhUe1PQeLbtX+cnX8fXB9cXt+PLq5uzv+Mrg+b+881Ju/3Kq1W7VyFZ5EnZNea80lU2tzCgXrJLu/svLx8d7yog4eBmZnnoOZo3G2p2rpDEutMi5qrzimlHRSlpEeCQIWk6ATkihKu1na63UD0g8DkgZJEvSDfodGlJE4TuOIJATjmGFvC7E9caFYFTJIIqD9IAvCiGKSpgnpYMZCjPokhS4kkBLoYTfENEydJ6FRlsUIYTeJk34vC96GNGXhQspSiMr2vOVjQedoK0I/t547XhUzoAvIa2Ynp1xa1BLEKXMJIAo0My+xnauEbZ5GLmlDECgZ6IFk3yYXl4qheONpLNiygmtPBrNZe2t/PnofdvbQgDniklscKvlSaE+AOc8s3cN1igyLpXWIMweDl0DnXKIv9DmKJcVq1QA/0JZTgV+5sUpvvKaahz3ZgHrGJejNRJQaxAiEIC6AL1ihNiDspgF5jCxH/T9SXYBeoOUyb2zqYgk8lwO6kGotqnMUrjJeU66IQb0CwgVvzMBECU43Y5W7PAyFMshuQZRoPvIMQ1VW7pdKKqskpx+JXbd6rqEY5LnGWvt9octa1fa0wBzpQjj8iXYnQGO8Fn8JG62EaCr7dvf8Aaj9jKykleR4K/pSNenbzGl5DTrQ3M4LNyHequ3uzUL9S+Mvzqgr8XN/uoZb+tN0jRlqJ13Nmr7dn+HvEiVFr2XWWLi3hrukGsCvUYK0053XmGdIN07h/cGvOK6bmVe7FV8NVWW8Mjeo9nK+mY6NezJMFVF27JTdW3cbofa+Gs5KJ+Z+ldwoyqFplmdULZHN3HOWla93tzfkUq9w0zjO9e4UjRsn47ex18BtJY4NyD9226OL0ZVnXDersvmuEsCL1zF6H3L1+XV48HzwB+/MSxp5DAAA"
  },
  "fuzz23.log": {
    "sha256": "7b3dd5b7b8d332da615765e49baa0e0a0bef93d96c4654a890c9397bada87a8b",
    "gzip_base64": "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"
  },
  "final-suite21.log": {
    "sha256": "dbb078326356cc0f7e911d0305f91e5d0553c9ae132a0af5ded33a95b2d3d966",
    "gzip_base64": "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"
  },
  "final-vet21.log": {
    "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "gzip_base64": "H4sIAAAAAAAC/wMAAAAAAAAAAAA="
  },
  "final-build21.log": {
    "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "gzip_base64": "H4sIAAAAAAAC/wMAAAAAAAAAAAA="
  },
  "final-deps21.log": {
    "sha256": "5fd6ed946c7f6b54a3fbe55eed6ed0ea46e9fe64608c7afef38040913b5a055c",
    "gzip_base64": "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"
  },
  "final-materialized21.json": {
    "sha256": "820c4b7823466f8de466473e656c3bbee137998f82b9e711885430060f5f456e",
    "gzip_base64": "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"
  },
  "source-comparisons6.log": {
    "sha256": "8cf75dce211037698423f8b0d55cd604c8e35502a1f1fa193f4c46aac580aa1d",
    "gzip_base64": "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"
  },
  "reference20.md": {
    "sha256": "8f59e419d553ecbefa9db1dd8921e836884f1b5a5dd3d6c3180f0da7dfb63cad",
    "gzip_base64": "H4sIAAAAAAAC/9WYW2/aMBSA35H4DxZ7m5qsdC/T3rJgaLQQIyeDoqqCkBxWr7mwxLRFiP8+OwmXvczZG5YiJfb5fPl07DjwAfm7jD8DZxEqYA0FZBEYeZbsUJSnm7xknOVZt9PtWEmCSp4XUAowDVmGxLXKt1kMMdoU7DXkgNbsnW8FYiIvR5BxKESkBBSFm3DFEsZ3NwjeRXUWJgjWa4h4ifJCtM/jbSSHQmGcsrKUTwX83kLJITbl+Mvl8lcpp7LvdhDqXcxuCoXke19Rr2/e9m6qeAlRAXwaFixcJVCK4P5QRzZh9CLLj7KE0L6+VQH+LDth6SaBVMw+lL0vKv7TiCwGln//jVh04C9sQrGZxvVYTePoxYll8xExzqQhyQvs9TzXW7Nv3l2ERAdZ/pZA/BPsPIsrNTlPXmzhTK1ZbfN4rBF1H3vHwtPFSH+p19WHm/91Hj6oXIcPCsfP1+44cobBwm6TV0kadou8Xr3zPXYnCxt7AaYqaYkaNar7anb6XzyVrmR093TJ3HKDuUq1wXS3HVv0Ow4cb6TyPYH/NL4z+9duTLzqdGmhfCZ1z/LEmlPiuirhBtP9/TwhysNIILqv4wklYxI4xFPLnkjdM0vxEFNquUrlE6j7zqV47HgDTFsYN6Du65riqYNnLXwrTPf8+pioTAWivaWrtnSJ7m8nn9iO5crjp80HVU0bDa37rvV/0CmeK3dtg+m+nmeWE7iOH6h0j5zuK3tmicPUkz+KXMsZ+2rtGjdq/MqyLW9P3c6h+lus2/kDDvLQG8YTAAA="
  },
  "composed20-suite.log": {
    "sha256": "c6d26a53af4abcbceb891d60bc241a9e346f45be61e9e54f04515c049355afb2",
    "gzip_base64": "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"
  },
  "candidate-denial.md": {
    "sha256": "9b076831b9140d3aace7aefafa852708a95275fe074e6f04ac30d801acadc402",
    "gzip_base64": "H4sIAAAAAAAC/1VQXWvCMBR9D+Q/hOxtzM7h295mESkMCm4URKSNyVWjbVKSbD4U//vyMbE+Xc5Hzr05TyRnSkjBHBADJ+BOakV6o/Ueo6ZpTlYrjAaMCKFcd722MjgqMNYP+k7oWzalL1G3wA24ihnJdi1YLw7XpPSMnwPeBETIkEYU3DGEyK5voQPlWEivo/91Wdbl/Guxqj7mxWfxva7zcrXIOpHW/b/n50KEhGU5eTBPgnnk/L1fPM1m2Wwk+QylLy2IA+TalxFOCNc68wN3116mP21ujOee6Q1sR5seCkh0HFuMrrFUjP4Amkce+ncBAAA="
  },
  "candidate-denial.log": {
    "sha256": "5140ca816bfe82c48fc260e9a70cd5c380d36f812b8295b1e41a514623c7ef47",
    "gzip_base64": "H4sIAAAAAAAC/6WPz0rDQBCH74G8wyA9eFkhzZ8muaVJkEBooEZB2Mt2d8W1SWbZXayHIj6ND+aTGGv1AfS7DDMw3/ymfuFSO4VTDmVOb600lm5aWmhdMcdoi5wNtJejpnJQTpLncBWRNEuSZbriQZzEUbgTu3gVZ0EaZpF4EGGYUY6jRqscGuoQB3seSDIysxd4mIg2+CS5u9I2yJdB5HutmiQcfQ/mDo5wwj0aPMCFZnwPTIzK2jkoCJQWJnSgpRnVXL5V8HP09MzicjGvwcfb+6z84qx8/S+/ur+kuu5It76pt3fFummb/p6U3bYGZaEsNlVTFX3te5/v1bV0kgEAAA=="
  },
  "canonical-changes.json": {
    "sha256": "9cc04751299981d087b5a9866708c18d41062ab2c81727a748080b3b22f7722d",
    "gzip_base64": "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"
  },
  "authority-receipt.json": {
    "sha256": "38fec5f9e3b50275b42f656be84e32102e9e910a768f8f6977b2aff9e3491320",
    "gzip_base64": "H4sIAAAAAAAC/83Xy27cNhQG4L0Bv8PA20YS7xd320foriiCQ55DW7BGHJAaO4Oi717KcVygE8GN6ybeGBJF/dTo8yHF3y4vdrs/1j+73dWIV9e7qwDj3M2HevXhqflYprX9dlkO9XoYHh4e+pmWQ8n7vFCpp7rQvo95P0DIx2XYE9RjGeeb7pSPpWtduy99uxpzoeE5OYdK5Z7w43cYApZ1BMGE6ZjvOPuV62uhrrnohbPa858Yu2bs+b56C0Kb9RbvbaAgIxmPDoQSUjsOgJSsNCwqF2TQiC45mThnaIFCMOBMasco3HNkOC1UW6KxVpgvjdMYaa70sT32gdbhfsnxuKd5gWXM86795gOUsbbDPE+nn3dz3i30aRliRhoQFtgVwrEuZQzHhXCXS2t4CsX+ah3lzw/nyDc530zU1UJbzO1S/7nX8JDLXcj5boCJyrK+9jx3dcr1ZclXpmxjec+sZZtYGgT3GsgLx9AL5TxFE3Ugbx2nxD1HK7m0misEGy03FCAkrayPKSb9FSwt1Q/GihOWLab1Wn+cx3WIPpebYZyRPj02d/VAcThMxwJTV44T1Re5/mPaFpvsOVfMb7OR1YwSQtApQQSrTSsRnxJaYV2MyqpWVjwqpnQAdBh8JO3QkTPeKmHP2IRlnL2TGhvb5AWHLUCke5rac5X6VCiP81wlKPF2wBzrEAs8TGu9PGK0g+Epsg7hOE54NsKW7v851Da9N1w5uUmfpFApGcZDcAZAGZW85MZEZkg4hSlwFRwqi5g0tPk0GCFCMk5zSdHTGT33zHP1PuyXcZno9fJwOLQzmCMNj0lds7l7A+Z/mbtlqnvprJdq09TxAIJUiD4qo1UIbX6NgNIaG2xop1Kmxk2GYWAWFRFL5NZVliEkUuemzhn2TkzrPLbXt7yJ6j+z3oR0M3TL0/TMSsn0piegZxC8TELaBhcVEjJmHQlF0fGIPnglY0jeOOOgVa1KSqNslkAov+bJnbQ/2LN9Lk7denu3b0HTFmh7prn1pD0t5dSP+fP7XlfDOqwRw3phbCd/R738XfQ2oVugtpdKtJrZBOUEIkiUToHhwFQrP5FatXotQ6OJQqnmShpS4hHauktIXEqnWZuekeAc1Bij1HsAra3TuIz39I2eFNuOYjkNtzDj4xL4HPT4P/Jq0m/P3VZtWxHr7aaqABd4cBG8822b0kozBad4SCkFaLsYE9GaaG2SCO0bWbXlNgnyxrvgkrbhTFVxKaz4nqqXF79fXvwFADkUzZUOAAA="
  },
  "before-hash-review.json": {
    "sha256": "a5338d955b09046ec0b16f3a9625b7955c763aae07dc722e474e6078745f932f",
    "gzip_base64": "H4sIAAAAAAAC/4uO5eUCAA99VwwEAAAA"
  }
}
```

## Prevención de deriva y correcciones

VERIFY_LIBRARY comprueba ahora el censo exacto de21, SHA del pack y sus42bloques, claim/versión/admisión, licencia, fuentes, condiciones y nombres de tests presentes. Su función canónica conserva13negativos: omisión, duplicado, alteración de hash/claim, promoción falsa, atribución upstream falsa, condición/prueba/fuente/licencia ausentes, escape y hashes intercambiados. Este gate enlaza evidencia; no ejecuta las suites Go ni sustituye freshness externa o aprobación de producto.

FAIL706/708 conservan las rutas/alias supuestos y su resolución desde owners observados. FAIL707 corrige sólo comentarios de implementación, con reconstrucción actual. FAIL709 conserva dos errores de transporte anteriores a ejecución; los scripts finales construyen fences desde carácter96. Ningún fallo se convierte en PASS por omisión.

```json
{
  "bind-audit374.py": {
    "sha256": "c2baac932d5f095978365391065b16598a550c9c41f9449d60ecde83f4113f58",
    "gzip_base64": "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"
  },
  "bind-file-hash374.py": {
    "sha256": "9d920cf3411b691443ae00b02fa4afae5282ca3fc0d655add1df7d9c8f76d540",
    "gzip_base64": "H4sIAAAAAAAC/6VSXUvDMBR9L/gfrt0giW6VCYpWypijw8IEKSIMHSOmd2uwS0qSsf180078fFAwLyGXc+895+QsjV5DzV1ZyWeQ61obB3f+GeRJc0ViW1B2VSf5CXlI82wyW0yz63yUz6LaDsiVTerIIC8WDneOskBXRRKGIbRnZC0a1783GwTabRDRWCvHpbKUjPN0dJ8CgWPoLmWFUUOCQZ+rAr5jbclPz85jCD/Q+5J/kpAwBkRog2D1xgiEktsSlIZKqxUaWHMnSrQguNJKCl7Bc6XFiyWeaKBw+0G4HX3j5Ui1ggRIp9OBSTZN45bmoyi5mV+ef2H8qR68DWmne/McGtVMoUPL2n6DK9zN4zi1gtdIP29jjRI6jOnw8H0pi9jR8BfpwT+dfueU2dvGpn1z76uIP/or3Q+PW4t5Sw58NEAqsD4z1memrrhA6os9/wW9gc9YtDXS4T5JtodK6MaZhGzcsn9BGlglFSbkSREWHASvvHMpc7oCAAA="
  },
  "permanent-regression374.py": {
    "sha256": "9200bb4034cacb6a96911fdaf9c4f32c098b8fa9fa1be13b473e3b19ce5186d4",
    "gzip_base64": "H4sIAAAAAAAC/6VVXW/bNhR9F7D/QAgGKMGSGrcDtsUQFtdxCxWtE8hut852BUa6dthKpEpSsbc0/32Xsh27XYIWmB9skvfr3HMu6aWSFamZuS75FeFVLZUhl7h10tj+RPm68Pz+pN14WbbkJWSZHynQsrwBz49qpkCYfh2nT+i7UZq8eJ+9Tp6ng/R9VOse7eu4Rm9WZAY2xvMdAVCUELsEPwOtQZlwqhog3pnXUXId2Qraj4ayEYaE8Jk89QnNpQKiZaNyIBUTfAnakIIvl6A0dZ38GvJPOnbdR7MqQFcQOSAMbTLBKtDkC5lgu+HF1UfIsdZbwT83cFz6odDcWhFS0dQlz5kB0oKzNoRW11ysqGNhLPGY5dfE61jbGCsSLsijaHxy24Z9g/8QHOYVM5iPfpji0WwQ/sXCf07C37JFt0NJyERBZgpWsFmcnib6jfW10RsTEOr9Xvkflo3ICSXdg9tI56yGQw0fjXQ+96iPDbKm4AaKbWdcEyGNbSBnQgpsvNzpsW32zkHyHdYiJ1uJrbNG+TW2W5csB297HmzV6u52Pd+pdexNnlBbKGyrhldcFEjks19+bofI/2qC0L0+JHXPbi1/MV0rKVYENQGhgbrBvYGLG5zQXSe0XzUGVYtvcXBZ5XWYTzosshrq2cniQWnwPKaW9GSXaQw3oEYbyBvc0Lu7YC7IYzB63wKmAlbM8BvQce8pDY63z6jf/x4Ta4WqbKmodYBIpbXHtDHL8FcaCFiXXCCOuaC+cyWLv7H2DMtjHth4tJMzjaXOUOPT+1P3D5s0vGhM3RhChxfpKBu8PU+m2fNkfJ6MX2aXg8nE9RfbhPbrwD/5Ktijoz8vR8Pp6DxLR69wkVyMCe22ZSPLkD8XbkARmx1Hw6WI3bnYr/ezP0Q5hiXj1cCSkMIK1dFo92ayvawL0mkteGmQetKRiq+4wJmMyVCiSMpMZfhKY8IwEQhrd8W3QSQ8h9pck97JCWLp2m66iOEON/sJpj8KiO6vhXaOJ/0Qn+Lghpcs//QGDCuYYd5MG4WKYQv2TcUOaLD37tIjKn4oEqdLloV9TR8Eiu+Hnex23b4I0UslGxwI+hHJoYvoHSvxlfmyZ+0F/hvseLvnyHfxzV7bEodcyPP/T2yVexj1IfejPgcJjrzv9UNO/vv84KG9HshZvz6+Rt+9RT85/wL8ndXsJgcAAA=="
  },
  "test-audit-binding374.ps1": {
    "sha256": "cf72b18468b513af3e4a68fd34ff5eac14011f01d786d9da1fc46c07d72b735a",
    "gzip_base64": "H4sIAAAAAAAC/51VbW/iOBD+nl/hD9E6uW1CU/q+ilQK6YoeSyuge7eqUOQmpnib2Fnb6Ys4/vuNk8ISqVeq+4SJZ55n5pnxjB1JKWQn0Uzwa0lnVFKe0BCPtSiwNabaG2vJEv1NpBR536lUYIgGRFOlLTtjd5LIl5EQOnS+gvFAJMRAuf410XPLVprcA1r3tHWjwLc1HLQ6RdEjmrSMadaa0Lxo0Yxp6j22j/a945PDw73joyQ4ODzYb9+ldwdHByfBcftkP52l7fYJtmyidHj7jXBAzinXfqfUIq9Y/QHh9yV8B3YJfNPT0+pwwTLqOJeCcc+EhTbjRvh7NOpf/IgH/fNRZ/TDL1SA3Z1b0GJq8zLLNo6uNROSkmTu2JzkFDGOzhzcUUClvYksKd5Z/euCYTcjLO+UKdPYddHCQjYHFcMzx6TgXzCedrLMWRREktyp7tzaBHlMoW0pXpS8qlqPzhhn5tRReoo8wtMaxR+aGL2E/kJVuMsde0YyRV3XQmxWE/pdUXKNPE5R4C70XIonhGevyIg+a0mqI15ayEfOrUokK/RdJpIHELcLYmhaI93uTv3oWZtoJ+AHJEvLFpLdM06y0PRGV3BzjbwBVFuSrCrFZlWqZkE4MbJ5JM2ZMs22F/g/FUTgIm9Enqw39UXOmgr9g4DoEWwupMi9S3BFXo8WQBDs7rqWnRBFFRTBQmcLI0uIKyJ+jyB5/CUvNeQUrqpCoCTET4BNVYWrj0AyphlNtHd19xN+kDd+YAVIuFzu/MZNyyJj8B7oNuTbYBquz7vTBkhBkgekSV5Q+S4EyK/mZO/gMMS7+I/D/QZIJekHUSrbEDPQEOqVIlLqOYirX3AzLlBXmPsEGo6lAPgu8NHUX1c0xKPoZtw5H0Txdaf7ZxO3alFUFkpDc+VbUxalTGjMVGyClYVkisYr59DW8CYb6JlQGiUCIq66epsUK0MV4maUq5YxY1Btg1mP1diYxwbBNFOzWUr+wMUTR3VGH8wbkA3Suli4ibmuoeH9P1HC9xBP4G//FWlI4WVFzzQpDVuD7EkKEAQannK1Nf5Xs/i/W7YwI4GqhBRbwYxpiH2/JUqtWEqbcaknUhSgwQyWAJoTNadvF2zGJCyWTVzjsfGwvrxz94Zf8I5f8Nuvol0urY3dYiaU2S31qHLN4iDhBwechT7Vfn6dIriCu6RmStE0rDeAhbR8Wbw9SW2yhJGVzBcbTtUzqpaGx2Fjrm9WG8PBnN7DenqkiCQJhELTU4Q/13GYKrjg/Zc0S/6q1EUJ8xpHf19H3UnUi0fRJRz6V8OmB6yPhgfuXo2iuHPT60/i8/6w1x9+hekxHqNCKGaowwCtolBh0MaW9S8jWlz/3AgAAA=="
  },
  "audit-binding13.log": {
    "sha256": "d20f423ada505e5f30cfa15215c7a6c30e398b79cfaf31b72ad8966057005307",
    "gzip_base64": "H4sIAAAAAAAC/32QwW7CMBBE70j8w/4C6pkDJVaVHgICKnGLVs5CLOxdy+uQ38c59FTT20jztDOz5no0+4tp+pP5LqI9dBCcquM7JJnXK/PXH6boncVMb4mI9gEZQ6RU9a1HF/4DYpIgmQawyIMbSlYVu6FXgilqToShinjRDFbKleyEq8jv3kyatUpM/GCZGVSmZOtVHD+Jl8bLlSoxJykh5XPESm/elkcgtRjrvs4YY4m4OU8woo5U6u4PJ9Pvfpr20n+2XdN2X/1xdz5DFC2Tn7TdANMdF6nbzcd69QK4G4+w8gEAAA=="
  },
  "render-report374.py": {
    "sha256": "cff1d0e71a51166553b9f21bdfe86139bc1e803272f0fc49d1aebef4f3118343",
    "gzip_base64": "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"
  }
}
```

## Final exact-source proof and Preflight

The final fuzz campaign runs from final21, whose42hashes equal the canonical reconstruction. It uses3seconds per target after the10second baseline campaign; every target passes, without source mutation. The initial campaign and its larger count remain historical supporting evidence. Product-comment-parity proves all21implementation bodies unchanged except whole-line comments in two files. All48test definitions/statuses were compared before closure; only TEST02 status changes below.

```json
{
  "preflight202.json": {
    "sha256": "c1f8f970cde649a160cc9c5c7df23ae207df286e032ec4e477a2df7509ac7cc9",
    "gzip_base64": "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"
  },
  "preflight202-process.json": {
    "sha256": "d795ae59110fd6af6d4a17d86e7e9ce2c70ad4ee9962a677fc1938112a4c505c",
    "gzip_base64": "H4sIAAAAAAAC/51TXU/bMBR9R+I/VH0mdew4H+YNSpiQuipCMA3NE/LHTZMtjT3HLUXT/vuSQAedtpc++Zzra59zr31/np5MJlPY1f5RGQ3T80l4NoaEW3U9+zKQns7P0X0HrkPLBZopoSpAQ/4ucJvW12voXql19Vq4530YabDQamhV3ae0wtdbQNY89TdV0DTIPnXVDHYwPXvVCZamcKasm3eh6/dsfs756ITz5YLzK+i+e2M5LzayqdUkb2oPk3nvpU/4lN/eXD885p/z+f3dxeUif1zcXN5e3D7MbIff7v84FL5nhYOyqVeVf9svnn1l2nwHauOFPPCCXvYiTJEd0WExH8y/T/1p5YW1V8ILtDBKNOgO1hbBUEGwJRkOolCVIYvTMGOaMqZlKksmicZJwmSqJVIb56D1gTemUZWoW7QySI7LoZEr45fgjzWjjW/BB1tMUvSC/34zDUcXGmEWsJiUGIchaI0poWVGiVCaxRlEgMs0Rm2vEJi2eR7RoXrR2vXx6gQHOKOS6piRBCjVKhaR0hiUlGFclmGaItsrBE9OWAvu8OsU3Y/mWO3SiVZVdQeB3NSNDsqIhGEaY0jKhFLFBGQqzERCdCIxhbifm86vHPSSveVZguyqh+Nz2x4cNiXf1sPUQSF89b/ZeXXG+WiN88Eb5/vGpDTIWJKQLFU4TmIaSS3jNGY4ixjVpY4ixrndTwsJyexbZ9rpIPX19OTX6clvgLrGWlsEAAA="
  },
  "final-fuzz23-receipt.json": {
    "sha256": "f3a8df38b1eb6bd5b827c553dfff6978d1cec4b987c577c0d73c3641e009a4f1",
    "gzip_base64": "H4sIAAAAAAAC/71W30/bMBB+R+J/qPrMj/xq0+6tdHRDK6NqgUlME7rYl9SqY2e20xIQ//uctCA2KRUPHk+R7y73+bs7f/bT4UGn09VkiTncr1FpJkX3U8c/auwZGLSr7per4++j64vb8+PJzd3d8ZfR9Xl3FyHf/GVXnd2qk0n/JOifDDobJqjc6FPIaT/a/ZWWj4/3huVN8lDvrEvQS9TW9FQvraFQMmW8CYowwL4XQBxHwzDxvGEMg57nRSmlfo+GPoZhHA4AvTAlgTccWh+kfi8Ih2SASRJvIbYbziWtU3pRAGTopZ4fEIx6vSjpY0p9DIZJD2KIoJfAAGMfezYPkIQEaRoi+HEURsNB6r1NqcvcphQl57XtecvHgMrQ1IR+biN3vGpmQFaQNcxOTpkwqATwU2oLkEhQVL/ktqECtmWa2JqNgaOgoEaCfptdXEqK/E2kNmDKGq47Gy0W3a39+eh92OlDC+aECWZwLMVLnx0BZiw1ZA/XOVLMC2MRFxYGL4EsmUBX6EvkBcF61QI/UoYRjl+ZNlJVTkvN/IFoQT1jAlQ146UCPgHOE5vAFSyXFXBTtSBPkWao/kepc1ArNExkrUOdF8AyMSIrITe83kduO+O05DLRqNaQMM5aKzCTnJFqKjNbhzGXGukt8BL1R+5hLMs6/FIKaaRg5COxm1HPFOSjLFPYSL8rdNGo2p4RWCJZcYs/U3YHqLXT5hdQKcl5W9u33vMHIOYz0pLUkuOs6YVs07eF1fIGdKSYWeb2hDjrtr03c/kvjb84o6rFz/5pB65wp+kKU1RWuto1fetf4O8SBUGnbVaY26eGvaRawK9RgDDzXdSUpUgqq/Du4NcMN+3Ma2/NV0HdGafMNcq9nG/mU22fDHOZSDO1yu5sujWXe18NZ6UVc7dKriVh0HaWF0QWSBf2NUvL17vbGXKp1li1HufGO0dtj5N2O9gbYKYWxxbkHzv35GJy5RjXnlXRfldxYPnrMXofcv35dXjwfPAHJFYQLHgMAAA="
  },
  "final-fuzz-profile.json": {
    "sha256": "4e2e602a77493b0097a85004fdd15d31e33738ae03fc2099500af15239c8ebb7",
    "gzip_base64": "H4sIAAAAAAAC/7VVTWvbQBC9B/wfjM+hIfRSenPcmoba1Nj9OJQQxrtjafBqR51d2VFC/ntHkhvoYU2lkpPYmdV7O19vnkYX4/EkmBwLuD+gBGI/eT++vmzt6GHr0KohSoWdbVc9Pt5HKlCtk7dh0lkjSIYxqO1ncx6Pn7qPukowe8ja62+uyEcUD+7KQsi3DGJPCO1VDx3sXDlm4NBbkKm3n1e3S7boJt3F58t/o9g9JKDn5CnijP2fgPvhZrSL5szL12ixKKMCbyJEXILJyWNPkhxdabA5JVimEsk4/EQhstRD8kPX73wC/IY8SL1ylYCbg3NbBeiJ7rgGF+sEwQJthvIf+SlA9hjJZ8n2KUqgzE/N3vPRNXSFpnNInngbUA6wJUfJeFbsyNQLzjSqmeOA9ju4CsMrUM24aq4v2XNkT+YVKNqmygSKaZYJZlqkviS+ne0z5cnR7J3SrESJMIQhhSmhFnYuVZLO+/EBTPyAtjLNRPYtSMmpKd+oPrXYU6GYF9qLfStRChfcPSoRAUojAfqn1rzsrVOCOxQd4LROdf4N/qrQGxxSAsGCvFUZTXB81RXi4/p0a0E7NLWqVm+WA+ExHUfjbV4v0KRzSBwB+WwE39aLoJtozVuOC1Wrvn0UHJ9dRjeVCtQgdQpsCFIzsDFcot3oerfVy67oS1DJAevkGLTeNQbtzzCohY5AsVGCBMGPk3t+O/8yDF573KfV1AEVL335F0HzuRtdPI8ufgMTzK4ZJgkAAA=="
  },
  "final-fuzz23.log": {
    "sha256": "476a1163e7be48e3aa3a52a23a0ed7341ed44b794fe71f87634cbb4d498a0bc0",
    "gzip_base64": "H4sIAAAAAAAC/62aTU/cSBCGz4m0/8FHLCWmqr8baQ8cstFedldKdg+5IDPTgBUzRraBkF+/ZTMfOGhm2tUc0VCPZ97XXV1d1c33LHsX6qoPRd0syvp00bThY7m8rbqualYfy/tl1Z9Wqz60K/p0WXY3l03ZLrt3UGgru/fNPMDVjyHQ6NmB19VVv9g8GGF2/E2o7xZh+IsAytjZgArdani2c7ND6+aprPunIdr72dG3Zfs99NXqmuINzpe8uexC+1BeVhTx9A4LMAzGanT9+UtoPf9H3JVPbVPXowRmfnQz2u7mvzZ3bXPb9PTBADBi/qPbcBXatqyf4zUj/rZaLUPbMYVrw0MVHsdoM/+t70IzCjf/de/qIdJINT+yWVTl4LRFMT/4vn0IT6PYMH+hPZZVX1ddP4QbwQgnp1fjQjUSZ4dPkqNxpNzV/c+fZ1moy7suLM8y6D5k12V/E1paSNll2RF7FbJF8xDa8jrQP5yiAvr79q4OfVjOjqfoKeFDtmoeswEzRDxW/U0msseGEkr76ttJoocfYdGdZUpY7VR2gkpokPQWLXIihcds/K2hG7IR/Vd20jd9WQ/PVXkEDiJJ/5x/+ZIgviTxBW9nMug4rmmVYJpWk3imZdJI6YAsE4jo9gktt0Jrm0fQIA7EMYz0loV1mFQJGCU5flmf4Jf1k3imX0L6UWHrrdtrl9iq7DCPgEEciGPXTnVyTbm0+mvMrvNtE8Lu9w2PxVP0lMB1DgQayljGSqv3Ce62ggup8wgaRJI43r3QXhUgHLf2paTAsc2bhNXmzSSe6RlKaYzPTpTSau9qg63S3uQRMIgDcRwbBaeF5kXCYcMM5SajBIGU7EjRUwLbMS+UJseMdPb42kC0eQQNIkkczzbCy8KBTTvl8XY1oVRKelRqSuAWIhYR3FCISONxn95ml9U05BE4iCRxnNtpP6RHlXzCRmU4/qUUkeo03Tk0QoLNTrQCb4/ntoMl/4YFURyOaVPRZSEspLdGPKsuScmX/vQNcqUFYamG0MZqE7Ep5REsiOK8iXEi3Tg69HCMky7BOekm8VzvEIWh05S04PZ6h1vN5WHz1jCIA73NsnNp3UTjX1cpGLHZgUvZ7MBNCXPcUy92J2roAtUViN7rvYWl3m1RKPMIHESSWAbuxFesTvqulTu2WRjpMqXF5afx7G6J8mosUujwhRGp7nCNsqZBHIjj2kZ0Ki6VZvbPnbCssiSpLlGnb1BTKiW8E0M/0hodYdeRduSaBnEgll3N0CdxFhMHFo51EECdcoDTk3juAvNSS6r/UCKaiO0ID564N7RIEMuxne6y8IxDwGRQ5IG1p4mkPU24KYG5pwkr6ACQnXiwbm9/C+FFW8rlETiIRXHc24mvCs+Y/Lyc0jneDMDZhEVHLQ6X3pmU6MbaHYGGZ+p4dnOHZwBrGsSBeLZtZKdkyRquboaj1vN6XV4krDiKnhLYLROF2gzViKKZ6fHaD/1h49Y4iCTxnHtWnhomVjPH0s7xpgCgkyp/PSWw5ze0I1ERYZ1S5viUjA4beQQNIkkcywbRVUFdOeZdACcVZ1NDkClLjMbPEwJzU6Mpm6EKXQojInqS1MB+5ZZ+BYNIEMssklwXgtH/316/cMg6naFLuoHgfiFwU6JGhcMBDSkjioiy7/CAdIODSBLLsmfheTvZ7t6LA83LiiYpK5opgZsVPYAjnb0xLqJuoEfmEbRYEsu1tfCqQOtTLhyNQx/O0EYkDW3ElMD1TdENHUpb9OoKE1Ooa5FH4CAWxXFuq70q6HyYctnLA+uuljjUyDpunf+FwM2UVBB4HIt+v78SwRedX8CDPcg1D2JRPO/W4lNBMhzXPv998df51z//+3Txx7/fvl18Pv/66WLgZn3ZXoe++13I397/D8K6WOCsLAAA"
  },
  "final-fuzz374.py": {
    "sha256": "ce71f3871c6b5555ccaaba7890b06318332017b29a3be04016191165bbd2ff57",
    "gzip_base64": "H4sIAAAAAAAC/41TWW/bMAx+D7D/4DfZWGwnTpo0DfyQ5mqxNClyYF27wpAlOtHmWKqkHOuvn+R2PYYO2INhHh8p8iOZSb51BNabnKUO2woutXNt1Mqz/EPxoqp2qZCcgFJVriqL2ALcJMlYDkniBRIUz/fgeoHAEgrdnZeIgByo61VEbHMEOcdUue4iRNnu8dE3+Wx4YH3IpsA00XDUrud1xV2JSTTbArqPUUOhbhnICpz7H4UfJNPwFF8+RndboVxRZQU19cSR9xl9N7gKFPuYMqJdrgIjM2maG8+Ws9mkf9G7nMYo5wTnyNi+zuZfYsSzzCrX89nNtxdtsboanD9r/fEsGU5755PhIEY1673q3Rh4fxGjyKqjSW9sZN8wg/Mc8ri09nv9i2GM+mfhSoFU4XQS9oQYYI3Dia0gXMJWhJCbrvx9o9nw69CEFHAtajaazbQZnQIl7ahOaLsenbTqOFxzn2CyAV8C4XuQQE27WK5VfPfumaBEhYRTOPpyV1iO1bMqJNti+euPOaQgwBJImIEUWLM9hIIfTKYN5HkoDmoTwBFQFflTfv00EauMyr/S8mXaP5k2FSblVNdYQ2ieSN4aAqHqyDOxJs0PIHrOuX5NYecevfizd+n/tRMGOwcCTOgPsFHD8mSdr+gxHx6B7DROyy7+bzbRad1v1EhW65y0a6cd2ux0aNpOs04a0Xqr1UnbNA3JTtqz8DXnOdlgVhguwrT8lfzdVw5Mb5y/SzQnsza7zc0QXHRMkedg5WRnD/HrOQaGR9dO2bRI+U7HmRVAyreYxXIwWy2rdvfNV7WjtdBWrWaOU7JCu2h0Oe1NktHq9jZqoOqDOUe9k4XdCq+LlaFBOwXXzltH5VPlNwOLtvk8BAAA"
  },
  "scope-review202.json": {
    "sha256": "35814ec84885b34c33b1dc9b2f3b7d4ed82e18ae62872c036ad4cdb155d56b72",
    "gzip_base64": "H4sIAAAAAAAC/42R3WrDMAyF7wN5h5Dr0j92EfYqpQilUhqBJ4PlBNKxd6+Tsm4tTeiFLBl952COv/OsKEp07qMC4kZUong1QCWwiLEzNuj01KKemcrPIoaOV5PIqxsgssXtHrgXYj0xINET9oWRg6CTCxM04pLfn9thRBIkmiBFtyG0tvYYyP6N67MvV2+R03te4rKrdDrmEV8bhx5rcRKHx9u8yNiPtQA4P9bSbkHchZ4H++13cOSOt3xrbnxgaFMGKdo+Zd3I/Qvy7CfPrubyCA/kAQAA"
  },
  "product-comment-parity.json": {
    "sha256": "b7b668b083b0096e584c703b002ea4f192ec06ed309a0daab434c40a1c5083c3",
    "gzip_base64": "H4sIAAAAAAAC/7VW34/bNgx+P+D+ByJPHRDnsLfh+tRixVasWw9Nhz0URUHLtC2cLGqSnCwb+r+P8o/Eaa7DtdWAOKJs6hP5kRL57voK4J/0B7ByGNvVLay0jeQtmpsKQ1sy+iosxE3Dq/W0QnHXkY0f2JrDB9WibSgIwrvx8xF50C2pZk/Lr8Prmxu4Q3WPDcFpD3Ced7qiAAiRLNpYBMWOKvjl7iUEiy60HEEAofZoVasDzUadcD059lHb5hYCGgprYF+Rl9HwAU08gGPxVeaedpr2stuOvFiyht/utoC2usREN6xJXss6bBpPDUYxrPbcQWzFC+5QWwhRvA3wxPVe8GrUplCGA1XfbVYn0PeLDVZYC++PJKjstakSPQn/E1aSJQqNIV+E3jmjxbwdmp7C5tKht2LyqAy8twEC914R/NmT14mykX7APrbs9d8YNdtEDThR4AowRq/LPr1+eokeJTLgZvNZCLFioSesPuFJIinc6/ogTkZMpkigWu3OyJrFj6MwzD+u/zOF67/k95iUfQxYo+uohuNwlHJBt2ScojRbiLnA9fc/2OEvF+B0fuYxF2yH/p7SgT1JuaC5DOR3WGqjxe6zWbYt7HA2k/knMRe4w4NncxyzwXJITzY4uXo4XQVhIeYC91ST92jCScoH3WmbSsNJygc9lJZ5zAUbiNOTDc5wev6/2i7gi6JuYfvq9Y3EkH1R9lVDEZQA3EoBlkI5nctUllOJkFN0WViGtTCtTeWo7L0FL6V4Dc70AX5ibgzB9s0L6HoTdbGXuPJ+Pc2SepHUHyjxUguHqwe2sTK6LBIDG3j2+9ufX7958ePXl+8zDmY40LboqGN/OHN+8GlJ0AN2eh3bTm5JtYHnszupJxIPUUiWvkBi6GSsYC+qC0qenhXmS+gWU6GGqDuaeAOja1IHZYTf7auXYrWkw9ALSBgqclIJZJsjdQ9W7fWXpc6vKVB/DLs/S7igO2do6LwyRPcWou8JUmxh35KFkoWi1MBJE+XHnDJsGxiRw5hgNQYJsqihfajZIe0laYMjFfWO5IVMWk6N2hPcsU6NLTayuWw7blMY0asgOH1P4Ss6wwuOpogHiHu+aAJVn1JjMG1MKsmSRk5YiJ9z5uTB5qJTnBLDsRtA2Uo/bHRjU4RG+uaU+Vxj+JzZkFwGqeMR1GN/eAy09LeLSM+hm/ItKXxrexhY6eOQ7TLt/Y4OYR5zwe5RR6NDPAr5gKWYWmmLZuHxwNdX76+v/gW3y1lqRQ4AAA=="
  },
  "comment-parity374.py": {
    "sha256": "b4e0e76d8a2c03c3ef19b6a7a5d3b78a24481dc7d8fe0aeeb51977c2caeccc69",
    "gzip_base64": "H4sIAAAAAAAC/31SwW4bIRS8r9R/4AbEG1b2cS0OUZJbHFu1e6hcC7G7rI3NAl7YxFHVf++j3qip1PYGvDfzZobX9q5DXsaD0RXSnXd9RCu4ZuP5GJzNG922UM/WPJWIEK02SgjKehWceVGEMi97ZeO8Pqj6pBq+3WWt65FH2iKyLnA99Kk+mxbaRtVbaTBle+Mqgm+KG7Z3mJYZ0i3yzMpOMWWb8KphFhZRhXhtqJ2N2g4qQ2cOnK0GmtkUFx50GBn1ixLR/TEN07nkqSwb4LlEEBq80dFoqwKh84qf/1nMUKNMlCF5QclMlPtc58f8lJtkawyFrdV5ULZWCxnBfU8kl3nFK7CnonC+dk1iA3fJHnBwjgEBAXzwg2QICtKWxpALMyH22ic5UfZxzKEoMP0l45KGy60uj7tJtT2VZkdzABAPikfJTHoPCZLvuFIAUbi89udYtpA+Lq+4H4AYP+w3Iu0CLv8XKZNBeBf0hdAc167r4F04a95EfZB2rwIurzKAf/RllCXjJMr5bJolSt+7Zqjj7UhxCyuk4xtLGwdDXuGirv+SXlgzdD68k+TaNgDhMzrB36A78z3sFcGrz8uHL/eb2VTcLxeLx+eNWD4/fRUPj0+bO7G6W6/RVWLD8SRlFoaOVM5B6Nu/O9l9yPzdAKWfsp/erc8RNwMAAA=="
  },
  "resume202.log": {
    "sha256": "93cb8da33e3943bed26f9142d2dc7a62c2e6ef28f8d90ac32ae7b0278edcff40",
    "gzip_base64": "H4sIAAAAAAAC/x2My2rDMBQF94X+g34gINuxHgstUkcUkxKXyi3pylxJV0QhkYOlFvr3dbs7MzBHn3T3PvbDcTLjbtTT684Ycl/mC7qi8BoLbq7RLrD8bG4QU8EEySG5zR6VPvVm7I/P5H6GjGrfm2740G+fJBcoX1k9vQzdQe/Jgt8xxzmpmtZk3R7XC1WJZgVMJf/7vwanfIa6ZYpxaKsgbNN6wa2VrOKhcuBFEEJKh5RuuQjMUoRGhBro1nkQVgrpgDngjw+/Jwfyj9gAAAA="
  },
  "plan202.log": {
    "sha256": "2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565",
    "gzip_base64": "H4sIAAAAAAAC/w3G0QrCIBQA0PfB/sEviHoN9iAmY7E5UfcQCMPaHVg3J13r+9t5Oppby8TZTwQf8qr3F6BX2bLX3zvGB5MYCzCxLUBem/EqhZulajslpelUO4tROcOFOzxpSwzhB9hkDKmu7DQM3NzYO6S4AhVqTiwHIlj2rCHinmNd/QHtCtYZgQAAAA=="
  }
}
```

Continuidad: reutilizar este censo/matriz por claim y sus hashes; el verificador rechaza deriva. No repetir la auditoría21ni la cualificación de entrenamiento V373 sin cambio material de fuente/claim/target. Próximo control: TEST03 sobre composición/grafos/runtimes actuales, preservando los hallazgos nativos y las exclusiones documentadas; TEST07 requiere después el artefacto final exacto. Los owners funcionales de FAIL385 se mantienen, no se duplican los dominios existentes para incorporar helpers aislados.

## Rectificación204

La distinción entre auditoría por claim e integración no permite marcar el control existente completo mientras conserva su bloqueo acordado. Checkpoint203 y las cifras de su cierre quedan como historial retractado. TEST02 vuelve a BLOCKED; las pruebas/correcciones canónicas se conservan y la investigación continúa en los owners funcionales. No se degrada un requisito ni se inventa un impedimento de negocio para justificar el cierre.

## Successor205 source-engine evidence

Portable views below normalize machine-local path prefixes to REFERENCE_STAGE. Listed SHA-256 values identify the original raw local receipts/logs, not the normalized display text. Original receipts remain outside the release tree. Preflight205 rejected the original display paths; successor206 reruns the unchanged gate. No execution result or binary digest was rewritten.


### low-traffic-red-receipt.json

SHA-256: c03171d826597d63592e8f71039d034f175739e1d234607a6d201fe6e2cd5f73

```json
{
  "exit_code": 1,
  "promtool": {
    "path": "REFERENCE_STAGE/elite-v372-7df147a09fe84a21962b3a63d2cca278\\promtool.exe",
    "bytes": 117974528,
    "sha256": "2d8428b4e7480a08e6c9b24c0e180669281192a361d92f63a3ba1421165eddf5"
  },
  "missing_alerts": [
    "low-traffic-all-errors",
    "counter-reset-still-all-errors"
  ],
  "log_sha256": "46fb615673300c0ecbc7207bb7d9631c5360c13b76c9b433d942f89db3fb23b8"
}

```

### low-traffic-canonical-receipt.json

SHA-256: bf0b47a3bbaf1720651fd312cb4bddfecb023379110e5c3324521e665358339e

```json
{
  "promtool": {
    "path": "REFERENCE_STAGE/elite-v372-7df147a09fe84a21962b3a63d2cca278\\promtool.exe",
    "bytes": 117974528,
    "sha256": "2d8428b4e7480a08e6c9b24c0e180669281192a361d92f63a3ba1421165eddf5"
  },
  "groups": 9,
  "assertions": 13,
  "canonical_files": {
    "ops/prometheus/platform.rules.yml": "cc718b7bbe35de3f230eeb47c854c270580dbb488101ee5f30c0b758e1d967b0",
    "ops/prometheus/platform.rules.test.yml": "3c5ed3e760f8a53cc9f6e4037e275177e1b6a3976e11e6dc6ca8d53485e3c643"
  }
}

```

### secure-ops-canonical-change.json

SHA-256: d85d0cd0ba30b45f7ce1a850a97300af0d9b4e4fcb863e7517662e6fc0775eeb

```json
{
  "pack": "SECURE-OPS-DELIVERY-CORE",
  "before_version": "1.1.3",
  "after_version": "1.1.4",
  "sha256": "3dc0b1179c3659d5803a63482ed6ac808cbcf9d8ce8910ba88992be24fdde697",
  "changed_consumers": [
    "ENTERPRISE_BACKEND_PACK_PLAN.md",
    "FRANCHISE_COMPLETE_PACK_PLAN.md",
    "FRANCHISE_SERVERLESS_PACK_PLAN.md"
  ]
}

```

### promtool-wrapper-focal.json

SHA-256: fb8a027cce8b3c6e003ea62e26d05eaaad551975d9703e4ae2f7722a09ba8435

```json
{
  "absent_engine_not_executed": true,
  "steps": [
    {
      "id": "prometheus-rule-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 1316.0
    },
    {
      "id": "prometheus-rules-source-engine-regression",
      "status": "PASS",
      "elapsed_ms": 482.0
    }
  ],
  "artifact": {
    "path": "REFERENCE_STAGE/elite-v372-7df147a09fe84a21962b3a63d2cca278\\promtool.exe",
    "sha256": "2d8428b4e7480a08e6c9b24c0e180669281192a361d92f63a3ba1421165eddf5",
    "rule_sha256": "cc718b7bbe35de3f230eeb47c854c270580dbb488101ee5f30c0b758e1d967b0",
    "fixture_sha256": "3c5ed3e760f8a53cc9f6e4037e275177e1b6a3976e11e6dc6ca8d53485e3c643",
    "production_admission": false
  },
  "negative_cases": 5
}

```

### low-traffic-red.log

SHA-256: 46fb615673300c0ecbc7207bb7d9631c5360c13b76c9b433d942f89db3fb23b8

```text
  FAILED:
    name: low-traffic-all-errors,
    alertname: PlatformHighErrorRate, time: 20m, 
        exp:[
            0:
              Labels:{alertname="PlatformHighErrorRate", severity="page"}
              Annotations:{runbook_url="https://replace.invalid/runbooks/api-availability", summary="User-visible server error rate exceeds one percent"}
            ], 
        got:[]

    name: counter-reset-still-all-errors,
    alertname: PlatformHighErrorRate, time: 35m, 
        exp:[
            0:
              Labels:{alertname="PlatformHighErrorRate", severity="page"}
              Annotations:{runbook_url="https://replace.invalid/runbooks/api-availability", summary="User-visible server error rate exceeds one percent"}
            ], 
        got:[]



```

### low-traffic-green-suite.log

SHA-256: b638c1ecba3a10900d6e3b081d45992d3c970b4271ac1f5bb8ca09a63085c4c8

```text
  FAILED:
    name: recovery-after-firing,
    alertname: PlatformHighErrorRate, time: 20m, 
        exp:[], 
        got:[
            0:
              Labels:{alertname="PlatformHighErrorRate", severity="page"}
              Annotations:{runbook_url="https://replace.invalid/runbooks/api-availability", summary="User-visible server error rate exceeds one percent"}
            ]



```

### low-traffic-canonical-check.log

SHA-256: 2797a7710be74703e0a64116bc2f17a4b152cd7a33f6bbad78ba4d9b965f4d0d

```text
Checking REFERENCE_STAGE/elite-v374-8966287c156543bdb575918394dfd339\ops-final\ops\prometheus\platform.rules.yml
  SUCCESS: 3 rules found


```

### low-traffic-canonical-suite.log

SHA-256: 8289a4919870ab5b28f5e90a86b4af3cdbd8c244516c02fca90d6e3f289a6711

```text
  SUCCESS


```


## Successor207 — final verification of the operations correction

FAIL712 is fixed and regression-proven. Current13assertion fixture against the original rule fails exactly low-traffic outage, counter reset and recovery-after-firing; the corrected canonical rule passes all13. The37to38file comparison finds one changed expression,36unchanged existing files and one new fixture. Pack1.1.4 preserves claim/admission, threshold,hold and other alerts. Original203 audit-only control promotion stays retracted; no state is advanced for this narrower fix.

### operations-final-delta.json

```json
{
  "canonical_fixture_on_original_rule": {
    "exit_code": 1,
    "failed_scenarios": [
      "low-traffic-all-errors",
      "counter-reset-still-all-errors",
      "recovery-after-firing"
    ],
    "raw_log_sha256": "362f78eb195dd5579e50aee81504498cd30838004e405f5443d5f5b2f4898570"
  },
  "canonical_file_delta": {
    "before": 37,
    "after": 38,
    "existing_changed": [
      "ops/prometheus/platform.rules.yml"
    ],
    "added": [
      "ops/prometheus/platform.rules.test.yml"
    ]
  },
  "other_existing_files_byte_identical": 36
}

```

### operations-preflight206-summary.json

```json
{
  "preflight": 206,
  "exit_code": 0,
  "executed_steps": 162,
  "passed": 162,
  "failed": 0,
  "overall_availability": "BLOCKED",
  "missing_tools": [
    "docker"
  ],
  "profiles": 55,
  "changed_compositions": {
    "ENTERPRISE_BACKEND_PACK_PLAN.md": "399",
    "FRANCHISE_COMPLETE_PACK_PLAN.md": "755",
    "FRANCHISE_SERVERLESS_PACK_PLAN.md": "558"
  },
  "raw_receipt_sha256": "55b6835f80050f1c1c21c8d00aba54e7073a4884d2278573717dd74e778b5c68",
  "raw_log_sha256": "1adcc35841e472bc263f68f120238e3105981b2c917c35f8ed5cec640e02359a",
  "prometheus_rule_sha256": "cc718b7bbe35de3f230eeb47c854c270580dbb488101ee5f30c0b758e1d967b0",
  "prometheus_fixture_sha256": "3c5ed3e760f8a53cc9f6e4037e275177e1b6a3976e11e6dc6ca8d53485e3c643",
  "rule_executed": true,
  "production_admitted": false
}

```


## Successor208 — canonical continuity map corrected

The lead-in still said currentV305 and67/743 despite later appended work. Header/count are corrected and dated sections explicitly historical; the observability alternate owner credits V372 without promoting the isolated core or whole-service equivalence. No execution code changed after Preflight206.

```json
{
  "path": "markdown_system/FRANCHISE_GAP_MAP.md",
  "before_sha256": "a432c3ca065c0dcc0254e73f97fe0ee4aa72194f9611ea3e0b691b131f305fee",
  "after_sha256": "c406a0d3fc165098937df6eb804017f8a36335bd97385dffed2b77d37791b0c2",
  "before_header": "Estado vigente V305, 2026-09-07. Owner: roadmap sección10, sin segundo plan.",
  "after_header": "Estado vigente V374 / checkpoint208, 2026-09-10. Owner: roadmap sección10.",
  "before_current_count": "67/743",
  "after_verified_count": "67/755",
  "evidence": "Preflight206 / V372 / V306-V312",
  "control_status_changed": false
}
```
