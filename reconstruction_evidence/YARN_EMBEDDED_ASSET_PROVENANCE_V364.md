# V364 — procedencia del cargador Yarn incorporado en pnpm

Fecha: 2026-09-09. Mantenimiento de biblioteca; análisis de bytes primarios ya
capturados. Entrada: checkpoint172 validado con resume/plan y hashes de owners.
No se adquirieron, instalaron ni ejecutaron fuentes de paquetes en esta revisión.

## Resultado demostrado

El literal Base64 completo del cargador ESM de @yarnpkg/pnp4.1.7 contenido en el
pnpm11.25.0 seleccionado coincide con el archivo de la revisión Yarn
4fe4d4bf45a13dca90181c5b7fee61376aa21794. Sus 16500caracteres codifican
12375bytes Brotli y producen **72463bytes descomprimidos idénticos**.
Se descomprimió como datos con límite de1MiB y se analizó su sintaxis; no se
ejecutó el cargador. La extracción independiente usa coincidencia léxica y
verifica el rango exacto del bundle, además de los hashes de ambas formas.

Los blobs de esm-loader/built-loader.js y loader/node-options.js se verificaron
contra el árbol sources3a129c0862c9942c1f7c68f188ae3ad81f357f18 capturado en V340.
Ese SHA coincide con la entrada de directorio consultada al commit fijado.
No se presenta la respuesta API como una firma criptográfica del commit.

Cinco funciones del source node-options coinciden en AST con el bundle:
getOptionValue, parseOptions, parseArgv, getNodeOptionsEnvArgv y
ParseNodeOptionsEnvVar. Mapas bundle→source explícitos: argv2→argv, arg_1→arg,
errors2→errors, index2→index y c3→c, sólo en sus funciones correspondientes.
Se excluyen nombres de propiedades no computadas. Las únicas normalizaciones
adicionales son concatenación de dos literales string, expansión de propiedad
shorthand y bloque if vacío a sentencia vacía. Se conservan operadores, valores,
condiciones, orden y raw de templates; offsets y grafía de literales se omiten.

Esto NO comprueba la identidad de los bindings externos arg/options, la
equivalencia del módulo completo, su build, reachability o comportamiento real.
Los nueve controles iniciales rechazan cambios de texto, tipo, cuerpo, propiedad,
predicado, incremento, bytes comprimidos/descomprimidos y sintaxis inválida.
Una comprobación independiente rechaza cinco mutaciones del source real:
lookup de opción, modo permissive, error de escape, incremento y condición final.
Son14controles del análisis, no14nuevos controles de producto ni14filas del contrato.

## Corrección de alcance y licencias

V340 sólo demostraba versión/tag/notices; la presencia de código no estaba
relacionada con esos sources. V364 demuestra el asset completo indicado y cinco
unidades del parser. El comentario MIT de Joyent se conserva en el draft47 y
pertenece al source de esas funciones. Su cadena exacta no aparece en pnpm.mjs;
las cadenas Node/Blake sí aparecen en el bundle completo. Esta observación no
dictamina incumplimiento ni afirma cobertura del conjunto distribuido.

Los imports estáticos del cargador decodificado son fs,url,path,crypto,os,module
y assert. No equivalen a un inventario transitivo: el código integrado y las
cargas dinámicas deben seguir incluyéndose en la calificación de seguridad.
No se elimina ninguna de las473identidades conservadoras ni se inventa versión
de una dependencia incorporada a partir de esos imports.

FAIL575 continúa abierto: el gitHead npm5761b03feb2146da8ce8cefeba6482f3cb6edb79
no queda sustituido por el tag. El endpoint de attestation fue rechazado por la
herramienta web y no se reintentó mediante otro transporte. La coincidencia del
asset no autentica toda la publicación npm. FAIL532 del pnpm nativo original
tampoco se cierra: V333–336 ya habían provisto la proyección442 y su routing;
no se cuentan como un nuevo trabajo. semver-utils/chownr y obligaciones de
source/relinking LGPL continúan por sus owners de V338–341.

FAIL652/653 registran errores de rutas de investigación recuperados; FAIL654
registra acceso web bloqueado sin inferir ausencia de fuentes ni permisos.
La investigación siguiente debe usar archivos descubiertos, conservar receipts
y resolver la publicación completa o una alternativa admitida por su propio gate.

## Evidencia y continuidad

Stage: %TEMP%/elite-v364-db32bcda0ab24dabacebf983f9e3dd08.
Parser existente Acorn8.14.0 de Next16.3.2, licencia MIT y SHA verificados;
Node24.20.0 existente. Harness de análisis local AUTHORED, sin nuevo runtime.
Los163Markdown de implementation_packs,442archivos del payload seleccionado y
47copias de notices conservan sus hashes; no hay nueva versión ni promoción.
EVID-74 es parcial para TEST03/07.43/48filas siguen PASS,5BLOCKED; eso no mide
esfuerzo restante. TEST02/03/05/07/09 y las10macro-tareas continúan abiertas.
No se infiere respuesta sobre corpus autorizado, target operativo o release final.
La evidencia del ZIP148 conserva su alcance histórico y no acepta un ZIP nuevo.

La verificación estructural se ejecuta después del checkpoint173; no se repiten
ni redatan los154checks del Preflight171 para un cambio sin código canónico.

| Recibo relativo al stage | SHA-256 |
|---|---|
| qualify_yarn.cjs | `7a6ee579d4f1db394a1e653b4980b9d291f6b594a712fd8af1faeb9ef6072af3` |
| independent_yarn.cjs | `d71f7113a04acaf7c4f03c7a3987ba0bc3ccf85e122f75408460b40f61425db2` |
| yarn-qualification.json | `25c7d1d04eba8384741d5c92586287480f02e36bd780f2cd797b11a77db3e53a` |
| independent-yarn.json | `6d36449af7b56e991f3fd66300c85cdfce14c2ae18a8fe5c125b64e3071dc720` |
| qualification173.json | `918f85eead063bb6380fa9ea4123d588a73f9b4fa0a781407c6c561b371a9dbb` |
| yarn-bundle-region.evidence | `33510fbab8397f66c0d1ea83316457b06803678806596ba4dcef68924429c385` |
| yarn-esm-loader-decoded.evidence | `d8fdc78e08361fd1ee851003498918b0b68fd5fd131226970201e2100000afdc` |
| resume172.log | `f978f285ae9a316f5a0ea6fce12373a7ae70823606874e33c17d449b3465c0b9` |
| plan172.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |

FAIL655: el generador inicial supuso12373bytes comprimidos y su assertion abortó antes de checkpoint173. Recibos independientes miden12375; texto canónico corregido conservando before/failure.

## Corrección de inventario antes del replay174

VERIFY_LIBRARY173 rechazó el encabezado vigente del preflight (801 frente a802Markdown). FAIL656 conserva el fallo; se corrigió únicamente el cursor actual aV364/802. El verifier sigue intacto. Checkpoint174 y replay estructural174 preceden el cierre.

Log fallido173 SHA256: `1a2526e5b9a98f25c00fba24ed6e41acf256d5bfbb9335ed687111cfb4299a7f`.

## Cierre175

VERIFY_LIBRARY174 PASS:162packs/1461archivos materializables/802Markdown/53perfiles. Resume/plan del checkpoint175 se ejecutan después de registrar el cierre. Los163Markdown de implementación y las48definiciones/requisitos/estados permanecen intactos.43PASS/5BLOCKED; no porcentaje global ni entrega final. FAIL656 queda REGRESSION_PROVEN por replay; FAIL652/653/655 recuperados y FAIL654 de acceso permanece. Asset72463bytes,5funciones y14controles constituyen evidencia parcial EVID-74, no readmisión de pnpm/Yarn completo.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| verify-library174.log | `ee3f73f61ccdd6ef593b26cd70ae8166e3d3d3250d722a7a18f5f425f5c8932b` |
| closure174.json | `dab63bbf21fb726d6ed492c8fdeafac299f167c2e5317ea8932a7b82d6d3c845` |
| checkpoint173.log | `8aee304d37bc40b6e55f081da65b60f4e3db12423f4275c8f1eedb7e22c19d30` |
| resume173.log | `b4ff3bdaa3c7484d8a028f656649045e6d356c72e838be807ebfedb3f5578e4f` |
| plan173.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
| checkpoint174.log | `b5eca84943637d14eb9126a125b59a82110a011a01ccb38b6ba94977d3c35e29` |
