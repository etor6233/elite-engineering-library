# V269 — presupuesto real y lector compartido de respuestas BFF

Fecha: 2026-09-06. Mantenimiento de biblioteca, no implementación de una franquicia target. Claim: corrección de la frontera HTTP pública/protegida seleccionada; `REBUILD_VERIFIED / CONDITIONED`, sin promoción global.

## Resultado y procedencia

Se sustituyeron tres lecturas completas de response.text() por un solo readBackendResponse en el owner público ya importado por el cliente protegido. Se conserva el contrato público y las llamadas protegidas GET/POST, bearer server-only, no-store, redirect error, deadline de cinco segundos y ausencia de retry. El límite existente de 1.048.576 ahora se mide sobre bytes del cuerpo descomprimido antes de decodificar, no sobre caracteres UTF-16 después de descargar todo. Overflow y media type inválido cancelan la lectura. JSON/UTF-8 inválidos producen INVALID_RESPONSE sin payload; un código de error upstream debe cumplir la forma acotada del contrato, o se conserva el HTTP status con UPSTREAM_ERROR.

Código y tests son **AUTHORED**. No son código copiado de Microsoft, Google ni Meta. No se agregó dependencia, framework, migración, servicio, cuenta ni costo externo. Se mantienen las versiones y el lock anteriores; TypeScript sigue siendo sólo el frontend seleccionado, no una base obligatoria para otros proyectos.

Autoridades oficiales consultadas:

- [Microsoft HttpCompletionOption, Remarks](https://learn.microsoft.com/en-us/dotnet/api/system.net.http.httpcompletionoption?view=net-10.0): documenta la diferencia entre completar headers y leer contenido, y advierte que timeout y buffer del contenido requieren atención explícita. Es autoridad de método; aquí no se ejecuta ni traslada HttpClient de .NET.
- [WHATWG Streams](https://streams.spec.whatwg.org/): primitivas read/cancel/locks y lectura incremental. La revisión observada del estándar vivo fue b9ba9f49d95b4280be0dc2372377a006c3a91c18, actualizada 2026-08-18; ese identificador registra la consulta, no convierte su snapshot histórico en estándar vigente para futuros agentes.
- [Node Web Streams, ReadableStreamDefaultReader](https://nodejs.org/api/webstreams.html#class-readablestreamdefaultreader): cancel devuelve una promesa y read entrega chunks. La implementación libera el lock y no espera indefinidamente una cancelación fallida o estancada.

No se copió código de esas páginas ni se promovió un upstream nuevo. La cifra del presupuesto proviene del contrato local anterior, no de una recomendación empresarial universal.

## Regresión y ejecución

- Red sobre los clientes anteriores: **11 FAIL, 16 PASS, 1 SKIP**, exit 1. Se reprodujeron exceso multibyte aceptado, descarga drenada antes de rechazar, media inválida no cancelada, parsing inseguro y ambos métodos protegidos. FAIL-20260906-371 / LIB-FAIL-2105.
- Green focal: **27 PASS, 1 SKIP**. Después se añadieron tres pruebas HTTP nativas para no depender sólo del fetch simulado.
- Web completa: **70 PASS, 1 SKIP**, diez archivos de tests; 16 casos nuevos respecto de V263. Typecheck y build PASS en staging.
- Instalación frozen, typecheck, los mismos **70 PASS/1 SKIP** y build repetidos desde la reconstrucción limpia, sin reutilizar node_modules ni .next del staging.
- Tres pruebas sin mock de fetch usan un servidor HTTP loopback real: gzip pequeño que expande más allá del límite; chunked excesivo cuyo servidor no envía EOF y observa cierre del cliente; headers inmediatos con cuerpo incompleto que alcanza el deadline. No contactan proveedores, no levantan PostgreSQL ni sustituyen navegador/IdP.
- Los tests aislados cubren además frontera exacta, Unicode partido byte a byte, Content-Length engañoso, cancelación sin drenar/reintentar, cancelación estancada, media inválida sin leer, error de stream y errores JSON/UTF-8 sin contenido sensible.
- El test conectado de confirmación PostgreSQL permanece **SKIP explícito** en esta suite sin ELITE_CONFIRMATION_E2E. Su PASS V263 es histórico y no se presenta como repetido aquí.

Toolchain observado: Node 24.14.1, pnpm del proyecto 11.19.0, Vitest 4.1.11, Next 16.3.2 y TypeScript 7.0.2. pnpm global informó 11.25.0 fuera del proyecto, pero el packageManager fijado seleccionó 11.19.0 dentro; instalación frozen/ignore-scripts sin cambiar lock. No hay SCA global nuevo acreditado por este expediente.

## Fuente canónica y reproducción

Packs TS-GO-API-WEB-BRIDGE 0.5.2 y TS-OIDC-PORTAL-ADAPTER 0.2.3; planes integral, web y serverless actualizados. Se incorpora sólo un archivo de regresiones, no otro módulo runtime. Perfil integral reconstruido: **67 packs / 716 archivos**.

SHA-256 idénticos entre staging probado y reconstrucción nueva:

| Archivo | SHA-256 |
|---|---|
| src/platform/backend/public-client.ts | 68f3b45038d22330b0d81aa7f3bf2795b7c847e0faa68d672302b5b009f08463 |
| src/platform/backend/response-boundary.test.ts | 4ae89980e12a768e92503b6b72c58e54a45f5bdb07bc6df16ebaaca8d5071e03 |
| src/platform/backend/protected-client.ts | bea2c997aa68238374432dcaeaf318f231cb02685058b18c221c28f44c403aa1 |
| src/platform/backend/protected-client.test.ts | ad49ce7d0dc2310ea88465271435ebfa4e3a0928930d97580941a5b77e6fd3c2 |
| pnpm-lock.yaml, sin cambios | b3f49d7f7660ea858ff0291bea56e82f0cb75e28f3e7b4fdfa51c4449b9ab5af |

Materializar el perfil en una carpeta ausente e instalar con pnpm install --frozen-lockfile --ignore-scripts. Ejecutar pnpm typecheck, pnpm test y pnpm build. La nueva suite requiere poder abrir loopback efímero y tarda unos cinco segundos por su prueba de deadline. Staging conservado bajo el temporal del sistema: elite-v269-50c3dcbf6f384343ab7a6a9a3896a9b3; red.log, green.log, web-tests.log y roundtrip-*.log preservan resultados. No se borraron datos de usuario ni se dejaron servidores de prueba deliberadamente activos.

## Límites y siguiente cierre real

Se limita el payload retenido por el lector, no RSS total, objetos JSON parseados ni buffers internos de fetch, descompresor, proxy o red. Los schemas de dominio siguen en cada consumidor. Un fallo de lectura tras un POST puede tener efecto ya comprometido: no habilita retry ni éxito supuesto; el flujo existente debe consultar estado/receipt.

La sincronización de conteos tuvo una propuesta mecánica rechazada antes de editar (arrays aplanados/salida truncada); se corrigió usando diccionarios y siete líneas exactas. FAIL-20260906-372 / LIB-FAIL-2106 conserva esa lección, sin ocultar el intento fallido.

Este cambio no completa la pantalla de notificación WhatsApp, el worker ni la reconciliación de estados del proveedor. Tampoco acredita nuevos journeys de cotización/pedido, help/training/support release-bound, telemetría sin PII, cuentas reales, IdP, carga, seguridad ofensiva, recovery y deploy. Se mantiene el roadmap vigente sin añadir un subsistema paralelo ni porcentaje de cierre inventado.

## Verificación general

VERIFY_LIBRARY_PASS, proceso exit 0: **160 packs / 1.407 archivos materializables / 702 Markdown / 51 perfiles**. Procedencia: AUTHORED 1.169, ADAPTED 133, VERBATIM 105. Se actualizaron inventarios actuales y memoria (2.106 lecciones locales + 209 condiciones upstream), sin reescribir los expedientes históricos ni cerrar condiciones upstream pendientes. Las cabeceras de compatibilidad web/OIDC se alinearon con la combinación reconstruida; no se afirma compatibilidad nueva con el frontend histórico. Este PASS valida integridad/composición de biblioteca, no es Audit ejecutable integral ni certificación del producto.
