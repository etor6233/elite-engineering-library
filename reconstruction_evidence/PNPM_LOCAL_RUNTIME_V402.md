# pnpm local runtime — V402 / execution311

PROVEN_LOCAL para instalación offline restringida de la biblioteca. Pack
PNPM-ARTIFACT-SELECTION-GATE0.9.0,14archivos. Biblioteca204/2250; el perfil
de producto109/1481 y todos sus archivos permanecen exactos a310.
T2803 integral sigue hasta cerrar SCA/lint de esa composición. No READY global.

La nueva ruta usa Node24.20.0 por SHA y --jitless --no-addons. WebAssembly
está ausente tanto en main como en worker; dlopen devuelve ERR_DLOPEN_DISABLED.
Se conserva el bundle ZIP-disabled aff07c57e598640266b0ca1c795879db1bbdd2eafe71e96be34026358ba1c42b.
ZIP sigue retirado y13archivos físicos excluidos. Los seis blobs WASM retenidos
se fijaron contra los artefactos originales, pero no se ejecutan en este modo.
Esto es contención local, no compilación nativa independiente ni validación
de esos parsers para red. No se investigó Daybreak/libxml2.

Tres instalaciones reales — enterprise-web, Playwright y Lighthouse — pasan
offline, frozen, copy, store-integrity, sin scripts/hooks y con entorno reemplazado.
Sus entradas permanecen intactas y los logs registran cero descargas positivas.
Pasaron11negativos de receta, entorno, notices, hash externo, target ocupado y
replay. Ocho tests focales cubren bytes de notices, payload faltante/alterado/extra,
fallo de proceso y timeout con recibo/log conservados. Los ensayosV401 sin estas
restricciones se conservan como evidencia histórica; no fueron reetiquetados.

Se retuvieron467archivos npm originales por tamaño y SHA512 fijados, más8archives
de dependencias embebidas de OpenPGP. Sus42archivos de sourcesContent coinciden
byte por byte con los miembros originales. SCA OSV2.5.1:472identidades previas
y5nuevas,477npm/Python en total, cero hallazgos. No se reescaneó la base sin delta.
El catálogo entrega475entradas de paquetes/279textos únicos y conserva los47textos
y complementos anteriores. README/declaraciones/avisos oficiales quedan separados
por locator y SHA; una ausencia de copyright no se rellena con un autor inventado.

La admisión es para uso local privado del tool. La redistribución del pnpm
alterado queda fuera: excluirlo de ZIP de biblioteca/producto, adquirir/adaptar
en la máquina consumidora. El alcance de uso interno está documentado por las
autoridades oficiales GPL/LGPL y MPL, sin trasladarlo a una autorización para
redistribuir el bundle: https://www.gnu.org/licenses/gpl.en.html y
https://www.mozilla.org/en-US/MPL/2.0/FAQ/ . Sus textos originales siguen retenidos.

Se puede crear el mismo runtime445archivos desde la proyección ya verificada;
no hace falta reconstruir una extracción manual auxiliar. Dos materializadores
independientes reconstruyen exactamente los14archivos del pack. El executor
verifica receta/entradas/notices antes de ejecutar y conserva start/log/result;
falla cerrado ante intento repetido o revisión distinta. Recuperación: preservar
el intento y materializar un nuevo consumer en destino ausente, sin borrar datos.

G0–G8 y hashes completos están en PNPM_LOCAL_RUNTIME_V402.json. Nuevos5archivos
AUTHORED son glue/catálogos; contienen textos upstream VERBATIM identificados,
sin atribuir el código local a ninguna empresa. FAIL917–920 resueltos localmente.
Integrar esta ruta en preflight/build portable corresponde a T2808/T2810;
el antiguo probing de pnpm no se habilita por presencia en PATH.
