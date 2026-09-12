# V395 — aviso QRCode completo y entrega verificable

2026-09-11. Mantenimiento de biblioteca; entrada checkpoint254 validada.
**Resultado:** cerrado el pendiente de generación y entrega local del aviso
MIT de QRCode en las3recetas pnpm.59tests PASS,12negativos reales rechazados,
6fuentes reconstruidas idénticas. La biblioteca conserva45/48controles:
integración integral, seguridad restante y release todavía pendientes.

## Del hallazgo al archivo entregado

V338 detectó la atribución MIT del componente QRCode dentro de qrcode-terminal0.12.0.
El baseline395 confirma diez módulos vendorizados en el bundle y ausencia del
nombre del autor; las tres recetas394 tampoco entregaban QRCODE-NOTICE.md.
La nueva versión0.4.0 lo genera y liga su SHA-256 a la receta v3.

El archivo contiene el encabezado original completo: copyright2009 Kazuhiko Arase,
declaración MIT, aviso de marca y modificación para Node. Se suma el texto MIT
completo publicado por OSI, con la identidad del copyright tomada del encabezado.
Se declara expresamente que es un suplemento ensamblado localmente. No se presenta
como LICENSE original del archivo pnpm ni se reemplaza la licencia Apache del
paquete contenedor. El aviso BlueOak anterior permanece byte-idéntico.

## Fuentes y alcance

- Fuente fijada: https://github.com/gtanner/qrcode-terminal/tree/90f66cf5c6b10bcb4358df96a9580f9eb383307b/vendor/QRCode
- Declaración original: vendor/QRCode/index.js, Git blob10eb8eb0a06aa50d0d3c508f886a7728f52e5d98.
- Texto de permiso: https://opensource.org/license/mit (consultado2026-09-11).
- Header SHA256f265b9225bb2a1a209d60f81be487f23257c3f8637498282212d04af7822c2b3.
- Source index SHA2567377be90fc61a40268acf7f30d5bd89c2fca99c57ef5391623de8c151b8da7df.

Se conservaron los10archivos del directorio fijado y se verificó cada Git blob
contra el árbol no truncado. El index coincide byte por byte con la evidencia338.
Nueve archivos auxiliares no repiten individualmente el encabezado del index;
el suplemento preserva la declaración del componente, sin inventar headers.
Las10regiones correspondientes del bundle tienen offsets y SHA-256 propios.
La comprobación distingue identidad de fuente, identidad del bundle y aplicabilidad
del aviso; no afirma igualdad AST, bindings completos ni build reproducible.

## Implementación y verificación

El pack mantiene6archivos. Cambian guía, preparador y pruebas; el selector,
su política, sus17pruebas, hashes de consumidores y Node permanecen intactos.
plan_install.py usa bloque ADAPTED con lógica AUTHORED y texto de aviso VERBATIM
declarado; licencia LicenseRef-Workspace-Owner AND MIT. No incorpora algoritmo
QRCode ejecutable ni una dependencia nueva.

Cada preparación verifica el bundle exacto y las diez regiones antes de generar
el suplemento. verify vuelve a calcular el aviso y rechaza ausencia, alteraciones
o falsificación conjunta de aviso/recibo. Ambas copias de notices quedan fuera
del payload442, que conserva sus22notices originales y receipt.

-59tests en candidato y reconstrucción:17selección +42routing.
-Seis archivos reconstruidos byte-idénticos al candidato; UTF-8/LF.
-Composición canónica fresca:6packs/126archivos.
-Tres recetas reales con el pnpm proyectado, Node y acquisition receipt fijados;
  no mocks en este recorrido y ningún pnpm ejecutado.
-Cada receta contiene encabezado618bytes exacto y los tres párrafos completos
  de permiso/condición/garantía MIT, además de los10paths y hashes.
-Cuatro negativos reales por receta: autor alterado, permiso sustituido por
  etiqueta MIT, archivo eliminado y aviso/recibo falsificados y rehasheados.
-Tras restaurar los bytes originales, las3verificaciones vuelven a pasar.
-BlueOak intacto,442archivos originales y receipt idénticos a394.

## Fallos y recuperación

FAIL719: consulta errónea a pnpm.cjs; inventario confirma pnpm.mjs. No se infiere
evidencia de un archivo ausente.
FAIL783: aviso identificado pero no entregado; corregido con artefacto usable
y verificador, no sólo un documento de intención.
FAIL784: materializer rechazó provenance MIXED en el pack pendiente. Se preservan
pack/log fallidos; schema de bloques acepta AUTHORED/ADAPTED/VERBATIM. El sucesor
ADAPTED declara exactamente sus dos procedencias y AND MIT. No se debilitó
el materializer ni se ocultó la atribución para pasar el gate.

## Límites, rollback y próximo pendiente

Este cierre demuestra entrega local de atribución/permiso. El planner no envía
el payload a terceros; cualquier entrega posterior debe llevar los notices.
No cierra licencia global pnpm, source/relinking LGPL, MPL/Artistic, semver-utils
ni publicación/correspondencia completa. runtime_admitted=false,
redistribution_admitted=false y executed=false en las3recetas. V386/libxml2
continúa diferido y no se reintentó por otra vía.
Rollback conserva0.3.0 y sus recetas; no se migra ni reemplaza un destino existente.
El siguiente alcance es completar la licencia exacta semver-utils o dejar su
acceso bloqueado con evidencia, y avanzar sobre las obligaciones restantes.

## Receipts

Raw evidencia en Temp/elite-v395-97bd97d95fc147689c0eb071b38190a6;
resolver elite-v395-current.txt. Sólo este reporte y los bloques canónicos son
portables; no se agregan binarios ni fuentes QRCode ejecutables al producto.

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline-missing-notice.json | 0d8d8ecbf6466bc5a5927eef076d1a7ef12338f9fa711aaeb36a217a26f042d3 |
| authority-receipts.json | e9b191af0b200903e5e940bb558d97c1544e7a3e05461c088b934a3baa9f655b |
| qrcode-bindings.json | 2c0befee20b26f6f6b8ec2381bc4b67486c2a9590891b0ba50034e6bd239a995 |
| qrcode-fixed-tree.json | 6fce3d853b993751fccd664e6c917fcfa783b2cdae83cdccb0f1bf455a2fe4c6 |
| source-header-inspection.json | 0514037cdaa2230bd2478020fee7e29f290b5f5d5bd38589d659c70d22546f43 |
| qrcode-vendor-header.txt | f265b9225bb2a1a209d60f81be487f23257c3f8637498282212d04af7822c2b3 |
| QRCode-MIT-LICENSE.txt | f94ed090c2a28dd590825d9a9f4891a938c52cb052f8be50a65cf6be38ab5780 |
| source-parity.json | 393919e06dd9867bb0b0defee71be6dab073b106641fc94259e5f30c62a43aff |
| rebuilt-tests.log | c9b92dbe3ca850547f5917f24a5d203083b44a738c21d7ee5c263297ea449ea7 |
| real-delivery-results.json | a40e05b9b9c41e6c709f43deed1c26ecac684cda72b835ebbd96d782ad4a0ae7 |
| projection-before-after.json | 4b56bb3ee5a5bfa16b7ac0832b67866ff5d5348e7357be81e1e3d200eca9f27b |
| rebuild.log | 26e124367b729bfaa982339b0f91fc5e82f1b8b1d67336409e44062afb3f6caf |
| pack-before-provenance-repair.md | 5c0aac0a08f56578946ac51f8086a8f4cd041c7d9b87c6542321ef7618fc485c |

Preflight255 rechazado por FAIL785: conteo de procedencia desactualizado en THIRD_PARTY_NOTICES. Corregido a1306AUTHORED/148ADAPTED/107VERBATIM=1561; original y log preservados. El sucesor Preflight256 completo pasa164pasos y56perfiles; disponibilidad global BLOCKED sólo por Docker ausente.

## Cierre de verificación general

Preflight256:164pasos ejecutados PASS;42regresiones routing más17selección,
comparación de los3consumers actuales y56perfiles. Inventario165packs,
1561archivos materializables,843Markdown. Procedencia1306AUTHORED/148ADAPTED/
107VERBATIM. Docker ausente mantiene el estado global BLOCKED; no se confunde
este resultado con admisión productiva.

FAIL783 y784 quedan reparados; FAIL785 se repara en THIRD_PARTY_NOTICES y pasa
el verificador sin cambiar su regla. Los48oráculos, acciones y estados siguen
idénticos. EVID-101 aporta sólo atribución/permiso y su entrega verificable.

La revisión siguiente se inició sin repetir el endpoint bloqueado: el gitHead
9f3dbc8d22ab93a0d56568b95b79257d58a58545 de semver-utils1.1.4 en el repositorio
GitHub previo del mantenedor devuelve404. Se conserva URL y resultado; no se
infiere que no exista fuente o licencia en otra autoridad. Continúa
RESEARCH_INCOMPLETE. El transporte admitido exige sidecar real, por lo que no
se reemplazó APACHEv2 por texto genérico. El siguiente paso es calificar una
adquisición del archivo oficial limitada a evidencia, previa a admisión,
o encontrar texto primario exacto disponible.

| Recibo final | SHA-256 |
|---|---|
| preflight255-process.json | 606e000f1884ccd1ad48a1bb37708d277b4f2019d6d711e66b9e1c254d5f8638 |
| preflight255.log | ec6548b18f0484221138d1349b691772c0ffb5935f0f61a3ea35e4469632f7eb |
| preflight256.json | 7073b07f700549bd3a4183c3b263ea9eb23e398a34981894effd175bf0e59c9c |
| preflight256-process.json | 192ad23bd010a3e54ff37d29a16e2719e5a1fc5cb7dc8a32d1aa0363f3438e2c |
| preflight256.log | 729f7deb27ef36af58d322bbf79781cda1fa76b835daa216c7dbdc9abeb03ba0 |
| acceptance-preservation.json | 303cb51b3b6edb398713047cc997dc407b88750caa4a0e1659f724f1cb333c96 |
| final-source-delta.json | f3782688303aa98ce399ab0330d39978c986813bd52654ed8129fd21e144edb9 |
| semver-next-source-observation.json | 7f89b67723b5efd76409a767e12be75bba0b42192ac5681cc03b52f76007c742 |

## Rectificación de alcance de ejecución (FAIL786)

Las3recetas395 se prepararon/verificaron sin ejecutar pnpm, al igual que las3de394.
El Preflight general252/256 sí usa el pnpm fijado en sus gates habituales:
instalaciones frozen/offline para Playwright (3paquetes) y Lighthouse (119).
Los logs y VERIFY_EXECUTABLE_LIBRARY.ps1 lo muestran. El informe394 había
generalizado indebidamente la no-ejecución a todo Preflight: se preservó su
original y se corrigió la fuente canónica. Los flags executed=false pertenecen
a los receipts del preparador, no a todas las operaciones de mantenimiento.
Esto no modifica código, fuente instalada ni admisión runtime/redistribución.

| Evidencia de corrección | SHA-256 |
|---|---|
| preflight-runtime-scope.json | b01f08538f73eca0e9c24d268ba0f8cbd9e054227b35bd6b7333a1d7891cb8c1 |
| V394-report-before-runtime-scope-correction.md | e51ed7fef9798428e69a8eb60dbbf69283aa26b173dd430808ef98d23afe595c |
