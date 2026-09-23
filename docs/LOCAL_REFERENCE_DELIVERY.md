# Entrega local de referencia

Scope: `LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES`, Windows x64. Estos comandos
son glue AUTHORED sobre los toolchains, contratos y owners ya fijados. No son
código de Microsoft, Google, xAI ni otra empresa. El artefacto incluye material
sintético público de sesión/Server Actions y rechaza activación productiva.

## Construcción

Materializar el perfil completo y los owners `PNPM-ARTIFACT-SELECTION-GATE` y
`WINDOWS-REFERENCE-TELEMETRY-RUNTIME` seleccionados por ese plan. Crear un JSON
`ruta relativa -> SHA-256` de las salidas exactas del compositor. El hash externo
de ese JSON es la autorización de la revisión; no se confía en un hash declarado
por el propio artefacto. Conservarlo con el receipt de composición.

`python ci/build_local_reference.py --help` enumera las rutas requeridas:
source, target ausente, inventory y su hash, Go, Node, caché de módulos y
install-inputs con su hash. `install-inputs` tiene exactamente projection,
contained, store, cache, acquisition_receipt y acquisition_sha256 del receipt
admitido de pnpm. Usa `pnpm_artifact_selection/local_runtime.py`: instalación
offline/frozen/copy, sin scripts, bajo Node jitless/no-addons. El payload local
adaptado de pnpm se adquiere/verifica fuera del artefacto y nunca se redistribuye
con él. Una caché ausente falla; no se descarga ni se promueve otra versión.

`--workspace` es una ruta local estable y ausente para el build. Next16.3.4 usa
rutas absolutas como identificadores RSC: los dos builds independientes usan esa
misma ruta lógica, cada vez con archivos nuevos e instalación nueva. Archivar el
workspace propio terminado antes de la segunda ejecución; el builder nunca lo
borra ni reutiliza un checkout/caché de compilación. El recibo conserva la ruta.
La cache de preview de Next recibe valores públicos de fixture fijados, y se
comprueba que el framework los consumió. El claim byte-idéntico es para esta
receta y toolchains del host, no para cualquier ruta o plataforma arbitraria.

Go se compila con trimpath, sin VCS ni build ID; se incluyen todos los módulos
locales de la composición. Next usa webpack, su salida standalone y el hash del
inventario como build ID. Sólo se normalizan sus localizadores de máquina en
server.js y required-server-files.json; los hashes previos quedan en el receipt.
Los arrays de archivos de los traces NFT se ordenan como metadata de packaging.
No se reescriben bundles de aplicación. El framework conserva enlaces a paquetes
en Windows: el artefacto almacena únicamente archivos normales y el mapa de
enlaces relativos; cada activación recrea junctions dentro de su copia privada.

El artefacto contiene API, Node fijado y su licencia original, frontend,
migraciones, fuentes correspondientes y los notices de los owners y dependencias
observadas. La igualdad de dos artefactos debe comprobarse sobre el inventario
completo de archivos. La firma y aceptación portable pertenecen a TEST07/T2810.

## Base vacía y reanudación

Preparar un PostgreSQL local admitido y una base de fixtures propia cuyo nombre
cumpla `elite_payment_connected_<32 hex>`. El runner exige loopback, usuario
postgres sin contraseña en el perfil de fixture, puerto explícito y sslmode
disable. Esto no prescribe autenticación para un despliegue productivo.

`python ci/migrate_local_reference.py --help` acepta el directorio migrations,
su inventario/hash, DSN local, psql/hash y directorio de receipt ausente. Aplica
las migraciones fijadas en orden con lock de sesión, transacción por migración y
ledger atómico de nombre/ordinal/hash. Reanudar no repite migraciones; un hash
distinto o una base preexistente sin baseline registrado falla cerrado. Una
migración fallida no autoriza tráfico y no ejecuta down migrations.

## Arranque y rollback

`ci/local_identity_fixture.cjs` es exclusivamente un emisor OIDC sintético
loopback para el ensayo. No sirve como IdP productivo. El runtime JSON local
contiene exactamente database_url, issuer, api_port y web_port. No acepta keys,
env arbitrario ni endpoints remotos. Pago/WhatsApp/IA/ARCA live no se habilitan.

`python ci/run_local_reference.py --help` recibe artefacto, inventario/hash,
runtime/hash, receipt/hash de migraciones y state-directory. Comprobará que la
base y las migraciones corresponden al artefacto. Arranca por un período finito
de hasta 480 segundos; conserva current.json y journals de activación. El owner
local exclusivo inicia API y frontend bajo el launcher de Job Objects existente.
El API adopta el mismo evento nativo ya probado por el worker de devolución.
Next conserva su propio cierre SIGTERM (exit 143); el wrapper termina con 0
sólo tras los cierres esperados y el supervisor comprueba árbol vacío.

`ci/local_release.py` expone Deployment.activate/recover/close para un operador
local. Verifica la revisión antes de detener la anterior; si falla la salud de
la candidata, reinicia el artefacto previo y conserva su puntero. Una activación
satisfactoria sustituye el puntero atómicamente. Los journals fallidos se
preservan. Recovery vuelve al último puntero confirmado; no adivina aceptación
desde un proceso vivo. Los handlers y PostgreSQL conservan idempotencia y datos.

Los probes locales prueban arranque/routing; los journeys completos conservan
sus suites por owner. La operación con carga/alertas/retención se cierra en
T2809. No se infiere HA, DR, SCM, aislamiento frente a otro usuario hostil ni
capacidad productiva por un Job Object o un health check.

## Otras plataformas

El blueprint selecciona esta referencia local. Contenedores/cloud/Kubernetes
siguen opt-in por target; el Dockerfile no aporta defaults de tags móviles y
copia el contexto para conservar los módulos locales. Sus imágenes, runtime,
TLS, permisos, registry, proveedor y ejecución requieren admisión propia. No
se etiqueta una ruta sin ejecutar como CONDITIONED sólo por credenciales.

Fuentes de configuración consultadas: [Next standalone](https://nextjs.org/docs/app/api-reference/config/next-config-js/output)
y [build ID](https://nextjs.org/docs/app/api-reference/config/next-config-js/generateBuildId).
La implementación consumida sigue fijada en Next16.3.4 y los locks vigentes;
consultar documentación actual no actualiza una dependencia automáticamente.

Firma local: docs/SIGNED_LOCAL_REFERENCE.md conecta este builder con el gate firmado. Cada invocación es limpia y archiva su workspace; la publicación verifica el recibo y el inventario real.
