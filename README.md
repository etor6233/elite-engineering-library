# Elite Engineering Library

Biblioteca Markdown portable para que Codex, Claude Code u otro agente diseñe, componga y materialice sistemas verificables sin imponer un backend TypeScript ni copiar código público sin admisión.

## V402: referencia local READY_FOR_LIBRARY_USE

**READY_FOR_LIBRARY_USE / LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES** corresponde al cierre local V402: ejecución **337 COMPLETE**, T2801–T2810 y ARCA_INFRA **PROVEN_LOCAL**, 48/48 controles de ese alcance. `production_authorized=false`.

**Alcance de este repositorio público:** este árbol publica la composición V402 que estaba pendiente en la carpeta canónica local (packs, evidencias, estado y scripts). Los ZIPs inmutables del cierre original **no** se adjuntan aquí y sus SHA no cambian. Para una instantánea byte-a-byte del READY original, usar el ZIP de biblioteca identificado abajo.

### Qué está listo / Qué NO está listo

- **Listo en la referencia local:** perfil de 116 packs / 1.653 archivos, integración con fixtures, reconstrucción independiente, producto firmado y ensayos locales de recuperación y materialización NEW/EXISTING.
- **Producto/target:** `PROJECT_READINESS_GATE.json` conserva `DISCOVERY` y `platform_mode=BLOCK`. No hereda el READY de biblioteca ni autorización para construir o desplegar un producto real.
- **Preflight genérico:** `BLOCKED` para pnpm general y Docker. El pnpm general está rechazado por admisión de seguridad; Docker está ausente. La referencia ya probada usa el instalador restringido y PostgreSQL nativo admitidos. No se afirma PASS del diagnóstico general.
- **Live:** proveedores condicionados a credenciales; reglas, jurisdicción y aceptación operativa pertenecen al target. El checklist de cuentas permanece vacío. ARCA tiene infraestructura/fixtures completos; Daybreak/libxml2 sigue `ACCESS_BLOCKED`, diferido sin investigación.

### Dos ZIPs: cuál usar para qué

| Artefacto inmutable | Cuál usar para qué | SHA-256 del ZIP completo |
|---|---|---|
| `elite-library-v402-source330-meta331-distribution336.zip` | Biblioteca portable: extraer y materializar el perfil aceptado en otro destino. Es la instantánea del cierre original. | `cd24757744d01dbc37b666495d659f76824a8a680015eee469187e975542591f` |
| `signed-release-v402-r330/artifact.zip` | Producto de referencia: verificar con su firma, manifiestos y política de confianza externa; usar en el alcance local/fixtures demostrado. | `dba7978d04d08dc2eedb112524708e220ba18cccf86c3a2142f8b1556ed1bf76` |

Son artefactos locales distintos. Este commit no los adjunta, regenera, reempaqueta ni cambia sus SHA. La carpeta durable está en `Desktop/Elite Franchise Reference V402`; dentro están `START_REFERENCE_V402.md` y `qualification/FINAL_LIBRARY_READY_V402.json`. Estas rutas son locales, no enlaces de descarga de GitHub. El recibo final tiene SHA-256 `01cc08637802fc01af8edca825eae0fc3d8b99e2dc9bb508528c957b1e41e5c5`.

**El árbol de este repo refleja el post-cierre V402; el ZIP sigue siendo la instantánea canónica del READY original**.

Leer el [expediente de gates y alcance](reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md), la [guía de arranque](START_FRANCHISE.md) y el [gap map con historia delimitada](markdown_system/FRANCHISE_GAP_MAP.md). El flujo es bridge → START_FRANCHISE → verify sobre la copia exacta elegida. No se afirma BENCH01, ahorro de tokens ni una franquicia en menos de una semana.


## Requisito de tooling

La biblioteca se puede leer sin instalar runtimes. Sus scripts de bootstrap y distribución requieren **PowerShell 7 o posterior** (`pwsh`); no son compatibles con Windows PowerShell 5.1. Los validadores de readiness y continuidad requieren además **Python 3.14+**, aunque aún no exista producto. Pueden usarse ejecutables aislados ya disponibles; no hace falta una instalación global. Los runtimes de producto —Go, PostgreSQL, Python, .NET para WCF, Node/pnpm u otros— se instalan sólo cuando los packs o fuentes seleccionados los requieren. Kiota 1.35.0 dispone de lane Windows x64 autocontenida verificada y no necesita .NET SDK para ese uso estrecho.

## Integrarla en un proyecto

`AGENTS.md` sólo gobierna su propio directorio y los descendientes. Copiar esta biblioteca como un subdirectorio no hace que sus instrucciones se apliquen automáticamente al proyecto padre. El instalador crea un bridge mínimo, una Skill de carga progresiva y conserva las instrucciones existentes; no crea Git, no instala dependencias, no pide secretos y no materializa producto.

### Modo A — proyecto descendiente, descubrimiento automático

1. Extraer la biblioteca en una ubicación estable.
2. Crear el proyecto dentro de un subdirectorio de la biblioteca, por ejemplo `projects/mi-proyecto/`.
3. Instalar el bridge desde la raíz de la biblioteca:

```powershell
pwsh -NoProfile -File .\INSTALL_AGENT_BRIDGE.ps1 -ProjectRoot .\projects\mi-proyecto -Agent Both
```

4. Abrir Codex o Claude Code en la raíz exacta del proyecto. Esto es obligatorio para descubrimiento Codex cuando no existe una raíz Git.
5. Expresar el objetivo del negocio: el bridge genera una ronda explicativa por vez, define términos, muestra un ejemplo no asumido y pide sólo decisiones, evidencia o accesos materiales faltantes.

### Modo B — biblioteca vendorizada, bridge explícito

1. Copiar la biblioteca, sin modificarla, dentro del proyecto; por ejemplo en `tools/elite-engineering-library/`.
2. Desde la raíz del proyecto, ejecutar:

```powershell
pwsh -NoProfile -File .\tools\elite-engineering-library\INSTALL_AGENT_BRIDGE.ps1 -ProjectRoot . -Agent Both
```

3. Abrir el agente en la raíz del proyecto y formular el objetivo. Codex encontrará `AGENTS.md` y `.agents/skills/elite-engineering-library/SKILL.md`; Claude Code encontrará el import administrado en `CLAUDE.md` y `.claude/skills/elite-engineering-library/SKILL.md`.

Repetir el comando actualiza sólo el bloque administrado y es byte-idempotente si la ubicación no cambió. Un bloque incompleto, una Skill homónima ajena o un `AGENTS.md` resultante mayor de 32 KiB se rechazan antes de escribir. Para previsualizar sin mutar usar `-WhatIf`; para instalar sólo uno de los dos agentes usar `-Agent Codex` o `-Agent Claude`.

Claude Code puede solicitar una aprobación la primera vez que un `CLAUDE.md` importa un archivo externo al proyecto; esa aprobación debe concederse en la interfaz para habilitar el import y no se sustituye con texto. El bridge y las Skills son orquestación local `AUTHORED`, basada en los mecanismos oficiales de descubrimiento/progressive disclosure de OpenAI y Anthropic; no se presentan como código de producto de esas empresas ni como sustituto de los packs admitidos.

Para preparar el tooling, interpretar el bloqueo inicial y reanudar con comandos
comprobados, seguir [Entrada CLI y recuperación](markdown_system/PROJECT_ENTRY_CLI_GUIDE.md).
La guía conserva el baseline NEW/EXISTING y no fabrica respuestas ni checkpoints.

## Verificar una copia

Antes de confiar en una copia o modificar la biblioteca, ejecutar desde su raíz:

```powershell
pwsh -NoProfile -File .\VERIFY_LIBRARY.ps1
```

Para reconstruir además el pack de adquisición oficial y comprobar qué toolchains faltan antes de iniciar un sistema:

```powershell
pwsh -NoProfile -File .\VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight
```

`Preflight` nunca convierte toolchains o cuentas ausentes en un PASS. `-Mode Foundation` compone los perfiles, ejecuta Go/Python y el typecheck/tests/build web; sólo se usa cuando los runtimes requeridos están disponibles. Amazon, Google Merchant, Meta Ads, TikTok, Meta WhatsApp, Firebase FCM y Mercado Libre se reconstruyen/validan estructuralmente en `Audit`; WhatsApp ejecuta además 22 tests stdlib, con scope WABA/teléfono y statuses v2; el bridge Go al fence existe, mientras la correlación de statuses sigue pendiente. Los adapters AWS S3/SES y Google Cloud Storage se materializan y validan siempre; si Go 1.26.7 está disponible —o se pasa su ruta exacta con `-GoExecutable`— Audit ejecuta también sus module verify, tests, vet y build offline. Las instalaciones/tests hash-fijados de los cuatro SDK Python providers, los tests Go del SDK Firebase exacto y el adapter HTTP Mercado Libre sin dependencias se ejecutan en Foundation sólo con red explícitamente autorizada o caches ya verificadas; TikTok exige source root+receipt aprobados. No existe adquisición, consumo de red o minutos implícito. Integraciones, documentos, PostgreSQL, seguridad, carga, recovery y deploy conservan gates propios del proyecto.

En Linux/macOS el mismo comando usa separadores POSIX:

```bash
pwsh -NoProfile -File ./VERIFY_LIBRARY.ps1
```

## Contenido del snapshot público base — histórico

- conocimiento y mapas de autoridad;
- contrato de blueprint y cobertura de 48 superficies;
- 204 implementation packs con 2250 archivos reconstruibles y 56 perfiles; la franquicia compone 109/1481 e incluye Mercado Libre Questions v4 hacia el owner omnicanal común y respuesta aprobada con fence/confirmación/reconciliación sin retry ciego. El perfil readiness 1/11 incluye catálogo A–H, prompts deterministas, alcance explícito de biblioteca y 102 regresiones; `EXECUTION-VALIDATOR 1.3.1` 1/15 añade NEW/EXISTING, baseline/delta e `implementation_assurance`; el resolver de gaps 1/4 obliga a investigar fuentes oficiales actuales; Kiota 1.35.0, fuzzing Go, DevSkim adaptado y release portable siguen disponibles sin Git ni CI pago;
- perfil email AWS 14/150 con receptor SES, raw Object Lock, MIME, GuardDuty exact-version, Magika→AWS IDP, Powertools oficial, handoff/review 1/8, strict-field offline 1/5, decision worker v2 1/8, PostgreSQL foundation 1/3, persistence boundary Aurora 1/11, runtime Debezium 1/10, consumidor inbox/proyección+CEL 1/24 y seguridad local. El consumer lee el payload revisado, exige mapping single-line ID/clase/version/SHA, schemas/corpus y claves exactas, rechaza eventos de otra clase documental y transacciona inbox+payload mapeado+receipt antes del ACK; 18/18 pruebas PostgreSQL, concurrencia 4×8 y SCA 170/0 PASS desde Markdown. Sigue condicionado a cuenta/MX/AWS/Aurora/Kafka/costo/canaries/config/corpus/HITL/load/DLQ/RLS/TLS/auth/rebalance/restart/redelivery/order/restore/rollback y mappings finales aprobados; no afirma exactly-once ni deploy live;
- adquisición sin Git de 124 fuentes oficiales fijadas por commit/SHA-256 —incluido el commit Microsoft BCApps firmado `2eae56d` con Inventory/Warehouse/SCM-Reservation inventariados, además de Kubernetes/kubectl 1.37.0, OpenID Foundation Conformance Suite 5.2.4, Cloudflare cloudflare-go 7.9.0, OWASP ZAP 2.17.0, PostgreSQL 18.6, Grafana k6 2.2.0, Microsoft Playwright 1.62.1 y GitHub Spec Kit 1.0.1. El runner usa diecisiete perfiles hash-linked. Kubernetes sólo produce `DEPLOY_ROLLBACK` después de digest inmutable, canary, rollout, rollback real, digest previo restaurado y migración compatible; BCApps conserva su requisito Business Central/AL/runtime/licencia; Spec Kit más Playwright sólo producen `BUSINESS_ACCEPTANCE` después de escenarios/effects reales, aprobaciones independientes, defectos gobernados y owner signoff. Prestigio, source disponible o exit 0 nunca sustituyen evidencia—;
- admisión de proveedores seleccionados sobre identidades exactas de Stripe, Mercado Pago, Amazon, Google, Meta, TikTok y Firebase: exige sandbox/test account por operación, webhook autenticado y durable o N/A con doble revisión, idempotencia/dedup y reconciliación, retry acotado ante throttling y términos/costo/scopes aprobados. TikTok Lead usa el wheel oficial 1.1.3 por URL/SHA, transporte genérico oficial, contratos v1.3 fijados, artefactos hash-linked e import Go al owner PostgreSQL/outbox compartido; el wrapper Lead es `ADAPTED`, no generado por TikTok, y todo webhook queda `UNAUTHENTICATED_PROVIDER_SIGNAL` hasta demostrar un mecanismo oficial. Mercado Libre se declara como glue `AUTHORED` sobre HTTP oficial porque su SDK archivado no se presenta como código reusable; Questions outbound exige approval durable, un solo POST y confirmación/reconciliación GET-only sobre el ledger común;
- Amazon Easy Ship v2022-03-23 materializa 2/34 para listar slots, programar pickup/dropoff y reconciliar el estado oficial. El schedule persiste intento antes del POST y bloquea retry ambiguo; `PickedUp` nunca completa automáticamente el handover interno. Cuenta, support table, sandbox y target real siguen condicionados;
- Amazon Fulfillment Outbound v2020-07-01 materializa 2/33 para recuperar foto, firma o PDF de entrega MCF reportados por `getFulfillmentOrder`: liga pedido/marketplace/paquete al receipt empresarial, valida host/tipo/status/bytes, guarda binarios por SHA-256 y excluye URLs firmadas, tracking e identidad. Sus 12 pruebas pasan con `amzn-sp-api==1.11.1`; evidencia provider nunca equivale a aceptación empresarial automática;
- Amazon Supply Sources v2020-07-01 materializa 2/34 para inventariar y gobernar tiendas/depósitos mediante list/get/create/update/status/archive oficiales. Sus 14 pruebas ligan cada mutación a inventario/snapshot, registran attempt antes del efecto, reconcilian o dejan `UNKNOWN_EFFECT` sin retry y excluyen dirección/contacto de receipts. No reemplaza la red multi-location interna ni activa ubicaciones automáticamente;
- Amazon Multi-Location Inventory materializa 3/43 combinando Supply Sources con Listings Items v2021-08-01. Sus 13 pruebas usan `PRODUCT` y `/attributes/fulfillment_availability` oficiales, ligan request/receipts/snapshot por SHA, detectan drift, escriben attempt antes del PATCH, reconcilian por GET y bloquean efecto ambiguo; no persiste seller/SKU/source IDs crudos ni hace commit automático al inventario interno;
- Amazon External Fulfillment Inventory v2024-09-11 materializa 2/34 para FETCH/UPDATE batch de inventario absoluto location-level desde POS/ERP/WMS. Sus 12 pruebas usan modelos oficiales, máximo 10 identidades, sequence por ubicación+SKU+marketplace+canal, approval/drift/attempt y reconciliación; rechazo parcial queda ambiguo y ningún acknowledgement hace commit interno automático;
- admisión productiva completa como maquinaria: los ocho controles tienen adaptador semántico —edge/CDN/WAF, identidad/autorización, proveedores, recovery PostgreSQL, carga/resiliencia, seguridad ofensiva, deploy/rollback y aceptación empresarial—. El template distribuido de cada uno permanece bloqueado hasta recibir evidencia real del proyecto;
- backend Go/PostgreSQL configurable y adapter web Next.js opcional;
- perfil backend 29/386: incorpora reserva serial/ATP, bulk lot/bin/FIFO-específico, UOM inmutable, composición exacta, breakbulk/gather y recepción→put-away→pick/reposición packaging-aware; conecta pedido colocado→pick→customer shipment→transporte carrier→reportes inmutables→entrega técnica del pedido, con replay, concurrencia, consumo atómico y cero aceptación automática del handover. Conserva FEFO, min/max, cross-docking y transferencias. Full Go, contratos .NET y PostgreSQL 18.6 0001–0039 pasan desde Markdown. Cuenta/webhook/polling carrier, evidencia/custodia, handover/checklist/aceptación, facturación/pago-crédito, customer shipment serial/específico, service/production, homologación ARCA, accounting/planning, calendarios, cuentas live y producción target siguen fail-closed;
- perfil web 6/84: Google SafeValues 1.2.0, CSP estricta con nonce fresco y 38 tests; agrega recepción/disposición de devoluciones, evidencia calculada server-side y rechazo de refund, stock, accounting, fiscal, IDs o estados browser-owned. Frozen install/typecheck/build/licencias pasaron desde Markdown; Playwright 4+8 y Lighthouse cinco runs permanecen gates independientes. Cada efecto se muestra como solicitado, nunca completado sin su owner real;
- identidad, autorización por recurso, workers, provider edge, puerta local Magika/ClamAV/YARA-X con política previa y receipt, landing SFTP privado Microsoft AVM con SSH-only/`cw`/Private Endpoint y retención separada obligatoria, samples byte-verbatim Microsoft Azure Document Intelligence y Azure Content Understanding `prebuilt-invoice`, lifecycle oficial Azure CU same-resource analyzer copy condicionado, samples Google Cloud Document AI `process_document` y Custom Document Extractor con `schema_override` variable por solicitud, lifecycle oficial processor/dataset/import/train/evaluate/deploy/default/undeploy, batch GCS y response handling OCR/form/table/entity/split/layout/custom, más fixtures oficiales hash-locked —invoice y packing-list input/output— y componente oficial Amazon Textractor 1.10.0 con wheel, tests Queries y fixtures fijados; evaluación estructurada con AWS Labs Stickler/schema cerrado/cero FA-FD-FN-FP sin storage authority, orquestación durable Microsoft desde seguridad hasta evidencia final/persistencia idempotente, transactional outbox Dapr 1.18.3 oficial condicionado —con SDK Go vulnerable rechazado—, evidencia inmutable condicionada en AWS S3 Object Lock, Google Cloud Storage Object Retention o Azure Blob version-level WORM, log tamper-evident local Transparency.dev Tessera sobre Linux/POSIX con checkpoints independientes obligatorios, Amazon Catalog read-only, Google Merchant product sync/status, Google/Meta/TikTok Ads reporting read-only, Meta WhatsApp template/webhook adaptado con provenance exacta, Firebase FCM dry-run/send condicionado, Mercado Libre item/order read + validator + notification fetch jobs + Questions outbound aprobado/fenced/reconciliado, observabilidad, CI, backup, packaging y evidencia de licencias;
- materializador, compositor, preflight/runner ejecutable, recuperación autocorrectiva, memoria persistente de fallos, lifecycle de dependencias y autocorrección de información obsoleta;
- instalador de bridge Codex/Claude sin Git y Skill Codex de carga progresiva, con preservación, idempotencia, límite de contexto y regresiones negativas;
- evidencia histórica separada de los claims vigentes.

El router `markdown_system/MARKDOWN_SYSTEM_READINESS.md` conserva el alcance de bootstrap del snapshot público base. Para el cierre local V402 manda la sección de gates anterior y su expediente; no trasladar estados históricos a esa aceptación. `READY_FOR_PROJECT_BOOTSTRAP` significa que la biblioteca puede iniciar y acelerar un proyecto; la admisión productiva depende siempre del país, reglas de negocio, proveedor, IdP, infraestructura, carga, datos y artefacto concreto.

## Crear una copia portable

```powershell
pwsh -NoProfile -File .\CREATE_PORTABLE_ARCHIVE.ps1 -OutputPath "C:\ruta\Elite-Engineering-Library.zip"
```

En Linux/macOS:

```bash
pwsh -NoProfile -File ./CREATE_PORTABLE_ARCHIVE.ps1 -OutputPath "/ruta/Elite-Engineering-Library.zip"
```

El ZIP sólo incluye los Markdown y scripts canónicos de la raíz y de las cuatro carpetas admitidas. Excluye `.git`; rechaza caches, dependencias instaladas, outputs generados, secretos, ZIP anidados y archivos desconocidos en vez de empaquetarlos silenciosamente. Incorpora `DISTRIBUTION_SHA256SUMS.txt`, verifica cada entrada contra ese manifest y genera un checksum SHA-256 lateral del ZIP completo.

Las entradas se guardan sin compresión, en orden ordinal y con fecha UTC fija. `-SourceDateEpoch` permite fijar un segundo Unix par entre 1980 y 2107; el valor predeterminado es `946684800` (2000-01-01). Con los mismos bytes, epoch y runtime verificado, el ZIP es reproducible aunque cambien las fechas de los archivos. Este formato ocupa más espacio que un ZIP comprimido.

La publicación reserva ambos destinos sin sobrescribir archivos existentes y retira sus propias salidas parciales ante una excepción normal. Una interrupción del proceso o del equipo puede dejar un par incompleto: verificar siempre el ZIP y su `.sha256` antes de usarlo. Reproducibilidad y checksum no sustituyen firma, admisión ni los gates del release final.

El archivo es una copia de fuentes, no un bundle offline de toolchains o dependencias. Guardar el ZIP y su `.sha256` juntos y conservar otra copia en un medio independiente si se necesita recuperación ante pérdida del equipo.
