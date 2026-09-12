# V397 — entrega reproducible de los47textos retenidos de pnpm

2026-09-11. Mantenimiento de biblioteca; checkpoint262 validado.
**Cerrado:** empaquetado y entrega local verificable del conjunto completo
de47textos previamente retenidos. Las tres recetas actuales los incluyen;
141copias verificadas byte por byte,82tests y12negativos reales PASS.
No equivale a47dependencias admitidas ni al cierre del grafo completo.
Los48oráculos no cambian y el total permanece45/48.

## Origen y preservación

La colección V340 se revalidó contra MANIFEST.json SHA256
27d666e86d769c368439b1cc20e0fa72a42f9366973216b42ced5434f003fa9d.
Sus47entradas conservan151253bytes originales:
22notices que ya estaban en el payload seleccionado y25textos de investigación.
Los hashes individuales y longitudes siguen coincidiendo con el manifest.
Se conserva el draft anterior, sus scopes y las distinciones de publicación/
fuente descritas en PNPM_NOTICE_COVERAGE_V338.md,
PNPM_NESTED_LICENSE_SOURCES_V339.md y PNPM_YARN_UNDICI_NOTICES_V340.md.

El nuevo catálogo portable sustituye sólo los locators home absolutos por
Temp/stage/archivo. No modifica los textos: base64 conserva bytes, CRLF y
terminadores originales. Cada entrada declara origin, scope, source_locator,
payload_path cuando aplica, bytes y SHA256. El status es RETAINED_TEXTS_ONLY
y release_ready=false. Las copias duplicadas por versión se conservan
intencionalmente; no se cuentan como dependencias únicas.

Los textos incluyen notices originales de pnpm/node-gyp/Undici y otros
paquetes, GPL/LGPL para OpenPGP, MPL next-path, Artistic npm-lifecycle,
atribución spdx-exceptions, BlueOak y fuentes anidadas Noble/QRCode/Yarn/Undici.
Se conservan las reservas originales: texto de licencia no prueba equivalencia
del bundle, fuente preferida, publicación exacta ni cumplimiento de relinking.
El suplemento semver-utils V396, BlueOak V394 y QRCode V395 siguen separados
y byte-idénticos; no se hace regresión de esos cierres.

## Implementación y garantías

PNPM-ARTIFACT-SELECTION-GATE0.6.0 tiene7archivos.
retained-notice-catalog.json es ADAPTED: envelope local con textos originales
de terceros bajo sus términos retenidos. LicenseRef-Pnpm-Bundled-Licenses
describe esa colección; no concede una licencia nueva ni atribuye el código
de control a terceros. El selector, su política y17tests no cambian.
Cambian guía, planner y routing tests; no hay nueva dependencia runtime.

Recipev5 entrega PNPM-RETAINED-NOTICES.md junto a los tres suplementos previos.
El loader fija SHA256 del catálogo, schema exacto,47entradas, paths seguros
y únicos, estados sin promoción, codificación base64 estricta, hashes y tamaños.
Los22notices originales deben coincidir además con el payload exacto.
Límites: catálogo524288bytes, texto individual65536, colección original200000,
documento entregado262144. No red ni ejecución de paquetes.
La publicación mantiene el destino ausente y la limpieza del stage fallido;
verify rechaza archivos extra/faltantes y reconstruye el documento esperado.

82tests:17selection+65routing;15regresiones nuevas incluyen entrega,
omisiones, cambio de permiso, falsificación conjunta del texto/receipt,
mutación de notice original, fallo de escritura, catálogo con path duplicado/
traversal, status/origin inflado, encoding/contenido/hash/count/fields alterados.
Se conservan todos los tests anteriores. Siete archivos reconstruidos por el
materializer canónico, idénticos al candidato probado.

## Recorrido real

Composición fresca del perfil enterprise web:6packs/126archivos.
Para enterprise-web, Playwright y Lighthouse se preparan recetas con el
Node, pnpm proyectado y receipt fijados. En cada documento se localiza cada
entrada y se extrae por su longitud original, independientemente del generador:
los47SHA256/textos coinciden,141entregas en total.
Cuatro ataques de integridad por receta se rechazan:
permiso omitido, scope falseado, colección ausente y documento/receipt
falseados y rehasheados juntos. Restaurar bytes devuelve PASS.
Los442archivos originales de pnpm y su receipt, y los nueve suplementos previos
entre las3recetas, permanecen byte-idénticos. executed/runtime_admitted/
redistribution_admitted permanecen false.

## Fallos, límites y continuidad

FAIL793: los47textos verificados seguían fuera de las recetas. Ahora existe
una entrega materializada y comprobable, con alcance conservado.
Un comando de arranque apuntó inicialmente al directorio candidato antes
de crearlo; fue rechazado sin ejecutar, se creó desde la raíz y se siguió
en el destino observado. No se infiere evidencia de la invocación rechazada.

Esto no cierra las obligaciones de fuente/relinking OpenPGP, MPL/Artistic,
la procedencia publicada completa de Yarn/pnpm, los nativos diferidos,
la seguridad del grafo ni TEST-02/03/07. Daybreak/V386 permanece diferido.
El siguiente trabajo es la cobertura de permisos/source/publicación que
todavía carece de evidencia; no repetir la entrega47 ya cerrada.

Rollback: snapshots del pack0.5.0/plan y expedientes en before; recetas
anteriores y drafts intactos. No se reemplaza un runtime ni se distribuye
el payload fuera de las carpetas locales de evidencia.

## Receipts

Stage: Temp/elite-v397-8b73338920d5; pointer elite-v397-current.txt.
Sólo el reporte y los bloques canónicos son portables; receipts y snapshots
del mantenimiento quedan locales. No se añade código de paquetes al producto.

| Evidencia relativa al stage | SHA256 |
|---|---|
| catalog-baseline.json | 46d6ac0ae7276aa0c660d38bcf42e3b2b1d4ec2db8c167365d8294f670096235 |
| source-parity.json | 2afbeca91bf375277e3242d5ec9221e573db88fd85679f3c13ecb55706025b36 |
| rebuild.log | be5407730fd15e155325fa415054e5c347db142dc8183472f359d486816d20c4 |
| rebuilt-tests.log | cad8691d79a7d1587d1e2d2cf17e3793ca637639f8abb07604305008298fbc78 |
| current-web.log | a8130b16e4e5bdb3129feb17191dfb7e9bc2a6a22fab09c6e75dd0f4c87bd800 |
| real-delivery.json | f3a05c71dcda31fe17fbfebf2f2af0f061af443c00eeda882a2d3497ccb4bc51 |
| projection-before-after.json | 4b56bb3ee5a5bfa16b7ac0832b67866ff5d5348e7357be81e1e3d200eca9f27b |

Preflight263 pendiente tras checkpoint; aún no se atribuye resultado global a esta revisión.

## Cierre comprobado — checkpoint264

Preflight263 pasó164pasos ejecutados y56composiciones; inventario165packs,
1563archivos materializables,845Markdown. Procedencia1307AUTHORED/149ADAPTED/
107VERBATIM. Sólo Docker está ausente; no se confunde disponibilidad global
BLOCKED con el PASS de las comprobaciones ejecutadas ni con admisión productiva.
La franquicia68/797 y referenciaHTTP69/805 no cambian.
El Preflight general ejecuta sus instalaciones offline congeladas previas;
las tres nuevas recetas verificadas no ejecutan pnpm.

La entrega del conjunto47 queda cerrada, además de los suplementos BlueOak,
QRCode y semver ya cerrados. Las7fuentes canónicas son idénticas a las
reconstrucciones probadas; selector, política y17tests permanecen byte-idénticos.
Los48oráculos/acciones/estados del contrato permanecen iguales:45passed y
TEST02/03/07blocked.

FAIL794 / STALE_ENTRYPOINT_SUMMARIES: README aún declaraba1561files y el
roadmap abría en V392/checkpoint248, pese a cierres posteriores conservados
al final. Se corrigieron el conteo vigente y el encabezado a V397/checkpoint264,
preservando V392 como corte previo y los snapshots originales. No es un cambio
de estado de admisión. Esto evita reanudar desde una entrada desactualizada.

| Evidencia final relativa al stage | SHA256 |
|---|---|
| preflight263.json | a77449be7e71f3c1e4c37cad228aa03e1c74ab7a56d5c316b80fa0a682fb95ba |
| preflight263-process.json | ad38966244012212879fd2320730218646b676081c79032c2731a72dc79b3609 |
| preflight263.log | 22b2f68bb10f9ee261bd87eaa933c7be76e140966465b2faf642bcd0d5dfc491 |
| delta-audit.json | 91f76ade218da1a6691666a2c188f7a3702cdce91acb044fc5a7f8c1913f73fe |
| current-entrypoint-correction.json | 9d030007ef9e1db02acfcd79becc7825f17edfa69e2659cb06de322d10f01a64 |
