# START_FRANCHISE — arrancar desde la copia elegida

## Qué está listo / Qué NO está listo

La **referencia local V402** está `READY_FOR_LIBRARY_USE`, ejecución 337 COMPLETE, en `LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES`. El gate de producto `PROJECT_READINESS_GATE.json` conserva `DISCOVERY` y `platform_mode=BLOCK`; no se lo promueve ni se transfieren las aprobaciones del mantenimiento al target.

Para reutilizar exactamente el perfil aceptado, extraer el ZIP de biblioteca distribution336 de la [tabla única de artefactos](reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md#dos-zips-cuál-usar-para-qué). El ZIP firmado de producto sirve para verificar y usar la referencia local. Los comandos siguientes se ejecutan desde la raíz de la copia elegida; el árbol Git ya publica los packs V402 y conserva aclaraciones posteriores; no equivale byte por byte al ZIP inmutable. El perfil vigente tiene 116 packs / 1653 archivos.

El Preflight genérico permanece BLOCKED por pnpm general rechazado y Docker ausente. La referencia V402 se probó con el instalador restringido y PostgreSQL nativo admitidos. Live y aceptación del target siguen pendientes; ARCA está completa en infraestructura/fixtures y espera credenciales para conectarse. Daybreak/libxml2 sigue diferido. El checklist de cuentas queda vacío: no se pide ningún secreto en este pulido.

En la instalación local, abrir `Desktop/Elite Franchise Reference V402/START_REFERENCE_V402.md` y `qualification/FINAL_LIBRARY_READY_V402.json` dentro de esa carpeta. La [separación de gates](reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md) explica FAIL384 y las fuentes. Ningún paso siguiente obliga a reabrir el mantenimiento337 ya completo; aplica al proyecto consumidor y a la revisión que elija.

## Paso 0 — Separar biblioteca y producto

Conservá la biblioteca en una ubicación estable y elegí un destino de producto separado. No hace falta Git. Ejemplo PowerShell 7 para una carpeta nueva:

```powershell
# Ejecutar desde la raiz de la copia de biblioteca que elegiste.
$libraryPath = (Get-Location).ProviderPath
if (-not (Test-Path -LiteralPath (Join-Path $libraryPath 'START_FRANCHISE.md') -PathType Leaf)) { throw 'Abrir PowerShell en la raiz de la biblioteca' }
$projectPath = Join-Path (Split-Path -Parent $libraryPath) 'Mi Franquicia'
New-Item -ItemType Directory -Path $projectPath -ErrorAction Stop
```

## Paso 1 — Conectar el protocolo operativo

La fuente única del paso a paso es [FRANCHISE_PROJECT_OPERATING_PROTOCOL](markdown_system/FRANCHISE_PROJECT_OPERATING_PROTOCOL.md): inicio/reanudación → alcance completo → autoridades/packs → incremento conectado → verificación → registros → próxima acción. Define carga mínima, reutilización de evidencia válida, trabajo independiente ante bloqueos y un solo escritor del estado.

El complemento vive en un pack Markdown para conservar intactos el instalador y los verificadores V402. Requiere PowerShell 7+ y Python 3.14+. Desde la copia elegida de la biblioteca, materializá sólo el tooling de conexión en un destino ausente:

```powershell
$connectionTools = Join-Path $projectPath '.elite/connection-tools'
& "$libraryPath/materialize_markdown_pack.ps1" -PackFile "$libraryPath/implementation_packs/PROJECT_OPERATING_CONNECTION.md" -Destination $connectionTools
if (-not $?) { throw 'No continuar: falló la materialización del complemento' }
& "$connectionTools/operating_connection/INSTALL_PROJECT_OPERATING_BRIDGE.ps1" -LibraryRoot $libraryPath -ProjectRoot $projectPath -Python 'python'
if (-not $?) { throw 'No continuar: falló la conexión operativa' }
```

El complemento **invoca INSTALL_AGENT_BRIDGE.ps1 original**, conserva sus bloques y agrega una entrada independiente a `PROJECT_AGENT_ENTRY.md`. También genera un receipt de rutas/SHA y una plantilla neutral de plataforma. No crea blueprint, estado ficticio, producto, cuentas ni aprobaciones. Para repetir una instalación ya materializada, ejecutar sólo el instalador complementario: verifica integridad y devuelve UNCHANGED. No rematerializar encima de archivos existentes.

El ZIP distribution336 permanece inmutable y no contiene esta extensión posterior. Para usar la conexión nueva, elegir el árbol vivo que incluye el protocolo y su pack. La [separación de gates y artefactos](reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md) sigue vigente. El informe de este delta es [PROJECT_OPERATING_CONNECTION_REVIEW](reconstruction_evidence/PROJECT_OPERATING_CONNECTION_REVIEW.md).

## Paso 2 — Abrir el proyecto y seguir sus archivos

Abrí el agente en la raíz exacta del proyecto. La entrada instalada es `PROJECT_AGENT_ENTRY.md`; AGENTS.md y CLAUDE.md conducen a ella. Para otro agente, usar `.elite/platform-adapter.template.md` y completar las rutas y capacidades realmente disponibles. Una instrucción inicial suficiente es «Leé PROJECT_AGENT_ENTRY.md y seguí el protocolo y el estado del proyecto». El método, los pendientes y las autorizaciones deben estar en archivos; no dependen de este mensaje.

Grok u otro host remoto necesita acceso real a esa biblioteca y proyecto en su entorno. La plantilla no instala conectores ni demuestra descubrimiento o ejecución de una plataforma. Todo adapter no ejecutado conserva NOT_TESTED.

## Paso 3 — Aplicar los gates del proyecto consumidor

El protocolo coordina los owners existentes: blueprint, readiness/rondas A–H, authority map, source lock, pack plan, spec/plan/tasks, execution state/events y failure lessons. NEW crea su baseline observado; EXISTING valida el cursor y registra sólo el delta. Nunca copiar el estado337 del mantenimiento ni heredar su READY como READY_TO_BUILD del negocio.

Planificar la franquicia completa y ejecutar incrementos sobre el mismo sistema. Elegir temprano el target de build, CI y staging; ampliar verificaciones según impacto y riesgo. El [perfil de referencia](markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md) sigue seleccionando116packs/1653archivos; el complemento de conexión es tooling separado. La referencia Windows/local y sus fixtures no autorizan otros runtimes, reglas de negocio o producción. Seguir los gates del protocolo antes de materializar producto.

## Paso 4 — Cuando quieras conectar tus cuentas

Usá los campos vacíos de `markdown_system/PROJECT_SECRETS_TEMPLATE.md` y guardá los valores en el gestor de secretos del destino; al chat y a los recibos sólo van referencias lógicas. Los adapters deben estar completos y probados localmente antes: código ausente no es una credencial pendiente. La infraestructura ARCA se prepara sin CUIT/certificados; su activación corresponde a este paso. Daybreak/libxml2 mantiene expediente y trigger de reapertura separados.

Para Cloud Run, el owner es `markdown_system/FRANCHISE_SERVERLESS_PACK_PLAN.md` con `GO-FINOPS-CORE`. El proyecto conserva sus verificaciones de reglas/país/corpus, permisos, callbacks, reconciliación, migraciones, rollout/rollback y aceptación del target antes de habilitar producción.

## Lo que NO es instantáneo (honesto)

- Responder las rondas A–H (negocio, reglas, supuestos).
- Obtener credenciales reales (cada consola tiene su trámite).

El agente te lleva de la mano en cada uno; no lo reemplaza, lo acelera.

## Partes condicionales

Sólo pueden diferirse si el blueprint las clasifica `OPTIONAL` o `NONE_WITH_REASON`. Si el negocio las clasifica `REQUIRED`, son gates y no se ocultan:

- **Fiscal / facturación** (ARCA u otra jurisdicción).
- **Homologación fiscal / notas de crédito / QR/PDF**.
- **Marketplaces** (Amazon/Google Merchant) y **ads** — solo si el canal se usa.
- **SMS** — solo si se elige ese canal.

El orden concreto lo determina el blueprint y las dependencias del journey; el agente debe entregar primero un vertical end-to-end verificable y luego ampliar sin duplicar owners.
